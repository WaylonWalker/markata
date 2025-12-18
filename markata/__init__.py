"""Markata is a tool for handling directories of markdown."""

# annotations needed to return self
from __future__ import annotations

import atexit
import datetime
import importlib
import logging
import os
import sys
import textwrap
from datetime import timedelta
from functools import lru_cache
from pathlib import Path
from typing import Any
from typing import Dict
from typing import Iterable
from typing import Optional

import pluggy
import pydantic
from diskcache import Cache
from rich.console import Console
from rich.progress import track
from rich.table import Table

from markata import hookspec
from markata import standard_config
from markata.__about__ import __version__
from markata.exceptions import NoPosts
from markata.exceptions import TooManyPosts
from markata.lifecycle import LifeCycle

logger = logging.getLogger("markata")


DEFAULT_MD_EXTENSIONS = [
    "codehilite",
    "markdown.extensions.admonition",
    "markdown.extensions.md_in_html",
    "markdown.extensions.tables",
    "markdown.extensions.toc",
    "markdown.extensions.wikilinks",
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
    "markata.plugins.mermaid",
    "markata.plugins.didyoumean",
    "markata.plugins.skip",
    "markata.plugins.md_it_wikilinks",
    "markata.plugins.copy_assets",
    # "markata.plugins.heading_link",
    "markata.plugins.pyinstrument",
    "markata.plugins.glob",
    "markata.plugins.load",
    "markata.plugins.auto_title",
    "markata.plugins.jinja_env",  # Add centralized jinja environment
    "markata.plugins.render_markdown",
    # "markata.plugins.manifest",
    # "markata.plugins.generator",
    "markata.plugins.feeds",
    "markata.plugins.auto_description",
    # "markata.plugins.seo",
    "markata.plugins.post_template",
    "markata.plugins.covers",
    "markata.plugins.publish_html",
    "markata.plugins.flat_slug",
    # "markata.plugins.rss",
    "markata.plugins.icon_resize",
    # "markata.plugins.sitemap",
    "markata.plugins.to_json",
    "markata.plugins.base_cli",
    "markata.cli.server",
    "markata.cli.runner",
    "markata.cli.plugins",
    "markata.cli.summary",
    "markata.plugins.tui",
    "markata.plugins.setup_logging",
    "markata.plugins.redirects",
    "markata.plugins.post_model",
    "markata.plugins.config_model",
    "markata.plugins.create_models",
    # "markata.plugins.jinja_md",
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


class HooksConfig(pydantic.BaseModel):
    model_config = pydantic.ConfigDict(
        validate_assignment=True,
        arbitrary_types_allowed=True,
        extra="allow",
        str_strip_whitespace=True,
        validate_default=True,
        coerce_numbers_to_str=True,
        populate_by_name=True,
    )
    hooks: list = ["default"]
    disabled_hooks: list = []


class Markata:
    def __init__(
        self: "Markata",
        console: Console = None,
        config=None,
        config_overrides: Optional[Dict[str, Any]] = None,
        config_file: Optional[Path] = None,
    ) -> None:
        self.__version__ = __version__
        self.stages_ran = set()
        self.threded = False
        self._cache = None
        self._precache = None
        self._map_cache_stats = {"hits": 0, "misses": 0, "total": 0}
        self.MARKATA_CACHE_DIR = Path(".") / ".markata.cache"
        self.MARKATA_CACHE_DIR.mkdir(exist_ok=True)
        self._pm = pluggy.PluginManager("markata")
        self._pm.add_hookspecs(hookspec.MarkataSpecs)

        # Store config overrides for later use in load_config hook
        self._config_overrides = config_overrides or {}
        self._config_file = config_file

        if config is not None:
            self.config = config
        with self.cache as cache:
            self.init_cache_stats = cache.stats()
        self.registered_attrs = hookspec.registered_attrs
        self.post_models = []
        self.config_models = []
        if config is not None:
            raw_hooks = config
        else:
            raw_hooks = standard_config.load(
                "markata",
                project_home=config_file.parent if config_file else ".",
                overrides=config_overrides or {},
                config_file=config_file,
            )
        self.hooks_conf = HooksConfig.parse_obj(raw_hooks)
        try:
            default_index = self.hooks_conf.hooks.index("default")
            hooks = [
                *self.hooks_conf.hooks[:default_index],
                *DEFAULT_HOOKS,
                *self.hooks_conf.hooks[default_index + 1 :],
            ]
            self.hooks_conf.hooks = [
                hook for hook in hooks if hook not in self.hooks_conf.disabled_hooks
            ]
        except ValueError:
            # 'default' is not in hooks , do not replace with default_hooks
            pass

        self._register_hooks()
        if console is not None:
            self._console = console
        atexit.register(self.teardown)
        self.precache

    @property
    def cache(self: "Markata") -> Cache:
        # if self.threded:
        #     FanoutCache(self.MARKATA_CACHE_DIR, statistics=True)
        if self._cache is not None:
            return self._cache
        self._cache = Cache(
            self.MARKATA_CACHE_DIR,
            statistics=True,
            size_limit=5 * 1024**3,  # 5GB to reduce culling frequency
            cull_limit=10,  # Evict fewer entries at a time (default is 100)
        )
        self._cache.expire()

        return self._cache

    @property
    def precache(self: "Markata") -> None:
        return self.cache
        if self._precache is None:
            self.cache.expire()
            self._precache = {k: self.cache.get(k) for k in self.cache.iterkeys()}
        return self._precache

    def __getattr__(self: "Markata", item: str) -> Any:
        if item in self._pm.hook.__dict__:
            # item is a hook, return a callable function
            return lambda: self.run(item)

        if item in self.__dict__:
            # item is an attribute, return it
            return self.__getitem__(item)

        elif item in self.registered_attrs:
            # item is created by a plugin, run it
            stage_to_run_to = max(
                [attr["lifecycle"] for attr in self.registered_attrs[item]],
            ).name
            self.console.log(
                f"Running to [purple]{stage_to_run_to}[/] to retrieve [purple]{item}[/]"
            )
            self.run(stage_to_run_to)
            # Check __dict__ directly to avoid infinite recursion
            if item in self.__dict__:
                return self.__dict__[item]
            else:
                raise AttributeError(
                    f"'Markata' object has no attribute '{item}' after running {stage_to_run_to}"
                )
        elif item == "precache":
            return self._precache or {}
        else:
            # Markata does not know what this is, raise
            raise AttributeError(f"'Markata' object has no attribute '{item}'")

    def __rich__(self: "Markata") -> Table:
        grid = Table.grid()
        grid.add_column("label")
        grid.add_column("value")

        for label, value in self.describe().items():
            grid.add_row(label, value)

        return grid

    def bust_cache(self: "Markata") -> Markata:
        with self.cache as cache:
            cache.clear()
        return self

    def configure(self) -> Markata:
        sys.path.append(os.getcwd())
        # self.config = {**DEFUALT_CONFIG, **standard_config.load("markata")}
        # if isinstance(self.config["glob_patterns"], str):
        #     self.config["glob_patterns"] = self.config["glob_patterns"].split(",")
        # elif isinstance(self.config["glob_patterns"], list):
        #     self.config["glob_patterns"] = list(self.config["glob_patterns"])
        # else:
        #     raise TypeError("glob_patterns must be list or str")
        # self.glob_patterns = self.config["glob_patterns"]

        # self.hooks = self.config["hooks"]

        # if "disabled_hooks" not in self.config:
        #     self.disabled_hooks = [""]
        # if isinstance(self.config["disabled_hooks"], str):
        #     self.disabled_hooks = self.config["disabled_hooks"].split(",")
        # if isinstance(self.config["disabled_hooks"], list):
        #     self.disabled_hooks = self.config["disabled_hooks"]

        # if not self.config.get("output_dir", "markout").endswith(
        #     self.config.get("path_prefix", "")
        # ):
        #     self.config["output_dir"] = (
        #         self.config.get("output_dir", "markout") +
        #         "/" +
        #         self.config.get("path_prefix", "").rstrip("/")
        #     )
        # if (
        #     len((output_split := self.config.get("output_dir", "markout").split("/"))) >
        #     1
        # ):
        #     if "path_prefix" not in self.config.keys():
        #         self.config["path_prefix"] = "/".join(output_split[1:]) + "/"
        # if not self.config.get("path_prefix", "").endswith("/"):
        #     self.config["path_prefix"] = self.config["path_prefix"] + "/"

        # self.config["output_dir"] = self.config["output_dir"].lstrip("/")
        # self.config["path_prefix"] = self.config["path_prefix"].lstrip("/")

        try:
            default_index = self.hooks_conf.hooks.index("default")
            hooks = [
                *self.hooks_conf.hooks[:default_index],
                *DEFAULT_HOOKS,
                *self.hooks_conf.hooks[default_index + 1 :],
            ]
            self.config.hooks = [
                hook for hook in hooks if hook not in self.config.disabled_hooks
            ]
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
        import xxhash

        str_keys = [str(key) for key in keys]
        hash = xxhash.xxh64("".join(str_keys).encode("utf-8")).hexdigest()
        return hash

    @property
    def content_dir_hash(self: "Markata") -> str:
        from checksumdir import dirhash

        hashes = [
            dirhash(dir)
            for dir in self.content_directories
            if dir.absolute() != Path(".").absolute()
        ]
        return self.make_hash(*hashes)

    @property
    def console(self: "Markata") -> Console:
        try:
            return self._console
        except AttributeError:
            self._console = Console()
            return self._console

    def describe(self: "Markata") -> dict[str, str]:
        return {"version": __version__}

    def _to_dict(self: "Markata") -> dict[str, Iterable]:
        return {"config": self.config, "articles": [a.to_dict() for a in self.articles]}

    def to_dict(self: "Markata") -> dict:
        return self._to_dict()

    def to_json(self: "Markata") -> str:
        import json

        return json.dumps(self.to_dict(), indent=4, sort_keys=True, default=str)

    def _register_hooks(self: "Markata") -> None:
        sys.path.append(os.getcwd())
        for hook in self.hooks_conf.hooks:
            try:
                # module style plugins
                plugin = importlib.import_module(hook)
            except ModuleNotFoundError as e:
                # class style plugins
                if "." in hook:
                    try:
                        mod = importlib.import_module(".".join(hook.split(".")[:-1]))
                        plugin = getattr(mod, hook.split(".")[-1])
                    except ModuleNotFoundError as e:
                        raise ModuleNotFoundError(
                            f"module {hook} not found\n{sys.path}"
                        ) from e
                else:
                    raise e

            self._pm.register(plugin)

    def __iter__(
        self: "Markata", description: str = "working..."
    ) -> Iterable["Markata.Post"]:
        articles: Iterable[Markata.Post] = track(
            self.articles,
            description=description,
            transient=False,
            console=self.console,
        )
        return articles

    def iter_articles(self: "Markata", description: str) -> Iterable[Markata.Post]:
        articles: Iterable[Markata.Post] = track(
            self.filter("skip == False"),
            description=description,
            transient=True,
            console=self.console,
        )
        return articles

    def teardown(self: "Markata"):
        """Cleanup and print statistics when Markata is done."""
        # Print map cache statistics if they exist
        if hasattr(self, "_map_cache_stats"):
            stats = self._map_cache_stats
            total = stats["total"]
            if total > 0:
                hit_rate = (stats["hits"] / total) * 100
                self.console.print("\n[yellow]Map Cache Statistics:[/yellow]")
                self.console.print(f"Total calls: {total}")
                self.console.print(f"Cache hits: {stats['hits']}")
                self.console.print(f"Cache misses: {stats['misses']}")
                self.console.print(f"Hit rate: {hit_rate:.1f}%")
                self.console.print(
                    f"Cache size: {len(getattr(self, '_filtered_cache', {}))}"
                )
        if self.stages_ran:
            self._pm.hook.teardown(markata=self)
        return self

    def run(self: "Markata", lifecycle: LifeCycle = None) -> Markata:
        if lifecycle is None:
            lifecycle = max(LifeCycle._member_map_.values())

        if isinstance(lifecycle, str):
            lifecycle = LifeCycle[lifecycle]

        stages_to_run = [
            m
            for m in LifeCycle._member_map_
            if (LifeCycle[m] <= lifecycle) and (m not in self.stages_ran)
        ]

        if not stages_to_run:
            self.console.log(f"{lifecycle.name} already ran")
            return self

        self.console.log(f"running {stages_to_run}")
        for stage in stages_to_run:
            self.console.log(f"{stage} running")
            getattr(self._pm.hook, stage)(markata=self)
            self.stages_ran.add(stage)
            self.console.log(f"{stage} complete")

        with self.cache as cache:
            hits, misses = cache.stats()

        if hits + misses > 0:
            self.console.log(
                f"lifetime cache hit rate {round(hits / (hits + misses) * 100, 2)}%",
            )

        if misses > 0:
            self.console.log(f"lifetime cache hits/misses {hits}/{misses}")

        hits -= self.init_cache_stats[0]
        misses -= self.init_cache_stats[1]

        if hits + misses > 0:
            self.console.log(
                f"run cache hit rate {round(hits / (hits + misses) * 100, 2)}%",
            )

        if misses > 0:
            self.console.log(f"run cache hits/misses {hits}/{misses}")

        return self

    def filter(self: "Markata", filter: str) -> list:
        def evalr(a: Markata.Post) -> Any:
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

    def _compile_sort_key(self, sort: str):
        """Compile a sort key function for better performance"""
        if "datetime" in sort.lower():
            return lambda a: a.get(sort, datetime.datetime(1970, 1, 1))
        if "date" in sort.lower():
            return lambda a: a.get(sort, datetime.date(1970, 1, 1))

        # Create a compiled function for complex sort expressions
        try:
            code = compile(sort, "<string>", "eval")

            def sort_key(a):
                try:
                    value = eval(code, a.to_dict(), {})
                    if isinstance(value, (int, float)):
                        return value
                    if hasattr(value, "timestamp"):
                        return value.timestamp()
                    if isinstance(value, datetime.date):
                        return datetime.datetime.combine(
                            value,
                            datetime.datetime.min.time(),
                        ).timestamp()
                    return sum(ord(c) for c in str(value))
                except Exception:
                    return -1

            return sort_key
        except Exception:
            return lambda _: -1

    @lru_cache(maxsize=32)
    def _get_sort_key(self, sort: str):
        """Cache compiled sort key functions"""
        return self._compile_sort_key(sort)

    def _get_eval_globals(self):
        """Get common globals used in eval operations"""
        if not hasattr(self, "_eval_globals"):
            self._eval_globals = {"timedelta": timedelta}
        return self._eval_globals

    def _eval_with_article(self, code, article, extra_globals=None):
        """Evaluate code with article context, reusing dict where possible"""
        if not hasattr(article, "_eval_dict"):
            article._eval_dict = article.to_dict()
            article._eval_dict.update({"post": article, "m": self})

        globals_dict = self._get_eval_globals()
        if extra_globals:
            globals_dict.update(extra_globals)

        try:
            return eval(code, article._eval_dict, globals_dict)
        except Exception:
            return None

    def map(
        self: "Markata",
        func: str = "title",
        filter: str = "True",
        sort: str = "True",
        reverse: bool = True,
        *args: tuple,
        **kwargs: dict,
    ) -> list:
        # Cache the filtered articles
        if not hasattr(self, "_filtered_cache"):
            self._filtered_cache = {}

        filter_key = (func, filter, sort, reverse, args, frozenset(kwargs.items()))

        self._map_cache_stats["total"] += 1

        if filter_key in self._filtered_cache:
            self._map_cache_stats["hits"] += 1
            articles = self._filtered_cache[filter_key]
        else:
            self._map_cache_stats["misses"] += 1
            filter_code = compile(filter, "<string>", "eval")
            eval_globals = (
                {"timedelta": timedelta, **kwargs}
                if kwargs
                else {"timedelta": timedelta}
            )

            # Filter in one pass
            articles = []
            for article in self.articles:
                try:
                    ctx = article.to_dict()
                    ctx["post"] = article
                    ctx["m"] = self
                    if eval(filter_code, ctx, eval_globals):
                        articles.append(article)
                except Exception:
                    continue

            # Sort if needed
            if sort != "True":
                sort_key = self._get_sort_key(sort)
                articles.sort(key=sort_key, reverse=reverse)
            elif reverse:
                articles.reverse()

            self._filtered_cache[filter_key] = articles

        # Map in one pass with same context structure
        map_code = compile(func, "<string>", "eval")
        eval_globals = (
            {"timedelta": timedelta, **kwargs} if kwargs else {"timedelta": timedelta}
        )
        return [
            eval(map_code, {"post": a, "m": self, **a.to_dict()}, eval_globals)
            for a in articles
        ]

    def first(
        self: "Markata",
        filter: str = "True",
        sort: str = "True",
        reverse: bool = True,
        *args: tuple,
        **kwargs: dict,
    ) -> list:
        return self.map("post", filter, sort, reverse, *args, **kwargs)[0]

    def last(
        self: "Markata",
        filter: str = "True",
        sort: str = "True",
        reverse: bool = True,
        *args: tuple,
        **kwargs: dict,
    ) -> list:
        return self.map("post", filter, sort, reverse, *args, **kwargs)[-1]

    def one(
        self: "Markata",
        filter: str = "True",
        *args: tuple,
        **kwargs: dict,
    ) -> list:
        posts = self.map("post", filter, *args, **kwargs)
        if len(posts) > 1:
            raise TooManyPosts(f"found {len(posts)} posts, expected 1. {posts}")
        if len(posts) == 0:
            raise NoPosts
        return posts[0]
