package plugins

import (
	"os"
	"path/filepath"
	"testing"
)

func writeHTML(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}

func TestCSSPurgeHTMLCache_ReusesAndInvalidatesSelectorSet(t *testing.T) {
	root := t.TempDir()
	output := filepath.Join(root, "output")
	cache := filepath.Join(root, "cache")
	page := filepath.Join(output, "a", "index.html")
	writeHTML(t, page, `<div class="alpha" id="first"></div>`)

	first, hit, err := scanHTMLFilesForSelectors(output, cache, 2, false)
	if err != nil || hit {
		t.Fatalf("first scan: hit=%t err=%v", hit, err)
	}
	if !first.Classes["alpha"] || !first.IDs["first"] {
		t.Fatalf("first selector set missing values: %+v", first)
	}

	second, hit, err := scanHTMLFilesForSelectors(output, cache, 2, false)
	if err != nil || !hit {
		t.Fatalf("warm scan: hit=%t err=%v", hit, err)
	}
	if !second.Classes["alpha"] {
		t.Fatalf("cached selector set missing class alpha: %+v", second)
	}

	writeHTML(t, page, `<section class="beta"></section>`)
	third, hit, err := scanHTMLFilesForSelectors(output, cache, 2, false)
	if err != nil || hit {
		t.Fatalf("changed page scan: hit=%t err=%v", hit, err)
	}
	if !third.Classes["beta"] || third.Classes["alpha"] {
		t.Fatalf("changed selector set is stale: %+v", third)
	}

	writeHTML(t, filepath.Join(output, "b", "index.html"), `<p class="gamma"></p>`)
	fourth, hit, err := scanHTMLFilesForSelectors(output, cache, 2, false)
	if err != nil || hit || !fourth.Classes["gamma"] {
		t.Fatalf("added page invalidation failed: hit=%t selectors=%+v err=%v", hit, fourth, err)
	}
	if err := os.Remove(page); err != nil {
		t.Fatal(err)
	}
	fifth, hit, err := scanHTMLFilesForSelectors(output, cache, 2, false)
	if err != nil || hit || fifth.Classes["beta"] || !fifth.Classes["gamma"] {
		t.Fatalf("removed page invalidation failed: hit=%t selectors=%+v err=%v", hit, fifth, err)
	}
}

func TestCSSPurgeHTMLCache_CorruptCacheFallsBackToScan(t *testing.T) {
	root := t.TempDir()
	output := filepath.Join(root, "output")
	cache := filepath.Join(root, "cache")
	writeHTML(t, filepath.Join(output, "index.html"), `<main class="fresh"></main>`)
	cachePath := filepath.Join(cache, "html-selectors.json")
	if err := os.MkdirAll(filepath.Dir(cachePath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(cachePath, []byte("invalid json"), 0o600); err != nil {
		t.Fatal(err)
	}
	selectors, hit, err := scanHTMLFilesForSelectors(output, cache, 1, false)
	if err != nil || hit || !selectors.Classes["fresh"] {
		t.Fatalf("corrupt cache fallback failed: hit=%t selectors=%+v err=%v", hit, selectors, err)
	}
}
