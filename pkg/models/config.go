package models

import (
	"encoding/json"
	"strings"
)

// NavItem represents a navigation link.
type NavItem struct {
	// Label is the display text for the nav link
	Label string `json:"label" yaml:"label" toml:"label"`

	// URL is the link destination (can be relative or absolute)
	URL string `json:"url" yaml:"url" toml:"url"`

	// External indicates if the link opens in a new tab (default: false)
	External bool `json:"external,omitempty" yaml:"external,omitempty" toml:"external,omitempty"`
}

// FooterConfig configures the site footer.
type FooterConfig struct {
	// Text is the footer text (supports template variables like {{ year }})
	Text string `json:"text,omitempty" yaml:"text,omitempty" toml:"text,omitempty"`

	// ShowCopyright shows the copyright line (default: true)
	ShowCopyright *bool `json:"show_copyright,omitempty" yaml:"show_copyright,omitempty" toml:"show_copyright,omitempty"`
}

// AuthorsConfig configures multi-author support for the site.
type AuthorsConfig struct {
	// GeneratePages enables automatic author bio page generation (default: false)
	GeneratePages bool `json:"generate_pages,omitempty" yaml:"generate_pages,omitempty" toml:"generate_pages,omitempty"`

	// URLPattern defines the URL pattern for author pages (default: "/authors/{author}/")
	URLPattern string `json:"url_pattern,omitempty" yaml:"url_pattern,omitempty" toml:"url_pattern,omitempty"`

	// FeedsEnabled enables author-specific RSS/Atom feeds (default: false)
	FeedsEnabled bool `json:"feeds_enabled,omitempty" yaml:"feeds_enabled,omitempty" toml:"feeds_enabled,omitempty"`

	// Authors is a map of author configurations keyed by author ID
	Authors map[string]Author `json:"authors,omitempty" yaml:"authors,omitempty" toml:"authors,omitempty"`
}

// ComponentsConfig configures the layout components system.
// This enables configuration-driven control over common UI elements.
type ComponentsConfig struct {
	// Nav configures the navigation component
	Nav NavComponentConfig `json:"nav" yaml:"nav" toml:"nav"`

	// Footer configures the footer component
	Footer FooterComponentConfig `json:"footer" yaml:"footer" toml:"footer"`

	// DocSidebar configures the document sidebar (table of contents)
	DocSidebar DocSidebarConfig `json:"doc_sidebar" yaml:"doc_sidebar" toml:"doc_sidebar"`

	// FeedSidebar configures the feed sidebar (series/collection navigation)
	FeedSidebar FeedSidebarConfig `json:"feed_sidebar" yaml:"feed_sidebar" toml:"feed_sidebar"`

	// ContentSidebar configures a content sidebar that renders a post as an aside
	ContentSidebar ContentSidebarConfig `json:"content_sidebar" yaml:"content_sidebar" toml:"content_sidebar"`

	// CardRouter configures the card template routing for feeds
	CardRouter CardRouterConfig `json:"card_router" yaml:"card_router" toml:"card_router"`

	// Share configures the per-post share component
	Share ShareComponentConfig `json:"share" yaml:"share" toml:"share"`

	// PostCopy configures the copy menu in post headers
	PostCopy PostCopyComponentConfig `json:"post_copy" yaml:"post_copy" toml:"post_copy"`

	// PostConnections configures inlink/outlink list and graph rendering on posts
	PostConnections PostConnectionsComponentConfig `json:"post_connections" yaml:"post_connections" toml:"post_connections"`
}

const (
	PostConnectionsDisplayGraph = "graph"
	PostConnectionsDisplayList  = "list"
)

// NavComponentConfig configures the navigation component.
type NavComponentConfig struct {
	// Enabled controls whether navigation is displayed (default: true)
	Enabled *bool `json:"enabled,omitempty" yaml:"enabled,omitempty" toml:"enabled,omitempty"`

	// Position controls where navigation appears: "header", "sidebar" (default: "header")
	Position string `json:"position,omitempty" yaml:"position,omitempty" toml:"position,omitempty"`

	// Style controls the navigation style: "horizontal", "vertical" (default: "horizontal")
	Style string `json:"style,omitempty" yaml:"style,omitempty" toml:"style,omitempty"`

	// Items are the navigation links (overrides top-level nav if set)
	Items []NavItem `json:"items,omitempty" yaml:"items,omitempty" toml:"items,omitempty"`
}

// FooterComponentConfig configures the footer component.
type FooterComponentConfig struct {
	// Enabled controls whether footer is displayed (default: true)
	Enabled *bool `json:"enabled,omitempty" yaml:"enabled,omitempty" toml:"enabled,omitempty"`

	// Text is the footer text (supports template variables)
	Text string `json:"text,omitempty" yaml:"text,omitempty" toml:"text,omitempty"`

	// ShowCopyright shows the copyright line (default: true)
	ShowCopyright *bool `json:"show_copyright,omitempty" yaml:"show_copyright,omitempty" toml:"show_copyright,omitempty"`

	// ShowBuiltWith shows the "built with markata-go" text (default: true)
	ShowBuiltWith *bool `json:"show_built_with,omitempty" yaml:"show_built_with,omitempty" toml:"show_built_with,omitempty"`

	// Links are additional footer links
	Links []NavItem `json:"links,omitempty" yaml:"links,omitempty" toml:"links,omitempty"`
}

// DocSidebarConfig configures the document sidebar (table of contents).
type DocSidebarConfig struct {
	// Enabled controls whether the TOC sidebar is displayed (default: false)
	Enabled *bool `json:"enabled,omitempty" yaml:"enabled,omitempty" toml:"enabled,omitempty"`

	// Position controls sidebar position: "left", "right" (default: "right")
	Position string `json:"position,omitempty" yaml:"position,omitempty" toml:"position,omitempty"`

	// Width is the sidebar width (default: "250px")
	Width string `json:"width,omitempty" yaml:"width,omitempty" toml:"width,omitempty"`

	// MinDepth is the minimum heading level to include (default: 2)
	MinDepth int `json:"min_depth,omitempty" yaml:"min_depth,omitempty" toml:"min_depth,omitempty"`

	// MaxDepth is the maximum heading level to include (default: 4)
	MaxDepth int `json:"max_depth,omitempty" yaml:"max_depth,omitempty" toml:"max_depth,omitempty"`
}

// FeedSidebarConfig configures the feed sidebar (series/collection navigation).
type FeedSidebarConfig struct {
	// Enabled controls whether the feed sidebar is displayed (default: false)
	Enabled *bool `json:"enabled,omitempty" yaml:"enabled,omitempty" toml:"enabled,omitempty"`

	// Position controls sidebar position: "left", "right" (default: "left")
	Position string `json:"position,omitempty" yaml:"position,omitempty" toml:"position,omitempty"`

	// Width is the sidebar width (default: "250px")
	Width string `json:"width,omitempty" yaml:"width,omitempty" toml:"width,omitempty"`

	// Title is the sidebar title (default: uses feed title or "In this series")
	Title string `json:"title,omitempty" yaml:"title,omitempty" toml:"title,omitempty"`

	// Feeds is the list of feed slugs to show navigation for
	Feeds []string `json:"feeds,omitempty" yaml:"feeds,omitempty" toml:"feeds,omitempty"`

	// MaxPosts caps the number of sidebar posts shown. When positive, markata-go
	// keeps a centered window around the current post instead of rendering the
	// entire feed list.
	MaxPosts int `json:"max_posts,omitempty" yaml:"max_posts,omitempty" toml:"max_posts,omitempty"`
}

// ContentSidebarConfig configures a content sidebar that renders a post
// by slug as an aside widget on every page. Unlike feed_sidebar and doc_sidebar
// which are navigational, this displays arbitrary rendered content (e.g. feeds,
// slash page grids, or any unpublished component post).
type ContentSidebarConfig struct {
	// Enabled controls whether the content sidebar is displayed (default: false)
	Enabled *bool `json:"enabled,omitempty" yaml:"enabled,omitempty" toml:"enabled,omitempty"`

	// Position controls sidebar position: "left", "right" (default: "right")
	Position string `json:"position,omitempty" yaml:"position,omitempty" toml:"position,omitempty"`

	// Width is the sidebar width (default: "320px")
	Width string `json:"width,omitempty" yaml:"width,omitempty" toml:"width,omitempty"`

	// Slug is the post slug whose rendered content is included in the sidebar
	Slug string `json:"slug,omitempty" yaml:"slug,omitempty" toml:"slug,omitempty"`
}

// CardRouterConfig configures the card template routing for feeds.
// Maps post template types to card templates, merged with defaults.
type CardRouterConfig struct {
	// Mappings maps post template names to card template names.
	// User mappings are merged with defaults, allowing overrides.
	// Example: {"daily": "article", "meeting": "note"}
	Mappings map[string]string `json:"mappings,omitempty" yaml:"mappings,omitempty" toml:"mappings,omitempty"`
}

// SharePlatformConfig defines a custom share button  entry.
type SharePlatformConfig struct {
	// Name is the accessible label for the platform button
	Name string `json:"name,omitempty" yaml:"name,omitempty" toml:"name,omitempty"`

	// Icon is the path to an icon file (relative to theme assets by default)
	Icon string `json:"icon,omitempty" yaml:"icon,omitempty" toml:"icon,omitempty"`

	// URL is the share URL template (supports {{title}}, {{url}}, {{excerpt}})
	URL string `json:"url,omitempty" yaml:"url,omitempty" toml:"url,omitempty"`
}

// ShareComponentConfig configures the share buttons that appear at the end of posts.
type ShareComponentConfig struct {
	// Enabled toggles the entire component (default: true)
	Enabled *bool `json:"enabled,omitempty" yaml:"enabled,omitempty" toml:"enabled,omitempty"`

	// Platforms controls which platform keys render and in what order
	Platforms []string `json:"platforms,omitempty" yaml:"platforms,omitempty" toml:"platforms,omitempty"`

	// Position adds `share-panel--<position>` modifier for CSS hooks
	Position string `json:"position,omitempty" yaml:"position,omitempty" toml:"position,omitempty"`

	// Title is the heading text displayed above the buttons
	Title string `json:"title,omitempty" yaml:"title,omitempty" toml:"title,omitempty"`

	// Custom maps platform keys to bespoke definitions
	Custom map[string]SharePlatformConfig `json:"custom,omitempty" yaml:"custom,omitempty" toml:"custom,omitempty"`
}

// PostCopyComponentConfig configures the copy menu in post headers.
type PostCopyComponentConfig struct {
	// Enabled shows the copy menu (default: true).
	Enabled *bool `json:"enabled,omitempty" yaml:"enabled,omitempty" toml:"enabled,omitempty"`
}

// IsEnabled reports whether the post copy menu is enabled.
func (c PostCopyComponentConfig) IsEnabled() bool {
	return c.Enabled == nil || *c.Enabled
}

// PostConnectionsComponentConfig configures post connection UI (list/graph).
type PostConnectionsComponentConfig struct {
	// Enabled toggles the entire post connections component (default: true)
	Enabled *bool `json:"enabled,omitempty" yaml:"enabled,omitempty" toml:"enabled,omitempty"`

	// Display selects which render modes are enabled: "graph", "list", or both.
	Display []string `json:"display,omitempty" yaml:"display,omitempty" toml:"display,omitempty"`

	// MinLinks is the minimum total inlinks+outlinks required to show any mode.
	MinLinks int `json:"min_links,omitempty" yaml:"min_links,omitempty" toml:"min_links,omitempty"`

	// MaxLinks is the maximum total inlinks+outlinks allowed to show any mode (0 = no max).
	MaxLinks int `json:"max_links,omitempty" yaml:"max_links,omitempty" toml:"max_links,omitempty"`

	// GraphMinLinks is the minimum total inlinks+outlinks required for graph mode.
	GraphMinLinks int `json:"graph_min_links,omitempty" yaml:"graph_min_links,omitempty" toml:"graph_min_links,omitempty"`

	// GraphMaxLinks is the maximum total inlinks+outlinks allowed for graph mode (0 = no max).
	GraphMaxLinks int `json:"graph_max_links,omitempty" yaml:"graph_max_links,omitempty" toml:"graph_max_links,omitempty"`

	// ListMinLinks is the minimum total inlinks+outlinks required for list mode.
	ListMinLinks int `json:"list_min_links,omitempty" yaml:"list_min_links,omitempty" toml:"list_min_links,omitempty"`

	// ListMaxLinks is the maximum total inlinks+outlinks allowed for list mode (0 = no max).
	ListMaxLinks int `json:"list_max_links,omitempty" yaml:"list_max_links,omitempty" toml:"list_max_links,omitempty"`

	// InlinksLimit limits rendered inlinks in list mode (default: 8).
	InlinksLimit int `json:"inlinks_limit,omitempty" yaml:"inlinks_limit,omitempty" toml:"inlinks_limit,omitempty"`

	// OutlinksLimit limits rendered outlinks in list mode (default: 8).
	OutlinksLimit int `json:"outlinks_limit,omitempty" yaml:"outlinks_limit,omitempty" toml:"outlinks_limit,omitempty"`
}

// NewPostConnectionsComponentConfig returns default post connections config.
func NewPostConnectionsComponentConfig() PostConnectionsComponentConfig {
	enabled := true
	return PostConnectionsComponentConfig{
		Enabled:       &enabled,
		Display:       []string{PostConnectionsDisplayGraph},
		MinLinks:      1,
		MaxLinks:      0,
		GraphMinLinks: 3,
		GraphMaxLinks: 0,
		ListMinLinks:  1,
		ListMaxLinks:  0,
		InlinksLimit:  8,
		OutlinksLimit: 8,
	}
}

// IsEnabled reports whether post connections component is enabled.
func (c PostConnectionsComponentConfig) IsEnabled() bool {
	if c.Enabled == nil {
		return true
	}
	return *c.Enabled
}

// Modes returns configured display modes with sane defaults.
func (c PostConnectionsComponentConfig) Modes() []string {
	if len(c.Display) == 0 {
		return []string{PostConnectionsDisplayGraph}
	}
	return append([]string{}, c.Display...)
}

// HasMode reports whether a display mode is enabled.
func (c PostConnectionsComponentConfig) HasMode(mode string) bool {
	for _, item := range c.Modes() {
		if item == mode {
			return true
		}
	}
	return false
}

// NewShareComponentConfig returns the default share component configuration.
func NewShareComponentConfig() ShareComponentConfig {
	enabled := true
	return ShareComponentConfig{
		Enabled:   &enabled,
		Platforms: append([]string{}, DefaultSharePlatformOrder...),
		Position:  "bottom",
		Title:     "Share this post",
		Custom:    map[string]SharePlatformConfig{},
	}
}

// IsEnabled reports whether the share component is enabled.
func (c ShareComponentConfig) IsEnabled() bool {
	if c.Enabled == nil {
		return true
	}
	return *c.Enabled
}

// DefaultCardMappings returns the default template-to-card mappings.
// These are the built-in mappings that user config merges with.
func DefaultCardMappings() map[string]string {
	return map[string]string{
		// Article card - high prominence for blog-style content
		"blog-post": "article",
		"article":   "article",
		"post":      "article",
		"essay":     "article",
		"tutorial":  "article",

		// Note card - low prominence for short updates
		"note":    "note",
		"ping":    "note",
		"thought": "note",
		"status":  "note",
		"tweet":   "note",

		// Photo card - image prominent
		"photo":   "photo",
		"shot":    "photo",
		"shots":   "photo",
		"image":   "photo",
		"gallery": "photo",

		// Video card - video/thumbnail prominent
		"video":  "video",
		"clip":   "video",
		"cast":   "video",
		"stream": "video",

		// Link card - URL preview style
		"link":     "link",
		"bookmark": "link",
		"til":      "inline",
		"stars":    "link",

		// Quote card - blockquote styling
		"quote":     "quote",
		"quotation": "quote",

		// Guide card - step/chapter indicator
		"guide":   "guide",
		"series":  "guide",
		"step":    "guide",
		"chapter": "guide",

		// Inline card - full rendered content
		"gratitude": "inline",
		"inline":    "inline",
		"micro":     "inline",

		// Contact card - person/character profile
		"contact":   "contact",
		"character": "contact",
		"person":    "contact",
	}
}

// GetCardTemplate returns the card template for a given post template.
// Returns the mapped card name, or "default" if not found.
func (c *CardRouterConfig) GetCardTemplate(postTemplate string) string {
	// First check user mappings
	if c.Mappings != nil {
		if card, ok := c.Mappings[postTemplate]; ok {
			return card
		}
	}

	// Fall back to defaults
	defaults := DefaultCardMappings()
	if card, ok := defaults[postTemplate]; ok {
		return card
	}

	return "default"
}

// MergedMappings returns the user mappings merged with defaults.
// User mappings take precedence over defaults.
func (c *CardRouterConfig) MergedMappings() map[string]string {
	result := DefaultCardMappings()

	// Overlay user mappings
	if c.Mappings != nil {
		for k, v := range c.Mappings {
			result[k] = v
		}
	}

	return result
}

// NewComponentsConfig creates a new ComponentsConfig with default values.
func NewComponentsConfig() ComponentsConfig {
	navEnabled := true
	postCopyEnabled := true
	footerEnabled := true
	docSidebarEnabled := false
	feedSidebarEnabled := false
	contentSidebarEnabled := false
	showCopyright := true
	showBuiltWith := true

	return ComponentsConfig{
		Nav: NavComponentConfig{
			Enabled:  &navEnabled,
			Position: "header",
			Style:    "horizontal",
		},
		Footer: FooterComponentConfig{
			Enabled:       &footerEnabled,
			ShowCopyright: &showCopyright,
			ShowBuiltWith: &showBuiltWith,
		},
		DocSidebar: DocSidebarConfig{
			Enabled:  &docSidebarEnabled,
			Position: "right",
			Width:    "250px",
			MinDepth: 2,
			MaxDepth: 4,
		},
		FeedSidebar: FeedSidebarConfig{
			Enabled:  &feedSidebarEnabled,
			Position: "left",
			Width:    "250px",
		},
		ContentSidebar: ContentSidebarConfig{
			Enabled:  &contentSidebarEnabled,
			Position: "right",
			Width:    "380px",
		},
		Share:           NewShareComponentConfig(),
		PostCopy:        PostCopyComponentConfig{Enabled: &postCopyEnabled},
		PostConnections: NewPostConnectionsComponentConfig(),
	}
}

// IsNavEnabled returns whether navigation is enabled.
func (c *ComponentsConfig) IsNavEnabled() bool {
	if c.Nav.Enabled == nil {
		return true
	}
	return *c.Nav.Enabled
}

// IsFooterEnabled returns whether footer is enabled.
func (c *ComponentsConfig) IsFooterEnabled() bool {
	if c.Footer.Enabled == nil {
		return true
	}
	return *c.Footer.Enabled
}

// IsDocSidebarEnabled returns whether the document sidebar is enabled.
func (c *ComponentsConfig) IsDocSidebarEnabled() bool {
	if c.DocSidebar.Enabled == nil {
		return false
	}
	return *c.DocSidebar.Enabled
}

// IsFeedSidebarEnabled returns whether the feed sidebar is enabled.
func (c *ComponentsConfig) IsFeedSidebarEnabled() bool {
	if c.FeedSidebar.Enabled == nil {
		return false
	}
	return *c.FeedSidebar.Enabled
}

// IsContentSidebarEnabled returns whether the content sidebar is enabled.
func (c *ComponentsConfig) IsContentSidebarEnabled() bool {
	if c.ContentSidebar.Enabled == nil {
		return false
	}
	return *c.ContentSidebar.Enabled
}

