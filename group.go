package godeconz

// Group is a group of lights as returned by GET /groups. deCONZ rooms are groups.
type Group struct {
	ID               string     `json:"id"`
	Name             string     `json:"name"`
	Type             string     `json:"type"`
	Class            string     `json:"class"`
	Hidden           bool       `json:"hidden"`
	Lights           []string   `json:"lights"`
	LightSequence    []string   `json:"lightsequence"`
	MultiDeviceIDs   []string   `json:"multideviceids"`
	DeviceMembership []string   `json:"devicemembership"`
	Scenes           []SceneRef `json:"scenes"`
	ETag             string     `json:"etag"`
	// Action is the last action sent to the group, not the actual state of its lights.
	Action LightState `json:"action"`
	State  GroupState `json:"state"`
}

// GroupState summarises the lights of a group.
type GroupState struct {
	AllOn bool `json:"all_on"`
	AnyOn bool `json:"any_on"`
}

// GroupActionUpdate is the body of PUT /groups/{id}/action. Only set fields are sent.
type GroupActionUpdate struct {
	LightStateUpdate
	// Toggle switches the lights on or off; it takes precedence over On.
	Toggle *bool `json:"toggle,omitempty"`
}

// GroupAttributesUpdate is the body of PUT /groups/{id}. Only set fields are sent, so an empty
// Lights slice cannot remove all lights from a group.
type GroupAttributesUpdate struct {
	Name           *string  `json:"name,omitempty"`
	Lights         []string `json:"lights,omitempty"`
	Hidden         *bool    `json:"hidden,omitempty"`
	LightSequence  []string `json:"lightsequence,omitempty"`
	MultiDeviceIDs []string `json:"multideviceids,omitempty"`
}

// SceneRef is a scene as listed in a group or in GET /groups/{id}/scenes.
type SceneRef struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	TransitionTime int    `json:"transitiontime"`
	LightCount     int    `json:"lightcount"`
}

// Scene is a stored scene with the state of each of its lights.
type Scene struct {
	ID     string       `json:"-"`
	Name   string       `json:"name"`
	State  int          `json:"state"`
	Lights []SceneLight `json:"lights"`
}

// SceneLight is the stored state of one light in a scene.
type SceneLight struct {
	ID             string   `json:"id"`
	On             *bool    `json:"on,omitempty"`
	Bri            *int     `json:"bri,omitempty"`
	CT             *int     `json:"ct,omitempty"`
	X              *float64 `json:"x,omitempty"`
	Y              *float64 `json:"y,omitempty"`
	ColorMode      *string  `json:"colormode,omitempty"`
	TransitionTime *int     `json:"transitiontime,omitempty"`
}
