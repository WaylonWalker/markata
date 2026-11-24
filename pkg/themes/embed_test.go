package themes

import (
	"io/fs"
	"strings"
	"testing"
)

func TestDefaultTemplates_ReturnsValidFS(t *testing.T) {
	templates := DefaultTemplates()
	if templates == nil {
		t.Fatal("DefaultTemplates() returned nil")
	}

	// Should be able to read the base template
	_, err := fs.ReadFile(templates, "base.html")
	if err != nil {
		t.Errorf("failed to read base.html from templates: %v", err)
	}
}

func TestDefaultTemplates_ContainsExpectedFiles(t *testing.T) {
	templates := DefaultTemplates()
	if templates == nil {
		t.Fatal("DefaultTemplates() returned nil")
	}

	expectedFiles := []string{
		"base.html",
		"post.html",
		"feed.html",
		"layouts/docs.html",
		"layouts/landing.html",
		"card.html",
		"reader.html",
	}

	for _, file := range expectedFiles {
		t.Run(file, func(t *testing.T) {
			_, err := fs.ReadFile(templates, file)
			if err != nil {
				t.Errorf("expected template %q not found: %v", file, err)
			}
		})
	}
}

func TestDefaultStatic_ReturnsValidFS(t *testing.T) {
	static := DefaultStatic()
	if static == nil {
		t.Fatal("DefaultStatic() returned nil")
	}

	// Should be able to read the main CSS file
	_, err := fs.ReadFile(static, "css/main.css")
	if err != nil {
		t.Errorf("failed to read css/main.css from static: %v", err)
	}
}

func TestDefaultStatic_ContainsExpectedFiles(t *testing.T) {
	static := DefaultStatic()
	if static == nil {
		t.Fatal("DefaultStatic() returned nil")
	}

	expectedFiles := []string{
		"css/main.css",
		"css/variables.css",
		"css/layouts.css",
		"css/code.css",
		"css/components.css",
		"css/admonitions.css",
	}

	for _, file := range expectedFiles {
		t.Run(file, func(t *testing.T) {
			_, err := fs.ReadFile(static, file)
			if err != nil {
				t.Errorf("expected static file %q not found: %v", file, err)
			}
		})
	}
}

func TestDefaultTheme_ReturnsValidFS(t *testing.T) {
	theme := DefaultTheme()
	if theme == nil {
		t.Fatal("DefaultTheme() returned nil")
	}

	// Should be able to read theme.toml
	_, err := fs.ReadFile(theme, "theme.toml")
	if err != nil {
		t.Errorf("failed to read theme.toml from theme: %v", err)
	}
}

func TestDefaultTheme_ContainsTemplatesAndStatic(t *testing.T) {
	theme := DefaultTheme()
	if theme == nil {
		t.Fatal("DefaultTheme() returned nil")
	}

	tests := []struct {
		name string
		path string
	}{
		{"theme config", "theme.toml"},
		{"templates directory", "templates/base.html"},
		{"static directory", "static/css/main.css"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := fs.ReadFile(theme, tt.path)
			if err != nil {
				t.Errorf("expected %q not found: %v", tt.path, err)
			}
		})
	}
}

func TestReadTemplate_ExistingFile(t *testing.T) {
	tests := []struct {
		name     string
		template string
	}{
		{"base template", "base.html"},
		{"post template", "post.html"},
		{"feed template", "feed.html"},
		{"card template", "card.html"},
		{"reader template", "reader.html"},
		{"nested component", "components/nav.html"},
		{"nested partial", "partials/head.html"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content, err := ReadTemplate(tt.template)
			if err != nil {
				t.Errorf("ReadTemplate(%q) returned error: %v", tt.template, err)
			}
			if len(content) == 0 {
				t.Errorf("ReadTemplate(%q) returned empty content", tt.template)
			}
		})
	}
}

func TestReadTemplate_NonExistentFile(t *testing.T) {
	_, err := ReadTemplate("nonexistent.html")
	if err == nil {
		t.Error("ReadTemplate(nonexistent.html) should return error")
	}
}

