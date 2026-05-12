package store

import (
	"context"
	"fmt"
	"strings"
	"time"

	"devfolio/backend/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const defaultMongoDatabase = "devfolio"

type MongoStore struct {
	client  *mongo.Client
	db      *mongo.Database
	seed    domain.Data
	timeout time.Duration
}

func NewMongo(ctx context.Context, uri string, seed domain.Data, databaseName string) (*MongoStore, error) {
	if databaseName == "" {
		databaseName = defaultMongoDatabase
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		_ = client.Disconnect(ctx)
		return nil, err
	}

	store := &MongoStore{
		client:  client,
		db:      client.Database(databaseName),
		seed:    seed,
		timeout: 5 * time.Second,
	}

	if err := store.seedIfNeeded(ctx); err != nil {
		_ = client.Disconnect(ctx)
		return nil, err
	}

	return store, nil
}

func (s *MongoStore) Site() domain.Site {
	ctx, cancel := s.context()
	defer cancel()

	var document siteDocument
	if err := s.db.Collection("site").FindOne(ctx, bson.M{"_id": "site"}).Decode(&document); err != nil {
		return s.seed.Site
	}
	return document.Site
}

func (s *MongoStore) UpdateSite(site domain.Site) domain.Site {
	ctx, cancel := s.context()
	defer cancel()

	_, _ = s.db.Collection("site").ReplaceOne(ctx, bson.M{"_id": "site"}, siteDocument{ID: "site", Site: site}, options.Replace().SetUpsert(true))
	return site
}

func (s *MongoStore) Stats() []domain.Stat {
	ctx, cancel := s.context()
	defer cancel()

	var document statsDocument
	if err := s.db.Collection("globals").FindOne(ctx, bson.M{"_id": "stats"}).Decode(&document); err != nil {
		return append([]domain.Stat(nil), s.seed.Stats...)
	}
	return document.Stats
}

func (s *MongoStore) Skills() []string {
	ctx, cancel := s.context()
	defer cancel()

	var document skillsDocument
	if err := s.db.Collection("globals").FindOne(ctx, bson.M{"_id": "skills"}).Decode(&document); err != nil {
		return append([]string(nil), s.seed.Skills...)
	}
	return document.Skills
}

func (s *MongoStore) UpdateSkills(skills []string) []string {
	ctx, cancel := s.context()
	defer cancel()

	_, _ = s.db.Collection("globals").ReplaceOne(ctx, bson.M{"_id": "skills"}, skillsDocument{ID: "skills", Skills: append([]string(nil), skills...)}, options.Replace().SetUpsert(true))
	return append([]string(nil), skills...)
}

func (s *MongoStore) ListProjects(featuredOnly bool, statusFilter string) []domain.Project {
	ctx, cancel := s.context()
	defer cancel()

	filter := bson.M{}
	if featuredOnly {
		filter["featured"] = true
	}
	if statusFilter != "" {
		filter["status"] = statusFilter
	}

	cursor, err := s.db.Collection("projects").Find(ctx, filter)
	if err != nil {
		return append([]domain.Project(nil), s.seed.Projects...)
	}
	defer cursor.Close(ctx)

	projects := make([]domain.Project, 0)
	if err := cursor.All(ctx, &projects); err != nil {
		return append([]domain.Project(nil), s.seed.Projects...)
	}
	return projects
}

func (s *MongoStore) GetProject(slug string) (domain.Project, bool) {
	ctx, cancel := s.context()
	defer cancel()

	var project domain.Project
	if err := s.db.Collection("projects").FindOne(ctx, bson.M{"slug": slug}).Decode(&project); err != nil {
		return domain.Project{}, false
	}
	return project, true
}

func (s *MongoStore) CreateProject(project domain.Project) domain.Project {
	ctx, cancel := s.context()
	defer cancel()

	project.Slug = s.uniqueSlug(ctx, "projects", project.Slug, project.Title)
	_, _ = s.db.Collection("projects").InsertOne(ctx, project)
	return project
}

func (s *MongoStore) UpdateProject(slug string, project domain.Project) (domain.Project, error) {
	ctx, cancel := s.context()
	defer cancel()

	project.Slug = slug
	result, err := s.db.Collection("projects").ReplaceOne(ctx, bson.M{"slug": slug}, project)
	if err != nil {
		return domain.Project{}, err
	}
	if result.MatchedCount == 0 {
		return domain.Project{}, ErrNotFound
	}
	return project, nil
}

func (s *MongoStore) DeleteProject(slug string) bool {
	ctx, cancel := s.context()
	defer cancel()

	result, err := s.db.Collection("projects").DeleteOne(ctx, bson.M{"slug": slug})
	return err == nil && result.DeletedCount > 0
}

func (s *MongoStore) ListPosts() []domain.Post {
	ctx, cancel := s.context()
	defer cancel()

	cursor, err := s.db.Collection("posts").Find(ctx, bson.M{})
	if err != nil {
		return append([]domain.Post(nil), s.seed.Posts...)
	}
	defer cursor.Close(ctx)

	posts := make([]domain.Post, 0)
	if err := cursor.All(ctx, &posts); err != nil {
		return append([]domain.Post(nil), s.seed.Posts...)
	}
	return posts
}

func (s *MongoStore) GetPost(slug string) (domain.Post, bool) {
	ctx, cancel := s.context()
	defer cancel()

	var post domain.Post
	if err := s.db.Collection("posts").FindOne(ctx, bson.M{"slug": slug}).Decode(&post); err != nil {
		return domain.Post{}, false
	}
	return post, true
}

func (s *MongoStore) CreatePost(post domain.Post) domain.Post {
	ctx, cancel := s.context()
	defer cancel()

	post.Slug = s.uniqueSlug(ctx, "posts", post.Slug, post.Title)
	_, _ = s.db.Collection("posts").InsertOne(ctx, post)
	return post
}

func (s *MongoStore) UpdatePost(slug string, post domain.Post) (domain.Post, error) {
	ctx, cancel := s.context()
	defer cancel()

	post.Slug = slug
	result, err := s.db.Collection("posts").ReplaceOne(ctx, bson.M{"slug": slug}, post)
	if err != nil {
		return domain.Post{}, err
	}
	if result.MatchedCount == 0 {
		return domain.Post{}, ErrNotFound
	}
	return post, nil
}

func (s *MongoStore) DeletePost(slug string) bool {
	ctx, cancel := s.context()
	defer cancel()

	result, err := s.db.Collection("posts").DeleteOne(ctx, bson.M{"slug": slug})
	return err == nil && result.DeletedCount > 0
}

func (s *MongoStore) seedIfNeeded(ctx context.Context) error {
	if err := s.upsertDocument(ctx, "site", "site", siteDocument{ID: "site", Site: s.seed.Site}); err != nil {
		return err
	}
	if err := s.upsertDocument(ctx, "globals", "stats", statsDocument{ID: "stats", Stats: append([]domain.Stat(nil), s.seed.Stats...)}); err != nil {
		return err
	}
	if err := s.upsertDocument(ctx, "globals", "skills", skillsDocument{ID: "skills", Skills: append([]string(nil), s.seed.Skills...)}); err != nil {
		return err
	}

	projectCollection := s.db.Collection("projects")
	projectCount, err := projectCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return err
	}
	if projectCount == 0 {
		documents := make([]any, 0, len(s.seed.Projects))
		for _, project := range s.seed.Projects {
			documents = append(documents, project)
		}
		if len(documents) > 0 {
			_, err = projectCollection.InsertMany(ctx, documents)
			if err != nil {
				return err
			}
		}
	}

	postCollection := s.db.Collection("posts")
	postCount, err := postCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return err
	}
	if postCount == 0 {
		documents := make([]any, 0, len(s.seed.Posts))
		for _, post := range s.seed.Posts {
			documents = append(documents, post)
		}
		if len(documents) > 0 {
			_, err = postCollection.InsertMany(ctx, documents)
			if err != nil {
				return err
			}
		}
	}

	// Ensure unique indexes on slug fields for projects and posts
	if err := s.ensureIndexes(ctx); err != nil {
		return err
	}

	return nil
}

