package plugins

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/WaylonWalker/markata/pkg/buildcache"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

func TestMentionsPlugin_Name(t *testing.T) {
	p := NewMentionsPlugin()
	if got := p.Name(); got != "mentions" {
		t.Errorf("Name() = %q, want %q", got, "mentions")
	}
}

func TestMentionsPlugin_ProcessMentions(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		handleMap map[string]*mentionEntry
		want      string
	}{
		{
			name:    "simple mention",
			content: "I was reading @daverupert's post",
			handleMap: map[string]*mentionEntry{
				"daverupert": {Handle: "daverupert", SiteURL: "https://daverupert.com", Title: "Dave Rupert"},
			},
			want: `I was reading <a href="https://daverupert.com" class="mention">@daverupert</a>'s post`,
		},
		{
			name:    "mention at start of line",
			content: "@daverupert wrote about CSS",
			handleMap: map[string]*mentionEntry{
				"daverupert": {Handle: "daverupert", SiteURL: "https://daverupert.com", Title: "Dave Rupert"},
			},
			want: `<a href="https://daverupert.com" class="mention">@daverupert</a> wrote about CSS`,
		},
		{
			name:    "mention at end of line",
			content: "Great article by @daverupert",
			handleMap: map[string]*mentionEntry{
				"daverupert": {Handle: "daverupert", SiteURL: "https://daverupert.com", Title: "Dave Rupert"},
			},
			want: `Great article by <a href="https://daverupert.com" class="mention">@daverupert</a>`,
		},
		{
			name:    "multiple mentions",
			content: "Both @alice and @bob are great",
			handleMap: map[string]*mentionEntry{
				"alice": {Handle: "alice", SiteURL: "https://alice.dev", Title: "Alice"},
				"bob":   {Handle: "bob", SiteURL: "https://bob.io", Title: "Bob"},
			},
			want: `Both <a href="https://alice.dev" class="mention">@alice</a> and <a href="https://bob.io" class="mention">@bob</a> are great`,
		},
		{
			name:    "unknown mention preserved",
			content: "I follow @unknown",
			handleMap: map[string]*mentionEntry{
				"known": {Handle: "known", SiteURL: "https://known.com", Title: "Known"},
			},
			want: "I follow @unknown",
		},
		{
			name:    "mention in code block preserved",
			content: "Check this:\n```\n@daverupert\n```\nAnd @daverupert outside",
			handleMap: map[string]*mentionEntry{
				"daverupert": {Handle: "daverupert", SiteURL: "https://daverupert.com", Title: "Dave Rupert"},
			},
			want: "Check this:\n```\n@daverupert\n```\nAnd <a href=\"https://daverupert.com\" class=\"mention\">@daverupert</a> outside",
		},
		{
			name:    "email address not replaced",
			content: "Contact me at test@example.com",
			handleMap: map[string]*mentionEntry{
				"example": {Handle: "example", SiteURL: "https://example.com", Title: "Example"},
			},
			want: "Contact me at test@example.com",
		},
		{
			name:    "case insensitive handle matching",
			content: "Follow @DaveRupert",
			handleMap: map[string]*mentionEntry{
				"daverupert": {Handle: "daverupert", SiteURL: "https://daverupert.com", Title: "Dave Rupert"},
			},
			want: `Follow <a href="https://daverupert.com" class="mention">@daverupert</a>`,
		},
		{
			name:    "handle with hyphen",
			content: "Check out @dave-rupert",
			handleMap: map[string]*mentionEntry{
				"dave-rupert": {Handle: "dave-rupert", SiteURL: "https://daverupert.com", Title: "Dave Rupert"},
			},
			want: `Check out <a href="https://daverupert.com" class="mention">@dave-rupert</a>`,
		},
		{
			name:    "handle with underscore",
			content: "Check out @dave_rupert",
			handleMap: map[string]*mentionEntry{
				"dave_rupert": {Handle: "dave_rupert", SiteURL: "https://daverupert.com", Title: "Dave Rupert"},
			},
			want: `Check out <a href="https://daverupert.com" class="mention">@dave_rupert</a>`,
		},
		{
			name:      "empty handleMap",
			content:   "Hello @world",
			handleMap: map[string]*mentionEntry{},
			want:      "Hello @world",
		},
		{
			name:    "special characters escaped in URL",
			content: "Check @test",
			handleMap: map[string]*mentionEntry{
				"test": {Handle: "test", SiteURL: "https://test.com/path?a=1&b=2", Title: "Test"},
			},
			want: `Check <a href="https://test.com/path?a=1&amp;b=2" class="mention">@test</a>`,
		},
		{
			name:    "domain-style handle",
			content: "I was reading @simonwillison.net's post",
			handleMap: map[string]*mentionEntry{
				"simonwillison.net": {Handle: "simonwillison.net", SiteURL: "https://simonwillison.net", Title: "Simon Willison"},
			},
			want: `I was reading <a href="https://simonwillison.net" class="mention">@simonwillison.net</a>'s post`,
		},
		{
			name:    "domain-style handle at start",
			content: "@simonwillison.net wrote about LLMs",
			handleMap: map[string]*mentionEntry{
				"simonwillison.net": {Handle: "simonwillison.net", SiteURL: "https://simonwillison.net", Title: "Simon Willison"},
			},
			want: `<a href="https://simonwillison.net" class="mention">@simonwillison.net</a> wrote about LLMs`,
		},
		{
			name:    "domain-style handle at end",
			content: "Great article by @simonwillison.net",
			handleMap: map[string]*mentionEntry{
				"simonwillison.net": {Handle: "simonwillison.net", SiteURL: "https://simonwillison.net", Title: "Simon Willison"},
			},
			want: `Great article by <a href="https://simonwillison.net" class="mention">@simonwillison.net</a>`,
		},
		{
			name:    "handle with multiple dots",
			content: "Follow @blog.example.co.uk",
			handleMap: map[string]*mentionEntry{
				"blog.example.co.uk": {Handle: "blog.example.co.uk", SiteURL: "https://blog.example.co.uk", Title: "Example Blog"},
			},
			want: `Follow <a href="https://blog.example.co.uk" class="mention">@blog.example.co.uk</a>`,
		},
		{
			name:    "mixed handles simple and domain-style",
			content: "Both @daverupert and @simonwillison.net are great",
			handleMap: map[string]*mentionEntry{
				"daverupert":        {Handle: "daverupert", SiteURL: "https://daverupert.com", Title: "Dave Rupert"},
				"simonwillison.net": {Handle: "simonwillison.net", SiteURL: "https://simonwillison.net", Title: "Simon Willison"},
			},
			want: `Both <a href="https://daverupert.com" class="mention">@daverupert</a> and <a href="https://simonwillison.net" class="mention">@simonwillison.net</a> are great`,
		},
	}

	p := NewMentionsPlugin()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := p.processMentionsWithMetadata(tt.content, tt.handleMap)
			if got != tt.want {
				t.Errorf("processMentionsWithMetadata() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestExtractHandleFromURL(t *testing.T) {
	tests := []struct {
		url  string
		want string
	}{
		{"https://daverupert.com", "daverupert"},
		{"https://www.example.com", "example"},
		{"https://blog.jane.dev", "jane"},
		{"https://blog.example.org/feed", "example"},
		{"http://test-site.io", "test-site"},
		{"https://mysite123.com", "mysite123"},
		{"https://www.blog.example.com", "example"}, // www. is stripped first, then blog.
		{"", ""},
		{"not-a-url", ""},
	}

	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			got := extractHandleFromURL(tt.url)
			if got != tt.want {
				t.Errorf("extractHandleFromURL(%q) = %q, want %q", tt.url, got, tt.want)
			}
		})
	}
}

