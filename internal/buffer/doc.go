// Package buffer implements a bounded ring buffer for log entries.
//
// The ring buffer is designed for use between the ingress reader and the
// routing pipeline. When the downstream pipeline is slower than the ingress
// source the buffer absorbs bursts up to its configured capacity. Once full,
// the oldest entry is silently dropped and a Dropped counter is incremented
// so operators can observe back-pressure via the /metrics endpoint.
//
// Usage:
//
//	buf := buffer.New(1024)
//
//	// producer
//	buf.Push(entry)
//
//	// consumer
//	if e, ok := buf.Pop(); ok {
//		// process e
//	}
package buffer
