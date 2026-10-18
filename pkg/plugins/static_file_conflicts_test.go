package plugins

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

func TestStaticFileConflictsPlugin_Name(t *testing.T) {
	plugin := NewStaticFileConflictsPlugin()
	if plugin.Name() != "static_file_conflicts" {
		t.Errorf("expected name 'static_file_conflicts', got %q", plugin.Name())
	}
}

func TestStaticFileConflictsPlugin_Priority(t *testing.T) {
	plugin := NewStaticFileConflictsPlugin()

	// Should return PriorityLate for Collect stage
	if got := plugin.Priority(lifecycle.StageCollect); got != lifecycle.PriorityLate {
		t.Errorf("expected PriorityLate for Collect, got %d", got)
	}

	// Should return PriorityDefault for other stages
	if got := plugin.Priority(lifecycle.StageRender); got != lifecycle.PriorityDefault {
		t.Errorf("expected PriorityDefault for Render, got %d", got)
	}
}

func TestStaticFileConflictsPlugin_NoStaticDir(t *testing.T) {
	// Test that plugin handles missing static directory gracefully
	plugin := NewStaticFileConflictsPlugin()
	plugin.SetStaticDir("nonexistent-static-dir")

	m := lifecycle.NewManager()
	m.AddPost(&models.Post{
		Path: "robots.md",
		Slug: "robots",
	})

	err := plugin.Collect(m)
	if err != nil {
		t.Errorf("expected no error with missing static dir, got %v", err)
	}

	if len(plugin.Conflicts()) != 0 {
		t.Errorf("expected no conflicts with missing static dir, got %d", len(plugin.Conflicts()))
	}
}

func TestStaticFileConflictsPlugin_Disabled(t *testing.T) {
	// Test that disabled plugin does nothing
	plugin := NewStaticFileConflictsPlugin()
	plugin.SetEnabled(false)

	m := lifecycle.NewManager()
	m.AddPost(&models.Post{
		Path: "robots.md",
		Slug: "robots",
	})

	err := plugin.Collect(m)
	if err != nil {
		t.Errorf("expected no error when disabled, got %v", err)
	}

	if len(plugin.Conflicts()) != 0 {
		t.Errorf("expected no conflicts when disabled, got %d", len(plugin.Conflicts()))
	}
}

func TestStaticFileConflictsPlugin_DetectsRobotsTxtConflict(t *testing.T) {
	// Create temp directory structure
	tmpDir := t.TempDir()
	staticDir := filepath.Join(tmpDir, "static")
	if err := os.MkdirAll(staticDir, 0o755); err != nil {
		t.Fatalf("failed to create static dir: %v", err)
	}

	// Create static/robots.txt
	robotsPath := filepath.Join(staticDir, "robots.txt")
	if err := os.WriteFile(robotsPath, []byte("User-agent: *\nDisallow: /"), 0o600); err != nil {
		t.Fatalf("failed to create robots.txt: %v", err)
	}

	// Create plugin and set static dir
	plugin := NewStaticFileConflictsPlugin()
	plugin.SetStaticDir(staticDir)

	// Create manager with a robots.md post
	m := lifecycle.NewManager()
	m.AddPost(&models.Post{
		Path:      "pages/robots.md",
		Slug:      "robots",
		Published: true,
	})

	err := plugin.Collect(m)

	// Should return a warning (error)
	if err == nil {
		t.Error("expected warning about conflict, got nil")
	}

	// Check that conflict was detected
	conflicts := plugin.Conflicts()
	if len(conflicts) == 0 {
		t.Error("expected at least one conflict, got none")
	}

	// Verify the conflict details
	found := false
	for _, c := range conflicts {
		if c.GeneratedSource == "pages/robots.md" && c.OutputPath == "/robots.txt" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected conflict for robots.md → /robots.txt, conflicts: %+v", conflicts)
	}
}

func TestStaticFileConflictsPlugin_DetectsSitemapConflict(t *testing.T) {
	tmpDir := t.TempDir()
	staticDir := filepath.Join(tmpDir, "static")
	if err := os.MkdirAll(staticDir, 0o755); err != nil {
		t.Fatalf("failed to create static dir: %v", err)
	}

	// Create static/sitemap.xml
	sitemapPath := filepath.Join(staticDir, "sitemap.xml")
	if err := os.WriteFile(sitemapPath, []byte("<?xml version=\"1.0\"?>"), 0o600); err != nil {
		t.Fatalf("failed to create sitemap.xml: %v", err)
	}

	plugin := NewStaticFileConflictsPlugin()
	plugin.SetStaticDir(staticDir)

	m := lifecycle.NewManager()
	m.AddPost(&models.Post{
		Path:      "sitemap.md",
		Slug:      "sitemap",
		Published: true,
	})

	err := plugin.Collect(m)

	if err == nil {
		t.Error("expected warning about conflict, got nil")
	}

	conflicts := plugin.Conflicts()
	found := false
	for _, c := range conflicts {
		if c.OutputPath == "/sitemap.xml" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected conflict for sitemap.xml, conflicts: %+v", conflicts)
	}
}

