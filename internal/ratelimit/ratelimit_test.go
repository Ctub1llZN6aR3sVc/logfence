package ratelimit_test

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/yourorg/logfence/internal/ratelimit"
)

func TestAllow_BurstConsumed(t *testing.T) {
	l := ratelimit.New(1, 3)
	for i := 0; i < 3; i++ {
		if !l.Allow("route-a") {
			t.Fatalf("expected allow on call %d", i)
		}
	}
	if l.Allow("route-a") {
		t.Fatal("expected deny after burst exhausted")
	}
}

func TestAllow_IndependentRoutes(t *testing.T) {
	l := ratelimit.New(1, 2)
	l.Allow("a")
	l.Allow("a")
	// route "b" should still have its own full burst
	if !l.Allow("b") {
		t.Fatal("expected allow for independent route")
	}
}

func TestAllow_Refill(t *testing.T) {
	l := ratelimit.New(100, 1)
	if !l.Allow("r") {
		t.Fatal("first allow should succeed")
	}
	if l.Allow("r") {
		t.Fatal("second immediate allow should fail")
	}
	time.Sleep(15 * time.Millisecond)
	if !l.Allow("r") {
		t.Fatal("allow after refill should succeed")
	}
}

func TestAllow_Reset(t *testing.T) {
	l := ratelimit.New(1, 1)
	l.Allow("x")
	l.Reset("x")
	if !l.Allow("x") {
		t.Fatal("expected allow after reset")
	}
}

func TestAllow_Concurrent(t *testing.T) {
	l := ratelimit.New(1000, 500)
	var allowed atomic.Int64
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if l.Allow("concurrent") {
				allowed.Add(1)
			}
		}()
	}
	wg.Wait()
	if allowed.Load() == 0 {
		t.Fatal("expected at least some allowed")
	}
}
