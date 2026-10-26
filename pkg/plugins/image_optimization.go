package plugins

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/logging"
	"github.com/WaylonWalker/markata/pkg/models"
	markatatemplates "github.com/WaylonWalker/markata/pkg/templates"
	xdraw "golang.org/x/image/draw"
	"golang.org/x/image/webp"
)

var imageOptimizationLog = logging.Component("image_optimization").Phase("write")

// ImageOptimizationPlugin generates modern image formats for local images.
// It rewrites HTML to use <picture> with AVIF/WebP sources and caches encodes.
type ImageOptimizationPlugin struct {
	config           ImageOptimizationConfig
	availableFormats []string
	avifencPath      string
	cwebpPath        string
	warnedEncoders   map[string]bool
	stats            imageOptimizationStats
}

type imageOptimizationStats struct {
	targets, missingSources, variantsChecked, cacheHits, cacheMisses, encoderCalls, encoderFailures int
}

// ImageOptimizationConfig holds configuration for image optimization.
type ImageOptimizationConfig struct {
	Enabled     bool
	Formats     []string
	Quality     int
	AvifQuality int
	WebpQuality int
	Widths      []int
	Sizes       string
	CacheDir    string
	AvifencPath string
	CwebpPath   string
}

type imageOptimizationTarget struct {
	Src      string
	PostSlug string
	Widths   []int
}

type imageOptimizationVariant struct {
	Width int
	Path  string
}

type imageOptimizationCacheEntry struct {
	SourcePath    string `json:"source_path"`
	SourceSize    int64  `json:"source_size"`
	SourceModTime int64  `json:"source_mod_time"`
	Format        string `json:"format"`
	Width         int    `json:"width"`
	Quality       int    `json:"quality"`
	Encoder       string `json:"encoder"`
}

var responsiveRemoteImageWidths = []int{32, 48, 64, 88, 96, 120, 160, 240, 320, 480, 600, 800, 960, 1200, 1600}

const responsiveRemoteImageSizes = "(max-width: 760px) 100vw, 760px"

const (
	formatAVIF = "avif"
	formatWebP = "webp"

	extAVIF = ".avif"
	extJPG  = ".jpg"
	extJPEG = ".jpeg"
	extPNG  = ".png"
	extWebP = ".webp"
)

func NewImageOptimizationPlugin() *ImageOptimizationPlugin {
	return &ImageOptimizationPlugin{
		config:         defaultImageOptimizationConfig(),
		warnedEncoders: make(map[string]bool),
	}
}

func (p *ImageOptimizationPlugin) Name() string {
	return "image_optimization"
}

func (p *ImageOptimizationPlugin) Priority(stage lifecycle.Stage) int {
	switch stage {
	case lifecycle.StageRender:
		return 75
	case lifecycle.StageWrite:
		return lifecycle.PriorityEarly + 1
	default:
		return lifecycle.PriorityDefault
	}
}

func (p *ImageOptimizationPlugin) Configure(m *lifecycle.Manager) error {
	p.config = parseImageOptimizationConfig(m.Config())
	p.detectAvailableFormats()
	return nil
}

func (p *ImageOptimizationPlugin) Render(m *lifecycle.Manager) error {
	if !p.config.Enabled {
		return nil
	}

	posts := m.FilterPosts(func(post *models.Post) bool {
		if post.Skip || post.ArticleHTML == "" {
			return false
		}
		return strings.Contains(post.ArticleHTML, "<img")
	})

	return m.ProcessPostsSliceConcurrently(posts, func(post *models.Post) error {
		return p.processPost(post)
	})
}

