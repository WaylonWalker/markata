package themes

import (
	"io/fs"
	"strings"
	"testing"
)

// The visual Markdown components must retain a usable heading scale. A past
// blanket H2-H6 rule flattened every card/layout heading to the same size,
// making an otherwise well-structured visual article feel typographically flat.
func TestVisualMarkdownContainerHeadingHierarchy(t *testing.T) {
	data, err := fs.ReadFile(DefaultStatic(), "css/components.css")
	if err != nil {
		t.Fatal(err)
	}
	css := string(data)
	if strings.Contains(css, `:is(h2, h3, h4, h5, h6) {
    font-size: var(--text-lg`) {
		t.Fatal("nested H2-H6 must not all share one font size")
	}
	for _, want := range []string{
		`:is(.md-layout, .md-box, .md-card) > h2`,
		`:is(.md-layout, .md-box, .md-card) > h3`,
		`:is(.md-layout, .md-box, .md-card) > h4`,
		`:is(.md-layout, .md-box, .md-card) > :is(h5, h6)`,
		`font-size: var(--text-2xl, 1.5rem)`,
		`font-size: var(--text-xl, 1.25rem)`,
		`font-size: var(--text-lg, 1.125rem)`,
		`color: var(--color-text-muted, var(--color-text))`,
		`.post-content .md-card > p`,
	} {
		if !strings.Contains(css, want) {
			t.Errorf("visual typography is missing %q", want)
		}
	}
}
