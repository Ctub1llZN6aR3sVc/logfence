package sink

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLokiSink_Write(t *testing.T) {
	var received lokiStream
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	sink, err := NewLokiSink(ts.URL, map[string]string{"app": "logfence"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entry := map[string]any{"level": "info", "msg": "hello loki"}
	if err := sink.Write(entry); err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	if len(received.Streams) != 1 {
		t.Fatalf("expected 1 stream, got %d", len(received.Streams))
	}
	if received.Streams[0].Stream["app"] != "logfence" {
		t.Errorf("expected label app=logfence, got %v", received.Streams[0].Stream)
	}
	if len(received.Streams[0].Values) != 1 {
		t.Errorf("expected 1 log value, got %d", len(received.Streams[0].Values))
	}
}

func TestLokiSink_Non2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	sink, _ := NewLokiSink(ts.URL, nil)
	err := sink.Write(map[string]any{"msg": "fail"})
	if err == nil {
		t.Fatal("expected error on non-2xx response")
	}
}

func TestLokiSink_EmptyURL(t *testing.T) {
	_, err := NewLokiSink("", nil)
	if err == nil {
		t.Fatal("expected error for empty URL")
	}
}

func TestLokiSink_Close(t *testing.T) {
	sink, _ := NewLokiSink("http://localhost:3100", nil)
	if err := sink.Close(); err != nil {
		t.Errorf("Close returned unexpected error: %v", err)
	}
}

func TestNew_Loki(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	s, err := New("loki", map[string]any{
		"url":    ts.URL,
		"labels": map[string]string{"env": "test"},
	})
	if err != nil {
		t.Fatalf("New loki sink: %v", err)
	}
	if err := s.Write(map[string]any{"msg": "via factory"}); err != nil {
		t.Fatalf("Write via factory: %v", err)
	}
}