func (p *ImageOptimizationPlugin) Write(m *lifecycle.Manager) error {
	p.stats = imageOptimizationStats{}
	if !p.config.Enabled {
		return nil
	}
	posts := m.FilterPosts(func(post *models.Post) bool {
		return !post.Skip && strings.Contains(post.HTML, "<img")
	})
	if err := m.ProcessPostsSliceConcurrently(posts, p.processRenderedHTML); err != nil {
		return err
	}
	if len(p.availableFormats) == 0 {
		// Even a build without encoders should report its cache activity.
		p.logSummary()
		return nil
	}

	targets := p.collectTargets(m)
	if len(targets) == 0 {
		p.logSummary()
		return nil
	}

	cacheDir := p.config.CacheDir
	if cacheDir == "" {
		cacheDir = ".markata/image-cache"
	}
	if err := os.MkdirAll(cacheDir, 0o755); err != nil {
		return fmt.Errorf("create image cache dir: %w", err)
	}

	outputDir := m.Config().OutputDir
	for _, target := range targets {
		if err := p.optimizeTarget(outputDir, cacheDir, target); err != nil {
			return err
		}
	}
	p.logSummary()

	return nil
}

func (p *ImageOptimizationPlugin) optimizeTarget(outputDir, cacheDir string, target imageOptimizationTarget) error {
	p.stats.targets++
	outputPath, err := resolveImageOutputPath(outputDir, target)
	if err != nil {
		p.stats.missingSources++
		imageOptimizationLog.Warnf("%v", err)
		return nil
	}

	info, err := os.Stat(outputPath)
	if err != nil {
		p.stats.missingSources++
		imageOptimizationLog.Warnf("source image not found: %s", outputPath)
		return nil
	}

	widths := target.Widths
	if len(widths) == 0 {
		widths = p.config.Widths
	}
	for _, format := range p.availableFormats {
		if err := p.optimizeImageFormat(outputPath, cacheDir, info, format, widths); err != nil {
			return err
		}
	}
	return nil
}

func (p *ImageOptimizationPlugin) logSummary() {
	imageOptimizationLog.Printf("targets=%d missing_sources=%d variants=%d cache_hits=%d cache_misses=%d encoder_calls=%d encoder_failures=%d",
		p.stats.targets, p.stats.missingSources, p.stats.variantsChecked, p.stats.cacheHits,
		p.stats.cacheMisses, p.stats.encoderCalls, p.stats.encoderFailures)
}

func (p *ImageOptimizationPlugin) optimizeImageFormat(outputPath, cacheDir string, info os.FileInfo, format string, widths []int) error {
	if shouldSkipFormat(outputPath, format) {
		return nil
	}
	quality := p.qualityForFormat(format)
	encoder := p.encoderPathForFormat(format)
	if encoder == "" {
		p.warnMissingEncoder(format)
		return nil
	}

	for _, variant := range buildImageVariants(outputPath, format, widths) {
		if variant.Path == "" {
			continue
		}
		if err := os.MkdirAll(filepath.Dir(variant.Path), 0o755); err != nil {
			return fmt.Errorf("create output dir: %w", err)
		}

		cachePath := imageOptimizationCachePath(cacheDir, outputPath, format, variant.Width, quality, encoder)
		p.stats.variantsChecked++
		if isImageCacheValid(cachePath, outputPath, info, format, variant.Width, quality, encoder, variant.Path) {
			p.stats.cacheHits++
			continue
		}
		p.stats.cacheMisses++
		if err := p.encodeImage(outputPath, variant.Path, format, quality, encoder, variant.Width); err != nil {
			p.stats.encoderFailures++
			imageOptimizationLog.Warnf("%v", err)
			continue
		}
		if err := writeImageCache(cachePath, outputPath, info, format, variant.Width, quality, encoder); err != nil {
			imageOptimizationLog.Warnf("cache write failed: %v", err)
		}
	}
	return nil
}

func defaultImageOptimizationConfig() ImageOptimizationConfig {
	return ImageOptimizationConfig{
		Enabled:     true,
		Formats:     []string{formatAVIF, formatWebP},
		Quality:     80,
		AvifQuality: 80,
		WebpQuality: 80,
		Widths:      []int{480, 768, 1200},
		Sizes:       "100vw",
		CacheDir:    ".markata/image-cache",
	}
}

