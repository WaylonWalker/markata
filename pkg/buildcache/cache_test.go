package buildcache

import (
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"
)

func TestNew(t *testing.T) {
	cache := New("")
	if cache == nil {
		t.Fatal("New returned nil")
	}
	if cache.Version != CacheVersion {
		t.Errorf("Version = %d, want %d", cache.Version, CacheVersion)
	}
	if cache.Posts == nil {
		t.Error("Posts map is nil")
	}
	if cache.Graph == nil {
		t.Error("Graph is nil")
	}
	if cache.changedSlugs == nil {
		t.Error("changedSlugs map is nil")
	}
}

func TestCache_SetAndGetDependencies(t *testing.T) {
	cache := New("")

	// Set dependencies
	cache.SetDependencies("pages/post-a.md", "post-a", []string{"post-b", "post-c"})

	// Verify through graph
	deps := cache.Graph.GetDependencies("pages/post-a.md")
	if len(deps) != 2 {
		t.Errorf("GetDependencies returned %d deps, want 2", len(deps))
	}
}

func TestCache_GetAffectedPosts(t *testing.T) {
	cache := New("")

	// Set up dependencies: post-a -> post-b -> post-c
	cache.SetDependencies("pages/post-a.md", "post-a", []string{"post-b"})
	cache.SetDependencies("pages/post-b.md", "post-b", []string{"post-c"})

	// When post-c changes, both post-a and post-b should be affected
	affected := cache.GetAffectedPosts([]string{"post-c"})
	if len(affected) != 2 {
		t.Errorf("GetAffectedPosts returned %d posts, want 2: %v", len(affected), affected)
	}
}

func TestCache_ShouldRebuildWithSlug_DependencyChanged(t *testing.T) {
	cache := New("")

	// Set up: post-a depends on post-b
	cache.SetDependencies("pages/post-a.md", "post-a", []string{"post-b"})

	// Mark post-a as previously built
	cache.MarkRebuilt("pages/post-a.md", "hash123", "output/post-a/index.html", "post.html")

	// Mark that post-b changed this build
	cache.MarkSlugChanged("post-b")

	// Even though post-a's hash matches, it should rebuild because post-b changed
	shouldRebuild := cache.ShouldRebuildWithSlug("pages/post-a.md", "post-a", "hash123", "post.html")
	if !shouldRebuild {
		t.Error("ShouldRebuildWithSlug = false, want true (dependency changed)")
	}
}

func TestCache_ShouldRebuildWithSlug_NoDependencyChange(t *testing.T) {
	cache := New("")

	// Set up: post-a depends on post-b
	cache.SetDependencies("pages/post-a.md", "post-a", []string{"post-b"})

	// Mark post-a as previously built
	cache.MarkRebuilt("pages/post-a.md", "hash123", "output/post-a/index.html", "post.html")

	// Don't mark post-b as changed

	// post-a should NOT rebuild (hash matches, no dependency changed)
	shouldRebuild := cache.ShouldRebuildWithSlug("pages/post-a.md", "post-a", "hash123", "post.html")
	if shouldRebuild {
		t.Error("ShouldRebuildWithSlug = true, want false (no dependency changed)")
	}
}

func TestCache_SetDependencies_UnchangedDoesNotMarkDirty(t *testing.T) {
	cache := New("")
	cache.SetDependencies("pages/post-a.md", "post-a", []string{"post-b"})
	cache.dirty = false

	cache.SetDependencies("pages/post-a.md", "post-a", []string{"post-b", "post-b"})

	if cache.dirty {
		t.Fatal("dirty = true, want false for unchanged dependencies")
	}
}

func TestCache_GetChangedSlugs(t *testing.T) {
	cache := New("")

	// No changes initially
	changed := cache.GetChangedSlugs()
	if len(changed) != 0 {
		t.Errorf("GetChangedSlugs = %v, want []", changed)
	}

	// Mark some slugs as changed
	cache.MarkSlugChanged("post-a")
	cache.MarkSlugChanged("post-b")

	changed = cache.GetChangedSlugs()
	if len(changed) != 2 {
		t.Errorf("GetChangedSlugs returned %d slugs, want 2", len(changed))
	}
}

