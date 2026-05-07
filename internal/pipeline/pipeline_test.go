package pipeline_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/logfence/internal/ingress"
	"github.com/yourorg/logfence/internal/metrics"
	"github.com/yourorg/logfence/internal/pipeline"
	"github.com/yourorg/logfence/internal/router"
)

func makeReader(lines []string) *ingress.Reader {
	return ingress.NewReader(strings.NewReader(strings.Join(lines, "\n")))
}

func makeRouter(t *testing.T, sink *bytes.Buffer) *router.Router {
	t.Helper()
	// A router that accepts everything and writes to a buffer sink.
	rt, err := router.NewPassthrough(sink)
	if err != nil {
		t.Fatalf("router.NewPassthrough: %v", err)
	}
	return rt
}

func TestPipeline_RoutesEntries(t *testing.T) {
	lines := []string{
		`{"level":"info","msg":"hello"}`,
		`{"level":"warn","msg":"world"}`,
	}
	r := makeReader(lines)
	var buf bytes.Buffer
	rt := makeRouter(t, &buf)
	c := metrics.New()
	p := pipeline.New(r, rt, c, slog.Default())

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := p.Run(ctx); err != nil && err != context.DeadlineExceeded {
		t.Fatalf("Run returned unexpected error: %v", err)
	}

	if got := c.Received(); got != 2 {
		t.Errorf("received = %d, want 2", got)
	}
	if got := c.Routed(); got != 2 {
		t.Errorf("routed = %d, want 2", got)
	}
	if got := c.Dropped(); got != 0 {
		t.Errorf("dropped = %d, want 0", got)
	}
}

func TestPipeline_DropsNonJSON(t *testing.T) {
	lines := []string{"not-json", `{"level":"info","msg":"ok"}`}
	r := makeReader(lines)
	var buf bytes.Buffer
	rt := makeRouter(t, &buf)
	c := metrics.New()
	p := pipeline.New(r, rt, c, slog.Default())

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	p.Run(ctx) //nolint:errcheck

	// Only the valid JSON line should be received.
	if got := c.Received(); got != 1 {
		t.Errorf("received = %d, want 1", got)
	}
}

func TestPipeline_ContextCancel(t *testing.T) {
	// Infinite source via a blocking reader — cancel should stop the pipeline.
	pr, pw := io.Pipe()
	r := ingress.NewReader(pr)
	var buf bytes.Buffer
	rt := makeRouter(t, &buf)
	c := metrics.New()
	p := pipeline.New(r, rt, c, slog.Default())

	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan error, 1)
	go func() { done <- p.Run(ctx) }()

	// Write one entry then cancel.
	entry, _ := json.Marshal(map[string]string{"level": "info", "msg": "hi"})
	pw.Write(append(entry, '\n'))
	time.Sleep(20 * time.Millisecond)
	cancel()
	pw.Close()

	select {
	case err := <-done:
		if err != context.Canceled {
			t.Errorf("expected context.Canceled, got %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("pipeline did not stop after context cancel")
	}
}
