package sink

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestRotatingSink_WriteAndClose(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.log")

	s, err := NewRotatingFileSink(path, 4096)
	if err != nil {
		t.Fatalf("NewRotatingFileSink: %v", err)
	}
	defer s.Close()

	entry := map[string]any{"level": "info", "msg": "hello"}
	if err := s.Write(entry); err != nil {
		t.Fatalf("Write: %v", err)
	}

	if err := s.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	var got map[string]any
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal written line: %v", err)
	}
	if got["msg"] != "hello" {
		t.Errorf("expected msg=hello, got %v", got["msg"])
	}
}

func TestRotatingSink_Rotation(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "rotate.log")

	// Use a tiny limit so the second write triggers rotation.
	s, err := NewRotatingFileSink(path, 50)
	if err != nil {
		t.Fatalf("NewRotatingFileSink: %v", err)
	}
	defer s.Close()

	for i := 0; i < 5; i++ {
		if err := s.Write(map[string]any{"i": i, "msg": "rotation test"}); err != nil {
			t.Fatalf("Write %d: %v", i, err)
		}
	}

	if err := s.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	// At least one rotated file should exist alongside the active file.
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("ReadDir: %v", err)
	}
	if len(entries) < 2 {
		t.Errorf("expected at least 2 files (active + rotated), got %d", len(entries))
	}
}

func TestRotatingSink_InvalidMaxBytes(t *testing.T) {
	dir := t.TempDir()
	_, err := NewRotatingFileSink(filepath.Join(dir, "x.log"), 0)
	if err == nil {
		t.Fatal("expected error for maxBytes=0, got nil")
	}
}

func TestRotatingSink_Dir(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "app.log")
	s, err := NewRotatingFileSink(path, 1024)
	if err != nil {
		t.Fatalf("NewRotatingFileSink: %v", err)
	}
	defer s.Close()
	if s.Dir() != dir {
		t.Errorf("Dir() = %q, want %q", s.Dir(), dir)
	}
}
