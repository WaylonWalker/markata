package plugins

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/WaylonWalker/markata/pkg/buildcache"
	"github.com/WaylonWalker/markata/pkg/encryption"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

func postSemanticHashes(post *models.Post) semanticHashes {
	return semanticHashes{
		feed: computePostFeedItemHash(post), tag: computePostTagIndexHash(post),
		garden: computePostGardenHash(post),
	}
}

func cachedSemanticHashes(cache *buildcache.Cache, path string) semanticHashes {
	feed, tag, garden := cache.GetPostSemanticHashes(path)
	return semanticHashes{feed: feed, tag: tag, garden: garden}
}

func loadBaselinePosts(t *testing.T, cache *buildcache.Cache, dir string, files []string, enabled bool) *lifecycle.Manager {
	t.Helper()
	m := lifecycle.NewManager()
	m.SetConcurrency(4)
	m.Config().ContentDir = dir
	m.Config().Extra["encryption_enabled"] = enabled
	m.Config().Extra["encryption_default_key"] = "baseline"
	m.Cache().Set("build_cache", cache)
	m.SetFiles(files)
	loader := NewLoadPlugin()
	if err := loader.Configure(m); err != nil {
		t.Fatal(err)
	}
	if err := loader.Load(m); err != nil {
		t.Fatal(err)
	}
	return m
}

func canonicalizeBaselineTitles(t *testing.T, m *lifecycle.Manager) {
	t.Helper()
	if err := NewAutoTitlePlugin().Transform(m); err != nil {
		t.Fatal(err)
	}
	if err := NewInlineTitlesPlugin().Transform(m); err != nil {
		t.Fatal(err)
	}
}

