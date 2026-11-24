package plugins

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/WaylonWalker/markata/pkg/buildcache"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

func TestPagefindPlugin_Name(t *testing.T) {
	plugin := NewPagefindPlugin()
	if got := plugin.Name(); got != "pagefind" {
		t.Errorf("Name() = %v, want pagefind", got)
	}
}

func TestPagefindPlugin_Priority(t *testing.T) {
	plugin := NewPagefindPlugin()

	tests := []struct {
		stage    lifecycle.Stage
		expected int
	}{
		{lifecycle.StageCleanup, lifecycle.PriorityLast},
		{lifecycle.StageWrite, lifecycle.PriorityDefault},
		{lifecycle.StageRender, lifecycle.PriorityDefault},
	}

	for _, tt := range tests {
		t.Run(string(tt.stage), func(t *testing.T) {
			if got := plugin.Priority(tt.stage); got != tt.expected {
				t.Errorf("Priority(%s) = %v, want %v", tt.stage, got, tt.expected)
			}
		})
	}
}

func TestPagefindPlugin_DisabledByConfig(t *testing.T) {
	plugin := NewPagefindPlugin()

	// Create a manager with search disabled
	m := lifecycle.NewManager()
	config := lifecycle.NewConfig()
	config.OutputDir = t.TempDir()

	// Create output directory with some content
	indexPath := filepath.Join(config.OutputDir, "index.html")
	if err := os.WriteFile(indexPath, []byte("<html><body>Test</body></html>"), 0o600); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Disable search
	disabled := false
	config.Extra["search"] = models.SearchConfig{
		Enabled: &disabled,
	}
	m.SetConfig(config)

	// Should not error, just skip
	if err := plugin.Cleanup(m); err != nil {
		t.Errorf("Cleanup() should not error when disabled, got: %v", err)
	}
}

func TestPagefindPlugin_EnabledByDefault(t *testing.T) {
	plugin := NewPagefindPlugin()

	// Create a manager with default config (search enabled by default)
	m := lifecycle.NewManager()
	config := lifecycle.NewConfig()
	config.OutputDir = t.TempDir()

	// Create output directory with some content
	indexPath := filepath.Join(config.OutputDir, "index.html")
	if err := os.WriteFile(indexPath, []byte("<html><body>Test</body></html>"), 0o600); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	m.SetConfig(config)

	// This will fail if pagefind is not installed, but that's okay
	// The plugin should handle that gracefully (we ignore the error here)
	_ = plugin.Cleanup(m) //nolint:errcheck // Pagefind may not be installed in test environment

	// Check that getSearchConfig returns enabled by default
	searchConfig := getSearchConfig(config)
	if !searchConfig.IsEnabled() {
		t.Error("Search should be enabled by default")
	}
}

func TestPagefindPlugin_SkipsBleveBackend(t *testing.T) {
	plugin := NewPagefindPlugin()

	m := lifecycle.NewManager()
	config := lifecycle.NewConfig()
	config.OutputDir = t.TempDir()

	indexPath := filepath.Join(config.OutputDir, "index.html")
	if err := os.WriteFile(indexPath, []byte("<html><body>Test</body></html>"), 0o600); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	config.Extra["search"] = models.SearchConfig{Backend: "bleve"}
	m.SetConfig(config)

	if err := plugin.Cleanup(m); err != nil {
		t.Fatalf("Cleanup() error = %v", err)
	}

	if _, err := os.Stat(filepath.Join(config.OutputDir, "_pagefind", "pagefind.js")); !os.IsNotExist(err) {
		t.Fatalf("expected pagefind index to be skipped, stat err = %v", err)
	}
}

func TestGetSearchConfig_Default(t *testing.T) {
	config := lifecycle.NewConfig()

	searchConfig := getSearchConfig(config)

	if !searchConfig.IsEnabled() {
		t.Error("IsEnabled() should return true by default")
	}

	if searchConfig.Position != "navbar" {
		t.Errorf("Position = %v, want navbar", searchConfig.Position)
	}

	if searchConfig.Placeholder != "Search…" {
		t.Errorf("Placeholder = %v, want Search…", searchConfig.Placeholder)
	}

	if !searchConfig.IsShowImages() {
		t.Error("IsShowImages() should return true by default")
	}

	if searchConfig.ExcerptLength != 200 {
		t.Errorf("ExcerptLength = %v, want 200", searchConfig.ExcerptLength)
	}

	if searchConfig.Pagefind.BundleDir != "_pagefind" {
		t.Errorf("Pagefind.BundleDir = %v, want _pagefind", searchConfig.Pagefind.BundleDir)
	}
}

