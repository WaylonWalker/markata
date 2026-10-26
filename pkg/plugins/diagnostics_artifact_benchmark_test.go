package plugins

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/WaylonWalker/markata/pkg/diagnostics"
)

func diagnosticsArtifactFeedFixture(posts, feeds int) diagnostics.ContentLedgerSnapshot {
	dispositions := make([]diagnostics.ContentFeedDisposition, feeds)
	for i := range dispositions {
		dispositions[i] = diagnostics.ContentFeedDisposition{
			Feed: fmt.Sprintf("feed-%03d", i), Included: i%10 == 0,
			Reasons: []string{diagnostics.ReasonContentFiltered},
		}
	}
	snapshot := diagnostics.ContentLedgerSnapshot{
		Summary: diagnostics.ContentSummary{
			Discovered: posts, Candidates: posts, Loaded: posts, FrontmatterValid: posts,
			Posts: posts, Eligible: posts, Rendered: posts, Emitted: posts, Warnings: posts,
		},
		TemplateCache: &diagnostics.TemplateCacheStats{
			Classified: posts, Cacheable: posts, Restored: posts, NavPreviewReset: true,
		},
		Entries: make([]diagnostics.ContentDisposition, posts),
	}
	for i := range snapshot.Entries {
		path := fmt.Sprintf("posts/%05d.md", i)
		snapshot.Entries[i] = diagnostics.ContentDisposition{
			Path: path, Candidate: true, Loaded: true, FrontmatterPresent: true,
			FrontmatterValid: true, PostCreated: true, Eligible: true, Rendered: true,
			Emitted: true, Disposition: diagnostics.DispositionEmitted,
			Diagnostics: []diagnostics.Issue{{
				File: path, Code: diagnostics.ReasonFrontmatterSuspiciousDelimiter,
				Severity: diagnostics.SeverityWarning,
			}},
			// Deliberately shared caller storage: neither buffer size may mutate it.
			Feeds: dispositions,
		}
	}
	return snapshot
}

func TestWriteDiagnosticsArtifact_BufferSizesParity(t *testing.T) {
	snapshot := diagnosticsArtifactFeedFixture(64, 64)
	info := diagnostics.ArtifactBuildInfo{BuiltAt: time.Unix(1, 0)}
	want, err := diagnostics.MarshalArtifact(snapshot, info)
	if err != nil {
		t.Fatal(err)
	}
	var calls [2]int64
	for i, size := range []int{4 << 10, diagnosticsArtifactBufferSize} {
		destination := filepath.Join(t.TempDir(), diagnostics.DefaultArtifactPath)
		var phases diagnosticsArtifactPhases
		err := writeDiagnosticsArtifactBuffered(destination, func(writer io.Writer) error {
			return diagnostics.WriteArtifact(writer, snapshot, info)
		}, &phases, size)
		if err != nil {
			t.Fatal(err)
		}
		got, err := os.ReadFile(destination)
		if err != nil || !bytes.Equal(got, want) {
			t.Fatalf("buffer %d changed complete bytes (including template_cache): %v", size, err)
		}
		if phases.fileWriteBytes != int64(len(want)) || phases.published != 1 {
			t.Fatalf("buffer %d: incorrect raw-write counts: %#v", size, phases)
		}
		calls[i] = phases.fileWriteCalls
	}
	if calls[1] >= calls[0] {
		t.Fatalf("64 KiB buffer did not reduce raw calls: %v", calls)
	}
}

// Both sizes use the actual same-directory temp/flush/sync/close/replace path.
// This measures publication only; snapshot construction is outside the timer.
func BenchmarkDiagnosticsArtifactPublication_CompleteFeeds(b *testing.B) {
	snapshot := diagnosticsArtifactFeedFixture(512, 418)
	info := diagnostics.ArtifactBuildInfo{BuiltAt: time.Unix(1, 0)}
	for _, size := range []int{4 << 10, diagnosticsArtifactBufferSize} {
		b.Run(fmt.Sprintf("%dKiB", size>>10), func(b *testing.B) {
			destination := filepath.Join(b.TempDir(), diagnostics.DefaultArtifactPath)
			var rawDuration time.Duration
			var rawCalls, rawBytes int64
			b.ReportAllocs()
			b.ResetTimer()
			for b.Loop() {
				var phases diagnosticsArtifactPhases
				err := writeDiagnosticsArtifactBuffered(destination, func(writer io.Writer) error {
					return diagnostics.WriteArtifact(writer, snapshot, info)
				}, &phases, size)
				if err != nil {
					b.Fatal(err)
				}
				rawDuration += phases.fileWrite
				rawCalls += phases.fileWriteCalls
				rawBytes += phases.fileWriteBytes
			}
			b.ReportMetric(float64(rawDuration.Nanoseconds())/float64(b.N), "raw-write-ns/op")
			b.ReportMetric(float64(rawCalls)/float64(b.N), "raw-write-calls/op")
			b.ReportMetric(float64(rawBytes)/float64(b.N), "raw-write-bytes/op")
		})
	}
}
