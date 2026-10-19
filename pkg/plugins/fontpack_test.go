package plugins

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/fontpacks"
	"github.com/WaylonWalker/markata/pkg/models"

	htmlparser "golang.org/x/net/html"
)

func TestConfiguredFontpackName_UsesCanonicalTheme(t *testing.T) {
	configured := &models.Config{}
	configured.Theme.Fontpack = "brush-poster"
	if got := configuredFontpackName(map[string]any{"models_config": configured, "fontpack": "system"}); got != "brush" {
		t.Fatalf("fontpack = %q", got)
	}
}

func TestMarkHTMLFontpackReplacesExactlyOneAttribute(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"exact", "<html>", `<html data-fontpack="field-notebook">`},
		{"attributes", `<html lang="en">`, `<html lang="en" data-fontpack="field-notebook">`},
		{"uppercase", `<HTML class="site">`, `<HTML class="site" data-fontpack="field-notebook">`},
		{"replace", `<html data-fontpack="old" lang="en">`, `<html data-fontpack="field-notebook" lang="en">`},
		{"already correct", `<html  data-fontpack="field-notebook" class="two  words">`, `<html  data-fontpack="field-notebook" class="two  words">`},
		{"quoted spaces", `<html class="two  words" data-fontpack="old">`, `<html class="two  words" data-fontpack="field-notebook">`},
		{"single quotes", `<html data-fontpack='old' title='a > b'>`, `<html data-fontpack="field-notebook" title='a > b'>`},
		{"mixed case", `<HtMl DATA-FONTPACK = 'old'>`, `<HtMl data-fontpack="field-notebook">`},
		{"boolean", `<html data-fontpack lang="en">`, `<html data-fontpack="field-notebook" lang="en">`},
		{"duplicate", `<html data-fontpack="old" DATA-FONTPACK="other">`, `<html data-fontpack="field-notebook" >`},
		{"self closing", `<html/>`, `<html data-fontpack="field-notebook"/>`},
		{"unquoted slash value", `<html lang=en data-base=/>`, `<html lang=en data-base=/ data-fontpack="field-notebook">`},
		{"unquoted trailing slash", `<html class=site/>`, `<html class=site/ data-fontpack="field-notebook">`},
		{"stray slash", `<html class/site data-fontpack="old">`, `<html class/site data-fontpack="field-notebook">`},
		{"leading equals", `<html =data-fontpack="unrelated">`, `<html =data-fontpack="unrelated" data-fontpack="field-notebook">`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := markHTMLFontpack(tt.in, "field-notebook")
			if got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
			tokenizer := htmlparser.NewTokenizer(strings.NewReader(got))
			tokenizer.Next()
			count := 0
			for _, attribute := range tokenizer.Token().Attr {
				if attribute.Key == "data-fontpack" {
					count++
					if attribute.Val != "field-notebook" {
						t.Fatalf("incorrect parsed fontpack: %q", attribute.Val)
					}
				}
			}
			if count != 1 {
				t.Fatalf("duplicate data-fontpack attribute: %q", got)
			}
		})
	}
}

func TestMarkHTMLFontpackPreservesFragmentsAndMalformedTags(t *testing.T) {
	for _, content := range []string{"", "<p>fragment</p>", `<html title="unfinished`, `<html data-fontpack="unfinished>`} {
		if got := markHTMLFontpack(content, "brush"); got != content {
			t.Fatalf("malformed or fragment HTML changed: %q", got)
		}
	}
}

