package plugins

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"syscall"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
)

const (
	assetMinifyRevision = "3"
	assetMinifyVersion  = "v2.24.21"
	assetParseVersion   = "v2.8.18"
)

var errUnsafeMinifyCacheLocation = errors.New("unsafe asset cache")

type minifyRecord struct {
	Version    int    `json:"version"`
	Scope      string `json:"scope"`
	Asset      string `json:"asset"`
	Recipe     string `json:"recipe"`
	OutputHash string `json:"output_hash"`
	SourceHash string `json:"source_hash,omitempty"`
}

type minifyCache struct {
	dir      string
	output   string
	scope    string
	root     *os.Root
	setupErr error
}

func minifyHash(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

func validMinifyHash(hash string) bool {
	if len(hash) != 64 {
		return false
	}
	for _, c := range hash {
		if !strings.ContainsRune("0123456789abcdef", c) {
			return false
		}
	}
	return true
}

func minifyRecipe(kind string, comments []string) string {
	// Comment matching is a set of substring predicates, independent of order.
	options := append([]string{}, comments...)
	sort.Strings(options)
	unique := options[:0]
	for _, option := range options {
		if len(unique) == 0 || unique[len(unique)-1] != option {
			unique = append(unique, option)
		}
	}
	recipe := []byte(kind + "\x00" + assetMinifyRevision + "\x00" + assetMinifyVersion + "\x00" + assetParseVersion + "\x00")
	for _, option := range unique {
		recipe = strconv.AppendQuote(recipe, option)
		recipe = append(recipe, 0)
	}
	return minifyHash(recipe)
}

// resolveMinifyPath resolves existing ancestors even when the final directory
// has not yet been created.
func resolveMinifyPath(path string) (string, error) {
	absolute, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	current := absolute
	var suffix []string
	for {
		resolved, resolveErr := filepath.EvalSymlinks(current)
		if resolveErr == nil {
			for i := len(suffix) - 1; i >= 0; i-- {
				resolved = filepath.Join(resolved, suffix[i])
			}
			return resolved, nil
		}
		if (!os.IsNotExist(resolveErr) && !errors.Is(resolveErr, syscall.ENOTDIR)) || filepath.Dir(current) == current {
			return "", resolveErr
		}
		// A dangling symlink is not a nonexistent plain directory.
		if info, statErr := os.Lstat(current); statErr == nil && info.Mode()&os.ModeSymlink != 0 {
			return "", resolveErr
		}
		suffix = append(suffix, filepath.Base(current))
		current = filepath.Dir(current)
	}
}

func minifyWithin(root, path string) bool {
	rel, err := filepath.Rel(root, path)
	return err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}

func minifyPathsOverlap(first, second string) bool {
	return minifyWithin(first, second) || minifyWithin(second, first)
}

func validateMinifyCacheLocation(base, dir, output string) error {
	resolvedOutput, err := resolveMinifyPath(output)
	if err != nil {
		return fmt.Errorf("resolve output: %w", err)
	}
	privateDirs := []string{
		filepath.Join(base, "asset-minify"), filepath.Join(base, "asset-minify", "v1"),
		dir, filepath.Join(dir, "records"), filepath.Join(dir, "blobs"),
	}
	for _, path := range privateDirs {
		resolved, err := resolveMinifyPath(path)
		if err != nil {
			return fmt.Errorf("resolve cache: %w", err)
		}
		if minifyPathsOverlap(output, path) || minifyPathsOverlap(resolvedOutput, resolved) {
			return fmt.Errorf("%w overlaps published output: %s", errUnsafeMinifyCacheLocation, path)
		}
	}
	// MkdirAll does not chmod existing base ancestors, but any missing ones
	// would be created private and must not contain the published output.
	for path := base; ; path = filepath.Dir(path) {
		_, err := os.Lstat(path)
		if err == nil {
			return nil
		}
		if !os.IsNotExist(err) {
			// On Windows ENOTDIR aliases ERROR_PATH_NOT_FOUND, so missing
			// ancestors must reach overlap validation before this fallback.
			if errors.Is(err, syscall.ENOTDIR) {
				return nil
			}
			return fmt.Errorf("inspect cache ancestor: %w", err)
		}
		resolved, err := resolveMinifyPath(path)
		if err != nil {
			return fmt.Errorf("resolve cache ancestor: %w", err)
		}
		if minifyPathsOverlap(output, path) || minifyPathsOverlap(resolvedOutput, resolved) {
			return fmt.Errorf("%w ancestor overlaps published output: %s", errUnsafeMinifyCacheLocation, path)
		}
	}
}

func newMinifyCache(config *lifecycle.Config, kind string) (*minifyCache, error) {
	site, err := filepath.Abs(config.ContentDir)
	if err != nil {
		return nil, err
	}
	output, err := filepath.Abs(config.OutputDir)
	if err != nil {
		return nil, err
	}
	base := filepath.Join(config.ContentDir, ".markata")
	if override, ok := config.Extra["cache_dir"].(string); ok && override != "" {
		base = override
	}
	base, err = filepath.Abs(base)
	if err != nil {
		return nil, err
	}
	scope := minifyHash([]byte(site + "\x00" + output + "\x00" + kind))
	dir := filepath.Join(base, "asset-minify", "v1", scope)
	cache := &minifyCache{dir: dir, output: output, scope: scope}
	if err := validateMinifyCacheLocation(base, dir, output); err != nil {
		if errors.Is(err, errUnsafeMinifyCacheLocation) {
			return nil, err
		}
		cache.setupErr = err
		return cache, err
	}
	// Refuse internal symlinks, including ones leading back into published output.
	if err := os.MkdirAll(base, 0o700); err != nil {
		cache.setupErr = err
		return cache, err
	}
	baseRoot, err := os.OpenRoot(base)
	if err != nil {
		cache.setupErr = err
		return cache, err
	}
	defer baseRoot.Close()
	for _, relative := range []string{"asset-minify", filepath.Join("asset-minify", "v1"), filepath.Join("asset-minify", "v1", scope)} {
		if err := privateMinifyDir(baseRoot, relative); err != nil {
			cache.setupErr = err
			return cache, err
		}
	}
	cache.root, err = os.OpenRoot(dir)
	if err == nil {
		for _, relative := range []string{"records", "blobs"} {
			if err = privateMinifyDir(cache.root, relative); err != nil {
				break
			}
		}
	}
	cache.setupErr = err
	return cache, err
}

func privateMinifyDir(root *os.Root, path string) error {
	if err := root.Mkdir(path, 0o700); err != nil && !os.IsExist(err) {
		return err
	}
	info, err := root.Lstat(path)
	if err != nil {
		return err
	}
	if !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("unsafe asset cache directory: %s", path)
	}
	if info.Mode().Perm() != 0o700 {
		return root.Chmod(path, 0o700)
	}
	return nil
}

