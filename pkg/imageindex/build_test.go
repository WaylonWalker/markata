package imageindex

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/WaylonWalker/markata/pkg/models"
	"golang.org/x/image/tiff"
)

func TestBuild_DiscoversPublicSourcesAndDeduplicatesUses(t *testing.T) {
	root := t.TempDir()
	assets := filepath.Join(root, "static")
	if err := os.MkdirAll(filepath.Join(assets, "images"), 0o755); err != nil {
		t.Fatal(err)
	}
	writePNG(t, filepath.Join(assets, "images", "photo.png"), 3, 2)
	writePNG(t, filepath.Join(assets, "images", "unused.png"), 4, 4)
	writePNG(t, filepath.Join(assets, "images", "private.png"), 5, 5)

	public := models.NewPost(filepath.Join(root, "posts", "hello.md"))
	public.Path = "posts/hello.md"
	public.Slug = "hello"
	public.Href = "/hello/"
	public.Published = true
	public.Content = "![Authored photo](images/photo.png)\n\n<img src=\"/static/images/photo.png\" alt=\"HTML photo\">\n\n![Inline](data:image/png;base64,abc) ![Blob](blob:abc) ![Script](javascript:alert(1)) ![Credential](https://user:SECRET@example.test/private.png)"
	public.ArticleHTML = `<p><img src="/images/photo.png" alt="Rendered photo"></p>`
	public.Extra["cover"] = "/static/images/photo.png"
	public.Extra["cover_image"] = "https://cdn.example.test/second.webp"
	public.Title = stringPointer("Hello")

	private := models.NewPost("posts/private.md")
	private.Path = "posts/private.md"
	private.Published = true
	private.Private = true
	private.Content = "![Secret](secret.png) ![Private asset](/images/private.png)"

	index, err := Build([]*models.Post{public, private}, BuildOptions{
		ContentDir:          root,
		AssetsDir:           assets,
		GeneratorVersion:    "test",
		IncludeUnreferenced: true,
	})
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}
	if index.ImageCount != 4 {
		t.Fatalf("ImageCount = %d, want 4 (%#v)", index.ImageCount, index.Images)
	}

	photo := imageBySrc(t, index, "/images/photo.png")
	if photo.Width != 3 || photo.Height != 2 || photo.MIMEType != "image/png" {
		t.Fatalf("photo metadata = %#v", photo)
	}
	if !photo.Cover || len(photo.Uses) != 1 || !photo.Uses[0].Cover {
		t.Fatalf("photo uses = %#v", photo.Uses)
	}
	if photo.Alt != "Authored photo" {
		t.Fatalf("photo alt = %q, want first authored alt", photo.Alt)
	}
	if len(imageBySrc(t, index, "https://cdn.example.test/second.webp").Uses) != 1 {
		t.Fatal("frontmatter cover_image was not indexed")
	}
	if len(imageBySrc(t, index, "https://cdn.example.test/second.webp").Uses) != 1 || imageBySrc(t, index, "https://cdn.example.test/second.webp").Cover {
		t.Fatal("cover_image should be a non-cover relationship when cover is present")
	}
	if len(imageBySrc(t, index, "/images/unused.png").Uses) != 0 {
		t.Fatal("unreferenced image has unexpected public uses")
	}
	if _, ok := findImage(index, "secret.png"); ok {
		t.Fatal("private image leaked into index")
	}
	privateImage, ok := findImage(index, "/images/private.png")
	if !ok {
		t.Fatal("unreferenced static image was omitted because a private body referenced it")
	}
	if len(privateImage.Uses) != 0 {
		t.Fatalf("private body created a public use: %#v", privateImage.Uses)
	}
	for _, source := range []string{"data:image/png;base64,abc", "blob:abc", "javascript:alert(1)", "https://user:SECRET@example.test/private.png"} {
		if _, ok := findImage(index, source); ok {
			t.Fatalf("unsafe image source %q leaked into index", source)
		}
	}
}

