package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/WaylonWalker/markata/pkg/builddag"
	"github.com/WaylonWalker/markata/pkg/buildstats"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/templates"
)

type dagPluginExpander func(*lifecycle.Manager, lifecycle.Stage, builddag.TaskSpec) []builddag.TaskSpec

func defaultDAGPluginExpander(_ *lifecycle.Manager, _ lifecycle.Stage, task builddag.TaskSpec) []builddag.TaskSpec {
	return []builddag.TaskSpec{task}
}

func isDAGPluginStage(stage lifecycle.Stage) bool {
	switch stage {
	case lifecycle.StageLoad, lifecycle.StageTransform, lifecycle.StageRender,
		lifecycle.StageCollect, lifecycle.StageWrite, lifecycle.StageCleanup:
		return true
	case lifecycle.StageConfigure, lifecycle.StageValidate, lifecycle.StageGlob:
		return false
	default:
		return false
	}
}

func executeDAGPluginStageSegments(
	executor *builddag.Executor,
	m *lifecycle.Manager,
	stage lifecycle.Stage,
	requires []builddag.ArtifactID,
	provided builddag.ArtifactID,
	observe func(lifecycle.Stage, bool, error),
	expand dagPluginExpander,
) ([]builddag.SegmentDigest, int, error) {
	if expand == nil {
		expand = defaultDAGPluginExpander
	}

	startedArtifact := builddag.ArtifactID{Kind: "lifecycle-stage-start", Key: string(stage)}
	var stageStart time.Time
	var skipStage bool

	startBuilder := builddag.NewBuilder()
	addDAGExternalInputs(startBuilder, requires)
	startBuilder.AddTask(builddag.TaskSpec{
		ID:        builddag.TaskID("lifecycle." + string(stage) + ".start"),
		Group:     string(stage),
		Requires:  append([]builddag.ArtifactID(nil), requires...),
		Provides:  []builddag.ArtifactID{startedArtifact},
		Scope:     builddag.ScopeSite,
		Version:   "legacy-plugin-stage-v2",
		Exclusive: true,
		Func: func(ctx context.Context) error {
			if err := ctx.Err(); err != nil {
				return err
			}
			stageStart = time.Now()
			if observe != nil {
				observe(stage, true, nil)
			}
			buildstats.SetActiveStage(string(stage))
			verbosef("  [%s] running as serial plugin DAG...", stage)
			templates.ClearAllCaches()
			skipStage = m.HasRun(stage)
			return nil
		},
	})

	segments := make([]builddag.SegmentDigest, 0)
	totalTasks := 0
	segment, count, err := compileExecuteDAGSegment(executor, startBuilder, string(stage)+".start")
	if err != nil {
		return nil, 0, err
	}
	segments = append(segments, segment)
	totalTasks += count

	lastRequires := []builddag.ArtifactID{startedArtifact}
	// Snapshot the ordered stage plugin tasks once, matching legacy hook-list
	// semantics. Each task's graph is expanded/compiled only when execution
	// reaches that plugin, after all preceding plugin work has completed.
	pluginTasks := builddag.LegacyTasks(m, stage, lastRequires)
	if !skipStage {
		for index := range pluginTasks {
			legacyTask := pluginTasks[index]
			expanded := expand(m, stage, legacyTask)
			if err := validateDAGPluginExpansionBoundary(legacyTask, expanded); err != nil {
				finishDAGPluginStageError(stage, err, observe)
				return nil, 0, err
			}

			pluginBuilder := builddag.NewBuilder()
			addDAGExternalInputs(pluginBuilder, legacyTask.Requires)
			for taskIndex := range expanded {
				pluginBuilder.AddTask(expanded[taskIndex])
			}

			segmentName := fmt.Sprintf("%s.plugin.%03d.%s", stage, index, legacyTask.ID)
			pluginSegment, pluginCount, executeErr := compileExecuteDAGSegment(executor, pluginBuilder, segmentName)
			if executeErr != nil {
				finishDAGPluginStageError(stage, executeErr, observe)
				return nil, 0, executeErr
			}
			segments = append(segments, pluginSegment)
			totalTasks += pluginCount
			lastRequires = append([]builddag.ArtifactID(nil), legacyTask.Provides...)
		}
	}

	completeBuilder := builddag.NewBuilder()
	addDAGExternalInputs(completeBuilder, lastRequires)
	completeBuilder.AddTask(builddag.TaskSpec{
		ID:        builddag.TaskID("lifecycle." + string(stage) + ".complete"),
		Group:     string(stage),
		Requires:  append([]builddag.ArtifactID(nil), lastRequires...),
		Provides:  []builddag.ArtifactID{provided},
		Scope:     builddag.ScopeSite,
		Version:   "legacy-plugin-stage-v2",
		Exclusive: true,
		Func: func(ctx context.Context) error {
			if err := ctx.Err(); err != nil {
				return err
			}
			if !skipStage {
				if err := m.MarkStageComplete(stage); err != nil {
					return err
				}
			}

			buildstats.SetActiveStage("")
			if observe != nil {
				observe(stage, false, nil)
			}
			stageElapsed := time.Since(stageStart)
			buildstats.RecordStage(string(stage), stageElapsed)
			if verbose {
				verbosef("  [%s] done in %s", stage, stageElapsed.Truncate(100*time.Microsecond))
			}
			return nil
		},
	})

	segment, count, err = compileExecuteDAGSegment(executor, completeBuilder, string(stage)+".complete")
	if err != nil {
		finishDAGPluginStageError(stage, err, observe)
		return nil, 0, err
	}
	segments = append(segments, segment)
	totalTasks += count
	return segments, totalTasks, nil
}