func TestExtractSiteURL(t *testing.T) {
	tests := []struct {
		feedURL string
		want    string
	}{
		{"https://daverupert.com/feed.xml", "https://daverupert.com"},
		{"https://example.com/blog/rss", "https://example.com"},
		{"http://test.org/atom.xml", "http://test.org"},
		{"", ""},
		{"not-a-url", ""},
	}

	for _, tt := range tests {
		t.Run(tt.feedURL, func(t *testing.T) {
			got := extractSiteURL(tt.feedURL)
			if got != tt.want {
				t.Errorf("extractSiteURL(%q) = %q, want %q", tt.feedURL, got, tt.want)
			}
		})
	}
}

func TestMentionsPlugin_LoadFromCache_UsesLastFetchedTimestamp(t *testing.T) {
	tmpDir := t.TempDir()
	p := NewMentionsPlugin()
	now := time.Now()
	metadata := &models.MentionMetadata{
		Domain:      "example.com",
		Name:        "Example",
		URL:         "https://example.com",
		LastFetched: now,
	}

	cachePath := filepath.Join(tmpDir, p.cacheKey(metadata.Domain)+".json")
	data, err := json.Marshal(metadata)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}
	if err := os.WriteFile(cachePath, data, 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	oldTime := now.Add(-14 * 24 * time.Hour)
	if err := os.Chtimes(cachePath, oldTime, oldTime); err != nil {
		t.Fatalf("Chtimes() error = %v", err)
	}

	cached := p.loadFromCache(metadata.Domain, tmpDir, 7*24*time.Hour)
	if cached == nil {
		t.Fatal("loadFromCache() = nil, want cached metadata")
	}
	if cached.Name != "Example" {
		t.Fatalf("cached.Name = %q, want Example", cached.Name)
	}
}

func TestMentionsPlugin_BuildHandleMap(t *testing.T) {
	p := NewMentionsPlugin()
	m := lifecycle.NewManager()

	// Set up blogroll config
	boolTrue := true
	config := m.Config()
	config.Extra = map[string]interface{}{
		"blogroll": models.BlogrollConfig{
			Enabled: true,
			Feeds: []models.ExternalFeedConfig{
				{
					URL:     "https://daverupert.com/feed.xml",
					Title:   "Dave Rupert",
					Handle:  "daverupert",
					SiteURL: "https://daverupert.com",
					Active:  &boolTrue,
				},
				{
					URL:     "https://example.com/rss",
					Title:   "Example Blog",
					SiteURL: "https://example.com",
					Active:  &boolTrue,
					// No explicit handle - should auto-generate from domain
				},
				{
					URL:    "https://blog.test.org/feed",
					Title:  "Test Blog",
					Active: &boolTrue,
					// No SiteURL - should extract from feed URL
				},
			},
		},
	}

	handleMap := p.buildHandleMap(m)

	// Check explicit handle
	if entry, ok := handleMap["daverupert"]; !ok {
		t.Error("expected 'daverupert' in handleMap")
	} else if entry.SiteURL != "https://daverupert.com" {
		t.Errorf("daverupert.SiteURL = %q, want %q", entry.SiteURL, "https://daverupert.com")
	}

	// Check auto-generated handle from domain
	if entry, ok := handleMap["example"]; !ok {
		t.Error("expected 'example' in handleMap (auto-generated from domain)")
	} else if entry.SiteURL != "https://example.com" {
		t.Errorf("example.SiteURL = %q, want %q", entry.SiteURL, "https://example.com")
	}

	// Check auto-generated handle from feed URL domain
	if entry, ok := handleMap["test"]; !ok {
		t.Error("expected 'test' in handleMap (auto-generated from feed URL domain)")
	} else if entry.SiteURL != "https://blog.test.org" {
		t.Errorf("test.SiteURL = %q, want %q", entry.SiteURL, "https://blog.test.org")
	}
}

func TestMentionsPlugin_Transform(t *testing.T) {
	p := NewMentionsPlugin()
	m := lifecycle.NewManager()

	// Set up blogroll config
	boolTrue := true
	config := m.Config()
	config.Extra = map[string]interface{}{
		"blogroll": models.BlogrollConfig{
			Enabled: true,
			Feeds: []models.ExternalFeedConfig{
				{
					URL:     "https://daverupert.com/feed.xml",
					Title:   "Dave Rupert",
					Handle:  "daverupert",
					SiteURL: "https://daverupert.com",
					Active:  &boolTrue,
				},
			},
		},
	}

	// Add a test post
	post := &models.Post{
		Path:    "test.md",
		Content: "I was reading @daverupert's latest post about CSS",
	}
	m.AddPost(post)

	// Run transform
	err := p.Transform(m)
	if err != nil {
		t.Fatalf("Transform() error = %v", err)
	}

	// Check the result
	posts := m.Posts()
	if len(posts) != 1 {
		t.Fatalf("expected 1 post, got %d", len(posts))
	}

	// The new implementation tries to fetch metadata, but since this is a test
	// without network access, it will likely have basic metadata or errors
	// For now, just verify that the basic link structure is present
	content := posts[0].Content
	expectedURL := `href="https://daverupert.com"`
	expectedHandle := `>@daverupert</a>`

	if !strings.Contains(content, expectedURL) {
		t.Errorf("Content missing expected URL: %q", content)
	}
	if !strings.Contains(content, expectedHandle) {
		t.Errorf("Content missing expected handle: %q", content)
	}
}

func TestMentionsPlugin_Interfaces(_ *testing.T) {
	p := NewMentionsPlugin()

	// Verify interface implementations
	var _ lifecycle.Plugin = p
	var _ lifecycle.ConfigurePlugin = p
	var _ lifecycle.TransformPlugin = p
	var _ lifecycle.PriorityPlugin = p
}

