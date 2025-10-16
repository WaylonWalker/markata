"""Post Skipping and Caching Plugin"""

import os
from typing import TYPE_CHECKING
from typing import Optional

import pydantic
from rich.console import Console

from markata.hookspec import hook_impl
from markata.hookspec import register_attr

if TYPE_CHECKING:
    from markata import Markata

console = Console()

MARKATA_PLUGIN_NAME = "Skip"
MARKATA_PLUGIN_PACKAGE_NAME = "skip"


class PostModel(pydantic.BaseModel):
    skip: Optional[bool] = False


@hook_impl
@register_attr("post_models")
def post_model(markata: "Markata") -> None:
    """Add skip attribute to post models."""
    markata.post_models.append(PostModel)


@hook_impl(trylast=True)
def load(markata: "Markata") -> None:
    """Runs after posts are loaded to check if they should be skipped."""
    raw_should_skip = os.environ.get("MARKATA_SKIP", "true")
    should_skip = raw_should_skip.lower() in ["true", "1", "t", "y", "yes", "on"]
    if not should_skip:
        return
    for post in markata.posts:
        if hasattr(post, "raw"):
            key = markata.make_hash("skip", post.raw)
            if markata.cache.get(key) == "done" and not post.get("jinja", False):
                post.skip = True
    markata.console.log(
        f"{len(markata.filter('skip'))}/{len(markata.posts)} posts skipped"
    )
    markata.console.log(
        f"{len(markata.filter('not skip'))}/{len(markata.posts)} posts not skipped"
    )


@hook_impl(trylast=True)
def save(markata: "Markata") -> None:
    """Save the 'done' status for processed posts."""
    # for post in markata.posts:
    raw_should_skip = os.environ.get("MARKATA_SKIP", "true")
    should_skip = raw_should_skip.lower() in ["true", "1", "t", "y", "yes", "on"]
    if not should_skip:
        return
    for post in markata.filter("not skip"):
        if hasattr(post, "raw"):
            if post.output_html.exists():
                key = markata.make_hash("skip", post.raw)
                markata.cache.set(key, "done")
