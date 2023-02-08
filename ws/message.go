package ws

import (
	"encoding/json"
	"time"
)

type EventType string

// EventTypeAdded resource has been added
const EventTypeAdded EventType = "added"

// EventTypeChanged resource attributes have changed;
const EventTypeChanged EventType = "changed"

// EventTypeDeleted resource has been deleted.
const EventTypeDeleted EventType = "deleted"

// EventTypeSceneCalled a scene has been recalled.
const EventTypeSceneCalled EventType = "scene-called"

type ResourceType string

// ResourceTypeGroups message relates to a group resource;
const ResourceTypeGroups ResourceType = "Groups"

// ResourceTypeLights message relates to a light resource;
const ResourceTypeLights ResourceType = "Lights"

// ResourceTypeScenes message relates to a scene under a group resource;
const ResourceTypeScenes ResourceType = "Scenes"

// ResourceTypeSensors message relates to a sensor resource.
const ResourceTypeSensors ResourceType = "Sensors"

type MessageType string

// MessageTypeEvent the message holds an event.
const MessageTypeEvent MessageType = "event"

// Attr is not documented but found in the wild. Here are probably other missing fields
type Attr struct {
	Id               string    `json:"id,omitempty"`
	Lastannounced    time.Time `json:"lastannounced,omitempty"`
	Lastseen         string    `json:"lastseen,omitempty"`
	Manufacturername string    `json:"manufacturername,omitempty"`
	Modelid          string    `json:"modelid,omitempty"`
	Name             string    `json:"name,omitempty"`
	Swversion        string    `json:"swversion,omitempty"`
	Type             string    `json:"type,omitempty"`
	Uniqueid         string    `json:"uniqueid,omitempty"`
}

type Message struct {
	// EventType the event type of the message
	EventType EventType `json:"e,omitempty"`
	// Id of the resource to which the message relates, e.g. 5 for /sensors/5. Not for scene-called events.
	Id string `json:"id,omitempty"`
	// ResourceType The resource type to which the message belongs
	ResourceType ResourceType `json:"r,omitempty"`
	// MessageType The type of the message
	MessageType MessageType `json:"t,omitempty"`
	// Depending on the websocketnotifyall setting: a map containing all or only the changed state attributes of a group, light, or sensor resource. Only for changed events.
	State json.RawMessage `json:"state,omitempty"`
	// Depending on the websocketnotifyall setting: a map containing all or only the changed config attributes of a sensor resource. Only for changed events.
	Config json.RawMessage `json:"config,omitempty"`
	// Not documented but found in the wild.
	Attr *Attr `json:"attr,omitempty"`
	// UniqueId of the resource to which the message relates, e.g. 00:0d:6f:00:10:65:8a:6e-01-1000. Only for light and sensor resources.
	UniqueId string `json:"uniqueid,omitempty"`
	// The (new) name of a resource. Only for changed events.
	Name string `json:"name,omitempty"`
	// The group id of the resource to which the message relates. Only for scene-called events.
	GroupId string `json:"GroupId,omitempty"`
	// The scene id of the resource to which the message relates. Only for scene-called events.
	SceneId string `json:"SceneId,omitempty"`
	// The full group resource. Only for added events of a group resource.
	Group json.RawMessage `json:"group,omitempty"`
	// The full light resource. Only for added events of a light resource.
	Light json.RawMessage `json:"light,omitempty"`
	//	The full sensor resource. Only for added events of a sensor resource.
	Sensor json.RawMessage `json:"sensor,omitempty"`
}

func (m Message) StateAs(data interface{}) error {
	if m.State == nil {
		return nil
	}
	return json.Unmarshal(m.State, data)
}

func (m Message) ConfigAs(data interface{}) error {
	if m.Config == nil {
		return nil
	}
	return json.Unmarshal(m.Config, data)
}

func (m Message) GroupAs(data interface{}) error {
	if m.Group == nil {
		return nil
	}
	return json.Unmarshal(m.Group, data)
}

func (m Message) LightAs(data interface{}) error {
	if m.Light == nil {
		return nil
	}
	return json.Unmarshal(m.Light, data)
}

func (m Message) SensorAs(data interface{}) error {
	if m.Sensor == nil {
		return nil
	}
	return json.Unmarshal(m.Sensor, data)
}
