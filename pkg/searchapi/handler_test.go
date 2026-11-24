package searchapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/search"
)

func TestHandler_Search(t *testing.T) {
	title1 := "Go Programming Guide"
	title2 := "Private Diary Entry"
	desc2 := "A private thought"
	date := time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)

	posts := []*models.Post{
		{
			Path:      "posts/go-guide.md",
			Title:     &title1,
			Content:   "Learn Go programming with examples and best practices.",
			Slug:      "go-guide",
			Href:      "/go-guide",
			Tags:      []string{"go", "programming"},
			Published: true,
			Date:      &date,
		},
		{
			Path:        "posts/diary.md",
			Title:       &title2,
			Description: &desc2,
			Content:     "Super secret encrypted content that should never be searchable.",
			Slug:        "diary",
			Href:        "/diary",
			Tags:        []string{"personal"},
			Published:   true,
			Private:     true,
			Date:        &date,
			Extra: map[string]interface{}{
				"_title_explicit":       true,
				"_description_explicit": true,
				"description":           descriptionValue(desc2),
				"cover_image":           "https://cdn.example.com/private-cover.webp",
				"thumbnail":             "https://cdn.example.com/private-thumb.webp",
				"word_count":            777,
			},
		},
	}

	cacheDir := t.TempDir()
	h := NewHandler(posts, cacheDir, DefaultConfig())
	defer h.Close()

	t.Run("basic search", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/search?q=go", http.NoBody)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
		}

		var resp SearchResponse
		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if resp.Total == 0 {
			t.Error("expected results for 'go'")
		}
		if resp.Query != "go" {
			t.Errorf("query = %q, want %q", resp.Query, "go")
		}
	})

	t.Run("missing query", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/search", http.NoBody)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})

	t.Run("method not allowed", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/search?q=test", http.NoBody)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("expected 405, got %d", w.Code)
		}
	})

	t.Run("private post content not searchable", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/search?q=searchable", http.NoBody)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		var resp SearchResponse
		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			t.Fatalf("decode: %v", err)
		}
		for _, r := range resp.Results {
			if r.Path == "posts/diary.md" {
				t.Error("private post should not be found via content-only terms")
			}
		}
	})

	t.Run("private post explicit title is searchable", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/search?q=diary", http.NoBody)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		var resp SearchResponse
		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			t.Fatalf("decode: %v", err)
		}
		found := false
		for _, r := range resp.Results {
			if r.Path != "posts/diary.md" {
				continue
			}
			found = true
			if !r.Private {
				t.Error("private flag should be set on result")
			}
			if r.Description != desc2 {
				t.Errorf("description = %q, want %q", r.Description, desc2)
			}
			if len(r.Tags) != 0 {
				t.Errorf("private search result should not expose tags, got %v", r.Tags)
			}
			if r.WordCount != 0 || r.ReadTime != "" {
				t.Errorf("private search result should not expose derived content stats, got word_count=%d read_time=%q", r.WordCount, r.ReadTime)
			}
			if r.MediaURL != "" || r.PosterURL != "" || r.VideoMIME != "" {
				t.Error("private search result should not expose media")
			}
		}
		if !found {
			t.Error("private post should be findable by explicit title")
		}
	})

	t.Run("fuzzy search", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/search?q=programing&fuzzy=true", http.NoBody)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		var resp SearchResponse
		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if !resp.Fuzzy {
			t.Error("expected fuzzy=true in response")
		}
	})

	t.Run("limit enforcement", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/search?q=go&limit=999", http.NoBody)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		var resp SearchResponse
		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if resp.Limit != 100 {
			t.Errorf("limit = %d, want 100 (max limit)", resp.Limit)
		}
	})

	t.Run("CORS preflight", func(t *testing.T) {
		req := httptest.NewRequest("OPTIONS", "/api/search", http.NoBody)
		req.Header.Set("Origin", "https://example.com")
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		if w.Code != http.StatusNoContent {
			t.Errorf("expected 204, got %d", w.Code)
		}
		if w.Header().Get("Access-Control-Allow-Origin") != "https://example.com" {
			t.Error("expected CORS origin header")
		}
	})
}

func descriptionValue(value string) string {
	return value
}

func TestReadOnlyHandler_Search(t *testing.T) {
	title := "Read **Only** Search"
	semanticTitle := "Read Only Search"
	description := "Index served without loading site content"
	date := time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC)

	posts := []*models.Post{
		{
			Path:             "posts/read-only.md",
			Title:            &title,
			TitleText:        semanticTitle,
			TitleTextDerived: true,
			Description:      &description,
			Content:          "This post validates read only bleve serving.",
			Slug:             "read-only",
			Href:             "/read-only",
			Tags:             []string{"search", "bleve"},
			Published:        true,
			Date:             &date,
		},
	}

	indexDir := filepath.Join(t.TempDir(), "search.bleve")
	idx, err := search.Build(indexDir, posts)
	if err != nil {
		t.Fatalf("build index: %v", err)
	}
	if err := idx.Close(); err != nil {
		t.Fatalf("close index: %v", err)
	}

	h := NewReadOnlyHandler(indexDir, DefaultConfig())
	defer h.Close()

	req := httptest.NewRequest("GET", "/api/search?q=read", http.NoBody)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp SearchResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Total == 0 {
		t.Fatal("expected read-only handler to return results")
	}
	if resp.Results[0].Title != semanticTitle {
		t.Fatalf("title = %q, want semantic title %q", resp.Results[0].Title, semanticTitle)
	}
	if resp.Results[0].Href != "/read-only" {
		t.Fatalf("href = %q, want %q", resp.Results[0].Href, "/read-only")
	}
}

func TestParseTags(t *testing.T) {
	got := parseTags(" go, programming , ,bleve ")
	want := []string{"go", "programming", "bleve"}
	if len(got) != len(want) {
		t.Fatalf("len(tags) = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("tags[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestHandler_UpdatePostsRebuildsOnContentChange(t *testing.T) {
	title := "Search Post"
	posts := []*models.Post{{
		Path:      "posts/search.md",
		Title:     &title,
		Content:   "original body",
		Slug:      "search",
		Href:      "/search",
		Published: true,
	}}

	h := NewHandler(posts, t.TempDir(), DefaultConfig())
	defer h.Close()

	req := httptest.NewRequest("GET", "/api/search?q=updated", http.NoBody)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	var before SearchResponse
	if err := json.NewDecoder(w.Body).Decode(&before); err != nil {
		t.Fatalf("decode before: %v", err)
	}
	if before.Total != 0 {
		t.Fatalf("expected no results before update, got %d", before.Total)
	}

	updated := *posts[0]
	updated.Content = "updated body"
	h.UpdatePosts([]*models.Post{&updated})

	req = httptest.NewRequest("GET", "/api/search?q=updated", http.NoBody)
	w = httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var after SearchResponse
	if err := json.NewDecoder(w.Body).Decode(&after); err != nil {
		t.Fatalf("decode after: %v", err)
	}
	if after.Total != 1 {
		t.Fatalf("expected 1 result after update, got %d", after.Total)
	}

	h.UpdatePosts(nil)
	req = httptest.NewRequest("GET", "/api/search?q=updated", http.NoBody)
	w = httptest.NewRecorder()
	h.ServeHTTP(w, req)
	var cleared SearchResponse
	if err := json.NewDecoder(w.Body).Decode(&cleared); err != nil {
		t.Fatalf("decode cleared response: %v", err)
	}
	if cleared.Total != 0 {
		t.Fatalf("expected no results after clearing posts, got %d", cleared.Total)
	}
}
