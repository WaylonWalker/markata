package plugins

import (
	"testing"

	"github.com/WaylonWalker/markata/pkg/buildcache"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

func TestLinkCollectorPlugin_Name(t *testing.T) {
	p := NewLinkCollectorPlugin()
	if p.Name() != "link_collector" {
		t.Errorf("expected name 'link_collector', got %q", p.Name())
	}
}

func TestLinkCollectorPlugin_Configure(t *testing.T) {
	p := NewLinkCollectorPlugin()
	m := lifecycle.NewManager()

	// Test default configuration
	if err := p.Configure(m); err != nil {
		t.Errorf("Configure returned error: %v", err)
	}
	if p.includeFeeds {
		t.Error("expected includeFeeds to be false by default")
	}
	if p.includeIndex {
		t.Error("expected includeIndex to be false by default")
	}

	// Test with custom configuration
	config := m.Config()
	config.Extra = map[string]interface{}{
		"link_collector": map[string]interface{}{
			"include_feeds": true,
			"include_index": true,
		},
		"url": "https://example.com",
	}
	m.SetConfig(config)

	p2 := NewLinkCollectorPlugin()
	if err := p2.Configure(m); err != nil {
		t.Errorf("Configure returned error: %v", err)
	}
	if !p2.includeFeeds {
		t.Error("expected includeFeeds to be true")
	}
	if !p2.includeIndex {
		t.Error("expected includeIndex to be true")
	}
	if p2.siteURL != "https://example.com" {
		t.Errorf("expected siteURL 'https://example.com', got %q", p2.siteURL)
	}
	if p2.siteDomain != "example.com" {
		t.Errorf("expected siteDomain 'example.com', got %q", p2.siteDomain)
	}
}

func TestLinkCollectorPlugin_Priority(t *testing.T) {
	p := NewLinkCollectorPlugin()

	// Should have late priority for render stage
	if p.Priority(lifecycle.StageRender) != lifecycle.PriorityLate {
		t.Errorf("expected PriorityLate for render stage, got %d", p.Priority(lifecycle.StageRender))
	}

	// Should have default priority for other stages
	if p.Priority(lifecycle.StageTransform) != lifecycle.PriorityDefault {
		t.Errorf("expected PriorityDefault for transform stage, got %d", p.Priority(lifecycle.StageTransform))
	}
}

func TestExtractHrefs(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		expected []string
	}{
		{
			name:     "single link",
			html:     `<p>Visit <a href="https://example.com">Example</a></p>`,
			expected: []string{"https://example.com"},
		},
		{
			name:     "multiple links",
			html:     `<a href="/page1">Page 1</a> and <a href="/page2">Page 2</a>`,
			expected: []string{"/page1", "/page2"},
		},
		{
			name:     "duplicate links",
			html:     `<a href="/page1">Link 1</a> <a href="/page1">Link 2</a>`,
			expected: []string{"/page1"},
		},
		{
			name:     "empty href",
			html:     `<a href="">Empty</a>`,
			expected: []string{},
		},
		{
			name:     "anchor only",
			html:     `<a href="#">Anchor</a>`,
			expected: []string{},
		},
		{
			name:     "no links",
			html:     `<p>No links here</p>`,
			expected: []string{},
		},
		{
			name:     "mixed links",
			html:     `<a href="https://external.com">External</a> <a href="/internal">Internal</a> <a href="relative">Relative</a>`,
			expected: []string{"https://external.com", "/internal", "relative"},
		},
		{
			name:     "link with attributes",
			html:     `<a class="btn" href="/page" target="_blank">Page</a>`,
			expected: []string{"/page"},
		},
		{
			name:     "single quoted href",
			html:     `<a href='/page'>Page</a>`,
			expected: []string{"/page"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractHrefs(tt.html)
			if len(got) != len(tt.expected) {
				t.Errorf("extractHrefs() = %v, want %v", got, tt.expected)
				return
			}
			for i, href := range got {
				if href != tt.expected[i] {
					t.Errorf("extractHrefs()[%d] = %q, want %q", i, href, tt.expected[i])
				}
			}
		})
	}
}