func TestCache_ResetStats_ClearsChangedSlugs(t *testing.T) {
	cache := New("")

	cache.MarkSlugChanged("post-a")
	cache.MarkSlugChanged("post-b")

	cache.ResetStats()

	changed := cache.GetChangedSlugs()
	if len(changed) != 0 {
		t.Errorf("GetChangedSlugs after ResetStats = %v, want []", changed)
	}
}

func TestCache_MarkRebuiltWithSlug_TracksChange(t *testing.T) {
	cache := New("")

	cache.MarkRebuiltWithSlug("pages/post-a.md", "post-a", "hash123", "output/post-a/index.html", "post.html")

	changed := cache.GetChangedSlugs()
	if len(changed) != 1 || changed[0] != "post-a" {
		t.Errorf("GetChangedSlugs = %v, want [post-a]", changed)
	}
	if got := cache.GetCachedPost("pages/post-a.md").Slug; got != "post-a" {
		t.Fatalf("cached slug = %q, want post-a", got)
	}
}

func TestCache_MarkRebuiltWithSlug_PreservesCachedHTMLPaths(t *testing.T) {
	cache := New("")
	cache.Posts["pages/post-a.md"] = &PostCache{
		ArticleHTMLPath:    "/tmp/article.html",
		FullHTMLPath:       "/tmp/full.html",
		FeedMembershipHash: "membership-hash",
	}

	cache.MarkRebuiltWithSlug("pages/post-a.md", "post-a", "hash123", "output/post-a/index.html", "post.html")

	got := cache.Posts["pages/post-a.md"]
	if got.ArticleHTMLPath != "/tmp/article.html" {
		t.Fatalf("ArticleHTMLPath = %q, want preserved path", got.ArticleHTMLPath)
	}
	if got.FullHTMLPath != "/tmp/full.html" {
		t.Fatalf("FullHTMLPath = %q, want preserved path", got.FullHTMLPath)
	}
	if got.FeedMembershipHash != "membership-hash" {
		t.Fatalf("FeedMembershipHash = %q, want preserved hash", got.FeedMembershipHash)
	}
}

func TestCache_MarkInputProcessed_TracksInputsWithoutCountingWrite(t *testing.T) {
	cache := New(t.TempDir())
	cache.MarkRebuilt("pages/post-a.md", "old-hash", "output/post-a/index.html", "post.html")
	cache.ResetStats()
	cache.dirty = false

	if !cache.MarkInputProcessed("pages/post-a.md", "post-a", "new-hash", "custom.html") {
		t.Fatal("changed processed input was not recorded")
	}

	if cache.ShouldRebuild("pages/post-a.md", "new-hash", "custom.html") {
		t.Fatal("processed input was still marked for rebuild")
	}
	if skipped, rebuilt := cache.Stats(); skipped != 0 || rebuilt != 0 {
		t.Fatalf("stats = (%d skipped, %d rebuilt), want (0, 0)", skipped, rebuilt)
	}
	if got := cache.Posts["pages/post-a.md"].OutputPath; got != "" {
		t.Fatalf("output path = %q, want cleared HTML ownership", got)
	}

	cache.dirty = false
	if cache.MarkInputProcessed("pages/post-a.md", "post-a", "new-hash", "custom.html") {
		t.Fatal("unchanged processed input was reported as changed")
	}
	if cache.dirty {
		t.Fatal("unchanged processed input dirtied the cache")
	}
}

func TestCache_MarkSlugChanged_IgnoresEmptySlug(t *testing.T) {
	cache := New(t.TempDir())
	cache.MarkSlugChanged("")
	cache.MarkFeedSlugChanged("")

	if changed := cache.GetChangedSlugs(); len(changed) != 0 {
		t.Fatalf("changed slugs = %v, want empty", changed)
	}
	if changed := cache.GetChangedFeedSlugs(); len(changed) != 0 {
		t.Fatalf("changed feed slugs = %v, want empty", changed)
	}
}

