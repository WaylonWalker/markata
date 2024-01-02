---
content: "Markata's base command line commands.\n\nThis plugin enables\n[`build`](https://markata.dev/markata/plugins/base_cli/#build-function)\nand\n[`list`](https://markata.dev/markata/plugins/base_cli/#list-function)\ncommands
  as part of the main markata cli.\n\n# Building Your Site with the Cli\n\nYour Markata
  Site can be build completely from the command line.\n\n``` bash\nmarkata build\n\n#
  or if you prefer pipx\npipx run markata build\n```\n\nsee the\n[`build`](https://markata.dev/markata/plugins/base_cli/#build-function)\nsection
  for more examples.\n\n# Listing your articles\n\nMarkata list is a tool to help
  list out artile attributes right to your\nterminal.  This is very helpful to find
  articles on larger sites, or\ndebug what is getting picked up by markata.\n\n```
  bash\nmarkata list --map 'str(date.year) + \",\" + title'\n```\n\nsee the\n[`list`](https://markata.dev/markata/plugins/base_cli/#list-function)\nsection
  for more examples.\n\n# Creating \"new\" things with the cli\n\nThe `new` cli is
  built on copier templates, and allows you to build a new blog\nfrom a starter repo,
  make new posts, and new plugins.  Before you start dumping\nnew things onto your
  site for the first time, make sure you have a clean git\nhistory fully backed up,
  or look at the template repos to fully understand\nthem.\n\n``` bash\n# create a
  new blog template\n# copier requires you to specify a directory\nmarkata new blog
  [directory]\n\n# create a new blog post\nmarkata new post\n\n# create a new plugin\nmarkata
  new plugin\n\n# for the most up to date help, just ask for help.\nmarkata new --help\n\nUsage:
  markata new [OPTIONS] COMMAND [ARGS]...\n\ncreate new things from templates\n\n\u256D\u2500
  Options \u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u256E\n\u2502
  --help          Show this message and exit.                                      \u2502\n\u2570\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u256F\n\u256D\u2500
  Commands \u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u256E\n\u2502
  blog    Create a new blog from using the template from                           \u2502\n\u2502
  \        https://github.com/WaylonWalker/markata-blog-starter.                    \u2502\n\u2502
  plugin  Create a new plugin using the template at                                \u2502\n\u2502
  \        https://github.com/WaylonWalker/markata-plugin-template.                 \u2502\n\u2502
  post    Create new blog post in the pages directory from the template at         \u2502\n\u2502
  \        https://github.com/WaylonWalker/markata-post-template.                   \u2502\n\u2570\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u256F\n\n\n```\n\n\n!!
  function <h2 id='make_pretty' class='admonition-title' style='margin:0;padding:.5rem
  1rem;'>make_pretty <em class='small'>function</em></h2>\n    This is a helper function
  that enables suppresses tracebacks from\n    frameworks like `click` that can make
  your traceback long and hard\n    to follow.  It also makes evrerything more colorful
  and easier to\n    follow.\n???+ source \"make_pretty <em class='small'>source</em>\"\n\n```python\n\n
  \       def make_pretty() -> None:\n            \"\"\"\n            This is a helper
  function that enables suppresses tracebacks from\n            frameworks like `click`
  that can make your traceback long and hard\n            to follow.  It also makes
  evrerything more colorful and easier to\n            follow.\n            \"\"\"\n
  \           import click\n            import pluggy\n            import typer\n
  \           from rich import pretty as _pretty\n            from rich import traceback\n\n
  \           _pretty.install()\n            traceback.install(\n                show_locals=True,\n
  \               suppress=[\n                    pluggy,\n                    click,\n
  \                   typer,\n                ],\n            )\n```\n\n\n!! function
  <h2 id='cli' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>cli <em
  class='small'>function</em></h2>\n    Markata hook to implement base cli commands.\n???+
  source \"cli <em class='small'>source</em>\"\n\n```python\n\n        def cli(app:
  typer.Typer, markata: \"Markata\") -> None:\n            \"\"\"\n            Markata
  hook to implement base cli commands.\n            \"\"\"\n\n            plugins_app
  = typer.Typer()\n            config_app = typer.Typer()\n            app.add_typer(plugins_app)\n
  \           app.add_typer(config_app)\n\n            @app.command()\n            def
  tui(ctx: typer.Context) -> None:\n                try:\n                    from
  trogon import Trogon\n                    from typer.main import get_group\n                except
  ImportError:\n                    typer.echo(\"trogon not installed\")\n                    typer.echo(\n
  \                       \"install markata with optional tui group to use tui `pip
  install 'markata[tui]'`\"\n                    )\n                    return\n\n
  \               Trogon(get_group(app), click_context=ctx).run()\n\n            @plugins_app.callback()\n
  \           def plugins():\n                \"create new things from templates\"\n\n
  \           @config_app.callback()\n            def config():\n                \"configuration
  management\"\n\n            @config_app.command()\n            def show(\n                verbose:
  bool = typer.Option(\n                    False,\n                    \"--verbose\",\n
  \                   \"-v\",\n                ),\n            ) -> None:\n                if
  verbose:\n                    markata.console.quiet = False\n                else:\n
  \                   markata.console.quiet = True\n                rich_print(markata.config)\n\n
  \           @config_app.command()\n            def generate(\n                verbose:
  bool = typer.Option(\n                    False,\n                    \"--verbose\",\n
  \                   \"-v\",\n                ),\n            ) -> None:\n                if
  verbose:\n                    markata.console.quiet = False\n                else:\n
  \                   markata.console.quiet = True\n\n                rich_print(toml.dumps(json.loads(markata.config.json())))\n\n
  \           @config_app.command()\n            def get(key: str) -> None:\n                keys
  = key.split(\".\")\n                markata.console.quiet = True\n                keys_processed
  = \"\"\n                value = markata.config\n                na = Literal[\"na\"]\n
  \               for key in keys:\n                    value = getattr(value, key,
  na)\n                    keys_processed = f\"{keys_processed}.{key}\".strip(\".\")\n
  \                   if value is na:\n                        rich_print(f\"{keys_processed}
  not found\")\n                        sys.exit(1)\n\n                rich_print(value)\n\n
  \           new_app = typer.Typer()\n            app.add_typer(new_app)\n\n            @new_app.callback()\n
  \           def new():\n                \"create new things from templates\"\n\n
  \           @new_app.command()\n            def blog(\n                directory:
  Path = typer.Argument(\n                    ...,\n                    help=\"The
  directory to create the blog in.\",\n                ),\n            ) -> None:\n
  \               \"\"\"\n                Create a new blog from using the template
  from\n                https://github.com/WaylonWalker/markata-blog-starter.\n                \"\"\"\n\n
  \               from copier import run_copy\n\n                typer.echo(f\"creating
  a new project in {directory.absolute()}\")\n                url = markata.config.get(\"starters\",
  {}).get(\n                    \"blog\",\n                    \"git+https://github.com/WaylonWalker/markata-blog-starter\",\n
  \               )\n                run_copy(url, directory)\n\n            @new_app.command()\n
  \           def post() -> None:\n                \"\"\"\n                Create
  new blog post in the pages directory from the template at\n                https://github.com/WaylonWalker/markata-post-template.\n
  \               \"\"\"\n\n                print(\"create a new post\")\n                from
  copier import run_copy\n\n                typer.echo(f\"creating a new post in {Path().absolute()}/posts\")\n
  \               url = markata.config.get(\"starters\", {}).get(\n                    \"post\",\n
  \                   \"git+https://github.com/WaylonWalker/markata-post-template\",\n
  \               )\n                run_copy(url, Path(\".\"))\n\n            @new_app.command()\n
  \           def plugin() -> None:\n                \"\"\"\n                Create
  a new plugin using the template at\n                https://github.com/WaylonWalker/markata-plugin-template.\n
  \               \"\"\"\n                from copier import run_copy\n\n                typer.echo(\n
  \                   f\"creating a new plugin in {Path().absolute()}\"\n                    f\"/<python-package-name>/plugins\",\n
  \               )\n                url = markata.config.get(\"starters\", {}).get(\n
  \                   \"post\",\n                    \"git+https://github.com/WaylonWalker/markata-plugin-template\",\n
  \               )\n                run_copy(url, Path(\".\"))\n\n            @app.command()\n
  \           def build(\n                pretty: bool = True,\n                quiet:
  bool = typer.Option(\n                    False,\n                    \"--quiet\",\n
  \                   \"-q\",\n                ),\n                verbose: bool =
  typer.Option(\n                    False,\n                    \"--verbose\",\n
  \                   \"-v\",\n                ),\n                should_pdb: bool
  = typer.Option(\n                    False,\n                    \"--pdb\",\n                ),\n
  \               profile: bool = True,\n            ) -> None:\n                \"\"\"\n
  \               Markata's primary way of building your site for production.\n                By
  default, running `markta build` will render your markdown to\n                the
  `./markout` directory.\n\n                ``` bash\n                markata build\n
  \               ```\n\n                If you are having an issue and want to pop
  immediately into a debugger\n                upon failure you can pass the `--pdb`
  flag to the build command.\n\n                ``` bash\n                markata
  build  --pdb\n                ```\n\n                If you do not like the way
  rich looks, or its suppressing tracebaks you\n                would like to remain
  visible you can use `--no-pretty`\n\n                ``` bash\n                markata
  build --no-pretty\n                ```\n\n                If you need to run without
  any console logging pass in the\n                `--quiet` flag.\n\n                ```
  bash\n                markata build --quiet\n                ```\n\n                `markta
  build` will automatically run the pyinstrument profiler\n                while building
  your site if you have pyinstrument installed.  It\n                will echo out
  your profile in the console as well as write it to\n                `/_profile`
  on your built site. If you prefer not to run\n                pyinstrument profiling,
  even when it is installed you can pass\n                in `--no-profile`\n\n                ```
  bash\n                markata build --no-profile\n                ```\n                \"\"\"\n\n
  \               if pretty:\n                    make_pretty()\n\n                if
  quiet:\n                    markata.console.quiet = True\n\n                if verbose:\n
  \                   markata.console.print(\"console options:\", markata.console.options)\n\n
  \               if profile:\n                    markata.should_profile_cli = True\n
  \                   markata.should_profile = True\n                    markata.configure()\n\n
  \               if should_pdb:\n                    pdb_run(markata.run)\n\n                else:\n
  \                   markata.run()\n\n            @app.command()\n            def
  list(\n                map: str = \"title\",\n                filter: str = \"True\",\n
  \               sort: str = \"True\",\n                head: Optional[int] = None,\n
  \               tail: Optional[int] = None,\n                include_empty: bool
  = False,\n                reverse: bool = False,\n                use_pager: bool
  = typer.Option(True, \"--pager\", \"--no-pager\"),\n            ) -> None:\n                \"\"\"\n
  \               Provides a way run markatas, map, filter, and sort from the\n                command
  line.  I personally use this more often than the build\n                command
  while I am writing on a site with a large number of\n                posts on it.
  \ It makes slicing in by `templatekey`, `tag`, or\n                `date` much easier.\n\n
  \               # default list\n\n                By default `markata list` will
  list all titles in a pager, for all posts\n                being loaded by markata.\n\n
  \               ``` bash\n                markata list\n                ```\n\n
  \               # Skip the pager\n\n                Markata uses rich for its pager,
  it's pretty smart about when to\n                use the pager or pass text to the
  next thing in the pipeline,\n                but if you don't want to run a pager
  you can pass  `--no-pager`\n\n                ``` bash\n                markata
  list --no-pager\n                ```\n\n                # List other attributes\n\n
  \               You can list any other attribute tied to your posts.  These are\n
  \               added through either your yaml frontmatter at the start of your\n
  \               post, or through the use of a plugin.\n\n\n                ``` bash\n
  \               # the filepath of the post\n                markata list --map path\n\n
  \               # the slug of the post (where it will show up on the site)\n                markata
  list --map slug\n\n                # the date of the post\n                markata
  list --map date\n\n                # the full raw content of the post\n                markata
  list --map content\n                ```\n\n                # List more than one
  attribute\n\n                You can create new attributes as you map to echo out
  by\n                combining existing attributes.\n\n                ``` bash\n
  \               markata list --map 'title + \" , \" + slug'\n                ```\n\n
  \               # Using Python objects as map\n\n                You can access
  attributes of each post attribute that you map\n                over.  For instance
  on my blog, each post has a date that is a\n                datetime object.  I
  can ask each post for its `date.year`\n\n                ``` bash\n                markata
  list --map date.year\n\n                # combining this with title\n                markata
  list --map 'str(date.year) + \",\" + title'\n                ```\n\n                #
  Filtering posts\n\n                Posts are filtered with python syntax, you will
  have all\n                attributes tied to your posts available to filter with.\n\n
  \               ``` bash\n                markata list --filter \"'__' not in title\"\n
  \               ```\n\n                # Filtering by dates\n\n                If
  your site has dates tied to your posts you can filter by\n                date.
  \ On my blog this makes a ton of sense and is quite useful.\n                On
  the Markata docs though it doesn't really make much sense,\n                since
  there really isn't the idea of a post date there.\n\n                ``` bash\n
  \               # listing today's posts\n                markata list --filter \"date==today\"\n\n
  \               # listing this year's posts\n                markata list --filter
  \"date.year==today.year\"\n                ```\n\n                # Full Content
  Search\n\n                You can also search the full content of each post for
  specific\n                words.\n                ``` bash\n\n                markata
  list --filter \"'python' in content\"\n                ```\n\n                #
  Filtering by frontmatter data\n\n                I use a templateKey on my personal
  blog to determine which\n                template to render the page with.  I can
  fitler my posts by a\n                `til` (today i learned) key.\n\n                ```
  bash\n                markata list --filter \"templateKey=='til'\"\n                ```\n\n
  \               # Combining filters\n\n                Filters can be combined together
  quite like maps can, it's all\n                just python syntax.\n\n                ```
  bash\n                markata list --filter \"templateKey=='til' and date == today\"\n
  \               ```\n\n                # Sorting posts\n\n                Posts
  can be sorted by attributes on your post, and they can\n                even be
  reversed.\n\n                ``` bash\n                markta list --sort date\n
  \               markta list --sort date --reverse\n                ```\n\n                #
  Putting it all together\n\n                The real power of all this comes when
  you combine them all into\n                lists that work for you and your workflow.
  \ This really makes\n                working on larger projects so much easier to
  find things.\n\n\n                # Making a fuzzy picker for your posts\n\n                Here
  is a bash command to open an fzf picker for todays posts,\n                then
  open it in your `$EDITOR`\n\n                ``` bash\n                markata list
  \\\n                        --map path\\\n                        --filter 'date==today'\\\n
  \                       --sort date\\\n                        --reverse |\\\n                        fzf
  --preview 'bat --color always {}' |\\\n                        xargs -I {} $EDITOR
  {}\n                ```\n\n                # Combining wtih nvim Telescope\n\n                Here
  is the same command setup as a Telescope picker for neovim.\n\n                ```
  vim\n                nnoremap <leader>et <cmd>Telescope find_files find_command=markata,list,--map,path,--filter,date==today<cr>\n
  \               ```\n\n                If you have another way to open posts in
  your editor with\n                `markata list` I would love to accept a PR to
  add it to the\n                examples here.\n                \"\"\"\n\n                markata.console.quiet
  = True\n\n                tail = -tail if tail else tail\n                filtered
  = markata.map(map, filter, sort)\n                if not include_empty:\n                    filtered
  = [a for a in filtered if a != \"\"]\n                filtered = filtered[tail:head]\n
  \               if reverse:\n                    filtered = reversed(filtered)\n\n
  \               markata.console.quiet = False\n                if markata.console.is_terminal
  and use_pager:\n                    with markata.console.pager():\n                        for
  a in filtered:\n                            markata.console.print(a, style=\"purple\")\n
  \               else:\n                    for a in filtered:\n                        markata.console.print(a)\n\n
  \           @app.command()\n            def clean(\n                quiet: bool
  = typer.Option(\n                    False,\n                    \"--quiet\",\n
  \                   \"-q\",\n                ),\n                dry_run: bool =
  typer.Option(\n                    False,\n                    \"--dry-run\",\n
  \               ),\n            ):\n                \"\"\"\n                Cleans
  up output generated by markata including both the output_dir and\n                the
  .markata_cache.\n\n                # Dry Run\n\n                You can run with
  `--dry-run` to see what markata is about to do.\n\n                ``` bash\n                markata
  clean --dry-run\n                [09:42:37] [DRYRUN] removing outptut directory:
  markout base_cli.py:371\n                           [DRYRUN] removing cache directory:
  .markata.cache base_cli.py:377\n\n                ```\n\n                # Running
  clean\n\n                Running markata clean will fully delete all of the directories
  created\n                by markata.\n\n                ``` bash\n                markata
  clean\n                [09:53:04]  removing outptut directory: markout base_cli.py:394\n
  \                           removing cache directory: .markata.cache base_cli.py:405\n
  \               ```\n\n                # Running Quietly\n\n                Running
  with `--quiet` will remove all of the directories created by\n                markata
  without announcing what it is doing.\n\n                ``` bash\n                markata
  clean --quiet\n                ```\n                \"\"\"\n                _clean(markata=markata,
  quiet=quiet, dry_run=dry_run)\n```\n\n\n!! function <h2 id='_clean' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>_clean <em class='small'>function</em></h2>\n\n???+
  source \"_clean <em class='small'>source</em>\"\n\n```python\n\n        def _clean(markata,
  quiet: bool = False, dry_run: bool = False):\n            if quiet:\n                markata.console.quiet
  = True\n\n            markata.console.log(\n                f'{\"[DRYRUN]\" if dry_run
  else \"\"}'\n                f\"removing outptut directory: {markata.config.output_dir}\",\n
  \           )\n            if not dry_run:\n                try:\n                    shutil.rmtree(str(markata.config.output_dir))\n
  \               except FileNotFoundError:\n                    warnings.warn(\n
  \                       f\"output directory: {markata.config.output_dir} does not
  exist\",\n                    )\n\n            markata.console.log(\n                f'{\"[DRYRUN]\"
  if dry_run else \"\"} removing cache directory: .markata.cache',\n            )\n
  \           if not dry_run:\n                try:\n                    shutil.rmtree(\".markata.cache\")\n
  \               except FileNotFoundError:\n                    warnings.warn(\"cache
  directory: .markata.cache does not exist\")\n```\n\n\n!! function <h2 id='pdb_run'
  class='admonition-title' style='margin:0;padding:.5rem 1rem;'>pdb_run <em class='small'>function</em></h2>\n
  \   Wraps a function call with a post_mortem pdb debugger.\n???+ source \"pdb_run
  <em class='small'>source</em>\"\n\n```python\n\n        def pdb_run(func: Callable)
  -> None:\n            \"\"\"\n            Wraps a function call with a post_mortem
  pdb debugger.\n            \"\"\"\n            try:\n                func()\n            except
  Exception:\n                extype, value, tb = sys.exc_info()\n                traceback.print_exc()\n
  \               pdb.post_mortem(tb)\n```\n\n\n!! function <h2 id='tui' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>tui <em class='small'>function</em></h2>\n\n???+
  source \"tui <em class='small'>source</em>\"\n\n```python\n\n        def tui(ctx:
  typer.Context) -> None:\n                try:\n                    from trogon import
  Trogon\n                    from typer.main import get_group\n                except
  ImportError:\n                    typer.echo(\"trogon not installed\")\n                    typer.echo(\n
  \                       \"install markata with optional tui group to use tui `pip
  install 'markata[tui]'`\"\n                    )\n                    return\n\n
  \               Trogon(get_group(app), click_context=ctx).run()\n```\n\n\n!! function
  <h2 id='plugins' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>plugins
  <em class='small'>function</em></h2>\n    create new things from templates\n???+
  source \"plugins <em class='small'>source</em>\"\n\n```python\n\n        def plugins():\n
  \               \"create new things from templates\"\n```\n\n\n!! function <h2 id='config'
  class='admonition-title' style='margin:0;padding:.5rem 1rem;'>config <em class='small'>function</em></h2>\n
  \   configuration management\n???+ source \"config <em class='small'>source</em>\"\n\n```python\n\n
  \       def config():\n                \"configuration management\"\n```\n\n\n!!
  function <h2 id='show' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>show
  <em class='small'>function</em></h2>\n\n???+ source \"show <em class='small'>source</em>\"\n\n```python\n\n
  \       def show(\n                verbose: bool = typer.Option(\n                    False,\n
  \                   \"--verbose\",\n                    \"-v\",\n                ),\n
  \           ) -> None:\n                if verbose:\n                    markata.console.quiet
  = False\n                else:\n                    markata.console.quiet = True\n
  \               rich_print(markata.config)\n```\n\n\n!! function <h2 id='generate'
  class='admonition-title' style='margin:0;padding:.5rem 1rem;'>generate <em class='small'>function</em></h2>\n\n???+
  source \"generate <em class='small'>source</em>\"\n\n```python\n\n        def generate(\n
  \               verbose: bool = typer.Option(\n                    False,\n                    \"--verbose\",\n
  \                   \"-v\",\n                ),\n            ) -> None:\n                if
  verbose:\n                    markata.console.quiet = False\n                else:\n
  \                   markata.console.quiet = True\n\n                rich_print(toml.dumps(json.loads(markata.config.json())))\n```\n\n\n!!
  function <h2 id='get' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>get
  <em class='small'>function</em></h2>\n\n???+ source \"get <em class='small'>source</em>\"\n\n```python\n\n
  \       def get(key: str) -> None:\n                keys = key.split(\".\")\n                markata.console.quiet
  = True\n                keys_processed = \"\"\n                value = markata.config\n
  \               na = Literal[\"na\"]\n                for key in keys:\n                    value
  = getattr(value, key, na)\n                    keys_processed = f\"{keys_processed}.{key}\".strip(\".\")\n
  \                   if value is na:\n                        rich_print(f\"{keys_processed}
  not found\")\n                        sys.exit(1)\n\n                rich_print(value)\n```\n\n\n!!
  function <h2 id='new' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>new
  <em class='small'>function</em></h2>\n    create new things from templates\n???+
  source \"new <em class='small'>source</em>\"\n\n```python\n\n        def new():\n
  \               \"create new things from templates\"\n```\n\n\n!! function <h2 id='blog'
  class='admonition-title' style='margin:0;padding:.5rem 1rem;'>blog <em class='small'>function</em></h2>\n
  \   Create a new blog from using the template from\n    https://github.com/WaylonWalker/markata-blog-starter.\n???+
  source \"blog <em class='small'>source</em>\"\n\n```python\n\n        def blog(\n
  \               directory: Path = typer.Argument(\n                    ...,\n                    help=\"The
  directory to create the blog in.\",\n                ),\n            ) -> None:\n
  \               \"\"\"\n                Create a new blog from using the template
  from\n                https://github.com/WaylonWalker/markata-blog-starter.\n                \"\"\"\n\n
  \               from copier import run_copy\n\n                typer.echo(f\"creating
  a new project in {directory.absolute()}\")\n                url = markata.config.get(\"starters\",
  {}).get(\n                    \"blog\",\n                    \"git+https://github.com/WaylonWalker/markata-blog-starter\",\n
  \               )\n                run_copy(url, directory)\n```\n\n\n!! function
  <h2 id='post' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>post
  <em class='small'>function</em></h2>\n    Create new blog post in the pages directory
  from the template at\n    https://github.com/WaylonWalker/markata-post-template.\n???+
  source \"post <em class='small'>source</em>\"\n\n```python\n\n        def post()
  -> None:\n                \"\"\"\n                Create new blog post in the pages
  directory from the template at\n                https://github.com/WaylonWalker/markata-post-template.\n
  \               \"\"\"\n\n                print(\"create a new post\")\n                from
  copier import run_copy\n\n                typer.echo(f\"creating a new post in {Path().absolute()}/posts\")\n
  \               url = markata.config.get(\"starters\", {}).get(\n                    \"post\",\n
  \                   \"git+https://github.com/WaylonWalker/markata-post-template\",\n
  \               )\n                run_copy(url, Path(\".\"))\n```\n\n\n!! function
  <h2 id='plugin' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>plugin
  <em class='small'>function</em></h2>\n    Create a new plugin using the template
  at\n    https://github.com/WaylonWalker/markata-plugin-template.\n???+ source \"plugin
  <em class='small'>source</em>\"\n\n```python\n\n        def plugin() -> None:\n
  \               \"\"\"\n                Create a new plugin using the template at\n
  \               https://github.com/WaylonWalker/markata-plugin-template.\n                \"\"\"\n
  \               from copier import run_copy\n\n                typer.echo(\n                    f\"creating
  a new plugin in {Path().absolute()}\"\n                    f\"/<python-package-name>/plugins\",\n
  \               )\n                url = markata.config.get(\"starters\", {}).get(\n
  \                   \"post\",\n                    \"git+https://github.com/WaylonWalker/markata-plugin-template\",\n
  \               )\n                run_copy(url, Path(\".\"))\n```\n\n\n!! function
  <h2 id='build' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>build
  <em class='small'>function</em></h2>\n    Markata's primary way of building your
  site for production.\n    By default, running `markta build` will render your markdown
  to\n    the `./markout` directory.\n\n    ``` bash\n    markata build\n    ```\n\n
  \   If you are having an issue and want to pop immediately into a debugger\n    upon
  failure you can pass the `--pdb` flag to the build command.\n\n    ``` bash\n    markata
  build  --pdb\n    ```\n\n    If you do not like the way rich looks, or its suppressing
  tracebaks you\n    would like to remain visible you can use `--no-pretty`\n\n    ```
  bash\n    markata build --no-pretty\n    ```\n\n    If you need to run without any
  console logging pass in the\n    `--quiet` flag.\n\n    ``` bash\n    markata build
  --quiet\n    ```\n\n    `markta build` will automatically run the pyinstrument profiler\n
  \   while building your site if you have pyinstrument installed.  It\n    will echo
  out your profile in the console as well as write it to\n    `/_profile` on your
  built site. If you prefer not to run\n    pyinstrument profiling, even when it is
  installed you can pass\n    in `--no-profile`\n\n    ``` bash\n    markata build
  --no-profile\n    ```\n???+ source \"build <em class='small'>source</em>\"\n\n```python\n\n
  \       def build(\n                pretty: bool = True,\n                quiet:
  bool = typer.Option(\n                    False,\n                    \"--quiet\",\n
  \                   \"-q\",\n                ),\n                verbose: bool =
  typer.Option(\n                    False,\n                    \"--verbose\",\n
  \                   \"-v\",\n                ),\n                should_pdb: bool
  = typer.Option(\n                    False,\n                    \"--pdb\",\n                ),\n
  \               profile: bool = True,\n            ) -> None:\n                \"\"\"\n
  \               Markata's primary way of building your site for production.\n                By
  default, running `markta build` will render your markdown to\n                the
  `./markout` directory.\n\n                ``` bash\n                markata build\n
  \               ```\n\n                If you are having an issue and want to pop
  immediately into a debugger\n                upon failure you can pass the `--pdb`
  flag to the build command.\n\n                ``` bash\n                markata
  build  --pdb\n                ```\n\n                If you do not like the way
  rich looks, or its suppressing tracebaks you\n                would like to remain
  visible you can use `--no-pretty`\n\n                ``` bash\n                markata
  build --no-pretty\n                ```\n\n                If you need to run without
  any console logging pass in the\n                `--quiet` flag.\n\n                ```
  bash\n                markata build --quiet\n                ```\n\n                `markta
  build` will automatically run the pyinstrument profiler\n                while building
  your site if you have pyinstrument installed.  It\n                will echo out
  your profile in the console as well as write it to\n                `/_profile`
  on your built site. If you prefer not to run\n                pyinstrument profiling,
  even when it is installed you can pass\n                in `--no-profile`\n\n                ```
  bash\n                markata build --no-profile\n                ```\n                \"\"\"\n\n
  \               if pretty:\n                    make_pretty()\n\n                if
  quiet:\n                    markata.console.quiet = True\n\n                if verbose:\n
  \                   markata.console.print(\"console options:\", markata.console.options)\n\n
  \               if profile:\n                    markata.should_profile_cli = True\n
  \                   markata.should_profile = True\n                    markata.configure()\n\n
  \               if should_pdb:\n                    pdb_run(markata.run)\n\n                else:\n
  \                   markata.run()\n```\n\n\n!! function <h2 id='list' class='admonition-title'
  style='margin:0;padding:.5rem 1rem;'>list <em class='small'>function</em></h2>\n
  \   Provides a way run markatas, map, filter, and sort from the\n    command line.
  \ I personally use this more often than the build\n    command while I am writing
  on a site with a large number of\n    posts on it.  It makes slicing in by `templatekey`,
  `tag`, or\n    `date` much easier.\n\n    # default list\n\n    By default `markata
  list` will list all titles in a pager, for all posts\n    being loaded by markata.\n\n
  \   ``` bash\n    markata list\n    ```\n\n    # Skip the pager\n\n    Markata uses
  rich for its pager, it's pretty smart about when to\n    use the pager or pass text
  to the next thing in the pipeline,\n    but if you don't want to run a pager you
  can pass  `--no-pager`\n\n    ``` bash\n    markata list --no-pager\n    ```\n\n
  \   # List other attributes\n\n    You can list any other attribute tied to your
  posts.  These are\n    added through either your yaml frontmatter at the start of
  your\n    post, or through the use of a plugin.\n\n\n    ``` bash\n    # the filepath
  of the post\n    markata list --map path\n\n    # the slug of the post (where it
  will show up on the site)\n    markata list --map slug\n\n    # the date of the
  post\n    markata list --map date\n\n    # the full raw content of the post\n    markata
  list --map content\n    ```\n\n    # List more than one attribute\n\n    You can
  create new attributes as you map to echo out by\n    combining existing attributes.\n\n
  \   ``` bash\n    markata list --map 'title + \" , \" + slug'\n    ```\n\n    #
  Using Python objects as map\n\n    You can access attributes of each post attribute
  that you map\n    over.  For instance on my blog, each post has a date that is a\n
  \   datetime object.  I can ask each post for its `date.year`\n\n    ``` bash\n
  \   markata list --map date.year\n\n    # combining this with title\n    markata
  list --map 'str(date.year) + \",\" + title'\n    ```\n\n    # Filtering posts\n\n
  \   Posts are filtered with python syntax, you will have all\n    attributes tied
  to your posts available to filter with.\n\n    ``` bash\n    markata list --filter
  \"'__' not in title\"\n    ```\n\n    # Filtering by dates\n\n    If your site has
  dates tied to your posts you can filter by\n    date.  On my blog this makes a ton
  of sense and is quite useful.\n    On the Markata docs though it doesn't really
  make much sense,\n    since there really isn't the idea of a post date there.\n\n
  \   ``` bash\n    # listing today's posts\n    markata list --filter \"date==today\"\n\n
  \   # listing this year's posts\n    markata list --filter \"date.year==today.year\"\n
  \   ```\n\n    # Full Content Search\n\n    You can also search the full content
  of each post for specific\n    words.\n    ``` bash\n\n    markata list --filter
  \"'python' in content\"\n    ```\n\n    # Filtering by frontmatter data\n\n    I
  use a templateKey on my personal blog to determine which\n    template to render
  the page with.  I can fitler my posts by a\n    `til` (today i learned) key.\n\n
  \   ``` bash\n    markata list --filter \"templateKey=='til'\"\n    ```\n\n    #
  Combining filters\n\n    Filters can be combined together quite like maps can, it's
  all\n    just python syntax.\n\n    ``` bash\n    markata list --filter \"templateKey=='til'
  and date == today\"\n    ```\n\n    # Sorting posts\n\n    Posts can be sorted by
  attributes on your post, and they can\n    even be reversed.\n\n    ``` bash\n    markta
  list --sort date\n    markta list --sort date --reverse\n    ```\n\n    # Putting
  it all together\n\n    The real power of all this comes when you combine them all
  into\n    lists that work for you and your workflow.  This really makes\n    working
  on larger projects so much easier to find things.\n\n\n    # Making a fuzzy picker
  for your posts\n\n    Here is a bash command to open an fzf picker for todays posts,\n
  \   then open it in your `$EDITOR`\n\n    ``` bash\n    markata list                 --map
  path                --filter 'date==today'                --sort date                --reverse
  |                fzf --preview 'bat --color always {}' |                xargs -I
  {} $EDITOR {}\n    ```\n\n    # Combining wtih nvim Telescope\n\n    Here is the
  same command setup as a Telescope picker for neovim.\n\n    ``` vim\n    nnoremap
  <leader>et <cmd>Telescope find_files find_command=markata,list,--map,path,--filter,date==today<cr>\n
  \   ```\n\n    If you have another way to open posts in your editor with\n    `markata
  list` I would love to accept a PR to add it to the\n    examples here.\n???+ source
  \"list <em class='small'>source</em>\"\n\n```python\n\n        def list(\n                map:
  str = \"title\",\n                filter: str = \"True\",\n                sort:
  str = \"True\",\n                head: Optional[int] = None,\n                tail:
  Optional[int] = None,\n                include_empty: bool = False,\n                reverse:
  bool = False,\n                use_pager: bool = typer.Option(True, \"--pager\",
  \"--no-pager\"),\n            ) -> None:\n                \"\"\"\n                Provides
  a way run markatas, map, filter, and sort from the\n                command line.
  \ I personally use this more often than the build\n                command while
  I am writing on a site with a large number of\n                posts on it.  It
  makes slicing in by `templatekey`, `tag`, or\n                `date` much easier.\n\n
  \               # default list\n\n                By default `markata list` will
  list all titles in a pager, for all posts\n                being loaded by markata.\n\n
  \               ``` bash\n                markata list\n                ```\n\n
  \               # Skip the pager\n\n                Markata uses rich for its pager,
  it's pretty smart about when to\n                use the pager or pass text to the
  next thing in the pipeline,\n                but if you don't want to run a pager
  you can pass  `--no-pager`\n\n                ``` bash\n                markata
  list --no-pager\n                ```\n\n                # List other attributes\n\n
  \               You can list any other attribute tied to your posts.  These are\n
  \               added through either your yaml frontmatter at the start of your\n
  \               post, or through the use of a plugin.\n\n\n                ``` bash\n
  \               # the filepath of the post\n                markata list --map path\n\n
  \               # the slug of the post (where it will show up on the site)\n                markata
  list --map slug\n\n                # the date of the post\n                markata
  list --map date\n\n                # the full raw content of the post\n                markata
  list --map content\n                ```\n\n                # List more than one
  attribute\n\n                You can create new attributes as you map to echo out
  by\n                combining existing attributes.\n\n                ``` bash\n
  \               markata list --map 'title + \" , \" + slug'\n                ```\n\n
  \               # Using Python objects as map\n\n                You can access
  attributes of each post attribute that you map\n                over.  For instance
  on my blog, each post has a date that is a\n                datetime object.  I
  can ask each post for its `date.year`\n\n                ``` bash\n                markata
  list --map date.year\n\n                # combining this with title\n                markata
  list --map 'str(date.year) + \",\" + title'\n                ```\n\n                #
  Filtering posts\n\n                Posts are filtered with python syntax, you will
  have all\n                attributes tied to your posts available to filter with.\n\n
  \               ``` bash\n                markata list --filter \"'__' not in title\"\n
  \               ```\n\n                # Filtering by dates\n\n                If
  your site has dates tied to your posts you can filter by\n                date.
  \ On my blog this makes a ton of sense and is quite useful.\n                On
  the Markata docs though it doesn't really make much sense,\n                since
  there really isn't the idea of a post date there.\n\n                ``` bash\n
  \               # listing today's posts\n                markata list --filter \"date==today\"\n\n
  \               # listing this year's posts\n                markata list --filter
  \"date.year==today.year\"\n                ```\n\n                # Full Content
  Search\n\n                You can also search the full content of each post for
  specific\n                words.\n                ``` bash\n\n                markata
  list --filter \"'python' in content\"\n                ```\n\n                #
  Filtering by frontmatter data\n\n                I use a templateKey on my personal
  blog to determine which\n                template to render the page with.  I can
  fitler my posts by a\n                `til` (today i learned) key.\n\n                ```
  bash\n                markata list --filter \"templateKey=='til'\"\n                ```\n\n
  \               # Combining filters\n\n                Filters can be combined together
  quite like maps can, it's all\n                just python syntax.\n\n                ```
  bash\n                markata list --filter \"templateKey=='til' and date == today\"\n
  \               ```\n\n                # Sorting posts\n\n                Posts
  can be sorted by attributes on your post, and they can\n                even be
  reversed.\n\n                ``` bash\n                markta list --sort date\n
  \               markta list --sort date --reverse\n                ```\n\n                #
  Putting it all together\n\n                The real power of all this comes when
  you combine them all into\n                lists that work for you and your workflow.
  \ This really makes\n                working on larger projects so much easier to
  find things.\n\n\n                # Making a fuzzy picker for your posts\n\n                Here
  is a bash command to open an fzf picker for todays posts,\n                then
  open it in your `$EDITOR`\n\n                ``` bash\n                markata list
  \\\n                        --map path\\\n                        --filter 'date==today'\\\n
  \                       --sort date\\\n                        --reverse |\\\n                        fzf
  --preview 'bat --color always {}' |\\\n                        xargs -I {} $EDITOR
  {}\n                ```\n\n                # Combining wtih nvim Telescope\n\n                Here
  is the same command setup as a Telescope picker for neovim.\n\n                ```
  vim\n                nnoremap <leader>et <cmd>Telescope find_files find_command=markata,list,--map,path,--filter,date==today<cr>\n
  \               ```\n\n                If you have another way to open posts in
  your editor with\n                `markata list` I would love to accept a PR to
  add it to the\n                examples here.\n                \"\"\"\n\n                markata.console.quiet
  = True\n\n                tail = -tail if tail else tail\n                filtered
  = markata.map(map, filter, sort)\n                if not include_empty:\n                    filtered
  = [a for a in filtered if a != \"\"]\n                filtered = filtered[tail:head]\n
  \               if reverse:\n                    filtered = reversed(filtered)\n\n
  \               markata.console.quiet = False\n                if markata.console.is_terminal
  and use_pager:\n                    with markata.console.pager():\n                        for
  a in filtered:\n                            markata.console.print(a, style=\"purple\")\n
  \               else:\n                    for a in filtered:\n                        markata.console.print(a)\n```\n\n\n!!
  function <h2 id='clean' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>clean
  <em class='small'>function</em></h2>\n    Cleans up output generated by markata
  including both the output_dir and\n    the .markata_cache.\n\n    # Dry Run\n\n
  \   You can run with `--dry-run` to see what markata is about to do.\n\n    ```
  bash\n    markata clean --dry-run\n    [09:42:37] [DRYRUN] removing outptut directory:
  markout base_cli.py:371\n               [DRYRUN] removing cache directory: .markata.cache
  base_cli.py:377\n\n    ```\n\n    # Running clean\n\n    Running markata clean will
  fully delete all of the directories created\n    by markata.\n\n    ``` bash\n    markata
  clean\n    [09:53:04]  removing outptut directory: markout base_cli.py:394\n                removing
  cache directory: .markata.cache base_cli.py:405\n    ```\n\n    # Running Quietly\n\n
  \   Running with `--quiet` will remove all of the directories created by\n    markata
  without announcing what it is doing.\n\n    ``` bash\n    markata clean --quiet\n
  \   ```\n???+ source \"clean <em class='small'>source</em>\"\n\n```python\n\n        def
  clean(\n                quiet: bool = typer.Option(\n                    False,\n
  \                   \"--quiet\",\n                    \"-q\",\n                ),\n
  \               dry_run: bool = typer.Option(\n                    False,\n                    \"--dry-run\",\n
  \               ),\n            ):\n                \"\"\"\n                Cleans
  up output generated by markata including both the output_dir and\n                the
  .markata_cache.\n\n                # Dry Run\n\n                You can run with
  `--dry-run` to see what markata is about to do.\n\n                ``` bash\n                markata
  clean --dry-run\n                [09:42:37] [DRYRUN] removing outptut directory:
  markout base_cli.py:371\n                           [DRYRUN] removing cache directory:
  .markata.cache base_cli.py:377\n\n                ```\n\n                # Running
  clean\n\n                Running markata clean will fully delete all of the directories
  created\n                by markata.\n\n                ``` bash\n                markata
  clean\n                [09:53:04]  removing outptut directory: markout base_cli.py:394\n
  \                           removing cache directory: .markata.cache base_cli.py:405\n
  \               ```\n\n                # Running Quietly\n\n                Running
  with `--quiet` will remove all of the directories created by\n                markata
  without announcing what it is doing.\n\n                ``` bash\n                markata
  clean --quiet\n                ```\n                \"\"\"\n                _clean(markata=markata,
  quiet=quiet, dry_run=dry_run)\n```\n"
date: 0001-01-01
description: Markata This plugin enables Your Markata Site can be build completely
  from the command line. see the Markata list is a tool to help list out artile attributes
  r
html: "<!DOCTYPE html>\n<html lang=\"en\">\n<head>\n<title>Base_Cli.Py</title>\n<meta
  charset=\"utf-8\"/>\n<meta content=\"width=device-width, initial-scale=1\" name=\"viewport\"/>\n<meta
  content=\"Markata This plugin enables Your Markata Site can be build completely
  from the command line. see the Markata list is a tool to help list out artile attributes
  r\" name=\"description\"/>\n<link href=\"/static/favicon.ico\" rel=\"icon\" type=\"image/png\"/>\n<script>\n
  \       function setTheme(theme) {\n            document.documentElement.setAttribute(\"data-theme\",
  theme);\n        }\n\n        function detectColorSchemeOnLoad() {\n            //local
  storage is used to override OS theme settings\n            if (localStorage.getItem(\"theme\"))
  {\n                if (localStorage.getItem(\"theme\") == \"dark\") {\n                    setTheme(\"dark\");\n
  \               } else if (localStorage.getItem(\"theme\") == \"light\") {\n                    setTheme(\"light\");\n
  \               }\n            } else if (!window.matchMedia) {\n                //matchMedia
  method not supported\n                setTheme(\"light\");\n                return
  false;\n            } else if (window.matchMedia(\"(prefers-color-scheme: dark)\").matches)
  {\n                //OS theme setting detected as dark\n                setTheme(\"dark\");\n
  \           } else {\n                setTheme(\"light\");\n            }\n        }\n
  \       detectColorSchemeOnLoad();\n        document.addEventListener(\n            \"DOMContentLoaded\",\n
  \           function () {\n                //identify the toggle switch HTML element\n
  \               const toggleSwitch = document.querySelector(\n                    '#theme-switch
  input[type=\"checkbox\"]',\n                );\n\n                //function that
  changes the theme, and sets a localStorage variable to track the theme between page
  loads\n                function switchTheme(e) {\n                    if (e.target.checked)
  {\n                        localStorage.setItem(\"theme\", \"dark\");\n                        document.documentElement.setAttribute(\"data-theme\",
  \"dark\");\n                        toggleSwitch.checked = true;\n                    }
  else {\n                        localStorage.setItem(\"theme\", \"light\");\n                        document.documentElement.setAttribute(\"data-theme\",
  \"light\");\n                        toggleSwitch.checked = false;\n                    }\n
  \               }\n\n                //listener for changing themes\n                toggleSwitch.addEventListener(\"change\",
  switchTheme, false);\n\n                //pre-check the dark-theme checkbox if dark-theme
  is set\n                if (document.documentElement.getAttribute(\"data-theme\")
  == \"dark\") {\n                    toggleSwitch.checked = true;\n                }\n
  \           },\n            false,\n        );\n    </script>\n<style>\n      :root
  {\n        --color-bg: #1f2022;\n        --color-bg-2: ;\n        --color-bg-code:
  #1f2022;\n        --color-text: #eefbfe;\n        --color-link: #fb30c4; \n        --color-accent:
  #e1bd00c9;\n        --overlay-brightness: .85;\n        --body-width: 800px;\n      }\n
  \     [data-theme=\"dark\"] {\n        --color-bg: #1f2022;\n        --color-bg-2:
  ;\n        --color-bg-code: #1f2022;\n        --color-text: #eefbfe;\n        --color-link:
  #fb30c4; \n        --color-accent: #e1bd00c9;\n        --overlay-brightness: .85;\n
  \       --body-width: 800px;\n      }\n      [data-theme=\"light\"] {\n        --color-bg:
  #eefbfe;\n        --color-bg-2: ;\n        --color-bg-code: #eefbfe;\n        --color-text:
  #1f2022;\n        --color-link: #fb30c4; \n        --color-accent: #ffeb00;\n        --overlay-brightness:
  .95;\n      }\n\n        html {\n            font-family: \"Space Mono\", monospace;\n
  \           background: var(--color-bg);\n            color: var(--color-text);\n
  \       }\n\n        a {\n            color: var(--color-link);\n        }\n\n        main
  a {\n            max-width: 100%;\n        }\n\n        .heading-permalink {\n            font-size:
  .7em;\n        }\n\n        body {\n            max-width: var(--body-width);\n
  \           margin: 5rem auto;\n            padding: 0 .5rem;\n            font-size:
  1rem;\n            line-height: 1.56;\n        }\n\n        blockquote {\n            background:
  var(--color-bg);\n            filter: brightness(var(--overlay-brightness));\n            border-left:
  4px solid var(--color-accent);\n            border-radius: 4px;\n            box-shadow:\n
  \               -0.8rem 0rem 1rem -1rem #f1fa8c,\n                0.2rem 0rem 1rem
  rgb(0, 0, 0, .4);\n            padding-left: 1rem;\n            margin: 1rem;\n
  \       }\n\n        li.post {\n            list-style-type: None;\n            padding:
  .2rem 0;\n        }\n\n        pre.wrapper {\n            padding: 0;\n            box-shadow:
  0.2rem 0rem 1rem rgb(0, 0, 0, .4);\n            display: flex;\n            flex-direction:
  column;\n            position: relative;\n            margin: 2rem;\n        }\n\n
  \       pre {\n            margin: 0;\n            padding: 1rem;\n            min-width:
  -webkit-fill-available;\n            max-width: fit-content;\n            overflow-x:
  auto;\n        }\n\n        pre .filepath {\n            margin: 0;\n            padding-left:
  1rem;\n            border-radius: 4px 4px 0 0;\n            background: black;\n
  \           display: flex;\n            justify-content: space-between;\n            align-items:
  center;\n        }\n\n        pre .filepath p {\n            margin: 0\n        }\n\n
  \       pre .filepath .right {\n            display: flex;\n            gap: .2rem;\n
  \           align-items: center;\n        }\n\n        pre::-webkit-scrollbar {\n
  \           height: 4px;\n            background-color: transparent;\n        }\n\n
  \       pre::-webkit-scrollbar-thumb {\n            background-color: #d3d3d32e;\n
  \           border-radius: 2px;\n        }\n\n        pre::-webkit-scrollbar-track
  {\n            background-color: transparent;\n        }\n\n        .copy-wrapper
  {\n            background: none;\n            position: absolute;\n            width:
  100%;\n            z-index: 100;\n            display: flex;\n            justify-content:
  flex-end;\n        }\n\n        button.copy {\n            z-index: 100;\n            background:
  none;\n            fill: #ffffff45;\n            border: none;\n            width:
  32px;\n            align-self: flex-end;\n            top: 0;\n            right:
  0;\n            margin: 0.5rem 0.2rem;\n\n        }\n\n        button.copy:hover
  {\n            fill: white\n        }\n\n        a.help {\n            fill: #ffffff45;\n
  \       }\n\n        a.help:hover {\n            fill: white;\n        }\n\n        a.help
  svg {\n            height: 24px;\n            width: 24px;\n        }\n\n        .highlight
  {\n            background: var(--color-bg-code);\n            color: var(--color-text);\n
  \           filter: brightness(var(--overlay-brightness));\n            border-radius:
  0 0 4px 4px;\n        }\n\n        .highlight .c {\n            color: #8b8b8b\n
  \       }\n\n        /* Comment */\n        .highlight .err {\n            color:
  #960050;\n            background-color: #1e0010\n        }\n\n        /* Error */\n
  \       .highlight .k {\n            color: #c678dd\n        }\n\n        /* Keyword
  */\n        .highlight .l {\n            color: #ae81ff\n        }\n\n        /*
  Literal */\n        .highlight .n {\n            color: #abb2bf\n        }\n\n        /*
  Name */\n        .highlight .o {\n            color: #c678dd\n        }\n\n        /*
  Operator */\n        .highlight .p {\n            color: #abb2bf\n        }\n\n
  \       /* Punctuation */\n        .highlight .ch {\n            color: #8b8b8b\n
  \       }\n\n        /* Comment.Hashbang */\n        .highlight .cm {\n            color:
  #8b8b8b\n        }\n\n        /* Comment.Multiline */\n        .highlight .cp {\n
  \           color: #8b8b8b\n        }\n\n        /* Comment.Preproc */\n        .highlight
  .cpf {\n            color: #8b8b8b\n        }\n\n        /* Comment.PreprocFile
  */\n        .highlight .c1 {\n            color: #8b8b8b\n        }\n\n        /*
  Comment.Single */\n        .highlight .cs {\n            color: #8b8b8b\n        }\n\n
  \       /* Comment.Special */\n        .highlight .gd {\n            color: #c678dd\n
  \       }\n\n        /* Generic.Deleted */\n        .highlight .ge {\n            font-style:
  italic\n        }\n\n        /* Generic.Emph */\n        .highlight .gi {\n            color:
  #a6e22e\n        }\n\n        /* Generic.Inserted */\n        .highlight .gs {\n
  \           font-weight: bold\n        }\n\n        /* Generic.Strong */\n        .highlight
  .gu {\n            color: #8b8b8b\n        }\n\n        /* Generic.Subheading */\n
  \       .highlight .kc {\n            color: #c678dd\n        }\n\n        /* Keyword.Constant
  */\n        .highlight .kd {\n            color: #c678dd\n        }\n\n        /*
  Keyword.Declaration */\n        .highlight .kn {\n            color: #c678dd\n        }\n\n
  \       /* Keyword.Namespace */\n        .highlight .kp {\n            color: #c678dd\n
  \       }\n\n        /* Keyword.Pseudo */\n        .highlight .kr {\n            color:
  #c678dd\n        }\n\n        /* Keyword.Reserved */\n        .highlight .kt {\n
  \           color: #c678dd\n        }\n\n        /* Keyword.Type */\n        .highlight
  .ld {\n            color: #e6db74\n        }\n\n        /* Literal.Date */\n        .highlight
  .m {\n            color: #ae81ff\n        }\n\n        /* Literal.Number */\n        .highlight
  .s {\n            color: #e6db74\n        }\n\n        /* Literal.String */\n        .highlight
  .na {\n            color: #a6e22e\n        }\n\n        /* Name.Attribute */\n        .highlight
  .nb {\n            color: #98c379\n        }\n\n        /* Name.Builtin */\n        .highlight
  .nc {\n            color: #abb2bf\n        }\n\n        /* Name.Class */\n        .highlight
  .no {\n            color: #c678dd\n        }\n\n        /* Name.Constant */\n        .highlight
  .nd {\n            color: #abb2bf\n        }\n\n        /* Name.Decorator */\n        .highlight
  .ni {\n            color: #abb2bf\n        }\n\n        /* Name.Entity */\n        .highlight
  .ne {\n            color: #a6e22e\n        }\n\n        /* Name.Exception */\n        .highlight
  .nf {\n            color: #61afef\n        }\n\n        /* Name.Function */\n        .highlight
  .nl {\n            color: #abb2bf\n        }\n\n        /* Name.Label */\n        .highlight
  .nn {\n            color: #abb2bf\n        }\n\n        /* Name.Namespace */\n        .highlight
  .nx {\n            color: #a6e22e\n        }\n\n        /* Name.Other */\n        .highlight
  .py {\n            color: #abb2bf\n        }\n\n        /* Name.Property */\n        .highlight
  .nt {\n            color: #c678dd\n        }\n\n        /* Name.Tag */\n        .highlight
  .nv {\n            color: #abb2bf\n        }\n\n        /* Name.Variable */\n        .highlight
  .ow {\n            color: #c678dd\n        }\n\n        /* Operator.Word */\n        .highlight
  .w {\n            color: #abb2bf\n        }\n\n        /* Text.Whitespace */\n        .highlight
  .mb {\n            color: #ae81ff\n        }\n\n        /* Literal.Number.Bin */\n
  \       .highlight .mf {\n            color: #ae81ff\n        }\n\n        /* Literal.Number.Float
  */\n        .highlight .mh {\n            color: #ae81ff\n        }\n\n        /*
  Literal.Number.Hex */\n        .highlight .mi {\n            color: #ae81ff\n        }\n\n
  \       /* Literal.Number.Integer */\n        .highlight .mo {\n            color:
  #ae81ff\n        }\n\n        /* Literal.Number.Oct */\n        .highlight .sa {\n
  \           color: #e6db74\n        }\n\n        /* Literal.String.Affix */\n        .highlight
  .sb {\n            color: #e6db74\n        }\n\n        /* Literal.String.Backtick
  */\n        .highlight .sc {\n            color: #e6db74\n        }\n\n        /*
  Literal.String.Char */\n        .highlight .dl {\n            color: #e6db74\n        }\n\n
  \       /* Literal.String.Delimiter */\n        .highlight .sd {\n            color:
  #98c379\n        }\n\n        /* Literal.String.Doc */\n        .highlight .s2 {\n
  \           color: #98c379\n        }\n\n        /* Literal.String.Double */\n        .highlight
  .se {\n            color: #ae81ff\n        }\n\n        /* Literal.String.Escape
  */\n        .highlight .sh {\n            color: #e6db74\n        }\n\n        /*
  Literal.String.Heredoc */\n        .highlight .si {\n            color: #e6db74\n
  \       }\n\n        /* Literal.String.Interpol */\n        .highlight .sx {\n            color:
  #e6db74\n        }\n\n        /* Literal.String.Other */\n        .highlight .sr
  {\n            color: #e6db74\n        }\n\n        /* Literal.String.Regex */\n
  \       .highlight .s1 {\n            color: #e6db74\n        }\n\n        /* Literal.String.Single
  */\n        .highlight .ss {\n            color: #e6db74\n        }\n\n        /*
  Literal.String.Symbol */\n        .highlight .bp {\n            color: #abb2bf\n
  \       }\n\n        /* Name.Builtin.Pseudo */\n        .highlight .fm {\n            color:
  #61afef\n        }\n\n        /* Name.Function.Magic */\n        .highlight .vc
  {\n            color: #abb2bf\n        }\n\n        /* Name.Variable.Class */\n
  \       .highlight .vg {\n            color: #abb2bf\n        }\n\n        /* Name.Variable.Global
  */\n        .highlight .vi {\n            color: #abb2bf\n        }\n\n        /*
  Name.Variable.Instance */\n        .highlight .vm {\n            color: #abb2bf\n
  \       }\n\n        /* Name.Variable.Magic */\n        .highlight .il {\n            color:
  #ae81ff\n        }\n\n        /* Literal.Number.Integer.Long */\n\n        /* Tab
  style starts here */\n        .tabbed-set {\n            position: relative;\n            display:
  flex;\n            flex-wrap: wrap;\n            margin: 1em 0;\n            border-radius:
  0.1rem;\n        }\n\n        .tabbed-set>input {\n            display: none;\n
  \       }\n\n        .tabbed-set label {\n            width: auto;\n            padding:
  0.9375em 1.25em 0.78125em;\n            font-weight: 700;\n            font-size:
  0.84em;\n            white-space: nowrap;\n            border-bottom: 0.15rem solid
  transparent;\n            border-top-left-radius: 0.1rem;\n            border-top-right-radius:
  0.1rem;\n            cursor: pointer;\n            transition: background-color
  250ms, color 250ms;\n        }\n\n        .tabbed-set .tabbed-content {\n            width:
  100%;\n            display: none;\n            box-shadow: 0 -.05rem #ddd;\n        }\n\n
  \       .tabbed-set input {\n            position: absolute;\n            opacity:
  0;\n        }\n\n        /* fonts */\n        h1 {\n            font-weight: 700;\n
  \       }\n\n        h1#title a {\n            font-size: 16px;\n        }\n\n        h1,\n
  \       h2,\n        h3,\n        h4,\n        h5,\n        h6 {\n            margin-top:
  3rem;\n        }\n\n        h1 {\n            font-size: 2.5em;\n            margin-top:
  5rem;\n        }\n\n        h2 {\n            font-size: 1.63rem;\n            margin-top:
  5rem;\n        }\n\n\n\n        p {\n            font-size: 21px;\n            font-style:
  normal;\n            font-variant: normal;\n            font-weight: 400;\n            line-height:
  1.5;\n        }\n\n        @media only screen and (max-width: 700px) {\n            p
  {\n                font-size: 18px;\n            }\n        }\n\n        @media
  only screen and (max-width: 600px) {\n            p {\n                font-size:
  16px;\n            }\n        }\n\n        @media only screen and (max-width: 500px)
  {\n            p {\n                font-size: 14px;\n            }\n        }\n\n
  \       @media only screen and (max-width: 400px) {\n            p {\n                font-size:
  12px;\n            }\n        }\n\n\n        pre {\n            font-style: normal;\n
  \           font-variant: normal;\n            font-weight: 400;\n            line-height:
  18.5714px;\n            */\n        }\n\n        a {\n            font-weight: 600;\n
  \           text-decoration-color: var(--color-accent);\n            color: var(--color-link);\n
  \           padding: .3rem .5rem;\n            display: inline-block;\n        }\n\n
  \       .admonition,\n        details {\n            box-shadow: 0.2rem 0rem 1rem
  rgb(0, 0, 0, .4);\n            margin: 5rem 0;\n            border: 1px solid transparent;\n
  \           border-radius: 4px;\n            text-align: left;\n            padding:
  0;\n            border: 0;\n\n        }\n\n        .admonition {\n            padding-bottom:
  1rem;\n        }\n\n        details[open] {\n            padding-bottom: .5rem;\n
  \       }\n\n        .admonition p {\n            padding: .2rem .6rem;\n        }\n\n
  \       .admonition-title,\n        .details-title,\n        summary {\n            background:
  var(--color-bg-2);\n            padding: 0;\n            margin: 0;\n            position:
  sticky;\n            top: 0;\n            z-index: 10;\n        }\n\n        summary:hover
  {\n            cursor: pointer;\n        }\n\n        summary.admonition-title,\n
  \       summary.details-title {\n            padding: .5rem;\n            padding-left:
  1rem;\n        }\n\n        .note {\n            border-left: 4px solid #f1fa8c;\n
  \           box-shadow:\n                -0.8rem 0rem 1rem -1rem #f1fa8c,\n                0.2rem
  0rem 1rem rgb(0, 0, 0, .4);\n        }\n\n        .note>.admonition-title {\n            border-bottom:
  1px solid #3c3d2d;\n        }\n\n        .abstract {\n            border-left: 4px
  solid #8be9fd;\n            box-shadow:\n                -0.8rem 0rem 1rem -1rem
  #8be9fd,\n                0.2rem 0rem 1rem rgb(0, 0, 0, .4);\n        }\n\n        .abstract>.admonition-title
  {\n            border-bottom: 1px solid #2c3a3f;\n        }\n\n        .info {\n
  \           border-left: 4px solid;\n            box-shadow:\n                -0.8rem
  0rem 1rem -1rem #8bb0fd,\n                0.2rem 0rem 1rem rgb(0, 0, 0, .4);\n        }\n\n
  \       .info>.admonition-title {\n            border-bottom: 1px solid #2c313f;\n
  \       }\n\n        .tip {\n            border-left: 4px solid #008080;\n            box-shadow:\n
  \               -0.8rem 0rem 1rem -1rem #008080,\n                0.2rem 0rem 1rem
  rgb(0, 0, 0, .4);\n        }\n\n        .tip>.admonition-title {\n            border-bottom:
  1px solid #1b2a2b;\n        }\n\n        .success {\n            border-left: 4px
  solid #50fa7b;\n            box-shadow:\n                -0.8rem 0rem 1rem -1rem
  #50fa7b,\n                0.2rem 0rem 1rem rgb(0, 0, 0, .4);\n        }\n\n        .success>.admonition-title
  {\n            border-bottom: 1px solid #263e2b;\n        }\n\n        .question
  {\n            border-left: 4px solid #a7fcbd;\n            box-shadow:\n                -0.8rem
  0rem 1rem -1rem #a7fcbd,\n                0.2rem 0rem 1rem rgb(0, 0, 0, .4);\n        }\n\n
  \       .question>.admonition-title {\n            border-bottom: 1px solid #303e35;\n
  \       }\n\n        .warning {\n            border-left: 4px solid #ffb86c;\n            box-shadow:\n
  \               -0.8rem 0rem 1rem -1rem #ffb86c,\n                0.2rem 0rem 1rem
  rgb(0, 0, 0, .4);\n        }\n\n        .warning>.admonition-title {\n            border-bottom:
  1px solid #3f3328;\n        }\n\n        .failure {\n            border-left: 4px
  solid #b23b3b;\n            box-shadow:\n                -0.8rem 0rem 1rem -1rem
  #b23b3b,\n                0.2rem 0rem 1rem rgb(0, 0, 0, .4);\n        }\n\n        .failure>.admonition-title
  {\n            border-bottom: 1px solid #34201f;\n        }\n\n        .danger {\n
  \           border-left: 4px solid #ff5555;\n            box-shadow:\n                -0.8rem
  0rem 1rem -1rem #ff5555,\n                0.2rem 0rem 1rem rgb(0, 0, 0, .4);\n        }\n\n
  \       .danger>.admonition-title {\n            border-bottom: 1px solid #402523;\n
  \       }\n\n        .bug {\n            border-left: 4px solid #b2548a;\n            box-shadow:\n
  \               -0.8rem 0rem 1rem -1rem #b2548a,\n                0.2rem 0rem 1rem
  rgb(0, 0, 0, .4);\n        }\n\n        .bug>.admonition-title {\n            border-bottom:
  1px solid #32232c;\n        }\n\n        .example {\n            border-left: 4px
  solid #bd93f9;\n            box-shadow:\n                -0.8rem 0rem 1rem -1rem
  #bd93f9,\n                0.2rem 0rem 1rem rgb(0, 0, 0, .4);\n        }\n\n        .example>.admonition-title
  {\n            border-bottom: 1px solid #332d3e;\n        }\n\n        .source {\n
  \           border-left: 4px solid #bd93f9;\n            box-shadow:\n                -0.8rem
  0rem 1rem -1rem #bd93f9,\n                0.2rem 0rem 1rem rgb(0, 0, 0, .4);\n        }\n\n
  \       .source>.admonition-title {\n            border-bottom: 1px solid #332d3e;\n
  \       }\n\n        .quote {\n            border-left: 4px solid #999;\n            box-shadow:\n
  \               -0.8rem 0rem 1rem -1rem #999,\n                0.2rem 0rem 1rem
  rgb(0, 0, 0, .4);\n        }\n\n        .quote>.admonition-title {\n            border-bottom:
  1px solid #2d2e2f;\n        }\n\n        table {\n            margin: 1rem 0;\n
  \           border-collapse: collapse;\n            border-spacing: 0;\n            display:
  block;\n            max-width: -moz-fit-content;\n            max-width: fit-content;\n
  \           overflow-x: auto;\n            white-space: nowrap;\n        }\n\n        table
  thead th {\n            border: solid 1px var(--color-text);\n            padding:
  10px;\n            text-align: left;\n        }\n\n        table tbody td {\n            border:
  solid 1px var(--color-text);\n            padding: 10px;\n        }\n\n        .theme-switch
  {\n            z-index: 10;\n            display: inline-block;\n            height:
  34px;\n            position: relative;\n            width: 60px;\n\n            display:
  flex;\n            justify-content: flex-end;\n            margin-right: 1rem;\n
  \           margin-left: auto;\n            position: fixed;\n            right:
  1rem;\n            top: 1rem;\n        }\n\n        .theme-switch input {\n            display:
  none;\n\n        }\n\n        .slider {\n            background-color: #ccc;\n            bottom:
  0;\n            cursor: pointer;\n            left: 0;\n            position: absolute;\n
  \           right: 0;\n            top: 0;\n            transition: .4s;\n        }\n\n
  \       .slider:before {\n            background-color: #fff;\n            bottom:
  4px;\n            content: \"\";\n            height: 26px;\n            left: 4px;\n
  \           position: absolute;\n            transition: .4s;\n            width:
  26px;\n        }\n\n        input:checked+.slider {\n            background-color:
  #343434;\n        }\n\n        input:checked+.slider:before {\n            background-color:
  #848484;\n        }\n\n        input:checked+.slider:before {\n            transform:
  translateX(26px);\n        }\n\n        .slider.round {\n            border-radius:
  34px;\n        }\n\n        .slider.round:before {\n            border-radius: 50%;\n
  \       }\n\n        main p img {\n            width: 100%;\n            width:
  -moz-available;\n            width: -webkit-fill-available;\n            width:
  fill-available;\n        }\n\n        details>* {\n            margin: 1rem;\n        }\n\n
  \       .admonition>* {\n            margin: 1rem;\n        }\n\n        p.admonition-title,\n
  \       summary {\n            margin: 0;\n            padding-left: 1.2rem;\n        }\n\n
  \       .small {\n            font-size: .9rem;\n            color: #888;\n        }\n\n
  \       admonition+admonition {\n            margin-top: 20rem;\n        }\n\n        ::-webkit-scrollbar
  {\n            height: 12px;\n            background-color: transparent;\n        }\n\n
  \       ::-webkit-scrollbar-thumb {\n            background-color: #d3d3d32e;\n
  \           border-radius: 6px;\n        }\n\n        ::-webkit-scrollbar-track
  {\n            background-color: transparent;\n        }\n    </style>\n<script>\n
  \       if (\"serviceWorker\" in navigator) {\n            navigator.serviceWorker.register(\"/service-worker.js\");\n
  \           navigator.serviceWorker.addEventListener(\"controllerchange\", () =>
  {\n                console.log(\"new worker\");\n                window.location.reload();\n
  \           });\n        }\n    </script>\n<meta content=\"waylon@waylonwalker.com\"
  name=\"og:author_email\"/>\n<meta content=\"waylon@waylonwalker.com\" name=\"og:author_email\"/>\n<meta
  content=\"Waylon Walker\" name=\"og:author\" property=\"og:author\"/><meta content=\"waylon@waylonwalaker.com\"
  name=\"og:author_email\" property=\"og:author_email\"/><meta content=\"website\"
  name=\"og:type\" property=\"og:type\"/><meta content=\"Markata This plugin enables
  Your Markata Site can be build completely from the command line. see the Markata
  list is a tool to help list out artile attributes r\" name=\"description\" property=\"description\"/><meta
  content=\"Markata This plugin enables Your Markata Site can be build completely
  from the command line. see the Markata list is a tool to help list out artile attributes
  r\" name=\"og:description\" property=\"og:description\"/><meta content=\"Markata
  This plugin enables Your Markata Site can be build completely from the command line.
  see the Markata list is a tool to help list out artile attributes r\" name=\"twitter:description\"
  property=\"twitter:description\"/><meta content=\"Base_Cli.Py | Markata\" name=\"og:title\"
  property=\"og:title\"/><meta content=\"Base_Cli.Py | Markata\" name=\"twitter:title\"
  property=\"twitter:title\"/><meta content=\"https://markata.dev//markata/plugins/base-cli-og.png\"
  name=\"og:image\" property=\"og:image\"/><meta content=\"https://markata.dev//markata/plugins/base-cli-og.png\"
  name=\"twitter:image\" property=\"twitter:image\"/><meta content=\"1600\" name=\"og:image:width\"
  property=\"og:image:width\"/><meta content=\"900\" name=\"og:image:width\" property=\"og:image:width\"/><meta
  content=\"summary_large_image\" name=\"twitter:card\" property=\"twitter:card\"/><meta
  content=\"Markata\" name=\"og:site_name\" property=\"og:site_name\"/><meta content=\"@_waylonwalker\"
  name=\"twitter:creator\" property=\"twitter:creator\"/><meta content=\"Base_Cli.Py\"
  name=\"title\" property=\"title\"/><meta content=\"markata 0.8.0.dev6\" name=\"generator\"
  property=\"generator\"/><link href=\"https://markata.dev//markata/plugins/base-cli/\"
  rel=\"canonical\"/><meta content=\"https://markata.dev//markata/plugins/base-cli/\"
  name=\"og:url\" property=\"og:url\"/></head>\n<body><nav>\n<a href=\"/\">markata</a>\n<a
  href=\"https://github.com/WaylonWalker/markata\">GitHub</a>\n</nav>\n<div>\n<label
  class=\"theme-switch\" for=\"checkbox-theme\" id=\"theme-switch\" title=\"light/dark
  mode toggle\">\n<input id=\"checkbox-theme\" type=\"checkbox\"/>\n<div class=\"slider
  round\"></div>\n</label>\n</div>\n<section class=\"title\">\n<h1 id=\"title\">\n
  \           Base_Cli.Py \n            \n        </h1>\n</section>\n<main><p>Markata's
  base command line commands.</p>\n<p>This plugin enables\n<a href=\"https://markata.dev/markata/plugins/base_cli/#build-function\"><code>build</code></a>\nand\n<a
  href=\"https://markata.dev/markata/plugins/base_cli/#list-function\"><code>list</code></a>\ncommands
  as part of the main markata cli.</p>\n<h1 id=\"building-your-site-with-the-cli\">Building
  Your Site with the Cli <a class=\"header-anchor\" href=\"#building-your-site-with-the-cli\"><svg
  aria-hidden=\"true\" class=\"heading-permalink\" fill=\"currentColor\" focusable=\"false\"
  height=\"1em\" viewbox=\"0 0 24 24\" width=\"1em\" xmlns=\"http://www.w3.org/2000/svg\"><path
  d=\"M9.199 13.599a5.99 5.99 0 0 0 3.949 2.345 5.987 5.987 0 0 0 5.105-1.702l2.995-2.994a5.992
  5.992 0 0 0 1.695-4.285 5.976 5.976 0 0 0-1.831-4.211 5.99 5.99 0 0 0-6.431-1.242
  6.003 6.003 0 0 0-1.905 1.24l-1.731 1.721a.999.999 0 1 0 1.41 1.418l1.709-1.699a3.985
  3.985 0 0 1 2.761-1.123 3.975 3.975 0 0 1 2.799 1.122 3.997 3.997 0 0 1 .111 5.644l-3.005
  3.006a3.982 3.982 0 0 1-3.395 1.126 3.987 3.987 0 0 1-2.632-1.563A1 1 0 0 0 9.201
  13.6zm5.602-3.198a5.99 5.99 0 0 0-3.949-2.345 5.987 5.987 0 0 0-5.105 1.702l-2.995
  2.994a5.992 5.992 0 0 0-1.695 4.285 5.976 5.976 0 0 0 1.831 4.211 5.99 5.99 0 0
  0 6.431 1.242 6.003 6.003 0 0 0 1.905-1.24l1.723-1.723a.999.999 0 1 0-1.414-1.414L9.836
  19.81a3.985 3.985 0 0 1-2.761 1.123 3.975 3.975 0 0 1-2.799-1.122 3.997 3.997 0
  0 1-.111-5.644l3.005-3.006a3.982 3.982 0 0 1 3.395-1.126 3.987 3.987 0 0 1 2.632
  1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<p>Your Markata Site can
  be build completely from the command line.</p>\n<pre class=\"wrapper\">\n\n<div
  class=\"copy-wrapper\">\n\n<button class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
  class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
  v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
  c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
  c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
  h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
  M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
  v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
  M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span>markata<span class=\"w\"> </span>build\n\n<span
  class=\"c1\"># or if you prefer pipx</span>\npipx<span class=\"w\"> </span>run<span
  class=\"w\"> </span>markata<span class=\"w\"> </span>build\n</pre></div>\n\n</pre>\n<p>see
  the\n<a href=\"https://markata.dev/markata/plugins/base_cli/#build-function\"><code>build</code></a>\nsection
  for more examples.</p>\n<h1 id=\"listing-your-articles\">Listing your articles <a
  class=\"header-anchor\" href=\"#listing-your-articles\"><svg aria-hidden=\"true\"
  class=\"heading-permalink\" fill=\"currentColor\" focusable=\"false\" height=\"1em\"
  viewbox=\"0 0 24 24\" width=\"1em\" xmlns=\"http://www.w3.org/2000/svg\"><path d=\"M9.199
  13.599a5.99 5.99 0 0 0 3.949 2.345 5.987 5.987 0 0 0 5.105-1.702l2.995-2.994a5.992
  5.992 0 0 0 1.695-4.285 5.976 5.976 0 0 0-1.831-4.211 5.99 5.99 0 0 0-6.431-1.242
  6.003 6.003 0 0 0-1.905 1.24l-1.731 1.721a.999.999 0 1 0 1.41 1.418l1.709-1.699a3.985
  3.985 0 0 1 2.761-1.123 3.975 3.975 0 0 1 2.799 1.122 3.997 3.997 0 0 1 .111 5.644l-3.005
  3.006a3.982 3.982 0 0 1-3.395 1.126 3.987 3.987 0 0 1-2.632-1.563A1 1 0 0 0 9.201
  13.6zm5.602-3.198a5.99 5.99 0 0 0-3.949-2.345 5.987 5.987 0 0 0-5.105 1.702l-2.995
  2.994a5.992 5.992 0 0 0-1.695 4.285 5.976 5.976 0 0 0 1.831 4.211 5.99 5.99 0 0
  0 6.431 1.242 6.003 6.003 0 0 0 1.905-1.24l1.723-1.723a.999.999 0 1 0-1.414-1.414L9.836
  19.81a3.985 3.985 0 0 1-2.761 1.123 3.975 3.975 0 0 1-2.799-1.122 3.997 3.997 0
  0 1-.111-5.644l3.005-3.006a3.982 3.982 0 0 1 3.395-1.126 3.987 3.987 0 0 1 2.632
  1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<p>Markata list is a tool
  to help list out artile attributes right to your\nterminal.  This is very helpful
  to find articles on larger sites, or\ndebug what is getting picked up by markata.</p>\n<pre
  class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
  class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
  v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
  c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
  c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
  h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
  M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
  v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
  M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span>markata<span class=\"w\"> </span>list<span
  class=\"w\"> </span>--map<span class=\"w\"> </span><span class=\"s1\">'str(date.year)
  + \",\" + title'</span>\n</pre></div>\n\n</pre>\n<p>see the\n<a href=\"https://markata.dev/markata/plugins/base_cli/#list-function\"><code>list</code></a>\nsection
  for more examples.</p>\n<h1 id=\"creating-new-things-with-the-cli\">Creating \"new\"
  things with the cli <a class=\"header-anchor\" href=\"#creating-new-things-with-the-cli\"><svg
  aria-hidden=\"true\" class=\"heading-permalink\" fill=\"currentColor\" focusable=\"false\"
  height=\"1em\" viewbox=\"0 0 24 24\" width=\"1em\" xmlns=\"http://www.w3.org/2000/svg\"><path
  d=\"M9.199 13.599a5.99 5.99 0 0 0 3.949 2.345 5.987 5.987 0 0 0 5.105-1.702l2.995-2.994a5.992
  5.992 0 0 0 1.695-4.285 5.976 5.976 0 0 0-1.831-4.211 5.99 5.99 0 0 0-6.431-1.242
  6.003 6.003 0 0 0-1.905 1.24l-1.731 1.721a.999.999 0 1 0 1.41 1.418l1.709-1.699a3.985
  3.985 0 0 1 2.761-1.123 3.975 3.975 0 0 1 2.799 1.122 3.997 3.997 0 0 1 .111 5.644l-3.005
  3.006a3.982 3.982 0 0 1-3.395 1.126 3.987 3.987 0 0 1-2.632-1.563A1 1 0 0 0 9.201
  13.6zm5.602-3.198a5.99 5.99 0 0 0-3.949-2.345 5.987 5.987 0 0 0-5.105 1.702l-2.995
  2.994a5.992 5.992 0 0 0-1.695 4.285 5.976 5.976 0 0 0 1.831 4.211 5.99 5.99 0 0
  0 6.431 1.242 6.003 6.003 0 0 0 1.905-1.24l1.723-1.723a.999.999 0 1 0-1.414-1.414L9.836
  19.81a3.985 3.985 0 0 1-2.761 1.123 3.975 3.975 0 0 1-2.799-1.122 3.997 3.997 0
  0 1-.111-5.644l3.005-3.006a3.982 3.982 0 0 1 3.395-1.126 3.987 3.987 0 0 1 2.632
  1.563 1 1 0 0 0 1.602-1.198z\"></path></svg></a></h1>\n<p>The <code>new</code> cli
  is built on copier templates, and allows you to build a new blog\nfrom a starter
  repo, make new posts, and new plugins.  Before you start dumping\nnew things onto
  your site for the first time, make sure you have a clean git\nhistory fully backed
  up, or look at the template repos to fully understand\nthem.</p>\n<pre class=\"wrapper\">\n\n<div
  class=\"copy-wrapper\">\n\n<button class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
  class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
  v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
  c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
  c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
  h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
  M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
  v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
  M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span><span class=\"c1\"># create
  a new blog template</span>\n<span class=\"c1\"># copier requires you to specify
  a directory</span>\nmarkata<span class=\"w\"> </span>new<span class=\"w\"> </span>blog<span
  class=\"w\"> </span><span class=\"o\">[</span>directory<span class=\"o\">]</span>\n\n<span
  class=\"c1\"># create a new blog post</span>\nmarkata<span class=\"w\"> </span>new<span
  class=\"w\"> </span>post\n\n<span class=\"c1\"># create a new plugin</span>\nmarkata<span
  class=\"w\"> </span>new<span class=\"w\"> </span>plugin\n\n<span class=\"c1\">#
  for the most up to date help, just ask for help.</span>\nmarkata<span class=\"w\">
  </span>new<span class=\"w\"> </span>--help\n\nUsage:<span class=\"w\"> </span>markata<span
  class=\"w\"> </span>new<span class=\"w\"> </span><span class=\"o\">[</span>OPTIONS<span
  class=\"o\">]</span><span class=\"w\"> </span>COMMAND<span class=\"w\"> </span><span
  class=\"o\">[</span>ARGS<span class=\"o\">]</span>...\n\ncreate<span class=\"w\">
  </span>new<span class=\"w\"> </span>things<span class=\"w\"> </span>from<span class=\"w\">
  </span>templates\n\n\u256D\u2500<span class=\"w\"> </span>Options<span class=\"w\">
  </span>\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u256E\n\u2502<span
  class=\"w\"> </span>--help<span class=\"w\">          </span>Show<span class=\"w\">
  </span>this<span class=\"w\"> </span>message<span class=\"w\"> </span>and<span class=\"w\">
  </span>exit.<span class=\"w\">                                      </span>\u2502\n\u2570\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u256F\n\u256D\u2500<span
  class=\"w\"> </span>Commands<span class=\"w\"> </span>\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u256E\n\u2502<span
  class=\"w\"> </span>blog<span class=\"w\">    </span>Create<span class=\"w\"> </span>a<span
  class=\"w\"> </span>new<span class=\"w\"> </span>blog<span class=\"w\"> </span>from<span
  class=\"w\"> </span>using<span class=\"w\"> </span>the<span class=\"w\"> </span>template<span
  class=\"w\"> </span>from<span class=\"w\">                           </span>\u2502\n\u2502<span
  class=\"w\">         </span>https://github.com/WaylonWalker/markata-blog-starter.<span
  class=\"w\">                    </span>\u2502\n\u2502<span class=\"w\"> </span>plugin<span
  class=\"w\">  </span>Create<span class=\"w\"> </span>a<span class=\"w\"> </span>new<span
  class=\"w\"> </span>plugin<span class=\"w\"> </span>using<span class=\"w\"> </span>the<span
  class=\"w\"> </span>template<span class=\"w\"> </span>at<span class=\"w\">                                </span>\u2502\n\u2502<span
  class=\"w\">         </span>https://github.com/WaylonWalker/markata-plugin-template.<span
  class=\"w\">                 </span>\u2502\n\u2502<span class=\"w\"> </span>post<span
  class=\"w\">    </span>Create<span class=\"w\"> </span>new<span class=\"w\"> </span>blog<span
  class=\"w\"> </span>post<span class=\"w\"> </span><span class=\"k\">in</span><span
  class=\"w\"> </span>the<span class=\"w\"> </span>pages<span class=\"w\"> </span>directory<span
  class=\"w\"> </span>from<span class=\"w\"> </span>the<span class=\"w\"> </span>template<span
  class=\"w\"> </span>at<span class=\"w\">         </span>\u2502\n\u2502<span class=\"w\">
  \        </span>https://github.com/WaylonWalker/markata-post-template.<span class=\"w\">
  \                  </span>\u2502\n\u2570\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u256F\n</pre></div>\n\n</pre>\n<p>!!
  function </p><h2 class=\"admonition-title\" id=\"make_pretty\" style=\"margin:0;padding:.5rem
  1rem;\">make_pretty <em class=\"small\">function</em></h2>\nThis is a helper function
  that enables suppresses tracebacks from\nframeworks like <code>click</code> that
  can make your traceback long and hard\nto follow.  It also makes evrerything more
  colorful and easier to\nfollow.\n<div class=\"admonition source is-collapsible collapsible-open\">\n<p
  class=\"admonition-title\">\"make_pretty <em class=\"small\">source</em>\"</p>\n</div>\n<pre
  class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
  class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
  v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
  c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
  c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
  h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
  M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
  v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
  M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
  <span class=\"nf\">make_pretty</span><span class=\"p\">()</span> <span class=\"o\">-&gt;</span>
  <span class=\"kc\">None</span><span class=\"p\">:</span>\n<span class=\"w\">            </span><span
  class=\"sd\">\"\"\"</span>\n<span class=\"sd\">            This is a helper function
  that enables suppresses tracebacks from</span>\n<span class=\"sd\">            frameworks
  like `click` that can make your traceback long and hard</span>\n<span class=\"sd\">
  \           to follow.  It also makes evrerything more colorful and easier to</span>\n<span
  class=\"sd\">            follow.</span>\n<span class=\"sd\">            \"\"\"</span>\n
  \           <span class=\"kn\">import</span> <span class=\"nn\">click</span>\n            <span
  class=\"kn\">import</span> <span class=\"nn\">pluggy</span>\n            <span class=\"kn\">import</span>
  <span class=\"nn\">typer</span>\n            <span class=\"kn\">from</span> <span
  class=\"nn\">rich</span> <span class=\"kn\">import</span> <span class=\"n\">pretty</span>
  <span class=\"k\">as</span> <span class=\"n\">_pretty</span>\n            <span
  class=\"kn\">from</span> <span class=\"nn\">rich</span> <span class=\"kn\">import</span>
  <span class=\"n\">traceback</span>\n\n            <span class=\"n\">_pretty</span><span
  class=\"o\">.</span><span class=\"n\">install</span><span class=\"p\">()</span>\n
  \           <span class=\"n\">traceback</span><span class=\"o\">.</span><span class=\"n\">install</span><span
  class=\"p\">(</span>\n                <span class=\"n\">show_locals</span><span
  class=\"o\">=</span><span class=\"kc\">True</span><span class=\"p\">,</span>\n                <span
  class=\"n\">suppress</span><span class=\"o\">=</span><span class=\"p\">[</span>\n
  \                   <span class=\"n\">pluggy</span><span class=\"p\">,</span>\n
  \                   <span class=\"n\">click</span><span class=\"p\">,</span>\n                    <span
  class=\"n\">typer</span><span class=\"p\">,</span>\n                <span class=\"p\">],</span>\n
  \           <span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n<p>!! function </p><h2
  class=\"admonition-title\" id=\"cli\" style=\"margin:0;padding:.5rem 1rem;\">cli
  <em class=\"small\">function</em></h2>\nMarkata hook to implement base cli commands.\n<div
  class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"cli
  <em class=\"small\">source</em>\"</p>\n</div>\n<pre class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button
  class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
  class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
  v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
  c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
  c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
  h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
  M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
  v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
  M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
  <span class=\"nf\">cli</span><span class=\"p\">(</span><span class=\"n\">app</span><span
  class=\"p\">:</span> <span class=\"n\">typer</span><span class=\"o\">.</span><span
  class=\"n\">Typer</span><span class=\"p\">,</span> <span class=\"n\">markata</span><span
  class=\"p\">:</span> <span class=\"s2\">\"Markata\"</span><span class=\"p\">)</span>
  <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n<span
  class=\"w\">            </span><span class=\"sd\">\"\"\"</span>\n<span class=\"sd\">
  \           Markata hook to implement base cli commands.</span>\n<span class=\"sd\">
  \           \"\"\"</span>\n\n            <span class=\"n\">plugins_app</span> <span
  class=\"o\">=</span> <span class=\"n\">typer</span><span class=\"o\">.</span><span
  class=\"n\">Typer</span><span class=\"p\">()</span>\n            <span class=\"n\">config_app</span>
  <span class=\"o\">=</span> <span class=\"n\">typer</span><span class=\"o\">.</span><span
  class=\"n\">Typer</span><span class=\"p\">()</span>\n            <span class=\"n\">app</span><span
  class=\"o\">.</span><span class=\"n\">add_typer</span><span class=\"p\">(</span><span
  class=\"n\">plugins_app</span><span class=\"p\">)</span>\n            <span class=\"n\">app</span><span
  class=\"o\">.</span><span class=\"n\">add_typer</span><span class=\"p\">(</span><span
  class=\"n\">config_app</span><span class=\"p\">)</span>\n\n            <span class=\"nd\">@app</span><span
  class=\"o\">.</span><span class=\"n\">command</span><span class=\"p\">()</span>\n
  \           <span class=\"k\">def</span> <span class=\"nf\">tui</span><span class=\"p\">(</span><span
  class=\"n\">ctx</span><span class=\"p\">:</span> <span class=\"n\">typer</span><span
  class=\"o\">.</span><span class=\"n\">Context</span><span class=\"p\">)</span> <span
  class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
  \               <span class=\"k\">try</span><span class=\"p\">:</span>\n                    <span
  class=\"kn\">from</span> <span class=\"nn\">trogon</span> <span class=\"kn\">import</span>
  <span class=\"n\">Trogon</span>\n                    <span class=\"kn\">from</span>
  <span class=\"nn\">typer.main</span> <span class=\"kn\">import</span> <span class=\"n\">get_group</span>\n
  \               <span class=\"k\">except</span> <span class=\"ne\">ImportError</span><span
  class=\"p\">:</span>\n                    <span class=\"n\">typer</span><span class=\"o\">.</span><span
  class=\"n\">echo</span><span class=\"p\">(</span><span class=\"s2\">\"trogon not
  installed\"</span><span class=\"p\">)</span>\n                    <span class=\"n\">typer</span><span
  class=\"o\">.</span><span class=\"n\">echo</span><span class=\"p\">(</span>\n                        <span
  class=\"s2\">\"install markata with optional tui group to use tui `pip install 'markata[tui]'`\"</span>\n
  \                   <span class=\"p\">)</span>\n                    <span class=\"k\">return</span>\n\n
  \               <span class=\"n\">Trogon</span><span class=\"p\">(</span><span class=\"n\">get_group</span><span
  class=\"p\">(</span><span class=\"n\">app</span><span class=\"p\">),</span> <span
  class=\"n\">click_context</span><span class=\"o\">=</span><span class=\"n\">ctx</span><span
  class=\"p\">)</span><span class=\"o\">.</span><span class=\"n\">run</span><span
  class=\"p\">()</span>\n\n            <span class=\"nd\">@plugins_app</span><span
  class=\"o\">.</span><span class=\"n\">callback</span><span class=\"p\">()</span>\n
  \           <span class=\"k\">def</span> <span class=\"nf\">plugins</span><span
  class=\"p\">():</span>\n                <span class=\"s2\">\"create new things from
  templates\"</span>\n\n            <span class=\"nd\">@config_app</span><span class=\"o\">.</span><span
  class=\"n\">callback</span><span class=\"p\">()</span>\n            <span class=\"k\">def</span>
  <span class=\"nf\">config</span><span class=\"p\">():</span>\n                <span
  class=\"s2\">\"configuration management\"</span>\n\n            <span class=\"nd\">@config_app</span><span
  class=\"o\">.</span><span class=\"n\">command</span><span class=\"p\">()</span>\n
  \           <span class=\"k\">def</span> <span class=\"nf\">show</span><span class=\"p\">(</span>\n
  \               <span class=\"n\">verbose</span><span class=\"p\">:</span> <span
  class=\"nb\">bool</span> <span class=\"o\">=</span> <span class=\"n\">typer</span><span
  class=\"o\">.</span><span class=\"n\">Option</span><span class=\"p\">(</span>\n
  \                   <span class=\"kc\">False</span><span class=\"p\">,</span>\n
  \                   <span class=\"s2\">\"--verbose\"</span><span class=\"p\">,</span>\n
  \                   <span class=\"s2\">\"-v\"</span><span class=\"p\">,</span>\n
  \               <span class=\"p\">),</span>\n            <span class=\"p\">)</span>
  <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
  \               <span class=\"k\">if</span> <span class=\"n\">verbose</span><span
  class=\"p\">:</span>\n                    <span class=\"n\">markata</span><span
  class=\"o\">.</span><span class=\"n\">console</span><span class=\"o\">.</span><span
  class=\"n\">quiet</span> <span class=\"o\">=</span> <span class=\"kc\">False</span>\n
  \               <span class=\"k\">else</span><span class=\"p\">:</span>\n                    <span
  class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">console</span><span
  class=\"o\">.</span><span class=\"n\">quiet</span> <span class=\"o\">=</span> <span
  class=\"kc\">True</span>\n                <span class=\"n\">rich_print</span><span
  class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
  class=\"n\">config</span><span class=\"p\">)</span>\n\n            <span class=\"nd\">@config_app</span><span
  class=\"o\">.</span><span class=\"n\">command</span><span class=\"p\">()</span>\n
  \           <span class=\"k\">def</span> <span class=\"nf\">generate</span><span
  class=\"p\">(</span>\n                <span class=\"n\">verbose</span><span class=\"p\">:</span>
  <span class=\"nb\">bool</span> <span class=\"o\">=</span> <span class=\"n\">typer</span><span
  class=\"o\">.</span><span class=\"n\">Option</span><span class=\"p\">(</span>\n
  \                   <span class=\"kc\">False</span><span class=\"p\">,</span>\n
  \                   <span class=\"s2\">\"--verbose\"</span><span class=\"p\">,</span>\n
  \                   <span class=\"s2\">\"-v\"</span><span class=\"p\">,</span>\n
  \               <span class=\"p\">),</span>\n            <span class=\"p\">)</span>
  <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
  \               <span class=\"k\">if</span> <span class=\"n\">verbose</span><span
  class=\"p\">:</span>\n                    <span class=\"n\">markata</span><span
  class=\"o\">.</span><span class=\"n\">console</span><span class=\"o\">.</span><span
  class=\"n\">quiet</span> <span class=\"o\">=</span> <span class=\"kc\">False</span>\n
  \               <span class=\"k\">else</span><span class=\"p\">:</span>\n                    <span
  class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">console</span><span
  class=\"o\">.</span><span class=\"n\">quiet</span> <span class=\"o\">=</span> <span
  class=\"kc\">True</span>\n\n                <span class=\"n\">rich_print</span><span
  class=\"p\">(</span><span class=\"n\">toml</span><span class=\"o\">.</span><span
  class=\"n\">dumps</span><span class=\"p\">(</span><span class=\"n\">json</span><span
  class=\"o\">.</span><span class=\"n\">loads</span><span class=\"p\">(</span><span
  class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
  class=\"o\">.</span><span class=\"n\">json</span><span class=\"p\">())))</span>\n\n
  \           <span class=\"nd\">@config_app</span><span class=\"o\">.</span><span
  class=\"n\">command</span><span class=\"p\">()</span>\n            <span class=\"k\">def</span>
  <span class=\"nf\">get</span><span class=\"p\">(</span><span class=\"n\">key</span><span
  class=\"p\">:</span> <span class=\"nb\">str</span><span class=\"p\">)</span> <span
  class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
  \               <span class=\"n\">keys</span> <span class=\"o\">=</span> <span class=\"n\">key</span><span
  class=\"o\">.</span><span class=\"n\">split</span><span class=\"p\">(</span><span
  class=\"s2\">\".\"</span><span class=\"p\">)</span>\n                <span class=\"n\">markata</span><span
  class=\"o\">.</span><span class=\"n\">console</span><span class=\"o\">.</span><span
  class=\"n\">quiet</span> <span class=\"o\">=</span> <span class=\"kc\">True</span>\n
  \               <span class=\"n\">keys_processed</span> <span class=\"o\">=</span>
  <span class=\"s2\">\"\"</span>\n                <span class=\"n\">value</span> <span
  class=\"o\">=</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
  class=\"n\">config</span>\n                <span class=\"n\">na</span> <span class=\"o\">=</span>
  <span class=\"n\">Literal</span><span class=\"p\">[</span><span class=\"s2\">\"na\"</span><span
  class=\"p\">]</span>\n                <span class=\"k\">for</span> <span class=\"n\">key</span>
  <span class=\"ow\">in</span> <span class=\"n\">keys</span><span class=\"p\">:</span>\n
  \                   <span class=\"n\">value</span> <span class=\"o\">=</span> <span
  class=\"nb\">getattr</span><span class=\"p\">(</span><span class=\"n\">value</span><span
  class=\"p\">,</span> <span class=\"n\">key</span><span class=\"p\">,</span> <span
  class=\"n\">na</span><span class=\"p\">)</span>\n                    <span class=\"n\">keys_processed</span>
  <span class=\"o\">=</span> <span class=\"sa\">f</span><span class=\"s2\">\"</span><span
  class=\"si\">{</span><span class=\"n\">keys_processed</span><span class=\"si\">}</span><span
  class=\"s2\">.</span><span class=\"si\">{</span><span class=\"n\">key</span><span
  class=\"si\">}</span><span class=\"s2\">\"</span><span class=\"o\">.</span><span
  class=\"n\">strip</span><span class=\"p\">(</span><span class=\"s2\">\".\"</span><span
  class=\"p\">)</span>\n                    <span class=\"k\">if</span> <span class=\"n\">value</span>
  <span class=\"ow\">is</span> <span class=\"n\">na</span><span class=\"p\">:</span>\n
  \                       <span class=\"n\">rich_print</span><span class=\"p\">(</span><span
  class=\"sa\">f</span><span class=\"s2\">\"</span><span class=\"si\">{</span><span
  class=\"n\">keys_processed</span><span class=\"si\">}</span><span class=\"s2\">
  not found\"</span><span class=\"p\">)</span>\n                        <span class=\"n\">sys</span><span
  class=\"o\">.</span><span class=\"n\">exit</span><span class=\"p\">(</span><span
  class=\"mi\">1</span><span class=\"p\">)</span>\n\n                <span class=\"n\">rich_print</span><span
  class=\"p\">(</span><span class=\"n\">value</span><span class=\"p\">)</span>\n\n
  \           <span class=\"n\">new_app</span> <span class=\"o\">=</span> <span class=\"n\">typer</span><span
  class=\"o\">.</span><span class=\"n\">Typer</span><span class=\"p\">()</span>\n
  \           <span class=\"n\">app</span><span class=\"o\">.</span><span class=\"n\">add_typer</span><span
  class=\"p\">(</span><span class=\"n\">new_app</span><span class=\"p\">)</span>\n\n
  \           <span class=\"nd\">@new_app</span><span class=\"o\">.</span><span class=\"n\">callback</span><span
  class=\"p\">()</span>\n            <span class=\"k\">def</span> <span class=\"nf\">new</span><span
  class=\"p\">():</span>\n                <span class=\"s2\">\"create new things from
  templates\"</span>\n\n            <span class=\"nd\">@new_app</span><span class=\"o\">.</span><span
  class=\"n\">command</span><span class=\"p\">()</span>\n            <span class=\"k\">def</span>
  <span class=\"nf\">blog</span><span class=\"p\">(</span>\n                <span
  class=\"n\">directory</span><span class=\"p\">:</span> <span class=\"n\">Path</span>
  <span class=\"o\">=</span> <span class=\"n\">typer</span><span class=\"o\">.</span><span
  class=\"n\">Argument</span><span class=\"p\">(</span>\n                    <span
  class=\"o\">...</span><span class=\"p\">,</span>\n                    <span class=\"n\">help</span><span
  class=\"o\">=</span><span class=\"s2\">\"The directory to create the blog in.\"</span><span
  class=\"p\">,</span>\n                <span class=\"p\">),</span>\n            <span
  class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
  class=\"p\">:</span>\n<span class=\"w\">                </span><span class=\"sd\">\"\"\"</span>\n<span
  class=\"sd\">                Create a new blog from using the template from</span>\n<span
  class=\"sd\">                https://github.com/WaylonWalker/markata-blog-starter.</span>\n<span
  class=\"sd\">                \"\"\"</span>\n\n                <span class=\"kn\">from</span>
  <span class=\"nn\">copier</span> <span class=\"kn\">import</span> <span class=\"n\">run_copy</span>\n\n
  \               <span class=\"n\">typer</span><span class=\"o\">.</span><span class=\"n\">echo</span><span
  class=\"p\">(</span><span class=\"sa\">f</span><span class=\"s2\">\"creating a new
  project in </span><span class=\"si\">{</span><span class=\"n\">directory</span><span
  class=\"o\">.</span><span class=\"n\">absolute</span><span class=\"p\">()</span><span
  class=\"si\">}</span><span class=\"s2\">\"</span><span class=\"p\">)</span>\n                <span
  class=\"n\">url</span> <span class=\"o\">=</span> <span class=\"n\">markata</span><span
  class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
  class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">\"starters\"</span><span
  class=\"p\">,</span> <span class=\"p\">{})</span><span class=\"o\">.</span><span
  class=\"n\">get</span><span class=\"p\">(</span>\n                    <span class=\"s2\">\"blog\"</span><span
  class=\"p\">,</span>\n                    <span class=\"s2\">\"git+https://github.com/WaylonWalker/markata-blog-starter\"</span><span
  class=\"p\">,</span>\n                <span class=\"p\">)</span>\n                <span
  class=\"n\">run_copy</span><span class=\"p\">(</span><span class=\"n\">url</span><span
  class=\"p\">,</span> <span class=\"n\">directory</span><span class=\"p\">)</span>\n\n
  \           <span class=\"nd\">@new_app</span><span class=\"o\">.</span><span class=\"n\">command</span><span
  class=\"p\">()</span>\n            <span class=\"k\">def</span> <span class=\"nf\">post</span><span
  class=\"p\">()</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
  class=\"p\">:</span>\n<span class=\"w\">                </span><span class=\"sd\">\"\"\"</span>\n<span
  class=\"sd\">                Create new blog post in the pages directory from the
  template at</span>\n<span class=\"sd\">                https://github.com/WaylonWalker/markata-post-template.</span>\n<span
  class=\"sd\">                \"\"\"</span>\n\n                <span class=\"nb\">print</span><span
  class=\"p\">(</span><span class=\"s2\">\"create a new post\"</span><span class=\"p\">)</span>\n
  \               <span class=\"kn\">from</span> <span class=\"nn\">copier</span>
  <span class=\"kn\">import</span> <span class=\"n\">run_copy</span>\n\n                <span
  class=\"n\">typer</span><span class=\"o\">.</span><span class=\"n\">echo</span><span
  class=\"p\">(</span><span class=\"sa\">f</span><span class=\"s2\">\"creating a new
  post in </span><span class=\"si\">{</span><span class=\"n\">Path</span><span class=\"p\">()</span><span
  class=\"o\">.</span><span class=\"n\">absolute</span><span class=\"p\">()</span><span
  class=\"si\">}</span><span class=\"s2\">/posts\"</span><span class=\"p\">)</span>\n
  \               <span class=\"n\">url</span> <span class=\"o\">=</span> <span class=\"n\">markata</span><span
  class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
  class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">\"starters\"</span><span
  class=\"p\">,</span> <span class=\"p\">{})</span><span class=\"o\">.</span><span
  class=\"n\">get</span><span class=\"p\">(</span>\n                    <span class=\"s2\">\"post\"</span><span
  class=\"p\">,</span>\n                    <span class=\"s2\">\"git+https://github.com/WaylonWalker/markata-post-template\"</span><span
  class=\"p\">,</span>\n                <span class=\"p\">)</span>\n                <span
  class=\"n\">run_copy</span><span class=\"p\">(</span><span class=\"n\">url</span><span
  class=\"p\">,</span> <span class=\"n\">Path</span><span class=\"p\">(</span><span
  class=\"s2\">\".\"</span><span class=\"p\">))</span>\n\n            <span class=\"nd\">@new_app</span><span
  class=\"o\">.</span><span class=\"n\">command</span><span class=\"p\">()</span>\n
  \           <span class=\"k\">def</span> <span class=\"nf\">plugin</span><span class=\"p\">()</span>
  <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n<span
  class=\"w\">                </span><span class=\"sd\">\"\"\"</span>\n<span class=\"sd\">
  \               Create a new plugin using the template at</span>\n<span class=\"sd\">
  \               https://github.com/WaylonWalker/markata-plugin-template.</span>\n<span
  class=\"sd\">                \"\"\"</span>\n                <span class=\"kn\">from</span>
  <span class=\"nn\">copier</span> <span class=\"kn\">import</span> <span class=\"n\">run_copy</span>\n\n
  \               <span class=\"n\">typer</span><span class=\"o\">.</span><span class=\"n\">echo</span><span
  class=\"p\">(</span>\n                    <span class=\"sa\">f</span><span class=\"s2\">\"creating
  a new plugin in </span><span class=\"si\">{</span><span class=\"n\">Path</span><span
  class=\"p\">()</span><span class=\"o\">.</span><span class=\"n\">absolute</span><span
  class=\"p\">()</span><span class=\"si\">}</span><span class=\"s2\">\"</span>\n                    <span
  class=\"sa\">f</span><span class=\"s2\">\"/&lt;python-package-name&gt;/plugins\"</span><span
  class=\"p\">,</span>\n                <span class=\"p\">)</span>\n                <span
  class=\"n\">url</span> <span class=\"o\">=</span> <span class=\"n\">markata</span><span
  class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
  class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">\"starters\"</span><span
  class=\"p\">,</span> <span class=\"p\">{})</span><span class=\"o\">.</span><span
  class=\"n\">get</span><span class=\"p\">(</span>\n                    <span class=\"s2\">\"post\"</span><span
  class=\"p\">,</span>\n                    <span class=\"s2\">\"git+https://github.com/WaylonWalker/markata-plugin-template\"</span><span
  class=\"p\">,</span>\n                <span class=\"p\">)</span>\n                <span
  class=\"n\">run_copy</span><span class=\"p\">(</span><span class=\"n\">url</span><span
  class=\"p\">,</span> <span class=\"n\">Path</span><span class=\"p\">(</span><span
  class=\"s2\">\".\"</span><span class=\"p\">))</span>\n\n            <span class=\"nd\">@app</span><span
  class=\"o\">.</span><span class=\"n\">command</span><span class=\"p\">()</span>\n
  \           <span class=\"k\">def</span> <span class=\"nf\">build</span><span class=\"p\">(</span>\n
  \               <span class=\"n\">pretty</span><span class=\"p\">:</span> <span
  class=\"nb\">bool</span> <span class=\"o\">=</span> <span class=\"kc\">True</span><span
  class=\"p\">,</span>\n                <span class=\"n\">quiet</span><span class=\"p\">:</span>
  <span class=\"nb\">bool</span> <span class=\"o\">=</span> <span class=\"n\">typer</span><span
  class=\"o\">.</span><span class=\"n\">Option</span><span class=\"p\">(</span>\n
  \                   <span class=\"kc\">False</span><span class=\"p\">,</span>\n
  \                   <span class=\"s2\">\"--quiet\"</span><span class=\"p\">,</span>\n
  \                   <span class=\"s2\">\"-q\"</span><span class=\"p\">,</span>\n
  \               <span class=\"p\">),</span>\n                <span class=\"n\">verbose</span><span
  class=\"p\">:</span> <span class=\"nb\">bool</span> <span class=\"o\">=</span> <span
  class=\"n\">typer</span><span class=\"o\">.</span><span class=\"n\">Option</span><span
  class=\"p\">(</span>\n                    <span class=\"kc\">False</span><span class=\"p\">,</span>\n
  \                   <span class=\"s2\">\"--verbose\"</span><span class=\"p\">,</span>\n
  \                   <span class=\"s2\">\"-v\"</span><span class=\"p\">,</span>\n
  \               <span class=\"p\">),</span>\n                <span class=\"n\">should_pdb</span><span
  class=\"p\">:</span> <span class=\"nb\">bool</span> <span class=\"o\">=</span> <span
  class=\"n\">typer</span><span class=\"o\">.</span><span class=\"n\">Option</span><span
  class=\"p\">(</span>\n                    <span class=\"kc\">False</span><span class=\"p\">,</span>\n
  \                   <span class=\"s2\">\"--pdb\"</span><span class=\"p\">,</span>\n
  \               <span class=\"p\">),</span>\n                <span class=\"n\">profile</span><span
  class=\"p\">:</span> <span class=\"nb\">bool</span> <span class=\"o\">=</span> <span
  class=\"kc\">True</span><span class=\"p\">,</span>\n            <span class=\"p\">)</span>
  <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n<span
  class=\"w\">                </span><span class=\"sd\">\"\"\"</span>\n<span class=\"sd\">
  \               Markata's primary way of building your site for production.</span>\n<span
  class=\"sd\">                By default, running `markta build` will render your
  markdown to</span>\n<span class=\"sd\">                the `./markout` directory.</span>\n\n<span
  class=\"sd\">                ``` bash</span>\n<span class=\"sd\">                markata
  build</span>\n<span class=\"sd\">                ```</span>\n\n<span class=\"sd\">
  \               If you are having an issue and want to pop immediately into a debugger</span>\n<span
  class=\"sd\">                upon failure you can pass the `--pdb` flag to the build
  command.</span>\n\n<span class=\"sd\">                ``` bash</span>\n<span class=\"sd\">
  \               markata build  --pdb</span>\n<span class=\"sd\">                ```</span>\n\n<span
  class=\"sd\">                If you do not like the way rich looks, or its suppressing
  tracebaks you</span>\n<span class=\"sd\">                would like to remain visible
  you can use `--no-pretty`</span>\n\n<span class=\"sd\">                ``` bash</span>\n<span
  class=\"sd\">                markata build --no-pretty</span>\n<span class=\"sd\">
  \               ```</span>\n\n<span class=\"sd\">                If you need to
  run without any console logging pass in the</span>\n<span class=\"sd\">                `--quiet`
  flag.</span>\n\n<span class=\"sd\">                ``` bash</span>\n<span class=\"sd\">
  \               markata build --quiet</span>\n<span class=\"sd\">                ```</span>\n\n<span
  class=\"sd\">                `markta build` will automatically run the pyinstrument
  profiler</span>\n<span class=\"sd\">                while building your site if
  you have pyinstrument installed.  It</span>\n<span class=\"sd\">                will
  echo out your profile in the console as well as write it to</span>\n<span class=\"sd\">
  \               `/_profile` on your built site. If you prefer not to run</span>\n<span
  class=\"sd\">                pyinstrument profiling, even when it is installed you
  can pass</span>\n<span class=\"sd\">                in `--no-profile`</span>\n\n<span
  class=\"sd\">                ``` bash</span>\n<span class=\"sd\">                markata
  build --no-profile</span>\n<span class=\"sd\">                ```</span>\n<span
  class=\"sd\">                \"\"\"</span>\n\n                <span class=\"k\">if</span>
  <span class=\"n\">pretty</span><span class=\"p\">:</span>\n                    <span
  class=\"n\">make_pretty</span><span class=\"p\">()</span>\n\n                <span
  class=\"k\">if</span> <span class=\"n\">quiet</span><span class=\"p\">:</span>\n
  \                   <span class=\"n\">markata</span><span class=\"o\">.</span><span
  class=\"n\">console</span><span class=\"o\">.</span><span class=\"n\">quiet</span>
  <span class=\"o\">=</span> <span class=\"kc\">True</span>\n\n                <span
  class=\"k\">if</span> <span class=\"n\">verbose</span><span class=\"p\">:</span>\n
  \                   <span class=\"n\">markata</span><span class=\"o\">.</span><span
  class=\"n\">console</span><span class=\"o\">.</span><span class=\"n\">print</span><span
  class=\"p\">(</span><span class=\"s2\">\"console options:\"</span><span class=\"p\">,</span>
  <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">console</span><span
  class=\"o\">.</span><span class=\"n\">options</span><span class=\"p\">)</span>\n\n
  \               <span class=\"k\">if</span> <span class=\"n\">profile</span><span
  class=\"p\">:</span>\n                    <span class=\"n\">markata</span><span
  class=\"o\">.</span><span class=\"n\">should_profile_cli</span> <span class=\"o\">=</span>
  <span class=\"kc\">True</span>\n                    <span class=\"n\">markata</span><span
  class=\"o\">.</span><span class=\"n\">should_profile</span> <span class=\"o\">=</span>
  <span class=\"kc\">True</span>\n                    <span class=\"n\">markata</span><span
  class=\"o\">.</span><span class=\"n\">configure</span><span class=\"p\">()</span>\n\n
  \               <span class=\"k\">if</span> <span class=\"n\">should_pdb</span><span
  class=\"p\">:</span>\n                    <span class=\"n\">pdb_run</span><span
  class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
  class=\"n\">run</span><span class=\"p\">)</span>\n\n                <span class=\"k\">else</span><span
  class=\"p\">:</span>\n                    <span class=\"n\">markata</span><span
  class=\"o\">.</span><span class=\"n\">run</span><span class=\"p\">()</span>\n\n
  \           <span class=\"nd\">@app</span><span class=\"o\">.</span><span class=\"n\">command</span><span
  class=\"p\">()</span>\n            <span class=\"k\">def</span> <span class=\"nf\">list</span><span
  class=\"p\">(</span>\n                <span class=\"nb\">map</span><span class=\"p\">:</span>
  <span class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">\"title\"</span><span
  class=\"p\">,</span>\n                <span class=\"nb\">filter</span><span class=\"p\">:</span>
  <span class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">\"True\"</span><span
  class=\"p\">,</span>\n                <span class=\"n\">sort</span><span class=\"p\">:</span>
  <span class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">\"True\"</span><span
  class=\"p\">,</span>\n                <span class=\"n\">head</span><span class=\"p\">:</span>
  <span class=\"n\">Optional</span><span class=\"p\">[</span><span class=\"nb\">int</span><span
  class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"kc\">None</span><span
  class=\"p\">,</span>\n                <span class=\"n\">tail</span><span class=\"p\">:</span>
  <span class=\"n\">Optional</span><span class=\"p\">[</span><span class=\"nb\">int</span><span
  class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"kc\">None</span><span
  class=\"p\">,</span>\n                <span class=\"n\">include_empty</span><span
  class=\"p\">:</span> <span class=\"nb\">bool</span> <span class=\"o\">=</span> <span
  class=\"kc\">False</span><span class=\"p\">,</span>\n                <span class=\"n\">reverse</span><span
  class=\"p\">:</span> <span class=\"nb\">bool</span> <span class=\"o\">=</span> <span
  class=\"kc\">False</span><span class=\"p\">,</span>\n                <span class=\"n\">use_pager</span><span
  class=\"p\">:</span> <span class=\"nb\">bool</span> <span class=\"o\">=</span> <span
  class=\"n\">typer</span><span class=\"o\">.</span><span class=\"n\">Option</span><span
  class=\"p\">(</span><span class=\"kc\">True</span><span class=\"p\">,</span> <span
  class=\"s2\">\"--pager\"</span><span class=\"p\">,</span> <span class=\"s2\">\"--no-pager\"</span><span
  class=\"p\">),</span>\n            <span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
  <span class=\"kc\">None</span><span class=\"p\">:</span>\n<span class=\"w\">                </span><span
  class=\"sd\">\"\"\"</span>\n<span class=\"sd\">                Provides a way run
  markatas, map, filter, and sort from the</span>\n<span class=\"sd\">                command
  line.  I personally use this more often than the build</span>\n<span class=\"sd\">
  \               command while I am writing on a site with a large number of</span>\n<span
  class=\"sd\">                posts on it.  It makes slicing in by `templatekey`,
  `tag`, or</span>\n<span class=\"sd\">                `date` much easier.</span>\n\n<span
  class=\"sd\">                # default list</span>\n\n<span class=\"sd\">                By
  default `markata list` will list all titles in a pager, for all posts</span>\n<span
  class=\"sd\">                being loaded by markata.</span>\n\n<span class=\"sd\">
  \               ``` bash</span>\n<span class=\"sd\">                markata list</span>\n<span
  class=\"sd\">                ```</span>\n\n<span class=\"sd\">                #
  Skip the pager</span>\n\n<span class=\"sd\">                Markata uses rich for
  its pager, it's pretty smart about when to</span>\n<span class=\"sd\">                use
  the pager or pass text to the next thing in the pipeline,</span>\n<span class=\"sd\">
  \               but if you don't want to run a pager you can pass  `--no-pager`</span>\n\n<span
  class=\"sd\">                ``` bash</span>\n<span class=\"sd\">                markata
  list --no-pager</span>\n<span class=\"sd\">                ```</span>\n\n<span class=\"sd\">
  \               # List other attributes</span>\n\n<span class=\"sd\">                You
  can list any other attribute tied to your posts.  These are</span>\n<span class=\"sd\">
  \               added through either your yaml frontmatter at the start of your</span>\n<span
  class=\"sd\">                post, or through the use of a plugin.</span>\n\n\n<span
  class=\"sd\">                ``` bash</span>\n<span class=\"sd\">                #
  the filepath of the post</span>\n<span class=\"sd\">                markata list
  --map path</span>\n\n<span class=\"sd\">                # the slug of the post (where
  it will show up on the site)</span>\n<span class=\"sd\">                markata
  list --map slug</span>\n\n<span class=\"sd\">                # the date of the post</span>\n<span
  class=\"sd\">                markata list --map date</span>\n\n<span class=\"sd\">
  \               # the full raw content of the post</span>\n<span class=\"sd\">                markata
  list --map content</span>\n<span class=\"sd\">                ```</span>\n\n<span
  class=\"sd\">                # List more than one attribute</span>\n\n<span class=\"sd\">
  \               You can create new attributes as you map to echo out by</span>\n<span
  class=\"sd\">                combining existing attributes.</span>\n\n<span class=\"sd\">
  \               ``` bash</span>\n<span class=\"sd\">                markata list
  --map 'title + \" , \" + slug'</span>\n<span class=\"sd\">                ```</span>\n\n<span
  class=\"sd\">                # Using Python objects as map</span>\n\n<span class=\"sd\">
  \               You can access attributes of each post attribute that you map</span>\n<span
  class=\"sd\">                over.  For instance on my blog, each post has a date
  that is a</span>\n<span class=\"sd\">                datetime object.  I can ask
  each post for its `date.year`</span>\n\n<span class=\"sd\">                ``` bash</span>\n<span
  class=\"sd\">                markata list --map date.year</span>\n\n<span class=\"sd\">
  \               # combining this with title</span>\n<span class=\"sd\">                markata
  list --map 'str(date.year) + \",\" + title'</span>\n<span class=\"sd\">                ```</span>\n\n<span
  class=\"sd\">                # Filtering posts</span>\n\n<span class=\"sd\">                Posts
  are filtered with python syntax, you will have all</span>\n<span class=\"sd\">                attributes
  tied to your posts available to filter with.</span>\n\n<span class=\"sd\">                ```
  bash</span>\n<span class=\"sd\">                markata list --filter \"'__' not
  in title\"</span>\n<span class=\"sd\">                ```</span>\n\n<span class=\"sd\">
  \               # Filtering by dates</span>\n\n<span class=\"sd\">                If
  your site has dates tied to your posts you can filter by</span>\n<span class=\"sd\">
  \               date.  On my blog this makes a ton of sense and is quite useful.</span>\n<span
  class=\"sd\">                On the Markata docs though it doesn't really make much
  sense,</span>\n<span class=\"sd\">                since there really isn't the idea
  of a post date there.</span>\n\n<span class=\"sd\">                ``` bash</span>\n<span
  class=\"sd\">                # listing today's posts</span>\n<span class=\"sd\">
  \               markata list --filter \"date==today\"</span>\n\n<span class=\"sd\">
  \               # listing this year's posts</span>\n<span class=\"sd\">                markata
  list --filter \"date.year==today.year\"</span>\n<span class=\"sd\">                ```</span>\n\n<span
  class=\"sd\">                # Full Content Search</span>\n\n<span class=\"sd\">
  \               You can also search the full content of each post for specific</span>\n<span
  class=\"sd\">                words.</span>\n<span class=\"sd\">                ```
  bash</span>\n\n<span class=\"sd\">                markata list --filter \"'python'
  in content\"</span>\n<span class=\"sd\">                ```</span>\n\n<span class=\"sd\">
  \               # Filtering by frontmatter data</span>\n\n<span class=\"sd\">                I
  use a templateKey on my personal blog to determine which</span>\n<span class=\"sd\">
  \               template to render the page with.  I can fitler my posts by a</span>\n<span
  class=\"sd\">                `til` (today i learned) key.</span>\n\n<span class=\"sd\">
  \               ``` bash</span>\n<span class=\"sd\">                markata list
  --filter \"templateKey=='til'\"</span>\n<span class=\"sd\">                ```</span>\n\n<span
  class=\"sd\">                # Combining filters</span>\n\n<span class=\"sd\">                Filters
  can be combined together quite like maps can, it's all</span>\n<span class=\"sd\">
  \               just python syntax.</span>\n\n<span class=\"sd\">                ```
  bash</span>\n<span class=\"sd\">                markata list --filter \"templateKey=='til'
  and date == today\"</span>\n<span class=\"sd\">                ```</span>\n\n<span
  class=\"sd\">                # Sorting posts</span>\n\n<span class=\"sd\">                Posts
  can be sorted by attributes on your post, and they can</span>\n<span class=\"sd\">
  \               even be reversed.</span>\n\n<span class=\"sd\">                ```
  bash</span>\n<span class=\"sd\">                markta list --sort date</span>\n<span
  class=\"sd\">                markta list --sort date --reverse</span>\n<span class=\"sd\">
  \               ```</span>\n\n<span class=\"sd\">                # Putting it all
  together</span>\n\n<span class=\"sd\">                The real power of all this
  comes when you combine them all into</span>\n<span class=\"sd\">                lists
  that work for you and your workflow.  This really makes</span>\n<span class=\"sd\">
  \               working on larger projects so much easier to find things.</span>\n\n\n<span
  class=\"sd\">                # Making a fuzzy picker for your posts</span>\n\n<span
  class=\"sd\">                Here is a bash command to open an fzf picker for todays
  posts,</span>\n<span class=\"sd\">                then open it in your `$EDITOR`</span>\n\n<span
  class=\"sd\">                ``` bash</span>\n<span class=\"sd\">                markata
  list \\</span>\n<span class=\"sd\">                        --map path\\</span>\n<span
  class=\"sd\">                        --filter 'date==today'\\</span>\n<span class=\"sd\">
  \                       --sort date\\</span>\n<span class=\"sd\">                        --reverse
  |\\</span>\n<span class=\"sd\">                        fzf --preview 'bat --color
  always {}' |\\</span>\n<span class=\"sd\">                        xargs -I {} $EDITOR
  {}</span>\n<span class=\"sd\">                ```</span>\n\n<span class=\"sd\">
  \               # Combining wtih nvim Telescope</span>\n\n<span class=\"sd\">                Here
  is the same command setup as a Telescope picker for neovim.</span>\n\n<span class=\"sd\">
  \               ``` vim</span>\n<span class=\"sd\">                nnoremap &lt;leader&gt;et
  &lt;cmd&gt;Telescope find_files find_command=markata,list,--map,path,--filter,date==today&lt;cr&gt;</span>\n<span
  class=\"sd\">                ```</span>\n\n<span class=\"sd\">                If
  you have another way to open posts in your editor with</span>\n<span class=\"sd\">
  \               `markata list` I would love to accept a PR to add it to the</span>\n<span
  class=\"sd\">                examples here.</span>\n<span class=\"sd\">                \"\"\"</span>\n\n
  \               <span class=\"n\">markata</span><span class=\"o\">.</span><span
  class=\"n\">console</span><span class=\"o\">.</span><span class=\"n\">quiet</span>
  <span class=\"o\">=</span> <span class=\"kc\">True</span>\n\n                <span
  class=\"n\">tail</span> <span class=\"o\">=</span> <span class=\"o\">-</span><span
  class=\"n\">tail</span> <span class=\"k\">if</span> <span class=\"n\">tail</span>
  <span class=\"k\">else</span> <span class=\"n\">tail</span>\n                <span
  class=\"n\">filtered</span> <span class=\"o\">=</span> <span class=\"n\">markata</span><span
  class=\"o\">.</span><span class=\"n\">map</span><span class=\"p\">(</span><span
  class=\"nb\">map</span><span class=\"p\">,</span> <span class=\"nb\">filter</span><span
  class=\"p\">,</span> <span class=\"n\">sort</span><span class=\"p\">)</span>\n                <span
  class=\"k\">if</span> <span class=\"ow\">not</span> <span class=\"n\">include_empty</span><span
  class=\"p\">:</span>\n                    <span class=\"n\">filtered</span> <span
  class=\"o\">=</span> <span class=\"p\">[</span><span class=\"n\">a</span> <span
  class=\"k\">for</span> <span class=\"n\">a</span> <span class=\"ow\">in</span> <span
  class=\"n\">filtered</span> <span class=\"k\">if</span> <span class=\"n\">a</span>
  <span class=\"o\">!=</span> <span class=\"s2\">\"\"</span><span class=\"p\">]</span>\n
  \               <span class=\"n\">filtered</span> <span class=\"o\">=</span> <span
  class=\"n\">filtered</span><span class=\"p\">[</span><span class=\"n\">tail</span><span
  class=\"p\">:</span><span class=\"n\">head</span><span class=\"p\">]</span>\n                <span
  class=\"k\">if</span> <span class=\"n\">reverse</span><span class=\"p\">:</span>\n
  \                   <span class=\"n\">filtered</span> <span class=\"o\">=</span>
  <span class=\"nb\">reversed</span><span class=\"p\">(</span><span class=\"n\">filtered</span><span
  class=\"p\">)</span>\n\n                <span class=\"n\">markata</span><span class=\"o\">.</span><span
  class=\"n\">console</span><span class=\"o\">.</span><span class=\"n\">quiet</span>
  <span class=\"o\">=</span> <span class=\"kc\">False</span>\n                <span
  class=\"k\">if</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
  class=\"n\">console</span><span class=\"o\">.</span><span class=\"n\">is_terminal</span>
  <span class=\"ow\">and</span> <span class=\"n\">use_pager</span><span class=\"p\">:</span>\n
  \                   <span class=\"k\">with</span> <span class=\"n\">markata</span><span
  class=\"o\">.</span><span class=\"n\">console</span><span class=\"o\">.</span><span
  class=\"n\">pager</span><span class=\"p\">():</span>\n                        <span
  class=\"k\">for</span> <span class=\"n\">a</span> <span class=\"ow\">in</span> <span
  class=\"n\">filtered</span><span class=\"p\">:</span>\n                            <span
  class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">console</span><span
  class=\"o\">.</span><span class=\"n\">print</span><span class=\"p\">(</span><span
  class=\"n\">a</span><span class=\"p\">,</span> <span class=\"n\">style</span><span
  class=\"o\">=</span><span class=\"s2\">\"purple\"</span><span class=\"p\">)</span>\n
  \               <span class=\"k\">else</span><span class=\"p\">:</span>\n                    <span
  class=\"k\">for</span> <span class=\"n\">a</span> <span class=\"ow\">in</span> <span
  class=\"n\">filtered</span><span class=\"p\">:</span>\n                        <span
  class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">console</span><span
  class=\"o\">.</span><span class=\"n\">print</span><span class=\"p\">(</span><span
  class=\"n\">a</span><span class=\"p\">)</span>\n\n            <span class=\"nd\">@app</span><span
  class=\"o\">.</span><span class=\"n\">command</span><span class=\"p\">()</span>\n
  \           <span class=\"k\">def</span> <span class=\"nf\">clean</span><span class=\"p\">(</span>\n
  \               <span class=\"n\">quiet</span><span class=\"p\">:</span> <span class=\"nb\">bool</span>
  <span class=\"o\">=</span> <span class=\"n\">typer</span><span class=\"o\">.</span><span
  class=\"n\">Option</span><span class=\"p\">(</span>\n                    <span class=\"kc\">False</span><span
  class=\"p\">,</span>\n                    <span class=\"s2\">\"--quiet\"</span><span
  class=\"p\">,</span>\n                    <span class=\"s2\">\"-q\"</span><span
  class=\"p\">,</span>\n                <span class=\"p\">),</span>\n                <span
  class=\"n\">dry_run</span><span class=\"p\">:</span> <span class=\"nb\">bool</span>
  <span class=\"o\">=</span> <span class=\"n\">typer</span><span class=\"o\">.</span><span
  class=\"n\">Option</span><span class=\"p\">(</span>\n                    <span class=\"kc\">False</span><span
  class=\"p\">,</span>\n                    <span class=\"s2\">\"--dry-run\"</span><span
  class=\"p\">,</span>\n                <span class=\"p\">),</span>\n            <span
  class=\"p\">):</span>\n<span class=\"w\">                </span><span class=\"sd\">\"\"\"</span>\n<span
  class=\"sd\">                Cleans up output generated by markata including both
  the output_dir and</span>\n<span class=\"sd\">                the .markata_cache.</span>\n\n<span
  class=\"sd\">                # Dry Run</span>\n\n<span class=\"sd\">                You
  can run with `--dry-run` to see what markata is about to do.</span>\n\n<span class=\"sd\">
  \               ``` bash</span>\n<span class=\"sd\">                markata clean
  --dry-run</span>\n<span class=\"sd\">                [09:42:37] [DRYRUN] removing
  outptut directory: markout base_cli.py:371</span>\n<span class=\"sd\">                           [DRYRUN]
  removing cache directory: .markata.cache base_cli.py:377</span>\n\n<span class=\"sd\">
  \               ```</span>\n\n<span class=\"sd\">                # Running clean</span>\n\n<span
  class=\"sd\">                Running markata clean will fully delete all of the
  directories created</span>\n<span class=\"sd\">                by markata.</span>\n\n<span
  class=\"sd\">                ``` bash</span>\n<span class=\"sd\">                markata
  clean</span>\n<span class=\"sd\">                [09:53:04]  removing outptut directory:
  markout base_cli.py:394</span>\n<span class=\"sd\">                            removing
  cache directory: .markata.cache base_cli.py:405</span>\n<span class=\"sd\">                ```</span>\n\n<span
  class=\"sd\">                # Running Quietly</span>\n\n<span class=\"sd\">                Running
  with `--quiet` will remove all of the directories created by</span>\n<span class=\"sd\">
  \               markata without announcing what it is doing.</span>\n\n<span class=\"sd\">
  \               ``` bash</span>\n<span class=\"sd\">                markata clean
  --quiet</span>\n<span class=\"sd\">                ```</span>\n<span class=\"sd\">
  \               \"\"\"</span>\n                <span class=\"n\">_clean</span><span
  class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">=</span><span
  class=\"n\">markata</span><span class=\"p\">,</span> <span class=\"n\">quiet</span><span
  class=\"o\">=</span><span class=\"n\">quiet</span><span class=\"p\">,</span> <span
  class=\"n\">dry_run</span><span class=\"o\">=</span><span class=\"n\">dry_run</span><span
  class=\"p\">)</span>\n</pre></div>\n\n</pre>\n<p>!! function </p><h2 class=\"admonition-title\"
  id=\"_clean\" style=\"margin:0;padding:.5rem 1rem;\">_clean <em class=\"small\">function</em></h2>\n<div
  class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"_clean
  <em class=\"small\">source</em>\"</p>\n</div>\n<pre class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button
  class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
  class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
  v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
  c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
  c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
  h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
  M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
  v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
  M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
  <span class=\"nf\">_clean</span><span class=\"p\">(</span><span class=\"n\">markata</span><span
  class=\"p\">,</span> <span class=\"n\">quiet</span><span class=\"p\">:</span> <span
  class=\"nb\">bool</span> <span class=\"o\">=</span> <span class=\"kc\">False</span><span
  class=\"p\">,</span> <span class=\"n\">dry_run</span><span class=\"p\">:</span>
  <span class=\"nb\">bool</span> <span class=\"o\">=</span> <span class=\"kc\">False</span><span
  class=\"p\">):</span>\n            <span class=\"k\">if</span> <span class=\"n\">quiet</span><span
  class=\"p\">:</span>\n                <span class=\"n\">markata</span><span class=\"o\">.</span><span
  class=\"n\">console</span><span class=\"o\">.</span><span class=\"n\">quiet</span>
  <span class=\"o\">=</span> <span class=\"kc\">True</span>\n\n            <span class=\"n\">markata</span><span
  class=\"o\">.</span><span class=\"n\">console</span><span class=\"o\">.</span><span
  class=\"n\">log</span><span class=\"p\">(</span>\n                <span class=\"sa\">f</span><span
  class=\"s1\">'</span><span class=\"si\">{</span><span class=\"s2\">\"[DRYRUN]\"</span><span
  class=\"w\"> </span><span class=\"k\">if</span><span class=\"w\"> </span><span class=\"n\">dry_run</span><span
  class=\"w\"> </span><span class=\"k\">else</span><span class=\"w\"> </span><span
  class=\"s2\">\"\"</span><span class=\"si\">}</span><span class=\"s1\">'</span>\n
  \               <span class=\"sa\">f</span><span class=\"s2\">\"removing outptut
  directory: </span><span class=\"si\">{</span><span class=\"n\">markata</span><span
  class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
  class=\"n\">output_dir</span><span class=\"si\">}</span><span class=\"s2\">\"</span><span
  class=\"p\">,</span>\n            <span class=\"p\">)</span>\n            <span
  class=\"k\">if</span> <span class=\"ow\">not</span> <span class=\"n\">dry_run</span><span
  class=\"p\">:</span>\n                <span class=\"k\">try</span><span class=\"p\">:</span>\n
  \                   <span class=\"n\">shutil</span><span class=\"o\">.</span><span
  class=\"n\">rmtree</span><span class=\"p\">(</span><span class=\"nb\">str</span><span
  class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
  class=\"n\">config</span><span class=\"o\">.</span><span class=\"n\">output_dir</span><span
  class=\"p\">))</span>\n                <span class=\"k\">except</span> <span class=\"ne\">FileNotFoundError</span><span
  class=\"p\">:</span>\n                    <span class=\"n\">warnings</span><span
  class=\"o\">.</span><span class=\"n\">warn</span><span class=\"p\">(</span>\n                        <span
  class=\"sa\">f</span><span class=\"s2\">\"output directory: </span><span class=\"si\">{</span><span
  class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
  class=\"o\">.</span><span class=\"n\">output_dir</span><span class=\"si\">}</span><span
  class=\"s2\"> does not exist\"</span><span class=\"p\">,</span>\n                    <span
  class=\"p\">)</span>\n\n            <span class=\"n\">markata</span><span class=\"o\">.</span><span
  class=\"n\">console</span><span class=\"o\">.</span><span class=\"n\">log</span><span
  class=\"p\">(</span>\n                <span class=\"sa\">f</span><span class=\"s1\">'</span><span
  class=\"si\">{</span><span class=\"s2\">\"[DRYRUN]\"</span><span class=\"w\"> </span><span
  class=\"k\">if</span><span class=\"w\"> </span><span class=\"n\">dry_run</span><span
  class=\"w\"> </span><span class=\"k\">else</span><span class=\"w\"> </span><span
  class=\"s2\">\"\"</span><span class=\"si\">}</span><span class=\"s1\"> removing
  cache directory: .markata.cache'</span><span class=\"p\">,</span>\n            <span
  class=\"p\">)</span>\n            <span class=\"k\">if</span> <span class=\"ow\">not</span>
  <span class=\"n\">dry_run</span><span class=\"p\">:</span>\n                <span
  class=\"k\">try</span><span class=\"p\">:</span>\n                    <span class=\"n\">shutil</span><span
  class=\"o\">.</span><span class=\"n\">rmtree</span><span class=\"p\">(</span><span
  class=\"s2\">\".markata.cache\"</span><span class=\"p\">)</span>\n                <span
  class=\"k\">except</span> <span class=\"ne\">FileNotFoundError</span><span class=\"p\">:</span>\n
  \                   <span class=\"n\">warnings</span><span class=\"o\">.</span><span
  class=\"n\">warn</span><span class=\"p\">(</span><span class=\"s2\">\"cache directory:
  .markata.cache does not exist\"</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n<p>!!
  function </p><h2 class=\"admonition-title\" id=\"pdb_run\" style=\"margin:0;padding:.5rem
  1rem;\">pdb_run <em class=\"small\">function</em></h2>\nWraps a function call with
  a post_mortem pdb debugger.\n<div class=\"admonition source is-collapsible collapsible-open\">\n<p
  class=\"admonition-title\">\"pdb_run <em class=\"small\">source</em>\"</p>\n</div>\n<pre
  class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
  class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
  v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
  c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
  c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
  h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
  M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
  v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
  M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
  <span class=\"nf\">pdb_run</span><span class=\"p\">(</span><span class=\"n\">func</span><span
  class=\"p\">:</span> <span class=\"n\">Callable</span><span class=\"p\">)</span>
  <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n<span
  class=\"w\">            </span><span class=\"sd\">\"\"\"</span>\n<span class=\"sd\">
  \           Wraps a function call with a post_mortem pdb debugger.</span>\n<span
  class=\"sd\">            \"\"\"</span>\n            <span class=\"k\">try</span><span
  class=\"p\">:</span>\n                <span class=\"n\">func</span><span class=\"p\">()</span>\n
  \           <span class=\"k\">except</span> <span class=\"ne\">Exception</span><span
  class=\"p\">:</span>\n                <span class=\"n\">extype</span><span class=\"p\">,</span>
  <span class=\"n\">value</span><span class=\"p\">,</span> <span class=\"n\">tb</span>
  <span class=\"o\">=</span> <span class=\"n\">sys</span><span class=\"o\">.</span><span
  class=\"n\">exc_info</span><span class=\"p\">()</span>\n                <span class=\"n\">traceback</span><span
  class=\"o\">.</span><span class=\"n\">print_exc</span><span class=\"p\">()</span>\n
  \               <span class=\"n\">pdb</span><span class=\"o\">.</span><span class=\"n\">post_mortem</span><span
  class=\"p\">(</span><span class=\"n\">tb</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n<p>!!
  function </p><h2 class=\"admonition-title\" id=\"tui\" style=\"margin:0;padding:.5rem
  1rem;\">tui <em class=\"small\">function</em></h2>\n<div class=\"admonition source
  is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"tui <em class=\"small\">source</em>\"</p>\n</div>\n<pre
  class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
  class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
  v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
  c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
  c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
  h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
  M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
  v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
  M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
  <span class=\"nf\">tui</span><span class=\"p\">(</span><span class=\"n\">ctx</span><span
  class=\"p\">:</span> <span class=\"n\">typer</span><span class=\"o\">.</span><span
  class=\"n\">Context</span><span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
  <span class=\"kc\">None</span><span class=\"p\">:</span>\n                <span
  class=\"k\">try</span><span class=\"p\">:</span>\n                    <span class=\"kn\">from</span>
  <span class=\"nn\">trogon</span> <span class=\"kn\">import</span> <span class=\"n\">Trogon</span>\n
  \                   <span class=\"kn\">from</span> <span class=\"nn\">typer.main</span>
  <span class=\"kn\">import</span> <span class=\"n\">get_group</span>\n                <span
  class=\"k\">except</span> <span class=\"ne\">ImportError</span><span class=\"p\">:</span>\n
  \                   <span class=\"n\">typer</span><span class=\"o\">.</span><span
  class=\"n\">echo</span><span class=\"p\">(</span><span class=\"s2\">\"trogon not
  installed\"</span><span class=\"p\">)</span>\n                    <span class=\"n\">typer</span><span
  class=\"o\">.</span><span class=\"n\">echo</span><span class=\"p\">(</span>\n                        <span
  class=\"s2\">\"install markata with optional tui group to use tui `pip install 'markata[tui]'`\"</span>\n
  \                   <span class=\"p\">)</span>\n                    <span class=\"k\">return</span>\n\n
  \               <span class=\"n\">Trogon</span><span class=\"p\">(</span><span class=\"n\">get_group</span><span
  class=\"p\">(</span><span class=\"n\">app</span><span class=\"p\">),</span> <span
  class=\"n\">click_context</span><span class=\"o\">=</span><span class=\"n\">ctx</span><span
  class=\"p\">)</span><span class=\"o\">.</span><span class=\"n\">run</span><span
  class=\"p\">()</span>\n</pre></div>\n\n</pre>\n<p>!! function </p><h2 class=\"admonition-title\"
  id=\"plugins\" style=\"margin:0;padding:.5rem 1rem;\">plugins <em class=\"small\">function</em></h2>\ncreate
  new things from templates\n<div class=\"admonition source is-collapsible collapsible-open\">\n<p
  class=\"admonition-title\">\"plugins <em class=\"small\">source</em>\"</p>\n</div>\n<pre
  class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
  class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
  v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
  c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
  c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
  h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
  M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
  v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
  M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
  <span class=\"nf\">plugins</span><span class=\"p\">():</span>\n                <span
  class=\"s2\">\"create new things from templates\"</span>\n</pre></div>\n\n</pre>\n<p>!!
  function </p><h2 class=\"admonition-title\" id=\"config\" style=\"margin:0;padding:.5rem
  1rem;\">config <em class=\"small\">function</em></h2>\nconfiguration management\n<div
  class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"config
  <em class=\"small\">source</em>\"</p>\n</div>\n<pre class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button
  class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
  class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
  v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
  c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
  c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
  h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
  M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
  v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
  M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
  <span class=\"nf\">config</span><span class=\"p\">():</span>\n                <span
  class=\"s2\">\"configuration management\"</span>\n</pre></div>\n\n</pre>\n<p>!!
  function </p><h2 class=\"admonition-title\" id=\"show\" style=\"margin:0;padding:.5rem
  1rem;\">show <em class=\"small\">function</em></h2>\n<div class=\"admonition source
  is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"show <em class=\"small\">source</em>\"</p>\n</div>\n<pre
  class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
  class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
  v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
  c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
  c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
  h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
  M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
  v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
  M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
  <span class=\"nf\">show</span><span class=\"p\">(</span>\n                <span
  class=\"n\">verbose</span><span class=\"p\">:</span> <span class=\"nb\">bool</span>
  <span class=\"o\">=</span> <span class=\"n\">typer</span><span class=\"o\">.</span><span
  class=\"n\">Option</span><span class=\"p\">(</span>\n                    <span class=\"kc\">False</span><span
  class=\"p\">,</span>\n                    <span class=\"s2\">\"--verbose\"</span><span
  class=\"p\">,</span>\n                    <span class=\"s2\">\"-v\"</span><span
  class=\"p\">,</span>\n                <span class=\"p\">),</span>\n            <span
  class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
  class=\"p\">:</span>\n                <span class=\"k\">if</span> <span class=\"n\">verbose</span><span
  class=\"p\">:</span>\n                    <span class=\"n\">markata</span><span
  class=\"o\">.</span><span class=\"n\">console</span><span class=\"o\">.</span><span
  class=\"n\">quiet</span> <span class=\"o\">=</span> <span class=\"kc\">False</span>\n
  \               <span class=\"k\">else</span><span class=\"p\">:</span>\n                    <span
  class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">console</span><span
  class=\"o\">.</span><span class=\"n\">quiet</span> <span class=\"o\">=</span> <span
  class=\"kc\">True</span>\n                <span class=\"n\">rich_print</span><span
  class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
  class=\"n\">config</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n<p>!!
  function </p><h2 class=\"admonition-title\" id=\"generate\" style=\"margin:0;padding:.5rem
  1rem;\">generate <em class=\"small\">function</em></h2>\n<div class=\"admonition
  source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"generate
  <em class=\"small\">source</em>\"</p>\n</div>\n<pre class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button
  class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
  class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
  v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
  c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
  c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
  h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
  M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
  v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
  M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
  <span class=\"nf\">generate</span><span class=\"p\">(</span>\n                <span
  class=\"n\">verbose</span><span class=\"p\">:</span> <span class=\"nb\">bool</span>
  <span class=\"o\">=</span> <span class=\"n\">typer</span><span class=\"o\">.</span><span
  class=\"n\">Option</span><span class=\"p\">(</span>\n                    <span class=\"kc\">False</span><span
  class=\"p\">,</span>\n                    <span class=\"s2\">\"--verbose\"</span><span
  class=\"p\">,</span>\n                    <span class=\"s2\">\"-v\"</span><span
  class=\"p\">,</span>\n                <span class=\"p\">),</span>\n            <span
  class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
  class=\"p\">:</span>\n                <span class=\"k\">if</span> <span class=\"n\">verbose</span><span
  class=\"p\">:</span>\n                    <span class=\"n\">markata</span><span
  class=\"o\">.</span><span class=\"n\">console</span><span class=\"o\">.</span><span
  class=\"n\">quiet</span> <span class=\"o\">=</span> <span class=\"kc\">False</span>\n
  \               <span class=\"k\">else</span><span class=\"p\">:</span>\n                    <span
  class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">console</span><span
  class=\"o\">.</span><span class=\"n\">quiet</span> <span class=\"o\">=</span> <span
  class=\"kc\">True</span>\n\n                <span class=\"n\">rich_print</span><span
  class=\"p\">(</span><span class=\"n\">toml</span><span class=\"o\">.</span><span
  class=\"n\">dumps</span><span class=\"p\">(</span><span class=\"n\">json</span><span
  class=\"o\">.</span><span class=\"n\">loads</span><span class=\"p\">(</span><span
  class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">config</span><span
  class=\"o\">.</span><span class=\"n\">json</span><span class=\"p\">())))</span>\n</pre></div>\n\n</pre>\n<p>!!
  function </p><h2 class=\"admonition-title\" id=\"get\" style=\"margin:0;padding:.5rem
  1rem;\">get <em class=\"small\">function</em></h2>\n<div class=\"admonition source
  is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"get <em class=\"small\">source</em>\"</p>\n</div>\n<pre
  class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
  class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
  v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
  c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
  c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
  h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
  M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
  v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
  M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
  <span class=\"nf\">get</span><span class=\"p\">(</span><span class=\"n\">key</span><span
  class=\"p\">:</span> <span class=\"nb\">str</span><span class=\"p\">)</span> <span
  class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span class=\"p\">:</span>\n
  \               <span class=\"n\">keys</span> <span class=\"o\">=</span> <span class=\"n\">key</span><span
  class=\"o\">.</span><span class=\"n\">split</span><span class=\"p\">(</span><span
  class=\"s2\">\".\"</span><span class=\"p\">)</span>\n                <span class=\"n\">markata</span><span
  class=\"o\">.</span><span class=\"n\">console</span><span class=\"o\">.</span><span
  class=\"n\">quiet</span> <span class=\"o\">=</span> <span class=\"kc\">True</span>\n
  \               <span class=\"n\">keys_processed</span> <span class=\"o\">=</span>
  <span class=\"s2\">\"\"</span>\n                <span class=\"n\">value</span> <span
  class=\"o\">=</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
  class=\"n\">config</span>\n                <span class=\"n\">na</span> <span class=\"o\">=</span>
  <span class=\"n\">Literal</span><span class=\"p\">[</span><span class=\"s2\">\"na\"</span><span
  class=\"p\">]</span>\n                <span class=\"k\">for</span> <span class=\"n\">key</span>
  <span class=\"ow\">in</span> <span class=\"n\">keys</span><span class=\"p\">:</span>\n
  \                   <span class=\"n\">value</span> <span class=\"o\">=</span> <span
  class=\"nb\">getattr</span><span class=\"p\">(</span><span class=\"n\">value</span><span
  class=\"p\">,</span> <span class=\"n\">key</span><span class=\"p\">,</span> <span
  class=\"n\">na</span><span class=\"p\">)</span>\n                    <span class=\"n\">keys_processed</span>
  <span class=\"o\">=</span> <span class=\"sa\">f</span><span class=\"s2\">\"</span><span
  class=\"si\">{</span><span class=\"n\">keys_processed</span><span class=\"si\">}</span><span
  class=\"s2\">.</span><span class=\"si\">{</span><span class=\"n\">key</span><span
  class=\"si\">}</span><span class=\"s2\">\"</span><span class=\"o\">.</span><span
  class=\"n\">strip</span><span class=\"p\">(</span><span class=\"s2\">\".\"</span><span
  class=\"p\">)</span>\n                    <span class=\"k\">if</span> <span class=\"n\">value</span>
  <span class=\"ow\">is</span> <span class=\"n\">na</span><span class=\"p\">:</span>\n
  \                       <span class=\"n\">rich_print</span><span class=\"p\">(</span><span
  class=\"sa\">f</span><span class=\"s2\">\"</span><span class=\"si\">{</span><span
  class=\"n\">keys_processed</span><span class=\"si\">}</span><span class=\"s2\">
  not found\"</span><span class=\"p\">)</span>\n                        <span class=\"n\">sys</span><span
  class=\"o\">.</span><span class=\"n\">exit</span><span class=\"p\">(</span><span
  class=\"mi\">1</span><span class=\"p\">)</span>\n\n                <span class=\"n\">rich_print</span><span
  class=\"p\">(</span><span class=\"n\">value</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n<p>!!
  function </p><h2 class=\"admonition-title\" id=\"new\" style=\"margin:0;padding:.5rem
  1rem;\">new <em class=\"small\">function</em></h2>\ncreate new things from templates\n<div
  class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"new
  <em class=\"small\">source</em>\"</p>\n</div>\n<pre class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button
  class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
  class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
  v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
  c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
  c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
  h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
  M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
  v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
  M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
  <span class=\"nf\">new</span><span class=\"p\">():</span>\n                <span
  class=\"s2\">\"create new things from templates\"</span>\n</pre></div>\n\n</pre>\n<p>!!
  function </p><h2 class=\"admonition-title\" id=\"blog\" style=\"margin:0;padding:.5rem
  1rem;\">blog <em class=\"small\">function</em></h2>\nCreate a new blog from using
  the template from\n<a href=\"https://github.com/WaylonWalker/markata-blog-starter\">https://github.com/WaylonWalker/markata-blog-starter</a>.\n<div
  class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"blog
  <em class=\"small\">source</em>\"</p>\n</div>\n<pre class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button
  class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
  class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
  v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
  c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
  c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
  h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
  M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
  v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
  M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
  <span class=\"nf\">blog</span><span class=\"p\">(</span>\n                <span
  class=\"n\">directory</span><span class=\"p\">:</span> <span class=\"n\">Path</span>
  <span class=\"o\">=</span> <span class=\"n\">typer</span><span class=\"o\">.</span><span
  class=\"n\">Argument</span><span class=\"p\">(</span>\n                    <span
  class=\"o\">...</span><span class=\"p\">,</span>\n                    <span class=\"n\">help</span><span
  class=\"o\">=</span><span class=\"s2\">\"The directory to create the blog in.\"</span><span
  class=\"p\">,</span>\n                <span class=\"p\">),</span>\n            <span
  class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
  class=\"p\">:</span>\n<span class=\"w\">                </span><span class=\"sd\">\"\"\"</span>\n<span
  class=\"sd\">                Create a new blog from using the template from</span>\n<span
  class=\"sd\">                https://github.com/WaylonWalker/markata-blog-starter.</span>\n<span
  class=\"sd\">                \"\"\"</span>\n\n                <span class=\"kn\">from</span>
  <span class=\"nn\">copier</span> <span class=\"kn\">import</span> <span class=\"n\">run_copy</span>\n\n
  \               <span class=\"n\">typer</span><span class=\"o\">.</span><span class=\"n\">echo</span><span
  class=\"p\">(</span><span class=\"sa\">f</span><span class=\"s2\">\"creating a new
  project in </span><span class=\"si\">{</span><span class=\"n\">directory</span><span
  class=\"o\">.</span><span class=\"n\">absolute</span><span class=\"p\">()</span><span
  class=\"si\">}</span><span class=\"s2\">\"</span><span class=\"p\">)</span>\n                <span
  class=\"n\">url</span> <span class=\"o\">=</span> <span class=\"n\">markata</span><span
  class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
  class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">\"starters\"</span><span
  class=\"p\">,</span> <span class=\"p\">{})</span><span class=\"o\">.</span><span
  class=\"n\">get</span><span class=\"p\">(</span>\n                    <span class=\"s2\">\"blog\"</span><span
  class=\"p\">,</span>\n                    <span class=\"s2\">\"git+https://github.com/WaylonWalker/markata-blog-starter\"</span><span
  class=\"p\">,</span>\n                <span class=\"p\">)</span>\n                <span
  class=\"n\">run_copy</span><span class=\"p\">(</span><span class=\"n\">url</span><span
  class=\"p\">,</span> <span class=\"n\">directory</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n<p>!!
  function </p><h2 class=\"admonition-title\" id=\"post\" style=\"margin:0;padding:.5rem
  1rem;\">post <em class=\"small\">function</em></h2>\nCreate new blog post in the
  pages directory from the template at\n<a href=\"https://github.com/WaylonWalker/markata-post-template\">https://github.com/WaylonWalker/markata-post-template</a>.\n<div
  class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"post
  <em class=\"small\">source</em>\"</p>\n</div>\n<pre class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button
  class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
  class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
  v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
  c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
  c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
  h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
  M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
  v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
  M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
  <span class=\"nf\">post</span><span class=\"p\">()</span> <span class=\"o\">-&gt;</span>
  <span class=\"kc\">None</span><span class=\"p\">:</span>\n<span class=\"w\">                </span><span
  class=\"sd\">\"\"\"</span>\n<span class=\"sd\">                Create new blog post
  in the pages directory from the template at</span>\n<span class=\"sd\">                https://github.com/WaylonWalker/markata-post-template.</span>\n<span
  class=\"sd\">                \"\"\"</span>\n\n                <span class=\"nb\">print</span><span
  class=\"p\">(</span><span class=\"s2\">\"create a new post\"</span><span class=\"p\">)</span>\n
  \               <span class=\"kn\">from</span> <span class=\"nn\">copier</span>
  <span class=\"kn\">import</span> <span class=\"n\">run_copy</span>\n\n                <span
  class=\"n\">typer</span><span class=\"o\">.</span><span class=\"n\">echo</span><span
  class=\"p\">(</span><span class=\"sa\">f</span><span class=\"s2\">\"creating a new
  post in </span><span class=\"si\">{</span><span class=\"n\">Path</span><span class=\"p\">()</span><span
  class=\"o\">.</span><span class=\"n\">absolute</span><span class=\"p\">()</span><span
  class=\"si\">}</span><span class=\"s2\">/posts\"</span><span class=\"p\">)</span>\n
  \               <span class=\"n\">url</span> <span class=\"o\">=</span> <span class=\"n\">markata</span><span
  class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
  class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">\"starters\"</span><span
  class=\"p\">,</span> <span class=\"p\">{})</span><span class=\"o\">.</span><span
  class=\"n\">get</span><span class=\"p\">(</span>\n                    <span class=\"s2\">\"post\"</span><span
  class=\"p\">,</span>\n                    <span class=\"s2\">\"git+https://github.com/WaylonWalker/markata-post-template\"</span><span
  class=\"p\">,</span>\n                <span class=\"p\">)</span>\n                <span
  class=\"n\">run_copy</span><span class=\"p\">(</span><span class=\"n\">url</span><span
  class=\"p\">,</span> <span class=\"n\">Path</span><span class=\"p\">(</span><span
  class=\"s2\">\".\"</span><span class=\"p\">))</span>\n</pre></div>\n\n</pre>\n<p>!!
  function </p><h2 class=\"admonition-title\" id=\"plugin\" style=\"margin:0;padding:.5rem
  1rem;\">plugin <em class=\"small\">function</em></h2>\nCreate a new plugin using
  the template at\n<a href=\"https://github.com/WaylonWalker/markata-plugin-template\">https://github.com/WaylonWalker/markata-plugin-template</a>.\n<div
  class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"plugin
  <em class=\"small\">source</em>\"</p>\n</div>\n<pre class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button
  class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
  class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
  v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
  c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
  c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
  h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
  M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
  v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
  M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
  <span class=\"nf\">plugin</span><span class=\"p\">()</span> <span class=\"o\">-&gt;</span>
  <span class=\"kc\">None</span><span class=\"p\">:</span>\n<span class=\"w\">                </span><span
  class=\"sd\">\"\"\"</span>\n<span class=\"sd\">                Create a new plugin
  using the template at</span>\n<span class=\"sd\">                https://github.com/WaylonWalker/markata-plugin-template.</span>\n<span
  class=\"sd\">                \"\"\"</span>\n                <span class=\"kn\">from</span>
  <span class=\"nn\">copier</span> <span class=\"kn\">import</span> <span class=\"n\">run_copy</span>\n\n
  \               <span class=\"n\">typer</span><span class=\"o\">.</span><span class=\"n\">echo</span><span
  class=\"p\">(</span>\n                    <span class=\"sa\">f</span><span class=\"s2\">\"creating
  a new plugin in </span><span class=\"si\">{</span><span class=\"n\">Path</span><span
  class=\"p\">()</span><span class=\"o\">.</span><span class=\"n\">absolute</span><span
  class=\"p\">()</span><span class=\"si\">}</span><span class=\"s2\">\"</span>\n                    <span
  class=\"sa\">f</span><span class=\"s2\">\"/&lt;python-package-name&gt;/plugins\"</span><span
  class=\"p\">,</span>\n                <span class=\"p\">)</span>\n                <span
  class=\"n\">url</span> <span class=\"o\">=</span> <span class=\"n\">markata</span><span
  class=\"o\">.</span><span class=\"n\">config</span><span class=\"o\">.</span><span
  class=\"n\">get</span><span class=\"p\">(</span><span class=\"s2\">\"starters\"</span><span
  class=\"p\">,</span> <span class=\"p\">{})</span><span class=\"o\">.</span><span
  class=\"n\">get</span><span class=\"p\">(</span>\n                    <span class=\"s2\">\"post\"</span><span
  class=\"p\">,</span>\n                    <span class=\"s2\">\"git+https://github.com/WaylonWalker/markata-plugin-template\"</span><span
  class=\"p\">,</span>\n                <span class=\"p\">)</span>\n                <span
  class=\"n\">run_copy</span><span class=\"p\">(</span><span class=\"n\">url</span><span
  class=\"p\">,</span> <span class=\"n\">Path</span><span class=\"p\">(</span><span
  class=\"s2\">\".\"</span><span class=\"p\">))</span>\n</pre></div>\n\n</pre>\n<p>!!
  function </p><h2 class=\"admonition-title\" id=\"build\" style=\"margin:0;padding:.5rem
  1rem;\">build <em class=\"small\">function</em></h2>\nMarkata's primary way of building
  your site for production.\nBy default, running <code>markta build</code> will render
  your markdown to\nthe <code>./markout</code> directory.\n<pre><code>``` bash\nmarkata
  build\n```\n\nIf you are having an issue and want to pop immediately into a debugger\nupon
  failure you can pass the `--pdb` flag to the build command.\n\n``` bash\nmarkata
  build  --pdb\n```\n\nIf you do not like the way rich looks, or its suppressing tracebaks
  you\nwould like to remain visible you can use `--no-pretty`\n\n``` bash\nmarkata
  build --no-pretty\n```\n\nIf you need to run without any console logging pass in
  the\n`--quiet` flag.\n\n``` bash\nmarkata build --quiet\n```\n\n`markta build` will
  automatically run the pyinstrument profiler\nwhile building your site if you have
  pyinstrument installed.  It\nwill echo out your profile in the console as well as
  write it to\n`/_profile` on your built site. If you prefer not to run\npyinstrument
  profiling, even when it is installed you can pass\nin `--no-profile`\n\n``` bash\nmarkata
  build --no-profile\n```\n</code></pre>\n<div class=\"admonition source is-collapsible
  collapsible-open\">\n<p class=\"admonition-title\">\"build <em class=\"small\">source</em>\"</p>\n</div>\n<pre
  class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
  class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
  v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
  c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
  c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
  h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
  M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
  v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
  M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
  <span class=\"nf\">build</span><span class=\"p\">(</span>\n                <span
  class=\"n\">pretty</span><span class=\"p\">:</span> <span class=\"nb\">bool</span>
  <span class=\"o\">=</span> <span class=\"kc\">True</span><span class=\"p\">,</span>\n
  \               <span class=\"n\">quiet</span><span class=\"p\">:</span> <span class=\"nb\">bool</span>
  <span class=\"o\">=</span> <span class=\"n\">typer</span><span class=\"o\">.</span><span
  class=\"n\">Option</span><span class=\"p\">(</span>\n                    <span class=\"kc\">False</span><span
  class=\"p\">,</span>\n                    <span class=\"s2\">\"--quiet\"</span><span
  class=\"p\">,</span>\n                    <span class=\"s2\">\"-q\"</span><span
  class=\"p\">,</span>\n                <span class=\"p\">),</span>\n                <span
  class=\"n\">verbose</span><span class=\"p\">:</span> <span class=\"nb\">bool</span>
  <span class=\"o\">=</span> <span class=\"n\">typer</span><span class=\"o\">.</span><span
  class=\"n\">Option</span><span class=\"p\">(</span>\n                    <span class=\"kc\">False</span><span
  class=\"p\">,</span>\n                    <span class=\"s2\">\"--verbose\"</span><span
  class=\"p\">,</span>\n                    <span class=\"s2\">\"-v\"</span><span
  class=\"p\">,</span>\n                <span class=\"p\">),</span>\n                <span
  class=\"n\">should_pdb</span><span class=\"p\">:</span> <span class=\"nb\">bool</span>
  <span class=\"o\">=</span> <span class=\"n\">typer</span><span class=\"o\">.</span><span
  class=\"n\">Option</span><span class=\"p\">(</span>\n                    <span class=\"kc\">False</span><span
  class=\"p\">,</span>\n                    <span class=\"s2\">\"--pdb\"</span><span
  class=\"p\">,</span>\n                <span class=\"p\">),</span>\n                <span
  class=\"n\">profile</span><span class=\"p\">:</span> <span class=\"nb\">bool</span>
  <span class=\"o\">=</span> <span class=\"kc\">True</span><span class=\"p\">,</span>\n
  \           <span class=\"p\">)</span> <span class=\"o\">-&gt;</span> <span class=\"kc\">None</span><span
  class=\"p\">:</span>\n<span class=\"w\">                </span><span class=\"sd\">\"\"\"</span>\n<span
  class=\"sd\">                Markata's primary way of building your site for production.</span>\n<span
  class=\"sd\">                By default, running `markta build` will render your
  markdown to</span>\n<span class=\"sd\">                the `./markout` directory.</span>\n\n<span
  class=\"sd\">                ``` bash</span>\n<span class=\"sd\">                markata
  build</span>\n<span class=\"sd\">                ```</span>\n\n<span class=\"sd\">
  \               If you are having an issue and want to pop immediately into a debugger</span>\n<span
  class=\"sd\">                upon failure you can pass the `--pdb` flag to the build
  command.</span>\n\n<span class=\"sd\">                ``` bash</span>\n<span class=\"sd\">
  \               markata build  --pdb</span>\n<span class=\"sd\">                ```</span>\n\n<span
  class=\"sd\">                If you do not like the way rich looks, or its suppressing
  tracebaks you</span>\n<span class=\"sd\">                would like to remain visible
  you can use `--no-pretty`</span>\n\n<span class=\"sd\">                ``` bash</span>\n<span
  class=\"sd\">                markata build --no-pretty</span>\n<span class=\"sd\">
  \               ```</span>\n\n<span class=\"sd\">                If you need to
  run without any console logging pass in the</span>\n<span class=\"sd\">                `--quiet`
  flag.</span>\n\n<span class=\"sd\">                ``` bash</span>\n<span class=\"sd\">
  \               markata build --quiet</span>\n<span class=\"sd\">                ```</span>\n\n<span
  class=\"sd\">                `markta build` will automatically run the pyinstrument
  profiler</span>\n<span class=\"sd\">                while building your site if
  you have pyinstrument installed.  It</span>\n<span class=\"sd\">                will
  echo out your profile in the console as well as write it to</span>\n<span class=\"sd\">
  \               `/_profile` on your built site. If you prefer not to run</span>\n<span
  class=\"sd\">                pyinstrument profiling, even when it is installed you
  can pass</span>\n<span class=\"sd\">                in `--no-profile`</span>\n\n<span
  class=\"sd\">                ``` bash</span>\n<span class=\"sd\">                markata
  build --no-profile</span>\n<span class=\"sd\">                ```</span>\n<span
  class=\"sd\">                \"\"\"</span>\n\n                <span class=\"k\">if</span>
  <span class=\"n\">pretty</span><span class=\"p\">:</span>\n                    <span
  class=\"n\">make_pretty</span><span class=\"p\">()</span>\n\n                <span
  class=\"k\">if</span> <span class=\"n\">quiet</span><span class=\"p\">:</span>\n
  \                   <span class=\"n\">markata</span><span class=\"o\">.</span><span
  class=\"n\">console</span><span class=\"o\">.</span><span class=\"n\">quiet</span>
  <span class=\"o\">=</span> <span class=\"kc\">True</span>\n\n                <span
  class=\"k\">if</span> <span class=\"n\">verbose</span><span class=\"p\">:</span>\n
  \                   <span class=\"n\">markata</span><span class=\"o\">.</span><span
  class=\"n\">console</span><span class=\"o\">.</span><span class=\"n\">print</span><span
  class=\"p\">(</span><span class=\"s2\">\"console options:\"</span><span class=\"p\">,</span>
  <span class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">console</span><span
  class=\"o\">.</span><span class=\"n\">options</span><span class=\"p\">)</span>\n\n
  \               <span class=\"k\">if</span> <span class=\"n\">profile</span><span
  class=\"p\">:</span>\n                    <span class=\"n\">markata</span><span
  class=\"o\">.</span><span class=\"n\">should_profile_cli</span> <span class=\"o\">=</span>
  <span class=\"kc\">True</span>\n                    <span class=\"n\">markata</span><span
  class=\"o\">.</span><span class=\"n\">should_profile</span> <span class=\"o\">=</span>
  <span class=\"kc\">True</span>\n                    <span class=\"n\">markata</span><span
  class=\"o\">.</span><span class=\"n\">configure</span><span class=\"p\">()</span>\n\n
  \               <span class=\"k\">if</span> <span class=\"n\">should_pdb</span><span
  class=\"p\">:</span>\n                    <span class=\"n\">pdb_run</span><span
  class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">.</span><span
  class=\"n\">run</span><span class=\"p\">)</span>\n\n                <span class=\"k\">else</span><span
  class=\"p\">:</span>\n                    <span class=\"n\">markata</span><span
  class=\"o\">.</span><span class=\"n\">run</span><span class=\"p\">()</span>\n</pre></div>\n\n</pre>\n<p>!!
  function </p><h2 class=\"admonition-title\" id=\"list\" style=\"margin:0;padding:.5rem
  1rem;\">list <em class=\"small\">function</em></h2>\nProvides a way run markatas,
  map, filter, and sort from the\ncommand line.  I personally use this more often
  than the build\ncommand while I am writing on a site with a large number of\nposts
  on it.  It makes slicing in by <code>templatekey</code>, <code>tag</code>, or\n<code>date</code>
  much easier.\n<pre><code># default list\n\nBy default `markata list` will list all
  titles in a pager, for all posts\nbeing loaded by markata.\n\n``` bash\nmarkata
  list\n```\n\n# Skip the pager\n\nMarkata uses rich for its pager, it's pretty smart
  about when to\nuse the pager or pass text to the next thing in the pipeline,\nbut
  if you don't want to run a pager you can pass  `--no-pager`\n\n``` bash\nmarkata
  list --no-pager\n```\n\n# List other attributes\n\nYou can list any other attribute
  tied to your posts.  These are\nadded through either your yaml frontmatter at the
  start of your\npost, or through the use of a plugin.\n\n\n``` bash\n# the filepath
  of the post\nmarkata list --map path\n\n# the slug of the post (where it will show
  up on the site)\nmarkata list --map slug\n\n# the date of the post\nmarkata list
  --map date\n\n# the full raw content of the post\nmarkata list --map content\n```\n\n#
  List more than one attribute\n\nYou can create new attributes as you map to echo
  out by\ncombining existing attributes.\n\n``` bash\nmarkata list --map 'title +
  \" , \" + slug'\n```\n\n# Using Python objects as map\n\nYou can access attributes
  of each post attribute that you map\nover.  For instance on my blog, each post has
  a date that is a\ndatetime object.  I can ask each post for its `date.year`\n\n```
  bash\nmarkata list --map date.year\n\n# combining this with title\nmarkata list
  --map 'str(date.year) + \",\" + title'\n```\n\n# Filtering posts\n\nPosts are filtered
  with python syntax, you will have all\nattributes tied to your posts available to
  filter with.\n\n``` bash\nmarkata list --filter \"'__' not in title\"\n```\n\n#
  Filtering by dates\n\nIf your site has dates tied to your posts you can filter by\ndate.
  \ On my blog this makes a ton of sense and is quite useful.\nOn the Markata docs
  though it doesn't really make much sense,\nsince there really isn't the idea of
  a post date there.\n\n``` bash\n# listing today's posts\nmarkata list --filter \"date==today\"\n\n#
  listing this year's posts\nmarkata list --filter \"date.year==today.year\"\n```\n\n#
  Full Content Search\n\nYou can also search the full content of each post for specific\nwords.\n```
  bash\n\nmarkata list --filter \"'python' in content\"\n```\n\n# Filtering by frontmatter
  data\n\nI use a templateKey on my personal blog to determine which\ntemplate to
  render the page with.  I can fitler my posts by a\n`til` (today i learned) key.\n\n```
  bash\nmarkata list --filter \"templateKey=='til'\"\n```\n\n# Combining filters\n\nFilters
  can be combined together quite like maps can, it's all\njust python syntax.\n\n```
  bash\nmarkata list --filter \"templateKey=='til' and date == today\"\n```\n\n# Sorting
  posts\n\nPosts can be sorted by attributes on your post, and they can\neven be reversed.\n\n```
  bash\nmarkta list --sort date\nmarkta list --sort date --reverse\n```\n\n# Putting
  it all together\n\nThe real power of all this comes when you combine them all into\nlists
  that work for you and your workflow.  This really makes\nworking on larger projects
  so much easier to find things.\n\n\n# Making a fuzzy picker for your posts\n\nHere
  is a bash command to open an fzf picker for todays posts,\nthen open it in your
  `$EDITOR`\n\n``` bash\nmarkata list                 --map path                --filter
  'date==today'                --sort date                --reverse |                fzf
  --preview 'bat --color always {}' |                xargs -I {} $EDITOR {}\n```\n\n#
  Combining wtih nvim Telescope\n\nHere is the same command setup as a Telescope picker
  for neovim.\n\n``` vim\nnnoremap &lt;leader&gt;et &lt;cmd&gt;Telescope find_files
  find_command=markata,list,--map,path,--filter,date==today&lt;cr&gt;\n```\n\nIf you
  have another way to open posts in your editor with\n`markata list` I would love
  to accept a PR to add it to the\nexamples here.\n</code></pre>\n<div class=\"admonition
  source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"list
  <em class=\"small\">source</em>\"</p>\n</div>\n<pre class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button
  class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
  class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
  v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
  c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
  c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
  h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
  M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
  v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
  M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
  <span class=\"nf\">list</span><span class=\"p\">(</span>\n                <span
  class=\"nb\">map</span><span class=\"p\">:</span> <span class=\"nb\">str</span>
  <span class=\"o\">=</span> <span class=\"s2\">\"title\"</span><span class=\"p\">,</span>\n
  \               <span class=\"nb\">filter</span><span class=\"p\">:</span> <span
  class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">\"True\"</span><span
  class=\"p\">,</span>\n                <span class=\"n\">sort</span><span class=\"p\">:</span>
  <span class=\"nb\">str</span> <span class=\"o\">=</span> <span class=\"s2\">\"True\"</span><span
  class=\"p\">,</span>\n                <span class=\"n\">head</span><span class=\"p\">:</span>
  <span class=\"n\">Optional</span><span class=\"p\">[</span><span class=\"nb\">int</span><span
  class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"kc\">None</span><span
  class=\"p\">,</span>\n                <span class=\"n\">tail</span><span class=\"p\">:</span>
  <span class=\"n\">Optional</span><span class=\"p\">[</span><span class=\"nb\">int</span><span
  class=\"p\">]</span> <span class=\"o\">=</span> <span class=\"kc\">None</span><span
  class=\"p\">,</span>\n                <span class=\"n\">include_empty</span><span
  class=\"p\">:</span> <span class=\"nb\">bool</span> <span class=\"o\">=</span> <span
  class=\"kc\">False</span><span class=\"p\">,</span>\n                <span class=\"n\">reverse</span><span
  class=\"p\">:</span> <span class=\"nb\">bool</span> <span class=\"o\">=</span> <span
  class=\"kc\">False</span><span class=\"p\">,</span>\n                <span class=\"n\">use_pager</span><span
  class=\"p\">:</span> <span class=\"nb\">bool</span> <span class=\"o\">=</span> <span
  class=\"n\">typer</span><span class=\"o\">.</span><span class=\"n\">Option</span><span
  class=\"p\">(</span><span class=\"kc\">True</span><span class=\"p\">,</span> <span
  class=\"s2\">\"--pager\"</span><span class=\"p\">,</span> <span class=\"s2\">\"--no-pager\"</span><span
  class=\"p\">),</span>\n            <span class=\"p\">)</span> <span class=\"o\">-&gt;</span>
  <span class=\"kc\">None</span><span class=\"p\">:</span>\n<span class=\"w\">                </span><span
  class=\"sd\">\"\"\"</span>\n<span class=\"sd\">                Provides a way run
  markatas, map, filter, and sort from the</span>\n<span class=\"sd\">                command
  line.  I personally use this more often than the build</span>\n<span class=\"sd\">
  \               command while I am writing on a site with a large number of</span>\n<span
  class=\"sd\">                posts on it.  It makes slicing in by `templatekey`,
  `tag`, or</span>\n<span class=\"sd\">                `date` much easier.</span>\n\n<span
  class=\"sd\">                # default list</span>\n\n<span class=\"sd\">                By
  default `markata list` will list all titles in a pager, for all posts</span>\n<span
  class=\"sd\">                being loaded by markata.</span>\n\n<span class=\"sd\">
  \               ``` bash</span>\n<span class=\"sd\">                markata list</span>\n<span
  class=\"sd\">                ```</span>\n\n<span class=\"sd\">                #
  Skip the pager</span>\n\n<span class=\"sd\">                Markata uses rich for
  its pager, it's pretty smart about when to</span>\n<span class=\"sd\">                use
  the pager or pass text to the next thing in the pipeline,</span>\n<span class=\"sd\">
  \               but if you don't want to run a pager you can pass  `--no-pager`</span>\n\n<span
  class=\"sd\">                ``` bash</span>\n<span class=\"sd\">                markata
  list --no-pager</span>\n<span class=\"sd\">                ```</span>\n\n<span class=\"sd\">
  \               # List other attributes</span>\n\n<span class=\"sd\">                You
  can list any other attribute tied to your posts.  These are</span>\n<span class=\"sd\">
  \               added through either your yaml frontmatter at the start of your</span>\n<span
  class=\"sd\">                post, or through the use of a plugin.</span>\n\n\n<span
  class=\"sd\">                ``` bash</span>\n<span class=\"sd\">                #
  the filepath of the post</span>\n<span class=\"sd\">                markata list
  --map path</span>\n\n<span class=\"sd\">                # the slug of the post (where
  it will show up on the site)</span>\n<span class=\"sd\">                markata
  list --map slug</span>\n\n<span class=\"sd\">                # the date of the post</span>\n<span
  class=\"sd\">                markata list --map date</span>\n\n<span class=\"sd\">
  \               # the full raw content of the post</span>\n<span class=\"sd\">                markata
  list --map content</span>\n<span class=\"sd\">                ```</span>\n\n<span
  class=\"sd\">                # List more than one attribute</span>\n\n<span class=\"sd\">
  \               You can create new attributes as you map to echo out by</span>\n<span
  class=\"sd\">                combining existing attributes.</span>\n\n<span class=\"sd\">
  \               ``` bash</span>\n<span class=\"sd\">                markata list
  --map 'title + \" , \" + slug'</span>\n<span class=\"sd\">                ```</span>\n\n<span
  class=\"sd\">                # Using Python objects as map</span>\n\n<span class=\"sd\">
  \               You can access attributes of each post attribute that you map</span>\n<span
  class=\"sd\">                over.  For instance on my blog, each post has a date
  that is a</span>\n<span class=\"sd\">                datetime object.  I can ask
  each post for its `date.year`</span>\n\n<span class=\"sd\">                ``` bash</span>\n<span
  class=\"sd\">                markata list --map date.year</span>\n\n<span class=\"sd\">
  \               # combining this with title</span>\n<span class=\"sd\">                markata
  list --map 'str(date.year) + \",\" + title'</span>\n<span class=\"sd\">                ```</span>\n\n<span
  class=\"sd\">                # Filtering posts</span>\n\n<span class=\"sd\">                Posts
  are filtered with python syntax, you will have all</span>\n<span class=\"sd\">                attributes
  tied to your posts available to filter with.</span>\n\n<span class=\"sd\">                ```
  bash</span>\n<span class=\"sd\">                markata list --filter \"'__' not
  in title\"</span>\n<span class=\"sd\">                ```</span>\n\n<span class=\"sd\">
  \               # Filtering by dates</span>\n\n<span class=\"sd\">                If
  your site has dates tied to your posts you can filter by</span>\n<span class=\"sd\">
  \               date.  On my blog this makes a ton of sense and is quite useful.</span>\n<span
  class=\"sd\">                On the Markata docs though it doesn't really make much
  sense,</span>\n<span class=\"sd\">                since there really isn't the idea
  of a post date there.</span>\n\n<span class=\"sd\">                ``` bash</span>\n<span
  class=\"sd\">                # listing today's posts</span>\n<span class=\"sd\">
  \               markata list --filter \"date==today\"</span>\n\n<span class=\"sd\">
  \               # listing this year's posts</span>\n<span class=\"sd\">                markata
  list --filter \"date.year==today.year\"</span>\n<span class=\"sd\">                ```</span>\n\n<span
  class=\"sd\">                # Full Content Search</span>\n\n<span class=\"sd\">
  \               You can also search the full content of each post for specific</span>\n<span
  class=\"sd\">                words.</span>\n<span class=\"sd\">                ```
  bash</span>\n\n<span class=\"sd\">                markata list --filter \"'python'
  in content\"</span>\n<span class=\"sd\">                ```</span>\n\n<span class=\"sd\">
  \               # Filtering by frontmatter data</span>\n\n<span class=\"sd\">                I
  use a templateKey on my personal blog to determine which</span>\n<span class=\"sd\">
  \               template to render the page with.  I can fitler my posts by a</span>\n<span
  class=\"sd\">                `til` (today i learned) key.</span>\n\n<span class=\"sd\">
  \               ``` bash</span>\n<span class=\"sd\">                markata list
  --filter \"templateKey=='til'\"</span>\n<span class=\"sd\">                ```</span>\n\n<span
  class=\"sd\">                # Combining filters</span>\n\n<span class=\"sd\">                Filters
  can be combined together quite like maps can, it's all</span>\n<span class=\"sd\">
  \               just python syntax.</span>\n\n<span class=\"sd\">                ```
  bash</span>\n<span class=\"sd\">                markata list --filter \"templateKey=='til'
  and date == today\"</span>\n<span class=\"sd\">                ```</span>\n\n<span
  class=\"sd\">                # Sorting posts</span>\n\n<span class=\"sd\">                Posts
  can be sorted by attributes on your post, and they can</span>\n<span class=\"sd\">
  \               even be reversed.</span>\n\n<span class=\"sd\">                ```
  bash</span>\n<span class=\"sd\">                markta list --sort date</span>\n<span
  class=\"sd\">                markta list --sort date --reverse</span>\n<span class=\"sd\">
  \               ```</span>\n\n<span class=\"sd\">                # Putting it all
  together</span>\n\n<span class=\"sd\">                The real power of all this
  comes when you combine them all into</span>\n<span class=\"sd\">                lists
  that work for you and your workflow.  This really makes</span>\n<span class=\"sd\">
  \               working on larger projects so much easier to find things.</span>\n\n\n<span
  class=\"sd\">                # Making a fuzzy picker for your posts</span>\n\n<span
  class=\"sd\">                Here is a bash command to open an fzf picker for todays
  posts,</span>\n<span class=\"sd\">                then open it in your `$EDITOR`</span>\n\n<span
  class=\"sd\">                ``` bash</span>\n<span class=\"sd\">                markata
  list \\</span>\n<span class=\"sd\">                        --map path\\</span>\n<span
  class=\"sd\">                        --filter 'date==today'\\</span>\n<span class=\"sd\">
  \                       --sort date\\</span>\n<span class=\"sd\">                        --reverse
  |\\</span>\n<span class=\"sd\">                        fzf --preview 'bat --color
  always {}' |\\</span>\n<span class=\"sd\">                        xargs -I {} $EDITOR
  {}</span>\n<span class=\"sd\">                ```</span>\n\n<span class=\"sd\">
  \               # Combining wtih nvim Telescope</span>\n\n<span class=\"sd\">                Here
  is the same command setup as a Telescope picker for neovim.</span>\n\n<span class=\"sd\">
  \               ``` vim</span>\n<span class=\"sd\">                nnoremap &lt;leader&gt;et
  &lt;cmd&gt;Telescope find_files find_command=markata,list,--map,path,--filter,date==today&lt;cr&gt;</span>\n<span
  class=\"sd\">                ```</span>\n\n<span class=\"sd\">                If
  you have another way to open posts in your editor with</span>\n<span class=\"sd\">
  \               `markata list` I would love to accept a PR to add it to the</span>\n<span
  class=\"sd\">                examples here.</span>\n<span class=\"sd\">                \"\"\"</span>\n\n
  \               <span class=\"n\">markata</span><span class=\"o\">.</span><span
  class=\"n\">console</span><span class=\"o\">.</span><span class=\"n\">quiet</span>
  <span class=\"o\">=</span> <span class=\"kc\">True</span>\n\n                <span
  class=\"n\">tail</span> <span class=\"o\">=</span> <span class=\"o\">-</span><span
  class=\"n\">tail</span> <span class=\"k\">if</span> <span class=\"n\">tail</span>
  <span class=\"k\">else</span> <span class=\"n\">tail</span>\n                <span
  class=\"n\">filtered</span> <span class=\"o\">=</span> <span class=\"n\">markata</span><span
  class=\"o\">.</span><span class=\"n\">map</span><span class=\"p\">(</span><span
  class=\"nb\">map</span><span class=\"p\">,</span> <span class=\"nb\">filter</span><span
  class=\"p\">,</span> <span class=\"n\">sort</span><span class=\"p\">)</span>\n                <span
  class=\"k\">if</span> <span class=\"ow\">not</span> <span class=\"n\">include_empty</span><span
  class=\"p\">:</span>\n                    <span class=\"n\">filtered</span> <span
  class=\"o\">=</span> <span class=\"p\">[</span><span class=\"n\">a</span> <span
  class=\"k\">for</span> <span class=\"n\">a</span> <span class=\"ow\">in</span> <span
  class=\"n\">filtered</span> <span class=\"k\">if</span> <span class=\"n\">a</span>
  <span class=\"o\">!=</span> <span class=\"s2\">\"\"</span><span class=\"p\">]</span>\n
  \               <span class=\"n\">filtered</span> <span class=\"o\">=</span> <span
  class=\"n\">filtered</span><span class=\"p\">[</span><span class=\"n\">tail</span><span
  class=\"p\">:</span><span class=\"n\">head</span><span class=\"p\">]</span>\n                <span
  class=\"k\">if</span> <span class=\"n\">reverse</span><span class=\"p\">:</span>\n
  \                   <span class=\"n\">filtered</span> <span class=\"o\">=</span>
  <span class=\"nb\">reversed</span><span class=\"p\">(</span><span class=\"n\">filtered</span><span
  class=\"p\">)</span>\n\n                <span class=\"n\">markata</span><span class=\"o\">.</span><span
  class=\"n\">console</span><span class=\"o\">.</span><span class=\"n\">quiet</span>
  <span class=\"o\">=</span> <span class=\"kc\">False</span>\n                <span
  class=\"k\">if</span> <span class=\"n\">markata</span><span class=\"o\">.</span><span
  class=\"n\">console</span><span class=\"o\">.</span><span class=\"n\">is_terminal</span>
  <span class=\"ow\">and</span> <span class=\"n\">use_pager</span><span class=\"p\">:</span>\n
  \                   <span class=\"k\">with</span> <span class=\"n\">markata</span><span
  class=\"o\">.</span><span class=\"n\">console</span><span class=\"o\">.</span><span
  class=\"n\">pager</span><span class=\"p\">():</span>\n                        <span
  class=\"k\">for</span> <span class=\"n\">a</span> <span class=\"ow\">in</span> <span
  class=\"n\">filtered</span><span class=\"p\">:</span>\n                            <span
  class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">console</span><span
  class=\"o\">.</span><span class=\"n\">print</span><span class=\"p\">(</span><span
  class=\"n\">a</span><span class=\"p\">,</span> <span class=\"n\">style</span><span
  class=\"o\">=</span><span class=\"s2\">\"purple\"</span><span class=\"p\">)</span>\n
  \               <span class=\"k\">else</span><span class=\"p\">:</span>\n                    <span
  class=\"k\">for</span> <span class=\"n\">a</span> <span class=\"ow\">in</span> <span
  class=\"n\">filtered</span><span class=\"p\">:</span>\n                        <span
  class=\"n\">markata</span><span class=\"o\">.</span><span class=\"n\">console</span><span
  class=\"o\">.</span><span class=\"n\">print</span><span class=\"p\">(</span><span
  class=\"n\">a</span><span class=\"p\">)</span>\n</pre></div>\n\n</pre>\n<p>!! function
  </p><h2 class=\"admonition-title\" id=\"clean\" style=\"margin:0;padding:.5rem 1rem;\">clean
  <em class=\"small\">function</em></h2>\nCleans up output generated by markata including
  both the output_dir and\nthe .markata_cache.\n<pre><code># Dry Run\n\nYou can run
  with `--dry-run` to see what markata is about to do.\n\n``` bash\nmarkata clean
  --dry-run\n[09:42:37] [DRYRUN] removing outptut directory: markout base_cli.py:371\n
  \          [DRYRUN] removing cache directory: .markata.cache base_cli.py:377\n\n```\n\n#
  Running clean\n\nRunning markata clean will fully delete all of the directories
  created\nby markata.\n\n``` bash\nmarkata clean\n[09:53:04]  removing outptut directory:
  markout base_cli.py:394\n            removing cache directory: .markata.cache base_cli.py:405\n```\n\n#
  Running Quietly\n\nRunning with `--quiet` will remove all of the directories created
  by\nmarkata without announcing what it is doing.\n\n``` bash\nmarkata clean --quiet\n```\n</code></pre>\n<div
  class=\"admonition source is-collapsible collapsible-open\">\n<p class=\"admonition-title\">\"clean
  <em class=\"small\">source</em>\"</p>\n</div>\n<pre class=\"wrapper\">\n\n<div class=\"copy-wrapper\">\n\n<button
  class=\"copy\" onclick=\"navigator.clipboard.writeText(this.parentElement.parentElement.querySelector('pre').textContent)\"
  title=\"copy code to clipboard\"><svg id=\"Layer_1\" style=\"enable-background:new
  0 0 115.77 122.88\" version=\"1.1\" viewbox=\"0 0 115.77 122.88\" x=\"0px\" xml:space=\"preserve\"
  xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"
  y=\"0px\"><style type=\"text/css\">.st0{fill-rule:evenodd;clip-rule:evenodd;}</style><g><path
  class=\"st0\" d=\"M89.62,13.96v7.73h12.19h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02v0.02
  v73.27v0.01h-0.02c-0.01,3.84-1.57,7.33-4.1,9.86c-2.51,2.5-5.98,4.06-9.82,4.07v0.02h-0.02h-61.7H40.1v-0.02
  c-3.84-0.01-7.34-1.57-9.86-4.1c-2.5-2.51-4.06-5.98-4.07-9.82h-0.02v-0.02V92.51H13.96h-0.01v-0.02c-3.84-0.01-7.34-1.57-9.86-4.1
  c-2.5-2.51-4.06-5.98-4.07-9.82H0v-0.02V13.96v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07V0h0.02h61.7
  h0.01v0.02c3.85,0.01,7.34,1.57,9.86,4.1c2.5,2.51,4.06,5.98,4.07,9.82h0.02V13.96L89.62,13.96z
  M79.04,21.69v-7.73v-0.02h0.02 c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v64.59v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h12.19V35.65
  v-0.01h0.02c0.01-3.85,1.58-7.34,4.1-9.86c2.51-2.5,5.98-4.06,9.82-4.07v-0.02h0.02H79.04L79.04,21.69z
  M105.18,108.92V35.65v-0.02 h0.02c0-0.91-0.39-1.75-1.01-2.37c-0.61-0.61-1.46-1-2.37-1v0.02h-0.01h-61.7h-0.02v-0.02c-0.91,0-1.75,0.39-2.37,1.01
  c-0.61,0.61-1,1.46-1,2.37h0.02v0.01v73.27v0.02h-0.02c0,0.91,0.39,1.75,1.01,2.37c0.61,0.61,1.46,1,2.37,1v-0.02h0.01h61.7h0.02
  v0.02c0.91,0,1.75-0.39,2.37-1.01c0.61-0.61,1-1.46,1-2.37h-0.02V108.92L105.18,108.92z\"></path></g></svg></button>\n</div>\n
  \       \n<div class=\"highlight\"><pre><span></span>        <span class=\"k\">def</span>
  <span class=\"nf\">clean</span><span class=\"p\">(</span>\n                <span
  class=\"n\">quiet</span><span class=\"p\">:</span> <span class=\"nb\">bool</span>
  <span class=\"o\">=</span> <span class=\"n\">typer</span><span class=\"o\">.</span><span
  class=\"n\">Option</span><span class=\"p\">(</span>\n                    <span class=\"kc\">False</span><span
  class=\"p\">,</span>\n                    <span class=\"s2\">\"--quiet\"</span><span
  class=\"p\">,</span>\n                    <span class=\"s2\">\"-q\"</span><span
  class=\"p\">,</span>\n                <span class=\"p\">),</span>\n                <span
  class=\"n\">dry_run</span><span class=\"p\">:</span> <span class=\"nb\">bool</span>
  <span class=\"o\">=</span> <span class=\"n\">typer</span><span class=\"o\">.</span><span
  class=\"n\">Option</span><span class=\"p\">(</span>\n                    <span class=\"kc\">False</span><span
  class=\"p\">,</span>\n                    <span class=\"s2\">\"--dry-run\"</span><span
  class=\"p\">,</span>\n                <span class=\"p\">),</span>\n            <span
  class=\"p\">):</span>\n<span class=\"w\">                </span><span class=\"sd\">\"\"\"</span>\n<span
  class=\"sd\">                Cleans up output generated by markata including both
  the output_dir and</span>\n<span class=\"sd\">                the .markata_cache.</span>\n\n<span
  class=\"sd\">                # Dry Run</span>\n\n<span class=\"sd\">                You
  can run with `--dry-run` to see what markata is about to do.</span>\n\n<span class=\"sd\">
  \               ``` bash</span>\n<span class=\"sd\">                markata clean
  --dry-run</span>\n<span class=\"sd\">                [09:42:37] [DRYRUN] removing
  outptut directory: markout base_cli.py:371</span>\n<span class=\"sd\">                           [DRYRUN]
  removing cache directory: .markata.cache base_cli.py:377</span>\n\n<span class=\"sd\">
  \               ```</span>\n\n<span class=\"sd\">                # Running clean</span>\n\n<span
  class=\"sd\">                Running markata clean will fully delete all of the
  directories created</span>\n<span class=\"sd\">                by markata.</span>\n\n<span
  class=\"sd\">                ``` bash</span>\n<span class=\"sd\">                markata
  clean</span>\n<span class=\"sd\">                [09:53:04]  removing outptut directory:
  markout base_cli.py:394</span>\n<span class=\"sd\">                            removing
  cache directory: .markata.cache base_cli.py:405</span>\n<span class=\"sd\">                ```</span>\n\n<span
  class=\"sd\">                # Running Quietly</span>\n\n<span class=\"sd\">                Running
  with `--quiet` will remove all of the directories created by</span>\n<span class=\"sd\">
  \               markata without announcing what it is doing.</span>\n\n<span class=\"sd\">
  \               ``` bash</span>\n<span class=\"sd\">                markata clean
  --quiet</span>\n<span class=\"sd\">                ```</span>\n<span class=\"sd\">
  \               \"\"\"</span>\n                <span class=\"n\">_clean</span><span
  class=\"p\">(</span><span class=\"n\">markata</span><span class=\"o\">=</span><span
  class=\"n\">markata</span><span class=\"p\">,</span> <span class=\"n\">quiet</span><span
  class=\"o\">=</span><span class=\"n\">quiet</span><span class=\"p\">,</span> <span
  class=\"n\">dry_run</span><span class=\"o\">=</span><span class=\"n\">dry_run</span><span
  class=\"p\">)</span>\n</pre></div>\n\n</pre>\n</main>\n<footer>\xA9 2024</footer>\n</body></html>"
published: true
slug: markata/plugins/base-cli
title: Base_Cli.Py


---

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


!! function <h2 id='make_pretty' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>make_pretty <em class='small'>function</em></h2>
    This is a helper function that enables suppresses tracebacks from
    frameworks like `click` that can make your traceback long and hard
    to follow.  It also makes evrerything more colorful and easier to
    follow.
???+ source "make_pretty <em class='small'>source</em>"

```python

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
```


!! function <h2 id='cli' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>cli <em class='small'>function</em></h2>
    Markata hook to implement base cli commands.
???+ source "cli <em class='small'>source</em>"

```python

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
```


!! function <h2 id='_clean' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>_clean <em class='small'>function</em></h2>

???+ source "_clean <em class='small'>source</em>"

```python

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
```


!! function <h2 id='pdb_run' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>pdb_run <em class='small'>function</em></h2>
    Wraps a function call with a post_mortem pdb debugger.
???+ source "pdb_run <em class='small'>source</em>"

```python

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
```


!! function <h2 id='tui' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>tui <em class='small'>function</em></h2>

???+ source "tui <em class='small'>source</em>"

```python

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
```


!! function <h2 id='plugins' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>plugins <em class='small'>function</em></h2>
    create new things from templates
???+ source "plugins <em class='small'>source</em>"

```python

        def plugins():
                "create new things from templates"
```


!! function <h2 id='config' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>config <em class='small'>function</em></h2>
    configuration management
???+ source "config <em class='small'>source</em>"

```python

        def config():
                "configuration management"
```


!! function <h2 id='show' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>show <em class='small'>function</em></h2>

???+ source "show <em class='small'>source</em>"

```python

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
```


!! function <h2 id='generate' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>generate <em class='small'>function</em></h2>

???+ source "generate <em class='small'>source</em>"

```python

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
```


!! function <h2 id='get' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>get <em class='small'>function</em></h2>

???+ source "get <em class='small'>source</em>"

```python

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
```


!! function <h2 id='new' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>new <em class='small'>function</em></h2>
    create new things from templates
???+ source "new <em class='small'>source</em>"

```python

        def new():
                "create new things from templates"
```


!! function <h2 id='blog' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>blog <em class='small'>function</em></h2>
    Create a new blog from using the template from
    https://github.com/WaylonWalker/markata-blog-starter.
???+ source "blog <em class='small'>source</em>"

```python

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
```


!! function <h2 id='post' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>post <em class='small'>function</em></h2>
    Create new blog post in the pages directory from the template at
    https://github.com/WaylonWalker/markata-post-template.
???+ source "post <em class='small'>source</em>"

```python

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
```


!! function <h2 id='plugin' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>plugin <em class='small'>function</em></h2>
    Create a new plugin using the template at
    https://github.com/WaylonWalker/markata-plugin-template.
???+ source "plugin <em class='small'>source</em>"

```python

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
```


!! function <h2 id='build' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>build <em class='small'>function</em></h2>
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
???+ source "build <em class='small'>source</em>"

```python

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
```


!! function <h2 id='list' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>list <em class='small'>function</em></h2>
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
    markata list                 --map path                --filter 'date==today'                --sort date                --reverse |                fzf --preview 'bat --color always {}' |                xargs -I {} $EDITOR {}
    ```

    # Combining wtih nvim Telescope

    Here is the same command setup as a Telescope picker for neovim.

    ``` vim
    nnoremap <leader>et <cmd>Telescope find_files find_command=markata,list,--map,path,--filter,date==today<cr>
    ```

    If you have another way to open posts in your editor with
    `markata list` I would love to accept a PR to add it to the
    examples here.
???+ source "list <em class='small'>source</em>"

```python

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
```


!! function <h2 id='clean' class='admonition-title' style='margin:0;padding:.5rem 1rem;'>clean <em class='small'>function</em></h2>
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
???+ source "clean <em class='small'>source</em>"

```python

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
```
