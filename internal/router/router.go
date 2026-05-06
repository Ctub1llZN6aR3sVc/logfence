package router

import (
	"fmt"

	"github.com/yourorg/logfence/internal/filter"
	"github.com/yourorg/logfence/internal/sink"
)

// Route pairs a filter chain with a destination sink.
type Route struct {
	Name   string
	Filter filter.Chain
	Sink   sink.Sink
}

// Router dispatches log entries to matching routes.
type Router struct {
	routes []Route
}

// New creates a Router from the provided routes.
func New(routes []Route) *Router {
	return &Router{routes: routes}
}

// Dispatch sends entry to every route whose filter chain accepts it.
// entry is expected to be a structured map with at least a "level" key.
func (r *Router) Dispatch(entry map[string]interface{}) error {
	for _, route := range r.routes {
		if route.Filter.Match(entry) {
			if err := route.Sink.Write(entry); err != nil {
				return fmt.Errorf("router: route %q write error: %w", route.Name, err)
			}
		}
	}
	return nil
}

// Close closes all sinks registered in the router.
func (r *Router) Close() error {
	var firstErr error
	for _, route := range r.routes {
		if err := route.Sink.Close(); err != nil && firstErr == nil {
			firstErr = fmt.Errorf("router: route %q close error: %w", route.Name, err)
		}
	}
	return firstErr
}
