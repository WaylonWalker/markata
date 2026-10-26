// Package plugins provides lifecycle plugins for markata-go.
package plugins

import (
	"encoding/hex"
	"fmt"
	"hash/fnv"
	"html"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/WaylonWalker/markata/pkg/assets"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/runtimeenv"
)

const (
	webAwesomeDefaultVersion = "3.5.0"
	webAwesomeDefaultCDNBase = "https://cdn.jsdelivr.net/npm/@awesome.me/webawesome@3.5.0/dist-cdn"
	webAwesomeAssetName      = "webawesome"
	webAwesomeSourceCDN      = "cdn"
	webAwesomeSourceVendor   = "vendor"
	webAwesomeComponentTag   = "tag"
)

// WebAwesomePlugin converts ergonomic markdown containers into Web Awesome components.
type WebAwesomePlugin struct {
	config webAwesomeConfig
}

type webAwesomeConfig struct {
	Enabled   bool
	Version   string
	Source    string
	CDNBase   string
	OutputDir string
	Theme     string
	Palette   string
	Brand     string
}

// NewWebAwesomePlugin creates a Web Awesome integration plugin.
func NewWebAwesomePlugin() *WebAwesomePlugin {
	return &WebAwesomePlugin{config: defaultWebAwesomeConfig()}
}

func defaultWebAwesomeConfig() webAwesomeConfig {
	return webAwesomeConfig{
		Enabled:   true,
		Version:   webAwesomeDefaultVersion,
		Source:    webAwesomeSourceVendor,
		CDNBase:   webAwesomeDefaultCDNBase,
		OutputDir: "assets/vendor/webawesome",
		Theme:     "default",
		Palette:   "default",
		Brand:     "blue",
	}
}

// Name returns the unique plugin name.
func (p *WebAwesomePlugin) Name() string {
	return "webawesome"
}

// Priority runs after markdown rendering but before templates wrap pages.
func (p *WebAwesomePlugin) Priority(stage lifecycle.Stage) int {
	if stage == lifecycle.StageConfigure {
		return lifecycle.PriorityFirst
	}
	if stage == lifecycle.StageRender {
		return 50
	}
	return lifecycle.PriorityDefault
}

// Configure reads [markata-go.webawesome] from config.Extra.
func (p *WebAwesomePlugin) Configure(m *lifecycle.Manager) error {
	config := m.Config()
	sourceExplicit := false
	if config.Extra != nil {
		webawesomeConfig, ok := config.Extra["webawesome"]
		if ok {
			cfgMap, ok := webawesomeConfig.(map[string]interface{})
			if ok {
				sourceExplicit = p.parseConfigMap(cfgMap)
			}
		}
	}
	if runtimeenv.OfflineEnabled() && !sourceExplicit && !hasExplicitSelfHostedAssetsConfig(config) {
		p.config.Source = webAwesomeSourceCDN
	}
	p.enableVendorAsset(config)

	return nil
}

func (p *WebAwesomePlugin) parseConfigMap(cfgMap map[string]interface{}) bool {
	sourceExplicit := false
	if enabled, ok := cfgMap["enabled"].(bool); ok {
		p.config.Enabled = enabled
	}
	if version, ok := cfgMap["version"].(string); ok && version != "" {
		p.config.Version = version
		p.config.CDNBase = "https://cdn.jsdelivr.net/npm/@awesome.me/webawesome@" + version + "/dist-cdn"
	}
	if source, ok := cfgMap["source"].(string); ok && source != "" {
		p.config.Source = source
		sourceExplicit = true
	}
	if cdnBase, ok := cfgMap["cdn_base_url"].(string); ok && cdnBase != "" {
		p.config.CDNBase = strings.TrimRight(cdnBase, "/")
	}
	if outputDir, ok := cfgMap["output_dir"].(string); ok && outputDir != "" {
		p.config.OutputDir = strings.Trim(outputDir, "/")
	}
	if theme, ok := cfgMap["theme"].(string); ok && theme != "" {
		p.config.Theme = theme
	}
	if palette, ok := cfgMap["palette"].(string); ok && palette != "" {
		p.config.Palette = palette
	}
	if brand, ok := cfgMap["brand"].(string); ok && brand != "" {
		p.config.Brand = brand
	}
	return sourceExplicit
}

