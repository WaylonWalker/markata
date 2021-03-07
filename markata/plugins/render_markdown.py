from pathlib import Path

from diskcache import Cache
from pymdownx import emoji
from tqdm import tqdm

from markata.hookspec import hook_impl

MARKATA_CACHE_DIR = Path(".") / ".markata.cache"
MARKATA_CACHE_DIR.mkdir(exist_ok=True)
cache = Cache(MARKATA_CACHE_DIR)


def make_renderer(md):
    @cache.memoize(tag="markata.render_markdown.render_article", expire=15 * 24 * 60)
    def render_article(content):
        return md.convert(content)

    return render_article


@hook_impl(tryfirst=True)
def render(markata):
    for article in tqdm(
        markata.articles, desc="rendering markdown", leave=False, colour="yellow"
    ):
        key = markata.make_hash("render_markdown", "render", article["content_hash"])
        html_from_cache = markata.cache.get(key)
        if html_from_cache is None:
            html = markata.md.convert(article.content)
            markata.cache.add(key, html, expire=15 * 24 * 60)
        else:
            html = html_from_cache
        article.html = html
        article.article_html = html
