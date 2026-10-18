package sidebar

import (
	"testing"

	"github.com/WaylonWalker/markata/pkg/models"
)

func boolPtr(v bool) *bool {
	return &v
}

func TestBuilder_BuildFromFeed(t *testing.T) {
	title1 := "Getting Started"
	title2 := "Installation"

	feed := &models.FeedConfig{
		Slug:  "docs",
		Title: "Documentation",
		Posts: []*models.Post{
			{Slug: "getting-started", Href: "/docs/getting-started/", Title: &title1},
			{Slug: "installation", Href: "/docs/installation/", Title: &title2},
		},
	}

	builder := NewBuilder(nil, nil, nil)
	items := builder.BuildFromFeed(feed)

	if len(items) != 2 {
		t.Errorf("BuildFromFeed() got %d items, want 2", len(items))
	}
	if items[0].Title != "Getting Started" {
		t.Errorf("BuildFromFeed() first title = %v, want 'Getting Started'", items[0].Title)
	}
	if items[0].Href != "/docs/getting-started/" {
		t.Errorf("BuildFromFeed() first href = %v, want '/docs/getting-started/'", items[0].Href)
	}
}

func TestBuilder_BuildFromFeed_Empty(t *testing.T) {
	builder := NewBuilder(nil, nil, nil)

	// nil feed
	items := builder.BuildFromFeed(nil)
	if items != nil {
		t.Errorf("BuildFromFeed(nil) got %v, want nil", items)
	}

	// empty posts
	feed := &models.FeedConfig{
		Slug:  "empty",
		Title: "Empty Feed",
		Posts: []*models.Post{},
	}
	items = builder.BuildFromFeed(feed)
	if items != nil {
		t.Errorf("BuildFromFeed(empty) got %v, want nil", items)
	}
}

func TestBuilder_BuildMultiFeed(t *testing.T) {
	title1, title2 := "Doc 1", "Guide 1"

	feeds := map[string]*models.FeedConfig{
		"docs": {
			Slug:  "docs",
			Title: "Documentation",
			Posts: []*models.Post{{Title: &title1, Href: "/docs/doc-1/"}},
		},
		"guides": {
			Slug:  "guides",
			Title: "Guides",
			Posts: []*models.Post{{Title: &title2, Href: "/guides/guide-1/"}},
		},
	}

	builder := NewBuilder(nil, feeds, nil)
	items := builder.BuildMultiFeed([]string{"docs", "guides"}, nil)

	if len(items) != 2 {
		t.Errorf("BuildMultiFeed() got %d sections, want 2", len(items))
	}
	if items[0].Title != "Documentation" {
		t.Errorf("BuildMultiFeed() first section title = %v, want 'Documentation'", items[0].Title)
	}
	if len(items[0].Children) != 1 {
		t.Errorf("BuildMultiFeed() first section has %d children, want 1", len(items[0].Children))
	}
}

func TestBuilder_BuildMultiFeed_WithSections(t *testing.T) {
	title1, title2, title3 := "Doc 1", "Doc 2", "Guide 1"

	feeds := map[string]*models.FeedConfig{
		"docs": {
			Slug:  "docs",
			Title: "Documentation",
			Posts: []*models.Post{
				{Title: &title1, Href: "/docs/doc-1/"},
				{Title: &title2, Href: "/docs/doc-2/"},
				{Title: &title3, Href: "/docs/doc-3/"},
			},
		},
	}

	sections := []models.MultiFeedSection{
		{
			Feed:     "docs",
			Title:    "Custom Title",
			MaxItems: 2,
		},
	}

	builder := NewBuilder(nil, feeds, nil)
	items := builder.BuildMultiFeed(nil, sections)

	if len(items) != 1 {
		t.Errorf("BuildMultiFeed() got %d sections, want 1", len(items))
	}
	if items[0].Title != "Custom Title" {
		t.Errorf("BuildMultiFeed() section title = %v, want 'Custom Title'", items[0].Title)
	}
	if len(items[0].Children) != 2 {
		t.Errorf("BuildMultiFeed() section has %d children, want 2 (MaxItems limit)", len(items[0].Children))
	}
}

func TestBuilder_BuildFromDirectory(t *testing.T) {
	title1 := "Getting Started"
	title2 := "Advanced"

	posts := []*models.Post{
		{Path: "docs/advanced.md", Href: "/docs/advanced/", Title: &title2, Extra: map[string]interface{}{"nav_order": 2}},
		{Path: "docs/getting-started.md", Href: "/docs/getting-started/", Title: &title1, Extra: map[string]interface{}{"nav_order": 1}},
	}

	builder := NewBuilder(nil, nil, posts)
	items := builder.BuildFromDirectory(&models.SidebarAutoGenerate{
		Directory: "docs",
		OrderBy:   "nav_order",
	})

	if len(items) != 2 {
		t.Errorf("BuildFromDirectory() got %d items, want 2", len(items))
	}
	// Should be sorted by nav_order
	if items[0].Title != "Getting Started" {
		t.Errorf("BuildFromDirectory() first title = %v, want 'Getting Started' (nav_order=1)", items[0].Title)
	}
}

func TestBuilder_BuildFromDirectory_Empty(t *testing.T) {
	builder := NewBuilder(nil, nil, nil)

	// nil config
	items := builder.BuildFromDirectory(nil)
	if items != nil {
		t.Errorf("BuildFromDirectory(nil) got %v, want nil", items)
	}

	// empty directory
	items = builder.BuildFromDirectory(&models.SidebarAutoGenerate{
		Directory: "",
	})
	if items != nil {
		t.Errorf("BuildFromDirectory(empty dir) got %v, want nil", items)
	}
}

