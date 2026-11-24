// Package plugins provides lifecycle plugins for markata-go.
package plugins

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/WaylonWalker/markata/pkg/buildcache"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/templates"
)

// TagInfo represents information about a single tag for the tags listing page.
type TagInfo struct {
	// Name is the display name of the tag
	Name string

	// Slug is the URL-safe slug of the tag
	Slug string

	// Count is the number of posts with this tag
	Count int

	// Href is the URL to the tag page
	Href string
}

// TagsListingPlugin generates a tags listing page at /tags showing all available tags.
// It supports blacklist and private tag configurations to control tag visibility.
type TagsListingPlugin struct {
	engineMu    sync.RWMutex
	engineCache map[string]*templates.Engine
}

// NewTagsListingPlugin creates a new TagsListingPlugin.
func NewTagsListingPlugin() *TagsListingPlugin {
	return &TagsListingPlugin{
		engineCache: make(map[string]*templates.Engine),
	}
}

// Name returns the unique name of the plugin.
func (p *TagsListingPlugin) Name() string {
	return "tags_listing"
}

// Priority returns the plugin's priority for a given stage.
func (p *TagsListingPlugin) Priority(stage lifecycle.Stage) int {
	switch stage {
	case lifecycle.StageWrite:
		// Run after publish_feeds so tag pages exist
		return lifecycle.PriorityLate
	default:
		return lifecycle.PriorityDefault
	}
}

// getOrCreateEngine returns a cached template engine, or creates one if not cached.
func (p *TagsListingPlugin) getOrCreateEngine(templatesDir, themeName string) (*templates.Engine, error) {
	cacheKey := templatesDir + ":" + themeName

	// Fast path: check cache with read lock
	p.engineMu.RLock()
	if engine, ok := p.engineCache[cacheKey]; ok {
		p.engineMu.RUnlock()
		return engine, nil
	}
	p.engineMu.RUnlock()

	// Slow path: create engine with write lock
	p.engineMu.Lock()
	defer p.engineMu.Unlock()

	// Double-check after acquiring write lock
	if engine, ok := p.engineCache[cacheKey]; ok {
		return engine, nil
	}

	engine, err := templates.NewEngineWithTheme(templatesDir, themeName)
	if err != nil {
		return nil, err
	}

	p.engineCache[cacheKey] = engine
	return engine, nil
}

// Write generates the tags listing page.
func (p *TagsListingPlugin) Write(m *lifecycle.Manager) error {
	config := m.Config()
	cache := GetBuildCache(m)

	// Get tags config
	tagsConfig := getTagsConfig(config)
	if !tagsConfig.IsEnabled() {
		return nil
	}

	if cache != nil {
		if !cache.TagsDirty() {
			return nil
		}
		tagsHash := computeTagsListingHash(m.Posts(), &tagsConfig, cache)
		if cached := cache.GetTagsListingHash(); cached != "" && cached == tagsHash {
			return nil
		}
		cache.SetTagsListingHash(tagsHash)
	}

	// Collect and filter tags
	tagInfos := p.collectTags(m.Posts(), &tagsConfig)
	if len(tagInfos) == 0 {
		log.Printf("[tags_listing] No tags found, skipping tags listing page")
		return nil
	}

	// Sort tags alphabetically
	sort.Slice(tagInfos, func(i, j int) bool {
		return tagInfos[i].Name < tagInfos[j].Name
	})

	// Generate the tags listing page
	return p.renderTagsPage(config, &tagsConfig, tagInfos)
}

// collectTags gathers all visible tags from posts with their counts.
func (p *TagsListingPlugin) collectTags(posts []*models.Post, tagsConfig *models.TagsConfig) []TagInfo {
	tagCounts := make(map[string]int)

	for _, post := range posts {
		// Skip draft/unpublished/private posts
		if post.Draft || !post.Published || post.Private || post.Skip {
			continue
		}

		for _, tag := range post.Tags {
			// Skip blacklisted tags
			if tagsConfig.IsBlacklisted(tag) {
				continue
			}
			tagCounts[tag]++
		}
	}

	// Filter out private tags and build TagInfo list
	slugPrefix := tagsConfig.SlugPrefix
	if slugPrefix == "" {
		slugPrefix = defaultTagsPrefix
	}

	tagInfos := make([]TagInfo, 0, len(tagCounts))
	for tag, count := range tagCounts {
		// Skip private tags from listing
		if tagsConfig.IsPrivate(tag) {
			continue
		}

		slug := models.Slugify(tag)
		tagInfos = append(tagInfos, TagInfo{
			Name:  tag,
			Slug:  slug,
			Count: count,
			Href:  "/" + slugPrefix + "/" + slug + "/",
		})
	}

	return tagInfos
}

