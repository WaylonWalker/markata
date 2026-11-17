package plugins

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"sync"
	"unicode/utf8"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/visualizations"
)

var excalidrawCodeBlockPattern = regexp.MustCompile(`(?s)<pre><code class="language-excalidraw"[^>]*>(.*?)</code></pre>`)

type excalidrawSVGRenderer interface {
	visualizations.Renderer
	Close() error
}

// ExcalidrawPlugin renders Excalidraw scene fences to published static SVG.
type ExcalidrawPlugin struct {
	config models.ExcalidrawConfig
	assets []visualizations.Asset
}

// NewExcalidrawPlugin creates the static Excalidraw renderer.
func NewExcalidrawPlugin() *ExcalidrawPlugin {
	return &ExcalidrawPlugin{config: models.NewExcalidrawConfig()}
}

func (p *ExcalidrawPlugin) Name() string { return "excalidraw" }

func (p *ExcalidrawPlugin) Priority(stage lifecycle.Stage) int {
	switch stage {
	case lifecycle.StageRender:
		return lifecycle.PriorityLate
	case lifecycle.StageWrite:
		return lifecycle.PriorityEarly
	default:
		return lifecycle.PriorityDefault
	}
}

func (p *ExcalidrawPlugin) Configure(m *lifecycle.Manager) error {
	p.config = parseExcalidrawConfig(m.Config().Extra)
	return nil
}

func (p *ExcalidrawPlugin) Render(m *lifecycle.Manager) (err error) {
	p.assets = nil
	if !p.config.Enabled {
		return nil
	}
	posts := m.FilterPosts(func(post *models.Post) bool {
		return !post.Skip && !post.Draft && !post.Private && excalidrawCodeBlockPattern.MatchString(post.ArticleHTML)
	})
	if len(posts) == 0 {
		return nil
	}

	renderer, err := newExcalidrawRenderer(p.config)
	if err != nil {
		return fmt.Errorf("initialize Excalidraw renderer: %w", err)
	}
	defer func() {
		if closeErr := renderer.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("close Excalidraw renderer: %w", closeErr)
		}
	}()

	var pipeline visualizations.Pipeline
	siteURL := ""
	if configured, ok := m.Config().Extra["models_config"].(*models.Config); ok {
		siteURL = configured.URL
	}
	basePath, err := excalidrawAssetBasePath(siteURL)
	if err != nil {
		return fmt.Errorf("resolve Excalidraw asset base path: %w", err)
	}
	var assetsMu sync.Mutex
	assetsByPath := make(map[string]visualizations.Asset)
	if err := m.ProcessPostsSliceConcurrently(posts, func(post *models.Post) error {
		updated, assets, renderErr := replaceExcalidrawBlocks(post.ArticleHTML, func(sceneSource string) (visualizations.Asset, error) {
			return pipeline.Compile(context.Background(), renderer, visualizations.Request{Source: sceneSource})
		}, func(asset visualizations.Asset, alt string) (string, error) {
			return visualizations.FigureHTMLWithBasePath(asset, alt, "", false, basePath)
		})
		if renderErr != nil {
			return fmt.Errorf("post %s: %w", post.Path, renderErr)
		}
		post.ArticleHTML = updated
		assetsMu.Lock()
		for _, asset := range assets {
			assetsByPath[asset.Path] = asset
		}
		assetsMu.Unlock()
		return nil
	}); err != nil {
		return fmt.Errorf("render Excalidraw scenes: %w", err)
	}

	paths := make([]string, 0, len(assetsByPath))
	for assetPath := range assetsByPath {
		paths = append(paths, assetPath)
	}
	sort.Strings(paths)
	for _, assetPath := range paths {
		p.assets = append(p.assets, assetsByPath[assetPath])
	}
	return nil
}

func excalidrawAssetBasePath(siteURL string) (string, error) {
	if strings.TrimSpace(siteURL) == "" {
		return "", nil
	}
	parsed, err := url.Parse(siteURL)
	if err != nil {
		return "", err
	}
	return strings.TrimRight(parsed.EscapedPath(), "/"), nil
}

func (p *ExcalidrawPlugin) Write(m *lifecycle.Manager) error {
	for _, asset := range p.assets {
		if _, err := visualizations.Publish(m.Config().OutputDir, asset); err != nil {
			return fmt.Errorf("publish Excalidraw SVG %s: %w", asset.Path, err)
		}
	}
	return nil
}

