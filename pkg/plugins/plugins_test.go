package plugins

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

// TestParseFrontmatter tests the frontmatter parsing functionality.
func TestParseFrontmatter(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		wantBody  string
		wantErr   bool
		checkMeta func(t *testing.T, meta map[string]interface{})
	}{
		{
			name: "basic YAML frontmatter with title and date",
			content: `---
title: My Post
date: 2024-01-15
---
# Content here`,
			wantBody: "# Content here",
			wantErr:  false,
			checkMeta: func(t *testing.T, meta map[string]interface{}) {
				if meta["title"] != "My Post" {
					t.Errorf("title = %v, want 'My Post'", meta["title"])
				}
				// Date is parsed as time.Time by YAML library
				if _, ok := meta["date"].(time.Time); !ok {
					t.Errorf("date should be time.Time, got %T", meta["date"])
				}
			},
		},
		{
			name: "frontmatter with tags list",
			content: `---
title: Tagged Post
tags:
  - go
  - programming
  - tutorial
---
Body content`,
			wantBody: "Body content",
			wantErr:  false,
			checkMeta: func(t *testing.T, meta map[string]interface{}) {
				if meta["title"] != "Tagged Post" {
					t.Errorf("title = %v, want 'Tagged Post'", meta["title"])
				}
				tags, ok := meta["tags"].([]interface{})
				if !ok {
					t.Errorf("tags should be []interface{}, got %T", meta["tags"])
					return
				}
				if len(tags) != 3 {
					t.Errorf("tags length = %d, want 3", len(tags))
				}
			},
		},
		{
			name: "frontmatter with nested objects",
			content: `---
title: Nested
author:
  name: John Doe
  email: john@example.com
---
Content`,
			wantBody: "Content",
			wantErr:  false,
			checkMeta: func(t *testing.T, meta map[string]interface{}) {
				if meta["title"] != "Nested" {
					t.Errorf("title = %v, want 'Nested'", meta["title"])
				}
				author, ok := meta["author"].(map[string]interface{})
				if !ok {
					t.Errorf("author should be map[string]interface{}, got %T", meta["author"])
					return
				}
				if author["name"] != "John Doe" {
					t.Errorf("author.name = %v, want 'John Doe'", author["name"])
				}
			},
		},
		{
			name:     "missing frontmatter",
			content:  "# Just a heading\n\nSome content",
			wantBody: "# Just a heading\n\nSome content",
			wantErr:  false,
			checkMeta: func(t *testing.T, meta map[string]interface{}) {
				if len(meta) != 0 {
					t.Errorf("meta should be empty, got %v", meta)
				}
			},
		},
		{
			name: "empty frontmatter",
			content: `---
---
Content after empty frontmatter`,
			wantBody: "Content after empty frontmatter",
			wantErr:  false,
			checkMeta: func(t *testing.T, meta map[string]interface{}) {
				if len(meta) != 0 {
					t.Errorf("meta should be empty, got %v", meta)
				}
			},
		},
		{
			name: "boolean values",
			content: `---
published: true
draft: false
---
Content`,
			wantBody: "Content",
			wantErr:  false,
			checkMeta: func(t *testing.T, meta map[string]interface{}) {
				if meta["published"] != true {
					t.Errorf("published = %v, want true", meta["published"])
				}
				if meta["draft"] != false {
					t.Errorf("draft = %v, want false", meta["draft"])
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMeta, gotBody, err := ParseFrontmatter(tt.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFrontmatter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if gotBody != tt.wantBody {
					t.Errorf("ParseFrontmatter() gotBody = %q, want %q", gotBody, tt.wantBody)
				}
				if tt.checkMeta != nil {
					tt.checkMeta(t, gotMeta)
				}
			}
		})
	}
}

