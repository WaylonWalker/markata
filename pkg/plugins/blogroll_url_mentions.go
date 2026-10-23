package plugins

import (
 "context"
 "log"
 "net/url"
 "sort"
 "strings"
 "time"

 "github.com/WaylonWalker/markata/pkg/urlmentions"
)

// readerURLMentions loads opt-in public provider indexes once per build.
// Network errors are non-fatal: a provider outage must not break the site build.
func readerURLMentions(raw interface{}) map[string][]map[string]string {
 endpoints := make([]string, 0)
 switch v := raw.(type) {
 case []string:
  endpoints = append(endpoints, v...)
 case []interface{}:
  for _, item := range v {
   if endpoint, ok := item.(string); ok { endpoints = append(endpoints, endpoint) }
  }
 }
 result := make(map[string][]map[string]string)
 client := urlmentions.Client{MaxPages: 8, MaxBytes: 2 << 20}
 ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
 defer cancel()
 for _, endpoint := range endpoints {
  parsed, err := url.Parse(endpoint)
  if err != nil || parsed.Scheme != "https" || parsed.Hostname() == "" || parsed.User != nil { continue }
  index, err := client.Fetch(ctx, endpoint)
  if err != nil { log.Printf("[url_mentions] provider %s unavailable: %v", parsed.Hostname(), err); continue }
  for target, entries := range index {
   for _, entry := range entries {
    source, err := url.Parse(entry.SourceURL)
    if err != nil || source.Scheme != "https" || source.Hostname() == "" { continue }
    provider := strings.TrimSpace(entry.Provider.Name)
    if provider == "" { provider = parsed.Hostname() }
    title := strings.TrimSpace(entry.Title)
    if title == "" { title = entry.SourceURL }
    result[target] = append(result[target], map[string]string{
     "provider_name": provider,
     "icon_url": "https://www.google.com/s2/favicons?domain="+url.QueryEscape(parsed.Hostname())+"&sz=32",
     "source_url": entry.SourceURL,
     "title": title,
    })
   }
  }
 }
 for target, mentions := range result {
  sort.SliceStable(mentions, func(i,j int) bool {
   if mentions[i]["provider_name"] == mentions[j]["provider_name"] { return mentions[i]["source_url"] < mentions[j]["source_url"] }
   return mentions[i]["provider_name"] < mentions[j]["provider_name"]
  })
  dedup := mentions[:0]
  seen := map[string]bool{}
  for _, mention := range mentions {
   key := mention["source_url"]
   if seen[key] { continue }
   seen[key] = true
   dedup = append(dedup, mention)
  }
  result[target] = dedup
 }
 return result
}

// Normalize conservatively: don't conflate paths, query strings or hosts.
func readerMentionsForURL(index map[string][]map[string]string, raw string) []map[string]string {
 parsed, err := url.Parse(strings.TrimSpace(raw))
 if err != nil || (parsed.Scheme != "https" && parsed.Scheme != "http") || parsed.Hostname() == "" { return nil }
 parsed.Fragment = ""
 parsed.Host = strings.ToLower(parsed.Host)
 key := parsed.String()
 if matches := index[key]; len(matches) != 0 { return matches }
 if strings.HasSuffix(parsed.Path, "/") && parsed.Path != "/" { parsed.Path = strings.TrimSuffix(parsed.Path, "/") } else if parsed.Path != "/" { parsed.Path += "/" }
 return index[parsed.String()]
}
