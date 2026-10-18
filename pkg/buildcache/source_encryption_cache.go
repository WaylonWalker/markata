package buildcache

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/WaylonWalker/markata/pkg/encryption"
)

const (
	SourceEncryptionCacheFileName = "source-encryption-cache.json"
	sourceEncryptionCacheVersion  = 1
)

// SourceEncryptionCache preserves authenticated ciphertext across a
// decrypt/edit/encrypt cycle. It deliberately stores no plaintext or
// password-derived verifier.
type SourceEncryptionCache struct {
	mu      sync.RWMutex
	path    string
	Version int                              `json:"version"`
	Entries map[string]SourceEncryptionEntry `json:"entries"`
}

// SourceEncryptionEntry is ciphertext associated with one source path.
type SourceEncryptionEntry struct {
	KeyName       string `json:"key_name"`
	EncryptedBody string `json:"encrypted_body"`
}

// LoadSourceEncryptionCache loads the dedicated source-encryption cache.
// Missing or malformed caches are treated as empty caches and returned as a
// warning error so callers can continue without sacrificing correctness.
func LoadSourceEncryptionCache(cacheDir string) (*SourceEncryptionCache, error) {
	if cacheDir == "" {
		cacheDir = DefaultCacheDir
	}
	cache := &SourceEncryptionCache{
		path:    filepath.Join(cacheDir, SourceEncryptionCacheFileName),
		Version: sourceEncryptionCacheVersion,
		Entries: make(map[string]SourceEncryptionEntry),
	}
	data, err := os.ReadFile(cache.path)
	if err != nil {
		if os.IsNotExist(err) {
			return cache, nil
		}
		return cache, fmt.Errorf("read source encryption cache: %w", err)
	}
	var disk SourceEncryptionCache
	if err := json.Unmarshal(data, &disk); err != nil || disk.Version != sourceEncryptionCacheVersion {
		if err == nil {
			err = fmt.Errorf("unsupported version %d", disk.Version)
		}
		return cache, fmt.Errorf("parse source encryption cache: %w", err)
	}
	cache.Entries = disk.Entries
	if cache.Entries == nil {
		cache.Entries = make(map[string]SourceEncryptionEntry)
	}
	return cache, nil
}

// Get returns a cache hit only after AES-GCM authentication and an exact
// plaintext comparison with the current source body succeed.
func (c *SourceEncryptionCache) Get(path, plaintext, keyName, password string) (string, bool) {
	c.mu.RLock()
	entry, ok := c.Entries[path]
	c.mu.RUnlock()
	if !ok || !strings.EqualFold(strings.TrimSpace(entry.KeyName), strings.TrimSpace(keyName)) {
		return "", false
	}
	decrypted, resolvedKey, err := encryption.DecryptSourceMarkdown(entry.EncryptedBody, password)
	if err != nil || !strings.EqualFold(strings.TrimSpace(resolvedKey), strings.TrimSpace(keyName)) || decrypted != plaintext {
		return "", false
	}
	return entry.EncryptedBody, true
}

// Put records an authenticated ciphertext for a source path.
func (c *SourceEncryptionCache) Put(path, keyName, encryptedBody string) {
	if path == "" || keyName == "" || encryptedBody == "" {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Entries[path] = SourceEncryptionEntry{KeyName: keyName, EncryptedBody: encryptedBody}
}

// Save atomically persists the cache with owner-only permissions.
func (c *SourceEncryptionCache) Save() error {
	c.mu.RLock()
	disk := struct {
		Version int                              `json:"version"`
		Entries map[string]SourceEncryptionEntry `json:"entries"`
	}{Version: sourceEncryptionCacheVersion, Entries: make(map[string]SourceEncryptionEntry, len(c.Entries))}
	for path, entry := range c.Entries {
		disk.Entries[path] = entry
	}
	c.mu.RUnlock()

	data, err := json.Marshal(disk)
	if err != nil {
		return fmt.Errorf("marshal source encryption cache: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(c.path), 0o700); err != nil {
		return fmt.Errorf("create source encryption cache directory: %w", err)
	}
	temporary, err := os.CreateTemp(filepath.Dir(c.path), ".source-encryption-cache-*")
	if err != nil {
		return fmt.Errorf("create source encryption cache temporary file: %w", err)
	}
	temporaryPath := temporary.Name()
	defer os.Remove(temporaryPath)
	if err := temporary.Chmod(0o600); err != nil {
		temporary.Close()
		return err
	}
	if _, err := temporary.Write(data); err != nil {
		temporary.Close()
		return fmt.Errorf("write source encryption cache: %w", err)
	}
	if err := temporary.Sync(); err != nil {
		temporary.Close()
		return fmt.Errorf("sync source encryption cache: %w", err)
	}
	if err := temporary.Close(); err != nil {
		return fmt.Errorf("close source encryption cache: %w", err)
	}
	if err := os.Rename(temporaryPath, c.path); err != nil {
		return fmt.Errorf("replace source encryption cache: %w", err)
	}
	return nil
}
