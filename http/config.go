package http

type Config struct {
	// (default: true)
	On bool `json:"on"`
	// (default: true)
	Reachable bool `json:"reachable"`
	// (0–100)
	Battery uint8 `json:"battery"`
}

// ConfigZHAPressure Not defined by API added for convenience
type ConfigZHAPressure struct {
	Config
}

// ConfigZHAOpenClose Not defined by API added for convenience
type ConfigZHAOpenClose struct {
	Config
}

type ConfigDaylight struct {
	Config
	// 	True if the daylight sensor is configured with coordinates.	R
	Daylight bool `json:"daylight"`
	// 	Latitude of the set location/timezone.	W
	Lat float32 `json:"lat"`
	// 	Longitude of the set location/timezone.	W
	Long float32 `json:"long"`
	// (-120–120)	Sunrise offset value for location/timezone in minutes.	R
	Sunriseoffset int8 `json:"sunriseoffset"`
	// (-120–120)	Sunset offset value for location/timezone in minutes.	R
	Sunsetoffset int8 `json:"sunsetoffset"`
}

type ConfigZHALightLevel struct {
	Config
	// Specifies at which lightlevel the dark attribute turns false. Default: 12000 RW
	Tholddark uint16 `json:"tholddark"`

	// Relative offset to tholddark. Sets daylight attribute to true when lightlevel is at or above
	// (tholddark + tholdoffset). Default: 7000 RW
	Tholdoffset uint16 `json:"tholdoffset"`
}

type ConfigZHAHumidity struct {
	Config
	// (-32768–32767)	Adds a signed offset value to measured state values. Values send by the REST-API are already
	// amended by the offset. RW
	Offset int16 `json:"offset"`
}

type ConfigZHAPresence struct {
	Config
	// (0–65535)	Timeout in seconds presence state is set to false again. RW
	Duration uint16 `json:"duration"`
	// (0–65535) The occupied to unoccupied delay in seconds. RW
	Delay uint16 `json:"delay"`
}

type ConfigZHATemperature struct {
	Config
	// Adds a signed offset value to measured state values. Values send by the REST-API are already amended by
	// the offset.	R
	Offset int16 `json:"offset"`
}

type ZHASwitchMode string

const ZHASwitchModeMomentary ZHASwitchMode = "momentary"
const ZHASwitchModeRocker ZHASwitchMode = "rocker"

type ConfigZHASwitch struct {
	Config
	// The associated Zigbee group the sensor controls. (only supported by some sensors) R
	Group uint16 `json:"group"`

	// For ubisys S1/S2, operation mode of the switch. RW
	Mode string `json:"mode"`
}

type ZHAThermostatMode string

const ZHAThermostatModeOff ZHAThermostatMode = "off"
const ZHAThermostatModeAuto ZHAThermostatMode = "auto"
const ZHAThermostatModeCool ZHAThermostatMode = "cool"
const ZHAThermostatModeHeat ZHAThermostatMode = "heat"
const ZHAThermostatModeEmergencyHeating ZHAThermostatMode = "emergency heating"
const ZHAThermostatModePrecooling ZHAThermostatMode = "precooling"
const ZHAThermostatModeFanOnly ZHAThermostatMode = "fan only"
const ZHAThermostatModeDry ZHAThermostatMode = "dry"
const ZHAThermostatModeSleep ZHAThermostatMode = "sleep"

type ZHAThermostatFanMode string

const ZHAThermostatFanModeOff ZHAThermostatFanMode = "off"
const ZHAThermostatFanModeLow ZHAThermostatFanMode = "low"
const ZHAThermostatFanModeMedium ZHAThermostatFanMode = "medium"
const ZHAThermostatFanModeHigh ZHAThermostatFanMode = "high"
const ZHAThermostatFanModeOn ZHAThermostatFanMode = "on"
const ZHAThermostatFanModeAuto ZHAThermostatFanMode = "auto"
const ZHAThermostatFanModeSmart ZHAThermostatFanMode = "smart"

type ZHAThermostatHostflags int32

const ZHAThermostatHostflagsDisplayFlipped ZHAThermostatHostflags = 0x0002
const ZHAThermostatHostflagsModeHeat ZHAThermostatHostflags = 0x0004
const ZHAThermostatHostflagsModeOff ZHAThermostatHostflags = 0x0010
const ZHAThermostatHostflagsLocked ZHAThermostatHostflags = 0x0080

