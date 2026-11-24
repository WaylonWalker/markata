package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/WaylonWalker/markata/pkg/blogroll"
	"github.com/WaylonWalker/markata/pkg/config"
	"github.com/WaylonWalker/markata/pkg/models"
)

var (
	// blogrollDryRun shows what would be updated without making changes.
	blogrollDryRun bool

	// blogrollForce overwrites existing metadata even if present.
	blogrollForce bool

	// blogrollFeed updates only a specific feed by handle/URL.
	blogrollFeed string
)

// blogrollCmd represents the blogroll command group.
var blogrollCmd = &cobra.Command{
	Use:   "blogroll",
	Short: "Manage blogroll feeds",
	Long: `Commands for managing blogroll feeds and metadata.

Subcommands:
  update     - Update feed metadata from external sources`,
}

// blogrollUpdateCmd updates blogroll metadata.
var blogrollUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update feed metadata from external sources",
	Long: `Automatically update blogroll feed metadata by fetching information from:

1. OpenGraph Protocol - og:title, og:description, og:image
2. HTML Meta Tags - description, keywords, author
3. RSS/Atom Feed Metadata - feed title, description, author

This command reads your config file, fetches metadata for each blogroll feed,
and updates the config with any missing or outdated information.

Example usage:
  markata-go blogroll update              # Update all feeds
  markata-go blogroll update --dry-run    # Preview changes without modifying
  markata-go blogroll update --force      # Overwrite existing metadata
  markata-go blogroll update --feed=dave  # Update only feed matching "dave"`,
	RunE: runBlogrollUpdate,
}

func init() {
	rootCmd.AddCommand(blogrollCmd)
	blogrollCmd.AddCommand(blogrollUpdateCmd)

	blogrollUpdateCmd.Flags().BoolVar(&blogrollDryRun, "dry-run", false, "show what would be updated without making changes")
	blogrollUpdateCmd.Flags().BoolVar(&blogrollForce, "force", false, "overwrite existing metadata even if present")
	blogrollUpdateCmd.Flags().StringVar(&blogrollFeed, "feed", "", "update only specific feed by handle or URL substring")
}

func runBlogrollUpdate(_ *cobra.Command, _ []string) error {
	// Discover config file
	configPath := cfgFile
	if configPath == "" {
		var err error
		configPath, err = config.Discover()
		if err != nil {
			return fmt.Errorf("no config file found: %w", err)
		}
	}

	// Load current configuration
	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Check if blogroll is configured
	if !cfg.Blogroll.Enabled {
		fmt.Println("Blogroll is not enabled in configuration.")
		fmt.Println("Add [blogroll] enabled = true to your config file.")
		return nil
	}

	if len(cfg.Blogroll.Feeds) == 0 {
		fmt.Println("No blogroll feeds configured.")
		fmt.Println("Add [[blogroll.feeds]] entries to your config file.")
		return nil
	}

	// Create the updater
	timeout := time.Duration(cfg.Blogroll.Timeout) * time.Second
	if timeout == 0 {
		timeout = 30 * time.Second
	}
	updater := blogroll.NewUpdater(timeout)

	// Track changes
	results := make([]blogroll.UpdateResult, 0, len(cfg.Blogroll.Feeds))
	feedsUpdated := 0
	feedsSkipped := 0
	feedsErrored := 0

	fmt.Printf("Updating blogroll metadata from %s...\n\n", configPath)

	// Process each feed
	for i := range cfg.Blogroll.Feeds {
		feed := &cfg.Blogroll.Feeds[i]

		// Filter by specific feed if requested
		if blogrollFeed != "" {
			if !matchesFeed(feed, blogrollFeed) {
				continue
			}
		}

		result := updateFeedMetadata(updater, feed, blogrollForce)
		results = append(results, result)

		switch {
		case result.Error != "":
			feedsErrored++
			fmt.Printf("  ✗ %s: %s\n", feed.URL, result.Error)
		case result.Updated:
			feedsUpdated++
			printUpdateResult(result, blogrollDryRun)
		default:
			feedsSkipped++
			if verbose {
				fmt.Printf("  - %s: no changes needed\n", feedDisplayName(feed))
			}
		}
	}

	fmt.Println()

	// Summary
	if blogrollFeed != "" && len(results) == 0 {
		fmt.Printf("No feeds matching '%s' found.\n", blogrollFeed)
		return nil
	}

	fmt.Printf("Summary: %d updated, %d skipped, %d errors\n", feedsUpdated, feedsSkipped, feedsErrored)

	// Write changes if not dry-run and there are updates
	if !blogrollDryRun && feedsUpdated > 0 {
		if err := writeConfigUpdate(configPath, cfg); err != nil {
			return fmt.Errorf("failed to write config: %w", err)
		}
		fmt.Printf("\nConfig updated: %s\n", configPath)
	} else if blogrollDryRun && feedsUpdated > 0 {
		fmt.Println("\n(dry-run mode - no changes written)")
	}

	return nil
}

