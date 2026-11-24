package plugins

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/buildcache"
	"github.com/WaylonWalker/markata/pkg/diagnostics"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
)

type failingHeadingMigrationMarkdown struct {
	goldmark.Markdown
	err error
}

func (p failingHeadingMigrationMarkdown) Convert(_ []byte, _ io.Writer, _ ...parser.ParseOption) error {
	return p.err
}

func TestHeadingHighlightMigration_ClassificationErrorIsReported(t *testing.T) {
	cache := buildcache.New(t.TempDir())
	post := headingMigrationPost("marked.md", "marked", "# ==Heading==")
	m, _ := headingMigrationManager(t, cache, post)
	lifecycle.SetServeIncremental(m, true)
	lifecycle.SetServeAffectedPaths(m, map[string]bool{"other.md": true})
	p := NewRenderMarkdownPlugin()
	if err := p.Configure(m); err != nil {
		t.Fatal(err)
	}
	expectedError := errors.New("classification conversion failed")
	p.md = failingHeadingMigrationMarkdown{Markdown: p.md, err: expectedError}
	if err := p.Render(m); !errors.Is(err, expectedError) {
		t.Fatalf("classification error = %v, want %v", err, expectedError)
	}
	snapshot := m.ContentLedger().Snapshot()
	if len(snapshot.Entries) != 1 || snapshot.Entries[0].Path != post.Path ||
		len(snapshot.Entries[0].Diagnostics) != 1 ||
		snapshot.Entries[0].Diagnostics[0].Code != diagnostics.ReasonContentRenderError {
		t.Fatalf("classification error missing from content ledger: %+v", snapshot)
	}
	if cache.GetHeadingHighlightRevision(post.Path) != 0 || post.ArticleHTML != "" || post.HTML != "" {
		t.Fatal("failed classification changed or certified output")
	}
}

func TestHeadingHighlightMigration_Classification(t *testing.T) {
	for _, tc := range []struct {
		html string
		want bool
	}{
		{`<h1><span class="heading-highlight"><mark>x</mark></span></h1>`, true},
		{`<H2 class="test"><MARK data-x="1">x</MARK></H2>`, true},
		{`<p><mark>x</mark></p>`, false},
		{`<h3><mark>x</mark></h3><h6><mark>x</mark></h6>`, false},
		{`<h1>&lt;mark&gt;x&lt;/mark&gt;</h1>`, false},
		{`<!-- <h1><mark>x</mark></h1> -->`, false},
		{`<h2><!-- <mark>x</mark> -->plain</h2>`, false},
		{`<h1><marker>x</marker></h1>`, false},
	} {
		if got := hasHeadingMarkHighlights(tc.html); got != tc.want {
			t.Errorf("classification(%q) = %v, want %v", tc.html, got, tc.want)
		}
	}
}

func TestHeadingHighlightMigration_OrdinaryLargeArticleClassificationAllocations(t *testing.T) {
	article := strings.Repeat("<h2>Ordinary HEADING</h2>\n<p>Mixed CASE prose with no highlighted text.</p>\n", 16384)
	if allocations := testing.AllocsPerRun(100, func() {
		if hasHeadingMarkHighlights(article) {
			t.Fatal("ordinary article classified as migration candidate")
		}
	}); allocations != 0 {
		t.Fatalf("ordinary large article classification allocated: %.0f allocations/run", allocations)
	}
}

func BenchmarkHeadingHighlightMigration_OrdinaryLargeArticleClassification(b *testing.B) {
	article := strings.Repeat("<h2>Ordinary HEADING</h2>\n<p>Mixed CASE prose with no highlighted text.</p>\n", 16384)
	if hasHeadingMarkHighlights(article) {
		b.Fatal("ordinary article classified as migration candidate")
	}
	b.ReportAllocs()
	b.SetBytes(int64(len(article)))
	b.ResetTimer()
	for b.Loop() {
		if hasHeadingMarkHighlights(article) {
			b.Fatal("ordinary article classified as migration candidate")
		}
	}
}