func TestMentionsPlugin_AliasResolution(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		handleMap map[string]*mentionEntry
		want      string
	}{
		{
			name:    "alias resolves to canonical handle",
			content: "I was reading @dave's post",
			handleMap: map[string]*mentionEntry{
				"daverupert": {Handle: "daverupert", SiteURL: "https://daverupert.com", Title: "Dave Rupert"},
				"dave":       {Handle: "daverupert", SiteURL: "https://daverupert.com", Title: "Dave Rupert"},
			},
			want: `I was reading <a href="https://daverupert.com" class="mention">@daverupert</a>'s post`,
		},
		{
			name:    "multiple aliases for same person",
			content: "@dave and @david are the same as @daverupert",
			handleMap: map[string]*mentionEntry{
				"daverupert": {Handle: "daverupert", SiteURL: "https://daverupert.com", Title: "Dave Rupert"},
				"dave":       {Handle: "daverupert", SiteURL: "https://daverupert.com", Title: "Dave Rupert"},
				"david":      {Handle: "daverupert", SiteURL: "https://daverupert.com", Title: "Dave Rupert"},
			},
			want: `<a href="https://daverupert.com" class="mention">@daverupert</a> and <a href="https://daverupert.com" class="mention">@daverupert</a> are the same as <a href="https://daverupert.com" class="mention">@daverupert</a>`,
		},
		{
			name:    "alias case insensitive",
			content: "Follow @DAVE",
			handleMap: map[string]*mentionEntry{
				"daverupert": {Handle: "daverupert", SiteURL: "https://daverupert.com", Title: "Dave Rupert"},
				"dave":       {Handle: "daverupert", SiteURL: "https://daverupert.com", Title: "Dave Rupert"},
			},
			want: `Follow <a href="https://daverupert.com" class="mention">@daverupert</a>`,
		},
		{
			name:    "canonical handle still works",
			content: "Check out @daverupert",
			handleMap: map[string]*mentionEntry{
				"daverupert": {Handle: "daverupert", SiteURL: "https://daverupert.com", Title: "Dave Rupert"},
				"dave":       {Handle: "daverupert", SiteURL: "https://daverupert.com", Title: "Dave Rupert"},
			},
			want: `Check out <a href="https://daverupert.com" class="mention">@daverupert</a>`,
		},
	}

	p := NewMentionsPlugin()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := p.processMentionsWithMetadata(tt.content, tt.handleMap)
			if got != tt.want {
				t.Errorf("processMentionsWithMetadata() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestMentionsPlugin_BuildHandleMapWithAliases(t *testing.T) {
	p := NewMentionsPlugin()
	m := lifecycle.NewManager()

	// Set up blogroll config with aliases
	boolTrue := true
	config := m.Config()
	config.Extra = map[string]interface{}{
		"blogroll": models.BlogrollConfig{
			Enabled: true,
			Feeds: []models.ExternalFeedConfig{
				{
					URL:     "https://daverupert.com/feed.xml",
					Title:   "Dave Rupert",
					Handle:  "daverupert",
					SiteURL: "https://daverupert.com",
					Active:  &boolTrue,
					Aliases: []string{"dave", "david", "rupert"},
				},
				{
					URL:     "https://example.com/rss",
					Title:   "Example Blog",
					Handle:  "example",
					SiteURL: "https://example.com",
					Active:  &boolTrue,
					Aliases: []string{"ex"},
				},
			},
		},
	}

	handleMap := p.buildHandleMap(m)

	// Check canonical handle
	if entry, ok := handleMap["daverupert"]; !ok {
		t.Error("expected 'daverupert' in handleMap")
	} else if entry.Handle != "daverupert" {
		t.Errorf("daverupert.Handle = %q, want %q", entry.Handle, "daverupert")
	}

	// Check aliases resolve to canonical handle
	for _, alias := range []string{"dave", "david", "rupert"} {
		if entry, ok := handleMap[alias]; !ok {
			t.Errorf("expected alias %q in handleMap", alias)
		} else if entry.Handle != "daverupert" {
			t.Errorf("alias %q resolves to Handle = %q, want %q", alias, entry.Handle, "daverupert")
		}
	}

	// Check second feed's alias
	if entry, ok := handleMap["ex"]; !ok {
		t.Error("expected alias 'ex' in handleMap")
	} else if entry.Handle != "example" {
		t.Errorf("alias 'ex' resolves to Handle = %q, want %q", entry.Handle, "example")
	}
}

func TestMentionsPlugin_DuplicateAliasFirstWins(t *testing.T) {
	p := NewMentionsPlugin()
	m := lifecycle.NewManager()

	// Set up blogroll config where two feeds have the same alias
	boolTrue := true
	config := m.Config()
	config.Extra = map[string]interface{}{
		"blogroll": models.BlogrollConfig{
			Enabled: true,
			Feeds: []models.ExternalFeedConfig{
				{
					URL:     "https://alice.com/feed.xml",
					Title:   "Alice",
					Handle:  "alice",
					SiteURL: "https://alice.com",
					Active:  &boolTrue,
					Aliases: []string{"al"}, // First feed with alias "al"
				},
				{
					URL:     "https://albert.com/rss",
					Title:   "Albert",
					Handle:  "albert",
					SiteURL: "https://albert.com",
					Active:  &boolTrue,
					Aliases: []string{"al"}, // Duplicate alias - should be ignored
				},
			},
		},
	}

	handleMap := p.buildHandleMap(m)

	// The alias "al" should resolve to alice (first entry wins)
	if entry, ok := handleMap["al"]; !ok {
		t.Error("expected alias 'al' in handleMap")
	} else if entry.Handle != "alice" {
		t.Errorf("alias 'al' resolves to Handle = %q, want %q (first entry wins)", entry.Handle, "alice")
	}
}

func TestMentionsPlugin_FromPosts(t *testing.T) {
	p := NewMentionsPlugin()
	m := lifecycle.NewManager()

	// Configure mentions with from_posts source
	config := m.Config()
	config.Extra = map[string]interface{}{
		"mentions": models.MentionsConfig{
			FromPosts: []models.MentionPostSource{
				{
					Filter:       "'contact' in tags",
					HandleField:  "handle",
					AliasesField: "aliases",
				},
			},
		},
	}

	// Add contact posts
	aliceTitle := "Alice Smith"
	post1 := &models.Post{
		Path:  "pages/contact/alice-smith.md",
		Slug:  "contact/alice-smith",
		Href:  "/contact/alice-smith/",
		Title: &aliceTitle,
		Tags:  []string{"contact", "team"},
		Extra: map[string]interface{}{
			"handle":  "alice",
			"aliases": []interface{}{"alices", "asmith"},
		},
	}

	bobTitle := "Bob Jones"
	post2 := &models.Post{
		Path:  "pages/contact/bob-jones.md",
		Slug:  "contact/bob-jones",
		Href:  "/contact/bob-jones/",
		Title: &bobTitle,
		Tags:  []string{"contact"},
		Extra: map[string]interface{}{
			"handle": "bob",
		},
	}

	// Add a non-contact post that shouldn't match
	otherTitle := "Other Post"
	post3 := &models.Post{
		Path:  "posts/other.md",
		Slug:  "other",
		Href:  "/other/",
		Title: &otherTitle,
		Tags:  []string{"blog"},
		Extra: map[string]interface{}{
			"handle": "other",
		},
	}

	m.AddPost(post1)
	m.AddPost(post2)
	m.AddPost(post3)

	handleMap := p.buildHandleMap(m)

	// Check alice's handle
	if entry, ok := handleMap["alice"]; !ok {
		t.Error("expected 'alice' in handleMap")
	} else {
		if entry.SiteURL != "/contact/alice-smith/" {
			t.Errorf("alice.SiteURL = %q, want %q", entry.SiteURL, "/contact/alice-smith/")
		}
		if entry.Title != "Alice Smith" {
			t.Errorf("alice.Title = %q, want %q", entry.Title, "Alice Smith")
		}
	}

	// Check alice's aliases
	if _, ok := handleMap["alices"]; !ok {
		t.Error("expected alias 'alices' in handleMap")
	}
	if _, ok := handleMap["asmith"]; !ok {
		t.Error("expected alias 'asmith' in handleMap")
	}

	// Check bob's handle
	if entry, ok := handleMap["bob"]; !ok {
		t.Error("expected 'bob' in handleMap")
	} else if entry.SiteURL != "/contact/bob-jones/" {
		t.Errorf("bob.SiteURL = %q, want %q", entry.SiteURL, "/contact/bob-jones/")
	}

	// Check that 'other' is NOT in the map (didn't match filter)
	if _, ok := handleMap["other"]; ok {
		t.Error("'other' should not be in handleMap (didn't match filter)")
	}
}

func TestMentionsPlugin_FromPosts_FallbackToSlug(t *testing.T) {
	p := NewMentionsPlugin()
	m := lifecycle.NewManager()

	// Configure mentions with from_posts but no handle_field
	config := m.Config()
	config.Extra = map[string]interface{}{
		"mentions": models.MentionsConfig{
			FromPosts: []models.MentionPostSource{
				{
					Filter: "'team' in tags",
					// No HandleField - should fall back to slug
				},
			},
		},
	}

	// Add a team post without explicit handle
	charlieTitle := "Charlie Brown"
	post := &models.Post{
		Path:  "team/charlie.md",
		Slug:  "charlie",
		Href:  "/team/charlie/",
		Title: &charlieTitle,
		Tags:  []string{"team"},
	}

	m.AddPost(post)

	handleMap := p.buildHandleMap(m)

	// Check that slug is used as handle
	entry, ok := handleMap["charlie"]
	if !ok {
		t.Error("expected 'charlie' (from slug) in handleMap")
	} else if entry.SiteURL != "/team/charlie/" {
		t.Errorf("charlie.SiteURL = %q, want %q", entry.SiteURL, "/team/charlie/")
	}
}

func TestMentionsPlugin_Transform_FromPosts(t *testing.T) {
	p := NewMentionsPlugin()
	m := lifecycle.NewManager()

	// Configure mentions with from_posts
	config := m.Config()
	config.Extra = map[string]interface{}{
		"mentions": models.MentionsConfig{
			CSSClass: "mention",
			FromPosts: []models.MentionPostSource{
				{
					Filter:      "'contact' in tags",
					HandleField: "handle",
				},
			},
		},
	}

	// Add a contact post
	aliceTitle := "Alice Smith"
	contactPost := &models.Post{
		Path:      "pages/contact/alice.md",
		Slug:      "contact/alice",
		Href:      "/contact/alice/",
		Title:     &aliceTitle,
		Tags:      []string{"contact"},
		Published: true,
		Extra: map[string]interface{}{
			"handle": "alice",
		},
	}

	// Add a blog post that mentions alice
	blogTitle := "Working with Alice"
	blogPost := &models.Post{
		Path:      "posts/working-with-alice.md",
		Slug:      "working-with-alice",
		Href:      "/working-with-alice/",
		Title:     &blogTitle,
		Tags:      []string{"blog"},
		Published: true,
		Content:   "I recently collaborated with @alice on a project.",
	}

	m.AddPost(contactPost)
	m.AddPost(blogPost)

	// Run transform
	err := p.Transform(m)
	if err != nil {
		t.Fatalf("Transform() error = %v", err)
	}

	// Check the blog post content was transformed
	posts := m.Posts()
	var transformedPost *models.Post
	for _, post := range posts {
		if post.Slug == "working-with-alice" {
			transformedPost = post
			break
		}
	}

	if transformedPost == nil {
		t.Fatal("could not find transformed blog post")
	}

	want := `I recently collaborated with <a href="/contact/alice/" class="mention" data-name="Alice Smith" data-handle="@alice">@alice</a> on a project.`
	if transformedPost.Content != want {
		t.Errorf("Content = %q, want %q", transformedPost.Content, want)
	}
}

func TestMentionsPlugin_Transform_CachesContent(t *testing.T) {
	p := NewMentionsPlugin()
	m := lifecycle.NewManager()
	cache := buildcache.New(t.TempDir())
	m.Cache().Set("build_cache", cache)

	config := m.Config()
	config.Extra = map[string]interface{}{
		"mentions": models.MentionsConfig{
			CSSClass: "mention",
			FromPosts: []models.MentionPostSource{
				{
					Filter:      "'contact' in tags",
					HandleField: "handle",
				},
			},
		},
	}

	aliceTitle := "Alice Smith"
	m.AddPost(&models.Post{
		Path:      "pages/contact/alice.md",
		Slug:      "contact/alice",
		Href:      "/contact/alice/",
		Title:     &aliceTitle,
		Tags:      []string{"contact"},
		Published: true,
		Extra: map[string]interface{}{
			"handle": "alice",
		},
	})

	blogPost := &models.Post{
		Path:      "posts/working-with-alice.md",
		Slug:      "working-with-alice",
		Href:      "/working-with-alice/",
		Title:     &aliceTitle,
		Tags:      []string{"blog"},
		Published: true,
		Content:   "I recently collaborated with @alice on a project.",
	}
	m.AddPost(blogPost)

	if err := p.Transform(m); err != nil {
		t.Fatalf("Transform() error = %v", err)
	}

	cachedPost, ok := cache.Posts[blogPost.Path]
	if !ok {
		t.Fatalf("expected cache entry for %s", blogPost.Path)
	}
	if cachedPost.MentionsHash == "" {
		t.Fatal("expected mentions hash to be cached")
	}
	if cachedPost.MentionsContent != blogPost.Content {
		t.Fatalf("cached mentions content = %q, want %q", cachedPost.MentionsContent, blogPost.Content)
	}
}

func TestMentionsPlugin_Transform_InvalidatesCacheWhenMentionMetadataChanges(t *testing.T) {
	cache := buildcache.New(t.TempDir())

	newManager := func(contactName string) (*lifecycle.Manager, *models.Post) {
		m := lifecycle.NewManager()
		m.Cache().Set("build_cache", cache)

		config := m.Config()
		config.Extra = map[string]interface{}{
			"mentions": models.MentionsConfig{
				CSSClass: "mention",
				FromPosts: []models.MentionPostSource{
					{
						Filter:      "'contact' in tags",
						HandleField: "handle",
					},
				},
			},
		}

		m.AddPost(&models.Post{
			Path:      "pages/contact/alice.md",
			Slug:      "contact/alice",
			Href:      "/contact/alice/",
			Title:     &contactName,
			Tags:      []string{"contact"},
			Published: true,
			Extra: map[string]interface{}{
				"handle": "alice",
			},
		})

		blogTitle := "Working with Alice"
		blogPost := &models.Post{
			Path:      "posts/working-with-alice.md",
			Slug:      "working-with-alice",
			Href:      "/working-with-alice/",
			Title:     &blogTitle,
			Tags:      []string{"blog"},
			Published: true,
			Content:   "I recently collaborated with @alice on a project.",
		}
		m.AddPost(blogPost)

		return m, blogPost
	}

	plugin := NewMentionsPlugin()
	firstManager, _ := newManager("Alice Smith")
	if err := plugin.Transform(firstManager); err != nil {
		t.Fatalf("first Transform() error = %v", err)
	}

	plugin = NewMentionsPlugin()
	secondManager, secondBlogPost := newManager("Alice Johnson")
	if err := plugin.Transform(secondManager); err != nil {
		t.Fatalf("second Transform() error = %v", err)
	}

	if !strings.Contains(secondBlogPost.Content, `data-name="Alice Johnson"`) {
		t.Fatalf("expected rebuilt content to include updated metadata, got %q", secondBlogPost.Content)
	}
	if strings.Contains(secondBlogPost.Content, `data-name="Alice Smith"`) {
		t.Fatalf("expected cache invalidation to drop stale metadata, got %q", secondBlogPost.Content)
	}
}

func TestMentionsPlugin_CombinedSources(t *testing.T) {
	p := NewMentionsPlugin()
	m := lifecycle.NewManager()

	// Configure both blogroll and from_posts
	boolTrue := true
	config := m.Config()
	config.Extra = map[string]interface{}{
		"blogroll": models.BlogrollConfig{
			Enabled: true,
			Feeds: []models.ExternalFeedConfig{
				{
					URL:     "https://external.example.com/feed.xml",
					Title:   "External Blog",
					Handle:  "external",
					SiteURL: "https://external.example.com",
					Active:  &boolTrue,
				},
			},
		},
		"mentions": models.MentionsConfig{
			FromPosts: []models.MentionPostSource{
				{
					Filter:      "'contact' in tags",
					HandleField: "handle",
				},
			},
		},
	}

	// Add an internal contact post
	aliceTitle := "Alice Smith"
	contactPost := &models.Post{
		Path:  "contact/alice.md",
		Slug:  "contact/alice",
		Href:  "/contact/alice/",
		Title: &aliceTitle,
		Tags:  []string{"contact"},
		Extra: map[string]interface{}{
			"handle": "alice",
		},
	}

	m.AddPost(contactPost)

	handleMap := p.buildHandleMap(m)

	// Check external handle from blogroll
	if entry, ok := handleMap["external"]; !ok {
		t.Error("expected 'external' from blogroll in handleMap")
	} else if entry.SiteURL != "https://external.example.com" {
		t.Errorf("external.SiteURL = %q, want %q", entry.SiteURL, "https://external.example.com")
	}

	// Check internal handle from from_posts
	if entry, ok := handleMap["alice"]; !ok {
		t.Error("expected 'alice' from from_posts in handleMap")
	} else if entry.SiteURL != "/contact/alice/" {
		t.Errorf("alice.SiteURL = %q, want %q", entry.SiteURL, "/contact/alice/")
	}
}

func TestMentionsPlugin_BuildPostMetadata(t *testing.T) {
	p := NewMentionsPlugin()

	tests := []struct {
		name       string
		post       *models.Post
		source     models.MentionPostSource
		wantName   string
		wantBio    string
		wantAvatar string
	}{
		{
			name: "all fields present with avatar",
			post: func() *models.Post {
				title := "Alice Smith"
				desc := "Software engineer"
				return &models.Post{
					Title:       &title,
					Description: &desc,
					Slug:        "alice",
					Href:        "/contact/alice/",
					Extra: map[string]interface{}{
						"avatar": "/images/alice.jpg",
					},
				}
			}(),
			source:     models.MentionPostSource{},
			wantName:   "Alice Smith",
			wantBio:    "Software engineer",
			wantAvatar: "/images/alice.jpg",
		},
		{
			name: "image field as avatar fallback",
			post: func() *models.Post {
				title := "Bob Jones"
				return &models.Post{
					Title: &title,
					Slug:  "bob",
					Href:  "/contact/bob/",
					Extra: map[string]interface{}{
						"image": "/images/bob.png",
					},
				}
			}(),
			source:     models.MentionPostSource{},
			wantName:   "Bob Jones",
			wantAvatar: "/images/bob.png",
		},
		{
			name: "icon field as avatar fallback",
			post: func() *models.Post {
				title := "Charlie"
				return &models.Post{
					Title: &title,
					Slug:  "charlie",
					Href:  "/contact/charlie/",
					Extra: map[string]interface{}{
						"icon": "/icons/charlie.svg",
					},
				}
			}(),
			source:     models.MentionPostSource{},
			wantName:   "Charlie",
			wantAvatar: "/icons/charlie.svg",
		},
		{
			name: "avatar_field config overrides default lookup",
			post: func() *models.Post {
				title := "Dana"
				return &models.Post{
					Title: &title,
					Slug:  "dana",
					Href:  "/contact/dana/",
					Extra: map[string]interface{}{
						"avatar":   "/images/dana-avatar.jpg",
						"portrait": "/images/dana-portrait.jpg",
					},
				}
			}(),
			source:     models.MentionPostSource{AvatarField: "portrait"},
			wantName:   "Dana",
			wantAvatar: "/images/dana-portrait.jpg",
		},
		{
			name: "avatar field priority: avatar over image",
			post: func() *models.Post {
				title := "Eve"
				return &models.Post{
					Title: &title,
					Slug:  "eve",
					Href:  "/contact/eve/",
					Extra: map[string]interface{}{
						"avatar": "/images/eve-avatar.jpg",
						"image":  "/images/eve-photo.jpg",
						"icon":   "/icons/eve.svg",
					},
				}
			}(),
			source:     models.MentionPostSource{},
			wantName:   "Eve",
			wantAvatar: "/images/eve-avatar.jpg",
		},
		{
			name: "no avatar fields present",
			post: func() *models.Post {
				title := "Frank"
				return &models.Post{
					Title: &title,
					Slug:  "frank",
					Href:  "/contact/frank/",
					Extra: map[string]interface{}{},
				}
			}(),
			source:     models.MentionPostSource{},
			wantName:   "Frank",
			wantAvatar: "",
		},
		{
			name: "no title falls back to slug",
			post: &models.Post{
				Slug: "ghost",
				Href: "/contact/ghost/",
			},
			source:   models.MentionPostSource{},
			wantName: "ghost",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metadata := p.buildPostMetadata(tt.post, tt.source)
			if metadata == nil {
				t.Fatal("buildPostMetadata returned nil")
			}
			if metadata.Name != tt.wantName {
				t.Errorf("Name = %q, want %q", metadata.Name, tt.wantName)
			}
			if metadata.Bio != tt.wantBio {
				t.Errorf("Bio = %q, want %q", metadata.Bio, tt.wantBio)
			}
			if metadata.Avatar != tt.wantAvatar {
				t.Errorf("Avatar = %q, want %q", metadata.Avatar, tt.wantAvatar)
			}
			if metadata.URL != tt.post.Href {
				t.Errorf("URL = %q, want %q", metadata.URL, tt.post.Href)
			}
		})
	}
}

