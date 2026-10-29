package plugins

import (
	"fmt"
	"html"
	"strings"

	"github.com/alecthomas/chroma/v2/styles"

	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/palettes"
	"github.com/WaylonWalker/markata/pkg/themes"

	xhtml "golang.org/x/net/html"
)

// StandaloneOptions configures rendering without a site build.
type StandaloneOptions struct {
	Unsafe       bool
	InlineAssets bool
	AssetBase    string
}

// RenderAsset is an ordered dependency for a rendered Markdown fragment.
type RenderAsset struct {
	ID        string `json:"id"`
	Kind      string `json:"kind"`
	Placement string `json:"placement"`
	Content   string `json:"content,omitempty"`
	URL       string `json:"url,omitempty"`
	Module    bool   `json:"module,omitempty"`
}

// StandaloneResult includes markup and everything needed to present it.
type StandaloneResult struct {
	SchemaVersion int             `json:"schema_version"`
	HTML          string          `json:"html"`
	Features      map[string]bool `json:"features"`
	Assets        []RenderAsset   `json:"assets"`
	ExternalURLs  []string        `json:"external_urls"`
	HeadHTML      string          `json:"head_html"`
	BodyEndHTML   string          `json:"body_end_html"`
}

// RenderStandalone uses the site's Markdown, Mermaid, and Chart.js renderers.
// It performs no network access or filesystem writes and loads no site config.
func RenderStandalone(source string, options StandaloneOptions) (*StandaloneResult, error) {
	p := &RenderMarkdownPlugin{md: createMarkdownRenderer(palettes.DefaultChromaThemeDark, false, DefaultMarkdownExtensionConfig(), options.Unsafe)}
	post := &models.Post{Content: source, Extra: map[string]interface{}{}}
	if err := p.renderPost(post); err != nil {
		return nil, fmt.Errorf("render markdown: %w", err)
	}
	r := &StandaloneResult{SchemaVersion: StandaloneSchemaVersion, Features: map[string]bool{"mermaid": false, "chartjs": false}, Assets: []RenderAsset{}, ExternalURLs: []string{}}
	if err := r.renderDiagrams(post); err != nil {
		return nil, err
	}
	post.Extra = map[string]interface{}{}
	p.detectCSSRequirements(post)
	r.HTML = expandStandaloneIcons(post.ArticleHTML)
	r.detectEnhancements()
	for signal, feature := range map[string]string{
		"needs_code_css": "code", "needs_admonitions_css": "admonitions",
		"needs_roads": "roads", "has_figure": "figures", "has_anchor_links": "anchors",
	} {
		r.Features[feature] = post.Extra[signal] == true
	}
	if err := r.addPresentationAssets(); err != nil {
		return nil, err
	}
	for _, asset := range r.Assets {
		tag := asset.tag(options)
		if asset.Placement == "head" {
			r.HeadHTML += tag + "\n"
		} else {
			r.BodyEndHTML += tag + "\n"
		}
	}
	return r, nil
}

func (r *StandaloneResult) renderDiagrams(post *models.Post) error {
	mermaid := NewMermaidPlugin()
	mermaid.config.Lightbox = false
	mermaid.escapeClientSource = true
	var mermaidScript string
	mermaid.scriptCollector = func(fragment string) {
		mermaidScript = strings.Replace(fragment, mermaid.mermaidImportJS(), "const mermaid = window.mermaid;\n", 1)
	}
	if err := mermaid.processPost(post); err != nil {
		return fmt.Errorf("render mermaid: %w", err)
	}
	if mermaidScript != "" {
		r.Features["mermaid"] = true
		if err := r.addBundled("assets/vendor/mermaid/mermaid.min.js", standaloneScriptKind, "body_end"); err != nil {
			return err
		}
		r.collectScripts(mermaidScript)
	}
	chart := NewChartJSPlugin()
	chart.scriptCollector = func(fragment string) { r.collectScripts(fragment) }
	if err := chart.processPost(post); err != nil {
		return fmt.Errorf("render chartjs: %w", err)
	}
	for i := range r.Assets {
		asset := &r.Assets[i]
		if asset.URL == "" {
			continue
		}
		r.Features["chartjs"] = true
		name := strings.TrimPrefix(asset.URL, "/")
		content, err := themes.ReadStatic(name)
		if err != nil {
			return fmt.Errorf("read chart library: %w", err)
		}
		asset.ID, asset.URL, asset.Content = name, "", string(content)
	}
	return nil
}

