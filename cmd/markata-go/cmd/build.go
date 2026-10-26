package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/WaylonWalker/markata/pkg/buildstats"
	"github.com/WaylonWalker/markata/pkg/diagnostics"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/pageweight"
	"github.com/spf13/cobra"
)

// Common string constants.
const (
	defaultOutputDir = "output"
)

var (
	// buildClean removes output directory and build cache before building.
	buildClean bool

	// buildCleanAll removes output, build cache, AND external plugin caches before building.
	buildCleanAll bool

	// buildDryRun shows what would be built without building.
	buildDryRun bool

	// buildFast skips expensive non-essential plugins (minification, CSS purging,
	// Tailwind rebuilds, and Pagefind indexing)
	// for faster development iteration.
	buildFast bool

	// buildBenchmarkJSON writes benchmark details as JSON. Use "-" for stdout.
	buildBenchmarkJSON string

	// buildBenchmarkSummaryJSON writes benchmark aggregates without content entries.
	buildBenchmarkSummaryJSON string

	// buildBenchmarkDetailed prints per-stage benchmark detail.
	buildBenchmarkDetailed bool

	// Build-time analysis of generated HTML (not browser transfer or Lighthouse).
	buildPageWeight       = true
	buildPageWeightHTMLKB = 160
)

// buildCmd represents the build command.
var buildCmd = &cobra.Command{
	Use:     "build [markdown-file]",
	Aliases: []string{"b"},
	Short:   "Build the static site",
	Long: `Build runs all lifecycle stages to generate the static site.

The build process includes:
  1. Configure - Load and validate configuration
  2. Glob - Discover content files
  3. Load - Parse markdown files and frontmatter
  4. Transform - Process Jinja templates in markdown
  5. Render - Convert markdown to HTML
  6. Collect - Build feeds and archives
  7. Write - Output files to disk

Clean modes:
  --clean      Remove output directory and build cache (.markata/).
               This is cheap to rebuild and fixes most caching issues.

  --clean-all  Everything --clean does, plus remove external plugin caches
               (blogroll feeds, embeds metadata, mentions, webmentions).
               These are expensive to re-fetch from remote servers.

	Fast mode:
	  --fast       Skip minification (JS/CSS), CSS purging, Tailwind rebuilds,
	               and Pagefind indexing for faster builds.
	               Useful during development iteration when you don't need
	               optimized output.

Example usage:
  markata-go build              # Standard build
  markata-go build pages/post.md # Render one Markdown file to <output>/index.html
  markata-go build --clean      # Clean build cache + output
  markata-go build --clean-all  # Also nuke external plugin caches
  markata-go build --fast       # Skip minification for faster builds
  markata-go build --dry-run    # Show what would be built
  markata-go build -v           # Build with verbose output`,
	Args: cobra.MaximumNArgs(1),
	RunE: runBuildCommand,
}

func init() {
	rootCmd.AddCommand(buildCmd)

	buildCmd.Flags().BoolVar(&buildClean, "clean", false, "clean output directory and build cache before build")
	buildCmd.Flags().BoolVar(&buildCleanAll, "clean-all", false, "clean everything including external plugin caches (blogroll, embeds, etc.)")
	buildCmd.Flags().BoolVar(&buildDryRun, "dry-run", false, "show what would be built without building")
	buildCmd.Flags().BoolVar(&buildFast, "fast", false, "skip minification, CSS purging, tailwind rebuilds, and pagefind indexing for faster builds")
	buildCmd.Flags().StringVar(&buildBenchmarkJSON, "benchmark-json", "", "write benchmark details as JSON (use '-' for stdout)")
	buildCmd.Flags().Lookup("benchmark-json").NoOptDefVal = "-"
	buildCmd.Flags().StringVar(&buildBenchmarkSummaryJSON, "benchmark-summary-json", "", "write compact benchmark JSON without content entries (use '-' for stdout)")
	buildCmd.Flags().Lookup("benchmark-summary-json").NoOptDefVal = "-"
	buildCmd.Flags().BoolVar(&buildBenchmarkDetailed, "benchmark-detailed", false, "print per-stage benchmark resource summaries")
	buildCmd.Flags().BoolVar(&buildPageWeight, "page-weight", true, "write per-page HTML weight telemetry and emit budget warnings")
	buildCmd.Flags().IntVar(&buildPageWeightHTMLKB, "page-weight-html-kb", 160, "warn when a post exceeds this many KiB of generated HTML")
}

