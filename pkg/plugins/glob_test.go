package plugins

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/WaylonWalker/markata/pkg/buildcache"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
)

func TestGlobPlugin_FastBuildRescansMovedFiles(t *testing.T) {
	tmpDir := t.TempDir()
	cache := buildcache.New(t.TempDir())

	oldPath := filepath.Join(tmpDir, "old.md")
	newPath := filepath.Join(tmpDir, "new.md")
	if err := os.WriteFile(oldPath, []byte("# old"), 0o600); err != nil {
		t.Fatalf("WriteFile(old.md) error = %v", err)
	}

	first := lifecycle.NewManager()
	first.Config().ContentDir = tmpDir
	first.Config().GlobPatterns = []string{"**/*.md"}
	first.Config().Extra = map[string]any{"fast_mode": true}
	first.Cache().Set("build_cache", cache)

	plugin := NewGlobPlugin()
	if err := plugin.Configure(first); err != nil {
		t.Fatalf("Configure() first build error = %v", err)
	}
	if err := plugin.Glob(first); err != nil {
		t.Fatalf("Glob() first build error = %v", err)
	}

	if err := os.Rename(oldPath, newPath); err != nil {
		t.Fatalf("Rename() error = %v", err)
	}

	second := lifecycle.NewManager()
	second.Config().ContentDir = tmpDir
	second.Config().GlobPatterns = []string{"**/*.md"}
	second.Config().Extra = map[string]any{"fast_mode": true}
	second.Cache().Set("build_cache", cache)

	plugin = NewGlobPlugin()
	if err := plugin.Configure(second); err != nil {
		t.Fatalf("Configure() second build error = %v", err)
	}
	if err := plugin.Glob(second); err != nil {
		t.Fatalf("Glob() second build error = %v", err)
	}

	got := second.Files()
	want := []string{"new.md"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Files() = %v, want %v", got, want)
	}
	if len(got) != 1 || got[0] == "old.md" {
		t.Fatalf("stale glob cache reused after move: %v", got)
	}
}

func TestGlobPlugin_StoresFileModTimes(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "post.md")
	if err := os.WriteFile(path, []byte("# post"), 0o600); err != nil {
		t.Fatalf("WriteFile(post.md) error = %v", err)
	}

	m := lifecycle.NewManager()
	m.Config().ContentDir = tmpDir
	m.Config().GlobPatterns = []string{"**/*.md"}

	plugin := NewGlobPlugin()
	if err := plugin.Configure(m); err != nil {
		t.Fatalf("Configure() error = %v", err)
	}
	if err := plugin.Glob(m); err != nil {
		t.Fatalf("Glob() error = %v", err)
	}

	cached, ok := m.Cache().Get(cacheKeyGlobFileModTimes)
	if !ok {
		t.Fatal("expected glob file modtimes in manager cache")
	}

	modTimes, ok := cached.(map[string]int64)
	if !ok {
		t.Fatalf("cached modtimes type = %T, want map[string]int64", cached)
	}

	if modTimes["post.md"] == 0 {
		t.Fatalf("expected non-zero modtime for post.md, got %d", modTimes["post.md"])
	}
}

func TestGlobPlugin_SkipsIgnoredDirectoriesDuringWalk(t *testing.T) {
	tmpDir := t.TempDir()
	ignoredDir := filepath.Join(tmpDir, "output")
	if err := os.MkdirAll(ignoredDir, 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	for path, contents := range map[string]string{
		filepath.Join(tmpDir, "post.md"):          "# post",
		filepath.Join(ignoredDir, "generated.md"): "# generated",
		filepath.Join(tmpDir, ".gitignore"):       "output\n",
	} {
		if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
			t.Fatalf("WriteFile(%s) error = %v", path, err)
		}
	}

	m := lifecycle.NewManager()
	m.Config().ContentDir = tmpDir
	m.Config().GlobPatterns = []string{"**/*.md"}

	plugin := NewGlobPlugin()
	if err := plugin.Configure(m); err != nil {
		t.Fatalf("Configure() error = %v", err)
	}
	if err := plugin.Glob(m); err != nil {
		t.Fatalf("Glob() error = %v", err)
	}

	want := []string{"post.md"}
	if got := m.Files(); !reflect.DeepEqual(got, want) {
		t.Fatalf("Files() = %v, want %v", got, want)
	}
}

func TestGlobPlugin_PropagatesContentDirectoryErrors(t *testing.T) {
	m := lifecycle.NewManager()
	m.Config().ContentDir = filepath.Join(t.TempDir(), "missing")
	m.Config().GlobPatterns = []string{"**/*.md"}

	plugin := NewGlobPlugin()
	if err := plugin.Configure(m); err != nil {
		t.Fatalf("Configure() error = %v", err)
	}
	err := plugin.Glob(m)
	if err == nil {
		t.Fatal("Glob() error = nil, want missing content directory error")
	}
	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("Glob() error = %v, want wrapped os.ErrNotExist", err)
	}
	if filepath.IsAbs(err.Error()) {
		t.Fatalf("Glob() error exposes absolute path: %v", err)
	}
}
