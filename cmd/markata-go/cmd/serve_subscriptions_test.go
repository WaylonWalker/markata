package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/models"
)

func TestSubscriptionBlocksPreserveNeighboringTables(t *testing.T) {
	source := "# hello\n[markata-go.blogroll]\nenabled = true\n\n" +
		"[[markata-go.blogroll.feeds]]\nurl = \"https://one.test/rss\"\ntitle = \"One\"\n# keep first comment\n\n" +
		"[[markata-go.blogroll.feeds]]\nurl = \"https://two.test/rss\"\ntitle = \"Two\"\n\n" +
		"[markata-go.theme]\ntext_size = \"medium\"\n"
	blocks := subscriptionBlocks(source)
	if len(blocks) != 2 {
		t.Fatalf("got %d blocks, want 2", len(blocks))
	}
	first, err := uniqueBlock(source, "https://one.test/rss")
	if err != nil {
		t.Fatal(err)
	}
	updated := source[:first.start] + subscriptionField(source[first.start:first.end], "title", "\"Edited\"") + source[first.end:]
	if !strings.Contains(updated, "title = \"Edited\"") || !strings.Contains(updated, "# keep first comment") ||
		!strings.Contains(updated, "url = \"https://two.test/rss\"") ||
		!strings.Contains(updated, "[markata-go.theme]") {
		t.Fatal("update unexpectedly changed another feed, comment, or section")
	}
	second, err := uniqueBlock(updated, "https://two.test/rss")
	if err != nil {
		t.Fatal(err)
	}
	removed := updated[:second.start] + updated[second.end:]
	if strings.Contains(removed, "https://two.test/rss") || !strings.Contains(removed, "https://one.test/rss") {
		t.Fatal("removal did not isolate selected feed")
	}
}

func TestEnableSubscriptions(t *testing.T) {
	for _, input := range []string{
		"[markata-go.blogroll]\nenabled = false\n",
		"[markata-go.blogroll]\n",
		"[markata-go]\ntitle = \"Blog\"\n",
	} {
		got := enableSubscriptions(input)
		if !strings.Contains(got, "enabled = true") || strings.Contains(got, "enabled = false") {
			t.Errorf("failed to enable blogroll for %q: %q", input, got)
		}
	}
}

func TestSubscriptionTargetEffectiveInclude(t *testing.T) {
	root := t.TempDir()
	files := []string{filepath.Join(root, "markata-go.toml"), filepath.Join(root, "base.toml"), filepath.Join(root, "overrides.toml")}
	for _, path := range files {
		if err := os.WriteFile(path, []byte("[markata-go]\n"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	state := &settingsState{
		cfg:     &models.Config{Blogroll: models.BlogrollConfig{Feeds: []models.ExternalFeedConfig{{URL: "https://last.test/feed"}}}},
		sources: files,
		raws: []map[string]any{
			{"markata-go": map[string]any{"blogroll": map[string]any{"enabled": true}}},
			{"markata-go": map[string]any{"blogroll": map[string]any{"feeds": []any{map[string]any{"url": "https://first.test/feed"}}}}},
			{"markata-go": map[string]any{"blogroll": map[string]any{"feeds": []any{map[string]any{"url": "https://last.test/feed"}}}}},
		},
		kinds: map[string]string{},
	}
	target, err := newFeedTarget(state)
	if err != nil || target != files[2] {
		t.Fatalf("new target = %q, %v", target, err)
	}
	target, err = feedOwner(state, "https://last.test/feed")
	if err != nil || target != files[2] {
		t.Fatalf("feed owner = %q, %v", target, err)
	}
	if _, err := feedOwner(state, "https://first.test/feed"); err == nil {
		t.Fatal("should not edit shadowed feed")
	}
	state.kinds[files[2]] = targetKindOverride
	if _, err := newFeedTarget(state); err == nil {
		t.Fatal("should not write into merge override")
	}
	if _, err := feedOwner(state, "https://last.test/feed"); err == nil {
		t.Fatal("should not edit override feed")
	}
}

func TestSubscriptionFieldsRejectBadType(t *testing.T) {
	_, err := subscriptionFields(subscriptionRequest{Type: "broken"})
	if err == nil {
		t.Fatal("unknown reader type accepted")
	}
}
