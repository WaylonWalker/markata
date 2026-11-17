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

Both `markata-go new` and the web workbench write a `date` and an explicit
`published` value on every new post. A missing, null, or blank date uses the
current local date. An explicit starter date is preserved. When the starter
does not set publication, new posts use `published: true`, except a starter
with `draft: true` uses `published: false`. Explicit starter publication values
are preserved.

These creation defaults also apply to CLI plain prompts and the TUI.

The editor includes the built-in types from `markata-go new` and the Shot and
Gratitude variants:

| Type | Default folder | Aliases |
|---|---|---|
| `post` | `pages/post` | — |
| `page` | `pages` | — |
| `docs` | `docs` | — |
| `article` | `pages/article` | `blog-post`, `essay`, `tutorial` |
| `note` | `pages/note` | `thought`, `status`, `tweet` |
| `ping` | `pages/ping` | — |
| `photo` | `pages/photo` | `shots`, `image`, `gallery` |
| `shot` | `pages/photo` | — |
| `video` | `pages/video` | `clip`, `cast`, `stream` |
| `link` | `pages/link` | `bookmark`, `til`, `stars` |
| `quote` | `pages/quote` | `quotation` |
| `guide` | `pages/guide` | `series`, `step`, `chapter` |
| `inline` | `pages/inline` | `micro` |
| `gratitude` | `pages/inline` | — |
| `contact` | `pages/contact` | `character`, `person` |
| `author` | `pages/author` | — |

The editor's Blog post and Thought starters remain available alongside these
types. Ping is listed separately from Note so the generated frontmatter keeps
`template: ping` for sites with a dedicated Ping feed. Ping and Note use
separate folders and share the same note card layout. Shot and Photo share a
card layout; Gratitude and Inline share a
card layout. Site-authored and configured starters can add or override starters.

The Link starter includes a title, date, unpublished status, source URL, and
image field. Enter an HTTP(S) source URL in the new-post form. The editor saves
it as `url` in frontmatter and does not fetch remote metadata.

The new-post form lets you set a post date. Leave it blank to keep a date from
the starter or use today's local date. Photo and Shot posts accept an image URL
or a pasted/dropped PNG, JPEG, GIF, or WebP file. URLs stay remote; pasted files
are uploaded according to `content_templates.image_upload`.

Ping and Gratitude use numbered filenames by default. Choose **Title slug** to
override this for one post. The new-post form can remember changed folder and
filename defaults for the selected type; those values are shared with
`markata-go new`.

To configure folder and filename defaults directly in TOML:

```toml
[content_templates.placement]
ping = "pages/ping"
note = "pages/note"

[content_templates.filename_mode]
ping = "numbered"
gratitude = "numbered"
note = "title"
```

```toml
[content_templates]
image_upload = "auto"
```

In `auto` mode, uploaded files go to Dropper when the editor server has
`DROPPER_TOKEN`; otherwise they are copied to `static/images/`. Choose `local`
or `dropper` to set the destination explicitly. URL images stay remote and are
only previewed in the form.

Site-authored starters can also set their numbering prefix and initial mode in
starter metadata. The editor honors these fields:

```yaml
starter:
  name: ping
  number_prefix: ping
  numbered_by_default: true
```

The editor and CLI scan the selected output folder and use the next available
`<type>-N` filename when numbering is enabled. The editor uses the optional
`number_prefix` for site-authored starters. Creation never replaces an existing
Markdown file or image.

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
