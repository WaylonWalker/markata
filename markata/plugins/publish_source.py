from markata.hookspec import hook_impl
from pathlib import Path
from more_itertools import flatten

import frontmatter
from tqdm import tqdm


@hook_impl
def save(markata):
    output_dir = Path(markata.config["output_dir"])
    output_dir.mkdir(parents=True, exist_ok=True)
    for article in tqdm(
        markata.articles, desc="saving source document", colour="yellow", leave=False
    ):
        with open(output_dir / Path(article["path"]).name, "w+") as f:
            f.write(frontmatter.dumps(article))
