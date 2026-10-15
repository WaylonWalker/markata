package plugins

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLemmyCommunityRSSPreservesArticleAndDiscussion(t *testing.T) {
	article := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		if _, err := w.Write([]byte(`<html><head><meta property="og:title" content="Fediverse article"></head></html>`)); err != nil {
			t.Errorf("write html: %v", err)
		}
	}))
	defer article.Close()

	discussion := "https://lemmy.world/post/1234567"
	rss := []byte(fmt.Sprintf(`<rss><channel><title>Selfhosted - Lemmy</title>
		<item><title>Original headline</title><link>%s</link>
		<guid>%s</guid><comments>%s</comments></item>
	</channel></rss>`, article.URL, discussion, discussion))
	_, entries, err := parseRSS2Feed(rss)
	if err != nil {
		t.Fatalf("parseRSS2Feed: %v", err)
	}
	if len(entries) != 1 || entries[0].URL != article.URL ||
		entries[0].OriginalURL != discussion || entries[0].Title != "Fediverse article" {
		t.Fatalf("Lemmy article/discussion missing: %+v", entries)
	}
	if got := aggregatorDiscussionLabel(entries[0].OriginalURL); got != "Lemmy" {
		t.Fatalf("discussion label = %q", got)
	}
}

func TestLemmyInstanceAwarePermalinks(t *testing.T) {
	for _, url := range []string{
		"https://lemmy.world/post/123",
		"https://programming.dev/post/987654",
		"https://feddit.example/post/42",
	} {
		if !isLemmyDiscussionURL(url) {
			t.Errorf("did not recognize Lemmy post %q", url)
		}
	}
	for _, url := range []string{
		"https://lemmy.world/c/selfhosted",
		"https://lemmy.world/post/not-an-id",
		"https://lemmy.world.evil.com/post/a2",
	} {
		if isLemmyDiscussionURL(url) {
			t.Errorf("incorrect Lemmy permalink match: %q", url)
		}
	}
}

func TestRedditAtomExternalLinkAndSelfPost(t *testing.T) {
	article := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		if _, err := w.Write([]byte(`<html><head><meta property="og:description" content="Linked article summary"></head></html>`)); err != nil {
			t.Errorf("write html: %v", err)
		}
	}))
	defer article.Close()

	discussion := "https://www.reddit.com/r/selfhosted/comments/abcdef/an_article/"
	// Reddit Atom includes an explicitly labelled [link] for external articles.
	content := fmt.Sprintf(`&lt;div&gt;&lt;a href="%s"&gt;[link]&lt;/a&gt; &lt;a href="%s"&gt;[comments]&lt;/a&gt;&lt;/div&gt;`, article.URL, discussion)
	feed := []byte(fmt.Sprintf(`<feed xmlns="http://www.w3.org/2005/Atom">
		<title>r/selfhosted</title>
		<entry><id>t3_abcdef</id><title>Article</title>
		<link rel="alternate" href="%s"/><content type="html">%s</content></entry>
		<entry><id>t3_ghijkl</id><title>Self post</title>
		<link rel="alternate" href="https://www.reddit.com/r/selfhosted/comments/ghijkl/self_post/"/>
		<content type="html">&lt;p&gt;Just discussion&lt;/p&gt;</content></entry>
	</feed>`, discussion, content))
	_, entries, err := parseAtomFeed(feed)
	if err != nil {
		t.Fatalf("parseAtomFeed: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected two entries, got %d", len(entries))
	}
	if entries[0].URL != article.URL || entries[0].OriginalURL != discussion ||
		entries[0].Description != "Linked article summary" {
		t.Fatalf("Reddit outbound post not resolved: %+v", entries[0])
	}
	if got := aggregatorDiscussionLabel(entries[0].OriginalURL); got != "Reddit" {
		t.Fatalf("discussion label = %q", got)
	}
	if entries[1].URL != "https://www.reddit.com/r/selfhosted/comments/ghijkl/self_post/" ||
		entries[1].OriginalURL != "" {
		t.Fatalf("Reddit self post should remain on Reddit: %+v", entries[1])
	}
}

func TestRedditOutboundExtractionRejectsUnrelatedLinks(t *testing.T) {
	discussion := "https://reddit.com/r/golang/comments/ab123/story/"
	for _, content := range []string{
		`<a href="https://bad.example">read me</a>`,
		`<a href="javascript:alert(1)">[link]</a>`,
		`<a href="https://www.reddit.com/r/test">[link]</a>`,
		`<a href="https://reddit.com.evil.example">read me</a>`,
	} {
		if got := redditOutboundArticleURL(discussion, content); got != "" {
			t.Errorf("promoted unrelated link %q to %q", content, got)
		}
	}
	allowed := `<a href="https://reddit.com.evil.example/article">[link]</a>`
	if got := redditOutboundArticleURL(discussion, allowed); got != "https://reddit.com.evil.example/article" {
		t.Fatalf("valid external domain was rejected: %q", got)
	}
	if isRedditDiscussionURL("https://reddit.com.evil.example/r/golang/comments/abc/test") {
		t.Fatal("subdomain spoofing must not match Reddit")
	}
	if !strings.Contains(aggregatorDiscussionLabel(discussion), "Reddit") {
		t.Fatal("Reddit discussion label missing")
	}
}
