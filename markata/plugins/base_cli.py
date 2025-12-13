"""
The `markata.plugins.base_cli` plugin provides Markata's core command-line interface
functionality, including essential commands like `build`, `list`, and `clean`.

## Installation

This plugin is built-in and enabled by default through the 'default' plugin.
If you want to be explicit, you can add it to your list of plugins:

```toml
hooks = [
    "markata.plugins.base_cli",
]
```

## Uninstallation

Since this plugin is included in the default plugin set, to disable it you must explicitly
add it to the disabled_hooks list if you are using the 'default' plugin:

```toml
disabled_hooks = [
    "markata.plugins.base_cli",
]
```

## Configuration

No explicit configuration is required. The plugin automatically registers CLI commands.

## Functionality

## Core Commands

### Build Command

Build your Markata site:
```bash
markata build [options]

# Options:
--clean         Clean output directory first
--watch         Watch for changes and rebuild
--serve         Start development server
--profile       Profile the build process
--debug         Enable debug mode
-c, --config    Path to alternate config file
-o, --output-dir  Override output directory
-s, --set       Set config values (key=value format)
```

### Configuration Overrides

Override configuration at runtime:

```bash
# Use alternate config file
markata build -c themes/catppuccin.toml

# Override output directory
markata build -o dist/theme-everforest

# Set multiple config values
markata build -s output_dir=dist -s style.theme=nord

# Combine multiple overrides
markata build -c base.toml -s output_dir=custom -s style.theme=gruvbox
```

### Environment Variable Overrides

All config can be overridden with environment variables prefixed with `MARKATA_`:

```bash
# Override output directory
MARKATA_OUTPUT_DIR=dist markata build

# Override theme
MARKATA_STYLE__THEME=nord markata build

# Use double underscore for nested config
MARKATA_STYLE__THEME=catppuccin MARKATA_OUTPUT_DIR=dist/catppuccin markata build
```

### List Command

List and filter articles:
```bash
markata list [options]

# Options:
--filter "post.published"  Filter posts by expression
--sort "post.date"        Sort posts by attribute
--reverse                 Reverse sort order
--format "{title}"        Custom output format
```

### Clean Command

Clean build artifacts:
```bash
markata clean [options]

# Options:
--dry-run      Show what would be deleted
--quiet        Suppress output
```

## Debug Features

The plugin provides:
- Pretty error formatting
- Post-mortem debugging
- Performance profiling
- Verbose logging

## Development Server

Features include:
- Live reload
- Asset serving
- Port configuration
- Host configuration

## Watch Mode

Supports:
- File watching
- Auto-rebuild
- Pattern matching
- Debouncing

## Dependencies

This plugin depends on:
- typer for CLI interface
- rich for terminal output
- watchfiles for watch mode
"""

import json
import pdb
import shutil
import sys
import traceback
import warnings
from pathlib import Path
from typing import TYPE_CHECKING
from typing import Any
from typing import Callable
from typing import Dict
from typing import List
from typing import Literal
from typing import Optional

import pydantic
import toml
import typer
from rich import print as rich_print

from markata.hookspec import hook_impl

if TYPE_CHECKING:
    from markata import Markata


def make_pretty() -> None:
    """
    This is a helper function that enables suppresses tracebacks from
    frameworks like `click` that can make your traceback long and hard
    to follow.  It also makes evrerything more colorful and easier to
    follow.
    """
    import click
    import pluggy
    import typer
    from rich import pretty as _pretty
    from rich import traceback

    _pretty.install()
    traceback.install(
        show_locals=True,
        suppress=[
            pluggy,
            click,
            typer,
        ],
    )


def parse_set_options(set_args: List[str]) -> Dict[str, Any]:
    """Parse --set key=value arguments into a nested config dict.
    
    Supports dot notation for nested keys:
    - output_dir=dist -> {"output_dir": "dist"}
    - style.theme=nord -> {"style": {"theme": "nord"}}
    """
    config = {}
    for arg in set_args:
        if "=" not in arg:
            raise ValueError(f"Invalid --set format: {arg}. Expected key=value")

        key, value = arg.split("=", 1)
        keys = key.split(".")

        # Navigate/create nested dict structure
        current = config
        for k in keys[:-1]:
            if k not in current:
                current[k] = {}
            current = current[k]

        # Set the value, attempting type conversion
        final_key = keys[-1]
        # Try to parse as JSON for complex types
        try:
            import json
            current[final_key] = json.loads(value)
        except (json.JSONDecodeError, ValueError):
            # Keep as string if not valid JSON
            current[final_key] = value

    return config


def _deep_merge(target: Dict, source: Dict) -> None:
    """Deep merge source dict into target dict."""
    for key, value in source.items():
        if key in target and isinstance(target[key], dict) and isinstance(value, dict):
            _deep_merge(target[key], value)
        else:
            target[key] = value


