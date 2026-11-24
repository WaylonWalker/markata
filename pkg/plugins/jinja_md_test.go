package plugins

import (
	"testing"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

func TestJinjaMdPlugin_Name(t *testing.T) {
	p := NewJinjaMdPlugin()
	if got := p.Name(); got != "jinja_md" {
		t.Errorf("Name() = %q, want %q", got, "jinja_md")
	}
}

func TestJinjaMdPlugin_Configure(t *testing.T) {
	p := NewJinjaMdPlugin()
	m := lifecycle.NewManager()

	err := p.Configure(m)
	if err != nil {
		t.Errorf("Configure() error = %v", err)
	}

	if p.engine == nil {
		t.Error("Configure() did not initialize engine")
	}
}

func TestJinjaMdPlugin_Transform_JinjaEnabled(t *testing.T) {
	p := NewJinjaMdPlugin()
	m := lifecycle.NewManager()

	err := p.Configure(m)
	if err != nil {
		t.Fatalf("Configure() error = %v", err)
	}

	// Create a post with jinja enabled
	title := "Test Post"
	post := &models.Post{
		Title:   &title,
		Content: "Title: {{ post.title }}",
		Extra:   map[string]interface{}{"jinja": true},
	}
	m.AddPost(post)

	// Transform
	err = p.Transform(m)
	if err != nil {
		t.Fatalf("Transform() error = %v", err)
	}

	// Check content was processed
	expected := "Title: Test Post"
	if post.Content != expected {
		t.Errorf("Transform() Content = %q, want %q", post.Content, expected)
	}
}

func TestJinjaMdPlugin_Transform_JinjaDisabled(t *testing.T) {
	p := NewJinjaMdPlugin()
	m := lifecycle.NewManager()

	err := p.Configure(m)
	if err != nil {
		t.Fatalf("Configure() error = %v", err)
	}

	// Create a post without jinja enabled
	title := "Test Post"
	originalContent := "Title: {{ post.title }}"
	post := &models.Post{
		Title:   &title,
		Content: originalContent,
	}
	m.AddPost(post)

	// Transform
	err = p.Transform(m)
	if err != nil {
		t.Fatalf("Transform() error = %v", err)
	}

	// Content should remain unchanged
	if post.Content != originalContent {
		t.Errorf("Transform() modified content when jinja disabled: %q", post.Content)
	}
}

func TestJinjaMdPlugin_Transform_ForLoop(t *testing.T) {
	p := NewJinjaMdPlugin()
	m := lifecycle.NewManager()

	err := p.Configure(m)
	if err != nil {
		t.Fatalf("Configure() error = %v", err)
	}

	// Create a post with tags loop
	post := &models.Post{
		Tags:    []string{"go", "templates"},
		Content: "Tags: {% for tag in post.tags %}{{ tag }} {% endfor %}",
		Extra:   map[string]interface{}{"jinja": true},
	}
	m.AddPost(post)

	// Transform
	err = p.Transform(m)
	if err != nil {
		t.Fatalf("Transform() error = %v", err)
	}

	expected := "Tags: go templates "
	if post.Content != expected {
		t.Errorf("Transform() Content = %q, want %q", post.Content, expected)
	}
}

func TestJinjaMdPlugin_Transform_Conditional(t *testing.T) {
	p := NewJinjaMdPlugin()
	m := lifecycle.NewManager()

	err := p.Configure(m)
	if err != nil {
		t.Fatalf("Configure() error = %v", err)
	}

	// Test with published=true
	post := &models.Post{
		Published: true,
		Content:   "{% if post.published %}Published{% else %}Draft{% endif %}",
		Extra:     map[string]interface{}{"jinja": true},
	}
	m.AddPost(post)

	err = p.Transform(m)
	if err != nil {
		t.Fatalf("Transform() error = %v", err)
	}

	if post.Content != "Published" {
		t.Errorf("Transform() Content = %q, want %q", post.Content, "Published")
	}

	// Test with published=false
	m.SetPosts([]*models.Post{})
	post2 := &models.Post{
		Published: false,
		Content:   "{% if post.published %}Published{% else %}Draft{% endif %}",
		Extra:     map[string]interface{}{"jinja": true},
	}
	m.AddPost(post2)

	err = p.Transform(m)
	if err != nil {
		t.Fatalf("Transform() error = %v", err)
	}

	if post2.Content != "Draft" {
		t.Errorf("Transform() Content = %q, want %q", post2.Content, "Draft")
	}
}

func TestJinjaMdPlugin_Transform_Filters(t *testing.T) {
	p := NewJinjaMdPlugin()
	m := lifecycle.NewManager()

	err := p.Configure(m)
	if err != nil {
		t.Fatalf("Configure() error = %v", err)
	}

	title := "hello world"
	post := &models.Post{
		Title:   &title,
		Content: "{{ post.title | upper }}",
		Extra:   map[string]interface{}{"jinja": true},
	}
	m.AddPost(post)

	err = p.Transform(m)
	if err != nil {
		t.Fatalf("Transform() error = %v", err)
	}

	expected := "HELLO WORLD"
	if post.Content != expected {
		t.Errorf("Transform() Content = %q, want %q", post.Content, expected)
	}
}

func TestJinjaMdPlugin_Transform_ConfigAccess(t *testing.T) {
	p := NewJinjaMdPlugin()
	m := lifecycle.NewManager()

	// Set config values
	config := m.Config()
	config.Extra["title"] = "My Site"

	err := p.Configure(m)
	if err != nil {
		t.Fatalf("Configure() error = %v", err)
	}

	post := &models.Post{
		Content: "Site: {{ config.title }}",
		Extra:   map[string]interface{}{"jinja": true},
	}
	m.AddPost(post)

	err = p.Transform(m)
	if err != nil {
		t.Fatalf("Transform() error = %v", err)
	}

	expected := "Site: My Site"
	if post.Content != expected {
		t.Errorf("Transform() Content = %q, want %q", post.Content, expected)
	}
}

func TestJinjaMdPlugin_Transform_SkippedPost(t *testing.T) {
	p := NewJinjaMdPlugin()
	m := lifecycle.NewManager()

	err := p.Configure(m)
	if err != nil {
		t.Fatalf("Configure() error = %v", err)
	}

	originalContent := "{{ post.title }}"
	post := &models.Post{
		Skip:    true,
		Content: originalContent,
		Extra:   map[string]interface{}{"jinja": true},
	}
	m.AddPost(post)

	err = p.Transform(m)
	if err != nil {
		t.Fatalf("Transform() error = %v", err)
	}

	// Content should remain unchanged for skipped posts
	if post.Content != originalContent {
		t.Errorf("Transform() modified skipped post content")
	}
}

func TestJinjaMdPlugin_Priority(t *testing.T) {
	p := NewJinjaMdPlugin()

	// Should run early in transform stage
	transformPriority := p.Priority(lifecycle.StageTransform)
	if transformPriority != lifecycle.PriorityEarly {
		t.Errorf("Priority(StageTransform) = %d, want %d", transformPriority, lifecycle.PriorityEarly)
	}

	// Default priority for other stages
	otherPriority := p.Priority(lifecycle.StageRender)
	if otherPriority != lifecycle.PriorityDefault {
		t.Errorf("Priority(StageRender) = %d, want %d", otherPriority, lifecycle.PriorityDefault)
	}
}

func TestCollectPrivatePathsForJinja(t *testing.T) {
	tests := []struct {
		name     string
		posts    []*models.Post
		expected []string
	}{
		{
			name:     "empty posts",
			posts:    []*models.Post{},
			expected: nil,
		},
		{
			name: "no private posts",
			posts: []*models.Post{
				{Slug: "public", Href: "/public/", Private: false},
			},
			expected: nil,
		},
		{
			name: "private post includes all variants",
			posts: []*models.Post{
				{Slug: "secret", Href: "/secret/", Private: true},
			},
			expected: []string{
				"/secret/",
				"/secret.txt",
				"/secret.ansi",
				"/secret.md",
				"/secret.og/",
			},
		},
		{
			name: "excludes robots post",
			posts: []*models.Post{
				{Slug: "robots", Href: "/robots/", Private: true},
				{Slug: "secret", Href: "/secret/", Private: true},
			},
			expected: []string{
				"/secret/",
				"/secret.txt",
				"/secret.ansi",
				"/secret.md",
				"/secret.og/",
			},
		},
		{
			name: "excludes drafts and skipped",
			posts: []*models.Post{
				{Slug: "draft-post", Href: "/draft-post/", Private: true, Draft: true},
				{Slug: "skipped-post", Href: "/skipped-post/", Private: true, Skip: true},
				{Slug: "real-private", Href: "/real-private/", Private: true},
			},
			expected: []string{
				"/real-private/",
				"/real-private.txt",
				"/real-private.ansi",
				"/real-private.md",
				"/real-private.og/",
			},
		},
		{
			name: "multiple private posts",
			posts: []*models.Post{
				{Slug: "private1", Href: "/private1/", Private: true},
				{Slug: "public", Href: "/public/", Private: false},
				{Slug: "private2", Href: "/private2/", Private: true},
			},
			expected: []string{
				"/private1/",
				"/private1.txt",
				"/private1.ansi",
				"/private1.md",
				"/private1.og/",
				"/private2/",
				"/private2.txt",
				"/private2.ansi",
				"/private2.md",
				"/private2.og/",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := collectPrivatePathsForJinja(tt.posts)
			if len(got) != len(tt.expected) {
				t.Errorf("collectPrivatePathsForJinja() returned %d paths, want %d\ngot: %v\nwant: %v",
					len(got), len(tt.expected), got, tt.expected)
				return
			}
			for i, path := range got {
				if path != tt.expected[i] {
					t.Errorf("collectPrivatePathsForJinja()[%d] = %q, want %q", i, path, tt.expected[i])
				}
			}
		})
	}
}

func TestIsJinjaEnabled(t *testing.T) {
	tests := []struct {
		name     string
		extra    map[string]interface{}
		expected bool
	}{
		{
			name:     "nil extra",
			extra:    nil,
			expected: false,
		},
		{
			name:     "jinja not set",
			extra:    map[string]interface{}{},
			expected: false,
		},
		{
			name:     "jinja true (bool)",
			extra:    map[string]interface{}{"jinja": true},
			expected: true,
		},
		{
			name:     "jinja false (bool)",
			extra:    map[string]interface{}{"jinja": false},
			expected: false,
		},
		{
			name:     "jinja 'true' (string)",
			extra:    map[string]interface{}{"jinja": "true"},
			expected: true,
		},
		{
			name:     "jinja 'yes' (string)",
			extra:    map[string]interface{}{"jinja": "yes"},
			expected: true,
		},
		{
			name:     "jinja 'false' (string)",
			extra:    map[string]interface{}{"jinja": "false"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			post := &models.Post{Extra: tt.extra}
			if got := isJinjaEnabled(post); got != tt.expected {
				t.Errorf("isJinjaEnabled() = %v, want %v", got, tt.expected)
			}
		})
	}
}