func TestResolveURL(t *testing.T) {
	tests := []struct {
		name     string
		baseURL  string
		href     string
		expected string
	}{
		{
			name:     "absolute https",
			baseURL:  "https://mysite.com/post/",
			href:     "https://example.com/page",
			expected: "https://example.com/page",
		},
		{
			name:     "absolute http",
			baseURL:  "https://mysite.com/post/",
			href:     "http://example.com/page",
			expected: "http://example.com/page",
		},
		{
			name:     "protocol relative",
			baseURL:  "https://mysite.com/post/",
			href:     "//example.com/page",
			expected: "https://example.com/page",
		},
		{
			name:     "root relative",
			baseURL:  "https://mysite.com/post/",
			href:     "/other-post/",
			expected: "https://mysite.com/other-post/",
		},
		{
			name:     "relative",
			baseURL:  "https://mysite.com/post/",
			href:     "sibling",
			expected: "https://mysite.com/post/sibling",
		},
		{
			name:     "parent relative",
			baseURL:  "https://mysite.com/post/sub/",
			href:     "../other",
			expected: "https://mysite.com/post/other",
		},
		{
			name:     "anchor",
			baseURL:  "https://mysite.com/post/",
			href:     "#section",
			expected: "https://mysite.com/post/#section",
		},
		{
			name:     "mailto",
			baseURL:  "https://mysite.com/post/",
			href:     "mailto:test@example.com",
			expected: "mailto:test@example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := resolveURL(tt.baseURL, tt.href)
			if got != tt.expected {
				t.Errorf("resolveURL(%q, %q) = %q, want %q", tt.baseURL, tt.href, got, tt.expected)
			}
		})
	}
}

func TestLinkCollectorPlugin_Render_BasicLinks(t *testing.T) {
	p := NewLinkCollectorPlugin()
	p.SetSiteURL("https://example.com")
	m := lifecycle.NewManager()

	title1 := "Post One"
	title2 := "Post Two"

	posts := []*models.Post{
		{
			Slug:        "post-one",
			Href:        "/post-one/",
			Title:       &title1,
			ArticleHTML: `<p>Link to <a href="/post-two/">Post Two</a></p>`,
		},
		{
			Slug:        "post-two",
			Href:        "/post-two/",
			Title:       &title2,
			ArticleHTML: `<p>Link to <a href="/post-one/">Post One</a></p>`,
		},
	}
	m.SetPosts(posts)

	err := p.Render(m)
	if err != nil {
		t.Fatalf("Render error: %v", err)
	}

	// Check post-one
	post1 := m.Posts()[0]
	if len(post1.Hrefs) != 1 {
		t.Errorf("post-one: expected 1 href, got %d", len(post1.Hrefs))
	}
	if len(post1.Outlinks) != 1 {
		t.Errorf("post-one: expected 1 outlink, got %d", len(post1.Outlinks))
	}
	if len(post1.Inlinks) != 1 {
		t.Errorf("post-one: expected 1 inlink, got %d", len(post1.Inlinks))
	}

	// Verify outlink properties
	if len(post1.Outlinks) > 0 {
		outlink := post1.Outlinks[0]
		if outlink.TargetPost == nil {
			t.Error("post-one outlink: expected target post to be set")
		} else if outlink.TargetPost.Slug != "post-two" {
			t.Errorf("post-one outlink: expected target slug 'post-two', got %q", outlink.TargetPost.Slug)
		}
		if !outlink.IsInternal {
			t.Error("post-one outlink: expected IsInternal to be true")
		}
		if outlink.IsSelf {
			t.Error("post-one outlink: expected IsSelf to be false")
		}
	}

	// Verify inlink properties
	if len(post1.Inlinks) > 0 {
		inlink := post1.Inlinks[0]
		if inlink.SourcePost == nil {
			t.Error("post-one inlink: expected source post to be set")
		} else if inlink.SourcePost.Slug != "post-two" {
			t.Errorf("post-one inlink: expected source slug 'post-two', got %q", inlink.SourcePost.Slug)
		}
	}

	// Check post-two
	post2 := m.Posts()[1]
	if len(post2.Outlinks) != 1 {
		t.Errorf("post-two: expected 1 outlink, got %d", len(post2.Outlinks))
	}
	if len(post2.Inlinks) != 1 {
		t.Errorf("post-two: expected 1 inlink, got %d", len(post2.Inlinks))
	}
}

