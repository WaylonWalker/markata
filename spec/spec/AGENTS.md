# Agent Integrations Specification

This document specifies how markata-go distributes agent-facing project guidance for site repositories.

## Goals

- Provide a bundled, installable skill for agents working on markata-go sites.
- Keep the skill portable across agent tools.
- Split guidance into focused topics instead of a single monolithic prompt.
- Leave room in the CLI for future tool-specific exports and MCP-backed integrations.

## Bundled Skill

Markata-go ships a bundled skill named `markata-go-site`.

The bundled skill MUST include:

- `SKILL.md` as the portable entrypoint
- a `topics/` directory with focused guidance files
- a `reference/` directory with quick-lookup reference material
- an `examples/` directory with starter config and template files
- an `evals/` directory with starter regression prompts for the bundled skill itself

The initial required topic files are:

- `configuration.md`
- `writing-frontmatter.md`
- `cli-usage.md`
- `build-deployment.md`
- `faster-builds.md`
- `theme-creation.md`
- `template-management.md`
- `plugin-creation.md`
- `roads.md`

The initial reference files are:

- `template-context.md`
- `feed-patterns.md`
- `palette-reference.md`

The initial example files are:

- `fast.toml`
- `markata-go.local.toml`
- `palettes/my-brand.toml`
- `templates/base.html`
- `templates/post.html`
- `templates/feed.html`

The initial eval files are:

- `evals/evals.json`

The entrypoint SHOULD tell agents to read only the topic files relevant to the current task. Reference and example files SHOULD be used when agents need exact shapes or starter material rather than narrative guidance.

The `SKILL.md` frontmatter description MUST name concrete site-repository tasks and decision contexts so the bundled skill triggers reliably for real work instead of only broad repository descriptions.

## Skill Content Requirements

The bundled skill MUST guide agents toward:

- inspecting the active config before editing behavior
- preferring content, frontmatter, config, templates, and CSS before custom plugin work
- using built-in CLI inspection commands where possible
- preserving the site's existing layout and conventions unless the task explicitly changes them
- using warm-build comparisons for performance work

The bundled skill MUST also include a starter eval set that covers at least:

- config changes or debugging
- frontmatter or content creation
- template or layout edits
- build or deploy debugging
- deciding whether work belongs in config, templates, feeds, CSS, or plugins

## Skill Maintenance Requirement

The bundled site skill is part of markata-go's user-facing product surface for coding agents.

Any change that affects how an agent should work in a markata-go site repository MUST review the bundled skill and update it when needed.

Common maintenance triggers include changes to:

- configuration behavior or recommended config patterns
- frontmatter semantics or content creation workflows
- CLI usage, flags, or command output that agents rely on
- build and deployment workflows
- performance guidance such as `--fast`, benchmarking, or cache behavior
- theme, palette, and aesthetic workflows
- template, layout, and template-context behavior
- plugin authoring constraints or extension workflows

When such a change does not require modifying the bundled skill, the implementation or PR SHOULD explicitly state why no skill change was needed.

## Installation Targets

The CLI MUST support agent-specific directory layouts using the same agent identifiers documented by `vercel-labs/skills`.

Required agent identifiers:

- `adal`
- `amp`
- `antigravity`
- `augment`
- `bob`
- `claude-code`
- `cline`
- `codebuddy`
- `codex`
- `command-code`
- `continue`
- `cortex`
- `crush`
- `cursor`
- `deepagents`
- `droid`
- `firebender`
- `gemini-cli`
- `github-copilot`
- `goose`
- `iflow-cli`
- `junie`
- `kimi-cli`
- `kilo`
- `kiro-cli`
- `kode`
- `mcpjam`
- `mistral-vibe`
- `mux`
- `neovate`
- `openclaw`
- `opencode`
- `openhands`
- `pi`
- `pochi`
- `qoder`
- `qwen-code`
- `replit`
- `roo`
- `trae`
- `trae-cn`
- `universal`
- `warp`
- `windsurf`
- `zencoder`

The CLI MUST resolve each agent to the same project and global skill directories used by `vercel-labs/skills`.

Legacy compatibility aliases MAY be accepted for existing users:

- `agents` -> `universal`
- `claude` -> `claude-code`

The installed file contents MUST remain the same across agents unless a future target explicitly requires generated wrappers.

## CLI

Markata-go MUST expose an `agent` command group.

Initial command:

```bash
markata-go agent install [site-path]
```

Additional required subcommand:

```bash
markata-go agent list-agents
```

### `list-agents` behavior

- The command MUST be read-only.
- The command MUST write primary results to `stdout`.
- The command MUST list each supported agent identifier.
- The command MUST include the project and global skill directories for each agent.
- When compatibility aliases exist, the command SHOULD show them.

Required flags:

- `--agent`
- `--name`
- `--force`
- `--dry-run`

Optional flags:

- `-g`, `--global`

### `install` behavior