func parseImageOptimizationConfig(cfg *lifecycle.Config) ImageOptimizationConfig {
	result := defaultImageOptimizationConfig()
	if cfg == nil || cfg.Extra == nil {
		return result
	}
	raw, ok := cfg.Extra["image_optimization"]
	if !ok {
		return result
	}

	if typed, ok := raw.(ImageOptimizationConfig); ok {
		return mergeImageOptimizationConfig(result, typed)
	}

	m, ok := raw.(map[string]any)
	if !ok {
		return result
	}

	if enabled, ok := m["enabled"].(bool); ok {
		result.Enabled = enabled
	}
	if formats, ok := m["formats"].([]any); ok {
		result.Formats = parseImageOptimizationFormats(formats)
	}
	if quality, ok := parseIntFromInterface(m["quality"]); ok {
		result.Quality = quality
	}
	if quality, ok := parseIntFromInterface(m["avif_quality"]); ok {
		result.AvifQuality = quality
	}
	if quality, ok := parseIntFromInterface(m["webp_quality"]); ok {
		result.WebpQuality = quality
	}
	if widths, ok := m["widths"].([]any); ok {
		result.Widths = parseIntSlice(widths)
	}
	if sizes, ok := m["sizes"].(string); ok && strings.TrimSpace(sizes) != "" {
		result.Sizes = strings.TrimSpace(sizes)
	}
	if cacheDir, ok := m["cache_dir"].(string); ok && cacheDir != "" {
		result.CacheDir = cacheDir
	}
	if path, ok := m["avifenc_path"].(string); ok {
		result.AvifencPath = path
	}
	if path, ok := m["cwebp_path"].(string); ok {
		result.CwebpPath = path
	}

	result.Widths = normalizeWidths(result.Widths)
	if strings.TrimSpace(result.Sizes) == "" {
		result.Sizes = "100vw"
	}

	return result
}

func mergeImageOptimizationConfig(base, override ImageOptimizationConfig) ImageOptimizationConfig {
	result := base
	result.Enabled = override.Enabled
	if len(override.Formats) > 0 {
		result.Formats = override.Formats
	}
	if override.Quality > 0 {
		result.Quality = override.Quality
	}
	if override.AvifQuality > 0 {
		result.AvifQuality = override.AvifQuality
	}
	if override.WebpQuality > 0 {
		result.WebpQuality = override.WebpQuality
	}
	if len(override.Widths) > 0 {
		result.Widths = override.Widths
	}
	if strings.TrimSpace(override.Sizes) != "" {
		result.Sizes = strings.TrimSpace(override.Sizes)
	}
	if override.CacheDir != "" {
		result.CacheDir = override.CacheDir
	}
	if override.AvifencPath != "" {
		result.AvifencPath = override.AvifencPath
	}
	if override.CwebpPath != "" {
		result.CwebpPath = override.CwebpPath
	}

	result.Widths = normalizeWidths(result.Widths)
	if strings.TrimSpace(result.Sizes) == "" {
		result.Sizes = "100vw"
	}

	return result
}

func parseImageOptimizationFormats(values []any) []string {
	result := make([]string, 0, len(values))
	for _, raw := range values {
		value, ok := raw.(string)
		if !ok {
			continue
		}
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		result = append(result, value)
	}
	return result
}

func parseIntSlice(values []any) []int {
	result := make([]int, 0, len(values))
	for _, raw := range values {
		if value, ok := parseIntFromInterface(raw); ok {
			result = append(result, value)
		}
	}
	return result
}

func (p *ImageOptimizationPlugin) detectAvailableFormats() {
	formats := make([]string, 0, len(p.config.Formats))
	seen := make(map[string]bool)

	for _, format := range p.config.Formats {
		normalized := strings.ToLower(strings.TrimSpace(format))
		if normalized == "" || seen[normalized] {
			continue
		}
		seen[normalized] = true

		switch normalized {
		case formatAVIF:
			path := p.config.AvifencPath
			if path == "" {
				if found, err := exec.LookPath("avifenc"); err == nil {
					path = found
				}
			}
			if path != "" {
				p.avifencPath = path
				formats = append(formats, normalized)
			} else {
				p.warnMissingEncoder(normalized)
			}
		case formatWebP:
			path := p.config.CwebpPath
			if path == "" {
				if found, err := exec.LookPath("cwebp"); err == nil {
					path = found
				}
			}
			if path != "" {
				p.cwebpPath = path
				formats = append(formats, normalized)
			} else {
				p.warnMissingEncoder(normalized)
			}
		}
	}

	p.availableFormats = formats
}

