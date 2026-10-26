package plugins

import (
	"fmt"
	"html"
	"math"
	"regexp"
	"strings"

	"github.com/WaylonWalker/markata/pkg/filter"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/logging"
	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/templates"
	"github.com/flosch/pongo2/v6"
)

var renderFeedLog = logging.Component("jinja_md").Phase("transform")

const defaultFeedEmbedTemplate = "partials/feed_preview.html"
const defaultFeedSortField = "date"

const (
	nilString             = "<nil>"
	templateTypePost      = "post"
	templateTypeLink      = "link"
	templateTypeBookmark  = "bookmark"
	templateTypeNote      = "note"
	templateTypeInline    = "inline"
	templateTypeDefault   = "default"
	templateTypePhoto     = "photo"
	templateTypeVideo     = "video"
	templateTypeGuide     = "guide"
	templateTypeQuote     = "quote"
	templateTypeContact   = "contact"
	templateTypeArticle   = "article"
	templateTypeGratitude = "gratitude"
	templateTypeTutorial  = "tutorial"
)

// ensureArticleHTML renders markdown on-demand for any feed post whose
// ArticleHTML has not yet been populated (e.g. when render_feed is called
// from jinja_md during StageTransform, before render_markdown runs).
// The rendered HTML is stored on the post so later stages can reuse it.
func ensureArticleHTML(posts []*models.Post, m *lifecycle.Manager) {
	cached, ok := m.Cache().Get(CacheKeyMarkdownRenderer)
	if !ok {
		return
	}
	renderFn, ok := cached.(MarkdownRenderFunc)
	if !ok {
		return
	}

	for _, post := range posts {
		if post.ArticleHTML == "" && post.Content != "" && !post.Skip {
			if rendered, err := renderFn(post.Content); err == nil {
				post.ArticleHTML = rendered
				templates.InvalidatePost(post)
			}
		}
	}
}

// createFeedPostsFunc exposes a helper that returns the latest posts for a feed.
func createFeedPostsFunc(m *lifecycle.Manager) func(slug string, args ...interface{}) ([]map[string]interface{}, error) {
	return func(slug string, args ...interface{}) ([]map[string]interface{}, error) {
		limit := parseLimitArg(args)
		posts, _ := getFeedPosts(slug, limit, m)
		if len(posts) == 0 {
			return []map[string]interface{}{}, nil
		}
		ensureArticleHTML(posts, m)
		return templates.PostsToMaps(posts), nil
	}
}

// createRenderFeedFunc renders a feed snippet as safe HTML.
func createRenderFeedFunc(m *lifecycle.Manager) func(slug string, args ...interface{}) (*pongo2.Value, error) {
	return func(slug string, args ...interface{}) (*pongo2.Value, error) {
		opts := parseRenderFeedArgs(args)
		posts, fc := getFeedPosts(slug, opts.limit, m)
		if fc == nil {
			renderFeedLog.Warnf("render_feed(%q): no feed with that slug is configured; rendering nothing", slug)
			return pongo2.AsSafeValue(""), nil
		}
		if len(posts) == 0 {
			return pongo2.AsSafeValue(""), nil
		}

		ensureArticleHTML(posts, m)
		postMaps := templates.PostsToMaps(posts)
		ctx := map[string]interface{}{
			"posts":   postMaps,
			"variant": opts.variant,
			"feed": map[string]interface{}{
				"slug":        fc.Slug,
				"title":       fc.Title,
				"description": fc.Description,
			},
		}

		// Include config so card-router can access card_router.mappings
		if cfg := m.Config(); cfg != nil {
			modelsConfig := ToModelsConfig(cfg)
			if modelsConfig != nil {
				ctx["config"] = templates.GetConfigMap(modelsConfig)
			}
		}

		htmlSnippet, tmplErr := renderFeedWithTemplate(ctx, opts, m)
		if tmplErr != nil {
			htmlSnippet = fallbackFeedHTML(postMaps)
		}

		return pongo2.AsSafeValue(htmlSnippet), nil
	}
}

