package plugins

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/buildcache"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

func headingMigrationManager(t *testing.T, cache *buildcache.Cache, posts ...*models.Post) (*lifecycle.Manager, *TemplatesPlugin) {
	t.Helper()
	m := lifecycle.NewManager()
	m.Cache().Set("build_cache", cache)
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "post.html"), []byte(`PAGE:{{ body | safe }}`), 0o600); err != nil {
		t.Fatal(err)
	}
	m.Config().Extra["templates_dir"] = dir
	for _, post := range posts {
		m.AddPost(post)
	}
	p := NewTemplatesPlugin()
	if err := p.Configure(m); err != nil {
		t.Fatal(err)
	}
	return m, p
}

func headingMigrationPost(path, slug, content string) *models.Post {
	post := models.NewPost(path)
	post.Slug = slug
	post.Href = "/" + slug + "/"
	post.Content = content
	post.InputHash = buildcache.ContentHash(path + content)
	post.Template = "post.html"
	post.Published = true
	return post
}

func renderHeadingMigrationMarkdown(t *testing.T, m *lifecycle.Manager) {
	t.Helper()
	p := NewRenderMarkdownPlugin()
	if err := p.Configure(m); err != nil {
		t.Fatal(err)
	}
	if err := p.Render(m); err != nil {
		t.Fatal(err)
	}
}

func seedHeadingMigrationPages(t *testing.T, cache *buildcache.Cache, posts []*models.Post) {
	t.Helper()
	for _, post := range posts {
		cache.MarkRebuiltWithSlug(post.Path, post.Slug, post.InputHash, "", post.Template)
		cache.SetHeadingHighlightRevision(post.Path, 0)
		if err := cache.CacheFullHTML(post.Path, "OLD:"+post.Path); err != nil {
			t.Fatal(err)
		}
		post.HTML = ""
	}
	cache.ResetStats()
}

