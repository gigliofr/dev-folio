package server

import (
	"encoding/json"
	"net/http"
	"time"

	"devfolio/backend/internal/domain"
)

type contactRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Message string `json:"message"`
}

func (s *Server) handleContact(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	var req contactRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	// Validate required fields
	if req.Name == "" || req.Email == "" || req.Message == "" {
		s.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "name, email, and message are required"})
		return
	}

	// Check rate limit by email
	if !s.rateLimiter.Allow(req.Email) {
		s.writeJSON(w, http.StatusTooManyRequests, map[string]string{"error": "too many submissions from this email, please try again later"})
		return
	}

	// Save submission
	submission := domain.ContactSubmission{
		Name:      req.Name,
		Email:     req.Email,
		Message:   req.Message,
		CreatedAt: time.Now(),
	}

	if err := s.store.SaveContactSubmission(submission); err != nil {
		s.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not save submission"})
		return
	}

	s.writeJSON(w, http.StatusCreated, map[string]string{"message": "submission received"})
}

func (s *Server) handleListContacts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	// Require admin authentication
	if !s.requireAdmin(w, r) {
		return
	}

	submissions := s.store.ListContactSubmissions()
	s.writeJSON(w, http.StatusOK, submissions)
}
