// Package config provides configuration loading and management for markata-go.
package config

import (
	"github.com/WaylonWalker/markata/pkg/models"
)

// DefaultConfig returns a Config with sensible default values.
func DefaultConfig() *models.Config {
	return &models.Config{
		Fontpack:     "system",
		License:      models.LicenseValue{Raw: models.DefaultLicenseKey},
		OutputDir:    "output",
		TemplatesDir: "templates",
		AssetsDir:    "static",
		Hooks:        []string{"default"},
		Extra:        make(map[string]any),
		GlobConfig: models.GlobConfig{
			Patterns:     []string{"*.md", "pages/**/*.md", "posts/**/*.md"},
			UseGitignore: true,
			SlugMode:     models.SlugModeFlat,
			SlugRules:    nil,
		},
		MarkdownConfig: models.MarkdownConfig{
			Highlight: models.NewHighlightConfig(),
		},
		Feeds: []models.FeedConfig{
			{
				Slug:        "archive",
				Title:       "Archive",
				Description: "All posts",
				Filter:      "published == true",
				Sort:        "date",
				Reverse:     true,
			},
		},
		FeedDefaults:  models.NewFeedDefaults(),
		WellKnown:     models.NewWellKnownConfig(),
		WebSub:        models.NewWebSubConfig(),
		ThemeCalendar: models.NewThemeCalendarConfig(),
		ErrorPages:    models.NewErrorPagesConfig(),
		ResourceHints: models.NewResourceHintsConfig(),
		Encryption:    models.NewEncryptionConfig(),
		Images:        models.NewImagesConfig(),
	}
}
