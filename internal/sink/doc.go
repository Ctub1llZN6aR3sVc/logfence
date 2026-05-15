// Package sink provides Sink implementations that write structured log entries
// (represented as map[string]any) to various backends.
//
// Available sinks:
//
//   - stdout      – writes JSON lines to os.Stdout
//   - file        – appends JSON lines to a file
//   - rotating    – appends JSON lines with automatic size-based rotation
//   - webhook     – HTTP POST JSON to an arbitrary URL
//   - kafka       – publishes via Kafka REST Proxy
//   - syslog      – forwards over UDP/TCP syslog
//   - redis       – RPUSH to a Redis list
//   - elasticsearch – indexes into Elasticsearch
//   - loki        – pushes to Grafana Loki
//   - splunk      – sends to a Splunk HTTP Event Collector (HEC)
//
// Use New(kind, opts) to construct any sink by name from a config map.
package sink
