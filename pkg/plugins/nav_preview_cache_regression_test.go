package plugins

import (
	"encoding/json"
	"testing"

	"github.com/WaylonWalker/markata/pkg/buildcache"
	"github.com/WaylonWalker/markata/pkg/models"
)

// TestArchiveNavPreviewBodyStatsDoNotInvalidateUnrelatedPostCache guards the
// production failure behind #1337. /archive/ is a global nav target on
// waylonwalker.com. Its preview must not depend on aggregate word/read-time
// statistics, otherwise an ordinary body edit changes the nav hash and clears
// every cached page.
func TestArchiveNavPreviewBodyStatsDoNotInvalidateUnrelatedPostCache(t *testing.T) {
	config := models.NewConfig()
	config.Nav = []models.NavItem{{Label: "archive", URL: "/archive/"}}

	changed := &models.Post{
		Path:      "pages/changed.md",
		Slug:      "changed",
		Href:      "/changed/",
		Published: true,
		InputHash: "changed-hash",
		Template:  "post.html",
		Extra:     map[string]interface{}{"word_count": 100, "reading_time": 1},
	}
	unrelated := &models.Post{
		Path:      "pages/unrelated.md",
		Slug:      "unrelated",
		Href:      "/unrelated/",
		Published: true,
		InputHash: "unrelated-hash",
		Template:  "post.html",
		Extra:     map[string]interface{}{"word_count": 200, "reading_time": 2},
	}
	posts := []*models.Post{changed, unrelated}
	feeds := []models.FeedConfig{{Slug: "archive", Title: "Archive", Posts: posts}}

	cache := buildcache.New(t.TempDir())
	cache.SetNavPreviewHash(navPreviewHashForTest(t, config, posts, feeds, RandomPostConfig{}))
	cache.MarkRebuiltWithSlug(changed.Path, changed.Slug, changed.InputHash, "output/changed/index.html", changed.Template)
	cache.MarkRebuiltWithSlug(unrelated.Path, unrelated.Slug, unrelated.InputHash, "output/unrelated/index.html", unrelated.Template)

	changed.Extra["word_count"] = 101
	changed.Extra["reading_time"] = 3
	if changedHash := cache.SetNavPreviewHash(navPreviewHashForTest(t, config, posts, feeds, RandomPostConfig{})); changedHash {
		t.Fatal("body-only stats changed the global archive nav-preview hash")
	}
	if cache.ShouldRebuild(unrelated.Path, unrelated.InputHash, unrelated.Template) {
		t.Fatal("unrelated post became non-cacheable after body-only archive edit")
	}

	// Feed membership remains visible in the nav preview. Adding a published
	// post should still change the shared hash so the displayed count is correct.
	added := &models.Post{Path: "pages/added.md", Slug: "added", Href: "/added/", Published: true}
	posts = append(posts, added)
	feeds[0].Posts = posts
	if changedHash := cache.SetNavPreviewHash(navPreviewHashForTest(t, config, posts, feeds, RandomPostConfig{})); !changedHash {
		t.Fatal("archive membership change did not invalidate shared nav preview")
	}
	if !cache.ShouldRebuild(unrelated.Path, unrelated.InputHash, unrelated.Template) {
		t.Fatal("visible archive membership change did not invalidate shared rendered pages")
	}
}

func TestRandomNavPreviewKeepsUnrelatedPageCachedWhenDrawChanges(t *testing.T) {
	config := models.NewConfig()
	config.Nav = []models.NavItem{{Label: "Random", URL: "/random/"}}
	random := RandomPostConfig{Enabled: true, Path: "random"}
	cached := &models.Post{Path: "pages/cached.md", Slug: "cached", Href: "/cached/", Published: true}
	photo := &models.Post{Path: "pages/photo/new.md", Slug: "photo", Href: "/photo/", Published: true}
	before := []*models.Post{cached}
	after := []*models.Post{cached, photo}

	// Verify that the addition changes the actual Random draw.
	if len(eligibleRandomPostHrefs(after, random.ExcludeTags)) != 2 {
		t.Fatal("photo was not eligible for the Random draw")
	}
	cache := buildcache.New(t.TempDir())
	cache.SetNavPreviewHash(navPreviewHashForTest(t, config, before, nil, random))
	cache.MarkRebuiltWithSlug(cached.Path, cached.Slug, "input", "output/cached/index.html", "post.html")
	for _, posts := range [][]*models.Post{after, before} {
		if cache.SetNavPreviewHash(navPreviewHashForTest(t, config, posts, nil, random)) {
			t.Fatal("Random draw change invalidated shared navigation previews")
		}
		if cache.ShouldRebuild(cached.Path, "input", "post.html") {
			t.Fatal("Random draw change invalidated an unrelated cached page")
		}
	}
}

func navPreviewHashForTest(t *testing.T, config *models.Config, posts []*models.Post, feeds []models.FeedConfig, random RandomPostConfig) string {
	t.Helper()
	previews := buildNavPreviews(config, posts, feeds, models.BlogrollConfig{}, random)
	encoded, err := json.Marshal(previews)
	if err != nil {
		t.Fatalf("marshal nav previews: %v", err)
	}
	return buildcache.ContentHash(string(encoded))
}
