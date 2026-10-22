// Package urlmentions consumes versioned public URL-reference indexes.
package urlmentions

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const httpsScheme = "https"

type Provider struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
type Reference struct {
	URL           string `json:"url"`
	NormalizedURL string `json:"normalized_url"`
	Type          string `json:"type"`
}
type Entry struct {
	ID           string      `json:"id"`
	CanonicalURL string      `json:"canonical_url"`
	SourceURL    string      `json:"source_url"`
	Title        string      `json:"title"`
	PublishedAt  string      `json:"published_at"`
	Provider     Provider    `json:"provider"`
	References   []Reference `json:"references"`
}
type Page struct {
	Protocol string   `json:"protocol"`
	Version  string   `json:"version"`
	Provider Provider `json:"provider"`
	Total    int      `json:"total"`
	Limit    int      `json:"limit"`
	Offset   int      `json:"offset"`
	Next     *string  `json:"next"`
	Items    []Entry  `json:"items"`
}
type Client struct {
	HTTP     *http.Client
	MaxPages int
	MaxBytes int64
}

func sameOrigin(a, b *url.URL) bool {
	return a.Scheme == httpsScheme && b.Scheme == httpsScheme && strings.EqualFold(a.Host, b.Host) && b.User == nil && b.Fragment == ""
}

func (c Client) httpClient(origin *url.URL) *http.Client {
	client := &http.Client{Timeout: 10 * time.Second}
	if c.HTTP != nil {
		clone := *c.HTTP
		client = &clone
		if client.Timeout == 0 {
			client.Timeout = 10 * time.Second
		}
	}
	previousRedirect := client.CheckRedirect
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		if !sameOrigin(origin, req.URL) {
			return errors.New("urlmentions: cross-origin redirect")
		}
		if len(via) >= 10 {
			return errors.New("urlmentions: too many redirects")
		}
		if previousRedirect != nil {
			return previousRedirect(req, via)
		}
		return nil
	}
	return client
}

func fetchPage(ctx context.Context, client *http.Client, target string, remaining int64) (Page, int64, error) {
	var page Page
	if remaining <= 0 {
		return page, 0, errors.New("urlmentions: response budget exceeded")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target, http.NoBody)
	if err != nil {
		return page, 0, err
	}
	req.Header.Set("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return page, 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return page, 0, fmt.Errorf("urlmentions: HTTP %d", resp.StatusCode)
	}
	raw, err := io.ReadAll(io.LimitReader(resp.Body, remaining+1))
	if err != nil {
		return page, 0, err
	}
	if int64(len(raw)) > remaining {
		return page, 0, errors.New("urlmentions: response budget exceeded")
	}
	if err = json.Unmarshal(raw, &page); err != nil {
		return page, 0, fmt.Errorf("urlmentions: decode page: %w", err)
	}
	if page.Version != "1.0" || page.Protocol != "url-mentions" {
		return page, 0, fmt.Errorf("urlmentions: unsupported protocol %q version %q", page.Protocol, page.Version)
	}
	if page.Offset < 0 || page.Total < 0 || page.Limit < 0 {
		return page, 0, errors.New("urlmentions: invalid pagination metadata")
	}
	return page, int64(len(raw)), nil
}

func indexEntries(page Page, result map[string][]Entry, seen map[string]map[string]bool) error {
	for i := range page.Items {
		entry := &page.Items[i]
		if entry.ID == "" || entry.SourceURL == "" {
			return errors.New("urlmentions: entry missing id or source_url")
		}
		source, err := url.Parse(entry.SourceURL)
		if err != nil || source.Scheme != httpsScheme || source.Host == "" {
			return errors.New("urlmentions: invalid source_url")
		}
		for _, ref := range entry.References {
			if ref.NormalizedURL == "" {
				return errors.New("urlmentions: reference missing normalized_url")
			}
			key := ref.NormalizedURL
			if seen[key] == nil {
				seen[key] = make(map[string]bool)
			}
			unique := entry.Provider.URL + "\x00" + entry.SourceURL + "\x00" + entry.ID
			if seen[key][unique] {
				continue
			}
			seen[key][unique] = true
			result[key] = append(result[key], *entry)
		}
	}
	return nil
}

// Fetch reads a bounded number of same-origin pages and indexes references by normalized URL.
func (c Client) Fetch(ctx context.Context, endpoint string) (map[string][]Entry, error) {
	origin, err := url.Parse(endpoint)
	if err != nil || origin.Scheme != httpsScheme || origin.Host == "" || origin.User != nil || origin.Fragment != "" {
		return nil, errors.New("urlmentions: endpoint must be an absolute HTTPS URL")
	}
	maxPages := c.MaxPages
	if maxPages <= 0 {
		maxPages = 100
	}
	maxBytes := c.MaxBytes
	if maxBytes <= 0 {
		maxBytes = 4 << 20
	}
	client := c.httpClient(origin)
	result := make(map[string][]Entry)
	seenPages := make(map[string]bool)
	seenEntries := make(map[string]map[string]bool)
	next := endpoint
	var total int64
	for pageNum := 0; next != ""; pageNum++ {
		if pageNum >= maxPages {
			return nil, fmt.Errorf("urlmentions: exceeded %d pages", maxPages)
		}
		target, err := url.Parse(next)
		if err != nil || !sameOrigin(origin, target) {
			return nil, errors.New("urlmentions: cross-origin or invalid next URL")
		}
		if seenPages[target.String()] {
			return nil, errors.New("urlmentions: pagination cycle")
		}
		seenPages[target.String()] = true
		page, n, err := fetchPage(ctx, client, target.String(), maxBytes-total)
		if err != nil {
			return nil, err
		}
		total += n
		if err = indexEntries(page, result, seenEntries); err != nil {
			return nil, err
		}
		if page.Next == nil {
			break
		}
		if *page.Next == "" {
			return nil, errors.New("urlmentions: empty next URL")
		}
		next = *page.Next
	}
	return result, nil
}
