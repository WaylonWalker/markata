package plugins

import (
	"fmt"
	"html"
	"regexp"
	"strings"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

// InlineTitlesPlugin derives safe rich and plain-text representations from
// authored title source after auto-title inference and before metadata
// plugins run.
type InlineTitlesPlugin struct{}

// NewInlineTitlesPlugin creates a title representation plugin.
func NewInlineTitlesPlugin() *InlineTitlesPlugin { return &InlineTitlesPlugin{} }

// Name returns the unique plugin name.
func (p *InlineTitlesPlugin) Name() string { return "inline_titles" }

// Priority places title rendering immediately after auto-title inference.
func (p *InlineTitlesPlugin) Priority(stage lifecycle.Stage) int {
	if stage == lifecycle.StageTransform {
		return lifecycle.PriorityFirst + 1
	}
	return lifecycle.PriorityDefault
}

// Transform populates TitleHTML and TitleText without changing Title, which
// remains the authored/source title for compatibility.
func (p *InlineTitlesPlugin) Transform(m *lifecycle.Manager) error {
	baseline := takeSemanticHashBaseline(m)
	value, ok := m.Cache().Get(CacheKeyInlineRenderer)
	if !ok {
		renderer := NewRenderMarkdownPlugin()
		if err := renderer.Configure(m); err != nil {
			return fmt.Errorf("configure inline title renderer: %w", err)
		}
		value, _ = m.Cache().Get(CacheKeyInlineRenderer)
	}
	renderInline, ok := value.(InlineRenderFunc)
	if !ok {
		return fmt.Errorf("inline title renderer is unavailable")
	}
	posts := m.FilterPosts(func(post *models.Post) bool {
		return !post.Skip && post.Title != nil
	})
	return m.ProcessPostsSliceConcurrently(posts, func(post *models.Post) error {
		result, err := renderInline(*post.Title)
		if err != nil {
			return err
		}
		post.TitleHTML = result.HTML
		post.TitleHTML = wrapHeadingHighlights(post.TitleHTML)
		post.TitleText = result.Text
		post.TitleTextDerived = true
		if post.Slug == "" && !post.Has("_slug_explicit") {
			post.GenerateSlugWithMode(configuredSlugMode(m.Config(), post.Path))
			post.GenerateHref()
		}
		if cache := GetBuildCache(m); cache != nil {
			cache.SetPostSlug(post.Path, post.Slug)
			current := semanticHashes{
				feed:   computePostFeedItemHash(post),
				tag:    computePostTagIndexHash(post),
				garden: computePostGardenHash(post),
			}
			feedChanged, tagChanged, gardenChanged := cache.UpdatePostSemanticHashes(
				post.Path, current.feed, current.tag, current.garden,
			)
			if previous := baseline[post.Path]; previous != nil {
				feedChanged = previous.feed != current.feed
				tagChanged = previous.tag != current.tag
				gardenChanged = previous.garden != current.garden
			}
			if feedChanged {
				cache.MarkFeedSlugChanged(post.Slug)
			}
			if tagChanged {
				cache.MarkSlugChanged(post.Slug)
			}
			if gardenChanged {
				cache.MarkSlugChanged(post.Slug)
			}
		}
		return nil
	})
}

var (
	openingMarkTagPattern = regexp.MustCompile(`(?i)<mark(\s[^>]*)?>`)
	closingMarkTagPattern = regexp.MustCompile(`(?i)</mark\s*>`)
	headingBlockPattern   = regexp.MustCompile(`(?is)<h([12])([^>]*)>(.*?)</h[12]>`)
	headingCommentPattern = regexp.MustCompile(`(?s)<!--.*?-->`)
)

// hasHeadingMarkHighlights classifies canonical article HTML, not title HTML
// or later derivatives. The fast gate avoids scanning ordinary articles.
func hasHeadingMarkHighlights(value string) bool {
	if !hasMarkTagPrefix(value) {
		return false
	}
	if strings.Contains(value, "<!--") {
		value = headingCommentPattern.ReplaceAllString(value, "")
	}
	for _, parts := range headingBlockPattern.FindAllStringSubmatch(value, -1) {
		if openingMarkTagPattern.MatchString(parts[3]) {
			return true
		}
	}
	return false
}

// hasMarkTagPrefix checks ASCII HTML tag names without copying the article or
// allocating regexp execution state (including under the race detector).
func hasMarkTagPrefix(value string) bool {
	for {
		start := strings.IndexByte(value, '<')
		if start < 0 || len(value)-start < len("<mark") {
			return false
		}
		if strings.EqualFold(value[start:start+len("<mark")], "<mark") {
			return true
		}
		value = value[start+1:]
	}
}

// wrapHeadingHighlights adds a semantic hook around marks used in title HTML.
// The mark remains the no-JavaScript fallback; the heading-highlights script
// can replace its background with a measured multi-line contour.
func wrapHeadingHighlights(value string) string {
	if !strings.Contains(strings.ToLower(value), "<mark") {
		return value
	}
	value = openingMarkTagPattern.ReplaceAllString(value, `<span class="heading-highlight">$0`)
	return closingMarkTagPattern.ReplaceAllString(value, `$0</span>`)
}

// wrapHeadingMarkHighlights adds the same semantic hook to Markdown headings
// in article HTML. Prose marks retain their ordinary fallback styling.
func wrapHeadingMarkHighlights(value string) string {
	return headingBlockPattern.ReplaceAllStringFunc(value, func(block string) string {
		parts := headingBlockPattern.FindStringSubmatch(block)
		if len(parts) != 4 {
			return block
		}
		return "<h" + parts[1] + parts[2] + ">" + wrapHeadingHighlights(parts[3]) + "</h" + parts[1] + ">"
	})
}

// PopulateTitleFallback provides a safe representation for manually-created
// posts when the normal lifecycle renderer is not installed.
func PopulateTitleFallback(post *models.Post) {
	if post == nil || post.Title == nil || post.TitleTextDerived {
		return
	}
	post.TitleText = html.UnescapeString(*post.Title)
	post.TitleHTML = html.EscapeString(post.TitleText)
	post.TitleTextDerived = true
}

var (
	_ lifecycle.Plugin          = (*InlineTitlesPlugin)(nil)
	_ lifecycle.TransformPlugin = (*InlineTitlesPlugin)(nil)
	_ lifecycle.PriorityPlugin  = (*InlineTitlesPlugin)(nil)
)
