package fontpacks

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"
)

func testCatalog(t *testing.T) *Catalog {
	t.Helper()
	return &Catalog{
		SystemStacks: map[string]SystemStack{"sans": {CSS: "system-ui, sans-serif"}, "mono": {CSS: "monospace"}},
		SubsetProfiles: map[string]SubsetProfile{
			"prose-core": {Unicode: []string{"U+0020-007E"}},
			"latin-ext":  {Unicode: []string{"U+0100-024F"}},
		},
		FontSources: map[string]FontSource{"demo": {Provider: "test", Family: "Demo"}},
		FontPacks: map[string]FontPack{
			"system":  {Performance: Performance{Class: "zero-download"}, Roles: map[string]Role{"body": {Stack: "sans"}}},
			"bundled": {Performance: Performance{Class: "bundled"}, Roles: map[string]Role{"body": {Source: "demo", Tier: "prose-core"}}},
		},
		Aliases: map[string]string{"default": "system"},
	}
}

func TestResolvePackAliasAndSystemHasNoAssets(t *testing.T) {
	c := testCatalog(t)
	name, pack, err := c.ResolvePack("default")
	if err != nil || name != "system" || pack.Performance.Class != "zero-download" {
		t.Fatalf("resolve alias = %q, %#v, %v", name, pack, err)
	}
	r, err := c.Resolve("default", t.TempDir(), t.TempDir(), "<p>Hello</p>")
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Assets) != 0 || strings.Contains(r.CSS, "@font-face") {
		t.Fatalf("system resolution emitted assets or font-face: %#v\n%s", r.Assets, r.CSS)
	}
}

func TestBuiltinBrushUsesCanonicalDocumentRoles(t *testing.T) {
	source, err := BuiltinSource()
	if err != nil {
		t.Fatal(err)
	}
	name, pack, err := source.Catalog.ResolvePack("brush-poster")
	if err != nil {
		t.Fatal(err)
	}
	if name != "brush" {
		t.Fatalf("resolved name = %q, want brush", name)
	}
	for role, want := range map[string]string{
		"body":    "space-grotesk",
		"heading": "knewave",
		"code":    "dm-mono",
	} {
		if got := pack.Roles[role].Source; got != want {
			t.Errorf("brush %s source = %q, want %q", role, got, want)
		}
	}
}

func TestBuiltinBrushAssetsHaveCanonicalDigests(t *testing.T) {
	source, err := BuiltinSource()
	if err != nil {
		t.Fatal(err)
	}
	resolved, err := source.Catalog.ResolveFS("brush", source.FS, source.Root, "<article><h1>Specimen</h1><p>Text</p><code>x</code></article>")
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]string{
		"space-grotesk": "5010c6e79c7e65a058dfd39a71d8d60f7b56070a3c32cc911d9dddbaec7f0851",
		"knewave":       "56fc586096f8e158fc214e3de31ee47c15ff80fc0abd5e0646ddff47a6af27ba",
		"dm-mono":       "90c42888e73fa88ebf3d0e27e0461535acd661e3892bbcb6eb21047f33b67c32",
	}
	if len(resolved.Assets) != len(want) {
		t.Fatalf("brush assets = %d, want %d", len(resolved.Assets), len(want))
	}
	for _, asset := range resolved.Assets {
		if got := want[asset.Source]; got == "" || asset.SHA256 != got {
			t.Errorf("asset %q digest = %q, want %q", asset.Source, asset.SHA256, got)
		}
	}
}

func TestRequiredTiersSelectExtendedAndFallback(t *testing.T) {
	c := testCatalog(t)
	_, pack, err := c.ResolvePack("bundled")
	if err != nil {
		t.Fatal(err)
	}
	if got := c.RequiredTiers(pack, "<p>Hello</p>")["demo"]; !got["prose-core"] || got["latin-ext"] {
		t.Fatalf("English tiers = %#v", got)
	}
	got := c.RequiredTiers(pack, "<p>Ā</p>")["demo"]
	if !got["latin-ext"] {
		t.Fatalf("extended tiers = %#v", got)
	}
	got = c.RequiredTiers(pack, "<p>Ж</p>")["demo"]
	if !got["prose-core"] || got["full"] || len(got) != 1 {
		t.Fatalf("unsupported rune should retain core tier and use CSS fallback: %#v", got)
	}
}

