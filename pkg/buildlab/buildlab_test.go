package buildlab

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/WaylonWalker/markata/pkg/config"
)

func TestManifest_DeterministicAndDiff(t *testing.T) {
	a, b := t.TempDir(), t.TempDir()
	for _, root := range []string{a, b} {
		if err := os.WriteFile(filepath.Join(root, "z.txt"), []byte("same"), 0o600); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(root, "a.txt"), []byte("one"), 0o600); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink("z.txt", filepath.Join(root, "link")); err != nil {
			t.Fatal(err)
		}
	}
	m1, err := BuildManifest(a, nil)
	if err != nil {
		t.Fatal(err)
	}
	m2, err := BuildManifest(b, nil)
	if err != nil {
		t.Fatal(err)
	}
	j1, err := m1.CanonicalJSON()
	if err != nil {
		t.Fatal(err)
	}
	j2, err := m2.CanonicalJSON()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(j1, j2) {
		t.Errorf("canonical manifests differ:\n%s\n%s", j1, j2)
	}
	d1, err := m1.Digest()
	if err != nil {
		t.Fatal(err)
	}
	d2, err := m2.Digest()
	if err != nil {
		t.Fatal(err)
	}
	if d1 != d2 {
		t.Error("digest differs")
	}
	if got := CompareManifests(m1, m2, nil); !got.Equal() {
		t.Fatalf("equal manifests differ: %+v", got)
	}
	if err := os.WriteFile(filepath.Join(b, "z.txt"), []byte("same!"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(filepath.Join(b, "a.txt")); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(b, "extra"), []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	m2, err = BuildManifest(b, nil)
	if err != nil {
		t.Fatal(err)
	}
	diff := CompareManifests(m1, m2, nil)
	if len(diff.Missing) != 1 || len(diff.Extra) != 1 || len(diff.Changed) != 1 {
		t.Fatalf("unexpected diff: %+v", diff)
	}
}

func TestManifest_ClassesPolicy(t *testing.T) {
	a, b := t.TempDir(), t.TempDir()
	for _, root := range []string{a, b} {
		if err := os.WriteFile(filepath.Join(root, "random"), []byte(root), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	ma, err := BuildManifest(a, map[string]OutputClass{"random": ClassVolatile})
	if err != nil {
		t.Fatal(err)
	}
	mb, err := BuildManifest(b, map[string]OutputClass{"random": ClassVolatile})
	if err != nil {
		t.Fatal(err)
	}
	if d := CompareManifests(ma, mb, nil); !d.Equal() {
		t.Fatalf("volatile output changed: %+v", d)
	}
	ma.Records[0].Class = ClassDeterministic
	mb.Records[0].Class = ClassDeterministic
	if d := CompareManifests(ma, mb, nil); d.Equal() {
		t.Fatal("deterministic byte change was ignored")
	}
}

func TestManifest_ExcludesBuildMetadata(t *testing.T) {
	root := t.TempDir()
	for _, path := range []string{
		"page.html",
		"assets/.markata-js_minify-cache",
		"assets/.markata-css_minify-cache",
		"assets/.markata-fontpack-cache",
		"assets/.markata-fonts.json",
		".markata-site-verification",
		"assets/.markata-theme/manifest.json",
	} {
		fullPath := filepath.Join(root, filepath.FromSlash(path))
		if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(fullPath, []byte(path), 0o600); err != nil {
			t.Fatal(err)
		}
	}

	manifest, err := BuildManifest(root, nil)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{".markata-site-verification", "assets/.markata-theme/manifest.json", "page.html"}
	if len(manifest.Records) != len(want) {
		t.Fatalf("manifest records = %+v, want paths %v", manifest.Records, want)
	}
	for i, record := range manifest.Records {
		if record.Path != want[i] {
			t.Fatalf("manifest record %d = %+v, want path %q", i, record, want[i])
		}
	}
}

func TestManifest_ExcludesFixtureBuildState(t *testing.T) {
	root := t.TempDir()
	for _, path := range []string{"content/post.md", "output/index.html", ".markata/build-cache.json"} {
		fullPath := filepath.Join(root, filepath.FromSlash(path))
		if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(fullPath, []byte(path), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	manifest, err := BuildManifest(root, nil, "output", ".markata")
	if err != nil {
		t.Fatal(err)
	}
	if len(manifest.Records) != 1 || manifest.Records[0].Path != "content/post.md" {
		t.Fatalf("fixture manifest = %+v", manifest.Records)
	}
}

func TestNewWorkspace_ExcludesBuildMetadata(t *testing.T) {
	fixture := t.TempDir()
	for _, name := range []string{
		".markata-css_minify-cache",
		".markata-fontpack-cache",
		".markata-fonts.json",
		".markata-js_minify-cache",
	} {
		path := filepath.Join(fixture, "nested", name)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte("generated"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(fixture, "source.md"), []byte("source"), 0o600); err != nil {
		t.Fatal(err)
	}

	w, err := NewWorkspace(fixture, t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := w.Remove(); err != nil {
			t.Errorf("remove workspace: %v", err)
		}
	})
	if _, err := os.Stat(filepath.Join(w.SiteDir, "source.md")); err != nil {
		t.Fatalf("source input was not copied: %v", err)
	}
	for _, name := range []string{
		".markata-css_minify-cache",
		".markata-fontpack-cache",
		".markata-fonts.json",
		".markata-js_minify-cache",
	} {
		if _, err := os.Stat(filepath.Join(w.SiteDir, "nested", name)); !os.IsNotExist(err) {
			t.Fatalf("build metadata %q was copied: %v", name, err)
		}
	}
}

func TestFixtureExclusionsKeepNestedSourceDirectories(t *testing.T) {
	root := t.TempDir()
	for _, path := range []string{
		"output/generated.html",
		"cache/generated.json",
		"content/output/source.md",
		"content/cache/source.md",
	} {
		fullPath := filepath.Join(root, filepath.FromSlash(path))
		if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(fullPath, []byte(path), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	manifest, err := BuildManifest(root, nil, "output", "cache")
	if err != nil {
		t.Fatal(err)
	}
	paths := make([]string, 0, len(manifest.Records))
	for _, record := range manifest.Records {
		paths = append(paths, record.Path)
	}
	want := []string{"content/cache/source.md", "content/output/source.md"}
	if !slices.Equal(paths, want) {
		t.Fatalf("fixture paths = %v, want %v", paths, want)
	}
}

func TestFixtureManifestExcludesBuildMetadata(t *testing.T) {
	fixture := t.TempDir()
	if err := os.WriteFile(filepath.Join(fixture, "source.md"), []byte("source"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(fixture, ".markata-fonts.json"), []byte("generated"), 0o600); err != nil {
		t.Fatal(err)
	}

	manifest, err := BuildManifest(fixture, nil, buildLabFixtureExclusions()...)
	if err != nil {
		t.Fatal(err)
	}
	if len(manifest.Records) != 1 || manifest.Records[0].Path != "source.md" {
		t.Fatalf("fixture manifest = %+v", manifest.Records)
	}
}

func TestWorkspace_IsolationConfig(t *testing.T) {
	fixture := t.TempDir()
	if err := os.WriteFile(filepath.Join(fixture, "markata-go.toml"), []byte("[markata-go]\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	w, err := NewWorkspace(fixture, t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := w.Remove(); err != nil {
			t.Errorf("remove workspace: %v", err)
		}
	})

	data, err := os.ReadFile(w.IsolationConfig)
	if err != nil {
		t.Fatal(err)
	}
	var config map[string]any
	if err := json.Unmarshal(data, &config); err != nil {
		t.Fatal(err)
	}
	markataGo, ok := config["markata-go"].(map[string]any)
	if !ok {
		t.Fatalf("isolation config = %+v", config)
	}
	wantCacheDirs := map[string]string{
		"cache_dir":          w.MarkataCache,
		"assets":             filepath.Join(w.MarkataCache, "assets"),
		"blogroll":           filepath.Join(w.MarkataCache, "blogroll"),
		"embeds":             filepath.Join(w.MarkataCache, "embeds"),
		"image_optimization": filepath.Join(w.MarkataCache, "image-optimization"),
		"mentions":           filepath.Join(w.MarkataCache, "mentions"),
		"tailwind":           filepath.Join(w.MarkataCache, "tailwind"),
		"webmentions":        filepath.Join(w.MarkataCache, "webmentions"),
	}
	for section, want := range wantCacheDirs {
		value := markataGo[section]
		if section == "cache_dir" {
			if value != want {
				t.Fatalf("isolation %s = %v, want %q", section, value, want)
			}
			continue
		}
		sectionMap, ok := value.(map[string]any)
		if !ok || sectionMap["cache_dir"] != want {
			t.Fatalf("isolation %s = %v, want cache_dir %q", section, value, want)
		}
	}
	search, ok := markataGo["search"].(map[string]any)
	if !ok {
		t.Fatalf("isolation search config = %v", markataGo["search"])
	}
	pagefind, ok := search["pagefind"].(map[string]any)
	if !ok || pagefind["cache_dir"] != filepath.Join(w.MarkataCache, "pagefind") {
		t.Fatalf("isolation pagefind config = %v", search["pagefind"])
	}
	if markataGo["cache_cleanup_async"] != false {
		t.Fatalf("isolation cache_cleanup_async = %v, want false", markataGo["cache_cleanup_async"])
	}
	if runtime.GOOS != buildLabWindowsOS {
		info, err := os.Stat(w.IsolationConfig)
		if err != nil {
			t.Fatal(err)
		}
		if info.Mode().Perm() != 0o600 {
			t.Fatalf("isolation config mode = %o, want 600", info.Mode().Perm())
		}
	}
}

func TestCopyFixtureReplacesExistingDestination(t *testing.T) {
	source := t.TempDir()
	destination := t.TempDir()
	if err := os.WriteFile(filepath.Join(source, "kept.txt"), []byte("source"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(destination, "stale.txt"), []byte("stale"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := CopyFixture(source, destination); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(destination, "stale.txt")); !os.IsNotExist(err) {
		t.Fatalf("stale destination file error = %v, want not exist", err)
	}
	data, err := os.ReadFile(filepath.Join(destination, "kept.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "source" {
		t.Fatalf("copied file = %q, want source", data)
	}
}

func TestCopyFixturePreservesRegularFileModificationTime(t *testing.T) {
	source := t.TempDir()
	destination := filepath.Join(t.TempDir(), "copy")
	path := filepath.Join(source, "source.md")
	if err := os.WriteFile(path, []byte("source"), 0o600); err != nil {
		t.Fatal(err)
	}
	want := time.Unix(123, 456)
	if err := os.Chtimes(path, want, want); err != nil {
		t.Fatal(err)
	}
	sourceInfo, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := CopyFixture(source, destination); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(filepath.Join(destination, "source.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !info.ModTime().Equal(sourceInfo.ModTime()) {
		t.Fatalf("copied modification time = %v, want source time %v", info.ModTime(), sourceInfo.ModTime())
	}
}

func TestCopyFixtureRejectsOverlappingDestination(t *testing.T) {
	source := t.TempDir()
	if err := os.WriteFile(filepath.Join(source, "source.md"), []byte("source"), 0o600); err != nil {
		t.Fatal(err)
	}
	for _, destination := range []string{filepath.Join(source, "nested"), filepath.Dir(source)} {
		err := CopyFixture(source, destination)
		if err == nil || !strings.Contains(err.Error(), "overlap") {
			t.Fatalf("destination %q was accepted: %v", destination, err)
		}
	}
	data, err := os.ReadFile(filepath.Join(source, "source.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "source" {
		t.Fatalf("source changed after rejected copies: %q", data)
	}
}

func TestCopyFixtureRejectsAliasedDestination(t *testing.T) {
	if runtime.GOOS == buildLabWindowsOS {
		t.Skip("symlink creation may require elevated permissions on Windows")
	}
	actualParent := t.TempDir()
	source := filepath.Join(actualParent, "fixture")
	if err := os.MkdirAll(source, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(source, "source.md"), []byte("source"), 0o600); err != nil {
		t.Fatal(err)
	}
	aliasParent := filepath.Join(t.TempDir(), "alias")
	if err := os.Symlink(actualParent, aliasParent); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	if err := CopyFixture(source, filepath.Join(aliasParent, "fixture")); err == nil || !strings.Contains(err.Error(), "overlap") {
		t.Fatalf("aliased destination was accepted: %v", err)
	}
	if _, err := os.Stat(filepath.Join(source, "source.md")); err != nil {
		t.Fatalf("source was removed by aliased destination: %v", err)
	}
}

func TestWorkspace_IsolationConfigOverridesBuiltInCacheSettings(t *testing.T) {
	fixture := t.TempDir()
	basePath := filepath.Join(fixture, "markata-go.json")
	base := `{
		"markata-go": {
			"assets": {"cache_dir": "/outside/assets"},
			"blogroll": {"cache_dir": "/outside/blogroll"},
    "embeds": {"cache_dir": "/outside/embeds"},
    "image_optimization": {"cache_dir": "/outside/image-optimization"},
    "mentions": {"cache_dir": "/outside/mentions"},
    "search": {"pagefind": {"cache_dir": "/outside/pagefind"}},
    "tailwind": {"cache_dir": "/outside/tailwind"},
    "webmentions": {"cache_dir": "/outside/webmentions"}
  }
}`
	if err := os.WriteFile(basePath, []byte(base), 0o600); err != nil {
		t.Fatal(err)
	}
	w, err := NewWorkspace(fixture, t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := w.Remove(); err != nil {
			t.Errorf("remove workspace: %v", err)
		}
	})

	loaded, err := config.LoadWithMerge(basePath, w.IsolationConfig)
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Blogroll.CacheDir != filepath.Join(w.MarkataCache, "blogroll") {
		t.Fatalf("blogroll cache_dir = %q, want workspace path", loaded.Blogroll.CacheDir)
	}
	if loaded.Assets.CacheDir != filepath.Join(w.MarkataCache, "assets") {
		t.Fatalf("assets cache_dir = %q, want workspace path", loaded.Assets.CacheDir)
	}
	if loaded.Mentions.CacheDir != filepath.Join(w.MarkataCache, "mentions") {
		t.Fatalf("mentions cache_dir = %q, want workspace path", loaded.Mentions.CacheDir)
	}
	if loaded.Search.Pagefind.CacheDir != filepath.Join(w.MarkataCache, "pagefind") {
		t.Fatalf("pagefind cache_dir = %q, want workspace path", loaded.Search.Pagefind.CacheDir)
	}
	for section, want := range map[string]string{
		"embeds":             filepath.Join(w.MarkataCache, "embeds"),
		"image_optimization": filepath.Join(w.MarkataCache, "image-optimization"),
		"tailwind":           filepath.Join(w.MarkataCache, "tailwind"),
		"webmentions":        filepath.Join(w.MarkataCache, "webmentions"),
	} {
		value, ok := loaded.Extra[section].(map[string]any)
		if !ok || value["cache_dir"] != want {
			t.Fatalf("%s cache config = %v, want cache_dir %q", section, loaded.Extra[section], want)
		}
	}
	if async, ok := loaded.Extra["cache_cleanup_async"].(bool); !ok || async {
		t.Fatalf("cache_cleanup_async = %v, want false", loaded.Extra["cache_cleanup_async"])
	}
}

func TestBuildCommand_IsolationConfigFollowsMergeArgs(t *testing.T) {
	w := Workspace{SiteDir: t.TempDir(), IsolationConfig: "/workspace/buildlab-isolation.json"}
	args, err := (BuildCommand{Args: []string{"build", "-m", "user.toml"}, OutputDir: "output"}).args(w, true)
	if err != nil {
		t.Fatal(err)
	}
	isolationIndex := slices.Index(args, "--merge-config")
	if isolationIndex < 0 || isolationIndex+1 >= len(args) || args[isolationIndex+1] != w.IsolationConfig {
		t.Fatalf("args = %v, missing isolation config", args)
	}
	if userIndex := slices.Index(args, "user.toml"); userIndex > isolationIndex {
		t.Fatalf("isolation config does not follow user merge config: %v", args)
	}
}

func TestWorkspace_ClearCacheClearsForcedCache(t *testing.T) {
	w := Workspace{SiteDir: t.TempDir(), MarkataCache: filepath.Join(t.TempDir(), "markata-cache"), HomeDir: filepath.Join(t.TempDir(), "home"), XDGCacheDir: filepath.Join(t.TempDir(), "xdg")}
	for _, path := range []string{
		filepath.Join(w.SiteDir, ".markata", "cache.json"),
		filepath.Join(w.SiteDir, ".markata.cache", "cache.json"),
		filepath.Join(w.SiteDir, ".markata-cache", "plugin.json"),
		filepath.Join(w.SiteDir, "cache", "plugin.json"),
		filepath.Join(w.MarkataCache, "build-cache.json"),
		filepath.Join(w.HomeDir, "keep"),
		filepath.Join(w.XDGCacheDir, "keep"),
	} {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte("state"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	if err := applyWorkspaceOperation(w, Operation{Type: OpClearCache}, BuildCommand{}); err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{filepath.Join(w.SiteDir, ".markata"), filepath.Join(w.SiteDir, ".markata.cache"), filepath.Join(w.SiteDir, ".markata-cache"), filepath.Join(w.SiteDir, "cache"), filepath.Join(w.MarkataCache, "build-cache.json")} {
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			t.Fatalf("cache path %q still exists: %v", path, err)
		}
	}
	for _, path := range []string{filepath.Join(w.HomeDir, "keep"), filepath.Join(w.XDGCacheDir, "keep")} {
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("non-build cache path %q was removed: %v", path, err)
		}
	}
	if _, err := os.Stat(w.MarkataCache); err != nil {
		t.Fatalf("forced cache directory was not recreated: %v", err)
	}
}

func TestScenario_PathSafetyAndPreconditions(t *testing.T) {
	root := t.TempDir()
	absolute := filepath.Join(filepath.VolumeName(root)+string(filepath.Separator), "absolute")
	for _, op := range []Operation{{Type: OpWriteFile, Path: "../escape"}, {Type: OpWriteFile, Path: absolute}} {
		if err := ApplyOperation(root, op); err == nil {
			t.Errorf("unsafe path accepted: %q", op.Path)
		}
	}
	if err := ApplyOperation(root, Operation{Type: OpWriteFile, Path: "x.txt", Content: "aaa"}); err != nil {
		t.Fatal(err)
	}
	if err := ApplyOperation(root, Operation{Type: OpReplaceExact, Path: "x.txt", Old: "aa", New: "b"}); err == nil || !strings.Contains(err.Error(), "precondition") {
		t.Fatalf("overlap precondition error = %v", err)
	}
	if err := (Scenario{Operations: []Operation{{Type: OpReplaceExact, Path: "x.txt", Old: "aaa", New: "bbb"}}}).Apply(root); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(root, "x.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "bbb" {
		t.Fatalf("got %q", data)
	}
}

func TestScenario_SetConfigAndDeterminism(t *testing.T) {
	r1, r2 := t.TempDir(), t.TempDir()
	s := Scenario{Operations: []Operation{{Type: OpSetConfig, Path: "config.toml", Key: "title", Value: `a"b`}, {Type: OpSetConfig, Path: "config.toml", Key: "title", Value: "updated"}}}
	if err := s.Apply(r1); err != nil {
		t.Fatal(err)
	}
	if err := s.Apply(r2); err != nil {
		t.Fatal(err)
	}
	a, err := os.ReadFile(filepath.Join(r1, "config.toml"))
	if err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(filepath.Join(r2, "config.toml"))
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(a, b) || strings.Count(string(a), "title =") != 1 {
		t.Fatalf("config not deterministic: %q / %q", a, b)
	}
}

func TestFixture_SameSeedAndStatefulSequence(t *testing.T) {
	cfg := FixtureConfig{Seed: 42, Posts: 3, Feeds: 2, Tags: 2, Wikilinks: 1, Embeds: 1, DependencyDepth: 2, Assets: 2, TemplateVariations: 2}
	a, b := t.TempDir(), t.TempDir()
	if _, err := GenerateFixture(a, cfg); err != nil {
		t.Fatal(err)
	}
	if _, err := GenerateFixture(b, cfg); err != nil {
		t.Fatal(err)
	}
	ma, err := BuildManifest(a, nil)
	if err != nil {
		t.Fatal(err)
	}
	mb, err := BuildManifest(b, nil)
	if err != nil {
		t.Fatal(err)
	}
	if d := CompareManifests(ma, mb, nil); !d.Equal() {
		t.Fatal(d)
	}
	s := GeneratedMutationScenario(cfg, 8)
	if err := s.Apply(a); err != nil {
		t.Fatal(err)
	}
	if len(s.Operations) != 17 || s.Operations[0].Type != OpBuild {
		t.Fatalf("missing build checkpoints: %v", s.Operations)
	}
}

func TestCopyFixture_Isolated(t *testing.T) {
	src, dst := t.TempDir(), t.TempDir()
	if err := os.WriteFile(filepath.Join(src, "file"), []byte("source"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("file", filepath.Join(src, "link")); err != nil {
		t.Fatal(err)
	}
	if err := CopyFixture(src, dst); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dst, "file"), []byte("changed"), 0o600); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(filepath.Join(src, "file"))
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != "source" {
		t.Fatal("source was modified")
	}
	if runtime.GOOS != buildLabWindowsOS {
		if x, e := os.Readlink(filepath.Join(dst, "link")); e != nil || x != "file" {
			t.Fatalf("link = %q, %v", x, e)
		}
	}
	if runtime.GOOS != buildLabWindowsOS {
		outside := filepath.Join(t.TempDir(), "outside")
		if err := os.Symlink(outside, filepath.Join(src, "escape")); err != nil {
			t.Fatal(err)
		}
		if err := CopyFixture(src, filepath.Join(t.TempDir(), "copy")); err == nil {
			t.Fatal("escaping fixture symlink was accepted")
		}
	}
}

func TestCopyFixture_AllowsSymlinkedParentPath(t *testing.T) {
	if runtime.GOOS == buildLabWindowsOS {
		t.Skip("symlink creation may require elevated permissions on Windows")
	}
	actualParent := t.TempDir()
	source := filepath.Join(actualParent, "fixture")
	if err := os.MkdirAll(source, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(source, "target"), []byte("source"), 0o600); err != nil {
		t.Fatal(err)
	}
	aliasParent := filepath.Join(t.TempDir(), "alias")
	if err := os.Symlink(actualParent, aliasParent); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	aliasedSource := filepath.Join(aliasParent, "fixture")
	if err := os.Symlink(filepath.Join(aliasedSource, "target"), filepath.Join(source, "link")); err != nil {
		t.Fatal(err)
	}

	destination := filepath.Join(t.TempDir(), "copy")
	if err := CopyFixture(aliasedSource, destination); err != nil {
		t.Fatal(err)
	}
	resolved, err := filepath.EvalSymlinks(filepath.Join(destination, "link"))
	if err != nil {
		t.Fatal(err)
	}
	want, err := filepath.EvalSymlinks(filepath.Join(destination, "target"))
	if err != nil {
		t.Fatal(err)
	}
	if resolved != want {
		t.Fatalf("resolved symlink = %q, want %q", resolved, want)
	}
}

func TestCopyFixture_RejectsRelativeEscapingSymlink(t *testing.T) {
	if runtime.GOOS == buildLabWindowsOS {
		t.Skip("symlink creation may require elevated permissions on Windows")
	}
	source := t.TempDir()
	outside := filepath.Join(t.TempDir(), "outside")
	target, err := filepath.Rel(source, outside)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(target, filepath.Join(source, "escape")); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	if err := CopyFixture(source, filepath.Join(t.TempDir(), "copy")); err == nil {
		t.Fatal("relative escaping fixture symlink was accepted")
	}
}

func TestCopyFixture_RejectsSymlinkRoot(t *testing.T) {
	if runtime.GOOS == buildLabWindowsOS {
		t.Skip("symlink creation may require elevated permissions on Windows")
	}
	source := filepath.Join(t.TempDir(), "source")
	if err := os.MkdirAll(source, 0o755); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(t.TempDir(), "fixture-link")
	if err := os.Symlink(source, link); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	if err := CopyFixture(link, filepath.Join(t.TempDir(), "copy")); err == nil || !strings.Contains(err.Error(), "symlink") {
		t.Fatalf("symlink fixture root was accepted: %v", err)
	}
}

func TestCopyFixture_RewritesInternalAbsoluteSymlink(t *testing.T) {
	if runtime.GOOS == buildLabWindowsOS {
		t.Skip("symlink creation may require elevated permissions on Windows")
	}
	src := t.TempDir()
	dst := filepath.Join(t.TempDir(), "copy")
	if err := os.WriteFile(filepath.Join(src, "target"), []byte("source"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(src, "target"), filepath.Join(src, "link")); err != nil {
		t.Fatal(err)
	}

	if err := CopyFixture(src, dst); err != nil {
		t.Fatal(err)
	}
	linkTarget, err := os.Readlink(filepath.Join(dst, "link"))
	if err != nil {
		t.Fatal(err)
	}
	if filepath.IsAbs(linkTarget) {
		t.Fatalf("copied symlink remained absolute: %q", linkTarget)
	}
	resolved, err := filepath.EvalSymlinks(filepath.Join(dst, "link"))
	if err != nil {
		t.Fatal(err)
	}
	want, err := filepath.EvalSymlinks(filepath.Join(dst, "target"))
	if err != nil {
		t.Fatal(err)
	}
	if resolved != want {
		t.Fatalf("resolved symlink = %q, want %q", resolved, want)
	}
}

func TestNewWorkspace_ExcludesBuildState(t *testing.T) {
	fixture := t.TempDir()
	for _, name := range []string{".markata", ".markata.cache", ".markata-cache", "cache", "output"} {
		path := filepath.Join(fixture, name, "generated")
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte("generated"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(fixture, "source.md"), []byte("source"), 0o600); err != nil {
		t.Fatal(err)
	}
	w, err := NewWorkspace(fixture, t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := w.Remove(); err != nil {
			t.Errorf("remove workspace: %v", err)
		}
	})
	if _, err := os.Stat(filepath.Join(w.SiteDir, "source.md")); err != nil {
		t.Fatalf("source input was not copied: %v", err)
	}
	for _, name := range []string{".markata", ".markata.cache", ".markata-cache", "cache", "output"} {
		if _, err := os.Stat(filepath.Join(w.SiteDir, name)); !os.IsNotExist(err) {
			t.Fatalf("build state %q was copied: %v", name, err)
		}
	}
}

func TestNewWorkspaceRejectsOverlappingParent(t *testing.T) {
	fixture := t.TempDir()
	if err := os.WriteFile(filepath.Join(fixture, "source.md"), []byte("source"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := NewWorkspace(fixture, fixture); err == nil || !strings.Contains(err.Error(), "overlap") {
		t.Fatalf("overlapping workspace parent was accepted: %v", err)
	}
	if _, err := os.Stat(filepath.Join(fixture, "source.md")); err != nil {
		t.Fatalf("fixture was changed by rejected workspace: %v", err)
	}
}

func TestRun_ExplicitWorkingDirectoryAndTimeout(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("Build Lab process groups are implemented on Linux")
	}
	root := t.TempDir()
	r := Run(context.Background(), RunConfig{Command: "sh", Args: []string{"-c", "pwd; printf out; printf err >&2"}, CWD: root, Timeout: time.Second})
	if !r.Successful() {
		t.Fatalf("run failed: %+v", r)
	}
	if !strings.HasPrefix(string(r.Stdout), root+"\n") || !strings.HasSuffix(string(r.Stdout), "out") || string(r.Stderr) != "err" {
		t.Fatalf("capture = %q/%q", r.Stdout, r.Stderr)
	}
	r = Run(context.Background(), RunConfig{Command: "sh", Args: []string{"-c", "sleep 2"}, CWD: root, Timeout: 20 * time.Millisecond})
	if !r.TimedOut {
		t.Fatalf("expected timeout: %+v", r)
	}
}

func TestRun_BoundsChildOutput(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("shell command is Linux-specific")
	}
	r := Run(context.Background(), RunConfig{
		Command:        "sh",
		Args:           []string{"-c", "printf 0123456789; printf abcdefghij >&2"},
		Timeout:        time.Second,
		MaxOutputBytes: 4,
	})
	if r.Successful() || !r.StdoutTruncated || !r.StderrTruncated {
		t.Fatalf("bounded output run = %+v", r)
	}
	if len(r.Stdout) != 4 || len(r.Stderr) != 4 {
		t.Fatalf("captured output lengths = %d/%d", len(r.Stdout), len(r.Stderr))
	}
	if r.Err == nil || !strings.Contains(r.Err.Error(), "exceeded") {
		t.Fatalf("bounded output error = %v", r.Err)
	}
}

func TestBuildCommand_OutputPathIsWorkspaceLocal(t *testing.T) {
	w := Workspace{SiteDir: filepath.Join(t.TempDir(), "site")}
	got, err := (BuildCommand{OutputDir: "output"}).outputPath(w)
	if err != nil {
		t.Fatal(err)
	}
	if want := filepath.Join(w.SiteDir, "output"); got != want {
		t.Fatalf("output path = %q, want %q", got, want)
	}

	if _, err := (BuildCommand{OutputDir: filepath.Join(t.TempDir(), "outside")}).outputPath(w); err == nil {
		t.Fatal("absolute output outside workspace was accepted")
	}
	if _, err := (BuildCommand{OutputDir: "."}).outputPath(w); err == nil {
		t.Fatal("site directory was accepted as the output directory")
	}
	if err := os.MkdirAll(filepath.Join(w.SiteDir, "links"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(t.TempDir(), filepath.Join(w.SiteDir, "links", "output")); err != nil {
		t.Fatal(err)
	}
	if _, err := (BuildCommand{OutputDir: "links/output"}).outputPath(w); err == nil {
		t.Fatal("symlink output path was accepted")
	}
}

func TestBuildCommand_OutputArgUsesLastFlag(t *testing.T) {
	args := []string{"build", "--output", "first", "--output=second", "-o", "last"}
	if got := outputArg(args); got != "last" {
		t.Fatalf("outputArg = %q, want last", got)
	}
}

func TestScenario_TouchUpdatesModificationTime(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "source.md")
	if err := os.WriteFile(path, []byte("source"), 0o600); err != nil {
		t.Fatal(err)
	}
	old := time.Unix(1, 0)
	if err := os.Chtimes(path, old, old); err != nil {
		t.Fatal(err)
	}
	if err := ApplyOperation(root, Operation{Type: OpTouch, Path: "source.md"}); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if !info.ModTime().After(old) {
		t.Fatalf("touched modification time = %v, want after %v", info.ModTime(), old)
	}
}

func TestBuildCommand_ArgsBindSiteDirAfterCallerFlags(t *testing.T) {
	w := Workspace{SiteDir: filepath.Join(t.TempDir(), "site"), IsolationConfig: filepath.Join(t.TempDir(), "isolation.json")}
	args, err := (BuildCommand{Args: []string{"build", "--site-dir", "/outside", "--output", "public"}, OutputDir: "public"}).args(w, false)
	if err != nil {
		t.Fatal(err)
	}
	lastSiteDir := -1
	for i, arg := range args {
		if arg == "--site-dir" {
			if i+1 >= len(args) {
				t.Fatalf("--site-dir has no value: %v", args)
			}
			lastSiteDir = i
		}
	}
	if lastSiteDir < 0 || args[lastSiteDir+1] != w.SiteDir {
		t.Fatalf("args = %v, final --site-dir does not bind workspace", args)
	}
}

func TestWorkspace_ClearOutputUsesCommandOutputDir(t *testing.T) {
	root := t.TempDir()
	output := filepath.Join(root, "public")
	keep := filepath.Join(root, "output", "keep")
	if err := os.MkdirAll(output, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Dir(keep), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(output, "stale"), []byte("stale"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(keep, []byte("keep"), 0o600); err != nil {
		t.Fatal(err)
	}

	w := Workspace{SiteDir: root}
	if err := applyWorkspaceOperation(w, Operation{Type: OpClearOutput}, BuildCommand{OutputDir: "public"}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(output); !os.IsNotExist(err) {
		t.Fatalf("configured output still exists: %v", err)
	}
	if _, err := os.Stat(keep); err != nil {
		t.Fatalf("default output was unexpectedly cleared: %v", err)
	}
}

func TestApplyOperation_AcceptsSpecCleanCacheName(t *testing.T) {
	root := t.TempDir()
	cache := filepath.Join(root, ".markata")
	if err := os.MkdirAll(cache, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(cache, "state"), []byte("state"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := ApplyOperation(root, Operation{Type: OperationType("clean-cache")}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(cache); !os.IsNotExist(err) {
		t.Fatalf("cache still exists: %v", err)
	}
}

func TestWorkspace_EnvironmentPinsIsolationAndOmitsSecrets(t *testing.T) {
	t.Setenv("BUILD_LAB_SECRET_TOKEN", "do-not-copy")
	t.Setenv("MARKATA_GO_SITE_DIR", "/unsafe/site")
	t.Setenv("MARKATA_GO_BUNDLED_ASSETS_CACHE_DIR", "/unsafe/bundled-assets")
	t.Setenv("MARKATA_GO_CACHE_DIR", "/unsafe/build-cache")
	t.Setenv("DATABASE_URL", "postgres://user:password@example.invalid/db")
	t.Setenv("GH_PAT", "do-not-copy")
	t.Setenv("KUBECONFIG", "/unsafe/kubeconfig")
	t.Setenv("HTTP_PROXY", "http://user:password@example.invalid")
	t.Setenv("XDG_CONFIG_HOME", "/unsafe/config")
	w := Workspace{
		Root:            "/workspace",
		HomeDir:         "/workspace/home",
		XDGCacheDir:     "/workspace/cache",
		XDGConfigDir:    "/workspace/config",
		XDGDataDir:      "/workspace/data",
		AppDataDir:      "/workspace/appdata",
		LocalAppDataDir: "/workspace/local-appdata",
		TempDir:         "/workspace/tmp",
	}
	env, err := w.Environment([]string{"PATH=/workspace/bin", "MARKATA_GO_ENCRYPTION_ENABLED=false"}, 1)
	if err != nil {
		t.Fatal(err)
	}
	values := make(map[string]string)
	for _, entry := range env {
		key, value, ok := strings.Cut(entry, "=")
		if ok {
			values[key] = value
		}
	}
	if values["HOME"] != w.HomeDir || values["XDG_CACHE_HOME"] != w.XDGCacheDir || values["XDG_CONFIG_HOME"] != w.XDGConfigDir || values["XDG_DATA_HOME"] != w.XDGDataDir || values["TMPDIR"] != w.TempDir || values["SOURCE_DATE_EPOCH"] != "0" || values["MARKATA_GO_DISABLE_DOTENV"] != "1" || values["PATH"] != "/workspace/bin" {
		t.Fatal("workspace environment was not pinned")
	}
	for _, key := range []string{"BUILD_LAB_SECRET_TOKEN", "CUSTOM", "DATABASE_URL", "GH_PAT", "HTTP_PROXY", "KUBECONFIG", "MARKATA_GO_CACHE_DIR", "MARKATA_GO_BUNDLED_ASSETS_CACHE_DIR", "MARKATA_GO_SITE_DIR"} {
		if values[key] != "" {
			t.Fatalf("untrusted environment key %q was copied", key)
		}
	}
	if values["MARKATA_GO_ENCRYPTION_ENABLED"] != "false" {
		t.Fatal("sensitive or isolated environment override was copied")
	}
}

func TestWorkspace_EnvironmentRejectsPathsOutsideRoot(t *testing.T) {
	root := t.TempDir()
	outside := t.TempDir()
	w := Workspace{
		Root:         root,
		HomeDir:      filepath.Join(root, "home"),
		XDGCacheDir:  filepath.Join(root, "cache"),
		XDGConfigDir: outside,
		XDGDataDir:   filepath.Join(root, "data"),
		TempDir:      filepath.Join(root, "tmp"),
	}
	if _, err := w.Environment(nil, 1); err == nil || !strings.Contains(err.Error(), "outside workspace root") {
		t.Fatalf("workspace path outside root was accepted: %v", err)
	}
}

func TestWorkspace_EnvironmentRejectsSymlinkedPaths(t *testing.T) {
	root := t.TempDir()
	outside := t.TempDir()
	link := filepath.Join(root, "home")
	if err := os.Symlink(outside, link); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	w := Workspace{
		Root:        root,
		HomeDir:     link,
		XDGCacheDir: filepath.Join(root, "cache"),
		XDGDataDir:  filepath.Join(root, "data"),
		TempDir:     filepath.Join(root, "tmp"),
	}
	if _, err := w.Environment(nil, 1); err == nil || !strings.Contains(err.Error(), "symlink") {
		t.Fatalf("symlinked workspace path was accepted: %v", err)
	}
}

func TestValidateEnvironmentRejectsSecretLookingKeys(t *testing.T) {
	for _, entry := range []string{"PATH=/tools", "MARKATA_GO_ENCRYPTION_ENABLED=false", "MARKATA_GO_OFFLINE=true", "PATH=/tools=with-equals"} {
		if err := ValidateEnvironment([]string{entry}); err != nil {
			t.Fatalf("allowed environment entry %q rejected: %v", entry, err)
		}
	}
	for _, entry := range []string{"SAFE=value", "API_TOKEN=value", "DATABASE_URL=super-secret", "GH_PAT=super-secret", "KUBECONFIG=/unsafe", "HTTP_PROXY=http://user:password@example.invalid", "HOME=/unsafe", "BROKEN", "path=/tools"} {
		err := ValidateEnvironment([]string{entry})
		if err == nil {
			t.Fatalf("environment entry %q was accepted", entry)
		}
		if strings.Contains(err.Error(), "super-secret") || strings.Contains(err.Error(), "password") {
			t.Fatalf("environment validation exposed a value: %v", err)
		}
	}
	if err := ValidateEnvironment([]string{"PATH=/one", "PATH=/two"}); err == nil {
		t.Fatal("duplicate environment key was accepted")
	}
}

func TestBuildCommandRunRejectsDisallowedEnvironmentBeforeStarting(t *testing.T) {
	run := (BuildCommand{Binary: "/path/that/must/not/be/started", Env: []string{"DATABASE_URL=super-secret"}}).run(context.Background(), Workspace{}, false, 1)
	if run.Err == nil || !strings.Contains(run.Err.Error(), "DATABASE_URL") {
		t.Fatalf("disallowed environment result = %+v", run)
	}
	if strings.Contains(run.Err.Error(), "super-secret") {
		t.Fatalf("environment value leaked in error: %v", run.Err)
	}
}

func TestRunScenarioRejectsInvalidScenarioContract(t *testing.T) {
	result, err := RunScenario(context.Background(), ScenarioRunConfig{
		Fixture:  t.TempDir(),
		Scenario: Scenario{Version: "1", Operations: []Operation{{Type: OpBuild}}},
	})
	if err == nil || !strings.Contains(err.Error(), "scenario id is required") {
		t.Fatalf("invalid scenario error = %v", err)
	}
	if result.FailureClass != FailureHarness || len(result.Diagnostics) != 1 || result.Diagnostics[0].Scope != "scenario" {
		t.Fatalf("invalid scenario result = %+v", result)
	}
}

func TestScenarioValidate_RequiresOperationPayload(t *testing.T) {
	tests := []struct {
		name string
		op   Operation
		want string
	}{
		{"write path", Operation{Type: OpWriteFile}, "path is required"},
		{"replace old text", Operation{Type: OpReplaceExact, Path: "x"}, "old text is required"},
		{"rename destination", Operation{Type: OpRename, Path: "x"}, "path and dest are required"},
		{"config key", Operation{Type: OpSetConfig, Path: "config.toml"}, "key is required"},
		{"build payload", Operation{Type: OpBuild, Path: "x"}, "does not accept payload fields"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := (Scenario{ID: "payload", Version: "1", Operations: []Operation{tt.op}}).Validate()
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("validation error = %v, want %q", err, tt.want)
			}
		})
	}
}

func TestObserve_ClassifiesHarnessAndProductFailures(t *testing.T) {
	harness := observe(RunResult{Err: errors.New("start failed")}, Manifest{Records: []FileRecord{}}, nil)
	if harness.FailureClass != FailureHarness {
		t.Fatalf("harness failure class = %q", harness.FailureClass)
	}
	product := observe(RunResult{ExitCode: 17, Err: errors.New("exit status 17")}, Manifest{Records: []FileRecord{}}, nil)
	if product.FailureClass != FailureProduct {
		t.Fatalf("product failure class = %q", product.FailureClass)
	}
	manifestFailure := observe(RunResult{}, Manifest{Records: []FileRecord{}}, errors.New("manifest failed"))
	if manifestFailure.FailureClass != FailureHarness {
		t.Fatalf("manifest failure class = %q", manifestFailure.FailureClass)
	}
}

func TestManifestForRun_RequiresOutputDirectory(t *testing.T) {
	w := Workspace{SiteDir: t.TempDir()}
	if _, err := manifestForRun(w, BuildCommand{OutputDir: "missing"}, nil); err == nil {
		t.Fatal("missing output directory was accepted")
	}
}

func TestRunScenario_PrimeAndIncrementalCheckpoints(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("Build Lab process groups are implemented on Linux")
	}
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "source.md"), []byte("source\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	binary := filepath.Join(root, "build.sh")
	const script = `#!/bin/sh
out=""
clean=0
 while [ "$#" -gt 0 ]; do
   case "$1" in
     --output) out="$2"; shift 2 ;;
     --clean) clean=1; shift ;;
     *) shift ;;
   esac
 done
 if [ "$clean" -eq 1 ]; then rm -rf "$out"; fi
 mkdir -p "$out"
 printf '%s' deterministic > "$out/result.txt"
`
	//nolint:gosec // The test fixture must be executable as a build command.
	if err := os.WriteFile(binary, []byte(script), 0o700); err != nil {
		t.Fatal(err)
	}
	result, err := RunScenario(context.Background(), ScenarioRunConfig{
		Fixture: root,
		Scenario: Scenario{ID: "prime-incremental", Version: "1", Operations: []Operation{
			{Type: OpBuild}, {Type: OpBuild},
		}},
		Baseline:  BuildCommand{Binary: binary, OutputDir: "output", Timeout: time.Second},
		Candidate: BuildCommand{Binary: binary, OutputDir: "output", Timeout: time.Second},
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Verdict != "pass" || len(result.Checkpoints) != 2 {
		t.Fatalf("result = %+v", result)
	}
	if result.Checkpoints[0].Correctness.IncrementalApplicable {
		t.Fatal("priming build was marked incremental")
	}
	if !result.Checkpoints[1].Correctness.IncrementalApplicable ||
		!result.Checkpoints[1].Correctness.IncrementalEqual {
		t.Fatalf("incremental checkpoint = %+v", result.Checkpoints[1].Correctness)
	}
}

func TestRunScenario_FailedDeterminismBuildFailsVerdict(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("Build Lab process groups are implemented on Linux")
	}
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "source.md"), []byte("source\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	binary := filepath.Join(root, "build.sh")
	const script = `#!/bin/sh
counter="$(dirname "$0")/invocations"
count=0
if [ -f "$counter" ]; then count=$(cat "$counter"); fi
count=$((count + 1))
printf '%s' "$count" > "$counter"
out=""
while [ "$#" -gt 0 ]; do
  case "$1" in
    --output) out="$2"; shift 2 ;;
    --clean) shift ;;
    *) shift ;;
  esac
done
if [ "$count" -eq 4 ]; then exit 23; fi
mkdir -p "$out"
printf '%s' deterministic > "$out/result.txt"
`
	//nolint:gosec // The test fixture must be executable as a build command.
	if err := os.WriteFile(binary, []byte(script), 0o700); err != nil {
		t.Fatal(err)
	}
	result, err := RunScenario(context.Background(), ScenarioRunConfig{
		Fixture:          root,
		Scenario:         Scenario{ID: "failed-determinism", Version: "1", Operations: []Operation{{Type: OpBuild}}},
		Baseline:         BuildCommand{Binary: binary, OutputDir: "output", Timeout: time.Second},
		Candidate:        BuildCommand{Binary: binary, OutputDir: "output", Timeout: time.Second},
		CheckDeterminism: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	checkpoint := result.Checkpoints[0]
	if result.Verdict != "fail" || result.FailureClass != FailureProduct || checkpoint.Correctness.DeterministicEqual {
		t.Fatalf("failed determinism was accepted: verdict=%q correctness=%+v", result.Verdict, checkpoint.Correctness)
	}
	if len(checkpoint.Diagnostics) == 0 || checkpoint.Diagnostics[0].Class != FailureProduct {
		t.Fatalf("failed determinism diagnostics = %+v", checkpoint.Diagnostics)
	}
	if checkpoint.CandidateDeterminism == nil || checkpoint.CandidateDeterminism.Successful {
		t.Fatalf("determinism observation = %+v", checkpoint.CandidateDeterminism)
	}
}

func TestRunScenario_RecordsDeterminismAfterFailedCandidateClean(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("Build Lab process groups are implemented on Linux")
	}
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "source.md"), []byte("source\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	binary := filepath.Join(root, "build.sh")
	const script = `#!/bin/sh
counter="$(dirname "$0")/invocations"
count=0
if [ -f "$counter" ]; then count=$(cat "$counter"); fi
count=$((count + 1))
printf '%s' "$count" > "$counter"
out=""
while [ "$#" -gt 0 ]; do
  case "$1" in
    --output) out="$2"; shift 2 ;;
    --clean) shift ;;
    *) shift ;;
  esac
done
if [ "$count" -eq 3 ]; then exit 23; fi
mkdir -p "$out"
printf '%s' deterministic > "$out/result.txt"
`
	//nolint:gosec // The test fixture must be executable as a build command.
	if err := os.WriteFile(binary, []byte(script), 0o700); err != nil {
		t.Fatal(err)
	}
	result, err := RunScenario(context.Background(), ScenarioRunConfig{
		Fixture:          root,
		Scenario:         Scenario{ID: "failed-candidate-clean", Version: "1", Operations: []Operation{{Type: OpBuild}}},
		Baseline:         BuildCommand{Binary: binary, OutputDir: "output", Timeout: time.Second},
		Candidate:        BuildCommand{Binary: binary, OutputDir: "output", Timeout: time.Second},
		CheckDeterminism: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	checkpoint := result.Checkpoints[0]
	if result.Verdict != "fail" || result.FailureClass != FailureProduct || checkpoint.CandidateClean.Successful || checkpoint.Correctness.DeterministicEqual {
		t.Fatalf("failed candidate clean was accepted: verdict=%q checkpoint=%+v", result.Verdict, checkpoint)
	}
	if checkpoint.CandidateDeterminism == nil || !checkpoint.CandidateDeterminism.Successful {
		t.Fatalf("determinism observation = %+v", checkpoint.CandidateDeterminism)
	}
}
