# Configuration

Use this topic when the task involves `markata-go.toml`, environment overrides, feed setup, content index output, image-library output, theme settings, or figuring out which config value is active.

## Preferred File

For standalone sites, prefer `markata-go.toml` unless the repo is already using YAML or JSON.

When a site has many feeds or environment-specific settings, keep `markata-go.toml` as the entrypoint and split the rest with `[markata-go].include`.

For blogroll work, remember that explicit feed config overrides such as `site_url`, `image_url`, `title`, `handle`, aliases, category, `type`, and tags should take effect even when cached feed data is reused. Do not assume a stale cache must be deleted before config fixes apply.

If the site uses `[markata-go].include`, treat the included file that already owns `[markata-go.blogroll]` as the canonical place to append or edit blogroll feeds instead of rewriting the root config entrypoint.

## Config Discovery

When `--config` is not passed, markata-go looks for config in this order:

1. `markata-go.toml`
2. `markata-go.yaml`
3. `markata-go.yml`
4. `markata-go.json`
5. `~/.config/markata-go/config.toml` (user-level fallback)

Without a config file, build and serve use built-in defaults plus
`MARKATA_GO_*` overrides. `config get` and `config set` need an existing file;
run `markata-go config init` to create a local one. The default output directory
is `output/`.

## High-Value Commands

- `markata-go config show`
- `markata-go config show` (annotated YAML by default; use `--no-annotate` for plain YAML)
- `markata-go config show --diff`
- `markata-go config get <key>`
- `markata-go config validate`

Bare `markata-go config` behaves like `markata-go config show`.

## Navigation previews

For a header or footer nav link that should show the default theme's hover preview, use the local URL of a generated feed or published post, such as `/blog/` or `/about/`. Set a feed's `description` in `[[markata-go.feeds]]`; set a post's `description` and `tags` in frontmatter. Enabled Reader, Blogroll, and Random routes also receive previews. The preview draws public counts and available reading statistics from the build. External links and unmatched paths render as ordinary links. Check the generated URL and plugin enablement when a preview is missing.

Post body links and wikilinks to public posts or feeds show the same metadata on hover and keyboard focus. A configured-site URL or absolute localhost URL resolves by its path during local development. Post tag links show feed counts and publication rhythm when a public `/tags/<slug>/` feed exists. Keep descriptions and dates in frontmatter so these previews have useful content.

## Core Namespace

Most settings live under `[markata-go]` and nested namespaces like `[markata-go.glob]` and `[markata-go.theme]`.

TOML, YAML, and JSON files support the same core configuration paths. Check
these settings when a feature appears to ignore the site config:

- `head`
- `template_presets`
- `default_templates`
- `theme_calendar`
- `error_pages`
- `resource_hints`
- `markdown.highlight`
- `components.post_copy.enabled` (set to `false` to remove the post-header copy menu; share buttons and post-format files remain available)

Use `markata-go config show` or `markata-go config get <key>` to inspect the
resolved typed values. Explicit `false` values are valid settings and must not
be replaced by defaults.

Full-site builds automatically write `/llms.txt` from the site description
and a short list of published, indexable pages. Configure it via
`[markata-go.llms_txt]` (`enabled = true`, `max_pages = 30` by default).
Prefer a published `llms.md` with `slug: llms` or a project
`static/llms.txt` when you need a hand-curated index; those override the
built-in generator. This file helps agents discover content; it is not an
AI training/crawling policy. Do not include draft, skipped, unpublished,
noindex, or private posts in publicly distributed indexes.

The optional Content Index is enabled with `[markata-go.content_index]`. It
writes resolved metadata and feed membership, including safe metadata-only
records for private posts. It never writes article bodies or encryption keys.
Check the output path and privacy behavior before publishing it.

The Image Library is enabled by default with `[markata-go.images]`. It writes a
searchable `/images/` page plus `/images/index.json` and `/images.json` below
`output_dir` without
fetching remote images or videos. Check `images.path`, `images.export_json`, and
`images.include_unreferenced` when media is missing or an inventory should
not be public. The page searches filenames, alt text, figure captions, and
public post titles. `added_at` is the earliest valid date from a public post
that uses the source, while `last_used_at` is the latest; filesystem timestamps
are not published. The `Recently added` view uses a 30-day browser-time window.
Private, draft, skipped, and unpublished posts must not create public usage
relationships. Private bodies are never scanned by the image library. Only a
private post's public-safe `cover` and `cover_alt` frontmatter may add a
metadata-only record, with no private use relationship or dates. The default
`images.include_unreferenced` is `false`; set it to `true` explicitly for a
local authoring inventory because it makes every supported file below
`assets_dir` centrally enumerable. Media below `assets_dir` is still public
static output; keep truly private files outside that directory.

Example:

```toml
[markata-go]
title = "My Site"
url = "https://example.com"
output_dir = "public"
assets_dir = "static"
templates_dir = "templates"

[markata-go.glob]
patterns = ["posts/**/*.md", "pages/*.md"]
use_gitignore = true
```

## Safe Defaults For New Sites

```toml
[markata-go]
title = "My Site"
url = "https://example.com"
output_dir = "public"
assets_dir = "static"
templates_dir = "templates"

[markata-go.glob]
patterns = ["posts/**/*.md", "pages/*.md"]
```

## Common Keys To Check

- `output_dir`
- `url`
- `templates_dir`
- `assets_dir`
- `concurrency`
- `hooks`
- `disabled_hooks`
- `assets.mode`
- `glob.patterns`
- `theme.palette`
- `theme.text_size`
- `theme.show_text_size_control`
- `head`
- `template_presets`
- `default_templates`
- `theme_calendar`
- `error_pages`
- `resource_hints`
- `markdown.highlight`
- `images.enabled`
- `images.path`
- `images.template`
- `images.export_json`
- `images.include_unreferenced`
- `layout.name`
- feed definitions under `[[markata-go.feeds]]`
- `templates.media.trusted_domains` when a site serves media from a CDN and needs trusted URLs normalized to `https`