func TestBuild_PrivateBodyHasNoInfluenceOnPublicInventory(t *testing.T) {
	root := t.TempDir()
	assets := filepath.Join(root, "static")
	if err := os.MkdirAll(assets, 0o755); err != nil {
		t.Fatal(err)
	}
	writePNG(t, filepath.Join(assets, "unreferenced.png"), 3, 2)

	public := models.NewPost("posts/public.md")
	public.Path = "posts/public.md"
	public.Slug = "public"
	public.Href = "/public/"
	public.Published = true
	public.Content = "![Public](https://example.test/public.png)"

	private := models.NewPost("posts/private.md")
	private.Path = "posts/private.md"
	private.Slug = "private"
	private.Href = "/private/"
	private.Published = true
	private.Private = true
	private.Content = "secret text"
	private.ArticleHTML = "<p>secret text</p>"

	build := func() []byte {
		t.Helper()
		index, err := Build([]*models.Post{private, public}, BuildOptions{
			ContentDir:          root,
			AssetsDir:           assets,
			GeneratorVersion:    "test",
			IncludeUnreferenced: true,
		})
		if err != nil {
			t.Fatalf("Build() error = %v", err)
		}
		data, err := Marshal(index)
		if err != nil {
			t.Fatalf("Marshal() error = %v", err)
		}
		return data
	}

	before := build()
	private.Content = `![classified markdown](/unreferenced.png)

<figure><img src="/private-html.png" alt="classified html"><figcaption>classified caption</figcaption></figure>

![remote](https://cdn.example.test/private.jpg?token=SECRET)`
	private.ArticleHTML = `<figure><img src="/private-rendered.png" alt="rendered secret"><figcaption>rendered caption</figcaption></figure>`
	after := build()
	if !bytes.Equal(before, after) {
		t.Fatalf("private body changed public inventory:\nbefore=%s\nafter=%s", before, after)
	}

	index, err := Build([]*models.Post{private, public}, BuildOptions{
		ContentDir:          root,
		AssetsDir:           assets,
		GeneratorVersion:    "test",
		IncludeUnreferenced: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if image, ok := findImage(index, "/unreferenced.png"); !ok || len(image.Uses) != 0 {
		t.Fatalf("private body affected unreferenced asset: found=%v image=%#v", ok, image)
	}
	for _, source := range []string{"/private-html.png", "/private-rendered.png", "https://cdn.example.test/private.jpg?token=SECRET"} {
		if _, ok := findImage(index, source); ok {
			t.Fatalf("private body source %q leaked into inventory", source)
		}
	}
}

func TestBuild_PrivateSafeCoverFrontmatterAddsMetadataWithoutUse(t *testing.T) {
	root := t.TempDir()
	assets := filepath.Join(root, "static")
	if err := os.MkdirAll(assets, 0o755); err != nil {
		t.Fatal(err)
	}
	writePNG(t, filepath.Join(assets, "cover.png"), 4, 3)
	writePNG(t, filepath.Join(assets, "replacement.png"), 5, 2)

	private := models.NewPost("posts/private.md")
	private.Path = "posts/private.md"
	private.Slug = "private"
	private.Href = "/private/"
	private.Published = true
	private.Private = true
	private.Content = "private body"
	private.Extra["cover"] = "/cover.png"
	private.Extra["cover_alt"] = "Public-safe cover"
	private.Extra["image"] = "/private-image.png"
	private.Extra["alt"] = "private alt that is not whitelisted"

	index, err := Build([]*models.Post{private}, BuildOptions{ContentDir: root, AssetsDir: assets})
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}
	cover := imageBySrc(t, index, "/cover.png")
	if !cover.Cover || cover.Alt != "Public-safe cover" || len(cover.Uses) != 0 {
		t.Fatalf("private safe cover metadata = %#v", cover)
	}
	if _, ok := findImage(index, "/private-image.png"); ok {
		t.Fatal("non-whitelisted private frontmatter image leaked")
	}
	data, err := Marshal(index)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Contains(data, []byte("posts/private.md")) || bytes.Contains(data, []byte("/private/")) || bytes.Contains(data, []byte("private body")) {
		t.Fatalf("private relationship or body leaked: %s", data)
	}

	private.Extra["cover"] = "/replacement.png"
	index, err = Build([]*models.Post{private}, BuildOptions{ContentDir: root, AssetsDir: assets})
	if err != nil {
		t.Fatalf("Build() after safe frontmatter change error = %v", err)
	}
	if _, ok := findImage(index, "/cover.png"); ok {
		t.Fatal("old safe cover remained after frontmatter change")
	}
	if replacement, ok := findImage(index, "/replacement.png"); !ok || len(replacement.Uses) != 0 || !replacement.Cover {
		t.Fatalf("new safe cover metadata = found=%v image=%#v", ok, replacement)
	}
}

func TestBuild_CanonicalizesLocalURLPathAndDeduplicatesEncodedReference(t *testing.T) {
	root := t.TempDir()
	assets := filepath.Join(root, "static")
	if err := os.MkdirAll(assets, 0o755); err != nil {
		t.Fatal(err)
	}
	filename := "my photo#1%.png"
	writePNG(t, filepath.Join(assets, filename), 2, 2)

	post := models.NewPost("posts/photo.md")
	post.Path = "posts/photo.md"
	post.Slug = "photo"
	post.Href = "/photo/"
	post.Published = true
	post.Content = "![Raw](/my photo#1%.png)\n![Encoded](/my%20photo%231%25.png)"
	index, err := Build([]*models.Post{post}, BuildOptions{ContentDir: root, AssetsDir: assets})
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}
	wantSource := "/my%20photo%231%25.png"
	image := imageBySrc(t, index, wantSource)
	if len(index.Images) != 1 || len(image.Uses) != 1 {
		t.Fatalf("canonical image records = %#v", index.Images)
	}
	if strings.Contains(image.Src, " ") || strings.Contains(image.Src, "#") || strings.Contains(image.Src, "%25%20") {
		t.Fatalf("canonical source is not a valid escaped path: %q", image.Src)
	}

	post.Content = ""
	post.Extra["cover"] = filename
	rawOnly, err := Build([]*models.Post{post}, BuildOptions{ContentDir: root, AssetsDir: assets})
	if err != nil {
		t.Fatalf("raw-only Build() error = %v", err)
	}
	if len(rawOnly.Images) != 1 || rawOnly.Images[0].Src != wantSource {
		t.Fatalf("raw-only canonical image = %#v", rawOnly.Images)
	}
}

