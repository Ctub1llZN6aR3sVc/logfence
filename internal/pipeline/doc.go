// Package pipeline provides the central processing loop for logfence.
//
// A Pipeline connects three subsystems:
//
//   - ingress.Reader  — parses newline-delimited JSON from an io.Reader
//   - router.Router   — matches each entry against configured routes and
//     forwards matching entries to the appropriate sinks
//   - metrics.Counters — tracks received, routed, and dropped entry counts
//
// Typical usage:
//
//	rd := ingress.NewReader(os.Stdin)
//	rt, _ := router.Build(cfg.Routes)
//	c  := metrics.New()
//	p  := pipeline.New(rd, rt, c, slog.Default())
//	if err := p.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
//	    log.Fatal(err)
//	}
//
// Run blocks until the context is cancelled or the underlying reader is
// exhausted, whichever comes first.
package pipeline
