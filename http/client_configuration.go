package http

import (
	"encoding/base64"
	"encoding/json"
	"strings"
)

type ConfigRequestAcquireAPIKey struct {
	// String (10–40 chars)	Will be used as username. If not specified a random key will be generated.	optional
	Username string `json:"username,omitempty"`
	// String (0–40 chars)	Name of the client application. required
	Devicetype string `json:"devicetype"`
}

// AcquireAPIKey Creates a new API key which provides authorized access to the REST-API.
// The request will only succeed if the gateway is unlocked or valid HTTP basic authentication
// credentials are provided in the HTTP request header see authorization. (TODO support basic auth also)
func (c *Client[R]) AcquireAPIKey(request ConfigRequestAcquireAPIKey) (R, error) {
	path := c.getBasePath("/api")
	return c.Post(path, request)
}

// DeleteAPIKey Deletes an API key, so it can no longer be used.
func (c *Client[R]) DeleteAPIKey(apiKeyToDelete string) (R, error) {
	return c.Delete("/config/whitelist/%s", nil, apiKeyToDelete)
}

type WhitelistedAPIKey struct {
	CreateDate  string `json:"create date"`
	LastUseDate string `json:"last use date"`
	Name        string `json:"name"`
}

type ConfigResponse struct {
	// 	The version of the deCONZ Rest API
	ApiVersion string `json:"apiversion"`
	// The unique identifier for the gateway.
	BridgeId string `json:"bridgeid"`
	// The product name of the gateway. Valid values are "ConBee", "RaspBee", "ConBee II" and "RaspBee II".
	DeviceName string `json:"devicename"`
	// Whether the IP address of the bridge is obtained with DHCP.
	Dhcp bool `json:"dhcp"`
	// The current Zigbee firmware version.
	FWVersion string `json:"fwversion"`
	// IPv4 address of the gateway.
	Gateway string `json:"gateway"`
	// IPv4 address of the gateway.
	Ipaddress string `json:"ipaddress"`
	// true if the gateway is unlocked.
	Linkbutton bool `json:"linkbutton"`
	// The localtime of the gateway
	Localtime string `json:"localtime"`
	// MAC address of the gateway.
	Mac string `json:"mac"`
	// Fixed string "deCONZ".
	ModelId string `json:"modelid"`
	// Name of the gateway.
	Name string `json:"name"`
	// Network mask of the gateway.
	Netmask string `json:"netmask"`
	// (1–65535) The duration in seconds used by lights and sensors search, see Modify configuration.
	NetworkOpenDuration uint16 `json:"networkopenduration"`
	// Only for gateways running on Linux. Tells if the NTP time is "synced" or "unsynced".
	Ntp string `json:"ntp"`
	// (0–65535) The Zigbee pan ID of the gateway.
	PanId uint16 `json:"panid"`
	// This indicates whether the bridge is registered to synchronize data with a portal account.
	PortalServices bool `json:"portalservices"`
	// Not supported
	ProxyAddress string `json:"proxyaddress"`
	// Not supported
	ProxyPort int `json:"proxyport"`
	// Is true when the deCONZ is connected with the firmware and the Zigbee network is up.
	RFConnected bool `json:"rfconnected"`
	// Contains information related to software updates.
	SWUpdate struct {
		Notify      bool   `json:"notify"`
		Text        string `json:"text"`
		UpdateState int    `json:"updatestate"`
		Url         string `json:"url"`
	} `json:"swupdate"`
	// The software version of the gateway.
	SWVersion string `json:"swversion"`
	// Stores a value of the time format that can be used by other applications. "12h" or "24h"
	TimeFormat string `json:"timeformat"`
	// Timezone used by the gateway (only on Raspberry Pi). "None" if not further specified.
	Timezone string `json:"timezone"`
	// Current UTC time of the gateway in ISO 8601 format.
	UTC string `json:"UTC"`
	// UPNP Unique ID of the gateway
	Uuid string `json:"uuid"`
	// When true all state changes will be signalled through the Websocket connection (default true).
	WebsocketNotifyAll bool `json:"websocketnotifyall"`
	// Port of the Websocket server.
	WebsocketPort int `json:"websocketport"`
	// An array of whitelisted API keys.
	Whitelist map[string]WhitelistedAPIKey `json:"whitelist"`
	// The current wireless frequency channel used by the Gateway. Supported channels: 11, 15, 20, 25.
	ZigbeeChannel int `json:"zigbeechannel"`
}

// GetConfig Returns the current gateway configuration.
func (c *Client[R]) GetConfig(config *ConfigResponse) (R, error) {
	return c.Get("/config", config)
}

type ConfigResponseFullState struct {
	// Configuration of the gateway.
	Config ConfigResponse `json:"config"`
	// All groups of the gateway.
	Groups map[string]GroupResponseAttribute `json:"groups"`
	// All lights of the gateway.
	Lights map[string]LightResponseState `json:"lights"`
	// All rules of the gateway. (as from deconz version > 2.04.12) TODO needs to be implemented
	Rules json.RawMessage `json:"rules"`
	// All schedules of the gateway. TODO needs to be implemented
	Schedules json.RawMessage `json:"schedules"`
}

// GetFullState Returns the full state of the gateway including all its lights, groups, scenes and schedules.
func (c *Client[R]) GetFullState(fullState *ConfigResponseFullState) (R, error) {
	return c.Get("", fullState)
}

type ConfigTimeFormat string

const (
	ConfigTimeFormat12h ConfigTimeFormat = "12h"
	ConfigTimeFormat24h ConfigTimeFormat = "24h"
)

type UpdateChannel string

const (
	UpdateChannelStable UpdateChannel = "stable"
	UpdateChannelAlpha  UpdateChannel = "alpha"
	UpdateChannelBeta   UpdateChannel = "beta"
)

