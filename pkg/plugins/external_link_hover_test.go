package plugins

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

func newExternalLinkHoverManager(t *testing.T, hover, extra map[string]any) (manager *lifecycle.Manager, cacheDir string) {
	t.Helper()
	cacheDir = filepath.Join(t.TempDir(), "embeds")
	cfgExtra := map[string]interface{}{
		"models_config":       &models.Config{URL: "https://mysite.test"},
		"external_link_hover": hover,
		"embeds": map[string]interface{}{
			"cache_dir":      cacheDir,
			"oembed_enabled": false,
		},
	}
	for k, v := range extra {
		cfgExtra[k] = v
	}
	m := lifecycle.NewManager()
	m.SetConfig(&lifecycle.Config{Extra: cfgExtra})
	return m, cacheDir
}

func seedEmbedsCache(t *testing.T, cacheDir, rawURL string, meta *OGMetadata) {
	t.Helper()
	e := NewEmbedsPlugin()
	e.config.CacheDir = cacheDir
	e.cacheMetadata(rawURL, "og", meta)
}

func runExternalLinkHover(t *testing.T, m *lifecycle.Manager, posts ...*models.Post) {
	t.Helper()
	m.SetPosts(posts)
	p := NewExternalLinkHoverPlugin()
	if err := p.Configure(m); err != nil {
		t.Fatalf("configure: %v", err)
	}
	if err := p.Render(m); err != nil {
		t.Fatalf("render: %v", err)
	}
}

func TestExternalLinkHover_DisabledByDefault(t *testing.T) {
	m, _ := newExternalLinkHoverManager(t, map[string]any{}, nil)
	post := &models.Post{Path: "a.md", ArticleHTML: `<p><a href="https://go.dev/">Go</a></p>`}
	original := post.ArticleHTML
	runExternalLinkHover(t, m, post)
	if post.ArticleHTML != original {
		t.Errorf("disabled plugin changed HTML: %s", post.ArticleHTML)
	}
}

func TestExternalLinkHover_CachedMetadata(t *testing.T) {
	m, cacheDir := newExternalLinkHoverManager(t, map[string]any{
		"enabled":         true,
		"favicon_service": "https://icons.test/{host}.ico",
	}, nil)
	seedEmbedsCache(t, cacheDir, "https://go.dev/doc/", &OGMetadata{
		Title:       "Documentation",
		Description: `The Go "docs" & more`,
		Image:       "/images/logo.svg",
		SiteName:    "Go",
	})

	post := &models.Post{Path: "a.md", ArticleHTML: `<p>Read <a href="https://go.dev/doc/">the docs</a>.</p>`}
	runExternalLinkHover(t, m, post)

	for _, want := range []string{
		`data-link-preview`,
		`data-link-host="go.dev"`,
		`data-link-title="Documentation"`,
		`data-link-description="The Go &#34;docs&#34; &amp; more"`,
		`data-link-site="Go"`,
		`data-link-image="https://go.dev/images/logo.svg"`,
		`data-link-icon="https://icons.test/go.dev.ico"`,
		`>the docs</a>`,
	} {
		if !strings.Contains(post.ArticleHTML, want) {
			t.Errorf("missing %s in %s", want, post.ArticleHTML)
		}
	}
}

func TestExternalLinkHover_Fallbacks(t *testing.T) {
	blogroll := models.NewBlogrollConfig()
	blogroll.Feeds = []models.ExternalFeedConfig{{
		URL:         "https://simonwillison.net/atom/everything/",
		SiteURL:     "https://simonwillison.net/",
		Title:       "Simon Willison",
		Description: "Blog about Python, SQLite, and AI tools",
		ImageURL:    "https://simonwillison.net/avatar.jpg",
	}}
	m, _ := newExternalLinkHoverManager(t, map[string]any{"enabled": true, "include_image": false},
		map[string]any{"blogroll": blogroll})

	post := &models.Post{Path: "a.md", ArticleHTML: `<p>` +
		`<a href="https://www.simonwillison.net/2024/">Simon</a> ` +
		`<a href="https://example.com/x" title="An example page">plain</a></p>`}
	runExternalLinkHover(t, m, post)

	for _, want := range []string{
		`data-link-host="simonwillison.net" data-link-description="Blog about Python, SQLite, and AI tools" data-link-site="Simon Willison">Simon</a>`,
		`data-link-host="example.com" data-link-description="An example page">plain</a>`,
	} {
		if !strings.Contains(post.ArticleHTML, want) {
			t.Errorf("missing %s in %s", want, post.ArticleHTML)
		}
	}
	if strings.Contains(post.ArticleHTML, "data-link-image") {
		t.Errorf("include_image=false still emitted image: %s", post.ArticleHTML)
	}
}

