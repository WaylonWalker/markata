// Package benchmarks provides end-to-end build performance benchmarks for markata-go.
//
// These benchmarks use a deterministic sample site fixture at benchmarks/site/
// to measure build performance in a reproducible way.
//
// Run benchmarks with:
//
//	go test -bench=. -run=^$ ./benchmarks/...
//
// For profiling:
//
//	go test -bench=BenchmarkBuild -run=^$ -cpuprofile=cpu.prof -memprofile=mem.prof ./benchmarks/...
package benchmarks

import (
	"io"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/WaylonWalker/markata/pkg/config"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/plugins"
)

// fixtureDir is the path to the benchmark site fixture.
const fixtureDir = "site"

// BenchmarkBuild_EndToEnd runs a complete build through all lifecycle stages.
// This is the primary benchmark for measuring overall build performance.
func BenchmarkBuild_EndToEnd(b *testing.B) {
	// Ensure fixture exists
	fixturePath := filepath.Join(fixtureDir, "markata-go.toml")
	if _, err := os.Stat(fixturePath); os.IsNotExist(err) {
		b.Skipf("Benchmark fixture not found at %s - run 'go run generate_posts.go' first", fixturePath)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		// Create a fresh temp directory for each iteration
		outputDir := b.TempDir()
		b.StartTimer()

		m, err := setupManager(fixturePath, outputDir)
		if err != nil {
			b.Fatalf("Failed to setup manager: %v", err)
		}

		if err := m.Run(); err != nil {
			b.Fatalf("Build failed: %v", err)
		}
	}
}

// BenchmarkBuild_Incremental measures incremental build performance after a single file change.
// It runs a full build to warm caches, then times a second build after touching one file.
func BenchmarkBuild_Incremental(b *testing.B) {
	fixturePath := filepath.Join(fixtureDir, "markata-go.toml")
	if _, err := os.Stat(fixturePath); os.IsNotExist(err) {
		b.Skipf("Benchmark fixture not found at %s - run 'go run generate_posts.go' first", fixturePath)
	}

	changedFile := filepath.Join("posts", "blog", "2024", "01", "post-001.md")

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		b.StopTimer()

		workDir := b.TempDir()
		if err := copyDir(fixtureDir, workDir); err != nil {
			b.Fatalf("Failed to copy fixture: %v", err)
		}

		configPath := filepath.Join(workDir, "markata-go.toml")
		outputDir := filepath.Join(workDir, "output")

		// Warm caches with a full build
		m, err := setupManager(configPath, outputDir)
		if err != nil {
			b.Fatalf("Failed to setup manager: %v", err)
		}
		if err := m.Run(); err != nil {
			b.Fatalf("Warm build failed: %v", err)
		}

		// Touch one file to trigger incremental rebuild
		if err := touchFile(filepath.Join(workDir, changedFile)); err != nil {
			b.Fatalf("Failed to touch file: %v", err)
		}

		b.StartTimer()

		m, err = setupManager(configPath, outputDir)
		if err != nil {
			b.Fatalf("Failed to setup manager: %v", err)
		}
		if err := m.Run(); err != nil {
			b.Fatalf("Incremental build failed: %v", err)
		}

		b.StopTimer()
	}
}

// BenchmarkBuild_Concurrency tests build performance at different concurrency levels.
func BenchmarkBuild_Concurrency(b *testing.B) {
	fixturePath := filepath.Join(fixtureDir, "markata-go.toml")
	if _, err := os.Stat(fixturePath); os.IsNotExist(err) {
		b.Skipf("Benchmark fixture not found at %s", fixturePath)
	}

	concurrencies := []int{1, 2, 4, 8}
	if runtime.NumCPU() >= 16 {
		concurrencies = append(concurrencies, 16)
	}

	for _, conc := range concurrencies {
		b.Run(concurrencyName(conc), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				b.StopTimer()
				outputDir := b.TempDir()
				b.StartTimer()

				m, err := setupManager(fixturePath, outputDir)
				if err != nil {
					b.Fatalf("Failed to setup manager: %v", err)
				}
				m.SetConcurrency(conc)

				if err := m.Run(); err != nil {
					b.Fatalf("Build failed: %v", err)
				}
			}
		})
	}
}

// BenchmarkStage_Glob measures the glob (file discovery) stage.
func BenchmarkStage_Glob(b *testing.B) {
	fixturePath := filepath.Join(fixtureDir, "markata-go.toml")
	if _, err := os.Stat(fixturePath); os.IsNotExist(err) {
		b.Skipf("Benchmark fixture not found at %s", fixturePath)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		outputDir := b.TempDir()
		m, err := setupManager(fixturePath, outputDir)
		if err != nil {
			b.Fatalf("Failed to setup manager: %v", err)
		}
		b.StartTimer()

		if err := m.RunTo(lifecycle.StageGlob); err != nil {
			b.Fatalf("Glob stage failed: %v", err)
		}
	}
}

// BenchmarkStage_Load measures the load (file reading + frontmatter parsing) stage.
func BenchmarkStage_Load(b *testing.B) {
	fixturePath := filepath.Join(fixtureDir, "markata-go.toml")
	if _, err := os.Stat(fixturePath); os.IsNotExist(err) {
		b.Skipf("Benchmark fixture not found at %s", fixturePath)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		outputDir := b.TempDir()
		m, err := setupManager(fixturePath, outputDir)
		if err != nil {
			b.Fatalf("Failed to setup manager: %v", err)
		}
		b.StartTimer()

		if err := m.RunTo(lifecycle.StageLoad); err != nil {
			b.Fatalf("Load stage failed: %v", err)
		}
	}
}

