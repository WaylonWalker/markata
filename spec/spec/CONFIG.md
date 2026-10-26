# Configuration System Specification

The configuration system is designed to be:
- **Discoverable** - Plugins declare their config with descriptions
- **Flexible** - Multiple file formats and locations
- **Hierarchical** - Namespaced under the tool name
- **Mergeable** - Multiple sources combine intelligently

## Configuration Discovery

### File Locations

Unless `--config` names an explicit file, markata-go searches for configuration
files in this order (first found wins):

```
1. Current directory:    ./markata-go.toml
2. Current directory:    ./markata-go.yaml
3. Current directory:    ./markata-go.yml
4. Current directory:    ./markata-go.json
5. User config dir:      ~/.config/markata-go/config.toml
```

Commands that load a full site configuration MUST use built-in defaults plus
`MARKATA_GO_*` environment overrides when no file is found. Commands that read
or edit a value directly in a file, such as `config get` and `config set`, MUST
require a file. An explicit `--config` path MUST be used as given; a missing
explicit file MUST NOT silently fall back to defaults.

### Supported Formats

| Extension | Format | Notes |
|-----------|--------|-------|
| `.toml` | TOML | Recommended, best for nested config |
| `.yaml`, `.yml` | YAML | Good for complex structures |
| `.json` | JSON | Strict, good for programmatic generation |

## Config Composition

markata-go supports config composition through an `include` key under `[markata-go]`.

```toml
[markata-go]
include = [
  "config/base/*.toml",
  "config/feeds/*.toml",
  "config/rss.toml",
]
```

`include` supports explicit file paths, glob patterns, recursive includes, and paths resolved relative to the file that declared the include.

### Composition Order

Resolved precedence is:

1. built-in defaults
2. the root config selected by discovery or `--config`
3. included files in declaration order
4. glob matches in lexicographic order
5. environment variable overrides

Later values win over earlier values.

### Repeated Includes And Cycles

- A file included more than once in the same resolution graph is loaded once.
- Include cycles are rejected with a clear error that shows the cycle path.

### Merge Semantics

- scalar values: last explicit value wins
- explicit `false`, `0`, and `""` count as real overrides
- tables/maps: deep merge
- arrays of scalars: replace
- `[[markata-go.feeds]]`: merge by `slug`

For feeds, a new `slug` appends a new feed. A repeated `slug` merges into the existing feed, and later fragments win on conflicts.

### Typed Projection Of Core Settings

TOML, YAML, and JSON configuration must materialize these supported settings into
their corresponding typed `models.Config` fields before lifecycle execution:

| Configuration path | Typed field | Lifecycle consumer |
|--------------------|-------------|--------------------|
| `head` | `Head` | HTML template head rendering |
| `template_presets` | `TemplatePresets` | Per-format template selection |
| `default_templates` | `DefaultTemplates` | Global per-format template fallback |
| `theme_calendar` | `ThemeCalendar` | Seasonal theme plugin |
| `error_pages` | `ErrorPages` | Static 404 generation |
| `resource_hints` | `ResourceHints` | Resource-hint generation when enabled |
| `markdown.highlight` | `MarkdownConfig.Highlight` | Markdown and Chroma rendering |

The same typed projection is used by `Load`, `LoadFromString`, and
`LoadWithMerge`. Configuration maps are deep-merged and scalar arrays are
replaced. Explicit `false`, `0`, and empty-string values remain overrides.

The affected defaults are: markdown highlighting enabled, theme calendar
disabled, 404 pages enabled, and resource hints enabled with auto-detection.
Template maps and head elements have empty defaults.

### Format Examples

**TOML (recommended):**
```toml
[my-ssg]
output_dir = "public"
url = "https://example.com"

[my-ssg.feeds.defaults]
items_per_page = 10

[[my-ssg.feeds]]
slug = "blog"
filter = "published == True"
```

**YAML:**
```yaml
my-ssg:
  output_dir: public
  url: https://example.com

  feeds:
    defaults:
      items_per_page: 10

    items:
      - slug: blog
        filter: "published == True"
```

**JSON:**
```json
{
  "my-ssg": {
    "output_dir": "public",
    "url": "https://example.com",
    "feeds": {
      "defaults": {
        "items_per_page": 10
      },
      "items": [
        {
          "slug": "blog",
          "filter": "published == True"
        }
      ]
    }
  }
}
```

---

## Configuration Namespacing

All configuration lives under the tool name namespace:

```toml
# Root namespace - minimal, mostly metadata
[my-ssg]
output_dir = "public"
url = "https://example.com"

# Plugin namespaces
[my-ssg.glob]
patterns = ["**/*.md"]

[my-ssg.markdown]
extensions = ["tables", "admonitions"]

[my-ssg.feeds]
# Feed-specific config

[my-ssg.serve]
port = 3000
```

### Why Namespacing?

1. **Avoids conflicts** with other tools in shared config files (`pyproject.toml`)
2. **Clear ownership** - each plugin owns its namespace
3. **Tooling friendly** - editors can provide completions per-namespace
4. **Discoverable** - `my-ssg config list` shows all namespaces

### Root-Level Fields

Only essential, cross-cutting concerns live at the root:

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `output_dir` | Path | `"output"` | Build output directory |
| `url` | URL? | null | Site base URL (needed by many plugins) |
| `title` | string? | null | Site title |
| `description` | string? | null | Site description |
| `author` | string? | null | Default author |
| `license` | string `\|` bool | `"cc-by-4.0"` | Select from the supported license keys below or set to `false` to disable the footer attribution. |
| `lang` | string | `"en"` | Site language |
| `hooks` | string[] | `["default"]` | Plugins to load |
| `disabled_hooks` | string[] | `[]` | Plugins to exclude |

### Configless Small Sites

