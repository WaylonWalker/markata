package plugins

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/WaylonWalker/markata/pkg/buildcache"
	"github.com/WaylonWalker/markata/pkg/diagnostics"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

func TestTemplateCacheReasonGatesAndPrecedence(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*models.Post, *buildcache.Cache, map[string]bool, map[string]string, map[string]bool)
		noCache bool
		want    templateCacheReason
	}{
		{name: "hit", want: templateCacheHit},
		{name: "affected masks unavailable", noCache: true, setup: func(post *models.Post, _ *buildcache.Cache, _ map[string]bool, _ map[string]string, affected map[string]bool) {
			affected[post.Path] = true
		}, want: templateCacheAffectedPath},
		{name: "cache unavailable masks input missing", noCache: true, setup: func(post *models.Post, _ *buildcache.Cache, _ map[string]bool, _ map[string]string, _ map[string]bool) {
			post.InputHash = ""
		}, want: templateCacheUnavailable},
		{name: "input missing masks entry missing", setup: func(post *models.Post, _ *buildcache.Cache, _ map[string]bool, _ map[string]string, _ map[string]bool) {
			post.InputHash = ""
			post.Path = "absent.md"
		}, want: templateCacheInputHashMissing},
		{name: "entry missing", setup: func(post *models.Post, _ *buildcache.Cache, _ map[string]bool, _ map[string]string, _ map[string]bool) {
			post.Path = "absent.md"
		}, want: templateCacheEntryMissing},
		{name: "input mismatch masks template", setup: func(post *models.Post, _ *buildcache.Cache, _ map[string]bool, _ map[string]string, _ map[string]bool) {
			post.InputHash = "new"
			post.Template = "new.html"
		}, want: templateCacheInputHashMismatch},
		{name: "partial cached input", setup: func(post *models.Post, cache *buildcache.Cache, _ map[string]bool, _ map[string]string, _ map[string]bool) {
			cache.MarkRebuilt(post.Path, "", "", post.Template)
		}, want: templateCacheInputHashMismatch},
		{name: "template masks dependency", setup: func(post *models.Post, _ *buildcache.Cache, changed map[string]bool, _ map[string]string, _ map[string]bool) {
			post.Template = "new.html"
			changed["dep"] = true
		}, want: templateCacheTemplateMismatch},
		{name: "dependency masks own slug", setup: func(post *models.Post, _ *buildcache.Cache, changed map[string]bool, _ map[string]string, _ map[string]bool) {
			changed["dep"] = true
			changed[post.Slug] = true
		}, want: templateCacheDependencyChanged},
		{name: "own slug masks feed", setup: func(post *models.Post, _ *buildcache.Cache, changed map[string]bool, membership map[string]string, _ map[string]bool) {
			changed[post.Slug] = true
			membership["tag"] = "new"
		}, want: templateCacheSlugChanged},
		{name: "feed masks local", setup: func(post *models.Post, cache *buildcache.Cache, _ map[string]bool, membership map[string]string, _ map[string]bool) {
			membership["tag"] = "new"
			cache.SetLocalPreviewHash(post.Path, "new")
		}, want: templateCacheFeedMembershipChanged},
		{name: "local last", setup: func(post *models.Post, cache *buildcache.Cache, _ map[string]bool, _ map[string]string, _ map[string]bool) {
			cache.SetLocalPreviewHash(post.Path, "new")
		}, want: templateCacheLocalPreviewChanged},
		{name: "empty current feed bypasses cached", setup: func(post *models.Post, cache *buildcache.Cache, _ map[string]bool, membership map[string]string, _ map[string]bool) {
			membership["tag"] = ""
			cache.SetFeedMembershipHash(post.Path, "old")
		}, want: templateCacheHit},
		{name: "unrelated changed slug", setup: func(_ *models.Post, _ *buildcache.Cache, changed map[string]bool, _ map[string]string, _ map[string]bool) {
			changed["other"] = true
		}, want: templateCacheHit},
		{name: "root empty slug matches changed set", setup: func(post *models.Post, _ *buildcache.Cache, changed map[string]bool, _ map[string]string, _ map[string]bool) {
			post.Slug = ""
			changed[""] = true
		}, want: templateCacheSlugChanged},
		{name: "empty path uses exact key", setup: func(post *models.Post, cache *buildcache.Cache, _ map[string]bool, _ map[string]string, _ map[string]bool) {
			post.Path = ""
			cache.MarkRebuilt("", post.InputHash, "", post.Template)
			cache.SetFeedMembershipHash("", "same")
			cache.SetLocalPreviewHash("", localPreviewHash(post, nil, ""))
		}, want: templateCacheHit},
		{name: "affected full build masks metadata", setup: func(post *models.Post, _ *buildcache.Cache, _ map[string]bool, _ map[string]string, affected map[string]bool) {
			affected[post.Path] = true
			post.InputHash = ""
		}, want: templateCacheAffectedPath},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			post := &models.Post{Path: "post.md", Slug: "post", InputHash: "hash", Template: "post.html", Dependencies: []string{"dep"}, Tags: []string{"tag"}}
			cache := buildcache.New(t.TempDir())
			cache.MarkRebuilt(post.Path, post.InputHash, "", post.Template)
			cache.SetFeedMembershipHash(post.Path, "same")
			cache.SetLocalPreviewHash(post.Path, localPreviewHash(post, nil, ""))
			changed := map[string]bool{}
			membership := map[string]string{"tag": "same"}
			affected := map[string]bool{}
			if test.setup != nil {
				test.setup(post, cache, changed, membership, affected)
			}
			if test.noCache {
				cache = nil
			}
			p := NewTemplatesPlugin()
			if got := p.templateCacheReason(post, cache, changed, membership, affected); got != test.want {
				t.Fatalf("reason = %v, want %v", got, test.want)
			}
			if got, want := canUseCachedHTML(post, cache, changed, membership), legacyCanUseCachedHTML(post, cache, changed, membership); got != want {
				t.Fatalf("bool changed: got %v, want %v", got, want)
			}
			legacy := !affected[post.Path] && legacyCanUseCachedHTML(post, cache, changed, membership) &&
				cache.GetLocalPreviewHash(post.Path) == localPreviewHash(post, nil, "")
			if legacy != (test.want == templateCacheHit) {
				t.Fatalf("outer bool changed: %v vs reason %v", legacy, test.want)
			}
		})
	}
}

