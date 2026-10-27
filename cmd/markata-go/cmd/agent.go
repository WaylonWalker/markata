package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/WaylonWalker/markata/pkg/agentskill"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Install agent integrations for markata-go sites",
	Long: `Manage installable agent integrations for markata-go sites.

The initial workflow installs a bundled, topic-based skill into either a project
repository or an agent-specific global skill directory. Agent identifiers and
install locations follow the same naming model documented by vercel-labs/skills.

Examples:
  markata-go agent install
  markata-go agent install --agent claude-code
  markata-go agent install --agent opencode --global`,
}

var agentInstallCmd = &cobra.Command{
	Use:     "install [site-path]",
	Aliases: []string{"setup"},
	Short:   "Install the bundled markata-go site skill",
	Long: `Install the bundled markata-go site skill.

When --agent is omitted, markata-go installs into the current agent's project
layout when it can detect one from the environment. Otherwise it falls back to
the portable universal layout under .agents/skills.

Use --global to install into an agent's user-level skill directory. On an
interactive terminal without --agent, markata-go detects installed agents and
offers to install for all of them. Declining opens a fuzzy picker. Scripts
retain the environment-detected target or universal project layout; global
installs require an explicit --agent when prompting is unavailable.

Examples:
  markata-go agent install
  markata-go agent setup
  markata-go agent install --agent claude-code
  markata-go agent install --global
  markata-go agent install --agent opencode --global
  markata-go agent install ../my-site --dry-run
  markata-go agent install . --force`,
	Args: cobra.MaximumNArgs(1),
	RunE: runAgentInstallCommand,
}

var agentListAgentsCmd = &cobra.Command{
	Use:   "list-agents",
	Short: "List supported agent identifiers and install paths",
	Long: `List the supported agent identifiers that markata-go can target.

The output includes each agent's project-level and global skill directory so you
can choose an explicit --agent value without scanning help text.

Examples:
  markata-go agent list-agents
  markata-go agent list-agents | grep opencode`,
	Args: cobra.NoArgs,
	RunE: runAgentListAgentsCommand,
}

var agentDoctorCmd = &cobra.Command{
	Use:   "doctor [site-path]",
	Short: "Check installed skill for drift against the current binary",
	Long: `Compare installed skill files against the versions bundled in the current
markata-go binary. Reports files that are modified, missing, or newly added
since the skill was installed.

Examples:
  markata-go agent doctor
  markata-go agent doctor --agent claude-code
  markata-go agent doctor --agent opencode --global
  markata-go agent doctor ../my-site`,
	Args: cobra.MaximumNArgs(1),
	RunE: runAgentDoctorCommand,
}

var agentUpdateCmd = &cobra.Command{
	Use:   "update [site-path]",
	Short: "Update an installed markata-go site skill",
	Long: `Update the installed markata-go site skill in place.

This is the user-friendly equivalent of running 'markata-go agent install --force'.

Examples:
  markata-go agent update
  markata-go agent update --agent claude-code
  markata-go agent update --agent opencode --global --dry-run
  markata-go agent update ../my-site`,
	Args: cobra.MaximumNArgs(1),
	RunE: runAgentUpdateCommand,
}

var agentRemoveCmd = &cobra.Command{
	Use:     "remove [site-path]",
	Aliases: []string{"uninstall"},
	Short:   "Remove an installed markata-go site skill",
	Long: `Remove the installed markata-go site skill.

Examples:
  markata-go agent remove
  markata-go agent uninstall
  markata-go agent remove --agent claude-code
  markata-go agent remove --agent opencode --global
  markata-go agent remove ../my-site`,
	Args: cobra.MaximumNArgs(1),
	RunE: runAgentRemoveCommand,
}