func TestMarkHTMLFontpackIgnoresCommentedTags(t *testing.T) {
	content := `<!-- <html data-fontpack="wrong"> --><html>`
	want := `<!-- <html data-fontpack="wrong"> --><html data-fontpack="brush">`
	if got := markHTMLFontpack(content, "brush"); got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestMarkHTMLFontpackEscapesPackName(t *testing.T) {
	name := `custom"&pack`
	got := markHTMLFontpack("<html>", name)
	want := `<html data-fontpack="custom&#34;&amp;pack">`
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
	if markedAgain := markHTMLFontpack(got, name); markedAgain != got {
		t.Fatalf("annotation is not idempotent: %q", markedAgain)
	}
}

func TestMarkPostFontpackKeepsHashedStylesheet(t *testing.T) {
	html := `<html><head><link rel="stylesheet" href="/css/fonts.12345678.css"></head><body></body></html>`
	got := markPostFontpack(html, "brush")
	if strings.Count(got, "fonts.") != 1 || strings.Contains(got, `href="/css/fonts.css"`) {
		t.Fatalf("duplicate font stylesheet: %s", got)
	}

	if !strings.Contains(got, `data-fontpack="brush"`) {
		t.Fatalf("missing resolved fontpack: %s", got)
	}
}

func BenchmarkMarkPostFontpack(b *testing.B) {
	for _, size := range []struct {
		name string
		body int
	}{{"small", 1024}, {"large", 1024 * 1024}} {
		b.Run(size.name, func(b *testing.B) {
			content := `<html lang="en" data-fontpack="brush"><head><link rel="stylesheet" href="/css/fonts.12345678.css"></head><body>` +
				strings.Repeat("content ", size.body/8) + `</body></html>`
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if got := markPostFontpack(content, "brush"); got != content {
					b.Fatal("unchanged fontpack annotation changed the page")
				}
			}
		})
	}
}

func TestFontpackCacheKeyTracksVisibleGlyphCoverage(t *testing.T) {
	names := []string{"system", "handwritten"}
	catalog, err := fontpacks.BuiltinSource()
	if err != nil {
		t.Fatal(err)
	}
	first := fontpackTestCacheKey(t, `<p>tone</p><script>Ж</script>`, "system", names, false, catalog)
	sameCoverage := fontpackTestCacheKey(t, `<strong>note</strong><style>Ж</style>`, "system", names, false, catalog)
	if first != sameCoverage {
		t.Fatal("fontpack cache key changed even though the visible rune set was unchanged")
	}
	newCoverage := fontpackTestCacheKey(t, `<p>note é</p>`, "system", names, false, catalog)
	if first == newCoverage {
		t.Fatal("fontpack cache key did not change when a new visible rune was introduced")
	}
}

