# Lifecycle Stages Specification

The build process runs through ordered stages. Each stage is a hook point where plugins can participate.

This specification offers **three lifecycle variants** to suit different implementation needs. Implementations MUST choose one variant and document which they use.

---

## Lifecycle Variants

### Choosing a Variant

| Variant | Stages | Best For |
|---------|--------|----------|
| **Full (13 stages)** | 13 | Plugin ecosystems, maximum flexibility, dynamic languages |
| **Standard (9 stages)** | 9 | Most implementations, good balance of flexibility and simplicity |
| **Minimal (6 stages)** | 6 | Simple implementations, compiled languages, embedded use |

**Decision guide:**

```
┌─────────────────────────────────────────────────────────────────────┐
│                    WHICH LIFECYCLE TO CHOOSE?                        │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  Do you need runtime schema extension by plugins?                    │
│    YES → Use FULL (13 stages)                                        │
│    NO  ↓                                                             │
│                                                                      │
│  Do you need separate pre/post render hooks?                         │
│    YES → Use STANDARD (9 stages)                                     │
│    NO  ↓                                                             │
│                                                                      │
│  Is simplicity/performance the priority?                             │
│    YES → Use MINIMAL (6 stages)                                      │
│                                                                      │
├─────────────────────────────────────────────────────────────────────┤
│  LANGUAGE RECOMMENDATIONS                                            │
│                                                                      │
│  Python/Ruby    → Full or Standard (dynamic, plugin-friendly)        │
│  TypeScript/JS  → Standard (good balance)                            │
│  Go             → Standard or Minimal (simplicity, compilation)      │
│  Rust           → Minimal or Standard (compile-time schemas)         │
└─────────────────────────────────────────────────────────────────────┘
```

---

## Variant A: Full Lifecycle (13 Stages)

**Best for:** Python, Ruby, or any implementation prioritizing plugin ecosystem extensibility.

**Characteristics:**
- Maximum hook points for plugins
- Runtime schema composition (plugins extend config/post models)
- Fine-grained error handling per stage
- Higher implementation complexity

```
┌─────────────────────────────────────────────────────────────────────┐
│                        CONFIGURATION PHASE                          │
├─────────────────────────────────────────────────────────────────────┤
│  config_model → post_model → create_models → load_config            │
│       │              │             │              │                 │
│       ▼              ▼             ▼              ▼                 │
│  [register      [register     [merge all    [load from             │
│   config         post          models]       files]                 │
│   schemas]       schemas]                                           │
├─────────────────────────────────────────────────────────────────────┤
│  configure → validate_config                                        │
│       │              │                                              │
│       ▼              ▼                                              │
│  [plugin         [validate                                          │
│   init]           config]                                           │
├─────────────────────────────────────────────────────────────────────┤
│                         CONTENT PHASE                               │
├─────────────────────────────────────────────────────────────────────┤
│  glob → load → pre_render → render → post_render                    │
│    │      │         │          │          │                         │
│    ▼      ▼         ▼          ▼          ▼                         │
│  [find  [parse   [process   [markdown  [enhance                     │
│   files] content] before]    → HTML]    HTML]                       │
├─────────────────────────────────────────────────────────────────────┤
│                         OUTPUT PHASE                                │
├─────────────────────────────────────────────────────────────────────┤
│  save → teardown                                                    │
│    │        │                                                       │
│    ▼        ▼                                                       │
│  [write  [cleanup]                                                  │
│   files]                                                            │
└─────────────────────────────────────────────────────────────────────┘
```

**Stages:**
1. `config_model` - Register config schema fragments
2. `post_model` - Register post schema fragments
3. `create_models` - Merge schemas into final models
4. `load_config` - Load configuration from files
5. `configure` - Plugin initialization
6. `validate_config` - Validate final configuration
7. `glob` - Discover content files
8. `load` - Parse files into posts
9. `pre_render` - Process content before rendering
10. `render` - Convert markdown to HTML
11. `post_render` - Enhance HTML, build feeds/navigation
12. `save` - Write output files
13. `teardown` - Cleanup resources

---

## Variant B: Standard Lifecycle (9 Stages)

**Best for:** TypeScript, Go, or balanced implementations.

**Characteristics:**
- Good plugin flexibility without over-engineering
- Config/post schemas defined at compile/init time (not runtime-composed)
- Separates collection-building from HTML post-processing
- Recommended for most implementations

