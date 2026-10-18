package templates

import (
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/models"
)

func editorialFixture() map[string]interface{} {
	return map[string]interface{}{
		"schema_version": 1, "author": "editor", "role": "copy-editor", "action": "edited",
		"base_commit": strings.Repeat("a", 40), "commit": strings.Repeat("b", 40),
		"path": "pages/a post's.md", "summary": "Fixed a typo.", "recorded_by": "importer",
		"file_url": "https://github.com/owner/repo/blob/" + strings.Repeat("b", 40) + "/pages/a%20post%27s.md",
		"diff_url": "javascript:alert(1)",
	}
}

func TestEditorialHistory_RevisionValidation(t *testing.T) {
	for _, tc := range []struct {
		name, field string
		value       interface{}
		wantLinks   bool
	}{
		{"valid encoded path", "summary", "Valid", true},
		{"unsafe URL", "file_url", "javascript:alert(1)", false},
		{"foreign host", "file_url", "https://github.com.evil.test/owner/repo/blob/" + strings.Repeat("b", 40) + "/pages/a%20post%27s.md", false},
		{"userinfo", "file_url", "https://evil@github.com/owner/repo/blob/" + strings.Repeat("b", 40) + "/pages/a%20post%27s.md", false},
		{"path traversal", "path", "../secret.md", false},
		{"absolute path", "path", "/secret.md", false},
		{"wrong path", "path", "pages/other.md", false},
		{"wrong revision", "commit", strings.Repeat("c", 40), false},
		{"short revision", "commit", "bbbbbbb", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			record := editorialFixture()
			record[tc.field] = tc.value
			post := &models.Post{Extra: map[string]interface{}{"editorial_history": []interface{}{record}}}
			rows := editorialHistoryToMaps(post)
			if len(rows) != 1 {
				t.Fatalf("rows = %#v", rows)
			}
			_, linked := rows[0]["file_url"]
			if linked != tc.wantLinks {
				t.Fatalf("links = %v, want %v: %#v", linked, tc.wantLinks, rows[0])
			}
			diffURL, diffOK := rows[0]["diff_url"].(string)
			if tc.wantLinks && !diffOK {
				t.Fatalf("missing diff URL: %#v", rows[0])
			}
			if tc.wantLinks && (!strings.HasPrefix(diffURL, "https://github.com/owner/repo/commit/") || rows[0]["before_url"] == nil) {
				t.Fatalf("missing validated revision links: %#v", rows[0])
			}
		})
	}
}

func TestEditorialHistory_MalformedAndReviewed(t *testing.T) {
	valid := editorialFixture()
	reviewed := editorialFixture()
	reviewed["action"] = "reviewed"
	reviewed["author"] = "human"
	unsupported := editorialFixture()
	unsupported["schema_version"] = 99
	post := &models.Post{Extra: map[string]interface{}{"editorial_history": []interface{}{valid, "bad", unsupported, reviewed}},
		AuthorObjects: []models.Author{{ID: "human", Name: "A Human"}}}
	rows := editorialHistoryToMaps(post)
	if len(rows) != 2 || rows[0]["name"] != "A Human" || rows[0]["before_url"] != nil || rows[0]["action"] != "reviewed" {
		t.Fatalf("unexpected newest-first review rows: %#v", rows)
	}
}

func TestEditorialHistory_BylineAndRendering(t *testing.T) {
	writer, editor, review := "Author", "copy editor", "reviewer"
	post := &models.Post{Extra: map[string]interface{}{"editorial_history": []interface{}{editorialFixture()}}, AuthorObjects: []models.Author{
		{ID: "writer", Name: "Main Writer", Role: &writer},
		{ID: "editor", Name: "Main Editor", Role: &editor},
		{ID: "reviewer", Name: "Other Editor", Role: &review},
	}}
	mapped := postToMapUncached(post)
	writers, writersOK := mapped["byline_authors"].([]map[string]interface{})
	editors, editorsOK := mapped["byline_editors"].([]map[string]interface{})
	if !writersOK || !editorsOK || len(writers) != 1 || len(editors) != 2 {
		t.Fatalf("incorrect byline grouping: %#v", mapped)
	}
	for _, directory := range []string{"", "../../templates"} {
		engine, err := NewEngine(directory)
		if err != nil {
			t.Fatal(err)
		}
		ctx := NewContext(post, "", &models.Config{})
		byline, err := engine.Render("components/post_byline.html", ctx)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(byline, "Main Writer") || !strings.Contains(byline, "+1") || !strings.Contains(byline, "Other Editor") {
			t.Fatalf("missing compact credits: %s", byline)
		}
		history, err := engine.Render("components/editorial_history.html", ctx)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(history, "Open diff in GitHub") || !strings.Contains(history, "data-editorial-comparison") || !strings.Contains(history, "data-before=") || !strings.Contains(history, "data-after=") || strings.Contains(history, "javascript:") {
			t.Fatalf("unsafe or missing history: %s", history)
		}
	}
}

