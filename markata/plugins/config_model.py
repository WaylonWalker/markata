from typing import TYPE_CHECKING, Optional

import pydantic
from polyfactory.factories.pydantic_factory import ModelFactory

from markata import standard_config
from markata.hookspec import hook_impl, register_attr
from pathlib import Path


if TYPE_CHECKING:
    from markata import Markata


class Config(pydantic.BaseModel):
    hooks: list[str] = ["default"]
    disabled_hooks: list[str] = []
    markdown_extensions: list[str] = []
    default_cache_expire: int = 3600
    output_dir: pydantic.DirectoryPath = "markout"
    assets_dir: Path = pydantic.Field(
        Path("static"),
        description="The directory to store static assets",
    )
    nav: dict = {"home": "/"}
    site_version: int = 1
    markdown_backend: str = "markdown-it-py"
    url: Optional[str] = None
    title: Optional[str] = None
    description: Optional[str] = None
    rss_description: Optional[str] = None
    author_name: Optional[str] = None
    author_email: Optional[str] = None
    icon: str = "icon.png"
    lang: str = "en"
    repo_url: str = None
    repo_branch: str = "main"
    theme_color: str = "#322D39"
    background_color: str = "#B73CF6"
    start_url: str = "/"
    site_name: Optional[str] = None
    short_name: Optional[str] = None
    display: str = "minimal-ui"
    twitter_card: str = "summary_large_image"
    twitter_creator: Optional[str] = None
    twitter_site: Optional[str] = None
    path_prefix: Optional[str] = None

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
            add_doc(doc, key, value)
            doc.add(tomlkit.comment(key))
            if value:
                doc[key] = value
        return tomlkit.dumps(doc)


# def add_doc(doc: pydantic.Document) -> None:


@hook_impl
@register_attr("post_models")
def config_model(markata: "Markata") -> None:
    markata.config_models.append(Config)


@hook_impl(tryfirst=True)
@register_attr("config")
def load_config(markata: "Markata") -> None:
    if markata.config is None:
        markata.config = markata.Config.parse_obj(standard_config.load("markata"))


class ConfigFactory(ModelFactory):
    __model__ = Config
