package plugins

import (
	"testing"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

// TestBlogrollPlugin_RegistersSyntheticPosts verifies that the BlogrollPlugin
// creates synthetic posts for blogroll and reader pages during Configure stage.
func TestBlogrollPlugin_RegistersSyntheticPosts(t *testing.T) {
	// Create manager
	m := lifecycle.NewManager()

	// Set config with blogroll configuration
	config := &lifecycle.Config{
		Extra: map[string]interface{}{
			"blogroll": models.BlogrollConfig{
				Enabled:       true,
				BlogrollSlug:  "blogroll",
				ReaderSlug:    "reader",
				Feeds:         []models.ExternalFeedConfig{},
				CacheDuration: "1h",
			},
		},
	}
	m.SetConfig(config)

	// Create and run plugin
	plugin := NewBlogrollPlugin()
	if err := plugin.Configure(m); err != nil {
		t.Fatalf("Configure() error = %v", err)
	}

	// Get posts
	posts := m.Posts()

	// Verify synthetic posts were created
	var blogrollPost, readerPost *models.Post
	for _, post := range posts {
		if post.Slug == "blogroll" {
			blogrollPost = post
		}
		if post.Slug == "reader" {
			readerPost = post
		}
	}

	// Check blogroll post
	if blogrollPost == nil {
		t.Fatal("Expected blogroll synthetic post to be created")
	}
	if !blogrollPost.Skip {
		t.Error("Expected blogroll post to have Skip=true")
	}
	if !blogrollPost.Published {
		t.Error("Expected blogroll post to have Published=true")
	}
	if blogrollPost.Title == nil || *blogrollPost.Title != "Blogroll" {
		t.Errorf("Expected blogroll title to be 'Blogroll', got %v", blogrollPost.Title)
	}
	if blogrollPost.Href != "/blogroll/" {
		t.Errorf("Expected blogroll href to be '/blogroll/', got %s", blogrollPost.Href)
	}

	// Check reader post
	if readerPost == nil {
		t.Fatal("Expected reader synthetic post to be created")
	}
	if !readerPost.Skip {
		t.Error("Expected reader post to have Skip=true")
	}
	if !readerPost.Published {
		t.Error("Expected reader post to have Published=true")
	}
	if readerPost.Title == nil || *readerPost.Title != "Reader" {
		t.Errorf("Expected reader title to be 'Reader', got %v", readerPost.Title)
	}
	if readerPost.Href != "/reader/" {
		t.Errorf("Expected reader href to be '/reader/', got %s", readerPost.Href)
	}
}

// TestBlogrollPlugin_CustomSlugs verifies that custom slugs are used for synthetic posts.
func TestBlogrollPlugin_CustomSlugs(t *testing.T) {
	// Create manager
	m := lifecycle.NewManager()

	// Set config with custom slugs
	config := &lifecycle.Config{
		Extra: map[string]interface{}{
			"blogroll": models.BlogrollConfig{
				Enabled:       true,
				BlogrollSlug:  "follows",
				ReaderSlug:    "feed",
				Feeds:         []models.ExternalFeedConfig{},
				CacheDuration: "1h",
			},
		},
	}
	m.SetConfig(config)

	// Create and run plugin
	plugin := NewBlogrollPlugin()
	if err := plugin.Configure(m); err != nil {
		t.Fatalf("Configure() error = %v", err)
	}

	// Get posts
	posts := m.Posts()

	// Verify custom slugs were used
	var followsPost, feedPost *models.Post
	for _, post := range posts {
		if post.Slug == "follows" {
			followsPost = post
		}
		if post.Slug == "feed" {
			feedPost = post
		}
	}

	if followsPost == nil {
		t.Fatal("Expected 'follows' synthetic post to be created")
	}
	if followsPost.Href != "/follows/" {
		t.Errorf("Expected follows href to be '/follows/', got %s", followsPost.Href)
	}

	if feedPost == nil {
		t.Fatal("Expected 'feed' synthetic post to be created")
	}
	if feedPost.Href != "/feed/" {
		t.Errorf("Expected feed href to be '/feed/', got %s", feedPost.Href)
	}
}

// TestPublishFeedsPlugin_RegistersSyntheticPosts verifies that the PublishFeedsPlugin
// creates synthetic posts for each feed during Configure stage.
func TestPublishFeedsPlugin_RegistersSyntheticPosts(t *testing.T) {
	// Create feed configs
	feedConfigs := []models.FeedConfig{
		{
			Slug:        "blog",
			Title:       "My Blog",
			Description: "All my blog posts",
		},
		{
			Slug:        "archive",
			Title:       "Archive",
			Description: "Archived posts",
		},
	}

	// Create manager
	m := lifecycle.NewManager()
	config := &lifecycle.Config{}
	m.SetConfig(config)

	// Set feed configs in cache
	m.Cache().Set("feed_configs", feedConfigs)

	// Create and run plugin
	plugin := NewPublishFeedsPlugin()
	if err := plugin.Configure(m); err != nil {
		t.Fatalf("Configure() error = %v", err)
	}

	// Get posts
	posts := m.Posts()

	// Verify synthetic posts were created
	var blogPost, archivePost *models.Post
	for _, post := range posts {
		if post.Slug == "blog" {
			blogPost = post
		}
		if post.Slug == "archive" {
			archivePost = post
		}
	}

	// Check blog post
	if blogPost == nil {
		t.Fatal("Expected blog synthetic post to be created")
	}
	if !blogPost.Skip {
		t.Error("Expected blog post to have Skip=true")
	}
	if !blogPost.Published {
		t.Error("Expected blog post to have Published=true")
	}
	if blogPost.Title == nil || *blogPost.Title != "My Blog" {
		t.Errorf("Expected blog title to be 'My Blog', got %v", blogPost.Title)
	}
	if blogPost.Href != "/blog/" {
		t.Errorf("Expected blog href to be '/blog/', got %s", blogPost.Href)
	}

	// Check archive post
	if archivePost == nil {
		t.Fatal("Expected archive synthetic post to be created")
	}
	if !archivePost.Skip {
		t.Error("Expected archive post to have Skip=true")
	}
	if archivePost.Title == nil || *archivePost.Title != "Archive" {
		t.Errorf("Expected archive title to be 'Archive', got %v", archivePost.Title)
	}
	if archivePost.Href != "/archive/" {
		t.Errorf("Expected archive href to be '/archive/', got %s", archivePost.Href)
	}
}

// TestPublishFeedsPlugin_NoFeedConfigs verifies that the plugin handles the case
// where no feed configs are available.
func TestPublishFeedsPlugin_NoFeedConfigs(t *testing.T) {
	// Create manager without feed configs
	m := lifecycle.NewManager()
	config := &lifecycle.Config{}
	m.SetConfig(config)

	// Create and run plugin
	plugin := NewPublishFeedsPlugin()
	if err := plugin.Configure(m); err != nil {
		t.Fatalf("Configure() error = %v", err)
	}

	// Verify no posts were added
	posts := m.Posts()
	if len(posts) != 0 {
		t.Errorf("Expected 0 posts, got %d", len(posts))
	}
}
