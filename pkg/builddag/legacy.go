package builddag

import (
	"context"
	"fmt"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
)

// LegacyTasks returns the runnable plugin hooks for one lifecycle stage as an
// explicit serial task chain. The first task depends on requires; every later
// task depends on the previous plugin's completion artifact.
//
// Legacy plugin tasks are conservative migration boundaries: they are site
// scoped, exclusive, and never marked parallel-safe because plugins may mutate
// Manager, Posts, Config, caches, or plugin-owned state.
func LegacyTasks(m *lifecycle.Manager, stage lifecycle.Stage, requires []ArtifactID) []TaskSpec {
	if m == nil {
		return nil
	}

	plugins := lifecycle.SortPluginsByPriority(m.Plugins(), stage)
	tasks := make([]TaskSpec, 0, len(plugins))
	previous := append([]ArtifactID(nil), requires...)
	taskIndex := 0

	for _, plugin := range plugins {
		if !lifecycle.PluginSupportsStage(plugin, stage) {
			continue
		}

		id := TaskID(fmt.Sprintf("legacy.%s.%03d.%s", stage, taskIndex, plugin.Name()))
		taskIndex++
		completion := ArtifactID{Kind: "legacy-completion", Key: string(id)}
		currentPlugin := plugin
		currentStage := stage

		tasks = append(tasks, TaskSpec{
			ID:           id,
			Group:        string(stage),
			Requires:     append([]ArtifactID(nil), previous...),
			Provides:     []ArtifactID{completion},
			Scope:        ScopeSite,
			Version:      "legacy-v1",
			Exclusive:    true,
			ParallelSafe: false,
			Func: func(ctx context.Context) error {
				if err := ctx.Err(); err != nil {
					return err
				}
				return lifecycle.ExecutePluginHook(m, currentPlugin, currentStage)
			},
		})
		previous = []ArtifactID{completion}
	}

	return tasks
}
