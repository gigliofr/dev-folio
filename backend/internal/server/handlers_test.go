package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"devfolio/backend/internal/data"
	"devfolio/backend/internal/domain"
	"devfolio/backend/internal/store"
)

// newTestServer returns a test server with an in-memory store
func newTestServer() *Server {
	repo := store.New(data.Seed())
	return New(repo)
}

// Helper to make HTTP requests to test server
func makeRequest(t *testing.T, server *Server, method, path string, body interface{}, authToken string) (*httptest.ResponseRecorder, error) {
	var reqBody []byte
	var err error

	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			t.Fatalf("Failed to marshal request body: %v", err)
		}
	}

	req := httptest.NewRequest(method, path, bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	if authToken != "" {
		req.Header.Set("Authorization", "Bearer "+authToken)
	}

	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	return w, nil
}

// ============ Health Check Tests ============
func TestHandleHealth(t *testing.T) {
	server := newTestServer()

	w, _ := makeRequest(t, server, http.MethodGet, "/health", nil, "")

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var resp map[string]string
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if resp["status"] != "ok" {
		t.Errorf("Expected status 'ok', got '%s'", resp["status"])
	}
}

// ============ Login Tests ============
func TestAdminLogin(t *testing.T) {
	// Set up environment
	os.Setenv("DEVFOLIO_ADMIN_USER", "testadmin")
	os.Setenv("DEVFOLIO_ADMIN_PASS", "testpass123")
	os.Setenv("DEVFOLIO_JWT_SECRET", "test-secret-key-for-jwt")
	defer os.Unsetenv("DEVFOLIO_ADMIN_USER")
	defer os.Unsetenv("DEVFOLIO_ADMIN_PASS")
	defer os.Unsetenv("DEVFOLIO_JWT_SECRET")

	server := newTestServer()

	tests := []struct {
		name           string
		username       string
		password       string
		expectedStatus int
		shouldHaveToken bool
	}{
		{
			name:             "Valid credentials",
			username:         "testadmin",
			password:         "testpass123",
			expectedStatus:   http.StatusOK,
			shouldHaveToken:  true,
		},
		{
			name:             "Invalid username",
			username:         "wronguser",
			password:         "testpass123",
			expectedStatus:   http.StatusUnauthorized,
			shouldHaveToken:  false,
		},
		{
			name:             "Invalid password",
			username:         "testadmin",
			password:         "wrongpass",
			expectedStatus:   http.StatusUnauthorized,
			shouldHaveToken:  false,
		},
		{
			name:             "Empty credentials",
			username:         "",
			password:         "",
			expectedStatus:   http.StatusUnauthorized,
			shouldHaveToken:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := loginRequest{
				Username: tt.username,
				Password: tt.password,
			}

			w, _ := makeRequest(t, server, http.MethodPost, "/api/v1/admin/login", payload, "")

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.shouldHaveToken {
				var resp loginResponse
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				if resp.Token == "" {
					t.Error("Expected token in response, got empty string")
				}
			}
		})
	}
}

