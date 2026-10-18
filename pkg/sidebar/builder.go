// Package sidebar provides functions for building sidebar navigation from various sources.
//
// # Overview
//
// The sidebar builder supports three main modes:
//   - Path-based sidebars: Different sidebars for different URL paths
//   - Feed-linked sidebars: Auto-generated from feed posts
//   - Multi-feed sidebars: Combined sidebars from multiple feeds
//
// # Usage
//
//	builder := sidebar.NewBuilder(config, feeds, posts)
//	items, title := builder.ResolveForPost(post)
package sidebar

import (
	"path/filepath"
	"sort"
	"strings"

	"github.com/WaylonWalker/markata/pkg/models"
)

// Builder constructs sidebar navigation from various sources.
type Builder struct {
	config *models.Config
	feeds  map[string]*models.FeedConfig
	posts  []*models.Post
}

// NewBuilder creates a new sidebar builder.
func NewBuilder(config *models.Config, feeds map[string]*models.FeedConfig, posts []*models.Post) *Builder {
	return &Builder{
		config: config,
		feeds:  feeds,
		posts:  posts,
	}
}

// ResolveForPost returns the sidebar items and title for a specific post.
// It resolves path-specific, multi-feed, or default sidebar configurations.
func (b *Builder) ResolveForPost(post *models.Post) (items []models.SidebarNavItem, title string) {
	if b.config == nil {
		return nil, ""
	}

	sidebar := &b.config.Sidebar

	// Check for path-specific sidebar
	if pathConfig, found := sidebar.ResolveForPath(post.Href); found {
		items := b.buildFromPathConfig(pathConfig)
		return items, pathConfig.Title
	}

	// Check for multi-feed mode
	if sidebar.IsMultiFeed() {
		items := b.BuildMultiFeed(sidebar.Feeds, sidebar.FeedSections)
		return items, sidebar.Title
	}

	// Check for auto-generation
	if sidebar.AutoGenerate != nil {
		items := b.BuildFromDirectory(sidebar.AutoGenerate)
		return items, sidebar.Title
	}

	// Return default nav
	return sidebar.Nav, sidebar.Title
}

// buildFromPathConfig builds sidebar items from a path-specific configuration.
func (b *Builder) buildFromPathConfig(pathConfig *models.PathSidebarConfig) []models.SidebarNavItem {
	// Manual items take precedence
	if len(pathConfig.Items) > 0 {
		return pathConfig.Items
	}

	// Try feed-linked sidebar
	if pathConfig.Feed != "" {
		if feed, ok := b.feeds[pathConfig.Feed]; ok {
			return b.BuildFromFeed(feed)
		}
	}

	// Try auto-generation
	if pathConfig.AutoGenerate != nil {
		return b.BuildFromDirectory(pathConfig.AutoGenerate)
	}

	return nil
}

// BuildFromFeed generates sidebar items from a feed's posts.
func (b *Builder) BuildFromFeed(feed *models.FeedConfig) []models.SidebarNavItem {
	if feed == nil || len(feed.Posts) == 0 {
		return nil
	}

	// If grouping is enabled, build grouped structure
	if feed.SidebarGroupBy != "" {
		return b.buildGroupedFromFeed(feed)
	}

	// Build flat list
	items := make([]models.SidebarNavItem, 0, len(feed.Posts))
	for _, post := range feed.Posts {
		items = append(items, postToNavItem(post))
	}

	return items
}

// buildGroupedFromFeed builds a grouped sidebar from a feed using a frontmatter field.
func (b *Builder) buildGroupedFromFeed(feed *models.FeedConfig) []models.SidebarNavItem {
	groups := make(map[string][]*models.Post)
	var ungrouped []*models.Post

	for _, post := range feed.Posts {
		groupValue := getExtraString(post, feed.SidebarGroupBy)
		if groupValue == "" {
			ungrouped = append(ungrouped, post)
		} else {
			groups[groupValue] = append(groups[groupValue], post)
		}
	}

	items := make([]models.SidebarNavItem, 0, len(groups)+len(ungrouped))

	// Add grouped items
	groupNames := make([]string, 0, len(groups))
	for name := range groups {
		groupNames = append(groupNames, name)
	}
	sort.Strings(groupNames)

	for _, groupName := range groupNames {
		children := make([]models.SidebarNavItem, len(groups[groupName]))
		for i, post := range groups[groupName] {
			children[i] = postToNavItem(post)
		}
		items = append(items, models.SidebarNavItem{
			Title:    groupName,
			Children: children,
		})
	}

	// Add ungrouped items at the end
	for _, post := range ungrouped {
		items = append(items, postToNavItem(post))
	}

	return items
}

