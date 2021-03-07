"""Default load plugin."""
import time
from pathlib import Path

import frontmatter
from more_itertools import flatten
from tqdm import tqdm
from yaml.parser import ParserError

from markata.background import task
from markata.hookspec import hook_impl


@hook_impl
def load(markata):
    # print("loading articles")

    futures = [get_post(article, markata) for article in markata.files]
    with tqdm(
        total=len(futures), desc="loading markdown", leave=False, colour="yellow"
    ) as pbar:
        while not all([f.done() for f in futures]):
            time.sleep(0.1)
            for _ in range(len([f for f in futures if f.done()]) - pbar.n):
                pbar.update()
    articles = [f.result() for f in futures]
    articles = [a for a in articles if a]
    # print(f"loaded {len(articles)}")

    # articles.sort(key=lambda post: int(post["date"].strftime("%Y%m%d")))
    # sorted(articles, key=lambda x: x["date"])
    markata.articles = articles


@task
def get_post(path: Path, markata):
    default = {
        "cover": "",
        "title": "",
        "tags": [],
        "status": "draft",
        "templateKey": "",
        "path": str(path),
        "description": "",
        "content": "",
    }
    try:
        post = frontmatter.load(path)
        post.metadata = {**default, **post.metadata}
    except ParserError:
        return
        post = default
    except ValueError:
        return
        post = default
    post.metadata["path"] = str(path)
    post["content_hash"] = markata.make_hash(post.content)
    return post
