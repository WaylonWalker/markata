package cmd

import (
	"net/url"
	"strings"

	"github.com/WaylonWalker/markata/pkg/models"
)

// blogrollMentionResolver deliberately uses only configured local data.
// Linting a site must never require a social-media API or live network lookup.
type blogrollMentionResolver struct {
	handles map[string]struct{}
}

func newBlogrollMentionResolver(cfg models.BlogrollConfig) *blogrollMentionResolver {
	if !cfg.Enabled {
		return nil
	}

	r := &blogrollMentionResolver{handles: make(map[string]struct{})}
	add := func(handle string) {
		handle = strings.ToLower(strings.TrimLeft(strings.TrimSpace(handle), "@~"))
		if handle != "" {
			r.handles[handle] = struct{}{}
		}
	}

	for i := range cfg.Feeds {
		feed := &cfg.Feeds[i]
		if !feed.IsActive() {
			continue
		}
		siteURL := feed.SiteURL
		if siteURL == "" {
			siteURL = feed.URL
		}
		parsed, err := url.Parse(siteURL)
		if err != nil || parsed.Hostname() == "" {
			continue
		}

		host := strings.ToLower(parsed.Hostname())
		add(feed.Handle)
		if feed.Handle == "" {
			bare := strings.TrimPrefix(strings.TrimPrefix(host, "www."), "blog.")
			add(strings.SplitN(bare, ".", 2)[0])
		}
		add(host)
		for _, alias := range feed.Aliases {
			add(alias)
		}
	}
	return r
}

func (r *blogrollMentionResolver) ResolveHandle(handle string) bool {
	if r == nil {
		return false
	}
	_, ok := r.handles[strings.ToLower(handle)]
	return ok
}

// The CLI lint command historically does not check wikilinks; only mentions
// are being enabled here. Returning true preserves the existing behavior.
func (r *blogrollMentionResolver) ResolveSlug(string) bool {
	return true
}