func assertTemplateCacheReconciles(t *testing.T, stats *diagnostics.TemplateCacheStats, postCount int) {
	t.Helper()
	if stats == nil {
		t.Fatal("template stats absent")
	}
	r := stats.MissReasons
	phase1a := r.AffectedPath + r.CacheUnavailable + r.InputHashMissing + r.EntryMissing +
		r.InputHashMismatch + r.TemplateMismatch + r.DependencyChanged + r.SlugChanged +
		r.FeedMembershipChanged + r.LocalPreviewChanged
	if postCount != stats.Skipped+stats.Classified ||
		stats.Classified-stats.Cacheable != phase1a ||
		stats.Cacheable != stats.Restored+r.FullHTMLUnavailable ||
		stats.RenderRequired != phase1a+r.FullHTMLUnavailable ||
		stats.Classified != stats.Restored+stats.RenderRequired ||
		stats.RenderRequired != stats.RenderSucceeded+stats.RenderFailed+stats.ServeDeferred {
		t.Fatalf("template counters do not reconcile for %d posts: %+v", postCount, stats)
	}
}

func TestTemplateCacheFirstMissCounterMapping(t *testing.T) {
	tests := []struct {
		reason templateCacheReason
		want   diagnostics.TemplateCacheMissReasons
	}{
		{templateCacheHit, diagnostics.TemplateCacheMissReasons{}},
		{templateCacheAffectedPath, diagnostics.TemplateCacheMissReasons{AffectedPath: 1}},
		{templateCacheUnavailable, diagnostics.TemplateCacheMissReasons{CacheUnavailable: 1}},
		{templateCacheInputHashMissing, diagnostics.TemplateCacheMissReasons{InputHashMissing: 1}},
		{templateCacheEntryMissing, diagnostics.TemplateCacheMissReasons{EntryMissing: 1}},
		{templateCacheInputHashMismatch, diagnostics.TemplateCacheMissReasons{InputHashMismatch: 1}},
		{templateCacheTemplateMismatch, diagnostics.TemplateCacheMissReasons{TemplateMismatch: 1}},
		{templateCacheDependencyChanged, diagnostics.TemplateCacheMissReasons{DependencyChanged: 1}},
		{templateCacheSlugChanged, diagnostics.TemplateCacheMissReasons{SlugChanged: 1}},
		{templateCacheFeedMembershipChanged, diagnostics.TemplateCacheMissReasons{FeedMembershipChanged: 1}},
		{templateCacheLocalPreviewChanged, diagnostics.TemplateCacheMissReasons{LocalPreviewChanged: 1}},
	}
	for _, test := range tests {
		var got diagnostics.TemplateCacheMissReasons
		recordTemplateCacheMiss(&got, test.reason)
		if got != test.want {
			t.Fatalf("reason %v: counters %+v, want %+v", test.reason, got, test.want)
		}
	}
}

