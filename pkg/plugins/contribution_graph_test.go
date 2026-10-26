package plugins

import (
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/models"
)

func TestContributionGraphPlugin_Name(t *testing.T) {
	p := NewContributionGraphPlugin()
	if got := p.Name(); got != "contribution_graph" {
		t.Errorf("Name() = %q, want %q", got, "contribution_graph")
	}
}

func TestContributionGraphPlugin_ProcessPost_NoContributionGraph(t *testing.T) {
	p := NewContributionGraphPlugin()

	post := &models.Post{ArticleHTML: "<p>Hello world</p>"}

	err := p.processPost(nil, post)
	if err != nil {
		t.Errorf("processPost() error = %v", err)
	}

	if post.ArticleHTML != "<p>Hello world</p>" {
		t.Errorf("HTML was modified when no contribution-graph blocks present")
	}
}

func TestContributionGraphPlugin_ProcessPost_ValidGraph(t *testing.T) {
	p := NewContributionGraphPlugin()

	post := &models.Post{ArticleHTML: `<p>Check out this contribution graph:</p>
<pre><code class="language-contribution-graph">{
  "data": [{"date": "2024-01-01", "value": 5}],
  "options": {"domain": "year", "subDomain": "day"}
}</code></pre>
<p>Nice graph!</p>`}

	err := p.processPost(nil, post)
	if err != nil {
		t.Errorf("processPost() error = %v", err)
	}

	if !strings.Contains(post.ArticleHTML, "<div id=") {
		t.Error("Expected div element in output")
	}
	if !strings.Contains(post.ArticleHTML, "cal-heatmap") {
		t.Error("Expected Cal-Heatmap CDN script")
	}
	if !strings.Contains(post.ArticleHTML, `class="contribution-graph-container"`) {
		t.Error("Expected contribution-graph-container class")
	}
	if !strings.Contains(post.ArticleHTML, "new CalHeatmap()") {
		t.Error("Expected CalHeatmap initialization script")
	}
	if !strings.Contains(post.ArticleHTML, "window._contributionGraphFitters") {
		t.Error("Expected built-in graph fitters")
	}
}

func TestContributionGraphPlugin_ProcessPost_InvalidJSON(t *testing.T) {
	p := NewContributionGraphPlugin()

	post := &models.Post{ArticleHTML: `<pre><code class="language-contribution-graph">{ invalid json }</code></pre>`}

	err := p.processPost(nil, post)
	if err != nil {
		t.Errorf("processPost() error = %v", err)
	}

	if !strings.Contains(post.ArticleHTML, "contribution-graph-error") {
		t.Error("Expected error class for invalid JSON")
	}
	if !strings.Contains(post.ArticleHTML, "Invalid JSON") {
		t.Error("Expected error message for invalid JSON")
	}
}

func TestContributionGraphPlugin_ProcessPost_MultipleGraphs(t *testing.T) {
	p := NewContributionGraphPlugin()

	post := &models.Post{ArticleHTML: `<pre><code class="language-contribution-graph">{"data": [], "options": {}}</code></pre>
<p>And another:</p>
<pre><code class="language-contribution-graph">{"data": [], "options": {"domain": "month"}}</code></pre>`}

	err := p.processPost(nil, post)
	if err != nil {
		t.Errorf("processPost() error = %v", err)
	}

	count := strings.Count(post.ArticleHTML, "<div id=")
	if count != 2 {
		t.Errorf("Expected 2 div elements, got %d", count)
	}

	cdnCount := strings.Count(post.ArticleHTML, "cal-heatmap.min.js")
	if cdnCount != 1 {
		t.Errorf("Expected 1 CDN script inclusion, got %d", cdnCount)
	}
}

func TestContributionGraphPlugin_ProcessPost_SkipPost(t *testing.T) {
	p := NewContributionGraphPlugin()

	post := &models.Post{
		Skip:        true,
		ArticleHTML: `<pre><code class="language-contribution-graph">{"data": []}</code></pre>`,
	}
	originalHTML := post.ArticleHTML

	err := p.processPost(nil, post)
	if err != nil {
		t.Errorf("processPost() error = %v", err)
	}

	if post.ArticleHTML != originalHTML {
		t.Error("Skip post HTML was modified")
	}
}