func parseLimitArg(args []interface{}) int {
	for _, arg := range args {
		if n, ok := numericValue(arg); ok {
			return nonNegativeInt(n)
		}
		if cfg, ok := arg.(map[string]interface{}); ok {
			if limitVal, ok := getNumberFromMap(cfg, "limit"); ok {
				return nonNegativeInt(limitVal)
			}
		}
	}
	return 0
}

type renderFeedArgs struct {
	limit    int
	variant  string
	template string
}

func parseRenderFeedArgs(args []interface{}) renderFeedArgs {
	opts := renderFeedArgs{variant: "card"}
	for _, arg := range args {
		switch v := arg.(type) {
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
			if n, ok := numericValue(v); ok {
				opts.limit = nonNegativeInt(n)
			}
		case string:
			opts.variant = strings.ToLower(v)
		case map[string]interface{}:
			if limitVal, ok := getNumberFromMap(v, "limit"); ok {
				opts.limit = nonNegativeInt(limitVal)
			}
			if variantVal, ok := v["variant"].(string); ok && variantVal != "" {
				opts.variant = strings.ToLower(variantVal)
			}
			if templateVal, ok := v["template"].(string); ok && templateVal != "" {
				opts.template = templateVal
			}
		}
	}
	if opts.variant == "" {
		opts.variant = "card"
	}
	return opts
}

func getFeedPosts(slug string, limit int, m *lifecycle.Manager) ([]*models.Post, *models.FeedConfig) {
	if slug == "" {
		slug = ""
	}

	if posts, fc, ok := postsFromCache(slug, limit, m); ok {
		return posts, fc
	}

	fc := getFeedBySlugFromConfig(slug, m.Config())
	if fc == nil {
		return nil, nil
	}

	posts := filterFeedPagePosts(computeFeedPosts(fc, m), fc.IncludesPrivate())

	if limit > 0 && limit < len(posts) {
		posts = posts[:limit]
	}

	return posts, fc
}

func postsFromCache(slug string, limit int, m *lifecycle.Manager) ([]*models.Post, *models.FeedConfig, bool) {
	cached, ok := m.Cache().Get("feed_configs")
	if !ok {
		return nil, nil, false
	}

	configs, ok := cached.([]models.FeedConfig)
	if !ok {
		return nil, nil, false
	}

	fc := GetFeedBySlug(slug, configs)
	if fc == nil {
		return nil, nil, false
	}

	posts := filterFeedPagePosts(fc.Posts, fc.IncludesPrivate())
	if len(posts) == 0 {
		return nil, fc, true
	}

	if limit > 0 && limit < len(posts) {
		posts = posts[:limit]
	}

	return posts, fc, true
}

func getFeedBySlugFromConfig(slug string, cfg *lifecycle.Config) *models.FeedConfig {
	if cfg == nil {
		return nil
	}
	configs := getFeedConfigs(cfg)
	if len(configs) == 0 {
		return nil
	}
	return GetFeedBySlug(slug, configs)
}

func computeFeedPosts(fc *models.FeedConfig, m *lifecycle.Manager) []*models.Post {
	if fc == nil {
		return nil
	}

	posts := filterPostsForFeed(fc, m)
	if len(posts) == 0 {
		return posts
	}

	sortField, reverse := feedSort(fc)
	sortPosts(posts, sortField, reverse)
	posts = applyFeedLimitOffset(posts, fc)
	return posts
}

func filterPostsForFeed(fc *models.FeedConfig, m *lifecycle.Manager) []*models.Post {
	candidate := make([]*models.Post, 0, len(m.Posts()))
	for _, post := range m.Posts() {
		if fc.IncludesPrivate() || !post.Private {
			candidate = append(candidate, post)
		}
	}

	if fc.Filter == "" {
		return cloneFeedPosts(candidate)
	}

	parsed, err := filter.Parse(fc.Filter)
	if err != nil {
		return cloneFeedPosts(candidate)
	}

	return cloneFeedPosts(parsed.MatchAll(candidate))
}