// Config represents the site configuration for markata-go.
type Config struct {
	// Fontpack selects a catalog typography pack. It defaults to the zero-download system pack.
	Fontpack string `json:"fontpack,omitempty" yaml:"fontpack,omitempty" toml:"fontpack,omitempty"`

	// FontpacksFile selects an explicit catalog; paths inside it are relative
	// to that catalog file.
	FontpacksFile string `json:"fontpacks_file,omitempty" yaml:"fontpacks_file,omitempty" toml:"fontpacks_file,omitempty"`

	// OutputDir is the directory where generated files are written (default: "output")
	OutputDir string `json:"output_dir" yaml:"output_dir" toml:"output_dir"`

	// URL is the base URL of the site
	URL string `json:"url" yaml:"url" toml:"url"`

	// Title is the site title
	Title string `json:"title" yaml:"title" toml:"title"`

	// Description is the site description
	Description string `json:"description" yaml:"description" toml:"description"`

	// Author is the site author
	Author string `json:"author" yaml:"author" toml:"author"`

	// Language is the site language used in feed and metadata outputs.
	Language string `json:"language,omitempty" yaml:"language,omitempty" toml:"language,omitempty"`

	// AuthorURL is the canonical URL for the site author.
	AuthorURL string `json:"author_url,omitempty" yaml:"author_url,omitempty" toml:"author_url,omitempty"`

	// ManagingEditor is the RSS managing editor contact.
	ManagingEditor string `json:"managing_editor,omitempty" yaml:"managing_editor,omitempty" toml:"managing_editor,omitempty"`

	// WebMaster is the RSS webmaster contact.
	WebMaster string `json:"webmaster,omitempty" yaml:"webmaster,omitempty" toml:"webmaster,omitempty"`

	// Copyright is the copyright notice for syndication outputs.
	Copyright string `json:"copyright,omitempty" yaml:"copyright,omitempty" toml:"copyright,omitempty"`

	// License controls the footer attribution (string key or false)
	License LicenseValue `json:"license,omitempty" yaml:"license,omitempty" toml:"license,omitempty"`

	// AssetsDir is the directory containing static assets (default: "static")
	AssetsDir string `json:"assets_dir" yaml:"assets_dir" toml:"assets_dir"`

	// TemplatesDir is the directory containing templates (default: "templates")
	TemplatesDir string `json:"templates_dir" yaml:"templates_dir" toml:"templates_dir"`

	// Templates config exposes helper-specific settings for rendering built-in cards
	Templates TemplatesConfig `json:"templates" yaml:"templates" toml:"templates"`

	// Nav is the list of navigation links
	Nav []NavItem `json:"nav" yaml:"nav" toml:"nav"`

	// Footer configures the site footer
	Footer FooterConfig `json:"footer" yaml:"footer" toml:"footer"`

	// Hooks is the list of hooks to run (default: ["default"])
	Hooks []string `json:"hooks" yaml:"hooks" toml:"hooks"`

	// DisabledHooks is the list of hooks to disable
	DisabledHooks []string `json:"disabled_hooks" yaml:"disabled_hooks" toml:"disabled_hooks"`

	// GlobConfig configures file globbing behavior
	GlobConfig GlobConfig `json:"glob" yaml:"glob" toml:"glob"`

	// MarkdownConfig configures markdown processing
	MarkdownConfig MarkdownConfig `json:"markdown" yaml:"markdown" toml:"markdown"`

	// Feeds is the list of feed configurations
	Feeds []FeedConfig `json:"feeds" yaml:"feeds" toml:"feeds"`

	// FeedDefaults provides default values for feed configurations
	FeedDefaults FeedDefaults `json:"feed_defaults" yaml:"feed_defaults" toml:"feed_defaults"`

	// Concurrency is the number of concurrent workers (default: 0 = auto)
	Concurrency int `json:"concurrency" yaml:"concurrency" toml:"concurrency"`

	// Theme configures the site theme
	Theme ThemeConfig `json:"theme" yaml:"theme" toml:"theme"`

	// ThemeCalendar configures automatic seasonal theme switching based on date ranges
	ThemeCalendar ThemeCalendarConfig `json:"theme_calendar" yaml:"theme_calendar" toml:"theme_calendar"`

	// PostFormats configures output formats for individual posts
	PostFormats PostFormatsConfig `json:"post_formats" yaml:"post_formats" toml:"post_formats"`

	// WellKnown configures auto-generated .well-known endpoints
	WellKnown WellKnownConfig `json:"well_known" yaml:"well_known" toml:"well_known"`

	// SEO configures SEO metadata generation
	SEO SEOConfig `json:"seo" yaml:"seo" toml:"seo"`

	// IndieAuth configures IndieAuth link tags for identity and authentication
	IndieAuth IndieAuthConfig `json:"indieauth" yaml:"indieauth" toml:"indieauth"`

	// Webmention configures Webmention endpoint for receiving mentions
	Webmention WebmentionConfig `json:"webmention" yaml:"webmention" toml:"webmention"`

	// WebSub configures WebSub discovery links for feeds
	WebSub WebSubConfig `json:"websub" yaml:"websub" toml:"websub"`

	// Components configures layout components (nav, footer, sidebar)
	Components ComponentsConfig `json:"components" yaml:"components" toml:"components"`

	// Head configures elements added to the HTML <head> section
	Head HeadConfig `json:"head" yaml:"head" toml:"head"`

	// Search configures site-wide search functionality using Pagefind
	Search SearchConfig `json:"search" yaml:"search" toml:"search"`

	// Layout configures the layout system for page structure
	Layout LayoutConfig `json:"layout" yaml:"layout" toml:"layout"`

	// Sidebar configures the sidebar navigation component
	Sidebar SidebarConfig `json:"sidebar" yaml:"sidebar" toml:"sidebar"`

	// Toc configures the table of contents component
	Toc TocConfig `json:"toc" yaml:"toc" toml:"toc"`

	// Header configures the header component for layouts
	Header HeaderLayoutConfig `json:"header" yaml:"header" toml:"header"`

	// FooterLayout configures the footer component for layouts
	FooterLayout FooterLayoutConfig `json:"footer_layout" yaml:"footer_layout" toml:"footer_layout"`

	// ContentTemplates configures the content template system for the new command
	ContentTemplates ContentTemplatesConfig `json:"content_templates" yaml:"content_templates" toml:"content_templates"`

	// Blogroll configures the blogroll and RSS reader functionality
	Blogroll BlogrollConfig `json:"blogroll" yaml:"blogroll" toml:"blogroll"`

	// Mentions configures the @mentions resolution plugin
	Mentions MentionsConfig `json:"mentions" yaml:"mentions" toml:"mentions"`

	// ErrorPages configures custom error pages (404, etc.)
	ErrorPages ErrorPagesConfig `json:"error_pages" yaml:"error_pages" toml:"error_pages"`

	// ResourceHints configures automatic resource hints generation (preconnect, dns-prefetch, etc.)
	ResourceHints ResourceHintsConfig `json:"resource_hints" yaml:"resource_hints" toml:"resource_hints"`

	// Encryption configures content encryption for private posts
	Encryption EncryptionConfig `json:"encryption" yaml:"encryption" toml:"encryption"`

	// Shortcuts configures user-defined keyboard shortcuts
	Shortcuts ShortcutsConfig `json:"shortcuts" yaml:"shortcuts" toml:"shortcuts"`

	// ViewTransitions configures client-side animated page navigation.
	ViewTransitions ViewTransitionsConfig `json:"view_transitions" yaml:"view_transitions" toml:"view_transitions"`

	// Tags configures the tags listing page at /tags
	Tags TagsConfig `json:"tags" yaml:"tags" toml:"tags"`

	// FeedsPage configures the feeds listing page at /feeds
	FeedsPage FeedsPageConfig `json:"feeds_page" yaml:"feeds_page" toml:"feeds_page"`

	// Garden configures the garden view plugin for knowledge graph export and visualization
	Garden GardenConfig `json:"garden" yaml:"garden" toml:"garden"`

	// Images configures the generated image inventory and authoring page.
	Images ImagesConfig `json:"images" yaml:"images" toml:"images"`

	// TagAggregator configures tag normalization and hierarchical expansion
	TagAggregator TagAggregatorConfig `json:"tag_aggregator" yaml:"tag_aggregator" toml:"tag_aggregator"`

	// Assets configures external CDN asset handling for self-hosting
	Assets AssetsConfig `json:"assets" yaml:"assets" toml:"assets"`

	// TemplatePresets defines named template preset configurations
	// Each preset specifies templates for all output formats
	TemplatePresets map[string]TemplatePreset `json:"template_presets,omitempty" yaml:"template_presets,omitempty" toml:"template_presets,omitempty"`

	// DefaultTemplates specifies default templates per output format
	// Keys: "html", "txt", "markdown", "og"
	// Values: template file names
	DefaultTemplates map[string]string `json:"default_templates,omitempty" yaml:"default_templates,omitempty" toml:"default_templates,omitempty"`

	// Authors configures multi-author support for the site
	Authors AuthorsConfig `json:"authors" yaml:"authors" toml:"authors"`

	// BuilderAdmin configures the long-lived builder-admin service.
	BuilderAdmin BuilderAdminConfig `json:"builder_admin" yaml:"builder_admin" toml:"builder_admin"`

	// Extra holds arbitrary plugin configurations that aren't part of the core config.
	// Plugin-specific configs like [markata-go.image_zoom] are stored here.
	Extra map[string]any `json:"-" yaml:"-" toml:"-"`
}

// BuilderAdminConfig configures builder-admin settings shared by CLI, environment, and site config.
type BuilderAdminConfig struct {
	Auth    BuilderAdminAuthConfig    `json:"auth" yaml:"auth" toml:"auth"`
	Webhook BuilderAdminWebhookConfig `json:"webhook" yaml:"webhook" toml:"webhook"`

	// Extra preserves unrecognized builder_admin settings without exposing them
	// through configuration serialization or template context.
	Extra map[string]any `json:"-" yaml:"-" toml:"-"`
}

// BuilderAdminWebhookConfig configures signed GitHub and Forgejo push webhooks.
type BuilderAdminWebhookConfig struct {
	Enabled *bool   `json:"enabled" yaml:"enabled" toml:"enabled"`
	Branch  *string `json:"branch" yaml:"branch" toml:"branch"`
	Secret  *string `json:"secret" yaml:"secret" toml:"secret"`
}

// MarshalJSON excludes the webhook secret from serialized configuration while
// retaining JSON decoding support for configuration files.
func (c BuilderAdminWebhookConfig) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Enabled *bool   `json:"enabled"`
		Branch  *string `json:"branch"`
	}{
		Enabled: c.Enabled,
		Branch:  c.Branch,
	})
}

// BuilderAdminAuthConfig configures identity assertions from a trusted ForwardAuth proxy.
type BuilderAdminAuthConfig struct {
	Headers BuilderAdminAuthHeadersConfig `json:"headers" yaml:"headers" toml:"headers"`
}

// BuilderAdminAuthHeadersConfig maps trusted proxy headers. Pointers preserve an explicitly
// configured empty optional header, which disables that display-only assertion.
type BuilderAdminAuthHeadersConfig struct {
	UserID      *string `json:"user_id" yaml:"user_id" toml:"user_id"`
	Username    *string `json:"username" yaml:"username" toml:"username"`
	DisplayName *string `json:"display_name" yaml:"display_name" toml:"display_name"`
	Email       *string `json:"email" yaml:"email" toml:"email"`
	Groups      *string `json:"groups" yaml:"groups" toml:"groups"`
	Roles       *string `json:"roles" yaml:"roles" toml:"roles"`
	Scopes      *string `json:"scopes" yaml:"scopes" toml:"scopes"`
	AuthMethod  *string `json:"auth_method" yaml:"auth_method" toml:"auth_method"`
}

// DefaultTrustedMediaDomains is the default allowlist for media helpers.
var DefaultTrustedMediaDomains = []string{
	"dropper.wayl.one",
	"dropper.waylonwalker.com",
	"dropper-dev.wayl.one",
}

// TemplatesMediaConfig configures template helper behavior around media URLs.
type TemplatesMediaConfig struct {
	// TrustedDomains is the allowlist of hosts that can receive auto sizing/poster helpers.
	TrustedDomains []string `json:"trusted_domains" yaml:"trusted_domains" toml:"trusted_domains"`
}

// TemplatesConfig contains template-specific configuration.
type TemplatesConfig struct {
	// Media holds settings that affect media helpers (with_size, poster_url, etc.)
	Media TemplatesMediaConfig `json:"media" yaml:"media" toml:"media"`
}

// NewTemplatesMediaConfig returns the default TemplatesMediaConfig.
func NewTemplatesMediaConfig() TemplatesMediaConfig {
	trusted := append([]string{}, DefaultTrustedMediaDomains...)
	return TemplatesMediaConfig{TrustedDomains: trusted}
}

// NewTemplatesConfig returns the default TemplatesConfig.
func NewTemplatesConfig() TemplatesConfig {
	return TemplatesConfig{Media: NewTemplatesMediaConfig()}
}

// HeadConfig configures elements added to the HTML <head> section.
type HeadConfig struct {
	// Text is raw HTML/text to include in the head (use with caution)
	Text string `json:"text,omitempty" yaml:"text,omitempty" toml:"text,omitempty"`

	// Meta is a list of meta tags to include
	Meta []MetaTag `json:"meta,omitempty" yaml:"meta,omitempty" toml:"meta,omitempty"`

	// Link is a list of link tags to include
	Link []LinkTag `json:"link,omitempty" yaml:"link,omitempty" toml:"link,omitempty"`

	// Script is a list of script tags to include
	Script []ScriptTag `json:"script,omitempty" yaml:"script,omitempty" toml:"script,omitempty"`

	// AlternateFeeds configures which feeds get <link rel="alternate"> tags
	// If empty, defaults to RSS and Atom feeds
	AlternateFeeds []AlternateFeed `json:"alternate_feeds,omitempty" yaml:"alternate_feeds,omitempty" toml:"alternate_feeds,omitempty"`
}

// MetaTag represents a <meta> tag configuration.
type MetaTag struct {
	Name     string `json:"name,omitempty" yaml:"name,omitempty" toml:"name,omitempty"`
	Property string `json:"property,omitempty" yaml:"property,omitempty" toml:"property,omitempty"`
	Content  string `json:"content" yaml:"content" toml:"content"`
}

// LinkTag represents a <link> tag configuration.
type LinkTag struct {
	Rel         string `json:"rel" yaml:"rel" toml:"rel"`
	Href        string `json:"href" yaml:"href" toml:"href"`
	Crossorigin bool   `json:"crossorigin,omitempty" yaml:"crossorigin,omitempty" toml:"crossorigin,omitempty"`
}

// ScriptTag represents a <script> tag configuration.
type ScriptTag struct {
	Src string `json:"src" yaml:"src" toml:"src"`
}

// AlternateFeed configures a <link rel="alternate"> tag for feed discovery.
type AlternateFeed struct {
	// Type is the feed type: "rss", "atom", or "json"
	Type string `json:"type" yaml:"type" toml:"type"`

	// Title is the human-readable feed title (e.g., "RSS Feed")
	Title string `json:"title" yaml:"title" toml:"title"`

	// Href is the URL path to the feed (e.g., "/rss.xml")
	Href string `json:"href" yaml:"href" toml:"href"`
}

// GetMIMEType returns the MIME type for this feed type.
func (f *AlternateFeed) GetMIMEType() string {
	switch f.Type {
	case "rss":
		return "application/rss+xml"
	case "atom":
		return "application/atom+xml"
	case "json":
		return "application/feed+json"
	default:
		return "application/xml"
	}
}

// SearchConfig configures site-wide search functionality.
type SearchConfig struct {
	// Enabled controls whether search is active (default: true)
	Enabled *bool `json:"enabled,omitempty" yaml:"enabled,omitempty" toml:"enabled,omitempty"`

	// Backend selects the search engine: "pagefind" (default) or "bleve"
	Backend string `json:"backend,omitempty" yaml:"backend,omitempty" toml:"backend,omitempty"`

	// Endpoint is the API path for bleve search (default: "/api/search")
	Endpoint string `json:"endpoint,omitempty" yaml:"endpoint,omitempty" toml:"endpoint,omitempty"`

	// Position controls where search UI appears: "navbar", "sidebar", "footer", "custom"
	Position string `json:"position,omitempty" yaml:"position,omitempty" toml:"position,omitempty"`

	// Placeholder is the search input placeholder text
	Placeholder string `json:"placeholder,omitempty" yaml:"placeholder,omitempty" toml:"placeholder,omitempty"`

	// ShowImages shows thumbnails in search results
	ShowImages *bool `json:"show_images,omitempty" yaml:"show_images,omitempty" toml:"show_images,omitempty"`

	// ExcerptLength is the character limit for result excerpts
	ExcerptLength int `json:"excerpt_length,omitempty" yaml:"excerpt_length,omitempty" toml:"excerpt_length,omitempty"`

	// Pagefind configures the Pagefind CLI options
	Pagefind PagefindConfig `json:"pagefind,omitempty" yaml:"pagefind,omitempty" toml:"pagefind,omitempty"`

	// Bleve configures the bleve search API options
	Bleve BleveSearchConfig `json:"bleve,omitempty" yaml:"bleve,omitempty" toml:"bleve,omitempty"`

	// Feeds configures feed-specific search instances
	Feeds []SearchFeedConfig `json:"feeds,omitempty" yaml:"feeds,omitempty" toml:"feeds,omitempty"`
}

// SearchBackend returns the configured search backend, defaulting to "pagefind".
func (s *SearchConfig) SearchBackend() string {
	if s.Backend == "" {
		return "pagefind"
	}
	return s.Backend
}

// SearchEndpoint returns the configured API endpoint, defaulting to "/api/search".
func (s *SearchConfig) SearchEndpoint() string {
	if s.Endpoint == "" {
		return "/api/search"
	}
	return s.Endpoint
}

// BleveSearchConfig configures bleve search API defaults.
type BleveSearchConfig struct {
	// Endpoint is the absolute or relative API URL used by the frontend bleve client.
	// When set, the navbar search prefers bleve over Pagefind.
	Endpoint string `json:"endpoint,omitempty" yaml:"endpoint,omitempty" toml:"endpoint,omitempty"`

	// Fuzzy enables fuzzy matching by default (default: false)
	Fuzzy *bool `json:"fuzzy,omitempty" yaml:"fuzzy,omitempty" toml:"fuzzy,omitempty"`

	// Limit is the default max results (default: 20)
	Limit int `json:"limit,omitempty" yaml:"limit,omitempty" toml:"limit,omitempty"`

	// MaxLimit is the maximum allowed limit from API requests (default: 100)
	MaxLimit int `json:"max_limit,omitempty" yaml:"max_limit,omitempty" toml:"max_limit,omitempty"`

	// CORSOrigins is the list of allowed CORS origins (default: ["*"] in dev)
	CORSOrigins []string `json:"cors_origins,omitempty" yaml:"cors_origins,omitempty" toml:"cors_origins,omitempty"`
}

// EndpointOrDefault returns the configured bleve endpoint or the top-level search endpoint.
func (b *BleveSearchConfig) EndpointOrDefault(searchEndpoint string) string {
	if strings.TrimSpace(b.Endpoint) != "" {
		return b.Endpoint
	}
	return searchEndpoint
}

// IsFuzzy returns whether fuzzy matching is enabled by default.
func (b *BleveSearchConfig) IsFuzzy() bool {
	if b.Fuzzy == nil {
		return false
	}
	return *b.Fuzzy
}

// DefaultLimit returns the default result limit.
func (b *BleveSearchConfig) DefaultLimit() int {
	if b.Limit <= 0 {
		return 20
	}
	return b.Limit
}

