package plugins

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

func TestCSSBundlePlugin_Name(t *testing.T) {
	p := NewCSSBundlePlugin()
	if got := p.Name(); got != "css_bundle" {
		t.Errorf("Name() = %q, want %q", got, "css_bundle")
	}
}

func TestCSSBundlePlugin_InterfaceCompliance(_ *testing.T) {
	var _ lifecycle.Plugin = (*CSSBundlePlugin)(nil)
	var _ lifecycle.ConfigurePlugin = (*CSSBundlePlugin)(nil)
	var _ lifecycle.WritePlugin = (*CSSBundlePlugin)(nil)
	var _ lifecycle.PriorityPlugin = (*CSSBundlePlugin)(nil)
}

func TestCSSBundlePlugin_Priority(t *testing.T) {
	p := NewCSSBundlePlugin()

	// Should have late priority in Write stage (after CSS generators)
	if got := p.Priority(lifecycle.StageWrite); got != lifecycle.PriorityLate {
		t.Errorf("Priority(StageWrite) = %d, want %d", got, lifecycle.PriorityLate)
	}

	// Should have default priority in other stages
	if got := p.Priority(lifecycle.StageConfigure); got != lifecycle.PriorityDefault {
		t.Errorf("Priority(StageConfigure) = %d, want %d", got, lifecycle.PriorityDefault)
	}
}

func TestCSSBundlePlugin_Configure_NotEnabled(t *testing.T) {
	p := NewCSSBundlePlugin()
	m := lifecycle.NewManager()

	err := p.Configure(m)
	if err != nil {
		t.Errorf("Configure() error = %v, want nil", err)
	}

	// Should not be enabled without config
	if p.config.Enabled {
		t.Error("config.Enabled = true, want false")
	}
}

func TestCSSBundlePlugin_Configure_FromMap(t *testing.T) {
	p := NewCSSBundlePlugin()

	// Create config with css_bundle in Extra as a map (like TOML parsing)
	cfg := lifecycle.NewConfig()
	cfg.Extra = map[string]interface{}{
		"css_bundle": map[string]interface{}{
			"enabled": true,
			"minify":  false,
			"exclude": []interface{}{"debug.css", "test-*.css"},
			"bundles": []interface{}{
				map[string]interface{}{
					"name":    "main",
					"sources": []interface{}{"css/variables.css", "css/main.css"},
					"output":  "css/bundle.css",
				},
			},
		},
	}
	m := lifecycle.NewManager()
	m.SetConfig(cfg)

	err := p.Configure(m)
	if err != nil {
		t.Fatalf("Configure() error = %v", err)
	}

	if !p.config.Enabled {
		t.Error("config.Enabled = false, want true")
	}

	if len(p.config.Bundles) != 1 {
		t.Errorf("len(config.Bundles) = %d, want 1", len(p.config.Bundles))
	}

	if len(p.config.Exclude) != 2 {
		t.Errorf("len(config.Exclude) = %d, want 2", len(p.config.Exclude))
	}

	bundle := p.config.Bundles[0]
	if bundle.Name != "main" {
		t.Errorf("bundle.Name = %q, want %q", bundle.Name, "main")
	}
	if bundle.Output != "css/bundle.css" {
		t.Errorf("bundle.Output = %q, want %q", bundle.Output, "css/bundle.css")
	}
	if len(bundle.Sources) != 2 {
		t.Errorf("len(bundle.Sources) = %d, want 2", len(bundle.Sources))
	}
}

func TestCSSBundlePlugin_Write_NotEnabled(t *testing.T) {
	p := NewCSSBundlePlugin()
	m := lifecycle.NewManager()

	// Should skip without error when not enabled
	err := p.Write(m)
	if err != nil {
		t.Errorf("Write() error = %v, want nil", err)
	}
}

