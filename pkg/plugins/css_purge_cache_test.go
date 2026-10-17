package plugins

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/WaylonWalker/markata-go/pkg/csspurge"
	"github.com/WaylonWalker/markata-go/pkg/models"
)

func TestCSSPurgeCache_ReusesOnlyMatchingInputs(t *testing.T) {
	root := t.TempDir()
	output := filepath.Join(root, "output")
	cache := filepath.Join(root, "cache")
	if err := os.MkdirAll(output, 0o755); err != nil {
		t.Fatal(err)
	}
	css := filepath.Join(output, "site.css")
	if err := os.WriteFile(css, []byte(".used{color:blue}.unused{color:red}"), 0o644); err != nil {
		t.Fatal(err)
	}
	selectors := csspurge.NewUsedSelectors()
	selectors.Classes["used"] = true
	config := models.NewCSSPurgeConfig()
	options := buildPurgeOptions(config, false)

	first := processCSSFiles([]string{css}, output, cache, selectors, options, config, false)
	if first.cacheHits != 0 {
		t.Fatalf("first run had %d cache hits", first.cacheHits)
	}
	second := processCSSFiles([]string{css}, output, cache, selectors, options, config, false)
	if second.cacheHits != 1 {
		t.Fatalf("second run had %d cache hits, want 1", second.cacheHits)
	}
	if err := os.WriteFile(css, []byte(".used{color:blue}.new{color:green}"), 0o644); err != nil {
		t.Fatal(err)
	}
	changed := processCSSFiles([]string{css}, output, cache, selectors, options, config, false)
	if changed.cacheHits != 0 {
		t.Fatal("stylesheet change reused stale cache")
	}

	selectors.Classes["other"] = true
	changed = processCSSFiles([]string{css}, output, cache, selectors, options, config, false)
	if changed.cacheHits != 0 {
		t.Fatal("selector change reused stale cache")
	}

	options.Preserve = append(options.Preserve, "dynamic-*")
	changed = processCSSFiles([]string{css}, output, cache, selectors, options, config, false)
	if changed.cacheHits != 0 {
		t.Fatal("purge option change reused stale cache")
	}
}

func TestCSSPurgeCache_CorruptEntryFallsBackToPurge(t *testing.T) {
	root := t.TempDir()
	output := filepath.Join(root, "output")
	cache := filepath.Join(root, "cache")
	if err := os.MkdirAll(output, 0o755); err != nil {
		t.Fatal(err)
	}
	css := filepath.Join(output, "site.css")
	if err := os.WriteFile(css, []byte(".used{color:blue}.unused{color:red}"), 0o644); err != nil {
		t.Fatal(err)
	}
	entry := cssPurgeCachePath(cache, output, css)
	if err := os.MkdirAll(filepath.Dir(entry), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(entry, []byte("not json"), 0o644); err != nil {
		t.Fatal(err)
	}
	selectors := csspurge.NewUsedSelectors()
	selectors.Classes["used"] = true
	config := models.NewCSSPurgeConfig()
	stats := processCSSFiles([]string{css}, output, cache, selectors, buildPurgeOptions(config, false), config, false)
	if stats.cacheHits != 0 {
		t.Fatal("corrupt cache entry was treated as a hit")
	}
	content, err := os.ReadFile(css)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != ".used{color:blue}\n" {
		t.Fatalf("CSS was not purged: %q", content)
	}
}
