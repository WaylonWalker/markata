package plugins

import (
	"strings"
	"testing"
	"time"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

func TestGenerateAtom_UpdatedDeterministic(t *testing.T) {
	config := lifecycle.NewConfig()
	config.Extra = map[string]interface{}{"url": "https://example.com"}

	old := time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC)
	newer := time.Date(2024, 2, 2, 12, 0, 0, 0, time.UTC)

	feed := &lifecycle.Feed{
		Name:  "test",
		Title: "Test Feed",
		Path:  "test",
		Posts: []*models.Post{
			{Slug: "old", Href: "/old/", Date: &old},
			{Slug: "new", Href: "/new/", Date: &newer},
		},
	}

	atom, err := GenerateAtom(feed, config)
	if err != nil {
		t.Fatalf("GenerateAtom error: %v", err)
	}

	if !strings.Contains(atom, "<updated>2024-02-02T12:00:00Z</updated>") {
		t.Fatalf("expected updated to use latest post date, got:\n%s", atom)
	}
}

func TestGenerateRSS_LastBuildDateDeterministic(t *testing.T) {
	config := lifecycle.NewConfig()
	config.Extra = map[string]interface{}{"url": "https://example.com"}

	old := time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC)
	newer := time.Date(2024, 2, 2, 12, 0, 0, 0, time.UTC)

	feed := &lifecycle.Feed{
		Name:  "test",
		Title: "Test Feed",
		Path:  "test",
		Posts: []*models.Post{
			{Slug: "old", Href: "/old/", Date: &old},
			{Slug: "new", Href: "/new/", Date: &newer},
		},
	}

	rss, err := GenerateRSS(feed, config)
	if err != nil {
		t.Fatalf("GenerateRSS error: %v", err)
	}

	if !strings.Contains(rss, "<lastBuildDate>Fri, 02 Feb 2024 12:00:00 +0000</lastBuildDate>") {
		t.Fatalf("expected lastBuildDate to use latest post date, got:\n%s", rss)
	}
}

func TestGenerateAtom_UpdatedFallsBackWhenNoDates(t *testing.T) {
	config := lifecycle.NewConfig()
	config.Extra = map[string]interface{}{"url": "https://example.com"}

	feed := &lifecycle.Feed{
		Name:  "test",
		Title: "Test Feed",
		Path:  "test",
		Posts: []*models.Post{
			{Slug: "one", Href: "/one/"},
			{Slug: "two", Href: "/two/"},
		},
	}

	atom, err := GenerateAtom(feed, config)
	if err != nil {
		t.Fatalf("GenerateAtom error: %v", err)
	}

	if !strings.Contains(atom, "<updated>1970-01-01T00:00:00Z</updated>") {
		t.Fatalf("expected updated fallback when no dates, got:\n%s", atom)
	}
}

func TestGenerateRSS_LastBuildDateFallsBackWhenNoDates(t *testing.T) {
	config := lifecycle.NewConfig()
	config.Extra = map[string]interface{}{"url": "https://example.com"}

	feed := &lifecycle.Feed{
		Name:  "test",
		Title: "Test Feed",
		Path:  "test",
		Posts: []*models.Post{
			{Slug: "one", Href: "/one/"},
			{Slug: "two", Href: "/two/"},
		},
	}

	rss, err := GenerateRSS(feed, config)
	if err != nil {
		t.Fatalf("GenerateRSS error: %v", err)
	}

	if !strings.Contains(rss, "<lastBuildDate>Thu, 01 Jan 1970 00:00:00 +0000</lastBuildDate>") {
		t.Fatalf("expected lastBuildDate fallback when no dates, got:\n%s", rss)
	}
}

func TestGenerateRSS_IncludePrivateUsesEncryptedHTMLAndPrivateDates(t *testing.T) {
	config := lifecycle.NewConfig()
	config.Extra = map[string]interface{}{"url": "https://example.com"}

	older := time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC)
	privateDate := time.Date(2024, 2, 2, 12, 0, 0, 0, time.UTC)

	feed := &lifecycle.Feed{
		Name:           "test",
		Title:          "Test Feed",
		Path:           "test",
		IncludePrivate: true,
		Posts: []*models.Post{
			{Slug: "public", Href: "/public/", Date: &older, Published: true, Title: testStringPtr("Public")},
			{Slug: "private", Href: "/private/", Date: &privateDate, Published: true, Private: true, Title: testStringPtr("Private"), Description: testStringPtr("secret summary"), Content: "secret body", ArticleHTML: `<div class="encrypted-content" data-encrypted="ciphertext" data-key-name="private-key">locked</div>`},
		},
	}

	rss, err := GenerateRSS(feed, config)
	if err != nil {
		t.Fatalf("GenerateRSS error: %v", err)
	}

	checks := []string{
		"<lastBuildDate>Fri, 02 Feb 2024 12:00:00 +0000</lastBuildDate>",
		"&lt;div class=&#34;encrypted-content&#34; data-encrypted=&#34;ciphertext&#34;&gt;locked&lt;/div&gt;",
	}
	for _, want := range checks {
		if !strings.Contains(rss, want) {
			t.Fatalf("expected rss to contain %q\n%s", want, rss)
		}
	}
	if strings.Contains(rss, "secret summary") || strings.Contains(rss, "secret body") {
		t.Fatalf("rss should not expose private plaintext\n%s", rss)
	}
	if strings.Contains(rss, "data-key-name") || strings.Contains(rss, "private-key") {
		t.Fatalf("rss should not expose private key names\n%s", rss)
	}
}

func TestGenerateAtom_IncludePrivateUsesEncryptedHTMLAndPrivateDates(t *testing.T) {
	config := lifecycle.NewConfig()
	config.Extra = map[string]interface{}{"url": "https://example.com"}

	older := time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC)
	privateDate := time.Date(2024, 2, 2, 12, 0, 0, 0, time.UTC)

	feed := &lifecycle.Feed{
		Name:           "test",
		Title:          "Test Feed",
		Path:           "test",
		IncludePrivate: true,
		Posts: []*models.Post{
			{Slug: "public", Href: "/public/", Date: &older, Published: true, Title: testStringPtr("Public")},
			{Slug: "private", Href: "/private/", Date: &privateDate, Published: true, Private: true, Title: testStringPtr("Private"), Description: testStringPtr("secret summary"), Content: "secret body", ArticleHTML: `<div class="encrypted-content" data-encrypted="ciphertext" data-key-name="private-key">locked</div>`},
		},
	}

	atom, err := GenerateAtom(feed, config)
	if err != nil {
		t.Fatalf("GenerateAtom error: %v", err)
	}

	checks := []string{
		"<updated>2024-02-02T12:00:00Z</updated>",
		"&lt;div class=&#34;encrypted-content&#34; data-encrypted=&#34;ciphertext&#34;&gt;locked&lt;/div&gt;",
	}
	for _, want := range checks {
		if !strings.Contains(atom, want) {
			t.Fatalf("expected atom to contain %q\n%s", want, atom)
		}
	}
	if strings.Contains(atom, "<summary>") || strings.Contains(atom, "secret body") || strings.Contains(atom, "secret summary") {
		t.Fatalf("atom should not expose private plaintext\n%s", atom)
	}
	if strings.Contains(atom, "data-key-name") || strings.Contains(atom, "private-key") {
		t.Fatalf("atom should not expose private key names\n%s", atom)
	}
}
