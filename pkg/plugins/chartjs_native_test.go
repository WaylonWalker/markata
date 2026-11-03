package plugins

import (
	"strings"
	"testing"
	"github.com/WaylonWalker/markata/pkg/models"
)

func TestChartJS_NativeAutoBasicBarNoBrowserBundle(t *testing.T) {
	p := NewChartJSPlugin()
	post := &models.Post{ArticleHTML: `<pre><code class="language-chartjs">{"type":"bar","data":{"labels":["A","B"],"datasets":[{"label":"Votes","data":[12,19]}]}}</code></pre>`}
	if err := p.processPost(post); err != nil { t.Fatal(err) }
	if !strings.Contains(post.ArticleHTML, "<svg") || !strings.Contains(post.ArticleHTML, "chartjs-data") {
		t.Fatalf("native chart should retain its data table: %s", post.ArticleHTML)
	}
	if strings.Contains(post.ArticleHTML, "<script") || strings.Contains(post.ArticleHTML, "chart.min.js") {
		t.Fatal("native chart must not ship Chart.js to readers")
	}
}

func TestChartJS_NativeAutoSupportedTypes(t *testing.T) {
	cases := map[string]string{
		"bar": `{"type":"bar","data":{"labels":["A","B"],"datasets":[{"data":[12,19]}]}}`,
		"line": `{"type":"line","data":{"labels":["A","B"],"datasets":[{"data":[12,19]}]}}`,
		"area": `{"type":"line","data":{"labels":["A","B"],"datasets":[{"data":[12,19],"fill":true}]}}`,
		"scatter": `{"type":"scatter","data":{"datasets":[{"data":[{"x":1,"y":3},{"x":2,"y":6}]}]}}`,
		"pie": `{"type":"pie","data":{"labels":["A","B"],"datasets":[{"data":[12,19]}]}}`,
		"doughnut": `{"type":"doughnut","data":{"labels":["A","B"],"datasets":[{"data":[12,19]}]}}`,
	}
	for name, source := range cases {
		t.Run(name, func(t *testing.T) {
			p := NewChartJSPlugin()
			post := &models.Post{ArticleHTML: "<pre><code class=\"language-chartjs\">" + source + "</code></pre>"}
			if err := p.processPost(post); err != nil { t.Fatal(err) }
			if !strings.Contains(post.ArticleHTML, "<svg") || strings.Contains(post.ArticleHTML, "<script") {
				t.Fatalf("%s was not rendered natively: %s", name, post.ArticleHTML)
			}
		})
	}
}

func TestChartJS_AutoFallbackAndStrictGo(t *testing.T) {
	source := `<pre><code class="language-chartjs">{"type":"bar","data":{"labels":["A"],"datasets":[{"data":[10]}]},"options":{"plugins":{"legend":{"display":false}}}}</code></pre>`
	auto := NewChartJSPlugin()
	post := &models.Post{ArticleHTML: source}
	if err := auto.processPost(post); err != nil { t.Fatal(err) }
	if !strings.Contains(post.ArticleHTML, "new Chart(ctx,") {
		t.Fatalf("advanced options must retain full Chart.js: %s", post.ArticleHTML)
	}
	strict := NewChartJSPlugin()
	strict.config.Mode = "go"
	other := &models.Post{ArticleHTML: source}
	if err := strict.processPost(other); err != nil { t.Fatal(err) }
	if !strings.Contains(other.ArticleHTML, "Native chart unsupported") || strings.Contains(other.ArticleHTML, "new Chart") {
		t.Fatalf("strict native mode must report unsupported chart options: %s", other.ArticleHTML)
	}
}

func TestRenderStandalone_DefaultAutoNativeChart(t *testing.T) {
	source := "```chartjs\n{\"type\":\"bar\",\"data\":{\"labels\":[\"A\"],\"datasets\":[{\"data\":[10]}]}}\n```"
	result, err := RenderStandalone(source, StandaloneOptions{})
	if err != nil { t.Fatal(err) }
	if !result.Features["chartjs"] || !strings.Contains(result.HTML, "<svg") {
		t.Fatalf("native chart feature not surfaced: %+v", result)
	}
	for _, asset := range result.Assets {
		if strings.Contains(asset.ID, "chart.min.js") {
			t.Fatalf("static chart shipped client library: %s", asset.ID)
		}
	}
}