func feedSort(fc *models.FeedConfig) (string, bool) {
	sortField := fc.Sort
	reverse := fc.Reverse
	if fc.Type == models.FeedTypeGuide {
		if sortField == "" {
			sortField = "guide_order"
			reverse = false
		}
	}
	if sortField == "" {
		sortField = defaultFeedSortField
		reverse = true
	}
	return sortField, reverse
}

func renderFeedWithTemplate(ctx map[string]interface{}, opts renderFeedArgs, m *lifecycle.Manager) (string, error) {
	engine, err := ensureTemplateEngine(m)
	if err != nil {
		return "", err
	}

	templateName := opts.template
	if templateName == "" {
		templateName = defaultFeedEmbedTemplate
	}

	if !engine.TemplateExists(templateName) {
		templateName = defaultFeedEmbedTemplate
		if !engine.TemplateExists(templateName) {
			return "", fmt.Errorf("template %q not found", templateName)
		}
	}

	result, err := engine.RenderToString(templateName, ctx)
	if err != nil {
		return "", err
	}
	return normalizeTemplateWhitespace(result), nil
}

// normalizeTemplateWhitespace cleans up whitespace artifacts from pongo2 template
// rendering. Pongo2 replaces {% %} tags with empty strings but preserves surrounding
// whitespace, creating blank lines and whitespace-only lines. When this HTML is
// injected into markdown (via render_feed), goldmark treats the snippet as an HTML
// block that ends at the first blank line; everything after it is parsed as
// Markdown again (list markers, ==marks==, indented code blocks, ...).
//
// This function:
//  1. Left-trims every line (removes leading whitespace) so no line starts with 4+ spaces
//  2. Drops blank lines entirely so the snippet stays a single HTML block
func normalizeTemplateWhitespace(s string) string {
	lines := strings.Split(s, "\n")
	out := make([]string, 0, len(lines))
	for _, line := range lines {
		trimmed := strings.TrimLeft(line, " \t")
		if trimmed == "" {
			continue
		}
		out = append(out, trimmed)
	}
	return strings.Join(out, "\n")
}

func getTemplateEngine(m *lifecycle.Manager) *templates.Engine {
	cached, ok := m.Cache().Get("templates.engine")
	if !ok {
		return nil
	}
	engine, ok := cached.(*templates.Engine)
	if !ok {
		return nil
	}
	return engine
}

func ensureTemplateEngine(m *lifecycle.Manager) (*templates.Engine, error) {
	if engine := getTemplateEngine(m); engine != nil {
		return engine, nil
	}

	cfg := m.Config()
	if cfg == nil {
		return nil, fmt.Errorf("config unavailable")
	}

	templatesDir := PluginNameTemplates
	if cfg.Extra != nil {
		if dir, ok := cfg.Extra["templates_dir"].(string); ok && dir != "" {
			templatesDir = dir
		}
	}

	themeName := resolveThemeName(cfg)
	engine, err := templates.NewEngineWithTheme(templatesDir, themeName)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize template engine: %w", err)
	}

	m.Cache().Set("templates.engine", engine)
	return engine, nil
}

func resolveThemeName(cfg *lifecycle.Config) string {
	themeName := ThemeDefault
	if cfg == nil || cfg.Extra == nil {
		return themeName
	}

	if raw, ok := cfg.Extra["theme"]; ok {
		switch typed := raw.(type) {
		case models.ThemeConfig:
			if typed.Name != "" {
				return typed.Name
			}
		case *models.ThemeConfig:
			if typed != nil && typed.Name != "" {
				return typed.Name
			}
		case map[string]interface{}:
			if name, ok := typed["name"].(string); ok && name != "" {
				return name
			}
		case string:
			if typed != "" {
				return typed
			}
		}
	}

	return themeName
}

