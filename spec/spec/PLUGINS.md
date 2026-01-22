# Plugin Development Guide

This guide explains how to create plugins for the static site generator.

## Philosophy

Plugins are the primary extension mechanism. A single file should be able to:
- Hook into any lifecycle stage
- Extend configuration with new fields
- Extend the post model with new fields
- Register attributes on the core instance
- Add CLI commands

## Plugin Structure

A minimal plugin is a single file (or module) that exports hook implementations:

```
my_plugin.[ext]
├── Configuration model (optional)
├── Post model extension (optional)
├── Hook implementations
└── Helper functions (private)
```

## Hook Implementation

### Basic Hook

Hooks are functions that receive the core orchestrator instance:

```python
# Python with pluggy
from your_ssg.hookspec import hook_impl

@hook_impl
def render(core):
    """Process posts during the render stage."""
    for post in core.filter("not skip"):
        post.custom_field = process(post)
```

```typescript
// TypeScript with tapable
export function render(core: Core): void {
  for (const post of core.filter("not skip")) {
    post.customField = process(post);
  }
}
```

```go
// Go with interfaces
func (p *MyPlugin) Render(core *Core) error {
    for _, post := range core.Filter("not skip") {
        post.CustomField = process(post)
    }
    return nil
}
```

### Hook Ordering

Control when your hook runs relative to other plugins:

| Modifier | Effect |
|----------|--------|
| `tryfirst` | Run before other plugins in this stage |
| `trylast` | Run after other plugins in this stage |
| (none) | Run in registration order |

```python
@hook_impl(tryfirst=True)
def configure(core):
    """Run early in configure stage."""
    pass

@hook_impl(trylast=True)
def render(core):
    """Run late in render stage."""
    pass
```

### Attribute Registration

Declare attributes your plugin creates so the system can track dependencies:

```python
@hook_impl
@register_attr("reading_times")
def configure(core):
    """Initialize reading times storage."""
    core.reading_times = {}
```

This allows:
- Documentation of what attributes exist
- Dependency checking between plugins
- IDE autocompletion (with proper typing)

---

## Available Hooks

### Configuration Phase

| Hook | Stage | Purpose |
|------|-------|---------|
| `config_model` | config_model | Register configuration schema |
| `post_model` | post_model | Register post field schema |
| `configure` | configure | Initialize plugin state |
| `validate_config` | validate_config | Validate configuration |

### Content Phase

| Hook | Stage | Purpose |
|------|-------|---------|
| `glob` | glob | Discover content files |
| `load` | load | Parse and load content |
| `pre_render` | pre_render | Process before rendering |
| `render` | render | Convert content to HTML |
| `post_render` | post_render | Process after rendering |

### Output Phase

| Hook | Stage | Purpose |
|------|-------|---------|
| `save` | save | Write files to disk |
| `teardown` | teardown | Cleanup resources |

### CLI Phase

| Hook | When | Purpose |
|------|------|---------|
| `cli` | CLI init | Add CLI commands |

---

## Extending Configuration

Plugins can add their own configuration sections:

### Python (Pydantic)

```python
import pydantic

class WordCountConfig(pydantic.BaseModel):
    """Configuration for word count plugin."""
    enabled: bool = True
    words_per_minute: int = 200
    count_code_blocks: bool = False

class Config(pydantic.BaseModel):
    """Adds word_count section to config."""
    word_count: WordCountConfig = WordCountConfig()

@hook_impl
def config_model(core):
    core.config_models.append(Config)
```

### TypeScript (Zod)

```typescript
import { z } from 'zod';

const WordCountConfig = z.object({
  enabled: z.boolean().default(true),
  words_per_minute: z.number().default(200),
  count_code_blocks: z.boolean().default(false),
});

export function configModel(core: Core): void {
  core.configModels.push({ word_count: WordCountConfig });
}
```

### Go (Struct Tags)

```go
type WordCountConfig struct {
    Enabled        bool `toml:"enabled" default:"true"`
    WordsPerMinute int  `toml:"words_per_minute" default:"200"`
    CountCodeBlocks bool `toml:"count_code_blocks" default:"false"`
}

func (p *WordCountPlugin) ConfigModel(core *Core) {
    core.RegisterConfigSection("word_count", WordCountConfig{})
}
```

### Rust (Serde)

```rust
#[derive(Debug, Deserialize, Default)]
pub struct WordCountConfig {
    #[serde(default = "default_true")]
    pub enabled: bool,
    #[serde(default = "default_wpm")]
    pub words_per_minute: u32,
    #[serde(default)]
    pub count_code_blocks: bool,
}

fn default_true() -> bool { true }
fn default_wpm() -> u32 { 200 }

impl Plugin for WordCount {
    fn config_model(&self, core: &mut Core) {
        core.register_config::<WordCountConfig>("word_count");
    }
}
```

