package buffer

import (
	"sync"
	"testing"

	"github.com/your-org/logfence/internal/ingress"
)

func entry(msg string) ingress.Entry {
	return ingress.Entry{Fields: map[string]any{"msg": msg}}
}

func TestRing_PushPop(t *testing.T) {
	r := New(3)
	r.Push(entry("a"))
	r.Push(entry("b"))

	e, ok := r.Pop()
	if !ok {
		t.Fatal("expected entry")
	}
	if e.Fields["msg"] != "a" {
		t.Fatalf("want a, got %v", e.Fields["msg"])
	}
}

func TestRing_EmptyPop(t *testing.T) {
	r := New(2)
	_, ok := r.Pop()
	if ok {
		t.Fatal("expected empty")
	}
}

func TestRing_Eviction(t *testing.T) {
	r := New(2)
	r.Push(entry("a"))
	r.Push(entry("b"))
	r.Push(entry("c")) // evicts "a"

	if r.Dropped != 1 {
		t.Fatalf("want Dropped=1, got %d", r.Dropped)
	}
	if r.Len() != 2 {
		t.Fatalf("want Len=2, got %d", r.Len())
	}

	e, _ := r.Pop()
	if e.Fields["msg"] != "b" {
		t.Fatalf("want b, got %v", e.Fields["msg"])
	}
}

func TestRing_Len(t *testing.T) {
	r := New(5)
	for i := 0; i < 3; i++ {
		r.Push(entry("x"))
	}
	if r.Len() != 3 {
		t.Fatalf("want 3, got %d", r.Len())
	}
	r.Pop()
	if r.Len() != 2 {
		t.Fatalf("want 2, got %d", r.Len())
	}
}

func TestRing_Concurrent(t *testing.T) {
	r := New(64)
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			r.Push(entry("concurrent"))
			r.Pop()
		}()
	}
	wg.Wait()
}

func TestNew_MinCapacity(t *testing.T) {
	r := New(0)
	if r.cap != 1 {
		t.Fatalf("expected cap=1 for zero input, got %d", r.cap)
	}
}
