package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/blogroll"
)

func TestImportBlogrollOPMLPreservesExisting(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	path := filepath.Join(dir, "markata-go.toml")
	initial := "[markata-go.blogroll]\nenabled = true\n\n" +
		"[[markata-go.blogroll.feeds]]\nurl = \"https://one.test/feed\"\ntitle = \"Original\"\n# keep this comment\n"
	if err := os.WriteFile(path, []byte(initial), 0o600); err != nil {
		t.Fatal(err)
	}
	oldConfig, oldMerge := cfgFile, mergeConfigFiles
	t.Cleanup(func() { cfgFile, mergeConfigFiles = oldConfig, oldMerge })
	cfgFile, mergeConfigFiles = path, nil
	state, err := loadSettingsState()
	if err != nil {
		t.Fatal(err)
	}
	feeds := []blogroll.OPMLFeed{
		{URL: "https://one.test/feed", Title: "Do not replace"},
		{URL: "https://two.test/feed", Title: "New", Category: "Tech", SiteURL: "https://two.test"},
		{URL: "https://two.test/feed", Title: "Duplicate"},
	}
	added, skipped, err := importBlogrollOPML(state, feeds, true)
	if err != nil || added != 1 || skipped != 2 {
		t.Fatalf("unexpected dry run: %d, %d, %v", added, skipped, err)
	}
	data, err := os.ReadFile(path)
	if err != nil || string(data) != initial {
		t.Fatalf("dry run modified source: %v", err)
	}
	added, skipped, err = importBlogrollOPML(state, feeds, false)
	if err != nil || added != 1 || skipped != 2 {
		t.Fatalf("unexpected import: %d, %d, %v", added, skipped, err)
	}
	data, err = os.ReadFile(path)
	if err != nil || !strings.Contains(string(data), strings.TrimSpace(initial)) || strings.Contains(string(data), "Do not replace") {
		t.Fatalf("existing config changed: %v, %s", err, data)
	}
	state, err = loadSettingsState()
	if err != nil {
		t.Fatal(err)
	}
	added, skipped, err = importBlogrollOPML(state, feeds, false)
	if err != nil || added != 0 || skipped != 3 {
		t.Fatalf("import not idempotent: %d, %d, %v", added, skipped, err)
	}
}

func TestRenderOPMLImportInvalidBatch(t *testing.T) {
	feeds := []blogroll.OPMLFeed{
		{URL: "https://one.test/feed", Title: "One"},
		{URL: "https://two.test/feed", Title: strings.Repeat("x", 201)},
	}
	if _, err := renderOPMLImport("", feeds, true); err == nil {
		t.Fatal("expected invalid batch to be rejected before writing")
	}
}
