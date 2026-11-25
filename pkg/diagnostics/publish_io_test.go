package diagnostics

import (
	"bytes"
	"errors"
	"sync"
	"testing"
	"time"
)

func TestPublishIORecorderCountersAndBoundedLatency(t *testing.T) {
	var recorder PublishIORecorder
	recorder.Consider(3)
	recorder.Skip()
	recorder.OutputWritten()
	recorder.FileUnchanged()
	recorder.RecordWrite(20, 2*time.Millisecond, nil)
	recorder.RecordWrite(30, 8*time.Millisecond, nil)
	recorder.RecordWrite(10, time.Millisecond, errors.New("disk full"))

	got := recorder.Snapshot()
	if got.OutputsConsidered != 3 || got.OutputsSkipped != 1 || got.OutputsWritten != 1 {
		t.Fatalf("unexpected output counts: %+v", got)
	}
	if got.FilesWritten != 2 || got.FilesUnchanged != 1 || got.WriteErrors != 1 || got.BytesWritten != 50 {
		t.Fatalf("unexpected file counts: %+v", got)
	}
	if got.WriteDurationNS != int64(10*time.Millisecond) || got.MaxWriteDurationNS != int64(8*time.Millisecond) {
		t.Fatalf("unexpected write durations: %+v", got)
	}
	if got.P95WriteDurationNS < got.MaxWriteDurationNS || got.P95WriteDurationNS > int64(16*time.Millisecond) {
		t.Fatalf("p95 must be the bounded upper edge of the slowest bucket: %+v", got)
	}
	if got.WriteThroughputBytesPerSec != 5000 {
		t.Fatalf("unexpected summed-write throughput: %v", got.WriteThroughputBytesPerSec)
	}

	var nilRecorder *PublishIORecorder
	nilRecorder.Consider(1)
	nilRecorder.Skip()
	nilRecorder.OutputWritten()
	nilRecorder.FileUnchanged()
	nilRecorder.RecordWrite(1, time.Second, nil)
	if nilRecorder.Snapshot() != (PublishIOCounters{}) {
		t.Fatal("nil recorder should be a no-op")
	}
}

func TestPublishIORecorderConcurrentFiles(t *testing.T) {
	var recorder PublishIORecorder
	const count = 128
	var wg sync.WaitGroup
	wg.Add(count)
	for i := 0; i < count; i++ {
		go func() {
			defer wg.Done()
			recorder.RecordWrite(4, time.Millisecond, nil)
		}()
	}
	wg.Wait()
	got := recorder.Snapshot()
	if got.FilesWritten != count || got.BytesWritten != count*4 ||
		got.WriteDurationNS != int64(count)*int64(time.Millisecond) {
		t.Fatalf("concurrent writes lost counts: %+v", got)
	}
}

func TestPublishIOArtifactRoundTripAndLedgerOwnership(t *testing.T) {
	ledger := NewContentLedger()
	ledger.Discover([]string{"content/example.md"})
	ledger.SetPublishIO("publish_html", PublishIOCounters{
		OutputsConsidered: 3, OutputsSkipped: 2, OutputsWritten: 1,
		FilesWritten: 1, BytesWritten: 42,
	})
	ledger.SetPublishIO("publish_feeds", PublishIOCounters{
		OutputsConsidered: 5, OutputsSkipped: 5,
	})
	ledger.SetPublishIO("unknown", PublishIOCounters{FilesWritten: 99})
	first := ledger.Snapshot()

	ledger.SetPublishIO("publish_html", PublishIOCounters{FilesWritten: 999})
	if first.PublishIO == nil || first.PublishIO.HTML.FilesWritten != 1 ||
		first.PublishIO.Feeds.OutputsSkipped != 5 {
		t.Fatalf("snapshot mutated after new publication: %+v", first.PublishIO)
	}
	first.PublishIO.HTML.FilesWritten = 123
	if got := ledger.Snapshot().PublishIO.HTML.FilesWritten; got != 999 {
		t.Fatalf("snapshot aliases live ledger: %d", got)
	}
	first.PublishIO.HTML.FilesWritten = 1

	info := ArtifactBuildInfo{MarkataVersion: "test", BuiltAt: time.Unix(0, 0).UTC(), Executor: ArtifactExecutorLegacy}
	expected, err := MarshalArtifact(first, info)
	if err != nil {
		t.Fatal(err)
	}
	var streamed bytes.Buffer
	if err := WriteArtifact(&streamed, first, info); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(expected, streamed.Bytes()) {
		t.Fatal("streamed artifact differs from marshaled artifact with publish I/O metrics")
	}
	decoded, err := ParseArtifact(streamed.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	if decoded.PublishIO == nil || decoded.PublishIO.HTML.BytesWritten != 42 ||
		decoded.PublishIO.Feeds.OutputsConsidered != 5 {
		t.Fatalf("publish I/O lost in artifact: %+v", decoded.PublishIO)
	}

	ledger.Reset()
	if ledger.Snapshot().PublishIO != nil {
		t.Fatal("reset must clear old publish metrics")
	}
	ledger.SetPublishIO("publish_feeds", PublishIOCounters{FilesWritten: 1})
	ledger.Discover([]string{"content/next.md"})
	if ledger.Snapshot().PublishIO != nil {
		t.Fatal("new discovery must clear old publish metrics")
	}
}
