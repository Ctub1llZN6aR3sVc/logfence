// Package server provides the HTTP management server for logfence,
// exposing health, readiness, and metrics endpoints.
package server

import (
	"context"
	"net/http"
	"time"

	"github.com/yourorg/logfence/internal/metrics"
)

const (
	defaultReadTimeout  = 5 * time.Second
	defaultWriteTimeout = 10 * time.Second
	defaultIdleTimeout  = 30 * time.Second
)

// Server is the HTTP management server.
type Server struct {
	httpServer *http.Server
	metrics    *metrics.Counters
}

// New creates a new Server bound to addr, wiring up all management routes.
func New(addr string, m *metrics.Counters) *Server {
	mux := http.NewServeMux()
	s := &Server{metrics: m}

	mux.HandleFunc("/healthz", s.handleHealth)
	mux.HandleFunc("/readyz", s.handleReady)
	mux.HandleFunc("/metrics", s.handleMetrics)

	s.httpServer = &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
		IdleTimeout:  defaultIdleTimeout,
	}
	return s
}

// Start begins listening and serving in the foreground.
// It returns http.ErrServerClosed on graceful shutdown.
func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully stops the server using the provided context.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