When no configuration file is present, the default content patterns discover
Markdown at the site root and in the `pages/` and `posts/` directories. The
implicit root feed renders an HTML homepage of published posts as well as
`/rss.xml` and `/atom.xml`; the default archive remains available at
`/archive/`. When no posts are published yet, both routes still render with
guidance to set `published: true`.
The default navigation links to `/` and `/archive/`; it does not assume a
custom `/blog/` feed exists.

An explicit root feed (`slug = ""`) remains authoritative. To render only one
Markdown file with the default theme and no generated home, archive, or feed
pages, pass it to the CLI. Single-file builds MUST:

- publish that file at the site root (`<output>/index.html`, href `/`),
  regardless of its path or frontmatter slug;
- still emit theme assets (CSS, JS, fonts, palettes);
- skip site-level generators: feeds, subscription feeds, series, blogroll,
  prev/next, archive and listing pages, sitemap, content index, `.well-known`,
  image library, redirects, random-post, 404 page, garden view, and Pagefind;
- hide search UI, feed `<link rel="alternate">` fallbacks, and the default
  Home/Archive navigation (configured nav items still render);
- bypass the build cache.

```bash
markata-go pages/sample.md
# equivalent:
markata-go build pages/sample.md
# serve accepts the same argument and keeps single-file mode on every rebuild:
markata-go serve pages/sample.md
```

Everything else goes in a plugin namespace.

### License configuration

The root `license` key lets you declare how visitors may reuse your content. It accepts either a string key or the literal `false` value:

- **String keys** register a license that is rendered in the footer and made available to templates via `config.license`. Use one of the supported values:
  - `all-rights-reserved` – All rights reserved (no reuse allowed).
  - `cc-by-4.0` (recommended) – Creative Commons Attribution 4.0 International (`https://creativecommons.org/licenses/by/4.0/`).
  - `cc-by-sa-4.0` – Creative Commons Attribution-ShareAlike 4.0 International (`https://creativecommons.org/licenses/by-sa/4.0/`).
  - `cc-by-nc-4.0` – Creative Commons Attribution-NonCommercial 4.0 International (`https://creativecommons.org/licenses/by-nc/4.0/`).
  - `cc-by-nd-4.0` – Creative Commons Attribution-NoDerivatives 4.0 International (`https://creativecommons.org/licenses/by-nd/4.0/`).
  - `cc-by-nc-sa-4.0` – Creative Commons Attribution-NonCommercial-ShareAlike 4.0 International (`https://creativecommons.org/licenses/by-nc-sa/4.0/`).
  - `mit` – MIT License (`https://opensource.org/licenses/MIT`).

- **Boolean `false`** suppresses the license footer and prevents the validation warning (useful for sites that intentionally publish without an explicit license).
- **Omitted key** inherits the default `cc-by-4.0` license. Set another
  supported string to override it, or `false` to suppress the footer
  attribution.

The default scaffolding details from `markata-go config init` include `license = "cc-by-4.0"`, so new sites ship with the recommended Creative Commons attribution out of the box.

---

## Plugin Configuration Declaration

Plugins declare their configuration schema, defaults, and descriptions. This enables:
- Config validation
- Auto-generated documentation
- Editor completions
- CLI config helpers

### Declaration Format

Each plugin exports a config schema:

**Python (Pydantic):**
```python
from pydantic import BaseModel, Field

class GlobConfig(BaseModel):
    """Configuration for the glob plugin."""

    patterns: list[str] = Field(
        default=["**/*.md"],
        description="Glob patterns to find content files"
    )
    use_gitignore: bool = Field(
        default=True,
        description="Respect .gitignore when finding files"
    )
    exclude: list[str] = Field(
        default=[],
        description="Patterns to exclude from results"
    )

class Config(BaseModel):
    """Adds glob section to config."""
    glob: GlobConfig = Field(
        default_factory=GlobConfig,
        description="File discovery settings"
    )

@hook_impl
def config_model(core):
    core.register_config(Config, namespace="glob")
```

**TypeScript (Zod):**
```typescript
import { z } from 'zod';

export const GlobConfig = z.object({
  patterns: z.array(z.string())
    .default(["**/*.md"])
    .describe("Glob patterns to find content files"),

  use_gitignore: z.boolean()
    .default(true)
    .describe("Respect .gitignore when finding files"),

  exclude: z.array(z.string())
    .default([])
    .describe("Patterns to exclude from results"),
});

export function configModel(core: Core) {
  core.registerConfig(GlobConfig, { namespace: "glob" });
}
```

**Go:**
```go
type GlobConfig struct {
    // Glob patterns to find content files
    Patterns []string `toml:"patterns" default:"[\"**/*.md\"]" description:"Glob patterns to find content files"`

    // Respect .gitignore when finding files
    UseGitignore bool `toml:"use_gitignore" default:"true" description:"Respect .gitignore when finding files"`

    // Patterns to exclude from results
    Exclude []string `toml:"exclude" description:"Patterns to exclude from results"`
}

func (p *GlobPlugin) ConfigModel(core *Core) {
    core.RegisterConfig("glob", GlobConfig{}, ConfigOptions{
        Description: "File discovery settings",
    })
}
```

**Rust:**
```rust
use serde::{Deserialize, Serialize};

/// Configuration for the glob plugin
#[derive(Debug, Deserialize, Serialize)]
#[serde(default)]
pub struct GlobConfig {
    /// Glob patterns to find content files
    #[serde(default = "default_patterns")]
    pub patterns: Vec<String>,

    /// Respect .gitignore when finding files
    #[serde(default = "default_true")]
    pub use_gitignore: bool,

    /// Patterns to exclude from results
    #[serde(default)]
    pub exclude: Vec<String>,
}

fn default_patterns() -> Vec<String> {
    vec!["**/*.md".to_string()]
}

impl Plugin for Glob {
    fn config_model(&self, core: &mut Core) {
        core.register_config::<GlobConfig>("glob");
    }
}
```

