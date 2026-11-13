# CLI `new` Command Specification

> **Shared starter format:** CLI, interactive CLI (TUI), and web workbench
> content creation MUST use the same site-authored Markdown starter files.
> The editor interface may differ, but template discovery and generated
> content MUST follow the contract in [Shared Markdown Starters](CONTENT_TEMPLATES.md).

This document specifies the `markata-go new` command behavior, including the
interactive TUI wizard powered by charmbracelet/huh.

## Overview

The `new` command creates new content files with frontmatter templates. It
operates in two modes:

1. **Non-interactive mode** - title provided as argument, flags control all options
2. **Interactive TUI mode** - no title argument, a huh-based wizard guides the user

## Command Signature

```
markata-go new [title-or-url] [flags]
```

### Flags

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--template` | `-t` | Content template to use | `post` |
| `--list` | `-l` | List available templates | `false` |
| `--dir` | | Output directory (overrides template) | Template default |
| `--draft` | | Create as a draft | `false` |
| `--tags` | | Comma-separated list of tags | `""` |
| `--plain` | | Use plain text prompts instead of TUI | Auto-detected |
| `--no-input` | | Disable prompts and interactive UI | `false` |
| `--title` | | Explicit title override, including URL-first link posts | `""` |

Every file created by `new` includes a `date` and an explicit `published` value.
This applies to all templates, aliases, URL-first links, plain prompts, and the TUI.
An absent, null, or blank date defaults to the current local date as `YYYY-MM-DD`.
An explicit starter date remains unchanged.
Without an explicit publication setting, `published` defaults to `true`.
A starter's `published: false` remains unpublished.
A starter's `draft: true` supplies `published: false` when `published` is absent.
The `--draft` flag overrides starter publication settings.
The command does not add a `draft` field unless the starter or an explicit override supplies it.

## URL-First Link Posts

```bash
markata new --template link 'https://example.com/my-article'
markata new --template link 'https://example.com/my-article' --title 'My title' --no-input
```

For `link` and its aliases, an HTTP(S) argument identifies the destination URL.
Existing non-URL title arguments retain their behavior.
Explicit user values take priority over metadata and inferred values.
The command preserves the entered destination, even after redirects or canonical metadata.

The command fetches the page with a bounded timeout and response size.
It extracts the title, description, and preview image.
A usable `og:title` supplies a trusted automatic title.
For recognized YouTube video URLs, the command first uses YouTube's official oEmbed endpoint.
A usable title from that endpoint is also trusted.
The response supplies the video title and thumbnail without a watch-page download.
This lookup supports watch URLs, `youtu.be` URLs, and Shorts URLs.
The command preserves the entered URL and explicit title overrides.
YouTube oEmbed does not supply a description, so a missing description remains optional.
If oEmbed fails, the command reports the error and tries the existing HTML metadata path.
Both paths share a bounded timeout and retain the response-size limit.
If both paths fail, the suggested title includes the video ID instead of the generic word `watch`.
Empty titles, URL-only titles, and hostname-only titles are not trusted.
Other title suggestions come from Twitter metadata, the HTML title, the main heading,
the decoded URL slug, and finally the hostname, in that order.
The command resolves relative preview images against the fetched page URL.
Missing optional metadata does not prevent creation.
The command does not invent a description or image.

Without a trusted title or an explicit title override, interactive creation asks for a title.
The prompt includes the best available suggestion.
If the fetch fails, the command explains the failure and continues with manual fields.
Without a terminal, or with `--no-input`, low-confidence creation fails before it writes a file.
The error includes a copy-pastable command with the suggested `--title` override.
The recovery command preserves explicit flags and the destination URL.
Invalid URLs produce an actionable error.

The interactive wizard selects the template before it asks for content.
For links, it asks for the URL before the title.
After metadata lookup and any required title prompt, **Create** is the default action.
An optional **Edit details** action exposes the title, description, image, directory,
tags, privacy, and authors.
Existing starter defaults remain intact unless metadata or explicit input replaces a field.
Explicit flags also take priority in the wizard.
Other content types retain their existing fields and creation behavior.
Plain prompts support URL-first links and the same title-confidence rules.

## Interactive TUI Wizard

When no title argument is provided, the command enters interactive mode using
charmbracelet/huh forms. The wizard falls back to plain text prompts when:

- The `--plain` flag is specified
- stdin is not a TTY

If `--no-input` is set, the command MUST NOT prompt or launch the TUI wizard.
When required input is missing, it MUST return an actionable error.
URL-first link errors include the title suggestion and a recovery command.

### Wizard Flow

```
┌─────────────────────────────┐
│  1. Template (select)       │  ← Pick from available post types
├─────────────────────────────┤
│  2. Title (text input)      │  ← Links use the URL-first flow instead
├─────────────────────────────┤
│  3. Directory (select)      │  ← Pick from existing dirs or use default
├─────────────────────────────┤
│  4. Tags (multi-select)     │  ← Existing site tags as suggestions
├─────────────────────────────┤
│  5. Private? (confirm)      │  ← Whether post is private
├─────────────────────────────┤
│  6. Authors (multi-select)  │  ← Only shown if multi-author configured
├─────────────────────────────┤
│  7. Summary + Confirm       │  ← Review before creation
└─────────────────────────────┘
```

### Step Details

#### Title (Required for Non-Link Content)

- `huh.Input` with placeholder text
- Cannot be empty (validated)
- Generates slug automatically

#### Template Selection (First)

- `huh.Select` showing all available templates
- Default: `post` (or value from `--template` flag)
- Templates sourced from: builtins, config, content-templates/ directory
- Built-in supported content types are shown with aliases when available, for example `note (aka: ping, thought, status, tweet)`
- Each option shows: `name (aka: aliases) -> directory/ (source)` when aliases exist, otherwise `name -> directory/ (source)`
- The command accepts either the canonical built-in type or any built-in alias via `--template`; aliases are written back to frontmatter as the selected `template`/`templateKey` value

#### 3. Directory Selection

- `huh.Select` presenting:
  1. The template's default directory (first option, marked as default)
  2. All existing content directories discovered from the filesystem
  3. A "Custom..." option that reveals a text input
- Default directories per template type are configurable in config under
  `[content_templates.placement]`
- Default directory pattern: `pages/<template>` for non-standard templates

#### 4. Tag Selection

- `huh.MultiSelect` populated with all existing tags from the site
- Tags are collected by scanning existing content files using the glob patterns
  from config
- Additional custom tags can be typed in (text input after multi-select)
- Tags sorted alphabetically

#### 5. Private Flag

- `huh.Confirm` asking "Is this post private?"
- Default: `No`
- When true, sets `private: true` in frontmatter

#### 6. Authors (Conditional)

Only shown when **all** of these conditions are met:
- The config has `[authors]` section with multiple authors defined
- At least 2 authors exist in config

When shown:
- First asks: "Use default author only?" (confirm)
  - If yes, uses the author marked `default: true` (or first active)
  - If no, shows `huh.MultiSelect` with all configured authors
- Author IDs are written to `authors:` frontmatter field

#### 7. Summary

- `huh.Note` showing a summary of all selections
- `huh.Confirm` to proceed or cancel

## Generated Frontmatter

The generated file includes these fields:

```yaml
---
title: "The Post Title"
slug: the-post-title
date: "2024-01-15"
published: true
draft: false
private: false
tags:
  - tag1
  - tag2
