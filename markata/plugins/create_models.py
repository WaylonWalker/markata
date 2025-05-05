"""
The `markata.plugins.create_models` plugin dynamically creates Markata's core Pydantic
models by combining model fragments from various plugins. This enables extensible and
type-safe data models.

## Installation

This plugin is built-in and enabled by default through the 'default' plugin.
If you want to be explicit, you can add it to your list of plugins:

```toml
hooks = [
    "markata.plugins.create_models",
]
```

## Uninstallation

Since this plugin is included in the default plugin set, to disable it you must explicitly
add it to the disabled_hooks list if you are using the 'default' plugin:

```toml
disabled_hooks = [
    "markata.plugins.create_models",
]
```

Note: Disabling this plugin will break most of Markata's functionality as it's
responsible for creating core models.

## Configuration

No explicit configuration is required. The plugin automatically creates models from
registered model fragments.

## Functionality

## Model Creation

The plugin:
1. Collects model fragments from plugins
2. Deduplicates model definitions
3. Creates composite models
4. Registers models with Markata

## Core Models

Creates these models:
- `Post`: Individual post data
- `Posts`: Collection of posts
- `Config`: Global configuration

## Model Configuration

Applies settings:
- Environment variable prefix
- Extra field handling
- Type validation
- Copy behavior
- Base classes

## Performance

Optimizations:
- Unique model deduplication
- Lazy model creation
- Minimal validation
- Efficient inheritance

## Dependencies

This plugin depends on:
- pydantic for model creation
- more-itertools for deduplication
"""

from typing import TYPE_CHECKING
from typing import List

from more_itertools import unique_everseen
from pydantic import create_model

from markata.hookspec import hook_impl
from markata.hookspec import register_attr

if TYPE_CHECKING:
    from markata import Markata


class Config:
    env_prefix = "markata_"
    extras = "allow"


class PostConfig:
    title = "Markata.Post"
    arbitrary_types_allowed = True
    copy_on_model_validation = False


@hook_impl
@register_attr("Post", "Posts", "Config")
def create_models(markata: "Markata") -> None:
    post_models = tuple(unique_everseen(markata.post_models))
    markata.Post = create_model(
        "Post",
        __base__=post_models,
    )
    markata.Posts = create_model(
        "Posts",
        posts=(List[markata.Post], ...),
    )
    markata.Post.markata = markata
    markata.Config = create_model(
        "Config",
        __base__=tuple(unique_everseen(markata.config_models)),
    )
