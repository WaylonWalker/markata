from typing import Dict, TYPE_CHECKING

from more_itertools import unique_everseen
import pydantic
from pydantic import create_model

from markata.hookspec import hook_impl, register_attr

if TYPE_CHECKING:
    from markata import Markata


class Config:
    env_prefix = "markata_"
    extras = "allow"


@hook_impl
@register_attr("Post", "Config")
def create_models(markata: "Markata") -> None:
    @pydantic.validator("markata", pre=True, always=True)
    def default_markata(cls, v, *, values: Dict) -> "Markata":
        if v is None:
            return markata

    post_models = tuple(unique_everseen(markata.post_models))
    markata.Post = create_model(
        "Post", __base__=post_models, __validators__={"markata": default_markata}
    )
    markata.Config = create_model(
        "Config",
        __base__=tuple(unique_everseen(markata.config_models)),
    )
