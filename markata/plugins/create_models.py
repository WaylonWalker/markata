from typing import TYPE_CHECKING

from pydantic import create_model

from markata.hookspec import hook_impl, register_attr

if TYPE_CHECKING:
    from markata import Markata


@hook_impl
@register_attr("Post", "Config")
def create_models(markata: "Markata") -> None:
    markata.Post = create_model("Post", __base__=tuple(set(markata.post_models)))
    markata.Config = create_model("Config", __base__=tuple(set(markata.config_models)))
