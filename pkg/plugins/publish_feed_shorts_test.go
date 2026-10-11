package plugins

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata-go/pkg/lifecycle"
	"github.com/WaylonWalker/markata-go/pkg/models"
)

func TestShortsPostItemUsesBoundedDropperImageSizes(t *testing.T) {
	imageURL := "https://dropper.waylonwalker.com/file/photo.webp?token=abc&width=900&height=600&h=300"
	item := shortsPostItem(&models.Post{
		Slug:  "shots/photo",
		Extra: map[string]interface{}{"image": imageURL},
	})

	if item["src"] != "https://dropper.waylonwalker.com/file/photo.webp?token=abc&w=1280" {
		t.Fatalf("unexpected main image URL: %v", item["src"])
	}
	if item["thumb"] != "https://dropper.waylonwalker.com/file/photo.webp?token=abc&w=240" {
		t.Fatalf("unexpected preview URL: %v", item["thumb"])
	}
	if item["placeholder"] != "https://dropper.waylonwalker.com/file/photo.webp?token=abc&w=72" {
		t.Fatalf("unexpected placeholder URL: %v", item["placeholder"])
	}
	if item["poster"] != "" {
		t.Fatalf("image-only item should not get a video poster: %v", item["poster"])
	}
}

func TestShortsPostItemDoesNotResizeUntrustedOrVideoMedia(t *testing.T) {
	imageURL := "https://images.example.test/photo.webp?token=abc"
	image := shortsPostItem(&models.Post{Extra: map[string]interface{}{"image": imageURL}})
	if image["src"] != imageURL || image["thumb"] != "" || image["placeholder"] != "" {
		t.Fatalf("untrusted image should stay unchanged without speculative variants: %+v", image)
	}

	videoURL := "https://dropper.waylonwalker.com/file/clip.mp4?token=video"
	posterURL := "https://dropper.waylonwalker.com/file/poster.webp?token=poster"
	video := shortsPostItem(&models.Post{Extra: map[string]interface{}{
		"video":  videoURL,
		"poster": posterURL,
	}})
	if video["src"] != videoURL {
		t.Fatalf("video source must not receive image resize parameters: %v", video["src"])
	}
	if video["thumb"] != "https://dropper.waylonwalker.com/file/poster.webp?token=poster&w=240" {
		t.Fatalf("unexpected video preview poster: %v", video["thumb"])
	}
	if video["placeholder"] != "https://dropper.waylonwalker.com/file/poster.webp?token=poster&w=72" {
		t.Fatalf("unexpected video placeholder: %v", video["placeholder"])
	}
	if video["poster"] != "https://dropper.waylonwalker.com/file/poster.webp?token=poster&w=720" {
		t.Fatalf("unexpected video poster: %v", video["poster"])
	}
}