func computeTagsListingHash(posts []*models.Post, tagsConfig *models.TagsConfig, cache *buildcache.Cache) string {
	if tagsConfig == nil {
		return ""
	}
	var b strings.Builder
	b.WriteString(tagsConfig.SlugPrefix)
	b.WriteByte('\x00')
	b.WriteString(tagsConfig.Title)
	b.WriteByte('\x00')
	b.WriteString(tagsConfig.Description)
	b.WriteByte('\x00')
	b.WriteString(tagsConfig.Template)
	b.WriteByte('\x00')
	blacklist := append([]string(nil), tagsConfig.Blacklist...)
	sort.Strings(blacklist)
	for _, tag := range blacklist {
		b.WriteString("blacklist:")
		b.WriteString(tag)
		b.WriteByte('\x00')
	}
	privateTags := append([]string(nil), tagsConfig.Private...)
	sort.Strings(privateTags)
	for _, tag := range privateTags {
		b.WriteString("private:")
		b.WriteString(tag)
		b.WriteByte('\x00')
	}
	b.WriteString(fmt.Sprintf("%t", tagsConfig.IsEnabled()))
	b.WriteByte('\x00')

	for _, post := range posts {
		if post.Skip || post.Draft || !post.Published || post.Private {
			continue
		}
		if cache != nil {
			if _, tagHash, _ := cache.GetPostSemanticHashes(post.Path); tagHash != "" {
				b.WriteString(tagHash)
				b.WriteByte('\x00')
				continue
			}
		}
		b.WriteString(computePostTagIndexHash(post))
		b.WriteByte('\x00')
	}

	return buildcache.ContentHash(b.String())
}

// renderTagsPage renders and writes the tags listing HTML page.
func (p *TagsListingPlugin) renderTagsPage(config *lifecycle.Config, tagsConfig *models.TagsConfig, tagInfos []TagInfo) error {
	slugPrefix := tagsConfig.SlugPrefix
	if slugPrefix == "" {
		slugPrefix = defaultTagsPrefix
	}

	// Create output directory
	outputDir := config.OutputDir
	tagsDir := filepath.Join(outputDir, slugPrefix)
	if err := os.MkdirAll(tagsDir, 0o755); err != nil {
		return fmt.Errorf("creating tags directory: %w", err)
	}

	// Get template engine
	engine, err := p.createTemplateEngine(config)
	if err != nil {
		return err
	}

	// Determine template to use
	templateName := tagsConfig.Template
	if templateName == "" {
		templateName = "tags.html"
	}

	// Check if template exists
	if !engine.TemplateExists(templateName) {
		log.Printf("[tags_listing] Warning: template %q not found, skipping tags listing page", templateName)
		return nil
	}

	// Build context for template - create a synthetic post for the tags page
	modelsConfig := ToModelsConfig(config)
	title := tagsConfig.Title
	description := tagsConfig.Description
	syntheticPost := &models.Post{
		Slug:        tagsConfig.SlugPrefix,
		Title:       &title,
		Description: &description,
	}

	ctx := templates.NewContext(syntheticPost, "", modelsConfig)
	ctx.Extra["tag_list"] = tagInfos
	ctx.Extra["total_tags"] = len(tagInfos)

	// Render template
	html, err := engine.Render(templateName, ctx)
	if err != nil {
		return fmt.Errorf("rendering tags template: %w", err)
	}

	// Write output file
	outputPath := filepath.Join(tagsDir, "index.html")
	//nolint:gosec // G306: Output files need 0644 for web serving
	if err := os.WriteFile(outputPath, []byte(html), 0o644); err != nil {
		return fmt.Errorf("writing tags listing page: %w", err)
	}

	log.Printf("[tags_listing] Generated /tags/ with %d tags", len(tagInfos))

	return nil
}

// createTemplateEngine creates or retrieves a cached template engine.
func (p *TagsListingPlugin) createTemplateEngine(config *lifecycle.Config) (*templates.Engine, error) {
	templatesDir := PluginNameTemplates
	if extra, ok := config.Extra["templates_dir"].(string); ok && extra != "" {
		templatesDir = extra
	}

	themeName := getThemeName(config)

	return p.getOrCreateEngine(templatesDir, themeName)
}

// getThemeName extracts the theme name from config.
func getThemeName(config *lifecycle.Config) string {
	if config.Extra == nil {
		return ThemeDefault
	}

	if theme, ok := config.Extra["theme"].(models.ThemeConfig); ok && theme.Name != "" {
		return theme.Name
	}
	if theme, ok := config.Extra["theme"].(map[string]interface{}); ok {
		if name, ok := theme["name"].(string); ok && name != "" {
			return name
		}
	}
	if name, ok := config.Extra["theme"].(string); ok && name != "" {
		return name
	}

	return ThemeDefault
}

// getTagsConfig retrieves tags configuration from the manager config.
func getTagsConfig(config *lifecycle.Config) models.TagsConfig {
	// Use ToModelsConfig to get the properly converted config
	modelsConfig := ToModelsConfig(config)
	if modelsConfig != nil {
		return modelsConfig.Tags
	}
	return models.NewTagsConfig()
}

// Ensure TagsListingPlugin implements the required interfaces.
var (
	_ lifecycle.Plugin         = (*TagsListingPlugin)(nil)
	_ lifecycle.WritePlugin    = (*TagsListingPlugin)(nil)
	_ lifecycle.PriorityPlugin = (*TagsListingPlugin)(nil)
)