func TestCache_MarkChangedPaths(t *testing.T) {
	cache := New("")

	cache.SetDependencies("pages/post-a.md", "post-a", []string{"post-b"})
	cache.MarkChangedPaths([]string{"pages/post-a.md"})

	changed := cache.GetChangedSlugs()
	if len(changed) != 1 || changed[0] != "post-a" {
		t.Errorf("GetChangedSlugs = %v, want [post-a]", changed)
	}
}

func TestCache_MarkMissingPaths_PreservesOwnershipUntilCleanup(t *testing.T) {
	cache := New("")
	cache.UpdateModTime("pages/keep.md", 1, "keep")
	cache.UpdateModTime("pages/deleted.md", 1, "deleted")

	got := cache.MarkMissingPaths([]string{"pages/keep.md"})
	if len(got) != 1 || got[0] != "pages/deleted.md" {
		t.Fatalf("MarkMissingPaths() = %v, want [pages/deleted.md]", got)
	}
	if cache.GetCachedPost("pages/deleted.md") == nil {
		t.Fatal("MarkMissingPaths removed the stale entry before cleanup")
	}
	changed := cache.GetChangedSlugs()
	if len(changed) != 1 || changed[0] != "deleted" {
		t.Fatalf("GetChangedSlugs() = %v, want [deleted]", changed)
	}
}

func TestCache_GlobalInvalidation_PreservesMissingPostOwnership(t *testing.T) {
	cache := New("")
	cache.AssetsHash = "old-assets"
	cache.Posts["pages/deleted.md"] = &PostCache{
		OutputPath: "output/deleted/index.html",
		Slug:       "deleted",
	}
	cache.SetPostSlug("pages/deleted.md", "deleted")

	if !cache.SetAssetsHash("new-assets") {
		t.Fatal("SetAssetsHash() did not report an invalidation")
	}
	if len(cache.Posts) != 0 {
		t.Fatalf("Posts after invalidation = %d, want 0", len(cache.Posts))
	}
	if cache.GetCachedPost("pages/deleted.md") == nil {
		t.Fatal("global invalidation discarded deleted post ownership")
	}

	missing := cache.MarkMissingPaths(nil)
	if len(missing) != 1 || missing[0] != "pages/deleted.md" {
		t.Fatalf("MarkMissingPaths() = %v, want [pages/deleted.md]", missing)
	}
	if removed := cache.RemoveStale(map[string]bool{}); removed != 1 {
		t.Fatalf("RemoveStale() = %d, want 1", removed)
	}
	if cache.GetCachedPost("pages/deleted.md") != nil {
		t.Fatal("stale post ownership survived cleanup")
	}
}

func TestCache_SetConfigHash_InvalidatesFeedsListing(t *testing.T) {
	cache := New("")
	cache.ConfigHash = "old-config"
	cache.SetFeedsListingHash("old-listing")
	cache.SetImageLibraryHash("old-images")

	if !cache.SetConfigHash("new-config") {
		t.Fatal("SetConfigHash() did not report an invalidation")
	}
	if got := cache.GetFeedsListingHash(); got != "" {
		t.Fatalf("FeedsListingHash after config invalidation = %q, want empty", got)
	}
	if got := cache.GetImageLibraryHash(); got != "" {
		t.Fatalf("ImageLibraryHash after config invalidation = %q, want empty", got)
	}
}