func fallbackFeedHTML(posts []map[string]interface{}) string {
	builder := &strings.Builder{}
	builder.WriteString("<div class=\"feed h-feed\">")
	builder.WriteString("<div class=\"posts posts-list\" id=\"posts-list\">")
	for _, post := range posts {
		templateType := cardTypeForPost(post)
		if templateType == templateTypePhoto {
			builder.WriteString(renderPhotoFigure(post))
			continue
		}
		fmt.Fprintf(builder, "<article class=\"card card-%s h-entry\">", html.EscapeString(templateType))
		if templateType == templateTypeLink {
			builder.WriteString(renderLinkHeader(post))
		} else {
			builder.WriteString(renderFeedLink(post))
		}
		if templateType == templateTypeLink {
			builder.WriteString(renderLinkSnippet(post))
		} else if snippet := renderGenericSnippet(post, templateType); snippet != "" {
			builder.WriteString(snippet)
		}
		builder.WriteString(renderFeedMeta(post))
		builder.WriteString("</article>")
	}
	builder.WriteString("</div>")
	builder.WriteString("</div>")
	return builder.String()
}

func renderFeedLink(post map[string]interface{}) string {
	href := fmt.Sprint(post["href"])
	if href == "" {
		href = "/"
	}
	title := fmt.Sprint(post["title"])
	if title == "" || title == nilString {
		title = fmt.Sprint(post["slug"])
	}
	builder := &strings.Builder{}
	fmt.Fprintf(builder, "<a href=\"%s\"><strong>%s</strong></a>", html.EscapeString(href), html.EscapeString(title))
	return builder.String()
}

func renderLinkHeader(post map[string]interface{}) string {
	href := fmt.Sprint(post["href"])
	if href == "" {
		href = "/"
	}
	title := fmt.Sprint(post["title"])
	if title == "" || title == nilString {
		title = fmt.Sprint(post["slug"])
	}
	b := &strings.Builder{}
	b.WriteString("<div class=\"card-link-wrapper\">")
	b.WriteString("<div class=\"card-link-content\">")
	fmt.Fprintf(b, "<a class=\"card-title p-name u-url\" href=\"%s\">%s</a>", html.EscapeString(href), html.EscapeString(title))
	b.WriteString("</div></div>")
	return b.String()
}

func renderFeedMeta(post map[string]interface{}) string {
	date := fmt.Sprint(post["date"])
	tags := extractTags(post)
	b := &strings.Builder{}
	b.WriteString("<footer class=\"card-meta\">")
	if date != "" && date != nilString {
		fmt.Fprintf(b, "<time>%s</time>", html.EscapeString(date))
	}
	if len(tags) > 0 {
		b.WriteString("<div class=\"card-tags\">")
		for _, tag := range tags {
			if tag == "" {
				continue
			}
			slug := models.Slugify(tag)
			fmt.Fprintf(b, "<a href=\"/tags/%s/\" class=\"tag p-category\">%s</a>", html.EscapeString(slug), html.EscapeString(tag))
		}
		b.WriteString("</div>")
	}
	b.WriteString("</footer>")
	return b.String()
}

func extractTags(post map[string]interface{}) []string {
	if raw, ok := post["tags"]; ok && raw != nil {
		switch typed := raw.(type) {
		case []string:
			return typed
		case []interface{}:
			result := make([]string, 0, len(typed))
			for _, item := range typed {
				if s, ok := item.(string); ok {
					result = append(result, s)
				}
			}
			return result
		case string:
			if typed != "" {
				return []string{typed}
			}
		}
	}
	return nil
}

