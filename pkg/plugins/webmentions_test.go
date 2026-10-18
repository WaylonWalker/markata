package plugins

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

func TestWebMentionsPlugin_Name(t *testing.T) {
	plugin := NewWebMentionsPlugin()
	if got := plugin.Name(); got != "webmentions" {
		t.Errorf("Name() = %q, want %q", got, "webmentions")
	}
}

func TestWebMentionsPlugin_InterfaceCompliance(_ *testing.T) {
	var _ lifecycle.Plugin = (*WebMentionsPlugin)(nil)
	var _ lifecycle.ConfigurePlugin = (*WebMentionsPlugin)(nil)
	var _ lifecycle.CollectPlugin = (*WebMentionsPlugin)(nil)
	var _ lifecycle.PriorityPlugin = (*WebMentionsPlugin)(nil)
}

func TestWebMentionsPlugin_Configure(t *testing.T) {
	plugin := NewWebMentionsPlugin()
	m := lifecycle.NewManager()

	// Set up config
	m.Config().Extra = map[string]interface{}{
		"url": "https://example.com",
		"webmentions": models.WebMentionsConfig{
			Enabled:   true,
			Outgoing:  true,
			UserAgent: "test-agent/1.0",
			Timeout:   "15s",
		},
	}

	err := plugin.Configure(m)
	if err != nil {
		t.Fatalf("Configure() error = %v", err)
	}

	if plugin.siteURL != "https://example.com" {
		t.Errorf("siteURL = %q, want %q", plugin.siteURL, "https://example.com")
	}

	if plugin.config.UserAgent != "test-agent/1.0" {
		t.Errorf("config.UserAgent = %q, want %q", plugin.config.UserAgent, "test-agent/1.0")
	}
}

func TestWebMentionsPlugin_Priority(t *testing.T) {
	plugin := NewWebMentionsPlugin()

	// Should run late in Collect stage
	collectPriority := plugin.Priority(lifecycle.StageCollect)
	if collectPriority <= lifecycle.PriorityLate {
		t.Errorf("Priority(StageCollect) = %d, want > %d", collectPriority, lifecycle.PriorityLate)
	}

	// Default priority for other stages
	renderPriority := plugin.Priority(lifecycle.StageRender)
	if renderPriority != lifecycle.PriorityDefault {
		t.Errorf("Priority(StageRender) = %d, want %d", renderPriority, lifecycle.PriorityDefault)
	}
}

func TestWebMentionsPlugin_Collect_Disabled(t *testing.T) {
	plugin := NewWebMentionsPlugin()
	m := lifecycle.NewManager()

	// Plugin disabled by default
	err := plugin.Collect(m)
	if err != nil {
		t.Fatalf("Collect() error = %v", err)
	}

	// Should have no mentions
	mentions := plugin.Mentions()
	if len(mentions) != 0 {
		t.Errorf("Mentions() len = %d, want 0", len(mentions))
	}
}

func TestWebMentionsPlugin_Collect_NoSiteURL(t *testing.T) {
	plugin := NewWebMentionsPlugin()
	plugin.config.Enabled = true
	plugin.config.Outgoing = true

	m := lifecycle.NewManager()

	// No site URL configured
	err := plugin.Collect(m)
	if err != nil {
		t.Fatalf("Collect() error = %v", err)
	}

	// Should have no mentions
	mentions := plugin.Mentions()
	if len(mentions) != 0 {
		t.Errorf("Mentions() len = %d, want 0", len(mentions))
	}
}

func TestWebMentionsPlugin_DiscoverEndpoint_HTTPHeader(t *testing.T) {
	// Create a test server that returns Link header
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Link", `<https://example.com/webmention>; rel="webmention"`)
		//nolint:errcheck // test handler
		w.Write([]byte("<html><body>Test</body></html>"))
	}))
	defer ts.Close()

	plugin := NewWebMentionsPlugin()
	plugin.config.UserAgent = "test-agent"
	plugin.httpClient = &http.Client{Timeout: 10 * time.Second}

	endpoint, err := plugin.discoverEndpoint(ts.URL)
	if err != nil {
		t.Fatalf("discoverEndpoint() error = %v", err)
	}

	if endpoint != "https://example.com/webmention" {
		t.Errorf("endpoint = %q, want %q", endpoint, "https://example.com/webmention")
	}
}