```
┌─────────────────────────────────────────────────────────────────────┐
│                        CONFIGURATION PHASE                          │
├─────────────────────────────────────────────────────────────────────┤
│  configure → validate                                               │
│       │           │                                                 │
│       ▼           ▼                                                 │
│  [load config, [validate                                            │
│   init plugins] config]                                             │
├─────────────────────────────────────────────────────────────────────┤
│                         CONTENT PHASE                               │
├─────────────────────────────────────────────────────────────────────┤
│  glob → load → transform → render → collect                         │
│    │      │        │          │         │                           │
│    ▼      ▼        ▼          ▼         ▼                           │
│  [find  [parse  [pre-proc  [markdown [feeds,                        │
│   files] posts]  content]   → HTML]   nav]                          │
├─────────────────────────────────────────────────────────────────────┤
│                         OUTPUT PHASE                                │
├─────────────────────────────────────────────────────────────────────┤
│  write → cleanup                                                    │
│    │        │                                                       │
│    ▼        ▼                                                       │
│  [output  [close                                                    │
│   files]   resources]                                               │
└─────────────────────────────────────────────────────────────────────┘
```

**Stages:**
1. `configure` - Load config, initialize plugins, set up resources
2. `validate` - Validate configuration
3. `glob` - Discover content files
4. `load` - Parse files into posts
5. `transform` - Pre-render processing (jinja-md, descriptions, etc.)
6. `render` - Convert markdown to HTML, apply templates
7. `collect` - Build feeds, series, prev/next, navigation
8. `write` - Write all output files
9. `cleanup` - Release resources

Title lifecycle ordering is part of the standard transform contract:

```text
auto_title → inline_titles → semantic-title consumers
```

`inline_titles` derives `title_html` and `title_text` before plugins that
consume titles for metadata, feeds, search, or title-derived slugs. Load may
leave an automatic fallback slug pending until this transform completes.

**Key differences from Full:**
- `configure` combines model creation + config loading + plugin init
- `transform` replaces `pre_render` (clearer name)
- `render` includes template application
- `collect` is new - separates feed/navigation building from HTML processing
- `write` + `cleanup` replace `save` + `teardown` (clearer names)

**Stage mapping (Full → Standard):**

| Full Stage | Standard Stage |
|------------|----------------|
| config_model | configure |
| post_model | configure |
| create_models | configure |
| load_config | configure |
| configure | configure |
| validate_config | validate |
| glob | glob |
| load | load |
| pre_render | transform |
| render | render |
| post_render | render + collect |
| save | write |
| teardown | cleanup |

### Manager post lookup index

The Go standard-lifecycle implementation exposes one shared `PostIndex` for
post lookups. The index is lazy: `Manager.PostIndex()` builds it from the
manager's current `Posts()` collection on first access and then returns the
cached index. It is not limited to a call after `load`; callers can request it
at any lifecycle stage.

The index provides these keys:

| Map | Key behavior |
| --- | --- |
| `BySlug` | Lowercase `Post.Slug` and lowercase string aliases |
| `BySlugified` | Slugified `Post.Slug` when distinct from its lowercase key, plus slugified aliases |
| `ByHref` | Exact, non-empty `Post.Href` |
| `ByPath` | Exact, non-empty `Post.Path` |

`LookupBySlug` checks the lowercase slug, then the slugified key, then the
slugified key in `BySlug`. Main slug, slugified slug, href, and path entries use
last-post-wins behavior when duplicate map keys occur. Alias entries do not
replace an existing slug or slugified entry. A nil index or missing key returns
nil.

`SetPosts`, `AddPost`, and `Reset` invalidate the manager's cached index while
holding the manager lock. Synthetic posts that are visible to later stages
must therefore be added with `AddPost`. Mutating `Slug`, `Href`, `Path`, or
alias data through an existing `*Post` cannot be observed by `Manager`; the
caller must call `PostIndex().Refresh(m)` before using that index. The index is
also not a replacement for filtered, cache-owned, validation, or aggregation
maps whose source collection or duplicate semantics differ.

---

## Variant C: Minimal Lifecycle (6 Stages)

**Best for:** Rust, embedded use, or maximum simplicity.

**Characteristics:**
- Minimal hook points
- Schemas defined at compile time
- Single content processing stage
- Easiest to implement and reason about

