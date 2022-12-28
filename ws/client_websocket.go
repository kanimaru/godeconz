package ws

import (
	"container/list"
	"encoding/json"
	"github.com/kanimaru/godeconz"
	"net/url"
)

type Client struct {
	callbacks list.List
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
func CreateClient(adapter Adapter, logger godeconz.Logger) *Client {
	return &Client{
		adapter: adapter,
		logger:  logger,
	}
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