template: post
description: ""
authors:              # Only if multi-author configured
  - author-id
---

Write your content here...
```

The generated markdown body MUST NOT include a default H1 heading. Page templates already render the post title as the document H1, so generated starter content begins directly with the template body.

### Defaults

| Field | Default |
|-------|---------|
| `published` | `true` |
| `draft` | `false` |
| `private` | `false` |
| `description` | `""` |
| `tags` | `[]` |
| `template` | Selected template name |

## Configuration

### Default Directories per Template Type

```toml
[content_templates.placement]
post = "pages/post"
page = "pages"
docs = "docs"
article = "pages/article"
note = "pages/note"
# etc.
```

When no placement is configured for a template, the default directory is
`pages/<template>`.

### Author Configuration

```toml
[authors]
generate_pages = true

[authors.authors.alice]
id = "alice"
name = "Alice"
default = true
active = true

[authors.authors.bob]
id = "bob"
name = "Bob"
active = true
```

## Tag Discovery

Tags are collected from existing content by:

1. Loading config to get glob patterns
2. Scanning all matching markdown files
3. Extracting tags from frontmatter
4. Deduplicating and sorting alphabetically

This runs once at wizard startup. Performance is acceptable for typical sites
(< 1000 files).

## Theme Integration

The TUI wizard uses the site's configured palette for styling (same as `init`
command), falling back to the default Charm theme if no palette is configured.

## Non-Interactive Mode

When title is provided as an argument:

```bash
markata-go new "My Post" --template page --tags "go,web" --dir custom-dir
```

All options come from flags. No prompts are shown. The defaults apply for any
unspecified flags.

When `--site-dir` or `MARKATA_GO_SITE_DIR` selects a site, the command resolves
templates and output directories from that site and prints the created file's
absolute path.

When `--no-input` is set and no title argument is provided, the command returns
an error instead of prompting.

## Error Handling

| Error | Behavior |
|-------|----------|
| File already exists | Error with path |
| Invalid template name | Error listing available templates |
| Empty title (interactive) | Validation prevents proceeding |
| Config not found | Use defaults, skip author features |
| No TTY + no --plain | Fall back to plain mode |
