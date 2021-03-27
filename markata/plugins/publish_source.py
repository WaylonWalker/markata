from pathlib import Path
from typing import TYPE_CHECKING

import frontmatter

from markata.hookspec import hook_impl

if TYPE_CHECKING:
    from markata import Markata


@hook_impl
def save(markata: "Markata") -> None:
    output_dir = Path(str(markata.config["output_dir"]))
    output_dir.mkdir(parents=True, exist_ok=True)
    for article in markata.iter_articles(description="saving source documents"):
        with open(output_dir / Path(article["path"]).name, "w+") as f:
            f.write(frontmatter.dumps(article))
