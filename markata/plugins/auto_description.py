"""
A Markata plugin to create automatic descriptions for markdown documents.  It
does this by grabbing the first `{len}` number of characters from the document
that are in a paragraph.

## Configuration

Open up your `markata.toml` file and add new entries for your
auto_descriptions.  You can have multiple desriptions, each one will be named
after the key you give it in your config.

``` toml
[markata]

# make sure its in your list of hooks
hooks=[
   "markata.plugins.auto_description",
   ]

[markata.auto_description.description]
len=160
[markata.auto_description.long_description]
len=250
[markata.auto_description.super_description]
len=500
```

!!! note
   Make sure that you have the auto_description plugin in your configured hooks.

In the above we will end up with three different descritpions,
(`description`, `long_description`, and `super_description`) each will be the
first number of characters from the document as specified in the config.

### Defaults

By default markata will set `description` to 160 and `long_description` to 250,
if they are not set in your config.

### Using the Description

Downstream hooks can now use the description for things such as seo, or feeds.
Here is a simple example that lists all of the descriptions in all posts.  This
is a handy thing you can do right from a repl.

``` python
from markata import Markata
m = Markata()
[p["description"] for p in m.articles]
```

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