// The disabled source case reproduces the 99-post mechanism with two posts:
// no parsed cache, no authored title, stable slug, raw hashes overwritten by
// Load before AutoTitle. The enabled case uses the actual decrypted H1.
func TestSemanticHashBaseline_FreshParseWarmCanonical(t *testing.T) {
	for _, mode := range []string{"source-disabled", "source-enabled", "plain-cache-loss"} {
		t.Run(mode, func(t *testing.T) {
			t.Setenv(EncryptionEnvPrefix+"BASELINE", "baseline-test-key")
			dir, cacheDir := t.TempDir(), t.TempDir()
			enabled := mode == "source-enabled"
			body := "# Exact inferred title\n\nPrivate body.\n"
			if mode != "plain-cache-loss" {
				var err error
				body, err = encryption.EncryptSourceMarkdown(body, "baseline", "baseline-test-key")
				if err != nil {
					t.Fatal(err)
				}
			}
			files := []string{"warm-secret.md", "second-secret.md"}
			for _, path := range files {
				if err := os.WriteFile(filepath.Join(dir, path), []byte("---\ntags: [garden]\n---\n"+body), 0o600); err != nil {
					t.Fatal(err)
				}
			}
			cache := buildcache.New(cacheDir)
			wantHashes := make(map[string]semanticHashes)
			for cycle := 0; cycle < 3; cycle++ {
				if cycle > 0 {
					if mode == "plain-cache-loss" {
						if err := os.RemoveAll(filepath.Join(cacheDir, buildcache.PostCacheDir)); err != nil {
							t.Fatal(err)
						}
					}
					var err error
					cache, err = buildcache.Load(cacheDir)
					if err != nil {
						t.Fatal(err)
					}
				}
				m := loadBaselinePosts(t, cache, dir, files, enabled)
				for _, post := range m.Posts() {
					if post.Title != nil && *post.Title != "" {
						t.Fatalf("cycle %d: fresh parse unexpectedly has title %q", cycle, *post.Title)
					}
				}
				if len(cache.GetChangedSlugs()) != 0 {
					t.Fatal("Load unexpectedly marked slugs")
				}
				canonicalizeBaselineTitles(t, m)
				posts := m.Posts()
				tagsConfig, gardenConfig := &models.TagsConfig{}, &models.GardenConfig{}
				if computeTagsListingHash(posts, tagsConfig, cache) != computeTagsListingHash(posts, tagsConfig, nil) ||
					computeGardenHash(posts, gardenConfig, cache) != computeGardenHash(posts, gardenConfig, nil) {
					t.Fatal("cached tag/garden consumers disagree with canonical posts")
				}
				for _, post := range m.Posts() {
					wantTitle := "Exact inferred title"
					if mode == "source-disabled" {
						wantTitle = map[string]string{"warm-secret.md": "Warm Secret", "second-secret.md": "Second Secret"}[post.Path]
					}
					if post.Title == nil || *post.Title != wantTitle || post.TitleText != wantTitle {
						t.Fatalf("title = %+v / %q, want %q", post.Title, post.TitleText, wantTitle)
					}
					current := postSemanticHashes(post)
					if got := cachedSemanticHashes(cache, post.Path); got != current {
						t.Fatalf("canonical consumer hashes = %+v, want %+v", got, current)
					}
					if cycle > 0 && current != wantHashes[post.Path] {
						t.Fatalf("cycle %d hashes changed: %+v != %+v", cycle, current, wantHashes[post.Path])
					}
					wantHashes[post.Path] = current
					if mode != "plain-cache-loss" && cache.GetCachedPostDataLatest(post.Path) != nil {
						t.Fatal("decrypted parsed post cached")
					}
				}
				if cycle > 0 && len(cache.GetChangedSlugs()) != 0 {
					t.Fatalf("cycle %d false InlineTitles changes: %v", cycle, cache.GetChangedSlugs())
				}
				if cycle == 0 && len(cache.GetChangedSlugs()) != len(files) {
					t.Fatal("empty-triplet new-post baseline suppressed real changes")
				}
				if _, ok := m.Cache().Get(semanticHashBaselineKey); ok {
					t.Fatal("completed InlineTitles retained baseline")
				}
				// No stale prior-build handoff on repeated InlineTitles execution.
				before := cache.GetChangedSlugs()
				if err := NewInlineTitlesPlugin().Transform(m); err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(before, cache.GetChangedSlugs()) {
					t.Fatal("repeated InlineTitles reused stale baseline")
				}
				if err := cache.Save(); err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}

func TestSemanticHashBaseline_IsolationResetFirstCapture(t *testing.T) {
	cache := buildcache.New(t.TempDir())
	m := lifecycle.NewManager()
	m.Cache().Set("build_cache", cache)
	cache.UpdatePostSemanticHashes("post.md", "feed", "tag", "garden")
	resetSemanticHashBaseline(m)
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			captureSemanticHashBaseline(m, cache, "post.md")
		}()
	}
	wg.Wait()
	cache.UpdatePostSemanticHashes("post.md", "intermediate", "", "")
	captureSemanticHashBaseline(m, cache, "post.md")
	if got := takeSemanticHashBaseline(m)["post.md"]; got == nil || *got != (semanticHashes{"feed", "tag", "garden"}) {
		t.Fatalf("first capture lost: %+v", got)
	}
	if takeSemanticHashBaseline(m) != nil {
		t.Fatal("handoff not consumed")
	}
	other := lifecycle.NewManager()
	other.Cache().Set("build_cache", cache)
	if takeSemanticHashBaseline(other) != nil {
		t.Fatal("baseline crossed manager boundary")
	}
	resetSemanticHashBaseline(m)
	captureSemanticHashBaseline(m, cache, "post.md")
	m.Cache().Set("build_cache", buildcache.New(t.TempDir()))
	if takeSemanticHashBaseline(m) != nil {
		t.Fatal("baseline crossed cache instance boundary")
	}
	m.Cache().Set("build_cache", cache)
	resetSemanticHashBaseline(m)
	captureSemanticHashBaseline(m, cache, "post.md")
	// An empty-file Load still resets a stale handoff before its early return.
	if err := NewLoadPlugin().Load(m); err != nil {
		t.Fatal(err)
	}
	if got := takeSemanticHashBaseline(m); len(got) != 0 {
		t.Fatalf("Load retained stale paths: %+v", got)
	}
	resetSemanticHashBaseline(m)
	captureSemanticHashBaseline(m, cache, "post.md")
	m.Cache().Set(CacheKeyInlineRenderer, "invalid renderer")
	if err := NewInlineTitlesPlugin().Transform(m); err == nil {
		t.Fatal("invalid renderer did not report error")
	}
	if _, ok := m.Cache().Get(semanticHashBaselineKey); ok {
		t.Fatal("renderer error retained handoff")
	}
}

func TestSemanticHashBaseline_UnavailableFirstCaptureNotRecaptured(t *testing.T) {
	cache := buildcache.New(t.TempDir())
	cache.MarkRebuilt("post.md", "prior-input", "", "post.html")
	m := lifecycle.NewManager()
	m.Cache().Set("build_cache", cache)
	resetSemanticHashBaseline(m)
	captureSemanticHashBaseline(m, cache, "post.md")
	cache.UpdatePostSemanticHashes("post.md", "intermediate-feed", "intermediate-tag", "intermediate-garden")
	captureSemanticHashBaseline(m, cache, "post.md")
	baseline := takeSemanticHashBaseline(m)
	if value, captured := baseline["post.md"]; !captured || value != nil {
		t.Fatalf("unavailable prior baseline was recaptured after mutation: %+v", baseline)
	}
	resetSemanticHashBaseline(m)
	captureSemanticHashBaseline(m, cache, "post.md")
	if value := takeSemanticHashBaseline(m)["post.md"]; value == nil ||
		*value != (semanticHashes{"intermediate-feed", "intermediate-tag", "intermediate-garden"}) {
		t.Fatal("next Load did not capture the now-complete baseline")
	}
}

func TestSemanticHashBaseline_NonstandardAndOtherProducerChanges(t *testing.T) {
	for _, mode := range []string{"missing-auto-title", "skip", "root", "synthetic", "no-load"} {
		t.Run(mode, func(t *testing.T) {
			cache := buildcache.New(t.TempDir())
			m := lifecycle.NewManager()
			m.Cache().Set("build_cache", cache)
			post := models.NewPost("post.md")
			post.Title = new("Title")
			post.Slug = "post"
			switch mode {
			case "missing-auto-title":
				post.Title = nil
			case "skip":
				post.Skip = true
			case "root":
				post.Slug = ""
				post.Set("_slug_explicit", true)
			case "synthetic":
				post.Path = ""
			}
			m.AddPost(post)
			if mode != "no-load" {
				resetSemanticHashBaseline(m)
				captureSemanticHashBaseline(m, cache, post.Path)
			}
			cache.MarkSlugChanged("other-producer")
			cache.MarkFeedSlugChanged("other-producer")
			lifecycle.SetServeAffectedPaths(m, map[string]bool{"other.md": true})
			if err := NewInlineTitlesPlugin().Transform(m); err != nil {
				t.Fatal(err)
			}
			if _, ok := m.Cache().Get(semanticHashBaselineKey); ok {
				t.Fatal("nonstandard pipeline retained handoff")
			}
			if !lifecycle.GetServeAffectedPaths(m)["other.md"] {
				t.Fatal("other producer affected paths lost")
			}
			found := false
			for _, slug := range cache.GetChangedSlugs() {
				found = found || slug == "other-producer"
			}
			if !found {
				t.Fatal("other producer slug lost")
			}
			if mode == "root" && post.Slug != "" {
				t.Fatal("explicit root slug changed")
			}
		})
	}
}

func TestSemanticHashBaseline_LegitimateCanonicalEdits(t *testing.T) {
	for _, edit := range []string{"title", "tags", "privacy", "slug"} {
		t.Run(edit, func(t *testing.T) {
			cacheDir := t.TempDir()
			cache := buildcache.New(cacheDir)
			m := lifecycle.NewManager()
			m.Cache().Set("build_cache", cache)
			post := models.NewPost("post.md")
			post.Title = new("Original")
			post.Slug = "post"
			m.AddPost(post)
			canonicalizeBaselineTitles(t, m)
			if err := cache.Save(); err != nil {
				t.Fatal(err)
			}

			var err error
			cache, err = buildcache.Load(cacheDir)
			if err != nil {
				t.Fatal(err)
			}
			m.Cache().Set("build_cache", cache)
			resetSemanticHashBaseline(m)
			captureSemanticHashBaseline(m, cache, post.Path)
			switch edit {
			case "title":
				post.Title = new("Changed")
			case "tags":
				post.Tags = []string{"changed"}
			case "privacy":
				post.Private = true
			case "slug":
				post.Slug = "changed"
			}
			if err := NewInlineTitlesPlugin().Transform(m); err != nil {
				t.Fatal(err)
			}
			if len(cache.GetChangedSlugs()) == 0 {
				t.Fatalf("%s edit not marked", edit)
			}
		})
	}
}

func TestSemanticHashBaseline_InputEditsRetainDependentInvalidation(t *testing.T) {
	original := "---\ntemplate: post.html\ntags: [original]\n---\n# Original heading\n\nOriginal body.\n"
	for _, edit := range []string{"body", "H1", "title", "tags", "privacy", "slug", "template"} {
		t.Run(edit, func(t *testing.T) {
			dir, cacheDir := t.TempDir(), t.TempDir()
			files := []string{"target.md", "direct.md", "transitive.md"}
			for _, path := range files {
				source := original
				if path != "target.md" {
					source = "---\ntemplate: post.html\n---\n# Reference\n\nUnchanged reference.\n"
				}
				if err := os.WriteFile(filepath.Join(dir, path), []byte(source), 0o600); err != nil {
					t.Fatal(err)
				}
			}
			cache := buildcache.New(cacheDir)
			cold := loadBaselinePosts(t, cache, dir, files, false)
			canonicalizeBaselineTitles(t, cold)
			for _, post := range cold.Posts() {
				cache.MarkRebuiltWithSlug(post.Path, post.Slug, post.InputHash, "", post.Template)
			}
			cache.Graph.SetDependencies("direct.md", "direct", []string{"target"})
			cache.Graph.SetDependencies("transitive.md", "transitive", []string{"direct"})
			if err := cache.Save(); err != nil {
				t.Fatal(err)
			}
			changed := original
			switch edit {
			case "body":
				changed = strings.ReplaceAll(changed, "Original body.", "Changed body.")
			case "H1":
				changed = strings.ReplaceAll(changed, "# Original heading", "# Changed heading")
			case "title":
				changed = strings.ReplaceAll(changed, "template: post.html", "template: post.html\ntitle: Explicit changed title")
			case "tags":
				changed = strings.ReplaceAll(changed, "[original]", "[changed]")
			case "privacy":
				changed = strings.ReplaceAll(changed, "template: post.html", "template: post.html\nprivate: true")
			case "slug":
				// Keep references to the old route resolvable through the
				// supported alias contract. Removing an old slug entirely has
				// a pre-existing dependency-closure gap outside this handoff.
				changed = strings.ReplaceAll(changed, "template: post.html", "template: post.html\nslug: changed\naliases: [target]")
			case "template":
				changed = strings.ReplaceAll(changed, "template: post.html", "template: changed.html")
			}
			if err := os.WriteFile(filepath.Join(dir, "target.md"), []byte(changed), 0o600); err != nil {
				t.Fatal(err)
			}
			var err error
			cache, err = buildcache.Load(cacheDir)
			if err != nil {
				t.Fatal(err)
			}
			// Force target parsing even on filesystems with coarse timestamps.
			cache.Posts["target.md"].ModTime = 0
			m := loadBaselinePosts(t, cache, dir, files, false)
			lifecycle.SetServeAffectedPaths(m, map[string]bool{"preaffected.md": true})
			bc := NewBuildCachePlugin()
			bc.cache = cache
			if err := bc.Load(m); err != nil {
				t.Fatal(err)
			}
			canonicalizeBaselineTitles(t, m)
			affected := lifecycle.GetServeAffectedPaths(m)
			for _, path := range append(files, "preaffected.md") {
				if !affected[path] {
					t.Fatalf("%s edit lost affected path %s: %v", edit, path, affected)
				}
			}
			for _, post := range m.Posts() {
				if cachedSemanticHashes(cache, post.Path) != postSemanticHashes(post) {
					t.Fatalf("%s edit persisted intermediate hashes for %s", edit, post.Path)
				}
			}
		})
	}
}

func TestSemanticHashBaseline_ColdNavResetKeepsUnrelatedPageCacheable(t *testing.T) {
	contentDir, cacheDir, outputDir, templatesDir := t.TempDir(), t.TempDir(), t.TempDir(), t.TempDir()
	if err := os.WriteFile(filepath.Join(templatesDir, "post.html"), []byte(`PAGE:{{ body | safe }}`), 0o600); err != nil {
		t.Fatal(err)
	}
	files := []string{"target.md", "dependent.md", "unrelated.md"}
	for _, path := range files {
		source := "---\ntitle: Explicit title\npublished: true\ntemplate: post.html\n---\nOriginal body.\n"
		if err := os.WriteFile(filepath.Join(contentDir, path), []byte(source), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	build := func() *lifecycle.Manager {
		t.Helper()
		m := lifecycle.NewManager()
		m.Config().ContentDir, m.Config().OutputDir = contentDir, outputDir
		m.Config().Extra["cache_dir"] = cacheDir
		m.Config().Extra["templates_dir"] = templatesDir
		m.SetFiles(files)
		bc, loader := NewBuildCachePlugin(), NewLoadPlugin()
		markdown, pages := NewRenderMarkdownPlugin(), NewTemplatesPlugin()
		for _, configure := range []func(*lifecycle.Manager) error{bc.Configure, loader.Configure, markdown.Configure, pages.Configure} {
			if err := configure(m); err != nil {
				t.Fatal(err)
			}
		}
		if err := loader.Load(m); err != nil {
			t.Fatal(err)
		}
		for _, post := range m.Posts() {
			if post.Path == "dependent.md" {
				post.Dependencies = []string{"target"}
			}
		}
		for _, stage := range []func(*lifecycle.Manager) error{
			bc.Load, NewAutoTitlePlugin().Transform, NewInlineTitlesPlugin().Transform,
			bc.Transform, markdown.Render, pages.Render, NewPublishHTMLPlugin().Write,
		} {
			if err := stage(m); err != nil {
				t.Fatal(err)
			}
		}
		if err := GetBuildCache(m).Save(); err != nil {
			t.Fatal(err)
		}
		return m
	}
	cold := build()
	cache := GetBuildCache(cold)
	if !cold.ContentDiagnostics().TemplateCache.NavPreviewReset {
		t.Fatal("fixture did not exercise the cold navigation reset")
	}
	for _, path := range files {
		entry := cache.GetCachedPost(path)
		if cachedSemanticHashes(cache, path) != (semanticHashes{}) || entry.InputHash == "" || entry.FullHTMLPath == "" || entry.ModTime != 0 {
			t.Fatalf("cold reset did not leave a rendered entry with missing semantic metadata: %s", path)
		}
	}
	source, err := os.ReadFile(filepath.Join(contentDir, "target.md"))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(contentDir, "target.md"), bytes.ReplaceAll(source, []byte("Original body."), []byte("Changed body.")), 0o600); err != nil {
		t.Fatal(err)
	}
	edited := build()
	stats := edited.ContentDiagnostics().TemplateCache
	if stats.NavPreviewReset || stats.Cacheable != 1 || stats.Restored != 1 || stats.RenderRequired != 2 ||
		stats.RenderSucceeded != 2 || stats.MissReasons.AffectedPath != 2 || stats.MissReasons.SlugChanged != 0 {
		t.Fatalf("cold -> edit expanded the template invalidation scope: %+v", stats)
	}
	if affected := lifecycle.GetServeAffectedPaths(edited); len(affected) != 2 || !affected["target.md"] || !affected["dependent.md"] {
		t.Fatalf("input/dependency affected paths changed: %v", affected)
	}
	for _, post := range edited.Posts() {
		if cachedSemanticHashes(GetBuildCache(edited), post.Path) != postSemanticHashes(post) {
			t.Fatalf("canonical hashes not repaired for %s", post.Path)
		}
	}
}

func TestSemanticHashBaseline_PartialEmptyTriplet(t *testing.T) {
	cache := buildcache.New(t.TempDir())
	m := lifecycle.NewManager()
	m.Cache().Set("build_cache", cache)
	post := models.NewPost("new.md")
	post.Title, post.Slug = new("New title"), "new"
	m.AddPost(post)
	// Partial entries still treat explicitly empty tag/garden hashes as changes.
	cache.UpdatePostSemanticHashes(post.Path, computePostFeedItemHash(post), "", "")
	resetSemanticHashBaseline(m)
	captureSemanticHashBaseline(m, cache, post.Path)
	current := postSemanticHashes(post)
	cache.UpdatePostSemanticHashes(post.Path, current.feed, current.tag, current.garden)
	if err := NewInlineTitlesPlugin().Transform(m); err != nil {
		t.Fatal(err)
	}
	if changed := cache.GetChangedSlugs(); len(changed) != 1 || changed[0] != post.Slug {
		t.Fatalf("partial baseline suppressed empty-hash changes: %v", changed)
	}
}

func TestSemanticHashBaseline_LoadOnlyRunToAndRepeatedLoad(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "post.md"), []byte("# Inferred\n\nBody.\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	cache := buildcache.New(t.TempDir())
	m := lifecycle.NewManager()
	m.Config().ContentDir = dir
	m.Cache().Set("build_cache", cache)
	m.SetFiles([]string{"post.md"})
	loader := NewLoadPlugin()
	m.RegisterPlugin(loader)
	if err := m.RunTo(lifecycle.StageLoad); err != nil {
		t.Fatal(err)
	}
	post := m.Posts()[0]
	if post.Title != nil || cachedSemanticHashes(cache, post.Path) != postSemanticHashes(post) {
		t.Fatal("RunTo(Load) no longer persists raw semantic hashes")
	}
	if len(cache.GetChangedFeedSlugs()) == 0 || !cache.TagsDirty() || !cache.GardenDirty() {
		t.Fatal("Load-only conservative invalidation signals were cleared")
	}
	canonicalizeBaselineTitles(t, m)
	canonical := postSemanticHashes(post)
	// Force a fresh parse without relying on filesystem timestamp granularity.
	cache.Posts[post.Path].ModTime = 0
	m.SetPosts(nil)
	if err := loader.Load(m); err != nil {
		t.Fatal(err)
	}
	if got := takeSemanticHashBaseline(m)[post.Path]; got == nil || *got != canonical {
		t.Fatalf("repeated Load baseline = %+v, want %+v", got, canonical)
	}
}
