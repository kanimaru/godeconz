package ws

// Callback that get called when websocket receives a message
type Callback interface {
	OnMessage(message Message)
}

type BaseCallback struct {
	Handler func(message Message)
}

func (b BaseCallback) OnMessage(message Message) {
	b.Handler(message)
}

type CallbackConfig struct {
	filter   Filter
	callback Callback
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
