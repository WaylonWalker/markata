package plugins

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

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
	html := `<link rel="stylesheet" href="/css/components.` + oldHash + `.css"><main class="new"></main>`
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
	cachedHTML := `<link href="/css/components.deadbeef.css" rel="stylesheet">`
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
	if strings.Contains(string(cached), "deadbeef") || !strings.Contains(string(cached), "components."+newHash+".css") {
		t.Fatalf("cached HTML retains stale CSS reference: %s", cached)
	}
}
