package plugins

import (
	"fmt"
	"html"
	"strings"

	"github.com/alecthomas/chroma/v2/styles"

	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/palettes"
	"github.com/WaylonWalker/markata/pkg/themes"
	"github.com/yuin/goldmark/text"

	xhtml "golang.org/x/net/html"
)

// StandaloneOptions configures rendering without a site build.
type StandaloneOptions struct {
	// Config applies site Markdown and highlighting settings when supplied.
	Config       *models.Config
	Unsafe       bool
	InlineAssets bool
	AssetBase    string
	// MermaidMode selects "client" (default) or "go" (native SVG).
	MermaidMode string
	// OmitAssetContent keeps bundled assets linked and omits their bytes.
	// Browser hosts must export and serve the native presentation bundle.
	OmitAssetContent bool
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
	SchemaVersion    int             `json:"schema_version"`
	HTML             string          `json:"html"`
	Features         map[string]bool `json:"features"`
	Assets           []RenderAsset   `json:"assets"`
	ExternalURLs     []string        `json:"external_urls"`
	HeadHTML         string          `json:"head_html"`
	BodyEndHTML      string          `json:"body_end_html"`
	omitAssetContent bool
}

// RenderStandalone uses the site's Markdown, Mermaid, and Chart.js renderers.
// It performs no network access or filesystem writes. Hosts may supply site config.
func RenderStandalone(source string, options StandaloneOptions) (*StandaloneResult, error) {
	return NewStandaloneRenderer(options).Render(source)
}

// StandaloneRenderer reuses a configured Markdown parser for sequential renders.
// It is not safe for concurrent use.
type StandaloneRenderer struct {
	parser  *RenderMarkdownPlugin
	options StandaloneOptions
}

// SourceLines returns the zero-based source line for each top-level rendered block.
// It uses the same parser as Render so custom Markdown blocks keep their position.
func (renderer *StandaloneRenderer) SourceLines(source string) []int {
	prepared := []byte(prepareMarkdownAttribution(source))
	root := renderer.parser.md.Parser().Parse(text.NewReader(prepared))
	lines := make([]int, 0)
	for node := root.FirstChild(); node != nil; node = node.NextSibling() {
		start := node
		for start.Lines().Len() == 0 && start.FirstChild() != nil {
			start = start.FirstChild()
		}
		if start.Lines().Len() == 0 {
			continue
		}
		offset := start.Lines().At(0).Start
		line := 0
		for _, b := range prepared[:offset] {
			if b == '\n' {
				line++
			}
		}
		lines = append(lines, line)
	}
	return lines
}

// NewStandaloneRenderer creates a renderer for a persistent host or worker.
func NewStandaloneRenderer(options StandaloneOptions) *StandaloneRenderer {
	p := &RenderMarkdownPlugin{md: createMarkdownRenderer(palettes.DefaultChromaThemeDark, false, DefaultMarkdownExtensionConfig(), options.Unsafe)}
	if options.Config != nil {
		extra := map[string]interface{}{"theme": options.Config.Theme, "models_config": options.Config, "markdown": options.Config.MarkdownConfig}
		chroma, numbers := p.resolveHighlightConfig(extra)
		p.md = createMarkdownRenderer(chroma, numbers, p.resolveExtensionConfig(extra), options.Unsafe)
	}
	return &StandaloneRenderer{parser: p, options: options}
}

// Render renders one independent document with fresh IDs and feature signals.
func (renderer *StandaloneRenderer) Render(source string) (*StandaloneResult, error) {
	p, options := renderer.parser, renderer.options
	post := &models.Post{Content: source, Extra: map[string]interface{}{}}
	article, err := p.renderArticleHTML(source)
	if err != nil {
		return nil, fmt.Errorf("render markdown: %w", err)
	}
	post.ArticleHTML = article
	r := &StandaloneResult{omitAssetContent: options.OmitAssetContent, SchemaVersion: StandaloneSchemaVersion, Features: map[string]bool{"mermaid": false, "chartjs": false}, Assets: []RenderAsset{}, ExternalURLs: []string{}}
	if err := r.renderDiagrams(post, options); err != nil {
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

func (r *StandaloneResult) renderDiagrams(post *models.Post, options StandaloneOptions) error {
	mermaid := NewMermaidPlugin()
	mermaid.config.Lightbox = false
	mermaid.escapeClientSource = true
	mode := options.MermaidMode
	if mode == "" {
		mode = mermaidModeClient
	}
	if mode != mermaidModeClient && mode != mermaidModeGo {
		return fmt.Errorf("unsupported standalone Mermaid mode %q", mode)
	}
	mermaid.config.Mode = mode
	if mode == mermaidModeGo {
		if options.Config != nil {
			mermaid.paletteColors = resolvePaletteColors(map[string]interface{}{"theme": options.Config.Theme})
			if raw, ok := options.Config.Extra["mermaid"].(map[string]interface{}); ok {
				mermaid.parseGoConfig(raw)
			}
		}
		renderer, err := newMermaidRenderer(mermaid.config, mermaid.paletteColors)
		if err != nil {
			return fmt.Errorf("initialize native mermaid: %w", err)
		}
		mermaid.renderer = renderer
		defer renderer.close()
	}
	var mermaidScript string
	mermaid.scriptCollector = func(fragment string) {
		mermaidScript = strings.Replace(fragment, mermaid.mermaidImportJS(), "const mermaid = window.mermaid;\n", 1)
	}
	if err := mermaid.processPost(post); err != nil {
		return fmt.Errorf("render mermaid: %w", err)
	}
	if mode == mermaidModeGo && strings.Contains(post.ArticleHTML, `<pre class="mermaid">`) {
		r.Features["mermaid"] = true
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
		var content []byte
		var err error
		if !r.omitAssetContent {
			content, err = themes.ReadStatic(name)
		}
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
	var content []byte
	var err error
	if !r.omitAssetContent {
		content, err = themes.ReadStatic(name)
	}
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
	if bundled && (!options.InlineAssets || options.OmitAssetContent) {
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