func (s *MongoStore) upsertDocument(ctx context.Context, collectionName string, id string, payload any) error {
	_, err := s.db.Collection(collectionName).ReplaceOne(ctx, bson.M{"_id": id}, payload, options.Replace().SetUpsert(true))
	return err
}

func (s *MongoStore) ensureIndexes(ctx context.Context) error {
	// projects.slug unique index
	projColl := s.db.Collection("projects")
	_, err := projColl.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "slug", Value: 1}},
		Options: options.Index().SetUnique(true).SetBackground(true),
	})
	if err != nil && !isIndexExistsError(err) {
		return err
	}

	// posts.slug unique index
	postColl := s.db.Collection("posts")
	_, err = postColl.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "slug", Value: 1}},
		Options: options.Index().SetUnique(true).SetBackground(true),
	})
	if err != nil && !isIndexExistsError(err) {
		return err
	}

	return nil
}

func isIndexExistsError(err error) bool {
	// Many drivers return a generic error when index exists; be permissive and treat as non-fatal.
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "already exists") || strings.Contains(err.Error(), "E11000")
}

func (s *MongoStore) uniqueSlug(ctx context.Context, collectionName string, currentSlug string, title string) string {
	baseSlug := strings.TrimSpace(currentSlug)
	if baseSlug == "" {
		baseSlug = slugify(title)
	}
	if baseSlug == "" {
		baseSlug = "item"
	}

	if !s.slugExists(ctx, collectionName, baseSlug) {
		return baseSlug
	}

	for suffix := 2; suffix < 1000; suffix++ {
		candidate := fmt.Sprintf("%s-%d", baseSlug, suffix)
		if !s.slugExists(ctx, collectionName, candidate) {
			return candidate
		}
	}

	return baseSlug
}

