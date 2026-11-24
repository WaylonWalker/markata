package cmd

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/WaylonWalker/markata/pkg/config"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/plugins"
	"github.com/spf13/cobra"
)

// Benchmark scenario names.
const (
	ScenarioSmall  = "small"
	ScenarioMedium = "medium"
	ScenarioLarge  = "large"
)

// Benchmark scenario post counts.
const (
	SmallPostCount  = 50
	MediumPostCount = 200
	LargePostCount  = 500
)

// Report format constants.
const (
	reportFormatJSON = "json"
)

var (
	// benchmarkScenario specifies a single scenario to run.
	benchmarkScenario string

	// benchmarkReport specifies the output report format.
	benchmarkReport string

	// benchmarkKeep keeps the generated test files after benchmark.
	benchmarkKeep bool
)

// BenchmarkScenario defines a benchmark test scenario.
type BenchmarkScenario struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	PostCount   int    `json:"post_count"`
	FeedCount   int    `json:"feed_count"`
}

// BenchmarkResult holds the metrics for a single benchmark run.
type BenchmarkResult struct {
	Scenario     string                   `json:"scenario"`
	PostCount    int                      `json:"post_count"`
	FeedCount    int                      `json:"feed_count"`
	TotalTime    time.Duration            `json:"total_time_ns"`
	TotalTimeSec float64                  `json:"total_time_sec"`
	PostsPerSec  float64                  `json:"posts_per_sec"`
	MemoryStats  BenchmarkMemoryStats     `json:"memory_stats"`
	StageTimes   map[string]time.Duration `json:"stage_times_ns"`
	StageTimeSec map[string]float64       `json:"stage_times_sec"`
	PluginCount  int                      `json:"plugin_count"`
	Success      bool                     `json:"success"`
	ErrorMessage string                   `json:"error_message,omitempty"`
}

// BenchmarkMemoryStats holds memory usage statistics.
type BenchmarkMemoryStats struct {
	PeakHeapMB   float64 `json:"peak_heap_mb"`
	AllocMB      float64 `json:"alloc_mb"`
	TotalAllocMB float64 `json:"total_alloc_mb"`
	GCCycles     uint32  `json:"gc_cycles"`
}

// BenchmarkReport holds the complete benchmark results.
type BenchmarkReport struct {
	Version    string             `json:"version"`
	Timestamp  time.Time          `json:"timestamp"`
	SystemInfo BenchmarkSystem    `json:"system_info"`
	Results    []*BenchmarkResult `json:"results"`
}

// BenchmarkSystem holds system information.
type BenchmarkSystem struct {
	OS         string `json:"os"`
	Arch       string `json:"arch"`
	NumCPU     int    `json:"num_cpu"`
	GoVersion  string `json:"go_version"`
	Gomaxprocs int    `json:"gomaxprocs"`
}

// predefinedScenarios defines the available benchmark scenarios.
var predefinedScenarios = map[string]BenchmarkScenario{
	ScenarioSmall: {
		Name:        ScenarioSmall,
		Description: "Small site (50 posts, 3 feeds)",
		PostCount:   SmallPostCount,
		FeedCount:   3,
	},
	ScenarioMedium: {
		Name:        ScenarioMedium,
		Description: "Medium site (200 posts, 5 feeds)",
		PostCount:   MediumPostCount,
		FeedCount:   5,
	},
	ScenarioLarge: {
		Name:        ScenarioLarge,
		Description: "Large site (500 posts, 8 feeds)",
		PostCount:   LargePostCount,
		FeedCount:   8,
	},
}

// benchmarkCmd represents the benchmark command.
var benchmarkCmd = &cobra.Command{
	Use:   "benchmark",
	Short: "Run performance benchmarks",
	Long: `Run comprehensive performance benchmarks with realistic content.

The benchmark generates test sites with rich content including:
- Comprehensive frontmatter (tags, categories, dates, custom fields)
- Complex markdown (code blocks, tables, lists, blockquotes)
- Multiple feeds with filtering and pagination

Available scenarios:
  small   - 50 posts, 3 feeds (personal blog)
  medium  - 200 posts, 5 feeds (professional blog)
  large   - 500 posts, 8 feeds (documentation site)

Example usage:
  markata-go benchmark                     # Run all scenarios
  markata-go benchmark --scenario small    # Run specific scenario
  markata-go benchmark --report json       # Generate JSON report
  markata-go benchmark --keep              # Keep generated test files
  markata-go benchmark -v                  # Verbose output with stage timings`,
	RunE: runBenchmarkCommand,
}

func init() {
	rootCmd.AddCommand(benchmarkCmd)

	benchmarkCmd.Flags().StringVar(&benchmarkScenario, "scenario", "", "run specific scenario (small, medium, large)")
	benchmarkCmd.Flags().StringVar(&benchmarkReport, "report", "", "output report format (json)")
	benchmarkCmd.Flags().BoolVar(&benchmarkKeep, "keep", false, "keep generated test files after benchmark")
}

