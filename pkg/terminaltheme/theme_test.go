package terminaltheme

import (
	"testing"

	"github.com/WaylonWalker/markata/pkg/palettes"
)

func TestFromPaletteUsesSharedFallbacks(t *testing.T) {
	theme := FromPalette(nil)
	if theme.Accent != "#60a5fa" || theme.Success != "#4ade80" || theme.TextMuted != "#9ca3af" {
		t.Fatalf("unexpected default theme: %+v", theme)
	}
}

func TestFromSparseLightPaletteUsesLightFallbacks(t *testing.T) {
	theme := FromPalette(palettes.NewPalette("sparse-light", palettes.VariantLight))
	if theme.TextPrimary != "#111827" || theme.Background != "#f9fafb" {
		t.Fatalf("sparse light palette has non-light fallbacks: %+v", theme)
	}
}

func TestFromPaletteResolvesRolesAndAliases(t *testing.T) {
	palette := palettes.NewPalette("test", palettes.VariantDark)
	palette.Colors["brand"] = "#123456"
	palette.Colors["text"] = "#eeeeee"
	palette.Semantic["accent"] = "brand"
	palette.Semantic["text-primary"] = "text"
	palette.Semantic["success"] = "#00aa00"

	theme := FromPalette(palette)
	if theme.Accent != "#123456" || theme.TextPrimary != "#eeeeee" || theme.Success != "#00aa00" {
		t.Fatalf("palette theme = %+v", theme)
	}
	if theme.Warning != Default().Warning {
		t.Fatalf("missing roles should use common defaults: warning = %q", theme.Warning)
	}
}
