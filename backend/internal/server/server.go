package server

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"devfolio/backend/internal/store"
)

type Server struct {
	store          store.Repository
	mux            *http.ServeMux
	contactService *ContactService
}

func New(repository store.Repository) *Server {
	server := &Server{
		store:          repository,
		mux:            http.NewServeMux(),
		contactService: NewContactService(repository, NewRateLimiter(1*time.Hour, 5), NewSMTPContactNotifierFromEnv()),
	}
	server.routes()
	return server
}

func (s *Server) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	// Security headers (applied to all responses)
	responseWriter.Header().Set("X-Content-Type-Options", "nosniff")
	responseWriter.Header().Set("X-Frame-Options", "DENY")
	responseWriter.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
	responseWriter.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
	// Basic CSP — adjust as needed for third-party assets
	responseWriter.Header().Set("Content-Security-Policy", "default-src 'self'; img-src 'self' data: https:; style-src 'self' 'unsafe-inline'; script-src 'self'; frame-ancestors 'none';")

	if request.Method == http.MethodOptions {
		s.writeCORS(responseWriter, request)
		responseWriter.WriteHeader(http.StatusNoContent)
		return
	}

	s.writeCORS(responseWriter, request)
	s.mux.ServeHTTP(responseWriter, request)
}

func (s *Server) routes() {
	s.mux.HandleFunc("GET /health", s.handleHealth)
	s.mux.HandleFunc("GET /api/v1/site", s.handleSite)
	s.mux.HandleFunc("PUT /api/v1/site", s.handleSite)
	s.mux.HandleFunc("GET /api/v1/stats", s.handleStats)
	s.mux.HandleFunc("GET /api/v1/skills", s.handleSkills)
	s.mux.HandleFunc("PUT /api/v1/skills", s.handleSkills)
	s.mux.HandleFunc("GET /api/v1/projects", s.handleProjects)
	s.mux.HandleFunc("POST /api/v1/projects", s.handleProjects)
	s.mux.HandleFunc("GET /api/v1/projects/", s.handleProjectBySlug)
	s.mux.HandleFunc("PUT /api/v1/projects/", s.handleProjectBySlug)
	s.mux.HandleFunc("DELETE /api/v1/projects/", s.handleProjectBySlug)
	s.mux.HandleFunc("GET /api/v1/posts", s.handlePosts)
	s.mux.HandleFunc("POST /api/v1/posts", s.handlePosts)
	s.mux.HandleFunc("GET /api/v1/posts/", s.handlePostBySlug)
	s.mux.HandleFunc("PUT /api/v1/posts/", s.handlePostBySlug)
	s.mux.HandleFunc("DELETE /api/v1/posts/", s.handlePostBySlug)
	// Admin auth endpoints
	s.mux.HandleFunc("POST /api/v1/admin/login", s.handleAdminLogin)
	s.mux.HandleFunc("GET /api/v1/admin/me", s.handleAdminMe)
	// Upload endpoint
	s.mux.HandleFunc("POST /api/v1/upload", s.handleUpload)
	// Contact form endpoint
	s.mux.HandleFunc("POST /api/v1/contact", s.handleContact)
	s.mux.HandleFunc("GET /api/v1/contact", s.handleListContacts)
}

func (s *Server) writeCORS(responseWriter http.ResponseWriter, request *http.Request) {
	origin := request.Header.Get("Origin")
	allowedEnv := os.Getenv("DEVFOLIO_ALLOWED_ORIGINS") // comma-separated list
	if allowedEnv == "" {
		// Default behavior: allow all
		if origin == "" {
			responseWriter.Header().Set("Access-Control-Allow-Origin", "*")
		} else {
			responseWriter.Header().Set("Access-Control-Allow-Origin", origin)
		}
	} else {
		allowed := strings.Split(allowedEnv, ",")
		found := false
		for _, a := range allowed {
			if strings.TrimSpace(a) == origin {
				found = true
				break
			}
		}
		if found {
			responseWriter.Header().Set("Access-Control-Allow-Origin", origin)
		}
		// if not found: intentionally do not set the header (browser will block)
	}
	responseWriter.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	responseWriter.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	responseWriter.Header().Set("Access-Control-Max-Age", "86400")
	responseWriter.Header().Set("Vary", "Origin")
}

func (s *Server) writeJSON(responseWriter http.ResponseWriter, statusCode int, payload any) {
	responseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
	responseWriter.WriteHeader(statusCode)
	_ = json.NewEncoder(responseWriter).Encode(payload)
}

