// Package sink provides writers that forward structured log entries to various
// destinations such as stdout, files, webhooks, Kafka, syslog, Redis,
// Elasticsearch, Loki, and Splunk.
package sink

import (
	"fmt"
)

// Sink is the interface implemented by every log destination.
type Sink interface {
	Write(entry map[string]any) error
	Close() error
}

// New constructs a Sink from a type name and a string-keyed options map.
// Recognised types: stdout, file, rotating, webhook, kafka, syslog, redis,
// elasticsearch, loki, splunk.
func New(kind string, opts map[string]string) (Sink, error) {
	switch kind {
	case "stdout":
		return NewStdoutSink(), nil
	case "file":
		path, ok := opts["path"]
		if !ok {
			return nil, fmt.Errorf("sink file: missing 'path' option")
		}
		return NewFileSink(path)
	case "rotating":
		path, ok := opts["path"]
		if !ok {
			return nil, fmt.Errorf("sink rotating: missing 'path' option")
		}
		return NewRotatingFileSink(path, 0, 0)
	case "webhook":
		return NewWebhookSink(opts["url"])
	case "kafka":
		return NewKafkaRestSink(opts["url"], opts["topic"])
	case "syslog":
		return NewSyslogSink(opts["network"], opts["addr"], opts["tag"])
	case "redis":
		return NewRedisSink(opts["addr"], opts["list"])
	case "elasticsearch":
		return NewElasticsearchSink(opts["url"])
	case "loki":
		return NewLokiSink(opts["url"])
	case "splunk":
		return NewSplunkSink(opts["url"], opts["token"])
	default:
		return nil, fmt.Errorf("sink: unknown type %q", kind)
	}
}