func TestMentionsPlugin_InternalMentionDisplaysHandle(t *testing.T) {
	p := NewMentionsPlugin()

	// Simulate a mention entry with metadata (as built from internal post)
	handleMap := map[string]*mentionEntry{
		"alice": {
			Handle:  "alice",
			SiteURL: "/contact/alice/",
			Title:   "Alice Smith",
			Metadata: &models.MentionMetadata{
				Name:   "Alice Smith",
				Bio:    "Software engineer",
				Avatar: "/images/alice.jpg",
				URL:    "/contact/alice/",
			},
		},
	}

	content := "I was working with @alice on this project."
	got := p.processMentionsWithMetadata(content, handleMap)

	// Should display @handle as link text (not the title)
	wantContains := `>@alice</a>`
	if !strings.Contains(got, wantContains) {
		t.Errorf("expected @handle display text, got: %q", got)
	}

	// Should include data attributes for hovercard
	if !strings.Contains(got, `data-name="Alice Smith"`) {
		t.Error("missing data-name attribute")
	}
	if !strings.Contains(got, `data-bio="Software engineer"`) {
		t.Error("missing data-bio attribute")
	}
	if !strings.Contains(got, `data-avatar="/images/alice.jpg"`) {
		t.Error("missing data-avatar attribute")
	}
	if !strings.Contains(got, `data-handle="@alice"`) {
		t.Error("missing data-handle attribute")
	}

	// Should link to the contact page
	if !strings.Contains(got, `href="/contact/alice/"`) {
		t.Error("missing href to contact page")
	}
}

