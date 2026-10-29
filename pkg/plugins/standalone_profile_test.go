package plugins

import (
	"strings"
	"testing"
)

func TestStandaloneProfile_WebRetainsExistingAssetBehavior(t *testing.T) {
	result, err := RenderStandalone("# Hello", StandaloneOptions{Profile: RenderProfileWeb, AssetBase: "/custom"})
	if err != nil {
		t.Fatal(err)
	}
	if result.RenderProfile != RenderProfileWeb || !strings.Contains(result.HeadHTML, "/custom/css/main.css") {
		t.Fatalf("web profile changed existing asset linking: %s", result.HeadHTML)
	}
}

func TestStandaloneProfile_OfflineInteractiveInlinesManagedAssets(t *testing.T) {
	result, err := RenderStandalone("# Local", StandaloneOptions{Profile: RenderProfileOfflineInteractive})
	if err != nil {
		t.Fatal(err)
	}
	if result.RenderProfile != RenderProfileOfflineInteractive {
		t.Fatalf("profile not reported: %q", result.RenderProfile)
	}
	if !strings.Contains(result.HeadHTML, "<style>") || strings.Contains(result.HeadHTML, "<link ") {
		t.Fatalf("bundled styles must be inline: %s", result.HeadHTML)
	}
	for _, asset := range result.Assets {
		if asset.URL != "" {
			t.Fatalf("offline profile exposed remote or linked asset %s: %q", asset.ID, asset.URL)
		}
	}
}

func TestStandaloneProfile_RejectsExternalMedia(t *testing.T) {
	for _, source := range []string{
		"![remote](https://example.com/secret.png)",
		"![remote](//example.com/remote.png)",
	} {
		if _, err := RenderStandalone(source, StandaloneOptions{Profile: RenderProfileOfflineInteractive}); err == nil {
			t.Fatalf("remote resource must be rejected: %s", source)
		}
	}
}

func TestStandaloneProfile_RejectsUnsafeAndUnsupportedStatic(t *testing.T) {
	for _, options := range []StandaloneOptions{
		{Profile: RenderProfileOfflineInteractive, Unsafe: true},
		{Profile: RenderProfileOfflineStatic},
		{Profile: "typo"},
	} {
		if _, err := RenderStandalone("# test", options); err == nil {
			t.Fatalf("expected explicit profile error: %+v", options)
		}
	}
}

func TestStandaloneProfile_AuditsGeneratedAssetURLs(t *testing.T) {
	if err := checkStandaloneOfflineResources(&StandaloneResult{
		Assets: []RenderAsset{{ID: "bad", Kind: "script", URL: "https://cdn.example.com/code.js"}},
	}); err == nil {
		t.Fatal("remote generated asset URL accepted")
	}
	if err := checkStandaloneOfflineResources(&StandaloneResult{
		Assets: []RenderAsset{{ID: "bad", Kind: "style", Content: "@import url('https://fonts.example/style.css');"}},
	}); err == nil {
		t.Fatal("remote CSS import accepted")
	}
	if err := checkStandaloneOfflineResources(&StandaloneResult{
		HTML: `<a href="https://example.com/">safe navigation</a>`,
	}); err != nil {
		t.Fatalf("ordinary user-initiated links must be allowed: %v", err)
	}
}
