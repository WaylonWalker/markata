package plugins

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

func TestStructuredDataPlugin_Name(t *testing.T) {
	plugin := NewStructuredDataPlugin()
	if got := plugin.Name(); got != "structured_data" {
		t.Errorf("Name() = %v, want %v", got, "structured_data")
	}
}

func TestStructuredDataPlugin_Priority(t *testing.T) {
	plugin := NewStructuredDataPlugin()

	// Should return default priority for transform stage
	if got := plugin.Priority(lifecycle.StageTransform); got != lifecycle.PriorityDefault {
		t.Errorf("Priority(StageTransform) = %v, want %v", got, lifecycle.PriorityDefault)
	}

	// Should return default priority for other stages
	if got := plugin.Priority(lifecycle.StageWrite); got != lifecycle.PriorityDefault {
		t.Errorf("Priority(StageWrite) = %v, want %v", got, lifecycle.PriorityDefault)
	}
}

func TestStructuredDataPlugin_Transform(t *testing.T) {
	plugin := NewStructuredDataPlugin()

	// Create a manager with a test post
	m := lifecycle.NewManager()
	m.SetConfig(&lifecycle.Config{
		OutputDir: "output",
		Extra: map[string]interface{}{
			"url":         "https://example.com",
			"title":       "Test Site",
			"description": "A test site",
			"author":      "Test Author",
			"seo":         models.NewSEOConfig(),
		},
	})

	title := "Test Post Title"
	description := "Test post description"
	postDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	post := &models.Post{
		Path:        "test.md",
		Title:       &title,
		Description: &description,
		Date:        &postDate,
		Slug:        "test-post",
		Href:        "/test-post/",
		Tags:        []string{"go", "testing"},
		Extra:       make(map[string]interface{}),
	}
	m.SetPosts([]*models.Post{post})

	// Run the transform
	err := plugin.Transform(m)
	if err != nil {
		t.Fatalf("Transform() error = %v", err)
	}

	// Check that structured data was added
	sd, ok := post.Extra["structured_data"].(*models.StructuredData)
	if !ok {
		t.Fatal("structured_data not found in post.Extra")
	}

	// Verify JSON-LD
	if sd.JSONLD == "" {
		t.Error("JSON-LD should not be empty")
	}

	// Parse JSON-LD to verify structure
	var jsonLD map[string]interface{}
	if err := json.Unmarshal([]byte(sd.JSONLD), &jsonLD); err != nil {
		t.Errorf("JSON-LD should be valid JSON: %v", err)
	}

	if jsonLD["@context"] != "https://schema.org" {
		t.Errorf("JSON-LD @context = %v, want https://schema.org", jsonLD["@context"])
	}
	if jsonLD["@type"] != "BlogPosting" {
		t.Errorf("JSON-LD @type = %v, want BlogPosting", jsonLD["@type"])
	}
	if jsonLD["headline"] != "Test Post Title" {
		t.Errorf("JSON-LD headline = %v, want Test Post Title", jsonLD["headline"])
	}

	// Verify OpenGraph tags
	if len(sd.OpenGraph) == 0 {
		t.Error("OpenGraph tags should not be empty")
	}

	// Check for required OG tags
	ogTags := make(map[string]string)
	for _, tag := range sd.OpenGraph {
		ogTags[tag.Property] = tag.Content
	}

	if ogTags["og:title"] != "Test Post Title" {
		t.Errorf("og:title = %v, want Test Post Title", ogTags["og:title"])
	}
	if ogTags["og:type"] != "article" {
		t.Errorf("og:type = %v, want article", ogTags["og:type"])
	}
	if ogTags["og:url"] != "https://example.com/test-post/" {
		t.Errorf("og:url = %v, want https://example.com/test-post/", ogTags["og:url"])
	}

	// Verify Twitter Card tags
	if len(sd.Twitter) == 0 {
		t.Error("Twitter tags should not be empty")
	}

	twitterTags := make(map[string]string)
	for _, tag := range sd.Twitter {
		twitterTags[tag.Name] = tag.Content
	}

	if twitterTags["twitter:card"] != "summary" {
		t.Errorf("twitter:card = %v, want summary", twitterTags["twitter:card"])
	}
	if twitterTags["twitter:title"] != "Test Post Title" {
		t.Errorf("twitter:title = %v, want Test Post Title", twitterTags["twitter:title"])
	}
}

