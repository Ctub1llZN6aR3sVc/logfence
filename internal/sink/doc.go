// Package sink provides output destination implementations for logfence.
//
// Available sink types:
//
//	"stdout"         – writes JSON log entries to standard output (default).
//	"file"           – appends JSON log entries to a file at a fixed path.
//	"rotating_file"  – like "file" but rotates when the file exceeds MaxBytes.
//	"webhook"        – POSTs each entry as a JSON body to an HTTP endpoint.
//
// All sinks satisfy the Sink interface:
//
//	type Sink interface {
//	    Write(entry map[string]any) error
//	    Close() error
//	}
//
// Use New(Config) to construct the appropriate Sink from a Config struct
// that is typically populated from the YAML configuration file.
package sink
