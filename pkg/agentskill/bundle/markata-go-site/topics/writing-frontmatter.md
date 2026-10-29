# Writing And Frontmatter

Use this topic when creating or editing posts, pages, docs, or frontmatter fields.

## Preferred Commands

- `markata-go new`
- `markata-go new --list`
- `markata-go list posts`
- `markata-go explain content`

## Frontmatter Format

Markata-go uses YAML frontmatter between `---` delimiters at the top of the file.

Example:

```markdown
---
title: "Getting Started with Go"
date: 2026-04-01
published: true
draft: false
tags:
  - go
  - tutorial
description: "A beginner-friendly guide."
template: post
---

Write your content here.
```

## Rich Markdown authoring

For icons, card layouts, feature rows, tables, status metrics and
progressive enhancements, read `topics/rich-markdown.md` before inventing
custom HTML or JavaScript. `:lucide-rocket:` and `:lucide-database:` are
built-in icon shortcodes; `md-feature` and `md-icon` are CSS-only layouts.
Keep visible headings alongside icons for accessibility.

## Important Default Behavior

- files without frontmatter are still valid
- `published` defaults to `false` if omitted in raw content
- the `new` command creates content with `published: true` and `draft: false`
- new starter content does not add a default H1 in the body because templates usually render the title as the page H1

## Behavior Rules Agents Should Remember

- `published: true` means the post is eligible for public feeds and listings
- `published: false` can still render a shadow page reachable by direct URL
- `draft: true` means the post should not render at all
- `skip: true` means the file is ignored during processing
- `slug: ""` or `slug: /` makes a content file the homepage
- duplicate slugs create output conflicts and should be treated as a build blocker

## Common Built-In Fields

- `title`
- `slug`
- `date` (also recognized as `publishdate` or `pubdate`)
- `published`
- `draft`
- `private`
- `tags`
- `description`
- `template`
- `layout`
- `skip`
- `authors`
- `author` (also recognized as `by` or `writer`)
- `modified` (also recognized as `lastmod`, `updated`, `updated_at`, `last_modified`)
- `aliases` (also recognized as `alias`, `handles`, `handle`)

## Date Field Aliases

Multiple frontmatter keys map to the same date fields. The first match wins:

- **Publication date**: `publishdate` > `date` > `pubdate`
- **Modified date**: `lastmod` > `modified` > `updated` > `updated_at` > `last_modified`

Dates are parsed flexibly: RFC3339, `2006-01-02`, `January 2, 2006`, and other common formats all work.

## Aliases Field

The `aliases` field (a string list) registers alternate names for a post. These are used by:

- **Wikilinks**: `[[JavaScript]]` resolves to a post with `aliases: ["JavaScript"]` even if its slug is `javascript-guide`
- **Glossary**: terms can match on aliases in addition to the primary slug
- **Mentions**: blogroll sources use aliases for `@mention` resolution

All four keys (`aliases`, `alias`, `handles`, `handle`) are merged and deduplicated.
Alias resolution must remain identical across cold, cached, and edited builds;
do not remove aliases or clear caches to work around missing wikilinks.

## Media Fields

Media fields are not typed on the Post model. They live in `post.Extra` and are read by various plugins. Common recognized keys:

- `image` (primary, most widely used across plugins)
- `cover_image` (cover image alias)
- `og_image` (Open Graph image)
- `featured_image` (wikilink hover previews)
- `thumbnail` (wikilink hover previews)
- `cover` (feed card helpers)
- `social_image` (social metadata and image-library usage)
- `hero_image` (hero image and image-library usage)
- `avatar` and `author_image` (profile media and image-library usage)
- `video` (video URL for OG cards)
- `poster_image`, `poster`, `video_poster`, `video_thumbnail`, `thumbnail`, and
  `thumb` (video poster candidates)

For OG card generation, the first non-empty value from `image`, `cover_image`, `og_image` is used.
The generated Image Library records the supported image fields above for public
posts. It never scans private post bodies. A published private post may expose
only its explicitly public-safe `cover` and `cover_alt` fields as a metadata-only
image record; it does not create a public usage relationship. Draft, skipped,
and unpublished posts do not participate.

The Image Library treats `cover`, `cover_image`, `image`, and `video` as cover
fallbacks in that order. It detects video media from the URL extension or MIME
type and exposes poster metadata for custom library templates. Render a video
media card with `<video>` and `<source>` instead of `<img>`.

## Extended Authors Format

The `authors` field supports three formats:

Simple string:

```yaml
authors: waylon
```

List of IDs:

```yaml
authors:
  - waylon
  - codex
```

Structured with per-post role overrides:

```yaml
authors:
  - id: waylon
    role: author
    details: "wrote the introduction"
  - id: codex
    role: pair programmer
    details: "wrote the code examples"
  - guest
```