## Useful Site Patterns

### Blog-like Site

```toml
[markata-go]
title = "My Blog"
url = "https://example.com"

[markata-go.glob]
patterns = ["posts/**/*.md", "pages/*.md"]

[[markata-go.feeds]]
slug = "blog"
title = "Blog"
filter = "published == True"
sort = "date"
reverse = true
```

### Docs-like Site

```toml
[markata-go]
title = "Project Docs"
url = "https://docs.example.com"

[markata-go.glob]
patterns = ["docs/**/*.md"]

[markata-go.layout]
name = "docs"
```

For feed-specific patterns, read `../reference/feed-patterns.md`.

## Layout Config Is Often The Right Answer

If the task sounds like page structure rather than full template replacement, inspect layout config first.

Example:

```toml
[markata-go.layout]
name = "blog"

[markata-go.layout.paths]
"/docs/" = "docs"
"/about/" = "landing"

[markata-go.layout.blog]
show_toc = true          # TOC sidebar on blog-layout (post.html) pages
show_prev_next = true

[markata-go.layout.docs]
sidebar_position = "left"
toc_position = "right"
```

## Hooks And Plugin Toggles

If a feature seems absent, check whether it was disabled in config before editing templates or content.

```toml
[markata-go]
hooks = ["default"]
disabled_hooks = []
```

Use `disabled_hooks` to isolate plugin-related issues. Use `markata-go explain plugins` to see available hook names.

Optional source-backed docs can also be enabled in config. For Python API reference pages:

```toml
[markata-go]
hooks = ["default", "python_docs"]

[markata-go.python_docs]
enabled = true
patterns = ["src/**/*.py"]
slug_prefix = "api"
template = "docs"
```

This generates posts from Python modules without replacing normal markdown content loading, but it only works when `python_docs` is explicitly listed in `hooks`.

Some optional hooks also depend on adjacent config namespaces. For example, when `webawesome` uses `source = "vendor"`, the shared `[markata-go.assets]` settings control the cache location and published vendor root for the downloaded Web Awesome tarball.

## Web Awesome

Use this when the site needs Web Awesome components such as tabs, details, comparisons, badges, tooltips, carousels, QR codes, or animated images.

Typical config:

```toml
[markata-go]
hooks = ["default", "webawesome"]

[markata-go.webawesome]
enabled = true
source = "vendor"
version = "3.5.0"
output_dir = "assets/vendor/webawesome"
theme = "default"
palette = "default"
brand = "blue"
```

Important rules:

- the `webawesome` hook must be enabled before expecting `wa-*` container shortcuts to render
- `source = "vendor"` is the default self-hosted mode and uses the shared `[markata-go.assets]` cache and vendor root
- `source = "cdn"` skips local publishing and serves the loader and CSS from `cdn_base_url`
- `webawesome.output_dir` controls the published Web Awesome subpath under the shared vendor root; if the site sets it to `assets/vendor/wa-kit`, pages will load `/assets/vendor/wa-kit/...`
- raw `<wa-*>` HTML also enables assets automatically, so do not add duplicate hardcoded script tags before checking whether the plugin already wires them in

Useful checks when debugging:

- `markata-go config get hooks`
- `markata-go config get webawesome`
- `markata-go config get assets`
- inspect whether the rendered page actually contains `<wa-*>` elements
- inspect whether the resolved config exposes the expected vendor root and Web Awesome `output_dir`

## Authors Config

If the site has multiple authors, configure them under `[markata-go.authors]`:

```toml
[markata-go.authors]
generate_pages = true
url_pattern = "/authors/{author}/"
feeds_enabled = true

[markata-go.authors.authors.waylon]
name = "Waylon Walker"
bio = "Software engineer"
avatar = "/images/waylon.jpg"
url = "https://waylonwalker.com"
default = true
active = true

[markata-go.authors.authors.waylon.social]
github = "https://github.com/waylonwalker"
```

Key fields per author: `name`, `bio`, `email`, `avatar`, `url`, `social` (map), `guest`, `active`, `default`, `role`, `contributions` (CReDiT roles).

When `generate_pages = true`, author profile pages are auto-generated. When `feeds_enabled = true`, per-author feeds are created.

## Environment Overrides

Environment variables use the `MARKATA_GO_` prefix.

Examples:

```bash
MARKATA_GO_URL=https://staging.example.com markata-go build
MARKATA_GO_OUTPUT_DIR=dist markata-go build
```

## Config Merging

Local override files like `markata-go.local.toml` are NOT auto-discovered. They must be passed explicitly via the `--merge-config` / `-m` flag:

```bash
markata-go build -m markata-go.local.toml
markata-go serve -m markata-go.local.toml
markata-go config show -m markata-go.local.toml
markata-go config validate -m markata-go.local.toml
```

Multiple merge files can be specified and are applied in order on top of the base config.

## Config Composition

The main config can include other config files:

```toml
[markata-go]
include = ["config/common/*.toml", "config/feeds/*.toml"]
```

Key rules:

- include paths are relative to the file that declared them
- glob matches load in lexicographic order
- repeated includes are loaded once
- feed entries merge by `slug`
- environment variables still win last

## Guidance

- Prefer editing the site's real config file over hardcoding behavior in templates.
- Confirm whether a value is user-defined or default before changing it.
- Keep related settings together under existing namespaces.
- For feed work, inspect both feed definitions and feed defaults.
- When debugging paths, remember CLI flags override config.
- If a change is site-wide, config is usually the right place before template logic or plugin logic.
