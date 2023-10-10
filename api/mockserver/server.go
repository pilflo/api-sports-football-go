// Package mockserver aims at mocking the API for unit tests.
// Call the instance with GetServer, add endpoints with AddJSONHandler.
package mockserver

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"sync"
	"testing"
)

var (
	server *httptest.Server
	mux    *http.ServeMux
	resMap = map[string]MockJSONResponse{}
)

// MockJSONResponse represents a fake http response with
//   - a response status code
//   - a user-provided json file path containing the payload.
type MockJSONResponse struct {
	Path         string
	ResponseCode int
	FilePath     string
	QueryParams  *url.Values
}

func initMockServer() {
	mux = http.NewServeMux()
	mux.Handle("/", httpHandler())
	server = httptest.NewServer(mux)
}

// GetServer returns the httptest Server used to fake the API.
func GetServer() *httptest.Server {
	once := sync.Once{}
	once.Do(func() {
		initMockServer()
	})

	return server
}

// AddJSONHandler adds a mocked handler corresponding to the MockJsonResponse.
func AddJSONHandler(t *testing.T, res MockJSONResponse) {
	t.Helper()

	queryPath := res.Path
	if res.QueryParams != nil {
		queryPath += "?" + res.QueryParams.Encode()
	}

	resMap[queryPath] = res
}

//nolint:varnamelen // (pilflo): w and r are common names for an http handler.
func httpHandler() http.HandlerFunc {
	handler := func(w http.ResponseWriter, r *http.Request) {
		queryPath := r.URL.Path
		if r.URL.RawQuery != "" {
			queryPath += "?" + r.URL.RawQuery
		}

		res := resMap[queryPath]

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(res.ResponseCode)

		file, _ := os.Open(res.FilePath)
		defer func() {
			_ = file.Close()
		}()

		bytes, _ := io.ReadAll(file)

		_, _ = w.Write(bytes)
	}

	return http.HandlerFunc(handler)
}
