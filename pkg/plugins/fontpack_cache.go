package plugins

import (
	"encoding/json"
	"fmt"
	"slices"

	"github.com/WaylonWalker/markata/pkg/buildcache"
	"github.com/WaylonWalker/markata/pkg/fontpacks"
)

func fontpackCacheKey(coverage, defaultName string, names []string, pickerEnabled bool, source *fontpacks.CatalogSource) (string, error) {
	if source == nil || source.Catalog == nil {
		return "", fmt.Errorf("fontpack cache identity requires a catalog")
	}
	if source.Builtin && source.ContentDigest == "" {
		return "", fmt.Errorf("fontpack cache identity requires bundled metadata digest")
	}
	resolvedDefault, _, err := source.Catalog.ResolvePack(defaultName)
	if err != nil {
		return "", fmt.Errorf("resolve default fontpack for cache identity: %w", err)
	}
	canonical := make([]string, 0, len(names))
	for _, name := range names {
		resolved, _, err := source.Catalog.ResolvePack(name)
		if err != nil {
			return "", fmt.Errorf("resolve fontpack for cache identity: %w", err)
		}
		canonical = append(canonical, resolved)
	}
	slices.Sort(canonical)
	canonical = slices.Compact(canonical)
	if !slices.Contains(canonical, resolvedDefault) {
		return "", fmt.Errorf("default fontpack %q is missing from effective packs", resolvedDefault)
	}
	identity := struct {
		Version       string
		Catalog       *fontpacks.Catalog
		ContentDigest string
		Coverage      string
		Default       string
		Packs         []string
		Picker        bool
	}{
		Version: fontpackCacheVersion, Catalog: source.Catalog,
		ContentDigest: source.ContentDigest, Coverage: coverage,
		Default: resolvedDefault, Packs: canonical, Picker: pickerEnabled,
	}
	data, err := json.Marshal(identity)
	if err != nil {
		return "", fmt.Errorf("encode fontpack cache identity: %w", err)
	}
	return buildcache.ContentHash(string(data)), nil
}
