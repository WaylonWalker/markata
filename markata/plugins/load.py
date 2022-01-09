"""Default load plugin."""
import time
from pathlib import Path
from typing import TYPE_CHECKING, Callable, Optional

import frontmatter
from rich.progress import BarColumn, Progress
from yaml.parser import ParserError

from markata.background import task
from markata.hookspec import hook_impl, register_attr

if TYPE_CHECKING:
    from markata import Markata, Post


@hook_impl
@register_attr("articles")
def load(markata: "Markata") -> None:
    progress = Progress(
        BarColumn(bar_width=None), transient=True, console=markata.console
    )

    futures = [get_post(article, markata) for article in markata.files]
    task_id = progress.add_task("loading markdown")
    progress.update(task_id, total=len(futures))
    with progress:
        while not all([f.done() for f in futures]):
            time.sleep(0.1)
            progress.update(task_id, total=len([f for f in futures if f.done()]))
    articles = [f.result() for f in futures]
    articles = [a for a in articles if a]
    markata.articles = articles


@task
def get_post(path: Path, markata: "Markata") -> Optional[Callable]:
    # -> Optional["Post"]:
    default = {
        "cover": "",
        "title": "",
        "tags": [],
        "status": "draft",
        "templateKey": "",
        "path": str(path),
        "description": "",
        "content": "",
    }
    try:
        post: "Post" = frontmatter.load(path)
        post.metadata = {**default, **post.metadata}
    except ParserError:
        return None
        post = default
    except ValueError:
        return None
        post = default
    post.metadata["path"] = str(path)
    post["content_hash"] = markata.make_hash(post.content)
    return post
