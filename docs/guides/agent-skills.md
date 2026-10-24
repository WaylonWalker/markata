---
title: "Give your agent the Markata skill"
description: "Teach OpenCode, Codex, Claude Code, Cursor, and other coding agents to build a Markata site."
date: 2026-10-18
published: true
tags: [documentation, agents]
---

# Give your agent the Markata skill

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

Run only the command for the agent you use. For Claude Code, the project skill goes in `.claude/skills/markata-go-site/`; for a universal project install, it goes in `.agents/skills/markata-go-site/`.

If you work across many site repositories, install globally with an explicit agent:

```bash
markata agent install --agent opencode --global
```

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

## Fixes

- If the wrong harness is detected, pass `--agent` explicitly.
- If no skill is found, run `markata agent install --dry-run` and check the reported destination.
- After upgrading the binary, run `markata agent doctor` and then `markata agent update`.
- Keep the installed skill with your site's source control when you want a reproducible project setup.

Next: [Write your first post](/docs/guides/authoring/) or [Quickstart](/docs/quickstart/).
