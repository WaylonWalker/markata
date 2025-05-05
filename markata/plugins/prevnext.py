"""
The `markata.plugins.prevnext` plugin adds previous and next navigation links to each
post, allowing readers to easily navigate between related content.

## Installation

This plugin is built-in and enabled by default through the 'default' plugin.
If you want to be explicit, you can add it to your list of plugins:

```toml
hooks = [
    "markata.plugins.prevnext",
]
```

## Uninstallation

Since this plugin is included in the default plugin set, to disable it you must explicitly
add it to the disabled_hooks list if you are using the 'default' plugin:

```toml
disabled_hooks = [
    "markata.plugins.prevnext",
]
```

## Configuration

Configure navigation behavior in `markata.toml`:

```toml
[markata.prevnext]
# Strategy for finding prev/next posts
# 'first': Use first map where post is found
# 'all': Use all maps
strategy = 'first'

# Custom colors (optional)
prevnext_color_text = "white"
prevnext_color_text_light = "black"
prevnext_color_angle = "white"
prevnext_color_angle_light = "black"

# Navigation maps
[[markata.prevnext.maps]]
# Map posts by category
category = "tutorials"
filter = "post.category == 'tutorials'"
sort = "post.date"
reverse = true

[[markata.prevnext.maps]]
# Map posts by series
category = "python-series"
filter = "post.series == 'python'"
sort = "post.part"
reverse = false
```

## Functionality

## Navigation Maps

The plugin allows you to:
1. Define multiple navigation maps
2. Filter posts by attributes
3. Sort posts by any field
4. Control sort direction
5. Group related content

## Navigation Strategies

Two navigation modes:
- `first`: Use first map containing the post
- `all`: Use all maps containing the post

## Template Integration

Adds to each post:
- Previous post link
- Next post link
- Navigation styling
- Responsive design

## Styling

Customizable colors:
- Text colors
- Arrow colors
- Light/dark mode support
- Hover effects

## Dependencies

This plugin depends on:
- jinja2 for templating
- pydantic for configuration
"""

import copy
from dataclasses import dataclass
from pathlib import Path
from typing import TYPE_CHECKING
from typing import Dict
from typing import List
from typing import Optional

from deepmerge import always_merger
from jinja2 import Template

from markata.hookspec import hook_impl
from markata.hookspec import register_attr

if TYPE_CHECKING:
    from frontmatter import Post

    from markata import Markata

SUPPORTED_STRATEGIES = ["first", "all"]


class UnsupportedPrevNextStrategy(NotImplementedError):
    """
    A custom error class to raise when an unsupporte prevnext strategy is
    defined.
    """


@dataclass
class PrevNext:
    prev: str
    next: str


def prevnext(
    markata: "Markata",
    post: "Post",
    conf: List[Dict[str, str]],
    strategy: str = "first",
) -> Optional[PrevNext]:
    posts = []
    for map_conf in conf.values():
        _posts = markata.map("post", **map_conf)
        # if the strategy is first, cycle back to the beginning after each map
        if strategy == "first" and _posts:
            _posts.append(_posts[0])
        posts.extend(_posts)
    # if the strategy is 'all', cycle back to the beginning after all of the maps.
    if strategy == "all":
        posts.append(posts[0])

    try:
        post_idx = posts.index(post)
        return PrevNext(prev=posts[post_idx - 1], next=posts[post_idx + 1])
    except ValueError:
        # post is not in posts
        return None


TEMPLATE = (Path(__file__).parent / "prevnext_template.html").read_text()


@hook_impl
@register_attr("prevnext")
def pre_render(markata: "Markata") -> None:
    config = markata.config.get("prevnext", {})
    feed_config = markata.config.get("feeds", {})
    strategy = config.get("strategy", "first")
    if strategy not in SUPPORTED_STRATEGIES:
        msg = f"""
        "{strategy}" is not a supported prevnext strategy

        configure prevnext in your markata.toml to use one of {SUPPORTED_STRATEGIES}
        """
        raise UnsupportedPrevNextStrategy(msg)
    template = config.get("template", None)
    if template is None:
        template = Template(TEMPLATE)
    else:
        template = Template(Path(template).read_text())

    _full_config = copy.deepcopy(markata.config)
    for article in set(markata.articles):
        article["prevnext"] = prevnext(
            markata,
            article,
            feed_config,
            strategy=strategy,
        )
        if "prevnext" not in article.content and article["prevnext"]:
            article.content += template.render(
                config=always_merger.merge(
                    _full_config,
                    copy.deepcopy(
                        article.get(
                            "config_overrides",
                            {},
                        ),
                    ),
                ),
                **article,
            )
