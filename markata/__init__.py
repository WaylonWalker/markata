"""Markata is a tool for handling directories of markdown.
"""
# annotations needed to return self
from __future__ import annotations

import atexit
import datetime
import hashlib
import importlib
import logging
import os
import sys
import textwrap
from datetime import timedelta
from pathlib import Path
from typing import Any, Callable, Dict, Iterable, List, Optional, Tuple

import frontmatter
import pluggy
from checksumdir import dirhash
from diskcache import FanoutCache
from rich.console import Console
from rich.progress import track
from rich.table import Table

from markata import hookspec, standard_config
from markata.__about__ import __version__
from markata.errors import MissingFrontMatter
from markata.lifecycle import LifeCycle

logger = logging.getLogger("markata")

DEFAULT_MD_EXTENSIONS = [
    "codehilite",
    "markdown.extensions.admonition",
    "markdown.extensions.md_in_html",
    "markdown.extensions.tables",
    "markdown.extensions.toc",
    "pymdownx.betterem",
    "pymdownx.details",
    "pymdownx.emoji",
    "pymdownx.extra",
    "pymdownx.highlight",
    "pymdownx.inlinehilite",
    "pymdownx.keys",
    "pymdownx.magiclink",
    "pymdownx.saneheaders",
    "pymdownx.superfences",
    "pymdownx.tabbed",
    "pymdownx.tabbed",
    "pymdownx.tasklist",
    "pymdownx.tilde",
]

DEFAULT_HOOKS = [
    "markata.plugins.heading_link",
    "markata.plugins.pyinstrument",
    "markata.plugins.glob",
    "markata.plugins.load",
    "markata.plugins.auto_title",
    "markata.plugins.render_markdown",
    "markata.plugins.manifest",
    "markata.plugins.feeds",
    # "markata.plugins.generator",
    "markata.plugins.jinja_md",
    "markata.plugins.auto_description",
    "markata.plugins.seo",
    "markata.plugins.post_template",
    "markata.plugins.covers",
    "markata.plugins.copy_assets",
    "markata.plugins.publish_html",
    "markata.plugins.flat_slug",
    "markata.plugins.datetime",
    "markata.plugins.rss",
    "markata.plugins.icon_resize",
    "markata.plugins.sitemap",
    "markata.plugins.to_json",
    "markata.plugins.base_cli",
    "markata.cli.server",
    "markata.cli.runner",
    "markata.cli.plugins",
    "markata.cli.summary",
    "markata.plugins.tui",
    "markata.plugins.setup_logging",
    "markata.plugins.redirects",
]

DEFUALT_CONFIG = {
    "glob_patterns": ["**/*.md"],
    "hooks": ["default"],
    "disabled_hooks": [""],
    "markdown_extensions": [],
    "default_cache_expire": 3600,
    "output_dir": "markout",
    "assets_dir": "static",
}


class Post(frontmatter.Post):
    html: str


def set_phase(function: Callable) -> Any:
    def wrapper(self: Markata, *args: Tuple, **kwargs: Dict) -> Any:
        self.phase = function.__name__
        self.phase_file.parent.mkdir(exist_ok=True)
        self.phase_file.write_text(self.phase)
        result = function(self, *args, **kwargs)
        self.phase = function.__name__
        self.phase_file.parent.mkdir(exist_ok=True)
        self.phase_file.write_text(self.phase)
        return result

    return wrapper


