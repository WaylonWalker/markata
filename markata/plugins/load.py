"""Default load plugin."""
import itertools
from pathlib import Path
from typing import Callable, List, Optional, TYPE_CHECKING

import frontmatter
import pydantic
from rich.progress import BarColumn, Progress
from yaml.parser import ParserError

from markata.hookspec import hook_impl, register_attr

if TYPE_CHECKING:
    from markata import Markata, Post

    class MarkataMarkdown(Markata):
        articles: List = []


class ValidationError(ValueError):
    ...


@hook_impl
@register_attr("articles", "posts")
def load(markata: "MarkataMarkdown") -> None:
    Progress(
        BarColumn(bar_width=None),
        transient=True,
        console=markata.console,
    )
    if not markata.config.get("repo_url", "https://github.com/").endswith("/"):
        markata.config.repo_url = (
            markata.config.get("repo_url", "https://github.com/") + "/"
        )
    Posts = pydantic.create_model(
        "Posts",
        posts=(List[markata.Post], ...),
    )
    markata.console.log(f"found {len(markata.files)} posts")
    markata.posts_obj = Posts.parse_obj(
        {"posts": [get_post(article, markata) for article in markata.files]},
    )
    markata.posts = markata.posts_obj.posts
    markata.articles = markata.posts


def get_post(path: Path, markata: "Markata") -> Optional[Callable]:
    if markata.Post:
        return pydantic_get_post(path=path, markata=markata)
    else:
        return legacy_get_post(path=path, markata=markata)


def get_models(markata: "Markata", error: pydantic.ValidationError) -> List:
    fields = []
    for err in error.errors():
        fields.extend(err["loc"])

    models = {field: f"{field} used by " for field in fields}

    for field, model in set(
        itertools.product(
            fields,
            markata.post_models,
        ),
    ):
        if field in model.__fields__:
            models[field] += f"'{model.__module__}.{model.__name__}'"

    return models


def pydantic_get_post(path: Path, markata: "Markata") -> Optional[Callable]:
    fm_post = frontmatter.load(path)
    fm_post["content"] = fm_post.content
    fm_post["path"] = str(path)
    fm_post["edit_link"] = (
        str(markata.config.get("repo_url", "https://github.com/"))
        + "edit/"
        + str(markata.config.get("repo_branch", "main"))
        + "/"
        + str(path)
    )

    try:
        post = markata.Post(**fm_post.metadata, markata=markata)

    except pydantic.ValidationError as e:

        models = get_models(markata=markata, error=e)
        models = list(models.values())
        models = "\n".join(models)
        raise ValidationError(f"{e}\n\n{models}\nfailed to load {path}") from e

    return post


def legacy_get_post(path: Path, markata: "Markata") -> Optional[Callable]:
    default = {
        "cover": "",
        "title": "",
        "tags": [],
        "published": "False",
        "templateKey": "",
        "path": str(path),
        "description": "",
        "content": "",
    }
    try:
        post: "Post" = frontmatter.load(path)
        post.metadata = {**default, **post.metadata}
        post["content"] = post.content
    except ParserError:
        return None
        post = default
    except ValueError:
        return None
        post = default
    post.metadata["path"] = str(path)
    post["edit_link"] = (
        str(markata.config.get("repo_url", "https://github.com/"))
        + "edit/"
        + str(markata.config.get("repo_branch", "main"))
        + "/"
        + str(post["path"])
    )
    return post
