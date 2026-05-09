package sink_test

import (
	"encoding/json"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/organisation/logfence/internal/sink"
)

// startUDPSyslog opens a UDP listener that collects received datagrams.
func startUDPSyslog(t *testing.T) (addr string, recv func() []string) {
	t.Helper()
	conn, err := net.ListenPacket("udp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen udp: %v", err)
	}
	t.Cleanup(func() { conn.Close() })

	var msgs []string
	go func() {
		buf := make([]byte, 4096)
		for {
			n, _, err := conn.ReadFrom(buf)
			if err != nil {
				return
			}
			msgs = append(msgs, string(buf[:n]))
		}
	}()
	return conn.LocalAddr().String(), func() []string { return msgs }
}

func TestSyslogSink_Write(t *testing.T) {
	addr, recv := startUDPSyslog(t)

	s, err := sink.NewSyslogSink("udp", addr, "logfence-test")
	if err != nil {
		t.Fatalf("NewSyslogSink: %v", err)
	}
	defer s.Close()

	entry := map[string]any{"level": "info", "msg": "hello syslog"}
	if err := s.Write(entry); err != nil {
		t.Fatalf("Write: %v", err)
	}

	time.Sleep(50 * time.Millisecond)
	msgs := recv()
	if len(msgs) == 0 {
		t.Fatal("expected at least one syslog message, got none")
	}
	if !strings.Contains(msgs[0], "hello syslog") {
		t.Errorf("message %q does not contain expected payload", msgs[0])
	}
}

func TestSyslogSink_LevelRouting(t *testing.T) {
	addr, recv := startUDPSyslog(t)

	s, err := sink.NewSyslogSink("udp", addr, "logfence-test")
	if err != nil {
		t.Fatalf("NewSyslogSink: %v", err)
	}
	defer s.Close()

	levels := []string{"debug", "info", "warn", "error", "fatal", "unknown"}
	for _, lvl := range levels {
		entry := map[string]any{"level": lvl, "msg": "test " + lvl}
		if err := s.Write(entry); err != nil {
			t.Errorf("Write level=%s: %v", lvl, err)
		}
	}

	time.Sleep(80 * time.Millisecond)
	if got := len(recv()); got < len(levels) {
		t.Errorf("expected %d messages, got %d", len(levels), got)
	}
}

func TestSyslogSink_JSONPayload(t *testing.T) {
	addr, recv := startUDPSyslog(t)

	s, err := sink.NewSyslogSink("udp", addr, "logfence-test")
	if err != nil {
		t.Fatalf("NewSyslogSink: %v", err)
	}
	defer s.Close()

	entry := map[string]any{"level": "info", "service": "api", "code": 200}
	_ = s.Write(entry)
	time.Sleep(50 * time.Millisecond)

	msgs := recv()
	if len(msgs) == 0 {
		t.Fatal("no messages received")
	}
	// Extract JSON portion — syslog frames include a header prefix.
	idx := strings.Index(msgs[0], "{")
	if idx == -1 {
		t.Fatalf("no JSON found in %q", msgs[0])
	}
	var got map[string]any
	if err := json.Unmarshal([]byte(msgs[0][idx:]), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got["service"] != "api" {
		t.Errorf("service field: got %v, want api", got["service"])
	}
}

func TestNewSyslogSink_InvalidNetwork(t *testing.T) {
	_, err := sink.NewSyslogSink("unix", "/tmp/fake.sock", "tag")
	if err == nil {
		t.Fatal("expected error for unsupported network")
	}
}

func TestNewSyslogSink_InvalidAddr(t *testing.T) {
	_, err := sink.NewSyslogSink("udp", "not-an-addr", "tag")
	if err == nil {
		t.Fatal("expected error for invalid addr")
	}
}
