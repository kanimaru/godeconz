package deconz

import (
	"fmt"
)

type Settings struct {
	// Address the base path of deconz without protocol example: 192.168.178.56
	Address string
	// Protocol that should be used default: http
	HttpProtocol string
	// ApiKey that should be used for authorization
	ApiKey string
}

type HttpClientAdapter[R any] interface {
	Get(path string, container interface{}) (R, error)
	Post(path string, container interface{}) (R, error)
	Put(path string, container interface{}) (R, error)
	Delete(path string, container interface{}) (R, error)
}

type Client[R any] struct {
	baseUrl    string
	httpClient HttpClientAdapter[R]
}

// CreateClient that can access the Deconz API it uses [R] as response type from http requests.
func CreateClient[R any](httpClient HttpClientAdapter[R], settings Settings) Client[R] {
	if settings.HttpProtocol == "" {
		settings.HttpProtocol = "http"
	}
	client := Client[R]{
		baseUrl:    fmt.Sprintf("%s://%s/api/%s/", settings.HttpProtocol, settings.Address, settings.ApiKey),
		httpClient: httpClient,
	}
	return client
}

func (c *Client[R]) getPath(path string, arguments []any) string {
	path = fmt.Sprintf(path, arguments...)
	url := c.baseUrl + path
	return url
}

func (c *Client[R]) Get(path string, container interface{}, pathArguments ...any) (R, error) {
	p := c.getPath(path, pathArguments)
	return c.httpClient.Get(p, container)
}

func (c *Client[R]) Post(path string, container interface{}, pathArguments ...any) (R, error) {
	p := c.getPath(path, pathArguments)
	return c.httpClient.Post(p, container)
}

func (c *Client[R]) Put(path string, container interface{}, pathArguments ...any) (R, error) {
	p := c.getPath(path, pathArguments)
	return c.httpClient.Put(p, container)
}

func (c *Client[R]) Delete(path string, container interface{}, pathArguments ...any) (R, error) {
	p := c.getPath(path, pathArguments)
	return c.httpClient.Delete(p, container)
}
