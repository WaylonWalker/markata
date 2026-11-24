// Package plugins provides lifecycle plugins for markata-go.
package plugins

import (
	"encoding/json"
	"path/filepath"
	"strings"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

func getSyndicationConfig(config *lifecycle.Config) models.SyndicationConfig {
	if config != nil && config.Extra != nil {
		if defaults, ok := config.Extra["feed_defaults"].(models.FeedDefaults); ok {
			return defaults.Syndication
		}
		if defaults, ok := config.Extra["feeds.defaults"].(models.FeedDefaults); ok {
			return defaults.Syndication
		}
	}
	return models.NewFeedDefaults().Syndication
}

func getFeedsPageConfig(config *lifecycle.Config) models.FeedsPageConfig {
	if config != nil && config.Extra != nil {
		if feedsPage, ok := config.Extra["feeds_page"].(models.FeedsPageConfig); ok {
			return feedsPage
		}
		if raw, ok := config.Extra["feeds_page"]; ok {
			var feedsPage models.FeedsPageConfig
			if data, err := json.Marshal(raw); err == nil {
				if err := json.Unmarshal(data, &feedsPage); err == nil {
					defaults := models.NewFeedsPageConfig()
					if feedsPage.Enabled == nil {
						feedsPage.Enabled = defaults.Enabled
					}
					if feedsPage.Title == "" {
						feedsPage.Title = defaults.Title
					}
					if feedsPage.Description == "" {
						feedsPage.Description = defaults.Description
					}
					if feedsPage.Template == "" {
						feedsPage.Template = defaults.Template
					}
					if feedsPage.SlugPrefix == "" {
						feedsPage.SlugPrefix = defaults.SlugPrefix
					}
					if feedsPage.ShowPrivateFeeds == nil {
						if rawFeedsPage, ok := raw.(map[string]interface{}); ok {
							if rawShowPrivate, ok := rawFeedsPage["show_private_feeds"].([]interface{}); ok {
								feedsPage.ShowPrivateFeeds = make([]string, 0, len(rawShowPrivate))
								for _, value := range rawShowPrivate {
									if s, ok := value.(string); ok {
										feedsPage.ShowPrivateFeeds = append(feedsPage.ShowPrivateFeeds, s)
									}
								}
							}
						}
					}
					return feedsPage
				}
			}
		}
	}
	return models.NewFeedsPageConfig()
}

func isRootFeed(fc *models.FeedConfig) bool {
	return fc != nil && fc.Slug == ""
}

func isArchiveFeed(fc *models.FeedConfig) bool {
	return fc != nil && fc.Slug == defaultArchivePrefix
}

func shouldGenerateFeedArchive(fc *models.FeedConfig, syndication models.SyndicationConfig) bool {
	if fc == nil {
		return false
	}
	if isRootFeed(fc) || isArchiveFeed(fc) {
		return false
	}
	if syndication.FeedArchivesDisabled || fc.ArchiveDisabled {
		return false
	}
	return fc.Formats.RSS || fc.Formats.Atom || fc.Formats.JSON
}

func feedArchiveDir(feedDir string) string {
	return filepath.Join(feedDir, defaultArchivePrefix)
}

func feedArchiveURL(slug, fileName string) string {
	cleanSlug := strings.Trim(slug, "/")
	if cleanSlug == "" {
		return "/" + defaultArchivePrefix + "/" + fileName
	}
	return "/" + cleanSlug + "/" + defaultArchivePrefix + "/" + fileName
}

func limitFeedPosts(posts []*models.Post, maxItems int) []*models.Post {
	if maxItems <= 0 || len(posts) <= maxItems {
		return posts
	}
	return posts[:maxItems]
}

func cloneFeedConfigWithPosts(fc *models.FeedConfig, posts []*models.Post) *models.FeedConfig {
	clone := *fc
	clone.Posts = posts
	return &clone
}
