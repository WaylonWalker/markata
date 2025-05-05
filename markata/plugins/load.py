"""
The `markata.plugins.load` plugin is responsible for loading and parsing markdown files
with frontmatter into Post objects. It provides parallel loading capabilities and
handles both modern Pydantic-based and legacy frontmatter validation.

## Installation

This plugin is built-in and enabled by default through the 'default' plugin.
If you want to be explicit, you can add it to your list of plugins:

```toml
hooks = [
    "markata.plugins.load",
]
```

## Uninstallation

Since this plugin is included in the default plugin set, to disable it you must explicitly
add it to the disabled_hooks list if you are using the 'default' plugin:

```toml
disabled_hooks = [
    "markata.plugins.load",
]
```

Note: Disabling this plugin will prevent Markata from loading any markdown files.
This will effectively disable most of Markata's functionality.

## Configuration

Configure loading behavior in your `markata.toml`:

```toml
[markata]
# Directories containing markdown content
content_directories = [
    "content",
    "posts"
]

# Optional: Set to true to use legacy frontmatter validation
legacy_frontmatter = false

# Optional: Number of worker processes for parallel loading
load_workers = 4
```

## Functionality

## File Loading

The plugin:
1. Discovers markdown files in content directories
2. Loads file content and frontmatter
3. Validates frontmatter against Post model
4. Creates Post objects for further processing

## Parallel Processing

Loading is parallelized using:
- Process pool for file reading
- Configurable number of workers
- Chunked file processing

## Validation Modes

Supports two validation approaches:
1. Modern Pydantic-based validation (default)
   - Strict type checking
   - Automatic type coercion
   - Detailed validation errors

2. Legacy frontmatter validation
   - Looser type checking
   - Compatible with older content
   - Less strict validation

## Error Handling

The plugin provides:
- Detailed validation error messages
- Per-file error reporting
- Graceful fallback to legacy mode
- Optional strict validation

## Registered Attributes

The plugin adds:
- `articles`: List of loaded Post objects
- `content_directories`: List of content source directories

## Dependencies

This plugin depends on:
- python-frontmatter for YAML parsing
- pydantic for validation
- multiprocessing for parallel loading
"""

import itertools
from concurrent.futures import ThreadPoolExecutor
from pathlib import Path
from typing import TYPE_CHECKING
from typing import Callable
from typing import List
from typing import Optional

import frontmatter
import pydantic
from rich.progress import BarColumn
from rich.progress import Progress
from yaml.parser import ParserError

from markata import background
from markata.hookspec import hook_impl
from markata.hookspec import register_attr

if TYPE_CHECKING:
    from markata import Markata
    from markata import Post

    class MarkataMarkdown(Markata):
        articles: List = []


class ValidationError(ValueError): ...


def load_file_content(path: Path) -> tuple[Path, dict]:
    """Load file content without validation."""
    try:
        with open(path, "r") as f:
            raw_content = f.read()
        try:
            content = frontmatter.loads(raw_content)
        except Exception:
            content = None
        content["raw"] = raw_content
        return path, content
    except Exception:
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
                        **content.metadata,
                    },
                    context={"markata": markata},
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