func TestShortsHTMLFirstFrameIsPublicAndStandalone(t *testing.T) {
	output := t.TempDir()
	cfg := lifecycle.NewConfig()
	cfg.OutputDir = output
	feed := &models.FeedConfig{
		Slug: "shots", Title: "Shots", ShortsPath: "/shorts/",
		Views: []string{models.FeedViewDefault, models.FeedViewShorts},
		Posts: []*models.Post{
			{Slug: "shots/private", Private: true, Published: true,
				Extra: map[string]interface{}{"image": "https://dropper.waylonwalker.com/file/private.webp"}},
			{Slug: "shots/first", Published: true,
				Extra: map[string]interface{}{"image": "https://dropper.waylonwalker.com/file/first.webp"}},
			{Slug: "shots/second", Published: true,
				Extra: map[string]interface{}{"image": "https://dropper.waylonwalker.com/file/second.webp"}},
		},
		Formats: models.FeedFormats{HTML: true},
	}
	plugin := NewPublishFeedsPlugin()
	if err := plugin.publishShortsPages(feed, cfg, output); err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(filepath.Join(output, "shorts", "index.html"))
	if err != nil {
		t.Fatal(err)
	}
	html := string(raw)
	for _, want := range []string{
		"<!doctype html>", "shorts-first-frame",
		`data-id="shots/first"`,
		"https://dropper.waylonwalker.com/file/first.webp?w=72",
		"https://dropper.waylonwalker.com/file/first.webp?w=1280",
		"shorts-first-frame-sharp",
		"feed-shorts.css", "feed-shorts.js",
	} {
		if !strings.Contains(html, want) {
			t.Errorf("Shorts HTML missing %q", want)
		}
	}
	for _, unwanted := range []string{
		"private.webp", "second.webp", "css/main.css", "css/fonts.css",
		"css/variables.css", "theme-bootstrap", "js/theme-bootstrap",
	} {
		if strings.Contains(html, unwanted) {
			t.Errorf("Shorts HTML unexpectedly contains %q", unwanted)
		}
	}
	// A deep-link target isn't known to the server. Keep the manifest intact
	// so the client can replace this first-frame preview after startup.
	raw, err = os.ReadFile(filepath.Join(output, "shorts", "data", "index.json"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(raw), "shots/second") || strings.Contains(string(raw), "shots/private") {
		t.Fatal("first-frame selection changed Shorts public manifest contents")
	}
}

func TestShortsHTMLFirstFrameVideoUsesPoster(t *testing.T) {
	output := t.TempDir()
	cfg := lifecycle.NewConfig()
	cfg.OutputDir = output
	feed := &models.FeedConfig{
		Slug: "shots", ShortsPath: "/shorts/",
		Views: []string{models.FeedViewDefault, models.FeedViewShorts},
		Posts: []*models.Post{
			{Slug: "shots/video", Published: true, Extra: map[string]interface{}{
				"video": "https://dropper.waylonwalker.com/file/movie.mp4",
				"poster": "https://dropper.waylonwalker.com/file/poster.webp",
			}},
		},
		Formats: models.FeedFormats{HTML: true},
	}
	if err := NewPublishFeedsPlugin().publishShortsPages(feed, cfg, output); err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(filepath.Join(output, "shorts", "index.html"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(raw), "poster.webp?w=72") || strings.Contains(string(raw), "movie.mp4") {
		t.Fatal("video-first Shorts must display a tiny poster, never preload video media")
	}
}

func TestShortsPublishesAllPagesAtRootRoute(t *testing.T) {
	output := t.TempDir()
	cfg := lifecycle.NewConfig()
	cfg.OutputDir = output
	posts := make([]*models.Post, 259)
	for i := range posts {
		slug := fmt.Sprintf("shots/p%03d", i)
		title := fmt.Sprintf("Post %d", i)
		posts[i] = &models.Post{
			Slug: slug, Href: "/" + slug + "/", Published: true, Title: &title,
			Extra: map[string]interface{}{"image": "/images/test.webp"},
		}
	}
	feed := &models.FeedConfig{
		Slug: "shots", Title: "Shots", ShortsPath: "/shorts/",
		Views: []string{models.FeedViewDefault, models.FeedViewShorts},
		Posts: posts, ItemsPerPage: 10,
		Formats: models.FeedFormats{HTML: true},
	}
	feed.Paginate("/shots")
	plugin := NewPublishFeedsPlugin()
	if err := plugin.publishFeed(feed, cfg, output); err != nil {
		t.Fatal(err)
	}
	indexPath := filepath.Join(output, "shorts", "data", "index.json")
	raw, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatal(err)
	}
	var index struct {
		Total     int      `json:"total"`
		ChunkSize int      `json:"chunk_size"`
		IDs       []string `json:"ids"`
	}
	if err := json.Unmarshal(raw, &index); err != nil {
		t.Fatal(err)
	}
	if index.Total != len(posts) || len(index.IDs) != len(posts) || index.ChunkSize != shortsChunkSize {
		t.Fatalf("incomplete Shorts manifest: %+v", index)
	}
	if index.IDs[258] != "shots/p258" {
		t.Fatalf("last item missing from Shorts manifest: %q", index.IDs[258])
	}
	for _, name := range []string{"0000.json", "0001.json", "0002.json"} {
		if _, err := os.Stat(filepath.Join(output, "shorts", "data", name)); err != nil {
			t.Fatalf("missing chunk %s: %v", name, err)
		}
	}
	if _, err := os.Stat(filepath.Join(output, "shorts", "index.html")); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(output, "shots", "index.html")); err != nil {
		t.Fatal(err)
	}
}