func TestVisibleTextSkipsCodeAndDecodesEntities(t *testing.T) {
	got := VisibleText(`<main>Hello &amp; <strong>world</strong><script>Ж</script><style>Ж</style></main>`)
	if !strings.Contains(got, "Hello &") || !strings.Contains(got, "world") || strings.Contains(got, "Ж") {
		t.Fatalf("VisibleText() = %q, want decoded visible text without script/style content", got)
	}
}

func TestRequiredTiersUseSystemFallbackWhenSourceLacksOptionalTier(t *testing.T) {
	c := testCatalog(t)
	_, pack, err := c.ResolvePack("bundled")
	if err != nil {
		t.Fatal(err)
	}
	manifests := map[string]Manifest{"demo": {Tiers: map[string]Tier{
		"prose-core": {Profile: "prose-core"}, "full": {Profile: "full"},
	}}}
	got := c.RequiredTiersForManifest(pack, "<p>Ā</p>", manifests)["demo"]
	if !got["prose-core"] || got["full"] || len(got) != 1 {
		t.Fatalf("source without latin-ext selected %#v, want core with CSS fallback", got)
	}
}

func TestResolveDeduplicatesSourceTierAndCopiesAsset(t *testing.T) {
	c := testCatalog(t)
	c.FontPacks["bundled"] = FontPack{Performance: Performance{Class: "bundled"}, Roles: map[string]Role{
		"body": {Source: "demo", Tier: "prose-core"}, "heading": {Source: "demo", Tier: "prose-core"},
	}}
	root := t.TempDir()
	family := filepath.Join(root, "demo")
	if err := os.MkdirAll(family, 0o755); err != nil {
		t.Fatal(err)
	}
	manifest := "schema: markata.font/v1\nid: demo\nfamily: Demo\ntiers:\n  prose-core:\n    file: demo.woff2\n    profile: prose-core\n"
	if err := os.WriteFile(filepath.Join(family, "manifest.yaml"), []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(family, "demo.woff2"), []byte("woff2"), 0o644); err != nil {
		t.Fatal(err)
	}
	r, err := c.Resolve("bundled", root, t.TempDir(), "<p>Hello</p>")
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Assets) != 1 {
		t.Fatalf("assets = %d, want 1", len(r.Assets))
	}
}

func TestBuiltinFieldNotebookResolvesStableBaseTiers(t *testing.T) {
	source, err := BuiltinSource()
	if err != nil {
		t.Fatal(err)
	}
	r, err := source.Catalog.ResolveFS("field-notebook", source.FS, source.Root, "<article><h1>Hello</h1><p>Ordinary prose.</p></article>")
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Assets) != 6 {
		t.Fatalf("assets = %d, want six deduplicated family/tier pairs", len(r.Assets))
	}
	if strings.Contains(r.CSS, "full.woff2") {
		t.Fatalf("ordinary English selected full tiers:\n%s", r.CSS)
	}
}

func TestBuiltinDisplayRoleUsesProseCoreWhenDisplayCoreIsAbsent(t *testing.T) {
	source, err := BuiltinSource()
	if err != nil {
		t.Fatal(err)
	}
	resolved, err := source.Catalog.ResolveFS("brutalist", source.FS, source.Root, "<h1>Ordinary heading</h1>")
	if err != nil {
		t.Fatal(err)
	}
	seen := make(map[string]bool)
	for _, asset := range resolved.Assets {
		seen[asset.Source+":"+asset.Tier] = true
	}
	if !seen["space-grotesk:prose-core"] || seen["space-grotesk:full"] {
		t.Fatalf("missing display-core did not resolve to compact prose-core: %v", seen)
	}
}

