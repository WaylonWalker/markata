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
		`data-editorial-relative-week data-week-start="{{ section.Start }}"`,
		`{{ section.RelativeLabel }}`,
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

func TestEditorialLeadAndPhotoEssayTemplateContract(t *testing.T) {
	content, err := ReadTemplate("partials/feed-editorial.html")
	if err != nil {
		t.Fatal(err)
	}
	tmpl := string(content)
	for _, want := range []string{
		`editorial-entry--lead`,
		`Featured this week`,
		`{{ item.TypeLabel }}`,
		`{{ item.DateLabel }}`,
		`editorial-shot-grid--mosaic`,
		`Photo journal`,
		`data-editorial-version="4"`,
	} {
		if !strings.Contains(tmpl, want) {
			t.Errorf("editorial template missing %q", want)
		}
	}
	raw, err := fs.ReadFile(DefaultStatic(), "css/feeds.css")
	if err != nil {
		t.Fatal(err)
	}
	css := string(raw)
	for _, want := range []string{
		".feed--editorial .editorial-entry--lead",
		".feed--editorial .editorial-entry-kicker",
		".feed--editorial .editorial-shot-grid--mosaic",
		"@media (min-width: 900px)",
		"@media (max-width: 480px)",
		"-webkit-line-clamp: 3;",
		"-webkit-line-clamp: 5;",
	} {
		if !strings.Contains(css, want) {
			t.Errorf("editorial CSS missing %q", want)
		}
	}
}

func TestEditorialPhotoDatesAreOnlyAddedInsideWeeklyGallery(t *testing.T) {
	for _, name := range []string{"partials/cards/photo-card.html", "partials/feed-editorial.html"} {
		tmpl, err := ReadTemplate(name)
		if err != nil {
			t.Fatal(err)
		}
		switch name {
		case "partials/cards/photo-card.html":
			if !strings.Contains(string(tmpl), "post.editorial_photo_date") ||
				!strings.Contains(string(tmpl), "editorial-shot-date") {
				t.Fatalf("photo card missing optional dated caption")
			}
		case "partials/feed-editorial.html":
			if !strings.Contains(string(tmpl), "data-editorial-version=\"4\"") {
				t.Fatal("editorial HTML version must invalidate old gallery caches")
			}
		}
	}
}

func TestEditorialArchiveDefersOnlyNonCriticalSearchCSS(t *testing.T) {
	base, err := ReadTemplate("base.html")
	if err != nil {
		t.Fatal(err)
	}
	s := strings.ReplaceAll(string(base), "\r\n", "\n")
	for _, want := range []string{
		`{% if feed_presentation == "editorial" %}`,
		`'css/search.css' | theme_asset_hashed }}" media="print" onload="this.media='all'"`,
		`<noscript><link rel="stylesheet" href="{{ 'css/search.css' | theme_asset_hashed }}"></noscript>`,
		`'css/cards.css' | theme_asset_hashed`,
		`'css/feeds.css' | theme_asset_hashed`,
	} {
		if !strings.Contains(s, want) {
			t.Errorf("expected deferred search and preserved initial archive styles: missing %q", want)
		}
	}
	if strings.Contains(s, `{% if feed_presentation == "editorial" %}`+"\n"+`   <script src="{{ 'js/heading-highlights.js' | theme_asset_hashed }}" defer></script>`) {
		t.Fatal("editorial archives must not load the heading highlight controller")
	}
	if !strings.Contains(s, `{% if feed_presentation != "editorial" %}`+"\n"+`   <script src="{{ 'js/heading-highlights.js' | theme_asset_hashed }}" defer></script>`+"\n"+`   {% endif %}`) {
		t.Fatal("heading highlight controller must stay available outside editorial feeds")
	}
}

func TestEditorialArchiveDefersCalendarStylesAndLeadImages(t *testing.T) {
	feed, err := ReadTemplate("feed.html")
	if err != nil {
		t.Fatal(err)
	}
	markup := string(feed)
	for _, want := range []string{
		`{% if feed_presentation == "editorial" %}`,
		`'css/calendar-feed.css' | theme_asset_hashed }}" media="print" onload="this.media='all'"`,
		`<noscript><link rel="stylesheet" href="{{ 'css/calendar-feed.css' | theme_asset_hashed }}"></noscript>`,
		`'js/calendar-feed.js' | theme_asset_hashed`,
	} {
		if !strings.Contains(markup, want) {
			t.Errorf("calendar nonblocking contract missing %q", want)
		}
	}
	link, err := ReadTemplate("partials/cards/link-card.html")
	if err != nil {
		t.Fatal(err)
	}
	card := string(link)
	for _, want := range []string{
		`loading="lazy" decoding="async"`,
		`editorial_card_image:item.Lead`,
	} {
		if !strings.Contains(card, want) {
			t.Errorf("link thumbnail priority contract missing %q", want)
		}
	}
	for _, unwanted := range []string{
		`loading="eager"`,
		`fetchpriority="high"`,
	} {
		if strings.Contains(card, unwanted) {
			t.Errorf("editorial link thumbnails must not include %q", unwanted)
		}
	}
}

func TestEditorialRelativeHeadingsScriptIsDeferred(t *testing.T) {
	for _, path := range []string{"feed.html", "partials/feed-editorial.html"} {
		tmpl, err := ReadTemplate(path)
		if err != nil {
			t.Fatal(err)
		}
		if path == "feed.html" && !strings.Contains(string(tmpl),
			`'js/editorial-relative-weeks.js' | theme_asset_hashed }}" defer`) {
			t.Fatal("editorial week labels need a deferred hashed script")
		}
		if path != "feed.html" && !strings.Contains(string(tmpl),
			`data-editorial-relative-week data-week-start="{{ section.Start }}"`) {
			t.Fatal("editorial week start must be present for runtime labels")
		}
	}
	js, err := fs.ReadFile(DefaultStatic(), "js/editorial-relative-weeks.js")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(js), "htmx:afterSwap") ||
		!strings.Contains(string(js), "view-transition-complete") {
		t.Fatal("relative labels must refresh after partial navigation")
	}
}