func (r *StandaloneResult) addPresentationAssets() error {
	// Put styles before scripts while preserving generated dependency order.
	scripts := r.Assets
	r.Assets = []RenderAsset{}
	files := []string{"css/variables.css", "css/main.css", "css/components.css"}
	for _, item := range []struct{ feature, path string }{
		{"admonitions", "css/admonitions.css"}, {"code", "css/code.css"},
		{"mermaid", "css/mermaid.css"}, {"roads", "css/roads.css"},
	} {
		if r.Features[item.feature] {
			files = append(files, item.path)
		}
	}
	for _, name := range files {
		if err := r.addBundled(name, "style", "head"); err != nil {
			return err
		}
	}
	if r.Features["code"] {
		css, err := NewChromaCSSPlugin().generateThemeCSS(styles.Get(palettes.DefaultChromaThemeDark))
		if err != nil {
			return fmt.Errorf("generate highlighting CSS: %w", err)
		}
		r.Assets = append(r.Assets, RenderAsset{ID: "css/chroma.css", Kind: "style", Placement: "head", Content: css})
		if err := r.addBundled("js/code-chrome.js", standaloneScriptKind, "body_end"); err != nil {
			return err
		}
	}
	for _, item := range []struct{ feature, path string }{
		{"roads", "js/roads.js"}, {"tables", "js/sortable-tables.js"},
		{"checklists", "js/markdown-checklists.js"},
	} {
		if r.Features[item.feature] {
			if err := r.addBundled(item.path, standaloneScriptKind, "body_end"); err != nil {
				return err
			}
		}
	}
	r.Assets = append(r.Assets, scripts...)
	return nil
}

func (r *StandaloneResult) addBundled(name, kind, placement string) error {
	content, err := themes.ReadStatic(name)
	if err != nil {
		return fmt.Errorf("read render asset %s: %w", name, err)
	}
	r.Assets = append(r.Assets, RenderAsset{ID: name, Kind: kind, Placement: placement, Content: string(content)})
	return nil
}

func (r *StandaloneResult) detectEnhancements() {
	r.Features["tables"], r.Features["checklists"] = false, false
	z := xhtml.NewTokenizer(strings.NewReader(r.HTML))
	for {
		t := z.Next()
		if t == xhtml.ErrorToken {
			return
		}
		if t != xhtml.StartTagToken {
			continue
		}
		token := z.Token()
		if token.Data == "table" {
			r.Features["tables"] = true
		}
		for _, attr := range token.Attr {
			if attr.Key == htmlAttrClass {
				for _, class := range strings.Fields(attr.Val) {
					if class == "md-checklist" {
						r.Features["checklists"] = true
					}
				}
			}
		}
	}
}

// collectScripts parses only script tags generated by the rendering plugins.
// Authored HTML is never moved into the dependency manifest.
func (r *StandaloneResult) collectScripts(fragment string) {
	z := xhtml.NewTokenizer(strings.NewReader(fragment))
	for {
		t := z.Next()
		if t == xhtml.ErrorToken {
			return
		}
		if t != xhtml.StartTagToken {
			continue
		}
		token := z.Token()
		if token.Data != standaloneScriptKind {
			continue
		}
		a := RenderAsset{ID: fmt.Sprintf("generated-%d", len(r.Assets)), Kind: standaloneScriptKind, Placement: "body_end"}
		for _, attr := range token.Attr {
			if attr.Key == "src" {
				a.URL = attr.Val
			}
			if attr.Key == "type" && attr.Val == "module" {
				a.Module = true
			}
		}
		for {
			t = z.Next()
			if t == xhtml.ErrorToken || t == xhtml.EndTagToken {
				break
			}
			if t == xhtml.TextToken {
				a.Content += string(z.Text())
			}
		}
		r.Assets = append(r.Assets, a)
	}
}

func (a RenderAsset) tag(options StandaloneOptions) string {
	url := a.URL
	bundled := strings.HasPrefix(a.ID, "css/") || strings.HasPrefix(a.ID, "js/") || strings.HasPrefix(a.ID, "assets/vendor/")
	if bundled && !options.InlineAssets {
		base := options.AssetBase
		if base == "" {
			base = "/markata-assets"
		}
		url = strings.TrimRight(base, "/") + "/" + a.ID
	}
	if a.Kind == "style" {
		if url != "" {
			return `<link rel="stylesheet" href="` + html.EscapeString(url) + `">`
		}
		return "<style>" + strings.ReplaceAll(a.Content, "</style", `<\/style`) + "</style>"
	}
	typeAttr := ""
	if a.Module {
		typeAttr = ` type="module"`
	}
	if url != "" {
		return `<script` + typeAttr + ` src="` + html.EscapeString(url) + `"></script>`
	}
	return "<script" + typeAttr + ">" + strings.ReplaceAll(a.Content, "</script", `<\/script`) + "</script>"
}