func TestCache_ImageLibraryHash_PersistsAndTemplateInvalidates(t *testing.T) {
	dir := t.TempDir()
	cache := New(dir)
	cache.SetImageLibraryHash("image-hash")
	pagePath := filepath.Join(dir, "output", "images", "index.html")
	jsonPath := filepath.Join(dir, "output", "images", "index.json")
	cache.SetImageLibraryOutputs(dir, dir, map[string]string{
		pagePath: "page-hash",
		jsonPath: "json-hash",
	})
	if got := cache.GetImageLibraryHash(); got != "image-hash" {
		t.Fatalf("GetImageLibraryHash() = %q, want image-hash", got)
	}
	if err := cache.Save(); err != nil {
		t.Fatalf("Save() error = %v", err)
	}
	loaded, err := Load(dir)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if got := loaded.GetImageLibraryHash(); got != "image-hash" {
		t.Fatalf("loaded ImageLibraryHash = %q, want image-hash", got)
	}
	if !loaded.OwnsImageLibraryOutput(dir, dir, pagePath) || !loaded.OwnsImageLibraryOutput(dir, dir, jsonPath) {
		contentRoot, outputRoot, outputs := loaded.GetImageLibraryOutputs()
		t.Fatalf("loaded image-library output ownership = (%q, %q, %v), want both paths", contentRoot, outputRoot, outputs)
	}

	loaded.TemplatesHash = "old-templates"
	if !loaded.SetTemplatesHash("new-templates") {
		t.Fatal("SetTemplatesHash() did not report an invalidation")
	}
	if got := loaded.GetImageLibraryHash(); got != "" {
		t.Fatalf("ImageLibraryHash after template invalidation = %q, want empty", got)
	}
	if !loaded.OwnsImageLibraryOutput(dir, dir, pagePath) {
		t.Fatal("template invalidation discarded image-library output ownership")
	}
}

func TestCache_NavPreviewHash_InvalidatesSharedPages(t *testing.T) {
	cache := New(t.TempDir())
	cache.Posts["post.md"] = &PostCache{}
	cache.Feeds["journal"] = &FeedCache{}
	if !cache.SetNavPreviewHash("first") {
		t.Fatal("initial preview hash should invalidate existing pages")
	}
	if len(cache.Posts) != 0 || len(cache.Feeds) != 0 {
		t.Fatal("preview change left stale pages in cache")
	}
	cache.Posts["post.md"] = &PostCache{}
	if cache.SetNavPreviewHash("first") || len(cache.Posts) != 1 {
		t.Fatal("unchanged preview hash invalidated a page")
	}
}

func TestCache_FontpackHashInvalidatesRenderedPagesOnlyWhenChanged(t *testing.T) {
	cache := New(t.TempDir())
	cache.Posts["post.md"] = &PostCache{}
	cache.Feeds["journal"] = &FeedCache{}
	if !cache.SetFontpackHash("first") {
		t.Fatal("initial fontpack hash should invalidate rendered pages")
	}
	if len(cache.Posts) != 0 || len(cache.Feeds) != 0 {
		t.Fatal("fontpack change left stale rendered pages in cache")
	}
	cache.Posts["post.md"] = &PostCache{}
	if cache.SetFontpackHash("first") || len(cache.Posts) != 1 {
		t.Fatal("unchanged fontpack hash invalidated a warm page")
	}
	if !cache.SetFontpackHash("second") || len(cache.Posts) != 0 {
		t.Fatal("changed fontpack hash did not invalidate rendered pages")
	}
}

func TestCache_ImageLibraryMedia_PersistsFingerprints(t *testing.T) {
	dir := t.TempDir()
	cache := New(dir)
	path := filepath.Join(dir, "static", "photo.webp")
	want := map[string]ImageLibraryMediaFingerprint{
		path: {Size: 42, ModTime: 123, ChangeTime: 456, ChangeTimeKnown: true, ContentHash: "content-hash"},
	}
	if !cache.SetImageLibraryMedia(want) {
		t.Fatal("SetImageLibraryMedia() did not report a new fingerprint")
	}
	if cache.SetImageLibraryMedia(want) {
		t.Fatal("SetImageLibraryMedia() reported a change for identical fingerprints")
	}
	if err := cache.Save(); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	loaded, err := Load(dir)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	got := loaded.GetImageLibraryMedia()
	if len(got) != 1 || got[filepath.ToSlash(filepath.Clean(path))] != want[path] {
		t.Fatalf("loaded image-library media = %#v, want %#v", got, want)
	}
}