```
┌─────────────────────────────────────────────────────────────────────┐
│                        CONFIGURATION PHASE                          │
├─────────────────────────────────────────────────────────────────────┤
│  init                                                               │
│    │                                                                │
│    ▼                                                                │
│  [load config, validate, init plugins]                              │
├─────────────────────────────────────────────────────────────────────┤
│                         CONTENT PHASE                               │
├─────────────────────────────────────────────────────────────────────┤
│  load → process → build                                             │
│    │       │         │                                              │
│    ▼       ▼         ▼                                              │
│  [find & [transform [feeds,                                         │
│   parse]  & render]  nav, HTML]                                     │
├─────────────────────────────────────────────────────────────────────┤
│                         OUTPUT PHASE                                │
├─────────────────────────────────────────────────────────────────────┤
│  output                                                             │
│    │                                                                │
│    ▼                                                                │
│  [write files, cleanup]                                             │
└─────────────────────────────────────────────────────────────────────┘
```

**Stages:**
1. `init` - Load config, validate, initialize everything
2. `load` - Find files and parse into posts
3. `process` - All content transformation and rendering
4. `build` - Build feeds, collections, navigation, final HTML
5. `output` - Write files and cleanup

**Key differences:**
- Single `init` stage for all setup
- `load` combines glob + load
- `process` combines transform + render
- `build` handles collections and final assembly
- `output` combines write + cleanup

**Stage mapping (Full → Minimal):**

| Full Stage | Minimal Stage |
|------------|---------------|
| config_model | init |
| post_model | init |
| create_models | init |
| load_config | init |
| configure | init |
| validate_config | init |
| glob | load |
| load | load |
| pre_render | process |
| render | process |
| post_render | build |
| save | output |
| teardown | output |

---

## Variant Comparison

### Hook Points

| Capability | Full | Standard | Minimal |
|------------|------|----------|---------|
| Extend config schema at runtime | ✓ | ✗ | ✗ |
| Extend post schema at runtime | ✓ | ✗ | ✗ |
| Separate config validation | ✓ | ✓ | ✗ |
| Pre-render hook | ✓ | ✓ | ✗ |
| Post-render hook | ✓ | ✗ | ✗ |
| Separate collection building | ✗ | ✓ | ✓ |
| Separate cleanup stage | ✓ | ✓ | ✗ |

### Implementation Complexity

| Aspect | Full | Standard | Minimal |
|--------|------|----------|---------|
| Lines of core code | ~2000 | ~1200 | ~600 |
| Plugin API surface | Large | Medium | Small |
| Mental model | Complex | Moderate | Simple |
| Testing surface | 13 stages | 9 stages | 6 stages |

### Plugin Capabilities by Variant

**Full (13 stages):**
```python
# Plugin can extend the Post model at runtime
@hook_impl
def post_model(core):
    core.post_models.append(MyFields)

# Plugin can hook pre and post render separately
@hook_impl
def pre_render(core):
    # Before markdown processing
    pass

@hook_impl  
def post_render(core):
    # After markdown, before save
    pass
```

**Standard (9 stages):**
```typescript
// Plugin registers with fixed schema (defined at init)
export const plugin: Plugin = {
  name: 'my-plugin',

  transform(core) {
    // Pre-render processing
  },

  collect(core) {
    // Build collections after render
  }
}
```

**Minimal (6 stages):**
```rust
// Plugin implements trait with fixed hooks
impl Plugin for MyPlugin {
    fn process(&self, core: &mut Core) {
        // All content processing here
    }

    fn build(&self, core: &mut Core) {
        // Collection building here
    }
}
```

---

## Stage Overview (Full Variant - Reference)

The following detailed stage documentation uses the **Full (13 stage)** variant. Implementations using Standard or Minimal variants should map stages according to the tables above.

---

## Stage 1: `config_model`

**Purpose:** Register configuration schema fragments.

**When:** First stage, before any configuration is loaded.

**What plugins do:**
- Append Pydantic/Zod/etc. models to `core.config_models`
- Define configuration sections with defaults

**Core actions:**
- Initialize `core.config_models` list
- Call all `config_model` hooks

**Example:**
```python
class MyPluginConfig(BaseModel):
    enabled: bool = True
    threshold: int = 100

class Config(BaseModel):
    my_plugin: MyPluginConfig = MyPluginConfig()

@hook_impl
def config_model(core):
    core.config_models.append(Config)
```

**After this stage:**
- `core.config_models` contains all config schema fragments

---

