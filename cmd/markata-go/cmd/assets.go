package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/WaylonWalker/markata/pkg/assets"
	"github.com/WaylonWalker/markata/pkg/config"
	"github.com/spf13/cobra"
)

// Constants for repeated strings.
const (
	statusYes = "yes"
	statusNo  = "no"
)

// assetsCmd represents the assets command group.
var assetsCmd = &cobra.Command{
	Use:   "assets",
	Short: "Manage external CDN assets for self-hosting",
	Long: `Manage external CDN assets (GLightbox, HTMX, Mermaid, etc.) for self-hosting.

When self-hosting is enabled (assets.mode = "self-hosted"), external assets are
downloaded from CDNs and served from your site's output directory.

Available subcommands:
  download - Download all external assets to the cache
  list     - List all assets and their status
  clean    - Remove the assets cache

Example usage:
  markata-go assets download    # Download all CDN assets
  markata-go assets list        # List asset status
  markata-go assets clean       # Clear the cache`,
}

// assetsDownloadCmd downloads all external assets.
var assetsDownloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download all external CDN assets",
	Long: `Download all external CDN assets to the local cache.

Assets are downloaded from their CDN URLs and stored in the cache directory
(default: .markata/assets-cache). These cached assets can then be served
from your site instead of loading from external CDNs.

Example:
  markata-go assets download`,
	RunE: runAssetsDownload,
}

// assetsListCmd lists all assets and their status.
var assetsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all external assets and their status",
	Long: `List all registered external assets and whether they are cached.

Shows the asset name, version, type (JS/CSS), cache status, and size.

Example:
  markata-go assets list`,
	RunE: runAssetsList,
}

// assetsCleanCmd cleans the assets cache.
var assetsCleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Remove the assets cache",
	Long: `Remove all cached assets from the local cache directory.

This forces assets to be re-downloaded on the next build when
self-hosting is enabled.

Example:
  markata-go assets clean`,
	RunE: runAssetsClean,
}

func init() {
	rootCmd.AddCommand(assetsCmd)
	assetsCmd.AddCommand(assetsDownloadCmd)
	assetsCmd.AddCommand(assetsListCmd)
	assetsCmd.AddCommand(assetsCleanCmd)
}

// getAssetsConfig loads the config and returns the assets configuration.
func getAssetsConfig() (downloader *assets.Downloader, cacheDir string) {
	cfg, err := config.Load(cfgFile)
	if err != nil {
		// Use defaults if no config
		cacheDir = ".markata/assets-cache"
		return assets.NewDownloader(cacheDir, true), cacheDir
	}

	cacheDir = cfg.Assets.GetCacheDir()
	verifyIntegrity := cfg.Assets.IsVerifyIntegrityEnabled()
	return assets.NewDownloader(cacheDir, verifyIntegrity), cacheDir
}

func runAssetsDownload(_ *cobra.Command, _ []string) error {
	downloader, _ := getAssetsConfig()

	fmt.Println("Downloading external CDN assets...")
	fmt.Println()

	ctx := context.Background()
	startTime := time.Now()
	results := downloader.DownloadAll(ctx, 4)

	// Print results
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ASSET\tSTATUS\tSIZE\tTIME")
	fmt.Fprintln(w, "-----\t------\t----\t----")

	var totalSize int64
	var successCount, cachedCount, errorCount int

	for i := range results {
		result := &results[i]
		var status string
		switch {
		case result.Error != nil:
			status = fmt.Sprintf("error: %v", result.Error)
			errorCount++
		case result.Cached:
			status = "cached"
			cachedCount++
		default:
			status = "downloaded"
			successCount++
		}

		sizeStr := formatSize(result.Size)
		timeStr := result.Duration.Truncate(time.Millisecond).String()
		if result.Cached {
			timeStr = "-"
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", result.Asset.Name, status, sizeStr, timeStr)
		totalSize += result.Size
	}

	w.Flush()

	duration := time.Since(startTime)
	fmt.Println()
	fmt.Printf("Total: %d downloaded, %d cached, %d errors (%s in %v)\n",
		successCount, cachedCount, errorCount, formatSize(totalSize), duration.Truncate(time.Millisecond))

	if errorCount > 0 {
		return fmt.Errorf("%d assets failed to download", errorCount)
	}

	return nil
}

func runAssetsList(_ *cobra.Command, _ []string) error {
	downloader, cacheDir := getAssetsConfig()

	fmt.Printf("Assets cache directory: %s\n\n", cacheDir)

	statuses := downloader.Status()

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ASSET\tVERSION\tTYPE\tCACHED\tSIZE")
	fmt.Fprintln(w, "-----\t-------\t----\t------\t----")

	var cachedCount int
	var totalSize int64

	for i := range statuses {
		status := &statuses[i]
		cached := statusNo
		sizeStr := "-"
		if status.Cached {
			cached = statusYes
			sizeStr = formatSize(status.Size)
			cachedCount++
			totalSize += status.Size
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			status.Asset.Name,
			status.Asset.Version,
			status.Asset.Type,
			cached,
			sizeStr,
		)
	}

	w.Flush()

	fmt.Println()
	fmt.Printf("Summary: %d/%d assets cached (%s total)\n", cachedCount, len(statuses), formatSize(totalSize))

	return nil
}

func runAssetsClean(_ *cobra.Command, _ []string) error {
	downloader, cacheDir := getAssetsConfig()

	// Check if cache exists
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		fmt.Printf("Cache directory does not exist: %s\n", cacheDir)
		return nil
	}

	if err := downloader.Clean(); err != nil {
		return fmt.Errorf("failed to clean cache: %w", err)
	}

	fmt.Printf("Removed assets cache: %s\n", cacheDir)
	return nil
}

// formatSize formats a byte count as a human-readable string.
func formatSize(bytes int64) string {
	const (
		kb = 1024
		mb = kb * 1024
	)

	switch {
	case bytes >= mb:
		return fmt.Sprintf("%.1f MB", float64(bytes)/mb)
	case bytes >= kb:
		return fmt.Sprintf("%.1f KB", float64(bytes)/kb)
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}
