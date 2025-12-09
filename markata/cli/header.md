---
date: 2025-12-09
description: None
published: false
slug: markata/cli/header
title: header.py


---

---

None

---

!!! class
    <h2 id="Header" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">Header <em class="small">class</em></h2>

    Display header with clock.

???+ source "Header <em class='small'>source</em>"
    ```python
    class Header:
        """Display header with clock."""

        def __rich__(self) -> Panel:
            grid = Table.grid(expand=True)
            grid.add_column(justify="center", ratio=1)
            grid.add_column(justify="right")
            grid.add_row(
                "[magenta][b]Markata[/b][/] [bright_black]Live Server[/]",
                datetime.now().ctime(),
            )
            return Panel(grid, style="yellow")
    ```