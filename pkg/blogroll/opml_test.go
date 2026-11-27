package blogroll

import (
	"bytes"
	"strings"
	"testing"
)

func TestReadOPMLNestedFolders(t *testing.T) {
	input := `<?xml version="1.0"?>
<opml version="2.0"><head><title>Reader</title></head><body>
  <outline text="Tech">
    <outline text="Go">
      <outline text="Go &amp; RSS" type="rss" xmlUrl="https://example.test/rss?a=1&amp;b=2" htmlUrl="https://example.test" />
    </outline>
    <outline title="Explicit" category="News" xmlUrl="https://example.test/atom" />
  </outline>
  <outline text="Not a feed" />
</body></opml>`
	feeds, err := ReadOPML(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	if len(feeds) != 2 || feeds[0].Category != "Tech/Go" || feeds[0].Title != "Go & RSS" ||
		feeds[0].URL != "https://example.test/rss?a=1&b=2" || feeds[1].Category != "News" {
		t.Fatalf("unexpected feeds: %+v", feeds)
	}
}

func TestOPMLRoundTrip(t *testing.T) {
	original := []OPMLFeed{
		{URL: "https://example.test/feed", Title: "A & B", Description: "<description>", Category: "Tech", SiteURL: "https://example.test"},
		{URL: "https://example.test/other", Title: "Other"},
	}
	var out bytes.Buffer
	if err := WriteOPML(&out, "Reader", original); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), `version="2.0"`) {
		t.Fatalf("missing OPML version: %s", out.String())
	}
	got, err := ReadOPML(&out)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != len(original) || got[0] != original[0] || got[1] != original[1] {
		t.Fatalf("round trip mismatch: %+v", got)
	}
}

func TestReadOPMLRejectsInvalidDocuments(t *testing.T) {
	for _, input := range []string{
		`<html/>`,
		`<opml><body><outline xmlUrl="file:///etc/passwd"/></body></opml>`,
		`<opml><body><outline xmlUrl="https://user:secret@example.test/feed"/></body></opml>`,
		`<opml><body>`,
		strings.Repeat(" ", (4<<20)+1),
	} {
		if _, err := ReadOPML(strings.NewReader(input)); err == nil {
			t.Fatal("expected malformed, unsupported or oversized OPML to fail")
		}
	}
}

func TestWriteOPMLRejectsCredentials(t *testing.T) {
	var out bytes.Buffer
	err := WriteOPML(&out, "Reader", []OPMLFeed{{URL: "https://user:secret@example.test/feed", Title: "Private"}})
	if err == nil || out.Len() != 0 || strings.Contains(err.Error(), "secret") {
		t.Fatal("credential-bearing URL must be rejected before writing")
	}
}