func (p *ImageOptimizationPlugin) warnMissingEncoder(format string) {
	if p.warnedEncoders[format] {
		return
	}
	p.warnedEncoders[format] = true
	switch format {
	case formatAVIF:
		imageOptimizationLog.Warnf("avifenc not found; skipping AVIF output")
	case formatWebP:
		imageOptimizationLog.Warnf("cwebp not found; skipping WebP output")
	}
}

func (p *ImageOptimizationPlugin) encoderPathForFormat(format string) string {
	switch format {
	case formatAVIF:
		return p.avifencPath
	case formatWebP:
		return p.cwebpPath
	default:
		return ""
	}
}

func (p *ImageOptimizationPlugin) qualityForFormat(format string) int {
	quality := p.config.Quality
	switch format {
	case formatAVIF:
		if p.config.AvifQuality > 0 {
			quality = p.config.AvifQuality
		}
	case formatWebP:
		if p.config.WebpQuality > 0 {
			quality = p.config.WebpQuality
		}
	}
	if quality <= 0 {
		quality = 80
	}
	if quality > 100 {
		quality = 100
	}
	return quality
}

func (p *ImageOptimizationPlugin) processPost(post *models.Post) error {
	context := &html.Node{Type: html.ElementNode, Data: "div", DataAtom: atom.Div}
	fragment, err := html.ParseFragment(strings.NewReader(post.ArticleHTML), context)
	if err != nil {
		return err
	}
	if len(fragment) == 0 {
		return nil
	}

	root := &html.Node{Type: html.ElementNode, Data: "div", DataAtom: atom.Div}
	for _, node := range fragment {
		root.AppendChild(node)
	}

	updatedHTML, targets := p.rewriteImages(root, post)
	if !updatedHTML {
		return nil
	}

	updated, err := renderFragment(root)
	if err != nil {
		return err
	}
	post.ArticleHTML = updated
	p.appendImageOptimizationTargets(post, targets)
	return nil
}

// processRenderedHTML applies image rewriting after the page template has been
// restored or rendered. This catches cached full-page HTML and images authored
// by templates, both of which are outside ArticleHTML during the Render stage.
func (p *ImageOptimizationPlugin) processRenderedHTML(post *models.Post) error {
	if post.HTML == "" {
		return nil
	}
	root, err := html.Parse(strings.NewReader(post.HTML))
	if err != nil {
		return err
	}
	updatedHTML, targets := p.rewriteImages(root, post)
	if !updatedHTML {
		return nil
	}
	updated, err := renderFragment(root)
	if err != nil {
		return err
	}
	post.HTML = updated
	p.appendImageOptimizationTargets(post, targets)
	return nil
}

func (p *ImageOptimizationPlugin) rewriteImages(root *html.Node, post *models.Post) (bool, []imageOptimizationTarget) {
	targets := make([]imageOptimizationTarget, 0)
	updatedHTML := false
	var walk func(*html.Node)
	walk = func(node *html.Node) {
		for child := node.FirstChild; child != nil; {
			next := child.NextSibling
			if child.Type == html.ElementNode && child.Data == "img" {
				changed, target := p.processImageNode(node, child, post)
				updatedHTML = updatedHTML || changed
				if target != nil {
					targets = append(targets, *target)
				}
				child = next
				continue
			}
			walk(child)
			child = next
		}
	}
	walk(root)
	return updatedHTML, targets
}

func (p *ImageOptimizationPlugin) appendImageOptimizationTargets(post *models.Post, targets []imageOptimizationTarget) {
	if len(targets) == 0 {
		return
	}
	if post.Extra == nil {
		post.Extra = make(map[string]any)
	}
	combined := imageOptimizationTargetsFrom(post.Extra["image_optimization"])
	post.Extra["image_optimization"] = append(combined, targets...)
}

func imageOptimizationTargetsFrom(raw any) []imageOptimizationTarget {
	switch items := raw.(type) {
	case []imageOptimizationTarget:
		return append([]imageOptimizationTarget(nil), items...)
	case []any:
		targets := make([]imageOptimizationTarget, 0, len(items))
		for _, item := range items {
			if target, ok := item.(imageOptimizationTarget); ok {
				targets = append(targets, target)
			}
		}
		return targets
	default:
		return nil
	}
}

