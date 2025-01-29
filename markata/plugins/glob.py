"""
The `markata.plugins.glob` plugin handles file discovery using glob patterns.
It determines which files should be processed by Markata based on configured patterns.

> A glob is a string that can be used to match files and directories.
> https://docs.python.org/3/library/pathlib.html#pathlib.Path.glob

# Installation

This plugin is built-in and enabled by default through the 'default' plugin.
If you want to be explicit, you can add it to your list of plugins:

```toml
hooks = [
    "markata.plugins.glob",
]
```

# Uninstallation

Since this plugin is included in the default plugin set, to disable it you must explicitly
add it to the disabled_hooks list if you are using the 'default' plugin:

```toml
disabled_hooks = [
    "markata.plugins.glob",
]
```

# Configuration

Configure glob patterns and behavior in your `markata.toml`:

```toml
[markata.glob]
# Single pattern
glob_patterns = "**/*.md"

# Or multiple patterns
glob_patterns = [
    "content/**/*.md",
    "pages/**/*.md",
    "posts/**/*.md"
]

# Control .gitignore integration
use_gitignore = true  # Set to false to process gitignored files
```

## Configuration Options

- `glob_patterns`: String or list of strings specifying which files to process
  - Default: `["**/*.md"]` (all markdown files in any subdirectory)
  - Supports any valid glob pattern
- `use_gitignore`: Boolean controlling whether to respect .gitignore rules
  - Default: `true`
  - When true, files matching .gitignore patterns are skipped

# Functionality

The glob plugin runs early in the Markata pipeline to discover files for processing.
It supports:
- Single or multiple glob patterns
- .gitignore integration
- Recursive directory searching
- Common glob patterns (*, **, ?, [...])

## Dependencies

This plugin has no dependencies on other Markata plugins.

"""

from markata import background
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

    @pydantic.field_validator("glob_patterns", mode="before")
    @classmethod
    def validate_glob_patterns(cls, v) -> List[str]:
        if isinstance(v, str):
            return [v]
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

        @background.task
        def check_spec(file: str) -> bool:
            key = markata.make_hash("glob", "check_spec", file)
            check = markata.precache.get(key)
            if check is not None:
                return check

            check = spec.match_file(str(file))
            with markata.cache as cache:
                cache.set(key, check)
            return check

        file_checks = [(file, check_spec(str(file))) for file in markata.files]
        [check.result() for _, check in file_checks]
        markata.files = [file for file, check in file_checks if check]
