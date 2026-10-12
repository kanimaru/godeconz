package rest

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/kanimaru/godeconz"
)

const testKey = "SECRETKEY1234"

// fakeGateway answers by "METHOD /path" (path without /api/<key>) and records the requests.
type fakeGateway struct {
	t        *testing.T
	mu       sync.Mutex
	routes   map[string]response
	requests []recorded
}

type response struct {
	status int
	body   string
}

type recorded struct {
	method string
	path   string
	body   string
}

func newFakeGateway(t *testing.T, routes map[string]response) (*Client, *fakeGateway) {
	t.Helper()
	g := &fakeGateway{t: t, routes: routes}
	srv := httptest.NewServer(g)
	t.Cleanup(srv.Close)
	client, err := New(srv.URL, testKey)
	if err != nil {
		t.Fatal(err)
	}
	return client, g
}

func (g *fakeGateway) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	path := r.URL.Path
	prefix := "/api/" + testKey
	if strings.HasPrefix(path, prefix) {
		path = strings.TrimPrefix(path, prefix)
	}
	g.mu.Lock()
	g.requests = append(g.requests, recorded{r.Method, path, string(body)})
	g.mu.Unlock()

	resp, ok := g.routes[r.Method+" "+path]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		_, _ = io.WriteString(w, `[{"error":{"type":3,"address":"`+path+`","description":"resource, `+path+`, not available"}}]`)
		return
	}
	if resp.status == 0 {
		resp.status = http.StatusOK
	}
	w.WriteHeader(resp.status)
	_, _ = io.WriteString(w, resp.body)
}

func (g *fakeGateway) last() recorded {
	g.mu.Lock()
	defer g.mu.Unlock()
	if len(g.requests) == 0 {
		g.t.Fatal("no request recorded")
	}
	return g.requests[len(g.requests)-1]
}

const sensorsJSON = `{
 "1": {"config":{"battery":100,"offset":0,"on":true,"reachable":true},"ep":1,"etag":"e1",
       "lastseen":"2024-05-01T10:00Z","manufacturername":"LUMI","modelid":"lumi.weather",
       "name":"BedroomSensor","state":{"lastupdated":"2024-05-01T10:00:01.123","temperature":2134},
       "swversion":"3000-0001","type":"ZHATemperature","uniqueid":"00:15:8d:00:04:aa:df:08-01-0402"},
 "2": {"config":{"on":true,"reachable":true},"name":"Plug Power","type":"ZHAPower",
       "state":{"current":120,"lastupdated":"2024-05-01T10:00:02.000","power":25,"voltage":231},
       "uniqueid":"a4:c1:38:00:00:00:00:01-01-0b04"},
 "3": {"config":{"on":true,"reachable":false},"name":"Plug Consumption","type":"ZHAConsumption",
       "state":{"consumption":123456,"lastupdated":"2024-05-01T10:00:02.000","power":25},
       "uniqueid":"a4:c1:38:00:00:00:00:01-01-0702"},
 "4": {"config":{"battery":87,"on":true,"reachable":true,"pending":[]},"name":"EntranceDoor",
       "type":"ZHAOpenClose","state":{"lastupdated":"2024-05-01T09:00:00.000","lowbattery":false,"open":true},
       "uniqueid":"00:15:8d:00:01:02:03:04-01-0006"},
 "5": {"config":{"on":true},"name":"Odd","type":"ZHAPower","state":{"power":12.5},"uniqueid":"x-01"}
}`

