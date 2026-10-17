---
title: "Configuration Reference"
description: "Complete reference for all markata-go configuration options"
date: 2026-01-24
published: true
slug: /docs/guides/configuration/
tags:
  - configuration
  - reference
---

# Configuration

markata-go uses a flexible configuration system that supports multiple file formats, environment variable overrides, and intelligent merging from multiple sources.

## Compose Config Files

Use `include` when one config file gets too large.

```toml
[markata-go]
title = "Example Site"
include = [
  "config/common/*.toml",
  "config/feeds/*.toml",
  "config/rss.toml",
]
```

Rules:

- paths are resolved relative to the file that declared the include
- explicit files are loaded as written
- glob matches are loaded in lexicographic order
- repeated includes are loaded once
- cycles fail with a clear error

## Edit Settings While Serving

While `markata-go serve` is running, every page has a small gear button in the bottom-left corner (or press `Alt+,`). It opens a **Site settings** sidebar listing every config option with its current value and description. Change a setting and the site rebuilds and reloads with it right away. The change is only a preview until you click **Bake** to write it into your config, or **Reset** to throw it away. `markata-go build` never includes the sidebar, so published sites don't have it. The sidebar and Bake only accept changes from the machine running `serve`; if you bind to `--host 0.0.0.0` to test on a phone, the phone can browse but cannot edit settings.

1. Run `markata-go serve` and open the site.
2. Click the gear. The **Common** section at the top has the settings most sites change (title, URL, palettes, style, font, nav, footer). Search to find anything else, for example `items_per_page`.
3. Change a value. Text fields apply when you press Enter or leave the field; dropdowns, checkboxes, and numbers apply right away. The status line shows **Rebuilding… 3s** until the build finishes, then the page reloads with the change and the sidebar stays open where you left it. If the build fails, the error is shown instead.
4. Keep going. Changed rows are highlighted, and **Revert** undoes one field. **Reset to default** removes a setting from your config so it falls back to the default shown on the row. While the sidebar is closed, a small **Previewing N unsaved changes · Reset · Review** bar next to the gear reminds you the site is showing unsaved changes.
5. Click **Bake** to keep them. The sidebar shows the exact diff for every file it will edit (secret values are hidden). Click **Write** to save, or **Cancel**. The server writes the files, rebuilds, and live-reloads.
6. Or click **Reset** to drop every unsaved change and go back to what your config files say.

Previews live in the serve process's memory only. Your files don't change until you bake. If you restart `markata-go serve`, the sidebar re-applies the edits you had staged in that browser tab. While a preview is active, builds use a separate cache in `.markata/serve-preview/`, so previewing never throws away your site's warm build cache.

If you changed the theme with the theme picker, the picker's choice is saved in your browser and would hide a palette, style, font, text size, or color mode preview. Previewing one of those settings clears the matching picker choice, and the status line tells you it did.

When you run `markata-go serve <file>` for a single page, sections for site-wide output such as feeds, tags, and the blogroll are hidden. Search or click **Show hidden sections** to see them.

Settings with a fixed set of values, such as `glob.slug_mode`, `theme.palette`, `theme.fontpack`, `nav.position`, or `markdown.highlight.theme`, are dropdowns. Numbers with a range, such as `theme.background.color_mix` (0-1), are limited to that range. A value that isn't allowed gets a red message under the field and is not previewed, so a typo can't break the site; the status line lists every change that is not being previewed. Leave a field empty to use the default.

Each row shows where a setting lives: **set in config/theme.toml** means that file defines it now, and **writes to ...** names the file a change will go to. The sidebar edits the file that owns the setting:

- A key that is already set is edited where it is set (the file that wins, when several set it).
- A new key goes to the file that already holds its group. For example, a new `feed_defaults.items_per_page` lands in the file that has `[markata-go.feed_defaults]`.
- Otherwise it goes to your root config, or a new `markata-go.toml` when you have none.
- A file passed with `--merge-config` is marked **override**, so you know the change won't be in your base config.
- If your site has no config of its own and serve fell back to your global `~/.config/markata-go/config.toml`, the rows are marked **global** and Bake is refused, so a preview never changes the defaults for all your sites. Create a `markata-go.toml` in the site to bake.

This works with `include = [...]` files, config directories, and `--merge-config` files, in TOML, YAML, or JSON. Comments, key order, and other settings are kept.

Some settings are shown but not editable in the sidebar; edit them in the file:

- Lists of tables and maps, such as `nav`, `feeds`, and `head.meta`, show a summary like "3 items".
- Secrets such as tokens and passwords are never sent to the browser.
- A few fields that the config loader does not read are marked **Not read from config files**.

Changes are checked before they stick. If a preview or bake makes the config invalid or does not load back as written, it is rejected (and a bake restores every file) and the sidebar shows the error. When a `MARKATA_GO_*` environment variable overrides a setting, the row says so, because the variable wins over both previews and files until it is unset.

On phones the sidebar opens full screen. `Alt+,` does nothing while you are typing in a field on the page.

The theme picker's **All settings** button opens this sidebar at the Theme section. Its **Bake** button writes only the look you picked, through the same checks and file selection as the sidebar (see [Themes Guide](/docs/guides/themes/#baking-a-look-into-your-config)).

## Content Index

Enable the optional metadata artifact with a plugin table:

```toml
[markata-go.content_index]
enabled = true
output = "content-index.json"
```

The output path is relative to `output_dir` unless it is absolute. The index
contains resolved public document metadata and feed membership, but never
complete article bodies. See the [Content Index guide](/docs/guides/content-index/)
for the versioned format and consumer rules. Environment overrides are also
available for isolated builds: `MARKATA_GO_CONTENT_INDEX_ENABLED=true` and
`MARKATA_GO_CONTENT_INDEX_OUTPUT=/tmp/content-index.json`.

Precedence:

1. built-in defaults
2. the root config file
3. included files in declaration order
4. glob matches in lexicographic order
5. environment variables

Later values win.

### Merge Behavior

- scalar values replace earlier values
- explicit `false` and `0` are preserved
- tables merge deeply
- arrays of scalars replace earlier arrays
- `[[markata-go.feeds]]` merges by `slug`

That last rule makes feed-heavy sites much easier to manage.

## Image Library

The image library plugin is enabled by default. It writes a searchable page to
`images/index.html` and versioned JSON inventories to `images/index.json` and
`images.json`, relative to `output_dir`.

```toml
[markata-go.images]
enabled = true
path = "images"
template = "images.html"
export_json = true
include_unreferenced = false
```

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `enabled` | bool | `true` | Generate the image page and inventory |
| `path` | string | `"images"` | Relative output directory below `output_dir` |
| `template` | string | `"images.html"` | HTML template name |
| `export_json` | bool | `true` | Write `path/index.json` and `images.json` |
| `include_unreferenced` | bool | `false` | Include image and video files found below `assets_dir` without public post usage; enable explicitly when a local authoring inventory should enumerate all supported assets |

The path must remain inside `output_dir`. Environment variables use the
`MARKATA_GO_IMAGES_` prefix. See the [Image Library guide](/docs/guides/image-library/)
for discovery rules and template context.

## Correcting Configuration Mistakes

Run `markata-go config validate` after editing a config file. If a key closely
resembles a built-in setting, markata-go reports the file and a possible
spelling. For example, `output_dr` suggests `output_dir`, and
`theme.palete` suggests `theme.palette`. Correct the key in the file and rerun
the command; markata-go never changes the file for you. Custom plugin sections
remain supported.

`markata-go config get <key>` also suggests close key names. If the key is
valid but absent from the file, use `markata-go config show` to see its resolved
value, including defaults. For a theme choice with a fixed list of values,
validation shows allowed values and, when close enough, a possible correction.

## Configuration File Locations

markata-go searches for configuration files in the following order (first found wins):

| Priority | Location | Description |
|----------|----------|-------------|
| 1 | `--config path/to/config.toml` | CLI-specified path |
| 2 | `./markata-go.toml` | Current directory (TOML) |
| 3 | `./markata-go.yaml` | Current directory (YAML) |
| 4 | `./markata-go.yml` | Current directory (YAML alternate) |
| 5 | `./markata-go.json` | Current directory (JSON) |
| 6 | `~/.config/markata-go/config.toml` | User config directory |

If no configuration file is found, markata-go uses default values with any environment variable overrides applied.

### Start With One Markdown File

A configuration file is optional for a small site. The default content
patterns include root Markdown plus the `pages/` and `posts/` directories, so
this minimal layout works:

```text
pages/
  sample.md
```

`markata-go build` and `markata-go serve` render `sample.md`, a homepage that
links to it, and an `/archive/` page. The homepage also keeps the default RSS
and Atom endpoints at `/rss.xml` and `/atom.xml`. Posts must be published to
appear in these collections: use `published: true` in frontmatter. An empty
site still renders `/` with that reminder.

To render only one Markdown file with the bundled default theme, without
generating a homepage, archive, or feeds, pass the file directly:

```bash
markata-go pages/sample.md
# equivalent:
markata-go build pages/sample.md
```

## Select A Site From Another Directory

Use `--site-dir` when the site is not your current directory. This is useful
for aliases, scripts, and agents that manage more than one site.

```bash
markata-go --site-dir ~/sites/blog config show
MARKATA_GO_SITE_DIR=~/notes/work markata-go list posts
```

`--site-dir` overrides `MARKATA_GO_SITE_DIR`; otherwise markata-go uses the
current directory. The selected site becomes the base for configuration
discovery, relative config flags, templates, content files, and caches.

## Supported Formats

| Extension | Format | Notes |
|-----------|--------|-------|
| `.toml` | TOML | **Recommended** - Best for nested config, readable |
| `.yaml`, `.yml` | YAML | Good for complex structures |
| `.json` | JSON | Strict, good for programmatic generation |

### TOML Example (Recommended)

```toml
[markata-go]
title = "My Site"
description = "A site built with markata-go"
url = "https://example.com"
output_dir = "public"

[markata-go.glob]
patterns = ["posts/**/*.md", "pages/*.md"]
use_gitignore = true

[markata-go.markdown]
extensions = ["tables", "strikethrough", "tasklist"]
```

### YAML Example

```yaml
markata-go:
  title: My Site
  description: A site built with markata-go
  url: https://example.com
  output_dir: public

  glob:
    patterns:
      - "posts/**/*.md"
      - "pages/*.md"
    use_gitignore: true

  markdown:
    extensions:
      - tables
      - strikethrough
      - tasklist
```

### JSON Example

```json
{
  "markata-go": {
    "title": "My Site",
    "description": "A site built with markata-go",
    "url": "https://example.com",
    "output_dir": "public",
    "glob": {
      "patterns": ["posts/**/*.md", "pages/*.md"],
      "use_gitignore": true
    },
    "markdown": {
      "extensions": ["tables", "strikethrough", "tasklist"]
    }
  }
}
```

## Configuration Namespacing

All configuration lives under the `[markata-go]` namespace. This:

- Avoids conflicts with other tools in shared config files
- Provides clear ownership for each plugin section
- Enables tooling-friendly editor completions

```toml
# Root namespace - site-wide settings
[markata-go]
title = "My Site"
output_dir = "public"

# Plugin namespaces
[markata-go.glob]
patterns = ["pages/**/*.md", "posts/**/*.md"]

[markata-go.markdown]
extensions = ["tables"]

[markata-go.feeds]
# Feed-specific config
```

## Configuration Options Reference

### Core Settings (`[markata-go]`)

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `output_dir` | string | `"output"` | Build output directory |
| `url` | string | `""` | Site base URL (for absolute links) |
| `title` | string | `""` | Site title |
| `description` | string | `""` | Site description |
| `author` | string | `""` | Default author |
| `license` | string `\|` bool | (unset) | Select a license key (see below). Set to `false` to skip the attribution and warning. |
| `assets_dir` | string | `"static"` | Static assets directory |
| `templates_dir` | string | `"templates"` | Templates directory |
| `hooks` | string[] | `["default"]` | Plugins to load |
| `disabled_hooks` | string[] | `[]` | Plugins to exclude |
| `concurrency` | int | `0` | Worker threads (0 = auto based on CPU cores) |

```toml
[markata-go]
output_dir = "public"
url = "https://example.com"
title = "My Site"
description = "A site built with markata-go"
author = "Jane Doe"
assets_dir = "static"
templates_dir = "templates"
hooks = ["default"]
disabled_hooks = ["sitemap"]
concurrency = 4
```

### Python Docs (`[markata-go.python_docs]`)

Generate API reference pages from Python source files.

