package plugins

import (
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"testing"

	"github.com/WaylonWalker/markata/pkg/buildcache"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/themes"
)

func TestConfigFilesHash_ChangesWhenOverlayChanges(t *testing.T) {
	dir := t.TempDir()
	base := filepath.Join(dir, "markata-go.toml")
	overlay := filepath.Join(dir, "tailwind.toml")

	if err := os.WriteFile(base, []byte("title = 'base'\n"), 0o600); err != nil {
		t.Fatalf("WriteFile(base) error = %v", err)
	}
	if err := os.WriteFile(overlay, []byte("[markata-go.tailwind]\nbuild = true\n"), 0o600); err != nil {
		t.Fatalf("WriteFile(overlay) error = %v", err)
	}

	first := buildcache.ContentHash(configFilesHash([]string{base, overlay}))

	if err := os.WriteFile(overlay, []byte("[markata-go.tailwind]\nbuild = false\n"), 0o600); err != nil {
		t.Fatalf("WriteFile(overlay update) error = %v", err)
	}

	second := buildcache.ContentHash(configFilesHash([]string{base, overlay}))
	if first == second {
		t.Fatal("expected config hash to change when overlay config changes")
	}
}

func TestConfigFilesHash_ChangesWhenIncludedFileChanges(t *testing.T) {
	dir := t.TempDir()
	base := filepath.Join(dir, "markata-go.toml")
	included := filepath.Join(dir, "config", "theme.toml")
	if err := os.MkdirAll(filepath.Dir(included), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(base, []byte("[markata-go]\ninclude = [\"config/*.toml\"]\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(included, []byte("[markata-go.theme]\npalette = \"a\"\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	first := buildcache.ContentHash(configFilesHash([]string{base}))
	if err := os.WriteFile(included, []byte("[markata-go.theme]\npalette = \"b\"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	second := buildcache.ContentHash(configFilesHash([]string{base}))
	if first == second {
		t.Fatal("expected config hash to change when an included config file changes")
	}
}

func TestConfigFilesHash_IsStableAcrossPathFormsAndOrder(t *testing.T) {
	dir := t.TempDir()
	base := filepath.Join(dir, "markata-go.toml")
	overlay := filepath.Join(dir, "tailwind.toml")

	if err := os.WriteFile(base, []byte("title = 'base'\n"), 0o600); err != nil {
		t.Fatalf("WriteFile(base) error = %v", err)
	}
	if err := os.WriteFile(overlay, []byte("[markata-go.tailwind]\nbuild = true\n"), 0o600); err != nil {
		t.Fatalf("WriteFile(overlay) error = %v", err)
	}

	first := buildcache.ContentHash(configFilesHash([]string{overlay, base, ""}))
	second := buildcache.ContentHash(configFilesHash([]string{filepath.Clean(base), filepath.Clean(overlay)}))
	third := buildcache.ContentHash(configFilesHash([]string{base, overlay}))

	if first != second || second != third {
		t.Fatalf("expected config hash to be stable across path forms and ordering: %q %q %q", first, second, third)
	}
}

func TestConfigHashInput_IsStableForEquivalentConfig(t *testing.T) {
	dir := t.TempDir()
	base := filepath.Join(dir, "markata-go.toml")
	if err := os.WriteFile(base, []byte("title = 'base'\n"), 0o600); err != nil {
		t.Fatalf("WriteFile(base) error = %v", err)
	}

	config := &lifecycle.Config{
		ContentDir:   ".",
		OutputDir:    "output",
		GlobPatterns: []string{"**/*.md"},
		Extra: map[string]interface{}{
			"title":       "Test Site",
			"url":         "https://example.com",
			"config_path": base,
			"config_paths": []string{
				base,
			},
		},
	}

	first := buildcache.ContentHash(configHashInput(config, []string{base}))
	second := buildcache.ContentHash(configHashInput(config, []string{filepath.Clean(base)}))
	if first != second {
		t.Fatalf("expected effective config hash to be stable for equivalent config: %q != %q", first, second)
	}
}

func TestConfigHashInput_IgnoresOutputDir(t *testing.T) {
	dir := t.TempDir()
	base := filepath.Join(dir, "markata-go.toml")
	if err := os.WriteFile(base, []byte("title = 'base'\n"), 0o600); err != nil {
		t.Fatalf("WriteFile(base) error = %v", err)
	}

	configA := &lifecycle.Config{
		ContentDir:   ".",
		OutputDir:    "output",
		GlobPatterns: []string{"**/*.md"},
		Extra: map[string]interface{}{
			"config_path":  base,
			"config_paths": []string{base},
		},
	}
	configB := &lifecycle.Config{
		ContentDir:   ".",
		OutputDir:    "/data/site/releases/20260614T000000Z-site.tmp",
		GlobPatterns: []string{"**/*.md"},
		Extra: map[string]interface{}{
			"config_path":  base,
			"config_paths": []string{base},
		},
	}

	first := buildcache.ContentHash(configHashInput(configA, []string{base}))
	second := buildcache.ContentHash(configHashInput(configB, []string{base}))
	if first != second {
		t.Fatalf("expected config hash to ignore output dir changes: %q != %q", first, second)
	}
}

func TestBuildCacheConfigure_DefaultsCacheDirToContentDir(t *testing.T) {
	contentDir := t.TempDir()
	outputDir := filepath.Join(t.TempDir(), "build", "site")

	config := &lifecycle.Config{
		ContentDir: contentDir,
		OutputDir:  outputDir,
		Extra: map[string]interface{}{
			"config_path": filepath.Join(contentDir, "markata-go.toml"),
		},
	}

	plugin := NewBuildCachePlugin()
	manager := lifecycle.NewManager()
	cfg := manager.Config()
	*cfg = *config

	if err := plugin.Configure(manager); err != nil {
		t.Fatalf("Configure() error = %v", err)
	}

	cache := GetBuildCache(manager)
	if cache == nil {
		t.Fatal("expected build cache to be stored on manager")
	}

	embeddedTemplatesHash, err := buildcache.HashFS(themes.DefaultTemplates(), []string{".html", ".txt", ".md"})
	if err != nil {
		t.Fatalf("hash embedded templates: %v", err)
	}
	wantTemplatesHash := buildcache.ContentHash("project:\nembedded:" + embeddedTemplatesHash)
	if cache.TemplatesHash != wantTemplatesHash {
		t.Fatalf("TemplatesHash = %q, want combined embedded template hash %q", cache.TemplatesHash, wantTemplatesHash)
	}

	want := filepath.Join(contentDir, ".markata", buildcache.CacheFileName)
	if cachePath := cachePathForTest(cache); cachePath != want {
		t.Fatalf("cache path = %q, want %q", cachePath, want)
	}
}

func cachePathForTest(c *buildcache.Cache) string {
	return filepath.Clean(reflect.ValueOf(c).Elem().FieldByName("path").String())
}

func TestBuildCacheLoad_InvalidatesOnlyChangedPostsAndTransitiveDependents(t *testing.T) {
	cache := buildcache.New(t.TempDir())
	cache.MarkRebuilt("content/target.md", "old-target", "output/target/index.html", "post.html")
	cache.MarkRebuilt("content/source.md", "source-hash", "output/source/index.html", "post.html")
	cache.MarkRebuilt("content/downstream.md", "downstream-hash", "output/downstream/index.html", "post.html")
	cache.MarkRebuilt("content/unrelated.md", "unrelated-hash", "output/unrelated/index.html", "post.html")
	cache.SetDependencies("content/source.md", "source", []string{"target"})
	cache.SetDependencies("content/downstream.md", "downstream", []string{"source"})
	cache.MarkRebuilt("content/root.md", "root-hash", "output/index.html", "home.html")
	cache.SetPostSlug("content/target.md", "target")
	cache.SetPostSlug("content/unrelated.md", "unrelated")
	cache.ResetStats()

	manager := lifecycle.NewManager()
	manager.Cache().Set("build_cache", cache)
	manager.SetPosts([]*models.Post{
		{Path: "content/target.md", Slug: "target", InputHash: "new-target", Template: "post.html"},
		{Path: "content/source.md", Slug: "source", InputHash: "source-hash", Template: "post.html"},
		{Path: "content/downstream.md", Slug: "downstream", InputHash: "downstream-hash", Template: "post.html"},
		{Path: "content/unrelated.md", Slug: "unrelated", InputHash: "unrelated-hash", Template: "post.html"},
		{Path: "content/root.md", Slug: "", InputHash: "root-hash", Template: "home.html"},
	})

	plugin := NewBuildCachePlugin()
	plugin.cache = cache
	if err := plugin.Load(manager); err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	affected := lifecycle.GetServeAffectedPaths(manager)
	if !affected["content/target.md"] || !affected["content/source.md"] || !affected["content/downstream.md"] {
		t.Fatalf("affected paths = %v, want changed target and transitive dependents", affected)
	}
	if len(affected) != 3 || affected["content/unrelated.md"] {
		t.Fatalf("unrelated post was invalidated: %v", affected)
	}
	if affected["content/root.md"] {
		t.Fatalf("unchanged empty-slug post was invalidated: %v", affected)
	}

	changedSlugs := cache.GetChangedSlugs()
	if want := []string{"downstream", "source", "target"}; !slices.Equal(changedSlugs, want) {
		t.Fatalf("changed slugs = %v, want %v", changedSlugs, want)
	}
}

func TestBuildCacheLoad_ChangedHomepageIsAffected(t *testing.T) {
	cache := buildcache.New(t.TempDir())
	cache.MarkRebuilt("content/other.md", "other-hash", "output/other/index.html", "post.html")
	cache.ResetStats()

	manager := lifecycle.NewManager()
	manager.Cache().Set("build_cache", cache)
	manager.SetPosts([]*models.Post{
		{Path: "content/other.md", Slug: "other", InputHash: "other-hash", Template: "post.html"},
		{Path: "content/index.md", Slug: "", InputHash: "home-hash", Template: "post.html"},
	})

	plugin := NewBuildCachePlugin()
	plugin.cache = cache
	if err := plugin.Load(manager); err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	affected := lifecycle.GetServeAffectedPaths(manager)
	if !affected["content/index.md"] {
		t.Fatalf("new empty-slug homepage must be affected so incremental builds render it: %v", affected)
	}
	if affected["content/other.md"] {
		t.Fatalf("unchanged post was invalidated: %v", affected)
	}
}