func hasExplicitSelfHostedAssetsConfig(config *lifecycle.Config) bool {
	if config == nil || config.Extra == nil {
		return false
	}

	switch assetsConfig := config.Extra["assets"].(type) {
	case models.AssetsConfig:
		return assetsConfig.IsSelfHosted()
	case *models.AssetsConfig:
		return assetsConfig != nil && assetsConfig.IsSelfHosted()
	case map[string]interface{}:
		mode, ok := assetsConfig["mode"].(string)
		return ok && (mode == "self-hosted" || mode == "auto")
	}
	return false
}

// Render processes Web Awesome markdown containers and marks pages that need assets.
func (p *WebAwesomePlugin) Render(m *lifecycle.Manager) error {
	if !p.config.Enabled {
		return nil
	}

	posts := m.FilterPosts(func(post *models.Post) bool {
		if post.Skip || post.ArticleHTML == "" {
			return false
		}
		return webAwesomeContainerOpenRegex.MatchString(post.ArticleHTML) ||
			webAwesomeElementRegex.MatchString(post.ArticleHTML)
	})

	if err := m.ProcessPostsSliceConcurrently(posts, p.processPost); err != nil {
		return err
	}

	needsWebAwesome := false
	for _, post := range m.Posts() {
		if post.ArticleHTML != "" && strings.Contains(post.ArticleHTML, "<wa-") {
			needsWebAwesome = true
			break
		}
	}

	if needsWebAwesome {
		p.enableAssets(m.Config())
	}

	return nil
}

var webAwesomeComparisonRegex = regexp.MustCompile(`(?s)<div([^>]*)class="([^"]*(?:\bwebawesome\b[^\"]*\bcomparison\b|\bwa-comparison\b)[^"]*)"([^>]*)>\s*(?:<p>\s*)?(?:<figure>\s*)?(?:<a\b[^>]*>\s*)?(<img\s+[^>]*>)\s*(?:</a>\s*)?(?:<a\b[^>]*>\s*)?(<img\s+[^>]*>)\s*(?:</a>\s*)?(?:</figure>\s*)?(?:</p>\s*)?</div>`)
var webAwesomeAttrRegex = regexp.MustCompile(`([A-Za-z_:][-A-Za-z0-9_:.]*)\s*=\s*("[^"]*"|'[^']*'|[^\s"'>]+)`)
var webAwesomeElementRegex = regexp.MustCompile(`<\s*(wa-[a-z0-9-]+)\b`)
var webAwesomeNestedTabsOpenRegex = regexp.MustCompile(`<div([^>]*)class="([^"]*(?:\bwebawesome\s+tabs\b|\bwa-tabs\b)[^"]*)"([^>]*)>`)
var webAwesomeNestedTabOpenRegex = regexp.MustCompile(`<div([^>]*)class="([^"]*(?:\bwebawesome\s+tab\b|\bwa-tab\b)[^"]*)"([^>]*)>`)
var webAwesomeContainerOpenRegex = regexp.MustCompile(`<div([^>]*)class="([^"]*(?:\bwebawesome\b|\bwa-[a-z0-9-]+\b)[^"]*)"([^>]*)>`)
var webAwesomeTabMarkerRegex = regexp.MustCompile(`(?s)(?:<hr>\s*<p>tab\s+&quot;([^&]+)&quot;</p>|<p>&mdash; tab &ldquo;([^&]+)&rdquo;</p>)\s*`)
var webAwesomeFirstParagraphRegex = regexp.MustCompile(`(?s)^\s*<p>(.*?)</p>\s*`)
var webAwesomeFigcaptionRegex = regexp.MustCompile(`(?s)<figcaption[^>]*>(.*?)</figcaption>`)
var webAwesomeImageRegex = regexp.MustCompile(`(?s)<img\s+[^>]*>`)
var webAwesomeCodeTextRegex = regexp.MustCompile(`(?s)<code[^>]*>(.*?)</code>`)
var webAwesomeTagRegex = regexp.MustCompile(`<[^>]+>`)

