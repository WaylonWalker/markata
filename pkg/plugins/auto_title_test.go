package plugins

import (
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

func TestAutoTitlePlugin_Name(t *testing.T) {
	plugin := NewAutoTitlePlugin()
	if plugin.Name() != "auto_title" {
		t.Errorf("expected name 'auto_title', got %q", plugin.Name())
	}
}

func TestAutoTitlePlugin_Priority(t *testing.T) {
	plugin := NewAutoTitlePlugin()

	// Should return PriorityFirst for Transform stage
	if got := plugin.Priority(lifecycle.StageTransform); got != lifecycle.PriorityFirst {
		t.Errorf("expected PriorityFirst for Transform, got %d", got)
	}

	// Should return PriorityDefault for other stages
	if got := plugin.Priority(lifecycle.StageRender); got != lifecycle.PriorityDefault {
		t.Errorf("expected PriorityDefault for Render, got %d", got)
	}
}

func TestAutoTitlePlugin_generateTitle(t *testing.T) {
	plugin := NewAutoTitlePlugin()

	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "hyphenated filename",
			path:     "posts/my-first-post.md",
			expected: "My First Post",
		},
		{
			name:     "underscored filename",
			path:     "posts/python_tips.md",
			expected: "Python Tips",
		},
		{
			name:     "mixed separators",
			path:     "posts/python_tips-and-tricks.md",
			expected: "Python Tips And Tricks",
		},
		{
			name:     "date prefixed - strips date",
			path:     "posts/2024-01-15-new-feature.md",
			expected: "New Feature",
		},
		{
			name:     "date prefixed with underscore",
			path:     "posts/2024-01-15_my-post.md",
			expected: "My Post",
		},
		{
			name:     "single word",
			path:     "readme.md",
			expected: "Readme",
		},
		{
			name:     "uppercase characters",
			path:     "posts/README.md",
			expected: "Readme",
		},
		{
			name:     "nested path",
			path:     "blog/tutorials/getting-started.md",
			expected: "Getting Started",
		},
		{
			name:     "special characters preserved",
			path:     "posts/c++_tutorial.md",
			expected: "C++ Tutorial",
		},
		{
			name:     "numbers in filename",
			path:     "posts/api-v2-reference.md",
			expected: "Api V2 Reference",
		},
		{
			name:     "empty stem after extension removal",
			path:     ".md",
			expected: "",
		},
		{
			name:     "date only filename returns empty",
			path:     "posts/2024-01-15.md",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := plugin.generateTitle(tt.path)
			if got != tt.expected {
				t.Errorf("generateTitle(%q) = %q, want %q", tt.path, got, tt.expected)
			}
		})
	}
}

func TestAutoTitlePlugin_Transform(t *testing.T) {
	plugin := NewAutoTitlePlugin()

	tests := []struct {
		name          string
		posts         []*models.Post
		expectedTitle map[string]string // path -> expected title
	}{
		{
			name: "generates title for post without title",
			posts: []*models.Post{
				{Path: "posts/my-first-post.md", Content: "Hello"},
			},
			expectedTitle: map[string]string{
				"posts/my-first-post.md": "My First Post",
			},
		},
		{
			name: "extracts H1 from content",
			posts: []*models.Post{
				{Path: "posts/my-post.md", Content: "# Awesome Title\n\nContent here"},
			},
			expectedTitle: map[string]string{
				"posts/my-post.md": "Awesome Title",
			},
		},
		{
			name: "skips post with existing title",
			posts: func() []*models.Post {
				title := "My Custom Title"
				return []*models.Post{
					{Path: "posts/my-post.md", Title: &title, Content: "Hello"},
				}
			}(),
			expectedTitle: map[string]string{
				"posts/my-post.md": "My Custom Title",
			},
		},
		{
			name: "skips post marked as skip - title remains nil",
			posts: []*models.Post{
				{Path: "posts/skipped-post.md", Skip: true, Content: "Hello"},
			},
			expectedTitle: map[string]string{
				"posts/skipped-post.md": "", // Should remain nil
			},
		},
		{
			name: "generates title for post with empty title",
			posts: func() []*models.Post {
				emptyTitle := ""
				return []*models.Post{
					{Path: "posts/empty-title.md", Title: &emptyTitle, Content: "Hello"},
				}
			}(),
			expectedTitle: map[string]string{
				"posts/empty-title.md": "Empty Title",
			},
		},
		{
			name: "index file uses directory name",
			posts: []*models.Post{
				{Path: "docs/index.md", Content: "Content without heading"},
			},
			expectedTitle: map[string]string{
				"docs/index.md": "Docs",
			},
		},
		{
			name: "date-prefixed filename strips date",
			posts: []*models.Post{
				{Path: "posts/2024-01-15-new-feature.md", Content: "No heading"},
			},
			expectedTitle: map[string]string{
				"posts/2024-01-15-new-feature.md": "New Feature",
			},
		},
		{
			name: "multiple posts mixed",
			posts: func() []*models.Post {
				customTitle := "Custom"
				return []*models.Post{
					{Path: "posts/auto-gen.md", Content: "Hello"},
					{Path: "posts/has-title.md", Title: &customTitle, Content: "World"},
					{Path: "posts/another-auto.md", Content: "!"},
				}
			}(),
			expectedTitle: map[string]string{
				"posts/auto-gen.md":     "Auto Gen",
				"posts/has-title.md":    "Custom",
				"posts/another-auto.md": "Another Auto",
			},
		},
		{
			name: "H1 takes priority over filename",
			posts: []*models.Post{
				{Path: "posts/wrong-name.md", Content: "# Correct Title\n\nBody"},
			},
			expectedTitle: map[string]string{
				"posts/wrong-name.md": "Correct Title",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a manager with the posts
			m := lifecycle.NewManager()
			for _, post := range tt.posts {
				m.AddPost(post)
			}

			// Run the transform
			err := plugin.Transform(m)
			if err != nil {
				t.Fatalf("Transform() error = %v", err)
			}

			// Check titles
			for _, post := range m.Posts() {
				expected := tt.expectedTitle[post.Path]
				var got string
				if post.Title != nil {
					got = *post.Title
				}
				if got != expected {
					t.Errorf("post %q: title = %q, want %q", post.Path, got, expected)
				}
			}
		})
	}
}

