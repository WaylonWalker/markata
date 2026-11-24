// Package templates provides template engine functionality for markata-go.
package templates

import (
	"sync"

	"github.com/WaylonWalker/markata/pkg/models"
)

// PostMapCache provides thread-safe caching for post-to-map conversions.
// This dramatically reduces memory allocations and CPU time when the same
// posts are converted multiple times (e.g., in feed pages, templates).
type PostMapCache struct {
	mu    sync.RWMutex
	cache map[*models.Post]map[string]interface{}
}

// Global cache instance for post map conversions.
// This is safe for concurrent access.
var globalPostMapCache = &PostMapCache{
	cache: make(map[*models.Post]map[string]interface{}),
}

// GetPostMap returns the cached map for a post, or converts and caches it.
// This is the primary entry point for getting post maps with caching.
func GetPostMap(p *models.Post) map[string]interface{} {
	if p == nil {
		return nil
	}
	return globalPostMapCache.GetOrCreate(p)
}

// GetOrCreate returns a cached map for the post, or creates and caches one.
func (c *PostMapCache) GetOrCreate(p *models.Post) map[string]interface{} {
	// Fast path: check if already cached (read lock)
	c.mu.RLock()
	if cached, ok := c.cache[p]; ok {
		c.mu.RUnlock()
		return cached
	}
	c.mu.RUnlock()

	// Slow path: need to create and cache (write lock)
	c.mu.Lock()
	defer c.mu.Unlock()

	// Double-check after acquiring write lock
	if cached, ok := c.cache[p]; ok {
		return cached
	}

	// Create the map (using the uncached conversion function)
	m := postToMapUncached(p)
	c.cache[p] = m
	return m
}

// Clear clears the entire cache. Call this between builds or when posts change.
func (c *PostMapCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache = make(map[*models.Post]map[string]interface{})
}

// Invalidate removes a specific post from the cache.
// Call this when a post is modified.
func (c *PostMapCache) Invalidate(p *models.Post) {
	if p == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.cache, p)
}

// Size returns the number of cached entries.
func (c *PostMapCache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.cache)
}

// ClearPostMapCache clears the global post map cache.
// This should be called at the start of each build to ensure fresh data.
func ClearPostMapCache() {
	globalPostMapCache.Clear()
}

// InvalidatePost removes a specific post from the global cache.
func InvalidatePost(p *models.Post) {
	globalPostMapCache.Invalidate(p)
}

// PostMapCacheSize returns the size of the global cache.
func PostMapCacheSize() int {
	return globalPostMapCache.Size()
}

// ConfigMapCache provides thread-safe caching for config-to-map conversions.
type ConfigMapCache struct {
	mu    sync.RWMutex
	cache map[*models.Config]map[string]interface{}
}

// Global cache instance for config map conversions.
var globalConfigMapCache = &ConfigMapCache{
	cache: make(map[*models.Config]map[string]interface{}),
}

// GetConfigMap returns the cached map for a config, or converts and caches it.
func GetConfigMap(c *models.Config) map[string]interface{} {
	if c == nil {
		return nil
	}
	return globalConfigMapCache.GetOrCreate(c)
}

// GetOrCreate returns a cached map for the config, or creates and caches one.
func (c *ConfigMapCache) GetOrCreate(cfg *models.Config) map[string]interface{} {
	c.mu.RLock()
	if cached, ok := c.cache[cfg]; ok {
		c.mu.RUnlock()
		return cached
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()

	if cached, ok := c.cache[cfg]; ok {
		return cached
	}

	m := configToMap(cfg)
	c.cache[cfg] = m
	return m
}

// Clear clears the config cache.
func (c *ConfigMapCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache = make(map[*models.Config]map[string]interface{})
}

// ClearConfigMapCache clears the global config map cache.
func ClearConfigMapCache() {
	globalConfigMapCache.Clear()
}

// ClearAllCaches clears all template caches.
// Call this at the start of each build.
func ClearAllCaches() {
	ClearPostMapCache()
	ClearConfigMapCache()
}
