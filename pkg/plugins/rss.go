// Package plugins provides lifecycle plugins for markata-go.
package plugins

import (
	"encoding/xml"
	"fmt"
	"strings"
	"time"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

// RSS represents an RSS 2.0 feed.
type RSS struct {
	XMLName xml.Name   `xml:"rss"`
	Version string     `xml:"version,attr"`
	Atom    string     `xml:"xmlns:atom,attr,omitempty"`
	FH      string     `xml:"xmlns:fh,attr,omitempty"`
	Channel RSSChannel `xml:"channel"`
}

// RSSChannel represents the channel element in an RSS feed.
type RSSChannel struct {
	Title          string       `xml:"title"`
	Link           string       `xml:"link"`
	Description    string       `xml:"description"`
	Language       string       `xml:"language,omitempty"`
	LastBuildDate  string       `xml:"lastBuildDate,omitempty"`
	ManagingEditor string       `xml:"managingEditor,omitempty"`
	WebMaster      string       `xml:"webMaster,omitempty"`
	Copyright      string       `xml:"copyright,omitempty"`
	Generator      string       `xml:"generator,omitempty"`
	Docs           string       `xml:"docs,omitempty"`
	AtomLinks      []AtomLink   `xml:"atom:link,omitempty"`
	Complete       *RSSComplete `xml:"fh:complete,omitempty"`
	Items          []RSSItem    `xml:"item"`
}

type RSSComplete struct{}

// AtomLink represents an atom:link element for RSS feed self-reference.
type AtomLink struct {
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr"`
	Type string `xml:"type,attr"`
}

// RSSItem represents an item element in an RSS feed.
type RSSItem struct {
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	Description string   `xml:"description"`
	PubDate     string   `xml:"pubDate,omitempty"`
	GUID        RSSGUID  `xml:"guid"`
	Author      string   `xml:"author,omitempty"`
	Categories  []string `xml:"category,omitempty"`
}

// RSSGUID represents a globally unique identifier for an RSS item.
type RSSGUID struct {
	Value       string `xml:",chardata"`
	IsPermaLink bool   `xml:"isPermaLink,attr"`
}

// GenerateRSS generates an RSS 2.0 feed from a lifecycle.Feed.
func GenerateRSS(feed *lifecycle.Feed, config *lifecycle.Config) (string, error) {
	meta := getSiteMetadata(config)
	feedURL := feedURLForFormat(meta.URL, feed.Path, "rss.xml")
	homeURL := feedHomePageURL(meta.URL, feed.Path)
	title := feedResolvedTitle(feed, meta)
	description := feedResolvedDescription(feed, meta)
	posts := filterFeedPagePosts(feed.Posts, feed.IncludePrivate)

	rss := RSS{
		Version: "2.0",
		Atom:    "http://www.w3.org/2005/Atom",
		Channel: RSSChannel{
			Title:          title,
			Link:           homeURL,
			Description:    description,
			Language:       meta.Language,
			ManagingEditor: meta.ManagingEditor,
			WebMaster:      meta.WebMaster,
			Copyright:      meta.Copyright,
			Generator:      "markata-go",
			Docs:           "https://www.rssboard.org/rss-specification",
			AtomLinks:      buildRSSAtomLinks(feedURL, config),
			Items:          make([]RSSItem, 0, len(posts)),
		},
	}
	if isArchiveFeedPath(feed.Path) {
		rss.FH = "http://purl.org/syndication/history/1.0"
		rss.Channel.Complete = &RSSComplete{}
		rss.Channel.AtomLinks = append(rss.Channel.AtomLinks, AtomLink{Href: feedArchiveCurrentURL(meta.URL, feed.Path, "rss.xml"), Rel: "current", Type: "application/rss+xml"})
	}

	// Set last build date based on most recent post date (deterministic)
	latest := latestFeedTime(posts)
	if !latest.IsZero() {
		rss.Channel.LastBuildDate = latest.Format(time.RFC1123Z)
	}

	// Add items
	for _, post := range posts {
		item := postToRSSItem(post, meta)
		rss.Channel.Items = append(rss.Channel.Items, item)
	}

	// Marshal to XML
	output, err := xml.MarshalIndent(rss, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal RSS: %w", err)
	}

	// Add XSL stylesheet processing instruction for human-readable display in browsers
	xslPI := `<?xml-stylesheet href="/rss.xsl" type="text/xsl"?>` + "\n"
	return xml.Header + xslPI + string(output), nil
}

func buildRSSAtomLinks(feedURL string, config *lifecycle.Config) []AtomLink {
	links := []AtomLink{
		{
			Href: feedURL,
			Rel:  "self",
			Type: "application/rss+xml",
		},
	}

	for _, hub := range getWebSubHubs(config) {
		links = append(links, AtomLink{Href: hub, Rel: "hub"})
	}

	return links
}

// GenerateRSSFromFeedConfig generates an RSS 2.0 feed from a FeedConfig.
func GenerateRSSFromFeedConfig(fc *models.FeedConfig, config *lifecycle.Config) (string, error) {
	feed := &lifecycle.Feed{
		Name:           fc.Slug,
		Title:          fc.Title,
		Description:    fc.Description,
		Posts:          fc.Posts,
		IncludePrivate: fc.IncludesPrivate(),
		Path:           fc.Slug,
	}
	return GenerateRSS(feed, config)
}

// postToRSSItem converts a Post to an RSSItem.
// Note: Do NOT manually escape XML here - xml.MarshalIndent handles escaping automatically.
// Manually escaping would cause double-encoding (e.g., & becomes &amp; then &amp;amp;).
func postToRSSItem(post *models.Post, meta siteMetadata) RSSItem {
	post = safeFeedPost(post)
	if post == nil {
		return RSSItem{}
	}
	// Build permalink
	permalink := meta.URL + post.Href

	// Get title
	title := post.PlainTitle()
	if title == "" {
		title = post.Slug
	}

	// Get description (use post description or truncated content)
	var description string
	switch {
	case post.Private && hasEncryptedFeedContent(post):
		description = post.ArticleHTML
	case post.Description != nil:
		description = *post.Description
	case post.ArticleHTML != "":
		// Use rendered HTML as description (truncated)
		description = truncateHTML(post.ArticleHTML, 500)
	default:
		description = truncateText(post.Content, 500)
	}

	// Get publication date
	pubDate := ""
	if post.Date != nil {
		pubDate = post.Date.Format(time.RFC1123Z)
	}

	item := RSSItem{
		Title:       title,
		Link:        permalink,
		Description: description,
		PubDate:     pubDate,
		GUID: RSSGUID{
			Value:       permalink,
			IsPermaLink: true,
		},
	}
	if author := firstAuthorForPost(post, meta); author != nil && !post.Private && author.Email != nil {
		item.Author = *author.Email
	}
	if len(post.Tags) > 0 {
		item.Categories = append([]string{}, post.Tags...)
	}
	return item
}

// getSiteURL retrieves the site URL from config.
func getSiteURL(config *lifecycle.Config) string {
	if config.Extra != nil {
		if url, ok := config.Extra["url"].(string); ok {
			return strings.TrimSuffix(url, "/")
		}
	}
	return DefaultSiteURL
}

// getSiteTitle retrieves the site title from config.
func getSiteTitle(config *lifecycle.Config) string {
	if config.Extra != nil {
		if title, ok := config.Extra["title"].(string); ok {
			return title
		}
	}
	return "Blog"
}

// getSiteDescription retrieves the site description from config.
func getSiteDescription(config *lifecycle.Config) string {
	if config.Extra != nil {
		if desc, ok := config.Extra["description"].(string); ok {
			return desc
		}
	}
	return ""
}

func getSiteAuthor(config *lifecycle.Config) string {
	if config.Extra != nil {
		if author, ok := config.Extra["author"].(string); ok {
			return author
		}
	}
	return ""
}

func getSiteLanguage(config *lifecycle.Config) string {
	if config.Extra != nil {
		if language, ok := config.Extra["language"].(string); ok {
			return language
		}
	}
	return ""
}

func getSiteAuthorURL(config *lifecycle.Config) string {
	if config.Extra != nil {
		if authorURL, ok := config.Extra["author_url"].(string); ok {
			return authorURL
		}
	}
	return ""
}

func getSiteManagingEditor(config *lifecycle.Config) string {
	if config.Extra != nil {
		if managingEditor, ok := config.Extra["managing_editor"].(string); ok {
			return managingEditor
		}
	}
	return ""
}

func getSiteWebMaster(config *lifecycle.Config) string {
	if config.Extra != nil {
		if webmaster, ok := config.Extra["webmaster"].(string); ok {
			return webmaster
		}
	}
	return ""
}

func getSiteCopyright(config *lifecycle.Config) string {
	if config.Extra != nil {
		if copyright, ok := config.Extra["copyright"].(string); ok {
			return copyright
		}
	}
	return ""
}

// truncateText truncates text to a maximum length, adding ellipsis if truncated.
func truncateText(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	// Find last space before maxLen
	lastSpace := strings.LastIndex(s[:maxLen], " ")
	if lastSpace > 0 {
		return s[:lastSpace] + "..."
	}
	return s[:maxLen] + "..."
}

// truncateHTML truncates HTML content while preserving anchor tags.
// It strips all HTML tags except <a> tags to keep wikilinks and mentions clickable.
func truncateHTML(html string, maxLen int) string {
	var result strings.Builder
	var textLen int
	i := 0
	runes := []rune(html)
	n := len(runes)
	openAnchors := 0 // Track open anchor tags

	for i < n && textLen < maxLen {
		if runes[i] == '<' {
			// Find the end of the tag
			tagStart := i
			i++
			for i < n && runes[i] != '>' {
				i++
			}
			if i >= n {
				break
			}
			tagEnd := i + 1
			tag := string(runes[tagStart:tagEnd])

			// Check if it's an anchor tag
			tagLower := strings.ToLower(tag)
			if strings.HasPrefix(tagLower, "<a ") || tagLower == "<a>" {
				// Opening anchor tag - preserve it
				result.WriteString(tag)
				openAnchors++
			} else if tagLower == "</a>" {
				// Closing anchor tag - preserve it
				result.WriteString(tag)
				if openAnchors > 0 {
					openAnchors--
				}
			}
			// Skip all other tags (don't write them)
			i++
		} else {
			// Regular character - write it and count toward text length
			result.WriteRune(runes[i])
			textLen++
			i++
		}
	}

	// Close any open anchor tags
	for j := 0; j < openAnchors; j++ {
		result.WriteString("</a>")
	}

	// Add ellipsis if truncated
	if textLen >= maxLen && i < n {
		// Find a good break point (last space)
		text := result.String()
		lastSpace := strings.LastIndex(text, " ")
		if lastSpace > maxLen/2 {
			// Truncate at last space, but we need to preserve unclosed anchors
			// Extract just the text portion and find anchors
			truncated := text[:lastSpace]
			// Count anchors in truncated portion
			truncatedOpenAnchors := strings.Count(strings.ToLower(truncated), "<a ") +
				strings.Count(strings.ToLower(truncated), "<a>")
			truncatedCloseAnchors := strings.Count(strings.ToLower(truncated), "</a>")
			unclosedAnchors := truncatedOpenAnchors - truncatedCloseAnchors
			// Add ellipsis and close any open anchors
			var final strings.Builder
			final.WriteString(truncated)
			final.WriteString("...")
			for j := 0; j < unclosedAnchors; j++ {
				final.WriteString("</a>")
			}
			return final.String()
		}
		result.WriteString("...")
	}

	return result.String()
}
