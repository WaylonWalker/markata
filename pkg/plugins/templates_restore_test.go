package plugins

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/WaylonWalker/markata/pkg/buildcache"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

func TestRestoreCachedFullHTMLOrderAndBytes(t *testing.T) {
	for _, concurrency := range []int{1, 2, 4, 8, 16} {
		t.Run(fmt.Sprint(concurrency), func(t *testing.T) {
			m := lifecycle.NewManager()
			m.SetConcurrency(concurrency)
			posts := []*models.Post{
				{Path: "hit", HTML: "OLD"}, {Path: "missing", HTML: "OLD"},
				{Path: "empty", HTML: "OLD"}, {Path: "other", HTML: "OLD"},
			}
			posts = append(posts, posts[0]) // repeated pointers remain race-safe
			expected := "FULL:\x00\xff<exact>"
			get := func(path string) string {
				if path == "hit" || path == "other" {
					return expected
				}
				return ""
			}
			restored, missing, err := restoreCachedFullHTML(m, posts, get, concurrency)
			if err != nil || restored != 3 || !reflect.DeepEqual(missing, posts[1:3]) {
				t.Fatalf("restore = %d, %v, %v", restored, missing, err)
			}
			for _, post := range posts {
				want := expected
				if post.Path == "missing" || post.Path == "empty" {
					want = "OLD"
				}
				if post.HTML != want {
					t.Fatalf("%s HTML = %q, want %q", post.Path, post.HTML, want)
				}
			}
			if m.Concurrency() != concurrency {
				t.Fatal("restoration changed manager concurrency")
			}
		})
	}
}

func TestRestoreCachedFullHTMLEmptyAndTiny(t *testing.T) {
	m := lifecycle.NewManager()
	restored, missing, err := restoreCachedFullHTML(m, nil, func(string) string {
		t.Fatal("empty restore called getter")
		return ""
	}, 16)
	if err != nil || restored != 0 || len(missing) != 0 {
		t.Fatalf("empty restore = %d, %v, %v", restored, missing, err)
	}
	post := &models.Post{Path: "one"}
	restored, missing, err = restoreCachedFullHTML(m, []*models.Post{post}, func(string) string { return "ONE" }, 16)
	if err != nil || restored != 1 || len(missing) != 0 || post.HTML != "ONE" {
		t.Fatalf("tiny restore = %d, %v, %v, %q", restored, missing, err, post.HTML)
	}
}

func TestRestoreCachedFullHTMLBoundAndJoin(t *testing.T) {
	for _, tc := range []struct{ concurrency, limit, workers int }{
		{16, 4, 4}, {2, 16, 2}, {1, 16, 1}, {16, 1, 1},
	} {
		t.Run(fmt.Sprintf("%d/%d", tc.concurrency, tc.limit), func(t *testing.T) {
			m := lifecycle.NewManager()
			m.SetConcurrency(tc.concurrency)
			posts := make([]*models.Post, 32)
			for i := range posts {
				posts[i] = &models.Post{Path: fmt.Sprint(i), HTML: "OLD"}
			}
			entered := make(chan struct{}, len(posts))
			release := make(chan struct{})
			defer close(release)
			done := make(chan error, 1)
			var active, peak, calls atomic.Int64
			go func() {
				_, _, err := restoreCachedFullHTML(m, posts, func(string) string {
					n := active.Add(1)
					for old := peak.Load(); n > old; old = peak.Load() {
						if peak.CompareAndSwap(old, n) {
							break
						}
					}
					calls.Add(1)
					entered <- struct{}{}
					<-release
					active.Add(-1)
					return "FULL"
				}, tc.limit)
				done <- err
			}()
			for i := 0; i < tc.workers; i++ {
				select {
				case <-entered:
				case <-time.After(5 * time.Second):
					t.Fatal("workers did not start")
				}
			}
			select {
			case err := <-done:
				t.Fatalf("restore returned before join: %v", err)
			default:
			}
			// Release without closing here so the deferred close remains safe.
			for i := 0; i < len(posts); i++ {
				release <- struct{}{}
			}
			if err := <-done; err != nil {
				t.Fatal(err)
			}
			if peak.Load() > int64(tc.workers) || calls.Load() != int64(len(posts)) {
				t.Fatalf("peak=%d calls=%d", peak.Load(), calls.Load())
			}
			for _, post := range posts {
				if post.HTML != "FULL" {
					t.Fatal("joined restore left an incomplete page")
				}
			}
		})
	}
}