func runBenchmarkCommand(_ *cobra.Command, _ []string) error {
	// Determine which scenarios to run
	scenarios := getScenarios()

	if len(scenarios) == 0 {
		return fmt.Errorf("no valid scenarios to run")
	}

	// Print header unless JSON output
	if benchmarkReport != reportFormatJSON {
		fmt.Println("Markata-go Performance Benchmark")
		fmt.Println("=================================")
		fmt.Printf("Version: %s\n", Version)
		fmt.Printf("Date: %s\n", time.Now().Format(time.RFC3339))
		fmt.Printf("System: %s/%s, %d CPUs\n", runtime.GOOS, runtime.GOARCH, runtime.NumCPU())
		fmt.Println()
	}

	results := make([]*BenchmarkResult, 0, len(scenarios))

	for _, scenario := range scenarios {
		if benchmarkReport != reportFormatJSON {
			if verbose {
				fmt.Printf("\n--- Running %s benchmark (%s) ---\n", scenario.Name, scenario.Description)
			} else {
				fmt.Printf("Running %s benchmark (%d posts)... ", scenario.Name, scenario.PostCount)
			}
		}

		result := runScenarioBenchmark(scenario)
		results = append(results, result)

		if benchmarkReport != reportFormatJSON {
			if result.Success {
				if verbose {
					fmt.Printf("Completed in %.2fs (%.1f posts/sec, %.1fMB peak memory)\n",
						result.TotalTimeSec, result.PostsPerSec, result.MemoryStats.PeakHeapMB)
				} else {
					fmt.Printf("%.2fs (%.1f posts/sec)\n", result.TotalTimeSec, result.PostsPerSec)
				}
			} else {
				fmt.Printf("FAILED: %s\n", result.ErrorMessage)
			}
		}
	}

	// Generate output
	if benchmarkReport == reportFormatJSON {
		return outputJSONReport(results)
	}

	// Print summary table
	fmt.Println()
	printBenchmarkSummary(results)

	return nil
}

// getScenarios returns the scenarios to run based on flags.
func getScenarios() []BenchmarkScenario {
	if benchmarkScenario != "" {
		if s, ok := predefinedScenarios[benchmarkScenario]; ok {
			return []BenchmarkScenario{s}
		}
		fmt.Printf("Warning: unknown scenario %q, running all scenarios\n", benchmarkScenario)
	}

	// Return scenarios in order: small, medium, large
	return []BenchmarkScenario{
		predefinedScenarios[ScenarioSmall],
		predefinedScenarios[ScenarioMedium],
		predefinedScenarios[ScenarioLarge],
	}
}

// runScenarioBenchmark runs a benchmark for a specific scenario.
func runScenarioBenchmark(scenario BenchmarkScenario) *BenchmarkResult {
	result := &BenchmarkResult{
		Scenario:     scenario.Name,
		PostCount:    scenario.PostCount,
		FeedCount:    scenario.FeedCount,
		StageTimes:   make(map[string]time.Duration),
		StageTimeSec: make(map[string]float64),
	}

	// Create temp directory for test site
	tempDir, err := os.MkdirTemp("", fmt.Sprintf("markata-benchmark-%s-*", scenario.Name))
	if err != nil {
		result.ErrorMessage = fmt.Sprintf("failed to create temp dir: %v", err)
		return result
	}

	if !benchmarkKeep {
		defer os.RemoveAll(tempDir)
	} else if verbose {
		fmt.Printf("Test site directory: %s\n", tempDir)
	}

	// Generate test content
	if verbose {
		fmt.Printf("Generating %d test posts with rich content...\n", scenario.PostCount)
	}

	postsDir := filepath.Join(tempDir, "posts")
	if err := os.MkdirAll(postsDir, 0o755); err != nil {
		result.ErrorMessage = fmt.Sprintf("failed to create posts dir: %v", err)
		return result
	}

	genStart := time.Now()
	if err := generateRichTestPosts(postsDir, scenario.PostCount); err != nil {
		result.ErrorMessage = fmt.Sprintf("failed to generate posts: %v", err)
		return result
	}

	if verbose {
		fmt.Printf("Generated %d posts in %.2fs\n", scenario.PostCount, time.Since(genStart).Seconds())
	}

	// Create comprehensive config with multiple feeds
	configPath := filepath.Join(tempDir, "markata-go.toml")
	configContent := generateBenchmarkConfig(tempDir, scenario)

	if err := os.WriteFile(configPath, []byte(configContent), 0o600); err != nil {
		result.ErrorMessage = fmt.Sprintf("failed to write config: %v", err)
		return result
	}

	// Force garbage collection before benchmark
	runtime.GC()
	var memBefore runtime.MemStats
	runtime.ReadMemStats(&memBefore)
	gcCountBefore := memBefore.NumGC

	peakAlloc := memBefore.HeapAlloc

	// Run build
	buildStart := time.Now()

	m, err := createBenchmarkManager(configPath, tempDir)
	if err != nil {
		result.ErrorMessage = fmt.Sprintf("failed to create manager: %v", err)
		return result
	}

	result.PluginCount = len(m.Plugins())

	// Run each stage and time it
	stages := []lifecycle.Stage{
		lifecycle.StageConfigure,
		lifecycle.StageValidate,
		lifecycle.StageGlob,
		lifecycle.StageLoad,
		lifecycle.StageTransform,
		lifecycle.StageRender,
		lifecycle.StageCollect,
		lifecycle.StageWrite,
		lifecycle.StageCleanup,
	}

	for _, stage := range stages {
		stageStart := time.Now()
		if err := m.RunTo(stage); err != nil {
			result.ErrorMessage = fmt.Sprintf("stage %s failed: %v", stage, err)
			return result
		}
		stageDuration := time.Since(stageStart)
		result.StageTimes[string(stage)] = stageDuration
		result.StageTimeSec[string(stage)] = stageDuration.Seconds()

		// Check peak memory after each stage
		var memStage runtime.MemStats
		runtime.ReadMemStats(&memStage)
		if memStage.HeapAlloc > peakAlloc {
			peakAlloc = memStage.HeapAlloc
		}

		if verbose {
			fmt.Printf("  [%s] %.3fs\n", stage, stageDuration.Seconds())
		}
	}

	result.TotalTime = time.Since(buildStart)
	result.TotalTimeSec = result.TotalTime.Seconds()

	// Collect final memory stats
	var memAfter runtime.MemStats
	runtime.ReadMemStats(&memAfter)

	result.MemoryStats = BenchmarkMemoryStats{
		PeakHeapMB:   float64(peakAlloc) / (1024 * 1024),
		AllocMB:      float64(memAfter.Alloc) / (1024 * 1024),
		TotalAllocMB: float64(memAfter.TotalAlloc-memBefore.TotalAlloc) / (1024 * 1024),
		GCCycles:     memAfter.NumGC - gcCountBefore,
	}

	// Calculate throughput
	if result.TotalTimeSec > 0 {
		result.PostsPerSec = float64(scenario.PostCount) / result.TotalTimeSec
	}

	result.Success = true
	return result
}