func TestCSSBundlePlugin_Write_NoBundles(t *testing.T) {
	p := NewCSSBundlePlugin()
	p.config.Enabled = true
	// No bundles configured

	m := lifecycle.NewManager()

	err := p.Write(m)
	if err != nil {
		t.Errorf("Write() error = %v, want nil", err)
	}
}

func TestCSSBundlePlugin_Write_CreatesBundles(t *testing.T) {
	// Create temp directory structure
	tmpDir := t.TempDir()
	cssDir := filepath.Join(tmpDir, "css")
	if err := os.MkdirAll(cssDir, 0o755); err != nil {
		t.Fatal(err)
	}

	// Create source CSS files
	varCSS := ":root { --color-primary: #007bff; }\n"
	mainCSS := "body { color: var(--color-primary); }\n"

	if err := os.WriteFile(filepath.Join(cssDir, "variables.css"), []byte(varCSS), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(cssDir, "main.css"), []byte(mainCSS), 0o600); err != nil {
		t.Fatal(err)
	}

	// Configure plugin
	p := NewCSSBundlePlugin()
	p.config = models.CSSBundleConfig{
		Enabled: true,
		Bundles: []models.BundleConfig{
			{
				Name:    "main",
				Sources: []string{"css/variables.css", "css/main.css"},
				Output:  "css/bundle.css",
			},
		},
	}

	addComments := true
	p.config.AddSourceComments = &addComments

	m := lifecycle.NewManager()
	cfg := lifecycle.NewConfig()
	cfg.OutputDir = tmpDir
	m.SetConfig(cfg)

	// Execute
	err := p.Write(m)
	if err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	// Verify bundle was created
	bundlePath := filepath.Join(tmpDir, "css", "bundle.css")
	content, err := os.ReadFile(bundlePath)
	if err != nil {
		t.Fatalf("Failed to read bundle: %v", err)
	}

	// Check bundle contains source CSS
	if !strings.Contains(string(content), "--color-primary: #007bff") {
		t.Error("Bundle missing variables.css content")
	}
	if !strings.Contains(string(content), "color: var(--color-primary)") {
		t.Error("Bundle missing main.css content")
	}

	// Check source comments are present
	if !strings.Contains(string(content), "Source: css/variables.css") {
		t.Error("Bundle missing source comment for variables.css")
	}
	if !strings.Contains(string(content), "Source: css/main.css") {
		t.Error("Bundle missing source comment for main.css")
	}

	// Check bundle paths in cache
	if bundles, ok := m.Cache().Get("css_bundles"); ok {
		bundleMap, ok := bundles.(map[string]string)
		if !ok {
			t.Error("css_bundles is not map[string]string")
		} else if bundleMap["main"] != "/css/bundle.css" {
			t.Errorf("cache css_bundles[main] = %q, want %q", bundleMap["main"], "/css/bundle.css")
		}
	} else {
		t.Error("css_bundles not in cache")
	}

	// Check bundling enabled flag in cache
	if enabled, ok := m.Cache().Get("css_bundling_enabled"); ok {
		if enabled != true {
			t.Error("css_bundling_enabled = false, want true")
		}
	} else {
		t.Error("css_bundling_enabled not in cache")
	}
}

