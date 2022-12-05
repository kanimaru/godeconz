package deconz

import (
	"github.com/go-resty/resty/v2"
)

type AdapterHttpClientResty struct {
	client *resty.Client
	logger Logger
	trace  bool
}

// CreateAdapterHttpClientResty to plugin deconz api
func CreateAdapterHttpClientResty(client *resty.Client, logger Logger, trace bool) HttpClientAdapter[*resty.Response] {
	return AdapterHttpClientResty{
		client: client,
		logger: logger,
		trace:  trace,
	}
}

func (c AdapterHttpClientResty) Get(path string, container interface{}) (*resty.Response, error) {
	r := c.client.R()
	if c.trace {
		r = r.EnableTrace()
	}
	response, err := r.
		SetResult(&container).
		Get(path)
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

func (c AdapterHttpClientResty) Post(path string, data interface{}) (*resty.Response, error) {
	r := c.client.R()
	if c.trace {
		r = r.EnableTrace()
	}
	if data != nil {
		r.SetBody(data)
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

func (c AdapterHttpClientResty) Put(path string, data interface{}) (*resty.Response, error) {
	r := c.client.R()
	if c.trace {
		r = r.EnableTrace()
	}
	if data != nil {
		r.SetBody(data)
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

func (c AdapterHttpClientResty) Delete(path string, data interface{}) (*resty.Response, error) {
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