func TestLinkCollectorPlugin_Render_ExternalLinks(t *testing.T) {
	p := NewLinkCollectorPlugin()
	p.SetSiteURL("https://mysite.com")
	m := lifecycle.NewManager()

	title := "My Post"
	posts := []*models.Post{
		{
			Slug:        "my-post",
			Href:        "/my-post/",
			Title:       &title,
			ArticleHTML: `<p>Check out <a href="https://external.com/page">External Site</a></p>`,
		},
	}
	m.SetPosts(posts)

	err := p.Render(m)
	if err != nil {
		t.Fatalf("Render error: %v", err)
	}

	post := m.Posts()[0]
	if len(post.Outlinks) != 1 {
		t.Fatalf("expected 1 outlink, got %d", len(post.Outlinks))
	}

	outlink := post.Outlinks[0]
	if outlink.IsInternal {
		t.Error("expected IsInternal to be false for external link")
	}
	if outlink.TargetDomain != "external.com" {
		t.Errorf("expected target domain 'external.com', got %q", outlink.TargetDomain)
	}
	if outlink.TargetPost != nil {
		t.Error("expected target post to be nil for external link")
	}
}

func TestLinkCollectorPlugin_Render_SelfLinks(t *testing.T) {
	p := NewLinkCollectorPlugin()
	p.SetSiteURL("https://example.com")
	m := lifecycle.NewManager()

	title := "My Post"
	posts := []*models.Post{
		{
			Slug:        "my-post",
			Href:        "/my-post/",
			Title:       &title,
			ArticleHTML: `<p>See <a href="#section">this section</a> and <a href="/my-post/#other">another</a></p>`,
		},
	}
	m.SetPosts(posts)

	err := p.Render(m)
	if err != nil {
		t.Fatalf("Render error: %v", err)
	}

	post := m.Posts()[0]
	// Self-links should not be in outlinks
	for _, link := range post.Outlinks {
		if link.IsSelf {
			t.Error("self-link should not be in outlinks")
		}
	}
	// Self-links should not be in inlinks
	if len(post.Inlinks) != 0 {
		t.Errorf("expected 0 inlinks for self-referencing post, got %d", len(post.Inlinks))
	}
}

func TestLinkCollectorPlugin_Render_SkippedPosts(t *testing.T) {
	p := NewLinkCollectorPlugin()
	m := lifecycle.NewManager()

	posts := []*models.Post{
		{
			Slug:        "skipped",
			Href:        "/skipped/",
			Skip:        true,
			ArticleHTML: `<p><a href="/other/">Link</a></p>`,
		},
	}
	m.SetPosts(posts)

	err := p.Render(m)
	if err != nil {
		t.Fatalf("Render error: %v", err)
	}

	post := m.Posts()[0]
	if len(post.Hrefs) != 0 {
		t.Errorf("expected 0 hrefs for skipped post, got %d", len(post.Hrefs))
	}
}

func TestLinkCollectorPlugin_Render_EmptyHTML(t *testing.T) {
	p := NewLinkCollectorPlugin()
	m := lifecycle.NewManager()

	posts := []*models.Post{
		{
			Slug:        "empty",
			Href:        "/empty/",
			ArticleHTML: "",
		},
	}
	m.SetPosts(posts)

	err := p.Render(m)
	if err != nil {
		t.Fatalf("Render error: %v", err)
	}

	post := m.Posts()[0]
	if len(post.Hrefs) != 0 {
		t.Errorf("expected 0 hrefs for empty post, got %d", len(post.Hrefs))
	}
}

