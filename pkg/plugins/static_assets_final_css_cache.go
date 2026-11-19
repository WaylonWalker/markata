package plugins

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/WaylonWalker/markata/internal/fileidentity"
	"github.com/WaylonWalker/markata/pkg/buildcache"
	"github.com/WaylonWalker/markata/pkg/logging"
)

const finalCSSLinkCacheLimit = 16 << 20

type finalCSSLinkCache struct {
	Rules    string
	Observed time.Time
	Files    map[string]fileidentity.Identity
}

func finalCSSLinkCachePath(contentDir, outputDir string) string {
	if contentDir == "" || !fileidentity.Available {
		return ""
	}
	output, err := filepath.Abs(outputDir)
	if err != nil {
		return ""
	}
	realOutput, err := filepath.EvalSymlinks(output)
	if err != nil {
		return ""
	}
	content, err := filepath.Abs(contentDir)
	if err != nil {
		return ""
	}
	content, err = filepath.EvalSymlinks(content)
	if err != nil {
		return ""
	}
	cacheRoot := filepath.Join(content, ".markata-cache")
	if resolved, err := filepath.EvalSymlinks(cacheRoot); err == nil {
		cacheRoot = resolved
	} else if !os.IsNotExist(err) {
		return ""
	}
	cacheRoot = filepath.Join(cacheRoot, "final-css-links-v1")
	if resolved, err := filepath.EvalSymlinks(cacheRoot); err == nil {
		cacheRoot = resolved
	} else if !os.IsNotExist(err) {
		return ""
	}
	cachePath := filepath.Join(cacheRoot, buildcache.ContentHash(realOutput)+".json")
	relative, err := filepath.Rel(realOutput, cachePath)
	if err == nil && filepath.IsLocal(relative) {
		return ""
	}
	return cachePath
}

func finalCSSLinkRules(changes []finalCSSReplacement) string {
	var rules strings.Builder
	for _, change := range changes {
		rules.WriteString(change.matcher.String())
		rules.WriteByte(0)
		rules.WriteString(change.destination)
		rules.WriteByte(0)
	}
	return buildcache.ContentHash(rules.String())
}

func loadFinalCSSLinkCache(path, rules string, now time.Time) finalCSSLinkCache {
	if path == "" {
		return finalCSSLinkCache{}
	}
	file, err := os.Open(path)
	if err != nil {
		return finalCSSLinkCache{}
	}
	defer file.Close()

	data, err := io.ReadAll(io.LimitReader(file, finalCSSLinkCacheLimit+1))
	var cached finalCSSLinkCache
	if err != nil || len(data) > finalCSSLinkCacheLimit || json.Unmarshal(data, &cached) != nil ||
		cached.Rules != rules || cached.Observed.IsZero() || now.Before(cached.Observed) {
		return finalCSSLinkCache{}
	}
	return cached
}

func saveFinalCSSLinkCache(path string, cached finalCSSLinkCache) {
	if path == "" {
		return
	}
	data, err := json.Marshal(cached)
	if err != nil || len(data) > finalCSSLinkCacheLimit || os.MkdirAll(filepath.Dir(path), 0o700) != nil {
		return
	}
	file, err := os.CreateTemp(filepath.Dir(path), ".verified-*")
	if err != nil {
		return
	}
	defer os.Remove(file.Name())

	_, writeErr := file.Write(data)
	closeErr := file.Close()
	if writeErr == nil && closeErr == nil {
		if err := replaceDiagnosticsArtifact(file.Name(), path); err != nil {
			logging.Component("static_assets").Phase("cleanup").Warnf("could not update final CSS link cache: %v", err)
		}
	}
}
