package plugins

import (
	"os"
	"path/filepath"
	"testing"
)

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
