package plugins

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/buildcache"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
)

func TestStaticAssets_CreateHashedCopies_UsesRegistryHash(t *testing.T) {
	plugin := NewStaticAssetsPlugin()

	outputDir := filepath.Join(t.TempDir(), "output")
	cssDir := filepath.Join(outputDir, "css")
	if err := os.MkdirAll(cssDir, 0o755); err != nil {
		t.Fatalf("create css dir: %v", err)
	}

	content := []byte("body { color: red; }")
	origPath := filepath.Join(cssDir, "main.css")
	if err := os.WriteFile(origPath, content, 0o600); err != nil {
		t.Fatalf("write original file: %v", err)
	}

	manager := lifecycle.NewManager()
	config := lifecycle.NewConfig()
	config.OutputDir = outputDir
	manager.SetConfig(config)

	actualHash := fmt.Sprintf("%x", sha256.Sum256(content))[:8]
	wantedHash := "deadbeef"
	if wantedHash == actualHash {
		wantedHash = "feedbeef"
	}
	manager.SetAssetHash("css/main.css", wantedHash)

	if err := plugin.createHashedCopies(manager, outputDir); err != nil {
		t.Fatalf("create hashed copies: %v", err)
	}

	hashedPath := filepath.Join(cssDir, fmt.Sprintf("main.%s.css", wantedHash))
	data, err := os.ReadFile(hashedPath)
	if err != nil {
		t.Fatalf("read hashed file: %v", err)
	}
	if !bytes.Equal(data, content) {
		t.Errorf("hashed content = %q, want %q", string(data), string(content))
	}

	if wantedHash != actualHash {
		unexpectedPath := filepath.Join(cssDir, fmt.Sprintf("main.%s.css", actualHash))
		if _, err := os.Stat(unexpectedPath); err == nil {
			t.Errorf("unexpected hashed file created: %s", unexpectedPath)
		} else if !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("stat unexpected hash: %v", err)
		}
	}
}

func TestStaticAssets_CopyDirSkipsSymbolicLinks(t *testing.T) {
	root := t.TempDir()
	source := filepath.Join(root, "static")
	destination := filepath.Join(root, "output")
	if err := os.MkdirAll(source, 0o755); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(root, "private.png")
	if err := os.WriteFile(target, []byte("private"), 0o600); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(source, "private.png")
	if err := os.Symlink(target, link); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}

	if err := NewStaticAssetsPlugin().copyDir(source, destination); err != nil {
		t.Fatalf("copyDir() error = %v", err)
	}
	if _, err := os.Lstat(filepath.Join(destination, "private.png")); !os.IsNotExist(err) {
		t.Fatalf("symbolic link was copied, stat error = %v", err)
	}
}

func TestStaticAssets_RejectsSymlinkedOutputRoot(t *testing.T) {
	root := t.TempDir()
	realOutput := filepath.Join(root, "real-output")
	output := filepath.Join(root, "output")
	if err := os.MkdirAll(realOutput, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(realOutput, output); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}

	manager := lifecycle.NewManager()
	config := lifecycle.NewConfig()
	config.ContentDir = filepath.Join(root, "content")
	config.OutputDir = output
	manager.SetConfig(config)

	if err := NewStaticAssetsPlugin().Write(manager); err == nil || !strings.Contains(err.Error(), "symlink") {
		t.Fatalf("Write() error = %v, want symlink-output-root error", err)
	}
}

func TestStaticAssets_RejectsSymlinkedDestinationComponent(t *testing.T) {
	root := t.TempDir()
	source := filepath.Join(root, "static")
	destination := filepath.Join(root, "output")
	outside := filepath.Join(root, "outside")
	if err := os.MkdirAll(filepath.Join(source, "css"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(outside, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(destination, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(destination, "css")); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	if err := os.WriteFile(filepath.Join(source, "css", "site.css"), []byte("body{}"), 0o600); err != nil {
		t.Fatal(err)
	}

	if err := NewStaticAssetsPlugin().copyDir(source, destination); err == nil || !strings.Contains(err.Error(), "symlink") {
		t.Fatalf("copyDir() error = %v, want symlink-destination error", err)
	}
	if _, err := os.Stat(filepath.Join(outside, "site.css")); !os.IsNotExist(err) {
		t.Fatalf("symlink target was written, stat error = %v", err)
	}
}

func TestVersionedFeedCSSIntegrity(t *testing.T) {
	rootDir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(rootDir, "css"), 0o755); err != nil {
		t.Fatal(err)
	}
	root, err := os.OpenRoot(rootDir)
	if err != nil {
		t.Fatal(err)
	}
	defer root.Close()
	if err := verifyVersionedFeedCSS(root, ""); err != nil {
		t.Fatalf("custom theme without feeds.css should be allowed: %v", err)
	}
	css := []byte(".feed--editorial .editorial-shot-grid { display: grid; }")
	hash := fmt.Sprintf("%x", sha256.Sum256(css))[:8]
	if err := os.WriteFile(filepath.Join(rootDir, "css", "feeds.css"), css, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := verifyVersionedFeedCSS(root, ""); err == nil {
		t.Fatal("unversioned default feed CSS should fail")
	}
	if err := verifyVersionedFeedCSS(root, "deadbeef"); err == nil {
		t.Fatal("incorrect content hash should fail")
	}
	if err := verifyVersionedFeedCSS(root, hash); err == nil {
		t.Fatal("missing hashed feed CSS should fail")
	}
	hashed := filepath.Join(rootDir, "css", "feeds."+hash+".css")
	if err := os.WriteFile(hashed, []byte("old stylesheet"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := verifyVersionedFeedCSS(root, hash); err == nil {
		t.Fatal("stale hashed CSS content should fail")
	}
	if err := os.WriteFile(hashed, css, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := verifyVersionedFeedCSS(root, hash); err != nil {
		t.Fatalf("correctly versioned feed CSS should pass: %v", err)
	}
}

// Generated CSS can change across binary revisions without any static-file or
// source changes. Only the final Configure registry can identify that change.
func TestStaticAssets_GeneratedAssetInvalidatesWarmFeed(t *testing.T) {
	cache := buildcache.New(t.TempDir())
	for index, hash := range []string{"old-palette", "new-palette", "new-palette"} {
		manager := lifecycle.NewManager()
		config := lifecycle.NewConfig()
		config.Extra = map[string]interface{}{"static_dir": t.TempDir()}
		manager.SetConfig(config)
		manager.Cache().Set("build_cache", cache)
		plugin := NewStaticAssetsPlugin()
		if err := plugin.Configure(manager); err != nil {
			t.Fatal(err)
		}
		// Palette Configure runs later than static_assets Configure.
		manager.SetAssetHash("css/palette-core.css", hash)
		if err := plugin.Validate(manager); err != nil {
			t.Fatal(err)
		}
		want := ""
		if index == 2 {
			want = "rendered-feed"
		}
		if got := cache.GetFeedHash("archive"); got != want {
			t.Fatalf("build %d feed cache = %q, want %q", index, got, want)
		}
		cache.SetFeedHash("archive", "rendered-feed")
	}
}
