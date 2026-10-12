// Package rest is a client for the deCONZ REST API.
//
//	client, err := rest.New("deconz.kani.home", apiKey)
//	lights, err := client.Lights(ctx)
//	err = client.SetLightState(ctx, "4", godeconz.LightStateUpdate{On: godeconz.Ptr(true)})
//
// Every call takes a context. Non-2xx responses and deCONZ error objects inside 2xx responses
// are returned as *Error; use errors.Is with ErrNotFound or ErrUnauthorized to classify them.
// Errors never contain the API key.
package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/kanimaru/godeconz"
)

const (
	defaultTimeout   = 10 * time.Second
	maxResponseBytes = 32 << 20
)

// Client talks to one deCONZ gateway. It is safe for concurrent use.
type Client struct {
	base       *url.URL
	apiKey     string
	httpClient *http.Client
	logger     *slog.Logger
}

// Option configures a Client.
type Option func(*Client)

// WithHTTPClient replaces the default HTTP client (10 s timeout).
func WithHTTPClient(c *http.Client) Option {
	return func(client *Client) {
		client.httpClient = c
	}
}

// WithLogger sets a logger for diagnostics such as attributes that could not be decoded.
// The default discards all output.
func WithLogger(l *slog.Logger) Option {
	return func(client *Client) {
		client.logger = l
	}
}

// New creates a client. address is "host", "host:port" or "http(s)://host[:port]".
func New(address, apiKey string, opts ...Option) (*Client, error) {
	base, err := godeconz.ParseAddress(address)
	if err != nil {
		return nil, err
	}
	if apiKey == "" {
		return nil, errors.New("deconz: empty API key")
	}
	c := &Client{
		base:       base,
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: defaultTimeout},
		logger:     slog.New(slog.DiscardHandler),
	}
	for _, opt := range opts {
		opt(c)
	}
	return c, nil
}

// BaseURL returns the gateway address (scheme and host) the client talks to.
func (c *Client) BaseURL() *url.URL {
	u := *c.base
	return &u
}

// do calls an endpoint below /api/<key>. path starts with "/".
func (c *Client) do(ctx context.Context, method, path string, body, out any) error {
	display := path
	if display == "" {
		display = "/"
	}
	return c.request(ctx, method, "/api/"+url.PathEscape(c.apiKey)+path, display, body, out)
}

// request performs the call. urlPath goes on the wire; displayPath is used in errors and logs
// so the API key never leaves this function.
func (c *Client) request(ctx context.Context, method, urlPath, displayPath string, body, out any) error {
	var reader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("deconz: %s %s: encode body: %w", method, displayPath, err)
		}
		reader = bytes.NewReader(data)
	}
	req, err := http.NewRequestWithContext(ctx, method, c.base.String()+urlPath, reader)
	if err != nil {
		return fmt.Errorf("deconz: %s %s: %w", method, displayPath, err)
	}
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		// *url.Error repeats the full URL, which contains the API key.
		var urlErr *url.Error
		if errors.As(err, &urlErr) {
			err = urlErr.Err
		}
		return fmt.Errorf("deconz: %s %s: %w", method, displayPath, err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseBytes))
	if err != nil {
		return fmt.Errorf("deconz: %s %s: read response: %w", method, displayPath, err)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return newError(method, displayPath, resp.StatusCode, data)
	}
	if apiErrors := findAPIErrors(data); len(apiErrors) > 0 {
		return &Error{Method: method, Path: displayPath, StatusCode: resp.StatusCode, Errors: apiErrors}
	}
	if out == nil {
		return nil
	}
	if err := json.Unmarshal(data, out); err != nil {
		var typeErr *json.UnmarshalTypeError
		if !errors.As(err, &typeErr) {
			return fmt.Errorf("deconz: %s %s: decode response: %w", method, displayPath, err)
		}
		// One odd attribute must not hide all other devices; the rest is decoded.
		c.logger.WarnContext(ctx, "deconz: attribute with unexpected type ignored",
			"method", method, "path", displayPath, "field", typeErr.Field, "value", typeErr.Value)
	}
	return nil
}

// create calls a POST endpoint that answers [{"success":{"id":"..."}}] and returns the id.
func (c *Client) create(ctx context.Context, path string, body any) (string, error) {
	var result []struct {
		Success struct {
			ID string `json:"id"`
		} `json:"success"`
	}
	if err := c.do(ctx, http.MethodPost, path, body, &result); err != nil {
		return "", err
	}
	if len(result) == 0 || result[0].Success.ID == "" {
		return "", fmt.Errorf("deconz: POST %s: response contains no id", path)
	}
	return result[0].Success.ID, nil
}

// getMap fetches a map of resources and copies each key into the resource via setID.
func getMap[T any](ctx context.Context, c *Client, path string, setID func(*T, string)) (map[string]T, error) {
	var resources map[string]T
	if err := c.do(ctx, http.MethodGet, path, nil, &resources); err != nil {
		return nil, err
	}
	if resources == nil {
		resources = map[string]T{}
	}
	for id, r := range resources {
		setID(&r, id)
		resources[id] = r
	}
	return resources, nil
}
