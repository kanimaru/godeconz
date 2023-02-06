package ws

import (
	"container/list"
	"encoding/json"
	"github.com/kanimaru/godeconz"
	"github.com/kanimaru/godeconz/http"
	"net/url"
	"strconv"
)

type Client struct {
	callbacks list.List
	url       *url.URL
	adapter   Adapter
	logger    godeconz.Logger
}

type Adapter interface {
	// Connect to the given host
	Connect(host *url.URL, messages chan<- []byte) error
}

// CreateClient creates a websocket client for Phoscon WSS.
// The Websocket functionality is still under development.
// Notably added and deleted notifications might not be issued under all circumstances.
func CreateClient(url *url.URL, adapter Adapter, logger godeconz.Logger) *Client {
	return &Client{
		adapter: adapter,
		logger:  logger,
		url:     url,
	}
}

// CreateClientFromConfig is an alias for CreateClient that uses the config from REST to determine the URL for websocket.
func CreateClientFromConfig[R any](client http.Client[R], adapter Adapter, logger godeconz.Logger) *Client {
	config := http.ConfigResponse{}
	_, err := client.GetConfig(&config)
	if err != nil {
		panic(err)
	}
	wsUrl, err := url.Parse("ws://" + config.Ipaddress + ":" + strconv.Itoa(config.WebsocketPort))
	return CreateClient(wsUrl, adapter, logger)
}

func (c *Client) Connect(host *url.URL) {
	messages := make(chan []byte)
	go func() {
		for msg := range messages {
			c.handleMessages(msg)
		}
	}()
	err := c.adapter.Connect(host, messages)
	if err != nil {
		panic(err)
	}
}

func (c *Client) handleMessages(messageData []byte) {
	message := c.parse(messageData)
	for el := c.callbacks.Front(); el != nil; el = el.Next() {
		cb := el.Value.(CallbackConfig)
		if !cb.filter.check(message) {
			continue
		}
		cb.callback.OnMessage(message)
	}
}

func (c *Client) parse(messageData []byte) Message {
	var message Message
	err := json.Unmarshal(messageData, &message)
	c.logger.Debugf("Read: %+v", message)
	if err != nil {
		c.logger.Errorf("Can't parse message: %v", err)
	}
	return message
}
