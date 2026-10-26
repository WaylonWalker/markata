package plugins

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/WaylonWalker/markata/pkg/csspurge"
	"github.com/WaylonWalker/markata/pkg/models"
)

func TestCSSPurgeCache_ReusesOnlyMatchingInputs(t *testing.T) {
	root := t.TempDir()
	output := filepath.Join(root, "output")
	cache := filepath.Join(root, "cache")
	if err := os.MkdirAll(output, 0o755); err != nil {
		t.Fatal(err)
	}
	css := filepath.Join(output, "site.css")
	if err := os.WriteFile(css, []byte(".used{color:blue}.unused{color:red}"), 0o600); err != nil {
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
	// A full rebuild rewrites the original (unpurged) CSS. This must be a
	// cache hit too, and must restore the previous purged bytes.
	if err := os.WriteFile(css, []byte(".used{color:blue}.unused{color:red}"), 0o600); err != nil {
		t.Fatal(err)
	}
	rebuilt := processCSSFiles([]string{css}, output, cache, selectors, options, config, false)
	if rebuilt.cacheHits != 1 {
		t.Fatalf("full rebuild had %d cache hits, want 1", rebuilt.cacheHits)
	}
	restored, err := os.ReadFile(css)
	if err != nil {
		t.Fatal(err)
	}
	if string(restored) != ".used{color:blue}\n" {
		t.Fatalf("cache did not restore purged CSS: %q", restored)
	}
	if err := os.WriteFile(css, []byte(".used{color:blue}.new{color:green}"), 0o600); err != nil {
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
	if err := os.WriteFile(css, []byte(".used{color:blue}.unused{color:red}"), 0o600); err != nil {
		t.Fatal(err)
	}
	entry := cssPurgeCachePath(cache, output, css)
	if err := os.MkdirAll(filepath.Dir(entry), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(entry, []byte("not json"), 0o600); err != nil {
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

func TestCSSPurgeCache_TamperedOutputFallsBackToPurge(t *testing.T) {
	root := t.TempDir()
	output, cache := filepath.Join(root, "output"), filepath.Join(root, "cache")
	if err := os.MkdirAll(output, 0o755); err != nil {
		t.Fatal(err)
	}
	css := filepath.Join(output, "site.css")
	original := []byte(".used{color:blue}.unused{color:red}")
	if err := os.WriteFile(css, original, 0o600); err != nil {
		t.Fatal(err)
	}
	selectors := csspurge.NewUsedSelectors()
	selectors.Classes["used"] = true
	config := models.NewCSSPurgeConfig()
	opts := buildPurgeOptions(config, false)
	if first := processCSSFiles([]string{css}, output, cache, selectors, opts, config, false); first.cacheHits != 0 {
		t.Fatal("first pass unexpectedly hit cache")
	}
	entryPath := cssPurgeCachePath(cache, output, css)
	entry, ok := loadCSSPurgeCache(entryPath)
	if !ok {
		t.Fatal("expected a valid cache entry")
	}
	entry.OutputCSS = ".tampered{color:red}"
	saveCSSPurgeCache(entryPath, entry)
	if err := os.WriteFile(css, original, 0o600); err != nil {
		t.Fatal(err)
	}
	if again := processCSSFiles([]string{css}, output, cache, selectors, opts, config, false); again.cacheHits != 0 {
		t.Fatal("tampered cached CSS was reused")
	}
	content, err := os.ReadFile(css)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != ".used{color:blue}\n" {
		t.Fatalf("bad output after corrupted cache: %q", content)
	}
}

func TestCSSPurgeCache_NoRemovalPreservesOriginalBytes(t *testing.T) {
	root := t.TempDir()
	output, cache := filepath.Join(root, "output"), filepath.Join(root, "cache")
	if err := os.MkdirAll(output, 0o755); err != nil {
		t.Fatal(err)
	}
	css := filepath.Join(output, "site.css")
	original := []byte(".used { color: blue; }")
	if err := os.WriteFile(css, original, 0o600); err != nil {
		t.Fatal(err)
	}
	selectors := csspurge.NewUsedSelectors()
	selectors.Classes["used"] = true
	config := models.NewCSSPurgeConfig()
	opts := buildPurgeOptions(config, false)
	processCSSFiles([]string{css}, output, cache, selectors, opts, config, false)
	if again := processCSSFiles([]string{css}, output, cache, selectors, opts, config, false); again.cacheHits != 1 {
		t.Fatal("expected cache hit on unchanged CSS")
	}
	content, err := os.ReadFile(css)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(content, original) {
		t.Fatalf("CSS changed without removed rules: %q", content)
	}
}
