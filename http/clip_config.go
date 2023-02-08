package http

type CLIPConfigHumidity struct {
	Offset int16 `json:"offset"`
}

type CLIPConfigLightLevel struct {
	Tholddark       uint16 `json:"tholddark"`
	Tholddarkoffset uint16 `json:"tholddarkoffset"`
}

type CLIPConfigTemperature struct {
	Offset int16 `json:"offset"`
}
