"""
The `markata.plugins.publish_html` plugin handles saving rendered HTML content to files.
It determines the output path for each article and ensures files are saved in the correct
location within the output directory.

## Installation

This plugin is built-in and enabled by default through the 'default' plugin.
If you want to be explicit, you can add it to your list of plugins:

```toml
hooks = [
    "markata.plugins.publish_html",
]
```

## Uninstallation

Since this plugin is included in the default plugin set, to disable it you must explicitly
add it to the disabled_hooks list if you are using the 'default' plugin:

```toml
disabled_hooks = [
    "markata.plugins.publish_html",
]
```

Note: Disabling this plugin will prevent HTML files from being written to disk.

## Configuration

Configure HTML output in `markata.toml`:

```toml
[markata]
# Base output directory
output_dir = "dist"

# Custom output paths
[[markata.output_paths]]
pattern = "blog/*"
output = "posts/{stem}.html"

[[markata.output_paths]]
pattern = "docs/*"
output = "documentation/{stem}/index.html"
```

## Functionality

## Path Resolution

The plugin:
1. Determines output path for each post
2. Creates necessary directories
3. Validates paths are within output_dir
4. Handles custom path mappings

## Output Model

Extends the base Post model with:
- output_html path
- Path validation
- Slug resolution
- Directory creation

## File Operations

Handles:
- Directory creation
- File writing
- Path validation
- Error logging

## Path Customization

Supports:
- Custom output paths
- Path patterns
- Directory structures
- Index files

## Safety Features

Includes:
- Path validation
- Directory verification
- Error handling
- Logging

## Dependencies

This plugin depends on:
- pathlib for path operations
- pydantic for model validation
"""

from pathlib import Path
from typing import TYPE_CHECKING
from typing import Any
from typing import Optional
from typing import Union

import pydantic
from pydantic import ConfigDict
from pydantic import Field
from pydantic import field_validator

from markata.hookspec import hook_impl
from markata.hookspec import register_attr

if TYPE_CHECKING:
    from markata import Markata


class OutputHTML(pydantic.BaseModel):
    markata: Any = Field(None, exclude=True)
    path: Path
    slug: str = None
    output_html: Optional[Path] = None

    class Config:
        model_config = ConfigDict(
            validate_assignment=False,
            arbitrary_types_allowed=True,
            extra="allow",
            str_strip_whitespace=True,
            validate_default=True,
            coerce_numbers_to_str=True,
            populate_by_name=True,
        )

    @field_validator("slug", mode="before")
    @classmethod
    def default_slug(cls, v, info) -> str:
        from slugify import slugify

        if v is None:
            path = info.data.get("path")
            if path is None:
                return ""
            return slugify(str(path.stem))
        return v

    @field_validator("output_html", mode="before")
    @classmethod
    def default_output_html(cls, v: Optional[Union[str, Path]], info) -> Optional[Path]:
        if isinstance(v, str):
            v = Path(v)
        if v is not None:
            return v

        markata = info.data.get("markata")
        if markata is None:
            raise ValueError("markata is required")

        slug = info.data.get("slug")
        if slug is None:
            slug = cls.default_slug(None, info)

        if slug in ["", "index",'/']:
            return markata.config.output_dir / "index.html"
        return markata.config.output_dir / slug / "index.html"

    @field_validator("output_html", mode="before")
    @classmethod
    def output_html_relative(cls, v: Optional[Path], info) -> Optional[Path]:
        if v is None:
            return cls.default_output_html(v, info)
        return v

    @field_validator("output_html", mode="before")
    @classmethod
    def output_html_exists(cls, v: Optional[Path], info) -> Optional[Path]:
        if v is None:
            return cls.default_output_html(v, info)
        return v


@hook_impl
@register_attr("post_models")
def post_model(markata: "Markata") -> None:
    markata.post_models.append(OutputHTML)


@hook_impl
def save(markata: "Markata") -> None:
    """
    Saves all the articles to their set `output_html` location if that location
    is relative to the specified `output_dir`.  If its not relative to the
    `output_dir` it will log an error and move on.
    """
    from slugify import slugify

    for article in markata.filter("not skip"):
        if article.html is None:
            continue

        if isinstance(article.html, str):
            # Create parent directories before writing
            current_html = article.output_html.read_text() if article.output_html.exists() else ""
            if current_html != article.html:
                article.output_html.parent.mkdir(parents=True, exist_ok=True)
                article.output_html.write_text(article.html)
        elif isinstance(article.html, dict):
            for slug, html in article.html.items():
                # Handle special case for index
                if slug == "index":
                    output_path = article.output_html
                # Handle files with extensions
                elif "." in slug:
                    output_path = article.output_html.parent / slug
                # Handle other slugs by creating subdirectories
                else:
                    slug_path = slugify(slug)
                    output_path = article.output_html.parent / slug_path / "index.html"

                # Create parent directories and write file
                output_path.parent.mkdir(parents=True, exist_ok=True)
                current_html = output_path.read_text() if output_path.exists() else ""
                if current_html != html:
                    output_path.write_text(html)
