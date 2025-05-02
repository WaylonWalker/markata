"""
The `markata.plugins.flat_slug` plugin handles URL slug generation for your articles.
It sanitizes special characters that don't work in browsers and provides a consistent
URL structure.

## Installation

This plugin is built-in and enabled by default through the 'default' plugin.
If you want to be explicit, you can add it to your list of plugins:

```toml
hooks = [
    "markata.plugins.flat_slug",
]
```

## Uninstallation

Since this plugin is included in the default plugin set, to disable it you must explicitly
add it to the disabled_hooks list if you are using the 'default' plugin:

```toml
disabled_hooks = [
    "markata.plugins.flat_slug",
]
```

Note: Disabling this plugin will prevent automatic URL slug generation and sanitization.

## Configuration

Configure the plugin behavior in your `markata.toml`:

```toml
[markata.flat_slug]
slugify = true  # Set to false to disable automatic slug sanitization
```

## Functionality

## Explicit Slug in Frontmatter

If you explicitly set the slug in the frontmatter of a post, markata will not
overwrite it:

```markdown
---
title: My First Post
slug: /my-post
---

This is my first post it will be at `<markata.config.url>/my-post/`
regardless of filename.
```

## Automatic Slug Based on Filename

By default the flat_slug plugin will use the `stem` of your filename (filename without
extension) unless you explicitly set your slug in frontmatter.

Examples:
* `/pages/my-post.md` becomes `<markata.config.url>/my-post/`
* `/pages/blog/a-blog-post.md` becomes `<markata.config.url>/a-blog-post/`

## Registered Attributes

The plugin registers the following attributes on Post objects:
- `should_slugify`: Boolean indicating if the post's slug should be sanitized
"""

from pathlib import Path
from typing import Any
from typing import Optional

import pydantic

from markata import Markata
from markata.hookspec import hook_impl
from markata.hookspec import register_attr


class FlatSlugConfig(pydantic.BaseModel):
    slugify: bool = True


class Config(pydantic.BaseModel):
    flat_slug: FlatSlugConfig = FlatSlugConfig()


class FlatSlugPost(pydantic.BaseModel):
    should_slugify: Optional[bool] = None
    markata: Any = pydantic.Field(None, exclude=True)

    model_config = pydantic.ConfigDict(
        validate_assignment=False,
        arbitrary_types_allowed=True,
        extra="allow",
        str_strip_whitespace=True,
        validate_default=True,
        coerce_numbers_to_str=True,
        populate_by_name=True,
    )

    @pydantic.field_validator("should_slugify", mode="before")
    @classmethod
    def default_slugify(cls, v: Optional[bool], info) -> bool:
        if v is None:
            markata = info.data.get("markata")
            if markata is None:
                return True  # Default to True if no markata instance
            return markata.config.flat_slug.slugify
        return v


@hook_impl()
@register_attr("config_models")
def config_model(markata: Markata) -> None:
    markata.config_models.append(Config)


@hook_impl
@register_attr("post_models")
def post_model(markata: "Markata") -> None:
    markata.post_models.append(FlatSlugPost)


@hook_impl(tryfirst=True)
def pre_render(markata: "Markata") -> None:
    """
    Sets the article slug if one is not already set in the frontmatter.
    """
    from slugify import slugify

    for article in markata.iter_articles(description="creating slugs"):
        stem = article.get(
            "slug",
            Path(article.get("path", article.get("title", ""))).stem,
        )
        if article.should_slugify:
            article.slug = "/".join([slugify(s) for s in stem.split("/")])
        else:
            article.slug = stem
