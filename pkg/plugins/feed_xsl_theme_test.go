package plugins

import (
	"encoding/xml"
	"errors"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

func newFeedXSLThemeConfig(t *testing.T) *lifecycle.Config {
	t.Helper()
	return &lifecycle.Config{
		OutputDir: t.TempDir(),
		Extra: map[string]interface{}{
			"templates_dir": filepath.Join(t.TempDir(), "missing-templates"),
			"theme": models.ThemeConfig{
				Palette:      "catppuccin-latte",
				PaletteLight: "catppuccin-latte",
				PaletteDark:  "catppuccin-mocha",
				FallbackMode: "dark",
				Aesthetic:    "minimal",
				Fontpack:     "brush",
				TextSize:     "x-large",
			},
			"palette_light_effective": "catppuccin-latte",
			"palette_dark_effective":  "catppuccin-mocha",
			"fontpack_css":            "/* fonts */",
		},
	}
}

func assertWellFormedXML(t *testing.T, name string, data []byte) {
	t.Helper()
	dec := xml.NewDecoder(strings.NewReader(string(data)))
	for {
		_, err := dec.Token()
		if errors.Is(err, io.EOF) {
			return
		}
		if err != nil {
			t.Fatalf("%s is not well-formed XML: %v", name, err)
		}
	}
}

func TestCopyXSLStylesheets_AppliesSiteTheme(t *testing.T) {
	config := newFeedXSLThemeConfig(t)
	p := NewPublishFeedsPlugin()
	if err := p.copyXSLStylesheets(config, config.OutputDir); err != nil {
		t.Fatalf("copyXSLStylesheets: %v", err)
	}

	for _, name := range []string{"rss.xsl", "atom.xsl"} {
		data, err := os.ReadFile(filepath.Join(config.OutputDir, name))
		if err != nil {
			t.Fatalf("reading %s: %v", name, err)
		}
		assertWellFormedXML(t, name, data)
		out := string(data)
		for _, want := range []string{
			`root.dataset.palette = 'catppuccin\u002Dmocha'`,
			"root.dataset.aesthetic = 'minimal'",
			"root.dataset.fontpack = 'brush'",
			`root.dataset.textSize = 'x\u002Dlarge'`,
			"localStorage.getItem('theme-palette-' + mode)",
			"<![CDATA[",
		} {
			if !strings.Contains(out, want) {
				t.Errorf("%s missing %q", name, want)
			}
		}
		// Other tests may register asset hashes, so accept hashed file names.
		for _, css := range []string{"palette", "fonts", "aesthetic"} {
			pattern := regexp.MustCompile(`href="/css/` + css + `(\.[0-9a-f]+)?\.css" />`)
			if !pattern.MatchString(out) {
				t.Errorf("%s missing self-closed link to css/%s.css", name, css)
			}
		}
		if strings.Contains(out, feedXSLThemeStart) {
			t.Errorf("%s still contains the theme marker", name)
		}
	}
}

func TestCopyXSLStylesheets_LeavesCustomXSLWithoutMarker(t *testing.T) {
	config := newFeedXSLThemeConfig(t)
	templatesDir := t.TempDir()
	config.Extra["templates_dir"] = templatesDir
	custom := `<?xml version="1.0"?><xsl:stylesheet version="1.0" xmlns:xsl="http://www.w3.org/1999/XSL/Transform"/>`
	if err := os.WriteFile(filepath.Join(templatesDir, "rss.xsl"), []byte(custom), 0o600); err != nil {
		t.Fatal(err)
	}

	p := NewPublishFeedsPlugin()
	if err := p.copyXSLStylesheets(config, config.OutputDir); err != nil {
		t.Fatalf("copyXSLStylesheets: %v", err)
	}
	got, err := os.ReadFile(filepath.Join(config.OutputDir, "rss.xsl"))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != custom {
		t.Errorf("custom rss.xsl was modified:\n%s", got)
	}
}

func TestXMLSafeHead(t *testing.T) {
	in := "<script>if (a < b && c) {}</script>\n<link rel=\"stylesheet\" href=\"/a.css\">\n<script src=\"/x.js\"></script>"
	got := xmlSafeHead(in)
	want := "<script><![CDATA[if (a < b && c) {}]]></script>\n<link rel=\"stylesheet\" href=\"/a.css\" />\n<script src=\"/x.js\"></script>"
	if got != want {
		t.Errorf("xmlSafeHead() =\n%s\nwant\n%s", got, want)
	}
	assertWellFormedXML(t, "head", []byte("<head>"+got+"</head>"))
}

func TestCopyXSLStylesheets_EscapesPaletteForScript(t *testing.T) {
	config := newFeedXSLThemeConfig(t)
	theme, ok := config.Extra["theme"].(models.ThemeConfig)
	if !ok {
		t.Fatal("theme is not a ThemeConfig")
	}
	theme.Palette = "st.-patrick's-day"
	config.Extra["theme"] = theme
	config.Extra["palette_dark_effective"] = "st.-patrick's-day"

	p := NewPublishFeedsPlugin()
	if err := p.copyXSLStylesheets(config, config.OutputDir); err != nil {
		t.Fatalf("copyXSLStylesheets: %v", err)
	}
	data, err := os.ReadFile(filepath.Join(config.OutputDir, "rss.xsl"))
	if err != nil {
		t.Fatalf("reading rss.xsl: %v", err)
	}
	assertWellFormedXML(t, "rss.xsl", data)
	out := string(data)
	if !strings.Contains(out, `root.dataset.palette = 'st\u002E\u002Dpatrick\u0027s\u002Dday'`) {
		t.Errorf("palette not JS-escaped:\n%s", out)
	}
	if strings.Contains(out, "patrick&#39;s") {
		t.Error("palette was HTML-escaped inside a script")
	}
}