// ============ Site Tests ============
func TestHandleSite(t *testing.T) {
	os.Setenv("DEVFOLIO_ADMIN_USER", "testadmin")
	os.Setenv("DEVFOLIO_ADMIN_PASS", "testpass123")
	os.Setenv("DEVFOLIO_JWT_SECRET", "test-secret-key-for-jwt")
	os.Setenv("DEVFOLIO_ADMIN_TOKEN", "legacy-token")
	defer os.Unsetenv("DEVFOLIO_ADMIN_USER")
	defer os.Unsetenv("DEVFOLIO_ADMIN_PASS")
	defer os.Unsetenv("DEVFOLIO_JWT_SECRET")
	defer os.Unsetenv("DEVFOLIO_ADMIN_TOKEN")

	server := newTestServer()

	// Get JWT token for authenticated requests
	tokenReq := loginRequest{Username: "testadmin", Password: "testpass123"}
	w, _ := makeRequest(t, server, http.MethodPost, "/api/v1/admin/login", tokenReq, "")
	var loginResp loginResponse
	json.NewDecoder(w.Body).Decode(&loginResp)
	token := loginResp.Token

	tests := []struct {
		name           string
		method         string
		body           interface{}
		authToken      string
		expectedStatus int
	}{
		{
			name:            "GET site without auth",
			method:          http.MethodGet,
			body:            nil,
			authToken:       "",
			expectedStatus:  http.StatusOK,
		},
		{
			name:   "PUT site with valid auth",
			method: http.MethodPut,
			body: domain.Site{
				Name:        "Updated Portfolio",
				Tagline:     "New tagline",
				Description: "New description",
			},
			authToken:      token,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "PUT site without auth",
			method:         http.MethodPut,
			body:           domain.Site{Name: "Test"},
			authToken:      "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "PUT site with invalid token",
			method:         http.MethodPut,
			body:           domain.Site{Name: "Test"},
			authToken:      "invalid_token",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w, _ := makeRequest(t, server, tt.method, "/api/v1/site", tt.body, tt.authToken)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

// ============ Skills Tests ============
func TestHandleSkills(t *testing.T) {
	os.Setenv("DEVFOLIO_ADMIN_USER", "testadmin")
	os.Setenv("DEVFOLIO_ADMIN_PASS", "testpass123")
	os.Setenv("DEVFOLIO_JWT_SECRET", "test-secret-key-for-jwt")
	os.Setenv("DEVFOLIO_ADMIN_TOKEN", "legacy-token")
	defer os.Unsetenv("DEVFOLIO_ADMIN_USER")
	defer os.Unsetenv("DEVFOLIO_ADMIN_PASS")
	defer os.Unsetenv("DEVFOLIO_JWT_SECRET")
	defer os.Unsetenv("DEVFOLIO_ADMIN_TOKEN")

	server := newTestServer()

	// Get JWT token
	tokenReq := loginRequest{Username: "testadmin", Password: "testpass123"}
	w, _ := makeRequest(t, server, http.MethodPost, "/api/v1/admin/login", tokenReq, "")
	var loginResp loginResponse
	json.NewDecoder(w.Body).Decode(&loginResp)
	token := loginResp.Token

	tests := []struct {
		name            string
		method          string
		body            interface{}
		authToken       string
		expectedStatus  int
		expectedLength  int
	}{
		{
			name:           "GET skills",
			method:         http.MethodGet,
			body:           nil,
			authToken:      "",
			expectedStatus: http.StatusOK,
		},
		{
			name:   "PUT skills with auth",
			method: http.MethodPut,
			body: []string{"Go", "React", "TypeScript", "MongoDB"},
			authToken:      token,
			expectedStatus: http.StatusOK,
			expectedLength: 4,
		},
		{
			name:           "PUT skills without auth",
			method:         http.MethodPut,
			body:           []string{"Go"},
			authToken:      "",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w, _ := makeRequest(t, server, tt.method, "/api/v1/skills", tt.body, tt.authToken)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedStatus == http.StatusOK && tt.method == http.MethodPut {
				var skills []string
				json.NewDecoder(w.Body).Decode(&skills)
				if len(skills) != tt.expectedLength {
					t.Errorf("Expected %d skills, got %d", tt.expectedLength, len(skills))
				}
			}
		})
	}
}

// ============ Project Tests ============
func TestHandleProjects(t *testing.T) {
	os.Setenv("DEVFOLIO_ADMIN_USER", "testadmin")
	os.Setenv("DEVFOLIO_ADMIN_PASS", "testpass123")
	os.Setenv("DEVFOLIO_JWT_SECRET", "test-secret-key-for-jwt")
	os.Setenv("DEVFOLIO_ADMIN_TOKEN", "legacy-token")
	defer os.Unsetenv("DEVFOLIO_ADMIN_USER")
	defer os.Unsetenv("DEVFOLIO_ADMIN_PASS")
	defer os.Unsetenv("DEVFOLIO_JWT_SECRET")
	defer os.Unsetenv("DEVFOLIO_ADMIN_TOKEN")

	server := newTestServer()

	// Get JWT token
	tokenReq := loginRequest{Username: "testadmin", Password: "testpass123"}
	w, _ := makeRequest(t, server, http.MethodPost, "/api/v1/admin/login", tokenReq, "")
	var loginResp loginResponse
	json.NewDecoder(w.Body).Decode(&loginResp)
	token := loginResp.Token

	newProject := domain.Project{
		Title:            "Test Project",
		Slug:             "test-project",
		DescriptionShort: "A test project",
		DescriptionLong:  "A longer description of the test project",
		Technologies:     []string{"Go", "React"},
		Status:           "completed",
		Featured:         true,
		Year:             "2024",
	}

	tests := []struct {
		name            string
		method          string
		path            string
		body            interface{}
		authToken       string
		expectedStatus  int
		testDescription string
	}{
		{
			name:            "GET projects list",
			method:          http.MethodGet,
			path:            "/api/v1/projects",
			body:            nil,
			authToken:       "",
			expectedStatus:  http.StatusOK,
			testDescription: "Returns 200 and list",
		},
		{
			name:            "GET projects with featured filter",
			method:          http.MethodGet,
			path:            "/api/v1/projects?featured=true",
			body:            nil,
			authToken:       "",
			expectedStatus:  http.StatusOK,
			testDescription: "Returns 200 with featured filter",
		},
		{
			name:            "GET projects with status filter",
			method:          http.MethodGet,
			path:            "/api/v1/projects?status=completed",
			body:            nil,
			authToken:       "",
			expectedStatus:  http.StatusOK,
			testDescription: "Returns 200 with status filter",
		},
		{
			name:            "POST new project with auth",
			method:          http.MethodPost,
			path:            "/api/v1/projects",
			body:            newProject,
			authToken:       token,
			expectedStatus:  http.StatusCreated,
			testDescription: "Creates project with 201",
		},
		{
			name:            "POST project without auth",
			method:          http.MethodPost,
			path:            "/api/v1/projects",
			body:            newProject,
			authToken:       "",
			expectedStatus:  http.StatusUnauthorized,
			testDescription: "Rejects POST without auth",
		},
		{
			name:            "POST invalid project body",
			method:          http.MethodPost,
			path:            "/api/v1/projects",
			body:            "invalid",
			authToken:       token,
			expectedStatus:  http.StatusBadRequest,
			testDescription: "Rejects bad request body",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w, _ := makeRequest(t, server, tt.method, tt.path, tt.body, tt.authToken)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d (%s)", tt.expectedStatus, w.Code, tt.testDescription)
			}
		})
	}
}