func TestReadTemplate_ContentValidity(t *testing.T) {
	// Read base.html and verify it contains expected HTML structure
	content, err := ReadTemplate("base.html")
	if err != nil {
		t.Fatalf("ReadTemplate(base.html) failed: %v", err)
	}

	contentStr := string(content)

	// Base template should have DOCTYPE and html tags
	if len(contentStr) < 50 {
		t.Error("base.html content seems too short to be valid")
	}
}

func TestReadStatic_ExistingFile(t *testing.T) {
	tests := []struct {
		name string
		file string
	}{
		{"main CSS", "css/main.css"},
		{"variables CSS", "css/variables.css"},
		{"text-size CSS", "css/text-size.css"},
		{"layouts CSS", "css/layouts.css"},
		{"code CSS", "css/code.css"},
		{"components CSS", "css/components.css"},
		{"admonitions CSS", "css/admonitions.css"},
		{"scroll spy JS", "js/scroll-spy.js"},
		{"text-size JS", "js/text-size.js"},
		{"pagination JS", "js/pagination.js"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content, err := ReadStatic(tt.file)
			if err != nil {
				t.Errorf("ReadStatic(%q) returned error: %v", tt.file, err)
			}
			if len(content) == 0 {
				t.Errorf("ReadStatic(%q) returned empty content", tt.file)
			}
		})
	}
}

func TestDefaultTheme_PhotoCardsAreIncludedInPagination(t *testing.T) {
	template, err := fs.ReadFile(DefaultTheme(), "templates/partials/cards/photo-card.html")
	if err != nil {
		t.Fatalf("read photo-card template: %v", err)
	}
	if !strings.Contains(string(template), "data-card") {
		t.Fatal("photo-card template must mark photo figures as cards")
	}

	pagination, err := ReadStatic("js/pagination.js")
	if err != nil {
		t.Fatalf("read pagination script: %v", err)
	}
	if !strings.Contains(string(pagination), "[data-card]") {
		t.Fatal("pagination script must include data-card elements")
	}
}

func TestReadStatic_NonExistentFile(t *testing.T) {
	_, err := ReadStatic("nonexistent.css")
	if err == nil {
		t.Error("ReadStatic(nonexistent.css) should return error")
	}
}

func TestListTemplates_ReturnsFiles(t *testing.T) {
	templates, err := ListTemplates()
	if err != nil {
		t.Fatalf("ListTemplates() returned error: %v", err)
	}

	if len(templates) == 0 {
		t.Error("ListTemplates() returned empty list")
	}

	// Verify some expected templates are in the list
	expectedTemplates := map[string]bool{
		"base.html":   false,
		"post.html":   false,
		"feed.html":   false,
		"card.html":   false,
		"reader.html": false,
	}

	for _, tmpl := range templates {
		if _, ok := expectedTemplates[tmpl]; ok {
			expectedTemplates[tmpl] = true
		}
	}

	for tmpl, found := range expectedTemplates {
		if !found {
			t.Errorf("expected template %q not found in ListTemplates() result", tmpl)
		}
	}
}

func TestListTemplates_IncludesNestedFiles(t *testing.T) {
	templates, err := ListTemplates()
	if err != nil {
		t.Fatalf("ListTemplates() returned error: %v", err)
	}

	// Check for nested files in components/ and partials/
	hasComponents := false
	hasPartials := false

	for _, tmpl := range templates {
		if len(tmpl) > 11 && tmpl[:11] == "components/" {
			hasComponents = true
		}
		if len(tmpl) > 9 && tmpl[:9] == "partials/" {
			hasPartials = true
		}
	}

	if !hasComponents {
		t.Error("ListTemplates() should include files from components/ subdirectory")
	}
	if !hasPartials {
		t.Error("ListTemplates() should include files from partials/ subdirectory")
	}
}

func TestListTemplates_NoDuplicates(t *testing.T) {
	templates, err := ListTemplates()
	if err != nil {
		t.Fatalf("ListTemplates() returned error: %v", err)
	}

	seen := make(map[string]bool)
	for _, tmpl := range templates {
		if seen[tmpl] {
			t.Errorf("duplicate template in list: %q", tmpl)
		}
		seen[tmpl] = true
	}
}

