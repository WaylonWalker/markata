package plugins

import (
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/models"
)

func TestReadingTimeTransformPostWritesOnlyReadingTimeFields(t *testing.T) {
	plugin := NewReadingTimePlugin()
	plugin.SetWordsPerMinute(100)
	post := models.NewPost("reading-time.md")
	post.Content = strings.TrimSpace(strings.Repeat("word ", 101))
	post.Set("sentinel", "keep")

	if err := plugin.TransformPost(post); err != nil {
		t.Fatalf("TransformPost() = %v", err)
	}
	assertReadingTimeFields(t, post, 101, 2, "2 min read")
	if got := post.Get("sentinel"); got != "keep" {
		t.Fatalf("unrelated post field changed: sentinel = %v", got)
	}
}

func TestReadingTimeTransformPostSkipsIneligiblePosts(t *testing.T) {
	plugin := NewReadingTimePlugin()
	tests := []*models.Post{
		nil,
		models.NewPost("empty.md"),
		func() *models.Post {
			post := models.NewPost("skip.md")
			post.Content = "should not be measured"
			post.Skip = true
			return post
		}(),
	}

	for _, post := range tests {
		if err := plugin.TransformPost(post); err != nil {
			t.Fatalf("TransformPost(%v) = %v", post, err)
		}
		if post != nil && GetReadingTime(post) != nil {
			t.Fatalf("ineligible post received reading-time fields: %+v", post.Extra)
		}
	}
}

func TestReadingTimePostEligibilityMatchesLegacyTransformFilter(t *testing.T) {
	eligible := models.NewPost("eligible.md")
	eligible.Content = "one two three"
	empty := models.NewPost("empty.md")
	skipped := models.NewPost("skip.md")
	skipped.Content = "one two three"
	skipped.Skip = true

	if !readingTimePostEligible(eligible) {
		t.Fatal("eligible post was rejected")
	}
	if readingTimePostEligible(empty) {
		t.Fatal("empty post was accepted")
	}
	if readingTimePostEligible(skipped) {
		t.Fatal("skipped post was accepted")
	}
	if readingTimePostEligible(nil) {
		t.Fatal("nil post was accepted")
	}
}
