"""
The `markata.plugins.publish_dev_to_source` plugin enables synchronization between
your Markata posts and dev.to articles. It handles frontmatter translation and content
formatting for dev.to compatibility.

## Installation

This plugin is built-in but not enabled by default. Add it to your plugins list:

```toml
hooks = [
    "markata.plugins.publish_dev_to_source",
]
```

## Uninstallation

Remove the plugin from your hooks list in `markata.toml`:

```toml
hooks = [
    # Remove or comment out the line below
    # "markata.plugins.publish_dev_to_source",
]
```

## Configuration

No explicit configuration is required. The plugin automatically processes posts with
dev.to-specific frontmatter.

## Functionality

## Frontmatter Mapping

Automatically maps these frontmatter fields:
- title
- published
- description
- tags
- canonical_url
- cover_image
- series

## Content Processing

The plugin:
1. Joins consecutive text lines for dev.to formatting
2. Preserves code blocks and lists
3. Maintains markdown compatibility
4. Updates canonical URLs

## Post Model

Extends the base Post model with:
- dev.to-specific fields
- Validation rules
- Performance optimizations

## Publishing Flow

1. Processes posts after rendering
2. Updates frontmatter
3. Formats content
4. Saves changes to source files

## Dependencies

This plugin depends on:
- python-frontmatter for YAML handling
- pydantic for model validation
"""

from pathlib import Path
from typing import TYPE_CHECKING
from typing import Any
from typing import Optional

import frontmatter
import pydantic
from pydantic import ConfigDict
from pydantic import Field

from markata.hookspec import hook_impl
from markata.hookspec import register_attr

if TYPE_CHECKING:
    from markata import Markata

DEV_TO_FRONTMATTER = [
    "title",
    "published",
    "description",
    "tags",
    "canonical_url",
    "cover_image",
    "series",
]


def should_join(line):
    if line == "":
        return False
    if line is None:
        return False
    if line[0].isalpha():
        return True
    if line[0].startswith("["):
        return True
    if line[0].startswith("!"):
        return True
    return False


def join_lines(article):
    lines = article.split("\n")
    line_number = 0
    while line_number + 1 < len(lines):
        line = lines[line_number]
        nextline = lines[line_number + 1]
        if should_join(line) and should_join(nextline):
            lines[line_number] = f"{line} {nextline}"
            lines.pop(line_number + 1)
        else:
            line_number += 1

    return "\n".join(lines)


class PublishDevToSourcePost(pydantic.BaseModel):
    markata: Any = Field(None, exclude=True)
    model_config = ConfigDict(
        validate_assignment=False,  # Performance model
        arbitrary_types_allowed=True,
        extra="allow",
        str_strip_whitespace=True,
        validate_default=True,
        coerce_numbers_to_str=True,
        populate_by_name=True,
    )
    canonical_url: Optional[str] = None


@hook_impl
@register_attr("post_models")
def post_model(markata: "Markata") -> None:
    markata.post_models.append(PublishDevToSourcePost)


@hook_impl
def post_render(markata: "Markata") -> None:
    for post in markata.iter_articles(description="saving source documents"):
        article = frontmatter.Post(
            post.content,
            **{k: v for k, v in post.metadata.items() if k in DEV_TO_FRONTMATTER},
        )

        article.content = join_lines(article.content)
        article.content = join_lines(article.content)

        if "canonical_url" not in article:
            article["canonical_url"] = f"{markata.config.url}/{post.slug}/"

        if "published" not in article:
            article["published"] = True

        if "cover_image" not in article:
            article["cover_image"] = f"{markata.config.images_url}/{post.slug}.png"
        post.dev_to = article


@hook_impl
def save(markata: "Markata") -> None:
    output_dir = Path(str(markata.config.output_dir))
    with markata.console.status("Saving source documents...") as status:
        for post in markata.iter_articles(description="saving source documents"):
            status.update(f"Saving {post['slug']}...")
            with open(output_dir / Path(post["slug"]) / "dev.md", "w+") as f:
                f.write(frontmatter.dumps(post.dev_to))
