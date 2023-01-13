package http

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestAdapterHttpClientResty_Get(t *testing.T) {
	t.Run("Test Etag", func(t *testing.T) {
		server := MockServer()
		defer server.Close()

		// GIVEN
		c := AdapterResty{
			client: resty.New(),
			logger: NewMockLogger(),
			trace:  false,
			cache:  make(map[string]*EtagCacheEntry),
		}

		// WHEN
		data := MockData{Val: false}
		_, _ = c.Get(server.URL+"/data", &data)
		response, _ := c.Get(server.URL+"/data", &data)

		// THEN
		expected := MockData{Val: true}
		if !reflect.DeepEqual(data, expected) {
			t.Errorf("Get() got = %v, want %v", data, expected)
		}
		assert.Equal(t, http.StatusNotModified, response.StatusCode())
	})
}

type MockData struct {
	Val bool `json:"Val"`
}

func MockServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch strings.TrimSpace(r.URL.Path) {
		case "/data":
			ReturnMockData(w, r)
		default:
			http.NotFoundHandler().ServeHTTP(w, r)
		}
	}))
}

func ReturnMockData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("ETag", "123")
	w.Header().Set("Content-Type", "application/json")
	if r.Header.Get("If-None-Match") == "123" {
		w.WriteHeader(http.StatusNotModified)
		return
	}
	w.WriteHeader(http.StatusOK)

	data := MockData{Val: true}
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		panic(err)
	}
}

type MockLogger struct {
	logger *log.Logger
}

func NewMockLogger() *MockLogger {
	return &MockLogger{logger: log.New(os.Stdout, "", 5)}
}

func (l *MockLogger) Debugf(format string, args ...interface{}) {
	l.logger.Printf("[Debug]"+format, args...)
}

func (l *MockLogger) Errorf(format string, args ...interface{}) {
	l.logger.Printf("[Error]"+format, args...)
}

func (l *MockLogger) Warnf(format string, args ...interface{}) {
	l.logger.Printf("[Warn]"+format, args...)
}

func (l *MockLogger) Infof(format string, args ...interface{}) {
	l.logger.Printf("[Info]"+format, args...)
}