func TestMentionsPlugin_ExternalMentionDisplaysHandle(t *testing.T) {
	p := NewMentionsPlugin()

	// External mention without metadata (e.g., fetch failed or not yet fetched)
	handleMap := map[string]*mentionEntry{
		"daverupert": {
			Handle:  "daverupert",
			SiteURL: "https://daverupert.com",
			Title:   "Dave Rupert",
		},
	}

	content := "Check out @daverupert"
	got := p.processMentionsWithMetadata(content, handleMap)

	// Without metadata, should display @handle
	wantContains := `>@daverupert</a>`
	if !strings.Contains(got, wantContains) {
		t.Errorf("expected @handle display text for external mention without metadata, got: %q", got)
	}
}

func TestMentionsPlugin_FromPosts_WithMetadata(t *testing.T) {
	p := NewMentionsPlugin()
	m := lifecycle.NewManager()

	// Configure mentions with from_posts
	config := m.Config()
	config.Extra = map[string]interface{}{
		"mentions": models.MentionsConfig{
			FromPosts: []models.MentionPostSource{
				{
					Filter:      "'contact' in tags",
					HandleField: "handle",
				},
			},
		},
	}

	// Add a contact post with avatar
	aliceTitle := "Alice Smith"
	aliceDesc := "Software engineer"
	contactPost := &models.Post{
		Path:        "pages/contact/alice.md",
		Slug:        "contact/alice",
		Href:        "/contact/alice/",
		Title:       &aliceTitle,
		Description: &aliceDesc,
		Tags:        []string{"contact"},
		Extra: map[string]interface{}{
			"handle": "alice",
			"avatar": "/images/alice.jpg",
		},
	}

	m.AddPost(contactPost)

	handleMap := p.buildHandleMap(m)

	// Check that metadata is populated from post
	entry, ok := handleMap["alice"]
	if !ok {
		t.Fatal("expected 'alice' in handleMap")
	}

	if entry.Metadata == nil {
		t.Fatal("expected Metadata to be populated from post")
	}

	if entry.Metadata.Name != "Alice Smith" {
		t.Errorf("Metadata.Name = %q, want %q", entry.Metadata.Name, "Alice Smith")
	}
	if entry.Metadata.Bio != "Software engineer" {
		t.Errorf("Metadata.Bio = %q, want %q", entry.Metadata.Bio, "Software engineer")
	}
	if entry.Metadata.Avatar != "/images/alice.jpg" {
		t.Errorf("Metadata.Avatar = %q, want %q", entry.Metadata.Avatar, "/images/alice.jpg")
	}
}

