package ws

// Callback that get called when websocket receives a message
type Callback interface {
	OnMessage(message Message)
}

type BaseCallback struct {
	Handler func(message Message)
}

// CallbackSendToChan sends messages to the given channel it registers itself as callback and will unregister itself as
// soon as the channel gets closed.
func (c *Client) CallbackSendToChan(messages chan Message, filter Filter) Callback {
	cb := BaseCallback{}
	cb.Handler = func(message Message) {
		select {
		case messages <- message:
			// Message send successfully
		default:
			// Channel got closed
			c.RemoveCallback(cb)
		}
	}
	c.AddCallback(cb, filter)
	return cb
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