## Stage 2: `post_model`

**Purpose:** Register post/content model fragments.

**When:** After config models are registered.

**What plugins do:**
- Append field definitions to `core.post_models`
- Define post attributes with types and defaults

**Core actions:**
- Initialize `core.post_models` list
- Call all `post_model` hooks

**Example:**
```python
class ReadingTimeFields(BaseModel):
    reading_time: int = 0
    word_count: int = 0

@hook_impl
def post_model(core):
    core.post_models.append(ReadingTimeFields)
```

**After this stage:**
- `core.post_models` contains all post schema fragments

---

## Stage 3: `create_models`

**Purpose:** Merge all schema fragments into final models.

**When:** After all model hooks have run.

**What plugins do:**
- Usually nothing (core handles this)
- Advanced: modify merged models

**Core actions:**
- Merge `core.config_models` into `core.Config`
- Merge `core.post_models` into `core.Post`
- Create factory functions for model instantiation

**After this stage:**
- `core.Config` - the merged configuration class
- `core.Post` - the merged post class

---

## Stage 4: `load_config`

**Purpose:** Load configuration from files and environment.

**When:** After models are created.

**What plugins do:**
- Usually nothing (core handles this)
- Advanced: provide custom config sources

**Core actions:**
- Find config files (local, global)
- Parse and merge configuration
- Apply environment variable overrides
- Instantiate `core.config` from `core.Config`

**Resolution order (highest precedence first):**
1. Environment variables (`TOOL_SECTION_KEY`)
2. CLI arguments
3. Local config file (`tool.toml` or `pyproject.toml`)
4. Global config file (`~/.tool.toml`)
5. Model defaults

**After this stage:**
- `core.config` - the loaded configuration object

---

## Stage 5: `configure`

**Purpose:** Plugin initialization with access to configuration.

**When:** After configuration is loaded.

**What plugins do:**
- Initialize plugin state
- Set up external connections
- Create shared resources
- Register attributes on core

**Core actions:**
- Initialize cache at `.tool.cache/`
- Set up logging
- Call all `configure` hooks

**Example:**
```python
@hook_impl
@register_attr("markdown_parser")
def configure(core):
    from markdown_it import MarkdownIt
    core.markdown_parser = MarkdownIt()
```

**After this stage:**
- `core.cache` - persistent cache instance
- Plugin-registered attributes available

---

## Stage 6: `validate_config`

**Purpose:** Validate configuration after all plugins have configured.

**When:** After all plugins have initialized.

**What plugins do:**
- Validate plugin-specific configuration
- Check for required fields
- Emit warnings for deprecated options
- Warn when the root `license` key is missing so that local builds remind authors to declare or opt out.
- Raise errors for invalid configurations

**Core actions:**
- Call all `validate_config` hooks
- Aggregate validation errors

**Example:**
```python
@hook_impl
def validate_config(core):
    if core.config.my_plugin.enabled:
        if not core.config.my_plugin.api_key:
            raise ConfigError("my_plugin.api_key is required when enabled")
```

**After this stage:**
- Configuration is validated and ready to use

## Live Serve Feedback

The development server injects the validate stage status into every HTML page via `build_status`. When the root `license` field remains unset (and is not explicitly `false`), the validation warning is propagated as a `license_warning` payload so the live reload banner can show a toast reminding you to set a license key or opt out. Once a license string is configured the toast disappears, keeping the warning scoped to the local workflow until metadata is complete.

---

## Stage 7: `glob`

**Purpose:** Discover content files.

**When:** First content phase stage.

**What plugins do:**
- Find files matching patterns
- Filter based on gitignore (optional)
- Add files to `core.files`

**Core actions:**
- Initialize `core.files` list
- Call all `glob` hooks

**Example:**
```python
@hook_impl
def glob(core):
    from pathlib import Path

    for pattern in core.config.glob.glob_patterns:
        core.files.extend(Path().glob(pattern))
```

**After this stage:**
- `core.files` - list of Path objects to process

---

## Stage 8: `load`

**Purpose:** Parse files into content objects.

**When:** After files are discovered.

**What plugins do:**
- Read file contents
- Parse frontmatter
- Create Post objects
- Handle encoding issues

**Core actions:**
- Initialize `core.articles` / `core.posts` list
- Call all `load` hooks

