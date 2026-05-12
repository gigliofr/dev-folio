package server

import (
    "encoding/json"
    "net/http"
    "os"
    "time"

    "devfolio/backend/internal/domain"
)

// login payloads
type loginRequest struct {
    Username string `json:"username,omitempty"`
    Password string `json:"password,omitempty"`
    Token    string `json:"token,omitempty"` // legacy token support
}

type loginResponse struct {
    Token string `json:"token"`
}

func (s *Server) handleAdminLogin(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        methodNotAllowed(w, http.MethodPost)
        return
    }

    var req loginRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        s.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request"})
        return
    }

    // Prefer explicit username/password
    envUser := os.Getenv("DEVFOLIO_ADMIN_USER")
    envPass := os.Getenv("DEVFOLIO_ADMIN_PASS")
    if envUser != "" && envPass != "" {
        if req.Username == envUser && req.Password == envPass {
            token, _ := createToken(envUser, 24*time.Hour)
            s.writeJSON(w, http.StatusOK, loginResponse{Token: token})
            return
        }
        s.writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
        return
    }

    // Fallback: legacy admin token
    legacy := os.Getenv("DEVFOLIO_ADMIN_TOKEN")
    if legacy != "" && req.Token == legacy {
        token, _ := createToken("admin", 24*time.Hour)
        s.writeJSON(w, http.StatusOK, loginResponse{Token: token})
        return
    }

    s.writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "no admin credentials configured"})
}

// Simple endpoint to return the current authenticated admin identity
func (s *Server) handleAdminMe(w http.ResponseWriter, r *http.Request) {
    // Only allow GET
    if r.Method != http.MethodGet {
        methodNotAllowed(w, http.MethodGet)
        return
    }

    // Try JWT
    auth := r.Header.Get("Authorization")
    if len(auth) > 7 && auth[:7] == "Bearer " {
        tokenStr := auth[7:]
        if claims, err := parseAndValidateToken(tokenStr); err == nil && claims != nil {
            s.writeJSON(w, http.StatusOK, domain.Site{Description: "admin:" + claims.Subject})
            return
        }
    }

    s.writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthenticated"})
}