var (
	agentInstallAgent        string
	agentInstallLegacyTarget string
	agentInstallName         string
	agentInstallForce        bool
	agentInstallDryRun       bool
	agentInstallGlobal       bool

	agentDoctorAgent        string
	agentDoctorLegacyTarget string
	agentDoctorName         string
	agentDoctorGlobal       bool

	agentUpdateAgent        string
	agentUpdateLegacyTarget string
	agentUpdateName         string
	agentUpdateDryRun       bool
	agentUpdateGlobal       bool

	agentRemoveAgent        string
	agentRemoveLegacyTarget string
	agentRemoveName         string
	agentRemoveGlobal       bool
)

func init() {
	rootCmd.AddCommand(agentCmd)
	agentCmd.AddCommand(agentInstallCmd)
	agentCmd.AddCommand(agentListAgentsCmd)
	agentCmd.AddCommand(agentUpdateCmd)
	agentCmd.AddCommand(agentDoctorCmd)
	agentCmd.AddCommand(agentRemoveCmd)

	installFlags := agentInstallCmd.Flags()
	installFlags.StringVar(&agentInstallAgent, "agent", "", fmt.Sprintf("target agent (%s)", strings.Join(supportedAgentNames(), ", ")))
	installFlags.StringVarP(&agentInstallLegacyTarget, "target", "", "", "legacy alias for --agent")
	installFlags.BoolVarP(&agentInstallGlobal, "global", "g", false, "install into the selected agent's user-level skill directory")
	installFlags.StringVar(&agentInstallName, "name", agentskill.SiteSkillName, "installed skill directory name")
	installFlags.BoolVar(&agentInstallForce, "force", false, "overwrite bundled skill files if they already exist")
	installFlags.BoolVar(&agentInstallDryRun, "dry-run", false, "show what would be installed without writing files")
	markTargetHidden(installFlags)

	updateFlags := agentUpdateCmd.Flags()
	updateFlags.StringVar(&agentUpdateAgent, "agent", "", fmt.Sprintf("target agent (%s)", strings.Join(supportedAgentNames(), ", ")))
	updateFlags.StringVarP(&agentUpdateLegacyTarget, "target", "", "", "legacy alias for --agent")
	updateFlags.BoolVarP(&agentUpdateGlobal, "global", "g", false, "update the selected agent's user-level skill directory")
	updateFlags.StringVar(&agentUpdateName, "name", agentskill.SiteSkillName, "installed skill directory name")
	updateFlags.BoolVar(&agentUpdateDryRun, "dry-run", false, "show what would be updated without writing files")
	markTargetHidden(updateFlags)

	doctorFlags := agentDoctorCmd.Flags()
	doctorFlags.StringVar(&agentDoctorAgent, "agent", "", fmt.Sprintf("target agent (%s)", strings.Join(supportedAgentNames(), ", ")))
	doctorFlags.StringVarP(&agentDoctorLegacyTarget, "target", "", "", "legacy alias for --agent")
	doctorFlags.BoolVarP(&agentDoctorGlobal, "global", "g", false, "check the selected agent's user-level skill directory")
	doctorFlags.StringVar(&agentDoctorName, "name", agentskill.SiteSkillName, "installed skill directory name")
	markTargetHidden(doctorFlags)

	removeFlags := agentRemoveCmd.Flags()
	removeFlags.StringVar(&agentRemoveAgent, "agent", "", fmt.Sprintf("target agent (%s)", strings.Join(supportedAgentNames(), ", ")))
	removeFlags.StringVarP(&agentRemoveLegacyTarget, "target", "", "", "legacy alias for --agent")
	removeFlags.BoolVarP(&agentRemoveGlobal, "global", "g", false, "remove from the selected agent's user-level skill directory")
	removeFlags.StringVar(&agentRemoveName, "name", agentskill.SiteSkillName, "installed skill directory name")
	markTargetHidden(removeFlags)
}

func markTargetHidden(flags *pflag.FlagSet) {
	if err := flags.MarkHidden("target"); err != nil {
		// Flag existence errors are non-fatal; log in debug mode only
		_ = err
	}
}