func TestLinkCollectorPlugin_Render_DeduplicateInlinks(t *testing.T) {
	p := NewLinkCollectorPlugin()
	p.SetSiteURL("https://example.com")
	m := lifecycle.NewManager()

	title1 := "Post One"
	title2 := "Post Two"

	posts := []*models.Post{
		{
			Slug:        "post-one",
			Href:        "/post-one/",
			Title:       &title1,
			ArticleHTML: `<p><a href="/post-two/">Link 1</a> and <a href="/post-two/">Link 2</a></p>`,
		},
		{
			Slug:        "post-two",
			Href:        "/post-two/",
			Title:       &title2,
			ArticleHTML: `<p>Content</p>`,
		},
	}
	m.SetPosts(posts)

	err := p.Render(m)
	if err != nil {
		t.Fatalf("Render error: %v", err)
	}

	post2 := m.Posts()[1]
	// Should have only 1 inlink (deduplicated by source URL)
	if len(post2.Inlinks) != 1 {
		t.Errorf("expected 1 inlink (deduplicated), got %d", len(post2.Inlinks))
	}
}

func TestLinkCollectorPlugin_Render_ExcludeIndex(t *testing.T) {
	p := NewLinkCollectorPlugin()
	p.SetSiteURL("https://example.com")
	p.SetIncludeIndex(false)
	m := lifecycle.NewManager()

	indexTitle := "Home"
	postTitle := "My Post"

	posts := []*models.Post{
		{
			Slug:        "index",
			Href:        "/index/",
			Title:       &indexTitle,
			ArticleHTML: `<p><a href="/my-post/">My Post</a></p>`,
		},
		{
			Slug:        "my-post",
			Href:        "/my-post/",
			Title:       &postTitle,
			ArticleHTML: `<p>Content</p>`,
		},
	}
	m.SetPosts(posts)

	err := p.Render(m)
	if err != nil {
		t.Fatalf("Render error: %v", err)
	}

	post := m.Posts()[1]
	// Index should be excluded from inlinks
	if len(post.Inlinks) != 0 {
		t.Errorf("expected 0 inlinks (index excluded), got %d", len(post.Inlinks))
	}

	// Now include index
	p2 := NewLinkCollectorPlugin()
	p2.SetSiteURL("https://example.com")
	p2.SetIncludeIndex(true)
	m2 := lifecycle.NewManager()
	m2.SetPosts(posts)

	err = p2.Render(m2)
	if err != nil {
		t.Fatalf("Render error: %v", err)
	}

	post2 := m2.Posts()[1]
	if len(post2.Inlinks) != 1 {
		t.Errorf("expected 1 inlink (index included), got %d", len(post2.Inlinks))
	}
}

func TestLinkCollectorPlugin_Render_ExcludeFeeds(t *testing.T) {
	p := NewLinkCollectorPlugin()
	p.SetSiteURL("https://example.com")
	p.SetIncludeFeeds(false)
	m := lifecycle.NewManager()

	feedTitle := "Blog Feed"
	postTitle := "My Post"

	posts := []*models.Post{
		{
			Slug:        "blog",
			Href:        "/blog/",
			Title:       &feedTitle,
			Template:    "feed.html",
			ArticleHTML: `<p><a href="/my-post/">My Post</a></p>`,
		},
		{
			Slug:        "my-post",
			Href:        "/my-post/",
			Title:       &postTitle,
			ArticleHTML: `<p>Content</p>`,
		},
	}
	m.SetPosts(posts)

	err := p.Render(m)
	if err != nil {
		t.Fatalf("Render error: %v", err)
	}

	post := m.Posts()[1]
	// Feed should be excluded from inlinks
	if len(post.Inlinks) != 0 {
		t.Errorf("expected 0 inlinks (feed excluded), got %d", len(post.Inlinks))
	}

	// Now include feeds
	p2 := NewLinkCollectorPlugin()
	p2.SetSiteURL("https://example.com")
	p2.SetIncludeFeeds(true)
	m2 := lifecycle.NewManager()
	m2.SetPosts(posts)

	err = p2.Render(m2)
	if err != nil {
		t.Fatalf("Render error: %v", err)
	}

	post2 := m2.Posts()[1]
	if len(post2.Inlinks) != 1 {
		t.Errorf("expected 1 inlink (feed included), got %d", len(post2.Inlinks))
	}
}

