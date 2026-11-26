package blogroll

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestParseFeedMetadata_RSSFeedTags(t *testing.T) {
	metadata := parseFeedMetadata(`
		<rss version="2.0">
		  <channel>
		    <title>Example Feed</title>
		    <description>Example Description</description>
		    <category>Go</category>
		    <category>Feeds</category>
		    <item>
		      <title>Post One</title>
		      <category>Ignored Entry Tag</category>
		    </item>
		  </channel>
		</rss>`)

	if metadata.FeedTitle != "Example Feed" {
		t.Fatalf("FeedTitle = %q, want %q", metadata.FeedTitle, "Example Feed")
	}
	if len(metadata.FeedTags) != 2 || metadata.FeedTags[0] != "Go" || metadata.FeedTags[1] != "Feeds" {
		t.Fatalf("FeedTags = %v, want [Go Feeds]", metadata.FeedTags)
	}
}

func TestParseFeedMetadata_AtomFeedTags(t *testing.T) {
	metadata := parseFeedMetadata(`
		<feed xmlns="http://www.w3.org/2005/Atom">
		  <title>Channel Feed</title>
		  <subtitle>Channel Description</subtitle>
		  <link rel="alternate" href="https://example.com/channel" />
		  <category term="Video" />
		  <category term="Tutorials" />
		  <entry>
		    <title>Ignored Entry</title>
		    <category term="Ignored Entry Tag" />
		  </entry>
		</feed>`)

	if metadata.FeedTitle != "Channel Feed" {
		t.Fatalf("FeedTitle = %q, want %q", metadata.FeedTitle, "Channel Feed")
	}
	if len(metadata.FeedTags) != 2 || metadata.FeedTags[0] != "Video" || metadata.FeedTags[1] != "Tutorials" {
		t.Fatalf("FeedTags = %v, want [Video Tutorials]", metadata.FeedTags)
	}
	if metadata.SiteURL != "https://example.com/channel" {
		t.Fatalf("SiteURL = %q, want %q", metadata.SiteURL, "https://example.com/channel")
	}
}

func TestParseHTMLMetadata_Tags(t *testing.T) {
	metadata := parseHTMLMetadata(`
		<html>
		  <head>
		    <title>Example Site</title>
		    <meta name="keywords" content="go, rss, feeds">
		    <meta property="article:tag" content="atom">
		    <meta property="article:tag" content="Go">
		  </head>
		</html>`, "https://example.com")

	if len(metadata.Tags) != 4 {
		t.Fatalf("Tags length = %d, want 4 (%v)", len(metadata.Tags), metadata.Tags)
	}
	if metadata.Tags[0] != "go" || metadata.Tags[1] != "rss" || metadata.Tags[2] != "feeds" || metadata.Tags[3] != "atom" {
		t.Fatalf("Tags = %v, want [go rss feeds atom]", metadata.Tags)
	}
}

func TestParseHTMLMetadata_DecodesEntitiesAndQuotedTags(t *testing.T) {
	metadata := parseHTMLMetadata(`
		<html>
		  <head>
		    <title>devtools-fm - YouTube</title>
		    <meta name="description" content="Hi, we&#39;re Andrew and Justin">
		    <meta name="keywords" content="podcast technology coding programming developer &quot;developer tools&quot;">
		  </head>
		</html>`, "https://example.com")

	if metadata.Description != "Hi, we're Andrew and Justin" {
		t.Fatalf("Description = %q, want %q", metadata.Description, "Hi, we're Andrew and Justin")
	}
	wantTags := []string{"podcast", "technology", "coding", "programming", "developer", "developer tools"}
	if len(metadata.Tags) != len(wantTags) {
		t.Fatalf("Tags length = %d, want %d (%v)", len(metadata.Tags), len(wantTags), metadata.Tags)
	}
	for i, want := range wantTags {
		if metadata.Tags[i] != want {
			t.Fatalf("Tags[%d] = %q, want %q (all=%v)", i, metadata.Tags[i], want, metadata.Tags)
		}
	}
}

func TestParseHTMLMetadata_StripsYouTubeTitleSuffix(t *testing.T) {
	metadata := parseHTMLMetadata(`
		<html>
		  <head>
		    <title>devtools-fm - YouTube</title>
		  </head>
		</html>`, "https://www.youtube.com/channel/UCFsRlOn7gODgv6WUriLrzXg")

	if metadata.Title != "devtools-fm" {
		t.Fatalf("Title = %q, want %q", metadata.Title, "devtools-fm")
	}
}

func TestFetchMetadata_UsesFeedSiteURLForSiteMetadata(t *testing.T) {
	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/feed":
			w.Header().Set("Content-Type", "application/atom+xml")
			if _, err := w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
				<feed xmlns="http://www.w3.org/2005/Atom">
				  <title>DevTools FM</title>
				  <link rel="alternate" href="` + server.URL + `/channel" />
				  <author>
				    <name>DevTools FM</name>
				    <uri>` + server.URL + `/channel</uri>
				  </author>
				  <entry>
				    <title>Episode</title>
				  </entry>
				</feed>`)); err != nil {
				t.Fatalf("write feed response: %v", err)
			}
		case "/channel":
			w.Header().Set("Content-Type", "text/html")
			if _, err := w.Write([]byte(`<html><head><meta name="description" content="Channel-specific description"></head><body></body></html>`)); err != nil {
				t.Fatalf("write channel response: %v", err)
			}
		case "/":
			w.Header().Set("Content-Type", "text/html")
			if _, err := w.Write([]byte(`<html><head><meta name="description" content="Generic site description"></head><body></body></html>`)); err != nil {
				t.Fatalf("write root response: %v", err)
			}
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	updater := NewUpdater(5 * time.Second)
	metadata, err := updater.FetchMetadata(context.Background(), server.URL+"/feed")
	if err != nil {
		t.Fatalf("FetchMetadata() error = %v", err)
	}
	if metadata.Description != "Channel-specific description" {
		t.Fatalf("Description = %q, want %q", metadata.Description, "Channel-specific description")
	}
	if metadata.SiteURL != server.URL+"/channel" {
		t.Fatalf("SiteURL = %q, want %q", metadata.SiteURL, server.URL+"/channel")
	}
	if metadata.FeedTitle != "DevTools FM" {
		t.Fatalf("FeedTitle = %q, want %q", metadata.FeedTitle, "DevTools FM")
	}
}

func TestFetchMetadataPreservesFeedAndSiteErrors(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/feed" {
			http.Error(w, "feed unavailable", http.StatusForbidden)
			return
		}
		http.Error(w, "site unavailable", http.StatusServiceUnavailable)
	}))
	defer server.Close()

	_, err := NewUpdater(time.Second).FetchMetadata(context.Background(), server.URL+"/feed")
	if err == nil || !strings.Contains(err.Error(), "feed: HTTP 403") || !strings.Contains(err.Error(), "website: HTTP 503") {
		t.Fatalf("expected both fetch failures, got %v", err)
	}
}