### Config Metadata

Each config field should have:

| Metadata | Purpose |
|----------|---------|
| `description` | Human-readable explanation |
| `default` | Default value if not specified |
| `type` | Data type for validation |
| `required` | Whether field must be provided |
| `deprecated` | Mark old config options |
| `env_var` | Environment variable override |
| `examples` | Example values |
| `see_also` | Related config fields |

### Example with Full Metadata

```python
class MarkdownConfig(BaseModel):
    """Markdown rendering configuration."""

    backend: str = Field(
        default="auto",
        description="Markdown parser backend",
        examples=["markdown-it", "commonmark", "mistune"],
    )

    extensions: list[str] = Field(
        default=["tables", "admonitions", "footnotes"],
        description="Markdown extensions to enable",
    )

    highlight_theme: str = Field(
        default="github-dark",
        description="Syntax highlighting theme",
        deprecated="Use markdown.highlight.theme instead",
        see_also=["markdown.highlight"],
    )

    class Highlight(BaseModel):
        """Syntax highlighting settings."""
        enabled: bool = Field(default=True, description="Enable syntax highlighting")
        theme: str = Field(default="github-dark", description="Color theme")
        line_numbers: bool = Field(default=False, description="Show line numbers")

    highlight: Highlight = Field(
        default_factory=Highlight,
        description="Syntax highlighting configuration",
    )
```

---

## Configuration Resolution

### Merge Order (lowest to highest precedence)

```
┌─────────────────────────────────────────────────────────────────────┐
│                    CONFIGURATION RESOLUTION                          │
├─────────────────────────────────────────────────────────────────────┤
│  1. Built-in defaults (from plugin declarations)                     │
│  2. Global config file (~/.config/my-ssg/config.toml)               │
│  3. Local config file (./my-ssg.toml)                               │
│  4. Environment variables (MY_SSG_SECTION_KEY)                      │
│  5. CLI arguments (--output-dir public)                             │
│                                                                      │
│  Later sources OVERRIDE earlier sources                              │
│  Nested objects are MERGED, not replaced                            │
└─────────────────────────────────────────────────────────────────────┘
```

### Environment Variables

Environment variables follow the pattern: `{NAME}_{SECTION}_{KEY}`

```bash
# Set output directory
MY_SSG_OUTPUT_DIR=public

# Set nested config
MY_SSG_FEEDS_DEFAULTS_ITEMS_PER_PAGE=20
MY_SSG_MARKDOWN_HIGHLIGHT_THEME=monokai

# Boolean values
MY_SSG_GLOB_USE_GITIGNORE=true
MY_SSG_GLOB_USE_GITIGNORE=1
MY_SSG_GLOB_USE_GITIGNORE=yes

# List values (comma-separated)
MY_SSG_GLOB_PATTERNS="posts/**/*.md,pages/*.md"
```

### CLI Arguments

Common config options have CLI flags:

```bash
my-ssg build --output-dir public --url https://example.com
my-ssg serve --port 8080 --host 0.0.0.0
```

Arbitrary config can be set with `--config` or `-c`:

```bash
my-ssg build -c feeds.defaults.items_per_page=20
my-ssg build -c markdown.highlight.theme=monokai
```

### Merge Behavior

**Scalar values:** Later wins
```toml
# Global: output_dir = "dist"
# Local:  output_dir = "public"
# Result: output_dir = "public"
```

**Objects:** Deep merge
```toml
# Global:
[my-ssg.feeds.defaults.formats]
html = true
rss = true

# Local:
[my-ssg.feeds.defaults.formats]
atom = true

# Result:
html = true   # from global
rss = true    # from global
atom = true   # from local
```

**Lists:** Replace (not append)
```toml
# Global: patterns = ["**/*.md"]
# Local:  patterns = ["posts/*.md", "pages/*.md"]
# Result: patterns = ["posts/*.md", "pages/*.md"]
```

**List append syntax (optional):**
```toml
# To append instead of replace
patterns = ["posts/*.md"]
patterns_append = ["pages/*.md"]
# Result: ["posts/*.md", "pages/*.md"]
```

---

## Configuration CLI

### `config show`

Display resolved configuration:

- bare `my-ssg config` MUST behave like `my-ssg config show`
- `config show` MUST honor the same config resolution path as build commands, including `--config` and `--merge-config`
- conflicting output flags such as `--json` with `--toml` MUST fail with a usage error (exit code `2`)

```bash
$ my-ssg config show
output_dir = "public"
url = "https://example.com"

[glob]
patterns = ["**/*.md"]
use_gitignore = true

[feeds.defaults]
items_per_page = 10
...
```

With source information:

```bash
$ my-ssg config show --sources
output_dir = "public"          # ./my-ssg.toml
url = "https://example.com"    # MY_SSG_URL env var

[glob]
patterns = ["**/*.md"]         # (default)
use_gitignore = true           # ~/.config/my-ssg/config.toml
```

### Source Diagnostics

`config show --sources` MUST emit valid YAML with a source comment for every
effective scalar and sequence item. Source comments MUST use this diagnostic
format:

```text
# <kind> <location>:<line>:<column> — <change hint>
```

Examples:

```yaml
output_dir: public  # file /site/markata-go.toml:12:1 — edit it or run: markata-go config set output_dir <value>
url: https://example.com  # environment MARKATA_GO_URL — export a new value
concurrency: 4  # CLI --output — pass a different flag value
use_gitignore: true  # default — run: markata-go config set glob.use_gitignore <value>
```

The resolver MUST record the winning source while applying configuration layers
in precedence order: defaults, root and included files, `--merge-config` files,
environment variables, then applicable CLI overrides. File sources MUST report
an absolute cleaned path and 1-indexed line and column. Environment comments
MUST include only variable names, never values.

