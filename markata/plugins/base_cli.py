"""
Markata's base command line commands.

This plugin enables
[`build`](https://markata.dev/markata/plugins/base_cli/#build-function)
and
[`list`](https://markata.dev/markata/plugins/base_cli/#list-function)
commands as part of the main markata cli.

# Building Your Site with the Cli

Your Markata Site can be build completely from the command line.

``` bash
markata build

# or if you prefer pipx
pipx run markata build
```

see the
[`build`](https://markata.dev/markata/plugins/base_cli/#build-function)
section for more examples.

# Listing your articles

Markata list is a tool to help list out artile attributes right to your
terminal.  This is very helpful to find articles on larger sites, or
debug what is getting picked up by markata.

``` bash
markata list --map 'str(date.year) + "," + title'
```

see the
[`list`](https://markata.dev/markata/plugins/base_cli/#list-function)
section for more examples.

# Creating "new" things with the cli

The `new` cli is built on copier templates, and allows you to build a new blog
from a starter repo, make new posts, and new plugins.  Before you start dumping
new things onto your site for the first time, make sure you have a clean git
history fully backed up, or look at the template repos to fully understand
them.

``` bash
# create a new blog template
# copier requires you to specify a directory
markata new blog [directory]

# create a new blog post
markata new post

# create a new plugin
markata new plugin

# for the most up to date help, just ask for help.
markata new --help

Usage: markata new [OPTIONS] COMMAND [ARGS]...

create new things from templates

╭─ Options ────────────────────────────────────────────────────────────────────────╮
│ --help          Show this message and exit.                                      │
╰──────────────────────────────────────────────────────────────────────────────────╯
╭─ Commands ───────────────────────────────────────────────────────────────────────╮
│ blog    Create a new blog from using the template from                           │
│         https://github.com/WaylonWalker/markata-blog-starter.                    │
│ plugin  Create a new plugin using the template at                                │
│         https://github.com/WaylonWalker/markata-plugin-template.                 │
│ post    Create new blog post in the pages directory from the template at         │
│         https://github.com/WaylonWalker/markata-post-template.                   │
╰──────────────────────────────────────────────────────────────────────────────────╯


```

"""
from pathlib import Path
import pdb
import shutil
import sys
import traceback
import toml
import json
from typing import Callable, Literal, Optional, TYPE_CHECKING
import warnings

from rich import print as rich_print
import typer

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


@hook_impl()
def cli(app: typer.Typer, markata: "Markata") -> None:
    """
    Markata hook to implement base cli commands.
    """

    plugins_app = typer.Typer()
    config_app = typer.Typer()
    app.add_typer(plugins_app)
    app.add_typer(config_app)

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

        rich_print(toml.dumps(json.loads(markata.config.json())))

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
    app.add_typer(new_app)

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
    ) -> None:
        """
        Markata's primary way of building your site for production.
        By default, running `markta build` will render your markdown to
        the `./markout` directory.

        ``` bash
        markata build
        ```

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

        if quiet:
            markata.console.quiet = True

        if verbose:
            markata.console.print("console options:", markata.console.options)

        if profile:
            markata.should_profile_cli = True
            markata.should_profile = True
            markata.configure()

        if should_pdb:
            pdb_run(markata.run)

        else:
            markata.run()

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
        f'{"[DRYRUN]" if dry_run else ""}'
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
        f'{"[DRYRUN]" if dry_run else ""} removing cache directory: .markata.cache',
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