class Markata:
    def __init__(self, console: Console = None) -> None:
        self.phase = "starting"
        self.MARKATA_CACHE_DIR = Path(".") / ".markata.cache"
        self.MARKATA_CACHE_DIR.mkdir(exist_ok=True)
        self.phase_file: Path = self.MARKATA_CACHE_DIR / "phase.txt"
        self.registered_attrs = hookspec.registered_attrs
        self.configure()
        if console is not None:
            self._console = console
        atexit.register(self.teardown)

        with self.cache as cache:
            self.init_cache_stats = cache.stats()

    @property
    def cache(self) -> FanoutCache:
        return FanoutCache(self.MARKATA_CACHE_DIR, statistics=True)

    def __getattr__(self, item: str) -> Any:
        if item in self._pm.hook.__dict__.keys():
            # item is a hook, return a callable function
            return lambda: self.run(item)

        if item in self.__dict__.keys():
            # item is an attribute, return it
            return self.__getitem__(item)

        elif item in self.registered_attrs.keys():
            # item is created by a plugin, run it
            stage_to_run_to = max(
                [attr["lifecycle"] for attr in self.registered_attrs[item]]
            ).name
            self.run(stage_to_run_to)
            return getattr(self, item)
        else:
            # Markata does not know what this is, raise
            raise AttributeError(f"'Markata' object has no attribute '{item}'")

    def __rich__(self) -> Table:

        grid = Table.grid()
        grid.add_column("label")
        grid.add_column("value")

        for label, value in self.describe().items():
            grid.add_row(label, value)

        return grid

    def bust_cache(self) -> Markata:
        with self.cache as cache:
            cache.clear()
        return self

    @set_phase
    def configure(self) -> Markata:
        sys.path.append(os.getcwd())
        self.config = {**DEFUALT_CONFIG, **standard_config.load("markata")}
        if isinstance(self.config["glob_patterns"], str):
            self.config["glob_patterns"] = self.config["glob_patterns"].split(",")
        elif isinstance(self.config["glob_patterns"], list):
            self.config["glob_patterns"] = list(self.config["glob_patterns"])
        else:
            raise TypeError("glob_patterns must be list or str")
        self.glob_patterns = self.config["glob_patterns"]

        if "hooks" not in self.config:
            self.hooks = [""]
        if isinstance(self.config["hooks"], str):
            self.hooks = self.config["hooks"].split(",")
        if isinstance(self.config["hooks"], list):
            self.hooks = self.config["hooks"]

        if "disabled_hooks" not in self.config:
            self.disabled_hooks = [""]
        if isinstance(self.config["disabled_hooks"], str):
            self.disabled_hooks = self.config["disabled_hooks"].split(",")
        if isinstance(self.config["disabled_hooks"], list):
            self.disabled_hooks = self.config["disabled_hooks"]

        if not self.config.get("output_dir", "markout").endswith(
            self.config.get("path_prefix", "")
        ):
            self.config["output_dir"] = (
                self.config.get("output_dir", "markout")
                + "/"
                + self.config.get("path_prefix", "").rstrip("/")
            )
        if (
            len((output_split := self.config.get("output_dir", "markout").split("/")))
            > 1
        ):
            if "path_prefix" not in self.config.keys():
                self.config["path_prefix"] = "/".join(output_split[1:]) + "/"
        if not self.config.get("path_prefix", "").endswith("/"):
            self.config["path_prefix"] = self.config.get("path_prefix", "") + "/"

        self.config["output_dir"] = self.config["output_dir"].lstrip("/")
        self.config["path_prefix"] = self.config["path_prefix"].lstrip("/")

        try:
            default_index = self.hooks.index("default")
            hooks = [
                *self.hooks[:default_index],
                *DEFAULT_HOOKS,
                *self.hooks[default_index + 1 :],
            ]
            self.hooks = [hook for hook in hooks if hook not in self.disabled_hooks]
        except ValueError:
            # 'default' is not in hooks , do not replace with default_hooks
            pass

        self._pm = pluggy.PluginManager("markata")
        self._pm.add_hookspecs(hookspec.MarkataSpecs)
        self._register_hooks()

        self._pm.hook.configure(markata=self)
        return self

    def get_plugin_config(self, path_or_name: str) -> Dict:

        key = Path(path_or_name).stem

        config = self.config.get(key, {})

        if not isinstance(config, dict):
            raise TypeError("must use dict")
        if "cache_expire" not in config.keys():
            config["cache_expire"] = self.config["default_cache_expire"]
        if "config_key" not in config.keys():
            config["config_key"] = key
        return config

    def get_config(
        self,
        key: str,
        default: str = "",
        warn: bool = True,
        suggested: Optional[str] = None,
    ) -> Any:
        if key in self.config.keys():
            return self.config[key]
        else:

            if suggested is None:
                suggested = textwrap.dedent(
                    f"""
                    [markata]
                    {key} = '{default}'
                    """
                )
            if warn:
                logger.warning(
                    textwrap.dedent(
                        f"""
                        Warning {key} is not set in markata config, sitemap will
                        be missing root site_name
                        to resolve this open your markata.toml and add

                        {suggested}

                        """
                    ),
                )
        return default

    def make_hash(self, *keys: str) -> str:
        str_keys = [str(key) for key in keys]
        return hashlib.md5("".join(str_keys).encode("utf-8")).hexdigest()

    @property
    def phase(self) -> str:
        return self._phase

    @phase.setter
    def phase(self, value: str) -> None:
        self._phase = value

    @property
    def content_dir_hash(self) -> str:
        hashes = [
            dirhash(dir)
            for dir in self.content_directories
            if dir.absolute() != Path(".").absolute()
        ]
        return self.make_hash(*hashes)

    @property
    def console(self) -> Console:
        try:
            return self._console
        except AttributeError:
            self._console = Console()
            return self._console

    def describe(self) -> dict[str, str]:
        return {"version": __version__, "phase": self.phase}

    def _to_dict(self) -> dict[str, Iterable]:
        return {"config": self.config, "articles": [a.to_dict() for a in self.articles]}

    def to_dict(self) -> dict:
        return self._to_dict()

    def to_json(self) -> str:
        import json

        return json.dumps(self.to_dict(), indent=4, sort_keys=True, default=str)

    def _register_hooks(self) -> None:
        for hook in self.hooks:
            try:
                # module style plugins
                plugin = importlib.import_module(hook)
            except ModuleNotFoundError as e:
                # class style plugins
                if "." in hook:
                    mod = importlib.import_module(".".join(hook.split(".")[:-1]))
                    plugin = getattr(mod, hook.split(".")[-1])
                else:
                    raise e

            self._pm.register(plugin)

    def __iter__(self, description: str = "working...") -> Iterable[frontmatter.Post]:
        articles: Iterable[frontmatter.Post] = track(
            self.articles, description=description, transient=True, console=self.console
        )
        return articles

    def iter_articles(self, description: str) -> Iterable[frontmatter.Post]:
        articles: Iterable[frontmatter.Post] = track(
            self.articles, description=description, transient=True, console=self.console
        )
        return articles

    def teardown(self) -> Markata:
        """give special access to the teardown lifecycle method"""
        self._pm.hook.teardown(markata=self)
        return self

    def run(self, lifecycle: LifeCycle = None) -> Markata:
        if lifecycle is None:
            lifecycle = getattr(LifeCycle, max(LifeCycle._member_map_))

        if isinstance(lifecycle, str):
            lifecycle = LifeCycle[lifecycle]

        stages_to_run = [m for m in LifeCycle._member_map_ if LifeCycle[m] <= lifecycle]

        self.console.log(f"running {stages_to_run}")
        for stage in stages_to_run:
            self.console.log(f"{stage} running")
            getattr(self._pm.hook, stage)(markata=self)
            self.console.log(f"{stage} complete")

        with self.cache as cache:
            hits, misses = cache.stats()

        if hits + misses > 0:
            self.console.log(
                f"lifetime cache hit rate {round(hits/ (hits + misses)*100, 2)}%"
            )

        if misses > 0:
            self.console.log(f"lifetime cache hits/misses {hits}/{misses}")

        hits -= self.init_cache_stats[0]
        misses -= self.init_cache_stats[1]

        if hits + misses > 0:
            self.console.log(
                f"run cache hit rate {round(hits/ (hits + misses)*100, 2)}%"
            )

        if misses > 0:
            self.console.log(f"run cache hits/misses {hits}/{misses}")

        return self

    def filter(self, filter: str) -> List:
        def evalr(a: Post) -> Any:
            try:
                return eval(
                    filter,
                    {**a.to_dict(), "timedelta": timedelta, "post": a, "m": self},
                    {},
                )
            except AttributeError:
                return eval(
                    filter,
                    {**a.to_dict(), "timedelta": timedelta, "post": a, "m": self},
                    {},
                )

        return [a for a in self.articles if evalr(a)]

    def map(
        self,
        func: str = "title",
        filter: str = "True",
        sort: str = "True",
        reverse: bool = True,
        *args: Tuple,
        **kwargs: Dict,
    ) -> List:
        import copy

        def try_sort(a: Any) -> int:

            if "datetime" in sort.lower():
                return a.get(sort, datetime.datetime(1970, 1, 1))

            if "date" in sort.lower():
                return a.get(sort, datetime.date(1970, 1, 1))

            try:
                value = eval(sort, a.to_dict(), {})
            except NameError:
                return -1
            return value
            try:
                return int(value)
            except TypeError:
                try:
                    return int(value.timestamp())
                except Exception:
                    try:
                        return int(
                            datetime.datetime.combine(
                                value, datetime.datetime.min.time()
                            ).timestamp()
                        )
                    except Exception:
                        try:
                            return sum([ord(c) for c in str(value)])
                        except Exception:
                            return -1

        articles = copy.copy(self.articles)
        articles.sort(key=try_sort)
        if reverse:
            articles.reverse()

        try:
            posts = [
                eval(
                    func,
                    {**a.to_dict(), "timedelta": timedelta, "post": a, "m": self},
                    {},
                )
                for a in articles
                if eval(
                    filter,
                    {**a.to_dict(), "timedelta": timedelta, "post": a, "m": self},
                    {},
                )
            ]

        except NameError as e:
            variable = str(e).split("'")[1]

            missing_in_posts = self.map(
                "path", filter=f'"{variable}" not in post.keys()'
            )
            message = (
                f"variable: '{variable}' is missing in {len(missing_in_posts)} posts"
            )
            if len(missing_in_posts) > 10:
                message += f"\nfirst 10 paths to posts missing {variable} [{','.join(missing_in_posts)}..."
            else:
                message += f"\npaths to posts missing {variable} {missing_in_posts}"

            raise MissingFrontMatter(message)

        return posts


def load_ipython_extension(ipython):
    ipython.user_ns["m"] = Markata()
    ipython.user_ns["markata"] = ipython.user_ns["m"]
