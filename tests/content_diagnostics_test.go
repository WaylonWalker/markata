package tests

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/diagnostics"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/plugins"
)

func TestContentDiagnostics_LifecycleKeepsValidSiblings(t *testing.T) {
	contentDir := t.TempDir()
	outputDir := t.TempDir()
	writeDiagnosticFixture(t, contentDir, "good.md", "---\ntitle: Good\npublished: true\n---\n# Good")
	writeDiagnosticFixture(t, contentDir, "broken.md", "---\ntitle: Broken\npublished: true\n--- \n# Broken")
	writeDiagnosticFixture(t, contentDir, "suspicious.md", "----\ntags:\n  - example\npublished: true\n---\n# Suspicious")
	writeDiagnosticFixture(t, contentDir, "plain.md", "# Plain content")

	m := lifecycle.NewManager()
	m.SetConfig(&lifecycle.Config{
		ContentDir:   contentDir,
		OutputDir:    outputDir,
		GlobPatterns: []string{"**/*.md"},
		Extra:        map[string]interface{}{},
	})
	m.RegisterPlugins(
		plugins.NewGlobPlugin(),
		plugins.NewLoadPlugin(),
		plugins.NewRenderMarkdownPlugin(),
		plugins.NewPublishHTMLPlugin(),
	)

	if err := m.RunTo(lifecycle.StageWrite); err != nil {
		t.Fatalf("RunTo(write) error = %v", err)
	}

	snapshot := m.ContentDiagnostics()
	if snapshot.Summary.Discovered != 4 || snapshot.Summary.Candidates != 4 {
		t.Fatalf("discovery summary = %+v", snapshot.Summary)
	}
	if snapshot.Summary.Loaded != 4 || snapshot.Summary.Posts != 3 {
		t.Fatalf("load summary = %+v", snapshot.Summary)
	}
	if snapshot.Summary.Eligible != 1 || snapshot.Summary.Emitted != 3 {
		t.Fatalf("publication summary = %+v", snapshot.Summary)
	}
	if snapshot.Summary.Excluded != 3 || snapshot.Summary.Errors != 1 {
		t.Fatalf("exclusion summary = %+v", snapshot.Summary)
	}

	broken := findContentEntry(snapshot, "broken.md")
	if broken == nil {
		t.Fatal("broken.md is missing from the content ledger")
	}
	if broken.Disposition != diagnostics.DispositionExcluded {
		t.Errorf("broken.md disposition = %q, want %q", broken.Disposition, diagnostics.DispositionExcluded)
	}
	if !containsString(broken.Reasons, diagnostics.ReasonFrontmatterMalformedClosing) {
		t.Errorf("broken.md reasons = %v", broken.Reasons)
	}

	suspicious := findContentEntry(snapshot, "suspicious.md")
	if suspicious == nil {
		t.Fatal("suspicious.md is missing from the content ledger")
	}
	if !containsDiagnostic(suspicious.Diagnostics, diagnostics.ReasonFrontmatterSuspiciousDelimiter) {
		t.Errorf("suspicious.md diagnostics = %+v", suspicious.Diagnostics)
	}
	if !containsString(suspicious.Reasons, diagnostics.ReasonContentPublishedFalse) {
		t.Errorf("suspicious.md reasons = %v", suspicious.Reasons)
	}

	plain := findContentEntry(snapshot, "plain.md")
	if plain == nil {
		t.Fatal("plain.md is missing from the content ledger")
	}
	if len(plain.Diagnostics) != 0 {
		t.Errorf("plain.md diagnostics = %+v, want none", plain.Diagnostics)
	}

	if _, err := os.Stat(filepath.Join(outputDir, "good", "index.html")); err != nil {
		t.Fatalf("valid sibling output is missing: %v", err)
	}
}

func TestContentDiagnostics_FullLifecyclePublishesSanitizedArtifact(t *testing.T) {
	contentDir := t.TempDir()
	outputDir := t.TempDir()
	writeDiagnosticFixture(t, contentDir, "good.md", "---\ntitle: Good\npublished: true\n---\n# Good\n\nprivate body marker")

	m := lifecycle.NewManager()
	m.SetConfig(&lifecycle.Config{
		ContentDir:   contentDir,
		OutputDir:    outputDir,
		GlobPatterns: []string{"**/*.md"},
		Extra: map[string]interface{}{
			"markata_version": "integration-test",
			"markata_commit":  "integration-commit",
		},
	})
	m.RegisterPlugins(
		plugins.NewGlobPlugin(),
		plugins.NewLoadPlugin(),
		plugins.NewRenderMarkdownPlugin(),
		plugins.NewPublishHTMLPlugin(),
		plugins.NewDiagnosticsArtifactPlugin(),
	)

	if err := m.Run(); err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	artifactPath := filepath.Join(outputDir, diagnostics.DefaultArtifactPath)
	data, err := os.ReadFile(artifactPath)
	if err != nil {
		t.Fatalf("read diagnostics artifact: %v", err)
	}
	artifact, err := diagnostics.ParseArtifact(data)
	if err != nil {
		t.Fatalf("ParseArtifact() error = %v", err)
	}
	if artifact.Generator.Version != "integration-test" || artifact.Generator.Commit != "integration-commit" {
		t.Fatalf("generator = %#v", artifact.Generator)
	}
	if len(artifact.Entries) != 1 || artifact.Entries[0].Path != "good.md" {
		t.Fatalf("entries = %#v", artifact.Entries)
	}
	if strings.Contains(string(data), "private body marker") {
		t.Fatal("diagnostics artifact contains raw content")
	}
}

func writeDiagnosticFixture(t *testing.T, root, name, content string) {
	t.Helper()
	path := filepath.Join(root, name)
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write fixture %s: %v", name, err)
	}
}

func findContentEntry(snapshot diagnostics.ContentLedgerSnapshot, path string) *diagnostics.ContentDisposition {
	for index := range snapshot.Entries {
		if snapshot.Entries[index].Path == path {
			return &snapshot.Entries[index]
		}
	}
	return nil
}

func containsString(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}

func containsDiagnostic(issues []diagnostics.Issue, code string) bool {
	for _, issue := range issues {
		if issue.Code == code {
			return true
		}
	}
	return false
}