func (p *WebAwesomePlugin) processPost(post *models.Post) error {
	if post.Skip || post.ArticleHTML == "" {
		return nil
	}

	needsWebAwesome := strings.Contains(post.ArticleHTML, "<wa-")
	post.ArticleHTML = webAwesomeComparisonRegex.ReplaceAllStringFunc(post.ArticleHTML, func(match string) string {
		needsWebAwesome = true
		parts := webAwesomeComparisonRegex.FindStringSubmatch(match)
		if len(parts) != 6 {
			return match
		}

		attrs := parseHTMLAttrs(parts[1] + " " + parts[3])
		beforeAttrs := parseHTMLAttrs(parts[4])
		afterAttrs := parseHTMLAttrs(parts[5])
		beforeSrc := beforeAttrs["src"]
		afterSrc := afterAttrs["src"]
		if beforeSrc == "" || afterSrc == "" {
			return match
		}

		position := normalizeWebAwesomePosition(attrs["position"])
		classAttr := strings.TrimSpace("markata-webawesome-comparison " + attrs["class"])

		var b strings.Builder
		caption := attrs["caption"]
		if caption != "" {
			b.WriteString(`<figure class="markata-webawesome-figure">`)
		}

		b.WriteString(`<wa-comparison class="`)
		b.WriteString(html.EscapeString(classAttr))
		b.WriteString(`"`)
		if id := attrs["id"]; id != "" {
			b.WriteString(` id="`)
			b.WriteString(html.EscapeString(id))
			b.WriteString(`"`)
		}
		if position != "" {
			b.WriteString(` position="`)
			b.WriteString(position)
			b.WriteString(`"`)
		}
		b.WriteString(`>`)
		// Map authored order left-to-right so the first image is initially on the left.
		b.WriteString(renderWebAwesomeComparisonImage("after", beforeAttrs))
		b.WriteString(renderWebAwesomeComparisonImage("before", afterAttrs))
		b.WriteString(`</wa-comparison>`)

		if caption != "" {
			b.WriteString(`<figcaption>`)
			b.WriteString(html.EscapeString(caption))
			b.WriteString(`</figcaption></figure>`)
		}

		return b.String()
	})
	post.ArticleHTML = p.processNestedTabs(post.ArticleHTML, &needsWebAwesome)
	post.ArticleHTML = p.processGenericContainers(post.ArticleHTML, &needsWebAwesome)
	if strings.Contains(post.ArticleHTML, "<wa-") {
		needsWebAwesome = true
	}
	if needsWebAwesome {
		if post.Extra == nil {
			post.Extra = make(map[string]interface{})
		}
		post.Extra["needs_webawesome"] = true
	}

	return nil
}

func (p *WebAwesomePlugin) processNestedTabs(content string, needsWebAwesome *bool) string {
	var b strings.Builder
	for {
		loc := webAwesomeNestedTabsOpenRegex.FindStringSubmatchIndex(content)
		if loc == nil {
			b.WriteString(content)
			break
		}

		b.WriteString(content[:loc[0]])
		openTag := content[loc[0]:loc[1]]
		closeIndex := matchingDivCloseIndex(content, loc[1])
		if closeIndex < 0 {
			b.WriteString(content[loc[0]:])
			break
		}

		body := content[loc[1]:closeIndex]
		attrs := parseHTMLAttrs(openTag)
		tabs, ok := extractNestedTabs(body)
		if !ok || len(tabs) == 0 {
			b.WriteString(content[loc[0] : closeIndex+len("</div>")])
			content = content[closeIndex+len("</div>"):]
			continue
		}

		var rendered strings.Builder
		rendered.WriteString(`<wa-tab-group`)
		writeAllowedAttrs(&rendered, attrs)
		rendered.WriteString(`>`)
		for i, tab := range tabs {
			label := tabLabel(tab.attrs, i)
			panel := slugWebAwesomeTabName(label, i)
			rendered.WriteString(`<wa-tab slot="nav" panel="`)
			rendered.WriteString(panel)
			rendered.WriteString(`">`)
			rendered.WriteString(html.EscapeString(label))
			rendered.WriteString(`</wa-tab>`)
		}
		for i, tab := range tabs {
			label := tabLabel(tab.attrs, i)
			panel := slugWebAwesomeTabName(label, i)
			rendered.WriteString(`<wa-tab-panel name="`)
			rendered.WriteString(panel)
			rendered.WriteString(`">`)
			rendered.WriteString(strings.TrimSpace(tab.body))
			rendered.WriteString(`</wa-tab-panel>`)
		}
		rendered.WriteString(`</wa-tab-group>`)
		*needsWebAwesome = true
		b.WriteString(rendered.String())
		content = content[closeIndex+len("</div>"):]
	}
	return b.String()
}

