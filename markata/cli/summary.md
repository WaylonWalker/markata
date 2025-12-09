---
date: 2025-12-09
description: "Run it Configuration There are two main things currently supported by
  summary, It can count the number of posts based on a filter ( grid_attr`). grid_attr\u2026"
published: false
slug: markata/cli/summary
title: summary.py


---

---

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

---

!!! function
    <h2 id="cli" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">cli <em class="small">function</em></h2>

    Markata hook to implement base cli commands.

???+ source "cli <em class='small'>source</em>"
    ```python
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
    ```
!!! method
    <h2 id="get_grid" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">get_grid <em class="small">method</em></h2>

    create a rich grid to display the summary

???+ source "get_grid <em class='small'>source</em>"
    ```python
    def get_grid(self) -> None:
            "create a rich grid to display the summary"
            self.grid = Table.grid(expand=True)

            for filter_count in self.m.config.summary.filter_count:
                self.filter_count(filter_count)

            for attr in self.m.config.summary.grid_attr:
                self.grid_attr(attr)

            return self.grid
    ```
!!! method
    <h2 id="filter_count" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">filter_count <em class="small">method</em></h2>

    add a row in the grid for the number of items in a filter config

???+ source "filter_count <em class='small'>source</em>"
    ```python
    def filter_count(
            self,
            fc: FilterCount,
        ) -> None:
            "add a row in the grid for the number of items in a filter config"
            self.grid.add_row(
                f"[{fc.color}]{len(self.m.map(filter=fc.filter))}[/] {fc.name}"
            )
    ```
!!! method
    <h2 id="grid_attr" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">grid_attr <em class="small">method</em></h2>

    add attribute the the object grid

???+ source "grid_attr <em class='small'>source</em>"
    ```python
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
                    self.grid.add_row(f"{count} {' ' * (3 - len(str(count)))} {post}")
    ```
!!! function
    <h2 id="summary" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">summary <em class="small">function</em></h2>

    show the application summary

???+ source "summary <em class='small'>source</em>"
    ```python
    def summary():
            "show the application summary"
            from rich import print

            markata.console.quiet = True

            print(Summary(markata, simple=True))
    ```