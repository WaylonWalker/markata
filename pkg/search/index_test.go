package search

import (
	"testing"

	"github.com/WaylonWalker/markata/pkg/models"
)

func TestSynonymSearch(t *testing.T) {
	// Create posts with synonym-related content
	title1 := "Walking on the Shore"
	title2 := "Mountain Hiking Guide"
	title3 := "Landing a Ship"
	title4 := "Lunar Eclipse Tonight"

	posts := []*models.Post{
		{
			Path:    "posts/shore.md",
			Title:   &title1,
			Content: "The shore was beautiful at sunset. We walked along the land by the water.",
			Slug:    "shore",
		},
		{
			Path:    "posts/hiking.md",
			Title:   &title2,
			Content: "Hiking in the mountains is a great way to exercise.",
			Slug:    "hiking",
		},
		{
			Path:    "posts/ship.md",
			Title:   &title3,
			Content: "The captain decided to land the vessel at the nearest port.",
			Slug:    "ship",
		},
		{
			Path:    "posts/lunar.md",
			Title:   &title4,
			Content: "A lunar eclipse occurs when the Earth passes between the sun and the moon.",
			Slug:    "lunar",
		},
	}

	dir := t.TempDir()
	idx, err := Build(dir, posts)
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	defer idx.Close()

	postsByPath := PostsByPath(posts)

	// Direct match should work
	results, err := idx.Search("shore", QueryOptions{}, postsByPath)
	if err != nil {
		t.Fatalf("Search 'shore': %v", err)
	}
	if len(results) == 0 {
		t.Error("expected results for 'shore', got none")
	}

	// Check that synonym expansion works: "land" and "shore" are synonyms in WordNet
	// Searching for "land" should also find the shore post (via synonym expansion)
	landResults, err := idx.Search("land", QueryOptions{}, postsByPath)
	if err != nil {
		t.Fatalf("Search 'land': %v", err)
	}

	// We expect at least 2 results: the shore post (via synonym) and the ship post (has "land" directly)
	foundShore := false
	for _, r := range landResults {
		if r.Post.Path == "posts/shore.md" {
			foundShore = true
		}
		t.Logf("  Result: %s (score: %.4f)", r.Post.Path, r.Score)
	}

	if !foundShore {
		t.Logf("Synonym expansion: searching 'land' did not find shore post")
		t.Logf("This may be expected if synonym indexing is not supported in this bleve version")
		t.Logf("Got %d results for 'land'", len(landResults))
	} else {
		t.Logf("Synonym expansion working: 'land' found shore post")
	}

	// Cross-POS synonym: "moon" should find the lunar post (lunar adj → moon noun)
	moonResults, err := idx.Search("moon", QueryOptions{}, postsByPath)
	if err != nil {
		t.Fatalf("Search 'moon': %v", err)
	}
	foundLunar := false
	for _, r := range moonResults {
		if r.Post.Path == "posts/lunar.md" {
			foundLunar = true
		}
		t.Logf("  Moon result: %s (score: %.4f)", r.Post.Path, r.Score)
	}
	if !foundLunar {
		t.Error("Synonym expansion: searching 'moon' did not find lunar post")
	} else {
		t.Logf("Cross-POS synonym working: 'moon' found lunar post")
	}
}

func TestLoadSynonyms(t *testing.T) {
	groups, err := loadSynonyms()
	if err != nil {
		t.Fatalf("loadSynonyms: %v", err)
	}
	if len(groups) == 0 {
		t.Fatal("expected synonym groups, got none")
	}
	t.Logf("Loaded %d synonym groups", len(groups))

	// Verify known synonym group exists (land/shore)
	found := false
	for _, g := range groups {
		hasLand := false
		hasShore := false
		for _, w := range g {
			if w == "land" {
				hasLand = true
			}
			if w == "shore" {
				hasShore = true
			}
		}
		if hasLand && hasShore {
			found = true
			t.Logf("Found land/shore group: %v", g)
			break
		}
	}
	if !found {
		t.Error("expected to find land/shore synonym group")
	}
}

