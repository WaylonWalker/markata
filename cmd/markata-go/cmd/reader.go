package cmd

import (
	"fmt"
	"sync"
	"time"

	"github.com/WaylonWalker/markata/pkg/config"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/plugins"
	"github.com/spf13/cobra"
)

var readerUpdateConcurrency int

var readerCmd = &cobra.Command{
	Use:   "reader",
	Short: "Manage reader cache data",
	Long: `Commands for refreshing the external feed data used by the /reader page.

These commands update cached reader data without running a full site build.
The next build will reuse the refreshed cache entries.`,
}

var readerUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Refresh cached reader feed data without building",
	Long: `Fetch the configured blogroll feeds and update the on-disk reader cache
without rendering the site.

This is useful when you want the next build to pick up fresh /reader data without
doing a clean build first.

Example usage:
  markata-go reader update
  markata-go reader update --concurrency 12
  markata-go reader update -c markata-go.toml`,
	RunE: runReaderUpdateCommand,
}

func init() {
	rootCmd.AddCommand(readerCmd)
	readerCmd.AddCommand(readerUpdateCmd)
	readerUpdateCmd.Flags().IntVar(&readerUpdateConcurrency, "concurrency", 0, "override feed refresh concurrency for this run (default: config value, else 5)")
}

func runReaderUpdateCommand(_ *cobra.Command, _ []string) error {
	if readerUpdateConcurrency < 0 {
		return fmt.Errorf("--concurrency must be 0 or greater")
	}

	manager, err := createManager(cfgFile)
	if err != nil {
		return fmt.Errorf("initialization failed: %w", err)
	}

	blogrollConfig, ok := readerBlogrollConfig(manager.Config())
	if !ok || !blogrollConfig.Enabled {
		outln("Blogroll reader is not enabled in configuration.")
		outln("Add [markata-go.blogroll] enabled = true to your config file.")
		return nil
	}

	if len(blogrollConfig.Feeds) == 0 {
		outln("No reader feeds configured.")
		outln("Add [[markata-go.blogroll.feeds]] entries to your config file.")
		return nil
	}

	if blogrollConfig.CacheDir == "" {
		blogrollConfig.CacheDir = config.DefaultConfig().Blogroll.CacheDir
	}
	if readerUpdateConcurrency > 0 {
		blogrollConfig.ConcurrentRequests = readerUpdateConcurrency
	}

	totalFeeds := countActiveReaderFeeds(blogrollConfig)
	concurrency := blogrollConfig.ConcurrentRequests
	if concurrency <= 0 {
		concurrency = config.DefaultConfig().Blogroll.ConcurrentRequests
	}
	reporter := newReaderRefreshReporter(totalFeeds, concurrency)
	defer reporter.Stop()

	result, err := plugins.RefreshBlogrollCacheWithProgress(blogrollConfig, reporter.Report)
	if err != nil {
		return fmt.Errorf("refresh reader cache: %w", err)
	}
	reporter.Finish()

	outlnf(
		"Reader cache updated: %d refreshed, %d stale fallback, %d failed, %d entries",
		result.FeedsRefreshed,
		result.FeedsStale,
		result.FeedsFailed,
		result.EntriesFetched,
	)
	outlnf("Cache directory: %s", result.CacheDir)

	return nil
}

func readerBlogrollConfig(cfg *lifecycle.Config) (models.BlogrollConfig, bool) {
	if cfg == nil || cfg.Extra == nil {
		return models.BlogrollConfig{}, false
	}
	blogrollConfig, ok := cfg.Extra["blogroll"].(models.BlogrollConfig)
	return blogrollConfig, ok
}

func countActiveReaderFeeds(config models.BlogrollConfig) int {
	count := 0
	for i := range config.Feeds {
		if config.Feeds[i].IsActive() {
			count++
		}
	}
	return count
}

type readerRefreshReporter struct {
	total       int
	concurrency int
	completed   int
	lastStatus  string
	lastTitle   string
	started     bool
	interactive bool

	mu   sync.Mutex
	done chan struct{}
	wg   sync.WaitGroup
}

func newReaderRefreshReporter(total, concurrency int) *readerRefreshReporter {
	r := &readerRefreshReporter{
		total:       total,
		concurrency: concurrency,
		interactive: errorOutputIsTerminal(),
		done:        make(chan struct{}),
	}
	if total <= 0 {
		return r
	}
	r.started = true
	if r.interactive {
		r.wg.Add(1)
		go r.runSpinner()
		return r
	}
	errlnf("Refreshing reader cache for %d feed(s) with concurrency %d...", total, concurrency)
	return r
}

func (r *readerRefreshReporter) Report(update plugins.BlogrollCacheRefreshProgress) {
	if !r.started {
		return
	}
	r.mu.Lock()
	r.completed++
	r.lastStatus = update.Status
	r.lastTitle = update.FeedTitle
	completed := r.completed
	total := r.total
	r.mu.Unlock()

	if r.interactive {
		return
	}
	if update.Error != "" && update.Status == "failed" {
		errlnf("  [%d/%d] %s %s (%s)", completed, total, update.Status, update.FeedTitle, update.Error)
		return
	}
	errlnf("  [%d/%d] %s %s", completed, total, update.Status, update.FeedTitle)
}

func (r *readerRefreshReporter) Finish() {
	if !r.started {
		return
	}
	r.mu.Lock()
	completed := r.completed
	total := r.total
	r.mu.Unlock()
	if r.interactive {
		errf("\r\033[2K")
	}
	errlnf("Reader cache refresh complete: %d/%d feeds processed", completed, total)
}

func (r *readerRefreshReporter) Stop() {
	if !r.started {
		return
	}
	close(r.done)
	r.wg.Wait()
	if r.interactive {
		errf("\r\033[2K")
	}
}

func (r *readerRefreshReporter) runSpinner() {
	defer r.wg.Done()
	frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	idx := 0
	for {
		select {
		case <-r.done:
			return
		case <-ticker.C:
			r.mu.Lock()
			completed := r.completed
			total := r.total
			lastStatus := r.lastStatus
			lastTitle := r.lastTitle
			r.mu.Unlock()
			message := fmt.Sprintf("%s Refreshing reader cache %d/%d (concurrency %d)", frames[idx], completed, total, r.concurrency)
			if lastTitle != "" {
				message += fmt.Sprintf(" (%s %s)", lastStatus, lastTitle)
			}
			errf("\r\033[2K%s", message)
			idx = (idx + 1) % len(frames)
		}
	}
}