// GetMaxLimit returns the maximum allowed result limit.
func (b *BleveSearchConfig) GetMaxLimit() int {
	if b.MaxLimit <= 0 {
		return 100
	}
	return b.MaxLimit
}

// IsEnabled returns whether search is enabled.
// Defaults to true if not explicitly set.
func (s *SearchConfig) IsEnabled() bool {
	if s.Enabled == nil {
		return true
	}
	return *s.Enabled
}

// IsShowImages returns whether to show images in search results.
// Defaults to true if not explicitly set.
func (s *SearchConfig) IsShowImages() bool {
	if s.ShowImages == nil {
		return true
	}
	return *s.ShowImages
}

// PagefindConfig configures Pagefind CLI behavior.
type PagefindConfig struct {
	// BundleDir is the output directory for search index (default: "_pagefind")
	BundleDir string `json:"bundle_dir,omitempty" yaml:"bundle_dir,omitempty" toml:"bundle_dir,omitempty"`

	// ExcludeSelectors are CSS selectors for elements to exclude from indexing
	ExcludeSelectors []string `json:"exclude_selectors,omitempty" yaml:"exclude_selectors,omitempty" toml:"exclude_selectors,omitempty"`

	// RootSelector is the CSS selector for the searchable content container
	RootSelector string `json:"root_selector,omitempty" yaml:"root_selector,omitempty" toml:"root_selector,omitempty"`

	// AutoInstall enables automatic Pagefind binary installation (default: true)
	AutoInstall *bool `json:"auto_install,omitempty" yaml:"auto_install,omitempty" toml:"auto_install,omitempty"`

	// Version is the Pagefind version to install (default: "latest")
	Version string `json:"version,omitempty" yaml:"version,omitempty" toml:"version,omitempty"`

	// CacheDir is the directory for caching Pagefind binaries (default: XDG cache)
	CacheDir string `json:"cache_dir,omitempty" yaml:"cache_dir,omitempty" toml:"cache_dir,omitempty"`

	// Verbose enables verbose output from Pagefind (default: false)
	// When false, only errors are shown. When true or when --verbose CLI flag is used, all output is shown.
	Verbose *bool `json:"verbose,omitempty" yaml:"verbose,omitempty" toml:"verbose,omitempty"`
}

// IsAutoInstallEnabled returns whether automatic Pagefind installation is enabled.
// Defaults to true if not explicitly set.
func (p *PagefindConfig) IsAutoInstallEnabled() bool {
	if p.AutoInstall == nil {
		return true
	}
	return *p.AutoInstall
}

// IsVerbose returns whether verbose output is enabled.
// Defaults to false if not explicitly set.
func (p *PagefindConfig) IsVerbose() bool {
	if p.Verbose == nil {
		return false
	}
	return *p.Verbose
}

// SearchFeedConfig configures a feed-specific search instance.
type SearchFeedConfig struct {
	// Name is the search instance identifier
	Name string `json:"name" yaml:"name" toml:"name"`

	// Filter is the filter expression for posts in this search
	Filter string `json:"filter" yaml:"filter" toml:"filter"`

	// Position controls where this search UI appears
	Position string `json:"position,omitempty" yaml:"position,omitempty" toml:"position,omitempty"`

	// Placeholder is the search input placeholder text
	Placeholder string `json:"placeholder,omitempty" yaml:"placeholder,omitempty" toml:"placeholder,omitempty"`
}

// NewSearchConfig creates a new SearchConfig with default values.
func NewSearchConfig() SearchConfig {
	enabled := true
	showImages := true
	autoInstall := true
	return SearchConfig{
		Enabled:       &enabled,
		Position:      "navbar",
		Placeholder:   "Search…",
		ShowImages:    &showImages,
		ExcerptLength: 200,
		Pagefind: PagefindConfig{
			BundleDir:        "_pagefind",
			ExcludeSelectors: []string{},
			RootSelector:     "",
			AutoInstall:      &autoInstall,
			Version:          "latest",
			CacheDir:         "",
		},
		Feeds: []SearchFeedConfig{},
	}
}

// FontConfig configures typography settings for the site.
type FontConfig struct {
	// Family is the primary font family for body text (default: "system-ui, -apple-system, sans-serif")
	Family string `json:"family,omitempty" yaml:"family,omitempty" toml:"family,omitempty"`

	// HeadingFamily is the font family for headings (default: inherits from Family)
	HeadingFamily string `json:"heading_family,omitempty" yaml:"heading_family,omitempty" toml:"heading_family,omitempty"`

	// CodeFamily is the font family for code blocks and inline code (default: "ui-monospace, monospace")
	CodeFamily string `json:"code_family,omitempty" yaml:"code_family,omitempty" toml:"code_family,omitempty"`

	// Size is the base font size (default: "16px")
	Size string `json:"size,omitempty" yaml:"size,omitempty" toml:"size,omitempty"`

	// LineHeight is the base line height (default: "1.6")
	LineHeight string `json:"line_height,omitempty" yaml:"line_height,omitempty" toml:"line_height,omitempty"`

	// GoogleFonts is a list of Google Fonts to load (e.g., ["Inter", "Fira Code"])
	GoogleFonts []string `json:"google_fonts,omitempty" yaml:"google_fonts,omitempty" toml:"google_fonts,omitempty"`

	// CustomURLs is a list of custom font CSS URLs to load
	CustomURLs []string `json:"custom_urls,omitempty" yaml:"custom_urls,omitempty" toml:"custom_urls,omitempty"`
}

// NewFontConfig creates a new FontConfig with default values.
func NewFontConfig() FontConfig {
	return FontConfig{
		Family:        "system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif",
		HeadingFamily: "",
		CodeFamily:    "ui-monospace, SFMono-Regular, 'SF Mono', Menlo, Consolas, monospace",
		Size:          "16px",
		LineHeight:    "1.6",
		GoogleFonts:   []string{},
		CustomURLs:    []string{},
	}
}

// GetHeadingFamily returns the heading font family, falling back to Family if not set.
func (f *FontConfig) GetHeadingFamily() string {
	if f.HeadingFamily != "" {
		return f.HeadingFamily
	}
	return f.Family
}

// GetGoogleFontsURL returns the Google Fonts CSS URL for the configured fonts.
// Returns empty string if no Google Fonts are configured.
func (f *FontConfig) GetGoogleFontsURL() string {
	if len(f.GoogleFonts) == 0 {
		return ""
	}
	// Build Google Fonts URL with all fonts
	// Format: https://fonts.googleapis.com/css2?family=Font+Name:wght@400;700&family=Other+Font
	families := make([]string, len(f.GoogleFonts))
	for i, font := range f.GoogleFonts {
		// Replace spaces with + for URL encoding
		encoded := ""
		for _, c := range font {
			if c == ' ' {
				encoded += "+"
			} else {
				encoded += string(c)
			}
		}
		families[i] = "family=" + encoded + ":wght@400;500;600;700"
	}
	return "https://fonts.googleapis.com/css2?" + joinStrings(families, "&") + "&display=swap"
}

// joinStrings joins strings with a separator (helper to avoid importing strings package).
func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}

// ThemeConfig configures the site theme.
type ThemeConfig struct {
	ContractVersion int    `json:"contract_version,omitempty" yaml:"contract_version,omitempty" toml:"contract_version,omitempty"`
	Fontpack        string `json:"fontpack,omitempty" yaml:"fontpack,omitempty" toml:"fontpack,omitempty"`
	// Name is the theme name (default: "default")
	Name string `json:"name" yaml:"name" toml:"name"`

	// Aesthetic configures the structural token set to use
	Aesthetic string `json:"aesthetic,omitempty" yaml:"aesthetic,omitempty" toml:"aesthetic,omitempty"`

	// Palette is the base color palette to use (default: "default-light")
	// When set to a base name like "everforest", the system will auto-detect
	// light/dark variants (e.g., "everforest-light" and "everforest-dark")
	Palette string `json:"palette" yaml:"palette" toml:"palette"`

	// PaletteLight is the palette to use for light mode (optional)
	// If not set, auto-detected from base Palette name
	PaletteLight string `json:"palette_light,omitempty" yaml:"palette_light,omitempty" toml:"palette_light,omitempty"`

	// PaletteDark is the palette to use for dark mode (optional)
	// If not set, auto-detected from base Palette name
	PaletteDark string `json:"palette_dark,omitempty" yaml:"palette_dark,omitempty" toml:"palette_dark,omitempty"`

	// FallbackMode controls which mode is used when no explicit user preference
	// is saved and the system preference is unavailable/unspecified.
	// Valid values: "dark", "light" (default: "dark").
	FallbackMode string `json:"fallback_mode,omitempty" yaml:"fallback_mode,omitempty" toml:"fallback_mode,omitempty"`

	// Seasonal makes the seasonal palette the default for visitors who have not
	// picked a theme. It follows northern hemisphere seasons and switches to a
	// holiday palette a few days before and on world holidays. It requires the
	// theme switcher, which ships the seasonal palettes.
	Seasonal bool `json:"seasonal,omitempty" yaml:"seasonal,omitempty" toml:"seasonal,omitempty"`

	// TextSize is the default reading-size preset for visitors without a saved
	// preference. Valid values are "small", "medium", "large", and "x-large".
	TextSize string `json:"text_size,omitempty" yaml:"text_size,omitempty" toml:"text_size,omitempty"`

	// ShowTextSizeControl controls whether the default theme renders the
	// visitor-facing reading-size selector. It defaults to true.
	ShowTextSizeControl *bool `json:"show_text_size_control,omitempty" yaml:"show_text_size_control,omitempty" toml:"show_text_size_control,omitempty"`

	// SeedColor is the hex color used to generate a triadic palette if Palette == "generated"
	SeedColor string `json:"seed_color,omitempty" yaml:"seed_color,omitempty" toml:"seed_color,omitempty"`

	// Variables allows overriding specific CSS variables
	Variables map[string]string `json:"variables" yaml:"variables" toml:"variables"`

	// CustomCSS is a path to a custom CSS file to include
	CustomCSS string `json:"custom_css" yaml:"custom_css" toml:"custom_css"`

	// Background configures multi-layered background decorations
	Background BackgroundConfig `json:"background,omitempty" yaml:"background,omitempty" toml:"background,omitempty"`

	// Font configures typography settings
	Font FontConfig `json:"font,omitempty" yaml:"font,omitempty" toml:"font,omitempty"`

	// Switcher configures the multi-palette theme switcher dropdown
	Switcher       ThemeSwitcherConfig       `json:"switcher,omitempty" yaml:"switcher,omitempty" toml:"switcher,omitempty"`
	Texture        ThemeTextureConfig        `json:"texture,omitempty" yaml:"texture,omitempty" toml:"texture,omitempty"`
	HeadingTexture ThemeHeadingTextureConfig `json:"heading_texture,omitempty" yaml:"heading_texture,omitempty" toml:"heading_texture,omitempty"`
	Motif          ThemeMotifConfig          `json:"motif,omitempty" yaml:"motif,omitempty" toml:"motif,omitempty"`
}

const (
	TextSizeSmall  = "small"
	TextSizeMedium = "medium"
	TextSizeLarge  = "large"
	TextSizeXLarge = "x-large"
)

// EffectiveTextSize returns a supported reading-size preset, falling back to
// the large preset when the configuration is empty or invalid.
func (c *ThemeConfig) EffectiveTextSize() string {
	if c == nil {
		return TextSizeLarge
	}

	switch c.TextSize {
	case TextSizeSmall, TextSizeMedium, TextSizeLarge, TextSizeXLarge:
		return c.TextSize
	default:
		return TextSizeLarge
	}
}

// IsTextSizeControlEnabled reports whether the visitor-facing reading-size
// selector should be rendered. It defaults to true.
func (c *ThemeConfig) IsTextSizeControlEnabled() bool {
	if c == nil || c.ShowTextSizeControl == nil {
		return true
	}
	return *c.ShowTextSizeControl
}

type ThemeTextureConfig struct {
	Kind        string  `json:"kind,omitempty" yaml:"kind,omitempty" toml:"kind,omitempty"`
	ColorMix    float64 `json:"color_mix,omitempty" yaml:"color_mix,omitempty" toml:"color_mix,omitempty"`
	Scale       float64 `json:"scale,omitempty" yaml:"scale,omitempty" toml:"scale,omitempty"`
	Scope       string  `json:"scope,omitempty" yaml:"scope,omitempty" toml:"scope,omitempty"`
	ColorMixSet bool    `json:"-" yaml:"-" toml:"-"`
	ScaleSet    bool    `json:"-" yaml:"-" toml:"-"`
}

type ThemeHeadingTextureConfig struct {
	Kind        string  `json:"kind,omitempty" yaml:"kind,omitempty" toml:"kind,omitempty"`
	ColorMix    float64 `json:"color_mix,omitempty" yaml:"color_mix,omitempty" toml:"color_mix,omitempty"`
	Scale       float64 `json:"scale,omitempty" yaml:"scale,omitempty" toml:"scale,omitempty"`
	ColorMixSet bool    `json:"-" yaml:"-" toml:"-"`
	ScaleSet    bool    `json:"-" yaml:"-" toml:"-"`
}

type ThemeMotifConfig struct {
	Kind         string  `json:"kind,omitempty" yaml:"kind,omitempty" toml:"kind,omitempty"`
	Glyph        string  `json:"glyph,omitempty" yaml:"glyph,omitempty" toml:"glyph,omitempty"`
	Size         string  `json:"size,omitempty" yaml:"size,omitempty" toml:"size,omitempty"`
	Gap          string  `json:"gap,omitempty" yaml:"gap,omitempty" toml:"gap,omitempty"`
	RowOffset    float64 `json:"row_offset,omitempty" yaml:"row_offset,omitempty" toml:"row_offset,omitempty"`
	Wobble       float64 `json:"wobble,omitempty" yaml:"wobble,omitempty" toml:"wobble,omitempty"`
	Scatter      float64 `json:"scatter,omitempty" yaml:"scatter,omitempty" toml:"scatter,omitempty"`
	Layer        string  `json:"layer,omitempty" yaml:"layer,omitempty" toml:"layer,omitempty"`
	Color        string  `json:"color,omitempty" yaml:"color,omitempty" toml:"color,omitempty"`
	ColorMix     float64 `json:"color_mix,omitempty" yaml:"color_mix,omitempty" toml:"color_mix,omitempty"`
	URL          string  `json:"url,omitempty" yaml:"url,omitempty" toml:"url,omitempty"`
	RowOffsetSet bool    `json:"-" yaml:"-" toml:"-"`
	WobbleSet    bool    `json:"-" yaml:"-" toml:"-"`
	ScatterSet   bool    `json:"-" yaml:"-" toml:"-"`
	ColorMixSet  bool    `json:"-" yaml:"-" toml:"-"`
}

// MarkThemeNumericPresence records numeric theme fields present in source config.
// It is used by config parsers because zero is a valid dial value.
func (c *ThemeConfig) MarkThemeNumericPresence(theme map[string]any) {
	if nested, ok := theme["theme"].(map[string]any); ok {
		theme = nested
	}
	if texture, ok := theme["texture"].(map[string]any); ok {
		_, c.Texture.ColorMixSet = texture["color_mix"]
		_, c.Texture.ScaleSet = texture["scale"]
	}
	if heading, ok := theme["heading_texture"].(map[string]any); ok {
		_, c.HeadingTexture.ColorMixSet = heading["color_mix"]
		_, c.HeadingTexture.ScaleSet = heading["scale"]
	}
	if motif, ok := theme["motif"].(map[string]any); ok {
		_, c.Motif.RowOffsetSet = motif["row_offset"]
		_, c.Motif.WobbleSet = motif["wobble"]
		_, c.Motif.ScatterSet = motif["scatter"]
		_, c.Motif.ColorMixSet = motif["color_mix"]
	}
}

// ThemeSwitcherConfig configures the visitor theme picker.
// When enabled (the default), visitors can preview and select any available
// palette at runtime in the browser.
type ThemeSwitcherConfig struct {
	// Enabled controls whether the visitor theme picker is shown (default: true)
	Enabled *bool `json:"enabled,omitempty" yaml:"enabled,omitempty" toml:"enabled,omitempty"`

	// ModeToggle controls whether the dark/light mode toggle is shown (default: true)
	ModeToggle *bool `json:"mode_toggle,omitempty" yaml:"mode_toggle,omitempty" toml:"mode_toggle,omitempty"`

	// IncludeAll includes all discovered palettes in the switcher (default: true)
	// When false, only palettes in the Include list are shown
	IncludeAll *bool `json:"include_all,omitempty" yaml:"include_all,omitempty" toml:"include_all,omitempty"`

	// Include is a list of palette names to include in the switcher
	// Only used when IncludeAll is false
	Include []string `json:"include,omitempty" yaml:"include,omitempty" toml:"include,omitempty"`

	// Exclude is a list of palette names to exclude from the switcher
	// Used when IncludeAll is true
	Exclude []string `json:"exclude,omitempty" yaml:"exclude,omitempty" toml:"exclude,omitempty"`

	// Position controls where the switcher appears: "header", "footer" (default: "header")
	Position string `json:"position,omitempty" yaml:"position,omitempty" toml:"position,omitempty"`
}

// NewThemeSwitcherConfig creates a new ThemeSwitcherConfig with default values.
func NewThemeSwitcherConfig() ThemeSwitcherConfig {
	enabled := true
	modeToggle := true
	includeAll := true
	return ThemeSwitcherConfig{
		Enabled:    &enabled,
		ModeToggle: &modeToggle,
		IncludeAll: &includeAll,
		Include:    []string{},
		Exclude:    []string{},
		Position:   "header",
	}
}

// IsEnabled returns whether the visitor theme picker is enabled.
// Defaults to true if not explicitly set.
func (s *ThemeSwitcherConfig) IsEnabled() bool {
	if s.Enabled == nil {
		return true
	}
	return *s.Enabled
}

// IsIncludeAll returns whether all palettes should be included.
// Defaults to true if not explicitly set.
func (s *ThemeSwitcherConfig) IsIncludeAll() bool {
	if s.IncludeAll == nil {
		return true
	}
	return *s.IncludeAll
}

// IsModeToggleEnabled returns whether the dark/light mode toggle is enabled.
// Defaults to true if not explicitly set.
func (s *ThemeSwitcherConfig) IsModeToggleEnabled() bool {
	if s.ModeToggle == nil {
		return true
	}
	return *s.ModeToggle
}

// ThemeCalendarConfig configures automatic theme switching based on date ranges.
// This enables seasonal themes, holiday themes, and event-specific styling.
type ThemeCalendarConfig struct {
	// Enabled controls whether the theme calendar is active (default: false)
	Enabled *bool `json:"enabled,omitempty" yaml:"enabled,omitempty" toml:"enabled,omitempty"`

	// Rules is the list of date-based theme rules
	Rules []ThemeCalendarRule `json:"rules,omitempty" yaml:"rules,omitempty" toml:"rules,omitempty"`

	// DefaultPalette is the fallback palette when no rules match (optional)
	// If not set, uses the base theme.palette value
	DefaultPalette string `json:"default_palette,omitempty" yaml:"default_palette,omitempty" toml:"default_palette,omitempty"`
}