func (p *ImageOptimizationPlugin) processImageNode(parent, img *html.Node, post *models.Post) (bool, *imageOptimizationTarget) {
	if parent.Data == "picture" {
		return false, nil
	}

	src := strings.TrimSpace(getAttr(img, "src"))
	if src == "" {
		return false, nil
	}
	if !isLocalImageSrc(src) {
		return addResponsiveRemoteImageAttrs(img, src), nil
	}
	if !isOptimizableImageSrc(src) {
		return false, nil
	}

	widths, sizes := localImageResponsiveSettings(img, p.config.Widths, p.config.Sizes)
	sources := buildPictureSources(src, p.availableFormats, widths, sizes)
	if len(sources) == 0 {
		return false, nil
	}

	picture := &html.Node{Type: html.ElementNode, Data: "picture"}
	for _, source := range sources {
		sourceNode, parseErr := parseSingleNode(source)
		if parseErr == nil {
			picture.AppendChild(sourceNode)
		}
	}

	parent.InsertBefore(picture, img)
	parent.RemoveChild(img)
	picture.AppendChild(img)

	if post.Slug == "" && !post.Has("_slug_explicit") {
		post.GenerateSlug()
	}
	return true, &imageOptimizationTarget{Src: src, PostSlug: post.Slug, Widths: widths}
}

func localImageResponsiveSettings(img *html.Node, configuredWidths []int, configuredSizes string) (widths []int, sizes string) {
	widths = normalizeWidths(configuredWidths)
	sizes = strings.TrimSpace(getAttr(img, "sizes"))
	if sizes == "" {
		sizes = strings.TrimSpace(configuredSizes)
	}
	if width, err := strconv.Atoi(strings.TrimSpace(getAttr(img, "width"))); err == nil && width > 0 {
		filtered := make([]int, 0, len(widths))
		for _, candidate := range widths {
			if candidate <= width {
				filtered = append(filtered, candidate)
			}
		}
		if len(filtered) == 0 {
			filtered = append(filtered, width)
		}
		widths = filtered
		if getAttr(img, "sizes") == "" {
			sizes = strconv.Itoa(width) + "px"
		}
	}
	return widths, sizes
}

func addResponsiveRemoteImageAttrs(img *html.Node, src string) bool {
	if getAttr(img, "srcset") != "" {
		return false
	}
	sizes := getAttr(img, "sizes")
	if sizes == "" {
		var ok bool
		sizes, ok = responsiveRemoteImageSizesFor(img)
		if !ok {
			return false
		}
	}
	srcset := markatatemplates.ResponsiveSrcset(src, responsiveRemoteImageWidths)
	if srcset == "" {
		return false
	}
	setAttr(img, "srcset", srcset)
	setAttr(img, "sizes", sizes)
	return true
}

func responsiveRemoteImageSizesFor(img *html.Node) (string, bool) {
	width := strings.TrimSpace(getAttr(img, "width"))
	if width != "" {
		if widthPx, err := strconv.Atoi(width); err == nil && widthPx > 0 {
			return strconv.Itoa(widthPx) + "px", true
		}
	}
	if height := strings.TrimSpace(getAttr(img, "height")); height != "" {
		if heightPx, err := strconv.Atoi(height); err == nil && heightPx > 0 && heightPx <= 64 {
			return "", false
		}
	}
	return responsiveRemoteImageSizes, true
}

func setAttr(node *html.Node, key, value string) {
	for i := range node.Attr {
		if node.Attr[i].Key == key {
			node.Attr[i].Val = value
			return
		}
	}
	node.Attr = append(node.Attr, html.Attribute{Key: key, Val: value})
}

func (p *ImageOptimizationPlugin) collectTargets(m *lifecycle.Manager) []imageOptimizationTarget {
	targets := make([]imageOptimizationTarget, 0)
	for _, post := range m.Posts() {
		if post.Extra == nil {
			continue
		}
		raw, ok := post.Extra["image_optimization"]
		if !ok {
			continue
		}
		if items, ok := raw.([]imageOptimizationTarget); ok {
			targets = append(targets, items...)
			continue
		}
		if items, ok := raw.([]any); ok {
			for _, item := range items {
				if target, ok := item.(imageOptimizationTarget); ok {
					targets = append(targets, target)
				}
			}
		}
	}
	return targets
}

