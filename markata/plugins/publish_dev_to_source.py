from pathlib import Path
from typing import TYPE_CHECKING

import frontmatter

from markata.hookspec import hook_impl

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
    joined = []
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


@hook_impl
def save(markata: "Markata") -> None:
    output_dir = Path(str(markata.config["output_dir"]))
    output_dir.mkdir(parents=True, exist_ok=True)

    for real_article in markata.iter_articles(description="saving source documents"):
        from copy import copy
        from copy import deepcopy

        article = deepcopy(real_article)

        before_keys = copy(list(article.keys()))
        for key in before_keys:
            if key not in DEV_TO_FRONTMATTER:
                del article[key]

        article.content = join_lines(article.content)
        article.content = join_lines(article.content)

        if "canonical_url" not in article:
            article["canonical_url"] = f'{markata.url}/{real_article["slug"]}/'

        if "published" not in article:
            article["published"] = True

        if "cover_image" not in article:
            article[
                "cover_image"
            ] = f"{markata.config['images_url']}/{real_article['slug']}.png"

        with open(output_dir / Path(real_article["slug"]) / "dev.md", "w+") as f:
            f.write(frontmatter.dumps(article))
