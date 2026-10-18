// Package plugins provides lifecycle plugins for markata-go.
package plugins

import (
	"log"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

// AuthorsPlugin resolves author IDs in posts against the site-wide authors
// configuration. For each post that has Authors or Author frontmatter fields,
// this plugin populates the computed AuthorObjects field with fully resolved
// Author structs. It also assigns the default author to posts with no authors.
type AuthorsPlugin struct{}

// NewAuthorsPlugin creates a new AuthorsPlugin.
func NewAuthorsPlugin() *AuthorsPlugin {
	return &AuthorsPlugin{}
}

// Name returns the unique name of the plugin.
func (p *AuthorsPlugin) Name() string {
	return "authors"
}

// Priority returns the plugin priority for the given stage.
// Authors should run very early in transform, right after auto_title,
// so that resolved author data is available for other plugins (e.g., structured_data).
func (p *AuthorsPlugin) Priority(stage lifecycle.Stage) int {
	if stage == lifecycle.StageTransform {
		return lifecycle.PriorityFirst + 1
	}
	return lifecycle.PriorityDefault
}

// Transform resolves author IDs for all posts.
// For each post with author references, it looks up the Author structs from
// config.Authors.Authors and populates post.AuthorObjects.
// Posts without any author references get the default author assigned.
// Per-post role overrides from AuthorRoleOverrides take precedence over config roles.
func (p *AuthorsPlugin) Transform(m *lifecycle.Manager) error {
	modelsConfig, ok := getModelsConfig(m.Config())
	if !ok {
		return nil
	}

	authorMap := modelsConfig.Authors.Authors
	if len(authorMap) == 0 {
		return nil
	}

	// Find the default author for posts with no author specified
	defaultAuthor, defaultID := models.GetDefaultAuthor(authorMap)

	posts := m.Posts()
	resolved := 0
	defaulted := 0

	for _, post := range posts {
		if post.Skip {
			continue
		}

		authorIDs := post.GetAuthors()

		if len(authorIDs) == 0 {
			// Assign default author if available
			if defaultAuthor != nil {
				post.Authors = []string{defaultID}
				post.AuthorObjects = []models.Author{*defaultAuthor}
				defaulted++
			}
			continue
		}

		// Resolve each author ID, applying per-post role and details overrides
		objects := make([]models.Author, 0, len(authorIDs))
		for _, id := range authorIDs {
			if author, exists := authorMap[id]; exists {
				hasRoleOverride := post.AuthorRoleOverrides != nil
				hasDetailsOverride := post.AuthorDetailsOverrides != nil

				roleOverride, applyRole := "", false
				if hasRoleOverride {
					roleOverride, applyRole = post.AuthorRoleOverrides[id]
				}

				detailsOverride, applyDetails := "", false
				if hasDetailsOverride {
					detailsOverride, applyDetails = post.AuthorDetailsOverrides[id]
				}

				if applyRole || applyDetails {
					// Clone the author and apply overrides
					authorCopy := author
					if applyRole {
						authorCopy.Role = &roleOverride
						// Clear Contribution so GetRoleDisplay() uses the overridden Role
						authorCopy.Contribution = nil
					}
					if applyDetails {
						authorCopy.Details = &detailsOverride
					}
					objects = append(objects, authorCopy)
					continue
				}

				objects = append(objects, author)
			} else {
				log.Printf("[authors] Warning: unknown author ID %q in post %s", id, post.Path)
			}
		}

		if len(objects) > 0 {
			post.AuthorObjects = objects
			resolved++
		}
	}

	if resolved > 0 || defaulted > 0 {
		log.Printf("[authors] Resolved authors for %d posts, assigned default author to %d posts",
			resolved, defaulted)
	}

	return nil
}

// Ensure AuthorsPlugin implements the required interfaces.
var (
	_ lifecycle.Plugin          = (*AuthorsPlugin)(nil)
	_ lifecycle.TransformPlugin = (*AuthorsPlugin)(nil)
	_ lifecycle.PriorityPlugin  = (*AuthorsPlugin)(nil)
)
