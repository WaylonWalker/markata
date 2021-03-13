"""Flat Slug Plugin

Creates a slug in article.metadata if missing based on filename.
"""
from pathlib import Path

from markata.hookspec import hook_impl

from typing import TYPE_CHECKING

if TYPE_CHECKING:
    from markata import Markata


@hook_impl(tryfirst=True)
def render(markata: "Markata") -> None:
    for article in markata.iter_articles(description="creating slugs"):
        try:
            article["slug"] = article.metadata["slug"]
        except KeyError:
            article["slug"] = Path(article["path"]).stem
