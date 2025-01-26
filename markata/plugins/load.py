"""Default load plugin."""

import itertools
from pathlib import Path
from typing import TYPE_CHECKING, Callable, List, Optional
from concurrent.futures import as_completed, ThreadPoolExecutor

from markata import background
import asyncio
import frontmatter
import pydantic
from rich.progress import BarColumn, Progress
from yaml.parser import ParserError

from markata.hookspec import hook_impl, register_attr

if TYPE_CHECKING:
    from markata import Markata, Post

    class MarkataMarkdown(Markata):
        articles: List = []


class ValidationError(ValueError): ...


def load_file_content(path: Path) -> tuple[Path, dict]:
    """Load file content without validation."""
    try:
        with open(path, 'r') as f:
            content = f.read()
        return path, frontmatter.loads(content)
    except Exception as e:
        return path, None


@hook_impl
@register_attr("articles", "posts")
def load(markata: "MarkataMarkdown") -> None:
    Progress(
        BarColumn(bar_width=None),
        transient=True,
        console=markata.console,
    )
    markata.console.log(f"found {len(markata.files)} posts")
    
    # Use ThreadPoolExecutor to load files in parallel
    with ThreadPoolExecutor() as executor:
        # Load all file contents first
        file_contents = list(executor.map(load_file_content, markata.files))
    
    posts = []
    errors = []
    
    # Bulk process and validate posts
    for path, content in file_contents:
        if content is None:
            continue
            
        try:
            if markata.Post:
                # Create post using model_validate to get proper type coercion
                post = markata.Post.model_validate(
                    {
                        "markata": markata,
                        "path": path,
                        "content": content.content,
                        **content.metadata
                    },
                    context={"markata": markata}
                )
                posts.append(post)
            else:
                post = legacy_get_post(path=path, markata=markata)
                if post is not None:
                    posts.append(post)
        except Exception as e:
            errors.append((path, str(e)))
            
    if errors:
        for path, error in errors:
            markata.console.log(f"Error loading {path}: {error}")

    markata.posts_obj = markata.Posts.parse_obj(
        {"posts": posts},
    )
    markata.posts = markata.posts_obj.posts
    markata.articles = markata.posts


@background.task
def get_post(path: Path, markata: "Markata") -> Optional[Callable]:
    if markata.Post:
        post = pydantic_get_post(path=path, markata=markata)
        return post
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
    try:
        post = markata.Post.parse_file(markata=markata, path=path)
        markata.Post.validate(post)

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
        markata.config.repo_url + "edit/" + markata.config.repo_branch + "/" + post.path
    )
    return post
