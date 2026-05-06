package ingress

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"log"
)

// Entry represents a single structured log entry read from an input source.
type Entry struct {
	Raw    []byte
	Fields map[string]interface{}
}

// Reader reads newline-delimited JSON log entries from an io.Reader and
// emits parsed Entry values on the returned channel. The channel is closed
// when the reader is exhausted or ctx is cancelled.
func Reader(ctx context.Context, r io.Reader) <-chan Entry {
	ch := make(chan Entry, 64)
	go func() {
		defer close(ch)
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			select {
			case <-ctx.Done():
				return
			default:
			}
			line := scanner.Bytes()
			if len(line) == 0 {
				continue
			}
			e := Entry{Raw: append([]byte(nil), line...)}
			if err := json.Unmarshal(line, &e.Fields); err != nil {
				// Non-JSON lines are forwarded with an empty field map so
				// downstream filters can still route them by raw content.
				e.Fields = map[string]interface{}{}
				log.Printf("ingress: non-JSON line skipped parsing: %v", err)
			}
			ch <- e
		}
		if err := scanner.Err(); err != nil {
			log.Printf("ingress: scanner error: %v", err)
		}
	}()
	return ch
}
