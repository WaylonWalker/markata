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