func isLocalImageSrc(src string) bool {
	if strings.HasPrefix(src, "//") || strings.HasPrefix(src, "data:") {
		return false
	}
	parsed, err := url.Parse(src)
	if err != nil {
		return false
	}
	if parsed.Scheme != "" || parsed.Host != "" {
		return false
	}
	return true
}

func isOptimizableImageSrc(src string) bool {
	parsed, err := url.Parse(src)
	if err != nil {
		return false
	}
	ext := strings.ToLower(filepath.Ext(parsed.Path))
	switch ext {
	case extJPG, extJPEG, extPNG, extWebP:
		return true
	default:
		return false
	}
}

func buildPictureSources(src string, formats []string, widths []int, sizes string) []string {
	widths = normalizeWidths(widths)
	sources := make([]string, 0, len(formats))
	for _, format := range formats {
		mime := formatToMime(format)
		if len(widths) == 0 {
			srcset := replaceImageExtension(src, format)
			if srcset == "" {
				continue
			}
			sources = append(sources, fmt.Sprintf(`<source type=%q srcset=%q>`, mime, srcset))
			continue
		}
		items := make([]string, 0, len(widths))
		for _, width := range widths {
			srcset := replaceImageExtensionWithWidth(src, width, format)
			if srcset == "" {
				continue
			}
			items = append(items, fmt.Sprintf("%s %dw", srcset, width))
		}
		if len(items) == 0 {
			continue
		}
		sizes = strings.TrimSpace(sizes)
		if sizes == "" {
			sources = append(sources, fmt.Sprintf(`<source type=%q srcset=%q>`, mime, strings.Join(items, ", ")))
			continue
		}
		sources = append(sources, fmt.Sprintf(`<source type=%q srcset=%q sizes=%q>`, mime, strings.Join(items, ", "), sizes))
	}
	return sources
}

func renderFragment(root *html.Node) (string, error) {
	var buf bytes.Buffer
	for node := root.FirstChild; node != nil; node = node.NextSibling {
		if err := html.Render(&buf, node); err != nil {
			return "", err
		}
	}
	return buf.String(), nil
}

func parseSingleNode(fragment string) (*html.Node, error) {
	context := &html.Node{Type: html.ElementNode, Data: "div", DataAtom: atom.Div}
	nodes, err := html.ParseFragment(strings.NewReader(fragment), context)
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, fmt.Errorf("empty fragment")
	}
	return nodes[0], nil
}

