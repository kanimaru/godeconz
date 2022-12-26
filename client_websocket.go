package godeconz

import (
	"github.com/gorilla/websocket"
	"net/url"
	"os"
	"os/signal"
	"time"
)

type WebsocketClient struct {
	callbacks []Callback
	adapter   WebsocketAdapter
	logger    Logger
}

type Callback interface {
}

type WebsocketAdapter interface {
	// Connect to the given host
	Connect(host url.URL) error
	// Close the connection dirty
	Close() error
	// CloseGracefully the connection with clean close message send
	CloseGracefully() error
	// ReadMessage of the buffer
	ReadMessage() ([]byte, error)
}

func CreateWebsocketClient(adapter WebsocketAdapter, logger Logger) *WebsocketClient {
	return &WebsocketClient{
		adapter: adapter,
		logger:  logger,
	}
}

func (w *WebsocketClient) Connect(host url.URL) error {
	err := w.adapter.Connect(host)
	if err != nil {
		return err
	}

	defer func() {
		err := w.adapter.Close()
		if err != nil {
			w.logger.Errorf("Can't close websocket: %v", err)
		}
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	for {
		w.logger.Infof("Connected to %v")
		select {
		case <-closed:
			log.Errorf("closed")
			return nil
		case <-interrupt:
			log.Errorf("interrupt")
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Errorf("write close:", err)
				return
			}
			select {
			case <-closed:
			case <-time.After(time.Second):
				close(closed)
			}
			return
		}
	}
	return nil
}

// TODO
func (w WebsocketClient) HandleMessage() {
	// Read messages
	go func() {
		log.Info("Read messages...")
		defer close(closed)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Errorf("err:", err)
				return
			}
			handler(message)
		}
	}()
}
