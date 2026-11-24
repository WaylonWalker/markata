package plugins

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/buildcache"
	"github.com/WaylonWalker/markata/pkg/encryption"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

type encryptionPipelineFixture struct {
	content, output, cache, templates string
}

func newEncryptionPipelineFixture(t *testing.T) encryptionPipelineFixture {
	t.Helper()
	f := encryptionPipelineFixture{t.TempDir(), t.TempDir(), t.TempDir(), t.TempDir()}
	if err := os.WriteFile(filepath.Join(f.templates, "post.html"), []byte(`PAGE:{{ body | safe }}`), 0o600); err != nil {
		t.Fatal(err)
	}
	return f
}

func (f encryptionPipelineFixture) source(t *testing.T, frontmatter, password string) {
	t.Helper()
	body, err := encryption.EncryptSourceMarkdown("# Exact inferred title\n\nPrivate body.\n", "rotation", password)
	if err != nil {
		t.Fatal(err)
	}
	source := "---\npublished: true\ntemplate: post.html\ntags: [garden]\n" + frontmatter + "---\n" + body
	if err := os.WriteFile(filepath.Join(f.content, "secret.md"), []byte(source), 0o600); err != nil {
		t.Fatal(err)
	}
}

// Uses the real cache configure/load, templates and publication gates. Only
// the plugin set is narrowed to avoid unrelated network and site work.
func (f encryptionPipelineFixture) build(t *testing.T, password string) (*lifecycle.Manager, *models.Post) {
	t.Helper()
	return f.buildWithFullCacheExpectation(t, password, true)
}

func (f encryptionPipelineFixture) buildWithFullCacheExpectation(t *testing.T, password string, fullCacheAvailable bool) (*lifecycle.Manager, *models.Post) {
	t.Helper()
	t.Setenv(EncryptionEnvPrefix+"ROTATION", password)
	m := lifecycle.NewManager()
	m.SetConcurrency(4)
	m.Config().ContentDir, m.Config().OutputDir = f.content, f.output
	m.Config().Extra["cache_dir"] = f.cache
	m.Config().Extra["templates_dir"] = f.templates
	m.Config().Extra["encryption_enabled"] = true
	m.Config().Extra["encryption_default_key"] = "rotation"
	m.Config().Extra["config_paths"] = []string{filepath.Join(f.content, "nonexistent-config.toml")}
	m.SetFiles([]string{"secret.md"})
	bc := NewBuildCachePlugin()
	loader := NewLoadPlugin()
	markdown := NewRenderMarkdownPlugin()
	crypt := NewEncryptionPlugin()
	pages := NewTemplatesPlugin()
	assets := NewStaticAssetsPlugin()
	for _, configure := range []func(*lifecycle.Manager) error{
		bc.Configure, assets.Configure, loader.Configure, markdown.Configure, crypt.Configure, pages.Configure,
	} {
		if err := configure(m); err != nil {
			t.Fatal(err)
		}
	}
	crypt.enforceStrength = false // Test-only keys; production policy is unchanged.
	for _, stage := range []func(*lifecycle.Manager) error{
		assets.Validate, loader.Load, bc.Load, crypt.Transform, NewAutoTitlePlugin().Transform,
		NewInlineTitlesPlugin().Transform, bc.Transform, markdown.Render,
		crypt.Render, pages.Render, assets.Write, NewPublishHTMLPlugin().Write, assets.Cleanup,
	} {
		if err := stage(m); err != nil {
			t.Fatal(err)
		}
	}
	if len(m.Posts()) != 1 {
		t.Fatalf("posts = %d, want 1", len(m.Posts()))
	}
	post := m.Posts()[0]
	cache := GetBuildCache(m)
	if err := cache.Save(); err != nil {
		t.Fatal(err)
	}
	published, err := os.ReadFile(filepath.Join(f.output, post.Slug, "index.html"))
	if err != nil {
		t.Fatal(err)
	}
	if string(published) != post.HTML {
		t.Fatal("published and rendered full pages differ")
	}
	decrypted, err := encryption.Decrypt(encryptedCiphertext(t, string(published)), password)
	if err != nil || !strings.Contains(string(decrypted), "Private body.") {
		t.Fatalf("disk-published page rejects the current source key: %v", err)
	}
	cachedFull := cache.GetCachedFullHTML(post.Path)
	if fullCacheAvailable && cachedFull != post.HTML {
		t.Fatal("rendered and cached full pages differ")
	}
	if !fullCacheAvailable && (cachedFull != "" || cache.Posts[post.Path].FullHTMLPath != "") {
		t.Fatal("failed fresh full-page cache write left stale full-page reference available")
	}
	if cache.GetCachedPostDataLatest(post.Path) != nil ||
		cache.GetCachedArticleHTML(post.Path, buildcache.ContentHash("# Exact inferred title\n\nPrivate body.\n")) != "" {
		t.Fatal("decrypted parsed/article source cache persisted")
	}
	return m, post
}

