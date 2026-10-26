package plugins

import (
	"path/filepath"
	"testing"

	"github.com/WaylonWalker/markata/pkg/models"
)

func TestNormalizeCustomSlug(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"empty string", "", ""},
		{"slash only", "/", ""},
		{"simple slug", "my-page", "my-page"},
		{"leading slash", "/my-page", "my-page"},
		{"trailing slash", "my-page/", "my-page"},
		{"both slashes", "/my-page/", "my-page"},
		{"nested path", "docs/guides/install", "docs/guides/install"},
		{"nested with leading slash", "/docs/guides/install", "docs/guides/install"},
		{"nested with trailing slash", "docs/guides/install/", "docs/guides/install"},
		{"nested with both slashes", "/docs/guides/install/", "docs/guides/install"},
		{"multiple leading slashes", "//my-page", "my-page"}, // strings.Trim removes all leading/trailing slashes
		{"whitespace", "  my-page  ", "my-page"},
		{"whitespace with slashes", "  /my-page/  ", "my-page"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := normalizeCustomSlug(tt.input)
			if got != tt.want {
				t.Errorf("normalizeCustomSlug(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestLoadPlugin_CustomSlugs(t *testing.T) {
	plugin := NewLoadPlugin()

	tests := []struct {
		name        string
		content     string
		wantSlug    string
		wantHref    string
		description string
	}{
		{
			name: "explicit empty slug for homepage",
			content: `---
title: My Homepage
slug: ""
published: true
---
Welcome to my site!`,
			wantSlug:    "",
			wantHref:    "/",
			description: "An explicit empty slug should become the homepage",
		},
		{
			name: "slash slug for homepage",
			content: `---
title: My Homepage
slug: /
published: true
---
Welcome to my site!`,
			wantSlug:    "",
			wantHref:    "/",
			description: "A slash slug should become the homepage",
		},
		{
			name: "custom absolute path",
			content: `---
title: About Me
slug: /about
published: true
---
About page content`,
			wantSlug:    "about",
			wantHref:    "/about/",
			description: "Leading slash should be stripped",
		},
		{
			name: "custom nested path",
			content: `---
title: Install Guide
slug: /docs/guides/install
published: true
---
Installation instructions`,
			wantSlug:    "docs/guides/install",
			wantHref:    "/docs/guides/install/",
			description: "Nested paths should preserve structure",
		},
		{
			name: "slug with trailing slash",
			content: `---
title: Blog Post
slug: blog/my-post/
published: true
---
Blog content`,
			wantSlug:    "blog/my-post",
			wantHref:    "/blog/my-post/",
			description: "Trailing slash should be stripped",
		},
		{
			name: "no slug uses auto-generated",
			content: `---
title: Auto Generated Slug
published: true
---
Content`,
			wantSlug:    "test",
			wantHref:    "/test/",
			description: "Without slug in frontmatter, basename is used (not title)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			post, err := plugin.parseFile("test.md", tt.content)
			if err != nil {
				t.Fatalf("parseFile() error = %v", err)
			}

			if post.Slug != tt.wantSlug {
				t.Errorf("Slug = %q, want %q (%s)", post.Slug, tt.wantSlug, tt.description)
			}

			if post.Href != tt.wantHref {
				t.Errorf("Href = %q, want %q (%s)", post.Href, tt.wantHref, tt.description)
			}
		})
	}
}

func TestLoadPlugin_SlugExplicitFlag(t *testing.T) {
	plugin := NewLoadPlugin()

	tests := []struct {
		name     string
		content  string
		explicit bool
	}{
		{
			name: "explicit empty slug sets flag",
			content: `---
slug: ""
---
Content`,
			explicit: true,
		},
		{
			name: "explicit slug sets flag",
			content: `---
slug: custom-slug
---
Content`,
			explicit: true,
		},
		{
			name: "no slug does not set flag",
			content: `---
title: No Slug
---
Content`,
			explicit: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			post, err := plugin.parseFile("test.md", tt.content)
			if err != nil {
				t.Fatalf("parseFile() error = %v", err)
			}

			hasFlag := post.Has("_slug_explicit")
			if hasFlag != tt.explicit {
				t.Errorf("_slug_explicit flag = %v, want %v", hasFlag, tt.explicit)
			}
		})
	}
}

