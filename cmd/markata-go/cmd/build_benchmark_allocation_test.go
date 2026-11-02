//go:build !race

package cmd

import (
	"io"
	"runtime"
	"slices"
	"testing"
)

// Measure performance without the race detector: instrumented allocations are
// not representative of production and made the former single-sample 75% gate
// flaky under CI's -race -covermode=atomic run. JSON equivalence and buffered
// write correctness still run with -race in build_benchmark_stream_test.go.
func TestBenchmarkJSONStreamAllocation(t *testing.T) {
	result := benchmarkStreamFixture(1500, 24)
	measure := func(write func(io.Writer, *BuildResult) error) uint64 {
		t.Helper()
		// Exclude fixture construction, warm the encoder once, and take the
		// median of multiple samples to reduce GC and runner noise.
		if err := write(io.Discard, result); err != nil {
			t.Fatal(err)
		}
		const sampleCount = 5
		var samples [sampleCount]uint64
		for i := range samples {
			runtime.GC()
			var before, after runtime.MemStats
			runtime.ReadMemStats(&before)
			if err := write(io.Discard, result); err != nil {
				t.Fatal(err)
			}
			runtime.ReadMemStats(&after)
			samples[i] = after.TotalAlloc - before.TotalAlloc
		}
		slices.Sort(samples[:])
		return samples[len(samples)/2]
	}
	legacy, streamed := measure(legacyBenchmarkJSON), measure(writeBenchmarkJSON)
	t.Logf("median of 5 samples: legacy = %d B, streamed = %d B; allocation reduction = %.2f%%",
		legacy, streamed, 100*(1-float64(streamed)/float64(legacy)))

	// Keep both the proportional improvement and an absolute allocation
	// ceiling for this fixed 1,500-entry / 24-feed fixture.
	if streamed >= legacy/4 {
		t.Fatalf("stream allocated %d bytes; want below 1/4 of legacy %d", streamed, legacy)
	}
	const maxStreamAllocationBytes = 48 << 20
	if streamed >= maxStreamAllocationBytes {
		t.Fatalf("stream allocated %d bytes; want below %d bytes", streamed, maxStreamAllocationBytes)
	}
	runtime.KeepAlive(result)
}
