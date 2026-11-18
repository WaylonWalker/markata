package plugins

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/WaylonWalker/markata/internal/fileidentity"
)

func TestFinalCSSLinkCacheReusesVerifiedFiles(t *testing.T) {
	output := t.TempDir()
	root, err := os.OpenRoot(output)
	if err != nil {
		t.Fatal(err)
	}
	defer root.Close()

	changes := []finalCSSReplacement{{
		matcher:     regexp.MustCompile(`^css/theme\.[0-9a-f]{8}\.css$`),
		destination: "css/theme.cafebabe.css",
	}}
	const current = `<link rel="stylesheet" href="/css/theme.cafebabe.css">`
	const stale = `<link rel="stylesheet" href="/css/theme.deadbeef.css">`
	if err := root.WriteFile("index.html", []byte(current), 0o600); err != nil {
		t.Fatal(err)
	}
	primed, err := rewriteOutputCSSLinks(root, output, changes, finalCSSLinkCache{}, true)
	if err != nil || primed.checked != 1 || primed.reused != 0 {
		t.Fatalf("initial verification = %+v, err = %v", primed, err)
	}

	// Put the deterministic observation time beyond the filesystem guard window.
	clock := time.Now().Add(time.Hour)
	rules := finalCSSLinkRules(changes)
	cachePath := filepath.Join(t.TempDir(), "verification.json")
	cache := finalCSSLinkCache{Rules: rules, Observed: clock, Files: primed.verified}
	saveFinalCSSLinkCache(cachePath, cache)
	previous := loadFinalCSSLinkCache(cachePath, rules, clock)
	if err := root.WriteFile("added.html", []byte(stale), 0o600); err != nil {
		t.Fatal(err)
	}
	result, err := rewriteOutputCSSLinks(root, output, changes, previous, true)
	if err != nil {
		t.Fatal(err)
	}
	if primed.verified["index.html"].Inode != 0 {
		if result.reused != 1 || result.checked != 1 || result.bytesRead != int64(len(stale)) {
			t.Fatalf("warm reconciliation = %+v", result)
		}
	} else if result.reused != 0 || result.checked != 2 {
		t.Fatalf("unsupported filesystem must recheck both files: %+v", result)
	}
	assertFinalCSSLink(t, root, "added.html", current)

	// A same-size replacement with the same mtime must miss the inode cache.
	info, err := root.Stat("index.html")
	if err != nil {
		t.Fatal(err)
	}
	if err := root.WriteFile("replacement", []byte(stale), 0o600); err != nil {
		t.Fatal(err)
	}
	replacement := filepath.Join(output, "replacement")
	if err := os.Chtimes(replacement, info.ModTime(), info.ModTime()); err != nil {
		t.Fatal(err)
	}
	if err := root.Rename("replacement", "index.html"); err != nil {
		t.Fatal(err)
	}
	result, err = rewriteOutputCSSLinks(root, output, changes, previous, true)
	if err != nil || result.reused != 0 || len(result.updatedHTMLHashes) != 1 {
		t.Fatalf("replacement reconciliation = %+v, err = %v", result, err)
	}
	assertFinalCSSLink(t, root, "index.html", current)

	// A different final digest must recheck all files even across stable inode hits.
	changes[0].destination = "css/theme.facefeed.css"
	previous = loadFinalCSSLinkCache(cachePath, finalCSSLinkRules(changes), clock)
	result, err = rewriteOutputCSSLinks(root, output, changes, previous, true)
	if err != nil || result.checked != 2 || result.reused != 0 || len(result.updatedHTMLHashes) != 2 {
		t.Fatalf("changed-rule reconciliation = %+v, err = %v", result, err)
	}
	assertFinalCSSLink(t, root, "index.html", strings.ReplaceAll(current, "cafebabe", "facefeed"))
}

func TestFinalCSSLinkCacheFallback(t *testing.T) {
	dir := t.TempDir()
	cachePath := filepath.Join(dir, "verification.json")
	now := time.Now()
	entry := finalCSSLinkCache{
		Rules:    "rules",
		Observed: now,
		Files:    map[string]fileidentity.Identity{"index.html": {Inode: 1}},
	}
	saveFinalCSSLinkCache(cachePath, entry)
	if got := loadFinalCSSLinkCache(cachePath, "rules", now.Add(-time.Hour)); len(got.Files) != 0 {
		t.Fatal("a backwards clock must discard the cache")
	}
	if err := os.WriteFile(cachePath, []byte("{"), 0o600); err != nil {
		t.Fatal(err)
	}
	if got := loadFinalCSSLinkCache(cachePath, "rules", now); len(got.Files) != 0 {
		t.Fatal("invalid JSON must discard the cache")
	}
	if path := finalCSSLinkCachePath(dir, dir); path != "" {
		t.Fatalf("verification cache must not enter output: %s", path)
	}
}

func assertFinalCSSLink(t *testing.T, root *os.Root, filename, want string) {
	t.Helper()
	data, err := root.ReadFile(filename)
	if err != nil || string(data) != want {
		t.Fatalf("%s = %q, want %q, err = %v", filename, data, want, err)
	}
}
