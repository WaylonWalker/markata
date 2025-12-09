---
date: 2025-12-09
description: "The plugin adds syntax highlighting to code blocks in your markdown
  content. It uses Pygments for highlighting and adds a copy button for easy code
  sharing.\u2026"
published: false
slug: markata/plugins/md-it-highlight-code
title: md_it_highlight_code.py


---

---

The `markata.plugins.md_it_highlight_code` plugin adds syntax highlighting to code blocks
in your markdown content. It uses Pygments for highlighting and adds a copy button for
easy code sharing.

## Installation

This plugin is built-in and enabled by default through the 'default' plugin.
If you want to be explicit, you can add it to your list of plugins:

```toml
hooks = [
    "markata.plugins.md_it_highlight_code",
]
```

## Uninstallation

Since this plugin is included in the default plugin set, to disable it you must explicitly
add it to the disabled_hooks list if you are using the 'default' plugin:

```toml
disabled_hooks = [
    "markata.plugins.md_it_highlight_code",
]
```

## Configuration

Configure syntax highlighting in your `markata.toml`:

```toml
[markata.highlight_code]
# Optional: Custom CSS class for code blocks
code_class = "highlight"

# Optional: Custom CSS class for copy button
copy_button_class = "copy-button"

# Optional: Custom copy button text
copy_button_text = "Copy"
copy_button_copied_text = "Copied!"
```

## Functionality

## Code Highlighting

The plugin:
1. Detects language from code block info
2. Applies Pygments syntax highlighting
3. Adds line numbers (optional)
4. Wraps code in proper HTML structure
5. Adds a copy-to-clipboard button

## Example Usage

In markdown:
````markdown
```python
def hello_world():
    print("Hello, World!")
```

```javascript
console.log('Hello, World!');
```
````

Generates HTML:
```html
<div class="highlight">
  <button class="copy-button">
    <svg><!-- Copy icon SVG --></svg>
  </button>
  <pre><code class="language-python">
    <span class="def">def</span> <span class="name">hello_world</span>():
        <span class="keyword">print</span>(<span class="string">"Hello, World!"</span>)
  </code></pre>
</div>
```

## Supported Languages

The plugin supports all languages that Pygments recognizes, including:
- Python
- JavaScript
- HTML/CSS
- Bash/Shell
- And many more

## Copy Button

Features:
- Hover-to-show button
- Click-to-copy functionality
- Success feedback
- Accessible keyboard support

## Dependencies

This plugin depends on:
- Pygments for syntax highlighting
- markdown-it-py for markdown parsing

---

!!! function
    <h2 id="highlight_code" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">highlight_code <em class="small">function</em></h2>

    Code highlighter for markdown-it-py.

???+ source "highlight_code <em class='small'>source</em>"
    ```python
    def highlight_code(code, name, attrs, markata=None):
        """Code highlighter for markdown-it-py."""
        # from pygments import highlight
        # from pygments.formatters import HtmlFormatter
        # from pygments.lexers import ClassNotFound, get_lexer_by_name

        try:
            lexer = get_lexer_by_name(name or "text")
        except ClassNotFound:
            lexer = get_lexer_by_name("text")

        import re

        pattern = r'(\w+)\s*=\s*(".*?"|\S+)'
        matches = re.findall(pattern, attrs)
        attrs = dict(matches)

        if attrs.get("hl_lines"):
            formatter = HtmlFormatter(hl_lines=attrs.get("hl_lines"))
        else:
            formatter = HtmlFormatter()

        copy_button = f"""<button class='copy' title='copy code to clipboard' onclick="navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)">{COPY_ICON}</button>"""

        from markdown_it import MarkdownIt

        md = MarkdownIt(
            "commonmark",
            {
                "html": True,
                "typographer": True,
            },
        )

        if attrs.get("help"):
            help = f"""
            <a href={attrs.get("help").strip("<").strip(">").strip('"').strip("'")} title='help link' class='help'>{HELP_ICON}</a>
            """
        else:
            help = ""
        if attrs.get("title"):
            file = f"""
    <div class='filepath'>
    {md.render(attrs.get("title").strip('"').strip("'"))}
    <div class='right'>
    {help}
    {copy_button}
    </div>
    </div>
    """
        else:
            file = f"""
    <div class='copy-wrapper'>
    {help}
    {copy_button}
    </div>
            """
        return f"""<pre class='wrapper'>
    {file}
    {highlight(code, lexer, formatter)}
    </pre>
    """
    ```