func TestLinkCollectorPlugin_Render_MultiplePosts(t *testing.T) {
	p := NewLinkCollectorPlugin()
	p.SetSiteURL("https://example.com")
	m := lifecycle.NewManager()

	titleA := "Post A"
	titleB := "Post B"
	titleC := "Post C"

	posts := []*models.Post{
		{
			Slug:        "post-a",
			Href:        "/post-a/",
			Title:       &titleA,
			ArticleHTML: `<p><a href="/post-b/">B</a> <a href="/post-c/">C</a></p>`,
		},
		{
			Slug:        "post-b",
			Href:        "/post-b/",
			Title:       &titleB,
			ArticleHTML: `<p><a href="/post-a/">A</a></p>`,
		},
		{
			Slug:        "post-c",
			Href:        "/post-c/",
			Title:       &titleC,
			ArticleHTML: `<p><a href="/post-a/">A</a> <a href="/post-b/">B</a></p>`,
		},
	}
	m.SetPosts(posts)

	err := p.Render(m)
	if err != nil {
		t.Fatalf("Render error: %v", err)
	}

	resultPosts := m.Posts()

	// Post A: outlinks to B and C, inlinks from B and C
	postA := resultPosts[0]
	if len(postA.Outlinks) != 2 {
		t.Errorf("post-a: expected 2 outlinks, got %d", len(postA.Outlinks))
	}
	if len(postA.Inlinks) != 2 {
		t.Errorf("post-a: expected 2 inlinks, got %d", len(postA.Inlinks))
	}

	// Post B: outlinks to A, inlinks from A and C
	postB := resultPosts[1]
	if len(postB.Outlinks) != 1 {
		t.Errorf("post-b: expected 1 outlink, got %d", len(postB.Outlinks))
	}
	if len(postB.Inlinks) != 2 {
		t.Errorf("post-b: expected 2 inlinks, got %d", len(postB.Inlinks))
	}

	// Post C: outlinks to A and B, inlinks from A
	postC := resultPosts[2]
	if len(postC.Outlinks) != 2 {
		t.Errorf("post-c: expected 2 outlinks, got %d", len(postC.Outlinks))
	}
	if len(postC.Inlinks) != 1 {
		t.Errorf("post-c: expected 1 inlink, got %d", len(postC.Inlinks))
	}
}

func TestLinkCollectorPlugin_Render_LinksStoredInCache(t *testing.T) {
	p := NewLinkCollectorPlugin()
	p.SetSiteURL("https://example.com")
	m := lifecycle.NewManager()

	title := "My Post"
	posts := []*models.Post{
		{
			Slug:        "my-post",
			Href:        "/my-post/",
			Title:       &title,
			ArticleHTML: `<p><a href="https://external.com">External</a></p>`,
		},
	}
	m.SetPosts(posts)

	err := p.Render(m)
	if err != nil {
		t.Fatalf("Render error: %v", err)
	}

	// Check that links are stored in cache
	links, ok := m.Cache().Get("links")
	if !ok {
		t.Error("expected links to be stored in cache")
	}
	if allLinks, ok := links.([]*models.Link); ok {
		if len(allLinks) != 1 {
			t.Errorf("expected 1 link in cache, got %d", len(allLinks))
		}
	} else {
		t.Error("links in cache have unexpected type")
	}
}

