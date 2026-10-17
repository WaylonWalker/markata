package config

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"

	"github.com/WaylonWalker/markata-go/pkg/models"
	"github.com/WaylonWalker/markata-go/pkg/renderingcontract"
)

// configSource defines the interface for format-specific config structs (TOML, YAML, JSON)
// to provide data for building a models.Config.
type configSource interface {
	getBaseConfig() baseConfigData
	getNavItems() []navItemData
	getFeeds() []feedConfigConverter
	getFeedDefaults() feedDefaultsConverter
	getPostFormats() postFormatsConverter
	getWellKnown() wellKnownConverter
	getSEO() seoConverter
	getIndieAuth() indieAuthConverter
	getWebmention() webmentionConverter
	getComponents() componentsConverter
	getSearch() models.SearchConfig
	getLayout() layoutConverter
	getSidebar() sidebarConverter
	getToc() tocConverter
	getHeader() headerConverter
	getBlogroll() blogrollConverter
	getTags() tagsConverter
	getFeedsPage() feedsPageConverter
	getEncryption() encryptionConverter
	getTagAggregator() tagAggregatorConverter
	getMentions() mentionsConverter
	getWebSub() webSubConverter
	getShortcuts() shortcutsConverter
	getViewTransitions() viewTransitionsConverter
	getAuthors() authorsConverter
	getGarden() gardenConverter
	getTemplates() templatesConverter
	getBuilderAdmin() models.BuilderAdminConfig
	getImages() models.ImagesConfig
}

type parsedThemeCalendar = models.ThemeCalendarConfig
type parsedHead struct {
	Text           any                    `json:"text,omitempty" yaml:"text,omitempty" toml:"text,omitempty"`
	Meta           []models.MetaTag       `json:"meta,omitempty" yaml:"meta,omitempty" toml:"meta,omitempty"`
	Link           []models.LinkTag       `json:"link,omitempty" yaml:"link,omitempty" toml:"link,omitempty"`
	Script         []models.ScriptTag     `json:"script,omitempty" yaml:"script,omitempty" toml:"script,omitempty"`
	AlternateFeeds []models.AlternateFeed `json:"alternate_feeds,omitempty" yaml:"alternate_feeds,omitempty" toml:"alternate_feeds,omitempty"`
}
type parsedTemplatePresets = map[string]models.TemplatePreset
type parsedDefaultTemplates = map[string]string
type parsedErrorPages = models.ErrorPagesConfig
type parsedResourceHints = models.ResourceHintsConfig
type parsedHighlight = models.HighlightConfig

func (h parsedHead) toHeadConfig() models.HeadConfig {
	return models.HeadConfig{
		Text:           normalizedHeadText(h.Text),
		Meta:           h.Meta,
		Link:           h.Link,
		Script:         h.Script,
		AlternateFeeds: h.AlternateFeeds,
	}
}

func normalizedHeadText(value any) string {
	text, err := parseHeadText(value)
	if err != nil {
		return ""
	}
	return text
}

func (h parsedHead) validate() error {
	_, err := parseHeadText(h.Text)
	return err
}

// parseHeadText accepts both the current string form and the legacy list form
// documented by HEAD_STYLE.md. The list form is used by existing site config
// files such as [[markata-go.head.text]] value = "...".
func parseHeadText(value any) (string, error) {
	switch value := value.(type) {
	case nil:
		return "", nil
	case string:
		return value, nil
	case []any:
		parts := make([]string, 0, len(value))
		for index, item := range value {
			text, err := parseHeadText(item)
			if err != nil {
				return "", fmt.Errorf("head.text block %d: %w", index, err)
			}
			if text != "" {
				parts = append(parts, text)
			}
		}
		return strings.Join(parts, "\n"), nil
	case []map[string]any:
		parts := make([]string, 0, len(value))
		for index, item := range value {
			text, err := parseHeadText(item)
			if err != nil {
				return "", fmt.Errorf("head.text block %d: %w", index, err)
			}
			if text != "" {
				parts = append(parts, text)
			}
		}
		return strings.Join(parts, "\n"), nil
	case map[string]any:
		text, ok := value["value"]
		if !ok {
			return "", fmt.Errorf("head.text block must contain a string value")
		}
		return parseHeadText(text)
	case map[string]string:
		text, ok := value["value"]
		if !ok {
			return "", fmt.Errorf("head.text block must contain a string value")
		}
		return text, nil
	case []string:
		return strings.Join(value, "\n"), nil
	case []map[string]string:
		parts := make([]string, 0, len(value))
		for index, item := range value {
			text, err := parseHeadText(item)
			if err != nil {
				return "", fmt.Errorf("head.text block %d: %w", index, err)
			}
			if text != "" {
				parts = append(parts, text)
			}
		}
		return strings.Join(parts, "\n"), nil
	default:
		return "", fmt.Errorf("head.text must be a string or list of string-valued blocks, got %T", value)
	}
}

// baseConfigData holds the basic config fields that are directly assignable.
type baseConfigData struct {
	Fontpack       string
	FontpacksFile  string
	OutputDir      string
	URL            string
	Title          string
	Description    string
	Author         string
	Language       string
	AuthorURL      string
	ManagingEditor string
	WebMaster      string
	Copyright      string
	License        interface{}
	AssetsDir      string
	Assets         models.AssetsConfig
	TemplatesDir   string
	Hooks          []string
	DisabledHooks  []string
	GlobPatterns   []string
	UseGitignore   *bool
	SlugMode       string
	SlugRules      []models.SlugRule
	Extensions     []string
	Concurrency    int
	Theme          models.ThemeConfig
	ThemeCalendar  parsedThemeCalendar
	Head           models.HeadConfig
	Presets        parsedTemplatePresets
	Defaults       parsedDefaultTemplates
	ErrorPages     parsedErrorPages
	ResourceHints  parsedResourceHints
	Highlight      parsedHighlight
	Footer         models.FooterConfig
}

// navItemData holds nav item fields.
type navItemData struct {
	Label    string
	URL      string
	External bool
}

// Converter interfaces for nested config types.
type feedConfigConverter interface {
	toFeedConfig() models.FeedConfig
}

type feedDefaultsConverter interface {
	toFeedDefaults() models.FeedDefaults
}

type postFormatsConverter interface {
	toPostFormatsConfig() models.PostFormatsConfig
}

type wellKnownConverter interface {
	toWellKnownConfig() models.WellKnownConfig
}

type seoConverter interface {
	toSEOConfig() models.SEOConfig
}

type indieAuthConverter interface {
	toIndieAuthConfig() models.IndieAuthConfig
}

type webmentionConverter interface {
	toWebmentionConfig() models.WebmentionConfig
}

type componentsConverter interface {
	toComponentsConfig() models.ComponentsConfig
}

type layoutConverter interface {
	toLayoutConfig() models.LayoutConfig
}

type sidebarConverter interface {
	toSidebarConfig() models.SidebarConfig
}

type tocConverter interface {
	toTocConfig() models.TocConfig
}

type headerConverter interface {
	toHeaderLayoutConfig() models.HeaderLayoutConfig
}

type blogrollConverter interface {
	toBlogrollConfig() models.BlogrollConfig
}

type tagsConverter interface {
	toTagsConfig() models.TagsConfig
}

type feedsPageConverter interface {
	toFeedsPageConfig() models.FeedsPageConfig
}

type encryptionConverter interface {
	toEncryptionConfig() models.EncryptionConfig
}

type tagAggregatorConverter interface {
	toTagAggregatorConfig() models.TagAggregatorConfig
}

type mentionsConverter interface {
	toMentionsConfig() models.MentionsConfig
}

type webSubConverter interface {
	toWebSubConfig() models.WebSubConfig
}

type shortcutsConverter interface {
	toShortcutsConfig() models.ShortcutsConfig
}

type viewTransitionsConverter interface {
	toViewTransitionsConfig() models.ViewTransitionsConfig
}

type authorsConverter interface {
	toAuthorsConfig() models.AuthorsConfig
}

type gardenConverter interface {
	toGardenConfig() models.GardenConfig
}

type templatesConverter interface {
	toTemplatesConfig() models.TemplatesConfig
}

// buildConfig constructs a models.Config from a configSource.
// This helper eliminates code duplication across TOML, YAML, and JSON config converters.
func buildConfig(src configSource) *models.Config {
	base := src.getBaseConfig()
	config := &models.Config{
		Fontpack:       base.Fontpack,
		FontpacksFile:  base.FontpacksFile,
		OutputDir:      base.OutputDir,
		URL:            base.URL,
		Title:          base.Title,
		Description:    base.Description,
		Author:         base.Author,
		Language:       base.Language,
		AuthorURL:      base.AuthorURL,
		ManagingEditor: base.ManagingEditor,
		WebMaster:      base.WebMaster,
		Copyright:      base.Copyright,
		AssetsDir:      base.AssetsDir,
		Assets:         base.Assets,
		TemplatesDir:   base.TemplatesDir,
		Hooks:          base.Hooks,
		DisabledHooks:  base.DisabledHooks,
		GlobConfig: models.GlobConfig{
			Patterns:  base.GlobPatterns,
			SlugMode:  base.SlugMode,
			SlugRules: append([]models.SlugRule{}, base.SlugRules...),
		},
		MarkdownConfig: models.MarkdownConfig{
			Extensions: base.Extensions,
			Highlight:  base.Highlight,
		},
		Concurrency: base.Concurrency,
		Theme:       base.Theme,

		ThemeCalendar:    base.ThemeCalendar,
		Head:             base.Head,
		TemplatePresets:  base.Presets,
		DefaultTemplates: base.Defaults,
		ErrorPages:       base.ErrorPages,
		ResourceHints:    base.ResourceHints,

		Footer: base.Footer,
	}
	config.License = models.LicenseValue{Raw: base.License}

	if base.UseGitignore != nil {
		config.GlobConfig.UseGitignore = *base.UseGitignore
	}

	// Convert nav items
	for _, nav := range src.getNavItems() {
		config.Nav = append(config.Nav, models.NavItem{
			Label:    nav.Label,
			URL:      nav.URL,
			External: nav.External,
		})
	}

	// Convert feeds
	for _, feed := range src.getFeeds() {
		config.Feeds = append(config.Feeds, feed.toFeedConfig())
	}

	// Convert feed defaults
	config.FeedDefaults = src.getFeedDefaults().toFeedDefaults()

	// Convert post formats
	config.PostFormats = src.getPostFormats().toPostFormatsConfig()

	// Convert well-known config
	config.WellKnown = src.getWellKnown().toWellKnownConfig()

	// Convert SEO config
	config.SEO = src.getSEO().toSEOConfig()

	// Convert IndieAuth and Webmention config
	config.IndieAuth = src.getIndieAuth().toIndieAuthConfig()
	config.Webmention = src.getWebmention().toWebmentionConfig()

	// Convert Components config
	config.Components = src.getComponents().toComponentsConfig()

	// Convert Search config
	config.Search = src.getSearch()

	// Convert Layout config
	config.Layout = src.getLayout().toLayoutConfig()

	// Convert Sidebar config
	config.Sidebar = src.getSidebar().toSidebarConfig()

	// Convert Toc config
	config.Toc = src.getToc().toTocConfig()

	// Convert Header config
	config.Header = src.getHeader().toHeaderLayoutConfig()

	// Convert Blogroll config
	config.Blogroll = src.getBlogroll().toBlogrollConfig()

	// Convert Tags config
	config.Tags = src.getTags().toTagsConfig()

	// Convert Feeds page config
	config.FeedsPage = src.getFeedsPage().toFeedsPageConfig()

	// Convert Encryption config
	config.Encryption = src.getEncryption().toEncryptionConfig()

	// Convert TagAggregator config
	config.TagAggregator = src.getTagAggregator().toTagAggregatorConfig()

	// Convert Mentions config
	config.Mentions = src.getMentions().toMentionsConfig()

	// Convert WebSub config
	config.WebSub = src.getWebSub().toWebSubConfig()

	// Convert Shortcuts config
	config.Shortcuts = src.getShortcuts().toShortcutsConfig()

	// Convert ViewTransitions config
	config.ViewTransitions = src.getViewTransitions().toViewTransitionsConfig()

	// Convert Authors config
	config.Authors = src.getAuthors().toAuthorsConfig()

	// Convert Garden config
	config.Garden = src.getGarden().toGardenConfig()

	// Convert Templates config
	config.Templates = src.getTemplates().toTemplatesConfig()

	// Convert image library config
	config.Images = src.getImages()

	// BuilderAdmin uses the shared model because it contains only scalar settings.
	config.BuilderAdmin = src.getBuilderAdmin()

	return config
}

