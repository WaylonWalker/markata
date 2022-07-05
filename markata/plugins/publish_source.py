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

!! note
    publish_source is included by default, but if you have not included the
    default set of hooks you will need to explicitly add it.
"""
import hashlib
from pathlib import Path
from typing import TYPE_CHECKING

import frontmatter
import yaml
from yaml.representer import RepresenterError

from markata import background
from markata.hookspec import hook_impl

if TYPE_CHECKING:
    from markata import Markata


def _strip_unserializable_values(article: frontmatter.Post) -> frontmatter.Post:
    """
    Returns an article with only yaml serializable frontmatter.
    """
    _article = frontmatter.Post(
        article.content, **{k: v for k, v in article.metadata.items() if k != "content"}
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
    return _article


def _is_relative_to(output_dir: Path, output_html: Path) -> bool:
    try:
        output_html.relative_to(output_dir)
        return True
    except ValueError:
        return False


@background.task
def _save(markata: "Markata", article: "Post"):
    article_path = (
        output_dir / Path(article["slug"]).parent / Path(article["path"]).name
    )

    try:
        content = frontmatter.dumps(article)
    except RepresenterError:
        article = _strip_unserializable_values(article)
        content = frontmatter.dumps(article)

    if (
        article_path.exists()
        and hashlib.sha256(content.encode("utf-8")).hexdigest()
        == hashlib.sha256(article_path.read_bytes()).hexdigest()
    ):
        ...
    elif _is_relative_to(output_dir, article_path):
        try:
            Path(article_path).write_text(content)
        except RepresenterError:
            _article = _strip_unserializable_values(article)

            _save(output_dir, _article)


@hook_impl
def save(markata: "Markata") -> None:
    """
    Saves the final modified post to the output site as markdown.

    !! note
      Any keys that are not yaml serializable will be stripped.
    """
    output_dir = Path(str(markata.config["output_dir"]))
    output_dir.mkdir(parents=True, exist_ok=True)
    futures = [_save(markata, article) for article in markata.articles]
    while not all([f.done() for f in futures]):
        ...
