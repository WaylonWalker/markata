package workbench

import (
	"html"
	"regexp"
	"strings"

	"github.com/WaylonWalker/markata/pkg/lsp"
)

var (
	workbenchWikilink  = regexp.MustCompile(`\[\[([^\]|]+)(?:\|([^\]]+))?\]\]`)
	workbenchCodeFence = regexp.MustCompile("(?s)(```[^`]*```|~~~[^~]*~~~)")
)

// renderWorkbenchWikilinks resolves links against the workspace snapshot, as
// the site build does during its transform stage. Unresolved links remain
// visible in source form so the LSP diagnostic can explain them.
func renderWorkbenchWikilinks(source string, index *lsp.Index) string {
	if index == nil || !strings.Contains(source, "[[") {
		return source
	}
	blocks := workbenchCodeFence.FindAllStringIndex(source, -1)
	if len(blocks) == 0 {
		return replaceWorkbenchWikilinks(source, index)
	}
	var result strings.Builder
	last := 0
	for _, block := range blocks {
		result.WriteString(replaceWorkbenchWikilinks(source[last:block[0]], index))
		result.WriteString(source[block[0]:block[1]])
		last = block[1]
	}
	result.WriteString(replaceWorkbenchWikilinks(source[last:], index))
	return result.String()
}

func replaceWorkbenchWikilinks(source string, index *lsp.Index) string {
	return workbenchWikilink.ReplaceAllStringFunc(source, func(match string) string {
		parts := workbenchWikilink.FindStringSubmatch(match)
		target := strings.TrimSpace(parts[1])
		fragment := ""
		if at := strings.IndexByte(target, '#'); at >= 0 {
			fragment, target = target[at:], target[:at]
		}
		post := index.GetBySlug(target)
		if post == nil {
			return match
		}
		label := strings.TrimSpace(parts[2])
		if label == "" {
			label = post.Title
		}
		if label == "" {
			label = post.Slug
		}
		return `<a class="wikilink" href="/` + html.EscapeString(post.Slug) + `/` + html.EscapeString(fragment) + `">` + html.EscapeString(label) + `</a>`
	})
}
