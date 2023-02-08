package ws

import (
	"github.com/PerformLine/go-stockutil/log"
	"time"
)

type StateCallback[S any] struct {
	stateChan chan S
	handler   CallbackHandler
}

// NewStateCallback registers a new callback that listen for the filtered state events and send the state to a chan.
func NewStateCallback[S any](stateChan chan S, handler CallbackHandler, filter Filter) *StateCallback[S] {
	cb := StateCallback[S]{
		stateChan: stateChan,
		handler:   handler,
	}
	filter.HasState = true
	handler.AddCallback(cb, filter)
	return &cb
}

func (s StateCallback[S]) OnMessage(message Message) {
	var state S
	err := message.StateAs(&state)
	log.Fatal("Can't convert state: %+v", err)
	select {
	case s.stateChan <- state:
		// Message send successfully
	case <-time.After(1 * time.Second):
		// Timout
		log.Warningf("State Chan of actor %v got not handled data lost!", message.Name)
	default:
		// Channel got closed
		s.handler.RemoveCallback(s)
	}
}