// ThemeCalendarRule defines a date range and theme overrides for that period.
// Rules are matched in order - the first matching rule is applied.
type ThemeCalendarRule struct {
	// Name is a descriptive name for the rule (e.g., "Christmas Season", "Winter Frost")
	Name string `json:"name" yaml:"name" toml:"name"`

	// StartDate is the start of the date range in MM-DD format (e.g., "12-15")
	StartDate string `json:"start_date" yaml:"start_date" toml:"start_date"`

	// EndDate is the end of the date range in MM-DD format (e.g., "12-26")
	// Ranges can cross year boundaries (e.g., start="12-01", end="02-28")
	EndDate string `json:"end_date" yaml:"end_date" toml:"end_date"`

	// Palette overrides theme.palette for this period (optional)
	Palette string `json:"palette,omitempty" yaml:"palette,omitempty" toml:"palette,omitempty"`

	// PaletteLight overrides theme.palette_light for this period (optional)
	PaletteLight string `json:"palette_light,omitempty" yaml:"palette_light,omitempty" toml:"palette_light,omitempty"`

	// PaletteDark overrides theme.palette_dark for this period (optional)
	PaletteDark string `json:"palette_dark,omitempty" yaml:"palette_dark,omitempty" toml:"palette_dark,omitempty"`

	// Background overrides theme.background for this period (optional)
	Background *BackgroundConfig `json:"background,omitempty" yaml:"background,omitempty" toml:"background,omitempty"`

	// Font overrides theme.font for this period (optional)
	Font *FontConfig `json:"font,omitempty" yaml:"font,omitempty" toml:"font,omitempty"`

	// Variables merges with theme.variables for this period (optional)
	// These are deep-merged with the base theme variables
	Variables map[string]string `json:"variables,omitempty" yaml:"variables,omitempty" toml:"variables,omitempty"`

	// CustomCSS overrides theme.custom_css for this period (optional)
	CustomCSS string `json:"custom_css,omitempty" yaml:"custom_css,omitempty" toml:"custom_css,omitempty"`
}

// IsEnabled returns whether the theme calendar is enabled.
// Defaults to false if not explicitly set.
func (c *ThemeCalendarConfig) IsEnabled() bool {
	if c.Enabled == nil {
		return false
	}
	return *c.Enabled
}

// BackgroundConfig configures multi-layered background decorations for pages.
// Background elements are rendered as fixed-position layers behind the main content.
type BackgroundConfig struct {
	// Enabled controls whether background decorations are active (default: false)
	Enabled *bool `json:"enabled,omitempty" yaml:"enabled,omitempty" toml:"enabled,omitempty"`

	// Backgrounds is the list of background elements to render
	Backgrounds []BackgroundElement `json:"backgrounds,omitempty" yaml:"backgrounds,omitempty" toml:"backgrounds,omitempty"`

	// Scripts is a list of script URLs to include for background functionality
	// Example: ["/static/js/snow-fall.js"]
	Scripts []string `json:"scripts,omitempty" yaml:"scripts,omitempty" toml:"scripts,omitempty"`

	// CSS is custom CSS for styling background elements
	CSS string `json:"css,omitempty" yaml:"css,omitempty" toml:"css,omitempty"`

	// ArticleBg is the background color for article/content areas (default: uses --color-background)
	// This helps ensure content is readable over decorative backgrounds.
	// Example: "rgba(255, 255, 255, 0.95)" or "#ffffff"
	ArticleBg string `json:"article_bg,omitempty" yaml:"article_bg,omitempty" toml:"article_bg,omitempty"`

	// ArticleBlurEnabled controls whether backdrop blur is applied to article areas (default: false)
	// When false, the ArticleBlur value is ignored even if set.
	ArticleBlurEnabled *bool `json:"article_blur_enabled,omitempty" yaml:"article_blur_enabled,omitempty" toml:"article_blur_enabled,omitempty"`

	// ArticleBlur is the backdrop blur amount for article areas (default: "0px")
	// Example: "8px" or "12px" for a frosted glass effect
	ArticleBlur string `json:"article_blur,omitempty" yaml:"article_blur,omitempty" toml:"article_blur,omitempty"`

	// ArticleShadow is the box-shadow for article areas
	// Example: "0 4px 20px rgba(0, 0, 0, 0.3)"
	ArticleShadow string `json:"article_shadow,omitempty" yaml:"article_shadow,omitempty" toml:"article_shadow,omitempty"`

	// ArticleBorder is the border style for article areas
	// Example: "1px solid rgba(255, 255, 255, 0.1)"
	ArticleBorder string `json:"article_border,omitempty" yaml:"article_border,omitempty" toml:"article_border,omitempty"`

	// ArticleRadius is the border-radius for article areas (default: uses --radius-lg)
	// Example: "12px" or "1rem"
	ArticleRadius string `json:"article_radius,omitempty" yaml:"article_radius,omitempty" toml:"article_radius,omitempty"`
}

// BackgroundElement represents a single background decoration layer.
type BackgroundElement struct {
	// HTML is the HTML content for this background layer
	// Example: '<snow-fall count="200"></snow-fall>'
	HTML string `json:"html" yaml:"html" toml:"html"`

	// ZIndex controls the stacking order of this layer (default: -1)
	// Negative values place the layer behind content, positive values in front
	ZIndex int `json:"z_index,omitempty" yaml:"z_index,omitempty" toml:"z_index,omitempty"`
}

// NewBackgroundConfig creates a new BackgroundConfig with default values.
func NewBackgroundConfig() BackgroundConfig {
	enabled := false
	articleBlurEnabled := false
	return BackgroundConfig{
		Enabled:            &enabled,
		Backgrounds:        []BackgroundElement{},
		Scripts:            []string{},
		CSS:                "",
		ArticleBg:          "",
		ArticleBlurEnabled: &articleBlurEnabled,
		ArticleBlur:        "",
		ArticleShadow:      "",
		ArticleBorder:      "",
		ArticleRadius:      "",
	}
}

// IsEnabled returns whether background decorations are enabled.
// Defaults to false if not explicitly set.
func (b *BackgroundConfig) IsEnabled() bool {
	if b.Enabled == nil {
		return false
	}
	return *b.Enabled
}

// IsArticleBlurEnabled returns whether backdrop blur is enabled for article areas.
// Defaults to false if not explicitly set.
func (b *BackgroundConfig) IsArticleBlurEnabled() bool {
	if b.ArticleBlurEnabled == nil {
		return false
	}
	return *b.ArticleBlurEnabled
}

// GlobConfig configures file globbing behavior.
type GlobConfig struct {
	// Patterns is the list of glob patterns to match source files
	Patterns []string `json:"patterns" yaml:"patterns" toml:"patterns"`

	// UseGitignore determines whether to respect .gitignore files
	UseGitignore bool `json:"use_gitignore" yaml:"use_gitignore" toml:"use_gitignore"`

	// SlugMode controls how slugs are derived from matched source paths.
	// Valid values: "path", "flat".
	SlugMode string `json:"slug_mode" yaml:"slug_mode" toml:"slug_mode"`

	// SlugRules override slug mode for specific content path prefixes.
	SlugRules []SlugRule `json:"slug_rules,omitempty" yaml:"slug_rules,omitempty" toml:"slug_rules,omitempty"`
}

// SlugRule selects a slug mode for files under a content path prefix.
type SlugRule struct {
	Prefix string `json:"prefix" yaml:"prefix" toml:"prefix"`
	Mode   string `json:"mode" yaml:"mode" toml:"mode"`
}

// MarkdownConfig configures markdown processing.
type MarkdownConfig struct {
	// Extensions is the list of markdown extensions to enable
	Extensions []string `json:"extensions" yaml:"extensions" toml:"extensions"`

	// Highlight configures syntax highlighting for code blocks
	Highlight HighlightConfig `json:"highlight" yaml:"highlight" toml:"highlight"`
}

// HighlightConfig configures syntax highlighting for code blocks.
type HighlightConfig struct {
	// Enabled controls whether syntax highlighting is active (default: true)
	Enabled *bool `json:"enabled,omitempty" yaml:"enabled,omitempty" toml:"enabled,omitempty"`

	// Theme is the Chroma theme to use for syntax highlighting.
	// If empty, the theme is automatically derived from the site's color palette.
	// See https://xyproto.github.io/splash/docs/ for available themes.
	Theme string `json:"theme,omitempty" yaml:"theme,omitempty" toml:"theme,omitempty"`

	// LineNumbers enables line numbers in code blocks (default: false)
	LineNumbers bool `json:"line_numbers" yaml:"line_numbers" toml:"line_numbers"`
}

// NewHighlightConfig creates a new HighlightConfig with default values.
func NewHighlightConfig() HighlightConfig {
	enabled := true
	return HighlightConfig{
		Enabled:     &enabled,
		Theme:       "", // Empty means auto-detect from palette
		LineNumbers: false,
	}
}

// IsEnabled returns whether syntax highlighting is enabled.
// Defaults to true if not explicitly set.
func (h *HighlightConfig) IsEnabled() bool {
	if h.Enabled == nil {
		return true
	}
	return *h.Enabled
}

// CSVFenceConfig configures the csv_fence plugin.
type CSVFenceConfig struct {
	// Enabled controls whether CSV blocks are converted to tables (default: true)
	Enabled bool `json:"enabled" yaml:"enabled" toml:"enabled"`

	// TableClass is the CSS class for generated tables (default: "csv-table")
	TableClass string `json:"table_class" yaml:"table_class" toml:"table_class"`

	// HasHeader indicates whether the first row is a header (default: true)
	HasHeader bool `json:"has_header" yaml:"has_header" toml:"has_header"`

	// Delimiter is the CSV field delimiter (default: ",")
	Delimiter string `json:"delimiter" yaml:"delimiter" toml:"delimiter"`
}

// NewCSVFenceConfig creates a new CSVFenceConfig with default values.
func NewCSVFenceConfig() CSVFenceConfig {
	return CSVFenceConfig{
		Enabled:    true,
		TableClass: "csv-table",
		HasHeader:  true,
		Delimiter:  ",",
	}
}

// MermaidConfig configures the mermaid plugin.
type MermaidConfig struct {
	// Enabled controls whether mermaid processing is active (default: true)
	Enabled bool `json:"enabled" yaml:"enabled" toml:"enabled"`

	// Mode specifies how to render diagrams: "client" (browser), "cli" (mmdc), or "chromium" (Chrome/Chromium)
	// Default: "client"
	Mode string `json:"mode" yaml:"mode" toml:"mode"`

	// CDNURL is the URL for the Mermaid.js library (client mode only)
	CDNURL string `json:"cdn_url" yaml:"cdn_url" toml:"cdn_url"`

	// Theme is the Mermaid theme to use (default, dark, forest, neutral)
	Theme string `json:"theme" yaml:"theme" toml:"theme"`

	// UseCSSVariables enables palette-aware theme variables from CSS custom properties.
	UseCSSVariables bool `json:"use_css_variables" yaml:"use_css_variables" toml:"use_css_variables"`

	// Lightbox enables GLightbox zoom for rendered Mermaid diagrams.
	Lightbox bool `json:"lightbox" yaml:"lightbox" toml:"lightbox"`

	// LightboxSelector is the CSS selector used for Mermaid lightbox links.
	LightboxSelector string `json:"lightbox_selector" yaml:"lightbox_selector" toml:"lightbox_selector"`

	// CLIConfig contains settings for CLI mode (npm mmdc)
	CLIConfig *CLIRendererConfig `json:"cli" yaml:"cli" toml:"cli"`

	// ChromiumConfig contains settings for Chromium mode.
	ChromiumConfig *ChromiumRendererConfig `json:"chromium" yaml:"chromium" toml:"chromium"`
}

// CLIRendererConfig configures the CLI-based renderer (npm mmdc)
type CLIRendererConfig struct {
	// MMDCPath is the path to the mmdc binary. If empty, looks for it in PATH.
	MMDCPath string `json:"mmdc_path" yaml:"mmdc_path" toml:"mmdc_path"`

	// ExtraArgs are additional command-line arguments passed to mmdc
	ExtraArgs string `json:"extra_args" yaml:"extra_args" toml:"extra_args"`
}

// ChromiumRendererConfig configures the Chromium-based renderer.
type ChromiumRendererConfig struct {
	// BrowserPath is the path to the Chrome/Chromium binary. If empty, auto-detects.
	BrowserPath string `json:"browser_path" yaml:"browser_path" toml:"browser_path"`

	// Timeout is the maximum time (in seconds) to wait for a diagram to render
	Timeout int `json:"timeout" yaml:"timeout" toml:"timeout"`

	// MaxConcurrent is the maximum number of concurrent diagram renders
	MaxConcurrent int `json:"max_concurrent" yaml:"max_concurrent" toml:"max_concurrent"`

	// NoSandbox disables the Chromium sandbox. Required in containers (Docker, Distrobox, etc.)
	NoSandbox bool `json:"no_sandbox" yaml:"no_sandbox" toml:"no_sandbox"`
}

// NewMermaidConfig creates a new MermaidConfig with default values.
func NewMermaidConfig() MermaidConfig {
	return MermaidConfig{
		Enabled:          true,
		Mode:             "client",
		CDNURL:           "/assets/vendor/mermaid/mermaid.min.js",
		Theme:            "default",
		UseCSSVariables:  true,
		Lightbox:         true,
		LightboxSelector: ".glightbox-mermaid",
		CLIConfig: &CLIRendererConfig{
			MMDCPath:  "",
			ExtraArgs: "",
		},
		ChromiumConfig: &ChromiumRendererConfig{
			BrowserPath:   "",
			Timeout:       30,
			MaxConcurrent: 4,
		},
	}
}

// MDVideoConfig configures the md_video plugin.
type MDVideoConfig struct {
	// Enabled controls whether video conversion is active (default: true)
	Enabled bool `json:"enabled" yaml:"enabled" toml:"enabled"`

	// VideoExtensions is the list of file extensions to treat as videos
	VideoExtensions []string `json:"video_extensions" yaml:"video_extensions" toml:"video_extensions"`

	// VideoClass is the CSS class added to video elements (default: "md-video")
	VideoClass string `json:"video_class" yaml:"video_class" toml:"video_class"`

	// Controls shows video controls (default: true)
	Controls bool `json:"controls" yaml:"controls" toml:"controls"`

	// Autoplay records authored playback intent; the client loading policy decides
	// whether unrequested video may start automatically (default: false).
	Autoplay bool `json:"autoplay" yaml:"autoplay" toml:"autoplay"`

	// Loop repeats the video (default: true for GIF-like behavior)
	Loop bool `json:"loop" yaml:"loop" toml:"loop"`

	// Muted mutes the video (default: true, required for autoplay in most browsers)
	Muted bool `json:"muted" yaml:"muted" toml:"muted"`

	// Playsinline enables inline playback on mobile (default: true)
	Playsinline bool `json:"playsinline" yaml:"playsinline" toml:"playsinline"`

	// Preload hints how much to preload: "none", "metadata", "auto" (default: "none")
	Preload string `json:"preload" yaml:"preload" toml:"preload"`
}

// NewMDVideoConfig creates a new MDVideoConfig with sensible defaults.
// Videos show their poster and begin loading only after playback is requested.
func NewMDVideoConfig() MDVideoConfig {
	return MDVideoConfig{
		Enabled:         true,
		VideoExtensions: []string{".mp4", ".webm", ".ogg", ".ogv", ".mov", ".m4v", ".avi", ".mkv"},
		VideoClass:      "md-video",
		Controls:        true,
		Autoplay:        false,
		Loop:            true,
		Muted:           true,
		Playsinline:     true,
		Preload:         "none",
	}
}

// ChartJSConfig configures the chartjs plugin.
type ChartJSConfig struct {
	// Enabled controls whether Chart.js processing is active (default: true)
	Enabled bool `json:"enabled" yaml:"enabled" toml:"enabled"`

	// CDNURL is the URL for the Chart.js library
	CDNURL string `json:"cdn_url" yaml:"cdn_url" toml:"cdn_url"`

	// ContainerClass is the CSS class for the chart container div (default: "chartjs-container")
	ContainerClass string `json:"container_class" yaml:"container_class" toml:"container_class"`

	// DataDisplay sets the initial disclosure state: "open" (default) or "collapsed".
	// An empty value also means open, preserving programmatic zero-value configs.
	DataDisplay string `json:"data_display" yaml:"data_display" toml:"data_display"`
}

// NewChartJSConfig creates a new ChartJSConfig with default values.
func NewChartJSConfig() ChartJSConfig {
	return ChartJSConfig{
		Enabled:        true,
		CDNURL:         "/assets/vendor/chartjs/chart.min.js",
		ContainerClass: "chartjs-container",
		DataDisplay:    "open",
	}
}

// ContributionGraphConfig configures the contribution_graph plugin.
type ContributionGraphConfig struct {
	// Enabled controls whether contribution graph processing is active (default: true)
	Enabled bool `json:"enabled" yaml:"enabled" toml:"enabled"`

	// CDNURL is the URL for the Cal-Heatmap library
	CDNURL string `json:"cdn_url" yaml:"cdn_url" toml:"cdn_url"`

	// ContainerClass is the CSS class for the graph container div (default: "contribution-graph-container")
	ContainerClass string `json:"container_class" yaml:"container_class" toml:"container_class"`

	// Theme is the Cal-Heatmap color theme (default: "light")
	Theme string `json:"theme" yaml:"theme" toml:"theme"`

	// ScaleMaxPercentile caps the color scale at this percentile when set to a value between 0 and 100.
	// A value of 0 disables percentile-based outlier control.
	ScaleMaxPercentile float64 `json:"scale_max_percentile,omitempty" yaml:"scale_max_percentile,omitempty" toml:"scale_max_percentile,omitempty"`
}

// NewContributionGraphConfig creates a new ContributionGraphConfig with default values.
func NewContributionGraphConfig() ContributionGraphConfig {
	return ContributionGraphConfig{
		Enabled:            true,
		CDNURL:             "/assets/vendor/cal-heatmap",
		ContainerClass:     "contribution-graph-container",
		Theme:              "light",
		ScaleMaxPercentile: 0,
	}
}

// OneLineLinkConfig configures the one_line_link plugin.
type OneLineLinkConfig struct {
	// Enabled controls whether link expansion is active (default: true)
	Enabled bool `json:"enabled" yaml:"enabled" toml:"enabled"`

	// CardClass is the CSS class for the link card (default: "link-card")
	CardClass string `json:"card_class" yaml:"card_class" toml:"card_class"`

	// FetchMetadata enables fetching title/description from URLs (default: false for performance)
	FetchMetadata bool `json:"fetch_metadata" yaml:"fetch_metadata" toml:"fetch_metadata"`

	// FallbackTitle is used when metadata fetch fails (default: "Link")
	FallbackTitle string `json:"fallback_title" yaml:"fallback_title" toml:"fallback_title"`

	// Timeout is the HTTP request timeout in seconds (default: 5)
	Timeout int `json:"timeout" yaml:"timeout" toml:"timeout"`

	// ExcludePatterns is a list of regex patterns for URLs to exclude from expansion
	ExcludePatterns []string `json:"exclude_patterns" yaml:"exclude_patterns" toml:"exclude_patterns"`
}

// NewOneLineLinkConfig creates a new OneLineLinkConfig with default values.
func NewOneLineLinkConfig() OneLineLinkConfig {
	return OneLineLinkConfig{
		Enabled:         true,
		CardClass:       "link-card",
		FetchMetadata:   false, // Disabled by default for build performance
		FallbackTitle:   "Link",
		Timeout:         5,
		ExcludePatterns: []string{},
	}
}

// WikilinkHoverConfig configures the wikilink_hover plugin.
type WikilinkHoverConfig struct {
	// Enabled controls whether hover previews are added (default: true)
	Enabled bool `json:"enabled" yaml:"enabled" toml:"enabled"`

	// PreviewLength is the maximum characters for the preview text (default: 200)
	PreviewLength int `json:"preview_length" yaml:"preview_length" toml:"preview_length"`

	// IncludeImage adds data-preview-image attribute if post has a featured image (default: true)
	IncludeImage bool `json:"include_image" yaml:"include_image" toml:"include_image"`

	// ScreenshotService is an optional URL prefix for screenshot generation
	// If set, adds data-preview-screenshot attribute with the URL
	ScreenshotService string `json:"screenshot_service" yaml:"screenshot_service" toml:"screenshot_service"`
}

