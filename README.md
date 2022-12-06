# godeconz

... is a binding for the deconz Rest and Websocket API
It has exchangeable HTTP Client (default Resty) and Logger (default standard go logger)

## Features
- [ ] Alarm Systems 
- [ ] Configuration 
- [ ] Groups 
- [x] Lights 
- [ ] Rules 
- [ ] Scenes 
- [ ] Schedules 
- [x] Sensors
- [ ] Button Events
- [ ] Touchlink
- [ ] Websocket

## Example

Query all lights:
```go
package main

import (
	"github.com/PerformLine/go-stockutil/log"
	"github.com/go-resty/resty/v2"
	"github.com/kanimaru/godeconz"
	"os"
)

func getEnv(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	if len(fallback) == 0 {
		log.Fatalf("Missing %q", key)
	}
	return fallback
}

func main() {
	setting := godeconz.Settings{
		Address:      getEnv("DECONZ_ADDRESS", ""),
		HttpProtocol: getEnv("DECONZ_PROTO", "http"),
		ApiKey:       getEnv("DECONZ_API_KEY", ""),
	}

	httpAdapter := godeconz.CreateAdapterHttpClientResty(resty.New(), Logger{}, false)
	deconzClient := godeconz.CreateClient(httpAdapter, setting)

	var container map[string]godeconz.LightResponseState
	_, err := deconzClient.GetAllLights(&container)
	if err != nil {
		log.Fatalf("Can't query lights", err)
	}

	for _, light := range container {
		log.Noticef("Name: %v", light.Name)
	}
}

```

## Resources

https://dresden-elektronik.github.io/deconz-rest-doc/