func TestShortsRouteCollisionAndValidation(t *testing.T) {
	feed := models.FeedConfig{
		Slug: "shots", ShortsPath: "/shorts/", Views: []string{models.FeedViewDefault, models.FeedViewShorts},
		Formats: models.FeedFormats{HTML: true},
	}
	if err := validateShortsRoutes([]models.FeedConfig{feed}, nil); err != nil {
		t.Fatal(err)
	}
	post := &models.Post{Slug: "shorts", Published: true}
	if err := validateShortsRoutes([]models.FeedConfig{feed}, []*models.Post{post}); err == nil {
		t.Fatal("shorts route collided with an existing post")
	}
	other := feed
	other.Slug = "other"
	if err := validateShortsRoutes([]models.FeedConfig{feed, other}, nil); err == nil {
		t.Fatal("two feeds claimed the same shorts route")
	}
	feed.ShortsPath = "/../outside/"
	if err := validateShortsRoutes([]models.FeedConfig{feed}, nil); err == nil {
		t.Fatal("unsafe shorts route accepted")
	}
}

func TestShortsIndexPrivacyAndCleanup(t *testing.T) {
	output := t.TempDir()
	cfg := lifecycle.NewConfig()
	cfg.OutputDir = output
	title := "Public"
	posts := []*models.Post{
		{Slug: "shots/public", Published: true, Title: &title, Extra: map[string]interface{}{"video": "https://example.com/movie.mp4", "poster": "/images/poster.webp"}},
		{Slug: "shots/secret", Published: true, Private: true, Title: &title, Extra: map[string]interface{}{"image": "/secret.png"}},
		{Slug: "shots/draft", Published: false, Draft: true},
	}
	feed := &models.FeedConfig{
		Slug: "shots", ShortsPath: "/shorts/", Views: []string{models.FeedViewDefault, models.FeedViewShorts},
		Posts: posts, Formats: models.FeedFormats{HTML: true},
	}
	plugin := NewPublishFeedsPlugin()
	if err := plugin.publishShortsPages(feed, cfg, output); err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(filepath.Join(output, "shorts", "data", "index.json"))
	if err != nil {
		t.Fatal(err)
	}
	if len(raw) == 0 || containsShortsPrivateData(string(raw)) {
		t.Fatal("private posts leaked into the Shorts manifest")
	}
	chunk, err := os.ReadFile(filepath.Join(output, "shorts", "data", "0000.json"))
	if err != nil {
		t.Fatal(err)
	}
	var entries []struct {
		ID     string `json:"id"`
		Kind   string `json:"kind"`
		Poster string `json:"poster"`
	}
	if err := json.Unmarshal(chunk, &entries); err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 || entries[0].Kind != "video" || entries[0].ID != "shots/public" {
		t.Fatalf("wrong safe selection: %+v", entries)
	}
	// A disabled root feed must never remove the Shots feed's /shorts/ route.
	if err := cleanupDisabledShorts(output, "/shorts/", ""); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(output, "shorts", "index.html")); err != nil {
		t.Fatalf("unrelated feed removed Shorts: %v", err)
	}
	feed.Views = []string{models.FeedViewDefault}
	if err := cleanupDisabledShorts(output, feed.ShortsURL(), feed.Slug); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(output, "shorts", "index.html")); !os.IsNotExist(err) {
		t.Fatalf("stale shorts page after disable: %v", err)
	}
}

func containsShortsPrivateData(s string) bool {
	return strings.Contains(s, "shots/secret") || strings.Contains(s, "shots/draft")
}
