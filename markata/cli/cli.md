---
date: 2025-12-09
description: None
published: false
slug: markata/cli/cli
title: cli.py


---

---

None

---

!!! function
    <h2 id="make_layout" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">make_layout <em class="small">function</em></h2>

    Define the layout.

???+ source "make_layout <em class='small'>source</em>"
    ```python
    def make_layout() -> Layout:
        """Define the layout."""
        layout = Layout(name="root")

        layout.split(
            Layout(name="header", size=3),
            Layout(name="main"),
        )
        layout["main"].split_row(
            Layout(name="side", ratio=50),
            Layout(name="mid", ratio=30),
            Layout(name="describe", ratio=20),
        )
        layout["mid"].split(
            Layout(name="server"),
            Layout(name="runner"),
        )
        layout["side"].split(
            Layout(name="plugins"),
        )
        return layout
    ```