func runBuildCommand(cmd *cobra.Command, args []string) error {
	if err := validateDAGEnvironment(cmd); err != nil {
		return err
	}
	if err := validateBuildOptions(); err != nil {
		return err
	}
	startTime := time.Now()

	verbosef("Starting build...")

	// Create the manager
	var (
		m   *lifecycle.Manager
		err error
	)
	if len(args) == 1 {
		m, err = createSinglePageManager(cfgFile, args[0])
	} else {
		m, err = createManager(cfgFile)
	}
	if err != nil {
		return fmt.Errorf("initialization failed: %w", err)
	}
	configureLoggerForManager(m)

	// Pass fast mode flag to plugins via config
	if buildFast {
		applyFastMode(m)
	}

	verbosef("Configuration loaded (output: %s, patterns: %v)", m.Config().OutputDir, m.Config().GlobPatterns)

	// Clean directories if requested
	if buildCleanAll || buildClean {
		if err := cleanBuildDirs(m); err != nil {
			return err
		}
	}

	// Dry run - just show what would be processed
	if buildDryRun {
		return runDryBuild(m)
	}

	// Run the build
	result, err := runBuild(m)
	if err != nil {
		return fmt.Errorf("build failed: %w", err)
	}

	// Calculate duration
	duration := time.Since(startTime)
	result.Duration = duration.Seconds()

	if err := writeBuildReport(result); err != nil {
		return err
	}
	if buildPageWeight {
		if err := reportGeneratedPageWeight(m.Config().OutputDir); err != nil {
			return err
		}
	}

	// Print warnings
	if len(result.Warnings) > 0 && verbose {
		errln("\nWarnings:")
		for _, w := range result.Warnings {
			errlnf("  - %s", w)
		}
	}

	return nil
}

func validateBuildOptions() error {
	if buildBenchmarkJSON != "" && buildBenchmarkSummaryJSON != "" {
		return fmt.Errorf("--benchmark-json and --benchmark-summary-json cannot be used together")
	}
	if buildPageWeight && (buildPageWeightHTMLKB < 1 || buildPageWeightHTMLKB > 10240) {
		return fmt.Errorf("--page-weight-html-kb must be between 1 and 10240")
	}
	return nil
}

func reportGeneratedPageWeight(outputDir string) error {
	if outputDir == "" {
		outputDir = defaultOutputDir
	}
	limits := pageweight.DefaultLimits()
	limits.PostHTMLBytes = buildPageWeightHTMLKB * 1024
	report, err := pageweight.Scan(outputDir, limits)
	if err != nil {
		return fmt.Errorf("page weight diagnostics: %w", err)
	}
	reportPath := filepath.Join(outputDir, ".markata", "page-weight.json")
	if err := pageweight.WriteJSON(reportPath, report); err != nil {
		return fmt.Errorf("write page weight diagnostics: %w", err)
	}
	errlnf("Page weight: %d HTML pages, %d over budget; report: %s",
		report.Summary.Pages, report.Summary.WarnedPages, reportPath)
	shown := 0
	for i := range report.Pages {
		page := &report.Pages[i]
		for _, finding := range page.Findings {
			if shown >= 5 {
				break
			}
			errlnf("  [page-weight:%s] %s: %d B (budget %d B) — %s",
				finding.Code, page.Path, finding.Bytes, finding.BudgetBytes, finding.Advice)
			shown++
		}
		if shown >= 5 {
			break
		}
	}
	if report.Summary.Findings > shown {
		errlnf("  ... %d more findings in %s", report.Summary.Findings-shown, reportPath)
	}
	return nil
}