// NewWikilinkHoverConfig creates a new WikilinkHoverConfig with default values.
func NewWikilinkHoverConfig() WikilinkHoverConfig {
	return WikilinkHoverConfig{
		Enabled:           true,
		PreviewLength:     200,
		IncludeImage:      true,
		ScreenshotService: "",
	}
}

// SEOConfig configures SEO metadata for the site.
type SEOConfig struct {
	// TwitterHandle is the Twitter/X username (without @) for twitter:site meta tag
	TwitterHandle string `json:"twitter_handle" yaml:"twitter_handle" toml:"twitter_handle"`

	// DefaultImage is the default Open Graph image URL for pages without a specific image
	DefaultImage string `json:"default_image" yaml:"default_image" toml:"default_image"`

	// LogoURL is the site logo URL for Schema.org structured data
	LogoURL string `json:"logo_url" yaml:"logo_url" toml:"logo_url"`

	// AuthorImage is the author's profile image URL for OG cards
	AuthorImage string `json:"author_image" yaml:"author_image" toml:"author_image"`

	// OGImageService is the URL for a screenshot service that generates OG images
	// from OG card pages. The URL should accept a `url` query parameter.
	// Example: "https://shots.example.com/shot/" generates URLs like:
	// "https://shots.example.com/shot/?url=https://site.com/post/og/&height=600&width=1200&format=jpg"
	OGImageService string `json:"og_image_service" yaml:"og_image_service" toml:"og_image_service"`

	// StructuredData configures JSON-LD Schema.org generation
	StructuredData StructuredDataConfig `json:"structured_data" yaml:"structured_data" toml:"structured_data"`
}

// StructuredDataConfig configures JSON-LD Schema.org structured data generation.
type StructuredDataConfig struct {
	// Enabled controls whether structured data is generated (default: true)
	Enabled *bool `json:"enabled,omitempty" yaml:"enabled,omitempty" toml:"enabled,omitempty"`

	// Publisher is the site publisher information for Schema.org
	Publisher *EntityConfig `json:"publisher,omitempty" yaml:"publisher,omitempty" toml:"publisher,omitempty"`

	// DefaultAuthor is the default author for posts without explicit author
	DefaultAuthor *EntityConfig `json:"default_author,omitempty" yaml:"default_author,omitempty" toml:"default_author,omitempty"`
}

// IsEnabled returns whether structured data generation is enabled.
// Defaults to true if not explicitly set.
func (s *StructuredDataConfig) IsEnabled() bool {
	if s.Enabled == nil {
		return true
	}
	return *s.Enabled
}

// EntityConfig represents a Schema.org Person or Organization entity.
type EntityConfig struct {
	// Type is "Person" or "Organization" (default: "Organization")
	Type string `json:"type" yaml:"type" toml:"type"`

	// Name is the entity name
	Name string `json:"name" yaml:"name" toml:"name"`

	// URL is the entity's web page
	URL string `json:"url,omitempty" yaml:"url,omitempty" toml:"url,omitempty"`

	// Logo is the logo URL (for Organizations only)
	Logo string `json:"logo,omitempty" yaml:"logo,omitempty" toml:"logo,omitempty"`
}

// NewStructuredDataConfig creates a new StructuredDataConfig with default values.
func NewStructuredDataConfig() StructuredDataConfig {
	enabled := true
	return StructuredDataConfig{
		Enabled: &enabled,
	}
}

// NewSEOConfig creates a new SEOConfig with default values.
func NewSEOConfig() SEOConfig {
	return SEOConfig{
		TwitterHandle:  "",
		DefaultImage:   "",
		LogoURL:        "",
		StructuredData: NewStructuredDataConfig(),
	}
}

// IndieAuthConfig configures IndieAuth link tags for identity and authentication.
// IndieAuth is a decentralized authentication protocol built on OAuth 2.0.
// See https://indieauth.spec.indieweb.org/ for the specification.
type IndieAuthConfig struct {
	// Enabled controls whether IndieAuth link tags are included (default: false)
	Enabled bool `json:"enabled" yaml:"enabled" toml:"enabled"`

	// AuthorizationEndpoint is the URL of your authorization endpoint
	// Example: "https://indieauth.com/auth"
	AuthorizationEndpoint string `json:"authorization_endpoint" yaml:"authorization_endpoint" toml:"authorization_endpoint"`

	// TokenEndpoint is the URL of your token endpoint
	// Example: "https://tokens.indieauth.com/token"
	TokenEndpoint string `json:"token_endpoint" yaml:"token_endpoint" toml:"token_endpoint"`

	// MeURL is your profile URL for rel="me" links (optional)
	// This links your site to other profiles (GitHub, Twitter, etc.)
	MeURL string `json:"me_url" yaml:"me_url" toml:"me_url"`
}

// NewIndieAuthConfig creates a new IndieAuthConfig with default values.
func NewIndieAuthConfig() IndieAuthConfig {
	return IndieAuthConfig{
		Enabled:               false,
		AuthorizationEndpoint: "",
		TokenEndpoint:         "",
		MeURL:                 "",
	}
}

// WebmentionConfig configures Webmention endpoint for receiving mentions.
// Webmention is a simple protocol for notifying URLs when you link to them.
// See https://www.w3.org/TR/webmention/ for the specification.
type WebmentionConfig struct {
	// Enabled controls whether Webmention link tag is included (default: false)
	Enabled bool `json:"enabled" yaml:"enabled" toml:"enabled"`

	// Endpoint is the URL of your Webmention endpoint
	// Example: "https://webmention.io/example.com/webmention"
	Endpoint string `json:"endpoint" yaml:"endpoint" toml:"endpoint"`
}

// NewWebmentionConfig creates a new WebmentionConfig with default values.
func NewWebmentionConfig() WebmentionConfig {
	return WebmentionConfig{
		Enabled:  false,
		Endpoint: "",
	}
}

// WebSubConfig configures WebSub discovery links for feeds.
// See https://www.w3.org/TR/websub/ for the specification.
type WebSubConfig struct {
	// Enabled controls whether WebSub discovery links are included (default: false)
	Enabled *bool `json:"enabled,omitempty" yaml:"enabled,omitempty" toml:"enabled,omitempty"`

	// Hubs is the list of WebSub hub URLs to advertise
	Hubs []string `json:"hubs,omitempty" yaml:"hubs,omitempty" toml:"hubs,omitempty"`
}

// NewWebSubConfig creates a new WebSubConfig with default values.
func NewWebSubConfig() WebSubConfig {
	enabled := false
	return WebSubConfig{
		Enabled: &enabled,
		Hubs:    []string{},
	}
}

// IsEnabled returns whether WebSub discovery links are enabled.
// Defaults to false if not explicitly set.
func (w WebSubConfig) IsEnabled() bool {
	if w.Enabled == nil {
		return false
	}
	return *w.Enabled
}

// HubsList returns a copy of configured hub URLs.
func (w WebSubConfig) HubsList() []string {
	return append([]string{}, w.Hubs...)
}

// WebMentionsConfig configures the webmentions plugin for sending outgoing mentions.
// This is separate from WebmentionConfig which handles receiving mentions.
type WebMentionsConfig struct {
	// Enabled controls whether the webmentions plugin is active (default: false)
	Enabled bool `json:"enabled" yaml:"enabled" toml:"enabled"`

	// Outgoing enables sending webmentions for external links (default: true when enabled)
	Outgoing bool `json:"outgoing" yaml:"outgoing" toml:"outgoing"`

	// UserAgent is the User-Agent string for HTTP requests
	UserAgent string `json:"user_agent" yaml:"user_agent" toml:"user_agent"`

	// Timeout is the HTTP request timeout (e.g., "30s")
	Timeout string `json:"timeout" yaml:"timeout" toml:"timeout"`

	// CacheDir is the directory for caching sent webmentions (default: ".cache/webmentions")
	CacheDir string `json:"cache_dir" yaml:"cache_dir" toml:"cache_dir"`

	// ConcurrentRequests is the max number of concurrent webmention requests (default: 5)
	ConcurrentRequests int `json:"concurrent_requests" yaml:"concurrent_requests" toml:"concurrent_requests"`

	// Bridges configures social media bridging for incoming webmentions
	Bridges BridgesConfig `json:"bridges" yaml:"bridges" toml:"bridges"`

	// WebmentionIOToken is the API token for webmention.io (for receiving mentions)
	WebmentionIOToken string `json:"webmention_io_token" yaml:"webmention_io_token" toml:"webmention_io_token"`
}

// BridgesConfig configures social media bridging services.
type BridgesConfig struct {
	// Enabled controls whether bridging detection is active (default: false)
	Enabled bool `json:"enabled" yaml:"enabled" toml:"enabled"`

	// BridgyFediverse enables Bridgy Fed integration (default: true when bridges enabled)
	BridgyFediverse bool `json:"bridgy_fediverse" yaml:"bridgy_fediverse" toml:"bridgy_fediverse"`

	// Platform-specific controls
	Bluesky  bool `json:"bluesky" yaml:"bluesky" toml:"bluesky"`
	Twitter  bool `json:"twitter" yaml:"twitter" toml:"twitter"`
	Mastodon bool `json:"mastodon" yaml:"mastodon" toml:"mastodon"`
	GitHub   bool `json:"github" yaml:"github" toml:"github"`
	Flickr   bool `json:"flickr" yaml:"flickr" toml:"flickr"`

	// Filters configures filtering of bridged mentions
	Filters BridgeFiltersConfig `json:"filters" yaml:"filters" toml:"filters"`
}

// BridgeFiltersConfig configures filtering for bridged webmentions.
type BridgeFiltersConfig struct {
	// Platforms limits which platforms to accept (empty = all enabled)
	Platforms []string `json:"platforms" yaml:"platforms" toml:"platforms"`

	// InteractionTypes limits which interaction types to accept (empty = all)
	// Valid values: "like", "repost", "reply", "bookmark", "mention"
	InteractionTypes []string `json:"interaction_types" yaml:"interaction_types" toml:"interaction_types"`

	// MinContentLength filters out mentions with content shorter than this
	MinContentLength int `json:"min_content_length" yaml:"min_content_length" toml:"min_content_length"`

	// BlockedDomains is a list of domains to reject mentions from
	BlockedDomains []string `json:"blocked_domains" yaml:"blocked_domains" toml:"blocked_domains"`
}

// NewBridgesConfig creates a new BridgesConfig with default values.
func NewBridgesConfig() BridgesConfig {
	return BridgesConfig{
		Enabled:         false,
		BridgyFediverse: true,
		Bluesky:         true,
		Twitter:         true,
		Mastodon:        true,
		GitHub:          true,
		Flickr:          false,
		Filters:         BridgeFiltersConfig{},
	}
}

// NewWebMentionsConfig creates a new WebMentionsConfig with default values.
func NewWebMentionsConfig() WebMentionsConfig {
	return WebMentionsConfig{
		Enabled:            false,
		Outgoing:           true,
		UserAgent:          "markata-go/1.0 (WebMention Sender; +https://github.com/WaylonWalker/markata)",
		Timeout:            "30s",
		CacheDir:           ".cache/webmentions",
		ConcurrentRequests: 5,
		Bridges:            NewBridgesConfig(),
	}
}

// ResourceHintsConfig configures automatic resource hints generation for network optimization.
// Resource hints (preconnect, dns-prefetch, preload, prefetch) help browsers prepare
// for external resources before they're needed, improving page load performance.
type ResourceHintsConfig struct {
	// Enabled controls whether resource hints are generated (default: true)
	Enabled *bool `json:"enabled,omitempty" yaml:"enabled,omitempty" toml:"enabled,omitempty"`

	// AutoDetect enables automatic detection of external domains in HTML/CSS (default: true)
	AutoDetect *bool `json:"auto_detect,omitempty" yaml:"auto_detect,omitempty" toml:"auto_detect,omitempty"`

	// Domains is a list of manually configured domain hints
	Domains []DomainHint `json:"domains,omitempty" yaml:"domains,omitempty" toml:"domains,omitempty"`

	// ExcludeDomains is a list of domains to exclude from auto-detection
	ExcludeDomains []string `json:"exclude_domains,omitempty" yaml:"exclude_domains,omitempty" toml:"exclude_domains,omitempty"`
}

// DomainHint represents a hint configuration for a specific domain.
type DomainHint struct {
	// Domain is the external domain (e.g., "fonts.googleapis.com")
	Domain string `json:"domain" yaml:"domain" toml:"domain"`

	// HintTypes specifies which hint types to generate for this domain
	// Valid values: "preconnect", "dns-prefetch", "preload", "prefetch"
	HintTypes []string `json:"hint_types" yaml:"hint_types" toml:"hint_types"`

	// CrossOrigin specifies the crossorigin attribute value
	// Valid values: "", "anonymous", "use-credentials" (default: "")
	CrossOrigin string `json:"crossorigin,omitempty" yaml:"crossorigin,omitempty" toml:"crossorigin,omitempty"`

	// As specifies the "as" attribute for preload hints (e.g., "font", "script", "style")
	As string `json:"as,omitempty" yaml:"as,omitempty" toml:"as,omitempty"`
}

// NewResourceHintsConfig creates a new ResourceHintsConfig with default values.
func NewResourceHintsConfig() ResourceHintsConfig {
	enabled := true
	autoDetect := true
	return ResourceHintsConfig{
		Enabled:        &enabled,
		AutoDetect:     &autoDetect,
		Domains:        []DomainHint{},
		ExcludeDomains: []string{},
	}
}

// IsEnabled returns whether resource hints generation is enabled.
// Defaults to true if not explicitly set.
func (r *ResourceHintsConfig) IsEnabled() bool {
	if r.Enabled == nil {
		return true
	}
	return *r.Enabled
}

// IsAutoDetectEnabled returns whether auto-detection is enabled.
// Defaults to true if not explicitly set.
func (r *ResourceHintsConfig) IsAutoDetectEnabled() bool {
	if r.AutoDetect == nil {
		return true
	}
	return *r.AutoDetect
}

// PostFormatsConfig configures the output formats for individual posts.
// This controls what file formats are generated for each post.
type PostFormatsConfig struct {
	// HTML enables standard HTML output (default: true)
	// Generates: /slug/index.html
	HTML *bool `json:"html,omitempty" yaml:"html,omitempty" toml:"html,omitempty"`

	// Markdown enables raw markdown output (default: false)
	// Generates: /slug.md (source with frontmatter)
	Markdown bool `json:"markdown" yaml:"markdown" toml:"markdown"`

	// Text enables plain text output (default: false)
	// Generates: /slug.txt (terminal-friendly plain text)
	Text bool `json:"text" yaml:"text" toml:"text"`

	// ANSI enables ANSI-styled terminal output (default: false)
	// Generates: /slug.ansi (terminal-friendly styled output)
	ANSI bool `json:"ansi" yaml:"ansi" toml:"ansi"`

	// OG enables OpenGraph card HTML output for social image generation (default: false)
	// Generates: /slug/og/index.html (1200x630 optimized for screenshots)
	OG bool `json:"og" yaml:"og" toml:"og"`
}

// WellKnownConfig configures auto-generated .well-known entries.
type WellKnownConfig struct {
	// Enabled controls whether .well-known generation runs (default: true)
	Enabled *bool `json:"enabled,omitempty" yaml:"enabled,omitempty" toml:"enabled,omitempty"`

	// AutoGenerate lists entries to generate from site metadata
	AutoGenerate []string `json:"auto_generate,omitempty" yaml:"auto_generate,omitempty" toml:"auto_generate,omitempty"`

	// SSHFingerprint is written to /.well-known/sshfp if set
	SSHFingerprint string `json:"ssh_fingerprint,omitempty" yaml:"ssh_fingerprint,omitempty" toml:"ssh_fingerprint,omitempty"`

	// KeybaseUsername is written to /.well-known/keybase.txt if set
	KeybaseUsername string `json:"keybase_username,omitempty" yaml:"keybase_username,omitempty" toml:"keybase_username,omitempty"`
}

// DefaultWellKnownAutoGenerate lists the default auto-generated entries.
var DefaultWellKnownAutoGenerate = []string{
	"host-meta",
	"host-meta.json",
	"webfinger",
	"nodeinfo",
	"time",
	"links",
}

// NewWellKnownConfig creates a WellKnownConfig with default values.
func NewWellKnownConfig() WellKnownConfig {
	enabled := true
	return WellKnownConfig{
		Enabled:      &enabled,
		AutoGenerate: append([]string{}, DefaultWellKnownAutoGenerate...),
	}
}

// IsEnabled returns whether .well-known generation is enabled.
// Defaults to true if not explicitly set.
func (w WellKnownConfig) IsEnabled() bool {
	if w.Enabled == nil {
		return true
	}
	return *w.Enabled
}

// AutoGenerateList returns the configured auto-generate list with defaults applied.
func (w WellKnownConfig) AutoGenerateList() []string {
	if w.AutoGenerate == nil {
		return append([]string{}, DefaultWellKnownAutoGenerate...)
	}
	return append([]string{}, w.AutoGenerate...)
}

// TemplatePreset defines templates for all output formats.
// This allows setting all format templates at once with a single preset name.
type TemplatePreset struct {
	// HTML template file for HTML output
	HTML string `json:"html" yaml:"html" toml:"html"`

	// Text template file for txt output
	Text string `json:"txt" yaml:"txt" toml:"txt"`

	// ANSI template file for ansi output
	ANSI string `json:"ansi,omitempty" yaml:"ansi,omitempty" toml:"ansi,omitempty"`

	// Markdown template file for markdown output
	Markdown string `json:"markdown" yaml:"markdown" toml:"markdown"`

	// OG template file for OpenGraph card output
	OG string `json:"og" yaml:"og" toml:"og"`
}

// TemplateForFormat returns the template for a specific format.
// Returns empty string if the format is not recognized.
func (p *TemplatePreset) TemplateForFormat(format string) string {
	switch format {
	case "html":
		return p.HTML
	case "txt", "text":
		return p.Text
	case "ansi":
		return p.ANSI
	case "markdown", "md":
		return p.Markdown
	case "og":
		return p.OG
	default:
		return ""
	}
}

// NewPostFormatsConfig creates a new PostFormatsConfig with default values.
// HTML, markdown, text, and OG are enabled by default. ANSI stays opt-in.
func NewPostFormatsConfig() PostFormatsConfig {
	enabled := true
	return PostFormatsConfig{
		HTML:     &enabled,
		Markdown: true,
		Text:     true,
		ANSI:     false,
		OG:       true,
	}
}

// IsHTMLEnabled returns whether HTML output is enabled.
// Defaults to true if not explicitly set.
func (p *PostFormatsConfig) IsHTMLEnabled() bool {
	if p.HTML == nil {
		return true
	}
	return *p.HTML
}

// QRCodeConfig configures the qrcode plugin.
type QRCodeConfig struct {
	// Enabled controls whether QR codes are generated (default: true)
	Enabled bool `json:"enabled" yaml:"enabled" toml:"enabled"`

	// Format is the output format: "svg" or "png" (default: "svg")
	Format string `json:"format" yaml:"format" toml:"format"`

	// Size is the QR code size in pixels (default: 200)
	Size int `json:"size" yaml:"size" toml:"size"`

	// OutputDir is the subdirectory in output for QR code files (default: "qrcodes")
	OutputDir string `json:"output_dir" yaml:"output_dir" toml:"output_dir"`

	// ErrorCorrection is the QR error correction level: L, M, Q, H (default: "M")
	ErrorCorrection string `json:"error_correction" yaml:"error_correction" toml:"error_correction"`

	// Foreground is the QR code foreground color in hex (default: "#000000")
	Foreground string `json:"foreground" yaml:"foreground" toml:"foreground"`

	// Background is the QR code background color in hex (default: "#ffffff")
	Background string `json:"background" yaml:"background" toml:"background"`
}

