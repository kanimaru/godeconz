package rest

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/kanimaru/godeconz"
)

func p(format string, ids ...string) string {
	args := make([]any, len(ids))
	for i, id := range ids {
		args[i] = url.PathEscape(id)
	}
	return fmt.Sprintf(format, args...)
}

// Gateway

// Config returns the gateway configuration (GET /config).
func (c *Client) Config(ctx context.Context) (godeconz.GatewayConfig, error) {
	var cfg godeconz.GatewayConfig
	err := c.do(ctx, http.MethodGet, "/config", nil, &cfg)
	return cfg, err
}

// UpdateConfig changes gateway settings (PUT /config), e.g. PermitJoin to pair devices.
func (c *Client) UpdateConfig(ctx context.Context, update godeconz.GatewayConfigUpdate) error {
	return c.do(ctx, http.MethodPut, "/config", update, nil)
}

// FullState returns config, lights, sensors and groups in one call (GET /), with IDs filled in.
func (c *Client) FullState(ctx context.Context) (godeconz.FullState, error) {
	var state godeconz.FullState
	if err := c.do(ctx, http.MethodGet, "", nil, &state); err != nil {
		return state, err
	}
	state.FillIDs()
	return state, nil
}

// DeleteAPIKey removes an API key from the gateway whitelist.
func (c *Client) DeleteAPIKey(ctx context.Context, key string) error {
	return c.do(ctx, http.MethodDelete, p("/config/whitelist/%s", key), nil, nil)
}

// AcquireAPIKey requests a new API key (POST /api). The gateway must be unlocked first, by the
// link button in the Phoscon app or by UpdateConfig with Unlock; otherwise the error matches
// ErrLinkButton. deviceType names the application, e.g. "hestia#collector".
func AcquireAPIKey(ctx context.Context, address, deviceType string, opts ...Option) (string, error) {
	c, err := New(address, "unused", opts...)
	if err != nil {
		return "", err
	}
	var result []struct {
		Success struct {
			Username string `json:"username"`
		} `json:"success"`
	}
	body := map[string]string{"devicetype": deviceType}
	if err := c.request(ctx, http.MethodPost, "/api", "/api", body, &result); err != nil {
		return "", err
	}
	if len(result) == 0 || result[0].Success.Username == "" {
		return "", fmt.Errorf("deconz: POST /api: response contains no API key")
	}
	return result[0].Success.Username, nil
}

// Lights

// Lights returns all lights keyed by id (GET /lights).
func (c *Client) Lights(ctx context.Context) (map[string]godeconz.Light, error) {
	return getMap(ctx, c, "/lights", func(l *godeconz.Light, id string) { l.ID = id })
}

// Light returns one light (GET /lights/{id}).
func (c *Client) Light(ctx context.Context, id string) (godeconz.Light, error) {
	var light godeconz.Light
	err := c.do(ctx, http.MethodGet, p("/lights/%s", id), nil, &light)
	light.ID = id
	return light, err
}

// SetLightState changes the state of a light (PUT /lights/{id}/state).
func (c *Client) SetLightState(ctx context.Context, id string, update godeconz.LightStateUpdate) error {
	return c.do(ctx, http.MethodPut, p("/lights/%s/state", id), update, nil)
}

// SetLightAttributes changes attributes such as the name (PUT /lights/{id}).
func (c *Client) SetLightAttributes(ctx context.Context, id string, update godeconz.LightAttributesUpdate) error {
	return c.do(ctx, http.MethodPut, p("/lights/%s", id), update, nil)
}

// DeleteLight removes a light from the gateway. reset also factory-resets the device.
func (c *Client) DeleteLight(ctx context.Context, id string, reset bool) error {
	return c.do(ctx, http.MethodDelete, p("/lights/%s", id), map[string]bool{"reset": reset}, nil)
}

// RemoveLightFromGroups removes a light from all groups (DELETE /lights/{id}/groups).
func (c *Client) RemoveLightFromGroups(ctx context.Context, id string) error {
	return c.do(ctx, http.MethodDelete, p("/lights/%s/groups", id), nil, nil)
}

// RemoveLightFromScenes removes a light from all scenes (DELETE /lights/{id}/scenes).
func (c *Client) RemoveLightFromScenes(ctx context.Context, id string) error {
	return c.do(ctx, http.MethodDelete, p("/lights/%s/scenes", id), nil, nil)
}

// Sensors

// Sensors returns all sensors keyed by id (GET /sensors).
func (c *Client) Sensors(ctx context.Context) (map[string]godeconz.Sensor, error) {
	return getMap(ctx, c, "/sensors", func(s *godeconz.Sensor, id string) { s.ID = id })
}

// Sensor returns one sensor (GET /sensors/{id}).
func (c *Client) Sensor(ctx context.Context, id string) (godeconz.Sensor, error) {
	var sensor godeconz.Sensor
	err := c.do(ctx, http.MethodGet, p("/sensors/%s", id), nil, &sensor)
	sensor.ID = id
	return sensor, err
}