Default interactive YAML output SHOULD include source comments and syntax
highlighting when stdout is a terminal. Piped output, `--no-color`, `NO_COLOR`,
and `TERM=dumb` MUST remain valid uncolored YAML. JSON and TOML output MUST
remain parseable and unannotated. `--sources` is the canonical explicit source
flag; `--annotate` MAY remain as a compatibility alias.

### `config list`

List all available configuration options:

```bash
$ my-ssg config list

[my-ssg] Core configuration
  output_dir     Path     "output"    Build output directory
  url            URL?     null        Site base URL
  title          string?  null        Site title
  hooks          string[] ["default"] Plugins to load

[my-ssg.glob] File discovery settings
  patterns       string[] ["**/*.md"] Glob patterns to find content files
  use_gitignore  bool     true        Respect .gitignore when finding files
  exclude        string[] []          Patterns to exclude from results

[my-ssg.markdown] Markdown rendering configuration
  backend        string   "auto"      Markdown parser backend
  extensions     string[] [...]       Markdown extensions to enable
  ...
```

### `config get`

Get a specific value:

```bash
$ my-ssg config get feeds.defaults.items_per_page
10

$ my-ssg config get glob.patterns
["**/*.md"]

$ my-ssg config get glob.patterns --json

Behavior:

- Reads values directly from the config file to preserve source-of-truth behavior.
- Supports TOML, YAML, and JSON.
- Uses tree-sitter parsing for TOML/YAML in CGO-enabled builds to locate byte ranges without reformatting.
- CGO-disabled builds parse TOML/YAML via full decode.
- JSON output may be re-emitted for structured values.
["**/*.md"]
```

### `config set`

Set a value (writes to local config file):

```bash
$ my-ssg config set output_dir public
$ my-ssg config set feeds.defaults.items_per_page 20
$ my-ssg config set glob.patterns '["posts/*.md", "pages/*.md"]'

Behavior:

- TOML/YAML are updated with byte-range edits to preserve formatting and comments in CGO-enabled builds.
- CGO-disabled builds re-emit TOML/YAML with standard encoders.
- JSON is re-emitted with stable indentation.
- File permissions are preserved.
```

### `config init`

Generate a starter config file:

```bash
$ my-ssg config init
Created my-ssg.toml with default configuration

$ my-ssg config init --format yaml
Created my-ssg.yaml with default configuration

$ my-ssg config init --full
Created my-ssg.toml with all options documented
```

### `config validate`

Validate configuration:

- `config validate` MUST honor the same config resolution path as build commands, including `--config` and `--merge-config`

```bash
$ my-ssg config validate
✓ Configuration is valid

$ my-ssg config validate
✗ Configuration errors:
  - feeds.defaults.items_per_page: must be >= 0, got -5
  - glob.patterns: must be non-empty array
  - unknown field: my-ssg.typo_field
```

### `config docs`

Generate configuration documentation:

```bash
$ my-ssg config docs
# my-ssg Configuration

## Core Settings

### output_dir
- Type: Path
- Default: "output"
- Environment: MY_SSG_OUTPUT_DIR

Build output directory. All generated files will be written here.

### url
- Type: URL (optional)
- Default: null
- Environment: MY_SSG_URL

Site base URL. Required for generating absolute URLs in feeds and sitemaps.
...
```

---

## Complete Configuration Reference

### Core (`[my-ssg]`)

```toml
[my-ssg]
# Build output directory
output_dir = "output"

# Site metadata
url = "https://example.com"      # Base URL for absolute links
title = "My Site"                 # Site title
description = "A great site"      # Site description
author = "Jane Doe"               # Default author
language = "en"                   # Site language for feeds and metadata
author_url = "https://example.com/about/"
managing_editor = "editor@example.com (Jane Doe)"
webmaster = "webmaster@example.com (Jane Doe)"
copyright = "Copyright 2026 Jane Doe"

# Plugin loading
hooks = ["default"]               # Plugins to load
disabled_hooks = []               # Plugins to exclude

# Build settings
concurrency = 0                   # Worker threads (0 = auto)
```

### Glob (`[my-ssg.glob]`)

```toml
[my-ssg.glob]
patterns = ["pages/**/*.md", "posts/**/*.md"] # File patterns to match
use_gitignore = true              # Respect .gitignore
exclude = ["node_modules/**"]     # Patterns to exclude
```

### Markdown (`[my-ssg.markdown]`)

```toml
[my-ssg.markdown]
backend = "auto"                  # Parser backend
extensions = ["tables", "admonitions", "footnotes"]

[my-ssg.markdown.highlight]
enabled = true
theme = "github-dark"
line_numbers = false
```

### Feeds (`[my-ssg.feeds]`)

```toml
[my-ssg.feeds.defaults]
items_per_page = 10
orphan_threshold = 3

[my-ssg.feeds.defaults.formats]
html = true
simple_html = true
rss = true
atom = true
json = true
sitemap = true

[my-ssg.feeds.syndication]
max_items = 20
include_content = false
site_archive_disabled = false
feed_archives_disabled = false

[my-ssg.feeds_page]
enabled = true
title = "Feeds"
description = "Browse the public feeds available on this site."
template = "feeds.html"
slug_prefix = "feeds"

[[my-ssg.feeds]]
slug = "archive"
title = "Archive"
description = "All posts"
filter = "published == True"
sort = "date"
reverse = true
limit = 0
offset = 0

[[my-ssg.feeds]]
slug = "blog"
archive_disabled = true
```

### Serve (`[my-ssg.serve]`)

```toml
[my-ssg.serve]
port = 3000
host = "localhost"
livereload = true
open_browser = false
debounce_ms = 100
```

