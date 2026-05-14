// Package sink provides output destinations for structured log entries.
//
// Each sink implements the Sink interface:
//
//	type Sink interface {
//		Write(entry map[string]any) error
//		Close() error
//	}
//
// Available sink types
//
//   - stdout   – writes JSON lines to standard output
//   - file     – appends JSON lines to a file
//   - rotating – like file but rotates when a maximum byte size is reached
//   - webhook  – HTTP POST JSON payload to a remote endpoint
//   - kafka    – publishes via Kafka REST Proxy
//   - syslog   – forwards entries to a syslog daemon over UDP or TCP
//   - redis    – pushes JSON entries onto a Redis list using RPUSH
//
// The New factory function selects the appropriate implementation from a
// sink-type string and a string→string configuration map, making it easy
// to wire sinks directly from the parsed YAML configuration.
package sink