func TestBuilder_BuildFromFeeds(t *testing.T) {
	title1, title2 := "Doc 1", "Guide 1"

	feeds := map[string]*models.FeedConfig{
		"docs": {
			Slug:         "docs",
			Title:        "Documentation",
			Sidebar:      boolPtr(true),
			SidebarOrder: 1,
			Posts:        []*models.Post{{Title: &title1, Href: "/docs/doc-1/"}},
		},
		"guides": {
			Slug:         "guides",
			Title:        "Guides",
			Sidebar:      boolPtr(true),
			SidebarOrder: 2,
			Posts:        []*models.Post{{Title: &title2, Href: "/guides/guide-1/"}},
		},
		"hidden": {
			Slug:    "hidden",
			Title:   "Hidden",
			Sidebar: boolPtr(false), // Not included
			Posts:   []*models.Post{},
		},
	}

	builder := NewBuilder(nil, feeds, nil)
	items := builder.BuildFromFeeds()

	if len(items) != 2 {
		t.Errorf("BuildFromFeeds() got %d items, want 2", len(items))
	}
	// Should be sorted by SidebarOrder
	if items[0].Title != "Documentation" {
		t.Errorf("BuildFromFeeds() first title = %v, want 'Documentation'", items[0].Title)
	}
	if items[1].Title != "Guides" {
		t.Errorf("BuildFromFeeds() second title = %v, want 'Guides'", items[1].Title)
	}
}

func TestSidebarConfig_ResolveForPath(t *testing.T) {
	tests := []struct {
		name      string
		paths     map[string]*models.PathSidebarConfig
		inputPath string
		wantFound bool
		wantTitle string
	}{
		{
			name:      "no paths configured",
			paths:     nil,
			inputPath: "/docs/getting-started/",
			wantFound: false,
		},
		{
			name: "exact match",
			paths: map[string]*models.PathSidebarConfig{
				"/docs/": {Title: "Documentation"},
			},
			inputPath: "/docs/",
			wantFound: true,
			wantTitle: "Documentation",
		},
		{
			name: "prefix match",
			paths: map[string]*models.PathSidebarConfig{
				"/docs/": {Title: "Documentation"},
			},
			inputPath: "/docs/guides/getting-started/",
			wantFound: true,
			wantTitle: "Documentation",
		},
		{
			name: "longest prefix wins",
			paths: map[string]*models.PathSidebarConfig{
				"/docs/":     {Title: "Documentation"},
				"/docs/api/": {Title: "API Reference"},
			},
			inputPath: "/docs/api/endpoints/",
			wantFound: true,
			wantTitle: "API Reference",
		},
		{
			name: "no match",
			paths: map[string]*models.PathSidebarConfig{
				"/docs/": {Title: "Documentation"},
			},
			inputPath: "/blog/my-post/",
			wantFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &models.SidebarConfig{Paths: tt.paths}
			got, found := s.ResolveForPath(tt.inputPath)

			if found != tt.wantFound {
				t.Errorf("ResolveForPath() found = %v, want %v", found, tt.wantFound)
			}
			if found && got.Title != tt.wantTitle {
				t.Errorf("ResolveForPath() title = %v, want %v", got.Title, tt.wantTitle)
			}
		})
	}
}

func TestSidebarConfig_GetEffectiveConfig(t *testing.T) {
	enabled := true
	collapsible := false

	s := &models.SidebarConfig{
		Enabled:  &enabled,
		Position: "left",
		Width:    "280px",
		Paths: map[string]*models.PathSidebarConfig{
			"/docs/": {
				Title:       "Documentation",
				Position:    "right",
				Collapsible: &collapsible,
			},
		},
	}

	// Path that matches
	effective := s.GetEffectiveConfig("/docs/intro/")
	if effective.Position != "right" {
		t.Errorf("GetEffectiveConfig() position = %v, want 'right'", effective.Position)
	}
	if effective.Title != "Documentation" {
		t.Errorf("GetEffectiveConfig() title = %v, want 'Documentation'", effective.Title)
	}
	if effective.IsCollapsible() != false {
		t.Errorf("GetEffectiveConfig() collapsible = %v, want false", effective.IsCollapsible())
	}

	// Path that doesn't match
	defaultConfig := s.GetEffectiveConfig("/blog/post/")
	if defaultConfig != s {
		t.Errorf("GetEffectiveConfig() for non-matching path should return original config")
	}
}

func TestFeedConfig_GetSidebarTitle(t *testing.T) {
	tests := []struct {
		name      string
		feed      *models.FeedConfig
		wantTitle string
	}{
		{
			name: "uses SidebarTitle when set",
			feed: &models.FeedConfig{
				Title:        "Documentation",
				SidebarTitle: "Docs",
			},
			wantTitle: "Docs",
		},
		{
			name: "falls back to Title",
			feed: &models.FeedConfig{
				Title:        "Documentation",
				SidebarTitle: "",
			},
			wantTitle: "Documentation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.feed.GetSidebarTitle()
			if got != tt.wantTitle {
				t.Errorf("GetSidebarTitle() = %v, want %v", got, tt.wantTitle)
			}
		})
	}
}

func TestTitleCase(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"getting-started", "Getting Started"},
		{"api_reference", "Api Reference"},
		{"hello", "Hello"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := titleCase(tt.input)
			if got != tt.want {
				t.Errorf("titleCase(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
