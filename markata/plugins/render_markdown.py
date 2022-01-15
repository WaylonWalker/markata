from typing import TYPE_CHECKING, List

import markdown

from markata import DEFAULT_MD_EXTENSIONS
from markata.hookspec import hook_impl, register_attr

if TYPE_CHECKING:
    from markata import Markata

    class MarkataMarkdown(Markata):
        articles: List = []
        md: markdown.Markdown = markdown.Markdown()
        markdown_extensions: List = []


@hook_impl(tryfirst=True)
@register_attr("md", "markdown_extensions")
def configure(markata: "MarkataMarkdown") -> None:
    if "markdown_extensions" not in markata.config:
        markdown_extensions = [""]
    if isinstance(markata.config["markdown_extensions"], str):
        markdown_extensions = [markata.config["markdown_extensions"]]
    if isinstance(markata.config["markdown_extensions"], list):
        markdown_extensions = markata.config["markdown_extensions"]
    else:
        raise TypeError("markdown_extensions should be List[str]")

    markata.markdown_extensions = [*DEFAULT_MD_EXTENSIONS, *markdown_extensions]
    markata.md = markdown.Markdown(extensions=markata.markdown_extensions)


@hook_impl(tryfirst=True)
def render(markata: "Markata") -> None:
    config = markata.get_plugin_config(__file__)
    with markata.cache as cache:
        for article in markata.iter_articles("rendering markdown"):
            key = markata.make_hash(
                "render_markdown",
                "render",
                article.content,
            )
            html_from_cache = cache.get(key)
            if html_from_cache is None:
                html = markata.md.convert(article.content)
                cache.add(key, html, expire=config["cache_expire"])
            else:
                html = html_from_cache
            article.html = html
            article.article_html = html
