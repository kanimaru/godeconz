package godeconz

import (
	"encoding/json"
	"errors"
	"maps"
	"strings"
	"time"
)

// SensorType is the deCONZ sensor type ("type" attribute).
type SensorType = string

// Common sensor types. One physical device usually exposes several sensors, e.g. an Aqara
// weather sensor appears as ZHATemperature, ZHAHumidity and ZHAPressure sharing a MAC.
const (
	SensorTemperature    SensorType = "ZHATemperature"
	SensorHumidity       SensorType = "ZHAHumidity"
	SensorPressure       SensorType = "ZHAPressure"
	SensorLightLevel     SensorType = "ZHALightLevel"
	SensorPresence       SensorType = "ZHAPresence"
	SensorOpenClose      SensorType = "ZHAOpenClose"
	SensorPower          SensorType = "ZHAPower"
	SensorConsumption    SensorType = "ZHAConsumption"
	SensorBattery        SensorType = "ZHABattery"
	SensorSwitch         SensorType = "ZHASwitch"
	SensorWater          SensorType = "ZHAWater"
	SensorFire           SensorType = "ZHAFire"
	SensorCarbonMonoxide SensorType = "ZHACarbonMonoxide"
	SensorVibration      SensorType = "ZHAVibration"
	SensorAlarm          SensorType = "ZHAAlarm"
	SensorAirQuality     SensorType = "ZHAAirQuality"
	SensorThermostat     SensorType = "ZHAThermostat"
	SensorDaylight       SensorType = "Daylight"
)

// Sensor is a sensor resource as returned by GET /sensors.
type Sensor struct {
	// ID is the REST resource id (the key in GET /sensors). It changes when a device is
	// re-paired; UniqueID does not.
	ID               string       `json:"-"`
	Name             string       `json:"name"`
	Type             SensorType   `json:"type"`
	ModelID          string       `json:"modelid"`
	ManufacturerName string       `json:"manufacturername"`
	UniqueID         string       `json:"uniqueid"`
	SWVersion        string       `json:"swversion"`
	ETag             string       `json:"etag"`
	LastSeen         string       `json:"lastseen"`
	LastAnnounced    string       `json:"lastannounced"`
	Ep               int          `json:"ep"`
	Mode             *int         `json:"mode,omitempty"`
	State            SensorState  `json:"state"`
	Config           SensorConfig `json:"config"`
}

// SensorState covers the state attributes of all common sensor types. Attributes the sensor
// does not report are nil. Raw keeps every attribute, including ones without a field here.
// Values are in deCONZ units; use the unit methods for SI values.
type SensorState struct {
	LastUpdated string `json:"lastupdated,omitempty"`

	Temperature *int `json:"temperature,omitempty"` // 0.01 °C
	Humidity    *int `json:"humidity,omitempty"`    // 0.01 %
	Pressure    *int `json:"pressure,omitempty"`    // hPa

	Lux        *int  `json:"lux,omitempty"`
	LightLevel *int  `json:"lightlevel,omitempty"` // 10000·log10(lux)+1
	Dark       *bool `json:"dark,omitempty"`
	Daylight   *bool `json:"daylight,omitempty"`

	Presence *bool `json:"presence,omitempty"`
	Open     *bool `json:"open,omitempty"`

	Power       *int   `json:"power,omitempty"`       // W
	Consumption *int64 `json:"consumption,omitempty"` // Wh
	Voltage     *int   `json:"voltage,omitempty"`     // V
	Current     *int   `json:"current,omitempty"`     // mA

	ButtonEvent   *int `json:"buttonevent,omitempty"`
	Gesture       *int `json:"gesture,omitempty"`
	EventDuration *int `json:"eventduration,omitempty"`

	Battery        *int  `json:"battery,omitempty"` // %, ZHABattery
	LowBattery     *bool `json:"lowbattery,omitempty"`
	Tampered       *bool `json:"tampered,omitempty"`
	Water          *bool `json:"water,omitempty"`
	Fire           *bool `json:"fire,omitempty"`
	CarbonMonoxide *bool `json:"carbonmonoxide,omitempty"`
	Vibration      *bool `json:"vibration,omitempty"`
	Alarm          *bool `json:"alarm,omitempty"`

	AirQuality    *string `json:"airquality,omitempty"`
	AirQualityPPB *int    `json:"airqualityppb,omitempty"`
	PM25          *int    `json:"pm2_5,omitempty"`

	Raw map[string]json.RawMessage `json:"-"`
}

// UnmarshalJSON decodes the known attributes and keeps all of them in Raw. An attribute with an
// unexpected JSON type (a device reporting a float where deCONZ normally sends an integer) is
// left nil instead of failing the whole sensor; its value is still in Raw.
func (s *SensorState) UnmarshalJSON(data []byte) error {
	type plain SensorState
	decoded, raw, err := decodeLenient[plain](data)
	if err != nil {
		return err
	}
	decoded.Raw = raw
	*s = SensorState(decoded)
	return nil
}

// TemperatureCelsius returns the temperature in °C.
func (s SensorState) TemperatureCelsius() (float64, bool) {
	return scaled(s.Temperature, 0.01)
}

