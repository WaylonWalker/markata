package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/WaylonWalker/markata/pkg/migrate"
	"github.com/spf13/cobra"
)

// defaultTemplatesDir is the default directory name for templates.
const defaultTemplatesDir = "templates"

var (
	// migrateInput is the input config file path.
	migrateInput string

	// migrateOutput is the output config file path.
	migrateOutput string

	// migrateDryRun shows changes without writing.
	migrateDryRun bool

	// migrateFormat is the output format (toml, yaml, json).
	migrateFormat string

	// migrateJSON outputs results as JSON.
	migrateJSON bool

	// migrateReport is the path to write a migration report file.
	migrateReport string

	// compareOldDir is the old site directory for comparison.
	compareOldDir string

	// compareNewDir is the new site directory for comparison.
	compareNewDir string

	// compareExtensions are file extensions to compare.
	compareExtensions []string

	// compareIgnore are patterns to ignore.
	compareIgnore []string
)

// migrateCmd represents the migrate command.
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate from Python markata to markata-go",
	Long: `Migrate helps you transition from Python markata to markata-go.

This command analyzes your existing Python markata configuration and:
  - Converts configuration namespaces and keys
  - Migrates filter expressions to markata-go syntax
  - Identifies template compatibility issues
  - Generates a detailed migration report

Example usage:
  markata-go migrate                    # Analyze and show migration report
  markata-go migrate --dry-run          # Show what would change
  markata-go migrate -o markata-go.toml # Write migrated config
  markata-go migrate -i pyproject.toml  # Use specific input file`,
	RunE: runMigrateCommand,
}

// migrateConfigCmd migrates configuration only.
var migrateConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Migrate configuration file only",
	Long: `Migrate configuration from Python markata to markata-go format.

Handles:
  - Namespace changes ([markata] -> [markata-go])
  - Key renames (glob_patterns -> patterns, etc.)
  - Nav map to array conversion
  - Feed configuration migration`,
	RunE: runMigrateConfigCommand,
}

// migrateFilterCmd migrates filter expressions.
var migrateFilterCmd = &cobra.Command{
	Use:   "filter [expression]",
	Short: "Check and migrate filter expressions",
	Long: `Migrate filter expressions from Python markata syntax to markata-go.

If an expression is provided, it will be analyzed and migrated.
If no expression is provided, all filters in the config will be checked.

Examples:
  markata-go migrate filter "published == 'True'"
  markata-go migrate filter "templateKey in ['blog-post', 'til']"
  markata-go migrate filter                        # Check all filters in config`,
	RunE: runMigrateFilterCommand,
}

// migrateTemplatesCmd checks template compatibility.
var migrateTemplatesCmd = &cobra.Command{
	Use:   "templates [path]",
	Short: "Validate template compatibility",
	Long: `Check templates for compatibility with markata-go (pongo2).

Identifies:
  - Unsupported Jinja2 features (macros, do statements, etc.)
  - Variable name changes (post.markata -> config)
  - Python expression usage
  - Filter syntax differences

Example:
  markata-go migrate templates
  markata-go migrate templates ./templates`,
	RunE: runMigrateTemplatesCommand,
}

// migrateCompareCmd compares two site output directories.
var migrateCompareCmd = &cobra.Command{
	Use:   "compare",
	Short: "Compare old and new site output directories",
	Long: `Compare two directories to verify migration completeness.

This command compares the file structure between your old site output
and the new markata-go output to identify:
  - Files missing in the new site
  - New files not present in the old site
  - Common files present in both

Example usage:
  markata-go migrate compare --old markout --new public
  markata-go migrate compare --old _site --new public --ext .html,.xml
  markata-go migrate compare --old dist --new public --ignore "assets/**"

Exit codes:
  0 = directories match (no missing files)
  1 = differences found`,
	RunE: runMigrateCompareCommand,
}

