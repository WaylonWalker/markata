from pathlib import Path
from typing import Optional

import pydantic

from markata import Markata
from markata.hookspec import hook_impl, register_attr


class SiteMapUrl(pydantic.BaseModel):
    slug: str = pydantic.Field(..., exclude=True)
    loc: str = pydantic.Field(None, include=True)
    changefreq: str = pydantic.Field("daily", include=True)
    priority: str = pydantic.Field("0.7", include=True)

    _loc = pydantic.field_validator("loc", pre=False, always=True)(
        lambda v, *, values, **kwargs: values["markata"].config.url
        + "/"
        + values["slug"]
        + "/"
        if v is None
        else v
    )

    def dict(self, *args, **kwargs):
        return {"url": {**super().dict(*args, **kwargs)}}


class SiteMapPost(pydantic.BaseModel):
    slug: str = None
    published: bool = True
    sitemap_url: Optional[SiteMapUrl] = None

    _sitemap_url = pydantic.field_validator("sitemap_url", mode="after")(
        lambda v, *, values, **kwargs: SiteMapUrl(
            markata=values["markata"], slug=values["slug"]
        )
        if v is None
        else SiteMapUrl(markata=values["markata"], slug=values["slug"])
        if v.markata is None
        else v
    )


@hook_impl()
@register_attr("post_models")
def post_model(markata: "Markata") -> None:
    markata.post_models.append(SiteMapPost)


@hook_impl
@register_attr("sitemap")
def render(markata: Markata) -> None:
    import anyconfig

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
    with open(Path(markata.config.output_dir) / "sitemap.xml", "w") as f:
        f.write(markata.sitemap)  # type: ignore
