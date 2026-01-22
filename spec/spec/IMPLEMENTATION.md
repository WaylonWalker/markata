# Implementation Guide

This guide provides step-by-step instructions for implementing the static site generator specification in any language.

## Overview

```
Implementation Phases
=====================

Phase 1: Foundation (Core)           Phase 2: Content Pipeline
─────────────────────────────        ─────────────────────────
1. Project structure                 5. Glob plugin (file discovery)
2. Configuration system              6. Load plugin (frontmatter + content)
3. Lifecycle manager                 7. Markdown rendering
4. Plugin manager                    8. Template processing

Phase 3: Output Generation           Phase 4: Developer Experience
─────────────────────────────        ─────────────────────────────
9. HTML publishing                   13. Dev server with live reload
10. Feed system                      14. CLI commands (build, serve, new)
11. Asset copying                    15. Caching layer
12. RSS/Atom/Sitemap                 16. Incremental builds
```

---

## Phase 1: Foundation

### Step 1.1: Project Structure

```
your-ssg/
├── src/
│   ├── core/
│   │   ├── config.{ext}        # Configuration loading & merging
│   │   ├── lifecycle.{ext}     # Stage orchestration
│   │   ├── plugin.{ext}        # Plugin loading & registration
│   │   ├── cache.{ext}         # Persistent cache
│   │   ├── query.{ext}         # Filter & map expressions
│   │   └── models.{ext}        # Post & Config models
│   ├── plugins/
│   │   ├── glob.{ext}
│   │   ├── load.{ext}
│   │   ├── render_markdown.{ext}
│   │   ├── templates.{ext}
│   │   ├── feeds.{ext}
│   │   ├── publish_html.{ext}
│   │   └── copy_assets.{ext}
│   └── cli/
│       ├── build.{ext}
│       ├── serve.{ext}
│       └── new.{ext}
├── tests/
│   └── spec/                   # Tests from tests.yaml
└── templates/                  # Default templates
```

### Step 1.2: Choose Your Lifecycle Variant

| Variant | Stages | Recommended For |
|---------|--------|-----------------|
| **Full (13)** | config_model → post_model → create_models → load_config → configure → validate_config → glob → load → pre_render → render → post_render → save → teardown | Python, Ruby |
| **Standard (9)** | configure → validate → glob → load → transform → render → collect → write → cleanup | TypeScript, Go |
| **Minimal (6)** | init → load → process → build → output | Rust, embedded |

**Recommendation:** Start with **Standard (9 stages)** unless you have specific needs for runtime schema composition (Full) or maximum simplicity (Minimal).

### Step 1.3: Core Data Structures

**Post Model:**

```
Post {
  // Required
  path: Path              // Source file path
  content: string         // Raw content (after frontmatter)
  slug: string            // URL-safe identifier
  href: string            // Relative URL (/{slug}/)
  
  // Standard Optional
  title: string?
  date: date?
  published: bool = false
  draft: bool = false
  skip: bool = false
  tags: string[] = []
  description: string?
  template: string = "post.html"
  
  // Set during rendering
  article_html: string?   // Rendered markdown
  html: string?           // Final HTML with template
}
```

**Config Model:**

```
Config {
  output_dir: Path = "output"
  url: URL?
  title: string?
  description: string?
  hooks: string[] = ["default"]
  disabled_hooks: string[] = []
  concurrency: int = 0    // 0 = auto
  
  // Plugin sections (added dynamically)
  glob: GlobConfig
  markdown: MarkdownConfig
  feeds: FeedConfig[]
  // ...
}
```

**Core Interface:**

```
Core {
  // Data
  config: Config
  posts: Post[]
  files: Path[]
  feeds: Feed[]
  cache: Cache
  
  // Query methods
  filter(expr: string) -> Post[]
  map(field: string, filter?: string, sort?: string) -> any[]
  first(filter?: string, sort?: string) -> Post?
  one(expr: string) -> Post  // Error if 0 or >1
  
  // Plugin communication
  set(key: string, value: any)
  get(key: string) -> any
  has(key: string) -> bool
  
  // Lifecycle
  run(until_stage?: string)
  register(plugin: Plugin)
}
```

### Step 1.4: Configuration Loading

