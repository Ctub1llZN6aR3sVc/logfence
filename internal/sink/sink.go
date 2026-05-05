package sink

import (
	"fmt"
	"io"
	"os"
)

// Sink represents a log output destination.
type Sink interface {
	Write(entry []byte) error
	Close() error
}

// StdoutSink writes log entries to stdout.
type StdoutSink struct {
	w io.Writer
}

// NewStdoutSink creates a new StdoutSink.
func NewStdoutSink() *StdoutSink {
	return &StdoutSink{w: os.Stdout}
}

func (s *StdoutSink) Write(entry []byte) error {
	_, err := fmt.Fprintf(s.w, "%s\n", entry)
	return err
}

func (s *StdoutSink) Close() error { return nil }

// FileSink writes log entries to a file.
type FileSink struct {
	f *os.File
}

// NewFileSink opens or creates a file sink at the given path.
func NewFileSink(path string) (*FileSink, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("sink: open file %q: %w", path, err)
	}
	return &FileSink{f: f}, nil
}

func (s *FileSink) Write(entry []byte) error {
	_, err := fmt.Fprintf(s.f, "%s\n", entry)
	return err
}

func (s *FileSink) Close() error {
	return s.f.Close()
}

// New constructs a Sink from a type string and optional path.
// Supported types: "stdout", "file".
func New(sinkType, path string) (Sink, error) {
	switch sinkType {
	case "stdout":
		return NewStdoutSink(), nil
	case "file":
		if path == "" {
			return nil, fmt.Errorf("sink: file sink requires a path")
		}
		return NewFileSink(path)
	default:
		return nil, fmt.Errorf("sink: unknown type %q", sinkType)
	}
}
