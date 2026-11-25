package builderadmin

import (
	"bytes"
	"compress/gzip"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/WaylonWalker/markata/pkg/diagnostics"
)

func writeCompressedArtifactFixture(t *testing.T, path string, data []byte) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	file, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	writer := gzip.NewWriter(file)
	if _, err := writer.Write(data); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestReadDiagnosticsArtifactBounded_Gzip(t *testing.T) {
	path := filepath.Join(t.TempDir(), diagnostics.DefaultArtifactPath)
	want, err := diagnostics.MarshalArtifact(diagnostics.ContentLedgerSnapshot{
		Summary: diagnostics.ContentSummary{Discovered: 1},
		Entries: []diagnostics.ContentDisposition{{Path: "posts/example.md", Emitted: true}},
	}, diagnostics.ArtifactBuildInfo{
		MarkataVersion: "test",
		MarkataCommit:  "0123456789abcdef0123456789abcdef01234567",
		SourceCommit:   "abcdef0123456789abcdef0123456789abcdef01",
		BuiltAt:        time.Unix(1, 0),
	})
	if err != nil {
		t.Fatal(err)
	}
	writeCompressedArtifactFixture(t, path+".gz", want)

	got, ok := readDiagnosticsArtifactBounded(path, int64(len(want)))
	if !ok || !bytes.Equal(got, want) {
		t.Fatalf("bounded gzip read changed bytes or failed: ok=%t len=%d want=%d", ok, len(got), len(want))
	}
	if _, ok := readDiagnosticsArtifactBounded(path, int64(len(want)-1)); ok {
		t.Fatal("gzip artifact larger than the uncompressed limit was accepted")
	}
}

func TestReadDiagnosticsArtifactIdentity_GzipSkipsEntries(t *testing.T) {
	path := filepath.Join(t.TempDir(), diagnostics.DefaultArtifactPath)
	want, err := diagnostics.MarshalArtifact(diagnostics.ContentLedgerSnapshot{
		Summary: diagnostics.ContentSummary{Discovered: 2},
		Entries: []diagnostics.ContentDisposition{
			{Path: "posts/one.md", Emitted: true},
			{Path: "posts/two.md", Excluded: true},
		},
	}, diagnostics.ArtifactBuildInfo{
		MarkataVersion: "test",
		MarkataCommit:  "0123456789abcdef0123456789abcdef01234567",
		SourceCommit:   "abcdef0123456789abcdef0123456789abcdef01",
		BuiltAt:        time.Unix(1, 0),
	})
	if err != nil {
		t.Fatal(err)
	}
	writeCompressedArtifactFixture(t, path+".gz", want)

	artifact, ok := readDiagnosticsArtifactIdentity(path)
	if !ok {
		t.Fatal("failed to read compressed artifact identity")
	}
	if artifact.Generator.Commit != "0123456789abcdef0123456789abcdef01234567" {
		t.Fatalf("generator commit = %q", artifact.Generator.Commit)
	}
	if artifact.Source == nil || artifact.Source.Commit != "abcdef0123456789abcdef0123456789abcdef01" {
		t.Fatalf("source = %#v", artifact.Source)
	}
	if len(artifact.Entries) != 0 {
		t.Fatalf("identity reader inflated entries: %#v", artifact.Entries)
	}
}

func TestReadDiagnosticsArtifactIdentity_GzipWithNoEntries(t *testing.T) {
	path := filepath.Join(t.TempDir(), diagnostics.DefaultArtifactPath)
	want, err := diagnostics.MarshalArtifact(diagnostics.ContentLedgerSnapshot{
		Summary: diagnostics.ContentSummary{Discovered: 0},
		Entries: []diagnostics.ContentDisposition{},
	}, diagnostics.ArtifactBuildInfo{
		MarkataVersion: "test",
		MarkataCommit:  "0123456789abcdef0123456789abcdef01234567",
		SourceCommit:   "abcdef0123456789abcdef0123456789abcdef01",
		BuiltAt:        time.Unix(1, 0),
	})
	if err != nil {
		t.Fatal(err)
	}
	writeCompressedArtifactFixture(t, path+".gz", want)

	artifact, ok := readDiagnosticsArtifactIdentity(path)
	if !ok {
		t.Fatal("failed to read identity from empty compressed artifact")
	}
	if artifact.Generator.Commit != "0123456789abcdef0123456789abcdef01234567" || artifact.Source == nil {
		t.Fatalf("identity = generator %#v source %#v", artifact.Generator, artifact.Source)
	}
	if artifact.Entries == nil || len(artifact.Entries) != 0 {
		t.Fatalf("entries = %#v, want empty non-nil array", artifact.Entries)
	}
}

func TestReadDiagnosticsArtifactBounded_GzipStopsAtLimit(t *testing.T) {
	path := filepath.Join(t.TempDir(), diagnostics.DefaultArtifactPath)
	writeCompressedArtifactFixture(t, path+".gz", bytes.Repeat([]byte("x"), 4096))
	if data, ok := readDiagnosticsArtifactBounded(path, 1024); ok || data != nil {
		t.Fatalf("oversized artifact accepted: ok=%t len=%d", ok, len(data))
	}
}
