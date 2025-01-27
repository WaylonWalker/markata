"""Post Skipping and Caching Plugin"""

from markata.hookspec import hook_impl, register_attr
import pydantic
from rich.console import Console
from typing import Optional, TYPE_CHECKING

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


@hook_impl
def load(markata: "Markata") -> None:
    """Runs after posts are loaded to check if they should be skipped."""
    for post in markata.posts:
        if hasattr(post, "raw"):
            key = markata.make_hash("skip", post.raw)
            if markata.cache.get(key) == "done":
                post.skip = True
    console.log(
        f"{len(markata.filter('skip==True'))}/{len(markata.posts)} posts skipped"
    )
    console.log(
        f"{len(markata.filter('skip==False'))}/{len(markata.posts)} posts not skipped"
    )


@hook_impl(trylast=True)
def save(markata: "Markata") -> None:
    """Save the 'done' status for processed posts."""
    # for post in markata.posts:
    for post in markata.filter("skip==False"):
        if hasattr(post, "raw"):
            if post.output_html.exists():
                key = markata.make_hash("skip", post.raw)
                markata.cache.set(key, "done")