type ZHAThermostatPreset string

const ZHAThermostatPresetHoliday ZHAThermostatPreset = "holiday"
const ZHAThermostatPresetAuto ZHAThermostatPreset = "auto"
const ZHAThermostatPresetManual ZHAThermostatPreset = "manual"
const ZHAThermostatPresetComfort ZHAThermostatPreset = "comfort"
const ZHAThermostatPresetEco ZHAThermostatPreset = "eco"
const ZHAThermostatPresetBoost ZHAThermostatPreset = "boost"
const ZHAThermostatPresetComplex ZHAThermostatPreset = "complex"

type ZHAThermostatSwingmode string

const ZHAThermostatSwingmodeFullyClosed ZHAThermostatSwingmode = "fully closed"
const ZHAThermostatSwingmodeFullyOpen ZHAThermostatSwingmode = "fully open"
const ZHAThermostatSwingmodeQuarterOpen ZHAThermostatSwingmode = "quarter open"
const ZHAThermostatSwingmodeHalfOpen ZHAThermostatSwingmode = "half open"
const ZHAThermostatSwingmodeThreeQuartersOpen ZHAThermostatSwingmode = "three quarters open"

type ZHAThermostatTemperatureMeasurement string

const ZHAThermostatTemperatureMeasurementAirSensor ZHAThermostatTemperatureMeasurement = "air sensor"
const ZHAThermostatTemperatureMeasurementFloorSensor ZHAThermostatTemperatureMeasurement = "floor sensor"
const ZHAThermostatTemperatureMeasurementFloorProtection ZHAThermostatTemperatureMeasurement = "floor protection"

type ConfigZHAThermostat struct {
	Config
	// Sets the current operating mode of a thermostat. (Supported modes are device dependent) RW
	Mode ZHAThermostatMode `json:"mode"`
	// Flip the display for TRVs supporting it. RW
	Displayflipped bool `json:"displayflipped"`
	// (-32768–32767) Allows to use the temperature value provided by an external sensor.
	// (device dependent and only exposed for devices supporting it) RW
	Externalsensortemp int16 `json:"externalsensortemp"`
	// Allows to use the open/close state from an external sensor.
	// (device dependent and only exposed for devices supporting it) RW
	Externalwindowopen bool `json:"externalwindowopen"`
	// Sets the mode of the fan. (device dependent and only exposed for devices supporting it) RW
	Fanmode ZHAThermostatFanMode `json:"fanmode"`
	// Eurotronic Spirit SPZB Only for debugging purpose. R
	Hostflags ZHAThermostatHostflags `json:"hostflags"`
	// Child lock active/inactive for thermostats/TRVs supporting it. RW
	Locked bool `json:"locked"`
	// Sets the operating mode for Tuya thermostats. (supported modes are device dependent) RW
	Preset ZHAThermostatPreset `json:"preset"`
	// Controls valve of thermostats. false — Close valve | true — Open valve (exposed for thermostats supporting it) RW
	Setvalve bool `json:"setvalve"`
	// Sets the AC louvers position. (exposed for thermostats supporting it) RW
	Swingmode ZHAThermostatSwingmode `json:"swingmode"`
	// Sets the mode of operation for Elko Super TR thermostat. RW
	TemperatureMeasurement ZHAThermostatTemperatureMeasurement `json:"temperature measurement"`
	// Sets if window open detection shall be active or inactive for Tuya thermostats. (support is device dependent) RW
	WindowOpenSet bool `json:"window open_set"`
	// A thermostat schedule. RW TODO not described in the https://dresden-elektronik.github.io/deconz-rest-doc/endpoints/sensors/#dev-sensor-config-attr
	Schedule []interface{} `json:"schedule"`
	// True when a thermostat schedule is enabled. RW
	ScheduleOn bool `json:"schedule_on"`
	// (700–3500)	Set the desired cooling temperature. RW
	Coolsetpoint uint8 `json:"coolsetpoint"`
	// (500–3200)	Set the desired heating temperature. RW
	Heatsetpoint uint8 `json:"heatsetpoint"`
	// 	Sets a TRV into mounting mode if supported (valve fully open position). RW
	Mountingmode bool `json:"mountingmode"`
}
