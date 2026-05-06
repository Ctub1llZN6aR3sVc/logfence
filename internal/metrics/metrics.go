// Package metrics provides lightweight counters for logfence runtime telemetry.
package metrics

import (
	"fmt"
	"io"
	"sync/atomic"
)

// Counters holds atomic counters for key logfence events.
type Counters struct {
	LinesReceived  atomic.Int64
	LinesDropped   atomic.Int64
	LinesRouted    atomic.Int64
	DecodeErrors   atomic.Int64
	SinkErrors     atomic.Int64
}

// New returns an initialised Counters instance.
func New() *Counters {
	return &Counters{}
}

// IncReceived increments the lines-received counter.
func (c *Counters) IncReceived() { c.LinesReceived.Add(1) }

// IncDropped increments the lines-dropped counter.
func (c *Counters) IncDropped() { c.LinesDropped.Add(1) }

// IncRouted increments the lines-routed counter.
func (c *Counters) IncRouted() { c.LinesRouted.Add(1) }

// IncDecodeError increments the decode-error counter.
func (c *Counters) IncDecodeError() { c.DecodeErrors.Add(1) }

// IncSinkError increments the sink-error counter.
func (c *Counters) IncSinkError() { c.SinkErrors.Add(1) }

// Snapshot returns a point-in-time copy of all counter values.
func (c *Counters) Snapshot() map[string]int64 {
	return map[string]int64{
		"lines_received": c.LinesReceived.Load(),
		"lines_dropped":  c.LinesDropped.Load(),
		"lines_routed":   c.LinesRouted.Load(),
		"decode_errors":  c.DecodeErrors.Load(),
		"sink_errors":    c.SinkErrors.Load(),
	}
}

// WriteTo writes a human-readable summary of the counters to w.
func (c *Counters) WriteTo(w io.Writer) (int64, error) {
	s := c.Snapshot()
	n, err := fmt.Fprintf(w,
		"lines_received=%d lines_routed=%d lines_dropped=%d decode_errors=%d sink_errors=%d\n",
		s["lines_received"], s["lines_routed"], s["lines_dropped"],
		s["decode_errors"], s["sink_errors"],
	)
	return int64(n), err
}
