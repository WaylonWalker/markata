package lifecycle

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/WaylonWalker/markata/pkg/models"
)

// BenchmarkManager_ColdStart measures performance with no cache (fresh manager).
func BenchmarkManager_ColdStart(b *testing.B) {
	// Setup test data directory
	testDir := setupBenchmarkTestData(b, 100)
	defer os.RemoveAll(testDir)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m := NewManager()
		m.config.ContentDir = testDir
		m.config.GlobPatterns = []string{"**/*.md"}
		m.config.OutputDir = filepath.Join(testDir, "output")

		// Run through stages up to render
		if err := m.RunTo(StageRender); err != nil {
			b.Fatal(err)
		}

		// Reset for next iteration (simulate cold start)
		m.Reset()
	}
}

// BenchmarkManager_HotCache measures performance with warm cache.
func BenchmarkManager_HotCache(b *testing.B) {
	testDir := setupBenchmarkTestData(b, 100)
	defer os.RemoveAll(testDir)

	// Create manager and warm up the cache
	m := NewManager()
	m.config.ContentDir = testDir
	m.config.GlobPatterns = []string{"**/*.md"}
	m.config.OutputDir = filepath.Join(testDir, "output")

	// Initial run to warm cache
	if err := m.RunTo(StageRender); err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Keep same manager but reset stages
		m.mu.Lock()
		m.stagesRun = make(map[Stage]bool)
		m.posts = make([]*models.Post, 0)
		m.files = make([]string, 0)
		m.feeds = make([]*Feed, 0)
		m.warnings = make([]*HookError, 0)
		m.currentStage = ""
		// Note: NOT clearing cache - that's the "hot" part
		m.mu.Unlock()

		if err := m.RunTo(StageRender); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkManager_GlobStage measures just the globbing phase.
func BenchmarkManager_GlobStage(b *testing.B) {
	testDir := setupBenchmarkTestData(b, 500)
	defer os.RemoveAll(testDir)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m := NewManager()
		m.config.ContentDir = testDir
		m.config.GlobPatterns = []string{"**/*.md"}

		if err := m.RunTo(StageGlob); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkManager_LoadStage measures file loading and parsing.
func BenchmarkManager_LoadStage(b *testing.B) {
	testDir := setupBenchmarkTestData(b, 100)
	defer os.RemoveAll(testDir)

	// Pre-glob to have files ready
	m := NewManager()
	m.config.ContentDir = testDir
	m.config.GlobPatterns = []string{"**/*.md"}
	if err := m.RunTo(StageGlob); err != nil {
		b.Fatal(err)
	}
	files := m.Files()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m := NewManager()
		m.config.ContentDir = testDir
		m.SetFiles(files)

		// Run just the load stage hooks
		hookErrors := runLoadHooks(m)
		if hookErrors.HasCritical() {
			b.Fatal(hookErrors)
		}
	}
}

// BenchmarkProcessPostsConcurrently measures concurrent post processing.
func BenchmarkProcessPostsConcurrently(b *testing.B) {
	// Create test posts
	posts := make([]*models.Post, 500)
	for i := 0; i < 500; i++ {
		title := fmt.Sprintf("Test Post %d", i)
		posts[i] = &models.Post{
			Path:    fmt.Sprintf("post-%d.md", i),
			Slug:    fmt.Sprintf("post-%d", i),
			Title:   &title,
			Content: generateBenchmarkContent(i),
		}
	}

	benchConcurrency := []int{1, 2, 4, 8, 16}
	for _, conc := range benchConcurrency {
		b.Run(fmt.Sprintf("concurrency-%d", conc), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				m := NewManager()
				m.SetConcurrency(conc)
				m.SetPosts(posts)

				err := m.ProcessPostsConcurrently(func(p *models.Post) error {
					// Simulate some work
					_ = len(p.Content)
					return nil
				})
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// BenchmarkProcessPostsConcurrentlyLargeScale measures worker pool scaling with large post counts.
// This benchmark validates that the worker pool approach scales efficiently to thousands of posts.
func BenchmarkProcessPostsConcurrentlyLargeScale(b *testing.B) {
	postCounts := []int{1000, 5000, 10000}

	for _, count := range postCounts {
		// Pre-allocate posts outside the benchmark loop
		posts := make([]*models.Post, count)
		for i := 0; i < count; i++ {
			title := fmt.Sprintf("Test Post %d", i)
			posts[i] = &models.Post{
				Path:    fmt.Sprintf("post-%d.md", i),
				Slug:    fmt.Sprintf("post-%d", i),
				Title:   &title,
				Content: generateBenchmarkContent(i % 100), // Reuse content patterns
			}
		}

		b.Run(fmt.Sprintf("posts-%d", count), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				m := NewManager()
				m.SetConcurrency(8) // Fixed concurrency for scaling comparison
				m.SetPosts(posts)

				err := m.ProcessPostsConcurrently(func(p *models.Post) error {
					// Simulate some work
					_ = len(p.Content)
					return nil
				})
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// BenchmarkFilter measures filter expression evaluation.
func BenchmarkFilter(b *testing.B) {
	// Create test posts with various attributes
	posts := make([]*models.Post, 1000)
	for i := 0; i < 1000; i++ {
		title := fmt.Sprintf("Test Post %d", i)
		published := i%2 == 0
		draft := i%5 == 0
		tags := []string{"go", "benchmark"}
		if i%3 == 0 {
			tags = append(tags, "performance")
		}
		posts[i] = &models.Post{
			Path:      fmt.Sprintf("post-%d.md", i),
			Slug:      fmt.Sprintf("post-%d", i),
			Title:     &title,
			Published: published,
			Draft:     draft,
			Tags:      tags,
		}
	}

	m := NewManager()
	m.SetPosts(posts)

	testCases := []struct {
		name string
		expr string
	}{
		{"simple_equals", "published==true"},
		{"simple_not_equals", "draft!=true"},
		{"contains", "tags contains performance"},
		{"and_condition", "published==true and draft!=true"},
		{"or_condition", "published==true or draft==true"},
		{"complex", "published==true and draft!=true and tags contains go"},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := m.Filter(tc.expr)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// BenchmarkMemoryCache measures cache operations.
func BenchmarkMemoryCache(b *testing.B) {
	cache := newMemoryCache()

	b.Run("Set", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			cache.Set(fmt.Sprintf("key-%d", i), i)
		}
	})

	// Pre-populate for Get benchmark
	for i := 0; i < 10000; i++ {
		cache.Set(fmt.Sprintf("key-%d", i), i)
	}

	b.Run("Get_Hit", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			cache.Get(fmt.Sprintf("key-%d", i%10000))
		}
	})

	b.Run("Get_Miss", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			cache.Get(fmt.Sprintf("nonexistent-%d", i))
		}
	})
}

// BenchmarkPluginSorting measures the overhead of sorting plugins by priority.
func BenchmarkPluginSorting(b *testing.B) {
	// Create mock plugins
	plugins := make([]Plugin, 50)
	for i := 0; i < 50; i++ {
		plugins[i] = &mockPlugin{name: fmt.Sprintf("plugin-%d", i)}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = sortPluginsByPriority(plugins, StageRender)
	}
}

// Helper types and functions for benchmarks

type mockPlugin struct {
	name string
}

func (p *mockPlugin) Name() string { return p.name }

// setupBenchmarkTestData creates a temporary directory with test markdown files.
func setupBenchmarkTestData(b *testing.B, numFiles int) string {
	b.Helper()

	dir, err := os.MkdirTemp("", "markata-bench-*")
	if err != nil {
		b.Fatal(err)
	}

	// Create nested structure
	subdirs := []string{"posts", "docs", "guides", "blog/2024", "blog/2023"}
	for _, subdir := range subdirs {
		if err := os.MkdirAll(filepath.Join(dir, subdir), 0o755); err != nil {
			b.Fatal(err)
		}
	}

	// Create markdown files
	for i := 0; i < numFiles; i++ {
		subdir := subdirs[i%len(subdirs)]
		filename := fmt.Sprintf("post-%d.md", i)
		path := filepath.Join(dir, subdir, filename)
		content := generateBenchmarkFile(i)
		//nolint:gosec // G306: test files don't need restrictive permissions
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			b.Fatal(err)
		}
	}

	return dir
}

// generateBenchmarkFile creates a realistic markdown file for benchmarking.
func generateBenchmarkFile(index int) string {
	date := time.Now().AddDate(0, 0, -index)
	codeBlock := "```go\npackage main\n\nfunc main() {\n    fmt.Println(\"Hello, benchmark!\")\n}\n```"
	return fmt.Sprintf("---\ntitle: \"Benchmark Test Post %d\"\ndescription: \"This is a test post for benchmarking markata-go performance\"\ndate: %s\npublished: %t\ndraft: %t\ntags:\n  - benchmark\n  - testing\n  - go\n  - post-%d\n---\n\n# Benchmark Test Post %d\n\nThis is the content of test post %d. It contains various markdown elements\nfor realistic benchmarking.\n\n## Section 1\n\nHere's a paragraph with some **bold** and *italic* text.\n\n%s\n\n## Section 2\n\nA list:\n- Item 1\n- Item 2\n- Item 3\n\n## Code Example\n\n%s\n\n## Conclusion\n\nThis concludes test post %d.\n",
		index, date.Format("2006-01-02"), index%2 == 0, index%10 == 0, index, index, index, generateExtraContent(index), codeBlock, index)
}

// generateBenchmarkContent creates markdown content for in-memory tests.
func generateBenchmarkContent(index int) string {
	return fmt.Sprintf("# Post %d\n\nThis is content for post %d with various **markdown** elements.\n\n## Details\n\nSome more content here with code and [links](https://example.com).\n\n%s\n",
		index, index, generateExtraContent(index))
}

// generateExtraContent adds varying amounts of content based on index.
func generateExtraContent(index int) string {
	// Add varying content to simulate real-world variety
	if index%5 == 0 {
		return "> This is a blockquote with some important information\n> that spans multiple lines for emphasis.\n\n| Column 1 | Column 2 | Column 3 |\n|----------|----------|----------|\n| Data 1   | Data 2   | Data 3   |\n| Data 4   | Data 5   | Data 6   |\n"
	}
	if index%3 == 0 {
		return "![Image alt text](https://example.com/image.png)\n\n---\n\nSome additional paragraph text here.\n"
	}
	return ""
}