func matchingDivCloseIndex(content string, start int) int {
	depth := 1
	for index := start; index < len(content); {
		nextOpen := strings.Index(content[index:], "<div")
		nextClose := strings.Index(content[index:], "</div>")
		if nextClose < 0 {
			return -1
		}
		if nextOpen >= 0 && nextOpen < nextClose {
			depth++
			index += nextOpen + len("<div")
			continue
		}
		depth--
		if depth == 0 {
			return index + nextClose
		}
		index += nextClose + len("</div>")
	}
	return -1
}

type webAwesomeTab struct {
	attrs map[string]string
	body  string
}

func extractNestedTabs(body string) ([]webAwesomeTab, bool) {
	remaining := body
	result := make([]webAwesomeTab, 0, 4)

	for strings.TrimSpace(remaining) != "" {
		loc := webAwesomeNestedTabOpenRegex.FindStringSubmatchIndex(remaining)
		if len(loc) == 0 || strings.TrimSpace(remaining[:loc[0]]) != "" {
			return nil, false
		}

		openTag := remaining[loc[0]:loc[1]]
		closeIndex := matchingDivCloseIndex(remaining, loc[1])
		if closeIndex < 0 {
			return nil, false
		}

		result = append(result, webAwesomeTab{
			attrs: parseHTMLAttrs(openTag),
			body:  remaining[loc[1]:closeIndex],
		})
		remaining = remaining[closeIndex+len("</div>"):]
	}

	return result, len(result) > 0
}

func tabLabel(attrs map[string]string, index int) string {
	label := attrs["label"]
	if label == "" {
		label = attrs["title"]
	}
	if label == "" {
		label = fmt.Sprintf("Tab %d", index+1)
	}
	return label
}

func (p *WebAwesomePlugin) processGenericContainers(content string, needsWebAwesome *bool) string {
	var b strings.Builder
	for {
		loc := webAwesomeContainerOpenRegex.FindStringSubmatchIndex(content)
		if loc == nil {
			b.WriteString(content)
			break
		}

		b.WriteString(content[:loc[0]])
		openTag := content[loc[0]:loc[1]]
		closeIndex := matchingDivCloseIndex(content, loc[1])
		if closeIndex < 0 {
			b.WriteString(content[loc[0]:])
			break
		}

		fullMatch := content[loc[0] : closeIndex+len("</div>")]
		attrs := parseHTMLAttrs(openTag)
		component := webAwesomeComponentName(strings.Fields(attrs[htmlAttrClass]))
		if component == "" || component == "comparison" || component == "tabs" {
			b.WriteString(fullMatch)
			content = content[closeIndex+len("</div>"):]
			continue
		}

		body := strings.TrimSpace(content[loc[1]:closeIndex])
		rendered := renderWebAwesomeContainer(component, attrs, body)
		if rendered == "" {
			b.WriteString(fullMatch)
			content = content[closeIndex+len("</div>"):]
			continue
		}

		*needsWebAwesome = true
		b.WriteString(rendered)
		content = content[closeIndex+len("</div>"):]
	}

	return b.String()
}

func webAwesomeComponentName(classes []string) string {
	for i, class := range classes {
		if class == "webawesome" && i+1 < len(classes) {
			return classes[i+1]
		}
		if strings.HasPrefix(class, "wa-") {
			return strings.TrimPrefix(class, "wa-")
		}
	}
	return ""
}

func renderWebAwesomeContainer(component string, attrs map[string]string, body string) string {
	switch component {
	case "details":
		return renderWebAwesomeDetails(attrs, body)
	case "tabs":
		return renderWebAwesomeTabs(attrs, body)
	case "copy", "copy-button":
		return renderWebAwesomeCopyButton(attrs, body)
	case "qr", "qr-code":
		return renderWebAwesomeQRCode(attrs, body)
	case "badge":
		return renderWebAwesomeSimpleComponent("badge", attrs, body)
	case webAwesomeComponentTag:
		return renderWebAwesomeSimpleComponent(webAwesomeComponentTag, attrs, body)
	case "tooltip":
		return renderWebAwesomeTooltip(attrs, body)
	case "carousel":
		return renderWebAwesomeCarousel(attrs, body)
	case "animated-image":
		return renderWebAwesomeAnimatedImage(attrs, body)
	default:
		return ""
	}
}

