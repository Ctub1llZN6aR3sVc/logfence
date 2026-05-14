package sink

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestElasticsearchSink_Write(t *testing.T) {
	var got map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &got)
		w.WriteHeader(http.StatusCreated)
	}))
	defer srv.Close()

	s, err := NewElasticsearchSink(srv.URL)
	if err != nil {
		t.Fatalf("NewElasticsearchSink: %v", err)
	}

	entry := map[string]any{"level": "info", "msg": "hello"}
	if err := s.Write(entry); err != nil {
		t.Fatalf("Write: %v", err)
	}
	if got["msg"] != "hello" {
		t.Errorf("expected msg=hello, got %v", got["msg"])
	}
}

func TestElasticsearchSink_Non2xx(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	s, err := NewElasticsearchSink(srv.URL)
	if err != nil {
		t.Fatalf("NewElasticsearchSink: %v", err)
	}
	if err := s.Write(map[string]any{"x": 1}); err == nil {
		t.Fatal("expected error on non-2xx response")
	}
}

func TestElasticsearchSink_EmptyURL(t *testing.T) {
	_, err := NewElasticsearchSink("")
	if err == nil {
		t.Fatal("expected error for empty URL")
	}
}

func TestElasticsearchSink_Close(t *testing.T) {
	s, _ := NewElasticsearchSink("http://localhost:9200/logs/_doc")
	if err := s.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
}

func TestNew_Elasticsearch(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	s, err := New("elasticsearch", map[string]string{"url": srv.URL})
	if err != nil {
		t.Fatalf("New elasticsearch: %v", err)
	}
	if err := s.Write(map[string]any{"level": "debug"}); err != nil {
		t.Fatalf("Write: %v", err)
	}
}
