// Package plugins provides lifecycle plugins for markata-go.
package plugins

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

// JSONFeed represents a JSON Feed 1.1 document.
// See https://jsonfeed.org/version/1.1
type JSONFeed struct {
	Version     string           `json:"version"`
	Title       string           `json:"title"`
	HomePageURL string           `json:"home_page_url,omitempty"`
	FeedURL     string           `json:"feed_url,omitempty"`
	Description string           `json:"description,omitempty"`
	UserComment string           `json:"user_comment,omitempty"`
	NextURL     string           `json:"next_url,omitempty"`
	Icon        string           `json:"icon,omitempty"`
	Favicon     string           `json:"favicon,omitempty"`
	Authors     []JSONFeedAuthor `json:"authors,omitempty"`
	Language    string           `json:"language,omitempty"`
	Expired     bool             `json:"expired,omitempty"`
	Items       []JSONFeedItem   `json:"items"`
}

// JSONFeedAuthor represents an author in a JSON Feed.
type JSONFeedAuthor struct {
	Name   string `json:"name,omitempty"`
	URL    string `json:"url,omitempty"`
	Avatar string `json:"avatar,omitempty"`
}

// JSONFeedItem represents an item in a JSON Feed.
type JSONFeedItem struct {
	ID            string           `json:"id"`
	URL           string           `json:"url,omitempty"`
	ExternalURL   string           `json:"external_url,omitempty"`
	Title         string           `json:"title,omitempty"`
	ContentHTML   string           `json:"content_html,omitempty"`
	ContentText   string           `json:"content_text,omitempty"`
	Summary       string           `json:"summary,omitempty"`
	Image         string           `json:"image,omitempty"`
	BannerImage   string           `json:"banner_image,omitempty"`
	DatePublished string           `json:"date_published,omitempty"`
	DateModified  string           `json:"date_modified,omitempty"`
	Authors       []JSONFeedAuthor `json:"authors,omitempty"`
	Tags          []string         `json:"tags,omitempty"`
	Language      string           `json:"language,omitempty"`
}

// JSONFeedVersion is the JSON Feed specification version.
const JSONFeedVersion = "https://jsonfeed.org/version/1.1"

// GenerateJSONFeed generates a JSON Feed 1.1 document from a lifecycle.Feed.
func GenerateJSONFeed(feed *lifecycle.Feed, config *lifecycle.Config) (string, error) {
	meta := getSiteMetadata(config)
	feedPath := feed.Path
	if feedPath == "" {
		feedPath = DefaultFeedPath
	}
	feedURL := feedURLForFormat(meta.URL, feedPath, "feed.json")
	homePageURL := feedHomePageURL(meta.URL, feedPath)
	title := feedResolvedTitle(feed, meta)
	description := feedResolvedDescription(feed, meta)
	posts := filterFeedPagePosts(feed.Posts, feed.IncludePrivate)

	jsonFeed := JSONFeed{
		Version:     JSONFeedVersion,
		Title:       title,
		HomePageURL: homePageURL,
		FeedURL:     feedURL,
		Description: description,
		Language:    meta.Language,
		Icon:        meta.LogoURL,
		Items:       make([]JSONFeedItem, 0, len(posts)),
	}

	// Add author if available
	if meta.Author != "" {
		jsonFeed.Authors = []JSONFeedAuthor{
			{Name: meta.Author, URL: meta.AuthorURL},
		}
	}

	// Add items
	for _, post := range posts {
		item := postToJSONFeedItem(post, meta)
		jsonFeed.Items = append(jsonFeed.Items, item)
	}

	// Marshal to JSON with indentation
	output, err := json.MarshalIndent(jsonFeed, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON feed: %w", err)
	}

	return string(output), nil
}

// GenerateJSONFeedFromFeedConfig generates a JSON Feed 1.1 from a FeedConfig.
func GenerateJSONFeedFromFeedConfig(fc *models.FeedConfig, config *lifecycle.Config) (string, error) {
	feed := &lifecycle.Feed{
		Name:           fc.Slug,
		Title:          fc.Title,
		Description:    fc.Description,
		Posts:          fc.Posts,
		IncludePrivate: fc.IncludesPrivate(),
		Path:           fc.Slug,
	}
	return GenerateJSONFeed(feed, config)
}

// postToJSONFeedItem converts a Post to a JSONFeedItem.
func postToJSONFeedItem(post *models.Post, meta siteMetadata) JSONFeedItem {
	post = safeFeedPost(post)
	if post == nil {
		return JSONFeedItem{}
	}
	// Build permalink
	permalink := meta.URL + post.Href

	item := JSONFeedItem{
		ID:  permalink,
		URL: permalink,
	}

	// Add title
	item.Title = post.PlainTitle()
	if item.Title == "" {
		item.Title = post.Slug
	}

	// Add content
	if post.ArticleHTML != "" {
		item.ContentHTML = post.ArticleHTML
	}
	if !post.Private && post.Content != "" {
		item.ContentText = post.Content
	}

	// Add summary/description
	if post.Description != nil && *post.Description != "" {
		item.Summary = *post.Description
	}

	// Add dates
	if post.Date != nil {
		dateStr := post.Date.Format(time.RFC3339)
		item.DatePublished = dateStr
		item.DateModified = dateStr
	} else {
		item.DateModified = stableFallbackTime.Format(time.RFC3339)
	}

	// Add tags
	if len(post.Tags) > 0 {
		item.Tags = post.Tags
	}

	// Add image if present in Extra
	if !post.Private && post.Extra != nil {
		if img, ok := post.Extra["image"].(string); ok {
			item.Image = img
		}
	}

	if author := firstAuthorForPost(post, meta); author != nil {
		jsonAuthor := JSONFeedAuthor{Name: author.Name}
		if author.URL != nil {
			jsonAuthor.URL = *author.URL
		}
		if author.Avatar != nil {
			jsonAuthor.Avatar = *author.Avatar
		}
		item.Authors = []JSONFeedAuthor{jsonAuthor}
	}
	if meta.Language != "" {
		item.Language = meta.Language
	}

	return item
}
