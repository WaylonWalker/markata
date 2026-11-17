package plugins

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// BenchmarkStandaloneFeatureFences models a long editor page with many diagrams.
// The renderer is reused, as it is by the local workbench.
func BenchmarkStandaloneFeatureFences(b *testing.B) {
	source := "# Preview\n\n" + strings.Repeat("```mermaid\ngraph LR\nStart-->End\n```\n\n", 12) +
		"```chartjs\n{\"type\":\"bar\",\"data\":{\"labels\":[\"A\"],\"datasets\":[{\"data\":[1]}]}}\n```\n"
	renderer := NewStandaloneRenderer(StandaloneOptions{Unsafe: true, MermaidMode: "client"})
	b.ReportAllocs()
	b.SetBytes(int64(len(source)))
	for b.Loop() {
		if _, err := renderer.Render(source); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkRenderStandaloneFixtures reuses the site benchmark post and visual fixtures.
// File reads are excluded; renderer setup and conditional asset generation are included.
func BenchmarkRenderStandaloneFixtures(b *testing.B) {
	fixtures, err := filepath.Glob("../../rendering-fixtures/*.md")
	if err != nil {
		b.Fatal(err)
	}
	fixtures = append([]string{"../../benchmarks/site/posts/blog/2024/01/post-001.md"}, fixtures...)
	for _, path := range fixtures {
		source, err := os.ReadFile(path)
		if err != nil {
			b.Fatal(err)
		}
		b.Run(filepath.Base(path), func(b *testing.B) {
			b.ReportAllocs()
			b.SetBytes(int64(len(source)))
			for b.Loop() {
				if _, err := RenderStandalone(string(source), StandaloneOptions{Unsafe: true, InlineAssets: true}); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