func TestWebMentionsPlugin_DiscoverEndpoint_HTMLLink(t *testing.T) {
	// Create a test server that returns HTML with link tag
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		//nolint:errcheck // test handler
		w.Write([]byte(`<html>
<head>
  <link rel="webmention" href="https://example.com/webmention">
</head>
<body>Test</body>
</html>`))
	}))
	defer ts.Close()

	plugin := NewWebMentionsPlugin()
	plugin.config.UserAgent = "test-agent"
	plugin.httpClient = &http.Client{Timeout: 10 * time.Second}

	endpoint, err := plugin.discoverEndpoint(ts.URL)
	if err != nil {
		t.Fatalf("discoverEndpoint() error = %v", err)
	}

	if endpoint != "https://example.com/webmention" {
		t.Errorf("endpoint = %q, want %q", endpoint, "https://example.com/webmention")
	}
}

func TestWebMentionsPlugin_DiscoverEndpoint_NoEndpoint(t *testing.T) {
	// Create a test server with no webmention endpoint
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		//nolint:errcheck // test handler
		w.Write([]byte("<html><body>Test</body></html>"))
	}))
	defer ts.Close()

	plugin := NewWebMentionsPlugin()
	plugin.config.UserAgent = "test-agent"
	plugin.httpClient = &http.Client{Timeout: 10 * time.Second}

	endpoint, err := plugin.discoverEndpoint(ts.URL)
	if err != nil {
		t.Fatalf("discoverEndpoint() error = %v", err)
	}

	if endpoint != "" {
		t.Errorf("endpoint = %q, want empty", endpoint)
	}
}

