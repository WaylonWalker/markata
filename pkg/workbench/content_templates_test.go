package workbench

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadSiteContentTemplates_SharedMarkdownStarter(t *testing.T) {
	root := t.TempDir()
	directory := filepath.Join(root, "content-templates")
	if err := os.MkdirAll(directory, 0o755); err != nil {
		t.Fatal(err)
	}
	content := "---\nstarter:\n  name: guide\n  description: A reusable guide\n  aliases:\n    - tutorial\n_directory: guides\ntemplateKey: guide\ndescription: Page description\ntitle: '{{title}}'\n---\n\n## {{title}}\n"
	if err := os.WriteFile(filepath.Join(directory, "guide.md"), []byte(content), 0o644); err != nil {
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
	if err := os.WriteFile(filepath.Join(directory, "guide.md"), []byte(content), 0o644); err != nil {
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
