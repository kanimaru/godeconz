package ws

// There is no documentation about the websocket configs so this will be user gathered configs here.
// Add your own as soon as you discover them.

type BaseConfig struct {
	// (default: true)
	On bool `json:"on"`
	// (default: true)
	Reachable *bool `json:"reachable,omitempty"`
	// (0–100)
	Battery *uint8 `json:"battery,omitempty"`
}

func (b BaseConfig) IsOn() bool {
	return b.On
}

func (b BaseConfig) IsReachable() *bool {
	return b.Reachable
}

func (b BaseConfig) GetBattery() *uint8 {
	return b.Battery
}

type ConfigDaylight struct {
	BaseConfig
	Configured    bool `json:"configured"`
	On            bool `json:"on"`
	Sunriseoffset int  `json:"sunriseoffset"`
	Sunsetoffset  int  `json:"sunsetoffset"`
}

type ConfigZHAPressure struct {
	BaseConfig
	Offset int `json:"offset"`
}

type ConfigZHAHumidity struct {
	BaseConfig
	Offset int `json:"offset"`
}

type ConfigZHATemperature struct {
	BaseConfig
	Offset int `json:"offset"`
}

type ConfigZHAOpenClose struct {
	BaseConfig
	Temperature int `json:"temperature"`
}

type ConfigZHALightLevel struct {
	BaseConfig
	Alert          string `json:"alert"`
	Ledindication  bool   `json:"ledindication"`
	Tholddark      int    `json:"tholddark"`
	Tholdoffset    int    `json:"tholdoffset"`
	Usertest       bool   `json:"usertest"`
	Delay          int    `json:"delay"`
	Sensitivity    int    `json:"sensitivity"`
	Sensitivitymax int    `json:"sensitivitymax"`
}

type ConfigZHAPresence struct {
	BaseConfig

	Duration    int `json:"duration"`
	Temperature int `json:"temperature"`
}
