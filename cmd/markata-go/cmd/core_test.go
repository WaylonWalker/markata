package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/WaylonWalker/markata/pkg/buildstats"
	"github.com/WaylonWalker/markata/pkg/diagnostics"
	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/plugins"
	"github.com/spf13/cobra"
)

func TestResolveConfigBaseDir(t *testing.T) {
	configDir := t.TempDir()
	configPath := filepath.Join(configDir, "markata-go.toml")
	if err := os.WriteFile(configPath, []byte("title = \"test\"\n"), 0o600); err != nil {
		t.Fatalf("WriteFile(config) error = %v", err)
	}

	got := resolveConfigBaseDir(configPath)
	want := filepath.Clean(configDir)
	if got != want {
		t.Fatalf("resolveConfigBaseDir() = %q, want %q", got, want)
	}
}

func TestResolveConfigRelativePath(t *testing.T) {
	baseDir := t.TempDir()
	absOut := filepath.Join(t.TempDir(), "out")

	tests := []struct {
		name string
		path string
		want string
	}{
		{name: "relative path", path: "output", want: filepath.Join(baseDir, "output")},
		{name: "absolute path", path: absOut, want: absOut},
		{name: "empty path", path: "", want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := resolveConfigRelativePath(baseDir, tt.path); got != tt.want {
				t.Fatalf("resolveConfigRelativePath(%q) = %q, want %q", tt.path, got, tt.want)
			}
		})
	}
}

func TestCreateManagerResolvesFontpacksFileRelativeToConfig(t *testing.T) {
	root := t.TempDir()
	site := filepath.Join(root, "site")
	if err := os.MkdirAll(filepath.Join(site, "fonts"), 0o755); err != nil {
		t.Fatal(err)
	}
	catalog := "schema: markata.fontpacks/v2\nsystem_stacks: {sans: {css: system-ui}}\nfont_sources: {}\nfontpacks: {system: {performance: {class: zero-download}, roles: {body: {stack: sans}}}}\ncatalog: {bundled_asset_root: .}\n"
	catalogPath := filepath.Join(site, "fonts", "catalog.yaml")
	if err := os.WriteFile(catalogPath, []byte(catalog), 0o600); err != nil {
		t.Fatal(err)
	}
	configPath := filepath.Join(site, "markata-go.yaml")
	config := "markata-go:\n  fontpacks_file: fonts/catalog.yaml\n"
	if err := os.WriteFile(configPath, []byte(config), 0o600); err != nil {
		t.Fatal(err)
	}
	manager, err := createManager(configPath)
	if err != nil {
		t.Fatal(err)
	}
	got, ok := manager.Config().Extra["fontpacks_file"].(string)
	if !ok || got != catalogPath {
		t.Fatalf("fontpacks_file = %q, want %q", got, catalogPath)
	}
}

func TestCreateManagerPreservesAbsoluteFontpacksFile(t *testing.T) {
	root := t.TempDir()
	catalogPath := filepath.Join(root, "catalog.yaml")
	catalog := "schema: markata.fontpacks/v2\nsystem_stacks: {sans: {css: system-ui}}\nfont_sources: {}\nfontpacks: {system: {performance: {class: zero-download}, roles: {body: {stack: sans}}}}\ncatalog: {bundled_asset_root: .}\n"
	if err := os.WriteFile(catalogPath, []byte(catalog), 0o600); err != nil {
		t.Fatal(err)
	}
	configPath := filepath.Join(root, "markata-go.yaml")
	if err := os.WriteFile(configPath, []byte("markata-go:\n  fontpacks_file: "+catalogPath+"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	manager, err := createManager(configPath)
	if err != nil {
		t.Fatal(err)
	}
	if got := manager.Config().Extra["fontpacks_file"]; got != catalogPath {
		t.Fatalf("fontpacks_file = %q, want %q", got, catalogPath)
	}
}

func TestCreateManagerRegistersConfiguredOptionalHook(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "markata-go.toml")
	config := "[markata-go]\nhooks = [\"default\", \"image_optimization\"]\n"
	if err := os.WriteFile(configPath, []byte(config), 0o600); err != nil {
		t.Fatal(err)
	}

	manager, err := createManager(configPath)
	if err != nil {
		t.Fatal(err)
	}
	for _, plugin := range manager.Plugins() {
		if plugin.Name() == "image_optimization" {
			if _, ok := plugin.(*plugins.ImageOptimizationPlugin); !ok {
				t.Fatalf("image_optimization plugin = %T, want *ImageOptimizationPlugin", plugin)
			}
			return
		}
	}
	t.Fatal("configured image_optimization hook was not registered")
}

func TestCreateManagerDoesNotRegisterUnconfiguredOptionalHook(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "markata-go.toml")
	if err := os.WriteFile(configPath, []byte("[markata-go]\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	manager, err := createManager(configPath)
	if err != nil {
		t.Fatal(err)
	}
	for _, plugin := range manager.Plugins() {
		if plugin.Name() == "image_optimization" {
			t.Fatal("unconfigured image_optimization hook was registered")
		}
	}
}

func TestCreateManagerHonorsDisabledOptionalHook(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "markata-go.toml")
	config := "[markata-go]\nhooks = [\"default\", \"image_optimization\"]\ndisabled_hooks = [\"image_optimization\"]\n"
	if err := os.WriteFile(configPath, []byte(config), 0o600); err != nil {
		t.Fatal(err)
	}

	manager, err := createManager(configPath)
	if err != nil {
		t.Fatal(err)
	}
	for _, plugin := range manager.Plugins() {
		if plugin.Name() == "image_optimization" {
			t.Fatal("disabled image_optimization hook was registered")
		}
	}
}

func TestCreateManagerProjectsAffectedConfig(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "markata-go.toml")
	config := `[markata-go]
[markata-go.head]
text = "typed head"
[markata-go.theme_calendar]
enabled = false
[markata-go.error_pages]
enable_404 = false
[markata-go.resource_hints]
enabled = false
[markata-go.markdown.highlight]
enabled = false
theme = "monokai"
`
	if err := os.WriteFile(configPath, []byte(config), 0o600); err != nil {
		t.Fatal(err)
	}

	manager, err := createManager(configPath)
	if err != nil {
		t.Fatal(err)
	}

	extra := manager.Config().Extra
	if extra["markata_version"] != Version || extra["markata_commit"] != Commit {
		t.Errorf("build identity = (%v, %v), want (%q, %q)", extra["markata_version"], extra["markata_commit"], Version, Commit)
	}
	if head, ok := extra["head"].(models.HeadConfig); !ok || head.Text != "typed head" {
		t.Errorf("head extra = %#v, want typed head config", extra["head"])
	}
	if calendar, ok := extra["theme_calendar"].(models.ThemeCalendarConfig); !ok || calendar.IsEnabled() {
		t.Errorf("theme_calendar extra = %#v, want disabled typed config", extra["theme_calendar"])
	}
	if errorPages, ok := extra["error_pages"].(models.ErrorPagesConfig); !ok || errorPages.Is404Enabled() {
		t.Errorf("error_pages extra = %#v, want disabled typed config", extra["error_pages"])
	}
	if resourceHints, ok := extra["resource_hints"].(models.ResourceHintsConfig); !ok || resourceHints.IsEnabled() {
		t.Errorf("resource_hints extra = %#v, want disabled typed config", extra["resource_hints"])
	}
	if markdown, ok := extra["markdown"].(models.MarkdownConfig); !ok || markdown.Highlight.IsEnabled() || markdown.Highlight.Theme != "monokai" {
		t.Errorf("markdown extra = %#v, want typed highlight config", extra["markdown"])
	}
}