func encryptedCiphertext(t *testing.T, html string) string {
	t.Helper()
	_, after, ok := strings.Cut(html, `data-encrypted="`)
	if !ok {
		t.Fatal("encrypted wrapper absent")
	}
	cipher, _, ok := strings.Cut(after, `"`)
	if !ok || cipher == "" {
		t.Fatal("ciphertext absent")
	}
	return cipher
}

func TestEncryptionInvalidation_RotationPublishesCanonicalInput(t *testing.T) {
	for _, tc := range []struct {
		name, frontmatter, title string
	}{
		{"inferred", "", "Exact inferred title"},
		{"explicit", "title: Public title\n", "Public title"},
		{"root", "title: Public root\nslug: ''\n", "Public root"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := newEncryptionPipelineFixture(t)
			f.source(t, tc.frontmatter, "rotation-password-A")
			// Settle cold navigation/asset metadata, then require two
			// consecutive warm full-page restorations before rotating.
			f.build(t, "rotation-password-A")
			cold, a := f.build(t, "rotation-password-A")
			hashes := cachedSemanticHashes(GetBuildCache(cold), a.Path)
			input, oldHTML := a.InputHash, a.HTML
			cache := GetBuildCache(cold)
			configHash, templateHash, assetHash, navHash := cache.ConfigHash, cache.GetTemplatesHash(), cache.AssetsHash, cache.NavPreviewHash
			if assetHash == "" {
				t.Fatal("static-assets cache identity missing")
			}
			for i := 0; i < 2; i++ {
				warm, hit := f.build(t, "rotation-password-A")
				cache := GetBuildCache(warm)
				stats := warm.ContentDiagnostics().TemplateCache
				if hit.HTML != oldHTML || len(lifecycle.GetServeAffectedPaths(warm)) != 0 ||
					len(cache.GetChangedSlugs()) != 0 || stats == nil || stats.Restored != 1 || stats.RenderRequired != 0 ||
					stats.NavPreviewReset || cache.ConfigHash != configHash || cache.GetTemplatesHash() != templateHash ||
					cache.AssetsHash != assetHash || cache.NavPreviewHash != navHash {
					t.Fatalf("warm A build %d has a restoration/invalidation confound: %+v", i+1, stats)
				}
			}
			// Rewrite only the randomized source envelope under a new password.
			f.source(t, tc.frontmatter, "rotation-password-B")
			rotated, b := f.build(t, "rotation-password-B")
			if b.InputHash != input || cachedSemanticHashes(GetBuildCache(rotated), b.Path) != hashes {
				t.Fatal("rotation changed canonical input or semantic hashes")
			}
			if b.Title == nil || *b.Title != tc.title {
				t.Fatalf("title = %+v, want %q", b.Title, tc.title)
			}
			if !lifecycle.GetServeAffectedPaths(rotated)[b.Path] || b.HTML == oldHTML {
				t.Fatal("rotation did not refresh and publish the real source path")
			}
			stats := rotated.ContentDiagnostics().TemplateCache
			rotatedCache := GetBuildCache(rotated)
			if stats == nil || stats.MissReasons.AffectedPath != 1 || stats.RenderSucceeded != 1 || stats.NavPreviewReset ||
				rotatedCache.ConfigHash != configHash || rotatedCache.GetTemplatesHash() != templateHash ||
				rotatedCache.AssetsHash != assetHash || rotatedCache.NavPreviewHash != navHash {
				t.Fatalf("rotation templates did not render via affected_path: %+v", stats)
			}
			cipher := encryptedCiphertext(t, GetBuildCache(rotated).GetCachedFullHTML(b.Path))
			decrypted, err := encryption.Decrypt(cipher, "rotation-password-B")
			if err != nil || !strings.Contains(string(decrypted), "Private body.") {
				t.Fatalf("new cached cipher does not decrypt under B: %v", err)
			}
			if _, err := encryption.Decrypt(cipher, "rotation-password-A"); err == nil {
				t.Fatal("rotated published cipher still accepts A")
			}
			stable, c := f.build(t, "rotation-password-B")
			if c.ArticleHTML != b.ArticleHTML || c.HTML != b.HTML ||
				len(lifecycle.GetServeAffectedPaths(stable)) != 0 ||
				len(GetBuildCache(stable).GetChangedSlugs()) != 0 {
				t.Fatal("warm post-rotation cache did not stabilize")
			}
			// Same plaintext/password with fresh randomized source ciphertext
			// must still reuse the same canonical wrapper and page.
			f.source(t, tc.frontmatter, "rotation-password-B")
			reEnvelope, d := f.build(t, "rotation-password-B")
			if d.HTML != c.HTML || len(lifecycle.GetServeAffectedPaths(reEnvelope)) != 0 {
				t.Fatal("random source envelope invalidated canonical wrapper")
			}
			f.source(t, tc.frontmatter, "rotation-password-A")
			back, e := f.build(t, "rotation-password-A")
			if e.InputHash != input || cachedSemanticHashes(GetBuildCache(back), e.Path) != hashes ||
				!lifecycle.GetServeAffectedPaths(back)[e.Path] || e.HTML == d.HTML {
				t.Fatal("B -> A rotation failed to refresh the published page")
			}
			if _, err := encryption.Decrypt(encryptedCiphertext(t, e.HTML), "rotation-password-B"); err == nil {
				t.Fatal("B -> A published page still accepts B")
			}
			for i := 0; i < 2; i++ {
				warm, restored := f.build(t, "rotation-password-A")
				if restored.ArticleHTML != e.ArticleHTML || restored.HTML != e.HTML ||
					len(lifecycle.GetServeAffectedPaths(warm)) != 0 || warm.ContentDiagnostics().TemplateCache.Restored != 1 {
					t.Fatalf("steady A build %d after reverse rotation reused stale output", i+1)
				}
			}
		})
	}
}

