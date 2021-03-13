"""Markata is a tool for handling directories of markdown.
"""
# annotations needed to return self
from __future__ import annotations

import hashlib
import importlib
import os
import sys
from pathlib import Path
from typing import TYPE_CHECKING, List

import markdown
import pluggy
from diskcache import Cache
from rich.progress import track
from rich.console import Console

from markata import hookspec, standard_config

__version__ = "0.0.1"

if TYPE_CHECKING:
    from markdown import Markdown

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

    def configure(self) -> None:
        sys.path.append(os.getcwd())
        self.config = {**DEFUALT_CONFIG, **standard_config.load("markata")}
        if isinstance(self.config["glob_patterns"], str):
            self.config["glob_patterns"] = self.config["glob_patterns"].split(",")
        if isinstance(self.config["hooks"], str):
            self.config["hooks"] = self.config["hooks"].split(",")

        try:
            default_index = self.config["hooks"].index("default")
            hooks = [
                *self.config["hooks"][:default_index],
                *DEFAULT_HOOKS,
                *self.config["hooks"][default_index + 1 :],
            ]
            self.config["hooks"] = [
                hook for hook in hooks if hook not in self.config["disabled_hooks"]
            ]
        except ValueError:
            # 'default' is not in hooks , do not replace with default_hooks
            pass

        self._pm = pluggy.PluginManager("markata")
        self._pm.add_hookspecs(hookspec.MarkataSpecs)
        self._register_hooks()

        extensions = self.config["markdown_extensions"]
        self.config["md_extensions"] = [*DEFAULT_MD_EXTENSIONS, *extensions]
        self.md = markdown.Markdown(extensions=self.config["md_extensions"])

        return self

    def make_hash(self, *keys: str):
        return hashlib.md5("".join(keys).encode("utf-8")).hexdigest()

    @property
    def phase(self) -> str:
        return self._phase

    @phase.setter
    def phase(self, phase: str) -> None:
        self._phase = phase

    @property
    def files(self) -> List["Path"]:
        return self._files

    @files.setter
    def files(self, files: List["Path"]) -> None:
        if self.phase == "glob":
            self._files = files
        else:
            raise RuntimeWarning("cannot set files outside of glob phase")

    @property
    def articles(self) -> List[Markdown]:
        return self._articles

    @articles.setter
    def articles(self, articles: List[Markdown]) -> None:
        if self.phase == "load":
            self._articles = articles
        else:
            raise RuntimeWarning("cannot set articles outside of load phase")

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
    def _to_dict(self) -> dict:
        return {"config": self.config, "articles": [a.to_dict() for a in self.articles]}

    def to_dict(self) -> list:
        try:
            return self._to_dict()
        except AttributeError:
            self.configure()
            self.glob()
            self.load()
            self.render()
            return self._to_dict()

    def _register_hooks(self) -> None:
        for hook in self.config["hooks"]:
            try:
                # module style plugins
                # print(f"importing hook as module style: {hook}")
                plugin = importlib.import_module(hook)
            except ModuleNotFoundError as e:
                # class style plugins
                # print(f"importing hook as class style: {hook}")
                if "." in hook:
                    mod = importlib.import_module(".".join(hook.split(".")[:-1]))
                    plugin = getattr(mod, hook.split(".")[-1])
                else:
                    raise e

            self._pm.register(plugin)

        # self._pm.register

    def __iter__(self, description="working..."):
        return track(
            self.articles, description=description, transient=True, console=self.console
        )

    def iter_articles(self, description: str) -> track:
        return track(
            self.articles, description=description, transient=True, console=self.console
        )

    def glob(self) -> Markata:
        """run glob hooks

        Glob hooks should append file lists to the markata object for later
        hooks to build from.  The default loader will utilize the `files`
        attribute for loading.
        """

        self.phase = "glob"
        self._pm.hook.glob(markata=self)
        return self

    def load(self) -> Markata:
        self.phase = "load"
        self._pm.hook.load(markata=self)
        return self

    def render(self) -> Markata:
        self.phase = "render"
        self._pm.hook.render(markata=self)
        return self

    def save(self) -> Markata:
        self.phase = "save"
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

    m.run()
