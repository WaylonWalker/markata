# Configuration System Specification

The configuration system is designed to be:
- **Discoverable** - Plugins declare their config with descriptions
- **Flexible** - Multiple file formats and locations
- **Hierarchical** - Namespaced under the tool name
- **Mergeable** - Multiple sources combine intelligently

## Configuration Discovery

### File Locations

The system searches for configuration files in this order (first found wins):

```
1. CLI-specified:        --config path/to/config.toml
2. Current directory:    ./[name].toml (or .yaml, .yml, .json, .jsonc)
3. Current directory:    ./[name]/config.toml
4. pyproject.toml:       ./pyproject.toml (under [tool.name] section)
5. package.json:         ./package.json (under "name" key)
6. User config dir:      ~/.config/[name]/config.toml
7. User home:            ~/.[name].toml
8. User home dotdir:     ~/.[name]/config.toml
```

### Supported Formats

| Extension | Format | Notes |
|-----------|--------|-------|
| `.toml` | TOML | Recommended, best for nested config |
| `.yaml`, `.yml` | YAML | Good for complex structures |
| `.json` | JSON | Strict, good for programmatic generation |
| `.jsonc` | JSON with comments | JSON + `//` and `/* */` comments |

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
| `lang` | string | `"en"` | Site language |
| `hooks` | string[] | `["default"]` | Plugins to load |
| `disabled_hooks` | string[] | `[]` | Plugins to exclude |

Everything else goes in a plugin namespace.

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
["**/*.md"]
```

### `config set`

Set a value (writes to local config file):

```bash
$ my-ssg config set output_dir public
$ my-ssg config set feeds.defaults.items_per_page 20
$ my-ssg config set glob.patterns '["posts/*.md", "pages/*.md"]'
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
lang = "en"                       # Site language

# Plugin loading
hooks = ["default"]               # Plugins to load
disabled_hooks = []               # Plugins to exclude

# Build settings
concurrency = 0                   # Worker threads (0 = auto)
```

### Glob (`[my-ssg.glob]`)

```toml
[my-ssg.glob]
patterns = ["**/*.md"]            # File patterns to match
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
rss = true
atom = false
json = false

[my-ssg.feeds.syndication]
max_items = 20
include_content = false

[[my-ssg.feeds]]
slug = "blog"
title = "Blog"
filter = "published == True"
sort = "date"
reverse = true
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

### Theme (`[my-ssg.theme]`)

```toml
[my-ssg.theme]
name = "default"              # Theme name (built-in or installed)
custom_css = ""               # Path to custom CSS file (loaded after theme CSS)

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

---

## See Also

- [SPEC.md](./SPEC.md) - Core specification
- [THEMES.md](./THEMES.md) - Theming system
- [PLUGINS.md](./PLUGINS.md) - Plugin development
