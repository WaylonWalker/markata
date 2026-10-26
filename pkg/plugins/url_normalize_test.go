package plugins

import (
	"testing"

	"github.com/WaylonWalker/markata/pkg/models"
)

func TestNormalizeExternalURL_UpgradesPublicHTTP(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "public domain", in: "http://example.com/image.png", want: "https://example.com/image.png"},
		{name: "protocol relative", in: "//example.com/image.png", want: "https://example.com/image.png"},
		{name: "already https", in: "https://example.com/image.png", want: "https://example.com/image.png"},
		{name: "localhost", in: "http://localhost:3000/image.png", want: "http://localhost:3000/image.png"},
		{name: "private ip", in: "http://192.168.1.10/image.png", want: "http://192.168.1.10/image.png"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normalizeExternalURL(tt.in); got != tt.want {
				t.Fatalf("normalizeExternalURL(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestNormalizeExternalFeed_NormalizesNestedURLs(t *testing.T) {
	feed := &models.ExternalFeed{
		FeedURL:   "http://example.com/feed.xml",
		SiteURL:   "http://example.com",
		ImageURL:  "http://example.com/logo.png",
		AvatarURL: "http://example.com/avatar.png",
		Entries: []*models.ExternalEntry{
			{FeedURL: "http://example.com/feed.xml", URL: "http://example.com/post", ImageURL: "http://example.com/post.png"},
		},
	}

	normalizeExternalFeed(feed)

	if feed.FeedURL != "https://example.com/feed.xml" {
		t.Fatalf("feed.FeedURL = %q", feed.FeedURL)
	}
	if feed.SiteURL != "https://example.com" {
		t.Fatalf("feed.SiteURL = %q", feed.SiteURL)
	}
	if feed.ImageURL != "https://example.com/logo.png" {
		t.Fatalf("feed.ImageURL = %q", feed.ImageURL)
	}
	if feed.AvatarURL != "https://example.com/avatar.png" {
		t.Fatalf("feed.AvatarURL = %q", feed.AvatarURL)
	}
	if feed.Entries[0].URL != "https://example.com/post" {
		t.Fatalf("entry.URL = %q", feed.Entries[0].URL)
	}
}