@hook_impl()
def cli(app: typer.Typer, markata: "Markata") -> None:
    """
    Markata hook to implement base cli commands.
    """

    plugins_app = typer.Typer()
    config_app = typer.Typer()
    app.add_typer(plugins_app, name="plugins")
    app.add_typer(config_app, name="config")

    @app.command()
    def tui(ctx: typer.Context) -> None:
        try:
            from trogon import Trogon
            from typer.main import get_group
        except ImportError:
            typer.echo("trogon not installed")
            typer.echo(
                "install markata with optional tui group to use tui `pip install 'markata[tui]'`"
            )
            return

        Trogon(get_group(app), click_context=ctx).run()

    @plugins_app.callback()
    def plugins():
        "create new things from templates"

    @config_app.callback()
    def config():
        "configuration management"

    @config_app.command()
    def show(
        verbose: bool = typer.Option(
            False,
            "--verbose",
            "-v",
        ),
    ) -> None:
        if verbose:
            markata.console.quiet = False
        else:
            markata.console.quiet = True
        rich_print(markata.config)

    @config_app.command()
    def generate(
        verbose: bool = typer.Option(
            False,
            "--verbose",
            "-v",
        ),
    ) -> None:
        if verbose:
            markata.console.quiet = False
        else:
            markata.console.quiet = True

        rich_print(toml.dumps(json.loads(markata.config.model_dump())))

    @config_app.command()
    def get(key: str) -> None:
        keys = key.split(".")
        markata.console.quiet = True
        keys_processed = ""
        value = markata.config
        na = Literal["na"]
        for key in keys:
            value = getattr(value, key, na)
            keys_processed = f"{keys_processed}.{key}".strip(".")
            if value is na:
                rich_print(f"{keys_processed} not found")
                sys.exit(1)

        rich_print(value)

    new_app = typer.Typer()
    app.add_typer(new_app, name="new")

    @new_app.callback()
    def new():
        "create new things from templates"

    @new_app.command()
    def blog(
        directory: Path = typer.Argument(
            ...,
            help="The directory to create the blog in.",
        ),
    ) -> None:
        """
        Create a new blog from using the template from
        https://github.com/WaylonWalker/markata-blog-starter.
        """

        from copier import run_copy

        typer.echo(f"creating a new project in {directory.absolute()}")
        url = markata.config.get("starters", {}).get(
            "blog",
            "git+https://github.com/WaylonWalker/markata-blog-starter",
        )
        run_copy(url, directory)

    @new_app.command()
    def post() -> None:
        """
        Create new blog post in the pages directory from the template at
        https://github.com/WaylonWalker/markata-post-template.
        """

        print("create a new post")
        from copier import run_copy

        typer.echo(f"creating a new post in {Path().absolute()}/posts")
        url = markata.config.get("starters", {}).get(
            "post",
            "git+https://github.com/WaylonWalker/markata-post-template",
        )
        run_copy(url, Path("."))

    @new_app.command()
    def plugin() -> None:
        """
        Create a new plugin using the template at
        https://github.com/WaylonWalker/markata-plugin-template.
        """
        from copier import run_copy

        typer.echo(
            f"creating a new plugin in {Path().absolute()}"
            f"/<python-package-name>/plugins",
        )
        url = markata.config.get("starters", {}).get(
            "post",
            "git+https://github.com/WaylonWalker/markata-plugin-template",
        )
        run_copy(url, Path("."))

    @app.command()
    def build(
        pretty: bool = True,
        quiet: bool = typer.Option(
            False,
            "--quiet",
            "-q",
        ),
        verbose: bool = typer.Option(
            False,
            "--verbose",
            "-v",
        ),
        should_pdb: bool = typer.Option(
            False,
            "--pdb",
        ),
        profile: bool = True,
        config_file: Optional[Path] = typer.Option(
            None,
            "-c",
            "--config",
            help="Path to alternate config file",
        ),
        output_dir: Optional[str] = typer.Option(
            None,
            "-o",
            "--output-dir",
            help="Override output directory",
        ),
        set_config: List[str] = typer.Option(
            [],
            "-s",
            "--set",
            help="Set config value (key=value, supports dot notation)",
        ),
    ) -> None:
        """
        Markata's primary way of building your site for production.
        By default, running `markta build` will render your markdown to
        the `./markout` directory.

        ``` bash
        markata build
        ```

        ## Configuration Overrides

        Override configuration at runtime using multiple methods:

        ### Alternate Config File
        Use a different config file with `-c` or `--config`:
        ``` bash
        markata build -c themes/catppuccin.toml
        ```

        ### Output Directory
        Override the output directory with `-o` or `--output-dir`:
        ``` bash
        markata build -o dist/theme-everforest
        ```

        ### Generic Config Override
        Set any config value using `-s` or `--set` with dot notation:
        ``` bash
        # Single value
        markata build -s output_dir=dist

        # Nested config
        markata build -s style.theme=nord

        # Multiple values
        markata build -s output_dir=dist -s style.theme=catppuccin

        # Complex values (use JSON)
        markata build -s 'nav={"home":"/","docs":"/docs"}'
        ```

        ### Environment Variables
        Override any config with environment variables:
        ``` bash
        # Simple value
        MARKATA_OUTPUT_DIR=dist markata build

        # Nested value (use double underscore)
        MARKATA_STYLE__THEME=nord markata build

        # Multiple values
        MARKATA_OUTPUT_DIR=dist MARKATA_STYLE__THEME=gruvbox markata build
        ```

        ### Combining Overrides
        All override methods can be combined (applied in order: file -> env -> cli):
        ``` bash
        MARKATA_STYLE__THEME=nord markata build -c base.toml -s output_dir=custom
        ```

        ## Debugging

        If you are having an issue and want to pop immediately into a debugger
        upon failure you can pass the `--pdb` flag to the build command.

        ``` bash
        markata build  --pdb
        ```

        If you do not like the way rich looks, or its suppressing tracebaks you
        would like to remain visible you can use `--no-pretty`

        ``` bash
        markata build --no-pretty
        ```

        If you need to run without any console logging pass in the
        `--quiet` flag.

        ``` bash
        markata build --quiet
        ```

        `markta build` will automatically run the pyinstrument profiler
        while building your site if you have pyinstrument installed.  It
        will echo out your profile in the console as well as write it to
        `/_profile` on your built site. If you prefer not to run
        pyinstrument profiling, even when it is installed you can pass
        in `--no-profile`

        ``` bash
        markata build --no-profile
        ```
        """

        if pretty:
            make_pretty()

        # Save console reference before potential reinit
        console = markata.console

        if quiet:
            console.quiet = True

        if verbose:
            console.print("console options:", console.options)

        # Build config overrides from CLI arguments
        config_overrides = {}

        # Add output_dir if specified
        if output_dir:
            config_overrides["output_dir"] = output_dir

        # Parse and merge --set options
        if set_config:
            set_overrides = parse_set_options(set_config)
            # Deep merge set_overrides into config_overrides
            _deep_merge(config_overrides, set_overrides)

        # Reinitialize markata with overrides if any were provided
        if config_file or config_overrides:
            from markata import Markata

            # Create a new instance with overrides
            markata_instance = Markata(
                console=console,
                config_file=config_file,
                config_overrides=config_overrides,
            )
        else:
            # Use the existing instance
            markata_instance = markata

        if not profile:
            markata_instance.config.profiler.should_profile = False

        if should_pdb:
            pdb_run(markata_instance.run)

        else:
            markata_instance.console.log("[purple]starting the build")
            markata_instance.run()

    @app.command()
    def list(
        map: str = "title",
        filter: str = "True",
        sort: str = "True",
        head: Optional[int] = None,
        tail: Optional[int] = None,
        include_empty: bool = False,
        reverse: bool = False,
        use_pager: bool = typer.Option(True, "--pager", "--no-pager"),
    ) -> None:
        """
        Provides a way run markatas, map, filter, and sort from the
        command line.  I personally use this more often than the build
        command while I am writing on a site with a large number of
        posts on it.  It makes slicing in by `templatekey`, `tag`, or
        `date` much easier.

        # default list

        By default `markata list` will list all titles in a pager, for all posts
        being loaded by markata.

        ``` bash
        markata list
        ```

        # Skip the pager

        Markata uses rich for its pager, it's pretty smart about when to
        use the pager or pass text to the next thing in the pipeline,
        but if you don't want to run a pager you can pass  `--no-pager`

        ``` bash
        markata list --no-pager
        ```

        # List other attributes

        You can list any other attribute tied to your posts.  These are
        added through either your yaml frontmatter at the start of your
        post, or through the use of a plugin.


        ``` bash
        # the filepath of the post
        markata list --map path

        # the slug of the post (where it will show up on the site)
        markata list --map slug

        # the date of the post
        markata list --map date

        # the full raw content of the post
        markata list --map content
        ```

        # List more than one attribute

        You can create new attributes as you map to echo out by
        combining existing attributes.

        ``` bash
        markata list --map 'title + " , " + slug'
        ```

        # Using Python objects as map

        You can access attributes of each post attribute that you map
        over.  For instance on my blog, each post has a date that is a
        datetime object.  I can ask each post for its `date.year`

        ``` bash
        markata list --map date.year

        # combining this with title
        markata list --map 'str(date.year) + "," + title'
        ```

        # Filtering posts

        Posts are filtered with python syntax, you will have all
        attributes tied to your posts available to filter with.

        ``` bash
        markata list --filter "'__' not in title"
        ```

        # Filtering by dates

        If your site has dates tied to your posts you can filter by
        date.  On my blog this makes a ton of sense and is quite useful.
        On the Markata docs though it doesn't really make much sense,
        since there really isn't the idea of a post date there.

        ``` bash
        # listing today's posts
        markata list --filter "date==today"

        # listing this year's posts
        markata list --filter "date.year==today.year"
        ```

        # Full Content Search

        You can also search the full content of each post for specific
        words.
        ``` bash

        markata list --filter "'python' in content"
        ```

        # Filtering by frontmatter data

        I use a templateKey on my personal blog to determine which
        template to render the page with.  I can fitler my posts by a
        `til` (today i learned) key.

        ``` bash
        markata list --filter "templateKey=='til'"
        ```

        # Combining filters

        Filters can be combined together quite like maps can, it's all
        just python syntax.

        ``` bash
        markata list --filter "templateKey=='til' and date == today"
        ```

        # Sorting posts

        Posts can be sorted by attributes on your post, and they can
        even be reversed.

        ``` bash
        markta list --sort date
        markta list --sort date --reverse
        ```

        # Putting it all together

        The real power of all this comes when you combine them all into
        lists that work for you and your workflow.  This really makes
        working on larger projects so much easier to find things.


        # Making a fuzzy picker for your posts

        Here is a bash command to open an fzf picker for todays posts,
        then open it in your `$EDITOR`

        ``` bash
        markata list \
                --map path\
                --filter 'date==today'\
                --sort date\
                --reverse |\
                fzf --preview 'bat --color always {}' |\
                xargs -I {} $EDITOR {}
        ```

        # Combining wtih nvim Telescope

        Here is the same command setup as a Telescope picker for neovim.

        ``` vim
        nnoremap <leader>et <cmd>Telescope find_files find_command=markata,list,--map,path,--filter,date==today<cr>
        ```

        If you have another way to open posts in your editor with
        `markata list` I would love to accept a PR to add it to the
        examples here.
        """

        markata.console.quiet = True

        tail = -tail if tail else tail
        filtered = markata.map(map, filter, sort)
        if not include_empty:
            filtered = [a for a in filtered if a != ""]
        filtered = filtered[tail:head]
        if reverse:
            filtered = reversed(filtered)

        class Posts(pydantic.RootModel):
            root: List[markata.Post]

        markata.console.quiet = False
        if markata.console.is_terminal and use_pager:
            with markata.console.pager():
                for a in filtered:
                    markata.console.print(a, style="purple")
        else:
            for a in filtered:
                markata.console.print(a)

    @app.command()
    def clean(
        quiet: bool = typer.Option(
            False,
            "--quiet",
            "-q",
        ),
        dry_run: bool = typer.Option(
            False,
            "--dry-run",
        ),
    ):
        """
        Cleans up output generated by markata including both the output_dir and
        the .markata_cache.

        # Dry Run

        You can run with `--dry-run` to see what markata is about to do.

        ``` bash
        markata clean --dry-run
        [09:42:37] [DRYRUN] removing outptut directory: markout base_cli.py:371
                   [DRYRUN] removing cache directory: .markata.cache base_cli.py:377

        ```

        # Running clean

        Running markata clean will fully delete all of the directories created
        by markata.

        ``` bash
        markata clean
        [09:53:04]  removing outptut directory: markout base_cli.py:394
                    removing cache directory: .markata.cache base_cli.py:405
        ```

        # Running Quietly

        Running with `--quiet` will remove all of the directories created by
        markata without announcing what it is doing.

        ``` bash
        markata clean --quiet
        ```
        """
        _clean(markata=markata, quiet=quiet, dry_run=dry_run)


def _clean(markata, quiet: bool = False, dry_run: bool = False):
    if quiet:
        markata.console.quiet = True

    markata.console.log(
        f"{'[DRYRUN]' if dry_run else ''}"
        f"removing outptut directory: {markata.config.output_dir}",
    )
    if not dry_run:
        try:
            shutil.rmtree(str(markata.config.output_dir))
        except FileNotFoundError:
            warnings.warn(
                f"output directory: {markata.config.output_dir} does not exist",
            )

    markata.console.log(
        f"{'[DRYRUN]' if dry_run else ''} removing cache directory: .markata.cache",
    )
    if not dry_run:
        try:
            shutil.rmtree(".markata.cache")
        except FileNotFoundError:
            warnings.warn("cache directory: .markata.cache does not exist")


def pdb_run(func: Callable) -> None:
    """
    Wraps a function call with a post_mortem pdb debugger.
    """
    try:
        func()
    except Exception:
        extype, value, tb = sys.exc_info()
        traceback.print_exc()
        pdb.post_mortem(tb)
