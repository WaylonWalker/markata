---
title: "Give your agent the Markata skill"
description: "Teach OpenCode, Codex, Claude Code, Cursor, and other coding agents to build a Markata site."
date: 2026-10-18
published: true
tags: [documentation, agents]
---


Let's be honest: a lot of websites get built alongside coding agents. **Give the agent the instructions that ship with Markata**, rather than a stale prompt found in a chat.

The skill is bundled with the **Go binary**; you don't need a separate installer or a proprietary cloud service.

## 1. Open your site repository

```bash
cd my-site
markata agent list-agents
```

This prints the agent identifiers and their project/global skill locations. Supported harnesses include **OpenCode, Codex, Claude Code, Cursor, GitHub Copilot, Gemini CLI**, and more.

## 2. Install the skill

Use auto-detection (fallback: portable `.agents/skills/markata-go-site/`):

```bash
markata agent install
```

Or choose your harness explicitly:

```bash
markata agent install --agent opencode
markata agent install --agent claude-code
markata agent install --agent codex
```

For a specific agent, use `--agent`. For Claude Code, the project skill goes in `.claude/skills/markata-go-site/`; for a universal project install, it goes in `.agents/skills/markata-go-site/`.

The `setup` alias also installs the skill. On an interactive terminal, either
command lists detected agents and asks whether to install for all of them. Press
Enter or `y` to accept. Enter `n` to choose an agent from the fuzzy picker.
If no agents are detected, the picker opens directly. Agents that share a skill
directory use one installation.

For a global installation, use:

```bash
markata agent install -g
```

The picker uses an fzf-style layout: the prompt stays at the bottom, with results above it.
The selected agent and matching characters use the terminal theme colors.

Type an agent name or abbreviation, such as `copilot` or `ghcp`. Use the arrow keys to navigate. Press Enter to select.
Escape or Ctrl+C cancels without an installation.
`--no-color` or `NO_COLOR` disables the colors.
If an install for multiple agents fails, the error identifies the failed agent.
Installations completed before the error remain in place.

For scripts, project installs keep the environment-detected agent or universal fallback. Global installs need an explicit agent:

```bash
markata agent install --agent opencode --global
markata agent install --agent github-copilot --global
```

The picker requires interactive input and output terminals. With `--no-input` or `TERM=dumb`, an explicit `--agent` is required.
Unknown names show suggestions, such as `Did you mean --agent github-copilot?` for `--agent copilot`.
Suggestions do not select an agent automatically.

## 3. Check it

```bash
markata agent doctor
```

On a Markata upgrade, bring your skill up to date:

```bash
markata agent update
markata agent doctor
```

Use `--dry-run` on `install` or `update` to preview paths, or `markata agent install --force` to overwrite local customizations.

## 4. Start with a small instruction

> Read the installed Markata site skill and the current configuration. Create a draft Markdown post about what this site is for. Do not change the theme, plugins, or deployment setup. Show me the file and run a fast build to check it.

Agents can inspect the site with:

```bash
markata config show
markata list posts
markata explain content
markata build --fast
```

For rich content with symbols like :lucide-rocket: or :lucide-database:,
responsive icon cards, or status/reporting blocks, the bundled skill includes
`topics/rich-markdown.md`. Ask the agent to read that topic before creating
new CSS, templates, or JS. The built-in icon syntax is documented under
[Icon Shortcodes](/docs/guides/icons/).

Keep a human in the loop for publishing and deployment. Markata's skill guides agents through frontmatter, themes, feed patterns, templates, site builds, and when **not** to create a plugin.

## Visualize a session or topic

After you install the skill, ask your agent:

> Visualize this session.

The agent creates a visual Markdown post from the conversation's work, decisions, and unresolved items.
It serves a local preview, checks the exact post URL, and attempts to open that URL in your browser.
The handoff includes the Markdown path, preview URL, browser-open result, and server handle or stop instruction.

For supplied material, use either phrase:

> Visualize this: our release process starts with a draft, then review, then approval.

> Make a visual post from these notes.

These are natural-language instructions for the agent, not Markata CLI commands.
Ordinary chart, diagram, or styling requests do not automatically start a preview server.

For Markdown without a preview, state that constraint:

> Visualize this session. Source only: create exactly one Markdown file.

The agent then creates the file without a server or preview scaffolding.
Outside a Markata site, a full preview uses an isolated temporary site.
The agent keeps that scaffolding outside your repository.

Local preview does not authorize public publication or deployment.
The agent preserves publication flags and excludes secrets.
The agent reports browser or persistent-server limitations instead of claiming that the page opened.
On a remote machine, the preview can require port forwarding to reach your browser.

The workflow lives in the bundled skill's `topics/visualize-session.md`.
For an existing installation, use `markata agent update` after you install a binary with the updated skill.

## Fixes

- If the wrong harness is detected, pass `--agent` explicitly.
- If no skill is found, run `markata agent install --dry-run` and check the reported destination.
- After upgrading the binary, run `markata agent doctor` and then `markata agent update`.
- Keep the installed skill with your site's source control when you want a reproducible project setup.

Next: [Write your first post](/docs/guides/authoring/) or [Quickstart](/docs/quickstart/).
