package ws

// There is no documentation about the websocket states so this will be only user gathered states here.
// Add your own as soon as you discover them.

type StateDaylight struct {
	// FIXME HERE ARE MISSING FIELDS
}

type StateBaseDevice struct {
	On        bool `json:"on,omitempty"`
	Reachable bool `json:"reachable,omitempty"`
}

func (s StateBaseDevice) IsOn() bool {
	return s.On
}

func (s StateBaseDevice) IsReachable() bool {
	return s.Reachable
}

type StateBaseLight struct {
	StateBaseDevice
	Alert interface{} `json:"alert,omitempty"`
	Bri   int         `json:"bri,omitempty"`
}

func (s StateBaseLight) GetAlert() interface{} {
	return s.Alert
}

func (s StateBaseLight) GetBrightness() int {
	return s.Bri
}

type StateDimmablelight struct {
	StateBaseLight
}

type StateColortemperaturelight struct {
	StateBaseLight
	Colormode string `json:"colormode"`
	Ct        int    `json:"ct"`
}

type StateExtendedcolorlight struct {
	// FIXME HERE ARE MISSING FIELDS
	StateBaseLight
	Colormode string `json:"colormode"`
	Ct        int    `json:"ct"`
}

type StateOnOffpluginunit struct {
	StateBaseDevice
	Alert interface{} `json:"alert"`
}

type StateZHATemperature struct {
	Lastupdated string `json:"lastupdated"`
	Temperature int    `json:"temperature"`
}

type StateZHAHumidity struct {
	Lastupdated string `json:"lastupdated"`
	Humidity    int    `json:"humidity"`
}

type StateZHAPressure struct {
	Lastupdated string `json:"lastupdated"`
	Pressure    int    `json:"pressure"`
}

type StateZHALightLevel struct {
	Dark        bool   `json:"dark"`
	Daylight    bool   `json:"daylight"`
	Lastupdated string `json:"lastupdated"`
	Lightlevel  int    `json:"lightlevel"`
	Lux         int    `json:"lux"`
}

type StateZHAPresence struct {
	Lastupdated string `json:"lastupdated"`
	Presence    bool   `json:"presence"`
}

type StateZHAOpenClose struct {
	Lastupdated string `json:"lastupdated"`
	Open        bool   `json:"open"`
}
