package server

import (
	"encoding/json"
	"net/http"
	"time"

	"devfolio/backend/internal/store"
)

type Server struct {
	store       store.Repository
	mux         *http.ServeMux
	rateLimiter *RateLimiter
}

func New(repository store.Repository) *Server {
	server := &Server{
		store:       repository,
		mux:         http.NewServeMux(),
		rateLimiter: NewRateLimiter(1*time.Hour, 5), // 5 submissions per hour per email
	}
	server.routes()
	return server
}

func (s *Server) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
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
	if origin == "" {
		responseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	} else {
		responseWriter.Header().Set("Access-Control-Allow-Origin", origin)
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

