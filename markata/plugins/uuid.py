"""
uuid

Enable this plugin by adding it to your `markata.toml` hooks list.

``` toml
[markata]
hooks=[
  # your hooks
  "markata.plugins.uuid",
]
```

"""
import uuid
from dataclasses import dataclass
from typing import TYPE_CHECKING

from markata import Markata
from markata.hookspec import hook_impl, register_attr

if TYPE_CHECKING:
    from frontmatter import Post


@dataclass
class PrevNext:
    prev: str
    curr: str
    next: str


@hook_impl
@register_attr("article_by_uuid")
def pre_render(markata: "Markata") -> None:
    markata.article_by_uuid = {}  # type: ignore
    for article in markata.articles:
        _uuid = uuid.uuid4()
        article["uuid"] = _uuid
        markata.article_by_uuid[_uuid] = article


@hook_impl
@register_attr("get_uuid")
def configure(markata: "Markata") -> None:
    """
    setup get_uuid
    """

    def get_uuid(self: "Markata", _uuid: uuid.UUID) -> "Post":
        return markata.article_by_uuid[_uuid]

    Markata.get_uuid = property(get_uuid)  # type: ignore
