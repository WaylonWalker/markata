package cmd

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/agentskill"
	"github.com/spf13/cobra"
)

func TestNormalizeAgentInstallTarget(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{name: "canonical", input: "opencode", want: "opencode"},
		{name: "legacy claude", input: "claude", want: "claude-code"},
		{name: "legacy agents", input: "agents", want: "universal"},
		{name: "unknown", input: "not-real", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := normalizeAgentInstallTarget(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("normalizeAgentInstallTarget() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if got.Name != tt.want {
				t.Fatalf("normalizeAgentInstallTarget().Name = %q, want %q", got.Name, tt.want)
			}
		})
	}
}

func TestNormalizeAgentInstallTarget_Suggestions(t *testing.T) {
	for _, tt := range []struct {
		input string
		want  string
	}{
		{input: "copilot", want: "github-copilot"},
		{input: " COPILOT ", want: "github-copilot"},
		{input: "opencod", want: "opencode"},
		{input: "cluade", want: "claude-code"},
	} {
		t.Run(tt.input, func(t *testing.T) {
			_, err := normalizeAgentInstallTarget(tt.input)
			if err == nil || !strings.Contains(err.Error(), "Did you mean") ||
				!strings.Contains(err.Error(), "--agent "+tt.want) {
				t.Fatalf("expected suggestion for %q, got %v", tt.want, err)
			}
		})
	}
	_, err := normalizeAgentInstallTarget("not-a-real-agent-xyz")
	if err == nil || strings.Contains(err.Error(), "Did you mean") ||
		!strings.Contains(err.Error(), "agent list-agents") {
		t.Fatalf("expected list guidance without an unrelated suggestion, got %v", err)
	}
}

func TestResolveAgentInstallTarget_DefaultsFromEnvironment(t *testing.T) {
	t.Setenv("OPENCODE", "1")
	t.Setenv("CLAUDECODE", "")
	t.Setenv("CLAUDE_CODE", "")
	t.Setenv("CODEX", "")
	t.Setenv("OPENAI_CODEX", "")
	t.Setenv("CURSOR_AGENT", "")
	t.Setenv("CURSOR_TRACE_ID", "")

	target, err := resolveAgentInstallTarget("", "", false)
	if err != nil {
		t.Fatalf("resolveAgentInstallTarget() error = %v", err)
	}
	if target.Name != "opencode" {
		t.Fatalf("resolveAgentInstallTarget().Name = %q, want %q", target.Name, "opencode")
	}
}

func TestResolveAgentInstallTarget_DefaultsToUniversal(t *testing.T) {
	t.Setenv("OPENCODE", "")
	t.Setenv("CLAUDECODE", "")
	t.Setenv("CLAUDE_CODE", "")
	t.Setenv("CODEX", "")
	t.Setenv("OPENAI_CODEX", "")
	t.Setenv("CURSOR_AGENT", "")
	t.Setenv("CURSOR_TRACE_ID", "")

	target, err := resolveAgentInstallTarget("", "", false)
	if err != nil {
		t.Fatalf("resolveAgentInstallTarget() error = %v", err)
	}
	if target.Name != "universal" {
		t.Fatalf("resolveAgentInstallTarget().Name = %q, want %q", target.Name, "universal")
	}
}

func TestResolveAgentInstallTarget_GlobalRequiresExplicitAgent(t *testing.T) {
	_, err := resolveAgentInstallTarget("", "", true)
	if err == nil {
		t.Fatal("expected error when --global is used without --agent")
	}
	if !strings.Contains(err.Error(), "--global requires --agent") {
		t.Fatalf("expected --global guidance, got %v", err)
	}
	for _, want := range []string{"--agent github-copilot", "agent list-agents"} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("expected %q guidance, got %v", want, err)
		}
	}
}

func TestResolveAgentInstallTarget_RejectsConflictingFlags(t *testing.T) {
	_, err := resolveAgentInstallTarget("opencode", "claude", false)
	if err == nil {
		t.Fatal("expected conflict error")
	}
	if !strings.Contains(err.Error(), "conflicts") {
		t.Fatalf("expected conflict error, got %v", err)
	}
}

