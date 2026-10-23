package plugins

import "testing"

func TestReaderMentionsForURL(t *testing.T) {
 index := map[string][]map[string]string{
  "https://example.com/post": {{"title":"Mention"}},
 }
 cases := []struct{ url string; want int }{
  {"https://example.com/post",1},
  {"https://EXAMPLE.COM/post/",1},
  {"https://example.com/post#section",1},
  {"https://example.com/other",0},
  {"https://example.com/post?tracking=1",0},
  {"javascript:alert(1)",0},
 }
 for _, tc := range cases {
  got := readerMentionsForURL(index, tc.url)
  if len(got) != tc.want { t.Errorf("%q: got %d, want %d",tc.url,len(got),tc.want) }
 }
}

func TestReaderURLMentionsDisabledByDefault(t *testing.T) {
 if got := readerURLMentions(nil); len(got) != 0 {
  t.Fatalf("unconfigured providers returned %d mention targets",len(got))
 }
}
