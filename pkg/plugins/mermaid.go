// Package plugins provides lifecycle plugins for markata-go.
package plugins

import (
	"fmt"
	"html"
	"log"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/WaylonWalker/markata/pkg/buildcache"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

const (
	mermaidModeClient   = "client"
	mermaidModeGo       = "go"
	mermaidModeCLI      = "cli"
	mermaidModeChromium = "chromium"
)

// MermaidPlugin converts Mermaid code blocks into rendered diagrams.
// It runs at the render stage (post_render, after markdown conversion).
type MermaidPlugin struct {
	scriptCollector    func(string) // standalone rendering collects generated scripts as assets
	escapeClientSource bool
	config             models.MermaidConfig
	renderer           mermaidRenderer
	paletteColors      *mermaidPaletteColors // resolved at configure time, nil if no palette
	assetURLs          map[string]string
	cache              *buildcache.Cache
	renderKey          string
}

// NewMermaidPlugin creates a new MermaidPlugin with default settings.
func NewMermaidPlugin() *MermaidPlugin {
	return &MermaidPlugin{
		config: models.NewMermaidConfig(),
	}
}

// Name returns the unique name of the plugin.
func (p *MermaidPlugin) Name() string {
	return "mermaid"
}

// Priority returns the plugin's priority for a given stage.
// This plugin runs after render_markdown (which has default priority 0).
func (p *MermaidPlugin) Priority(stage lifecycle.Stage) int {
	if stage == lifecycle.StageRender {
		return lifecycle.PriorityLate // Run after render_markdown
	}
	return lifecycle.PriorityDefault
}

// Configure reads configuration options for the plugin from config.Extra.
// Configuration is expected under the "mermaid" key.
// It also validates that the selected rendering mode has its dependencies installed.
func (p *MermaidPlugin) Configure(m *lifecycle.Manager) error {
	config := m.Config()
	if config.Extra == nil {
		return nil
	}

	p.assetURLs = assetURLsFromConfig(config)

	// Check for mermaid config in Extra
	mermaidConfig, ok := config.Extra["mermaid"]
	if !ok {
		return nil
	}

	// Handle map configuration
	if cfgMap, ok := mermaidConfig.(map[string]interface{}); ok {
		p.parseMainConfig(cfgMap)
		p.parseCLIConfig(cfgMap)
		p.parseChromiumConfig(cfgMap)
		p.parseGoConfig(cfgMap)
	}

	// Validate the selected mode and check dependencies
	if !p.config.Enabled {
		return nil
	}

	// For pre-render modes, resolve palette colors at build time so
	// mermaid diagrams use the site's color scheme instead of defaults.
	if p.config.Mode != mermaidModeClient && p.config.UseCSSVariables {
		p.paletteColors = resolvePaletteColors(config.Extra)
	}

	return p.validateMode()
}

// parseMainConfig extracts main mermaid configuration from the config map.
func (p *MermaidPlugin) parseMainConfig(cfgMap map[string]interface{}) {
	if enabled, ok := cfgMap["enabled"].(bool); ok {
		p.config.Enabled = enabled
	}
	if mode, ok := cfgMap["mode"].(string); ok && mode != "" {
		p.config.Mode = mode
	}
	if cdnURL, ok := cfgMap["cdn_url"].(string); ok && cdnURL != "" {
		p.config.CDNURL = cdnURL
	}
	if theme, ok := cfgMap["theme"].(string); ok && theme != "" {
		p.config.Theme = theme
	}
	if useCSSVariables, ok := cfgMap["use_css_variables"].(bool); ok {
		p.config.UseCSSVariables = useCSSVariables
	}
	if lightbox, ok := cfgMap["lightbox"].(bool); ok {
		p.config.Lightbox = lightbox
	}
	if selector, ok := cfgMap["lightbox_selector"].(string); ok && selector != "" {
		p.config.LightboxSelector = selector
	}
}

// parseCLIConfig extracts CLI renderer configuration from the config map.
func (p *MermaidPlugin) parseCLIConfig(cfgMap map[string]interface{}) {
	cliCfg, ok := cfgMap["cli"].(map[string]interface{})
	if !ok {
		return
	}

	if p.config.CLIConfig == nil {
		p.config.CLIConfig = &models.CLIRendererConfig{}
	}

	if mmdc, ok := cliCfg["mmdc_path"].(string); ok && mmdc != "" {
		p.config.CLIConfig.MMDCPath = mmdc
	}
	if extraArgs, ok := cliCfg["extra_args"].(string); ok && extraArgs != "" {
		p.config.CLIConfig.ExtraArgs = extraArgs
	}
}

// parseGoConfig reads the opt-in native Go diagram renderer's style settings.
func (p *MermaidPlugin) parseGoConfig(cfgMap map[string]interface{}) {
	cfg, ok := cfgMap["go"].(map[string]interface{})
	if !ok {
		return
	}
	if p.config.GoConfig == nil {
		p.config.GoConfig = &models.GoMermaidRendererConfig{}
	}
	n := p.config.GoConfig
	if font, ok := cfg["font_family"].(string); ok && font != "" {
		n.FontFamily = font
	}
	for _, item := range []struct {
		key    string
		target *float64
	}{
		{"font_size", &n.FontSize},
		{"node_spacing", &n.NodeSpacing},
		{"rank_spacing", &n.RankSpacing},
		{"padding", &n.Padding},
	} {
		if value, ok := mermaidFloat(cfg[item.key]); ok && value > 0 {
			*item.target = value
		}
	}
	if v, ok := cfg["curved_edges"].(bool); ok {
		n.CurvedEdges = v
	}
	if v, ok := cfg["transparent_background"].(bool); ok {
		n.TransparentBackground = v
	}
}

func mermaidFloat(value interface{}) (float64, bool) {
	switch v := value.(type) {
	case float64:
		return v, true
	case float32:
		return float64(v), true
	case int:
		return float64(v), true
	case int64:
		return float64(v), true
	default:
		return 0, false
	}
}

// parseChromiumConfig extracts Chromium renderer configuration from the config map.
func (p *MermaidPlugin) parseChromiumConfig(cfgMap map[string]interface{}) {
	chromCfg, ok := cfgMap["chromium"].(map[string]interface{})
	if !ok {
		return
	}

	if p.config.ChromiumConfig == nil {
		p.config.ChromiumConfig = &models.ChromiumRendererConfig{}
	}

	if browserPath, ok := chromCfg["browser_path"].(string); ok && browserPath != "" {
		p.config.ChromiumConfig.BrowserPath = browserPath
	}
	if timeout, ok := parseIntFromInterface(chromCfg["timeout"]); ok && timeout > 0 {
		p.config.ChromiumConfig.Timeout = timeout
	}
	if maxConcurrent, ok := parseIntFromInterface(chromCfg["max_concurrent"]); ok && maxConcurrent > 0 {
		p.config.ChromiumConfig.MaxConcurrent = maxConcurrent
	}
	if noSandbox, ok := chromCfg["no_sandbox"].(bool); ok {
		p.config.ChromiumConfig.NoSandbox = noSandbox
	}
}

// validateMode validates the selected rendering mode and checks for required dependencies.
func (p *MermaidPlugin) validateMode() error {
	// Validate mode value
	switch p.config.Mode {
	case mermaidModeClient, mermaidModeCLI, mermaidModeChromium, mermaidModeGo:
		// valid modes
	default:
		return models.NewConfigValidationError("mermaid.mode", p.config.Mode,
			"invalid mode: must be 'client', 'go', 'cli', or 'chromium'")
	}

	// Check dependencies for non-client modes
	if p.config.Mode == mermaidModeCLI {
		info := checkCLIDependency(p.config.CLIConfig.MMDCPath)
		if !info.IsInstalled {
			err := models.NewMermaidRenderError("", p.config.Mode, "mmdc binary not found", nil)
			err.Suggestion = info.InstallInstructions + "\n\n" + info.FallbackSuggestion
			return err
		}
	} else if p.config.Mode == mermaidModeChromium {
		info := checkChromiumDependency(p.config.ChromiumConfig.BrowserPath)
		if !info.IsInstalled {
			err := models.NewMermaidRenderError("", p.config.Mode, "browser not found", nil)
			err.Suggestion = info.InstallInstructions + "\n\n" + info.FallbackSuggestion
			return err
		}
	}

	return nil
}

// Render processes mermaid code blocks in the rendered HTML for all posts.
func (p *MermaidPlugin) Render(m *lifecycle.Manager) (err error) {
	if !p.config.Enabled {
		return nil
	}
	p.cache = GetBuildCache(m)
	if p.cache != nil {
		p.cache.ResetMermaidSVGUsage()
	}
	p.renderKey = p.rendererCacheKey()

	if p.config.Mode != mermaidModeClient {
		renderer, createErr := newMermaidRenderer(p.config, p.paletteColors)
		if createErr != nil {
			if p.config.Mode == mermaidModeGo {
				return fmt.Errorf("native Mermaid renderer: %w", createErr)
			}
			p.config.Mode = mermaidModeClient
		} else {
			p.renderer = renderer
			defer func() {
				if closeErr := renderer.close(); closeErr != nil && err == nil {
					err = closeErr
				}
			}()
		}
	}

	posts := m.FilterPosts(func(post *models.Post) bool {
		if post.Skip || post.ArticleHTML == "" {
			return false
		}
		return strings.Contains(post.ArticleHTML, `class="language-mermaid"`) ||
			strings.Contains(post.ArticleHTML, `class="mermaid"`)
	})

	if err := m.ProcessPostsSliceConcurrently(posts, p.processPost); err != nil {
		return err
	}

	if p.config.Lightbox {
		needsLightbox := false
		for _, post := range m.Posts() {
			if post.ArticleHTML == "" {
				continue
			}
			if strings.Contains(post.ArticleHTML, `class="mermaid"`) || strings.Contains(post.ArticleHTML, `class="language-mermaid"`) {
				needsLightbox = true
				break
			}
		}

		if needsLightbox {
			config := m.Config()
			if config.Extra == nil {
				config.Extra = make(map[string]interface{})
			}
			// Only enable GLightbox loading; do NOT set glightbox_options.
			// The mermaid lightbox uses a separate programmatic GLightbox
			// instance (selector: false) and does not need the template's
			// shared instance. Overwriting glightbox_options would clobber
			// settings from image_zoom or other plugins.
			config.Extra["glightbox_enabled"] = true
			config.Extra["glightbox_cdn"] = false
		}
	}

	return nil
}

func (p *MermaidPlugin) renderDiagram(post *models.Post, diagramCode string) (string, error) {
	renderHash := p.diagramRenderHash(diagramCode)
	if p.cache != nil && renderHash != "" {
		if cachedSVG := p.cache.GetCachedMermaidSVG(post.Path, renderHash); cachedSVG != "" {
			return cachedSVG, nil
		}
	}

	svgOutput, err := p.renderer.render(diagramCode)
	if err != nil {
		return "", err
	}

	if p.cache != nil && renderHash != "" {
		//nolint:errcheck // caching is best-effort; rendering already succeeded.
		p.cache.CacheMermaidSVG(post.Path, renderHash, svgOutput)
	}
	return svgOutput, nil
}

func (p *MermaidPlugin) diagramRenderHash(diagramCode string) string {
	if p.config.Mode == mermaidModeClient {
		return ""
	}

	var builder strings.Builder
	builder.WriteString("mermaid-svg-v1\n")
	builder.WriteString("renderer:")
	builder.WriteString(p.renderKey)
	builder.WriteString("mode:")
	builder.WriteString(p.config.Mode)
	builder.WriteString("\ntheme:")
	builder.WriteString(p.config.Theme)
	builder.WriteString("\nuse_css_variables:")
	if p.config.UseCSSVariables {
		builder.WriteString("true")
	} else {
		builder.WriteString("false")
	}
	builder.WriteString("\nlightbox:")
	if p.config.Lightbox {
		builder.WriteString("true")
	} else {
		builder.WriteString("false")
	}
	if p.config.CLIConfig != nil {
		builder.WriteString("\ncli.mmdc_path:")
		builder.WriteString(p.config.CLIConfig.MMDCPath)
		builder.WriteString("\ncli.extra_args:")
		builder.WriteString(p.config.CLIConfig.ExtraArgs)
		writeMermaidCLIFileInputHashes(&builder, p.config.CLIConfig.ExtraArgs)
	}
	if p.config.ChromiumConfig != nil {
		builder.WriteString("\nchromium.browser_path:")
		builder.WriteString(p.config.ChromiumConfig.BrowserPath)
		builder.WriteString("\nchromium.timeout:")
		builder.WriteString(fmt.Sprint(p.config.ChromiumConfig.Timeout))
		builder.WriteString("\nchromium.max_concurrent:")
		builder.WriteString(fmt.Sprint(p.config.ChromiumConfig.MaxConcurrent))
		builder.WriteString("\nchromium.no_sandbox:")
		if p.config.ChromiumConfig.NoSandbox {
			builder.WriteString("true")
		} else {
			builder.WriteString("false")
		}
	}
	if p.config.Mode == mermaidModeGo {
		builder.WriteString("\ngo.version:")
		builder.WriteString(nativeMermaidVersion)
		builder.WriteString("\ngo.config:")
		builder.WriteString(fmt.Sprintf("%+v", p.config.GoConfig))
	}
	if p.paletteColors != nil {
		builder.WriteString("\npalette:")
		builder.WriteString(p.paletteColors.themeVariablesJSON())
	}
	builder.WriteString("\ndiagram:\n")
	builder.WriteString(diagramCode)
	return buildcache.ContentHash(builder.String())
}

var mermaidCLIFileArgFlags = map[string]struct{}{
	"--configFile":          {},
	"--cssFile":             {},
	"--puppeteerConfigFile": {},
	"-c":                    {},
	"-C":                    {},
	"-p":                    {},
}

func writeMermaidCLIFileInputHashes(builder *strings.Builder, extraArgs string) {
	if extraArgs == "" {
		return
	}

	args := strings.Fields(extraArgs)
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if _, ok := mermaidCLIFileArgFlags[arg]; ok {
			if i+1 < len(args) {
				writeMermaidCLIFileInputHash(builder, args[i+1])
				i++
			}
			continue
		}

		for flag := range mermaidCLIFileArgFlags {
			prefix := flag + "="
			if strings.HasPrefix(arg, prefix) {
				writeMermaidCLIFileInputHash(builder, strings.TrimPrefix(arg, prefix))
				break
			}
		}
	}
}

