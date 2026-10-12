package godeconz

// Alert is a temporary alert effect of a light or group.
type Alert string

const (
	AlertNone    Alert = "none"
	AlertSelect  Alert = "select"  // blink once
	AlertLSelect Alert = "lselect" // blink for a longer time
)

// Effect is a dynamic effect of a light or group.
type Effect string

const (
	EffectNone      Effect = "none"
	EffectColorLoop Effect = "colorloop"
)

// Light is a light or switchable plug as returned by GET /lights.
type Light struct {
	// ID is the REST resource id (the key in GET /lights). It changes when a device is re-paired;
	// UniqueID does not.
	ID                string     `json:"-"`
	Name              string     `json:"name"`
	Type              string     `json:"type"`
	ModelID           string     `json:"modelid"`
	ManufacturerName  string     `json:"manufacturername"`
	UniqueID          string     `json:"uniqueid"`
	SWVersion         string     `json:"swversion"`
	ETag              string     `json:"etag"`
	LastAnnounced     string     `json:"lastannounced"`
	LastSeen          string     `json:"lastseen"`
	HasColor          *bool      `json:"hascolor,omitempty"`
	ColorCapabilities *int       `json:"colorcapabilities,omitempty"`
	CTMin             *int       `json:"ctmin,omitempty"`
	CTMax             *int       `json:"ctmax,omitempty"`
	State             LightState `json:"state"`
}

// LightState is the state of a light, the last action of a group, or the changed part of a
// light state in a websocket event. Attributes a device does not report are nil.
type LightState struct {
	On        *bool     `json:"on,omitempty"`
	Bri       *int      `json:"bri,omitempty"` // 0–255
	Hue       *int      `json:"hue,omitempty"` // 0–65535
	Sat       *int      `json:"sat,omitempty"` // 0–255
	CT        *int      `json:"ct,omitempty"`  // mired
	XY        []float64 `json:"xy,omitempty"`  // CIE x, y
	Alert     *string   `json:"alert,omitempty"`
	ColorMode *string   `json:"colormode,omitempty"`
	Effect    *string   `json:"effect,omitempty"`
	Speed     *int      `json:"speed,omitempty"` // fans, 0–6
	Reachable *bool     `json:"reachable,omitempty"`
}

// IsOn reports whether the light is on; unknown counts as off.
func (s LightState) IsOn() bool {
	return s.On != nil && *s.On
}

// IsReachable reports whether the light is reachable; unknown counts as reachable, because
// group actions and partial events do not carry the attribute.
func (s LightState) IsReachable() bool {
	return s.Reachable == nil || *s.Reachable
}

// Brightness returns the brightness, 0 if unknown.
func (s LightState) Brightness() int {
	if s.Bri == nil {
		return 0
	}
	return *s.Bri
}

// LightStateUpdate is the body of PUT /lights/{id}/state. Only set fields are sent.
type LightStateUpdate struct {
	On             *bool     `json:"on,omitempty"`
	Bri            *uint8    `json:"bri,omitempty"`
	Hue            *uint16   `json:"hue,omitempty"`
	Sat            *uint8    `json:"sat,omitempty"`
	CT             *uint16   `json:"ct,omitempty"`
	XY             []float64 `json:"xy,omitempty"`
	Alert          Alert     `json:"alert,omitempty"`
	Effect         Effect    `json:"effect,omitempty"`
	ColorLoopSpeed *uint8    `json:"colorloopspeed,omitempty"` // 1 fast – 255 slow
	TransitionTime *uint16   `json:"transitiontime,omitempty"` // 1/10 s
}

// LightAttributesUpdate is the body of PUT /lights/{id}.
type LightAttributesUpdate struct {
	Name *string `json:"name,omitempty"`
}
