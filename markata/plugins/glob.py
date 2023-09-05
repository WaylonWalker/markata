"""Default glob plugin"""
from pathlib import Path
from typing import List, TYPE_CHECKING, Union

from more_itertools import flatten
import pydantic

from markata.hookspec import hook_impl, register_attr

if TYPE_CHECKING:
    from markata import Markata


class GlobConfig(pydantic.BaseModel):
    glob_patterns: Union[List[str], str] = ["**/*.md"]
    use_gitignore: bool = True

    @pydantic.validator("glob_patterns")
    def convert_to_list(cls, v):
        if not isinstance(v, list):
            return v.split(",")
        return v


class Config(pydantic.BaseModel):
    glob: GlobConfig = GlobConfig()


@hook_impl
@register_attr("post_models")
def config_model(markata: "Markata") -> None:
    markata.config_models.append(Config)


@hook_impl
@register_attr("content_directories", "files")
def glob(markata: "Markata") -> None:
    markata.files = list(
        flatten(
            [
                Path().glob(str(pattern))
                for pattern in markata.config.glob.glob_patterns
            ],
        ),
    )
    markata.content_directories = list({f.parent for f in markata.files})

    try:
        ignore = True
    except KeyError:
        ignore = True

    if ignore and (Path(".gitignore").exists() or Path(".markataignore").exists()):
        import pathspec

        lines = []

        if Path(".gitignore").exists():
            lines.extend(Path(".gitignore").read_text().splitlines())

        if Path(".markataignore").exists():
            lines.extend(Path(".markataignore").read_text().splitlines())

        key = markata.make_hash("glob", "spec", lines)
        spec = markata.precache.get(key)
        if spec is None:
            spec = pathspec.PathSpec.from_lines("gitwildmatch", lines)
            with markata.cache as cache:
                cache.set(key, spec)

        def check_spec(file: str) -> bool:
            key = markata.make_hash("glob", "check_spec", file)
            check = markata.precache.get(key)
            if check is not None:
                return check

            check = spec.match_file(str(file))
            with markata.cache as cache:
                cache.set(key, check)
            return check

        markata.files = [file for file in markata.files if not check_spec(str(file))]
