"""
Sets the articles `output_html` path, and saves the article's `html` to the
`output_html` file.

##  Ouptut Directory

Output will always be written inside of the configured `output_dir`

```toml
[markata]
# markout is the default, but you can override it in your markata.toml file
output_dir = "markout"
```

## Explicityly set the output

markata will save the articles `html` to the `output_html` specified in the
articles metadata, loaded from frontmatter.

## 404 example use case

Here is an example use case of explicitly setting the output_html.  By default
markata will turn `pages/404.md` into `markout/404/index.html`, but many
hosting providers look for a 404.html to redirect the user to when a page is
not found.

```markdown
---
title: Whoops that page was not found
description: 404, looks like we can't find the page you are looking for
output_html: 404.html

---

404, looks like we can't find the page you are looking for.  Try one of these
pages.

<ul>
{% for post in
    markata.map(
        'post',
        filter='"markata" not in slug and "tests" not in slug and "404" not in slug'
        )
 %}
    <li><a href="{{ post.slug }}">{{ post.title or "CHANGELOG" }}</a></li>
{% endfor %}
</ul>
```

## Index.md is the one special case

If you have a file `pages/index.md` it will become `markout/index.html` rather
than `markout/index/inject.html` This is one of the primary ways that markata
lets you [make your home page](https://markata.dev/home-page/)

"""

from pathlib import Path
from typing import Any, Dict, Optional, TYPE_CHECKING

import pydantic
from pydantic import Field, field_validator, ConfigDict

from markata.hookspec import hook_impl, register_attr

if TYPE_CHECKING:
    from markata import Markata


class OutputHTML(pydantic.BaseModel):
    markata: Any = Field(None, exclude=True)
    path: Path
    slug: str = None
    output_html: Optional[str] = None

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
    def default_slug(cls, v, info):
        from slugify import slugify

        if v is None:
            return slugify(str(info.data.get("path", "").stem))
        return v

    @field_validator("output_html", mode="before")
    def default_output_html(cls: "OutputHTML", v: Optional[str], info) -> Path:
        if isinstance(v, str):
            v = Path(v)
        if v is not None:
            return v

        slug = info.data.get("slug")
        if slug is None:
            slug = cls.default_slug(None, info)

        if slug == "index":
            return cls.markata.config.output_dir / "index.html"
        return cls.markata.config.output_dir / slug / "index.html"

    @field_validator("output_html", mode="before")
    def output_html_relative(cls: "OutputHTML", v: Optional[Path], info) -> Path:
        if v is None:
            return cls.default_output_html(v, info)
        return v

    @field_validator("output_html", mode="before")
    def output_html_exists(cls: "OutputHTML", v: Optional[Path], info) -> Path:
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

    for article in markata.filter("skip==False"):
        if article.html is None:
            continue
        if isinstance(article.html, str):
            # Create parent directories before writing
            article.output_html.parent.mkdir(parents=True, exist_ok=True)
            article.output_html.write_text(article.html)
        if isinstance(article.html, Dict):
            for slug, html in article.html.items():
                if slug == "index":
                    slug = ""
                    output_html = article.output_html
                elif "." in slug:
                    output_html = article.output_html.parent / slug
                else:
                    slug = slugify(slug)
                    output_html = article.output_html.parent / slug / "index.html"

                output_html.parent.mkdir(parents=True, exist_ok=True)
                output_html.write_text(html)
