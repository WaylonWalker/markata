package plugins

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/WaylonWalker/markata/pkg/buildcache"
	"github.com/WaylonWalker/markata/pkg/imageindex"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/templates"
)

// ImageLibraryConfig controls the generated image library.
type ImageLibraryConfig = models.ImagesConfig

// ImageLibraryPlugin writes image inventory JSON artifacts and an accessible
// browser-based image library page.
type ImageLibraryPlugin struct {
	engineMu    sync.RWMutex
	engineCache map[string]*templates.Engine
}

// NewImageLibraryPlugin creates a new image library plugin.
func NewImageLibraryPlugin() *ImageLibraryPlugin {
	return &ImageLibraryPlugin{engineCache: make(map[string]*templates.Engine)}
}

// NewImagesPlugin is a compatibility constructor for the short plugin name.
func NewImagesPlugin() *ImageLibraryPlugin {
	return NewImageLibraryPlugin()
}

// Name returns the unique plugin name.
func (p *ImageLibraryPlugin) Name() string {
	return "images"
}

// Priority runs the aggregate writer after post output and before the final
// content index writer.
func (p *ImageLibraryPlugin) Priority(stage lifecycle.Stage) int {
	if stage == lifecycle.StageWrite {
		return lifecycle.PriorityLate
	}
	return lifecycle.PriorityDefault
}

// Write builds the canonical image index and renders the image library page.
func (p *ImageLibraryPlugin) Write(m *lifecycle.Manager) error {
	if m == nil || m.Config() == nil {
		return fmt.Errorf("image library requires a lifecycle manager and config")
	}

	config, err := parseImageLibraryConfig(m.Config().Extra)
	if err != nil {
		return err
	}
	pathPrefix, err := imageLibraryPath(config.Path)
	if err != nil {
		return err
	}
	templateName, err := imageLibraryTemplate(config.Template)
	if err != nil {
		return err
	}
	config.Template = templateName

	contentDir := m.Config().ContentDir
	if contentDir == "" {
		contentDir = "."
	}
	contentRoot, err := filepath.Abs(contentDir)
	if err != nil {
		return fmt.Errorf("resolve image library content directory: %w", err)
	}
	assetsDir := imageLibraryAssetsDir(m.Config(), contentDir)
	outputDir := m.Config().OutputDir
	if outputDir == "" {
		outputDir = defaultOutputDir
	}
	outputRoot, err := filepath.Abs(outputDir)
	if err != nil {
		return fmt.Errorf("resolve image library output directory: %w", err)
	}
	pagePath := filepath.Join(outputDir, filepath.FromSlash(pathPrefix), "index.html")
	jsonPath := filepath.Join(outputDir, filepath.FromSlash(pathPrefix), "index.json")
	flatJSONPath := filepath.Join(outputDir, "images.json")
	if !config.IsEnabled() {
		return removeImageLibraryOutputs(contentRoot, GetBuildCache(m))
	}
	for _, path := range []string{pagePath, jsonPath, flatJSONPath} {
		if err := imageLibraryValidateOutputPath(outputRoot, path); err != nil {
			return err
		}
	}
	adoptPublishedImageLibraryOutputs(GetBuildCache(m), contentRoot, outputRoot)
	if err := imageLibraryOutputConflict(m, contentRoot, outputRoot, pagePath, jsonPath, flatJSONPath, assetsDir, pathPrefix, config.ShouldExportJSON()); err != nil {
		return err
	}
	if !config.ShouldExportJSON() {
		for _, path := range []string{jsonPath, flatJSONPath} {
			if imageLibraryOutputMatches(GetBuildCache(m), contentRoot, outputRoot, path) {
				if err := removeImageLibraryFile(outputRoot, path); err != nil {
					return fmt.Errorf("remove disabled image index export: %w", err)
				}
			}
		}
	}
	return p.writeImageLibrary(m, &config, contentDir, assetsDir, contentRoot, outputRoot, pathPrefix, pagePath, jsonPath, flatJSONPath)
}

// A staged clean build records ownership under its temporary output root.
// After publication, transfer that ownership only when every recorded file
// exists unchanged under the final output root.
func adoptPublishedImageLibraryOutputs(cache *buildcache.Cache, contentRoot, outputRoot string) {
	if cache == nil {
		return
	}
	previousContentRoot, previousRoot, outputs := cache.GetImageLibraryOutputs()
	if filepath.Clean(previousContentRoot) != filepath.Clean(contentRoot) ||
		filepath.Clean(previousRoot) == filepath.Clean(outputRoot) ||
		filepath.Dir(previousRoot) != filepath.Dir(outputRoot) ||
		!strings.HasPrefix(filepath.Base(previousRoot), "."+filepath.Base(outputRoot)+".markata-build-") ||
		len(outputs) == 0 {
		return
	}
	hashes := make(map[string]string, len(outputs))
	for relative, expected := range outputs {
		path, err := imageLibraryOutputPath(outputRoot, relative)
		if err != nil || !isRegularImageLibraryFile(path) {
			return
		}
		actual, err := buildcache.HashFile(path)
		if err != nil || actual != expected {
			return
		}
		hashes[path] = expected
	}
	cache.SetImageLibraryOutputs(contentRoot, outputRoot, hashes)
}

func imageLibraryOutputConflict(m *lifecycle.Manager, contentRoot, outputRoot, pagePath, jsonPath, flatJSONPath, assetsDir, pathPrefix string, exportJSON bool) error {
	if exportJSON && imageLibraryPathsOverlap(pagePath, flatJSONPath) {
		return fmt.Errorf("image library path %q conflicts with root images.json", pathPrefix)
	}
	outputDir := m.Config().OutputDir
	if outputDir == "" {
		outputDir = defaultOutputDir
	}
	if err := imageLibraryPostOutputConflict(m, outputDir, pagePath, flatJSONPath, exportJSON); err != nil {
		return err
	}
	if err := imageLibraryFeedOutputConflict(m, outputDir, pagePath, jsonPath, flatJSONPath, exportJSON); err != nil {
		return err
	}
	if err := imageLibraryStaticOutputConflict(assetsDir, pathPrefix, exportJSON); err != nil {
		return err
	}
	cache := GetBuildCache(m)
	for _, path := range []string{pagePath, jsonPath, flatJSONPath} {
		if err := imageLibraryExistingOutputConflict(cache, contentRoot, outputRoot, path, true); err != nil {
			return err
		}
	}
	return imageLibraryStaleOutputConflict(cache, contentRoot, pagePath, jsonPath, flatJSONPath)
}

func imageLibraryPostOutputConflict(m *lifecycle.Manager, outputDir, pagePath, flatJSONPath string, exportJSON bool) error {
	for _, post := range m.Posts() {
		if post == nil || post.Skip || post.Draft {
			continue
		}
		postOutput := filepath.Join(outputDir, post.Slug, "index.html")
		if imageLibraryPathsOverlap(postOutput, pagePath) || (exportJSON && imageLibraryPathsOverlap(postOutput, flatJSONPath)) {
			return fmt.Errorf("image library output conflicts with post %q at %s", post.Path, pagePath)
		}
	}
	return nil
}

func imageLibraryFeedOutputConflict(m *lifecycle.Manager, outputDir, pagePath, jsonPath, flatJSONPath string, exportJSON bool) error {
	if cached, ok := m.Cache().Get("feed_configs"); ok {
		if feedConfigs, ok := cached.([]models.FeedConfig); ok {
			checker := NewOverwriteCheckPlugin()
			for index := range feedConfigs {
				for _, feedOutput := range checker.getFeedOutputPaths(outputDir, &feedConfigs[index]) {
					if filepath.Clean(feedOutput) == filepath.Clean(pagePath) || (exportJSON && (filepath.Clean(feedOutput) == filepath.Clean(jsonPath) || filepath.Clean(feedOutput) == filepath.Clean(flatJSONPath))) {
						return fmt.Errorf("image library output conflicts with feed %q at %s", feedConfigs[index].Slug, feedOutput)
					}
				}
			}
		}
	}
	return nil
}