func renderWebAwesomeDetails(attrs map[string]string, body string) string {
	summary := attrs["summary"]
	if summary == "" {
		summary = attrs["label"]
	}
	if summary == "" {
		summary = "Details"
	}

	var b strings.Builder
	b.WriteString(`<wa-details summary="`)
	b.WriteString(html.EscapeString(summary))
	b.WriteString(`"`)
	writeAllowedAttrs(&b, attrs, "summary", "label")
	b.WriteString(`>`)
	b.WriteString(body)
	b.WriteString(`</wa-details>`)
	return b.String()
}

func renderWebAwesomeTabs(attrs map[string]string, body string) string {
	matches := webAwesomeTabMarkerRegex.FindAllStringSubmatchIndex(body, -1)
	if len(matches) == 0 {
		return ""
	}

	var b strings.Builder
	b.WriteString(`<wa-tab-group`)
	writeAllowedAttrs(&b, attrs)
	b.WriteString(`>`)
	for i, match := range matches {
		name := tabMatchName(body, match)
		panel := slugWebAwesomeTabName(name, i)
		b.WriteString(`<wa-tab slot="nav" panel="`)
		b.WriteString(panel)
		b.WriteString(`">`)
		b.WriteString(html.EscapeString(name))
		b.WriteString(`</wa-tab>`)
	}
	for i, match := range matches {
		start := match[1]
		end := len(body)
		if i+1 < len(matches) {
			end = matches[i+1][0]
		}
		name := tabMatchName(body, match)
		panel := slugWebAwesomeTabName(name, i)
		b.WriteString(`<wa-tab-panel name="`)
		b.WriteString(panel)
		b.WriteString(`">`)
		b.WriteString(strings.TrimSpace(body[start:end]))
		b.WriteString(`</wa-tab-panel>`)
	}
	b.WriteString(`</wa-tab-group>`)
	return b.String()
}

func tabMatchName(body string, match []int) string {
	for i := 2; i+1 < len(match); i += 2 {
		if match[i] >= 0 && match[i+1] >= 0 {
			return html.UnescapeString(body[match[i]:match[i+1]])
		}
	}
	return "Tab"
}

func slugWebAwesomeTabName(name string, index int) string {
	var b strings.Builder
	lastDash := false
	for _, r := range strings.ToLower(name) {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
			lastDash = false
		} else if b.Len() > 0 && !lastDash {
			b.WriteByte('-')
			lastDash = true
		}
	}
	slug := strings.Trim(b.String(), "-")
	if slug == "" {
		slug = fmt.Sprintf("tab-%d", index+1)
	}
	return slug
}

func renderWebAwesomeCopyButton(attrs map[string]string, body string) string {
	value := attrs["value"]
	if value == "" {
		value = plainTextFromHTML(body)
	}
	if value == "" {
		return ""
	}
	var b strings.Builder
	b.WriteString(`<wa-copy-button value="`)
	b.WriteString(html.EscapeString(value))
	b.WriteString(`"`)
	writeAllowedAttrs(&b, attrs, "value")
	b.WriteString(`></wa-copy-button>`)
	return b.String()
}

func renderWebAwesomeQRCode(attrs map[string]string, body string) string {
	value := attrs["value"]
	if value == "" {
		value = plainTextFromHTML(body)
	}
	if value == "" {
		return ""
	}
	var b strings.Builder
	b.WriteString(`<wa-qr-code value="`)
	b.WriteString(html.EscapeString(value))
	b.WriteString(`"`)
	writeAllowedAttrs(&b, attrs, "value")
	b.WriteString(`></wa-qr-code>`)
	return b.String()
}