func TestBuild_TreatsLiteralQuestionMarkInLocalFilenameAsPath(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Windows does not support question marks in filenames")
	}

	root := t.TempDir()
	assets := filepath.Join(root, "static")
	if err := os.MkdirAll(assets, 0o755); err != nil {
		t.Fatal(err)
	}
	writePNG(t, filepath.Join(assets, "question?mark.png"), 2, 2)

	post := models.NewPost("posts/question.md")
	post.Path = "posts/question.md"
	post.Slug = "question"
	post.Href = "/question/"
	post.Published = true
	post.Content = "![Question](/question?mark.png)"

	index, err := Build([]*models.Post{post}, BuildOptions{ContentDir: root, AssetsDir: assets})
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}
	if len(index.Images) != 1 || index.Images[0].Src != "/question%3Fmark.png" {
		t.Fatalf("literal question-mark image = %#v", index.Images)
	}
	post.Content = ""
	post.Extra["cover"] = "question?mark.png"
	frontmatter, err := Build([]*models.Post{post}, BuildOptions{ContentDir: root, AssetsDir: assets})
	if err != nil {
		t.Fatalf("frontmatter Build() error = %v", err)
	}
	if len(frontmatter.Images) != 1 || frontmatter.Images[0].Src != "/question%3Fmark.png" {
		t.Fatalf("frontmatter question-mark image = %#v", frontmatter.Images)
	}
}

func TestBuild_ResolvesWindowsAbsoluteLocalReference(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Windows drive-letter paths require a Windows filesystem")
	}

	root := t.TempDir()
	assets := filepath.Join(root, "static")
	if err := os.MkdirAll(assets, 0o755); err != nil {
		t.Fatal(err)
	}
	imagePath := filepath.Join(assets, "photo.png")
	writePNG(t, imagePath, 3, 2)

	post := models.NewPost("posts/photo.md")
	post.Path = "posts/photo.md"
	post.Slug = "photo"
	post.Href = "/photo/"
	post.Published = true
	post.Extra["cover"] = filepath.ToSlash(imagePath)

	index, err := Build([]*models.Post{post}, BuildOptions{
		ContentDir:          root,
		AssetsDir:           assets,
		IncludeUnreferenced: true,
	})
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}
	if len(index.Images) != 1 {
		t.Fatalf("absolute Windows reference created duplicate records: %#v", index.Images)
	}
	image := imageBySrc(t, index, "/photo.png")
	if image.Width != 3 || image.Height != 2 || !image.Cover || len(image.Uses) != 1 {
		t.Fatalf("absolute Windows reference = %#v", image)
	}
}

func TestMediaExtension_RecognizesWindowsDrivePaths(t *testing.T) {
	for _, path := range []string{`C:\site\static\photo.png`, `D:/site/static/photo.webp`} {
		if got := mediaExtension(path); got != filepath.Ext(path) {
			t.Errorf("mediaExtension(%q) = %q, want %q", path, got, filepath.Ext(path))
		}
	}
}

func TestMediaExtension_UsesURLPathBeforeQuery(t *testing.T) {
	want := filepath.Ext("/site/photo.png")
	if got := mediaExtension(`/site/photo.png?size=large`); got != want {
		t.Fatalf("mediaExtension() = %q, want %q", got, want)
	}
}

func TestBuild_OnlyReferencedImagesWhenUnreferencedDisabled(t *testing.T) {
	root := t.TempDir()
	assets := filepath.Join(root, "static")
	if err := os.MkdirAll(assets, 0o755); err != nil {
		t.Fatal(err)
	}
	writePNG(t, filepath.Join(assets, "only-static.png"), 1, 1)

	post := models.NewPost("posts/one.md")
	post.Path = "posts/one.md"
	post.Slug = "one"
	post.Published = true
	post.Content = "![Referenced](https://example.test/remote.jpg)"

	index, err := Build([]*models.Post{post}, BuildOptions{ContentDir: root, AssetsDir: assets})
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}
	if index.ImageCount != 1 || index.Images[0].Src != "https://example.test/remote.jpg" {
		t.Fatalf("Build() = %#v", index.Images)
	}
}

