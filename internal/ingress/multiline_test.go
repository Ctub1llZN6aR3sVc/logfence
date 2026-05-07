package ingress

import (
	"strings"
	"testing"
	"time"
)

func TestMultilineAssembler_SingleLine(t *testing.T) {
	a := NewMultilineAssembler(MultilineConfig{StartPattern: "{"})

	// First line — buffered, not flushed yet.
	entry, ok := a.Feed(`{"level":"info","msg":"hello"}`)
	if ok {
		t.Fatalf("expected no flush on first line, got %q", entry)
	}

	// Second start-pattern line triggers flush of the first.
	entry, ok = a.Feed(`{"level":"warn","msg":"world"}`)
	if !ok {
		t.Fatal("expected flush when new start-pattern line arrives")
	}
	if !strings.Contains(entry, "hello") {
		t.Errorf("flushed entry should contain 'hello', got %q", entry)
	}
}

func TestMultilineAssembler_Continuation(t *testing.T) {
	a := NewMultilineAssembler(MultilineConfig{StartPattern: "ERROR"})

	a.Feed("ERROR something went wrong")
	a.Feed("  at foo.go:42")
	a.Feed("  at bar.go:7")

	entry, ok := a.Flush()
	if !ok {
		t.Fatal("expected flush to return buffered content")
	}
	lines := strings.Split(entry, "\n")
	if len(lines) != 3 {
		t.Errorf("expected 3 lines, got %d: %v", len(lines), lines)
	}
}

func TestMultilineAssembler_MaxLines(t *testing.T) {
	a := NewMultilineAssembler(MultilineConfig{StartPattern: "START", MaxLines: 3})

	a.Feed("START entry")
	a.Feed("line 2")
	entry, ok := a.Feed("line 3") // should auto-flush at MaxLines
	if !ok {
		t.Fatal("expected auto-flush at MaxLines")
	}
	if !strings.Contains(entry, "START entry") {
		t.Errorf("unexpected entry content: %q", entry)
	}
}

func TestMultilineAssembler_MaxAge(t *testing.T) {
	a := NewMultilineAssembler(MultilineConfig{
		StartPattern: "START",
		MaxAge:       10 * time.Millisecond,
	})

	a.Feed("START entry")
	time.Sleep(20 * time.Millisecond)

	entry, ok := a.Feed("continuation") // age exceeded, should flush
	if !ok {
		t.Fatal("expected auto-flush when MaxAge exceeded")
	}
	if !strings.Contains(entry, "START entry") {
		t.Errorf("unexpected entry content: %q", entry)
	}
}

func TestMultilineAssembler_FlushEmpty(t *testing.T) {
	a := NewMultilineAssembler(MultilineConfig{})
	entry, ok := a.Flush()
	if ok || entry != "" {
		t.Errorf("expected empty flush on empty buffer, got ok=%v entry=%q", ok, entry)
	}
}