// createBenchmarkManager creates a manager for benchmarking without relying on global state.
func createBenchmarkManager(cfgPath, workDir string) (*lifecycle.Manager, error) {
	// Change to work directory temporarily
	originalDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("getting current dir: %w", err)
	}

	if err := os.Chdir(workDir); err != nil {
		return nil, fmt.Errorf("changing to work dir: %w", err)
	}
	defer func() {
		_ = os.Chdir(originalDir) //nolint:errcheck // best-effort restore of original directory
	}()

	// Load config
	cfg, err := config.Load(cfgPath)
	if err != nil {
		return nil, fmt.Errorf("loading config: %w", err)
	}

	// Create manager
	m := lifecycle.NewManager()

	// Convert models.Config to lifecycle.Config
	lcConfig := &lifecycle.Config{
		ContentDir:   ".",
		OutputDir:    cfg.OutputDir,
		GlobPatterns: cfg.GlobConfig.Patterns,
		Extra:        make(map[string]interface{}),
	}

	lcConfig.Extra["url"] = cfg.URL
	lcConfig.Extra["title"] = cfg.Title
	lcConfig.Extra["description"] = cfg.Description
	lcConfig.Extra["author"] = cfg.Author
	lcConfig.Extra["language"] = cfg.Language
	lcConfig.Extra["author_url"] = cfg.AuthorURL
	lcConfig.Extra["managing_editor"] = cfg.ManagingEditor
	lcConfig.Extra["webmaster"] = cfg.WebMaster
	lcConfig.Extra["copyright"] = cfg.Copyright
	lcConfig.Extra["templates_dir"] = cfg.TemplatesDir
	lcConfig.Extra["assets_dir"] = cfg.AssetsDir
	lcConfig.Extra["feeds"] = cfg.Feeds
	lcConfig.Extra["feed_defaults"] = cfg.FeedDefaults
	lcConfig.Extra["use_gitignore"] = cfg.GlobConfig.UseGitignore
	lcConfig.Extra["nav"] = cfg.Nav
	lcConfig.Extra["footer"] = cfg.Footer
	lcConfig.Extra["post_formats"] = cfg.PostFormats
	lcConfig.Extra["websub"] = cfg.WebSub
	lcConfig.Extra["well_known"] = cfg.WellKnown
	lcConfig.Extra["theme"] = map[string]interface{}{
		"name":       cfg.Theme.Name,
		"palette":    cfg.Theme.Palette,
		"variables":  cfg.Theme.Variables,
		"custom_css": cfg.Theme.CustomCSS,
	}

	// Pass layout configuration for automatic layout selection
	lcConfig.Extra["layout"] = &cfg.Layout

	// Pass blogroll configuration
	lcConfig.Extra["blogroll"] = cfg.Blogroll

	// Pass mentions configuration
	lcConfig.Extra["mentions"] = cfg.Mentions

	// Pass view transitions configuration
	lcConfig.Extra["view_transitions"] = cfg.ViewTransitions
	for key, value := range cfg.Extra {
		lcConfig.Extra[key] = value
	}
	// Build identity is supplied by the binary, not by benchmark site configuration.
	lcConfig.Extra["markata_version"] = Version
	lcConfig.Extra["markata_commit"] = Commit

	m.SetConfig(lcConfig)

	if cfg.Concurrency > 0 {
		m.SetConcurrency(cfg.Concurrency)
	}

	// Register default plugins
	m.RegisterPlugins(plugins.DefaultPlugins()...)

	return m, nil
}

