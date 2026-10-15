package themes

import (
	"regexp"
	"strings"
	"testing"
)

// TestFeedPhotoGeometry guards against shrink-to-fit responsive images. With
// width-descriptor srcset candidates, fit-content can use the density-adjusted
// intrinsic width (600w at 2x -> 300 CSS px) instead of the feed's intended
// media width. The media must also retain its natural aspect ratio.
func TestFeedPhotoGeometry(t *testing.T) {
	raw, err := ReadStatic("css/cards.css")
	if err != nil {
		t.Fatal(err)
	}
	css := strings.ReplaceAll(string(raw), "\r\n", "\n")

	rule := func(selector string) string {
		t.Helper()
		pattern := regexp.MustCompile(regexp.QuoteMeta(selector) + `\s*\{([^}]*)\}`)
		match := pattern.FindStringSubmatch(css)
		if match == nil {
			t.Fatalf("missing CSS rule %s", selector)
		}
		return match[1]
	}

	figure := rule(".feed .posts-list .photo-figure")
	if !strings.Contains(figure, "width: min(100%, 600px);") || strings.Contains(figure, "fit-content") {
		t.Errorf("feed photo figure should have stable responsive width, got %q", figure)
	}

	for _, selector := range []string{
		".feed .posts-list .photo-figure > a",
		".feed .posts-list .photo-figure picture",
	} {
		if !strings.Contains(rule(selector), "width: 100%;") {
			t.Errorf("%s must fill the stable photo figure width", selector)
		}
	}

	media := rule(".feed .posts-list .photo-figure img,\n.feed .posts-list .photo-figure video")
	for _, want := range []string{"width: 100%;", "max-width: 100%;", "height: auto;", "object-fit: contain;"} {
		if !strings.Contains(media, want) {
			t.Errorf("feed photo media must preserve natural geometry; missing %q", want)
		}
	}
}
