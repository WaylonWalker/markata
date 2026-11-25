package plugins

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/WaylonWalker/markata/internal/fileidentity"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
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

func TestFinalCSSLinkCacheAcrossCleanupColdWarmAndRuleChange(t *testing.T) {
	contentDir := t.TempDir()
	output := t.TempDir()
	if err := os.MkdirAll(filepath.Join(output, "css"), 0o755); err != nil {
		t.Fatal(err)
	}
	cssPath := filepath.Join(output, "css", "site.css")
	firstCSS := []byte("body { color: red }")
	secondCSS := []byte("body { color: blue }")
	firstHash := fmt.Sprintf("%x", sha256.Sum256(firstCSS))[:8]
	secondHash := fmt.Sprintf("%x", sha256.Sum256(secondCSS))[:8]
	if err := os.WriteFile(cssPath, firstCSS, 0o600); err != nil {
		t.Fatal(err)
	}
	const staleHash = "deadbeef"
	const body = "<main>Authored text about site.deadbeef.css</main>"
	firstHTML := "<link rel=\"stylesheet\" href=\"/css/site." + firstHash + ".css\">" + body
	original := filepath.Join(output, "index.html")
	if err := os.WriteFile(original, []byte(firstHTML), 0o600); err != nil {
		t.Fatal(err)
	}
	root, err := os.OpenRoot(output)
	if err != nil {
		t.Fatal(err)
	}
	defer root.Close()

	manager := lifecycle.NewManager()
	config := lifecycle.NewConfig()
	config.ContentDir = contentDir
	config.OutputDir = output
	manager.SetConfig(config)
	runCleanup := func() {
		t.Helper()
		// Force reconciliation as a build with a Configure-time CSS hash would.
		manager.SetAssetHash("css/site.css", staleHash)
		if err := reconcileFinalCSSAssets(manager, root, output); err != nil {
			t.Fatal(err)
		}
	}
	assertHTML := func(path, hash string) {
		t.Helper()
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(data), "/css/site."+hash+".css") {
			t.Fatalf("%s has wrong final CSS reference: %s", path, data)
		}
		if !strings.Contains(string(data), body) {
			t.Fatalf("%s modified authored body: %s", path, data)
		}
	}
	// Cold cleanup verifies unchanged HTML and persists its identity cache.
	runCleanup()
	assertHTML(original, firstHash)
	cachePath := finalCSSLinkCachePath(contentDir, output)
	if fileidentity.Available && cachePath == "" {
		t.Fatal("Linux identity cache unexpectedly disabled")
	}
	if cachePath != "" {
		changes, _, err := collectFinalCSSChanges(map[string]string{"css/site.css": staleHash}, root, output)
		if err != nil {
			t.Fatal(err)
		}
		verified := loadFinalCSSLinkCache(cachePath, finalCSSLinkRules(changes), time.Now())
		if _, ok := verified.Files["index.html"]; !ok {
			t.Fatal("cold cleanup did not persist verification for unchanged HTML")
		}
	}
	// A warm build must still inspect a freshly added cached page.
	added := filepath.Join(output, "added.html")
	addedHTML := "<link rel=\"stylesheet\" href=\"/css/site." + staleHash + ".css\">" + body
	if err := os.WriteFile(added, []byte(addedHTML), 0o600); err != nil {
		t.Fatal(err)
	}
	runCleanup()
	assertHTML(original, firstHash)
	assertHTML(added, firstHash)

	// A new final CSS digest invalidates any prior verified links, even when
	// neither HTML input has changed.
	if err := os.WriteFile(cssPath, secondCSS, 0o600); err != nil {
		t.Fatal(err)
	}
	runCleanup()
	assertHTML(original, secondHash)
	assertHTML(added, secondHash)
	if firstHash == secondHash {
		t.Fatal("test CSS fixtures must have different fingerprints")
	}
}

func assertFinalCSSLink(t *testing.T, root *os.Root, filename, want string) {
	t.Helper()
	data, err := root.ReadFile(filename)
	if err != nil || string(data) != want {
		t.Fatalf("%s = %q, want %q, err = %v", filename, data, want, err)
	}
}