func TestContributionGraphPlugin_ProcessPost_EmptyHTML(t *testing.T) {
	p := NewContributionGraphPlugin()

	post := &models.Post{ArticleHTML: ""}

	err := p.processPost(nil, post)
	if err != nil {
		t.Errorf("processPost() error = %v", err)
	}

	if post.ArticleHTML != "" {
		t.Error("Empty HTML was modified")
	}
}

func TestContributionGraphPlugin_ProcessPost_Disabled(t *testing.T) {
	p := NewContributionGraphPlugin()
	p.SetConfig(models.ContributionGraphConfig{
		Enabled: false,
		CDNURL:  "https://cdn.jsdelivr.net/npm/cal-heatmap@4",
	})

	if p.Config().Enabled {
		t.Error("Expected plugin to be disabled")
	}
}

func TestContributionGraphPlugin_ProcessPost_CustomContainerClass(t *testing.T) {
	p := NewContributionGraphPlugin()
	p.SetConfig(models.ContributionGraphConfig{
		Enabled:        true,
		CDNURL:         "https://example.com/cal-heatmap",
		ContainerClass: "my-custom-graph",
		Theme:          "dark",
	})

	post := &models.Post{ArticleHTML: `<pre><code class="language-contribution-graph">{"data": []}</code></pre>`}

	err := p.processPost(nil, post)
	if err != nil {
		t.Errorf("processPost() error = %v", err)
	}

	if !strings.Contains(post.ArticleHTML, `class="my-custom-graph"`) {
		t.Error("Expected custom container class")
	}
	if !strings.Contains(post.ArticleHTML, ".my-custom-graph") {
		t.Error("Expected custom class selector in injected CSS")
	}
	if !strings.Contains(post.ArticleHTML, "example.com/cal-heatmap") {
		t.Error("Expected custom CDN URL")
	}
}

func TestContributionGraphPlugin_ProcessPost_HTMLEntities(t *testing.T) {
	p := NewContributionGraphPlugin()

	post := &models.Post{ArticleHTML: `<pre><code class="language-contribution-graph">{&quot;data&quot;: [], &quot;options&quot;: {}}</code></pre>`}

	err := p.processPost(nil, post)
	if err != nil {
		t.Errorf("processPost() error = %v", err)
	}

	if strings.Contains(post.ArticleHTML, "contribution-graph-error") {
		t.Error("Should handle HTML entities correctly")
	}
	if !strings.Contains(post.ArticleHTML, "<div id=") {
		t.Error("Expected div element after HTML entity decoding")
	}
}

func TestContributionGraphPlugin_ProcessPost_PreserveSurroundingContent(t *testing.T) {
	p := NewContributionGraphPlugin()

	post := &models.Post{ArticleHTML: `<h1>Title</h1>
<pre><code class="language-contribution-graph">{"data": []}</code></pre>
<p>Footer content</p>`}

	err := p.processPost(nil, post)
	if err != nil {
		t.Errorf("processPost() error = %v", err)
	}

	if !strings.Contains(post.ArticleHTML, "<h1>Title</h1>") {
		t.Error("Lost header content")
	}
	if !strings.Contains(post.ArticleHTML, "<p>Footer content</p>") {
		t.Error("Lost footer content")
	}
}

func TestContributionGraphPlugin_ProcessPost_WithOptions(t *testing.T) {
	p := NewContributionGraphPlugin()

	post := &models.Post{ArticleHTML: `<pre><code class="language-contribution-graph">{
  "data": [{"date": "2024-01-01", "value": 5}],
  "options": {
    "domain": "month",
    "subDomain": "day",
    "cellSize": 15,
    "range": 12
  }
}</code></pre>`}

	err := p.processPost(nil, post)
	if err != nil {
		t.Errorf("processPost() error = %v", err)
	}

	if !strings.Contains(post.ArticleHTML, "domain:") {
		t.Error("Expected domain configuration in output")
	}
	if !strings.Contains(post.ArticleHTML, "cal.paint") {
		t.Error("Expected cal.paint call in output")
	}
}

