package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/WaylonWalker/markata/pkg/buildstats"
	"github.com/WaylonWalker/markata/pkg/diagnostics"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
)

// Keep the original encoder and custom MarshalJSON as an independent wire-format
// oracle, not a second implementation of the streaming layout.
func legacyBenchmarkJSON(w io.Writer, result *BuildResult) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(benchmarkJSONOutput{
		Executor: result.Executor, PostsProcessed: result.PostsProcessed,
		FeedsGenerated: result.FeedsGenerated, Duration: result.Duration,
		Warnings: result.Warnings, Benchmark: result.Benchmark,
		Blogroll: result.BlogrollStatus, Content: result.Content,
	})
}

func benchmarkStreamFixture(count, feeds int) *BuildResult {
	const tricky = "<script>&\"\\\n\t\r\x00 caf\u00e9 \u4e16\u754c \U0001f600 \u2028\u2029 \xff"
	resources := buildstats.ResourceBreakdown{
		CPU: time.Second, NetworkWait: 2 * time.Second,
		DiskReadWait: 3 * time.Second, DiskWriteWait: 4 * time.Second, Idle: 5 * time.Second,
	}
	result := &BuildResult{
		PostsProcessed: count, FeedsGenerated: feeds, Duration: 1.23456789,
		Warnings:       []string{tricky, "", "warning"},
		BlogrollStatus: BlogrollStatus{Configured: true, Enabled: true, FeedsConfigured: 17, FeedsFetched: 13},
		Benchmark: buildstats.Summary{
			Total: 15 * time.Second, Resources: resources,
			Spans: []buildstats.SpanTiming{{
				ID: "span-000001", Name: tricky, Status: "ok",
				StartOffset: 123, Duration: 456, Attributes: map[string]string{"operation": tricky},
			}},
			Hotspots: []buildstats.Hotspot{{Stage: "render", Plugin: tricky, Duration: 12345}},
			Requests: []buildstats.RequestTiming{{
				Stage: "collect", Plugin: "blogroll", Method: "GET", Host: "example.com",
				URL: "https://example.com/?q=" + tricky, Status: 503, Error: tricky, Duration: 67890,
			}},
			Stages: []buildstats.StageTiming{{Stage: "write", Duration: 45678, Resources: resources}},
		},
		Content: diagnostics.ContentLedgerSnapshot{
			Summary: diagnostics.ContentSummary{
				Discovered: count + 2, Candidates: count, Loaded: count - 1,
				FrontmatterValid: count - 2, Posts: count - 3, Eligible: count - 4,
				Rendered: count - 5, Emitted: count - 6, Excluded: 6, Warnings: 7, Errors: 8,
			},
			Entries: make([]diagnostics.ContentDisposition, count),
		},
	}
	for i := range result.Content.Entries {
		entry := diagnostics.ContentDisposition{
			// Deliberately descending, non-normalized paths and unsorted reasons.
			Path:      fmt.Sprintf("../posts\\%06d-%s.md", count-i, tricky),
			Candidate: true, Loaded: true, FrontmatterPresent: true, FrontmatterValid: true,
			PostCreated: true, Eligible: i%2 == 0, Rendered: true,
			Emitted: i%3 == 0, Excluded: i%2 != 0, Disposition: diagnostics.DispositionShadow,
			Reasons: []string{diagnostics.ReasonContentPrivate, tricky, diagnostics.ReasonContentDraft, tricky},
			Diagnostics: []diagnostics.Issue{{
				File: "/raw/" + tricky, Range: diagnostics.Range{StartLine: 1, StartCol: 2, EndLine: 3, EndCol: 4},
				Code: diagnostics.ReasonFrontmatterParseError, Severity: diagnostics.SeverityWarning,
				Message: tricky, Fixable: true,
			}},
		}
		for j := 0; j < feeds; j++ {
			entry.Feeds = append(entry.Feeds, diagnostics.ContentFeedDisposition{
				Feed: fmt.Sprintf("feed-%03d-%s", feeds-j, tricky), Included: j%3 == 0,
				Reasons: []string{diagnostics.ReasonFeedLimit, diagnostics.ReasonContentFiltered, diagnostics.ReasonFeedOffset, tricky},
			})
		}
		result.Content.Entries[i] = entry
	}
	return result
}

