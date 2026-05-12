package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"devfolio/backend/internal/domain"
	"devfolio/backend/internal/store"
)

func (s *Server) handleHealth(responseWriter http.ResponseWriter, _ *http.Request) {
	s.writeJSON(responseWriter, http.StatusOK, map[string]string{
		"status":  "ok",
		"service": "devfolio-api",
	})
}

func (s *Server) handleSite(responseWriter http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		s.writeJSON(responseWriter, http.StatusOK, s.store.Site())
	case http.MethodPut:
		if !s.requireAdmin(responseWriter, request) {
			return
		}
		var payload domain.Site
		if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
			s.writeJSON(responseWriter, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
			return
		}
		s.writeJSON(responseWriter, http.StatusOK, s.store.UpdateSite(payload))
	default:
		methodNotAllowed(responseWriter, http.MethodGet, http.MethodPut)
	}
}

func (s *Server) handleStats(responseWriter http.ResponseWriter, _ *http.Request) {
	s.writeJSON(responseWriter, http.StatusOK, s.store.Stats())
}

func (s *Server) handleSkills(responseWriter http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		s.writeJSON(responseWriter, http.StatusOK, s.store.Skills())
	case http.MethodPut:
		if !s.requireAdmin(responseWriter, request) {
			return
		}
		var payload []string
		if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
			s.writeJSON(responseWriter, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
			return
		}
		s.writeJSON(responseWriter, http.StatusOK, s.store.UpdateSkills(payload))
	default:
		methodNotAllowed(responseWriter, http.MethodGet, http.MethodPut)
	}
}

func (s *Server) handleProjects(responseWriter http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		featuredOnly := request.URL.Query().Get("featured") == "true"
		statusFilter := strings.TrimSpace(request.URL.Query().Get("status"))
		s.writeJSON(responseWriter, http.StatusOK, s.store.ListProjects(featuredOnly, statusFilter))
	case http.MethodPost:
		if !s.requireAdmin(responseWriter, request) {
			return
		}
		var payload domain.Project
		if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
			s.writeJSON(responseWriter, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
			return
		}
		s.writeJSON(responseWriter, http.StatusCreated, s.store.CreateProject(payload))
	default:
		methodNotAllowed(responseWriter, http.MethodGet, http.MethodPost)
	}
}

func (s *Server) handleProjectBySlug(responseWriter http.ResponseWriter, request *http.Request) {
	slug := strings.TrimPrefix(request.URL.Path, "/api/v1/projects/")
	if slug == "" || strings.Contains(slug, "/") {
		http.NotFound(responseWriter, request)
		return
	}

	switch request.Method {
	case http.MethodGet:
		project, found := s.store.GetProject(slug)
		if !found {
			http.NotFound(responseWriter, request)
			return
		}
		s.writeJSON(responseWriter, http.StatusOK, project)
	case http.MethodPut:
		if !s.requireAdmin(responseWriter, request) {
			return
		}
		var payload domain.Project
		if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
			s.writeJSON(responseWriter, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
			return
		}
		updated, err := s.store.UpdateProject(slug, payload)
		if err != nil {
			if errors.Is(err, store.ErrNotFound) {
				http.NotFound(responseWriter, request)
				return
			}
			s.writeJSON(responseWriter, http.StatusInternalServerError, map[string]string{"error": "could not update project"})
			return
		}
		s.writeJSON(responseWriter, http.StatusOK, updated)
	case http.MethodDelete:
		if !s.requireAdmin(responseWriter, request) {
			return
		}
		if !s.store.DeleteProject(slug) {
			http.NotFound(responseWriter, request)
			return
		}
		responseWriter.WriteHeader(http.StatusNoContent)
	default:
		methodNotAllowed(responseWriter, http.MethodGet, http.MethodPut, http.MethodDelete)
	}
}

func (s *Server) handlePosts(responseWriter http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		s.writeJSON(responseWriter, http.StatusOK, s.store.ListPosts())
	case http.MethodPost:
		if !s.requireAdmin(responseWriter, request) {
			return
		}
		var payload domain.Post
		if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
			s.writeJSON(responseWriter, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
			return
		}
		s.writeJSON(responseWriter, http.StatusCreated, s.store.CreatePost(payload))
	default:
		methodNotAllowed(responseWriter, http.MethodGet, http.MethodPost)
	}
}

func (s *Server) handlePostBySlug(responseWriter http.ResponseWriter, request *http.Request) {
	slug := strings.TrimPrefix(request.URL.Path, "/api/v1/posts/")
	if slug == "" || strings.Contains(slug, "/") {
		http.NotFound(responseWriter, request)
		return
	}

	switch request.Method {
	case http.MethodGet:
		post, found := s.store.GetPost(slug)
		if !found {
			http.NotFound(responseWriter, request)
			return
		}
		s.writeJSON(responseWriter, http.StatusOK, post)
	case http.MethodPut:
		if !s.requireAdmin(responseWriter, request) {
			return
		}
		var payload domain.Post
		if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
			s.writeJSON(responseWriter, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
			return
		}
		updated, err := s.store.UpdatePost(slug, payload)
		if err != nil {
			if errors.Is(err, store.ErrNotFound) {
				http.NotFound(responseWriter, request)
				return
			}
			s.writeJSON(responseWriter, http.StatusInternalServerError, map[string]string{"error": "could not update post"})
			return
		}
		s.writeJSON(responseWriter, http.StatusOK, updated)
	case http.MethodDelete:
		if !s.requireAdmin(responseWriter, request) {
			return
		}
		if !s.store.DeletePost(slug) {
			http.NotFound(responseWriter, request)
			return
		}
		responseWriter.WriteHeader(http.StatusNoContent)
	default:
		methodNotAllowed(responseWriter, http.MethodGet, http.MethodPut, http.MethodDelete)
	}
}

func methodNotAllowed(responseWriter http.ResponseWriter, allowed ...string) {
	responseWriter.Header().Set("Allow", strings.Join(allowed, ", "))
	responseWriter.WriteHeader(http.StatusMethodNotAllowed)
}