func imageLibraryStaticOutputConflict(assetsDir, pathPrefix string, exportJSON bool) error {
	staticPage := filepath.Join(assetsDir, filepath.FromSlash(pathPrefix), "index.html")
	if imageLibraryPathExists(staticPage) {
		return fmt.Errorf("image library output conflicts with static file %s", staticPage)
	}
	if exportJSON {
		staticJSON := filepath.Join(assetsDir, filepath.FromSlash(pathPrefix), "index.json")
		if imageLibraryPathExists(staticJSON) {
			return fmt.Errorf("image library output conflicts with static file %s", staticJSON)
		}
		flatJSON := filepath.Join(assetsDir, "images.json")
		if imageLibraryPathExists(flatJSON) {
			return fmt.Errorf("image library output conflicts with static file %s", flatJSON)
		}
	}
	return nil
}

func imageLibraryExistingOutputConflict(cache *buildcache.Cache, contentRoot, outputRoot, path string, required bool) error {
	if !imageLibraryPathExists(path) {
		return nil
	}
	if imageLibraryOutputMatches(cache, contentRoot, outputRoot, path) {
		return nil
	}
	if imageLibraryOwnsOutput(cache, contentRoot, outputRoot, path) {
		return fmt.Errorf("image library output was modified outside the images plugin: %s", path)
	}
	if required {
		return fmt.Errorf("image library output already exists and is not owned by the images plugin: %s", path)
	}
	return nil
}

func (p *ImageLibraryPlugin) writeImageLibrary(m *lifecycle.Manager, config *ImageLibraryConfig, contentDir, assetsDir, contentRoot, outputRoot, pathPrefix, pagePath, jsonPath, flatJSONPath string) error {
	cache := GetBuildCache(m)
	inputHash, err := p.imageLibraryInputHash(m, config, contentDir, assetsDir, outputRoot, pathPrefix)
	if err != nil {
		return fmt.Errorf("hashing image library inputs: %w", err)
	}
	previousContentRoot := ""
	previousRoot := ""
	previousOutputs := map[string]string(nil)
	if cache != nil {
		previousContentRoot, previousRoot, previousOutputs = cache.GetImageLibraryOutputs()
	}
	if cache != nil && cache.GetImageLibraryHash() == inputHash && imageLibraryOutputsExist(pagePath, jsonPath, flatJSONPath, config.ShouldExportJSON()) && imageLibraryTrackedOutputsMatch(cache, contentRoot, outputRoot, previousOutputs) {
		return nil
	}

	index, err := imageindex.Build(m.Posts(), imageindex.BuildOptions{
		ContentDir:          contentDir,
		AssetsDir:           assetsDir,
		OutputDir:           outputRoot,
		GeneratorVersion:    imageLibraryGeneratorVersion(m.Config()),
		IncludeUnreferenced: config.ShouldIncludeUnreferenced(),
	})
	if err != nil {
		return fmt.Errorf("building image index: %w", err)
	}
	pages := imageLibraryPaginate(newImageLibraryPage(index), config.PageSize, pathPrefix)
	pagePaths, err := imageLibraryPagePaths(m, cache, contentRoot, outputRoot, pagePath, jsonPath, flatJSONPath, len(pages), config.ShouldExportJSON())
	if err != nil {
		return err
	}
	ownedPaths := append([]string(nil), pagePaths...)
	if config.ShouldExportJSON() {
		ownedPaths = append(ownedPaths, jsonPath, flatJSONPath)
	}
	outputFS, err := openOutputRoot(outputRoot)
	if err != nil {
		return fmt.Errorf("open image library output directory: %w", err)
	}
	defer outputFS.Close()

	if config.ShouldExportJSON() {
		data, err := imageindex.Marshal(index)
		if err != nil {
			return fmt.Errorf("marshal image index: %w", err)
		}
		if err := writeImageLibraryFile(outputFS, outputRoot, jsonPath, data); err != nil {
			return fmt.Errorf("write image index: %w", err)
		}
		if err := writeImageLibraryFile(outputFS, outputRoot, flatJSONPath, data); err != nil {
			return fmt.Errorf("write root image index: %w", err)
		}
	}

	for i := range pages {
		if err := p.renderPage(m, config, pathPrefix, index, pages[i], outputFS, outputRoot, pagePaths[i]); err != nil {
			return err
		}
	}
	if cache != nil {
		if err := removeStaleImageLibraryOutputs(previousContentRoot, previousRoot, previousOutputs, contentRoot, ownedPaths); err != nil {
			return fmt.Errorf("remove stale image library outputs: %w", err)
		}
		outputHashes, err := imageLibraryOutputHashesForPaths(ownedPaths)
		if err != nil {
			return fmt.Errorf("hash image library outputs: %w", err)
		}
		cache.SetImageLibraryHash(inputHash)
		cache.SetImageLibraryOutputs(contentRoot, outputRoot, outputHashes)
	}
	log.Printf("[images] Generated /%s/ with %d images", pathPrefix, len(index.Images))
	return nil
}

func parseImageLibraryConfig(extra map[string]interface{}) (ImageLibraryConfig, error) {
	config := models.NewImagesConfig()
	if extra == nil {
		return config, nil
	}
	if modelsConfig, ok := extra["models_config"].(*models.Config); ok && modelsConfig != nil {
		return normalizeImageLibraryConfig(modelsConfig.Images), nil
	}
	raw, ok := extra["images"]
	if !ok {
		if raw, ok = extra["image_library"]; !ok {
			return config, nil
		}
	}
	if typed, ok := raw.(ImageLibraryConfig); ok {
		return normalizeImageLibraryConfig(typed), nil
	}
	if typed, ok := raw.(*ImageLibraryConfig); ok && typed != nil {
		return normalizeImageLibraryConfig(*typed), nil
	}
	values, ok := raw.(map[string]interface{})
	if !ok {
		return ImageLibraryConfig{}, fmt.Errorf("images must be a table")
	}
	return parseImageLibraryConfigFields(config, values)
}

// Keep individual field validation separate from the dispatch path so the
// parser stays straightforward as optional image-library features grow.
func parseImageLibraryConfigFields(config ImageLibraryConfig, values map[string]interface{}) (ImageLibraryConfig, error) {
	if value, exists := values["enabled"]; exists {
		parsed, ok := value.(bool)
		if !ok {
			return ImageLibraryConfig{}, fmt.Errorf("images.enabled must be boolean")
		}
		config.Enabled = &parsed
	}
	if value, exists := values["path"]; exists {
		parsed, ok := value.(string)
		if !ok {
			return ImageLibraryConfig{}, fmt.Errorf("images.path must be a string")
		}
		config.Path = parsed
	}
	if value, exists := values["template"]; exists {
		parsed, ok := value.(string)
		if !ok {
			return ImageLibraryConfig{}, fmt.Errorf("images.template must be a string")
		}
		config.Template = parsed
	}
	if value, exists := values["export_json"]; exists {
		parsed, ok := value.(bool)
		if !ok {
			return ImageLibraryConfig{}, fmt.Errorf("images.export_json must be boolean")
		}
		config.ExportJSON = &parsed
	}
	if value, exists := values["include_unreferenced"]; exists {
		parsed, ok := value.(bool)
		if !ok {
			return ImageLibraryConfig{}, fmt.Errorf("images.include_unreferenced must be boolean")
		}
		config.IncludeUnreferenced = &parsed
	}
	if value, exists := values["page_size"]; exists {
		size, err := parseImageLibraryPageSize(value)
		if err != nil {
			return ImageLibraryConfig{}, err
		}
		config.PageSize = size
	}
	return normalizeImageLibraryConfig(config), nil
}

