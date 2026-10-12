package rest

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// deCONZ error types, see https://dresden-elektronik.github.io/deconz-rest-doc/errors/
const (
	ErrorTypeUnauthorized        = 1
	ErrorTypeInvalidJSON         = 2
	ErrorTypeResourceUnavailable = 3
	ErrorTypeMethodUnavailable   = 4
	ErrorTypeMissingParameter    = 5
	ErrorTypeParameterNotAvail   = 6
	ErrorTypeInvalidValue        = 7
	ErrorTypeReadOnly            = 8
	ErrorTypeLinkButtonNotActive = 101
	ErrorTypeDeviceOff           = 201
	ErrorTypeInternal            = 901
)

var (
	// ErrNotFound matches errors for resources that do not exist.
	ErrNotFound = errors.New("deconz: resource not available")
	// ErrUnauthorized matches errors caused by a missing or invalid API key.
	ErrUnauthorized = errors.New("deconz: unauthorized")
	// ErrLinkButton matches AcquireAPIKey failures because the gateway is not unlocked.
	ErrLinkButton = errors.New("deconz: link button not pressed")
)

// APIError is one error object returned by deCONZ.
type APIError struct {
	Type        int    `json:"type"`
	Address     string `json:"address"`
	Description string `json:"description"`
}

// Error is returned when deCONZ answers with a non-2xx status or reports errors in its
// response body. PUT requests can partly succeed; Errors then lists the failed parts.
type Error struct {
	Method     string
	Path       string
	StatusCode int
	Errors     []APIError
	// Body holds the start of the response when it contained no deCONZ error objects.
	Body string
}

func (e *Error) Error() string {
	var b strings.Builder
	fmt.Fprintf(&b, "deconz: %s %s: status %d", e.Method, e.Path, e.StatusCode)
	for _, apiErr := range e.Errors {
		fmt.Fprintf(&b, "; %s (type %d)", apiErr.Description, apiErr.Type)
	}
	if len(e.Errors) == 0 && e.Body != "" {
		fmt.Fprintf(&b, ": %s", e.Body)
	}
	return b.String()
}

// Is lets errors.Is match ErrNotFound, ErrUnauthorized and ErrLinkButton.
func (e *Error) Is(target error) bool {
	switch target {
	case ErrNotFound:
		return e.StatusCode == http.StatusNotFound || e.hasType(ErrorTypeResourceUnavailable)
	case ErrUnauthorized:
		return e.StatusCode == http.StatusUnauthorized || e.hasType(ErrorTypeUnauthorized)
	case ErrLinkButton:
		return e.hasType(ErrorTypeLinkButtonNotActive)
	}
	return false
}

func (e *Error) hasType(t int) bool {
	for _, apiErr := range e.Errors {
		if apiErr.Type == t {
			return true
		}
	}
	return false
}

const maxErrorBody = 200

func newError(method, path string, status int, body []byte) *Error {
	e := &Error{Method: method, Path: path, StatusCode: status, Errors: findAPIErrors(body)}
	if len(e.Errors) == 0 {
		text := strings.TrimSpace(string(body))
		if len(text) > maxErrorBody {
			text = text[:maxErrorBody] + "…"
		}
		e.Body = text
	}
	return e
}

// findAPIErrors extracts error objects from a deCONZ result array such as
// [{"success":{...}},{"error":{...}}]. Anything else yields nil.
func findAPIErrors(body []byte) []APIError {
	body = bytes.TrimSpace(body)
	if len(body) == 0 || body[0] != '[' {
		return nil
	}
	var items []struct {
		Error *APIError `json:"error"`
	}
	if err := json.Unmarshal(body, &items); err != nil {
		return nil
	}
	var found []APIError
	for _, item := range items {
		if item.Error != nil {
			found = append(found, *item.Error)
		}
	}
	return found
}
