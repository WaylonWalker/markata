package builderadmin

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"io"
	"os"

	"github.com/WaylonWalker/markata/pkg/diagnostics"
)

const diagnosticsArtifactHeaderLimit = 64 << 10

// readDiagnosticsArtifactBounded reads either the plain or compressed release
// artifact, with the same uncompressed-size limit for both formats.
func readDiagnosticsArtifactBounded(path string, maxBytes int64) ([]byte, bool) {
	if info, err := os.Stat(path); err == nil {
		if info.Size() > maxBytes {
			return nil, false
		}
		data, err := os.ReadFile(path)
		return data, err == nil && int64(len(data)) <= maxBytes
	} else if !os.IsNotExist(err) {
		return nil, false
	}

	return readGzipArtifactBounded(path+".gz", maxBytes)
}

func readGzipArtifactBounded(path string, maxBytes int64) ([]byte, bool) {
	if maxBytes < 0 {
		return nil, false
	}
	if size, ok := gzipArtifactUncompressedSize(path); ok && uint64(size) > uint64(maxBytes) {
		return nil, false
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, false
	}
	defer file.Close()
	reader, err := gzip.NewReader(file)
	if err != nil {
		return nil, false
	}
	data, readErr := io.ReadAll(io.LimitReader(reader, maxBytes+1))
	closeErr := reader.Close()
	if readErr != nil || closeErr != nil || int64(len(data)) > maxBytes {
		return nil, false
	}
	return data, true
}

// gzipArtifactUncompressedSize reads the gzip ISIZE trailer. The trailer wraps
// at 4 GiB; callers still use a bounded decompression fallback when it appears
// to fit under their limit.
func gzipArtifactUncompressedSize(path string) (uint32, bool) {
	file, err := os.Open(path)
	if err != nil {
		return 0, false
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil || info.Size() < 4 {
		return 0, false
	}
	if _, err := file.Seek(-4, io.SeekEnd); err != nil {
		return 0, false
	}
	var trailer [4]byte
	if _, err := io.ReadFull(file, trailer[:]); err != nil {
		return 0, false
	}
	return binary.LittleEndian.Uint32(trailer[:]), true
}

// readDiagnosticsArtifactIdentity reads provenance from a plain artifact or
// from the bounded JSON header of a gzip artifact. The content entries are the
// final JSON field, so Builder Admin does not need to inflate them to discover
// which revisions produced the release.
func readDiagnosticsArtifactIdentity(path string) (diagnostics.Artifact, bool) {
	if info, err := os.Stat(path); err == nil {
		if info.Size() > 2<<20 {
			return diagnostics.Artifact{}, false
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return diagnostics.Artifact{}, false
		}
		artifact, err := diagnostics.ParseArtifact(data)
		return artifact, err == nil
	} else if !os.IsNotExist(err) {
		return diagnostics.Artifact{}, false
	}
	return readGzipArtifactIdentity(path + ".gz")
}

func readGzipArtifactIdentity(path string) (diagnostics.Artifact, bool) {
	file, err := os.Open(path)
	if err != nil {
		return diagnostics.Artifact{}, false
	}
	defer file.Close()
	reader, err := gzip.NewReader(file)
	if err != nil {
		return diagnostics.Artifact{}, false
	}
	data, readErr := io.ReadAll(io.LimitReader(reader, diagnosticsArtifactHeaderLimit))
	_ = reader.Close()
	if readErr != nil || len(data) == diagnosticsArtifactHeaderLimit {
		return diagnostics.Artifact{}, false
	}
	const entriesPrefix = "\n  \"entries\": ["
	index := bytes.Index(data, []byte(entriesPrefix))
	if index < 0 {
		return diagnostics.Artifact{}, false
	}
	header := make([]byte, 0, index+len("\n  \"entries\": []\n}"))
	header = append(header, data[:index]...)
	header = append(header, []byte("\n  \"entries\": []\n}")...)
	artifact, err := diagnostics.ParseArtifact(header)
	if err != nil {
		return diagnostics.Artifact{}, false
	}
	if artifact.Entries == nil {
		return diagnostics.Artifact{}, false
	}
	return artifact, true
}
