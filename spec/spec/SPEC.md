# Static Site Generator Specification

A specification for building plugin-driven static site generators. Inspired by [markata](https://github.com/waylonwalker/markata) and the [whenwords](https://github.com/dbreunig/whenwords) spec-as-product philosophy.

## Design Principles

### 1. The Core Provides Only What Plugins Cannot

The core should be minimal, providing only:
- **Lifecycle orchestration** - Running build stages in order
- **Plugin loading** - Discovering, registering, and calling plugins
- **Caching** - Persistent cache for expensive operations
- **Data access** - Querying and filtering content
- **Configuration merging** - Combining config from multiple sources

Everything else is a plugin: file discovery, content loading, markdown rendering, templating, feed generation, asset copying, etc.

### 2. Plugins Are First-Class

A single file should be able to:
- Hook into any build stage
- Extend the configuration schema
- Extend the post/content model
- Register new attributes on the core instance
- Add CLI commands

Plugins should be trivial for humans or AI agents to write.

### 3. Content Is Queryable Data

All content becomes structured data that can be:
- Filtered with expressions: `filter("published == True and date <= today")`
- Mapped to extract fields: `map("title", filter="'python' in tags")`
- Sorted, limited, and paginated
- Accessed from templates and markdown

### 4. Configuration Is Hierarchical and Mergeable

Configuration comes from (in order of precedence):
1. CLI arguments / environment variables
2. Local config file (`[name].toml` or `pyproject.toml`)
3. Global config file (`~/.[name].toml`)
4. Plugin defaults

### 5. Markdown Is Enhanced, Not Replaced

Writers use markdown with modern extensions:
- Frontmatter metadata
- Admonitions/callouts
- Tables
- Code blocks with syntax highlighting
- Template expressions within content
- Wikilinks / internal links

### 6. The Build Is Deterministic and Cacheable

Given the same inputs, produce the same outputs. Cache expensive operations (markdown parsing, image processing) keyed by content hash.

---

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                         CLI                                  │
├─────────────────────────────────────────────────────────────┤
│                    Core Orchestrator                         │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │  Lifecycle  │  │   Plugin    │  │    Data     │         │
│  │  Manager    │  │   Manager   │  │   Access    │         │
│  └─────────────┘  └─────────────┘  └─────────────┘         │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │   Config    │  │    Cache    │  │  Attribute  │         │
│  │   Loader    │  │   Manager   │  │   Storage   │         │
│  └─────────────┘  └─────────────┘  └─────────────┘         │
├─────────────────────────────────────────────────────────────┤
│                       Plugins                                │
│  ┌──────┐ ┌──────┐ ┌──────┐ ┌──────┐ ┌──────┐ ┌──────┐    │
│  │ glob │ │ load │ │render│ │feeds │ │ save │ │ ...  │    │
│  └──────┘ └──────┘ └──────┘ └──────┘ └──────┘ └──────┘    │
└─────────────────────────────────────────────────────────────┘
```

---

## Core Interface

The central orchestrator that plugins interact with.

### Properties

| Property | Type | Description |
|----------|------|-------------|
| `config` | Config | Merged configuration object |
| `posts` | List[Post] | All loaded content items |
| `files` | List[Path] | Discovered content files |
| `cache` | Cache | Persistent key-value cache |

### Query Methods

| Method | Signature | Description |
|--------|-----------|-------------|
| `filter` | `(expr: string) -> List[Post]` | Query posts with expression |
| `map` | `(field: string, filter?: string, sort?: string, reverse?: bool) -> List[Any]` | Extract field from filtered posts |
| `first` | `(filter?: string, sort?: string) -> Post?` | First matching post |
| `last` | `(filter?: string, sort?: string) -> Post?` | Last matching post |
| `one` | `(expr: string) -> Post` | Exactly one match (error if 0 or >1) |

### Lifecycle Methods

| Method | Description |
|--------|-------------|
| `register(plugin)` | Register a plugin |
| `run(stage?)` | Run build up to optional stage |

### Attribute Storage

For plugin communication:

| Method | Description |
|--------|-------------|
| `set(key, value)` | Store a value |
| `get(key) -> Any` | Retrieve a value |
| `has(key) -> bool` | Check if key exists |

### Model Registration

Called during `config_model` / `post_model` stages:

| Method | Description |
|--------|-------------|
| `register_config_model(model)` | Add config schema fragment |
| `register_post_fields(fields)` | Add post schema fragment |

---

## Lifecycle Stages

The build runs through ordered stages. Plugins implement hooks for stages they care about.

**This spec offers three lifecycle variants.** Choose based on your implementation language and plugin ecosystem needs:

| Variant | Stages | Best For |
|---------|--------|----------|
| **Full** | 13 | Python, Ruby - maximum plugin flexibility |
| **Standard** | 9 | TypeScript, Go - balanced flexibility/simplicity |
| **Minimal** | 6 | Rust, embedded - maximum simplicity |

See [LIFECYCLE.md](./LIFECYCLE.md) for detailed variant documentation and selection guide.

### Full Lifecycle (13 Stages) - Reference

| Stage | Purpose | Example Plugin Actions |
|-------|---------|----------------------|
| `config_model` | Register configuration schema fragments | Add `feeds` config section |
| `post_model` | Register content model fragments | Add `reading_time` field |
| `create_models` | Merge all model fragments | Core merges schemas |
| `load_config` | Load and validate configuration | Core loads from files |
| `configure` | Plugin initialization | Set up markdown parser |
| `validate_config` | Final config validation | Check required fields |
| `glob` | Discover content files | Find `**/*.md` files |
| `load` | Parse files into content objects | Parse frontmatter + content |
| `pre_render` | Process content before rendering | Expand template expressions |
| `render` | Convert content to HTML | Markdown → HTML |
| `post_render` | Post-process rendered content | Add heading anchors |
| `save` | Write output files | Write HTML to disk |
| `teardown` | Cleanup | Close connections |

### Standard Lifecycle (9 Stages)

| Stage | Purpose |
|-------|---------|
| `configure` | Load config, init plugins, set up resources |
| `validate` | Validate configuration |
| `glob` | Discover content files |
| `load` | Parse files into posts |
| `transform` | Pre-render processing |
| `render` | Markdown → HTML + templates |
| `collect` | Build feeds, series, navigation |
| `write` | Write output files |
| `cleanup` | Release resources |

### Minimal Lifecycle (6 Stages)

| Stage | Purpose |
|-------|---------|
| `init` | Load config, validate, init plugins |
| `load` | Find files and parse into posts |
| `process` | All content transformation and rendering |
| `build` | Build feeds, collections, final HTML |
| `output` | Write files and cleanup |

### Stage Guarantees

- Stages run in order
- All plugins for a stage complete before the next stage starts
- Within a stage, plugins can specify ordering (`tryfirst`, `trylast`)
- Plugins can access attributes registered by earlier stages

---

## Concurrency

### Post Processing

- `pre_render`, `render`, `post_render` stages SHOULD process posts concurrently
- Use worker pool pattern with configurable concurrency limit
- Default: number of CPU cores

### Thread Safety

- Plugins MUST NOT modify shared state without synchronization
- Each post can be processed independently
- `core.posts` list is read-only during concurrent stages

### Configuration

```toml
[name]
concurrency = 4  # 0 = auto (CPU count)
```

---

## Configuration

Configuration is namespaced under the tool name and supports multiple file formats (TOML, YAML, JSON) from multiple locations.

**Resolution Order (highest to lowest precedence):**
1. CLI arguments (`--output-dir public`)
2. Environment variables (`[NAME]_[SECTION]_[KEY]`)
3. Local config file (`./[name].toml`)
4. Global config file (`~/.config/[name]/config.toml`)
5. Plugin defaults

### Core Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `output_dir` | Path | `"output"` | Build output directory |
| `url` | URL? | null | Site base URL |
| `title` | string? | null | Site title |
| `description` | string? | null | Site description |
| `hooks` | List[str] | `["default"]` | Plugins to load |
| `disabled_hooks` | List[str] | `[]` | Plugins to exclude |
| `concurrency` | int | 0 | Worker count (0 = auto) |

### Example Configuration

```toml
[name]
output_dir = "public"
url = "https://example.com"
title = "My Site"
hooks = ["default"]

[name.glob]
patterns = ["posts/**/*.md", "pages/*.md"]

[name.markdown]
extensions = ["tables", "admonitions"]

[[name.feeds]]
slug = "blog"
title = "Blog"
filter = "published == True"
sort = "date"
reverse = true
```

**For complete configuration documentation, see [CONFIG.md](./CONFIG.md).**

---

## Data Access Layer

Enables querying content as structured data.

### Filter Expression Syntax

```
# Boolean fields
"published == True"
"draft == False"

# Comparisons
"date <= today"
"word_count > 1000"

# String operations
"'python' in tags"
"title.startswith('How to')"

# Compound expressions
"published == True and date <= today"
"status == 'draft' or 'wip' in tags"

# Negation
"not skip"
"not draft"
```

### Built-in Variables

| Variable | Description |
|----------|-------------|
| `today` | Current date |
| `now` | Current datetime |
| `True` / `False` | Boolean literals |

---

## Cache Manager

Persistent caching for expensive operations.

### Interface

```
cache.get(key) -> value | None
cache.set(key, value, expire?) -> None
cache.delete(key) -> None
cache.get_or_set(key, factory_fn) -> value
```

### Cache Keys

```
post:{path}:mtime     → file modification time
post:{path}:hash      → content hash
post:{path}:rendered  → cached article_html
template:{path}:mtime → template modification time
```

### Location

`.[name].cache/` in project root

---

## Standard Plugins

Implementations SHOULD provide these plugins:

### File Discovery (`glob`)
- Find content files matching patterns
- Respect `.gitignore` optionally
- Configuration: `glob_patterns`, `use_gitignore`

### Content Loading (`load`)
- Parse frontmatter and content
- Create Post objects
- Handle encoding

### Markdown Rendering (`render_markdown`)
- Convert markdown to HTML
- Support extensions: tables, code blocks, admonitions
- Configuration: `markdown_extensions`

### Template Processing (`templates`)
- Render posts through templates
- Provide template context: `post`, `config`, `core`
- Support template inheritance

### Feed Generation (`feeds`)
- Create filtered/sorted post collections
- Generate index pages
- Support pagination
- Configuration: `feeds` list with `filter`, `sort`, `template`

### HTML Publishing (`publish_html`)
- Write rendered HTML to output directory
- Create directory structure

### Asset Copying (`copy_assets`)
- Copy static files to output
- Preserve directory structure
- Configuration: `assets_dir`, `ignore`

#### Asset Fingerprinting

For cache busting, assets can be fingerprinted with content hashes:

**Configuration:**
```toml
[name.assets]
dir = "static"
fingerprint = true                    # Enable fingerprinting
fingerprint_algorithm = "sha256"      # or "md5"
fingerprint_length = 8                # Characters of hash to use
exclude_fingerprint = ["robots.txt", "favicon.ico"]  # Files to skip
```

**How It Works:**
1. Read asset file content
2. Compute hash of content
3. Rename file: `style.css` → `style.a1b2c3d4.css`
4. Generate manifest mapping original → fingerprinted names
5. Provide `asset_url()` function/filter for templates

**Manifest Format:**
```json
{
  "css/style.css": "css/style.a1b2c3d4.css",
  "js/app.js": "js/app.e5f6g7h8.js",
  "images/logo.png": "images/logo.i9j0k1l2.png"
}
```

**Template Usage:**
```jinja2
<link rel="stylesheet" href="{{ 'css/style.css' | asset_url }}">
<script src="{{ 'js/app.js' | asset_url }}"></script>
```

**Output:**
```html
<link rel="stylesheet" href="/css/style.a1b2c3d4.css">
<script src="/js/app.js.e5f6g7h8.js"></script>
```

**Implementation Notes:**
- Only fingerprint on production builds (not during `serve`)
- Store manifest in cache for incremental builds
- Only re-fingerprint assets that changed
- Support CDN prefix: `{{ 'css/style.css' | asset_url(cdn=True) }}`

### RSS/Atom Generation

RSS/Atom feeds are generated through the **feeds** plugin, not a separate plugin. Any feed can output RSS, Atom, JSON, or other formats by enabling them in the feed configuration.

See [FEEDS.md](./FEEDS.md) for complete feed format configuration.

### Sitemap Generation

Sitemaps are generated through the **feeds** plugin. Create a feed with `sitemap = true` format to generate a sitemap for that collection of posts.

See [FEEDS.md](./FEEDS.md) for sitemap format configuration.

---

## Plugin Resolution

When loading plugins from the `hooks` config:

1. **"default"** - Expands to the built-in plugin set:
   - glob, load, render_markdown, jinja_md, templates, feeds, publish_html, publish_feeds, copy_assets

2. **Built-in name** (e.g., `"glob"`, `"render_markdown"`)
   - Resolves to internal plugins

3. **Full module path** (language-specific)
   - Python: `"my_package.plugins.custom"`
   - TypeScript: `"my-package/plugins/custom"`
   - Go: `"github.com/user/ssg-plugin"`
   - Rust: `"my_crate::plugins::custom"`

4. **Local path** (e.g., `"./plugins/myplugin"`)
   - Loads from project-relative path

### Example

```toml
[name]
hooks = [
    "default",                    # Built-in set
    "./plugins/reading_time",     # Local plugin
]
disabled_hooks = ["toc", "wikilinks"]  # Disable specific plugins
```

---

## CLI Specification

### `build`

Build the static site.

```bash
[name] build [flags]

Flags:
  -c, --config string   Config file path (default "[name].toml")
  -o, --output string   Output directory (overrides config)
  -v, --verbose         Verbose output
  --clean               Remove output dir before build
  --drafts              Include draft posts
  --future              Include future-dated posts
```

### `serve`

Start development server with live reload.

```bash
[name] serve [flags]

Flags:
  -p, --port int        Server port (default 3000)
  -H, --host string     Bind address (default "localhost")
  --no-reload           Disable live reload
  --no-watch            Disable file watching
  --open                Open browser automatically
```

### `new <path>`

Create a new post with frontmatter template.

```bash
[name] new posts/my-new-post.md [flags]

Flags:
  -t, --title string    Post title
  --edit                Open in $EDITOR after creation
```

### `validate`

Validate configuration and content without building.

```bash
[name] validate [flags]

Flags:
  --strict              Treat warnings as errors
```

### Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | Build error |
| 2 | Config error |
| 3 | Validation error |

---

## Development Server

### Features

- HTTP server for previewing the built site
- File watching for automatic rebuilds
- Live reload for browser refresh

### Configuration

```toml
[name.serve]
port = 3000
host = "localhost"
livereload = true
open_browser = false
debounce_ms = 100  # wait before rebuilding after file change

# Watch patterns (defaults shown)
[name.serve.watch]
content = ["**/*.md"]                    # Post files
config = ["[name].toml", "pyproject.toml", "*.yaml", "*.json"]
templates = ["templates/**/*.html", "templates/**/*.jinja2"]
assets = ["static/**/*"]
plugins = ["plugins/**/*.py", "plugins/**/*.ts"]  # Local plugins
```

### Watch Patterns and Cache Busting

File changes trigger different rebuild behaviors based on what changed:

```
┌─────────────────────────────────────────────────────────────────────┐
│                    WATCH PATTERN CACHE INVALIDATION                  │
├─────────────────────────────────────────────────────────────────────┤
│  CONTENT FILES (*.md)                                                │
│    ├─ Invalidate: Post cache for changed file only                  │
│    ├─ Rebuild: Changed posts + feeds containing them                │
│    └─ Action: Incremental rebuild                                   │
│                                                                      │
│  CONFIG FILES ([name].toml, etc.)                                    │
│    ├─ Invalidate: ALL caches (config affects everything)            │
│    ├─ Rebuild: Full site                                            │
│    └─ Action: Full rebuild                                          │
│                                                                      │
│  TEMPLATE FILES (templates/**/*.html)                                │
│    ├─ Invalidate: Template cache + all posts using that template    │
│    ├─ Rebuild: All posts that use the changed template              │
│    └─ Action: Smart rebuild (template dependency tracking)          │
│                                                                      │
│  ASSET FILES (static/**/*)                                           │
│    ├─ Invalidate: Asset manifest only                               │
│    ├─ Rebuild: None (just copy)                                     │
│    └─ Action: Copy changed assets                                   │
│                                                                      │
│  PLUGIN FILES (plugins/**/*.py)                                      │
│    ├─ Invalidate: ALL caches (plugin may affect any stage)          │
│    ├─ Rebuild: Full site                                            │
│    └─ Action: Full rebuild with plugin reload                       │
└─────────────────────────────────────────────────────────────────────┘
```

### Cache Key Dependencies

Cache entries are keyed to include all dependencies:

| Cache Entry | Key Components |
|-------------|----------------|
| Post content | `post:{path}:{content_hash}` |
| Post rendered | `post:{path}:{content_hash}:{template_hash}:{config_hash}` |
| Template | `template:{path}:{content_hash}:{parent_hash}` |
| Feed | `feed:{slug}:{filter_hash}:{post_hashes}:{config_hash}` |
| Asset | `asset:{path}:{content_hash}` |

When any component of the key changes, the cache entry is invalidated.

### Template Dependency Tracking

Templates can extend or include other templates. The server tracks these dependencies:

```python
# Example dependency graph
base.html → []  # No dependencies
post.html → [base.html]  # Extends base
feed.html → [base.html, partials/card.html]  # Extends base, includes card
partials/card.html → []  # No dependencies

# When base.html changes:
# Invalidate: base.html, post.html, feed.html (all dependents)
```

### Rebuild Triggers

| Change Type | Action | Cache Impact |
|-------------|--------|--------------|
| Content file (*.md) | Incremental rebuild | Post + downstream feeds |
| Template file | Smart rebuild | Posts using template |
| Config file | Full rebuild | Everything |
| Static asset | Copy only | Asset manifest |
| Plugin file | Full rebuild | Everything |
| Theme file | Full rebuild | Everything |

### Live Reload

- Inject reload script into HTML pages during serve
- Use WebSocket or Server-Sent Events
- Trigger reload after successful build

---

## Output Structure

```
[output_dir]/
├── index.html              # Home page (from feed)
├── [post-slug]/
│   └── index.html          # Individual posts
├── [feed-slug]/
│   ├── index.html          # Feed page 1
│   ├── page/2/index.html   # Feed page 2
│   └── ...
├── rss.xml                 # RSS feed
├── sitemap.xml             # Sitemap
└── [assets]/               # Copied static files
```

---

## Error Handling

### Error Types

Errors should include contextual information:
- `stage` - which lifecycle stage
- `plugin` - which plugin caused it
- `path` - file path if applicable
- `line` - line number if applicable
- `message` - human-readable description

### Standard Errors

| Error | When |
|-------|------|
| `FrontmatterParseError` | Invalid YAML in frontmatter |
| `FilterExpressionError` | Invalid filter syntax |
| `TemplateNotFoundError` | Template file doesn't exist |
| `TemplateSyntaxError` | Invalid template syntax |
| `ConfigValidationError` | Invalid configuration |
| `PluginNotFoundError` | Plugin couldn't be loaded |
| `CircularTemplateError` | Template inheritance cycle |

### Graceful Degradation

- Missing optional fields use defaults
- Unknown frontmatter fields are preserved
- Plugin errors log but don't halt build (configurable)

---

## Implementation Checklist

An implementation is complete when it:

### Core
- [ ] Runs lifecycle stages in order
- [ ] Supports concurrent post processing
- [ ] Loads plugins from configuration
- [ ] Resolves plugin names to implementations
- [ ] Merges config from multiple sources
- [ ] Provides thread-safe attribute storage

### Content
- [ ] Discovers content files via glob patterns
- [ ] Parses YAML frontmatter
- [ ] Handles unicode filenames
- [ ] Supports `filter()` with expressions
- [ ] Supports `map()` for field extraction

### Rendering
- [ ] Renders markdown to HTML
- [ ] Processes templates
- [ ] Supports template inheritance
- [ ] Detects circular template inheritance

### Output
- [ ] Generates feeds with filtering/sorting/pagination
- [ ] Writes output to configured directory
- [ ] Copies static assets
- [ ] Generates RSS feed
- [ ] Generates sitemap

### Caching
- [ ] Caches expensive operations
- [ ] Supports incremental builds
- [ ] Invalidates cache on config change

### CLI
- [ ] `build` command with all flags
- [ ] `serve` command with live reload
- [ ] `new` command for post creation
- [ ] `validate` command for checking

### Plugins
- [ ] Allows single-file plugin creation
- [ ] Plugin can extend config schema
- [ ] Plugin can extend post model
- [ ] Plugin can add CLI commands
- [ ] Plugin errors don't crash build (configurable)

### Testing
- [ ] Passes all tests in `tests.yaml`
- [ ] Handles edge cases (large files, unicode, etc.)

---

## See Also

- [IMPLEMENTATION.md](./IMPLEMENTATION.md) - Step-by-step implementation guide
- [INSTALL.md](./INSTALL.md) - Quick start with tech choices
- [CONFIG.md](./CONFIG.md) - Complete configuration documentation
- [THEMES.md](./THEMES.md) - Theming and customization
- [LIFECYCLE.md](./LIFECYCLE.md) - Detailed lifecycle specification
- [FEEDS.md](./FEEDS.md) - Feed system specification (core differentiator)
- [DEFAULT_PLUGINS.md](./DEFAULT_PLUGINS.md) - Built-in plugins
- [PLUGINS.md](./PLUGINS.md) - Plugin development guide
- [DATA_MODEL.md](./DATA_MODEL.md) - Content model specification
- [CONTENT.md](./CONTENT.md) - Markdown processing specification
- [TEMPLATES.md](./TEMPLATES.md) - Template system specification
- [tests.yaml](./tests.yaml) - Test cases for implementations