func TestHeadingHighlightMigration_CanonicalArticlePipeline(t *testing.T) {
	p := NewRenderMarkdownPlugin()
	const content = "# ==Heading==\n\n![Example](example.png)\n\n> Figure caption.\n\n> Quotation.\n-- Author\n"
	raw, err := p.doRender(content)
	if err != nil {
		t.Fatal(err)
	}
	want := wrapHeadingMarkHighlights(mergeBlockquoteAttributions(mergeFigureBlockquoteCaptions(raw)))
	if want == raw || !strings.Contains(want, "<figcaption>") ||
		!strings.Contains(want, `class="heading-highlight"`) ||
		!strings.Contains(want, "<footer>") {
		t.Fatalf("fixture does not exercise every canonical post-processor: %s", want)
	}
	canonical, err := p.renderArticleHTML(content)
	if err != nil {
		t.Fatal(err)
	}
	if canonical != want {
		t.Fatalf("canonical probe pipeline differs:\ngot: %s\nwant: %s", canonical, want)
	}
	post := headingMigrationPost("pipeline.md", "pipeline", content)
	if err := p.renderPost(post); err != nil {
		t.Fatal(err)
	}
	if post.ArticleHTML != canonical {
		t.Fatalf("fresh article differs from canonical probe:\ngot: %s\nwant: %s", post.ArticleHTML, canonical)
	}
}

func TestHeadingHighlightMigration_NavResetStampsReplacementMetadata(t *testing.T) {
	cache := buildcache.New(t.TempDir())
	marked := headingMigrationPost("marked.md", "marked", "# ==Heading==")
	control := headingMigrationPost("control.md", "control", "Plain.")
	m, page := headingMigrationManager(t, cache, marked, control)
	renderHeadingMigrationMarkdown(t, m)
	if err := page.Render(m); err != nil {
		t.Fatal(err)
	}
	seedHeadingMigrationPages(t, cache, m.Posts())
	previous := cache.Posts[marked.Path]
	// Canonical articles restore first; Templates then invalidates metadata.
	lifecycle.SetServeAffectedPaths(m, nil)
	renderHeadingMigrationMarkdown(t, m)
	cache.NavPreviewHash = "obsolete-nav"
	if err := page.Render(m); err != nil {
		t.Fatal(err)
	}
	for _, post := range m.Posts() {
		if cache.Posts[post.Path] == previous || previous.HeadingHighlightRevision != 0 ||
			cache.GetHeadingHighlightRevision(post.Path) != headingHighlightRevision {
			t.Fatalf("fresh page did not stamp replacement metadata: %s", post.Path)
		}
		cache.MarkRebuiltWithSlug(post.Path, post.Slug, post.InputHash, "", post.Template)
	}
	cache.ResetStats()
	lifecycle.SetServeAffectedPaths(m, nil)
	renderHeadingMigrationMarkdown(t, m)
	if len(lifecycle.GetServeAffectedPaths(m)) != 0 {
		t.Fatal("nav-reset replacement metadata was migrated again")
	}
	if err := page.Render(m); err != nil {
		t.Fatal(err)
	}
	if marked.HTML != "PAGE:"+marked.ArticleHTML || control.HTML != "PAGE:"+control.ArticleHTML {
		t.Fatal("second warm render changed freshly cached pages")
	}
}

func TestHeadingHighlightMigration_MissingCachesAndIncrementalSelection(t *testing.T) {
	const fullPageOnly = "full-page"
	for _, missing := range []string{"article", fullPageOnly, "both"} {
		t.Run(missing, func(t *testing.T) {
			cache := buildcache.New(t.TempDir())
			marked := headingMigrationPost("marked.md", "marked", "# ==Heading==")
			control := headingMigrationPost("control.md", "control", "Plain.")
			m, page := headingMigrationManager(t, cache, marked, control)
			renderHeadingMigrationMarkdown(t, m)
			if err := page.Render(m); err != nil {
				t.Fatal(err)
			}
			seedHeadingMigrationPages(t, cache, m.Posts())
			if missing != fullPageOnly {
				cache.Posts[marked.Path].ArticleHTMLPath = ""
				cache.Posts[control.Path].ArticleHTMLPath = ""
			}
			if missing != "article" {
				cache.Posts[marked.Path].FullHTMLPath = ""
				cache.Posts[control.Path].FullHTMLPath = ""
			}
			marked.ArticleHTML, control.ArticleHTML = "", ""
			lifecycle.SetServeIncremental(m, true)
			lifecycle.SetServeAffectedPaths(m, map[string]bool{"selected.md": true})
			renderHeadingMigrationMarkdown(t, m)
			affected := lifecycle.GetServeAffectedPaths(m)
			if !affected[marked.Path] || affected[control.Path] {
				t.Fatalf("missing-cache classification expanded unrelated selection: %v", affected)
			}
			if marked.ArticleHTML != "<h1 id=\"heading\"><span class=\"heading-highlight\"><mark style=\"background-color:var(--color-highlight)!important;color:var(--color-highlight-text)!important\">Heading</mark></span></h1>\n" {
				t.Fatalf("missing-cache heading not canonically rendered: %q", marked.ArticleHTML)
			}
			if missing != fullPageOnly && control.ArticleHTML != "" {
				t.Fatal("omitted unrelated missing article was rendered")
			}
			if err := page.Render(m); err != nil {
				t.Fatal(err)
			}
			if marked.HTML != "PAGE:"+marked.ArticleHTML || cache.GetHeadingHighlightRevision(marked.Path) != headingHighlightRevision {
				t.Fatal("missing-cache migration was not refreshed and certified")
			}
			if missing != fullPageOnly && cache.GetHeadingHighlightRevision(control.Path) != 0 {
				t.Fatal("omitted unrelated article was falsely certified")
			}
			if missing == "both" && control.HTML != "" {
				t.Fatalf("omitted unrelated article produced a fresh blank page: %q", control.HTML)
			}
			if missing == "article" && control.HTML != "OLD:"+control.Path {
				t.Fatal("omitted unrelated article lost its usable full-page hit")
			}
		})
	}
}