**Resolution order (highest precedence first):**
1. CLI arguments
2. Environment variables (`[NAME]_[SECTION]_[KEY]`)
3. Local config (`./[name].toml`)
4. Global config (`~/.config/[name]/config.toml`)
5. Plugin defaults

**Implementation:**

```python
# Pseudocode
def load_config():
    config = {}
    
    # 1. Plugin defaults
    for plugin in registered_plugins:
        merge(config, plugin.defaults)
    
    # 2. Global config
    if exists("~/.config/[name]/config.toml"):
        merge(config, parse_toml(read("~/.config/[name]/config.toml")))
    
    # 3. Local config
    for path in ["[name].toml", "pyproject.toml"]:
        if exists(path):
            merge(config, parse_toml(read(path)))
            break
    
    # 4. Environment variables
    for key, value in environ:
        if key.startswith("[NAME]_"):
            set_nested(config, parse_env_key(key), value)
    
    # 5. CLI arguments
    merge(config, cli_args)
    
    return validate(config)
```

### Step 1.5: Plugin Manager

**Plugin Interface:**

```
Plugin {
  name: string
  
  // Hooks (all optional)
  configure?(core: Core)
  validate?(core: Core)
  glob?(core: Core)
  load?(core: Core)
  transform?(core: Core)      // or pre_render for Full lifecycle
  render?(core: Core)
  collect?(core: Core)        // or post_render for Full lifecycle
  write?(core: Core)          // or save for Full lifecycle
  cleanup?(core: Core)        // or teardown for Full lifecycle
}
```

**Plugin Resolution:**

```python
def resolve_plugin(name: str) -> Plugin:
    if name == "default":
        return DEFAULT_PLUGINS
    
    if name in BUILTIN_PLUGINS:
        return BUILTIN_PLUGINS[name]
    
    if name.startswith("./"):
        return load_local_plugin(name)
    
    # Module path (language-specific)
    return import_plugin(name)
```

### Step 1.6: Lifecycle Manager

```python
class LifecycleManager:
    stages = [
        "configure", "validate", "glob", "load",
        "transform", "render", "collect", "write", "cleanup"
    ]
    current_stage = None
    completed_stages = set()
    
    def run(self, until_stage=None):
        for stage in self.stages:
            if stage in self.completed_stages:
                continue
            
            self.current_stage = stage
            self.run_stage(stage)
            self.completed_stages.add(stage)
            
            if stage == until_stage:
                break
    
    def run_stage(self, stage):
        # Get all hooks for this stage
        hooks = []
        for plugin in self.plugins:
            if hasattr(plugin, stage):
                hook = getattr(plugin, stage)
                hooks.append((plugin, hook))
        
        # Sort by priority (tryfirst, normal, trylast)
        hooks.sort(key=lambda h: get_hook_priority(h))
        
        # Execute hooks
        for plugin, hook in hooks:
            try:
                hook(self.core)
            except Exception as e:
                handle_error(e, stage, plugin)
```

---

## Phase 2: Content Pipeline

### Step 2.1: Glob Plugin

**Purpose:** Discover content files.

**Implementation:**

```python
def glob(core):
    patterns = core.config.glob.glob_patterns  # default: ["**/*.md"]
    exclude = core.config.glob.exclude_patterns
    use_gitignore = core.config.glob.use_gitignore
    
    files = []
    for pattern in patterns:
        for path in glob_files(pattern):
            if not matches_any(path, exclude):
                if not use_gitignore or not is_gitignored(path):
                    files.append(path)
    
    core.files = deduplicate(files)
```

**Test cases from tests.yaml:**
- `glob_finds_markdown_files`
- `glob_respects_exclude_patterns`
- `glob_respects_gitignore`

### Step 2.2: Load Plugin

**Purpose:** Parse frontmatter and create Post objects.

**Frontmatter Parsing:**

```python
def parse_frontmatter(content: str) -> (dict, str):
    if not content.startswith("---"):
        return ({}, content)
    
    # Find closing delimiter
    end = content.find("---", 3)
    if end == -1:
        return ({}, content)
    
    yaml_str = content[3:end]
    body = content[end+3:].lstrip("\n")
    
    frontmatter = yaml.safe_load(yaml_str) or {}
    return (frontmatter, body)
```

