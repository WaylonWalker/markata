package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/WaylonWalker/markata/pkg/config"
	"github.com/WaylonWalker/markata/pkg/lint"
	"github.com/bmatcuk/doublestar/v4"
	"github.com/spf13/cobra"
)

var (
	// lintFix automatically fixes issues when possible.
	lintFix bool

	// lintDryRun shows which files would be checked without actually linting them.
	lintDryRun bool
)

// lintCmd represents the lint command.
var lintCmd = &cobra.Command{
	Use:   "lint [files...]",
	Short: "Lint markdown files for common issues",
	Long: `Lint checks markdown files for common issues that can cause build failures.

The linter detects:
  - Duplicate YAML keys in frontmatter
  - Invalid date formats (non-ISO 8601)
  - Malformed image links (missing alt text)
  - Protocol-less URLs (should use https://)

When run without arguments, lints all files matching the configured glob patterns
(defaults to **/*.md). Explicit file arguments override config patterns.

Use --fix to automatically fix detected issues.
Use --dry-run to see which files would be checked without actually linting them.

Example usage:
  markata-go lint                        # Lint all configured input files
  markata-go lint --dry-run              # Show which files would be checked
  markata-go lint posts/**/*.md          # Lint specific pattern (overrides config)
  markata-go lint posts/**/*.md --fix    # Lint and auto-fix issues
  markata-go lint pages/about.md         # Lint a specific file`,
	RunE: runLintCommand,
}

func init() {
	rootCmd.AddCommand(lintCmd)

	lintCmd.Flags().BoolVar(&lintFix, "fix", false, "automatically fix issues")
	lintCmd.Flags().BoolVar(&lintDryRun, "dry-run", false, "show which files would be checked without linting")
}

// lintStats tracks linting statistics.
type lintStats struct {
	totalFiles      int
	totalIssues     int
	fixableIssues   int
	totalFixed      int
	filesWithIssues int
	hasErrors       bool
}

func runLintCommand(cmd *cobra.Command, args []string) error {
	currentCmd = cmd
	var files []string
	var err error

	if len(args) > 0 {
		// Explicit file arguments override config patterns
		files, err = expandGlobPatterns(args)
		if err != nil {
			return err
		}
	} else {
		// No arguments: use configured glob patterns
		files, err = getFilesFromConfig()
		if err != nil {
			return err
		}
	}

	if len(files) == 0 {
		return fmt.Errorf("no files to lint")
	}

	// Sort files for consistent output
	sort.Strings(files)

	// Dry run: just show which files would be checked
	if lintDryRun {
		outlnf("Would lint %d file(s):", len(files))
		for _, f := range files {
			outlnf("  %s", f)
		}
		return nil
	}

	// Match against the configured blogroll, without making network requests.
	// A disabled blogroll preserves the CLI's historical lint behavior.
	var mentionResolver *blogrollMentionResolver
	cfg, cfgErr := config.Load(cfgFile)
	if cfgErr != nil {
		return fmt.Errorf("loading config for mention lint: %w", cfgErr)
	}
	mentionResolver = newBlogrollMentionResolver(cfg.Blogroll)

	stats := &lintStats{}
	for _, file := range files {
		processFile(file, stats, mentionResolver)
	}
	processEncryptionPolicyLint(stats)

	printSummary(stats)

	// Exit with error code if there are errors (not just warnings)
	if stats.hasErrors && !lintFix {
		return fmt.Errorf("linting found errors")
	}

	return nil
}

func processEncryptionPolicyLint(stats *lintStats) {
	cfg, err := config.Load(cfgFile)
	if err != nil {
		errlnf("Error loading config for encryption lint: %v", err)
		stats.hasErrors = true
		return
	}

	if !cfg.Encryption.Enabled {
		return
	}

	results, _, _, err := evaluateEncryptionKeyPolicy(cfg, "")
	if err != nil {
		errlnf("Error checking encryption keys: %v", err)
		stats.hasErrors = true
		return
	}

	issues := make([]lint.Issue, 0, len(results)+1)
	if !cfg.Encryption.EnforceStrength {
		issues = append(issues, lint.Issue{
			Line:     1,
			Severity: lint.SeverityWarning,
			Message:  "encryption.enforce_strength is false; build will not fail on weak encryption keys",
		})
	}

	for _, result := range results {
		if result.Err == nil {
			continue
		}
		issues = append(issues, lint.Issue{
			Line:     1,
			Severity: lint.SeverityError,
			Message:  fmt.Sprintf("encryption key %q failed policy (%s): %v", result.KeyName, result.EnvName, result.Err),
		})
	}

	if len(issues) == 0 {
		return
	}

	stats.filesWithIssues++
	stats.totalIssues += len(issues)
	outln("\n[encryption-config]:")
	for _, issue := range issues {
		printIssue(issue)
		if issue.Severity == lint.SeverityError {
			stats.hasErrors = true
		}
	}
}