func init() {
	rootCmd.AddCommand(migrateCmd)

	// Add subcommands
	migrateCmd.AddCommand(migrateConfigCmd)
	migrateCmd.AddCommand(migrateFilterCmd)
	migrateCmd.AddCommand(migrateTemplatesCmd)
	migrateCmd.AddCommand(migrateCompareCmd)

	// Flags for migrate command
	migrateCmd.Flags().StringVarP(&migrateInput, "input", "i", "", "input config file (default: auto-detect)")
	migrateCmd.Flags().StringVarP(&migrateOutput, "output", "o", "", "output config file")
	migrateCmd.Flags().BoolVarP(&migrateDryRun, "dry-run", "n", false, "show changes without writing")
	migrateCmd.Flags().StringVarP(&migrateFormat, "format", "f", "toml", "output format (toml, yaml)")
	migrateCmd.Flags().BoolVar(&migrateJSON, "json", false, "output results as JSON")
	migrateCmd.Flags().StringVar(&migrateReport, "report", "", "write migration report to file")

	// Flags for config subcommand
	migrateConfigCmd.Flags().StringVarP(&migrateInput, "input", "i", "", "input config file (default: auto-detect)")
	migrateConfigCmd.Flags().StringVarP(&migrateOutput, "output", "o", "", "output config file")
	migrateConfigCmd.Flags().BoolVarP(&migrateDryRun, "dry-run", "n", false, "show changes without writing")
	migrateConfigCmd.Flags().StringVarP(&migrateFormat, "format", "f", "toml", "output format (toml, yaml)")
	migrateConfigCmd.Flags().BoolVar(&migrateJSON, "json", false, "output results as JSON")

	// Flags for filter subcommand
	migrateFilterCmd.Flags().StringVarP(&migrateInput, "input", "i", "", "input config file (default: auto-detect)")
	migrateFilterCmd.Flags().BoolVar(&migrateJSON, "json", false, "output results as JSON")

	// Flags for templates subcommand
	migrateTemplatesCmd.Flags().BoolVar(&migrateJSON, "json", false, "output results as JSON")

	// Flags for compare subcommand
	migrateCompareCmd.Flags().StringVar(&compareOldDir, "old", "", "old site output directory (required)")
	migrateCompareCmd.Flags().StringVar(&compareNewDir, "new", "", "new site output directory (required)")
	migrateCompareCmd.Flags().StringSliceVar(&compareExtensions, "ext", []string{".html"}, "file extensions to compare (comma-separated)")
	migrateCompareCmd.Flags().StringSliceVar(&compareIgnore, "ignore", []string{}, "glob patterns to ignore (comma-separated)")
	migrateCompareCmd.Flags().BoolVar(&migrateJSON, "json", false, "output results as JSON")

	// Mark required flags for compare (errors ignored as they only occur if flag doesn't exist)
	//nolint:errcheck // errors only occur if flag doesn't exist, which we know it does
	migrateCompareCmd.MarkFlagRequired("old")
	//nolint:errcheck // errors only occur if flag doesn't exist, which we know it does
	migrateCompareCmd.MarkFlagRequired("new")
}

// runMigrateCommand runs the full migration analysis.
func runMigrateCommand(_ *cobra.Command, _ []string) error {
	// Find input file
	inputPath, err := findInputConfig()
	if err != nil {
		return err
	}

	// Determine output path
	outputPath := migrateOutput
	if migrateDryRun {
		outputPath = "" // Don't write on dry run
	} else if outputPath == "" && migrateOutput == "" {
		// Default output path based on format
		outputPath = ""
	}

	// Apply format to output path if specified but not in output
	if outputPath != "" && migrateFormat != "" {
		ext := filepath.Ext(outputPath)
		if ext == "" {
			outputPath = outputPath + "." + migrateFormat
		}
	}

	// Run migration
	result, err := migrate.Config(inputPath, outputPath)
	if err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	// Check templates if directory exists
	templatesDir := defaultTemplatesDir
	if info, err := os.Stat(templatesDir); err == nil && info.IsDir() {
		templateIssues, err := migrate.CheckTemplates(templatesDir)
		if err != nil {
			result.Warnings = append(result.Warnings, migrate.Warning{
				Category: "template",
				Message:  fmt.Sprintf("Failed to check templates: %v", err),
			})
		} else {
			result.TemplateIssues = templateIssues
		}
	}

	// Output results
	if migrateJSON {
		return outputJSON(result.JSONReport())
	}

	report := result.Report()
	fmt.Print(report)

	// Write report to file if requested
	if migrateReport != "" {
		if err := os.WriteFile(migrateReport, []byte(report), 0o600); err != nil {
			return fmt.Errorf("failed to write report: %w", err)
		}
		fmt.Printf("\nReport written to: %s\n", migrateReport)
	}

	// Exit with appropriate code
	if result.HasErrors() {
		return newExitCodeError(result.ExitCode(), nil)
	}

	return nil
}

