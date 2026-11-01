//go:build !js

package plugins

import (
	"strings"
	"sync"
	"testing"

	"github.com/WaylonWalker/markata/pkg/models"
)

func TestNativeMermaidBackend_SVGAndTheme(t *testing.T) {
	cfg := models.NewMermaidConfig()
	cfg.Mode = mermaidModeGo
	cfg.GoConfig.FontFamily = "DejaVu Sans"
	cfg.GoConfig.FontSize = 19
	colors := &mermaidPaletteColors{
		Background: "#10241b", Surface: "#184a35",
		TextColor: "#f7fff8", Accent: "#82eac2",
	}
	renderer, err := newMermaidRenderer(cfg, colors)
	if err != nil {
		t.Fatal(err)
	}
	defer renderer.close()
	svg, err := renderer.render("flowchart LR\n  A[Load] --> B[Serve]")
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"<svg", "Load", "Serve", "#f7fff8", "#82eac2", "DejaVu Sans"} {
		if !strings.Contains(svg, want) {
			t.Fatalf("SVG missing %q: %s", want, svg)
		}
	}
	if strings.Contains(svg, "<script") {
		t.Fatal("native SVG must not ship script")
	}
}

func TestNativeMermaidBackend_InvalidSyntaxIsError(t *testing.T) {
	cfg := models.NewMermaidConfig()
	cfg.Mode = mermaidModeGo
	renderer, err := newMermaidRenderer(cfg, nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = renderer.render("not a mermaid diagram\nA->B"); err == nil {
		t.Fatal("unexpected success for unsupported Mermaid input")
	}
}

func TestNativeMermaidBackend_ConcurrentPalettesStayIsolated(t *testing.T) {
	a := &mermaidPaletteColors{Background: "#080c0a", Surface: "#114422", TextColor: "#fff8dd", Accent: "#71df80"}
	b := &mermaidPaletteColors{Background: "#fdf8fc", Surface: "#fff3f9", TextColor: "#222244", Accent: "#ae158c"}
	cfg := models.NewMermaidConfig()
	cfg.Mode = mermaidModeGo
	renderA := newNativeMermaidRenderer(cfg, a)
	renderB := newNativeMermaidRenderer(cfg, b)
	var wg sync.WaitGroup
	for i := 0; i < 32; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			renderer, color := renderA, "#fff8dd"
			if i%2 == 1 {
				renderer, color = renderB, "#222244"
			}
			svg, err := renderer.render("flowchart LR\n A[Start] --> B[End]")
			if err != nil {
				t.Errorf("render: %v", err)
				return
			}
			if !strings.Contains(svg, color) {
				t.Errorf("palette leaked: wanted %s", color)
			}
		}(i)
	}
	wg.Wait()
}

func TestRenderStandalone_NativeMermaidNoClientLibrary(t *testing.T) {
	text := "```mermaid\nflowchart LR\n  A[Load] --> B[Serve]\n```"
	result, err := RenderStandalone(text, StandaloneOptions{MermaidMode: "go"})
	if err != nil {
		t.Fatal(err)
	}
	if !result.Features["mermaid"] || !strings.Contains(result.HTML, "<svg") {
		t.Fatalf("native Mermaid not embedded: %+v", result)
	}
	if !strings.Contains(result.HTML, "Load") {
		t.Fatal("diagram lost source label")
	}
	for _, asset := range result.Assets {
		if strings.Contains(asset.ID, "mermaid.min.js") {
			t.Fatalf("native SVG should not ship Mermaid JS: %s", asset.ID)
		}
	}
	if result.BodyEndHTML != "" {
		t.Fatalf("native diagram unexpectedly requires JS: %q", result.BodyEndHTML)
	}
	if !strings.Contains(result.HeadHTML, "mermaid.css") {
		t.Fatalf("native diagram should use Mermaid presentation CSS: %s", result.HeadHTML)
	}
}

func TestRenderStandalone_NativeMermaidRejectsUnknownMode(t *testing.T) {
	if _, err := RenderStandalone("# Plain", StandaloneOptions{MermaidMode: "not-real"}); err == nil {
		t.Fatal("unexpected success for unsupported standalone Mermaid mode")
	}
}

func TestNativeMermaidBackend_FamiliesAndMalformedInput(t *testing.T) {
	cfg := models.NewMermaidConfig()
	cfg.Mode = mermaidModeGo
	renderer, err := newMermaidRenderer(cfg, nil)
	if err != nil {
		t.Fatal(err)
	}
	cases := map[string]string{
		"flowchart": "flowchart TD\n A[Start] --> B[End]",
		"sequence":  "sequenceDiagram\n Alice->>Bob: Hello",
		"class":     "classDiagram\n class Person {\n +String name\n }",
		"state":     "stateDiagram-v2\n [*] --> Active",
		"er":        "erDiagram\n CUSTOMER ||--o{ ORDER : places",
		"pie":       "pie\n \"Cats\" : 4\n \"Dogs\" : 3",
		"gitgraph":  "gitGraph\n commit\n branch dev\n commit",
		"mindmap":   "mindmap\n root((idea))\n  Design",
	}
	for name, source := range cases {
		t.Run(name, func(t *testing.T) {
			svg, err := renderer.render(source)
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(svg, "<svg") {
				t.Fatal("missing svg element")
			}
		})
	}
	if _, err = renderer.render("absolutely-not-a-mermaid-diagram"); err == nil {
		t.Fatal("invalid syntax succeeded")
	}
}