// writeBuildReport writes the selected benchmark report or human-readable results.
func writeBuildReport(result *BuildResult) error {
	switch {
	case buildBenchmarkJSON == "-":
		if err := writeBenchmarkJSON(outWriter(), result); err != nil {
			return fmt.Errorf("writing benchmark json: %w", err)
		}
	case buildBenchmarkSummaryJSON == "-":
		if err := writeBenchmarkSummaryJSON(outWriter(), result); err != nil {
			return fmt.Errorf("writing benchmark summary json: %w", err)
		}
	default:
		// Print results
		printBuildResult(result)
	}

	if buildBenchmarkJSON != "" && buildBenchmarkJSON != "-" {
		if err := writeBenchmarkJSONFile(buildBenchmarkJSON, result); err != nil {
			return fmt.Errorf("writing benchmark json file: %w", err)
		}
	}
	if buildBenchmarkSummaryJSON != "" && buildBenchmarkSummaryJSON != "-" {
		if err := writeBenchmarkSummaryJSONFile(buildBenchmarkSummaryJSON, result); err != nil {
			return fmt.Errorf("writing benchmark summary json file: %w", err)
		}
	}

	return nil
}

// cleanBuildDirs removes build artifacts before a fresh build.
// --clean removes output dir and .markata/ build cache (Tier 1).
// --clean-all additionally removes external plugin caches (Tier 2).
func cleanBuildDirs(m *lifecycle.Manager) error {
	outputPath := m.Config().OutputDir
	if outputPath == "" {
		outputPath = defaultOutputDir
	}

	// Build cache directory (.markata/) is a sibling of the output directory.
	// Without cleaning it, the glob cache retains stale file lists.
	cacheDir := filepath.Join(outputPath, "..", ".markata")
	if extra := m.Config().Extra; extra != nil {
		if dir, ok := extra["cache_dir"].(string); ok && dir != "" {
			cacheDir = dir
		}
	}

	if verbose {
		verbosef("Cleaning output directory: %s", outputPath)
		verbosef("Cleaning cache directory: %s", cacheDir)
	}

	if !buildDryRun {
		if err := os.RemoveAll(outputPath); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to clean output directory: %w", err)
		}
		if err := os.RemoveAll(cacheDir); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to clean cache directory: %w", err)
		}
	}

	// --clean-all: also remove external plugin caches (Tier 2)
	if buildCleanAll {
		if modelsConfig, ok := m.Config().Extra["models_config"].(*models.Config); ok {
			for _, dir := range modelsConfig.ExternalCacheDirs() {
				if verbose {
					verbosef("Cleaning external cache: %s", dir)
				}
				if !buildDryRun {
					if err := os.RemoveAll(dir); err != nil && !os.IsNotExist(err) {
						return fmt.Errorf("failed to clean external cache %s: %w", dir, err)
					}
				}
			}
		}
	}

	return nil
}

