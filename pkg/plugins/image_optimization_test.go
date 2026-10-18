package plugins

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"golang.org/x/net/html"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

func TestImageOptimization_Render_WrapsLocalImages(t *testing.T) {
	plugin := NewImageOptimizationPlugin()
	plugin.config = defaultImageOptimizationConfig()
	plugin.availableFormats = []string{"avif", "webp"}
	plugin.config.Widths = []int{480, 960}
	plugin.config.Sizes = "(max-width: 960px) 100vw, 960px"

	post := &models.Post{
		ArticleHTML: `<p><img src="/images/cat.jpg" alt="Cat"></p><p><img src="https://example.com/dog.jpg" alt="Dog"></p>`,
		Slug:        "post",
	}

	if !isLocalImageSrc("/images/cat.jpg") {
		t.Fatalf("expected local image src to be recognized")
	}
	if !isOptimizableImageSrc("/images/cat.jpg") {
		t.Fatalf("expected optimizable image src to be recognized")
	}
	if sources := buildPictureSources("/images/cat.jpg", plugin.availableFormats, plugin.config.Widths, plugin.config.Sizes); len(sources) == 0 {
		t.Fatalf("expected picture sources to be generated")
	}

	if err := plugin.processPost(post); err != nil {
		t.Fatalf("processPost error: %v", err)
	}

	if !strings.Contains(post.ArticleHTML, "<picture>") {
		t.Fatalf("expected picture wrapper, got: %s", post.ArticleHTML)
	}
	if !strings.Contains(post.ArticleHTML, `srcset="/images/cat-480w.avif 480w, /images/cat-960w.avif 960w"`) {
		t.Fatalf("expected AVIF source, got: %s", post.ArticleHTML)
	}
	if !strings.Contains(post.ArticleHTML, `srcset="/images/cat-480w.webp 480w, /images/cat-960w.webp 960w"`) {
		t.Fatalf("expected WebP source, got: %s", post.ArticleHTML)
	}
	if strings.Contains(post.ArticleHTML, "example.com/dog") && strings.Contains(post.ArticleHTML, "<picture>") {
		if strings.Contains(post.ArticleHTML, "dog.webp") || strings.Contains(post.ArticleHTML, "dog.avif") {
			t.Fatalf("expected external image to remain unchanged")
		}
	}

	if post.Extra == nil {
		t.Fatalf("expected post.Extra to be set")
	}
	if _, ok := post.Extra["image_optimization"]; !ok {
		t.Fatalf("expected image_optimization targets in post.Extra")
	}
}

func TestImageOptimization_StatsCountMissingSource(t *testing.T) {
	plugin := NewImageOptimizationPlugin()
	plugin.config = defaultImageOptimizationConfig()
	if err := plugin.optimizeTarget(t.TempDir(), t.TempDir(), imageOptimizationTarget{Src: "/images/missing.jpg"}); err != nil {
		t.Fatalf("optimizeTarget error: %v", err)
	}
	if plugin.stats.targets != 1 || plugin.stats.missingSources != 1 {
		t.Fatalf("unexpected image stats: %+v", plugin.stats)
	}
}

