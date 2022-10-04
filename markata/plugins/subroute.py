from pathlib import Path
from typing import TYPE_CHECKING

from markata.hookspec import hook_impl

if TYPE_CHECKING:
    from markata import Markata


@hook_impl(tryfirst=True)
def pre_render(markata: "Markata") -> None:
    """
    Sets the article slug if one is not already set in the frontmatter.
    """
    subroute = markata.config.get("subroute", "")
    for article in markata.iter_articles(description="creating slugs"):
        slug = Path(article.get("slug", Path(article["path"]).stem))
        article["slug"] = slug.parent / Path(subroute) / slug.stem
