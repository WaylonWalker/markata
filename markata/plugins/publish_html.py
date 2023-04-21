"""
Sets the articles `output_html` path, and saves the article's `html` to the
`output_html` file.

# Ouptut Directory

Output will always be written inside of the configured `output_dir`

```toml
[markata]
# markout is the default, but you can override it in your markata.toml file
output_dir = "markout"
```

# Explicityly set the output

markata will save the articles `html` to the `output_html` specified in the
articles metadata, loaded from frontmatter.

# 404 example use case

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

# Index.md is the one special case

If you have a file `pages/index.md` it will become `markout/index.html` rather
than `markout/index/inject.html` This is one of the primary ways that markata
lets you [make your home page](https://markata.dev/home-page/)

"""
from pathlib import Path
from typing import TYPE_CHECKING, Dict, Union

import pydantic

# if TYPE_CHECKING:
from markata import Markata
from markata.hookspec import hook_impl, register_attr


def _is_relative_to(output_dir: Path, output_html: Path):
    try:
        output_html.relative_to(output_dir)
        return True
    except ValueError:
        return False


# @hook_impl
# def pre_render(markata: "Markata") -> None:
#     """
#     Sets the `output_html` in the articles metadata.  If the output is
#     explicitly given, it will make sure its in the `output_dir`, if it is not
#     explicitly set it will use the articles slug.
#     """
#     output_dir = Path(markata.config["output_dir"])  # type: ignore
#     # output_dir.mkdir(parents=True, exist_ok=True)

#     for article in markata.articles:
#         if article.output_html:
#             article_path = Path(article["output_html"])
#             if not _is_relative_to(output_dir, article_path):
#                 article["output_html"] = output_dir / article["output_html"]
#         elif article["slug"] == "index":
#             article["output_html"] = output_dir / "index.html"
#         else:
#             article["output_html"] = output_dir / article["slug"] / "index.html"


class OutputHTML(pydantic.BaseModel):
    markata: Markata
    slug: str
    output_html: Union[str, Path] = None

    class Config:
        validate_assignment = True
        arbitrary_types_allowed = True

    @pydantic.validator("output_html", pre=True, always=True)
    @classmethod
    def output_html_path(cls, v, *, values: Dict) -> Path:
        if v:
            v = Path(v)
        return values["markata"].config.output_dir / values["slug"] / "index.html"

    @pydantic.validator("output_html")
    @classmethod
    def output_html_relative(cls, v, *, values: Dict) -> Path:
        if not v.relative_to(values["markata"].config.output_dir):
            return values["markata"].config.output_dir / v
        return v

    @pydantic.validator("output_html")
    @classmethod
    def output_html_exists(cls, v, *, values: Dict) -> Path:
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
