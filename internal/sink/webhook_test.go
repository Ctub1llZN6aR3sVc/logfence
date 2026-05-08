package sink_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yourorg/logfence/internal/sink"
)

func TestWebhookSink_Write(t *testing.T) {
	received := make(chan map[string]any, 1)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var m map[string]any
		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
			t.Errorf("decode body: %v", err)
		}
		received <- m
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	ws, err := sink.NewWebhookSink(sink.WebhookConfig{URL: ts.URL, TimeoutSeconds: 2})
	if err != nil {
		t.Fatalf("new webhook sink: %v", err)
	}

	entry := map[string]any{"level": "info", "msg": "hello"}
	if err := ws.Write(entry); err != nil {
		t.Fatalf("write: %v", err)
	}
	got := <-received
	if got["msg"] != "hello" {
		t.Errorf("expected msg=hello, got %v", got["msg"])
	}
}

func TestWebhookSink_Non2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	ws, _ := sink.NewWebhookSink(sink.WebhookConfig{URL: ts.URL})
	if err := ws.Write(map[string]any{"x": 1}); err == nil {
		t.Fatal("expected error for non-2xx response")
	}
}

func TestWebhookSink_EmptyURL(t *testing.T) {
	_, err := sink.NewWebhookSink(sink.WebhookConfig{})
	if err == nil {
		t.Fatal("expected error for empty URL")
	}
}

func TestWebhookSink_Close(t *testing.T) {
	ws, _ := sink.NewWebhookSink(sink.WebhookConfig{URL: "http://localhost"})
	if err := ws.Close(); err != nil {
		t.Errorf("close should be no-op, got: %v", err)
	}
}

func TestNew_Webhook(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	s, err := sink.New(sink.Config{Type: "webhook", WebhookURL: ts.URL})
	if err != nil {
		t.Fatalf("New webhook: %v", err)
	}
	if err := s.Write(map[string]any{"msg": "test"}); err != nil {
		t.Errorf("write: %v", err)
	}
}
