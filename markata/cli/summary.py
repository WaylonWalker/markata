from typing import TYPE_CHECKING, Union

from rich.panel import Panel
from rich.table import Table

if TYPE_CHECKING:
    from markata import Markata


class Summary:
    def __init__(self, m: "Markata", simple: bool = False) -> None:
        self.m = m
        self.simple = simple

    def __rich__(self) -> Union[Panel, Table]:
        grid = Table.grid(expand=True)
        grid.add_row(f"[bright_blue]{len(self.m.articles)}[/] articles")
        grid.add_row(
            f"[green]{len([a for a in self.m.articles if a['status'] =='published'])}[/] published"
        )
        grid.add_row(
            f"[gold1]{len([a for a in self.m.articles if a['status'] =='draft'])}[/] drafts"
        )
        grid.add_row("")
        grid.add_row("[bold gold1]TAGS[/]")
        from collections import Counter

        from more_itertools import flatten

        try:
            for tag, count in Counter(
                list(flatten([a["tags"] for a in self.m.articles]))
            ).most_common():
                grid.add_row(f'{count} {" "*(3-len(str(count)))} {tag}')
        except KeyError:
            ...

        try:
            grid.add_row("[bold gold1]Series[/]")
            for series, count in Counter(
                [a["templateKey"] for a in self.m.articles]
            ).most_common():
                grid.add_row(f'{count} {" "*(3-len(str(count)))} {series}')
        except KeyError:
            ...
        if self.simple:
            return grid
        else:
            return Panel(grid, title="[gold1]summary[/]", border_style="magenta")


if __name__ == "__main__":
    from rich import print

    m = Markata()
    print(Summary(m, simple=True))
