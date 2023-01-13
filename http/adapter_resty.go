package http

import (
	"github.com/go-resty/resty/v2"
	"github.com/kanimaru/godeconz"
	"reflect"
)

type AdapterResty struct {
	client *resty.Client
	logger godeconz.Logger
	trace  bool
	cache  map[string]*EtagCacheEntry
}

// CreateAdapterHttpClientResty to plugin deconz api
func CreateAdapterHttpClientResty(client *resty.Client, logger godeconz.Logger, trace bool) ClientAdapter[*resty.Response] {
	return AdapterResty{
		client: client,
		logger: logger,
		trace:  trace,
		cache:  make(map[string]*EtagCacheEntry),
	}
}

func (c AdapterResty) Get(path string, container interface{}) (*resty.Response, error) {
	r := c.client.R()
	if c.trace {
		r = r.EnableTrace()
	}

	content, ok := c.cache[path]
	if ok {
		r.SetHeader("If-None-Match", content.etag)
	}

	response, err := r.
		SetResult(container).
		Get(path)
	if err != nil {
		return response, err
	}

	if c.HandleEtag(response, path, container) {
		c.logger.Debugf("Request cached")
	}

	if response.IsError() {
		c.logger.Errorf(response.String())
	} else {
		c.logger.Debugf("Request successfully")
	}
	return response, nil
}

func (c AdapterResty) HandleEtag(response *resty.Response, path string, container interface{}) bool {
	entry, ok := c.cache[path]

	if response.StatusCode() == 304 {
		data := reflect.Indirect(reflect.ValueOf(entry.content))
		source := reflect.Indirect(reflect.ValueOf(container))
		source.Set(data)
		return true
	}

	etag := response.Header().Get("ETag")
	if etag != "" {
		if !ok {
			entry = &EtagCacheEntry{
				etag: etag,
			}
			c.cache[path] = entry
		}
		entry.content = container
	}
	return false
}

func (c AdapterResty) Post(path string, data interface{}, container interface{}) (*resty.Response, error) {
	r := c.client.R()
	if c.trace {
		r = r.EnableTrace()
	}
	if data != nil {
		r.SetBody(data)
	}
	if container != nil {
		r.SetResult(container)
	}
	response, err := r.
		Put(path)
	if err != nil {
		return response, err
	}
	if response.IsError() {
		c.logger.Errorf("Resp: %s", response.String())
	} else {
		c.logger.Debugf("Request successfully")
	}
	return response, nil
}

func (c AdapterResty) Put(path string, data interface{}, container interface{}) (*resty.Response, error) {
	r := c.client.R()
	if c.trace {
		r = r.EnableTrace()
	}
	if data != nil {
		r.SetBody(data)
	}
	if container != nil {
		r.SetResult(container)
	}
	response, err := r.
		Put(path)
	if err != nil {
		return response, err
	}
	if !response.IsSuccess() {
		c.logger.Errorf("Resp: %s", response.String())
	} else {
		c.logger.Debugf("Request successfully")
	}
	return response, nil
}

func (c AdapterResty) Delete(path string, data interface{}) (*resty.Response, error) {
	r := c.client.R()
	if c.trace {
		r = r.EnableTrace()
	}
	if data != nil {
		r.SetBody(data)
	}
	response, err := r.Delete(path)
	if err != nil {
		return response, err
	}

	if response.IsError() {
		c.logger.Errorf(response.String())
	} else {
		c.logger.Debugf("Request successfully")
	}
	return response, nil
}

//
// Etag handling
//

type EtagCacheEntry struct {
	etag    string
	content interface{}
}
