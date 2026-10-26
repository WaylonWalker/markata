package services

import (
	"context"
	"sort"
	"strings"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

// tagService implements TagService using lifecycle.Manager.
type tagService struct {
	manager *lifecycle.Manager
}

// newTagService creates a new TagService.
func newTagService(m *lifecycle.Manager) TagService {
	return &tagService{manager: m}
}

// List returns all tags with their post counts.
func (s *tagService) List(_ context.Context) ([]TagInfo, error) {
	posts := s.manager.Posts()
	tagCounts := make(map[string]int)

	for _, p := range posts {
		for _, tag := range p.Tags {
			tagCounts[tag]++
		}
	}

	tags := make([]TagInfo, 0, len(tagCounts))
	for name, count := range tagCounts {
		tags = append(tags, TagInfo{
			Name:  name,
			Count: count,
			Slug:  slugify(name),
		})
	}

	// Sort by count descending, then name ascending
	sort.Slice(tags, func(i, j int) bool {
		if tags[i].Count != tags[j].Count {
			return tags[i].Count > tags[j].Count
		}
		return tags[i].Name < tags[j].Name
	})

	return tags, nil
}

// GetPosts returns posts with a specific tag.
func (s *tagService) GetPosts(_ context.Context, tag string, opts ListOptions) ([]*models.Post, error) {
	posts := s.manager.Posts()

	var result []*models.Post
	for _, p := range posts {
		for _, t := range p.Tags {
			if strings.EqualFold(t, tag) {
				result = append(result, p)
				break
			}
		}
	}

	// Apply sorting
	if opts.SortBy != "" {
		sortPosts(result, opts.SortBy, opts.SortOrder)
	}

	// Apply pagination
	if opts.Offset > 0 && opts.Offset < len(result) {
		result = result[opts.Offset:]
	}
	if opts.Limit > 0 && opts.Limit < len(result) {
		result = result[:opts.Limit]
	}

	return result, nil
}

// slugify converts a string to a URL-safe slug.
// This is a convenience wrapper around models.Slugify.
func slugify(s string) string {
	return models.Slugify(s)
}