func TestTemplateCacheRenderRestorationMissingEmptyAndSkip(t *testing.T) {
	cache := buildcache.New(t.TempDir())
	hit := headingMigrationPost("hit.md", "", "")
	hit.Href = "/"
	missing := headingMigrationPost("missing.md", "missing", "")
	empty := headingMigrationPost("empty.md", "empty", "")
	skip := headingMigrationPost("skip.md", "skip", "")
	skip.Skip = true
	m, p := headingMigrationManager(t, cache, hit, missing, empty, skip)
	m.Config().Extra["components"] = models.ComponentsConfig{
		Nav: models.NavComponentConfig{Items: []models.NavItem{{Label: "Home", URL: hit.Href}}},
	}
	// Establish the shared navigation hash, then seed metadata without recording
	// changed slugs. A second invocation must replace, not accumulate, counters.
	if err := p.Render(m); err != nil {
		t.Fatal(err)
	}
	if stats := m.ContentDiagnostics().TemplateCache; stats == nil || !stats.NavPreviewReset {
		t.Fatalf("initial nav reset absent: %+v", stats)
	}
	for _, post := range []*models.Post{hit, missing, empty, skip} {
		cache.MarkRebuilt(post.Path, post.InputHash, "", post.Template)
		post.HTML = "UNTOUCHED"
	}
	// Drop the initial rendered metadata before creating a genuinely missing
	// full-page entry, while retaining the shared navigation identity.
	cache.RemoveStale(map[string]bool{hit.Path: true, empty.Path: true, skip.Path: true})
	cache.MarkRebuilt(missing.Path, missing.InputHash, "", missing.Template)
	expected := "PAGE:\x00\xff<exact>"
	if err := cache.CacheFullHTML(hit.Path, expected); err != nil {
		t.Fatal(err)
	}
	if err := cache.CacheFullHTML(empty.Path, ""); err != nil {
		t.Fatal(err)
	}
	// Missing has an entry but no full page, empty has a zero-length full page.
	if err := p.Render(m); err != nil {
		t.Fatal(err)
	}
	stats := m.ContentDiagnostics().TemplateCache
	assertTemplateCacheReconciles(t, stats, 4)
	want := diagnostics.TemplateCacheStats{
		Classified: 3, Skipped: 1, Cacheable: 3, Restored: 1,
		RenderRequired: 2, RenderSucceeded: 2,
		MissReasons: diagnostics.TemplateCacheMissReasons{FullHTMLUnavailable: 2},
	}
	if *stats != want {
		t.Fatalf("stats = %+v, want %+v", *stats, want)
	}
	if hit.HTML != expected || missing.HTML != "PAGE:" || empty.HTML != "PAGE:" || skip.HTML != "UNTOUCHED" {
		t.Fatalf("HTML behavior changed: %q %q %q %q", hit.HTML, missing.HTML, empty.HTML, skip.HTML)
	}
	// Changing public metadata resets navigation, but reports actual metadata
	// misses rather than attributing each miss to navigation.
	hit.Description = new("new public description")
	if err := p.Render(m); err != nil {
		t.Fatal(err)
	}
	stats = m.ContentDiagnostics().TemplateCache
	assertTemplateCacheReconciles(t, stats, 4)
	if !stats.NavPreviewReset || stats.MissReasons.EntryMissing != 3 {
		t.Fatalf("navigation context did not reflect reset: %+v", stats)
	}
}

