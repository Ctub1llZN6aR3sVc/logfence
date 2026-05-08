package sink

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// KafkaRestSink publishes log entries to a Kafka topic via the Confluent REST
// Proxy (or any compatible HTTP/JSON Kafka proxy endpoint).
type KafkaRestSink struct {
	url    string
	topic  string
	client *http.Client
}

type kafkaRecord struct {
	Value map[string]interface{} `json:"value"`
}

type kafkaPayload struct {
	Records []kafkaRecord `json:"records"`
}

// NewKafkaRestSink creates a KafkaRestSink that POSTs to baseURL/topics/<topic>.
func NewKafkaRestSink(baseURL, topic string, timeout time.Duration) (*KafkaRestSink, error) {
	if baseURL == "" {
		return nil, fmt.Errorf("kafka: base URL must not be empty")
	}
	if topic == "" {
		return nil, fmt.Errorf("kafka: topic must not be empty")
	}
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	return &KafkaRestSink{
		url:    fmt.Sprintf("%s/topics/%s", baseURL, topic),
		topic:  topic,
		client: &http.Client{Timeout: timeout},
	}, nil
}

// Write encodes the entry as a Kafka REST Proxy record and POSTs it.
func (k *KafkaRestSink) Write(entry map[string]interface{}) error {
	payload := kafkaPayload{
		Records: []kafkaRecord{{Value: entry}},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("kafka: marshal: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, k.url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("kafka: new request: %w", err)
	}
	req.Header.Set("Content-Type", "application/vnd.kafka.json.v2+json")
	resp, err := k.client.Do(req)
	if err != nil {
		return fmt.Errorf("kafka: post: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("kafka: unexpected status %d", resp.StatusCode)
	}
	return nil
}

// Close is a no-op for the HTTP-backed sink.
func (k *KafkaRestSink) Close() error { return nil }