This plugin is strict opt-in: add `python_docs` to `[markata-go].hooks` and set `enabled = true`.

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `enabled` | bool | `false` | Enable source-backed Python docs generation |
| `patterns` | string[] | `["**/*.py"]` | Python file glob patterns |
| `directories` | string[] | `[]` | Directory roots expanded to `dir/**/*.py` |
| `exclude` | string[] | venv/cache defaults | Globs to skip |
| `slug_prefix` | string | `"api"` | Prefix for generated module pages |
| `template` | string | `""` | Template used for generated docs posts |
| `published` | bool | `false` | Include generated docs in feeds/search/sitemap |
| `include_private` | bool | `false` | Include private `_symbols` |
| `include_source` | bool | `true` | Include implementation snippets |
| `include_module_code` | bool | `false` | Include module-level implementation source |
| `tags` | string[] | `["python", "docs"]` | Tags added to generated docs posts |
| `interpreter` | string | auto-detect | Python executable used for AST extraction |

```toml
[markata-go]
hooks = ["default", "python_docs"]

[markata-go.python_docs]
enabled = true
patterns = ["src/**/*.py"]
slug_prefix = "api"
template = "docs"
include_source = true
```

See [Python Source Docs](./python-docs.md) for the full workflow.

### Embeds (`[embeds]`)

Controls internal and external embed behavior, including oEmbed resolution and Open Graph fallback.

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `enabled` | bool | `true` | Enable/disable embeds plugin |
| `internal_card_class` | string | `"embed-card"` | CSS class for internal embed cards |
| `external_card_class` | string | `"embed-card embed-card-external"` | CSS class for external embed cards |
| `fetch_external` | bool | `true` | Fetch Open Graph metadata for external embeds |
| `oembed_enabled` | bool | `true` | Enable oEmbed resolution |
| `oembed_auto_discover` | bool | `false` | Discover oEmbed endpoints from HTML link tags |
| `default_mode` | string | `"card"` | Default embed mode (`card`, `rich`, `performance`, `hover`, `image_only`) |
| `oembed_providers_url` | string | `"https://oembed.com/providers.json"` | providers.json URL (empty disables) |
| `resolution_strategy` | string | `"oembed_first"` | `oembed_first`, `og_first`, `oembed_only` |
| `cache_dir` | string | `".cache/embeds"` | Cache directory for metadata |
| `cache_ttl` | int | `604800` | Cache TTL (seconds) |
| `timeout` | int | `10` | HTTP timeout (seconds) |
| `fallback_title` | string | `"External Link"` | Title when metadata unavailable |
| `show_image` | bool | `true` | Show images for external embeds |
| `providers.<name>.mode` | string | (unset) | Provider mode override (`card`, `rich`, `performance`, `hover`, `image_only`) |

```toml
[embeds]
enabled = true
internal_card_class = "embed-card"
external_card_class = "embed-card embed-card-external"
fetch_external = true
oembed_enabled = true
oembed_auto_discover = false
default_mode = "card"
oembed_providers_url = "https://oembed.com/providers.json"
resolution_strategy = "oembed_first"
cache_dir = ".cache/embeds"
cache_ttl = 604800
timeout = 10
fallback_title = "External Link"
show_image = true

[embeds.providers]
youtube = { enabled = true, mode = "rich" }
vimeo = { enabled = true }
tiktok = { enabled = true }
flickr = { enabled = true }
spotify = { enabled = true }
soundcloud = { enabled = true }
codepen = { enabled = true }
codesandbox = { enabled = true }
jsfiddle = { enabled = true }
observable = { enabled = true }
github = { enabled = true }
slideshare = { enabled = true }
prezi = { enabled = true }
speakerdeck = { enabled = true }
issuu = { enabled = true }
datawrapper = { enabled = true }
flourish = { enabled = true }
infogram = { enabled = true }
reddit = { enabled = true }
dailymotion = { enabled = true }
```

### External Link Hover (`[markata-go.external_link_hover]`)

