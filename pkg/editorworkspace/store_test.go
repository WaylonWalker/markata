package editorworkspace

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func storeFixture(t *testing.T) (store *Store, root string) {
	t.Helper()
	root = t.TempDir()
	writeStoreFile(t, root, "markata-go.toml", "[markata-go.glob]\npatterns=['*.md']\nuse_gitignore=true\n")
	writeStoreFile(t, root, "one.md", "---\ntitle: One\n---\n# First")
	writeStoreFile(t, root, "two.md", "---\ntitle: Two\n---\n# Second")
	store, err := Open(context.Background(), root)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := store.Close(); err != nil {
			t.Error(err)
		}
	})
	return store, root
}
func writeStoreFile(t *testing.T, root, path, source string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(root, path), []byte(source), 0o600); err != nil {
		t.Fatal(err)
	}
}
func waitStore(t *testing.T, condition func() bool) {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for !condition() {
		if time.Now().After(deadline) {
			t.Fatal("workspace update did not arrive")
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func TestStore_ReusesDiscoveryAndMetadata(t *testing.T) {
	store, _ := storeFixture(t)
	if _, err := store.Files(); err != nil {
		t.Fatal(err)
	}
	before := store.Stats()
	for i := 0; i < 10; i++ {
		if _, err := store.Paths(); err != nil {
			t.Fatal(err)
		}
		if _, err := store.Files(); err != nil {
			t.Fatal(err)
		}
		source, err := store.Read("one.md")
		if err != nil || !strings.Contains(source, "# First") {
			t.Fatalf("source=%q err=%v", source, err)
		}
	}
	after := store.Stats()
	if after.GlobScans != before.GlobScans || after.MetadataLoads != before.MetadataLoads {
		t.Fatalf("cached reads repeated work: before=%+v after=%+v", before, after)
	}
	if after.SourceReads-before.SourceReads != 10 {
		t.Fatal("source reads were not counted")
	}
}

func TestStore_CreateRefreshesPathsImmediately(t *testing.T) {
	store, _ := storeFixture(t)
	if err := store.Create("three.md", "---\ntitle: Three\n---\n# Third"); err != nil {
		t.Fatal(err)
	}
	files, err := store.Paths()
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, file := range files {
		found = found || file.Path == "three.md"
	}
	if len(files) != 3 || !found {
		t.Fatalf("paths after create = %+v, want three files including three.md", files)
	}
	source, err := store.Read("three.md")
	if err != nil || !strings.Contains(source, "# Third") {
		t.Fatalf("created source = %q err=%v", source, err)
	}
}

func TestStore_UpdatesOnlyChangedMetadata(t *testing.T) {
	store, root := storeFixture(t)
	if _, err := store.Files(); err != nil {
		t.Fatal(err)
	}
	store.mu.RLock()
	unchanged := store.entries["two.md"].info
	store.mu.RUnlock()
	writeStoreFile(t, root, "one.md", "---\ntitle: Updated one\n---\n# Changed")
	waitStore(t, func() bool {
		files, err := store.Files()
		return err == nil && len(files) == 2 && files[0].Title != nil && *files[0].Title == "Updated one"
	})
	store.mu.RLock()
	same := store.entries["two.md"].info == unchanged
	store.mu.RUnlock()
	if !same {
		t.Fatal("unchanged post was reparsed")
	}
	index, _ := store.Snapshot()
	if post := index.GetBySlug("one"); post == nil || post.Title != "Updated one" {
		t.Fatal("shared LSP index was not updated")
	}
	writeStoreFile(t, root, "three.md", "# Third")
	waitStore(t, func() bool { files, err := store.Paths(); return err == nil && len(files) == 3 })
	if err := os.Remove(filepath.Join(root, "one.md")); err != nil {
		t.Fatal(err)
	}
	waitStore(t, func() bool { files, err := store.Paths(); return err == nil && len(files) == 2 })
	if _, err := store.Read("one.md"); err == nil {
		t.Fatal("removed source remained readable")
	}
}

func TestStore_ConfigAndIgnoreInvalidation(t *testing.T) {
	store, root := storeFixture(t)
	revision := store.ConfigRevision()
	writeStoreFile(t, root, ".gitignore", "two.md\n")
	waitStore(t, func() bool { files, err := store.Paths(); return err == nil && len(files) == 1 })
	if _, err := store.Read("two.md"); err == nil {
		t.Fatal("ignored source remained readable")
	}
	writeStoreFile(t, root, "markata-go.toml", "[markata-go.glob]\npatterns=['two.md']\nuse_gitignore=false\n[markata-go.theme]\ntext_size='small'\n")
	waitStore(t, func() bool { return store.ConfigRevision() > revision })
	files, err := store.Paths()
	if err != nil || len(files) != 1 || files[0].Path != "two.md" {
		t.Fatalf("paths=%+v error=%v", files, err)
	}
}

func TestStore_CloseStopsWatcher(t *testing.T) {
	store, _ := storeFixture(t)
	if err := store.Close(); err != nil {
		t.Fatal(err)
	}
	if _, err := store.Paths(); !errors.Is(err, context.Canceled) {
		t.Fatalf("closed paths error=%v", err)
	}
	if _, err := store.Read("one.md"); !errors.Is(err, context.Canceled) {
		t.Fatalf("closed read error=%v", err)
	}
}

func TestStore_WatchesNestedAppearanceAndConfig(t *testing.T) {
	store, root := storeFixture(t)
	if err := os.Mkdir(filepath.Join(root, "config"), 0o700); err != nil {
		t.Fatal(err)
	}
	writeStoreFile(t, root, "config/theme.toml", "[markata-go.theme]\ntext_size='small'\n")
	writeStoreFile(t, root, "markata-go.toml", "[markata-go]\ninclude=['config/theme.toml']\n[markata-go.glob]\npatterns=['*.md']\n")
	waitStore(t, func() bool {
		cfg, err := store.Config()
		return err == nil && cfg != nil && cfg.Theme.TextSize == "small"
	})
	revision := store.ConfigRevision()
	writeStoreFile(t, root, "config/theme.toml", "[markata-go.theme]\ntext_size='large'\n")
	waitStore(t, func() bool {
		cfg, err := store.Config()
		return err == nil && store.ConfigRevision() > revision && cfg != nil && cfg.Theme.TextSize == "large"
	})
}
