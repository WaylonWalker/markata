"""Default load plugin."""
import sys
from pathlib import Path
from typing import TYPE_CHECKING, Callable, List, Optional

import frontmatter
import pydantic
from rich.progress import BarColumn, Progress

from markata.hookspec import hook_impl, register_attr

if TYPE_CHECKING:
    from markata import Markata, Post

    class MarkataMarkdown(Markata):
        articles: List = []


class ValidationError(ValueError):
    ...

    # def __iter__(self):
    #     return iter(self.__root__)

    # def __getitem__(self, item):
    #     return self.__root__[item]


@hook_impl
@register_attr("articles", "posts")
def load(markata: "MarkataMarkdown") -> None:
    progress = Progress(
        BarColumn(bar_width=None), transient=True, console=markata.console
    )
    if not markata.config.get("repo_url", "https://github.com/").endswith("/"):
        markata.config["repo_url"] = (
            markata.config.get("repo_url", "https://github.com/") + "/"
        )
    Posts = pydantic.create_model(
        "Posts",
        posts=(List[markata.Post], ...),
    )
    markata.posts_obj = Posts.parse_obj(
        {"posts": [get_post(article, markata) for article in markata.files]}
    )
    markata.posts = markata.posts_obj.posts
    markata.articles = markata.posts

    # markata.console.log(f"loaded {len(markata.articles)} posts")


def get_post(path: Path, markata: "Markata") -> Optional[Callable]:
    breakpoint()
    if markata.post_model:
        return pydantic_get_post(path=path, markata=markata)
    else:
        return legacy_get_post(path=path, markata=markata)


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

    # try:
    breakpoint()
    post = markata.Post(**fm_post.metadata, markata=markata)

    # except pydantic.ValidationError as e:
    #     markata.console.log(str(e).replace("Post", str(path)))
    #     sys.exit(1)

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