func TestContributionGraphPlugin_ProcessPost_DataParsing(t *testing.T) {
	p := NewContributionGraphPlugin()

	post := &models.Post{ArticleHTML: `<pre><code class="language-contribution-graph">{
  "data": [
    {"date": "2024-01-01", "value": 5},
    {"date": "2024-01-02", "value": 10},
    {"date": "2024-01-03", "value": 3}
  ],
  "options": {}
}</code></pre>`}

	err := p.processPost(nil, post)
	if err != nil {
		t.Errorf("processPost() error = %v", err)
	}

	if !strings.Contains(post.ArticleHTML, "2024-01-01") {
		t.Error("Expected data to be included in output")
	}
	if !strings.Contains(post.ArticleHTML, "const maxValue = 10;") {
		t.Error("Expected scale max to be emitted into the script")
	}
}

func TestContributionGraphPlugin_ProcessPost_MaxPercentileCapsScale(t *testing.T) {
	p := NewContributionGraphPlugin()

	post := &models.Post{ArticleHTML: `<pre><code class="language-contribution-graph">{
  "data": [
    {"date": "2024-01-01", "value": 1},
    {"date": "2024-01-02", "value": 2},
    {"date": "2024-01-03", "value": 3},
    {"date": "2024-01-04", "value": 20}
  ],
  "options": {
    "maxPercentile": 95
  }
}</code></pre>`}

	err := p.processPost(nil, post)
	if err != nil {
		t.Fatalf("processPost() error = %v", err)
	}

	if !strings.Contains(post.ArticleHTML, "const maxValue = 20;") {
		t.Error("expected maxValue to be emitted into the graph script")
	}
}

func TestContributionGraphPlugin_ProcessPost_MaxValueOverridesScale(t *testing.T) {
	p := NewContributionGraphPlugin()

	post := &models.Post{ArticleHTML: `<pre><code class="language-contribution-graph">{
  "data": [
    {"date": "2024-01-01", "value": 1},
    {"date": "2024-01-02", "value": 2},
    {"date": "2024-01-03", "value": 20}
  ],
  "options": {
    "maxValue": 8
  }
}</code></pre>`}

	err := p.processPost(nil, post)
	if err != nil {
		t.Fatalf("processPost() error = %v", err)
	}

	if !strings.Contains(post.ArticleHTML, "const maxValue = 8;") {
		t.Error("expected explicit maxValue override to affect scale max")
	}
	if !strings.Contains(post.ArticleHTML, "source: displayData") {
		t.Error("expected maxValue override to clamp rendered data")
	}
}

func TestContributionGraphPlugin_ProcessPost_GlobalScaleMaxPercentile(t *testing.T) {
	p := NewContributionGraphPlugin()
	p.SetConfig(models.ContributionGraphConfig{
		Enabled:            true,
		CDNURL:             "https://example.com/cal-heatmap",
		ContainerClass:     "contribution-graph-container",
		Theme:              "light",
		ScaleMaxPercentile: 50,
	})

	post := &models.Post{ArticleHTML: `<pre><code class="language-contribution-graph">{
  "data": [
    {"date": "2024-01-01", "value": 1},
    {"date": "2024-01-02", "value": 2},
    {"date": "2024-01-03", "value": 9}
  ]
}</code></pre>`}

	err := p.processPost(nil, post)
	if err != nil {
		t.Fatalf("processPost() error = %v", err)
	}

	if !strings.Contains(post.ArticleHTML, "const maxValue = 2;") {
		t.Error("expected global percentile cap to affect scale max")
	}
	if !strings.Contains(post.ArticleHTML, "window._contributionGraphFitters") {
		t.Error("expected built-in fitters to be registered")
	}
}
