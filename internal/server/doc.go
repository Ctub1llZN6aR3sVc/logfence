// Package server implements the logfence HTTP management server.
//
// It exposes three endpoints:
//
//	 GET /healthz  – liveness probe; returns {"status":"ok"} when the
//	                 process is running.
//
//	 GET /readyz   – readiness probe; returns {"status":"ready"} once
//	                 the daemon has finished initialising.
//
//	 GET /metrics  – Prometheus-compatible plain-text exposition of the
//	                 internal counters provided by internal/metrics.
//
// Typical usage:
//
//	srv := server.New(cfg.ManagementAddr, counters)
//	go func() {
//	    if err := srv.Start(); err != nil && err != http.ErrServerClosed {
//	        log.Fatal(err)
//	    }
//	}()
//	// … on shutdown …
//	srv.Shutdown(ctx)
package server
