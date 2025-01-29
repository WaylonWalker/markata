"""
The `markata.plugins.post_model` plugin defines the core Post model used throughout
Markata. It provides robust validation, serialization, and configuration options for
all post attributes.

# Installation

This plugin is built-in and enabled by default through the 'default' plugin.
If you want to be explicit, you can add it to your list of plugins:

```toml
hooks = [
    "markata.plugins.post_model",
]
```

# Uninstallation

Since this plugin is included in the default plugin set, to disable it you must explicitly
add it to the disabled_hooks list if you are using the 'default' plugin:

```toml
disabled_hooks = [
    "markata.plugins.post_model",
]
```

Note: Disabling this plugin will break most of Markata's functionality as the Post
model is fundamental to the system.

# Configuration

Configure post model behavior in your `markata.toml`:

```toml
[markata.post_model]
# Attributes to include when serializing posts
include = [
    "date",
    "description",
    "published",
    "slug",
    "title",
    "content",
    "html"
]

# Attributes to show in post representations
repr_include = [
    "date",
    "description",
    "published",
    "slug",
    "title"
]

# Attributes to include when exporting
export_include = [
    "date",
    "description",
    "published",
    "slug",
    "title"
]
```

# Functionality

## Post Model

Core attributes:
- `path`: Path to source file
- `slug`: URL-friendly identifier
- `href`: Full URL path
- `published`: Publication status
- `description`: Post summary
- `content`: Raw markdown content
- `html`: Rendered HTML content
- `tags`: List of post tags
- `date`: Publication date
- `title`: Post title

## Validation

The model provides:
- Type checking and coercion
- Required field validation
- Custom field validators
- Default values
- Rich error messages

## Serialization

Supports multiple output formats:
- Full serialization (all fields)
- Representation (subset for display)
- Export (subset for external use)
- JSON/YAML compatible

## Performance

Uses optimized Pydantic config:
- Disabled assignment validation
- Arbitrary types allowed
- Extra fields allowed
- String whitespace stripping
- Default value validation
- Number to string coercion
- Alias population

## Dependencies

This plugin depends on:
- pydantic for model definition
- rich for console output
- yaml for YAML handling
"""

import datetime
import logging
from pathlib import Path
from rich.jupyter import JupyterMixin
from rich.pretty import Pretty
from typing import Any, Dict, List, Optional, TYPE_CHECKING, Union

import pydantic
from pydantic import ConfigDict, Field, field_validator
import yaml

from markata.hookspec import hook_impl, register_attr

logger = logging.getLogger(__name__)

if TYPE_CHECKING:
    from pydantic.typing import ReprArgs

    from markata import Markata