func TestBuild_ScansUnreferencedVideoAssets(t *testing.T) {
	root := t.TempDir()
	assets := filepath.Join(root, "static", "videos")
	if err := os.MkdirAll(assets, 0o755); err != nil {
		t.Fatal(err)
	}
	videoPath := filepath.Join(assets, "clip.mp4")
	if err := os.WriteFile(videoPath, []byte("test video"), 0o600); err != nil {
		t.Fatal(err)
	}

	index, err := Build(nil, BuildOptions{
		ContentDir:          root,
		AssetsDir:           filepath.Join(root, "static"),
		IncludeUnreferenced: true,
	})
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}
	video := imageBySrc(t, index, "/videos/clip.mp4")
	if video.MIMEType != "video/mp4" || video.Cover || len(video.Uses) != 0 {
		t.Fatalf("unreferenced video = %#v", video)
	}
}

func TestBuild_ReportsDimensionsForSupportedLocalFormats(t *testing.T) {
	root := t.TempDir()
	assets := filepath.Join(root, "static")
	if err := os.MkdirAll(assets, 0o755); err != nil {
		t.Fatal(err)
	}

	writeWebP(t, filepath.Join(assets, "web.webp"), 17, 11)
	writeBMP(t, filepath.Join(assets, "bitmap.bmp"), 13, 7)
	writeTIFF(t, filepath.Join(assets, "scan.tiff"), 9, 5)
	if err := os.WriteFile(filepath.Join(assets, "vector.svg"), []byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 120 80"></svg>`), 0o600); err != nil {
		t.Fatal(err)
	}
	writeICO(t, filepath.Join(assets, "icon.ico"), 64, 32)
	if err := os.WriteFile(filepath.Join(assets, "photo.avif"), []byte("ftypavif without a decoder"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(assets, "photo.heic"), []byte("ftypheic without a decoder"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(assets, "photo.heif"), []byte("ftypheif without a decoder"), 0o600); err != nil {
		t.Fatal(err)
	}

	index, err := Build(nil, BuildOptions{ContentDir: root, AssetsDir: assets, IncludeUnreferenced: true})
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}
	for _, test := range []struct {
		src      string
		width    int
		height   int
		mimeType string
	}{
		{src: "/web.webp", width: 17, height: 11, mimeType: "image/webp"},
		{src: "/bitmap.bmp", width: 13, height: 7, mimeType: "image/bmp"},
		{src: "/scan.tiff", width: 9, height: 5, mimeType: "image/tiff"},
		{src: "/vector.svg", width: 120, height: 80, mimeType: "image/svg+xml"},
		{src: "/icon.ico", width: 64, height: 32, mimeType: "image/x-icon"},
		{src: "/photo.avif", mimeType: "image/avif"},
		{src: "/photo.heic", mimeType: "image/heic"},
		{src: "/photo.heif", mimeType: "image/heif"},
	} {
		image := imageBySrc(t, index, test.src)
		if image.Width != test.width || image.Height != test.height || image.MIMEType != test.mimeType {
			t.Errorf("%s metadata = (%d, %d, %q), want (%d, %d, %q)", test.src, image.Width, image.Height, image.MIMEType, test.width, test.height, test.mimeType)
		}
	}
}

func TestSVGIntrinsicDimensions_UsesAbsoluteValuesAndViewBox(t *testing.T) {
	tests := []struct {
		name                  string
		widthValue            string
		heightValue           string
		viewBox               string
		wantWidth, wantHeight int
	}{
		{name: "pixel dimensions", widthValue: "120px", heightValue: "80px", wantWidth: 120, wantHeight: 80},
		{name: "viewBox fallback", viewBox: "0 0 200 100", wantWidth: 200, wantHeight: 100},
		{name: "percentages use viewBox", widthValue: "100%", heightValue: "50%", viewBox: "0 0 200 100", wantWidth: 200, wantHeight: 100},
		{name: "malformed", widthValue: "wide", heightValue: "80px", wantWidth: 0, wantHeight: 0},
		{name: "overflow", widthValue: "9223372036854775807px", heightValue: "80px", wantWidth: 0, wantHeight: 0},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			width, height := svgIntrinsicDimensions(test.widthValue, test.heightValue, test.viewBox)
			if width != test.wantWidth || height != test.wantHeight {
				t.Fatalf("svgIntrinsicDimensions() = (%d, %d), want (%d, %d)", width, height, test.wantWidth, test.wantHeight)
			}
		})
	}
}

