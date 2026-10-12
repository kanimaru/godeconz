package ws

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/coder/websocket"

	"github.com/kanimaru/godeconz/rest"
)

// Real-world shaped messages: lowercase resource names, partial state on "changed".
const (
	tempEvent  = `{"e":"changed","id":"1","r":"sensors","state":{"lastupdated":"2024-05-01T10:00:01.123","temperature":2150},"t":"event","uniqueid":"00:15:8d:00:04:aa:df:08-01-0402"}`
	lightEvent = `{"e":"changed","id":"4","r":"lights","state":{"on":false},"t":"event","uniqueid":"00:0d:6f:00:10:65:8a:6e-01"}`
	sceneEvent = `{"e":"scene-called","gid":"1","r":"scenes","scid":"2","t":"event"}`
	attrEvent  = `{"attr":{"id":"4","lastannounced":null,"lastseen":"2024-05-01T10:01Z","name":"Workfloor-1","type":"Dimmable light"},"e":"changed","id":"4","r":"lights","t":"event"}`
)

// server sends the given messages on each connection, then closes it, so clients reconnect.
func server(t *testing.T, perConnection ...[]string) (url string, connections *atomic.Int32) {
	t.Helper()
	connections = &atomic.Int32{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Accept(w, r, nil)
		if err != nil {
			return
		}
		defer conn.CloseNow()
		n := int(connections.Add(1)) - 1
		if n < len(perConnection) {
			for _, msg := range perConnection[n] {
				if err := conn.Write(r.Context(), websocket.MessageText, []byte(msg)); err != nil {
					return
				}
			}
			conn.Close(websocket.StatusNormalClosure, "bye")
			return
		}
		// Later connections stay open until the client goes away.
		ctx := conn.CloseRead(r.Context())
		<-ctx.Done()
	}))
	t.Cleanup(srv.Close)
	return "ws" + strings.TrimPrefix(srv.URL, "http"), connections
}

func TestRunDeliversEventsAndReconnects(t *testing.T) {
	url, connections := server(t, []string{tempEvent, lightEvent}, []string{sceneEvent, attrEvent})

	var mu sync.Mutex
	var events []Event
	var onConnectCalls atomic.Int32
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client := New(url,
		WithBackoff(10*time.Millisecond, 50*time.Millisecond),
		WithOnConnect(func(context.Context) { onConnectCalls.Add(1) }),
	)
	done := make(chan error, 1)
	go func() {
		done <- client.Run(ctx, func(e Event) {
			mu.Lock()
			events = append(events, e)
			if len(events) == 4 {
				defer func() { go waitThenCancel(connections, cancel) }()
			}
			mu.Unlock()
		})
	}()

	err := <-done
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Run returned %v, want context.Canceled", err)
	}
	mu.Lock()
	defer mu.Unlock()
	if len(events) != 4 {
		t.Fatalf("got %d events, want 4", len(events))
	}
	if onConnectCalls.Load() < 2 {
		t.Errorf("onConnect called %d times, want at least 2", onConnectCalls.Load())
	}

	temp := events[0]
	if temp.Kind != EventChanged || temp.Resource != ResourceSensors || temp.ID != "1" || string(temp.Raw) != tempEvent {
		t.Errorf("event 0 = %+v", temp)
	}
	state, err := temp.SensorState()
	if err != nil {
		t.Fatal(err)
	}
	if c, ok := state.TemperatureCelsius(); !ok || c != 21.5 || state.Humidity != nil {
		t.Errorf("sensor state = %+v", state)
	}

	light, err := events[1].LightState()
	if err != nil || light.On == nil || *light.On || light.Bri != nil {
		t.Errorf("light state = %+v, %v", light, err)
	}

	if events[2].Kind != EventSceneCalled || events[2].GroupID != "1" || events[2].SceneID != "2" {
		t.Errorf("scene event = %+v", events[2])
	}

	attr, err := events[3].Attributes()
	if err != nil || attr.LastSeen != "2024-05-01T10:01Z" || events[3].HasState() {
		t.Errorf("attr event = %+v, %v", attr, err)
	}
}

// waitThenCancel lets the third connection come up before stopping, so the test also covers
// shutting down an idle, healthy connection.
func waitThenCancel(connections *atomic.Int32, cancel context.CancelFunc) {
	deadline := time.Now().Add(5 * time.Second)
	for connections.Load() < 3 && time.Now().Before(deadline) {
		time.Sleep(5 * time.Millisecond)
	}
	cancel()
}

func TestRunRetriesWhenGatewayIsDown(t *testing.T) {
	srv := httptest.NewServer(http.NotFoundHandler())
	url := "ws" + strings.TrimPrefix(srv.URL, "http")
	srv.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	var connects atomic.Int32
	err := New(url, WithBackoff(10*time.Millisecond, 20*time.Millisecond),
		WithOnConnect(func(context.Context) { connects.Add(1) })).Run(ctx, func(Event) {})
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("Run returned %v, want deadline exceeded", err)
	}
	if connects.Load() != 0 {
		t.Error("onConnect called without a connection")
	}
}

func TestFilter(t *testing.T) {
	sensor := Event{Kind: EventChanged, Resource: ResourceSensors, ID: "1", State: []byte(`{"open":true}`)}
	light := Event{Kind: EventChanged, Resource: ResourceLights, ID: "4", Attr: []byte(`{}`)}

	f := Filter{Resources: []Resource{ResourceSensors}, WithState: true}
	if !f.Match(sensor) || f.Match(light) {
		t.Error("resource/state filter wrong")
	}
	if (Filter{IDs: []string{"4", "5"}}).Match(sensor) || !(Filter{IDs: []string{"4", "5"}}).Match(light) {
		t.Error("id filter wrong")
	}
	if !(Filter{}).Match(light) {
		t.Error("empty filter must match everything")
	}
	var got int
	f.Handler(func(Event) { got++ })(light)
	f.Handler(func(Event) { got++ })(sensor)
	if got != 1 {
		t.Errorf("filtered handler called %d times", got)
	}
}

func TestURLFromREST(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"ipaddress":"172.17.0.5","websocketport":8443}`))
	}))
	defer srv.Close()
	client, err := rest.New(srv.URL, "KEY")
	if err != nil {
		t.Fatal(err)
	}
	got, err := URLFromREST(context.Background(), client)
	if err != nil {
		t.Fatal(err)
	}
	if got != "ws://127.0.0.1:8443" {
		t.Errorf("URL = %q (must use the REST host, not the gateway's ipaddress)", got)
	}
}
