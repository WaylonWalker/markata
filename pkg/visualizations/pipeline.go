// Package visualizations provides engine-independent static SVG asset handling.
// Rendering engines (D2, Vega-Lite, Mermaid, etc.) implement Renderer; this
// package does not add their dependencies or execute browser JavaScript.
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

// Request contains all presentation inputs that can change rendered SVG bytes.
// Engines must not depend on unstated global theme/configuration when caching.
type Request struct {
	Source    string
	Palette   string
	Font      string
	ThemeHash string // Fingerprint the effective theme/font content, not merely its name.
	Options   map[string]string
}

// Renderer is a backend-only compiler. ID and Version must describe the
// implementation, including layout engine versions that affect its output.
type Renderer interface {
	ID() string
	Version() string
	RenderSVG(ctx context.Context, req Request) ([]byte, error)
}

// Asset represents ready-to-serve SVG. Width and Height reserve layout space
// in the HTML document. Path is a content-addressed site-relative URL path.
type Asset struct {
	Path   string
	SHA256 string
	SVG    []byte
	Width  int
	Height int
}

// Pipeline deduplicates successful renders within a process. It never caches
// failures; separate builds can use the content-addressed published assets.
const (
	maxCacheEntries = 128
	maxCacheBytes   = 32 << 20
)

type Pipeline struct {
	mu         sync.Mutex
	entries    map[string]Asset
	pending    map[string]*pending
	order      []string
	cacheBytes int
}

type pending struct {
	done  chan struct{}
	asset Asset
	err   error
}

// Compile renders once for each renderer/version/source/palette/font/options
// identity, even when identical requests arrive concurrently.
func (p *Pipeline) Compile(ctx context.Context, engine Renderer, req Request) (Asset, error) {
	if engine == nil {
		return Asset{}, errors.New("visualization renderer is nil")
	}
	if engine.ID() == "" || engine.Version() == "" {
		return Asset{}, errors.New("visualization renderer ID and version are required")
	}
	key := renderKey(engine, req)
	p.mu.Lock()
	if cached, ok := p.entries[key]; ok {
		p.mu.Unlock()
		return clone(cached), nil
	}
	if p.pending == nil {
		p.pending = make(map[string]*pending)
	}
	if active, ok := p.pending[key]; ok {
		p.mu.Unlock()
		select {
		case <-ctx.Done():
			return Asset{}, ctx.Err()
		case <-active.done:
			return clone(active.asset), active.err
		}
	}
	active := &pending{done: make(chan struct{})}
	p.pending[key] = active
	p.mu.Unlock()

	// Copy mutable request options so concurrent callers cannot alter them while
	// the engine is rendering. Callers must not mutate req during the call.
	req.Options = copyOptions(req.Options)
	data, err := engine.RenderSVG(ctx, req)
	var asset Asset
	if err == nil {
		asset, err = NewSVGAsset(data)
		if err != nil {
			err = fmt.Errorf("%s SVG: %w", engine.ID(), err)
		}
	}
	p.mu.Lock()
	if err == nil {
		if p.entries == nil {
			p.entries = make(map[string]Asset)
		}
		if len(asset.SVG) <= maxCacheBytes {
			for len(p.order) >= maxCacheEntries || p.cacheBytes+len(asset.SVG) > maxCacheBytes {
				oldest := p.order[0]
				p.order = p.order[1:]
				p.cacheBytes -= len(p.entries[oldest].SVG)
				delete(p.entries, oldest)
			}
			p.entries[key] = clone(asset)
			p.order = append(p.order, key)
			p.cacheBytes += len(asset.SVG)
		}
	}
	active.asset, active.err = clone(asset), err
	delete(p.pending, key)
	close(active.done)
	p.mu.Unlock()
	return clone(asset), err
}

func clone(asset Asset) Asset {
	asset.SVG = append([]byte(nil), asset.SVG...)
	return asset
}

func copyOptions(source map[string]string) map[string]string {
	if source == nil {
		return nil
	}
	result := make(map[string]string, len(source))
	for k, v := range source {
		result[k] = v
	}
	return result
}

// renderKey hashes length-prefixed fields to avoid delimiter collisions and
// sorts option names so map iteration does not affect cache identity.
func renderKey(engine Renderer, req Request) string {
	h := sha256.New()
	add := func(v string) {
		var length [8]byte
		binary.BigEndian.PutUint64(length[:], uint64(len(v)))
		_, _ = h.Write(length[:])
		_, _ = h.Write([]byte(v))
	}
	add("markata-visualization-render-v1")
	add(engine.ID())
	add(engine.Version())
	add(req.Source)
	add(req.Palette)
	add(req.Font)
	add(req.ThemeHash)
	keys := make([]string, 0, len(req.Options))
	for key := range req.Options {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		add(key)
		add(req.Options[key])
	}
	return hex.EncodeToString(h.Sum(nil))
}
