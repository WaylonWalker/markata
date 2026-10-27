package plugins

import (
 "testing"

 "github.com/WaylonWalker/markata/pkg/lifecycle"
 "github.com/WaylonWalker/markata/pkg/models"
)

func TestMentionsReferencedHandlesDemandDriven(t *testing.T) {
 p := NewMentionsPlugin()
 alice := &mentionEntry{Handle:"alice",SiteURL:"https://alice.example"}
 bob := &mentionEntry{Handle:"bob",SiteURL:"https://bob.example"}
 handles := map[string]*mentionEntry{
  "alice":alice, "alias":alice, "bob":bob,
 }
 tests := []struct{
  name string
  posts []*models.Post
  want []string
 }{
  {"no mentions", []*models.Post{{Content:"Just a normal article",Published:true}}, nil},
  {"unknown mention", []*models.Post{{Content:"Hello @missing",Published:true}}, nil},
  {"known with unused contact", []*models.Post{{Content:"Hello @alice",Published:true}}, []string{"alice"}},
  {"alias and punctuation", []*models.Post{{Content:"Thanks @ALIAS.",Published:true}}, []string{"alias"}},
  {"chat title", []*models.Post{{Content:"!!! chat \"@alice\"\n    Hello",Published:true}}, []string{"alice"}},
  {"attribution", []*models.Post{{Content:"Quoted from ~bob",Published:true}}, []string{"bob"}},
  {"skip other post", []*models.Post{{Content:"@alice",Skip:true},{Content:"@bob",Published:true}}, []string{"bob"}},
 }
 for _,tc:=range tests {
  t.Run(tc.name,func(t *testing.T){
   m:=lifecycle.NewManager()
   m.SetPosts(tc.posts)
   got:=p.referencedHandles(m,handles)
   if len(got)!=len(tc.want){t.Fatalf("got %d handles, want %d: %#v",len(got),len(tc.want),got)}
   for _,handle:=range tc.want {
    if got[handle]!=handles[handle]{t.Errorf("missing referenced handle %s",handle)}
   }
  })
 }
}
