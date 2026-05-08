// Package buffer provides a bounded, thread-safe ring buffer for log entries.
// When the buffer is full, the oldest entry is dropped to make room for the newest.
package buffer

import (
	"sync"

	"github.com/your-org/logfence/internal/ingress"
)

// Dropped tracks the number of entries dropped due to a full buffer.
type Ring struct {
	mu      sync.Mutex
	items   []ingress.Entry
	head    int
	tail    int
	count   int
	cap     int
	Dropped uint64
}

// New creates a new Ring buffer with the given capacity.
// capacity must be >= 1.
func New(capacity int) *Ring {
	if capacity < 1 {
		capacity = 1
	}
	return &Ring{
		items: make([]ingress.Entry, capacity),
		cap:   capacity,
	}
}

// Push adds an entry to the buffer. If the buffer is full, the oldest entry
// is evicted and Dropped is incremented.
func (r *Ring) Push(e ingress.Entry) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.count == r.cap {
		// Evict oldest
		r.head = (r.head + 1) % r.cap
		r.count--
		r.Dropped++
	}

	r.items[r.tail] = e
	r.tail = (r.tail + 1) % r.cap
	r.count++
}

// Pop removes and returns the oldest entry from the buffer.
// The second return value is false when the buffer is empty.
func (r *Ring) Pop() (ingress.Entry, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.count == 0 {
		return ingress.Entry{}, false
	}

	e := r.items[r.head]
	r.head = (r.head + 1) % r.cap
	r.count--
	return e, true
}

// Len returns the current number of entries in the buffer.
func (r *Ring) Len() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.count
}
