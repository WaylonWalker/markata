"""Markata is a tool for handling directories of markdown.
"""
# annotations needed to return self
from __future__ import annotations

import hashlib
import importlib
import os
from pathlib import Path
import sys
from typing import Any, Callable, Dict, Iterable, List, Tuple

from checksumdir import dirhash
from diskcache import FanoutCache
import frontmatter
import markdown
import pluggy
from rich.console import Console
from rich.progress import track
from rich.table import Table

from datetime import timedelta
from markata import hookspec, standard_config
from markata.errors import MarkataConfigError


__version__ = "0.0.1"


DEFAULT_MD_EXTENSIONS = [
    "markdown.extensions.toc",
    "markdown.extensions.admonition",
    "markdown.extensions.tables",
    "markdown.extensions.md_in_html",
    "pymdownx.magiclink",
    "pymdownx.betterem",
    "pymdownx.tilde",
    "pymdownx.emoji",
    "pymdownx.tasklist",
    "pymdownx.superfences",
    "pymdownx.highlight",
    "pymdownx.inlinehilite",
    "pymdownx.keys",
    "pymdownx.saneheaders",
    "codehilite",
]

DEFAULT_HOOKS = [
    "markata.plugins.pyinstrument",
    "markata.plugins.glob",
    "markata.plugins.load",
    "markata.plugins.render_markdown",
    "markata.plugins.manifest",
    # "markata.plugins.generator",
    "markata.plugins.long_description",
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
]

DEFUALT_CONFIG = {
    "glob_patterns": ["**/*.md"],
    "hooks": ["default"],
    "markdown_extensions": [],
    "disabled_hooks": "",
    "default_cache_expire": 3600,
}


class Post(frontmatter.Post):
    html: str


def set_phase(function: Callable) -> Any:
    def wrapper(self: Markata, *args: Tuple, **kwargs: Dict) -> Any:
        self.phase = function.__name__
        result = function(self, *args, **kwargs)
        self.phase = function.__name__
        self.phase_file.write_text(self.phase)
        return result

    return wrapper