// runDryBuild performs a dry run build showing what would be processed.
func runDryBuild(m *lifecycle.Manager) error {
	errln("Dry run mode - no files will be written")

	// Run stages up to Write (but not Write)
	stages := []lifecycle.Stage{
		lifecycle.StageConfigure,
		lifecycle.StageValidate,
		lifecycle.StageGlob,
		lifecycle.StageLoad,
		lifecycle.StageTransform,
		lifecycle.StageRender,
		lifecycle.StageCollect,
	}

	for _, stage := range stages {
		verbosef("Running stage: %s", stage)
		if err := m.RunTo(stage); err != nil {
			return fmt.Errorf("stage %s failed: %w", stage, err)
		}
	}

	// Print what would be written
	outlnf("Files discovered: %d", len(m.Files()))
	outlnf("Posts to process: %d", len(m.Posts()))
	outlnf("Feeds to generate: %d", len(m.Feeds()))
	printContentSummary(m.ContentDiagnostics())

	if verbose {
		errln("\nFiles that would be processed:")
		for _, f := range m.Files() {
			errlnf("  - %s", f)
		}

		errln("\nPosts that would be generated:")
		for _, p := range m.Posts() {
			title := p.Slug
			if p.PlainTitle() != "" {
				title = p.PlainTitle()
			}
			status := "published"
			if p.Draft {
				status = "draft"
			} else if !p.Published {
				status = "unpublished"
			}
			errlnf("  - %s (%s) [%s]", title, p.Slug, status)
		}

		errln("\nFeeds that would be generated:")
		for _, f := range m.Feeds() {
			errlnf("  - %s (%d posts)", f.Name, len(f.Posts))
		}
	}

	outputDir := m.Config().OutputDir
	if outputDir == "" {
		outputDir = defaultOutputDir
	}
	outlnf("\nOutput directory: %s", filepath.Clean(outputDir))

	return nil
}

// printBuildResult prints a summary of the build result.
func printBuildResult(result *BuildResult) {
	outln("\n" + colorizeOutput("Build completed successfully!", currentLogTheme.Component))
	outlnf("  %s %d", buildLabel("Posts processed:"), result.PostsProcessed)

	// Only show feeds if any were generated
	if result.FeedsGenerated > 0 {
		outlnf("  %s %d", buildLabel("Feeds generated:"), result.FeedsGenerated)
	}

	// Show blogroll status if configured
	printBlogrollStatus(result.BlogrollStatus)
	printContentSummary(result.Content)
	printBuildBenchmarkSummary(result.Benchmark)

	outlnf("  %s %.2fs", buildLabel("Duration:"), result.Duration)
	if result.Executor == lifecycle.BuildExecutorDAG {
		outlnf("  %s %s", buildLabel("Executor:"), "serial DAG (experimental)")
	}
}

func printContentSummary(snapshot diagnostics.ContentLedgerSnapshot) {
	summary := snapshot.Summary
	outln("  Content:")
	outlnf("    Discovered: %d  Candidates: %d  Loaded: %d", summary.Discovered, summary.Candidates, summary.Loaded)
	outlnf("    Frontmatter valid: %d  Posts: %d  Eligible: %d", summary.FrontmatterValid, summary.Posts, summary.Eligible)
	outlnf("    Rendered: %d  Emitted: %d  Excluded: %d", summary.Rendered, summary.Emitted, summary.Excluded)
	outlnf("    Warnings: %d  Errors: %d", summary.Warnings, summary.Errors)

	for _, entry := range snapshot.Entries {
		for _, issue := range entry.Diagnostics {
			if issue.Severity != diagnostics.SeverityWarning && issue.Severity != diagnostics.SeverityError {
				continue
			}
			errlnf("Content %s: %s:%d %s: %s", issue.Severity.String(), issue.File, issue.Range.StartLine+1, issue.Code, issue.Message)
		}
	}
	if !verbose {
		return
	}

	errln("\nContent dispositions:")
	for _, entry := range snapshot.Entries {
		if entry.Disposition != diagnostics.DispositionExcluded && entry.Disposition != diagnostics.DispositionShadow {
			continue
		}
		reasons := strings.Join(entry.Reasons, ", ")
		if reasons == "" {
			reasons = "none"
		}
		errlnf("  %s [%s]: %s", entry.Path, entry.Disposition, reasons)
	}
}

