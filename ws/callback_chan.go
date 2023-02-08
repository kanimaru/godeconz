package ws

import (
	"github.com/PerformLine/go-stockutil/log"
	"time"
)

type ChanCallback struct {
	handler     CallbackHandler
	messageChan chan Message
}

// NewChanCallback registers a new callback that listen for the filtered events and send the message to a chan.
func NewChanCallback(messageChan chan Message, handler CallbackHandler, filter Filter) *ChanCallback {
	cb := ChanCallback{
		handler:     handler,
		messageChan: messageChan,
	}
	handler.AddCallback(cb, filter)
	return &cb
}

func (c ChanCallback) OnMessage(message Message) {
	select {
	case c.messageChan <- message:
		// Message send successfully
	case <-time.After(1 * time.Second):
		// Timout
		log.Warningf("State Chan of actor %v got not handled data lost!", message.Name)
	default:
		// Channel got closed
		c.handler.RemoveCallback(c)
	}
}
