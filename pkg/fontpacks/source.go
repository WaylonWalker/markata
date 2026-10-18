package fontpacks

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/WaylonWalker/markata/internal/fontcatalog"
	"gopkg.in/yaml.v3"
)

// BuiltinSource returns the catalog and assets embedded in markata-go.
func BuiltinSource() (*CatalogSource, error) {
	data, err := fs.ReadFile(fontcatalog.FS(), "markata-fontpacks.yaml")
	if err != nil {
		return nil, fmt.Errorf("read embedded font catalog: %w", err)
	}
	c, err := loadBytes(data)
	if err != nil {
		return nil, err
	}
	digest, err := fontcatalog.Fingerprint()
	if err != nil {
		return nil, fmt.Errorf("fingerprint embedded font catalog: %w", err)
	}
	return &CatalogSource{
		Catalog:       c,
		FS:            fontcatalog.FS(),
		Root:          ".",
		LockFS:        fontcatalog.FS(),
		Lock:          "markata-fonts.lock.yaml",
		Builtin:       true,
		ContentDigest: digest,
	}, nil
}

// Open opens an asset from this source. It is useful to callers that need to
// inspect or copy an asset without knowing whether it is embedded.
func (s *CatalogSource) Open(name string) (fs.File, error) {
	return s.FS.Open(filepath.ToSlash(filepath.Join(s.Root, name)))
}

func loadBytes(data []byte) (*Catalog, error) {
	var c Catalog
	if err := yaml.Unmarshal(data, &c); err != nil {
		return nil, fmt.Errorf("parse font catalog: %w", err)
	}
	if err := c.Validate(); err != nil {
		return nil, err
	}
	return &c, nil
}

func loadLock(s *CatalogSource) (Lockfile, error) {
	data, err := fs.ReadFile(s.LockFS, s.Lock)
	if err != nil {
		return Lockfile{}, fmt.Errorf("read font lockfile: %w", err)
	}
	var lock Lockfile
	if err := yaml.Unmarshal(data, &lock); err != nil {
		return Lockfile{}, fmt.Errorf("parse font lockfile: %w", err)
	}
	return lock, nil
}