func TestEncryptionInvalidation_FullPageWriteFailureReloadRecovery(t *testing.T) {
	for _, frontmatter := range []string{"", "title: Public title\n", "title: Root\nslug: ''\n"} {
		t.Run(strings.TrimSpace(frontmatter), func(t *testing.T) {
			f := newEncryptionPipelineFixture(t)
			f.source(t, frontmatter, "rotation-password-A")
			f.build(t, "rotation-password-A")
			before, old := f.build(t, "rotation-password-A")
			cache := GetBuildCache(before)
			oldFull := cache.Posts[old.Path].FullHTMLPath
			fullDir := filepath.Join(f.cache, buildcache.FullHTMLCacheDir)
			parkedDir := filepath.Join(f.cache, "retained-old-fullhtml")
			if err := os.Rename(fullDir, parkedDir); err != nil {
				t.Fatal(err)
			}
			// Retain a genuinely readable old full page while deterministically
			// blocking the normal full-page write directory (also works as root).
			cache.Posts[old.Path].FullHTMLPath = filepath.Join(parkedDir, filepath.Base(oldFull))
			cache.SetLocalPreviewHash(old.Path, cache.GetLocalPreviewHash(old.Path))
			if err := cache.Save(); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(fullDir, []byte("blocks full-page cache writes"), 0o600); err != nil {
				t.Fatal(err)
			}
			reloaded, err := buildcache.Load(f.cache)
			if err != nil || reloaded.GetCachedFullHTML(old.Path) != old.HTML {
				t.Fatalf("old A page was not readable before rotation: %v", err)
			}
			f.source(t, frontmatter, "rotation-password-B")
			rotated, b := f.buildWithFullCacheExpectation(t, "rotation-password-B", false)
			if b.InputHash != old.InputHash || cachedSemanticHashes(GetBuildCache(rotated), b.Path) != cachedSemanticHashes(reloaded, old.Path) {
				t.Fatal("write-failure rotation changed canonical hashes")
			}
			if b.HTML == old.HTML || !lifecycle.GetServeAffectedPaths(rotated)[b.Path] {
				t.Fatal("write-failure rotation did not publish B")
			}
			persisted, err := buildcache.Load(f.cache)
			if err != nil || persisted.GetCachedFullHTML(b.Path) != "" || persisted.Posts[b.Path].FullHTMLPath != "" {
				t.Fatalf("stale full-page reference survived persistent reload: %v", err)
			}
			if err := os.Remove(fullDir); err != nil {
				t.Fatal(err)
			}
			if err := os.Remove(filepath.Join(f.output, b.Slug, "index.html")); err != nil {
				t.Fatal(err)
			}
			recovered, c := f.build(t, "rotation-password-B")
			stats := recovered.ContentDiagnostics().TemplateCache
			if c.ArticleHTML != b.ArticleHTML || c.HTML != b.HTML || len(lifecycle.GetServeAffectedPaths(recovered)) != 0 ||
				stats == nil || stats.MissReasons.FullHTMLUnavailable != 1 || stats.RenderSucceeded != 1 {
				t.Fatalf("unchanged B wrapper did not re-render/publish missing full page: %+v", stats)
			}
			cipher := encryptedCiphertext(t, c.HTML)
			if _, err := encryption.Decrypt(cipher, "rotation-password-B"); err != nil {
				t.Fatalf("recovered page does not decrypt with B: %v", err)
			}
			if _, err := encryption.Decrypt(cipher, "rotation-password-A"); err == nil {
				t.Fatal("recovered page still accepts A")
			}
			stable, d := f.build(t, "rotation-password-B")
			if d.HTML != c.HTML || stable.ContentDiagnostics().TemplateCache.Restored != 1 {
				t.Fatal("post-recovery warm page did not stabilize")
			}
		})
	}
}