**Load Implementation:**

```python
def load(core):
    for path in core.files:
        try:
            raw = read_file(path, encoding="utf-8")
            frontmatter, content = parse_frontmatter(raw)
            
            # Generate slug
            slug = frontmatter.get("slug") or \
                   slugify(frontmatter.get("title")) or \
                   path.stem
            
            post = Post(
                path=path,
                content=content,
                slug=slug,
                href=f"/{slug}/",
                **frontmatter
            )
            core.posts.append(post)
        
        except Exception as e:
            log.warning(f"Failed to load {path}: {e}")
```

**Slug Generation:**

```python
def slugify(text: str) -> str:
    text = text.lower()
    text = re.sub(r"[^\w\s-]", "", text)    # Remove special chars
    text = re.sub(r"[-\s]+", "-", text)      # Replace spaces/hyphens
    text = text.strip("-")                    # Strip leading/trailing
    return text
```

### Step 2.3: Markdown Rendering

**Configuration:**

```toml
[name.markdown]
extensions = ["tables", "admonitions", "footnotes"]

[name.markdown.highlight]
enabled = true
theme = "github-dark"
```

**Implementation:**

```python
def render(core):
    md = create_markdown_parser(core.config.markdown)
    
    for post in core.filter("not skip"):
        post.article_html = md.render(post.content)
```

**Required Markdown Features:**
- CommonMark base
- GFM tables
- Fenced code blocks
- Syntax highlighting
- Admonitions (optional)

### Step 2.4: Template Processing

**Template Context:**

```python
context = {
    "post": post,
    "body": post.article_html,
    "config": core.config,
    "core": core,
    "today": date.today(),
    "now": datetime.now(),
}
```

**Template Rendering:**

```python
def render_post(post, core):
    template_name = post.template or "post.html"
    template = core.template_env.get_template(template_name)
    
    post.html = template.render(
        post=post,
        body=post.article_html,
        config=core.config,
        core=core,
    )
```

**Default Template (post.html):**

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{ post.title or post.slug }} | {{ config.title }}</title>
    {% if post.description %}
    <meta name="description" content="{{ post.description }}">
    {% endif %}
</head>
<body>
    <article>
        <header>
            <h1>{{ post.title }}</h1>
            {% if post.date %}
            <time datetime="{{ post.date.isoformat() }}">
                {{ post.date.strftime("%B %d, %Y") }}
            </time>
            {% endif %}
        </header>
        
        <div class="content">
            {{ body | safe }}
        </div>
    </article>
</body>
</html>
```

---

## Phase 3: Output Generation

### Step 3.1: HTML Publishing

**Implementation:**

```python
def write(core):
    output_dir = core.config.output_dir
    
    for post in core.filter("not skip"):
        # Render with template if not already done
        if not post.html:
            post.html = render_post(post, core)
        
        # Write to output
        output_path = output_dir / post.slug / "index.html"
        output_path.parent.mkdir(parents=True, exist_ok=True)
        output_path.write_text(post.html)
```

### Step 3.2: Feed System

**Feed Configuration:**

```toml
[[name.feeds]]
slug = "blog"
title = "Blog"
filter = "published == True"
sort = "date"
reverse = true
items_per_page = 10
template = "feed.html"

[name.feeds.blog.formats]
html = true
rss = true
atom = false
json = false
sitemap = false
```

**Feed Data Structure:**

```
Feed {
  slug: string
  title: string
  filter: string
  sort: string
  reverse: bool
  items_per_page: int
  template: string
  formats: {html, rss, atom, json, sitemap}
  
  // Computed
  posts: Post[]           // Filtered + sorted posts
  pages: Post[][]         // Paginated posts
  page_count: int
}
```

**Feed Generation:**

```python
def collect(core):
    core.feeds = []
    
    for feed_config in core.config.feeds:
        # Filter posts
        posts = core.filter(feed_config.filter)
        
        # Sort
        posts = sorted(
            posts,
            key=lambda p: getattr(p, feed_config.sort, ""),
            reverse=feed_config.reverse
        )
        
        # Paginate
        per_page = feed_config.items_per_page
        pages = [posts[i:i+per_page] for i in range(0, len(posts), per_page)]
        
        feed = Feed(
            slug=feed_config.slug,
            title=feed_config.title,
            posts=posts,
            pages=pages,
            page_count=len(pages),
            formats=feed_config.formats,
            template=feed_config.template,
        )
        core.feeds.append(feed)
