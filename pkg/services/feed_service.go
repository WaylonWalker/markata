package services

import (
	"context"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

// feedService implements FeedService using lifecycle.Manager.
type feedService struct {
	manager *lifecycle.Manager
}

// newFeedService creates a new FeedService.
func newFeedService(m *lifecycle.Manager) FeedService {
	return &feedService{manager: m}
}

// List returns all configured feeds.
func (s *feedService) List(_ context.Context) ([]*lifecycle.Feed, error) {
	return s.manager.Feeds(), nil
}

// Get returns a single feed by name.
func (s *feedService) Get(_ context.Context, name string) (*lifecycle.Feed, error) {
	feeds := s.manager.Feeds()
	for _, f := range feeds {
		if f.Name == name {
			return f, nil
		}
	}
	return nil, nil
}

// GetPosts returns posts belonging to a feed.
func (s *feedService) GetPosts(ctx context.Context, feedName string, opts ListOptions) ([]*models.Post, error) {
	feed, err := s.Get(ctx, feedName)
	if err != nil {
		return nil, err
	}
	if feed == nil {
		return nil, nil
	}

	posts := feed.Posts

	// Apply sorting from opts
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