- `site-path` defaults to the current directory.
- When `--agent` is omitted for project installs, the command MUST default to the current agent when it can detect one from the environment. If no current agent can be detected, it MUST default to `universal`.
- When `--global` is set without `--agent` or the legacy `--target`, `install` MUST open a fuzzy agent picker on interactive input and output terminals.
- The picker MUST list supported agent identifiers and global directories. Typing MUST filter names by ordered character matches.
- The picker MUST use fzf's default bottom-up layout: the prompt at the bottom, the first result nearest the prompt, and additional results above it. Up MUST move toward results above the current selection; Down MUST move toward the prompt.
- The picker MUST color the selection, matched characters, and prompt with the active terminal theme. It MUST respect the CLI's color controls, including `--no-color` and `NO_COLOR`.
- Enter MUST select an agent. Escape or Ctrl+C MUST cancel without writing files.
- With `--no-input`, `TERM=dumb`, or non-terminal input or output, `--global` MUST require an explicit agent and show an example plus `agent list-agents` guidance.
- Unknown agent names MUST remain errors and suggest close identifiers, including `github-copilot` for `copilot`. Suggestions MUST include short names, case-insensitive matches, and minor spelling errors without selecting an agent automatically.
- `--global` MUST install into the selected agent's user-level skill directory instead of a repository path.
- `site-path` MUST be rejected when `--global` is set.
- The command MUST install bundled files into the selected agent layout.
- The command MUST fail clearly when destination files already exist and `--force` is not set.
- `--dry-run` MUST report what would be written without modifying the filesystem.
- Primary results MUST be written to `stdout`.
- Errors MUST suggest `--force` when overwrite conflicts are the only blocker.

---

```bash
markata-go agent update [site-path]
```

Required flags:

- `--agent`
- `--name`
- `--dry-run`

Optional flags:

- `-g`, `--global`

### `update` behavior

- `site-path` defaults to the current directory.
- The command MUST update the installed skill in place using the currently bundled files.
- `update` MUST behave like `install --force` with the same agent and scope resolution.
- `--global` MUST require an explicit `--agent`.
- `site-path` MUST be rejected when `--global` is set.
- `--dry-run` MUST report what would be updated without modifying the filesystem.
- Primary results MUST be written to `stdout`.

---

```bash
markata-go agent remove [site-path]
markata-go agent uninstall [site-path]
```

Required flags:

- `--agent`
- `--name`

Optional flags:

- `-g`, `--global`

### `remove` behavior

- `site-path` defaults to the current directory.
- The command MUST remove the installed skill directory for the selected agent, scope, and name.
- `--global` MUST require an explicit `--agent`.
- `site-path` MUST be rejected when `--global` is set.
- `uninstall` MUST be a direct alias for `remove`.
- The command MUST fail clearly when the skill is not installed at the resolved location.
- Primary results MUST be written to `stdout`.

### Install Manifest

After writing skill files, `install` MUST write a manifest file named `.manifest.json` inside the installed skill directory.

The manifest MUST contain:

- `version` — the markata-go binary version string (from `cmd.Version`)
- `installed_at` — RFC 3339 timestamp of installation
- `target` — the selected agent identifier (for example `opencode` or `claude-code`)
- `scope` — `project` or `global`
- `files` — a map of relative file paths to their SHA-256 hex digests

The manifest MUST NOT be included in the overwrite-protection check. It is always overwritten on install.

The manifest MUST NOT appear in the `--dry-run` output file list or the installed file count reported to the user. It is an internal bookkeeping file.

Example manifest:

```json
{
  "version": "0.5.0",
  "installed_at": "2026-04-01T12:00:00Z",
  "target": "opencode",
  "scope": "project",
  "files": {
    "SKILL.md": "a1b2c3...",
    "topics/configuration.md": "d4e5f6..."
  }
}
```

---

```bash
markata-go agent doctor [site-path]
```

Required flags:

- `--agent`
- `--name`

Optional flags:

- `-g`, `--global`

### `doctor` behavior

The `doctor` command detects drift between installed skill files and the bundled versions in the current binary.

- `site-path` defaults to the current directory.
- The command MUST locate the installed skill directory using the same agent/scope/name resolution as `install`.
- `--global` MUST require an explicit `--agent`.
- `site-path` MUST be rejected when `--global` is set.
- If no manifest file exists, the command MUST report that the skill was installed without a manifest and recommend re-installing with the current binary.
- If the manifest exists, the command MUST compare:
  1. The `version` field against the current binary version.
  2. Each file in the manifest `files` map against the SHA-256 digest of the corresponding bundled file.
  3. Whether new files exist in the bundle that are not present in the manifest.

### `doctor` output

The command MUST report a per-file status using these categories:

| Category | Meaning |
|----------|---------|
| `ok` | File hash matches the bundled version |
| `modified` | File exists on disk but hash differs from the bundled version |
| `new` | File exists in the current bundle but was not present when the skill was installed |
| `missing` | File is listed in the manifest but does not exist on disk |

The command MUST print a summary line:

- If all files are `ok` and no `new` files exist: `Skill is up to date.`
- Otherwise: `Skill has N issue(s). Run 'markata-go agent install --force' to update.`

### `doctor` exit codes

| Code | Meaning |
|------|---------|
| `0` | Skill is up to date (all files `ok`, no `new` files) |
| `1` | Drift detected (any `modified`, `new`, or `missing` files) |
| `2` | Error (skill not installed, manifest unreadable, etc.) |

### `doctor` dry-run and machine-readable output

- `doctor` does not modify the filesystem. It is read-only.
- Future iterations MAY add `--json` for machine-readable output.

## Future Extensibility

The command group is reserved for future additions such as:

- `agent export`
- `agent inspect`
- `agent mcp ...`

Future subcommands MUST be able to reuse the same bundled skill source without requiring a rewrite of the portable `SKILL.md`, `topics/`, `reference/`, or `examples/` content.