func TestFinalCSSLinkCacheEmptyVerificationReplacesOldCache(t *testing.T) {
	cachePath := filepath.Join(t.TempDir(), "verification.json")
	rules := "unchanged-rules"
	now := time.Now()
	old := finalCSSLinkCache{
		Rules:    rules,
		Observed: now.Add(-time.Minute),
		Files: map[string]fileidentity.Identity{
			"stale.html": {Device: 1, Inode: 1, Size: 42},
		},
	}
	saveFinalCSSLinkCache(cachePath, old)
	if got := loadFinalCSSLinkCache(cachePath, rules, now); len(got.Files) != 1 {
		t.Fatalf("initial cache should contain the previous verified file: %+v", got)
	}

	// A subsequent successful cleanup may have no eligible HTML identities.
	// It must replace the old snapshot instead of keeping stale verification.
	empty := finalCSSLinkCache{
		Rules:    rules,
		Observed: now,
		Files:    map[string]fileidentity.Identity{},
	}
	saveFinalCSSLinkCache(cachePath, empty)
	got := loadFinalCSSLinkCache(cachePath, rules, now.Add(time.Second))
	if len(got.Files) != 0 || !got.Observed.Equal(now) {
		t.Fatalf("empty verification should replace the old cache: %+v", got)
	}
}

func TestFinalCSSLinkCacheOutputPromotionInvalidates(t *testing.T) {
	if !fileidentity.Available {
		t.Skip("filesystem identity caching is Linux-only")
	}
	base := t.TempDir()
	contentDir := filepath.Join(base, "content")
	output := filepath.Join(base, "published")
	replacement := filepath.Join(base, "replacement")
	if err := os.MkdirAll(contentDir, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(output, 0o700); err != nil {
		t.Fatal(err)
	}
	changes := []finalCSSReplacement{{
		matcher:     regexp.MustCompile(`^css/theme\.[0-9a-f]{8}\.css$`),
		destination: "css/theme.cafebabe.css",
	}}
	const current = `<link rel="stylesheet" href="/css/theme.cafebabe.css">`
	const stale = `<link rel="stylesheet" href="/css/theme.deadbeef.css">`
	root, err := os.OpenRoot(output)
	if err != nil {
		t.Fatal(err)
	}
	defer root.Close()
	if err := root.WriteFile("index.html", []byte(current), 0o600); err != nil {
		t.Fatal(err)
	}
	oldInfo, err := root.Stat("index.html")
	if err != nil {
		t.Fatal(err)
	}
	first, err := rewriteOutputCSSLinks(root, output, changes, finalCSSLinkCache{}, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(first.verified) != 1 {
		t.Skip("the filesystem cannot provide a reusable file identity")
	}

	// Promote a new output tree at the same path. A size- and mtime-matched
	// HTML file must still miss the saved cache because its inode changed.
	if err := os.MkdirAll(replacement, 0o700); err != nil {
		t.Fatal(err)
	}
	replacedFile := filepath.Join(replacement, "index.html")
	if err := os.WriteFile(replacedFile, []byte(stale), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Chtimes(replacedFile, oldInfo.ModTime(), oldInfo.ModTime()); err != nil {
		t.Fatal(err)
	}
	rules := finalCSSLinkRules(changes)
	clock := time.Now().Add(time.Hour)
	cachePath := finalCSSLinkCachePath(contentDir, output)
	if cachePath == "" {
		t.Fatal("expected a cache path outside published output")
	}
	saveFinalCSSLinkCache(cachePath, finalCSSLinkCache{Rules: rules, Observed: clock, Files: first.verified})
	if err := os.Rename(output, filepath.Join(base, "previous-release")); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(replacement, output); err != nil {
		t.Fatal(err)
	}
	if nextPath := finalCSSLinkCachePath(contentDir, output); nextPath != cachePath {
		t.Fatalf("cache path changed after output promotion: %q != %q", nextPath, cachePath)
	}
	newRoot, err := os.OpenRoot(output)
	if err != nil {
		t.Fatal(err)
	}
	defer newRoot.Close()
	previous := loadFinalCSSLinkCache(cachePath, rules, clock)
	second, err := rewriteOutputCSSLinks(newRoot, output, changes, previous, true)
	if err != nil {
		t.Fatal(err)
	}
	if second.reused != 0 || second.checked != 1 || len(second.updatedHTMLHashes) != 1 {
		t.Fatalf("promoted HTML must be rechecked and rewritten: %+v", second)
	}
	assertFinalCSSLink(t, newRoot, "index.html", current)
}
