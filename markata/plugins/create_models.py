from pydantic import create_model

from markata.hookspec import hook_impl, register_attr


@hook_impl
@register_attr("Post")
def create_models(markata: "Markata") -> None:
    markata.Post = create_model("Post", __base__=tuple(set(markata.post_models)))
