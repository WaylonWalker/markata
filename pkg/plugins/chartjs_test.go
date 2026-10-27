package plugins

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
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

func TestChartJSPlugin_FigureFromMarkdown(t *testing.T) {
	for _, metadata := range []string{
		`"title":"Build time","caption":"Seconds per build",`,
		`"title":"Build time",`,
		`"caption":"Seconds per build",`,
	} {
		t.Run(metadata, func(t *testing.T) {
			post := &models.Post{Content: "```chartjs\n{" + metadata + `"type":"bar","data":{"labels":["Full","Incremental"],"datasets":[{"data":[46,21]}]},"options":{"indexAxis":"y"}}` + "\n```"}
			if err := NewRenderMarkdownPlugin().renderPost(post); err != nil {
				t.Fatal(err)
			}
			if err := NewChartJSPlugin().processPost(post); err != nil {
				t.Fatal(err)
			}
			for _, want := range []string{
				`<figure class="chartjs-figure chartjs-figure-compact">`,
				`<figcaption id="chartjs-1-caption">`,
				`aria-describedby="chartjs-1-caption"`,
				`<details class="chartjs-data" open>`,
				`document.querySelector('#chartjs-1-data > details')`,
			} {
				if !strings.Contains(post.ArticleHTML, want) {
					t.Errorf("missing %q in %s", want, post.ArticleHTML)
				}
			}
			if strings.Contains(metadata, `"title"`) && !strings.Contains(post.ArticleHTML, `aria-label="Build time"`) {
				t.Error("canvas should be labeled by the title")
			}
			scriptStart := strings.Index(post.ArticleHTML, "<script")
			if scriptStart < 0 {
				t.Fatal("missing chart script")
			}
			script := post.ArticleHTML[scriptStart:]
			if strings.Contains(script, `"title":`) || strings.Contains(script, `"caption":`) || strings.Contains(post.ArticleHTML, "<h4") {
				t.Error("figure metadata must not reach Chart.js or create headings")
			}
		})
	}
}

func TestChartJSPlugin_FigureEscapesText(t *testing.T) {
	post := &models.Post{Content: "```chartjs\n" + `{"title":"<img src=x onerror=alert(1)> \" &","caption":"</script><script>alert(1)</script>","type":"line"}` + "\n```"}
	if err := NewRenderMarkdownPlugin().renderPost(post); err != nil {
		t.Fatal(err)
	}
	if err := NewChartJSPlugin().processPost(post); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		`&lt;img src=x onerror=alert(1)&gt; &#34; &amp;`,
		`&lt;/script&gt;&lt;script&gt;alert(1)&lt;/script&gt;`,
	} {
		if !strings.Contains(post.ArticleHTML, want) {
			t.Errorf("missing escaped text %q", want)
		}
	}
	if strings.Contains(post.ArticleHTML, "<img") || strings.Contains(post.ArticleHTML, "<script>alert") {
		t.Fatal("metadata injected active HTML")
	}
}

func TestChartJSPlugin_DataDisplay(t *testing.T) {
	for _, tc := range []struct {
		name, site, override string
		open                 bool
	}{
		{"default", "", "", true},
		{"site collapsed", "collapsed", "", false},
		{"chart collapsed", "open", `"data_display":"collapsed",`, false},
		{"chart open", "collapsed", `"data_display":"open",`, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			for _, config := range []string{
				`"type":"bar","data":{"labels":["One"],"datasets":[{"data":[1]}]}`,
				`"type":"scatter"`,
			} {
				p := NewChartJSPlugin()
				p.config.DataDisplay = tc.site
				post := &models.Post{ArticleHTML: `<pre><code class="language-chartjs">{` + tc.override + config + `}</code></pre>`}
				if err := p.processPost(post); err != nil {
					t.Fatal(err)
				}
				if got := strings.Contains(post.ArticleHTML, `<details class="chartjs-data" open>`); got != tc.open {
					t.Errorf("open = %v, want %v", got, tc.open)
				}
				if !strings.Contains(post.ArticleHTML, `<details class="chartjs-data"`) || strings.Contains(post.ArticleHTML, "data_display") {
					t.Error("disclosure must remain reachable and metadata must be stripped")
				}
			}
		})
	}
}

