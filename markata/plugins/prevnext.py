"""prevnext plugin"""
from dataclasses import dataclass
from pathlib import Path
from typing import TYPE_CHECKING, Dict, List

from jinja2 import Template

from markata.hookspec import hook_impl, register_attr

if TYPE_CHECKING:
    from frontmatter import Post

    from markata import Markata


@dataclass
class PrevNext:
    prev: str
    next: str


def prevnext(markata: "Markata", post: "Post", conf: List[Dict[str, str]]):
    posts = []
    for name, map_conf in conf.items():
        _posts = markata.map("post", **map_conf)
        # make sure the last post has a next
        if _posts:
            _posts.append(_posts[0])
        posts.extend(_posts)

    try:
        post_idx = posts.index(post)
        return PrevNext(prev=posts[post_idx - 1], next=posts[post_idx + 1])
    except ValueError:
        # post is not in posts
        return None


#     _prev = posts[-1]
#     _post = next(iposts)
#     try:
#         while _post is not post:
#             if post is _post:
#             next(iposts)
#     except StopIteration:
#         return PrevNext(prev=posts[-1], next=posts[1])

TEMPLATE = """
<div class='nextprev'>

    <style type='text/css'>

    :root {
        --prevnext-color-text: {{ config.get('prevnext_color_text', 'white') }};
        --prevnext-color-angle: {{ config.get('prevnext_color_angle', 'white') }};
        --prevnext-subtitle-brightness: 3;
    }
    [data-theme="light"] {
        --prevnext-color-text: {{ config.get('prevnext_color_text_light', '#343434') }};
        --prevnext-color-angle: {{ config.get('prevnext_color_angle_light', '#343434') }};
        --prevnext-subtitle-brightness: 3;
    }
    [data-theme="dark"] {
        --prevnext-color-text: {{ config.get('prevnext_color_text', 'white') }};
        --prevnext-color-angle: {{ config.get('prevnext_color_angle', 'white') }};
        --prevnext-subtitle-brightness: 3;
    }
    .nextprev {
    display: flex;
    justify-content: space-around;
    }
    .nextprev a {
    display: flex;
    align-items: center;
    width: 100%;
    text-decoration: none;
    }
    a.next {
    justify-content: flex-end;
    }
    .nextprev a:hover {
    background: #00000006;
    }
    .prevnext-subtitle {
    color: var(--prevnext-color-text);
    filter: brightness(var(--prevnext-subtitle-brightness));
    font-size: .8rem;
    }
    .prevnext-title {
    color: var(--prevnext-color-text);
    font-size: 1rem;
    }
    </style>
    {% if prevnext.prev['slug'] == 'index' %}
    <a class='prev' href='/'>
    {% else %}
    <a class='prev' href='/{{ prevnext.prev['slug'] }}'>
    {% endif %}

        <svg width="50px" height="50px" viewbox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M13.5 8.25L9.75 12L13.5 15.75" stroke="var(--prevnext-color-angle)" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"> </path>
        </svg>
        <div>
            <p class='prevnext-subtitle'>prev</p>
            <p class='prevnext-title'>{{ prevnext.prev['title'] }}</p>
        </div>
    </a> 
    {% if prevnext.next['slug'] == 'index' %}
    <a class='next' href='/'>
    {% else %}
    <a class='next' href='/{{ prevnext.next['slug'] }}'>
    {% endif %}
        <div>
            <p class='prevnext-subtitle'>next</p>
            <p class='prevnext-title'>{{ prevnext.next['title'] }}</p>
        </div>
        <svg width="50px" height="50px" viewbox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M10.5 15.75L14.25 12L10.5 8.25" stroke="var(--prevnext-color-angle)" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"></path>
        </svg>
    </a>
  </div>
"""


@hook_impl
@register_attr("prevnext")
def pre_render(markata: "Markata") -> None:

    config = markata.get_plugin_config("prevnext")
    template = config.get("template", None)
    if template is None:
        template = Template(TEMPLATE)
    else:
        template = Template(Path(template).read_text())
    for article in set(markata.articles):
        article["prevnext"] = prevnext(markata, article, config.get("maps", {}))
        if "prevnext" not in article.content and article["prevnext"]:
            article.content += template.render(config=config, **article)
