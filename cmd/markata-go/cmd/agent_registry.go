package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

const (
	agentDefaultUniversal = "universal"
	agentLegacyAgents     = "agents"
	agentLegacyClaude     = "claude"

	agentScopeProject = "project"
	agentScopeGlobal  = "global"

	agentNameOpencode = "opencode"
	agentNameCodex    = "codex"
	envValueDisabled  = "false"
)

type agentInstallTarget struct {
	Name        string
	ProjectPath string
	GlobalPath  string
	Aliases     []string
}

var supportedAgentTargets = []agentInstallTarget{
	{Name: "adal", ProjectPath: ".adal/skills", GlobalPath: ".adal/skills"},
	{Name: "amp", ProjectPath: ".agents/skills", GlobalPath: ".config/agents/skills"},
	{Name: "antigravity", ProjectPath: ".agents/skills", GlobalPath: ".gemini/antigravity/skills"},
	{Name: "augment", ProjectPath: ".augment/skills", GlobalPath: ".augment/skills"},
	{Name: "bob", ProjectPath: ".bob/skills", GlobalPath: ".bob/skills"},
	{Name: "claude-code", ProjectPath: ".claude/skills", GlobalPath: ".claude/skills", Aliases: []string{agentLegacyClaude}},
	{Name: "cline", ProjectPath: ".agents/skills", GlobalPath: ".agents/skills"},
	{Name: "codebuddy", ProjectPath: ".codebuddy/skills", GlobalPath: ".codebuddy/skills"},
	{Name: agentNameCodex, ProjectPath: ".agents/skills", GlobalPath: ".codex/skills"},
	{Name: "command-code", ProjectPath: ".commandcode/skills", GlobalPath: ".commandcode/skills"},
	{Name: "continue", ProjectPath: ".continue/skills", GlobalPath: ".continue/skills"},
	{Name: "cortex", ProjectPath: ".cortex/skills", GlobalPath: ".snowflake/cortex/skills"},
	{Name: "crush", ProjectPath: ".crush/skills", GlobalPath: ".config/crush/skills"},
	{Name: "cursor", ProjectPath: ".agents/skills", GlobalPath: ".cursor/skills"},
	{Name: "deepagents", ProjectPath: ".agents/skills", GlobalPath: ".deepagents/agent/skills"},
	{Name: "droid", ProjectPath: ".factory/skills", GlobalPath: ".factory/skills"},
	{Name: "firebender", ProjectPath: ".agents/skills", GlobalPath: ".firebender/skills"},
	{Name: "gemini-cli", ProjectPath: ".agents/skills", GlobalPath: ".gemini/skills"},
	{Name: "github-copilot", ProjectPath: ".agents/skills", GlobalPath: ".copilot/skills"},
	{Name: "goose", ProjectPath: ".goose/skills", GlobalPath: ".config/goose/skills"},
	{Name: "iflow-cli", ProjectPath: ".iflow/skills", GlobalPath: ".iflow/skills"},
	{Name: "junie", ProjectPath: ".junie/skills", GlobalPath: ".junie/skills"},
	{Name: "kimi-cli", ProjectPath: ".agents/skills", GlobalPath: ".config/agents/skills"},
	{Name: "kilo", ProjectPath: ".kilocode/skills", GlobalPath: ".kilocode/skills"},
	{Name: "kiro-cli", ProjectPath: ".kiro/skills", GlobalPath: ".kiro/skills"},
	{Name: "kode", ProjectPath: ".kode/skills", GlobalPath: ".kode/skills"},
	{Name: "mcpjam", ProjectPath: ".mcpjam/skills", GlobalPath: ".mcpjam/skills"},
	{Name: "mistral-vibe", ProjectPath: ".vibe/skills", GlobalPath: ".vibe/skills"},
	{Name: "mux", ProjectPath: ".mux/skills", GlobalPath: ".mux/skills"},
	{Name: "neovate", ProjectPath: ".neovate/skills", GlobalPath: ".neovate/skills"},
	{Name: "openclaw", ProjectPath: "skills", GlobalPath: ".openclaw/skills"},
	{Name: "opencode", ProjectPath: ".agents/skills", GlobalPath: ".config/opencode/skills"},
	{Name: "openhands", ProjectPath: ".openhands/skills", GlobalPath: ".openhands/skills"},
	{Name: "pi", ProjectPath: ".pi/skills", GlobalPath: ".pi/agent/skills"},
	{Name: "pochi", ProjectPath: ".pochi/skills", GlobalPath: ".pochi/skills"},
	{Name: "qoder", ProjectPath: ".qoder/skills", GlobalPath: ".qoder/skills"},
	{Name: "qwen-code", ProjectPath: ".qwen/skills", GlobalPath: ".qwen/skills"},
	{Name: "replit", ProjectPath: ".agents/skills", GlobalPath: ".config/agents/skills"},
	{Name: "roo", ProjectPath: ".roo/skills", GlobalPath: ".roo/skills"},
	{Name: "trae", ProjectPath: ".trae/skills", GlobalPath: ".trae/skills"},
	{Name: "trae-cn", ProjectPath: ".trae/skills", GlobalPath: ".trae-cn/skills"},
	{Name: "universal", ProjectPath: ".agents/skills", GlobalPath: ".config/agents/skills", Aliases: []string{agentLegacyAgents}},
	{Name: "warp", ProjectPath: ".agents/skills", GlobalPath: ".agents/skills"},
	{Name: "windsurf", ProjectPath: ".windsurf/skills", GlobalPath: ".codeium/windsurf/skills"},
	{Name: "zencoder", ProjectPath: ".zencoder/skills", GlobalPath: ".zencoder/skills"},
}