func TestGetSearchConfig_FromExtra(t *testing.T) {
	config := lifecycle.NewConfig()

	enabled := true
	showImages := false
	config.Extra["search"] = models.SearchConfig{
		Enabled:       &enabled,
		Position:      "sidebar",
		Placeholder:   "Find...",
		ShowImages:    &showImages,
		ExcerptLength: 300,
		Pagefind: models.PagefindConfig{
			BundleDir:        "_search",
			ExcludeSelectors: []string{".no-search"},
			RootSelector:     "main",
		},
	}

	searchConfig := getSearchConfig(config)

	if !searchConfig.IsEnabled() {
		t.Error("IsEnabled() should return true")
	}

	if searchConfig.Position != "sidebar" {
		t.Errorf("Position = %v, want sidebar", searchConfig.Position)
	}

	if searchConfig.Placeholder != "Find..." {
		t.Errorf("Placeholder = %v, want Find...", searchConfig.Placeholder)
	}

	if searchConfig.IsShowImages() {
		t.Error("IsShowImages() should return false")
	}

	if searchConfig.ExcerptLength != 300 {
		t.Errorf("ExcerptLength = %v, want 300", searchConfig.ExcerptLength)
	}

	if searchConfig.Pagefind.BundleDir != "_search" {
		t.Errorf("Pagefind.BundleDir = %v, want _search", searchConfig.Pagefind.BundleDir)
	}

	if len(searchConfig.Pagefind.ExcludeSelectors) != 1 || searchConfig.Pagefind.ExcludeSelectors[0] != ".no-search" {
		t.Errorf("Pagefind.ExcludeSelectors = %v, want [.no-search]", searchConfig.Pagefind.ExcludeSelectors)
	}

	if searchConfig.Pagefind.RootSelector != "main" {
		t.Errorf("Pagefind.RootSelector = %v, want main", searchConfig.Pagefind.RootSelector)
	}
}

func TestSearchConfig_Methods(t *testing.T) {
	t.Run("IsEnabled_nil", func(t *testing.T) {
		sc := models.SearchConfig{}
		if !sc.IsEnabled() {
			t.Error("IsEnabled() should return true when Enabled is nil")
		}
	})

	t.Run("IsEnabled_true", func(t *testing.T) {
		enabled := true
		sc := models.SearchConfig{Enabled: &enabled}
		if !sc.IsEnabled() {
			t.Error("IsEnabled() should return true when Enabled is true")
		}
	})

	t.Run("IsEnabled_false", func(t *testing.T) {
		enabled := false
		sc := models.SearchConfig{Enabled: &enabled}
		if sc.IsEnabled() {
			t.Error("IsEnabled() should return false when Enabled is false")
		}
	})

	t.Run("IsShowImages_nil", func(t *testing.T) {
		sc := models.SearchConfig{}
		if !sc.IsShowImages() {
			t.Error("IsShowImages() should return true when ShowImages is nil")
		}
	})

	t.Run("IsShowImages_true", func(t *testing.T) {
		showImages := true
		sc := models.SearchConfig{ShowImages: &showImages}
		if !sc.IsShowImages() {
			t.Error("IsShowImages() should return true when ShowImages is true")
		}
	})

	t.Run("IsShowImages_false", func(t *testing.T) {
		showImages := false
		sc := models.SearchConfig{ShowImages: &showImages}
		if sc.IsShowImages() {
			t.Error("IsShowImages() should return false when ShowImages is false")
		}
	})
}

func TestNewSearchConfig(t *testing.T) {
	sc := models.NewSearchConfig()

	if !sc.IsEnabled() {
		t.Error("New config should have search enabled")
	}

	if sc.Position != "navbar" {
		t.Errorf("Position = %v, want navbar", sc.Position)
	}

	if sc.Placeholder != "Search…" {
		t.Errorf("Placeholder = %v, want Search…", sc.Placeholder)
	}

	if !sc.IsShowImages() {
		t.Error("New config should have show_images enabled")
	}

	if sc.ExcerptLength != 200 {
		t.Errorf("ExcerptLength = %v, want 200", sc.ExcerptLength)
	}

	if sc.Pagefind.BundleDir != "_pagefind" {
		t.Errorf("Pagefind.BundleDir = %v, want _pagefind", sc.Pagefind.BundleDir)
	}

	if len(sc.Pagefind.ExcludeSelectors) != 0 {
		t.Errorf("Pagefind.ExcludeSelectors should be empty, got %v", sc.Pagefind.ExcludeSelectors)
	}

	if len(sc.Feeds) != 0 {
		t.Errorf("Feeds should be empty, got %v", sc.Feeds)
	}
}

func TestPagefindPlugin_CleanupSkipsWhenCorpusUnchanged(t *testing.T) {
	plugin := NewPagefindPlugin()
	m := lifecycle.NewManager()
	config := lifecycle.NewConfig()
	config.OutputDir = t.TempDir()
	config.Extra["search"] = models.NewSearchConfig()
	m.SetConfig(config)

	post := models.NewPost("pages/test.md")
	post.Published = true
	post.HTML = `<html><body><main data-pagefind-body>Test</main></body></html>`
	m.SetPosts([]*models.Post{post})

	bundleDir := defaultBundleDir
	indexPath := filepath.Join(config.OutputDir, bundleDir, "pagefind.js")
	if err := os.MkdirAll(filepath.Dir(indexPath), 0o755); err != nil {
		t.Fatalf("MkdirAll(indexPath) error = %v", err)
	}
	if err := os.WriteFile(indexPath, []byte("ok"), 0o600); err != nil {
		t.Fatalf("WriteFile(indexPath) error = %v", err)
	}

	cache := buildcache.New(t.TempDir())
	cache.SetPagefindCorpusHash(plugin.computeCorpusHash(m, models.NewSearchConfig()))
	m.Cache().Set("build_cache", cache)

	if err := plugin.Cleanup(m); err != nil {
		t.Fatalf("Cleanup() error = %v", err)
	}
}

// TestPagefindPlugin_InterfaceConformance verifies the plugin implements required interfaces.
func TestPagefindPlugin_InterfaceConformance(_ *testing.T) {
	plugin := NewPagefindPlugin()

	// Test lifecycle.Plugin
	var _ lifecycle.Plugin = plugin

	// Test lifecycle.CleanupPlugin
	var _ lifecycle.CleanupPlugin = plugin

	// Test lifecycle.PriorityPlugin
	var _ lifecycle.PriorityPlugin = plugin
}