func renderWebAwesomeSimpleComponent(name string, attrs map[string]string, body string) string {
	body = unwrapSingleParagraph(body)
	if body == "" && attrs["label"] != "" {
		body = html.EscapeString(attrs["label"])
	}
	if body == "" {
		return ""
	}
	var b strings.Builder
	b.WriteString(`<wa-`)
	b.WriteString(name)
	writeAllowedAttrs(&b, attrs, "label")
	b.WriteString(`>`)
	b.WriteString(body)
	b.WriteString(`</wa-`)
	b.WriteString(name)
	b.WriteString(`>`)
	return b.String()
}

func renderWebAwesomeTooltip(attrs map[string]string, body string) string {
	content := attrs["content"]
	if content == "" {
		content = attrs["text"]
	}
	if content == "" {
		return ""
	}
	body = unwrapSingleParagraph(body)
	if body == "" && attrs["label"] != "" {
		body = html.EscapeString(attrs["label"])
	}
	if body == "" {
		return ""
	}
	// WA tooltip's default slot is the popup body; the trigger must be a
	// separate element referenced by `for=`. Render an inline anchor span
	// followed by a sibling <wa-tooltip for=...> whose default slot holds
	// the tooltip text.
	h := fnv.New64a()
	h.Write([]byte(body + "|" + content))
	sum := h.Sum(nil)
	id := "wa-tt-" + hex.EncodeToString(sum[:8])
	if seed := attrs["id"]; seed != "" {
		id = seed
	} else if seed := attrs["data-tooltip-id"]; seed != "" {
		id = seed
	} else {
		id = uniqueWebAwesomeTooltipID(id, body, content)
	}
	var b strings.Builder
	b.WriteString(`<span class="markata-wa-tooltip-anchor" id="`)
	b.WriteString(id)
	b.WriteString(`" tabindex="0">`)
	b.WriteString(body)
	b.WriteString(`</span><wa-tooltip for="`)
	b.WriteString(id)
	b.WriteString(`"`)
	writeAllowedAttrs(&b, attrs, "content", "text", "label")
	b.WriteString(`>`)
	b.WriteString(html.EscapeString(content))
	b.WriteString(`</wa-tooltip>`)
	return b.String()
}

func uniqueWebAwesomeTooltipID(base, body, content string) string {
	h := fnv.New64a()
	_, _ = h.Write([]byte(base))
	_, _ = h.Write([]byte("|"))
	_, _ = h.Write([]byte(body))
	_, _ = h.Write([]byte("|"))
	_, _ = h.Write([]byte(content))
	_, _ = h.Write([]byte("|"))
	_, _ = h.Write([]byte(strconv.FormatInt(nextWebAwesomeTooltipOrdinal(), 10)))
	sum := h.Sum(nil)
	return "wa-tt-" + hex.EncodeToString(sum[:8])
}

var webAwesomeTooltipOrdinal int64

func nextWebAwesomeTooltipOrdinal() int64 {
	return atomic.AddInt64(&webAwesomeTooltipOrdinal, 1)
}

func renderWebAwesomeCarousel(attrs map[string]string, body string) string {
	imageMatches := webAwesomeImageRegex.FindAllStringIndex(body, -1)
	if len(imageMatches) == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString(`<wa-carousel`)
	writeAllowedAttrs(&b, attrs)
	b.WriteString(`>`)
	for i, imageMatch := range imageMatches {
		image := body[imageMatch[0]:imageMatch[1]]
		captionEnd := len(body)
		if i+1 < len(imageMatches) {
			captionEnd = imageMatches[i+1][0]
		}
		caption := webAwesomeCarouselCaption(body[imageMatch[1]:captionEnd])

		b.WriteString(`<wa-carousel-item>`)
		if caption != "" {
			b.WriteString(`<figure class="markata-wa-carousel-figure">`)
			b.WriteString(renderWebAwesomeCarouselImage(image))
			b.WriteString(`<figcaption>`)
			b.WriteString(html.EscapeString(caption))
			b.WriteString(`</figcaption></figure>`)
		} else {
			b.WriteString(renderWebAwesomeCarouselImage(image))
		}
		b.WriteString(`</wa-carousel-item>`)
	}
	b.WriteString(`</wa-carousel>`)
	return b.String()
}

