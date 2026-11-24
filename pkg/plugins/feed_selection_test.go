package plugins

import (
	"reflect"
	"testing"

	"github.com/WaylonWalker/markata/pkg/buildcache"
	"github.com/WaylonWalker/markata/pkg/diagnostics"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

// Independent pre-batch selection oracle: do not use the recorder's helpers.
func legacyFeedSelection(ledger *diagnostics.ContentLedger, feed string, considered, matched, selected []*models.Post, filter string, private bool, offset, limit int) {
	matchedPaths := make(map[string]int)
	for i, post := range matched {
		if post != nil && post.Path != "" {
			if _, ok := matchedPaths[post.Path]; !ok {
				matchedPaths[post.Path] = i
			}
		}
	}
	selectedPaths := make(map[string]bool)
	for _, post := range selected {
		if post != nil && post.Path != "" {
			selectedPaths[post.Path] = true
		}
	}
	for _, post := range considered {
		if post == nil || post.Path == "" {
			continue
		}
		var reasons []string
		if !post.Published {
			reasons = append(reasons, diagnostics.ReasonContentPublishedFalse)
		}
		if post.Skip {
			reasons = append(reasons, diagnostics.ReasonContentSkip)
		}
		if post.Draft {
			reasons = append(reasons, diagnostics.ReasonContentDraft)
		}
		if post.Private && !private {
			reasons = append(reasons, diagnostics.ReasonContentPrivate)
		}
		index, match := matchedPaths[post.Path]
		included := selectedPaths[post.Path]
		if !match {
			included = false
			if filter != "" && len(reasons) == 0 {
				reasons = append(reasons, diagnostics.ReasonContentFiltered)
			}
		} else if !included {
			o, l := max(offset, 0), max(limit, 0)
			if index < o || o >= len(matched) {
				reasons = append(reasons, diagnostics.ReasonFeedOffset)
			} else if l > 0 && index >= o+l {
				reasons = append(reasons, diagnostics.ReasonFeedLimit)
			}
		}
		if (!match || !included) && len(reasons) == 0 {
			reasons = append(reasons, diagnostics.ReasonContentNoOutput)
		}
		ledger.RecordFeed(post.Path, feed, included && len(reasons) == 0, reasons...)
	}
}

func TestFeedSelectionRecorderDifferential(t *testing.T) {
	post := &models.Post{Path: "a.md", Published: true}
	posts := []*models.Post{
		nil, {}, post, post, {Path: "a.md", Draft: true},
		{Path: "./a.md", Published: true, Private: true},
		{Path: "b.md", Published: true, Skip: true, Draft: true, Private: true},
		{Path: "c.md", Published: true}, {Path: "../outside.md"},
		{Path: "/absolute.md", Published: true}, {Path: "unknown.txt"},
		{Path: "unknown.png", Published: true},
	}
	matched := []*models.Post{nil, post, post, posts[6], posts[7], posts[9]}
	for _, private := range []bool{false, true} {
		for _, filter := range []string{"", "tags contains 'x'"} {
			for _, window := range [][2]int{{0, 0}, {2, 1}, {99, 1}, {-1, -1}, {0, 2}} {
				got, want := diagnostics.NewContentLedger(), diagnostics.NewContentLedger()
				recorder := newFeedSelectionRecorder(got, posts)
				for _, feed := range []string{"", "same", "same", "other"} {
					selected := matched
					if window[0] > 0 {
						selected = matched[min(window[0], len(matched)):]
					}
					if window[1] > 0 {
						selected = selected[:min(window[1], len(selected))]
					}
					recorder.record(feed, recorder.sources, matched, selected, filter, private, window[0], window[1])
					legacyFeedSelection(want, feed, posts, matched, selected, filter, private, window[0], window[1])
					if !reflect.DeepEqual(got.Snapshot(), want.Snapshot()) {
						t.Fatalf("private=%v filter=%q window=%v feed=%q", private, filter, window, feed)
					}
				}
				preset := []*models.Post{posts[6], post, posts[4], post}
				recorder.record("series", prepareFeedSources(preset), preset, preset[1:], filter, private, 1, 0)
				legacyFeedSelection(want, "series", preset, preset, preset[1:], filter, private, 1, 0)
				if !reflect.DeepEqual(got.Snapshot(), want.Snapshot()) {
					t.Fatal("preset series differs")
				}
			}
		}
	}
}

func TestFeedCollectorsSuccessiveDifferential(t *testing.T) {
	m := lifecycle.NewManager()
	want := diagnostics.NewContentLedger()
	post := &models.Post{Path: "a.md", Published: true, Tags: []string{"Go", "go"}}
	posts := []*models.Post{post, post, {Path: "a.md", Draft: true, Tags: []string{"go"}},
		{Path: "./a.md", Published: true, Private: true, Tags: []string{"go"}},
		{Path: "b.md", Published: true, Tags: []string{"other"}},
		{Path: "skip.md", Published: true, Skip: true, Tags: []string{"go"}}}
	for iteration := 0; iteration < 3; iteration++ {
		if iteration == 1 {
			post.Private, post.Draft = true, true
			posts[4].Tags = []string{"go"}
		}
		if iteration == 2 {
			post.Private, post.Draft = false, false
			post.Tags = []string{"new"}
			posts = posts[:len(posts)-1]
		}
		m.SetPosts(posts)
		config := lifecycle.NewConfig()
		configured := []models.FeedConfig{
			{Slug: "tags/go", Filter: "tags contains 'go'", IncludePrivate: iteration == 1, Offset: iteration, Limit: 1},
			{Slug: "all", IncludePrivate: true},
			{Slug: "series", Type: models.FeedTypeSeries, Posts: []*models.Post{posts[3], post, posts[2], post}, Offset: 1, Limit: 2},
		}
		auto := AutoFeedsConfig{Tags: AutoFeedTypeConfig{Enabled: true, SlugPrefix: "tags"}}
		config.Extra = map[string]interface{}{
			"feeds": configured, "auto_feeds": auto,
			"models_config": &models.Config{Encryption: models.EncryptionConfig{PrivateTags: map[string]string{"go": "test-key"}}},
		}
		m.SetConfig(config)
		expected := append([]models.FeedConfig(nil), configured...)
		autoPlugin := NewAutoFeedsPlugin()
		expected = append(expected, autoPlugin.generateTagFeeds(posts, auto.Tags, getPrivateTagSlugs(config))...)
		for _, fc := range expected {
			fc.ApplyDefaults(getFeedDefaults(config))
			considered := posts
			var matched []*models.Post
			if fc.Type == models.FeedTypeSeries && len(fc.Posts) > 0 {
				considered = fc.Posts
				matched = filterFeedPagePosts(fc.Posts, fc.IncludesPrivate())
			} else {
				var err error
				matched, err = filterPosts(posts, fc.Filter, fc.IncludesPrivate())
				if err != nil {
					t.Fatal(err)
				}
				sortPosts(matched, "date", true)
			}
			selected := applyFeedLimitOffset(matched, &fc)
			legacyFeedSelection(want, fc.Slug, considered, matched, selected, fc.Filter, fc.IncludesPrivate(), fc.Offset, fc.Limit)
		}
		if err := NewFeedsPlugin().Collect(m); err != nil {
			t.Fatal(err)
		}
		if err := autoPlugin.Collect(m); err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(m.ContentLedger().Snapshot(), want.Snapshot()) {
			t.Fatalf("iteration %d configured/auto snapshots differ", iteration)
		}
	}
}

func TestFeedSelectionSkippedCollectFreshness(t *testing.T) {
	m := lifecycle.NewManager()
	selected := &models.Post{Path: "selected.md", Published: true}
	excluded := &models.Post{Path: "excluded.md", Published: true, Draft: true}
	posts := []*models.Post{selected, excluded}
	cache := buildcache.New(t.TempDir())
	cache.UpdatePostSemanticHashes(selected.Path, computePostFeedItemHash(selected), "", "")
	m.Cache().Set("build_cache", cache)
	want := diagnostics.NewContentLedger()
	plugin := NewFeedsPlugin()
	for iteration := 0; iteration < 2; iteration++ {
		if iteration == 1 {
			excluded.Draft, excluded.Private = false, true
		}
		m.SetPosts(posts)
		fc := models.FeedConfig{Slug: "cached", Filter: "tags contains 'x'", Posts: []*models.Post{selected}}
		m.Config().Extra = map[string]interface{}{
			"feeds_incremental": true,
			"feeds":             []models.FeedConfig{fc},
		}
		if !plugin.shouldSkipFeedCollect(&fc, m) {
			t.Fatal("fixture did not take skipped collection branch")
		}
		legacyFeedSelection(want, fc.Slug, posts, fc.Posts, fc.Posts, fc.Filter, false, 0, 0)
		if err := plugin.Collect(m); err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(m.ContentLedger().Snapshot(), want.Snapshot()) {
			t.Fatalf("skipped iteration %d used stale source flags", iteration)
		}
	}
}

func TestFeedSelectionOneShotFreshness(t *testing.T) {
	m := lifecycle.NewManager()
	post := &models.Post{Path: "post.md", Published: true}
	m.SetPosts([]*models.Post{post, post})
	want := diagnostics.NewContentLedger()
	for _, private := range []bool{true, false} {
		post.Private = private
		recordFeedSelection(m, "", []*models.Post{post}, "", false)
		legacyFeedSelection(want, "", m.Posts(), []*models.Post{post}, []*models.Post{post}, "", false, 0, 0)
		if !reflect.DeepEqual(m.ContentLedger().Snapshot(), want.Snapshot()) {
			t.Fatal("one-shot helper retained prepared flags")
		}
	}
	recordFeedSelection(nil, "", nil, "", false)
	recordFeedSelectionForPosts(nil, "", nil, nil, nil, "", false, 0, 0)
}

func TestFeedSelectionReasonOrderAndWindowPrecedence(t *testing.T) {
	post := &models.Post{Path: "post.md", Skip: true, Draft: true, Private: true}
	recorder := newFeedSelectionRecorder(diagnostics.NewContentLedger(), []*models.Post{post})
	recorder.record("offset", recorder.sources, []*models.Post{nil, post, post}, nil, "", false, 2, 1)
	expected := []string{
		diagnostics.ReasonContentPublishedFalse, diagnostics.ReasonContentSkip,
		diagnostics.ReasonContentDraft, diagnostics.ReasonContentPrivate, diagnostics.ReasonFeedOffset,
	}
	if !reflect.DeepEqual(recorder.observations[0].Reasons, expected) {
		t.Fatalf("reasons = %v, want %v", recorder.observations[0].Reasons, expected)
	}
	recorder.record("limit", recorder.sources, []*models.Post{nil, post, post}, nil, "", true, 0, 1)
	expected = []string{
		diagnostics.ReasonContentPublishedFalse, diagnostics.ReasonContentSkip,
		diagnostics.ReasonContentDraft, diagnostics.ReasonFeedLimit,
	}
	if !reflect.DeepEqual(recorder.observations[0].Reasons, expected) {
		t.Fatalf("reasons = %v, want %v", recorder.observations[0].Reasons, expected)
	}
}
