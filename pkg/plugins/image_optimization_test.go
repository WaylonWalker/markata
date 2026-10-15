package plugins

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata-go/pkg/lifecycle"
	"github.com/WaylonWalker/markata-go/pkg/models"
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
	if !strings.Contains(post.HTML, `<picture>`) || !strings.Contains(post.HTML, `srcset="/8bitcc-480w.avif 480w`) {
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