func webAwesomeCarouselCaption(segment string) string {
	if match := webAwesomeFigcaptionRegex.FindStringSubmatch(segment); len(match) == 2 {
		return normalizeWebAwesomeCaption(match[1])
	}
	return normalizeWebAwesomeCaption(segment)
}

func normalizeWebAwesomeCaption(caption string) string {
	caption = webAwesomeTagRegex.ReplaceAllString(caption, " ")
	caption = html.UnescapeString(caption)
	return strings.Join(strings.Fields(caption), " ")
}

func renderWebAwesomeCarouselImage(image string) string {
	attrs := parseHTMLAttrs(image)
	if attrs["src"] == "" {
		return image
	}

	var b strings.Builder
	b.WriteString(`<img`)
	writeImageAttr(&b, "src", attrs["src"])
	writeImageAttr(&b, "alt", attrs["alt"])
	writeImageAttr(&b, "title", attrs["title"])
	writeImageAttr(&b, "loading", defaultString(attrs["loading"], "lazy"))
	writeImageAttr(&b, "decoding", defaultString(attrs["decoding"], "async"))
	writeImageAttr(&b, "style", mergeCarouselImageStyle(attrs["style"]))
	b.WriteString(`>`)
	return b.String()
}

func writeImageAttr(b *strings.Builder, name, value string) {
	if value == "" {
		return
	}
	b.WriteByte(' ')
	b.WriteString(name)
	b.WriteString(`="`)
	b.WriteString(html.EscapeString(value))
	b.WriteByte('"')
}

func defaultString(value, fallback string) string {
	if value != "" {
		return value
	}
	return fallback
}

func mergeCarouselImageStyle(style string) string {
	fitStyle := "display: block; max-width: 100%; max-height: 100%; width: auto !important; height: auto !important; object-fit: contain !important; border-radius: var(--radius-lg);"
	if strings.TrimSpace(style) == "" {
		return fitStyle
	}
	return strings.TrimSpace(style) + " " + fitStyle
}

func renderWebAwesomeAnimatedImage(attrs map[string]string, body string) string {
	image := webAwesomeImageRegex.FindString(body)
	if image == "" {
		return ""
	}
	imageAttrs := parseHTMLAttrs(image)
	if imageAttrs["src"] == "" {
		return ""
	}
	var b strings.Builder
	b.WriteString(`<wa-animated-image src="`)
	b.WriteString(html.EscapeString(imageAttrs["src"]))
	b.WriteString(`" alt="`)
	b.WriteString(html.EscapeString(imageAttrs["alt"]))
	b.WriteString(`"`)
	writeAllowedAttrs(&b, attrs)
	b.WriteString(`></wa-animated-image>`)
	return b.String()
}

func writeAllowedAttrs(b *strings.Builder, attrs map[string]string, skip ...string) {
	skipMap := make(map[string]bool, len(skip))
	for _, key := range skip {
		skipMap[key] = true
	}
	keys := make([]string, 0, len(attrs))
	for key := range attrs {
		if skipMap[key] || key == htmlAttrClass {
			continue
		}
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		b.WriteByte(' ')
		b.WriteString(html.EscapeString(key))
		b.WriteString(`="`)
		b.WriteString(html.EscapeString(attrs[key]))
		b.WriteString(`"`)
	}
}

func unwrapSingleParagraph(body string) string {
	matches := webAwesomeFirstParagraphRegex.FindStringSubmatch(body)
	if len(matches) == 2 && strings.TrimSpace(matches[0]) == strings.TrimSpace(body) {
		return strings.TrimSpace(matches[1])
	}
	return strings.TrimSpace(body)
}

func plainTextFromHTML(body string) string {
	if matches := webAwesomeCodeTextRegex.FindStringSubmatch(body); len(matches) == 2 {
		return strings.TrimSpace(html.UnescapeString(matches[1]))
	}
	text := webAwesomeTagRegex.ReplaceAllString(body, "")
	return strings.TrimSpace(html.UnescapeString(text))
}

func parseHTMLAttrs(tag string) map[string]string {
	attrs := make(map[string]string)
	for _, match := range webAwesomeAttrRegex.FindAllStringSubmatch(tag, -1) {
		if len(match) != 3 {
			continue
		}
		value := strings.Trim(match[2], `"'`)
		attrs[strings.ToLower(match[1])] = html.UnescapeString(value)
	}
	return attrs
}