func TestCache_SetAssetsHash_InvalidatesImageLibrary(t *testing.T) {
	cache := New("")
	cache.AssetsHash = "old-assets"
	cache.SetImageLibraryHash("old-images")

	if !cache.SetAssetsHash("new-assets") {
		t.Fatal("SetAssetsHash() did not report an invalidation")
	}
	if got := cache.GetImageLibraryHash(); got != "" {
		t.Fatalf("ImageLibraryHash after asset invalidation = %q, want empty", got)
	}
}

func TestCache_ImageLibraryOutputs_AreScopedToOutputRoot(t *testing.T) {
	root := t.TempDir()
	otherRoot := t.TempDir()
	outside := filepath.Join(root, "..", "outside.html")
	cache := New("")
	cache.SetImageLibraryOutputs(root, root, map[string]string{
		filepath.Join(root, "images", "index.html"): "page-hash",
		outside: "outside-hash",
	})

	if !cache.OwnsImageLibraryOutput(root, root, filepath.Join(root, "images", "index.html")) {
		t.Fatal("recorded output was not owned")
	}
	if cache.OwnsImageLibraryOutput(otherRoot, otherRoot, filepath.Join(otherRoot, "images", "index.html")) {
		t.Fatal("output ownership crossed output roots")
	}
	_, _, outputs := cache.GetImageLibraryOutputs()
	if len(outputs) != 1 {
		t.Fatalf("outputs = %v, want only paths inside root", outputs)
	}
}

func TestCache_MarkAffectedDependents(t *testing.T) {
	cache := New("")

	cache.SetDependencies("pages/post-a.md", "post-a", []string{"post-b"})
	cache.SetDependencies("pages/post-b.md", "post-b", []string{"post-c"})

	cache.MarkAffectedDependents([]string{"post-c"})

	changed := cache.GetChangedSlugs()
	if len(changed) != 2 {
		t.Fatalf("GetChangedSlugs returned %d slugs, want 2", len(changed))
	}
	want := map[string]bool{"post-a": true, "post-b": true}
	for _, slug := range changed {
		if !want[slug] {
			t.Errorf("unexpected changed slug %q", slug)
		}
	}
}

func TestCache_SaveLoad_PreservesGraph(t *testing.T) {
	dir := t.TempDir()
	cacheDir := filepath.Join(dir, ".markata")

	// Create and populate cache
	cache := New(cacheDir)
	cache.SetDependencies("pages/post-a.md", "post-a", []string{"post-b", "post-c"})
	cache.SetDependencies("pages/post-b.md", "post-b", []string{"post-c"})
	cache.MarkRebuilt("pages/post-a.md", "hash-a", "output/post-a/index.html", "post.html")

	// Save
	if err := cache.Save(); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Load into new cache
	loaded, err := Load(cacheDir)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Verify graph was preserved
	deps := loaded.Graph.GetDependencies("pages/post-a.md")
	if len(deps) != 2 {
		t.Errorf("Loaded cache has %d deps for post-a, want 2", len(deps))
	}

	// Verify PathToSlug was preserved (needed for transitive lookups)
	slug := loaded.Graph.PathToSlug["pages/post-a.md"]
	if slug != "post-a" {
		t.Errorf("PathToSlug[post-a.md] = %q, want %q", slug, "post-a")
	}

	// Verify Dependents was rebuilt
	if !loaded.Graph.HasDependents("post-c") {
		t.Error("post-c should have dependents after load")
	}

	// Verify affected posts calculation works after load
	affected := loaded.GetAffectedPosts([]string{"post-c"})
	if len(affected) != 2 {
		t.Errorf("GetAffectedPosts after load returned %d posts, want 2: %v", len(affected), affected)
	}
}

func TestCache_LoadMissingGraph(t *testing.T) {
	dir := t.TempDir()
	cacheDir := filepath.Join(dir, ".markata")

	// Create cache file manually without graph field
	if err := os.MkdirAll(cacheDir, 0o755); err != nil {
		t.Fatalf("MkdirAll failed: %v", err)
	}
	cacheFile := filepath.Join(cacheDir, CacheFileName)
	data := `{"version":1,"config_hash":"abc","templates_hash":"def","posts":{}}`
	//nolint:gosec // G306: Test file, 0644 is fine
	if err := os.WriteFile(cacheFile, []byte(data), 0o644); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	// Load - should initialize graph
	loaded, err := Load(cacheDir)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if loaded.Graph == nil {
		t.Error("Graph should be initialized when missing from cache file")
	}
}

