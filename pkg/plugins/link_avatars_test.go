package plugins

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

func TestLinkAvatars_Render_LocalModeInjectsAndCaches(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		if _, err := w.Write([]byte("fake-png")); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}))
	defer server.Close()

	outputDir := t.TempDir()
	config := &lifecycle.Config{
		OutputDir: outputDir,
		Extra: map[string]interface{}{
			"models_config": &models.Config{URL: "https://mysite.test"},
			"link_avatars": map[string]any{
				"enabled":  true,
				"mode":     "local",
				"service":  "custom",
				"template": server.URL + "/favicon/{host}.png",
				"size":     16,
			},
		},
	}

	m := lifecycle.NewManager()
	m.SetConfig(config)

	post := &models.Post{Path: "post.md", ArticleHTML: `<p><a href="https://example.com/path">Example</a></p>`}
	m.SetPosts([]*models.Post{post})

	plugin := NewLinkAvatarsPlugin()
	if err := plugin.Configure(m); err != nil {
		t.Fatalf("configure error: %v", err)
	}

	if err := plugin.Render(m); err != nil {
		t.Fatalf("render error: %v", err)
	}

	if !strings.Contains(post.ArticleHTML, "has-avatar") {
		t.Errorf("expected injected avatar class, got %s", post.ArticleHTML)
	}
	if !strings.Contains(post.ArticleHTML, "data-favicon=") {
		t.Errorf("expected data-favicon attribute, got %s", post.ArticleHTML)
	}

	iconPath := filepath.Join(outputDir, "assets", "markata", "link-avatars", "example.com.png")
	if _, err := os.Stat(iconPath); err != nil {
		t.Fatalf("expected cached icon at %s: %v", iconPath, err)
	}
}

func TestLinkAvatars_Render_LocalModeCachesHTTP404Misses(t *testing.T) {
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		requests++
		http.NotFound(w, nil)
	}))
	defer server.Close()

	outputDir := t.TempDir()
	config := &lifecycle.Config{
		OutputDir: outputDir,
		Extra: map[string]interface{}{
			"models_config": &models.Config{URL: "https://mysite.test"},
			"link_avatars": map[string]any{
				"enabled":  true,
				"mode":     "local",
				"service":  "custom",
				"template": server.URL + "/favicon/{host}.png",
				"size":     16,
			},
		},
	}

	plugin := NewLinkAvatarsPlugin()

	newManager := func() *lifecycle.Manager {
		m := lifecycle.NewManager()
		m.SetConfig(config)
		m.SetPosts([]*models.Post{{Path: "post.md", ArticleHTML: `<p><a href="https://example.com/path">Example</a></p>`}})
		return m
	}

	m1 := newManager()
	if err := plugin.Configure(m1); err != nil {
		t.Fatalf("configure error: %v", err)
	}
	if err := plugin.Render(m1); err != nil {
		t.Fatalf("first render error: %v", err)
	}

	m2 := newManager()
	if err := plugin.Configure(m2); err != nil {
		t.Fatalf("configure error: %v", err)
	}
	if err := plugin.Render(m2); err != nil {
		t.Fatalf("second render error: %v", err)
	}

	if requests != 1 {
		t.Fatalf("requests = %d, want 1", requests)
	}

	missPath := filepath.Join(outputDir, "assets", "markata", "link-avatars", "example.com.missing")
	if _, err := os.Stat(missPath); err != nil {
		t.Fatalf("expected cached miss marker at %s: %v", missPath, err)
	}
}

func TestLinkAvatars_ConfigureHostedRequiresBaseURL(t *testing.T) {
	config := &lifecycle.Config{
		Extra: map[string]interface{}{
			"models_config": &models.Config{URL: "https://mysite.test"},
			"link_avatars": map[string]any{
				"enabled": true,
				"mode":    "hosted",
			},
		},
	}

	m := lifecycle.NewManager()
	m.SetConfig(config)

	plugin := NewLinkAvatarsPlugin()
	if err := plugin.Configure(m); err == nil {
		t.Fatal("expected error when hosted_base_url is missing")
	}
}

