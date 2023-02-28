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
    Example::

        a  b
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
        link = text
        possible_pages = markata.filter(
            f'path.split("/")[-1].split(".")[0] == "{text}"'
        )
        if len(possible_pages) == 1:
            link = possible_pages[0].get("slug", f"/{text}")
        elif len(possible_pages) > 1:
            logger.warning(
                f"wikilink [[{text}]] has duplicate matches, defaulting to the first"
            )
            link = possible_pages[0].get("slug", f"/{text}")
        else:
            logger.warning(
                f"wikilink [[{text}]] has no matches, defaulting to '/{text}'"
            )
            link = text

        token.attrSet("href", f"/{link}")
        content_token = state.push("text", "", 0)
        content_token.content = text

        token = state.push("link_close", "a", -1)
        token.content = text

        return True

    md.inline.ruler.before("escape", "wikilinks_inline", _wikilinks_inline)
