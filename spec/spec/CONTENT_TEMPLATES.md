# Shared Markdown Content Starters Specification

This specification defines site-authored Markdown starters used to create
content from `markata-go new`, its interactive TUI, and the web workbench.
Each editor may provide a different interaction, but they share discovery,
metadata, defaults, and generated Markdown behavior.

## Goals

- Let a site define a content starter once, in Markdown.
- Make the starter available to CLI, TUI, and web workbench editors.
- Keep starter files readable and useful when opened as ordinary Markdown.
- Preserve existing CLI built-ins and configuration-based templates.
- Make editor-specific actions optional extensions, not required for a starter.

## Starter files

The default directory is `content-templates/`. A starter is one `.md` file at
the top level of that directory. The file basename is its canonical name
unless the optional `name` metadata field supplies a canonical name.
The configured `content_templates.directory` setting may change this location.
Nested starter discovery is not required by this specification.

```markdown
---
starter:
  name: guide
  description: A structured guide with an introduction, steps, and summary.
  aliases:
    - tutorial
    - how-to
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

The opening YAML block defines starter metadata and default content
frontmatter. The Markdown after the closing delimiter is the starter body.
The optional `starter` mapping may define `name`, `description`, and `aliases`
for discovery. Those fields and `_directory` MUST NOT appear in generated
content. A top-level `description` remains ordinary page frontmatter.
`_directory` is the default output directory.
When `_directory` is omitted, editors use the canonical starter name as the
output directory. When `starter.name` is omitted, the basename supplies it.
When `starter.description` is omitted, editors may use a concise generated
description. `starter.aliases` MUST be a list of strings.

## Shared creation behavior

All editors MUST:

1. Discover configured Markdown starters and merge them with built-in starters.
2. Resolve CLI/TUI selections by canonical name or alias. In the web workbench,
   show aliases as searchable terms and select by the visible canonical starter.
3. Ask for or accept a title when the starter requires one.
4. Generate a collision-safe slug-based Markdown path inside the selected site.
5. Expand the standard variables below in both starter frontmatter and body.
6. Omit starter metadata from generated frontmatter.
7. Preserve all other starter frontmatter values and body content.
8. Serialize only fields already present in the selected starter or values
   explicitly overridden during creation; do not inject implicit defaults.
9. Reject paths that escape the site root and report invalid starter metadata.

Site-authored starters override a config starter or built-in starter with
the same canonical name. Canonical names MUST remain selectable even when a
different starter uses that name as an alias. CLI/TUI MUST report an ambiguous
alias instead of silently choosing a starter. The web workbench may show multiple canonical results for an alias
search; it MUST NOT silently select one. Existing built-in alias catalogs may
differ between editor surfaces for compatibility.
Starter output directories MUST be relative paths that remain inside the site
root. A CLI user's explicit `--dir` override retains the CLI's existing
directory override behavior.

## Standard variables

| Variable | Meaning |
|---|---|
| `{{title}}` | User-provided title, YAML-escaped when emitted in YAML |
| `{{slug}}` | Lowercase, URL-safe slug derived from the title |
| `{{date}}` | Current local date as `YYYY-MM-DD` |
| `{{datetime}}` | Current time in UTC as RFC 3339 |

Unknown variables MUST remain literal so a starter can contain template syntax
for another processor. Standard variables are supported in YAML scalar values,
including scalar items in block or flow sequences, but not in YAML keys. Editors
MUST escape substitutions for the surrounding YAML scalar context; they MUST
NOT parse and re-serialize the body as YAML.

Editors may provide standard creation fields such as tags, draft/published,
private, and authors. The precedence order is:

1. Values explicitly supplied by the user (including `false`, cleared values,
   and empty lists).
2. Values present in a site-authored Markdown starter.
3. Site-configured template and site-level defaults.
4. Markata built-in defaults.

An editor MAY show inherited or computed defaults in its controls, but MUST NOT
serialize them into new content unless the user explicitly selects an override.
Unselected fields preserve the starter's values or remain absent. For example,
when a starter specifies only `title: "{{title}}"`, creating a page should
write that title but MUST NOT insert `slug`, `date`, `published`, `draft`,
`tags`, or `description` keys. The title supplied for a file path does not
require an absent `title` frontmatter field to be created. A missing date
resolves to the current local date at runtime; missing publication fields use
the site's effective defaults without being written into the new file.
Explicit CLI flags and TUI/web changes override starter values.

## Editor extensions

The Markdown starter is the portable baseline. Editors MAY offer optional
named capabilities such as URL metadata lookup or image upload. A starter may
request a capability through a validated metadata field, but a missing
capability MUST NOT prevent the starter from appearing or its Markdown from
being created. Arbitrary code execution, shell commands, and editor-specific
JavaScript embedded in a starter are out of scope.

## Compatibility and precedence

- Existing built-in starter names and CLI aliases remain available.
- Existing TOML `content_templates.templates` definitions remain supported.
- Markdown starters in the configured `content_templates.directory` take
  precedence over same-named config templates and built-ins.
- TOML/JSON/YAML configured templates take precedence over built-ins, and
  configured placements apply where no higher-precedence starter specifies
  its own output directory.
- Web workbench discovery MUST use the same configured/default directory and
  metadata contract as the CLI.
- The interactive TUI uses the same resolved catalog as non-interactive CLI.

## Errors

Editors MUST identify the starter file for YAML, alias, variable, or path
errors. Invalid starter metadata MUST be reported before a file is created.
Starters MUST NOT be interpreted as executable code.

## Conformance

An editor conforms when the same starter file, site root, and explicitly
selected creation inputs produce equivalent persisted frontmatter, substituted
values, body, and output directory. Editors MAY display different effective
fallbacks without writing them to the file. Differences in prompts and visual presentation are
allowed.