func TestAgentSkillInstallDir_ProjectAndGlobal(t *testing.T) {
	projectTarget, err := normalizeAgentInstallTarget("claude-code")
	if err != nil {
		t.Fatalf("normalizeAgentInstallTarget() error = %v", err)
	}
	projectDir, err := agentSkillInstallDir("/tmp/site", agentskill.SiteSkillName, projectTarget, false)
	if err != nil {
		t.Fatalf("agentSkillInstallDir(project) error = %v", err)
	}
	if projectDir != filepath.Join(string(filepath.Separator), "tmp", "site", ".claude", "skills", agentskill.SiteSkillName) {
		t.Fatalf("projectDir = %q", projectDir)
	}

	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	t.Setenv("HOMEDRIVE", filepath.VolumeName(home))
	t.Setenv("HOMEPATH", strings.TrimPrefix(home, filepath.VolumeName(home)))
	globalTarget, err := normalizeAgentInstallTarget("opencode")
	if err != nil {
		t.Fatalf("normalizeAgentInstallTarget() error = %v", err)
	}
	globalDir, err := agentSkillInstallDir("", agentskill.SiteSkillName, globalTarget, true)
	if err != nil {
		t.Fatalf("agentSkillInstallDir(global) error = %v", err)
	}
	if globalDir != filepath.Join(home, ".config", "opencode", "skills", agentskill.SiteSkillName) {
		t.Fatalf("globalDir = %q", globalDir)
	}
}

