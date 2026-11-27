package cmd

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

func TestVerifyDirectBlogrollFeed(t *testing.T) {
	const address = "http://localhost:8080/feed"
	for _, tt := range []struct {
		name string
		body string
		ok   bool
	}{
		{"RSS", "<rss><channel><title>Example</title></channel></rss>", true},
		{"HTML", "<html><title>Example</title></html>", false},
		{"JSON Feed", `{"version":"https://jsonfeed.org/version/1.1","items":[]}`, false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			client := &http.Client{Transport: subscriptionFixtureTransport{address: tt.body}}
			err := verifyDirectBlogrollFeed(context.Background(), client, address)
			if (err == nil) != tt.ok {
				t.Fatalf("unexpected validation result: %v", err)
			}
		})
	}
}

func TestChooseBlogrollCandidateNoninteractive(t *testing.T) {
	const first = "https://example.test/rss"
	const second = "https://example.test/atom"
	got, err := chooseBlogrollCandidate([]string{first}, false)
	if err != nil || got != first {
		t.Fatalf("single candidate: %q, %v", got, err)
	}
	if _, err := chooseBlogrollCandidate([]string{first, second}, false); err == nil || !strings.Contains(FormatError(err), second) {
		t.Fatalf("expected candidate choices, got %v", err)
	}
}

func TestDiscoverSubscriptionMultipleCandidates(t *testing.T) {
	const site = "https://example.test"
	client := &http.Client{Transport: subscriptionFixtureTransport{
		site: `<html><head>
<link rel="alternate" type="application/rss+xml" href="/rss">
<link rel="alternate" type="application/rss+xml" href="/rss">
<link rel="alternate" type="application/atom+xml" href="/atom"></head></html>`,
		site + "/rss":  "<rss><channel><title>Example</title></channel></rss>",
		site + "/atom": `<feed xmlns="http://www.w3.org/2005/Atom"><title>Example</title></feed>`,
	}}
	got, err := discoverSubscriptionWithClient(context.Background(), client, site)
	if err != nil || len(got.Candidates) != 2 || got.URL != site+"/rss" || got.Candidates[1] != site+"/atom" {
		t.Fatalf("unexpected candidates: %+v, %v", got, err)
	}
}