// HumidityPercent returns the relative humidity in %.
func (s SensorState) HumidityPercent() (float64, bool) {
	return scaled(s.Humidity, 0.01)
}

// PressureHPa returns the air pressure in hPa.
func (s SensorState) PressureHPa() (float64, bool) {
	return scaled(s.Pressure, 1)
}

// PowerWatts returns the current power draw in W.
func (s SensorState) PowerWatts() (float64, bool) {
	return scaled(s.Power, 1)
}

// EnergyWattHours returns the energy counter in Wh.
func (s SensorState) EnergyWattHours() (float64, bool) {
	if s.Consumption == nil {
		return 0, false
	}
	return float64(*s.Consumption), true
}

// VoltageVolts returns the voltage in V.
func (s SensorState) VoltageVolts() (float64, bool) {
	return scaled(s.Voltage, 1)
}

// CurrentAmperes returns the current in A.
func (s SensorState) CurrentAmperes() (float64, bool) {
	return scaled(s.Current, 0.001)
}

// LastUpdatedTime returns LastUpdated as time.
func (s SensorState) LastUpdatedTime() (time.Time, bool) {
	return ParseTime(s.LastUpdated)
}

// decodeLenient decodes a JSON object into T and also returns all of its attributes. An
// attribute whose JSON type does not fit its field is dropped and decoding repeated, so the
// field stays nil. (encoding/json would leave such a pointer allocated with value 0, which reads
// like a real measurement.)
func decodeLenient[T any](data []byte) (T, map[string]json.RawMessage, error) {
	var zero T
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return zero, nil, err
	}
	fields := maps.Clone(raw)
	for {
		body, err := json.Marshal(fields)
		if err != nil {
			return zero, raw, err
		}
		var decoded T
		err = json.Unmarshal(body, &decoded)
		var typeErr *json.UnmarshalTypeError
		if err == nil || !errors.As(err, &typeErr) {
			return decoded, raw, err
		}
		key, _, _ := strings.Cut(typeErr.Field, ".")
		if _, ok := fields[key]; !ok {
			return decoded, raw, nil
		}
		delete(fields, key)
	}
}

func scaled(v *int, factor float64) (float64, bool) {
	if v == nil {
		return 0, false
	}
	return float64(*v) * factor, true
}

// SensorConfig covers the config attributes of common sensor types. Attributes the sensor does
// not have are nil; Raw keeps every attribute. All fields are pointers with omitempty, so a
// SensorConfig with only the fields to change can be sent as update.
type SensorConfig struct {
	On        *bool `json:"on,omitempty"`
	Reachable *bool `json:"reachable,omitempty"`
	Battery   *int  `json:"battery,omitempty"` // %

	Offset         *int  `json:"offset,omitempty"`      // temperature/humidity calibration
	TholdDark      *int  `json:"tholddark,omitempty"`   // light level below which "dark" is true
	TholdOffset    *int  `json:"tholdoffset,omitempty"` // added to TholdDark for "daylight"
	Duration       *int  `json:"duration,omitempty"`    // presence timeout in s
	Delay          *int  `json:"delay,omitempty"`       // presence delay in s
	Sensitivity    *int  `json:"sensitivity,omitempty"`
	SensitivityMax *int  `json:"sensitivitymax,omitempty"`
	LEDIndication  *bool `json:"ledindication,omitempty"`

	HeatSetpoint *int    `json:"heatsetpoint,omitempty"` // 0.01 °C
	CoolSetpoint *int    `json:"coolsetpoint,omitempty"` // 0.01 °C
	Mode         *string `json:"mode,omitempty"`
	Locked       *bool   `json:"locked,omitempty"`

	SunriseOffset *int `json:"sunriseoffset,omitempty"`
	SunsetOffset  *int `json:"sunsetoffset,omitempty"`

	Raw map[string]json.RawMessage `json:"-"`
}

// UnmarshalJSON decodes the known attributes and keeps all of them in Raw. Attributes with an
// unexpected JSON type are left nil, as for SensorState.
func (c *SensorConfig) UnmarshalJSON(data []byte) error {
	type plain SensorConfig
	decoded, raw, err := decodeLenient[plain](data)
	if err != nil {
		return err
	}
	decoded.Raw = raw
	*c = SensorConfig(decoded)
	return nil
}

// IsReachable reports whether the sensor is reachable; unknown counts as reachable.
func (c SensorConfig) IsReachable() bool {
	return c.Reachable == nil || *c.Reachable
}

// SensorAttributesUpdate is the body of PUT /sensors/{id}.
type SensorAttributesUpdate struct {
	Name *string `json:"name,omitempty"`
	Mode *int    `json:"mode,omitempty"`
}

// SensorCreate is the body of POST /sensors, which creates a CLIP (software) sensor.
type SensorCreate struct {
	Name             string        `json:"name"`
	ModelID          string        `json:"modelid"`
	SWVersion        string        `json:"swversion"`
	Type             SensorType    `json:"type"`
	UniqueID         string        `json:"uniqueid"`
	ManufacturerName string        `json:"manufacturername"`
	State            *SensorState  `json:"state,omitempty"`
	Config           *SensorConfig `json:"config,omitempty"`
}