// BenchmarkStage_Transform measures the transform (pre-rendering processing) stage.
func BenchmarkStage_Transform(b *testing.B) {
	fixturePath := filepath.Join(fixtureDir, "markata-go.toml")
	if _, err := os.Stat(fixturePath); os.IsNotExist(err) {
		b.Skipf("Benchmark fixture not found at %s", fixturePath)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		outputDir := b.TempDir()
		m, err := setupManager(fixturePath, outputDir)
		if err != nil {
			b.Fatalf("Failed to setup manager: %v", err)
		}
		b.StartTimer()

		if err := m.RunTo(lifecycle.StageTransform); err != nil {
			b.Fatalf("Transform stage failed: %v", err)
		}
	}
}

// BenchmarkStage_Render measures the render (markdown to HTML) stage.
func BenchmarkStage_Render(b *testing.B) {
	fixturePath := filepath.Join(fixtureDir, "markata-go.toml")
	if _, err := os.Stat(fixturePath); os.IsNotExist(err) {
		b.Skipf("Benchmark fixture not found at %s", fixturePath)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		outputDir := b.TempDir()
		m, err := setupManager(fixturePath, outputDir)
		if err != nil {
			b.Fatalf("Failed to setup manager: %v", err)
		}
		b.StartTimer()

		if err := m.RunTo(lifecycle.StageRender); err != nil {
			b.Fatalf("Render stage failed: %v", err)
		}
	}
}

// BenchmarkStage_Collect measures the collect (feed generation) stage.
func BenchmarkStage_Collect(b *testing.B) {
	fixturePath := filepath.Join(fixtureDir, "markata-go.toml")
	if _, err := os.Stat(fixturePath); os.IsNotExist(err) {
		b.Skipf("Benchmark fixture not found at %s", fixturePath)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		outputDir := b.TempDir()
		m, err := setupManager(fixturePath, outputDir)
		if err != nil {
			b.Fatalf("Failed to setup manager: %v", err)
		}
		b.StartTimer()

		if err := m.RunTo(lifecycle.StageCollect); err != nil {
			b.Fatalf("Collect stage failed: %v", err)
		}
	}
}

// BenchmarkStage_Write measures the write (output generation) stage.
func BenchmarkStage_Write(b *testing.B) {
	fixturePath := filepath.Join(fixtureDir, "markata-go.toml")
	if _, err := os.Stat(fixturePath); os.IsNotExist(err) {
		b.Skipf("Benchmark fixture not found at %s", fixturePath)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		outputDir := b.TempDir()
		m, err := setupManager(fixturePath, outputDir)
		if err != nil {
			b.Fatalf("Failed to setup manager: %v", err)
		}
		b.StartTimer()

		if err := m.RunTo(lifecycle.StageWrite); err != nil {
			b.Fatalf("Write stage failed: %v", err)
		}
	}
}

// setupManager creates a lifecycle manager configured for the benchmark fixture.
func setupManager(configPath, outputDir string) (*lifecycle.Manager, error) {
	// Get absolute path for config
	absConfigPath, err := filepath.Abs(configPath)
	if err != nil {
		return nil, err
	}

	// Load config
	cfg, err := config.Load(absConfigPath)
	if err != nil {
		return nil, err
	}

	// Override output directory
	cfg.OutputDir = outputDir

	// Create manager
	m := lifecycle.NewManager()

	// Convert models.Config to lifecycle.Config
	lcConfig := &lifecycle.Config{
		ContentDir:   filepath.Dir(absConfigPath),
		OutputDir:    outputDir,
		GlobPatterns: cfg.GlobConfig.Patterns,
		Extra:        make(map[string]interface{}),
	}

	// Copy config values to Extra for plugins to access
	lcConfig.Extra["url"] = cfg.URL
	lcConfig.Extra["title"] = cfg.Title
	lcConfig.Extra["description"] = cfg.Description
	lcConfig.Extra["author"] = cfg.Author
	lcConfig.Extra["templates_dir"] = cfg.TemplatesDir
	lcConfig.Extra["assets_dir"] = cfg.AssetsDir
	lcConfig.Extra["feeds"] = cfg.Feeds
	lcConfig.Extra["feed_defaults"] = cfg.FeedDefaults

	m.SetConfig(lcConfig)

	// Set concurrency from config if specified
	if cfg.Concurrency > 0 {
		m.SetConcurrency(cfg.Concurrency)
	}

	// Register all default plugins
	m.RegisterPlugins(plugins.DefaultPlugins()...)

	return m, nil
}

func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		targetPath := filepath.Join(dst, relPath)
		if info.IsDir() {
			return os.MkdirAll(targetPath, info.Mode())
		}

		return copyFile(path, targetPath, info.Mode())
	})
}

func copyFile(src, dst string, mode os.FileMode) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}

	dstFile, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

func touchFile(path string) error {
	now := time.Now()
	return os.Chtimes(path, now, now)
}

// concurrencyName returns a descriptive name for the concurrency level.
func concurrencyName(n int) string {
	switch n {
	case 1:
		return "serial"
	case 2:
		return "02workers"
	case 4:
		return "04workers"
	case 8:
		return "08workers"
	case 16:
		return "16workers"
	default:
		return "auto"
	}
}
