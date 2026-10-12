package ws

import (
	"encoding/json"
	"fmt"

	"github.com/kanimaru/godeconz"
)

// EventKind is the "e" field of a websocket message.
type EventKind string

const (
	EventAdded       EventKind = "added"
	EventChanged     EventKind = "changed"
	EventDeleted     EventKind = "deleted"
	EventSceneCalled EventKind = "scene-called"
)

// Resource is the "r" field of a websocket message.
type Resource string

const (
	ResourceLights       Resource = "lights"
	ResourceSensors      Resource = "sensors"
	ResourceGroups       Resource = "groups"
	ResourceScenes       Resource = "scenes"
	ResourceAlarmSystems Resource = "alarmsystems"
)

// Event is one websocket message. "changed" events carry only the attributes that changed
// unless websocketnotifyall is enabled in the gateway config.
type Event struct {
	Type     string    `json:"t"` // always "event"
	Kind     EventKind `json:"e"`
	Resource Resource  `json:"r"`
	ID       string    `json:"id"`
	UniqueID string    `json:"uniqueid"`
	GroupID  string    `json:"gid"`  // scene-called only
	SceneID  string    `json:"scid"` // scene-called only
	Name     string    `json:"name"` // set on rename

	State  json.RawMessage `json:"state,omitempty"`
	Config json.RawMessage `json:"config,omitempty"`
	Attr   json.RawMessage `json:"attr,omitempty"`
	Light  json.RawMessage `json:"light,omitempty"`  // added lights only
	Sensor json.RawMessage `json:"sensor,omitempty"` // added sensors only
	Group  json.RawMessage `json:"group,omitempty"`  // added groups only

	// Raw is the complete message as received.
	Raw json.RawMessage `json:"-"`
}

// Attr holds the attribute updates deCONZ sends when e.g. lastseen or the name changes.
type Attr struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Type             string `json:"type"`
	ModelID          string `json:"modelid"`
	ManufacturerName string `json:"manufacturername"`
	SWVersion        string `json:"swversion"`
	UniqueID         string `json:"uniqueid"`
	LastSeen         string `json:"lastseen"`
	LastAnnounced    string `json:"lastannounced"`
}

// HasState reports whether the event carries state attributes.
func (e Event) HasState() bool { return len(e.State) > 0 }

// HasConfig reports whether the event carries config attributes.
func (e Event) HasConfig() bool { return len(e.Config) > 0 }

// HasAttr reports whether the event carries attribute updates.
func (e Event) HasAttr() bool { return len(e.Attr) > 0 }

// SensorState decodes the state of a sensor event. Only changed attributes are non-nil.
func (e Event) SensorState() (godeconz.SensorState, error) {
	var s godeconz.SensorState
	return s, decode(e.State, &s, "state")
}

// SensorConfig decodes the config of a sensor event. Only changed attributes are non-nil.
func (e Event) SensorConfig() (godeconz.SensorConfig, error) {
	var c godeconz.SensorConfig
	return c, decode(e.Config, &c, "config")
}

// LightState decodes the state of a light event. Only changed attributes are non-nil.
func (e Event) LightState() (godeconz.LightState, error) {
	var s godeconz.LightState
	return s, decode(e.State, &s, "state")
}

// GroupState decodes the state of a group event.
func (e Event) GroupState() (godeconz.GroupState, error) {
	var s godeconz.GroupState
	return s, decode(e.State, &s, "state")
}

// Attributes decodes the attr part of an event.
func (e Event) Attributes() (Attr, error) {
	var a Attr
	return a, decode(e.Attr, &a, "attr")
}

// AddedSensor decodes the full sensor of an "added" event.
func (e Event) AddedSensor() (godeconz.Sensor, error) {
	var s godeconz.Sensor
	err := decode(e.Sensor, &s, "sensor")
	s.ID = e.ID
	return s, err
}

// AddedLight decodes the full light of an "added" event.
func (e Event) AddedLight() (godeconz.Light, error) {
	var l godeconz.Light
	err := decode(e.Light, &l, "light")
	l.ID = e.ID
	return l, err
}

// AddedGroup decodes the full group of an "added" event.
func (e Event) AddedGroup() (godeconz.Group, error) {
	var g godeconz.Group
	err := decode(e.Group, &g, "group")
	g.ID = e.ID
	return g, err
}

func decode(raw json.RawMessage, target any, part string) error {
	if len(raw) == 0 {
		return nil
	}
	if err := json.Unmarshal(raw, target); err != nil {
		return fmt.Errorf("deconz: decode event %s: %w", part, err)
	}
	return nil
}

// Filter selects events. Within a field the values are alternatives; set fields must all match.
// Empty fields match everything.
type Filter struct {
	Kinds     []EventKind
	Resources []Resource
	IDs       []string
	UniqueIDs []string
	// WithState and WithConfig require the event to carry state or config attributes.
	WithState  bool
	WithConfig bool
}

// Match reports whether the event passes the filter.
func (f Filter) Match(e Event) bool {
	return matches(f.Kinds, e.Kind) &&
		matches(f.Resources, e.Resource) &&
		matches(f.IDs, e.ID) &&
		matches(f.UniqueIDs, e.UniqueID) &&
		(!f.WithState || e.HasState()) &&
		(!f.WithConfig || e.HasConfig())
}

// Handler wraps handle so it only sees events that pass the filter.
func (f Filter) Handler(handle func(Event)) func(Event) {
	return func(e Event) {
		if f.Match(e) {
			handle(e)
		}
	}
}

func matches[T comparable](allowed []T, value T) bool {
	if len(allowed) == 0 {
		return true
	}
	for _, a := range allowed {
		if a == value {
			return true
		}
	}
	return false
}