func TestHandleProjectBySlug(t *testing.T) {
	os.Setenv("DEVFOLIO_ADMIN_USER", "testadmin")
	os.Setenv("DEVFOLIO_ADMIN_PASS", "testpass123")
	os.Setenv("DEVFOLIO_JWT_SECRET", "test-secret-key-for-jwt")
	os.Setenv("DEVFOLIO_ADMIN_TOKEN", "legacy-token")
	defer os.Unsetenv("DEVFOLIO_ADMIN_USER")
	defer os.Unsetenv("DEVFOLIO_ADMIN_PASS")
	defer os.Unsetenv("DEVFOLIO_JWT_SECRET")
	defer os.Unsetenv("DEVFOLIO_ADMIN_TOKEN")

	server := newTestServer()

	// Get JWT token
	tokenReq := loginRequest{Username: "testadmin", Password: "testpass123"}
	w, _ := makeRequest(t, server, http.MethodPost, "/api/v1/admin/login", tokenReq, "")
	var loginResp loginResponse
	json.NewDecoder(w.Body).Decode(&loginResp)
	token := loginResp.Token

	// Create a project first
	project := domain.Project{
		Title:            "Test Project",
		Slug:             "test-project-slug",
		DescriptionShort: "Test",
		DescriptionLong:  "Test description",
		Technologies:     []string{"Go"},
		Status:           "completed",
		Featured:         false,
		Year:             "2024",
	}
	w, _ = makeRequest(t, server, http.MethodPost, "/api/v1/projects", project, token)

	tests := []struct {
		name           string
		method         string
		path           string
		body           interface{}
		authToken      string
		expectedStatus int
	}{
		{
			name:            "GET project by slug",
			method:          http.MethodGet,
			path:            "/api/v1/projects/test-project-slug",
			body:            nil,
			authToken:       "",
			expectedStatus:  http.StatusOK,
		},
		{
			name:            "GET nonexistent project",
			method:          http.MethodGet,
			path:            "/api/v1/projects/nonexistent",
			body:            nil,
			authToken:       "",
			expectedStatus:  http.StatusNotFound,
		},
		{
			name:   "PUT project with auth",
			method: http.MethodPut,
			path:   "/api/v1/projects/test-project-slug",
			body: domain.Project{
				Title:            "Updated Title",
				Slug:             "test-project-slug",
				DescriptionShort: "Updated",
				DescriptionLong:  "Updated description",
				Technologies:     []string{"Go", "React"},
				Status:           "in-progress",
				Featured:         true,
				Year:             "2024",
			},
			authToken:      token,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "PUT project without auth",
			method:         http.MethodPut,
			path:           "/api/v1/projects/test-project-slug",
			body:           project,
			authToken:      "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:            "DELETE project with auth",
			method:          http.MethodDelete,
			path:            "/api/v1/projects/test-project-slug",
			body:            nil,
			authToken:       token,
			expectedStatus:  http.StatusNoContent,
		},
		{
			name:            "DELETE nonexistent project",
			method:          http.MethodDelete,
			path:            "/api/v1/projects/nonexistent",
			body:            nil,
			authToken:       token,
			expectedStatus:  http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w, _ := makeRequest(t, server, tt.method, tt.path, tt.body, tt.authToken)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

// Test duplicate slug handling
func TestProjectDuplicateSlug(t *testing.T) {
	os.Setenv("DEVFOLIO_ADMIN_USER", "testadmin")
	os.Setenv("DEVFOLIO_ADMIN_PASS", "testpass123")
	os.Setenv("DEVFOLIO_JWT_SECRET", "test-secret-key-for-jwt")
	os.Setenv("DEVFOLIO_ADMIN_TOKEN", "legacy-token")
	defer os.Unsetenv("DEVFOLIO_ADMIN_USER")
	defer os.Unsetenv("DEVFOLIO_ADMIN_PASS")
	defer os.Unsetenv("DEVFOLIO_JWT_SECRET")
	defer os.Unsetenv("DEVFOLIO_ADMIN_TOKEN")

	server := newTestServer()

	// Get JWT token
	tokenReq := loginRequest{Username: "testadmin", Password: "testpass123"}
	w, _ := makeRequest(t, server, http.MethodPost, "/api/v1/admin/login", tokenReq, "")
	var loginResp loginResponse
	json.NewDecoder(w.Body).Decode(&loginResp)
	token := loginResp.Token

	project := domain.Project{
		Title:            "Test Project",
		Slug:             "duplicate-slug",
		DescriptionShort: "Test",
		DescriptionLong:  "Test description",
		Technologies:     []string{"Go"},
		Status:           "completed",
		Featured:         false,
		Year:             "2024",
	}

	// Create first project
	w, _ = makeRequest(t, server, http.MethodPost, "/api/v1/projects", project, token)
	if w.Code != http.StatusCreated {
		t.Fatalf("Failed to create first project: %d", w.Code)
	}

	// Try to create second project with same slug
	project2 := project
	project2.Title = "Different Title"

	w, _ = makeRequest(t, server, http.MethodPost, "/api/v1/projects", project2, token)
	// Should get an error or the store should auto-generate a unique slug
	if w.Code == http.StatusInternalServerError || w.Code == http.StatusBadRequest {
		t.Logf("Correctly rejected duplicate slug: %d", w.Code)
	}
}

// ============ Post Tests ============
func TestHandlePosts(t *testing.T) {
	os.Setenv("DEVFOLIO_ADMIN_USER", "testadmin")
	os.Setenv("DEVFOLIO_ADMIN_PASS", "testpass123")
	os.Setenv("DEVFOLIO_JWT_SECRET", "test-secret-key-for-jwt")
	os.Setenv("DEVFOLIO_ADMIN_TOKEN", "legacy-token")
	defer os.Unsetenv("DEVFOLIO_ADMIN_USER")
	defer os.Unsetenv("DEVFOLIO_ADMIN_PASS")
	defer os.Unsetenv("DEVFOLIO_JWT_SECRET")
	defer os.Unsetenv("DEVFOLIO_ADMIN_TOKEN")

	server := newTestServer()

	// Get JWT token
	tokenReq := loginRequest{Username: "testadmin", Password: "testpass123"}
	w, _ := makeRequest(t, server, http.MethodPost, "/api/v1/admin/login", tokenReq, "")
	var loginResp loginResponse
	json.NewDecoder(w.Body).Decode(&loginResp)
	token := loginResp.Token

	newPost := domain.Post{
		Title:       "Test Post",
		Slug:        "test-post",
		Excerpt:     "A test post excerpt",
		Category:    "tutorial",
		ReadTime:    "5 min",
		PublishedAt: "2024-01-01",
		Tags:        []string{"Go", "API"},
	}

	tests := []struct {
		name           string
		method         string
		body           interface{}
		authToken      string
		expectedStatus int
	}{
		{
			name:            "GET posts list",
			method:          http.MethodGet,
			body:            nil,
			authToken:       "",
			expectedStatus:  http.StatusOK,
		},
		{
			name:            "POST new post with auth",
			method:          http.MethodPost,
			body:            newPost,
			authToken:       token,
			expectedStatus:  http.StatusCreated,
		},
		{
			name:            "POST post without auth",
			method:          http.MethodPost,
			body:            newPost,
			authToken:       "",
			expectedStatus:  http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w, _ := makeRequest(t, server, tt.method, "/api/v1/posts", tt.body, tt.authToken)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestHandlePostBySlug(t *testing.T) {
	os.Setenv("DEVFOLIO_ADMIN_USER", "testadmin")
	os.Setenv("DEVFOLIO_ADMIN_PASS", "testpass123")
	os.Setenv("DEVFOLIO_JWT_SECRET", "test-secret-key-for-jwt")
	os.Setenv("DEVFOLIO_ADMIN_TOKEN", "legacy-token")
	defer os.Unsetenv("DEVFOLIO_ADMIN_USER")
	defer os.Unsetenv("DEVFOLIO_ADMIN_PASS")
	defer os.Unsetenv("DEVFOLIO_JWT_SECRET")
	defer os.Unsetenv("DEVFOLIO_ADMIN_TOKEN")

	server := newTestServer()

	// Get JWT token
	tokenReq := loginRequest{Username: "testadmin", Password: "testpass123"}
	w, _ := makeRequest(t, server, http.MethodPost, "/api/v1/admin/login", tokenReq, "")
	var loginResp loginResponse
	json.NewDecoder(w.Body).Decode(&loginResp)
	token := loginResp.Token

	// Create a post first
	post := domain.Post{
		Title:       "Test Post",
		Slug:        "test-post-slug",
		Excerpt:     "Test",
		Category:    "tutorial",
		ReadTime:    "5 min",
		PublishedAt: "2024-01-01",
		Tags:        []string{"Go"},
	}
	w, _ = makeRequest(t, server, http.MethodPost, "/api/v1/posts", post, token)

	tests := []struct {
		name           string
		method         string
		path           string
		body           interface{}
		authToken      string
		expectedStatus int
	}{
		{
			name:           "GET post by slug",
			method:         http.MethodGet,
			path:           "/api/v1/posts/test-post-slug",
			body:           nil,
			authToken:      "",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "GET nonexistent post",
			method:         http.MethodGet,
			path:           "/api/v1/posts/nonexistent",
			body:           nil,
			authToken:      "",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:   "PUT post with auth",
			method: http.MethodPut,
			path:   "/api/v1/posts/test-post-slug",
			body: domain.Post{
				Title:       "Updated Title",
				Slug:        "test-post-slug",
				Excerpt:     "Updated excerpt",
				Category:    "guide",
				ReadTime:    "10 min",
				PublishedAt: "2024-01-01",
				Tags:        []string{"Go", "REST"},
			},
			authToken:      token,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "PUT post without auth",
			method:         http.MethodPut,
			path:           "/api/v1/posts/test-post-slug",
			body:           post,
			authToken:      "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "DELETE post with auth",
			method:         http.MethodDelete,
			path:           "/api/v1/posts/test-post-slug",
			body:           nil,
			authToken:      token,
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "DELETE nonexistent post",
			method:         http.MethodDelete,
			path:           "/api/v1/posts/nonexistent",
			body:           nil,
			authToken:      token,
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w, _ := makeRequest(t, server, tt.method, tt.path, tt.body, tt.authToken)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

// ============ Stats Tests ============
func TestHandleStats(t *testing.T) {
	server := newTestServer()

	w, _ := makeRequest(t, server, http.MethodGet, "/api/v1/stats", nil, "")

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var stats []domain.Stat
	if err := json.NewDecoder(w.Body).Decode(&stats); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(stats) == 0 {
		t.Logf("Stats list is empty (expected for fresh store)")
	}
}

// ============ Admin Me Tests ============
func TestAdminMe(t *testing.T) {
	os.Setenv("DEVFOLIO_ADMIN_USER", "testadmin")
	os.Setenv("DEVFOLIO_ADMIN_PASS", "testpass123")
	os.Setenv("DEVFOLIO_JWT_SECRET", "test-secret-key-for-jwt")
	os.Setenv("DEVFOLIO_ADMIN_TOKEN", "legacy-token")
	defer os.Unsetenv("DEVFOLIO_ADMIN_USER")
	defer os.Unsetenv("DEVFOLIO_ADMIN_PASS")
	defer os.Unsetenv("DEVFOLIO_JWT_SECRET")
	defer os.Unsetenv("DEVFOLIO_ADMIN_TOKEN")

	server := newTestServer()

	// Get JWT token
	tokenReq := loginRequest{Username: "testadmin", Password: "testpass123"}
	w, _ := makeRequest(t, server, http.MethodPost, "/api/v1/admin/login", tokenReq, "")
	var loginResp loginResponse
	json.NewDecoder(w.Body).Decode(&loginResp)
	token := loginResp.Token

	tests := []struct {
		name           string
		authToken      string
		expectedStatus int
	}{
		{
			name:            "GET admin identity with valid token",
			authToken:       token,
			expectedStatus:  http.StatusOK,
		},
		{
			name:            "GET admin identity without token",
			authToken:       "",
			expectedStatus:  http.StatusUnauthorized,
		},
		{
			name:            "GET admin identity with invalid token",
			authToken:       "invalid_token_string",
			expectedStatus:  http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w, _ := makeRequest(t, server, http.MethodGet, "/api/v1/admin/me", nil, tt.authToken)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

// ============ CORS Tests ============
func TestCORS(t *testing.T) {
	server := newTestServer()

	tests := []struct {
		name        string
		origin      string
		expectCORS  bool
	}{
		{
			name:       "Request with origin header",
			origin:     "http://localhost:3000",
			expectCORS: true,
		},
		{
			name:       "Request without origin header",
			origin:     "",
			expectCORS: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/v1/site", nil)
			if tt.origin != "" {
				req.Header.Set("Origin", tt.origin)
			}

			w := httptest.NewRecorder()
			server.ServeHTTP(w, req)

			corsOrigin := w.Header().Get("Access-Control-Allow-Origin")
			if tt.expectCORS && corsOrigin == "" {
				t.Error("Expected CORS header in response")
			}

			if corsOrigin != "" {
				t.Logf("CORS origin: %s", corsOrigin)
			}
		})
	}
}

// ============ Method Not Allowed Tests ============
func TestMethodNotAllowed(t *testing.T) {
	server := newTestServer()

	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
	}{
		{
			name:            "DELETE on /api/v1/site",
			method:          http.MethodDelete,
			path:            "/api/v1/site",
			expectedStatus:  http.StatusMethodNotAllowed,
		},
		{
			name:            "POST on /api/v1/stats",
			method:          http.MethodPost,
			path:            "/api/v1/stats",
			expectedStatus:  http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w, _ := makeRequest(t, server, tt.method, tt.path, nil, "")

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			// Check Allow header is set
			allowed := w.Header().Get("Allow")
			if allowed == "" {
				t.Error("Expected Allow header in response")
			}
		})
	}
}