// BuildFromDirectory generates sidebar items from posts in a directory.
func (b *Builder) BuildFromDirectory(config *models.SidebarAutoGenerate) []models.SidebarNavItem {
	if config == nil || config.Directory == "" {
		return nil
	}

	// Filter posts by directory
	dirPrefix := config.Directory
	if !strings.HasSuffix(dirPrefix, "/") {
		dirPrefix += "/"
	}

	var dirPosts []*models.Post
	for _, post := range b.posts {
		postPath := post.Path
		// Normalize path for comparison
		postPath = strings.TrimPrefix(postPath, "./")

		if strings.HasPrefix(postPath, dirPrefix) || strings.HasPrefix(postPath, config.Directory+"/") {
			// Check exclusions
			if !b.isExcluded(postPath, config.Exclude) {
				dirPosts = append(dirPosts, post)
			}
		}
	}

	// Sort posts
	sortPosts(dirPosts, config.OrderBy, config.IsReverse())

	// Build hierarchy based on subdirectories
	return b.buildHierarchy(dirPosts, dirPrefix, config.MaxDepth)
}

// isExcluded checks if a path matches any of the exclude patterns.
func (b *Builder) isExcluded(path string, patterns []string) bool {
	for _, pattern := range patterns {
		if matched, err := filepath.Match(pattern, filepath.Base(path)); err == nil && matched {
			return true
		}
		if matched, err := filepath.Match(pattern, path); err == nil && matched {
			return true
		}
	}
	return false
}

// hierarchyNode represents a node in the directory hierarchy tree.
type hierarchyNode struct {
	post     *models.Post
	children map[string]*hierarchyNode
	order    int
}

// buildHierarchy creates hierarchical nav items from posts based on directory structure.
func (b *Builder) buildHierarchy(posts []*models.Post, baseDir string, maxDepth int) []models.SidebarNavItem {
	root := &hierarchyNode{children: make(map[string]*hierarchyNode)}

	for _, post := range posts {
		relPath := strings.TrimPrefix(post.Path, baseDir)
		relPath = strings.TrimPrefix(relPath, "./"+baseDir)
		relPath = strings.TrimSuffix(relPath, ".md")
		relPath = strings.TrimSuffix(relPath, "/index")

		parts := strings.Split(relPath, "/")
		if len(parts) == 0 || (len(parts) == 1 && parts[0] == "") {
			continue
		}

		// Respect maxDepth
		if maxDepth > 0 && len(parts) > maxDepth {
			parts = parts[:maxDepth]
		}

		current := root
		for i, part := range parts {
			if part == "" {
				continue
			}
			if current.children[part] == nil {
				current.children[part] = &hierarchyNode{
					children: make(map[string]*hierarchyNode),
					order:    getNavOrder(post),
				}
			}
			// If this is the last part, assign the post
			if i == len(parts)-1 {
				current.children[part].post = post
			}
			current = current.children[part]
		}
	}

	return nodeToItems(root)
}

// nodeToItems converts a tree node to sidebar nav items.
func nodeToItems(n *hierarchyNode) []models.SidebarNavItem {
	if n == nil || len(n.children) == 0 {
		return nil
	}

	// Sort children by order then by name
	type childEntry struct {
		name string
		node *hierarchyNode
	}
	entries := make([]childEntry, 0, len(n.children))
	for name, child := range n.children {
		entries = append(entries, childEntry{name, child})
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].node.order != entries[j].node.order {
			return entries[i].node.order < entries[j].node.order
		}
		return entries[i].name < entries[j].name
	})

	items := make([]models.SidebarNavItem, 0, len(entries))
	for _, entry := range entries {
		item := models.SidebarNavItem{}

		if entry.node.post != nil {
			item = postToNavItem(entry.node.post)
		} else {
			// Directory without an index post
			item.Title = titleCase(entry.name)
		}

		// Add children
		if len(entry.node.children) > 0 {
			item.Children = nodeToItems(entry.node)
		}

		items = append(items, item)
	}

	return items
}

