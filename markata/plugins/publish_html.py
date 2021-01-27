from markata.hookspec import hook_impl
from pathlib import Path
from more_itertools import flatten

from tqdm import tqdm


@hook_impl
def save(markata):
    output_dir = Path(markata.config["output_dir"])
    output_dir.mkdir(parents=True, exist_ok=True)
    for article in tqdm(
        markata.articles, desc="saving html", colour="yellow", leave=False
    ):
        article_path = output_dir / article["slug"] / "index.html"
        article_path.parent.mkdir(parents=True, exist_ok=True)
        with open(article_path, "w+") as f:
            f.write(article.html)
