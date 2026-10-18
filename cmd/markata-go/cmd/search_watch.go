package cmd

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/fsnotify/fsnotify"
)

var (
	searchWatchInternalPaths = []string{".markata", ".markata-cache", "cache", "markout", "public", "output"}
	searchWatchGlobMagic     = regexp.MustCompile(`[*?\[{]`)
)

func searchContentWatchRoots(config *lifecycle.Config) []string {
	contentDir := config.ContentDir
	if contentDir == "" {
		contentDir = "."
	}

	if contentDir != "." || len(config.GlobPatterns) == 0 {
		return []string{contentDir}
	}

	seen := map[string]struct{}{}
	roots := make([]string, 0, len(config.GlobPatterns))
	for _, pattern := range config.GlobPatterns {
		root := searchWatchRootFromPattern(pattern)
		if root == "" {
			continue
		}
		if _, ok := seen[root]; ok {
			continue
		}
		seen[root] = struct{}{}
		roots = append(roots, root)
	}

	if len(roots) == 0 {
		return []string{contentDir}
	}

	sort.Strings(roots)
	return roots
}

func searchWatchRootFromPattern(pattern string) string {
	cleaned := filepath.Clean(filepath.FromSlash(pattern))
	if cleaned == "." || cleaned == string(filepath.Separator) {
		return "."
	}

	parts := strings.Split(cleaned, string(filepath.Separator))
	rootParts := make([]string, 0, len(parts))
	for _, part := range parts {
		if part == "" || part == "." {
			continue
		}
		if searchWatchGlobMagic.MatchString(part) {
			break
		}
		rootParts = append(rootParts, part)
	}

	if len(rootParts) == 0 {
		return "."
	}
	return filepath.Join(rootParts...)
}

func searchShouldIgnorePath(pathname string) bool {
	absPath, err := filepath.Abs(pathname)
	if err != nil {
		absPath = pathname
	}

	for _, p := range searchWatchInternalPaths {
		absInternal, absErr := filepath.Abs(p)
		if absErr != nil {
			absInternal = p
		}
		if isPathWithinDir(absPath, absInternal) {
			return true
		}
	}

	baseName := filepath.Base(pathname)
	return strings.HasSuffix(pathname, "~") ||
		strings.HasPrefix(baseName, ".") ||
		strings.HasSuffix(pathname, ".swp") ||
		strings.HasSuffix(pathname, ".swo") ||
		strings.HasSuffix(pathname, ".tmp")
}

func searchAddDirRecursive(watcher *fsnotify.Watcher, root string) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		absPath, pathErr := filepath.Abs(path)
		if pathErr != nil {
			absPath = path
		}

		for _, internalPath := range searchWatchInternalPaths {
			absInternal, absErr := filepath.Abs(internalPath)
			if absErr != nil {
				absInternal = internalPath
			}
			if isPathWithinDir(absPath, absInternal) {
				if d.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
		}

		if d.IsDir() && strings.HasPrefix(d.Name(), ".") && d.Name() != "." {
			return filepath.SkipDir
		}

		if d.IsDir() {
			return watcher.Add(path)
		}

		return nil
	})
}

func searchHandleNewDirectory(watcher *fsnotify.Watcher, event fsnotify.Event) {
	if event.Op&fsnotify.Create == 0 {
		return
	}

	info, err := os.Stat(event.Name)
	if err != nil || !info.IsDir() {
		return
	}

	if watchErr := searchAddDirRecursive(watcher, event.Name); watchErr != nil && verbose {
		verbosef("Failed to watch new directory %s: %v", event.Name, watchErr)
	}
}