func TestLinkAvatars_LocalMode_WithURLPathPrefix(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		if _, err := w.Write([]byte("fake-png")); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}))
	defer server.Close()

	outputDir := t.TempDir()
	config := &lifecycle.Config{
		OutputDir: outputDir,
		Extra: map[string]interface{}{
			"models_config": &models.Config{URL: "https://mysite.test/blog"},
			"link_avatars": map[string]any{
				"enabled":  true,
				"mode":     "local",
				"service":  "custom",
				"template": server.URL + "/favicon/{host}.png",
				"size":     16,
			},
		},
	}

	m := lifecycle.NewManager()
	m.SetConfig(config)

	post := &models.Post{Path: "post.md", ArticleHTML: `<p><a href="https://example.com/path">Example</a></p>`}
	m.SetPosts([]*models.Post{post})

	plugin := NewLinkAvatarsPlugin()
	if err := plugin.Configure(m); err != nil {
		t.Fatalf("configure error: %v", err)
	}

	if err := plugin.Render(m); err != nil {
		t.Fatalf("render error: %v", err)
	}

	if !strings.Contains(post.ArticleHTML, "has-avatar") {
		t.Errorf("expected injected avatar class, got %s", post.ArticleHTML)
	}
	if !strings.Contains(post.ArticleHTML, `data-favicon="/assets/markata/link-avatars/`) {
		t.Errorf("expected data-favicon with absolute path, got %s", post.ArticleHTML)
	}

	iconPath := filepath.Join(outputDir, "assets", "markata", "link-avatars", "example.com.png")
	if _, err := os.Stat(iconPath); err != nil {
		t.Fatalf("expected cached icon at %s: %v", iconPath, err)
	}
}

func TestLinkAvatars_Configure_InsertsRootAbsoluteAssets(t *testing.T) {
	config := &lifecycle.Config{
		Extra: map[string]interface{}{
			"models_config": &models.Config{URL: "https://mysite.test/blog"},
			"link_avatars": map[string]any{
				"enabled": true,
				"mode":    "js",
			},
		},
	}

	m := lifecycle.NewManager()
	m.SetConfig(config)

	plugin := NewLinkAvatarsPlugin()
	if err := plugin.Configure(m); err != nil {
		t.Fatalf("configure error: %v", err)
	}

	modelsConfig, ok := config.Extra["models_config"].(*models.Config)
	if !ok || modelsConfig == nil {
		t.Fatalf("expected models_config to be set")
	}
	if len(modelsConfig.Head.Link) == 0 {
		t.Fatalf("expected head link tags to be injected")
	}

	var cssHref string
	for _, link := range modelsConfig.Head.Link {
		if link.Rel == "stylesheet" && strings.Contains(link.Href, "link-avatars") {
			cssHref = link.Href
			break
		}
	}
	if cssHref == "" {
		t.Fatalf("expected link-avatars stylesheet link")
	}
	if !strings.HasPrefix(cssHref, "/css/") {
		t.Fatalf("expected root-absolute css href, got %q", cssHref)
	}

	var jsSrc string
	for _, script := range modelsConfig.Head.Script {
		if strings.Contains(script.Src, "link-avatars") {
			jsSrc = script.Src
			break
		}
	}
	if jsSrc == "" {
		t.Fatalf("expected link-avatars script tag")
	}
	if !strings.HasPrefix(jsSrc, "/js/") {
		t.Fatalf("expected root-absolute js src, got %q", jsSrc)
	}
}

func TestLinkAvatars_Write_JSMode_WritesAssetsToCssAndJs(t *testing.T) {
	outputDir := t.TempDir()
	config := &lifecycle.Config{
		OutputDir: outputDir,
		Extra: map[string]interface{}{
			"models_config": &models.Config{URL: "https://mysite.test"},
			"link_avatars": map[string]any{
				"enabled": true,
				"mode":    "js",
			},
		},
	}

	m := lifecycle.NewManager()
	m.SetConfig(config)

	plugin := NewLinkAvatarsPlugin()
	if err := plugin.Configure(m); err != nil {
		t.Fatalf("configure error: %v", err)
	}
	if err := plugin.Write(m); err != nil {
		t.Fatalf("write error: %v", err)
	}

	cssPath := filepath.Join(outputDir, "css", "link-avatars.css")
	if _, err := os.Stat(cssPath); err != nil {
		t.Fatalf("expected css asset at %s: %v", cssPath, err)
	}
	jsPath := filepath.Join(outputDir, "js", "link-avatars.js")
	if _, err := os.Stat(jsPath); err != nil {
		t.Fatalf("expected js asset at %s: %v", jsPath, err)
	}

	if hash := m.GetAssetHash("css/link-avatars.css"); hash == "" {
		t.Fatalf("expected css hash to be registered")
	}
	if hash := m.GetAssetHash("js/link-avatars.js"); hash == "" {
		t.Fatalf("expected js hash to be registered")
	}
}
