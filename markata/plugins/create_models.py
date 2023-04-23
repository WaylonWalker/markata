from typing import TYPE_CHECKING

from more_itertools import unique_everseen
from pydantic import create_model

from markata.hookspec import hook_impl, register_attr

if TYPE_CHECKING:
    from markata import Markata


@hook_impl
@register_attr("Post", "Config")
def create_models(markata: "Markata") -> None:
    post_models = tuple(list(unique_everseen(markata.post_models)))
    markata.Post = create_model("Post", __base__=post_models)
    markata.Config = create_model(
        "Config", __base__=tuple(unique_everseen(markata.config_models))
    )
