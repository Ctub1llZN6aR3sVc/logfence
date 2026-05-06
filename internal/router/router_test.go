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
	if !errors.Is(err, s.wErr) {
		t.Fatalf("expected underlying error %q, got %q", s.wErr, err)
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

// TestDispatch_MultipleRoutes verifies that an entry matching more than one
// route is delivered to all matching sinks.
func TestDispatch_MultipleRoutes(t *testing.T) {
	s1 := &captureSink{}
	s2 := &captureSink{}
	r := router.New([]router.Route{
		{Name: "warnings-and-above", Filter: parseChain(t, "warn"), Sink: s1},
		{Name: "errors-only", Filter: parseChain(t, "error"), Sink: s2},
	})

	entry := map[string]interface{}{"level": "error", "msg": "critical failure"}
	if err := r.Dispatch(entry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s1.entries) != 1 {
		t.Fatalf("s1: expected 1 entry, got %d", len(s1.entries))
	}
	if len(s2.entries) != 1 {
		t.Fatalf("s2: expected 1 entry, got %d", len(s2.entries))
	}
}

// Ensure the package compiles with json/bytes (used by real sinks).
var _ = json.Marshal
var _ = bytes.NewBuffer
