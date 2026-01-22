# Build Your Static Site Generator

> Fill in the two values below, then hand this spec to an AI agent.

## Your Choices

```
Name:     [name]
Language: [language]
```

**Name:** What to call your SSG (e.g., `blogkit`, `sitegen`, `pylon`, `grove`)

**Language:** `Python`, `TypeScript`, `Go`, or `Rust`

That's it. The agent will select appropriate libraries for your language.

---

## Prompt for AI Agent

Copy everything below and give it to your AI coding assistant:

---

**Build a static site generator called `[name]` following this specification.**

### Language & Stack

Use **[language]** with idiomatic libraries:

| Component | Python | TypeScript | Go | Rust |
|-----------|--------|------------|-----|------|
| Plugin system | pluggy | tapable | interfaces | traits |
| Markdown | markdown-it-py | markdown-it | goldmark | pulldown-cmark |
| Templates | Jinja2 | Nunjucks | html/template | tera |
| Config | tomllib + pydantic | @iarna/toml + zod | BurntSushi/toml | toml + serde |
| CLI | typer | commander | cobra | clap |
| Cache | diskcache | keyv | badger | sled |

Choose the standard/popular option for each component in the selected language.

### What to Build

Create a plugin-driven static site generator with:

1. **Core Orchestrator**
   - Manages 13 lifecycle stages (see LIFECYCLE.md)
   - Loads plugins from configuration
   - Provides `filter(expr)` and `map(field, filter?, sort?)` for querying posts
   - Caches expensive operations

2. **Feed System** (the differentiator - see FEEDS.md)
   - Define feeds with filter/sort/pagination
   - Output to multiple formats from one definition:
     - HTML (paginated index pages)
     - RSS 2.0
     - Atom
     - JSON Feed
     - Markdown
     - Plain text
   - Auto-generate tag/category/date archive feeds

3. **Standard Plugins** (see DEFAULT_PLUGINS.md)
   - `glob` - Find markdown files
   - `load` - Parse frontmatter + content into Post objects
   - `render_markdown` - Convert markdown to HTML (tables, code blocks, admonitions)
   - `jinja_md` - Process template expressions in markdown content
   - `templates` - Wrap content in HTML templates
   - `feeds` - Generate feed collections
   - `publish_feeds` - Write HTML/RSS/Atom/JSON/MD/TXT
   - `publish_html` - Write individual post HTML files

4. **Data Models**
   - `Post`: path, content, slug, href, title, date, published, tags, html
   - `Config`: output_dir, hooks, url, title, glob_patterns, feeds
   - `Feed`: slug, title, filter, sort, posts, pagination, formats

5. **CLI Commands**
   - `build` - Run full build
   - `serve` - Build + local dev server with file watching
   - `new <title>` - Create new post from template

### Key Behaviors

**Plugin hooks:**
```
@hook_impl
def render(core):
    for post in core.filter("not skip"):
        post.html = process(post)
```

**Querying posts:**
```
core.filter("published == True and date <= today")
core.filter("'python' in tags")
core.map("title", filter="published == True", sort="date", reverse=True)
```

**Jinja in markdown** (when `jinja: true` in frontmatter):
```markdown
---
title: All Posts
jinja: true
---
{% for post in core.filter("published == True") %}
- [{{ post.title }}]({{ post.href }})
{% endfor %}
```

**Configuration** (`[name].toml`):
```toml
[name]
output_dir = "public"
url = "https://example.com"

[name.glob]
glob_patterns = ["posts/**/*.md"]

# Feeds - the core feature
[[name.feeds]]
slug = "blog"
title = "Blog"
filter = "published == True"
sort = "date"
reverse = true
items_per_page = 10

[name.feeds.formats]
html = true
rss = true
atom = true
json = true

# Auto-generate tag pages
[name.feeds.auto_tags]
enabled = true
slug_prefix = "tags"
```

### Output Structure

```
public/
├── index.html               # From default feed
├── my-post/
│   └── index.html           # Individual post
├── blog/
│   ├── index.html           # Feed page 1
│   ├── page/2/index.html    # Feed page 2
│   ├── rss.xml              # RSS feed
│   ├── atom.xml             # Atom feed
│   └── feed.json            # JSON Feed
├── tags/
│   ├── python/
│   │   ├── index.html
│   │   └── rss.xml
│   └── rust/
│       ├── index.html
│       └── rss.xml
└── sitemap.xml
```

### Specification Files

Read these for complete details:
- `spec/SPEC.md` - Architecture overview
- `spec/LIFECYCLE.md` - Build stages
- `spec/FEEDS.md` - Feed system (multi-format output)
- `spec/DEFAULT_PLUGINS.md` - Built-in plugins
- `spec/PLUGINS.md` - Plugin development
- `spec/DATA_MODEL.md` - Post/Config schemas
- `spec/CONTENT.md` - Markdown processing
- `spec/TEMPLATES.md` - Template system
- `spec/tests.yaml` - Test cases to pass

### Deliverables

Generate:
1. Complete source code with all plugins
2. Default templates (base.html, post.html, feed.html)
3. Example config file
4. README with usage instructions
5. Tests that verify behavior matches `tests.yaml`

The implementation is complete when `tests.yaml` test cases pass.

---

## Example Usage (After Generation)

```bash
# Install
pip install ./[name]  # or npm install, go install, cargo install

# Create a post
[name] new "Hello World"

# Build
[name] build

# Serve locally
[name] serve
```

---

## See Also

- [SPEC.md](./SPEC.md) - Core architecture and CLI
- [CONFIG.md](./CONFIG.md) - Configuration system
- [THEMES.md](./THEMES.md) - Theming and customization
- [LIFECYCLE.md](./LIFECYCLE.md) - Build lifecycle stages
- [FEEDS.md](./FEEDS.md) - Feed system (core differentiator)
- [DEFAULT_PLUGINS.md](./DEFAULT_PLUGINS.md) - Built-in plugins
- [tests.yaml](./tests.yaml) - Test cases for verification
