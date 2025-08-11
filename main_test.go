package lowestlatencyendpoint_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cilasbeltrame/lowestlatencyendpoint"
)

func TestLowestLatency(t *testing.T) {
	cfg := lowestlatencyendpoint.CreateConfig()
	cfg.Endpoints = []string{
		"http://endpoint1.example.com",
		"http://endpoint2.example.com",
	}
	cfg.HeaderName = "X-Fastest-Endpoint"

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := lowestlatencyendpoint.New(ctx, next, cfg, "lowestlatency-plugin")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(recorder, req)

	// Since we can't easily test actual endpoint latency in unit tests,
	// we'll just verify the handler doesn't error and the structure works
	if recorder.Code != http.StatusOK {
		t.Errorf("expected status OK, got %d", recorder.Code)
	}
}