func TestWebMentionsPlugin_ExtractEndpointFromHeader(t *testing.T) {
	plugin := NewWebMentionsPlugin()

	tests := []struct {
		name     string
		headers  map[string][]string
		baseURL  string
		expected string
	}{
		{
			name: "standard link header",
			headers: map[string][]string{
				"Link": {`<https://example.com/webmention>; rel="webmention"`},
			},
			baseURL:  "https://example.com/post",
			expected: "https://example.com/webmention",
		},
		{
			name: "link header with single quotes",
			headers: map[string][]string{
				"Link": {`<https://example.com/webmention>; rel='webmention'`},
			},
			baseURL:  "https://example.com/post",
			expected: "https://example.com/webmention",
		},
		{
			name: "link header no quotes",
			headers: map[string][]string{
				"Link": {`<https://example.com/webmention>; rel=webmention`},
			},
			baseURL:  "https://example.com/post",
			expected: "https://example.com/webmention",
		},
		{
			name: "relative URL in header",
			headers: map[string][]string{
				"Link": {`</webmention>; rel="webmention"`},
			},
			baseURL:  "https://example.com/post/",
			expected: "https://example.com/webmention",
		},
		{
			name: "no webmention header",
			headers: map[string][]string{
				"Link": {`<https://example.com/feed>; rel="alternate"`},
			},
			baseURL:  "https://example.com/post",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			headers := http.Header{}
			for k, v := range tt.headers {
				for _, val := range v {
					headers.Add(k, val)
				}
			}

			got := plugin.extractEndpointFromHeader(headers, tt.baseURL)
			if got != tt.expected {
				t.Errorf("extractEndpointFromHeader() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestWebMentionsPlugin_ExtractEndpointFromHTML(t *testing.T) {
	plugin := NewWebMentionsPlugin()

	tests := []struct {
		name     string
		html     string
		baseURL  string
		expected string
	}{
		{
			name:     "link tag rel first",
			html:     `<html><head><link rel="webmention" href="https://example.com/webmention"></head></html>`,
			baseURL:  "https://example.com/post",
			expected: "https://example.com/webmention",
		},
		{
			name:     "link tag href first",
			html:     `<html><head><link href="https://example.com/webmention" rel="webmention"></head></html>`,
			baseURL:  "https://example.com/post",
			expected: "https://example.com/webmention",
		},
		{
			name:     "relative href",
			html:     `<html><head><link rel="webmention" href="/webmention"></head></html>`,
			baseURL:  "https://example.com/post/",
			expected: "https://example.com/webmention",
		},
		{
			name:     "a tag webmention",
			html:     `<html><body><a rel="webmention" href="https://example.com/webmention">Webmention</a></body></html>`,
			baseURL:  "https://example.com/post",
			expected: "https://example.com/webmention",
		},
		{
			name:     "no webmention link",
			html:     `<html><head><link rel="stylesheet" href="/style.css"></head></html>`,
			baseURL:  "https://example.com/post",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := plugin.extractEndpointFromHTML(tt.html, tt.baseURL)
			if got != tt.expected {
				t.Errorf("extractEndpointFromHTML() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestWebMentionsPlugin_SendWebMention(t *testing.T) {
	// Create a test server that accepts webmentions
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}

		if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
			t.Errorf("expected Content-Type application/x-www-form-urlencoded, got %s", r.Header.Get("Content-Type"))
		}

		source := r.FormValue("source")
		target := r.FormValue("target")

		if source == "" || target == "" {
			w.WriteHeader(http.StatusBadRequest)
			//nolint:errcheck // test handler
			w.Write([]byte("missing source or target"))
			return
		}

		w.WriteHeader(http.StatusAccepted)
		//nolint:errcheck // test handler
		w.Write([]byte("Webmention received"))
	}))
	defer ts.Close()

	plugin := NewWebMentionsPlugin()
	plugin.config.UserAgent = "test-agent"
	plugin.httpClient = &http.Client{Timeout: 10 * time.Second}

	mention := &WebMention{
		Source:   "https://my-site.com/post/",
		Target:   "https://target-site.com/article/",
		Endpoint: ts.URL,
	}

	err := plugin.sendWebMention(mention)
	if err != nil {
		t.Fatalf("sendWebMention() error = %v", err)
	}

	if mention.StatusCode != http.StatusAccepted {
		t.Errorf("StatusCode = %d, want %d", mention.StatusCode, http.StatusAccepted)
	}
}

func TestWebMentionsPlugin_CacheKey(t *testing.T) {
	plugin := NewWebMentionsPlugin()

	key1 := plugin.cacheKey("https://source.com/a", "https://target.com/b")
	key2 := plugin.cacheKey("https://source.com/a", "https://target.com/b")
	key3 := plugin.cacheKey("https://source.com/x", "https://target.com/y")

	// Same inputs should produce same key
	if key1 != key2 {
		t.Errorf("cacheKey should be deterministic: %q != %q", key1, key2)
	}

	// Different inputs should produce different keys
	if key1 == key3 {
		t.Errorf("cacheKey should be different for different inputs: %q == %q", key1, key3)
	}

	// Key should be a valid hex string
	if len(key1) != 32 {
		t.Errorf("cacheKey length = %d, want 32", len(key1))
	}
}

func TestWebMentionsPlugin_ResolveURL(t *testing.T) {
	plugin := NewWebMentionsPlugin()

	tests := []struct {
		name     string
		baseURL  string
		href     string
		expected string
	}{
		{
			name:     "absolute URL unchanged",
			baseURL:  "https://example.com/post/",
			href:     "https://other.com/webmention",
			expected: "https://other.com/webmention",
		},
		{
			name:     "relative path",
			baseURL:  "https://example.com/post/",
			href:     "/webmention",
			expected: "https://example.com/webmention",
		},
		{
			name:     "relative to current directory",
			baseURL:  "https://example.com/post/",
			href:     "webmention",
			expected: "https://example.com/post/webmention",
		},
		{
			name:     "empty href",
			baseURL:  "https://example.com/post/",
			href:     "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := plugin.resolveURL(tt.baseURL, tt.href)
			if got != tt.expected {
				t.Errorf("resolveURL() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestWebMentionsPlugin_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	// Create a mock target server with webmention endpoint
	targetServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Link", fmt.Sprintf(`<%s/webmention>; rel="webmention"`, r.Host))
		//nolint:errcheck // test handler
		w.Write([]byte("<html><body>Test article</body></html>"))
	}))
	defer targetServer.Close()

	// Create a mock webmention receiver
	receivedMentions := make(chan *WebMention, 10)
	mentionServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mention := &WebMention{
			Source: r.FormValue("source"),
			Target: r.FormValue("target"),
		}
		receivedMentions <- mention
		w.WriteHeader(http.StatusAccepted)
	}))
	defer mentionServer.Close()

	// Override the target server response to point to our mention server
	targetWithEndpoint := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Link", fmt.Sprintf(`<%s>; rel="webmention"`, mentionServer.URL))
		//nolint:errcheck // test handler
		w.Write([]byte("<html><body>Test article</body></html>"))
	}))
	defer targetWithEndpoint.Close()

	// Create plugin and manager
	plugin := NewWebMentionsPlugin()
	plugin.config.Enabled = true
	plugin.config.Outgoing = true
	plugin.config.UserAgent = "test-agent"
	plugin.siteURL = "https://my-site.com"
	plugin.httpClient = &http.Client{}

	m := lifecycle.NewManager()

	// Create a post with an external link
	title := "Test Post"
	post := &models.Post{
		Path:        "test.md",
		Slug:        "test",
		Href:        "/test/",
		Title:       &title,
		ArticleHTML: fmt.Sprintf(`<p>Check out <a href=%q>this article</a>.</p>`, targetWithEndpoint.URL),
		Outlinks: []*models.Link{
			{
				TargetURL:  targetWithEndpoint.URL,
				IsInternal: false,
			},
		},
	}

	m.SetPosts([]*models.Post{post})

	// Run the plugin
	err := plugin.Collect(m)
	if err != nil {
		t.Fatalf("Collect() error = %v", err)
	}

	// Check mentions were processed
	mentions := plugin.Mentions()
	if len(mentions) != 1 {
		t.Errorf("Mentions() len = %d, want 1", len(mentions))
	}
}

