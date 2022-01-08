"""Define hook specs."""
import pluggy

HOOK_NAMESPACE = "markata"

hook_spec = pluggy.HookspecMarker(HOOK_NAMESPACE)
hook_impl = pluggy.HookimplMarker(HOOK_NAMESPACE)

from typing import TYPE_CHECKING

if TYPE_CHECKING:
    from markata import Markata


class MarkataSpecs:
    """
    Namespace that defines all specifications for Load hooks.

    configure -> glob -> load -> render -> save
    """

    @hook_spec
    def configure(self, markata: "Markata") -> None:
        """Glob for files to load."""
        pass

    @hook_spec
    def glob(self, markata: "Markata") -> None:
        """Glob for files to load."""
        pass

    @hook_spec
    def load(self, markata: "Markata") -> None:
        """Load list of files."""
        pass

    @hook_spec
    def pre_render(self, markata: "Markata") -> None:
        """Pre render content from loaded data."""
        pass

    @hook_spec
    def render(self, markata: "Markata") -> None:
        """Render content from loaded data."""
        pass

    @hook_spec
    def post_render(self, markata: "Markata") -> None:
        """Post render content from loaded data."""
        pass

    @hook_spec
    def save(self, markata: "Markata") -> None:
        """Save content from data."""
        pass

    @hook_spec
    def cli(self, markata: "Markata", app) -> None:
        """cli"""
        pass
