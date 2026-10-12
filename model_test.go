package godeconz

import (
	"encoding/json"
	"testing"
	"time"
)

func TestParseAddress(t *testing.T) {
	cases := map[string]string{
		"deconz.kani.home":           "http://deconz.kani.home",
		"192.168.1.10:8080":          "http://192.168.1.10:8080",
		"https://deconz.kani.home/":  "https://deconz.kani.home",
		" http://deconz:80/api/foo ": "http://deconz:80",
	}
	for in, want := range cases {
		got, err := ParseAddress(in)
		if err != nil {
			t.Errorf("ParseAddress(%q): %v", in, err)
			continue
		}
		if got.String() != want {
			t.Errorf("ParseAddress(%q) = %q, want %q", in, got, want)
		}
	}
	for _, bad := range []string{"", "ftp://deconz", "http://"} {
		if _, err := ParseAddress(bad); err == nil {
			t.Errorf("ParseAddress(%q) succeeded, want error", bad)
		}
	}
}

func TestDeviceMAC(t *testing.T) {
	if got := DeviceMAC("00:15:8d:00:06:d4:5b:dd-01-0402"); got != "00:15:8d:00:06:d4:5b:dd" {
		t.Errorf("DeviceMAC = %q", got)
	}
	if got := DeviceMAC("00:15:8d:00:06:d4:5b:dd"); got != "00:15:8d:00:06:d4:5b:dd" {
		t.Errorf("DeviceMAC without suffix = %q", got)
	}
}

func TestParseTime(t *testing.T) {
	cases := map[string]time.Time{
		"2024-05-01T10:00:01.123":   time.Date(2024, 5, 1, 10, 0, 1, 123e6, time.UTC),
		"2024-05-01T10:00:01":       time.Date(2024, 5, 1, 10, 0, 1, 0, time.UTC),
		"2024-05-01T10:00Z":         time.Date(2024, 5, 1, 10, 0, 0, 0, time.UTC),
		"2024-05-01T12:00:00+02:00": time.Date(2024, 5, 1, 10, 0, 0, 0, time.UTC),
	}
	for in, want := range cases {
		got, ok := ParseTime(in)
		if !ok || !got.Equal(want) {
			t.Errorf("ParseTime(%q) = %v, %v; want %v", in, got, ok, want)
		}
	}
	for _, in := range []string{"", "none", "yesterday"} {
		if _, ok := ParseTime(in); ok {
			t.Errorf("ParseTime(%q) succeeded", in)
		}
	}
}

func TestSensorStateUnits(t *testing.T) {
	var s SensorState
	data := `{"temperature":2134,"humidity":4567,"pressure":1003,"power":25,"consumption":123456,
		"voltage":231,"current":120,"lastupdated":"2024-05-01T10:00:01.123"}`
	if err := json.Unmarshal([]byte(data), &s); err != nil {
		t.Fatal(err)
	}
	checks := []struct {
		name string
		get  func() (float64, bool)
		want float64
	}{
		{"temperature", s.TemperatureCelsius, 21.34},
		{"humidity", s.HumidityPercent, 45.67},
		{"pressure", s.PressureHPa, 1003},
		{"power", s.PowerWatts, 25},
		{"energy", s.EnergyWattHours, 123456},
		{"voltage", s.VoltageVolts, 231},
		{"current", s.CurrentAmperes, 0.12},
	}
	for _, c := range checks {
		got, ok := c.get()
		if !ok || abs(got-c.want) > 1e-9 {
			t.Errorf("%s = %v, %v; want %v", c.name, got, ok, c.want)
		}
	}
	if _, ok := s.LastUpdatedTime(); !ok {
		t.Error("LastUpdatedTime not parsed")
	}
	if len(s.Raw) != 8 {
		t.Errorf("Raw has %d keys, want 8", len(s.Raw))
	}
	if _, ok := (SensorState{}).TemperatureCelsius(); ok {
		t.Error("TemperatureCelsius on empty state reported ok")
	}
}

func TestSensorStateToleratesUnexpectedTypes(t *testing.T) {
	var s SensorState
	if err := json.Unmarshal([]byte(`{"power":25.5,"open":true,"vendor_x":"y"}`), &s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Power != nil {
		t.Errorf("Power = %v, want nil for a float value", *s.Power)
	}
	if s.Open == nil || !*s.Open {
		t.Error("Open not decoded next to the odd attribute")
	}
	if string(s.Raw["power"]) != "25.5" || string(s.Raw["vendor_x"]) != `"y"` {
		t.Errorf("Raw = %v", s.Raw)
	}
}

func TestSensorConfigAsUpdateSendsOnlySetFields(t *testing.T) {
	data, err := json.Marshal(SensorConfig{Offset: Ptr(-50)})
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != `{"offset":-50}` {
		t.Errorf("body = %s", data)
	}
}

func TestUpdatesSendOnlySetFields(t *testing.T) {
	cases := []struct {
		update any
		want   string
	}{
		{LightStateUpdate{}, `{}`},
		{LightStateUpdate{On: Ptr(false)}, `{"on":false}`},
		{LightStateUpdate{On: Ptr(true), Bri: Ptr[uint8](128), TransitionTime: Ptr[uint16](10)},
			`{"on":true,"bri":128,"transitiontime":10}`},
		{GroupActionUpdate{Toggle: Ptr(true)}, `{"toggle":true}`},
		{GroupActionUpdate{LightStateUpdate: LightStateUpdate{CT: Ptr[uint16](300)}}, `{"ct":300}`},
		{GatewayConfigUpdate{PermitJoin: Ptr[uint8](0)}, `{"permitjoin":0}`},
	}
	for _, c := range cases {
		data, err := json.Marshal(c.update)
		if err != nil {
			t.Fatal(err)
		}
		if string(data) != c.want {
			t.Errorf("%T: body = %s, want %s", c.update, data, c.want)
		}
	}
}

func TestLightStateHelpers(t *testing.T) {
	var s LightState
	if s.IsOn() || !s.IsReachable() || s.Brightness() != 0 {
		t.Error("zero LightState helpers wrong")
	}
	s = LightState{On: Ptr(true), Bri: Ptr(77), Reachable: Ptr(false)}
	if !s.IsOn() || s.IsReachable() || s.Brightness() != 77 {
		t.Error("LightState helpers wrong")
	}
}

func abs(f float64) float64 {
	if f < 0 {
		return -f
	}
	return f
}