func TestImageOptimization_StatsCountEncoderCallsAndCacheHits(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("test fake encoder requires a POSIX shell")
	}
	root := t.TempDir()
	source := filepath.Join(root, "photo.jpg")
	cache := filepath.Join(root, "cache")
	if err := os.WriteFile(source, []byte("fake JPEG for subprocess test"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(cache, 0o700); err != nil {
		t.Fatal(err)
	}
	encoder := filepath.Join(root, "fake-cwebp")
	script := "#!/bin/sh\nwhile [ \"$#\" -gt 1 ]; do shift; done\nprintf 'webp' > \"$1\"\n"
	if err := os.WriteFile(encoder, []byte(script), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(encoder, 0o700); err != nil {
		t.Fatal(err)
	}
	plugin := NewImageOptimizationPlugin()
	plugin.cwebpPath = encoder
	info, err := os.Stat(source)
	if err != nil {
		t.Fatal(err)
	}
	if err := plugin.optimizeImageFormat(source, cache, info, formatWebP, nil); err != nil {
		t.Fatal(err)
	}
	if err := plugin.optimizeImageFormat(source, cache, info, formatWebP, nil); err != nil {
		t.Fatal(err)
	}
	if plugin.stats.variantsChecked != 2 || plugin.stats.cacheMisses != 1 ||
		plugin.stats.cacheHits != 1 || plugin.stats.encoderCalls != 1 || plugin.stats.encoderFailures != 0 {
		t.Fatalf("unexpected image optimization stats: %+v", plugin.stats)
	}
}

func TestImageOptimization_Render_SkipsExistingPicture(t *testing.T) {
	plugin := NewImageOptimizationPlugin()
	plugin.config = defaultImageOptimizationConfig()
	plugin.availableFormats = []string{"avif"}

	post := &models.Post{
		ArticleHTML: `<picture><img src="/images/cat.jpg" alt="Cat"></picture>`,
		Slug:        "post",
	}

	if err := plugin.processPost(post); err != nil {
		t.Fatalf("processPost error: %v", err)
	}

	if strings.Count(post.ArticleHTML, "<picture>") != 1 {
		t.Fatalf("expected existing picture untouched, got: %s", post.ArticleHTML)
	}
	if post.Extra != nil {
		if _, ok := post.Extra["image_optimization"]; ok {
			t.Fatalf("did not expect image_optimization targets")
		}
	}
}

func TestImageOptimization_RenderAddsResponsiveCandidatesToTrustedRemoteImages(t *testing.T) {
	plugin := NewImageOptimizationPlugin()
	post := &models.Post{
		ArticleHTML: `<figure><a href="https://dropper.waylonwalker.com/file/photo.webp"><img src="https://dropper.waylonwalker.com/file/photo.webp" alt="Photo"></a></figure><img src="https://example.com/photo.webp" alt="Other">`,
		Slug:        "post",
	}

	if err := plugin.processPost(post); err != nil {
		t.Fatalf("processPost error: %v", err)
	}

	if !strings.Contains(post.ArticleHTML, `href="https://dropper.waylonwalker.com/file/photo.webp"`) {
		t.Fatalf("expected zoom link to keep original URL, got: %s", post.ArticleHTML)
	}
	if !strings.Contains(post.ArticleHTML, `https://dropper.waylonwalker.com/file/photo.webp?w=320 320w`) {
		t.Fatalf("expected responsive candidates for trusted remote image, got: %s", post.ArticleHTML)
	}
	if !strings.Contains(post.ArticleHTML, `sizes="(max-width: 760px) 100vw, 760px"`) {
		t.Fatalf("expected responsive image sizes, got: %s", post.ArticleHTML)
	}
	if strings.Contains(post.ArticleHTML, `srcset="https://example.com/photo.webp`) {
		t.Fatalf("did not expect candidates for untrusted image URLs, got: %s", post.ArticleHTML)
	}
	if _, ok := post.Extra["image_optimization"]; ok {
		t.Fatalf("remote image should not create local optimization targets")
	}
}

func TestImageOptimization_RenderHandlesRemoteImagesWithoutLocalEncoders(t *testing.T) {
	plugin := NewImageOptimizationPlugin()
	plugin.config = defaultImageOptimizationConfig()
	plugin.availableFormats = nil
	manager := lifecycle.NewManager()
	post := &models.Post{
		ArticleHTML: `<figure><a href="https://dropper.waylonwalker.com/file/photo.webp"><img src="https://dropper.waylonwalker.com/file/photo.webp" alt="Photo"></a></figure>`,
	}
	manager.SetPosts([]*models.Post{post})

	if err := plugin.Render(manager); err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	if !strings.Contains(post.ArticleHTML, `?w=320 320w`) {
		t.Fatalf("expected responsive candidates without local encoders, got: %s", post.ArticleHTML)
	}
}

func TestImageOptimizationWriteRunsBeforeHTMLPublication(t *testing.T) {
	outputDir := t.TempDir()
	manager := lifecycle.NewManager()
	config := manager.Config()
	config.OutputDir = outputDir
	config.Extra["image_optimization"] = map[string]any{
		"enabled": true,
		"formats": []any{""},
	}

	post := models.NewPost("post.md")
	post.Slug = "post"
	post.Published = true
	post.HTML = `<html><body><img src="https://dropper.waylonwalker.com/photo.webp" alt="Photo"></body></html>`
	manager.SetPosts([]*models.Post{post})
	manager.RegisterPlugins(NewImageOptimizationPlugin(), NewPublishHTMLPlugin())

	if err := manager.RunTo(lifecycle.StageWrite); err != nil {
		t.Fatalf("RunTo(write) error = %v", err)
	}

	output, err := os.ReadFile(filepath.Join(outputDir, "post", "index.html"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(output), `srcset="https://dropper.waylonwalker.com/photo.webp?w=32 32w`) {
		t.Fatalf("published HTML is missing responsive image candidates: %s", output)
	}
}

func TestImageOptimizationWritePriorityFollowsStaticAssetsAndPrecedesPublishHTML(t *testing.T) {
	staticAssets := NewStaticAssetsPlugin().Priority(lifecycle.StageWrite)
	imageOptimization := NewImageOptimizationPlugin().Priority(lifecycle.StageWrite)
	if !(staticAssets < imageOptimization && imageOptimization < lifecycle.PriorityDefault) {
		t.Fatalf("write priorities are static_assets=%d image_optimization=%d publish_html=%d, want static_assets < image_optimization < publish_html", staticAssets, imageOptimization, lifecycle.PriorityDefault)
	}
}

func TestLocalImageResponsiveSettingsRespectDeclaredWidth(t *testing.T) {
	image := &html.Node{Type: html.ElementNode, Data: "img"}
	setAttr(image, "width", "160")

	widths, sizes := localImageResponsiveSettings(image, []int{480, 768, 1200}, "100vw")
	if !reflect.DeepEqual(widths, []int{160}) {
		t.Fatalf("responsive widths = %v, want [160]", widths)
	}
	if sizes != "160px" {
		t.Fatalf("sizes = %q, want 160px", sizes)
	}
}

func TestPrepareAVIFInputResizesToPNG(t *testing.T) {
	sourcePath := filepath.Join(t.TempDir(), "source.png")
	source := image.NewNRGBA(image.Rect(0, 0, 4, 2))
	for y := 0; y < 2; y++ {
		for x := 0; x < 4; x++ {
			source.Set(x, y, color.NRGBA{R: 255, A: 255})
		}
	}
	file, err := os.Create(sourcePath)
	if err != nil {
		t.Fatal(err)
	}
	if err := png.Encode(file, source); err != nil {
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}

	resizedPath, cleanup, err := prepareAVIFInput(sourcePath, 2)
	if err != nil {
		t.Fatalf("prepareAVIFInput() error = %v", err)
	}
	defer cleanup()
	resized, err := os.Open(resizedPath)
	if err != nil {
		t.Fatal(err)
	}
	config, err := png.DecodeConfig(resized)
	if err != nil {
		t.Fatal(err)
	}
	if err := resized.Close(); err != nil {
		t.Fatal(err)
	}
	if config.Width != 2 || config.Height != 1 {
		t.Fatalf("resized dimensions = %dx%d, want 2x1", config.Width, config.Height)
	}
}

func TestAVIFEncoderArgsUseSupportedQualityFlag(t *testing.T) {
	args := avifEncoderArgs("input.png", "output.avif", 80)
	want := []string{"-q", "80", "input.png", "output.avif"}
	if !reflect.DeepEqual(args, want) {
		t.Fatalf("avifEncoderArgs() = %v, want %v", args, want)
	}
}

func TestImageOptimization_RenderKeepsExistingRemoteSrcset(t *testing.T) {
	plugin := NewImageOptimizationPlugin()
	post := &models.Post{ArticleHTML: `<img src="https://dropper.waylonwalker.com/file/photo.webp" srcset="small.webp 320w, large.webp 640w" sizes="320px">`}

	if err := plugin.processPost(post); err != nil {
		t.Fatalf("processPost error: %v", err)
	}
	if !strings.Contains(post.ArticleHTML, `srcset="small.webp 320w, large.webp 640w"`) {
		t.Fatalf("existing srcset should remain unchanged, got: %s", post.ArticleHTML)
	}
}

func TestImageOptimization_RenderUsesExplicitRemoteImageWidth(t *testing.T) {
	plugin := NewImageOptimizationPlugin()
	post := &models.Post{ArticleHTML: `<img src="https://dropper.waylonwalker.com/file/icon.webp" width="44" height="44">`}

	if err := plugin.processPost(post); err != nil {
		t.Fatalf("processPost error: %v", err)
	}
	if !strings.Contains(post.ArticleHTML, `sizes="44px"`) || !strings.Contains(post.ArticleHTML, `?w=88 88w`) {
		t.Fatalf("expected width-aware 2x candidates for small images, got: %s", post.ArticleHTML)
	}
}

func TestImageOptimization_RenderSkipsSmallRemoteImagesWithoutWidth(t *testing.T) {
	plugin := NewImageOptimizationPlugin()
	post := &models.Post{ArticleHTML: `<img src="https://dropper.waylonwalker.com/file/icon.webp" height="24">`}

	if err := plugin.processPost(post); err != nil {
		t.Fatalf("processPost error: %v", err)
	}
	if strings.Contains(post.ArticleHTML, `srcset=`) {
		t.Fatalf("small remote images without a known width should remain untouched, got: %s", post.ArticleHTML)
	}
}

func TestImageOptimization_WriteAddsResponsiveCandidatesToCachedFullPageHTML(t *testing.T) {
	plugin := NewImageOptimizationPlugin()
	plugin.config = defaultImageOptimizationConfig()
	plugin.availableFormats = nil
	post := &models.Post{
		ArticleHTML: `<p>original article body</p>`,
		HTML:        `<!doctype html><html><head><title>Home</title></head><body><img src="https://dropper.waylonwalker.com/file/home.webp" width="160" height="160"><p>Home body</p></body></html>`,
	}
	manager := lifecycle.NewManager()
	manager.SetPosts([]*models.Post{post})

	if err := plugin.Write(manager); err != nil {
		t.Fatalf("Write() error = %v", err)
	}
	if !strings.Contains(post.HTML, `https://dropper.waylonwalker.com/file/home.webp?w=320 320w`) {
		t.Fatalf("expected responsive candidates on rendered page image, got: %s", post.HTML)
	}
	if !strings.Contains(post.HTML, `<title>Home</title>`) || !strings.Contains(post.HTML, `<p>Home body</p>`) {
		t.Fatalf("expected full-page document content to be preserved, got: %s", post.HTML)
	}
	if post.ArticleHTML != `<p>original article body</p>` {
		t.Fatalf("ArticleHTML changed while processing cached full-page HTML: %s", post.ArticleHTML)
	}
}

func TestImageOptimization_WriteAddsLocalTemplateImagesToOptimizationTargets(t *testing.T) {
	plugin := NewImageOptimizationPlugin()
	plugin.config = defaultImageOptimizationConfig()
	plugin.availableFormats = []string{formatAVIF}
	post := &models.Post{
		Slug: "home",
		HTML: `<!doctype html><html><body><img class="home-avatar" src="/8bitcc.png" width="160" height="160"></body></html>`,
	}

	if err := plugin.processRenderedHTML(post); err != nil {
		t.Fatalf("processRenderedHTML() error = %v", err)
	}
	if !strings.Contains(post.HTML, `<picture>`) || !strings.Contains(post.HTML, `srcset="/8bitcc-160w.avif 160w`) {
		t.Fatalf("expected optimized picture markup for local template image, got: %s", post.HTML)
	}
	targets := imageOptimizationTargetsFrom(post.Extra["image_optimization"])
	if len(targets) != 1 || targets[0].Src != "/8bitcc.png" || targets[0].PostSlug != "home" {
		t.Fatalf("local template optimization targets = %#v, want one homepage avatar target", targets)
	}
}

func TestImageOptimization_ResolveOutputPath_Relative(t *testing.T) {
	target := imageOptimizationTarget{
		Src:      "images/cat.jpg",
		PostSlug: "post",
	}
	path, err := resolveImageOutputPath("output", target)
	if err != nil {
		t.Fatalf("resolveImageOutputPath error: %v", err)
	}
	normalized := filepath.ToSlash(path)
	if !strings.HasSuffix(normalized, "output/post/images/cat.jpg") {
		t.Fatalf("unexpected output path: %s", path)
	}
}

func TestImageOptimization_PreservesExplicitHomepageSlug(t *testing.T) {
	plugin := NewImageOptimizationPlugin()
	plugin.config = defaultImageOptimizationConfig()
	plugin.availableFormats = []string{formatAVIF}
	post := &models.Post{
		Path: "pages/post/index.md",
		HTML: `<html><body><img src="/8bitcc.png" width="160" height="160"></body></html>`,
	}
	post.Set("_slug_explicit", true)
	if err := plugin.processRenderedHTML(post); err != nil {
		t.Fatal(err)
	}
	if post.Slug != "" {
		t.Fatalf("explicit homepage slug changed to %q", post.Slug)
	}
	targets := imageOptimizationTargetsFrom(post.Extra["image_optimization"])
	if len(targets) != 1 || targets[0].PostSlug != "" {
		t.Fatalf("homepage target must retain empty slug: %#v", targets)
	}
}