func writeMermaidCLIFileInputHash(builder *strings.Builder, path string) {
	if path == "" {
		return
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		absPath = filepath.Clean(path)
	}

	builder.WriteString("\ncli.file_arg:")
	builder.WriteString(absPath)
	builder.WriteString(":")
	if hash, err := buildcache.HashFile(absPath); err == nil && hash != "" {
		builder.WriteString(hash)
		return
	}
	builder.WriteString("unreadable")
}

func (p *MermaidPlugin) rendererCacheKey() string {
	if p.config.Mode == mermaidModeClient {
		return ""
	}

	switch p.config.Mode {
	case mermaidModeGo:
		return nativeMermaidVersion
	case mermaidModeCLI:
		info := checkCLIDependency(p.config.CLIConfig.MMDCPath)
		if !info.IsInstalled {
			return "cli:unavailable"
		}
		return "cli:" + info.BinaryPath + ":" + info.InstalledVersion
	case mermaidModeChromium:
		info := checkChromiumDependency(p.config.ChromiumConfig.BrowserPath)
		if !info.IsInstalled {
			return "chromium:unavailable"
		}
		return "chromium:" + info.BinaryPath + ":" + info.InstalledVersion + ":mermaidjs-v" + mermaidJSVersion
	default:
		return p.config.Mode
	}
}

