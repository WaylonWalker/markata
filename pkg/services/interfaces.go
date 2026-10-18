// Package services provides business logic interfaces that can be reused
// across TUI, CLI, and future web interfaces.
package services

import (
	"context"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

// PostService provides business logic for post operations.
type PostService interface {
	// List returns posts matching the given options.
	List(ctx context.Context, opts ListOptions) ([]*models.Post, error)

	// Get returns a single post by path.
	Get(ctx context.Context, path string) (*models.Post, error)

	// Search returns posts matching a text query.
	Search(ctx context.Context, query string, opts SearchOptions) ([]*models.Post, error)

	// Count returns the total number of posts matching options.
	Count(ctx context.Context, opts ListOptions) (int, error)
}

// FeedService provides business logic for feed operations.
type FeedService interface {
	// List returns all configured feeds.
	List(ctx context.Context) ([]*lifecycle.Feed, error)

	// Get returns a single feed by name.
	Get(ctx context.Context, name string) (*lifecycle.Feed, error)

	// GetPosts returns posts belonging to a feed.
	GetPosts(ctx context.Context, feedName string, opts ListOptions) ([]*models.Post, error)
}

// TagService provides business logic for tag operations.
type TagService interface {
	// List returns all tags with their post counts.
	List(ctx context.Context) ([]TagInfo, error)

	// GetPosts returns posts with a specific tag.
	GetPosts(ctx context.Context, tag string, opts ListOptions) ([]*models.Post, error)
}

// BuildService provides build orchestration.
type BuildService interface {
	// Build runs the build process.
	Build(ctx context.Context, opts BuildOptions) (*BuildResult, error)

	// LoadOnly runs only the load stage (for TUI browsing without full build).
	LoadOnly(ctx context.Context) error

	// LoadForTUI runs through Collect stage for TUI browsing.
	// This includes Transform (for stats, titles) and Collect (for feeds).
	LoadForTUI(ctx context.Context) error

	// Subscribe returns a channel for build progress events.
	Subscribe() <-chan BuildEvent
}

// App bundles all services together for easy dependency injection.
type App struct {
	Posts PostService
	Feeds FeedService
	Tags  TagService
	Build BuildService

	// Manager is the underlying lifecycle manager (for advanced access)
	Manager *lifecycle.Manager
}
