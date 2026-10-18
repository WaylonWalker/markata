package plugins

import (
	"testing"
	"time"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
)

func TestThemeCalendarPlugin_Name(t *testing.T) {
	p := NewThemeCalendarPlugin()
	if p.Name() != "theme_calendar" {
		t.Errorf("expected name 'theme_calendar', got %q", p.Name())
	}
}

func TestThemeCalendarPlugin_Interfaces(_ *testing.T) {
	p := NewThemeCalendarPlugin()

	// Verify plugin implements required interfaces
	var _ lifecycle.Plugin = p
	var _ lifecycle.ConfigurePlugin = p
	var _ lifecycle.PriorityPlugin = p
}

func TestThemeCalendarPlugin_Priority(t *testing.T) {
	p := NewThemeCalendarPlugin()

	// Should have high priority (negative) in Configure stage
	if got := p.Priority(lifecycle.StageConfigure); got >= 0 {
		t.Errorf("expected negative priority for Configure stage, got %d", got)
	}

	// Should have default priority in other stages
	if got := p.Priority(lifecycle.StageRender); got != 0 {
		t.Errorf("expected 0 priority for Render stage, got %d", got)
	}
}

func TestThemeCalendarPlugin_ParseMMDD(t *testing.T) {
	p := NewThemeCalendarPlugin()

	tests := []struct {
		name      string
		input     string
		wantMonth int
		wantDay   int
		wantErr   bool
	}{
		{
			name:      "valid date",
			input:     "12-25",
			wantMonth: 12,
			wantDay:   25,
		},
		{
			name:      "january first",
			input:     "01-01",
			wantMonth: 1,
			wantDay:   1,
		},
		{
			name:      "february end",
			input:     "02-28",
			wantMonth: 2,
			wantDay:   28,
		},
		{
			name:    "invalid format no dash",
			input:   "1225",
			wantErr: true,
		},
		{
			name:    "invalid month 13",
			input:   "13-01",
			wantErr: true,
		},
		{
			name:    "invalid month 0",
			input:   "00-01",
			wantErr: true,
		},
		{
			name:    "invalid day 0",
			input:   "01-00",
			wantErr: true,
		},
		{
			name:    "invalid day 32",
			input:   "01-32",
			wantErr: true,
		},
		{
			name:    "non-numeric month",
			input:   "ab-01",
			wantErr: true,
		},
		{
			name:    "non-numeric day",
			input:   "01-cd",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			month, day, err := p.parseMMDD(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseMMDD() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if month != tt.wantMonth {
					t.Errorf("parseMMDD() month = %d, want %d", month, tt.wantMonth)
				}
				if day != tt.wantDay {
					t.Errorf("parseMMDD() day = %d, want %d", day, tt.wantDay)
				}
			}
		})
	}
}

