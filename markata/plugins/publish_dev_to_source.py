from pathlib import Path
from typing import Optional, TYPE_CHECKING

import frontmatter
import pydantic

from markata.hookspec import hook_impl, register_attr

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
    markata: Markata
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
    for post in markata.iter_articles(description="saving source documents"):
        with open(output_dir / Path(post["slug"]) / "dev.md", "w+") as f:
            f.write(frontmatter.dumps(post.dev_to))
