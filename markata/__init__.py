"""Markata is a tool for handling directories of markdown.
"""
# annotations needed to return self
from __future__ import annotations

import hashlib
import importlib
import os
import sys
from pathlib import Path
from typing import List, Iterable, Any

import markdown
import frontmatter
import pluggy
from diskcache import Cache
from rich.progress import track
from rich.console import Console

from markata import hookspec, standard_config
from markata.errors import MarkataConfigError

from checksumdir import dirhash


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
]

DEFUALT_CONFIG = {
    "glob_patterns": ["**/*.md"],
    "hooks": ["default"],
    "markdown_extensions": [],
    "disabled_hooks": "",
}


class Markata:
    def __init__(self) -> None:
        self.configure()
        self.MARKATA_CACHE_DIR = Path(".") / ".markata.cache"
        self.MARKATA_CACHE_DIR.mkdir(exist_ok=True)
        self.cache = Cache(self.MARKATA_CACHE_DIR, statistics=True)

    # def __setattr__(self, name, value):
    #     self.

    def bust_cache(self) -> Markata:
        self.cache.clear()
        # self = Markata()
        return self

    def configure(self) -> Markata:
        sys.path.append(os.getcwd())
        self.config = {**DEFUALT_CONFIG, **standard_config.load("markata")}
        if isinstance(self.config["glob_patterns"], str):
            self.config["glob_patterns"] = self.config["glob_patterns"].split(",")

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
        if isinstance(self.seo, str):
            self.seo = self.config["seo"].split(",")
        if isinstance(self.config["seo"], list):
            self.seo = self.config["seo"]

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

        return self

    def make_hash(self, *keys: str) -> str:
        return hashlib.md5("".join(keys).encode("utf-8")).hexdigest()

    @property
    def phase(self) -> str:
        return self._phase

    @phase.setter
    def phase(self, phase: str) -> None:
        self._phase = phase

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
        # if self.phase == "glob":
        #     self._files = files
        # else:
        #     raise RuntimeWarning("cannot set files outside of glob phase")

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
        # if self.phase == "load":
        #     self._articles = articles
        # else:
        #     raise RuntimeWarning("cannot set articles outside of load phase")

    @property
    def console(self) -> Console:
        try:
            return self._console
        except AttributeError:
            self._console: Console = Console()
            return self._console

    # @property
    # def html(self) -> List[str]:
    #     return self._html

    # @html.setter
    # def html(self, html):
    #     if self.phase == "load":
    #         self._html = html
    #     else:
    #         raise RuntimeWarning("cannot set html outside of render phase")
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

    def glob(self) -> Markata:
        """run glob hooks

        Glob hooks should append file lists to the markata object for later
        hooks to build from.  The default loader will utilize the `files`
        attribute for loading.
        """

        self.phase = "glob"
        try:
            self._pm.hook.glob(markata=self)
        except AttributeError:
            self.configure()
            self._pm.hook.glob(markata=self)

        return self

    def load(self) -> Markata:
        self.phase = "load"
        try:
            self._pm.hook.load(markata=self)
        except AttributeError:
            print("missed glob")
            self.glob()
            self._pm.hook.load(markata=self)
        return self

    def render(self) -> Markata:
        self.phase = "render"
        try:
            self._pm.hook.render(markata=self)
        except AttributeError:
            print("missed load")
            self.load()
            self._pm.hook.render(markata=self)
        return self

    def save(self) -> Markata:
        self.phase = "save"
        try:
            self._pm.hook.save(markata=self)
        except AttributeError:
            print("missed render")
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

        if self.cache.hits + self.cache.misses > 0:
            self.console.log(
                f"cache hit rate {round(self.cache.hits/ (self.cache.hits + self.cache.misses)*100, 2)}%"
            )
        self.console.log(f"cache hits/misses {self.cache.hits}/{self.cache.misses}")

        return self


def cli() -> None:
    from rich import pretty, traceback

    import sys
    import time

    if "--no-rich" not in sys.argv:
        pretty.install()
        traceback.install()

    m = Markata()

    print("console options:", m.console.options)

    if "--quiet" in sys.argv or "-q" in sys.argv:
        m.console.quiet = True

    if "--to-dict" in sys.argv:
        m.console.quiet = True
        data = m.to_dict()
        m.console.quiet = False
        m.console.print(data)
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
        return

    m.run()