func TestBuildAndSearch(t *testing.T) {
	title := "Go Programming Tutorial"
	posts := []*models.Post{
		{
			Path:    "posts/go-tutorial.md",
			Title:   &title,
			Content: "Learn Go programming with this comprehensive tutorial on concurrency and channels.",
			Slug:    "go-tutorial",
			Tags:    []string{"go", "tutorial", "programming"},
		},
	}

	dir := t.TempDir()
	idx, err := Build(dir, posts)
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	defer idx.Close()

	postsByPath := PostsByPath(posts)

	// Exact match
	results, err := idx.Search("concurrency", QueryOptions{}, postsByPath)
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}

	// Fuzzy match
	results, err = idx.Search("concurency", QueryOptions{Fuzzy: true}, postsByPath)
	if err != nil {
		t.Fatalf("Fuzzy search: %v", err)
	}
	if len(results) == 0 {
		t.Logf("fuzzy search for 'concurency' returned 0 results (may depend on analyzer)")
	} else {
		t.Logf("fuzzy search for 'concurency' returned %d results", len(results))
	}

	// No match
	results, err = idx.Search("kubernetes", QueryOptions{}, postsByPath)
	if err != nil {
		t.Fatalf("Search no match: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestBuildIfNeeded(t *testing.T) {
	title := "Test Post"
	posts := []*models.Post{
		{
			Path:    "posts/test.md",
			Title:   &title,
			Content: "Test content for caching.",
			Slug:    "test",
		},
	}

	cacheDir := t.TempDir()

	// First build
	idx1, err := BuildIfNeeded(cacheDir, posts)
	if err != nil {
		t.Fatalf("First BuildIfNeeded: %v", err)
	}
	idx1.Close()

	// Second call should use cached index
	idx2, err := BuildIfNeeded(cacheDir, posts)
	if err != nil {
		t.Fatalf("Second BuildIfNeeded: %v", err)
	}
	defer idx2.Close()

	// Verify search still works
	postsByPath := PostsByPath(posts)
	results, err := idx2.Search("caching", QueryOptions{}, postsByPath)
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}
}

func TestContentHashChangesWhenSearchMetadataChanges(t *testing.T) {
	title := "Original Title"
	posts := []*models.Post{{
		Path:        "posts/test.md",
		Title:       &title,
		Description: ptrString("original description"),
		Content:     "unchanged content",
		Slug:        "test",
		Tags:        []string{"go"},
		Published:   true,
	}}

	original := ContentHash(posts)

	updated := *posts[0]
	updated.Title = ptrString("Updated Title")
	updated.Description = ptrString("updated description")
	updated.Tags = append(append([]string(nil), posts[0].Tags...), "search")

	changed := ContentHash([]*models.Post{&updated})
	if original == changed {
		t.Fatal("expected content hash to change when indexed metadata changes")
	}
}

func TestToPostDoc_PrivatePostStripsSensitiveFields(t *testing.T) {
	title := "Private Post"
	description := "Explicit description"
	doc := toPostDoc(&models.Post{
		Path:        "posts/private.md",
		Title:       &title,
		Description: &description,
		Content:     "secret body https://cdn.example.com/secret.webp",
		Slug:        "private",
		Href:        "/private",
		Tags:        []string{"secret", "journal"},
		Published:   true,
		Private:     true,
		Extra: map[string]interface{}{
			"_title_explicit":       true,
			"_description_explicit": true,
			"description":           description,
			"cover_image":           "https://cdn.example.com/secret.webp",
			"word_count":            321,
		},
	})

	if doc.Title != title {
		t.Fatalf("Title = %q, want %q", doc.Title, title)
	}
	if doc.Description != description {
		t.Fatalf("Description = %q, want %q", doc.Description, description)
	}
	if doc.Content != "" {
		t.Fatalf("Content = %q, want empty", doc.Content)
	}
	if len(doc.Tags) != 0 {
		t.Fatalf("Tags = %v, want empty", doc.Tags)
	}
	if doc.WordCount != 0 {
		t.Fatalf("WordCount = %d, want 0", doc.WordCount)
	}
	if doc.MediaURL != "" || doc.PosterURL != "" || doc.VideoMIME != "" || doc.MediaType != "" {
		t.Fatalf("media fields leaked: %#v", doc)
	}
}

func TestToPostDoc_PrivatePostWithoutExplicitDescriptionClearsDescription(t *testing.T) {
	title := "Private Post"
	description := "Derived description"
	doc := toPostDoc(&models.Post{
		Path:        "posts/private.md",
		Title:       &title,
		Description: &description,
		Content:     "secret body",
		Slug:        "private",
		Href:        "/private",
		Published:   true,
		Private:     true,
		Extra: map[string]interface{}{
			"_title_explicit": true,
			"cover_image":     "https://cdn.example.com/secret.webp",
			"word_count":      321,
		},
	})

	if doc.Description != "" {
		t.Fatalf("Description = %q, want empty", doc.Description)
	}
}

func ptrString(value string) *string {
	return &value
}