Users configure in their config file:

```toml
[your-ssg.word_count]
enabled = true
words_per_minute = 250
count_code_blocks = true
```

Access in your plugin:

```python
@hook_impl
def render(core):
    if not core.config.word_count.enabled:
        return
    
    wpm = core.config.word_count.words_per_minute
    # ...
```

---

## Extending the Post Model

Add fields to every post:

```python
import pydantic
from typing import Optional

class WordCountPostFields(pydantic.BaseModel):
    """Fields added to Post model."""
    word_count: int = 0
    reading_time_minutes: int = 0
    reading_time_display: Optional[str] = None

@hook_impl
def post_model(core):
    core.post_models.append(WordCountPostFields)
```

Now every post has these fields with defaults. Populate them:

```python
@hook_impl
def pre_render(core):
    wpm = core.config.word_count.words_per_minute
    
    for post in core.filter("not skip"):
        words = len(post.content.split())
        post.word_count = words
        post.reading_time_minutes = max(1, words // wpm)
        post.reading_time_display = f"{post.reading_time_minutes} min read"
```

---

## Accessing Core Features

### Filtering Posts

```python
# All posts
all_posts = core.posts

# Published posts
published = core.filter("published == True")

# Posts with specific tag
python_posts = core.filter("'python' in tags")

# Complex filter
recent_published = core.filter(
    "published == True and date <= today and date >= today - timedelta(days=30)"
)
```

### Mapping Posts

```python
# Get all titles
titles = core.map("title")

# Get titles of published posts, sorted by date
recent_titles = core.map(
    "title",
    filter="published == True",
    sort="date",
    reverse=True
)

# Get post objects (not just a field)
posts = core.map("post", filter="True")
```

### Using the Cache

```python
@hook_impl
def render(core):
    for post in core.filter("not skip"):
        cache_key = f"expensive_op:{post.path}:{hash(post.content)}"
        
        cached = core.cache.get(cache_key)
        if cached is not None:
            post.expensive_result = cached
            continue
        
        result = expensive_operation(post)
        core.cache.set(cache_key, result)
        post.expensive_result = result
```

### Accessing Configuration

```python
output_dir = core.config.output_dir
site_url = core.config.url
custom_setting = core.config.my_plugin.setting
```

---

## Adding CLI Commands

```python
@hook_impl
def cli(app):
    @app.command()
    def word_stats():
        """Show word count statistics for all posts."""
        from your_ssg import Core
        
        core = Core()
        core.run("pre_render")  # Run up to pre_render stage
        
        total_words = sum(p.word_count for p in core.posts)
        avg_words = total_words // len(core.posts) if core.posts else 0
        
        print(f"Total posts: {len(core.posts)}")
        print(f"Total words: {total_words:,}")
        print(f"Average words: {avg_words:,}")
```

---

## Complete Plugin Example

Here's a full plugin that adds reading time to posts:

```python
"""
Reading Time Plugin

Calculates and displays estimated reading time for posts.

Configuration:
    [your-ssg.reading_time]
    enabled = true
    words_per_minute = 200
    format = "{minutes} min read"

Usage in templates:
    {{ post.reading_time }}
"""
from typing import TYPE_CHECKING, Optional
import pydantic

from your_ssg.hookspec import hook_impl, register_attr

if TYPE_CHECKING:
    from your_ssg import Core


# =============================================================================
# Configuration
# =============================================================================

class ReadingTimeConfig(pydantic.BaseModel):
    """Plugin configuration."""
    enabled: bool = True
    words_per_minute: int = 200
    format: str = "{minutes} min read"


class Config(pydantic.BaseModel):
    """Config model fragment."""
    reading_time: ReadingTimeConfig = ReadingTimeConfig()


# =============================================================================
# Post Model
# =============================================================================

class ReadingTimePostFields(pydantic.BaseModel):
    """Post model fragment."""
    word_count: int = 0
    reading_time: Optional[str] = None


# =============================================================================
# Hooks
# =============================================================================

@hook_impl
def config_model(core: "Core") -> None:
    """Register configuration model."""
    core.config_models.append(Config)


@hook_impl
def post_model(core: "Core") -> None:
    """Register post model fields."""
    core.post_models.append(ReadingTimePostFields)


@hook_impl
def pre_render(core: "Core") -> None:
    """Calculate reading time for each post."""
    config = core.config.reading_time
    
    if not config.enabled:
        return
    
    for post in core.filter("not skip"):
        # Count words
        words = len(post.content.split())
        post.word_count = words
        
        # Calculate reading time
        minutes = max(1, words // config.words_per_minute)
        post.reading_time = config.format.format(minutes=minutes)


@hook_impl
def cli(app) -> None:
    """Add CLI command."""
    @app.command()
    def reading_stats():
        """Show reading time statistics."""
        from your_ssg import Core
        
        core = Core()
        core.run("pre_render")
        
        if not core.posts:
            print("No posts found")
            return
        
        times = [p.word_count // core.config.reading_time.words_per_minute 
                 for p in core.posts]
        
        print(f"Posts: {len(core.posts)}")
        print(f"Shortest: {min(times)} min")
        print(f"Longest: {max(times)} min")
        print(f"Average: {sum(times) // len(times)} min")
```