func printBuildBenchmarkSummary(summary buildstats.Summary) {
	if summary.Total <= 0 {
		return
	}

	outlnf("  %s %s", buildLabel("Resource profile:"), colorizeOutput("estimated wall time", currentLogTheme.Component))
	printResourceLine("CPU", summary.Resources.CPU, summary.Total)
	printResourceLine("Network wait", summary.Resources.NetworkWait, summary.Total)
	printResourceLine("Disk read", summary.Resources.DiskReadWait, summary.Total)
	printResourceLine("Disk write", summary.Resources.DiskWriteWait, summary.Total)
	printResourceLine("Idle", summary.Resources.Idle, summary.Total)

	if verbose || buildBenchmarkDetailed {
		printStageBenchmarkSummary(summary)
	}

	hotspots := topHotspots(summary.Hotspots, 5)
	if len(hotspots) == 0 {
		printNetworkRequestSummary(summary.Requests)
		return
	}

	outlnf("  %s", buildLabel("Hotspots:"))
	for _, hotspot := range hotspots {
		outlnf("    %s %s", hotspotKey(hotspot.Stage, hotspot.Plugin), formatDuration(hotspot.Duration))
	}

	printNetworkRequestSummary(summary.Requests)
}

func printResourceLine(label string, duration, total time.Duration) {
	outlnf("    %-12s %8s (%5.1f%%)", colorizeOutput(label, currentLogTheme.Warning), formatDuration(duration), percent(duration, total))
}

func printStageBenchmarkSummary(summary buildstats.Summary) {
	stages := nonZeroStages(summary.Stages)
	if len(stages) == 0 {
		return
	}

	outlnf("  %s %s", buildLabel("Stage resources:"), colorizeOutput("estimated wall time", currentLogTheme.Component))
	for _, stage := range stages {
		outlnf("    %s total=%-8s cpu=%-8s net=%-8s read=%-8s write=%-8s idle=%-8s",
			colorizeOutput(stage.Stage, stageThemeColor(stage.Stage)),
			formatDuration(stage.Duration),
			formatDuration(stage.Resources.CPU),
			formatDuration(stage.Resources.NetworkWait),
			formatDuration(stage.Resources.DiskReadWait),
			formatDuration(stage.Resources.DiskWriteWait),
			formatDuration(stage.Resources.Idle),
		)
	}
}

func nonZeroStages(stages []buildstats.StageTiming) []buildstats.StageTiming {
	filtered := make([]buildstats.StageTiming, 0, len(stages))
	for _, stage := range stages {
		if stage.Duration <= 0 {
			continue
		}
		filtered = append(filtered, stage)
	}
	return filtered
}

func topHotspots(hotspots []buildstats.Hotspot, limit int) []buildstats.Hotspot {
	if len(hotspots) == 0 || limit <= 0 {
		return nil
	}
	if len(hotspots) < limit {
		limit = len(hotspots)
	}
	return hotspots[:limit]
}

func topRequests(requests []buildstats.RequestTiming, limit int) []buildstats.RequestTiming {
	if len(requests) == 0 || limit <= 0 {
		return nil
	}
	if len(requests) < limit {
		limit = len(requests)
	}
	return requests[:limit]
}

func formatDuration(duration time.Duration) string {
	text := duration.Round(10 * time.Millisecond).String()
	if text != "0s" && strings.HasSuffix(text, "0s") {
		return strings.TrimSuffix(text, "0s") + "s"
	}
	return text
}

func percent(duration, total time.Duration) float64 {
	if total <= 0 {
		return 0
	}
	return float64(duration) / float64(total) * 100
}

func hotspotKey(stage, plugin string) string {
	return colorizeOutput(stage+"/"+plugin, stageThemeColor(stage))
}

func printNetworkRequestSummary(requests []buildstats.RequestTiming) {
	requests = topRequests(requests, 10)
	if len(requests) == 0 {
		return
	}

	outlnf("  %s", buildLabel("Slowest requests:"))
	for _, request := range requests {
		outlnf("    %s %s %s", requestKey(request), formatDuration(request.Duration), requestStatus(request))
	}
}

func requestKey(request buildstats.RequestTiming) string {
	plugin := hotspotKey(request.Stage, request.Plugin)
	if request.Plugin == "" {
		plugin = colorizeOutput(request.Stage, stageThemeColor(request.Stage))
	}
	return plugin + " " + request.Method + " " + request.URL
}