func TestThemeCalendarPlugin_IsDateInRange(t *testing.T) {
	p := NewThemeCalendarPlugin()

	tests := []struct {
		name      string
		month     int
		day       int
		startDate string
		endDate   string
		want      bool
	}{
		// Simple ranges (no year boundary)
		{
			name:      "in range middle",
			month:     3,
			day:       15,
			startDate: "03-01",
			endDate:   "03-31",
			want:      true,
		},
		{
			name:      "in range start boundary",
			month:     3,
			day:       1,
			startDate: "03-01",
			endDate:   "03-31",
			want:      true,
		},
		{
			name:      "in range end boundary",
			month:     3,
			day:       31,
			startDate: "03-01",
			endDate:   "03-31",
			want:      true,
		},
		{
			name:      "before range",
			month:     2,
			day:       28,
			startDate: "03-01",
			endDate:   "03-31",
			want:      false,
		},
		{
			name:      "after range",
			month:     4,
			day:       1,
			startDate: "03-01",
			endDate:   "03-31",
			want:      false,
		},
		// Year boundary ranges (e.g., winter season Dec-Feb)
		{
			name:      "year boundary in december",
			month:     12,
			day:       15,
			startDate: "12-01",
			endDate:   "02-28",
			want:      true,
		},
		{
			name:      "year boundary in january",
			month:     1,
			day:       15,
			startDate: "12-01",
			endDate:   "02-28",
			want:      true,
		},
		{
			name:      "year boundary in february",
			month:     2,
			day:       15,
			startDate: "12-01",
			endDate:   "02-28",
			want:      true,
		},
		{
			name:      "year boundary at start",
			month:     12,
			day:       1,
			startDate: "12-01",
			endDate:   "02-28",
			want:      true,
		},
		{
			name:      "year boundary at end",
			month:     2,
			day:       28,
			startDate: "12-01",
			endDate:   "02-28",
			want:      true,
		},
		{
			name:      "year boundary outside in march",
			month:     3,
			day:       1,
			startDate: "12-01",
			endDate:   "02-28",
			want:      false,
		},
		{
			name:      "year boundary outside in november",
			month:     11,
			day:       30,
			startDate: "12-01",
			endDate:   "02-28",
			want:      false,
		},
		// Christmas season
		{
			name:      "christmas eve",
			month:     12,
			day:       24,
			startDate: "12-15",
			endDate:   "12-26",
			want:      true,
		},
		{
			name:      "after christmas season",
			month:     12,
			day:       27,
			startDate: "12-15",
			endDate:   "12-26",
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := p.isDateInRange(tt.month, tt.day, tt.startDate, tt.endDate)
			if got != tt.want {
				t.Errorf("isDateInRange(%d, %d, %q, %q) = %v, want %v",
					tt.month, tt.day, tt.startDate, tt.endDate, got, tt.want)
			}
		})
	}
}

func TestThemeCalendarPlugin_Configure(t *testing.T) {
	tests := []struct {
		name         string
		date         time.Time
		extra        map[string]interface{}
		wantPalette  string
		wantNoChange bool
	}{
		{
			name:         "no calendar config",
			date:         time.Date(2024, 12, 25, 0, 0, 0, 0, time.UTC),
			extra:        nil,
			wantNoChange: true,
		},
		{
			name: "disabled calendar",
			date: time.Date(2024, 12, 25, 0, 0, 0, 0, time.UTC),
			extra: map[string]interface{}{
				"theme_calendar": map[string]interface{}{
					"enabled": false,
					"rules": []interface{}{
						map[string]interface{}{
							"name":       "Christmas",
							"start_date": "12-15",
							"end_date":   "12-26",
							"palette":    "christmas",
						},
					},
				},
			},
			wantNoChange: true,
		},
		{
			name: "matching rule applies palette",
			date: time.Date(2024, 12, 25, 0, 0, 0, 0, time.UTC),
			extra: map[string]interface{}{
				"theme_calendar": map[string]interface{}{
					"enabled": true,
					"rules": []interface{}{
						map[string]interface{}{
							"name":       "Christmas",
							"start_date": "12-15",
							"end_date":   "12-26",
							"palette":    "christmas",
						},
					},
				},
			},
			wantPalette: "christmas",
		},
		{
			name: "no matching rule",
			date: time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC),
			extra: map[string]interface{}{
				"theme_calendar": map[string]interface{}{
					"enabled": true,
					"rules": []interface{}{
						map[string]interface{}{
							"name":       "Christmas",
							"start_date": "12-15",
							"end_date":   "12-26",
							"palette":    "christmas",
						},
					},
				},
			},
			wantNoChange: true,
		},
		{
			name: "year boundary winter theme in january",
			date: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			extra: map[string]interface{}{
				"theme_calendar": map[string]interface{}{
					"enabled": true,
					"rules": []interface{}{
						map[string]interface{}{
							"name":       "Winter",
							"start_date": "12-01",
							"end_date":   "02-28",
							"palette":    "winter-frost",
						},
					},
				},
			},
			wantPalette: "winter-frost",
		},
		{
			name: "first matching rule wins",
			date: time.Date(2024, 12, 20, 0, 0, 0, 0, time.UTC),
			extra: map[string]interface{}{
				"theme_calendar": map[string]interface{}{
					"enabled": true,
					"rules": []interface{}{
						map[string]interface{}{
							"name":       "Christmas",
							"start_date": "12-15",
							"end_date":   "12-26",
							"palette":    "christmas",
						},
						map[string]interface{}{
							"name":       "Winter",
							"start_date": "12-01",
							"end_date":   "02-28",
							"palette":    "winter-frost",
						},
					},
				},
			},
			wantPalette: "christmas", // First matching rule wins
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewThemeCalendarPlugin()
			p.nowFunc = func() time.Time { return tt.date }

			config := &lifecycle.Config{
				Extra: tt.extra,
			}

			m := lifecycle.NewManager()
			m.SetConfig(config)

			err := p.Configure(m)
			if err != nil {
				t.Fatalf("Configure() error = %v", err)
			}

			cfg := m.Config()

			if tt.wantNoChange {
				if cfg.Extra != nil {
					if theme, ok := cfg.Extra["theme"].(map[string]interface{}); ok {
						if palette, ok := theme["palette"].(string); ok && palette != "" {
							// Only fail if palette was set by the plugin (not pre-existing)
							if tt.extra == nil || tt.extra["theme"] == nil {
								t.Errorf("expected no theme change, but palette set to %q", palette)
							}
						}
					}
				}
				return
			}

			if cfg.Extra == nil {
				t.Fatal("expected Extra to be set")
			}

			theme, ok := cfg.Extra["theme"].(map[string]interface{})
			if !ok {
				t.Fatal("expected theme map in Extra")
			}

			palette, ok := theme["palette"].(string)
			if !ok {
				t.Fatal("expected palette in theme")
			}

			if palette != tt.wantPalette {
				t.Errorf("palette = %q, want %q", palette, tt.wantPalette)
			}
		})
	}
}