func TestCache_GraphSize(t *testing.T) {
	cache := New("")

	if cache.GraphSize() != 0 {
		t.Errorf("GraphSize = %d, want 0", cache.GraphSize())
	}

	cache.SetDependencies("pages/post-a.md", "post-a", []string{"post-b"})
	cache.SetDependencies("pages/post-b.md", "post-b", []string{"post-c"})

	if cache.GraphSize() != 2 {
		t.Errorf("GraphSize = %d, want 2", cache.GraphSize())
	}
}

func TestComputeFeedMembershipHash_Empty(t *testing.T) {
	hash := ComputeFeedMembershipHash(nil)
	if hash != "" {
		t.Errorf("ComputeFeedMembershipHash(nil) = %q, want empty", hash)
	}
	hash = ComputeFeedMembershipHash([]string{})
	if hash != "" {
		t.Errorf("ComputeFeedMembershipHash([]) = %q, want empty", hash)
	}
}

func TestComputeFeedMembershipHash_Deterministic(t *testing.T) {
	// Same slugs in different order should produce the same hash
	hash1 := ComputeFeedMembershipHash([]string{"post-c", "post-a", "post-b"})
	hash2 := ComputeFeedMembershipHash([]string{"post-a", "post-b", "post-c"})
	hash3 := ComputeFeedMembershipHash([]string{"post-b", "post-c", "post-a"})

	if hash1 == "" {
		t.Fatal("ComputeFeedMembershipHash returned empty for non-empty input")
	}
	if hash1 != hash2 {
		t.Errorf("hash mismatch: %q != %q (different order should produce same hash)", hash1, hash2)
	}
	if hash1 != hash3 {
		t.Errorf("hash mismatch: %q != %q (different order should produce same hash)", hash1, hash3)
	}
}

func TestComputeFeedMembershipHash_ChangesOnMembershipChange(t *testing.T) {
	hash1 := ComputeFeedMembershipHash([]string{"post-a", "post-b"})
	hash2 := ComputeFeedMembershipHash([]string{"post-a", "post-b", "post-c"})
	hash3 := ComputeFeedMembershipHash([]string{"post-a"})

	if hash1 == hash2 {
		t.Error("hash should change when a member is added")
	}
	if hash1 == hash3 {
		t.Error("hash should change when a member is removed")
	}
}

func TestComputeFeedMembershipHash_DoesNotMutateInput(t *testing.T) {
	input := []string{"post-c", "post-a", "post-b"}
	original := make([]string, len(input))
	copy(original, input)

	ComputeFeedMembershipHash(input)

	for i, v := range input {
		if v != original[i] {
			t.Errorf("input[%d] = %q, want %q (input was mutated)", i, v, original[i])
		}
	}
}

func TestCache_SetGetFeedMembershipHash(t *testing.T) {
	cache := New("")

	// Getting hash for non-existent post returns empty string
	hash := cache.GetFeedMembershipHash("pages/post-a.md")
	if hash != "" {
		t.Errorf("GetFeedMembershipHash for missing post = %q, want empty", hash)
	}

	// Create a post entry first
	cache.MarkRebuilt("pages/post-a.md", "hash123", "output/post-a/index.html", "post.html")

	// Set and get the membership hash
	cache.SetFeedMembershipHash("pages/post-a.md", "membership-hash-abc")
	hash = cache.GetFeedMembershipHash("pages/post-a.md")
	if hash != "membership-hash-abc" {
		t.Errorf("GetFeedMembershipHash = %q, want %q", hash, "membership-hash-abc")
	}

	// Update the membership hash
	cache.SetFeedMembershipHash("pages/post-a.md", "membership-hash-def")
	hash = cache.GetFeedMembershipHash("pages/post-a.md")
	if hash != "membership-hash-def" {
		t.Errorf("GetFeedMembershipHash after update = %q, want %q", hash, "membership-hash-def")
	}
}

