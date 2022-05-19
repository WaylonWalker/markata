"""
Writes the final modified markdown and frontmatter to the output directory.
"""
from copy import deepcopy
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
    with open(
        output_dir / Path(article["slug"]).parent / Path(article["path"]).name, "w+"
    ) as f:
        f.write(frontmatter.dumps(article))


def _strip_unserializable_values(article: frontmatter.Post) -> frontmatter.Post:
    """
    Returns an article with only yaml serializable frontmatter.
    """
    _article = deepcopy(article)
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


@hook_impl
def save(markata: "Markata") -> None:
    """
    Saves the final modified post to the output site as markdown.

    !! note
      Any keys that are not yaml serializable will be stripped.
    """
    output_dir = Path(str(markata.config["output_dir"]))
    output_dir.mkdir(parents=True, exist_ok=True)
    for article in markata.iter_articles(description="saving source documents"):
        try:
            _save(output_dir, article)
        except RepresenterError:
            _article = _strip_unserializable_values(article)

            _save(output_dir, _article)
