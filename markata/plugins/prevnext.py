"""
The prevnext plugin, creates previous and next links inside each post.

## Example config

In this example we have two maps of posts to look through.  prevnext will look
through each of these lists of posts for the current post, then return the post
before and after this post as the prevnext posts.

``` toml

[markata]
# default colors will be taken from markata's color_text and color_accent
color_text=white
color_text_light=black
color_accent=white
color_accent_light=black

[markata.prevnext]
# strategy can be 'first' or 'all'
# 'first' will cycle through the first map the post is found in.
# 'all' will cycle through all of the maps
strategy='first'

# if you want different colors than your main color_text and color_accent, then
# you can override it here
# colors can be any valid css color format

prevnext_color_text=white
prevnext_color_text_light=black
prevnext_color_angle=white
prevnext_color_angle_light=black


# you can have multiple maps, the order they appear will determine their preference
[markata.feeds.python]
filter='"python" in tags'
sort='slug'

[markata.feeds.others]
filter='"python" not in tags'
sort='slug'
```

The configuration below will setup two maps, one where "python" is in the list
of tags, and another where it is not.  This will link all python posts together
with a prevnext cycle, and all non-python posts in a separate prevnext cycle.

## strategy

There are currently two supported strategies.

* first
* all

### first

`first` will cycle through only the posts contained within the first map that
contains the post.

### all

`all` will cycle through all of the posts aggregated from any prevnext map.

"""
import copy
from dataclasses import dataclass
from pathlib import Path
from typing import TYPE_CHECKING, Dict, List, Optional

from deepmerge import always_merger
from jinja2 import Template

from markata.hookspec import hook_impl, register_attr

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
