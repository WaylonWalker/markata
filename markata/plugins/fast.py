"""
fast

Enable this plugin by adding it to your `markata.toml` hooks list.

``` toml
[markata]
hooks=[
  # your hooks
  "markata.plugins.fast",
]
```

"""
import copy
from typing import Dict

import frontmatter

from markata.hookspec import hook_impl


def keys_on_file(post: frontmatter.Post) -> Dict[str, str]:
    if post.get("path") is not None:
        try:
            return {
                key: value
                for key, value in frontmatter.load(post.get("path")).metadata.items()
            }
        except FileNotFoundError:
            return


@hook_impl
def teardown(markata) -> None:
    with markata.cache as cache:
        articles = copy.deepcopy(markata.articles)
        cache.add("articles", articles)