func TestEncryptionInvalidation_MissingWrapperPublishesRoot(t *testing.T) {
	f := newEncryptionPipelineFixture(t)
	f.source(t, "title: Root\nslug: ''\n", "rotation-password-A")
	f.build(t, "rotation-password-A") // Settle the unrelated cold nav reset.
	_, old := f.build(t, "rotation-password-A")
	if err := os.RemoveAll(filepath.Join(f.cache, buildcache.EncryptedHTMLCacheDir)); err != nil {
		t.Fatal(err)
	}
	m, fresh := f.build(t, "rotation-password-A")
	if !lifecycle.GetServeAffectedPaths(m)[fresh.Path] || fresh.HTML == old.HTML {
		t.Fatal("missing encrypted cache did not refresh existing root FullHTML/output")
	}
}

func TestEncryptionInvalidation_IdentityAndGenerationResult(t *testing.T) {
	cacheDir := t.TempDir()
	cache := buildcache.New(cacheDir)
	p := NewEncryptionPlugin()
	p.enabled, p.enforceStrength = true, false
	p.defaultKey, p.keys["one"], p.keys["two"] = "one", "test-key", "test-key"
	newPost := func(path string) *models.Post {
		post := headingMigrationPost(path, "secret", "plaintext")
		post.Private, post.ArticleHTML = true, "<p>Canonical article.</p>"
		return post
	}
	post := newPost("secret.md")
	changed, err := p.encryptPostWithCacheResult(post, cache)
	if err != nil || !changed {
		t.Fatalf("fresh encryption = %v, %v", changed, err)
	}
	oldWrapper := post.ArticleHTML
	post = newPost("secret.md")
	changed, err = p.encryptPostWithCacheResult(post, cache)
	if err != nil || changed || post.ArticleHTML != oldWrapper || post.Content != "" {
		t.Fatalf("valid hit = %v, %v (scrubbed=%v)", changed, err, post.Content == "")
	}
	for _, change := range []string{"key-name", "hint", "password", "path"} {
		t.Run(change, func(t *testing.T) {
			plugin := *p
			plugin.keys = map[string]string{"one": "test-key", "two": "test-key"}
			post := newPost("secret.md")
			switch change {
			case "key-name":
				post.SecretKey = "two"
			case "hint":
				plugin.decryptionHint = "Public hint"
			case "password":
				plugin.keys["one"] = "different-key"
			case "path":
				post.Path = "moved.md"
				// Even if an entry is copied across paths, IDs require a miss.
				copied := *cache.Posts["secret.md"]
				cache.Posts[post.Path] = &copied
			}
			changed, err := plugin.encryptPostWithCacheResult(post, cache)
			if err != nil || !changed || post.ArticleHTML == oldWrapper {
				t.Fatalf("%s change reused obsolete wrapper: %v, %v", change, changed, err)
			}
			if change == "path" && strings.Contains(post.ArticleHTML, fmt.Sprintf(`id="decrypt-input-%d"`, hashString("secret.md"))) {
				t.Fatal("moved path reused original accessibility IDs")
			}
		})
	}
	// Seed the pre-revision/path identity; it must miss, then hit current identity.
	legacyHash := buildcache.ContentHash("<p>Canonical article.</p>\x00one\x00test-key\x00")
	if err := cache.CacheEncryptedHTML("legacy.md", legacyHash, "obsolete-wrapper"); err != nil {
		t.Fatal(err)
	}
	post = newPost("legacy.md")
	if changed, err := p.encryptPostWithCacheResult(post, cache); err != nil || !changed {
		t.Fatalf("old format identity hit: %v, %v", changed, err)
	}
	post = newPost("legacy.md")
	if changed, err := p.encryptPostWithCacheResult(post, cache); err != nil || changed {
		t.Fatalf("current format identity missed: %v, %v", changed, err)
	}
	// Cache storage failure must not suppress real page invalidation.
	if err := cache.CacheFullHTML("uncached.md", "old full page"); err != nil {
		t.Fatal(err)
	}
	blocked := filepath.Join(cacheDir, buildcache.EncryptedHTMLCacheDir)
	if err := os.RemoveAll(blocked); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(blocked, []byte("blocks cache directory"), 0o600); err != nil {
		t.Fatal(err)
	}
	m := lifecycle.NewManager()
	m.Cache().Set("build_cache", cache)
	post = newPost("uncached.md")
	m.AddPost(post)
	if err := p.Render(m); err != nil {
		t.Fatal(err)
	}
	if !lifecycle.GetServeAffectedPaths(m)[post.Path] {
		t.Fatal("successful generation with failed cache write not invalidated")
	}
	if cache.GetCachedFullHTML(post.Path) != "" {
		t.Fatal("failed encrypted-wrapper write retained old full-page reference")
	}
	if err := cache.Save(); err != nil {
		t.Fatal(err)
	}
	reloaded, err := buildcache.Load(cacheDir)
	if err != nil || reloaded.GetCachedFullHTML(post.Path) != "" {
		t.Fatalf("failed wrapper write did not persist full-page invalidation: %v", err)
	}
	p.keys = nil
	post = newPost("missing-key.md")
	if changed, err := p.encryptPostWithCacheResult(post, cache); err == nil || changed {
		t.Fatalf("missing key bypassed: %v, %v", changed, err)
	}
	m.SetPosts([]*models.Post{newPost("secret.md")})
	lifecycle.SetServeAffectedPaths(m, map[string]bool{"preaffected.md": true})
	if err := p.Render(m); err == nil {
		t.Fatal("Render bypassed missing-key validation using a cached wrapper")
	}
	if affected := lifecycle.GetServeAffectedPaths(m); len(affected) != 1 || !affected["preaffected.md"] {
		t.Fatalf("failed validation changed existing affected paths: %v", affected)
	}
}