func cardTypeForPost(post map[string]interface{}) string {
	template := strings.ToLower(firstNonEmpty(post, "template", "templateKey"))
	switch template {
	case "blog-post", templateTypeArticle, templateTypePost, "essay", templateTypeTutorial:
		return templateTypeArticle
	case templateTypeNote, "ping", "thought", "status", "tweet":
		return templateTypeNote
	case templateTypePhoto, "shot", "shots", "image", templateTypeGallery:
		return templateTypePhoto
	case templateTypeVideo, "clip", "cast", "stream":
		return templateTypeVideo
	case templateTypeLink, templateTypeBookmark, "til", "stars":
		return templateTypeLink
	case templateTypeQuote, "quotation":
		return templateTypeQuote
	case templateTypeGuide, seriesKey, "step", "chapter":
		return templateTypeGuide
	case templateTypeGratitude, templateTypeInline, "micro":
		return templateTypeInline
	case templateTypeContact, "character", "person":
		return templateTypeContact
	default:
		return templateTypeDefault
	}
}

var htmlTagStripper = regexp.MustCompile(`(?s)<[^>]*>`)

func renderLinkSnippet(post map[string]interface{}) string {
	preview := firstNonEmpty(post, "article_html", "html", "content", "description")
	if preview == "" {
		return ""
	}
	cleaned := htmlTagStripper.ReplaceAllString(preview, "")
	cleaned = html.UnescapeString(cleaned)
	cleaned = strings.TrimSpace(cleaned)
	if cleaned == "" {
		return ""
	}
	cleaned = truncateWords(cleaned, 60)
	b := &strings.Builder{}
	b.WriteString("<div class=\"card-link-body\">")
	fmt.Fprintf(b, "<p class=\"card-link-snippet p-summary\">%s</p>", html.EscapeString(cleaned))
	b.WriteString("</div>")
	return b.String()
}

func renderGenericSnippet(post map[string]interface{}, templateType string) string {
	var source string
	switch templateType {
	case templateTypeArticle:
		source = firstNonEmpty(post, "article_html", "html", "description", "content")
	default:
		source = firstNonEmpty(post, "description", "article_html", "html", "content")
	}
	if source == "" {
		return ""
	}
	cleaned := htmlTagStripper.ReplaceAllString(source, "")
	cleaned = html.UnescapeString(cleaned)
	cleaned = strings.TrimSpace(cleaned)
	if cleaned == "" {
		return ""
	}
	limit := 45
	if templateType == templateTypeArticle {
		limit = 70
	}
	cleaned = truncateWords(cleaned, limit)
	b := &strings.Builder{}
	b.WriteString("<div class=\"card-body\">")
	fmt.Fprintf(b, "<p class=\"card-description p-summary\">%s</p>", html.EscapeString(cleaned))
	b.WriteString("</div>")
	return b.String()
}

func truncateWords(text string, limit int) string {
	if limit <= 0 {
		return text
	}
	words := strings.Fields(text)
	if len(words) <= limit {
		return text
	}
	return strings.Join(words[:limit], " ") + "..."
}

func renderPhotoFigure(post map[string]interface{}) string {
	href := fmt.Sprint(post["href"])
	if href == "" {
		href = "/"
	}
	image := firstNonEmpty(post, "image", "cover", "cover_image", "og_image")
	if image == "" {
		slug := firstNonEmpty(post, "slug")
		if slug != "" {
			image = "/" + strings.Trim(slug, "/") + "/og/"
		}
	}
	if image == "" {
		return ""
	}
	caption := firstNonEmpty(post, "description", "title", "slug")
	// Photo previews should keep the source aspect ratio, unlike intentional cover crops.
	image = templates.WithSize(image, 1200, 0)
	builder := &strings.Builder{}
	builder.WriteString("<figure class=\"photo-figure h-entry\">")
	fmt.Fprintf(builder, "<a href=\"%s\" aria-label=\"%s\"><img src=\"%s\" alt=\"%s\" loading=\"lazy\" width=\"1200\"></a>", html.EscapeString(href), html.EscapeString(caption), html.EscapeString(image), html.EscapeString(caption))
	fmt.Fprintf(builder, "<figcaption class=\"p-summary\">%s</figcaption>", html.EscapeString(caption))
	builder.WriteString("</figure>")
	return builder.String()
}

