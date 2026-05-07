// Package ingress handles log entry ingestion from various sources.
package ingress

import (
	"strings"
	"time"
)

// MultilineConfig defines how multiline log entries are assembled.
type MultilineConfig struct {
	// StartPattern is a prefix that indicates the start of a new log entry.
	// Lines not matching this pattern are appended to the previous entry.
	StartPattern string
	// MaxAge is the maximum time to wait before flushing an incomplete entry.
	MaxAge time.Duration
	// MaxLines is the maximum number of lines to accumulate before forcing a flush.
	MaxLines int
}

// MultilineAssembler accumulates lines into complete log entries.
type MultilineAssembler struct {
	cfg     MultilineConfig
	buf     []string
	started time.Time
}

// NewMultilineAssembler creates an assembler with the given configuration.
func NewMultilineAssembler(cfg MultilineConfig) *MultilineAssembler {
	if cfg.MaxAge == 0 {
		cfg.MaxAge = 5 * time.Second
	}
	if cfg.MaxLines == 0 {
		cfg.MaxLines = 100
	}
	return &MultilineAssembler{cfg: cfg}
}

// Feed adds a raw line to the assembler. It returns a completed entry and true
// when the buffer should be flushed, or empty string and false otherwise.
func (a *MultilineAssembler) Feed(line string) (string, bool) {
	isStart := a.cfg.StartPattern == "" || strings.HasPrefix(line, a.cfg.StartPattern)

	if isStart && len(a.buf) > 0 {
		// Flush the previous entry before starting a new one.
		entry := strings.Join(a.buf, "\n")
		a.buf = []string{line}
		a.started = time.Now()
		return entry, true
	}

	if len(a.buf) == 0 {
		a.started = time.Now()
	}
	a.buf = append(a.buf, line)

	// Force flush if limits exceeded.
	if len(a.buf) >= a.cfg.MaxLines || time.Since(a.started) >= a.cfg.MaxAge {
		return a.Flush()
	}
	return "", false
}

// Flush returns the current buffer contents as a single entry and resets state.
func (a *MultilineAssembler) Flush() (string, bool) {
	if len(a.buf) == 0 {
		return "", false
	}
	entry := strings.Join(a.buf, "\n")
	a.buf = nil
	a.started = time.Time{}
	return entry, true
}