func TestCache_FeedMembershipHash_SaveLoad(t *testing.T) {
	dir := t.TempDir()
	cacheDir := filepath.Join(dir, ".markata")

	// Create cache with feed membership hash
	cache := New(cacheDir)
	cache.MarkRebuilt("pages/post-a.md", "hash-a", "output/post-a/index.html", "post.html")
	cache.SetFeedMembershipHash("pages/post-a.md", "membership-hash-123")

	// Save
	if err := cache.Save(); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Load into new cache
	loaded, err := Load(cacheDir)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Verify feed membership hash was preserved
	hash := loaded.GetFeedMembershipHash("pages/post-a.md")
	if hash != "membership-hash-123" {
		t.Errorf("GetFeedMembershipHash after load = %q, want %q", hash, "membership-hash-123")
	}
}

func TestCache_SetTemplatesFingerprint_UnchangedDoesNotMarkDirty(t *testing.T) {
	cache := New("")
	cache.TemplatesFingerprint = "fingerprint-123"
	cache.dirty = false

	if changed := cache.SetTemplatesFingerprint("fingerprint-123"); changed {
		t.Fatal("SetTemplatesFingerprint returned true for unchanged fingerprint")
	}
	if cache.dirty {
		t.Fatal("dirty = true, want false for unchanged fingerprint")
	}
}

func TestHashDirectoryState_Deterministic(t *testing.T) {
	dir := t.TempDir()
	templatesDir := filepath.Join(dir, "templates")
	if err := os.MkdirAll(templatesDir, 0o755); err != nil {
		t.Fatalf("MkdirAll failed: %v", err)
	}

	filePath := filepath.Join(templatesDir, "base.html")
	if err := os.WriteFile(filePath, []byte("<html>one</html>"), 0o600); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	hash1, err := HashDirectoryState(templatesDir, []string{".html"})
	if err != nil {
		t.Fatalf("HashDirectoryState first call failed: %v", err)
	}
	hash2, err := HashDirectoryState(templatesDir, []string{".html"})
	if err != nil {
		t.Fatalf("HashDirectoryState second call failed: %v", err)
	}
	if hash1 == "" {
		t.Fatal("HashDirectoryState returned empty hash for non-empty directory")
	}
	if hash1 != hash2 {
		t.Fatalf("HashDirectoryState = %q then %q, want deterministic output", hash1, hash2)
	}

	if err := os.WriteFile(filePath, []byte("<html>two and more</html>"), 0o600); err != nil {
		t.Fatalf("WriteFile update failed: %v", err)
	}
	hash3, err := HashDirectoryState(templatesDir, []string{".html"})
	if err != nil {
		t.Fatalf("HashDirectoryState updated call failed: %v", err)
	}
	if hash1 == hash3 {
		t.Fatalf("HashDirectoryState = %q after file change, want different hash", hash3)
	}
}

func TestHashFS_DeterministicAndContentSensitive(t *testing.T) {
	first := fstest.MapFS{
		"post.html":          &fstest.MapFile{Data: []byte("<article>one</article>")},
		"partials/card.html": &fstest.MapFile{Data: []byte("<div>card</div>")},
	}
	second := fstest.MapFS{
		"partials/card.html": &fstest.MapFile{Data: []byte("<div>card</div>")},
		"post.html":          &fstest.MapFile{Data: []byte("<article>one</article>")},
	}

	hash, err := HashFS(first, []string{".html"})
	if err != nil {
		t.Fatalf("HashFS() failed: %v", err)
	}
	if got, err := HashFS(second, []string{".html"}); err != nil || got != hash {
		t.Fatalf("HashFS() = %q, %v; want deterministic hash %q", got, err, hash)
	}
	if changed, err := HashFS(fstest.MapFS{
		"post.html":          &fstest.MapFile{Data: []byte("<article>two</article>")},
		"partials/card.html": &fstest.MapFile{Data: []byte("<div>card</div>")},
	}, []string{".html"}); err != nil || changed == hash {
		t.Fatalf("HashFS() after template change = %q, %v; want changed hash", changed, err)
	}
}