Adds hover preview cards to links that point to other sites. Disabled by
default. See the [plugin reference](../reference/plugins.md#external_link_hover)
for data sources and output.

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `enabled` | bool | `false` | Enable/disable external link cards |
| `fetch` | bool | `false` | Fetch metadata for URLs missing from the embeds cache (no network when `false`) |
| `include_image` | bool | `true` | Include the page's preview image |
| `favicon_service` | string | `""` | Icon URL template with `{host}` / `{origin}`; empty disables icons |
| `ignore_domains` | list | `[]` | Hosts (and subdomains) that never get cards |

```toml
[markata-go.external_link_hover]
enabled = true
fetch = true
favicon_service = "https://icons.duckduckgo.com/ip3/{host}.ico"
ignore_domains = ["localhost"]
```

### Vendor Assets (`[markata-go.assets]`)

markata-go can self-host common third-party JS/CSS dependencies (HTMX, GLightbox, Mermaid, Chart.js, Cal-Heatmap, D3, Lite YouTube). When enabled, assets are downloaded into a cache directory and copied to `/assets/vendor` in the output. Templates use the `asset_urls` mapping injected by the CDN assets plugin.

```toml
[markata-go.assets]
mode = "self-hosted"           # "self-hosted" (default), "cdn", or "auto"
cache_dir = ".markata/assets-cache"
output_dir = "assets/vendor"
verify_integrity = true
```

### License configuration

The `license` key controls the attribution shown in the footer and whether the dev server reminds you to pick a license. It accepts either a string key (selects the attribution) or the literal `false` (hides the line and silences the warning). When the footer shows the copyright line, the license appears on the same line next to the copyright symbol; if copyright is disabled, the license renders as its own line.

Supported license keys:

- `all-rights-reserved` – All rights reserved.
- `cc-by-4.0` (recommended) – Creative Commons Attribution 4.0 International (`https://creativecommons.org/licenses/by/4.0/`). The default when running `markata-go config init` or picking the recommended option in the init wizard.
- `cc-by-sa-4.0` – Creative Commons Attribution-ShareAlike 4.0 International (`https://creativecommons.org/licenses/by-sa/4.0/`).
- `cc-by-nc-4.0` – Creative Commons Attribution-NonCommercial 4.0 International (`https://creativecommons.org/licenses/by-nc/4.0/`).
- `cc-by-nd-4.0` – Creative Commons Attribution-NoDerivatives 4.0 International (`https://creativecommons.org/licenses/by-nd/4.0/`).
- `cc-by-nc-sa-4.0` – Creative Commons Attribution-NonCommercial-ShareAlike 4.0 International (`https://creativecommons.org/licenses/by-nc-sa/4.0/`).
- `mit` – MIT License (`https://opensource.org/licenses/MIT`).

When omitted, `license` defaults to `cc-by-4.0`. Set another supported string
to override it, or set `license = false` to suppress footer attribution.

### Navigation Links (`[[markata-go.nav]]`)

Navigation links appear in the site header and define your site's main navigation.

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `label` | string | Required | Display text for the link |
| `url` | string | Required | Link destination (relative or absolute) |
| `external` | bool | `false` | Opens link in new tab with noopener |

```toml
[[markata-go.nav]]
label = "Home"
url = "/"

[[markata-go.nav]]
label = "Blog"
url = "/blog/"

[[markata-go.nav]]
label = "Docs"
url = "/docs/"

[[markata-go.nav]]
label = "GitHub"
url = "https://github.com/WaylonWalker/markata-go"
external = true
```

The default theme adds a small preview when a local header or footer navigation URL matches a generated feed or published post. For a feed such as `/blog/`, hovering the link or focusing it with a keyboard shows its description, public post count, total words, reading time, and publication rhythm when those values are available. For a post, the preview uses its description, tags, word count, and reading time. Enabled Reader and Blogroll links show the number of active followed sources; an enabled Random post link shows how many posts are in its draw. External and unmatched links stay as ordinary links. Previews remain open while the pointer moves from the link into the panel, close when it leaves, and appear above page content. They are hidden on narrow or touch screens, where the navigation is a scrollable strip.

Links inside a post can use the same preview. Write a normal Markdown link such as `[Start here](/start/)` or a wikilink such as `[[start]]`. When its destination is a public post or feed, hover or keyboard focus shows the destination title and available metadata. An absolute URL for the configured site also works; `http://localhost:8000/start/` works during local development. Tags shown on a post get a collection preview when their `/tags/<tag>/` feed exists. Links to other sites and links without a matching public page keep their usual behavior.

To give a feed link useful preview text, set `description` on its `[[markata-go.feeds]]` entry. For a post link, set `description` and `tags` in that post's frontmatter. The link URL must match the generated page path, such as `/blog/` or `/about/`.

To add a built-in **Random post** link:

```toml
[markata-go.random_post]
enabled = true
path = "random"                # generates /random/
exclude_tags = ["draft"]       # optional

[[markata-go.nav]]
label = "Random"
url = "/random/"
```

If you use sidebar navigation, add the same `url = "/random/"` link to `[[markata-go.sidebar.nav]]` (or to `[[markata-go.components.nav.items]]` if you're using the components system).

### SEO Configuration (`[markata-go.seo]`)

Configure SEO metadata and structured data generation for better search engine visibility.

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `twitter_handle` | string | `""` | Twitter/X username (without @) for twitter:site |
| `default_image` | string | `""` | Default Open Graph image URL |
| `logo_url` | string | `""` | Site logo URL for Schema.org |

```toml
[markata-go.seo]
twitter_handle = "waylonwalker"
default_image = "/static/og-default.png"
logo_url = "/static/logo.png"

[markata-go.seo.structured_data]
enabled = true

[markata-go.seo.structured_data.publisher]
type = "Organization"
name = "My Company"
url = "https://example.com"
logo = "/static/logo.png"

[markata-go.seo.structured_data.default_author]
type = "Person"
name = "Jane Doe"
url = "https://example.com/about/"
```

#### Structured Data Configuration

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `enabled` | bool | `true` | Enable JSON-LD Schema.org generation |
| `publisher` | object | `nil` | Publisher information |
| `default_author` | object | `nil` | Default author for posts |

When enabled, markata-go generates JSON-LD structured data for:
- `WebSite` schema on the home page
- `Article` or `BlogPosting` schema on posts
- Breadcrumb navigation schema

### Theme Settings (`[markata-go.theme]`)

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `name` | string | `"default"` | Theme name |
| `palette` | string | `"ayu-dark"` | Color palette to use |
| `palette_dark` | string | `""` | Dark mode palette |
| `fallback_mode` | string | `"dark"` | Mode for visitors who have not toggled it (`"dark"` or `"light"`); the OS preference is not used |
| `seasonal` | boolean | `false` | Default visitors to the seasonal palette (northern hemisphere seasons plus world holidays). Needs the theme switcher; `palette` stays the fallback. See [Themes: Seasonal](/docs/guides/themes/#seasonal) |
| `text_size` | string | `"large"` | Default reading-size preset: `"small"`, `"medium"`, `"large"`, or `"x-large"`. Article text additionally scales up on viewports ≥1800px. |
| `show_text_size_control` | bool | `true` | Show the visitor-facing reading-size selector |
| `custom_css` | string | `""` | Custom CSS file path (relative to static/) |
| `variables` | map | `{}` | CSS variable overrides |
| `fontpack` | string | `"brush"` | Font pack: Knewave headings, Space Grotesk body, DM Mono code by default |
| `aesthetic` | string | `"minimal"` | Surface style: `minimal`, `balanced`, `elevated`, `precision`, or `brutal` |
| `motif.kind` | string | `"off"` | Background motif (`off`, `block-w`, `letter`) |
| `switcher.enabled` | bool | `true` | Show the live theme picker so visitors can choose any palette, style, and font; under `markata-go serve` its **Bake** button writes the current choices into your config (see [Themes Guide](/docs/guides/themes/)) |
| `switcher.mode_toggle` | bool | `true` | Show the light/dark toggle next to the picker |

The default font pack is `brush`. Built-in defaults do not cause migration warnings.
If your configuration supplies different values for `fontpack` and `theme.fontpack`, the nested value wins.
The loader reports that conflict once. Remove the legacy top-level `fontpack` to resolve the warning.

```toml
[markata-go.theme]
name = "default"
palette = "catppuccin-mocha"

# Optional: different palette for dark mode
palette_dark = "catppuccin-mocha"

# Optional: fallback when system preference is unavailable
fallback_mode = "dark"  # or "light"

# Optional: default reading size and visitor control
text_size = "large"              # small, medium, large, or x-large
show_text_size_control = true

# Optional: override specific CSS variables
[markata-go.theme.variables]
"--color-primary" = "#8b5cf6"
"--content-width" = "800px"

# Optional: add custom CSS file
custom_css = "my-styles.css"
```

**Available palettes:** `default-light`, `default-dark`, `catppuccin-mocha`, `catppuccin-latte`, `nord-dark`, `gruvbox-dark`, `dracula`, `rose-pine`, `solarized-dark`, `tokyo-night`

See the [[themes-and-styling|Themes Guide]] for detailed customization options.

### Page CSS Bundles (`[markata-go.css_bundle]`)

```toml
[markata-go.css_bundle]
enabled = true
pages = ["index.html"]
```

`pages` defaults to an empty list. Each entry names an HTML file relative to
`output_dir`. Production builds bundle eligible head styles in cascade order
after CSS optimization. Fast builds skip page bundling. Existing `bundles`
definitions continue to work independently. See [CSS optimization](themes/#css-optimization)
for stylesheet boundaries and runtime behavior.

### CSS Purge (`[markata-go.css_purge]`)

```toml
[markata-go.css_purge]
enabled = false
verbose = false
preserve = ["js-*", "htmx-*", "theme-*", "palette-*"]
preserve_attributes = ["data-theme", "data-palette", "data-text-size"]
skip_files = ["vendor/*", "normalize.css"]
warning_threshold = 0
```

CSS purge removes unused rules by scanning generated HTML and keeping only selectors
that are actually present. The purge logic always preserves key @-rules and keeps
pseudo-only selectors like `:root` or `::selection` to avoid dropping base/theme styles.

### Tailwind (`[markata-go.tailwind]`)

markata-go can run the Tailwind standalone CLI automatically and wire the output
into your theme so you don't need a separate build step.

```toml
[markata-go.tailwind]
include = "css"                # "css", "js", or false (default: "css")
preflight = false               # Enable Tailwind Preflight reset styles (default: false)
input = "tailwind.css"          # Input CSS (relative to assets_dir)
output = "markata-tailwind.css" # Output CSS (relative to assets_dir)
config_file = ""                # Optional tailwind.config.js path
build = true                     # Run Tailwind CLI during build
minify = true                    # Pass --minify to Tailwind CLI
auto_install = true              # Auto-download Tailwind CLI (default: true)
version = "v3.4.19"             # Managed Tailwind CLI version tag
cache_dir = ""                  # Cache dir for Tailwind CLI
binary = ""                     # Optional path to tailwindcss binary
extra_args = []                  # Optional extra CLI arguments
verbose = false                  # Verbose installer/build logs
```

**How it works:**

- `input`/`output` resolve relative to `assets_dir` and are copied/fingerprinted
  like any other static asset.
- If `input` does not exist, markata-go generates a default Tailwind entry file for
  you, so basic utility usage works without creating `tailwind.css` first.
- If `extra_args` is empty and `config_file` is unset, markata-go generates a
  temporary Tailwind config that scans a generated token manifest built from
  rendered page HTML plus local JS/template sources. That generated config
  disables Tailwind Preflight by default so utilities do not override markata-go's
  built-in theme styles. Set `preflight = true` to opt back into the reset. The
  manifest is hashed so Tailwind is skipped when the effective utility set has not changed.
- `include = "css"` loads the compiled CSS. If `theme.custom_css` is unset, the
  Tailwind output is assigned automatically. Explicit `theme.custom_css` always wins.
- `include = "js"` injects `https://cdn.tailwindcss.com` in the document head.
  When `[markata-go.assets].mode = "self-hosted"`, markata-go vendors this script
  and uses the local URL from `asset_urls`.
- `include = false` disables auto-inclusion; the build can still run.
- Fast mode skips Tailwind rebuilds when the compiled asset already exists.

By default markata-go uses its managed Tailwind CLI version instead of a global
`tailwindcss` on your `PATH`, which keeps builds consistent across machines. Set
`auto_install = false` or `binary = "/path/to/tailwindcss"` only when you want to
override that behavior.

**Validation warning:** If `include = "css"` and `[markata-go.css_purge].enabled = false`,
markata-go emits a warning because Tailwind outputs an intentionally large CSS bundle.

For Tailwind-first sites with a custom reset strategy, either set `preflight = true`
or provide your own `config_file`. User-supplied Tailwind config files are respected as-is.

### Aesthetic Settings (`[markata-go]`)

Aesthetics control the non-color visual aspects of your site: border radius, spacing, shadows, and borders. While palettes control *what colors* you use, aesthetics control *how things are shaped*.

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `aesthetic` | string | `"balanced"` | Aesthetic preset to use |

```toml
[markata-go]
aesthetic = "elevated"
```

**Available aesthetics:**

| Aesthetic | Description |
|-----------|-------------|
| `balanced` | Default. Comfortable rounding, subtle shadows, normal spacing |
| `brutal` | Sharp corners, thick borders, tight spacing, no shadows |
| `minimal` | No rounding, maximum whitespace, no shadows, hairline borders |
| `elevated` | Generous rounding, layered shadows, generous spacing |
| `precision` | Subtle corners, compact spacing, hairline borders, minimal shadows |

#### Aesthetic Overrides (`[markata-go.aesthetic_overrides]`)

Fine-tune individual aesthetic values without creating a full custom aesthetic:

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `border_radius` | string | varies | CSS border-radius value (e.g., `"8px"`, `"0"`) |
| `spacing_scale` | float | varies | Spacing multiplier (e.g., `1.0`, `1.5`) |
| `border_width` | string | varies | CSS border-width value (e.g., `"1px"`, `"3px"`) |
| `border_style` | string | varies | CSS border-style value (e.g., `"solid"`, `"none"`) |
| `shadow_intensity` | float | varies | Shadow opacity multiplier (0.0-1.0) |
| `shadow_size` | string | varies | Shadow size: `"none"`, `"sm"`, `"md"`, `"lg"`, `"xl"` |

```toml
[markata-go]
aesthetic = "balanced"

# Override specific values
[markata-go.aesthetic_overrides]
border_radius = "16px"
shadow_intensity = 0.8
spacing_scale = 1.2
```

**Combining with palettes:**

```toml
[markata-go]
aesthetic = "elevated"

[markata-go.theme]
palette = "catppuccin-mocha"
```

See the [[themes-and-styling|Themes Guide]] for detailed aesthetic customization and the full list of built-in presets.

### Layout Components (`[markata-go.components]`)

markata-go includes a configurable layout components system for navigation, footer, and sidebars. Each component can be enabled/disabled and positioned via configuration.

#### Navigation Component (`[markata-go.components.nav]`)

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `enabled` | bool | `true` | Show navigation |
| `position` | string | `"header"` | Position: `"header"` or `"sidebar"` |
| `style` | string | `"horizontal"` | Style: `"horizontal"` or `"vertical"` |
| `items` | array | `[]` | Navigation links (overrides top-level nav) |

```toml
[markata-go.components.nav]
enabled = true
position = "header"
style = "horizontal"

[[markata-go.components.nav.items]]
label = "Home"
url = "/"

[[markata-go.components.nav.items]]
label = "Blog"
url = "/blog/"

[[markata-go.components.nav.items]]
label = "GitHub"
url = "https://github.com/user"
external = true
```

#### Footer Component (`[markata-go.components.footer]`)

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `enabled` | bool | `true` | Show footer |
| `text` | string | `""` | Custom footer text |
| `show_copyright` | bool | `true` | Show copyright line |
| `links` | array | `[]` | Footer links |

```toml
[markata-go.components.footer]
enabled = true
text = "Thanks for reading!"
show_copyright = true

[[markata-go.components.footer.links]]
label = "RSS"
url = "/rss.xml"

[[markata-go.components.footer.links]]
label = "GitHub"
url = "https://github.com/user"
external = true
```

#### Document Sidebar (`[markata-go.components.doc_sidebar]`)

Table of contents sidebar for long-form content.

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `enabled` | bool | `false` | Show TOC sidebar |
| `position` | string | `"right"` | Position: `"left"` or `"right"` |
| `width` | string | `"250px"` | Sidebar width |
| `min_depth` | int | `2` | Minimum heading level (h2 = 2) |
| `max_depth` | int | `4` | Maximum heading level (h4 = 4) |

```toml
[markata-go.components.doc_sidebar]
enabled = true
position = "right"
width = "280px"
min_depth = 2
max_depth = 4
```

To show the TOC only on blog-layout posts (pages rendered with `post.html`), use `show_toc` under `[markata-go.layout.blog]` instead. When `doc_sidebar.enabled = true`, it takes precedence and its `position` and `width` are used. See the Blog Layout section below.

#### Feed Sidebar (`[markata-go.components.feed_sidebar]`)

Series/collection navigation sidebar for posts in the same feed.

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `enabled` | bool | `false` | Show feed sidebar |
| `position` | string | `"left"` | Position: `"left"` or `"right"` |
| `width` | string | `"250px"` | Sidebar width |
| `title` | string | `""` | Custom title (defaults to feed title) |
| `feeds` | array | `[]` | Feed slugs to show navigation for |
| `max_posts` | integer | `0` | Initial server-rendered sidebar window; public HTML feeds can progressively load the rest from static chunks |

```toml
[markata-go.components.feed_sidebar]
enabled = true
position = "left"
width = "250px"
title = "In this series"
feeds = ["tutorials", "guides"]
max_posts = 51
```

#### Sidebar behavior

The document and feed sidebars are drawers at every width. Each drawer has a
small handle ("Series" on the left, "On this page" on the right). Clicking the
handle, or pressing `b` (left) / `Shift+B` (right), opens the drawer; doing it
again closes it. Drawers never open on hover, so they cannot pop in and out
while you read. Hovering or focusing a handle shows its hotkey (`b` or `⇧B`),
and the handle's tooltip names it too.

**Wide screens (`1201px` and up):** an open drawer docks beside the article
and pushes the content aside instead of covering it. With both drawers open
the article sits between them, and it narrows if the viewport cannot fit all
three at full width. The open/closed choice is remembered per side in browser
storage and restored on later pages without animating. While the site header
is on screen the drawer starts just below it; once you scroll past the header
the drawer fills the full viewport height. Each drawer has a single themed
scrollbar.

**Narrow screens (`1200px` and below):** the handles become small pills in the
bottom corners. An open drawer slides over the page with a dimmed backdrop,
and only one drawer is open at a time. Tapping the backdrop, pressing
`Escape`, following a link inside the drawer, or tapping the pill again closes
it, and focus returns to the handle. Overlay drawers always start closed and
don't change the remembered wide-screen state. A content sidebar
(`[markata-go.content_sidebar]`) stacks below the article at this width.

The feed drawer also carries the feed picker, a `current/total` counter, and
previous/next feed buttons (`{` / `}`) for switching between primary feeds.

TOC entries are built from the heading's visible text: resolved wikilinks,
Markdown links, emphasis and inline code are stripped, so the sidebar never
shows raw HTML or markup. TOC anchors always match the rendered heading ids,
including duplicate headings (`setup`, `setup-1`, `setup-2`), headings with
inline code such as `` `update_meta` ``, headings with smart quotes, and
explicit `{#custom-id}` attributes. Lines starting with `#` inside fenced code
blocks are never treated as headings.

Clicking a TOC entry moves quickly to the heading with a short glide (under a
quarter of a second, starting close to the target so long jumps don't
stream the whole page past). The URL hash updates and the heading briefly
highlights. Scrolling, touching or pressing a key during the glide hands
control back to you immediately. With `prefers-reduced-motion`, the jump is
instant. As you read, a marker slides along the TOC rail to the current
section, and the TOC keeps that entry in view.

**Header on narrow screens:** below `1201px` the header switches to a compact
layout (title and controls on one row, search below, then a single scrollable
nav strip).

#### Share Component (`[markata-go.components.share]`)

Add a "Share this post" grid to the end of every article so readers can send your content to social platforms or copy the link.

Post pages also include a separate `Copy this post` control near the article title. That top-of-post menu exposes `URL`, `Rich`, `Markdown`, and `Text` clipboard options based on the current post output. To remove the menu without removing the share grid or disabling post formats, set:

```toml
[markata-go.components.post_copy]
enabled = false
```

`components.post_copy.enabled` defaults to `true`. The setting applies to post pages, including custom post templates that include the built-in copy component.

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `enabled` | bool | `true` | Show the component. |
| `title` | string | `"Share this post"` | Heading text above the buttons. |
| `position` | string | `"bottom"` | Modifier used by CSS (applied as `share-panel--<position>`). |
| `platforms` | array | `["twitter","bluesky","linkedin","whatsapp","facebook","telegram","pinterest","reddit","hacker_news","email","copy"]` | Ordered list of platform keys to render. Omit entries to hide them. |
| `custom` | table | `{}` | Custom platform definitions (`name`, `icon`, `url`). Icon paths without `/`, `http`, or `data:` are resolved via `theme_asset_hashed` (e.g., `icons/share/mastodon.svg`). |

```toml
[markata-go.components.share]
enabled = true
title = "Share this post"
  platforms = ["twitter", "mastodon", "linkedin", "copy"]

[markata-go.components.share.custom]
mastodon = { name = "Mastodon", icon = "mastodon.svg", url = "https://mastodon.social/share?text={{title}}&url={{url}}" }
```

Supported placeholders for share URLs:

- `{{title}}` → URL-encoded post title (falls back to site title or slug).
- `{{url}}` → URL-encoded absolute post URL (`config.url` + `post.href`).
- `{{excerpt}}` → URL-encoded post description or excerpt when provided.

Copying uses the Clipboard API with a DOM fallback, updates the button label/tooltip to "Link copied" for 2 seconds, and shows a small desktop-only toast. Built-in platforms include Twitter, Bluesky, LinkedIn, WhatsApp, Facebook, Telegram, Pinterest, Reddit, Hacker News, Email, and Copy Link, each with a coordinated icon from `pkg/themes/default/static/icons/share/`.

#### Post Connections Component (`[markata-go.components.post_connections]`)

Configure backlinks/outlinks UI on post pages. You can render a graph, a list, or both, and tune thresholds independently.

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `enabled` | bool | `true` | Enable post connections UI. |
| `display` | array | `['graph']` | Modes to render: `graph`, `list`, or both. |
| `min_links` | int | `1` | Lower bound (inlinks + outlinks) to render any mode. |
| `max_links` | int | `0` | Upper bound to render any mode (`0` = no max). |
| `graph_min_links` | int | `3` | Lower bound for graph mode. |
| `graph_max_links` | int | `0` | Upper bound for graph mode (`0` = no max). |
| `list_min_links` | int | `1` | Lower bound for list mode. |
| `list_max_links` | int | `0` | Upper bound for list mode (`0` = no max). |
| `inlinks_limit` | int | `8` | Max inlinks shown in list mode. |
| `outlinks_limit` | int | `8` | Max outlinks shown in list mode. |

The graph loader is only emitted on pages that render the graph preview, so pages that stay in list-only mode avoid the graph JavaScript entirely.

```toml
[markata-go.components.post_connections]
enabled = true
display = ["list", "graph"]

# Show list first, then graph for denser pages
list_min_links = 1
list_max_links = 5
graph_min_links = 6
graph_max_links = 0

inlinks_limit = 6
outlinks_limit = 10
```

### IndieAuth Settings (`[markata-go.indieauth]`)

[IndieAuth](https://indieauth.net/) is a decentralized identity protocol that allows you to use your own domain to sign in to websites. markata-go can add the necessary `<link>` tags to your site's HTML head.

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `enabled` | bool | `false` | Enable IndieAuth link tags in HTML head |
| `authorization_endpoint` | string | `""` | URL of your authorization endpoint |
| `token_endpoint` | string | `""` | URL of your token endpoint |
| `me_url` | string | `""` | Your profile URL for `rel="me"` links |

```toml
[markata-go.indieauth]
enabled = true
authorization_endpoint = "https://indieauth.com/auth"
token_endpoint = "https://tokens.indieauth.com/token"
me_url = "https://github.com/yourusername"
```

When enabled, this adds the following link tags to your site's `<head>`:

```html
<link rel="authorization_endpoint" href="https://indieauth.com/auth">
<link rel="token_endpoint" href="https://tokens.indieauth.com/token">
<link rel="me" href="https://github.com/yourusername">
```

### Webmention Settings (`[markata-go.webmention]`)

[Webmention](https://www.w3.org/TR/webmention/) is a web standard for conversations and interactions across websites. markata-go can add the webmention endpoint link tag to your site.

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `enabled` | bool | `false` | Enable Webmention link tag in HTML head |
| `endpoint` | string | `""` | URL of your Webmention endpoint |

```toml
[markata-go.webmention]
enabled = true
endpoint = "https://webmention.io/example.com/webmention"
```

When enabled, this adds the following link tag to your site's `<head>`:

```html
<link rel="webmention" href="https://webmention.io/example.com/webmention">
```

**Popular Webmention services:**
- [webmention.io](https://webmention.io/) - Free hosted webmention service
- [Bridgy](https://brid.gy/) - Connects social media interactions to webmentions

### WebSub Settings (`[markata-go.websub]`)

[WebSub](https://www.w3.org/TR/websub/) enables near-real-time feed delivery by advertising hub URLs. When enabled, markata-go adds discovery links to HTML pages and RSS/Atom feeds.

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `enabled` | bool | `false` | Enable WebSub discovery links |
| `hubs` | array | `[]` | List of WebSub hub URLs |

```toml
[markata-go.websub]
enabled = true
hubs = ["https://hub.example.com/"]
```

**Output:**
- HTML `<head>` includes `<link rel="hub" href="...">` for each hub
- RSS/Atom feeds include `rel="hub"` and `rel="self"` links

### Head Configuration (`[markata-go.head]`)

The head configuration allows you to customize elements in the HTML `<head>` section, including custom meta tags, links, scripts, and feed alternate links.

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `text` | string | `""` | Raw HTML to include in head (use with caution) |
| `meta` | array | `[]` | Custom meta tags |
| `link` | array | `[]` | Custom link tags |
| `script` | array | `[]` | Custom script tags |
| `alternate_feeds` | array | `[]` | Feed alternate links (RSS, Atom, JSON) |

#### Alternate Feed Links

By default, markata-go includes RSS and Atom feed links in the `<head>`. You can customize which feeds are advertised:

```toml
[markata-go.head]
# Customize which feeds get <link rel="alternate"> tags
[[markata-go.head.alternate_feeds]]
type = "rss"
title = "RSS Feed"
href = "/rss.xml"

[[markata-go.head.alternate_feeds]]
type = "atom"
title = "Atom Feed"
href = "/atom.xml"

[[markata-go.head.alternate_feeds]]
type = "json"
title = "JSON Feed"
href = "/feed.json"
```

**Supported feed types:**
| Type | MIME Type | Description |
|------|-----------|-------------|
| `rss` | `application/rss+xml` | RSS 2.0 feed |
| `atom` | `application/atom+xml` | Atom 1.0 feed |
| `json` | `application/feed+json` | JSON Feed |

**Example: Only advertise JSON Feed:**

```toml
[[markata-go.head.alternate_feeds]]
type = "json"
title = "JSON Feed"
href = "/feed.json"
```

**Example: Feed per section:**

```toml
[[markata-go.head.alternate_feeds]]
type = "rss"
title = "Blog RSS"
href = "/blog/rss.xml"

[[markata-go.head.alternate_feeds]]
type = "rss"
title = "Tutorials RSS"
href = "/tutorials/rss.xml"
```

#### Custom Meta Tags

```toml
[[markata-go.head.meta]]
name = "author"
content = "Jane Doe"

[[markata-go.head.meta]]
property = "og:site_name"
content = "My Site"
```

#### Custom Link Tags

```toml
[[markata-go.head.link]]
rel = "icon"
href = "/favicon.ico"

[[markata-go.head.link]]
rel = "preconnect"
href = "https://fonts.googleapis.com"
crossorigin = true
```

#### Custom Script Tags

```toml
[[markata-go.head.script]]
src = "/js/analytics.js"
```

### Glob Settings (`[markata-go.glob]`)

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `patterns` | string[] | `["pages/**/*.md", "posts/**/*.md"]` | Glob patterns to find content files |
| `use_gitignore` | bool | `true` | Respect .gitignore when finding files |
| `slug_mode` | string | `"flat"` | Slug derivation mode: `flat` or `path` |

```toml
[markata-go.glob]
patterns = ["posts/**/*.md", "pages/*.md", "docs/**/*.md"]
use_gitignore = true
slug_mode = "flat"
```

`slug_mode` controls how markata derives a slug when frontmatter does not explicitly set `slug`:

- `flat` keeps the current default behavior and uses the filename only. `posts/2026/hello-world.md` becomes `hello-world`.
- `path` derives the slug from the content path, strips leading `posts/` or `pages/`, and treats `index.md` or `README.md` as the directory root. `posts/notes/today.md` becomes `notes/today`, and `pages/docs/README.md` becomes `docs`.

Frontmatter `slug` still overrides either mode.

You can also mix both behaviors by directory with `slug_rules`:

```toml
[markata-go.glob]
slug_mode = "flat"

[[markata-go.glob.slug_rules]]
prefix = "posts/blog"
mode = "flat"

[[markata-go.glob.slug_rules]]
prefix = "posts/notes"
mode = "path"
```

Rules match by content path prefix, and the longest matching prefix wins. That lets you keep primary content flat while letting notes or docs stay nested.

### Markdown Settings (`[markata-go.markdown]`)

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `extensions` | string[] | `[]` | Markdown extensions to enable |

Available extensions:
- `tables` - GFM tables
- `strikethrough` - ~~strikethrough~~ text
- `autolinks` - Automatic URL linking
- `tasklist` - Task list checkboxes

```toml
[markata-go.markdown]
extensions = ["tables", "strikethrough", "autolinks", "tasklist"]
```

#### Syntax Highlighting (`[markata-go.markdown.highlight]`)

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `enabled` | bool | `true` | Enable syntax highlighting |
| `theme` | string | `""` | Chroma theme (empty = auto from palette) |
| `line_numbers` | bool | `false` | Show line numbers in code blocks |

```toml
[markata-go.markdown.highlight]
enabled = true
theme = "github-dark"    # Or leave empty for auto-detection
line_numbers = false
```

### Template Selection

Named presets and global per-format fallbacks are loaded from the
`[markata-go]` namespace and work in TOML, YAML, and JSON configuration files.

```toml
[markata-go.template_presets.blog]
html = "blog.html"
txt = "blog.txt"
markdown = "blog.md"
og = "blog-og.html"

[markata-go.default_templates]
html = "post.html"
txt = "default.txt"
```

When a post selects `template = "blog"`, the matching preset is used first.
`default_templates` is used after post, preset, and layout selection when no
earlier template supplies the requested format.

### Error Pages

Configure the generated 404 page with `[markata-go.error_pages]`:

```toml
[markata-go.error_pages]
enable_404 = true
custom_404_template = "404.html"
max_suggestions = 5
```

The built-in `404.html` uses the site's theme and offers a search field, the
requested address, and links back to the home page and archive. To customize its
layout, provide a `404.html` template in your `templates/` directory. It extends
the post page by default, so it can use the same navigation, palette, and theme
controls as the rest of the site.

### Resource Hints

Resource-hint settings use `[markata-go.resource_hints]`:

```toml
[markata-go.resource_hints]
enabled = true
auto_detect = true
exclude_domains = ["example.com"]
```

The feature is configured independently from plugin registration. Its default
settings enable hints and automatic detection; set `enabled = false` to keep
the typed setting disabled.

### Mermaid Settings (`[markata-go.mermaid]`)

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `enabled` | bool | `true` | Whether mermaid diagram processing is active |
| `mode` | string | `"client"` | Rendering mode: `client` (browser-side), `cli` (mmdc), or `chromium` (headless Chrome) |
| `cdn_url` | string | `/assets/vendor/mermaid/mermaid.min.js` | Mermaid JS URL (local by default) |
| `theme` | string | `"default"` | Mermaid theme: `default`, `dark`, `forest`, `neutral`. Ignored when `use_css_variables` is true. |
| `use_css_variables` | bool | `true` | Derive diagram colors from site CSS custom properties. Reads `--color-background`, `--color-text`, `--color-primary`, `--color-code-bg`, and `--color-surface` with hardcoded fallbacks. |
| `lightbox` | bool | `true` | Enable click-to-zoom lightbox overlay with interactive pan and zoom via svg-pan-zoom. |
| `lightbox_selector` | string | `".glightbox-mermaid"` | CSS selector for lightbox links (backward compatibility; programmatic API does not use it). |

#### CLI Mode Settings (`[markata-go.mermaid.cli]`)

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `mmdc_path` | string | `""` | Path to the mmdc binary. Auto-detected if empty. |
| `extra_args` | string | `""` | Additional command-line arguments passed to mmdc. |

#### Chromium Mode Settings (`[markata-go.mermaid.chromium]`)

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `browser_path` | string | `""` | Path to Chrome/Chromium binary. Auto-detected if empty. |
| `timeout` | int | `30` | Maximum seconds to wait for a diagram to render. |
| `max_concurrent` | int | `4` | Maximum number of parallel diagram renders. |
| `no_sandbox` | bool | `false` | Disable Chromium sandbox. Required in containers (Docker, Distrobox, etc.). |

```toml
[markata-go.mermaid]
enabled = true
theme = "default"
use_css_variables = true    # Diagrams automatically match your site palette
lightbox = true             # Click any diagram to zoom, pan, and explore
```

To pre-render diagrams to static SVGs using Chromium (no client-side JS needed):

```toml
[markata-go.mermaid]
mode = "chromium"

[markata-go.mermaid.chromium]
timeout = 30
max_concurrent = 4
# no_sandbox = true         # Required inside Docker or other containers
```

Pre-rendered Mermaid SVGs are cached automatically. Hot builds reuse cached SVGs
when the diagram source and rendering inputs match, avoiding repeated Chromium or
`mmdc` rendering for unchanged diagrams.

#### Chromium in Containers (Docker, Distrobox, Podman)

Chromium's sandbox requires kernel capabilities that most containers restrict.
Set `no_sandbox = true` when running inside Docker, Distrobox, Podman, or
similar environments:

```toml
[markata-go.mermaid.chromium]
no_sandbox = true
```

**Installing Chrome without sudo:** If your container does not have Chrome or
Chromium installed, download `chrome-headless-shell` from
[Chrome for Testing](https://googlechromelabs.github.io/chrome-for-testing/).
This is a lightweight headless binary (~180 MB) that requires no GUI
dependencies and no root access:

```bash
# Download and extract (adjust version as needed)
cd /tmp
curl -fsSL https://storage.googleapis.com/chrome-for-testing-public/LATEST_RELEASE_STABLE -o version.txt
VERSION=$(cat version.txt)
curl -fsSL "https://storage.googleapis.com/chrome-for-testing-public/${VERSION}/linux64/chrome-headless-shell-linux64.zip" -o headless.zip
unzip headless.zip

# Install to ~/.local (no sudo needed)
mkdir -p ~/.local/lib/chrome-headless-shell
mv chrome-headless-shell-linux64/* ~/.local/lib/chrome-headless-shell/

# Create a wrapper script so markata-go auto-detects it
mkdir -p ~/.local/bin
cat > ~/.local/bin/headless-shell << 'EOF'
#!/bin/sh
exec "$HOME/.local/lib/chrome-headless-shell/chrome-headless-shell" "$@"
EOF
chmod +x ~/.local/bin/headless-shell

# Verify
headless-shell --version
```

markata-go auto-detects `headless-shell` in your `PATH`. You can also set the
path explicitly:

```toml
[markata-go.mermaid.chromium]
browser_path = "/home/you/.local/lib/chrome-headless-shell/chrome-headless-shell"
no_sandbox = true
```

To disable the lightbox overlay (diagrams render inline only):

```toml
[markata-go.mermaid]
lightbox = false
```

To use a fixed Mermaid theme instead of CSS variables:

```toml
[markata-go.mermaid]
use_css_variables = false
theme = "dark"
```

See the [plugin reference](../reference/plugins.md#mermaid) for full details on lightbox interaction, CSS classes, and theming.

### Image Optimization Settings (`[markata-go.image_optimization]`)

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `enabled` | bool | `true` | Enable image optimization |
| `formats` | string[] | `"avif", "webp"` | Output formats to generate |
| `quality` | int | `80` | Default quality for all formats |
| `avif_quality` | int | `80` | AVIF quality override |
| `webp_quality` | int | `80` | WebP quality override |
| `widths` | int[] | `[480, 768, 1200]` | Responsive widths to generate |
| `sizes` | string | `"100vw"` | Sizes attribute for responsive sources |
| `cache_dir` | string | `.markata/image-cache` | Cache directory for encode metadata |
| `avifenc_path` | string | `""` | Path to `avifenc` (auto-detect if empty) |
| `cwebp_path` | string | `""` | Path to `cwebp` (auto-detect if empty) |

```toml
[markata-go.image_optimization]
enabled = true
formats = ["avif", "webp"]
quality = 80
widths = [480, 768, 1200]
sizes = "100vw"
```

The plugin generates local AVIF/WebP variants when an encoder is available. Trusted
Dropper image URLs receive responsive `srcset` candidates that use Dropper's width
parameter; the build does not download or generate those remote variants. Other
external URLs and data URIs remain unchanged. Missing local encoders do not fail
the build.

### Layout System (`[markata-go.layout]`)

The layout system controls page structure including sidebars, table of contents, headers, and footers. Different layouts can be assigned to different content types.

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `name` | string | `"blog"` | Default layout: `"docs"`, `"blog"`, `"landing"`, `"bare"` |
| `paths` | map | `{}` | Map URL prefixes to layouts |
| `feeds` | map | `{}` | Map feed slugs to layouts |

```toml
[markata-go.layout]
name = "docs"  # Default layout for all pages

# Path-based layout selection
[markata-go.layout.paths]
"/docs/" = "docs"
"/blog/" = "blog"
"/about/" = "landing"

# Feed-based layout selection
[markata-go.layout.feeds]
"docs" = "docs"
"blog" = "blog"
```

#### Documentation Layout (`[markata-go.layout.docs]`)

Three-panel layout with sidebar navigation, content, and table of contents.

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `sidebar_position` | string | `"left"` | Sidebar position: `"left"` or `"right"` |
| `sidebar_width` | string | `"280px"` | Sidebar width |
| `sidebar_collapsible` | bool | `true` | Allow sidebar collapse |
| `sidebar_default_open` | bool | `true` | Sidebar open by default |
| `toc_position` | string | `"right"` | TOC position: `"left"` or `"right"` |
| `toc_width` | string | `"220px"` | TOC width |
| `toc_collapsible` | bool | `true` | Allow TOC collapse |
| `toc_default_open` | bool | `true` | TOC open by default |
| `content_max_width` | string | `"800px"` | Maximum content width |
| `header_style` | string | `"minimal"` | Header: `"full"`, `"minimal"`, `"transparent"`, `"none"` |
| `footer_style` | string | `"minimal"` | Footer: `"full"`, `"minimal"`, `"none"` |

```toml
[markata-go.layout.docs]
sidebar_position = "left"
sidebar_width = "280px"
toc_position = "right"
toc_width = "220px"
content_max_width = "800px"
header_style = "minimal"
footer_style = "minimal"
```

#### Blog Layout (`[markata-go.layout.blog]`)

Single-column layout optimized for reading long-form content.

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `content_max_width` | string | `"720px"` | Maximum content width |
| `show_toc` | bool | `false` | Show the table of contents sidebar on blog-layout posts that have headings |
| `toc_position` | string | `"right"` | TOC position if enabled: `"left"` or `"right"` |
| `toc_width` | string | `"200px"` | TOC width |
| `header_style` | string | `"full"` | Header style |
| `footer_style` | string | `"full"` | Footer style |
| `show_author` | bool | `true` | Display post author |
| `show_date` | bool | `true` | Display post date |
| `show_tags` | bool | `true` | Display post tags |
| `show_reading_time` | bool | `true` | Display estimated reading time |
| `show_prev_next` | bool | `true` | Display prev/next navigation |

```toml
[markata-go.layout.blog]
content_max_width = "720px"
show_toc = true
show_author = true
show_date = true
show_tags = true
show_reading_time = true
show_prev_next = true
```

#### Landing Layout (`[markata-go.layout.landing]`)

Full-width layout for marketing pages and home pages.

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `content_max_width` | string | `"100%"` | Maximum content width |
| `header_style` | string | `"transparent"` | Header style |
| `header_sticky` | bool | `true` | Sticky header |
| `footer_style` | string | `"full"` | Footer style |
| `hero_enabled` | bool | `true` | Enable hero section |

```toml
[markata-go.layout.landing]
content_max_width = "100%"
header_style = "transparent"
header_sticky = true
hero_enabled = true
```

### Sidebar Configuration (`[markata-go.sidebar]`)

Configure the sidebar navigation component for documentation and guides.

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `enabled` | bool | `true` | Show sidebar |
| `position` | string | `"left"` | Position: `"left"` or `"right"` |
| `width` | string | `"280px"` | Sidebar width |
| `collapsible` | bool | `true` | Allow collapse |
| `default_open` | bool | `true` | Open by default |
| `title` | string | `""` | Sidebar title/header |
| `multi_feed` | bool | `false` | Multi-feed mode with sections |
| `feeds` | string[] | `[]` | Feed slugs for multi-feed mode |

```toml
[markata-go.sidebar]
enabled = true
position = "left"
width = "280px"
title = "Documentation"

# Manual navigation structure
[[markata-go.sidebar.nav]]
title = "Getting Started"
href = "/docs/"

[[markata-go.sidebar.nav]]
title = "Guides"
children = [
    { title = "Configuration", href = "/docs/guides/configuration/" },
    { title = "Templates", href = "/docs/guides/templates/" },
    { title = "Themes", href = "/docs/guides/themes/" },
    { title = "Feeds", href = "/docs/guides/feeds/" },
]

[[markata-go.sidebar.nav]]
title = "Reference"
children = [
    { title = "CLI", href = "/docs/reference/cli/" },
    { title = "Plugins", href = "/docs/reference/plugins/" },
]
```

#### Path-Specific Sidebars (`[markata-go.sidebar.paths]`)

Configure different sidebar content for different URL paths.

```toml
[markata-go.sidebar.paths."/docs/"]
title = "Documentation"
feed = "docs"  # Auto-generate from this feed

[markata-go.sidebar.paths."/guides/"]
title = "Guides"
[markata-go.sidebar.paths."/guides/".auto_generate]
directory = "guides"
order_by = "nav_order"
```

#### Multi-Feed Sidebars

Show posts from multiple feeds in collapsible sections.

```toml
[markata-go.sidebar]
multi_feed = true
feeds = ["docs", "guides", "tutorials"]

[[markata-go.sidebar.feed_sections]]
feed = "docs"
title = "Documentation"
collapsed = false

[[markata-go.sidebar.feed_sections]]
feed = "guides"
title = "Guides"
collapsed = true
max_items = 10
```

### Table of Contents (`[markata-go.toc]`)

Configure the table of contents component for document navigation.

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `enabled` | bool | `true` | Show TOC |
| `auto_enable` | bool | `false` | Auto-show based on thresholds |
| `min_toc_links` | int | `3` | Min headings for auto-show |
| `min_word_count` | int | `500` | Min words for auto-show |
| `position` | string | `"right"` | Position: `"left"` or `"right"` |
| `width` | string | `"220px"` | TOC width |
| `min_depth` | int | `2` | Minimum heading level (h2 = 2) |
| `max_depth` | int | `4` | Maximum heading level (h4 = 4) |
| `title` | string | `"On this page"` | TOC section title |
| `collapsible` | bool | `true` | Allow collapse |
| `default_open` | bool | `true` | Open by default |
| `scroll_spy` | bool | `true` | Highlight current section |

```toml
[markata-go.toc]
enabled = true
auto_enable = true          # Auto-show TOC for long posts
min_toc_links = 3           # Show if 3+ headings
min_word_count = 500        # Show if 500+ words
position = "right"
width = "220px"
title = "On this page"
min_depth = 2
max_depth = 4
scroll_spy = true
```

### Header Configuration (`[markata-go.header]`)

Configure the header component for layouts.

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `style` | string | `"full"` | Style: `"full"`, `"minimal"`, `"transparent"`, `"none"` |
| `sticky` | bool | `true` | Stick to top when scrolling |
| `show_logo` | bool | `true` | Display site logo |
| `show_title` | bool | `true` | Display site title |
| `show_nav` | bool | `true` | Display navigation links |
| `show_search` | bool | `true` | Display search box |
| `show_theme_toggle` | bool | `true` | Display dark/light toggle in header (legacy gate; also controlled by `theme.switcher.mode_toggle`) |

```toml
[markata-go.header]
style = "full"
sticky = true
show_logo = true
show_title = true
show_nav = true
show_search = true
show_theme_toggle = true
```

### Footer Configuration (`[markata-go.footer_layout]`)

Configure the footer component for layouts.

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `style` | string | `"full"` | Style: `"full"`, `"minimal"`, `"none"` |
| `sticky` | bool | `false` | Stick to bottom |
| `show_copyright` | bool | `true` | Display copyright notice |
| `copyright_text` | string | `""` | Custom copyright text |
| `show_social_links` | bool | `true` | Display social media links |
| `show_nav_links` | bool | `true` | Display navigation links |

```toml
[markata-go.footer_layout]
style = "full"
sticky = false
show_copyright = true
copyright_text = "© 2024 My Site. All rights reserved."
show_social_links = true
show_nav_links = true
```

### Link Avatars (`[markata-go.link_avatars]`)

Add favicon icons next to external links.

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `enabled` | bool | `false` | Enable link avatars |
| `mode` | string | `"js"` | `"js"`, `"local"`, or `"hosted"` |
| `selector` | string | `"a[href^='http']"` | CSS selector for links to enhance |
| `service` | string | `"duckduckgo"` | Favicon service (`duckduckgo`, `google`, `custom`) |
| `template` | string | `""` | Custom service template (used when `service = "custom"`) |
| `size` | int | `16` | Icon size in pixels |
| `position` | string | `"before"` | `"before"` or `"after"` |
| `ignore_domains` | string[] | `[]` | Domains to skip |
| `ignore_origins` | string[] | `[]` | Full origins to skip |
| `ignore_selectors` | string[] | `[]` | CSS selectors to skip |
| `ignore_classes` | string[] | `[]` | CSS classes to skip |
| `ignore_ids` | string[] | `[]` | Element IDs to skip |
| `hosted_base_url` | string | `""` | Base URL for hosted mode assets |

```toml
[markata-go.link_avatars]
enabled = true
mode = "local"
service = "duckduckgo"
size = 16
position = "before"
ignore_domains = ["localhost", "127.0.0.1"]
```

### Blogroll Configuration (`[markata-go.blogroll]`)

Configure the blogroll and RSS reader functionality to display feeds from blogs you follow.
Blogroll fetches are cached by fetch timestamp, and stale cached results are reused when
refresh attempts fail.

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `enabled` | bool | `false` | Enable blogroll plugin |
| `cache_dir` | string | `"cache/blogroll"` | Cache directory |
| `cache_duration` | string | `"24h"` | Cache TTL (Go duration format) |
| `timeout` | int | `30` | HTTP timeout in seconds |
| `concurrent_requests` | int | `5` | Max parallel feed fetches |
| `max_entries_per_feed` | int | `50` | Max entries per feed |
| `fallback_image_service` | string | `""` | Last-resort screenshot template for Reader previews after article and source images |

```toml
[markata-go.blogroll]
enabled = true
cache_dir = "cache/blogroll"
cache_duration = "24h"
timeout = 30
concurrent_requests = 5
max_entries_per_feed = 50

# Custom templates
[markata-go.blogroll.templates]
blogroll = "blogroll.html"
reader = "reader.html"
```

#### Adding Feeds (`[[markata-go.blogroll.feeds]]`)

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `url` | string | Yes | RSS or Atom feed URL |
| `title` | string | No | Display name (auto-fetched if not set) |
| `description` | string | No | Short description |
| `category` | string | No | Group feeds together |
| `tags` | string[] | No | Additional labels |
| `site_url` | string | No | Main website URL |
| `image_url` | string | No | Logo or icon URL |
| `active` | bool | No | Set `false` to disable |

```toml
# Technology blogs
[[markata-go.blogroll.feeds]]
url = "https://simonwillison.net/atom/everything/"
title = "Simon Willison"
description = "Creator of Datasette, Django co-creator"
category = "Technology"
tags = ["python", "ai", "llm"]

[[markata-go.blogroll.feeds]]
url = "https://jvns.ca/atom.xml"
title = "Julia Evans"
description = "Making hard things easy to understand"
category = "Technology"
tags = ["linux", "networking"]

# Design blogs
[[markata-go.blogroll.feeds]]
url = "https://css-tricks.com/feed/"
title = "CSS-Tricks"
category = "Design"
tags = ["css", "frontend"]
```

**Generated pages:**
- `/blogroll/` - Directory of all feeds grouped by category
- `/reader/` - Curated reader page with day grouping, source favicons, and columns

See the [Blogroll Guide](/docs/guides/blogroll/) for detailed configuration and customization.

### Post Output Formats (`[markata-go.post_formats]`)

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `html` | bool | `true` | Generate standard HTML pages |
| `markdown` | bool | `true` | Generate raw markdown source |
| `text` | bool | `true` | Generate plain terminal-friendly text output |
| `ansi` | bool | `false` | Generate ANSI-styled terminal output |
| `og` | bool | `true` | Generate OpenGraph card HTML |

```toml
[markata-go.post_formats]
html = true       # /slug/index.html (default)
markdown = true   # /slug.md (canonical) - enabled by default
text = true       # /slug.txt (canonical) - enabled by default
ansi = true       # /slug.ansi (canonical) - opt-in rich terminal output
og = true         # /slug/og/index.html (social card)
```

**Reversed Redirects for txt/md/ansi**: For `.txt`, `.md`, and `.ansi` formats, content is placed at the canonical URL (`/slug.txt`, `/slug.md`, `/slug.ansi`) with backwards-compatible redirects at `/slug/index.<ext>` and `/slug/index.<ext>/index.html`. This supports standard web txt files like `robots.txt`, `llms.txt`, and `humans.txt` while adding an explicit ANSI terminal view.

**Canonical URLs and Alternate Links**: When post formats are enabled, markata-go automatically adds:

- `<link rel="canonical">` pointing to the post's primary URL (for SEO)
- `<link rel="alternate">` for each enabled format:
  - Markdown: `type="text/markdown"` linking to `/slug.md`
  - ANSI terminal: `type="text/plain"` linking to `/slug.ansi`
  - OG Card: `type="text/html"` linking to `og/`

**Visible Format Links**: When alternate formats are enabled, posts and feeds display visible links allowing visitors to access content in their preferred format.

Feed sidebar variant links use the same canonical short URLs (`/slug.md`, `/slug.txt`, `/slug.ansi`) rather than nested `/slug/index.<ext>` paths.
If a post format is disabled, its sidebar link is omitted.

**Per-post frontmatter overrides**: A post can override any of these format flags in frontmatter with a `post_formats` object. Overrides merge with the site defaults, so omitted keys continue to inherit from `[markata-go.post_formats]`.

```yaml
---
title: "One-off terminal post"
post_formats:
  ansi: true
  og: false
---
```

OG card pages automatically include:
- `<link rel="canonical">` pointing back to the original post
- `<meta name="robots" content="noindex, nofollow">` to prevent search engine indexing

See the [[post-formats|Post Output Formats Guide]] for detailed usage including social image generation and content negotiation.

### LLM-Friendly Site Index (`[markata-go.llms_txt]`)

Full-site builds generate a concise `/llms.txt` Markdown index by default.
It lists the site overview, up to 30 public documentation/recent pages, and
a sitemap link if available. Private, draft, unpublished, skipped, and
`noindex` posts are excluded. Links prefer the generated `.md` version
of pages when enabled and fall back to HTML URLs otherwise.

```toml
[markata-go.llms_txt]
enabled = true
max_pages = 30
```

- `enabled`: defaults to `true`; set `false` to disable the generator.
- `max_pages`: defaults to `30`, clamped to `1-100` (invalid values
  use the default).
- A hand-authored `llms.md` (published with text format), project
  `static/llms.txt`, or unmanaged `output/llms.txt` takes precedence.
- Single-page builds never generate this site-wide index.
- `llms.txt` is a content guide, **not** a robots/training permission policy.

See [Standard Web Txt Files](/docs/guides/standard-txt-files/).

### Well-Known Files (`[markata-go.well_known]`)

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `enabled` | bool | `true` | Enable auto-generated `.well-known` endpoints |
| `auto_generate` | array | `[...]` | Entries to auto-generate from site metadata |
| `ssh_fingerprint` | string | `""` | SSH fingerprint for `/.well-known/sshfp` |
| `keybase_username` | string | `""` | Keybase username for `/.well-known/keybase.txt` |

```toml
[markata-go.well_known]
enabled = true
auto_generate = ["host-meta", "host-meta.json", "webfinger", "nodeinfo", "time", "links"]

# Optional entries requiring config
ssh_fingerprint = "SHA256:abcdef..."
keybase_username = "username"
```

**Generated files:**
- `/.well-known/host-meta`
- `/.well-known/host-meta.json`
- `/.well-known/webfinger`
- `/.well-known/nodeinfo` and `/nodeinfo/2.0`
- `/.well-known/time`
- `/.well-known/links` (grouped outbound links by target domain; enabled by default)
- `/.well-known/internal-links` (grouped internal links by target URL; enabled by default)
- `/external-links/` (HTML view of the same grouped outbound link data)
- `/internal-links/` (HTML view of grouped internal link destinations)
- `/.well-known/sshfp` (when `ssh_fingerprint` is set)
- `/.well-known/keybase.txt` (when `keybase_username` is set)

Templates live under `templates/well-known/` plus `templates/external-links.html` and `templates/internal-links.html` for the HTML page overrides.

### Content Templates (`[content_templates]`)

Content templates configure the `markata-go new` command, controlling default frontmatter and output directories for different content types.

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `directory` | string | `"content-templates"` | Directory for user-defined template files |
| `placement` | object | `{post:"pages/post",page:"pages",docs:"docs"}` | Map of template names to output directories |
| `templates` | array | `[]` | Custom template definitions |

```toml
[content_templates]
directory = "content-templates"

# Override default directory placement
[content_templates.placement]
post = "blog"          # markata-go new -t post creates in blog/
page = "pages"
docs = "documentation"

# Define custom templates
[[content_templates.templates]]
name = "tutorial"
directory = "tutorials"
body = "## Prerequisites\n\n## Steps\n\n## Summary"

[content_templates.templates.frontmatter]
templateKey = "tutorial"
series = ""

[[content_templates.templates]]
name = "recipe"
directory = "recipes"
body = "## Ingredients\n\n## Instructions"

[content_templates.templates.frontmatter]
templateKey = "recipe"
prep_time = ""
cook_time = ""
servings = 4
```

**File-based Templates:**

Create markdown files in the `content-templates/` directory (or your configured directory):

```markdown
---
# content-templates/changelog.md
templateKey: changelog
_directory: changelogs
version: ""
---

## Added

## Changed

## Fixed
```

The `_directory` field in frontmatter specifies the output directory and is removed from generated content.

**Usage:**

```bash
markata-go new --list                     # List all templates
markata-go new "My Post"                  # Use default (post) template
markata-go new "Tutorial" -t tutorial     # Use custom template
markata-go new "Recipe" -t recipe --dir custom  # Override directory
```

Built-in content templates are listed with aliases in `markata-go new --list`, such as `note (aka: ping, thought, status, tweet)`. Built-in aliases are also accepted directly with `--template`.

Generated files include frontmatter plus the configured body template, but they do not include a starter `# H1` heading. markata-go already renders the page title from frontmatter, so starter content should begin with body text or `##` sections.

See [[cli-reference|CLI Reference]] for complete `new` command documentation.

### Search Settings (`[search]`)

Site-wide search is enabled by default using [Pagefind](https://pagefind.app/). You can also use the built-in bleve search engine for a self-hosted API.

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `enabled` | bool | `true` | Enable/disable search |
| `backend` | string | `"pagefind"` | Search engine: `"pagefind"` or `"bleve"` |
| `endpoint` | string | `"/api/search"` | API endpoint path for bleve search |
| `position` | string | `"navbar"` | Where to show search: `navbar`, `sidebar`, `footer`, `custom` |
| `placeholder` | string | `"Search..."` | Search input placeholder text |
| `show_images` | bool | `true` | Show thumbnails in search results |
| `excerpt_length` | int | `200` | Characters for result excerpts |

```toml
[search]
enabled = true
backend = "pagefind"        # or "bleve" for self-hosted API
endpoint = "/api/search"    # API path (bleve only)
position = "navbar"
placeholder = "Search..."
show_images = true
excerpt_length = 200

# Pagefind CLI options
[search.pagefind]
bundle_dir = "_pagefind"    # Output directory for search index
root_selector = "main"       # CSS selector for searchable content
exclude_selectors = [".no-search", "nav", "footer"]  # Elements to exclude

# Bleve search API options
[search.bleve]
fuzzy = false               # Default fuzzy matching
limit = 20                  # Default result limit
max_limit = 100             # Maximum allowed limit
cors_origins = ["*"]        # Allowed CORS origins
```

**Privacy:** Private posts are searchable by public metadata such as title, description, tags, and date. Their body content and media are never included in the search index or exposed by the search API.

**Dev mode:** During `markata-go serve`, the bleve search API is automatically mounted at the configured endpoint. A dev-only script injects `window.__markataSearchEndpoint` so the frontend can use the local API — this is never included in production builds.

**Requirements:** Pagefind CLI must be installed. Install the standalone binary from [GitHub releases](https://github.com/Pagefind/pagefind/releases), or enable `auto_install = true` under `[search.pagefind]` to let markata-go download it automatically. If not installed, search is skipped with a warning.

See the [[search|Search Guide]] for detailed usage and customization.

### Feed Defaults (`[markata-go.feed_defaults]`)

Set `views = ["default", "simple", "calendar"]` to choose the available HTML
presentations. All three are enabled by default. Per-feed `views` replaces this
list; omitted values inherit it. Every explicit list must include `default`,
and an empty list or unknown name is invalid. Disabling `simple` suppresses its
HTML output even when `formats.simple_html` is true; disabling `calendar` removes
its controls and assets. Syndication formats are independent.
See [Feed Views](./feed-views.md) for examples.

The implicit [Pins feed](./pins.md) uses 100 items per page and manual pagination.
To change those values, define a feed with `slug = "pins"`; its configuration
takes precedence over the built-in Pins definition.

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `items_per_page` | int | `10` | Default items per page |
| `orphan_threshold` | int | `3` | Minimum items for a separate page |

```toml
[markata-go.feed_defaults]
items_per_page = 10
orphan_threshold = 3

[markata-go.feed_defaults.formats]
html = true
simple_html = true
rss = true
atom = true
json = true
sitemap = true
markdown = false
text = false

[markata-go.feed_defaults.templates]
html = "feed.html"
simple_html = "simple-feed.html"
rss = "feed.xml"
atom = "atom.xml"
json = "feed.json"
card = "card.html"
sitemap = "sitemap.xml"

[markata-go.feed_defaults.syndication]
max_items = 20
include_content = true
site_archive_disabled = false
feed_archives_disabled = false

[markata-go.feeds_page]
enabled = true
title = "Feeds"
description = "Browse the public feeds available on this site."
template = "feeds.html"
slug_prefix = "feeds"
```

### Feed Format Options

| Format | Output Path | Description |
|--------|-------------|-------------|
| `html` | `/{slug}/index.html` | Paginated HTML pages |
| `simple_html` | `/{slug}/simple/index.html` | Compact HTML list |
| `rss` | `/{slug}/rss.xml` | RSS 2.0 feed |
| `atom` | `/{slug}/atom.xml` | Atom 1.0 feed |
| `json` | `/{slug}/feed.json` | JSON Feed |
| `sitemap` | `/{slug}/sitemap.xml` | Sitemap XML |
| `markdown` | `/{slug}.md` | Markdown output |
| `text` | `/{slug}.txt` | Plain text output |

### Feed Configuration (`[[markata-go.feeds]]`)

Each feed is defined as an array item. The default config ships with an archive feed at `/archive/`.

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `slug` | string | `""` | URL-safe identifier (empty = root index) |
| `title` | string | `""` | Feed title |
| `description` | string | `""` | Feed description |
| `filter` | string | `""` | Filter expression for selecting posts |
| `sort` | string | `""` | Field to sort posts by |
| `reverse` | bool | `false` | Reverse sort order |
| `items_per_page` | int | inherited | Items per page (inherits from defaults) |
| `orphan_threshold` | int | inherited | Orphan threshold (inherits from defaults) |
| `limit` | int | `0` | Hard cap on total items (0 = unlimited) |
| `offset` | int | `0` | Skip the first N items (0 = none) |
| `formats` | object | inherited | Output formats (inherits from defaults) |
| `templates` | object | inherited | Templates (inherits from defaults) |
| `archive_disabled` | bool | `false` | Disable `/{slug}/archive/*` syndication endpoints for this feed |

```toml
# Main blog feed
[[markata-go.feeds]]
slug = "blog"
title = "Blog"
description = "All blog posts"
filter = "published == True"
sort = "date"
reverse = true
items_per_page = 10

[markata-go.feeds.formats]
html = true
rss = true
atom = true
json = true

# Home page feed (empty slug = root index.html)
[[markata-go.feeds]]
slug = ""
title = "Latest Posts"
filter = "published == True"
sort = "date"
reverse = true
items_per_page = 5

[markata-go.feeds.formats]
html = true

# Featured posts feed
[[markata-go.feeds]]
slug = "featured"
title = "Featured"
description = "Featured posts"
filter = "published == True and featured == True"
sort = "date"
reverse = true
```

### Series Settings (`[markata-go.series]`)

Series automatically generate ordered feed pages and navigation from `series` frontmatter.

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `slug_prefix` | string | `"series"` | URL prefix for series feeds (e.g., `/series/<name>/`) |
| `auto_sidebar` | bool | `true` | Auto-enable the series sidebar on series posts |
| `defaults` | object | see below | Default settings for all series feeds |
| `overrides` | map[string]object | `{}` | Per-series overrides (keyed by series slug or name) |

```toml
[markata-go.series]
slug_prefix = "series"
auto_sidebar = true

[markata-go.series.defaults]
items_per_page = 0
sidebar = true

[markata-go.series.defaults.formats]
html = true
simple_html = true
rss = true
atom = true
json = true
sitemap = true

[markata-go.series.overrides."building-a-rest-api"]
title = "Building a REST API with Go"
description = "A complete guide from zero to production"

[markata-go.series.overrides."building-a-rest-api".formats]
markdown = true
```

## Environment Variable Overrides

All configuration options can be overridden via environment variables using the `MARKATA_GO_` prefix.

### Variable Naming Convention

Environment variables follow the pattern: `MARKATA_GO_{SECTION}_{KEY}`

- Keys are UPPERCASE
- Underscores separate nested keys
- Section names use underscores

### Examples

```bash
# Core settings
export MARKATA_GO_OUTPUT_DIR=dist
export MARKATA_GO_URL=https://example.com
export MARKATA_GO_TITLE="My Site"
export MARKATA_GO_CONCURRENCY=4

# Glob settings
export MARKATA_GO_GLOB_PATTERNS="posts/**/*.md,pages/*.md"
export MARKATA_GO_GLOB_USE_GITIGNORE=true

# Markdown settings
export MARKATA_GO_MARKDOWN_EXTENSIONS="tables,strikethrough"

# Feed defaults
export MARKATA_GO_FEED_DEFAULTS_ITEMS_PER_PAGE=20
export MARKATA_GO_FEEDS_DEFAULTS_ORPHAN_THRESHOLD=5
export MARKATA_GO_FEED_DEFAULTS_FORMATS_HTML=true
export MARKATA_GO_FEED_DEFAULTS_FORMATS_RSS=true
export MARKATA_GO_FEED_DEFAULTS_SYNDICATION_MAX_ITEMS=50
export MARKATA_GO_FEED_DEFAULTS_SYNDICATION_INCLUDE_CONTENT=false

# Blogroll settings
export MARKATA_GO_BLOGROLL_ENABLED=false
```

### Value Formats

| Type | Format | Examples |
|------|--------|----------|
| String | Plain text | `MARKATA_GO_TITLE="My Site"` |
| Integer | Numeric | `MARKATA_GO_CONCURRENCY=4` |
| Boolean | `true`/`false`/`1`/`0`/`yes`/`no` | `MARKATA_GO_GLOB_USE_GITIGNORE=true` |
| List | Comma-separated | `MARKATA_GO_HOOKS="glob,load,render"` |

### Build with Environment Overrides

```bash
# Override output directory for a single build
MARKATA_GO_OUTPUT_DIR=dist markata-go build

# Build for staging environment
MARKATA_GO_URL=https://staging.example.com markata-go build

# Increase concurrency
MARKATA_GO_CONCURRENCY=8 markata-go build
```

## Multiple Config Files (Config Merging)

You can merge multiple configuration files together using the `--merge-config` (or `-m`) flag. This is useful for:

- **Environment-specific overrides** (e.g., `markata-go.local.toml`)
- **Fast build configurations** (e.g., `fast-markata-go.toml`)
- **Temporary overrides** without modifying your main config

### How Config Merging Works

Configs are merged in order, with later configs taking precedence:

1. **Defaults** - Built-in default values
2. **Base config** - Your main config file (`markata-go.toml`)
3. **Merge configs** - Each `--merge-config` file, in order
4. **Environment variables** - `MARKATA_GO_*` vars (highest precedence)

### Example: Fast Build Config

Create a `fast-markata-go.toml` for quick development builds:

```toml
[markata-go]
# Only process specific files for fast builds
glob_patterns = ["posts/draft-*.md"]

[markata-go.blogroll]
# Disable blogroll for faster builds
enabled = false
```

Use it with:

```bash
# Merge fast config with base config
markata-go build --merge-config fast-markata-go.toml

# Or use the short flag
markata-go build -m fast-markata-go.toml
```

### Example: Multiple Merge Configs

You can specify multiple merge configs. They are applied in order:

```bash
# Base config + local overrides + debug settings
markata-go build -m markata-go.local.toml -m debug.toml
```

### Example: No Base Config

You can use merge configs without a base config (uses defaults as base):

```bash
# Just use defaults + your override file
markata-go build -c "" -m fast-markata-go.toml
```

### Common Merge Config Patterns

**Local Development** (`markata-go.local.toml` - add to `.gitignore`):

```toml
[markata-go]
url = "http://localhost:8080"
concurrency = 2

[markata-go.blogroll]
enabled = false
```

**Fast Build** (`fast.toml`):

```toml
[markata-go]
# Limit to specific files
glob_patterns = ["posts/current/**/*.md"]

[markata-go.blogroll]
enabled = false
```

**Note on Boolean Fields**: Merging `enabled = false` into a base config with `enabled = true` requires environment variables. Use `MARKATA_GO_BLOGROLL_ENABLED=false` instead.

## Configuration CLI Commands

### `config show`

Display the resolved configuration:

```bash
# Same as `markata-go config show`
markata-go config

# Show as YAML (default)
markata-go config show

# Show as JSON
markata-go config show --json

# Show as TOML
markata-go config show --toml

# Include merged override files
markata-go config show -m fast-markata-go.toml
```

`config show` uses the same config resolution path as `build` and `serve`, including any `--merge-config` overrides.

YAML output includes source comments by default. Comments identify whether an
effective value came from a config file (with its path and line), an environment
variable, a CLI override, or a default. Use `--no-annotate` for plain YAML;
JSON and TOML output remain unannotated and parseable. In an interactive
terminal, YAML is syntax highlighted using the configured palette.

Default and file-backed comments include a copyable `markata-go config set
<key> <value>` command so you can override a setting without locating its
configuration file manually.

Conflicting format requests such as `markata-go config show --json --toml` fail with usage exit code `2`.

### `config get`

Get a specific configuration value:

```bash
# Get top-level value
markata-go config get output_dir

# Get nested value with dot notation
markata-go config get glob.patterns
markata-go config get feed_defaults.items_per_page
markata-go config get feed_defaults.formats.html
```

### `config set`

Update a specific configuration value in-place:

```bash
# Set top-level value
markata-go config set url "https://example.com"

# Set nested value
markata-go config set glob.patterns '["posts/**/*.md", "pages/**/*.md"]'
```

`config set` preserves file permissions for TOML, YAML, and JSON files. In CGO-enabled builds, TOML/YAML edits use tree-sitter byte ranges so existing comments, ordering, and formatting stay intact.

### `config validate`

Validate the configuration file:

```bash
# Validate default config file
markata-go config validate

# Validate specific config file
markata-go config validate -c custom.toml

# Validate with merged overrides
markata-go config validate -m fast-markata-go.toml
```

`config validate` uses the same merged configuration resolution as `build` and `serve`.

Example output:

```
✓ Configuration is valid
```

Or with errors:

```
✗ Configuration errors:
  - url: URL must include a scheme (e.g., https://)
  - concurrency: must be >= 0 (0 means auto-detect)
  - glob.patterns: no glob patterns specified, no files will be processed (warning)
```

### `config init`

Generate a starter configuration file:

```bash
# Create markata-go.toml with defaults
markata-go config init

# Create YAML config
markata-go config init site.yaml

# Overwrite existing file
markata-go config init --force
```

## Configuration Merging

markata-go merges configuration from multiple sources in order of increasing precedence:

```
┌─────────────────────────────────────────────────────────────────────┐
│                    CONFIGURATION RESOLUTION                          │
├─────────────────────────────────────────────────────────────────────┤
│  1. Built-in defaults                                                │
│  2. User config file (~/.config/markata-go/config.toml)             │
│  3. Local config file (./markata-go.toml)                           │
│  4. Environment variables (MARKATA_GO_*)                            │
│  5. CLI arguments (--output-dir, --config, etc.)                    │
│                                                                      │
│  Later sources OVERRIDE earlier sources                              │
│  Nested objects are DEEP MERGED                                     │
│  Arrays REPLACE (not append)                                        │
└─────────────────────────────────────────────────────────────────────┘
```

### Merge Behavior

**Scalar values:** Later wins

```toml
# User config: output_dir = "dist"
# Local config: output_dir = "public"
# Result: output_dir = "public"
```

**Objects:** Deep merge

```toml
# User config:
[markata-go.feed_defaults.formats]
html = true
rss = true

# Local config:
[markata-go.feed_defaults.formats]
atom = true

# Result:
# html = true   (from user)
# rss = true    (from user)
# atom = true   (from local)
```

**Arrays:** Replace (not append)

```toml
# User config: patterns = ["**/*.md"]
# Local config: patterns = ["posts/*.md", "pages/*.md"]
# Result: patterns = ["posts/*.md", "pages/*.md"]
```

## Default Values

When no configuration is provided, markata-go uses these defaults:

```toml
[markata-go]
output_dir = "output"
templates_dir = "templates"
assets_dir = "static"
hooks = ["default"]
concurrency = 0  # Auto-detect based on CPU cores

[markata-go.glob]
patterns = ["pages/**/*.md", "posts/**/*.md"]
use_gitignore = true

[markata-go.feed_defaults]
items_per_page = 10
orphan_threshold = 3

[markata-go.feed_defaults.formats]
html = true
rss = true
atom = false
json = false
markdown = false
text = false
```

## Authors Configuration

The authors system lets you define author profiles in your configuration and reference them from post frontmatter. This enables multi-author support where each post can have one or more authors with rich metadata.

### Basic Setup

```toml
[markata-go.authors]
generate_pages = false    # Generate individual author profile pages
feeds_enabled = false     # Generate per-author feeds

[markata-go.authors.authors.jane]
name = "Jane Doe"
role = "author"
url = "https://janedoe.com"
active = true
default = true            # Assigned to posts with no explicit author

[markata-go.authors.authors.bob]
name = "Bob Smith"
role = "editor"
active = true
```

### Author Fields

Each author entry supports these fields:

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `name` | string | **required** | Display name |
| `bio` | string | - | Short biography |
| `email` | string | - | Contact email (must contain `@` if set) |
| `avatar` | string | - | Avatar image URL or path |
| `url` | string | - | Personal website URL |
| `role` | string | - | Role label (e.g., "author", "editor") |
| `guest` | bool | `false` | Whether this is a guest author |
| `active` | bool | `false` | Whether the author is currently active |
| `default` | bool | `false` | Default author for posts without explicit authors |
| `contribution` | string | - | Free-text contribution description |
| `contributions` | []string | `[]` | CReDiT taxonomy roles |

### Referencing Authors in Posts

Use author IDs in post frontmatter:

```yaml
---
title: "Our Collaborative Post"
authors:
  - jane
  - bob
---
```

The legacy single-author field is also supported:

```yaml
---
title: "Solo Post"
author: "Jane Doe"
---
```

**Priority:** `authors` array takes precedence over `author` string. If neither is specified, the default author (the one with `default = true`) is assigned automatically.

### Author Configuration Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `generate_pages` | bool | `false` | Generate individual author profile pages |
| `url_pattern` | string | `""` | URL pattern for author pages (e.g., `/authors/{id}/`) |
| `feeds_enabled` | bool | `false` | Generate per-author feeds |

### Full Author Example

```toml
[markata-go.authors]
generate_pages = false
feeds_enabled = false

[markata-go.authors.authors.waylon]
name = "Waylon Walker"
bio = "Python and Go developer, blogger"
email = "waylon@example.com"
avatar = "/images/waylon.jpg"
url = "https://waylonwalker.com"
role = "author"
active = true
default = true
contributions = ["Writing", "Original Draft"]

[markata-go.authors.authors.guest]
name = "Guest Writer"
role = "editor"
guest = true
active = true
contribution = "Technical review"
```

**Related:** See the [Frontmatter Guide](frontmatter.md) for details on the `authors` and `author` fields.

## Encryption Settings

Encryption protects private post content with AES-256-GCM client-side encryption. It is **enabled by default**.

```toml
[encryption]
enabled = true                             # default: true
default_key = "default"                    # default: "default"
decryption_hint = "DM me for access"       # optional hint shown to visitors

[encryption.private_tags]
diary = "personal"                         # tag "diary" auto-encrypts with key "personal"
draft-ideas = "default"                    # tag "draft-ideas" encrypts with default key
```

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `enabled` | bool | `true` | Enable/disable encryption processing |
| `default_key` | string | `"default"` | Key name when post doesn't specify one |
| `decryption_hint` | string | `""` | Help text shown next to password prompt |
| `private_tags` | map | `{}` | Maps tag names to encryption key names |
| `enforce_strength` | bool | `true` | Require key passwords to meet the strength policy |
| `min_estimated_crack_time` | string | `"10y"` | Minimum estimated crack time per key (supports `y`, `d`, `h`, `m`, `s` units) |
| `min_password_length` | int | `14` | Minimum password length enforced for every encryption key |

Encryption keys are loaded from environment variables:

```bash
MARKATA_GO_ENCRYPTION_KEY_DEFAULT=my-password
MARKATA_GO_ENCRYPTION_KEY_PERSONAL=another-password
```

You can also override the enforcement policy:

```bash
MARKATA_GO_ENCRYPTION_ENFORCE_STRENGTH=false
MARKATA_GO_ENCRYPTION_MIN_ESTIMATED_CRACK_TIME=5d
MARKATA_GO_ENCRYPTION_MIN_PASSWORD_LENGTH=20
```

**Related:** See the [Encryption Guide](encryption.md) for complete documentation on making posts private.

---

## Complete Configuration Example

Here's a comprehensive example showing all available options:

```toml
[markata-go]
# Site metadata
title = "My Awesome Blog"
description = "A blog about technology and life"
url = "https://example.com"
author = "Jane Doe"
language = "en"
author_url = "https://example.com/about/"
managing_editor = "editor@example.com (Jane Doe)"
webmaster = "webmaster@example.com (Jane Doe)"
copyright = "Copyright 2026 Jane Doe"

# Directory configuration
output_dir = "public"
templates_dir = "templates"
assets_dir = "static"

# Plugin configuration
hooks = ["default"]
disabled_hooks = []

# Performance
concurrency = 0  # Auto-detect

# Content discovery
[markata-go.glob]
patterns = [
    "posts/**/*.md",
    "pages/*.md",
    "docs/**/*.md"
]
use_gitignore = true

# Markdown processing
[markata-go.markdown]
extensions = ["tables", "strikethrough", "autolinks", "tasklist"]

# Feed defaults (inherited by all feeds)
[markata-go.feed_defaults]
items_per_page = 10
orphan_threshold = 3

[markata-go.feed_defaults.formats]
html = true
rss = true
atom = false
json = false

[markata-go.feed_defaults.templates]
html = "feed.html"
rss = "rss.xml"
atom = "atom.xml"
json = "feed.json"
card = "partials/card.html"

[markata-go.feed_defaults.syndication]
max_items = 20
include_content = true

# Home page feed
[[markata-go.feeds]]
slug = ""
title = "Latest Posts"
filter = "published == True"
sort = "date"
reverse = true
items_per_page = 5

[markata-go.feeds.formats]
html = true

# Blog archive
[[markata-go.feeds]]
slug = "blog"
title = "Blog"
description = "All blog posts"
filter = "published == True"
sort = "date"
reverse = true

[markata-go.feeds.formats]
html = true
rss = true
atom = true
json = true

# Featured posts
[[markata-go.feeds]]
slug = "featured"
title = "Featured Posts"
description = "Hand-picked featured content"
filter = "published == True and featured == True"
sort = "date"
reverse = true
items_per_page = 6

[markata-go.feeds.formats]
html = true
rss = true

# Authors
[markata-go.authors]
generate_pages = false
feeds_enabled = false

[markata-go.authors.authors.jane]
name = "Jane Doe"
role = "author"
url = "https://janedoe.com"
active = true
default = true

# CSS minification
[markata-go.css_minify]
enabled = true
exclude = []
preserve_comments = ["Copyright"]

# JS minification
[markata-go.js_minify]
enabled = true
exclude = []
```

## Common Configuration Patterns

### Minimal Blog

```toml
[markata-go]
title = "My Blog"
url = "https://myblog.com"
output_dir = "public"

[markata-go.glob]
patterns = ["posts/**/*.md"]

[[markata-go.feeds]]
slug = "blog"
title = "Blog"
filter = "published == True"
sort = "date"
reverse = true

[markata-go.feeds.formats]
html = true
rss = true
```

### Documentation Site

```toml
[markata-go]
title = "Project Docs"
url = "https://docs.example.com"
output_dir = "site"

[markata-go.glob]
patterns = ["docs/**/*.md"]

[markata-go.markdown]
extensions = ["tables", "tasklist"]

[[markata-go.feeds]]
slug = ""
title = "Documentation"
filter = "true"
sort = "title"

[markata-go.feeds.formats]
html = true
```

### Multi-Section Site

```toml
[markata-go]
title = "My Site"
url = "https://example.com"
output_dir = "public"

[markata-go.glob]
patterns = ["content/**/*.md"]

# Blog section
[[markata-go.feeds]]
slug = "blog"
title = "Blog"
filter = "'blog' in tags"
sort = "date"
reverse = true

[markata-go.feeds.formats]
html = true
rss = true

# Tutorials section
[[markata-go.feeds]]
slug = "tutorials"
title = "Tutorials"
filter = "'tutorial' in tags"
sort = "date"
reverse = true

[markata-go.feeds.formats]
html = true
rss = true

# Projects section
[[markata-go.feeds]]
slug = "projects"
title = "Projects"
filter = "'project' in tags"
sort = "title"

[markata-go.feeds.formats]
html = true
```

### CI/CD Deployment

```toml
[markata-go]
title = "My Site"
# URL set via environment variable in CI
output_dir = "dist"

[markata-go.glob]
patterns = ["content/**/*.md"]
```

Then in CI:

```bash
MARKATA_GO_URL=https://example.com markata-go build
```

## CSS and JS Minification

markata-go automatically minifies CSS and JS files in the output directory to reduce file sizes.

### CSS Minification

```toml
[markata-go.css_minify]
enabled = true                      # Enable CSS minification (default: true)
exclude = ["variables.css"]         # Files to skip
preserve_comments = ["Copyright"]   # Comment substrings to preserve
```

| Option | Default | Description |
|--------|---------|-------------|
| `enabled` | `true` | Enable/disable CSS minification |
| `exclude` | `[]` | File names or glob patterns to skip (e.g., `"vendor-*.css"`) |
| `preserve_comments` | `[]` | Substrings that identify comments to keep |

### JS Minification

```toml
[markata-go.js_minify]
enabled = true                      # Enable JS minification (default: true)
exclude = ["pagefind-ui.js"]        # Files to skip
```

| Option | Default | Description |
|--------|---------|-------------|
| `enabled` | `true` | Enable/disable JS minification |
| `exclude` | `[]` | File names or glob patterns to skip (e.g., `"vendor-*.js"`) |

Files with `.min.js` extension are automatically skipped (already minified).

Both plugins run at the very end of the Write stage (`PriorityLast`), after all CSS/JS-generating plugins have finished.

Both plugins memoize exact current-input minification results in the private
`<content_dir>/.markata/asset-minify/v1/` cache. A nonempty top-level
`[markata-go] cache_dir` override is used as supplied, independently of whether
the build cache is enabled. Keep the cache outside published output; unsafe
locations warn and disable reuse. CSS `preserve_comments` changes invalidate
the transform recipe; changed recipes always minify current stage-input bytes.
Comments already removed require regeneration from authoritative source, not
recovery from a historical cache snapshot.
Reuse requires the exact transform input and matching recipe, never a prior
output or retry digest. Hits still publish atomically with normal `0644`
permissions (except empty JavaScript); repeated retained CSS output may need
fresh transformations because comment handling can be non-idempotent.
Fast mode, exclusions, disabled plugins, and `_pagefind` remain skipped.
See [Persistent CSS and JavaScript Minification](performance.md#persistent-css-and-javascript-minification)
for recovery, upgrade, privacy, and measurement guidance.

---

## Validation

markata-go validates configuration on load and can report both errors and warnings:

**Errors** (prevent build):
- Invalid URL format (missing scheme or host)
- Negative `items_per_page` or `orphan_threshold`
- Negative `concurrency`

**Warnings** (build continues):
- Empty glob patterns (no files will be processed)
- Feed with no output formats enabled

Run validation explicitly:

```bash
markata-go config validate
```

## See Also

- [Getting Started](/docs/getting-started/) - Quick start guide
- [Themes Guide](/docs/guides/themes/) - Theme and palette customization
- [Feeds Guide](/docs/guides/feeds/) - Detailed feed configuration
- [Blogroll Guide](/docs/guides/blogroll/) - RSS reader and blogroll setup
- [Templates Guide](/docs/guides/templates/) - Template configuration and usage
- [Sidebars Guide](/docs/guides/sidebars/) - Sidebar navigation setup
- [Frontmatter Guide](/docs/guides/frontmatter/) - Post metadata reference
- [CLI Reference](/docs/reference/cli/) - Command-line interface reference
- [Plugin Reference](/docs/reference/plugins/) - Plugin configuration and development


## Video loading preferences

```toml
[markata-go.md_video]
autoplay = false
preload = "none"
controls = true
```

These defaults show native controls without automatic video downloads. Setting
`autoplay = true` preserves author intent in generated markup. Full Quality or
already measured fast Auto at page entry can start unrequested authored videos; Save Data and
constrained Auto suppress expensive automatic playback. Raw HTML follows its
own autoplay attribute rather than the Markdown setting. See
[Adaptive loading](../adaptive-loading/) for user modes and no-JS behavior, and
[md_video](../plugins/#md_video) for the full plugin configuration.