```

**Writing Feed Outputs:**

```python
def write_feeds(core):
    for feed in core.feeds:
        # HTML pages
        if feed.formats.html:
            for page_num, page_posts in enumerate(feed.pages, 1):
                html = render_feed_page(feed, page_posts, page_num, core)
                
                if page_num == 1:
                    path = f"{feed.slug}/index.html"
                else:
                    path = f"{feed.slug}/page/{page_num}/index.html"
                
                write_output(core.config.output_dir / path, html)
        
        # RSS
        if feed.formats.rss:
            rss = render_rss(feed, core)
            write_output(core.config.output_dir / f"{feed.slug}/rss.xml", rss)
        
        # Atom
        if feed.formats.atom:
            atom = render_atom(feed, core)
            write_output(core.config.output_dir / f"{feed.slug}/atom.xml", atom)
```

### Step 3.3: Asset Copying

```python
def copy_assets(core):
    assets_dir = Path(core.config.assets.dir)
    output_dir = core.config.output_dir
    
    if not assets_dir.exists():
        return
    
    for src in assets_dir.rglob("*"):
        if src.is_file():
            rel_path = src.relative_to(assets_dir)
            dst = output_dir / rel_path
            dst.parent.mkdir(parents=True, exist_ok=True)
            shutil.copy2(src, dst)
```

---

## Phase 4: Developer Experience

### Step 4.1: Dev Server

**Features:**
- HTTP server on localhost
- File watching
- Automatic rebuilds
- Live reload via WebSocket/SSE

**Implementation Sketch:**

```python
class DevServer:
    def __init__(self, core, port=3000):
        self.core = core
        self.port = port
        self.watcher = FileWatcher()
    
    def start(self):
        # Initial build
        self.core.run()
        
        # Start HTTP server
        self.server = HTTPServer(self.core.config.output_dir, self.port)
        self.server.start()
        
        # Start file watcher
        self.watcher.watch(
            patterns=[
                "**/*.md",
                "templates/**/*.html",
                "[name].toml",
            ],
            callback=self.on_change
        )
    
    def on_change(self, path):
        # Determine what changed
        if path.endswith(".toml"):
            # Config change: full rebuild
            self.core = Core()  # Reload
            self.core.run()
        elif path.endswith(".md"):
            # Content change: incremental rebuild
            self.core.run()  # Smart rebuild
        elif "templates" in path:
            # Template change: re-render affected posts
            self.core.run(from_stage="render")
        
        # Trigger live reload
        self.notify_clients()
```

### Step 4.2: CLI Commands

**Build:**

```bash
[name] build [flags]
  -c, --config string   Config file path
  -o, --output string   Output directory
  -v, --verbose         Verbose output
  --clean               Remove output before build
  --drafts              Include drafts
```

**Serve:**

```bash
[name] serve [flags]
  -p, --port int        Port (default 3000)
  -H, --host string     Host (default localhost)
  --no-reload           Disable live reload
  --open                Open browser
```

**New:**

```bash
[name] new <path> [flags]
  -t, --title string    Post title
  --edit                Open in $EDITOR
```

### Step 4.3: Caching

**Cache Interface:**

```
Cache {
  get(key: string) -> any?
  set(key: string, value: any, ttl?: int)
  delete(key: string)
  clear()
}
```

**Cache Keys:**

```
post:{path}:content_hash    → cached article_html
template:{path}:hash        → compiled template
feed:{slug}:{filter_hash}   → cached feed data
```

**File-based Cache:**

```python
class FileCache:
    def __init__(self, cache_dir=".name.cache"):
        self.cache_dir = Path(cache_dir)
        self.cache_dir.mkdir(exist_ok=True)
    
    def get(self, key):
        path = self.cache_dir / hash_key(key)
        if path.exists():
            data = json.loads(path.read_text())
            if not data.get("expires") or data["expires"] > time.time():
                return data["value"]
        return None
    
    def set(self, key, value, ttl=None):
        path = self.cache_dir / hash_key(key)
        data = {
            "value": value,
            "expires": time.time() + ttl if ttl else None
        }
        path.write_text(json.dumps(data))