// mermaidCodeBlockRegex matches <pre><code class="language-mermaid"> blocks.
// It captures the diagram code inside.
var mermaidCodeBlockRegex = regexp.MustCompile(
	`<pre><code class="language-mermaid"[^>]*>([\s\S]*?)</code></pre>`,
)

// processPost processes a single post's HTML for mermaid code blocks.
// For client mode: converts code blocks to mermaid pre tags and injects script
// For cli/chromium modes: renders diagrams to SVGs and embeds them
func (p *MermaidPlugin) processPost(post *models.Post) error {
	// Skip posts marked as skip or with no HTML content
	if post.Skip || post.ArticleHTML == "" {
		return nil
	}

	// Check if there are any mermaid code blocks (language-mermaid needs conversion)
	// or pre-rendered mermaid blocks (class="mermaid" just needs the script)
	hasLanguageMermaid := strings.Contains(post.ArticleHTML, `class="language-mermaid"`)
	hasPreRendered := strings.Contains(post.ArticleHTML, `class="mermaid"`)
	if !hasLanguageMermaid && !hasPreRendered {
		return nil
	}

	// Track if we found any mermaid blocks
	foundMermaid := false
	result := post.ArticleHTML

	// Replace language-mermaid code blocks with proper mermaid pre tags or rendered SVGs
	renderer := p.renderer
	mode := p.config.Mode
	if mode != mermaidModeClient && renderer == nil {
		mode = mermaidModeClient
	}
	var renderErr error
	if hasLanguageMermaid {
		result = mermaidCodeBlockRegex.ReplaceAllStringFunc(result, func(match string) string {
			foundMermaid = true
			if renderErr != nil {
				return match
			}

			// Extract the diagram code
			submatches := mermaidCodeBlockRegex.FindStringSubmatch(match)
			if len(submatches) < 2 {
				return match
			}

			// Decode HTML entities in the diagram code (goldmark encodes them)
			diagramCode := html.UnescapeString(submatches[1])

			// Trim whitespace from the diagram code
			diagramCode = strings.TrimSpace(diagramCode)

			// For client mode: return as mermaid pre block
			if mode == mermaidModeClient {
				if p.escapeClientSource {
					diagramCode = html.EscapeString(diagramCode)
				}
				return `<pre class="mermaid">` + "\n" + diagramCode + "\n</pre>"
			}

			// For pre-rendering modes (cli/chromium): render to SVG
			svgOutput, err := p.renderDiagram(post, diagramCode)
			if err != nil {
				log.Printf("[mermaid] render error for %s: %v", post.Path, err)
				renderErr = models.NewMermaidRenderError(post.Path, p.config.Mode, "failed to render diagram", err)
				return match
			}

			// Return the SVG wrapped in a mermaid container
			return `<pre class="mermaid">` + "\n" + svgOutput + "\n</pre>"
		})
	}

	if renderErr != nil {
		return renderErr
	}

	// If we found mermaid blocks, inject appropriate scripts
	if foundMermaid || strings.Contains(result, `class="mermaid"`) {
		if mode == mermaidModeClient {
			result = p.injectMermaidScript(result)
		} else if p.config.Lightbox {
			// Pre-rendered modes (cli/chromium): SVGs are already in the HTML,
			// but we need the lightbox JS to wire up click-to-zoom.
			result = p.injectPrerenderedLightboxScript(result)
		}

		// Signal that this post needs GLightbox so the base template
		// loads the GLightbox JS/CSS. The template condition is:
		//   config.Extra.glightbox_enabled AND needs_image_zoom
		// image_zoom sets this for zoomable images; we reuse the same
		// flag for mermaid diagrams that also need lightbox zoom.
		if p.config.Lightbox {
			if post.Extra == nil {
				post.Extra = make(map[string]interface{})
			}
			post.Extra["needs_image_zoom"] = true
		}
	}

	post.ArticleHTML = result
	return nil
}

