package agentskill

import (
	"io/fs"
	"strings"
	"testing"
)

func TestSiteSkill_ReturnsValidFS(t *testing.T) {
	skillFS, err := SiteSkill()
	if err != nil {
		t.Fatalf("SiteSkill() error = %v", err)
	}

	f, err := skillFS.Open("SKILL.md")
	if err != nil {
		t.Fatalf("expected SKILL.md in bundled skill FS: %v", err)
	}
	f.Close()
}

func TestSiteSkill_ContainsAllRequiredTopics(t *testing.T) {
	requiredTopics := []string{
		"topics/configuration.md",
		"topics/writing-frontmatter.md",
		"topics/rich-markdown.md",
		"topics/visual-authoring.md",
		"topics/chart-authoring.md",
		"topics/cli-usage.md",
		"topics/build-deployment.md",
		"topics/faster-builds.md",
		"topics/theme-creation.md",
		"topics/template-management.md",
		"topics/plugin-creation.md",
		"topics/roads.md",
	}

	skillFS, err := SiteSkill()
	if err != nil {
		t.Fatalf("SiteSkill() error = %v", err)
	}

	for _, topic := range requiredTopics {
		t.Run(topic, func(t *testing.T) {
			info, err := fs.Stat(skillFS, topic)
			if err != nil {
				t.Fatalf("missing required topic file %q: %v", topic, err)
			}
			if info.Size() == 0 {
				t.Fatalf("topic file %q is empty", topic)
			}
		})
	}
}

func TestSiteSkill_ContainsReferenceFiles(t *testing.T) {
	requiredRefs := []string{
		"reference/template-context.md",
		"reference/feed-patterns.md",
		"reference/palette-reference.md",
	}

	skillFS, err := SiteSkill()
	if err != nil {
		t.Fatalf("SiteSkill() error = %v", err)
	}

	for _, ref := range requiredRefs {
		t.Run(ref, func(t *testing.T) {
			info, err := fs.Stat(skillFS, ref)
			if err != nil {
				t.Fatalf("missing required reference file %q: %v", ref, err)
			}
			if info.Size() == 0 {
				t.Fatalf("reference file %q is empty", ref)
			}
		})
	}
}

func TestSiteSkill_ContainsExampleFiles(t *testing.T) {
	requiredExamples := []string{
		"examples/fast.toml",
		"examples/markata-go.local.toml",
		"examples/palettes/my-brand.toml",
		"examples/templates/base.html",
		"examples/templates/post.html",
		"examples/templates/feed.html",
		"examples/editorial-visual-composition.md",
		"examples/chart-composition.md",
	}

	skillFS, err := SiteSkill()
	if err != nil {
		t.Fatalf("SiteSkill() error = %v", err)
	}

	for _, example := range requiredExamples {
		t.Run(example, func(t *testing.T) {
			info, err := fs.Stat(skillFS, example)
			if err != nil {
				t.Fatalf("missing required example file %q: %v", example, err)
			}
			if info.Size() == 0 {
				t.Fatalf("example file %q is empty", example)
			}
		})
	}
}

func TestSiteSkill_ContainsEvalFiles(t *testing.T) {
	requiredEvals := []string{
		"evals/evals.json",
	}

	skillFS, err := SiteSkill()
	if err != nil {
		t.Fatalf("SiteSkill() error = %v", err)
	}

	for _, evalFile := range requiredEvals {
		t.Run(evalFile, func(t *testing.T) {
			info, err := fs.Stat(skillFS, evalFile)
			if err != nil {
				t.Fatalf("missing required eval file %q: %v", evalFile, err)
			}
			if info.Size() == 0 {
				t.Fatalf("eval file %q is empty", evalFile)
			}
		})
	}
}

func TestListFiles_ReturnsExpectedFiles(t *testing.T) {
	files, err := ListFiles()
	if err != nil {
		t.Fatalf("ListFiles() error = %v", err)
	}

	// SKILL.md + 8 topic files + 3 reference files + 6 example files + 1 eval file = 19 minimum
	if len(files) < 19 {
		t.Fatalf("ListFiles() returned %d files, expected at least 19", len(files))
	}

	hasSkill := false
	for _, f := range files {
		if f == "SKILL.md" {
			hasSkill = true
			break
		}
	}
	if !hasSkill {
		t.Fatal("ListFiles() did not include SKILL.md")
	}
}

func TestListFiles_IsSorted(t *testing.T) {
	files, err := ListFiles()
	if err != nil {
		t.Fatalf("ListFiles() error = %v", err)
	}

	for i := 1; i < len(files); i++ {
		if files[i] < files[i-1] {
			t.Fatalf("ListFiles() not sorted: %q came after %q", files[i], files[i-1])
		}
	}
}

func TestSiteSkillName_IsExpectedValue(t *testing.T) {
	if SiteSkillName != "markata-go-site" {
		t.Fatalf("SiteSkillName = %q, want %q", SiteSkillName, "markata-go-site")
	}
}

func TestSiteSkill_EditorialVisualComposition(t *testing.T) {
	skillFS, err := SiteSkill()
	if err != nil {
		t.Fatalf("SiteSkill() error = %v", err)
	}
	data, err := fs.ReadFile(skillFS, "examples/editorial-visual-composition.md")
	if err != nil {
		t.Fatalf("cannot read visual workflow example: %v", err)
	}
	src := string(data)

	// Source tests validate structure, not a design quota. A previous
	// exact-one-feature/exact-two-cards assertion pushed agents toward
	// under-designed walls of prose rather than balanced visual storytelling.
	// Require distinct orientation, journey, choice and explanation roles.
	for _, want := range []string{
		"# From uncertainty to a clear decision",
		"## The essentials",
		"## The journey",
		"## Understand",
		"## Choose",
		"## Act",
		"## Verify",
		"::: road steps",
		"::: md-columns",
		"::: md-feature",
		"::: md-status",
		"md-card",
		"md-icon",
		":lucide-",
	} {
		if !strings.Contains(src, want) {
			t.Errorf("editorial visual example is missing %q", want)
		}
	}
	if strings.Count(src, ":::: md-card") < 2 {
		t.Error("generic visual example lost its meaningful card-based orientation or choice")
	}
	if strings.Contains(strings.ToLower(src), "repo-steward") {
		t.Error("visual composition example must remain repository-neutral")
	}

	// Markata nested containers use strictly longer fences and must close
	// with exactly the matching number of bare colons.
	stack := make([]int, 0)
	for i, line := range strings.Split(src, "\n") {
		line = strings.TrimSpace(line)
		depth := len(line) - len(strings.TrimLeft(line, ":"))
		if depth < 3 {
			continue
		}
		remaining := strings.TrimSpace(line[depth:])
		if remaining == "" {
			if len(stack) == 0 || stack[len(stack)-1] != depth {
				t.Fatalf("line %d: unmatched %d-colon closing fence", i+1, depth)
			}
			stack = stack[:len(stack)-1]
			continue
		}
		if !strings.HasPrefix(line[depth:], " ") {
			continue
		}
		if len(stack) > 0 && depth <= stack[len(stack)-1] {
			t.Fatalf("line %d: nested fence must be longer than its parent", i+1)
		}
		stack = append(stack, depth)
	}
	if len(stack) > 0 {
		t.Fatalf("unclosed Markdown containers with depths %v", stack)
	}
}
