package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/logfence/internal/config"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "logfence-*.yaml")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestLoad_Valid(t *testing.T) {
	raw := `
listen_addr: ":6000"
sinks:
  - name: console
    type: stdout
routes:
  - name: all-errors
    levels: [error, fatal]
    sink: console
`
	path := writeTemp(t, raw)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.ListenAddr != ":6000" {
		t.Errorf("listen_addr = %q, want :6000", cfg.ListenAddr)
	}
	if len(cfg.Sinks) != 1 || cfg.Sinks[0].Name != "console" {
		t.Errorf("unexpected sinks: %+v", cfg.Sinks)
	}
}

func TestLoad_DefaultListenAddr(t *testing.T) {
	raw := `
sinks:
  - name: out
    type: stdout
routes:
  - name: r1
    sink: out
`
	cfg, err := config.Load(writeTemp(t, raw))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.ListenAddr != ":5170" {
		t.Errorf("default listen_addr = %q, want :5170", cfg.ListenAddr)
	}
}

func TestLoad_UnknownSink(t *testing.T) {
	raw := `
sinks:
  - name: real
    type: stdout
routes:
  - name: bad-route
    sink: ghost
`
	_, err := config.Load(writeTemp(t, raw))
	if err == nil {
		t.Fatal("expected error for unknown sink reference, got nil")
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := config.Load(filepath.Join(t.TempDir(), "missing.yaml"))
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
