package plugins

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/WaylonWalker/markata/pkg/buildcache"
	"github.com/WaylonWalker/markata/pkg/imageindex"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

func TestParseImageLibraryConfig_DefaultsAndExplicitFalse(t *testing.T) {
	defaults, err := parseImageLibraryConfig(nil)
	if err != nil {
		t.Fatalf("parseImageLibraryConfig(nil) error = %v", err)
	}
	if !defaults.IsEnabled() || !defaults.ShouldExportJSON() || defaults.ShouldIncludeUnreferenced() {
		t.Fatalf("defaults = %#v", defaults)
	}
	normalized := normalizeImageLibraryConfig(models.ImagesConfig{})
	if normalized.ShouldIncludeUnreferenced() {
		t.Fatalf("normalized zero-value config includes unreferenced media: %#v", normalized)
	}

	configured, err := parseImageLibraryConfig(map[string]interface{}{
		"images": map[string]interface{}{
			"enabled":              false,
			"path":                 "media",
			"template":             "media.html",
			"export_json":          false,
			"include_unreferenced": false,
		},
	})
	if err != nil {
		t.Fatalf("parseImageLibraryConfig() error = %v", err)
	}
	if configured.IsEnabled() || configured.ShouldExportJSON() || configured.ShouldIncludeUnreferenced() || configured.Path != "media" || configured.Template != "media.html" {
		t.Fatalf("configured = %#v", configured)
	}
}

func TestImageLibrary_AdoptsOnlyUnchangedPublishedStageOutputs(t *testing.T) {
	for _, tc := range []struct {
		name, stageName, contents string
		wantAdopt                 bool
	}{
		{name: "published stage", stageName: ".output.markata-build-123", contents: "generated", wantAdopt: true},
		{name: "modified output", stageName: ".output.markata-build-123", contents: "modified"},
		{name: "unrelated root", stageName: "unrelated", contents: "generated"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			output := filepath.Join(root, "output")
			stage := filepath.Join(root, tc.stageName)
			path := filepath.Join(output, "images", "index.html")
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(path, []byte(tc.contents), 0o600); err != nil {
				t.Fatal(err)
			}
			cache := buildcache.New(filepath.Join(root, "cache"))
			cache.SetImageLibraryOutputs(root, stage, map[string]string{
				filepath.Join(stage, "images", "index.html"): buildcache.ContentHash("generated"),
			})
			adoptPublishedImageLibraryOutputs(cache, root, output)
			_, gotRoot, _ := cache.GetImageLibraryOutputs()
			if got := gotRoot == output; got != tc.wantAdopt {
				t.Fatalf("adopted = %v, want %v", got, tc.wantAdopt)
			}
		})
	}
}

func TestImageLibraryPage_UsesCanonicalDropperDerivatives(t *testing.T) {
	index := imageIndexForTest()
	page := newImageLibraryPage(index)
	if len(page.Images) != 2 {
		t.Fatalf("page images = %#v", page.Images)
	}

	remote := page.Images[1]
	if !strings.Contains(remote.PreviewSrc, "w=640") || !strings.Contains(remote.Srcset, "w=320") || !strings.Contains(remote.Srcset, "w=1280") {
		t.Fatalf("remote presentation URLs = preview %q, srcset %q", remote.PreviewSrc, remote.Srcset)
	}
	if remote.Src != "https://dropper.wayl.one/file/photo.webp?token=keep" {
		t.Fatalf("canonical source changed = %q", remote.Src)
	}
	if !strings.Contains(remote.Markdown, remote.Src) {
		t.Fatalf("copy Markdown does not use canonical source: %q", remote.Markdown)
	}
	if page.UsedCount != 1 || page.UnusedCount != 1 || page.CoverCount != 1 {
		t.Fatalf("page counts = %#v", page)
	}
}

func TestImageLibraryPage_RendersVideoMetadataAndEmbedLabel(t *testing.T) {
	index := imageindex.Index{Images: []imageindex.Image{{
		Src:       "https://dropper.wayl.one/file/clip.mp4?token=keep",
		MIMEType:  "video/mp4",
		PosterSrc: "http://dropper.wayl.one/file/clip.webp",
		Embed:     true,
		Uses:      []imageindex.Use{{Href: "/clip/", Title: "Clip post", Caption: "A useful clip", Embed: true}},
	}}}

	page := newImageLibraryPage(index)
	if len(page.Images) != 1 {
		t.Fatalf("page images = %#v", page.Images)
	}
	card := page.Images[0]
	if !card.IsVideo || !card.Embed {
		t.Fatalf("video card metadata = %#v", card)
	}
	if !strings.Contains(card.PosterSrc, "clip.webp") || !strings.Contains(card.PosterSrc, "w=640") {
		t.Fatalf("video poster = %q", card.PosterSrc)
	}
	if card.Srcset != "" {
		t.Fatalf("video srcset = %q, want empty", card.Srcset)
	}
	if !strings.Contains(card.SearchText, "video") || !strings.Contains(card.SearchText, "embed") {
		t.Fatalf("video search text = %q", card.SearchText)
	}
	if !strings.Contains(card.SearchText, "a useful clip") || !strings.Contains(card.SearchText, "clip post") {
		t.Fatalf("video search text does not include caption and title = %q", card.SearchText)
	}
}

func TestImageLibraryPage_FormatsAddedAtInUTC(t *testing.T) {
	addedAt := time.Date(2026, time.January, 2, 23, 30, 0, 0, time.FixedZone("west", -8*60*60))
	page := newImageLibraryPage(imageindex.Index{Images: []imageindex.Image{{Src: "/images/photo.png", AddedAt: &addedAt}}})
	if got := page.Images[0].AddedAt; got != "Jan 3, 2026" {
		t.Fatalf("AddedAt = %q, want UTC date", got)
	}
}

func TestImageLibraryPage_FormatsLastUsedAtInUTC(t *testing.T) {
	lastUsedAt := time.Date(2026, time.January, 2, 23, 30, 0, 0, time.FixedZone("west", -8*60*60))
	page := newImageLibraryPage(imageindex.Index{Images: []imageindex.Image{{Src: "/images/photo.png", LastUsedAt: &lastUsedAt}}})
	if got := page.Images[0].LastUsedAt; got != "Jan 3, 2026" {
		t.Fatalf("LastUsedAt = %q, want UTC date", got)
	}
	if page.Images[0].LastUsedAtUnix != lastUsedAt.Unix() {
		t.Fatalf("LastUsedAtUnix = %d, want %d", page.Images[0].LastUsedAtUnix, lastUsedAt.Unix())
	}
}

func TestImageLibraryRecentlyAddedFilterUsesTimestampWindow(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("..", "themes", "default", "static", "js", "image-library.js"))
	if err != nil {
		t.Fatalf("read image-library.js: %v", err)
	}
	source := string(data)
	for _, marker := range []string{
		"var recentWindowSeconds = 30 * 24 * 60 * 60;",
		"var addedAt = Number(card.getAttribute(\"data-added\") || 0);",
		"if (!(addedAt > 0)) return false;",
		"var now = Math.floor(Date.now() / 1000);",
		"addedAt >= now - recentWindowSeconds && addedAt <= now;",
	} {
		if !strings.Contains(source, marker) {
			t.Errorf("image-library.js missing recent-filter marker %q", marker)
		}
	}
	if strings.Contains(source, "return !!card.getAttribute(\"data-added\")") {
		t.Fatal("recent filter tests the data-added string instead of its timestamp")
	}
}

