package themes

import (
	"os"
	"strings"
	"testing"
)

func TestPostTemplatesLabelHiddenCanonicalURL(t *testing.T) {
	// The root template is used by local development; the embedded default
	// theme is shipped in the binary. Both must avoid empty-link warnings.
	root, err := os.ReadFile("../../templates/post.html")
	if err != nil {
		t.Fatalf("read root post template: %v", err)
	}
	bundled, err := ReadTemplate("post.html")
	if err != nil {
		t.Fatalf("read bundled post template: %v", err)
	}

	const labeled = `<a class="u-url" href="{{ config.url }}{{ post.href }}" hidden>Canonical URL for this post</a>`
	const empty = `<a class="u-url" href="{{ config.url }}{{ post.href }}" hidden></a>`
	for name, source := range map[string]string{
		"root":    string(root),
		"bundled": string(bundled),
	} {
		t.Run(name, func(t *testing.T) {
			if strings.Contains(source, empty) {
				t.Fatal("hidden microformats canonical URL link has no accessible text")
			}
			if strings.Count(source, labeled) != 1 {
				t.Fatalf("expected exactly one labeled, hidden canonical u-url anchor; found %d", strings.Count(source, labeled))
			}
		})
	}
}
