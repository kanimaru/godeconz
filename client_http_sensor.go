package godeconz

type SensorRequestCreateConfig struct {
	// (default: true)
	On bool `json:"on"`
	// (default: true)
	Reachable bool `json:"reachable"`
	// (0–100)
	Battery uint8 `json:"battery"`
}

type SensorRequestCreate[T any] struct {
	// Name The name of the sensor.	required
	Name string `json:"name"`
	// Modelid The model identifier of the sensor.	required
	Modelid string `json:"modelid"`
	// The software version of the sensor.	required
	Swversion string `json:"swversion"`
	// The type of the sensor, see Sensor types and states.	required
	Type string `json:"type"`
	// The unique id of the sensor. Should be the MAC address of the device.	required
	Uniqueid string `json:"uniqueid"`
	// The manufacturer name of the sensor.	required
	Manufacturername string `json:"manufacturername"`
	// The state of the sensor, see Sensor types and states.	optional
	State *T `json:"state,omitempty"`
	// The config of the sensor. optional
	Config *SensorRequestCreateConfig `json:"config,omitempty"`
}

// SensorMode (only available for dresden elektronik Lighting Switch)
type SensorMode int

const SensorModeScene SensorMode = 1
const SensorModeTwoGroup SensorMode = 2
const SensorModeTemperature SensorMode = 3

type SensorResponse[C, S any] struct {
	// The config of the sensor. Refer to Change sensor config for further details.
	Config C `json:"config"`
	// The Endpoint of the sensor.
	Ep int `json:"ep"`
	// HTTP etag which changes whenever the sensor changes.
	Etag string `json:"etag"`
	// (ISO 8601 timestamp) Timestamp representing the last time a message from the sensor was received. UTC with resolution of minutes.
	Lastseen string `json:"lastseen"`
	// The manufacturer name of the sensor.
	Manufacturername string `json:"manufacturername"`
	// The mode of the sensor.
	Mode SensorMode `json:"mode"`
	// The model id of the sensor.
	Modelid string `json:"modelid"`
	// The name of the sensor.
	Name string `json:"name"`
	// The state of the sensor.
	State S `json:"state"`
	// Software version of the sensor.
	Swversion string `json:"swversion"`
	// The type of the sensor.
	Type string `json:"type"`
	// The unique identifiers including the MAC address of the sensor.
	Uniqueid string `json:"uniqueid"`
}

type SensorRequestUpdate struct {
	// The name of the sensor. optional
	Name *string `json:"name,omitempty"`
	// Only available for dresden elektronik Lighting Switch. Set the mode of the switch. optional
	Mode *SensorMode `json:"mode,omitempty"`
}

// CreateSensor Creates a new CLIP sensor.
func (c *Client[R]) CreateSensor(create SensorRequestCreate[any]) (R, error) {
	return c.Post("/sensors", create)
}

// GetAllSensors Returns a list of all sensors. If there are no sensors in the system an empty object {} is
// returned.
func (c *Client[R]) GetAllSensors(container *map[string]SensorResponse[any, any]) (R, error) {
	return c.Get("/sensors", container)
}

// GetSensor Returns the sensor with the specified id.
func (c *Client[R]) GetSensor(id string, container *SensorResponse[any, any]) (R, error) {
	return c.Get("/sensors/%s", container, id)
}

// UpdateSensor Update a sensor with the specified parameters.
func (c *Client[R]) UpdateSensor(id string, data SensorRequestUpdate) (R, error) {
	return c.Put("/sensors/%s", data, id)
}

// ChangeSensorConfig Update a sensor config with the specified parameters.
// Sensors expose certain configuration parameters depending on their defined or known capabilities.
// To get an overview on which parameters are available for a particular device, get the sensor state of either
// all (https://dresden-elektronik.github.io/deconz-rest-doc/endpoints/sensors/#getall)
// or a single sensor https://dresden-elektronik.github.io/deconz-rest-doc/endpoints/sensors/#getsensor.
func (c *Client[R]) ChangeSensorConfig(id string, config interface{}) (R, error) {
	return c.Put("/sensors/%s/config", config, id)
}

// ChangeSensorState Update a sensor state with the specified parameters.
// Changing the sensor state is only allowed for CLIP sensors.
func (c *Client[R]) ChangeSensorState(id string, state interface{}) (R, error) {
	return c.Put("/sensors/%s/state", state, id)
}

// DeleteSensor Delete a sensor.
func (c *Client[R]) DeleteSensor(id string) (R, error) {
	return c.Delete("/sensors/%s", nil, id)
}