var supportedAgentTargetsByName = newSupportedAgentTargetsByName()

func newSupportedAgentTargetsByName() map[string]agentInstallTarget {
	lookup := make(map[string]agentInstallTarget, len(supportedAgentTargets)*2)
	for _, target := range supportedAgentTargets {
		lookup[target.Name] = target
		for _, alias := range target.Aliases {
			lookup[alias] = target
		}
	}
	return lookup
}

func normalizeAgentInstallTarget(name string) (agentInstallTarget, error) {
	canonical := strings.ToLower(strings.TrimSpace(name))
	target, ok := supportedAgentTargetsByName[canonical]
	if !ok {
		message := fmt.Sprintf("unknown agent %q.", name)
		for _, suggestion := range suggestAgentNames(canonical) {
			message += fmt.Sprintf("\nDid you mean --agent %s?", suggestion)
		}
		return agentInstallTarget{}, fmt.Errorf("%s\nRun 'markata-go agent list-agents' for supported values", message)
	}
	return target, nil
}

func suggestAgentNames(name string) []string {
	command := &cobra.Command{Use: "agent", SuggestionsMinimumDistance: 2}
	for _, target := range supportedAgentTargets {
		for _, candidate := range append([]string{target.Name}, target.Aliases...) {
			child := &cobra.Command{
				Use: candidate,
				Run: func(_ *cobra.Command, _ []string) {},
			}
			if len(name) >= 3 && strings.Contains(candidate, name) {
				child.SuggestFor = []string{name}
			}
			command.AddCommand(child)
		}
	}
	seen := make(map[string]bool)
	var suggestions []string
	for _, suggestion := range command.SuggestionsFor(name) {
		canonical := supportedAgentTargetsByName[suggestion].Name
		if !seen[canonical] {
			suggestions = append(suggestions, canonical)
			seen[canonical] = true
		}
	}
	sort.Strings(suggestions)
	return suggestions
}

func resolveAgentInstallTarget(agentFlag, legacyTarget string, global bool) (agentInstallTarget, error) {
	explicit, explicitSet, err := selectedAgentFlagValue(agentFlag, legacyTarget)
	if err != nil {
		return agentInstallTarget{}, err
	}
	if global && !explicitSet {
		return agentInstallTarget{}, fmt.Errorf("--global requires --agent without an interactive picker; use '--agent github-copilot' or run 'markata-go agent list-agents' for supported values")
	}
	if explicitSet {
		return normalizeAgentInstallTarget(explicit)
	}
	return normalizeAgentInstallTarget(detectDefaultAgent())
}

func selectedAgentFlagValue(agentFlag, legacyTarget string) (selected string, explicit bool, err error) {
	agentFlag = strings.TrimSpace(agentFlag)
	legacyTarget = strings.TrimSpace(legacyTarget)
	if agentFlag == "" && legacyTarget == "" {
		return "", false, nil
	}
	if agentFlag == "" {
		return legacyTarget, true, nil
	}
	if legacyTarget == "" {
		return agentFlag, true, nil
	}

	normalizedAgent, agentErr := normalizeAgentInstallTarget(agentFlag)
	if agentErr != nil {
		return "", false, agentErr
	}
	normalizedLegacy, legacyErr := normalizeAgentInstallTarget(legacyTarget)
	if legacyErr != nil {
		return "", false, legacyErr
	}
	if normalizedAgent.Name != normalizedLegacy.Name {
		return "", false, fmt.Errorf("--agent %q conflicts with legacy --target %q", agentFlag, legacyTarget)
	}
	return agentFlag, true, nil
}

func detectDefaultAgent() string {
	switch {
	case envVarEnabled("OPENCODE"):
		return agentNameOpencode
	case envVarEnabled("CLAUDECODE"), envVarEnabled("CLAUDE_CODE"):
		return "claude-code"
	case envVarEnabled("CODEX"), envVarEnabled("OPENAI_CODEX"):
		return agentNameCodex
	case envVarEnabled("CURSOR_AGENT"), os.Getenv("CURSOR_TRACE_ID") != "":
		return "cursor"
	case envVarEnabled("GEMINI_CLI"):
		return "gemini-cli"
	case envVarEnabled("QWEN_CODE"):
		return "qwen-code"
	case envVarEnabled("KIMI_CLI"):
		return "kimi-cli"
	default:
		return agentDefaultUniversal
	}
}

func envVarEnabled(name string) bool {
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		return false
	}
	value = strings.ToLower(value)
	return value != "0" && value != envValueDisabled && value != "no"
}

func agentSkillInstallDir(root, skillName string, target agentInstallTarget, global bool) (string, error) {
	base := root
	relPath := target.ProjectPath
	if global {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("resolve user home directory: %w", err)
		}
		base = home
		relPath = target.GlobalPath
	}
	return filepath.Join(base, filepath.FromSlash(relPath), skillName), nil
}

func agentScopeName(global bool) string {
	if global {
		return agentScopeGlobal
	}
	return agentScopeProject
}

func supportedAgentNames() []string {
	names := make([]string, 0, len(supportedAgentTargets))
	for _, target := range supportedAgentTargets {
		names = append(names, target.Name)
	}
	sort.Strings(names)
	return names
}
