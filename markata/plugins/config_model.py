from typing import TYPE_CHECKING

import pydantic

from markata.hookspec import hook_impl, register_attr

if TYPE_CHECKING:
    from markata import Markata


class Config(pydantic.BaseModel):
    glob_patterns: list
    hooks: list
    disabled_hooks: list
    markdown_extensions: list
    default_cache_expire: int
    output_dir: str
    assets_dir: str
    nav: dict
    site_version: int
    markdown_backend: str
    url: str
    title: str
    description: str
    rss_description: str
    author_name: str
    author_email: str
    icon: str
    lang: str
    repo_url: str
    repo_branch: str
    theme_color: str
    background_color: str
    start_url: str
    site_name: str
    short_name: str
    display: str
    twitter_card: str
    twitter_creator: str
    twitter_site: str
    markdown_it_py: dict
    feeds_config: dict
    feeds: dict
    jinja_md: dict
    head: dict
    tui: dict
    summary: dict
    post_model: dict
    path_prefix: str


@hook_impl
@register_attr("post_models")
def post_model(markata: "Markata") -> None:
    markata.config_models.append(Config)
