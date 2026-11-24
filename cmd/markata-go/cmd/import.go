package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/WaylonWalker/markata/pkg/importer"
	"github.com/spf13/cobra"
)

var (
	// importOutputDir is the output directory for imported posts.
	importOutputDir string

	// importSince filters posts to only those after this date.
	importSince string

	// importDryRun previews imports without writing files.
	importDryRun bool

	// importAddTags are additional tags to add to all imported posts.
	importAddTags string
)

// importCmd represents the import command group.
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import content from external sources (PESOS pattern)",
	Long: `Import content from external sources into your markata-go site.

This command implements the PESOS (Publish Elsewhere, Syndicate to Own Site)
pattern, allowing you to import your content from various platforms and
consolidate it on your own site.

Subcommands:
  rss       - Import from RSS/Atom feeds
  jsonfeed  - Import from JSON Feed format

Example usage:
  markata-go import rss https://example.com/feed.xml
  markata-go import jsonfeed https://example.com/feed.json
  markata-go import rss https://example.com/feed.xml --output posts/imported
  markata-go import rss https://example.com/feed.xml --since 2024-01-01
  markata-go import rss https://example.com/feed.xml --dry-run`,
}

// importRSSCmd imports content from RSS/Atom feeds.
var importRSSCmd = &cobra.Command{
	Use:   "rss <url>",
	Short: "Import content from an RSS/Atom feed",
	Long: `Import content from an RSS or Atom feed.

The command fetches the feed, parses the entries, and creates markdown files
with appropriate frontmatter for each post.

Example usage:
  markata-go import rss https://example.com/feed.xml
  markata-go import rss https://example.com/feed.xml --output posts/blog
  markata-go import rss https://example.com/feed.xml --since 2024-01-01 --dry-run`,
	Args: cobra.ExactArgs(1),
	RunE: runImportRSSCommand,
}

// importJSONFeedCmd imports content from JSON Feed format.
var importJSONFeedCmd = &cobra.Command{
	Use:   "jsonfeed <url>",
	Short: "Import content from a JSON Feed",
	Long: `Import content from a JSON Feed (https://jsonfeed.org).

The command fetches the feed, parses the entries, and creates markdown files
with appropriate frontmatter for each post.

Example usage:
  markata-go import jsonfeed https://example.com/feed.json
  markata-go import jsonfeed https://example.com/feed.json --output posts/external
  markata-go import jsonfeed https://example.com/feed.json --dry-run`,
	Args: cobra.ExactArgs(1),
	RunE: runImportJSONFeedCommand,
}

func init() {
	rootCmd.AddCommand(importCmd)

	// Add subcommands
	importCmd.AddCommand(importRSSCmd)
	importCmd.AddCommand(importJSONFeedCmd)

	// Common flags for all import subcommands
	importCmd.PersistentFlags().StringVarP(&importOutputDir, "output", "o", "posts/imported", "Output directory for imported posts")
	importCmd.PersistentFlags().StringVar(&importSince, "since", "", "Only import posts published after this date (YYYY-MM-DD)")
	importCmd.PersistentFlags().BoolVar(&importDryRun, "dry-run", false, "Preview imports without writing files")
	importCmd.PersistentFlags().StringVar(&importAddTags, "tags", "", "Additional tags to add to all imported posts (comma-separated)")
}

// runImportRSSCommand runs the RSS import.
func runImportRSSCommand(_ *cobra.Command, args []string) error {
	url := args[0]

	// Create importer
	imp, err := importer.NewRSSImporter(url)
	if err != nil {
		return err
	}

	return runImport(imp)
}

// runImportJSONFeedCommand runs the JSON Feed import.
func runImportJSONFeedCommand(_ *cobra.Command, args []string) error {
	url := args[0]

	// Create importer
	imp, err := importer.NewJSONFeedImporter(url)
	if err != nil {
		return err
	}

	return runImport(imp)
}

// runImport executes the import with the given importer.
func runImport(imp importer.Importer) error {
	// Parse options
	opts := importer.ImportOptions{
		DryRun:    importDryRun,
		OutputDir: importOutputDir,
	}

	// Parse --since flag
	if importSince != "" {
		since, err := time.Parse("2006-01-02", importSince)
		if err != nil {
			return fmt.Errorf("invalid date format for --since (use YYYY-MM-DD): %w", err)
		}
		opts.Since = &since
	}

	// Parse --tags flag
	if importAddTags != "" {
		for _, tag := range strings.Split(importAddTags, ",") {
			tag = strings.TrimSpace(tag)
			if tag != "" {
				opts.AddTags = append(opts.AddTags, tag)
			}
		}
	}

	// Print header
	fmt.Printf("Importing from %s feed: %s\n", imp.Name(), imp.SourceURL())
	if importDryRun {
		fmt.Println("(dry-run mode - no files will be written)")
	}
	fmt.Println(strings.Repeat("-", 60))

	// Fetch and parse feed
	posts, err := imp.Import(opts)
	if err != nil {
		return fmt.Errorf("import failed: %w", err)
	}

	if len(posts) == 0 {
		fmt.Println("No posts found to import.")
		if opts.Since != nil {
			fmt.Printf("(filtered to posts after %s)\n", opts.Since.Format("2006-01-02"))
		}
		return nil
	}

	fmt.Printf("Found %d post(s) to import\n", len(posts))
	if opts.Since != nil {
		fmt.Printf("(filtered to posts after %s)\n", opts.Since.Format("2006-01-02"))
	}
	fmt.Println()

	// Write posts
	writer := importer.NewWriter(opts.OutputDir)
	result, err := writer.Write(posts, opts.DryRun)
	if err != nil {
		return fmt.Errorf("write failed: %w", err)
	}

	// Print results
	for _, path := range result.Paths {
		if opts.DryRun {
			fmt.Printf("  [dry-run] Would write: %s\n", path)
		} else {
			fmt.Printf("  Created: %s\n", path)
		}
	}

	if result.Skipped > 0 {
		fmt.Printf("\nSkipped %d post(s) (already exist)\n", result.Skipped)
	}

	for _, e := range result.Errors {
		fmt.Fprintf(os.Stderr, "  Warning: %v\n", e)
	}

	fmt.Println(strings.Repeat("-", 60))
	if opts.DryRun {
		fmt.Printf("Dry run complete: %d post(s) would be imported\n", result.Written)
	} else {
		fmt.Printf("Import complete: %d post(s) imported to %s/\n", result.Written, opts.OutputDir)
	}

	return nil
}