func parseImageLibraryPageSize(value interface{}) (int, error) {
	var size int
	switch parsed := value.(type) {
	case int:
		size = parsed
	case int64:
		if parsed > int64(int(^uint(0)>>1)) {
			return 0, fmt.Errorf("images.page_size exceeds supported integer range")
		}
		size = int(parsed)
	case float64:
		if parsed != math.Trunc(parsed) || parsed > float64(int(^uint(0)>>1)) || parsed < 0 {
			return 0, fmt.Errorf("images.page_size must be a nonnegative integer")
		}
		size = int(parsed)
	default:
		return 0, fmt.Errorf("images.page_size must be a nonnegative integer")
	}
	if size < 0 {
		return 0, fmt.Errorf("images.page_size must be nonnegative")
	}
	return size, nil
}

func normalizeImageLibraryConfig(config ImageLibraryConfig) ImageLibraryConfig {
	defaults := models.NewImagesConfig()
	if config.Enabled == nil {
		config.Enabled = defaults.Enabled
	}
	if config.Path == "" {
		config.Path = defaults.Path
	}
	if config.Template == "" {
		config.Template = defaults.Template
	}
	if config.ExportJSON == nil {
		config.ExportJSON = defaults.ExportJSON
	}
	if config.IncludeUnreferenced == nil {
		config.IncludeUnreferenced = defaults.IncludeUnreferenced
	}
	return config
}

func imageLibraryPath(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		raw = "images"
	}
	if filepath.IsAbs(raw) || strings.HasPrefix(raw, "/") {
		return "", fmt.Errorf("images.path must stay within output_dir: %q", raw)
	}
	clean := filepath.ToSlash(filepath.Clean(raw))
	if clean == "." || clean == ".." || strings.HasPrefix(clean, "../") {
		return "", fmt.Errorf("images.path must stay within output_dir: %q", raw)
	}
	return strings.Trim(clean, "/"), nil
}

func imageLibraryTemplate(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		raw = "images.html"
	}
	if filepath.IsAbs(raw) || strings.HasPrefix(raw, "/") {
		return "", fmt.Errorf("images.template must stay within the template directory: %q", raw)
	}
	clean := filepath.ToSlash(filepath.Clean(raw))
	if clean == "." || clean == ".." || strings.HasPrefix(clean, "../") {
		return "", fmt.Errorf("images.template must stay within the template directory: %q", raw)
	}
	return clean, nil
}

func imageLibraryAssetsDir(config *lifecycle.Config, contentDir string) string {
	assetsDir := ""
	if config != nil && config.Extra != nil {
		if value, ok := config.Extra["assets_dir"].(string); ok {
			assetsDir = value
		}
	}
	if assetsDir == "" {
		if configured, ok := getModelsConfig(config); ok && configured != nil {
			assetsDir = configured.AssetsDir
		}
	}
	if assetsDir == "" {
		assetsDir = "static"
	}
	if !filepath.IsAbs(assetsDir) {
		assetsDir = filepath.Join(contentDir, assetsDir)
	}
	return filepath.Clean(assetsDir)
}

func imageLibraryGeneratorVersion(config *lifecycle.Config) string {
	if config != nil && config.Extra != nil {
		if version, ok := config.Extra["markata_version"].(string); ok && version != "" {
			return version
		}
	}
	return "dev"
}

// imageLibraryRenderConfigHash fingerprints the configuration exposed to the
// image-library template. The generated page uses the full site config, not
// only the images section, so changes to site metadata, navigation, theme, or
// plugin-provided template values must invalidate the cached page.
func imageLibraryRenderConfigHash(config *lifecycle.Config) (string, error) {
	modelsConfig, ok := getModelsConfig(config)
	if !ok || modelsConfig == nil {
		modelsConfig = ToModelsConfig(config)
	}
	if modelsConfig == nil {
		return "", nil
	}

	serialized, err := json.Marshal(modelsConfig)
	if err != nil {
		return "", fmt.Errorf("marshal template configuration: %w", err)
	}
	extra, err := json.Marshal(modelsConfig.Extra)
	if err != nil {
		return "", fmt.Errorf("marshal template configuration extras: %w", err)
	}
	return buildcache.ContentHash(string(serialized) + "\x00" + string(extra)), nil
}

func (p *ImageLibraryPlugin) imageLibraryInputHash(m *lifecycle.Manager, config *ImageLibraryConfig, contentDir, assetsDir, outputDir, pathPrefix string) (string, error) {
	var input strings.Builder
	input.WriteString(pathPrefix)
	input.WriteByte('\x00')
	input.WriteString(imageLibraryGeneratorVersion(m.Config()))
	input.WriteByte('\x00')
	input.WriteString(config.Template)
	input.WriteByte('\x00')
	templateHash, err := p.imageLibraryTemplateHash(m, config.Template)
	if err != nil {
		return "", err
	}
	input.WriteString(templateHash)
	input.WriteByte('\x00')
	input.WriteString(fmt.Sprintf("%t", config.IsEnabled()))
	input.WriteByte('\x00')
	input.WriteString(fmt.Sprintf("%t", config.ShouldExportJSON()))
	input.WriteByte('\x00')
	input.WriteString(fmt.Sprintf("%t", config.ShouldIncludeUnreferenced()))
	input.WriteByte('\x00')
	input.WriteString(strconv.Itoa(config.PageSize))
	input.WriteByte('\x00')
	renderConfigHash, err := imageLibraryRenderConfigHash(m.Config())
	if err != nil {
		return "", err
	}
	input.WriteString(renderConfigHash)
	input.WriteByte('\x00')

	publicPosts := make([]*models.Post, 0, len(m.Posts()))
	privateFrontmatterInputs := make(map[string]struct{})
	for _, post := range m.Posts() {
		if post == nil {
			continue
		}
		if post.Published && !post.Draft && !post.Private && !post.Skip {
			publicPosts = append(publicPosts, post)
			continue
		}
		if privateInput := imageindex.PrivateFrontmatterHashInput(post); privateInput != "" {
			privateFrontmatterInputs[privateInput] = struct{}{}
		}
	}
	sort.SliceStable(publicPosts, func(i, j int) bool {
		left, right := publicPosts[i], publicPosts[j]
		leftContent := left.Content + "\x00" + left.ArticleHTML + "\x00" + left.RawFrontmatter
		rightContent := right.Content + "\x00" + right.ArticleHTML + "\x00" + right.RawFrontmatter
		if left.Path != right.Path {
			return left.Path < right.Path
		}
		if left.Slug != right.Slug {
			return left.Slug < right.Slug
		}
		if left.Href != right.Href {
			return left.Href < right.Href
		}
		return leftContent < rightContent
	})
	for _, post := range publicPosts {
		input.WriteString(post.Path)
		input.WriteByte('\x00')
		input.WriteString(post.Slug)
		input.WriteByte('\x00')
		input.WriteString(post.Href)
		input.WriteByte('\x00')
		input.WriteString(post.PlainTitle())
		input.WriteByte('\x00')
		if post.Date != nil {
			input.WriteString(post.Date.UTC().Format(time.RFC3339Nano))
		}
		input.WriteByte('\x00')
		input.WriteString(buildcache.ContentHash(strings.Join([]string{
			post.InputHash,
			post.Content,
			post.ArticleHTML,
			post.RawFrontmatter,
			frontmatterImageHash(post),
		}, "\x00")))
		input.WriteByte('\x00')
	}
	privateInputs := make([]string, 0, len(privateFrontmatterInputs))
	for privateInput := range privateFrontmatterInputs {
		privateInputs = append(privateInputs, privateInput)
	}
	sort.Strings(privateInputs)
	for _, privateInput := range privateInputs {
		input.WriteString("private-frontmatter")
		input.WriteByte('\x00')
		input.WriteString(privateInput)
		input.WriteByte('\x00')
	}

	localMediaPaths := imageindex.LocalMediaPaths(m.Posts(), imageindex.BuildOptions{
		ContentDir: contentDir,
		AssetsDir:  assetsDir,
		OutputDir:  outputDir,
	})
	mediaHash, _, err := hashImageLibrarySources(contentDir, assetsDir, outputDir, config.ShouldIncludeUnreferenced(), localMediaPaths, imageHashExtensions(), GetBuildCache(m))
	if err != nil {
		return "", err
	}
	input.WriteString(mediaHash)
	return buildcache.ContentHash(input.String()), nil
}

