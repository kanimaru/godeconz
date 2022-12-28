package ws

import "encoding/json"

type Message struct {
	// recv: {"e":"changed","id":"43","r":"sensors","state":{"lastupdated":"2022-02-14T22:46:22.317","presence":true},"t":"event","uniqueid":"00:15:8d:00:06:f4:5f:99-01-0406"}
	Emit     string          `json:"e,omitempty"`
	Id       string          `json:"id,omitempty"`
	Role     string          `json:"r,omitempty"`
	Type     string          `json:"t,omitempty"`
	Attr     json.RawMessage `json:"attr,omitempty"`
	State    json.RawMessage `json:"state,omitempty"`
	Config   json.RawMessage `json:"config,omitempty"`
	UniqueId string          `json:"uniqueid,omitempty"`
}