func TestResolveManyScopesFullFaceToItsExplicitPack(t *testing.T) {
	c := testCatalog(t)
	c.FontPacks["explicit"] = FontPack{Name: "Explicit", Performance: Performance{Class: "bundled"}, Roles: map[string]Role{
		"body": {Source: "demo", Tier: "full", Fallback: "sans"},
	}}
	assets := fstest.MapFS{
		"demo/manifest.yaml": {Data: []byte("tiers:\n  prose-core: {file: core.woff2, profile: prose-core}\n  full: {file: full.woff2, profile: full}\n")},
		"demo/core.woff2":    {Data: []byte("wOF2-core")},
		"demo/full.woff2":    {Data: []byte("wOF2-full")},
	}
	resolved, err := c.ResolveManyFSWithCoverage([]string{"bundled", "explicit"}, assets, ".", coverageFromString("<p>Ordinary text</p>"), ResolveOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(resolved.CSS, `font-family: "Demo [Markata explicit full]";`) {
		t.Fatalf("full face lacks its pack-specific family:\n%s", resolved.CSS)
	}
	if !strings.Contains(resolved.CSS, `--font-body: "Demo [Markata explicit full]", system-ui, sans-serif;`) {
		t.Fatalf("explicit pack role does not select its full face alias:\n%s", resolved.CSS)
	}
	if strings.Contains(resolved.CSS, `font-family: "Demo";\n  src: url('/assets/fonts/full.woff2')`) {
		t.Fatalf("unrestricted full face competes with the shared core family:\n%s", resolved.CSS)
	}
}

func TestBuiltinPaintedSignUsesExpressiveAndReadableRoles(t *testing.T) {
	source, err := BuiltinSource()
	if err != nil {
		t.Fatal(err)
	}
	pack := source.Catalog.FontPacks["painted-sign"]
	if pack.Roles["display"].Source != "finger-paint" || pack.Roles["heading"].Source != "rock-salt" {
		t.Fatalf("expressive roles = %#v", pack.Roles)
	}
	if pack.Roles["body"].Source != "source-sans-3" || pack.Roles["code"].Source != "dm-mono" {
		t.Fatalf("readability roles = %#v", pack.Roles)
	}
	resolved, err := source.Catalog.ResolveFS("painted-sign", source.FS, source.Root, "<article><h1>Painted sign</h1><p>Readable prose.</p></article>")
	if err != nil {
		t.Fatal(err)
	}
	if len(resolved.Assets) != 5 || !strings.Contains(resolved.CSS, "Finger Paint") || !strings.Contains(resolved.CSS, "Rock Salt") {
		t.Fatalf("painted-sign assets/CSS = %d\n%s", len(resolved.Assets), resolved.CSS)
	}
	if !strings.Contains(resolved.CSS, "h1 {\n  font-family: var(--font-display)") {
		t.Fatalf("display role is not consumed by h1:\n%s", resolved.CSS)
	}
	if !strings.Contains(resolved.CSS, "h2, h3, h4, h5, h6 {\n  font-family: var(--font-heading)") {
		t.Fatalf("heading role is not consumed by h2-h6:\n%s", resolved.CSS)
	}
}

func TestBuiltinHandwrittenMapsLegacyMonoTokenToItsCodeRole(t *testing.T) {
	source, err := BuiltinSource()
	if err != nil {
		t.Fatal(err)
	}
	resolved, err := source.Catalog.ResolveFS("handwritten", source.FS, source.Root, "<kbd>Key</kbd>")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(resolved.CSS, `--font-mono: var(--font-code);`) ||
		!strings.Contains(resolved.CSS, `--font-code: "DM Mono"`) {
		t.Fatalf("legacy mono token does not use the pack's code role:\n%s", resolved.CSS)
	}
}

func TestBuiltinEditorialQuoteDoesNotInheritBodyOpticalSize(t *testing.T) {
	source, err := BuiltinSource()
	if err != nil {
		t.Fatal(err)
	}
	resolved, err := source.Catalog.ResolveManyFS([]string{"editorial"}, source.FS, source.Root, "<blockquote>Quote</blockquote>")
	if err != nil {
		t.Fatal(err)
	}
	start := strings.Index(resolved.CSS, `[data-fontpack="editorial"] blockquote`)
	if start < 0 {
		t.Fatalf("missing editorial quote rule:\n%s", resolved.CSS)
	}
	end := strings.Index(resolved.CSS[start:], "}\n")
	if end < 0 {
		t.Fatalf("unterminated editorial quote rule")
	}
	rule := resolved.CSS[start : start+end]
	if !strings.Contains(rule, "font-family: var(--font-quote)") || !strings.Contains(rule, "font-size: var(--font-quote-size)") {
		t.Fatalf("editorial quote role lost its configured properties: %s", rule)
	}
	if strings.Contains(rule, "font-variation-settings:") || strings.Contains(rule, "opsz") {
		t.Fatalf("editorial quote inherited body optical size: %s", rule)
	}
}

func TestRoleCSSDoesNotResetUnspecifiedHeadingSize(t *testing.T) {
	c := testCatalog(t)
	c.FontPacks["bundled"] = FontPack{Performance: Performance{Class: "bundled"}, Roles: map[string]Role{
		"heading": {Source: "demo", Tier: "prose-core", Weight: 400},
	}}
	root := t.TempDir()
	family := filepath.Join(root, "demo")
	if err := os.MkdirAll(family, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(family, "manifest.yaml"), []byte("tiers:\n  prose-core: {file: demo.woff2, profile: prose-core}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(family, "demo.woff2"), []byte("wOF2"), 0o644); err != nil {
		t.Fatal(err)
	}
	r, err := c.Resolve("bundled", root, t.TempDir(), `<style>h1 { font-size: 3rem; } h2 { font-size: 2rem; }</style><h1>Title</h1><h2>Section</h2>`)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(r.CSS, "font-size:") {
		t.Fatalf("unspecified heading size reset existing hierarchy:\n%s", r.CSS)
	}
}

func TestDisplayFallsBackToHeadingWhenAbsent(t *testing.T) {
	c := testCatalog(t)
	c.FontPacks["bundled"] = FontPack{Performance: Performance{Class: "bundled"}, Roles: map[string]Role{
		"heading": {Source: "demo", Tier: "prose-core", Weight: 400},
	}}
	root := t.TempDir()
	family := filepath.Join(root, "demo")
	if err := os.MkdirAll(family, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(family, "manifest.yaml"), []byte("tiers:\n  prose-core: {file: demo.woff2, profile: prose-core}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(family, "demo.woff2"), []byte("wOF2"), 0o644); err != nil {
		t.Fatal(err)
	}
	r, err := c.Resolve("bundled", root, t.TempDir(), "<h1>Title</h1>")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(r.CSS, "h1 {\n  font-family: var(--font-heading)") || strings.Contains(r.CSS, "var(--font-display)") {
		t.Fatalf("h1 did not fall back to heading:\n%s", r.CSS)
	}
}

func TestBuiltinPaintedSignUsesFallbackForDisplayFontsWithoutExtendedTier(t *testing.T) {
	source, err := BuiltinSource()
	if err != nil {
		t.Fatal(err)
	}
	resolved, err := source.Catalog.ResolveFS("painted-sign", source.FS, source.Root, "<article><p>Ǎ</p></article>")
	if err != nil {
		t.Fatal(err)
	}
	seen := make(map[string]bool)
	for _, asset := range resolved.Assets {
		seen[asset.Source+":"+asset.Tier] = true
	}
	for _, want := range []string{"finger-paint:display-core", "rock-salt:display-core", "source-sans-3:prose-basic", "source-sans-3:prose-extra", "dm-mono:latin-ext"} {
		if !seen[want] {
			t.Errorf("extended painted-sign coverage missing %s (assets=%v)", want, seen)
		}
	}
	if seen["finger-paint:full"] || seen["rock-salt:full"] {
		t.Fatalf("display fonts without extended tiers should use CSS fallback: %v", seen)
	}
}

func TestRoleTypographyPropertiesReachGeneratedCSS(t *testing.T) {
	c := testCatalog(t)
	c.FontPacks["bundled"].Roles["body"] = Role{Source: "demo", Tier: "prose-core", Weight: 450, Style: "italic", Size: "1.2rem", OpticalSize: 18}
	root := t.TempDir()
	family := filepath.Join(root, "demo")
	if err := os.MkdirAll(family, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(family, "manifest.yaml"), []byte("faces:\n  italic:\n    style: italic\n    variable: true\n    weight: [300, 800]\n    axes:\n      wght: [300, 800]\n      opsz: [9, 144]\ntiers:\n  prose-core: {file: demo.woff2, profile: prose-core}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(family, "demo.woff2"), []byte("wOF2"), 0o644); err != nil {
		t.Fatal(err)
	}
	r, err := c.Resolve("bundled", root, t.TempDir(), "<p>Hello</p>")
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"--font-body-weight: 450", "--font-body-style: italic", "--font-body-size: 1.2rem", "--font-body-variation: \"opsz\" 18", "font-weight: var(--font-body-weight"} {
		if !strings.Contains(r.CSS, want) {
			t.Errorf("CSS missing %q:\n%s", want, r.CSS)
		}
	}
}

func TestResolveManyGeneratesPerPackRoleSelectors(t *testing.T) {
	c := testCatalog(t)
	root := t.TempDir()
	family := filepath.Join(root, "demo")
	if err := os.MkdirAll(family, 0o755); err != nil {
		t.Fatal(err)
	}
	manifest := "schema: markata.font/v1\nid: demo\nfamily: Demo\ntiers:\n  prose-core:\n    file: demo.woff2\n    profile: prose-core\n"
	if err := os.WriteFile(filepath.Join(family, "manifest.yaml"), []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(family, "demo.woff2"), []byte("wOF2"), 0o644); err != nil {
		t.Fatal(err)
	}
	c.FontPacks["alternate"] = FontPack{Performance: Performance{Class: "bundled"}, Roles: map[string]Role{"body": {Source: "demo", Tier: "prose-core"}}}
	resolved, err := c.ResolveMany([]string{"system", "alternate"}, root, "<p>Hello</p>")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(resolved.CSS, `[data-fontpack="system"]`) || !strings.Contains(resolved.CSS, `[data-fontpack="alternate"]`) {
		t.Fatalf("missing per-pack selectors:\n%s", resolved.CSS)
	}
	if len(resolved.Assets) != 1 {
		t.Fatalf("assets = %d, want one deduplicated asset", len(resolved.Assets))
	}
}

func TestResolveManyScopesOptionalRolePropertiesPerPack(t *testing.T) {
	c := testCatalog(t)
	c.FontPacks["pack-a"] = FontPack{Performance: Performance{Class: "bundled"}, Roles: map[string]Role{
		"quote": {Source: "demo", Tier: "prose-core", Size: "1.25rem", Style: "italic"},
	}}
	c.FontPacks["pack-b"] = FontPack{Performance: Performance{Class: "bundled"}, Roles: map[string]Role{
		"quote": {Source: "demo", Tier: "prose-core"},
	}}
	root := t.TempDir()
	family := filepath.Join(root, "demo")
	if err := os.MkdirAll(family, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(family, "manifest.yaml"), []byte("tiers:\n  prose-core: {file: demo.woff2, profile: prose-core}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(family, "demo.woff2"), []byte("wOF2"), 0o644); err != nil {
		t.Fatal(err)
	}
	resolved, err := c.ResolveManyFS([]string{"pack-a", "pack-b"}, os.DirFS(root), ".", "<blockquote>Quote</blockquote>")
	if err != nil {
		t.Fatal(err)
	}
	for name, wantSize := range map[string]bool{"pack-a": true, "pack-b": false} {
		start := strings.Index(resolved.CSS, `[data-fontpack="`+name+`"] blockquote`)
		if start < 0 {
			t.Fatalf("missing scoped quote rule for %s:\n%s", name, resolved.CSS)
		}
		end := strings.Index(resolved.CSS[start:], "}\n")
		if end < 0 {
			t.Fatalf("unterminated scoped quote rule for %s", name)
		}
		rule := resolved.CSS[start : start+end]
		if strings.Contains(rule, "font-size:") != wantSize || strings.Contains(rule, "font-style:") != wantSize {
			t.Errorf("%s optional properties leaked or disappeared: %s", name, rule)
		}
	}
}

func TestResolveManyDoesNotInheritBodyOptionalPropertiesForPresentRole(t *testing.T) {
	c := testCatalog(t)
	c.FontPacks["specialized"] = FontPack{Performance: Performance{Class: "bundled"}, Roles: map[string]Role{
		"body":  {Source: "demo", Tier: "prose-core", OpticalSize: 18},
		"quote": {Source: "demo", Tier: "prose-core", Weight: 600, Size: "1.35rem"},
	}}
	c.FontPacks["absent"] = FontPack{Performance: Performance{Class: "bundled"}, Roles: map[string]Role{
		"body": {Source: "demo", Tier: "prose-core", OpticalSize: 18},
	}}
	root := t.TempDir()
	family := filepath.Join(root, "demo")
	if err := os.MkdirAll(family, 0o755); err != nil {
		t.Fatal(err)
	}
	manifest := "faces:\n  normal:\n    style: normal\n    variable: true\n    weight: [300, 800]\n    axes:\n      wght: [300, 800]\n      opsz: [8, 72]\ntiers:\n  prose-core: {file: demo.woff2, profile: prose-core}\n"
	if err := os.WriteFile(filepath.Join(family, "manifest.yaml"), []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(family, "demo.woff2"), []byte("wOF2"), 0o644); err != nil {
		t.Fatal(err)
	}

	resolved, err := c.ResolveManyFS([]string{"specialized", "absent"}, os.DirFS(root), ".", "<blockquote>Quote</blockquote>")
	if err != nil {
		t.Fatal(err)
	}
	for name, want := range map[string][]string{
		"specialized": {
			"font-family: var(--font-quote)",
			"font-weight: var(--font-quote-weight)",
			"font-size: var(--font-quote-size)",
		},
		"absent": {
			"font-family: var(--font-body)",
			"font-variation-settings: var(--font-body-variation)",
		},
	} {
		start := strings.Index(resolved.CSS, `[data-fontpack="`+name+`"] blockquote`)
		if start < 0 {
			t.Fatalf("missing scoped quote rule for %s:\n%s", name, resolved.CSS)
		}
		end := strings.Index(resolved.CSS[start:], "}\n")
		if end < 0 {
			t.Fatalf("unterminated scoped quote rule for %s", name)
		}
		rule := resolved.CSS[start : start+end]
		for _, property := range want {
			if !strings.Contains(rule, property) {
				t.Errorf("%s quote rule missing %q: %s", name, property, rule)
			}
		}
		if name == "specialized" && strings.Contains(rule, "font-variation-settings:") {
			t.Errorf("%s quote rule inherited body optical size: %s", name, rule)
		}
	}
}

func TestResolveManyPreservesDisplayHeadingMapping(t *testing.T) {
	c := testCatalog(t)
	c.FontPacks["sign"] = FontPack{Performance: Performance{Class: "bundled"}, Roles: map[string]Role{
		"display": {Source: "demo", Tier: "prose-core"},
		"heading": {Source: "demo", Tier: "prose-core"},
	}}
	root := t.TempDir()
	family := filepath.Join(root, "demo")
	if err := os.MkdirAll(family, 0o755); err != nil {
		t.Fatal(err)
	}
	manifest := "tiers:\n  prose-core: {file: demo.woff2, profile: prose-core}\n"
	if err := os.WriteFile(filepath.Join(family, "manifest.yaml"), []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(family, "demo.woff2"), []byte("wOF2"), 0o644); err != nil {
		t.Fatal(err)
	}
	resolved, err := c.ResolveManyFS([]string{"sign"}, os.DirFS(root), ".", "<h1>Title</h1>")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(resolved.CSS, "[data-fontpack=\"sign\"] h1 {\n  font-family: var(--font-display)") || !strings.Contains(resolved.CSS, `[data-fontpack="sign"] h2,`) || !strings.Contains(resolved.CSS, `[data-fontpack="sign"] h3`) {
		t.Fatalf("display/heading mapping regressed:\n%s", resolved.CSS)
	}
}

func TestLoadSourceResolvesAssetsRelativeToNestedCatalog(t *testing.T) {
	root := t.TempDir()
	catalogDir := filepath.Join(root, "config", "fonts")
	assetDir := filepath.Join(catalogDir, "assets", "demo")
	if err := os.MkdirAll(assetDir, 0o755); err != nil {
		t.Fatal(err)
	}
	catalog := `schema: markata.fontpacks/v2
system_stacks: {sans: {css: system-ui}}
subset_profiles: {prose-core: {unicode: [U+0020-007E]}}
font_sources: {demo: {provider: test, family: Demo}}
fontpacks: {custom: {performance: {class: bundled}, roles: {body: {source: demo, tier: prose-core}}}}
catalog: {bundled_asset_root: assets}
`
	path := filepath.Join(catalogDir, "catalog.yaml")
	if err := os.WriteFile(path, []byte(catalog), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(assetDir, "manifest.yaml"), []byte("tiers:\n  prose-core: {file: demo.woff2, profile: prose-core}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(assetDir, "demo.woff2"), []byte("wOF2"), 0o644); err != nil {
		t.Fatal(err)
	}
	source, err := LoadSource(path)
	if err != nil {
		t.Fatal(err)
	}
	resolved, err := source.Catalog.ResolveFS("custom", source.FS, source.Root, "<p>Hello</p>")
	if err != nil || len(resolved.Assets) != 1 {
		t.Fatalf("nested catalog resolution = %#v, %v", resolved, err)
	}
}

func TestAssetSHA256IsFullHash(t *testing.T) {
	path := filepath.Join(t.TempDir(), "font.woff2")
	if err := os.WriteFile(path, []byte("wOF2"), 0o644); err != nil {
		t.Fatal(err)
	}
	hash, _, err := AssetSHA256(path)
	if err != nil || len(hash) != 64 {
		t.Fatalf("hash = %q, err = %v", hash, err)
	}
}

func TestCopyFSRemovesOnlyStaleManagedFonts(t *testing.T) {
	root := t.TempDir()
	family := filepath.Join(root, "demo")
	if err := os.MkdirAll(family, 0o755); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"old.woff2", "new.woff2"} {
		if err := os.WriteFile(filepath.Join(family, name), []byte("wOF2"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	out := filepath.Join(t.TempDir(), "output")
	first := &Resolved{Assets: []Asset{{Source: "demo", File: "old.woff2"}}}
	if err := first.Copy(root, out); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(out, "assets/fonts/user.woff2"), []byte("user"), 0o644); err != nil {
		t.Fatal(err)
	}
	second := &Resolved{Assets: []Asset{{Source: "demo", File: "new.woff2"}}}
	if err := second.Copy(root, out); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(out, "assets/fonts/old.woff2")); !os.IsNotExist(err) {
		t.Fatalf("stale managed font still exists: %v", err)
	}
	if _, err := os.Stat(filepath.Join(out, "assets/fonts/user.woff2")); err != nil {
		t.Fatalf("user font was removed: %v", err)
	}
}
