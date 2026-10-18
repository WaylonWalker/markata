package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/WaylonWalker/markata/pkg/plugins"
)

var (
	// webmentionsCacheDir overrides the configured cache directory.
	webmentionsCacheDir string

	// webmentionsVerbose enables verbose output.
	webmentionsVerbose bool
)

// webmentionsCmd represents the webmentions command group.
var webmentionsCmd = &cobra.Command{
	Use:   "webmentions",
	Short: "Manage webmentions",
	Long: `Commands for managing webmentions (incoming and outgoing).

Subcommands:
  fetch      - Fetch incoming webmentions from webmention.io`,
}

// webmentionsFetchCmd fetches incoming webmentions.
var webmentionsFetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch incoming webmentions from webmention.io",
	Long: `Fetch all incoming webmentions from webmention.io API.

This command requires a webmention.io API token to be configured.
You can get your token from: https://webmention.io/settings

Configuration options:
  1. Add to your config file:
     [markata-go.webmentions]
     webmention_io_token = "your_token_here"

  2. Set environment variable:
     export WEBMENTION_IO_TOKEN="your_token_here"

  3. Add to .env file:
     WEBMENTION_IO_TOKEN=your_token_here

The fetched mentions will be cached in .cache/webmentions/received_mentions.json

Example usage:
  markata-go webmentions fetch                           # Fetch all mentions
  markata-go webmentions fetch --verbose                 # Show detailed output
  markata-go webmentions fetch --cache-dir=/tmp/cache    # Use custom cache dir`,
	RunE: runWebmentionsFetch,
}

func init() {
	rootCmd.AddCommand(webmentionsCmd)
	webmentionsCmd.AddCommand(webmentionsFetchCmd)

	webmentionsFetchCmd.Flags().StringVar(&webmentionsCacheDir, "cache-dir", "", "override configured cache directory")
	webmentionsFetchCmd.Flags().BoolVarP(&webmentionsVerbose, "verbose", "v", false, "enable verbose output")
}

func runWebmentionsFetch(_ *cobra.Command, _ []string) error {
	// Use createManager to properly set up the manager with config
	manager, err := createManager(cfgFile)
	if err != nil {
		return fmt.Errorf("failed to create manager: %w", err)
	}

	if webmentionsVerbose {
		fmt.Printf("Using site URL: %s\n", manager.Config().Extra["url"])
	}

	// Create and configure the fetch plugin
	fetchPlugin := plugins.NewWebmentionsFetchPlugin()
	if err := fetchPlugin.Configure(manager); err != nil {
		return fmt.Errorf("failed to configure plugin: %w", err)
	}

	fmt.Println("Fetching webmentions from webmention.io...")

	// Fetch mentions
	if err := fetchPlugin.FetchMentions(); err != nil {
		return fmt.Errorf("failed to fetch mentions: %w", err)
	}

	mentions := fetchPlugin.GetMentions()
	fmt.Printf("✓ Fetched %d webmentions\n", len(mentions))

	if webmentionsVerbose {
		// Group by type
		typeCount := make(map[string]int)
		for i := range mentions {
			typeCount[mentions[i].WMProperty]++
		}

		fmt.Println("\nBreakdown by type:")
		for wmType, count := range typeCount {
			fmt.Printf("  %s: %d\n", wmType, count)
		}

		// Group by URL
		urlGroups := fetchPlugin.GroupMentionsByURL()
		fmt.Printf("\nMentions across %d unique URLs\n", len(urlGroups))
	}

	// Show cache location
	cacheDir := ".cache/webmentions"
	if webmentionsCacheDir != "" {
		cacheDir = webmentionsCacheDir
	}
	fmt.Printf("\nMentions cached to: %s/received_mentions.json\n", cacheDir)

	return nil
}
