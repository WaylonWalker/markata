package plugins

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/WaylonWalker/markata/pkg/palettes"
)

func TestEasterDate_KnownYears(t *testing.T) {
	tests := map[int]string{2024: "2024-03-31", 2025: "2025-04-20", 2026: "2026-04-05", 2027: "2027-03-28", 2030: "2030-04-21"}
	for year, want := range tests {
		if got := easterDate(year).Format("2006-01-02"); got != want {
			t.Errorf("easterDate(%d) = %s, want %s", year, got, want)
		}
	}
}

func seasonalTestManifest() []PaletteManifestEntry {
	pair := func(dark, light string) []PaletteManifestEntry {
		return []PaletteManifestEntry{
			{Name: dark, Variant: themeModeDark, Counterpart: light},
			{Name: light, Variant: themeModeLight, Counterpart: dark},
		}
	}
	var manifest []PaletteManifestEntry
	manifest = append(manifest, pair("winter-frost", "winter-frost-light")...)
	manifest = append(manifest, pair("pollen8-dark", "pollen8")...)
	manifest = append(manifest, pair("summer-beach-dark", "summer-beach")...)
	manifest = append(manifest, pair("autumn", "autumn-light")...)
	manifest = append(manifest, pair("christmas", "christmas-light")...)
	manifest = append(manifest, pair("black-gold", "white-gold")...)
	return manifest
}

func TestBuildSeasonalSchedule_ResolvesModesAndDropsMissingPalettes(t *testing.T) {
	schedule := buildSeasonalSchedule(seasonalTestManifest(), time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC))
	if schedule == nil {
		t.Fatal("expected a schedule")
	}
	if len(schedule.Seasons) != 4 {
		t.Fatalf("seasons = %d, want 4", len(schedule.Seasons))
	}
	spring := schedule.Seasons[0]
	if spring.Light != "pollen8" || spring.Dark != "pollen8-dark" {
		t.Errorf("spring = %+v, want pollen8 / pollen8-dark", spring)
	}
	labels := map[string]seasonalEntry{}
	for _, h := range schedule.Holidays {
		labels[h.Label] = h
	}
	if _, ok := labels["Halloween"]; ok {
		t.Error("holidays whose palettes are not shipped must be dropped")
	}
	newYear, ok := labels["New Year"]
	if !ok || newYear.Light != "white-gold" || newYear.Dark != "black-gold" {
		t.Errorf("new year = %+v, want white-gold / black-gold", newYear)
	}
	if xmas := labels["Christmas"]; len(xmas.Dates) != 1 || xmas.Dates[0] != "12-25" {
		t.Errorf("christmas dates = %v, want [12-25]", xmas.Dates)
	}
}

func TestBuildSeasonalSchedule_NoSeasonsReturnsNil(t *testing.T) {
	manifest := []PaletteManifestEntry{{Name: "christmas", Variant: themeModeDark}}
	if schedule := buildSeasonalSchedule(manifest, time.Now()); schedule != nil {
		t.Errorf("expected nil schedule without season palettes, got %+v", schedule)
	}
}

func TestSeasonalScheduleJSON_BuiltinPalettesCoverAllHolidays(t *testing.T) {
	manifest := testFullPaletteManifest(t)
	raw := seasonalScheduleJSON(manifest, time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC))
	var schedule seasonalSchedule
	if err := json.Unmarshal([]byte(raw), &schedule); err != nil {
		t.Fatalf("invalid seasonal JSON %q: %v", raw, err)
	}
	if len(schedule.Seasons) != len(seasonalSeasons) {
		t.Errorf("seasons = %d, want %d", len(schedule.Seasons), len(seasonalSeasons))
	}
	if len(schedule.Holidays) != len(seasonalHolidays) {
		var got []string
		for _, h := range schedule.Holidays {
			got = append(got, h.Label)
		}
		t.Errorf("holidays = %s, want all %d", strings.Join(got, ", "), len(seasonalHolidays))
	}
	for _, entry := range append(schedule.Seasons, schedule.Holidays...) {
		if entry.Light == "" || entry.Dark == "" || entry.Light == entry.Dark {
			t.Errorf("%s: light %q / dark %q must be distinct palettes", entry.Label, entry.Light, entry.Dark)
		}
	}
}

func testFullPaletteManifest(t *testing.T) []PaletteManifestEntry {
	t.Helper()
	infos, err := palettes.NewLoader().Discover()
	if err != nil {
		t.Fatalf("discover palettes: %v", err)
	}
	return linkCounterparts((&PaletteCSSPlugin{}).generatePaletteManifest(infos, "ayu-light", "ayu-dark"))
}
