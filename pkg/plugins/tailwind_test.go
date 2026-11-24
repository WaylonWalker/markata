package plugins

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

type stubTailwindInstaller struct {
	path string
	err  error
}

func (s stubTailwindInstaller) Install() (string, error) {
	return s.path, s.err
}

func TestTailwindPlugin_GeneratedTailwindContentPaths(t *testing.T) {
	plugin := NewTailwindPlugin()
	tmpDir := t.TempDir()
	assetsDir := filepath.Join(tmpDir, "static")
	templatesDir := filepath.Join(tmpDir, "templates")
	if err := os.MkdirAll(assetsDir, 0o755); err != nil {
		t.Fatalf("MkdirAll(assetsDir) error = %v", err)
	}
	if err := os.MkdirAll(templatesDir, 0o755); err != nil {
		t.Fatalf("MkdirAll(templatesDir) error = %v", err)
	}

	config := &lifecycle.Config{
		Extra: map[string]interface{}{
			"assets_dir":    assetsDir,
			"templates_dir": templatesDir,
		},
	}

	patterns := plugin.generatedTailwindContentPaths(config, "/tmp/manifest.txt")
	want := []string{"/tmp/manifest.txt", filepath.ToSlash(filepath.Join(assetsDir, "**", "*.js")), filepath.ToSlash(filepath.Join(templatesDir, "**", "*.html")), filepath.ToSlash(filepath.Join(templatesDir, "**", "*.js")), filepath.ToSlash(filepath.Join(templatesDir, "**", "*.md"))}
	for _, needle := range want {
		found := false
		for _, pattern := range patterns {
			if pattern == needle {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("patterns %v missing %q", patterns, needle)
		}
	}
}

func TestTailwindPlugin_ResolveBuildConfigFile_GeneratesContentConfig(t *testing.T) {
	plugin := NewTailwindPlugin()
	plugin.config = models.NewTailwindConfig()

	config := &lifecycle.Config{}

	configPath, cleanup, err := plugin.resolveBuildConfigFile(config, []string{"/tmp/manifest.txt", "/tmp/templates/**/*.html"})
	if err != nil {
		t.Fatalf("resolveBuildConfigFile() error = %v", err)
	}
	defer cleanup()

	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", configPath, err)
	}
	text := string(data)
	for _, needle := range []string{"module.exports", "/tmp/manifest.txt", "/tmp/templates/**/*.html"} {
		if !strings.Contains(text, needle) {
			t.Fatalf("generated config missing %q:\n%s", needle, text)
		}
	}
	if !strings.Contains(text, "preflight: false") {
		t.Fatalf("generated config should disable preflight by default:\n%s", text)
	}
}

func TestTailwindPlugin_ResolveBuildConfigFile_AllowsPreflightOptIn(t *testing.T) {
	plugin := NewTailwindPlugin()
	plugin.config = models.NewTailwindConfig()
	preflight := true
	plugin.config.Preflight = &preflight

	configPath, cleanup, err := plugin.resolveBuildConfigFile(&lifecycle.Config{}, []string{"/tmp/manifest.txt"})
	if err != nil {
		t.Fatalf("resolveBuildConfigFile() error = %v", err)
	}
	defer cleanup()

	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", configPath, err)
	}
	text := string(data)
	if !strings.Contains(text, "preflight: true") {
		t.Fatalf("generated config should preserve explicit preflight opt-in:\n%s", text)
	}
}

func TestTailwindPlugin_ManifestHashChangesWhenPreflightChanges(t *testing.T) {
	tmpDir := t.TempDir()
	plugin := NewTailwindPlugin()
	plugin.config = models.NewTailwindConfig()
	config := &lifecycle.Config{Extra: map[string]interface{}{"assets_dir": tmpDir}}

	baseHash := plugin.computeTailwindManifestHash(config, "tokens")

	preflight := true
	plugin.config.Preflight = &preflight
	preflightHash := plugin.computeTailwindManifestHash(config, "tokens")

	if baseHash == preflightHash {
		t.Fatal("expected manifest hash to change when preflight setting changes")
	}
}

func TestExtractTailwindTokens(t *testing.T) {
	html := `<main class="prose prose-zinc dark:prose-invert"><div class='grid grid-cols-2 gap-4'></div><p class="prose prose-zinc"></p></main>`
	got := extractTailwindTokens(html)
	for _, needle := range []string{"dark:prose-invert", "gap-4", "grid", "grid-cols-2", "prose", "prose-zinc"} {
		if !strings.Contains(got, needle) {
			t.Fatalf("extractTailwindTokens() missing %q in %q", needle, got)
		}
	}
}

func TestTailwindPlugin_IncludeAssetPath_RejectsAbsoluteOutput(t *testing.T) {
	plugin := NewTailwindPlugin()
	config := &lifecycle.Config{Extra: map[string]interface{}{"assets_dir": "/repo/static"}}

	got := plugin.includeAssetPath(config, "/tmp/markata-tailwind.css")
	if got != "" {
		t.Fatalf("includeAssetPath() = %q, want empty string for absolute output outside assets dir", got)
	}
}

