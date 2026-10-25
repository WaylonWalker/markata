package templates

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"testing/fstest"
	"time"

	"github.com/WaylonWalker/markata/pkg/models"
)

func TestNewEngine(t *testing.T) {
	tests := []struct {
		name         string
		templatesDir string
		wantErr      bool
	}{
		{
			name:         "empty directory",
			templatesDir: "",
			wantErr:      false,
		},
		{
			name:         "non-existent directory",
			templatesDir: "/nonexistent/path",
			wantErr:      false, // Should not error, just won't find templates
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine, err := NewEngine(tt.templatesDir)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewEngine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && engine == nil {
				t.Error("NewEngine() returned nil engine")
			}
		})
	}
}

func TestSearchPathLoader_EmbeddedIncludesPreferProjectOverride(t *testing.T) {
	directory := t.TempDir()
	projectTemplate := filepath.Join(directory, "partials", "card.html")
	if err := os.MkdirAll(filepath.Dir(projectTemplate), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(projectTemplate, []byte("project"), 0o600); err != nil {
		t.Fatal(err)
	}

	loader := &searchPathLoader{
		searchPaths: []string{directory},
		embeddedFS: fstest.MapFS{
			"partials/card.html": &fstest.MapFile{Data: []byte("embedded")},
		},
	}

	got := loader.Abs("embedded:feed.html", "partials/card.html")
	want, err := filepath.Abs(projectTemplate)
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("Abs() = %q, want project override %q", got, want)
	}
}

func TestSearchPathLoader_EmbeddedSiblingIncludesPreferProjectOverride(t *testing.T) {
	directory := t.TempDir()
	projectTemplate := filepath.Join(directory, "partials", "card.html")
	if err := os.MkdirAll(filepath.Dir(projectTemplate), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(projectTemplate, []byte("project"), 0o600); err != nil {
		t.Fatal(err)
	}

	loader := &searchPathLoader{
		searchPaths: []string{directory},
		embeddedFS: fstest.MapFS{
			"partials/feed.html": &fstest.MapFile{Data: []byte("{% include \"card.html\" %}")},
			"partials/card.html": &fstest.MapFile{Data: []byte("embedded")},
		},
	}

	got := loader.Abs("embedded:partials/feed.html", "card.html")
	want, err := filepath.Abs(projectTemplate)
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("Abs() = %q, want project sibling override %q", got, want)
	}
}

func TestSearchPathLoader_EmbeddedIncludesFallBackToEmbeddedTemplate(t *testing.T) {
	loader := &searchPathLoader{
		searchPaths: []string{t.TempDir()},
		embeddedFS: fstest.MapFS{
			"partials/card.html": &fstest.MapFile{Data: []byte("embedded")},
		},
	}

	if got, want := loader.Abs("embedded:feed.html", "partials/card.html"), "embedded:partials/card.html"; got != want {
		t.Fatalf("Abs() = %q, want embedded fallback %q", got, want)
	}
}

func TestSearchPathLoader_GetClosesFilesystemFile(t *testing.T) {
	directory := t.TempDir()
	path := filepath.Join(directory, "template.html")
	if err := os.WriteFile(path, []byte("template"), 0o600); err != nil {
		t.Fatal(err)
	}

	loader := &searchPathLoader{searchPaths: []string{directory}}
	reader, err := loader.Get("template.html")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	content, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("ReadAll() error = %v", err)
	}
	if string(content) != "template" {
		t.Fatalf("template content = %q", content)
	}
	if err := os.Remove(path); err != nil {
		t.Fatalf("remove loaded template: %v", err)
	}
}

