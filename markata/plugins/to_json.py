import json
from typing import TYPE_CHECKING

from markata.hookspec import hook_impl

if TYPE_CHECKING:
    from markata import Markata


@hook_impl
def save(markata: "Markata") -> None:
    output_file = markata.config.output_dir / "markata.json"
    new_content = json.dumps(markata.to_dict(), default=str)
    current_content = output_file.read_text() if output_file.exists() else ""
    if current_content != new_content:
        output_file.write_text(new_content)
