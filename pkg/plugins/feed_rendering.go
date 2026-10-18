package plugins

import "github.com/WaylonWalker/markata/pkg/models"

// splitFeedRenderablePosts partitions posts into:
// 1) pagePosts: posts visible on feed HTML pages
// 2) outputPosts: posts with content suitable for RSS/Atom/JSON/MD/TXT output
func splitFeedRenderablePosts(posts []*models.Post, includePrivate bool) (pagePosts, outputPosts []*models.Post) {
	pagePosts = make([]*models.Post, 0, len(posts))
	outputPosts = make([]*models.Post, 0, len(posts))
	for _, post := range posts {
		if post == nil || post.Skip || post.Draft {
			continue
		}
		if post.Private && !includePrivate {
			continue
		}

		safePost := safeFeedPost(post)
		pagePosts = append(pagePosts, safePost)
		if safePost.Content == "" && safePost.ArticleHTML == "" {
			continue
		}
		outputPosts = append(outputPosts, safePost)
	}

	return pagePosts, outputPosts
}

func filterFeedPagePosts(posts []*models.Post, includePrivate bool) []*models.Post {
	visible, _ := splitFeedRenderablePosts(posts, includePrivate)
	return visible
}

func filterFeedOutputPosts(posts []*models.Post, includePrivate bool) []*models.Post {
	_, renderable := splitFeedRenderablePosts(posts, includePrivate)
	return renderable
}