func TestExternalLinkHover_SkipsIneligibleLinks(t *testing.T) {
	m, _ := newExternalLinkHoverManager(t, map[string]any{
		"enabled":        true,
		"ignore_domains": []interface{}{"ignored.test"},
	}, nil)
	links := []string{
		`<a href="/local/">relative</a>`,
		`<a href="https://mysite.test/post/">own site</a>`,
		`<a href="mailto:me@example.com">mail</a>`,
		`<a href="https://cdn.ignored.test/x">ignored subdomain</a>`,
		`<a href="https://example.com/" class="wikilink">wikilink</a>`,
		`<a href="https://example.com/" class="mention">mention</a>`,
		`<a href="https://example.com/" class="embed-card-link">embed</a>`,
		`<a href="https://example.com/" data-no-preview>opted out</a>`,
		`<a href="https://example.com/" data-hover-card="#x">custom</a>`,
	}
	post := &models.Post{Path: "a.md", ArticleHTML: "<p>" + strings.Join(links, " ") + "</p>"}
	runExternalLinkHover(t, m, post)
	if strings.Contains(post.ArticleHTML, "data-link-preview") {
		t.Errorf("ineligible link annotated: %s", post.ArticleHTML)
	}
}

func TestExternalLinkHover_SkipsPrivatePosts(t *testing.T) {
	m, _ := newExternalLinkHoverManager(t, map[string]any{"enabled": true}, nil)
	post := &models.Post{Path: "a.md", Private: true, ArticleHTML: `<a href="https://go.dev/">Go</a>`}
	runExternalLinkHover(t, m, post)
	if strings.Contains(post.ArticleHTML, "data-link-preview") {
		t.Errorf("private post annotated: %s", post.ArticleHTML)
	}
}

func TestExternalLinkHover_FetchOptIn(t *testing.T) {
	var hits int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		atomic.AddInt32(&hits, 1)
		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write([]byte(`<html><head><meta property="og:title" content="Fetched Title"></head></html>`)) //nolint:errcheck // test server
	}))
	defer server.Close()

	link := `<a href="` + server.URL + `/page">x</a>`

	m, _ := newExternalLinkHoverManager(t, map[string]any{"enabled": true}, nil)
	post := &models.Post{Path: "a.md", ArticleHTML: link}
	runExternalLinkHover(t, m, post)
	if atomic.LoadInt32(&hits) != 0 {
		t.Fatalf("fetch=false made %d requests", hits)
	}
	if strings.Contains(post.ArticleHTML, "data-link-title") {
		t.Errorf("unexpected title without fetch: %s", post.ArticleHTML)
	}

	m, _ = newExternalLinkHoverManager(t, map[string]any{"enabled": true, "fetch": true}, nil)
	posts := []*models.Post{
		{Path: "a.md", ArticleHTML: link},
		{Path: "b.md", ArticleHTML: link + link},
	}
	runExternalLinkHover(t, m, posts...)
	if got := atomic.LoadInt32(&hits); got != 1 {
		t.Errorf("expected 1 request for a repeated URL, got %d", got)
	}
	for _, p := range posts {
		if !strings.Contains(p.ArticleHTML, `data-link-title="Fetched Title"`) {
			t.Errorf("missing fetched title: %s", p.ArticleHTML)
		}
	}
}

func TestExternalLinkHover_IgnoresFallbackTitle(t *testing.T) {
	m, cacheDir := newExternalLinkHoverManager(t, map[string]any{"enabled": true}, nil)
	seedEmbedsCache(t, cacheDir, "https://down.test/", &OGMetadata{Title: models.NewEmbedsConfig().FallbackTitle})
	post := &models.Post{Path: "a.md", ArticleHTML: `<a href="https://down.test/">down</a>`}
	runExternalLinkHover(t, m, post)
	if strings.Contains(post.ArticleHTML, "data-link-title") {
		t.Errorf("fallback title leaked into card: %s", post.ArticleHTML)
	}
	if !strings.Contains(post.ArticleHTML, `data-link-host="down.test"`) {
		t.Errorf("expected host-only annotation: %s", post.ArticleHTML)
	}
}
