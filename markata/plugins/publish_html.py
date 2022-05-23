"""
Sets the articles `output_html` path, and saves the article's `html` to the
`output_html` file.

## Ouptut Directory

Output will always be written inside of the configured `output_dir`

```toml
[markata]
# markout is the default, but you can override it in your markata.toml file
output_dir = "markout"
```

## Explicityly set the output

markata will save the articles `html` to the `output_html` specified in the
articles metadata, loaded from frontmatter.

### 404 example use case

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
{% for post in markata.map('post', filter='"markata" not in slug and "tests" not in slug and "404" not in slug') %}
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
from typing import TYPE_CHECKING

from markata.hookspec import hook_impl

if TYPE_CHECKING:
    from markata import Markata


def _is_relative_to(output_dir: Path, output_html: Path):
    try:
        output_html.relative_to(output_dir)
        return True
    except ValueError:
        return False


@hook_impl
def pre_render(markata: "Markata") -> None:
    """
    Sets the `output_html` in the articles metadata.  If the output is
    explicitly given, it will make sure its in the `output_dir`, if it is not
    explicitly set it will use the articles slug.
    """
    output_dir = Path(markata.config["output_dir"])  # type: ignore
    output_dir.mkdir(parents=True, exist_ok=True)

    for article in markata.articles:
        if "output_html" in article.metadata:
            article_path = Path(article["output_html"])
            if not _is_relative_to(output_dir, article_path):
                article["output_html"] = output_dir / article["output_html"]
        elif article["slug"] == "index":
            article["output_html"] = output_dir / "index.html"
        else:
            article["output_html"] = output_dir / article["slug"] / "index.html"


@hook_impl
def save(markata: "Markata") -> None:
    """
    Saves all the articles to their set `output_html` location if that location
    is relative to the specified `output_dir`.  If its not relative to the
    `output_dir` it will log an error and move on.
    """
    output_dir = Path(markata.config["output_dir"])  # type: ignore

    for article in markata.articles:
        article_path = Path(article["output_html"])
        if _is_relative_to(output_dir, article_path):
            article_path.parent.mkdir(parents=True, exist_ok=True)
            with open(article_path, "w+") as f:
                f.write(article.html)
        else:
            markata.console.log(
                f'article "{article["path"]}" attempted to write to "{article["output_html"]}" outside of the configured output_dir "{output_dir}"'
            )
