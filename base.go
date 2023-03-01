package godeconz

// All stuff that are same in WS and HTTP

type Config interface {
	IsOn() bool
	IsReachable() *bool
	GetBattery() *uint8
}

type StateLight interface {
	GetAlert() interface{}
	GetBrightness() int
	IsOn() bool
}
