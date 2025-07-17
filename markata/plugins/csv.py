"""
CSV Fence to Table Plugin

Converts fenced code blocks with ```csv into rendered HTML tables.

# Installation

```toml
[markata]
hooks = ["markata.plugins.csv"]
```

# Configuration

None required.

# Usage

Use markdown fences like this:

```csv
name,age,city
Alice,30,New York
Bob,25,San Francisco
Charlie,35,Chicago
```

This will be replaced with an HTML table.

# Notes

- Assumes CSV is well-formed.
- No external dependencies required.
"""

import csv
import io
import re
from typing import TYPE_CHECKING

from markata.hookspec import hook_impl

if TYPE_CHECKING:
    from markata import Markata


CSV_BLOCK_RE = re.compile(r"```[\s]*csv\n(.*?)\n```", re.DOTALL)


def csv_to_html_table(csv_string: str) -> str:
    f = io.StringIO(csv_string)
    reader = csv.reader(f)
    rows = list(reader)
    if not rows:
        return "<table></table>"

    headers = rows[0]
    html = "<table><thead><tr>"
    html += "".join(f"<th>{h}</th>" for h in headers)
    html += "</tr></thead><tbody>"

    for row in rows[1:]:
        html += "<tr>" + "".join(f"<td>{cell}</td>" for cell in row) + "</tr>"
    html += "</tbody></table>"

    return html


@hook_impl
def pre_render(markata: "Markata") -> None:
    for post in markata.posts:
        if hasattr(post, "content") and (
            "```csv" in post.content.lower() or "``` csv" in post.content.lower()
        ):
            post.content = CSV_BLOCK_RE.sub(
                lambda m: csv_to_html_table(m.group(1).strip()), post.content
            )
