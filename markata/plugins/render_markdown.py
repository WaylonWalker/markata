import copy
from typing import List, TYPE_CHECKING

import markdown
from markdown_it import MarkdownIt
from mdit_py_plugins.admon import admon_plugin
from mdit_py_plugins.footnote import footnote_plugin
from mdit_py_plugins.front_matter import front_matter_plugin

from markata import DEFAULT_MD_EXTENSIONS
from markata.hookspec import hook_impl, register_attr
import importlib

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

    if (
        markata.config.get("markdown_backend", "")
        .lower()
        .replace(" ", "-")
        .replace("_", "-")
        == "markdown-it-py"
    ):
        markata.md = MarkdownIt().disable("image").enable("table")
        plugins = (front_matter_plugin, footnote_plugin, admon_plugin)
        plugins = (getattr(importlib.import_module(".".join(plugin.split(".")[:-1])), plugin.split('.')[-1]) for plugin in ('mdit_py_plugins.admon.admon_plugin',
                                                                                                                            'mdit_py_plugins.footnote.footnote_plugin',
                                                                                                                            'mdit_py_plugins.front_matter.front_matter_plugin',
                                                                                                                            ))
        for plugin in plugins:
            markata.md = markata.md.use(plugin)
        markata.md.convert = markata.md.render
        markata.md.toc = ''
    else:
        markata.md = markdown.Markdown(extensions=markata.markdown_extensions)


@hook_impl(tryfirst=True)
@register_attr("articles")
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
            article.article_html = copy.deepcopy(html)

            article["html"] = html
            article["article_html"] = article.article_html