func TestEncryptionInvalidation_FullPageReferencePreservedOnHitOrError(t *testing.T) {
	cache := buildcache.New(t.TempDir())
	p := NewEncryptionPlugin()
	p.defaultKey, p.keys["test"] = "test", "test-key"
	makePost := func() *models.Post {
		post := headingMigrationPost("post.md", "post", "")
		post.Private, post.ArticleHTML = true, "<p>Canonical article.</p>"
		return post
	}
	if changed, err := p.encryptPostWithCacheResult(makePost(), cache); err != nil || !changed {
		t.Fatalf("initial wrapper generation = %v, %v", changed, err)
	}
	if err := cache.CacheFullHTML("post.md", "full page"); err != nil {
		t.Fatal(err)
	}
	oldPath := cache.Posts["post.md"].FullHTMLPath
	if changed, err := p.encryptPostWithCacheResult(makePost(), cache); err != nil || changed {
		t.Fatalf("wrapper hit = %v, %v", changed, err)
	}
	if cache.Posts["post.md"].FullHTMLPath != oldPath || cache.GetCachedFullHTML("post.md") != "full page" {
		t.Fatal("wrapper hit invalidated valid full page")
	}
	p.keys = nil
	if changed, err := p.encryptPostWithCacheResult(makePost(), cache); err == nil || changed {
		t.Fatalf("missing-key result = %v, %v", changed, err)
	}
	if cache.Posts["post.md"].FullHTMLPath != oldPath || cache.GetCachedFullHTML("post.md") != "full page" {
		t.Fatal("encryption error invalidated existing full page")
	}
}

