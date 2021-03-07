"""Default load plugin."""
from markata.hookspec import hook_impl
from markata.background import task
from pathlib import Path
from more_itertools import flatten

from yaml.parser import ParserError
import time
import frontmatter

from tqdm import tqdm


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
