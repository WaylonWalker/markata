from pathlib import Path
from typing import Optional, Any

import pydantic
from pydantic import Field

from markata import Markata
from markata.hookspec import hook_impl, register_attr


class SiteMapUrl(pydantic.BaseModel):
    """A model representing a URL entry in the sitemap.xml file.

    To configure the base URL for your site, set the 'url' field in your markata config:
    ```yaml
    url: https://example.com
    ```

    If no base URL is set, relative URLs will be used.
    """

    slug: str = Field(..., exclude=True)
    loc: str = Field(
        None,
        include=True,
        description="The full URL for this page in the sitemap. Generated automatically from config.url + slug.",
    )
    changefreq: str = Field("daily", include=True)
    priority: str = Field("0.7", include=True)
    markata: Any = Field(None, exclude=True)

    model_config = pydantic.ConfigDict(
        validate_assignment=False,
        arbitrary_types_allowed=True,
        extra="allow",
        str_strip_whitespace=True,
        validate_default=True,
        coerce_numbers_to_str=True,
        populate_by_name=True,
    )

    @pydantic.field_validator("loc", mode="before")
    @classmethod
    def validate_loc(cls, v, info) -> str:
        """Generate the URL for the sitemap entry.

        Uses markata.config.url as the base URL if set, otherwise uses relative URLs.
        Example: https://example.com/my-page/ or /my-page/
        """
        if v is None:
            markata = info.data.get("markata")
            slug = info.data.get("slug")
            if markata is None or slug is None:
                raise ValueError(
                    "Could not generate sitemap URL: markata and slug are required. "
                    "This usually means the Post model is missing required fields. "
                    "Check that your post has a valid slug and markata instance."
                )

            # Get base URL from config, default to empty string if not set
            base_url = getattr(markata.config, "url", "")
            if not base_url:
                return f"/{slug}/"

            # Ensure URL has a trailing slash for consistency
            return f"{base_url.rstrip('/')}/{slug}/"
        return v

    def dict(self, *args, **kwargs):
        return {"url": {**super().dict(*args, **kwargs)}}


class SiteMapPost(pydantic.BaseModel):
    """A model for posts that will be included in the sitemap.

    To configure the base URL for your site, set the 'url' field in your markata config:
    ```yaml
    url: https://example.com
    ```
    """

    slug: str = None
    published: bool = True
    sitemap_url: Optional[SiteMapUrl] = None
    markata: Any = Field(None, exclude=True)

    model_config = pydantic.ConfigDict(
        validate_assignment=False,
        arbitrary_types_allowed=True,
        extra="allow",
        str_strip_whitespace=True,
        validate_default=True,
        coerce_numbers_to_str=True,
        populate_by_name=True,
    )

    @pydantic.field_validator("sitemap_url", mode="before")
    @classmethod
    def validate_sitemap_url(cls, v, info) -> Optional[SiteMapUrl]:
        """Initialize sitemap_url if not provided."""
        markata = info.data.get("markata")
        slug = info.data.get("slug")
        if markata is None or slug is None:
            raise ValueError(
                "Could not create sitemap entry: markata and slug are required. "
                "This usually means the Post model is missing required fields. "
                "Check that your post has a valid slug and markata instance."
            )

        if v is None:
            return SiteMapUrl(markata=markata, slug=slug)
        if isinstance(v, dict):
            return SiteMapUrl(**v, markata=markata, slug=slug)
        if v.markata is None:
            v.markata = markata
        if v.slug is None:
            v.slug = slug
        return v


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
