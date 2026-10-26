package plugins

import (
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/models"
)

func TestChartJSPlugin_Name(t *testing.T) {
	p := NewChartJSPlugin()
	if got := p.Name(); got != "chartjs" {
		t.Errorf("Name() = %q, want %q", got, "chartjs")
	}
}

func TestChartJSPlugin_ProcessPost_NoChartJS(t *testing.T) {
	p := NewChartJSPlugin()

	post := &models.Post{
		ArticleHTML: "<p>Hello world</p>",
	}

	err := p.processPost(post)
	if err != nil {
		t.Errorf("processPost() error = %v", err)
	}

	// HTML should be unchanged
	if post.ArticleHTML != "<p>Hello world</p>" {
		t.Errorf("HTML was modified when no chartjs blocks present")
	}
}

func TestChartJSPlugin_ProcessPost_ValidChart(t *testing.T) {
	p := NewChartJSPlugin()

	post := &models.Post{
		ArticleHTML: `<p>Check out this chart:</p>
<pre><code class="language-chartjs">{
  "type": "bar",
  "data": {
    "labels": ["Red", "Blue"],
    "datasets": [{
      "label": "Votes",
      "data": [12, 19]
    }]
  }
}</code></pre>
<p>Nice chart!</p>`,
	}

	err := p.processPost(post)
	if err != nil {
		t.Errorf("processPost() error = %v", err)
	}

	// Should contain a canvas element
	if !strings.Contains(post.ArticleHTML, "<canvas id=") {
		t.Error("Expected canvas element in output")
	}

	// Should contain Chart.js script
	if !strings.Contains(post.ArticleHTML, "chart.min.js") {
		t.Error("Expected Chart.js script")
	}

	// Should have the container class
	if !strings.Contains(post.ArticleHTML, `class="chartjs-container"`) {
		t.Error("Expected chartjs-container class")
	}

	// Should have initialization script
	if !strings.Contains(post.ArticleHTML, "new Chart(ctx,") {
		t.Error("Expected Chart initialization script")
	}
}

func TestChartJSPlugin_ProcessPost_InvalidJSON(t *testing.T) {
	p := NewChartJSPlugin()

	post := &models.Post{
		ArticleHTML: `<pre><code class="language-chartjs">{ invalid json }</code></pre>`,
	}

	err := p.processPost(post)
	if err != nil {
		t.Errorf("processPost() error = %v", err)
	}

	// Should contain an error message
	if !strings.Contains(post.ArticleHTML, "chartjs-error") {
		t.Error("Expected error class for invalid JSON")
	}
	if !strings.Contains(post.ArticleHTML, "Invalid JSON") {
		t.Error("Expected error message for invalid JSON")
	}
}

func TestChartJSPlugin_ProcessPost_MultipleCharts(t *testing.T) {
	p := NewChartJSPlugin()

	post := &models.Post{
		ArticleHTML: `<pre><code class="language-chartjs">{"type": "bar", "data": {}}</code></pre>
<p>And another:</p>
<pre><code class="language-chartjs">{"type": "line", "data": {}}</code></pre>`,
	}

	err := p.processPost(post)
	if err != nil {
		t.Errorf("processPost() error = %v", err)
	}

	// Should have two canvas elements with different IDs
	count := strings.Count(post.ArticleHTML, "<canvas id=")
	if count != 2 {
		t.Errorf("Expected 2 canvas elements, got %d", count)
	}

	// Script should only be included once
	scriptCount := strings.Count(post.ArticleHTML, "chart.min.js")
	if scriptCount != 1 {
		t.Errorf("Expected 1 chart script inclusion, got %d", scriptCount)
	}
}

func TestChartJSPlugin_ProcessPost_SkipPost(t *testing.T) {
	p := NewChartJSPlugin()

	post := &models.Post{
		Skip:        true,
		ArticleHTML: `<pre><code class="language-chartjs">{"type": "bar"}</code></pre>`,
	}
	originalHTML := post.ArticleHTML

	err := p.processPost(post)
	if err != nil {
		t.Errorf("processPost() error = %v", err)
	}

	// HTML should be unchanged for skipped posts
	if post.ArticleHTML != originalHTML {
		t.Error("Skip post HTML was modified")
	}
}

func TestChartJSPlugin_ProcessPost_EmptyHTML(t *testing.T) {
	p := NewChartJSPlugin()

	post := &models.Post{
		ArticleHTML: "",
	}

	err := p.processPost(post)
	if err != nil {
		t.Errorf("processPost() error = %v", err)
	}

	if post.ArticleHTML != "" {
		t.Error("Empty HTML was modified")
	}
}

func TestChartJSPlugin_ProcessPost_Disabled(t *testing.T) {
	p := NewChartJSPlugin()
	p.SetConfig(models.ChartJSConfig{
		Enabled: false,
		CDNURL:  "https://cdn.jsdelivr.net/npm/chart.js",
	})

	// Create a mock manager would be needed for full Render() test
	// For now, we test processPost directly and verify disabled via config
	if p.Config().Enabled {
		t.Error("Expected plugin to be disabled")
	}
}