func TestThemeCalendarPlugin_ApplyRule_Variables(t *testing.T) {
	p := NewThemeCalendarPlugin()
	p.nowFunc = func() time.Time { return time.Date(2024, 12, 25, 0, 0, 0, 0, time.UTC) }

	config := &lifecycle.Config{
		Extra: map[string]interface{}{
			"theme_calendar": map[string]interface{}{
				"enabled": true,
				"rules": []interface{}{
					map[string]interface{}{
						"name":       "Christmas",
						"start_date": "12-15",
						"end_date":   "12-26",
						"palette":    "christmas",
						"variables": map[string]interface{}{
							"--accent":     "#ff0000",
							"--background": "#00ff00",
						},
					},
				},
			},
			"theme": map[string]interface{}{
				"variables": map[string]interface{}{
					"--text": "#ffffff",
				},
			},
		},
	}

	m := lifecycle.NewManager()
	m.SetConfig(config)

	err := p.Configure(m)
	if err != nil {
		t.Fatalf("Configure() error = %v", err)
	}

	cfg := m.Config()
	theme, ok := cfg.Extra["theme"].(map[string]interface{})
	if !ok {
		t.Fatal("theme not found in config")
	}
	vars, ok := theme["variables"].(map[string]interface{})
	if !ok {
		t.Fatal("variables not found in theme")
	}

	// Check merged variables
	if vars["--accent"] != "#ff0000" {
		t.Errorf("--accent = %v, want #ff0000", vars["--accent"])
	}
	if vars["--background"] != "#00ff00" {
		t.Errorf("--background = %v, want #00ff00", vars["--background"])
	}
	if vars["--text"] != "#ffffff" {
		t.Errorf("--text = %v, want #ffffff (should preserve existing)", vars["--text"])
	}
}

