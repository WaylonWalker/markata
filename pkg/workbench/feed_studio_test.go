package workbench

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/models"
)

func TestFeedStudioMatchAndFilterErrors(t *testing.T) {
	title := "Hello"
	posts := []*models.Post{
		{Path: "pages/hello.md", Title: &title, Published: true},
		{Path: "pages/draft.md", Draft: true},
	}
	matched, err := feedStudioMatch(posts, "published == true")
	if err != nil {
		t.Fatal(err)
	}
	if len(matched) != 1 || matched[0].Path != posts[0].Path {
		t.Fatalf("matches = %#v, want only published post", matched)
	}
	if _, err := feedStudioMatch(posts, "published =="); err == nil {
		t.Fatal("invalid filter unexpectedly accepted")
	}
}

func TestFeedStudioSaveFormatPreservationAndConflicts(t *testing.T) {
	tests := []struct {
		name, filename, source, preserved string
	}{
		{
			name: "toml", filename: "markata-go.toml",
			source:    "# Feed comments should survive\n[[markata-go.feeds]]\nslug = \"blog\"\ntitle = \"Original\"\nfilter = \"published == true\"\nextra_setting = \"preserve-me\"\n",
			preserved: "# Feed comments should survive",
		},
		{
			name: "yaml", filename: "markata-go.yaml",
			source:    "markata-go:\n  feeds:\n    - slug: blog\n      title: Original\n      filter: published == true\n      extra_setting: preserve-me\n",
			preserved: "extra_setting: preserve-me",
		},
		{
			name: "json", filename: "markata-go.json",
			source:    "{\"markata-go\":{\"feeds\":[{\"slug\":\"blog\",\"title\":\"Original\",\"filter\":\"published == true\",\"extra_setting\":\"preserve-me\"}]}}",
			preserved: "preserve-me",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			filename := filepath.Join(root, tc.filename)
			if err := os.WriteFile(filename, []byte(tc.source), 0o600); err != nil {
				t.Fatal(err)
			}
			h := &Handler{
				root: root,
				siteConfig: &models.Config{Feeds: []models.FeedConfig{{
					Slug:   "blog",
					Title:  "Original",
					Filter: "published == true",
				}}},
			}
			expected := feedStudioVersion([]byte(tc.source))
			payload := feedStudioUpdate{
				Slug: "blog", ExpectedVersion: expected,
				Title: "Edited", Description: "New description",
				Filter: "published == true", Sort: "date",
				Reverse: true, Presentation: "gallery",
				Template: "feed-photo-grid.html",
			}
			body, err := json.Marshal(payload)
			if err != nil {
				t.Fatal(err)
			}
			post := func() *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodPost, "http://localhost/workspace/feed-studio", strings.NewReader(string(body)))
				req.Header.Set("Origin", "http://localhost")
				req.Header.Set("Content-Type", "application/json")
				recorder := httptest.NewRecorder()
				h.serveFeedStudio(recorder, req)
				return recorder
			}
			first := post()
			if first.Code != http.StatusOK {
				t.Fatalf("save: status %d: %s", first.Code, first.Body.String())
			}
			saved, err := os.ReadFile(filename)
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(string(saved), tc.preserved) {
				t.Fatalf("unrelated authored config lost: %s", string(saved))
			}
			if !strings.Contains(string(saved), "feed-photo-grid.html") || !strings.Contains(string(saved), "Edited") {
				t.Fatalf("edited fields not saved: %s", string(saved))
			}
			if feedStudioVersion(saved) == expected {
				t.Fatal("version did not advance")
			}
			second := post()
			if second.Code != http.StatusConflict {
				t.Fatalf("stale editor save: got %d, want 409: %s", second.Code, second.Body.String())
			}
			after, err := os.ReadFile(filename)
			if err != nil {
				t.Fatal(err)
			}
			if string(after) != string(saved) {
				t.Fatal("stale edit changed the site config")
			}
		})
	}
}

func TestFeedStudioEditsActualIncludedFeedFile(t *testing.T) {
	root := t.TempDir()
	configDir := filepath.Join(root, "config")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatal(err)
	}
	siteSource := "[markata-go]\ntitle = \"Test site\"\ninclude = [\"config/feeds.toml\"]\n"
	sitePath := filepath.Join(root, "markata-go.toml")
	if err := os.WriteFile(sitePath, []byte(siteSource), 0o600); err != nil {
		t.Fatal(err)
	}
	feedPath := filepath.Join(configDir, "feeds.toml")
	feedSource := "# Preserve this exact comment\n[[markata-go.feeds]]\nslug = \"blog\"\ntitle = \"Original\"\nfilter = \"published == true\"\n"
	if err := os.WriteFile(feedPath, []byte(feedSource), 0o600); err != nil {
		t.Fatal(err)
	}
	sources, err := feedStudioSources(root)
	if err != nil {
		t.Fatal(err)
	}
	source, ok := sources["blog"]
	if !ok || source.path != feedPath || source.index != 0 || source.label != "config/feeds.toml" {
		t.Fatalf("included feed source not discovered: %#v, found=%v", source, ok)
	}
	payload := feedStudioUpdate{
		Slug: "blog", ExpectedVersion: source.version,
		Title: "Updated from Feed Studio", Filter: "published == true",
		Sort: "date", Reverse: true,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	handler := &Handler{root: root}
	req := httptest.NewRequest(http.MethodPost, "http://localhost/workspace/feed-studio", strings.NewReader(string(body)))
	req.Header.Set("Origin", "http://localhost")
	w := httptest.NewRecorder()
	handler.serveFeedStudio(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("included feed save failed: %d %s", w.Code, w.Body.String())
	}
	updated, err := os.ReadFile(feedPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(updated), "Updated from Feed Studio") || !strings.Contains(string(updated), "# Preserve this exact comment") {
		t.Fatalf("included file lost edited fields or comments: %s", string(updated))
	}
	rootAfter, err := os.ReadFile(sitePath)
	if err != nil {
		t.Fatal(err)
	}
	if string(rootAfter) != siteSource {
		t.Fatalf("unexpected change to root config: %s", rootAfter)
	}
	if !strings.Contains(w.Body.String(), "config/feeds.toml") {
		t.Fatalf("save response did not identify the edited file: %s", w.Body.String())
	}
}

func TestFeedStudioSnapshotConfinedToOutput(t *testing.T) {
	root := t.TempDir()
	feedDir := filepath.Join(root, "output", "blog")
	if err := os.MkdirAll(feedDir, 0o755); err != nil {
		t.Fatal(err)
	}
	page := filepath.Join(feedDir, "index.html")
	if err := os.WriteFile(page, []byte("<h1>Actual site output</h1>"), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg := &models.Config{OutputDir: "output"}
	got, err := feedSnapshotFile(root, cfg, "blog/index.html")
	if err != nil || got != page {
		t.Fatalf("built HTML path = %q, %v", got, err)
	}
	if _, err := feedSnapshotFile(root, cfg, "../anything"); err == nil {
		t.Fatal("directory traversal accepted")
	}
	outside := t.TempDir()
	if err := os.WriteFile(filepath.Join(outside, "secret.html"), []byte("secret"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(outside, "secret.html"), filepath.Join(feedDir, "leak.html")); err == nil {
		if _, err := feedSnapshotFile(root, cfg, "blog/leak.html"); err == nil {
			t.Fatal("outside-workspace symlink accepted")
		}
	}
	if _, err := feedSnapshotFile(root, &models.Config{OutputDir: outside}, "secret.html"); err == nil {
		t.Fatal("arbitrary external output directory accepted")
	}
}