func (p *ImageLibraryPlugin) imageLibraryTemplateHash(m *lifecycle.Manager, templateName string) (string, error) {
	var input strings.Builder
	input.WriteString(imageLibraryGeneratorVersion(m.Config()))
	input.WriteByte('\x00')
	if cache := GetBuildCache(m); cache != nil {
		input.WriteString(cache.GetTemplatesHash())
		input.WriteByte('\x00')
	}

	engine, err := p.templateEngine(m.Config())
	if err != nil {
		return "", err
	}
	paths := engine.SearchPaths()
	for _, path := range paths {
		hash, err := buildcache.HashDirectory(path, []string{".html", ".txt", ".md"})
		if err != nil {
			return "", err
		}
		input.WriteString(filepath.Clean(path))
		input.WriteByte('=')
		input.WriteString(hash)
		input.WriteByte('\x00')
	}
	if len(paths) == 0 && engine.HasEmbeddedTemplate(templateName) {
		input.WriteString("embedded:")
		input.WriteString(templateName)
	}
	return buildcache.ContentHash(input.String()), nil
}

func frontmatterImageHash(post *models.Post) string {
	if post == nil || post.Extra == nil {
		return ""
	}
	keys := make([]string, 0, len(post.Extra))
	for key := range post.Extra {
		for _, field := range []string{"image", "video", "cover", "cover_image", "og_image", "social_image", "thumbnail", "featured_image", "hero_image", "avatar", "author_image", "poster_image", "poster", "video_poster", "video_thumbnail", "thumb", "image_alt", "video_alt", "cover_alt", "cover_image_alt", "alt"} {
			if key == field {
				keys = append(keys, key)
				break
			}
		}
	}
	sort.Strings(keys)
	var result strings.Builder
	for _, key := range keys {
		result.WriteString(key)
		result.WriteByte('=')
		result.WriteString(fmt.Sprint(post.Extra[key]))
		result.WriteByte('\x00')
	}
	return buildcache.ContentHash(result.String())
}

func imageHashExtensions() []string {
	return []string{
		".avif", ".avi", ".bmp", ".gif", ".heic", ".heif", ".ico", ".jpeg", ".jpg", ".m4v", ".mkv", ".mov", ".mp4", ".ogv", ".ogg", ".png", ".svg", ".tif", ".tiff", ".webm", ".webp",
		".AVIF", ".AVI", ".BMP", ".GIF", ".HEIC", ".HEIF", ".ICO", ".JPEG", ".JPG", ".M4V", ".MKV", ".MOV", ".MP4", ".OGV", ".OGG", ".PNG", ".SVG", ".TIF", ".TIFF", ".WEBM", ".WEBP",
	}
}

// hashImageDirectory fingerprints regular image and video files while ignoring
// symbolic links. It is retained as an uncached helper for callers and tests
// that need to hash one directory.
func hashImageDirectory(dir string, extensions []string) (contentHash, stateHash string, err error) {
	return hashImageDirectoriesWithHasher(dir, "", extensions, nil, buildcache.HashFile)
}

// hashImageLibrarySources hashes the configured public asset root when the
// inventory includes unreferenced files, plus local media files referenced by
// eligible posts. It intentionally does not walk the whole content root:
// private-only media and unrelated content-local files must not become
// image-library cache inputs. Referenced files use stable assets/ or content/
// logical paths regardless of their absolute filesystem location.
func hashImageLibrarySources(contentDir, assetsDir, outputDir string, includeUnreferenced bool, sourceFiles, extensions []string, cache *buildcache.Cache) (contentHash, stateHash string, err error) {
	return hashImageLibrarySourcesWithHasher(contentDir, assetsDir, outputDir, includeUnreferenced, sourceFiles, extensions, cache, buildcache.HashFile)
}

func hashImageLibrarySourcesWithHasher(contentDir, assetsDir, outputDir string, includeUnreferenced bool, sourceFiles, extensions []string, cache *buildcache.Cache, hashFile func(string) (string, error)) (contentHash, stateHash string, err error) {
	if hashFile == nil {
		hashFile = buildcache.HashFile
	}
	extensionsByName := imageLibraryExtensionsByName(extensions)
	excludedRoots := imageLibraryExcludedRoots([]string{outputDir})
	filesByPath := make(map[string]imageLibraryFile, len(sourceFiles)+16)
	if includeUnreferenced {
		for _, root := range imageLibraryScanRoots("", assetsDir) {
			if walkErr := collectImageLibraryFiles(root, extensionsByName, excludedRoots, filesByPath); os.IsNotExist(walkErr) {
				continue
			} else if walkErr != nil {
				return "", "", walkErr
			}
		}
	}

	contentRoot := normalizedImageLibraryDirectory(contentDir)
	assetsRoot := ""
	if assetsDir != "" {
		assetsRoot = normalizedImageLibraryDirectory(assetsDir)
	}
	for _, sourcePath := range sourceFiles {
		absolute := normalizedImageLibraryPath(sourcePath)
		if imageLibraryPathExcluded(absolute, excludedRoots) || (includeUnreferenced && imageLibraryPathWithin(assetsRoot, absolute)) {
			continue
		}
		logicalPath, ok := imageLibraryLogicalSourcePath(contentRoot, assetsRoot, absolute)
		if !ok {
			continue
		}
		if err := collectImageLibraryFile(absolute, logicalPath, extensionsByName, excludedRoots, filesByPath); err != nil {
			return "", "", err
		}
	}

	files := make([]imageLibraryFile, 0, len(filesByPath))
	for _, file := range filesByPath {
		files = append(files, file)
	}
	sort.Slice(files, func(i, j int) bool { return files[i].hashPath < files[j].hashPath })
	previous := map[string]buildcache.ImageLibraryMediaFingerprint(nil)
	if cache != nil {
		previous = cache.GetImageLibraryMedia()
	}
	contentHash, stateHash, next, err := hashImageLibraryFiles(files, previous, hashFile)
	if err != nil {
		return "", "", err
	}
	if cache != nil {
		cache.SetImageLibraryMedia(next)
	}
	return contentHash, stateHash, nil
}

type imageLibraryScanRoot struct {
	path  string
	label string
}

type imageLibraryFile struct {
	cacheKey        string
	hashPath        string
	fullPath        string
	size            int64
	modTime         int64
	changeTime      int64
	changeTimeKnown bool
}

func hashImageDirectoriesWithHasher(contentDir, assetsDir string, extensions []string, cache *buildcache.Cache, hashFile func(string) (string, error), excludedDirs ...string) (contentHash, stateHash string, err error) {
	if hashFile == nil {
		hashFile = buildcache.HashFile
	}

	roots := imageLibraryScanRoots(contentDir, assetsDir)
	excludedRoots := imageLibraryExcludedRoots(excludedDirs)
	extensionsByName := imageLibraryExtensionsByName(extensions)
	filesByPath := make(map[string]imageLibraryFile, 64)
	for _, root := range roots {
		if walkErr := collectImageLibraryFiles(root, extensionsByName, excludedRoots, filesByPath); os.IsNotExist(walkErr) {
			continue
		} else if walkErr != nil {
			return "", "", walkErr
		}
	}

	files := make([]imageLibraryFile, 0, len(filesByPath))
	for _, file := range filesByPath {
		files = append(files, file)
	}
	sort.Slice(files, func(i, j int) bool { return files[i].hashPath < files[j].hashPath })

	previous := map[string]buildcache.ImageLibraryMediaFingerprint(nil)
	if cache != nil {
		previous = cache.GetImageLibraryMedia()
	}
	contentHash, stateHash, next, err := hashImageLibraryFiles(files, previous, hashFile)
	if err != nil {
		return "", "", err
	}
	if cache != nil {
		cache.SetImageLibraryMedia(next)
	}
	return contentHash, stateHash, nil
}