type ZigbeeChannel uint8

const (
	Channel11 ZigbeeChannel = 11
	Channel15 ZigbeeChannel = 15
	Channel20 ZigbeeChannel = 20
	Channel25 ZigbeeChannel = 25
)

type ConfigRequest struct {
	// 	Set gateway discovery over the internet active or inactive.	optional
	Discovery *bool `json:"discovery,omitempty"`
	//  (0–5000)	Time between two group commands in milliseconds.	optional
	GroupDelay *uint16 `json:"groupdelay,omitempty"`
	//  (1–65535) Default: 60; Sets the number of seconds where the timestamp for "lastseen" is updated at the earliest for light resources. For any such update, a seperate websocket event will be triggered.
	LightLastSeenInterval *uint16 `json:"lightlastseeninterval,omitempty"`
	// (0–16 chars)	Name of the gateway.	optional
	Name *string `json:"name,omitempty"`
	// (1–65535) Sets the lights and sensors search duration in seconds.	optional
	NetworkOpenDuration *uint16 `json:"networkopenduration,omitempty"`
	// Set OTAU active or inactive.	optional
	OTAUActive *bool `json:"otauactive,omitempty"`
	//(0–255) Open the network so that other zigbee devices can join.
	//0 = network closed,
	//255 = network open,
	//1–254 = time in seconds the network remains open. The value will decrement automatically.	optional
	PermitJoin *uint8 `json:"permitjoin,omitempty"`
	// Set to true to bring the Zigbee network up and false to bring it down.
	// This has the same effect as using the Join and Leave buttons in deCONZ.	optional
	RFConnected *bool `json:"rfconnected,omitempty"`
	// 	Can be used to store the timeformat permanently. It can be either "12h" or "24h".	optional
	TimeFormat *ConfigTimeFormat `json:"timeformat,omitempty"`
	// Set the timezone of the gateway (only on Raspberry Pi).
	// Format: tzdatabase e.g. “Europe/Berlin” Wikipedia:ListOfTimeZones	optional
	Timezone *string `json:"timezone,omitempty"`
	// (0–600)	Unlock the gateway so that apps can register themselves to the gateway (time in seconds).	optional
	Unlock *uint16 `json:"unlock,omitempty"`
	// Set update channel ("stable"|"alpha"|"beta").	optional
	UpdateChannel *UpdateChannel `json:"updatechannel,omitempty"`
	// 	Set the UTC time of the gateway (only on Raspberry Pi) in ISO 8601 format (yyyy-MM-ddTHH:mm:ss). optional
	Utc string `json:"utc,omitempty"`
	// Set the zigbee channel of the gateway. Notify other Zigbee devices also to change their channel.	optional
	ZigbeeChannel *ZigbeeChannel `json:"zigbeechannel,omitempty"`
	// When true all state changes will be signalled through the Websocket connection (default true).	optional
	WebsocketNotifyAll *bool `json:"websocketnotifyall,omitempty"`
}

// ModifyConfiguration Modify configuration parameters.
func (c *Client[R]) ModifyConfiguration(request ConfigRequest) (R, error) {
	return c.Put("/config", request)
}

// UpdateSoftware Returns the newest software version available. Starts the update if available (only on Raspberry Pi).
func (c *Client[R]) UpdateSoftware() (R, error) {
	return c.Post("/config/update", nil)
}

// UpdateFirmware Starts the update firmware process if newer version is available.
func (c *Client[R]) UpdateFirmware() (R, error) {
	return c.Post("/config/updatefirmware", nil)
}

// ConfigRequestReset At least one parameter is required!
type ConfigRequestReset struct {
	// 	Set the network settings of the gateway to factory new.	optional
	ResetGW *bool `json:"resetGW,omitempty"`
	// 	Delete the Database.	optional
	DeleteDB *bool `json:"deleteDB,omitempty"`
}

// ResetGateway Reset the gateway network settings to factory new and/or delete the deCONZ database
// (config, lights, scenes, groups, schedules, devices, rules).
func (c *Client[R]) ResetGateway(request ConfigRequestReset) (R, error) {
	return c.Post("/config/reset", request)
}

type ConfigRequestChangePassword struct {
	// 	The username (currently only “delight” is supported).	required
	Username string `json:"username"`
	// 	The Base64 encoded combination of “username:old password”.	required
	OldHash string `json:"oldhash"`
	// 	The Base64 encoded combination of “username:new password”.	required
	NewHash string `json:"newhash"`
}

// ChangePasswordX Change the Password of the Gateway.
func (c *Client[R]) ChangePasswordX(username string, oldPassword string, newPassword string) (R, error) {
	builder := strings.Builder{}
	encoder := base64.NewEncoder(base64.StdEncoding, &builder)
	_, err := encoder.Write([]byte(username + ":" + oldPassword))
	if err != nil {
		panic(err)
	}
	oldHash := builder.String()
	builder.Reset()
	_, err = encoder.Write([]byte(username + ":" + newPassword))
	if err != nil {
		panic(err)
	}
	newHash := builder.String()
	request := ConfigRequestChangePassword{
		Username: username,
		OldHash:  oldHash,
		NewHash:  newHash,
	}
	return c.ChangePassword(request)
}

// ChangePassword Change the Password of the Gateway. The parameter must be a Base64 encoded string of
// <username>:<password>.
func (c *Client[R]) ChangePassword(request ConfigRequestChangePassword) (R, error) {
	return c.Put("/config/password", request)
}

// ResetPassword Resets the username and password to default username = “delight” and password = “delight”.
// The request can only succeed within 10 minutes after gateway start.
func (c *Client[R]) ResetPassword() (R, error) {
	return c.Delete("/config/password", nil)
}