func TestImageLibraryPlugin_WriteOutputsPageAndJSON(t *testing.T) {
	root := t.TempDir()
	output := filepath.Join(root, "output")
	assets := filepath.Join(root, "static")
	if err := os.MkdirAll(assets, 0o755); err != nil {
		t.Fatal(err)
	}
	post := models.NewPost("posts/hello.md")
	post.Slug = "hello"
	post.Href = "/hello/"
	post.Published = true
	post.Title = stringPointerImages("Hello")
	post.Content = "![Hello image](https://dropper.wayl.one/file/hello.webp)\nHello figure caption"
	videoPost := models.NewPost("posts/video.md")
	videoPost.Slug = "video"
	videoPost.Href = "/video/"
	videoPost.Published = true
	videoPost.Title = stringPointerImages("Video")
	videoDate := time.Date(2026, time.January, 10, 12, 0, 0, 0, time.UTC)
	videoPost.Date = &videoDate
	videoPost.Extra["image"] = "https://dropper.wayl.one/file/video.mp4"
	videoPost.Extra["poster_image"] = "https://dropper.wayl.one/file/video.webp"

	modelsConfig := models.NewConfig()
	modelsConfig.OutputDir = output
	modelsConfig.AssetsDir = assets
	modelsConfig.Images = models.NewImagesConfig()
	manager := lifecycle.NewManager()
	manager.SetConfig(&lifecycle.Config{
		ContentDir: root,
		OutputDir:  output,
		Extra: map[string]interface{}{
			"assets_dir":      assets,
			"markata_version": "test",
			"models_config":   modelsConfig,
		},
	})
	manager.SetPosts([]*models.Post{post, videoPost})
	manager.RegisterPlugin(NewImageLibraryPlugin())

	if err := manager.RunTo(lifecycle.StageWrite); err != nil {
		t.Fatalf("RunTo(write) error = %v", err)
	}

	pagePath := filepath.Join(output, "images", "index.html")
	jsonPath := filepath.Join(output, "images", "index.json")
	flatJSONPath := filepath.Join(output, "images.json")
	page, err := os.ReadFile(pagePath)
	if err != nil {
		t.Fatalf("read generated page: %v", err)
	}
	if !strings.Contains(string(page), "Copy Markdown") || !strings.Contains(string(page), "hello.webp") || !strings.Contains(string(page), "Hello figure caption") || !strings.Contains(string(page), "data-image-filter=\"used\"") || !strings.Contains(string(page), `data-added="0"`) {
		t.Fatalf("generated page missing image-library controls/content: %s", page)
	}
	if !strings.Contains(string(page), "<video") || !strings.Contains(string(page), `preload="none"`) || !strings.Contains(string(page), `poster="https://dropper.wayl.one/file/video.webp?w=640"`) || !strings.Contains(string(page), `data-video-src="https://dropper.wayl.one/file/video.mp4?w=640"`) || !strings.Contains(string(page), `data-video-type="video/mp4"`) || !strings.Contains(string(page), ">Video</span>") || !strings.Contains(string(page), `data-added="1768046400"`) || !strings.Contains(string(page), `data-last-used="1768046400"`) || !strings.Contains(string(page), ">Latest used</option>") || strings.Contains(string(page), `<source src="`) || strings.Contains(string(page), `<img src="https://dropper.wayl.one/file/video.mp4`) {
		t.Fatalf("generated page missing video presentation: %s", page)
	}
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		t.Fatalf("read generated JSON: %v", err)
	}
	var artifact map[string]interface{}
	if err := json.Unmarshal(data, &artifact); err != nil {
		t.Fatalf("generated JSON is invalid: %v", err)
	}
	if artifact["image_count"] != float64(2) {
		t.Fatalf("generated image_count = %#v", artifact["image_count"])
	}
	if !strings.Contains(string(data), `"caption":"Hello figure caption"`) {
		t.Fatalf("generated JSON missing figure caption: %s", data)
	}
	flatData, err := os.ReadFile(flatJSONPath)
	if err != nil {
		t.Fatalf("read root image index: %v", err)
	}
	if !bytes.Equal(flatData, data) {
		t.Fatalf("root image index differs from nested artifact")
	}
}

func TestImageLibraryPlugin_PrivateBodyDoesNotChangePublicArtifactsOrHash(t *testing.T) {
	root := t.TempDir()
	output := filepath.Join(root, "output")
	assets := filepath.Join(root, "static")
	if err := os.MkdirAll(assets, 0o755); err != nil {
		t.Fatal(err)
	}
	privateMediaPath := filepath.Join(root, "secret-private.png")
	if err := os.WriteFile(privateMediaPath, []byte("private media"), 0o600); err != nil {
		t.Fatal(err)
	}

	public := models.NewPost("posts/public.md")
	public.Path = "posts/public.md"
	public.Slug = "public"
	public.Href = "/public/"
	public.Published = true
	public.Title = stringPointerImages("Public")
	public.Content = "![Public](https://example.test/public.png)"

	private := models.NewPost("posts/private.md")
	private.Path = "posts/private.md"
	private.Slug = "private"
	private.Href = "/private/"
	private.Published = true
	private.Private = true
	private.Content = "secret before"
	private.ArticleHTML = "<p>secret before</p>"

	modelsConfig := models.NewConfig()
	modelsConfig.OutputDir = output
	modelsConfig.AssetsDir = assets
	manager := lifecycle.NewManager()
	manager.SetConfig(&lifecycle.Config{
		ContentDir: root,
		OutputDir:  output,
		Extra: map[string]interface{}{
			"assets_dir":    assets,
			"models_config": modelsConfig,
		},
	})
	manager.SetPosts([]*models.Post{private, public})
	cache := buildcache.New(filepath.Join(root, "cache"))
	manager.Cache().Set("build_cache", cache)
	plugin := NewImageLibraryPlugin()
	manager.RegisterPlugin(plugin)

	if err := manager.RunTo(lifecycle.StageWrite); err != nil {
		t.Fatalf("first RunTo(write) error = %v", err)
	}
	paths := []string{
		filepath.Join(output, "images", "index.html"),
		filepath.Join(output, "images", "index.json"),
		filepath.Join(output, "images.json"),
	}
	before := make(map[string][]byte, len(paths))
	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read first artifact %s: %v", path, err)
		}
		before[path] = data
	}
	beforeHash := cache.GetImageLibraryHash()

	private.Content = `![classified markdown](/secret-private.png)

![classified markdown](/secret-a.png)

<figure><img src="/secret-b.png" alt="classified alt"><figcaption>classified caption</figcaption></figure>

![remote](https://cdn.example.test/private.jpg?token=SECRET)`
	private.ArticleHTML = `<figure><img src="/secret-c.png" alt="rendered secret"><figcaption>rendered caption</figcaption></figure>`
	if err := plugin.Write(manager); err != nil {
		t.Fatalf("second Write() error = %v", err)
	}
	if afterHash := cache.GetImageLibraryHash(); afterHash != beforeHash {
		t.Fatalf("private body changed image-library hash: before=%q after=%q", beforeHash, afterHash)
	}
	for _, path := range paths {
		after, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read second artifact %s: %v", path, err)
		}
		if !bytes.Equal(after, before[path]) {
			t.Fatalf("private body changed public artifact %s", path)
		}
	}

	if err := os.WriteFile(privateMediaPath, []byte("changed private media"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := plugin.Write(manager); err != nil {
		t.Fatalf("private-only media Write() error = %v", err)
	}
	if afterHash := cache.GetImageLibraryHash(); afterHash != beforeHash {
		t.Fatalf("private-only content media changed image-library hash: before=%q after=%q", beforeHash, afterHash)
	}
	for _, path := range paths {
		after, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read private-only media artifact %s: %v", path, err)
		}
		if !bytes.Equal(after, before[path]) {
			t.Fatalf("private-only media changed public artifact %s", path)
		}
	}
}

