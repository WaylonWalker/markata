package services

import (
	"github.com/WaylonWalker/markata/pkg/lifecycle"
)

// NewApp creates a new App with all services initialized.
func NewApp(manager *lifecycle.Manager) *App {
	return &App{
		Posts:   newPostService(manager),
		Feeds:   newFeedService(manager),
		Tags:    newTagService(manager),
		Build:   newBuildService(manager),
		Manager: manager,
	}
}
