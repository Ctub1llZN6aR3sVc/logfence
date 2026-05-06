// Package ingress provides primitives for reading structured log entries
// from arbitrary io.Reader sources (stdin, TCP connections, Unix sockets, etc.).
//
// Entries are expected to be newline-delimited JSON (NDJSON). Lines that
// cannot be parsed as JSON are still forwarded downstream with an empty
// field map so that routing rules based on raw content remain possible.
//
// The Reader function returns a read-only channel that is closed automatically
// when the provided context is cancelled or the underlying reader reaches EOF.
// Callers should range over the channel to process entries:
//
//	ch := ingress.Reader(ctx, os.Stdin)
//	for entry := range ch {
//		router.Dispatch(entry.Fields)
//	}
//
// For network sources, prefer wrapping accepted connections with a
// bufio.Reader before passing them to Reader to amortise syscall overhead.
package ingress
