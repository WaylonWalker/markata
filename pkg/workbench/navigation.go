package workbench

import (
	"encoding/json"
	"sort"

	"github.com/WaylonWalker/markata/pkg/filter"
	"github.com/WaylonWalker/markata/pkg/lsp"
	"github.com/WaylonWalker/markata/pkg/models"
)

type feedEntry struct {
	Slug        string   `json:"slug"`
	Title       string   `json:"title"`
	Description string   `json:"description,omitempty"`
	Paths       []string `json:"paths"`
}

func (h *Handler) workspaceFeeds() []feedEntry {
	if h.siteConfig == nil {
		return []feedEntry{}
	}
	var index *lsp.Index
	if snapshot, ok := h.workspace.(interface{ Snapshot() (*lsp.Index, uint64) }); ok {
		index, _ = snapshot.Snapshot()
	}
	if index == nil {
		if _, err := h.workspace.Files(); err != nil {
			return []feedEntry{}
		}
		if snapshot, ok := h.workspace.(interface{ Snapshot() (*lsp.Index, uint64) }); ok {
			index, _ = snapshot.Snapshot()
		}
	}
	posts := make([]*models.Post, 0)
	if index != nil {
		for _, info := range index.AllPosts() {
			var post models.Post
			data, err := json.Marshal(info.Metadata)
			if err != nil || json.Unmarshal(data, &post) != nil {
				continue
			}
			post.Path, post.Slug = info.Path, info.Slug
			post.Tags = append([]string(nil), getMetadataStrings(info.Metadata["tags"])...)
			title := info.Title
			post.Title = &title
			posts = append(posts, &post)
		}
	}
	feeds := make([]feedEntry, 0, len(h.siteConfig.Feeds))
	for i := range h.siteConfig.Feeds {
		config := &h.siteConfig.Feeds[i]
		entry := feedEntry{Slug: config.Slug, Title: config.Title, Description: config.Description, Paths: []string{}}
		matches := posts
		if config.Filter != "" {
			parsed, err := filter.Parse(config.Filter)
			if err != nil {
				continue
			}
			matches = parsed.MatchAll(posts)
		}
		for _, post := range matches {
			if post.Published && !post.Private && !post.Draft {
				entry.Paths = append(entry.Paths, post.Path)
			}
		}
		sort.Strings(entry.Paths)
		feeds = append(feeds, entry)
	}
	sort.Slice(feeds, func(i, j int) bool { return feeds[i].Title < feeds[j].Title })
	return feeds
}

func getMetadataStrings(value any) []string {
	switch values := value.(type) {
	case []string:
		return values
	case []any:
		result := make([]string, 0, len(values))
		for _, value := range values {
			if text, ok := value.(string); ok {
				result = append(result, text)
			}
		}
		return result
	default:
		return nil
	}
}