Serve mode starts the HTTP server immediately while the initial build runs in the background.
During builds, a status banner is injected into HTML responses, and a minimal 404 page is
served until the generated 404.html is available.

### Assets (`[my-ssg.assets]`)

```toml
[my-ssg.assets]
dir = "static"
output_subdir = ""

[my-ssg.assets.fingerprint]
enabled = false
algorithm = "sha256"
length = 8
exclude = ["robots.txt", "favicon.ico"]
```

markata-go also supports self-hosting third-party CDN assets (HTMX, GLightbox, Mermaid, Chart.js, Cal-Heatmap, D3, Lite YouTube, Reveal.js) through the assets config. When enabled, assets are downloaded into a cache directory and copied into the output under the vendor directory. Templates can use the `asset_urls` mapping to reference the local paths.

```toml
[markata-go.assets]
mode = "self-hosted"           # default: "self-hosted"
cache_dir = ".markata/assets-cache"
output_dir = "assets/vendor"
verify_integrity = true
```

### CSS Purge (`[my-ssg.css_purge]`)

```toml
[my-ssg.css_purge]
enabled = false
verbose = false
preserve = ["js-*", "htmx-*", "theme-*", "palette-*"]
preserve_attributes = ["data-theme", "data-palette"]
skip_files = ["vendor/*", "normalize.css"]
warning_threshold = 0
```

CSS purge removes unused rules by scanning generated HTML and keeping only selectors
that are actually present.
CSS purge MUST preserve the default theme's `md-code-*` classes even when
absent from generated HTML, because the code-copy enhancement adds them at
browser runtime. The generated `code.css` MUST retain positioning and keyboard
focus rules for `.md-code-copy` and `pre.md-code-enhanced`. This protection
MUST NOT turn off purging for unrelated selectors. Sites overriding
`preserve` should include `md-code-*` in their custom list.
 The purge logic always preserves key @-rules and keeps
pseudo-only selectors like `:root` or `::selection` to avoid dropping base/theme styles.
Full builds MAY cache the aggregate selector set and purged CSS under
`.markata-cache/css-purge` in the content directory. Selector cache validity MUST include
the sorted generated HTML paths and exact file bytes. Purged CSS cache validity MUST
include stylesheet bytes, the selector set, and purge options. Missing or invalid cache
entries MUST fall back to the normal scan or purge path.

### Tailwind (`[my-ssg.tailwind]`)

```toml
[my-ssg.tailwind]
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

Tailwind automation runs a markata-managed standalone Tailwind CLI and writes the
compiled CSS into your assets directory so it is copied and fingerprinted like any
other static file. The zero-setup default is: enable Tailwind in config and build.
If `tailwind.css` is missing, markata-go generates a default input containing
`@tailwind base`, `@tailwind components`, and `@tailwind utilities`.

Behavior:

- `build = true` injects the Tailwind output during configure and performs the
  actual Tailwind rebuild in cleanup when needed.
- `preflight = false` is the default for markata-go-managed Tailwind configs so
  Tailwind utilities work without resetting the built-in theme's typography and
  spacing. Set `preflight = true` for Tailwind-first sites that want the reset.
- `input`/`output` resolve relative to `assets_dir` (absolute paths are respected).
- If `extra_args` is empty and `config_file` is unset, markata-go generates a
  temporary Tailwind config that scans a generated token manifest derived from
  rendered page HTML plus local JS/template sources. The generated config also
  sets `corePlugins.preflight` from `tailwind.preflight`. The manifest is hashed
  and cached so Tailwind is skipped when the effective utility set is unchanged.
- `include = "css"` ensures the output CSS is included in templates. If
  `theme.custom_css` is unset, it is set to the output path. The plugin does not
  override explicit `theme.custom_css` values.
- `include = "js"` injects `<script src="https://cdn.tailwindcss.com"></script>`
  into the document head. If `[my-ssg.assets].mode` is self-hosted, the JS is
  pulled from the vendor asset registry and served locally.
- `include = false` disables automatic inclusion; build can still run.
- If `include = "css"` and CSS purge is disabled, a validation warning is emitted.
- `auto_install = true` downloads and uses the managed Tailwind CLI (versioned,
  checksum verified) into the cache directory when needed. This is preferred over
  `PATH` for consistent builds. If disabled, `binary` or `PATH` is used.
- Fast mode skips Tailwind rebuilds when the compiled CSS asset already exists,
  keeping development builds fast without requiring a separate output directory.

### Theme (`[my-ssg.theme]`)

```toml
[my-ssg.theme]
name = "default"              # Theme name (built-in or installed)
custom_css = ""               # Path to custom CSS file (loaded after theme CSS)
fallback_mode = "dark"        # Fallback when system preference is unavailable: "dark" or "light"

# Theme-specific options (defined by theme)
[my-ssg.theme.options]
primary_color = "#3b82f6"     # Varies by theme
font_family = "system-ui"
show_toc = true

# CSS custom property overrides
[my-ssg.theme.variables]
"--color-primary" = "#8b5cf6"
"--color-primary-dark" = "#7c3aed"
"--font-body" = "Inter, system-ui"
"--content-width" = "70ch"
```

See [THEMES.md](./THEMES.md) for complete theming documentation.

### Post Formats (`[my-ssg.post_formats]`)

```toml
[my-ssg.post_formats]
html = true       # /slug/index.html (default: true)
markdown = true   # /slug.md - raw source with frontmatter (default: true)
text = true       # /slug.txt - plain terminal-friendly content (default: true)
ansi = true       # /slug.ansi - ANSI-styled terminal output (default: false)
og = true         # /slug/og/index.html - social card for screenshots
```

This section controls what output formats are generated for each post:

| Format | Default | Output Path | Description |
|--------|---------|-------------|-------------|
| `html` | `true` | `/slug/index.html` | Standard rendered HTML page |
| `markdown` | `true` | `/slug.md` | Raw markdown with reconstructed frontmatter |
| `text` | `true` | `/slug.txt` | Plain terminal-friendly content with no ANSI escapes |
| `ansi` | `false` | `/slug.ansi` | ANSI-styled terminal page output |
| `og` | `true` | `/slug/og/index.html` | OpenGraph card HTML (1200x630) for social screenshots |

Posts MAY override these site defaults in frontmatter with a `post_formats` mapping. Per-post overrides merge with `[my-ssg.post_formats]` key-by-key; omitted keys inherit the site setting.

```yaml
---
title: "Terminal-first post"
post_formats:
  ansi: true
  og: false
