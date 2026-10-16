package themes

import (
	"io/fs"
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

func TestEditorialArchiveHasDistinctDesktopAndMobileComposition(t *testing.T) {
	raw, err := fs.ReadFile(DefaultStatic(), "css/feeds.css")
	if err != nil {
		t.Fatal(err)
	}
	css := string(raw)
	for _, want := range []string{
		".page-wrapper.page-wrapper--feed-only:has(> .main-content > .feed.h-feed.feed--editorial)",
		".main-content > .feed.h-feed.feed--editorial",
		".feed.h-feed.feed--editorial > .editorial-feed.posts-list",
		".feed--editorial .editorial-week-entries",
		"grid-template-columns: repeat(2, minmax(0, 1fr));",
		".editorial-entry--quiet",
		".editorial-entry--featured",
		".editorial-shot-collection",
		"@media (min-width: 880px)",
		"grid-template-columns: minmax(180px, 22%) minmax(0, 1fr);",
		"@media (max-width: 760px)",
		"grid-template-columns: minmax(0, 1fr);",
		"grid-template-columns: repeat(3, minmax(0, 1fr));",
	} {
		if !strings.Contains(css, want) {
			t.Errorf("editorial layout missing %q", want)
		}
	}
	if strings.Contains(css, ".page-wrapper.page-wrapper--feed-only:has(> .main-content > .feed.h-feed):not(:has(.feed-sidebar") {
		t.Fatal("editorial sizing must not depend on whether sidebars are present")
	}
}