func TestPrintBuildResult_IncludesBenchmarkSummary(t *testing.T) {
	stdout := bytes.NewBuffer(nil)
	command := &cobra.Command{Use: "build"}
	command.SetOut(stdout)
	currentCmd = command
	defer func() { currentCmd = nil }()

	printBuildResult(&BuildResult{
		PostsProcessed: 12,
		FeedsGenerated: 3,
		Duration:       9.87,
		Benchmark: buildstats.Summary{
			Total: 10 * time.Second,
			Resources: buildstats.ResourceBreakdown{
				CPU:           2 * time.Second,
				NetworkWait:   5 * time.Second,
				DiskReadWait:  750 * time.Millisecond,
				DiskWriteWait: 1250 * time.Millisecond,
				Idle:          1 * time.Second,
			},
			Hotspots: []buildstats.Hotspot{
				{Stage: "collect", Plugin: "blogroll", Duration: 3 * time.Second},
				{Stage: "render", Plugin: "link_avatars", Duration: 2 * time.Second},
			},
			Requests: []buildstats.RequestTiming{
				{Stage: "collect", Plugin: "blogroll", Method: "GET", URL: "https://example.com/feed.xml", Duration: 4 * time.Second, Status: 200},
			},
		},
	})

	output := stdout.String()
	for _, want := range []string{
		"Build completed successfully!",
		"Posts processed: 12",
		"Feeds generated: 3",
		"Content:",
		"Resource profile: estimated wall time",
		"CPU",
		"Network wait",
		"Disk read",
		"Disk write",
		"Idle",
		"Hotspots:",
		"collect/blogroll",
		"render/link_avatars",
		"Slowest requests:",
		"GET https://example.com/feed.xml",
		"Duration: 9.87s",
	} {
		if !strings.Contains(output, want) {
			t.Fatalf("output missing %q:\n%s", want, output)
		}
	}
}

func TestWriteBenchmarkJSON_IncludesContentDiagnostics(t *testing.T) {
	var output bytes.Buffer
	result := &BuildResult{
		Content: diagnostics.ContentLedgerSnapshot{
			Summary: diagnostics.ContentSummary{Discovered: 1, Candidates: 1, Warnings: 1},
			Entries: []diagnostics.ContentDisposition{{
				Path:        "post.md",
				Candidate:   true,
				Disposition: diagnostics.DispositionShadow,
				Reasons:     []string{diagnostics.ReasonContentPublishedFalse},
			}},
		},
	}

	if err := writeBenchmarkJSON(&output, result); err != nil {
		t.Fatalf("writeBenchmarkJSON() error = %v", err)
	}

	var decoded benchmarkJSONOutput
	if err := json.Unmarshal(output.Bytes(), &decoded); err != nil {
		t.Fatalf("benchmark JSON = %v", err)
	}
	if decoded.Content.Summary != result.Content.Summary {
		t.Fatalf("content summary = %+v, want %+v", decoded.Content.Summary, result.Content.Summary)
	}
	if len(decoded.Content.Entries) != 1 || decoded.Content.Entries[0].Path != "post.md" {
		t.Fatalf("content entries = %+v", decoded.Content.Entries)
	}
}
