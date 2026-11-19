package workbench

import (
	"log"
	"path/filepath"
	"testing"

	"github.com/WaylonWalker/markata/pkg/lsp"
	"github.com/WaylonWalker/markata/pkg/models"
)

type navigationTestWorkspace struct {
	workspaceSource
	index *lsp.Index
}

func (w navigationTestWorkspace) Snapshot() (*lsp.Index, uint64) {
	return w.index, 1
}

func TestWorkspaceFeeds_MatchesRelativePathFilters(t *testing.T) {
	root := t.TempDir()
	postPath := filepath.Join(root, "pages", "blog", "post.md")
	index := lsp.NewWorkspaceIndex(log.Default(), root, []*lsp.PostInfo{{
		Path:  postPath,
		Slug:  "post",
		Title: "Post",
		Metadata: map[string]interface{}{
			"published": true,
		},
	}})
	handler := &Handler{
		root:      root,
		workspace: navigationTestWorkspace{index: index},
		siteConfig: &models.Config{Feeds: []models.FeedConfig{{
			Slug:   "blog",
			Title:  "Blog",
			Filter: "path.startswith('pages/blog/')",
		}}},
	}

	feeds := handler.workspaceFeeds()
	if len(feeds) != 1 {
		t.Fatalf("workspaceFeeds() returned %d feeds, want 1", len(feeds))
	}
	if len(feeds[0].Paths) != 1 || feeds[0].Paths[0] != "pages/blog/post.md" {
		t.Fatalf("feed paths = %v, want [pages/blog/post.md]", feeds[0].Paths)
	}
}

func TestWorkspaceRelativePath(t *testing.T) {
	tests := []struct {
		name, root, path, want string
	}{
		{"absolute under root", "/site", "/site/pages/blog/a.md", "pages/blog/a.md"},
		{"already relative", "/site", "pages/blog/a.md", "pages/blog/a.md"},
		{"outside root keeps path", "/site", "/elsewhere/a.md", "/elsewhere/a.md"},
		{"no root", "", "pages/a.md", "pages/a.md"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := workspaceRelativePath(tt.root, tt.path); got != tt.want {
				t.Errorf("workspaceRelativePath(%q, %q) = %q, want %q", tt.root, tt.path, got, tt.want)
			}
		})
	}
}