func TestHeadingHighlightMigration_IncrementalMissingDependentArticlesSelectedBeforeWorkers(t *testing.T) {
	for _, sourceMissing := range []bool{false, true} {
		name := "source-hit"
		if sourceMissing {
			name = "source-missing"
		}
		t.Run(name, func(t *testing.T) {
			cache := buildcache.New(t.TempDir())
			source := headingMigrationPost("source.md", "source", "# ==Heading==")
			dependent := headingMigrationPost("dependent.md", "dependent", "Dependent article.")
			transitive := headingMigrationPost("transitive.md", "transitive", "Transitive article.")
			control := headingMigrationPost("control.md", "control", "Unrelated.")
			dependent.Dependencies = []string{"source"}
			transitive.Dependencies = []string{"dependent"}
			m, page := headingMigrationManager(t, cache, source, dependent, transitive, control)
			renderHeadingMigrationMarkdown(t, m)
			if err := page.Render(m); err != nil {
				t.Fatal(err)
			}
			seedHeadingMigrationPages(t, cache, m.Posts())
			cache.Graph.SetDependencies(dependent.Path, dependent.Slug, dependent.Dependencies)
			cache.Graph.SetDependencies(transitive.Path, transitive.Slug, transitive.Dependencies)
			for _, post := range []*models.Post{dependent, transitive, control} {
				// Already certified, nonheading dependents cannot enqueue
				// themselves through legacy missing-article classification.
				cache.SetHeadingHighlightRevision(post.Path, headingHighlightRevision)
				cache.Posts[post.Path].ArticleHTMLPath = ""
				post.ArticleHTML = ""
			}
			if sourceMissing {
				cache.Posts[source.Path].ArticleHTMLPath = ""
			}
			source.ArticleHTML = ""
			lifecycle.SetServeIncremental(m, true)
			lifecycle.SetServeAffectedPaths(m, map[string]bool{"selected.md": true})
			renderHeadingMigrationMarkdown(t, m)
			affected := lifecycle.GetServeAffectedPaths(m)
			if len(affected) != 4 || !affected["selected.md"] ||
				!affected[source.Path] || !affected[dependent.Path] || !affected[transitive.Path] {
				t.Fatalf("dependency closure was not selected: %v", affected)
			}
			for _, tc := range []struct {
				post *models.Post
				want string
			}{{dependent, "<p>Dependent article.</p>\n"}, {transitive, "<p>Transitive article.</p>\n"}} {
				if tc.post.ArticleHTML != tc.want {
					t.Fatalf("affected dependent missed canonical rendering: %s: %q, want %q", tc.post.Path, tc.post.ArticleHTML, tc.want)
				}
			}
			if control.ArticleHTML != "" {
				t.Fatal("unrelated missing article bypassed incremental selection")
			}
			if err := page.Render(m); err != nil {
				t.Fatal(err)
			}
			for _, post := range []*models.Post{source, dependent, transitive} {
				if post.HTML != "PAGE:"+post.ArticleHTML {
					t.Fatalf("affected page restored obsolete full HTML: %s: %q", post.Path, post.HTML)
				}
			}
			if control.HTML != "OLD:"+control.Path {
				t.Fatal("unrelated full-page cache hit was lost")
			}
		})
	}
}

