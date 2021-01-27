from markata.hookspec import hook_impl
import markdown
from tqdm import tqdm

from pymdownx import emoji

default_extensions = [
    "markdown.extensions.toc",
    "markdown.extensions.admonition",
    "markdown.extensions.tables",
    "markdown.extensions.md_in_html",
    "pymdownx.magiclink",
    "pymdownx.betterem",
    "pymdownx.tilde",
    "pymdownx.emoji",
    "pymdownx.tasklist",
    "pymdownx.superfences",
    "pymdownx.highlight",
    "pymdownx.inlinehilite",
    "pymdownx.keys",
    "pymdownx.saneheaders",
    # "codehilite",
]


@hook_impl(tryfirst=True)
def render(markata):
    for article in tqdm(
        markata.articles, desc="rendering markdown", leave=False, colour="yellow"
    ):
        extensions = markata.config["markdown_extensions"]
        md = markdown.Markdown(extensions=[*default_extensions, *extensions])
        html = md.convert(article.content)
        article.html = html
        article.md = md
        article.article_html = html
