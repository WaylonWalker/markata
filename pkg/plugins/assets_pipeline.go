// Package plugins provides lifecycle plugins for markata-go.
package plugins

import (
	"github.com/WaylonWalker/markata/pkg/assets"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
)

func assetURLsFromConfig(config *lifecycle.Config) map[string]string {
	if config == nil || config.Extra == nil {
		return nil
	}
	return assetURLsFromExtra(config.Extra)
}

func assetURLsFromExtra(extra map[string]interface{}) map[string]string {
	if extra == nil {
		return nil
	}

	if direct, ok := extra["asset_urls"].(map[string]string); ok {
		return direct
	}

	if anyMap, ok := extra["asset_urls"].(map[string]interface{}); ok {
		assetURLs := make(map[string]string, len(anyMap))
		for key, value := range anyMap {
			if v, ok := value.(string); ok && v != "" {
				assetURLs[key] = v
			}
		}
		return assetURLs
	}

	return nil
}

func resolveAssetURL(assetURLs map[string]string, name, fallback string) string {
	if assetURLs == nil {
		return fallback
	}
	if url := assetURLs[name]; url != "" {
		return url
	}
	return fallback
}

func appendExtraAsset(config *lifecycle.Config, asset assets.Asset) {
	if config.Extra == nil {
		config.Extra = make(map[string]interface{})
	}

	existing, ok := config.Extra["cdn_assets_extra"].([]assets.Asset)
	if !ok {
		existing = nil
	}
	for i := range existing {
		current := existing[i]
		if current.Name == asset.Name {
			return
		}
	}
	config.Extra["cdn_assets_extra"] = append(existing, asset)
}
