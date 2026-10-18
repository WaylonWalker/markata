package lifecycle

import (
	"path/filepath"

	"github.com/WaylonWalker/markata/pkg/models"
)

func normalizeServePath(path string) string {
	return filepath.ToSlash(filepath.Clean(path))
}

// Serve incremental rebuild cache keys.
const (
	CacheKeyServeChangedPaths = "serve.changed_paths"
	CacheKeyServeAffectedPath = "serve.affected_paths"
	CacheKeyServeFullRebuild  = "serve.full_rebuild"
	CacheKeyServeCachedPosts  = "serve.cached_posts"
	CacheKeyServeRemovedPaths = "serve.removed_paths"
	CacheKeyServeGlobDirty    = "serve.glob_dirty"
	CacheKeyServeIncremental  = "serve.incremental"
)

// SetServeChangedPaths stores normalized, relative content paths that changed.
func SetServeChangedPaths(m *Manager, paths []string) {
	if m == nil {
		return
	}
	if len(paths) == 0 {
		m.Cache().Delete(CacheKeyServeChangedPaths)
		return
	}
	clean := make([]string, 0, len(paths))
	seen := make(map[string]struct{}, len(paths))
	for _, p := range paths {
		if p == "" {
			continue
		}
		cp := normalizeServePath(p)
		if _, ok := seen[cp]; ok {
			continue
		}
		seen[cp] = struct{}{}
		clean = append(clean, cp)
	}
	m.Cache().Set(CacheKeyServeChangedPaths, clean)
}

// SetServeIncremental stores whether the current rebuild may reuse unchanged
// posts without enabling fast-mode output skips.
func SetServeIncremental(m *Manager, enabled bool) {
	if m == nil {
		return
	}
	m.Cache().Set(CacheKeyServeIncremental, enabled)
}

// IsServeIncremental reports whether the current rebuild may reuse cached
// post results. Fast mode also enables this behavior.
func IsServeIncremental(m *Manager) bool {
	if m == nil {
		return false
	}
	if IsServeFastMode(m) {
		return true
	}
	if raw, ok := m.Cache().Get(CacheKeyServeIncremental); ok {
		if enabled, ok := raw.(bool); ok {
			return enabled
		}
	}
	return false
}

// IsServeIncrementalFromConfig is the config-only form used by plugins during
// output cleanup.
func IsServeIncrementalFromConfig(config *Config) bool {
	if config == nil {
		return false
	}
	if IsServeFastModeFromConfig(config) {
		return true
	}
	if enabled, ok := config.Extra["incremental_mode"].(bool); ok {
		return enabled
	}
	return false
}

// GetServeChangedPaths returns the changed content paths if present.
func GetServeChangedPaths(m *Manager) []string {
	if m == nil {
		return nil
	}
	if raw, ok := m.Cache().Get(CacheKeyServeChangedPaths); ok {
		if paths, ok := raw.([]string); ok {
			return paths
		}
	}
	return nil
}

// SetServeAffectedPaths stores a lookup map of content paths that should be rebuilt.
func SetServeAffectedPaths(m *Manager, paths map[string]bool) {
	if m == nil {
		return
	}
	if len(paths) == 0 {
		m.Cache().Delete(CacheKeyServeAffectedPath)
		return
	}
	m.Cache().Set(CacheKeyServeAffectedPath, paths)
}

// GetServeAffectedPaths returns a lookup map of paths to rebuild.
func GetServeAffectedPaths(m *Manager) map[string]bool {
	if m == nil {
		return nil
	}
	if raw, ok := m.Cache().Get(CacheKeyServeAffectedPath); ok {
		if paths, ok := raw.(map[string]bool); ok {
			return paths
		}
	}
	return nil
}

// SetServeFullRebuild stores whether the current serve rebuild should be full.
func SetServeFullRebuild(m *Manager, full bool) {
	if m == nil {
		return
	}
	m.Cache().Set(CacheKeyServeFullRebuild, full)
}

// IsServeFullRebuild returns whether the current serve rebuild is full.
func IsServeFullRebuild(m *Manager) bool {
	if m == nil {
		return false
	}
	if raw, ok := m.Cache().Get(CacheKeyServeFullRebuild); ok {
		if full, ok := raw.(bool); ok {
			return full
		}
	}
	return false
}

// SetServeRemovedPaths stores normalized, relative content paths that were removed.
func SetServeRemovedPaths(m *Manager, paths []string) {
	if m == nil {
		return
	}
	if len(paths) == 0 {
		m.Cache().Delete(CacheKeyServeRemovedPaths)
		return
	}
	clean := make([]string, 0, len(paths))
	seen := make(map[string]struct{}, len(paths))
	for _, p := range paths {
		if p == "" {
			continue
		}
		cp := normalizeServePath(p)
		if _, ok := seen[cp]; ok {
			continue
		}
		seen[cp] = struct{}{}
		clean = append(clean, cp)
	}
	m.Cache().Set(CacheKeyServeRemovedPaths, clean)
}

// GetServeRemovedPaths returns removed content paths if present.
func GetServeRemovedPaths(m *Manager) []string {
	if m == nil {
		return nil
	}
	if raw, ok := m.Cache().Get(CacheKeyServeRemovedPaths); ok {
		if paths, ok := raw.([]string); ok {
			return paths
		}
	}
	return nil
}

// SetServeGlobDirty stores whether the glob cache should be refreshed.
func SetServeGlobDirty(m *Manager, dirty bool) {
	if m == nil {
		return
	}
	if !dirty {
		m.Cache().Delete(CacheKeyServeGlobDirty)
		return
	}
	m.Cache().Set(CacheKeyServeGlobDirty, true)
}

// IsServeGlobDirty returns true if glob cache should be refreshed.
func IsServeGlobDirty(m *Manager) bool {
	if m == nil {
		return false
	}
	if raw, ok := m.Cache().Get(CacheKeyServeGlobDirty); ok {
		if dirty, ok := raw.(bool); ok {
			return dirty
		}
	}
	return false
}

// IsServeFastMode returns true when fast_mode is enabled in config.
func IsServeFastMode(m *Manager) bool {
	if m == nil {
		return false
	}
	if extra := m.Config().Extra; extra != nil {
		if fast, ok := extra["fast_mode"].(bool); ok {
			return fast
		}
	}
	return false
}

// IsServeFastModeFromConfig returns true when fast_mode is enabled on the config.
func IsServeFastModeFromConfig(config *Config) bool {
	if config == nil {
		return false
	}
	if extra := config.Extra; extra != nil {
		if fast, ok := extra["fast_mode"].(bool); ok {
			return fast
		}
	}
	return false
}

// SetServeCachedPosts stores cached posts for incremental serve rebuilds.
func SetServeCachedPosts(m *Manager, posts map[string]*models.Post) {
	if m == nil {
		return
	}
	if len(posts) == 0 {
		m.Cache().Delete(CacheKeyServeCachedPosts)
		return
	}
	m.Cache().Set(CacheKeyServeCachedPosts, posts)
}

// GetServeCachedPosts returns cached posts for incremental serve rebuilds.
func GetServeCachedPosts(m *Manager) map[string]*models.Post {
	if m == nil {
		return nil
	}
	if raw, ok := m.Cache().Get(CacheKeyServeCachedPosts); ok {
		if posts, ok := raw.(map[string]*models.Post); ok {
			return posts
		}
	}
	return nil
}
