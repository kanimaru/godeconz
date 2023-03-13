package http

type LightAlertMode string

// LightAlertNone light is not performing an alert
const LightAlertNone LightAlertMode = "none"

// LightAlertSelect light is blinking a short time
const LightAlertSelect LightAlertMode = "select"

// LightAlertLselect light is blinking a longer time
const LightAlertLselect LightAlertMode = "lselect"

type ColorMode string

// ColorModeHS hue and saturation
const ColorModeHS ColorMode = "hs"

// ColorModeXY CIE xy values
const ColorModeXY ColorMode = "xy"

// ColorModeCT color temperature
const ColorModeCT ColorMode = "ct"

type LightEffectMode string

// LightEffectNone no effect
const LightEffectNone LightEffectMode = "none"

// LightEffectColorLoop cycle through hue values 0–360
const LightEffectColorLoop LightEffectMode = "colorloop"

type LightResponseStateDetail struct {
	// true if the light is on.
	On *bool `json:"on"`
	// (0–255)	Brightness of the light. Depending on the light type 0 might not mean visible "off" but minimum
	// brightness.
	Bri *int `json:"bri"`
	// (0–65535)	Color hue of the light. The hue parameter in the HSV color model is between 0°–360° and is mapped
	// to 0–65535 to get 16-bit resolution.
	Hue *int `json:"hue"`
	// (0–255)	Color saturation of the light. There 0 means no color at all and 255 is the greatest saturation of
	// the color.
	Sat *int `json:"sat"`
	// (ctmin–ctmax)	Mired color temperature of the light.	Where Mired is 1000000 / color temperature (in kelvins).
	Ct *int `json:"ct"`
	// CIE xy color space coordinates as array [x, y] of real values (0–1).
	Xy *[]float32 `json:"xy"`
	// Temporary alert effect.
	Alert *LightAlertMode `json:"alert"`
	// The current color mode of the light
	Colormode *ColorMode `json:"colormode"`
	// Effect of the light
	Effect *LightEffectMode `json:"effect"`
	// (0–6)	SETTABLE. Sets the speed of fans/ventilators.
	Speed *uint8 `json:"speed"`
	// true if the light is reachable and accepts commands.
	Reachable *bool `json:"reachable"`
}

func (l LightResponseStateDetail) IsOn() bool {
	if l.On != nil {
		return *l.On
	} else {
		return false
	}
}

func (l LightResponseStateDetail) GetAlert() interface{} {
	return l.Alert
}

func (l LightResponseStateDetail) GetBrightness() int {
	if l.Bri != nil {
		return *l.Bri
	} else {
		return 0
	}
}

type LightResponseState[STATE any] struct {

	// The color capabilities as reported by the light.
	Colorcapabilities *int `json:"colorcapabilities"`
	// The maximum mired color temperature value a device supports.
	Ctmax *int `json:"ctmax"`
	// The minimum mired color temperature value a device supports.
	Ctmin *int `json:"ctmin"`
	// Last time the device announced itself to the network.
	Lastannounced string `json:"lastannounced"`
	// Last time the device has transmitted any data.
	Lastseen string `json:"lastseen"`
	// HTTP etag which changes on any action to the light.
	Etag string `json:"etag"`
	// Indicates if the light can change color. Deprecated - use state instead: if light has no color colormode, hue
	// and xy will not be shown.
	Hascolor *bool `json:"hascolor"`
	// The manufacturer of the light device.
	Manufacturer string `json:"manufacturer"`
	// Name of a light.
	Name string `json:"name"`
	// An identifier unique to the product.
	Modelid string `json:"modelid"`
	// Not used in the current version.
	Pointsymbol *interface{} `json:"pointsymbol"`
	// SETTABLE. Brightness to set after power on (limited to DE devices).
	Powerup *int `json:"powerup"`
	// Firmware version.
	Swversion string `json:"swversion"`
	// Human-readable type of the light.
	Type string `json:"type"`
	// The current state of the light.
	State STATE `json:"state"`
	// The unique id of the light. It consists of the MAC address of the light followed by a dash and a unique endpoint
	// identifier in the range 01 to FF.
	Uniqueid string `json:"uniqueid"`
}