// matchesFeed checks if a feed matches the filter string.
func matchesFeed(feed *configFeedRef, filter string) bool {
	filter = strings.ToLower(filter)

	// Match by URL
	if strings.Contains(strings.ToLower(feed.URL), filter) {
		return true
	}

	// Match by title
	if strings.Contains(strings.ToLower(feed.Title), filter) {
		return true
	}

	return false
}

// configFeedRef is a type alias for cleaner code.
type configFeedRef = models.ExternalFeedConfig

// updateFeedMetadata fetches and applies metadata updates to a feed.
func updateFeedMetadata(updater *blogroll.Updater, feed *configFeedRef, force bool) blogroll.UpdateResult {
	result := blogroll.UpdateResult{
		FeedURL: feed.URL,
		Handle:  feedDisplayName(feed),
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Fetch metadata
	metadata, err := updater.FetchMetadata(ctx, feed.URL)
	if err != nil {
		result.Error = err.Error()
		return result
	}

	// Build old metadata for comparison
	result.OldMetadata = &blogroll.Metadata{
		Title:       feed.Title,
		Description: feed.Description,
		ImageURL:    feed.ImageURL,
		SiteURL:     feed.SiteURL,
	}

	// Apply updates
	result.NewMetadata = &blogroll.Metadata{}

	// Title
	if feed.Title == "" || force {
		if metadata.Title != "" && metadata.Title != feed.Title {
			feed.Title = metadata.Title
			result.NewMetadata.Title = metadata.Title
			result.Updated = true
		}
	}

	// Description
	if feed.Description == "" || force {
		if metadata.Description != "" && metadata.Description != feed.Description {
			feed.Description = metadata.Description
			result.NewMetadata.Description = metadata.Description
			result.Updated = true
		}
	}

	// ImageURL
	if feed.ImageURL == "" || force {
		if metadata.ImageURL != "" && metadata.ImageURL != feed.ImageURL {
			feed.ImageURL = metadata.ImageURL
			result.NewMetadata.ImageURL = metadata.ImageURL
			result.Updated = true
		}
	}

	// SiteURL
	if feed.SiteURL == "" || force {
		if metadata.SiteURL != "" && metadata.SiteURL != feed.SiteURL {
			feed.SiteURL = metadata.SiteURL
			result.NewMetadata.SiteURL = metadata.SiteURL
			result.Updated = true
		}
	}

	return result
}

// feedDisplayName returns a display name for a feed.
func feedDisplayName(feed *configFeedRef) string {
	if feed.Title != "" {
		return feed.Title
	}
	return feed.URL
}

// printUpdateResult prints the changes for a feed.
func printUpdateResult(result blogroll.UpdateResult, dryRun bool) {
	prefix := "  ✓"
	if dryRun {
		prefix = "  →"
	}

	fmt.Printf("%s %s:\n", prefix, result.Handle)

	if result.NewMetadata.Title != "" {
		fmt.Printf("      title: %q\n", result.NewMetadata.Title)
	}
	if result.NewMetadata.Description != "" {
		desc := result.NewMetadata.Description
		if len(desc) > 60 {
			desc = desc[:57] + "..."
		}
		fmt.Printf("      description: %q\n", desc)
	}
	if result.NewMetadata.ImageURL != "" {
		fmt.Printf("      image_url: %s\n", result.NewMetadata.ImageURL)
	}
	if result.NewMetadata.SiteURL != "" {
		fmt.Printf("      site_url: %s\n", result.NewMetadata.SiteURL)
	}
}

// writeConfigUpdate writes the updated config back to the file.
func writeConfigUpdate(configPath string, cfg *models.Config) error {
	// Determine format from extension
	ext := strings.ToLower(filepath.Ext(configPath))

	// Read the original file to preserve structure
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("read config: %w", err)
	}

	// Parse to a generic map to preserve other settings
	var configMap map[string]interface{}

	switch ext {
	case formatTOML:
		if err := toml.Unmarshal(data, &configMap); err != nil {
			return fmt.Errorf("parse TOML: %w", err)
		}
	case extYAML, extYML:
		if err := yaml.Unmarshal(data, &configMap); err != nil {
			return fmt.Errorf("parse YAML: %w", err)
		}
	default:
		return fmt.Errorf("unsupported config format: %s", ext)
	}

	// Update the blogroll section
	updateBlogrollInMap(configMap, cfg.Blogroll)

	// Write back
	switch ext {
	case formatTOML:
		return writeTOMLConfig(configPath, configMap)
	case extYAML, extYML:
		return writeYAMLConfig(configPath, configMap)
	default:
		return fmt.Errorf("unsupported config format: %s", ext)
	}
}