func TestSensorsDecodesAllCommonTypes(t *testing.T) {
	client, gw := newFakeGateway(t, map[string]response{"GET /sensors": {body: sensorsJSON}})
	sensors, err := client.Sensors(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if r := gw.last(); r.method != http.MethodGet || r.path != "/sensors" {
		t.Errorf("request = %+v", r)
	}
	if len(sensors) != 5 {
		t.Fatalf("got %d sensors", len(sensors))
	}
	for id, s := range sensors {
		if s.ID != id {
			t.Errorf("sensor %s has ID %q", id, s.ID)
		}
	}
	if temp, ok := sensors["1"].State.TemperatureCelsius(); !ok || temp != 21.34 {
		t.Errorf("temperature = %v, %v", temp, ok)
	}
	if b := sensors["1"].Config.Battery; b == nil || *b != 100 {
		t.Errorf("battery = %v", b)
	}
	if w, _ := sensors["2"].State.PowerWatts(); w != 25 {
		t.Errorf("power = %v", w)
	}
	if a, _ := sensors["2"].State.CurrentAmperes(); a != 0.12 {
		t.Errorf("current = %v", a)
	}
	if wh, _ := sensors["3"].State.EnergyWattHours(); wh != 123456 {
		t.Errorf("energy = %v", wh)
	}
	if sensors["3"].Config.IsReachable() {
		t.Error("sensor 3 should be unreachable")
	}
	if o := sensors["4"].State.Open; o == nil || !*o {
		t.Errorf("open = %v", o)
	}
	if _, ok := sensors["4"].Config.Raw["pending"]; !ok {
		t.Error("config attribute without field missing from Raw")
	}
	if sensors["5"].State.Power != nil || string(sensors["5"].State.Raw["power"]) != "12.5" {
		t.Error("odd float attribute not handled leniently")
	}
	if godeconz.DeviceMAC(sensors["2"].UniqueID) != godeconz.DeviceMAC(sensors["3"].UniqueID) {
		t.Error("plug resources should share a MAC")
	}
}

func TestLightsAndSetLightState(t *testing.T) {
	client, gw := newFakeGateway(t, map[string]response{
		"GET /lights": {body: `{"4":{"name":"Workfloor-1","type":"Dimmable light","uniqueid":"00:0d:6f:00:10:65:8a:6e-01",
			"lastannounced":null,"state":{"on":true,"bri":178,"alert":"none","reachable":true}}}`},
		"PUT /lights/4/state": {body: `[{"success":{"/lights/4/state/on":true}},{"success":{"/lights/4/state/bri":10}}]`},
	})
	ctx := context.Background()
	lights, err := client.Lights(ctx)
	if err != nil {
		t.Fatal(err)
	}
	l := lights["4"]
	if l.ID != "4" || !l.State.IsOn() || l.State.Brightness() != 178 || l.LastAnnounced != "" {
		t.Errorf("light = %+v", l)
	}

	err = client.SetLightState(ctx, "4", godeconz.LightStateUpdate{On: godeconz.Ptr(true), Bri: godeconz.Ptr[uint8](10)})
	if err != nil {
		t.Fatal(err)
	}
	r := gw.last()
	if r.method != http.MethodPut || r.path != "/lights/4/state" || r.body != `{"on":true,"bri":10}` {
		t.Errorf("request = %+v", r)
	}
}

func TestCreateUsesPOSTAndReturnsID(t *testing.T) {
	client, gw := newFakeGateway(t, map[string]response{
		"POST /groups":          {body: `[{"success":{"id":"7"}}]`},
		"POST /groups/7/scenes": {body: `[{"success":{"id":"2"}}]`},
	})
	ctx := context.Background()
	id, err := client.CreateGroup(ctx, "Kitchen")
	if err != nil || id != "7" {
		t.Fatalf("CreateGroup = %q, %v", id, err)
	}
	if r := gw.last(); r.method != http.MethodPost || r.body != `{"name":"Kitchen"}` {
		t.Errorf("request = %+v", r)
	}
	sceneID, err := client.CreateScene(ctx, "7", "Evening")
	if err != nil || sceneID != "2" {
		t.Fatalf("CreateScene = %q, %v", sceneID, err)
	}
}

func TestGroupsSceneRecallAndAction(t *testing.T) {
	client, gw := newFakeGateway(t, map[string]response{
		"GET /groups": {body: `{"1":{"name":"Kitchen","type":"Room","class":"Kitchen","lights":["10","12"],
			"state":{"all_on":false,"any_on":true},"action":{"on":true,"bri":127},
			"scenes":[{"id":"2","name":"Evening","lightcount":2,"transitiontime":10}]}}`},
		"PUT /groups/1/scenes/2/recall": {body: `[{"success":{"id":"2"}}]`},
		"PUT /groups/1/action":          {body: `[{"success":{"/groups/1/action/toggle":true}}]`},
	})
	ctx := context.Background()
	groups, err := client.Groups(ctx)
	if err != nil {
		t.Fatal(err)
	}
	g := groups["1"]
	if g.ID != "1" || !g.State.AnyOn || g.State.AllOn || len(g.Lights) != 2 || g.Scenes[0].Name != "Evening" {
		t.Errorf("group = %+v", g)
	}
	if err := client.RecallScene(ctx, "1", "2"); err != nil {
		t.Fatal(err)
	}
	if r := gw.last(); r.method != http.MethodPut || r.path != "/groups/1/scenes/2/recall" {
		t.Errorf("request = %+v", r)
	}
	if err := client.SetGroupAction(ctx, "1", godeconz.GroupActionUpdate{Toggle: godeconz.Ptr(true)}); err != nil {
		t.Fatal(err)
	}
	if r := gw.last(); r.body != `{"toggle":true}` {
		t.Errorf("action body = %s", r.body)
	}
}

func TestErrors(t *testing.T) {
	client, _ := newFakeGateway(t, map[string]response{
		"GET /lights": {status: http.StatusForbidden,
			body: `[{"error":{"type":1,"address":"/","description":"unauthorized user"}}]`},
		"PUT /lights/4/state": {body: `[{"success":{"/lights/4/state/on":true}},
			{"error":{"type":201,"address":"/lights/4/state/bri","description":"parameter, bri, is not modifiable. Device is set to off."}}]`},
		"GET /config": {status: http.StatusBadGateway, body: "<html>proxy error</html>"},
	})
	ctx := context.Background()

	_, err := client.Lights(ctx)
	if !errors.Is(err, ErrUnauthorized) || errors.Is(err, ErrNotFound) {
		t.Errorf("Lights error = %v", err)
	}

	_, err = client.Sensor(ctx, "99")
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("missing sensor error = %v", err)
	}

	err = client.SetLightState(ctx, "4", godeconz.LightStateUpdate{On: godeconz.Ptr(true), Bri: godeconz.Ptr[uint8](10)})
	var apiErr *Error
	if !errors.As(err, &apiErr) || len(apiErr.Errors) != 1 || apiErr.Errors[0].Type != ErrorTypeDeviceOff {
		t.Errorf("partial failure error = %v", err)
	}

	_, err = client.Config(ctx)
	if !errors.As(err, &apiErr) || apiErr.StatusCode != http.StatusBadGateway || !strings.Contains(err.Error(), "proxy error") {
		t.Errorf("non-JSON error = %v", err)
	}
}

