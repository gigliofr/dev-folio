package server

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"devfolio/backend/internal/domain"
	"devfolio/backend/internal/store"
)

func TestRateLimiter(t *testing.T) {
	rl := NewRateLimiter(1*time.Second, 2)

	tests := []struct {
		name     string
		key      string
		call     int
		expected bool
	}{
		{"first call", "user@example.com", 1, true},
		{"second call", "user@example.com", 2, true},
		{"third call (over limit)", "user@example.com", 3, false},
		{"different user", "other@example.com", 1, true},
	}

	for _, tc := range tests {
		result := rl.Allow(tc.key)
		if result != tc.expected {
			t.Errorf("%s: expected %v, got %v", tc.name, tc.expected, result)
		}
	}

	// Wait for window to expire
	time.Sleep(1100 * time.Millisecond)

	// After window expiry, rate limit should reset
	if !rl.Allow("user@example.com") {
		t.Error("expected rate limit to reset after window expiry")
	}
}

func TestContactSubmissionRateLimit(t *testing.T) {
	repo := store.New(domain.Data{})
	s := New(repo)

	email := "ratelimit@example.com"

	// First 5 submissions should succeed
	for i := 0; i < 5; i++ {
		req := httptest.NewRequest("POST", "/api/v1/contact", bytes.NewReader([]byte(`{
			"name": "Test", "email": "`+email+`", "message": "Test message"
		}`)))
		w := httptest.NewRecorder()
		s.handleContact(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("submission %d: expected status %d, got %d", i+1, http.StatusCreated, w.Code)
		}
	}

	// 6th submission should be rate limited
	req := httptest.NewRequest("POST", "/api/v1/contact", bytes.NewReader([]byte(`{
		"name": "Test", "email": "`+email+`", "message": "Test message"
	}`)))
	w := httptest.NewRecorder()
	s.handleContact(w, req)

	if w.Code != http.StatusTooManyRequests {
		t.Errorf("expected status %d, got %d", http.StatusTooManyRequests, w.Code)
	}
}
