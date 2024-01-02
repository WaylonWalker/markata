"""
Writes the final modified markdown and frontmatter to the output directory.
Replacing the trailing slash if its there and adding .md will bring up the raw
source.

## Configuration

The only configuration for the publish_source plugin is to make sure its in
your list of hooks.


``` toml
[markata]

# make sure its in your list of hooks
hooks=[
   "markata.plugins.publish_source",
   ]
```

!!! note
    publish_source is included by default, but if you have not included the
    default set of hooks you will need to explicitly add it.
"""
from pathlib import Path
from typing import TYPE_CHECKING

import frontmatter
import yaml
from yaml.representer import RepresenterError

from markata.hookspec import hook_impl

if TYPE_CHECKING:
    from markata import Markata


def _save(output_dir: Path, article: frontmatter.Post) -> None:
    """
    saves the article to the output directory at its specified slug.
    """
    path = Path(
        output_dir / Path(article["slug"]).parent / Path(article["path"]).name,
    )
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(article.dumps())


def _strip_unserializable_values(
    markata: "Markata",
    article: frontmatter.Post,
) -> frontmatter.Post:
    """
    Returns an article with only yaml serializable frontmatter.
    """
    _article = frontmatter.Post(
        article.content,
        **{k: v for k, v in article.metadata.items() if k != "content"},
    )
    kwargs = {
        "Dumper": yaml.cyaml.CSafeDumper,
        "default_flow_style": False,
        "allow_unicode": True,
    }
    for key, value in article.metadata.items():
        try:
            yaml.dump({key: value}, **kwargs)
        except RepresenterError:
            del _article[key]
    if markata.Post:
        _article = markata.Post(**_article.metadata, path=str(article.path))
    return _article


@hook_impl
def save(markata: "Markata") -> None:
    """
    Saves the final modified post to the output site as markdown.

    !!! note
        Any keys that are not yaml serializable will be stripped.

    """
    output_dir = Path(str(markata.config["output_dir"]))
    for (
        article
    ) in markata.articles:  # iter_articles(description="saving source documents"):
        try:
            _save(output_dir, article)
        except RepresenterError:
            _article = _strip_unserializable_values(markata, article)

            _save(output_dir, _article)