func firstNonEmpty(post map[string]interface{}, keys ...string) string {
	for _, key := range keys {
		v := fmt.Sprint(post[key])
		if v != "" && v != nilString {
			return v
		}
	}
	return ""
}

func getNumberFromMap(opts map[string]interface{}, key string) (int, bool) {
	if raw, ok := opts[key]; ok {
		return numericValue(raw)
	}
	return 0, false
}

func numericValue(value interface{}) (int, bool) {
	if n, ok := toInt(value); ok {
		return n, true
	}
	switch v := value.(type) {
	case uint:
		return clampUint64ToInt(uint64(v)), true
	case uint8:
		return clampUint64ToInt(uint64(v)), true
	case uint16:
		return clampUint64ToInt(uint64(v)), true
	case uint32:
		return clampUint64ToInt(uint64(v)), true
	case uint64:
		return clampUint64ToInt(v), true
	default:
		return 0, false
	}
}

func clampUint64ToInt(value uint64) int {
	if value > math.MaxInt {
		return math.MaxInt
	}
	return int(value)
}

func nonNegativeInt(value int) int {
	if value < 0 {
		return 0
	}
	return value
}

// createRenderSlashesFunc returns a Jinja-callable that renders slash pages as
// a grid of pill links.  It queries the "slashes" feed (posts tagged "slash")
// and emits <div class="home-slashes">…</div>.
func createRenderSlashesFunc(m *lifecycle.Manager) func(args ...interface{}) (*pongo2.Value, error) {
	return func(args ...interface{}) (*pongo2.Value, error) {
		// Allow an optional limit argument (default 0 = all)
		limit := 0
		for _, arg := range args {
			if n, ok := numericValue(arg); ok {
				limit = nonNegativeInt(n)
			}
		}

		posts, _ := getFeedPosts("slashes", limit, m)
		if len(posts) == 0 {
			return pongo2.AsSafeValue(""), nil
		}

		var sb strings.Builder
		sb.WriteString(`<div class="home-slashes">`)
		for _, p := range posts {
			// Derive display name: prefer the href (e.g. /now/) to look like a slash page
			display := p.Href
			if display == "" {
				display = "/" + p.Slug + "/"
			}
			sb.WriteString(`<a href="`)
			sb.WriteString(html.EscapeString(p.Href))
			sb.WriteString(`" class="home-slash-link">`)
			sb.WriteString(html.EscapeString(display))
			sb.WriteString(`</a>`)
		}
		sb.WriteString(`</div>`)
		return pongo2.AsSafeValue(sb.String()), nil
	}
}

// createIncludePostFunc returns a Jinja-callable that looks up a post by slug
// and returns its rendered ArticleHTML.  This lets templates and jinja-enabled
// markdown compose pages from other posts:
//
//	{{ include_post("components/sidebar") }}
func createIncludePostFunc(m *lifecycle.Manager) func(slug string) (*pongo2.Value, error) {
	return func(slug string) (*pongo2.Value, error) {
		slug = strings.TrimSpace(slug)
		if slug == "" {
			return pongo2.AsSafeValue(""), nil
		}
		// Normalize: strip leading/trailing slashes
		slug = strings.Trim(slug, "/")

		for _, p := range m.Posts() {
			if p.Slug == slug {
				if p.ArticleHTML != "" {
					return pongo2.AsSafeValue(p.ArticleHTML), nil
				}
				// Fallback to raw content if not yet rendered
				if p.Content != "" {
					return pongo2.AsSafeValue(p.Content), nil
				}
				return pongo2.AsSafeValue(""), nil
			}
		}
		return pongo2.AsSafeValue("<!-- include_post: post not found: " + html.EscapeString(slug) + " -->"), nil
	}
}
