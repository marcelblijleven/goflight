package goflight_test

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
)

// HTTPTestClient returns a http client which can be used to stub requests
// it will also return a method to close the internal test server
func HTTPTestClient(handler http.Handler) (httpTestClient *http.Client, close func()) {
	server := httptest.NewServer(handler)
	client := &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, network, _ string) (net.Conn, error) {
				return net.Dial(network, server.Listener.Addr().String())
			},
		},
	}

	return client, server.Close
}

// CreateTestHandler returns a handler that returns the provided status and body
func CreateTestHandler(status int, body []byte) http.HandlerFunc {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		w.Header().Set("provided-timestamp", r.URL.Query().Get("time"))
		w.Header().Set("provided-icao24", r.URL.Query().Get("icao24"))
		w.Write(body)
	})

	return handler
}