func normalizeWebAwesomePosition(position string) string {
	if position == "" {
		return ""
	}
	parsed, err := strconv.Atoi(position)
	if err != nil || parsed < 0 || parsed > 100 {
		return ""
	}
	return strconv.Itoa(parsed)
}

func renderWebAwesomeComparisonImage(slot string, attrs map[string]string) string {
	var b strings.Builder
	b.WriteString(`<img slot="`)
	b.WriteString(slot)
	b.WriteString(`" src="`)
	b.WriteString(html.EscapeString(attrs["src"]))
	b.WriteString(`" alt="`)
	b.WriteString(html.EscapeString(attrs["alt"]))
	b.WriteString(`" loading="lazy"`)
	if title := attrs["title"]; title != "" {
		b.WriteString(` title="`)
		b.WriteString(html.EscapeString(title))
		b.WriteString(`"`)
	}
	b.WriteString(`>`)
	return b.String()
}

func (p *WebAwesomePlugin) enableAssets(config *lifecycle.Config) {
	if config.Extra == nil {
		config.Extra = make(map[string]interface{})
	}

	cssURL := webAwesomeThemeCSSURL(p.config.CDNBase, p.config.Theme)
	loaderURL := p.config.CDNBase + "/webawesome.loader.js"
	basePath := p.config.CDNBase
	if p.config.Source == webAwesomeSourceVendor {
		if assetURL := resolveAssetURL(assetURLsFromConfig(config), webAwesomeAssetName, ""); assetURL != "" {
			basePath = assetURL
			cssURL = webAwesomeThemeCSSURL(basePath, p.config.Theme)
			loaderURL = basePath + "/webawesome.loader.js"
		} else {
			basePath = p.config.CDNBase
		}
	}

	config.Extra["webawesome_enabled"] = true
	config.Extra["webawesome_css_url"] = cssURL
	config.Extra["webawesome_loader_url"] = loaderURL
	config.Extra["webawesome_base_path"] = basePath
	config.Extra["webawesome_theme_class"] = fmt.Sprintf("wa-theme-%s wa-palette-%s wa-brand-%s", p.config.Theme, p.config.Palette, p.config.Brand)
}

func webAwesomeThemeCSSURL(basePath, theme string) string {
	theme = strings.Trim(theme, "/")
	if theme == "" {
		theme = ThemeDefault
	}
	return basePath + "/styles/themes/" + theme + ".css"
}

func (p *WebAwesomePlugin) enableVendorAsset(config *lifecycle.Config) {
	if !p.config.Enabled || p.config.Source != webAwesomeSourceVendor {
		return
	}
	if config.Extra == nil {
		config.Extra = make(map[string]interface{})
	}

	appendExtraAsset(config, assets.NewWebAwesomeAsset(p.config.Version, p.config.OutputDir))
}

func (p *WebAwesomePlugin) componentModules(htmlContent string) []string {
	seen := make(map[string]bool)
	modules := []string{}
	for _, match := range webAwesomeElementRegex.FindAllStringSubmatch(htmlContent, -1) {
		if len(match) != 2 {
			continue
		}
		name := strings.TrimPrefix(match[1], "wa-")
		if name == "" || seen[name] {
			continue
		}
		seen[name] = true
		modules = append(modules, p.componentModuleURL(name))
	}
	return modules
}

// componentModuleURL is retained for compatibility and tests; the runtime now
// uses webawesome.loader.js (autoloader) which discovers components via DOM
// observation, so per-component module URLs are not embedded in the page.
func (p *WebAwesomePlugin) componentModuleURL(name string) string {
	basePath := p.config.CDNBase
	if p.config.Source == webAwesomeSourceVendor {
		basePath = "/" + strings.Trim(p.config.OutputDir, "/")
	}
	return fmt.Sprintf("%s/components/%s/%s.js", basePath, name, name)
}

var (
	_ lifecycle.Plugin          = (*WebAwesomePlugin)(nil)
	_ lifecycle.ConfigurePlugin = (*WebAwesomePlugin)(nil)
	_ lifecycle.RenderPlugin    = (*WebAwesomePlugin)(nil)
	_ lifecycle.PriorityPlugin  = (*WebAwesomePlugin)(nil)
)
