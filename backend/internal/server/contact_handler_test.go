package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"devfolio/backend/internal/domain"
	"devfolio/backend/internal/store"
)

func TestHandleContact(t *testing.T) {
	repo := store.New(domain.Data{})
	s := New(repo)

	tests := []struct {
		name           string
		method         string
		body           map[string]string
		expectedStatus int
	}{
		{
			name:           "POST valid contact",
			method:         "POST",
			body:           map[string]string{"name": "John", "email": "john@example.com", "message": "Hello"},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "POST missing name",
			method:         "POST",
			body:           map[string]string{"email": "john@example.com", "message": "Hello"},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "POST missing email",
			method:         "POST",
			body:           map[string]string{"name": "John", "message": "Hello"},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "POST missing message",
			method:         "POST",
			body:           map[string]string{"name": "John", "email": "john@example.com"},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "GET not allowed",
			method:         "GET",
			expectedStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var bodyBytes []byte
			if tc.body != nil {
				bodyBytes, _ = json.Marshal(tc.body)
			}
			req := httptest.NewRequest(tc.method, "/api/v1/contact", bytes.NewReader(bodyBytes))
			w := httptest.NewRecorder()

			s.handleContact(w, req)

			if w.Code != tc.expectedStatus {
				t.Errorf("expected status %d, got %d", tc.expectedStatus, w.Code)
			}
		})
	}
}

func TestHandleListContacts(t *testing.T) {
	// Set up JWT secret for testing
	os.Setenv("DEVFOLIO_JWT_SECRET", "test-secret-key-for-testing")
	defer os.Unsetenv("DEVFOLIO_JWT_SECRET")

	repo := store.New(domain.Data{})
	s := New(repo)

	// Add a test contact
	repo.SaveContactSubmission(domain.ContactSubmission{
		Name:      "Alice",
		Email:     "alice@example.com",
		Message:   "Test message",
		CreatedAt: time.Now(),
	})

	t.Run("GET without auth", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/contact", nil)
		w := httptest.NewRecorder()

		s.handleListContacts(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("GET with valid auth token", func(t *testing.T) {
		token, _ := createToken("admin", 1*time.Hour)
		req := httptest.NewRequest("GET", "/api/v1/contact", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		s.handleListContacts(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var resp map[string]interface{}
		json.NewDecoder(w.Body).Decode(&resp)
		submissions, ok := resp["submissions"].([]interface{})
		if !ok {
			t.Errorf("expected submissions array in response")
		}
		if len(submissions) != 1 {
			t.Errorf("expected 1 submission, got %d", len(submissions))
		}
	})

	t.Run("POST not allowed", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/v1/contact", nil)
		w := httptest.NewRecorder()

		s.handleListContacts(w, req)

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
		}
	})
}
