package urlmentions

import (
 "context"
 "fmt"
 "net/http"
 "net/http/httptest"
 "strings"
 "testing"
)

func TestFetchPaginationAndDedup(t *testing.T) {
 var server *httptest.Server
 server = httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  next := "null"
  if r.URL.Query().Get("offset") == "" { next = fmt.Sprintf("%q", server.URL+"/.well-known/url-mentions.json?offset=1") }
  fmt.Fprintf(w, `{"protocol":"url-mentions","version":"1.0","total":2,"next":%s,"items":[{"id":"a","source_url":"https://source.example/a","provider":{"name":"Listen","url":"https://listen.example"},"references":[{"url":"https://example.com/post?utm_source=x","normalized_url":"https://example.com/post","type":"link"},{"url":"https://example.com/post","normalized_url":"https://example.com/post","type":"text"}]}]}`, next)
 }))
 defer server.Close()
 c := Client{HTTP:server.Client()}
 // TLS test server is intentionally used to exercise the HTTPS-only boundary.
 got, err := c.Fetch(context.Background(),server.URL+"/.well-known/url-mentions.json")
 if err != nil { t.Fatal(err) }
 if len(got["https://example.com/post"]) != 1 { t.Fatalf("expected deduplicated entry, got %#v",got) }
}
func TestFetchRejectsInsecureURL(t *testing.T) {
 _,err := (Client{}).Fetch(context.Background(),"http://example.com")
 if err == nil || !strings.Contains(err.Error(),"HTTPS") { t.Fatalf("expected HTTPS error, got %v",err) }
}
func TestFetchRejectsCrossOriginNext(t *testing.T) {
 server:=httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request) {
  fmt.Fprint(w,`{"version":"1.0","next":"https://evil.example/page","items":[]}`)
 }))
 defer server.Close()
 _,err:=(Client{HTTP:server.Client()}).Fetch(context.Background(),server.URL)
 if err==nil || !strings.Contains(err.Error(),"cross-origin") { t.Fatalf("expected cross-origin error, got %v",err) }
}
