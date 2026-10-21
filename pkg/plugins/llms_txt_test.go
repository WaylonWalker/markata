package plugins

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/WaylonWalker/markata-go/pkg/lifecycle"
	"github.com/WaylonWalker/markata-go/pkg/models"
)

func llmsTXTTestPost(slug, title string, date *time.Time) *models.Post {
	return &models.Post{
		Slug: slug, Href: "/" + slug + "/", Title: &title,
		Published: true, Date: date,
	}
}

func TestLLMsTXTPlugin_GeneratesPublicIndex(t *testing.T) {
	root := t.TempDir()
	cfg := lifecycle.NewConfig()
	cfg.OutputDir = filepath.Join(root, "output")
	cfg.ContentDir = root
	cfg.Extra["url"] = "https://example.com/"
	cfg.Extra["title"] = "Example Blog"
	cfg.Extra["description"] = "A useful website"
	m := lifecycle.NewManager()
	m.SetConfig(cfg)
	recentDate := time.Date(2026, 10, 17, 0, 0, 0, 0, time.UTC)
	oldDate := recentDate.AddDate(-1, 0, 0)
	recent := llmsTXTTestPost("latest-post", "Latest [post]", &recentDate)
	note := "An overview \n with helpful  details."
	recent.Description = &note
	docs := llmsTXTTestPost("docs/guide", "Guide", &oldDate)
	private := llmsTXTTestPost("private", "Private secret", &recentDate)
	private.Private = true
	draft := llmsTXTTestPost("draft", "Draft secret", &recentDate)
	draft.Draft = true
	unpublished := llmsTXTTestPost("unpublished", "Unpublished secret", &recentDate)
	unpublished.Published = false
	skipped := llmsTXTTestPost("skipped", "Skipped secret", &recentDate)
	skipped.Skip = true
	noindex := llmsTXTTestPost("noindex", "Noindex secret", &recentDate)
	noindex.Extra = map[string]interface{}{"robots": "noindex, follow"}
	duplicate := llmsTXTTestPost("latest-post", "Duplicate title", &oldDate)
	m.SetPosts([]*models.Post{private, recent, docs, draft, unpublished, skipped, noindex, duplicate})

	if err := os.MkdirAll(cfg.OutputDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(cfg.OutputDir, "sitemap.xml"), []byte("xml"), 0o600); err != nil {
		t.Fatal(err)
	}
	plugin := NewLLMsTXTPlugin()
	if plugin.Name() != "llms_txt" || plugin.Priority(lifecycle.StageWrite) != lifecycle.PriorityLast {
		t.Fatal("unexpected plugin registration or stage priority")
	}
	if err := plugin.Write(m); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(cfg.OutputDir, "llms.txt")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	got := string(data)
	for _, want := range []string{
		"# Example Blog", "> A useful website",
		"## Site", "[Home](https://example.com/)",
		"## Documentation", "[Guide](https://example.com/docs/guide.md)",
		"## Recent Content", "[Latest \\[post\\]](https://example.com/latest-post.md): An overview with helpful details.",
		"## Optional", "[Sitemap](https://example.com/sitemap.xml)",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q in:\n%s", want, got)
		}
	}
	for _, forbidden := range []string{"Private secret", "Draft secret", "Unpublished secret", "Skipped secret", "Noindex secret", "Duplicate title"} {
		if strings.Contains(got, forbidden) {
			t.Errorf("leaked %q in index", forbidden)
		}
	}

	// A warm build must refresh an existing generated file, but not change
	// bytes when the source metadata has not changed.
	if err := plugin.Write(m); err != nil {
		t.Fatal(err)
	}
	again, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(again) != got {
		t.Fatal("unchanged content generated a different llms.txt on warm build")
	}
}

func TestLLMsTXTPlugin_ConfigLimitsAndHTMLFallback(t *testing.T) {
	cfg := lifecycle.NewConfig()
	cfg.OutputDir = t.TempDir()
	cfg.Extra["llms_txt"] = map[string]interface{}{"max_pages": int64(1)}
	cfg.Extra["post_formats"] = models.PostFormatsConfig{Markdown: false}
	m := lifecycle.NewManager()
	m.SetConfig(cfg)
	m.SetPosts([]*models.Post{
		llmsTXTTestPost("first", "First", nil),
		llmsTXTTestPost("second", "Second", nil),
		llmsTXTTestPost("third", "Third", nil),
	})
	if err := NewLLMsTXTPlugin().Write(m); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(cfg.OutputDir, "llms.txt"))
	if err != nil {
		t.Fatal(err)
	}
	got := string(data)
	if !strings.Contains(got, "[First](/first/)") || strings.Contains(got, "[Second]") || strings.Contains(got, "[Third]") {
		t.Fatalf("unexpected max_pages selection or HTML fallback:\n%s", got)
	}
	if strings.Contains(got, "example.com") {
		t.Fatal("used the placeholder domain instead of relative links")
	}
}