func TestTailwindPlugin_IncludeAssetPath_AcceptsAbsoluteOutputInsideAssetsDir(t *testing.T) {
	plugin := NewTailwindPlugin()
	tmpDir := t.TempDir()
	assetsDir := filepath.Join(tmpDir, "static")
	config := &lifecycle.Config{Extra: map[string]interface{}{"assets_dir": assetsDir}}

	got := plugin.includeAssetPath(config, filepath.Join(assetsDir, "css", "markata-tailwind.css"))
	if got != "css/markata-tailwind.css" {
		t.Fatalf("includeAssetPath() = %q, want css/markata-tailwind.css", got)
	}
}

func TestIsAbsoluteOrRootedPath(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{name: "unix absolute", input: "/repo/pages/**/*.md", want: true},
		{name: "windows absolute", input: `C:\repo\pages\**\*.md`, want: true},
		{name: "windows rooted", input: `\repo\pages\**\*.md`, want: true},
		{name: "relative", input: "pages/**/*.md", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isAbsoluteOrRootedPath(tt.input); got != tt.want {
				t.Fatalf("isAbsoluteOrRootedPath(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestTailwindPlugin_ResolveBuildInput_GeneratesDefaultInput(t *testing.T) {
	plugin := NewTailwindPlugin()
	plugin.config = models.NewTailwindConfig()

	tmpDir := t.TempDir()
	config := &lifecycle.Config{Extra: map[string]interface{}{"assets_dir": filepath.Join(tmpDir, "static")}}

	inputPath, cleanup, err := plugin.resolveBuildInput(config)
	if err != nil {
		t.Fatalf("resolveBuildInput() error = %v", err)
	}
	defer cleanup()

	data, err := os.ReadFile(inputPath)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", inputPath, err)
	}
	if string(data) != tailwindDefaultInputCSS {
		t.Fatalf("generated input = %q, want %q", string(data), tailwindDefaultInputCSS)
	}
}

func TestTailwindPlugin_FindOrInstallTailwind_PrefersPathWhenAvailable(t *testing.T) {
	plugin := NewTailwindPlugin()
	plugin.config = models.NewTailwindConfig()
	plugin.config.Version = "v3.4.19"

	origLookPath := tailwindLookPath
	origInstaller := newTailwindInstaller
	t.Cleanup(func() {
		tailwindLookPath = origLookPath
		newTailwindInstaller = origInstaller
	})

	lookups := 0
	tailwindLookPath = func(_ string) (string, error) {
		lookups++
		return "/usr/bin/tailwindcss", nil
	}

	installed := false
	newTailwindInstaller = func(config TailwindInstallerConfig) tailwindInstaller {
		installed = true
		if config.Version != "v3.4.19" {
			t.Fatalf("installer version = %q, want v3.4.19", config.Version)
		}
		return stubTailwindInstaller{path: "/managed/tailwindcss"}
	}

	path, err := plugin.findOrInstallTailwind()
	if err != nil {
		t.Fatalf("findOrInstallTailwind() error = %v", err)
	}
	if path != "/usr/bin/tailwindcss" {
		t.Fatalf("findOrInstallTailwind() = %q, want /usr/bin/tailwindcss", path)
	}
	if installed {
		t.Fatal("expected system tailwindcss binary to be used before auto-install")
	}
	if lookups != 1 {
		t.Fatalf("expected one PATH lookup when auto-install enabled, got %d lookups", lookups)
	}
}

func TestTailwindPlugin_FindOrInstallTailwind_UsesPathWhenAutoInstallDisabled(t *testing.T) {
	plugin := NewTailwindPlugin()
	plugin.config = models.NewTailwindConfig()
	autoInstall := false
	plugin.config.AutoInstall = &autoInstall

	origLookPath := tailwindLookPath
	origInstaller := newTailwindInstaller
	t.Cleanup(func() {
		tailwindLookPath = origLookPath
		newTailwindInstaller = origInstaller
	})

	tailwindLookPath = func(_ string) (string, error) {
		return "/usr/bin/tailwindcss", nil
	}
	newTailwindInstaller = func(_ TailwindInstallerConfig) tailwindInstaller {
		t.Fatal("installer should not be used when auto_install is false")
		return stubTailwindInstaller{}
	}

	path, err := plugin.findOrInstallTailwind()
	if err != nil {
		t.Fatalf("findOrInstallTailwind() error = %v", err)
	}
	if path != "/usr/bin/tailwindcss" {
		t.Fatalf("findOrInstallTailwind() = %q, want /usr/bin/tailwindcss", path)
	}
}

func TestTailwindPlugin_FindOrInstallTailwind_UsesManagedInstallerWhenPathMissing(t *testing.T) {
	plugin := NewTailwindPlugin()
	plugin.config = models.NewTailwindConfig()
	plugin.config.Version = "v3.4.19"

	origLookPath := tailwindLookPath
	origInstaller := newTailwindInstaller
	t.Cleanup(func() {
		tailwindLookPath = origLookPath
		newTailwindInstaller = origInstaller
	})

	tailwindLookPath = func(_ string) (string, error) {
		return "", os.ErrNotExist
	}

	newTailwindInstaller = func(config TailwindInstallerConfig) tailwindInstaller {
		if config.Version != "v3.4.19" {
			t.Fatalf("installer version = %q, want v3.4.19", config.Version)
		}
		return stubTailwindInstaller{path: "/managed/tailwindcss"}
	}

	path, err := plugin.findOrInstallTailwind()
	if err != nil {
		t.Fatalf("findOrInstallTailwind() error = %v", err)
	}
	if path != "/managed/tailwindcss" {
		t.Fatalf("findOrInstallTailwind() = %q, want /managed/tailwindcss", path)
	}
}