func TestLoadPlugin_PathSlugMode(t *testing.T) {
	plugin := NewLoadPlugin()
	plugin.slugMode = models.SlugModePath

	tests := []struct {
		name     string
		path     string
		wantSlug string
		wantHref string
	}{
		{name: "posts nested file", path: "posts/notes/today.md", wantSlug: "notes/today", wantHref: "/notes/today/"},
		{name: "posts index file", path: "posts/notes/index.md", wantSlug: "notes", wantHref: "/notes/"},
		{name: "pages readme file", path: "pages/docs/README.md", wantSlug: "docs", wantHref: "/docs/"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			post, err := plugin.parseFile(tt.path, "---\npublished: true\n---\ncontent")
			if err != nil {
				t.Fatalf("parseFile() error = %v", err)
			}
			if post.Slug != tt.wantSlug {
				t.Errorf("Slug = %q, want %q", post.Slug, tt.wantSlug)
			}
			if post.Href != tt.wantHref {
				t.Errorf("Href = %q, want %q", post.Href, tt.wantHref)
			}
		})
	}
}

func TestLoadPlugin_SlugRules(t *testing.T) {
	plugin := NewLoadPlugin()
	plugin.slugMode = models.SlugModeFlat
	plugin.slugRules = []models.SlugRule{
		{Prefix: "posts/blog", Mode: models.SlugModeFlat},
		{Prefix: "posts/notes", Mode: models.SlugModePath},
	}

	tests := []struct {
		name     string
		path     string
		wantSlug string
	}{
		{name: "blog stays flat", path: "posts/blog/launch-post.md", wantSlug: "launch-post"},
		{name: "notes become nested", path: "posts/notes/today.md", wantSlug: "notes/today"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			post, err := plugin.parseFile(tt.path, "---\npublished: true\n---\ncontent")
			if err != nil {
				t.Fatalf("parseFile() error = %v", err)
			}
			if post.Slug != tt.wantSlug {
				t.Errorf("Slug = %q, want %q", post.Slug, tt.wantSlug)
			}
		})
	}
}

func TestOverwriteCheckPlugin_DetectsConflicts(t *testing.T) {
	// Create posts with explicitly set empty slugs (both would be homepage)
	posts := []*models.Post{
		{Path: "pages/home.md", Slug: "", Published: true},
		{Path: "pages/index.md", Slug: "", Published: true},
	}

	plugin := NewOverwriteCheckPlugin()

	// Test the getPostOutputPath helper - both should produce the same path
	outputDir := "public"

	path1 := plugin.getPostOutputPath(outputDir, posts[0])
	path2 := plugin.getPostOutputPath(outputDir, posts[1])

	if path1 != path2 {
		t.Errorf("Expected same output path for conflicting slugs, got %q and %q", path1, path2)
	}

	expected := filepath.Join("public", "index.html")
	if path1 != expected {
		t.Errorf("Expected output path %q, got %q", expected, path1)
	}
}

func TestOverwriteCheckPlugin_NoConflict(t *testing.T) {
	// Create posts with different slugs
	posts := []*models.Post{
		{Path: "pages/about.md", Slug: "about", Published: true},
		{Path: "pages/contact.md", Slug: "contact", Published: true},
	}

	plugin := NewOverwriteCheckPlugin()
	outputDir := "public"

	path1 := plugin.getPostOutputPath(outputDir, posts[0])
	path2 := plugin.getPostOutputPath(outputDir, posts[1])

	if path1 == path2 {
		t.Errorf("Expected different output paths, but both are %q", path1)
	}
}
