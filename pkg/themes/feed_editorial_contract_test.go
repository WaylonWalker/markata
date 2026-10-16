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