class Markata:
    def __init__(self) -> None:
        self.phase = "starting"
        self.MARKATA_CACHE_DIR = Path(".") / ".markata.cache"
        self.MARKATA_CACHE_DIR.mkdir(exist_ok=True)
        self.phase_file = self.MARKATA_CACHE_DIR / "phase.txt"
        self.registered_attrs = hookspec.registered_attrs
        self.configure()

    @property
    def cache(self):
        return FanoutCache(self.MARKATA_CACHE_DIR, statistics=True)

    def __getattr__(self, item):
        if item in self.__dict__.keys():
            return self.__getitem__(item)
        elif item in self.registered_attrs.keys():
            stage_to_run_to = max(
                [attr["lifecycle"] for attr in self.registered_attrs[item]]
            ).name
            runner = getattr(self, stage_to_run_to)
            runner()
            return getattr(self, item)
        else:
            raise AttributeError(item)

    @property
    def server(self):
        try:
            return self._server
        except AttributeError:
            from markata.cli.server import Server

            self._server = Server()
            return self.server

    @property
    def runner(self):
        try:
            return self._runner
        except AttributeError:
            from markata.cli.runner import Runner

            self._runner = Runner()
            return self.runner

    @property
    def plugins(self):
        try:
            return self._plugins
        except AttributeError:
            from markata.cli.plugins import Plugins

            self._plugins = Plugins(self)
        return self.plugins

    @property
    def summary(self):
        try:
            return self._summary
        except AttributeError:
            from markata.cli.summary import Summary

            self._summary = Summary(self)
            return self.summary

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

        if "seo" not in self.config:
            self.seo = [""]
        if isinstance(self.config["seo"], str):
            self.seo = self.config["seo"].split(",")
        if isinstance(self.config["seo"], list):
            self.seo = self.config["seo"]

        if "covers" not in self.config:
            self.covers = [""]
        if isinstance(self.config["covers"], str):
            self.covers = self.config["covers"].split(",")
        if isinstance(self.config["covers"], list):
            self.covers = self.config["covers"]
        else:
            raise TypeError("covers must be a string or list of dicts")

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

        if "markdown_extensions" not in self.config:
            markdown_extensions = [""]
        if isinstance(self.config["markdown_extensions"], str):
            markdown_extensions = [self.config["markdown_extensions"]]
        if isinstance(self.config["markdown_extensions"], list):
            markdown_extensions = self.config["markdown_extensions"]
        else:
            raise TypeError("markdown_extensions should be List[str]")

        self.markdown_extensions = [*DEFAULT_MD_EXTENSIONS, *markdown_extensions]
        self.md = markdown.Markdown(extensions=self.markdown_extensions)

        if "url" not in self.config:
            self.url = ""
        else:
            self.url = str(self.config["url"])

        if "output_dir" not in self.config:
            raise MarkataConfigError("output_dir must be specified in markata config")
        else:
            self.output_dir = Path(str(self.config["output_dir"]))

        if "assets_dir" not in self.config:
            raise MarkataConfigError("assets_dir must be specified in markata config")
        else:
            self.assets_dir = Path(str(self.config["assets_dir"]))

        if "site_name" not in self.config:
            self.site_name = ""
        else:
            self.site_name = str(self.config["site_name"])

        if "twitter_card" not in self.config:
            self.twitter_card = ""
        else:
            self.twitter_card = str(self.config["twitter_card"])

        if "site_name" not in self.config:
            self.site_name = ""
        else:
            self.site_name = str(self.config["site_name"])

        if "title" not in self.config:
            self.title = ""
        else:
            self.title = str(self.config["title"])

        if "author_name" not in self.config:
            self.author_name = ""
        else:
            self.author_name = str(self.config["author_name"])

        if "author_email" not in self.config:
            self.author_email = ""
        else:
            self.author_email = str(self.config["author_email"])

        if "icon" not in self.config:
            self.icon = ""
        else:
            self.icon = str(self.config["icon"])

        if "description" not in self.config:
            self.description = ""
        else:
            self.description = str(self.config["description"])

        if "rss_description" not in self.config:
            self.rss_description = ""
        else:
            self.rss_description = str(self.config["rss_description"])

        if "lang" not in self.config:
            self.lang = ""
        else:
            self.lang = str(self.config["lang"])

        if "post_template" not in self.config:
            self.post_template = ""
        else:
            self.post_template = str(self.config["post_template"])

        return self

    def get_plugin_config(self, plugin_path: str):

        key = Path(plugin_path).stem
        try:
            config = self.config[key]
        except KeyError:
            config = {}
        if "cache_expire" not in config.keys():
            config["cache_expire"] = self.config["default_cache_expire"]
        if "config_key" not in config.keys():
            config["config_key"] = key
        return config

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
    def hooks(self) -> List[str]:
        return self._hooks

    @hooks.setter
    def hooks(self, hooks: List[str]) -> None:
        self._hooks = hooks

    @property
    def disabled_hooks(self) -> List[str]:
        return self._disabled_hooks

    @disabled_hooks.setter
    def disabled_hooks(self, hooks: List[str]) -> None:
        self._disabled_hooks = hooks

    @property
    def files(self) -> List["Path"]:
        try:
            return self._files
        except AttributeError:
            self.glob()
            return self._files

    @files.setter
    def files(self, files: List["Path"]) -> None:
        self._files = files

    @property
    def content_directories(self) -> List["Path"]:
        try:
            return self._content_directories
        except AttributeError:
            self.glob()
            return self._content_directories

    @content_directories.setter
    def content_directories(self, files: List["Path"]) -> None:
        if self.phase == "glob":
            self._content_directories = files
        else:
            raise RuntimeWarning("cannot set content_directories outside of glob phase")

    @property
    def content_dir_hash(self) -> str:
        hashes = [dirhash(dir) for dir in self.content_directories]
        return self.make_hash(*hashes)

    @property
    def articles(self) -> List[frontmatter.Post]:
        try:
            return self._articles
        except AttributeError:
            self.load()
            return self._articles

    @articles.setter
    def articles(self, articles: List[frontmatter.Post]) -> None:
        self._articles = articles

    @property
    def console(self) -> Console:
        try:
            return self._console
        except AttributeError:
            self._console: Console = Console()
            return self._console

    def describe(self) -> dict[str, str]:
        return {"version": __version__, "phase": self.phase}

    def _to_dict(self) -> dict[str, Iterable]:
        return {"config": self.config, "articles": [a.to_dict() for a in self.articles]}

    def to_dict(self) -> dict:
        try:
            return self._to_dict()
        except AttributeError:
            self.render()
            return self._to_dict()

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

    @set_phase
    def glob(self) -> Markata:
        """run glob hooks

        Glob hooks should append file lists to the markata object for later
        hooks to build from.  The default loader will utilize the `files`
        attribute for loading.
        """

        try:
            self._pm.hook.glob(markata=self)
        except AttributeError:
            self.configure()
            self._pm.hook.glob(markata=self)

        return self

    @set_phase
    def load(self) -> Markata:
        try:
            self._pm.hook.load(markata=self)
        except AttributeError:
            self.glob()
            self._pm.hook.load(markata=self)
        return self

    @set_phase
    def render(self) -> Markata:
        try:
            self._pm.hook.pre_render(markata=self)
            self._pm.hook.render(markata=self)
            self._pm.hook.post_render(markata=self)
        except AttributeError:
            self.load()
            self._pm.hook.pre_render(markata=self)
            self._pm.hook.render(markata=self)
            self._pm.hook.post_render(markata=self)
        return self

    @set_phase
    def save(self) -> Markata:
        try:
            self._pm.hook.save(markata=self)
        except AttributeError:
            self.render()
            self._pm.hook.save(markata=self)
        return self

    def run(
        self,
    ) -> Markata:
        self.configure()
        self.console.log("configure complete")
        self.glob()
        self.console.log("glob complete")
        self.load()
        self.console.log("load complete")
        self.render()
        self.console.log("render complete")
        self.save()
        self.console.log("save complete")

        with self.cache as cache:
            hits, misses = cache.stats()

        if hits + misses > 0:
            self.console.log(f"cache hit rate {round(hits/ (hits + misses)*100, 2)}%")
        self.console.log(f"cache hits/misses {hits}/{misses}")

        return self

    def filter(self, filter: str):
        return [
            a
            for a in self.articles
            if eval(filter, {**a.to_dict(), "timedelta": timedelta}, {})
        ]

    def map(self, func: str = "title", filter: str = "True", sort: str = "True"):
        import copy

        articles = copy.copy(self.articles)
        articles.sort(key=lambda a: eval(sort, a.to_dict(), {}))
        return [
            eval(func, {**a.to_dict(), "timedelta": timedelta}, {})
            for a in articles
            if eval(filter, {**a.to_dict(), "timedelta": timedelta}, {})
        ]


