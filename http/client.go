package http

import (
	"fmt"
	"github.com/kanimaru/godeconz"
)

type ClientAdapter[R any] interface {
	Get(path string, container interface{}) (R, error)
	Post(path string, data interface{}, container interface{}) (R, error)
	Put(path string, data interface{}, container interface{}) (R, error)
	Delete(path string, container interface{}) (R, error)
}

type Client[R any] struct {
	baseUrl string
	adapter ClientAdapter[R]
}

// CreateClient that can access the Deconz API it uses [R] as response type from http requests.
func CreateClient[R any](adapter ClientAdapter[R], settings godeconz.Settings) Client[R] {
	if settings.HttpProtocol == "" {
		settings.HttpProtocol = "http"
	}
	client := Client[R]{
		baseUrl: fmt.Sprintf("%s://%s/api/%s/", settings.HttpProtocol, settings.Address, settings.ApiKey),
		adapter: adapter,
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
	return c.adapter.Get(p, container)
}

func (c *Client[R]) PostWithResult(path string, body interface{}, container interface{}, pathArguments ...any) (R, error) {
	p := c.getPath(path, pathArguments)
	return c.adapter.Post(p, body, container)
}

func (c *Client[R]) Post(path string, body interface{}, pathArguments ...any) (R, error) {
	p := c.getPath(path, pathArguments)
	return c.adapter.Post(p, body, nil)
}

func (c *Client[R]) PutWithResult(path string, body interface{}, container interface{}, pathArguments ...any) (R, error) {
	p := c.getPath(path, pathArguments)
	return c.adapter.Put(p, body, container)
}

func (c *Client[R]) Put(path string, body interface{}, pathArguments ...any) (R, error) {
	p := c.getPath(path, pathArguments)
	return c.adapter.Put(p, body, nil)
}

func (c *Client[R]) Delete(path string, container interface{}, pathArguments ...any) (R, error) {
	p := c.getPath(path, pathArguments)
	return c.adapter.Delete(p, container)
}
