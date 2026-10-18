package plugins

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

func TestPublishHTMLPlugin_ExplicitHomepageWritesAllEnabledFormats(t *testing.T) {
	outputDir := t.TempDir()
	config := &lifecycle.Config{
		OutputDir: outputDir,
		Extra: map[string]interface{}{
			"url":   "https://example.com",
			"title": "Example",
			"post_formats": models.PostFormatsConfig{
				Markdown: true,
				Text:     true,
				ANSI:     true,
				OG:       true,
			},
		},
	}
	title := "Homepage"
	post := &models.Post{
		Path:        "pages/post/index.md",
		Slug:        "",
		Title:       &title,
		Content:     "# Homepage\n",
		HTML:        "<html><body>Homepage</body></html>",
		ArticleHTML: "<h1>Homepage</h1>",
		Published:   true,
	}
	post.Set("_slug_explicit", true)

	plugin := NewPublishHTMLPlugin()
	if err := plugin.writePost(post, config, nil, createTestManager(t, config)); err != nil {
		t.Fatalf("writePost() error = %v", err)
	}

	for _, relative := range []string{
		"index.html",
		"index.md",
		"index.txt",
		"index.ansi",
		filepath.Join("og", "index.html"),
	} {
		if info, err := os.Stat(filepath.Join(outputDir, relative)); err != nil || !info.Mode().IsRegular() {
			t.Errorf("expected root output %q: info=%v err=%v", relative, info, err)
		}
	}
}

func TestPublishHTMLPlugin_HomepageMigratesLegacyFormatRedirects(t *testing.T) {
	outputDir := t.TempDir()
	plugin := NewPublishHTMLPlugin()
	for _, ext := range []string{"md", "txt"} {
		// Reproduce the pre-fix writer's empty-slug output, as copied from a release.
		if err := plugin.writeRegularFormatOutput("", ext, "old", outputDir, true); err != nil {
			t.Fatal(err)
		}
		if err := plugin.writeReversedFormatOutput("", ext, "new homepage", outputDir, true); err != nil {
			t.Fatal(err)
		}
		data, err := os.ReadFile(filepath.Join(outputDir, "index."+ext))
		if err != nil || string(data) != "new homepage" {
			t.Fatalf("%s output = %q, error = %v", ext, data, err)
		}
	}
}

func TestPublishHTMLPlugin_HomepagePreservesUnrecognizedDirectories(t *testing.T) {
	for _, scenario := range []string{"extra file", "different HTML", "symlink"} {
		t.Run(scenario, func(t *testing.T) {
			outputDir := t.TempDir()
			plugin := NewPublishHTMLPlugin()
			if err := plugin.writeRegularFormatOutput("", "md", "old", outputDir, true); err != nil {
				t.Fatal(err)
			}
			dir := filepath.Join(outputDir, "index.md")
			redirect := filepath.Join(dir, "index.html")
			switch scenario {
			case "extra file":
				if err := os.WriteFile(filepath.Join(dir, "notes.txt"), []byte("keep"), 0o600); err != nil {
					t.Fatal(err)
				}
			case "different HTML":
				if err := os.WriteFile(redirect, []byte("user page"), 0o600); err != nil {
					t.Fatal(err)
				}
			case "symlink":
				data, err := os.ReadFile(redirect)
				if err != nil {
					t.Fatal(err)
				}
				target := filepath.Join(outputDir, "user.html")
				if err := os.WriteFile(target, data, 0o600); err != nil {
					t.Fatal(err)
				}
				if err := os.Remove(redirect); err != nil {
					t.Fatal(err)
				}
				if err := os.Symlink(target, redirect); err != nil {
					t.Skipf("symlinks unavailable: %v", err)
				}
			}
			before, err := os.ReadFile(redirect)
			if err != nil {
				t.Fatal(err)
			}
			if err := plugin.writeReversedFormatOutput("", "md", "new", outputDir, true); err == nil {
				t.Fatal("expected conflicting directory to fail")
			}
			after, err := os.ReadFile(redirect)
			if err != nil || !bytes.Equal(after, before) {
				t.Fatalf("conflicting contents changed: %q, %v", after, err)
			}
			if scenario == "extra file" {
				data, err := os.ReadFile(filepath.Join(dir, "notes.txt"))
				if err != nil || string(data) != "keep" {
					t.Fatal("user file changed")
				}
			}
		})
	}
}
