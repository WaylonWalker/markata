import json
from pathlib import Path
from typing import TYPE_CHECKING

from markata.hookspec import hook_impl

if TYPE_CHECKING:
    from markata import Markata


@hook_impl
def save(markata: "Markata") -> None:
    output_file = Path(markata.config["output_dir"]) / "markata.json"
    output_file.parent.mkdir(parents=True, exist_ok=True)
    output_file.write_text(json.dumps(markata.to_dict(), default=str))
