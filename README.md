# godeconz

Go client for the [deCONZ](https://dresden-elektronik.github.io/deconz-rest-doc/) REST and
websocket API (Phoscon / ConBee / RaspBee Zigbee gateways).

- `godeconz` – shared data model: lights, sensors, groups, scenes, gateway config, with unit
  helpers (`TemperatureCelsius`, `PowerWatts`, …) and timestamp parsing
- `rest` – REST client on `net/http`; every call takes a `context`, errors are typed
- `ws` – websocket client that reconnects, pings and lets you resync after every reconnect

Dependencies: the standard library and [`github.com/coder/websocket`](https://github.com/coder/websocket).
Requires Go 1.25.

## Features

- [x] Configuration (read, update, API key acquire/delete, full state)
- [x] Lights
- [x] Sensors (all common ZHA types in one `SensorState`, unknown attributes kept in `Raw`)
- [x] Groups
- [x] Scenes (list, create, store, recall, delete)
- [x] Websocket events with reconnect
- [ ] Rules, schedules, resource links
- [ ] Alarm systems, touchlink, `/devices`

## REST

```go
client, err := rest.New("deconz.kani.home", os.Getenv("DECONZ_API_KEY"))
if err != nil {
	log.Fatal(err)
}

sensors, err := client.Sensors(ctx)
if err != nil {
	log.Fatal(err)
}
for _, s := range sensors {
	if t, ok := s.State.TemperatureCelsius(); ok {
		fmt.Printf("%s: %.1f °C\n", s.Name, t)
	}
}

err = client.SetLightState(ctx, "4", godeconz.LightStateUpdate{
	On:             godeconz.Ptr(true),
	Bri:            godeconz.Ptr[uint8](180),
	TransitionTime: godeconz.Ptr[uint16](10), // 1 s
})
if errors.Is(err, rest.ErrNotFound) {
	// light 4 does not exist
}
```

The address can be `host`, `host:port` or a full `http(s)://host:port` URL. Update types only
send the fields you set. deCONZ can accept part of a PUT and reject the rest; the error then
lists the rejected parts (`*rest.Error`). Error messages never contain the API key.

New API key (unlock the gateway in Phoscon first):

```go
key, err := rest.AcquireAPIKey(ctx, "deconz.kani.home", "myapp#host")
```

## Websocket

```go
url, err := ws.URLFromREST(ctx, client) // REST host + websocket port from /config
if err != nil {
	log.Fatal(err)
}

listener := ws.New(url,
	ws.WithOnConnect(func(ctx context.Context) {
		// Events sent while disconnected are lost: re-read state here.
	}),
)

sensorEvents := ws.Filter{Resources: []ws.Resource{ws.ResourceSensors}, WithState: true}
err = listener.Run(ctx, sensorEvents.Handler(func(e ws.Event) {
	state, err := e.SensorState() // only changed attributes are set
	if err != nil {
		return
	}
	if open := state.Open; open != nil {
		fmt.Printf("sensor %s open=%v\n", e.ID, *open)
	}
}))
// Run returns ctx.Err() once ctx ends; until then it reconnects with backoff (1 s … 1 min).
```

## Identity

REST ids (`/sensors/12`) change when a device is re-paired. Use `UniqueID` for stable identity
and `godeconz.DeviceMAC(uniqueID)` to group the resources of one physical device (e.g. the
temperature, humidity and pressure sensors of one weather sensor).

## Testing

```bash
go test ./...
```

Read-only checks against a real gateway (skipped without the variables):

```bash
DECONZ_ADDRESS=deconz.local DECONZ_API_KEY=... DECONZ_WS_SECONDS=60 go test -run Gateway -v .
```

Upgrading from v0.0.x: see [CHANGELOG.md](CHANGELOG.md).
