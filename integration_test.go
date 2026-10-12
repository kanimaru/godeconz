package godeconz_test

// Read-only checks against a real gateway. Skipped unless DECONZ_ADDRESS and DECONZ_API_KEY are
// set. DECONZ_WS_SECONDS > 0 also listens on the websocket for that long.
//
//	DECONZ_ADDRESS=deconz.kani.home DECONZ_API_KEY=... DECONZ_WS_SECONDS=60 go test -run Gateway -v .

import (
	"cmp"
	"context"
	"os"
	"slices"
	"strconv"
	"testing"
	"time"

	"github.com/kanimaru/godeconz"
	"github.com/kanimaru/godeconz/rest"
	"github.com/kanimaru/godeconz/ws"
)

func gatewayClient(t *testing.T) *rest.Client {
	t.Helper()
	address, key := os.Getenv("DECONZ_ADDRESS"), os.Getenv("DECONZ_API_KEY")
	if address == "" || key == "" {
		t.Skip("DECONZ_ADDRESS and DECONZ_API_KEY not set")
	}
	client, err := rest.New(address, key)
	if err != nil {
		t.Fatal(err)
	}
	return client
}

func TestGatewayREST(t *testing.T) {
	client := gatewayClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cfg, err := client.Config(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("gateway %q, deCONZ %s, firmware %s, websocket port %d", cfg.Name, cfg.SWVersion, cfg.FWVersion, cfg.WebsocketPort)

	lights, err := client.Lights(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%d lights", len(lights))
	for _, l := range sortedByID(lights) {
		t.Logf("  light %-3s %-28q %-26s on=%v bri=%d reachable=%v", l.ID, l.Name, l.Type, l.State.IsOn(), l.State.Brightness(), l.State.IsReachable())
	}

	sensors, err := client.Sensors(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%d sensors", len(sensors))
	unmodelled := map[string][]string{}
	for _, s := range sortedByID(sensors) {
		t.Logf("  sensor %-3s %-28q %-18s reachable=%v", s.ID, s.Name, s.Type, s.Config.IsReachable())
		for key := range s.State.Raw {
			if !knownStateKeys[key] && !slices.Contains(unmodelled[s.Type], key) {
				unmodelled[s.Type] = append(unmodelled[s.Type], key)
			}
		}
	}
	for sensorType, keys := range unmodelled {
		t.Logf("state attributes without a field (still in Raw): %s %v", sensorType, keys)
	}

	groups, err := client.Groups(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%d groups", len(groups))
	for _, g := range sortedByID(groups) {
		t.Logf("  group %-3s %-20q %-10s lights=%v any_on=%v scenes=%d", g.ID, g.Name, g.Type, g.Lights, g.State.AnyOn, len(g.Scenes))
	}
}

func TestGatewayWebsocket(t *testing.T) {
	client := gatewayClient(t)
	seconds, _ := strconv.Atoi(os.Getenv("DECONZ_WS_SECONDS"))
	if seconds <= 0 {
		t.Skip("DECONZ_WS_SECONDS not set")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(seconds)*time.Second)
	defer cancel()

	url, err := ws.URLFromREST(ctx, client)
	if err != nil {
		t.Fatal(err)
	}
	connects := 0
	counts := map[string]int{}
	err = ws.New(url, ws.WithOnConnect(func(context.Context) { connects++ })).Run(ctx, func(e ws.Event) {
		counts[string(e.Resource)+"/"+string(e.Kind)]++
		if e.Resource == ws.ResourceSensors && e.HasState() {
			if _, err := e.SensorState(); err != nil {
				t.Errorf("sensor %s: %v (%s)", e.ID, err, e.Raw)
			}
		}
	})
	if err != context.DeadlineExceeded {
		t.Errorf("Run returned %v", err)
	}
	if connects == 0 {
		t.Fatalf("never connected to %s", url)
	}
	t.Logf("%s: %d connect(s), events in %ds: %v", url, connects, seconds, counts)
}

var knownStateKeys = map[string]bool{
	"lastupdated": true, "temperature": true, "humidity": true, "pressure": true, "lux": true,
	"lightlevel": true, "dark": true, "daylight": true, "presence": true, "open": true,
	"power": true, "consumption": true, "voltage": true, "current": true, "buttonevent": true,
	"gesture": true, "eventduration": true, "battery": true, "lowbattery": true, "tampered": true,
	"water": true, "fire": true, "carbonmonoxide": true, "vibration": true, "alarm": true,
	"airquality": true, "airqualityppb": true, "pm2_5": true,
}

type identified interface {
	godeconz.Light | godeconz.Sensor | godeconz.Group
}

func sortedByID[T identified](m map[string]T) []T {
	ids := make([]string, 0, len(m))
	for id := range m {
		ids = append(ids, id)
	}
	slices.SortFunc(ids, func(a, b string) int {
		ai, aErr := strconv.Atoi(a)
		bi, bErr := strconv.Atoi(b)
		if aErr == nil && bErr == nil {
			return ai - bi
		}
		return cmp.Compare(a, b)
	})
	out := make([]T, 0, len(ids))
	for _, id := range ids {
		out = append(out, m[id])
	}
	return out
}