// injectPrerenderedLightboxScript adds click-to-lightbox behavior for
// pre-rendered (cli/chromium) SVG diagrams. Unlike the client-mode script,
// this does NOT import or initialize MermaidJS -- the SVGs are already
// rendered into the HTML at build time.
func (p *MermaidPlugin) injectPrerenderedLightboxScript(htmlContent string) string {
	script := `
<script>
(function() {
` + p.mermaidLightboxJS() + `
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', () => ensureMermaidLightbox());
  } else {
    ensureMermaidLightbox();
  }
})();
</script>`
	return htmlContent + script
}

// injectMermaidScript adds the Mermaid.js initialization script to the HTML.
// The script is only injected once per post.
func (p *MermaidPlugin) injectMermaidScript(htmlContent string) string {
	var script string
	if p.config.UseCSSVariables {
		script = p.cssVariablesScript()
	} else {
		script = `
<script type="module">
` + p.mermaidImportJS() + p.mermaidLightboxJS() + `
  mermaid.initialize({ startOnLoad: false, theme: '` + p.config.Theme + `' });
  window.initMermaid = async () => {
    try {
      await mermaid.run();
    } catch (e) {
      console.error('mermaid.run failed:', e);
    }
    ensureMermaidLightbox();
  };
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', () => window.initMermaid());
  } else {
    window.initMermaid();
  }
</script>`
	}

	// Append the script to the end of the content
	if p.scriptCollector != nil {
		p.scriptCollector(script)
		return htmlContent
	}
	return htmlContent + script
}

