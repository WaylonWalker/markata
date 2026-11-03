package config

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCanonicalSpecimenContract(t *testing.T) {
	template, err := os.ReadFile(filepath.Join("..", "..", "templates", "post.html"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(template), `{{ body | safe }}`) {
		t.Fatal("canonical specimen template is missing body rendering")
	}
	base, err := os.ReadFile(filepath.Join("..", "..", "templates", "base.html"))
	if err != nil {
		t.Fatal(err)
	}
	for _, marker := range []string{`data-rendering-specimen="canonical-headings"`, `data-rendering-profile="canonical-document-v1"`, `data-rendering-content-width="68ch"`, `data-rendering-viewport="1280x1000"`, `data-rendering-layer-roles="surface,document,heading-wear,motif,texture"`} {
		if !strings.Contains(string(base), marker) {
			t.Errorf("canonical specimen root is missing %q", marker)
		}
	}
	if !strings.Contains(string(base), `canonical_rendering_fixture`) || !strings.Contains(string(base), `css/canonical-specimen.css`) {
		t.Fatal("canonical specimen stylesheet must load only for its fixture page")
	}
}

func TestCanonicalSpecimenCaptureGeometry(t *testing.T) {
	base, err := os.ReadFile(filepath.Join("..", "..", "pkg", "themes", "default", "static", "css", "canonical-specimen.css"))
	if err != nil {
		t.Fatal(err)
	}
	markup := string(base)
	for _, marker := range []string{
		"width: 1280px !important;",
		"height: 1000px !important;",
		"min-height: 1000px !important;",
		"width: 680px;",
		"padding: 24px;",
		"body:has([data-rendering-specimen=\"canonical-headings\"]) > .site-header",
	} {
		if !strings.Contains(markup, marker) {
			t.Errorf("canonical capture geometry is missing %q", marker)
		}
	}
	if strings.Contains(markup, "width: 100vw;") || strings.Contains(markup, "height: 100vh;") {
		t.Fatal("canonical capture geometry must not depend on viewport units")
	}
}

func TestCanonicalTemplateTreesHaveEquivalentHeadingProfile(t *testing.T) {
	cssPath := filepath.Join("..", "..", "pkg", "themes", "default", "static", "css", "canonical-specimen.css")
	data, err := os.ReadFile(cssPath)
	if err != nil {
		t.Fatal(err)
	}
	profile := string(data)
	if strings.Contains(profile, "first-of-type") || strings.Contains(profile, "font-size: 56px") {
		t.Fatalf("stale H1 workaround remains in %s", cssPath)
	}
	for _, selector := range []string{".post-content h1", ".post-content h2"} {
		if !strings.Contains(profile, `[data-rendering-specimen="canonical-headings"] `+selector+" {") {
			t.Errorf("%s missing %s rule", cssPath, selector)
		}
	}
	for _, templatePath := range []string{filepath.Join("..", "..", "templates", "base.html"), filepath.Join("..", "..", "pkg", "themes", "default", "templates", "base.html")} {
		template, err := os.ReadFile(templatePath)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(template), `{% if canonical_rendering_fixture %}`) || !strings.Contains(string(template), `css/canonical-specimen.css`) {
			t.Errorf("%s must load the shared specimen stylesheet conditionally", templatePath)
		}
	}
}

func TestCanonicalDocumentProfileExists(t *testing.T) {
	profilePath := filepath.Join("..", "..", "spec", "rendering-contract", "profiles", "canonical-document-v1.json")
	profile, err := os.ReadFile(profilePath)
	if err != nil {
		t.Fatal(err)
	}
	var document struct {
		SchemaVersion string `json:"schema_version"`
		Internal      bool   `json:"internal"`
		Viewport      struct {
			Width  int `json:"width"`
			Height int `json:"height"`
		} `json:"viewport"`
		ContentWidth string `json:"content_width"`
	}
	if err := json.Unmarshal(profile, &document); err != nil {
		t.Fatal(err)
	}
	if document.SchemaVersion != "canonical-document-v1" || !document.Internal || document.Viewport.Width != 1280 || document.Viewport.Height != 1000 || document.ContentWidth != "68ch" {
		t.Fatalf("unexpected canonical profile metadata: %+v", document)
	}
}

func TestCanonicalSpecimenExternalFixtureIdentity(t *testing.T) {
	local, err := os.ReadFile(filepath.Join("..", "..", "rendering-fixtures", "test-headings.md"))
	if err != nil {
		t.Fatal(err)
	}
	local = canonicalFixtureBytes(local)
	if got := sha256Digest(local); got != "58bc0ab1e501adbe498c2902f1e0ede5112c8deede94e8b3a278f16613e9cedb" {
		t.Fatalf("local canonical fixture hash drifted: %s", got)
	}
	fixture := filepath.Join("..", "..", "..", "fixtures", "canonical-headings", "test-headings.md")
	shared, err := os.ReadFile(fixture)
	if err != nil {
		if os.IsNotExist(err) {
			t.Skipf("shared canonical fixture is unavailable: %v", err)
		}
		t.Fatal(err)
	}
	shared = canonicalFixtureBytes(shared)
	if sha256Digest(shared) != sha256Digest(local) || sha256Digest(shared) != "58bc0ab1e501adbe498c2902f1e0ede5112c8deede94e8b3a278f16613e9cedb" {
		t.Fatalf("canonical fixture hash drifted: shared=%s local=%s", sha256Digest(shared), sha256Digest(local))
	}

	theme, err := os.ReadFile(filepath.Join("..", "..", "..", "fixtures", "canonical-headings", "theme.toml"))
	if err != nil {
		t.Fatal(err)
	}
	if got := sha256Digest(theme); got != "c751df3b26ae476615d22107adfd5d00a4b76c661b613cbc44bd8499f7076806" {
		t.Fatalf("canonical theme fixture hash drifted: %s", got)
	}
}

func sha256Digest(data []byte) string {
	digest := sha256.Sum256(data)
	return hex.EncodeToString(digest[:])
}