func imageLibraryScanRoots(contentDir, assetsDir string) []imageLibraryScanRoot {
	roots := make([]imageLibraryScanRoot, 0, 2)
	contentRoot := ""
	if contentDir != "" {
		contentRoot = normalizedImageLibraryDirectory(contentDir)
		roots = append(roots, imageLibraryScanRoot{path: contentRoot, label: "content"})
	}
	if assetsDir == "" {
		return roots
	}
	assetsRoot := normalizedImageLibraryDirectory(assetsDir)
	if contentRoot == "" {
		return []imageLibraryScanRoot{{path: assetsRoot, label: "assets"}}
	}
	switch {
	case imageLibraryPathWithin(contentRoot, assetsRoot):
		// The content walk already includes the configured asset directory.
	case imageLibraryPathWithin(assetsRoot, contentRoot):
		// Prefer the outer asset root when it contains the content root.
		roots[0] = imageLibraryScanRoot{path: assetsRoot, label: "assets"}
	default:
		roots = append(roots, imageLibraryScanRoot{path: assetsRoot, label: "assets"})
	}
	return roots
}

func imageLibraryExtensionsByName(extensions []string) map[string]struct{} {
	byName := make(map[string]struct{}, len(extensions))
	for _, extension := range extensions {
		byName[strings.ToLower(extension)] = struct{}{}
	}
	return byName
}

func collectImageLibraryFiles(root imageLibraryScanRoot, extensionsByName map[string]struct{}, excludedRoots []string, filesByPath map[string]imageLibraryFile) error {
	return filepath.WalkDir(root.path, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if imageLibraryPathExcluded(path, excludedRoots) {
			if entry.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		relative, relativeErr := filepath.Rel(root.path, path)
		if relativeErr != nil {
			relative = path
		}
		return collectImageLibraryFile(path, root.label+"/"+filepath.ToSlash(relative), extensionsByName, excludedRoots, filesByPath)
	})
}

func collectImageLibraryFile(path, hashPath string, extensionsByName map[string]struct{}, excludedRoots []string, filesByPath map[string]imageLibraryFile) error {
	if imageLibraryPathExcluded(path, excludedRoots) {
		return nil
	}
	info, err := os.Lstat(path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	if info.Mode()&os.ModeSymlink != 0 || !info.Mode().IsRegular() {
		return nil
	}
	if _, ok := extensionsByName[strings.ToLower(filepath.Ext(path))]; !ok {
		return nil
	}
	absolute := normalizedImageLibraryPath(path)
	cacheKey := filepath.ToSlash(absolute)
	if _, exists := filesByPath[cacheKey]; exists {
		return nil
	}
	changeTime, changeTimeKnown := imageLibraryChangeTime(info)
	filesByPath[cacheKey] = imageLibraryFile{
		cacheKey:        cacheKey,
		hashPath:        hashPath,
		fullPath:        absolute,
		size:            info.Size(),
		modTime:         info.ModTime().UnixNano(),
		changeTime:      changeTime,
		changeTimeKnown: changeTimeKnown,
	}
	return nil
}

func imageLibraryExcludedRoots(paths []string) []string {
	result := make([]string, 0, len(paths))
	seen := make(map[string]struct{}, len(paths))
	for _, path := range paths {
		if path == "" {
			continue
		}
		normalized := normalizedImageLibraryDirectory(path)
		if _, exists := seen[normalized]; exists {
			continue
		}
		seen[normalized] = struct{}{}
		result = append(result, normalized)
	}
	return result
}

func imageLibraryPathExcluded(path string, excludedRoots []string) bool {
	normalized := normalizedImageLibraryPath(path)
	for _, root := range excludedRoots {
		if imageLibraryPathWithin(root, normalized) {
			return true
		}
	}
	return false
}

func hashImageLibraryFiles(files []imageLibraryFile, previous map[string]buildcache.ImageLibraryMediaFingerprint, hashFile func(string) (string, error)) (contentHash, stateHash string, next map[string]buildcache.ImageLibraryMediaFingerprint, err error) {
	next = make(map[string]buildcache.ImageLibraryMediaFingerprint, len(files))
	contentHasher := sha256.New()
	stateHasher := sha256.New()
	for _, file := range files {
		fingerprint, cached := previous[file.cacheKey]
		contentFileHash := fingerprint.ContentHash
		if imageLibraryMediaNeedsHash(cached, fingerprint, file) {
			contentFileHash, err = hashFile(file.fullPath)
			if err != nil {
				return "", "", nil, fmt.Errorf("hash media file %s: %w", file.fullPath, err)
			}
		}
		next[file.cacheKey] = buildcache.ImageLibraryMediaFingerprint{
			Size:            file.size,
			ModTime:         file.modTime,
			ChangeTime:      file.changeTime,
			ChangeTimeKnown: file.changeTimeKnown,
			ContentHash:     contentFileHash,
		}

		contentHasher.Write([]byte(file.hashPath))
		contentHasher.Write([]byte{0})
		contentHasher.Write([]byte(contentFileHash))
		contentHasher.Write([]byte{0})
		stateHasher.Write([]byte(file.hashPath))
		stateHasher.Write([]byte{0})
		stateHasher.Write([]byte(strconv.FormatInt(file.size, 10)))
		stateHasher.Write([]byte{0})
		stateHasher.Write([]byte(strconv.FormatInt(file.modTime, 10)))
		stateHasher.Write([]byte{0})
	}
	return hex.EncodeToString(contentHasher.Sum(nil)), hex.EncodeToString(stateHasher.Sum(nil)), next, nil
}

func imageLibraryMediaNeedsHash(cached bool, fingerprint buildcache.ImageLibraryMediaFingerprint, file imageLibraryFile) bool {
	return !cached || fingerprint.Size != file.size || fingerprint.ModTime != file.modTime || fingerprint.ChangeTimeKnown != file.changeTimeKnown || (file.changeTimeKnown && fingerprint.ChangeTime != file.changeTime) || fingerprint.ContentHash == ""
}

func normalizedImageLibraryDirectory(path string) string {
	if path == "" {
		path = "."
	}
	return normalizedImageLibraryPath(path)
}

func normalizedImageLibraryPath(path string) string {
	absolute, err := filepath.Abs(path)
	if err != nil {
		absolute = filepath.Clean(path)
	}
	return filepath.Clean(absolute)
}

func imageLibraryPathWithin(root, path string) bool {
	if root == "" {
		return false
	}
	relative, err := filepath.Rel(root, path)
	return err == nil && relative != ".." && !strings.HasPrefix(relative, ".."+string(filepath.Separator)) && !filepath.IsAbs(relative)
}

// imageLibraryLogicalSourcePath maps a resolved local media path to the stable
// namespace used by the semantic aggregate hash. Asset paths take precedence
// when the configured roots overlap so an asset keeps the same identity whether
// it is discovered by an asset walk or by an explicit public reference.
func imageLibraryLogicalSourcePath(contentRoot, assetsRoot, path string) (string, bool) {
	for _, root := range []struct {
		path      string
		namespace string
	}{
		{path: assetsRoot, namespace: "assets"},
		{path: contentRoot, namespace: "content"},
	} {
		if root.path == "" {
			continue
		}
		relative, err := filepath.Rel(root.path, path)
		if err != nil || relative == "." || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) || filepath.IsAbs(relative) {
			continue
		}
		return root.namespace + "/" + filepath.ToSlash(relative), true
	}
	return "", false
}

// imageLibraryChangeTime extracts the portable shape shared by Unix stat
// structures without importing platform-specific syscall types. Filesystems
// that do not expose a change time return false and use size plus mtime only.
func imageLibraryChangeTime(info os.FileInfo) (int64, bool) {
	if info == nil || info.Sys() == nil {
		return 0, false
	}
	value := reflect.ValueOf(info.Sys())
	for value.Kind() == reflect.Ptr {
		if value.IsNil() {
			return 0, false
		}
		value = value.Elem()
	}
	if value.Kind() != reflect.Struct {
		return 0, false
	}
	for _, fieldName := range []string{"Ctim", "Ctimespec", "ChangeTime"} {
		field := value.FieldByName(fieldName)
		if !field.IsValid() || field.Kind() != reflect.Struct {
			continue
		}
		seconds, secondsOK := imageLibraryIntField(field.FieldByName("Sec"))
		nanos, nanosOK := imageLibraryIntField(field.FieldByName("Nsec"))
		if !secondsOK || !nanosOK {
			continue
		}
		return seconds*int64(time.Second) + nanos, true
	}
	return 0, false
}