// generateBenchmarkConfig creates a comprehensive configuration for benchmarking.
func generateBenchmarkConfig(tempDir string, scenario BenchmarkScenario) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf(`url = "https://benchmark.example.com"
title = "Benchmark Site - %s"
description = "A comprehensive benchmark site for markata-go performance testing"
author = "Benchmark Generator"
output_dir = "%s/output"

[glob]
patterns = ["posts/**/*.md"]

[feed_defaults]
items_per_page = 20
orphan_threshold = 5

[feed_defaults.formats]
html = true
rss = true
atom = false
json = false

`, scenario.Name, tempDir))

	// Generate feed configurations based on scenario
	feeds := generateFeedConfigs(scenario.FeedCount)
	for _, feed := range feeds {
		sb.WriteString(feed)
	}

	return sb.String()
}

// generateFeedConfigs creates feed configuration entries.
func generateFeedConfigs(count int) []string {
	feedConfigs := []struct {
		slug        string
		title       string
		description string
		filter      string
		sort        string
		reverse     bool
	}{
		{"blog", "All Posts", "All published blog posts", "published == true", "date", true},
		{"tutorials", "Tutorials", "Tutorial posts", "published == true and tags contains tutorial", "date", true},
		{"guides", "Guides", "Guide posts", "published == true and tags contains guide", "date", true},
		{"programming", "Programming", "Programming posts", "published == true and category == programming", "date", true},
		{"devops", "DevOps", "DevOps related posts", "published == true and tags contains devops", "date", true},
		{"featured", "Featured", "Featured content", "published == true and featured == true", "weight", false},
		{"recent", "Recent Posts", "Recently updated posts", "published == true", "updated", true},
		{"archive", "Archive", "Full archive", "published == true", "date", true},
	}

	result := make([]string, 0, count)
	for i := 0; i < count && i < len(feedConfigs); i++ {
		fc := feedConfigs[i]
		result = append(result, fmt.Sprintf(`
[[feeds]]
slug = %q
title = %q
description = %q
filter = %q
sort = %q
reverse = %v
`, fc.slug, fc.title, fc.description, fc.filter, fc.sort, fc.reverse))
	}

	return result
}

// generateRichTestPosts creates test markdown files with comprehensive content.
func generateRichTestPosts(dir string, count int) error {
	//nolint:gosec // Using math/rand is fine for benchmark content generation
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < count; i++ {
		title := generateRichTitle(rng, i)
		slug := benchmarkSlug(title)
		filename := filepath.Join(dir, fmt.Sprintf("%s.md", slug))

		content := generateRichPostMarkdown(rng, title, i, count)

		if err := os.WriteFile(filename, []byte(content), 0o600); err != nil {
			return fmt.Errorf("writing post %d: %w", i, err)
		}
	}

	return nil
}

// generateRichTitle creates a realistic post title.
func generateRichTitle(rng *rand.Rand, index int) string {
	adjectives := []string{
		"Comprehensive", "Ultimate", "Complete", "Essential", "Practical",
		"Modern", "Advanced", "Beginner's", "Expert", "In-Depth",
		"Definitive", "Professional", "Step-by-Step", "Quick", "Deep Dive",
	}
	topics := []string{
		"Guide to Go Programming", "Tutorial on Testing Strategies",
		"Introduction to REST APIs", "Overview of Database Design",
		"Walkthrough of Docker Containers", "Exploration of Kubernetes",
		"Analysis of Performance Optimization", "Review of Best Practices",
		"Comparison of Web Frameworks", "Journey into Microservices",
		"Mastering Concurrency Patterns", "Understanding Error Handling",
		"Building CLI Applications", "Implementing Authentication",
		"Working with JSON and YAML", "Creating Custom Middleware",
	}

	adj := adjectives[rng.Intn(len(adjectives))]
	topic := topics[rng.Intn(len(topics))]

	return fmt.Sprintf("%s %s Part %d", adj, topic, index+1)
}

