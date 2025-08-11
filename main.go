// Package lowestlatencyendpoint is a plugin that sets the header with the lowest latency endpoint.
package lowestlatencyendpoint

import (
	"context"
	"net/http"
	"time"
)

// Config holds the plugin configuration.
type Config struct {
	Endpoints  []string `json:"endpoints,omitempty"`
	HeaderName string   `json:"headerName,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		Endpoints:  []string{},
		HeaderName: "X-Lowest-Latency",
	}
}

// LowestLatency middleware structure.
type LowestLatency struct {
	next       http.Handler
	endpoints  []string
	headerName string
}

// New creates a new instance of the middleware.
func New(_ context.Context, next http.Handler, config *Config, _ string) (http.Handler, error) {
	return &LowestLatency{
		next:       next,
		endpoints:  config.Endpoints,
		headerName: config.HeaderName,
	}, nil
}

func (m *LowestLatency) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	fastest := m.checkFastestEndpoint(req.Context())
	if fastest != "" {
		req.Header.Set(m.headerName, fastest)
	}
	m.next.ServeHTTP(rw, req)
}

type result struct {
	url     string
	latency time.Duration
}

func (m *LowestLatency) checkFastestEndpoint(ctx context.Context) string {
	if len(m.endpoints) == 0 {
		return ""
	}

	results := make(chan result, len(m.endpoints))
	client := m.createHTTPClient()

	for _, ep := range m.endpoints {
		go m.checkEndpointLatency(ctx, client, ep, results)
	}

	fastest := m.findFastestResult(results)
	if fastest.latency == time.Hour {
		return ""
	}

	return fastest.url
}

func (m *LowestLatency) createHTTPClient() *http.Client {
	return &http.Client{
		Timeout:       2 * time.Second,
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           nil,
	}
}

func (m *LowestLatency) checkEndpointLatency(ctx context.Context, client *http.Client, endpoint string, results chan<- result) {
	start := time.Now()
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, endpoint, http.NoBody)
	if err != nil {
		results <- result{url: endpoint, latency: time.Hour}

		return
	}

	resp, err := client.Do(req)
	if err != nil {
		results <- result{url: endpoint, latency: time.Hour}

		return
	}

	_ = resp.Body.Close() //nolint:errcheck // Ignore close error
	results <- result{url: endpoint, latency: time.Since(start)}
}

func (m *LowestLatency) findFastestResult(results <-chan result) result {
	fastest := result{url: "", latency: time.Hour}
	for i := 0; i < len(m.endpoints); i++ {
		r := <-results
		if r.latency < fastest.latency {
			fastest = r
		}
	}

	return fastest
}
