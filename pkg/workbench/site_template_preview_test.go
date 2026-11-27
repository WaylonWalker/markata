package workbench

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/models"
)

func TestSiteTemplatePreviewUsesUnsavedMediaAndLink(t *testing.T) {
	h := newHandler(t)
	for _, tc := range []struct {
		name, source, body, want string
	}{
		{"photo", "---\ntemplate: photo\nimage: https://dropper.waylonwalker.com/file/skate.webp\n---\n", "", `src="https://dropper.waylonwalker.com/file/skate.webp"`},
		{"shot file", "---\ntemplate: shot.html\nimage: /images/shot.webp\n---\n", "Caption", `src="/workspace/image?path=images%2Fshot.webp"`},
		{"legacy", "---\ntemplateKey: shots\nimage: https://example.com/shot.webp\n---\n", "", `src="https://example.com/shot.webp"`},
		{"link", "---\ntemplate: link.html\nurl: https://example.com/article\n---\n", "", `href="https://example.com/article"`},
		{"per-format", "---\ntemplates:\n  html: photo.html\nimage: /images/cover.webp\n---\n", "", `src="/workspace/image?path=images%2Fcover.webp"`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			for _, initial := range []bool{true, false} {
				input, err := json.Marshal(map[string]interface{}{
					"markdown": tc.body, "frontmatter_source": tc.source,
					"source_path": "pages/photo/draft.md", "include_document": initial,
				})
				if err != nil {
					t.Fatal(err)
				}
				response := request(h, string(input))
				if response.Code != http.StatusOK {
					t.Fatalf("render failed: %d %s", response.Code, response.Body.String())
				}
				var result renderResponse
				if err := json.Unmarshal(response.Body.Bytes(), &result); err != nil {
					t.Fatal(err)
				}
				if !strings.Contains(result.HTML, tc.want) {
					t.Fatalf("missing %q: %s", tc.want, result.HTML)
				}
				if tc.body == "" && len(result.SourceLines) != 0 {
					t.Fatalf("template generated a source line: %v", result.SourceLines)
				}
			}
		})
	}
}

func TestSiteTemplatePreviewCompanionAndFallback(t *testing.T) {
	h := newHandler(t)
	h.root = t.TempDir()
	h.siteConfig = &models.Config{TemplatesDir: "templates", Title: "Test Site"}
	dir := filepath.Join(h.root, "templates", "editor")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	fragment := filepath.Join(dir, "project.html")
	if err := os.WriteFile(fragment, []byte(`<header>{{ post.title }}: {{ post.Extra.status }} ({{ config.title }})</header>`), 0o600); err != nil {
		t.Fatal(err)
	}
	source := "---\ntemplate: project\ntitle: My project\nstatus: In progress\n---\n"
	got := h.renderSiteFrontmatterPreview(source, "Body", "pages/project/test.md")
	if !strings.Contains(got, "My project: In progress (Test Site)") || !strings.Contains(got, "workbench-template-preview") {
		t.Fatalf("custom fragment missing: %s", got)
	}
	got = h.renderSiteFrontmatterPreview(strings.ReplaceAll(source, "In progress", "Done"), "Body", "pages/project/test.md")
	if !strings.Contains(got, "My project: Done") || strings.Contains(got, "In progress") {
		t.Fatalf("stale custom fragment: %s", got)
	}
	if err := os.WriteFile(filepath.Join(dir, "status.html"), []byte(`<strong>{{ post.Extra.status }}</strong>`), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(fragment, []byte(`{% include "status.html" %}`), 0o600); err != nil {
		t.Fatal(err)
	}
	got = h.renderSiteFrontmatterPreview(source, "Body", "pages/project/test.md")
	if !strings.Contains(got, "<strong>In progress</strong>") {
		t.Fatalf("local include missing: %s", got)
	}
	if err := os.WriteFile(fragment, []byte(`{% include "project.html" %}`), 0o600); err != nil {
		t.Fatal(err)
	}
	got = h.renderSiteFrontmatterPreview(source, "Body", "pages/project/test.md")
	if !strings.Contains(got, "Site template preview unavailable") {
		t.Fatalf("recursive include did not fall back: %s", got)
	}
	if err := os.WriteFile(fragment, []byte(`{% include "../outside.html" %}`), 0o600); err != nil {
		t.Fatal(err)
	}
	got = h.renderSiteFrontmatterPreview(source, "Body", "pages/project/test.md")
	if !strings.Contains(got, "Site template preview unavailable") || !strings.Contains(got, "My project") {
		t.Fatalf("invalid custom fragment did not fall back: %s", got)
	}
}
