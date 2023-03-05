"""
Renders markdown content as html.  This may be markdown files loaded in by way
of the [[load]] plugin.


"""
import copy
import importlib
from typing import List, TYPE_CHECKING

import markdown

from markata import DEFAULT_MD_EXTENSIONS
from markata.hookspec import hook_impl, register_attr
from markata.plugins.md_it_highlight_code import highlight_code

if TYPE_CHECKING:
    from markata import Markata

    class MarkataMarkdown(Markata):
        articles: List = []
        md: markdown.Markdown = markdown.Markdown()
        markdown_extensions: List = []


@hook_impl(tryfirst=True)
@register_attr("md", "markdown_extensions")
def configure(markata: "MarkataMarkdown") -> None:
    "Sets up a markdown instance as md"
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
        from markdown_it import MarkdownIt

        config_update = markata.config.get("markdown_it_py", {}).get(
            "options_update",
            {
                "linkify": True,
                "html": True,
                "typographer": True,
                "highlight": highlight_code,
            },
        )
        if isinstance(config_update.get("highlight"), str):
            module = config_update["highlight"].split(":")[0]
            func = config_update["highlight"].split(":")[1]
            config_update["highlight"] = getattr(
                importlib.import_module(module),
                func,
            )

        markata.md = MarkdownIt(
            markata.config.get("markdown_it_py", {}).get("config", "gfm-like"),
            config_update,
        )
        for plugin in markata.config.get("markdown_it_py", {}).get("enable", []):
            markata.md.enable(plugin)
        for plugin in markata.config.get("markdown_it_py", {}).get("disable", []):
            markata.md.disable(plugin)

        plugins = copy.deepcopy(
            markata.config.get("markdown_it_py", {}).get("plugins", [])
        )
        for plugin in plugins:
            if isinstance(plugin["plugin"], str):
                plugin["plugin_str"] = plugin["plugin"]
                plugin_module = plugin["plugin"].split(":")[0]
                plugin_func = plugin["plugin"].split(":")[1]
                plugin["plugin"] = getattr(
                    importlib.import_module(plugin_module),
                    plugin_func,
                )
            plugin["config"] = plugin.get("config", {})
            for k, v in plugin["config"].items():
                if k == "markata":
                    plugin["config"][k] = markata

            markata.md = markata.md.use(plugin["plugin"], **plugin["config"])

        markata.md.convert = markata.md.render
        markata.md.toc = ""
    elif (
        markata.config.get("markdown_backend", "")
        .lower()
        .replace(" ", "-")
        .replace("_", "-")
        == "markdown2"
    ):
        import markdown2

        markata.md = markdown2.Markdown(extras=markata.markdown_extensions)
        markata.md.toc = ""
    else:
        import markdown

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
