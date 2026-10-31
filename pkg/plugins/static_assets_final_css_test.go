package plugins

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/buildcache"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
)

func TestVersionedCSSFilesAreExcludedFromFurtherTransforms(t *testing.T) {
	tests := []struct {
		path    string
		skipped bool
	}{
		{"css/components.bd9a4951.css", true},
		{"css/components.52A9D553.css", true},
		{"css/components.css", false},
		{"css/components.min.css", false},
	}
	minifier := NewCSSMinifyPlugin()
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			if got := isVersionedCSSFile(tt.path); got != tt.skipped {
				t.Errorf("versioned CSS detection = %v, want %v", got, tt.skipped)
			}
			if got := minifier.isExcluded(tt.path); got != tt.skipped {
				t.Errorf("minify exclusion = %v, want %v", got, tt.skipped)
			}
			if got := shouldSkipCSSFile(tt.path, nil); got != tt.skipped {
				t.Errorf("purge exclusion = %v, want %v", got, tt.skipped)
			}
		})
	}
}

func TestReconcileFinalCSSAssetsUsesPublishedBytes(t *testing.T) {
	out := t.TempDir()
	if err := os.MkdirAll(filepath.Join(out, "css"), 0o755); err != nil {
		t.Fatal(err)
	}
	initial := []byte(".old { display:block }")
	final := []byte(".new { display:grid }")
	oldHash := fmt.Sprintf("%x", sha256.Sum256(initial))[:8]
	newHash := fmt.Sprintf("%x", sha256.Sum256(final))[:8]
	if err := os.WriteFile(filepath.Join(out, "css", "components.css"), final, 0o600); err != nil {
		t.Fatal(err)
	}
	html := `<link rel="stylesheet" href="/css/components.` + oldHash + `.css"><main class="new"><p>css/components.` + oldHash + `.css</p><pre><code>&lt;link href="/css/components.` + oldHash + `.css"&gt;</code></pre><!-- css/components.` + oldHash + `.css --></main>`
	if err := os.WriteFile(filepath.Join(out, "index.html"), []byte(html), 0o600); err != nil {
		t.Fatal(err)
	}
	m := lifecycle.NewManager()
	m.SetAssetHash("css/components.css", oldHash)
	root, err := os.OpenRoot(out)
	if err != nil {
		t.Fatal(err)
	}
	defer root.Close()
	if err := reconcileFinalCSSAssets(m, root, out); err != nil {
		t.Fatal(err)
	}
	if got := m.GetAssetHash("css/components.css"); got != newHash {
		t.Errorf("hash=%s want=%s", got, newHash)
	}
	rendered, err := os.ReadFile(filepath.Join(out, "index.html"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(rendered), "/css/components."+newHash+".css") {
		t.Fatalf("stale asset reference: %s", rendered)
	}
	// Literal authored content must not be rewritten to the new asset URL.
	if !strings.Contains(string(rendered), "<p>css/components."+oldHash+".css</p>") ||
		!strings.Contains(string(rendered), `&lt;link href="/css/components.`+oldHash+`.css"&gt;`) ||
		!strings.Contains(string(rendered), "<!-- css/components."+oldHash+".css -->") {
		t.Fatalf("authored content was modified: %s", rendered)
	}
	plugin := NewStaticAssetsPlugin()
	if err := plugin.createHashedCopiesWithRoot(m, out, root); err != nil {
		t.Fatal(err)
	}
	published, err := os.ReadFile(filepath.Join(out, "css", "components."+newHash+".css"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(published, final) {
		t.Fatalf("versioned CSS differs from final bytes")
	}
	// Re-running without any CSS change must preserve URLs and file content.
	if err := reconcileFinalCSSAssets(m, root, out); err != nil {
		t.Fatal(err)
	}
	again, err := os.ReadFile(filepath.Join(out, "index.html"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(again, rendered) {
		t.Fatal("unchanged output must be stable")
	}

	// Incremental builds can restore HTML referencing a prior final hash.
	// The reconciliation must not assume the source Configure hash appears.
	cachedHTML := `<link href="/css/components.deadbeef.css" rel="stylesheet"><p>css/components.deadbeef.css</p>`
	if err := os.WriteFile(filepath.Join(out, "cached.html"), []byte(cachedHTML), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := reconcileFinalCSSAssets(m, root, out); err != nil {
		t.Fatal(err)
	}
	cached, err := os.ReadFile(filepath.Join(out, "cached.html"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(cached), `<link href="/css/components.`+newHash+`.css" rel="stylesheet">`) ||
		!strings.Contains(string(cached), "<p>css/components.deadbeef.css</p>") {
		t.Fatalf("cached HTML retains stale CSS reference: %s", cached)
	}
}

// CSS links are finalized after image_library.Write records output ownership.
// An incremental build must not reject Markata's own CSS rewrite as an
// externally modified image-library page.
func TestReconcileFinalCSSAssetsPreservesImageLibraryOwnership(t *testing.T) {
	out := t.TempDir()
	if err := os.MkdirAll(filepath.Join(out, "css"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(out, "images"), 0o755); err != nil {
		t.Fatal(err)
	}
	beforeCSS := []byte(".before { color: red }")
	finalCSS := []byte(".after { color: blue }")
	beforeHash := fmt.Sprintf("%x", sha256.Sum256(beforeCSS))[:8]
	finalHash := fmt.Sprintf("%x", sha256.Sum256(finalCSS))[:8]
	if err := os.WriteFile(filepath.Join(out, "css", "components.css"), finalCSS, 0o600); err != nil {
		t.Fatal(err)
	}
	page := filepath.Join(out, "images", "index.html")
	inputHTML := []byte(`<html><head><link rel="stylesheet" href="/css/components.` + beforeHash + `.css"></head><body>Images</body></html>`)
	if err := os.WriteFile(page, inputHTML, 0o600); err != nil {
		t.Fatal(err)
	}
	cacheDir := t.TempDir()
	cachePath := filepath.Join(cacheDir, buildcache.CacheFileName)
	cache := buildcache.New(cacheDir)
	cache.SetImageLibraryOutputs(out, out, map[string]string{page: buildcache.ContentHash(string(inputHTML))})
	if err := cache.Save(); err != nil {
		t.Fatal(err)
	}

	m := lifecycle.NewManager()
	m.Cache().Set("build_cache", cache)
	m.SetAssetHash("css/components.css", beforeHash)
	root, err := os.OpenRoot(out)
	if err != nil {
		t.Fatal(err)
	}
	defer root.Close()
	if err := reconcileFinalCSSAssets(m, root, out); err != nil {
		t.Fatal(err)
	}
	actual, err := os.ReadFile(page)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(actual, []byte("components."+finalHash+".css")) {
		t.Fatalf("HTML not reconciled: %s", actual)
	}
	want := buildcache.ContentHash(string(actual))
	if tracked, ok := cache.ImageLibraryOutputHash(out, out, page); !ok || tracked != want {
		t.Fatalf("tracked image library hash = %q (owned=%v), want %q", tracked, ok, want)
	}
	persisted, err := os.ReadFile(cachePath)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(persisted, []byte(want)) {
		t.Fatal("rewritten ownership hash was not persisted for the next warm build")
	}
	unowned := filepath.Join(out, "other.html")
	if err := os.WriteFile(unowned, []byte("outside"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := refreshImageLibraryCSSOwnership(cache, out, map[string]string{"other.html": buildcache.ContentHash("outside")}); err != nil {
		t.Fatal(err)
	}
	if cache.OwnsImageLibraryOutput(out, out, unowned) {
		t.Fatal("must not adopt an unrelated HTML file")
	}
	// Configured output directories may be relative to the working directory.
	// Their normalized root must still match the cache's absolute output root.
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	relativeOutput, err := filepath.Rel(wd, out)
	if err != nil {
		// Windows CI may check out on D: while t.TempDir is on C:.
		// Cross-volume paths have no relative representation.
		if filepath.VolumeName(wd) != filepath.VolumeName(out) {
			t.Logf("skipping cross-volume relative output subcheck: %v", err)
			return
		}
		t.Fatal(err)
	}
	if err := refreshImageLibraryCSSOwnership(cache, relativeOutput, map[string]string{
		"images/index.html": want,
	}); err != nil {
		t.Fatal(err)
	}
	if tracked, ok := cache.ImageLibraryOutputHash(out, out, page); !ok || tracked != want {
		t.Fatalf("relative output root lost ownership: %q (owned=%v)", tracked, ok)
	}
}

func TestRewriteCSSLinkReferencesOnlyLocalStylesheets(t *testing.T) {
	matcher := regexp.MustCompile(`^css/components\.[0-9a-fA-F]{8}\.css$`)
	changes := []finalCSSReplacement{{matcher: matcher, destination: "css/components.12345678.css"}}
	input := []byte(`<html><head>
<link rel="stylesheet" href="https://cdn.example/css/components.deadbeef.css">
<link rel="stylesheet" href="//cdn.example/css/components.deadbeef.css">
<link rel="canonical" href="/css/components.deadbeef.css">
<link rel="stylesheet" href="/css/components.deadbeef.css?ref=1#theme">
<link rel="preload" as="style" href="../css/components.deadbeef.css">
<link rel="preload" as="script" href="/css/components.deadbeef.css">
<link rel="stylesheet" href="css/components.deadbeef.css">
<link rel="stylesheet" href="https://cdn.example/css/components.deadbeef.css">
</head><body>css/components.deadbeef.css</body></html>`)
	got, err := rewriteCSSLinkReferences(input, "nested/index.html", changes)
	if err != nil {
		t.Fatal(err)
	}
	wantChanged := []string{
		`href="/css/components.12345678.css?ref=1#theme"`,
		`href="../css/components.12345678.css"`,
	}
	for _, value := range wantChanged {
		if !bytes.Contains(got, []byte(value)) {
			t.Fatalf("expected rewritten local stylesheet %q: %s", value, got)
		}
	}
	wantUntouched := []string{
		`href="https://cdn.example/css/components.deadbeef.css"`,
		`href="//cdn.example/css/components.deadbeef.css"`,
		`rel="canonical" href="/css/components.deadbeef.css"`,
		`as="script" href="/css/components.deadbeef.css"`,
		`href="css/components.deadbeef.css"`,
		`<body>css/components.deadbeef.css</body>`,
	}
	for _, value := range wantUntouched {
		if !bytes.Contains(got, []byte(value)) {
			t.Fatalf("unrelated reference changed %q: %s", value, got)
		}
	}
}
