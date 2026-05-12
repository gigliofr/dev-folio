package store

import "devfolio/backend/internal/domain"

type Repository interface {
	Site() domain.Site
	UpdateSite(site domain.Site) domain.Site
	Stats() []domain.Stat
	Skills() []string
	UpdateSkills(skills []string) []string
	ListProjects(featuredOnly bool, statusFilter string) []domain.Project
	GetProject(slug string) (domain.Project, bool)
	CreateProject(project domain.Project) domain.Project
	UpdateProject(slug string, project domain.Project) (domain.Project, error)
	DeleteProject(slug string) bool
	ListPosts() []domain.Post
	GetPost(slug string) (domain.Post, bool)
	CreatePost(post domain.Post) domain.Post
	UpdatePost(slug string, post domain.Post) (domain.Post, error)
	DeletePost(slug string) bool
	SaveContactSubmission(submission domain.ContactSubmission) error
	ListContactSubmissions() []domain.ContactSubmission
}