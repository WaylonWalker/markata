package builderadmin

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/WaylonWalker/markata/pkg/buildstats"
	"github.com/WaylonWalker/markata/pkg/diagnostics"
)

const benchmarkArtifactDirName = "benchmarks"

// BuildBenchmarkArtifact is the stable subset of `markata-go build
// --benchmark-json` that Builder Admin retains for diagnostics. Unknown fields
// are intentionally ignored so CLI presentation can evolve independently.
type BuildBenchmarkArtifact struct {
	PostsProcessed int                               `json:"posts_processed"`
	FeedsGenerated int                               `json:"feeds_generated"`
	Duration       float64                           `json:"duration_seconds"`
	Warnings       []string                          `json:"warnings,omitempty"`
	Benchmark      buildstats.Summary                `json:"benchmark"`
	Content        diagnostics.ContentLedgerSnapshot `json:"content"`
}

func benchmarkArtifactPath(historyDir, buildID string) (string, error) {
	buildID = strings.TrimSpace(buildID)
	if buildID == "" || filepath.Base(buildID) != buildID || strings.ContainsAny(buildID, `/\\`) {
		return "", fmt.Errorf("invalid build id %q", buildID)
	}
	return filepath.Join(historyDir, benchmarkArtifactDirName, buildID+".json"), nil
}

func readBuildBenchmarkArtifact(path string) (*BuildBenchmarkArtifact, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var artifact BuildBenchmarkArtifact
	if err := json.Unmarshal(data, &artifact); err != nil {
		return nil, fmt.Errorf("decode benchmark artifact: %w", err)
	}
	return &artifact, nil
}

func ensureBenchmarkArtifactDir(historyDir string) (string, error) {
	dir := filepath.Join(historyDir, benchmarkArtifactDirName)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("create benchmark artifact dir: %w", err)
	}
	return dir, nil
}