func (p *MermaidPlugin) cssVariablesScript() string {
	return `
<script type="module">
` + p.mermaidImportJS() + `
  const diagramSources = new WeakMap();
  let renderQueue = Promise.resolve();
  const renderDiagrams = async (root = document) => {
    const roots = Array.isArray(root) ? root : [root];
    const blocks = roots.flatMap(node => node.matches?.('pre.mermaid') ? [node] : [...(node.querySelectorAll?.('pre.mermaid') || [])]);
    if (blocks.length === 0) return;

    for (const block of blocks) {
      if (!diagramSources.has(block)) {
        const source = block.textContent || '';
        diagramSources.set(block, source);
        block.dataset.markataMermaidSource = source;
      }
      block.removeAttribute('data-processed');
      block.textContent = diagramSources.get(block);
    }

    const css = (name, fallback) => {
      const probe = document.createElement('span');
      probe.style.backgroundColor = ` + "`" + `var(${name}, ${fallback})` + "`" + `;
      document.documentElement.appendChild(probe);
      const value = getComputedStyle(probe).backgroundColor;
      probe.remove();
      // Mermaid 10 does not accept color-mix(), color(), or other CSS color
      // functions. Resolve the active CSS color through canvas to an RGB value.
      const canvas = document.createElement('canvas');
      const context = canvas.getContext('2d');
      if (!context) return value;
      context.fillStyle = value;
      context.fillRect(0, 0, 1, 1);
      const [red, green, blue] = context.getImageData(0, 0, 1, 1).data;
      return ` + "`" + `rgb(${red}, ${green}, ${blue})` + "`" + `;
    };
    // The active palette is the source of truth: OS preference may differ
    // from the user-selected theme.
    const background = css('--color-background', '#ffffff');
    const accent = css('--color-primary', '#64748b');
    const surface = css('--color-surface', background);
    const foreground = css('--color-text', '#1f2937');
    const border = css('--color-border', accent);
    const flowchart = { nodeSpacing: 60, rankSpacing: 90, padding: 12 };
    const themeCSS = ` + "`" + `
    .label foreignObject > div { padding: 14px 14px 10px; line-height: 1.2; }
    .nodeLabel { padding: 14px 14px 10px; line-height: 1.2; }

  ` + "`" + `;
    const themeVariables = {
      background,
      primaryColor: surface,
      primaryTextColor: foreground,
      primaryBorderColor: border,
      // Diagram edges encode direction and relationships, so keep them at
      // foreground contrast instead of the intentionally subtle node border.
      lineColor: foreground,
      textColor: foreground,
      nodeBkg: surface,
      nodeBorder: border,
      nodeTextColor: foreground,
      fontSize: '16px',
      nodePadding: 20,
      nodeTextMargin: 14,
      clusterBkg: background,
      clusterBorder: border,
      clusterTextColor: foreground,
      titleColor: foreground,
      edgeLabelBackground: surface,
    };

    mermaid.initialize({ startOnLoad: false, theme: 'base', themeVariables, flowchart, themeCSS });
    await mermaid.run({ nodes: blocks });
    ensureMermaidLightbox();
  };

  window.initMermaid = (root = document) => {
    renderQueue = renderQueue.then(() => renderDiagrams(root)).catch(e => {
      console.error('mermaid.run failed:', e);
    });
    return renderQueue;
  };
  const rerenderMermaid = () => window.initMermaid();
  window.addEventListener('theme-change', rerenderMermaid);
  window.addEventListener('palette-change', rerenderMermaid);
  window.addEventListener('color-mode-change', rerenderMermaid);
  window.addEventListener('view-transition-complete', rerenderMermaid);
` + p.mermaidLightboxJS() + `
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', () => window.initMermaid());
  } else {
    window.initMermaid();
  }
</script>`
}

