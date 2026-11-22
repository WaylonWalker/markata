package plugins

import (
	"crypto/sha256"
	"encoding/json"
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

func TestCSSPurgeHTMLCache_ParsesOnlyChangedFiles(t *testing.T) {
	root := t.TempDir()
	output := filepath.Join(root, "output")
	cache := filepath.Join(root, "cache")
	changedPage := filepath.Join(output, "changed.html")
	stablePage := filepath.Join(output, "stable.html")
	changedInitial := `<main class="alpha"></main>`
	changedNext := `<main class="bravo"></main>`
	stable := `<aside class="stable" id="keep"></aside>`
	if len(changedInitial) != len(changedNext) {
		t.Fatalf("test contents must have the same length: %d != %d", len(changedInitial), len(changedNext))
	}
	writeHTML(t, changedPage, changedInitial)
	writeHTML(t, stablePage, stable)
	if fileidentity.Available {
		// Keep fixture identities outside the coarse-filesystem timestamp guard.
		time.Sleep(1100 * time.Millisecond)
	}

	first, firstStats, err := scanHTMLFilesForSelectorsWithStats(output, cache, 2, false)
	if err != nil || firstStats.FilesParsed != 2 || firstStats.FilesReused != 0 || firstStats.CacheMode != htmlSelectorCacheModeScan {
		t.Fatalf("initial scan: stats=%+v err=%v", firstStats, err)
	}
	if !first.Classes["alpha"] || !first.Classes["stable"] || !first.IDs["keep"] {
		t.Fatalf("initial selector set missing values: %+v", first)
	}

	// A same-size replacement must use the exact content hash to identify only
	// the changed page, then replace rather than union its old selectors.
	changedInfo, err := os.Stat(changedPage)
	if err != nil {
		t.Fatal(err)
	}
	writeHTML(t, changedPage, changedNext)
	if err := os.Chtimes(changedPage, changedInfo.ModTime(), changedInfo.ModTime()); err != nil {
		t.Fatal(err)
	}
	if fileidentity.Available {
		time.Sleep(1100 * time.Millisecond)
	}
	second, secondStats, err := scanHTMLFilesForSelectorsWithStats(output, cache, 2, false)
	wantPartialMode := htmlSelectorCacheModePartial
	if fileidentity.Available {
		wantPartialMode = htmlSelectorCacheModeIdentityPartial
	}
	if err != nil || secondStats.CacheHit || secondStats.CacheMode != wantPartialMode ||
		secondStats.FilesParsed != 1 || secondStats.FilesReused != 1 {
		t.Fatalf("partial scan: stats=%+v err=%v", secondStats, err)
	}
	if fileidentity.Available && secondStats.BytesHashed != int64(len(changedNext)) {
		t.Fatalf("partial identity scan hashed %d bytes, want only changed page bytes %d", secondStats.BytesHashed, len(changedNext))
	}
	if !second.Classes["bravo"] || !second.Classes["stable"] || second.Classes["alpha"] || !second.IDs["keep"] {
		t.Fatalf("partial selector set is stale or incomplete: %+v", second)
	}

	if err := os.Remove(changedPage); err != nil {
		t.Fatal(err)
	}
	third, thirdStats, err := scanHTMLFilesForSelectorsWithStats(output, cache, 2, false)
	if err != nil || thirdStats.CacheHit || thirdStats.CacheMode != wantPartialMode ||
		thirdStats.FilesParsed != 0 || thirdStats.FilesReused != 1 {
		t.Fatalf("removed-page update: stats=%+v err=%v", thirdStats, err)
	}
	if third.Classes["alpha"] || third.Classes["bravo"] || !third.Classes["stable"] || !third.IDs["keep"] {
		t.Fatalf("removed page selectors remain in union: %+v", third)
	}

	clean, _, err := scanHTMLFilesForSelectors(output, filepath.Join(root, "clean-cache"), 2, false)
	if err != nil {
		t.Fatalf("clean full scan: %v", err)
	}
	if !reflect.DeepEqual(third, clean) {
		t.Fatalf("incremental selectors differ from clean scan: incremental=%+v clean=%+v", third, clean)
	}

	// Replacing a page with identical bytes changes its verified identity. The
	// exact hash should preserve its cached selectors without parsing the page.
	writeHTML(t, stablePage, stable)
	sameBytes, sameBytesStats, err := scanHTMLFilesForSelectorsWithStats(output, cache, 2, false)
	if err != nil || !sameBytesStats.CacheHit || sameBytesStats.CacheMode != wantPartialMode ||
		sameBytesStats.FilesParsed != 0 || sameBytesStats.FilesReused != 1 || !reflect.DeepEqual(third, sameBytes) {
		t.Fatalf("identical-content replacement: stats=%+v selectors=%+v err=%v", sameBytesStats, sameBytes, err)
	}
	if fileidentity.Available && sameBytesStats.BytesHashed != int64(len(stable)) {
		t.Fatalf("identical-content replacement hashed %d bytes, want only replaced page bytes %d", sameBytesStats.BytesHashed, len(stable))
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

func TestCSSPurgeHTMLCache_RejectsValidJSONWithAlteredAggregateSelectors(t *testing.T) {
	root := t.TempDir()
	output := filepath.Join(root, "output")
	cache := filepath.Join(root, "cache")
	writeHTML(t, filepath.Join(output, "index.html"), `<main class="protected"></main>`)

	initial, _, err := scanHTMLFilesForSelectors(output, cache, 1, false)
	if err != nil || !initial.Classes["protected"] {
		t.Fatalf("initial scan: selectors=%+v err=%v", initial, err)
	}
	cachePath := filepath.Join(cache, "html-selectors.json")
	encoded, err := os.ReadFile(cachePath)
	if err != nil {
		t.Fatal(err)
	}
	var corrupted cssPurgeHTMLCache
	if err := json.Unmarshal(encoded, &corrupted); err != nil {
		t.Fatal(err)
	}
	corrupted.Selectors.Classes["protected"] = false
	encoded, err = json.Marshal(corrupted)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(cachePath, encoded, 0o600); err != nil {
		t.Fatal(err)
	}

	recovered, stats, err := scanHTMLFilesForSelectorsWithStats(output, cache, 1, false)
	if err != nil || stats.CacheMode == htmlSelectorCacheModeIdentity || !recovered.Classes["protected"] {
		t.Fatalf("altered aggregate cache was trusted: stats=%+v selectors=%+v err=%v", stats, recovered, err)
	}
}

func TestCSSPurgeHTMLScan_ParsesBytesThatWereHashed(t *testing.T) {
	root := t.TempDir()
	output := filepath.Join(root, "output")
	page := filepath.Join(output, "index.html")
	writeHTML(t, page, `<main class="current"></main>`)
	expected := map[string]string{"index.html": "a stale fingerprint"}

	parsed, complete, hashesMatch := scanHTMLFilesWithExpectedHashesConcurrently([]string{page}, output, expected, 1)
	if !complete || hashesMatch || !parsed[page].Classes["current"] {
		t.Fatalf("changed bytes were not parsed consistently: complete=%t hashes_match=%t selectors=%+v", complete, hashesMatch, parsed[page])
	}
	if err := validateCSSPurgeHTMLScan(complete, hashesMatch); err == nil {
		t.Fatal("changed HTML snapshot was allowed to proceed to CSS purging")
	}

	content, err := os.ReadFile(page)
	if err != nil {
		t.Fatal(err)
	}
	if got := cssPurgeBytesHash(content); got == expected["index.html"] {
		t.Fatal("test expected hash unexpectedly matches current page")
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
	if err != nil || firstStats.BytesHashed != int64(len(initial)) || firstStats.CacheMode != htmlSelectorCacheModeScan {
		t.Fatalf("initial scan: stats=%+v err=%v", firstStats, err)
	}
	if err := os.WriteFile(filepath.Join(cache, htmlSelectorFileCacheName), []byte("corrupt page cache"), 0o600); err != nil {
		t.Fatal(err)
	}
	warm, warmStats, err := scanHTMLFilesForSelectorsWithStats(output, cache, 2, false)
	if err != nil || !warmStats.CacheHit || warmStats.CacheMode != htmlSelectorCacheModeIdentity || warmStats.BytesHashed != 0 {
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
	if err != nil || !warm.Classes["linked"] || !warmStats.CacheHit || warmStats.CacheMode != htmlSelectorCacheModeContent || warmStats.BytesHashed == 0 {
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
