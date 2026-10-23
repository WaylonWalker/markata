package themes

import (
	"os"
	"strings"
	"testing"
)

// Pagefind accepts metadata outside data-pagefind-body, including <head>.
// HTML <meta> elements inside an article are invalid and trigger Squirrel.
func TestPostPagefindImageMetadataIsInHead(t *testing.T) {
	root, err := os.ReadFile("../../templates/post.html")
	if err != nil {
		t.Fatal(err)
	}
	bundled, err := ReadTemplate("post.html")
	if err != nil {
		t.Fatal(err)
	}

	const marker = `data-pagefind-meta="image[content]"`
	for name, source := range map[string]string{
		"root":    string(root),
		"bundled": string(bundled),
	} {
		t.Run(name, func(t *testing.T) {
			metaStart := strings.Index(source, "{% block meta %}")
			if metaStart < 0 {
				t.Fatal("could not find head metadata template block")
			}
			metaEndOffset := strings.Index(source[metaStart:], "{% endblock %}")
			bodyStart := strings.Index(source, "{% block content %}")
			if metaEndOffset < 0 || bodyStart <= metaStart+metaEndOffset {
				t.Fatal("could not identify head metadata and content template blocks")
			}
			head := source[metaStart:metaStart+metaEndOffset]
			body := source[bodyStart:]
			if count := strings.Count(head, marker); count != 4 {
				t.Fatalf("expected four prioritized image sources in head, found %d", count)
			}
			if strings.Contains(body, marker) {
				t.Fatal("Pagefind image metadata must not emit <meta> in the article body")
			}
			sources := []string{
				"post.cover_image | media_url",
				"post.Extra.image | media_url",
				"post.Extra.thumbnail | media_url",
				"config.seo.default_image | media_url",
			}
			last := -1
			for _, expression := range sources {
				at := strings.Index(head, expression)
				if at <= last {
					t.Fatalf("missing or out-of-order image fallback: %s", expression)
				}
				last = at
			}
		})
	}
}