func imageLibraryIntField(value reflect.Value) (int64, bool) {
	if !value.IsValid() || !value.CanInt() {
		return 0, false
	}
	return value.Int(), true
}

func imageLibraryOutputsExist(pagePath, jsonPath, flatJSONPath string, exportJSON bool) bool {
	if !isRegularImageLibraryFile(pagePath) {
		return false
	}
	if exportJSON {
		if !isRegularImageLibraryFile(jsonPath) || !isRegularImageLibraryFile(flatJSONPath) {
			return false
		}
	}
	return true
}

func isRegularImageLibraryFile(path string) bool {
	info, err := os.Lstat(path)
	return err == nil && info.Mode().IsRegular()
}

func imageLibraryPathExists(path string) bool {
	_, err := os.Lstat(path)
	return err == nil
}

func imageLibraryPathsOverlap(left, right string) bool {
	left = filepath.Clean(left)
	right = filepath.Clean(right)
	if left == right {
		return true
	}
	for _, pair := range [][2]string{{left, right}, {right, left}} {
		relative, err := filepath.Rel(pair[1], pair[0])
		if err != nil || relative == "." || relative == ".." || filepath.IsAbs(relative) || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
			continue
		}
		return true
	}
	return false
}

func imageLibraryOwnsOutput(cache *buildcache.Cache, contentRoot, outputRoot, path string) bool {
	return cache != nil && cache.OwnsImageLibraryOutput(contentRoot, outputRoot, path)
}

func imageLibraryOutputMatches(cache *buildcache.Cache, contentRoot, outputRoot, path string) bool {
	if !imageLibraryOwnsOutput(cache, contentRoot, outputRoot, path) || !isRegularImageLibraryFile(path) {
		return false
	}
	expected, ok := cache.ImageLibraryOutputHash(contentRoot, outputRoot, path)
	if !ok || expected == "" {
		return false
	}
	actual, err := buildcache.HashFile(path)
	return err == nil && actual == expected
}

func imageLibraryOutputHashesForPaths(paths []string) (map[string]string, error) {
	hashes := make(map[string]string, len(paths))
	for _, path := range paths {
		if !isRegularImageLibraryFile(path) {
			continue
		}
		hash, err := buildcache.HashFile(path)
		if err != nil {
			return nil, err
		}
		hashes[path] = hash
	}
	return hashes, nil
}

// imageLibraryPagePaths validates every new route before any index or page is
// published. Page 1 was already checked by the caller.
func imageLibraryPagePaths(m *lifecycle.Manager, cache *buildcache.Cache, contentRoot, outputRoot, pagePath, jsonPath, flatJSONPath string, count int, exportJSON bool) ([]string, error) {
	paths := []string{pagePath}
	for number := 2; number <= count; number++ {
		path := filepath.Join(filepath.Dir(pagePath), "page", strconv.Itoa(number), "index.html")
		if err := imageLibraryValidateOutputPath(outputRoot, path); err != nil {
			return nil, err
		}
		if err := imageLibraryPostOutputConflict(m, outputRoot, path, flatJSONPath, exportJSON); err != nil {
			return nil, err
		}
		if err := imageLibraryFeedOutputConflict(m, outputRoot, path, jsonPath, flatJSONPath, exportJSON); err != nil {
			return nil, err
		}
		if err := imageLibraryExistingOutputConflict(cache, contentRoot, outputRoot, path, true); err != nil {
			return nil, err
		}
		paths = append(paths, path)
	}
	return paths, nil
}

func imageLibraryTrackedOutputsMatch(cache *buildcache.Cache, contentRoot, outputRoot string, previous map[string]string) bool {
	if cache == nil || len(previous) == 0 {
		return false
	}
	for relative := range previous {
		path, err := imageLibraryOutputPath(outputRoot, relative)
		if err != nil || !imageLibraryOutputMatches(cache, contentRoot, outputRoot, path) {
			return false
		}
	}
	return true
}

func imageLibraryStaleOutputConflict(cache *buildcache.Cache, contentRoot, pagePath, jsonPath, flatJSONPath string) error {
	if cache == nil {
		return nil
	}
	previousContentRoot, previousRoot, outputHashes := cache.GetImageLibraryOutputs()
	activeContentRoot, err := filepath.Abs(contentRoot)
	if err != nil || filepath.Clean(previousContentRoot) != filepath.Clean(activeContentRoot) {
		return nil
	}
	if previousRoot == "" {
		return nil
	}
	previousRoot, err = filepath.Abs(previousRoot)
	if err != nil {
		return nil
	}
	currentPage, err := filepath.Abs(pagePath)
	if err != nil {
		return err
	}
	currentJSON, err := filepath.Abs(jsonPath)
	if err != nil {
		return err
	}
	currentFlatJSON, err := filepath.Abs(flatJSONPath)
	if err != nil {
		return err
	}
	current := map[string]struct{}{
		filepath.Clean(currentPage):     {},
		filepath.Clean(currentJSON):     {},
		filepath.Clean(currentFlatJSON): {},
	}
	for relative, expected := range outputHashes {
		path, err := imageLibraryOutputPath(previousRoot, relative)
		if err != nil {
			return fmt.Errorf("invalid cached image library output %q: %w", relative, err)
		}
		if _, ok := current[filepath.Clean(path)]; ok || !imageLibraryPathExists(path) {
			continue
		}
		if !isRegularImageLibraryFile(path) {
			return fmt.Errorf("image library output is no longer a regular file: %s", path)
		}
		actual, err := buildcache.HashFile(path)
		if err != nil {
			return fmt.Errorf("hash image library output %s: %w", path, err)
		}
		if actual != expected {
			return fmt.Errorf("image library output was modified outside the images plugin: %s", path)
		}
	}
	return nil
}

func (p *ImageLibraryPlugin) renderPage(m *lifecycle.Manager, config *ImageLibraryConfig, pathPrefix string, index imageindex.Index, page ImageLibraryPage, outputFS *os.Root, outputRoot, pagePath string) error {
	engine, err := p.templateEngine(m.Config())
	if err != nil {
		return err
	}
	engine.ClearCache()
	// Config maps are cached by pointer. Clear the cache because configuration
	// can be updated in place between incremental writes.
	templates.ClearConfigMapCache()
	if !engine.TemplateExists(config.Template) {
		log.Printf("[images] Warning: template %q not found, skipping image library page", config.Template)
		if err := removeImageLibraryFileFromRoot(outputFS, outputRoot, pagePath); err != nil {
			return fmt.Errorf("remove image library page: %w", err)
		}
		return nil
	}

	if config.ShouldExportJSON() {
		page.JSONHref = "/" + strings.Trim(pathPrefix, "/") + "/index.json"
	}
	title := "Image Library"
	description := "Browse the images and videos used by this site."
	syntheticSlug := pathPrefix
	if page.PageNumber > 1 {
		syntheticSlug += "/page/" + strconv.Itoa(page.PageNumber)
	}
	syntheticPost := &models.Post{Slug: syntheticSlug, Href: page.PageLinks[page.PageNumber-1].Href, Published: true, Extra: make(map[string]interface{})}
	syntheticPost.Title = &title
	syntheticPost.Description = &description
	modelsConfig, ok := getModelsConfig(m.Config())
	if !ok || modelsConfig == nil {
		modelsConfig = ToModelsConfig(m.Config())
	}
	if modelsConfig == nil {
		modelsConfig = &models.Config{}
	}
	ctx := templates.NewContext(syntheticPost, "", modelsConfig)
	ctx.Extra["title"] = title
	ctx.Extra["description"] = description
	ctx.Extra["image_library"] = page
	ctx.Extra["image_index"] = index
	ctx.Extra["image_library_config"] = *config
	ctx.Extra["needs_image_zoom"] = false

	html, err := engine.Render(config.Template, ctx)
	if err != nil {
		return fmt.Errorf("rendering image library template: %w", err)
	}
	if err := writeImageLibraryFile(outputFS, outputRoot, pagePath, []byte(html)); err != nil {
		return fmt.Errorf("writing image library page: %w", err)
	}
	return nil
}

