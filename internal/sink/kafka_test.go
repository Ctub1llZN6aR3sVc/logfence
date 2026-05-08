package sink

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestKafkaSink_Write(t *testing.T) {
	var received kafkaPayload
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		body, _ := io.ReadAll(r.Body)
		if err := json.Unmarshal(body, &received); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	k, err := NewKafkaRestSink(srv.URL, "logs", time.Second)
	if err != nil {
		t.Fatalf("new: %v", err)
	}
	entry := map[string]interface{}{"level": "info", "msg": "hello"}
	if err := k.Write(entry); err != nil {
		t.Fatalf("write: %v", err)
	}
	if len(received.Records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(received.Records))
	}
	if received.Records[0].Value["msg"] != "hello" {
		t.Errorf("unexpected value: %v", received.Records[0].Value)
	}
}

func TestKafkaSink_Non2xx(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer srv.Close()

	k, _ := NewKafkaRestSink(srv.URL, "logs", time.Second)
	if err := k.Write(map[string]interface{}{"x": 1}); err == nil {
		t.Error("expected error on non-2xx response")
	}
}

func TestKafkaSink_EmptyURL(t *testing.T) {
	if _, err := NewKafkaRestSink("", "logs", time.Second); err == nil {
		t.Error("expected error for empty URL")
	}
}

func TestKafkaSink_EmptyTopic(t *testing.T) {
	if _, err := NewKafkaRestSink("http://localhost", "", time.Second); err == nil {
		t.Error("expected error for empty topic")
	}
}

func TestKafkaSink_Close(t *testing.T) {
	k, _ := NewKafkaRestSink("http://localhost", "logs", time.Second)
	if err := k.Close(); err != nil {
		t.Errorf("close: %v", err)
	}
}

func TestNew_Kafka(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	s, err := New(map[string]interface{}{
		"type":      "kafka",
		"base_url":  srv.URL,
		"topic":     "events",
		"timeout_s": 2,
	})
	if err != nil {
		t.Fatalf("New kafka: %v", err)
	}
	if err := s.Write(map[string]interface{}{"msg": "test"}); err != nil {
		t.Errorf("write: %v", err)
	}
	_ = s.Close()
}
