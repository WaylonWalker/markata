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
from typing import TYPE_CHECKING, Any, Dict, Optional

import pydantic
from slugify import slugify

from markata.hookspec import hook_impl, register_attr

if TYPE_CHECKING:
    from markata import Markata


class OutputHTML(pydantic.BaseModel):
    markata: Any = None
    path: Path
    slug: str = None
    output_html: Path = None

    @pydantic.validator("slug", pre=True, always=True)
    @classmethod
    def default_slug(cls, v, *, values):
        if v is None:
            return slugify(str(values["path"].stem))
        return v

    @pydantic.validator("output_html", pre=True, always=True)
    def default_output_html(
        cls: "OutputHTML", v: Optional[Path], *, values: Dict
    ) -> Path:
        if isinstance(v, str):
            v = Path(v)
        if v is not None:
            return v
        if "slug" not in values:
            for validator in cls.__validators__["slug"]:
                values["slug"] = validator.func(cls, v, values=values)

        if values["slug"] == "index":
            return cls.markata.config.output_dir / "index.html"
        return cls.markata.config.output_dir / values["slug"] / "index.html"

    @pydantic.validator("output_html")
    def output_html_relative(
        cls: "OutputHTML", v: Optional[Path], *, values: Dict
    ) -> Path:
        if isinstance(v, str):
            v = Path(v)
        if cls.markata.config.output_dir.absolute() not in v.absolute().parents:
            return cls.markata.config.output_dir / v
        return v

    @pydantic.validator("output_html")
    def output_html_exists(
        cls: "OutputHTML", v: Optional[Path], *, values: Dict
    ) -> Path:
        if isinstance(v, str):
            v = Path(v)
        if not v.parent.exists():
            v.parent.mkdir(parents=True, exist_ok=True)
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

    for article in markata.articles:
        article.output_html.write_text(article.html)