def clif() -> None:
    import sys
    import time

    from rich import pretty, traceback

    if "--no-rich" not in sys.argv:
        pretty.install()
        traceback.install()

    m = Markata()

    if "--quiet" in sys.argv or "-q" in sys.argv:
        m.console.quiet = True
    else:
        m.console.print("console options:", m.console.options)

    if "--to-dict" in sys.argv:
        m.console.quiet = True
        data = m.to_dict()
        m.console.quiet = False
        m.console.print(data)
        return

    if "--draft" in sys.argv:
        print("\n".join([a["path"] for a in m.articles if a["status"] == "draft"]))

        return

    if "--today" in sys.argv:
        print("\n".join([a["path"] for a in m.articles if a["date"] == a["today"]]))
        return

    if "--scheduled" in sys.argv:
        print("\n".join([a["path"] for a in m.articles if a["date"] > a["today"]]))
        return

    if "--back-days" in sys.argv:
        print("\n".join([a["path"] for a in m.articles if a["date"] > a["today"]]))
        return

    if "--watch" in sys.argv:

        hash = m.content_dir_hash
        m.run()
        console = Console()
        with console.status("waiting for change", spinner="aesthetic", speed=0.2):
            while True:
                if m.content_dir_hash != hash:
                    hash = m.content_dir_hash
                    m.run()
                time.sleep(0.1)

    m.run()