func TestStaticFileConflictsPlugin_SkipsDraftPosts(t *testing.T) {
	tmpDir := t.TempDir()
	staticDir := filepath.Join(tmpDir, "static")
	if err := os.MkdirAll(staticDir, 0o755); err != nil {
		t.Fatalf("failed to create static dir: %v", err)
	}

	// Create static/robots.txt
	robotsPath := filepath.Join(staticDir, "robots.txt")
	if err := os.WriteFile(robotsPath, []byte("User-agent: *"), 0o600); err != nil {
		t.Fatalf("failed to create robots.txt: %v", err)
	}

	plugin := NewStaticFileConflictsPlugin()
	plugin.SetStaticDir(staticDir)

	m := lifecycle.NewManager()
	m.AddPost(&models.Post{
		Path:  "pages/robots.md",
		Slug:  "robots",
		Draft: true, // Draft post should be skipped
	})

	err := plugin.Collect(m)

	if err != nil {
		t.Errorf("expected no error for draft post, got %v", err)
	}

	if len(plugin.Conflicts()) != 0 {
		t.Errorf("expected no conflicts for draft post, got %d", len(plugin.Conflicts()))
	}
}

func TestStaticFileConflictsPlugin_SkipsSkippedPosts(t *testing.T) {
	tmpDir := t.TempDir()
	staticDir := filepath.Join(tmpDir, "static")
	if err := os.MkdirAll(staticDir, 0o755); err != nil {
		t.Fatalf("failed to create static dir: %v", err)
	}

	robotsPath := filepath.Join(staticDir, "robots.txt")
	if err := os.WriteFile(robotsPath, []byte("User-agent: *"), 0o600); err != nil {
		t.Fatalf("failed to create robots.txt: %v", err)
	}

	plugin := NewStaticFileConflictsPlugin()
	plugin.SetStaticDir(staticDir)

	m := lifecycle.NewManager()
	m.AddPost(&models.Post{
		Path: "pages/robots.md",
		Slug: "robots",
		Skip: true, // Skipped post should be skipped
	})

	err := plugin.Collect(m)

	if err != nil {
		t.Errorf("expected no error for skipped post, got %v", err)
	}

	if len(plugin.Conflicts()) != 0 {
		t.Errorf("expected no conflicts for skipped post, got %d", len(plugin.Conflicts()))
	}
}

func TestStaticFileConflictsPlugin_NoConflictForNonRootPosts(t *testing.T) {
	tmpDir := t.TempDir()
	staticDir := filepath.Join(tmpDir, "static")
	if err := os.MkdirAll(staticDir, 0o755); err != nil {
		t.Fatalf("failed to create static dir: %v", err)
	}

	// Create static/robots.txt
	robotsPath := filepath.Join(staticDir, "robots.txt")
	if err := os.WriteFile(robotsPath, []byte("User-agent: *"), 0o600); err != nil {
		t.Fatalf("failed to create robots.txt: %v", err)
	}

	plugin := NewStaticFileConflictsPlugin()
	plugin.SetStaticDir(staticDir)

	m := lifecycle.NewManager()
	// Post in a subdirectory - should not conflict with root static files
	m.AddPost(&models.Post{
		Path:      "blog/posts/robots.md",
		Slug:      "blog/posts/robots",
		Published: true,
	})

	err := plugin.Collect(m)

	if err != nil {
		t.Errorf("expected no error for non-root post, got %v", err)
	}

	if len(plugin.Conflicts()) != 0 {
		t.Errorf("expected no conflicts for non-root post, got %d", len(plugin.Conflicts()))
	}
}

func TestStaticFileConflictsPlugin_MultipleConflicts(t *testing.T) {
	tmpDir := t.TempDir()
	staticDir := filepath.Join(tmpDir, "static")
	if err := os.MkdirAll(staticDir, 0o755); err != nil {
		t.Fatalf("failed to create static dir: %v", err)
	}

	// Create multiple conflicting static files
	files := map[string]string{
		"robots.txt":  "User-agent: *",
		"sitemap.xml": "<?xml version=\"1.0\"?>",
		"humans.txt":  "/* TEAM */",
	}
	for name, content := range files {
		path := filepath.Join(staticDir, name)
		if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
			t.Fatalf("failed to create %s: %v", name, err)
		}
	}

	plugin := NewStaticFileConflictsPlugin()
	plugin.SetStaticDir(staticDir)

	m := lifecycle.NewManager()
	m.AddPost(&models.Post{
		Path:      "robots.md",
		Slug:      "robots",
		Published: true,
	})
	m.AddPost(&models.Post{
		Path:      "pages/sitemap.md",
		Slug:      "sitemap",
		Published: true,
	})
	m.AddPost(&models.Post{
		Path:      "humans.md",
		Slug:      "humans",
		Published: true,
	})

	err := plugin.Collect(m)

	if err == nil {
		t.Error("expected warning about conflicts, got nil")
	}

	conflicts := plugin.Conflicts()
	if len(conflicts) < 3 {
		t.Errorf("expected at least 3 conflicts, got %d: %+v", len(conflicts), conflicts)
	}
}

