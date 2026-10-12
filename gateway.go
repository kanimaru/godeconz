package godeconz

// GatewayConfig is the gateway configuration from GET /config.
type GatewayConfig struct {
	Name                string                    `json:"name"`
	APIVersion          string                    `json:"apiversion"`
	SWVersion           string                    `json:"swversion"`
	FWVersion           string                    `json:"fwversion"`
	BridgeID            string                    `json:"bridgeid"`
	ModelID             string                    `json:"modelid"`
	UUID                string                    `json:"uuid"`
	Mac                 string                    `json:"mac"`
	IPAddress           string                    `json:"ipaddress"`
	Netmask             string                    `json:"netmask"`
	Gateway             string                    `json:"gateway"`
	DHCP                bool                      `json:"dhcp"`
	Timezone            string                    `json:"timezone"`
	UTC                 string                    `json:"UTC"`
	LocalTime           string                    `json:"localtime"`
	ZigbeeChannel       int                       `json:"zigbeechannel"`
	PanID               int                       `json:"panid"`
	RFConnected         bool                      `json:"rfconnected"`
	PermitJoin          int                       `json:"permitjoin"`
	NetworkOpenDuration int                       `json:"networkopenduration"`
	LinkButton          bool                      `json:"linkbutton"`
	WebsocketPort       int                       `json:"websocketport"`
	WebsocketNotifyAll  bool                      `json:"websocketnotifyall"`
	Whitelist           map[string]WhitelistEntry `json:"whitelist"`
}

// WhitelistEntry is an API key known to the gateway.
type WhitelistEntry struct {
	Name        string `json:"name"`
	CreateDate  string `json:"create date"`
	LastUseDate string `json:"last use date"`
}

// GatewayConfigUpdate is the body of PUT /config. Only set fields are sent.
type GatewayConfigUpdate struct {
	Name                  *string `json:"name,omitempty"`
	PermitJoin            *uint8  `json:"permitjoin,omitempty"` // seconds, 0 closes
	NetworkOpenDuration   *uint16 `json:"networkopenduration,omitempty"`
	WebsocketNotifyAll    *bool   `json:"websocketnotifyall,omitempty"`
	Timezone              *string `json:"timezone,omitempty"`
	LightLastSeenInterval *uint16 `json:"lightlastseeninterval,omitempty"`
	GroupDelay            *uint16 `json:"groupdelay,omitempty"`
	OTAUActive            *bool   `json:"otauactive,omitempty"`
	Discovery             *bool   `json:"discovery,omitempty"`
	RFConnected           *bool   `json:"rfconnected,omitempty"`
	Unlock                *uint16 `json:"unlock,omitempty"` // seconds the gateway accepts new API keys
}

// FullState is the response of GET /, the whole gateway state in one call.
type FullState struct {
	Config  GatewayConfig     `json:"config"`
	Lights  map[string]Light  `json:"lights"`
	Sensors map[string]Sensor `json:"sensors"`
	Groups  map[string]Group  `json:"groups"`
}

// FillIDs copies the map keys into the ID fields of lights, sensors and groups.
func (f *FullState) FillIDs() {
	for id, l := range f.Lights {
		l.ID = id
		f.Lights[id] = l
	}
	for id, s := range f.Sensors {
		s.ID = id
		f.Sensors[id] = s
	}
	for id, g := range f.Groups {
		g.ID = id
		f.Groups[id] = g
	}
}