---
```

In this example, the post inherits the site defaults for `html`, `markdown`, and `text`, enables `.ansi` for this post only, and suppresses OG output for this post only.

`text` and `ansi` are separate explicit variants:

- `text` MUST emit readable plain text with no ANSI escape sequences.
- `ansi` MUST emit the same terminal-oriented structure with ANSI styling for capable clients.
- ANSI output is opt-in via the `.ansi` path; markata-go MUST NOT inject ANSI escape sequences into `.txt` output.

Terminal rendering for both `text` and `ansi` variants MUST:

- derive structure from rendered page content when rendered HTML is available
- preserve headings, emphasis, links, blockquotes, lists, horizontal rules, tables, admonitions, and code fences in terminal-safe form
- preserve image and video references as readable labeled URLs when media appears inline or in frontmatter
- degrade cleanly to plain text when ANSI styling is disabled
- keep canonical `.txt` endpoints readable in clients that do not support ANSI

Theme-aware ANSI rendering MUST derive colors from the active site palette when possible. Palette resolution for ANSI output follows:

1. `theme.palette`
2. `theme.palette_dark`
3. `theme.palette_light`
4. built-in dark fallback palette

**Directory-based Redirects for txt/md/ansi:**

For `.txt`, `.md`, and `.ansi` formats, content is placed at the canonical short URL (`/slug.txt`, `/slug.md`, `/slug.ansi`). Redirects are provided at `/slug.<ext>/index.html` (for hosts that serve `index.html` in a directory) and `/slug/index.<ext>/index.html` (for backwards compatibility).

Rendered feed/sidebar variant links MUST use the same canonical short URLs rather than nested `/slug/index.<ext>` paths.
If a post format is disabled in the resolved config, the corresponding sidebar link MUST be omitted.

**Special Files (robots, llms, humans, security, ads):**

Special web files have an inverted structure to serve content at their expected root-level locations:
- Content at `/slug.txt` (e.g., `/robots.txt`)
- HTML redirect at `/slug/index.txt/index.html` pointing to `/slug.txt`

This enables standard web txt files to be served at their expected locations:
- `/robots.txt` - Robot exclusion standard
- `/llms.txt` - AI/LLM guidance file  
- `/humans.txt` - Human-readable site info
- `/security.txt` - Security contact information
- `/ads.txt` - Authorized digital sellers

**Use cases:**
- **markdown**: API consumers, "view source" links, copy-paste code
- **text**: Standard web txt files, plain text readers, CLI tools
- **ansi**: `curl`, pagers, and intentional terminal reading experiences
- **og**: Automated social image generation with puppeteer/playwright

**Example:**
```toml
[markata-go.post_formats]
html = true
markdown = true  # Enable raw markdown output at /slug.md
text = true      # Enable plain terminal-friendly output at /slug.txt
ansi = true      # Enable ANSI terminal output at /slug.ansi
og = true        # Enable social card HTML for screenshot tools
```

---

### Glob Settings (`[my-ssg.glob]`)

```toml
[my-ssg.glob]
patterns = ["posts/**/*.md", "pages/**/*.md"]
use_gitignore = true
slug_mode = "flat"    # "flat" (default) or "path"
```

The glob section controls both file discovery and the default slug derivation strategy.

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `patterns` | string[] | `["pages/**/*.md", "posts/**/*.md"]` | Glob patterns to find content files |
| `use_gitignore` | bool | `true` | Respect `.gitignore` when discovering content |
| `slug_mode` | string | `"flat"` | How to derive slugs when frontmatter does not set `slug` |

`slug_mode` values:

- `flat` - Current markata behavior. Slugs come from the filename only. Examples: `posts/2026/hello.md -> hello`, `docs/index.md -> docs`.
- `path` - Derive slugs from the relative content path. Leading `posts/` and `pages/` path segments are removed. `index.md`, `README.md`, and `readme.md` become the root of their containing directory. Examples: `posts/notes/today.md -> notes/today`, `pages/docs/README.md -> docs`.

Optional path-specific overrides can be declared with `slug_rules`:

```toml
[my-ssg.glob]
slug_mode = "flat"

[[my-ssg.glob.slug_rules]]
prefix = "posts/blog"
mode = "flat"

[[my-ssg.glob.slug_rules]]
prefix = "posts/notes"
mode = "path"
```

Rules use the longest matching prefix, so more specific directories win.

Explicit frontmatter `slug` always wins over `slug_mode`.

---

### Well-Known Files (`[my-ssg.well_known]`)

```toml
[my-ssg.well_known]
enabled = true
auto_generate = ["host-meta", "host-meta.json", "webfinger", "nodeinfo", "time", "links"]

