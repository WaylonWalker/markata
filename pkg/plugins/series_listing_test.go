package plugins

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

func TestSeriesPluginCollectSeriesListing(t *testing.T) {
	plugin := NewSeriesPlugin()
	posts := []*models.Post{
		{Slug: "part-1", Published: true, Extra: map[string]interface{}{"series": "go tutorial"}},
		{Slug: "part-2", Published: true, Extra: map[string]interface{}{"series": "go tutorial"}},
		{Slug: "draft", Published: true, Draft: true, Extra: map[string]interface{}{"series": "go tutorial"}},
		{Slug: "private", Published: true, Private: true, Extra: map[string]interface{}{"series": "secret notes"}},
		{Slug: "another", Published: true, Extra: map[string]interface{}{"series": "another series"}},
	}

	got := plugin.collectSeriesListing(posts, seriesConfig{SlugPrefix: "series", Overrides: map[string]*seriesOverride{}})
	if len(got) != 2 {
		t.Fatalf("series count = %d, want 2", len(got))
	}
	if got[0].Title != "Another Series" || got[0].Href != "/series/another-series/" || got[0].Count != 1 {
		t.Fatalf("first series = %#v, want Another Series with one post", got[0])
	}
	if got[1].Title != "Go Tutorial" || got[1].Href != "/series/go-tutorial/" || got[1].Count != 2 {
		t.Fatalf("second series = %#v, want Go Tutorial with two posts", got[1])
	}
}

func TestSeriesPluginWriteGeneratesSeriesIndex(t *testing.T) {
	outDir := t.TempDir()
	manager := lifecycle.NewManager()
	manager.SetConfig(&lifecycle.Config{
		OutputDir: outDir,
		Extra:     map[string]interface{}{},
	})
	manager.SetPosts([]*models.Post{
		{Slug: "part-1", Published: true, Extra: map[string]interface{}{"series": "go tutorial"}},
		{Slug: "part-2", Published: true, Extra: map[string]interface{}{"series": "go tutorial"}},
	})

	if err := NewSeriesPlugin().Write(manager); err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	outputPath := filepath.Join(outDir, "series", "index.html")
	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", outputPath, err)
	}
	page := string(content)
	for _, want := range []string{"Go Tutorial", `href="/series/go-tutorial/"`, "2 posts"} {
		if !strings.Contains(page, want) {
			t.Errorf("series index missing %q", want)
		}
	}
}

func TestSeriesPluginWriteRespectsSlugPrefixAndOverride(t *testing.T) {
	outDir := t.TempDir()
	manager := lifecycle.NewManager()
	manager.SetConfig(&lifecycle.Config{
		OutputDir: outDir,
		Extra: map[string]interface{}{
			"series": map[string]interface{}{
				"slug_prefix": "collections",
				"overrides": map[string]interface{}{
					"go tutorial": map[string]interface{}{
						"title":       "Go From Scratch",
						"description": "A practical Go series.",
					},
				},
			},
		},
	})
	manager.SetPosts([]*models.Post{
		{Slug: "part-1", Published: true, Extra: map[string]interface{}{"series": "go tutorial"}},
	})

	if err := NewSeriesPlugin().Write(manager); err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	outputPath := filepath.Join(outDir, "collections", "index.html")
	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", outputPath, err)
	}
	page := string(content)
	for _, want := range []string{"Go From Scratch", "A practical Go series.", `href="/collections/go-tutorial/"`} {
		if !strings.Contains(page, want) {
			t.Errorf("series index missing %q", want)
		}
	}
}
