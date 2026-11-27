package blogroll

import (
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/models"
)

func TestSourceSetOwnership(t *testing.T) {
	const firstURL = "https://first.test/feed"
	const lastURL = "https://last.test/feed"
	sources := SourceSet{
		{Path: "base.toml", Raw: map[string]any{
			"blogroll": map[string]any{"feeds": []any{map[string]any{"url": firstURL}}},
		}},
		{Path: "site.toml", Raw: map[string]any{
			"markata-go": map[string]any{"blogroll": map[string]any{"feeds": []map[string]any{{"url": lastURL}}}},
		}},
	}
	owner, err := sources.Owner(lastURL)
	if err != nil || owner.Path != "site.toml" {
		t.Fatalf("owner: %s, %v", owner.Path, err)
	}
	if _, err := sources.Owner(firstURL); err == nil {
		t.Fatal("shadowed source must not be editable")
	}
	target, err := sources.AddTarget(nil)
	if err != nil || target.Path != "site.toml" {
		t.Fatalf("target: %s, %v", target.Path, err)
	}
	sources[1].Raw = map[string]any{"blogroll": map[string]any{"feeds": []any{}}}
	if _, err := sources.Owner(firstURL); err == nil {
		t.Fatal("explicit empty list must shadow lower feeds")
	}
	target, err = sources.AddTarget(nil)
	if err != nil || target.Path != "site.toml" {
		t.Fatalf("empty list target: %s, %v", target.Path, err)
	}
	sources[1].Kind = "override"
	if _, err := sources.AddTarget(nil); err == nil {
		t.Fatal("cannot add to a read-only source")
	}
	sources[1].Kind = ""
	sources[1].Raw = map[string]any{"blogroll": map[string]any{"feeds": "invalid"}}
	if _, err := sources.AddTarget(nil); err == nil {
		t.Fatal("must not fall through a malformed definition")
	}
}

func TestSourceSetFallbackAndInventory(t *testing.T) {
	var sources SourceSet
	target, err := sources.AddTarget(func() (SubscriptionSource, error) {
		return SubscriptionSource{Path: "markata-go.toml"}, nil
	})
	if err != nil || target.Path != "markata-go.toml" {
		t.Fatalf("fallback: %s, %v", target.Path, err)
	}
	const address = "https://example.test/feed"
	feeds := []models.ExternalFeedConfig{{URL: address}}
	sources = SourceSet{{
		Path: "site.toml", DisplayPath: "site config", Raw: map[string]any{
			"blogroll": map[string]any{"feeds": []any{map[string]any{"url": address}}},
		},
	}}
	items := sources.Inventory(feeds)
	if len(items) != 1 || !items[0].Editable || items[0].Source != "site config" {
		t.Fatalf("unexpected inventory: %+v", items)
	}
	sources[0].Kind = "global"
	items = sources.Inventory(feeds)
	if items[0].Editable || !strings.Contains(items[0].Source, "global") {
		t.Fatalf("expected read-only inventory: %+v", items)
	}
}

func TestSelectSubscriptionExact(t *testing.T) {
	feeds := []models.ExternalFeedConfig{
		{URL: "https://first.test/feed", Handle: "Engineering"},
		{URL: "https://second.test/feed", Handle: "engineering"},
	}
	if _, err := SelectSubscription(feeds, "engineering"); err == nil {
		t.Fatal("ambiguous handle must be rejected")
	}
	if _, err := SelectSubscription(feeds, "engineer"); err == nil {
		t.Fatal("substring must not match")
	}
	got, err := SelectSubscription(feeds, feeds[0].URL)
	if err != nil || got.URL != feeds[0].URL {
		t.Fatalf("exact URL: %q, %v", got.URL, err)
	}
	got, err = SelectSubscription(feeds[:1], "ENGINEERING")
	if err != nil || got.URL != feeds[0].URL {
		t.Fatalf("case-insensitive handle: %q, %v", got.URL, err)
	}
}
