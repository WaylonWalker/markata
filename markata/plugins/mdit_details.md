---
date: 2025-12-09
description: "The plugin adds support for collapsible details/summary sections in
  markdown using the syntax. This is similar to HTML's element but with a more\u2026"
published: false
slug: markata/plugins/mdit-details
title: mdit_details.py


---

---

The `markata.plugins.mdit_details` plugin adds support for collapsible details/summary
sections in markdown using the `???` syntax. This is similar to HTML's `<details>` element
but with a more markdown-friendly syntax.

## Installation

This plugin is built-in and enabled by default through the 'default' plugin.
If you want to be explicit, you can add it to your list of plugins:

```toml
hooks = [
    "markata.plugins.mdit_details",
]
```

## Uninstallation

Since this plugin is included in the default plugin set, to disable it you must explicitly
add it to the disabled_hooks list if you are using the 'default' plugin:

```toml
disabled_hooks = [
    "markata.plugins.mdit_details",
]
```

## Configuration

This plugin requires no explicit configuration. It automatically processes details
blocks in your markdown content.

## Functionality

## Basic Usage

Create collapsible sections:
```markdown
??? note "Optional Title"
    This content will be collapsible.
    It can contain *any* markdown.

???+ note "Open by Default"
    The + symbol makes this section expanded by default.
```

## Supported Styles

Default styles include:
- note
- info
- warning
- danger
- success
- question
- abstract
- example

## Syntax Options

The plugin supports:
- Basic: `??? type`
- With title: `??? type "Custom Title"`
- Open by default: `???+ type`
- No title: Title defaults to capitalized type

## HTML Output

Generated HTML structure:
```html
<details class="details type">
  <summary>Title Text</summary>
  <div class="details-content">
    <!-- Markdown content -->
  </div>
</details>
```

## Nesting Support

Details blocks can be nested:
```markdown
??? outer "Outer Section"
    Some content

    ??? inner "Inner Section"
        Nested content
```

## Dependencies

This plugin depends on:
- markdown-it-py for markdown parsing
- The `render_markdown` plugin for final HTML output

---

!!! function
    <h2 id="details_plugin" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">details_plugin <em class="small">function</em></h2>

    Plugin to use
    `python-markdown style detailss
    <https://python-markdown.github.io/extensions/details>`_.

    .. code-block:: md

        ??? note
            *content*

    Note, this is ported from
    `markdown-it-admon
    <https://github.com/commenthol/markdown-it-admon>`_.

???+ source "details_plugin <em class='small'>source</em>"
    ```python
    def details_plugin(md: MarkdownIt, render: Optional[Callable] = None) -> None:
        """Plugin to use
        `python-markdown style detailss
        <https://python-markdown.github.io/extensions/details>`_.

        .. code-block:: md

            ??? note
                *content*

        Note, this is ported from
        `markdown-it-admon
        <https://github.com/commenthol/markdown-it-admon>`_.
        """

        def renderDefault(self, tokens, idx, _options, env):
            return self.renderToken(tokens, idx, _options, env)

        render = render or renderDefault

        md.add_render_rule("details_open", render)
        md.add_render_rule("details_close", render)
        md.add_render_rule("details_title_open", render)
        md.add_render_rule("details_title_close", render)

        md.block.ruler.before(
            "fence",
            "details",
            details,
            {"alt": ["paragraph", "reference", "blockquote", "list"]},
        )
    ```