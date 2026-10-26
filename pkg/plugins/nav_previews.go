package plugins

import (
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

// previewFeeds includes auto tag pages before their feed configs are published
// in the Collect stage, which runs after post templates render.
func previewFeeds(config *lifecycle.Config, posts []*models.Post, configured []models.FeedConfig) []models.FeedConfig {
	auto := getAutoFeedsConfig(config)
	if !auto.Tags.Enabled || !auto.Tags.Formats.HTML {
		return configured
	}
	prefix := autoFeedSlugPrefix(auto.Tags.SlugPrefix, defaultTagsPrefix)
	privateTags := getPrivateTagSlugs(config)
	seen := make(map[string]bool, len(configured))
	for i := range configured {
		seen[configured[i].Slug] = true
	}
	result := append([]models.FeedConfig{}, configured...)
	postsByTag := make(map[string][]*models.Post)
	for _, post := range posts {
		if !eligibleNavPost(post) {
			continue
		}
		seenTag := make(map[string]bool, len(post.Tags))
		for _, tag := range post.Tags {
			slug := models.Slugify(tag)
			if slug != "" && !seenTag[slug] {
				postsByTag[slug] = append(postsByTag[slug], post)
				seenTag[slug] = true
			}
		}
	}
	for _, group := range collectAutoTagGroups(posts) {
		slug := prefix + "/" + group.SlugPart
		if seen[slug] || privateTags[group.SlugPart] {
			continue
		}
		feed := models.FeedConfig{Slug: slug, Title: "Posts tagged: " + group.Display, Description: fmt.Sprintf("All posts with the tag %q", group.Display), Posts: postsByTag[group.SlugPart]}
		if len(feed.Posts) > 0 {
			result = append(result, feed)
		}
	}
	return result
}

// buildNavPreviews resolves configured local navigation links once per build.
// The result is read-only while templates render concurrently.
func buildNavPreviews(config *models.Config, posts []*models.Post, feeds []models.FeedConfig, blogroll models.BlogrollConfig, random RandomPostConfig) map[string]map[string]interface{} {
	if config == nil {
		return nil
	}
	items := config.Components.Nav.Items
	if len(items) == 0 {
		items = config.Nav
	}
	items = append(append([]models.NavItem{}, items...), config.Components.Footer.Links...)
	if len(items) == 0 {
		return nil
	}

	localPreviews := buildLocalPreviews(posts, feeds)

	previews := make(map[string]map[string]interface{}, len(items))
	for _, item := range items {
		key, ok := localNavPath(item)
		if !ok {
			continue
		}
		if preview := localPreviews[key]; preview != nil {
			previews[item.URL] = preview
		} else if special := specialNavPreview(key, posts, blogroll, random); special != nil {
			previews[item.URL] = special
		}
	}
	return previews
}

// buildLocalPreviews indexes every public feed and post for article links.
// Feeds take precedence when a feed and post publish the same route.
func buildLocalPreviews(posts []*models.Post, feeds []models.FeedConfig) map[string]map[string]interface{} {
	previews := make(map[string]map[string]interface{}, len(posts)+len(feeds))
	for _, post := range posts {
		if eligibleNavPost(post) {
			previews[cleanNavPath(post.Href)] = postNavPreview(post)
		}
	}
	for i := range feeds {
		if !feeds[i].IncludesPrivate() {
			previews[cleanNavPath("/"+feeds[i].Slug)] = feedNavPreview(&feeds[i])
		}
	}
	return previews
}

func specialNavPreview(key string, posts []*models.Post, blogroll models.BlogrollConfig, random RandomPostConfig) map[string]interface{} {
	if blogroll.Enabled {
		blogrollSlug := blogroll.BlogrollSlug
		if blogrollSlug == "" {
			blogrollSlug = defaultBlogrollSlug
		}
		readerSlug := blogroll.ReaderSlug
		if readerSlug == "" {
			readerSlug = defaultReaderSlug
		}
		blogrollPath := cleanNavPath("/" + blogrollSlug)
		readerPath := cleanNavPath("/" + readerSlug)
		if key == blogrollPath || key == readerPath {
			count := 0
			for i := range blogroll.Feeds {
				if blogroll.Feeds[i].IsActive() {
					count++
				}
			}
			if key == readerPath {
				return map[string]interface{}{
					"kind": "reader", "description": "Recent articles, videos, and podcasts from followed RSS feeds.",
					"count": count, "count_label": "sources",
				}
			}
			return map[string]interface{}{
				"kind": "followed", "description": "The blogs and feeds I follow.",
				"count": count, "count_label": "sources",
			}
		}
	}
	if random.Enabled && key == cleanNavPath("/"+normalizeRandomPostPath(random.Path)) {
		return map[string]interface{}{
			"kind": "surprise", "description": "Find a post you might have missed.",
			"count": len(eligibleRandomPostHrefs(posts, random.ExcludeTags)), "count_label": "posts in the draw",
		}
	}
	return nil
}

func eligibleNavPost(post *models.Post) bool {
	return post != nil && !post.Skip && !post.Private && post.Published && !post.Draft && post.Href != ""
}

func localNavPath(item models.NavItem) (string, bool) {
	if item.External || item.URL == "" {
		return "", false
	}
	u, err := url.Parse(item.URL)
	if err != nil || u.IsAbs() || u.Host != "" || !strings.HasPrefix(u.Path, "/") {
		return "", false
	}
	return cleanNavPath(u.Path), true
}

func cleanNavPath(raw string) string {
	cleaned := path.Clean("/" + strings.Trim(raw, "/"))
	if cleaned == "/" {
		return cleaned
	}
	return cleaned + "/"
}

// feedNavPreview intentionally excludes body-derived aggregate statistics from
// global nav state. Total words/read time can change on any ordinary body edit,
// which previously turned a single archive edit into a whole-site cache miss.
// Membership count and publication rhythm remain because they depend on visible
// feed membership/dates, which legitimately change the shared nav preview.
func feedNavPreview(feed *models.FeedConfig) map[string]interface{} {
	preview := map[string]interface{}{
		"kind":        "feed",
		"title":       feed.Title,
		"description": feed.Description,
	}
	count := 0
	for _, post := range feed.Posts {
		if post == nil || post.Private || post.Skip || !post.Published || post.Draft {
			continue
		}
		count++
	}
	preview["count"] = count
	window := computeSparklineWindow(feed.Posts, false)
	preview["sparkline"] = buildFeedSparkline(feed.Posts, window, false)
	preview["sparkline_title"] = buildFeedSparklineTitle(feed.Posts, window, false)
	if strings.HasPrefix(feed.Slug, "tags/") {
		preview["tags"] = []string{strings.TrimPrefix(feed.Slug, "tags/")}
	}
	return preview
}

// postNavPreview likewise avoids word/read-time statistics. A nav-targeted
// article still invalidates the global shell when visible title, description,
// or tags change, while an ordinary body-only edit stays incremental.
func postNavPreview(post *models.Post) map[string]interface{} {
	preview := map[string]interface{}{
		"kind": "post",
	}
	if post.Title != nil {
		preview["title"] = *post.Title
	}
	if post.Description != nil {
		preview["description"] = *post.Description
	}
	if len(post.Tags) > 0 {
		preview["tags"] = post.Tags[:min(len(post.Tags), 3)]
	}
	return preview
}
