package themes

import (
	"os"
	"strings"
	"testing"
)

// Standalone posts may receive the generic page-wide "posts" collection.
// It must not be treated as a resolved feed or series navigation context.
func TestFeedSidebarRequiresResolvedFeedContext(t *testing.T) {
	root, err := os.ReadFile("../../templates/components/feed_sidebar.html")
	if err != nil {
		t.Fatal(err)
	}
	bundled, err := ReadTemplate("components/feed_sidebar.html")
	if err != nil {
		t.Fatal(err)
	}

	for name, source := range map[string]string{
		"root":    string(root),
		"bundled": string(bundled),
	} {
		t.Run(name, func(t *testing.T) {
			if !strings.Contains(source, "{% if sidebar.enabled and (sidebar_posts or feed.posts) %}") {
				t.Fatal("sidebar must require enabled state and a resolved sidebar or feed")
			}
			if strings.Contains(source, "{% elif posts %}") ||
				strings.Contains(source, "or posts)") ||
				strings.Contains(source, "{% for p in posts %}") {
				t.Fatal("generic page-wide posts must not trigger or populate the series sidebar")
			}
			if !strings.Contains(source, "{% if sidebar_posts %}") ||
				!strings.Contains(source, "{% for p in sidebar_posts %}") {
				t.Fatal("resolved post sidebar must remain available")
			}
			if !strings.Contains(source, "{% elif feed.posts %}") ||
				!strings.Contains(source, "{% for p in feed.posts %}") {
				t.Fatal("feed pages must retain their own navigation")
			}
			if !strings.Contains(source, "aria-current=\"page\"") {
				t.Fatal("active series item should be marked as current")
			}
		})
	}
}
