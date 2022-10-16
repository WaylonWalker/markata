from collections import Counter
from typing import TYPE_CHECKING, Union

from more_itertools import flatten
from rich.panel import Panel
from rich.table import Table

if TYPE_CHECKING:
    from markata import Markata


class Summary:
    def __init__(self, m: "Markata", simple: bool = False) -> None:
        self.m = m
        self.simple = simple

    def get_grid(self):
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

    def filter_count(self, title, filter="True", color="white") -> None:
        self.grid.add_row(f"[{color}]{len(self.m.map(filter=filter))}[/] {title}")

    def grid_attr(self, attr) -> None:
        posts = list(
            flatten(
                [
                    tags if isinstance(tags, list) else [tags]
                    for a in self.m.articles
                    if (tags := a.get(attr, None)) is not None
                ]
            )
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
                grid, title="[gold1]summary[/]", border_style="magenta", expand=False
            )


if __name__ == "__main__":
    from rich import print

    m = Markata()
    print(Summary(m, simple=True))
