// Package ws receives deCONZ websocket events.
//
//	url, err := ws.URLFromREST(ctx, restClient)
//	client := ws.New(url, ws.WithOnConnect(resync))
//	err = client.Run(ctx, func(e ws.Event) { ... })
//
// Run keeps the connection alive: it pings the gateway, reconnects with exponential backoff and
// only returns when the context ends. Events sent while disconnected are lost, so use
// WithOnConnect to re-read state over REST after every (re)connect.
package ws

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/coder/websocket"

	"github.com/kanimaru/godeconz/rest"
)

const (
	defaultMinBackoff   = time.Second
	defaultMaxBackoff   = time.Minute
	defaultPingInterval = 30 * time.Second
	pingTimeout         = 10 * time.Second
	// "added" events contain whole resources; 1 MiB leaves plenty of room.
	readLimit = 1 << 20
)

// Client receives events from one deCONZ websocket.
type Client struct {
	url          string
	logger       *slog.Logger
	minBackoff   time.Duration
	maxBackoff   time.Duration
	pingInterval time.Duration
	onConnect    func(context.Context)
	dialOptions  *websocket.DialOptions
}

// Option configures a Client.
type Option func(*Client)

// WithLogger sets the logger for connection state changes. The default discards all output.
func WithLogger(l *slog.Logger) Option {
	return func(c *Client) {
		c.logger = l
	}
}

// WithBackoff sets the reconnect delays. The delay doubles from min up to max and resets to min
// once a connection succeeded. Defaults: 1 s and 1 min.
func WithBackoff(min, max time.Duration) Option {
	return func(c *Client) {
		c.minBackoff = min
		c.maxBackoff = max
	}
}

// WithPingInterval sets how often the connection is checked. A failed ping reconnects.
// Default 30 s; 0 disables pings.
func WithPingInterval(d time.Duration) Option {
	return func(c *Client) {
		c.pingInterval = d
	}
}

// WithOnConnect registers a function that runs after every successful (re)connect and before
// the first event of that connection is delivered, e.g. to re-read state over REST.
func WithOnConnect(f func(ctx context.Context)) Option {
	return func(c *Client) {
		c.onConnect = f
	}
}

// WithDialOptions passes options to the websocket dialer, e.g. a custom HTTP client.
func WithDialOptions(o *websocket.DialOptions) Option {
	return func(c *Client) {
		c.dialOptions = o
	}
}

// New creates a client for a websocket URL such as "ws://deconz.kani.home:8443".
func New(url string, opts ...Option) *Client {
	c := &Client{
		url:          url,
		logger:       slog.New(slog.DiscardHandler),
		minBackoff:   defaultMinBackoff,
		maxBackoff:   defaultMaxBackoff,
		pingInterval: defaultPingInterval,
	}
	for _, opt := range opts {
		opt(c)
	}
	if c.minBackoff <= 0 {
		c.minBackoff = defaultMinBackoff
	}
	if c.maxBackoff < c.minBackoff {
		c.maxBackoff = c.minBackoff
	}
	return c
}

// URLFromREST builds the websocket URL from the REST client's host and the websocket port in
// the gateway config. The host the REST client uses is kept on purpose: the gateway's own
// "ipaddress" is wrong when deCONZ runs in a container.
func URLFromREST(ctx context.Context, client *rest.Client) (string, error) {
	cfg, err := client.Config(ctx)
	if err != nil {
		return "", err
	}
	if cfg.WebsocketPort == 0 {
		return "", errors.New("deconz: gateway config has no websocket port")
	}
	base := client.BaseURL()
	scheme := "ws"
	if base.Scheme == "https" {
		scheme = "wss"
	}
	return fmt.Sprintf("%s://%s:%s", scheme, base.Hostname(), strconv.Itoa(cfg.WebsocketPort)), nil
}

// Run connects and calls handle for every event until ctx ends, reconnecting as needed. It
// returns ctx.Err(). handle runs on the reading goroutine; a slow handler delays reading.
func (c *Client) Run(ctx context.Context, handle func(Event)) error {
	backoff := c.minBackoff
	for {
		connected, err := c.session(ctx, handle)
		if ctx.Err() != nil {
			return ctx.Err()
		}
		if connected {
			backoff = c.minBackoff
		}
		c.logger.WarnContext(ctx, "deconz websocket disconnected", "url", c.url, "error", err, "retry_in", backoff)

		timer := time.NewTimer(backoff)
		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		case <-timer.C:
		}
		backoff = min(backoff*2, c.maxBackoff)
	}
}

// session runs one connection. connected reports whether the dial succeeded.
func (c *Client) session(ctx context.Context, handle func(Event)) (connected bool, err error) {
	conn, _, err := websocket.Dial(ctx, c.url, c.dialOptions)
	if err != nil {
		return false, err
	}
	defer conn.CloseNow()
	conn.SetReadLimit(readLimit)
	c.logger.InfoContext(ctx, "deconz websocket connected", "url", c.url)

	sessionCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	if c.onConnect != nil {
		c.onConnect(sessionCtx)
	}
	if c.pingInterval > 0 {
		go c.keepAlive(sessionCtx, conn, cancel)
	}

	for {
		msgType, data, err := conn.Read(sessionCtx)
		if err != nil {
			return true, err
		}
		if msgType != websocket.MessageText {
			continue
		}
		var event Event
		if err := json.Unmarshal(data, &event); err != nil {
			c.logger.WarnContext(ctx, "deconz websocket: undecodable message", "error", err)
			continue
		}
		event.Raw = data
		handle(event)
	}
}

// keepAlive pings until the session ends and cancels it when a ping fails.
func (c *Client) keepAlive(ctx context.Context, conn *websocket.Conn, cancel context.CancelFunc) {
	ticker := time.NewTicker(c.pingInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
		pingCtx, pingCancel := context.WithTimeout(ctx, pingTimeout)
		err := conn.Ping(pingCtx)
		pingCancel()
		if err != nil && ctx.Err() == nil {
			c.logger.WarnContext(ctx, "deconz websocket ping failed", "error", err)
			cancel()
			return
		}
	}
}
