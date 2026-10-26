package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
)

// MarshalJSON adds executor identity to the existing benchmark payload without
// changing the human build result or the benchmark writer call sites.
func (output benchmarkJSONOutput) MarshalJSON() ([]byte, error) {
	executor := output.Executor
	if executor == "" {
		executor = lifecycle.BuildExecutorLegacy
	}
	if !executor.Valid() {
		return nil, fmt.Errorf("unsupported benchmark executor %q", executor)
	}
	type benchmarkJSONAlias benchmarkJSONOutput
	return json.Marshal(struct {
		Executor lifecycle.BuildExecutor `json:"executor"`
		benchmarkJSONAlias
	}{
		Executor:           executor,
		benchmarkJSONAlias: benchmarkJSONAlias(output),
	})
}
