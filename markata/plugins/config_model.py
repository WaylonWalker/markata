import datetime
from pathlib import Path
from typing import TYPE_CHECKING, Optional

import pydantic
from polyfactory.factories.pydantic_factory import ModelFactory
from pydantic import AnyUrl, ConfigDict, PositiveInt
from pydantic_extra_types.color import Color
from pydantic_settings import BaseSettings
from rich.jupyter import JupyterMixin
from rich.pretty import Pretty

from markata import standard_config
from markata.hookspec import hook_impl, register_attr

if TYPE_CHECKING:
    from markata import Markata


class Config(BaseSettings, JupyterMixin):
    hooks: list[str] = ["default"]
    disabled_hooks: list[str] = []
    markdown_extensions: list[str] = []
    default_cache_expire: PositiveInt = 3600
    output_dir: pydantic.DirectoryPath = Path("markout")
    assets_dir: Path = pydantic.Field(
        Path("static"),
        description="The directory to store static assets",
    )
    nav: dict[str, str] = {"home": "/"}
    site_version: int = 1
    markdown_backend: str = "markdown-it-py"
    url: Optional[AnyUrl] = None
    title: Optional[str] = "Markata Site"
    description: Optional[str] = None
    rss_description: Optional[str] = None
    author_name: Optional[str] = None
    author_email: Optional[str] = None
    lang: str = "en"
    repo_url: Optional[AnyUrl] = None
    repo_branch: str = "main"
    theme_color: Color = "#322D39"
    background_color: Color = "#B73CF6"
    start_url: str = "/"
    site_name: Optional[str] = None
    short_name: Optional[str] = None
    display: str = "minimal-ui"
    twitter_card: str = "summary_large_image"
    twitter_creator: Optional[str] = None
    twitter_site: Optional[str] = None
    path_prefix: Optional[str] = ""
    model_config = ConfigDict(env_prefix="markata_", extra="allow")
    today: datetime.date = pydantic.Field(default_factory=datetime.date.today)

    def __getitem__(self, item):
        "for backwards compatability"
        return getattr(self, item)

    def __setitem__(self, key, item):
        "for backwards compatability"
        return setattr(self, key, item)

    def get(self, item, default):
        "for backwards compatability"
        return getattr(self, item, default)

    def keys(self):
        "for backwards compatability"
        return self.__dict__.keys()

    def toml(self: "Config") -> str:
        import tomlkit

        doc = tomlkit.document()

        for key, value in self.dict().items():
            doc.add(key, value)
            doc.add(tomlkit.comment(key))
            if value:
                doc[key] = value
        return tomlkit.dumps(doc)

    @pydantic.validator("output_dir", pre=True, always=True)
    def validate_output_dir_exists(cls, value: Path) -> Path:
        if not isinstance(value, Path):
            value = Path(value)
        value.mkdir(parents=True, exist_ok=True)
        return value

    @property
    def __rich__(self) -> Pretty:
        return lambda: Pretty(self)


@hook_impl
@register_attr("post_models")
def config_model(markata: "Markata") -> None:
    markata.config_models.append(Config)


@hook_impl(tryfirst=True)
@register_attr("config")
def load_config(markata: "Markata") -> None:
    if "config" not in markata.__dict__.keys():
        config = standard_config.load("markata")
        if config == {}:
            markata.config = markata.Config()
        else:
            markata.config = markata.Config.parse_obj(config)


class ConfigFactory(ModelFactory):
    __model__ = Config