**Example:**
```python
@hook_impl
def load(core):
    for path in core.files:
        content = path.read_text()
        frontmatter, body = parse_frontmatter(content)

        post = core.Post(
            path=path,
            content=body,
            **frontmatter
        )
        core.articles.append(post)
```

**After this stage:**
- `core.articles` / `core.posts` - list of Post objects

---

## Stage 9: `pre_render`

**Purpose:** Process content before markdown rendering.

**When:** After content is loaded, before rendering.

**What plugins do:**
- Expand template expressions in content (Jinja-in-Markdown)
- Process shortcodes
- Calculate derived fields (reading time, etc.)
- Set up prev/next links
- Auto-generate descriptions

**Example:**
```python
@hook_impl
def pre_render(core):
    for post in core.filter("jinja == True"):
        template = core.jinja_env.from_string(post.content)
        post.content = template.render(
            post=post,
            core=core,
            config=core.config
        )
```

**After this stage:**
- Post `content` fields are processed and ready for rendering

---

## Stage 10: `render`

**Purpose:** Convert content to HTML.

**When:** After pre-render processing.

**What plugins do:**
- Convert markdown to HTML
- Apply syntax highlighting
- Process admonitions
- Expand internal links

**Example:**
```python
@hook_impl
def render(core):
    for post in core.filter("not skip"):
        post.article_html = core.markdown_parser.render(post.content)
```

**After this stage:**
- `post.article_html` - rendered HTML content (without template)

**Which posts are rendered:** every post with `skip == false`. An empty body
is not a reason to skip; a published post whose Markdown body is empty still
gets a full page (title, metadata, template chrome) so that feed cards and
links pointing at it never 404. Only `skip` (set explicitly or by the
publishing rules) removes a post from rendering.

---

## Stage 11: `post_render`

**Purpose:** Post-process rendered HTML.

**When:** After markdown is rendered.

**What plugins do:**
- Add heading anchors
- Process images (lazy loading, etc.)
- Inject scripts/styles
- Apply templates to wrap content
- Build navigation elements

**Example:**
```python
@hook_impl
def post_render(core):
    for post in core.filter("not skip"):
        # Wrap in template
        template = core.jinja_env.get_template(post.template)
        post.html = template.render(
            post=post,
            body=post.article_html,
            config=core.config,
            core=core
        )
```

**After this stage:**
- `post.html` - final HTML ready to write

---

## Stage 12: `save`

**Purpose:** Write output files to disk.

**When:** After all rendering is complete.

**What plugins do:**
- Write post HTML files
- Write feed/index pages
- Write RSS/Atom feeds
- Write sitemap
- Copy static assets
- Generate service workers

**Example:**
```python
@hook_impl
def save(core):
    for post in core.filter("not skip"):
        output_path = core.config.output_dir / post.slug / "index.html"
        output_path.parent.mkdir(parents=True, exist_ok=True)
        output_path.write_text(post.html)
```

**After this stage:**
- All output files written to `core.config.output_dir`

---

## Stage 13: `teardown`

**Purpose:** Cleanup resources.

**When:** Final stage, after save.

**What plugins do:**
- Close database connections
- Flush caches
- Clean up temporary files
- Log final statistics

**Example:**
```python
@hook_impl
def teardown(core):
    core.cache.close()
    logger.info(f"Built {len(core.posts)} posts")
```

**After this stage:**
- Build complete
- All resources released

---

## Running Specific Stages

The core `run()` method accepts an optional stage parameter:

```python
# Run full build
core.run()

# Run up to (and including) a specific stage
core.run("render")  # Stops after render, doesn't save

# Run from current point to a stage
core.run("pre_render")  # Run through pre_render
core.run("render")      # Continue to render (won't re-run earlier stages)
```

**Use cases:**
- Testing: `core.run("pre_render")` to check content processing
- CLI commands: Run partial builds for specific operations
- Development: Incremental builds

---

## Stage Dependencies

Each stage has implicit dependencies on previous stages:

| Stage | Requires |
|-------|----------|
| `config_model` | (none) |
| `post_model` | `config_model` |
| `create_models` | `post_model` |
| `load_config` | `create_models` |
| `configure` | `load_config` |
| `validate_config` | `configure` |
| `glob` | `validate_config` |
| `load` | `glob` |
| `pre_render` | `load` |
| `render` | `pre_render` |
| `post_render` | `render` |
| `save` | `post_render` |
| `teardown` | `save` |

Running a stage automatically runs all preceding stages that haven't run yet.

