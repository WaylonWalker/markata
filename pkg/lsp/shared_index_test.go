package lsp

import (
	"io"
	"log"
	"path/filepath"
	"testing"
)

func TestSharedIndex_DraftsArePrivateAndSnapshotsRefresh(t *testing.T) {
	logger := log.New(io.Discard, "", 0)
	root := t.TempDir()
	globalURI := pathToURI(filepath.Join(root, "global.md"))
	draftURI := pathToURI(filepath.Join(root, "draft.md"))
	base := NewWorkspaceIndex(logger, root, []*PostInfo{{URI: globalURI, Path: filepath.Join(root, "global.md"), Slug: "global", Title: "Original"}})
	revision := uint64(1)
	provider := func() (*Index, uint64) { return base, revision }
	one, two := New(logger), New(logger)
	one.sharedIndex = provider
	two.sharedIndex = provider
	one.documents[draftURI] = &Document{URI: draftURI, Content: "---\nslug: private-one\ntitle: One\n---\n# One"}
	two.documents[draftURI] = &Document{URI: draftURI, Content: "---\nslug: private-two\ntitle: Two\n---\n# Two"}
	one.refreshSharedIndex()
	two.refreshSharedIndex()
	if base.GetByURI(draftURI) != nil || one.index.GetBySlug("private-two") != nil || two.index.GetBySlug("private-one") != nil {
		t.Fatal("draft leaked between sessions or into shared snapshot")
	}
	base = NewWorkspaceIndex(logger, root, []*PostInfo{{URI: globalURI, Path: filepath.Join(root, "global.md"), Slug: "global", Title: "Changed"}})
	revision++
	one.refreshSharedIndex()
	if one.index.GetBySlug("global").Title != "Changed" || one.index.GetBySlug("private-one") == nil {
		t.Fatal("snapshot update lost global changes or private draft")
	}
}
