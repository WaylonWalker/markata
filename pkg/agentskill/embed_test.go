package agentskill

import (
	"io/fs"
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
