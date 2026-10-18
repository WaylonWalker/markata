package plugins

import (
	"io"
	"os"
	"strings"
)

// editorialHTMLMarkupValid is the HTML presentation contract. Checking output
// rather than only requested config makes silent template override/fallback
// regressions visible to authors and agents.
func editorialHTMLMarkupValid(html string) bool {
	return strings.Contains(html, `data-feed-presentation="editorial"`) &&
		strings.Contains(html, `class="editorial-feed posts-list"`) &&
		strings.Contains(html, `data-editorial-version="4"`)
}

// editorialHTMLPublished inspects only the beginning of the existing HTML
// output when deciding whether an incremental cached feed is reusable.
// It does not read the entire archive or alter its build cache keys.
func editorialHTMLPublished(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		return false
	}
	defer file.Close()

	const maxHeaderBytes = 512 << 10
	head, err := io.ReadAll(io.LimitReader(file, maxHeaderBytes))
	if err != nil {
		return false
	}
	return editorialHTMLMarkupValid(string(head))
}