func TestStructuredDataPlugin_SkipDraft(t *testing.T) {
	plugin := NewStructuredDataPlugin()

	m := lifecycle.NewManager()
	m.SetConfig(&lifecycle.Config{
		OutputDir: "output",
		Extra: map[string]interface{}{
			"url":   "https://example.com",
			"title": "Test Site",
			"seo":   models.NewSEOConfig(),
		},
	})

	title := "Draft Post"
	post := &models.Post{
		Path:  "draft.md",
		Title: &title,
		Slug:  "draft",
		Href:  "/draft/",
		Draft: true,
		Extra: make(map[string]interface{}),
	}
	m.SetPosts([]*models.Post{post})

	err := plugin.Transform(m)
	if err != nil {
		t.Fatalf("Transform() error = %v", err)
	}

	// Draft posts should not have structured data
	if _, ok := post.Extra["structured_data"]; ok {
		t.Error("Draft posts should not have structured data")
	}
}

func TestStructuredDataPlugin_SkipNoTitle(t *testing.T) {
	plugin := NewStructuredDataPlugin()

	m := lifecycle.NewManager()
	m.SetConfig(&lifecycle.Config{
		OutputDir: "output",
		Extra: map[string]interface{}{
			"url":   "https://example.com",
			"title": "Test Site",
			"seo":   models.NewSEOConfig(),
		},
	})

	post := &models.Post{
		Path:  "no-title.md",
		Title: nil, // No title
		Slug:  "no-title",
		Href:  "/no-title/",
		Extra: make(map[string]interface{}),
	}
	m.SetPosts([]*models.Post{post})

	err := plugin.Transform(m)
	if err != nil {
		t.Fatalf("Transform() error = %v", err)
	}

	// Posts without titles should not have structured data
	if _, ok := post.Extra["structured_data"]; ok {
		t.Error("Posts without titles should not have structured data")
	}
}

func TestStructuredDataPlugin_WithImage(t *testing.T) {
	plugin := NewStructuredDataPlugin()

	m := lifecycle.NewManager()
	m.SetConfig(&lifecycle.Config{
		OutputDir: "output",
		Extra: map[string]interface{}{
			"url":   "https://example.com",
			"title": "Test Site",
			"seo":   models.NewSEOConfig(),
		},
	})

	title := "Post With Image"
	post := &models.Post{
		Path:  "image-post.md",
		Title: &title,
		Slug:  "image-post",
		Href:  "/image-post/",
		Extra: map[string]interface{}{
			"image": "/images/post.jpg",
		},
	}
	m.SetPosts([]*models.Post{post})

	err := plugin.Transform(m)
	if err != nil {
		t.Fatalf("Transform() error = %v", err)
	}

	sd, ok := post.Extra["structured_data"].(*models.StructuredData)
	if !ok {
		t.Fatal("structured_data is not *models.StructuredData")
	}

	// Check OG image
	ogTags := make(map[string]string)
	for _, tag := range sd.OpenGraph {
		ogTags[tag.Property] = tag.Content
	}

	if !strings.Contains(ogTags["og:image"], "/images/post.jpg") {
		t.Errorf("og:image should contain /images/post.jpg, got %v", ogTags["og:image"])
	}

	// Twitter card should be summary_large_image when image is present
	twitterTags := make(map[string]string)
	for _, tag := range sd.Twitter {
		twitterTags[tag.Name] = tag.Content
	}

	if twitterTags["twitter:card"] != "summary_large_image" {
		t.Errorf("twitter:card = %v, want summary_large_image", twitterTags["twitter:card"])
	}
}

func TestStructuredDataPlugin_Disabled(t *testing.T) {
	plugin := NewStructuredDataPlugin()

	// Create SEO config with structured data disabled
	enabled := false
	seoConfig := models.SEOConfig{
		StructuredData: models.StructuredDataConfig{
			Enabled: &enabled,
		},
	}

	m := lifecycle.NewManager()
	m.SetConfig(&lifecycle.Config{
		OutputDir: "output",
		Extra: map[string]interface{}{
			"url":   "https://example.com",
			"title": "Test Site",
			"seo":   seoConfig,
		},
	})

	title := "Test Post"
	post := &models.Post{
		Path:  "test.md",
		Title: &title,
		Slug:  "test",
		Href:  "/test/",
		Extra: make(map[string]interface{}),
	}
	m.SetPosts([]*models.Post{post})

	err := plugin.Transform(m)
	if err != nil {
		t.Fatalf("Transform() error = %v", err)
	}

	// When disabled, no structured data should be added
	if _, ok := post.Extra["structured_data"]; ok {
		t.Error("Structured data should not be added when disabled")
	}
}