func requestStatus(request buildstats.RequestTiming) string {
	if request.Error != "" {
		return "(" + request.Error + ")"
	}
	if request.Status > 0 {
		return fmt.Sprintf("(HTTP %d)", request.Status)
	}
	return ""
}

func buildLabel(name string) string {
	return colorizeOutput(name, currentLogTheme.Warning)
}

func stageThemeColor(stage string) string {
	if stageColor, ok := currentLogTheme.PhaseColor[stage]; ok && stageColor != "" {
		return stageColor
	}
	return currentLogTheme.Component
}

type benchmarkJSONOutput struct {
	Executor       lifecycle.BuildExecutor           `json:"-"`
	PostsProcessed int                               `json:"posts_processed"`
	FeedsGenerated int                               `json:"feeds_generated"`
	Duration       float64                           `json:"duration_seconds"`
	Warnings       []string                          `json:"warnings,omitempty"`
	Benchmark      buildstats.Summary                `json:"benchmark"`
	Blogroll       BlogrollStatus                    `json:"blogroll"`
	Content        diagnostics.ContentLedgerSnapshot `json:"content"`
}

type benchmarkSummaryContent struct {
	Summary       diagnostics.ContentSummary      `json:"summary"`
	TemplateCache *diagnostics.TemplateCacheStats `json:"template_cache,omitempty"`
}

type benchmarkSummaryJSONOutput struct {
	Executor       lifecycle.BuildExecutor `json:"executor"`
	PostsProcessed int                     `json:"posts_processed"`
	FeedsGenerated int                     `json:"feeds_generated"`
	Duration       float64                 `json:"duration_seconds"`
	Warnings       []string                `json:"warnings,omitempty"`
	Benchmark      buildstats.Summary      `json:"benchmark"`
	Blogroll       BlogrollStatus          `json:"blogroll"`
	Content        benchmarkSummaryContent `json:"content"`
}

func writeBenchmarkJSONFile(path string, result *BuildResult) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil && filepath.Dir(path) != "." {
		return err
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	return writeBenchmarkJSONAndClose(file, result)
}

func writeBenchmarkSummaryJSONFile(path string, result *BuildResult) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil && filepath.Dir(path) != "." {
		return err
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	return writeBenchmarkSummaryJSONAndClose(file, result)
}

func writeBenchmarkSummaryJSONAndClose(w io.WriteCloser, result *BuildResult) (err error) {
	defer func() {
		if closeErr := w.Close(); err == nil && closeErr != nil {
			err = fmt.Errorf("close benchmark summary JSON: %w", closeErr)
		}
	}()
	return writeBenchmarkSummaryJSON(w, result)
}

func writeBenchmarkSummaryJSON(w io.Writer, result *BuildResult) error {
	executor := result.Executor
	if executor == "" {
		executor = lifecycle.BuildExecutorLegacy
	}
	if !executor.Valid() {
		return fmt.Errorf("unsupported benchmark executor %q", executor)
	}
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(benchmarkSummaryJSONOutput{
		Executor:       executor,
		PostsProcessed: result.PostsProcessed,
		FeedsGenerated: result.FeedsGenerated,
		Duration:       result.Duration,
		Warnings:       result.Warnings,
		Benchmark:      result.Benchmark,
		Blogroll:       result.BlogrollStatus,
		Content: benchmarkSummaryContent{
			Summary:       result.Content.Summary,
			TemplateCache: result.Content.TemplateCache,
		},
	})
}

func writeBenchmarkJSONAndClose(w io.WriteCloser, result *BuildResult) (err error) {
	defer func() {
		if closeErr := w.Close(); err == nil && closeErr != nil {
			err = fmt.Errorf("close benchmark JSON: %w", closeErr)
		}
	}()
	return writeBenchmarkJSON(w, result)
}

