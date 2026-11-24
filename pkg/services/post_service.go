package services

import (
	"context"
	"sort"
	"strings"
	"time"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

// postService implements PostService using lifecycle.Manager.
type postService struct {
	manager *lifecycle.Manager
}

// newPostService creates a new PostService.
func newPostService(m *lifecycle.Manager) PostService {
	return &postService{manager: m}
}

// List returns posts matching the given options.
func (s *postService) List(_ context.Context, opts ListOptions) ([]*models.Post, error) {
	posts := s.manager.Posts()

	// Apply filter expression
	if opts.Filter != "" {
		filtered, err := s.manager.Filter(opts.Filter)
		if err != nil {
			return nil, err
		}
		posts = filtered
	}

	// Apply tag filter (AND logic)
	if len(opts.Tags) > 0 {
		posts = filterByTags(posts, opts.Tags)
	}

	// Apply published filter
	if opts.Published != nil {
		posts = filterByPublished(posts, *opts.Published)
	}

	// Apply draft filter
	if opts.Draft != nil {
		posts = filterByDraft(posts, *opts.Draft)
	}

	// Apply date range filter
	if opts.DateRange != nil {
		posts = filterByDateRange(posts, opts.DateRange)
	}

	// Apply sorting
	if opts.SortBy != "" {
		sortPosts(posts, opts.SortBy, opts.SortOrder)
	}

	// Apply pagination
	if opts.Offset > 0 && opts.Offset < len(posts) {
		posts = posts[opts.Offset:]
	}
	if opts.Limit > 0 && opts.Limit < len(posts) {
		posts = posts[:opts.Limit]
	}

	return posts, nil
}

// Get returns a single post by path.
func (s *postService) Get(_ context.Context, path string) (*models.Post, error) {
	posts := s.manager.Posts()
	for _, p := range posts {
		if p.Path == path {
			return p, nil
		}
	}
	return nil, nil
}

// Search returns posts matching a text query.
func (s *postService) Search(_ context.Context, query string, opts SearchOptions) ([]*models.Post, error) {
	posts := s.manager.Posts()

	var results []*models.Post
	for _, p := range posts {
		if matchesSearch(p, query, opts) {
			results = append(results, p)
		}
	}

	if opts.Limit > 0 && len(results) > opts.Limit {
		results = results[:opts.Limit]
	}

	return results, nil
}

// Count returns the total number of posts matching options.
func (s *postService) Count(ctx context.Context, opts ListOptions) (int, error) {
	posts, err := s.List(ctx, opts)
	if err != nil {
		return 0, err
	}
	return len(posts), nil
}

// Helper functions

func filterByTags(posts []*models.Post, tags []string) []*models.Post {
	var result []*models.Post
	for _, p := range posts {
		if hasAllTags(p.Tags, tags) {
			result = append(result, p)
		}
	}
	return result
}

func hasAllTags(postTags, required []string) bool {
	tagSet := make(map[string]bool)
	for _, t := range postTags {
		tagSet[strings.ToLower(t)] = true
	}
	for _, req := range required {
		if !tagSet[strings.ToLower(req)] {
			return false
		}
	}
	return true
}

func filterByPublished(posts []*models.Post, published bool) []*models.Post {
	var result []*models.Post
	for _, p := range posts {
		if p.Published == published {
			result = append(result, p)
		}
	}
	return result
}

func filterByDraft(posts []*models.Post, draft bool) []*models.Post {
	var result []*models.Post
	for _, p := range posts {
		if p.Draft == draft {
			result = append(result, p)
		}
	}
	return result
}

func filterByDateRange(posts []*models.Post, dr *DateRange) []*models.Post {
	result := make([]*models.Post, 0, len(posts))
	for _, p := range posts {
		if p.Date == nil {
			continue
		}
		if dr.Start != nil && p.Date.Before(*dr.Start) {
			continue
		}
		if dr.End != nil && p.Date.After(*dr.End) {
			continue
		}
		result = append(result, p)
	}
	return result
}

func sortPosts(posts []*models.Post, field string, order SortOrder) {
	sort.SliceStable(posts, func(i, j int) bool {
		var cmp int
		switch strings.ToLower(field) {
		case "date":
			cmp = compareDates(posts[i].Date, posts[j].Date)
		case "title":
			cmp = compareTitles(posts[i], posts[j])
		case "path":
			cmp = strings.Compare(posts[i].Path, posts[j].Path)
		case "words":
			cmp = compareWordCounts(posts[i], posts[j])
		case "reading_time":
			cmp = compareReadingTimes(posts[i], posts[j])
		case "tags":
			cmp = compareTags(posts[i], posts[j])
		default:
			return false
		}

		// Use path as tie-breaker for deterministic ordering
		if cmp == 0 && field != "path" {
			cmp = strings.Compare(posts[i].Path, posts[j].Path)
		}

		if order == SortDesc {
			return cmp > 0
		}
		return cmp < 0
	})
}

func compareDates(a, b *time.Time) int {
	if a == nil && b == nil {
		return 0
	}
	if a == nil {
		return -1
	}
	if b == nil {
		return 1
	}
	if a.Before(*b) {
		return -1
	}
	if a.After(*b) {
		return 1
	}
	return 0
}

func compareTitles(a, b *models.Post) int {
	as, bs := "", ""
	if a != nil {
		as = a.PlainTitle()
	}
	if b != nil {
		bs = b.PlainTitle()
	}
	return strings.Compare(strings.ToLower(as), strings.ToLower(bs))
}

func compareWordCounts(a, b *models.Post) int {
	wcA := getWordCount(a)
	wcB := getWordCount(b)
	if wcA < wcB {
		return -1
	}
	if wcA > wcB {
		return 1
	}
	return 0
}

func getWordCount(p *models.Post) int {
	if p == nil || p.Extra == nil {
		return 0
	}
	if wc, ok := p.Extra["word_count"].(int); ok {
		return wc
	}
	// Fallback: count words in content
	if p.Content != "" {
		return len(strings.Fields(p.Content))
	}
	return 0
}

func compareReadingTimes(a, b *models.Post) int {
	rtA := getReadingTime(a)
	rtB := getReadingTime(b)
	if rtA < rtB {
		return -1
	}
	if rtA > rtB {
		return 1
	}
	return 0
}

func getReadingTime(p *models.Post) int {
	if p == nil || p.Extra == nil {
		return 0
	}
	if rt, ok := p.Extra["reading_time"].(int); ok {
		return rt
	}
	return 0
}

func compareTags(a, b *models.Post) int {
	// Compare by first tag alphabetically, then by number of tags
	tagsA := ""
	if len(a.Tags) > 0 {
		tagsA = strings.Join(a.Tags, ", ")
	}
	tagsB := ""
	if len(b.Tags) > 0 {
		tagsB = strings.Join(b.Tags, ", ")
	}
	return strings.Compare(strings.ToLower(tagsA), strings.ToLower(tagsB))
}

func matchesSearch(p *models.Post, query string, opts SearchOptions) bool {
	fields := opts.Fields
	searchAll := len(fields) == 0

	contains := func(haystack, needle string) bool {
		if opts.CaseSensitive {
			return strings.Contains(haystack, needle)
		}
		return strings.Contains(strings.ToLower(haystack), needle)
	}

	if !opts.CaseSensitive {
		query = strings.ToLower(query)
	}

	if (searchAll || hasField(fields, "title")) && contains(p.PlainTitle(), query) {
		return true
	}
	if (searchAll || hasField(fields, "description")) && p.Description != nil && contains(*p.Description, query) {
		return true
	}
	if (searchAll || hasField(fields, "content")) && contains(p.Content, query) {
		return true
	}
	if searchAll || hasField(fields, "tags") {
		for _, tag := range p.Tags {
			if contains(tag, query) {
				return true
			}
		}
	}
	return false
}

func hasField(fields []string, name string) bool {
	for _, f := range fields {
		if strings.EqualFold(f, name) {
			return true
		}
	}
	return false
}
