package sink

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// WebhookSink sends log entries as JSON POST requests to an HTTP endpoint.
type WebhookSink struct {
	url    string
	client *http.Client
}

// WebhookConfig holds configuration for a WebhookSink.
type WebhookConfig struct {
	URL            string
	TimeoutSeconds int
}

// NewWebhookSink creates a new WebhookSink.
func NewWebhookSink(cfg WebhookConfig) (*WebhookSink, error) {
	if cfg.URL == "" {
		return nil, fmt.Errorf("webhook url must not be empty")
	}
	timeout := time.Duration(cfg.TimeoutSeconds) * time.Second
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	return &WebhookSink{
		url: cfg.URL,
		client: &http.Client{Timeout: timeout},
	}, nil
}

// Write serialises the entry and POSTs it to the configured URL.
func (w *WebhookSink) Write(entry map[string]any) error {
	body, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("webhook marshal: %w", err)
	}
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, w.url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("webhook request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := w.client.Do(req)
	if err != nil {
		return fmt.Errorf("webhook send: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook non-2xx response: %d", resp.StatusCode)
	}
	return nil
}

// Close is a no-op for WebhookSink.
func (w *WebhookSink) Close() error { return nil }