// mermaidScriptURL returns the URL the client-mode script loads Mermaid from.
// A self-hosted asset URL (from the cdn_assets plugin) takes precedence over
// the configured cdn_url.
func (p *MermaidPlugin) mermaidScriptURL() string {
	if url := p.resolveAssetURL("mermaid", ""); url != "" {
		return url
	}
	return p.resolveAssetURL("mermaid-esm", p.config.CDNURL)
}

// isESModuleURL reports whether a Mermaid URL points at an ES module build.
// The default vendored bundle is the self-contained UMD build; the ESM build
// lazily imports per-diagram chunks and only works when the full dist
// directory is served alongside it (e.g. from a CDN).
func isESModuleURL(url string) bool {
	path := url
	if i := strings.IndexAny(path, "?#"); i >= 0 {
		path = path[:i]
	}
	return strings.HasSuffix(path, ".mjs") || strings.Contains(path, ".esm.")
}

// mermaidImportJS returns JavaScript that binds a `mermaid` constant inside a
// module script, using a static import for ESM URLs and a classic script tag
// (resolving to window.mermaid) for UMD bundles.
func (p *MermaidPlugin) mermaidImportJS() string {
	url := p.mermaidScriptURL()
	if isESModuleURL(url) {
		return `  import mermaid from '` + url + `';
`
	}
	return `  const mermaid = await (async () => {
    if (window.mermaid) return window.mermaid;
    await new Promise((resolve, reject) => {
      const s = document.createElement('script');
      s.src = '` + url + `';
      s.onload = resolve;
      s.onerror = () => reject(new Error('failed to load mermaid from ' + s.src));
      document.head.appendChild(s);
    });
    return window.mermaid;
  })();
`
}