func parseExcalidrawConfig(extra map[string]any) models.ExcalidrawConfig {
	config := models.NewExcalidrawConfig()
	if extra == nil {
		return config
	}
	raw, ok := extra["excalidraw"]
	if !ok {
		return config
	}
	if typed, ok := raw.(models.ExcalidrawConfig); ok {
		config = typed
		return normalizeExcalidrawConfig(config)
	}
	values, ok := raw.(map[string]any)
	if !ok {
		return config
	}
	if enabled, ok := values["enabled"].(bool); ok {
		config.Enabled = enabled
	}
	if chromium, ok := values["chromium"].(map[string]any); ok {
		if browserPath, ok := chromium["browser_path"].(string); ok {
			config.Chromium.BrowserPath = browserPath
		}
		if timeout, ok := parseIntFromInterface(chromium["timeout"]); ok {
			config.Chromium.Timeout = timeout
		}
		if noSandbox, ok := chromium["no_sandbox"].(bool); ok {
			config.Chromium.NoSandbox = noSandbox
		}
	}
	return normalizeExcalidrawConfig(config)
}

func normalizeExcalidrawConfig(config models.ExcalidrawConfig) models.ExcalidrawConfig {
	defaults := models.NewExcalidrawConfig()
	if config.Chromium.Timeout <= 0 {
		config.Chromium.Timeout = defaults.Chromium.Timeout
	}
	return config
}

type excalidrawScene struct {
	Type     string              `json:"type"`
	Version  int                 `json:"version"`
	Elements []excalidrawElement `json:"elements"`
}

type excalidrawElement struct {
	Type      string `json:"type"`
	Text      string `json:"text"`
	Link      string `json:"link"`
	IsDeleted bool   `json:"isDeleted"`
}

func validateExcalidrawScene(source string) (string, error) {
	var scene excalidrawScene
	if err := json.Unmarshal([]byte(source), &scene); err != nil {
		return "", fmt.Errorf("invalid scene JSON: %w", err)
	}
	if scene.Type != "excalidraw" || scene.Version != 2 || scene.Elements == nil {
		return "", errors.New("expected an Excalidraw version 2 scene with an elements array")
	}
	texts := make([]string, 0)
	activeCount := 0
	for index, element := range scene.Elements {
		if element.IsDeleted {
			continue
		}
		activeCount++
		if strings.TrimSpace(element.Link) != "" {
			return "", fmt.Errorf("element %d has a link; linked elements are not supported in static SVG yet", index+1)
		}
		switch element.Type {
		case "rectangle", "diamond", "ellipse", "line", "arrow", "freedraw", "text", "frame", "magicframe":
		case "image":
			return "", fmt.Errorf("element %d is an embedded image; image elements are not supported in static SVG yet", index+1)
		case "iframe", "embeddable":
			return "", fmt.Errorf("element %d is an embedded web element; embeds are not supported in static SVG yet", index+1)
		default:
			return "", fmt.Errorf("element %d has unsupported Excalidraw type %q", index+1, element.Type)
		}
		if element.Type == "text" && strings.TrimSpace(element.Text) != "" {
			texts = append(texts, strings.Join(strings.Fields(element.Text), " "))
		}
	}
	alt := strings.Join(texts, "; ")
	if alt == "" {
		alt = fmt.Sprintf("Excalidraw diagram with %d elements", activeCount)
	}
	if utf8.RuneCountInString(alt) > 180 {
		runes := []rune(alt)
		alt = string(runes[:177]) + "…"
	}
	return alt, nil
}

func replaceExcalidrawBlocks(
	article string,
	render func(string) (visualizations.Asset, error),
	figure func(visualizations.Asset, string) (string, error),
) (string, []visualizations.Asset, error) {
	matches := excalidrawCodeBlockPattern.FindAllStringSubmatchIndex(article, -1)
	if len(matches) == 0 {
		return article, nil, nil
	}
	var output strings.Builder
	output.Grow(len(article))
	assets := make([]visualizations.Asset, 0, len(matches))
	last := 0
	for _, match := range matches {
		output.WriteString(article[last:match[0]])
		sceneSource := html.UnescapeString(article[match[2]:match[3]])
		alt, err := validateExcalidrawScene(sceneSource)
		if err != nil {
			return "", nil, fmt.Errorf("Excalidraw block: %w", err)
		}
		asset, err := render(sceneSource)
		if err != nil {
			return "", nil, fmt.Errorf("render scene: %w", err)
		}
		markup, err := figure(asset, alt)
		if err != nil {
			return "", nil, fmt.Errorf("create static figure: %w", err)
		}
		output.WriteString(markup)
		assets = append(assets, asset)
		last = match[1]
	}
	output.WriteString(article[last:])
	return output.String(), assets, nil
}

var (
	_ lifecycle.Plugin          = (*ExcalidrawPlugin)(nil)
	_ lifecycle.ConfigurePlugin = (*ExcalidrawPlugin)(nil)
	_ lifecycle.RenderPlugin    = (*ExcalidrawPlugin)(nil)
	_ lifecycle.WritePlugin     = (*ExcalidrawPlugin)(nil)
	_ lifecycle.PriorityPlugin  = (*ExcalidrawPlugin)(nil)
)
