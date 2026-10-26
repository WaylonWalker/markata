---
title: "Quickstart"
description: "Build your first Markata site, install the agent skill, and write a post."
date: 2026-10-18
published: true
tags: [documentation, getting-started]
---

# Quickstart

Already installed? This is the shortest path from an empty directory to a local website.

If not, start with [Installation](/docs/installation/).

## 1. Create your site

```bash
mkdir my-site && cd my-site
markata config init
markata new "Hello World"
```

`config init` creates **`markata-go.toml`**. That's the correct filename for the Go project today, despite the new command name.

## 2. Preview it

```bash
markata serve
```

Open [http://localhost:8000](http://localhost:8000). The development server watches for changes and rebuilds your site as you edit.

## 3. Teach your agent

```bash
markata agent list-agents
markata agent install
markata agent doctor
```

The default project install detects a supported agent or uses `.agents/skills/markata-go-site/`. For a specific harness:

```bash
markata agent install --agent opencode
# or: markata agent install --agent claude-code
```

See the [agent guide](/docs/guides/agent-skills/) for Codex, Cursor, Claude Code, OpenCode, and global installs.

## 4. Write

Open `pages/post/hello-world.md`. Edit the Markdown below its YAML frontmatter and save. Your browser should update.

Prefer a draft while working? Run:

```bash
markata new "My draft" --draft
```

See [Authoring](/docs/guides/authoring/) for a complete walk-through.

## 5. Build for production

```bash
markata build
```

The output is in `output/` by default. Deploy that **generated** folder to any static host. Start with the [deployment guide](/docs/guides/deployment/).

## Need help?

- `markata --help`
- `markata config validate`
- `markata agent doctor`
- [Troubleshooting](/docs/troubleshooting/)