func TestStructuredDataPlugin_WithTwitterHandle(t *testing.T) {
	plugin := NewStructuredDataPlugin()

	seoConfig := models.SEOConfig{
		TwitterHandle:  "testhandle",
		StructuredData: models.NewStructuredDataConfig(),
	}

	m := lifecycle.NewManager()
	m.SetConfig(&lifecycle.Config{
		OutputDir: "output",
		Extra: map[string]interface{}{
			"url":   "https://example.com",
			"title": "Test Site",
			"seo":   seoConfig,
		},
	})

	title := "Test Post"
	post := &models.Post{
		Path:  "test.md",
		Title: &title,
		Slug:  "test",
		Href:  "/test/",
		Extra: make(map[string]interface{}),
	}
	m.SetPosts([]*models.Post{post})

	err := plugin.Transform(m)
	if err != nil {
		t.Fatalf("Transform() error = %v", err)
	}

	sd, ok := post.Extra["structured_data"].(*models.StructuredData)
	if !ok {
		t.Fatal("structured_data is not *models.StructuredData")
	}

	twitterTags := make(map[string]string)
	for _, tag := range sd.Twitter {
		twitterTags[tag.Name] = tag.Content
	}

	if twitterTags["twitter:site"] != "@testhandle" {
		t.Errorf("twitter:site = %v, want @testhandle", twitterTags["twitter:site"])
	}
}

func TestStructuredDataPlugin_MakeAbsoluteURL(t *testing.T) {
	plugin := NewStructuredDataPlugin()

	tests := []struct {
		name    string
		url     string
		siteURL string
		want    string
	}{
		{
			name:    "already absolute http",
			url:     "http://other.com/image.jpg",
			siteURL: "https://example.com",
			want:    "http://other.com/image.jpg",
		},
		{
			name:    "already absolute https",
			url:     "https://other.com/image.jpg",
			siteURL: "https://example.com",
			want:    "https://other.com/image.jpg",
		},
		{
			name:    "protocol-relative",
			url:     "//cdn.example.com/image.jpg",
			siteURL: "https://example.com",
			want:    "https://cdn.example.com/image.jpg",
		},
		{
			name:    "relative with leading slash",
			url:     "/images/post.jpg",
			siteURL: "https://example.com",
			want:    "https://example.com/images/post.jpg",
		},
		{
			name:    "relative without leading slash",
			url:     "images/post.jpg",
			siteURL: "https://example.com",
			want:    "https://example.com/images/post.jpg",
		},
		{
			name:    "empty url",
			url:     "",
			siteURL: "https://example.com",
			want:    "",
		},
		{
			name:    "site url with trailing slash",
			url:     "/images/post.jpg",
			siteURL: "https://example.com/",
			want:    "https://example.com/images/post.jpg",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := plugin.makeAbsoluteURL(tt.url, tt.siteURL)
			if got != tt.want {
				t.Errorf("makeAbsoluteURL(%q, %q) = %q, want %q", tt.url, tt.siteURL, got, tt.want)
			}
		})
	}
}

func TestStructuredDataPlugin_ArticleTags(t *testing.T) {
	plugin := NewStructuredDataPlugin()

	m := lifecycle.NewManager()
	m.SetConfig(&lifecycle.Config{
		OutputDir: "output",
		Extra: map[string]interface{}{
			"url":   "https://example.com",
			"title": "Test Site",
			"seo":   models.NewSEOConfig(),
		},
	})

	title := "Tagged Post"
	postDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	post := &models.Post{
		Path:  "tagged.md",
		Title: &title,
		Date:  &postDate,
		Slug:  "tagged",
		Href:  "/tagged/",
		Tags:  []string{"go", "testing", "seo"},
		Extra: make(map[string]interface{}),
	}
	m.SetPosts([]*models.Post{post})

	err := plugin.Transform(m)
	if err != nil {
		t.Fatalf("Transform() error = %v", err)
	}

	sd, ok := post.Extra["structured_data"].(*models.StructuredData)
	if !ok {
		t.Fatal("structured_data is not *models.StructuredData")
	}

	// Count article:tag entries
	tagCount := 0
	for _, tag := range sd.OpenGraph {
		if tag.Property == "article:tag" {
			tagCount++
		}
	}

	if tagCount != 3 {
		t.Errorf("Expected 3 article:tag entries, got %d", tagCount)
	}
}