```

---

## Implementation Checklist

### Core (Required)

- [ ] Configuration loading from TOML
- [ ] Configuration merging (precedence order)
- [ ] Environment variable overrides
- [ ] Plugin loading and registration
- [ ] Lifecycle stage execution
- [ ] Post model with required fields
- [ ] Filter expression evaluation
- [ ] Map function for field extraction

### Content Pipeline (Required)

- [ ] Glob plugin - file discovery
- [ ] Load plugin - frontmatter parsing
- [ ] Slug generation (title → slug)
- [ ] Markdown rendering (CommonMark + extensions)
- [ ] Template processing (Jinja2 or equivalent)
- [ ] Template inheritance

### Output (Required)

- [ ] HTML publishing (post files)
- [ ] Feed generation with filtering/sorting
- [ ] Feed pagination
- [ ] Asset copying

### CLI (Required)

- [ ] `build` command
- [ ] `serve` command with HTTP server
- [ ] File watching for dev server

### Nice to Have

- [ ] Live reload (WebSocket/SSE)
- [ ] RSS/Atom feed formats
- [ ] Sitemap generation
- [ ] JSON feed format
- [ ] Incremental builds
- [ ] Asset fingerprinting
- [ ] `new` command for post creation
- [ ] `validate` command

---

## Testing Against Spec

Run your implementation against `tests.yaml`:

```bash
# Run all tests
your-ssg test --spec tests.yaml

# Run specific test category
your-ssg test --spec tests.yaml --filter "glob_*"
```

**Test Categories:**

| Category | Tests | Description |
|----------|-------|-------------|
| `config_*` | 15 | Configuration loading |
| `glob_*` | 8 | File discovery |
| `load_*` | 12 | Content loading |
| `render_*` | 20 | Markdown rendering |
| `filter_*` | 18 | Query expressions |
| `feed_*` | 25 | Feed generation |
| `template_*` | 15 | Template processing |
| `build_*` | 10 | Full build |

---

## Language-Specific Guidance

### Python

**Recommended Libraries:**
- **Config:** `tomli` (TOML), `pydantic` (validation)
- **Templates:** `jinja2`
- **Markdown:** `markdown-it-py` or `mistune`
- **CLI:** `typer` or `click`
- **Server:** `uvicorn` with `starlette`

**Plugin System:** Use `pluggy` for hook-based plugins.

### TypeScript

**Recommended Libraries:**
- **Config:** `@iarna/toml`, `zod` (validation)
- **Templates:** `nunjucks` or `liquid`
- **Markdown:** `markdown-it`
- **CLI:** `commander` or `yargs`
- **Server:** `express` with `chokidar` for watching

**Plugin System:** Use tapable or custom hook system.

### Go

**Recommended Libraries:**
- **Config:** `BurntSushi/toml`, `go-playground/validator`
- **Templates:** `html/template`
- **Markdown:** `goldmark`
- **CLI:** `cobra`
- **Server:** `net/http` with `fsnotify`

**Plugin System:** Interface-based plugins or yaegi for scripting.

### Rust

**Recommended Libraries:**
- **Config:** `toml`, `serde`
- **Templates:** `tera` or `minijinja`
- **Markdown:** `pulldown-cmark`
- **CLI:** `clap`
- **Server:** `axum` or `actix-web`

**Plugin System:** Trait-based plugins, consider WASM for external plugins.

---

## Common Pitfalls

### 1. Slug Collisions

Multiple posts can generate the same slug. Handle by:
- Appending numbers: `post`, `post-1`, `post-2`
- Warning and skipping
- Error and abort

### 2. Circular Template Inheritance

Detect cycles during template loading:

```python
def get_template(name, seen=None):
    if seen is None:
        seen = set()
    
    if name in seen:
        raise CircularTemplateError(f"Cycle detected: {' -> '.join(seen)} -> {name}")
    
    seen.add(name)
    template = load_template(name)
    
    if template.extends:
        get_template(template.extends, seen)
    
    return template
