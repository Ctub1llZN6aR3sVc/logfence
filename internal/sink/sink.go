package sink

import (
	"fmt"
	"os"
)

// Sink is the interface all output destinations must implement.
type Sink interface {
	Write(entry map[string]any) error
	Close() error
}

// Config describes a sink in the YAML configuration.
type Config struct {
	Type       string `yaml:"type"`
	Path       string `yaml:"path"`
	MaxBytes   int64  `yaml:"max_bytes"`
	MaxBackups int    `yaml:"max_backups"`
	WebhookURL string `yaml:"webhook_url"`
	TimeoutSec int    `yaml:"timeout_sec"`
}

// New constructs a Sink from a Config.
func New(cfg Config) (Sink, error) {
	switch cfg.Type {
	case "stdout", "":
		return NewStdoutSink(), nil
	case "file":
		return NewFileSink(cfg.Path)
	case "rotating_file":
		return NewRotatingFileSink(RotatingConfig{
			Path:       cfg.Path,
			MaxBytes:   cfg.MaxBytes,
			MaxBackups: cfg.MaxBackups,
		})
	case "webhook":
		return NewWebhookSink(WebhookConfig{
			URL:            cfg.WebhookURL,
			TimeoutSeconds: cfg.TimeoutSec,
		})
	default:
		return nil, fmt.Errorf("unknown sink type: %q", cfg.Type)
	}
}

// StdoutSink writes JSON entries to stdout.
type StdoutSink struct{}

func NewStdoutSink() *StdoutSink { return &StdoutSink{} }

func (s *StdoutSink) Write(entry map[string]any) error {
	return writeJSON(os.Stdout, entry)
}

func (s *StdoutSink) Close() error { return nil }

// FileSink writes JSON entries to a file.
type FileSink struct {
	f *os.File
}

func NewFileSink(path string) (*FileSink, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return nil, fmt.Errorf("open file sink %q: %w", path, err)
	}
	return &FileSink{f: f}, nil
}

func (s *FileSink) Write(entry map[string]any) error { return writeJSON(s.f, entry) }
func (s *FileSink) Close() error                    { return s.f.Close() }
