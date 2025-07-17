"""
Markata Plugin: Mermaid Diagram Renderer

This plugin converts Mermaid code blocks in Markdown files into rendered Mermaid diagrams.

# Installation

Ensure Mermaid.js is available in your site. If serving locally, add the script to your template:

```html
<script type="module">
  import mermaid from 'https://cdn.jsdelivr.net/npm/mermaid@10/dist/mermaid.esm.min.mjs';
  mermaid.initialize({ startOnLoad: true });
</script>
```

# Configuration

Enable the plugin in `markata.toml`:

```toml
[markata]
hooks = ["markata.plugins.mermaid"]
```

# Usage

Use Mermaid code blocks in your Markdown content:

```markdown
```mermaid
graph TD;
    A-->B;
    A-->C;
    B-->D;
    C-->D;
```
```

# Notes

- Requires the Markata markdown-it-py backend with the `html` option enabled.
"""

import re
from typing import TYPE_CHECKING

from markata.hookspec import hook_impl

if TYPE_CHECKING:
    from markata import Markata

MERMAID_BLOCK_RE = re.compile(r"```[\s]*mermaid\n(.*?)\n```", re.DOTALL)


@hook_impl
def pre_render(markata: "Markata") -> None:
    for article in markata.iter_articles("processing mermaid blocks"):
        key = markata.make_hash("mermaid", article.content)
        if "mermaid" in article.content:
            article.content = MERMAID_BLOCK_RE.sub(
                replace_mermaid_block, article.content
            )


def replace_mermaid_block(match: re.Match) -> str:
    mermaid_code = match.group(1).strip()
    mermaid_block = f'<pre class="mermaid">{mermaid_code}</pre>'
    return mermaid_block


MERMAID_SCRIPT = """
"""