// generateRichPostMarkdown creates a complete post with comprehensive frontmatter and rich content.
func generateRichPostMarkdown(rng *rand.Rand, title string, index, totalPosts int) string {
	tags := generateRichTags(rng)
	category := generateCategory(rng)
	date := generateRichDate(rng, index, totalPosts)
	updated := date.Add(time.Duration(rng.Intn(30*24)) * time.Hour)

	var sb strings.Builder

	// Comprehensive frontmatter
	sb.WriteString("---\n")
	sb.WriteString(fmt.Sprintf("title: %q\n", title))
	sb.WriteString(fmt.Sprintf("slug: %q\n", benchmarkSlug(title)))
	sb.WriteString(fmt.Sprintf("date: %s\n", date.Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("updated: %s\n", updated.Format(time.RFC3339)))
	sb.WriteString("published: true\n")
	sb.WriteString("draft: false\n")
	sb.WriteString("tags:\n")
	for _, tag := range tags {
		sb.WriteString(fmt.Sprintf("  - %s\n", tag))
	}
	sb.WriteString(fmt.Sprintf("category: %s\n", category))
	sb.WriteString(fmt.Sprintf("author: \"Author %d\"\n", (index%5)+1))
	sb.WriteString(fmt.Sprintf("description: \"A comprehensive benchmark test post covering %s\"\n", title))

	// Additional rich frontmatter fields
	if rng.Float32() < 0.3 {
		sb.WriteString("featured: true\n")
	}
	sb.WriteString(fmt.Sprintf("weight: %d\n", rng.Intn(100)))
	sb.WriteString(fmt.Sprintf("reading_time: %d\n", 5+rng.Intn(15)))

	// Custom metadata
	sb.WriteString("custom:\n")
	sb.WriteString(fmt.Sprintf("  difficulty: %s\n", []string{"beginner", "intermediate", "advanced"}[rng.Intn(3)]))
	sb.WriteString(fmt.Sprintf("  estimated_time: \"%d minutes\"\n", 10+rng.Intn(50)))
	sb.WriteString(fmt.Sprintf("  version: \"1.%d.%d\"\n", rng.Intn(10), rng.Intn(10)))

	sb.WriteString("---\n\n")

	// Rich body content
	sb.WriteString(generateRichBody(rng))

	return sb.String()
}

// generateRichTags creates a diverse set of tags.
func generateRichTags(rng *rand.Rand) []string {
	allTags := []string{
		"go", "golang", "programming", "tutorial", "guide",
		"testing", "api", "web", "backend", "devops",
		"docker", "kubernetes", "cloud", "performance", "best-practices",
		"security", "database", "microservices", "cli", "tools",
		"architecture", "patterns", "debugging", "deployment", "automation",
	}

	numTags := 3 + rng.Intn(5) // 3-7 tags
	tags := make([]string, 0, numTags)
	used := make(map[int]bool)

	for len(tags) < numTags {
		idx := rng.Intn(len(allTags))
		if !used[idx] {
			used[idx] = true
			tags = append(tags, allTags[idx])
		}
	}

	return tags
}

// generateCategory returns a category name.
func generateCategory(rng *rand.Rand) string {
	categories := []string{
		"programming", "devops", "architecture", "tutorials", "guides", "tools",
	}
	return categories[rng.Intn(len(categories))]
}

// generateRichDate creates a date spread across the past year.
func generateRichDate(rng *rand.Rand, index, totalPosts int) time.Time {
	now := time.Now()
	// Spread posts evenly across the past year
	daysRange := 365
	daysPerPost := float64(daysRange) / float64(totalPosts)
	baseDays := int(float64(index) * daysPerPost)
	jitter := rng.Intn(int(daysPerPost) + 1)
	daysAgo := baseDays + jitter
	return now.AddDate(0, 0, -daysAgo)
}

// generateRichBody creates comprehensive markdown content with various elements.
func generateRichBody(rng *rand.Rand) string {
	var sb strings.Builder

	// Introduction with emphasis
	sb.WriteString("## Introduction\n\n")
	sb.WriteString(generateRichParagraph(rng))
	sb.WriteString("\n\n")

	// Add a blockquote
	sb.WriteString("> ")
	sb.WriteString(generateQuote(rng))
	sb.WriteString("\n\n")

	// Prerequisites section with task list
	sb.WriteString("## Prerequisites\n\n")
	sb.WriteString("Before starting, ensure you have:\n\n")
	sb.WriteString(generateTaskList(rng))
	sb.WriteString("\n")

	// Main sections with varied content
	sections := []string{"Getting Started", "Core Concepts", "Implementation Details", "Advanced Topics", "Best Practices"}
	for i, section := range sections {
		sb.WriteString(fmt.Sprintf("## %s\n\n", section))
		sb.WriteString(generateRichParagraph(rng))
		sb.WriteString("\n\n")

		// Add code block with language variety
		if rng.Float32() < 0.7 {
			lang := []string{"go", "bash", "yaml", "json", "python"}[rng.Intn(5)]
			sb.WriteString(fmt.Sprintf("```%s\n", lang))
			sb.WriteString(generateRichCodeBlock(rng, lang))
			sb.WriteString("```\n\n")
		}

		// Add a table for some sections
		if i == 2 && rng.Float32() < 0.8 {
			sb.WriteString(generateTable(rng))
			sb.WriteString("\n")
		}

		// Add a list
		if rng.Float32() < 0.6 {
			if rng.Float32() < 0.5 {
				sb.WriteString(generateOrderedList(rng))
			} else {
				sb.WriteString(generateUnorderedList(rng))
			}
			sb.WriteString("\n")
		}

		// Add subsection
		if rng.Float32() < 0.4 {
			sb.WriteString(fmt.Sprintf("### %s Details\n\n", section))
			sb.WriteString(generateRichParagraph(rng))
			sb.WriteString("\n\n")
		}
	}

	// Summary section
	sb.WriteString("## Summary\n\n")
	sb.WriteString(generateRichParagraph(rng))
	sb.WriteString("\n\n")

	// Key takeaways
	sb.WriteString("### Key Takeaways\n\n")
	sb.WriteString(generateUnorderedList(rng))
	sb.WriteString("\n")

	// Conclusion
	sb.WriteString("## Conclusion\n\n")
	sb.WriteString(generateRichParagraph(rng))
	sb.WriteString("\n")

	return sb.String()
}

// generateRichParagraph creates a paragraph with inline formatting.
func generateRichParagraph(rng *rand.Rand) string {
	sentences := []string{
		"This is a **critical aspect** of modern software development.",
		"Understanding these concepts is _essential_ for building scalable applications.",
		"Many developers struggle with this topic initially, but **practice makes perfect**.",
		"With dedication and _continuous learning_, mastery is achievable.",
		"The key is to start with **fundamentals** and build from there.",
		"Best practices have _evolved significantly_ over the years.",
		"Performance considerations should **never** be overlooked.",
		"Testing is _crucial_ for maintaining high code quality.",
		"Documentation helps teams collaborate **effectively**.",
		"Continuous improvement leads to _better outcomes_ over time.",
		"The ecosystem provides many **powerful tools** to help with this.",
		"Community support is _invaluable_ when learning new technologies.",
		"Real-world experience often **differs** from tutorials.",
		"Edge cases require _careful consideration_ and proper handling.",
		"Security should be a **primary concern** from the start.",
		"This approach enables `clean code` that is easy to maintain.",
		"Using the right **patterns** can significantly improve your architecture.",
	}

	numSentences := 4 + rng.Intn(4) // 4-7 sentences
	result := make([]string, 0, numSentences)
	for i := 0; i < numSentences; i++ {
		result = append(result, sentences[rng.Intn(len(sentences))])
	}

	return strings.Join(result, " ")
}

// generateQuote returns an inspirational/technical quote.
func generateQuote(rng *rand.Rand) string {
	quotes := []string{
		"Simplicity is the ultimate sophistication. — Leonardo da Vinci",
		"First, solve the problem. Then, write the code. — John Johnson",
		"Code is like humor. When you have to explain it, it's bad. — Cory House",
		"The best error message is the one that never shows up. — Thomas Fuchs",
		"Make it work, make it right, make it fast. — Kent Beck",
	}
	return quotes[rng.Intn(len(quotes))]
}

// generateTaskList creates a markdown task list.
func generateTaskList(rng *rand.Rand) string {
	items := []struct {
		done bool
		text string
	}{
		{true, "Go 1.21 or later installed"},
		{true, "Basic understanding of Go syntax"},
		{true, "Familiarity with the command line"},
		{rng.Float32() < 0.7, "Docker installed (optional)"},
		{rng.Float32() < 0.5, "IDE or text editor configured"},
		{false, "Database setup (covered in this guide)"},
	}

	var sb strings.Builder
	numItems := 4 + rng.Intn(3)
	for i := 0; i < numItems && i < len(items); i++ {
		check := " "
		if items[i].done {
			check = "x"
		}
		sb.WriteString(fmt.Sprintf("- [%s] %s\n", check, items[i].text))
	}
	return sb.String()
}

// generateRichCodeBlock creates code snippets in various languages.
func generateRichCodeBlock(rng *rand.Rand, lang string) string {
	snippets := map[string][]string{
		"go": {
			`package main

import (
    "fmt"
    "log"
)

func main() {
    result, err := processData("input")
    if err != nil {
        log.Fatalf("Error: %v", err)
    }
    fmt.Printf("Result: %s\n", result)
}

func processData(input string) (string, error) {
    // Process the input data
    return fmt.Sprintf("Processed: %s", input), nil
}
`,
			`type Config struct {
    Name        string            ` + "`json:\"name\" yaml:\"name\"`" + `
    Port        int               ` + "`json:\"port\" yaml:\"port\"`" + `
    Debug       bool              ` + "`json:\"debug\" yaml:\"debug\"`" + `
    Features    []string          ` + "`json:\"features\" yaml:\"features\"`" + `
    Settings    map[string]string ` + "`json:\"settings\" yaml:\"settings\"`" + `
}

func NewConfig() *Config {
    return &Config{
        Name:  "default",
        Port:  8080,
        Debug: false,
    }
}
`,
			`func (s *Server) handleRequest(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    // Parse request body
    var req RequestBody
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    // Process the request
    result, err := s.service.Process(ctx, req)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Send response
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(result)
}
`,
		},
		"bash": {
			`#!/bin/bash
set -euo pipefail

# Build the application
echo "Building application..."
go build -o bin/app ./cmd/app

# Run tests
echo "Running tests..."
go test -v -race ./...

# Build Docker image
echo "Building Docker image..."
docker build -t myapp:latest .

echo "Done!"
`,
			`# Set environment variables
export APP_ENV=production
export LOG_LEVEL=info

# Start the application
./bin/app serve \
    --port 8080 \
    --config /etc/app/config.yaml \
    --workers 4
`,
		},
		"yaml": {
			`name: myapp
version: "1.0.0"

server:
  host: localhost
  port: 8080
  timeout: 30s

database:
  driver: postgres
  host: localhost
  port: 5432
  name: myapp_db

logging:
  level: info
  format: json
  output: stdout
`,
		},
		"json": {
			`{
  "name": "benchmark-app",
  "version": "1.0.0",
  "dependencies": {
    "express": "^4.18.0",
    "typescript": "^5.0.0"
  },
  "scripts": {
    "build": "tsc",
    "start": "node dist/index.js",
    "test": "jest"
  }
}
`,
		},
		"python": {
			`from typing import List, Optional
import asyncio

class DataProcessor:
    def __init__(self, config: dict):
        self.config = config
        self.results: List[str] = []

    async def process(self, items: List[str]) -> List[str]:
        """Process items concurrently."""
        tasks = [self._process_item(item) for item in items]
        return await asyncio.gather(*tasks)

    async def _process_item(self, item: str) -> str:
        await asyncio.sleep(0.1)  # Simulate work
        return f"Processed: {item}"
`,
		},
	}

	langSnippets, ok := snippets[lang]
	if !ok {
		langSnippets = snippets["go"]
	}
	return langSnippets[rng.Intn(len(langSnippets))]
}

// generateTable creates a markdown table.
func generateTable(rng *rand.Rand) string {
	tables := []string{
		`| Feature | Status | Notes |
|---------|--------|-------|
| Authentication | Implemented | OAuth2 + JWT |
| Authorization | Implemented | Role-based |
| Rate Limiting | In Progress | Redis-backed |
| Caching | Planned | Multi-tier |
| Monitoring | Implemented | Prometheus |
`,
		`| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/users | List all users |
| POST | /api/users | Create new user |
| GET | /api/users/:id | Get user by ID |
| PUT | /api/users/:id | Update user |
| DELETE | /api/users/:id | Delete user |
`,
		`| Metric | Value | Target |
|--------|-------|--------|
| Response Time | 45ms | <100ms |
| Throughput | 1000 req/s | >500 req/s |
| Error Rate | 0.1% | <1% |
| Availability | 99.9% | >99.5% |
`,
	}
	return tables[rng.Intn(len(tables))]
}

// generateOrderedList creates an ordered markdown list.
func generateOrderedList(rng *rand.Rand) string {
	items := []string{
		"Initialize the project with `go mod init`",
		"Define the data models and interfaces",
		"Implement the core business logic",
		"Add comprehensive error handling",
		"Write unit tests for all components",
		"Configure logging and monitoring",
		"Set up CI/CD pipeline",
		"Deploy to staging environment",
		"Run integration tests",
		"Deploy to production",
	}

	numItems := 4 + rng.Intn(4)
	var sb strings.Builder
	for i := 0; i < numItems && i < len(items); i++ {
		sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, items[i]))
	}
	return sb.String()
}