func TestGetStringField(t *testing.T) {
	tests := []struct {
		name  string
		extra map[string]interface{}
		key   string
		want  string
	}{
		{"string value", map[string]interface{}{"avatar": "/img.jpg"}, "avatar", "/img.jpg"},
		{"missing key", map[string]interface{}{"other": "val"}, "avatar", ""},
		{"nil map", nil, "avatar", ""},
		{"non-string value", map[string]interface{}{"avatar": 42}, "avatar", ""},
		{"empty string", map[string]interface{}{"avatar": ""}, "avatar", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getStringField(tt.extra, tt.key)
			if got != tt.want {
				t.Errorf("getStringField() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestMentionsPlugin_TrailingPunctuation(t *testing.T) {
	p := NewMentionsPlugin()

	handleMap := map[string]*mentionEntry{
		"alice": {
			Handle:  "alice",
			SiteURL: "https://alice.dev",
			Title:   "Alice",
			Metadata: &models.MentionMetadata{
				Name: "Alice",
			},
		},
		"simonwillison.net": {
			Handle:  "simonwillison.net",
			SiteURL: "https://simonwillison.net",
			Title:   "Simon Willison",
			Metadata: &models.MentionMetadata{
				Name: "Simon Willison",
			},
		},
	}

	tests := []struct {
		name    string
		content string
		want    string
	}{
		{
			name:    "handle with trailing period",
			content: "Talk to @alice.",
			want:    `Talk to <a href="https://alice.dev" class="mention" data-name="Alice" data-handle="@alice">@alice</a>.`,
		},
		{
			name:    "handle with trailing comma",
			content: "Hey @alice, welcome!",
			want:    `Hey <a href="https://alice.dev" class="mention" data-name="Alice" data-handle="@alice">@alice</a>, welcome!`,
		},
		{
			name:    "domain handle exact match preserved",
			content: "Check @simonwillison.net for more",
			want:    `Check <a href="https://simonwillison.net" class="mention" data-name="Simon Willison" data-handle="@simonwillison.net">@simonwillison.net</a> for more`,
		},
		{
			name:    "domain handle with trailing period",
			content: "Visit @simonwillison.net.",
			want:    `Visit <a href="https://simonwillison.net" class="mention" data-name="Simon Willison" data-handle="@simonwillison.net">@simonwillison.net</a>.`,
		},
		{
			name:    "handle without punctuation unchanged",
			content: "Follow @alice on social",
			want:    `Follow <a href="https://alice.dev" class="mention" data-name="Alice" data-handle="@alice">@alice</a> on social`,
		},
		{
			name:    "unknown handle with punctuation stays plain",
			content: "Hello @unknown.",
			want:    "Hello @unknown.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := p.processMentionsWithMetadata(tt.content, handleMap)
			if got != tt.want {
				t.Errorf("processMentionsWithMetadata() =\n  %q\nwant:\n  %q", got, tt.want)
			}
		})
	}
}

func TestMentionsPlugin_RegisterAuthors(t *testing.T) {
	p := NewMentionsPlugin()

	bio := "Go developer"
	avatar := "/images/waylon.jpg"
	url := "https://waylonwalker.com"

	config := &lifecycle.Config{
		Extra: map[string]interface{}{
			"models_config": &models.Config{
				Authors: models.AuthorsConfig{
					Authors: map[string]models.Author{
						"waylon": {
							ID:     "waylon",
							Name:   "Waylon Walker",
							Bio:    &bio,
							Avatar: &avatar,
							URL:    &url,
						},
						"guest": {
							ID:   "guest",
							Name: "Guest Writer",
							// No URL — should not be registered
						},
					},
				},
			},
		},
	}

	handleMap := make(map[string]*mentionEntry)
	p.registerAuthors(config, handleMap)

	// Waylon should be registered (has URL)
	if entry, exists := handleMap["waylon"]; !exists {
		t.Error("waylon should be registered in handleMap")
	} else {
		if entry.SiteURL != url {
			t.Errorf("waylon SiteURL = %q, want %q", entry.SiteURL, url)
		}
		if entry.Metadata == nil {
			t.Fatal("waylon metadata should not be nil")
		}
		if entry.Metadata.Name != "Waylon Walker" {
			t.Errorf("waylon metadata.Name = %q, want %q", entry.Metadata.Name, "Waylon Walker")
		}
		if entry.Metadata.Bio != bio {
			t.Errorf("waylon metadata.Bio = %q, want %q", entry.Metadata.Bio, bio)
		}
		if entry.Metadata.Avatar != avatar {
			t.Errorf("waylon metadata.Avatar = %q, want %q", entry.Metadata.Avatar, avatar)
		}
	}

	// Guest should not be registered (no URL)
	if _, exists := handleMap["guest"]; exists {
		t.Error("guest should not be registered (no URL)")
	}
}

func TestMentionsPlugin_RegisterAuthors_FirstEntryWins(t *testing.T) {
	p := NewMentionsPlugin()

	url := "https://waylonwalker.com"

	config := &lifecycle.Config{
		Extra: map[string]interface{}{
			"models_config": &models.Config{
				Authors: models.AuthorsConfig{
					Authors: map[string]models.Author{
						"alice": {
							ID:   "alice",
							Name: "Alice Author",
							URL:  &url,
						},
					},
				},
			},
		},
	}

	// Pre-populate handleMap with existing entry
	handleMap := map[string]*mentionEntry{
		"alice": {Handle: "alice", SiteURL: "https://existing.com", Title: "Existing"},
	}

	p.registerAuthors(config, handleMap)

	// Existing entry should win
	if entry := handleMap["alice"]; entry.SiteURL != "https://existing.com" {
		t.Errorf("first entry should win, got SiteURL = %q", entry.SiteURL)
	}
}

func TestMentionsPlugin_ChatAdmonitionTitles(t *testing.T) {
	p := NewMentionsPlugin()

	handleMap := map[string]*mentionEntry{
		"alice": {
			Handle:  "alice",
			SiteURL: "/contact/alice/",
			Title:   "Alice Smith",
			Metadata: &models.MentionMetadata{
				Name:   "Alice Smith",
				Avatar: "/images/alice.jpg",
				Bio:    "Software engineer",
			},
		},
		"bob": {
			Handle:  "bob",
			SiteURL: "/contact/bob/",
			Title:   "Bob Jones",
			Metadata: &models.MentionMetadata{
				Name:   "Bob Jones",
				Avatar: "/images/bob.jpg",
			},
		},
	}

	avatar := "/images/waylon.jpg"
	post := &models.Post{
		Content: "",
		AuthorObjects: []models.Author{
			{
				ID:     "waylon",
				Name:   "Waylon Walker",
				Avatar: &avatar,
			},
		},
	}

	tests := []struct {
		name    string
		content string
		check   func(t *testing.T, result string)
	}{
		{
			name:    "chat with @handle gets enriched",
			content: `!!! chat "@alice"` + "\n    Hello there!",
			check: func(t *testing.T, result string) {
				t.Helper()
				if !strings.Contains(result, `chat-contact-avatar`) {
					t.Error("should contain avatar img")
				}
				if !strings.Contains(result, `class="mention"`) {
					t.Error("should contain mention link")
				}
				if !strings.Contains(result, `@alice`) {
					t.Error("should contain @alice text")
				}
				if !strings.Contains(result, `/contact/alice/`) {
					t.Error("should contain alice's URL")
				}
				if !strings.Contains(result, "\n    Hello there!") {
					t.Error("expected content to start on new line")
				}
			},
		},
		{
			name:    "chat-reply with @handle gets enriched",
			content: `!!! chat-reply "@bob"` + "\n    Great, thanks!",
			check: func(t *testing.T, result string) {
				t.Helper()
				if !strings.Contains(result, `chat-contact-avatar`) {
					t.Error("should contain avatar img for bob")
				}
				if !strings.Contains(result, `@bob`) {
					t.Error("should contain @bob text")
				}
				if !strings.Contains(result, "\n    Great, thanks!") {
					t.Error("expected content to start on new line")
				}
			},
		},
		{
			name:    "chat-reply without handle uses author",
			content: "!!! chat-reply\n    Thanks!",
			check: func(t *testing.T, result string) {
				t.Helper()
				if !strings.Contains(result, `chat-contact-name`) {
					t.Error("should contain author name span")
				}
				if !strings.Contains(result, `Waylon Walker`) {
					t.Error("should contain author name")
				}
				if !strings.Contains(result, `chat-contact-avatar`) {
					t.Error("should contain author avatar")
				}
				if !strings.Contains(result, "\n    Thanks!") {
					t.Error("expected content to start on new line")
				}
			},
		},
		{
			name:    "chat with unknown handle stays unchanged",
			content: `!!! chat "@unknown"` + "\n    Some message",
			check: func(t *testing.T, result string) {
				t.Helper()
				if strings.Contains(result, `chat-contact`) {
					t.Error("unknown handle should not be enriched")
				}
				if !strings.Contains(result, `@unknown`) {
					t.Error("should preserve original @unknown")
				}
			},
		},
		{
			name:    "non-chat admonition not affected",
			content: `!!! note "@alice"` + "\n    This is a note",
			check: func(t *testing.T, result string) {
				t.Helper()
				if strings.Contains(result, `chat-contact`) {
					t.Error("note admonition should not be affected")
				}
			},
		},
		{
			name:    "collapsible chat admonition with handle",
			content: `??? chat "@alice"` + "\n    Collapsible chat",
			check: func(t *testing.T, result string) {
				t.Helper()
				if !strings.Contains(result, `chat-contact-avatar`) {
					t.Error("collapsible chat should also get enriched")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			post.Content = tt.content
			result := p.processChatAdmonitionTitles(post, handleMap)
			tt.check(t, result)
		})
	}
}

func TestNewMentionsConfig_IncludesAuthorTemplate(t *testing.T) {
	config := models.NewMentionsConfig()

	foundContact := false
	foundAuthor := false
	for _, source := range config.FromPosts {
		if source.Filter == "template == 'contact'" {
			foundContact = true
		}
		if source.Filter == "template == 'author'" {
			foundAuthor = true
		}
	}

	if !foundContact {
		t.Error("default from_posts should include template == 'contact'")
	}
	if !foundAuthor {
		t.Error("default from_posts should include template == 'author'")
	}
}

func TestMentionsPlugin_ChatAdmonitionTitles_Unquoted(t *testing.T) {
	p := NewMentionsPlugin()

	handleMap := map[string]*mentionEntry{
		"alice": {
			Handle:  "alice",
			SiteURL: "/contact/alice/",
			Title:   "Alice Smith",
			Metadata: &models.MentionMetadata{
				Name:   "Alice Smith",
				Avatar: "/images/alice.jpg",
				Bio:    "Software engineer",
			},
		},
		"bob": {
			Handle:  "bob",
			SiteURL: "/contact/bob/",
			Title:   "Bob Jones",
			Metadata: &models.MentionMetadata{
				Name:   "Bob Jones",
				Avatar: "/images/bob.jpg",
			},
		},
	}

	avatar := "/images/waylon.jpg"
	post := &models.Post{
		Content: "",
		AuthorObjects: []models.Author{
			{
				ID:     "waylon",
				Name:   "Waylon Walker",
				Avatar: &avatar,
			},
		},
	}

	tests := []struct {
		name    string
		content string
		check   func(t *testing.T, result string)
	}{
		{
			name:    "unquoted chat @handle gets enriched",
			content: "!!! chat @alice\n    Hello there!",
			check: func(t *testing.T, result string) {
				t.Helper()
				if !strings.Contains(result, `chat-contact-avatar`) {
					t.Error("should contain avatar img")
				}
				if !strings.Contains(result, `class="mention"`) {
					t.Error("should contain mention link")
				}
				if !strings.Contains(result, `@alice`) {
					t.Error("should contain @alice text")
				}
				// Enriched title should NOT be wrapped in quotes (quotes break goldmark parsing)
				if strings.Contains(result, `!!! chat "`) {
					t.Error("enriched title should not be wrapped in quotes")
				}
				if !strings.Contains(result, "\n    Hello there!") {
					t.Error("expected content to start on new line")
				}
			},
		},
		{
			name:    "unquoted chat-reply @handle gets enriched",
			content: "!!! chat-reply @bob\n    Thanks!",
			check: func(t *testing.T, result string) {
				t.Helper()
				if !strings.Contains(result, `chat-contact-avatar`) {
					t.Error("should contain avatar img for bob")
				}
				if !strings.Contains(result, `@bob`) {
					t.Error("should contain @bob text")
				}
				if !strings.Contains(result, "\n    Thanks!") {
					t.Error("expected content to start on new line")
				}
			},
		},
		{
			name:    "unquoted chat-reply without handle uses author",
			content: "!!! chat-reply\n    Thanks!",
			check: func(t *testing.T, result string) {
				t.Helper()
				if !strings.Contains(result, `chat-contact-name`) {
					t.Error("should contain author name span")
				}
				if !strings.Contains(result, `Waylon Walker`) {
					t.Error("should contain author name")
				}
				if !strings.Contains(result, `chat-contact-avatar`) {
					t.Error("should contain author avatar")
				}
				if !strings.Contains(result, "\n    Thanks!") {
					t.Error("expected content to start on new line")
				}
			},
		},
		{
			name:    "unquoted collapsible chat @handle gets enriched",
			content: "??? chat @alice\n    Collapsible",
			check: func(t *testing.T, result string) {
				t.Helper()
				if !strings.Contains(result, `chat-contact-avatar`) {
					t.Error("collapsible unquoted chat should also get enriched")
				}
			},
		},
		{
			name:    "unquoted expanded collapsible chat @handle gets enriched",
			content: "???+ chat @alice\n    Expanded collapsible",
			check: func(t *testing.T, result string) {
				t.Helper()
				if !strings.Contains(result, `chat-contact-avatar`) {
					t.Error("expanded collapsible unquoted chat should also get enriched")
				}
			},
		},
		{
			name:    "unquoted unknown handle stays unchanged",
			content: "!!! chat @unknown\n    Message",
			check: func(t *testing.T, result string) {
				t.Helper()
				if strings.Contains(result, `chat-contact`) {
					t.Error("unknown handle should not be enriched")
				}
				if result != "!!! chat @unknown\n    Message" {
					t.Errorf("should preserve original content, got: %q", result)
				}
			},
		},
		{
			name:    "unquoted non-chat admonition not affected",
			content: "!!! note @alice\n    This is a note",
			check: func(t *testing.T, result string) {
				t.Helper()
				if strings.Contains(result, `chat-contact`) {
					t.Error("note admonition should not be affected")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			post.Content = tt.content
			result := p.processChatAdmonitionTitles(post, handleMap)
			tt.check(t, result)
		})
	}
}

func TestMentionsPlugin_AdmonitionLinesProtected(t *testing.T) {
	p := NewMentionsPlugin()

	handleMap := map[string]*mentionEntry{
		"alice": {
			Handle:  "alice",
			SiteURL: "/contact/alice/",
			Title:   "Alice Smith",
			Metadata: &models.MentionMetadata{
				Name: "Alice Smith",
			},
		},
	}

	tests := []struct {
		name    string
		content string
		check   func(t *testing.T, result string)
	}{
		{
			name:    "admonition header line @mention not linkified",
			content: "!!! note @alice\n    Some note content",
			check: func(t *testing.T, result string) {
				t.Helper()
				// The admonition header line should NOT have the @alice transformed
				if strings.Contains(result, `<a href=`) {
					t.Error("@mention on admonition header line should not be linkified")
				}
				if !strings.Contains(result, "!!! note @alice") {
					t.Error("admonition header should be preserved")
				}
			},
		},
		{
			name:    "mention in body after admonition is still processed",
			content: "!!! note \"Title\"\n    Content mentioning @alice here.\n\nText outside with @alice too",
			check: func(t *testing.T, result string) {
				t.Helper()
				// The body content and text outside should have mentions processed
				// But the mention inside an indented admonition body is fine (it's not a header line)
				if !strings.Contains(result, `<a href="/contact/alice/"`) {
					t.Error("@mention outside admonition header should still be linkified")
				}
			},
		},
		{
			name:    "collapsible admonition header protected",
			content: "??? warning @alice\n    Warning content",
			check: func(t *testing.T, result string) {
				t.Helper()
				if strings.Contains(result, `<a href=`) {
					t.Error("@mention on ??? header should not be linkified")
				}
			},
		},
		{
			name:    "expanded collapsible admonition header protected",
			content: "???+ info @alice\n    Info content",
			check: func(t *testing.T, result string) {
				t.Helper()
				if strings.Contains(result, `<a href=`) {
					t.Error("@mention on ???+ header should not be linkified")
				}
			},
		},
		{
			name:    "quoted admonition title with @handle protected",
			content: `!!! tip "@alice"` + "\n    Tip content",
			check: func(t *testing.T, result string) {
				t.Helper()
				if strings.Contains(result, `<a href=`) {
					t.Error("@mention in quoted admonition title should not be linkified by general pass")
				}
			},
		},
		{
			name:    "multiple admonitions with mentions between",
			content: "!!! chat @alice\n    Chat content\n\nHello @alice!\n\n!!! note @alice\n    Note content",
			check: func(t *testing.T, result string) {
				t.Helper()
				// The text between admonitions should have the mention linkified
				if !strings.Contains(result, `Hello <a href="/contact/alice/"`) {
					t.Error("mention between admonitions should be linkified")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := p.processMentionsWithMetadata(tt.content, handleMap)
			tt.check(t, result)
		})
	}
}

func TestMentionsPlugin_NoDoubleTransformation(t *testing.T) {
	p := NewMentionsPlugin()

	handleMap := map[string]*mentionEntry{
		"alice": {
			Handle:  "alice",
			SiteURL: "/contact/alice/",
			Title:   "Alice Smith",
			Metadata: &models.MentionMetadata{
				Name:   "Alice Smith",
				Avatar: "/images/alice.jpg",
				Bio:    "Software engineer",
			},
		},
	}

	tests := []struct {
		name    string
		content string
		check   func(t *testing.T, result string)
	}{
		{
			name:    "enriched chat title not re-processed by general mention pass",
			content: `!!! chat "@alice"` + "\n    Hello there!\n\nSome text mentioning @alice.",
			check: func(t *testing.T, result string) {
				t.Helper()
				// First, enrich the chat title
				enriched := p.processChatAdmonitionTitles(&models.Post{Content: result}, handleMap)
				// Then run the general mention pass
				final := p.processMentionsWithMetadata(enriched, handleMap)

				// The chat title line should NOT have nested <a> tags
				lines := strings.Split(final, "\n")
				chatLine := lines[0]
				// Count <a href occurrences in the chat title line
				aCount := strings.Count(chatLine, `<a href=`)
				if aCount > 1 {
					t.Errorf("chat title line has %d <a> tags (double transformation), got: %q", aCount, chatLine)
				}

				// The mention in the body text should still be linkified
				if !strings.Contains(final, `Some text mentioning <a href="/contact/alice/"`) {
					t.Error("mention in body text should still be linkified")
				}
			},
		},
		{
			name:    "unquoted enriched chat title not re-processed",
			content: "!!! chat @alice\n    Hello there!\n\nSome text mentioning @alice.",
			check: func(t *testing.T, result string) {
				t.Helper()
				enriched := p.processChatAdmonitionTitles(&models.Post{Content: result}, handleMap)
				final := p.processMentionsWithMetadata(enriched, handleMap)

				// The chat title line should NOT have nested <a> tags
				lines := strings.Split(final, "\n")
				chatLine := lines[0]
				aCount := strings.Count(chatLine, `<a href=`)
				if aCount > 1 {
					t.Errorf("unquoted chat title line has %d <a> tags (double transformation), got: %q", aCount, chatLine)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.check(t, tt.content)
		})
	}
}