func TestBuildCachePlugin_HeadingMigrationTargetsPagesDependenciesAndFeeds(t *testing.T) {
	cacheDir := t.TempDir()
	cache := buildcache.New(cacheDir)
	root := headingMigrationPost("root.md", "", "# ==Kirby==\n")
	marked := headingMigrationPost("marked.md", "marked", "## ==Kirby==\n")
	marked.Set("aliases", []string{"Heading Alias"})
	dependent := headingMigrationPost("dependent.md", "dependent", "Dependent.")
	dependent.Dependencies = []string{"heading-alias"}
	prose := headingMigrationPost("prose.md", "prose", "Prose ==ordinary==.")
	h3 := headingMigrationPost("h3.md", "h3", "### ==third==\n")
	title := headingMigrationPost("title.md", "title", "Ordinary body.")
	title.TitleHTML = `<span class="heading-highlight"><mark>Title</mark></span>`
	unrelated := headingMigrationPost("unrelated.md", "unrelated", "Unrelated.")
	definition := headingMigrationPost("definition.md", "kirby", "Definition.")
	definition.Title = new("Kirby")
	definition.Set("templateKey", "glossary")
	posts := []*models.Post{root, marked, dependent, prose, h3, title, unrelated, definition}
	m, page := headingMigrationManager(t, cache, posts...)
	renderHeadingMigrationMarkdown(t, m)
	if err := page.Render(m); err != nil {
		t.Fatal(err)
	}
	canonical := root.ArticleHTML
	seedHeadingMigrationPages(t, cache, posts)
	cache.Graph.SetDependencies(dependent.Path, dependent.Slug, dependent.Dependencies)
	glossary := NewGlossaryPlugin()
	if err := glossary.Configure(m); err != nil {
		t.Fatal(err)
	}
	if err := glossary.buildGlossary(posts); err != nil {
		t.Fatal(err)
	}
	oldDerivative := wrapHeadingMarkHighlights(canonical)
	cache.CacheGlossaryHTML(root.Path, buildcache.ContentHash(oldDerivative+glossary.computeTermsHash()), oldDerivative)
	oldPage := "PAGE:" + oldDerivative
	if err := cache.CacheFullHTML(root.Path, oldPage); err != nil {
		t.Fatal(err)
	}
	if strings.Count(oldPage, `class="heading-highlight"`) != 2 {
		t.Fatal("fixture must contain duplicated legacy heading wrappers")
	}
	lifecycle.SetServeAffectedPaths(m, map[string]bool{"existing.md": true})
	renderHeadingMigrationMarkdown(t, m)
	affected := lifecycle.GetServeAffectedPaths(m)
	wantAffected := map[string]bool{"root.md": true, "marked.md": true, "dependent.md": true, "existing.md": true}
	if !reflect.DeepEqual(affected, wantAffected) {
		t.Fatalf("affected paths = %v, want %v", affected, wantAffected)
	}
	if got := cache.GetChangedFeedSlugs(); !reflect.DeepEqual(got, []string{"marked"}) {
		t.Fatalf("changed feed identities = %v", got)
	}
	changed := getChangedSlugsMap(cache)
	for _, identity := range []string{"marked", "Heading Alias", "heading alias", "heading-alias", "dependent"} {
		if !changed[identity] {
			t.Fatalf("migration did not propagate dependency identity %q: %v", identity, changed)
		}
	}
	if root.ArticleHTML != canonical || cache.GetHeadingHighlightRevision(root.Path) != 0 {
		t.Fatal("canonical restore mutated article or prematurely certified page")
	}
	expectedArticle := glossary.linkTerms(canonical, root)
	if err := glossary.Render(m); err != nil {
		t.Fatal(err)
	}
	if root.ArticleHTML != expectedArticle || strings.Count(root.ArticleHTML, `class="heading-highlight"`) != 1 {
		t.Fatalf("legacy glossary derivative survived: %s", root.ArticleHTML)
	}

	publisher := NewPublishFeedsPlugin()
	for _, tc := range []struct {
		slug string
		post *models.Post
		skip bool
	}{{"root-member", root, false}, {"marked-member", marked, false}, {"dependent-member", dependent, false}, {"control", unrelated, true}} {
		feed := &models.FeedConfig{Slug: tc.slug, Posts: []*models.Post{tc.post}}
		hash := publisher.computeFeedHashWithConfigAndCache(feed, nil, cache)
		cache.SetFeedHash(feed.Slug, hash)
		if skip, gotHash := publisher.shouldSkipFeedWithConfigAndChanges(feed, cache, t.TempDir(), nil, affected); skip != tc.skip || gotHash != hash {
			t.Fatalf("feed %s: skip=%v hash=%s, want skip=%v unchanged hash=%s", tc.slug, skip, gotHash, tc.skip, hash)
		}
	}
	if err := page.Render(m); err != nil {
		t.Fatal(err)
	}
	for _, post := range posts {
		if affected[post.Path] {
			expectedPage := "PAGE:" + annotateLocalPreviewLinks(post.ArticleHTML, post.Href, page.siteURL, page.localPreviews)
			if post.HTML != expectedPage || cache.GetHeadingHighlightRevision(post.Path) != headingHighlightRevision {
				t.Fatalf("affected page not freshly cached and certified: %s: %s", post.Path, post.HTML)
			}
		} else if post.HTML != "OLD:"+post.Path || cache.GetHeadingHighlightRevision(post.Path) != 0 {
			t.Fatalf("unrelated page lost its cache hit: %s: %s", post.Path, post.HTML)
		}
	}
	if err := cache.Save(); err != nil {
		t.Fatal(err)
	}
	reloaded, err := buildcache.Load(cacheDir)
	if err != nil {
		t.Fatal(err)
	}
	if reloaded.GetHeadingHighlightRevision(root.Path) != headingHighlightRevision {
		t.Fatal("page revision did not survive disk reload")
	}
	m.Cache().Set("build_cache", reloaded)
	lifecycle.SetServeAffectedPaths(m, nil)
	renderHeadingMigrationMarkdown(t, m)
	if len(lifecycle.GetServeAffectedPaths(m)) != 0 {
		t.Fatal("completed disk-loaded migration was enqueued again")
	}
}

func TestBuildCachePlugin_HeadingMigrationPreservesConfigIdentity(t *testing.T) {
	dir := t.TempDir()
	m := lifecycle.NewManager()
	m.Config().ContentDir = dir
	m.Config().GlobPatterns = []string{"posts/**/*.md"}
	m.Config().Extra["config_overlay"] = "preview"
	want := dir + "\nposts/**/*.md\noverlay:preview"
	if got := configHashInput(m.Config(), nil); got != want {
		t.Fatalf("configuration identity changed: %q, want %q", got, want)
	}
}