// TestExtractFrontmatter tests the frontmatter extraction.
func TestExtractFrontmatter(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		wantFM   string
		wantBody string
		wantErr  bool
	}{
		{
			name: "standard frontmatter",
			content: `---
key: value
---
body`,
			wantFM:   "key: value",
			wantBody: "body",
			wantErr:  false,
		},
		{
			name:     "no frontmatter",
			content:  "just content",
			wantFM:   "",
			wantBody: "just content",
			wantErr:  false,
		},
		{
			name: "unclosed frontmatter",
			content: `---
key: value
no closing`,
			wantFM:   "",
			wantBody: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFM, gotBody, err := ExtractFrontmatter(tt.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractFrontmatter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if gotFM != tt.wantFM {
					t.Errorf("ExtractFrontmatter() gotFM = %q, want %q", gotFM, tt.wantFM)
				}
				if gotBody != tt.wantBody {
					t.Errorf("ExtractFrontmatter() gotBody = %q, want %q", gotBody, tt.wantBody)
				}
			}
		})
	}
}

// TestGlobPlugin tests the glob plugin functionality.
func TestGlobPlugin(t *testing.T) {
	// Create a temporary directory structure for testing
	tmpDir := t.TempDir()

	// Create test files
	testFiles := []string{
		"posts/post1.md",
		"posts/post2.md",
		"posts/nested/post3.md",
		"pages/about.md",
		"README.md",
		"ignored/draft.md",
	}

	for _, f := range testFiles {
		path := filepath.Join(tmpDir, f)
		_ = os.MkdirAll(filepath.Dir(path), 0o755) //nolint:errcheck // test setup
		//nolint:gosec,errcheck // test file
		_ = os.WriteFile(path, []byte("# Test"), 0o644)
	}

	// Create a .gitignore
	gitignore := `ignored/
*.txt
`
	//nolint:gosec,errcheck // test setup
	_ = os.WriteFile(filepath.Join(tmpDir, ".gitignore"), []byte(gitignore), 0o644)

	tests := []struct {
		name         string
		patterns     []string
		useGitignore bool
		wantCount    int
	}{
		{
			name:         "all markdown files",
			patterns:     []string{"**/*.md"},
			useGitignore: true,
			wantCount:    5, // Excludes ignored/draft.md
		},
		{
			name:         "posts only",
			patterns:     []string{"posts/**/*.md"},
			useGitignore: false,
			wantCount:    3,
		},
		{
			name:         "multiple patterns",
			patterns:     []string{"posts/*.md", "pages/*.md"},
			useGitignore: false,
			wantCount:    3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plugin := NewGlobPlugin()
			plugin.SetPatterns(tt.patterns)
			plugin.SetUseGitignore(tt.useGitignore)

			m := lifecycle.NewManager()
			cfg := m.Config()
			cfg.ContentDir = tmpDir
			cfg.GlobPatterns = tt.patterns

			// Configure the plugin
			if err := plugin.Configure(m); err != nil {
				t.Fatalf("Configure() error = %v", err)
			}

			// Run glob
			if err := plugin.Glob(m); err != nil {
				t.Fatalf("Glob() error = %v", err)
			}

			files := m.Files()
			if len(files) != tt.wantCount {
				t.Errorf("Glob() found %d files, want %d. Files: %v", len(files), tt.wantCount, files)
			}
		})
	}
}