func runAgentInstallCommand(cmd *cobra.Command, args []string) error {
	currentCmd = cmd
	if _, err := resolveAgentProjectRoot(args, agentInstallGlobal); err != nil {
		return err
	}
	selected, explicit, err := selectedAgentFlagValue(agentInstallAgent, agentInstallLegacyTarget)
	if err != nil {
		return err
	}
	if explicit {
		target, err := normalizeAgentInstallTarget(selected)
		if err != nil {
			return err
		}
		return runAgentWriteCommand(args, target, agentInstallGlobal, agentInstallName, agentInstallDryRun, agentInstallForce, "install")
	}
	interactive := !noInput && os.Getenv("TERM") != triageDumb &&
		streamIsTerminal(cmd.InOrStdin()) && streamIsTerminal(cmd.OutOrStdout())
	if !interactive {
		if agentInstallGlobal {
			return fmt.Errorf("--global requires --agent without an interactive picker; use '--agent github-copilot' or run 'markata-go agent list-agents' for supported values")
		}
		target, err := resolveAgentInstallTarget("", "", false)
		if err != nil {
			return err
		}
		return runAgentWriteCommand(args, target, false, agentInstallName, agentInstallDryRun, agentInstallForce, "install")
	}
	detected := detectInstalledAgentTargets()
	if len(detected) > 0 {
		outlnf("Detected %d installed agent(s):", len(detected))
		for _, target := range detected {
			outlnf("  - %s", target.Name)
		}
		outlnf("Install Markata skills for all detected agents? [Y/n]")
		reader := bufio.NewReader(cmd.InOrStdin())
		answer, readErr := reader.ReadString('\n')
		if readErr != nil {
			if errors.Is(readErr, io.EOF) {
				return fmt.Errorf("agent setup cancelled: input ended before confirmation")
			}
			return readErr
		}
		switch strings.ToLower(strings.TrimSpace(answer)) {
		case "", "y", "yes":
			seen := make(map[string]bool)
			for _, target := range detected {
				path := target.ProjectPath
				if agentInstallGlobal {
					path = target.GlobalPath
				}
				if seen[path] {
					continue
				}
				seen[path] = true
				if err := runAgentWriteCommand(args, target, agentInstallGlobal, agentInstallName, agentInstallDryRun, agentInstallForce, "install"); err != nil {
					return fmt.Errorf("%s: %w", target.Name, err)
				}
			}
			return nil
		case "n", "no":
			// Continue to the existing searchable agent picker.
		default:
			return fmt.Errorf("invalid response %q; enter y or n", strings.TrimSpace(answer))
		}
	} else {
		outlnf("No supported installed agents detected; choose an agent manually.")
	}
	selected, err = pickInstallAgent(cmd)
	if err != nil {
		return err
	}
	target, err := normalizeAgentInstallTarget(selected)
	if err != nil {
		return err
	}
	return runAgentWriteCommand(args, target, agentInstallGlobal, agentInstallName, agentInstallDryRun, agentInstallForce, "install")
}

// Detection uses agent-specific executables and user configuration, not the
// universal .agents directory shared by multiple agent implementations.
func detectInstalledAgentTargets() []agentInstallTarget {
	binaries := map[string][]string{
		"opencode":       {"opencode"},
		"claude-code":    {"claude"},
		"codex":          {"codex"},
		"cursor":         {"cursor", "cursor-agent"},
		"gemini-cli":     {"gemini"},
		"github-copilot": {"copilot"},
		"goose":          {"goose"},
		"cline":          {"cline"},
		"amp":            {"amp"},
		"qwen-code":      {"qwen"},
		"kimi-cli":       {"kimi"},
		"crush":          {"crush"},
		"continue":       {"cn"},
	}
	home, homeErr := os.UserHomeDir()
	if homeErr != nil {
		home = ""
	}
	var found []agentInstallTarget
	for _, target := range supportedAgentTargets {
		if target.Name == agentDefaultUniversal {
			continue
		}
		present := false
		for _, binary := range binaries[target.Name] {
			if _, err := exec.LookPath(binary); err == nil {
				present = true
				break
			}
		}
		if !present && home != "" {
			// Ignore generic/shared paths: they do not identify an agent.
			if target.GlobalPath != ".agents/skills" && target.GlobalPath != ".config/agents/skills" {
				_, err := os.Stat(filepath.Join(home, filepath.Dir(target.GlobalPath)))
				present = err == nil
			}
		}
		if present {
			found = append(found, target)
		}
	}
	return found
}

