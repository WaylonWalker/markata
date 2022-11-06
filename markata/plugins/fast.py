"""
fast is a plugin that allows markata to load fast.  It will store all articles in
cache on exit, so when running commands like `markata list` markata only needs to
pull from cache, and not process any files.

Enable this plugin by adding it to your `markata.toml` hooks list.

``` toml
[markata]
hooks=[
  # your hooks
  "markata.plugins.fast",
]
```

## Configuration

If you want a different cache expiration on the articles stored in fast you can
configure this in your `markata.toml` configuration file.

[markata.fast]
default_cache_expire = 12

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
    if getattr(markata, "going_fast", False):
        markata.console.log(f"going_fast")
    else:
        with markata.cache as cache:
            markata.console.log(f"not going_fast")
            articles = copy.deepcopy(markata.articles)
            expire = markata.config.get("fast", markata.config).get(
                "default_cache_expire", 1209600
            )
            cache.add("articles", articles, expire=expire)
