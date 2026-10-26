package plugins

import "github.com/WaylonWalker/markata/pkg/models"

// BlogrollCacheRefreshProgress reports per-feed refresh progress.
type BlogrollCacheRefreshProgress struct {
	FeedURL   string
	FeedTitle string
	Status    string
	Error     string
}

// BlogrollCacheRefreshResult summarizes a direct reader/blogroll cache refresh.
type BlogrollCacheRefreshResult struct {
	FeedsConfigured int
	FeedsRefreshed  int
	FeedsStale      int
	FeedsFailed     int
	EntriesFetched  int
	CacheDir        string
}

// RefreshBlogrollCache refreshes the external feed cache without running a site
// build. The next build will reuse the refreshed cache entries.
func RefreshBlogrollCache(config models.BlogrollConfig) (*BlogrollCacheRefreshResult, error) {
	return RefreshBlogrollCacheWithProgress(config, nil)
}

// RefreshBlogrollCacheWithProgress refreshes the external feed cache and emits
// per-feed progress updates when a callback is provided.
func RefreshBlogrollCacheWithProgress(config models.BlogrollConfig, onProgress func(BlogrollCacheRefreshProgress)) (*BlogrollCacheRefreshResult, error) {
	plugin := NewBlogrollPlugin()
	feeds, entries, summary, err := plugin.fetchFeedsWithOptions(config, blogrollFetchOptions{
		forceRefresh:   true,
		onFeedComplete: onProgress,
	})
	if err != nil {
		return nil, err
	}

	return &BlogrollCacheRefreshResult{
		FeedsConfigured: len(feeds),
		FeedsRefreshed:  summary.refreshed,
		FeedsStale:      summary.stale,
		FeedsFailed:     summary.failed,
		EntriesFetched:  len(entries),
		CacheDir:        config.CacheDir,
	}, nil
}
