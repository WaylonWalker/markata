package plugins

import (
	"net/url"
	"strings"

	"github.com/WaylonWalker/markata/pkg/buildcache"
	"github.com/WaylonWalker/markata/pkg/models"
)

// linkPostDestination returns an external destination for bookmark-style posts.
// The primary "link" field and the "url" alias are both valid frontmatter.
// Do not dereference destinations belonging to private posts.
func linkPostDestination(post *models.Post) string {
	if post == nil || post.Private || post.Extra == nil {
		return ""
	}
	switch strings.ToLower(strings.TrimSpace(post.Template)) {
	case templateTypeLink, templateTypeBookmark, "stars", "til":
	default:
		return ""
	}

	for _, key := range []string{"link", "url"} {
		destination, ok := post.Extra[key].(string)
		if !ok {
			continue
		}
		destination = strings.TrimSpace(destination)
		if strings.ContainsAny(destination, "\r\n\t<>") {
			continue
		}
		parsed, err := url.Parse(destination)
		if err != nil || (parsed.Scheme != schemeHTTPS && parsed.Scheme != schemeHTTP) || parsed.Host == "" || parsed.User != nil {
			continue
		}
		return destination
	}
	return ""
}

// A changed frontmatter destination must invalidate transformed Markdown even
// if the post's authored body did not change. Keep non-link cache keys intact.
func linkPostEmbedCacheHash(post *models.Post, cacheSignature string) string {
	input := post.Content + cacheSignature
	if destination := linkPostDestination(post); destination != "" {
		input += "\x00link-post-destination:" + destination
	}
	return buildcache.ContentHash(input)
}

// Do not duplicate embeds already authored explicitly, or a standalone URL
// that the existing YouTube / one-line-link plugins will embed during Render.
func hasAuthoredLinkPostEmbed(content, destination string) bool {
	if strings.Contains(content, "![embed]("+destination) ||
		strings.Contains(content, "[!embed]("+destination) ||
		strings.Contains(content, "![["+destination) {
		return true
	}
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if line == destination || line == "<"+destination+">" {
			return true
		}
	}
	return false
}

// embedLinkPostDestination reuses the existing provider/oEmbed/OG resolution,
// fallback card, metadata cache, and Lite YouTube code path. The authored notes
// remain untouched and follow the preview instead of being overwritten.
func (p *EmbedsPlugin) embedLinkPostDestination(post *models.Post, content string) string {
	destination := linkPostDestination(post)
	if destination == "" || hasAuthoredLinkPostEmbed(post.Content, destination) {
		return content
	}
	resolvedURL, sourceURL := resolveHackerNewsReference(destination)
	parsed, err := url.Parse(resolvedURL)
	if err != nil || (parsed.Scheme != schemeHTTP && parsed.Scheme != schemeHTTPS) || parsed.Host == "" {
		return content
	}
	metadata := p.fetchExternalMetadata(destination)
	card := p.buildExternalEmbedCard(resolvedURL, parsed, metadata, EmbedOptions{}, sourceURL)
	if strings.TrimSpace(card) == "" {
		return content
	}
	if strings.TrimSpace(content) == "" {
		return card
	}
	return card + "\n\n" + content
}
