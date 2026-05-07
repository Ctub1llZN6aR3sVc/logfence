// Package ratelimit provides per-route token-bucket rate limiting for log entries.
package ratelimit

import (
	"sync"
	"time"
)

// Limiter is a token-bucket rate limiter keyed by route name.
type Limiter struct {
	mu      sync.Mutex
	buckets map[string]*bucket
	rate    int // tokens per second
	burst   int
}

type bucket struct {
	tokens   float64
	lastSeen time.Time
}

// New creates a Limiter with the given sustained rate (per second) and burst size.
func New(rate, burst int) *Limiter {
	return &Limiter{
		buckets: make(map[string]*bucket),
		rate:    rate,
		burst:   burst,
	}
}

// Allow returns true if the named route is within its rate limit.
func (l *Limiter) Allow(route string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	b, ok := l.buckets[route]
	if !ok {
		b = &bucket{tokens: float64(l.burst), lastSeen: now}
		l.buckets[route] = b
	}

	elapsed := now.Sub(b.lastSeen).Seconds()
	b.lastSeen = now
	b.tokens += elapsed * float64(l.rate)
	if b.tokens > float64(l.burst) {
		b.tokens = float64(l.burst)
	}

	if b.tokens >= 1 {
		b.tokens--
		return true
	}
	return false
}

// Reset clears the bucket for the given route.
func (l *Limiter) Reset(route string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.buckets, route)
}