func TestCache_CleanupMermaidSVG_RemovesUnusedFiles(t *testing.T) {
	dir := t.TempDir()
	cache := New(filepath.Join(dir, ".markata"))

	if err := cache.CacheMermaidSVG("pages/keep.md", "keep-hash", "<svg>keep</svg>"); err != nil {
		t.Fatalf("CacheMermaidSVG keep failed: %v", err)
	}
	if err := cache.CacheMermaidSVG("pages/drop.md", "drop-hash", "<svg>drop</svg>"); err != nil {
		t.Fatalf("CacheMermaidSVG drop failed: %v", err)
	}

	cache = New(filepath.Join(dir, ".markata"))
	if got := cache.GetCachedMermaidSVG("pages/keep.md", "keep-hash"); got != "<svg>keep</svg>" {
		t.Fatalf("GetCachedMermaidSVG keep = %q, want keep svg", got)
	}

	removed, err := cache.CleanupMermaidSVG()
	if err != nil {
		t.Fatalf("CleanupMermaidSVG failed: %v", err)
	}
	if removed != 1 {
		t.Fatalf("CleanupMermaidSVG removed %d files, want 1", removed)
	}

	if got := cache.GetCachedMermaidSVG("pages/drop.md", "drop-hash"); got != "" {
		t.Fatalf("GetCachedMermaidSVG drop = %q, want empty", got)
	}
	if got := cache.GetCachedMermaidSVG("pages/keep.md", "keep-hash"); got != "<svg>keep</svg>" {
		t.Fatalf("GetCachedMermaidSVG keep after cleanup = %q, want keep svg", got)
	}
}

func TestCache_CleanupMermaidSVG_AfterUsageResetRemovesPreviouslyUsedFiles(t *testing.T) {
	dir := t.TempDir()
	cache := New(filepath.Join(dir, ".markata"))

	if err := cache.CacheMermaidSVG("pages/keep.md", "keep-hash", "<svg>keep</svg>"); err != nil {
		t.Fatalf("CacheMermaidSVG keep failed: %v", err)
	}
	if err := cache.CacheMermaidSVG("pages/drop.md", "drop-hash", "<svg>drop</svg>"); err != nil {
		t.Fatalf("CacheMermaidSVG drop failed: %v", err)
	}

	if got := cache.GetCachedMermaidSVG("pages/keep.md", "keep-hash"); got != "<svg>keep</svg>" {
		t.Fatalf("GetCachedMermaidSVG keep = %q, want keep svg", got)
	}
	if got := cache.GetCachedMermaidSVG("pages/drop.md", "drop-hash"); got != "<svg>drop</svg>" {
		t.Fatalf("GetCachedMermaidSVG drop = %q, want drop svg", got)
	}

	cache.ResetMermaidSVGUsage()
	if got := cache.GetCachedMermaidSVG("pages/keep.md", "keep-hash"); got != "<svg>keep</svg>" {
		t.Fatalf("GetCachedMermaidSVG keep after reset = %q, want keep svg", got)
	}

	removed, err := cache.CleanupMermaidSVG()
	if err != nil {
		t.Fatalf("CleanupMermaidSVG failed: %v", err)
	}
	if removed != 1 {
		t.Fatalf("CleanupMermaidSVG removed %d files, want 1", removed)
	}

	if got := cache.GetCachedMermaidSVG("pages/drop.md", "drop-hash"); got != "" {
		t.Fatalf("GetCachedMermaidSVG drop after cleanup = %q, want empty", got)
	}
	if got := cache.GetCachedMermaidSVG("pages/keep.md", "keep-hash"); got != "<svg>keep</svg>" {
		t.Fatalf("GetCachedMermaidSVG keep after cleanup = %q, want keep svg", got)
	}
}
