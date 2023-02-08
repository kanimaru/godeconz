package ws

// There is no documentation about the websocket states so this will be only user gathered states here.
// Add your own as soon as you discover them.

type ZHATemperatureConfig struct {
	Battery   int  `json:"battery"`
	Offset    int  `json:"offset"`
	On        bool `json:"on"`
	Reachable bool `json:"reachable"`
}

type ZHALightLevelConfig struct {
	Alert         string `json:"alert"`
	Battery       int    `json:"battery"`
	Ledindication bool   `json:"ledindication"`
	On            bool   `json:"on"`
	Reachable     bool   `json:"reachable"`
	Tholddark     int    `json:"tholddark"`
	Tholdoffset   int    `json:"tholdoffset"`
	Usertest      bool   `json:"usertest"`
}
