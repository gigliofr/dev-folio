package store

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"devfolio/backend/internal/domain"
)

var ErrNotFound = errors.New("resource not found")

type Store struct {
	mu       sync.RWMutex
	data     domain.Data
	contacts []domain.ContactSubmission
}

func New(seed domain.Data) *Store {
	return &Store{data: seed}
}

func (s *Store) Snapshot() domain.Data {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return cloneData(s.data)
}

func (s *Store) Site() domain.Site {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.data.Site
}

func (s *Store) UpdateSite(site domain.Site) domain.Site {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data.Site = site
	return s.data.Site
}

func (s *Store) Stats() []domain.Stat {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]domain.Stat(nil), s.data.Stats...)
}

func (s *Store) Skills() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]string(nil), s.data.Skills...)
}

func (s *Store) UpdateSkills(skills []string) []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data.Skills = append([]string(nil), skills...)
	return append([]string(nil), s.data.Skills...)
}

func (s *Store) ListProjects(featuredOnly bool, statusFilter string) []domain.Project {
	s.mu.RLock()
	defer s.mu.RUnlock()
	filtered := make([]domain.Project, 0, len(s.data.Projects))
	for _, project := range s.data.Projects {
		if featuredOnly && !project.Featured {
			continue
		}
		if statusFilter != "" && project.Status != statusFilter {
			continue
		}
		filtered = append(filtered, project)
	}
	return filtered
}

func (s *Store) GetProject(slug string) (domain.Project, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, project := range s.data.Projects {
		if project.Slug == slug {
			return project, true
		}
	}
	return domain.Project{}, false
}

func (s *Store) CreateProject(project domain.Project) domain.Project {
	s.mu.Lock()
	defer s.mu.Unlock()
	project.Slug = uniqueSlug(project.Slug, project.Title, projectExists(s.data.Projects))
	s.data.Projects = append([]domain.Project{project}, s.data.Projects...)
	return project
}

func (s *Store) UpdateProject(slug string, project domain.Project) (domain.Project, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for index, current := range s.data.Projects {
		if current.Slug != slug {
			continue
		}
		project.Slug = slug
		s.data.Projects[index] = project
		return project, nil
	}
	return domain.Project{}, ErrNotFound
}

func (s *Store) DeleteProject(slug string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	for index, project := range s.data.Projects {
		if project.Slug != slug {
			continue
		}
		s.data.Projects = append(s.data.Projects[:index], s.data.Projects[index+1:]...)
		return true
	}
	return false
}

func (s *Store) ListPosts() []domain.Post {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]domain.Post(nil), s.data.Posts...)
}

func (s *Store) GetPost(slug string) (domain.Post, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, post := range s.data.Posts {
		if post.Slug == slug {
			return post, true
		}
	}
	return domain.Post{}, false
}

func (s *Store) CreatePost(post domain.Post) domain.Post {
	s.mu.Lock()
	defer s.mu.Unlock()
	post.Slug = uniqueSlug(post.Slug, post.Title, postExists(s.data.Posts))
	s.data.Posts = append([]domain.Post{post}, s.data.Posts...)
	return post
}

func (s *Store) UpdatePost(slug string, post domain.Post) (domain.Post, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for index, current := range s.data.Posts {
		if current.Slug != slug {
			continue
		}
		post.Slug = slug
		s.data.Posts[index] = post
		return post, nil
	}
	return domain.Post{}, ErrNotFound
}

func (s *Store) DeletePost(slug string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	for index, post := range s.data.Posts {
		if post.Slug != slug {
			continue
		}
		s.data.Posts = append(s.data.Posts[:index], s.data.Posts[index+1:]...)
		return true
	}
	return false
}

func (s *Store) SaveContactSubmission(submission domain.ContactSubmission) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.contacts = append(s.contacts, submission)
	return nil
}

func (s *Store) ListContactSubmissions() []domain.ContactSubmission {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]domain.ContactSubmission(nil), s.contacts...)
}

func cloneData(data domain.Data) domain.Data {
	cloned := data
	cloned.Projects = append([]domain.Project(nil), data.Projects...)
	cloned.Posts = append([]domain.Post(nil), data.Posts...)
	cloned.Stats = append([]domain.Stat(nil), data.Stats...)
	cloned.Skills = append([]string(nil), data.Skills...)
	return cloned
}

func projectExists(projects []domain.Project) func(string) bool {
	return func(slug string) bool {
		for _, project := range projects {
			if project.Slug == slug {
				return true
			}
		}
		return false
	}
}

func postExists(posts []domain.Post) func(string) bool {
	return func(slug string) bool {
		for _, post := range posts {
			if post.Slug == slug {
				return true
			}
		}
		return false
	}
}

func uniqueSlug(currentSlug string, title string, exists func(string) bool) string {
	baseSlug := strings.TrimSpace(currentSlug)
	if baseSlug == "" {
		baseSlug = slugify(title)
	}
	if baseSlug == "" {
		baseSlug = "item"
	}
	if !exists(baseSlug) {
		return baseSlug
	}

	for suffix := 2; suffix < 1000; suffix++ {
		candidate := fmt.Sprintf("%s-%d", baseSlug, suffix)
		if !exists(candidate) {
			return candidate
		}
	}

	return baseSlug
}

func slugify(value string) string {
	var builder strings.Builder
	lastDash := false
	for _, r := range strings.ToLower(value) {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			builder.WriteRune(r)
			lastDash = false
		case r == ' ' || r == '-' || r == '_' || r == '/':
			if !lastDash && builder.Len() > 0 {
				builder.WriteByte('-')
				lastDash = true
			}
		}
	}
	return strings.Trim(builder.String(), "-")
}