func runAgentListAgentsCommand(cmd *cobra.Command, _ []string) error {
	currentCmd = cmd

	outlnf("Supported agents:")
	for _, target := range supportedAgentTargets {
		outlnf("- %s", target.Name)
		outlnf("  project: %s", target.ProjectPath)
		outlnf("  global:  %s", target.GlobalPath)
		if len(target.Aliases) > 0 {
			outlnf("  aliases: %s", strings.Join(target.Aliases, ", "))
		}
	}

	return nil
}

func runAgentUpdateCommand(cmd *cobra.Command, args []string) error {
	currentCmd = cmd
	target, err := resolveAgentInstallTarget(agentUpdateAgent, agentUpdateLegacyTarget, agentUpdateGlobal)
	if err != nil {
		return err
	}
	return runAgentWriteCommand(args, target, agentUpdateGlobal, agentUpdateName, agentUpdateDryRun, true, "update")
}

func runAgentWriteCommand(args []string, target agentInstallTarget, global bool, name string, dryRun, force bool, operation string) error {
	root, err := resolveAgentProjectRoot(args, global)
	if err != nil {
		return err
	}

	if err := validateSkillName(name); err != nil {
		return err
	}

	destination, err := agentSkillInstallDir(root, name, target, global)
	if err != nil {
		return err
	}

	installedFiles, err := installAgentSkill(destination, target.Name, agentScopeName(global), dryRun, force, Version)
	if err != nil {
		return err
	}

	if dryRun {
		outlnf("%s would %s %d file(s) %s %s", agentYellow("Dry run:"), operation, len(installedFiles), prepositionForAgentOperation(operation), destination)
	} else {
		outlnf("%s %d file(s) %s %s", agentGreen("%s", pastTenseForAgentOperation(operation)), len(installedFiles), prepositionForAgentOperation(operation), destination)
	}

	for _, filePath := range installedFiles {
		if !global && root != "" {
			relPath, relErr := filepath.Rel(root, filePath)
			if relErr == nil {
				filePath = filepath.ToSlash(relPath)
			} else {
				filePath = filepath.ToSlash(filePath)
			}
		} else {
			filePath = filepath.ToSlash(filePath)
		}
		outlnf("- %s", filePath)
	}

	return nil
}

