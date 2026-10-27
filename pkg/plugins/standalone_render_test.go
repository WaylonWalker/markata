package plugins

import (
	"strings"
	"testing"

	xhtml "golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func TestRenderStandalone_ConditionalAssets(t *testing.T) {
	tests := []struct{ name, source, feature, asset string }{
		{"code", "```go\nfmt.Println(1)\n```", "code", "css/chroma.css"},
		{"admonition", "!!! note\n    Hello", "admonitions", "css/admonitions.css"},
		{"table", "| A | B |\n|---|---|\n|1|2|", "tables", "js/sortable-tables.js"},
		{"checklist", "::: md-checklist\n\n- [ ] Todo\n\n:::", "checklists", "js/markdown-checklists.js"},
		{"road", "::: road\nHello\n:::", "roads", "js/roads.js"},
		{"mermaid", "```mermaid\nflowchart TD\n A-->B\n```", "mermaid", "assets/vendor/mermaid/mermaid.min.js"},
		{"chartjs", "```chartjs\n{\"type\":\"bar\",\"data\":{\"labels\":[\"A\"],\"datasets\":[{\"data\":[1]}]}}\n```", "chartjs", "assets/vendor/chartjs/chart.min.js"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := RenderStandalone(tt.source, StandaloneOptions{InlineAssets: true})
			if err != nil {
				t.Fatal(err)
			}
			if !r.Features[tt.feature] {
				t.Fatalf("missing feature %s: %v", tt.feature, r.Features)
			}
			found := false
			for _, asset := range r.Assets {
				if asset.ID == tt.asset && asset.Content != "" {
					found = true
				}
			}
			if !found {
				t.Fatalf("missing asset %s", tt.asset)
			}
			if len(r.ExternalURLs) != 0 {
				t.Fatalf("unexpected external dependencies: %v", r.ExternalURLs)
			}
			if strings.Contains(r.HTML, "<script") {
				t.Fatal("generated scripts leaked into fragment")
			}
			assertStandaloneInlineTags(t, r.HeadHTML+r.BodyEndHTML)
		})
	}
}

func TestRenderStandalone_PlainAndLinked(t *testing.T) {
	r, err := RenderStandalone("# Hello", StandaloneOptions{AssetBase: "/static/md"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(r.HTML, "Hello") || !strings.Contains(r.HeadHTML, `/static/md/css/main.css`) {
		t.Fatalf("unexpected result: %+v", r)
	}
	if r.BodyEndHTML != "" {
		t.Fatal("plain markdown must not load scripts")
	}
	for _, feature := range r.Features {
		if feature {
			t.Fatal("plain text requires optional assets")
		}
	}
}

func TestRenderStandalone_RawHTMLPolicy(t *testing.T) {
	source := `<script>alert(1)</script>`
	safe, err := RenderStandalone(source, StandaloneOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(safe.HTML, "<script>") {
		t.Fatal("raw HTML must be disabled by default")
	}
	trusted, err := RenderStandalone(source, StandaloneOptions{Unsafe: true})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(trusted.HTML, source) {
		t.Fatal("trusted HTML missing")
	}
}

func TestRenderStandalone_MermaidEscapesMarkup(t *testing.T) {
	r, err := RenderStandalone("```mermaid\ngraph TD\nA[\"</pre><script>alert(1)</script>\"]\n```", StandaloneOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(r.HTML, "<script>") || !strings.Contains(r.HTML, "&lt;/pre&gt;") {
		t.Fatal("diagram escaped its pre element")
	}
}

func TestRenderStandalone_DependenciesPrecedeInitializers(t *testing.T) {
	r, err := RenderStandalone("```mermaid\ngraph TD\nA-->B\n```", StandaloneOptions{})
	if err != nil {
		t.Fatal(err)
	}
	vendor, init := -1, -1
	for i, asset := range r.Assets {
		if asset.ID == "assets/vendor/mermaid/mermaid.min.js" {
			vendor = i
		}
		if asset.Module {
			init = i
		}
	}
	if vendor < 0 || init <= vendor {
		t.Fatal("mermaid initialization must follow its library")
	}
	if strings.Contains(r.BodyEndHTML, "import mermaid") {
		t.Fatal("initializer must use the bundled library")
	}
}

func assertStandaloneInlineTags(t *testing.T, source string) {
	t.Helper()
	z := xhtml.NewTokenizer(strings.NewReader(source))
	for {
		kind := z.Next()
		if kind == xhtml.ErrorToken {
			return
		}
		if kind != xhtml.StartTagToken && kind != xhtml.SelfClosingTagToken {
			continue
		}
		tag := z.Token()
		if tag.DataAtom == atom.Link {
			t.Fatal("inline output contains external styles")
		}
		if tag.Data == "script" {
			for _, attr := range tag.Attr {
				if attr.Key == "src" {
					t.Fatal("inline output contains external scripts")
				}
			}
		}
	}
}
