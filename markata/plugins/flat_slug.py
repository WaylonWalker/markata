"""Flat Slug Plugin

Creates a slug in article.metadata if missing based on filename.
"""
from pathlib import Path

from more_itertools import flatten
from tqdm import tqdm

from markata.hookspec import hook_impl


@hook_impl(tryfirst=True)
def render(markata):
    for article in tqdm(
        markata.articles, desc="creating slugs", leave=False, colour="yellow"
    ):
        try:
            article["slug"] = article.metadata["slug"]
        except KeyError:
            article["slug"] = Path(article["path"]).stem
