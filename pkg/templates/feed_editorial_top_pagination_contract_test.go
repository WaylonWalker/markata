package templates

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEditorialFeedsExposeTopAndBottomManualPagination(t *testing.T) {
	for _, root := range []string{"../../templates", "../themes/default/templates"} {
		data, err := os.ReadFile(filepath.Join(root, "feed.html"))
		if err != nil {
			t.Fatal(err)
		}
		html := string(data)
		for _, expected := range []string{
			`feed_presentation == "editorial" and page.total_pages > 1 and page.pagination_type == "manual"`,
			`class="pagination pagination--top"`,
			`aria-label="Archive top pagination"`,
			`href="{{ page.prev_url }}"`,
			`href="{{ page.next_url }}"`,
			`Page {{ page.number }} of {{ page.total_pages }}`,
			`{% include "partials/pagination.html" %}`,
		} {
			if !strings.Contains(html, expected) {
				t.Errorf("%s missing pagination element %q", root, expected)
			}
		}
		if top, gallery := strings.Index(html, "pagination--top"), strings.Index(html, `{% include "partials/feed-editorial.html" %}`); top < 0 || gallery < 0 || top > gallery {
			t.Errorf("%s must render the navigation before the editorial reading content", root)
		}
	}
}