func compileExecuteDAGSegment(
	executor *builddag.Executor,
	builder *builddag.Builder,
	name string,
) (builddag.SegmentDigest, int, error) {
	graph, err := builder.Compile()
	if err != nil {
		return builddag.SegmentDigest{}, 0, fmt.Errorf("compile DAG segment %q: %w", name, err)
	}
	digest, err := graph.Digest()
	if err != nil {
		return builddag.SegmentDigest{}, 0, fmt.Errorf("digest DAG segment %q: %w", name, err)
	}
	execution, err := executor.Execute(context.Background(), graph)
	if err != nil {
		return builddag.SegmentDigest{}, 0, fmt.Errorf("execute DAG segment %q: %w", name, err)
	}
	return builddag.SegmentDigest{Name: name, Digest: digest, TaskCount: execution.TaskCount}, execution.TaskCount, nil
}

func addDAGExternalInputs(builder *builddag.Builder, artifacts []builddag.ArtifactID) {
	for _, artifact := range artifacts {
		builder.AddExternal(artifact)
	}
}

func validateDAGPluginExpansionBoundary(legacyTask builddag.TaskSpec, expanded []builddag.TaskSpec) error {
	if len(expanded) == 0 {
		return fmt.Errorf("DAG plugin expansion for %q produced no tasks", legacyTask.ID)
	}

	required := make(map[builddag.ArtifactID]bool, len(legacyTask.Requires))
	provided := make(map[builddag.ArtifactID]bool, len(legacyTask.Provides))
	for _, artifact := range legacyTask.Requires {
		required[artifact] = false
	}
	for _, artifact := range legacyTask.Provides {
		provided[artifact] = false
	}
	for index := range expanded {
		for _, artifact := range expanded[index].Requires {
			if _, ok := required[artifact]; ok {
				required[artifact] = true
			}
		}
		for _, artifact := range expanded[index].Provides {
			if _, ok := provided[artifact]; ok {
				provided[artifact] = true
			}
		}
	}
	for artifact, seen := range required {
		if !seen {
			return fmt.Errorf("DAG plugin expansion for %q dropped required boundary artifact %s", legacyTask.ID, artifact.String())
		}
	}
	for artifact, seen := range provided {
		if !seen {
			return fmt.Errorf("DAG plugin expansion for %q dropped provided boundary artifact %s", legacyTask.ID, artifact.String())
		}
	}
	return nil
}
