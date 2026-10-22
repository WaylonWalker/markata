// Package urlmentions consumes the public, versioned URL-mentions protocol.
// It has no dependency on either provider's application code.
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

type Provider struct { Name string `json:"name"`; URL string `json:"url"` }
type Reference struct { URL string `json:"url"`; NormalizedURL string `json:"normalized_url"`; Type string `json:"type"` }
type Entry struct {
 ID string `json:"id"`
 CanonicalURL string `json:"canonical_url"`
 SourceURL string `json:"source_url"`
 Title string `json:"title"`
 PublishedAt string `json:"published_at"`
 Provider Provider `json:"provider"`
 References []Reference `json:"references"`
}
type Page struct {
 Protocol string `json:"protocol"`
 Version string `json:"version"`
 Provider Provider `json:"provider"`
 Total int `json:"total"`
 Limit int `json:"limit"`
 Offset int `json:"offset"`
 Next *string `json:"next"`
 Items []Entry `json:"items"`
}
type Client struct {
 HTTP *http.Client
 MaxPages int
 MaxBytes int64
}
func (c Client) httpClient() *http.Client {
 if c.HTTP != nil { return c.HTTP }
 return &http.Client{Timeout: 10*time.Second}
}
// Fetch follows bounded same-origin pagination and returns an inverted index
// keyed by the provider-supplied normalized_url.
func (c Client) Fetch(ctx context.Context, endpoint string) (map[string][]Entry, error) {
 first, err := url.Parse(endpoint)
 if err != nil || first.Scheme != "https" || first.Host == "" || first.User != nil {
  return nil, errors.New("urlmentions: endpoint must be an absolute HTTPS URL")
 }
 maxPages := c.MaxPages; if maxPages <= 0 { maxPages = 100 }
 maxBytes := c.MaxBytes; if maxBytes <= 0 { maxBytes = 4 << 20 }
 result := make(map[string][]Entry)
 seenPages := make(map[string]bool)
 seenEntries := make(map[string]map[string]bool)
 next := endpoint
 for pageNum := 0; next != ""; pageNum++ {
  if pageNum >= maxPages { return nil, fmt.Errorf("urlmentions: exceeded %d pages", maxPages) }
  target, err := url.Parse(next)
  if err != nil || target.Scheme != first.Scheme || !strings.EqualFold(target.Host, first.Host) || target.User != nil {
   return nil, errors.New("urlmentions: cross-origin or invalid next URL")
  }
  if seenPages[target.String()] { return nil, errors.New("urlmentions: pagination cycle") }
  seenPages[target.String()] = true
  req, err := http.NewRequestWithContext(ctx, http.MethodGet, target.String(), nil)
  if err != nil { return nil, err }
  req.Header.Set("Accept", "application/json")
  resp, err := c.httpClient().Do(req)
  if err != nil { return nil, err }
  if resp.StatusCode != http.StatusOK { resp.Body.Close(); return nil, fmt.Errorf("urlmentions: %s returned HTTP %d", target, resp.StatusCode) }
  var p Page
  dec := json.NewDecoder(io.LimitReader(resp.Body, maxBytes+1))
  err = dec.Decode(&p)
  resp.Body.Close()
  if err != nil { return nil, fmt.Errorf("urlmentions: decode page: %w", err) }
  if p.Version != "1.0" { return nil, fmt.Errorf("urlmentions: unsupported version %q", p.Version) }
  for _, entry := range p.Items {
   if entry.ID == "" || entry.SourceURL == "" { continue }
   for _, ref := range entry.References {
    key := ref.NormalizedURL
    if key == "" { continue }
    if seenEntries[key] == nil { seenEntries[key] = make(map[string]bool) }
    unique := entry.Provider.URL + "\x00" + entry.ID
    if seenEntries[key][unique] { continue }
    seenEntries[key][unique] = true
    result[key] = append(result[key], entry)
   }
  }
  if p.Next == nil { break }
  next = *p.Next
 }
 return result, nil
}
