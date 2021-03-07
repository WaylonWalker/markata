from pathlib import Path

import frontmatter
from more_itertools import flatten
from tqdm import tqdm

from markata.hookspec import hook_impl


@hook_impl
def save(markata):
    output_dir = Path(markata.config["output_dir"])
    output_dir.mkdir(parents=True, exist_ok=True)
    for article in tqdm(
        markata.articles, desc="saving source document", colour="yellow", leave=False
    ):
        with open(output_dir / Path(article["path"]).name, "w+") as f:
            f.write(frontmatter.dumps(article))
