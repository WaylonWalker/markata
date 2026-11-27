package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/models"
)

func TestSelectBlogrollFeed(t *testing.T) {
	feeds := []models.ExternalFeedConfig{
		{URL: "https://one.test/feed", Handle: "Engineering"},
		{URL: "https://two.test/feed", Handle: "engineering"},
	}
	for _, selector := range []string{"", "engineer", "engineering"} {
		if _, err := selectBlogrollFeed(feeds, selector); err == nil {
			t.Errorf("expected missing or ambiguous selector error for %q", selector)
		}
	}
	got, err := selectBlogrollFeed(feeds, "https://one.test/feed")
	if err != nil || got.URL != feeds[0].URL {
		t.Fatalf("exact URL: %v, %v", got.URL, err)
	}
	got, err = selectBlogrollFeed(feeds[:1], "ENGINEERING")
	if err != nil || got.URL != feeds[0].URL {
		t.Fatalf("exact handle: %v, %v", got.URL, err)
	}
}

func TestBlogrollActiveChangePreservesOtherFields(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	path := filepath.Join(dir, "markata-go.toml")
	initial := "[markata-go.blogroll]\nenabled = true\n\n" +
		"[[markata-go.blogroll.feeds]]\nurl = \"https://one.test/feed\"\ntitle = \"One\"\nactive = true\n# Keep comment\n" +
		"tags = [\"go\", \"feeds\"]\nmax_entries = 3\n\n" +
		"[[markata-go.blogroll.feeds]]\nurl = \"https://two.test/feed\"\nactive = true\n"
	if err := os.WriteFile(path, []byte(initial), 0o600); err != nil {
		t.Fatal(err)
	}
	oldConfig, oldMerge, oldRebuild := cfgFile, mergeConfigFiles, serveRequestFullRebuild
	t.Cleanup(func() {
		cfgFile, mergeConfigFiles, serveRequestFullRebuild = oldConfig, oldMerge, oldRebuild
	})
	cfgFile, mergeConfigFiles, serveRequestFullRebuild = path, nil, nil
	state, err := loadSettingsState()
	if err != nil {
		t.Fatal(err)
	}
	req := subscriptionRequest{Action: subscriptionActionActive, URL: "https://one.test/feed", Active: false}
	if err := changeSubscriptionWithOptions(state, req, true); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil || string(data) != initial {
		t.Fatalf("dry run changed config: %v, %s", err, data)
	}
	if err := changeSubscriptionWithOptions(state, req, false); err != nil {
		t.Fatal(err)
	}
	data, err = os.ReadFile(path)
	want := strings.Replace(initial, "active = true", "active = false", 1)
	if err != nil || string(data) != want {
		t.Fatalf("unexpected config changes: %v\n%s", err, data)
	}
}
