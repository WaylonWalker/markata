from pathlib import Path

from markata.hookspec import hook_impl, register_attr


@hook_impl
@register_attr("articles", "posts")
def pre_render(markata) -> None:
    for article in markata.filter('title==""'):
        article["title"] = (
            Path(article["path"]).stem.replace("-", " ").replace("_", " ").title()
        )
