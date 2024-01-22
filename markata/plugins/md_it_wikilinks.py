"""
Wikilinks depicted with `\[[wikilinks]]` can be enabled for sites using the
`markdown-it-py` plugin markata will look through all posts matching up the
file stem to the wiki link and inserting the slug as the href.

???+ note normal wikilink

    Here is a link to a markdown file `docs/nav.md`, and the url becomes
    [/nav](/nav).

    ```md
    [[nav]]
    ```

    When rendered out this wikilink will become an anchor link.

    ```html
    <a class="wikilink" href="/nav">load</a>
    ```

    > This behaves just like the standard wikilink

A special feature that markata brings is slug lookup.  It is able to not only
blindly link to the route specified, but will look up the slug of an article.


???+ note markata slug lookup
    Markata has a load plugin that is generated with the [[docs]] plugin.  It's
    filepath is `markata/plugins/load.py`, so it can be referenced by the file
    stem `load`.

    ```md
    [[load]]
    ```

    Markata will look up the article by the file stem, grab the first article,
    and use its slug as the href.  This turns it into an anchor link that looks
    like this.

    ```html
    <a class="wikilink" href="/markata/plugins/load">load</a>
    ```
"""
import logging

from markdown_it import MarkdownIt
from markdown_it.rules_inline import StateInline

logger = logging.getLogger("markata")


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
        possible_pages = markata.filter(
            f'str(path).split("/")[-1].split(".")[0].replace("_", "-") == "{link.replace("_", "-")}"',
        )
        if len(possible_pages) == 1:
            link = possible_pages[0].get("slug", f"/{text}")
        elif len(possible_pages) > 1:
            logger.warning(
                f"wikilink [[{text}]] ({link}, {id}) has duplicate matches, defaulting to the first",
            )
            link = possible_pages[0].get("slug", f"/{text}")
        else:
            logger.warning(
                f"wikilink [[{text}]] ({link}, {id}) no matches, defaulting to '/{text}'",
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
