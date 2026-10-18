package plugins

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
)

func TestPythonDocsPlugin_Name(t *testing.T) {
	p := NewPythonDocsPlugin()
	if got := p.Name(); got != "python_docs" {
		t.Errorf("Name() = %q, want %q", got, "python_docs")
	}
}

func TestPythonDocsPlugin_DefaultConfig(t *testing.T) {
	p := NewPythonDocsPlugin()
	if p.config.Enabled {
		t.Error("default Enabled = true, want false")
	}
	if p.config.SlugPrefix != "api" {
		t.Errorf("default SlugPrefix = %q, want %q", p.config.SlugPrefix, "api")
	}
	if !p.config.IncludeSource {
		t.Error("default IncludeSource = false, want true")
	}
	if p.config.IncludeModuleCode {
		t.Error("default IncludeModuleCode = true, want false")
	}
}

func TestPythonDocsPlugin_LoadCreatesPosts(t *testing.T) {
	requirePythonInterpreter(t)

	root := t.TempDir()
	writeTestFile(t, root, "pkg/util.py", "\"\"\"Utility helpers.\"\"\"\n\n"+
		"def greet(name: str) -> str:\n"+
		"    \"\"\"Return a greeting from `pkg.util.greet`.\"\"\"\n"+
		"    message = f\"Hello {name}\"\n"+
		"    return message\n")
	writeTestFile(t, root, "pkg/core.py", "\"\"\"Core module.\n\n"+
		"Uses `pkg.util` for shared helpers.\n"+
		"\"\"\"\n\n"+
		"from pkg.util import greet\n\n\n"+
		"class Runner:\n"+
		"    \"\"\"Coordinates greeting.\"\"\"\n\n"+
		"    def run(self, name: str) -> str:\n"+
		"        \"\"\"Run the greeter.\"\"\"\n"+
		"        return greet(name)\n")

	m := lifecycle.NewManager()
	m.SetConfig(&lifecycle.Config{
		ContentDir: root,
		Extra: map[string]any{
			"hooks": []any{"default", "python_docs"},
			"python_docs": map[string]any{
				"enabled":     true,
				"patterns":    []any{"pkg/**/*.py"},
				"slug_prefix": "reference",
			},
		},
	})

	p := NewPythonDocsPlugin()
	if err := p.Configure(m); err != nil {
		t.Fatalf("Configure() error = %v", err)
	}
	if err := p.Load(m); err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	posts := m.Posts()
	if len(posts) != 2 {
		t.Fatalf("len(posts) = %d, want 2", len(posts))
	}

	byModule := make(map[string]string, len(posts))
	for _, post := range posts {
		value := post.Get("python_module")
		module, ok := value.(string)
		if !ok {
			t.Fatalf("python_module = %T, want string", value)
		}
		byModule[module] = post.Content
		if !strings.HasPrefix(post.Slug, "reference/") {
			t.Errorf("post slug = %q, want prefix reference/", post.Slug)
		}
	}

	utilContent := byModule["pkg.util"]
	if !strings.Contains(utilContent, "Utility helpers.") {
		t.Errorf("util content missing module docstring: %q", utilContent)
	}
	if !strings.Contains(utilContent, "Function `greet`") {
		t.Errorf("util content missing function heading: %q", utilContent)
	}
	if !strings.Contains(utilContent, "message = f\"Hello {name}\"") {
		t.Errorf("util content missing implementation source: %q", utilContent)
	}
	if strings.Contains(utilContent, `"""Return a greeting`) {
		t.Errorf("util implementation should omit inline docstring, got: %q", utilContent)
	}

	coreContent := byModule["pkg.core"]
	if !strings.Contains(coreContent, "[`pkg.util`](/reference/pkg/util/)") {
		t.Errorf("core content missing internal module link: %q", coreContent)
	}
	if !strings.Contains(coreContent, "from [`pkg.util`](/reference/pkg/util/) import [`pkg.util.greet`](/reference/pkg/util/#symbol-greet)") {
		t.Errorf("core content missing internal symbol import link: %q", coreContent)
	}
	if !strings.Contains(coreContent, "Method `run`") {
		t.Errorf("core content missing method heading: %q", coreContent)
	}
	if !strings.Contains(utilContent, "Related: [`pkg.util.greet`](/reference/pkg/util/#symbol-greet)") {
		t.Errorf("util content missing related symbol links: %q", utilContent)
	}
}

