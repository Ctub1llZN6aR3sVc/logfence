package sink

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSplunkSink_Write(t *testing.T) {
	var received splunkEvent
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Splunk test-token" {
			t.Errorf("missing or wrong Authorization header: %s", r.Header.Get("Authorization"))
		}
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	s, err := NewSplunkSink(srv.URL, "test-token")
	if err != nil {
		t.Fatalf("NewSplunkSink: %v", err)
	}
	entry := map[string]any{"level": "info", "msg": "hello splunk"}
	if err := s.Write(entry); err != nil {
		t.Fatalf("Write: %v", err)
	}
	if received.Event["msg"] != "hello splunk" {
		t.Errorf("expected msg 'hello splunk', got %v", received.Event["msg"])
	}
	if received.Time == 0 {
		t.Error("expected non-zero time")
	}
}

func TestSplunkSink_Non2xx(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer srv.Close()

	s, _ := NewSplunkSink(srv.URL, "tok")
	if err := s.Write(map[string]any{"msg": "x"}); err == nil {
		t.Fatal("expected error on non-2xx response")
	}
}

func TestSplunkSink_EmptyURL(t *testing.T) {
	_, err := NewSplunkSink("", "tok")
	if err == nil {
		t.Fatal("expected error for empty URL")
	}
}

func TestSplunkSink_EmptyToken(t *testing.T) {
	_, err := NewSplunkSink("http://splunk:8088", "")
	if err == nil {
		t.Fatal("expected error for empty token")
	}
}

func TestSplunkSink_Close(t *testing.T) {
	s, _ := NewSplunkSink("http://splunk:8088", "tok")
	if err := s.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
}

func TestNew_Splunk(t *testing.T) {
	s, err := New("splunk", map[string]string{
		"url":   "http://splunk:8088/services/collector",
		"token": "abc123",
	})
	if err != nil {
		t.Fatalf("New splunk: %v", err)
	}
	if _, ok := s.(*SplunkSink); !ok {
		t.Fatalf("expected *SplunkSink, got %T", s)
	}
}