func getAttr(node *html.Node, key string) string {
	for _, attr := range node.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

func formatToMime(format string) string {
	switch format {
	case formatAVIF:
		return "image/avif"
	case formatWebP:
		return "image/webp"
	default:
		return "application/octet-stream"
	}
}

func normalizeWidths(widths []int) []int {
	if len(widths) == 0 {
		return nil
	}
	seen := make(map[int]bool)
	filtered := make([]int, 0, len(widths))
	for _, width := range widths {
		if width <= 0 || seen[width] {
			continue
		}
		seen[width] = true
		filtered = append(filtered, width)
	}
	if len(filtered) == 0 {
		return nil
	}
	sort.Ints(filtered)
	return filtered
}

func replaceImageExtension(src, format string) string {
	parsed, err := url.Parse(src)
	if err != nil {
		return ""
	}
	ext := filepath.Ext(parsed.Path)
	if ext == "" {
		return ""
	}
	parsed.Path = strings.TrimSuffix(parsed.Path, ext) + "." + format
	return parsed.String()
}

func replaceImageExtensionWithWidth(src string, width int, format string) string {
	parsed, err := url.Parse(src)
	if err != nil {
		return ""
	}
	ext := filepath.Ext(parsed.Path)
	if ext == "" {
		return ""
	}
	base := strings.TrimSuffix(parsed.Path, ext)
	parsed.Path = fmt.Sprintf("%s-%dw.%s", base, width, format)
	return parsed.String()
}

func replaceImageFileExtension(path, format string) string {
	ext := filepath.Ext(path)
	if ext == "" {
		return ""
	}
	return strings.TrimSuffix(path, ext) + "." + format
}

func replaceImageFileExtensionWithWidth(path string, width int, format string) string {
	ext := filepath.Ext(path)
	if ext == "" {
		return ""
	}
	base := strings.TrimSuffix(path, ext)
	return fmt.Sprintf("%s-%dw.%s", base, width, format)
}

func shouldSkipFormat(sourcePath, format string) bool {
	ext := strings.ToLower(filepath.Ext(sourcePath))
	switch format {
	case formatAVIF:
		return ext == extAVIF
	case formatWebP:
		return ext == extWebP
	default:
		return true
	}
}

func buildImageVariants(outputPath, format string, widths []int) []imageOptimizationVariant {
	widths = normalizeWidths(widths)
	if len(widths) == 0 {
		return []imageOptimizationVariant{{Width: 0, Path: replaceImageFileExtension(outputPath, format)}}
	}
	variants := make([]imageOptimizationVariant, 0, len(widths))
	for _, width := range widths {
		variants = append(variants, imageOptimizationVariant{
			Width: width,
			Path:  replaceImageFileExtensionWithWidth(outputPath, width, format),
		})
	}
	return variants
}

func resolveImageOutputPath(outputDir string, target imageOptimizationTarget) (string, error) {
	parsed, err := url.Parse(target.Src)
	if err != nil {
		return "", fmt.Errorf("invalid image src %q", target.Src)
	}
	srcPath := filepath.FromSlash(parsed.Path)
	if srcPath == "" {
		return "", fmt.Errorf("empty image path for %q", target.Src)
	}

	var fullPath string
	if strings.HasPrefix(parsed.Path, "/") {
		fullPath = filepath.Join(outputDir, strings.TrimPrefix(srcPath, string(filepath.Separator)))
	} else {
		fullPath = filepath.Join(outputDir, target.PostSlug, srcPath)
	}

	cleaned := filepath.Clean(fullPath)
	outputDir = filepath.Clean(outputDir)
	if !isWithinDir(outputDir, cleaned) {
		return "", fmt.Errorf("image path escapes output dir: %s", cleaned)
	}

	return cleaned, nil
}

func isWithinDir(base, target string) bool {
	baseWithSep := base + string(filepath.Separator)
	targetWithSep := target + string(filepath.Separator)
	return strings.HasPrefix(targetWithSep, baseWithSep)
}

func imageOptimizationCachePath(cacheDir, sourcePath, format string, width, quality int, encoder string) string {
	hasher := sha256.New()
	hasher.Write([]byte(sourcePath))
	hasher.Write([]byte("|"))
	hasher.Write([]byte(format))
	hasher.Write([]byte("|"))
	fmt.Fprintf(hasher, "%d", width)
	hasher.Write([]byte("|"))
	fmt.Fprintf(hasher, "%d", quality)
	hasher.Write([]byte("|"))
	hasher.Write([]byte(encoder))
	return filepath.Join(cacheDir, hex.EncodeToString(hasher.Sum(nil))+".json")
}

func isImageCacheValid(cachePath, sourcePath string, info os.FileInfo, format string, width, quality int, encoder, destPath string) bool {
	if _, err := os.Stat(destPath); err != nil {
		return false
	}

	data, err := os.ReadFile(cachePath)
	if err != nil {
		return false
	}
	var entry imageOptimizationCacheEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return false
	}
	if entry.SourcePath != sourcePath || entry.Format != format || entry.Width != width || entry.Quality != quality || entry.Encoder != encoder {
		return false
	}
	if entry.SourceSize != info.Size() || entry.SourceModTime != info.ModTime().UnixNano() {
		return false
	}
	return true
}

func writeImageCache(cachePath, sourcePath string, info os.FileInfo, format string, width, quality int, encoder string) error {
	entry := imageOptimizationCacheEntry{
		SourcePath:    sourcePath,
		SourceSize:    info.Size(),
		SourceModTime: info.ModTime().UnixNano(),
		Format:        format,
		Width:         width,
		Quality:       quality,
		Encoder:       encoder,
	}
	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	return os.WriteFile(cachePath, data, 0o600)
}