Accepted key aliases in structured entries:

| Logical field | Accepted keys |
|--------------|---------------|
| id | `id`, `name`, `handle` |
| role | `role`, `job`, `position`, `part`, `title` |
| details | `details`, `detail`, `description` |

Author IDs are resolved against `[markata-go.authors]` config to produce full `post.author_objects` for templates.

## Slug Guidance

- if omitted, slug is generated from title or filename
- `index.md` gets a directory-based slug
- `slug: ""` or `slug: /` creates the homepage
- duplicate slugs cause output conflicts

## Field Effects

- `title`: affects page title, feed listings, and slug generation
- `slug`: affects final URL path and output location
- `date`: affects sorting, scheduling, and feed ordering
- `published`: affects discoverability in feeds, sitemaps, and public listings
- `draft`: affects whether output is rendered at all
- `tags`: affects filtering and tag/archive pages
- `description`: affects meta tags, SEO, and some feed metadata, but not the main excerpt generation path
- `template`: selects a specific template file
- `layout`: selects a layout-driven page shell when the site uses layout config
- `authors` or `author`: affect bylines and author objects in templates
- `post_formats`: can override site-level post output formats per post, such as enabling `ansi` for one article or disabling `og`/`markdown` for a specific page
- custom fields: appear in `post.Extra` for templates and plugin logic

Example per-post output override:

```yaml
---
title: "CLI cheat sheet"
post_formats:
  ansi: true
  markdown: false
---
```

## Recommended Content Workflow

Use `markata-go new` when possible because it knows current built-in content templates such as:

- `post`
- `page`
- `docs`
- `article`
- `note`
- `photo`
- `video`
- `link`
- `quote`
- `guide`
- `inline`
- `contact`
- `author`

## Common Site Patterns

### Public blog post

```yaml
---
title: "My Post"
date: 2026-04-01
published: true
draft: false
tags:
  - go
description: "Short SEO summary"
---
```

### Shadow page

```yaml
---
title: "Reviewer Notes"
published: false
draft: false
---
```

### Truly private draft

```yaml
---
title: "Unpublished Draft"
draft: true
---
```

### Homepage from content

```yaml
---
title: "Welcome"
slug: ""
published: true
layout: landing
---
```

### Page with custom template

```yaml
---
title: "API Landing"
published: true
template: "landing.html"
---
```

## Guidance

- Prefer `markata-go new` when creating content so the site's defaults stay consistent.
- To edit existing content, use `markata-go search <query> --format path` or `markata-go list posts --format path`, then edit the returned Markdown file. With `--site-dir`, those paths are absolute.
- Keep frontmatter YAML valid and minimal.
- Use existing frontmatter conventions from nearby content before introducing new fields.
- Do not duplicate the title with an extra H1 in the body unless the repo already does that intentionally.
- If you introduce a new custom field, check whether templates and feeds need to read it.
- Use `draft: true` for content that must not render at all; use `published: false` for shadow pages.
- Prefer `authors` over legacy `author` when the site already uses structured author config.

## Before Adding New Metadata

1. Search nearby content for the same field.
2. Check whether templates or feeds already read it.
3. Prefer existing naming over introducing another alias.
4. Decide whether the field belongs in frontmatter, config, or computed plugin output.

## Editorial attribution

For copy editing, internal linking, or recording a human's edits, use
`markata-go editorial start <post...> --author ID --role copy-editor` before
editing clean posts, then `editorial finish --session UUID --summary "..."`.
For explicitly owned existing changes, use `editorial commit <post...>` with
`--author`, `--role`, and `--summary`. Inspect `--dry-run` before committing.
Use `--json` for the session and complete list of resulting commit SHAs.

Use the actual contributor ID, even when an agent runs commands on a human's
behalf. Preserve previous `authors` and append-only `editorial_history`.
The first commit contains edits and credit; the second records exact file and
diff links to the first. Report every resulting commit to the user. Push those
commits before expecting the generated GitHub links to resolve.

The default theme shows editorial history and groups editors under the author.
Inline comparisons fetch public GitHub revision files only after a reader
requests them. Custom post templates need the editorial history component and
its conditional CSS/bootstrap assets. Full rendering and authoring contracts
are documented in the site's Editorial history and Editorial authoring guides.

## Workbench snippet commands

The Python, Go, and WASM workbenches offer **Commands** above the editor.
Use **Ctrl+Shift+P** (**Cmd+Shift+P** on macOS) to search the snippet catalog.
Commands preserve selected text and insert Markdown in one undo step.
The catalog includes admonitions, Mermaid, Chart.js, and inline formatting.
These standalone editors do not load a site index for post embeds.

The Python and Go renderer examples include a bundled example site with Files and
Feeds navigation. These examples save per-file drafts in browser storage, not in
the site repository. Download changes before applying them to a real site.
