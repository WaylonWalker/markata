package workbench

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/models"
)

func TestLoadSiteContentTemplates_SharedMarkdownStarter(t *testing.T) {
	root := t.TempDir()
	directory := filepath.Join(root, "content-templates")
	if err := os.MkdirAll(directory, 0o755); err != nil {
		t.Fatal(err)
	}
	content := "---\nstarter:\n  name: guide\n  description: A reusable guide\n  aliases:\n    - tutorial\n_directory: guides\ntemplateKey: guide\ndescription: Page description\ntitle: '{{title}}'\n---\n\n## {{title}}\n"
	if err := os.WriteFile(filepath.Join(directory, "guide.md"), []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	templates, err := loadSiteContentTemplates(root, "content-templates")
	if err != nil {
		t.Fatal(err)
	}
	if len(templates) != 1 {
		t.Fatalf("got %d templates, want 1", len(templates))
	}
	template := templates[0]
	if template.Name != "guide" || template.Description != "A reusable guide" || template.Directory != "guides" {
		t.Fatalf("unexpected template metadata: %#v", template)
	}
	if len(template.Aliases) != 1 || template.Aliases[0] != "tutorial" {
		t.Fatalf("aliases = %#v", template.Aliases)
	}
	if strings.Contains(template.Source, "_directory:") || strings.Contains(template.Source, "starter:") || strings.Contains(template.Source, "aliases:") {
		t.Fatalf("discovery metadata leaked into starter source:\n%s", template.Source)
	}
	if !strings.Contains(template.Source, "templateKey: guide") || !strings.Contains(template.Source, "description: Page description") || !strings.Contains(template.Source, "## {{title}}") {
		t.Fatalf("content fields were not preserved:\n%s", template.Source)
	}
}

func TestLoadSiteContentTemplates_DefaultDirectoryUsesCanonicalName(t *testing.T) {
	root := t.TempDir()
	directory := filepath.Join(root, "content-templates")
	if err := os.MkdirAll(directory, 0o755); err != nil {
		t.Fatal(err)
	}
	content := "---\nstarter:\n  name: how-to\n---\n\nBody\n"
	if err := os.WriteFile(filepath.Join(directory, "guide.md"), []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	templates, err := loadSiteContentTemplates(root, "content-templates")
	if err != nil {
		t.Fatal(err)
	}
	if len(templates) != 1 || templates[0].Name != "how-to" || templates[0].Directory != "how-to" {
		t.Fatalf("template name/directory = %#v", templates)
	}
}

func TestValidateContentTemplateNames_RejectsDuplicateCanonicalName(t *testing.T) {
	templates := []contentTemplate{
		{ID: "local:guide.md", Name: "Guide"},
		{ID: "local:duplicate.md", Name: "guide"},
	}
	if err := validateContentTemplateNames(templates); err == nil {
		t.Fatal("expected duplicate canonical name error")
	}
}

func TestContentTemplates_ConfigPrecedenceAndPlacement(t *testing.T) {
	root := t.TempDir()
	cfg := &models.Config{}
	cfg.ContentTemplates.Placement = map[string]string{"note": "custom-notes"}
	cfg.ContentTemplates.Templates = []models.ContentTemplateConfig{
		{Name: "note", Directory: "from-config", Frontmatter: map[string]interface{}{"title": "{{title}}"}, Body: "configured"},
		{Name: "custom", Directory: "from-config", Frontmatter: map[string]interface{}{"published": false}, Body: "custom-body"},
	}
	handler := &Handler{root: root, siteConfig: cfg}
	templates, err := handler.contentTemplates()
	if err != nil {
		t.Fatal(err)
	}
	byName := make(map[string]contentTemplate)
	for _, item := range templates {
		byName[contentTemplateKey(item)] = item
	}
	if got := byName["note"]; got.Directory != "from-config" || !strings.Contains(got.Source, "configured") {
		t.Fatalf("config template must override built-in and placement: %#v", got)
	}
	if got := byName["custom"]; got.Directory != "from-config" || !strings.Contains(got.Source, "published: false") {
		t.Fatalf("custom config template missing: %#v", got)
	}

	dir := filepath.Join(root, "content-templates")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	source := "---\nstarter:\n  name: note\n_directory: from-markdown\ntitle: '{{title}}'\n---\nMarkdown winner"
	if err := os.WriteFile(filepath.Join(dir, "note.md"), []byte(source), 0o600); err != nil {
		t.Fatal(err)
	}
	templates, err = handler.contentTemplates()
	if err != nil {
		t.Fatal(err)
	}
	for _, item := range templates {
		if contentTemplateKey(item) == "note" {
			if item.Directory != "from-markdown" || !strings.Contains(item.Source, "Markdown winner") {
				t.Fatalf("Markdown starter must override config: %#v", item)
			}
			return
		}
	}
	t.Fatal("missing note starter")
}

func TestSiteContentTemplate_RejectsUnclosedFrontmatter(t *testing.T) {
	_, err := siteContentTemplateFromSource("broken.md", "broken", []byte("---\nstarter:\n  name: broken\nbody"))
	if err == nil || !strings.Contains(err.Error(), "broken.md") {
		t.Fatalf("expected named malformed frontmatter error, got %v", err)
	}
}

func TestSiteContentTemplate_PreservesBodyWithoutFrontmatter(t *testing.T) {
	template, err := siteContentTemplateFromSource("plain.md", "plain", []byte("Just {{title}}"))
	if err != nil || template.Source != "Just {{title}}" {
		t.Fatalf("plain starter must remain sparse: %#v %v", template, err)
	}
}
