package server

import (
	"errors"
	"encoding/json"
	"net/http"
)

func (s *Server) handleContact(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	var req ContactRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	submission, err := s.contactService.Submit(req.Name, req.Email, req.Message)
	if err != nil {
		status := http.StatusInternalServerError
		message := "could not save submission"
		switch {
		case errors.Is(err, ErrContactInvalid):
			status = http.StatusBadRequest
			message = "name, email, and message are required"
		case errors.Is(err, ErrContactRateLimit):
			status = http.StatusTooManyRequests
			message = "too many submissions from this email, please try again later"
		case errors.Is(err, ErrContactNotify):
			status = http.StatusBadGateway
			message = "message saved but email notification failed"
		}
		s.writeJSON(w, status, map[string]string{"error": message})
		return
	}

	s.writeJSON(w, http.StatusCreated, map[string]string{"message": "submission received", "name": submission.Name})
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

	submissions := s.contactService.List()
	s.writeJSON(w, http.StatusOK, map[string]any{"submissions": submissions, "count": len(submissions)})
}
