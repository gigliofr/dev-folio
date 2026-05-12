package server

import (
	"net/http"
	"os"
	"strings"
)

func (s *Server) requireAdmin(responseWriter http.ResponseWriter, request *http.Request) bool {
	// First, try JWT bearer
	auth := request.Header.Get("Authorization")
	if strings.HasPrefix(strings.ToLower(auth), "bearer ") {
		tokenStr := strings.TrimSpace(auth[len("bearer "):])
		if claims, err := parseAndValidateToken(tokenStr); err == nil && claims != nil {
			return true
		}
	}

	// Fallback: legacy X-Admin-Token environment variable
	adminToken := strings.TrimSpace(os.Getenv("DEVFOLIO_ADMIN_TOKEN"))
	if adminToken == "" {
		// No admin guard configured in this environment
		return true
	}
	if request.Header.Get("X-Admin-Token") == adminToken {
		return true
	}

	s.writeJSON(responseWriter, http.StatusUnauthorized, map[string]string{"error": "admin token missing or invalid"})
	return false
}