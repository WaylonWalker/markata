"""Markata's hook specification system for plugin development.

# Overview

Markata uses pluggy to define hooks that plugins can implement. These hooks allow plugins
to modify Markata's behavior at specific points in the build process.

# Hook Types

## Configuration Hooks

Used to set up plugin configuration and models:

```python
from markata.hookspec import hook_impl, register_attr

@hook_impl
@register_attr("config")
def config_model(markata):
    \"\"\"Add plugin-specific config.\"\"\"
    from pydantic import BaseModel

    class MyConfig(BaseModel):
        enabled: bool = True
        output_file: str = "output.html"

    return {"my_plugin": MyConfig()}

@hook_impl
@register_attr("my_data")
def configure(markata):
    \"\"\"Initialize plugin using config.\"\"\"
    if markata.config.my_plugin.enabled:
        # Set up plugin resources
        markata.my_data = []
```

## Content Model Hooks

Define how content is structured:

```python
@hook_impl
@register_attr("post_models")
def post_model(markata):
    \"\"\"Add fields to post model.\"\"\"
    from pydantic import BaseModel, Field

    class MyPostFields(BaseModel):
        custom_date: str = Field(None, description="Custom date field")
        tags: list[str] = Field(default_factory=list)

    return MyPostFields
```

## Content Processing Hooks

Handle content transformation:

```python
@hook_impl(trylast=True)  # Run after other render hooks
def render(markata):
    \"\"\"Process each article.\"\"\"
    for article in markata.filter("not skip"):
        # Add custom processing
        if article.tags:
            article.tag_links = [f"<a href='/tags/{tag}'>{tag}</a>"
                               for tag in article.tags]
```

## Output Generation Hooks

Control how content is saved:

```python
@hook_impl
def save(markata):
    \"\"\"Save processed content.\"\"\"
    output_dir = Path(markata.config.output_dir)

    # Save custom index
    if markata.config.my_plugin.enabled:
        index = generate_custom_index(markata.articles)
        (output_dir / "custom.html").write_text(index)
```

# Hook Ordering

Control execution order with decorators:

```python
# Run first in configure stage
@hook_impl(tryfirst=True)
def configure(markata): ...

# Run in middle (default)
@hook_impl
def configure(markata): ...

# Run last in configure stage
@hook_impl(trylast=True)
def configure(markata): ...
```

# Attribute Registration

Register data on the Markata instance:

```python
# Single attribute
@register_attr("articles")
def my_hook(markata):
    markata.articles = []

# Multiple attributes
@register_attr("articles", "tags", "categories")
def my_hook(markata):
    markata.articles = []
    markata.tags = {}
    markata.categories = {}

# Access in other hooks
@hook_impl
def render(markata):
    print(markata.articles)  # Access registered data
```

# Complex Example

Here's a complete plugin example combining multiple hooks:

```python
from pathlib import Path
from typing import List, Optional
from pydantic import BaseModel, Field
from markata.hookspec import hook_impl, register_attr

class TagConfig(BaseModel):
    \"\"\"Configuration for tag handling.\"\"\"
    enabled: bool = True
    min_posts: int = 2
    output_dir: str = "tags"

class TaggedPost(BaseModel):
    \"\"\"Add tag fields to posts.\"\"\"
    tags: List[str] = Field(default_factory=list)
    tag_links: Optional[str] = None

@hook_impl
@register_attr("config_models")
def config_model(markata):
    \"\"\"Add tag configuration.\"\"\"
    markata.config_models.append(TagConfig)

@hook_impl
@register_attr("post_models")
def post_model(markata):
    \"\"\"Add tag fields to posts.\"\"\"
    markata.post_models.append(TaggedPost)


@hook_impl(trylast=True)
def render(markata):
    \"\"\"Add tag links to articles.\"\"\"
    if not markata.config.tags.enabled:
        return

    for article in markata.filter("not skip"):
        article.tag_links = " ".join(
            f"<a href='/tags/{tag}'>{tag}</a>"
            for tag in article.tags
        )

@hook_impl
def save(markata):
    \"\"\"Generate tag pages.\"\"\"
    if not markata.config.tags.enabled:
        return

    output_dir = Path(markata.config.output_dir)
    tag_dir = output_dir / markata.config.tags.output_dir
    tag_dir.mkdir(exist_ok=True)

    for tag, articles in markata.tags.items():
        if len(articles) >= markata.config.tags.min_posts:
            content = generate_tag_page(tag, articles)
            (tag_dir / f"{tag}.html").write_text(content)
```

This example shows:
1. Configuration definition
2. Model extension
3. Data processing
4. Content generation
5. Output handling

See [[ markata/lifecycle ]] for the exact order hooks are executed.

"""

import functools
from typing import Any, Callable, Dict, List, TYPE_CHECKING

import pluggy

from markata.lifecycle import LifeCycle

if TYPE_CHECKING:
    import typer

    from markata import Markata

HOOK_NAMESPACE = "markata"
hook_spec = pluggy.HookspecMarker(HOOK_NAMESPACE)
hook_impl = pluggy.HookimplMarker(HOOK_NAMESPACE)


class MarkataSpecs:
    """
    Namespace that defines all specifications for Load hooks.

    configure -> glob -> load -> render -> save
    """


@hook_spec
def generic_lifecycle_method(
    markata: "Markata",
) -> Any: ...


@hook_spec
def cli_lifecycle_method(markata: "Markata", app: "typer.Typer") -> Any:
    "A Markata lifecycle methos that includes a typer app used for cli's"


for method in LifeCycle._member_map_:
    if "cli" in method:
        setattr(MarkataSpecs, method, cli_lifecycle_method)
    else:
        setattr(MarkataSpecs, method, generic_lifecycle_method)

registered_attrs: Dict[str, List[Dict[str, Any]]] = {}


def register_attr(*attrs: Any) -> Callable:
    def decorator_register(
        func: Callable,
    ) -> Callable:
        for attr in attrs:
            if attr not in registered_attrs:
                registered_attrs[attr] = []
            registered_attrs[attr].append(
                {
                    "func": func,
                    "funcname": func.__code__.co_name,
                    "lifecycle": getattr(LifeCycle, func.__code__.co_name),
                },
            )

        @functools.wraps(func)
        def wrapper_register(markata: "Markata", *args: Any, **kwargs: Any) -> Any:
            return func(markata, *args, **kwargs)

        return wrapper_register

    return decorator_register