---

## Hook Execution Order

Within a stage, hooks execute in this order:

1. Hooks marked `tryfirst=True` (in registration order)
2. Hooks with no modifier (in registration order)
3. Hooks marked `trylast=True` (in registration order)

**Example:**
```python
# These hooks are in the same stage
@hook_impl(tryfirst=True)
def render(core):  # Runs first
    pass

@hook_impl
def render(core):  # Runs second
    pass

@hook_impl(trylast=True)
def render(core):  # Runs last
    pass
```

---

## Error Handling by Stage

| Stage | Error Behavior |
|-------|----------------|
| `config_model` | Fatal - cannot proceed without models |
| `post_model` | Fatal |
| `create_models` | Fatal |
| `load_config` | Fatal - invalid config stops build |
| `configure` | Fatal by default, can be configured to warn |
| `validate_config` | Fatal - validation errors stop build |
| `glob` | Warn if no files found |
| `load` | Skip invalid files, log warnings |
| `pre_render` | Skip post on error, continue others |
| `render` | Skip post on error, continue others |
| `post_render` | Skip post on error, continue others |
| `save` | Log error, continue other files |
| `teardown` | Log errors, attempt all cleanup |

---

## Concurrent Post Processing

Stages that process posts (transform, render, collect) use a **bounded worker pool** for concurrent execution.

### Worker Pool Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                    BOUNDED WORKER POOL                               │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  ┌─────────┐    ┌───────────────────┐    ┌────────────────────┐     │
│  │ Posts   │───►│   Jobs Channel    │───►│  Worker Goroutines │     │
│  │ (N)     │    │   (buffered)      │    │  (Concurrency)     │     │
│  └─────────┘    └───────────────────┘    └────────────────────┘     │
│                                                                      │
│  Post count: unlimited                    Worker count: bounded      │
│                                           (default: CPU cores)       │
│                                                                      │
├─────────────────────────────────────────────────────────────────────┤
│  GUARANTEES:                                                         │
│  • Never spawns more than Concurrency() + k goroutines              │
│  • Constant memory overhead regardless of post count                 │
│  • Deterministic completion (all posts processed before return)      │
│  • Error aggregation (count + first error preserved)                 │
└─────────────────────────────────────────────────────────────────────┘
```

### Why Worker Pool (not Goroutine-per-Post)

The naive approach of spawning one goroutine per post with a semaphore:
- Creates N goroutines immediately (N = post count)
- Causes scheduler overhead for large sites (5k+ posts)
- Results in memory churn from goroutine stack allocations

The worker pool approach:
- Creates exactly `Concurrency()` goroutines
- Streams posts through a jobs channel
- Scales to any post count without additional overhead

### Configuration

```toml
[markata-go]
concurrency = 0  # 0 = auto (CPU cores), or explicit worker count
```

### Error Handling

- If a post fails to process, the error is recorded
- Processing continues for remaining posts
- Returns aggregated error: `"N posts failed to process; first error: ..."`

---

## Incremental Builds

Incremental builds only process content that has changed since the last build, dramatically improving build times for large sites.

### Change Detection

The system detects changes using multiple signals:

| Signal | Detection Method | Triggers |
|--------|-----------------|----------|
| Content file modified | File mtime or content hash | Rebuild that post |
| Content file added | File not in cache | Build new post |
| Content file deleted | File in cache but not on disk | Remove from output |
| Template modified | Template mtime | Rebuild all posts using that template |
| Config modified | Config file mtime or hash | Full rebuild |
| Plugin changed | Plugin file mtime | Full rebuild |

### Serve Incremental Rebuilds

The development server (`serve`) prioritizes deterministic rebuilds by default and a faster
incremental mode when `--fast` is enabled.

**Default (safe) mode:**
- Rebuilds are full and deterministic on every change event.
- Build cache still skips expensive per-post work (render/write) when inputs are unchanged.

**Fast mode (`serve --fast`):**
- File watcher events are normalized to content-relative paths.
- Only changed posts are reloaded from disk; the rest are restored from cache.
- The dependency graph expands the rebuild set to include transitive dependents.
- If dependency changes are detected (links, embeds), a second pass reprocesses affected posts.
- Changes outside the content directory (templates, config, assets) force a full rebuild.

`build --fast` is different from `serve --fast`. A normal build still rescans content globs on
every invocation so added, removed, and moved files are detected without clearing caches.
`--fast` on `build` only skips expensive non-essential work; it does not enable serve-style
incremental file discovery.

This ensures `serve` remains accurate by default, while `serve --fast` provides a
quick iteration loop without requiring config hacks.

### Cache Keys

Each post's cache entry includes:

```
post:{path}:
  content_hash: "sha256:abc123..."   # Hash of raw file content
  mtime: 1706123456                  # File modification time
  frontmatter_hash: "sha256:def..."  # Hash of parsed frontmatter
  dependencies:                      # Files this post depends on
    - templates/post.html
    - templates/base.html
    - partials/header.html
  rendered_hash: "sha256:ghi..."     # Hash of final HTML output
  article_html: "<p>...</p>"         # Cached rendered markdown