// mermaidLightboxJS returns the shared JavaScript for mermaid diagram lightbox
// with svg-pan-zoom support. This is used by both the standard and CSS variables
// code paths.
func (p *MermaidPlugin) mermaidLightboxJS() string {
	svgPanZoomURL := p.resolveAssetURL("svg-pan-zoom", "/assets/vendor/svg-pan-zoom/svg-pan-zoom.min.js")
	return `
  const SVG_PAN_ZOOM_CDN = '` + svgPanZoomURL + `';
  let mermaidLightbox = null;
  let activePanZoom = null;

  // Inject lightbox styles once
  const injectLightboxStyles = () => {
    if (document.getElementById('mermaid-lightbox-css')) return;
    const style = document.createElement('style');
    style.id = 'mermaid-lightbox-css';
    style.textContent = ` + "`" + `
      /* Container fills the GLightbox slide */
      .mermaid-lightbox-wrap {
        width: 100%;
        height: 100%;
        display: flex;
        align-items: center;
        justify-content: center;
        background: transparent;
        position: relative;
      }
      .mermaid-lightbox-wrap svg {
        width: 100% !important;
        height: 100% !important;
        max-width: 100%;
        max-height: 100%;
      }
      /* Hide GLightbox prev/next arrows (single-slide lightbox) */
      .glightbox-container .gprev,
      .glightbox-container .gnext {
        display: none !important;
      }
      /* Hide description area that renders as a white box */
      .glightbox-container .gslide-description,
      .glightbox-container .gslide-title,
      .glightbox-container .gdesc-inner,
      .glightbox-container .gslide-desc {
        display: none !important;
      }
      /* Remove white background from inline slide content */
      .glightbox-container .gslide-inline {
        background: transparent !important;
      }
      /* Make the inline content area fill the slide */
      .glightbox-container .ginlined-content {
        max-width: none !important;
        max-height: none !important;
        width: 100%;
        height: 100%;
        padding: 0 !important;
      }
      /* Remove box-shadow from the media container */
      .glightbox-container .gslide-media {
        box-shadow: none !important;
      }
      /* Toolbar styling */
      .mermaid-lightbox-toolbar {
        position: absolute;
        top: 8px;
        right: 8px;
        z-index: 10;
        display: flex;
        gap: 4px;
      }
      .mermaid-pz-btn {
        background: rgba(0,0,0,0.6);
        color: #fff;
        border: 1px solid rgba(255,255,255,0.3);
        border-radius: 4px;
        padding: 4px 10px;
        cursor: pointer;
        font-size: 14px;
        line-height: 1;
      }
      .mermaid-pz-btn:hover {
        background: rgba(0,0,0,0.8);
        border-color: rgba(255,255,255,0.6);
      }
    ` + "`" + `;
    document.head.appendChild(style);
  };

  // Lazy-load svg-pan-zoom from CDN, returns a promise
  const loadSvgPanZoom = () => {
    if (typeof svgPanZoom !== 'undefined') return Promise.resolve();
    return new Promise((resolve, reject) => {
      const s = document.createElement('script');
      s.src = SVG_PAN_ZOOM_CDN;
      s.onload = resolve;
      s.onerror = reject;
      document.head.appendChild(s);
    });
  };

  // Initialize svg-pan-zoom on the SVG inside the lightbox.
  // Retries until the lightbox container has settled dimensions.
  let _pzRetries = 0;
  const initPanZoom = () => {
    if (activePanZoom) return;
    const container = document.querySelector('.glightbox-container .gslide.current .mermaid-lightbox-wrap');
    if (!container) return;
    const svgEl = container.querySelector('svg');
    if (!svgEl) return;

    // Ensure the container has layout dimensions before initializing.
    const cRect = container.getBoundingClientRect();
    if (cRect.width < 10 || cRect.height < 10) {
      if (_pzRetries < 20) { _pzRetries++; setTimeout(initPanZoom, 50); }
      return;
    }

    // svg-pan-zoom needs a viewBox. Pre-rendered SVGs from mermaid
    // usually have one; browser-rendered ones may not.
    if (!svgEl.getAttribute('viewBox')) {
      let w = parseFloat(svgEl.getAttribute('width'));
      let h = parseFloat(svgEl.getAttribute('height'));
      if (!w && svgEl.style.maxWidth) w = parseFloat(svgEl.style.maxWidth);
      if (!w || !h) {
        const r = svgEl.getBoundingClientRect();
        if (!w) w = r.width;
        if (!h) h = r.height;
      }
      if (w > 0 && h > 0) {
        svgEl.setAttribute('viewBox', '0 0 ' + w + ' ' + h);
      } else if (_pzRetries < 20) {
        _pzRetries++; setTimeout(initPanZoom, 50); return;
      }
    }
    _pzRetries = 0;

    // Clear inline dimensions so SVG can be sized by the container
    // and svg-pan-zoom can manage transforms.
    svgEl.removeAttribute('width');
    svgEl.removeAttribute('height');
    svgEl.style.cssText = 'width:100%;height:100%;';

    try {
      activePanZoom = svgPanZoom(svgEl, {
        zoomEnabled: true,
        panEnabled: true,
        controlIconsEnabled: false,
        fit: true,
        center: true,
        contain: false,
        minZoom: 0.3,
        maxZoom: 10,
        zoomScaleSensitivity: 0.3,
        mouseWheelZoomEnabled: true,
        preventMouseEventsDefault: true,
      });
      // Double-check fit after a frame in case dimensions shifted
      requestAnimationFrame(() => {
        if (!activePanZoom) return;
        activePanZoom.resize();
        activePanZoom.fit();
        activePanZoom.center();
      });
    } catch (_) {
      activePanZoom = null;
    }

    // Add reset/fit buttons
    let toolbar = container.querySelector('.mermaid-lightbox-toolbar');
    if (!toolbar) {
      toolbar = document.createElement('div');
      toolbar.className = 'mermaid-lightbox-toolbar';
      toolbar.innerHTML =
        '<button class="mermaid-pz-btn" data-action="fit" title="Fit to view">Fit</button>' +
        '<button class="mermaid-pz-btn" data-action="zoomin" title="Zoom in">+</button>' +
        '<button class="mermaid-pz-btn" data-action="zoomout" title="Zoom out">&minus;</button>';
      toolbar.addEventListener('click', (ev) => {
        const btn = ev.target.closest('[data-action]');
        if (!btn || !activePanZoom) return;
        ev.preventDefault();
        ev.stopPropagation();
        const action = btn.dataset.action;
        if (action === 'fit') { activePanZoom.resize(); activePanZoom.fit(); activePanZoom.center(); }
        else if (action === 'zoomin') { activePanZoom.zoomIn(); }
        else if (action === 'zoomout') { activePanZoom.zoomOut(); }
      });
      container.prepend(toolbar);
    }
  };

  // Destroy pan-zoom on lightbox close
  const destroyPanZoom = () => {
    if (activePanZoom) {
      try { activePanZoom.destroy(); } catch (_) { /* no-op */ }
      activePanZoom = null;
    }
  };

  let _lbRetries = 0;
  const ensureMermaidLightbox = () => {
    const diagrams = document.querySelectorAll('.mermaid svg');
    if (!diagrams.length) {
      // Mermaid ESM may still be rendering -- retry up to 2s
      if (_lbRetries < 20) { _lbRetries++; setTimeout(ensureMermaidLightbox, 100); }
      return;
    }
    _lbRetries = 0;
    injectLightboxStyles();
    diagrams.forEach((svg) => {
      if (svg.dataset.lightboxBound) return;
      svg.dataset.lightboxBound = 'true';
      svg.style.cursor = 'pointer';
      svg.addEventListener('click', (e) => {
        e.preventDefault();
        e.stopPropagation();
        const svgHtml = svg.outerHTML;
        const openLightbox = () => {
          if (!mermaidLightbox) {
            mermaidLightbox = GLightbox({
              selector: false,
              openEffect: 'fade',
              closeEffect: 'fade',
              zoomable: false,
              draggable: false,
              skin: 'clean',
            });
            mermaidLightbox.on('slide_after_load', () => {
              destroyPanZoom();
              _pzRetries = 0;
              loadSvgPanZoom().then(() => initPanZoom());
            });
            mermaidLightbox.on('close', destroyPanZoom);
          }
          mermaidLightbox.setElements([{
            content: '<div class="mermaid-lightbox-wrap">' + svgHtml + '</div>',
            width: '90vw',
            height: '90vh'
          }]);
          mermaidLightbox.open();
          loadSvgPanZoom();
        };
        if (typeof GLightbox !== 'undefined') {
          openLightbox();
        } else if (window.initGLightbox) {
          window.initGLightbox();
          openLightbox();
        } else {
          window.addEventListener('glightbox-ready', () => { openLightbox(); }, { once: true });
        }
      });
    });
  };
`
}

func (p *MermaidPlugin) resolveAssetURL(name, fallback string) string {
	return resolveAssetURL(p.assetURLs, name, fallback)
}

// SetConfig sets the mermaid configuration directly.
// This is useful for testing or programmatic configuration.
func (p *MermaidPlugin) SetConfig(config models.MermaidConfig) {
	p.config = config
}

// Config returns the current mermaid configuration.
func (p *MermaidPlugin) Config() models.MermaidConfig {
	return p.config
}

// Ensure MermaidPlugin implements the required interfaces.
var (
	_ lifecycle.Plugin          = (*MermaidPlugin)(nil)
	_ lifecycle.ConfigurePlugin = (*MermaidPlugin)(nil)
	_ lifecycle.RenderPlugin    = (*MermaidPlugin)(nil)
	_ lifecycle.PriorityPlugin  = (*MermaidPlugin)(nil)
)