func normalizeImagesConfig(config models.ImagesConfig) models.ImagesConfig {
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

// ParseTOML parses TOML configuration data into a Config struct.
// The TOML data is expected to have a top-level [markata-go] section.
func ParseTOML(data []byte) (*models.Config, error) {
	// Wrapper struct for the markata-go section
	var wrapper struct {
		MarkataGo tomlConfig `toml:"markata-go"`
	}

	if err := toml.Unmarshal(data, &wrapper); err != nil {
		return nil, err
	}
	if err := wrapper.MarkataGo.Head.validate(); err != nil {
		return nil, err
	}

	config := wrapper.MarkataGo.toConfig()

	// Also parse into a raw map to capture unknown plugin sections
	var rawWrapper map[string]any
	if err := toml.Unmarshal(data, &rawWrapper); err != nil {
		return config, nil // Non-fatal: continue with parsed config
	}

	// Extract the markata-go section as a map
	if markataGoRaw, ok := rawWrapper["markata-go"].(map[string]any); ok {
		normalizeParsedTheme(config, rawWrapper)
		if globRaw, ok := markataGoRaw["glob"].(map[string]any); ok && len(config.GlobConfig.SlugRules) == 0 {
			config.GlobConfig.SlugRules = extractSlugRules(globRaw["slug_rules"])
		}
		populateBuilderAdminExtra(config, markataGoRaw)

		// Initialize Extra if needed
		if config.Extra == nil {
			config.Extra = make(map[string]any)
		}

		// List of known top-level keys that are already parsed into struct fields
		knownKeys := map[string]bool{
			"output_dir": true, "url": true, "title": true, "description": true,
			"author": true, "license": true, "assets_dir": true, "templates_dir": true,
			"nav": true, "footer": true, "hooks": true, "disabled_hooks": true,
			"glob": true, "markdown": true, "feeds": true, "feed_defaults": true,
			"concurrency": true, "theme": true, "post_formats": true, "well_known": true,
			"seo": true, "indieauth": true, "webmention": true, "components": true,
			"layout": true, "sidebar": true, "toc": true, "header": true,
			"blogroll": true, "mentions": true, "template_presets": true,
			"slug_conflicts":    false,
			"default_templates": true, "auto_feeds": true, "head": true,
			"theme_calendar": true, "error_pages": true, "resource_hints": true,
			"content_templates": true, "footer_layout": true, "search": true,
			"plugins": true, "thoughts": true, "wikilinks": true, "tags": true,
			"tag_aggregator": true, "websub": true, "shortcuts": true, "view_transitions": true, "encryption": true,
			"authors": true, "garden": true, "builder_admin": true, "include": true, "tailwind": false, "css_purge": false,
			"assets": true, "images": true,
		}

		// Copy unknown sections to Extra
		for key, value := range markataGoRaw {
			if !knownKeys[key] {
				config.Extra[key] = value
			}
		}
	}

	return config, nil
}

// populateBuilderAdminExtra retains unrecognized builder_admin settings while
// excluding recognized authentication and webhook fields, including the secret.
func populateBuilderAdminExtra(config *models.Config, markataGoRaw map[string]any) {
	builderAdmin, ok := markataGoRaw["builder_admin"].(map[string]any)
	if !ok {
		return
	}

	extra := unknownMap(builderAdmin, map[string]bool{"auth": true, "webhook": true})
	if auth, ok := builderAdmin["auth"].(map[string]any); ok {
		if authExtra := unknownMap(auth, map[string]bool{"headers": true}); len(authExtra) > 0 {
			extra["auth"] = authExtra
		}
		if headers, ok := auth["headers"].(map[string]any); ok {
			if headersExtra := unknownMap(headers, map[string]bool{
				"user_id": true, "username": true, "display_name": true, "email": true,
				"groups": true, "roles": true, "scopes": true,
			}); len(headersExtra) > 0 {
				authExtra, ok := extra["auth"].(map[string]any)
				if !ok {
					authExtra = make(map[string]any)
					extra["auth"] = authExtra
				}
				authExtra["headers"] = headersExtra
			}
		}
	}
	if webhook, ok := builderAdmin["webhook"].(map[string]any); ok {
		if webhookExtra := unknownMap(webhook, map[string]bool{"enabled": true, "branch": true, "secret": true}); len(webhookExtra) > 0 {
			extra["webhook"] = webhookExtra
		}
	}
	if len(extra) > 0 {
		config.BuilderAdmin.Extra = extra
	}
}

func unknownMap(values map[string]any, known map[string]bool) map[string]any {
	result := make(map[string]any)
	for key, value := range values {
		if !known[key] {
			result[key] = value
		}
	}
	return result
}

func extractSlugRules(raw any) []models.SlugRule {
	entries, ok := raw.([]any)
	if !ok {
		return nil
	}

	rules := make([]models.SlugRule, 0, len(entries))
	for _, entry := range entries {
		m, ok := entry.(map[string]any)
		if !ok {
			continue
		}
		rules = append(rules, models.SlugRule{
			Prefix: stringValue(m["prefix"]),
			Mode:   stringValue(m["mode"]),
		})
	}
	return rules
}

func stringValue(v any) string {
	s, ok := v.(string)
	if !ok {
		return ""
	}
	return s
}

// ParseYAML parses YAML configuration data into a Config struct.
// The YAML data is expected to have a top-level markata-go key.
func ParseYAML(data []byte) (*models.Config, error) {
	// Wrapper struct for the markata-go section
	var wrapper struct {
		MarkataGo yamlConfig `yaml:"markata-go"`
	}

	if err := yaml.Unmarshal(data, &wrapper); err != nil {
		return nil, err
	}
	if err := wrapper.MarkataGo.Head.validate(); err != nil {
		return nil, err
	}

	config := wrapper.MarkataGo.toConfig()
	var raw map[string]any
	if err := yaml.Unmarshal(data, &raw); err == nil {
		normalizeParsedTheme(config, raw)
		copyLLMsTXTConfig(config, raw)
	}
	return config, nil
}

// copyLLMsTXTConfig retains the untyped llms_txt plugin options in YAML/JSON
// as well as TOML. Unlike TOML, these parsers do not automatically copy
// unknown [markata-go] plugin sections to Config.Extra.
func copyLLMsTXTConfig(config *models.Config, raw map[string]any) {
	root, ok := raw["markata-go"].(map[string]any)
	if !ok {
		return
	}
	llmsConfig, ok := root["llms_txt"].(map[string]any)
	if !ok {
		return
	}
	if config.Extra == nil {
		config.Extra = make(map[string]any)
	}
	config.Extra["llms_txt"] = llmsConfig
}

// ParseJSON parses JSON configuration data into a Config struct.
// The JSON data is expected to have a top-level "markata-go" key.
func ParseJSON(data []byte) (*models.Config, error) {
	// Wrapper struct for the markata-go section
	var wrapper struct {
		MarkataGo jsonConfig `json:"markata-go"`
	}

	if err := json.Unmarshal(data, &wrapper); err != nil {
		return nil, err
	}
	if err := wrapper.MarkataGo.Head.validate(); err != nil {
		return nil, err
	}

	config := wrapper.MarkataGo.toConfig()
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err == nil {
		normalizeParsedTheme(config, raw)
		copyLLMsTXTConfig(config, raw)
	}
	return config, nil
}

func normalizeParsedTheme(config *models.Config, wrapper map[string]any) {
	markata, ok := wrapper["markata-go"].(map[string]any)
	if !ok {
		return
	}
	config.Theme.MarkThemeNumericPresence(markata)
	theme, warnings := renderingcontract.NormalizeTheme(markata)
	encoded, err := json.Marshal(theme)
	if err == nil {
		if err := json.Unmarshal(encoded, &config.Theme); err != nil {
			return
		}
	}
	config.Theme.MarkThemeNumericPresence(markata)
	if len(warnings) > 0 {
		if config.Extra == nil {
			config.Extra = make(map[string]any)
		}
		config.Extra["theme_migration_warnings"] = warnings
	}
}

// tomlConfig is an internal struct for parsing TOML configuration.
type tomlConfig struct {
	Fontpack        string                    `toml:"fontpack"`
	FontpacksFile   string                    `toml:"fontpacks_file"`
	OutputDir       string                    `toml:"output_dir"`
	URL             string                    `toml:"url"`
	Title           string                    `toml:"title"`
	Description     string                    `toml:"description"`
	Author          string                    `toml:"author"`
	Language        string                    `toml:"language"`
	AuthorURL       string                    `toml:"author_url"`
	ManagingEditor  string                    `toml:"managing_editor"`
	WebMaster       string                    `toml:"webmaster"`
	Copyright       string                    `toml:"copyright"`
	License         interface{}               `toml:"license"`
	AssetsDir       string                    `toml:"assets_dir"`
	Assets          models.AssetsConfig       `toml:"assets"`
	TemplatesDir    string                    `toml:"templates_dir"`
	Nav             []tomlNavItem             `toml:"nav"`
	Footer          tomlFooterConfig          `toml:"footer"`
	Hooks           []string                  `toml:"hooks"`
	DisabledHooks   []string                  `toml:"disabled_hooks"`
	Glob            tomlGlobConfig            `toml:"glob"`
	Markdown        tomlMarkdownConfig        `toml:"markdown"`
	Feeds           []tomlFeedConfig          `toml:"feeds"`
	FeedDefaults    tomlFeedDefaults          `toml:"feed_defaults"`
	Concurrency     int                       `toml:"concurrency"`
	Theme           tomlThemeConfig           `toml:"theme"`
	ThemeCalendar   parsedThemeCalendar       `toml:"theme_calendar"`
	Head            parsedHead                `toml:"head"`
	Presets         parsedTemplatePresets     `toml:"template_presets"`
	Defaults        parsedDefaultTemplates    `toml:"default_templates"`
	ErrorPages      parsedErrorPages          `toml:"error_pages"`
	ResourceHints   parsedResourceHints       `toml:"resource_hints"`
	PostFormats     tomlPostFormatsConfig     `toml:"post_formats"`
	WellKnown       tomlWellKnownConfig       `toml:"well_known"`
	SEO             tomlSEOConfig             `toml:"seo"`
	IndieAuth       tomlIndieAuthConfig       `toml:"indieauth"`
	Webmention      tomlWebmentionConfig      `toml:"webmention"`
	Search          models.SearchConfig       `toml:"search"`
	Components      tomlComponentsConfig      `toml:"components"`
	Layout          tomlLayoutConfig          `toml:"layout"`
	Sidebar         tomlSidebarConfig         `toml:"sidebar"`
	Toc             tomlTocConfig             `toml:"toc"`
	Header          tomlHeaderLayoutConfig    `toml:"header"`
	Blogroll        tomlBlogrollConfig        `toml:"blogroll"`
	Tags            tomlTagsConfig            `toml:"tags"`
	FeedsPage       tomlFeedsPageConfig       `toml:"feeds_page"`
	Encryption      tomlEncryptionConfig      `toml:"encryption"`
	TagAggregator   tomlTagAggregatorConfig   `toml:"tag_aggregator"`
	Mentions        tomlMentionsConfig        `toml:"mentions"`
	WebSub          tomlWebSubConfig          `toml:"websub"`
	Shortcuts       tomlShortcutsConfig       `toml:"shortcuts"`
	ViewTransitions tomlViewTransitionsConfig `toml:"view_transitions"`
	Authors         tomlAuthorsConfig         `toml:"authors"`
	Garden          tomlGardenConfig          `toml:"garden"`
	Templates       tomlTemplatesConfig       `toml:"templates"`
	Images          models.ImagesConfig       `toml:"images"`
	BuilderAdmin    models.BuilderAdminConfig `toml:"builder_admin"`
	UnknownFields   map[string]any            `toml:"-"`
}

type tomlNavItem struct {
	Label    string `toml:"label"`
	URL      string `toml:"url"`
	External bool   `toml:"external"`
}

type tomlTemplatesConfig struct {
	Media tomlTemplatesMediaConfig `toml:"media"`
}

type tomlTemplatesMediaConfig struct {
	TrustedDomains []string `toml:"trusted_domains"`
}

func (t *tomlTemplatesConfig) toTemplatesConfig() models.TemplatesConfig {
	cfg := models.NewTemplatesConfig()
	if t.Media.TrustedDomains != nil {
		cfg.Media.TrustedDomains = append([]string{}, t.Media.TrustedDomains...)
	}
	return cfg
}

type tomlFooterConfig struct {
	Text          string `toml:"text"`
	ShowCopyright *bool  `toml:"show_copyright"`
}

type tomlThemeConfig struct {
	ContractVersion     int                      `toml:"contract_version"`
	Fontpack            string                   `toml:"fontpack"`
	Name                string                   `toml:"name"`
	Aesthetic           string                   `toml:"aesthetic"`
	Palette             string                   `toml:"palette"`
	PaletteLight        string                   `toml:"palette_light"`
	PaletteDark         string                   `toml:"palette_dark"`
	FallbackMode        string                   `toml:"fallback_mode"`
	TextSize            string                   `toml:"text_size"`
	Seasonal            bool                     `toml:"seasonal"`
	ShowTextSizeControl *bool                    `toml:"show_text_size_control"`
	SeedColor           string                   `toml:"seed_color"`
	Variables           map[string]string        `toml:"variables"`
	CustomCSS           string                   `toml:"custom_css"`
	Background          tomlBackgroundConfig     `toml:"background"`
	Font                tomlFontConfig           `toml:"font"`
	Switcher            tomlThemeSwitcherConfig  `toml:"switcher"`
	Texture             tomlTextureConfig        `toml:"texture"`
	HeadingTexture      tomlHeadingTextureConfig `toml:"heading_texture"`
	Motif               tomlMotifConfig          `toml:"motif"`
}

type tomlTextureConfig struct {
	Kind     string  `toml:"kind"`
	ColorMix float64 `toml:"color_mix"`
	Scale    float64 `toml:"scale"`
	Scope    string  `toml:"scope"`
}
type tomlHeadingTextureConfig struct {
	Kind     string  `toml:"kind"`
	ColorMix float64 `toml:"color_mix"`
	Scale    float64 `toml:"scale"`
}
type tomlMotifConfig struct {
	Kind      string  `toml:"kind"`
	Glyph     string  `toml:"glyph"`
	Size      string  `toml:"size"`
	Gap       string  `toml:"gap"`
	RowOffset float64 `toml:"row_offset"`
	Wobble    float64 `toml:"wobble"`
	Scatter   float64 `toml:"scatter"`
	Layer     string  `toml:"layer"`
	Color     string  `toml:"color"`
	ColorMix  float64 `toml:"color_mix"`
	URL       string  `toml:"url"`
}

type tomlThemeSwitcherConfig struct {
	Enabled    *bool    `toml:"enabled"`
	ModeToggle *bool    `toml:"mode_toggle"`
	IncludeAll *bool    `toml:"include_all"`
	Include    []string `toml:"include"`
	Exclude    []string `toml:"exclude"`
	Position   string   `toml:"position"`
}

type tomlBackgroundConfig struct {
	Enabled            *bool                   `toml:"enabled"`
	Backgrounds        []tomlBackgroundElement `toml:"backgrounds"`
	Scripts            []string                `toml:"scripts"`
	CSS                string                  `toml:"css"`
	ArticleBg          string                  `toml:"article_bg"`
	ArticleBlurEnabled *bool                   `toml:"article_blur_enabled"`
	ArticleBlur        string                  `toml:"article_blur"`
	ArticleShadow      string                  `toml:"article_shadow"`
	ArticleBorder      string                  `toml:"article_border"`
	ArticleRadius      string                  `toml:"article_radius"`
}

type tomlBackgroundElement struct {
	HTML   string `toml:"html"`
	ZIndex int    `toml:"z_index"`
}

type tomlFontConfig struct {
	Family        string   `toml:"family"`
	HeadingFamily string   `toml:"heading_family"`
	CodeFamily    string   `toml:"code_family"`
	Size          string   `toml:"size"`
	LineHeight    string   `toml:"line_height"`
	GoogleFonts   []string `toml:"google_fonts"`
	CustomURLs    []string `toml:"custom_urls"`
}

type tomlGlobConfig struct {
	Patterns     []string       `toml:"patterns"`
	UseGitignore *bool          `toml:"use_gitignore"`
	SlugMode     string         `toml:"slug_mode"`
	SlugRules    []tomlSlugRule `toml:"slug_rules"`
}

type tomlSlugRule struct {
	Prefix string `toml:"prefix"`
	Mode   string `toml:"mode"`
}

type tomlMarkdownConfig struct {
	Extensions []string        `toml:"extensions"`
	Highlight  parsedHighlight `toml:"highlight"`
}

type tomlFeedConfig struct {
	Slug            string            `toml:"slug"`
	Presentation    string            `toml:"presentation"`
	Title           string            `toml:"title"`
	Description     string            `toml:"description"`
	Robots          string            `toml:"robots"`
	Filter          string            `toml:"filter"`
	Sort            string            `toml:"sort"`
	Reverse         bool              `toml:"reverse"`
	Views           []string          `toml:"views"`
	ShortsPath      string            `toml:"shorts_path"`
	Primary         bool              `toml:"primary"`
	IncludePrivate  bool              `toml:"include_private"`
	Private         bool              `toml:"private"`
	Sidebar         *bool             `toml:"sidebar"`
	ItemsPerPage    int               `toml:"items_per_page"`
	OrphanThreshold int               `toml:"orphan_threshold"`
	Limit           int               `toml:"limit"`
	Offset          int               `toml:"offset"`
	PaginationType  string            `toml:"pagination_type"`
	ArchiveDisabled bool              `toml:"archive_disabled"`
	Formats         tomlFeedFormats   `toml:"formats"`
	Templates       tomlFeedTemplates `toml:"templates"`
}

type tomlFeedFormats struct {
	HTML       *bool `toml:"html"`
	SimpleHTML *bool `toml:"simple_html"`
	RSS        *bool `toml:"rss"`
	Atom       *bool `toml:"atom"`
	JSON       *bool `toml:"json"`
	Markdown   *bool `toml:"markdown"`
	Text       *bool `toml:"text"`
	Sitemap    *bool `toml:"sitemap"`
}

type tomlFeedTemplates struct {
	HTML       string `toml:"html"`
	SimpleHTML string `toml:"simple_html"`
	RSS        string `toml:"rss"`
	Atom       string `toml:"atom"`
	JSON       string `toml:"json"`
	Card       string `toml:"card"`
	Sitemap    string `toml:"sitemap"`
}

type tomlFeedDefaults struct {
	ItemsPerPage    int                   `toml:"items_per_page"`
	Presentation    string                `toml:"presentation"`
	OrphanThreshold int                   `toml:"orphan_threshold"`
	PaginationType  string                `toml:"pagination_type"`
	Views           []string              `toml:"views"`
	Formats         tomlFeedFormats       `toml:"formats"`
	Templates       tomlFeedTemplates     `toml:"templates"`
	Syndication     tomlSyndicationConfig `toml:"syndication"`
}

type tomlSyndicationConfig struct {
	MaxItems             int  `toml:"max_items"`
	IncludeContent       bool `toml:"include_content"`
	SiteArchiveDisabled  bool `toml:"site_archive_disabled"`
	FeedArchivesDisabled bool `toml:"feed_archives_disabled"`
}

type tomlPostFormatsConfig struct {
	HTML     *bool `toml:"html"`
	Markdown bool  `toml:"markdown"`
	Text     bool  `toml:"text"`
	ANSI     bool  `toml:"ansi"`
	OG       bool  `toml:"og"`
}

type tomlWellKnownConfig struct {
	Enabled         *bool    `toml:"enabled"`
	AutoGenerate    []string `toml:"auto_generate"`
	SSHFingerprint  string   `toml:"ssh_fingerprint"`
	KeybaseUsername string   `toml:"keybase_username"`
}

type tomlSEOConfig struct {
	TwitterHandle  string `toml:"twitter_handle"`
	DefaultImage   string `toml:"default_image"`
	LogoURL        string `toml:"logo_url"`
	AuthorImage    string `toml:"author_image"`
	OGImageService string `toml:"og_image_service"`
}

type tomlIndieAuthConfig struct {
	Enabled               bool   `toml:"enabled"`
	AuthorizationEndpoint string `toml:"authorization_endpoint"`
	TokenEndpoint         string `toml:"token_endpoint"`
	MeURL                 string `toml:"me_url"`
}

type tomlWebmentionConfig struct {
	Enabled  bool   `toml:"enabled"`
	Endpoint string `toml:"endpoint"`
}

type tomlTagsConfig struct {
	Enabled     *bool    `toml:"enabled"`
	Blacklist   []string `toml:"blacklist"`
	Private     []string `toml:"private"`
	Title       string   `toml:"title"`
	Description string   `toml:"description"`
	Template    string   `toml:"template"`
	SlugPrefix  string   `toml:"slug_prefix"`
}

type tomlFeedsPageConfig struct {
	Enabled          *bool    `toml:"enabled"`
	Title            string   `toml:"title"`
	Description      string   `toml:"description"`
	Template         string   `toml:"template"`
	SlugPrefix       string   `toml:"slug_prefix"`
	Robots           string   `toml:"robots"`
	ShowPrivateFeeds []string `toml:"show_private_feeds"`
}

func (t *tomlTagsConfig) toTagsConfig() models.TagsConfig {
	defaults := models.NewTagsConfig()

	config := models.TagsConfig{
		Enabled:     t.Enabled,
		Blacklist:   t.Blacklist,
		Private:     t.Private,
		Title:       t.Title,
		Description: t.Description,
		Template:    t.Template,
		SlugPrefix:  t.SlugPrefix,
	}

	// Apply defaults if not set
	if config.Enabled == nil {
		config.Enabled = defaults.Enabled
	}
	if config.Title == "" {
		config.Title = defaults.Title
	}
	if config.Template == "" {
		config.Template = defaults.Template
	}
	if config.SlugPrefix == "" {
		config.SlugPrefix = defaults.SlugPrefix
	}

	return config
}

func (t *tomlFeedsPageConfig) toFeedsPageConfig() models.FeedsPageConfig {
	defaults := models.NewFeedsPageConfig()

	config := models.FeedsPageConfig{
		Enabled:          t.Enabled,
		Title:            t.Title,
		Description:      t.Description,
		Template:         t.Template,
		SlugPrefix:       t.SlugPrefix,
		Robots:           t.Robots,
		ShowPrivateFeeds: append([]string{}, t.ShowPrivateFeeds...),
	}

	if config.Enabled == nil {
		config.Enabled = defaults.Enabled
	}
	if config.Title == "" {
		config.Title = defaults.Title
	}
	if config.Description == "" {
		config.Description = defaults.Description
	}
	if config.Template == "" {
		config.Template = defaults.Template
	}
	if config.SlugPrefix == "" {
		config.SlugPrefix = defaults.SlugPrefix
	}

	return config
}

type tomlTagAggregatorConfig struct {
	Enabled        *bool               `toml:"enabled"`
	Synonyms       map[string][]string `toml:"synonyms"`
	Additional     map[string][]string `toml:"additional"`
	GenerateReport bool                `toml:"generate_report"`
}

func (t *tomlTagAggregatorConfig) toTagAggregatorConfig() models.TagAggregatorConfig {
	defaults := models.NewTagAggregatorConfig()

	config := models.TagAggregatorConfig{
		Enabled:        t.Enabled,
		Synonyms:       t.Synonyms,
		Additional:     t.Additional,
		GenerateReport: t.GenerateReport,
	}

	// Apply defaults if not set
	if config.Enabled == nil {
		config.Enabled = defaults.Enabled
	}

	return config
}

type tomlWebSubConfig struct {
	Enabled *bool    `toml:"enabled"`
	Hubs    []string `toml:"hubs"`
}

func (w *tomlWebSubConfig) toWebSubConfig() models.WebSubConfig {
	defaults := models.NewWebSubConfig()
	config := models.WebSubConfig{
		Enabled: w.Enabled,
		Hubs:    w.Hubs,
	}
	if config.Enabled == nil {
		config.Enabled = defaults.Enabled
	}
	if config.Hubs == nil {
		config.Hubs = defaults.Hubs
	}
	return config
}

type tomlShortcutsConfig struct {
	Navigation map[string]string `toml:"navigation"`
}

type tomlViewTransitionsConfig struct {
	Enabled       *bool    `toml:"enabled"`
	Native        *bool    `toml:"native"`
	Debug         *bool    `toml:"debug"`
	Duration      int      `toml:"duration"`
	UpdateMeta    *bool    `toml:"update_meta"`
	ScrollToTop   *bool    `toml:"scroll_to_top"`
	SkipClasses   []string `toml:"skip_classes"`
	SkipSelectors []string `toml:"skip_selectors"`
}

func (s *tomlShortcutsConfig) toShortcutsConfig() models.ShortcutsConfig {
	config := models.ShortcutsConfig{
		Navigation: s.Navigation,
	}
	if config.Navigation == nil {
		config.Navigation = make(map[string]string)
	}
	return config
}

func (v *tomlViewTransitionsConfig) toViewTransitionsConfig() models.ViewTransitionsConfig {
	defaults := models.NewViewTransitionsConfig()
	config := models.ViewTransitionsConfig{
		Enabled:       v.Enabled,
		Native:        v.Native,
		Debug:         v.Debug,
		Duration:      v.Duration,
		UpdateMeta:    v.UpdateMeta,
		ScrollToTop:   v.ScrollToTop,
		SkipClasses:   v.SkipClasses,
		SkipSelectors: v.SkipSelectors,
	}
	if config.Enabled == nil {
		config.Enabled = defaults.Enabled
	}
	if config.Debug == nil {
		config.Debug = defaults.Debug
	}
	if config.Duration == 0 {
		config.Duration = defaults.Duration
	}
	if config.UpdateMeta == nil {
		config.UpdateMeta = defaults.UpdateMeta
	}
	if config.ScrollToTop == nil {
		config.ScrollToTop = defaults.ScrollToTop
	}
	if config.SkipClasses == nil {
		config.SkipClasses = defaults.SkipClasses
	}
	if config.SkipSelectors == nil {
		config.SkipSelectors = defaults.SkipSelectors
	}
	return config
}

type tomlAuthorsConfig struct {
	GeneratePages bool                  `toml:"generate_pages"`
	URLPattern    string                `toml:"url_pattern"`
	FeedsEnabled  bool                  `toml:"feeds_enabled"`
	Authors       map[string]tomlAuthor `toml:"authors"`
}

type tomlAuthor struct {
	Name          string            `toml:"name"`
	Bio           string            `toml:"bio"`
	Email         string            `toml:"email"`
	Avatar        string            `toml:"avatar"`
	URL           string            `toml:"url"`
	Social        map[string]string `toml:"social"`
	Guest         bool              `toml:"guest"`
	Active        bool              `toml:"active"`
	Default       bool              `toml:"default"`
	Contributions []string          `toml:"contributions"`
	Role          string            `toml:"role"`
	Contribution  string            `toml:"contribution"`
}

func (a *tomlAuthorsConfig) toAuthorsConfig() models.AuthorsConfig {
	config := models.AuthorsConfig{
		GeneratePages: a.GeneratePages,
		URLPattern:    a.URLPattern,
		FeedsEnabled:  a.FeedsEnabled,
	}
	if a.Authors != nil {
		config.Authors = make(map[string]models.Author, len(a.Authors))
		for id := range a.Authors {
			author := a.Authors[id]
			ma := models.Author{
				ID:      id,
				Name:    author.Name,
				Guest:   author.Guest,
				Active:  author.Active,
				Default: author.Default,
				Social:  author.Social,
			}
			if author.Bio != "" {
				bio := author.Bio
				ma.Bio = &bio
			}
			if author.Email != "" {
				email := author.Email
				ma.Email = &email
			}
			if author.Avatar != "" {
				avatar := author.Avatar
				ma.Avatar = &avatar
			}
			if author.URL != "" {
				url := author.URL
				ma.URL = &url
			}
			if author.Role != "" {
				role := author.Role
				ma.Role = &role
			}
			if author.Contribution != "" {
				contribution := author.Contribution
				ma.Contribution = &contribution
			}
			if len(author.Contributions) > 0 {
				ma.Contributions = author.Contributions
			}
			config.Authors[id] = ma
		}
	}
	return config
}

type tomlGardenConfig struct {
	Enabled      *bool    `toml:"enabled"`
	Path         string   `toml:"path"`
	ExportJSON   *bool    `toml:"export_json"`
	RenderPage   *bool    `toml:"render_page"`
	IncludeTags  *bool    `toml:"include_tags"`
	IncludePosts *bool    `toml:"include_posts"`
	MaxNodes     int      `toml:"max_nodes"`
	ExcludeTags  []string `toml:"exclude_tags"`
	Template     string   `toml:"template"`
	Title        string   `toml:"title"`
	Description  string   `toml:"description"`
}

func (g *tomlGardenConfig) toGardenConfig() models.GardenConfig {
	defaults := models.NewGardenConfig()

	config := models.GardenConfig{
		Enabled:      g.Enabled,
		Path:         g.Path,
		ExportJSON:   g.ExportJSON,
		RenderPage:   g.RenderPage,
		IncludeTags:  g.IncludeTags,
		IncludePosts: g.IncludePosts,
		MaxNodes:     g.MaxNodes,
		ExcludeTags:  g.ExcludeTags,
		Template:     g.Template,
		Title:        g.Title,
		Description:  g.Description,
	}

	// Apply defaults for unset fields
	if config.Enabled == nil {
		config.Enabled = defaults.Enabled
	}
	if config.ExportJSON == nil {
		config.ExportJSON = defaults.ExportJSON
	}
	if config.RenderPage == nil {
		config.RenderPage = defaults.RenderPage
	}
	if config.IncludeTags == nil {
		config.IncludeTags = defaults.IncludeTags
	}
	if config.IncludePosts == nil {
		config.IncludePosts = defaults.IncludePosts
	}
	if config.Path == "" {
		config.Path = defaults.Path
	}
	if config.Template == "" {
		config.Template = defaults.Template
	}
	if config.Title == "" {
		config.Title = defaults.Title
	}

	return config
}

func (s *tomlSEOConfig) toSEOConfig() models.SEOConfig {
	return models.SEOConfig{
		TwitterHandle:  s.TwitterHandle,
		DefaultImage:   s.DefaultImage,
		LogoURL:        s.LogoURL,
		AuthorImage:    s.AuthorImage,
		OGImageService: s.OGImageService,
	}
}

func (i *tomlIndieAuthConfig) toIndieAuthConfig() models.IndieAuthConfig {
	return models.IndieAuthConfig{
		Enabled:               i.Enabled,
		AuthorizationEndpoint: i.AuthorizationEndpoint,
		TokenEndpoint:         i.TokenEndpoint,
		MeURL:                 i.MeURL,
	}
}

func (w *tomlWebmentionConfig) toWebmentionConfig() models.WebmentionConfig {
	return models.WebmentionConfig{
		Enabled:  w.Enabled,
		Endpoint: w.Endpoint,
	}
}

type tomlComponentsConfig struct {
	Nav            tomlNavComponentConfig      `toml:"nav"`
	Footer         tomlFooterComponentConfig   `toml:"footer"`
	DocSidebar     tomlDocSidebarConfig        `toml:"doc_sidebar"`
	FeedSidebar    tomlFeedSidebarConfig       `toml:"feed_sidebar"`
	ContentSidebar tomlContentSidebarConfig    `toml:"content_sidebar"`
	CardRouter     tomlCardRouterConfig        `toml:"card_router"`
	Share          tomlShareComponentConfig    `toml:"share"`
	PostCopy       tomlPostCopyComponentConfig `toml:"post_copy"`
	PostConn       tomlPostConnectionsConfig   `toml:"post_connections"`
}

type tomlPostCopyComponentConfig struct {
	Enabled *bool `toml:"enabled"`
}

type tomlPostConnectionsConfig struct {
	Enabled       *bool    `toml:"enabled"`
	Display       []string `toml:"display"`
	MinLinks      int      `toml:"min_links"`
	MaxLinks      int      `toml:"max_links"`
	GraphMinLinks int      `toml:"graph_min_links"`
	GraphMaxLinks int      `toml:"graph_max_links"`
	ListMinLinks  int      `toml:"list_min_links"`
	ListMaxLinks  int      `toml:"list_max_links"`
	InlinksLimit  int      `toml:"inlinks_limit"`
	OutlinksLimit int      `toml:"outlinks_limit"`
}

type tomlNavComponentConfig struct {
	Enabled  *bool         `toml:"enabled"`
	Position string        `toml:"position"`
	Style    string        `toml:"style"`
	Items    []tomlNavItem `toml:"items"`
}

type tomlFooterComponentConfig struct {
	Enabled       *bool         `toml:"enabled"`
	Text          string        `toml:"text"`
	ShowCopyright *bool         `toml:"show_copyright"`
	Links         []tomlNavItem `toml:"links"`
}

type tomlDocSidebarConfig struct {
	Enabled  *bool  `toml:"enabled"`
	Position string `toml:"position"`
	Width    string `toml:"width"`
	MinDepth int    `toml:"min_depth"`
	MaxDepth int    `toml:"max_depth"`
}

type tomlFeedSidebarConfig struct {
	Enabled  *bool    `toml:"enabled"`
	Position string   `toml:"position"`
	Width    string   `toml:"width"`
	Title    string   `toml:"title"`
	Feeds    []string `toml:"feeds"`
	MaxPosts int      `toml:"max_posts"`
}

type tomlContentSidebarConfig struct {
	Enabled  *bool  `toml:"enabled"`
	Position string `toml:"position"`
	Width    string `toml:"width"`
	Slug     string `toml:"slug"`
}

type tomlCardRouterConfig struct {
	Mappings map[string]string `toml:"mappings"`
}

type tomlShareComponentConfig struct {
	Enabled   *bool                             `toml:"enabled"`
	Platforms []string                          `toml:"platforms"`
	Position  string                            `toml:"position"`
	Title     string                            `toml:"title"`
	Custom    map[string]tomlSharePlatformEntry `toml:"custom"`
}

type tomlSharePlatformEntry struct {
	Name string `toml:"name"`
	Icon string `toml:"icon"`
	URL  string `toml:"url"`
}

// Layout-related TOML structs

type tomlLayoutConfig struct {
	Name     string                  `toml:"name"`
	Paths    map[string]string       `toml:"paths"`
	Feeds    map[string]string       `toml:"feeds"`
	Docs     tomlDocsLayoutConfig    `toml:"docs"`
	Blog     tomlBlogLayoutConfig    `toml:"blog"`
	Landing  tomlLandingLayoutConfig `toml:"landing"`
	Bare     tomlBareLayoutConfig    `toml:"bare"`
	Defaults tomlLayoutDefaults      `toml:"defaults"`
}

type tomlLayoutDefaults struct {
	ContentMaxWidth string `toml:"content_max_width"`
	HeaderSticky    *bool  `toml:"header_sticky"`
	FooterSticky    *bool  `toml:"footer_sticky"`
}

type tomlDocsLayoutConfig struct {
	SidebarPosition    string `toml:"sidebar_position"`
	SidebarWidth       string `toml:"sidebar_width"`
	SidebarCollapsible *bool  `toml:"sidebar_collapsible"`
	SidebarDefaultOpen *bool  `toml:"sidebar_default_open"`
	TocPosition        string `toml:"toc_position"`
	TocWidth           string `toml:"toc_width"`
	TocCollapsible     *bool  `toml:"toc_collapsible"`
	TocDefaultOpen     *bool  `toml:"toc_default_open"`
	ContentMaxWidth    string `toml:"content_max_width"`
	HeaderStyle        string `toml:"header_style"`
	FooterStyle        string `toml:"footer_style"`
}

type tomlBlogLayoutConfig struct {
	ContentMaxWidth string `toml:"content_max_width"`
	ShowToc         *bool  `toml:"show_toc"`
	TocPosition     string `toml:"toc_position"`
	TocWidth        string `toml:"toc_width"`
	HeaderStyle     string `toml:"header_style"`
	FooterStyle     string `toml:"footer_style"`
	ShowAuthor      *bool  `toml:"show_author"`
	ShowDate        *bool  `toml:"show_date"`
	ShowTags        *bool  `toml:"show_tags"`
	ShowReadingTime *bool  `toml:"show_reading_time"`
	ShowPrevNext    *bool  `toml:"show_prev_next"`
}

type tomlLandingLayoutConfig struct {
	ContentMaxWidth string `toml:"content_max_width"`
	HeaderStyle     string `toml:"header_style"`
	HeaderSticky    *bool  `toml:"header_sticky"`
	FooterStyle     string `toml:"footer_style"`
	HeroEnabled     *bool  `toml:"hero_enabled"`
}

type tomlBareLayoutConfig struct {
	ContentMaxWidth string `toml:"content_max_width"`
}

func (l *tomlLayoutConfig) toLayoutConfig() models.LayoutConfig {
	return models.LayoutConfig{
		Name:  l.Name,
		Paths: l.Paths,
		Feeds: l.Feeds,
		Docs: models.DocsLayoutConfig{
			SidebarPosition:    l.Docs.SidebarPosition,
			SidebarWidth:       l.Docs.SidebarWidth,
			SidebarCollapsible: l.Docs.SidebarCollapsible,
			SidebarDefaultOpen: l.Docs.SidebarDefaultOpen,
			TocPosition:        l.Docs.TocPosition,
			TocWidth:           l.Docs.TocWidth,
			TocCollapsible:     l.Docs.TocCollapsible,
			TocDefaultOpen:     l.Docs.TocDefaultOpen,
			ContentMaxWidth:    l.Docs.ContentMaxWidth,
			HeaderStyle:        l.Docs.HeaderStyle,
			FooterStyle:        l.Docs.FooterStyle,
		},
		Blog: models.BlogLayoutConfig{
			ContentMaxWidth: l.Blog.ContentMaxWidth,
			ShowToc:         l.Blog.ShowToc,
			TocPosition:     l.Blog.TocPosition,
			TocWidth:        l.Blog.TocWidth,
			HeaderStyle:     l.Blog.HeaderStyle,
			FooterStyle:     l.Blog.FooterStyle,
			ShowAuthor:      l.Blog.ShowAuthor,
			ShowDate:        l.Blog.ShowDate,
			ShowTags:        l.Blog.ShowTags,
			ShowReadingTime: l.Blog.ShowReadingTime,
			ShowPrevNext:    l.Blog.ShowPrevNext,
		},
		Landing: models.LandingLayoutConfig{
			ContentMaxWidth: l.Landing.ContentMaxWidth,
			HeaderStyle:     l.Landing.HeaderStyle,
			HeaderSticky:    l.Landing.HeaderSticky,
			FooterStyle:     l.Landing.FooterStyle,
			HeroEnabled:     l.Landing.HeroEnabled,
		},
		Bare: models.BareLayoutConfig{
			ContentMaxWidth: l.Bare.ContentMaxWidth,
		},
		Defaults: models.LayoutDefaults{
			ContentMaxWidth: l.Defaults.ContentMaxWidth,
			HeaderSticky:    l.Defaults.HeaderSticky,
			FooterSticky:    l.Defaults.FooterSticky,
		},
	}
}

// Sidebar-related TOML structs

type tomlSidebarConfig struct {
	Enabled      *bool                             `toml:"enabled"`
	Position     string                            `toml:"position"`
	Width        string                            `toml:"width"`
	Collapsible  *bool                             `toml:"collapsible"`
	DefaultOpen  *bool                             `toml:"default_open"`
	Nav          []tomlSidebarNavItem              `toml:"nav"`
	Title        string                            `toml:"title"`
	Paths        map[string]*tomlPathSidebarConfig `toml:"paths"`
	MultiFeed    *bool                             `toml:"multi_feed"`
	Feeds        []string                          `toml:"feeds"`
	FeedSections []tomlMultiFeedSection            `toml:"feed_sections"`
	AutoGenerate *tomlSidebarAutoGenerate          `toml:"auto_generate"`
}

type tomlSidebarNavItem struct {
	Title    string               `toml:"title"`
	Href     string               `toml:"href"`
	Children []tomlSidebarNavItem `toml:"children"`
}

type tomlPathSidebarConfig struct {
	Title        string                   `toml:"title"`
	AutoGenerate *tomlSidebarAutoGenerate `toml:"auto_generate"`
	Items        []tomlSidebarNavItem     `toml:"items"`
	Feed         string                   `toml:"feed"`
	Position     string                   `toml:"position"`
	Collapsible  *bool                    `toml:"collapsible"`
}

type tomlSidebarAutoGenerate struct {
	Directory string   `toml:"directory"`
	OrderBy   string   `toml:"order_by"`
	Reverse   *bool    `toml:"reverse"`
	MaxDepth  int      `toml:"max_depth"`
	Exclude   []string `toml:"exclude"`
}

type tomlMultiFeedSection struct {
	Feed      string `toml:"feed"`
	Title     string `toml:"title"`
	Collapsed *bool  `toml:"collapsed"`
	MaxItems  int    `toml:"max_items"`
}

func convertTomlSidebarNavItems(items []tomlSidebarNavItem) []models.SidebarNavItem {
	result := make([]models.SidebarNavItem, len(items))
	for i, item := range items {
		result[i] = models.SidebarNavItem{
			Title:    item.Title,
			Href:     item.Href,
			Children: convertTomlSidebarNavItems(item.Children),
		}
	}
	return result
}

func (s *tomlSidebarConfig) toSidebarConfig() models.SidebarConfig {
	config := models.SidebarConfig{
		Enabled:     s.Enabled,
		Position:    s.Position,
		Width:       s.Width,
		Collapsible: s.Collapsible,
		DefaultOpen: s.DefaultOpen,
		Nav:         convertTomlSidebarNavItems(s.Nav),
		Title:       s.Title,
		MultiFeed:   s.MultiFeed,
		Feeds:       s.Feeds,
	}

	// Convert paths
	if len(s.Paths) > 0 {
		config.Paths = make(map[string]*models.PathSidebarConfig)
		for path, pathConfig := range s.Paths {
			var autoGen *models.SidebarAutoGenerate
			if pathConfig.AutoGenerate != nil {
				autoGen = &models.SidebarAutoGenerate{
					Directory: pathConfig.AutoGenerate.Directory,
					OrderBy:   pathConfig.AutoGenerate.OrderBy,
					Reverse:   pathConfig.AutoGenerate.Reverse,
					MaxDepth:  pathConfig.AutoGenerate.MaxDepth,
					Exclude:   pathConfig.AutoGenerate.Exclude,
				}
			}
			config.Paths[path] = &models.PathSidebarConfig{
				Title:        pathConfig.Title,
				AutoGenerate: autoGen,
				Items:        convertTomlSidebarNavItems(pathConfig.Items),
				Feed:         pathConfig.Feed,
				Position:     pathConfig.Position,
				Collapsible:  pathConfig.Collapsible,
			}
		}
	}

	// Convert feed sections
	if len(s.FeedSections) > 0 {
		config.FeedSections = make([]models.MultiFeedSection, len(s.FeedSections))
		for i, section := range s.FeedSections {
			config.FeedSections[i] = models.MultiFeedSection{
				Feed:      section.Feed,
				Title:     section.Title,
				Collapsed: section.Collapsed,
				MaxItems:  section.MaxItems,
			}
		}
	}

	// Convert auto-generate
	if s.AutoGenerate != nil {
		config.AutoGenerate = &models.SidebarAutoGenerate{
			Directory: s.AutoGenerate.Directory,
			OrderBy:   s.AutoGenerate.OrderBy,
			Reverse:   s.AutoGenerate.Reverse,
			MaxDepth:  s.AutoGenerate.MaxDepth,
			Exclude:   s.AutoGenerate.Exclude,
		}
	}

	return config
}

// TOC-related TOML structs

type tomlTocConfig struct {
	Enabled     *bool  `toml:"enabled"`
	Position    string `toml:"position"`
	Width       string `toml:"width"`
	MinDepth    int    `toml:"min_depth"`
	MaxDepth    int    `toml:"max_depth"`
	Title       string `toml:"title"`
	Collapsible *bool  `toml:"collapsible"`
	DefaultOpen *bool  `toml:"default_open"`
	ScrollSpy   *bool  `toml:"scroll_spy"`
}

func (t *tomlTocConfig) toTocConfig() models.TocConfig {
	return models.TocConfig{
		Enabled:     t.Enabled,
		Position:    t.Position,
		Width:       t.Width,
		MinDepth:    t.MinDepth,
		MaxDepth:    t.MaxDepth,
		Title:       t.Title,
		Collapsible: t.Collapsible,
		DefaultOpen: t.DefaultOpen,
		ScrollSpy:   t.ScrollSpy,
	}
}

// Header layout TOML structs

type tomlHeaderLayoutConfig struct {
	Style           string `toml:"style"`
	Sticky          *bool  `toml:"sticky"`
	ShowLogo        *bool  `toml:"show_logo"`
	ShowTitle       *bool  `toml:"show_title"`
	ShowNav         *bool  `toml:"show_nav"`
	ShowSearch      *bool  `toml:"show_search"`
	ShowThemeToggle *bool  `toml:"show_theme_toggle"`
}

func (h *tomlHeaderLayoutConfig) toHeaderLayoutConfig() models.HeaderLayoutConfig {
	return models.HeaderLayoutConfig{
		Style:           h.Style,
		Sticky:          h.Sticky,
		ShowLogo:        h.ShowLogo,
		ShowTitle:       h.ShowTitle,
		ShowNav:         h.ShowNav,
		ShowSearch:      h.ShowSearch,
		ShowThemeToggle: h.ShowThemeToggle,
	}
}

// Blogroll-related TOML structs

type tomlBlogrollConfig struct {
	Enabled              bool                     `toml:"enabled"`
	BlogrollSlug         string                   `toml:"blogroll_slug"`
	ReaderSlug           string                   `toml:"reader_slug"`
	CacheDir             string                   `toml:"cache_dir"`
	CacheDuration        string                   `toml:"cache_duration"`
	RefreshOnBuild       *bool                    `toml:"refresh_on_build"`
	Timeout              int                      `toml:"timeout"`
	ConcurrentRequests   int                      `toml:"concurrent_requests"`
	MaxEntriesPerFeed    int                      `toml:"max_entries_per_feed"`
	FallbackImageService string                   `toml:"fallback_image_service"`
	Feeds                []tomlExternalFeedConfig `toml:"feeds"`
	Templates            tomlBlogrollTemplates    `toml:"templates"`
}

type tomlExternalFeedConfig struct {
	URL           string   `toml:"url"`
	Title         string   `toml:"title"`
	Description   string   `toml:"description"`
	Category      string   `toml:"category"`
	Type          string   `toml:"type"`
	Tags          []string `toml:"tags"`
	Active        *bool    `toml:"active"`
	SiteURL       string   `toml:"site_url"`
	ImageURL      string   `toml:"image_url"`
	Handle        string   `toml:"handle"`
	Aliases       []string `toml:"aliases,omitempty"`
	MaxEntries    *int     `toml:"max_entries,omitempty"`
	Primary       *bool    `toml:"primary,omitempty"`
	PrimaryPerson string   `toml:"primary_person"`
}

type tomlBlogrollTemplates struct {
	Blogroll string `toml:"blogroll"`
	Reader   string `toml:"reader"`
}

type tomlEncryptionConfig struct {
	Enabled               *bool             `toml:"enabled"`
	DefaultKey            string            `toml:"default_key"`
	DecryptionHint        string            `toml:"decryption_hint"`
	PrivateTags           map[string]string `toml:"private_tags"`
	EnforceStrength       *bool             `toml:"enforce_strength"`
	MinEstimatedCrackTime string            `toml:"min_estimated_crack_time"`
	MinPasswordLength     int               `toml:"min_password_length"`
}

func (e *tomlEncryptionConfig) toEncryptionConfig() models.EncryptionConfig {
	defaults := models.NewEncryptionConfig()

	config := models.EncryptionConfig{
		DefaultKey:            e.DefaultKey,
		DecryptionHint:        e.DecryptionHint,
		PrivateTags:           e.PrivateTags,
		MinEstimatedCrackTime: e.MinEstimatedCrackTime,
		MinPasswordLength:     e.MinPasswordLength,
	}

	// Apply defaults for unset values
	if e.Enabled != nil {
		config.Enabled = *e.Enabled
	} else {
		config.Enabled = defaults.Enabled
	}
	if config.DefaultKey == "" {
		config.DefaultKey = defaults.DefaultKey
	}
	if e.EnforceStrength != nil {
		config.EnforceStrength = *e.EnforceStrength
	} else {
		config.EnforceStrength = defaults.EnforceStrength
	}
	if config.MinEstimatedCrackTime == "" {
		config.MinEstimatedCrackTime = defaults.MinEstimatedCrackTime
	}
	if config.MinPasswordLength == 0 {
		config.MinPasswordLength = defaults.MinPasswordLength // pragma: allowlist secret
	}

	return config
}

// Mentions-related TOML structs

type tomlMentionsConfig struct {
	Enabled            *bool                   `toml:"enabled"`
	CSSClass           string                  `toml:"css_class"`
	FromPosts          []tomlMentionPostSource `toml:"from_posts"`
	CacheDir           string                  `toml:"cache_dir"`
	CacheDuration      string                  `toml:"cache_duration"`
	Timeout            int                     `toml:"timeout"`
	ConcurrentRequests int                     `toml:"concurrent_requests"`
}

type tomlMentionPostSource struct {
	Filter       string `toml:"filter"`
	HandleField  string `toml:"handle_field"`
	AliasesField string `toml:"aliases_field"`
	AvatarField  string `toml:"avatar_field"`
}

func (m *tomlMentionsConfig) toMentionsConfig() models.MentionsConfig {
	defaults := models.NewMentionsConfig()

	config := models.MentionsConfig{
		Enabled:            m.Enabled,
		CSSClass:           m.CSSClass,
		CacheDir:           m.CacheDir,
		CacheDuration:      m.CacheDuration,
		Timeout:            m.Timeout,
		ConcurrentRequests: m.ConcurrentRequests,
	}

	// Apply defaults for unset values
	if config.Enabled == nil {
		config.Enabled = defaults.Enabled
	}
	if config.CSSClass == "" {
		config.CSSClass = defaults.CSSClass
	}
	if config.CacheDir == "" {
		config.CacheDir = defaults.CacheDir
	}
	if config.CacheDuration == "" {
		config.CacheDuration = defaults.CacheDuration
	}
	if config.Timeout == 0 {
		config.Timeout = defaults.Timeout
	}
	if config.ConcurrentRequests == 0 {
		config.ConcurrentRequests = defaults.ConcurrentRequests
	}

	// Convert from_posts sources
	for _, src := range m.FromPosts {
		aliasesField := src.AliasesField
		if aliasesField == "" {
			aliasesField = defaultAliasesField
		}
		config.FromPosts = append(config.FromPosts, models.MentionPostSource{
			Filter:       src.Filter,
			HandleField:  src.HandleField,
			AliasesField: aliasesField,
			AvatarField:  src.AvatarField,
		})
	}

	// Fall back to defaults when user has mentions config but no from_posts
	if len(config.FromPosts) == 0 {
		config.FromPosts = defaults.FromPosts
	}

	return config
}

const defaultAliasesField = "aliases"

func (b *tomlBlogrollConfig) toBlogrollConfig() models.BlogrollConfig {
	// Get default values
	defaults := models.NewBlogrollConfig()

	// Apply defaults for template names if not specified
	blogrollTemplate := b.Templates.Blogroll
	if blogrollTemplate == "" {
		blogrollTemplate = defaults.Templates.Blogroll
	}
	readerTemplate := b.Templates.Reader
	if readerTemplate == "" {
		readerTemplate = defaults.Templates.Reader
	}

	config := models.BlogrollConfig{
		Enabled:              b.Enabled,
		BlogrollSlug:         b.BlogrollSlug,
		ReaderSlug:           b.ReaderSlug,
		CacheDir:             b.CacheDir,
		CacheDuration:        b.CacheDuration,
		RefreshOnBuild:       b.RefreshOnBuild,
		Timeout:              b.Timeout,
		ConcurrentRequests:   b.ConcurrentRequests,
		MaxEntriesPerFeed:    b.MaxEntriesPerFeed,
		FallbackImageService: b.FallbackImageService,
		Templates: models.BlogrollTemplates{
			Blogroll: blogrollTemplate,
			Reader:   readerTemplate,
		},
	}

	for i := range b.Feeds {
		fc := &b.Feeds[i]
		config.Feeds = append(config.Feeds, models.ExternalFeedConfig{
			URL:           fc.URL,
			Title:         fc.Title,
			Description:   fc.Description,
			Category:      fc.Category,
			Type:          models.ReaderFeedType(fc.Type),
			Tags:          fc.Tags,
			Active:        fc.Active,
			SiteURL:       fc.SiteURL,
			ImageURL:      fc.ImageURL,
			Handle:        fc.Handle,
			Aliases:       fc.Aliases,
			MaxEntries:    fc.MaxEntries,
			Primary:       fc.Primary,
			PrimaryPerson: fc.PrimaryPerson,
		})
	}

	return config
}

func (c *tomlComponentsConfig) toComponentsConfig() models.ComponentsConfig {
	config := models.ComponentsConfig{
		Nav: models.NavComponentConfig{
			Enabled:  c.Nav.Enabled,
			Position: c.Nav.Position,
			Style:    c.Nav.Style,
		},
		Footer: models.FooterComponentConfig{
			Enabled:       c.Footer.Enabled,
			Text:          c.Footer.Text,
			ShowCopyright: c.Footer.ShowCopyright,
		},
		DocSidebar: models.DocSidebarConfig{
			Enabled:  c.DocSidebar.Enabled,
			Position: c.DocSidebar.Position,
			Width:    c.DocSidebar.Width,
			MinDepth: c.DocSidebar.MinDepth,
			MaxDepth: c.DocSidebar.MaxDepth,
		},
		FeedSidebar: models.FeedSidebarConfig{
			Enabled:  c.FeedSidebar.Enabled,
			Position: c.FeedSidebar.Position,
			Width:    c.FeedSidebar.Width,
			Title:    c.FeedSidebar.Title,
			Feeds:    c.FeedSidebar.Feeds,
			MaxPosts: c.FeedSidebar.MaxPosts,
		},
		ContentSidebar: models.ContentSidebarConfig{
			Enabled:  c.ContentSidebar.Enabled,
			Position: c.ContentSidebar.Position,
			Width:    c.ContentSidebar.Width,
			Slug:     c.ContentSidebar.Slug,
		},
		CardRouter: models.CardRouterConfig{
			Mappings: c.CardRouter.Mappings,
		},
		Share: models.ShareComponentConfig{
			Enabled:   c.Share.Enabled,
			Platforms: append([]string{}, c.Share.Platforms...),
			Position:  c.Share.Position,
			Title:     c.Share.Title,
			Custom:    map[string]models.SharePlatformConfig{},
		},
		PostCopy: models.PostCopyComponentConfig{Enabled: c.PostCopy.Enabled},
		PostConnections: models.PostConnectionsComponentConfig{
			Enabled:       c.PostConn.Enabled,
			Display:       append([]string{}, c.PostConn.Display...),
			MinLinks:      c.PostConn.MinLinks,
			MaxLinks:      c.PostConn.MaxLinks,
			GraphMinLinks: c.PostConn.GraphMinLinks,
			GraphMaxLinks: c.PostConn.GraphMaxLinks,
			ListMinLinks:  c.PostConn.ListMinLinks,
			ListMaxLinks:  c.PostConn.ListMaxLinks,
			InlinksLimit:  c.PostConn.InlinksLimit,
			OutlinksLimit: c.PostConn.OutlinksLimit,
		},
	}

	if len(c.Share.Platforms) == 0 {
		config.Share.Platforms = nil
	}

	if len(c.Share.Platforms) == 0 {
		config.Share.Platforms = nil
	}

	if len(c.Share.Platforms) == 0 {
		config.Share.Platforms = nil
	}

	if len(c.Share.Platforms) == 0 {
		config.Share.Platforms = nil
	}

	if len(c.Share.Platforms) == 0 {
		config.Share.Platforms = nil
	}

	for key, custom := range c.Share.Custom {
		config.Share.Custom[key] = models.SharePlatformConfig{
			Name: custom.Name,
			Icon: custom.Icon,
			URL:  custom.URL,
		}
	}
	if len(c.Share.Custom) == 0 {
		config.Share.Custom = nil
	}
	if len(c.Share.Custom) == 0 {
		config.Share.Custom = nil
	}
	if len(c.Share.Custom) == 0 {
		config.Share.Custom = nil
	}
	if len(c.Share.Custom) == 0 {
		config.Share.Custom = nil
	}
	if len(c.Share.Custom) == 0 {
		config.Share.Custom = nil
	}

	// Convert nav items
	for _, item := range c.Nav.Items {
		config.Nav.Items = append(config.Nav.Items, models.NavItem{
			Label:    item.Label,
			URL:      item.URL,
			External: item.External,
		})
	}

	// Convert footer links
	for _, link := range c.Footer.Links {
		config.Footer.Links = append(config.Footer.Links, models.NavItem{
			Label:    link.Label,
			URL:      link.URL,
			External: link.External,
		})
	}

	return config
}

func (p *tomlPostFormatsConfig) toPostFormatsConfig() models.PostFormatsConfig {
	return models.PostFormatsConfig{
		HTML:     p.HTML,
		Markdown: p.Markdown,
		Text:     p.Text,
		ANSI:     p.ANSI,
		OG:       p.OG,
	}
}

func (w *tomlWellKnownConfig) toWellKnownConfig() models.WellKnownConfig {
	defaults := models.NewWellKnownConfig()
	config := models.WellKnownConfig{
		Enabled:         w.Enabled,
		AutoGenerate:    w.AutoGenerate,
		SSHFingerprint:  w.SSHFingerprint,
		KeybaseUsername: w.KeybaseUsername,
	}

	if config.Enabled == nil {
		config.Enabled = defaults.Enabled
	}
	if config.AutoGenerate == nil {
		config.AutoGenerate = defaults.AutoGenerate
	}

	return config
}

// configSource interface implementation for tomlConfig.
func (c *tomlConfig) getBaseConfig() baseConfigData {
	return baseConfigData{
		Fontpack:       c.Fontpack,
		FontpacksFile:  c.FontpacksFile,
		OutputDir:      c.OutputDir,
		URL:            c.URL,
		Title:          c.Title,
		Description:    c.Description,
		Author:         c.Author,
		Language:       c.Language,
		AuthorURL:      c.AuthorURL,
		ManagingEditor: c.ManagingEditor,
		WebMaster:      c.WebMaster,
		Copyright:      c.Copyright,
		License:        c.License,
		AssetsDir:      c.AssetsDir,
		Assets:         c.Assets,
		TemplatesDir:   c.TemplatesDir,
		Hooks:          c.Hooks,
		DisabledHooks:  c.DisabledHooks,
		GlobPatterns:   c.Glob.Patterns,
		UseGitignore:   c.Glob.UseGitignore,
		SlugMode:       c.Glob.SlugMode,
		SlugRules:      c.Glob.toSlugRules(),
		Extensions:     c.Markdown.Extensions,
		Concurrency:    c.Concurrency,
		Theme:          c.Theme.toThemeConfig(),
		ThemeCalendar:  c.ThemeCalendar,
		Head:           c.Head.toHeadConfig(),
		Presets:        c.Presets,
		Defaults:       c.Defaults,
		ErrorPages:     c.ErrorPages,
		ResourceHints:  c.ResourceHints,
		Highlight:      c.Markdown.Highlight,
		Footer:         c.Footer.toFooterConfig(),
	}
}

func (c *tomlConfig) getNavItems() []navItemData {
	items := make([]navItemData, len(c.Nav))
	for i, nav := range c.Nav {
		items[i] = navItemData(nav)
	}
	return items
}

func (c *tomlConfig) getFeeds() []feedConfigConverter {
	feeds := make([]feedConfigConverter, len(c.Feeds))
	for i := range c.Feeds {
		feeds[i] = &c.Feeds[i]
	}
	return feeds
}

func (c *tomlConfig) getFeedDefaults() feedDefaultsConverter       { return &c.FeedDefaults }
func (c *tomlConfig) getPostFormats() postFormatsConverter         { return &c.PostFormats }
func (c *tomlConfig) getWellKnown() wellKnownConverter             { return &c.WellKnown }
func (c *tomlConfig) getSEO() seoConverter                         { return &c.SEO }
func (c *tomlConfig) getIndieAuth() indieAuthConverter             { return &c.IndieAuth }
func (c *tomlConfig) getWebmention() webmentionConverter           { return &c.Webmention }
func (c *tomlConfig) getSearch() models.SearchConfig               { return c.Search }
func (c *tomlConfig) getComponents() componentsConverter           { return &c.Components }
func (c *tomlConfig) getLayout() layoutConverter                   { return &c.Layout }
func (c *tomlConfig) getSidebar() sidebarConverter                 { return &c.Sidebar }
func (c *tomlConfig) getToc() tocConverter                         { return &c.Toc }
func (c *tomlConfig) getHeader() headerConverter                   { return &c.Header }
func (c *tomlConfig) getBlogroll() blogrollConverter               { return &c.Blogroll }
func (c *tomlConfig) getTags() tagsConverter                       { return &c.Tags }
func (c *tomlConfig) getFeedsPage() feedsPageConverter             { return &c.FeedsPage }
func (c *tomlConfig) getEncryption() encryptionConverter           { return &c.Encryption }
func (c *tomlConfig) getTagAggregator() tagAggregatorConverter     { return &c.TagAggregator }
func (c *tomlConfig) getWebSub() webSubConverter                   { return &c.WebSub }
func (c *tomlConfig) getMentions() mentionsConverter               { return &c.Mentions }
func (c *tomlConfig) getShortcuts() shortcutsConverter             { return &c.Shortcuts }
func (c *tomlConfig) getViewTransitions() viewTransitionsConverter { return &c.ViewTransitions }
func (c *tomlConfig) getAuthors() authorsConverter                 { return &c.Authors }
func (c *tomlConfig) getGarden() gardenConverter                   { return &c.Garden }
func (c *tomlConfig) getTemplates() templatesConverter             { return &c.Templates }
func (c *tomlConfig) getBuilderAdmin() models.BuilderAdminConfig   { return c.BuilderAdmin }
func (c *tomlConfig) getImages() models.ImagesConfig               { return normalizeImagesConfig(c.Images) }

func (c *tomlConfig) toConfig() *models.Config {
	return buildConfig(c)
}

func (f *tomlFooterConfig) toFooterConfig() models.FooterConfig {
	return models.FooterConfig{
		Text:          f.Text,
		ShowCopyright: f.ShowCopyright,
	}
}

func (t *tomlThemeConfig) toThemeConfig() models.ThemeConfig {
	variables := t.Variables
	if variables == nil {
		variables = make(map[string]string)
	}
	return models.ThemeConfig{
		ContractVersion:     t.ContractVersion,
		Fontpack:            t.Fontpack,
		Name:                t.Name,
		Aesthetic:           t.Aesthetic,
		Palette:             t.Palette,
		PaletteLight:        t.PaletteLight,
		PaletteDark:         t.PaletteDark,
		FallbackMode:        t.FallbackMode,
		TextSize:            t.TextSize,
		Seasonal:            t.Seasonal,
		ShowTextSizeControl: t.ShowTextSizeControl,
		SeedColor:           t.SeedColor,
		Variables:           variables,
		CustomCSS:           t.CustomCSS,
		Background:          t.Background.toBackgroundConfig(),
		Font:                t.Font.toFontConfig(),
		Switcher:            t.Switcher.toThemeSwitcherConfig(),
		Texture:             models.ThemeTextureConfig{Kind: t.Texture.Kind, ColorMix: t.Texture.ColorMix, Scale: t.Texture.Scale, Scope: t.Texture.Scope},
		HeadingTexture:      models.ThemeHeadingTextureConfig{Kind: t.HeadingTexture.Kind, ColorMix: t.HeadingTexture.ColorMix, Scale: t.HeadingTexture.Scale},
		Motif:               models.ThemeMotifConfig{Kind: t.Motif.Kind, Glyph: t.Motif.Glyph, Size: t.Motif.Size, Gap: t.Motif.Gap, RowOffset: t.Motif.RowOffset, Wobble: t.Motif.Wobble, Scatter: t.Motif.Scatter, Layer: t.Motif.Layer, Color: t.Motif.Color, ColorMix: t.Motif.ColorMix, URL: t.Motif.URL},
	}
}

func (s *tomlThemeSwitcherConfig) toThemeSwitcherConfig() models.ThemeSwitcherConfig {
	return models.ThemeSwitcherConfig{
		Enabled:    s.Enabled,
		ModeToggle: s.ModeToggle,
		IncludeAll: s.IncludeAll,
		Include:    s.Include,
		Exclude:    s.Exclude,
		Position:   s.Position,
	}
}

func (b *tomlBackgroundConfig) toBackgroundConfig() models.BackgroundConfig {
	backgrounds := make([]models.BackgroundElement, len(b.Backgrounds))
	for i, bg := range b.Backgrounds {
		backgrounds[i] = models.BackgroundElement{
			HTML:   bg.HTML,
			ZIndex: bg.ZIndex,
		}
	}
	return models.BackgroundConfig{
		Enabled:            b.Enabled,
		Backgrounds:        backgrounds,
		Scripts:            b.Scripts,
		CSS:                b.CSS,
		ArticleBg:          b.ArticleBg,
		ArticleBlurEnabled: b.ArticleBlurEnabled,
		ArticleBlur:        b.ArticleBlur,
		ArticleShadow:      b.ArticleShadow,
		ArticleBorder:      b.ArticleBorder,
		ArticleRadius:      b.ArticleRadius,
	}
}

func (f *tomlFontConfig) toFontConfig() models.FontConfig {
	return models.FontConfig{
		Family:        f.Family,
		HeadingFamily: f.HeadingFamily,
		CodeFamily:    f.CodeFamily,
		Size:          f.Size,
		LineHeight:    f.LineHeight,
		GoogleFonts:   f.GoogleFonts,
		CustomURLs:    f.CustomURLs,
	}
}

func (f *tomlFeedConfig) toFeedConfig() models.FeedConfig {
	return models.FeedConfig{
		Slug:            f.Slug,
		Presentation:    f.Presentation,
		Title:           f.Title,
		Description:     f.Description,
		Robots:          f.Robots,
		Filter:          f.Filter,
		Sort:            f.Sort,
		Reverse:         f.Reverse,
		Views:           f.Views,
		ShortsPath:      f.ShortsPath,
		Primary:         f.Primary,
		IncludePrivate:  f.IncludePrivate,
		Private:         f.Private,
		Sidebar:         f.Sidebar,
		ItemsPerPage:    f.ItemsPerPage,
		OrphanThreshold: f.OrphanThreshold,
		Limit:           f.Limit,
		Offset:          f.Offset,
		PaginationType:  models.PaginationType(f.PaginationType),
		ArchiveDisabled: f.ArchiveDisabled,
		Formats:         f.Formats.toFeedFormats(),
		Templates:       f.Templates.toFeedTemplates(),
	}
}

func (f *tomlFeedFormats) toFeedFormats() models.FeedFormats {
	formats := models.FeedFormats{}
	if f.HTML != nil {
		formats.HTML = *f.HTML
	}
	if f.SimpleHTML != nil {
		formats.SimpleHTML = *f.SimpleHTML
	}
	if f.RSS != nil {
		formats.RSS = *f.RSS
	}
	if f.Atom != nil {
		formats.Atom = *f.Atom
	}
	if f.JSON != nil {
		formats.JSON = *f.JSON
	}
	if f.Markdown != nil {
		formats.Markdown = *f.Markdown
	}
	if f.Text != nil {
		formats.Text = *f.Text
	}
	if f.Sitemap != nil {
		formats.Sitemap = *f.Sitemap
	}
	return formats
}

func (t *tomlFeedTemplates) toFeedTemplates() models.FeedTemplates {
	return models.FeedTemplates{
		HTML:       t.HTML,
		SimpleHTML: t.SimpleHTML,
		RSS:        t.RSS,
		Atom:       t.Atom,
		JSON:       t.JSON,
		Card:       t.Card,
		Sitemap:    t.Sitemap,
	}
}

func (d *tomlFeedDefaults) toFeedDefaults() models.FeedDefaults {
	return models.FeedDefaults{
		ItemsPerPage:    d.ItemsPerPage,
		Presentation:    d.Presentation,
		OrphanThreshold: d.OrphanThreshold,
		PaginationType:  models.PaginationType(d.PaginationType),
		Views:           d.Views,
		Formats:         d.Formats.toFeedFormats(),
		Templates:       d.Templates.toFeedTemplates(),
		Syndication: models.SyndicationConfig{
			MaxItems:             d.Syndication.MaxItems,
			IncludeContent:       d.Syndication.IncludeContent,
			SiteArchiveDisabled:  d.Syndication.SiteArchiveDisabled,
			FeedArchivesDisabled: d.Syndication.FeedArchivesDisabled,
		},
	}
}

// yamlConfig is an internal struct for parsing YAML configuration.
type yamlConfig struct {
	Fontpack        string                    `yaml:"fontpack"`
	FontpacksFile   string                    `yaml:"fontpacks_file"`
	OutputDir       string                    `yaml:"output_dir"`
	URL             string                    `yaml:"url"`
	Title           string                    `yaml:"title"`
	Description     string                    `yaml:"description"`
	Author          string                    `yaml:"author"`
	Language        string                    `yaml:"language"`
	AuthorURL       string                    `yaml:"author_url"`
	ManagingEditor  string                    `yaml:"managing_editor"`
	WebMaster       string                    `yaml:"webmaster"`
	Copyright       string                    `yaml:"copyright"`
	License         interface{}               `yaml:"license"`
	AssetsDir       string                    `yaml:"assets_dir"`
	Assets          models.AssetsConfig       `yaml:"assets"`
	TemplatesDir    string                    `yaml:"templates_dir"`
	Nav             []yamlNavItem             `yaml:"nav"`
	Footer          yamlFooterConfig          `yaml:"footer"`
	Hooks           []string                  `yaml:"hooks"`
	DisabledHooks   []string                  `yaml:"disabled_hooks"`
	Glob            yamlGlobConfig            `yaml:"glob"`
	Markdown        yamlMarkdownConfig        `yaml:"markdown"`
	Feeds           []yamlFeedConfig          `yaml:"feeds"`
	FeedDefaults    yamlFeedDefaults          `yaml:"feed_defaults"`
	Concurrency     int                       `yaml:"concurrency"`
	Theme           yamlThemeConfig           `yaml:"theme"`
	ThemeCalendar   parsedThemeCalendar       `yaml:"theme_calendar"`
	Head            parsedHead                `yaml:"head"`
	Presets         parsedTemplatePresets     `yaml:"template_presets"`
	Defaults        parsedDefaultTemplates    `yaml:"default_templates"`
	ErrorPages      parsedErrorPages          `yaml:"error_pages"`
	ResourceHints   parsedResourceHints       `yaml:"resource_hints"`
	PostFormats     yamlPostFormatsConfig     `yaml:"post_formats"`
	WellKnown       yamlWellKnownConfig       `yaml:"well_known"`
	IndieAuth       yamlIndieAuthConfig       `yaml:"indieauth"`
	Webmention      yamlWebmentionConfig      `yaml:"webmention"`
	Search          models.SearchConfig       `yaml:"search"`
	SEO             yamlSEOConfig             `yaml:"seo"`
	Components      yamlComponentsConfig      `yaml:"components"`
	Layout          yamlLayoutConfig          `yaml:"layout"`
	Sidebar         yamlSidebarConfig         `yaml:"sidebar"`
	Toc             yamlTocConfig             `yaml:"toc"`
	Header          yamlHeaderLayoutConfig    `yaml:"header"`
	Blogroll        yamlBlogrollConfig        `yaml:"blogroll"`
	Tags            yamlTagsConfig            `yaml:"tags"`
	FeedsPage       yamlFeedsPageConfig       `yaml:"feeds_page"`
	Encryption      yamlEncryptionConfig      `yaml:"encryption"`
	TagAggregator   yamlTagAggregatorConfig   `yaml:"tag_aggregator"`
	Mentions        yamlMentionsConfig        `yaml:"mentions"`
	WebSub          yamlWebSubConfig          `yaml:"websub"`
	Shortcuts       yamlShortcutsConfig       `yaml:"shortcuts"`
	ViewTransitions yamlViewTransitionsConfig `yaml:"view_transitions"`
	Authors         yamlAuthorsConfig         `yaml:"authors"`
	Garden          yamlGardenConfig          `yaml:"garden"`
	Templates       yamlTemplatesConfig       `yaml:"templates"`
	Images          models.ImagesConfig       `yaml:"images"`
	BuilderAdmin    models.BuilderAdminConfig `yaml:"builder_admin"`
}

type yamlNavItem struct {
	Label    string `yaml:"label"`
	URL      string `yaml:"url"`
	External bool   `yaml:"external"`
}

type yamlTemplatesConfig struct {
	Media yamlTemplatesMediaConfig `yaml:"media"`
}

type yamlTemplatesMediaConfig struct {
	TrustedDomains []string `yaml:"trusted_domains"`
}

func (t *yamlTemplatesConfig) toTemplatesConfig() models.TemplatesConfig {
	cfg := models.NewTemplatesConfig()
	if t.Media.TrustedDomains != nil {
		cfg.Media.TrustedDomains = append([]string{}, t.Media.TrustedDomains...)
	}
	return cfg
}

type yamlFooterConfig struct {
	Text          string `yaml:"text"`
	ShowCopyright *bool  `yaml:"show_copyright"`
}

type yamlGlobConfig struct {
	Patterns     []string       `yaml:"patterns"`
	UseGitignore *bool          `yaml:"use_gitignore"`
	SlugMode     string         `yaml:"slug_mode"`
	SlugRules    []yamlSlugRule `yaml:"slug_rules"`
}

type yamlSlugRule struct {
	Prefix string `yaml:"prefix"`
	Mode   string `yaml:"mode"`
}

type yamlMarkdownConfig struct {
	Extensions []string        `yaml:"extensions"`
	Highlight  parsedHighlight `yaml:"highlight"`
}

type yamlFeedConfig struct {
	Slug            string            `yaml:"slug"`
	Presentation    string            `yaml:"presentation"`
	Title           string            `yaml:"title"`
	Description     string            `yaml:"description"`
	Robots          string            `yaml:"robots"`
	Filter          string            `yaml:"filter"`
	Sort            string            `yaml:"sort"`
	Reverse         bool              `yaml:"reverse"`
	Views           []string          `yaml:"views"`
	ShortsPath      string            `yaml:"shorts_path"`
	Primary         bool              `yaml:"primary"`
	IncludePrivate  bool              `yaml:"include_private"`
	Private         bool              `yaml:"private"`
	Sidebar         *bool             `yaml:"sidebar"`
	ItemsPerPage    int               `yaml:"items_per_page"`
	OrphanThreshold int               `yaml:"orphan_threshold"`
	Limit           int               `yaml:"limit"`
	Offset          int               `yaml:"offset"`
	PaginationType  string            `yaml:"pagination_type"`
	ArchiveDisabled bool              `yaml:"archive_disabled"`
	Formats         yamlFeedFormats   `yaml:"formats"`
	Templates       yamlFeedTemplates `yaml:"templates"`
}

type yamlFeedFormats struct {
	HTML       *bool `yaml:"html"`
	SimpleHTML *bool `yaml:"simple_html"`
	RSS        *bool `yaml:"rss"`
	Atom       *bool `yaml:"atom"`
	JSON       *bool `yaml:"json"`
	Markdown   *bool `yaml:"markdown"`
	Text       *bool `yaml:"text"`
	Sitemap    *bool `yaml:"sitemap"`
}

type yamlFeedTemplates struct {
	HTML       string `yaml:"html"`
	SimpleHTML string `yaml:"simple_html"`
	RSS        string `yaml:"rss"`
	Atom       string `yaml:"atom"`
	JSON       string `yaml:"json"`
	Card       string `yaml:"card"`
	Sitemap    string `yaml:"sitemap"`
}

type yamlFeedDefaults struct {
	ItemsPerPage    int                   `yaml:"items_per_page"`
	Presentation    string                `yaml:"presentation"`
	OrphanThreshold int                   `yaml:"orphan_threshold"`
	PaginationType  string                `yaml:"pagination_type"`
	Views           []string              `yaml:"views"`
	Formats         yamlFeedFormats       `yaml:"formats"`
	Templates       yamlFeedTemplates     `yaml:"templates"`
	Syndication     yamlSyndicationConfig `yaml:"syndication"`
}

type yamlSyndicationConfig struct {
	MaxItems             int  `yaml:"max_items"`
	IncludeContent       bool `yaml:"include_content"`
	SiteArchiveDisabled  bool `yaml:"site_archive_disabled"`
	FeedArchivesDisabled bool `yaml:"feed_archives_disabled"`
}

type yamlPostFormatsConfig struct {
	HTML     *bool `yaml:"html"`
	Markdown bool  `yaml:"markdown"`
	Text     bool  `yaml:"text"`
	ANSI     bool  `yaml:"ansi"`
	OG       bool  `yaml:"og"`
}

type yamlWellKnownConfig struct {
	Enabled         *bool    `yaml:"enabled"`
	AutoGenerate    []string `yaml:"auto_generate"`
	SSHFingerprint  string   `yaml:"ssh_fingerprint"`
	KeybaseUsername string   `yaml:"keybase_username"`
}

type yamlSEOConfig struct {
	TwitterHandle  string `yaml:"twitter_handle"`
	DefaultImage   string `yaml:"default_image"`
	LogoURL        string `yaml:"logo_url"`
	AuthorImage    string `yaml:"author_image"`
	OGImageService string `yaml:"og_image_service"`
}

type yamlIndieAuthConfig struct {
	Enabled               bool   `yaml:"enabled"`
	AuthorizationEndpoint string `yaml:"authorization_endpoint"`
	TokenEndpoint         string `yaml:"token_endpoint"`
	MeURL                 string `yaml:"me_url"`
}

type yamlWebmentionConfig struct {
	Enabled  bool   `yaml:"enabled"`
	Endpoint string `yaml:"endpoint"`
}

type yamlTagsConfig struct {
	Enabled     *bool    `yaml:"enabled"`
	Blacklist   []string `yaml:"blacklist"`
	Private     []string `yaml:"private"`
	Title       string   `yaml:"title"`
	Description string   `yaml:"description"`
	Template    string   `yaml:"template"`
	SlugPrefix  string   `yaml:"slug_prefix"`
}

func (t *yamlTagsConfig) toTagsConfig() models.TagsConfig {
	defaults := models.NewTagsConfig()

	config := models.TagsConfig{
		Enabled:     t.Enabled,
		Blacklist:   t.Blacklist,
		Private:     t.Private,
		Title:       t.Title,
		Description: t.Description,
		Template:    t.Template,
		SlugPrefix:  t.SlugPrefix,
	}

	// Apply defaults if not set
	if config.Enabled == nil {
		config.Enabled = defaults.Enabled
	}
	if config.Title == "" {
		config.Title = defaults.Title
	}
	if config.Template == "" {
		config.Template = defaults.Template
	}
	if config.SlugPrefix == "" {
		config.SlugPrefix = defaults.SlugPrefix
	}

	return config
}

type yamlFeedsPageConfig struct {
	Enabled          *bool    `yaml:"enabled"`
	Title            string   `yaml:"title"`
	Description      string   `yaml:"description"`
	Template         string   `yaml:"template"`
	SlugPrefix       string   `yaml:"slug_prefix"`
	Robots           string   `yaml:"robots"`
	ShowPrivateFeeds []string `yaml:"show_private_feeds"`
}

func (t *yamlFeedsPageConfig) toFeedsPageConfig() models.FeedsPageConfig {
	defaults := models.NewFeedsPageConfig()

	config := models.FeedsPageConfig{
		Enabled:          t.Enabled,
		Title:            t.Title,
		Description:      t.Description,
		Template:         t.Template,
		SlugPrefix:       t.SlugPrefix,
		Robots:           t.Robots,
		ShowPrivateFeeds: append([]string{}, t.ShowPrivateFeeds...),
	}

	if config.Enabled == nil {
		config.Enabled = defaults.Enabled
	}
	if config.Title == "" {
		config.Title = defaults.Title
	}
	if config.Description == "" {
		config.Description = defaults.Description
	}
	if config.Template == "" {
		config.Template = defaults.Template
	}
	if config.SlugPrefix == "" {
		config.SlugPrefix = defaults.SlugPrefix
	}

	return config
}

type yamlEncryptionConfig struct {
	Enabled               *bool             `yaml:"enabled"`
	DefaultKey            string            `yaml:"default_key"`
	DecryptionHint        string            `yaml:"decryption_hint"`
	PrivateTags           map[string]string `yaml:"private_tags"`
	EnforceStrength       *bool             `yaml:"enforce_strength"`
	MinEstimatedCrackTime string            `yaml:"min_estimated_crack_time"`
	MinPasswordLength     int               `yaml:"min_password_length"`
}

func (e *yamlEncryptionConfig) toEncryptionConfig() models.EncryptionConfig {
	defaults := models.NewEncryptionConfig()

	config := models.EncryptionConfig{
		DefaultKey:            e.DefaultKey,
		DecryptionHint:        e.DecryptionHint,
		PrivateTags:           e.PrivateTags,
		MinEstimatedCrackTime: e.MinEstimatedCrackTime,
		MinPasswordLength:     e.MinPasswordLength,
	}

	// Apply defaults for unset values
	if e.Enabled != nil {
		config.Enabled = *e.Enabled
	} else {
		config.Enabled = defaults.Enabled
	}
	if config.DefaultKey == "" {
		config.DefaultKey = defaults.DefaultKey
	}
	if e.EnforceStrength != nil {
		config.EnforceStrength = *e.EnforceStrength
	} else {
		config.EnforceStrength = defaults.EnforceStrength
	}
	if config.MinEstimatedCrackTime == "" {
		config.MinEstimatedCrackTime = defaults.MinEstimatedCrackTime
	}
	if config.MinPasswordLength == 0 {
		config.MinPasswordLength = defaults.MinPasswordLength // pragma: allowlist secret
	}

	return config
}

type yamlTagAggregatorConfig struct {
	Enabled        *bool               `yaml:"enabled"`
	Synonyms       map[string][]string `yaml:"synonyms"`
	Additional     map[string][]string `yaml:"additional"`
	GenerateReport bool                `yaml:"generate_report"`
}

func (t *yamlTagAggregatorConfig) toTagAggregatorConfig() models.TagAggregatorConfig {
	defaults := models.NewTagAggregatorConfig()

	config := models.TagAggregatorConfig{
		Enabled:        t.Enabled,
		Synonyms:       t.Synonyms,
		Additional:     t.Additional,
		GenerateReport: t.GenerateReport,
	}

	// Apply defaults if not set
	if config.Enabled == nil {
		config.Enabled = defaults.Enabled
	}

	return config
}

// Mentions-related YAML structs

type yamlMentionsConfig struct {
	Enabled            *bool                   `yaml:"enabled"`
	CSSClass           string                  `yaml:"css_class"`
	FromPosts          []yamlMentionPostSource `yaml:"from_posts"`
	CacheDir           string                  `yaml:"cache_dir"`
	CacheDuration      string                  `yaml:"cache_duration"`
	Timeout            int                     `yaml:"timeout"`
	ConcurrentRequests int                     `yaml:"concurrent_requests"`
}

type yamlMentionPostSource struct {
	Filter       string `yaml:"filter"`
	HandleField  string `yaml:"handle_field"`
	AliasesField string `yaml:"aliases_field"`
	AvatarField  string `yaml:"avatar_field"`
}

func (m *yamlMentionsConfig) toMentionsConfig() models.MentionsConfig {
	defaults := models.NewMentionsConfig()

	config := models.MentionsConfig{
		Enabled:            m.Enabled,
		CSSClass:           m.CSSClass,
		CacheDir:           m.CacheDir,
		CacheDuration:      m.CacheDuration,
		Timeout:            m.Timeout,
		ConcurrentRequests: m.ConcurrentRequests,
	}

	// Apply defaults for unset values
	if config.Enabled == nil {
		config.Enabled = defaults.Enabled
	}
	if config.CSSClass == "" {
		config.CSSClass = defaults.CSSClass
	}
	if config.CacheDir == "" {
		config.CacheDir = defaults.CacheDir
	}
	if config.CacheDuration == "" {
		config.CacheDuration = defaults.CacheDuration
	}
	if config.Timeout == 0 {
		config.Timeout = defaults.Timeout
	}
	if config.ConcurrentRequests == 0 {
		config.ConcurrentRequests = defaults.ConcurrentRequests
	}

	// Convert from_posts sources
	for _, src := range m.FromPosts {
		aliasesField := src.AliasesField
		if aliasesField == "" {
			aliasesField = defaultAliasesField
		}
		config.FromPosts = append(config.FromPosts, models.MentionPostSource{
			Filter:       src.Filter,
			HandleField:  src.HandleField,
			AliasesField: aliasesField,
			AvatarField:  src.AvatarField,
		})
	}

	// Fall back to defaults when user has mentions config but no from_posts
	if len(config.FromPosts) == 0 {
		config.FromPosts = defaults.FromPosts
	}

	return config
}

type yamlWebSubConfig struct {
	Enabled *bool    `yaml:"enabled"`
	Hubs    []string `yaml:"hubs"`
}

func (w *yamlWebSubConfig) toWebSubConfig() models.WebSubConfig {
	defaults := models.NewWebSubConfig()
	config := models.WebSubConfig{
		Enabled: w.Enabled,
		Hubs:    w.Hubs,
	}
	if config.Enabled == nil {
		config.Enabled = defaults.Enabled
	}
	if config.Hubs == nil {
		config.Hubs = defaults.Hubs
	}
	return config
}

type yamlShortcutsConfig struct {
	Navigation map[string]string `yaml:"navigation"`
}

type yamlViewTransitionsConfig struct {
	Enabled       *bool    `yaml:"enabled"`
	Native        *bool    `yaml:"native"`
	Debug         *bool    `yaml:"debug"`
	Duration      int      `yaml:"duration"`
	UpdateMeta    *bool    `yaml:"update_meta"`
	ScrollToTop   *bool    `yaml:"scroll_to_top"`
	SkipClasses   []string `yaml:"skip_classes"`
	SkipSelectors []string `yaml:"skip_selectors"`
}

func (s *yamlShortcutsConfig) toShortcutsConfig() models.ShortcutsConfig {
	config := models.ShortcutsConfig{
		Navigation: s.Navigation,
	}
	if config.Navigation == nil {
		config.Navigation = make(map[string]string)
	}
	return config
}

func (v *yamlViewTransitionsConfig) toViewTransitionsConfig() models.ViewTransitionsConfig {
	defaults := models.NewViewTransitionsConfig()
	config := models.ViewTransitionsConfig{
		Enabled:       v.Enabled,
		Native:        v.Native,
		Debug:         v.Debug,
		Duration:      v.Duration,
		UpdateMeta:    v.UpdateMeta,
		ScrollToTop:   v.ScrollToTop,
		SkipClasses:   v.SkipClasses,
		SkipSelectors: v.SkipSelectors,
	}
	if config.Enabled == nil {
		config.Enabled = defaults.Enabled
	}
	if config.Debug == nil {
		config.Debug = defaults.Debug
	}
	if config.Duration == 0 {
		config.Duration = defaults.Duration
	}
	if config.UpdateMeta == nil {
		config.UpdateMeta = defaults.UpdateMeta
	}
	if config.ScrollToTop == nil {
		config.ScrollToTop = defaults.ScrollToTop
	}
	if config.SkipClasses == nil {
		config.SkipClasses = defaults.SkipClasses
	}
	if config.SkipSelectors == nil {
		config.SkipSelectors = defaults.SkipSelectors
	}
	return config
}

type yamlAuthorsConfig struct {
	GeneratePages bool                  `yaml:"generate_pages"`
	URLPattern    string                `yaml:"url_pattern"`
	FeedsEnabled  bool                  `yaml:"feeds_enabled"`
	Authors       map[string]yamlAuthor `yaml:"authors"`
}

type yamlAuthor struct {
	Name          string            `yaml:"name"`
	Bio           string            `yaml:"bio"`
	Email         string            `yaml:"email"`
	Avatar        string            `yaml:"avatar"`
	URL           string            `yaml:"url"`
	Social        map[string]string `yaml:"social"`
	Guest         bool              `yaml:"guest"`
	Active        bool              `yaml:"active"`
	Default       bool              `yaml:"default"`
	Contributions []string          `yaml:"contributions"`
	Role          string            `yaml:"role"`
	Contribution  string            `yaml:"contribution"`
}

func (a *yamlAuthorsConfig) toAuthorsConfig() models.AuthorsConfig {
	config := models.AuthorsConfig{
		GeneratePages: a.GeneratePages,
		URLPattern:    a.URLPattern,
		FeedsEnabled:  a.FeedsEnabled,
	}
	if a.Authors != nil {
		config.Authors = make(map[string]models.Author, len(a.Authors))
		for id := range a.Authors {
			author := a.Authors[id]
			ma := models.Author{
				ID:      id,
				Name:    author.Name,
				Guest:   author.Guest,
				Active:  author.Active,
				Default: author.Default,
				Social:  author.Social,
			}
			if author.Bio != "" {
				bio := author.Bio
				ma.Bio = &bio
			}
			if author.Email != "" {
				email := author.Email
				ma.Email = &email
			}
			if author.Avatar != "" {
				avatar := author.Avatar
				ma.Avatar = &avatar
			}
			if author.URL != "" {
				url := author.URL
				ma.URL = &url
			}
			if author.Role != "" {
				role := author.Role
				ma.Role = &role
			}
			if author.Contribution != "" {
				contribution := author.Contribution
				ma.Contribution = &contribution
			}
			if len(author.Contributions) > 0 {
				ma.Contributions = author.Contributions
			}
			config.Authors[id] = ma
		}
	}
	return config
}

type yamlGardenConfig struct {
	Enabled      *bool    `yaml:"enabled"`
	Path         string   `yaml:"path"`
	ExportJSON   *bool    `yaml:"export_json"`
	RenderPage   *bool    `yaml:"render_page"`
	IncludeTags  *bool    `yaml:"include_tags"`
	IncludePosts *bool    `yaml:"include_posts"`
	MaxNodes     int      `yaml:"max_nodes"`
	ExcludeTags  []string `yaml:"exclude_tags"`
	Template     string   `yaml:"template"`
	Title        string   `yaml:"title"`
	Description  string   `yaml:"description"`
}

func (g *yamlGardenConfig) toGardenConfig() models.GardenConfig {
	defaults := models.NewGardenConfig()

	config := models.GardenConfig{
		Enabled:      g.Enabled,
		Path:         g.Path,
		ExportJSON:   g.ExportJSON,
		RenderPage:   g.RenderPage,
		IncludeTags:  g.IncludeTags,
		IncludePosts: g.IncludePosts,
		MaxNodes:     g.MaxNodes,
		ExcludeTags:  g.ExcludeTags,
		Template:     g.Template,
		Title:        g.Title,
		Description:  g.Description,
	}

	// Apply defaults for unset fields
	if config.Enabled == nil {
		config.Enabled = defaults.Enabled
	}
	if config.ExportJSON == nil {
		config.ExportJSON = defaults.ExportJSON
	}
	if config.RenderPage == nil {
		config.RenderPage = defaults.RenderPage
	}
	if config.IncludeTags == nil {
		config.IncludeTags = defaults.IncludeTags
	}
	if config.IncludePosts == nil {
		config.IncludePosts = defaults.IncludePosts
	}
	if config.Path == "" {
		config.Path = defaults.Path
	}
	if config.Template == "" {
		config.Template = defaults.Template
	}
	if config.Title == "" {
		config.Title = defaults.Title
	}

	return config
}

type yamlThemeConfig struct {
	ContractVersion     int                      `yaml:"contract_version"`
	Fontpack            string                   `yaml:"fontpack"`
	Name                string                   `yaml:"name"`
	Aesthetic           string                   `yaml:"aesthetic"`
	Palette             string                   `yaml:"palette"`
	PaletteLight        string                   `yaml:"palette_light"`
	PaletteDark         string                   `yaml:"palette_dark"`
	FallbackMode        string                   `yaml:"fallback_mode"`
	TextSize            string                   `yaml:"text_size"`
	Seasonal            bool                     `yaml:"seasonal"`
	ShowTextSizeControl *bool                    `yaml:"show_text_size_control"`
	SeedColor           string                   `yaml:"seed_color"`
	Variables           map[string]string        `yaml:"variables"`
	CustomCSS           string                   `yaml:"custom_css"`
	Background          yamlBackgroundConfig     `yaml:"background"`
	Font                yamlFontConfig           `yaml:"font"`
	Switcher            yamlThemeSwitcherConfig  `yaml:"switcher"`
	Texture             yamlTextureConfig        `yaml:"texture"`
	HeadingTexture      yamlHeadingTextureConfig `yaml:"heading_texture"`
	Motif               yamlMotifConfig          `yaml:"motif"`
}

type yamlTextureConfig struct {
	Kind     string  `yaml:"kind"`
	ColorMix float64 `yaml:"color_mix"`
	Scale    float64 `yaml:"scale"`
	Scope    string  `yaml:"scope"`
}
type yamlHeadingTextureConfig struct {
	Kind     string  `yaml:"kind"`
	ColorMix float64 `yaml:"color_mix"`
	Scale    float64 `yaml:"scale"`
}
type yamlMotifConfig struct {
	Kind      string  `yaml:"kind"`
	Glyph     string  `yaml:"glyph"`
	Size      string  `yaml:"size"`
	Gap       string  `yaml:"gap"`
	RowOffset float64 `yaml:"row_offset"`
	Wobble    float64 `yaml:"wobble"`
	Scatter   float64 `yaml:"scatter"`
	Layer     string  `yaml:"layer"`
	Color     string  `yaml:"color"`
	ColorMix  float64 `yaml:"color_mix"`
	URL       string  `yaml:"url"`
}

type yamlThemeSwitcherConfig struct {
	Enabled    *bool    `yaml:"enabled"`
	ModeToggle *bool    `yaml:"mode_toggle"`
	IncludeAll *bool    `yaml:"include_all"`
	Include    []string `yaml:"include"`
	Exclude    []string `yaml:"exclude"`
	Position   string   `yaml:"position"`
}

type yamlBackgroundConfig struct {
	Enabled            *bool                   `yaml:"enabled"`
	Backgrounds        []yamlBackgroundElement `yaml:"backgrounds"`
	Scripts            []string                `yaml:"scripts"`
	CSS                string                  `yaml:"css"`
	ArticleBg          string                  `yaml:"article_bg"`
	ArticleBlurEnabled *bool                   `yaml:"article_blur_enabled"`
	ArticleBlur        string                  `yaml:"article_blur"`
	ArticleShadow      string                  `yaml:"article_shadow"`
	ArticleBorder      string                  `yaml:"article_border"`
	ArticleRadius      string                  `yaml:"article_radius"`
}

type yamlBackgroundElement struct {
	HTML   string `yaml:"html"`
	ZIndex int    `yaml:"z_index"`
}

type yamlFontConfig struct {
	Family        string   `yaml:"family"`
	HeadingFamily string   `yaml:"heading_family"`
	CodeFamily    string   `yaml:"code_family"`
	Size          string   `yaml:"size"`
	LineHeight    string   `yaml:"line_height"`
	GoogleFonts   []string `yaml:"google_fonts"`
	CustomURLs    []string `yaml:"custom_urls"`
}

func (t *yamlThemeConfig) toThemeConfig() models.ThemeConfig {
	variables := t.Variables
	if variables == nil {
		variables = make(map[string]string)
	}
	return models.ThemeConfig{
		ContractVersion:     t.ContractVersion,
		Fontpack:            t.Fontpack,
		Name:                t.Name,
		Aesthetic:           t.Aesthetic,
		Palette:             t.Palette,
		PaletteLight:        t.PaletteLight,
		PaletteDark:         t.PaletteDark,
		FallbackMode:        t.FallbackMode,
		TextSize:            t.TextSize,
		Seasonal:            t.Seasonal,
		ShowTextSizeControl: t.ShowTextSizeControl,
		SeedColor:           t.SeedColor,
		Variables:           variables,
		CustomCSS:           t.CustomCSS,
		Background:          t.Background.toBackgroundConfig(),
		Font:                t.Font.toFontConfig(),
		Switcher:            t.Switcher.toThemeSwitcherConfig(),
		Texture:             models.ThemeTextureConfig{Kind: t.Texture.Kind, ColorMix: t.Texture.ColorMix, Scale: t.Texture.Scale, Scope: t.Texture.Scope},
		HeadingTexture:      models.ThemeHeadingTextureConfig{Kind: t.HeadingTexture.Kind, ColorMix: t.HeadingTexture.ColorMix, Scale: t.HeadingTexture.Scale},
		Motif:               models.ThemeMotifConfig{Kind: t.Motif.Kind, Glyph: t.Motif.Glyph, Size: t.Motif.Size, Gap: t.Motif.Gap, RowOffset: t.Motif.RowOffset, Wobble: t.Motif.Wobble, Scatter: t.Motif.Scatter, Layer: t.Motif.Layer, Color: t.Motif.Color, ColorMix: t.Motif.ColorMix, URL: t.Motif.URL},
	}
}

func (s *yamlThemeSwitcherConfig) toThemeSwitcherConfig() models.ThemeSwitcherConfig {
	return models.ThemeSwitcherConfig{
		Enabled:    s.Enabled,
		ModeToggle: s.ModeToggle,
		IncludeAll: s.IncludeAll,
		Include:    s.Include,
		Exclude:    s.Exclude,
		Position:   s.Position,
	}
}

func (b *yamlBackgroundConfig) toBackgroundConfig() models.BackgroundConfig {
	backgrounds := make([]models.BackgroundElement, len(b.Backgrounds))
	for i, bg := range b.Backgrounds {
		backgrounds[i] = models.BackgroundElement{
			HTML:   bg.HTML,
			ZIndex: bg.ZIndex,
		}
	}
	return models.BackgroundConfig{
		Enabled:            b.Enabled,
		Backgrounds:        backgrounds,
		Scripts:            b.Scripts,
		CSS:                b.CSS,
		ArticleBg:          b.ArticleBg,
		ArticleBlurEnabled: b.ArticleBlurEnabled,
		ArticleBlur:        b.ArticleBlur,
		ArticleShadow:      b.ArticleShadow,
		ArticleBorder:      b.ArticleBorder,
		ArticleRadius:      b.ArticleRadius,
	}
}

func (f *yamlFontConfig) toFontConfig() models.FontConfig {
	return models.FontConfig{
		Family:        f.Family,
		HeadingFamily: f.HeadingFamily,
		CodeFamily:    f.CodeFamily,
		Size:          f.Size,
		LineHeight:    f.LineHeight,
		GoogleFonts:   f.GoogleFonts,
		CustomURLs:    f.CustomURLs,
	}
}

func (s *yamlSEOConfig) toSEOConfig() models.SEOConfig {
	return models.SEOConfig{
		TwitterHandle:  s.TwitterHandle,
		DefaultImage:   s.DefaultImage,
		LogoURL:        s.LogoURL,
		AuthorImage:    s.AuthorImage,
		OGImageService: s.OGImageService,
	}
}

func (i *yamlIndieAuthConfig) toIndieAuthConfig() models.IndieAuthConfig {
	return models.IndieAuthConfig{
		Enabled:               i.Enabled,
		AuthorizationEndpoint: i.AuthorizationEndpoint,
		TokenEndpoint:         i.TokenEndpoint,
		MeURL:                 i.MeURL,
	}
}

func (w *yamlWebmentionConfig) toWebmentionConfig() models.WebmentionConfig {
	return models.WebmentionConfig{
		Enabled:  w.Enabled,
		Endpoint: w.Endpoint,
	}
}

type yamlComponentsConfig struct {
	Nav            yamlNavComponentConfig      `yaml:"nav"`
	Footer         yamlFooterComponentConfig   `yaml:"footer"`
	DocSidebar     yamlDocSidebarConfig        `yaml:"doc_sidebar"`
	FeedSidebar    yamlFeedSidebarConfig       `yaml:"feed_sidebar"`
	ContentSidebar yamlContentSidebarConfig    `yaml:"content_sidebar"`
	CardRouter     yamlCardRouterConfig        `yaml:"card_router"`
	Share          yamlShareComponentConfig    `yaml:"share"`
	PostCopy       yamlPostCopyComponentConfig `yaml:"post_copy"`
	PostConn       yamlPostConnectionsConfig   `yaml:"post_connections"`
}

type yamlPostCopyComponentConfig struct {
	Enabled *bool `yaml:"enabled"`
}

type yamlPostConnectionsConfig struct {
	Enabled       *bool    `yaml:"enabled"`
	Display       []string `yaml:"display"`
	MinLinks      int      `yaml:"min_links"`
	MaxLinks      int      `yaml:"max_links"`
	GraphMinLinks int      `yaml:"graph_min_links"`
	GraphMaxLinks int      `yaml:"graph_max_links"`
	ListMinLinks  int      `yaml:"list_min_links"`
	ListMaxLinks  int      `yaml:"list_max_links"`
	InlinksLimit  int      `yaml:"inlinks_limit"`
	OutlinksLimit int      `yaml:"outlinks_limit"`
}

type yamlNavComponentConfig struct {
	Enabled  *bool         `yaml:"enabled"`
	Position string        `yaml:"position"`
	Style    string        `yaml:"style"`
	Items    []yamlNavItem `yaml:"items"`
}

type yamlFooterComponentConfig struct {
	Enabled       *bool         `yaml:"enabled"`
	Text          string        `yaml:"text"`
	ShowCopyright *bool         `yaml:"show_copyright"`
	Links         []yamlNavItem `yaml:"links"`
}

type yamlDocSidebarConfig struct {
	Enabled  *bool  `yaml:"enabled"`
	Position string `yaml:"position"`
	Width    string `yaml:"width"`
	MinDepth int    `yaml:"min_depth"`
	MaxDepth int    `yaml:"max_depth"`
}

type yamlFeedSidebarConfig struct {
	Enabled  *bool    `yaml:"enabled"`
	Position string   `yaml:"position"`
	Width    string   `yaml:"width"`
	Title    string   `yaml:"title"`
	Feeds    []string `yaml:"feeds"`
	MaxPosts int      `yaml:"max_posts"`
}

type yamlContentSidebarConfig struct {
	Enabled  *bool  `yaml:"enabled"`
	Position string `yaml:"position"`
	Width    string `yaml:"width"`
	Slug     string `yaml:"slug"`
}

type yamlCardRouterConfig struct {
	Mappings map[string]string `yaml:"mappings"`
}

type yamlShareComponentConfig struct {
	Enabled   *bool                            `yaml:"enabled"`
	Platforms []string                         `yaml:"platforms"`
	Position  string                           `yaml:"position"`
	Title     string                           `yaml:"title"`
	Custom    map[string]yamlSharePlatformItem `yaml:"custom"`
}

type yamlSharePlatformItem struct {
	Name string `yaml:"name"`
	Icon string `yaml:"icon"`
	URL  string `yaml:"url"`
}

// Layout-related YAML structs

type yamlLayoutConfig struct {
	Name     string                  `yaml:"name"`
	Paths    map[string]string       `yaml:"paths"`
	Feeds    map[string]string       `yaml:"feeds"`
	Docs     yamlDocsLayoutConfig    `yaml:"docs"`
	Blog     yamlBlogLayoutConfig    `yaml:"blog"`
	Landing  yamlLandingLayoutConfig `yaml:"landing"`
	Bare     yamlBareLayoutConfig    `yaml:"bare"`
	Defaults yamlLayoutDefaults      `yaml:"defaults"`
}

type yamlLayoutDefaults struct {
	ContentMaxWidth string `yaml:"content_max_width"`
	HeaderSticky    *bool  `yaml:"header_sticky"`
	FooterSticky    *bool  `yaml:"footer_sticky"`
}

type yamlDocsLayoutConfig struct {
	SidebarPosition    string `yaml:"sidebar_position"`
	SidebarWidth       string `yaml:"sidebar_width"`
	SidebarCollapsible *bool  `yaml:"sidebar_collapsible"`
	SidebarDefaultOpen *bool  `yaml:"sidebar_default_open"`
	TocPosition        string `yaml:"toc_position"`
	TocWidth           string `yaml:"toc_width"`
	TocCollapsible     *bool  `yaml:"toc_collapsible"`
	TocDefaultOpen     *bool  `yaml:"toc_default_open"`
	ContentMaxWidth    string `yaml:"content_max_width"`
	HeaderStyle        string `yaml:"header_style"`
	FooterStyle        string `yaml:"footer_style"`
}

type yamlBlogLayoutConfig struct {
	ContentMaxWidth string `yaml:"content_max_width"`
	ShowToc         *bool  `yaml:"show_toc"`
	TocPosition     string `yaml:"toc_position"`
	TocWidth        string `yaml:"toc_width"`
	HeaderStyle     string `yaml:"header_style"`
	FooterStyle     string `yaml:"footer_style"`
	ShowAuthor      *bool  `yaml:"show_author"`
	ShowDate        *bool  `yaml:"show_date"`
	ShowTags        *bool  `yaml:"show_tags"`
	ShowReadingTime *bool  `yaml:"show_reading_time"`
	ShowPrevNext    *bool  `yaml:"show_prev_next"`
}

type yamlLandingLayoutConfig struct {
	ContentMaxWidth string `yaml:"content_max_width"`
	HeaderStyle     string `yaml:"header_style"`
	HeaderSticky    *bool  `yaml:"header_sticky"`
	FooterStyle     string `yaml:"footer_style"`
	HeroEnabled     *bool  `yaml:"hero_enabled"`
}

type yamlBareLayoutConfig struct {
	ContentMaxWidth string `yaml:"content_max_width"`
}

func (l *yamlLayoutConfig) toLayoutConfig() models.LayoutConfig {
	return models.LayoutConfig{
		Name:  l.Name,
		Paths: l.Paths,
		Feeds: l.Feeds,
		Docs: models.DocsLayoutConfig{
			SidebarPosition:    l.Docs.SidebarPosition,
			SidebarWidth:       l.Docs.SidebarWidth,
			SidebarCollapsible: l.Docs.SidebarCollapsible,
			SidebarDefaultOpen: l.Docs.SidebarDefaultOpen,
			TocPosition:        l.Docs.TocPosition,
			TocWidth:           l.Docs.TocWidth,
			TocCollapsible:     l.Docs.TocCollapsible,
			TocDefaultOpen:     l.Docs.TocDefaultOpen,
			ContentMaxWidth:    l.Docs.ContentMaxWidth,
			HeaderStyle:        l.Docs.HeaderStyle,
			FooterStyle:        l.Docs.FooterStyle,
		},
		Blog: models.BlogLayoutConfig{
			ContentMaxWidth: l.Blog.ContentMaxWidth,
			ShowToc:         l.Blog.ShowToc,
			TocPosition:     l.Blog.TocPosition,
			TocWidth:        l.Blog.TocWidth,
			HeaderStyle:     l.Blog.HeaderStyle,
			FooterStyle:     l.Blog.FooterStyle,
			ShowAuthor:      l.Blog.ShowAuthor,
			ShowDate:        l.Blog.ShowDate,
			ShowTags:        l.Blog.ShowTags,
			ShowReadingTime: l.Blog.ShowReadingTime,
			ShowPrevNext:    l.Blog.ShowPrevNext,
		},
		Landing: models.LandingLayoutConfig{
			ContentMaxWidth: l.Landing.ContentMaxWidth,
			HeaderStyle:     l.Landing.HeaderStyle,
			HeaderSticky:    l.Landing.HeaderSticky,
			FooterStyle:     l.Landing.FooterStyle,
			HeroEnabled:     l.Landing.HeroEnabled,
		},
		Bare: models.BareLayoutConfig{
			ContentMaxWidth: l.Bare.ContentMaxWidth,
		},
		Defaults: models.LayoutDefaults{
			ContentMaxWidth: l.Defaults.ContentMaxWidth,
			HeaderSticky:    l.Defaults.HeaderSticky,
			FooterSticky:    l.Defaults.FooterSticky,
		},
	}
}

// Sidebar-related YAML structs

type yamlSidebarConfig struct {
	Enabled      *bool                             `yaml:"enabled"`
	Position     string                            `yaml:"position"`
	Width        string                            `yaml:"width"`
	Collapsible  *bool                             `yaml:"collapsible"`
	DefaultOpen  *bool                             `yaml:"default_open"`
	Nav          []yamlSidebarNavItem              `yaml:"nav"`
	Title        string                            `yaml:"title"`
	Paths        map[string]*yamlPathSidebarConfig `yaml:"paths"`
	MultiFeed    *bool                             `yaml:"multi_feed"`
	Feeds        []string                          `yaml:"feeds"`
	FeedSections []yamlMultiFeedSection            `yaml:"feed_sections"`
	AutoGenerate *yamlSidebarAutoGenerate          `yaml:"auto_generate"`
}

type yamlSidebarNavItem struct {
	Title    string               `yaml:"title"`
	Href     string               `yaml:"href"`
	Children []yamlSidebarNavItem `yaml:"children"`
}

type yamlPathSidebarConfig struct {
	Title        string                   `yaml:"title"`
	AutoGenerate *yamlSidebarAutoGenerate `yaml:"auto_generate"`
	Items        []yamlSidebarNavItem     `yaml:"items"`
	Feed         string                   `yaml:"feed"`
	Position     string                   `yaml:"position"`
	Collapsible  *bool                    `yaml:"collapsible"`
}

type yamlSidebarAutoGenerate struct {
	Directory string   `yaml:"directory"`
	OrderBy   string   `yaml:"order_by"`
	Reverse   *bool    `yaml:"reverse"`
	MaxDepth  int      `yaml:"max_depth"`
	Exclude   []string `yaml:"exclude"`
}

type yamlMultiFeedSection struct {
	Feed      string `yaml:"feed"`
	Title     string `yaml:"title"`
	Collapsed *bool  `yaml:"collapsed"`
	MaxItems  int    `yaml:"max_items"`
}

func convertYamlSidebarNavItems(items []yamlSidebarNavItem) []models.SidebarNavItem {
	result := make([]models.SidebarNavItem, len(items))
	for i, item := range items {
		result[i] = models.SidebarNavItem{
			Title:    item.Title,
			Href:     item.Href,
			Children: convertYamlSidebarNavItems(item.Children),
		}
	}
	return result
}

func (s *yamlSidebarConfig) toSidebarConfig() models.SidebarConfig {
	config := models.SidebarConfig{
		Enabled:     s.Enabled,
		Position:    s.Position,
		Width:       s.Width,
		Collapsible: s.Collapsible,
		DefaultOpen: s.DefaultOpen,
		Nav:         convertYamlSidebarNavItems(s.Nav),
		Title:       s.Title,
		MultiFeed:   s.MultiFeed,
		Feeds:       s.Feeds,
	}

	// Convert paths
	if len(s.Paths) > 0 {
		config.Paths = make(map[string]*models.PathSidebarConfig)
		for path, pathConfig := range s.Paths {
			var autoGen *models.SidebarAutoGenerate
			if pathConfig.AutoGenerate != nil {
				autoGen = &models.SidebarAutoGenerate{
					Directory: pathConfig.AutoGenerate.Directory,
					OrderBy:   pathConfig.AutoGenerate.OrderBy,
					Reverse:   pathConfig.AutoGenerate.Reverse,
					MaxDepth:  pathConfig.AutoGenerate.MaxDepth,
					Exclude:   pathConfig.AutoGenerate.Exclude,
				}
			}
			config.Paths[path] = &models.PathSidebarConfig{
				Title:        pathConfig.Title,
				AutoGenerate: autoGen,
				Items:        convertYamlSidebarNavItems(pathConfig.Items),
				Feed:         pathConfig.Feed,
				Position:     pathConfig.Position,
				Collapsible:  pathConfig.Collapsible,
			}
		}
	}

	// Convert feed sections
	if len(s.FeedSections) > 0 {
		config.FeedSections = make([]models.MultiFeedSection, len(s.FeedSections))
		for i, section := range s.FeedSections {
			config.FeedSections[i] = models.MultiFeedSection{
				Feed:      section.Feed,
				Title:     section.Title,
				Collapsed: section.Collapsed,
				MaxItems:  section.MaxItems,
			}
		}
	}

	// Convert auto-generate
	if s.AutoGenerate != nil {
		config.AutoGenerate = &models.SidebarAutoGenerate{
			Directory: s.AutoGenerate.Directory,
			OrderBy:   s.AutoGenerate.OrderBy,
			Reverse:   s.AutoGenerate.Reverse,
			MaxDepth:  s.AutoGenerate.MaxDepth,
			Exclude:   s.AutoGenerate.Exclude,
		}
	}

	return config
}

// TOC-related YAML structs

type yamlTocConfig struct {
	Enabled     *bool  `yaml:"enabled"`
	Position    string `yaml:"position"`
	Width       string `yaml:"width"`
	MinDepth    int    `yaml:"min_depth"`
	MaxDepth    int    `yaml:"max_depth"`
	Title       string `yaml:"title"`
	Collapsible *bool  `yaml:"collapsible"`
	DefaultOpen *bool  `yaml:"default_open"`
	ScrollSpy   *bool  `yaml:"scroll_spy"`
}

func (t *yamlTocConfig) toTocConfig() models.TocConfig {
	return models.TocConfig{
		Enabled:     t.Enabled,
		Position:    t.Position,
		Width:       t.Width,
		MinDepth:    t.MinDepth,
		MaxDepth:    t.MaxDepth,
		Title:       t.Title,
		Collapsible: t.Collapsible,
		DefaultOpen: t.DefaultOpen,
		ScrollSpy:   t.ScrollSpy,
	}
}

// Header layout YAML structs

type yamlHeaderLayoutConfig struct {
	Style           string `yaml:"style"`
	Sticky          *bool  `yaml:"sticky"`
	ShowLogo        *bool  `yaml:"show_logo"`
	ShowTitle       *bool  `yaml:"show_title"`
	ShowNav         *bool  `yaml:"show_nav"`
	ShowSearch      *bool  `yaml:"show_search"`
	ShowThemeToggle *bool  `yaml:"show_theme_toggle"`
}

func (h *yamlHeaderLayoutConfig) toHeaderLayoutConfig() models.HeaderLayoutConfig {
	return models.HeaderLayoutConfig{
		Style:           h.Style,
		Sticky:          h.Sticky,
		ShowLogo:        h.ShowLogo,
		ShowTitle:       h.ShowTitle,
		ShowNav:         h.ShowNav,
		ShowSearch:      h.ShowSearch,
		ShowThemeToggle: h.ShowThemeToggle,
	}
}

// Blogroll-related YAML structs

type yamlBlogrollConfig struct {
	Enabled              bool                     `yaml:"enabled"`
	BlogrollSlug         string                   `yaml:"blogroll_slug"`
	ReaderSlug           string                   `yaml:"reader_slug"`
	CacheDir             string                   `yaml:"cache_dir"`
	CacheDuration        string                   `yaml:"cache_duration"`
	RefreshOnBuild       *bool                    `yaml:"refresh_on_build"`
	Timeout              int                      `yaml:"timeout"`
	ConcurrentRequests   int                      `yaml:"concurrent_requests"`
	MaxEntriesPerFeed    int                      `yaml:"max_entries_per_feed"`
	FallbackImageService string                   `yaml:"fallback_image_service"`
	Feeds                []yamlExternalFeedConfig `yaml:"feeds"`
	Templates            yamlBlogrollTemplates    `yaml:"templates"`
}

type yamlExternalFeedConfig struct {
	URL           string   `yaml:"url"`
	Title         string   `yaml:"title"`
	Description   string   `yaml:"description"`
	Category      string   `yaml:"category"`
	Type          string   `yaml:"type"`
	Tags          []string `yaml:"tags"`
	Active        *bool    `yaml:"active"`
	SiteURL       string   `yaml:"site_url"`
	ImageURL      string   `yaml:"image_url"`
	Handle        string   `yaml:"handle"`
	Aliases       []string `yaml:"aliases,omitempty"`
	MaxEntries    *int     `yaml:"max_entries,omitempty"`
	Primary       *bool    `yaml:"primary,omitempty"`
	PrimaryPerson string   `yaml:"primary_person"`
}

type yamlBlogrollTemplates struct {
	Blogroll string `yaml:"blogroll"`
	Reader   string `yaml:"reader"`
}

func (b *yamlBlogrollConfig) toBlogrollConfig() models.BlogrollConfig {
	// Get default values
	defaults := models.NewBlogrollConfig()

	// Apply defaults for template names if not specified
	blogrollTemplate := b.Templates.Blogroll
	if blogrollTemplate == "" {
		blogrollTemplate = defaults.Templates.Blogroll
	}
	readerTemplate := b.Templates.Reader
	if readerTemplate == "" {
		readerTemplate = defaults.Templates.Reader
	}

	config := models.BlogrollConfig{
		Enabled:              b.Enabled,
		BlogrollSlug:         b.BlogrollSlug,
		ReaderSlug:           b.ReaderSlug,
		CacheDir:             b.CacheDir,
		CacheDuration:        b.CacheDuration,
		RefreshOnBuild:       b.RefreshOnBuild,
		Timeout:              b.Timeout,
		ConcurrentRequests:   b.ConcurrentRequests,
		MaxEntriesPerFeed:    b.MaxEntriesPerFeed,
		FallbackImageService: b.FallbackImageService,
		Templates: models.BlogrollTemplates{
			Blogroll: blogrollTemplate,
			Reader:   readerTemplate,
		},
	}

	for i := range b.Feeds {
		fc := &b.Feeds[i]
		config.Feeds = append(config.Feeds, models.ExternalFeedConfig{
			URL:           fc.URL,
			Title:         fc.Title,
			Description:   fc.Description,
			Category:      fc.Category,
			Type:          models.ReaderFeedType(fc.Type),
			Tags:          fc.Tags,
			Active:        fc.Active,
			SiteURL:       fc.SiteURL,
			ImageURL:      fc.ImageURL,
			Handle:        fc.Handle,
			Aliases:       fc.Aliases,
			MaxEntries:    fc.MaxEntries,
			Primary:       fc.Primary,
			PrimaryPerson: fc.PrimaryPerson,
		})
	}

	return config
}

func (c *yamlComponentsConfig) toComponentsConfig() models.ComponentsConfig {
	config := models.ComponentsConfig{
		Nav: models.NavComponentConfig{
			Enabled:  c.Nav.Enabled,
			Position: c.Nav.Position,
			Style:    c.Nav.Style,
		},
		Footer: models.FooterComponentConfig{
			Enabled:       c.Footer.Enabled,
			Text:          c.Footer.Text,
			ShowCopyright: c.Footer.ShowCopyright,
		},
		DocSidebar: models.DocSidebarConfig{
			Enabled:  c.DocSidebar.Enabled,
			Position: c.DocSidebar.Position,
			Width:    c.DocSidebar.Width,
			MinDepth: c.DocSidebar.MinDepth,
			MaxDepth: c.DocSidebar.MaxDepth,
		},
		FeedSidebar: models.FeedSidebarConfig{
			Enabled:  c.FeedSidebar.Enabled,
			Position: c.FeedSidebar.Position,
			Width:    c.FeedSidebar.Width,
			Title:    c.FeedSidebar.Title,
			Feeds:    c.FeedSidebar.Feeds,
			MaxPosts: c.FeedSidebar.MaxPosts,
		},
		ContentSidebar: models.ContentSidebarConfig{
			Enabled:  c.ContentSidebar.Enabled,
			Position: c.ContentSidebar.Position,
			Width:    c.ContentSidebar.Width,
			Slug:     c.ContentSidebar.Slug,
		},
		CardRouter: models.CardRouterConfig{
			Mappings: c.CardRouter.Mappings,
		},
		Share: models.ShareComponentConfig{
			Enabled:   c.Share.Enabled,
			Platforms: append([]string{}, c.Share.Platforms...),
			Position:  c.Share.Position,
			Title:     c.Share.Title,
			Custom:    map[string]models.SharePlatformConfig{},
		},
		PostCopy: models.PostCopyComponentConfig{Enabled: c.PostCopy.Enabled},
		PostConnections: models.PostConnectionsComponentConfig{
			Enabled:       c.PostConn.Enabled,
			Display:       append([]string{}, c.PostConn.Display...),
			MinLinks:      c.PostConn.MinLinks,
			MaxLinks:      c.PostConn.MaxLinks,
			GraphMinLinks: c.PostConn.GraphMinLinks,
			GraphMaxLinks: c.PostConn.GraphMaxLinks,
			ListMinLinks:  c.PostConn.ListMinLinks,
			ListMaxLinks:  c.PostConn.ListMaxLinks,
			InlinksLimit:  c.PostConn.InlinksLimit,
			OutlinksLimit: c.PostConn.OutlinksLimit,
		},
	}

	if len(c.Share.Platforms) == 0 {
		config.Share.Platforms = nil
	}

	for key, custom := range c.Share.Custom {
		config.Share.Custom[key] = models.SharePlatformConfig{
			Name: custom.Name,
			Icon: custom.Icon,
			URL:  custom.URL,
		}
	}
	if len(c.Share.Custom) == 0 {
		config.Share.Custom = nil
	}

	// Convert nav items
	for _, item := range c.Nav.Items {
		config.Nav.Items = append(config.Nav.Items, models.NavItem{
			Label:    item.Label,
			URL:      item.URL,
			External: item.External,
		})
	}

	// Convert footer links
	for _, link := range c.Footer.Links {
		config.Footer.Links = append(config.Footer.Links, models.NavItem{
			Label:    link.Label,
			URL:      link.URL,
			External: link.External,
		})
	}

	return config
}

func (p *yamlPostFormatsConfig) toPostFormatsConfig() models.PostFormatsConfig {
	return models.PostFormatsConfig{
		HTML:     p.HTML,
		Markdown: p.Markdown,
		Text:     p.Text,
		ANSI:     p.ANSI,
		OG:       p.OG,
	}
}

func (w *yamlWellKnownConfig) toWellKnownConfig() models.WellKnownConfig {
	defaults := models.NewWellKnownConfig()
	config := models.WellKnownConfig{
		Enabled:         w.Enabled,
		AutoGenerate:    w.AutoGenerate,
		SSHFingerprint:  w.SSHFingerprint,
		KeybaseUsername: w.KeybaseUsername,
	}

	if config.Enabled == nil {
		config.Enabled = defaults.Enabled
	}
	if config.AutoGenerate == nil {
		config.AutoGenerate = defaults.AutoGenerate
	}

	return config
}

// configSource interface implementation for yamlConfig.
func (c *yamlConfig) getBaseConfig() baseConfigData {
	return baseConfigData{
		Fontpack:       c.Fontpack,
		FontpacksFile:  c.FontpacksFile,
		OutputDir:      c.OutputDir,
		URL:            c.URL,
		Title:          c.Title,
		Description:    c.Description,
		Author:         c.Author,
		Language:       c.Language,
		AuthorURL:      c.AuthorURL,
		ManagingEditor: c.ManagingEditor,
		WebMaster:      c.WebMaster,
		Copyright:      c.Copyright,
		License:        c.License,
		AssetsDir:      c.AssetsDir,
		Assets:         c.Assets,
		TemplatesDir:   c.TemplatesDir,
		Hooks:          c.Hooks,
		DisabledHooks:  c.DisabledHooks,
		GlobPatterns:   c.Glob.Patterns,
		UseGitignore:   c.Glob.UseGitignore,
		SlugMode:       c.Glob.SlugMode,
		SlugRules:      c.Glob.toSlugRules(),
		Extensions:     c.Markdown.Extensions,
		Concurrency:    c.Concurrency,
		Theme:          c.Theme.toThemeConfig(),
		ThemeCalendar:  c.ThemeCalendar,
		Head:           c.Head.toHeadConfig(),
		Presets:        c.Presets,
		Defaults:       c.Defaults,
		ErrorPages:     c.ErrorPages,
		ResourceHints:  c.ResourceHints,
		Highlight:      c.Markdown.Highlight,
		Footer:         c.Footer.toFooterConfig(),
	}
}

func (c *yamlConfig) getNavItems() []navItemData {
	items := make([]navItemData, len(c.Nav))
	for i, nav := range c.Nav {
		items[i] = navItemData(nav)
	}
	return items
}

func (c *yamlConfig) getFeeds() []feedConfigConverter {
	feeds := make([]feedConfigConverter, len(c.Feeds))
	for i := range c.Feeds {
		feeds[i] = &c.Feeds[i]
	}
	return feeds
}

func (c *yamlConfig) getFeedDefaults() feedDefaultsConverter       { return &c.FeedDefaults }
func (c *yamlConfig) getPostFormats() postFormatsConverter         { return &c.PostFormats }
func (c *yamlConfig) getWellKnown() wellKnownConverter             { return &c.WellKnown }
func (c *yamlConfig) getSEO() seoConverter                         { return &c.SEO }
func (c *yamlConfig) getIndieAuth() indieAuthConverter             { return &c.IndieAuth }
func (c *yamlConfig) getWebmention() webmentionConverter           { return &c.Webmention }
func (c *yamlConfig) getSearch() models.SearchConfig               { return c.Search }
func (c *yamlConfig) getComponents() componentsConverter           { return &c.Components }
func (c *yamlConfig) getLayout() layoutConverter                   { return &c.Layout }
func (c *yamlConfig) getSidebar() sidebarConverter                 { return &c.Sidebar }
func (c *yamlConfig) getToc() tocConverter                         { return &c.Toc }
func (c *yamlConfig) getHeader() headerConverter                   { return &c.Header }
func (c *yamlConfig) getBlogroll() blogrollConverter               { return &c.Blogroll }
func (c *yamlConfig) getTags() tagsConverter                       { return &c.Tags }
func (c *yamlConfig) getFeedsPage() feedsPageConverter             { return &c.FeedsPage }
func (c *yamlConfig) getEncryption() encryptionConverter           { return &c.Encryption }
func (c *yamlConfig) getTagAggregator() tagAggregatorConverter     { return &c.TagAggregator }
func (c *yamlConfig) getWebSub() webSubConverter                   { return &c.WebSub }
func (c *yamlConfig) getMentions() mentionsConverter               { return &c.Mentions }
func (c *yamlConfig) getShortcuts() shortcutsConverter             { return &c.Shortcuts }
func (c *yamlConfig) getViewTransitions() viewTransitionsConverter { return &c.ViewTransitions }
func (c *yamlConfig) getAuthors() authorsConverter                 { return &c.Authors }
func (c *yamlConfig) getGarden() gardenConverter                   { return &c.Garden }
func (c *yamlConfig) getTemplates() templatesConverter             { return &c.Templates }
func (c *yamlConfig) getBuilderAdmin() models.BuilderAdminConfig   { return c.BuilderAdmin }
func (c *yamlConfig) getImages() models.ImagesConfig               { return normalizeImagesConfig(c.Images) }

func (c *yamlConfig) toConfig() *models.Config {
	return buildConfig(c)
}

func (f *yamlFooterConfig) toFooterConfig() models.FooterConfig {
	return models.FooterConfig{
		Text:          f.Text,
		ShowCopyright: f.ShowCopyright,
	}
}

func (f *yamlFeedConfig) toFeedConfig() models.FeedConfig {
	return models.FeedConfig{
		Slug:            f.Slug,
		Presentation:    f.Presentation,
		Title:           f.Title,
		Description:     f.Description,
		Robots:          f.Robots,
		Filter:          f.Filter,
		Sort:            f.Sort,
		Reverse:         f.Reverse,
		Views:           f.Views,
		ShortsPath:      f.ShortsPath,
		Primary:         f.Primary,
		IncludePrivate:  f.IncludePrivate,
		Private:         f.Private,
		Sidebar:         f.Sidebar,
		ItemsPerPage:    f.ItemsPerPage,
		OrphanThreshold: f.OrphanThreshold,
		Limit:           f.Limit,
		Offset:          f.Offset,
		PaginationType:  models.PaginationType(f.PaginationType),
		ArchiveDisabled: f.ArchiveDisabled,
		Formats:         f.Formats.toFeedFormats(),
		Templates:       f.Templates.toFeedTemplates(),
	}
}

func (f *yamlFeedFormats) toFeedFormats() models.FeedFormats {
	formats := models.FeedFormats{}
	if f.HTML != nil {
		formats.HTML = *f.HTML
	}
	if f.SimpleHTML != nil {
		formats.SimpleHTML = *f.SimpleHTML
	}
	if f.RSS != nil {
		formats.RSS = *f.RSS
	}
	if f.Atom != nil {
		formats.Atom = *f.Atom
	}
	if f.JSON != nil {
		formats.JSON = *f.JSON
	}
	if f.Markdown != nil {
		formats.Markdown = *f.Markdown
	}
	if f.Text != nil {
		formats.Text = *f.Text
	}
	if f.Sitemap != nil {
		formats.Sitemap = *f.Sitemap
	}
	return formats
}

func (t *yamlFeedTemplates) toFeedTemplates() models.FeedTemplates {
	return models.FeedTemplates{
		HTML:       t.HTML,
		SimpleHTML: t.SimpleHTML,
		RSS:        t.RSS,
		Atom:       t.Atom,
		JSON:       t.JSON,
		Card:       t.Card,
		Sitemap:    t.Sitemap,
	}
}

func (d *yamlFeedDefaults) toFeedDefaults() models.FeedDefaults {
	return models.FeedDefaults{
		ItemsPerPage:    d.ItemsPerPage,
		Presentation:    d.Presentation,
		OrphanThreshold: d.OrphanThreshold,
		PaginationType:  models.PaginationType(d.PaginationType),
		Views:           d.Views,
		Formats:         d.Formats.toFeedFormats(),
		Templates:       d.Templates.toFeedTemplates(),
		Syndication: models.SyndicationConfig{
			MaxItems:             d.Syndication.MaxItems,
			IncludeContent:       d.Syndication.IncludeContent,
			SiteArchiveDisabled:  d.Syndication.SiteArchiveDisabled,
			FeedArchivesDisabled: d.Syndication.FeedArchivesDisabled,
		},
	}
}

// jsonConfig is an internal struct for parsing JSON configuration.
type jsonConfig struct {
	Fontpack        string                    `json:"fontpack"`
	FontpacksFile   string                    `json:"fontpacks_file"`
	OutputDir       string                    `json:"output_dir"`
	URL             string                    `json:"url"`
	Title           string                    `json:"title"`
	Description     string                    `json:"description"`
	Author          string                    `json:"author"`
	Language        string                    `json:"language"`
	AuthorURL       string                    `json:"author_url"`
	ManagingEditor  string                    `json:"managing_editor"`
	WebMaster       string                    `json:"webmaster"`
	Copyright       string                    `json:"copyright"`
	License         interface{}               `json:"license"`
	AssetsDir       string                    `json:"assets_dir"`
	Assets          models.AssetsConfig       `json:"assets"`
	TemplatesDir    string                    `json:"templates_dir"`
	Nav             []jsonNavItem             `json:"nav"`
	Footer          jsonFooterConfig          `json:"footer"`
	Hooks           []string                  `json:"hooks"`
	DisabledHooks   []string                  `json:"disabled_hooks"`
	Glob            jsonGlobConfig            `json:"glob"`
	Markdown        jsonMarkdownConfig        `json:"markdown"`
	Feeds           []jsonFeedConfig          `json:"feeds"`
	FeedDefaults    jsonFeedDefaults          `json:"feed_defaults"`
	Concurrency     int                       `json:"concurrency"`
	Theme           jsonThemeConfig           `json:"theme"`
	ThemeCalendar   parsedThemeCalendar       `json:"theme_calendar"`
	Head            parsedHead                `json:"head"`
	Presets         parsedTemplatePresets     `json:"template_presets"`
	Defaults        parsedDefaultTemplates    `json:"default_templates"`
	ErrorPages      parsedErrorPages          `json:"error_pages"`
	ResourceHints   parsedResourceHints       `json:"resource_hints"`
	PostFormats     jsonPostFormatsConfig     `json:"post_formats"`
	WellKnown       jsonWellKnownConfig       `json:"well_known"`
	IndieAuth       jsonIndieAuthConfig       `json:"indieauth"`
	Webmention      jsonWebmentionConfig      `json:"webmention"`
	Search          models.SearchConfig       `json:"search"`
	SEO             jsonSEOConfig             `json:"seo"`
	Components      jsonComponentsConfig      `json:"components"`
	Layout          jsonLayoutConfig          `json:"layout"`
	Sidebar         jsonSidebarConfig         `json:"sidebar"`
	Toc             jsonTocConfig             `json:"toc"`
	Header          jsonHeaderLayoutConfig    `json:"header"`
	Blogroll        jsonBlogrollConfig        `json:"blogroll"`
	Tags            jsonTagsConfig            `json:"tags"`
	FeedsPage       jsonFeedsPageConfig       `json:"feeds_page"`
	Encryption      jsonEncryptionConfig      `json:"encryption"`
	TagAggregator   jsonTagAggregatorConfig   `json:"tag_aggregator"`
	Mentions        jsonMentionsConfig        `json:"mentions"`
	WebSub          jsonWebSubConfig          `json:"websub"`
	Shortcuts       jsonShortcutsConfig       `json:"shortcuts"`
	ViewTransitions jsonViewTransitionsConfig `json:"view_transitions"`
	Authors         jsonAuthorsConfig         `json:"authors"`
	Garden          jsonGardenConfig          `json:"garden"`
	Templates       jsonTemplatesConfig       `json:"templates"`
	Images          models.ImagesConfig       `json:"images"`
	BuilderAdmin    models.BuilderAdminConfig `json:"builder_admin"`
}

type jsonNavItem struct {
	Label    string `json:"label"`
	URL      string `json:"url"`
	External bool   `json:"external"`
}

type jsonTemplatesConfig struct {
	Media jsonTemplatesMediaConfig `json:"media"`
}

type jsonTemplatesMediaConfig struct {
	TrustedDomains []string `json:"trusted_domains"`
}

func (t *jsonTemplatesConfig) toTemplatesConfig() models.TemplatesConfig {
	cfg := models.NewTemplatesConfig()
	if t.Media.TrustedDomains != nil {
		cfg.Media.TrustedDomains = append([]string{}, t.Media.TrustedDomains...)
	}
	return cfg
}

type jsonFooterConfig struct {
	Text          string `json:"text"`
	ShowCopyright *bool  `json:"show_copyright"`
}

type jsonGlobConfig struct {
	Patterns     []string       `json:"patterns"`
	UseGitignore *bool          `json:"use_gitignore"`
	SlugMode     string         `json:"slug_mode"`
	SlugRules    []jsonSlugRule `json:"slug_rules"`
}

type jsonSlugRule struct {
	Prefix string `json:"prefix"`
	Mode   string `json:"mode"`
}

type jsonMarkdownConfig struct {
	Extensions []string        `json:"extensions"`
	Highlight  parsedHighlight `json:"highlight"`
}

func (g tomlGlobConfig) toSlugRules() []models.SlugRule {
	rules := make([]models.SlugRule, 0, len(g.SlugRules))
	for _, rule := range g.SlugRules {
		rules = append(rules, models.SlugRule{Prefix: rule.Prefix, Mode: rule.Mode})
	}
	return rules
}

func (g yamlGlobConfig) toSlugRules() []models.SlugRule {
	rules := make([]models.SlugRule, 0, len(g.SlugRules))
	for _, rule := range g.SlugRules {
		rules = append(rules, models.SlugRule{Prefix: rule.Prefix, Mode: rule.Mode})
	}
	return rules
}

func (g jsonGlobConfig) toSlugRules() []models.SlugRule {
	rules := make([]models.SlugRule, 0, len(g.SlugRules))
	for _, rule := range g.SlugRules {
		rules = append(rules, models.SlugRule{Prefix: rule.Prefix, Mode: rule.Mode})
	}
	return rules
}

type jsonFeedConfig struct {
	Slug            string            `json:"slug"`
	Presentation    string            `json:"presentation"`
	Title           string            `json:"title"`
	Description     string            `json:"description"`
	Robots          string            `json:"robots"`
	Filter          string            `json:"filter"`
	Sort            string            `json:"sort"`
	Reverse         bool              `json:"reverse"`
	Views           []string          `json:"views"`
	ShortsPath      string            `json:"shorts_path"`
	Primary         bool              `json:"primary"`
	IncludePrivate  bool              `json:"include_private"`
	Private         bool              `json:"private"`
	Sidebar         *bool             `json:"sidebar"`
	ItemsPerPage    int               `json:"items_per_page"`
	OrphanThreshold int               `json:"orphan_threshold"`
	Limit           int               `json:"limit"`
	Offset          int               `json:"offset"`
	PaginationType  string            `json:"pagination_type"`
	ArchiveDisabled bool              `json:"archive_disabled"`
	Formats         jsonFeedFormats   `json:"formats"`
	Templates       jsonFeedTemplates `json:"templates"`
}

type jsonFeedFormats struct {
	HTML       *bool `json:"html"`
	SimpleHTML *bool `json:"simple_html"`
	RSS        *bool `json:"rss"`
	Atom       *bool `json:"atom"`
	JSON       *bool `json:"json"`
	Markdown   *bool `json:"markdown"`
	Text       *bool `json:"text"`
	Sitemap    *bool `json:"sitemap"`
}

type jsonFeedTemplates struct {
	HTML       string `json:"html"`
	SimpleHTML string `json:"simple_html"`
	RSS        string `json:"rss"`
	Atom       string `json:"atom"`
	JSON       string `json:"json"`
	Card       string `json:"card"`
	Sitemap    string `json:"sitemap"`
}

type jsonFeedDefaults struct {
	ItemsPerPage    int                   `json:"items_per_page"`
	Presentation    string                `json:"presentation"`
	OrphanThreshold int                   `json:"orphan_threshold"`
	PaginationType  string                `json:"pagination_type"`
	Views           []string              `json:"views"`
	Formats         jsonFeedFormats       `json:"formats"`
	Templates       jsonFeedTemplates     `json:"templates"`
	Syndication     jsonSyndicationConfig `json:"syndication"`
}

type jsonSyndicationConfig struct {
	MaxItems             int  `json:"max_items"`
	IncludeContent       bool `json:"include_content"`
	SiteArchiveDisabled  bool `json:"site_archive_disabled"`
	FeedArchivesDisabled bool `json:"feed_archives_disabled"`
}

type jsonPostFormatsConfig struct {
	HTML     *bool `json:"html"`
	Markdown bool  `json:"markdown"`
	Text     bool  `json:"text"`
	ANSI     bool  `json:"ansi"`
	OG       bool  `json:"og"`
}

type jsonWellKnownConfig struct {
	Enabled         *bool    `json:"enabled"`
	AutoGenerate    []string `json:"auto_generate"`
	SSHFingerprint  string   `json:"ssh_fingerprint"`
	KeybaseUsername string   `json:"keybase_username"`
}

type jsonSEOConfig struct {
	TwitterHandle  string `json:"twitter_handle"`
	DefaultImage   string `json:"default_image"`
	LogoURL        string `json:"logo_url"`
	AuthorImage    string `json:"author_image"`
	OGImageService string `json:"og_image_service"`
}

type jsonIndieAuthConfig struct {
	Enabled               bool   `json:"enabled"`
	AuthorizationEndpoint string `json:"authorization_endpoint"`
	TokenEndpoint         string `json:"token_endpoint"`
	MeURL                 string `json:"me_url"`
}

type jsonWebmentionConfig struct {
	Enabled  bool   `json:"enabled"`
	Endpoint string `json:"endpoint"`
}

type jsonTagsConfig struct {
	Enabled     *bool    `json:"enabled"`
	Blacklist   []string `json:"blacklist"`
	Private     []string `json:"private"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Template    string   `json:"template"`
	SlugPrefix  string   `json:"slug_prefix"`
}

type jsonFeedsPageConfig struct {
	Enabled          *bool    `json:"enabled"`
	Title            string   `json:"title"`
	Description      string   `json:"description"`
	Template         string   `json:"template"`
	SlugPrefix       string   `json:"slug_prefix"`
	Robots           string   `json:"robots"`
	ShowPrivateFeeds []string `json:"show_private_feeds"`
}

func (t *jsonTagsConfig) toTagsConfig() models.TagsConfig {
	defaults := models.NewTagsConfig()

	config := models.TagsConfig{
		Enabled:     t.Enabled,
		Blacklist:   t.Blacklist,
		Private:     t.Private,
		Title:       t.Title,
		Description: t.Description,
		Template:    t.Template,
		SlugPrefix:  t.SlugPrefix,
	}

	// Apply defaults if not set
	if config.Enabled == nil {
		config.Enabled = defaults.Enabled
	}
	if config.Title == "" {
		config.Title = defaults.Title
	}
	if config.Template == "" {
		config.Template = defaults.Template
	}
	if config.SlugPrefix == "" {
		config.SlugPrefix = defaults.SlugPrefix
	}

	return config
}

func (t *jsonFeedsPageConfig) toFeedsPageConfig() models.FeedsPageConfig {
	defaults := models.NewFeedsPageConfig()

	config := models.FeedsPageConfig{
		Enabled:          t.Enabled,
		Title:            t.Title,
		Description:      t.Description,
		Template:         t.Template,
		SlugPrefix:       t.SlugPrefix,
		Robots:           t.Robots,
		ShowPrivateFeeds: append([]string{}, t.ShowPrivateFeeds...),
	}

	if config.Enabled == nil {
		config.Enabled = defaults.Enabled
	}
	if config.Title == "" {
		config.Title = defaults.Title
	}
	if config.Description == "" {
		config.Description = defaults.Description
	}
	if config.Template == "" {
		config.Template = defaults.Template
	}
	if config.SlugPrefix == "" {
		config.SlugPrefix = defaults.SlugPrefix
	}

	return config
}

type jsonEncryptionConfig struct {
	Enabled               *bool             `json:"enabled"`
	DefaultKey            string            `json:"default_key"`
	DecryptionHint        string            `json:"decryption_hint"`
	PrivateTags           map[string]string `json:"private_tags"`
	EnforceStrength       *bool             `json:"enforce_strength"`
	MinEstimatedCrackTime string            `json:"min_estimated_crack_time"`
	MinPasswordLength     int               `json:"min_password_length"`
}

func (e *jsonEncryptionConfig) toEncryptionConfig() models.EncryptionConfig {
	defaults := models.NewEncryptionConfig()

	config := models.EncryptionConfig{
		DefaultKey:            e.DefaultKey,
		DecryptionHint:        e.DecryptionHint,
		PrivateTags:           e.PrivateTags,
		MinEstimatedCrackTime: e.MinEstimatedCrackTime,
		MinPasswordLength:     e.MinPasswordLength,
	}

	// Apply defaults for unset values
	if e.Enabled != nil {
		config.Enabled = *e.Enabled
	} else {
		config.Enabled = defaults.Enabled
	}
	if config.DefaultKey == "" {
		config.DefaultKey = defaults.DefaultKey
	}
	if e.EnforceStrength != nil {
		config.EnforceStrength = *e.EnforceStrength
	} else {
		config.EnforceStrength = defaults.EnforceStrength
	}
	if config.MinEstimatedCrackTime == "" {
		config.MinEstimatedCrackTime = defaults.MinEstimatedCrackTime
	}
	if config.MinPasswordLength == 0 {
		config.MinPasswordLength = defaults.MinPasswordLength // pragma: allowlist secret
	}

	return config
}

type jsonTagAggregatorConfig struct {
	Enabled        *bool               `json:"enabled"`
	Synonyms       map[string][]string `json:"synonyms"`
	Additional     map[string][]string `json:"additional"`
	GenerateReport bool                `json:"generate_report"`
}

func (t *jsonTagAggregatorConfig) toTagAggregatorConfig() models.TagAggregatorConfig {
	defaults := models.NewTagAggregatorConfig()

	config := models.TagAggregatorConfig{
		Enabled:        t.Enabled,
		Synonyms:       t.Synonyms,
		Additional:     t.Additional,
		GenerateReport: t.GenerateReport,
	}

	// Apply defaults if not set
	if config.Enabled == nil {
		config.Enabled = defaults.Enabled
	}

	return config
}

// Mentions-related JSON structs

type jsonMentionsConfig struct {
	Enabled            *bool                   `json:"enabled"`
	CSSClass           string                  `json:"css_class"`
	FromPosts          []jsonMentionPostSource `json:"from_posts"`
	CacheDir           string                  `json:"cache_dir"`
	CacheDuration      string                  `json:"cache_duration"`
	Timeout            int                     `json:"timeout"`
	ConcurrentRequests int                     `json:"concurrent_requests"`
}

type jsonMentionPostSource struct {
	Filter       string `json:"filter"`
	HandleField  string `json:"handle_field"`
	AliasesField string `json:"aliases_field"`
	AvatarField  string `json:"avatar_field"`
}

func (m *jsonMentionsConfig) toMentionsConfig() models.MentionsConfig {
	defaults := models.NewMentionsConfig()

	config := models.MentionsConfig{
		Enabled:            m.Enabled,
		CSSClass:           m.CSSClass,
		CacheDir:           m.CacheDir,
		CacheDuration:      m.CacheDuration,
		Timeout:            m.Timeout,
		ConcurrentRequests: m.ConcurrentRequests,
	}

	// Apply defaults for unset values
	if config.Enabled == nil {
		config.Enabled = defaults.Enabled
	}
	if config.CSSClass == "" {
		config.CSSClass = defaults.CSSClass
	}
	if config.CacheDir == "" {
		config.CacheDir = defaults.CacheDir
	}
	if config.CacheDuration == "" {
		config.CacheDuration = defaults.CacheDuration
	}
	if config.Timeout == 0 {
		config.Timeout = defaults.Timeout
	}
	if config.ConcurrentRequests == 0 {
		config.ConcurrentRequests = defaults.ConcurrentRequests
	}

	// Convert from_posts sources
	for _, src := range m.FromPosts {
		aliasesField := src.AliasesField
		if aliasesField == "" {
			aliasesField = defaultAliasesField
		}
		config.FromPosts = append(config.FromPosts, models.MentionPostSource{
			Filter:       src.Filter,
			HandleField:  src.HandleField,
			AliasesField: aliasesField,
			AvatarField:  src.AvatarField,
		})
	}

	// Fall back to defaults when user has mentions config but no from_posts
	if len(config.FromPosts) == 0 {
		config.FromPosts = defaults.FromPosts
	}

	return config
}

type jsonWebSubConfig struct {
	Enabled *bool    `json:"enabled"`
	Hubs    []string `json:"hubs"`
}

func (w *jsonWebSubConfig) toWebSubConfig() models.WebSubConfig {
	defaults := models.NewWebSubConfig()
	config := models.WebSubConfig{
		Enabled: w.Enabled,
		Hubs:    w.Hubs,
	}
	if config.Enabled == nil {
		config.Enabled = defaults.Enabled
	}
	if config.Hubs == nil {
		config.Hubs = defaults.Hubs
	}
	return config
}

type jsonShortcutsConfig struct {
	Navigation map[string]string `json:"navigation"`
}

type jsonViewTransitionsConfig struct {
	Enabled       *bool    `json:"enabled"`
	Native        *bool    `json:"native"`
	Debug         *bool    `json:"debug"`
	Duration      int      `json:"duration"`
	UpdateMeta    *bool    `json:"update_meta"`
	ScrollToTop   *bool    `json:"scroll_to_top"`
	SkipClasses   []string `json:"skip_classes"`
	SkipSelectors []string `json:"skip_selectors"`
}

func (s *jsonShortcutsConfig) toShortcutsConfig() models.ShortcutsConfig {
	config := models.ShortcutsConfig{
		Navigation: s.Navigation,
	}
	if config.Navigation == nil {
		config.Navigation = make(map[string]string)
	}
	return config
}

func (v *jsonViewTransitionsConfig) toViewTransitionsConfig() models.ViewTransitionsConfig {
	defaults := models.NewViewTransitionsConfig()
	config := models.ViewTransitionsConfig{
		Enabled:       v.Enabled,
		Native:        v.Native,
		Debug:         v.Debug,
		Duration:      v.Duration,
		UpdateMeta:    v.UpdateMeta,
		ScrollToTop:   v.ScrollToTop,
		SkipClasses:   v.SkipClasses,
		SkipSelectors: v.SkipSelectors,
	}
	if config.Enabled == nil {
		config.Enabled = defaults.Enabled
	}
	if config.Debug == nil {
		config.Debug = defaults.Debug
	}
	if config.Duration == 0 {
		config.Duration = defaults.Duration
	}
	if config.UpdateMeta == nil {
		config.UpdateMeta = defaults.UpdateMeta
	}
	if config.ScrollToTop == nil {
		config.ScrollToTop = defaults.ScrollToTop
	}
	if config.SkipClasses == nil {
		config.SkipClasses = defaults.SkipClasses
	}
	if config.SkipSelectors == nil {
		config.SkipSelectors = defaults.SkipSelectors
	}
	return config
}

type jsonAuthorsConfig struct {
	GeneratePages bool                  `json:"generate_pages"`
	URLPattern    string                `json:"url_pattern"`
	FeedsEnabled  bool                  `json:"feeds_enabled"`
	Authors       map[string]jsonAuthor `json:"authors"`
}

type jsonAuthor struct {
	Name          string            `json:"name"`
	Bio           string            `json:"bio"`
	Email         string            `json:"email"`
	Avatar        string            `json:"avatar"`
	URL           string            `json:"url"`
	Social        map[string]string `json:"social"`
	Guest         bool              `json:"guest"`
	Active        bool              `json:"active"`
	Default       bool              `json:"default"`
	Contributions []string          `json:"contributions"`
	Role          string            `json:"role"`
	Contribution  string            `json:"contribution"`
}

func (a *jsonAuthorsConfig) toAuthorsConfig() models.AuthorsConfig {
	config := models.AuthorsConfig{
		GeneratePages: a.GeneratePages,
		URLPattern:    a.URLPattern,
		FeedsEnabled:  a.FeedsEnabled,
	}
	if a.Authors != nil {
		config.Authors = make(map[string]models.Author, len(a.Authors))
		for id := range a.Authors {
			author := a.Authors[id]
			ma := models.Author{
				ID:      id,
				Name:    author.Name,
				Guest:   author.Guest,
				Active:  author.Active,
				Default: author.Default,
				Social:  author.Social,
			}
			if author.Bio != "" {
				bio := author.Bio
				ma.Bio = &bio
			}
			if author.Email != "" {
				email := author.Email
				ma.Email = &email
			}
			if author.Avatar != "" {
				avatar := author.Avatar
				ma.Avatar = &avatar
			}
			if author.URL != "" {
				url := author.URL
				ma.URL = &url
			}
			if author.Role != "" {
				role := author.Role
				ma.Role = &role
			}
			if author.Contribution != "" {
				contribution := author.Contribution
				ma.Contribution = &contribution
			}
			if len(author.Contributions) > 0 {
				ma.Contributions = author.Contributions
			}
			config.Authors[id] = ma
		}
	}
	return config
}

type jsonGardenConfig struct {
	Enabled      *bool    `json:"enabled"`
	Path         string   `json:"path"`
	ExportJSON   *bool    `json:"export_json"`
	RenderPage   *bool    `json:"render_page"`
	IncludeTags  *bool    `json:"include_tags"`
	IncludePosts *bool    `json:"include_posts"`
	MaxNodes     int      `json:"max_nodes"`
	ExcludeTags  []string `json:"exclude_tags"`
	Template     string   `json:"template"`
	Title        string   `json:"title"`
	Description  string   `json:"description"`
}

func (g *jsonGardenConfig) toGardenConfig() models.GardenConfig {
	defaults := models.NewGardenConfig()

	config := models.GardenConfig{
		Enabled:      g.Enabled,
		Path:         g.Path,
		ExportJSON:   g.ExportJSON,
		RenderPage:   g.RenderPage,
		IncludeTags:  g.IncludeTags,
		IncludePosts: g.IncludePosts,
		MaxNodes:     g.MaxNodes,
		ExcludeTags:  g.ExcludeTags,
		Template:     g.Template,
		Title:        g.Title,
		Description:  g.Description,
	}

	// Apply defaults for unset fields
	if config.Enabled == nil {
		config.Enabled = defaults.Enabled
	}
	if config.ExportJSON == nil {
		config.ExportJSON = defaults.ExportJSON
	}
	if config.RenderPage == nil {
		config.RenderPage = defaults.RenderPage
	}
	if config.IncludeTags == nil {
		config.IncludeTags = defaults.IncludeTags
	}
	if config.IncludePosts == nil {
		config.IncludePosts = defaults.IncludePosts
	}
	if config.Path == "" {
		config.Path = defaults.Path
	}
	if config.Template == "" {
		config.Template = defaults.Template
	}
	if config.Title == "" {
		config.Title = defaults.Title
	}

	return config
}

type jsonThemeConfig struct {
	ContractVersion     int                      `json:"contract_version"`
	Fontpack            string                   `json:"fontpack"`
	Name                string                   `json:"name"`
	Aesthetic           string                   `json:"aesthetic"`
	Palette             string                   `json:"palette"`
	PaletteLight        string                   `json:"palette_light"`
	PaletteDark         string                   `json:"palette_dark"`
	FallbackMode        string                   `json:"fallback_mode"`
	TextSize            string                   `json:"text_size"`
	Seasonal            bool                     `json:"seasonal"`
	ShowTextSizeControl *bool                    `json:"show_text_size_control"`
	SeedColor           string                   `json:"seed_color"`
	Variables           map[string]string        `json:"variables"`
	CustomCSS           string                   `json:"custom_css"`
	Background          jsonBackgroundConfig     `json:"background"`
	Font                jsonFontConfig           `json:"font"`
	Switcher            jsonThemeSwitcherConfig  `json:"switcher"`
	Texture             jsonTextureConfig        `json:"texture"`
	HeadingTexture      jsonHeadingTextureConfig `json:"heading_texture"`
	Motif               jsonMotifConfig          `json:"motif"`
}

type jsonTextureConfig struct {
	Kind     string  `json:"kind"`
	ColorMix float64 `json:"color_mix"`
	Scale    float64 `json:"scale"`
	Scope    string  `json:"scope"`
}
type jsonHeadingTextureConfig struct {
	Kind     string  `json:"kind"`
	ColorMix float64 `json:"color_mix"`
	Scale    float64 `json:"scale"`
}
type jsonMotifConfig struct {
	Kind      string  `json:"kind"`
	Glyph     string  `json:"glyph"`
	Size      string  `json:"size"`
	Gap       string  `json:"gap"`
	RowOffset float64 `json:"row_offset"`
	Wobble    float64 `json:"wobble"`
	Scatter   float64 `json:"scatter"`
	Layer     string  `json:"layer"`
	Color     string  `json:"color"`
	ColorMix  float64 `json:"color_mix"`
	URL       string  `json:"url"`
}

type jsonThemeSwitcherConfig struct {
	Enabled    *bool    `json:"enabled"`
	ModeToggle *bool    `json:"mode_toggle"`
	IncludeAll *bool    `json:"include_all"`
	Include    []string `json:"include"`
	Exclude    []string `json:"exclude"`
	Position   string   `json:"position"`
}

type jsonBackgroundConfig struct {
	Enabled            *bool                   `json:"enabled"`
	Backgrounds        []jsonBackgroundElement `json:"backgrounds"`
	Scripts            []string                `json:"scripts"`
	CSS                string                  `json:"css"`
	ArticleBg          string                  `json:"article_bg"`
	ArticleBlurEnabled *bool                   `json:"article_blur_enabled"`
	ArticleBlur        string                  `json:"article_blur"`
	ArticleShadow      string                  `json:"article_shadow"`
	ArticleBorder      string                  `json:"article_border"`
	ArticleRadius      string                  `json:"article_radius"`
}

type jsonBackgroundElement struct {
	HTML   string `json:"html"`
	ZIndex int    `json:"z_index"`
}

type jsonFontConfig struct {
	Family        string   `json:"family"`
	HeadingFamily string   `json:"heading_family"`
	CodeFamily    string   `json:"code_family"`
	Size          string   `json:"size"`
	LineHeight    string   `json:"line_height"`
	GoogleFonts   []string `json:"google_fonts"`
	CustomURLs    []string `json:"custom_urls"`
}

func (t *jsonThemeConfig) toThemeConfig() models.ThemeConfig {
	variables := t.Variables
	if variables == nil {
		variables = make(map[string]string)
	}
	return models.ThemeConfig{
		ContractVersion:     t.ContractVersion,
		Fontpack:            t.Fontpack,
		Name:                t.Name,
		Aesthetic:           t.Aesthetic,
		Palette:             t.Palette,
		PaletteLight:        t.PaletteLight,
		PaletteDark:         t.PaletteDark,
		FallbackMode:        t.FallbackMode,
		TextSize:            t.TextSize,
		Seasonal:            t.Seasonal,
		ShowTextSizeControl: t.ShowTextSizeControl,
		SeedColor:           t.SeedColor,
		Variables:           variables,
		CustomCSS:           t.CustomCSS,
		Background:          t.Background.toBackgroundConfig(),
		Font:                t.Font.toFontConfig(),
		Switcher:            t.Switcher.toThemeSwitcherConfig(),
		Texture:             models.ThemeTextureConfig{Kind: t.Texture.Kind, ColorMix: t.Texture.ColorMix, Scale: t.Texture.Scale, Scope: t.Texture.Scope},
		HeadingTexture:      models.ThemeHeadingTextureConfig{Kind: t.HeadingTexture.Kind, ColorMix: t.HeadingTexture.ColorMix, Scale: t.HeadingTexture.Scale},
		Motif:               models.ThemeMotifConfig{Kind: t.Motif.Kind, Glyph: t.Motif.Glyph, Size: t.Motif.Size, Gap: t.Motif.Gap, RowOffset: t.Motif.RowOffset, Wobble: t.Motif.Wobble, Scatter: t.Motif.Scatter, Layer: t.Motif.Layer, Color: t.Motif.Color, ColorMix: t.Motif.ColorMix, URL: t.Motif.URL},
	}
}

func (s *jsonThemeSwitcherConfig) toThemeSwitcherConfig() models.ThemeSwitcherConfig {
	return models.ThemeSwitcherConfig{
		Enabled:    s.Enabled,
		ModeToggle: s.ModeToggle,
		IncludeAll: s.IncludeAll,
		Include:    s.Include,
		Exclude:    s.Exclude,
		Position:   s.Position,
	}
}

func (b *jsonBackgroundConfig) toBackgroundConfig() models.BackgroundConfig {
	backgrounds := make([]models.BackgroundElement, len(b.Backgrounds))
	for i, bg := range b.Backgrounds {
		backgrounds[i] = models.BackgroundElement{
			HTML:   bg.HTML,
			ZIndex: bg.ZIndex,
		}
	}
	return models.BackgroundConfig{
		Enabled:            b.Enabled,
		Backgrounds:        backgrounds,
		Scripts:            b.Scripts,
		CSS:                b.CSS,
		ArticleBg:          b.ArticleBg,
		ArticleBlurEnabled: b.ArticleBlurEnabled,
		ArticleBlur:        b.ArticleBlur,
		ArticleShadow:      b.ArticleShadow,
		ArticleBorder:      b.ArticleBorder,
		ArticleRadius:      b.ArticleRadius,
	}
}

func (f *jsonFontConfig) toFontConfig() models.FontConfig {
	return models.FontConfig{
		Family:        f.Family,
		HeadingFamily: f.HeadingFamily,
		CodeFamily:    f.CodeFamily,
		Size:          f.Size,
		LineHeight:    f.LineHeight,
		GoogleFonts:   f.GoogleFonts,
		CustomURLs:    f.CustomURLs,
	}
}

func (s *jsonSEOConfig) toSEOConfig() models.SEOConfig {
	return models.SEOConfig{
		TwitterHandle:  s.TwitterHandle,
		DefaultImage:   s.DefaultImage,
		LogoURL:        s.LogoURL,
		AuthorImage:    s.AuthorImage,
		OGImageService: s.OGImageService,
	}
}

func (i *jsonIndieAuthConfig) toIndieAuthConfig() models.IndieAuthConfig {
	return models.IndieAuthConfig{
		Enabled:               i.Enabled,
		AuthorizationEndpoint: i.AuthorizationEndpoint,
		TokenEndpoint:         i.TokenEndpoint,
		MeURL:                 i.MeURL,
	}
}

func (w *jsonWebmentionConfig) toWebmentionConfig() models.WebmentionConfig {
	return models.WebmentionConfig{
		Enabled:  w.Enabled,
		Endpoint: w.Endpoint,
	}
}

type jsonComponentsConfig struct {
	Nav            jsonNavComponentConfig      `json:"nav"`
	Footer         jsonFooterComponentConfig   `json:"footer"`
	DocSidebar     jsonDocSidebarConfig        `json:"doc_sidebar"`
	FeedSidebar    jsonFeedSidebarConfig       `json:"feed_sidebar"`
	ContentSidebar jsonContentSidebarConfig    `json:"content_sidebar"`
	CardRouter     jsonCardRouterConfig        `json:"card_router"`
	Share          jsonShareComponentConfig    `json:"share"`
	PostCopy       jsonPostCopyComponentConfig `json:"post_copy"`
	PostConn       jsonPostConnectionsConfig   `json:"post_connections"`
}

type jsonPostCopyComponentConfig struct {
	Enabled *bool `json:"enabled"`
}

type jsonPostConnectionsConfig struct {
	Enabled       *bool    `json:"enabled"`
	Display       []string `json:"display"`
	MinLinks      int      `json:"min_links"`
	MaxLinks      int      `json:"max_links"`
	GraphMinLinks int      `json:"graph_min_links"`
	GraphMaxLinks int      `json:"graph_max_links"`
	ListMinLinks  int      `json:"list_min_links"`
	ListMaxLinks  int      `json:"list_max_links"`
	InlinksLimit  int      `json:"inlinks_limit"`
	OutlinksLimit int      `json:"outlinks_limit"`
}

type jsonNavComponentConfig struct {
	Enabled  *bool         `json:"enabled"`
	Position string        `json:"position"`
	Style    string        `json:"style"`
	Items    []jsonNavItem `json:"items"`
}

type jsonFooterComponentConfig struct {
	Enabled       *bool         `json:"enabled"`
	Text          string        `json:"text"`
	ShowCopyright *bool         `json:"show_copyright"`
	Links         []jsonNavItem `json:"links"`
}

type jsonDocSidebarConfig struct {
	Enabled  *bool  `json:"enabled"`
	Position string `json:"position"`
	Width    string `json:"width"`
	MinDepth int    `json:"min_depth"`
	MaxDepth int    `json:"max_depth"`
}

type jsonFeedSidebarConfig struct {
	Enabled  *bool    `json:"enabled"`
	Position string   `json:"position"`
	Width    string   `json:"width"`
	Title    string   `json:"title"`
	Feeds    []string `json:"feeds"`
	MaxPosts int      `json:"max_posts"`
}

type jsonContentSidebarConfig struct {
	Enabled  *bool  `json:"enabled"`
	Position string `json:"position"`
	Width    string `json:"width"`
	Slug     string `json:"slug"`
}

type jsonCardRouterConfig struct {
	Mappings map[string]string `json:"mappings"`
}

type jsonShareComponentConfig struct {
	Enabled   *bool                            `json:"enabled"`
	Platforms []string                         `json:"platforms"`
	Position  string                           `json:"position"`
	Title     string                           `json:"title"`
	Custom    map[string]jsonSharePlatformItem `json:"custom"`
}

type jsonSharePlatformItem struct {
	Name string `json:"name"`
	Icon string `json:"icon"`
	URL  string `json:"url"`
}

// Layout-related JSON structs

type jsonLayoutConfig struct {
	Name     string                  `json:"name"`
	Paths    map[string]string       `json:"paths"`
	Feeds    map[string]string       `json:"feeds"`
	Docs     jsonDocsLayoutConfig    `json:"docs"`
	Blog     jsonBlogLayoutConfig    `json:"blog"`
	Landing  jsonLandingLayoutConfig `json:"landing"`
	Bare     jsonBareLayoutConfig    `json:"bare"`
	Defaults jsonLayoutDefaults      `json:"defaults"`
}

type jsonLayoutDefaults struct {
	ContentMaxWidth string `json:"content_max_width"`
	HeaderSticky    *bool  `json:"header_sticky"`
	FooterSticky    *bool  `json:"footer_sticky"`
}

type jsonDocsLayoutConfig struct {
	SidebarPosition    string `json:"sidebar_position"`
	SidebarWidth       string `json:"sidebar_width"`
	SidebarCollapsible *bool  `json:"sidebar_collapsible"`
	SidebarDefaultOpen *bool  `json:"sidebar_default_open"`
	TocPosition        string `json:"toc_position"`
	TocWidth           string `json:"toc_width"`
	TocCollapsible     *bool  `json:"toc_collapsible"`
	TocDefaultOpen     *bool  `json:"toc_default_open"`
	ContentMaxWidth    string `json:"content_max_width"`
	HeaderStyle        string `json:"header_style"`
	FooterStyle        string `json:"footer_style"`
}

type jsonBlogLayoutConfig struct {
	ContentMaxWidth string `json:"content_max_width"`
	ShowToc         *bool  `json:"show_toc"`
	TocPosition     string `json:"toc_position"`
	TocWidth        string `json:"toc_width"`
	HeaderStyle     string `json:"header_style"`
	FooterStyle     string `json:"footer_style"`
	ShowAuthor      *bool  `json:"show_author"`
	ShowDate        *bool  `json:"show_date"`
	ShowTags        *bool  `json:"show_tags"`
	ShowReadingTime *bool  `json:"show_reading_time"`
	ShowPrevNext    *bool  `json:"show_prev_next"`
}

type jsonLandingLayoutConfig struct {
	ContentMaxWidth string `json:"content_max_width"`
	HeaderStyle     string `json:"header_style"`
	HeaderSticky    *bool  `json:"header_sticky"`
	FooterStyle     string `json:"footer_style"`
	HeroEnabled     *bool  `json:"hero_enabled"`
}

type jsonBareLayoutConfig struct {
	ContentMaxWidth string `json:"content_max_width"`
}

func (l *jsonLayoutConfig) toLayoutConfig() models.LayoutConfig {
	return models.LayoutConfig{
		Name:  l.Name,
		Paths: l.Paths,
		Feeds: l.Feeds,
		Docs: models.DocsLayoutConfig{
			SidebarPosition:    l.Docs.SidebarPosition,
			SidebarWidth:       l.Docs.SidebarWidth,
			SidebarCollapsible: l.Docs.SidebarCollapsible,
			SidebarDefaultOpen: l.Docs.SidebarDefaultOpen,
			TocPosition:        l.Docs.TocPosition,
			TocWidth:           l.Docs.TocWidth,
			TocCollapsible:     l.Docs.TocCollapsible,
			TocDefaultOpen:     l.Docs.TocDefaultOpen,
			ContentMaxWidth:    l.Docs.ContentMaxWidth,
			HeaderStyle:        l.Docs.HeaderStyle,
			FooterStyle:        l.Docs.FooterStyle,
		},
		Blog: models.BlogLayoutConfig{
			ContentMaxWidth: l.Blog.ContentMaxWidth,
			ShowToc:         l.Blog.ShowToc,
			TocPosition:     l.Blog.TocPosition,
			TocWidth:        l.Blog.TocWidth,
			HeaderStyle:     l.Blog.HeaderStyle,
			FooterStyle:     l.Blog.FooterStyle,
			ShowAuthor:      l.Blog.ShowAuthor,
			ShowDate:        l.Blog.ShowDate,
			ShowTags:        l.Blog.ShowTags,
			ShowReadingTime: l.Blog.ShowReadingTime,
			ShowPrevNext:    l.Blog.ShowPrevNext,
		},
		Landing: models.LandingLayoutConfig{
			ContentMaxWidth: l.Landing.ContentMaxWidth,
			HeaderStyle:     l.Landing.HeaderStyle,
			HeaderSticky:    l.Landing.HeaderSticky,
			FooterStyle:     l.Landing.FooterStyle,
			HeroEnabled:     l.Landing.HeroEnabled,
		},
		Bare: models.BareLayoutConfig{
			ContentMaxWidth: l.Bare.ContentMaxWidth,
		},
		Defaults: models.LayoutDefaults{
			ContentMaxWidth: l.Defaults.ContentMaxWidth,
			HeaderSticky:    l.Defaults.HeaderSticky,
			FooterSticky:    l.Defaults.FooterSticky,
		},
	}
}

// Sidebar-related JSON structs

type jsonSidebarConfig struct {
	Enabled      *bool                             `json:"enabled"`
	Position     string                            `json:"position"`
	Width        string                            `json:"width"`
	Collapsible  *bool                             `json:"collapsible"`
	DefaultOpen  *bool                             `json:"default_open"`
	Nav          []jsonSidebarNavItem              `json:"nav"`
	Title        string                            `json:"title"`
	Paths        map[string]*jsonPathSidebarConfig `json:"paths"`
	MultiFeed    *bool                             `json:"multi_feed"`
	Feeds        []string                          `json:"feeds"`
	FeedSections []jsonMultiFeedSection            `json:"feed_sections"`
	AutoGenerate *jsonSidebarAutoGenerate          `json:"auto_generate"`
}

type jsonSidebarNavItem struct {
	Title    string               `json:"title"`
	Href     string               `json:"href"`
	Children []jsonSidebarNavItem `json:"children"`
}

type jsonPathSidebarConfig struct {
	Title        string                   `json:"title"`
	AutoGenerate *jsonSidebarAutoGenerate `json:"auto_generate"`
	Items        []jsonSidebarNavItem     `json:"items"`
	Feed         string                   `json:"feed"`
	Position     string                   `json:"position"`
	Collapsible  *bool                    `json:"collapsible"`
}

type jsonSidebarAutoGenerate struct {
	Directory string   `json:"directory"`
	OrderBy   string   `json:"order_by"`
	Reverse   *bool    `json:"reverse"`
	MaxDepth  int      `json:"max_depth"`
	Exclude   []string `json:"exclude"`
}

type jsonMultiFeedSection struct {
	Feed      string `json:"feed"`
	Title     string `json:"title"`
	Collapsed *bool  `json:"collapsed"`
	MaxItems  int    `json:"max_items"`
}

func convertJSONSidebarNavItems(items []jsonSidebarNavItem) []models.SidebarNavItem {
	result := make([]models.SidebarNavItem, len(items))
	for i, item := range items {
		result[i] = models.SidebarNavItem{
			Title:    item.Title,
			Href:     item.Href,
			Children: convertJSONSidebarNavItems(item.Children),
		}
	}
	return result
}

func (s *jsonSidebarConfig) toSidebarConfig() models.SidebarConfig {
	config := models.SidebarConfig{
		Enabled:     s.Enabled,
		Position:    s.Position,
		Width:       s.Width,
		Collapsible: s.Collapsible,
		DefaultOpen: s.DefaultOpen,
		Nav:         convertJSONSidebarNavItems(s.Nav),
		Title:       s.Title,
		MultiFeed:   s.MultiFeed,
		Feeds:       s.Feeds,
	}

	// Convert paths
	if len(s.Paths) > 0 {
		config.Paths = make(map[string]*models.PathSidebarConfig)
		for path, pathConfig := range s.Paths {
			var autoGen *models.SidebarAutoGenerate
			if pathConfig.AutoGenerate != nil {
				autoGen = &models.SidebarAutoGenerate{
					Directory: pathConfig.AutoGenerate.Directory,
					OrderBy:   pathConfig.AutoGenerate.OrderBy,
					Reverse:   pathConfig.AutoGenerate.Reverse,
					MaxDepth:  pathConfig.AutoGenerate.MaxDepth,
					Exclude:   pathConfig.AutoGenerate.Exclude,
				}
			}
			config.Paths[path] = &models.PathSidebarConfig{
				Title:        pathConfig.Title,
				AutoGenerate: autoGen,
				Items:        convertJSONSidebarNavItems(pathConfig.Items),
				Feed:         pathConfig.Feed,
				Position:     pathConfig.Position,
				Collapsible:  pathConfig.Collapsible,
			}
		}
	}

	// Convert feed sections
	if len(s.FeedSections) > 0 {
		config.FeedSections = make([]models.MultiFeedSection, len(s.FeedSections))
		for i, section := range s.FeedSections {
			config.FeedSections[i] = models.MultiFeedSection{
				Feed:      section.Feed,
				Title:     section.Title,
				Collapsed: section.Collapsed,
				MaxItems:  section.MaxItems,
			}
		}
	}

	// Convert auto-generate
	if s.AutoGenerate != nil {
		config.AutoGenerate = &models.SidebarAutoGenerate{
			Directory: s.AutoGenerate.Directory,
			OrderBy:   s.AutoGenerate.OrderBy,
			Reverse:   s.AutoGenerate.Reverse,
			MaxDepth:  s.AutoGenerate.MaxDepth,
			Exclude:   s.AutoGenerate.Exclude,
		}
	}

	return config
}

// TOC-related JSON structs

type jsonTocConfig struct {
	Enabled     *bool  `json:"enabled"`
	Position    string `json:"position"`
	Width       string `json:"width"`
	MinDepth    int    `json:"min_depth"`
	MaxDepth    int    `json:"max_depth"`
	Title       string `json:"title"`
	Collapsible *bool  `json:"collapsible"`
	DefaultOpen *bool  `json:"default_open"`
	ScrollSpy   *bool  `json:"scroll_spy"`
}

func (t *jsonTocConfig) toTocConfig() models.TocConfig {
	return models.TocConfig{
		Enabled:     t.Enabled,
		Position:    t.Position,
		Width:       t.Width,
		MinDepth:    t.MinDepth,
		MaxDepth:    t.MaxDepth,
		Title:       t.Title,
		Collapsible: t.Collapsible,
		DefaultOpen: t.DefaultOpen,
		ScrollSpy:   t.ScrollSpy,
	}
}

// Header layout JSON structs

type jsonHeaderLayoutConfig struct {
	Style           string `json:"style"`
	Sticky          *bool  `json:"sticky"`
	ShowLogo        *bool  `json:"show_logo"`
	ShowTitle       *bool  `json:"show_title"`
	ShowNav         *bool  `json:"show_nav"`
	ShowSearch      *bool  `json:"show_search"`
	ShowThemeToggle *bool  `json:"show_theme_toggle"`
}

func (h *jsonHeaderLayoutConfig) toHeaderLayoutConfig() models.HeaderLayoutConfig {
	return models.HeaderLayoutConfig{
		Style:           h.Style,
		Sticky:          h.Sticky,
		ShowLogo:        h.ShowLogo,
		ShowTitle:       h.ShowTitle,
		ShowNav:         h.ShowNav,
		ShowSearch:      h.ShowSearch,
		ShowThemeToggle: h.ShowThemeToggle,
	}
}

// Blogroll-related JSON structs

type jsonBlogrollConfig struct {
	Enabled              bool                     `json:"enabled"`
	BlogrollSlug         string                   `json:"blogroll_slug"`
	ReaderSlug           string                   `json:"reader_slug"`
	CacheDir             string                   `json:"cache_dir"`
	CacheDuration        string                   `json:"cache_duration"`
	RefreshOnBuild       *bool                    `json:"refresh_on_build"`
	Timeout              int                      `json:"timeout"`
	ConcurrentRequests   int                      `json:"concurrent_requests"`
	MaxEntriesPerFeed    int                      `json:"max_entries_per_feed"`
	FallbackImageService string                   `json:"fallback_image_service"`
	Feeds                []jsonExternalFeedConfig `json:"feeds"`
	Templates            jsonBlogrollTemplates    `json:"templates"`
}

type jsonExternalFeedConfig struct {
	URL           string   `json:"url"`
	Title         string   `json:"title"`
	Description   string   `json:"description"`
	Category      string   `json:"category"`
	Type          string   `json:"type"`
	Tags          []string `json:"tags"`
	Active        *bool    `json:"active"`
	SiteURL       string   `json:"site_url"`
	ImageURL      string   `json:"image_url"`
	Handle        string   `json:"handle"`
	Aliases       []string `json:"aliases,omitempty"`
	MaxEntries    *int     `json:"max_entries,omitempty"`
	Primary       *bool    `json:"primary,omitempty"`
	PrimaryPerson string   `json:"primary_person"`
}

type jsonBlogrollTemplates struct {
	Blogroll string `json:"blogroll"`
	Reader   string `json:"reader"`
}

func (b *jsonBlogrollConfig) toBlogrollConfig() models.BlogrollConfig {
	// Get default values
	defaults := models.NewBlogrollConfig()

	// Apply defaults for template names if not specified
	blogrollTemplate := b.Templates.Blogroll
	if blogrollTemplate == "" {
		blogrollTemplate = defaults.Templates.Blogroll
	}
	readerTemplate := b.Templates.Reader
	if readerTemplate == "" {
		readerTemplate = defaults.Templates.Reader
	}

	config := models.BlogrollConfig{
		Enabled:              b.Enabled,
		BlogrollSlug:         b.BlogrollSlug,
		ReaderSlug:           b.ReaderSlug,
		CacheDir:             b.CacheDir,
		CacheDuration:        b.CacheDuration,
		RefreshOnBuild:       b.RefreshOnBuild,
		Timeout:              b.Timeout,
		ConcurrentRequests:   b.ConcurrentRequests,
		MaxEntriesPerFeed:    b.MaxEntriesPerFeed,
		FallbackImageService: b.FallbackImageService,
		Templates: models.BlogrollTemplates{
			Blogroll: blogrollTemplate,
			Reader:   readerTemplate,
		},
	}

	for i := range b.Feeds {
		fc := &b.Feeds[i]
		config.Feeds = append(config.Feeds, models.ExternalFeedConfig{
			URL:           fc.URL,
			Title:         fc.Title,
			Description:   fc.Description,
			Category:      fc.Category,
			Type:          models.ReaderFeedType(fc.Type),
			Tags:          fc.Tags,
			Active:        fc.Active,
			SiteURL:       fc.SiteURL,
			ImageURL:      fc.ImageURL,
			Handle:        fc.Handle,
			Aliases:       fc.Aliases,
			MaxEntries:    fc.MaxEntries,
			Primary:       fc.Primary,
			PrimaryPerson: fc.PrimaryPerson,
		})
	}

	return config
}

func (c *jsonComponentsConfig) toComponentsConfig() models.ComponentsConfig {
	config := models.ComponentsConfig{
		Nav: models.NavComponentConfig{
			Enabled:  c.Nav.Enabled,
			Position: c.Nav.Position,
			Style:    c.Nav.Style,
		},
		Footer: models.FooterComponentConfig{
			Enabled:       c.Footer.Enabled,
			Text:          c.Footer.Text,
			ShowCopyright: c.Footer.ShowCopyright,
		},
		DocSidebar: models.DocSidebarConfig{
			Enabled:  c.DocSidebar.Enabled,
			Position: c.DocSidebar.Position,
			Width:    c.DocSidebar.Width,
			MinDepth: c.DocSidebar.MinDepth,
			MaxDepth: c.DocSidebar.MaxDepth,
		},
		FeedSidebar: models.FeedSidebarConfig{
			Enabled:  c.FeedSidebar.Enabled,
			Position: c.FeedSidebar.Position,
			Width:    c.FeedSidebar.Width,
			Title:    c.FeedSidebar.Title,
			Feeds:    c.FeedSidebar.Feeds,
			MaxPosts: c.FeedSidebar.MaxPosts,
		},
		ContentSidebar: models.ContentSidebarConfig{
			Enabled:  c.ContentSidebar.Enabled,
			Position: c.ContentSidebar.Position,
			Width:    c.ContentSidebar.Width,
			Slug:     c.ContentSidebar.Slug,
		},
		CardRouter: models.CardRouterConfig{
			Mappings: c.CardRouter.Mappings,
		},
		Share: models.ShareComponentConfig{
			Enabled:   c.Share.Enabled,
			Platforms: append([]string{}, c.Share.Platforms...),
			Position:  c.Share.Position,
			Title:     c.Share.Title,
			Custom:    map[string]models.SharePlatformConfig{},
		},
		PostCopy: models.PostCopyComponentConfig{Enabled: c.PostCopy.Enabled},
		PostConnections: models.PostConnectionsComponentConfig{
			Enabled:       c.PostConn.Enabled,
			Display:       append([]string{}, c.PostConn.Display...),
			MinLinks:      c.PostConn.MinLinks,
			MaxLinks:      c.PostConn.MaxLinks,
			GraphMinLinks: c.PostConn.GraphMinLinks,
			GraphMaxLinks: c.PostConn.GraphMaxLinks,
			ListMinLinks:  c.PostConn.ListMinLinks,
			ListMaxLinks:  c.PostConn.ListMaxLinks,
			InlinksLimit:  c.PostConn.InlinksLimit,
			OutlinksLimit: c.PostConn.OutlinksLimit,
		},
	}

	if len(c.Share.Platforms) == 0 {
		config.Share.Platforms = nil
	}

	for key, custom := range c.Share.Custom {
		config.Share.Custom[key] = models.SharePlatformConfig{
			Name: custom.Name,
			Icon: custom.Icon,
			URL:  custom.URL,
		}
	}
	if len(c.Share.Custom) == 0 {
		config.Share.Custom = nil
	}

	// Convert nav items
	for _, item := range c.Nav.Items {
		config.Nav.Items = append(config.Nav.Items, models.NavItem{
			Label:    item.Label,
			URL:      item.URL,
			External: item.External,
		})
	}

	// Convert footer links
	for _, link := range c.Footer.Links {
		config.Footer.Links = append(config.Footer.Links, models.NavItem{
			Label:    link.Label,
			URL:      link.URL,
			External: link.External,
		})
	}

	return config
}

func (p *jsonPostFormatsConfig) toPostFormatsConfig() models.PostFormatsConfig {
	return models.PostFormatsConfig{
		HTML:     p.HTML,
		Markdown: p.Markdown,
		Text:     p.Text,
		ANSI:     p.ANSI,
		OG:       p.OG,
	}
}

func (w *jsonWellKnownConfig) toWellKnownConfig() models.WellKnownConfig {
	defaults := models.NewWellKnownConfig()
	config := models.WellKnownConfig{
		Enabled:         w.Enabled,
		AutoGenerate:    w.AutoGenerate,
		SSHFingerprint:  w.SSHFingerprint,
		KeybaseUsername: w.KeybaseUsername,
	}

	if config.Enabled == nil {
		config.Enabled = defaults.Enabled
	}
	if config.AutoGenerate == nil {
		config.AutoGenerate = defaults.AutoGenerate
	}

	return config
}

// configSource interface implementation for jsonConfig.
func (c *jsonConfig) getBaseConfig() baseConfigData {
	return baseConfigData{
		Fontpack:       c.Fontpack,
		FontpacksFile:  c.FontpacksFile,
		OutputDir:      c.OutputDir,
		URL:            c.URL,
		Title:          c.Title,
		Description:    c.Description,
		Author:         c.Author,
		Language:       c.Language,
		AuthorURL:      c.AuthorURL,
		ManagingEditor: c.ManagingEditor,
		WebMaster:      c.WebMaster,
		Copyright:      c.Copyright,
		License:        c.License,
		AssetsDir:      c.AssetsDir,
		Assets:         c.Assets,
		TemplatesDir:   c.TemplatesDir,
		Hooks:          c.Hooks,
		DisabledHooks:  c.DisabledHooks,
		GlobPatterns:   c.Glob.Patterns,
		UseGitignore:   c.Glob.UseGitignore,
		SlugMode:       c.Glob.SlugMode,
		SlugRules:      c.Glob.toSlugRules(),
		Extensions:     c.Markdown.Extensions,
		Concurrency:    c.Concurrency,
		Theme:          c.Theme.toThemeConfig(),
		ThemeCalendar:  c.ThemeCalendar,
		Head:           c.Head.toHeadConfig(),
		Presets:        c.Presets,
		Defaults:       c.Defaults,
		ErrorPages:     c.ErrorPages,
		ResourceHints:  c.ResourceHints,
		Highlight:      c.Markdown.Highlight,
		Footer:         c.Footer.toFooterConfig(),
	}
}

func (c *jsonConfig) getNavItems() []navItemData {
	items := make([]navItemData, len(c.Nav))
	for i, nav := range c.Nav {
		items[i] = navItemData(nav)
	}
	return items
}

func (c *jsonConfig) getFeeds() []feedConfigConverter {
	feeds := make([]feedConfigConverter, len(c.Feeds))
	for i := range c.Feeds {
		feeds[i] = &c.Feeds[i]
	}
	return feeds
}

func (c *jsonConfig) getFeedDefaults() feedDefaultsConverter       { return &c.FeedDefaults }
func (c *jsonConfig) getPostFormats() postFormatsConverter         { return &c.PostFormats }
func (c *jsonConfig) getWellKnown() wellKnownConverter             { return &c.WellKnown }
func (c *jsonConfig) getSEO() seoConverter                         { return &c.SEO }
func (c *jsonConfig) getIndieAuth() indieAuthConverter             { return &c.IndieAuth }
func (c *jsonConfig) getWebmention() webmentionConverter           { return &c.Webmention }
func (c *jsonConfig) getSearch() models.SearchConfig               { return c.Search }
func (c *jsonConfig) getComponents() componentsConverter           { return &c.Components }
func (c *jsonConfig) getLayout() layoutConverter                   { return &c.Layout }
func (c *jsonConfig) getSidebar() sidebarConverter                 { return &c.Sidebar }
func (c *jsonConfig) getToc() tocConverter                         { return &c.Toc }
func (c *jsonConfig) getHeader() headerConverter                   { return &c.Header }
func (c *jsonConfig) getBlogroll() blogrollConverter               { return &c.Blogroll }
func (c *jsonConfig) getTags() tagsConverter                       { return &c.Tags }
func (c *jsonConfig) getFeedsPage() feedsPageConverter             { return &c.FeedsPage }
func (c *jsonConfig) getEncryption() encryptionConverter           { return &c.Encryption }
func (c *jsonConfig) getTagAggregator() tagAggregatorConverter     { return &c.TagAggregator }
func (c *jsonConfig) getMentions() mentionsConverter               { return &c.Mentions }
func (c *jsonConfig) getWebSub() webSubConverter                   { return &c.WebSub }
func (c *jsonConfig) getShortcuts() shortcutsConverter             { return &c.Shortcuts }
func (c *jsonConfig) getViewTransitions() viewTransitionsConverter { return &c.ViewTransitions }
func (c *jsonConfig) getAuthors() authorsConverter                 { return &c.Authors }
func (c *jsonConfig) getGarden() gardenConverter                   { return &c.Garden }
func (c *jsonConfig) getTemplates() templatesConverter             { return &c.Templates }
func (c *jsonConfig) getBuilderAdmin() models.BuilderAdminConfig   { return c.BuilderAdmin }
func (c *jsonConfig) getImages() models.ImagesConfig               { return normalizeImagesConfig(c.Images) }

func (c *jsonConfig) toConfig() *models.Config {
	return buildConfig(c)
}

func (f *jsonFooterConfig) toFooterConfig() models.FooterConfig {
	return models.FooterConfig{
		Text:          f.Text,
		ShowCopyright: f.ShowCopyright,
	}
}

func (f *jsonFeedConfig) toFeedConfig() models.FeedConfig {
	return models.FeedConfig{
		Slug:            f.Slug,
		Presentation:    f.Presentation,
		Title:           f.Title,
		Description:     f.Description,
		Robots:          f.Robots,
		Filter:          f.Filter,
		Sort:            f.Sort,
		Reverse:         f.Reverse,
		Views:           f.Views,
		ShortsPath:      f.ShortsPath,
		Primary:         f.Primary,
		IncludePrivate:  f.IncludePrivate,
		Private:         f.Private,
		Sidebar:         f.Sidebar,
		ItemsPerPage:    f.ItemsPerPage,
		OrphanThreshold: f.OrphanThreshold,
		Limit:           f.Limit,
		Offset:          f.Offset,
		PaginationType:  models.PaginationType(f.PaginationType),
		ArchiveDisabled: f.ArchiveDisabled,
		Formats:         f.Formats.toFeedFormats(),
		Templates:       f.Templates.toFeedTemplates(),
	}
}

func (f *jsonFeedFormats) toFeedFormats() models.FeedFormats {
	formats := models.FeedFormats{}
	if f.HTML != nil {
		formats.HTML = *f.HTML
	}
	if f.SimpleHTML != nil {
		formats.SimpleHTML = *f.SimpleHTML
	}
	if f.RSS != nil {
		formats.RSS = *f.RSS
	}
	if f.Atom != nil {
		formats.Atom = *f.Atom
	}
	if f.JSON != nil {
		formats.JSON = *f.JSON
	}
	if f.Markdown != nil {
		formats.Markdown = *f.Markdown
	}
	if f.Text != nil {
		formats.Text = *f.Text
	}
	if f.Sitemap != nil {
		formats.Sitemap = *f.Sitemap
	}
	return formats
}

func (t *jsonFeedTemplates) toFeedTemplates() models.FeedTemplates {
	return models.FeedTemplates{
		HTML:       t.HTML,
		SimpleHTML: t.SimpleHTML,
		RSS:        t.RSS,
		Atom:       t.Atom,
		JSON:       t.JSON,
		Card:       t.Card,
		Sitemap:    t.Sitemap,
	}
}

func (d *jsonFeedDefaults) toFeedDefaults() models.FeedDefaults {
	return models.FeedDefaults{
		ItemsPerPage:    d.ItemsPerPage,
		Presentation:    d.Presentation,
		OrphanThreshold: d.OrphanThreshold,
		PaginationType:  models.PaginationType(d.PaginationType),
		Views:           d.Views,
		Formats:         d.Formats.toFeedFormats(),
		Templates:       d.Templates.toFeedTemplates(),
		Syndication: models.SyndicationConfig{
			MaxItems:             d.Syndication.MaxItems,
			IncludeContent:       d.Syndication.IncludeContent,
			SiteArchiveDisabled:  d.Syndication.SiteArchiveDisabled,
			FeedArchivesDisabled: d.Syndication.FeedArchivesDisabled,
		},
	}
}
