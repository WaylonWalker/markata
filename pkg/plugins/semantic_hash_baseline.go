package plugins

import (
	"sync"

	"github.com/WaylonWalker/markata/pkg/buildcache"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
)

const semanticHashBaselineKey = "load_semantic_hash_baseline"

// semanticHashes contains only existing one-way cache hashes, never title,
// plaintext content, or key material.
type semanticHashes struct {
	feed, tag, garden string
}

// semanticHashBaseline is a single Load -> InlineTitles handoff. It is local
// to a manager and bound to the cache whose metadata Load is about to update.
type semanticHashBaseline struct {
	mu     sync.Mutex
	cache  *buildcache.Cache
	hashes map[string]*semanticHashes
}

func resetSemanticHashBaseline(m *lifecycle.Manager) {
	m.Cache().Set(semanticHashBaselineKey, &semanticHashBaseline{
		cache: GetBuildCache(m), hashes: make(map[string]*semanticHashes),
	})
}

func captureSemanticHashBaseline(m *lifecycle.Manager, cache *buildcache.Cache, path string) {
	value, ok := m.Cache().Get(semanticHashBaselineKey)
	if !ok || cache == nil {
		return
	}
	baseline, ok := value.(*semanticHashBaseline)
	if !ok || baseline.cache != cache {
		return
	}
	baseline.mu.Lock()
	defer baseline.mu.Unlock()
	if _, captured := baseline.hashes[path]; captured {
		return
	}
	feed, tag, garden, available := cache.GetPostSemanticHashBaseline(path)
	if !available {
		// Keep a first-capture marker without retaining incomplete metadata.
		// A later capture must not mistake Load's own hashes for prior values.
		baseline.hashes[path] = nil
		return
	}
	baseline.hashes[path] = &semanticHashes{feed: feed, tag: tag, garden: garden}
}

// Taking the handoff before any renderer setup also discards it on errors.
// InlineTitles reads the immutable snapshot only after Load workers have joined.
func takeSemanticHashBaseline(m *lifecycle.Manager) map[string]*semanticHashes {
	value, ok := m.Cache().Get(semanticHashBaselineKey)
	m.Cache().Delete(semanticHashBaselineKey)
	if !ok {
		return nil
	}
	baseline, ok := value.(*semanticHashBaseline)
	if !ok || baseline.cache != GetBuildCache(m) {
		return nil
	}
	baseline.mu.Lock()
	defer baseline.mu.Unlock()
	hashes := baseline.hashes
	baseline.hashes = nil
	return hashes
}
