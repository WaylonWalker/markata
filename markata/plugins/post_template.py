from markata.hookspec import hook_impl
from tqdm import tqdm
from jinja2 import Template


@hook_impl
def render(markata):
    template_file = markata.config["post_template"]
    with open(template_file) as f:
        template = Template(f.read())
    for article in markata.iter_articles("apply template"):

        article.html = template.render(
            body=article.html,
            title=article.metadata["title"],
            slug=article.metadata["slug"],
            toc=markata.md.toc,
            date=article.metadata["date"],
        )