func TestTemplatesRestorePeerVisibilityAndOutputRepair(t *testing.T) {
	for _, concurrency := range []int{1, 4, 16} {
		t.Run(fmt.Sprint(concurrency), func(t *testing.T) {
			cache := buildcache.New(t.TempDir())
			miss := headingMigrationPost("miss.md", "miss", "")
			hit := headingMigrationPost("hit.md", "hit", "")
			m, p := headingMigrationManager(t, cache, miss, hit)
			m.SetConcurrency(concurrency)
			m.Config().OutputDir = filepath.Join(t.TempDir(), "output")
			if err := p.Render(m); err != nil {
				t.Fatal(err)
			}
			cache.MarkRebuilt(hit.Path, hit.InputHash, "", hit.Template)
			expected := "<main>EXACT PEER</main>"
			if err := cache.CacheFullHTML(hit.Path, expected); err != nil {
				t.Fatal(err)
			}
			// Removing the disk file must not invalidate the existing memory hit.
			if err := os.Remove(cache.Posts[hit.Path].FullHTMLPath); err != nil {
				t.Fatal(err)
			}
			miss.Template = "peer.html"
			dir, ok := m.Config().Extra["templates_dir"].(string)
			if !ok {
				t.Fatal("fixture template directory missing")
			}
			if err := os.WriteFile(filepath.Join(dir, "peer.html"),
				[]byte(`{% for peer in core.Posts() %}{% if peer.Slug == "hit" %}{{ peer.HTML | safe }}{% endif %}{% endfor %}`), 0o600); err != nil {
				t.Fatal(err)
			}
			hit.HTML = "STALE"
			if err := p.Render(m); err != nil {
				t.Fatal(err)
			}
			stats := m.ContentDiagnostics().TemplateCache
			assertTemplateCacheReconciles(t, stats, 2)
			if stats.Restored != 1 || stats.RenderSucceeded != 1 || hit.HTML != expected || miss.HTML != expected {
				t.Fatalf("peer not restored before rendering: %+v hit=%q miss=%q", stats, hit.HTML, miss.HTML)
			}
			publisher := NewPublishHTMLPlugin()
			if err := publisher.Write(m); err != nil {
				t.Fatal(err)
			}
			if err := os.RemoveAll(m.Config().OutputDir); err != nil {
				t.Fatal(err)
			}
			cache.ResetStats()
			hit.HTML = ""
			if err := p.Render(m); err != nil {
				t.Fatal(err)
			}
			assertTemplateCacheReconciles(t, m.ContentDiagnostics().TemplateCache, 2)
			if err := publisher.Write(m); err != nil {
				t.Fatal(err)
			}
			output, err := os.ReadFile(filepath.Join(m.Config().OutputDir, "hit", "index.html"))
			if err != nil || string(output) != expected {
				t.Fatalf("output repair = %q, %v", output, err)
			}
		})
	}
}

func TestTemplatesRestoreFallbackOrderAndCounters(t *testing.T) {
	for _, concurrency := range []int{1, 4, 16} {
		t.Run(fmt.Sprint(concurrency), func(t *testing.T) {
			cache := buildcache.New(t.TempDir())
			firstHit := headingMigrationPost("first.md", "first", "")
			missing := headingMigrationPost("missing.md", "missing", "")
			classifiedMiss := headingMigrationPost("classified.md", "classified", "")
			empty := headingMigrationPost("empty.md", "empty", "")
			m, p := headingMigrationManager(t, cache, firstHit, missing, classifiedMiss, empty)
			m.SetConcurrency(concurrency)
			if err := p.Render(m); err != nil {
				t.Fatal(err)
			}
			cache.RemoveStale(nil)
			for _, post := range []*models.Post{firstHit, missing, empty} {
				cache.MarkRebuilt(post.Path, post.InputHash, "", post.Template)
				post.HTML = "OLD"
			}
			if err := cache.CacheFullHTML(firstHit.Path, "FULL"); err != nil {
				t.Fatal(err)
			}
			if err := cache.CacheFullHTML(empty.Path, ""); err != nil {
				t.Fatal(err)
			}
			dir, ok := m.Config().Extra["templates_dir"].(string)
			if !ok {
				t.Fatal("fixture template directory missing")
			}
			// All misses fail, so the pool's first error exposes selection order,
			// independent of concurrent completion order.
			if err := os.WriteFile(filepath.Join(dir, "post.html"),
				[]byte(`{% include "absent.html" %}`), 0o600); err != nil {
				t.Fatal(err)
			}
			p.engine.ClearCache()
			err := p.Render(m)
			if err == nil || !strings.Contains(err.Error(), "processing classified.md:") {
				t.Fatalf("classification miss must precede full-page fallbacks: %v", err)
			}
			stats := m.ContentDiagnostics().TemplateCache
			assertTemplateCacheReconciles(t, stats, 4)
			if stats.Cacheable != 3 || stats.Restored != 1 || stats.RenderFailed != 3 ||
				stats.MissReasons.FullHTMLUnavailable != 2 || firstHit.HTML != "FULL" ||
				missing.HTML != "OLD" || empty.HTML != "OLD" {
				t.Fatalf("mixed restore counters or failure HTML changed: %+v", stats)
			}
		})
	}
}

