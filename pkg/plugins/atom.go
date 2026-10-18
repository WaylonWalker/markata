// Package plugins provides lifecycle plugins for markata-go.
package plugins

import (
	"encoding/xml"
	"fmt"
	"time"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

// AtomFeed represents an Atom feed.
type AtomFeed struct {
	XMLName   xml.Name       `xml:"feed"`
	XMLNS     string         `xml:"xmlns,attr"`
	XMLLang   string         `xml:"xml:lang,attr,omitempty"`
	XMLNSFH   string         `xml:"xmlns:fh,attr,omitempty"`
	Title     string         `xml:"title"`
	ID        string         `xml:"id"`
	Updated   string         `xml:"updated,omitempty"`
	Subtitle  string         `xml:"subtitle,omitempty"`
	Links     []AtomFeedLink `xml:"link"`
	Author    *AtomAuthor    `xml:"author,omitempty"`
	Generator *AtomGenerator `xml:"generator,omitempty"`
	Complete  *AtomComplete  `xml:"fh:complete,omitempty"`
	Entries   []AtomEntry    `xml:"entry"`
}

type AtomGenerator struct {
	URI     string `xml:"uri,attr,omitempty"`
	Version string `xml:"version,attr,omitempty"`
	Value   string `xml:",chardata"`
}

type AtomComplete struct{}

// AtomFeedLink represents a link element in an Atom feed.
type AtomFeedLink struct {
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr,omitempty"`
	Type string `xml:"type,attr,omitempty"`
}

// AtomAuthor represents an author element in an Atom feed.
type AtomAuthor struct {
	Name  string `xml:"name"`
	Email string `xml:"email,omitempty"`
	URI   string `xml:"uri,omitempty"`
}

// AtomEntry represents an entry element in an Atom feed.
type AtomEntry struct {
	Title     string         `xml:"title"`
	ID        string         `xml:"id"`
	Updated   string         `xml:"updated,omitempty"`
	Published string         `xml:"published,omitempty"`
	Links     []AtomFeedLink `xml:"link"`
	Summary   *AtomContent   `xml:"summary,omitempty"`
	Content   *AtomContent   `xml:"content,omitempty"`
	Author    *AtomAuthor    `xml:"author,omitempty"`
}

// AtomContent represents content with a type attribute.
type AtomContent struct {
	Type  string `xml:"type,attr"`
	Value string `xml:",chardata"`
}

// GenerateAtom generates an Atom feed from a lifecycle.Feed.
func GenerateAtom(feed *lifecycle.Feed, config *lifecycle.Config) (string, error) {
	meta := getSiteMetadata(config)
	siteURL := meta.URL
	feedPath := feed.Path
	if feedPath == "" {
		feedPath = DefaultFeedPath
	}
	feedURL := feedURLForFormat(siteURL, feedPath, "atom.xml")
	homeURL := feedHomePageURL(siteURL, feedPath)
	posts := filterFeedPagePosts(feed.Posts, feed.IncludePrivate)

	title := feedResolvedTitle(feed, meta)
	description := feedResolvedDescription(feed, meta)
	updatedTime := latestFeedTime(posts)

	atomFeed := AtomFeed{
		XMLNS:    "http://www.w3.org/2005/Atom",
		XMLLang:  meta.Language,
		Title:    title,
		ID:       feedURL,
		Updated:  updatedTime.Format(time.RFC3339),
		Subtitle: description,
		Links: []AtomFeedLink{
			{Href: homeURL, Rel: "alternate", Type: "text/html"},
			{Href: feedURL, Rel: "self", Type: "application/atom+xml"},
		},
		Author:    &AtomAuthor{Name: meta.Author, URI: meta.AuthorURL},
		Generator: &AtomGenerator{URI: "https://github.com/WaylonWalker/markata", Value: "markata-go"},
		Entries:   make([]AtomEntry, 0, len(posts)),
	}
	if isArchiveFeedPath(feedPath) {
		atomFeed.XMLNSFH = "http://purl.org/syndication/history/1.0"
		atomFeed.Complete = &AtomComplete{}
	}

	for _, hub := range getWebSubHubs(config) {
		atomFeed.Links = append(atomFeed.Links, AtomFeedLink{Href: hub, Rel: "hub"})
	}
	if isArchiveFeedPath(feedPath) {
		atomFeed.Links = append(atomFeed.Links, AtomFeedLink{
			Href: feedArchiveCurrentURL(siteURL, feedPath, "atom.xml"),
			Rel:  "current",
			Type: "application/atom+xml",
		})
	}

	// Add entries
	for _, post := range posts {
		entry := postToAtomEntry(post, meta, updatedTime)
		atomFeed.Entries = append(atomFeed.Entries, entry)
	}

	// Marshal to XML
	output, err := xml.MarshalIndent(atomFeed, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal Atom feed: %w", err)
	}

	// Add XSL stylesheet processing instruction for human-readable display in browsers
	xslPI := `<?xml-stylesheet href="/atom.xsl" type="text/xsl"?>` + "\n"
	return xml.Header + xslPI + string(output), nil
}

// GenerateAtomFromFeedConfig generates an Atom feed from a FeedConfig.
func GenerateAtomFromFeedConfig(fc *models.FeedConfig, config *lifecycle.Config) (string, error) {
	feed := &lifecycle.Feed{
		Name:           fc.Slug,
		Title:          fc.Title,
		Description:    fc.Description,
		Posts:          fc.Posts,
		IncludePrivate: fc.IncludesPrivate(),
		Path:           fc.Slug,
	}
	return GenerateAtom(feed, config)
}

// postToAtomEntry converts a Post to an AtomEntry.
func postToAtomEntry(post *models.Post, meta siteMetadata, fallback time.Time) AtomEntry {
	post = safeFeedPost(post)
	if post == nil {
		return AtomEntry{}
	}
	// Build permalink
	permalink := meta.URL + post.Href

	// Get title
	title := post.PlainTitle()
	if title == "" {
		title = post.Slug
	}

	// Determine dates
	updatedTime := postUpdatedTime(post, fallback)
	updated := updatedTime.Format(time.RFC3339)
	published := ""
	if post.Date != nil {
		published = post.Date.UTC().Format(time.RFC3339)
	}

	entry := AtomEntry{
		Title:     title,
		ID:        permalink,
		Updated:   updated,
		Published: published,
		Links: []AtomFeedLink{
			{Href: permalink, Rel: "alternate", Type: "text/html"},
		},
	}

	if author := firstAuthorForPost(post, meta); author != nil {
		entry.Author = &AtomAuthor{Name: author.Name}
		if author.URL != nil {
			entry.Author.URI = *author.URL
		}
		if !post.Private && author.Email != nil {
			entry.Author.Email = *author.Email
		}
	}

	// Add summary
	if post.Description != nil && *post.Description != "" {
		entry.Summary = &AtomContent{
			Type:  "text",
			Value: *post.Description,
		}
	}

	// Add content
	if post.ArticleHTML != "" {
		entry.Content = &AtomContent{
			Type:  "html",
			Value: post.ArticleHTML,
		}
	} else if !post.Private && post.Content != "" {
		entry.Content = &AtomContent{
			Type:  "text",
			Value: post.Content,
		}
	}
	return entry
}
