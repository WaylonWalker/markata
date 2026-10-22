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
type Client struct { HTTP *http.Client; MaxPages int; MaxBytes int64 }
func sameOrigin(a,b *url.URL) bool {
 return a.Scheme=="https" && b.Scheme=="https" && strings.EqualFold(a.Host,b.Host) && b.User==nil && b.Fragment==""
}
func (c Client) Fetch(ctx context.Context, endpoint string) (map[string][]Entry,error) {
 origin,err:=url.Parse(endpoint)
 if err!=nil || origin.Scheme!="https" || origin.Host=="" || origin.User!=nil || origin.Fragment!="" { return nil,errors.New("urlmentions: endpoint must be an absolute HTTPS URL") }
 maxPages:=c.MaxPages; if maxPages<=0 {maxPages=100}
 maxBytes:=c.MaxBytes; if maxBytes<=0 {maxBytes=4<<20}
 client:=&http.Client{Timeout:10*time.Second}
 if c.HTTP!=nil { clone:=*c.HTTP; client=&clone; if client.Timeout==0 {client.Timeout=10*time.Second} }
 previousRedirect:=client.CheckRedirect
 client.CheckRedirect=func(req *http.Request, via []*http.Request) error {
  if !sameOrigin(origin,req.URL) {return errors.New("urlmentions: cross-origin redirect")}
  if len(via)>=10 {return errors.New("urlmentions: too many redirects")}
  if previousRedirect!=nil {return previousRedirect(req,via)}
  return nil
 }
 result:=make(map[string][]Entry)
 seenPages:=make(map[string]bool)
 seenEntries:=make(map[string]map[string]bool)
 next:=endpoint
 var total int64
 for pageNum:=0;next!="";pageNum++ {
  if pageNum>=maxPages {return nil,fmt.Errorf("urlmentions: exceeded %d pages",maxPages)}
  target,err:=url.Parse(next)
  if err!=nil || !sameOrigin(origin,target) {return nil,errors.New("urlmentions: cross-origin or invalid next URL")}
  if seenPages[target.String()] {return nil,errors.New("urlmentions: pagination cycle")}
  seenPages[target.String()]=true
  req,err:=http.NewRequestWithContext(ctx,http.MethodGet,target.String(),nil);if err!=nil{return nil,err}
  req.Header.Set("Accept","application/json")
  resp,err:=client.Do(req);if err!=nil{return nil,err}
  if resp.StatusCode!=http.StatusOK {resp.Body.Close();return nil,fmt.Errorf("urlmentions: HTTP %d",resp.StatusCode)}
  remaining:=maxBytes-total
  if remaining<=0 {resp.Body.Close();return nil,errors.New("urlmentions: response budget exceeded")}
  raw,err:=io.ReadAll(io.LimitReader(resp.Body,remaining+1));resp.Body.Close()
  if err!=nil{return nil,err}
  if int64(len(raw))>remaining {return nil,errors.New("urlmentions: response budget exceeded")}
  total+=int64(len(raw))
  var p Page
  if err=json.Unmarshal(raw,&p);err!=nil{return nil,fmt.Errorf("urlmentions: decode page: %w",err)}
  if p.Version!="1.0" || p.Protocol!="url-mentions" {return nil,fmt.Errorf("urlmentions: unsupported protocol %q version %q",p.Protocol,p.Version)}
  if p.Offset<0 || p.Total<0 || p.Limit<0 {return nil,errors.New("urlmentions: invalid pagination metadata")}
  for _,entry:=range p.Items {
   if entry.ID=="" || entry.SourceURL=="" {return nil,errors.New("urlmentions: entry missing id or source_url")}
   if u,e:=url.Parse(entry.SourceURL);e!=nil || u.Scheme!="https" || u.Host=="" {return nil,errors.New("urlmentions: invalid source_url")}
   for _,ref:=range entry.References {
    if ref.NormalizedURL=="" {return nil,errors.New("urlmentions: reference missing normalized_url")}
    key:=ref.NormalizedURL
    if seenEntries[key]==nil {seenEntries[key]=make(map[string]bool)}
    unique:=entry.Provider.URL+"\x00"+entry.SourceURL+"\x00"+entry.ID
    if seenEntries[key][unique] {continue}
    seenEntries[key][unique]=true
    result[key]=append(result[key],entry)
   }
  }
  if p.Next==nil {break}
  if *p.Next=="" {return nil,errors.New("urlmentions: empty next URL")}
  next=*p.Next
 }
 return result,nil
}
