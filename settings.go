package godeconz

type Settings struct {
	// Address the base path of deconz without protocol example: 192.168.178.56
	Address string
	// Protocol that should be used default: http
	HttpProtocol string
	// ApiKey that should be used for authorization
	ApiKey string
}
