package sink

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestStdoutSink_Write(t *testing.T) {
	var buf bytes.Buffer
	s := &StdoutSink{w: &buf}

	if err := s.Write([]byte(`{"level":"info","msg":"hello"}`)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := buf.String(); got != "{\"level\":\"info\",\"msg\":\"hello\"}\n" {
		t.Errorf("unexpected output: %q", got)
	}
}

func TestStdoutSink_Close(t *testing.T) {
	s := NewStdoutSink()
	if err := s.Close(); err != nil {
		t.Fatalf("Close() returned error: %v", err)
	}
}

func TestFileSink_WriteAndClose(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.log")

	s, err := NewFileSink(path)
	if err != nil {
		t.Fatalf("NewFileSink: %v", err)
	}

	entry := []byte(`{"level":"error","msg":"boom"}`)
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
	expected := string(entry) + "\n"
	if string(data) != expected {
		t.Errorf("file content = %q, want %q", string(data), expected)
	}
}

func TestNew_Stdout(t *testing.T) {
	s, err := New("stdout", "")
	if err != nil {
		t.Fatalf("New stdout: %v", err)
	}
	if _, ok := s.(*StdoutSink); !ok {
		t.Errorf("expected *StdoutSink, got %T", s)
	}
}

func TestNew_File(t *testing.T) {
	dir := t.TempDir()
	s, err := New("file", filepath.Join(dir, "out.log"))
	if err != nil {
		t.Fatalf("New file: %v", err)
	}
	defer s.Close()
	if _, ok := s.(*FileSink); !ok {
		t.Errorf("expected *FileSink, got %T", s)
	}
}

func TestNew_FileMissingPath(t *testing.T) {
	_, err := New("file", "")
	if err == nil {
		t.Error("expected error for file sink with empty path")
	}
}

func TestNew_UnknownType(t *testing.T) {
	_, err := New("kafka", "")
	if err == nil {
		t.Error("expected error for unknown sink type")
	}
}
