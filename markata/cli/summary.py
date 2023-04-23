"""

## Run it

``` bash
python -m markata.cli.summary
```

## Configuration

There are two main things currently supported by summary, It can count the
number of posts based on a filter (`filter_count'), and it can automatically
list all the values of an attribute and the number of posts that have that
attribute (`grid_attr`).

### grid_attr

`grid_attr` will map over all posts, find all values for each attribute
configured, then report the number of posts for each value.

``` toml
[markata.summary]
grid_attr = ['tags', 'series']
```

Example output that will be shown in the summary, counting all posts with each
tag value.

```
TAGS
247  python
90   linux
68   cli
49   kedro
46   bash
```

### filter_counts

`filter_count` will pass a given filter into `markata.map` and return the number
of posts.

```
[markata.summary.filter_count.drafts]
filter="published == 'False'"
color='red'

[markata.summary.filter_count.articles]
color='dark_orange'

[markata.summary.filter_count.py_modules]
filter='"plugin" not in slug and "docs" not in path'
color="yellow1"

[markata.summary.filter_count.published]
filter="published == 'True'"
color='green1'

[markata.summary.filter_count.plugins]
filter='"plugin" in slug and "docs" not in path'
color="blue"

[markata.summary.filter_count.docs]
filter="'docs' in path"
color='purple'
```

Example output might look like this, showing the number of posts that contained
within each filter specified.

```
8 drafts
66 articles
20 py_modules
58 published
38 plugins
8 docs
```
"""
from collections import Counter
from typing import TYPE_CHECKING, Union

from more_itertools import flatten
from rich.panel import Panel
from rich.table import Table

from markata.hookspec import hook_impl, register_attr

if TYPE_CHECKING:
    from markata import Markata


class Summary:
    def __init__(self, m: "Markata", simple: bool = False) -> None:
        self.m = m
        self.simple = simple

    def get_grid(self) -> None:
        "create a rich grid to display the summary"
        self.grid = Table.grid(expand=True)

        for name, config in (
            self.m.config.get("summary", {})
            .get("filter_count", {"aricles": {"color": "purple", "filter": "True"}})
            .items()
        ):
            self.filter_count(name, **config)

        for attr in self.m.config.get("summary", {}).get("grid_attr", []):
            self.grid_attr(attr)

        return self.grid

    def filter_count(
        self, title: str, filter: str = "True", color: str = "white",
    ) -> None:
        "add a row in the grid for the number of items in a filter config"
        self.grid.add_row(f"[{color}]{len(self.m.map(filter=filter))}[/] {title}")

    def grid_attr(self, attr: str) -> None:
        "add attribute the the object grid"
        posts = list(
            flatten(
                [
                    tags if isinstance(tags, list) else [tags]
                    for a in self.m.articles
                    if (tags := a.get(attr, None)) is not None
                ],
            ),
        )
        if len(posts) > 0:
            self.grid.add_row()
            self.grid.add_row(f"[bold gold1]{attr.upper()}[/]")
            for post, count in Counter(posts).most_common():
                self.grid.add_row(f'{count} {" "*(3-len(str(count)))} {post}')

    def __rich__(self) -> Union[Panel, Table]:
        try:
            grid = self.get_grid()
        except Exception:
            grid = "Error"
        if self.simple:
            return grid
        else:
            return Panel(
                grid, title="[gold1]summary[/]", border_style="magenta", expand=False,
            )


@hook_impl
@register_attr("summary")
def configure(markata: "Markata") -> None:
    def get_summary(self):
        try:
            return self._summary
        except AttributeError:
            self._summary: Summary = Summary(self)
            return self._summary

    from markata import Markata

    Markata.summary = property(get_summary)


if __name__ == "__main__":
    from rich import print

    m = Markata()
    print(Summary(m, simple=True))
