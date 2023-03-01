"""
Renders markdown content as html.  This may be markdown files loaded in by way
of the [[load]] plugin.


"""
import copy
import importlib
from typing import List, TYPE_CHECKING

import markdown
from pygments import highlight
from pygments.formatters import HtmlFormatter
from pygments.lexers import get_lexer_by_name

from markata import DEFAULT_MD_EXTENSIONS
from markata.hookspec import hook_impl, register_attr

COPY_ICON = '<svg version="1.1" id="Layer_1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" x="0px" y="0px" viewBox="0 0 115.77 122.88" style="enable-background:new 0 0 115.77 122.88" xml:space="preserve"><style type="text/css">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path class="st0" d="M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02 v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02 c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1 c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7 h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01 c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65 v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01 c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02 v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z"/></g></svg>'
HELP_ICON = '<svg height="92px" id="Capa_1" style="enable-background:new 0 0 91.999 92;" version="1.1" viewBox="0 0 91.999 92" width="91.999px" xml:space="preserve" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink"><path d="M45.385,0.004C19.982,0.344-0.334,21.215,0.004,46.619c0.34,25.393,21.209,45.715,46.611,45.377  c25.398-0.342,45.718-21.213,45.38-46.615C91.655,19.986,70.785-0.335,45.385,0.004z M45.249,74l-0.254-0.004  c-3.912-0.116-6.67-2.998-6.559-6.852c0.109-3.788,2.934-6.538,6.717-6.538l0.227,0.004c4.021,0.119,6.748,2.972,6.635,6.937  C51.903,71.346,49.122,74,45.249,74z M61.704,41.341c-0.92,1.307-2.943,2.93-5.492,4.916l-2.807,1.938  c-1.541,1.198-2.471,2.325-2.82,3.434c-0.275,0.873-0.41,1.104-0.434,2.88l-0.004,0.451H39.429l0.031-0.907  c0.131-3.728,0.223-5.921,1.768-7.733c2.424-2.846,7.771-6.289,7.998-6.435c0.766-0.577,1.412-1.234,1.893-1.936  c1.125-1.551,1.623-2.772,1.623-3.972c0-1.665-0.494-3.205-1.471-4.576c-0.939-1.323-2.723-1.993-5.303-1.993  c-2.559,0-4.311,0.812-5.359,2.478c-1.078,1.713-1.623,3.512-1.623,5.35v0.457H27.935l0.02-0.477  c0.285-6.769,2.701-11.643,7.178-14.487C37.946,18.918,41.446,18,45.53,18c5.346,0,9.859,1.299,13.412,3.861  c3.6,2.596,5.426,6.484,5.426,11.556C64.368,36.254,63.472,38.919,61.704,41.341z"/><g/><g/><g/><g/><g/><g/><g/><g/><g/><g/><g/><g/><g/><g/><g/></svg>'
LINK_ICON = '<svg class="heading-permalink" aria-hidden="true" fill="currentColor" focusable="false" height="1em" viewBox="0 0 24 24" width="1em" xmlns="http://www.w3.org/2000/svg"><path d="M9.199 13.599a5.99 5.99 0 0 0 3.949 2.345 5.987 5.987 0 0 0 5.105-1.702l2.995-2.994a5.992 5.992 0 0 0 1.695-4.285 5.976 5.976 0 0 0-1.831-4.211 5.99 5.99 0 0 0-6.431-1.242 6.003 6.003 0 0 0-1.905 1.24l-1.731 1.721a.999.999 0 1 0 1.41 1.418l1.709-1.699a3.985 3.985 0 0 1 2.761-1.123 3.975 3.975 0 0 1 2.799 1.122 3.997 3.997 0 0 1 .111 5.644l-3.005 3.006a3.982 3.982 0 0 1-3.395 1.126 3.987 3.987 0 0 1-2.632-1.563A1 1 0 0 0 9.201 13.6zm5.602-3.198a5.99 5.99 0 0 0-3.949-2.345 5.987 5.987 0 0 0-5.105 1.702l-2.995 2.994a5.992 5.992 0 0 0-1.695 4.285 5.976 5.976 0 0 0 1.831 4.211 5.99 5.99 0 0 0 6.431 1.242 6.003 6.003 0 0 0 1.905-1.24l1.723-1.723a.999.999 0 1 0-1.414-1.414L9.836 19.81a3.985 3.985 0 0 1-2.761 1.123 3.975 3.975 0 0 1-2.799-1.122 3.997 3.997 0 0 1-.111-5.644l3.005-3.006a3.982 3.982 0 0 1 3.395-1.126 3.987 3.987 0 0 1 2.632 1.563 1 1 0 0 0 1.602-1.198z"></path></svg>'


if TYPE_CHECKING:
    from markata import Markata

    class MarkataMarkdown(Markata):
        articles: List = []
        md: markdown.Markdown = markdown.Markdown()
        markdown_extensions: List = []


def highlight_code(code, name, attrs, markata=None):
    """Code highlighter for markdown-it-py."""

    lexer = get_lexer_by_name(name or "text")
    import re

    pattern = r'(\w+)\s*=\s*(".*?"|\S+)'
    matches = re.findall(pattern, attrs)
    attrs = dict(matches)

    if attrs.get("hl_lines"):
        formatter = HtmlFormatter(hl_lines=attrs.get("hl_lines"))
    else:
        formatter = HtmlFormatter()

    copy_button = f"""<button class='copy' title='copy code to clipboard' onclick="navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)">{COPY_ICON}</button>"""

    from markdown_it import MarkdownIt

    md = MarkdownIt(
        "commonmark",
        {
            "html": True,
            "typographer": True,
        },
    )

    if attrs.get("help"):
        help = f"""
        <a href={attrs.get('help').strip('<').strip('>').strip('"').strip("'")} title='help link' class='help'>{HELP_ICON}</a>
        """
    else:
        help = ""
    if attrs.get("title"):
        file = f"""
<div class='filepath'>
{md.render(attrs.get('title').strip('"').strip("'"))}
<div class='right'>
{help}
{copy_button}
</div>
</div>
"""
    else:
        file = f"""
<div class='copy-wrapper'>
{help}
{copy_button}
</div>
        """
    return f"""<pre class='wrapper'>
{file}
{highlight(code, lexer, formatter)}
</pre>
"""


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

        markata.md = (
            MarkdownIt(
                "gfm-like",
                {
                    "linkify": True,
                    "html": True,
                    "typographer": True,
                    "highlight": highlight_code,
                },
            )
            .disable("image")
            .enable("table")
        )
        plugins = (
            (
                getattr(
                    importlib.import_module(plugin[0].split(":")[0]),
                    plugin[0].split(":")[1],
                ),
                plugin[1],
            )
            for plugin in (
                ("mdit_py_plugins.admon:admon_plugin", {}),
                ("mdit_py_plugins.attrs:attrs_plugin", {"spans": True}),
                ("mdit_py_plugins.attrs:attrs_block_plugin", {}),
                ("markata.plugins.mdit_details:details_plugin", {}),
                (
                    "mdit_py_plugins.anchors:anchors_plugin",
                    {"permalink": True, "permalinkSymbol": LINK_ICON},
                ),
                (
                    "markata.plugins.md_it_wikilinks:wikilinks_plugin",
                    {"markata": markata},
                ),
            )
        )
        for plugin in plugins:
            if plugin[1] == {}:
                markata.md = markata.md.use(plugin[0])
            else:
                markata.md = markata.md.use(plugin[0], **plugin[1])
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