```

Incremental cleanup plugins may also persist global hashes for derived artifacts,
such as Tailwind utility manifests and the Pagefind search corpus, so expensive
whole-site rebuilds can be skipped when the effective output is unchanged.

### Dependency Tracking

Posts may depend on:

1. **Templates** - The template chain used to render
2. **Includes** - Partials included by templates
3. **Other posts** - Via wikilinks or queries in Jinja-in-Markdown
4. **Static assets** - Images or files referenced in content

When a dependency changes, all dependent posts are marked for rebuild.
Change detection and transitive dependent expansion MUST finish after Load and
before Transform. This ordering lets affected transforms and renderers reject
stale cached derivatives. Feeds that contain an affected post MUST also be
refreshed, even when that post's source file did not change.

### Rebuild Strategy

```
┌─────────────────────────────────────────────────────────────┐
│                    INCREMENTAL BUILD                         │
├─────────────────────────────────────────────────────────────┤
│  1. Load previous build cache                                │
│  2. Scan content directory for files                         │
│  3. For each file:                                           │
│     ├─ If new: mark for full processing                      │
│     ├─ If deleted: mark output for removal                   │
│     ├─ If modified: mark for rebuild                         │
│     └─ If unchanged: load from cache                         │
│  4. Check template/config changes                            │
│     ├─ If template changed: mark dependent posts             │
│     └─ If config changed: mark all posts                     │
│  5. Process only marked posts through render stages          │
│  6. Merge cached + newly rendered posts                      │
│  7. Run save stage (write changed files only)                │
│  8. Update cache with new state                              │
└─────────────────────────────────────────────────────────────┘
```

When an unchanged post is restored from the build cache, its cached full-page
HTML MUST be available to the write stage before output materialization. This
also applies when the output directory is new or was removed after the cache
was created. A published post with an empty Markdown body still has a full
page, so an empty `post.article_html` value MUST NOT prevent that cached page
from being written. If the full-page cache entry is unavailable, the post MUST
be rendered instead of being treated as successfully restored.

When a post is rendered or restored for publication, the publisher MUST leave
an existing regular output file untouched when its bytes already equal the
desired output. This preserves file identity for downstream consumers such as
CSS selector caching. Missing or different files are written normally. A
comparison error MUST fall through to the normal write path, and non-regular
paths such as symlinks retain the normal write behavior. This comparison does
not change the existing mode behavior for files that already exist.

### Cache Invalidation Rules

| Change | Invalidation Scope |
|--------|-------------------|
| Single post content | That post only |
| Linked post content | That post and its transitive dependents |
| Post frontmatter (no template change) | That post only |
| Post template assignment | That post only |
| Template file | All posts using that template |
| Base template | All posts (cascades through inheritance) |
| Partial/include | All posts whose templates use it |
| Feed membership change | All posts in the affected feed |
| Config change | Full rebuild |
| Plugin code change | Full rebuild |
| Markdown extension settings | Full rebuild |

**Feed Membership Tracking:**

Posts that belong to a feed sidebar (configured via `[components.feed_sidebar]`) track a
`FeedMembershipHash` -- a hash of the sorted slugs of all co-members in the feed. When a
post is added to or removed from a tag-based feed, this hash changes and all existing
members are rebuilt with the updated sidebar. This ensures sidebar navigation stays
consistent even when only one post's content changes.

### Implementation Notes

**Hash vs Mtime:**
- Use content hash for correctness (handles git checkout, rsync, etc.)
- Use mtime as fast-path optimization (skip hash if mtime unchanged)
- Always fall back to hash comparison when mtime differs

**Query Dependencies:**
Posts with Jinja-in-Markdown queries like:
```jinja2
{% for p in core.filter("'python' in tags")[:5] %}
```
Depend on the query result set. Track which posts match the query; if that set changes, rebuild the querying post.

**Atomic Cache Updates:**
- Write cache updates atomically (write to temp, then rename)
- On failed build, don't update cache
- Store cache format version; invalidate on version mismatch

### Configuration

```toml
[name]
# Enable/disable incremental builds
incremental = true

