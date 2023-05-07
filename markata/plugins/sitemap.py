from pathlib import Path

import anyconfig
import pydantic

from markata import Markata
from markata.hookspec import hook_impl, register_attr


class SiteMapUrl(pydantic.BaseModel):
    markata: "Markata" = pydantic.Field(..., exclude=True)
    slug: str = pydantic.Field(..., exclude=True)
    loc: str = pydantic.Field(None, include=True)
    changefreq: str = pydantic.Field("daily", include=True)
    priority: str = pydantic.Field("0.7", include=True)

    class Config:
        arbitrary_types_allowed = True

    @pydantic.validator("loc", always=True)
    def default_loc(cls, v, *, values):
        if v is None:
            return values["markata"].config.url + "/" + values["slug"] + "/"
        return v

    def dict(self, *args, **kwargs):
        return {"url": {**super().dict(*args, **kwargs)}}


class SiteMapPost(pydantic.BaseModel):
    markata: "Markata"
    slug: str
    published: bool = True
    sitemap_url: SiteMapUrl = None

    class Config:
        arbitrary_types_allowed = True

    @pydantic.validator("sitemap_url", always=True)
    def default_loc(cls, v, *, values):
        if v is None:
            return SiteMapUrl(markata=values["markata"], slug=values["slug"])
        return v


@hook_impl()
@register_attr("post_models")
def post_model(markata: "Markata") -> None:
    markata.post_models.append(SiteMapPost)


@hook_impl
@register_attr("sitemap")
def render(markata: Markata) -> None:

    sitemap = {
        "urlset": markata.map("post.sitemap_url.dict()", filter="post.published")
    }

    sitemap = (
        anyconfig.dumps(sitemap, "xml")
        .decode("utf-8")
        .replace(
            "<urlset>",
            (
                '<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"\n'
                'xmlns:news="http://www.google.com/schemas/sitemap-news/0.9"\n'
                'xmlns:xhtml="http://www.w3.org/1999/xhtml"\n'
                'xmlns:mobile="http://www.google.com/schemas/sitemap-mobile/1.0"\n'
                'xmlns:image="http://www.google.com/schemas/sitemap-image/1.1"\n'
                'xmlns:video="http://www.google.com/schemas/sitemap-video/1.1">\n'
            ),
        )
        .replace("</url>", "</url>\n")
    )
    markata.sitemap = sitemap


@hook_impl
def save(markata: Markata) -> None:
    with open(Path(markata.config["output_dir"]) / "sitemap.xml", "w") as f:
        f.write(markata.sitemap)  # type: ignore
