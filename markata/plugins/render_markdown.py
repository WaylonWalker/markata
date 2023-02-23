from pygments.lexers import get_lexer_by_name
from pygments.formatters import HtmlFormatter
from pygments import highlight
import copy
from typing import List, TYPE_CHECKING


from markata import DEFAULT_MD_EXTENSIONS
from markata.hookspec import hook_impl, register_attr
import importlib

COPY_ICON = '<svg version="1.1" id="Layer_1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" x="0px" y="0px" viewBox="0 0 115.77 122.88" style="enable-background:new 0 0 115.77 122.88" xml:space="preserve"><style type="text/css">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path class="st0" d="M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02 v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02 c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1 c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7 h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01 c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65 v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01 c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02 v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z"/></g></svg>'


if TYPE_CHECKING:
    from markata import Markata

    class MarkataMarkdown(Markata):
        articles: List = []
        md: markdown.Markdown = markdown.Markdown()
        markdown_extensions: List = []


def highlight_code(code, name, attrs):
    """Highlight a block of code"""

    import rich
    if attrs:
        rich.print(f"Ignoring {attrs=}")

    lexer = get_lexer_by_name(name or 'text')
    # rich.print(highlight(code, lexer, formatter))
    import re
    pattern = r'(\w+)\s*=\s*(".*?"|\S+)'
    matches = re.findall(pattern, attrs)
    attrs = dict(matches)

    if attrs.get('hl_lines'):
        formatter = HtmlFormatter(hl_lines=attrs.get('hl_lines'))
    else:
        formatter = HtmlFormatter()

    copy_button = f'''<button class='copy' title='copy code to clipboard' onclick="navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)">{COPY_ICON}</button>'''

    if attrs.get('title'):
        file = f'''
<div class='filepath'>
{attrs.get('title')}
{copy_button}
</div>
'''
    else:
        file = f'''
<div class='copy-wrapper'>
{copy_button}
</div>
        '''
    return f'''<pre class='wrapper'>
{file}
{highlight(code, lexer, formatter)}
</pre>
'''


@ hook_impl(tryfirst=True)
@ register_attr("md", "markdown_extensions")
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
        from markdown_it import MarkdownIt
        markata.md = MarkdownIt(
            "js-default",
            {
                "linkify": True,
                "html": True,
                "typographer": True,
                "highlight": highlight_code,
            },
        ).disable("image").enable("table")
        plugins = (getattr(importlib.import_module(plugin.split(":")[0]), plugin.split(':')[1]) for plugin in (
            'mdit_py_plugins.admon:admon_plugin',
            'markata.plugins.mdit_details:details_plugin',
        ))
        for plugin in plugins:
            markata.md = markata.md.use(plugin)
        markata.md.convert = markata.md.render
        markata.md.toc = ''
    elif (
        markata.config.get("markdown_backend", "")
        .lower()
        .replace(" ", "-")
        .replace("_", "-")
        == "markdown2"
    ):
        import markdown2
        markata.md = markdown2.Markdown(extras=markata.markdown_extensions)
        markata.md.toc = ''
    else:
        import markdown
        markata.md = markdown.Markdown(extensions=markata.markdown_extensions)


@ hook_impl(tryfirst=True)
@ register_attr("articles")
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