// generateUnorderedList creates an unordered markdown list.
func generateUnorderedList(rng *rand.Rand) string {
	items := []string{
		"Keep functions small and focused",
		"Use meaningful variable names",
		"Write tests before implementation",
		"Document public APIs thoroughly",
		"Handle errors explicitly",
		"Use interfaces for abstraction",
		"Prefer composition over inheritance",
		"Keep dependencies minimal",
		"Follow consistent formatting",
		"Review code before merging",
	}

	numItems := 3 + rng.Intn(4)
	var sb strings.Builder
	for i := 0; i < numItems; i++ {
		sb.WriteString(fmt.Sprintf("- %s\n", items[rng.Intn(len(items))]))
	}
	return sb.String()
}

// benchmarkSlug creates a URL-safe slug from a title for benchmark content.
func benchmarkSlug(title string) string {
	slug := strings.ToLower(title)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "'", "")
	slug = strings.ReplaceAll(slug, "\"", "")

	// Remove any non-alphanumeric characters except hyphens
	var result strings.Builder
	for _, r := range slug {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result.WriteRune(r)
		}
	}

	// Clean up multiple hyphens
	cleaned := result.String()
	for strings.Contains(cleaned, "--") {
		cleaned = strings.ReplaceAll(cleaned, "--", "-")
	}
	cleaned = strings.Trim(cleaned, "-")

	return cleaned
}

