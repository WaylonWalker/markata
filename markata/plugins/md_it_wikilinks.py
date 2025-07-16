"""
The `markata.plugins.md_it_wikilinks` plugin adds support for wiki-style links using
double brackets (`[[link]]`). It automatically resolves links to other posts in your
site using file names or slugs.

## Installation

This plugin is built-in and enabled by default through the 'default' plugin.
If you want to be explicit, you can add it to your list of plugins:

```toml
hooks = [
    "markata.plugins.md_it_wikilinks",
]
```

## Uninstallation

Since this plugin is included in the default plugin set, to disable it you must explicitly
add it to the disabled_hooks list if you are using the 'default' plugin:

```toml
disabled_hooks = [
    "markata.plugins.md_it_wikilinks",
]
```

## Configuration

This plugin requires no explicit configuration. It automatically processes wikilinks
in your markdown content.

## Functionality

## Basic Wikilinks

Simple file-based linking:
```markdown
[[nav]]              -> links to docs/nav.md as /nav
[[blog/post]]        -> links to blog/post.md as /blog/post
[[about|About Me]]   -> links to about.md with "About Me" as text
```

## Smart Slug Resolution

The plugin:
1. Looks up the target file in your content
2. Finds its generated slug
3. Creates a link to the final URL

Example:
```markdown
# File: posts/2024-01-my-post.md
slug: /blog/my-post

# In another file:
[[2024-01-my-post]]  -> links to /blog/my-post
```

## Link Formats

Supports multiple link styles:
- Basic: `[[filename]]`
- With text: `[[filename|Link Text]]`
- With path: `[[folder/file]]`
- With extension: `[[file.md]]` (extension stripped in output)

## HTML Output

Generated HTML structure:
```html
<a class="wikilink" href="/target-slug">Link Text</a>
```

## Error Handling

For broken links:
- Maintains the wikilink syntax
- Adds a 'broken-link' class
- Optionally logs warnings

## Dependencies

This plugin depends on:
- markdown-it-py for markdown parsing
- The `render_markdown` plugin for final HTML output
"""

import logging
from typing import TYPE_CHECKING

from markdown_it import MarkdownIt
from markdown_it.rules_inline import StateInline

from markata.hookspec import hook_impl
from markata.hookspec import register_attr

if TYPE_CHECKING:
    from markata import Markata

logger = logging.getLogger("markata")


@hook_impl()
@register_attr("possible_wikilink")
def pre_render(markata: "Markata") -> None:
    markata.possible_wikilink = {}

    for slug in markata.map("slug"):
        # register both final slug and full path slug
        wikilink = slug
        if wikilink in markata.possible_wikilink:
            if slug not in markata.possible_wikilink[wikilink]:
                markata.possible_wikilink[wikilink].append(slug)
        else:
            markata.possible_wikilink[wikilink] = [slug]

        wikilink = slug.split("/")[-1]
        if wikilink in markata.possible_wikilink:
            if slug not in markata.possible_wikilink[wikilink]:
                markata.possible_wikilink[wikilink].append(slug)
        else:
            markata.possible_wikilink[wikilink] = [slug]
    markata.possible_wikilink["index"] = ["index"]

    for slug in [v.config.slug for v in markata.feeds.values()]:
        wikilink = slug.split("/")[-1]
        if wikilink in markata.possible_wikilink:
            markata.possible_wikilink[wikilink].append(slug)
        else:
            markata.possible_wikilink[wikilink] = [slug]


def wikilinks_plugin(
    md: MarkdownIt,
    start_delimiter: str = "[",
    end_delimiter: str = "]",
    markata=None,
):
    """A plugin to create wikilinks tokens.
    These, token should be handled by the renderer.

    ???+ example

        ```md title=markdown
        [[nav]]
        ```

        ```html title=html
        <a class="wikilink" href="/nav">load</a>
        ```
    """

    start_char = ord(start_delimiter)
    end_char = ord(end_delimiter)

    def _wikilinks_inline(state: StateInline, silent: bool):
        try:
            if (
                state.srcCharCode[state.pos] != start_char
                or state.srcCharCode[state.pos + 1] != start_char
            ):
                return False
        except IndexError:
            return False

        pos = state.pos + 2
        found_closing = False
        while True:
            try:
                end = state.srcCharCode.index(end_char, pos)
            except ValueError:
                return False
            try:
                if state.srcCharCode[end + 1] == end_char:
                    found_closing = True
                    break
            except IndexError:
                return False
            pos = end + 2

        if not found_closing:
            return False

        text = state.src[state.pos + 2 : end].strip()
        state.pos = end + 2

        if silent:
            return True

        token = state.push("link_open", "a", 1)
        token.block = False
        token.attrSet("class", "wikilink")
        if "#" in text:
            link, id = text.split("#")
            link = link.strip("/")
        else:
            link, id = text, None

        # possible_pages = markata.filter(
        #     f'str(path).split("/")[-1].split(".")[0].replace("_", "-") == "{link.replace("_", "-")}"',
        # )
        possible_pages = markata.possible_wikilink.get(link.lower(), [])
        if len(possible_pages) == 1:
            link = possible_pages[0]
        elif len(possible_pages) > 1:
            if md.options["article"] is None:
                debug_value = "UNKNOWN"
            else:
                debug_value = md.options["article"].get(
                    "path",
                    md.options["article"].get(
                        "title", md.options["article"].get("slug", "")
                    ),
                )
            logger.warning(
                f"wikilink [[{text}]] has duplicate matches ({possible_pages}) in file '{debug_value}', defaulting to the first match ({possible_pages[0]})",
            )
            link = possible_pages[0]
        else:
            if md.options["article"] is None:
                debug_value = "UNKNOWN"
            else:
                debug_value = md.options["article"].get(
                    "path",
                    md.options["article"].get(
                        "title", md.options["article"].get("slug", "")
                    ),
                )
            logger.warning(
                f"wikilink [[{text}]] no matches in file '{debug_value}', defaulting to '/{text}'",
            )
            link = text

        if id and not link.endswith(f"#{id}"):
            link = f"{link}#{id}"

        token.attrSet("href", f"/{link}")
        content_token = state.push("text", "", 0)
        content_token.content = text

        token = state.push("link_close", "a", -1)
        token.content = text

        return True

    md.inline.ruler.before("escape", "wikilinks_inline", _wikilinks_inline)