// getFilesFromConfig discovers files using the configured glob patterns.
func getFilesFromConfig() ([]string, error) {
	// Load config (will use defaults if no config file found)
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return nil, fmt.Errorf("loading config: %w", err)
	}

	// Get glob patterns from config, default to **/*.md
	patterns := cfg.GlobConfig.Patterns
	if len(patterns) == 0 {
		patterns = []string{"**/*.md"}
	}

	// Get gitignore setting
	useGitignore := cfg.GlobConfig.UseGitignore

	// Parse .gitignore if enabled
	var gitignorePatterns []string
	if useGitignore {
		gitignorePatterns = loadGitignore(".")
	}

	// Get absolute path for base directory
	absBaseDir, err := filepath.Abs(".")
	if err != nil {
		return nil, fmt.Errorf("getting absolute path: %w", err)
	}

	// Expand glob patterns
	fileSet := make(map[string]struct{})
	for _, pattern := range patterns {
		fullPattern := pattern
		if !filepath.IsAbs(pattern) {
			fullPattern = filepath.Join(absBaseDir, pattern)
		}

		matches, err := doublestar.FilepathGlob(fullPattern)
		if err != nil {
			return nil, fmt.Errorf("glob pattern %q: %w", pattern, err)
		}

		for _, match := range matches {
			// Get relative path for consistency
			relPath, err := filepath.Rel(absBaseDir, match)
			if err != nil {
				relPath = match
			}

			// Skip if ignored by gitignore
			if useGitignore && isIgnored(relPath, gitignorePatterns) {
				continue
			}

			// Skip directories
			info, err := os.Stat(match)
			if err != nil {
				continue
			}
			if info.IsDir() {
				continue
			}

			// Only include markdown files
			ext := filepath.Ext(match)
			if ext != markdownShortExtension && ext != markdownExtension {
				continue
			}

			fileSet[relPath] = struct{}{}
		}
	}

	// Convert to sorted slice
	files := make([]string, 0, len(fileSet))
	for file := range fileSet {
		files = append(files, file)
	}
	sort.Strings(files)

	return files, nil
}

// loadGitignore reads .gitignore patterns from the specified directory.
func loadGitignore(baseDir string) []string {
	gitignorePath := filepath.Join(baseDir, ".gitignore")
	file, err := os.Open(gitignorePath)
	if err != nil {
		return nil
	}
	defer file.Close()

	var patterns []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		patterns = append(patterns, line)
	}

	return patterns
}

// isIgnored checks if a path matches any gitignore pattern.
func isIgnored(path string, gitignorePatterns []string) bool {
	if len(gitignorePatterns) == 0 {
		return false
	}

	normalizedPath := filepath.ToSlash(path)

	for _, pattern := range gitignorePatterns {
		// Skip negation patterns
		if strings.HasPrefix(pattern, "!") {
			continue
		}

		normalizedPattern := filepath.ToSlash(pattern)
		normalizedPattern = strings.TrimSuffix(normalizedPattern, "/")

		// Try direct match
		matched, err := doublestar.Match(normalizedPattern, normalizedPath)
		if err == nil && matched {
			return true
		}

		// Try as prefix (for directories)
		if strings.HasPrefix(normalizedPath, normalizedPattern+"/") {
			return true
		}

		// Match against filename
		filename := filepath.Base(normalizedPath)
		matched, err = doublestar.Match(normalizedPattern, filename)
		if err == nil && matched {
			return true
		}

		// Try with **/ prefix
		if !strings.HasPrefix(normalizedPattern, "**/") && !strings.HasPrefix(normalizedPattern, "/") {
			matched, err = doublestar.Match("**/"+normalizedPattern, normalizedPath)
			if err == nil && matched {
				return true
			}
		}
	}

	return false
}