func TestBuild_ClassifiesCoversEmbedsAndVideos(t *testing.T) {
	root := t.TempDir()
	post := models.NewPost(filepath.Join(root, "posts", "media.md"))
	post.Path = "posts/media.md"
	post.Slug = "media"
	post.Href = "/media/"
	post.Published = true
	post.Title = stringPointer("Media")
	postDate := time.Date(2026, time.January, 15, 12, 0, 0, 0, time.FixedZone("test", 2*60*60))
	post.Date = &postDate
	post.Extra["image"] = "https://dropper.wayl.one/file/clip.mp4?token=keep"
	post.Extra["poster_image"] = "http://dropper.wayl.one/file/clip.webp"
	post.Content = `![Inline video](https://dropper.wayl.one/file/inline.mp4)

<div class="custom-external-card" data-markata-embed="true"><img src="https://example.test/og-image.jpg" alt="External OG image"></div>`
	post.ArticleHTML = `<video poster="https://dropper.wayl.one/file/body.webp"><source src="https://dropper.wayl.one/file/body.webm" type="video/webm"></video>`

	index, err := Build([]*models.Post{post}, BuildOptions{ContentDir: root, GeneratorVersion: "test"})
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}
	if index.ImageCount != 4 {
		t.Fatalf("ImageCount = %d, want 4 (%#v)", index.ImageCount, index.Images)
	}

	coverVideo := imageBySrc(t, index, "https://dropper.wayl.one/file/clip.mp4?token=keep")
	if coverVideo.MIMEType != "video/mp4" || !coverVideo.Cover || coverVideo.PosterSrc != "https://dropper.wayl.one/file/clip.webp" || coverVideo.LastUsedAt == nil || !coverVideo.LastUsedAt.Equal(postDate.UTC()) {
		t.Fatalf("frontmatter video = %#v", coverVideo)
	}
	if len(coverVideo.Uses) != 1 || !coverVideo.Uses[0].Cover || coverVideo.Uses[0].Embed {
		t.Fatalf("frontmatter video uses = %#v", coverVideo.Uses)
	}

	bodyVideo := imageBySrc(t, index, "https://dropper.wayl.one/file/body.webm")
	if bodyVideo.MIMEType != "video/webm" || bodyVideo.PosterSrc != "https://dropper.wayl.one/file/clip.webp" {
		t.Fatalf("HTML video = %#v", bodyVideo)
	}

	inlineVideo := imageBySrc(t, index, "https://dropper.wayl.one/file/inline.mp4")
	if inlineVideo.MIMEType != "video/mp4" || inlineVideo.PosterSrc != "https://dropper.wayl.one/file/clip.webp" {
		t.Fatalf("Markdown video = %#v", inlineVideo)
	}

	embedImage := imageBySrc(t, index, "https://example.test/og-image.jpg")
	if !embedImage.Embed || embedImage.Cover || len(embedImage.Uses) != 1 || !embedImage.Uses[0].Embed {
		t.Fatalf("embedded OG image = %#v", embedImage)
	}
}

func TestBuild_DoesNotMarkInternalEmbedMediaAsExternal(t *testing.T) {
	post := models.NewPost("posts/internal.md")
	post.Path = "posts/internal.md"
	post.Slug = "internal"
	post.Href = "/internal/"
	post.Published = true
	post.Content = `<div class="embed-card"><img src="https://example.test/internal.jpg" alt="Internal image"></div>`

	index, err := Build([]*models.Post{post}, BuildOptions{GeneratorVersion: "test"})
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}
	image := imageBySrc(t, index, "https://example.test/internal.jpg")
	if image.Embed || len(image.Uses) != 1 || image.Uses[0].Embed {
		t.Fatalf("internal embed media = %#v", image)
	}
}

func TestBuild_ClassifiesDirectHTMLVideoWithPoster(t *testing.T) {
	post := models.NewPost("posts/direct-video.md")
	post.Path = "posts/direct-video.md"
	post.Slug = "direct-video"
	post.Href = "/direct-video/"
	post.Published = true
	post.ArticleHTML = `<video src="https://example.test/direct.mp4" type="video/mp4" poster="https://example.test/direct-poster.jpg"></video>`

	index, err := Build([]*models.Post{post}, BuildOptions{GeneratorVersion: "test"})
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}
	video := imageBySrc(t, index, "https://example.test/direct.mp4")
	if video.MIMEType != "video/mp4" || video.PosterSrc != "https://example.test/direct-poster.jpg" {
		t.Fatalf("direct HTML video = %#v", video)
	}
}

func TestBuild_RejectsCredentialBearingPosterURLs(t *testing.T) {
	post := models.NewPost("posts/private-poster.md")
	post.Path = "posts/private-poster.md"
	post.Slug = "private-poster"
	post.Href = "/private-poster/"
	post.Published = true
	post.Extra["poster_image"] = "https://user:SECRET@example.test/poster.webp"
	post.ArticleHTML = `<video src="https://example.test/clip.mp4"></video>`

	index, err := Build([]*models.Post{post}, BuildOptions{GeneratorVersion: "test"})
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}
	video := imageBySrc(t, index, "https://example.test/clip.mp4")
	data, err := Marshal(index)
	if err != nil {
		t.Fatal(err)
	}
	if video.PosterSrc != "" || bytes.Contains(data, []byte("SECRET")) {
		t.Fatalf("credential-bearing poster leaked: %#v", video)
	}
}