// outputJSONReport outputs the results as JSON.
func outputJSONReport(results []*BenchmarkResult) error {
	report := BenchmarkReport{
		Version:   Version,
		Timestamp: time.Now(),
		SystemInfo: BenchmarkSystem{
			OS:         runtime.GOOS,
			Arch:       runtime.GOARCH,
			NumCPU:     runtime.NumCPU(),
			GoVersion:  runtime.Version(),
			Gomaxprocs: runtime.GOMAXPROCS(0),
		},
		Results: results,
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(report)
}

// printBenchmarkSummary outputs the results in a formatted table.
func printBenchmarkSummary(results []*BenchmarkResult) {
	printResultsTable(results)
	fmt.Println()

	// Print stage breakdown if verbose
	if verbose {
		printStageBreakdown(results)
		printMemoryUsage(results)
	}

	// Performance insights
	printPerformanceAnalysis(results)

	// Notes
	printBenchmarkNotes()
}

// printResultsTable prints the main results table.
func printResultsTable(results []*BenchmarkResult) {
	fmt.Println("## Benchmark Results")
	fmt.Println()
	fmt.Println("| Scenario | Posts | Feeds | Build Time | Posts/sec | Peak Memory |")
	fmt.Println("|----------|-------|-------|------------|-----------|-------------|")

	for _, r := range results {
		if r.Success {
			fmt.Printf("| %s | %d | %d | %.2fs | %.1f | %.1fMB |\n",
				r.Scenario,
				r.PostCount,
				r.FeedCount,
				r.TotalTimeSec,
				r.PostsPerSec,
				r.MemoryStats.PeakHeapMB,
			)
		} else {
			fmt.Printf("| %s | %d | %d | FAILED | - | - |\n",
				r.Scenario, r.PostCount, r.FeedCount)
		}
	}
}

// printStageBreakdown prints the stage timing breakdown table.
func printStageBreakdown(results []*BenchmarkResult) {
	fmt.Println("### Stage Breakdown (seconds)")
	fmt.Println()

	stageOrder := []string{"configure", "validate", "glob", "load", "transform", "render", "collect", "write", "cleanup"}

	// Print header
	fmt.Print("| Stage |")
	for _, r := range results {
		fmt.Printf(" %s |", r.Scenario)
	}
	fmt.Println()

	// Print divider
	fmt.Print("|-------|")
	for range results {
		fmt.Print("--------|")
	}
	fmt.Println()

	// Print stage times
	for _, stage := range stageOrder {
		fmt.Printf("| %s |", stage)
		for _, r := range results {
			if r.Success {
				if t, ok := r.StageTimeSec[stage]; ok {
					fmt.Printf(" %.3f |", t)
				} else {
					fmt.Printf(" - |")
				}
			} else {
				fmt.Printf(" - |")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

// printMemoryUsage prints the memory usage table.
func printMemoryUsage(results []*BenchmarkResult) {
	fmt.Println("### Memory Usage")
	fmt.Println()
	fmt.Println("| Scenario | Peak Heap | Total Alloc | GC Cycles |")
	fmt.Println("|----------|-----------|-------------|-----------|")
	for _, r := range results {
		if r.Success {
			fmt.Printf("| %s | %.1fMB | %.1fMB | %d |\n",
				r.Scenario,
				r.MemoryStats.PeakHeapMB,
				r.MemoryStats.TotalAllocMB,
				r.MemoryStats.GCCycles,
			)
		}
	}
	fmt.Println()
}

// printPerformanceAnalysis prints performance insights.
func printPerformanceAnalysis(results []*BenchmarkResult) {
	fmt.Println("### Performance Analysis")
	fmt.Println()

	printSlowestStages(results)
	printScalingAnalysis(results)
}

// printSlowestStages prints the slowest stages across scenarios.
func printSlowestStages(results []*BenchmarkResult) {
	stageAvgTimes := make(map[string]float64)
	stageCount := make(map[string]int)
	for _, r := range results {
		if r.Success {
			for stage, t := range r.StageTimeSec {
				stageAvgTimes[stage] += t
				stageCount[stage]++
			}
		}
	}

	type stageStat struct {
		name    string
		avgTime float64
	}
	var stats []stageStat
	for stage, total := range stageAvgTimes {
		if count := stageCount[stage]; count > 0 {
			stats = append(stats, stageStat{stage, total / float64(count)})
		}
	}
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].avgTime > stats[j].avgTime
	})

	if len(stats) > 0 {
		fmt.Println("**Slowest stages (avg across scenarios):**")
		for i := 0; i < 3 && i < len(stats); i++ {
			fmt.Printf("- %s: %.3fs\n", stats[i].name, stats[i].avgTime)
		}
		fmt.Println()
	}
}

// printScalingAnalysis prints scaling analysis for multiple scenarios.
func printScalingAnalysis(results []*BenchmarkResult) {
	if len(results) < 2 {
		return
	}

	var successResults []*BenchmarkResult
	for _, r := range results {
		if r.Success {
			successResults = append(successResults, r)
		}
	}

	if len(successResults) < 2 {
		return
	}

	first := successResults[0]
	last := successResults[len(successResults)-1]
	postScale := float64(last.PostCount) / float64(first.PostCount)
	timeScale := last.TotalTimeSec / first.TotalTimeSec

	fmt.Printf("**Scaling:** %.1fx posts resulted in %.1fx build time (%.2f scaling factor)\n",
		postScale, timeScale, timeScale/postScale)
	fmt.Println()
}

// printBenchmarkNotes prints benchmark notes.
func printBenchmarkNotes() {
	fmt.Println("**Notes:**")
	fmt.Println("- Peak memory shows maximum heap allocation during build")
	fmt.Println("- Results may vary based on system load and hardware")
	fmt.Println("- Rich content includes comprehensive frontmatter, code blocks, tables, and lists")
	if benchmarkKeep {
		fmt.Println("- Test files were kept in temp directories")
	}
}