// expandGlobPatterns expands glob patterns from args into a list of files.
func expandGlobPatterns(args []string) ([]string, error) {
	var files []string
	for _, pattern := range args {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			return nil, fmt.Errorf("invalid pattern %q: %w", pattern, err)
		}
		if len(matches) == 0 {
			// Try as a literal file path
			if _, err := os.Stat(pattern); err == nil {
				files = append(files, pattern)
			} else {
				warnf("no files match pattern %q", pattern)
			}
		} else {
			files = append(files, matches...)
		}
	}
	return files, nil
}

// processFile lints a single file and updates stats.
func processFile(file string, stats *lintStats, mentionResolver *blogrollMentionResolver) {
	// Skip non-markdown files
	ext := filepath.Ext(file)
	if ext != markdownShortExtension && ext != markdownExtension {
		return
	}

	stats.totalFiles++

	content, err := os.ReadFile(file)
	if err != nil {
		errlnf("Error reading %s: %v", file, err)
		return
	}

	var result *lint.Result
	switch {
	case lintFix:
		result = lint.Fix(file, string(content))
		if mentionResolver != nil {
			for _, issue := range lint.WithResolver(file, string(content), mentionResolver).Issues {
				if issue.Type == "unknown-mention" {
					result.Issues = append(result.Issues, issue)
				}
			}
		}
	case mentionResolver != nil:
		result = lint.WithResolver(file, string(content), mentionResolver)
	default:
		result = lint.Lint(file, string(content))
	}

	if len(result.Issues) == 0 {
		return
	}

	stats.filesWithIssues++
	stats.totalIssues += len(result.Issues)

	// Count fixable issues
	for _, issue := range result.Issues {
		if issue.Fixable {
			stats.fixableIssues++
		}
	}

	outlnf("\n%s:", file)
	for _, issue := range result.Issues {
		printIssue(issue)
		if issue.Severity == lint.SeverityError {
			stats.hasErrors = true
		}
	}

	// Write fixed content if --fix was used
	if lintFix && result.Fixed != result.Content {
		if err := os.WriteFile(file, []byte(result.Fixed), 0o600); err != nil {
			errlnf("Error writing %s: %v", file, err)
		} else {
			stats.totalFixed++
			outlnf("  -> Fixed %d issue(s)", len(result.Issues))
		}
	}
}

// printIssue prints a single lint issue with colors.
func printIssue(issue lint.Issue) {
	severityText := severityLabel(issue.Severity)

	location := fmt.Sprintf("line %d", issue.Line)
	if issue.Column > 0 {
		location += fmt.Sprintf(", col %d", issue.Column)
	}

	out("  %s [%s]: %s\n", severityText, location, issue.Message)
}

// severityLabel returns a themed severity label.
func severityLabel(severity lint.Severity) string {
	switch severity {
	case lint.SeverityError:
		return colorizeOutput(severity.String(), currentLogTheme.Error)
	case lint.SeverityWarning:
		return colorizeOutput(severity.String(), currentLogTheme.Warning)
	case lint.SeverityInfo:
		return colorizeOutput(severity.String(), currentLogTheme.Component)
	default:
		return severity.String()
	}
}

// printSummary prints the linting summary.
func printSummary(stats *lintStats) {
	outln()
	if stats.totalIssues == 0 {
		outlnf("OK %d file(s) linted, no issues found", stats.totalFiles)
	} else {
		outlnf("FAIL %d file(s) linted, %d issue(s) in %d file(s)",
			stats.totalFiles, stats.totalIssues, stats.filesWithIssues)

		if lintFix {
			outlnf("  -> Fixed %d file(s)", stats.totalFixed)
		} else if stats.fixableIssues > 0 {
			// Show fixable count and suggest fix command
			nonFixable := stats.totalIssues - stats.fixableIssues
			out("INFO %d issue(s) can be automatically fixed", stats.fixableIssues)
			if nonFixable > 0 {
				out(", %d cannot", nonFixable)
			}
			outln()
			outln()
			outln("Run 'markata-go lint --fix' to automatically fix fixable issues")
		}
	}
}
