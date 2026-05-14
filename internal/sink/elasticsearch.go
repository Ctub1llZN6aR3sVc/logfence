package sink

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// ElasticsearchSink ships log entries to an Elasticsearch index via the
// Bulk API. Each Write call sends a single index action.
type ElasticsearchSink struct {
	url    string // e.g. http://localhost:9200/logs/_doc
	client *http.Client
}

// NewElasticsearchSink creates a sink that POSTs entries to the given
// Elasticsearch document endpoint. url must be non-empty.
func NewElasticsearchSink(url string) (*ElasticsearchSink, error) {
	if url == "" {
		return nil, fmt.Errorf("elasticsearch: url must not be empty")
	}
	return &ElasticsearchSink{
		url: url,
		client: &http.Client{Timeout: 5 * time.Second},
	}, nil
}

// Write serialises entry and sends it to Elasticsearch.
func (s *ElasticsearchSink) Write(entry map[string]any) error {
	body, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("elasticsearch: marshal: %w", err)
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, s.url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("elasticsearch: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("elasticsearch: send: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("elasticsearch: unexpected status %d", resp.StatusCode)
	}
	return nil
}

// Close is a no-op; the HTTP client does not hold persistent resources.
func (s *ElasticsearchSink) Close() error { return nil }
