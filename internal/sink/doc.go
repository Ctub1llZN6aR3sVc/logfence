// Package sink provides output destinations for log entries processed by
// logfence. Each sink implements a common Write/Close interface so the router
// can dispatch entries without knowing the underlying transport.
//
// Available sinks:
//
//   - StdoutSink  – writes JSON lines to standard output.
//   - FileSink    – appends JSON lines to a plain file.
//   - RotatingFileSink – like FileSink but rotates when the file exceeds a
//     configurable byte threshold, keeping a bounded number of backups.
//   - WebhookSink – POSTs JSON payloads to an HTTP endpoint.
//   - KafkaRestSink – publishes entries to a Kafka topic via the Confluent
//     REST Proxy HTTP API.
//   - SyslogSink  – forwards entries to a syslog daemon over UDP or TCP,
//     mapping the entry "level" field to the appropriate syslog priority.
//
// New sinks are registered in New() (sink.go) so that the router builder can
// instantiate them from configuration by name.
package sink
