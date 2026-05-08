package sink

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// RotatingFileSink writes log entries to a file and rotates it when it
// exceeds a configured size limit. Rotated files are renamed with a
// timestamp suffix.
type RotatingFileSink struct {
	mu      sync.Mutex
	path    string
	maxSize int64 // bytes
	file    *os.File
	size    int64
}

// NewRotatingFileSink creates a RotatingFileSink that rotates the file at
// path once it grows beyond maxBytes.
func NewRotatingFileSink(path string, maxBytes int64) (*RotatingFileSink, error) {
	if maxBytes <= 0 {
		return nil, fmt.Errorf("rotating sink: maxBytes must be positive, got %d", maxBytes)
	}
	s := &RotatingFileSink{path: path, maxSize: maxBytes}
	if err := s.openOrCreate(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *RotatingFileSink) openOrCreate() error {
	f, err := os.OpenFile(s.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("rotating sink: open %s: %w", s.path, err)
	}
	info, err := f.Stat()
	if err != nil {
		_ = f.Close()
		return fmt.Errorf("rotating sink: stat %s: %w", s.path, err)
	}
	s.file = f
	s.size = info.Size()
	return nil
}

// Write appends the JSON entry (with a trailing newline) to the current log
// file, rotating first if the size limit would be exceeded.
func (s *RotatingFileSink) Write(entry map[string]any) error {
	line, err := marshalEntry(entry)
	if err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.size+int64(len(line)) > s.maxSize {
		if err := s.rotate(); err != nil {
			return err
		}
	}
	n, err := s.file.Write(line)
	if err != nil {
		return fmt.Errorf("rotating sink: write: %w", err)
	}
	s.size += int64(n)
	return nil
}

func (s *RotatingFileSink) rotate() error {
	_ = s.file.Close()
	ts := time.Now().UTC().Format("20060102T150405Z")
	dest := fmt.Sprintf("%s.%s", s.path, ts)
	if err := os.Rename(s.path, dest); err != nil {
		return fmt.Errorf("rotating sink: rename: %w", err)
	}
	return s.openOrCreate()
}

// Close flushes and closes the underlying file.
func (s *RotatingFileSink) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.file == nil {
		return nil
	}
	err := s.file.Close()
	s.file = nil
	return err
}

// Dir returns the directory that will contain rotated files (same as the
// directory of the active log file).
func (s *RotatingFileSink) Dir() string { return filepath.Dir(s.path) }