func TestBenchmarkJSONStreamLegacyBytes(t *testing.T) {
	for _, executor := range []lifecycle.BuildExecutor{"", lifecycle.BuildExecutorLegacy, lifecycle.BuildExecutorDAG} {
		for _, count := range []int{-1, 0, 1, 3} {
			t.Run(fmt.Sprintf("%s/entries=%d", executor, count), func(t *testing.T) {
				makeResult := func() *BuildResult {
					result := benchmarkStreamFixture(max(count, 0), 4)
					result.Executor = executor
					if count < 0 {
						result.Content.Entries = nil
						result.Warnings = nil
					}
					return result
				}
				result, original := makeResult(), makeResult()
				var want, got bytes.Buffer
				if err := legacyBenchmarkJSON(&want, result); err != nil {
					t.Fatal(err)
				}
				if err := writeBenchmarkJSON(&got, result); err != nil {
					t.Fatal(err)
				}
				if !bytes.Equal(got.Bytes(), want.Bytes()) {
					t.Fatalf("stream differs from legacy:\ngot:\n%s\nwant:\n%s", got.Bytes(), want.Bytes())
				}
				if !reflect.DeepEqual(result, original) {
					t.Fatal("serialization mutated the build result")
				}
			})
		}
	}
	t.Run("zero result", func(t *testing.T) {
		var want, got bytes.Buffer
		if err := legacyBenchmarkJSON(&want, &BuildResult{}); err != nil {
			t.Fatal(err)
		}
		if err := writeBenchmarkJSON(&got, &BuildResult{}); err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(got.Bytes(), want.Bytes()) {
			t.Fatal("zero result differs from legacy")
		}
	})
}

func TestBenchmarkSummaryJSONOmitsContentEntries(t *testing.T) {
	result := benchmarkStreamFixture(3, 4)
	result.Executor = lifecycle.BuildExecutorDAG
	result.Content.TemplateCache = &diagnostics.TemplateCacheStats{Classified: 3, Restored: 2, RenderRequired: 1}

	var output bytes.Buffer
	if err := writeBenchmarkSummaryJSON(&output, result); err != nil {
		t.Fatal(err)
	}
	var decoded map[string]json.RawMessage
	if err := json.Unmarshal(output.Bytes(), &decoded); err != nil {
		t.Fatal(err)
	}
	var content map[string]json.RawMessage
	if err := json.Unmarshal(decoded["content"], &content); err != nil {
		t.Fatal(err)
	}
	if _, exists := content["entries"]; exists {
		t.Fatal("compact benchmark contains content entries")
	}
	if _, exists := content["summary"]; !exists {
		t.Fatal("compact benchmark omitted content summary")
	}
	if _, exists := content["template_cache"]; !exists {
		t.Fatal("compact benchmark omitted template cache statistics")
	}
	if string(decoded["executor"]) != `"dag"` {
		t.Fatalf("executor = %s, want dag", decoded["executor"])
	}
	if bytes.Contains(output.Bytes(), []byte("../posts")) {
		t.Fatal("compact benchmark serialized entry data")
	}
}

func TestBenchmarkSummaryJSONRejectsInvalidExecutor(t *testing.T) {
	if err := writeBenchmarkSummaryJSON(io.Discard, &BuildResult{Executor: "invalid"}); err == nil {
		t.Fatal("invalid executor was serialized as a successful compact benchmark")
	}
}

type benchmarkFailureWriter struct {
	err      error
	short    bool
	calls    int
	failAt   int
	closeErr error
	closed   bool
}

func (w *benchmarkFailureWriter) Write(p []byte) (int, error) {
	w.calls++
	if w.calls < w.failAt {
		return len(p), nil
	}
	if w.short {
		return len(p) - 1, nil
	}
	return 0, w.err
}

func (w *benchmarkFailureWriter) Close() error {
	w.closed = true
	return w.closeErr
}

