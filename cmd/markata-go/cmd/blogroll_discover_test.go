package cmd

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

type subscriptionFixtureTransport map[string]string

func (f subscriptionFixtureTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	body, ok := f[req.URL.String()]
	status := http.StatusOK
	if !ok {
		status = http.StatusNotFound
	}
	return &http.Response{
		StatusCode: status,
		Status:     http.StatusText(status),
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     make(http.Header),
		Request:    req,
	}, nil
}

func TestDiscoverSubscriptionValidatesAlternatives(t *testing.T) {
	const site = "https://example.test"
	client := &http.Client{Transport: subscriptionFixtureTransport{
		site: `<html><head><title>Example</title>
<link rel="alternate" type="application/rss+xml" href="/bad">
<link rel="alternate" type="application/atom+xml" href="/valid"></head></html>`,
		site + "/bad":   "<html>not a feed</html>",
		site + "/valid": `<feed xmlns="http://www.w3.org/2005/Atom"><title>Example</title></feed>`,
	}}
	got, err := discoverSubscriptionWithClient(context.Background(), client, site)
	if err != nil || got.URL != site+"/valid" {
		t.Fatalf("expected validated alternate, got %+v, %v", got, err)
	}
}

func TestDiscoverSubscriptionRejectsMasqueradingFeed(t *testing.T) {
	const site = "https://example.test"
	client := &http.Client{Transport: subscriptionFixtureTransport{
		site: `<?xml version="1.0"?><html><title>Access denied</title></html>`,
	}}
	if _, err := discoverSubscriptionWithClient(context.Background(), client, site); err == nil {
		t.Fatal("unexpected success for HTML with XML declaration")
	}
}