func TestTemplateCacheRenderFailuresAndServeDeferred(t *testing.T) {
	cache := buildcache.New(t.TempDir())
	good := headingMigrationPost("good.md", "good", "")
	bad := headingMigrationPost("bad.md", "bad", "")
	bad.Template = "bad.html"
	bad.HTML = "OLD"
	deferred := headingMigrationPost("deferred.md", "deferred", "")
	hit := headingMigrationPost("hit.md", "hit", "")
	m, p := headingMigrationManager(t, cache, good, bad, deferred, hit)
	// A runtime template failure, not a setup error, permits successful siblings.
	dir, ok := m.Config().Extra["templates_dir"].(string)
	if !ok {
		t.Fatal("fixture template directory missing")
	}
	if err := os.WriteFile(filepath.Join(dir, "bad.html"),
		[]byte(`{% include "does-not-exist.html" %}`), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := p.Render(m); err == nil {
		t.Fatal("bad template should fail")
	}
	cache.MarkRebuilt(hit.Path, hit.InputHash, "", hit.Template)
	if err := cache.CacheFullHTML(hit.Path, "HIT"); err != nil {
		t.Fatal(err)
	}
	lifecycle.SetServeIncremental(m, true)
	lifecycle.SetServeAffectedPaths(m, map[string]bool{good.Path: true, bad.Path: true, deferred.Path: true})
	m.Cache().Set(canonicalArticlePathsKey, map[string]bool{good.Path: true, bad.Path: true})
	if err := p.Render(m); err == nil {
		t.Fatal("bad template should fail")
	}
	stats := m.ContentDiagnostics().TemplateCache
	assertTemplateCacheReconciles(t, stats, 4)
	want := diagnostics.TemplateCacheStats{
		Classified: 4, Cacheable: 1, Restored: 1, RenderRequired: 3,
		RenderSucceeded: 1, RenderFailed: 1, ServeDeferred: 1,
		MissReasons: diagnostics.TemplateCacheMissReasons{AffectedPath: 3},
	}
	if *stats != want {
		t.Fatalf("stats = %+v, want %+v", *stats, want)
	}
	if good.HTML != "PAGE:" || bad.HTML != "OLD" || hit.HTML != "HIT" {
		t.Fatalf("render or noncanonical restore changed: %q %q %q", good.HTML, bad.HTML, hit.HTML)
	}
}

func TestTemplateCacheMissingFullHTMLServeDeferred(t *testing.T) {
	cache := buildcache.New(t.TempDir())
	post := headingMigrationPost("deferred.md", "deferred", "")
	m, p := headingMigrationManager(t, cache, post)
	if err := p.Render(m); err != nil {
		t.Fatal(err)
	}
	cache.RemoveStale(nil)
	cache.MarkRebuilt(post.Path, post.InputHash, "", post.Template)
	lifecycle.SetServeIncremental(m, true)
	m.Cache().Set(canonicalArticlePathsKey, map[string]bool{})
	post.HTML = "UNCHANGED"
	if err := p.Render(m); err != nil {
		t.Fatal(err)
	}
	stats := m.ContentDiagnostics().TemplateCache
	assertTemplateCacheReconciles(t, stats, 1)
	want := diagnostics.TemplateCacheStats{
		Classified: 1, Cacheable: 1, RenderRequired: 1, ServeDeferred: 1,
		MissReasons: diagnostics.TemplateCacheMissReasons{FullHTMLUnavailable: 1},
	}
	if *stats != want || post.HTML != "UNCHANGED" {
		t.Fatalf("full-page fallback should defer once: %+v, %q", stats, post.HTML)
	}
}

func TestTemplateCacheRenderStorageFailureRemainsSuccess(t *testing.T) {
	dir := t.TempDir()
	cache := buildcache.New(dir)
	post := headingMigrationPost("post.md", "post", "")
	m, p := headingMigrationManager(t, cache, post)
	if err := os.WriteFile(filepath.Join(dir, buildcache.FullHTMLCacheDir), []byte("block mkdir"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := p.Render(m); err != nil {
		t.Fatal(err)
	}
	stats := m.ContentDiagnostics().TemplateCache
	assertTemplateCacheReconciles(t, stats, 1)
	if stats.RenderSucceeded != 1 || stats.RenderFailed != 0 || post.HTML != "PAGE:" {
		t.Fatalf("storage failure changed render semantics: %+v, %q", stats, post.HTML)
	}
	if cache.GetHeadingHighlightRevision(post.Path) != 0 {
		t.Fatal("failed storage certified heading revision")
	}
}

func TestTemplateCacheRenderClearsEarlyErrorsAndManagerIsolation(t *testing.T) {
	first := lifecycle.NewManager()
	second := lifecycle.NewManager()
	old := diagnostics.TemplateCacheStats{Classified: 100}
	first.ContentLedger().SetTemplateCache(&old)
	second.ContentLedger().SetTemplateCache(&old)
	if err := NewTemplatesPlugin().Render(first); err == nil {
		t.Fatal("unconfigured engine should fail")
	}
	if first.ContentDiagnostics().TemplateCache != nil || second.ContentDiagnostics().TemplateCache == nil {
		t.Fatal("early error did not clear only its own manager")
	}
	m, p := headingMigrationManager(t, nil)
	m.ContentLedger().SetTemplateCache(&old)
	m.Cache().Set(canonicalArticlePathsKey, "wrong type")
	if err := p.Render(m); err == nil {
		t.Fatal("invalid canonical type should fail")
	}
	if m.ContentDiagnostics().TemplateCache != nil {
		t.Fatal("setup failure retained stale stats")
	}
	m.Cache().Set(canonicalArticlePathsKey, map[string]bool{})
	if err := p.Render(m); err != nil {
		t.Fatal(err)
	}
	if got := m.ContentDiagnostics().TemplateCache; got == nil || !reflect.DeepEqual(*got, diagnostics.TemplateCacheStats{}) {
		t.Fatalf("empty completed invocation must publish zero stats: %+v", got)
	}
}

// legacyCanUseCachedHTML preserves the pre-telemetry boolean implementation for
// equivalence tests and a direct same-input benchmark, not a reason-wrapper proxy.
func legacyCanUseCachedHTML(post *models.Post, cache *buildcache.Cache, changedSlugs map[string]bool, feedMembershipHashes map[string]string) bool {
	if cache == nil || post.InputHash == "" || cache.ShouldRebuild(post.Path, post.InputHash, post.Template) {
		return false
	}
	if len(changedSlugs) > 0 {
		for _, dep := range post.Dependencies {
			if changedSlugs[dep] {
				return false
			}
		}
		if changedSlugs[post.Slug] {
			return false
		}
	}
	if current := lookupFeedMembershipHash(post, feedMembershipHashes); current != "" && cache.GetFeedMembershipHash(post.Path) != current {
		return false
	}
	return true
}

var templateClassificationSink int

func BenchmarkTemplateCacheClassification(b *testing.B) {
	cache := buildcache.New(b.TempDir())
	p := NewTemplatesPlugin()
	posts := make([]*models.Post, 128)
	for i := range posts {
		posts[i] = &models.Post{Path: "same.md", Slug: "same", InputHash: "hash", Template: "post.html"}
		if i%8 == 0 {
			posts[i].InputHash = "changed"
		}
		if i%8 == 1 {
			posts[i].Path = "absent.md"
		}
	}
	cache.MarkRebuilt("same.md", "hash", "", "post.html")
	for _, mode := range []string{"bool", "reason_and_counts"} {
		b.Run(mode, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				var hits int
				var reasons diagnostics.TemplateCacheMissReasons
				for _, post := range posts {
					if mode == "bool" {
						if legacyCanUseCachedHTML(post, cache, nil, nil) &&
							cache.GetLocalPreviewHash(post.Path) == localPreviewHash(post, nil, "") {
							hits++
						}
					} else {
						reason := p.templateCacheReason(post, cache, nil, nil, nil)
						if reason == templateCacheHit {
							hits++
						} else {
							recordTemplateCacheMiss(&reasons, reason)
						}
					}
				}
				templateClassificationSink = hits + reasons.InputHashMismatch + reasons.EntryMissing
			}
		})
	}
}
