package plugins

import (
	"testing"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

func authorsTestStrPtr(s string) *string {
	return &s
}

func TestAuthorsPlugin_Transform_DetailsOverrides(t *testing.T) {
	authorMap := map[string]models.Author{
		"waylon": {
			ID:      "waylon",
			Name:    "Waylon Walker",
			Role:    authorsTestStrPtr("author"),
			Active:  true,
			Default: true,
		},
		"codex": {
			ID:   "codex",
			Name: "Codex",
			Role: authorsTestStrPtr("contributor"),
		},
	}

	modelsConfig := &models.Config{
		Authors: models.AuthorsConfig{
			Authors: authorMap,
		},
	}

	tests := []struct {
		name            string
		post            *models.Post
		wantDetails     map[string]*string // author ID -> expected Details
		wantRoles       map[string]string  // author ID -> expected role_display
		wantAuthorCount int
	}{
		{
			name: "details override applied",
			post: &models.Post{
				Path:    "test.md",
				Authors: []string{"waylon", "codex"},
				AuthorDetailsOverrides: map[string]string{
					"waylon": "wrote the introduction",
					"codex":  "wrote the code examples",
				},
			},
			wantDetails: map[string]*string{
				"waylon": authorsTestStrPtr("wrote the introduction"),
				"codex":  authorsTestStrPtr("wrote the code examples"),
			},
			wantRoles: map[string]string{
				"waylon": "author",
				"codex":  "contributor",
			},
			wantAuthorCount: 2,
		},
		{
			name: "details and role override together",
			post: &models.Post{
				Path:    "test.md",
				Authors: []string{"codex"},
				AuthorRoleOverrides: map[string]string{
					"codex": "pair programmer",
				},
				AuthorDetailsOverrides: map[string]string{
					"codex": "wrote all the tests",
				},
			},
			wantDetails: map[string]*string{
				"codex": authorsTestStrPtr("wrote all the tests"),
			},
			wantRoles: map[string]string{
				"codex": "pair programmer",
			},
			wantAuthorCount: 1,
		},
		{
			name: "no details override - details stays nil",
			post: &models.Post{
				Path:    "test.md",
				Authors: []string{"waylon"},
			},
			wantDetails: map[string]*string{
				"waylon": nil,
			},
			wantRoles: map[string]string{
				"waylon": "author",
			},
			wantAuthorCount: 1,
		},
		{
			name: "partial details override - only one author gets details",
			post: &models.Post{
				Path:    "test.md",
				Authors: []string{"waylon", "codex"},
				AuthorDetailsOverrides: map[string]string{
					"codex": "reviewed the draft",
				},
			},
			wantDetails: map[string]*string{
				"waylon": nil,
				"codex":  authorsTestStrPtr("reviewed the draft"),
			},
			wantAuthorCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := lifecycle.NewManager()
			config := &lifecycle.Config{
				Extra: map[string]interface{}{
					"models_config": modelsConfig,
				},
			}
			m.SetConfig(config)
			m.SetPosts([]*models.Post{tt.post})

			plugin := NewAuthorsPlugin()
			if err := plugin.Transform(m); err != nil {
				t.Fatalf("Transform() error = %v", err)
			}

			if len(tt.post.AuthorObjects) != tt.wantAuthorCount {
				t.Fatalf("AuthorObjects length = %d, want %d", len(tt.post.AuthorObjects), tt.wantAuthorCount)
			}

			for _, author := range tt.post.AuthorObjects {
				// Check details
				if wantDetails, exists := tt.wantDetails[author.ID]; exists {
					if wantDetails == nil {
						if author.Details != nil {
							t.Errorf("author %q Details = %q, want nil", author.ID, *author.Details)
						}
					} else {
						if author.Details == nil {
							t.Errorf("author %q Details = nil, want %q", author.ID, *wantDetails)
						} else if *author.Details != *wantDetails {
							t.Errorf("author %q Details = %q, want %q", author.ID, *author.Details, *wantDetails)
						}
					}
				}

				// Check role display
				if wantRole, exists := tt.wantRoles[author.ID]; exists {
					got := author.GetRoleDisplay()
					if got != wantRole {
						t.Errorf("author %q GetRoleDisplay() = %q, want %q", author.ID, got, wantRole)
					}
				}
			}
		})
	}
}

func TestAuthorsPlugin_DefaultConfig(t *testing.T) {
	plugin := NewAuthorsPlugin()
	if plugin.Name() != "authors" {
		t.Errorf("Name() = %q, want %q", plugin.Name(), "authors")
	}
	if plugin.Priority(lifecycle.StageTransform) != lifecycle.PriorityFirst+1 {
		t.Errorf("Priority(StageTransform) = %d, want %d", plugin.Priority(lifecycle.StageTransform), lifecycle.PriorityFirst+1)
	}
}
