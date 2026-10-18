package plugins

import (
	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

// restoreCachedFullHTML joins all reads before assigning public HTML or
// collecting unavailable pages. The getter's empty-result fallback policy is
// unchanged. Results own immutable strings, not another copy of page bodies.
func restoreCachedFullHTML(m *lifecycle.Manager, posts []*models.Post, get func(string) string, limit int) (int, []*models.Post, error) {
	if len(posts) == 0 {
		return 0, nil, nil
	}
	if limit <= 1 || m.Concurrency() == 1 || len(posts) == 1 {
		var restored int
		var unavailable []*models.Post
		for _, post := range posts {
			html := get(post.Path)
			if html == "" {
				unavailable = append(unavailable, post)
			} else {
				post.HTML = html
				restored++
			}
		}
		return restored, unavailable, nil
	}
	indices := make(map[*models.Post]int, len(posts))
	unique := make([]*models.Post, 0, len(posts))
	for i, post := range posts {
		// A repeated pointer shares one read/result, avoiding concurrent writes
		// to a result slot while preserving input multiplicity on aggregation.
		if _, exists := indices[post]; !exists {
			indices[post] = i
			unique = append(unique, post)
		}
	}
	results := make([]string, len(posts))
	if err := m.ProcessPostsSliceConcurrentlyWithLimit(unique, limit, func(post *models.Post) error {
		results[indices[post]] = get(post.Path)
		return nil
	}); err != nil {
		return 0, nil, err
	}

	var restored int
	var unavailable []*models.Post
	for _, post := range posts {
		html := results[indices[post]]
		if html == "" {
			unavailable = append(unavailable, post)
			continue
		}
		post.HTML = html
		restored++
	}
	return restored, unavailable, nil
}
