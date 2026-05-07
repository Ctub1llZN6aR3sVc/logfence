// export_test.go exposes internal fields for white-box testing.
package server

import "net/http"

// Handler returns the underlying http.Handler so tests can wrap it
// in httptest.NewServer without binding a real port.
func (s *Server) Handler() http.Handler {
	return s.httpServer.Handler
}