func TestBuild_CollectsFigureCaptionsForImagesAndVideos(t *testing.T) {
	post := models.NewPost("posts/captions.md")
	post.Path = "posts/captions.md"
	post.Slug = "captions"
	post.Href = "/captions/"
	post.Published = true
	post.Title = stringPointer("Caption post")
	post.Content = `![Markdown image](https://example.test/markdown.jpg)
![Second image](https://example.test/second.jpg)
Markdown figure caption`
	post.ArticleHTML = `<figure><a href="https://example.test/html.jpg"><img src="https://example.test/html.jpg" alt="HTML image"></a><figcaption><strong>HTML image</strong> caption</figcaption></figure>
<figure><video src="https://example.test/direct-clip.mp4"></video><figcaption>Direct video caption</figcaption></figure>
<figure><video><source src="https://example.test/clip.mp4" type="video/mp4"></video><figcaption>Video <em>figure</em> caption</figcaption></figure>`

	index, err := Build([]*models.Post{post}, BuildOptions{GeneratorVersion: "test"})
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}

	for _, test := range []struct {
		src     string
		caption string
	}{
		{src: "https://example.test/markdown.jpg", caption: "Markdown figure caption"},
		{src: "https://example.test/second.jpg", caption: "Markdown figure caption"},
		{src: "https://example.test/html.jpg", caption: "HTML image caption"},
		{src: "https://example.test/direct-clip.mp4", caption: "Direct video caption"},
		{src: "https://example.test/clip.mp4", caption: "Video figure caption"},
	} {
		image := imageBySrc(t, index, test.src)
		if len(image.Uses) != 1 || image.Uses[0].Caption != test.caption {
			t.Fatalf("caption for %q = %#v, want %q", test.src, image.Uses, test.caption)
		}
	}
}

func TestBuild_TracksLatestPublicUse(t *testing.T) {
	olderDate := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)
	newerDate := time.Date(2026, time.February, 1, 0, 0, 0, 0, time.UTC)
	older := models.NewPost("posts/older.md")
	older.Path = "posts/older.md"
	older.Slug = "older"
	older.Href = "/older/"
	older.Published = true
	older.Date = &olderDate
	older.Content = "![Shared](https://example.test/shared.jpg)"
	newer := models.NewPost("posts/newer.md")
	newer.Path = "posts/newer.md"
	newer.Slug = "newer"
	newer.Href = "/newer/"
	newer.Published = true
	newer.Date = &newerDate
	newer.Content = "![Shared](https://example.test/shared.jpg)"

	index, err := Build([]*models.Post{older, newer}, BuildOptions{GeneratorVersion: "test"})
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}
	shared := imageBySrc(t, index, "https://example.test/shared.jpg")
	if shared.LastUsedAt == nil || !shared.LastUsedAt.Equal(newerDate) {
		t.Fatalf("latest use = %#v, want %s", shared.LastUsedAt, newerDate.Format(time.RFC3339))
	}
}

func TestBuild_TracksEarliestDatedPublicUse(t *testing.T) {
	oldestDate := time.Date(2024, time.March, 10, 0, 0, 0, 0, time.UTC)
	newestDate := time.Date(2026, time.September, 12, 0, 0, 0, 0, time.UTC)
	sharedSource := "https://example.test/shared.webp"

	oldest := models.NewPost("posts/oldest.md")
	oldest.Path = "posts/oldest.md"
	oldest.Slug = "oldest"
	oldest.Href = "/oldest/"
	oldest.Published = true
	oldest.Date = &oldestDate
	oldest.Content = "![Shared](" + sharedSource + ")"

	newest := models.NewPost("posts/newest.md")
	newest.Path = "posts/newest.md"
	newest.Slug = "newest"
	newest.Href = "/newest/"
	newest.Published = true
	newest.Date = &newestDate
	newest.Content = "![Shared](" + sharedSource + ")"

	undated := models.NewPost("posts/undated.md")
	undated.Path = "posts/undated.md"
	undated.Slug = "undated"
	undated.Href = "/undated/"
	undated.Published = true
	undated.Content = "![Shared](" + sharedSource + ")"

	privateDate := time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)
	private := models.NewPost("posts/private.md")
	private.Path = "posts/private.md"
	private.Slug = "private"
	private.Href = "/private/"
	private.Published = true
	private.Private = true
	private.Date = &privateDate
	private.Content = "![Shared](" + sharedSource + ")"

	index, err := Build([]*models.Post{private, newest, undated, oldest}, BuildOptions{GeneratorVersion: "test"})
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}
	shared := imageBySrc(t, index, sharedSource)
	if shared.AddedAt == nil || !shared.AddedAt.Equal(oldestDate) {
		t.Fatalf("earliest public use = %#v, want %s", shared.AddedAt, oldestDate.Format(time.RFC3339))
	}
	if shared.LastUsedAt == nil || !shared.LastUsedAt.Equal(newestDate) {
		t.Fatalf("latest public use = %#v, want %s", shared.LastUsedAt, newestDate.Format(time.RFC3339))
	}

	undatedOnly := models.NewPost("posts/undated-only.md")
	undatedOnly.Path = "posts/undated-only.md"
	undatedOnly.Slug = "undated-only"
	undatedOnly.Href = "/undated-only/"
	undatedOnly.Published = true
	undatedOnly.Content = "![Undated](https://example.test/undated.webp)"
	index, err = Build([]*models.Post{undatedOnly}, BuildOptions{GeneratorVersion: "test"})
	if err != nil {
		t.Fatalf("undated Build() error = %v", err)
	}
	undatedImage := imageBySrc(t, index, "https://example.test/undated.webp")
	if undatedImage.AddedAt != nil || undatedImage.LastUsedAt != nil {
		t.Fatalf("undated public use invented metadata: %#v", undatedImage)
	}
}

