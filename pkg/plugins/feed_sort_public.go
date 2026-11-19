package plugins

import "github.com/WaylonWalker/markata/pkg/models"

// SortPostsForFeed applies the same field comparison and stable tie-breaking
// used by the site's real feed publisher. The editor uses this to keep its
// live membership preview ordering aligned with generated feed pages.
func SortPostsForFeed(posts []*models.Post, field string, reverse bool) {
	sortPosts(posts, field, reverse)
}
