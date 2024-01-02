from typing import TYPE_CHECKING, List

from more_itertools import unique_everseen
from pydantic import create_model

from markata.hookspec import hook_impl, register_attr

if TYPE_CHECKING:
    from markata import Markata


class Config:
    env_prefix = "markata_"
    extras = "allow"


class PostConfig:
    title = "Markata.Post"
    arbitrary_types_allowed = True
    copy_on_model_validation = False


@hook_impl
@register_attr("Post", "Posts", "Config")
def create_models(markata: "Markata") -> None:
    post_models = tuple(unique_everseen(markata.post_models))
    markata.Post = create_model(
        "Post",
        __base__=post_models,
    )
    markata.Posts = create_model(
        "Posts",
        posts=(List[markata.Post], ...),
    )
    markata.Post.markata = markata
    markata.Config = create_model(
        "Config",
        __base__=tuple(unique_everseen(markata.config_models)),
    )