func TestErrorsNeverContainAPIKey(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	addr := srv.URL
	srv.Close() // connection refused from now on

	client, err := New(addr, testKey)
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.Lights(context.Background())
	if err == nil {
		t.Fatal("expected an error from a closed server")
	}
	if strings.Contains(err.Error(), testKey) {
		t.Errorf("error leaks the API key: %v", err)
	}

	client, _ = newFakeGateway(t, nil)
	_, err = client.Light(context.Background(), "1")
	if err == nil || strings.Contains(err.Error(), testKey) {
		t.Errorf("API error leaks the key or is missing: %v", err)
	}
}

func TestContextCancellation(t *testing.T) {
	client, _ := newFakeGateway(t, map[string]response{"GET /lights": {body: `{}`}})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := client.Lights(ctx); !errors.Is(err, context.Canceled) {
		t.Errorf("err = %v, want context.Canceled", err)
	}
}

func TestEmptyMapsAndFullState(t *testing.T) {
	client, _ := newFakeGateway(t, map[string]response{
		"GET /groups": {body: `{}`},
		"GET ": {body: `{"config":{"name":"Phoscon-GW","websocketport":8443},
			"lights":{"1":{"name":"L"}},"sensors":{"2":{"name":"S","state":{"open":false}}},"groups":{"3":{"name":"G"}}}`},
	})
	ctx := context.Background()
	groups, err := client.Groups(ctx)
	if err != nil || groups == nil || len(groups) != 0 {
		t.Errorf("Groups = %v, %v", groups, err)
	}
	full, err := client.FullState(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if full.Config.WebsocketPort != 8443 || full.Lights["1"].ID != "1" || full.Sensors["2"].ID != "2" || full.Groups["3"].ID != "3" {
		t.Errorf("full state = %+v", full)
	}
}

func TestAcquireAPIKey(t *testing.T) {
	var gotBody map[string]string
	locked := true
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		_ = json.NewDecoder(r.Body).Decode(&gotBody)
		if locked {
			w.WriteHeader(http.StatusForbidden)
			_, _ = io.WriteString(w, `[{"error":{"type":101,"address":"/","description":"link button not pressed"}}]`)
			return
		}
		_, _ = io.WriteString(w, `[{"success":{"username":"NEWKEY"}}]`)
	}))
	defer srv.Close()
	ctx := context.Background()

	if _, err := AcquireAPIKey(ctx, srv.URL, "hestia#test"); !errors.Is(err, ErrLinkButton) {
		t.Errorf("locked gateway error = %v", err)
	}
	locked = false
	key, err := AcquireAPIKey(ctx, srv.URL, "hestia#test")
	if err != nil || key != "NEWKEY" {
		t.Fatalf("AcquireAPIKey = %q, %v", key, err)
	}
	if gotBody["devicetype"] != "hestia#test" {
		t.Errorf("body = %v", gotBody)
	}
}

func TestPathEscaping(t *testing.T) {
	client, gw := newFakeGateway(t, nil)
	_ = client.DeleteSensor(context.Background(), "../config")
	if r := gw.last(); r.path == "/config" {
		t.Errorf("id escaped the resource path: %+v", r)
	}
}
