package plugins

import (
	"crypto/sha256"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/WaylonWalker/markata/internal/fileidentity"
	"github.com/WaylonWalker/markata/pkg/csspurge"
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

func TestCSSPurgeHTMLCache_StableIdentityHitAvoidsHTMLReads(t *testing.T) {
	if !fileidentity.Available {
		t.Skip("verified file identity is available on Linux only")
	}
	root := t.TempDir()
	output := filepath.Join(root, "output")
	cache := filepath.Join(root, "cache")
	page := filepath.Join(output, "index.html")
	initial := `<div class="alpha" id="first"></div>`
	writeHTML(t, page, initial)

	// The identity cache rejects files written within its one-second timestamp
	// safety window. Age the output before priming the cache to test the warm hit.
	time.Sleep(1100 * time.Millisecond)
	first, firstStats, err := scanHTMLFilesForSelectorsWithStats(output, cache, 2, false)
	if err != nil || firstStats.BytesRead != int64(len(initial)) || firstStats.CacheMode != htmlSelectorCacheModeScan {
		t.Fatalf("initial scan: stats=%+v err=%v", firstStats, err)
	}
	warm, warmStats, err := scanHTMLFilesForSelectorsWithStats(output, cache, 2, false)
	if err != nil || !warmStats.CacheHit || warmStats.CacheMode != htmlSelectorCacheModeIdentity || warmStats.BytesRead != 0 || warmStats.BytesReused != int64(len(initial)) {
		t.Fatalf("identity hit: stats=%+v err=%v", warmStats, err)
	}
	if !reflect.DeepEqual(first, warm) {
		t.Fatalf("identity cache changed selectors: first=%+v warm=%+v", first, warm)
	}

	// An in-place same-size edit with a restored mtime still changes ctime and
	// must invalidate the identity fast path.
	info, err := os.Stat(page)
	if err != nil {
		t.Fatal(err)
	}
	replacement := `<div class="omega" id="first"></div>`
	if len(replacement) != len(initial) {
		t.Fatalf("replacement size = %d, want %d", len(replacement), len(initial))
	}
	writeHTML(t, page, replacement)
	if err := os.Chtimes(page, info.ModTime(), info.ModTime()); err != nil {
		t.Fatal(err)
	}
	changed, changedStats, err := scanHTMLFilesForSelectorsWithStats(output, cache, 2, false)
	if err != nil || changedStats.CacheHit || !changed.Classes["omega"] || changed.Classes["alpha"] {
		t.Fatalf("same-size edit invalidation: stats=%+v selectors=%+v err=%v", changedStats, changed, err)
	}

	addedPath := filepath.Join(output, "added.html")
	writeHTML(t, addedPath, `<p class="added"></p>`)
	added, addedStats, err := scanHTMLFilesForSelectorsWithStats(output, cache, 2, false)
	if err != nil || addedStats.CacheHit || !added.Classes["added"] {
		t.Fatalf("added page invalidation: stats=%+v selectors=%+v err=%v", addedStats, added, err)
	}
	if err := os.Remove(addedPath); err != nil {
		t.Fatal(err)
	}
	removed, removedStats, err := scanHTMLFilesForSelectorsWithStats(output, cache, 2, false)
	if err != nil || removedStats.CacheHit || removed.Classes["added"] || !removed.Classes["omega"] {
		t.Fatalf("removed page invalidation: stats=%+v selectors=%+v err=%v", removedStats, removed, err)
	}
}

func TestCSSPurgeHTMLCache_SymlinkFallsBackToExactContentHash(t *testing.T) {
	if !fileidentity.Available {
		t.Skip("verified file identity is available on Linux only")
	}
	root := t.TempDir()
	output := filepath.Join(root, "output")
	cache := filepath.Join(root, "cache")
	target := filepath.Join(output, "page-source")
	link := filepath.Join(output, "index.html")
	writeHTML(t, target, `<main class="linked"></main>`)
	if err := os.Symlink(target, link); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}

	first, firstStats, err := scanHTMLFilesForSelectorsWithStats(output, cache, 1, false)
	if err != nil || !first.Classes["linked"] || firstStats.CacheMode != htmlSelectorCacheModeScan {
		t.Fatalf("symlink initial scan: stats=%+v selectors=%+v err=%v", firstStats, first, err)
	}
	warm, warmStats, err := scanHTMLFilesForSelectorsWithStats(output, cache, 1, false)
	if err != nil || !warm.Classes["linked"] || !warmStats.CacheHit || warmStats.CacheMode != htmlSelectorCacheModeContent || warmStats.BytesRead == 0 {
		t.Fatalf("symlink fallback: stats=%+v selectors=%+v err=%v", warmStats, warm, err)
	}
}

func TestCSSPurgeHTMLCache_RejectsFutureIdentityObservation(t *testing.T) {
	if !fileidentity.Available {
		t.Skip("verified file identity is available on Linux only")
	}
	root := t.TempDir()
	page := filepath.Join(root, "index.html")
	writeHTML(t, page, `<main class="present"></main>`)
	info, err := os.Lstat(page)
	if err != nil {
		t.Fatal(err)
	}
	identity := fileidentity.Read(info)
	cached := cssPurgeHTMLCache{
		Schema:      cssPurgeHTMLCacheSchema,
		Fingerprint: "known-content-hash",
		Observed:    time.Now().Add(time.Hour),
		HTMLBytes:   info.Size(),
		Files:       map[string]fileidentity.Identity{"index.html": identity},
		Selectors:   *newSelectorSetForTest("present"),
	}
	if cssPurgeHTMLIdentityCacheMatches(cached, map[string]fileidentity.Identity{"index.html": identity}, time.Now()) {
		t.Fatal("future cache observation was trusted")
	}
}

func newSelectorSetForTest(class string) *csspurge.UsedSelectors {
	selectors := csspurge.NewUsedSelectors()
	selectors.Classes[class] = true
	return selectors
}

func TestCSSPurgeHTMLFingerprint_StreamsExactBytes(t *testing.T) {
	root := t.TempDir()
	page := filepath.Join(root, "index.html")
	content := `<main class="streamed"><p>Exact bytes.</p></main>`
	writeHTML(t, page, content)

	result := fingerprintHTMLFile(page)
	if result.err != nil {
		t.Fatal(result.err)
	}
	want := sha256.Sum256([]byte(content))
	if result.sum != want || result.size != int64(len(content)) {
		t.Fatalf("stream fingerprint = %x/%d, want %x/%d", result.sum, result.size, want, len(content))
	}
}