func TestEncryptionInvalidation_WorkerErrorStillMergesSuccessfulSiblings(t *testing.T) {
	cache := buildcache.New(t.TempDir())
	m := lifecycle.NewManager()
	m.SetConcurrency(4)
	m.Cache().Set("build_cache", cache)
	p := NewEncryptionPlugin()
	p.enabled, p.enforceStrength = true, false
	p.keys["good"], p.keys["bad"] = "test-key", ""
	good := headingMigrationPost("good.md", "good", "")
	good.Private, good.SecretKey, good.ArticleHTML = true, "good", "<p>Private.</p>"
	bad := headingMigrationPost("bad.md", "bad", "")
	bad.Private, bad.SecretKey, bad.ArticleHTML = true, "bad", "<p>Private.</p>"
	m.SetPosts([]*models.Post{bad, good})
	cache.Graph.SetDependencies("direct.md", "direct", []string{"good"})
	cache.Graph.SetDependencies("transitive.md", "transitive", []string{"direct"})
	lifecycle.SetServeAffectedPaths(m, map[string]bool{"preaffected.md": true})
	// Fault at the encryption boundary without replacing global randomness:
	// direct test key injection supplies an empty password, which reaches
	// Encrypt when strength enforcement is disabled and returns its usual error.
	err := p.Render(m)
	if !errors.Is(err, encryption.ErrEmptyPassword) {
		t.Fatalf("worker error not preserved: %v", err)
	}
	affected := lifecycle.GetServeAffectedPaths(m)
	for _, path := range []string{"good.md", "direct.md", "transitive.md", "preaffected.md"} {
		if !affected[path] {
			t.Errorf("successful sibling closure missing %s: %v", path, affected)
		}
	}
	if affected[bad.Path] {
		t.Fatal("failed encryption reported regenerated")
	}
	for _, slug := range cache.GetChangedSlugs() {
		if slug == bad.Slug {
			t.Fatal("failed encryption marked its slug")
		}
	}
}

func TestEncryptionInvalidation_ConcurrentClosurePreservesAffected(t *testing.T) {
	cache := buildcache.New(t.TempDir())
	m := lifecycle.NewManager()
	m.SetConcurrency(4)
	m.Cache().Set("build_cache", cache)
	p := NewEncryptionPlugin()
	p.enabled, p.enforceStrength = true, false
	p.defaultKey, p.keys["test"] = "test", "test-key"
	for _, path := range []string{"one.md", "two.md", "three.md", "root.md"} {
		post := headingMigrationPost(path, strings.TrimSuffix(path, ".md"), "")
		post.Private, post.ArticleHTML = true, "<p>Private.</p>"
		if path == "root.md" {
			post.Slug = ""
			post.Set("_slug_explicit", true)
		}

		m.AddPost(post)
	}
	cache.Graph.SetDependencies("direct.md", "direct", []string{"one"})
	cache.Graph.SetDependencies("transitive.md", "transitive", []string{"direct"})
	lifecycle.SetServeAffectedPaths(m, map[string]bool{"preaffected.md": true})
	if err := p.Render(m); err != nil {
		t.Fatal(err)
	}
	affected := lifecycle.GetServeAffectedPaths(m)
	for _, path := range []string{"one.md", "two.md", "three.md", "root.md", "direct.md", "transitive.md", "preaffected.md"} {
		if !affected[path] {
			t.Errorf("affected closure missing %s: %v", path, affected)
		}
	}
	cache.ResetStats()
	lifecycle.SetServeAffectedPaths(m, map[string]bool{"preaffected.md": true})
	for _, post := range m.Posts() {
		post.ArticleHTML = "<p>Private.</p>"
	}
	if err := p.Render(m); err != nil {
		t.Fatal(err)
	}
	if affected := lifecycle.GetServeAffectedPaths(m); len(affected) != 1 || !affected["preaffected.md"] ||
		len(cache.GetChangedSlugs()) != 0 {
		t.Fatalf("valid wrapper hits added affected/dependency marks: %v, %v", affected, cache.GetChangedSlugs())
	}
}
