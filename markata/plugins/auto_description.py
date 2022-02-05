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
from pathlib import Path
from typing import TYPE_CHECKING, Dict

import commonmark

from markata.hookspec import hook_impl

if TYPE_CHECKING:
    from diskcache import FanoutCache
    from frontmatter import Post

    from markata import Markata

_parser = commonmark.Parser()


def get_description(article: "Post") -> str:
    """
    Get the full-length description for a single post using the commonmark
    parser.  Only paragraph nodes will count as text towards the description.
    """
    ast = _parser.parse(article.content)
    return " ".join(
        [
            node[0].first_child.literal
            for node in ast.walker()
            if node[0].t == "paragraph" and node[0].first_child.literal is not None
        ]
    )


def set_description(
    markata: "Markata", article: "Post", cache: "FanoutCache", config: Dict
) -> None:
    """
    For a given `article`, find the description, put it in the cache, and set
    the configured descriptions for the article.
    """
    key = markata.make_hash(
        "long_description",
        "render",
        Path(__file__).read_text(),
        article.content,
    )

    description_from_cache = cache.get(key)
    if description_from_cache is None:
        description = get_description(article)
        markata.cache.add(key, description, expire=config["cache_expire"])
    else:
        description = description_from_cache

    for description_key in config:
        if description_key not in ["cache_expire", "config_key"]:
            if description_key not in article.metadata.keys():
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
    with markata.cache as cache:
        for article in markata.iter_articles("setting auto description"):
            set_description(
                markata,
                article,
                cache,
                config,
            )
