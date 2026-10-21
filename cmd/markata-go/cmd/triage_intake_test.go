package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func resetTriageIntakeFlags(t *testing.T) {
	t.Helper()
	beforeCategory, beforeScope, beforePost := triageCategory, triageScope, triagePost
	beforeAction, beforeFeature, beforeGoal := triageAction, triageFeature, triageGoal
	beforeIssue, beforeTo, beforeAgent, beforeNoInput := triageIssue, triageTo, triageAgent, noInput
	t.Cleanup(func() {
		triageCategory, triageScope, triagePost = beforeCategory, beforeScope, beforePost
		triageAction, triageFeature, triageGoal = beforeAction, beforeFeature, beforeGoal
		triageIssue, triageTo, triageAgent, noInput = beforeIssue, beforeTo, beforeAgent, beforeNoInput
	})
	triageCategory, triageScope, triagePost = "", "", ""
	triageAction, triageFeature, triageGoal = "", "", ""
	triageIssue, triageTo, triageAgent, noInput = "", triageStdout, "", true
}

func TestTriageIntakeInteractiveQuestionsInFixedOrder(t *testing.T) {
	resetTriageIntakeFlags(t)
	var errOutput bytes.Buffer
	cmd := &cobra.Command{Use: "triage"}
	cmd.SetIn(strings.NewReader("post\nposts/story.md\nenhance\nroad\nturn existing history into a visual timeline\n"))
	cmd.SetErr(&errOutput)
	req, issue, err := resolveTriageIntake(cmd, true)
	if err != nil {
		t.Fatalf("intake failed: %v", err)
	}
	if issue != "" || req.Category != "post" || req.Post != "posts/story.md" || req.Action != "enhance" || req.Feature != "road" || req.Goal != "turn existing history into a visual timeline" {
		t.Fatalf("unexpected intake: %#v; issue=%q", req, issue)
	}
	if len(req.Missing) != 0 {
		t.Fatalf("interactive intake left missing questions: %v", req.Missing)
	}
	questions := []string{
		"What kind of help do you need?",
		"Which Markdown post should I work on",
		"What do you want to do with this post?",
		"Which post feature is involved?",
		"What do you want to add or change?",
	}
	prev := -1
	for _, question := range questions {
		at := strings.Index(errOutput.String(), question)
		if at <= prev {
			t.Fatalf("questions out of order or absent: %s\n%s", question, errOutput.String())
		}
		prev = at
	}
}

func TestTriageIntakeSkipsSuppliedFlags(t *testing.T) {
	resetTriageIntakeFlags(t)
	triagePost = "posts/chart.md"
	triageAction = "debug"
	triageFeature = "chart"
	triageIssue = "Chart doesn't appear in the output"
	var errOutput bytes.Buffer
	cmd := &cobra.Command{Use: "post"}
	cmd.SetIn(strings.NewReader(""))
	cmd.SetErr(&errOutput)
	req, issue, err := resolveTriageIntake(cmd, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(req.Missing) != 0 || errOutput.Len() != 0 || req.Category != "post" || issue != triageIssue {
		t.Fatalf("should be fully specified: %#v, issue=%q, stderr=%q", req, issue, errOutput.String())
	}
	prompt := buildTriageIntakePrompt(triageFacts{Version: "v0.14.0"}, issue, req)
	for _, want := range []string{"- Action: debug", "- Feature: chart", "posts/chart.md", "Chart doesn't appear", "Do not repeat answered questions"} {
		if !strings.Contains(prompt, want) {
			t.Fatalf("missing %q in prompt", want)
		}
	}
}

func TestTriageIntakeNoninteractiveIsPromptOnly(t *testing.T) {
	resetTriageIntakeFlags(t)
	var stdout, stderr bytes.Buffer
	cmd := &cobra.Command{Use: "post"}
	cmd.SetIn(strings.NewReader("SHOULD NOT BE CONSUMED\n"))
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)
	if err := runTriageCommand(cmd, nil); err != nil {
		t.Fatal(err)
	}
	prompt := stdout.String()
	if stderr.Len() != 0 || !strings.HasPrefix(prompt, "# Markata troubleshooting handoff\n") {
		t.Fatalf("prompt-only contract failed: stdout=%q, stderr=%q", prompt, stderr.String())
	}
	if !strings.Contains(prompt, "Unanswered intake questions") || !strings.Contains(prompt, "Which Markdown post should I work on") ||
		!strings.Contains(prompt, "What do you want to do with this post?") || strings.Contains(prompt, "SHOULD NOT BE CONSUMED") {
		t.Fatalf("unanswered deterministic questions missing: %q", prompt)
	}
}

func TestTriageIntakeLintScopeAndConflicts(t *testing.T) {
	resetTriageIntakeFlags(t)
	triagePost = "posts/example.md"
	cmd := &cobra.Command{Use: "lint"}
	req, _, err := resolveTriageIntake(cmd, false)
	if err != nil || req.Category != "lint" || req.Scope != "post" || req.Post != "posts/example.md" || len(req.Missing) != 0 {
		t.Fatalf("lint post inference failed: %#v, %v", req, err)
	}
	triageScope = "site"
	if _, _, err := resolveTriageIntake(cmd, false); err == nil || !strings.Contains(err.Error(), "conflict") {
		t.Fatalf("expected conflicting lint scope to fail: %v", err)
	}
	triagePost, triageScope, triageAction = "", "", "enhance"
	if _, _, err := resolveTriageIntake(cmd, false); err == nil || !strings.Contains(err.Error(), "not lint") {
		t.Fatalf("expected action to be rejected for lint: %v", err)
	}
}

func TestTriageIntakeValidatesFlagsAndAnswers(t *testing.T) {
	resetTriageIntakeFlags(t)
	cmd := &cobra.Command{Use: "post"}
	triagePost, triageAction, triageFeature = "posts/x.md", "wrong", "road"
	if _, _, err := resolveTriageIntake(cmd, false); err == nil || !strings.Contains(err.Error(), "invalid --action") {
		t.Fatalf("invalid action accepted: %v", err)
	}
	triageAction, triageFeature, triageGoal = "enhance", "road", strings.Repeat("x", 4001)
	if _, _, err := resolveTriageIntake(cmd, false); err == nil || !strings.Contains(err.Error(), "4000") {
		t.Fatalf("oversized goal accepted: %v", err)
	}
	triageGoal = "visualize release history"
	triageFeature = "roads"
	if _, _, err := resolveTriageIntake(cmd, false); err == nil || !strings.Contains(err.Error(), "invalid --feature") {
		t.Fatalf("unsupported feature accepted: %v", err)
	}
}

func TestTriageSubcommandsInheritFlagsAndBypassBrokenSite(t *testing.T) {
	for _, name := range []string{"lint", "post"} {
		child, _, err := triageCmd.Find([]string{name})
		if err != nil || child == triageCmd {
			t.Fatalf("missing triage %s subcommand: %v", name, err)
		}
		if !isTriageCommand(child) {
			t.Fatalf("root pre-run would reject a broken site directory for %s", name)
		}
		for _, flag := range []string{"to", "issue", "post", "category", "scope", "action", "feature", "goal"} {
			if child.InheritedFlags().Lookup(flag) == nil {
				t.Fatalf("%s did not inherit --%s", name, flag)
			}
		}
	}
}
