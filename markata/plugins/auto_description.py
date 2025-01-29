"""
The `markata.plugins.auto_description` plugin automatically generates descriptions for your
posts by extracting text from the first paragraphs of your markdown content. It can create
multiple descriptions of different lengths for different purposes (e.g., SEO, previews).

# Installation

This plugin is built-in and enabled by default through the 'default' plugin.
If you want to be explicit, you can add it to your list of plugins:

```toml
hooks = [
    "markata.plugins.auto_description",
]
```

# Uninstallation

Since this plugin is included in the default plugin set, to disable it you must explicitly
add it to the disabled_hooks list if you are using the 'default' plugin:

```toml
disabled_hooks = [
    "markata.plugins.auto_description",
]
```

# Configuration

Configure multiple description types in your `markata.toml`:

```toml
# Standard description (160 characters)
[markata.auto_description.description]
len = 160

# Longer description (250 characters)
[markata.auto_description.long_description]
len = 250

# Full description (500 characters)
[markata.auto_description.super_description]
len = 500
```

Each description configuration:
- Must be under `markata.auto_description.[name]`
- Requires a `len` parameter specifying maximum character length
- Creates an attribute named after the configuration key

## Frontmatter Override

You can override automatic descriptions by setting them in frontmatter:

```markdown
---
title: My Post
description: My custom description
long_description: A longer custom description
---
```

# Functionality

## Description Generation

The plugin:
1. Converts markdown to plain text
2. Finds the first meaningful paragraphs
3. Truncates to the specified length
4. Ensures clean word breaks
5. Caches results for performance

## Registered Attributes

For each configured description (e.g., `description`, `long_description`), the plugin adds:
- The description attribute with the generated/specified text
- Truncated to the configured length
- Preserving complete words

## Dependencies

This plugin depends on:
- markdown-it-py for markdown parsing

"""

import html
from pathlib import Path
from typing import Any, Dict, TYPE_CHECKING


from markata.hookspec import hook_impl

if TYPE_CHECKING:
    from diskcache import FanoutCache
    from frontmatter import Post

    from markata import Markata


def get_description(article: "Post") -> str:
    """
    Get the full-length description for a single post by converting markdown to plain text.
    Uses markdown-it-py to parse the markdown and extracts text content from all nodes.
    """
    from markdown_it import MarkdownIt

    md = MarkdownIt("commonmark")
    tokens = md.parse(article.content)

    # Extract text content from all inline tokens
    text_chunks = []
    for token in tokens:
        if token.type == "inline" and token.content:
            text_chunks.append(token.content)

    description = " ".join(text_chunks)
    return html.escape(description)


def set_description(
    markata: "Markata",
    article: "Post",
    cache: "FanoutCache",
    config: Dict,
    max_description: int = 500,
    plugin_text: None = "",
) -> None:
    """
    For a given `article`, find the description, put it in the cache, and set
    the configured descriptions for the article.
    """
    key = markata.make_hash(
        "auto_description",
        article.content,
        plugin_text,
        config,
    )

    description_from_cache = markata.precache.get(key)
    if description_from_cache is None:
        description = get_description(article)[:max_description]
        markata.cache.set(key, description, expire=markata.config.default_cache_expire)
    else:
        description = description_from_cache

    for description_key in config:
        if description_key not in ["cache_expire", "config_key"]:
            # overwrites missing (None) and empty ('')
            if not article.metadata.get(description_key):
                article.metadata[description_key] = description[
                    : config[description_key]["len"]
                ]


@hook_impl
def pre_render(markata: "Markata") -> None:
    """
    The Markata hook that will set descriptions for all posts in the pre-render phase.
    """
    config = markata.get_plugin_config(__file__)

    if "description" not in config.keys():
        config["description"] = {}
        config["description"]["len"] = 160

    if "long_description" not in config.keys():
        config["long_description"] = {}
        config["long_description"]["len"] = 250

    def try_config_get(key: str) -> Any:
        try:
            return config.get(key).get("len") or None
        except AttributeError:
            return None

    max_description = max(
        [
            value
            for description_key in config
            if (value := try_config_get(description_key))
        ],
    )

    with markata.cache as cache:
        for article in markata.iter_articles("setting auto description"):
            set_description(
                markata=markata,
                article=article,
                cache=cache,
                config=config,
                max_description=max_description,
                plugin_text=Path(__file__).read_text(),
            )