func TestStaticFileConflictsPlugin_FeedConflicts(t *testing.T) {
	tmpDir := t.TempDir()
	staticDir := filepath.Join(tmpDir, "static")
	if err := os.MkdirAll(staticDir, 0o755); err != nil {
		t.Fatalf("failed to create static dir: %v", err)
	}

	// Create static/rss.xml
	rssPath := filepath.Join(staticDir, "rss.xml")
	if err := os.WriteFile(rssPath, []byte("<?xml version=\"1.0\"?>"), 0o600); err != nil {
		t.Fatalf("failed to create rss.xml: %v", err)
	}

	plugin := NewStaticFileConflictsPlugin()
	plugin.SetStaticDir(staticDir)

	m := lifecycle.NewManager()
	// Even with no posts, feed generation creates rss.xml
	// The plugin should detect this conflict

	err := plugin.Collect(m)

	if err == nil {
		t.Error("expected warning about RSS conflict, got nil")
	}

	conflicts := plugin.Conflicts()
	found := false
	for _, c := range conflicts {
		if c.OutputPath == "/rss.xml" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected conflict for /rss.xml, conflicts: %+v", conflicts)
	}
}

func TestStaticFileConflictsPlugin_WarningType(t *testing.T) {
	tmpDir := t.TempDir()
	staticDir := filepath.Join(tmpDir, "static")
	if err := os.MkdirAll(staticDir, 0o755); err != nil {
		t.Fatalf("failed to create static dir: %v", err)
	}

	robotsPath := filepath.Join(staticDir, "robots.txt")
	if err := os.WriteFile(robotsPath, []byte("User-agent: *"), 0o600); err != nil {
		t.Fatalf("failed to create robots.txt: %v", err)
	}

	plugin := NewStaticFileConflictsPlugin()
	plugin.SetStaticDir(staticDir)

	m := lifecycle.NewManager()
	m.AddPost(&models.Post{
		Path:      "pages/robots.md",
		Slug:      "robots",
		Published: true,
	})

	err := plugin.Collect(m)

	// Verify the error is a warning type using errors.As
	var warning *StaticFileConflictWarning
	if !errors.As(err, &warning) {
		t.Errorf("expected *StaticFileConflictWarning, got %T", err)
	}

	if warning != nil && !warning.IsWarning() {
		t.Error("expected IsWarning() to return true")
	}
}

func TestStaticFileConflictsPlugin_NestedStaticFiles(t *testing.T) {
	tmpDir := t.TempDir()
	staticDir := filepath.Join(tmpDir, "static")
	wellKnownDir := filepath.Join(staticDir, ".well-known")
	if err := os.MkdirAll(wellKnownDir, 0o755); err != nil {
		t.Fatalf("failed to create .well-known dir: %v", err)
	}

	// Create static/.well-known/security.txt
	securityPath := filepath.Join(wellKnownDir, "security.txt")
	if err := os.WriteFile(securityPath, []byte("Contact: security@example.com"), 0o600); err != nil {
		t.Fatalf("failed to create security.txt: %v", err)
	}

	plugin := NewStaticFileConflictsPlugin()
	plugin.SetStaticDir(staticDir)

	m := lifecycle.NewManager()
	// security.md at root level would generate /security.txt, not /.well-known/security.txt
	// So this should NOT conflict
	m.AddPost(&models.Post{
		Path:      "security.md",
		Slug:      "security",
		Published: true,
	})

	// First call - no conflict expected for .well-known path
	if err := plugin.Collect(m); err != nil {
		// If there's an error, it should not be about the .well-known path
		conflicts := plugin.Conflicts()
		for _, c := range conflicts {
			if c.OutputPath == "/.well-known/security.txt" {
				t.Errorf("unexpected conflict for .well-known path: %+v", c)
			}
		}
	}

	// The .well-known/security.txt should not conflict with /security.txt
	// But we should still detect /security.txt if it existed
	conflicts := plugin.Conflicts()
	for _, c := range conflicts {
		if c.OutputPath == "/.well-known/security.txt" {
			t.Errorf("unexpected conflict for .well-known path: %+v", c)
		}
	}

	// Now add static/security.txt to create a real conflict
	directSecurityPath := filepath.Join(staticDir, "security.txt")
	if err := os.WriteFile(directSecurityPath, []byte("Contact: security@example.com"), 0o600); err != nil {
		t.Fatalf("failed to create direct security.txt: %v", err)
	}

	// Reset and re-run
	err := plugin.Collect(m)
	if err == nil {
		t.Error("expected warning about /security.txt conflict")
	}
}

func TestStaticFileConflictWarning_Error(t *testing.T) {
	warning := &StaticFileConflictWarning{
		Message: "test warning message",
		Conflicts: []StaticFileConflict{
			{OutputPath: "/robots.txt"},
		},
	}

	if warning.Error() != "test warning message" {
		t.Errorf("expected message 'test warning message', got %q", warning.Error())
	}

	if !warning.IsWarning() {
		t.Error("expected IsWarning() to return true")
	}
}
