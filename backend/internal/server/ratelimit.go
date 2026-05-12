package server

import (
	"sync"
	"time"
)

type RateLimiter struct {
	mu        sync.Mutex
	window    time.Duration
	maxCalls  int
	lastCalls map[string][]time.Time
}

// NewRateLimiter creates a new rate limiter with the specified window and max calls per window.
func NewRateLimiter(window time.Duration, maxCalls int) *RateLimiter {
	return &RateLimiter{
		window:    window,
		maxCalls:  maxCalls,
		lastCalls: make(map[string][]time.Time),
	}
}

// Allow checks if a request from the given key is allowed based on rate limits.
// It returns true if the request is allowed, false if it exceeds the limit.
func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-rl.window)

	// Clean up old timestamps for this key
	calls := rl.lastCalls[key]
	var recentCalls []time.Time
	for _, t := range calls {
		if t.After(cutoff) {
			recentCalls = append(recentCalls, t)
		}
	}

	// Check if under limit
	if len(recentCalls) < rl.maxCalls {
		recentCalls = append(recentCalls, now)
		rl.lastCalls[key] = recentCalls
		return true
	}

	return false
}
