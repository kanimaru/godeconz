# Changelog

## v0.1.0 (2026-10-12)

Rewrite. The API is not compatible with v0.0.x.

### Fixed

- POST requests were sent as PUT, so creating groups, sensors and API keys never worked.
- HTTP error statuses and deCONZ error objects were only logged; they are now returned as
  `*rest.Error` (`errors.Is` with `ErrNotFound`, `ErrUnauthorized`, `ErrLinkButton`). Partly
  failed PUTs report the rejected parts.
- Transport errors included the request URL and therefore the API key.
- Websocket: crashed on a failed dial, busy-looped after a disconnect, never reconnected,
  hijacked SIGINT, and the state callback exited the process on the first message.
- Websocket resource filters never matched: deCONZ sends lowercase resource names
  (`"sensors"`), the constants were capitalised. Scene events used the wrong field names
  (`gid`/`scid`).
- Request bodies serialised wrongly: unexported fields (`{}` for light rename/delete),
  group actions sent `null`/`""` for unset fields, read structs reused as update bodies sent
  `"on": false`.
- Thermostat setpoints were `uint8` and failed to decode real values.
- Lights read `manufacturer` instead of `manufacturername`.
- `go vet` failed (logger and callback format strings), so `go test ./...` failed.

### Changed

- Packages: `http` → `rest` (no longer shadows `net/http`); one shared model in the root package
  instead of separate structs in `http` and `ws`.
- Every REST call takes a `context.Context` and returns typed values instead of filling a
  caller-provided container.
- Sensor state and config are one struct each (`SensorState`, `SensorConfig`) covering all
  common types, with unit helpers. Unknown attributes are kept in `Raw`; an attribute with an
  unexpected JSON type is left nil instead of failing the whole response.
- Resource maps fill in `ID` from the map key.
- Websocket: `ws.New(url).Run(ctx, handler)` with reconnect, backoff, ping and an on-connect
  hook; `ws.URLFromREST` uses the REST host instead of the gateway's self-reported IP.
- Logging uses `log/slog` and is off by default.
- Requires Go 1.25. Dependencies: only `github.com/coder/websocket`.

### Added

- Scenes: list, get, create, store, recall, delete.
- Group state (`all_on`, `any_on`), sensor type constants, `ZHABattery`, air quality and PM2.5
  attributes, `ParseTime`, `DeviceMAC`, `Ptr`.
- CI workflow (gofmt, vet, race-enabled tests) and an opt-in integration test against a real
  gateway.

### Removed

- Resty and go-stockutil dependencies, the `Logger` interface and `Settings`.
- ETag caching: it was not safe for concurrent use and could return another caller's data.
  Use websocket events for change notification.
- Firmware/software update, gateway reset and password endpoints. They never worked (POST bug)
  and are risky to ship untested; they can come back with tests.
- Per-type state structs (`StateZHATemperature`, …) and CLIP structs; use `SensorState`.

### Migration

| v0.0.x | v0.1.0 |
|---|---|
| `http.CreateClient(http.CreateAdapterHttpClientResty(resty.New(), logger, false), settings)` | `rest.New(address, apiKey)` |
| `client.GetAllLights(&m)` | `m, err := client.Lights(ctx)` |
| `client.GetLightState(id, &l)` | `l, err := client.Light(ctx, id)` |
| `client.SetLightState(id, http.LightRequestState{...})` | `client.SetLightState(ctx, id, godeconz.LightStateUpdate{...})` |
| `client.SetGroupState(id, ...)` | `client.SetGroupAction(ctx, id, godeconz.GroupActionUpdate{...})` |
| `client.GetAllSensors(&m)` + `s.StateAs(&http.StateZHATemperature{})` | `m, err := client.Sensors(ctx)`; `m[id].State.TemperatureCelsius()` |
| `client.GetConfig(&cfg)` | `cfg, err := client.Config(ctx)` |
| `ws.CreateClientFromConfig(client, ws.CreateAdapterWebsocketClientGorilla(...), logger)` + `Connect()` | `url, err := ws.URLFromREST(ctx, client)`; `ws.New(url).Run(ctx, handler)` |
| `ws.NewStateCallback` / `ws.NewChanCallback` + `ws.Filter` | `ws.Filter{...}.Handler(handler)`; `event.SensorState()` |
| `godeconz.StateLight` / `godeconz.Config` interfaces | `godeconz.LightState` / `godeconz.SensorConfig` helpers (`IsOn`, `Brightness`, `IsReachable`) |