func TestImageLibraryPlugin_PrivateSafeCoverFrontmatterChangesInventory(t *testing.T) {
	root := t.TempDir()
	output := filepath.Join(root, "output")
	assets := filepath.Join(root, "static")
	if err := os.MkdirAll(assets, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(assets, "first.png"), []byte("first"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(assets, "second.png"), []byte("second"), 0o600); err != nil {
		t.Fatal(err)
	}

	private := models.NewPost("posts/private.md")
	private.Path = "posts/private.md"
	private.Slug = "private"
	private.Href = "/private/"
	private.Published = true
	private.Private = true
	private.Content = "private body"
	private.Extra["cover"] = "/first.png"

	modelsConfig := models.NewConfig()
	modelsConfig.OutputDir = output
	modelsConfig.AssetsDir = assets
	manager := lifecycle.NewManager()
	manager.SetConfig(&lifecycle.Config{
		ContentDir: root,
		OutputDir:  output,
		Extra: map[string]interface{}{
			"assets_dir":    assets,
			"models_config": modelsConfig,
		},
	})
	manager.SetPosts([]*models.Post{private})
	cache := buildcache.New(filepath.Join(root, "cache"))
	manager.Cache().Set("build_cache", cache)
	plugin := NewImageLibraryPlugin()
	if err := plugin.Write(manager); err != nil {
		t.Fatalf("first Write() error = %v", err)
	}
	firstHash := cache.GetImageLibraryHash()
	firstJSON, err := os.ReadFile(filepath.Join(output, "images.json"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(firstJSON, []byte(`"src":"/first.png"`)) || !bytes.Contains(firstJSON, []byte(`"uses":[]`)) {
		t.Fatalf("first private safe-cover JSON = %s", firstJSON)
	}

	private.Extra["cover"] = "/second.png"
	if err := plugin.Write(manager); err != nil {
		t.Fatalf("second Write() error = %v", err)
	}
	secondHash := cache.GetImageLibraryHash()
	secondJSON, err := os.ReadFile(filepath.Join(output, "images.json"))
	if err != nil {
		t.Fatal(err)
	}
	if secondHash == firstHash || bytes.Equal(firstJSON, secondJSON) || !bytes.Contains(secondJSON, []byte(`"src":"/second.png"`)) || bytes.Contains(secondJSON, []byte(`"src":"/first.png"`)) {
		t.Fatalf("private safe-cover change was not reflected: first=%s second=%s hashes=%q/%q", firstJSON, secondJSON, firstHash, secondHash)
	}
}

func TestImageLibraryPlugin_FullWriteExcludesOutputMediaFromInputs(t *testing.T) {
	root := t.TempDir()
	output := filepath.Join(root, "output")
	assets := filepath.Join(root, "static")
	if err := os.MkdirAll(assets, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(assets, "photo.webp"), []byte("source media"), 0o600); err != nil {
		t.Fatal(err)
	}

	modelsConfig := models.NewConfig()
	modelsConfig.OutputDir = output
	modelsConfig.AssetsDir = assets
	includeUnreferenced := true
	modelsConfig.Images.IncludeUnreferenced = &includeUnreferenced
	manager := lifecycle.NewManager()
	manager.SetConfig(&lifecycle.Config{
		ContentDir: root,
		OutputDir:  output,
		Extra: map[string]interface{}{
			"assets_dir":    assets,
			"models_config": modelsConfig,
		},
	})
	cache := buildcache.New(filepath.Join(root, "cache"))
	manager.Cache().Set("build_cache", cache)
	manager.RegisterPlugins(NewStaticAssetsPlugin(), NewImageLibraryPlugin())

	if err := manager.RunTo(lifecycle.StageWrite); err != nil {
		t.Fatalf("full lifecycle write error = %v", err)
	}
	if _, err := os.Stat(filepath.Join(output, "photo.webp")); err != nil {
		t.Fatalf("static asset was not copied before image-library write: %v", err)
	}
	outputPrefix := filepath.ToSlash(filepath.Clean(output)) + "/"
	for path := range cache.GetImageLibraryMedia() {
		if strings.HasPrefix(path, outputPrefix) || path == strings.TrimSuffix(outputPrefix, "/") {
			t.Fatalf("generated output participated in image-library media inputs: %q", path)
		}
	}

	firstHash := cache.GetImageLibraryHash()
	if err := os.WriteFile(filepath.Join(output, "generated.webp"), []byte("generated output"), 0o600); err != nil {
		t.Fatal(err)
	}
	plugin := NewImageLibraryPlugin()
	if err := plugin.Write(manager); err != nil {
		t.Fatalf("second image-library write error = %v", err)
	}
	if secondHash := cache.GetImageLibraryHash(); secondHash != firstHash {
		t.Fatalf("output-only media changed image-library hash: before=%q after=%q", firstHash, secondHash)
	}
}

func TestImageLibraryPlugin_ExportJSONFalseRemovesStaleArtifact(t *testing.T) {
	root := t.TempDir()
	output := filepath.Join(root, "output")
	assets := filepath.Join(root, "static")
	if err := os.MkdirAll(filepath.Join(output, "images"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(output, "images", "index.json"), []byte("stale"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(output, "images.json"), []byte("stale-root"), 0o600); err != nil {
		t.Fatal(err)
	}

	modelsConfig := models.NewConfig()
	modelsConfig.OutputDir = output
	modelsConfig.AssetsDir = assets
	modelsConfig.Images = models.NewImagesConfig()
	noJSON := false
	modelsConfig.Images.ExportJSON = &noJSON
	manager := lifecycle.NewManager()
	manager.SetConfig(&lifecycle.Config{
		ContentDir: root,
		OutputDir:  output,
		Extra: map[string]interface{}{
			"assets_dir":    assets,
			"models_config": modelsConfig,
		},
	})
	cache := buildcache.New(filepath.Join(root, "cache"))
	cache.SetImageLibraryOutputs(root, output, map[string]string{
		filepath.Join(output, "images", "index.html"): "stale-page",
		filepath.Join(output, "images", "index.json"): buildcache.ContentHash("stale"),
		filepath.Join(output, "images.json"):          buildcache.ContentHash("stale-root"),
	})
	// The stale JSON is the only file that exists. The page hash is retained to
	// verify that missing plugin-owned files do not block cleanup.
	manager.Cache().Set("build_cache", cache)
	manager.RegisterPlugin(NewImageLibraryPlugin())

	if err := manager.RunTo(lifecycle.StageWrite); err != nil {
		t.Fatalf("RunTo(write) error = %v", err)
	}
	if _, err := os.Stat(filepath.Join(output, "images", "index.json")); !os.IsNotExist(err) {
		t.Fatalf("stale JSON artifact still exists, stat error = %v", err)
	}
	if _, err := os.Stat(filepath.Join(output, "images.json")); !os.IsNotExist(err) {
		t.Fatalf("stale root JSON artifact still exists, stat error = %v", err)
	}
	page, err := os.ReadFile(filepath.Join(output, "images", "index.html"))
	if err != nil {
		t.Fatalf("read generated page: %v", err)
	}
	if strings.Contains(string(page), "Download JSON index") {
		t.Fatal("page links to JSON after export_json=false")
	}
}

func TestImageLibraryPlugin_ExportJSONFalseReportsUnownedArtifact(t *testing.T) {
	root := t.TempDir()
	output := filepath.Join(root, "output")
	assets := filepath.Join(root, "static")
	jsonPath := filepath.Join(output, "images", "index.json")
	if err := os.MkdirAll(filepath.Dir(jsonPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(jsonPath, []byte("site-owned"), 0o600); err != nil {
		t.Fatal(err)
	}

	modelsConfig := models.NewConfig()
	modelsConfig.OutputDir = output
	modelsConfig.AssetsDir = assets
	noJSON := false
	modelsConfig.Images.ExportJSON = &noJSON
	manager := lifecycle.NewManager()
	manager.SetConfig(&lifecycle.Config{
		ContentDir: root,
		OutputDir:  output,
		Extra: map[string]interface{}{
			"assets_dir":    assets,
			"models_config": modelsConfig,
		},
	})

	err := NewImageLibraryPlugin().Write(manager)
	if err == nil || !strings.Contains(err.Error(), "not owned") {
		t.Fatalf("Write() error = %v, want unowned-output conflict", err)
	}
	if data, readErr := os.ReadFile(jsonPath); readErr != nil || string(data) != "site-owned" {
		t.Fatalf("unowned JSON changed: data=%q error=%v", data, readErr)
	}
}

func TestImageLibraryPlugin_RejectsUnownedRootJSONArtifact(t *testing.T) {
	root := t.TempDir()
	output := filepath.Join(root, "output")
	assets := filepath.Join(root, "static")
	flatJSONPath := filepath.Join(output, "images.json")
	if err := os.MkdirAll(output, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(flatJSONPath, []byte("site-owned"), 0o600); err != nil {
		t.Fatal(err)
	}

	modelsConfig := models.NewConfig()
	modelsConfig.OutputDir = output
	modelsConfig.AssetsDir = assets
	manager := lifecycle.NewManager()
	manager.SetConfig(&lifecycle.Config{
		ContentDir: root,
		OutputDir:  output,
		Extra: map[string]interface{}{
			"assets_dir":    assets,
			"models_config": modelsConfig,
		},
	})

	err := NewImageLibraryPlugin().Write(manager)
	if err == nil || !strings.Contains(err.Error(), "not owned") {
		t.Fatalf("Write() error = %v, want unowned root output conflict", err)
	}
	if data, readErr := os.ReadFile(flatJSONPath); readErr != nil || string(data) != "site-owned" {
		t.Fatalf("unowned root JSON changed: data=%q error=%v", data, readErr)
	}
}

func TestImageLibraryPlugin_RejectsUserReplacementOfCachedOutput(t *testing.T) {
	root := t.TempDir()
	output := filepath.Join(root, "output")
	assets := filepath.Join(root, "static")
	if err := os.MkdirAll(assets, 0o755); err != nil {
		t.Fatal(err)
	}

	modelsConfig := models.NewConfig()
	modelsConfig.OutputDir = output
	modelsConfig.AssetsDir = assets
	manager := lifecycle.NewManager()
	manager.SetConfig(&lifecycle.Config{
		ContentDir: root,
		OutputDir:  output,
		Extra: map[string]interface{}{
			"assets_dir":      assets,
			"markata_version": "test",
			"models_config":   modelsConfig,
		},
	})
	manager.Cache().Set("build_cache", buildcache.New(filepath.Join(root, "cache")))
	plugin := NewImageLibraryPlugin()
	manager.RegisterPlugin(plugin)

	if err := manager.RunTo(lifecycle.StageWrite); err != nil {
		t.Fatalf("first RunTo(write) error = %v", err)
	}
	pagePath := filepath.Join(output, "images", "index.html")
	before, err := os.Stat(pagePath)
	if err != nil {
		t.Fatal(err)
	}
	if err := plugin.Write(manager); err != nil {
		t.Fatalf("cached Write() error = %v", err)
	}
	after, err := os.Stat(pagePath)
	if err != nil {
		t.Fatal(err)
	}
	if !after.ModTime().Equal(before.ModTime()) {
		t.Fatal("cached image-library page was rewritten")
	}
	if err := os.WriteFile(pagePath, []byte("sentinel"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := plugin.Write(manager); err == nil || !strings.Contains(err.Error(), "was modified") {
		t.Fatalf("second Write() error = %v, want modified-output error", err)
	}
	page, err := os.ReadFile(pagePath)
	if err != nil {
		t.Fatal(err)
	}
	if string(page) != "sentinel" {
		t.Fatalf("cached page was rewritten: %q", page)
	}
}

func TestImageLibraryPlugin_RebuildsWhenTemplateChanges(t *testing.T) {
	root := t.TempDir()
	output := filepath.Join(root, "output")
	assets := filepath.Join(root, "static")
	templatesDir := filepath.Join(root, "templates")
	if err := os.MkdirAll(assets, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(templatesDir, 0o755); err != nil {
		t.Fatal(err)
	}
	templatePath := filepath.Join(templatesDir, "images.html")
	if err := os.WriteFile(templatePath, []byte("first"), 0o600); err != nil {
		t.Fatal(err)
	}

	modelsConfig := models.NewConfig()
	modelsConfig.OutputDir = output
	modelsConfig.AssetsDir = assets
	manager := lifecycle.NewManager()
	manager.SetConfig(&lifecycle.Config{
		ContentDir: root,
		OutputDir:  output,
		Extra: map[string]interface{}{
			"assets_dir":    assets,
			"models_config": modelsConfig,
			"templates_dir": templatesDir,
		},
	})
	manager.Cache().Set("build_cache", buildcache.New(filepath.Join(root, "cache")))
	plugin := NewImageLibraryPlugin()
	manager.RegisterPlugin(plugin)

	if err := manager.RunTo(lifecycle.StageWrite); err != nil {
		t.Fatalf("first RunTo(write) error = %v", err)
	}
	pagePath := filepath.Join(output, "images", "index.html")
	first, err := os.ReadFile(pagePath)
	if err != nil {
		t.Fatal(err)
	}
	if string(first) != "first" {
		t.Fatalf("first page = %q", first)
	}

	if err := os.WriteFile(templatePath, []byte("second"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := plugin.Write(manager); err != nil {
		t.Fatalf("second Write() error = %v", err)
	}
	second, err := os.ReadFile(pagePath)
	if err != nil {
		t.Fatal(err)
	}
	if string(second) != "second" {
		t.Fatalf("template change was not rendered: %q", second)
	}
}

func TestImageLibraryPlugin_RebuildsWhenRenderConfigChanges(t *testing.T) {
	root := t.TempDir()
	output := filepath.Join(root, "output")
	assets := filepath.Join(root, "static")
	templatesDir := filepath.Join(root, "templates")
	if err := os.MkdirAll(assets, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(templatesDir, 0o755); err != nil {
		t.Fatal(err)
	}
	templatePath := filepath.Join(templatesDir, "images.html")
	if err := os.WriteFile(templatePath, []byte("{{ site_title }}"), 0o600); err != nil {
		t.Fatal(err)
	}

	modelsConfig := models.NewConfig()
	modelsConfig.OutputDir = output
	modelsConfig.AssetsDir = assets
	modelsConfig.Title = "First title"
	manager := lifecycle.NewManager()
	manager.SetConfig(&lifecycle.Config{
		ContentDir: root,
		OutputDir:  output,
		Extra: map[string]interface{}{
			"assets_dir":    assets,
			"models_config": modelsConfig,
			"templates_dir": templatesDir,
		},
	})
	manager.Cache().Set("build_cache", buildcache.New(filepath.Join(root, "cache")))
	plugin := NewImageLibraryPlugin()
	manager.RegisterPlugin(plugin)

	if err := manager.RunTo(lifecycle.StageWrite); err != nil {
		t.Fatalf("first RunTo(write) error = %v", err)
	}
	pagePath := filepath.Join(output, "images", "index.html")
	first, err := os.ReadFile(pagePath)
	if err != nil {
		t.Fatal(err)
	}
	if string(first) != "First title" {
		t.Fatalf("first page = %q", first)
	}

	modelsConfig.Title = "Second title"
	if err := plugin.Write(manager); err != nil {
		t.Fatalf("second Write() error = %v", err)
	}
	second, err := os.ReadFile(pagePath)
	if err != nil {
		t.Fatal(err)
	}
	if string(second) != "Second title" {
		t.Fatalf("configuration change was not rendered: %q", second)
	}
}

func TestImageLibraryPlugin_RebuildsWhenReferencedContentImageChanges(t *testing.T) {
	root := t.TempDir()
	output := filepath.Join(root, "output")
	assets := filepath.Join(root, "static")
	postsDir := filepath.Join(root, "posts")
	if err := os.MkdirAll(postsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(assets, 0o755); err != nil {
		t.Fatal(err)
	}
	imagePath := filepath.Join(postsDir, "local.png")
	if err := os.WriteFile(imagePath, []byte("first image"), 0o600); err != nil {
		t.Fatal(err)
	}

	modelsConfig := models.NewConfig()
	modelsConfig.OutputDir = output
	modelsConfig.AssetsDir = assets
	post := models.NewPost("posts/example.md")
	post.Slug = "example"
	post.Href = "/example/"
	post.Published = true
	post.Content = "![Local image](local.png)"
	manager := lifecycle.NewManager()
	manager.SetConfig(&lifecycle.Config{
		ContentDir: root,
		OutputDir:  output,
		Extra: map[string]interface{}{
			"assets_dir":    assets,
			"models_config": modelsConfig,
		},
	})
	manager.SetPosts([]*models.Post{post})
	cache := buildcache.New(filepath.Join(root, "cache"))
	manager.Cache().Set("build_cache", cache)
	plugin := NewImageLibraryPlugin()

	if err := plugin.Write(manager); err != nil {
		t.Fatalf("first Write() error = %v", err)
	}
	firstHash := cache.GetImageLibraryHash()
	if firstHash == "" {
		t.Fatal("first image-library hash is empty")
	}

	if err := os.WriteFile(imagePath, []byte("second image"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := plugin.Write(manager); err != nil {
		t.Fatalf("second Write() error = %v", err)
	}
	if secondHash := cache.GetImageLibraryHash(); secondHash == firstHash {
		t.Fatalf("image-library hash did not change after referenced content image changed: %q", secondHash)
	}
}

func TestHashImageLibrarySources_ReusesReferencedExternalAsset(t *testing.T) {
	root := t.TempDir()
	contentDir := filepath.Join(root, "site", "content")
	assetsDir := filepath.Join(root, "shared", "static")
	outputDir := filepath.Join(root, "site", "output")
	if err := os.MkdirAll(contentDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(assetsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	assetPath := filepath.Join(assetsDir, "photo.png")
	writeImageLibraryPNG(t, assetPath, 3, 2)

	post := models.NewPost("posts/photo.md")
	post.Path = "posts/photo.md"
	post.Slug = "photo"
	post.Href = "/photo/"
	post.Published = true
	post.Content = "![Photo](/photo.png)"
	paths := imageindex.LocalMediaPaths([]*models.Post{post}, imageindex.BuildOptions{
		ContentDir: contentDir,
		AssetsDir:  assetsDir,
		OutputDir:  outputDir,
	})
	if len(paths) != 1 || filepath.Clean(paths[0]) != filepath.Clean(assetPath) {
		t.Fatalf("referenced external media paths = %#v, want %q", paths, assetPath)
	}

	cache := buildcache.New(filepath.Join(root, "cache"))
	reads := 0
	hashFile := func(path string) (string, error) {
		reads++
		return buildcache.HashFile(path)
	}
	sourceFiles := append(append([]string(nil), paths...), paths...)
	firstContentHash, firstStateHash, err := hashImageLibrarySourcesWithHasher(contentDir, assetsDir, outputDir, false, sourceFiles, imageHashExtensions(), cache, hashFile)
	if err != nil {
		t.Fatalf("first external asset hash error = %v", err)
	}
	if reads != 1 {
		t.Fatalf("first external asset read count = %d, want one deduplicated read", reads)
	}

	secondContentHash, secondStateHash, err := hashImageLibrarySourcesWithHasher(contentDir, assetsDir, outputDir, false, sourceFiles, imageHashExtensions(), cache, hashFile)
	if err != nil {
		t.Fatalf("second external asset hash error = %v", err)
	}
	if reads != 1 {
		t.Fatalf("unchanged external asset read count = %d, want 1", reads)
	}
	if firstContentHash != secondContentHash || firstStateHash != secondStateHash {
		t.Fatalf("unchanged external asset hashes differ: first=(%q, %q), second=(%q, %q)", firstContentHash, firstStateHash, secondContentHash, secondStateHash)
	}

	writeImageLibraryPNG(t, assetPath, 5, 4)
	thirdContentHash, _, err := hashImageLibrarySourcesWithHasher(contentDir, assetsDir, outputDir, false, sourceFiles, imageHashExtensions(), cache, hashFile)
	if err != nil {
		t.Fatalf("changed external asset hash error = %v", err)
	}
	if reads != 2 {
		t.Fatalf("changed external asset read count = %d, want changed file reread only", reads)
	}
	if thirdContentHash == firstContentHash {
		t.Fatal("changed external asset retained the cached content hash")
	}
}

func TestImageLibraryPlugin_RebuildsWhenReferencedExternalAssetChanges(t *testing.T) {
	root := t.TempDir()
	contentDir := filepath.Join(root, "site", "content")
	assetsDir := filepath.Join(root, "shared", "static")
	outputDir := filepath.Join(root, "site", "output")
	if err := os.MkdirAll(contentDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(assetsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	assetPath := filepath.Join(assetsDir, "photo.png")
	writeImageLibraryPNG(t, assetPath, 3, 2)

	post := models.NewPost("posts/photo.md")
	post.Path = "posts/photo.md"
	post.Slug = "photo"
	post.Href = "/photo/"
	post.Published = true
	post.Content = "![Photo](/photo.png)"

	modelsConfig := models.NewConfig()
	modelsConfig.OutputDir = outputDir
	modelsConfig.AssetsDir = assetsDir
	// Keep this explicit in the regression fixture: the referenced file must be
	// hashed even when the asset directory is not enumerated wholesale.
	includeUnreferenced := false
	modelsConfig.Images.IncludeUnreferenced = &includeUnreferenced
	manager := lifecycle.NewManager()
	manager.SetConfig(&lifecycle.Config{
		ContentDir: contentDir,
		OutputDir:  outputDir,
		Extra: map[string]interface{}{
			"assets_dir":    assetsDir,
			"models_config": modelsConfig,
		},
	})
	manager.SetPosts([]*models.Post{post})
	cache := buildcache.New(filepath.Join(root, "cache"))
	manager.Cache().Set("build_cache", cache)
	plugin := NewImageLibraryPlugin()

	if err := plugin.Write(manager); err != nil {
		t.Fatalf("first external asset Write() error = %v", err)
	}
	jsonPath := filepath.Join(outputDir, "images.json")
	firstJSON, err := os.ReadFile(jsonPath)
	if err != nil {
		t.Fatal(err)
	}
	firstIndex, err := imageindex.Parse(firstJSON)
	if err != nil {
		t.Fatalf("parse first external asset index: %v", err)
	}
	firstImage := imageLibraryTestImageBySrc(t, firstIndex, "/photo.png")
	if firstImage.Width != 3 || firstImage.Height != 2 {
		t.Fatalf("first external asset metadata = %#v, want 3x2", firstImage)
	}
	firstHash := cache.GetImageLibraryHash()
	pagePath := filepath.Join(outputDir, "images", "index.html")
	firstPageInfo, err := os.Stat(pagePath)
	if err != nil {
		t.Fatal(err)
	}

	if err := plugin.Write(manager); err != nil {
		t.Fatalf("unchanged external asset Write() error = %v", err)
	}
	secondPageInfo, err := os.Stat(pagePath)
	if err != nil {
		t.Fatal(err)
	}
	if !secondPageInfo.ModTime().Equal(firstPageInfo.ModTime()) {
		t.Fatal("unchanged external asset caused the image-library page to be rewritten")
	}
	if secondHash := cache.GetImageLibraryHash(); secondHash != firstHash {
		t.Fatalf("unchanged external asset changed image-library hash: before=%q after=%q", firstHash, secondHash)
	}

	writeImageLibraryPNG(t, assetPath, 5, 4)
	if err := plugin.Write(manager); err != nil {
		t.Fatalf("changed external asset Write() error = %v", err)
	}
	secondJSON, err := os.ReadFile(jsonPath)
	if err != nil {
		t.Fatal(err)
	}
	secondIndex, err := imageindex.Parse(secondJSON)
	if err != nil {
		t.Fatalf("parse changed external asset index: %v", err)
	}
	secondImage := imageLibraryTestImageBySrc(t, secondIndex, "/photo.png")
	if secondImage.Width != 5 || secondImage.Height != 4 {
		t.Fatalf("changed external asset metadata = %#v, want 5x4", secondImage)
	}
	if bytes.Equal(firstJSON, secondJSON) || cache.GetImageLibraryHash() == firstHash {
		t.Fatal("changed external asset did not invalidate image-library output")
	}
}

func TestHashImageDirectories_ReusesUnchangedMediaContent(t *testing.T) {
	root := t.TempDir()
	assets := filepath.Join(root, "static")
	if err := os.MkdirAll(filepath.Join(root, "posts"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(assets, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(assets, "asset.webp"), []byte("asset bytes"), 0o600); err != nil {
		t.Fatal(err)
	}
	contentMedia := filepath.Join(root, "posts", "content.webp")
	if err := os.WriteFile(contentMedia, []byte("content bytes"), 0o600); err != nil {
		t.Fatal(err)
	}

	cache := buildcache.New(filepath.Join(root, "cache"))
	reads := 0
	hashFile := func(path string) (string, error) {
		reads++
		return buildcache.HashFile(path)
	}
	firstContentHash, firstStateHash, err := hashImageDirectoriesWithHasher(root, assets, imageHashExtensions(), cache, hashFile)
	if err != nil {
		t.Fatalf("first hashImageDirectories() error = %v", err)
	}
	if reads != 2 {
		t.Fatalf("first media read count = %d, want one read per file and one overlapping traversal", reads)
	}

	secondContentHash, secondStateHash, err := hashImageDirectoriesWithHasher(root, assets, imageHashExtensions(), cache, hashFile)
	if err != nil {
		t.Fatalf("second hashImageDirectories() error = %v", err)
	}
	if reads != 2 {
		t.Fatalf("unchanged second media read count = %d, want %d", reads, 2)
	}
	if firstContentHash != secondContentHash || firstStateHash != secondStateHash {
		t.Fatalf("unchanged hashes differ: first=(%q, %q), second=(%q, %q)", firstContentHash, firstStateHash, secondContentHash, secondStateHash)
	}

	if err := os.WriteFile(contentMedia, []byte("changed content bytes"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, _, err := hashImageDirectoriesWithHasher(root, assets, imageHashExtensions(), cache, hashFile); err != nil {
		t.Fatalf("changed hashImageDirectories() error = %v", err)
	}
	if reads != 3 {
		t.Fatalf("changed media read count = %d, want changed file reread only", reads)
	}
}

func TestHashImageDirectories_ExcludesOutputDirectory(t *testing.T) {
	root := t.TempDir()
	assets := filepath.Join(root, "static")
	output := filepath.Join(root, "output")
	if err := os.MkdirAll(assets, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(output, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(assets, "source.webp"), []byte("source"), 0o600); err != nil {
		t.Fatal(err)
	}
	outputMedia := filepath.Join(output, "generated.webp")
	if err := os.WriteFile(outputMedia, []byte("generated one"), 0o600); err != nil {
		t.Fatal(err)
	}

	cache := buildcache.New(filepath.Join(root, "cache"))
	reads := 0
	hashFile := func(path string) (string, error) {
		reads++
		return buildcache.HashFile(path)
	}
	firstContentHash, _, err := hashImageDirectoriesWithHasher(root, assets, imageHashExtensions(), cache, hashFile, output)
	if err != nil {
		t.Fatalf("first hashImageDirectories() error = %v", err)
	}
	if reads != 1 {
		t.Fatalf("first media read count = %d, want only source media", reads)
	}
	for path := range cache.GetImageLibraryMedia() {
		if filepath.Clean(path) == filepath.Clean(outputMedia) {
			t.Fatal("output media was cached as source input")
		}
	}

	if err := os.WriteFile(outputMedia, []byte("generated two"), 0o600); err != nil {
		t.Fatal(err)
	}
	secondContentHash, _, err := hashImageDirectoriesWithHasher(root, assets, imageHashExtensions(), cache, hashFile, output)
	if err != nil {
		t.Fatalf("second hashImageDirectories() error = %v", err)
	}
	if reads != 1 || firstContentHash != secondContentHash {
		t.Fatalf("output-only change affected source hash: reads=%d first=%q second=%q", reads, firstContentHash, secondContentHash)
	}
}

func TestImageLibraryPlugin_TouchingMediaDoesNotChangeCanonicalHash(t *testing.T) {
	root := t.TempDir()
	output := filepath.Join(root, "output")
	assets := filepath.Join(root, "static")
	if err := os.MkdirAll(assets, 0o755); err != nil {
		t.Fatal(err)
	}
	mediaPath := filepath.Join(assets, "photo.webp")
	if err := os.WriteFile(mediaPath, []byte("same bytes"), 0o600); err != nil {
		t.Fatal(err)
	}

	modelsConfig := models.NewConfig()
	modelsConfig.OutputDir = output
	modelsConfig.AssetsDir = assets
	includeUnreferenced := true
	modelsConfig.Images.IncludeUnreferenced = &includeUnreferenced
	manager := lifecycle.NewManager()
	manager.SetConfig(&lifecycle.Config{
		ContentDir: root,
		OutputDir:  output,
		Extra: map[string]interface{}{
			"assets_dir":    assets,
			"models_config": modelsConfig,
		},
	})
	cache := buildcache.New(filepath.Join(root, "cache"))
	manager.Cache().Set("build_cache", cache)
	plugin := NewImageLibraryPlugin()

	if err := plugin.Write(manager); err != nil {
		t.Fatalf("first Write() error = %v", err)
	}
	firstHash := cache.GetImageLibraryHash()
	firstJSON, err := os.ReadFile(filepath.Join(output, "images.json"))
	if err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(mediaPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chtimes(mediaPath, info.ModTime().Add(2*time.Second), info.ModTime().Add(2*time.Second)); err != nil {
		t.Fatal(err)
	}
	if err := plugin.Write(manager); err != nil {
		t.Fatalf("second Write() error = %v", err)
	}
	if secondHash := cache.GetImageLibraryHash(); secondHash != firstHash {
		t.Fatalf("touch changed canonical image-library hash: before=%q after=%q", firstHash, secondHash)
	}
	secondJSON, err := os.ReadFile(filepath.Join(output, "images.json"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(firstJSON, secondJSON) {
		t.Fatal("touch changed semantic image-library JSON")
	}
}

func TestMarkdownForImage_UsesEscapedCanonicalDestination(t *testing.T) {
	markdown := markdownForImage("/my%20photo%231%25.png", "Photo")
	if markdown != "![Photo](/my%20photo%231%25.png)" {
		t.Fatalf("markdownForImage() = %q", markdown)
	}
}

func TestHashImageDirectories_RehashesSameSizeMediaWithPreservedMtime(t *testing.T) {
	root := t.TempDir()
	mediaPath := filepath.Join(root, "photo.webp")
	if err := os.WriteFile(mediaPath, []byte("first image"), 0o600); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(mediaPath)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := imageLibraryChangeTime(info); !ok {
		t.Skip("filesystem does not expose a change time")
	}

	cache := buildcache.New(filepath.Join(root, "cache"))
	reads := 0
	hashFile := func(path string) (string, error) {
		reads++
		return buildcache.HashFile(path)
	}
	firstContentHash, _, err := hashImageDirectoriesWithHasher(root, "", imageHashExtensions(), cache, hashFile)
	if err != nil {
		t.Fatalf("first hashImageDirectories() error = %v", err)
	}

	if err := os.WriteFile(mediaPath, []byte("second img!"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Chtimes(mediaPath, info.ModTime(), info.ModTime()); err != nil {
		t.Fatal(err)
	}
	secondContentHash, _, err := hashImageDirectoriesWithHasher(root, "", imageHashExtensions(), cache, hashFile)
	if err != nil {
		t.Fatalf("second hashImageDirectories() error = %v", err)
	}
	if reads != 2 {
		t.Fatalf("same-size preserved-mtime media read count = %d, want changed file reread", reads)
	}
	if firstContentHash == secondContentHash {
		t.Fatal("same-size media replacement retained the cached content hash")
	}
}

func TestHashImageDirectory_NormalizesMediaExtensionCase(t *testing.T) {
	root := t.TempDir()
	mediaPath := filepath.Join(root, "clip.Mp4")
	if err := os.WriteFile(mediaPath, []byte("first video"), 0o600); err != nil {
		t.Fatal(err)
	}

	first, _, err := hashImageDirectory(root, imageHashExtensions())
	if err != nil {
		t.Fatalf("first hashImageDirectory() error = %v", err)
	}
	if err := os.WriteFile(mediaPath, []byte("second video"), 0o600); err != nil {
		t.Fatal(err)
	}
	second, _, err := hashImageDirectory(root, imageHashExtensions())
	if err != nil {
		t.Fatalf("second hashImageDirectory() error = %v", err)
	}
	if first == second {
		t.Fatalf("mixed-case media content hash did not change: %q", first)
	}
}

func TestImageLibraryPlugin_IgnoresDanglingImageSymlinkInContentHash(t *testing.T) {
	root := t.TempDir()
	output := filepath.Join(root, "output")
	assets := filepath.Join(root, "static")
	link := filepath.Join(root, "ignored.png")
	if err := os.Symlink(filepath.Join(root, "missing.png"), link); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}

	modelsConfig := models.NewConfig()
	modelsConfig.OutputDir = output
	modelsConfig.AssetsDir = assets
	manager := lifecycle.NewManager()
	manager.SetConfig(&lifecycle.Config{
		ContentDir: root,
		OutputDir:  output,
		Extra: map[string]interface{}{
			"assets_dir":    assets,
			"models_config": modelsConfig,
		},
	})

	if err := NewImageLibraryPlugin().Write(manager); err != nil {
		t.Fatalf("Write() error = %v", err)
	}
}

func TestImageLibraryPlugin_RejectsSymlinkedOutputRoot(t *testing.T) {
	root := t.TempDir()
	realOutput := filepath.Join(root, "real-output")
	output := filepath.Join(root, "output")
	if err := os.MkdirAll(realOutput, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(realOutput, output); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	assets := filepath.Join(root, "static")
	modelsConfig := models.NewConfig()
	modelsConfig.OutputDir = output
	modelsConfig.AssetsDir = assets
	manager := lifecycle.NewManager()
	manager.SetConfig(&lifecycle.Config{
		ContentDir: root,
		OutputDir:  output,
		Extra: map[string]interface{}{
			"assets_dir":    assets,
			"models_config": modelsConfig,
		},
	})

	err := NewImageLibraryPlugin().Write(manager)
	if err == nil || !strings.Contains(err.Error(), "symlink") {
		t.Fatalf("Write() error = %v, want symlink-output-root error", err)
	}
	if _, err := os.Stat(filepath.Join(realOutput, "images", "index.html")); !os.IsNotExist(err) {
		t.Fatalf("symlink target was written, stat error = %v", err)
	}
}

func TestImageLibraryPlugin_PathChangePreservesModifiedOldOutput(t *testing.T) {
	root := t.TempDir()
	output := filepath.Join(root, "output")
	assets := filepath.Join(root, "static")
	oldPage := filepath.Join(output, "images", "index.html")
	if err := os.MkdirAll(filepath.Dir(oldPage), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(oldPage, []byte("user content"), 0o600); err != nil {
		t.Fatal(err)
	}

	modelsConfig := models.NewConfig()
	modelsConfig.OutputDir = output
	modelsConfig.AssetsDir = assets
	modelsConfig.Images.Path = "gallery"
	manager := lifecycle.NewManager()
	manager.SetConfig(&lifecycle.Config{
		ContentDir: root,
		OutputDir:  output,
		Extra: map[string]interface{}{
			"assets_dir":    assets,
			"models_config": modelsConfig,
		},
	})
	cache := buildcache.New(filepath.Join(root, "cache"))
	cache.SetImageLibraryOutputs(root, output, map[string]string{oldPage: buildcache.ContentHash("generated content")})
	manager.Cache().Set("build_cache", cache)
	plugin := NewImageLibraryPlugin()

	if err := plugin.Write(manager); err == nil || !strings.Contains(err.Error(), "was modified") {
		t.Fatalf("Write() error = %v, want modified-output error", err)
	}
	page, err := os.ReadFile(oldPage)
	if err != nil {
		t.Fatal(err)
	}
	if string(page) != "user content" {
		t.Fatalf("modified old output changed: %q", page)
	}
	if _, err := os.Stat(filepath.Join(output, "gallery", "index.html")); !os.IsNotExist(err) {
		t.Fatalf("new output was written after collision, stat error = %v", err)
	}
}

func TestImageLibraryPlugin_PathChangeCleansPreviousOutputs(t *testing.T) {
	root := t.TempDir()
	output := filepath.Join(root, "output")
	assets := filepath.Join(root, "static")
	modelsConfig := models.NewConfig()
	modelsConfig.OutputDir = output
	modelsConfig.AssetsDir = assets
	manager := lifecycle.NewManager()
	manager.SetConfig(&lifecycle.Config{
		ContentDir: root,
		OutputDir:  output,
		Extra: map[string]interface{}{
			"assets_dir":    assets,
			"models_config": modelsConfig,
		},
	})
	manager.Cache().Set("build_cache", buildcache.New(filepath.Join(root, "cache")))
	plugin := NewImageLibraryPlugin()

	if err := plugin.Write(manager); err != nil {
		t.Fatalf("first Write() error = %v", err)
	}
	oldPage := filepath.Join(output, "images", "index.html")
	if _, err := os.Stat(oldPage); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(output, "images.json")); err != nil {
		t.Fatalf("root image index missing: %v", err)
	}
	modelsConfig.Images.Path = "gallery"
	if err := plugin.Write(manager); err != nil {
		t.Fatalf("second Write() error = %v", err)
	}
	if _, err := os.Stat(oldPage); !os.IsNotExist(err) {
		t.Fatalf("old image-library page was not removed, stat error = %v", err)
	}
	if _, err := os.Stat(filepath.Join(output, "gallery", "index.html")); err != nil {
		t.Fatalf("new image-library page missing: %v", err)
	}
	if _, err := os.Stat(filepath.Join(output, "images.json")); err != nil {
		t.Fatalf("root image index was not retained: %v", err)
	}
}

func TestImageLibraryPlugin_RejectsOutputCollisions(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(t *testing.T, root string)
		posts      []*models.Post
		wantPhrase string
	}{
		{
			name:       "post",
			posts:      []*models.Post{{Path: "posts/images.md", Slug: "images", Published: true}},
			wantPhrase: "conflicts with post",
		},
		{
			name: "static file",
			setup: func(t *testing.T, root string) {
				t.Helper()
				path := filepath.Join(root, "static", "images", "index.html")
				if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(path, []byte("owned by site"), 0o600); err != nil {
					t.Fatal(err)
				}
			},
			wantPhrase: "conflicts with static file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()
			if tt.setup != nil {
				tt.setup(t, root)
			}
			modelsConfig := models.NewConfig()
			modelsConfig.OutputDir = filepath.Join(root, "output")
			modelsConfig.AssetsDir = filepath.Join(root, "static")
			manager := lifecycle.NewManager()
			manager.SetConfig(&lifecycle.Config{
				ContentDir: root,
				OutputDir:  modelsConfig.OutputDir,
				Extra: map[string]interface{}{
					"assets_dir":    modelsConfig.AssetsDir,
					"models_config": modelsConfig,
				},
			})
			manager.SetPosts(tt.posts)
			manager.RegisterPlugin(NewImageLibraryPlugin())

			if err := manager.RunTo(lifecycle.StageWrite); err != nil {
				t.Fatalf("RunTo(write) error = %v", err)
			}
			warnings := manager.Warnings()
			if len(warnings) != 1 || !strings.Contains(warnings[0].Error(), tt.wantPhrase) {
				t.Fatalf("warnings = %v, want %q", warnings, tt.wantPhrase)
			}
		})
	}
}

func imageIndexForTest() imageindex.Index {
	return imageindex.Index{
		Generator: imageindex.Generator{Name: imageindex.GeneratorName, Version: "test"},
		Images: []imageindex.Image{
			{Src: "/static/unused.png", Alt: "unused"},
			{Src: "https://dropper.wayl.one/file/photo.webp?token=keep", Alt: "photo", Cover: true, Uses: []imageindex.Use{{Href: "/photo/", Title: "Photo", Cover: true}}},
		},
	}
}

func stringPointerImages(value string) *string {
	return &value
}

func writeImageLibraryPNG(t *testing.T, path string, width, height int) {
	t.Helper()
	file, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	canvas := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			canvas.Set(x, y, color.RGBA{R: 20, G: 90, B: 70, A: 255})
		}
	}
	if err := png.Encode(file, canvas); err != nil {
		t.Fatal(err)
	}
}

func imageLibraryTestImageBySrc(t *testing.T, index imageindex.Index, src string) imageindex.Image {
	t.Helper()
	for i := range index.Images {
		image := &index.Images[i]
		if image.Src == src {
			return *image
		}
	}
	t.Fatalf("image %q not found in %#v", src, index.Images)
	return imageindex.Image{}
}

func TestImageLibraryPaginate_PreservesInventoryAndNativeLinks(t *testing.T) {
	index := imageindex.Index{Images: make([]imageindex.Image, 105)}
	for i := range index.Images {
		index.Images[i].Src = fmt.Sprintf("/static/picture-%03d.webp", i)
	}
	full := newImageLibraryPage(index)
	pages := imageLibraryPaginate(full, 48, "images")
	if len(pages) != 3 || len(pages[0].Images) != 48 || len(pages[1].Images) != 48 || len(pages[2].Images) != 9 {
		t.Fatalf("page sizes: %v", []int{len(pages[0].Images), len(pages[1].Images), len(pages[2].Images)})
	}
	for i, page := range pages {
		if page.Count != 105 || page.PageNumber != i+1 || page.TotalPages != 3 || len(page.PageLinks) != 3 {
			t.Fatalf("page %d metadata: %#v", i+1, page)
		}
		if !page.PageLinks[i].Current {
			t.Fatalf("page %d current marker missing", i+1)
		}
	}
	if pages[1].PreviousURL != "/images/" || pages[1].NextURL != "/images/page/3/" {
		t.Fatalf("page 2 navigation: %#v", pages[1])
	}
	if pages[1].Images[0].Anchor != "image-49" || pages[2].Images[0].Anchor != "image-97" {
		t.Fatalf("unstable anchors across page boundaries")
	}
	if len(full.Images) != 105 {
		t.Fatal("pagination mutated the full source inventory")
	}
	legacy := imageLibraryPaginate(full, 0, "images")
	if len(legacy) != 1 || len(legacy[0].Images) != 105 || legacy[0].TotalPages != 1 {
		t.Fatal("legacy unpaginated behavior changed")
	}
}

func TestParseImageLibraryConfig_PageSize(t *testing.T) {
	got, err := parseImageLibraryConfig(map[string]interface{}{
		"images": map[string]interface{}{"page_size": int64(48)},
	})
	if err != nil || got.PageSize != 48 {
		t.Fatalf("page_size parse: %d, %v", got.PageSize, err)
	}
	for _, invalid := range []interface{}{-1, 1.5, "48", true} {
		_, err := parseImageLibraryConfig(map[string]interface{}{
			"images": map[string]interface{}{"page_size": invalid},
		})
		if err == nil {
			t.Errorf("accepted invalid page_size %#v", invalid)
		}
	}
}

func TestImageLibraryPaginationTemplateContract(t *testing.T) {
	for _, path := range []string{
		filepath.Join("..", "..", "templates", "images.html"),
		filepath.Join("..", "themes", "default", "templates", "images.html"),
	} {
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		for _, marker := range []string{"data-image-total-pages", "data-image-pagination", "image.Anchor", "image_library.PageLinks"} {
			if !strings.Contains(string(data), marker) {
				t.Errorf("%s missing pagination contract %q", path, marker)
			}
		}
	}
	data, err := os.ReadFile(filepath.Join("..", "themes", "default", "static", "js", "image-library.js"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "hydrateAllPages") || !strings.Contains(string(data), "scheduleApply") {
		t.Fatal("image library has no on-demand cross-page search hydration")
	}
}

func TestImageLibraryPlugin_PaginatedOutputsAndStaleCleanup(t *testing.T) {
	root := t.TempDir()
	output := filepath.Join(root, "output")
	assets := filepath.Join(root, "static")
	if err := os.MkdirAll(assets, 0o755); err != nil {
		t.Fatal(err)
	}
	post := models.NewPost("posts/album.md")
	post.Slug = "album"
	post.Href = "/album/"
	post.Published = true
	post.Title = stringPointerImages("Album")
	post.Content = strings.Join([]string{
		"![A](https://dropper.wayl.one/file/a.webp)",
		"![B](https://dropper.wayl.one/file/b.webp)",
		"![C](https://dropper.wayl.one/file/c.webp)",
		"![D](https://dropper.wayl.one/file/d.webp)",
		"![E](https://dropper.wayl.one/file/e.webp)",
	}, "\n")

	modelsConfig := models.NewConfig()
	modelsConfig.OutputDir = output
	modelsConfig.AssetsDir = assets
	modelsConfig.Images = models.NewImagesConfig()
	modelsConfig.Images.PageSize = 2
	manager := lifecycle.NewManager()
	manager.SetConfig(&lifecycle.Config{
		ContentDir: root, OutputDir: output,
		Extra: map[string]interface{}{
			"assets_dir": assets,
			"models_config": modelsConfig,
		},
	})
	manager.SetPosts([]*models.Post{post})
	cache := buildcache.New(filepath.Join(root, "cache"))
	manager.Cache().Set("build_cache", cache)
	plugin := NewImageLibraryPlugin()

	if err := plugin.Write(manager); err != nil {
		t.Fatalf("first paginated build: %v", err)
	}
	paths := []string{
		filepath.Join(output, "images", "index.html"),
		filepath.Join(output, "images", "page", "2", "index.html"),
		filepath.Join(output, "images", "page", "3", "index.html"),
	}
	for index, path := range paths {
		page, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("page %d missing: %v", index+1, err)
		}
		want := 2
		if index == 2 {
			want = 1
		}
		if got := strings.Count(string(page), "data-image-card"); got != want {
			t.Errorf("page %d has %d cards, want %d", index+1, got, want)
		}
		if !strings.Contains(string(page), "data-image-pagination") {
			t.Errorf("page %d has no no-JS pagination links", index+1)
		}
	}
	jsonData, err := os.ReadFile(filepath.Join(output, "images.json"))
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.Count(string(jsonData), `"src":`); got != 5 {
		t.Fatalf("canonical index lost entries: %d", got)
	}
	if err := plugin.Write(manager); err != nil {
		t.Fatalf("same-input rebuild: %v", err)
	}
	if _, err := os.Stat(paths[2]); err != nil {
		t.Fatalf("cached page 3 vanished: %v", err)
	}

	modelsConfig.Images.PageSize = 3
	if err := plugin.Write(manager); err != nil {
		t.Fatalf("page size reduction rebuild: %v", err)
	}
	if _, err := os.Stat(paths[2]); !os.IsNotExist(err) {
		t.Fatalf("stale page 3 still present, stat error=%v", err)
	}
	page2, err := os.ReadFile(paths[1])
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.Count(string(page2), "data-image-card"); got != 2 {
		t.Fatalf("page 2 with page_size=3 has %d cards, want 2", got)
	}
}