func TestThemeCalendarPlugin_ApplyRule_Font(t *testing.T) {
	p := NewThemeCalendarPlugin()
	p.nowFunc = func() time.Time { return time.Date(2024, 12, 25, 0, 0, 0, 0, time.UTC) }

	config := &lifecycle.Config{
		Extra: map[string]interface{}{
			"theme_calendar": map[string]interface{}{
				"enabled": true,
				"rules": []interface{}{
					map[string]interface{}{
						"name":       "Christmas",
						"start_date": "12-15",
						"end_date":   "12-26",
						"font": map[string]interface{}{
							"family":         "Mountains of Christmas",
							"heading_family": "Snowburst One",
							"google_fonts":   []interface{}{"Mountains of Christmas", "Snowburst One"},
						},
					},
				},
			},
		},
	}

	m := lifecycle.NewManager()
	m.SetConfig(config)

	err := p.Configure(m)
	if err != nil {
		t.Fatalf("Configure() error = %v", err)
	}

	cfg := m.Config()
	theme, ok := cfg.Extra["theme"].(map[string]interface{})
	if !ok {
		t.Fatal("theme not found in config")
	}
	font, ok := theme["font"].(map[string]interface{})
	if !ok {
		t.Fatal("font not found in theme")
	}

	if font["family"] != "Mountains of Christmas" {
		t.Errorf("font family = %v, want 'Mountains of Christmas'", font["family"])
	}
	if font["heading_family"] != "Snowburst One" {
		t.Errorf("heading_family = %v, want 'Snowburst One'", font["heading_family"])
	}
}

func TestThemeCalendarPlugin_ApplyRule_LightDarkPalettes(t *testing.T) {
	p := NewThemeCalendarPlugin()
	p.nowFunc = func() time.Time { return time.Date(2024, 12, 25, 0, 0, 0, 0, time.UTC) }

	config := &lifecycle.Config{
		Extra: map[string]interface{}{
			"theme_calendar": map[string]interface{}{
				"enabled": true,
				"rules": []interface{}{
					map[string]interface{}{
						"name":          "Christmas",
						"start_date":    "12-15",
						"end_date":      "12-26",
						"palette_light": "christmas-light",
						"palette_dark":  "christmas-dark",
					},
				},
			},
		},
	}

	m := lifecycle.NewManager()
	m.SetConfig(config)

	err := p.Configure(m)
	if err != nil {
		t.Fatalf("Configure() error = %v", err)
	}

	cfg := m.Config()
	theme, ok := cfg.Extra["theme"].(map[string]interface{})
	if !ok {
		t.Fatal("theme not found in config")
	}

	if theme["palette_light"] != "christmas-light" {
		t.Errorf("palette_light = %v, want 'christmas-light'", theme["palette_light"])
	}
	if theme["palette_dark"] != "christmas-dark" {
		t.Errorf("palette_dark = %v, want 'christmas-dark'", theme["palette_dark"])
	}
}

func TestThemeCalendarPlugin_InvalidDateFormats(t *testing.T) {
	p := NewThemeCalendarPlugin()
	p.nowFunc = func() time.Time { return time.Date(2024, 12, 25, 0, 0, 0, 0, time.UTC) }

	// Rules with invalid date formats should not match
	config := &lifecycle.Config{
		Extra: map[string]interface{}{
			"theme_calendar": map[string]interface{}{
				"enabled": true,
				"rules": []interface{}{
					map[string]interface{}{
						"name":       "Bad Start Date",
						"start_date": "invalid",
						"end_date":   "12-26",
						"palette":    "should-not-apply",
					},
					map[string]interface{}{
						"name":       "Bad End Date",
						"start_date": "12-15",
						"end_date":   "invalid",
						"palette":    "also-should-not-apply",
					},
				},
			},
		},
	}

	m := lifecycle.NewManager()
	m.SetConfig(config)

	// Should not error, just skip invalid rules
	err := p.Configure(m)
	if err != nil {
		t.Fatalf("Configure() should not error on invalid rules: %v", err)
	}

	cfg := m.Config()
	if cfg.Extra["theme"] != nil {
		theme, ok := cfg.Extra["theme"].(map[string]interface{})
		if ok {
			if palette, ok := theme["palette"].(string); ok && palette != "" {
				t.Errorf("expected no palette to be set with invalid rules, got %q", palette)
			}
		}
	}
}
