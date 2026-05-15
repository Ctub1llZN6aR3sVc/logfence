package sink

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// splunkEvent wraps a log entry in the Splunk HEC envelope.
type splunkEvent struct {
	Time  float64        `json:"time"`
	Event map[string]any `json:"event"`
}

// SplunkSink forwards log entries to a Splunk HTTP Event Collector endpoint.
type SplunkSink struct {
	url    string
	token  string
	client *http.Client
}

// NewSplunkSink creates a SplunkSink that posts to the given HEC url using
// the provided token.  url must be non-empty (e.g. "https://splunk:8088/services/collector").
func NewSplunkSink(url, token string) (*SplunkSink, error) {
	if url == "" {
		return nil, fmt.Errorf("splunk: url must not be empty")
	}
	if token == "" {
		return nil, fmt.Errorf("splunk: token must not be empty")
	}
	return &SplunkSink{
		url:   url,
		token: token,
		client: &http.Client{Timeout: 5 * time.Second},
	}, nil
}

// Write encodes entry as a Splunk HEC event and POSTs it.
func (s *SplunkSink) Write(entry map[string]any) error {
	env := splunkEvent{
		Time:  float64(time.Now().UnixNano()) / 1e9,
		Event: entry,
	}
	body, err := json.Marshal(env)
	if err != nil {
		return fmt.Errorf("splunk: marshal: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, s.url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("splunk: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Splunk "+s.token)

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("splunk: post: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("splunk: unexpected status %d", resp.StatusCode)
	}
	return nil
}

// Close is a no-op for the HTTP-based sink.
func (s *SplunkSink) Close() error { return nil }