func TestFontpackOutputCachedRequiresMatchingKeyAndStylesheet(t *testing.T) {
	output := t.TempDir()
	if fontpackOutputCached(output, "key") {
		t.Fatal("missing marker was treated as a cache hit")
	}

	if err := os.WriteFile(filepath.Join(output, fontpackCacheFile), []byte("key"), 0o600); err != nil {
		t.Fatal(err)
	}

	if fontpackOutputCached(output, "key") {
		t.Fatal("missing stylesheet was treated as a cache hit")
	}
	if err := os.Mkdir(filepath.Join(output, "css"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(output, "css", "fonts.css"), []byte(""), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(output, "assets", "fonts"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(output, "assets", "fonts", ".markata-fonts.json"), []byte(`{"files":[]}`), 0o600); err != nil {
		t.Fatal(err)
	}
	if !fontpackOutputCached(output, "key") {
		t.Fatal("matching marker and stylesheet were not treated as a cache hit")
	}
	if fontpackOutputCached(output, "different") {
		t.Fatal("mismatched marker was treated as a cache hit")
	}
}

func TestFontpackPreloadURLs(t *testing.T) {
	assets := []fontpacks.Asset{
		{Source: "body", Tier: "prose-core", URL: "/assets/fonts/body-core.woff2"},
		{Source: "body", Tier: "latin-ext", URL: "/assets/fonts/body-ext.woff2"},
		{Source: "display", Tier: "full", URL: "/assets/fonts/display-full.woff2"},
		{Source: "code", Tier: "full", URL: "/assets/fonts/code-full.woff2"},
	}

	pack := fontpacks.FontPack{Roles: map[string]fontpacks.Role{
		"body": {Source: "body", Tier: "prose-core"}, "display": {Source: "display", Tier: "full"},
		"heading": {Source: "body"}, "code": {Source: "code"},
	}}
	got := fontpackPreloadURLs(pack, assets)
	want := []string{"/assets/fonts/body-core.woff2", "/assets/fonts/display-full.woff2"}
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Fatalf("preloads = %v, want %v", got, want)
	}

	assets = append(assets, fontpacks.Asset{Source: "heading", Tier: "display-core", URL: "/assets/fonts/heading-core.woff2"})
	pack.Roles["heading"] = fontpacks.Role{Source: "heading", Tier: "display-core"}
	if got := fontpackPreloadURLs(pack, assets); strings.Join(got, ",") != strings.Join(append(want, "/assets/fonts/heading-core.woff2"), ",") {
		t.Fatalf("distinct heading face was suppressed by display role: %v", got)
	}
	pack.Roles["heading"] = fontpacks.Role{Source: "body", Tier: "prose-core"}

	if got := fontpackPreloadURLs(fontpacks.FontPack{Roles: map[string]fontpacks.Role{"body": {Stack: "sans"}}}, assets); len(got) != 0 {
		t.Fatalf("system pack preloads = %v", got)
	}
	assets = append(assets, fontpacks.Asset{Source: "body", Tier: "full", URL: "/assets/fonts/body-full.woff2"})
	if got := fontpackPreloadURLs(pack, assets); got[0] != "/assets/fonts/body-core.woff2" {
		t.Fatalf("unrelated full tier replaced the requested core preload: %v", got)
	}
	fullPack := fontpacks.FontPack{Roles: map[string]fontpacks.Role{"body": {Source: "body", Tier: "full"}}}
	if got := fontpackPreloadURLs(fullPack, assets); got[0] != "/assets/fonts/body-full.woff2" {
		t.Fatalf("explicit full role did not preload its requested tier: %v", got)
	}
	assets = []fontpacks.Asset{{Source: "body", Tier: "full", URL: "/assets/fonts/body-full.woff2"}}
	missingTierPack := fontpacks.FontPack{Roles: map[string]fontpacks.Role{"body": {Source: "body", Tier: "prose-core"}}}
	if got := fontpackPreloadURLs(missingTierPack, assets); len(got) != 1 || got[0] != "/assets/fonts/body-full.woff2" {
		t.Fatalf("full fallback was not preloaded when no subset exists: %v", got)
	}
}

func TestFontpackHTMLCacheHashTracksPreloadSelection(t *testing.T) {
	first := fontpackPreloadCache{Hash: "1234abcd", URLs: map[string][]string{"handwritten": {"/fonts/body-core.woff2"}}}
	second := fontpackPreloadCache{Hash: "1234abcd", URLs: map[string][]string{"handwritten": {"/fonts/body-full.woff2"}}}
	if fontpackHTMLCacheHash(first) == fontpackHTMLCacheHash(second) {
		t.Fatal("preload URL change did not invalidate cached HTML")
	}
	same := fontpackPreloadCache{Hash: first.Hash, URLs: map[string][]string{"handwritten": {"/fonts/body-core.woff2"}}}
	if fontpackHTMLCacheHash(first) != fontpackHTMLCacheHash(same) {
		t.Fatal("same font CSS and preload URLs changed the HTML cache identity")
	}
}

func TestFontpackHTMLCacheHashIsIndependentOfMapOrder(t *testing.T) {
	first := fontpackPreloadCache{
		Hash: "1234abcd",
		URLs: map[string][]string{
			"handwritten": {"/fonts/body-core.woff2"},
			"brutalist":   {"/fonts/heading-core.woff2"},
		},
	}
	second := fontpackPreloadCache{Hash: "1234abcd", URLs: make(map[string][]string)}
	second.URLs["brutalist"] = []string{"/fonts/heading-core.woff2"}
	second.URLs["handwritten"] = []string{"/fonts/body-core.woff2"}
	if got, want := fontpackHTMLCacheHash(second), fontpackHTMLCacheHash(first); got != want {
		t.Fatalf("identical preload data in different map orders changed HTML cache hash: %s != %s", got, want)
	}
	second.Hash = "5678efab"
	if fontpackHTMLCacheHash(second) == fontpackHTMLCacheHash(first) {
		t.Fatal("font CSS hash change did not invalidate HTML cache")
	}
}

func TestValidFontpackPreloadCache(t *testing.T) {
	output := t.TempDir()
	catalog := &fontpacks.Catalog{FontPacks: map[string]fontpacks.FontPack{"system": {}}}
	cache := fontpackPreloadCache{Hash: "1234abcd", URLs: map[string][]string{"system": {}}, EditorialURLs: map[string][]string{"system": {}}}
	if !validFontpackPreloadCache(output, []string{"system"}, catalog, cache) {
		t.Fatal("empty system pack should have a valid cached preload list")
	}
	cache.URLs["system"] = []string{"/assets/fonts/missing.woff2"}
	if validFontpackPreloadCache(output, []string{"system"}, catalog, cache) {
		t.Fatal("missing font file accepted")
	}
	cache.URLs["system"] = []string{"/assets/fonts/../bad.woff2"}
	if validFontpackPreloadCache(output, []string{"system"}, catalog, cache) {
		t.Fatal("unsafe font path accepted")
	}
	cache.URLs["system"] = nil
	cache.Hash = "invalid!"
	if validFontpackPreloadCache(output, []string{"system"}, catalog, cache) {
		t.Fatal("invalid hash accepted")
	}
}

func TestFontpackEditorialPreloadsCodeWithoutDuplicatingSharedSources(t *testing.T) {
	assets := []fontpacks.Asset{
		{Source: "body", Tier: "prose-core", URL: "/assets/fonts/body.woff2"},
		{Source: "code", Tier: "code-core", URL: "/assets/fonts/code.woff2"},
		{Source: "code", Tier: "full", URL: "/assets/fonts/code-full.woff2"},
	}
	pack := fontpacks.FontPack{Roles: map[string]fontpacks.Role{
		"body": {Source: "body", Tier: "prose-core"},
		"code": {Source: "code", Tier: "code-core"},
	}}
	if got := fontpackEditorialPreloadURLs(pack, assets); strings.Join(got, ",") != "/assets/fonts/body.woff2,/assets/fonts/code.woff2" {
		t.Fatalf("editorial preloads = %v", got)
	}
	if got := fontpackPreloadURLs(pack, assets); len(got) != 1 {
		t.Fatalf("ordinary pages unexpectedly preload code: %v", got)
	}
	pack.Roles["code"] = fontpacks.Role{Source: "code", Tier: "full"}
	if got := fontpackEditorialPreloadURLs(pack, assets); got[1] != "/assets/fonts/code-full.woff2" {
		t.Fatalf("requested full code tier not selected: %v", got)
	}
	pack.Roles["code"] = fontpacks.Role{Source: "body", Tier: "code-core"}
	if got := fontpackEditorialPreloadURLs(pack, assets); len(got) != 1 {
		t.Fatalf("shared source was duplicated: %v", got)
	}
	pack.Roles["code"] = fontpacks.Role{Stack: "mono"}
	if got := fontpackEditorialPreloadURLs(pack, assets); len(got) != 1 {
		t.Fatalf("system code stack added a preload: %v", got)
	}
}

func TestFontpackEditorialMetadataInvalidatesHTMLAndRequiresIntactOutput(t *testing.T) {
	cache := fontpackPreloadCache{Hash: "1234abcd", URLs: map[string][]string{"system": {}}, EditorialURLs: map[string][]string{"system": {}}}
	before := fontpackHTMLCacheHash(cache)
	cache.EditorialURLs["system"] = []string{"/assets/fonts/code.woff2"}
	if before == fontpackHTMLCacheHash(cache) {
		t.Fatal("changed editorial preload failed to invalidate HTML")
	}
	output := t.TempDir()
	catalog := &fontpacks.Catalog{FontPacks: map[string]fontpacks.FontPack{"system": {}}}
	if validFontpackPreloadCache(output, []string{"system"}, catalog, cache) {
		t.Fatal("missing editorial font output accepted")
	}
	fonts := filepath.Join(output, "assets", "fonts")
	if err := os.MkdirAll(fonts, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(fonts, "code.woff2"), []byte("font"), 0o600); err != nil {
		t.Fatal(err)
	}
	if !validFontpackPreloadCache(output, []string{"system"}, catalog, cache) {
		t.Fatal("intact editorial font output rejected")
	}
	delete(cache.EditorialURLs, "system")
	if validFontpackPreloadCache(output, []string{"system"}, catalog, cache) {
		t.Fatal("old or incomplete editorial metadata accepted")
	}
}