func TestPythonDocsPlugin_RelativeImportsInInitModule(t *testing.T) {
	requirePythonInterpreter(t)

	root := t.TempDir()
	writeTestFile(t, root, "pkg/__init__.py", "from .config import Settings\n")
	writeTestFile(t, root, "pkg/config.py", "class Settings:\n    \"\"\"Package settings.\"\"\"\n    pass\n")

	m := lifecycle.NewManager()
	m.SetConfig(&lifecycle.Config{
		ContentDir: root,
		Extra: map[string]any{
			"hooks": []any{"default", "python_docs"},
			"python_docs": map[string]any{
				"enabled":     true,
				"patterns":    []any{"pkg/**/*.py"},
				"slug_prefix": "reference",
			},
		},
	})

	p := NewPythonDocsPlugin()
	if err := p.Configure(m); err != nil {
		t.Fatalf("Configure() error = %v", err)
	}
	if err := p.Load(m); err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	for _, post := range m.Posts() {
		value := post.Get("python_module")
		module, ok := value.(string)
		if !ok {
			t.Fatalf("python_module = %T, want string", value)
		}
		if module != "pkg" {
			continue
		}
		if !strings.Contains(post.Content, "from [`pkg.config`](/reference/pkg/config/) import [`pkg.config.Settings`](/reference/pkg/config/#symbol-settings)") {
			t.Fatalf("init module relative import not linked: %q", post.Content)
		}
		return
	}

	t.Fatal("pkg post not found")
}

func TestPythonDocsPlugin_RequiresExplicitHook(t *testing.T) {
	requirePythonInterpreter(t)

	root := t.TempDir()
	writeTestFile(t, root, "pkg/util.py", "def greet(name: str) -> str:\n    return name\n")

	m := lifecycle.NewManager()
	m.SetConfig(&lifecycle.Config{
		ContentDir: root,
		Extra: map[string]any{
			"python_docs": map[string]any{
				"enabled":  true,
				"patterns": []any{"pkg/**/*.py"},
			},
		},
	})

	p := NewPythonDocsPlugin()
	if err := p.Configure(m); err != nil {
		t.Fatalf("Configure() error = %v", err)
	}
	if err := p.Load(m); err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if got := len(m.Posts()); got != 0 {
		t.Fatalf("len(posts) = %d, want 0 when hook not explicitly enabled", got)
	}
	if p.strictHookOptIn {
		t.Fatal("strictHookOptIn = true, want false")
	}
	if !p.config.Enabled {
		t.Fatal("config.Enabled = false, want true")
	}
}

func TestPythonDocsPlugin_CustomSymbolTemplate(t *testing.T) {
	requirePythonInterpreter(t)

	root := t.TempDir()
	writeTestFile(t, root, "pkg/util.py", "def greet(name: str) -> str:\n    \"\"\"Uses pkg.types.User.\"\"\"\n    return name\n")
	writeTestFile(t, root, "pkg/types.py", "class User:\n    pass\n")

	m := lifecycle.NewManager()
	m.SetConfig(&lifecycle.Config{
		ContentDir: root,
		Extra: map[string]any{
			"hooks": []any{"default", "python_docs"},
			"python_docs": map[string]any{
				"enabled":         true,
				"patterns":        []any{"pkg/**/*.py"},
				"slug_prefix":     "reference",
				"symbol_template": "{{ heading_prefix }} Custom `{{ item.name }}`\n\nRefs: {{ references | join:\", \" }}\n",
			},
		},
	})

	p := NewPythonDocsPlugin()
	if err := p.Configure(m); err != nil {
		t.Fatalf("Configure() error = %v", err)
	}
	if err := p.Load(m); err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	posts := m.Posts()
	if len(posts) != 2 {
		t.Fatalf("len(posts) = %d, want 2", len(posts))
	}

	for _, post := range posts {
		value := post.Get("python_module")
		module, ok := value.(string)
		if !ok {
			t.Fatalf("python_module = %T, want string", value)
		}
		if module != "pkg.util" {
			continue
		}
		if !strings.Contains(post.Content, "## Custom `greet`") {
			t.Fatalf("custom template not used: %q", post.Content)
		}
		if !strings.Contains(post.Content, "[`pkg.types.User`](/reference/pkg/types/#symbol-user)") {
			t.Fatalf("custom template missing resolved references: %q", post.Content)
		}
		return
	}

	t.Fatal("pkg.util post not found")
}

func requirePythonInterpreter(t *testing.T) {
	t.Helper()
	for _, candidate := range []string{"python3", "python"} {
		if _, err := exec.LookPath(candidate); err == nil {
			return
		}
	}
	t.Skip("python interpreter not available")
}

func writeTestFile(t *testing.T, root, relPath, content string) {
	t.Helper()
	fullPath := filepath.Join(root, relPath)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(fullPath, []byte(content), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
}
