package plugins

import (
	"encoding/base64"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/visualizations"
)

func TestInlineExcalidrawSVGFonts_EmbedsBundledFonts(t *testing.T) {
	font := "Excalifont/Excalifont-Regular-a88b72a24fb54c9f94e3b5fdaa7481c9.woff2"
	svg := []byte(`<svg xmlns="http://www.w3.org/2000/svg" width="20" height="10"><style>@font-face{font-family:Excalifont;src:url("` + excalidrawFontURLPrefix + font + `")}</style><rect width="20" height="10"/></svg>`)

	got, err := inlineExcalidrawSVGFonts(svg)
	if err != nil {
		t.Fatalf("inlineExcalidrawSVGFonts() error = %v", err)
	}
	if strings.Contains(string(got), "https://esm.sh/") || !strings.Contains(string(got), "url(data:font/woff2;base64,") {
		t.Fatalf("inlineExcalidrawSVGFonts() did not embed its local font: %s", got)
	}
	if _, err := visualizations.NewSVGAsset(got); err != nil {
		t.Fatalf("inlined Excalidraw SVG failed the shared validator: %v", err)
	}
}

func TestInlineExcalidrawSVGFonts_LeavesOtherRemoteCSSForValidator(t *testing.T) {
	svg := []byte(`<svg xmlns="http://www.w3.org/2000/svg" width="20" height="10"><style>@font-face{src:url(https://fonts.example.test/font.woff2)}</style><rect width="20" height="10"/></svg>`)

	got, err := inlineExcalidrawSVGFonts(svg)
	if err != nil {
		t.Fatalf("inlineExcalidrawSVGFonts() error = %v", err)
	}
	if _, err := visualizations.NewSVGAsset(got); err == nil || !strings.Contains(err.Error(), "external or active CSS") {
		t.Fatalf("shared validator error = %v, want external CSS rejection", err)
	}
}

func TestInlineExcalidrawSVGFonts_RejectsUnlistedBundledFont(t *testing.T) {
	svg := []byte(`<svg xmlns="http://www.w3.org/2000/svg" width="20" height="10"><style>@font-face{src:url(` + excalidrawFontURLPrefix + `Unknown/Unknown.woff2)}</style></svg>`)

	if _, err := inlineExcalidrawSVGFonts(svg); err == nil || !strings.Contains(err.Error(), "not listed in the bundled manifest") {
		t.Fatalf("inlineExcalidrawSVGFonts() error = %v, want manifest rejection", err)
	}
}

func TestInlineExcalidrawSVGFonts_EnforcesExpandedSVGLimit(t *testing.T) {
	font := "Excalifont/Excalifont-Regular-a88b72a24fb54c9f94e3b5fdaa7481c9.woff2"
	fontData, err := visualizations.ReadExcalidrawAsset("excalidraw-fonts/" + font)
	if err != nil {
		t.Fatal(err)
	}
	fontURL := excalidrawFontURLPrefix + font
	inputURLSize := len("url(") + len(fontURL) + len(")")
	outputURLSize := len("url(data:font/woff2;base64,") + base64.StdEncoding.EncodedLen(len(fontData)) + len(")")
	fontRefs := strings.Repeat("src:url("+fontURL+");", visualizations.MaxSVGBytes/(outputURLSize-inputURLSize)+1)
	svg := []byte(`<svg xmlns="http://www.w3.org/2000/svg" width="20" height="10"><style>` + fontRefs + `</style></svg>`)
	if len(svg) >= visualizations.MaxSVGBytes {
		t.Fatalf("test input unexpectedly exceeds the SVG limit: %d bytes", len(svg))
	}
	if _, err := inlineExcalidrawSVGFonts(svg); err == nil || !strings.Contains(err.Error(), "after font embedding") {
		t.Fatalf("inlineExcalidrawSVGFonts() error = %v, want expanded-size rejection", err)
	}
}