// CreateSensor creates a CLIP sensor and returns its id (POST /sensors).
func (c *Client) CreateSensor(ctx context.Context, sensor godeconz.SensorCreate) (string, error) {
	return c.create(ctx, "/sensors", sensor)
}

// SetSensorAttributes changes name or mode (PUT /sensors/{id}).
func (c *Client) SetSensorAttributes(ctx context.Context, id string, update godeconz.SensorAttributesUpdate) error {
	return c.do(ctx, http.MethodPut, p("/sensors/%s", id), update, nil)
}

// SetSensorConfig changes sensor configuration (PUT /sensors/{id}/config). update is a
// godeconz.SensorConfig with only the fields to change, or a map for attributes without a field.
func (c *Client) SetSensorConfig(ctx context.Context, id string, update any) error {
	return c.do(ctx, http.MethodPut, p("/sensors/%s/config", id), update, nil)
}

// SetSensorState changes the state of a CLIP sensor (PUT /sensors/{id}/state).
func (c *Client) SetSensorState(ctx context.Context, id string, update any) error {
	return c.do(ctx, http.MethodPut, p("/sensors/%s/state", id), update, nil)
}

// DeleteSensor removes a sensor (DELETE /sensors/{id}).
func (c *Client) DeleteSensor(ctx context.Context, id string) error {
	return c.do(ctx, http.MethodDelete, p("/sensors/%s", id), nil, nil)
}

// Groups

// Groups returns all groups keyed by id (GET /groups).
func (c *Client) Groups(ctx context.Context) (map[string]godeconz.Group, error) {
	return getMap(ctx, c, "/groups", func(g *godeconz.Group, id string) { g.ID = id })
}

// Group returns one group with its lights and scenes (GET /groups/{id}).
func (c *Client) Group(ctx context.Context, id string) (godeconz.Group, error) {
	var group godeconz.Group
	err := c.do(ctx, http.MethodGet, p("/groups/%s", id), nil, &group)
	group.ID = id
	return group, err
}

// CreateGroup creates a group and returns its id. If the name exists, deCONZ returns the id of
// the existing group.
func (c *Client) CreateGroup(ctx context.Context, name string) (string, error) {
	return c.create(ctx, "/groups", map[string]string{"name": name})
}

// SetGroupAttributes changes name, members or visibility (PUT /groups/{id}).
func (c *Client) SetGroupAttributes(ctx context.Context, id string, update godeconz.GroupAttributesUpdate) error {
	return c.do(ctx, http.MethodPut, p("/groups/%s", id), update, nil)
}

// SetGroupAction changes the state of all lights in a group (PUT /groups/{id}/action).
func (c *Client) SetGroupAction(ctx context.Context, id string, update godeconz.GroupActionUpdate) error {
	return c.do(ctx, http.MethodPut, p("/groups/%s/action", id), update, nil)
}

// DeleteGroup deletes a group (DELETE /groups/{id}).
func (c *Client) DeleteGroup(ctx context.Context, id string) error {
	return c.do(ctx, http.MethodDelete, p("/groups/%s", id), nil, nil)
}

// Scenes

// Scenes returns the scenes of a group keyed by scene id (GET /groups/{id}/scenes).
func (c *Client) Scenes(ctx context.Context, groupID string) (map[string]godeconz.SceneRef, error) {
	return getMap(ctx, c, p("/groups/%s/scenes", groupID), func(s *godeconz.SceneRef, id string) { s.ID = id })
}

// Scene returns a scene with its light states (GET /groups/{gid}/scenes/{sid}).
func (c *Client) Scene(ctx context.Context, groupID, sceneID string) (godeconz.Scene, error) {
	var scene godeconz.Scene
	err := c.do(ctx, http.MethodGet, p("/groups/%s/scenes/%s", groupID, sceneID), nil, &scene)
	scene.ID = sceneID
	return scene, err
}

// CreateScene creates a scene from the current state of the group's lights and returns its id.
func (c *Client) CreateScene(ctx context.Context, groupID, name string) (string, error) {
	return c.create(ctx, p("/groups/%s/scenes", groupID), map[string]string{"name": name})
}

// StoreScene stores the current state of the group's lights in the scene.
func (c *Client) StoreScene(ctx context.Context, groupID, sceneID string) error {
	return c.do(ctx, http.MethodPut, p("/groups/%s/scenes/%s/store", groupID, sceneID), nil, nil)
}

// RecallScene applies a scene to its lights.
func (c *Client) RecallScene(ctx context.Context, groupID, sceneID string) error {
	return c.do(ctx, http.MethodPut, p("/groups/%s/scenes/%s/recall", groupID, sceneID), nil, nil)
}

// DeleteScene deletes a scene.
func (c *Client) DeleteScene(ctx context.Context, groupID, sceneID string) error {
	return c.do(ctx, http.MethodDelete, p("/groups/%s/scenes/%s", groupID, sceneID), nil, nil)
}
