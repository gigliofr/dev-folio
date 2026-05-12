package server

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func (s *Server) handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	// Require admin authentication
	if !s.requireAdmin(w, r) {
		return
	}

	// Parse multipart form (max 10MB)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		s.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid form data"})
		return
	}

	// Get the file from form
	file, handler, err := r.FormFile("file")
	if err != nil {
		s.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "file field required"})
		return
	}
	defer file.Close()

	// Validate file type
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/webp": true,
		"image/gif":  true,
	}

	if !allowedTypes[handler.Header.Get("Content-Type")] {
		s.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "only image files allowed"})
		return
	}

	// Create uploads directory if it doesn't exist
	uploadDir := "public/uploads"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		s.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not create upload directory"})
		return
	}

	// Generate unique filename
	ext := filepath.Ext(handler.Filename)
	filename := fmt.Sprintf("upload_%d%s", time.Now().UnixNano(), ext)
	filepath := filepath.Join(uploadDir, filename)

	// Save file
	dst, err := os.Create(filepath)
	if err != nil {
		s.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not save file"})
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		s.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not write file"})
		return
	}

	// Return relative URL
	url := strings.TrimPrefix(strings.ReplaceAll(filepath, "\\", "/"), "public")
	s.writeJSON(w, http.StatusOK, map[string]string{"url": url})
}
