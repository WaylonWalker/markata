package workbench

import (
	"html"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/WaylonWalker/markata/pkg/lsp"
)

var (
	workbenchInternalEmbed = regexp.MustCompile(`!\[\[([^\]|]+)(?:\|([^\]]+))?\]\]`)
	workbenchMarkdownEmbed = regexp.MustCompile(`!\[embed\]\(([^)]+)\)`)
)

// renderWorkbenchEmbeds resolves Markata post embeds against the current site
// index. It supports the site's ![[slug]] form and the Markdown image form
// used by the editor's file completion, ![embed](path/to/post.md).
func renderWorkbenchEmbeds(source string, index *lsp.Index, root string) string {
	if index == nil || (!strings.Contains(source, "![embed](") && !strings.Contains(source, "![[")) {
		return source
	}

	return replaceOutsideWorkbenchCode(source, func(text string) string {
		text = workbenchMarkdownEmbed.ReplaceAllStringFunc(text, func(match string) string {
			parts := workbenchMarkdownEmbed.FindStringSubmatch(match)
			return renderWorkbenchEmbedCard(match, parts[1], "", root, index)
		})
		return workbenchInternalEmbed.ReplaceAllStringFunc(text, func(match string) string {
			parts := workbenchInternalEmbed.FindStringSubmatch(match)
			display := ""
			if len(parts) > 2 {
				display = strings.TrimSpace(parts[2])
			}
			return renderWorkbenchEmbedCard(match, strings.TrimSpace(parts[1]), display, root, index)
		})
	})
}

func renderWorkbenchEmbedCard(original, target, display, root string, index *lsp.Index) string {
	post := resolveWorkbenchEmbedTarget(target, root, index)
	if post == nil || post.Slug == "" {
		return original
	}
	title := strings.TrimSpace(display)
	if title == "" {
		title = strings.TrimSpace(post.Title)
	}
	if title == "" {
		title = post.Slug
	}
	var card strings.Builder
	card.WriteString(`<aside class="embed-card workbench-post-embed"><a class="wikilink embed-card-title" href="/`)
	card.WriteString(html.EscapeString(post.Slug))
	card.WriteString(`/"><strong>`)
	card.WriteString(html.EscapeString(title))
	card.WriteString(`</strong></a>`)
	if description := strings.TrimSpace(post.Description); description != "" {
		card.WriteString(`<p class="embed-card-description">`)
		card.WriteString(html.EscapeString(description))
		card.WriteString(`</p>`)
	}
	card.WriteString(`</aside>`)
	return card.String()
}

func resolveWorkbenchEmbedTarget(target, root string, index *lsp.Index) *lsp.PostInfo {
	target = strings.TrimSpace(target)
	if target == "" || strings.Contains(target, "://") {
		return nil
	}
	if parsed, err := url.Parse(target); err == nil {
		target = parsed.Path
	}
	if decoded, err := url.PathUnescape(target); err == nil {
		target = decoded
	}
	cleanTarget := filepath.Clean(filepath.FromSlash(strings.TrimPrefix(target, "/")))
	if cleanTarget != "." && cleanTarget != ".." && !strings.HasPrefix(cleanTarget, ".."+string(filepath.Separator)) {
		for _, post := range index.AllPosts() {
			relative, err := filepath.Rel(root, post.Path)
			if err != nil {
				continue
			}
			if filepath.Clean(relative) == cleanTarget {
				return post
			}
		}
	}
	if filepath.Ext(target) == "" {
		return index.GetBySlug(filepath.Base(target))
	}
	return nil
}

func replaceOutsideWorkbenchCode(source string, replace func(string) string) string {
	blocks := workbenchCodeFence.FindAllStringIndex(source, -1)
	if len(blocks) == 0 {
		return replace(source)
	}
	var result strings.Builder
	last := 0
	for _, block := range blocks {
		result.WriteString(replace(source[last:block[0]]))
		result.WriteString(source[block[0]:block[1]])
		last = block[1]
	}
	result.WriteString(replace(source[last:]))
	return result.String()
}
