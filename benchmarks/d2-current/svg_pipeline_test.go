package d2native

import (
	"context"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/visualizations"
)

// Exercise the production Markata asset contract, not just XML decoding.
// Any SVG feature that the conservative validator rejects must be handled
// explicitly by a future D2 adapter rather than silently falling back to JS.
func TestD2CurrentSVGCanPublishAsMarkataAsset(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}
	for _, fixture := range []struct {
		name   string
		source string
	}{
		{name: "simple", source: graph(3)},
		{name: "medium", source: graph(18)},
		{name: "large", source: graph(72)},
		{name: "labels", source: "a: Hello\nb: World\na -> b: Diagram\n"},
		{name: "shapes", source: "a: {shape: cylinder}\nb: {shape: diamond}\na -> b\n"},
	} {
		t.Run(fixture.name, func(t *testing.T) {
			svg, err := engine.Render(context.Background(), fixture.source)
			if err != nil {
				t.Fatalf("D2 render: %v", err)
			}
			asset, err := visualizations.NewSVGAsset(svg)
			if err != nil {
				t.Fatalf("Markata SVG policy rejects D2 output: %v", err)
			}
			if asset.Width <= 0 || asset.Height <= 0 {
				t.Fatalf("missing dimensions: %d x %d", asset.Width, asset.Height)
			}
			dir := t.TempDir()
			name, err := visualizations.Publish(dir, asset)
			if err != nil {
				t.Fatalf("publish SVG: %v", err)
			}
			if !strings.HasSuffix(name, ".svg") {
				t.Fatalf("unexpected published path: %q", name)
			}
			markup, err := visualizations.FigureHTML(asset, "D2 "+fixture.name, "Backend-rendered diagram", false)
			if err != nil {
				t.Fatalf("figure markup: %v", err)
			}
			if !strings.Contains(markup, "width=") || !strings.Contains(markup, "height=") {
				t.Fatalf("figure missing dimensions: %s", markup)
			}
			t.Logf("svg_bytes=%d dimensions=%dx%d asset=%s", len(svg), asset.Width, asset.Height, name)
		})
	}
}