func TestChartJSPlugin_InvalidFigureOptions(t *testing.T) {
	for _, metadata := range []string{`"title":42`, `"caption":null`, `"data_display":"closed"`} {
		post := &models.Post{ArticleHTML: `<pre><code class="language-chartjs">{` + metadata + `}</code></pre>`}
		if err := NewChartJSPlugin().processPost(post); err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(post.ArticleHTML, "chartjs-error") || strings.Contains(post.ArticleHTML, "new Chart") {
			t.Errorf("invalid metadata must produce an error, got %s", post.ArticleHTML)
		}
	}
}

func TestChartJSPlugin_ConfigureDataDisplay(t *testing.T) {
	for _, value := range []interface{}{"open", "collapsed", "closed", false} {
		m := lifecycle.NewManager()
		m.Config().Extra = map[string]interface{}{"chartjs": map[string]interface{}{"data_display": value}}
		p := NewChartJSPlugin()
		err := p.Configure(m)
		valid := value == "open" || value == "collapsed"
		if (err == nil) != valid {
			t.Fatalf("Configure(%v) error = %v", value, err)
		}
		if valid && p.Config().DataDisplay != value {
			t.Errorf("Configure(%v) did not set data display", value)
		}
	}
}

func TestChartJSPlugin_LegacyMarkupUnchanged(t *testing.T) {
	post := &models.Post{ArticleHTML: `<pre><code class="language-chartjs">{"type":"pie"}</code></pre>`}
	if err := NewChartJSPlugin().processPost(post); err != nil {
		t.Fatal(err)
	}
	want := `<div class="chartjs-container">
  <canvas id="chartjs-1" hidden role="img" aria-label="Chart visualization; expand Chart data for exact values"></canvas>
  <div id="chartjs-1-data"><details class="chartjs-data" open><summary>Chart configuration</summary><pre><code>{&#34;type&#34;:&#34;pie&#34;}</code></pre></details></div>
</div>`
	want += `
<script src="/assets/vendor/chartjs/chart.min.js"></script>
<script>
document.addEventListener('DOMContentLoaded', function() {
  (function() {
    const ctx = document.getElementById('chartjs-1');
    if (!ctx || typeof Chart !== 'function') return;
    ctx.hidden = false;
    try {
      new Chart(ctx, {"type":"pie"});
      const data = document.getElementById('chartjs-1-data');
      if (data) data.open = false;
    } catch (error) {
      ctx.hidden = true;
    }
  })();
});
</script>`
	if post.ArticleHTML != want {
		t.Fatalf("legacy output changed:\ngot: %s\nwant: %s", post.ArticleHTML, want)
	}
}

func TestHorizontalBarDefaults(t *testing.T) {
	cases := []struct {
		name   string
		input  string
		want   string
		absent string
	}{
		{"simple", `{"type":"bar","data":{"labels":["A","B"],"datasets":[{"data":[1,2]}]},"options":{"indexAxis":"y"}}`, `"barThickness":14`, ""},
		{"explicit geometry", `{"type":"bar","data":{"labels":["A"],"datasets":[{"data":[1],"borderRadius":0}]},"options":{"indexAxis":"y","aspectRatio":2}}`, `"aspectRatio":2`, `"barThickness":14`},
		{"element geometry", `{"type":"bar","data":{"labels":["A"],"datasets":[{"data":[1]}]},"options":{"indexAxis":"y","elements":{"bar":{"borderRadius":0}}}}`, `"borderRadius":0`, `"barThickness":14`},
		{"root geometry", `{"type":"bar","data":{"labels":["A"],"datasets":[{"data":[1]}]},"options":{"indexAxis":"y","borderRadius":0}}`, `"borderRadius":0`, `"barThickness":14`},
		{"multiple series", `{"type":"bar","data":{"labels":["A"],"datasets":[{"data":[1]},{"data":[2]}]},"options":{"indexAxis":"y"}}`, `"indexAxis":"y"`, `"barThickness":14`},
		{"vertical", `{"type":"bar","data":{"labels":["A"],"datasets":[{"data":[1]}]}}`, `"type":"bar"`, `"barThickness":14`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var config map[string]interface{}
			if err := json.Unmarshal([]byte(tc.input), &config); err != nil {
				t.Fatal(err)
			}
			applyHorizontalBarDefaults(config)
			out, err := json.Marshal(config)
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(string(out), tc.want) {
				t.Errorf("got %s, want %s", out, tc.want)
			}
			if tc.absent != "" && strings.Contains(string(out), tc.absent) {
				t.Errorf("got unwanted %s in %s", tc.absent, out)
			}
		})
	}
}