// updateBlogrollInMap updates the blogroll section in a config map.
func updateBlogrollInMap(configMap map[string]interface{}, blogrollCfg models.BlogrollConfig) {
	// Find the markata-go wrapper or use root
	var target map[string]interface{}
	if mg, ok := configMap["markata-go"].(map[string]interface{}); ok {
		target = mg
	} else {
		target = configMap
	}

	// Get or create blogroll section
	var blogrollMap map[string]interface{}
	if br, ok := target["blogroll"].(map[string]interface{}); ok {
		blogrollMap = br
	} else {
		blogrollMap = make(map[string]interface{})
		target["blogroll"] = blogrollMap
	}

	// Convert feeds to map format
	feeds := make([]map[string]interface{}, len(blogrollCfg.Feeds))
	for i := range blogrollCfg.Feeds {
		feed := &blogrollCfg.Feeds[i]
		feedMap := map[string]interface{}{
			"url": feed.URL,
		}
		if feed.Title != "" {
			feedMap["title"] = feed.Title
		}
		if feed.Description != "" {
			feedMap["description"] = feed.Description
		}
		if feed.Category != "" {
			feedMap["category"] = feed.Category
		}
		if len(feed.Tags) > 0 {
			feedMap["tags"] = feed.Tags
		}
		if feed.SiteURL != "" {
			feedMap["site_url"] = feed.SiteURL
		}
		if feed.ImageURL != "" {
			feedMap["image_url"] = feed.ImageURL
		}
		if feed.Active != nil {
			feedMap["active"] = *feed.Active
		}
		feeds[i] = feedMap
	}

	blogrollMap["feeds"] = feeds
}

// writeTOMLConfig writes a config map as TOML.
func writeTOMLConfig(path string, configMap map[string]interface{}) error {
	var buf strings.Builder
	encoder := toml.NewEncoder(&buf)
	if err := encoder.Encode(configMap); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(buf.String()), 0o644) //nolint:gosec // config files should be readable
}

// writeYAMLConfig writes a config map as YAML.
func writeYAMLConfig(path string, configMap map[string]interface{}) error {
	data, err := yaml.Marshal(configMap)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644) //nolint:gosec // config files should be readable
}