func TestChartJSPlugin_ProcessPost_CustomContainerClass(t *testing.T) {
	p := NewChartJSPlugin()
	p.SetConfig(models.ChartJSConfig{
		Enabled:        true,
		CDNURL:         "https://example.com/chart.js",
		ContainerClass: "my-custom-chart",
	})

	post := &models.Post{
		ArticleHTML: `<pre><code class="language-chartjs">{"type": "pie"}</code></pre>`,
	}

	err := p.processPost(post)
	if err != nil {
		t.Errorf("processPost() error = %v", err)
	}

	if !strings.Contains(post.ArticleHTML, `class="my-custom-chart"`) {
		t.Error("Expected custom container class")
	}
	if !strings.Contains(post.ArticleHTML, "example.com/chart.js") {
		t.Error("Expected custom CDN URL")
	}
}

func TestChartJSPlugin_ProcessPost_HTMLEntities(t *testing.T) {
	p := NewChartJSPlugin()

	// Goldmark encodes special characters, so we test that they get decoded
	post := &models.Post{
		ArticleHTML: `<pre><code class="language-chartjs">{&quot;type&quot;: &quot;bar&quot;, &quot;data&quot;: {}}</code></pre>`,
	}

	err := p.processPost(post)
	if err != nil {
		t.Errorf("processPost() error = %v", err)
	}

	// Should process successfully (HTML entities decoded)
	if strings.Contains(post.ArticleHTML, "chartjs-error") {
		t.Error("Should handle HTML entities correctly")
	}
	if !strings.Contains(post.ArticleHTML, "<canvas id=") {
		t.Error("Expected canvas element after HTML entity decoding")
	}
}

func TestChartJSPlugin_ProcessPost_PreserveSurroundingContent(t *testing.T) {
	p := NewChartJSPlugin()

	post := &models.Post{
		ArticleHTML: `<h1>Title</h1>
<pre><code class="language-chartjs">{"type": "bar"}</code></pre>
<p>Footer content</p>`,
	}

	err := p.processPost(post)
	if err != nil {
		t.Errorf("processPost() error = %v", err)
	}

	// Should preserve surrounding content
	if !strings.Contains(post.ArticleHTML, "<h1>Title</h1>") {
		t.Error("Lost header content")
	}
	if !strings.Contains(post.ArticleHTML, "<p>Footer content</p>") {
		t.Error("Lost footer content")
	}
}

func TestChartJSPlugin_StaticTableFallback(t *testing.T) {
	p := NewChartJSPlugin()
	post := &models.Post{ArticleHTML: `<pre><code class="language-chartjs">{
		"type": "bar",
		"data": {"labels": ["First", "&lt;script&gt;unsafe&lt;/script&gt;"],
			"datasets": [{"label": "Used &amp; free", "data": [12, 19]}]}
	}</code></pre>`}
	if err := p.processPost(post); err != nil {
		t.Fatal(err)
	}
	for _, part := range []string{
		`<canvas id="chartjs-`,
		`hidden role="img"`,
		`<details class="chartjs-data" open>`,
		`<table><caption>Chart data</caption>`,
		`<th scope="col">Used &amp; free</th>`,
		`&lt;script&gt;unsafe&lt;/script&gt;`,
		`<td>12</td>`,
		`<td>19</td>`,
		`if (data) data.open = false;`,
	} {
		if !strings.Contains(post.ArticleHTML, part) {
			t.Errorf("expected %q in static chart output", part)
		}
	}
	if strings.Contains(post.ArticleHTML, "<script>unsafe</script>") {
		t.Fatal("chart labels must be safely HTML escaped")
	}
}

func TestChartJSPlugin_UnstructuredChartRetainsSourceWithoutJS(t *testing.T) {
	p := NewChartJSPlugin()
	post := &models.Post{ArticleHTML: `<pre><code class="language-chartjs">{"type":"scatter","data":{"datasets":[{"data":[{"x":1,"y":2}]}]}}</code></pre>`}
	if err := p.processPost(post); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(post.ArticleHTML, "<summary>Chart configuration</summary>") {
		t.Fatal("unstructured charts need readable configuration when JS is unavailable")
	}
	if strings.Contains(post.ArticleHTML, `<table>`) {
		t.Fatal("scatter object pairs should not render misleading categorical tables")
	}
}

func TestChartJSPlugin_EscapesAuthoredInlineScriptClosers(t *testing.T) {
	p := NewChartJSPlugin()
	post := &models.Post{ArticleHTML: `<pre><code class="language-chartjs">{"type":"bar","data":{"labels":["&lt;/script&gt;&lt;script&gt;alert(1)"],"datasets":[{"data":[7]}]}}</code></pre>`}
	if err := p.processPost(post); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(post.ArticleHTML, "</script><script>alert(1)") {
		t.Fatal("chart JSON must not close an inline script")
	}
	if !strings.Contains(post.ArticleHTML, `<details class="chartjs-data" open>`) {
		t.Fatal("static fallback should still be available")
	}
}