func TestLinkCollectorPlugin_Render_CachesAndInvalidatesHrefs(t *testing.T) {
	p := NewLinkCollectorPlugin()
	p.SetSiteURL("https://example.com")

	cache := buildcache.New(t.TempDir())
	m := lifecycle.NewManager()
	m.Cache().Set("build_cache", cache)

	posts := []*models.Post{
		{
			Path:        "posts/one.md",
			Slug:        "one",
			Href:        "/one/",
			ArticleHTML: `<p><a href="https://example.com/a">A</a></p>`,
		},
	}
	m.SetPosts(posts)

	if err := p.Render(m); err != nil {
		t.Fatalf("Render error: %v", err)
	}

	if len(posts[0].Hrefs) != 1 || posts[0].Hrefs[0] != "https://example.com/a" {
		t.Fatalf("expected initial hrefs, got %v", posts[0].Hrefs)
	}

	posts[0].ArticleHTML = `<p><a href="https://example.com/b">B</a></p>`
	posts[0].Hrefs = nil

	if err := p.Render(m); err != nil {
		t.Fatalf("Render error: %v", err)
	}

	if len(posts[0].Hrefs) != 1 || posts[0].Hrefs[0] != "https://example.com/b" {
		t.Fatalf("expected updated hrefs after change, got %v", posts[0].Hrefs)
	}
}

func TestDeduplicateLinksBySource(t *testing.T) {
	links := []*models.Link{
		{SourceURL: "https://example.com/post-1/", TargetURL: "https://example.com/target/"},
		{SourceURL: "https://example.com/post-1/", TargetURL: "https://example.com/target/"}, // duplicate
		{SourceURL: "https://example.com/post-2/", TargetURL: "https://example.com/target/"},
	}

	result := deduplicateLinksBySource(links)
	if len(result) != 2 {
		t.Errorf("expected 2 deduplicated links, got %d", len(result))
	}
}

func TestDeduplicateLinksByTarget(t *testing.T) {
	links := []*models.Link{
		{SourceURL: "https://example.com/post/", TargetURL: "https://example.com/target-1/"},
		{SourceURL: "https://example.com/post/", TargetURL: "https://example.com/target-1/"}, // duplicate
		{SourceURL: "https://example.com/post/", TargetURL: "https://example.com/target-2/"},
	}

	result := deduplicateLinksByTarget(links)
	if len(result) != 2 {
		t.Errorf("expected 2 deduplicated links, got %d", len(result))
	}
}

func TestDeduplicateLinksBySource_Empty(t *testing.T) {
	result := deduplicateLinksBySource(nil)
	if result != nil {
		t.Error("expected nil for empty input")
	}

	result = deduplicateLinksBySource([]*models.Link{})
	if result != nil {
		t.Error("expected nil for empty slice")
	}
}

func TestDeduplicateLinksByTarget_Empty(t *testing.T) {
	result := deduplicateLinksByTarget(nil)
	if result != nil {
		t.Error("expected nil for empty input")
	}

	result = deduplicateLinksByTarget([]*models.Link{})
	if result != nil {
		t.Error("expected nil for empty slice")
	}
}

func TestIsFeedPost(t *testing.T) {
	tests := []struct {
		name     string
		post     *models.Post
		expected bool
	}{
		{
			name:     "feed template",
			post:     &models.Post{Template: "feed.html"},
			expected: true,
		},
		{
			name:     "archive template",
			post:     &models.Post{Template: "archive.html"},
			expected: true,
		},
		{
			name:     "post template",
			post:     &models.Post{Template: "post.html"},
			expected: false,
		},
		{
			name: "is_feed extra",
			post: &models.Post{
				Template: "post.html",
				Extra:    map[string]interface{}{"is_feed": true},
			},
			expected: true,
		},
		{
			name: "is_feed false",
			post: &models.Post{
				Template: "post.html",
				Extra:    map[string]interface{}{"is_feed": false},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isFeedPost(tt.post)
			if got != tt.expected {
				t.Errorf("isFeedPost() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// Interface compliance tests

func TestLinkCollectorPlugin_Interfaces(_ *testing.T) {
	p := NewLinkCollectorPlugin()

	// Verify interface compliance
	var _ lifecycle.Plugin = p
	var _ lifecycle.ConfigurePlugin = p
	var _ lifecycle.RenderPlugin = p
	var _ lifecycle.PriorityPlugin = p
}
