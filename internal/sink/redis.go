package sink

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"time"
)

// RedisSink publishes log entries to a Redis list via the RPUSH command
// using the inline command protocol (no external dependency required).
type RedisSink struct {
	addr    string
	list    string
	timeout time.Duration
}

// NewRedisSink creates a RedisSink that pushes JSON entries to the given
// Redis list. addr is "host:port", list is the Redis key name.
func NewRedisSink(addr, list string, timeout time.Duration) (*RedisSink, error) {
	if addr == "" {
		return nil, fmt.Errorf("redis: addr must not be empty")
	}
	if list == "" {
		return nil, fmt.Errorf("redis: list must not be empty")
	}
	if timeout <= 0 {
		timeout = 3 * time.Second
	}
	return &RedisSink{addr: addr, list: list, timeout: timeout}, nil
}

// Write serialises entry as JSON and pushes it onto the Redis list.
func (r *RedisSink) Write(entry map[string]any) error {
	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("redis: marshal: %w", err)
	}

	conn, err := net.DialTimeout("tcp", r.addr, r.timeout)
	if err != nil {
		return fmt.Errorf("redis: dial %s: %w", r.addr, err)
	}
	defer conn.Close()

	_ = conn.SetDeadline(time.Now().Add(r.timeout))

	// Build RESP inline RPUSH command.
	cmd := fmt.Sprintf("*3\r\n$5\r\nRPUSH\r\n$%d\r\n%s\r\n$%d\r\n%s\r\n",
		len(r.list), r.list, len(data), data)

	if _, err = fmt.Fprint(conn, cmd); err != nil {
		return fmt.Errorf("redis: send: %w", err)
	}

	// Read at least the first byte of the reply to confirm delivery.
	buf := make([]byte, 1)
	if _, err = conn.Read(buf); err != nil {
		return fmt.Errorf("redis: read reply: %w", err)
	}
	if buf[0] == '-' {
		return fmt.Errorf("redis: server returned error")
	}
	return nil
}

// Close is a no-op; connections are short-lived per write.
func (r *RedisSink) Close() error { return nil }

// newRedisSink wires RedisSink into the generic New factory.
func newRedisSink(cfg map[string]string) (Sink, error) {
	timeout, _ := time.ParseDuration(cfg["timeout"])
	return NewRedisSink(cfg["addr"], cfg["list"], timeout)
}

// ensure compile-time interface satisfaction.
var _ Sink = (*RedisSink)(nil)

// contextKey is used internally to carry a context through dial when needed.
type contextKey struct{}

// dialContext wraps net.DialTimeout with context cancellation awareness.
func dialContext(ctx context.Context, addr string, timeout time.Duration) (net.Conn, error) {
	d := net.Dialer{Timeout: timeout}
	return d.DialContext(ctx, "tcp", addr)
}
