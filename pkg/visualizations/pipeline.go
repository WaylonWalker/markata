// Package visualizations provides engine-independent backend SVG rendering and caching.
// It deliberately has no dependency on D2, Vega, Mermaid, or browser JavaScript.
package visualizations

import (
 "context"
 "crypto/sha256"
 "encoding/binary"
 "encoding/hex"
 "errors"
 "fmt"
 "sort"
 "sync"
)

// Request includes every effective input that may change renderer output.
// ThemeHash fingerprints theme and font contents, not just their names.
type Request struct {
 Source string
 Palette string
 Font string
 ThemeHash string
 Options map[string]string
}

type Renderer interface {
 ID() string
 Version() string
 RenderSVG(context.Context, Request) ([]byte, error)
}

// Asset is a ready-to-serve, content-addressed static SVG with layout dimensions.
type Asset struct {
 Path string
 SHA256 string
 SVG []byte
 Width int
 Height int
}

type pending struct {
 done chan struct{}
 asset Asset
 err error
}

// Pipeline memoizes successes and coalesces concurrent renders of one input.
// Failures are never cached.
type Pipeline struct {
 mu sync.Mutex
 entries map[string]Asset
 pending map[string]*pending
}

func (p *Pipeline) Compile(ctx context.Context, engine Renderer, req Request) (Asset, error) {
 if engine == nil { return Asset{}, errors.New("visualization renderer is nil") }
 if engine.ID() == "" || engine.Version() == "" {
  return Asset{}, errors.New("visualization renderer ID and version are required")
 }
 key := renderKey(engine, req)
 p.mu.Lock()
 if cached, ok := p.entries[key]; ok {
  p.mu.Unlock()
  return clone(cached), nil
 }
 if p.pending == nil { p.pending = make(map[string]*pending) }
 if active, ok := p.pending[key]; ok {
  p.mu.Unlock()
  select {
  case <-ctx.Done(): return Asset{}, ctx.Err()
  case <-active.done: return clone(active.asset), active.err
  }
 }
 active := &pending{done: make(chan struct{})}
 p.pending[key] = active
 p.mu.Unlock()

 // Callers must not mutate req concurrently with this method.
 req.Options = copyOptions(req.Options)
 data, err := engine.RenderSVG(ctx, req)
 var asset Asset
 if err == nil {
  asset, err = NewSVGAsset(data)
  if err != nil { err = fmt.Errorf("%s SVG: %w", engine.ID(), err) }
 }
 p.mu.Lock()
 if err == nil {
  if p.entries == nil { p.entries = make(map[string]Asset) }
  p.entries[key] = clone(asset)
 }
 active.asset, active.err = clone(asset), err
 delete(p.pending, key)
 close(active.done)
 p.mu.Unlock()
 return clone(asset), err
}

func clone(a Asset) Asset {
 a.SVG = append([]byte(nil), a.SVG...)
 return a
}

func copyOptions(m map[string]string) map[string]string {
 if m == nil { return nil }
 out := make(map[string]string, len(m))
 for k,v := range m { out[k] = v }
 return out
}

// renderKey uses length-prefixed fields so strings cannot collide at boundaries.
func renderKey(engine Renderer, req Request) string {
 h := sha256.New()
 add := func(value string) {
  var length [8]byte
  binary.BigEndian.PutUint64(length[:], uint64(len(value)))
  _, _ = h.Write(length[:])
  _, _ = h.Write([]byte(value))
 }
 add("markata-visualization-render-v1")
 add(engine.ID())
 add(engine.Version())
 add(req.Source)
 add(req.Palette)
 add(req.Font)
 add(req.ThemeHash)
 keys := make([]string, 0, len(req.Options))
 for k := range req.Options { keys = append(keys, k) }
 sort.Strings(keys)
 for _, k := range keys { add(k); add(req.Options[k]) }
 return hex.EncodeToString(h.Sum(nil))
}
