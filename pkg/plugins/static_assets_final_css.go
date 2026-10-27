package plugins

import (
	"crypto/sha256"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/templates"
)

// reconcileFinalCSSAssets reconciles final published CSS with already-rendered
// HTML after purge and minification, including HTML restored from build cache.
func reconcileFinalCSSAssets(m *lifecycle.Manager, outputFS *os.Root, outputDir string) error {
	type replacement struct {
		matcher     *regexp.Regexp
		destination string
	}
	hashes := m.AssetHashes()
	paths := make([]string, 0, len(hashes))
	for path := range hashes {
		if strings.EqualFold(filepath.Ext(path), ".css") {
			paths = append(paths, path)
		}
	}
	sort.Strings(paths)
	changes := make([]replacement, 0, len(paths))
	finalHashes := make(map[string]string)
	for _, path := range paths {
		relative, err := outputRelativePath(outputDir, filepath.Join(outputDir, path))
		if err != nil {
			return err
		}
		data, err := outputFS.ReadFile(relative)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return fmt.Errorf("reading final stylesheet %s: %w", path, err)
		}
		hash := fmt.Sprintf("%x", sha256.Sum256(data))[:8]
		ext := filepath.Ext(path)
		base := strings.TrimSuffix(filepath.ToSlash(path), ext)
		matcher, err := regexp.Compile(regexp.QuoteMeta(base) + `\.[0-9a-fA-F]{8}` + regexp.QuoteMeta(ext))
		if err != nil {
			return err
		}
		changes = append(changes, replacement{matcher: matcher, destination: base + "." + hash + ext})
		if hash != hashes[path] {
			finalHashes[path] = hash
		}
	}
	if len(changes) == 0 {
		return nil
	}

	// Cached HTML may reference either the original Configure hash or a final
	// hash from an earlier build. Normalize every matching URL, not only the
	// source-hash reference from this build.
	err := filepath.WalkDir(outputDir, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			return nil
		}
		if entry.Type()&os.ModeSymlink != 0 || !strings.EqualFold(filepath.Ext(path), ".html") {
			return nil
		}
		relative, err := outputRelativePath(outputDir, path)
		if err != nil {
			return err
		}
		data, err := outputFS.ReadFile(relative)
		if err != nil {
			return err
		}
		before := string(data)
		after := before
		for _, change := range changes {
			after = change.matcher.ReplaceAllString(after, change.destination)
		}
		if before == after {
			return nil
		}
		return outputFS.WriteFile(relative, []byte(after), 0o644)
	})
	if err != nil {
		return fmt.Errorf("updating HTML CSS URLs: %w", err)
	}
	for path, hash := range finalHashes {
		m.SetAssetHash(path, hash)
	}
	templates.SetAssetHashes(finalHashes)
	return nil
}
