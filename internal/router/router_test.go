package router_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"testing"

	"github.com/yourorg/logfence/internal/filter"
	"github.com/yourorg/logfence/internal/router"
)

// captureSink records written entries for assertions.
type captureSink struct {
	entries []map[string]interface{}
	closed  bool
	wErr    error
}

func (c *captureSink) Write(entry map[string]interface{}) error {
	if c.wErr != nil {
		return c.wErr
	}
	c.entries = append(c.entries, entry)
	return nil
}
func (c *captureSink) Close() error { c.closed = true; return nil }

func parseChain(t *testing.T, level string) filter.Chain {
	t.Helper()
	lvl, err := filter.ParseLevel(level)
	if err != nil {
		t.Fatalf("ParseLevel: %v", err)
	}
	return filter.Chain{Rules: []filter.Rule{{MinLevel: lvl}}}
}

func TestDispatch_MatchingRoute(t *testing.T) {
	s := &captureSink{}
	r := router.New([]router.Route{
		{Name: "errors", Filter: parseChain(t, "error"), Sink: s},
	})

	entry := map[string]interface{}{"level": "error", "msg": "boom"}
	if err := r.Dispatch(entry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s.entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(s.entries))
	}
}

func TestDispatch_NoMatch(t *testing.T) {
	s := &captureSink{}
	r := router.New([]router.Route{
		{Name: "errors", Filter: parseChain(t, "error"), Sink: s},
	})

	entry := map[string]interface{}{"level": "debug", "msg": "verbose"}
	if err := r.Dispatch(entry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s.entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(s.entries))
	}
}

func TestDispatch_SinkWriteError(t *testing.T) {
	s := &captureSink{wErr: errors.New("disk full")}
	r := router.New([]router.Route{
		{Name: "all", Filter: filter.Chain{}, Sink: s},
	})

	err := r.Dispatch(map[string]interface{}{"level": "info"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestClose_CallsSinkClose(t *testing.T) {
	s := &captureSink{}
	r := router.New([]router.Route{
		{Name: "all", Filter: filter.Chain{}, Sink: s},
	})

	if err := r.Close(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !s.closed {
		t.Fatal("expected sink to be closed")
	}
}

// Ensure the package compiles with json/bytes (used by real sinks).
var _ = json.Marshal
var _ = bytes.NewBuffer
