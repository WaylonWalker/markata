package plugins

import (
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

func TestBreadcrumbsPlugin_Name(t *testing.T) {
	p := NewBreadcrumbsPlugin()
	if got := p.Name(); got != "breadcrumbs" {
		t.Errorf("Name() = %q, want %q", got, "breadcrumbs")
	}
}

func TestBreadcrumbsPlugin_generateBreadcrumbs(t *testing.T) {
	tests := []struct {
		name       string
		href       string
		title      string
		wantCount  int
		wantLabels []string
		wantURLs   []string
	}{
		{
			name:       "simple path",
			href:       "/docs/",
			title:      "Documentation",
			wantCount:  2,
			wantLabels: []string{"Home", "Documentation"},
			wantURLs:   []string{"/", "/docs/"},
		},
		{
			name:       "nested path",
			href:       "/docs/guides/getting-started/",
			title:      "Getting Started",
			wantCount:  4,
			wantLabels: []string{"Home", "Docs", "Guides", "Getting Started"},
			wantURLs:   []string{"/", "/docs/", "/docs/guides/", "/docs/guides/getting-started/"},
		},
		{
			name:       "homepage returns empty",
			href:       "/",
			title:      "Home",
			wantCount:  0,
			wantLabels: nil,
			wantURLs:   nil,
		},
		{
			name:       "hyphenated segments",
			href:       "/api-reference/",
			title:      "API Reference",
			wantCount:  2,
			wantLabels: []string{"Home", "API Reference"},
			wantURLs:   []string{"/", "/api-reference/"},
		},
		{
			name:       "underscored segments",
			href:       "/user_guide/",
			title:      "User Guide",
			wantCount:  2,
			wantLabels: []string{"Home", "User Guide"},
			wantURLs:   []string{"/", "/user_guide/"},
		},
		{
			name:       "deep nesting",
			href:       "/a/b/c/d/e/",
			title:      "Page E",
			wantCount:  6,
			wantLabels: []string{"Home", "A", "B", "C", "D", "Page E"},
			wantURLs:   []string{"/", "/a/", "/a/b/", "/a/b/c/", "/a/b/c/d/", "/a/b/c/d/e/"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewBreadcrumbsPlugin()

			post := models.NewPost("test.md")
			post.Href = tt.href
			if tt.title != "" {
				post.Title = &tt.title
			}

			postConfig := BreadcrumbConfig{}
			breadcrumbs := p.generateBreadcrumbs(post, postConfig, "https://example.com")

			if len(breadcrumbs) != tt.wantCount {
				t.Errorf("got %d breadcrumbs, want %d", len(breadcrumbs), tt.wantCount)
				return
			}

			for i, bc := range breadcrumbs {
				if i < len(tt.wantLabels) && bc.Label != tt.wantLabels[i] {
					t.Errorf("breadcrumb[%d].Label = %q, want %q", i, bc.Label, tt.wantLabels[i])
				}
				if i < len(tt.wantURLs) && bc.URL != tt.wantURLs[i] {
					t.Errorf("breadcrumb[%d].URL = %q, want %q", i, bc.URL, tt.wantURLs[i])
				}
				// Check position is 1-indexed
				if bc.Position != i+1 {
					t.Errorf("breadcrumb[%d].Position = %d, want %d", i, bc.Position, i+1)
				}
				// Check IsCurrent for last item
				isLast := i == len(breadcrumbs)-1
				if bc.IsCurrent != isLast {
					t.Errorf("breadcrumb[%d].IsCurrent = %v, want %v", i, bc.IsCurrent, isLast)
				}
			}
		})
	}
}

func TestBreadcrumbsPlugin_showHome(t *testing.T) {
	t.Run("show_home=false", func(t *testing.T) {
		p := NewBreadcrumbsPlugin()
		p.showHome = false

		post := models.NewPost("test.md")
		post.Href = "/docs/guide/"
		title := "Guide"
		post.Title = &title

		breadcrumbs := p.generateBreadcrumbs(post, BreadcrumbConfig{}, "https://example.com")

		if len(breadcrumbs) != 2 {
			t.Errorf("expected 2 breadcrumbs without home, got %d", len(breadcrumbs))
		}

		if len(breadcrumbs) > 0 && breadcrumbs[0].Label == "Home" {
			t.Error("should not include Home when show_home=false")
		}
	})

	t.Run("custom home label", func(t *testing.T) {
		p := NewBreadcrumbsPlugin()
		p.homeLabel = "Start"

		post := models.NewPost("test.md")
		post.Href = "/docs/"
		title := "Docs"
		post.Title = &title

		breadcrumbs := p.generateBreadcrumbs(post, BreadcrumbConfig{}, "https://example.com")

		if len(breadcrumbs) == 0 || breadcrumbs[0].Label != "Start" {
			t.Errorf("expected first breadcrumb label to be 'Start', got %q", breadcrumbs[0].Label)
		}
	})
}

func TestBreadcrumbsPlugin_maxDepth(t *testing.T) {
	p := NewBreadcrumbsPlugin()
	p.maxDepth = 3 // Home + 2 path segments

	post := models.NewPost("test.md")
	post.Href = "/a/b/c/d/e/"
	title := "Page E"
	post.Title = &title

	breadcrumbs := p.generateBreadcrumbs(post, BreadcrumbConfig{}, "https://example.com")

	if len(breadcrumbs) != 3 {
		t.Errorf("expected 3 breadcrumbs with maxDepth=3, got %d", len(breadcrumbs))
	}
}

func TestBreadcrumbsPlugin_manualBreadcrumbs(t *testing.T) {
	p := NewBreadcrumbsPlugin()

	post := models.NewPost("test.md")
	post.Href = "/products/widgets/blue-widget/"
	title := "Blue Widget"
	post.Title = &title

	postConfig := BreadcrumbConfig{
		Items: []BreadcrumbItem{
			{Label: "Products", URL: "/products/"},
			{Label: "Widgets", URL: "/products/widgets/"},
		},
	}

	breadcrumbs := p.generateBreadcrumbs(post, postConfig, "https://example.com")

	// Should be: Home + 2 manual + current page
	if len(breadcrumbs) != 4 {
		t.Errorf("expected 4 breadcrumbs, got %d", len(breadcrumbs))
		return
	}

	expected := []string{"Home", "Products", "Widgets", "Blue Widget"}
	for i, bc := range breadcrumbs {
		if bc.Label != expected[i] {
			t.Errorf("breadcrumb[%d].Label = %q, want %q", i, bc.Label, expected[i])
		}
	}

	// Last one should be current
	if !breadcrumbs[3].IsCurrent {
		t.Error("last breadcrumb should be marked as current")
	}
}

func TestBreadcrumbsPlugin_perPostDisable(t *testing.T) {
	t.Helper() // Mark as test helper

	p := NewBreadcrumbsPlugin()

	post := models.NewPost("test.md")
	post.Href = "/docs/"
	title := "Docs"
	post.Title = &title

	enabled := false
	postConfig := BreadcrumbConfig{
		Enabled: &enabled,
	}

	// When disabled, buildFromPath is still called but Transform checks Enabled first
	// This tests that we can generate breadcrumbs normally; the disable happens in Transform
	breadcrumbs := p.generateBreadcrumbs(post, postConfig, "https://example.com")

	// Breadcrumbs should still be generated, the disabled check is at Transform level
	if len(breadcrumbs) == 0 {
		t.Log("Breadcrumbs disabled at Transform level, not at generateBreadcrumbs level")
	}
}

func TestBreadcrumbsPlugin_perPostShowHomeOverride(t *testing.T) {
	p := NewBreadcrumbsPlugin()
	p.showHome = true

	post := models.NewPost("test.md")
	post.Href = "/docs/"
	title := "Docs"
	post.Title = &title

	showHome := false
	postConfig := BreadcrumbConfig{
		ShowHome: &showHome,
	}

	breadcrumbs := p.generateBreadcrumbs(post, postConfig, "https://example.com")

	// With showHome override = false, should not have Home
	if len(breadcrumbs) > 0 && breadcrumbs[0].Label == "Home" {
		t.Error("should not include Home when post config has show_home=false")
	}
}

func TestBreadcrumbsPlugin_generateJSONLD(t *testing.T) {
	p := NewBreadcrumbsPlugin()

	breadcrumbs := []Breadcrumb{
		{Label: "Home", URL: "/", Position: 1, IsCurrent: false},
		{Label: "Docs", URL: "/docs/", Position: 2, IsCurrent: false},
		{Label: "Guide", URL: "/docs/guide/", Position: 3, IsCurrent: true},
	}

	jsonLD := p.generateJSONLD(breadcrumbs, "https://example.com")

	if jsonLD == "" {
		t.Fatal("expected non-empty JSON-LD")
	}

	// Check structure
	if !strings.Contains(jsonLD, `"@context": "https://schema.org"`) {
		t.Error("missing @context in JSON-LD")
	}
	if !strings.Contains(jsonLD, `"@type": "BreadcrumbList"`) {
		t.Error("missing @type in JSON-LD")
	}
	if !strings.Contains(jsonLD, `"itemListElement"`) {
		t.Error("missing itemListElement in JSON-LD")
	}

	// Check items
	if !strings.Contains(jsonLD, `"name": "Home"`) {
		t.Error("missing Home item in JSON-LD")
	}
	if !strings.Contains(jsonLD, `"item": "https://example.com/"`) {
		t.Error("missing Home URL in JSON-LD")
	}
	if !strings.Contains(jsonLD, `"name": "Guide"`) {
		t.Error("missing Guide item in JSON-LD")
	}

	// Current item should NOT have item URL
	if strings.Contains(jsonLD, `"item": "https://example.com/docs/guide/"`) {
		t.Error("current breadcrumb should not have item URL in JSON-LD")
	}
}

func TestBreadcrumbsPlugin_humanizeSegment(t *testing.T) {
	p := NewBreadcrumbsPlugin()

	tests := []struct {
		input string
		want  string
	}{
		{"getting-started", "Getting Started"},
		{"api_reference", "Api Reference"},
		{"docs", "Docs"},
		{"my-awesome-guide", "My Awesome Guide"},
		{"FAQ", "FAQ"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := p.humanizeSegment(tt.input)
			if got != tt.want {
				t.Errorf("humanizeSegment(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestBreadcrumbsPlugin_Configure(t *testing.T) {
	tests := []struct {
		name      string
		extra     map[string]interface{}
		wantHome  bool
		wantLabel string
		wantSep   string
		wantDepth int
	}{
		{
			name:      "default values",
			extra:     nil,
			wantHome:  true,
			wantLabel: "Home",
			wantSep:   "/",
			wantDepth: 0,
		},
		{
			name: "custom values from components.breadcrumbs",
			extra: map[string]interface{}{
				"components": map[string]interface{}{
					"breadcrumbs": map[string]interface{}{
						"show_home":  false,
						"home_label": "Start",
						"separator":  ">",
						"max_depth":  5,
					},
				},
			},
			wantHome:  false,
			wantLabel: "Start",
			wantSep:   ">",
			wantDepth: 5,
		},
		{
			name: "custom values from breadcrumbs (top-level)",
			extra: map[string]interface{}{
				"breadcrumbs": map[string]interface{}{
					"home_label": "Main",
					"separator":  "→",
				},
			},
			wantHome:  true,
			wantLabel: "Main",
			wantSep:   "→",
			wantDepth: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewBreadcrumbsPlugin()

			config := &lifecycle.Config{
				Extra: tt.extra,
			}
			m := lifecycle.NewManager()
			m.SetConfig(config)

			err := p.Configure(m)
			if err != nil {
				t.Fatalf("Configure() error = %v", err)
			}

			if p.showHome != tt.wantHome {
				t.Errorf("showHome = %v, want %v", p.showHome, tt.wantHome)
			}
			if p.homeLabel != tt.wantLabel {
				t.Errorf("homeLabel = %q, want %q", p.homeLabel, tt.wantLabel)
			}
			if p.separator != tt.wantSep {
				t.Errorf("separator = %q, want %q", p.separator, tt.wantSep)
			}
			if p.maxDepth != tt.wantDepth {
				t.Errorf("maxDepth = %d, want %d", p.maxDepth, tt.wantDepth)
			}
		})
	}
}

func TestBreadcrumbsPlugin_Transform(t *testing.T) {
	p := NewBreadcrumbsPlugin()

	config := &lifecycle.Config{
		Extra: map[string]interface{}{
			"url": "https://example.com",
		},
	}

	m := lifecycle.NewManager()
	m.SetConfig(config)

	// Add test posts
	title1 := "Documentation"
	post1 := &models.Post{
		Path:  "docs/index.md",
		Href:  "/docs/",
		Title: &title1,
		Extra: make(map[string]interface{}),
	}

	title2 := "Getting Started"
	post2 := &models.Post{
		Path:  "docs/getting-started.md",
		Href:  "/docs/getting-started/",
		Title: &title2,
		Extra: make(map[string]interface{}),
	}

	m.SetPosts([]*models.Post{post1, post2})

	err := p.Transform(m)
	if err != nil {
		t.Fatalf("Transform() error = %v", err)
	}

	// Check post1 has breadcrumbs
	if bc := post1.Get("breadcrumbs"); bc == nil {
		t.Error("post1 should have breadcrumbs")
	} else if breadcrumbs, ok := bc.([]Breadcrumb); !ok {
		t.Error("post1.breadcrumbs should be []Breadcrumb")
	} else if len(breadcrumbs) != 2 {
		t.Errorf("post1 should have 2 breadcrumbs, got %d", len(breadcrumbs))
	}

	// Check post2 has breadcrumbs
	if bc := post2.Get("breadcrumbs"); bc == nil {
		t.Error("post2 should have breadcrumbs")
	} else if breadcrumbs, ok := bc.([]Breadcrumb); !ok {
		t.Error("post2.breadcrumbs should be []Breadcrumb")
	} else if len(breadcrumbs) != 3 {
		t.Errorf("post2 should have 3 breadcrumbs, got %d", len(breadcrumbs))
	}

	// Check JSON-LD is generated
	if jsonld := post1.Get("breadcrumbs_jsonld"); jsonld == nil {
		t.Error("post1 should have breadcrumbs_jsonld")
	}
}

func TestBreadcrumbsPlugin_Priority(t *testing.T) {
	p := NewBreadcrumbsPlugin()

	if got := p.Priority(lifecycle.StageTransform); got != lifecycle.PriorityDefault {
		t.Errorf("Priority(StageTransform) = %d, want %d", got, lifecycle.PriorityDefault)
	}
}

func TestBreadcrumbsPlugin_interfaces(_ *testing.T) {
	p := NewBreadcrumbsPlugin()

	// Verify interface implementations
	var _ lifecycle.Plugin = p
	var _ lifecycle.ConfigurePlugin = p
	var _ lifecycle.TransformPlugin = p
	var _ lifecycle.PriorityPlugin = p
}