// runMigrateConfigCommand runs config-only migration.
func runMigrateConfigCommand(_ *cobra.Command, _ []string) error {
	// Find input file
	inputPath, err := findInputConfig()
	if err != nil {
		return err
	}

	// Determine output path
	outputPath := migrateOutput
	if migrateDryRun {
		outputPath = ""
	}

	// Run migration
	result, err := migrate.Config(inputPath, outputPath)
	if err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	// Output results
	if migrateJSON {
		return outputJSON(result.JSONReport())
	}

	fmt.Print(result.Report())
	return nil
}

// runMigrateFilterCommand checks/migrates filter expressions.
func runMigrateFilterCommand(_ *cobra.Command, args []string) error {
	// If expression provided, migrate it directly
	if len(args) > 0 {
		expr := args[0]
		migrated, changes := migrate.Filter(expr)

		if migrateJSON {
			return outputJSON(map[string]interface{}{
				"original": expr,
				"migrated": migrated,
				"changes":  changes,
				"valid":    migrate.ValidateFilter(migrated) == nil,
			})
		}

		fmt.Printf("Original: %s\n", expr)
		fmt.Printf("Migrated: %s\n", migrated)

		if len(changes) == 0 {
			fmt.Println("No changes needed")
		} else {
			fmt.Println("Changes:")
			for _, c := range changes {
				fmt.Printf("  - %s\n", c)
			}
		}

		if err := migrate.ValidateFilter(migrated); err != nil {
			fmt.Printf("\nWarning: Migrated expression may have issues: %v\n", err)
		} else {
			fmt.Println("\nMigrated expression is valid")
		}

		return nil
	}

	// Otherwise, check all filters in config
	inputPath, err := findInputConfig()
	if err != nil {
		return err
	}

	result, err := migrate.Config(inputPath, "")
	if err != nil {
		return fmt.Errorf("failed to analyze config: %w", err)
	}

	if migrateJSON {
		return outputJSON(map[string]interface{}{
			"filter_migrations": result.FilterMigrations,
		})
	}

	if len(result.FilterMigrations) == 0 {
		fmt.Println("No filter expressions found in configuration")
		return nil
	}

	fmt.Printf("Found %d filter expressions:\n\n", len(result.FilterMigrations))
	for _, fm := range result.FilterMigrations {
		feedName := fm.Feed
		if feedName == "" {
			feedName = "(unnamed)"
		}

		fmt.Printf("Feed: %s\n", feedName)
		fmt.Printf("  Original: %s\n", fm.Original)
		fmt.Printf("  Migrated: %s\n", fm.Migrated)

		if len(fm.Changes) > 0 {
			fmt.Println("  Changes:")
			for _, c := range fm.Changes {
				fmt.Printf("    - %s\n", c)
			}
		} else {
			fmt.Println("  No changes needed")
		}

		if !fm.Valid {
			fmt.Printf("  Warning: %s\n", fm.Error)
		}
		fmt.Println()
	}

	return nil
}

