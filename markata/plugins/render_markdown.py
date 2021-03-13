from pathlib import Path

from markata.hookspec import hook_impl

from typing import TYPE_CHECKING

if TYPE_CHECKING:
    from markata import Markata


@hook_impl(tryfirst=True)
def render(markata: "Markata") -> None:
    for article in markata.iter_articles("rendering markdown"):
        key = markata.make_hash("render_markdown", "render", article["content_hash"])
        html_from_cache = markata.cache.get(key)
        if html_from_cache is None:
            html = markata.md.convert(article.content)
            markata.cache.add(key, html, expire=15 * 24 * 60)
        else:
            html = html_from_cache
        article.html = html
        article.article_html = html