// TestLoadPlugin tests the load plugin functionality.
func TestLoadPlugin(t *testing.T) {
	// Create a temporary directory with test files
	tmpDir := t.TempDir()

	// Create test markdown files
	files := map[string]string{
		"post1.md": `---
title: First Post
date: 2024-01-15
published: true
tags:
  - go
  - test
---
# First Post Content

This is the body.`,
		"post2.md": `---
title: Draft Post
draft: true
slug: custom-slug
---
Draft content here.`,
		"post3.md": `No frontmatter here.

Just content.`,
		"post4.md": `---
title: Extra Fields
custom_field: custom_value
nested:
  key: value
---
Content with extras.`,
	}

	for name, content := range files {
		//nolint:gosec,errcheck // test setup
		_ = os.WriteFile(filepath.Join(tmpDir, name), []byte(content), 0o644)
	}

	// Setup manager with files
	m := lifecycle.NewManager()
	cfg := m.Config()
	cfg.ContentDir = tmpDir
	m.SetFiles([]string{"post1.md", "post2.md", "post3.md", "post4.md"})

	// Run load plugin
	plugin := NewLoadPlugin()
	if err := plugin.Load(m); err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	posts := m.Posts()
	if len(posts) != 4 {
		t.Fatalf("Load() created %d posts, want 4", len(posts))
	}

	// Test post1
	post1 := findPostByPath(posts, "post1.md")
	if post1 == nil {
		t.Fatal("post1.md not found")
	}
	if post1.Title == nil || *post1.Title != "First Post" {
		t.Errorf("post1 title = %v, want 'First Post'", post1.Title)
	}
	if !post1.Published {
		t.Error("post1 should be published")
	}
	if len(post1.Tags) != 2 {
		t.Errorf("post1 has %d tags, want 2", len(post1.Tags))
	}
	expectedDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	if post1.Date == nil || !post1.Date.Equal(expectedDate) {
		t.Errorf("post1 date = %v, want %v", post1.Date, expectedDate)
	}

	// Test post2 (custom slug)
	post2 := findPostByPath(posts, "post2.md")
	if post2 == nil {
		t.Fatal("post2.md not found")
	}
	if post2.Slug != "custom-slug" {
		t.Errorf("post2 slug = %q, want 'custom-slug'", post2.Slug)
	}
	if !post2.Draft {
		t.Error("post2 should be a draft")
	}

	// Test post3 (no frontmatter)
	post3 := findPostByPath(posts, "post3.md")
	if post3 == nil {
		t.Fatal("post3.md not found")
	}
	if post3.Title != nil {
		t.Errorf("post3 should have no title, got %v", post3.Title)
	}
	if post3.Content != "No frontmatter here.\n\nJust content." {
		t.Errorf("post3 content mismatch: %q", post3.Content)
	}

	// Test post4 (extra fields)
	post4 := findPostByPath(posts, "post4.md")
	if post4 == nil {
		t.Fatal("post4.md not found")
	}
	if post4.Get("custom_field") != "custom_value" {
		t.Errorf("post4 custom_field = %v, want 'custom_value'", post4.Get("custom_field"))
	}
	nested, ok := post4.Get("nested").(map[string]interface{})
	if !ok {
		t.Error("post4 nested should be a map")
	} else if nested["key"] != "value" {
		t.Errorf("post4 nested.key = %v, want 'value'", nested["key"])
	}
}

// findPostByPath finds a post by its path in a slice of posts.
func findPostByPath(posts []*models.Post, path string) *models.Post {
	for _, p := range posts {
		if p.Path == path {
			return p
		}
	}
	return nil
}

// TestGetBool tests the boolean extraction helper.
func TestGetBool(t *testing.T) {
	tests := []struct {
		name       string
		metadata   map[string]interface{}
		key        string
		defaultVal bool
		want       bool
	}{
		{
			name:       "true boolean",
			metadata:   map[string]interface{}{"key": true},
			key:        "key",
			defaultVal: false,
			want:       true,
		},
		{
			name:       "false boolean",
			metadata:   map[string]interface{}{"key": false},
			key:        "key",
			defaultVal: true,
			want:       false,
		},
		{
			name:       "yes string",
			metadata:   map[string]interface{}{"key": "yes"},
			key:        "key",
			defaultVal: false,
			want:       true,
		},
		{
			name:       "no string",
			metadata:   map[string]interface{}{"key": "no"},
			key:        "key",
			defaultVal: true,
			want:       false,
		},
		{
			name:       "missing key",
			metadata:   map[string]interface{}{},
			key:        "missing",
			defaultVal: true,
			want:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetBool(tt.metadata, tt.key, tt.defaultVal)
			if got != tt.want {
				t.Errorf("GetBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestGetStringSlice tests the string slice extraction helper.
func TestGetStringSlice(t *testing.T) {
	tests := []struct {
		name     string
		metadata map[string]interface{}
		key      string
		want     []string
	}{
		{
			name:     "string slice from interface slice",
			metadata: map[string]interface{}{"tags": []interface{}{"a", "b", "c"}},
			key:      "tags",
			want:     []string{"a", "b", "c"},
		},
		{
			name:     "direct string slice",
			metadata: map[string]interface{}{"tags": []string{"x", "y"}},
			key:      "tags",
			want:     []string{"x", "y"},
		},
		{
			name:     "missing key",
			metadata: map[string]interface{}{},
			key:      "missing",
			want:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetStringSlice(tt.metadata, tt.key)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetStringSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}