func TestHeadingHighlightMigration_FailuresDoNotStamp(t *testing.T) {
	for _, failure := range []string{"template", "full-cache"} {
		t.Run(failure, func(t *testing.T) {
			cacheDir := t.TempDir()
			cache := buildcache.New(cacheDir)
			post := headingMigrationPost("marked.md", "marked", "# ==Heading==")
			m, page := headingMigrationManager(t, cache, post)
			renderHeadingMigrationMarkdown(t, m)
			// Stabilize nav state before installing the failure.
			if err := page.Render(m); err != nil {
				t.Fatal(err)
			}
			seedHeadingMigrationPages(t, cache, m.Posts())
			renderHeadingMigrationMarkdown(t, m)
			if failure == "template" {
				dir := t.TempDir()
				if err := os.WriteFile(filepath.Join(dir, "post.html"), []byte(`{% invalid_tag %}`), 0o600); err != nil {
					t.Fatal(err)
				}
				m.Config().Extra["templates_dir"] = dir
				if err := page.Configure(m); err != nil {
					t.Fatal(err)
				}
			} else {
				dir := filepath.Join(cacheDir, buildcache.FullHTMLCacheDir)
				if err := os.RemoveAll(dir); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(dir, []byte("blocks directory"), 0o600); err != nil {
					t.Fatal(err)
				}
			}
			err := page.Render(m)
			if failure == "template" && err == nil {
				t.Fatal("template failure was hidden")
			}
			if failure == "full-cache" && err != nil {
				t.Fatalf("best-effort cache failure became fatal: %v", err)
			}
			if cache.GetHeadingHighlightRevision(post.Path) != 0 {
				t.Fatal("failed page/cache write certified migration")
			}
			if got := cache.GetCachedFullHTML(post.Path); got != "OLD:"+post.Path {
				t.Fatalf("failed write replaced old cached page: %q", got)
			}
		})
	}
}

func TestHeadingHighlightMigration_SkippedEncryptedAndSyntheticUntouched(t *testing.T) {
	cache := buildcache.New(t.TempDir())
	skipped := headingMigrationPost("skip.md", "skip", "# ==Heading==")
	skipped.Skip = true
	encrypted := headingMigrationPost("encrypted.md", "encrypted", "# ==Heading==")
	encrypted.Set(sourceEncryptedPostKey, true)
	synthetic := headingMigrationPost("", "synthetic", "# ==Heading==")
	posts := []*models.Post{skipped, encrypted, synthetic}
	m, page := headingMigrationManager(t, cache, posts...)
	renderHeadingMigrationMarkdown(t, m)
	if len(lifecycle.GetServeAffectedPaths(m)) != 0 {
		t.Fatal("ineligible source enqueued for migration")
	}
	if cache.GetCachedArticleHTML(encrypted.Path, buildcache.ContentHash(encrypted.Content)) != "" {
		t.Fatal("plaintext source-encrypted article was cached")
	}
	if err := page.Render(m); err != nil {
		t.Fatal(err)
	}
	for _, post := range posts {
		if cache.GetHeadingHighlightRevision(post.Path) != 0 {
			t.Fatalf("ineligible source was certified: %q", post.Path)
		}
	}
	if skipped.ArticleHTML != "" || skipped.HTML != "" {
		t.Fatal("skipped post was rendered")
	}
	// An unrelated incremental edit must not probe these sources for migration.
	// A failing converter makes any accidental plaintext probe observable.
	lifecycle.SetServeIncremental(m, true)
	lifecycle.SetServeAffectedPaths(m, map[string]bool{"other.md": true})
	p := NewRenderMarkdownPlugin()
	if err := p.Configure(m); err != nil {
		t.Fatal(err)
	}
	p.md = failingHeadingMigrationMarkdown{Markdown: p.md, err: errors.New("ineligible source was probed")}
	if err := p.Render(m); err != nil {
		t.Fatal(err)
	}
	affected := lifecycle.GetServeAffectedPaths(m)
	if len(affected) != 1 || !affected["other.md"] {
		t.Fatalf("ineligible source expanded incremental migration: %v", affected)
	}
}