---

## Plugin Best Practices

### 1. Check Configuration

Always check if your plugin is enabled:

```python
@hook_impl
def render(core):
    if not core.config.my_plugin.enabled:
        return
    # ... rest of plugin
```

### 2. Use `not skip` Filter

Respect the skip flag that other plugins set:

```python
for post in core.filter("not skip"):  # Good
    process(post)

for post in core.posts:  # May process posts that should be skipped
    process(post)
```

### 3. Handle Missing Fields Gracefully

Posts may not have all optional fields:

```python
# Good
title = getattr(post, 'title', None) or post.slug

# Also good (with Pydantic defaults)
title = post.title or post.slug
```

### 4. Cache Expensive Operations

Use content hashes for cache keys:

```python
import hashlib

def get_cache_key(post):
    content_hash = hashlib.md5(post.content.encode()).hexdigest()[:8]
    return f"my_plugin:{post.path}:{content_hash}"
```

### 5. Log Appropriately

```python
import logging

logger = logging.getLogger(__name__)

@hook_impl
def render(core):
    logger.debug(f"Processing {len(core.posts)} posts")
    
    for post in core.filter("not skip"):
        logger.debug(f"Processing: {post.path}")
        # ...
    
    logger.info(f"Processed {len(core.posts)} posts")
```

### 6. Document Your Plugin

Include a docstring with:
- What the plugin does
- Configuration options
- Template usage examples

---

## Testing Plugins

### Unit Test Example

```python
import pytest
from your_ssg import Core

def test_reading_time_calculation():
    """Test that reading time is calculated correctly."""
    core = Core()
    core.config.reading_time.words_per_minute = 200
    
    # Create a mock post with 400 words
    post = core.Post(
        path="test.md",
        content=" ".join(["word"] * 400),
        slug="test"
    )
    core.articles = [post]
    
    # Run the pre_render stage (which calculates reading time)
    core.run("pre_render")
    
    assert post.word_count == 400
    assert post.reading_time == "2 min read"


def test_reading_time_disabled():
    """Test that plugin respects enabled flag."""
    core = Core()
    core.config.reading_time.enabled = False
    
    post = core.Post(path="test.md", content="words", slug="test")
    core.articles = [post]
    
    core.run("pre_render")
    
    assert post.reading_time is None
```

### Integration Test

```python
def test_full_build_with_reading_time(tmp_path):
    """Test reading time in full build."""
    # Create test content
    (tmp_path / "posts").mkdir()
    (tmp_path / "posts" / "test.md").write_text("""---
title: Test Post
published: true
---

This is test content with enough words to calculate reading time.
""" + " word" * 200)
    
    # Create config
    (tmp_path / "config.toml").write_text("""
[your-ssg]
output_dir = "output"

[your-ssg.glob]
glob_patterns = ["posts/*.md"]
""")
    
    # Run build
    from your_ssg import Core
    core = Core(root=tmp_path)
    core.run()
    
    # Verify output
    output = (tmp_path / "output" / "test-post" / "index.html").read_text()
    assert "1 min read" in output
```

---

## Registering Your Plugin

### Option 1: Local Plugin

Put your plugin file in your project:

```
my-site/
├── plugins/
│   └── reading_time.py
├── posts/
└── config.toml
```

Configure:

```toml
[your-ssg]
hooks = ["default", "plugins.reading_time"]
```

### Option 2: Installed Package

Create a package and install it:

```
pip install your-ssg-reading-time
```

Configure:

```toml
[your-ssg]
hooks = ["default", "your_ssg_reading_time"]
```

### Option 3: Entry Points (Python)

In your package's `pyproject.toml`:

```toml
[project.entry-points."your_ssg"]
reading_time = "your_ssg_reading_time:plugin"
```

The plugin is auto-discovered if the system supports entry points.

---

## See Also

- [SPEC.md](./SPEC.md) - Full specification
- [CONFIG.md](./CONFIG.md) - Plugin configuration declaration
- [LIFECYCLE.md](./LIFECYCLE.md) - Lifecycle stages detail
- [DATA_MODEL.md](./DATA_MODEL.md) - Post and config models
