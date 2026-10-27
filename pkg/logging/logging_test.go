package logging

import (
	"bytes"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/palettes"
)

func TestWriterNormalizesBracketedComponent(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	writer := NewWriter(Options{Writer: buf, Format: FormatPlain})

	if _, err := writer.Write([]byte("[mentions] Processing 3 posts\n")); err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "[mentions] Processing 3 posts") {
		t.Fatalf("output = %q", output)
	}
}

func TestWriterNormalizesColonComponent(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	writer := NewWriter(Options{Writer: buf, Format: FormatPlain})

	if _, err := writer.Write([]byte("hashtag_tags: processing 10 posts\n")); err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "[hashtag_tags] processing 10 posts") {
		t.Fatalf("output = %q", output)
	}
}

func TestWriterDoesNotTreatWarningAsComponent(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	writer := NewWriter(Options{Writer: buf, Format: FormatPlain})

	if _, err := writer.Write([]byte("warning: duplicate alias\n")); err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	output := buf.String()
	if strings.Contains(output, "[warning]") {
		t.Fatalf("output = %q", output)
	}
	if !strings.Contains(output, "warning: duplicate alias") {
		t.Fatalf("output = %q", output)
	}
}

func TestAllowColorForceColor(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	t.Setenv("FORCE_COLOR", "")
	t.Setenv("TERM", "xterm-256color")

	if !allowColor(Options{Format: FormatRich, ForceColor: true}) {
		t.Fatal("expected force color to enable color")
	}
}

func TestColorEnabledMatchesCLIForcePrecedence(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	t.Setenv("FORCE_COLOR", "")
	t.Setenv("TERM", "dumb")
	if !ColorEnabled(false, FormatRich, true, false) {
		t.Fatal("explicit force color should override automatic suppression")
	}
}

func TestColorEnabledWhitespaceNoColor(t *testing.T) {
	t.Setenv("NO_COLOR", " ")
	t.Setenv("FORCE_COLOR", "")
	t.Setenv("TERM", "xterm-256color")
	if ColorEnabled(true, FormatAuto, false, false) {
		t.Fatal("any non-empty NO_COLOR value should disable automatic color")
	}
}

func TestRichComponentUsesPhaseColor(t *testing.T) {
	t.Setenv("NO_COLOR", "")
	t.Setenv("TERM", "xterm-256color")
	buf := bytes.NewBuffer(nil)
	writer := NewWriter(Options{Writer: buf, Format: FormatRich, ForceColor: true, Theme: DefaultTheme()})

	if _, err := writer.Write([]byte(encodeEntry(Entry{Component: "lifecycle", Phase: "render"}, "templates took 1s") + "\n")); err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	if !strings.Contains(buf.String(), ansiCyan+"[lifecycle]"+ansiReset) {
		t.Fatalf("output = %q", buf.String())
	}
}

func TestThemeFromPaletteUsesPaletteColors(t *testing.T) {
	palette := &palettes.Palette{
		Colors: map[string]string{
			"primary":    "#112233",
			"secondary":  "#223344",
			"tertiary":   "#334455",
			"warning":    "#445566",
			"error":      "#556677",
			"text-muted": "#667788",
			"success":    "#778899",
		},
		Semantic: map[string]string{},
	}

	theme := ThemeFromPalette(palette)
	if theme.Component != "#112233" {
		t.Fatalf("Component = %q, want %q", theme.Component, "#112233")
	}
	if theme.Success != "#778899" {
		t.Fatalf("Success = %q, want %q", theme.Success, "#778899")
	}
	if theme.PhaseColor["transform"] != "#334455" {
		t.Fatalf("transform color = %q, want %q", theme.PhaseColor["transform"], "#334455")
	}
	if theme.PhaseColor["collect"] != "#778899" {
		t.Fatalf("collect color = %q, want %q", theme.PhaseColor["collect"], "#778899")
	}
}

func TestLoggerEncodesMetadata(t *testing.T) {
	entry, message := decodeEntry(encodeEntry(Entry{Component: "lifecycle", Phase: "render", Level: "info"}, "templates took 1s"))
	if entry.Component != "lifecycle" || entry.Phase != "render" || entry.Level != "info" {
		t.Fatalf("entry = %+v", entry)
	}
	if message != "templates took 1s" {
		t.Fatalf("message = %q", message)
	}
}

func TestStylePrefixedMessage(t *testing.T) {
	tests := []struct {
		name    string
		message string
		want    string
	}{
		{"already prefixed", "Warning: disk full", "Warning: disk full"},
		{"lowercase prefix", "warning: disk full", "Warning: disk full"},
		{"unprefixed keeps full message", "config warning: license: missing", "Warning: config warning: license: missing"},
		{"short unprefixed", "oops", "Warning: oops"},
		{"empty", "", "Warning:"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stylePrefixedMessage(tt.message, "Warning:", "", false); got != tt.want {
				t.Fatalf("stylePrefixedMessage(%q) = %q, want %q", tt.message, got, tt.want)
			}
		})
	}
}