func TestTemplatesRestoreZeroRenderCounters(t *testing.T) {
	cache := buildcache.New(t.TempDir())
	posts := []*models.Post{
		headingMigrationPost("one.md", "one", ""),
		headingMigrationPost("two.md", "two", ""),
		headingMigrationPost("three.md", "three", ""),
	}
	m, p := headingMigrationManager(t, cache, posts...)
	if err := p.Render(m); err != nil {
		t.Fatal(err)
	}
	for _, post := range posts {
		cache.MarkRebuilt(post.Path, post.InputHash, "", post.Template)
	}
	for i := 0; i < 3; i++ {
		for _, post := range posts {
			post.HTML = ""
		}
		if err := p.Render(m); err != nil {
			t.Fatal(err)
		}
		stats := m.ContentDiagnostics().TemplateCache
		assertTemplateCacheReconciles(t, stats, len(posts))
		if stats.Restored != 3 || stats.RenderRequired != 0 || stats.RenderSucceeded != 0 {
			t.Fatalf("zero-render counters changed: %+v", stats)
		}
		for _, post := range posts {
			if post.HTML != "PAGE:" {
				t.Fatalf("zero-render invocation lost full page: %q", post.HTML)
			}
		}
	}
}

// BenchmarkTemplatesFullHTMLRestore measures real getter disk fallback and
// memory hits separately. Disk samples retain the OS page cache but start with
// a fresh build-cache instance, never an already-populated fullHTMLMemory map.
func BenchmarkTemplatesFullHTMLRestore(b *testing.B) {
	const count = 256
	body := strings.Repeat("<p>fixture full page</p>\n", 11000)
	dir := b.TempDir()
	posts := make([]*models.Post, count)
	metadata := make(map[string]*buildcache.PostCache, count)
	for i := range posts {
		path := fmt.Sprintf("post-%d.md", i)
		htmlPath := filepath.Join(dir, fmt.Sprintf("%d.html", i))
		if err := os.WriteFile(htmlPath, []byte(body), 0o600); err != nil {
			b.Fatal(err)
		}
		posts[i] = &models.Post{Path: path}
		metadata[path] = &buildcache.PostCache{FullHTMLPath: htmlPath}
	}
	for _, mode := range []string{"disk", "memory"} {
		for _, limit := range []int{1, 2, 4, 8, 16} {
			b.Run(fmt.Sprintf("%s/limit%d", mode, limit), func(b *testing.B) {
				m := lifecycle.NewManager()
				m.SetConcurrency(16)
				b.ReportAllocs()
				b.SetBytes(int64(count * len(body)))
				for i := 0; i < b.N; i++ {
					b.StopTimer()
					cache := buildcache.New(dir)
					cache.Posts = metadata
					for _, post := range posts {
						post.HTML = ""
						if mode == "memory" {
							cache.GetCachedFullHTML(post.Path)
						}
					}
					b.StartTimer()
					restored, missing, err := restoreCachedFullHTML(m, posts, cache.GetCachedFullHTML, limit)
					if err != nil || restored != count || len(missing) != 0 {
						b.Fatalf("restore = %d, %d, %v", restored, len(missing), err)
					}
				}
			})
		}
	}
}
