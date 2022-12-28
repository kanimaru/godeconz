package ws

import (
	"github.com/gorilla/websocket"
	"github.com/kanimaru/godeconz"
	"net/url"
	"os"
	"os/signal"
	"time"
)

type AdapterWebsocketClientGorilla struct {
	host       url.URL
	logger     godeconz.Logger
	connection *websocket.Conn
	closed     chan interface{}
}

func CreateAdapterWebsocketClientGorilla(logger godeconz.Logger) Adapter {
	return &AdapterWebsocketClientGorilla{
		logger: logger,
		closed: make(chan interface{}),
	}
}

func (a *AdapterWebsocketClientGorilla) Connect(host url.URL) error {
	a.logger.Infof("Connect to %v", host)
	a.host = host
	connection, _, err := websocket.DefaultDialer.Dial(host.String(), nil)
	a.connection = connection
	if err != nil {
		a.logger.Errorf("Problems with connection: %v", err)
		return err
	}

	defer func() {
		err := a.Close()
		if err != nil {
			a.logger.Errorf("Can't close websocket: %v", err)
		}
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	for {
		a.logger.Infof("Connected to %v")
		select {
		case <-a.closed:
			a.logger.Infof("Websocket closed by program")
		case <-interrupt:
			a.logger.Infof("Websocket closed by OS")
			err := a.CloseGracefully()
			if err != nil {
				a.logger.Errorf("Can't close gracefully:", err)
				panic(err)
			}

			select {
			case <-a.closed:
			case <-time.After(time.Second):
				close(a.closed)
			}
		}
	}
}

func (a *AdapterWebsocketClientGorilla) Close() error {
	a.logger.Infof("Closing websocket to %v", a.host)
	return a.connection.Close()
}

func (a *AdapterWebsocketClientGorilla) CloseGracefully() error {
	a.logger.Infof("Closing websocket to gracefully %v", a.host)
	return a.connection.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
}

func (a *AdapterWebsocketClientGorilla) ReadMessages(handler func([]byte)) {
	defer close(a.closed)
	messageType, data, err := a.connection.ReadMessage()
	if err != nil {
		a.logger.Errorf("Can't read message: %v", err)
	}
	if messageType == websocket.TextMessage {
		a.logger.Debugf("Receiving message %v", string(data))
		handler(data)
	} else {
		a.logger.Debugf("Receiving unknown message")
	}
}
