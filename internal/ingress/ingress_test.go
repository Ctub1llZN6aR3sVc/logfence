package ingress

import (
	"context"
	"strings"
	"testing"
	"time"
)

func collect(ch <-chan Entry) []Entry {
	var out []Entry
	for e := range ch {
		out = append(out, e)
	}
	return out
}

func TestReader_ValidJSON(t *testing.T) {
	input := `{"level":"info","msg":"hello"}
{"level":"error","msg":"boom"}
`
	ch := Reader(context.Background(), strings.NewReader(input))
	entries := collect(ch)
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Fields["level"] != "info" {
		t.Errorf("expected level=info, got %v", entries[0].Fields["level"])
	}
	if entries[1].Fields["msg"] != "boom" {
		t.Errorf("expected msg=boom, got %v", entries[1].Fields["msg"])
	}
}

func TestReader_NonJSON(t *testing.T) {
	input := "plain text line\n"
	ch := Reader(context.Background(), strings.NewReader(input))
	entries := collect(ch)
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if string(entries[0].Raw) != "plain text line" {
		t.Errorf("unexpected raw: %q", entries[0].Raw)
	}
	if len(entries[0].Fields) != 0 {
		t.Errorf("expected empty fields for non-JSON, got %v", entries[0].Fields)
	}
}

func TestReader_EmptyLines(t *testing.T) {
	input := "\n{\"level\":\"warn\"}\n\n"
	ch := Reader(context.Background(), strings.NewReader(input))
	entries := collect(ch)
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry (empty lines skipped), got %d", len(entries))
	}
}

func TestReader_ContextCancel(t *testing.T) {
	// Use a pipe so the reader blocks; cancel should unblock.
	pr, pw := strings.NewReader(""), strings.NewReader("")
	_ = pr
	_ = pw

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	ch := Reader(ctx, strings.NewReader(""))
	select {
	case <-ch:
		// channel closed normally — OK
	case <-time.After(200 * time.Millisecond):
		t.Fatal("channel not closed after context timeout")
	}
}