func writeBenchmarkJSON(w io.Writer, result *BuildResult) error {
	// Encode metadata using the existing custom marshaler, retaining executor
	// normalization/validation and field order without encoding the entry array.
	content := result.Content
	if content.Entries != nil {
		content.Entries = content.Entries[:0]
	}
	header, err := json.MarshalIndent(benchmarkJSONOutput{
		Executor:       result.Executor,
		PostsProcessed: result.PostsProcessed,
		FeedsGenerated: result.FeedsGenerated,
		Duration:       result.Duration,
		Warnings:       result.Warnings,
		Benchmark:      result.Benchmark,
		Blogroll:       result.BlogrollStatus,
		Content:        content,
	}, "", "  ")
	if err != nil {
		return fmt.Errorf("encode benchmark JSON header: %w", err)
	}

	writer := bufio.NewWriter(w)
	if err := streamBenchmarkJSON(writer, header, result.Content.Entries); err != nil {
		return err
	}
	if err := writer.Flush(); err != nil {
		return fmt.Errorf("flush benchmark JSON: %w", err)
	}
	return nil
}

func streamBenchmarkJSON(w io.Writer, header []byte, entries []diagnostics.ContentDisposition) error {
	if len(entries) == 0 {
		if err := writeBenchmarkJSONBytes(w, header); err != nil {
			return err
		}
		return writeBenchmarkJSONBytes(w, []byte("\n"))
	}

	// Content and its entries are the final fields. Check the layout so a future
	// field addition cannot silently corrupt the streamed document.
	const suffix = "    \"entries\": []\n  }\n}"
	if !bytes.HasSuffix(header, []byte(suffix)) {
		return fmt.Errorf("benchmark JSON header must end with empty content entries")
	}
	if err := writeBenchmarkJSONBytes(w, header[:len(header)-len("]\n  }\n}")]); err != nil {
		return err
	}

	// Match the artifact writer's reusable per-entry buffers and json.Indent
	// layout, but preserve the raw snapshot's order and values without cloning.
	var compact, indented bytes.Buffer
	encoder := json.NewEncoder(&compact)
	for index := range entries {
		compact.Reset()
		if err := encoder.Encode(&entries[index]); err != nil {
			return fmt.Errorf("encode benchmark JSON entry %d: %w", index, err)
		}
		indented.Reset()
		if err := json.Indent(&indented, bytes.TrimSuffix(compact.Bytes(), []byte("\n")), "      ", "  "); err != nil {
			return fmt.Errorf("indent benchmark JSON entry %d: %w", index, err)
		}
		separator := "\n      "
		if index > 0 {
			separator = ",\n      "
		}
		if err := writeBenchmarkJSONBytes(w, []byte(separator)); err != nil {
			return err
		}
		if err := writeBenchmarkJSONBytes(w, indented.Bytes()); err != nil {
			return err
		}
	}
	return writeBenchmarkJSONBytes(w, []byte("\n    ]\n  }\n}\n"))
}

func writeBenchmarkJSONBytes(w io.Writer, data []byte) error {
	n, err := w.Write(data)
	if err != nil {
		return fmt.Errorf("write benchmark JSON: %w", err)
	}
	if n != len(data) {
		return fmt.Errorf("write benchmark JSON: %w", io.ErrShortWrite)
	}
	return nil
}

// printBlogrollStatus prints the blogroll feature status.
func printBlogrollStatus(status BlogrollStatus) {
	if !status.Configured {
		return
	}

	if status.Enabled {
		// Active blogroll - show pages and feed count
		outlnf("  %s /blogroll, /reader (%d %s)", buildLabel("Blogroll:"), status.FeedsFetched, pluralize(status.FeedsFetched, "feed", "feeds"))
	} else if status.FeedsConfigured > 0 {
		// Configured but disabled - show warning
		warnf("Blogroll: feeds configured but enabled=false")
	}
}

// pluralize returns singular or plural form based on count.
func pluralize(count int, singular, plural string) string {
	if count == 1 {
		return singular
	}
	return plural
}
