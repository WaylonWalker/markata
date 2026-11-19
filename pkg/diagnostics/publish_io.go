package diagnostics

import (
	"math"
	"math/bits"
	"sync/atomic"
	"time"
)

// PublishIOStats summarizes the two core publishers without retaining filenames,
// output contents, or per-file latency samples. The fields are optional in v1
// artifacts so older artifacts and builds without these plugins remain valid.
type PublishIOStats struct {
	HTML  *PublishIOCounters `json:"publish_html,omitempty"`
	Feeds *PublishIOCounters `json:"publish_feeds,omitempty"`
}

// PublishIOCounters separates logical outputs (posts or feeds) from actual file
// writes: a single output can create several files, including redirects.
// WriteDurationNS is the sum of successful file-write syscall durations, NOT
// elapsed wall time; parallel writes can therefore exceed stage wall time.
// Throughput is bytes / summed write seconds, not a filesystem bandwidth test.
// P95WriteDurationNS is the upper edge of a power-of-two latency bucket.
type PublishIOCounters struct {
	OutputsConsidered          int64   `json:"outputs_considered"`
	OutputsSkipped             int64   `json:"outputs_skipped"`
	OutputsWritten             int64   `json:"outputs_written"`
	FilesWritten               int64   `json:"files_written"`
	FilesUnchanged             int64   `json:"files_unchanged"`
	WriteErrors                int64   `json:"write_errors"`
	BytesWritten               int64   `json:"bytes_written"`
	WriteDurationNS            int64   `json:"write_duration_ns"`
	P95WriteDurationNS         int64   `json:"p95_write_duration_ns"`
	MaxWriteDurationNS         int64   `json:"max_write_duration_ns"`
	WriteThroughputBytesPerSec float64 `json:"write_throughput_bytes_per_sec"`
}

// PublishIORecorder is build-local and safe for parallel publishing. All data
// is fixed-size; file writes need no allocations or path retention.
type PublishIORecorder struct {
	considered   atomic.Int64
	skipped      atomic.Int64
	written      atomic.Int64
	files        atomic.Int64
	unchanged    atomic.Int64
	errors       atomic.Int64
	bytes        atomic.Int64
	durationNS   atomic.Int64
	maxNS        atomic.Int64
	latencyCount [64]atomic.Int64
}

func (r *PublishIORecorder) Consider(count int) {
	if r != nil && count > 0 {
		r.considered.Add(int64(count))
	}
}

func (r *PublishIORecorder) Skip() {
	if r != nil {
		r.skipped.Add(1)
	}
}

func (r *PublishIORecorder) OutputWritten() {
	if r != nil {
		r.written.Add(1)
	}
}

func (r *PublishIORecorder) FileUnchanged() {
	if r != nil {
		r.unchanged.Add(1)
	}
}

// RecordWrite records only completed os.WriteFile calls. Failures count
// separately and do not contribute misleading successful bytes or latency.
func (r *PublishIORecorder) RecordWrite(size int, elapsed time.Duration, err error) {
	if r == nil {
		return
	}
	if err != nil {
		r.errors.Add(1)
		return
	}
	ns := elapsed.Nanoseconds()
	if ns < 1 {
		ns = 1
	}
	r.files.Add(1)
	r.bytes.Add(int64(size))
	r.durationNS.Add(ns)
	for old := r.maxNS.Load(); ns > old; old = r.maxNS.Load() {
		if r.maxNS.CompareAndSwap(old, ns) {
			break
		}
	}
	r.latencyCount[bits.Len64(uint64(ns))-1].Add(1)
}

// Snapshot returns standalone scalar metrics; it never exposes atomic fields.
func (r *PublishIORecorder) Snapshot() PublishIOCounters {
	if r == nil {
		return PublishIOCounters{}
	}
	stats := PublishIOCounters{
		OutputsConsidered: r.considered.Load(),
		OutputsSkipped:    r.skipped.Load(),
		OutputsWritten:    r.written.Load(),
		FilesWritten:     r.files.Load(),
		FilesUnchanged:   r.unchanged.Load(),
		WriteErrors:      r.errors.Load(),
		BytesWritten:     r.bytes.Load(),
		WriteDurationNS:  r.durationNS.Load(),
		MaxWriteDurationNS: r.maxNS.Load(),
	}
	if stats.WriteDurationNS > 0 {
		stats.WriteThroughputBytesPerSec = float64(stats.BytesWritten) * float64(time.Second) / float64(stats.WriteDurationNS)
	}
	// Ceiling of 95% of successful file writes. A fixed log2 histogram
	// prevents allocations and bounds diagnostic overhead under load.
	target := (stats.FilesWritten*95 + 99) / 100
	if target > 0 {
		var seen int64
		for bucket := range r.latencyCount {
			seen += r.latencyCount[bucket].Load()
			if seen >= target {
				if bucket >= 62 {
					stats.P95WriteDurationNS = math.MaxInt64
				} else {
					stats.P95WriteDurationNS = (int64(1) << (bucket + 1)) - 1
				}
				break
			}
		}
	}
	return stats
}

func clonePublishIOStats(stats *PublishIOStats) *PublishIOStats {
	if stats == nil {
		return nil
	}
	cloned := *stats
	if stats.HTML != nil {
		value := *stats.HTML
		cloned.HTML = &value
	}
	if stats.Feeds != nil {
		value := *stats.Feeds
		cloned.Feeds = &value
	}
	return &cloned
}

// SetPublishIO copies completed plugin counters into the content ledger.
// Unknown names are ignored to keep the artifact schema bounded.
func (l *ContentLedger) SetPublishIO(plugin string, counters PublishIOCounters) {
	if l == nil {
		return
	}
	if plugin != "publish_html" && plugin != "publish_feeds" {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.publishIO == nil {
		l.publishIO = &PublishIOStats{}
	}
	value := counters
	if plugin == "publish_html" {
		l.publishIO.HTML = &value
	} else {
		l.publishIO.Feeds = &value
	}
}
