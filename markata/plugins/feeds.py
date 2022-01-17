import datetime
import shutil
import textwrap
from pathlib import Path

from jinja2 import Template

from markata.hookspec import hook_impl


class MarkataFilterError(RuntimeError):
    ...


@hook_impl
def save(markata):
    config = markata.get_plugin_config("feeds")
    if "feeds" not in config.keys():
        config["feeds"] = dict()
    if "archive" not in config["feeds"].keys():
        config["feeds"]["archive"] = dict()
        config["feeds"]["archive"][
            "filter"
        ] = "date<=today and status.lower()=='published'"

    description = markata.get_config("description") or ""
    url = markata.get_config("url") or ""
    title = markata.get_config("title") or ""

    template = Path(__file__).parent / "default_post_template.html"

    for page, page_conf in config["feeds"].items():
        create_page(
            markata,
            page,
            description=description,
            url=url,
            title=title,
            template=template,
            **page_conf,
        )

    home = Path(markata.config["output_dir"]) / "index.html"
    archive = Path(markata.config["output_dir"]) / "archive" / "index.html"
    if not home.exists() and archive.exists():
        shutil.copy(str(archive), str(home))


def create_page(
    markata,
    page,
    tags=None,
    status="published",
    template=None,
    card_template=None,
    filter=None,
    description=None,
    url=None,
    title=None,
):
    all_posts = reversed(sorted(markata.articles, key=lambda x: x["date"]))

    if filter is not None:
        posts = reversed(sorted(markata.articles, key=lambda x: x["date"]))
        try:
            posts = [post for post in posts if eval(filter, post.to_dict(), {})]
        except BaseException as e:
            msg = textwrap.dedent(
                f"""
                    While processing {page =} markata hit the following exception
                    during {filter =}
                    {e}
                    """
            )
            raise MarkataFilterError(msg)

    cards = [create_card(post, card_template) for post in posts]

    with open(template) as f:
        template = Template(f.read())
    output_file = Path(markata.config["output_dir"]) / page / "index.html"
    canonical_url = f"{url}/{page}/"
    output_file.parent.mkdir(exist_ok=True, parents=True)

    with open(output_file, "w+") as f:
        f.write(
            template.render(
                body="".join(cards),
                url=url,
                description=description,
                title=title,
                canonical_url=canonical_url,
                today=datetime.datetime.today(),
            )
        )


def create_card(post, template=None):
    if template is None:
        return textwrap.dedent(
            f"""
            <li class='post'>
            <a href="/{post['slug']}/">
                {post['title']} {post['date'].year}-{post['date'].month}-{post['date'].day}
            </a>
            </li>
            """
        )
    with open(template) as f:
        template = Template(f.read())
    post["article_html"] = post.article_html

    return template.safe_substitute(**post.to_dict())