func (s *MongoStore) slugExists(ctx context.Context, collectionName string, slug string) bool {
	count, err := s.db.Collection(collectionName).CountDocuments(ctx, bson.M{"slug": slug})
	return err == nil && count > 0
}

func (s *MongoStore) SaveContactSubmission(submission domain.ContactSubmission) error {
	ctx, cancel := s.context()
	defer cancel()
	
	submission.CreatedAt = time.Now()
	_, err := s.db.Collection("contact_submissions").InsertOne(ctx, submission)
	return err
}

func (s *MongoStore) ListContactSubmissions() []domain.ContactSubmission {
	ctx, cancel := s.context()
	defer cancel()
	
	cursor, err := s.db.Collection("contact_submissions").Find(ctx, bson.M{})
	if err != nil {
		return []domain.ContactSubmission{}
	}
	defer cursor.Close(ctx)
	
	var submissions []domain.ContactSubmission
	if err := cursor.All(ctx, &submissions); err != nil {
		return []domain.ContactSubmission{}
	}
	return submissions
}

func (s *MongoStore) context() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), s.timeout)
}

func (s *MongoStore) Close(ctx context.Context) error {
	if s.client == nil {
		return nil
	}
	return s.client.Disconnect(ctx)
}

type siteDocument struct {
	ID   string      `bson:"_id" json:"-"`
	Site domain.Site `bson:"site" json:"site"`
}

type statsDocument struct {
	ID    string       `bson:"_id" json:"-"`
	Stats []domain.Stat `bson:"stats" json:"stats"`
}

type skillsDocument struct {
	ID     string   `bson:"_id" json:"-"`
	Skills []string `bson:"skills" json:"skills"`
}

var _ Repository = (*MongoStore)(nil)