// NewQRCodeConfig creates a new QRCodeConfig with default values.
func NewQRCodeConfig() QRCodeConfig {
	return QRCodeConfig{
		Enabled:         true,
		Format:          "svg",
		Size:            200,
		OutputDir:       "qrcodes",
		ErrorCorrection: "M",
		Foreground:      "#000000",
		Background:      "#ffffff",
	}
}

// YouTubeConfig configures the youtube plugin.
type YouTubeConfig struct {
	// Enabled controls whether YouTube URL conversion is active (default: true)
	Enabled bool `json:"enabled" yaml:"enabled" toml:"enabled"`

	// PrivacyEnhanced uses youtube-nocookie.com for enhanced privacy (default: true)
	PrivacyEnhanced bool `json:"privacy_enhanced" yaml:"privacy_enhanced" toml:"privacy_enhanced"`

	// ContainerClass is the CSS class for the embed container (default: "youtube-embed")
	ContainerClass string `json:"container_class" yaml:"container_class" toml:"container_class"`

	// LazyLoad enables lazy loading of iframe (default: true)
	LazyLoad bool `json:"lazy_load" yaml:"lazy_load" toml:"lazy_load"`
}

// NewYouTubeConfig creates a new YouTubeConfig with default values.
func NewYouTubeConfig() YouTubeConfig {
	return YouTubeConfig{
		Enabled:         true,
		PrivacyEnhanced: true,
		ContainerClass:  "youtube-embed",
		LazyLoad:        true,
	}
}

// CSSPurgeConfig configures the css_purge plugin for removing unused CSS rules.
type CSSPurgeConfig struct {
	// Enabled controls whether CSS purging is active (default: false)
	Enabled bool `json:"enabled" yaml:"enabled" toml:"enabled"`

	// Verbose enables detailed logging of purge operations (default: false)
	Verbose bool `json:"verbose" yaml:"verbose" toml:"verbose"`

	// Preserve is a list of glob patterns for CSS selectors to always keep.
	// These patterns match against class names and IDs.
	// Example: ["js-*", "htmx-*", "active", "hidden"]
	Preserve []string `json:"preserve" yaml:"preserve" toml:"preserve"`

	// PreserveAttributes is a list of attribute names to always preserve.
	// Selectors with these attributes (e.g., [data-theme], [data-palette]) will not be purged.
	// This is essential for runtime theming where attributes are set via JavaScript.
	// Example: ["data-theme", "data-palette"]
	PreserveAttributes []string `json:"preserve_attributes" yaml:"preserve_attributes" toml:"preserve_attributes"`

	// SkipFiles is a list of CSS file patterns to skip during purging.
	// Useful for third-party CSS that should not be modified.
	// Example: ["vendor/*", "normalize.css"]
	SkipFiles []string `json:"skip_files" yaml:"skip_files" toml:"skip_files"`

	// WarningThreshold is the minimum percentage of CSS removed before showing a warning.
	// If more than this percentage is removed, a warning is shown (might indicate overly aggressive purging).
	// Set to 0 to disable. Default: 0 (disabled)
	WarningThreshold int `json:"warning_threshold" yaml:"warning_threshold" toml:"warning_threshold"`
}

// NewCSSPurgeConfig creates a new CSSPurgeConfig with default values.
func NewCSSPurgeConfig() CSSPurgeConfig {
	return CSSPurgeConfig{
		Enabled:            false, // Disabled by default - opt-in feature
		Verbose:            false,
		Preserve:           []string{}, // Uses csspurge.DefaultPreservePatterns() when empty
		PreserveAttributes: []string{}, // Uses csspurge.DefaultPreserveAttributes() when empty
		SkipFiles:          []string{},
		WarningThreshold:   0,
	}
}

// TailwindConfig configures Tailwind CSS automation and inclusion.
type TailwindConfig struct {
	// Include controls how Tailwind is injected: "css", "js", or false.
	// Empty means use the default (css).
	Include *string `json:"include,omitempty" yaml:"include,omitempty" toml:"include,omitempty"`

	// Preflight controls whether Tailwind's base reset styles are enabled.
	// Defaults to false for markata-go's built-in theme so utility classes work
	// without overriding existing theme typography and spacing.
	Preflight *bool `json:"preflight,omitempty" yaml:"preflight,omitempty" toml:"preflight,omitempty"`

	// Input is the Tailwind input CSS file path. Relative paths are resolved
	// against assets_dir unless they already start with assets_dir.
	Input string `json:"input,omitempty" yaml:"input,omitempty" toml:"input,omitempty"`

	// Output is the Tailwind output CSS file path. Relative paths are resolved
	// against assets_dir unless they already start with assets_dir.
	Output string `json:"output,omitempty" yaml:"output,omitempty" toml:"output,omitempty"`

	// ConfigFile is the Tailwind config file path (optional).
	ConfigFile string `json:"config_file,omitempty" yaml:"config_file,omitempty" toml:"config_file,omitempty"`

	// Build controls whether Tailwind should be executed automatically (default: true).
	Build *bool `json:"build,omitempty" yaml:"build,omitempty" toml:"build,omitempty"`

	// Minify enables Tailwind's built-in minification (default: true).
	Minify *bool `json:"minify,omitempty" yaml:"minify,omitempty" toml:"minify,omitempty"`

	// AutoInstall enables automatic Tailwind CLI installation (default: true).
	AutoInstall *bool `json:"auto_install,omitempty" yaml:"auto_install,omitempty" toml:"auto_install,omitempty"`

	// Version is the Tailwind CLI version to install (default: "latest").
	Version string `json:"version,omitempty" yaml:"version,omitempty" toml:"version,omitempty"`

	// CacheDir is the directory for caching Tailwind binaries (default: XDG cache).
	CacheDir string `json:"cache_dir,omitempty" yaml:"cache_dir,omitempty" toml:"cache_dir,omitempty"`

	// Binary is the explicit path to the Tailwind CLI binary (optional).
	Binary string `json:"binary,omitempty" yaml:"binary,omitempty" toml:"binary,omitempty"`

	// ExtraArgs are additional arguments passed to the Tailwind CLI.
	ExtraArgs []string `json:"extra_args,omitempty" yaml:"extra_args,omitempty" toml:"extra_args,omitempty"`

	// Verbose enables verbose Tailwind CLI output (default: false).
	Verbose *bool `json:"verbose,omitempty" yaml:"verbose,omitempty" toml:"verbose,omitempty"`
}

// NewTailwindConfig creates a new TailwindConfig with default values.
func NewTailwindConfig() TailwindConfig {
	build := true
	autoInstall := true
	minify := true
	include := "css"
	preflight := false
	return TailwindConfig{
		Include:     &include,
		Preflight:   &preflight,
		Input:       "tailwind.css",
		Output:      "markata-tailwind.css",
		ConfigFile:  "",
		Build:       &build,
		Minify:      &minify,
		AutoInstall: &autoInstall,
		Version:     "v3.4.19",
		CacheDir:    "",
		Binary:      "",
		ExtraArgs:   []string{},
	}
}

// IsBuildEnabled returns whether Tailwind auto-build is enabled.
// Defaults to true if not explicitly set.
func (t *TailwindConfig) IsBuildEnabled() bool {
	if t.Build == nil {
		return true
	}
	return *t.Build
}

// IsMinifyEnabled returns whether Tailwind minification is enabled.
// Defaults to true if not explicitly set.
func (t *TailwindConfig) IsMinifyEnabled() bool {
	if t.Minify == nil {
		return true
	}
	return *t.Minify
}

// IsAutoInstallEnabled returns whether Tailwind auto-install is enabled.
// Defaults to true if not explicitly set.
func (t *TailwindConfig) IsAutoInstallEnabled() bool {
	if t.AutoInstall == nil {
		return true
	}
	return *t.AutoInstall
}

// IsVerbose returns whether Tailwind CLI should be verbose.
// Defaults to false if not explicitly set.
func (t *TailwindConfig) IsVerbose() bool {
	if t.Verbose == nil {
		return false
	}
	return *t.Verbose
}

// IncludeMode resolves the Tailwind include mode with defaults.
func (t *TailwindConfig) IncludeMode() string {
	if t.Include == nil {
		return "css"
	}
	return *t.Include
}

// IsPreflightEnabled returns whether Tailwind preflight reset styles are enabled.
// Defaults to false if not explicitly set.
func (t *TailwindConfig) IsPreflightEnabled() bool {
	if t.Preflight == nil {
		return false
	}
	return *t.Preflight
}

// ImageZoomConfig configures the image_zoom plugin for lightbox functionality.
type ImageZoomConfig struct {
	// Enabled controls whether image zoom is active (default: false)
	Enabled bool `json:"enabled" yaml:"enabled" toml:"enabled"`

	// Library is the lightbox library to use (default: "glightbox")
	// Supported: "glightbox"
	Library string `json:"library" yaml:"library" toml:"library"`

	// Selector is the CSS selector for images to make zoomable (default: ".glightbox")
	Selector string `json:"selector" yaml:"selector" toml:"selector"`

	// CDN uses CDN for library files instead of local (default: true)
	CDN bool `json:"cdn" yaml:"cdn" toml:"cdn"`

	// AutoAllImages enables zoom on all images without explicit marking (default: false)
	AutoAllImages bool `json:"auto_all_images" yaml:"auto_all_images" toml:"auto_all_images"`

	// OpenEffect is the effect when opening the lightbox (default: "zoom")
	// Options: "zoom", "fade", "none"
	OpenEffect string `json:"open_effect" yaml:"open_effect" toml:"open_effect"`

	// CloseEffect is the effect when closing the lightbox (default: "zoom")
	// Options: "zoom", "fade", "none"
	CloseEffect string `json:"close_effect" yaml:"close_effect" toml:"close_effect"`

	// SlideEffect is the effect when sliding between images (default: "slide")
	// Options: "slide", "fade", "zoom", "none"
	SlideEffect string `json:"slide_effect" yaml:"slide_effect" toml:"slide_effect"`

	// TouchNavigation enables touch/swipe navigation (default: true)
	TouchNavigation bool `json:"touch_navigation" yaml:"touch_navigation" toml:"touch_navigation"`

	// Loop enables looping through images in a gallery (default: false)
	Loop bool `json:"loop" yaml:"loop" toml:"loop"`

	// Draggable enables dragging images to navigate (default: true)
	Draggable bool `json:"draggable" yaml:"draggable" toml:"draggable"`
}

// NewImageZoomConfig creates a new ImageZoomConfig with default values.
func NewImageZoomConfig() ImageZoomConfig {
	return ImageZoomConfig{
		Enabled:         false, // Disabled by default
		Library:         "glightbox",
		Selector:        ".post-content .glightbox",
		CDN:             true,
		AutoAllImages:   false,
		OpenEffect:      "zoom",
		CloseEffect:     "zoom",
		SlideEffect:     "slide",
		TouchNavigation: true,
		Loop:            false,
		Draggable:       true,
	}
}

// CriticalCSSConfig configures the critical CSS extraction and inlining plugin.
// Critical CSS optimization inlines above-the-fold styles and async loads the rest,
// improving First Contentful Paint (FCP) by 200-800ms.
type CriticalCSSConfig struct {
	// Enabled controls whether critical CSS optimization is active (default: false)
	Enabled *bool `json:"enabled,omitempty" yaml:"enabled,omitempty" toml:"enabled,omitempty"`

	// ViewportWidth is reserved for future viewport-based critical CSS detection.
	// NOTE: Currently NOT implemented. The critical CSS extractor uses a selector-based
	// approach (see ExtraSelectors/ExcludeSelectors) rather than viewport simulation.
	// True viewport-based detection would require a headless browser.
	// This field is retained for configuration compatibility and future enhancement.
	// Default: 1300
	ViewportWidth int `json:"viewport_width,omitempty" yaml:"viewport_width,omitempty" toml:"viewport_width,omitempty"`

	// ViewportHeight is reserved for future viewport-based critical CSS detection.
	// NOTE: Currently NOT implemented. The critical CSS extractor uses a selector-based
	// approach (see ExtraSelectors/ExcludeSelectors) rather than viewport simulation.
	// True viewport-based detection would require a headless browser.
	// This field is retained for configuration compatibility and future enhancement.
	// Default: 900
	ViewportHeight int `json:"viewport_height,omitempty" yaml:"viewport_height,omitempty" toml:"viewport_height,omitempty"`

	// Minify controls whether to minify the critical CSS output (default: true)
	Minify *bool `json:"minify,omitempty" yaml:"minify,omitempty" toml:"minify,omitempty"`

	// PreloadNonCritical uses link rel="preload" for non-critical CSS (default: true)
	// This async loads the full stylesheet without blocking render
	PreloadNonCritical *bool `json:"preload_non_critical,omitempty" yaml:"preload_non_critical,omitempty" toml:"preload_non_critical,omitempty"`

	// ExtraSelectors is a list of additional CSS selectors to always include as critical
	// Useful for JavaScript-injected content that may appear above the fold
	ExtraSelectors []string `json:"extra_selectors,omitempty" yaml:"extra_selectors,omitempty" toml:"extra_selectors,omitempty"`

	// ExcludeSelectors is a list of CSS selectors to always exclude from critical CSS
	// Useful for content that should never be inlined (e.g., large animations)
	ExcludeSelectors []string `json:"exclude_selectors,omitempty" yaml:"exclude_selectors,omitempty" toml:"exclude_selectors,omitempty"`

	// InlineThreshold is the maximum size (in bytes) for the critical CSS before giving up inlining (default: 50000)
	// If critical CSS exceeds this threshold, the optimization is skipped for that page
	InlineThreshold int `json:"inline_threshold,omitempty" yaml:"inline_threshold,omitempty" toml:"inline_threshold,omitempty"`
}

// NewCriticalCSSConfig creates a new CriticalCSSConfig with default values.
// Note: ViewportWidth and ViewportHeight are set for future compatibility but
// are not currently used by the selector-based critical CSS extractor.
func NewCriticalCSSConfig() CriticalCSSConfig {
	enabled := false
	minify := true
	preloadNonCritical := true
	return CriticalCSSConfig{
		Enabled:            &enabled,
		ViewportWidth:      1300,
		ViewportHeight:     900,
		Minify:             &minify,
		PreloadNonCritical: &preloadNonCritical,
		ExtraSelectors:     []string{},
		ExcludeSelectors:   []string{},
		InlineThreshold:    50000,
	}
}

// IsEnabled returns whether critical CSS optimization is enabled.
// Defaults to false if not explicitly set.
func (c *CriticalCSSConfig) IsEnabled() bool {
	if c.Enabled == nil {
		return false
	}
	return *c.Enabled
}

// IsMinify returns whether critical CSS minification is enabled.
// Defaults to true if not explicitly set.
func (c *CriticalCSSConfig) IsMinify() bool {
	if c.Minify == nil {
		return true
	}
	return *c.Minify
}

// IsPreloadNonCritical returns whether non-critical CSS should be preloaded.
// Defaults to true if not explicitly set.
func (c *CriticalCSSConfig) IsPreloadNonCritical() bool {
	if c.PreloadNonCritical == nil {
		return true
	}
	return *c.PreloadNonCritical
}

// EmbedsConfig configures the embeds plugin for embedding internal and external content.
type EmbedsConfig struct {
	// Enabled controls whether embed processing is active (default: true)
	Enabled bool `json:"enabled" yaml:"enabled" toml:"enabled"`

	// InternalCardClass is the CSS class for internal embed cards (default: "embed-card")
	InternalCardClass string `json:"internal_card_class" yaml:"internal_card_class" toml:"internal_card_class"`

	// ExternalCardClass is the CSS class for external embed cards (default: "embed-card embed-card-external")
	ExternalCardClass string `json:"external_card_class" yaml:"external_card_class" toml:"external_card_class"`

	// FetchExternal enables fetching OG metadata for external embeds (default: true)
	FetchExternal bool `json:"fetch_external" yaml:"fetch_external" toml:"fetch_external"`

	// OEmbedEnabled enables oEmbed resolution for external embeds (default: true)
	OEmbedEnabled bool `json:"oembed_enabled" yaml:"oembed_enabled" toml:"oembed_enabled"`

	// ResolutionStrategy controls how external embeds are resolved (default: "oembed_first")
	ResolutionStrategy string `json:"resolution_strategy" yaml:"resolution_strategy" toml:"resolution_strategy"`

	// OEmbedProviders configures provider-specific overrides (default: empty)
	OEmbedProviders map[string]OEmbedProviderConfig `json:"providers" yaml:"providers" toml:"providers"`

	// OEmbedAutoDiscover enables auto-discovery via HTML link tags (default: false)
	OEmbedAutoDiscover bool `json:"oembed_auto_discover" yaml:"oembed_auto_discover" toml:"oembed_auto_discover"`

	// DefaultEmbedMode is the default mode for embeds when not specified (default: "card")
	// Options: rich, card, performance, hover
	DefaultEmbedMode string `json:"default_mode" yaml:"default_mode" toml:"default_mode"`

	// OEmbedProvidersURL is the URL to fetch oEmbed providers from (default: "https://oembed.com/providers.json")
	// Set to empty string to disable fetching from providers.json
	OEmbedProvidersURL string `json:"oembed_providers_url" yaml:"oembed_providers_url" toml:"oembed_providers_url"`

	// CacheDir is the directory for caching external embed metadata (default: ".cache/embeds")
	CacheDir string `json:"cache_dir" yaml:"cache_dir" toml:"cache_dir"`

	// Timeout is the HTTP request timeout in seconds for external fetches (default: 10)
	Timeout int `json:"timeout" yaml:"timeout" toml:"timeout"`

	// CacheTTL is the cache time-to-live in seconds for external metadata (default: 604800)
	CacheTTL int `json:"cache_ttl" yaml:"cache_ttl" toml:"cache_ttl"`

	// FallbackTitle is used when OG title cannot be fetched (default: "External Link")
	FallbackTitle string `json:"fallback_title" yaml:"fallback_title" toml:"fallback_title"`

	// ShowImage controls whether to display OG images in external embeds (default: true)
	ShowImage bool `json:"show_image" yaml:"show_image" toml:"show_image"`

	// AttachmentsPrefix is the URL prefix for attachment embeds (default: "/static/")
	// Used for Obsidian-style ![[image.jpg]] syntax
	AttachmentsPrefix string `json:"attachments_prefix" yaml:"attachments_prefix" toml:"attachments_prefix"`
}

// OEmbedProviderConfig configures a single oEmbed provider.
type OEmbedProviderConfig struct {
	// Enabled controls whether the provider is used (default: true)
	Enabled bool `json:"enabled" yaml:"enabled" toml:"enabled"`

	// Mode overrides the default embed mode for this provider.
	// Options: rich, card, performance, hover, image_only
	// When not set, the default is determined by the provider's oEmbed type:
	//   - photo: image_only
	//   - video: rich
	//   - rich: rich
	//   - link: card
	Mode string `json:"mode" yaml:"mode" toml:"mode"`
}