// runMigrateTemplatesCommand checks template compatibility.
func runMigrateTemplatesCommand(_ *cobra.Command, args []string) error {
	templatesDir := defaultTemplatesDir
	if len(args) > 0 {
		templatesDir = args[0]
	}

	// Check if directory exists
	info, err := os.Stat(templatesDir)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("templates directory not found: %s", templatesDir)
		}
		return fmt.Errorf("failed to access templates directory: %w", err)
	}

	if !info.IsDir() {
		return fmt.Errorf("not a directory: %s", templatesDir)
	}

	// Check templates
	issues, err := migrate.CheckTemplates(templatesDir)
	if err != nil {
		return fmt.Errorf("failed to check templates: %w", err)
	}

	if migrateJSON {
		return outputJSON(map[string]interface{}{
			"templates_dir": templatesDir,
			"issues":        issues,
			"issues_count":  len(issues),
		})
	}

	if len(issues) == 0 {
		fmt.Printf("No compatibility issues found in %s\n", templatesDir)
		return nil
	}

	fmt.Printf("Found %d template issues in %s:\n\n", len(issues), templatesDir)

	// Group by severity
	var errorIssues, warningIssues, infoIssues []migrate.TemplateIssue
	for _, issue := range issues {
		switch issue.Severity {
		case buildStatusError:
			errorIssues = append(errorIssues, issue)
		case buildStatusWarning:
			warningIssues = append(warningIssues, issue)
		default:
			infoIssues = append(infoIssues, issue)
		}
	}

	if len(errorIssues) > 0 {
		fmt.Println("ERRORS:")
		for _, issue := range errorIssues {
			printTemplateIssue(issue)
		}
		fmt.Println()
	}

	if len(warningIssues) > 0 {
		fmt.Println("WARNINGS:")
		for _, issue := range warningIssues {
			printTemplateIssue(issue)
		}
		fmt.Println()
	}

	if len(infoIssues) > 0 {
		fmt.Println("INFO:")
		for _, issue := range infoIssues {
			printTemplateIssue(issue)
		}
	}

	// Print variable migration reference
	fmt.Println("\nVariable Migration Reference:")
	for _, v := range migrate.GetVariableMigrations() {
		fmt.Printf("  %s -> %s\n", v.Old, v.New)
	}

	fmt.Println("\nFilter Migration Reference:")
	for _, v := range migrate.GetFilterMigrations() {
		fmt.Printf("  %s -> %s\n", v.Old, v.New)
	}

	return nil
}

// runMigrateCompareCommand compares two site output directories.
func runMigrateCompareCommand(_ *cobra.Command, _ []string) error {
	// Build compare options
	opts := migrate.CompareOptions{
		Extensions:     compareExtensions,
		IgnorePatterns: compareIgnore,
	}

	// If no extensions specified, use default
	if len(opts.Extensions) == 0 {
		opts.Extensions = []string{".html"}
	}

	// Run comparison
	result, err := migrate.Compare(compareOldDir, compareNewDir, opts)
	if err != nil {
		return fmt.Errorf("comparison failed: %w", err)
	}

	// Output results
	if migrateJSON {
		return outputJSON(result.JSONReport())
	}

	fmt.Print(result.Report())

	// Exit with appropriate code for scripting
	if result.HasDifferences() {
		return newExitCodeError(result.ExitCode(), nil)
	}

	return nil
}

// findInputConfig finds the input configuration file.
func findInputConfig() (string, error) {
	if migrateInput != "" {
		if _, err := os.Stat(migrateInput); err != nil {
			return "", fmt.Errorf("input file not found: %s", migrateInput)
		}
		return migrateInput, nil
	}

	// Try common Python markata config files
	candidates := []string{
		"markata.toml",
		"pyproject.toml",
		"markata.yaml",
		"markata.yml",
	}

	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
	}

	return "", fmt.Errorf("no Python markata config file found (tried: %v)", candidates)
}

// printTemplateIssue prints a template issue.
func printTemplateIssue(issue migrate.TemplateIssue) {
	fmt.Printf("  %s:%d\n", issue.File, issue.Line)
	fmt.Printf("    %s\n", issue.Issue)
	if issue.Suggestion != "" {
		fmt.Printf("    Suggestion: %s\n", issue.Suggestion)
	}
}

// outputJSON outputs data as JSON.
func outputJSON(data interface{}) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(data)
}