func TestEditorialHistory_PrimaryWriterKeepsProminence(t *testing.T) {
	role, contribution := "Author", "writer and editor"
	if isEditorialContributor(&models.Author{Role: &role, Contribution: &contribution}) {
		t.Fatal("explicit author was demoted because their description mentions editing")
	}
	if isEditorialContributor(&models.Author{Contributions: []string{"writing-original-draft", "writing-review-editing"}}) {
		t.Fatal("original writer was demoted because they also edited")
	}
}

func TestEditorialHistory_AssetsOnlyOnPostsWithHistory(t *testing.T) {
	for _, directory := range []string{"", "../../templates"} {
		engine, err := NewEngine(directory)
		if err != nil {
			t.Fatal(err)
		}
		for _, hasHistory := range []bool{false, true} {
			title := "A post"
			post := &models.Post{Title: &title, Slug: "a-post", Href: "/a-post/"}
			if hasHistory {
				post.Extra = map[string]interface{}{"editorial_history": []interface{}{editorialFixture()}}
			}
			html, err := engine.Render("post.html", NewContext(post, "<p>Body</p>", &models.Config{}))
			if err != nil {
				t.Fatal(err)
			}
			if strings.Contains(html, "editorial-history.js") != hasHistory || strings.Contains(html, "editorial-history.css") != hasHistory {
				t.Fatalf("unexpected history assets: directory=%q has_history=%v", directory, hasHistory)
			}
			if strings.Contains(html, `<script src="/js/editorial-diff.js"`) {
				t.Fatal("heavy viewer loaded eagerly")
			}
			if hasHistory && strings.Index(html, `class="editorial-history"`) > strings.Index(html, `class="post-content e-content`) {
				t.Fatalf("editorial history should be visible before the post body: directory=%q", directory)
			}
		}
	}
}

func TestEditorialHistory_NoScriptFallbackAndMissingAuthorObjects(t *testing.T) {
	for _, directory := range []string{"", "../../templates"} {
		engine, err := NewEngine(directory)
		if err != nil {
			t.Fatal(err)
		}
		// Imported records may have no resolved author objects; the history
		// entry point must still appear in this byline fallback.
		post := &models.Post{Extra: map[string]interface{}{
			"editorial_history": []interface{}{editorialFixture()},
		}}
		ctx := NewContext(post, "", &models.Config{})
		byline, err := engine.Render("components/post_byline.html", ctx)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(byline, "data-editorial-open") {
			t.Fatalf("history opener lost without resolved authors: directory=%q:\n%s", directory, byline)
		}
		history, err := engine.Render("components/editorial_history.html", ctx)
		if err != nil {
			t.Fatal(err)
		}
		start := strings.Index(history, `<details class="editorial-history__fallback"`)
		end := strings.LastIndex(history, "</details>")
		if start < 0 || end <= start {
			t.Fatalf("missing script-independent history: directory=%q:\n%s", directory, history)
		}
		fallback := history[start:end]
		if !strings.Contains(fallback, "data-editorial-fallback") || strings.Contains(history, "<noscript>") {
			t.Fatalf("fallback must work when JavaScript is enabled but fails: directory=%q", directory)
		}
		for _, expected := range []string{"<details", "editorial-history__fallback", "Fixed a typo.", "File at this revision", "Open diff in GitHub"} {
			if !strings.Contains(fallback, expected) {
				t.Fatalf("no-script fallback missing %q: directory=%q:\n%s", expected, directory, fallback)
			}
		}
		if strings.Contains(fallback, "javascript:") {
			t.Fatalf("unsafe fallback URL: directory=%q:\n%s", directory, fallback)
		}
	}
}