```

### 3. Filter Expression Safety

Don't use `eval()` directly. Parse and validate expressions:

```python
# UNSAFE
result = eval(f"post.{filter_expr}")

# SAFE - use a parser
ast = parse_filter_expr(filter_expr)
result = evaluate_ast(ast, context={"post": post, "today": today})
```

### 4. Unicode Handling

- Always read/write files as UTF-8
- Slugify unicode correctly: `"Über"` → `"uber"`
- Handle BOM in frontmatter

### 5. Large Site Performance

- Process posts concurrently
- Cache rendered markdown
- Use incremental builds
- Stream large files instead of loading into memory

---

## Example: Minimal Working Implementation

Here's a ~200 line Python implementation covering core functionality:

```python
# mini_ssg.py - Minimal SSG implementation
import tomli
import yaml
from pathlib import Path
from dataclasses import dataclass, field
from typing import Optional
import re
from jinja2 import Environment, FileSystemLoader
import markdown

@dataclass
class Post:
    path: Path
    content: str
    slug: str
    title: Optional[str] = None
    date: Optional[str] = None
    published: bool = False
    tags: list = field(default_factory=list)
    template: str = "post.html"
    article_html: str = ""
    html: str = ""
    
    @property
    def href(self):
        return f"/{self.slug}/"

@dataclass
class Config:
    output_dir: str = "output"
    glob_patterns: list = field(default_factory=lambda: ["**/*.md"])
    url: str = ""
    title: str = ""

class Core:
    def __init__(self, config_path="ssg.toml"):
        self.config = self._load_config(config_path)
        self.posts = []
        self.files = []
        self.jinja = Environment(loader=FileSystemLoader("templates"))
        self.md = markdown.Markdown(extensions=["tables", "fenced_code"])
    
    def _load_config(self, path):
        if Path(path).exists():
            data = tomli.loads(Path(path).read_text())
            return Config(**data.get("ssg", {}))
        return Config()
    
    def filter(self, expr):
        # Simple filter implementation
        result = []
        for post in self.posts:
            ctx = {"post": post, "True": True, "False": False}
            if eval(expr, {"__builtins__": {}}, ctx):
                result.append(post)
        return result
    
    def run(self):
        self._glob()
        self._load()
        self._render()
        self._write()
    
    def _glob(self):
        for pattern in self.config.glob_patterns:
            self.files.extend(Path(".").glob(pattern))
    
    def _load(self):
        for path in self.files:
            raw = path.read_text()
            fm, content = self._parse_frontmatter(raw)
            slug = fm.get("slug") or self._slugify(fm.get("title", path.stem))
            self.posts.append(Post(path=path, content=content, slug=slug, **fm))
    
    def _parse_frontmatter(self, content):
        if not content.startswith("---"):
            return {}, content
        end = content.find("---", 3)
        if end == -1:
            return {}, content
        fm = yaml.safe_load(content[3:end]) or {}
        body = content[end+3:].lstrip("\n")
        return fm, body
    
    def _slugify(self, text):
        text = text.lower()
        text = re.sub(r"[^\w\s-]", "", text)
        text = re.sub(r"[-\s]+", "-", text)
        return text.strip("-")
    
    def _render(self):
        for post in self.posts:
            post.article_html = self.md.convert(post.content)
            self.md.reset()
            
            template = self.jinja.get_template(post.template)
            post.html = template.render(post=post, config=self.config)
    
    def _write(self):
        output = Path(self.config.output_dir)
        for post in self.posts:
            path = output / post.slug / "index.html"
            path.parent.mkdir(parents=True, exist_ok=True)
            path.write_text(post.html)

if __name__ == "__main__":
    Core().run()
```

This minimal implementation covers:
- Configuration loading
- Frontmatter parsing
- Markdown rendering
- Template processing
- HTML output

From here, add features incrementally:
1. CLI with argparse/click
2. Feed generation
3. Dev server
4. Caching

---

## See Also

- [SPEC.md](./SPEC.md) - Full specification
- [LIFECYCLE.md](./LIFECYCLE.md) - Lifecycle stage details
- [PLUGINS.md](./PLUGINS.md) - Plugin development
- [FEEDS.md](./FEEDS.md) - Feed system
- [tests.yaml](./tests.yaml) - Test cases