func resolveAgentProjectRoot(args []string, global bool) (string, error) {
	if global {
		if len(args) > 0 {
			return "", fmt.Errorf("site path is not supported with --global")
		}
		return "", nil
	}

	sitePath := "."
	if len(args) > 0 {
		sitePath = args[0]
	}

	root, err := filepath.Abs(sitePath)
	if err != nil {
		return "", fmt.Errorf("resolve site path %q: %w", sitePath, err)
	}

	info, err := os.Stat(root)
	if err != nil {
		return "", fmt.Errorf("stat site path %q: %w", root, err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("site path %q is not a directory", root)
	}

	return root, nil
}

func prepositionForAgentOperation(operation string) string {
	switch operation {
	case "update":
		return "in"
	default:
		return "to"
	}
}

func pastTenseForAgentOperation(operation string) string {
	switch operation {
	case "update":
		return "Updated"
	default:
		return "Installed"
	}
}

// validateSkillName rejects names containing path separators or traversal components.
func validateSkillName(name string) error {
	if name == "" {
		return fmt.Errorf("skill name must not be empty")
	}
	if strings.Contains(name, "/") || strings.Contains(name, "\\") {
		return fmt.Errorf("skill name %q must not contain path separators", name)
	}
	if name == ".." || strings.HasPrefix(name, "..") {
		return fmt.Errorf("skill name %q must not contain path traversal components", name)
	}
	if name == "." {
		return fmt.Errorf("skill name %q is not valid", name)
	}
	return nil
}

func installAgentSkill(destinationRoot, agentName, scope string, dryRun, force bool, version string) ([]string, error) {
	bundledFS, err := agentskill.SiteSkill()
	if err != nil {
		return nil, fmt.Errorf("load bundled skill: %w", err)
	}

	relativeFiles, err := agentskill.ListFiles()
	if err != nil {
		return nil, fmt.Errorf("list bundled skill files: %w", err)
	}

	installedFiles := make([]string, 0, len(relativeFiles))
	for _, relativePath := range relativeFiles {
		installedFiles = append(installedFiles, filepath.Join(destinationRoot, filepath.FromSlash(relativePath)))
	}

	if !force {
		for _, filePath := range installedFiles {
			_, statErr := os.Stat(filePath)
			if statErr == nil {
				return nil, fmt.Errorf("refusing to overwrite existing file %q; re-run with --force", filePath)
			}
			if statErr != nil && !errors.Is(statErr, os.ErrNotExist) {
				return nil, fmt.Errorf("stat destination %q: %w", filePath, statErr)
			}
		}
	}

	if dryRun {
		return installedFiles, nil
	}

	for i, relativePath := range relativeFiles {
		data, readErr := fs.ReadFile(bundledFS, relativePath)
		if readErr != nil {
			return nil, fmt.Errorf("read bundled skill file %q: %w", relativePath, readErr)
		}

		filePath := installedFiles[i]
		if err := os.MkdirAll(filepath.Dir(filePath), 0o755); err != nil {
			return nil, fmt.Errorf("create destination directory for %q: %w", filePath, err)
		}
		if err := os.WriteFile(filePath, data, 0o644); err != nil { //nolint:gosec // installed skill files are non-sensitive markdown
			return nil, fmt.Errorf("write bundled skill file %q: %w", filePath, err)
		}
	}

	manifest, err := agentskill.ComputeBundledManifest(version, agentName, scope)
	if err != nil {
		return nil, fmt.Errorf("compute manifest: %w", err)
	}
	if err := agentskill.WriteManifest(manifest, destinationRoot); err != nil {
		return nil, fmt.Errorf("write manifest: %w", err)
	}

	return installedFiles, nil
}

// doctorExitDrift is the exit code when drift is detected.
const doctorExitDrift = 1

// doctorExitError is the exit code when doctor cannot complete.
const doctorExitError = 2

func runAgentDoctorCommand(cmd *cobra.Command, args []string) error {
	currentCmd = cmd

	target, err := resolveAgentInstallTarget(agentDoctorAgent, agentDoctorLegacyTarget, agentDoctorGlobal)
	if err != nil {
		return newExitCodeError(doctorExitError, err)
	}

	root, err := resolveAgentProjectRoot(args, agentDoctorGlobal)
	if err != nil {
		return newExitCodeError(doctorExitError, err)
	}

	if err := validateSkillName(agentDoctorName); err != nil {
		return newExitCodeError(doctorExitError, err)
	}

	skillDir, err := agentSkillInstallDir(root, agentDoctorName, target, agentDoctorGlobal)
	if err != nil {
		return newExitCodeError(doctorExitError, err)
	}
	info, err := os.Stat(skillDir)
	if err != nil || !info.IsDir() {
		return newExitCodeError(doctorExitError, fmt.Errorf("skill not installed at %s; run 'markata-go agent install' first", skillDir))
	}

	manifest, err := agentskill.ReadManifest(skillDir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			outlnf("Skill at %s has no manifest.", skillDir)
			outlnf("Re-install with 'markata-go agent install --force' to enable drift detection.")
			return newExitCodeError(doctorExitDrift, nil)
		}
		return newExitCodeError(doctorExitError, fmt.Errorf("read manifest: %w", err))
	}

	report, err := agentskill.ComputeDrift(manifest, skillDir, Version)
	if err != nil {
		return newExitCodeError(doctorExitError, fmt.Errorf("compute drift: %w", err))
	}

	outlnf("Skill:     %s", agentDoctorName)
	outlnf("Agent:     %s", target.Name)
	outlnf("Scope:     %s", agentScopeName(agentDoctorGlobal))
	outlnf("Location:  %s", skillDir)
	outlnf("Installed: %s", manifest.Version)
	outlnf("Current:   %s", Version)
	if report.VersionMismatch {
		outlnf("Version:   %s", agentYellow("mismatch"))
	}
	outlnf("")

	for _, entry := range report.Files {
		switch entry.Status {
		case agentskill.FileOK:
			outlnf("  %s        %s", agentGreen("ok"), entry.Path)
		case agentskill.FileModified:
			outlnf("  %s  %s", agentYellow("modified"), entry.Path)
		case agentskill.FileNew:
			outlnf("  %s       %s", agentCyan("new"), entry.Path)
		case agentskill.FileMissing:
			outlnf("  %s   %s", agentRed("missing"), entry.Path)
		}
	}

	outlnf("")
	if report.HasDrift() {
		outlnf("%s Run 'markata-go agent update' to update.", agentYellow("Skill has %d issue(s).", report.IssueCount()))
		return newExitCodeError(doctorExitDrift, nil)
	}

	outlnf("%s", agentGreen("Skill is up to date."))
	return nil
}