class Post(pydantic.BaseModel, JupyterMixin):
    markata: Any = Field(None, exclude=True)
    path: Path
    slug: Optional[str] = None
    href: Optional[str] = None
    published: bool = False
    description: Optional[str] = None
    content: str = ""  # Default to empty string instead of None
    tags: List[str] = Field(default_factory=list)
    raw_date: Optional[Any] = Field(
        None, alias="date", exclude=True
    )  # Accept any type for raw_date
    date: Optional[datetime.date] = None
    date_time: Optional[datetime.datetime] = None
    today: datetime.date = pydantic.Field(default_factory=datetime.date.today)
    now: datetime.datetime = pydantic.Field(default_factory=datetime.datetime.utcnow)
    load_time: float = 0
    profile: Optional[str] = None
    title: str = None
    template: Optional[str | Dict[str, str]] = "post.html"
    sidebar: Optional[Any] = None

    model_config = ConfigDict(
        validate_assignment=False,  # Skip validation on assignment for performance
        arbitrary_types_allowed=True,
        extra="allow",
        str_strip_whitespace=True,
        validate_default=True,
        coerce_numbers_to_str=True,
        populate_by_name=True,
    )

    def __rich__(self) -> Pretty:
        return Pretty(self)

    def __repr_args__(self: "Post") -> "ReprArgs":
        return [
            (key, value)
            for key, value in self.__dict__.items()
            if key in self.markata.config.post_model.repr_include
        ]

    @property
    def key(self: "Post") -> List[str]:
        return self.markata.make_hash(
            self.slug,
            self.href,
            self.published,
            self.description,
            self.content,
            self.date,
            self.title,
            self.template,
        )

    @property
    def metadata(self: "Post") -> Dict:
        "for backwards compatability"
        return self.__dict__

    def to_dict(self: "Post") -> Dict:
        "for backwards compatability"
        return self.__dict__

    def __getitem__(self: "Post", item: str) -> Any:
        "for backwards compatability"
        return getattr(self, item)

    def __setitem__(self: "Post", key: str, item: Any) -> None:
        "for backwards compatability"
        setattr(self, key, item)

    def get(self: "Post", item: str, default: Any) -> Any:
        "for backwards compatability"
        return getattr(self, item, default)

    def keys(self: "Post") -> List[str]:
        "for backwards compatability"
        return self.__dict__.keys()

    # def json(
    #     self: "Post",
    #     include: Iterable = None,
    #     all: bool = False,
    #     **kwargs,
    # ) -> str:
    #     """
    #     override function to give a default include value that will include
    #     user configured includes.
    #     """
    #     if all:
    #         return pydantic.create_model("Post", **self)(**self).json(
    #             **kwargs,
    #         )
    #     if include:
    #         return pydantic.create_model("Post", **self)(**self).json(
    #             include=include,
    #             **kwargs,
    #         )
    #     return pydantic.create_model("Post", **self)(**self).json(
    #         include={i: True for i in self.markata.config.post_model.include},
    #         **kwargs,
    #     )

    def yaml(self: "Post") -> str:
        """
        dump model to yaml
        """
        import yaml

        return yaml.dump(
            self.dict(
                include={i: True for i in self.markata.config.post_model.include}
            ),
            Dumper=yaml.CDumper,
        )

    def markdown(self: "Post") -> str:
        """
        dump model to markdown
        """

        import yaml

        frontmatter = yaml.dump(
            self.dict(
                include={
                    i: True
                    for i in [
                        _i
                        for _i in self.markata.config.post_model.include
                        if _i != "content"
                    ]
                }
            ),
            Dumper=yaml.CDumper,
        )
        post = "---\n"
        post += frontmatter
        post += "---\n\n"

        if self.content:
            post += self.content
        return post

    @classmethod
    def parse_file(cls, markata, path: Union[Path, str], **kwargs) -> "Post":
        if isinstance(path, Path):
            if path.suffix in [".md", ".markdown"]:
                return cls.parse_markdown(markata=markata, path=path, **kwargs)
        elif isinstance(path, str):
            if path.endswith(".md") or path.endswith(".markdown"):
                return cls.parse_markdown(markata=markata, path=path, **kwargs)
        return super(Post, cls).parse_file(path, **kwargs)

    @classmethod
    def parse_markdown(cls, markata, path: Union[Path, str], **kwargs) -> "Post":
        if isinstance(path, str):
            path = Path(path)
        text = path.read_text()
        try:
            _, fm, *content = text.split("---\n")
            content = "---\n".join(content)
            try:
                fm = yaml.load(fm, Loader=yaml.CBaseLoader)
            except yaml.YAMLError:
                fm = {}
        except ValueError:
            fm = {}
            content = text
        if fm is None or isinstance(fm, str):
            fm = {}

        post_args = {
            "markata": markata,
            "path": path,
            "content": content,
            "raw": text,
            **fm,
        }

        return markata.Post.parse_obj(post_args)

    def dumps(self):
        """
        dumps raw article back out
        """
        return f"---\n{self.yaml()}\n\n---\n\n{self.content}"

    @field_validator("slug", mode="before")
    def default_slug(cls, v, info):
        from slugify import slugify

        return v or slugify(str(info.data.get("path", "").stem))

    @field_validator("slug", mode="before")
    def index_slug_is_empty(cls, v, info):
        if v == "index":
            return ""
        return v

    @field_validator("slug", mode="before")
    def no_double_slash_in_slug(cls, v, info):
        if v is None:
            return v
        return v.replace("//", "/")

    @field_validator("href", mode="before")
    def default_href(cls, v, info):
        if v:
            return v
        return f"/{info.data.get('slug', '').strip('/')}/".replace("//", "/")

    @field_validator("title", mode="before")
    def title_title(cls, v, info):
        if v:
            return v
        return info.data.get("path", "").stem.replace("-", " ").title()

    @field_validator("description", mode="before")
    def default_description(cls, v, info):
        if v:
            return v
        content = info.data.get("content", "")
        if content:
            return " ".join(content.split()[:50])
        return None

    @field_validator("tags", mode="before")
    def tags_not_none(cls, v, info):
        if v is None:
            return []
        return v

    @field_validator("date", mode="before")
    def date_is_date(cls, v, info):
        # print(f"date_is_date received value: {v!r} of type {type(v)}")  # Debug log
        if v is None:
            # return None
            return info.data.get("markata").config.default_date

        if isinstance(v, datetime.date) and not isinstance(v, datetime.datetime):
            return v
        if isinstance(v, datetime.datetime):
            # Ensure zero time when converting datetime to date
            zero_time = v.replace(hour=0, minute=0, second=0, microsecond=0)
            return zero_time.date()
        if isinstance(v, str):
            try:
                # Try ISO format first
                return datetime.datetime.fromisoformat(v.replace("Z", "+00:00")).date()
            except ValueError:
                try:
                    # Try parsing with time
                    dt = datetime.datetime.strptime(v, "%Y-%m-%d %H:%M")
                    # Ensure zero time
                    zero_time = dt.replace(hour=0, minute=0, second=0, microsecond=0)
                    return zero_time.date()
                except ValueError:
                    try:
                        # Fallback to date only
                        return datetime.datetime.strptime(v, "%Y-%m-%d").date()
                    except ValueError:
                        return datetime.date.today()
        # If we get here, try to convert to string
        try:
            return cls.date_is_date(str(v), info)
        except (ValueError, TypeError):
            print(f"Failed to convert {v!r} to date")  # Debug log
            return datetime.date.today()

    @field_validator("date_time", mode="before")
    def parse_date_time(cls, v, info):
        """Single validator to handle all date_time parsing cases"""
        # If we have an explicit date_time value
        if v is not None:
            if isinstance(v, datetime.datetime):
                return v
            if isinstance(v, datetime.date):
                return datetime.datetime.combine(v, datetime.time.min)
            if isinstance(v, str):
                try:
                    # Try ISO format first
                    return datetime.datetime.fromisoformat(v.replace("Z", "+00:00"))
                except ValueError:
                    try:
                        return datetime.datetime.strptime(v, "%Y-%m-%d %H:%M")
                    except ValueError:
                        try:
                            return datetime.datetime.strptime(v, "%Y-%m-%d")
                        except ValueError:
                            # Try dateparser as last resort for explicit date_time
                            import dateparser

                            parsed = dateparser.parse(v)
                            if parsed:
                                return parsed
                            return datetime.datetime.now()

        # Get the raw date string directly from raw_date field
        raw_date = info.data.get("raw_date")
        if raw_date and isinstance(raw_date, str):
            try:
                # Try ISO format first
                return datetime.datetime.fromisoformat(raw_date.replace("Z", "+00:00"))
            except ValueError:
                try:
                    # Try parsing raw_date with time first
                    return datetime.datetime.strptime(raw_date, "%Y-%m-%d %H:%M")
                except ValueError:
                    try:
                        # Fallback to date only
                        return datetime.datetime.strptime(raw_date, "%Y-%m-%d")
                    except ValueError:
                        # Try dateparser as last resort
                        import dateparser

                        parsed = dateparser.parse(raw_date)
                        if parsed:
                            return parsed

        # If no raw_date, try to derive from date field
        date = info.data.get("date")
        if date:
            if isinstance(date, datetime.datetime):
                return date
            if isinstance(date, str):
                try:
                    # Try ISO format first
                    return datetime.datetime.fromisoformat(date.replace("Z", "+00:00"))
                except ValueError:
                    try:
                        # Try parsing date with time first
                        return datetime.datetime.strptime(date, "%Y-%m-%d %H:%M")
                    except ValueError:
                        try:
                            # Fallback to date only
                            return datetime.datetime.strptime(date, "%Y-%m-%d")
                        except ValueError:
                            # Try dateparser as last resort
                            import dateparser

                            parsed = dateparser.parse(date)
                            if parsed:
                                return parsed
            if isinstance(date, datetime.date):
                return datetime.datetime.combine(date, datetime.time.min)

        # If we still don't have a date, use now
        return datetime.datetime.now()