# Force full rebuild (ignore cache)
# Useful for CI or when cache may be stale
force_rebuild = false

# Cache directory
cache_dir = ".name.cache"

# What to use for change detection
change_detection = "hash"  # "hash", "mtime", or "both"
```

### CLI Integration

```bash
# Normal build (incremental if enabled)
[name] build

# Force full rebuild
[name] build --clean

# Show what would be rebuilt
[name] build --dry-run

# Clear cache
[name] cache clear
```

Clean builds MUST finish in a staging directory beside the configured output
directory and publish the completed generation before reporting success. A
failed clean build MUST leave the previously published output available.
Generated-file ownership recorded in the build cache MUST remain valid after
the staged directory is published at its configured output path. An
incremental build following a clean build MUST be able to update those files
without treating its own generated output as an unrelated user file.

---

## Production-Safe Incremental Mode

Production-safe incremental mode is separate from fast mode. It MAY reuse
unchanged posts, dependency results, and glob state, but it MUST continue to
run output-quality and search steps that can affect the published site.

Safe incremental mode MUST NOT imply any of the following fast-mode skips:

- JavaScript or CSS minification
- CSS purging or Tailwind generation
- Pagefind indexing
- redirect generation
- encryption or other global output validation

The implementation MUST fall back to a full rebuild when configuration,
templates, global styles, redirects, search configuration, or other global
inputs change. A changed content file MAY use the dependency graph to limit
post rendering, but feed, search, and asset outputs MUST be refreshed when
their inputs changed.

## Nested build spans

`pkg/buildstats.StartSpan(ctx, name, attributes...)` MUST return a context carrying
its span ID and a completion handle. Children MUST retain a parent ID only when
that context belongs to the same active profile. Without an active profile, the
API MUST return the original context and a safe no-op handle; nil contexts use
`context.Background()`.

Completed spans MUST appear in `Summary.Spans` and benchmark JSON at
`benchmark.Spans`, ordered by start offset, then ID. Each span MUST contain an ID,
name, relative start offset, non-negative duration, and `ok` or `error` status.
Stage and plugin attribution MUST capture the active values at start. `End` and
`EndError` MUST record at most one completion, including concurrent calls.
`EndError` MUST discard error text. Callers MUST finish spans before stopping the
profile; unfinished spans are omitted.

Names MUST default to `unnamed` when blank and be limited to 128 bytes.
Attributes MUST contain at most 16 entries with trimmed keys and values, limited
to 64 and 256 bytes respectively. Blank keys and keys containing authorization,
cookie, password, passwd, secret, token, or credential (case-insensitive) MUST be
dropped. Callers MUST supply safe values: this key filter does not redact URLs,
paths, or credentials embedded in other attributes. This foundation does not
require operation instrumentation or critical-path computation.

## See Also

- [SPEC.md](./SPEC.md) - Full specification
- [CONFIG.md](./CONFIG.md) - Configuration system
- [PLUGINS.md](./PLUGINS.md) - Plugin development guide
- [DATA_MODEL.md](./DATA_MODEL.md) - Post and config models

## Optional OTLP configuration foundation

`ResolveOTLPConfig(getenv)` MUST resolve a trace-specific endpoint/protocol before
the generic OTLP endpoint/protocol. Whitespace is trimmed; the default protocol
is `http/protobuf`. Supported protocols are `grpc`, `http/protobuf`, and
`http/json`. Endpoints MUST be HTTP(S) URLs with a host. A generic HTTP endpoint
MUST receive a `/v1/traces` path suffix; generic gRPC and trace-specific endpoints
MUST retain their paths. Query strings MUST survive HTTP path construction.
Invalid endpoints MUST return errors without echoing the supplied URL.

An absent endpoint MUST leave `Enabled()` false. A nil environment reader MUST
return an empty disabled config. Unsupported protocols MUST return an error,
even without an endpoint. Resolution MUST perform no network I/O. This API is
configuration groundwork, not an exporter; it does not enable trace export in
the build CLI.
