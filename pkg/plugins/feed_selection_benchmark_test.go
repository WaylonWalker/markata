package plugins

import (
	"fmt"
	"testing"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

func BenchmarkFeedSelectionComplete(b *testing.B) {
	for _, size := range []struct {
		posts, feeds int
	}{{64, 8}, {512, 418}, {3911, 418}} {
		b.Run(fmt.Sprintf("%dx%d", size.posts, size.feeds), func(b *testing.B) {
			posts := make([]*models.Post, size.posts)
			feeds := make([]string, size.feeds)
			for i := range posts {
				posts[i] = &models.Post{Path: fmt.Sprintf("posts/%d.md", i), Published: i%11 != 0, Private: i%13 == 0, Draft: i%17 == 0, Skip: i%19 == 0}
			}
			for i := range feeds {
				feeds[i] = fmt.Sprintf("feed-%d", i)
			}
			matched := posts[:len(posts)/8]
			selected := matched[2:]
			b.ReportAllocs()
			b.ResetTimer()
			for b.Loop() {
				m := lifecycle.NewManager()
				m.SetPosts(posts)
				benchmarkRecordSelections(m, feeds, posts, matched, selected)
			}
		})
		b.Run(fmt.Sprintf("Updates/%dx%d", size.posts, size.feeds), func(b *testing.B) {
			posts := make([]*models.Post, size.posts)
			feeds := make([]string, size.feeds)
			for i := range posts {
				posts[i] = &models.Post{Path: fmt.Sprintf("posts/%d.md", i), Published: i%11 != 0, Private: i%13 == 0, Draft: i%17 == 0, Skip: i%19 == 0}
			}
			for i := range feeds {
				feeds[i] = fmt.Sprintf("feed-%d", i)
			}
			matched, selected := posts[:len(posts)/8], posts[2:len(posts)/8]
			m := lifecycle.NewManager()
			m.SetPosts(posts)
			benchmarkRecordSelections(m, feeds, posts, matched, selected)
			b.ReportAllocs()
			b.ResetTimer()
			for b.Loop() {
				benchmarkRecordSelections(m, feeds, posts, matched, selected)
			}
		})
	}
}

func benchmarkRecordSelections(m *lifecycle.Manager, feeds []string, posts, matched, selected []*models.Post) {
	recorder := newFeedSelectionRecorder(m.ContentLedger(), posts)
	for _, feed := range feeds {
		recorder.record(feed, recorder.sources, matched, selected, "tags contains 'test'", false, 2, 0)
	}
}
