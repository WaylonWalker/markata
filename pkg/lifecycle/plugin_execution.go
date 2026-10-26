package lifecycle

import (
	"fmt"
	"time"

	"github.com/WaylonWalker/markata/pkg/buildstats"
	"github.com/WaylonWalker/markata/pkg/logging"
)

// SortPluginsByPriority returns a stable priority-ordered copy of plugins for
// stage. It exposes the lifecycle's existing ordering contract to schedulers
// without exposing Manager's plugin slice for mutation.
func SortPluginsByPriority(plugins []Plugin, stage Stage) []Plugin {
	return sortPluginsByPriority(plugins, stage)
}

// PluginSupportsStage reports whether plugin implements the hook interface for
// stage. It is useful to graph compilers that need to enumerate only runnable
// hooks while preserving lifecycle order.
func PluginSupportsStage(plugin Plugin, stage Stage) bool {
	if plugin == nil {
		return false
	}
	switch stage {
	case StageConfigure:
		_, ok := plugin.(ConfigurePlugin)
		return ok
	case StageValidate:
		_, ok := plugin.(ValidatePlugin)
		return ok
	case StageGlob:
		_, ok := plugin.(GlobPlugin)
		return ok
	case StageLoad:
		_, ok := plugin.(LoadPlugin)
		return ok
	case StageTransform:
		_, ok := plugin.(TransformPlugin)
		return ok
	case StageRender:
		_, ok := plugin.(RenderPlugin)
		return ok
	case StageCollect:
		_, ok := plugin.(CollectPlugin)
		return ok
	case StageWrite:
		_, ok := plugin.(WritePlugin)
		return ok
	case StageCleanup:
		_, ok := plugin.(CleanupPlugin)
		return ok
	default:
		return false
	}
}

// ExecutePluginHook executes exactly one plugin hook using the lifecycle's
// current-stage, error-severity, warning collection, timing, and active-plugin
// semantics.
//
// It deliberately does not mark a lifecycle stage complete. Schedulers can use
// this function to compose a stage from explicit plugin tasks, then decide when
// the stage boundary itself has completed.
func ExecutePluginHook(m *Manager, plugin Plugin, stage Stage) error {
	if m == nil {
		return fmt.Errorf("lifecycle: manager is nil")
	}
	if plugin == nil {
		return fmt.Errorf("lifecycle: plugin is nil")
	}
	if !IsValidStage(stage) {
		return fmt.Errorf("lifecycle: invalid stage %q", stage)
	}
	if !PluginSupportsStage(plugin, stage) {
		return fmt.Errorf("lifecycle: plugin %q does not support stage %s", plugin.Name(), stage)
	}

	m.mu.Lock()
	m.currentStage = stage
	m.mu.Unlock()

	start := time.Now()
	buildstats.SetActivePlugin(plugin.Name())
	err := executePluginStage(m, plugin, stage)
	buildstats.SetActivePlugin("")

	if err != nil {
		critical := isCriticalStage(stage) || isCriticalError(err) || pluginTreatsStageErrorsAsCritical(plugin, stage)
		hookErr := &HookError{
			Stage:    stage,
			Plugin:   plugin.Name(),
			Err:      err,
			Critical: critical,
		}
		if critical {
			// Match executeHooks: critical hooks stop immediately before the
			// normal successful/non-critical timing bookkeeping.
			return hookErr
		}

		m.mu.Lock()
		m.warnings = append(m.warnings, hookErr)
		m.mu.Unlock()
	}

	elapsed := time.Since(start)
	buildstats.RecordPlugin(string(stage), plugin.Name(), elapsed)
	if elapsed > 50*time.Millisecond {
		logging.Component(plugin.Name()).Phase(string(stage)).Printf("took %v", elapsed)
	}
	return nil
}

func executePluginStage(m *Manager, plugin Plugin, stage Stage) error {
	var err error
	switch stage {
	case StageConfigure:
		hook, ok := plugin.(ConfigurePlugin)
		if !ok {
			return fmt.Errorf("lifecycle: plugin %q does not support stage %s", plugin.Name(), stage)
		}
		err = hook.Configure(m)
	case StageValidate:
		hook, ok := plugin.(ValidatePlugin)
		if !ok {
			return fmt.Errorf("lifecycle: plugin %q does not support stage %s", plugin.Name(), stage)
		}
		err = hook.Validate(m)
	case StageGlob:
		hook, ok := plugin.(GlobPlugin)
		if !ok {
			return fmt.Errorf("lifecycle: plugin %q does not support stage %s", plugin.Name(), stage)
		}
		err = hook.Glob(m)
	case StageLoad:
		hook, ok := plugin.(LoadPlugin)
		if !ok {
			return fmt.Errorf("lifecycle: plugin %q does not support stage %s", plugin.Name(), stage)
		}
		err = hook.Load(m)
	case StageTransform:
		hook, ok := plugin.(TransformPlugin)
		if !ok {
			return fmt.Errorf("lifecycle: plugin %q does not support stage %s", plugin.Name(), stage)
		}
		err = hook.Transform(m)
	case StageRender:
		hook, ok := plugin.(RenderPlugin)
		if !ok {
			return fmt.Errorf("lifecycle: plugin %q does not support stage %s", plugin.Name(), stage)
		}
		err = hook.Render(m)
	case StageCollect:
		hook, ok := plugin.(CollectPlugin)
		if !ok {
			return fmt.Errorf("lifecycle: plugin %q does not support stage %s", plugin.Name(), stage)
		}
		err = hook.Collect(m)
	case StageWrite:
		hook, ok := plugin.(WritePlugin)
		if !ok {
			return fmt.Errorf("lifecycle: plugin %q does not support stage %s", plugin.Name(), stage)
		}
		err = hook.Write(m)
	case StageCleanup:
		hook, ok := plugin.(CleanupPlugin)
		if !ok {
			return fmt.Errorf("lifecycle: plugin %q does not support stage %s", plugin.Name(), stage)
		}
		err = hook.Cleanup(m)
	default:
		return fmt.Errorf("lifecycle: invalid stage %q", stage)
	}
	return err
}
