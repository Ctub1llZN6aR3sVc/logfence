// Package pipeline wires ingress, router, and metrics into a single
// processing loop that can be started and stopped via context.
package pipeline

import (
	"context"
	"log/slog"

	"github.com/yourorg/logfence/internal/ingress"
	"github.com/yourorg/logfence/internal/metrics"
	"github.com/yourorg/logfence/internal/router"
)

// Pipeline reads log entries from an ingress.Reader, dispatches each entry
// through a router.Router, and records counters via a metrics.Counters.
type Pipeline struct {
	reader  *ingress.Reader
	router  *router.Router
	counters *metrics.Counters
	log     *slog.Logger
}

// New creates a Pipeline. All arguments must be non-nil.
func New(r *ingress.Reader, rt *router.Router, c *metrics.Counters, log *slog.Logger) *Pipeline {
	if log == nil {
		log = slog.Default()
	}
	return &Pipeline{
		reader:   r,
		router:   rt,
		counters: c,
		log:      log,
	}
}

// Run processes log entries until the context is cancelled or the reader is
// exhausted. It returns only context errors or unexpected reader failures.
func (p *Pipeline) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		entry, err := p.reader.Next(ctx)
		if err != nil {
			// EOF / context cancellation are normal shutdown paths.
			if ctx.Err() != nil {
				return ctx.Err()
			}
			return err
		}
		if entry == nil {
			// Reader exhausted (e.g. pipe closed).
			return nil
		}

		p.counters.IncReceived()

		if p.router.Dispatch(entry) {
			p.counters.IncRouted()
		} else {
			p.counters.IncDropped()
			p.log.Debug("entry dropped: no matching route", "entry", entry)
		}
	}
}
