package workbench

import (
	"strings"
	"testing"
)

func TestRenderFrontmatterPreview_LinkUsesMetadata(t *testing.T) {
	source := `---
template: link
title: "A <strong>title</strong>"
url: "https://example.com/watch?v=1&source=editor"
image: "https://cdn.example.com/preview.jpg"
description: "A useful description."
site_name: Example
provider_name: Example Provider
type: video.other
author_name: A Creator
author_url: "https://example.com/creator"
---
`

	got := renderFrontmatterPreview(source, "A note about this link.")
	for _, want := range []string{
		`class="embed-card embed-card-external workbench-link-preview"`,
		`href="https://example.com/watch?v=1&amp;source=editor"`,
		`src="https://cdn.example.com/preview.jpg"`,
		`A &lt;strong&gt;title&lt;/strong&gt;`,
		`A useful description.`,
		`Example`,
		`video.other`,
		`A Creator`,
		`href="https://example.com/creator"`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("preview does not contain %q:\n%s", want, got)
		}
	}
	if strings.Contains(got, "<strong>") {
		t.Error("frontmatter title was not escaped")
	}
}

func TestRenderFrontmatterPreview_PhotoUsesRemoteImage(t *testing.T) {
	source := `---
template: shot
title: A photo
image: "https://dropper.example.com/photo.webp"
description: A caption
---
`

	got := renderFrontmatterPreview(source, "Add a caption.")
	for _, want := range []string{
		`<h1 class="p-name">A photo</h1>`,
		`<figure class="photo-figure workbench-frontmatter-image">`,
		`src="https://dropper.example.com/photo.webp"`,
		`<figcaption class="p-summary">A caption</figcaption>`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("preview does not contain %q:\n%s", want, got)
		}
	}
}

func TestRenderFrontmatterPreview_LocalPhotoDoesNotDuplicateSiteTreeImage(t *testing.T) {
	source := `---
template: photo
image: /images/photo.webp
---
`

	got := renderFrontmatterPreview(source, "![Post image](/workspace/image?path=images%2Fphoto.webp)\n\nCaption")
	if got != "" {
		t.Errorf("preview duplicated a site-tree image: %s", got)
	}
}

func TestRenderFrontmatterPreview_InvalidYAMLKeepsWarning(t *testing.T) {
	source := "---\ntitle: broken\nimage: \"https://example.com/photo.webp\"template: shot\n---\n"

	got := renderFrontmatterPreview(source, "Body remains available.")
	if !strings.Contains(got, "Frontmatter preview is unavailable") {
		t.Fatalf("preview warning missing: %s", got)
	}
}

func TestRenderFrontmatterPreview_RejectsUnsafeLinkDestination(t *testing.T) {
	source := "---\ntemplate: link\ntitle: unsafe\nurl: 'javascript:alert(1)'\n---\n"
	if got := renderFrontmatterPreview(source, "Body"); got != "" {
		t.Errorf("unsafe link destination was rendered: %s", got)
	}
}

func TestRenderFrontmatterPreview_AuthoredLeadingHeadingTakesPrecedence(t *testing.T) {
	source := "---\ntitle: Frontmatter title\n---\n"
	got := renderFrontmatterPreview(source, "# Different body title\n\nBody")
	if strings.Contains(got, "<h1") {
		t.Fatalf("frontmatter preview must not add a second H1: %s", got)
	}
}
