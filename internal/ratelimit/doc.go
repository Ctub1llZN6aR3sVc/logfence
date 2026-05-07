// Package ratelimit implements per-route token-bucket rate limiting.
//
// Each named route maintains an independent bucket. Tokens refill at the
// configured sustained rate (tokens/second) up to the burst ceiling.
//
// Usage:
//
//	limiter := ratelimit.New(100, 200) // 100 t/s, burst 200
//	if limiter.Allow(routeName) {
//	    // forward the log entry
//	} else {
//	    // drop or sample
//	}
package ratelimit
