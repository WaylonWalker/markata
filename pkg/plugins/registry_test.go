package plugins

import (
	"testing"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

func TestDefaultPlugins_HaveUniqueNames(t *testing.T) {
	seen := make(map[string]struct{})

	for _, plugin := range DefaultPlugins() {
		name := plugin.Name()
		if _, exists := seen[name]; exists {
			t.Fatalf("duplicate default plugin registration: %q", name)
		}
		seen[name] = struct{}{}
	}
}

func TestDiagnosticsArtifactPlugin_IsRegisteredByDefault(t *testing.T) {
	defaultManager := lifecycle.NewManager()
	defaultManager.RegisterPlugins(DefaultPlugins()...)
	plugin := defaultPluginByName(t, defaultManager, "diagnostics_artifact")
	if _, ok := plugin.(*DiagnosticsArtifactPlugin); !ok {
		t.Fatalf("default diagnostics artifact plugin = %T, want *DiagnosticsArtifactPlugin", plugin)
	}

	foundMinimal := false
	for _, plugin := range MinimalPlugins() {
		if plugin.Name() == "diagnostics_artifact" {
			foundMinimal = true
			break
		}
	}
	if !foundMinimal {
		t.Fatal("MinimalPlugins() does not include diagnostics_artifact")
	}
}

func TestDefaultPlugins_MultiStageCoverage(t *testing.T) {
	manager := lifecycle.NewManager()
	manager.RegisterPlugins(DefaultPlugins()...)

	stats := defaultPluginByName(t, manager, "stats")
	if _, ok := stats.(*StatsPlugin); !ok {
		t.Fatalf("stats registered as %T, want *StatsPlugin", stats)
	}
	if _, ok := stats.(lifecycle.ConfigurePlugin); !ok {
		t.Error("stats does not implement configure stage")
	}
	if _, ok := stats.(lifecycle.TransformPlugin); !ok {
		t.Error("stats does not implement transform stage")
	}
	if _, ok := stats.(lifecycle.CollectPlugin); !ok {
		t.Error("stats does not implement collect stage")
	}
	assertPluginPriority(t, stats, lifecycle.StageTransform, lifecycle.PriorityEarly)
	assertPluginPriority(t, stats, lifecycle.StageCollect, lifecycle.PriorityLate)
	readingTime := defaultPluginByName(t, manager, "reading_time")
	if _, ok := readingTime.(lifecycle.PriorityPlugin); ok {
		t.Error("reading_time should not add lifecycle ordering requirements")
	}
	jinja := defaultPluginByName(t, manager, "jinja_md")
	assertPluginPriority(t, jinja, lifecycle.StageTransform, lifecycle.PriorityEarly)

	blogroll := defaultPluginByName(t, manager, "blogroll")
	if _, ok := blogroll.(*BlogrollPlugin); !ok {
		t.Fatalf("blogroll registered as %T, want *BlogrollPlugin", blogroll)
	}
	if _, ok := blogroll.(lifecycle.ConfigurePlugin); !ok {
		t.Error("blogroll does not implement configure stage")
	}
	if _, ok := blogroll.(lifecycle.CollectPlugin); !ok {
		t.Error("blogroll does not implement collect stage")
	}
	if _, ok := blogroll.(lifecycle.WritePlugin); !ok {
		t.Error("blogroll does not implement write stage")
	}
	assertPluginPriority(t, blogroll, lifecycle.StageConfigure, lifecycle.PriorityDefault)
	assertPluginPriority(t, blogroll, lifecycle.StageCollect, lifecycle.PriorityLate+10)
	assertPluginPriority(t, blogroll, lifecycle.StageWrite, lifecycle.PriorityLate+20)
}

func TestDefaultPlugins_StatsLifecycleExecution(t *testing.T) {
	defaults := lifecycle.NewManager()
	defaults.RegisterPlugins(DefaultPlugins()...)
	stats, ok := defaultPluginByName(t, defaults, "stats").(*StatsPlugin)
	if !ok {
		t.Fatal("stats default registration is not a *StatsPlugin")
	}

	manager := lifecycle.NewManager()
	manager.Config().Extra["stats"] = map[string]interface{}{"words_per_minute": 1}
	post := models.NewPost("post.md")
	post.Content = "one two"
	manager.SetPosts([]*models.Post{post})
	manager.SetFeeds([]*lifecycle.Feed{{Name: "posts", Posts: []*models.Post{post}}})
	manager.RegisterPlugin(stats)

	if err := manager.RunTo(lifecycle.StageCollect); err != nil {
		t.Fatalf("RunTo(collect) error = %v", err)
	}
	if got := post.Get("word_count"); got != 2 {
		t.Errorf("word_count = %v, want 2", got)
	}
	if got := post.Get("reading_time"); got != 2 {
		t.Errorf("reading_time = %v, want 2 after Configure", got)
	}
	if _, ok := manager.Cache().Get("feed_stats.posts"); !ok {
		t.Error("Collect did not store feed stats")
	}
	if _, ok := manager.Cache().Get("site_stats"); !ok {
		t.Error("Collect did not store site stats")
	}
}

func defaultPluginByName(t *testing.T, manager *lifecycle.Manager, name string) lifecycle.Plugin {
	t.Helper()
	var found lifecycle.Plugin
	for _, plugin := range manager.Plugins() {
		if plugin.Name() != name {
			continue
		}
		if found != nil {
			t.Fatalf("registered %q more than once", name)
		}
		found = plugin
	}
	if found == nil {
		t.Fatalf("default plugin %q not registered", name)
	}
	return found
}

func assertPluginPriority(t *testing.T, plugin lifecycle.Plugin, stage lifecycle.Stage, want int) {
	t.Helper()
	priorityPlugin, ok := plugin.(lifecycle.PriorityPlugin)
	if !ok {
		t.Fatalf("%q does not implement lifecycle.PriorityPlugin", plugin.Name())
	}
	if got := priorityPlugin.Priority(stage); got != want {
		t.Errorf("%q %s priority = %d, want %d", plugin.Name(), stage, got, want)
	}
}