func TestBenchmarkJSONStreamFailures(t *testing.T) {
	sentinel := errors.New("writer failed")
	closeErr := errors.New("close failed")
	for _, test := range []struct {
		name   string
		count  int
		short  bool
		failAt int
	}{
		{"flush error", 0, false, 1},
		{"flush short write", 0, true, 1},
		{"entry write error", 10, false, 1},
		{"entry short write", 10, true, 1},
		{"later write error", 10, false, 3},
	} {
		t.Run(test.name, func(t *testing.T) {
			w := &benchmarkFailureWriter{err: sentinel, short: test.short, failAt: test.failAt, closeErr: closeErr}
			err := writeBenchmarkJSONAndClose(w, benchmarkStreamFixture(test.count, 12))
			want := sentinel
			if test.short {
				want = io.ErrShortWrite
			}
			if !errors.Is(err, want) || errors.Is(err, closeErr) || !w.closed {
				t.Fatalf("error = %v, closed = %v; want primary %v", err, w.closed, want)
			}
		})
	}
	t.Run("close error", func(t *testing.T) {
		w := &benchmarkFailureWriter{failAt: math.MaxInt, closeErr: closeErr}
		if err := writeBenchmarkJSONAndClose(w, &BuildResult{}); !errors.Is(err, closeErr) || !w.closed {
			t.Fatalf("close error = %v, closed = %v", err, w.closed)
		}
	})
	for _, value := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
		t.Run(fmt.Sprint(value), func(t *testing.T) {
			var output bytes.Buffer
			err := writeBenchmarkJSON(&output, &BuildResult{Duration: value})
			var unsupported *json.UnsupportedValueError
			if !errors.As(err, &unsupported) || output.Len() != 0 {
				t.Fatalf("invalid header: error = %v, output bytes = %d", err, output.Len())
			}
		})
	}
	t.Run("invalid executor", func(t *testing.T) {
		w := &benchmarkFailureWriter{closeErr: closeErr}
		err := writeBenchmarkJSONAndClose(w, &BuildResult{Executor: "unsupported"})
		if err == nil || !strings.Contains(err.Error(), "unsupported benchmark executor") || w.calls != 0 || !w.closed {
			t.Fatalf("invalid executor: error = %v, writes = %d, closed = %v", err, w.calls, w.closed)
		}
	})
	t.Run("layout guard", func(t *testing.T) {
		if err := streamBenchmarkJSON(io.Discard, []byte("{}"), benchmarkStreamFixture(1, 1).Content.Entries); err == nil {
			t.Fatal("expected header layout error")
		}
	})
}

func TestBenchmarkJSONStreamFile(t *testing.T) {
	result := benchmarkStreamFixture(3, 4)
	path := filepath.Join(t.TempDir(), "nested", "benchmark.json")
	if err := writeBenchmarkJSONFile(path, result); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var want bytes.Buffer
	if err := legacyBenchmarkJSON(&want, result); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, want.Bytes()) {
		t.Fatal("file differs from legacy")
	}
	// A directory is an invalid file destination on every supported platform.
	if err := writeBenchmarkJSONFile(filepath.Dir(path), result); err == nil {
		t.Fatal("expected create failure")
	}
}

type benchmarkCountingWriter struct {
	bytes int64
	calls int
}

func (w *benchmarkCountingWriter) Write(p []byte) (int, error) {
	w.bytes += int64(len(p))
	w.calls++
	return len(p), nil
}

func TestBenchmarkJSONStreamAllocation(t *testing.T) {
	result := benchmarkStreamFixture(1500, 24)
	measure := func(write func(io.Writer, *BuildResult) error) uint64 {
		t.Helper()
		// Warm reusable encoding/json pools before measuring. GC bounds leftover
		// buffers between implementations; fixture allocation is excluded.
		if err := write(io.Discard, result); err != nil {
			t.Fatal(err)
		}
		runtime.GC()
		var before, after runtime.MemStats
		runtime.ReadMemStats(&before)
		if err := write(io.Discard, result); err != nil {
			t.Fatal(err)
		}
		runtime.ReadMemStats(&after)
		return after.TotalAlloc - before.TotalAlloc
	}
	legacy, streamed := measure(legacyBenchmarkJSON), measure(writeBenchmarkJSON)
	t.Logf("legacy = %d B, streamed = %d B; allocation reduction = %.2f%%",
		legacy, streamed, 100*(1-float64(streamed)/float64(legacy)))
	// Generous relative budget survives pool/GC/toolchain variation while
	// rejecting any implementation that still materializes the full report.
	if streamed >= legacy/8 {
		t.Fatalf("stream allocated %d bytes; want below 1/8 of legacy %d", streamed, legacy)
	}
	var output benchmarkCountingWriter
	if err := writeBenchmarkJSON(&output, result); err != nil {
		t.Fatal(err)
	}
	if output.calls > int(output.bytes/4096)+2 {
		t.Fatalf("unbuffered writes: %d calls for %d bytes", output.calls, output.bytes)
	}
	runtime.KeepAlive(result)
}

func BenchmarkBenchmarkJSON(b *testing.B) {
	result := benchmarkStreamFixture(1500, 24)
	for _, test := range []struct {
		name  string
		write func(io.Writer, *BuildResult) error
	}{
		{"Legacy", legacyBenchmarkJSON},
		{"Stream", writeBenchmarkJSON},
	} {
		b.Run(test.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				if err := test.write(io.Discard, result); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
