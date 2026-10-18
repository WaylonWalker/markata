package plugins

import (
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/assets"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/runtimeenv"
)

func TestCDNAssetsPlugin_OfflineMissingAssetFails(t *testing.T) {
	t.Setenv(runtimeenv.EnvOffline, "true")
	t.Setenv(runtimeenv.EnvBundledAssetsCacheDir, t.TempDir())

	config := &lifecycle.Config{
		OutputDir: t.TempDir(),
		Extra: map[string]interface{}{
			"assets": models.AssetsConfig{Mode: "cdn"},
			"cdn_assets_extra": []assets.Asset{
				{
					Name:      "missing-extra",
					URL:       "https://example.invalid/missing.js",
					LocalPath: "missing/missing.js",
					Type:      "js",
				},
			},
		},
	}

	manager := lifecycle.NewManager()
	manager.SetConfig(config)

	err := NewCDNAssetsPlugin().Configure(manager)
	if err == nil {
		t.Fatal("expected offline CDN assets configure to fail when asset is missing")
	}
	if !strings.Contains(err.Error(), "missing-extra") {
		t.Fatalf("expected missing asset name in error, got %v", err)
	}
}