func TestWebMentionsPlugin_SaveLoadCache(t *testing.T) {
	// Create temp directory for cache
	tempDir, err := os.MkdirTemp("", "webmentions-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create plugin with cache
	plugin := NewWebMentionsPlugin()
	plugin.config.CacheDir = tempDir

	// Add some cache entries
	plugin.sentCache["abc123"] = true
	plugin.sentCache["def456"] = true

	// Save cache
	plugin.saveSentCache()

	// Check file was created
	cacheFile := filepath.Join(tempDir, "webmentions_sent.json")
	if _, err := os.Stat(cacheFile); os.IsNotExist(err) {
		t.Fatal("cache file was not created")
	}

	// Create new plugin and load cache
	plugin2 := NewWebMentionsPlugin()
	plugin2.config.CacheDir = tempDir
	plugin2.loadSentCache()

	// Verify cache was loaded
	if !plugin2.sentCache["abc123"] {
		t.Error("cache entry abc123 was not loaded")
	}
	if !plugin2.sentCache["def456"] {
		t.Error("cache entry def456 was not loaded")
	}
}

func TestNewWebMentionsConfig(t *testing.T) {
	config := models.NewWebMentionsConfig()

	if config.Enabled {
		t.Error("Enabled should default to false")
	}
	if !config.Outgoing {
		t.Error("Outgoing should default to true")
	}
	if config.UserAgent == "" {
		t.Error("UserAgent should have a default value")
	}
	if config.Timeout != "30s" {
		t.Errorf("Timeout = %q, want %q", config.Timeout, "30s")
	}
	if config.CacheDir != ".cache/webmentions" {
		t.Errorf("CacheDir = %q, want %q", config.CacheDir, ".cache/webmentions")
	}
	if config.ConcurrentRequests != 5 {
		t.Errorf("ConcurrentRequests = %d, want 5", config.ConcurrentRequests)
	}
}
