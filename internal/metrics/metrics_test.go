package metrics_test

import (
	"bytes"
	"strings"
	"sync"
	"testing"

	"github.com/yourorg/logfence/internal/metrics"
)

func TestCounters_InitialZero(t *testing.T) {
	c := metrics.New()
	for k, v := range c.Snapshot() {
		if v != 0 {
			t.Errorf("expected counter %q to be 0, got %d", k, v)
		}
	}
}

func TestCounters_Inc(t *testing.T) {
	c := metrics.New()
	c.IncReceived()
	c.IncReceived()
	c.IncRouted()
	c.IncDropped()
	c.IncDecodeError()
	c.IncSinkError()

	s := c.Snapshot()
	if s["lines_received"] != 2 {
		t.Errorf("lines_received: want 2, got %d", s["lines_received"])
	}
	if s["lines_routed"] != 1 {
		t.Errorf("lines_routed: want 1, got %d", s["lines_routed"])
	}
	if s["lines_dropped"] != 1 {
		t.Errorf("lines_dropped: want 1, got %d", s["lines_dropped"])
	}
	if s["decode_errors"] != 1 {
		t.Errorf("decode_errors: want 1, got %d", s["decode_errors"])
	}
	if s["sink_errors"] != 1 {
		t.Errorf("sink_errors: want 1, got %d", s["sink_errors"])
	}
}

func TestCounters_ConcurrentInc(t *testing.T) {
	c := metrics.New()
	const goroutines = 50
	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			c.IncReceived()
			c.IncRouted()
		}()
	}
	wg.Wait()
	s := c.Snapshot()
	if s["lines_received"] != goroutines {
		t.Errorf("lines_received: want %d, got %d", goroutines, s["lines_received"])
	}
}

func TestCounters_WriteTo(t *testing.T) {
	c := metrics.New()
	c.IncReceived()
	c.IncRouted()

	var buf bytes.Buffer
	_, err := c.WriteTo(&buf)
	if err != nil {
		t.Fatalf("WriteTo error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"lines_received=1", "lines_routed=1", "lines_dropped=0"} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q; got: %s", want, out)
		}
	}
}