// NewEmbedsConfig creates a new EmbedsConfig with default values.
func NewEmbedsConfig() EmbedsConfig {
	return EmbedsConfig{
		Enabled:            true,
		InternalCardClass:  "embed-card",
		ExternalCardClass:  "embed-card embed-card-external",
		FetchExternal:      true,
		OEmbedEnabled:      true,
		ResolutionStrategy: "oembed_first",
		OEmbedProviders: map[string]OEmbedProviderConfig{
			"giphy":   {Enabled: true, Mode: "image_only"},
			"tiktok":  {Enabled: true, Mode: "card"},
			"youtube": {Enabled: true, Mode: "rich"},
		},
		OEmbedAutoDiscover: false,
		DefaultEmbedMode:   "card",
		OEmbedProvidersURL: "https://oembed.com/providers.json",
		CacheDir:           ".cache/embeds",
		Timeout:            10,
		CacheTTL:           7 * 24 * 60 * 60,
		FallbackTitle:      "External Link",
		ShowImage:          true,
		AttachmentsPrefix:  "/static/",
	}
}

// EncryptionConfig configures content encryption for private posts.
// When enabled and a post has private: true, the post content will be encrypted
// using AES-256-GCM and require client-side decryption.
//
// Encryption is enabled by default with default_key="default". To use it,
// set MARKATA_GO_ENCRYPTION_KEY_DEFAULT in your environment or .env file.
type EncryptionConfig struct {
	// Enabled controls whether encryption processing is active (default: true)
	Enabled bool `json:"enabled" yaml:"enabled" toml:"enabled"`

	// DefaultKey is the default encryption key name to use when a post doesn't specify one.
	// Maps to environment variable MARKATA_GO_ENCRYPTION_KEY_{DefaultKey}
	// Example: default_key: "blog" reads from MARKATA_GO_ENCRYPTION_KEY_BLOG
	// Default: "default" (reads from MARKATA_GO_ENCRYPTION_KEY_DEFAULT)
	DefaultKey string `json:"default_key,omitempty" yaml:"default_key,omitempty" toml:"default_key,omitempty"`

	// DecryptionHint is a hint shown to users about how to get the decryption password
	// Example: "Contact me on Twitter @user to get access"
	DecryptionHint string `json:"decryption_hint,omitempty" yaml:"decryption_hint,omitempty" toml:"decryption_hint,omitempty"`

	// PrivateTags maps tag names to encryption key names.
	// Any post with a matching tag is automatically treated as private and encrypted
	// with the specified key. Frontmatter secret_key overrides the tag-level key.
	// Example: {"diary": "personal", "draft-ideas": "default"}
	// This means any post tagged "diary" is encrypted with MARKATA_GO_ENCRYPTION_KEY_PERSONAL
	PrivateTags map[string]string `json:"private_tags,omitempty" yaml:"private_tags,omitempty" toml:"private_tags,omitempty"`

	// EnforceStrength requires all encryption keys to satisfy the configured strength policy.
	EnforceStrength bool `json:"enforce_strength" yaml:"enforce_strength" toml:"enforce_strength"`

	// MinEstimatedCrackTime is the minimum allowed crack-time estimate per key (supports y/d/h/m/s).
	MinEstimatedCrackTime string `json:"min_estimated_crack_time,omitempty" yaml:"min_estimated_crack_time,omitempty" toml:"min_estimated_crack_time,omitempty"`

	// MinPasswordLength is the minimum number of characters for every encryption password.
	MinPasswordLength int `json:"min_password_length,omitempty" yaml:"min_password_length,omitempty" toml:"min_password_length,omitempty"`
}

// NewEncryptionConfig creates a new EncryptionConfig with default values.
// Encryption is enabled by default with default_key="default" so that
// any post with private: true is automatically encrypted. Users only need
// to set MARKATA_GO_ENCRYPTION_KEY_DEFAULT in their environment or .env file.
func NewEncryptionConfig() EncryptionConfig {
	return EncryptionConfig{
		Enabled:               true,
		DefaultKey:            "default",
		EnforceStrength:       true,
		MinEstimatedCrackTime: "10y",
		MinPasswordLength:     14,
	}
}

// ShortcutsConfig configures user-defined keyboard shortcuts.
// Shortcuts are organized by group (e.g., "navigation") and map key sequences to URLs.
type ShortcutsConfig struct {
	// Navigation contains shortcuts for navigating to specific pages.
	// Keys are key sequences (e.g., "g t"), values are destination URLs.
	// Example: {"g t": "/tags/", "g a": "/about/"}
	Navigation map[string]string `json:"navigation,omitempty" yaml:"navigation,omitempty" toml:"navigation,omitempty"`
}

// NewShortcutsConfig creates a new ShortcutsConfig with default values.
func NewShortcutsConfig() ShortcutsConfig {
	return ShortcutsConfig{
		Navigation: make(map[string]string),
	}
}

// ViewTransitionsConfig configures client-side page transition behavior.
type ViewTransitionsConfig struct {
	// Enabled controls whether page transitions are enabled.
	Enabled *bool `json:"enabled,omitempty" yaml:"enabled,omitempty" toml:"enabled,omitempty"`

	// Native enables browser-native cross-document transitions instead of SPA interception.
	// This is experimental and opt-in; unsupported browsers use normal links.
	Native *bool `json:"native,omitempty" yaml:"native,omitempty" toml:"native,omitempty"`

	// Debug enables browser console logging and timing output.
	Debug *bool `json:"debug,omitempty" yaml:"debug,omitempty" toml:"debug,omitempty"`

	// Duration is the reference transition duration in milliseconds.
	Duration int `json:"duration,omitempty" yaml:"duration,omitempty" toml:"duration,omitempty"`

	// UpdateMeta controls whether head meta tags are synchronized on navigation.
	UpdateMeta *bool `json:"update_meta,omitempty" yaml:"update_meta,omitempty" toml:"update_meta,omitempty"`

	// ScrollToTop controls whether navigation scrolls to the top unless a hash is present.
	ScrollToTop *bool `json:"scroll_to_top,omitempty" yaml:"scroll_to_top,omitempty" toml:"scroll_to_top,omitempty"`

	// SkipClasses are additional CSS classes that should bypass transitions.
	SkipClasses []string `json:"skip_classes,omitempty" yaml:"skip_classes,omitempty" toml:"skip_classes,omitempty"`

	// SkipSelectors are additional selectors that should bypass transitions.
	SkipSelectors []string `json:"skip_selectors,omitempty" yaml:"skip_selectors,omitempty" toml:"skip_selectors,omitempty"`
}

// NewViewTransitionsConfig creates a new ViewTransitionsConfig with default values.
func NewViewTransitionsConfig() ViewTransitionsConfig {
	enabled := true
	debug := false
	updateMeta := true
	scrollToTop := true

	return ViewTransitionsConfig{
		Enabled:       &enabled,
		Debug:         &debug,
		Duration:      300,
		UpdateMeta:    &updateMeta,
		ScrollToTop:   &scrollToTop,
		SkipClasses:   []string{},
		SkipSelectors: []string{},
	}
}

// HasCustomShortcuts returns true if any custom shortcuts are defined.
func (s *ShortcutsConfig) HasCustomShortcuts() bool {
	return len(s.Navigation) > 0
}

// ContentTemplateConfig defines a single content template.
type ContentTemplateConfig struct {
	// Name is the template identifier (e.g., "post", "page", "docs")
	Name string `json:"name" yaml:"name" toml:"name"`

	// Directory is the output directory for this content type
	Directory string `json:"directory" yaml:"directory" toml:"directory"`

	// Frontmatter contains default frontmatter fields for this template
	Frontmatter map[string]interface{} `json:"frontmatter,omitempty" yaml:"frontmatter,omitempty" toml:"frontmatter,omitempty"`

	// Body is the default body content (markdown) for this template
	Body string `json:"body,omitempty" yaml:"body,omitempty" toml:"body,omitempty"`
}

// ContentTemplatesConfig configures the content template system for the new command.
type ContentTemplatesConfig struct {
	// Directory is where user-defined templates are stored (default: "content-templates")
	Directory string `json:"directory" yaml:"directory" toml:"directory"`

	// Placement maps template names to output directories
	Placement map[string]string `json:"placement" yaml:"placement" toml:"placement"`

	// Templates is a list of custom template configurations
	Templates []ContentTemplateConfig `json:"templates,omitempty" yaml:"templates,omitempty" toml:"templates,omitempty"`
}

// NewContentTemplatesConfig creates a new ContentTemplatesConfig with default values.
func NewContentTemplatesConfig() ContentTemplatesConfig {
	return ContentTemplatesConfig{
		Directory: "content-templates",
		Placement: map[string]string{
			"post": "posts",
			"page": "pages",
			"docs": "docs",
		},
		Templates: []ContentTemplateConfig{},
	}
}

// GetPlacement returns the output directory for a template name.
// Returns the template name itself if no explicit placement is configured.
func (c *ContentTemplatesConfig) GetPlacement(templateName string) string {
	if dir, ok := c.Placement[templateName]; ok {
		return dir
	}
	return templateName
}

// ErrorPagesConfig configures custom error pages (404, etc.).
type ErrorPagesConfig struct {
	// Enable404 enables the built-in 404 page (default: true)
	Enable404 *bool `json:"enable_404,omitempty" yaml:"enable_404,omitempty" toml:"enable_404,omitempty"`

	// Custom404Template is the path to a custom 404 template (default: "404.html")
	Custom404Template string `json:"custom_404_template,omitempty" yaml:"custom_404_template,omitempty" toml:"custom_404_template,omitempty"`

	// MaxSuggestions is the maximum number of similar posts to suggest (default: 5)
	MaxSuggestions int `json:"max_suggestions,omitempty" yaml:"max_suggestions,omitempty" toml:"max_suggestions,omitempty"`
}

// NewErrorPagesConfig creates a new ErrorPagesConfig with default values.
func NewErrorPagesConfig() ErrorPagesConfig {
	enabled := true
	return ErrorPagesConfig{
		Enable404:         &enabled,
		Custom404Template: "404.html",
		MaxSuggestions:    5,
	}
}

// Is404Enabled returns whether the 404 page is enabled.
// Defaults to true if not explicitly set.
func (e *ErrorPagesConfig) Is404Enabled() bool {
	if e.Enable404 == nil {
		return true
	}
	return *e.Enable404
}

// TagsConfig configures the tags listing page at /tags.
// The tags listing page shows all available tags with post counts and links to tag pages.
type TagsConfig struct {
	// Enabled controls whether the tags listing page is generated (default: true)
	Enabled *bool `json:"enabled,omitempty" yaml:"enabled,omitempty" toml:"enabled,omitempty"`

	// Blacklist is a list of tag names to completely exclude from the tags listing.
	// These tags will not appear on the /tags page and won't be visible publicly.
	// Example: ["draft", "wip", "internal"]
	Blacklist []string `json:"blacklist,omitempty" yaml:"blacklist,omitempty" toml:"blacklist,omitempty"`

	// Private is a list of tag names that exist but should be hidden from the listing.
	// Posts with these tags are still accessible via direct URL, but the tags
	// won't appear on the /tags overview page.
	// Example: ["personal", "unlisted"]
	Private []string `json:"private,omitempty" yaml:"private,omitempty" toml:"private,omitempty"`

	// Title is the title for the tags listing page (default: "Tags")
	Title string `json:"title,omitempty" yaml:"title,omitempty" toml:"title,omitempty"`

	// Description is the description for the tags listing page
	Description string `json:"description,omitempty" yaml:"description,omitempty" toml:"description,omitempty"`

	// Template is the template file to use (default: "tags.html")
	Template string `json:"template,omitempty" yaml:"template,omitempty" toml:"template,omitempty"`

	// SlugPrefix is the URL prefix for the tags listing (default: "tags")
	SlugPrefix string `json:"slug_prefix,omitempty" yaml:"slug_prefix,omitempty" toml:"slug_prefix,omitempty"`
}

// NewTagsConfig creates a new TagsConfig with default values.
func NewTagsConfig() TagsConfig {
	enabled := true
	return TagsConfig{
		Enabled:     &enabled,
		Blacklist:   []string{},
		Private:     []string{},
		Title:       "Tags",
		Description: "",
		Template:    "tags.html",
		SlugPrefix:  "tags",
	}
}

// IsEnabled returns whether the tags listing page is enabled.
// Defaults to true if not explicitly set.
func (t *TagsConfig) IsEnabled() bool {
	if t.Enabled == nil {
		return true
	}
	return *t.Enabled
}

// IsBlacklisted returns whether a tag is in the blacklist.
func (t *TagsConfig) IsBlacklisted(tag string) bool {
	for _, b := range t.Blacklist {
		if b == tag {
			return true
		}
	}
	return false
}

// IsPrivate returns whether a tag is in the private list.
func (t *TagsConfig) IsPrivate(tag string) bool {
	for _, p := range t.Private {
		if p == tag {
			return true
		}
	}
	return false
}

// FeedsPageConfig configures the feeds listing page at /feeds.
type FeedsPageConfig struct {
	// Enabled controls whether the feeds listing page is generated (default: true)
	Enabled *bool `json:"enabled,omitempty" yaml:"enabled,omitempty" toml:"enabled,omitempty"`

	// Title is the title for the feeds listing page (default: "Feeds")
	Title string `json:"title,omitempty" yaml:"title,omitempty" toml:"title,omitempty"`

	// Description is the description for the feeds listing page.
	Description string `json:"description,omitempty" yaml:"description,omitempty" toml:"description,omitempty"`

	// Template is the template file to use (default: "feeds.html")
	Template string `json:"template,omitempty" yaml:"template,omitempty" toml:"template,omitempty"`

	// SlugPrefix is the URL prefix for the feeds listing (default: "feeds")
	SlugPrefix string `json:"slug_prefix,omitempty" yaml:"slug_prefix,omitempty" toml:"slug_prefix,omitempty"`

	// Robots controls the robots meta tag for the feeds listing HTML pages.
	// Example: "noindex,follow"
	Robots string `json:"robots,omitempty" yaml:"robots,omitempty" toml:"robots,omitempty"`

	// ShowPrivateFeeds lists private-enabled feed slugs that should appear on /feeds/.
	// Public feeds are always listed; private-enabled feeds remain hidden unless their
	// slug appears here.
	ShowPrivateFeeds []string `json:"show_private_feeds,omitempty" yaml:"show_private_feeds,omitempty" toml:"show_private_feeds,omitempty"`
}

// NewFeedsPageConfig creates a new FeedsPageConfig with default values.
func NewFeedsPageConfig() FeedsPageConfig {
	enabled := true
	return FeedsPageConfig{
		Enabled:          &enabled,
		Title:            "Feeds",
		Description:      "Browse the public feeds available on this site.",
		Template:         "feeds.html",
		SlugPrefix:       "feeds",
		Robots:           "",
		ShowPrivateFeeds: nil,
	}
}

// IsEnabled returns whether the feeds listing page is enabled.
func (f *FeedsPageConfig) IsEnabled() bool {
	if f.Enabled == nil {
		return true
	}
	return *f.Enabled
}

// TagAggregatorConfig configures the tag aggregator plugin for normalizing and expanding tags.
type TagAggregatorConfig struct {
	// Enabled controls whether tag aggregation is active (default: true)
	Enabled *bool `json:"enabled,omitempty" yaml:"enabled,omitempty" toml:"enabled,omitempty"`

	// Synonyms maps canonical tags to their synonyms/variants.
	// When a post has a synonym tag, it's replaced with the canonical tag.
	// Example: {"kubernetes": ["k8s"], "javascript": ["js"]}
	Synonyms map[string][]string `json:"synonyms,omitempty" yaml:"synonyms,omitempty" toml:"synonyms,omitempty"`

	// Additional maps tags to additional tags that should be automatically added.
	// These are applied recursively to create tag hierarchies.
	// Example: {"pandas": ["data", "python"], "docker": ["containers"]}
	Additional map[string][]string `json:"additional,omitempty" yaml:"additional,omitempty" toml:"additional,omitempty"`

	// GenerateReport controls whether to generate a debug report page (default: false)
	GenerateReport bool `json:"generate_report,omitempty" yaml:"generate_report,omitempty" toml:"generate_report,omitempty"`
}

// NewTagAggregatorConfig creates a new TagAggregatorConfig with default values.
func NewTagAggregatorConfig() TagAggregatorConfig {
	enabled := true
	return TagAggregatorConfig{
		Enabled:        &enabled,
		Synonyms:       make(map[string][]string),
		Additional:     make(map[string][]string),
		GenerateReport: false,
	}
}

// IsEnabled returns whether tag aggregation is enabled.
// Defaults to true if not explicitly set.
func (t *TagAggregatorConfig) IsEnabled() bool {
	if t.Enabled == nil {
		return true
	}
	return *t.Enabled
}

// GardenConfig configures the garden view plugin for knowledge graph export and visualization.
// The garden provides a way to explore content by relationships (tags, links) rather than
// chronology. It exports a graph.json file and optionally renders an interactive garden page.
type GardenConfig struct {
	// Enabled controls whether the garden view is generated (default: true)
	Enabled *bool `json:"enabled,omitempty" yaml:"enabled,omitempty" toml:"enabled,omitempty"`

	// Path is the output path prefix for garden files (default: "garden")
	Path string `json:"path,omitempty" yaml:"path,omitempty" toml:"path,omitempty"`

	// ExportJSON enables emitting graph.json with the relationship graph (default: true)
	ExportJSON *bool `json:"export_json,omitempty" yaml:"export_json,omitempty" toml:"export_json,omitempty"`

	// RenderPage enables generating garden/index.html (default: true)
	RenderPage *bool `json:"render_page,omitempty" yaml:"render_page,omitempty" toml:"render_page,omitempty"`

	// IncludeTags includes tag nodes in the graph (default: true)
	IncludeTags *bool `json:"include_tags,omitempty" yaml:"include_tags,omitempty" toml:"include_tags,omitempty"`

	// IncludePosts includes post nodes in the graph (default: true)
	IncludePosts *bool `json:"include_posts,omitempty" yaml:"include_posts,omitempty" toml:"include_posts,omitempty"`

	// MaxNodes is the maximum number of nodes in the graph (default: 2000)
	// When exceeded, least-used tags are removed first
	MaxNodes int `json:"max_nodes,omitempty" yaml:"max_nodes,omitempty" toml:"max_nodes,omitempty"`

	// ExcludeTags is a list of tag names to exclude from the graph
	ExcludeTags []string `json:"exclude_tags,omitempty" yaml:"exclude_tags,omitempty" toml:"exclude_tags,omitempty"`

	// Template is the template file to use for the garden page (default: "garden.html")
	Template string `json:"template,omitempty" yaml:"template,omitempty" toml:"template,omitempty"`

	// Title is the title for the garden page (default: "Garden")
	Title string `json:"title,omitempty" yaml:"title,omitempty" toml:"title,omitempty"`

	// Description is the description for the garden page
	Description string `json:"description,omitempty" yaml:"description,omitempty" toml:"description,omitempty"`
}

// ImagesConfig configures the generated image inventory and authoring page.
// The image library is enabled by default and writes to /images/.
type ImagesConfig struct {
	// Enabled controls whether the image index and authoring page are generated.
	Enabled *bool `json:"enabled,omitempty" yaml:"enabled,omitempty" toml:"enabled,omitempty"`

	// Path is the output path prefix, relative to output_dir.
	Path string `json:"path,omitempty" yaml:"path,omitempty" toml:"path,omitempty"`

	// Template is the template used for the authoring page.
	Template string `json:"template,omitempty" yaml:"template,omitempty" toml:"template,omitempty"`

	// ExportJSON controls whether path/index.json and output_dir/images.json are written.
	ExportJSON *bool `json:"export_json,omitempty" yaml:"export_json,omitempty" toml:"export_json,omitempty"`

	// IncludeUnreferenced includes image and video files found in the static asset
	// directory even when no public post references them. It is opt-in because
	// enabling it makes every supported asset centrally enumerable. The default,
	// including a nil configuration value, is false.
	IncludeUnreferenced *bool `json:"include_unreferenced,omitempty" yaml:"include_unreferenced,omitempty" toml:"include_unreferenced,omitempty"`

	// PageSize opts into server-rendered pagination of the image library.
	// Zero retains the historical complete-page behavior.
	PageSize int `json:"page_size,omitempty" yaml:"page_size,omitempty" toml:"page_size,omitempty"`
}