func TestListStatic_ReturnsFiles(t *testing.T) {
	files, err := ListStatic()
	if err != nil {
		t.Fatalf("ListStatic() returned error: %v", err)
	}

	if len(files) == 0 {
		t.Error("ListStatic() returned empty list")
	}

	// Verify some expected static files are in the list
	expectedFiles := map[string]bool{
		"css/main.css":      false,
		"css/variables.css": false,
		"css/layouts.css":   false,
	}

	for _, file := range files {
		if _, ok := expectedFiles[file]; ok {
			expectedFiles[file] = true
		}
	}

	for file, found := range expectedFiles {
		if !found {
			t.Errorf("expected static file %q not found in ListStatic() result", file)
		}
	}
}

func TestListStatic_IncludesNestedFiles(t *testing.T) {
	files, err := ListStatic()
	if err != nil {
		t.Fatalf("ListStatic() returned error: %v", err)
	}

	// Check for files in css/ and js/ subdirectories
	hasCSS := false
	hasJS := false

	for _, file := range files {
		if len(file) > 4 && file[:4] == "css/" {
			hasCSS = true
		}
		if len(file) > 3 && file[:3] == "js/" {
			hasJS = true
		}
	}

	if !hasCSS {
		t.Error("ListStatic() should include files from css/ subdirectory")
	}
	if !hasJS {
		t.Error("ListStatic() should include files from js/ subdirectory")
	}
}

func TestListStatic_NoDuplicates(t *testing.T) {
	files, err := ListStatic()
	if err != nil {
		t.Fatalf("ListStatic() returned error: %v", err)
	}

	seen := make(map[string]bool)
	for _, file := range files {
		if seen[file] {
			t.Errorf("duplicate static file in list: %q", file)
		}
		seen[file] = true
	}
}

func TestListStatic_NoDirectories(t *testing.T) {
	files, err := ListStatic()
	if err != nil {
		t.Fatalf("ListStatic() returned error: %v", err)
	}

	static := DefaultStatic()
	for _, file := range files {
		info, err := fs.Stat(static, file)
		if err != nil {
			t.Errorf("could not stat file %q: %v", file, err)
			continue
		}
		if info.IsDir() {
			t.Errorf("ListStatic() should not include directories, but found: %q", file)
		}
	}
}

func TestListTemplates_NoDirectories(t *testing.T) {
	templates, err := ListTemplates()
	if err != nil {
		t.Fatalf("ListTemplates() returned error: %v", err)
	}

	templatesFS := DefaultTemplates()
	for _, tmpl := range templates {
		info, err := fs.Stat(templatesFS, tmpl)
		if err != nil {
			t.Errorf("could not stat file %q: %v", tmpl, err)
			continue
		}
		if info.IsDir() {
			t.Errorf("ListTemplates() should not include directories, but found: %q", tmpl)
		}
	}
}

func TestReadTemplate_AllListedTemplatesReadable(t *testing.T) {
	templates, err := ListTemplates()
	if err != nil {
		t.Fatalf("ListTemplates() returned error: %v", err)
	}

	for _, tmpl := range templates {
		t.Run(tmpl, func(t *testing.T) {
			content, err := ReadTemplate(tmpl)
			if err != nil {
				t.Errorf("ReadTemplate(%q) failed: %v", tmpl, err)
			}
			if len(content) == 0 {
				t.Errorf("ReadTemplate(%q) returned empty content", tmpl)
			}
		})
	}
}

func TestReadStatic_AllListedFilesReadable(t *testing.T) {
	files, err := ListStatic()
	if err != nil {
		t.Fatalf("ListStatic() returned error: %v", err)
	}

	for _, file := range files {
		t.Run(file, func(t *testing.T) {
			content, err := ReadStatic(file)
			if err != nil {
				t.Errorf("ReadStatic(%q) failed: %v", file, err)
			}
			if len(content) == 0 {
				t.Errorf("ReadStatic(%q) returned empty content", file)
			}
		})
	}
}
