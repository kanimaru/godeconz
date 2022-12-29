package http

type GroupRequestCreate struct {
	Name string `json:"name"`
}

type GroupResponseCreate struct {
	Success struct {
		// The unique identifier of the group.
		Id string `json:"id"`
	} `json:"success"`
}

type GroupResponse struct {
	// If this group was created by a device (switch or sensor) this list contains the device ids.
	Devicemembership []string `json:"devicemembership"`
	// HTTP etag which changes on any action to the group.
	Etag string `json:"etag"`
	// Indicates if this group is hidden.
	Hidden bool `json:"hidden"`
	// Name of a group.
	Name string `json:"name"`
}

type GroupEffectMode string

// GroupEffectNone no effect
const GroupEffectNone GroupEffectMode = "none"

// GroupEffectColorLoop the lights of the group will cycle continously through all colors with the speed specified by colorloopspeed
const GroupEffectColorLoop GroupEffectMode = "colorloop"

type GroupResponseAttribute struct {
	// The last action which was sent to the group.
	action struct {
		// true if the group was turned on.
		On bool `json:"on"`
		// (0–255)	Brightness of the group. Depending on the lights 0 might not mean visible "off" but minimum brightness.
		Bri uint8 `json:"bri"`
		// (0–65535)	The hue parameter in the HSV color model is between 0°–360° and is mapped to 0–65535 to get 16-bit resolution.
		Hue uint16 `json:"hue"`
		// (0–255)	Color saturation there 0 means no color at all and 255 is the greatest saturation of the color.
		Sat uint8 `json:"sat"`
		// (153–500)	Mired color temperature. (2000K–6500K)
		Ct uint16 `json:"ct"`
		// CIE xy color space coordinates as array [x, y] of real values (0–1).
		Xy []float32 `json:"xy"`
		// Dynamic effect:
		Effect GroupEffectMode `json:"effect"`
	}
	// A list of device ids (sensors) if this group was created by a device.
	Devicemembership []string `json:"devicemembership"`
	// HTTP etag which changes on any action to the group.
	Etag string `json:"etag"`
	// Indicates the hidden status of the group. Has no effect at the gateway but apps can use this to hide groups.
	Hidden bool `json:"hidden"`
	// The id of the group.
	Id string `json:"id"`
	// A list of all light ids of this group. Sequence is defined by the gateway.
	Lights []string `json:"lights"`
	// A list of light ids of this group that can be sorted by the user. Need not contain all light ids of this group.
	Lightsequence []string `json:"lightsequence"`
	// A list of light ids of this group that are subsequent ids from multi devices with multiple endpoints like the FLS-PP.
	Mulitdeviceids []string `json:"mulitdeviceids"`
	// Name of the group.
	Name string `json:"name"`
	// A list of scenes of the group.
	Scenes []struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	}
}

type GroupRequestAttribute struct {
	// (0–32)	The name of the group	optional
	Name string `json:"name,omitempty"`
	// IDs of the lights which are members of the group.	optional
	Lights []string `json:"lights,omitempty"`
	// Indicates the hidden status of the group. Has no effect at the gateway but apps can use this to hide groups.	optional
	Hidden *bool `json:"hidden,omitempty"`
	// Specify a sorted list of light ids that can be used in apps.	optional
	Lightsequence []string `json:"lightsequence,omitempty"`
	// Append the subsequential light ids of multi-devices like the FLS-PP if the app should handle that light differently.	optional
	Mulitdeviceids []string `json:"mulitdeviceids,omitempty"`
}

type GroupAlertMode string

// GroupAlertModeNone lights are not performing an alert
const GroupAlertModeNone GroupAlertMode = "none"

// GroupAlertModeSelect lights are blinking a short time
const GroupAlertModeSelect GroupAlertMode = "select"

// GroupAlertModeLselect lights are blinking a longer time
const GroupAlertModeLselect GroupAlertMode = "lselect"

type GroupRequestState struct {
	// Set to true to turn the lights on, false to turn them off.	optional
	On *bool `json:"on"`
	// Set to true toggles the lights of that group from on to off or vice versa, false has no effect. **Notice:** This setting supersedes the `on` parameter!	optional
	Toggle *bool `json:"toggle"`
	// (0–255)	Set the brightness of the group. Depending on the lights 0 might not mean visible "off" but minimum brightness. If the lights are off and the value is greater 0 a on=true shall also be provided.	optional
	Bri *uint8 `json:"bri"`
	// (0–65535)	Set the color hue of the group. The hue parameter in the HSV color model is between 0°–360° and is mapped to 0–65535 to get 16-bit resolution.	optional
	Hue *uint16 `json:"hue"`
	// (0–255)	Set the color saturation of the group. There 0 means no color at all and 255 is the highest saturation of the color.	optional
	Sat *uint8 `json:"sat"`
	// (153–500)	Set the Mired color temperature of the group. (2000K–6500K)	optional
	Ct *uint16 `json:"ct"`
	// Set the CIE xy color space coordinates as array [x, y] of real values (0–1).	optional
	Xy []float32 `json:"xy"`
	// Trigger a temporary alert effect. optional
	Alert GroupAlertMode `json:"alert"`
	// Trigger an effect of the group. optional
	Effect GroupEffectMode `json:"effect"`
	//  (1–255)	Specifies the speed of a colorloop. 1 = very fast, 255 = very slow (default: 15). This parameter only has an effect when it is called together with effect colorloop.	optional
	Colorloopspeed *uint8 `json:"colorloopspeed"`
	// Transition time in 1/10 seconds between two states. Note that not all states support a transition time.
	// For example, a transition time when setting on will be ignored as the Zigbee On and Off commands do not support
	//transition times. In general, light attributes that support a range of values support transition times, while
	//boolean values do not.	optional
	Transitiontime *uint `json:"transitiontime"`
}

// CreateGroup creates a new empty group. Creating a group with a name which already exists will not create a new group
// or fail. Such a call does only return the id of the existing group.
func (c *Client[R]) CreateGroup(create GroupRequestCreate, createResponse *GroupResponseCreate) (R, error) {
	return c.PostWithResult("/groups", create, createResponse)
}

// GetAllGroups returns a list of all groups.
func (c *Client[R]) GetAllGroups(groups *map[string]GroupResponse) (R, error) {
	return c.Get("/groups", groups)
}

// GetGroupAttributes returns the full state of a group.
func (c *Client[R]) GetGroupAttributes(id string, attribute *GroupResponseAttribute) (R, error) {
	return c.Get("/groups/%s", attribute, id)
}

// SetGroupAttributes of a group which are not related to its state. In order to add or remove lights to the
// group the lights must be powered on.
func (c *Client[R]) SetGroupAttributes(id string, attribute GroupRequestAttribute) (R, error) {
	return c.Put("/groups/%s", attribute, id)
}

// SetGroupState Sets the state of a group.
func (c *Client[R]) SetGroupState(id string, state GroupRequestState) (R, error) {
	return c.Put("/groups/%s/action", state, id)
}

// DeleteGroup Deletes a group. In order to delete the group and therefore remove all lights from the group the
// lights must be powered on.
func (c *Client[R]) DeleteGroup(id string) (R, error) {
	return c.Delete("/groups/%s", nil, id)

}
