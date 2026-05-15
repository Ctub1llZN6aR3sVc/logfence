package sink

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// lokiStream represents a single Loki stream payload.
type lokiStream struct {
	Streams []lokiStreamEntry `json:"streams"`
}

type lokiStreamEntry struct {
	Stream map[string]string `json:"stream"`
	Values [][]string        `json:"values"`
}

// LokiSink pushes log entries to a Grafana Loki push API endpoint.
type LokiSink struct {
	url    string
	labels map[string]string
	client *http.Client
}

// NewLokiSink creates a new LokiSink.
// url should point to the Loki push endpoint, e.g. http://localhost:3100/loki/api/v1/push.
// labels are static key/value pairs attached to every stream.
func NewLokiSink(url string, labels map[string]string) (*LokiSink, error) {
	if url == "" {
		return nil, fmt.Errorf("loki: url must not be empty")
	}
	if labels == nil {
		labels = map[string]string{}
	}
	return &LokiSink{
		url:    url,
		labels: labels,
		client: &http.Client{Timeout: 5 * time.Second},
	}, nil
}

// Write serialises the log entry and pushes it to Loki.
func (s *LokiSink) Write(entry map[string]any) error {
	line, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("loki: marshal entry: %w", err)
	}

	ts := strconv.FormatInt(time.Now().UnixNano(), 10)
	payload := lokiStream{
		Streams: []lokiStreamEntry{
			{
				Stream: s.labels,
				Values: [][]string{{ts, string(line)}},
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("loki: marshal payload: %w", err)
	}

	resp, err := s.client.Post(s.url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("loki: http post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("loki: unexpected status %d", resp.StatusCode)
	}
	return nil
}

// Close is a no-op for LokiSink.
func (s *LokiSink) Close() error { return nil }
