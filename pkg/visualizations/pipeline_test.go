package visualizations

import (
 "context"
 "errors"
 "strings"
 "sync"
 "sync/atomic"
 "testing"
 "time"
)

const sampleSVG = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 960 540"><rect width="960" height="540" fill="#fefefe"/><text x="10" y="20">Markata</text></svg>`

type fakeRenderer struct {
 name, version string
 calls atomic.Int32
 fail bool
 pause chan struct{}
}
func (f *fakeRenderer) ID() string { return f.name }
func (f *fakeRenderer) Version() string { return f.version }
func (f *fakeRenderer) RenderSVG(ctx context.Context, req Request) ([]byte,error) {
 f.calls.Add(1)
 if f.pause != nil {
  select {
  case <-f.pause:
  case <-ctx.Done(): return nil, ctx.Err()
  }
 }
 if f.fail { return nil, errors.New("render failed") }
 return []byte(sampleSVG),nil
}

func TestPipelineSingleflightAndCopies(t *testing.T) {
 var p Pipeline
 engine := &fakeRenderer{name:"d2",version:"1",pause:make(chan struct{})}
 req := Request{Source:"a -> b",Palette:"ayu-dark",Font:"system",ThemeHash:"theme-sha",Options:map[string]string{"layout":"dagre"}}
 const workers = 25
 var wg sync.WaitGroup
 start := make(chan struct{})
 results := make([]Asset,workers)
 errs := make([]error,workers)
 for i:=0;i<workers;i++ {
  wg.Add(1)
  go func(i int) { defer wg.Done(); <-start; results[i],errs[i]=p.Compile(context.Background(),engine,req) }(i)
 }
 close(start)
 deadline:=time.After(time.Second)
 for engine.calls.Load()==0 {
  select {
  case <-deadline: t.Fatal("render did not start")
  default: time.Sleep(time.Millisecond)
  }
 }
 close(engine.pause)
 wg.Wait()
 for i:=range results {
  if errs[i]!=nil || results[i].Width!=960 || results[i].Height!=540 {
   t.Fatalf("result %d: %+v %v",i,results[i],errs[i])
  }
 }
 if engine.calls.Load()!=1 { t.Fatalf("renders=%d, want one",engine.calls.Load()) }
 results[0].SVG[0]='X'
 cached,err:=p.Compile(context.Background(),engine,req)
 if err!=nil || cached.SVG[0]!='<' {t.Fatalf("cached asset mutation: %v",err)}
}

func TestPipelineKeyAndFailures(t *testing.T) {
 var p Pipeline
 engine:=&fakeRenderer{name:"d2",version:"1"}
 a:=Request{Source:"a",Palette:"light",Options:map[string]string{"z":"2","a":"1"}}
 b:=Request{Source:"a",Palette:"light",Options:map[string]string{"a":"1","z":"2"}}
 if renderKey(engine,a)!=renderKey(engine,b) { t.Fatal("option order affected key") }
 for _,req:=range []Request{a,b,{Source:"a",Palette:"dark"},{Source:"b",Palette:"light"},{Source:"a",Palette:"light",ThemeHash:"changed-css"}} {
  if _,err:=p.Compile(context.Background(),engine,req);err!=nil{t.Fatal(err)}
 }
 if engine.calls.Load()!=4 { t.Fatalf("cache calls=%d",engine.calls.Load()) }
 engine.version="2"
 if _,err:=p.Compile(context.Background(),engine,a);err!=nil { t.Fatal(err) }
 if engine.calls.Load()!=5 { t.Fatal("renderer version must invalidate cache") }
 broken:=&fakeRenderer{name:"broken",version:"1",fail:true}
 for i:=0;i<2;i++ {
  if _,err:=p.Compile(context.Background(),broken,Request{});err==nil || !strings.Contains(err.Error(),"render failed") {t.Fatalf("error=%v",err)}
 }
 if broken.calls.Load()!=2 {t.Fatal("failures were cached")}
}