func TestResolveAgentProjectRoot_GlobalRejectsSitePath(t *testing.T) {
	_, err := resolveAgentProjectRoot([]string{"/tmp/site"}, true)
	if err == nil {
		t.Fatal("expected error for site path with --global")
	}
	if !strings.Contains(err.Error(), "site path is not supported") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestInstallAgentSkill_DryRunDoesNotWriteFiles(t *testing.T) {
	root := t.TempDir()
	target, err := normalizeAgentInstallTarget("universal")
	if err != nil {
		t.Fatalf("normalizeAgentInstallTarget() error = %v", err)
	}
	destination, err := agentSkillInstallDir(root, agentskill.SiteSkillName, target, false)
	if err != nil {
		t.Fatalf("agentSkillInstallDir() error = %v", err)
	}

	installedFiles, err := installAgentSkill(destination, target.Name, agentScopeProject, true, false, "test")
	if err != nil {
		t.Fatalf("installAgentSkill() error = %v", err)
	}
	if len(installedFiles) == 0 {
		t.Fatal("expected bundled skill files")
	}

	if _, err := os.Stat(filepath.Join(root, ".agents")); !os.IsNotExist(err) {
		t.Fatalf("expected dry run to avoid creating files, got err=%v", err)
	}
}

func TestInstallAgentSkill_WritesBundledFiles(t *testing.T) {
	root := t.TempDir()
	target, err := normalizeAgentInstallTarget("claude-code")
	if err != nil {
		t.Fatalf("normalizeAgentInstallTarget() error = %v", err)
	}
	destination, err := agentSkillInstallDir(root, agentskill.SiteSkillName, target, false)
	if err != nil {
		t.Fatalf("agentSkillInstallDir() error = %v", err)
	}

	installedFiles, err := installAgentSkill(destination, target.Name, agentScopeProject, false, false, "test")
	if err != nil {
		t.Fatalf("installAgentSkill() error = %v", err)
	}
	if len(installedFiles) == 0 {
		t.Fatal("expected bundled skill files")
	}

	skillFile := filepath.Join(root, ".claude", "skills", agentskill.SiteSkillName, "SKILL.md")
	content, err := os.ReadFile(skillFile)
	if err != nil {
		t.Fatalf("ReadFile(SKILL.md) error = %v", err)
	}
	if !strings.Contains(string(content), "markata-go-site") {
		t.Fatalf("expected installed skill content, got %q", string(content))
	}
	if _, err := os.Stat(filepath.Join(root, ".claude", "skills", agentskill.SiteSkillName, "topics", "template-management.md")); err != nil {
		t.Fatalf("expected topic file to be installed: %v", err)
	}
}

func TestInstallAgentSkill_ExistingFileRequiresForce(t *testing.T) {
	root := t.TempDir()
	target, err := normalizeAgentInstallTarget("universal")
	if err != nil {
		t.Fatalf("normalizeAgentInstallTarget() error = %v", err)
	}
	destination, err := agentSkillInstallDir(root, agentskill.SiteSkillName, target, false)
	if err != nil {
		t.Fatalf("agentSkillInstallDir() error = %v", err)
	}
	if err := os.MkdirAll(destination, 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(destination, "SKILL.md"), []byte("existing"), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	_, err = installAgentSkill(destination, target.Name, agentScopeProject, false, false, "test")
	if err == nil {
		t.Fatal("expected conflict error when files already exist")
	}
	if !strings.Contains(err.Error(), "--force") {
		t.Fatalf("expected --force guidance, got %v", err)
	}
}

func TestRunAgentInstallCommand_DryRunUsesCommandWriter(t *testing.T) {
	stdout := bytes.NewBuffer(nil)
	stderr := bytes.NewBuffer(nil)
	command := &cobra.Command{Use: "install"}
	command.SetOut(stdout)
	command.SetErr(stderr)

	originalAgent := agentInstallAgent
	originalLegacyTarget := agentInstallLegacyTarget
	originalName := agentInstallName
	originalDryRun := agentInstallDryRun
	originalForce := agentInstallForce
	originalGlobal := agentInstallGlobal
	defer func() {
		agentInstallAgent = originalAgent
		agentInstallLegacyTarget = originalLegacyTarget
		agentInstallName = originalName
		agentInstallDryRun = originalDryRun
		agentInstallForce = originalForce
		agentInstallGlobal = originalGlobal
		currentCmd = nil
	}()

	agentInstallAgent = "opencode"
	agentInstallLegacyTarget = ""
	agentInstallName = agentskill.SiteSkillName
	agentInstallDryRun = true
	agentInstallForce = false
	agentInstallGlobal = false

	if err := runAgentInstallCommand(command, []string{t.TempDir()}); err != nil {
		t.Fatalf("runAgentInstallCommand() error = %v", err)
	}

	if !strings.Contains(stdout.String(), "Dry run: would install") {
		t.Fatalf("expected output in command writer, got %q", stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr output, got %q", stderr.String())
	}
}

func TestRunAgentListAgentsCommand_UsesCommandWriter(t *testing.T) {
	stdout := bytes.NewBuffer(nil)
	stderr := bytes.NewBuffer(nil)
	command := &cobra.Command{Use: "list-agents"}
	command.SetOut(stdout)
	command.SetErr(stderr)

	if err := runAgentListAgentsCommand(command, nil); err != nil {
		t.Fatalf("runAgentListAgentsCommand() error = %v", err)
	}

	output := stdout.String()
	if !strings.Contains(output, "Supported agents:") {
		t.Fatalf("expected heading in output, got %q", output)
	}
	if !strings.Contains(output, "- opencode") {
		t.Fatalf("expected opencode in output, got %q", output)
	}
	if !strings.Contains(output, "project: .agents/skills") {
		t.Fatalf("expected project path in output, got %q", output)
	}
	if !strings.Contains(output, "global:  .config/opencode/skills") {
		t.Fatalf("expected global path in output, got %q", output)
	}
	if !strings.Contains(output, "aliases: claude") {
		t.Fatalf("expected alias information in output, got %q", output)
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr output, got %q", stderr.String())
	}
}

func TestRunAgentInstallCommand_GlobalRequiresAgent(t *testing.T) {
	command := &cobra.Command{Use: "install"}

	originalAgent := agentInstallAgent
	originalLegacyTarget := agentInstallLegacyTarget
	originalName := agentInstallName
	originalDryRun := agentInstallDryRun
	originalForce := agentInstallForce
	originalGlobal := agentInstallGlobal
	defer func() {
		agentInstallAgent = originalAgent
		agentInstallLegacyTarget = originalLegacyTarget
		agentInstallName = originalName
		agentInstallDryRun = originalDryRun
		agentInstallForce = originalForce
		agentInstallGlobal = originalGlobal
		currentCmd = nil
	}()

	agentInstallAgent = ""
	agentInstallLegacyTarget = ""
	agentInstallName = agentskill.SiteSkillName
	agentInstallDryRun = true
	agentInstallForce = false
	agentInstallGlobal = true

	err := runAgentInstallCommand(command, nil)
	if err == nil {
		t.Fatal("expected error when --global is used without --agent")
	}
	if !strings.Contains(err.Error(), "--global requires --agent") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunAgentUpdateCommand_DryRunUsesCommandWriter(t *testing.T) {
	stdout := bytes.NewBuffer(nil)
	stderr := bytes.NewBuffer(nil)
	command := &cobra.Command{Use: "update"}
	command.SetOut(stdout)
	command.SetErr(stderr)

	originalAgent := agentUpdateAgent
	originalLegacyTarget := agentUpdateLegacyTarget
	originalName := agentUpdateName
	originalDryRun := agentUpdateDryRun
	originalGlobal := agentUpdateGlobal
	defer func() {
		agentUpdateAgent = originalAgent
		agentUpdateLegacyTarget = originalLegacyTarget
		agentUpdateName = originalName
		agentUpdateDryRun = originalDryRun
		agentUpdateGlobal = originalGlobal
		currentCmd = nil
	}()

	agentUpdateAgent = "universal"
	agentUpdateLegacyTarget = ""
	agentUpdateName = agentskill.SiteSkillName
	agentUpdateDryRun = true
	agentUpdateGlobal = false

	if err := runAgentUpdateCommand(command, []string{t.TempDir()}); err != nil {
		t.Fatalf("runAgentUpdateCommand() error = %v", err)
	}

	if !strings.Contains(stdout.String(), "Dry run: would update") {
		t.Fatalf("expected update dry-run output, got %q", stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr output, got %q", stderr.String())
	}
}

func TestRunAgentUpdateCommand_OverwritesInstalledFiles(t *testing.T) {
	root := t.TempDir()
	target, err := normalizeAgentInstallTarget("universal")
	if err != nil {
		t.Fatalf("normalizeAgentInstallTarget() error = %v", err)
	}
	destination, err := agentSkillInstallDir(root, agentskill.SiteSkillName, target, false)
	if err != nil {
		t.Fatalf("agentSkillInstallDir() error = %v", err)
	}
	_, err = installAgentSkill(destination, target.Name, agentScopeProject, false, false, Version)
	if err != nil {
		t.Fatalf("installAgentSkill() error = %v", err)
	}

	skillFile := filepath.Join(destination, "SKILL.md")
	if err := os.WriteFile(skillFile, []byte("modified"), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	stdout := bytes.NewBuffer(nil)
	stderr := bytes.NewBuffer(nil)
	command := &cobra.Command{Use: "update"}
	command.SetOut(stdout)
	command.SetErr(stderr)

	originalAgent := agentUpdateAgent
	originalLegacyTarget := agentUpdateLegacyTarget
	originalName := agentUpdateName
	originalDryRun := agentUpdateDryRun
	originalGlobal := agentUpdateGlobal
	defer func() {
		agentUpdateAgent = originalAgent
		agentUpdateLegacyTarget = originalLegacyTarget
		agentUpdateName = originalName
		agentUpdateDryRun = originalDryRun
		agentUpdateGlobal = originalGlobal
		currentCmd = nil
	}()

	agentUpdateAgent = "universal"
	agentUpdateLegacyTarget = ""
	agentUpdateName = agentskill.SiteSkillName
	agentUpdateDryRun = false
	agentUpdateGlobal = false

	if err := runAgentUpdateCommand(command, []string{root}); err != nil {
		t.Fatalf("runAgentUpdateCommand() error = %v", err)
	}

	content, err := os.ReadFile(skillFile)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if !strings.Contains(string(content), "markata-go-site") {
		t.Fatalf("expected SKILL.md to be restored by update, got %q", string(content))
	}
	if !strings.Contains(stdout.String(), "Updated ") {
		t.Fatalf("expected update output, got %q", stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr output, got %q", stderr.String())
	}
}

func TestInstallAgentSkill_WritesManifest(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	target, err := normalizeAgentInstallTarget("opencode")
	if err != nil {
		t.Fatalf("normalizeAgentInstallTarget() error = %v", err)
	}
	destination, err := agentSkillInstallDir("", agentskill.SiteSkillName, target, true)
	if err != nil {
		t.Fatalf("agentSkillInstallDir() error = %v", err)
	}
	_, err = installAgentSkill(destination, target.Name, agentScopeGlobal, false, false, "0.5.0-test")
	if err != nil {
		t.Fatalf("installAgentSkill() error = %v", err)
	}

	manifest, err := agentskill.ReadManifest(destination)
	if err != nil {
		t.Fatalf("ReadManifest() error = %v", err)
	}
	if manifest.Version != "0.5.0-test" {
		t.Errorf("Version = %q, want %q", manifest.Version, "0.5.0-test")
	}
	if manifest.Target != "opencode" {
		t.Errorf("Target = %q, want %q", manifest.Target, "opencode")
	}
	if manifest.Scope != agentScopeGlobal {
		t.Errorf("Scope = %q, want %q", manifest.Scope, agentScopeGlobal)
	}
	if len(manifest.Files) == 0 {
		t.Error("expected manifest to contain file hashes")
	}
}

func TestInstallAgentSkill_DryRunDoesNotWriteManifest(t *testing.T) {
	root := t.TempDir()
	target, err := normalizeAgentInstallTarget("universal")
	if err != nil {
		t.Fatalf("normalizeAgentInstallTarget() error = %v", err)
	}
	destination, err := agentSkillInstallDir(root, agentskill.SiteSkillName, target, false)
	if err != nil {
		t.Fatalf("agentSkillInstallDir() error = %v", err)
	}
	_, err = installAgentSkill(destination, target.Name, agentScopeProject, true, false, "0.5.0-test")
	if err != nil {
		t.Fatalf("installAgentSkill() error = %v", err)
	}

	_, err = agentskill.ReadManifest(destination)
	if err == nil {
		t.Fatal("expected manifest to not exist after dry run")
	}
}

func TestRunAgentDoctorCommand_DriftReturnsExitCode(t *testing.T) {
	root := t.TempDir()
	target, err := normalizeAgentInstallTarget("universal")
	if err != nil {
		t.Fatalf("normalizeAgentInstallTarget() error = %v", err)
	}
	destination, err := agentSkillInstallDir(root, agentskill.SiteSkillName, target, false)
	if err != nil {
		t.Fatalf("agentSkillInstallDir() error = %v", err)
	}
	_, err = installAgentSkill(destination, target.Name, agentScopeProject, false, false, "older-version")
	if err != nil {
		t.Fatalf("installAgentSkill() error = %v", err)
	}

	stdout := bytes.NewBuffer(nil)
	stderr := bytes.NewBuffer(nil)
	command := &cobra.Command{Use: "doctor"}
	command.SetOut(stdout)
	command.SetErr(stderr)

	originalAgent := agentDoctorAgent
	originalLegacyTarget := agentDoctorLegacyTarget
	originalName := agentDoctorName
	originalGlobal := agentDoctorGlobal
	defer func() {
		agentDoctorAgent = originalAgent
		agentDoctorLegacyTarget = originalLegacyTarget
		agentDoctorName = originalName
		agentDoctorGlobal = originalGlobal
		currentCmd = nil
	}()

	agentDoctorAgent = "universal"
	agentDoctorLegacyTarget = ""
	agentDoctorName = agentskill.SiteSkillName
	agentDoctorGlobal = false

	err = runAgentDoctorCommand(command, []string{root})
	if err == nil {
		t.Fatal("expected drift exit code error")
	}

	var exitErr *ExitCodeError
	if !errors.As(err, &exitErr) {
		t.Fatalf("expected ExitCodeError, got %T", err)
	}
	if exitErr.ExitCode() != doctorExitDrift {
		t.Fatalf("ExitCode() = %d, want %d", exitErr.ExitCode(), doctorExitDrift)
	}
	if !strings.Contains(stdout.String(), "Skill has") {
		t.Fatalf("expected drift output, got %q", stdout.String())
	}
}

func TestRunAgentDoctorCommand_MissingManifestReturnsExitCode(t *testing.T) {
	root := t.TempDir()
	skillDir := filepath.Join(root, ".agents", "skills", agentskill.SiteSkillName)
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte("test"), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	stdout := bytes.NewBuffer(nil)
	stderr := bytes.NewBuffer(nil)
	command := &cobra.Command{Use: "doctor"}
	command.SetOut(stdout)
	command.SetErr(stderr)

	originalAgent := agentDoctorAgent
	originalLegacyTarget := agentDoctorLegacyTarget
	originalName := agentDoctorName
	originalGlobal := agentDoctorGlobal
	defer func() {
		agentDoctorAgent = originalAgent
		agentDoctorLegacyTarget = originalLegacyTarget
		agentDoctorName = originalName
		agentDoctorGlobal = originalGlobal
		currentCmd = nil
	}()

	agentDoctorAgent = "universal"
	agentDoctorLegacyTarget = ""
	agentDoctorName = agentskill.SiteSkillName
	agentDoctorGlobal = false

	err := runAgentDoctorCommand(command, []string{root})
	if err == nil {
		t.Fatal("expected exit code error for missing manifest")
	}

	var exitErr *ExitCodeError
	if !errors.As(err, &exitErr) {
		t.Fatalf("expected ExitCodeError, got %T", err)
	}
	if exitErr.ExitCode() != doctorExitDrift {
		t.Fatalf("ExitCode() = %d, want %d", exitErr.ExitCode(), doctorExitDrift)
	}
	if !strings.Contains(stdout.String(), "has no manifest") {
		t.Fatalf("expected missing manifest output, got %q", stdout.String())
	}
}

func TestRunAgentDoctorCommand_UpToDate(t *testing.T) {
	root := t.TempDir()
	target, err := normalizeAgentInstallTarget("universal")
	if err != nil {
		t.Fatalf("normalizeAgentInstallTarget() error = %v", err)
	}
	destination, err := agentSkillInstallDir(root, agentskill.SiteSkillName, target, false)
	if err != nil {
		t.Fatalf("agentSkillInstallDir() error = %v", err)
	}
	_, err = installAgentSkill(destination, target.Name, agentScopeProject, false, false, Version)
	if err != nil {
		t.Fatalf("installAgentSkill() error = %v", err)
	}

	stdout := bytes.NewBuffer(nil)
	stderr := bytes.NewBuffer(nil)
	command := &cobra.Command{Use: "doctor"}
	command.SetOut(stdout)
	command.SetErr(stderr)

	originalAgent := agentDoctorAgent
	originalLegacyTarget := agentDoctorLegacyTarget
	originalName := agentDoctorName
	originalGlobal := agentDoctorGlobal
	defer func() {
		agentDoctorAgent = originalAgent
		agentDoctorLegacyTarget = originalLegacyTarget
		agentDoctorName = originalName
		agentDoctorGlobal = originalGlobal
		currentCmd = nil
	}()

	agentDoctorAgent = "universal"
	agentDoctorLegacyTarget = ""
	agentDoctorName = agentskill.SiteSkillName
	agentDoctorGlobal = false

	if err := runAgentDoctorCommand(command, []string{root}); err != nil {
		t.Fatalf("runAgentDoctorCommand() error = %v", err)
	}

	if !strings.Contains(stdout.String(), "up to date") {
		t.Errorf("expected 'up to date' in output, got %q", stdout.String())
	}
	if !strings.Contains(stdout.String(), "Agent:     universal") {
		t.Errorf("expected agent line in output, got %q", stdout.String())
	}
}

func TestRunAgentRemoveCommand_RemovesInstalledSkill(t *testing.T) {
	root := t.TempDir()
	target, err := normalizeAgentInstallTarget("claude-code")
	if err != nil {
		t.Fatalf("normalizeAgentInstallTarget() error = %v", err)
	}
	destination, err := agentSkillInstallDir(root, agentskill.SiteSkillName, target, false)
	if err != nil {
		t.Fatalf("agentSkillInstallDir() error = %v", err)
	}
	_, err = installAgentSkill(destination, target.Name, agentScopeProject, false, false, Version)
	if err != nil {
		t.Fatalf("installAgentSkill() error = %v", err)
	}

	stdout := bytes.NewBuffer(nil)
	stderr := bytes.NewBuffer(nil)
	command := &cobra.Command{Use: "remove"}
	command.SetOut(stdout)
	command.SetErr(stderr)

	originalAgent := agentRemoveAgent
	originalLegacyTarget := agentRemoveLegacyTarget
	originalName := agentRemoveName
	originalGlobal := agentRemoveGlobal
	defer func() {
		agentRemoveAgent = originalAgent
		agentRemoveLegacyTarget = originalLegacyTarget
		agentRemoveName = originalName
		agentRemoveGlobal = originalGlobal
		currentCmd = nil
	}()

	agentRemoveAgent = "claude-code"
	agentRemoveLegacyTarget = ""
	agentRemoveName = agentskill.SiteSkillName
	agentRemoveGlobal = false

	if err := runAgentRemoveCommand(command, []string{root}); err != nil {
		t.Fatalf("runAgentRemoveCommand() error = %v", err)
	}

	if _, err := os.Stat(destination); !os.IsNotExist(err) {
		t.Fatalf("expected skill dir to be removed, got err=%v", err)
	}
	if !strings.Contains(stdout.String(), "Removed .claude/skills/") {
		t.Fatalf("expected remove output, got %q", stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr output, got %q", stderr.String())
	}
}

func TestAgentRemoveCmd_HasUninstallAlias(t *testing.T) {
	found := false
	for _, alias := range agentRemoveCmd.Aliases {
		if alias == "uninstall" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("expected uninstall alias on agent remove command")
	}
}

func TestValidateSkillName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{name: "valid name", input: "markata-go-site", wantErr: false},
		{name: "valid simple", input: "my-skill", wantErr: false},
		{name: "empty", input: "", wantErr: true},
		{name: "dot dot", input: "..", wantErr: true},
		{name: "dot dot prefix", input: "..foo", wantErr: true},
		{name: "single dot", input: ".", wantErr: true},
		{name: "forward slash", input: "foo/bar", wantErr: true},
		{name: "backslash", input: "foo\\bar", wantErr: true},
		{name: "traversal path", input: "../../etc", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateSkillName(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("validateSkillName(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestRunAgentRemoveCommand_RejectsTraversalName(t *testing.T) {
	root := t.TempDir()
	command := &cobra.Command{Use: "remove"}

	originalAgent := agentRemoveAgent
	originalLegacyTarget := agentRemoveLegacyTarget
	originalName := agentRemoveName
	originalGlobal := agentRemoveGlobal
	defer func() {
		agentRemoveAgent = originalAgent
		agentRemoveLegacyTarget = originalLegacyTarget
		agentRemoveName = originalName
		agentRemoveGlobal = originalGlobal
		currentCmd = nil
	}()

	agentRemoveAgent = "universal"
	agentRemoveLegacyTarget = ""
	agentRemoveName = "../.."
	agentRemoveGlobal = false

	err := runAgentRemoveCommand(command, []string{root})
	if err == nil {
		t.Fatal("expected error for path traversal name")
	}
	if !strings.Contains(err.Error(), "path separator") && !strings.Contains(err.Error(), "path traversal") {
		t.Fatalf("expected path validation error, got %v", err)
	}
}

func TestRunAgentRemoveCommand_RejectsNonSkillDirectory(t *testing.T) {
	root := t.TempDir()
	notASkill := filepath.Join(root, ".agents", "skills", "not-a-skill")
	if err := os.MkdirAll(notASkill, 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(notASkill, "random.txt"), []byte("hi"), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	command := &cobra.Command{Use: "remove"}

	originalAgent := agentRemoveAgent
	originalLegacyTarget := agentRemoveLegacyTarget
	originalName := agentRemoveName
	originalGlobal := agentRemoveGlobal
	defer func() {
		agentRemoveAgent = originalAgent
		agentRemoveLegacyTarget = originalLegacyTarget
		agentRemoveName = originalName
		agentRemoveGlobal = originalGlobal
		currentCmd = nil
	}()

	agentRemoveAgent = "universal"
	agentRemoveLegacyTarget = ""
	agentRemoveName = "not-a-skill"
	agentRemoveGlobal = false

	err := runAgentRemoveCommand(command, []string{root})
	if err == nil {
		t.Fatal("expected error for non-skill directory")
	}
	if !strings.Contains(err.Error(), "does not appear to be an installed skill") {
		t.Fatalf("expected skill marker error, got %v", err)
	}

	if _, statErr := os.Stat(notASkill); os.IsNotExist(statErr) {
		t.Fatal("expected non-skill directory to be preserved")
	}
}
