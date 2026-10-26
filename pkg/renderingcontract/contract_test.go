package renderingcontract

import "testing"

func TestContract_UsesZeroDownloadSystemFontByDefault(t *testing.T) {
	c, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if c.Defaults.Fontpack != "system" {
		t.Fatalf("default fontpack = %q, want system", c.Defaults.Fontpack)
	}
	pack, ok := c.Fontpacks["system"]
	if !ok || pack["body"] != "system-ui" || pack["heading"] != "system-ui" || pack["mono"] != "ui-monospace" {
		t.Fatalf("system typography roles not configured: %#v", pack)
	}
	if _, ok := c.Fontpacks["brush"]; !ok {
		t.Fatal("opt-in brush fontpack must remain available")
	}
}

func TestContract_PaletteFamiliesArePaired(t *testing.T) {
	c, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	pairs := map[string]map[string]bool{}
	for _, palette := range c.Palettes {
		if palette.Variant != "light" && palette.Variant != "dark" {
			t.Fatalf("%s has invalid variant %q", palette.ID, palette.Variant)
		}
		if pairs[palette.Family] == nil {
			pairs[palette.Family] = map[string]bool{}
		}
		if pairs[palette.Family][palette.Variant] {
			t.Fatalf("duplicate %s variant %s", palette.Family, palette.Variant)
		}
		pairs[palette.Family][palette.Variant] = true
	}
	if len(pairs) == 0 {
		t.Fatal("expected palette catalog")
	}
	classified := map[string]bool{}
	for _, family := range c.IncompletePaletteFamilies {
		classified[family] = true
	}
	for family, variants := range pairs {
		if len(variants) < 2 && !classified[family] {
			t.Errorf("incomplete palette family %q is not explicitly classified", family)
		}
	}
}

func TestContract_FinalRenderSweepCoversEveryCanonicalPalette(t *testing.T) {
	c, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(c.Palettes) != 120 {
		t.Fatalf("canonical palette count = %d, want 120", len(c.Palettes))
	}
	audited := 0
	for _, palette := range c.Palettes {
		audited++
		for _, result := range FinalRenderContrast(palette) {
			if !result.Passed {
				t.Fatalf("final render projection failed for %s: %s on %s (%.2f:1, need %.1f)", palette.ID, result.Foreground, result.Background, result.Ratio, result.Required)
			}
		}
	}
	if audited != len(c.Palettes) {
		t.Fatalf("audited %d canonical palettes, want %d", audited, len(c.Palettes))
	}
}

func TestNormalizeMix(t *testing.T) {
	for input, want := range map[any]float64{"35%": .35, 0.3519: .352, -1: 0, 2: 1} {
		if got := NormalizeMix(input, .5); got != want {
			t.Errorf("NormalizeMix(%v) = %v, want %v", input, got, want)
		}
	}
}

func TestNormalizeTheme_CanonicalWinsLegacy(t *testing.T) {
	got, warnings := NormalizeTheme(map[string]any{"texture_strength": "10%", "theme": map[string]any{"texture": map[string]any{"color_mix": 0.9}}})
	texture, ok := got["texture"].(map[string]any)
	if !ok {
		t.Fatalf("texture has unexpected type: %#v", got["texture"])
	}
	if texture["color_mix"] != 0.9 || len(warnings) != 1 {
		t.Fatalf("got %#v warnings %#v", got, warnings)
	}
}

func TestNormalizeTheme_EquivalentLegacyMixDoesNotWarn(t *testing.T) {
	got, warnings := NormalizeTheme(map[string]any{
		"texture_strength": "35%",
		"theme":            map[string]any{"texture": map[string]any{"color_mix": .35}},
	})
	texture, ok := got["texture"].(map[string]any)
	if !ok {
		t.Fatalf("texture has unexpected type: %#v", got["texture"])
	}
	if texture["color_mix"] != .35 || len(warnings) != 0 {
		t.Fatalf("equivalent migration = %#v warnings %#v", got, warnings)
	}
}

func TestNormalizeTheme_CanonicalScopeConflictsLegacy(t *testing.T) {
	_, warnings := NormalizeTheme(map[string]any{
		"texture_scope": "all",
		"theme":         map[string]any{"texture": map[string]any{"scope": "quiet"}},
	})
	if len(warnings) != 1 {
		t.Fatalf("warnings = %#v", warnings)
	}
}

func TestNormalizeTheme_HeadingTextureScope(t *testing.T) {
	got, warnings := NormalizeTheme(map[string]any{"texture": "splatter", "texture_scope": "headings"})
	texture, ok := got["texture"].(map[string]any)
	if !ok {
		t.Fatalf("texture has unexpected type: %#v", got["texture"])
	}
	heading, ok := got["heading_texture"].(map[string]any)
	if !ok {
		t.Fatalf("heading texture has unexpected type: %#v", got["heading_texture"])
	}
	if texture["scope"] != "quiet" || texture["kind"] != "none" || heading["kind"] != "splatter" {
		t.Fatalf("migration = %#v", got)
	}
	if len(warnings) != 1 {
		t.Fatalf("warnings = %#v", warnings)
	}
}

func TestNormalizeTheme_CanonicalTextureKindWinsLegacyHeadingScope(t *testing.T) {
	got, warnings := NormalizeTheme(map[string]any{
		"texture_scope": "headings",
		"theme":         map[string]any{"texture": map[string]any{"kind": "screenprint"}},
	})
	texture, ok := got["texture"].(map[string]any)
	if !ok {
		t.Fatalf("texture has unexpected type: %#v", got["texture"])
	}
	if texture["kind"] != "screenprint" || len(warnings) == 0 {
		t.Fatalf("canonical kind was not preserved: %#v warnings %#v", got, warnings)
	}
}

func TestNormalizeTheme_CanonicalTextureScopeWinsLegacyHeadingScope(t *testing.T) {
	got, warnings := NormalizeTheme(map[string]any{
		"texture_scope": "headings",
		"theme":         map[string]any{"texture": map[string]any{"scope": "all"}},
	})
	texture, ok := got["texture"].(map[string]any)
	if !ok {
		t.Fatalf("texture has unexpected type: %#v", got["texture"])
	}
	if texture["scope"] != "all" || len(warnings) == 0 {
		t.Fatalf("canonical scope was not preserved: %#v warnings %#v", got, warnings)
	}
}

func TestNormalizeMixAndScale_Int64(t *testing.T) {
	if got := NormalizeMix(int64(1), 0); got != 1 {
		t.Fatalf("mix = %v", got)
	}
	if got := NormalizeScale(int64(2), 0); got != 2 {
		t.Fatalf("scale = %v", got)
	}
}