func TestBuild_DoesNotUseMediaMtimeAsMetadata(t *testing.T) {
	root := t.TempDir()
	assets := filepath.Join(root, "static")
	if err := os.MkdirAll(filepath.Join(assets, "images"), 0o755); err != nil {
		t.Fatal(err)
	}
	imagePath := filepath.Join(assets, "images", "photo.png")
	writePNG(t, imagePath, 7, 5)

	postDate := time.Date(2024, time.March, 10, 0, 0, 0, 0, time.UTC)
	post := models.NewPost("posts/photo.md")
	post.Path = "posts/photo.md"
	post.Slug = "photo"
	post.Href = "/photo/"
	post.Published = true
	post.Date = &postDate
	post.Content = "![Photo](/images/photo.png)"

	build := func() []byte {
		t.Helper()
		index, err := Build([]*models.Post{post}, BuildOptions{ContentDir: root, AssetsDir: assets, GeneratorVersion: "test", IncludeUnreferenced: true})
		if err != nil {
			t.Fatalf("Build() error = %v", err)
		}
		image := imageBySrc(t, index, "/images/photo.png")
		if image.AddedAt == nil || !image.AddedAt.Equal(postDate) {
			t.Fatalf("added_at = %#v, want post date %s", image.AddedAt, postDate.Format(time.RFC3339))
		}
		data, err := Marshal(index)
		if err != nil {
			t.Fatalf("Marshal() error = %v", err)
		}
		return data
	}

	first := build()
	mtime := time.Date(1999, time.December, 31, 23, 59, 0, 0, time.UTC)
	if err := os.Chtimes(imagePath, mtime, mtime); err != nil {
		t.Fatal(err)
	}
	second := build()
	if !bytes.Equal(first, second) {
		t.Fatalf("media mtime changed generated output:\nfirst=%s\nsecond=%s", first, second)
	}
}

func TestBuild_IsDeterministicAcrossPostOrderAndLocalAliases(t *testing.T) {
	root := t.TempDir()
	assets := filepath.Join(root, "static", "images")
	if err := os.MkdirAll(assets, 0o755); err != nil {
		t.Fatal(err)
	}
	writePNG(t, filepath.Join(assets, "photo.png"), 2, 2)

	first := models.NewPost("posts/a.md")
	first.Path = "posts/a.md"
	first.Slug = "a"
	first.Href = "/a/"
	first.Published = true
	first.Content = "![First alt](/static/images/photo.png)"
	first.Title = stringPointer("A")

	second := models.NewPost("posts/b.md")
	second.Path = "posts/b.md"
	second.Slug = "b"
	second.Href = "/b/"
	second.Published = true
	second.Content = "![Second alt](/images/photo.png)"
	second.Title = stringPointer("B")

	build := func(posts []*models.Post) []byte {
		t.Helper()
		index, err := Build(posts, BuildOptions{ContentDir: root, AssetsDir: filepath.Join(root, "static")})
		if err != nil {
			t.Fatalf("Build() error = %v", err)
		}
		data, err := json.Marshal(index)
		if err != nil {
			t.Fatalf("marshal index: %v", err)
		}
		return data
	}

	forward := build([]*models.Post{first, second})
	reverse := build([]*models.Post{second, first})
	if !bytes.Equal(forward, reverse) {
		t.Fatalf("post ordering changed output:\nforward=%s\nreverse=%s", forward, reverse)
	}
	index, err := Build([]*models.Post{second, first}, BuildOptions{ContentDir: root, AssetsDir: filepath.Join(root, "static")})
	if err != nil {
		t.Fatal(err)
	}
	if len(index.Images) != 1 || index.Images[0].Src != "/images/photo.png" || index.Images[0].Alt != "First alt" {
		t.Fatalf("canonical image = %#v", index.Images)
	}
}