func (p *ImageLibraryPlugin) templateEngine(config *lifecycle.Config) (*templates.Engine, error) {
	if config != nil && config.Extra != nil {
		if value, ok := config.Extra["templates.engine"]; ok {
			if engine, ok := value.(*templates.Engine); ok && engine != nil {
				return engine, nil
			}
		}
	}
	templatesDir := PluginNameTemplates
	if config != nil && config.Extra != nil {
		if value, ok := config.Extra["templates_dir"].(string); ok && value != "" {
			templatesDir = value
		}
	}
	themeName := getThemeName(config)
	cacheKey := templatesDir + ":" + themeName
	p.engineMu.RLock()
	engine := p.engineCache[cacheKey]
	p.engineMu.RUnlock()
	if engine != nil {
		return engine, nil
	}
	p.engineMu.Lock()
	defer p.engineMu.Unlock()
	if cached := p.engineCache[cacheKey]; cached != nil {
		return cached, nil
	}
	var err error
	engine, err = templates.NewEngineWithTheme(templatesDir, themeName)
	if err != nil {
		return nil, err
	}
	p.engineCache[cacheKey] = engine
	return engine, nil
}

// ImageLibraryPage is the template view model for the generated page.
type ImageLibraryPage struct {
	Images      []ImageCard
	PageNumber  int
	TotalPages  int
	BaseHref    string
	PageLinks   []ImageLibraryPageLink
	PreviousURL string
	NextURL     string
	Count       int
	UsedCount   int
	UnusedCount int
	CoverCount  int
	JSONHref    string
}

// ImageLibraryPageLink is a native, crawlable pagination destination.
type ImageLibraryPageLink struct {
	Number  int
	Href    string
	Current bool
}

// imageLibraryPaginate keeps the canonical index complete while emitting
// a bounded number of server-rendered cards on each public HTML page.
func imageLibraryPaginate(full ImageLibraryPage, pageSize int, pathPrefix string) []ImageLibraryPage {
	total := len(full.Images)
	if pageSize <= 0 || pageSize >= total {
		pageSize = total
	}
	if pageSize == 0 {
		pageSize = 1
	}
	totalPages := (total + pageSize - 1) / pageSize
	if totalPages == 0 {
		totalPages = 1
	}
	base := "/" + strings.Trim(pathPrefix, "/") + "/"
	links := make([]ImageLibraryPageLink, totalPages)
	for index := range links {
		href := base
		if index > 0 {
			href += "page/" + strconv.Itoa(index+1) + "/"
		}
		links[index] = ImageLibraryPageLink{Number: index + 1, Href: href}
	}
	pages := make([]ImageLibraryPage, 0, totalPages)
	for index := 0; index < totalPages; index++ {
		start := index * pageSize
		end := start + pageSize
		if end > total {
			end = total
		}
		page := full
		page.Images = full.Images[start:end]
		page.PageNumber = index + 1
		page.TotalPages = totalPages
		page.BaseHref = base
		page.PageLinks = append([]ImageLibraryPageLink(nil), links...)
		page.PageLinks[index].Current = true
		if index > 0 {
			page.PreviousURL = links[index-1].Href
		}
		if index+1 < totalPages {
			page.NextURL = links[index+1].Href
		}
		pages = append(pages, page)
	}
	return pages
}

// ImageCard contains the presentation-only values derived from one image.
type ImageCard struct {
	Anchor         string
	Src            string
	PreviewSrc     string
	PosterSrc      string
	Srcset         string
	Name           string
	Alt            string
	MIMEType       string
	Width          int
	Height         int
	AspectRatio    string
	AddedAt        string
	AddedAtUnix    int64
	LastUsedAt     string
	LastUsedAtUnix int64
	Markdown       string
	SearchText     string
	Cover          bool
	Embed          bool
	IsVideo        bool
	Used           bool
	Uses           []imageindex.Use
	MoreUses       []imageindex.Use
	MoreCount      int
	UseCount       int
}

func newImageLibraryPage(index imageindex.Index) ImageLibraryPage {
	page := ImageLibraryPage{Images: make([]ImageCard, 0, len(index.Images)), Count: len(index.Images)}
	for i := range index.Images {
		image := &index.Images[i]
		isVideo := imageIsVideo(*image)
		srcset := ""
		if !isVideo {
			srcset = imageSrcset(image.Src)
		}
		card := ImageCard{
			Anchor:      "image-" + strconv.Itoa(i+1),
			Src:         image.Src,
			PreviewSrc:  imagePreviewURL(image.Src),
			PosterSrc:   imagePosterURL(*image, isVideo),
			Srcset:      srcset,
			Name:        imageName(image.Src),
			Alt:         image.Alt,
			MIMEType:    image.MIMEType,
			Width:       image.Width,
			Height:      image.Height,
			AspectRatio: imageAspectRatio(image.Width, image.Height),
			Markdown:    markdownForImage(image.Src, image.Alt),
			SearchText:  imageSearchText(*image),
			Cover:       image.Cover,
			Embed:       image.Embed,
			IsVideo:     isVideo,
			Used:        len(image.Uses) > 0,
			UseCount:    len(image.Uses),
		}
		if image.AddedAt != nil {
			card.AddedAt = image.AddedAt.UTC().Format("Jan 2, 2006")
			card.AddedAtUnix = image.AddedAt.Unix()
		}
		if image.LastUsedAt != nil {
			card.LastUsedAt = image.LastUsedAt.UTC().Format("Jan 2, 2006")
			card.LastUsedAtUnix = image.LastUsedAt.Unix()
		}
		if len(image.Uses) > 3 {
			card.Uses = append([]imageindex.Use(nil), image.Uses[:3]...)
			card.MoreUses = append([]imageindex.Use(nil), image.Uses[3:]...)
			card.MoreCount = len(card.MoreUses)
		} else {
			card.Uses = append([]imageindex.Use(nil), image.Uses...)
		}
		page.Images = append(page.Images, card)
		if card.Used {
			page.UsedCount++
		} else {
			page.UnusedCount++
		}
		if card.Cover {
			page.CoverCount++
		}
	}
	return page
}

func imagePreviewURL(src string) string {
	parsed, err := url.Parse(src)
	if err != nil || parsed.Host == "" || !templates.IsTrustedMediaURL(src) {
		return src
	}
	return templates.WithSize(src, 640, 0)
}

func imageIsVideo(image imageindex.Image) bool {
	return strings.HasPrefix(strings.ToLower(strings.TrimSpace(image.MIMEType)), "video/") || templates.IsVideoURL(image.Src)
}

func imagePosterURL(image imageindex.Image, isVideo bool) string {
	if !isVideo {
		return ""
	}
	poster := strings.TrimSpace(image.PosterSrc)
	if poster == "" {
		poster = templates.PosterURLFromMap(map[string]interface{}{}, image.Src)
	}
	if poster == "" {
		return ""
	}
	return templates.WithSize(poster, 640, 0)
}

func imageSrcset(src string) string {
	if templates.IsVideoURL(src) {
		return ""
	}
	parsed, err := url.Parse(src)
	if err != nil || parsed.Host == "" || !templates.IsTrustedMediaURL(src) {
		return ""
	}
	widths := []int{320, 640, 960, 1280}
	values := make([]string, 0, len(widths))
	for _, width := range widths {
		values = append(values, fmt.Sprintf("%s %dw", templates.WithSize(src, width, 0), width))
	}
	return strings.Join(values, ", ")
}

func imageName(src string) string {
	parsed, err := url.Parse(src)
	if err == nil && parsed.Path != "" {
		name := filepath.Base(parsed.Path)
		if name != "." && name != string(filepath.Separator) && name != "" {
			return name
		}
	}
	return src
}

