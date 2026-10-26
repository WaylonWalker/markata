package plugins

import (
	"testing"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

func TestTypedConfigProjectionConsumers(t *testing.T) {
	highlightEnabled := false
	calendarEnabled := true
	resourceHintsEnabled := false
	autoDetect := false

	highlight := models.MarkdownConfig{
		Highlight: models.HighlightConfig{
			Enabled:     &highlightEnabled,
			Theme:       "monokai",
			LineNumbers: true,
		},
	}
	calendar := models.ThemeCalendarConfig{
		Enabled: &calendarEnabled,
		Rules: []models.ThemeCalendarRule{{
			Name:      "Winter",
			StartDate: "12-01",
			EndDate:   "02-28",
			Palette:   "winter-frost",
		}},
	}
	resourceHints := models.ResourceHintsConfig{
		Enabled:    &resourceHintsEnabled,
		AutoDetect: &autoDetect,
	}

	extra := map[string]interface{}{
		"markdown":       highlight,
		"theme_calendar": calendar,
		"resource_hints": resourceHints,
	}

	markdownPlugin := NewRenderMarkdownPlugin()
	gotTheme, gotLineNumbers := markdownPlugin.resolveHighlightConfig(extra)
	if gotTheme != disabledHighlightTheme || gotLineNumbers {
		t.Errorf("resolveHighlightConfig() = (%q, %v), want disabled typed config", gotTheme, gotLineNumbers)
	}

	chromaEnabled := true
	chromaExtra := map[string]interface{}{
		"markdown": models.MarkdownConfig{Highlight: models.HighlightConfig{
			Enabled: &chromaEnabled,
			Theme:   "monokai",
		}},
	}
	if gotTheme, explicit := NewChromaCSSPlugin().getExplicitHighlightTheme(chromaExtra); !explicit || gotTheme != "monokai" {
		t.Errorf("getExplicitHighlightTheme() = (%q, %v), want typed explicit theme", gotTheme, explicit)
	}

	if gotTheme, gotLineNumbers, gotEnabled := resolveEmbedHighlightConfig(extra); gotTheme != "monokai" || !gotLineNumbers || gotEnabled {
		t.Errorf("resolveEmbedHighlightConfig() = (%q, %v, %v), want typed values", gotTheme, gotLineNumbers, gotEnabled)
	}

	calendarPlugin := NewThemeCalendarPlugin()
	gotCalendar := calendarPlugin.getCalendarConfig(&lifecycle.Config{Extra: extra})
	if !gotCalendar.IsEnabled() || len(gotCalendar.Rules) != 1 || gotCalendar.Rules[0].Palette != "winter-frost" {
		t.Errorf("getCalendarConfig() = %+v, want typed calendar", gotCalendar)
	}

	gotHints := getResourceHintsConfig(&lifecycle.Config{Extra: extra})
	if gotHints.IsEnabled() || gotHints.IsAutoDetectEnabled() {
		t.Errorf("getResourceHintsConfig() = %+v, want explicit false values", gotHints)
	}
}
