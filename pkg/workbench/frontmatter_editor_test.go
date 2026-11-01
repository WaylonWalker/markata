package workbench

import (
	"strings"
	"testing"
)

func TestFrontmatterEditorRoundTripsUnchangedSource(t *testing.T) {
	source := "---\n# authored header\ntitle: 'Original' # important\npublished: false\ncustom:\n  nested: [a, b]\n# keep this note\ntags: [one]\n---\n\n# Body\n\nA --- delimiter in prose.\n"
	edited, err := replaceFrontmatterProperty(source, "published", true, false)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(edited, "published: true\n") || !strings.Contains(edited, "title: 'Original' # important") {
		t.Fatalf("edit lost existing values: %q", edited)
	}
	if !strings.Contains(edited, "custom:\n  nested: [a, b]\n# keep this note\ntags: [one]") {
		t.Fatalf("unrelated custom metadata changed: %q", edited)
	}
	if !strings.HasSuffix(edited, "\n# Body\n\nA --- delimiter in prose.\n") {
		t.Fatalf("Markdown body was changed: %q", edited)
	}
}

func TestFrontmatterEditorAddsRemovesAndPreservesCRLF(t *testing.T) {
	source := "# Hello\r\n\r\nBody\r\n"
	added, err := replaceFrontmatterProperty(source, "tags", []string{"one", "two"}, false)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(added, "---\r\ntags:\r\n  - one\r\n  - two\r\n---\r\n") {
		t.Fatalf("frontmatter creation failed: %q", added)
	}
	if !strings.HasSuffix(added, source) {
		t.Fatal("body not preserved")
	}
	added, err = replaceFrontmatterProperty(added, "published", true, false)
	if err != nil {
		t.Fatal(err)
	}
	removed, err := replaceFrontmatterProperty(added, "tags", nil, true)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(removed, "tags:") || !strings.Contains(removed, "published: true") {
		t.Fatalf("field remove failed: %q", removed)
	}
}

func TestFrontmatterEditorRefusesMalformedAndUnsafeFields(t *testing.T) {
	for _, source := range []string{"---\ntitle: missing close", "---\nfoo: [\n---\ntext", "---\nfoo: 1\nfoo: 2\n---\n"} {
		if _, err := replaceFrontmatterProperty(source, "published", true, false); err == nil {
			t.Fatalf("accepted malformed source %q", source)
		}
	}
	for _, name := range []string{"", "a: b", "key\n---", "a.b"} {
		if _, err := replaceFrontmatterProperty("# hello", name, true, false); err == nil {
			t.Fatalf("accepted invalid field %q", name)
		}
	}
}
