package attribute

type HAAirQualityValues string

const HAAirQualityValuesExcellent = "excellent"
const HAAirQualityValuesGood = "good"
const HAAirQualityValuesModerate = "moderate"
const HAAirQualityValuesPoor = "poor"
const HAAirQualityValuesUnhealthy = "unhealthy"
const HAAirQualityValuesOutOfScale = "out of scale"

type StateHAAirQuality struct {
	Airquality    HAAirQualityValues `json:"airquality"`
	Airqualityppb uint16             `json:"airqualityppb"`
}

type StateZHAAlarm struct {
	Alarm bool `json:"alarm"`
	// ISO 8601 timestamp
	Lastupdated string `json:"lastupdated"`
	Lowbattery  bool   `json:"lowbattery"`
	Tampered    bool   `json:"tampered"`
}

type StateZHACarbonMonoxide struct {
	Carbonmonoxide bool `json:"carbonmonoxide"`
	// ISO 8601 timestamp
	Lastupdated string `json:"lastupdated"`
	Lowbattery  bool   `json:"lowbattery"`
	Tampered    bool   `json:"tampered"`
}

type StateZHAConsumption struct {
	Consumption uint `json:"consumption"`
	// ISO 8601 timestamp
	Lastupdated string `json:"lastupdated"`
	Power       int    `json:"power"`
}

type StateZHAFire struct {
	Fire bool `json:"fire"`
	// ISO 8601 timestamp
	Lastupdated string `json:"lastupdated"`
	Lowbattery  bool   `json:"lowbattery"`
	Tampered    bool   `json:"tampered"`
}

type StateZHAHumidity struct {
	Humidity int `json:"humidity"`
	//ISO 8601 timestamp
	Lastupdated string `json:"lastupdated"`
}

type StateZHALightLevel struct {
	Lux int `json:"lux"`
	// ISO 8601 timestamp
	Lastupdated string `json:"lastupdated"`
	Lightlevel  int    `json:"lightlevel"`
	Dark        bool   `json:"dark"`
	Daylight    bool   `json:"daylight"`
}

type StateZHAOpenClose struct {
	// ISO 8601 timestamp
	Lastupdated string `json:"lastupdated"`
	Lowbattery  bool   `json:"lowbattery"`
	Open        bool   `json:"open"`
	Tampered    bool   `json:"tampered"`
}

type StateZHAPower struct {
	Current int `json:"current"`
	// ISO 8601 timestamp
	Lastupdated string `json:"lastupdated"`
	Power       int    `json:"power"`
	Voltage     int    `json:"voltage"`
}

type StateZHAPresence struct {
	// ISO 8601 timestamp
	Lastupdated string `json:"lastupdated"`
	Lowbattery  bool   `json:"lowbattery"`
	Presence    bool   `json:"presence"`
	Tampered    bool   `json:"tampered"`
}

type StateZHASwitch struct {
	// Refer to https://dresden-elektronik.github.io/deconz-rest-doc/endpoints/sensors/button_events for device specific
	// values.
	Buttonevent int `json:"buttonevent"`

	// ISO 8601 timestamp
	Lastupdated   string `json:"lastupdated"`
	Gesture       int    `json:"gesture"`
	Eventduration int    `json:"eventduration"`
	X             int    `json:"x"`
	Y             int    `json:"y"`
	Angle         int    `json:"angle"`
}

type StateZHAPressure struct {
	Pressure int `json:"pressure"`
	// ISO 8601 timestamp
	Lastupdated string `json:"lastupdated"`
}

type StateZHATemperature struct {
	Temperature int `json:"temperature"`
	// ISO 8601 timestamp
	Lastupdated string `json:"lastupdated"`
}

type StateZHATime struct {
	Lastset string `json:"lastset"`
	// ISO 8601 timestamp
	Lastupdated string `json:"lastupdated"`
	// ISO 8601 timestamp
	Localtime string `json:"localtime"`
	Utc       uint   `json:"utc"`
}

type StateZHAThermostat struct {
	On               bool   `json:"on"`
	Errorcode        string `json:"errorcode"`
	Fanmode          string `json:"fanmode"`
	Floortemperature int    `json:"floortemperature"`
	Heating          bool   `json:"heating"`
	// ISO 8601 timestamp
	Lastupdated        string `json:"lastupdated"`
	Mountingmodeactive bool   `json:"mountingmodeactive"`
	Temperature        int    `json:"temperature"`
	Valve              int    `json:"valve"`
	Windowopen         string `json:"windowopen"`
}

type StateZHAVibration struct {
	Vibration bool `json:"vibration"`
	// ISO 8601 timestamp
	Lastupdated string `json:"lastupdated"`
	// Array of 3 Numbers
	Orientation       []int `json:"orientation"`
	Tiltangle         int   `json:"tiltangle"`
	Vibrationstrength int   `json:"vibrationstrength"`
}

type StateZHAWater struct {
	Water bool `json:"water"`
	// ISO 8601 timestamp
	Lastupdated string `json:"lastupdated"`
	Lowbattery  bool   `json:"lowbattery"`
	Tampered    bool   `json:"tampered"`
}
