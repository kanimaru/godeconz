package godeconz

import (
	"github.com/gorilla/websocket"
	"net/url"
)

type AdapterWebsocketClientGorilla struct {
	host       url.URL
	logger     Logger
	connection *websocket.Conn
}

func CreateAdapterWebsocketClientGorilla(logger Logger) WebsocketAdapter {
	return &AdapterWebsocketClientGorilla{
		logger: logger,
	}
}

func (a *AdapterWebsocketClientGorilla) Connect(host url.URL) error {
	a.logger.Infof("Connect to %v", host)
	a.host = host
	connection, _, err := websocket.DefaultDialer.Dial(host.String(), nil)
	a.connection = connection
	if err != nil {
		a.logger.Errorf("Problems with connection: %v", err)
	}
	return err
}

func (a *AdapterWebsocketClientGorilla) Close() error {
	a.logger.Infof("Closing websocket to %v", a.host)
	return a.connection.Close()
}

func (a *AdapterWebsocketClientGorilla) CloseGracefully() error {
	return a.connection.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
}

func (a *AdapterWebsocketClientGorilla) ReadMessage() ([]byte, error) {
	messageType, data, err := a.connection.ReadMessage()
	if messageType == websocket.TextMessage {
		return data, err
	} else {
		return nil, nil
	}
}
