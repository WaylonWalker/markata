---
title: "Triage a Markata problem with a coding agent"
description: "Collect safe site diagnostics and start an interactive troubleshooting conversation, copy a prompt, or pipe it to an agent"
date: 2026-10-19
published: true
tags:
  - documentation
  - troubleshooting
  - agents
  - cli
---

# Troubleshooting with `markata triage`

When your Markata site behaves unexpectedly, start with:

```sh
markata triage
```

Markata collects lightweight local diagnostics, then displays a **fuzzy picker** with the coding agents it finds on your PATH. It also offers **Copy to clipboard** and **Print to stdout**. Pick a coding agent and the agent opens with a troubleshooting prompt: it should **ask you what happened first**, instead of making guesses or immediately modifying your site.

The command works with the new `markata` binary and the backwards-compatible `markata-go` binary.

## Task-specific intake before agent handoff

Triage can help with lint errors, broken Markdown features, or turning a written
post into a richer presentation. In an interactive terminal, **missing intake
flags become deterministic fuzzy-search choices** before the agent picker appears.
Category, scope, action, and feature are searchable menus with descriptions.
When a post is needed, triage offers a fuzzy picker of local Markdown filenames
from common content directories, plus **Enter a Markdown path manually**.
Issue descriptions and enhancement goals remain free-text questions.
Answers already supplied as flags are never asked again.

```sh
# Lint the whole site, or just a specific post
markata triage lint --scope site
markata triage lint --scope post --post posts/example.md

# Fix a chart or Mermaid diagram that does not render
markata triage post --post posts/example.md --action debug \
  --feature chart --issue "The axis labels overlap"
markata triage post --post posts/example.md --action debug \
  --feature mermaid --issue "The diagram fails to parse"

# Enhance existing writing with a road or other complex component
markata triage post --post posts/history.md --action enhance \
  --feature road --goal "Show the project's milestones as a static road"

# Ask the CLI for the category, post, intent and feature first
markata triage
markata triage post
```

Available categories: `general`, `lint`, `post`. Lint asks for `--scope site|post`
and a `--post` path when needed. Post tasks ask for a `--post` path, `--action
debug|enhance|review`, and a `--feature chart|road|mermaid|markdown|other`
for debugging or enhancement. Debugging asks for an `--issue`, while enhancement
asks for a `--goal`. Only filenames are scanned for the post selector;
post contents and configuration contents are not read by triage. A selected path
is provided to the agent. Files in custom content directories can always be
entered manually. An unavailable site directory also falls back to manual entry.

With `--no-input` or redirected streams, no interactive question is
presented. Missing intake answers are listed **in a fixed order inside the clean
Markdown prompt**, so it remains safe to pipe or save. Existing generic
`markata triage --issue ...` commands continue to work.
Questions and destination choices are separate: triage collects what you want
first and then lets you explicitly select the coding harness. Press **Escape**
in any intake picker to cancel triage without generating a prompt or launching
an agent. Supplied `--post` flags skip the filename scan entirely.

## Where the prompt goes

| Option | Behavior |
| --- | --- |
| `markata triage` | Fuzzy picker in an interactive terminal; prompt-only stdout when piped or redirected |
| `markata triage --to opencode` | Start an interactive OpenCode session with an initial prompt |
| `markata triage --to claude` | Start an interactive Claude Code session with an initial prompt |
| `markata triage --to codex` | Start interactive Codex with an initial prompt |
| `markata triage --to gemini` | Start interactive Gemini CLI |
| `markata triage --to cursor` | Start interactive Cursor Agent |
| `markata triage --to clipboard` | Copy the prompt using the platform's clipboard helper |
| `markata triage --to stdout` | Print only the Markdown prompt (ideal for pipes) |
| `markata triage --agent codex` | Alias for selecting a specific installed agent |
| `markata triage --issue "Images are missing"` | Include a short first-hand problem description |

Harness detection only checks that an executable is on PATH. It does not guarantee the agent is authenticated or configured with a model. Unsupported or missing agent names produce an error; clipboard and stdout still work.

For Linux clipboard handoffs, install `wl-copy`, `xclip`, or `xsel`. macOS uses `pbcopy`; Windows uses `clip.exe`. A missing helper is reported as an error; use stdout instead.

## Piping to an agent

`--to stdout` is clean Markdown with no banner, terminal formatting, or interactive picker. It can be redirected to a file or piped directly into another CLI:

```sh
# Interactive picker disabled automatically in a pipe
markata triage --to stdout | codex exec -

# Use OpenCode's noninteractive prompt runner
markata triage --to stdout | opencode run

# Use Claude Code's print mode
markata triage --to stdout | claude -p

# Save and inspect the context before sharing it
markata triage --to stdout > triage-prompt.md
```

Piped agent commands may be **noninteractive**. If you want the agent to ask questions back-and-forth in a terminal, use `--to opencode`, `--to claude`, etc., or select an agent from the fuzzy picker.

With `--no-input`, or if stdin or stdout is not a terminal, the default `--to auto` prints to stdout; it never unexpectedly launches an agent. Explicit interactive agent launch requires a terminal. Press Escape to cancel the fuzzy picker without launching anything. Closing a pipe early (for example, `markata triage | head`) exits quietly. The primary prompt remains plain and pipeable with `--no-color`, `NO_COLOR`, and `--quiet`.

Triage remains available even when a requested `--site-dir` or `MARKATA_GO_SITE_DIR` directory is missing or inaccessible; it records that problem in its diagnostics instead of failing before collecting context. Config file contents are never parsed as part of triage.

## What is collected

The prompt describes:
- The **running triage CLI's** embedded version, commit, and build date
- Operating system, architecture, Go runtime, and (on Linux) OS distribution
- Location of the Markata executable and an **inferred**, not guaranteed, installation method
- Working/site directory and names of local Markata config files (or explicit config path)
- Presence of a small allowlist of environment variables, **not their values**
- Links to Markata's source, documentation, and issue tracker

It does **not** read your configuration contents, logs, passwords, authentication tokens, arbitrary environment values, or other user files. Home-directory prefixes in paths are abbreviated when possible. Because paths may still be personal, **inspect the prompt before sharing it publicly or sending it to a new provider**.

A key distinction: the Markata version used to generate the prompt may be different from the version that failed in production, Docker, CI, or a remote builder. The agent is instructed to ask which version and execution environment actually experienced the problem.

## What the agent should do

The prompt asks the agent to first understand your symptoms, intended behavior, reproduction steps, and exact errors. It should check whether the issue is with site content, configuration, templates, a plugin, installation, or Markata itself; investigate read-only first; and ask permission before changing files, running mutating commands, or publishing an upstream issue.

If you already have a clear problem statement, add `--issue` and the agent should not make you repeat it.

For common manual troubleshooting, see the [Troubleshooting guide](../troubleshooting.md) and the [CLI reference](../reference/cli.md). The [source repository](https://github.com/WaylonWalker/markata) and [issue tracker](https://github.com/WaylonWalker/markata/issues) are available when a fix belongs upstream.
