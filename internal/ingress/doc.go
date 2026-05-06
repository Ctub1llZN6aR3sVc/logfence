// Package ingress provides primitives for reading structured log entries
// from arbitrary io.Reader sources (stdin, TCP connections, Unix sockets, etc.).
//
// Entries are expected to be newline-delimited JSON (NDJSON). Lines that
// cannot be parsed as JSON are still forwarded downstream with an empty
// field map so that routing rules based on raw content remain possible.
//
// Usage:
//
//	ch := ingress.Reader(ctx, os.Stdin)
//	for entry := range ch {
//		router.Dispatch(entry.Fields)
//	}
package ingress
