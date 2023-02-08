package http

type CLIPStateAlarm struct {
	Alarm bool `json:"alarm"`
}

type CLIPStateBattery struct {
	Battery uint8 `json:"battery"`
}

type CLIPStateCarbonMonoxide struct {
	Carbonmonoxide bool `json:"carbonmonoxide"`
}

type CLIPStateConsumption struct {
	Consumption uint64 `json:"consumption"`
}

type CLIPStateFire struct {
	Fire bool `json:"fire"`
}

type CLIPStateGenericFlag struct {
	Flag bool `json:"flag"`
}

type CLIPStateGenericStatus struct {
	Status int32 `json:"status"`
}

type CLIPStateHumidity struct {
	Humidity uint16 `json:"humidity"`
}

type CLIPStateLightLevel struct {
	Lightlevel uint16 `json:"lightlevel"`
	Lux        uint32 `json:"lux"`
	Dark       bool   `json:"dark"`
	Daylight   bool   `json:"daylight"`
}

type CLIPStateOpenClose struct {
	Open bool `json:"open"`
}

type CLIPStatePower struct {
	Power   int16  `json:"power"`
	Voltage uint16 `json:"voltage"`
	Current uint16 `json:"current"`
}

type CLIPStatePresence struct {
	Presence bool   `json:"presence"`
	Duration uint16 `json:"duration"`
}

type CLIPStatePressure struct {
	Pressure int16 `json:"pressure"`
}

type CLIPStateSwitch struct {
	Buttonevent uint32 `json:"buttonevent"`
}

type CLIPStateTemperature struct {
	Temperature int16 `json:"temperature"`
}

type CLIPStateVibration struct {
	Vibration bool `json:"vibration"`
}

type CLIPStateWater struct {
	Water bool `json:"water"`
}