// NewImagesConfig returns the default image-library configuration.
func NewImagesConfig() ImagesConfig {
	enabled := true
	exportJSON := true
	includeUnreferenced := false
	return ImagesConfig{
		Enabled:             &enabled,
		Path:                "images",
		Template:            "images.html",
		ExportJSON:          &exportJSON,
		IncludeUnreferenced: &includeUnreferenced,
	}
}

// IsEnabled reports whether the image library is enabled.
func (c *ImagesConfig) IsEnabled() bool {
	return c == nil || c.Enabled == nil || *c.Enabled
}

// ShouldExportJSON reports whether the JSON artifact should be written.
func (c *ImagesConfig) ShouldExportJSON() bool {
	return c == nil || c.ExportJSON == nil || *c.ExportJSON
}

// ShouldIncludeUnreferenced reports whether unreferenced static image files are
// included in the inventory. A nil receiver or field uses the safe default of
// false.
func (c *ImagesConfig) ShouldIncludeUnreferenced() bool {
	if c == nil || c.IncludeUnreferenced == nil {
		return false
	}
	return *c.IncludeUnreferenced
}

// NewGardenConfig creates a new GardenConfig with default values.
func NewGardenConfig() GardenConfig {
	enabled := true
	exportJSON := true
	renderPage := true
	includeTags := true
	includePosts := true
	return GardenConfig{
		Enabled:      &enabled,
		Path:         "garden",
		ExportJSON:   &exportJSON,
		RenderPage:   &renderPage,
		IncludeTags:  &includeTags,
		IncludePosts: &includePosts,
		MaxNodes:     2000,
		ExcludeTags:  []string{},
		Template:     "garden.html",
		Title:        "Garden",
		Description:  "",
	}
}

// IsEnabled returns whether the garden view is enabled.
// Defaults to true if not explicitly set.
func (g *GardenConfig) IsEnabled() bool {
	if g.Enabled == nil {
		return true
	}
	return *g.Enabled
}

// IsExportJSON returns whether graph.json export is enabled.
// Defaults to true if not explicitly set.
func (g *GardenConfig) IsExportJSON() bool {
	if g.ExportJSON == nil {
		return true
	}
	return *g.ExportJSON
}

// IsRenderPage returns whether the garden HTML page should be rendered.
// Defaults to true if not explicitly set.
func (g *GardenConfig) IsRenderPage() bool {
	if g.RenderPage == nil {
		return true
	}
	return *g.RenderPage
}

// IsIncludeTags returns whether tag nodes should be included in the graph.
// Defaults to true if not explicitly set.
func (g *GardenConfig) IsIncludeTags() bool {
	if g.IncludeTags == nil {
		return true
	}
	return *g.IncludeTags
}

// IsIncludePosts returns whether post nodes should be included in the graph.
// Defaults to true if not explicitly set.
func (g *GardenConfig) IsIncludePosts() bool {
	if g.IncludePosts == nil {
		return true
	}
	return *g.IncludePosts
}

// IsTagExcluded returns whether a tag is in the exclude list.
func (g *GardenConfig) IsTagExcluded(tag string) bool {
	for _, t := range g.ExcludeTags {
		if t == tag {
			return true
		}
	}
	return false
}

// GetPath returns the output path, with default if not set.
func (g *GardenConfig) GetPath() string {
	if g.Path == "" {
		return "garden"
	}
	return g.Path
}

// GetTemplate returns the template name, with default if not set.
func (g *GardenConfig) GetTemplate() string {
	if g.Template == "" {
		return "garden.html"
	}
	return g.Template
}

// GetMaxNodes returns the max nodes, with default if not set.
func (g *GardenConfig) GetMaxNodes() int {
	if g.MaxNodes <= 0 {
		return 2000
	}
	return g.MaxNodes
}

// CSSBundleConfig configures css_bundle plugin for combining CSS files.
type CSSBundleConfig struct {
	// Enabled controls whether CSS bundling is active (default: false)
	Enabled bool `json:"enabled" yaml:"enabled" toml:"enabled"`

	// Bundles is list of bundle configurations
	Bundles []BundleConfig `json:"bundles" yaml:"bundles" toml:"bundles"`

	// Pages lists output-relative HTML files whose head styles are bundled after cleanup.
	Pages []string `json:"pages" yaml:"pages" toml:"pages"`

	// Exclude is a list of CSS file patterns to exclude from bundling
	Exclude []string `json:"exclude" yaml:"exclude" toml:"exclude"`

	// Minify controls whether bundled CSS is minified (default: false)
	// Note: minification is not yet implemented
	Minify bool `json:"minify" yaml:"minify" toml:"minify"`

	// AddSourceComments adds comments indicating source files in bundles (default: true)
	AddSourceComments *bool `json:"add_source_comments,omitempty" yaml:"add_source_comments,omitempty" toml:"add_source_comments,omitempty"`
}

// BundleConfig defines a single CSS bundle.
type BundleConfig struct {
	// Name is the bundle identifier (e.g., "main", "critical")
	Name string `json:"name" yaml:"name" toml:"name"`

	// Sources is a list of CSS file paths or glob patterns to include
	// Files are concatenated in the order specified
	Sources []string `json:"sources" yaml:"sources" toml:"sources"`

	// Output is the output file path relative to output_dir (e.g., "css/bundle.css")
	Output string `json:"output" yaml:"output" toml:"output"`
}

// NewCSSBundleConfig creates a new CSSBundleConfig with default values.
func NewCSSBundleConfig() CSSBundleConfig {
	addSourceComments := true
	return CSSBundleConfig{
		Enabled:           false,
		Bundles:           []BundleConfig{},
		Exclude:           []string{},
		Minify:            false,
		AddSourceComments: &addSourceComments,
	}
}

// IsAddSourceComments returns whether source comments should be added to bundles.
// Defaults to true if not explicitly set.
func (c *CSSBundleConfig) IsAddSourceComments() bool {
	if c.AddSourceComments == nil {
		return true
	}
	return *c.AddSourceComments
}

// CSSMinifyConfig configures the css_minify plugin for minifying CSS files.
// CSS minification reduces file sizes by 15-30%, improving page load performance
// and Lighthouse scores.
type CSSMinifyConfig struct {
	// Enabled controls whether CSS minification is active (default: true)
	Enabled bool `json:"enabled" yaml:"enabled" toml:"enabled"`

	// Exclude is a list of CSS file patterns to skip during minification.
	// Useful for files that should not be modified (e.g., already minified vendor CSS).
	// Example: ["variables.css", "vendor/*.css"]
	Exclude []string `json:"exclude" yaml:"exclude" toml:"exclude"`

	// PreserveComments is a list of comment patterns to preserve during minification.
	// Comments containing any of these strings will not be removed.
	// Example: ["/*! Copyright */", "/*! License */"]
	PreserveComments []string `json:"preserve_comments" yaml:"preserve_comments" toml:"preserve_comments"`
}

// NewCSSMinifyConfig creates a new CSSMinifyConfig with default values.
func NewCSSMinifyConfig() CSSMinifyConfig {
	return CSSMinifyConfig{
		Enabled:          true, // Enabled by default for performance
		Exclude:          []string{},
		PreserveComments: []string{},
	}
}

// JSMinifyConfig configures the js_minify plugin for minifying JavaScript files.
// JS minification reduces file sizes by 40-60%, improving page load performance
// and Lighthouse scores. Already-minified files (*.min.js) are automatically skipped.
type JSMinifyConfig struct {
	// Enabled controls whether JS minification is active (default: true)
	Enabled bool `json:"enabled" yaml:"enabled" toml:"enabled"`

	// Exclude is a list of JS file patterns to skip during minification.
	// Useful for files that should not be modified (e.g., vendor scripts).
	// Files ending in .min.js are always excluded automatically.
	// Example: ["pagefind-ui.js", "vendor/*.js"]
	Exclude []string `json:"exclude" yaml:"exclude" toml:"exclude"`
}

// NewJSMinifyConfig creates a new JSMinifyConfig with default values.
func NewJSMinifyConfig() JSMinifyConfig {
	return JSMinifyConfig{
		Enabled: true, // Enabled by default for performance
		Exclude: []string{},
	}
}

// ResourcesConfig configures conditional resource loading behavior.
// By default, CSS and JS resources are loaded conditionally based on page content
// detection (e.g., admonitions.css only loads when a page has admonitions).
// This config allows overriding that behavior.
type ResourcesConfig struct {
	// Strategy controls the overall resource loading approach:
	// - "conditional" (default): Load resources only when page content requires them
	// - "eager": Load all resources on every page (legacy behavior, simpler but slower)
	Strategy string `json:"strategy,omitempty" yaml:"strategy,omitempty" toml:"strategy,omitempty"`

	// ForceLoad is a list of resource base names to always load regardless of content.
	// Useful when custom templates use CSS classes that build-time detection can't see.
	// Example: ["admonitions", "code", "cards"]
	// Valid names: admonitions, code, chroma, cards, webmentions, encryption
	ForceLoad []string `json:"force_load,omitempty" yaml:"force_load,omitempty" toml:"force_load,omitempty"`

	// Search configures search resource loading behavior.
	Search ResourceLoadingMode `json:"search,omitempty" yaml:"search,omitempty" toml:"search,omitempty"`

	// GLightbox configures GLightbox resource loading behavior.
	GLightbox ResourceLoadingMode `json:"glightbox,omitempty" yaml:"glightbox,omitempty" toml:"glightbox,omitempty"`
}

// ResourceLoadingMode configures when a specific resource is loaded.
type ResourceLoadingMode struct {
	// Loading controls when the resource is fetched:
	// - "lazy" (default): Load on first user interaction (hover, focus, keyboard shortcut)
	// - "eager": Load immediately on page load
	Loading string `json:"loading,omitempty" yaml:"loading,omitempty" toml:"loading,omitempty"`
}

// NewResourcesConfig creates a new ResourcesConfig with default values.
func NewResourcesConfig() ResourcesConfig {
	return ResourcesConfig{
		Strategy:  "conditional",
		ForceLoad: []string{},
		Search:    ResourceLoadingMode{Loading: "lazy"},
		GLightbox: ResourceLoadingMode{Loading: "lazy"},
	}
}

// IsConditional returns true if resources should be loaded conditionally.
func (r *ResourcesConfig) IsConditional() bool {
	return r.Strategy != "eager"
}

// IsForceLoaded returns true if a resource should always be loaded.
func (r *ResourcesConfig) IsForceLoaded(name string) bool {
	for _, n := range r.ForceLoad {
		if n == name {
			return true
		}
	}
	return false
}

// AssetsConfig configures external CDN asset handling for self-hosting.
// When mode is "self-hosted", external assets (GLightbox, HTMX, Mermaid, etc.)
// are downloaded at build time and served from the site itself.
type AssetsConfig struct {
	// Mode controls how external assets are handled:
	// - "cdn": Always load from external CDN (default, no download)
	// - "self-hosted": Download and serve from local output/assets/vendor/
	// - "auto": Use self-hosted if assets are cached, fall back to CDN
	Mode string `json:"mode,omitempty" yaml:"mode,omitempty" toml:"mode,omitempty"`

	// CacheDir is the directory for caching downloaded assets (default: ".markata/assets-cache")
	CacheDir string `json:"cache_dir,omitempty" yaml:"cache_dir,omitempty" toml:"cache_dir,omitempty"`

	// VerifyIntegrity enables SRI hash verification for downloaded assets (default: true)
	VerifyIntegrity *bool `json:"verify_integrity,omitempty" yaml:"verify_integrity,omitempty" toml:"verify_integrity,omitempty"`

	// OutputDir is the subdirectory in output for vendor assets (default: "assets/vendor")
	OutputDir string `json:"output_dir,omitempty" yaml:"output_dir,omitempty" toml:"output_dir,omitempty"`
}

// NewAssetsConfig creates a new AssetsConfig with default values.
func NewAssetsConfig() AssetsConfig {
	verifyIntegrity := true
	return AssetsConfig{
		Mode:            "self-hosted",
		CacheDir:        ".markata/assets-cache",
		VerifyIntegrity: &verifyIntegrity,
		OutputDir:       "assets/vendor",
	}
}

// IsSelfHosted returns true if assets should be self-hosted.
func (a *AssetsConfig) IsSelfHosted() bool {
	return a.Mode == "self-hosted" || a.Mode == "auto"
}

// IsVerifyIntegrityEnabled returns whether integrity verification is enabled.
// Defaults to true if not explicitly set.
func (a *AssetsConfig) IsVerifyIntegrityEnabled() bool {
	if a.VerifyIntegrity == nil {
		return true
	}
	return *a.VerifyIntegrity
}

// GetCacheDir returns the cache directory, with default if not set.
func (a *AssetsConfig) GetCacheDir() string {
	if a.CacheDir == "" {
		return ".markata/assets-cache"
	}
	return a.CacheDir
}

// GetOutputDir returns the output directory for vendor assets.
func (a *AssetsConfig) GetOutputDir() string {
	if a.OutputDir == "" {
		return "assets/vendor"
	}
	return a.OutputDir
}

// NewConfig creates a new Config with default values.
func NewConfig() *Config {
	return &Config{
		OutputDir:     "output",
		AssetsDir:     "static",
		TemplatesDir:  "templates",
		Templates:     NewTemplatesConfig(),
		Hooks:         []string{"default"},
		DisabledHooks: []string{},
		GlobConfig: GlobConfig{
			Patterns:     []string{"pages/**/*.md", "posts/**/*.md"},
			UseGitignore: true,
		},
		MarkdownConfig: MarkdownConfig{
			Extensions: []string{},
			Highlight:  NewHighlightConfig(),
		},
		Feeds:            []FeedConfig{},
		FeedDefaults:     NewFeedDefaults(),
		Concurrency:      0,
		Theme:            NewThemeConfig(),
		ThemeCalendar:    NewThemeCalendarConfig(),
		PostFormats:      NewPostFormatsConfig(),
		SEO:              NewSEOConfig(),
		IndieAuth:        NewIndieAuthConfig(),
		Webmention:       NewWebmentionConfig(),
		WebSub:           NewWebSubConfig(),
		Components:       NewComponentsConfig(),
		Search:           NewSearchConfig(),
		Layout:           NewLayoutConfig(),
		Sidebar:          NewSidebarConfig(),
		Toc:              NewTocConfig(),
		Header:           NewHeaderLayoutConfig(),
		FooterLayout:     NewFooterLayoutConfig(),
		ContentTemplates: NewContentTemplatesConfig(),
		Blogroll:         NewBlogrollConfig(),
		Mentions:         NewMentionsConfig(),
		ErrorPages:       NewErrorPagesConfig(),
		ResourceHints:    NewResourceHintsConfig(),
		Encryption:       NewEncryptionConfig(),
		Shortcuts:        NewShortcutsConfig(),
		ViewTransitions:  NewViewTransitionsConfig(),
		Tags:             NewTagsConfig(),
		FeedsPage:        NewFeedsPageConfig(),
		Garden:           NewGardenConfig(),
		Images:           NewImagesConfig(),
		Assets:           NewAssetsConfig(),
	}
}

// LicenseOption returns the configured license metadata when a valid key is present.
func (c *Config) LicenseOption() (LicenseOption, bool) {
	if c == nil {
		return LicenseOption{}, false
	}
	if key, ok := c.License.Key(); ok {
		return GetLicenseOption(key)
	}
	return LicenseOption{}, false
}

// NeedsLicenseWarning reports whether the license reminder should be shown to users.
func (c *Config) NeedsLicenseWarning() bool {
	if c == nil {
		return false
	}
	if c.License.IsDisabled() {
		return false
	}
	if _, ok := c.License.Key(); ok {
		return false
	}
	return !c.License.HasValue()
}

// NewThemeConfig creates a new ThemeConfig with default values.
func NewThemeConfig() ThemeConfig {
	return ThemeConfig{
		ContractVersion: 1,
		Name:            "default",
		Palette:         "ayu-dark",
		Aesthetic:       "minimal",
		FallbackMode:    "dark",
		TextSize:        TextSizeLarge,
		Variables:       make(map[string]string),
		Texture:         ThemeTextureConfig{Kind: "screenprint", ColorMix: 0.35, Scale: 1, Scope: "all"},
		HeadingTexture:  ThemeHeadingTextureConfig{Kind: "inherit", ColorMix: 0.45, Scale: 1},
		Motif:           ThemeMotifConfig{Kind: "off", Glyph: "W", Size: "78px", Gap: "10px", RowOffset: 0.24, Wobble: 0.18, Layer: "sandwich", Color: "ink", ColorMix: 0.01},
		Font:            NewFontConfig(),
		Switcher:        NewThemeSwitcherConfig(),
	}
}

// NewThemeCalendarConfig creates a new ThemeCalendarConfig with default values.
func NewThemeCalendarConfig() ThemeCalendarConfig {
	enabled := false
	return ThemeCalendarConfig{
		Enabled: &enabled,
		Rules:   []ThemeCalendarRule{},
	}
}

// ExternalCacheDirs returns the list of Tier 2 external plugin cache directories.
// These directories contain expensive-to-rebuild data (fetched RSS feeds, embed
// metadata, webmentions) and are only cleaned by --clean-all, not --clean.
// For plugins configured via the Extra map (embeds, webmentions), defaults are
// included when the plugin isn't explicitly configured, since the plugins still
// create these directories at runtime.
func (c *Config) ExternalCacheDirs() []string {
	var dirs []string

	// Blogroll cache (typed config field, falls back to default)
	blogrollDir := c.Blogroll.CacheDir
	if blogrollDir == "" {
		blogrollDir = NewBlogrollConfig().CacheDir
	}
	dirs = append(dirs, blogrollDir)

	// Mentions cache (typed config field with getter that handles default)
	if dir := c.Mentions.GetCacheDir(); dir != "" {
		dirs = append(dirs, dir)
	}

	// Embeds cache (untyped Extra map — stored as map[string]interface{})
	embedsDir := NewEmbedsConfig().CacheDir // default: ".cache/embeds"
	if c.Extra != nil {
		if embedsCfg, ok := c.Extra["embeds"]; ok {
			if cfgMap, ok := embedsCfg.(map[string]interface{}); ok {
				if dir, ok := cfgMap["cache_dir"].(string); ok && dir != "" {
					embedsDir = dir
				}
			}
		}
	}
	dirs = append(dirs, embedsDir)

	// Webmentions cache (untyped Extra map — stored as WebMentionsConfig)
	wmDir := NewWebMentionsConfig().CacheDir // default: ".cache/webmentions"
	if c.Extra != nil {
		if wmCfg, ok := c.Extra["webmentions"]; ok {
			if wm, ok := wmCfg.(WebMentionsConfig); ok {
				if wm.CacheDir != "" {
					wmDir = wm.CacheDir
				}
			}
		}
	}
	dirs = append(dirs, wmDir)

	return dirs
}

// IsHookEnabled checks if a hook is enabled (in Hooks and not in DisabledHooks).
func (c *Config) IsHookEnabled(name string) bool {
	// Check if disabled
	for _, h := range c.DisabledHooks {
		if h == name {
			return false
		}
	}

	// Check if enabled
	for _, h := range c.Hooks {
		if h == name || h == "default" {
			return true
		}
	}

	return false
}
