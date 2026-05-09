package sink

import (
	"encoding/json"
	"fmt"
	"log/syslog"
	"net"
)

// SyslogSink forwards log entries to a syslog daemon over UDP or TCP.
type SyslogSink struct {
	writer *syslog.Writer
}

// NewSyslogSink dials the given network/address and returns a SyslogSink.
// network must be "udp" or "tcp"; addr is host:port.
func NewSyslogSink(network, addr, tag string) (*SyslogSink, error) {
	if network != "udp" && network != "tcp" {
		return nil, fmt.Errorf("syslog: unsupported network %q (want udp or tcp)", network)
	}
	if _, _, err := net.SplitHostPort(addr); err != nil {
		return nil, fmt.Errorf("syslog: invalid addr %q: %w", addr, err)
	}
	w, err := syslog.Dial(network, addr, syslog.LOG_INFO|syslog.LOG_DAEMON, tag)
	if err != nil {
		return nil, fmt.Errorf("syslog: dial %s://%s: %w", network, addr, err)
	}
	return &SyslogSink{writer: w}, nil
}

// Write serialises entry as JSON and emits it at the appropriate syslog priority.
func (s *SyslogSink) Write(entry map[string]any) error {
	line, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("syslog: marshal: %w", err)
	}
	msg := string(line)

	level, _ := entry["level"].(string)
	switch level {
	case "error", "fatal":
		return s.writer.Err(msg)
	case "warn", "warning":
		return s.writer.Warning(msg)
	case "debug":
		return s.writer.Debug(msg)
	default:
		return s.writer.Info(msg)
	}
}

// Close closes the underlying syslog connection.
func (s *SyslogSink) Close() error { return s.writer.Close() }
