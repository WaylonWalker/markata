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
[[markata.summary.filter_count]]
name='drafts'
filter="published == 'False'"
color='red'

[[markata.summary.filter_count]]
name='articles'
color='dark_orange'

[[markata.summary.filter_count]]
name='py_modules'
filter='"plugin" not in slug and "docs" not in str(path)'
color="yellow1"

[markata.summary.filter_count.published]
filter="published == 'True'"
color='green1'

[markata.summary.filter_count.plugins]
filter='"plugin" in slug and "docs" not in str(path)'
color="blue"

[markata.summary.filter_count.docs]
filter="'docs' in str(path)"
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
from typing import List, TYPE_CHECKING, Union

from more_itertools import flatten
import pydantic
from rich.panel import Panel
from rich.table import Table
import typer

from markata.hookspec import hook_impl, register_attr

if TYPE_CHECKING:
    from markata import Markata


class FilterCount(pydantic.BaseModel):
    name: str
    filter: str = "True"
    color: str = "white"


class SummaryConfig(pydantic.BaseModel):
    grid_attr: List[str] = ["tags", "series"]
    filter_count: List[FilterCount] = FilterCount(
        name="drafts", filter="published == 'False'", color="red"
    )


class Config(pydantic.BaseModel):
    summary: SummaryConfig = SummaryConfig()


@hook_impl()
@register_attr("config_models")
def config_model(markata: "Markata") -> None:
    markata.config_models.append(Config)


class Summary:
    def __init__(self, m: "Markata", simple: bool = False) -> None:
        self.m = m
        self.simple = simple

    def get_grid(self) -> None:
        "create a rich grid to display the summary"
        self.grid = Table.grid(expand=True)

        for filter_count in self.m.config.summary.filter_count:
            self.filter_count(filter_count)

        for attr in self.m.config.summary.grid_attr:
            self.grid_attr(attr)

        return self.grid

    def filter_count(
        self,
        fc: FilterCount,
    ) -> None:
        "add a row in the grid for the number of items in a filter config"
        self.grid.add_row(
            f"[{fc.color}]{len(self.m.map(filter=fc.filter))}[/] {fc.name}"
        )

    def grid_attr(self, attr: str) -> None:
        "add attribute the the object grid"
        posts = list(
            flatten(
                [
                    tags if isinstance(tags, list) else [tags]
                    for a in self.m.posts
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
        grid = self.get_grid()

        if self.simple:
            return grid
        else:
            return Panel(
                grid,
                title="[gold1]summary[/]",
                border_style="magenta",
                expand=False,
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


@hook_impl()
def cli(app: typer.Typer, markata: "Markata") -> None:
    """
    Markata hook to implement base cli commands.
    """
    summary_app = typer.Typer()
    app.add_typer(summary_app, name="summary")

    @summary_app.callback(invoke_without_command=True)
    def summary():
        "show the application summary"
        from rich import print

        markata.console.quiet = True

        print(Summary(markata, simple=True))


if __name__ == "__main__":
    from rich import print

    m = Markata()
    print(Summary(m, simple=True))
