package domain

type Site struct {
	Name        string `json:"name" bson:"name"`
	Tagline     string `json:"tagline" bson:"tagline"`
	Description string `json:"description" bson:"description"`
	Accent      string `json:"accent" bson:"accent"`
}

type Project struct {
	Title            string   `json:"title" bson:"title"`
	Slug             string   `json:"slug" bson:"slug"`
	DescriptionShort  string   `json:"descriptionShort" bson:"descriptionShort"`
	DescriptionLong   string   `json:"descriptionLong" bson:"descriptionLong"`
	Technologies     []string `json:"technologies" bson:"technologies"`
	Status           string   `json:"status" bson:"status"`
	Featured         bool     `json:"featured" bson:"featured"`
	Year             string   `json:"year" bson:"year"`
	Image            string   `json:"image" bson:"image"`
	LiveURL          string   `json:"liveUrl,omitempty" bson:"liveUrl,omitempty"`
	GitHubURL        string   `json:"githubUrl,omitempty" bson:"githubUrl,omitempty"`
}

type Post struct {
	Title       string   `json:"title" bson:"title"`
	Slug        string   `json:"slug" bson:"slug"`
	Excerpt     string   `json:"excerpt" bson:"excerpt"`
	Category    string   `json:"category" bson:"category"`
	ReadTime    string   `json:"readTime" bson:"readTime"`
	PublishedAt string   `json:"publishedAt" bson:"publishedAt"`
	Tags        []string `json:"tags" bson:"tags"`
	Sections    []PostSection `json:"sections,omitempty" bson:"sections,omitempty"`
}

type PostSection struct {
	Heading    string   `json:"heading" bson:"heading"`
	Paragraphs []string `json:"paragraphs" bson:"paragraphs"`
	Bullets    []string `json:"bullets,omitempty" bson:"bullets,omitempty"`
}

type Stat struct {
	Label string `json:"label" bson:"label"`
	Value string `json:"value" bson:"value"`
}

type Data struct {
	Site     Site      `json:"site"`
	Projects []Project `json:"projects"`
	Posts    []Post    `json:"posts"`
	Stats    []Stat    `json:"stats"`
	Skills   []string  `json:"skills"`
}