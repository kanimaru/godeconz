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

// Callback that get called when websocket receives a message
type Callback interface {
	OnMessage(message Message)
}

type CallbackConfig struct {
	filter   Filter
	callback Callback
}

type Adapter interface {
	// Connect to the given host
	Connect(host url.URL) error
	// Close the connection dirty
	Close() error
	// ReadMessages of the buffer
	ReadMessages(handler func([]byte))
}

func CreateClient(adapter Adapter, logger godeconz.Logger) *Client {
	return &Client{
		adapter: adapter,
		logger:  logger,
	}
}

func (c *Client) Connect(host url.URL) {
	go c.adapter.ReadMessages(c.handleMessages)
	err := c.adapter.Connect(host)
	if err != nil {
		panic(err)
	}
}

func (c *Client) AddCallback(callback Callback, filter Filter) {
	c.callbacks.PushBack(CallbackConfig{
		filter:   filter,
		callback: callback,
	})
}

func (c *Client) RemoveCallback(callback Callback) bool {
	for el := c.callbacks.Front(); el != nil; el = el.Next() {
		if el.Value.(CallbackConfig).callback == callback {
			c.callbacks.Remove(el)
			return true
		}
	}
	return false
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
