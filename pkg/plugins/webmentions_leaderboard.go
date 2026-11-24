// Package plugins provides lifecycle plugins for markata-go.
package plugins

import (
	"sort"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

// Compile-time interface verification.
var (
	_ lifecycle.Plugin          = (*WebmentionsLeaderboardPlugin)(nil)
	_ lifecycle.TransformPlugin = (*WebmentionsLeaderboardPlugin)(nil)
	_ lifecycle.PriorityPlugin  = (*WebmentionsLeaderboardPlugin)(nil)
)

// LeaderboardEntry represents a post with its webmention counts.
type LeaderboardEntry struct {
	Post      *models.Post `json:"post"`
	Href      string       `json:"href"`
	Title     string       `json:"title"`
	Likes     int          `json:"likes"`
	Reposts   int          `json:"reposts"`
	Replies   int          `json:"replies"`
	Bookmarks int          `json:"bookmarks"`
	Mentions  int          `json:"mentions"`
	Total     int          `json:"total"`
}

// WebmentionLeaderboard holds the sorted leaderboard lists.
type WebmentionLeaderboard struct {
	TopLiked      []LeaderboardEntry `json:"top_liked"`
	TopReposted   []LeaderboardEntry `json:"top_reposted"`
	TopReplied    []LeaderboardEntry `json:"top_replied"`
	TopTotal      []LeaderboardEntry `json:"top_total"`
	TotalLikes    int                `json:"total_likes"`
	TotalReposts  int                `json:"total_reposts"`
	TotalReplies  int                `json:"total_replies"`
	TotalMentions int                `json:"total_mentions"`
}

// WebmentionsLeaderboardPlugin calculates top posts by webmention engagement.
type WebmentionsLeaderboardPlugin struct {
	maxEntries int
}

// NewWebmentionsLeaderboardPlugin creates a new WebmentionsLeaderboardPlugin.
func NewWebmentionsLeaderboardPlugin() *WebmentionsLeaderboardPlugin {
	return &WebmentionsLeaderboardPlugin{
		maxEntries: 20, // Default to top 20
	}
}

// Name returns the unique name of the plugin.
func (p *WebmentionsLeaderboardPlugin) Name() string {
	return "webmentions_leaderboard"
}

// Priority returns the execution priority for this plugin.
// Run after webmentions_fetch (-200) but before jinja_md (-100).
func (p *WebmentionsLeaderboardPlugin) Priority(stage lifecycle.Stage) int {
	if stage == lifecycle.StageTransform {
		return -150 // After webmentions_fetch (-200), before jinja_md (-100)
	}
	return lifecycle.PriorityDefault
}

// Transform calculates the webmention leaderboard and stores it in config.
func (p *WebmentionsLeaderboardPlugin) Transform(m *lifecycle.Manager) error {
	posts := m.Posts()

	// Build entries for all posts with webmentions
	entries := make([]LeaderboardEntry, 0, len(posts))
	totalLikes := 0
	totalReposts := 0
	totalReplies := 0
	totalMentions := 0

	for _, post := range posts {
		if post.Extra == nil {
			continue
		}

		webmentions, ok := post.Extra["webmentions"]
		if !ok {
			continue
		}

		// Type assert to slice of ReceivedWebMention
		mentions, ok := webmentions.([]ReceivedWebMention)
		if !ok {
			continue
		}

		if len(mentions) == 0 {
			continue
		}

		entry := LeaderboardEntry{
			Post:  post,
			Href:  post.Href,
			Title: getPostTitle(post),
		}

		// Count by type
		for i := range mentions {
			switch mentions[i].WMProperty {
			case "like-of":
				entry.Likes++
				totalLikes++
			case "repost-of":
				entry.Reposts++
				totalReposts++
			case "in-reply-to":
				entry.Replies++
				totalReplies++
			case "bookmark-of":
				entry.Bookmarks++
			case "mention-of":
				entry.Mentions++
				totalMentions++
			}
		}

		entry.Total = entry.Likes + entry.Reposts + entry.Replies + entry.Bookmarks + entry.Mentions
		entries = append(entries, entry)
	}

	// Create sorted leaderboards
	leaderboard := WebmentionLeaderboard{
		TotalLikes:    totalLikes,
		TotalReposts:  totalReposts,
		TotalReplies:  totalReplies,
		TotalMentions: totalMentions,
	}

	// Top by likes
	leaderboard.TopLiked = p.sortAndLimit(entries, func(a, b LeaderboardEntry) bool {
		return a.Likes > b.Likes
	}, func(e LeaderboardEntry) bool {
		return e.Likes > 0
	})

	// Top by reposts
	leaderboard.TopReposted = p.sortAndLimit(entries, func(a, b LeaderboardEntry) bool {
		return a.Reposts > b.Reposts
	}, func(e LeaderboardEntry) bool {
		return e.Reposts > 0
	})

	// Top by replies
	leaderboard.TopReplied = p.sortAndLimit(entries, func(a, b LeaderboardEntry) bool {
		return a.Replies > b.Replies
	}, func(e LeaderboardEntry) bool {
		return e.Replies > 0
	})

	// Top by total
	leaderboard.TopTotal = p.sortAndLimit(entries, func(a, b LeaderboardEntry) bool {
		return a.Total > b.Total
	}, func(e LeaderboardEntry) bool {
		return e.Total > 0
	})

	// Store in config.Extra
	config := m.Config()
	if config.Extra == nil {
		config.Extra = make(map[string]interface{})
	}
	config.Extra["webmention_leaderboard"] = leaderboard

	return nil
}

// sortAndLimit sorts entries by the given comparison function and returns top N.
func (p *WebmentionsLeaderboardPlugin) sortAndLimit(
	entries []LeaderboardEntry,
	less func(a, b LeaderboardEntry) bool,
	filter func(e LeaderboardEntry) bool,
) []LeaderboardEntry {
	// Copy and filter
	filtered := make([]LeaderboardEntry, 0, len(entries))
	for _, e := range entries {
		if filter(e) {
			filtered = append(filtered, e)
		}
	}

	// Sort
	sort.Slice(filtered, func(i, j int) bool {
		return less(filtered[i], filtered[j])
	})

	// Limit
	if len(filtered) > p.maxEntries {
		filtered = filtered[:p.maxEntries]
	}

	return filtered
}

// getPostTitle extracts the title from a post.
func getPostTitle(post *models.Post) string {
	if post.PlainTitle() != "" {
		return post.PlainTitle()
	}
	return post.Slug
}
