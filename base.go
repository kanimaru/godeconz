package godeconz

// All stuff that are same in WS and HTTP

type Config interface {
	IsOn() bool
	IsReachable() *bool
	GetBattery() *uint8
}

// StateLight is the base for all lights
type StateLight interface {
	GetAlert() interface{}
	GetBrightness() int
	IsOn() bool
}

// StateDevice is the base for all "light" devices also for OnOff  Sockets
type StateDevice interface {
	IsOn() bool
}