func (p *ImageOptimizationPlugin) encodeImage(sourcePath, destPath, format string, quality int, encoder string, width int) error {
	if encoder == "" {
		return fmt.Errorf("missing encoder for format %s", format)
	}
	switch format {
	case formatAVIF:
		inputPath, cleanup, err := prepareAVIFInput(sourcePath, width)
		if err != nil {
			return fmt.Errorf("prepare AVIF input %s: %w", sourcePath, err)
		}
		defer cleanup()
		args := avifEncoderArgs(inputPath, destPath, quality)
		// #nosec G204 -- encoder comes from config or LookPath and is validated.
		cmd := exec.Command(encoder, args...)
		p.stats.encoderCalls++
		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("avifenc failed for %s: %w (output: %s)", sourcePath, err, string(output))
		}
		return nil
	case formatWebP:
		args := []string{"-q", fmt.Sprintf("%d", quality)}
		if width > 0 {
			args = append(args, "-resize", fmt.Sprintf("%d", width), "0")
		}
		args = append(args, sourcePath, "-o", destPath)
		// #nosec G204 -- encoder comes from config or LookPath and is validated.
		cmd := exec.Command(encoder, args...)
		p.stats.encoderCalls++
		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("cwebp failed for %s: %w (output: %s)", sourcePath, err, string(output))
		}
		return nil
	default:
		return fmt.Errorf("unsupported image format: %s", format)
	}
}

func avifEncoderArgs(sourcePath, destPath string, quality int) []string {
	return []string{"-q", strconv.Itoa(quality), sourcePath, destPath}
}

func prepareAVIFInput(sourcePath string, width int) (inputPath string, cleanup func(), err error) {
	noCleanup := func() {}
	file, err := os.Open(sourcePath)
	if err != nil {
		return "", noCleanup, err
	}
	source, decodeErr := decodeImageForAVIF(file, sourcePath)
	closeErr := file.Close()
	if decodeErr != nil {
		return "", noCleanup, decodeErr
	}
	if closeErr != nil {
		return "", noCleanup, closeErr
	}

	bounds := source.Bounds()
	sourceWidth := bounds.Dx()
	sourceHeight := bounds.Dy()
	if sourceWidth <= 0 || sourceHeight <= 0 {
		return "", noCleanup, fmt.Errorf("invalid image bounds %v", bounds)
	}
	if width <= 0 || width >= sourceWidth {
		if strings.EqualFold(filepath.Ext(sourcePath), extWebP) {
			width = 0
		} else {
			return sourcePath, noCleanup, nil
		}
	}

	resizedWidth := sourceWidth
	resizedHeight := sourceHeight
	if width > 0 {
		resizedWidth = width
		resizedHeight = (sourceHeight*width + sourceWidth/2) / sourceWidth
		if resizedHeight < 1 {
			resizedHeight = 1
		}
	}
	resized := image.NewNRGBA(image.Rect(0, 0, resizedWidth, resizedHeight))
	xdraw.CatmullRom.Scale(resized, resized.Bounds(), source, bounds, xdraw.Src, nil)

	temp, err := os.CreateTemp(filepath.Dir(sourcePath), ".markata-avif-*.png")
	if err != nil {
		return "", noCleanup, err
	}
	tempPath := temp.Name()
	removeTemp := func() { _ = os.Remove(tempPath) }
	if err := png.Encode(temp, resized); err != nil {
		_ = temp.Close()
		removeTemp()
		return "", noCleanup, err
	}
	if err := temp.Close(); err != nil {
		removeTemp()
		return "", noCleanup, err
	}
	return tempPath, removeTemp, nil
}

func decodeImageForAVIF(reader io.Reader, sourcePath string) (image.Image, error) {
	switch strings.ToLower(filepath.Ext(sourcePath)) {
	case extJPG, extJPEG:
		return jpeg.Decode(reader)
	case extPNG:
		return png.Decode(reader)
	case extWebP:
		return webp.Decode(reader)
	default:
		decoded, _, err := image.Decode(reader)
		return decoded, err
	}
}
