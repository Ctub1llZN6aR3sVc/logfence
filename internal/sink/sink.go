// Package sink provides writers that accept a structured log entry
// (map[string]interface{}) and forward it to a destination.
package sink

import (
	"fmt"
	"os"
	"time"
)

// Sink is the common interface implemented by every destination.
type Sink interface {
	Write(entry map[string]interface{}) error
	Close() error
}

// New constructs a Sink from a generic configuration map.
// Required key: "type" — one of "stdout", "file", "rotating", "webhook", "kafka".
func New(cfg map[string]interface{}) (Sink, error) {
	typ, _ := cfg["type"].(string)
	switch typ {
	case "stdout":
		return NewStdoutSink(), nil
	case "file":
		path, _ := cfg["path"].(string)
		if path == "" {
			return nil, fmt.Errorf("sink file: missing path")
		}
		return NewFileSink(path)
	case "rotating":
		path, _ := cfg["path"].(string)
		if path == "" {
			return nil, fmt.Errorf("sink rotating: missing path")
		}
		maxBytes, _ := cfg["max_bytes"].(int)
		maxFiles, _ := cfg["max_files"].(int)
		return NewRotatingFileSink(path, int64(maxBytes), maxFiles)
	case "webhook":
		url, _ := cfg["url"].(string)
		timeoutS, _ := cfg["timeout_s"].(int)
		if timeoutS <= 0 {
			timeoutS = 5
		}
		return NewWebhookSink(url, time.Duration(timeoutS)*time.Second)
	case "kafka":
		baseURL, _ := cfg["base_url"].(string)
		topic, _ := cfg["topic"].(string)
		timeoutS, _ := cfg["timeout_s"].(int)
		if timeoutS <= 0 {
			timeoutS = 5
		}
		return NewKafkaRestSink(baseURL, topic, time.Duration(timeoutS)*time.Second)
	default:
		return nil, fmt.Errorf("sink: unknown type %q", typ)
	}
}

// stdoutSink writes JSON lines to standard output.
type stdoutSink struct{}

func NewStdoutSink() Sink { return &stdoutSink{} }

func (s *stdoutSink) Write(entry map[string]interface{}) error {
	return writeJSON(os.Stdout, entry)
}
func (s *stdoutSink) Close() error { return nil }

// fileSink appends JSON lines to a file.
type fileSink struct{ f *os.File }

func NewFileSink(path string) (Sink, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return nil, err
	}
	return &fileSink{f: f}, nil
}

func (s *fileSink) Write(entry map[string]interface{}) error {
	return writeJSON(s.f, entry)
}
func (s *fileSink) Close() error { return s.f.Close() }