type LightRequestState struct {
	// Trigger a temporary alert effect. optional
	Alert LightAlertMode `json:"alert,omitempty"`
	// Number (0–255)	Set the brightness of the light. Depending on the light type 0 might not mean visible "off"
	// but minimum brightness. If the light is off and the value is greater 0 a on=true shall also be provided. optional
	Bri *uint8 `json:"bri,omitempty"`
	// Specifies the speed of a colorloop (default: 15). 1 = very fast 255 = very slow This parameter only has an
	// effect when it is called together with effect colorloop. optional
	Colorloopspeed *uint8 `json:"colorloopspeed,omitempty"`
	// (ctmin–ctmax)	Set the Mired color temperature of the light. Where Mired is 1000000 / color temperature
	// (in kelvins).	optional
	Ct *int `json:"ct,omitempty"`
	// Trigger an effect of the light. optional
	Effect LightEffectMode `json:"effect,omitempty"`
	// (0–65535)	Set the color hue of the light. The hue parameter in the HSV color model is between 0°–360° and is
	// mapped to 0–65535 to get 16-bit resolution.	optional
	Hue *uint16 `json:"hue,omitempty"`
	// Set to true to turn the light on, false to turn it off.	optional
	On *bool `json:"on,omitempty"`
	// Number (0–255)	Set the color saturation of the light. There 0 means no color at all and 255 is the greatest
	// saturation of the color.	optional
	Sat *uint8 `json:"sat,omitempty"`
	// Transition time in 1/10 seconds between two states. Note that not all states support a transition time.
	// For example, a transition time when setting on will be ignored as the Zigbee On and Off commands do not support
	// transition times. In general, light attributes that support a range of values support transition times, while
	// boolean values do not.	optional
	Transitiontime *uint `json:"transitiontime,omitempty"`
	// Set the CIE xy color space coordinates as array [x, y] of real values (0–1).	optional
	Xy []float32 `json:"xy,omitempty"`
}

type LightRequestAttribute struct {
	name string
}

type LightRequestDelete struct {
	reset bool
}

// GetAllLights Returns a list of all lights.
func (c *Client[R]) GetAllLights(container *map[string]LightResponseState[LightResponseStateDetail]) (R, error) {
	return c.Get("/lights", container)
}

// GetLightState is almost the same as Client.GetLightState but you can add generics for the response
func GetLightState[R any, S any](c *Client[R], id string, container *LightResponseState[S]) (R, error) {
	return c.Get("/lights/%s", container, id)
}

// GetLightState Returns the full state of a light.
func (c *Client[R]) GetLightState(id string, container *LightResponseState[LightResponseStateDetail]) (R, error) {
	return c.Get("/lights/%s", container, id)
}

// SetLightState Sets the state of a light.
func (c *Client[R]) SetLightState(id string, state LightRequestState) (R, error) {
	return c.Put("/lights/%s/state", state, id)
}

// SetLightAttributes Sets attributes of a light which are not related to its state.
func (c *Client[R]) SetLightAttributes(id string, attr LightRequestAttribute) (R, error) {
	return c.Put("/lights/%s", attr, id)
}

// DeleteLight Removes the light from the gateway. It will not be shown in any REST-API call. Also deletes all
// groups and scenes on the light device.
func (c *Client[R]) DeleteLight(id string, reset LightRequestDelete) (R, error) {
	return c.Delete("/lights/%s", reset, id)
}

// RemoveFromAllGroups Remove the light from all groups it is a member of.
func (c *Client[R]) RemoveFromAllGroups(id string) (R, error) {
	return c.Delete("/lights/%s/groups", nil, id)
}

// RemoveFromAllScenes Remove the light from all scenes it is a member of.
func (c *Client[R]) RemoveFromAllScenes(id string) (R, error) {
	return c.Delete("/lights/%s/scenes", nil, id)
}