func TestCSSBundlePlugin_Write_NoSourceComments(t *testing.T) {
	// Create temp directory structure
	tmpDir := t.TempDir()
	cssDir := filepath.Join(tmpDir, "css")
	if err := os.MkdirAll(cssDir, 0o755); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join(cssDir, "test.css"), []byte("body { margin: 0; }\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	// Configure plugin with source comments disabled
	p := NewCSSBundlePlugin()
	addComments := false
	p.config = models.CSSBundleConfig{
		Enabled:           true,
		AddSourceComments: &addComments,
		Bundles: []models.BundleConfig{
			{
				Name:    "test",
				Sources: []string{"css/test.css"},
				Output:  "css/test-bundle.css",
			},
		},
	}

	m := lifecycle.NewManager()
	cfg := lifecycle.NewConfig()
	cfg.OutputDir = tmpDir
	m.SetConfig(cfg)
	// Config set

	// (config already set above)

	err := p.Write(m)
	if err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	// Check that source comments are NOT present
	content, err := os.ReadFile(filepath.Join(tmpDir, "css", "test-bundle.css"))
	if err != nil {
		t.Fatal(err)
	}

	if strings.Contains(string(content), "=== Source:") {
		t.Error("Bundle should not contain source comments when disabled")
	}
}

func TestCSSBundlePlugin_Write_ExcludesFiles(t *testing.T) {
	// Create temp directory structure
	tmpDir := t.TempDir()
	cssDir := filepath.Join(tmpDir, "css")
	if err := os.MkdirAll(cssDir, 0o755); err != nil {
		t.Fatal(err)
	}

	// Create CSS files including one that should be excluded
	if err := os.WriteFile(filepath.Join(cssDir, "main.css"), []byte("main { }\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(cssDir, "debug.css"), []byte("debug { }\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	// Configure plugin with exclusion
	p := NewCSSBundlePlugin()
	p.config = models.CSSBundleConfig{
		Enabled: true,
		Exclude: []string{"debug.css"},
		Bundles: []models.BundleConfig{
			{
				Name:    "main",
				Sources: []string{"css/*.css"},
				Output:  "css/bundle.css",
			},
		},
	}

	// Build exclude map
	p.exclude["debug.css"] = true

	m := lifecycle.NewManager()
	cfg := lifecycle.NewConfig()
	cfg.OutputDir = tmpDir
	m.SetConfig(cfg)
	// Config set

	// (config already set above)

	err := p.Write(m)
	if err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	// Check bundle content
	content, err := os.ReadFile(filepath.Join(tmpDir, "css", "bundle.css"))
	if err != nil {
		t.Fatal(err)
	}

	// Should have main.css
	if !strings.Contains(string(content), "main {") {
		t.Error("Bundle missing main.css content")
	}

	// Should NOT have debug.css
	if strings.Contains(string(content), "debug {") {
		t.Error("Bundle should not contain excluded debug.css content")
	}
}

func TestCSSBundlePlugin_Write_GlobPattern(t *testing.T) {
	// Create temp directory structure
	tmpDir := t.TempDir()
	cssDir := filepath.Join(tmpDir, "css")
	if err := os.MkdirAll(cssDir, 0o755); err != nil {
		t.Fatal(err)
	}

	// Create multiple CSS files
	if err := os.WriteFile(filepath.Join(cssDir, "a.css"), []byte(".a { }\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(cssDir, "b.css"), []byte(".b { }\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(cssDir, "c.css"), []byte(".c { }\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	p := NewCSSBundlePlugin()
	p.config = models.CSSBundleConfig{
		Enabled: true,
		Bundles: []models.BundleConfig{
			{
				Name:    "all",
				Sources: []string{"css/*.css"},
				Output:  "css/all.css",
			},
		},
	}

	m := lifecycle.NewManager()
	cfg := lifecycle.NewConfig()
	cfg.OutputDir = tmpDir
	m.SetConfig(cfg)
	// Config set

	// (config already set above)

	err := p.Write(m)
	if err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	// Check bundle contains all files
	content, err := os.ReadFile(filepath.Join(tmpDir, "css", "all.css"))
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(content), ".a {") {
		t.Error("Bundle missing a.css")
	}
	if !strings.Contains(string(content), ".b {") {
		t.Error("Bundle missing b.css")
	}
	if !strings.Contains(string(content), ".c {") {
		t.Error("Bundle missing c.css")
	}
}

func TestCSSBundlePlugin_Write_MultipleBundles(t *testing.T) {
	// Create temp directory structure
	tmpDir := t.TempDir()
	cssDir := filepath.Join(tmpDir, "css")
	if err := os.MkdirAll(cssDir, 0o755); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join(cssDir, "critical.css"), []byte("critical { }\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(cssDir, "main.css"), []byte("main { }\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	p := NewCSSBundlePlugin()
	p.config = models.CSSBundleConfig{
		Enabled: true,
		Bundles: []models.BundleConfig{
			{
				Name:    "critical",
				Sources: []string{"css/critical.css"},
				Output:  "css/critical-bundle.css",
			},
			{
				Name:    "main",
				Sources: []string{"css/main.css"},
				Output:  "css/main-bundle.css",
			},
		},
	}

	m := lifecycle.NewManager()
	cfg := lifecycle.NewConfig()
	cfg.OutputDir = tmpDir
	m.SetConfig(cfg)
	// Config set

	// (config already set above)

	err := p.Write(m)
	if err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	// Check both bundles exist
	if _, err := os.Stat(filepath.Join(tmpDir, "css", "critical-bundle.css")); os.IsNotExist(err) {
		t.Error("critical-bundle.css not created")
	}
	if _, err := os.Stat(filepath.Join(tmpDir, "css", "main-bundle.css")); os.IsNotExist(err) {
		t.Error("main-bundle.css not created")
	}

	// Check cache has both bundles
	if bundles, ok := m.Cache().Get("css_bundles"); ok {
		bundleMap, ok := bundles.(map[string]string)
		if !ok {
			t.Error("css_bundles is not map[string]string")
		} else if len(bundleMap) != 2 {
			t.Errorf("cache has %d bundles, want 2", len(bundleMap))
		}
	}
}

func TestCSSBundlePlugin_Write_MissingSourceFile(t *testing.T) {
	tmpDir := t.TempDir()
	cssDir := filepath.Join(tmpDir, "css")
	if err := os.MkdirAll(cssDir, 0o755); err != nil {
		t.Fatal(err)
	}

	// Create only one of the two source files
	if err := os.WriteFile(filepath.Join(cssDir, "exists.css"), []byte("exists { }\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	p := NewCSSBundlePlugin()
	p.config = models.CSSBundleConfig{
		Enabled: true,
		Bundles: []models.BundleConfig{
			{
				Name:    "partial",
				Sources: []string{"css/exists.css", "css/missing.css"},
				Output:  "css/partial.css",
			},
		},
	}

	m := lifecycle.NewManager()
	cfg := lifecycle.NewConfig()
	cfg.OutputDir = tmpDir
	m.SetConfig(cfg)
	// Config set

	// (config already set above)

	// Should not fail - just skip missing files
	err := p.Write(m)
	if err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	// Bundle should still be created with the file that exists
	content, err := os.ReadFile(filepath.Join(tmpDir, "css", "partial.css"))
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(content), "exists {") {
		t.Error("Bundle missing content from existing file")
	}
}

func TestCSSBundlePlugin_parseConfigFromMap_Empty(t *testing.T) {
	p := NewCSSBundlePlugin()
	cfg := p.parseConfigFromMap(map[string]interface{}{})

	if cfg.Enabled {
		t.Error("Enabled should default to false")
	}
	if len(cfg.Bundles) != 0 {
		t.Error("Bundles should be empty")
	}
}

func TestCSSBundlePlugin_isExcluded_ExactMatch(t *testing.T) {
	p := NewCSSBundlePlugin()
	p.exclude["debug.css"] = true

	if !p.isExcluded("/path/to/debug.css") {
		t.Error("debug.css should be excluded")
	}
	if p.isExcluded("/path/to/main.css") {
		t.Error("main.css should not be excluded")
	}
}

func TestCSSBundlePlugin_isExcluded_WildcardMatch(t *testing.T) {
	p := NewCSSBundlePlugin()
	p.exclude["test-*.css"] = true

	if !p.isExcluded("/path/to/test-something.css") {
		t.Error("test-something.css should be excluded")
	}
	if p.isExcluded("/path/to/main.css") {
		t.Error("main.css should not be excluded")
	}
}