# Optional entries requiring config
ssh_fingerprint = "SHA256:abcdef..."
keybase_username = "username"
```

This section controls auto-generated `.well-known` endpoints derived from site metadata:

| Entry | Output Path | Description |
|-------|-------------|-------------|
| `host-meta` | `/.well-known/host-meta` | XRD host metadata for discovery |
| `host-meta.json` | `/.well-known/host-meta.json` | JSON host metadata (JRD) |
| `webfinger` | `/.well-known/webfinger` | WebFinger endpoint response |
| `nodeinfo` | `/.well-known/nodeinfo` + `/nodeinfo/2.0` | NodeInfo discovery + instance metadata |
| `time` | `/.well-known/time` | Build timestamp (RFC3339) |
| `links` | `/.well-known/links` + `/.well-known/internal-links` + `/external-links/` + `/internal-links/` | Outbound links grouped by target domain and internal links grouped by target URL, each with HTML view |
| `sshfp` | `/.well-known/sshfp` | SSH fingerprint text (requires `ssh_fingerprint`) |
| `keybase` | `/.well-known/keybase.txt` | Keybase verification (requires `keybase_username`) |

**Defaults:**
- `enabled` defaults to `true`
- `auto_generate` defaults to `host-meta`, `host-meta.json`, `webfinger`, `nodeinfo`, `time`, and `links`

**Notes:**
- If `auto_generate` is empty, only optional entries with explicit config are generated.
- `nodeinfo` generates both the discovery document and a minimal `/nodeinfo/2.0` payload.
- `time` is rebuilt on each build and is always UTC RFC3339.
- `links` follows the Jim Nielsen-style grouped JSON format for external links, and additionally writes internal links plus `/external-links/` and `/internal-links/`.

---

### WebSub (`[my-ssg.websub]`)

```toml
[my-ssg.websub]
enabled = true
hubs = ["https://hub.example.com/"]
```

When enabled, markata-go emits WebSub discovery links:

- HTML pages include `<link rel="hub" href="...">` for each hub
- RSS/Atom feeds include `rel="hub"` and `rel="self"` links

**Defaults:**
- `enabled` defaults to `false`
- `hubs` defaults to an empty list

---

### Templates (`[my-ssg.templates]`)

```toml
[my-ssg.templates.media]
trusted_domains = [
  "dropper.wayl.one",
  "dropper.waylonwalker.com",
  "dropper-dev.wayl.one",
]
```

- `media.trusted_domains` controls which hosts the built-in template helpers will decorate with `w`/`h` sizing parameters, derived posters, and `https` normalization. Relative URLs are always treated as trusted.
- The default values match the dropper CDN. Override this list when you serve media through a different host so that video posters and cached previews stay consistent.

## Serve Settings Sidebar

`markata-go serve` (never `build`) injects a settings sidebar into served HTML:
`window.__markataSettingsEndpoint = "/__markata/settings"` in `<head>` and
`<script src="/__markata/settings.js" defer>` before `</body>`. Static output
MUST NOT contain either. The script mounts a `<markata-dev-settings>` element
with a Shadow DOM, so site CSS cannot style it and it cannot style the site.

### Schema

The field list is derived by reflection over `models.Config` TOML tags, so new
config fields appear automatically. Each field has:

| Field | Meaning |
|-------|---------|
| `key` | Dotted path under `markata-go`, e.g. `theme.palette` |
| `section` | First path segment for nested keys, `""` for top-level keys |
| `kind` | `string`, `bool`, `int`, `float`, `list` (of strings), or `complex` |
| `value` | Effective value; `null` for unset optional (pointer) fields and for sensitive fields |
| `doc` | The Go doc comment of the model field (generated into `pkg/config/settings_docs_gen.go` by `go generate ./pkg/config`; a test fails when it is stale) |
| `options` | Suggested or allowed values: palettes, fontpacks, chroma styles, rendering-contract enums, the curated enum registry (`settingEnums`), or values parsed from a `Valid values:` / `Options:` doc line |
| `closed` | True when `options` is the complete set; the client renders a dropdown and the server rejects other non-empty values |
| `min`, `max` | Inclusive numeric range for `int`/`float` fields, when one applies (for example `0`-`1` for `theme.background.color_mix`) |
| `default` | Value the key has when no config file sets it (display only) |
| `source` | Last config file that defines the key, relative to the site |
| `sources` | Every config file that defines the key, in load order |
| `target` | File a change will be written to |
| `target_kind` | `override` for a `--merge-config` file, `global` for `~/.config/markata-go/config.toml`, omitted otherwise |
| `env` | Name of a set `MARKATA_GO_*` variable that overrides the file |
| `editable` | False for `complex`, `sensitive`, and `unsupported` fields |

- **complex**: maps, slices of structs, and types from other packages. Shown
  read-only with a summary such as "3 items".
- **sensitive**: the leaf key matches
  `(^|_)(secret|token|password|passphrase|api_key|private_key|credentials?)$`.
  The value is never sent to the browser.
- **unsupported**: model fields the config loader does not read from files.
  Each editable field is probed at startup by parsing a TOML document that sets
  it; fields that do not round-trip are marked unsupported.

### Target File

Candidates are the same load-order list as theme bake (root, `include` files,
`--merge-config` files). A change to `a.b.c` is written to the candidate that
defines the deepest prefix of the path (`a.b.c`, then `a.b`, then `a`); ties go
to the later file. With no match it goes to the root config, or a new
`markata-go.toml` when there is none. So a setting stays in the file that owns
it, and a new key lands next to its siblings.

- A `--merge-config` target is reported as `target_kind: "override"` and the
  client labels it before writing.
- When no `--config` is given and discovery fell back to the global
  `~/.config/markata-go/config.toml`, every non-override candidate is
  `global`. Bakes (settings and theme) into a global target are refused with
  409 so a site preview never edits the user's shared defaults.
- An unset (reset to default) targets every file that defines the key.

### Closed Options

- An empty string is always accepted and means "use the default".
- Palette keys accept any name `KnownPalette` accepts (contract IDs, aliases,
  loader palettes), even when not listed.
- Fontpack options are closed unless `theme.fontpacks_file` is set.
- A guard test fails when a field's doc lists quoted values but the field has
  no closed options (allowlist: free-form fields such as
  `components.share.position`).

### Preview

`POST /__markata/settings/preview` with the same body replaces the whole
preview set; an empty `changes` list resets it. A change may be
`{"key": "...", "unset": true}` (no `value`) to preview the key at its
default. Preview changes are held only
in the serve process's memory and are never written to disk:

1. Changes are validated as in Apply step 1, plus closed options and ranges
   (400 on failure).
2. The preview is loaded as a config overlay (`LoadOptions.Overlay`), merged
   after all config files and before defaults normalization and environment
   variables, so `MARKATA_GO_*` still wins. Unset keys are deleted from the
   merged raw config (`LoadOptions.Remove`) after the overlay. A new validation error or a value
   that does not load back is rejected with 422 and the previous preview is kept.
3. A full rebuild is queued. The preview fingerprint
   (`{"set": {...}, "reset": [...]}`) is stored in `Extra["config_overlay"]`
   and included in the build-cache config hash so cached pages re-render.
   While a preview is active the build cache lives in
   `<cache_dir>/serve-preview` (default `.markata/serve-preview`), so
   previewing, resetting, and stopping serve never invalidate the site's main
   build cache.

The response and `GET /__markata/settings` include `preview`, the current
list of previewed changes, `session`, a random ID for this serve process, and
`single_page` (true for `markata-go serve <file>`). Restarting serve discards
the server preview; the client keeps staged edits in `sessionStorage` and, when
`session` changes, re-sends them as a preview.

### Apply

`POST /__markata/settings` with `{"changes": [{"key": "...", "value": ...}], "dry_run": false}`
(max 200 changes, 256 KiB, unknown JSON fields rejected). A change with
`"unset": true` and no value removes the key from every file that defines it;
unsetting a key no file defines, or sending a value with `unset`, is a 400.
With `dry_run: true` nothing is written and the response lists `diffs`:
`[{target, kind, created, diff}]`, one unified diff (2 lines of context) per
file, with the values of sensitive keys replaced by `"…"`. The preview
endpoint rejects `dry_run`.

1. Each key MUST be editable, appear once, and have a value of its kind.
   Strings are at most 4096 characters with no control characters other than
   newline and tab. Lists have at most 256 items. Ints must fit the Go type.
2. Changes are grouped by target file and written with the format-preserving
   editor used by theme bake (comments and order kept, TOML tables, dotted
   keys, and multi-line arrays, nested YAML mappings, JSON objects).
3. The config is reloaded. If it gains a validation error that was not
   present before, or a key does not load back as the requested value (an
   unset key must equal its value in a load with the key removed), unless an
   env variable overrides it, every file is restored atomically and the
   response is 422. The whole apply, including rollback, holds the serve config
   lock, which manager creation for rebuilds also takes.
4. On success the baked keys are dropped from the preview, a full rebuild is
   queued, and live reload shows the result. Bake reads config from disk only,
   never from the preview overlay.

Status codes: 400 invalid request, 403 non-loopback client, cross-origin, or non-local `Host`,
409 layout the editor refuses (see THEMES.md Bake) or a global target, 422 rolled back.
Both settings and theme-bake endpoints require a loopback client address
(`RemoteAddr`), a same-origin request, and a `Host` of `localhost`,
`*.localhost`, or an IP literal, as a defense against DNS rebinding. The
loopback requirement means binding `serve` to `0.0.0.0` or a LAN address never
lets other hosts read or change the config.

### Client

- Toggle: a gear button at the bottom left (`Alt+,`, ignored while focus is in
  a page input, textarea, select, or contenteditable), with a badge that shows
  the number of unsaved changes. While the panel is closed and a preview is
  active, a chip reads "Previewing N unsaved changes · Reset · Review".
- The panel has search, a filter (all / set in config / changed), sections that
  can be collapsed, and per-field Revert. A pinned **Common** section repeats
  frequently edited keys (site identity, light/dark palettes, aesthetic,
  fontpack, text size, nav, footer, layout, search) above the full schema; both
  copies stay in sync. In single-page serve, sections for site-wide output
  (`glob`, feeds, `blogroll`, `tags`, `garden`, and similar) are hidden unless
  searching or after "Show hidden sections".
  It is full screen at 640px wide and below.
- Each field set in a config file has **Reset to default**, which stages an
  unset; the row shows the default it will fall back to. Rows show the
  `target` and flag `override` and `global` targets.
- Previewing a `theme.*` key that the theme picker stores in `localStorage`
  (palette, aesthetic, fontpack, text size, color mode) clears that stored pick
  so the preview is visible, and says so in the status line.
- The server preview is the source of truth. Edits commit on change, blur, or
  Enter, are checked client-side (options, ranges, kinds), and sent to the
  preview endpoint after a 300ms debounce. A rejected key is marked invalid
  inline and the rest are re-sent without it; the status line lists every
  change that is not being previewed and why.
- After a preview or bake the status line shows "Rebuilding… Ns" until the dev
  server reports the build finished (`markata:build-status` window event from
  the live reload client) or failed (the build error is shown).
- Closed options render as a `<select>`; open options as an input with a
  datalist; numbers get `min`/`max`/`step`.
- **Reset** clears the preview. Panel state (open, scroll, collapsed sections)
  lives in `sessionStorage` and survives live reload.
- **Bake** first sends a dry run and shows the per-file diffs with **Write
  N files** and **Cancel**. Changing any staged edit closes the confirmation.
  Global targets cannot be baked.
- `window.markataDevSettings.open(section?)`, `.close()`, `.isOpen()`, and `.reset()`.

## See Also

- [SPEC.md](./SPEC.md) - Core specification
- [THEMES.md](./THEMES.md) - Theming system
- [PLUGINS.md](./PLUGINS.md) - Plugin development