func (c *minifyCache) close() {
	if c != nil && c.root != nil {
		_ = c.root.Close()
	}
}

func (c *minifyCache) identity(path string) (string, error) {
	absolute, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	relative, err := filepath.Rel(c.output, absolute)
	if err != nil || !minifyWithin(c.output, absolute) || relative == "." {
		return "", fmt.Errorf("asset is outside output: %s", path)
	}
	if err := validateOutputPath(c.output, absolute); err != nil {
		return "", err
	}
	return filepath.ToSlash(relative), nil
}

func minifyRecordPath(asset string) string {
	return filepath.Join("records", minifyHash([]byte(asset))+".json")
}

func (c *minifyCache) readRecord(asset string) (*minifyRecord, error) {
	if c.root == nil {
		return nil, c.setupErr
	}
	data, err := c.root.ReadFile(minifyRecordPath(asset))
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var record minifyRecord
	if err := json.Unmarshal(data, &record); err != nil {
		return nil, fmt.Errorf("invalid asset record: %w", err)
	}
	if record.Version != 1 || record.Scope != c.scope || record.Asset != asset ||
		!validMinifyHash(record.Recipe) || !validMinifyHash(record.OutputHash) ||
		(record.SourceHash != "" && !validMinifyHash(record.SourceHash)) {
		return nil, errors.New("invalid asset record identity, schema or digest")
	}
	return &record, nil
}

func (c *minifyCache) readBlob(hash string) ([]byte, error) {
	if !validMinifyHash(hash) {
		return nil, errors.New("invalid asset blob digest")
	}
	data, err := c.root.ReadFile(filepath.Join("blobs", hash))
	if err != nil {
		return nil, err
	}
	if minifyHash(data) != hash {
		return nil, errors.New("asset blob checksum mismatch")
	}
	return data, nil
}

func (c *minifyCache) persist(record *minifyRecord, result []byte) error {
	if c.setupErr != nil {
		return c.setupErr
	}
	if err := c.storeBlob(record.OutputHash, result); err != nil {
		return err
	}
	data, err := json.Marshal(record)
	if err != nil {
		return err
	}
	return c.writePrivate(minifyRecordPath(record.Asset), data)
}

func (c *minifyCache) storeBlob(hash string, data []byte) error {
	if !validMinifyHash(hash) {
		return errors.New("invalid asset payload digest")
	}
	if _, err := c.readBlob(hash); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		log.Printf("[asset_minify] Warning: repairing cached blob %s: %v", hash, err)
		// Repair corrupt payloads without truncating potential hard links.
		// The caller reports corruption when a payload was needed for reuse.
		info, statErr := c.root.Lstat(filepath.Join("blobs", hash))
		if statErr != nil || !info.Mode().IsRegular() {
			return err
		}
	}
	return c.writePrivate(filepath.Join("blobs", hash), data)
}

func (c *minifyCache) writePrivate(path string, data []byte) error {
	var random [16]byte
	if _, err := rand.Read(random[:]); err != nil {
		return err
	}
	temp := filepath.Join(filepath.Dir(path), ".tmp-"+hex.EncodeToString(random[:]))
	file, err := c.root.OpenFile(temp, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}
	defer c.root.Remove(temp) //nolint:errcheck // best-effort temporary cleanup
	if _, err := file.Write(data); err != nil {
		_ = file.Close()
		return err
	}
	if err := file.Sync(); err != nil {
		_ = file.Close()
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	if err := c.root.Rename(temp, path); err != nil {
		return err
	}
	return syncMinifyDir(c.root, filepath.Dir(path))
}

func syncMinifyDir(root *os.Root, path string) error {
	if runtime.GOOS == "windows" { //nolint:goconst // Cache durability is independent of the Tailwind plugin.
		return nil // Windows does not support syncing directory handles.
	}
	dir, err := root.Open(path)
	if err != nil {
		return err
	}
	defer dir.Close()
	return dir.Sync()
}

func (c *minifyCache) invalidate(asset string) error {
	if c.root == nil {
		// Absence is safe; inaccessible existing provenance is not.
		_, err := os.Lstat(c.dir)
		if os.IsNotExist(err) || errors.Is(err, syscall.ENOTDIR) {
			return nil
		}
		return fmt.Errorf("cannot safely invalidate asset cache: %w", c.setupErr)
	}
	path := minifyRecordPath(asset)
	if err := c.root.Remove(path); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return syncMinifyDir(c.root, "records")
}
