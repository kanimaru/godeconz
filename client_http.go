package godeconz

import (
	"fmt"
)

type HttpClientAdapter[R any] interface {
	Get(path string, container interface{}) (R, error)
	Post(path string, data interface{}, container interface{}) (R, error)
	Put(path string, data interface{}, container interface{}) (R, error)
	Delete(path string, container interface{}) (R, error)
}

type HttpClient[R any] struct {
	baseUrl    string
	httpClient HttpClientAdapter[R]
}

// CreateClient that can access the Deconz API it uses [R] as response type from http requests.
func CreateClient[R any](httpClient HttpClientAdapter[R], settings Settings) HttpClient[R] {
	if settings.HttpProtocol == "" {
		settings.HttpProtocol = "http"
	}
	client := HttpClient[R]{
		baseUrl:    fmt.Sprintf("%s://%s/api/%s/", settings.HttpProtocol, settings.Address, settings.ApiKey),
		httpClient: httpClient,
	}
	return client
}

func (c *HttpClient[R]) getPath(path string, arguments []any) string {
	path = fmt.Sprintf(path, arguments...)
	url := c.baseUrl + path
	return url
}

func (c *HttpClient[R]) Get(path string, container interface{}, pathArguments ...any) (R, error) {
	p := c.getPath(path, pathArguments)
	return c.httpClient.Get(p, container)
}

func (c *HttpClient[R]) PostWithResult(path string, body interface{}, container interface{}, pathArguments ...any) (R, error) {
	p := c.getPath(path, pathArguments)
	return c.httpClient.Post(p, body, container)
}

func (c *HttpClient[R]) Post(path string, body interface{}, pathArguments ...any) (R, error) {
	p := c.getPath(path, pathArguments)
	return c.httpClient.Post(p, body, nil)
}

func (c *HttpClient[R]) PutWithResult(path string, body interface{}, container interface{}, pathArguments ...any) (R, error) {
	p := c.getPath(path, pathArguments)
	return c.httpClient.Put(p, body, container)
}

func (c *HttpClient[R]) Put(path string, body interface{}, pathArguments ...any) (R, error) {
	p := c.getPath(path, pathArguments)
	return c.httpClient.Put(p, body, nil)
}

func (c *HttpClient[R]) Delete(path string, container interface{}, pathArguments ...any) (R, error) {
	p := c.getPath(path, pathArguments)
	return c.httpClient.Delete(p, container)
}