func imageAspectRatio(width, height int) string {
	if width > 0 && height > 0 {
		return fmt.Sprintf("%d / %d", width, height)
	}
	return "4 / 3"
}

func markdownForImage(src, alt string) string {
	alt = strings.ReplaceAll(strings.ReplaceAll(alt, `\`, `\\`), "]", `\]`)
	src = strings.ReplaceAll(strings.ReplaceAll(src, `\`, `\\`), ")", `\)`)
	return fmt.Sprintf("![%s](%s)", alt, src)
}

func imageSearchText(image imageindex.Image) string {
	values := []string{image.Src, image.Alt}
	if imageIsVideo(image) {
		values = append(values, "video")
	}
	if image.Embed {
		values = append(values, "embed")
	}
	for _, use := range image.Uses {
		values = append(values, use.Href, use.Title, use.Caption)
	}
	return strings.ToLower(strings.Join(values, " "))
}

var imageLibraryTempCounter uint64

func writeImageLibraryFile(outputFS *os.Root, outputRoot, path string, data []byte) error {
	relative, err := outputRelativePath(outputRoot, path)
	if err != nil {
		return err
	}
	directory := filepath.Dir(relative)
	if directory != "." {
		if err := outputFS.MkdirAll(directory, 0o755); err != nil {
			return err
		}
	}
	temporary, temporaryName, err := createImageLibraryTemp(outputFS, directory)
	if err != nil {
		return err
	}
	defer outputFS.Remove(temporaryName) //nolint:errcheck // best-effort cleanup after rename or failure
	if err := temporary.Chmod(0o644); err != nil {
		_ = temporary.Close()
		return err
	}
	if _, err := temporary.Write(data); err != nil {
		_ = temporary.Close()
		return err
	}
	if err := temporary.Close(); err != nil {
		return err
	}
	if runtime.GOOS == tailwindOSWindows {
		if err := removeImageLibraryFileFromRoot(outputFS, outputRoot, path); err != nil {
			return err
		}
	}
	return outputFS.Rename(temporaryName, relative)
}

func createImageLibraryTemp(outputFS *os.Root, directory string) (*os.File, string, error) {
	for attempt := 0; attempt < 100; attempt++ {
		name := fmt.Sprintf(".images-%d-%d.tmp", os.Getpid(), atomic.AddUint64(&imageLibraryTempCounter, 1))
		if directory != "." {
			name = filepath.Join(directory, name)
		}
		file, err := outputFS.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
		if err == nil {
			return file, name, nil
		}
		if !os.IsExist(err) {
			return nil, "", err
		}
	}
	return nil, "", fmt.Errorf("create temporary image library output: too many existing temporary files")
}

func removeImageLibraryFile(outputRoot, path string) error {
	outputFS, err := openExistingOutputRoot(outputRoot)
	if err != nil || outputFS == nil {
		return err
	}
	defer outputFS.Close()
	return removeImageLibraryFileFromRoot(outputFS, outputRoot, path)
}

func removeImageLibraryFileFromRoot(outputFS *os.Root, outputRoot, path string) error {
	relative, err := outputRelativePath(outputRoot, path)
	if err != nil {
		return err
	}
	if err := outputFS.Remove(relative); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func removeImageLibraryOutputs(contentRoot string, cache *buildcache.Cache) error {
	if cache == nil {
		return nil
	}
	previousContentRoot, previousRoot, outputHashes := cache.GetImageLibraryOutputs()
	activeContentRoot, err := filepath.Abs(contentRoot)
	if err != nil {
		return fmt.Errorf("resolve image library content directory: %w", err)
	}
	if filepath.Clean(previousContentRoot) != filepath.Clean(activeContentRoot) {
		cache.ClearImageLibraryOutputs()
		cache.SetImageLibraryHash("")
		return nil
	}
	if previousRoot == "" {
		cache.ClearImageLibraryOutputs()
		cache.SetImageLibraryHash("")
		return nil
	}
	previousRoot, err = filepath.Abs(previousRoot)
	if err != nil {
		return fmt.Errorf("resolve cached image library output directory: %w", err)
	}
	outputFS, err := openExistingOutputRoot(previousRoot)
	if err != nil {
		return fmt.Errorf("open cached image library output directory: %w", err)
	}
	if outputFS == nil {
		cache.ClearImageLibraryOutputs()
		cache.SetImageLibraryHash("")
		return nil
	}
	defer outputFS.Close()
	for relative, expected := range outputHashes {
		path, err := imageLibraryOutputPath(previousRoot, relative)
		if err != nil {
			return fmt.Errorf("invalid cached image library output %q: %w", relative, err)
		}
		if !imageLibraryPathExists(path) {
			continue
		}
		if !isRegularImageLibraryFile(path) {
			return fmt.Errorf("image library output is no longer a regular file: %s", path)
		}
		actual, err := buildcache.HashFile(path)
		if err != nil {
			return fmt.Errorf("hash disabled image library output %s: %w", path, err)
		}
		if actual != expected {
			return fmt.Errorf("image library output was modified outside the images plugin: %s", path)
		}
		if err := removeImageLibraryFileFromRoot(outputFS, previousRoot, path); err != nil {
			return fmt.Errorf("remove disabled image library output %s: %w", path, err)
		}
	}
	cache.ClearImageLibraryOutputs()
	cache.SetImageLibraryHash("")
	return nil
}

func imageLibraryOutputPath(outputRoot, relative string) (string, error) {
	if relative == "" || filepath.IsAbs(relative) {
		return "", fmt.Errorf("path must be relative to output root")
	}
	cleaned := filepath.Clean(relative)
	if cleaned == "." || cleaned == ".." || strings.HasPrefix(cleaned, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("path escapes output root")
	}
	path := filepath.Join(outputRoot, cleaned)
	if err := imageLibraryValidateOutputPath(outputRoot, path); err != nil {
		return "", err
	}
	return path, nil
}

func imageLibraryValidateOutputPath(outputRoot, path string) error {
	return validateOutputPath(outputRoot, path)
}

func removeStaleImageLibraryOutputs(previousContentRoot, previousRoot string, previous map[string]string, contentRoot string, pagePaths []string) error {
	activeContentRoot, err := filepath.Abs(contentRoot)
	if err != nil {
		return err
	}
	if filepath.Clean(previousContentRoot) != filepath.Clean(activeContentRoot) {
		return nil
	}
	if previousRoot == "" {
		return nil
	}
	previousRoot, err = filepath.Abs(previousRoot)
	if err != nil {
		return err
	}
	outputFS, err := openExistingOutputRoot(previousRoot)
	if err != nil {
		return fmt.Errorf("open cached image library output directory: %w", err)
	}
	if outputFS == nil {
		return nil
	}
	defer outputFS.Close()
	current := make(map[string]struct{}, len(pagePaths))
	for _, path := range pagePaths {
		absolute, err := filepath.Abs(path)
		if err != nil {
			return err
		}
		current[filepath.Clean(absolute)] = struct{}{}
	}
	for relative, expected := range previous {
		path, err := imageLibraryOutputPath(previousRoot, relative)
		if err != nil {
			return err
		}
		if _, ok := current[path]; ok {
			continue
		}
		if !imageLibraryPathExists(path) {
			continue
		}
		if !isRegularImageLibraryFile(path) {
			return fmt.Errorf("image library output is no longer a regular file: %s", path)
		}
		actual, err := buildcache.HashFile(path)
		if err != nil {
			return err
		}
		if actual != expected {
			return fmt.Errorf("image library output was modified outside the images plugin: %s", path)
		}
		if err := removeImageLibraryFileFromRoot(outputFS, previousRoot, path); err != nil {
			return err
		}
	}
	return nil
}

// Ensure ImageLibraryPlugin implements the lifecycle contracts.
var (
	_ lifecycle.Plugin         = (*ImageLibraryPlugin)(nil)
	_ lifecycle.WritePlugin    = (*ImageLibraryPlugin)(nil)
	_ lifecycle.PriorityPlugin = (*ImageLibraryPlugin)(nil)
)