class PostModelConfig(pydantic.BaseModel):
    """Configuration for the Post model"""

    def __init__(self, **data) -> None:
        """

        include: post attributes to include by default in Post
        model serialization.
        repr_include: post attributes to include by default in Post
        repr.  If `repr_include` is None, it will default to
        `include`, but it is likely that you want less in the repr
        than serialized output.

        example:

        ``` toml title='markata.toml'
        [markata.post_model]
        include = ['date', 'description', 'published',
            'slug', 'title', 'content', 'html']
        repr_include = ['date', 'description', 'published', 'slug', 'title']
        ```
        """
        super().__init__(**data)

    default_date: datetime.date = datetime.date.today()
    include: List[str] = [
        "date",
        "description",
        "published",
        "slug",
        "title",
        "content",
        "html",
    ]
    repr_include: Optional[List[str]] = [
        "date",
        "description",
        "published",
        "slug",
        "title",
    ]
    export_include: Optional[List[str]] = [
        "date",
        "description",
        "published",
        "slug",
        "title",
    ]

    model_config = ConfigDict(
        validate_assignment=True,  # Config model
        arbitrary_types_allowed=True,
        extra="allow",
        str_strip_whitespace=True,
        validate_default=True,
        coerce_numbers_to_str=True,
        populate_by_name=True,
    )

    @field_validator("repr_include", mode="before")
    @classmethod
    def repr_include_validator(cls, v, info) -> Optional[List[str]]:
        if v:
            return v
        return info.data.get("include")


class Config(pydantic.BaseModel):
    post_model: PostModelConfig = pydantic.Field(default_factory=PostModelConfig)


@hook_impl(trylast=True)
@register_attr("post_models")
def post_model(markata: "Markata") -> None:
    markata.post_models.append(Post)


@hook_impl(tryfirst=True)
def config_model(markata: "Markata") -> None:
    markata.config_models.append(Config)


# from polyfactory.factories.pydantic_factory import ModelFactory
# class PostFactory(ModelFactory):
#     __model__ = Post
#     __model__ = Post
