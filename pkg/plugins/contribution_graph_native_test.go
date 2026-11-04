package plugins

import (
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/models"
)

func TestNativeContributionGraph_DefaultNoBrowserAssets(t *testing.T) {
	p := NewContributionGraphPlugin()
	post := &models.Post{ArticleHTML: `<pre><code class="language-contribution-graph">{"data":[{"date":"2024-01-01","value":5},{"date":"2024-02-29","value":2}],"options":{"year":2024}}</code></pre>`}
	if err := p.processPost(nil, post); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(post.ArticleHTML, "<svg") || !strings.Contains(post.ArticleHTML, "Feb 29, 2024") {
		t.Fatalf("expected native year SVG: %s", post.ArticleHTML)
	}
	if !strings.Contains(post.ArticleHTML, "var(--color-primary") {
		t.Fatal("native SVG must honor the active CSS palette")
	}
	if strings.Contains(post.ArticleHTML, "<script") || strings.Contains(post.ArticleHTML, "cal-heatmap") {
		t.Fatal("native calendar unexpectedly shipped browser rendering dependencies")
	}
	if !strings.Contains(post.ArticleHTML, "5 posts on Jan 1, 2024") {
		t.Fatal("calendar cells must include accessible value/date tooltips")
	}
}

func TestNativeContributionGraph_StrictNativeRejectsUnsupportedOptions(t *testing.T) {
	p := NewContributionGraphPlugin()
	p.config.Mode = "go"
	post := &models.Post{ArticleHTML: `<pre><code class="language-contribution-graph">{"data":[],"options":{"domain":"month","range":12}}</code></pre>`}
	if err := p.processPost(nil, post); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(post.ArticleHTML, "contribution-graph-error") ||
		!strings.Contains(post.ArticleHTML, "mode = &#34;client&#34;") {
		t.Fatalf("unsupported native options should give an actionable error: %s", post.ArticleHTML)
	}
	if strings.Contains(post.ArticleHTML, "<script") {
		t.Fatal("unsupported native options must not silently load Cal-Heatmap")
	}
}

func TestNativeContributionGraph_ExplicitClientCompatibility(t *testing.T) {
	p := NewContributionGraphPlugin()
	p.config.Mode = "client"
	post := &models.Post{ArticleHTML: `<pre><code class="language-contribution-graph">{"data":[{"date":"2024-01-01","value":5}],"options":{"domain":"month","range":12}}</code></pre>`}
	if err := p.processPost(nil, post); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(post.ArticleHTML, "new CalHeatmap()") {
		t.Fatal("client mode must preserve Cal-Heatmap")
	}
}

func TestNativeContributionGraph_AutoPreservesDocumentedLegacyLayouts(t *testing.T) {
	for _, options := range []string{
		`{"year":2024,"domain":"month","subDomain":"day","range":3}`,
		`{"year":2024,"domain":"year","subDomain":"day","range":3}`,
		`{"year":2024,"domain":"year","subDomain":"day","cellSize":12}`,
	} {
		t.Run(options, func(t *testing.T) {
			p := NewContributionGraphPlugin()
			src := `{"data":[{"date":"2024-01-01","value":1}],"options":` + options + `}`
			post := &models.Post{ArticleHTML: `<pre><code class="language-contribution-graph">` + src + `</code></pre>`}
			if err := p.processPost(nil, post); err != nil { t.Fatal(err) }
			if !strings.Contains(post.ArticleHTML, "new CalHeatmap()") || strings.Contains(post.ArticleHTML, "contribution-graph-error") {
				t.Fatalf("legacy config should use Cal-Heatmap in auto mode: %s", post.ArticleHTML)
			}
		})
	}
}
func TestNativeContributionGraph_AutoAcceptsYearRangeOne(t *testing.T) {
	p := NewContributionGraphPlugin()
	post := &models.Post{ArticleHTML: `<pre><code class="language-contribution-graph">{"data":[{"date":"2024-01-01","value":1}],"options":{"year":2024,"domain":"year","range":1}}</code></pre>`}
	if err := p.processPost(nil, post); err != nil { t.Fatal(err) }
	if !strings.Contains(post.ArticleHTML, "<svg") || strings.Contains(post.ArticleHTML, "new CalHeatmap()") {
		t.Fatalf("default year range one should be browser-free: %s", post.ArticleHTML)
	}
}