func TestBuild_PrivateSymlinkBodyDoesNotSuppressPublicTarget(t *testing.T) {
	root := t.TempDir()
	assets := filepath.Join(root, "static", "images")
	if err := os.MkdirAll(assets, 0o755); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(assets, "private-target.png")
	alias := filepath.Join(assets, "private-alias.png")
	writePNG(t, target, 2, 2)
	if err := os.Symlink("private-target.png", alias); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}

	private := models.NewPost(filepath.Join(root, "posts", "private.md"))
	private.Path = "posts/private.md"
	private.Published = true
	private.Private = true
	private.Content = "![Secret](images/private-alias.png)"

	index, err := Build([]*models.Post{private}, BuildOptions{
		ContentDir:          root,
		AssetsDir:           filepath.Join(root, "static"),
		IncludeUnreferenced: true,
	})
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}
	if index.ImageCount != 1 || index.Images[0].Src != "/images/private-target.png" || len(index.Images[0].Uses) != 0 {
		t.Fatalf("private symlink body affected unreferenced target: %#v", index.Images)
	}
}

func writePNG(t *testing.T, path string, width, height int) {
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

func writeWebP(t *testing.T, path string, width, height int) {
	t.Helper()
	if width < 1 || height < 1 || width > 16384 || height > 16384 {
		t.Fatalf("unsupported test WebP dimensions: %dx%d", width, height)
	}
	widthMinusOne, heightMinusOne := uint32(width-1), uint32(height-1) //nolint:gosec // dimensions are bounded above.
	payload := []byte{
		0x2f,
		byte(widthMinusOne),
		byte(widthMinusOne>>8) | byte(heightMinusOne<<6),
		byte(heightMinusOne >> 2),
		byte(heightMinusOne >> 10),
	}
	data := make([]byte, 0, 12+8+len(payload)+1)
	data = append(data, []byte("RIFF")...)
	data = append(data, 0, 0, 0, 0)
	data = append(data, []byte("WEBPVP8L")...)
	chunkLength := make([]byte, 4)
	binary.LittleEndian.PutUint32(chunkLength, uint32(len(payload))) //nolint:gosec // test payload is bounded.
	data = append(data, chunkLength...)
	data = append(data, payload...)
	data = append(data, 0)                                        // RIFF chunks are padded to an even length.
	binary.LittleEndian.PutUint32(data[4:8], uint32(len(data)-8)) //nolint:gosec // test payload is bounded.
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
}

func writeBMP(t *testing.T, path string, width, height int) {
	t.Helper()
	rowSize := (width*3 + 3) &^ 3
	data := make([]byte, 54+rowSize*height)
	copy(data, "BM")
	binary.LittleEndian.PutUint32(data[2:6], uint32(len(data))) //nolint:gosec // test bitmap is bounded by its fixture dimensions.
	binary.LittleEndian.PutUint32(data[10:14], 54)
	binary.LittleEndian.PutUint32(data[14:18], 40)
	binary.LittleEndian.PutUint32(data[18:22], uint32(width))  //nolint:gosec // test dimensions are controlled positive values.
	binary.LittleEndian.PutUint32(data[22:26], uint32(height)) //nolint:gosec // test dimensions are controlled positive values.
	binary.LittleEndian.PutUint16(data[26:28], 1)
	binary.LittleEndian.PutUint16(data[28:30], 24)
	binary.LittleEndian.PutUint32(data[34:38], uint32(rowSize*height)) //nolint:gosec // test bitmap is bounded by its fixture dimensions.
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
}

func writeTIFF(t *testing.T, path string, width, height int) {
	t.Helper()
	file, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	canvas := image.NewRGBA(image.Rect(0, 0, width, height))
	if err := tiff.Encode(file, canvas, nil); err != nil {
		_ = file.Close()
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}
}

func writeICO(t *testing.T, path string, width, height int) {
	t.Helper()
	data := make([]byte, 22+40)
	binary.LittleEndian.PutUint16(data[2:4], 1)
	binary.LittleEndian.PutUint16(data[4:6], 1)
	data[6] = byte(width % 256)
	data[7] = byte(height % 256)
	binary.LittleEndian.PutUint16(data[10:12], 1)
	binary.LittleEndian.PutUint16(data[12:14], 32)
	binary.LittleEndian.PutUint32(data[14:18], 40)
	binary.LittleEndian.PutUint32(data[18:22], 22)
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
}

func imageBySrc(t *testing.T, index Index, src string) Image {
	t.Helper()
	image, ok := findImage(index, src)
	if !ok {
		t.Fatalf("image %q not found in %#v", src, index.Images)
	}
	return image
}

func findImage(index Index, src string) (Image, bool) {
	for i := range index.Images {
		image := &index.Images[i]
		if image.Src == src {
			return *image, true
		}
	}
	return Image{}, false
}

func stringPointer(value string) *string {
	return &value
}