func runAgentRemoveCommand(cmd *cobra.Command, args []string) error {
	currentCmd = cmd

	target, err := resolveAgentInstallTarget(agentRemoveAgent, agentRemoveLegacyTarget, agentRemoveGlobal)
	if err != nil {
		return err
	}

	root, err := resolveAgentProjectRoot(args, agentRemoveGlobal)
	if err != nil {
		return err
	}

	if err := validateSkillName(agentRemoveName); err != nil {
		return err
	}

	skillDir, err := agentSkillInstallDir(root, agentRemoveName, target, agentRemoveGlobal)
	if err != nil {
		return err
	}
	info, err := os.Stat(skillDir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("skill not installed at %s; run 'markata-go agent install' first", skillDir)
		}
		return fmt.Errorf("stat skill directory %q: %w", skillDir, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("skill path %q is not a directory", skillDir)
	}

	// Safety check: verify the directory looks like an installed skill before
	// removing it. Require SKILL.md or .manifest.json to be present.
	hasSkillMarker := false
	for _, marker := range []string{"SKILL.md", agentskill.ManifestFileName} {
		if _, statErr := os.Stat(filepath.Join(skillDir, marker)); statErr == nil {
			hasSkillMarker = true
			break
		}
	}
	if !hasSkillMarker {
		return fmt.Errorf("directory %q does not appear to be an installed skill (missing SKILL.md and %s)", skillDir, agentskill.ManifestFileName)
	}

	if err := os.RemoveAll(skillDir); err != nil {
		return fmt.Errorf("remove skill directory %q: %w", skillDir, err)
	}

	if !agentRemoveGlobal && root != "" {
		relPath, relErr := filepath.Rel(root, skillDir)
		if relErr == nil {
			skillDir = relPath
		}
	}
	outlnf("Removed %s", filepath.ToSlash(skillDir))
	return nil
}

// Agent status colors follow the shared CLI color controls for the output stream.
func agentColor(code, message string) string {
	enabled := currentCmd != nil && colorEnabledFor(streamIsTerminal(currentCmd.OutOrStdout()))
	return colorize(message, "\x1b["+code+"m", enabled)
}

func agentGreen(format string, args ...any) string {
	return agentColor("32", fmt.Sprintf(format, args...))
}
func agentYellow(format string, args ...any) string {
	return agentColor("33", fmt.Sprintf(format, args...))
}
func agentRed(format string, args ...any) string {
	return agentColor("31", fmt.Sprintf(format, args...))
}
func agentCyan(format string, args ...any) string {
	return agentColor("36", fmt.Sprintf(format, args...))
}
