// Package terminaltheme resolves semantic colors shared by markata-go's CLI
// logs, prompts, and terminal user interfaces.
package terminaltheme

import "github.com/WaylonWalker/markata/pkg/palettes"

// Theme contains semantic colors used across terminal surfaces.
type Theme struct {
	Accent              string
	AccentHover         string
	Link                string
	TextPrimary         string
	TextSecondary       string
	TextMuted           string
	Success             string
	Warning             string
	Error               string
	Info                string
	Border              string
	BorderFocus         string
	Background          string
	Surface             string
	Elevated            string
	ButtonPrimaryBg     string
	ButtonPrimaryText   string
	ButtonSecondaryBg   string
	ButtonSecondaryText string
}

// Default returns the common fallback palette for terminal interfaces.
func Default() Theme {
	return DefaultForVariant(palettes.VariantDark)
}

// DefaultForVariant returns shared fallback colors for a light or dark palette.
func DefaultForVariant(variant palettes.Variant) Theme {
	if variant == palettes.VariantLight {
		return Theme{
			Accent:              "#3b82f6",
			AccentHover:         "#2563eb",
			Link:                "#1f2937",
			TextPrimary:         "#111827",
			TextSecondary:       "#1f2937",
			TextMuted:           "#4b5563",
			Success:             "#047857",
			Warning:             "#b45309",
			Error:               "#dc2626",
			Info:                "#2563eb",
			Border:              "#d1d5db",
			BorderFocus:         "#3b82f6",
			Background:          "#f9fafb",
			Surface:             "#f3f4f6",
			Elevated:            "#e5e7eb",
			ButtonPrimaryBg:     "#2563eb",
			ButtonPrimaryText:   "#f9fafb",
			ButtonSecondaryBg:   "#e5e7eb",
			ButtonSecondaryText: "#1f2937",
		}
	}

	return Theme{
		Accent:              "#60a5fa",
		AccentHover:         "#93c5fd",
		Link:                "#60a5fa",
		TextPrimary:         "#f3f4f6",
		TextSecondary:       "#d1d5db",
		TextMuted:           "#9ca3af",
		Success:             "#4ade80",
		Warning:             "#facc15",
		Error:               "#f87171",
		Info:                "#60a5fa",
		Border:              "#374151",
		BorderFocus:         "#60a5fa",
		Background:          "#111827",
		Surface:             "#1f2937",
		Elevated:            "#374151",
		ButtonPrimaryBg:     "#2563eb",
		ButtonPrimaryText:   "#f9fafb",
		ButtonSecondaryBg:   "#374151",
		ButtonSecondaryText: "#f3f4f6",
	}
}

// FromPalette resolves the shared semantic colors from a site palette.
// Missing semantic roles use the common fallback values.
func FromPalette(palette *palettes.Palette) Theme {
	base := Default()
	if palette == nil {
		return base
	}
	base = DefaultForVariant(palette.Variant)
	resolve := func(name, fallback string, aliases ...string) string {
		if value := palette.Resolve(name); value != "" {
			return value
		}
		for _, alias := range aliases {
			if value := palette.Resolve(alias); value != "" {
				return value
			}
		}
		return fallback
	}
	return Theme{
		Accent:              resolve("accent", base.Accent, "primary", "info"),
		AccentHover:         resolve("accent-hover", base.AccentHover, "link-hover"),
		Link:                resolve("link", base.Link, "accent", "primary", "info"),
		TextPrimary:         resolve("text-primary", base.TextPrimary, "text"),
		TextSecondary:       resolve("text-secondary", base.TextSecondary, "text-primary", "text"),
		TextMuted:           resolve("text-muted", base.TextMuted, "text-secondary"),
		Success:             resolve("success", base.Success, "info"),
		Warning:             resolve("warning", base.Warning),
		Error:               resolve("error", base.Error),
		Info:                resolve("info", base.Info, "accent"),
		Border:              resolve("border", base.Border),
		BorderFocus:         resolve("border-focus", base.BorderFocus, "accent"),
		Background:          resolve("bg-primary", base.Background, "bg-secondary"),
		Surface:             resolve("bg-surface", base.Surface, "bg-secondary"),
		Elevated:            resolve("bg-elevated", base.Elevated, "bg-surface"),
		ButtonPrimaryBg:     resolve("button-primary-bg", base.ButtonPrimaryBg, "accent"),
		ButtonPrimaryText:   resolve("button-primary-text", base.ButtonPrimaryText, "text-primary"),
		ButtonSecondaryBg:   resolve("button-secondary-bg", base.ButtonSecondaryBg, "bg-elevated"),
		ButtonSecondaryText: resolve("button-secondary-text", base.ButtonSecondaryText, "text-primary"),
	}
}
