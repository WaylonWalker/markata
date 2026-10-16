package themes

import (
	"strings"
	"testing"
)

// HTMX manual and infinite pagination both target .posts-list. Editorial
// sections must preserve that wrapper while retaining the calendar view hook.
func TestEditorialFeedPreservesPaginationAndCalendarHooks(t *testing.T) {
	content, err := ReadTemplate("partials/feed-editorial.html")
	if err != nil {
		t.Fatal(err)
	}
	template := string(content)
	for _, want := range []string{
		`class="editorial-feed posts-list" id="posts-list"`,
		`data-calendar-primary`,
		`{% include "partials/cards/card-router.html" %}`,
	} {
		if !strings.Contains(template, want) {
			t.Errorf("editorial template missing %q", want)
		}
	}
}

func TestEditorialFeedKeepsMobileGridWithoutExternalCSS(t *testing.T) {
	content, err := ReadTemplate("feed.html")
	if err != nil {
		t.Fatal(err)
	}
	markup := string(content)
	for _, want := range []string{
		`{% if feed_presentation == "editorial" %}`,
		`id="editorial-gallery-critical"`,
		`grid-template-columns: repeat(3, minmax(0, 1fr))`,
		`aspect-ratio: 1; object-fit: cover;`,
	} {
		if !strings.Contains(markup, want) {
			t.Errorf("mobile editorial grid missing critical rule %q", want)
		}
	}
}