func TestLLMsTXTPlugin_PreservesManualAndStaticSources(t *testing.T) {
	t.Run("manual post", func(t *testing.T) {
		cfg := lifecycle.NewConfig()
		cfg.OutputDir = t.TempDir()
		m := lifecycle.NewManager()
		m.SetConfig(cfg)
		m.SetPosts([]*models.Post{llmsTXTTestPost("llms", "Manual", nil)})
		file := filepath.Join(cfg.OutputDir, "llms.txt")
		if err := os.WriteFile(file, []byte("# My hand-written file\n"), 0o600); err != nil {
			t.Fatal(err)
		}
		if err := NewLLMsTXTPlugin().Write(m); err != nil {
			t.Fatal(err)
		}
		got, err := os.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}
		if string(got) != "# My hand-written file\n" {
			t.Fatal("overwrote an authored llms.md")
		}
	})
	t.Run("project static", func(t *testing.T) {
		root := t.TempDir()
		cfg := lifecycle.NewConfig()
		cfg.ContentDir = root
		cfg.OutputDir = filepath.Join(root, "output")
		cfg.Extra["assets_dir"] = "public-assets"
		if err := os.MkdirAll(filepath.Join(root, "public-assets"), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(root, "public-assets", "llms.txt"), []byte("# Static\n"), 0o600); err != nil {
			t.Fatal(err)
		}
		m := lifecycle.NewManager()
		m.SetConfig(cfg)
		if err := NewLLMsTXTPlugin().Write(m); err != nil {
			t.Fatal(err)
		}
		if _, err := os.Stat(filepath.Join(cfg.OutputDir, "llms.txt")); !os.IsNotExist(err) {
			t.Fatalf("expected auto generator to defer to project static asset; stat error = %v", err)
		}
	})
	t.Run("unmanaged output", func(t *testing.T) {
		cfg := lifecycle.NewConfig()
		cfg.OutputDir = t.TempDir()
		file := filepath.Join(cfg.OutputDir, "llms.txt")
		if err := os.WriteFile(file, []byte("# Hand-managed output\n"), 0o600); err != nil {
			t.Fatal(err)
		}
		m := lifecycle.NewManager()
		m.SetConfig(cfg)
		if err := NewLLMsTXTPlugin().Write(m); err != nil {
			t.Fatal(err)
		}
		data, err := os.ReadFile(file)
		if err != nil || string(data) != "# Hand-managed output\n" {
			t.Fatalf("overwrote hand-managed output: %q, %v", data, err)
		}
	})
}

func TestLLMsTXTPlugin_DisabledCleansOnlyGeneratedFile(t *testing.T) {
	cfg := lifecycle.NewConfig()
	cfg.OutputDir = t.TempDir()
	m := lifecycle.NewManager()
	m.SetConfig(cfg)
	plugin := NewLLMsTXTPlugin()
	if err := plugin.Write(m); err != nil {
		t.Fatal(err)
	}
	output := filepath.Join(cfg.OutputDir, "llms.txt")
	cfg.Extra["llms_txt"] = map[string]interface{}{"enabled": false}
	if err := plugin.Write(m); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(output); !os.IsNotExist(err) {
		t.Fatalf("generated file should be removed when disabled: %v", err)
	}
	if err := os.WriteFile(output, []byte("# Manual\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := plugin.Write(m); err != nil {
		t.Fatal(err)
	}
	if got, err := os.ReadFile(output); err != nil || string(got) != "# Manual\n" {
		t.Fatalf("disabled generator changed manual file: %q, %v", got, err)
	}
}

func TestLLMsTXTPlugin_InDefaultButNotSinglePageBuilds(t *testing.T) {
	foundDefault := false
	for _, plugin := range DefaultPlugins() {
		if plugin.Name() == "llms_txt" {
			foundDefault = true
		}
	}
	if !foundDefault {
		t.Fatal("llms_txt missing from DefaultPlugins")
	}
	for _, plugin := range SinglePagePlugins() {
		if plugin.Name() == "llms_txt" {
			t.Fatal("llms_txt should not run for a single-page build")
		}
	}
}
