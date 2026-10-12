package godeconz

import (
	"fmt"
	"net/url"
	"strings"
)

// ParseAddress accepts a gateway address as "host", "host:port" or "http(s)://host[:port]" and
// returns it as URL with scheme and host only. Bare hosts default to http.
func ParseAddress(address string) (*url.URL, error) {
	address = strings.TrimSpace(address)
	if address == "" {
		return nil, fmt.Errorf("deconz: empty gateway address")
	}
	if !strings.Contains(address, "://") {
		address = "http://" + address
	}
	u, err := url.Parse(address)
	if err != nil {
		return nil, fmt.Errorf("deconz: invalid gateway address %q: %w", address, err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return nil, fmt.Errorf("deconz: unsupported scheme %q in gateway address", u.Scheme)
	}
	if u.Host == "" {
		return nil, fmt.Errorf("deconz: gateway address %q has no host", address)
	}
	return &url.URL{Scheme: u.Scheme, Host: u.Host}, nil
}

// DeviceMAC returns the MAC part of a deCONZ uniqueid. All resources of one physical device
// share it, e.g. "00:15:8d:00:06:d4:5b:dd-01-0402" and "...-01-0405" both return
// "00:15:8d:00:06:d4:5b:dd".
func DeviceMAC(uniqueID string) string {
	mac, _, _ := strings.Cut(uniqueID, "-")
	return mac
}
