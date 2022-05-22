"""
The articles slug is what determines the url to your page.  It should be
sanitized of special characters that do not work in the browser.


## Explicit Slug in Frontmatter

If you explicitly set the slug in the frontmatter of a post, markata will not
overwrite it.

```markdown
---
title: My First Post
slug: /my-post

---

This is my first post it will be at `<markata.config['url']>/my-post/`
reguardless of filename.
```

## Automatic Slug Based on Filename

By default the flat_slug plugin will use the `stem` of your filename, which is
the filename without the extension, unless you explicitly set your slug in
frontmatter.

* `/pages/my-post.md` becomes `<markata.config['url']>/my-post/`
* `/pages/blog/a-blog-post.md` becomes `<markata.config['url']>/a-blog-post/`
"""
from pathlib import Path
from typing import TYPE_CHECKING

from markata.hookspec import hook_impl

if TYPE_CHECKING:
    from markata import Markata


@hook_impl(tryfirst=True)
def pre_render(markata: "Markata") -> None:
    """
    Sets the article slug if one is not already set in the frontmatter.
    """
    for article in markata.iter_articles(description="creating slugs"):
        article["slug"] = article.get("slug", Path(article["path"]).stem)
