package ws

// There is no documentation about the websocket states so this will be only user gathered states here.
// Add your own as soon as you discover them.

type StateDaylight struct {
	// FIXME HERE ARE MISSING FIELDS
}

type StateBaseLight struct {
	Alert interface{} `json:"alert,omitempty"`
	Bri   int         `json:"bri,omitempty"`
	On    bool        `json:"on,omitempty"`
}

type StateDimmablelight struct {
	StateBaseLight
	Reachable bool `json:"reachable,omitempty"`
}

type StateColortemperaturelight struct {
	StateBaseLight
	Colormode string `json:"colormode"`
	Ct        int    `json:"ct"`
	Reachable bool   `json:"reachable"`
}

type StateExtendedcolorlight struct {
	// FIXME HERE ARE MISSING FIELDS
	StateBaseLight
	Colormode string `json:"colormode"`
	Ct        int    `json:"ct"`
	Reachable bool   `json:"reachable"`
}

type StateOnOffpluginunit struct {
	Alert     interface{} `json:"alert"`
	On        bool        `json:"on"`
	Reachable bool        `json:"reachable"`
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
