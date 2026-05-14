package sink_test

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/your-org/logfence/internal/sink"
)

// startFakeRedis listens on a random TCP port, reads one RPUSH command and
// replies with ":1\r\n". It returns the listener address and a channel that
// receives the raw command string when it arrives.
func startFakeRedis(t *testing.T, reply string) (addr string, cmdCh <-chan string) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	ch := make(chan string, 1)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		scanner := bufio.NewScanner(conn)
		var lines []string
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
			if len(lines) >= 6 { // *3 + 3 pairs of len+value
				break
			}
		}
		ch <- strings.Join(lines, "\n")
		fmt.Fprint(conn, reply)
	}()
	t.Cleanup(func() { ln.Close(); wg.Wait() })
	return ln.Addr().String(), ch
}

func TestRedisSink_Write(t *testing.T) {
	addr, cmdCh := startFakeRedis(t, ":1\r\n")
	s, err := sink.NewRedisSink(addr, "logs", time.Second)
	if err != nil {
		t.Fatalf("NewRedisSink: %v", err)
	}
	entry := map[string]any{"level": "info", "msg": "hello"}
	if err := s.Write(entry); err != nil {
		t.Fatalf("Write: %v", err)
	}
	cmd := <-cmdCh
	if !strings.Contains(cmd, "RPUSH") {
		t.Errorf("expected RPUSH in command, got: %s", cmd)
	}
	if !strings.Contains(cmd, "logs") {
		t.Errorf("expected list name 'logs' in command, got: %s", cmd)
	}
}

func TestRedisSink_ServerError(t *testing.T) {
	addr, _ := startFakeRedis(t, "-ERR unknown command\r\n")
	s, _ := sink.NewRedisSink(addr, "logs", time.Second)
	if err := s.Write(map[string]any{"x": 1}); err == nil {
		t.Fatal("expected error on Redis error reply")
	}
}

func TestRedisSink_EmptyAddr(t *testing.T) {
	_, err := sink.NewRedisSink("", "logs", time.Second)
	if err == nil {
		t.Fatal("expected error for empty addr")
	}
}

func TestRedisSink_EmptyList(t *testing.T) {
	_, err := sink.NewRedisSink("127.0.0.1:6379", "", time.Second)
	if err == nil {
		t.Fatal("expected error for empty list")
	}
}

func TestRedisSink_Close(t *testing.T) {
	s, _ := sink.NewRedisSink("127.0.0.1:6379", "logs", time.Second)
	if err := s.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
}

func TestNew_Redis(t *testing.T) {
	addr, _ := startFakeRedis(t, ":1\r\n")
	s, err := sink.New("redis", map[string]string{
		"addr":    addr,
		"list":    "applogs",
		"timeout": "2s",
	})
	if err != nil {
		t.Fatalf("New redis: %v", err)
	}
	if s == nil {
		t.Fatal("expected non-nil sink")
	}
}