// BuildMultiFeed generates a multi-feed sidebar with collapsible sections.
func (b *Builder) BuildMultiFeed(feedSlugs []string, sections []models.MultiFeedSection) []models.SidebarNavItem {
	var items []models.SidebarNavItem

	// Use detailed sections if provided
	if len(sections) > 0 {
		for _, section := range sections {
			if feed, ok := b.feeds[section.Feed]; ok {
				sectionItem := b.buildFeedSection(feed, &section)
				items = append(items, sectionItem)
			}
		}
		return items
	}

	// Otherwise use feed slugs with default settings
	for _, slug := range feedSlugs {
		if feed, ok := b.feeds[slug]; ok {
			sectionItem := b.buildFeedSection(feed, nil)
			items = append(items, sectionItem)
		}
	}

	return items
}

// buildFeedSection creates a sidebar section for a single feed.
func (b *Builder) buildFeedSection(feed *models.FeedConfig, section *models.MultiFeedSection) models.SidebarNavItem {
	title := feed.GetSidebarTitle()
	if section != nil && section.Title != "" {
		title = section.Title
	}

	feedItems := b.BuildFromFeed(feed)

	// Apply max items limit
	if section != nil && section.MaxItems > 0 && len(feedItems) > section.MaxItems {
		feedItems = feedItems[:section.MaxItems]
	}

	return models.SidebarNavItem{
		Title:    title,
		Children: feedItems,
	}
}

// BuildFromFeeds builds sidebar items from all feeds that have Sidebar enabled.
// Feeds are sorted by SidebarOrder.
func (b *Builder) BuildFromFeeds() []models.SidebarNavItem {
	// Collect feeds with sidebar enabled
	var sidebarFeeds []*models.FeedConfig
	for _, feed := range b.feeds {
		if feed.Sidebar == nil || *feed.Sidebar {
			sidebarFeeds = append(sidebarFeeds, feed)
		}
	}

	// Sort by SidebarOrder
	sort.Slice(sidebarFeeds, func(i, j int) bool {
		return sidebarFeeds[i].SidebarOrder < sidebarFeeds[j].SidebarOrder
	})

	// Build items for each feed
	items := make([]models.SidebarNavItem, 0, len(sidebarFeeds))
	for _, feed := range sidebarFeeds {
		sectionItem := b.buildFeedSection(feed, nil)
		items = append(items, sectionItem)
	}

	return items
}

// Helper functions

// postToNavItem converts a Post to a SidebarNavItem.
func postToNavItem(post *models.Post) models.SidebarNavItem {
	title := post.Slug
	if post.PlainTitle() != "" {
		title = post.PlainTitle()
	}
	return models.SidebarNavItem{
		Title: title,
		Href:  post.Href,
	}
}

// sortPosts sorts posts by the specified field.
func sortPosts(posts []*models.Post, orderBy string, reverse bool) {
	sort.Slice(posts, func(i, j int) bool {
		var less bool
		switch orderBy {
		case "title":
			ti, tj := posts[i].PlainTitle(), posts[j].PlainTitle()
			less = ti < tj
		case "date":
			switch {
			case posts[i].Date != nil && posts[j].Date != nil:
				less = posts[i].Date.Before(*posts[j].Date)
			case posts[i].Date != nil:
				less = true
			default:
				less = false
			}
		case "nav_order":
			oi := getNavOrder(posts[i])
			oj := getNavOrder(posts[j])
			less = oi < oj
		default: // filename
			less = posts[i].Path < posts[j].Path
		}
		if reverse {
			return !less
		}
		return less
	})
}

// getNavOrder extracts the nav_order from a post's extra fields.
func getNavOrder(post *models.Post) int {
	if post.Extra == nil {
		return 999
	}
	if order, ok := post.Extra["nav_order"].(int); ok {
		return order
	}
	if order, ok := post.Extra["nav_order"].(float64); ok {
		return int(order)
	}
	return 999
}

// getExtraString extracts a string value from a post's extra fields.
func getExtraString(post *models.Post, key string) string {
	if post.Extra == nil {
		return ""
	}
	if val, ok := post.Extra[key].(string); ok {
		return val
	}
	return ""
}

// titleCase converts a slug to title case.
func titleCase(s string) string {
	s = strings.ReplaceAll(s, "-", " ")
	s = strings.ReplaceAll(s, "_", " ")
	words := strings.Fields(s)
	for i, word := range words {
		if word != "" {
			words[i] = strings.ToUpper(word[:1]) + word[1:]
		}
	}
	return strings.Join(words, " ")
}