func TestAutoTitlePlugin_Transform_NoNilTitles(t *testing.T) {
	// This test specifically verifies that no post ever has a nil title
	// after the auto_title plugin runs (unless Skip is true)
	plugin := NewAutoTitlePlugin()

	posts := []*models.Post{
		{Path: "normal.md", Content: ""},
		{Path: ".md", Content: ""},                 // Edge case: no filename
		{Path: "posts/2024-01-01.md", Content: ""}, // Date only
		{Path: "index.md", Content: ""},            // Root index with no directory
		{Path: "", Content: ""},                    // Empty path
	}

	m := lifecycle.NewManager()
	for _, post := range posts {
		m.AddPost(post)
	}

	err := plugin.Transform(m)
	if err != nil {
		t.Fatalf("Transform() error = %v", err)
	}

	for _, post := range m.Posts() {
		if post.Title == nil {
			t.Errorf("post %q has nil title - all non-skipped posts should have a title", post.Path)
		}
		if post.Title != nil && *post.Title == "" {
			t.Errorf("post %q has empty title - should have a fallback title", post.Path)
		}
	}
}

func TestToTitleCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello world", "Hello World"},
		{"HELLO WORLD", "Hello World"},
		{"hElLo WoRlD", "Hello World"},
		{"hello", "Hello"},
		{"", ""},
		{"multiple   spaces", "Multiple Spaces"},
		{"c++ tutorial", "C++ Tutorial"},
		{"api v2 reference", "Api V2 Reference"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := toTitleCase(tt.input)
			if got != tt.expected {
				t.Errorf("toTitleCase(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestAutoTitlePlugin_extractFromContent(t *testing.T) {
	plugin := NewAutoTitlePlugin()

	tests := []struct {
		name     string
		content  string
		expected string
	}{
		{
			name:     "simple H1",
			content:  "# My Great Title\n\nSome content here.",
			expected: "My Great Title",
		},
		{
			name:     "H1 with trailing hashes",
			content:  "# Title With Hashes #\n\nContent",
			expected: "Title With Hashes",
		},
		{
			name:     "H1 after blank lines",
			content:  "\n\n# Title After Blanks\n\nContent",
			expected: "Title After Blanks",
		},
		{
			name:     "H1 with extra spaces",
			content:  "#   Spaced Out Title   \n\nContent",
			expected: "Spaced Out Title",
		},
		{
			name:     "no H1 heading",
			content:  "## H2 Heading\n\nJust some text",
			expected: "",
		},
		{
			name:     "empty content",
			content:  "",
			expected: "",
		},
		{
			name:     "H1 not at start of line - no match",
			content:  "text # Not a heading\n",
			expected: "",
		},
		{
			name:     "multiple H1s - returns first",
			content:  "# First Title\n\n# Second Title",
			expected: "First Title",
		},
		{
			name:     "H1 with special characters",
			content:  "# C++ Programming Guide\n\nContent",
			expected: "C++ Programming Guide",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := plugin.extractFromContent(tt.content)
			if got != tt.expected {
				t.Errorf("extractFromContent() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestAutoTitlePlugin_extractFromDirectory(t *testing.T) {
	plugin := NewAutoTitlePlugin()

	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "simple directory",
			path:     "docs/index.md",
			expected: "Docs",
		},
		{
			name:     "hyphenated directory",
			path:     "getting-started/index.md",
			expected: "Getting Started",
		},
		{
			name:     "underscored directory",
			path:     "api_reference/index.md",
			expected: "Api Reference",
		},
		{
			name:     "nested directory",
			path:     "guides/advanced/index.md",
			expected: "Advanced",
		},
		{
			name:     "root index",
			path:     "index.md",
			expected: "",
		},
		{
			name:     "current dir index",
			path:     "./index.md",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := plugin.extractFromDirectory(tt.path)
			if got != tt.expected {
				t.Errorf("extractFromDirectory(%q) = %q, want %q", tt.path, got, tt.expected)
			}
		})
	}
}

func TestAutoTitlePlugin_stripDatePrefix(t *testing.T) {
	plugin := NewAutoTitlePlugin()

	tests := []struct {
		name         string
		filename     string
		expectedName string
		expectedHas  bool
	}{
		{
			name:         "date with hyphen separator",
			filename:     "2024-01-15-my-post",
			expectedName: "my-post",
			expectedHas:  true,
		},
		{
			name:         "date with underscore separator",
			filename:     "2024-01-15_my-post",
			expectedName: "my-post",
			expectedHas:  true,
		},
		{
			name:         "date without separator",
			filename:     "2024-01-15my-post",
			expectedName: "my-post",
			expectedHas:  true,
		},
		{
			name:         "no date prefix",
			filename:     "my-post",
			expectedName: "my-post",
			expectedHas:  false,
		},
		{
			name:         "date only",
			filename:     "2024-01-15",
			expectedName: "",
			expectedHas:  true,
		},
		{
			name:         "partial date - not stripped",
			filename:     "2024-01-post",
			expectedName: "2024-01-post",
			expectedHas:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotName, gotHas := plugin.stripDatePrefix(tt.filename)
			if gotName != tt.expectedName || gotHas != tt.expectedHas {
				t.Errorf("stripDatePrefix(%q) = (%q, %v), want (%q, %v)",
					tt.filename, gotName, gotHas, tt.expectedName, tt.expectedHas)
			}
		})
	}
}

func TestAutoTitlePlugin_generateFallback(t *testing.T) {
	plugin := NewAutoTitlePlugin()

	tests := []struct {
		name     string
		path     string
		contains string // Check contains because timestamp varies
	}{
		{
			name:     "path with filename",
			path:     "posts/weird.md",
			contains: "Untitled weird",
		},
		{
			name:     "empty path",
			path:     "",
			contains: "Untitled 20", // Will contain timestamp starting with year
		},
		{
			name:     "path with stem",
			path:     "test.md",
			contains: "Untitled test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := plugin.generateFallback(tt.path)
			if !strings.Contains(got, tt.contains) {
				t.Errorf("generateFallback(%q) = %q, want to contain %q", tt.path, got, tt.contains)
			}
		})
	}
}

func TestAutoTitlePlugin_inferTitle(t *testing.T) {
	plugin := NewAutoTitlePlugin()

	tests := []struct {
		name     string
		post     *models.Post
		expected string
	}{
		{
			name: "H1 takes priority",
			post: &models.Post{
				Path:    "posts/some-file.md",
				Content: "# Title From H1\n\nContent",
			},
			expected: "Title From H1",
		},
		{
			name: "index.md uses directory when no H1",
			post: &models.Post{
				Path:    "my-docs/index.md",
				Content: "No heading here",
			},
			expected: "My Docs",
		},
		{
			name: "index.md with H1 uses H1",
			post: &models.Post{
				Path:    "my-docs/index.md",
				Content: "# Documentation Home\n\nContent",
			},
			expected: "Documentation Home",
		},
		{
			name: "filename fallback when no H1 and not index",
			post: &models.Post{
				Path:    "posts/my-great-article.md",
				Content: "Just some text without heading",
			},
			expected: "My Great Article",
		},
		{
			name: "date-prefixed filename",
			post: &models.Post{
				Path:    "posts/2024-01-15-new-feature.md",
				Content: "No heading",
			},
			expected: "New Feature",
		},
		{
			name: "fallback for edge case",
			post: &models.Post{
				Path:    ".md",
				Content: "",
			},
			expected: "Untitled ", // Will have timestamp, but starts with this
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := plugin.inferTitle(tt.post)
			if tt.name == "fallback for edge case" {
				if !strings.HasPrefix(got, tt.expected) {
					t.Errorf("inferTitle() = %q, want prefix %q", got, tt.expected)
				}
			} else if got != tt.expected {
				t.Errorf("inferTitle() = %q, want %q", got, tt.expected)
			}
		})
	}
}