func TestEngine_RenderString(t *testing.T) {
	engine, err := NewEngine("")
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	title := "Hello World"
	post := &models.Post{
		Title:     &title,
		Slug:      "hello-world",
		Tags:      []string{"go", "programming"},
		Published: true,
	}

	config := &models.Config{
		Title:       "My Site",
		URL:         "https://example.com",
		Description: "A test site",
	}

	tests := []struct {
		name     string
		template string
		post     *models.Post
		config   *models.Config
		want     string
		wantErr  bool
	}{
		{
			name:     "simple variable",
			template: "{{ post.title }}",
			post:     post,
			want:     "Hello World",
		},
		{
			name:     "title shortcut",
			template: "{{ title }}",
			post:     post,
			want:     "Hello World",
		},
		{
			name:     "nested config variable",
			template: "{{ config.title }}",
			post:     post,
			config:   config,
			want:     "My Site",
		},
		{
			name:     "site_title shortcut",
			template: "{{ site_title }}",
			post:     post,
			config:   config,
			want:     "My Site",
		},
		{
			name:     "for loop with tags",
			template: "{% for tag in post.tags %}{{ tag }},{% endfor %}",
			post:     post,
			want:     "go,programming,",
		},
		{
			name:     "if condition true",
			template: "{% if post.published %}yes{% endif %}",
			post:     post,
			want:     "yes",
		},
		{
			name:     "if condition false",
			template: "{% if post.draft %}yes{% else %}no{% endif %}",
			post:     post,
			want:     "no",
		},
		{
			name:     "filter upper",
			template: "{{ post.slug | upper }}",
			post:     post,
			want:     "HELLO-WORLD",
		},
		{
			name:     "filter lower",
			template: "{{ post.title | lower }}",
			post:     post,
			want:     "hello world",
		},
		{
			name:     "filter default with value",
			template: "{{ post.title | default_if_none:\"Untitled\" }}",
			post:     post,
			want:     "Hello World",
		},
		{
			name:     "filter default with nil",
			template: "{{ post.description | default_if_none:\"No description\" }}",
			post:     post,
			want:     "No description",
		},
		{
			name:     "filter truncate",
			template: "{{ \"This is a long string that should be truncated\" | truncate:20 }}",
			post:     post,
			want:     "This is a long...",
		},
		{
			name:     "filter slugify",
			template: "{{ \"Hello World!\" | slugify }}",
			post:     post,
			want:     "hello-world",
		},
		{
			name:     "filter length string",
			template: "{{ \"hello\" | length }}",
			post:     post,
			want:     "5",
		},
		{
			name:     "filter length slice",
			template: "{{ post.tags | length }}",
			post:     post,
			want:     "2",
		},
		{
			name:     "filter join",
			template: "{{ post.tags | join:\", \" }}",
			post:     post,
			want:     "go, programming",
		},
		{
			name:     "invalid template syntax",
			template: "{{ post.title",
			post:     post,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := NewContext(tt.post, "", tt.config)
			got, err := engine.RenderString(tt.template, ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("RenderString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("RenderString() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestEngine_RenderString_DateFilters(t *testing.T) {
	engine, err := NewEngine("")
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	date := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	post := &models.Post{
		Date: &date,
	}

	tests := []struct {
		name     string
		template string
		want     string
	}{
		{
			name:     "rss_date filter",
			template: "{{ post.date | rss_date }}",
			want:     "Mon, 15 Jan 2024 10:30:00 +0000",
		},
		{
			name:     "atom_date filter",
			template: "{{ post.date | atom_date }}",
			want:     "2024-01-15T10:30:00Z",
		},
		{
			name:     "custom date format",
			template: "{{ post.date | date_format:\"2006-01-02\" }}",
			want:     "2024-01-15",
		},
		{
			name:     "human date format",
			template: "{{ post.date | human_date }}",
			want:     "Jan 15, 2024",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := NewContext(post, "", nil)
			got, err := engine.RenderString(tt.template, ctx)
			if err != nil {
				t.Errorf("RenderString() error = %v", err)
				return
			}
			if got != tt.want {
				t.Errorf("RenderString() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestContext_ToPongo2(t *testing.T) {
	title := "Test Post"
	avatar := "/8bitcc.png"
	post := &models.Post{
		Title:     &title,
		Slug:      "test-post",
		Tags:      []string{"test"},
		Published: true,
		Extra:     map[string]interface{}{"custom_field": "custom_value"},
	}

	config := &models.Config{
		Title: "Test Site",
		URL:   "https://test.com",
		Authors: models.AuthorsConfig{Authors: map[string]models.Author{
			"waylon": {
				Name:    "Waylon Walker",
				Avatar:  &avatar,
				Default: true,
				Active:  true,
				Social: map[string]string{
					"github": "https://github.com/WaylonWalker",
				},
			},
		}},
	}

	ctx := NewContext(post, "<p>Body HTML</p>", config)
	ctx.Set("extra_key", "extra_value")

	p2ctx := ctx.ToPongo2()

	// Check post is converted to map
	postMap, ok := p2ctx["post"].(map[string]interface{})
	if !ok {
		t.Error("post not set in context as map")
	}
	if p2ctx["body"] != "<p>Body HTML</p>" {
		t.Error("body not set correctly")
	}

	// Check shortcuts (now strings, not pointers)
	if p2ctx["title"] != "Test Post" {
		t.Errorf("title shortcut not set correctly, got %v", p2ctx["title"])
	}
	if p2ctx["slug"] != "test-post" {
		t.Error("slug shortcut not set correctly")
	}

	// Also check via post map
	if postMap["title"] != "Test Post" {
		t.Error("post.title not set correctly in map")
	}

	// Check config shortcuts
	if p2ctx["site_title"] != "Test Site" {
		t.Error("site_title shortcut not set correctly")
	}
	if p2ctx["site_url"] != "https://test.com" {
		t.Error("site_url shortcut not set correctly")
	}

	defaultAuthor, ok := p2ctx["default_author"].(map[string]interface{})
	if !ok {
		t.Error("default_author not set in context as map")
	} else {
		if defaultAuthor["name"] != "Waylon Walker" {
			t.Errorf("default_author.name not set correctly, got %v", defaultAuthor["name"])
		}
		if defaultAuthor["avatar"] != "/8bitcc.png" {
			t.Errorf("default_author.avatar not set correctly, got %v", defaultAuthor["avatar"])
		}
	}

	if p2ctx["default_author_id"] != "waylon" {
		t.Errorf("default_author_id not set correctly, got %v", p2ctx["default_author_id"])
	}

	resolvedSidebar, ok := p2ctx["resolved_content_sidebar"].(map[string]interface{})
	if !ok {
		t.Error("resolved_content_sidebar not set in context as map")
	} else if resolvedSidebar["enabled"] != false {
		t.Errorf("resolved_content_sidebar.enabled default incorrect, got %v", resolvedSidebar["enabled"])
	}

	// Check extra fields from post
	if p2ctx["custom_field"] != "custom_value" {
		t.Error("custom field from post.Extra not set")
	}

	// Check extra context values
	if p2ctx["extra_key"] != "extra_value" {
		t.Error("extra context value not set")
	}
}

func TestContext_ToPongo2_ConcurrentConfigIsPerRender(t *testing.T) {
	config := &models.Config{Title: "Test Site", Theme: models.ThemeConfig{Name: "test-theme"}}
	ctx := NewContext(nil, "", config)
	const renders = 128
	var wg sync.WaitGroup
	for i := 0; i < renders; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			pongoCtx := ctx.ToPongo2()
			configMap, ok := pongoCtx["config"].(map[string]interface{})
			if !ok || configMap["title"] != "Test Site" {
				t.Errorf("config context = %#v", pongoCtx["config"])
			}
		}()
	}
	wg.Wait()
}

func TestContext_ToPongo2_ResolvedContentSidebar(t *testing.T) {
	title := "Test Post"
	enabled := true
	post := &models.Post{
		Title: &title,
		Extra: map[string]interface{}{"sidebar_slug": "components/page-sidebar"},
	}

	config := &models.Config{
		Components: models.ComponentsConfig{
			ContentSidebar: models.ContentSidebarConfig{
				Enabled:  &enabled,
				Position: "left",
				Width:    "280px",
				Slug:     "components/global-sidebar",
			},
		},
	}

	ctx := NewContext(post, "", config)
	p2ctx := ctx.ToPongo2()

	resolvedSidebar, ok := p2ctx["resolved_content_sidebar"].(map[string]interface{})
	if !ok {
		t.Fatal("resolved_content_sidebar not set in context as map")
	}

	if resolvedSidebar["enabled"] != true {
		t.Errorf("resolved_content_sidebar.enabled incorrect, got %v", resolvedSidebar["enabled"])
	}
	if resolvedSidebar["position"] != "left" {
		t.Errorf("resolved_content_sidebar.position incorrect, got %v", resolvedSidebar["position"])
	}
	if resolvedSidebar["width"] != "280px" {
		t.Errorf("resolved_content_sidebar.width incorrect, got %v", resolvedSidebar["width"])
	}
	if resolvedSidebar["slug"] != "components/page-sidebar" {
		t.Errorf("resolved_content_sidebar.slug incorrect, got %v", resolvedSidebar["slug"])
	}
	if p2ctx["sidebar_slug"] != "components/page-sidebar" {
		t.Errorf("sidebar_slug extra field not exposed, got %v", p2ctx["sidebar_slug"])
	}
}

func TestResolveDocSidebar(t *testing.T) {
	enabled := true
	disabled := false

	tests := []struct {
		name         string
		config       *models.Config
		blog         bool
		wantEnabled  bool
		wantPosition string
		wantWidth    string
	}{
		{
			name:         "nil config",
			config:       nil,
			blog:         true,
			wantEnabled:  false,
			wantPosition: "right",
			wantWidth:    "250px",
		},
		{
			name: "blog show_toc enables sidebar with blog settings",
			config: &models.Config{
				Components: models.ComponentsConfig{DocSidebar: models.DocSidebarConfig{Enabled: &disabled, Position: "right", Width: "250px"}},
				Layout:     models.LayoutConfig{Blog: models.BlogLayoutConfig{ShowToc: &enabled, TocPosition: "left", TocWidth: "200px"}},
			},
			blog:         true,
			wantEnabled:  true,
			wantPosition: "left",
			wantWidth:    "200px",
		},
		{
			name: "blog show_toc ignored outside blog layout",
			config: &models.Config{
				Layout: models.LayoutConfig{Blog: models.BlogLayoutConfig{ShowToc: &enabled, TocPosition: "left"}},
			},
			blog:         false,
			wantEnabled:  false,
			wantPosition: "right",
			wantWidth:    "250px",
		},
		{
			name: "blog show_toc false keeps sidebar disabled",
			config: &models.Config{
				Layout: models.LayoutConfig{Blog: models.BlogLayoutConfig{ShowToc: &disabled}},
			},
			blog:         true,
			wantEnabled:  false,
			wantPosition: "right",
			wantWidth:    "250px",
		},
		{
			name: "doc_sidebar enabled takes precedence",
			config: &models.Config{
				Components: models.ComponentsConfig{DocSidebar: models.DocSidebarConfig{Enabled: &enabled, Position: "left", Width: "300px"}},
				Layout:     models.LayoutConfig{Blog: models.BlogLayoutConfig{ShowToc: &enabled, TocPosition: "right", TocWidth: "200px"}},
			},
			blog:         true,
			wantEnabled:  true,
			wantPosition: "left",
			wantWidth:    "300px",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ResolveDocSidebar(tt.config, tt.blog)
			if got["enabled"] != tt.wantEnabled {
				t.Errorf("enabled = %v, want %v", got["enabled"], tt.wantEnabled)
			}
			if got["position"] != tt.wantPosition {
				t.Errorf("position = %v, want %v", got["position"], tt.wantPosition)
			}
			if got["width"] != tt.wantWidth {
				t.Errorf("width = %v, want %v", got["width"], tt.wantWidth)
			}
			if got["min_depth"] != 2 || got["max_depth"] != 4 {
				t.Errorf("depth = %v-%v, want 2-4", got["min_depth"], got["max_depth"])
			}
		})
	}
}

func TestContext_ToPongo2_ResolvedDocSidebar(t *testing.T) {
	enabled := true
	config := &models.Config{
		Layout: models.LayoutConfig{Blog: models.BlogLayoutConfig{ShowToc: &enabled}},
	}

	ctx := NewContext(&models.Post{}, "", config).WithBlogLayout(true)
	sidebar, ok := ctx.ToPongo2()["resolved_doc_sidebar"].(map[string]interface{})
	if !ok {
		t.Fatal("resolved_doc_sidebar not set in context as map")
	}
	if sidebar["enabled"] != true {
		t.Errorf("resolved_doc_sidebar.enabled = %v, want true", sidebar["enabled"])
	}
	if !ctx.Clone().BlogLayout {
		t.Error("Clone() did not copy BlogLayout")
	}
}

func TestContext_Clone(t *testing.T) {
	title := "Original"
	post := &models.Post{Title: &title}
	config := &models.Config{Title: "Site"}

	original := NewContext(post, "body", config)
	original.Set("key", "value")

	clone := original.Clone()

	// Modify original
	original.Set("key", "modified")
	original.Set("new_key", "new_value")

	// Clone should not be affected
	if clone.Get("key") != "value" {
		t.Error("Clone was affected by modification to original")
	}
	if clone.Get("new_key") != nil {
		t.Error("Clone has key that was added after cloning")
	}
}

func TestContext_Merge(t *testing.T) {
	title1 := "Post 1"
	title2 := "Post 2"
	post1 := &models.Post{Title: &title1}
	post2 := &models.Post{Title: &title2}
	config := &models.Config{Title: "Site"}

	ctx1 := NewContext(post1, "body1", nil)
	ctx1.Set("key1", "value1")

	ctx2 := NewContext(post2, "body2", config)
	ctx2.Set("key2", "value2")

	ctx1.Merge(ctx2)

	// Post should be overwritten
	if ctx1.Post != post2 {
		t.Error("Post not merged")
	}

	// Body should be overwritten
	if ctx1.Body != "body2" {
		t.Error("Body not merged")
	}

	// Config should be set
	if ctx1.Config != config {
		t.Error("Config not merged")
	}

	// Both extra keys should exist
	if ctx1.Get("key1") != "value1" {
		t.Error("Original extra key lost")
	}
	if ctx1.Get("key2") != "value2" {
		t.Error("Merged extra key not added")
	}
}

func TestEngine_Render_SlidesTemplate(t *testing.T) {
	engine, err := NewEngine("../../templates")
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	title := "Slides"
	post := &models.Post{
		Title: &title,
		Slug:  "slides",
		Href:  "/slides/",
	}
	config := models.NewConfig()
	config.Title = "My Site"
	config.URL = "https://example.com"

	body := `<h2>Intro</h2><p>Welcome</p><h3>Details</h3><p>More</p><hr/><h2>Wrap Up</h2>`
	ctx := NewContext(post, body, config)

	rendered, err := engine.Render("slides.html", ctx)
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}

	checks := []string{
		`<body class="slides-page">`,
		`https://cdn.jsdelivr.net/npm/reveal.js@5/dist/reveal.css`,
		`<div class="reveal">`,
		`<div class="slides">`,
		`<section><section><div class="slide-content"><h2>Intro</h2><p>Welcome</p></div></section><section><div class="slide-content"><h3>Details</h3><p>More</p></div></section></section>`,
		`<section><div class="slide-content"><h2>Wrap Up</h2></div></section>`,
		`function slideDeckOptions() {`,
		`Reveal.initialize(slideDeckOptions()).then(function() {`,
		`min-height: 100dvh;`,
	}

	for _, check := range checks {
		if !strings.Contains(rendered, check) {
			t.Fatalf("slides template output missing %q", check)
		}
	}

	for _, hidden := range []string{`<header class="site-header">`, `<footer class="site-footer">`} {
		if strings.Contains(rendered, hidden) {
			t.Fatalf("slides template should not render %q, output: %q", hidden, rendered)
		}
	}
}

func TestEngine_Render_SlidesTemplate_UsesVendoredRevealAssets(t *testing.T) {
	engine, err := NewEngine("../../templates")
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	title := "Slides"
	post := &models.Post{Title: &title, Slug: "slides", Href: "/slides/"}
	config := models.NewConfig()
	config.Title = "My Site"
	config.URL = "https://example.com"
	config.Extra = map[string]interface{}{
		"asset_urls": map[string]interface{}{
			"revealjs-css": "/assets/vendor/revealjs/reveal.css",
			"revealjs-js":  "/assets/vendor/revealjs/reveal.js",
		},
	}

	ctx := NewContext(post, `<h2>Intro</h2>`, config)
	rendered, err := engine.Render("slides.html", ctx)
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}

	checks := []string{
		`href="/assets/vendor/revealjs/reveal.css"`,
		`src="/assets/vendor/revealjs/reveal.js"`,
	}

	for _, check := range checks {
		if !strings.Contains(rendered, check) {
			t.Fatalf("slides template vendored output missing %q", check)
		}
	}
}
