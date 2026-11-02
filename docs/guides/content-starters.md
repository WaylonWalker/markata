---
title: "Shared Markdown content starters"
description: "Create reusable Markdown starters for the CLI, TUI, and web workbench."
date: 2026-11-01
published: true
tags:
  - content
  - templates
  - editors
---

Use a Markdown starter to define how a new page begins. The same starter is
available to `markata-go new`, its interactive TUI, and the web workbench.

Create `content-templates/guide.md`:

```markdown
---
starter:
  name: guide
  description: A structured guide with an introduction, steps, and summary.
  aliases:
    - tutorial
_directory: guides
templateKey: guide
title: "{{title}}"
date: "{{date}}"
published: false
tags: []
---

## Introduction

## Steps

## Summary
```

The file name supplies the starter name when `name` is omitted. `_directory`
sets its default output directory. `description` and `aliases` help users find
it in the editor. These metadata keys are not copied to the new page.

The shared variables are `{{title}}`, `{{slug}}`, `{{date}}`, and
`{{datetime}}`. They work in frontmatter and body text. Unknown variables stay
as written, which lets another renderer use its own syntax.

Run `markata-go new --list` to see starters, then create content with:

```sh
markata-go new "My Guide" --template guide
```

The web workbench discovers starters from the same configured
`content_templates.directory` (default: `content-templates`). The CLI and
interactive TUI use that same directory. Existing TOML-defined templates and
built-ins remain available.

See the [CLI reference](../reference/cli.md#template-system) and
[starter specification](../../spec/spec/CONTENT_TEMPLATES.md) for the full
format and precedence rules.
