import datetime
import logging
from pathlib import Path
from typing import TYPE_CHECKING, Any, Dict, List, Optional, Union

import dateparser
import pydantic
import yaml
from polyfactory.factories.pydantic_factory import ModelFactory
from pydantic import ConfigDict
from slugify import slugify

from markata.hookspec import hook_impl, register_attr

logger = logging.getLogger(__name__)

if TYPE_CHECKING:
    from pydantic.typing import ReprArgs

    from markata import Markata


class Post(pydantic.BaseModel):
    markata: Any = None
    path: Path
    slug: Optional[str] = None
    href: Optional[str] = None
    published: bool = False
    description: Optional[str] = None
    content: str = None
    # date: Union[datetime.date, str]=None
    date: Optional[Union[datetime.date, str]] = None
    # pydantic.Field(
    # default_factory=lambda: datetime.date.min
    # )
    date_time: Optional[datetime.datetime] = None
    today: datetime.date = pydantic.Field(default_factory=datetime.date.today)
    now: datetime.datetime = pydantic.Field(default_factory=datetime.datetime.utcnow)
    load_time: float = 0
    profile: Optional[str] = None
    title: str = None
    model_config = ConfigDict(
        validate_assignment=True,
        arbitrary_types_allowed=True,
        extra="allow",
    )

    def __repr_args__(self: "Post") -> "ReprArgs":
        return [
            (key, value)
            for key, value in self.__dict__.items()
            if key in self.markata.config.post_model.repr_include
        ]

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
            **fm,
        }

        return markata.Post(**post_args)

    def dumps(self):
        """
        dumps raw article back out
        """
        return f"---\n{self.yaml()}\n\n---\n\n{self.content}"

    @pydantic.validator("slug", pre=True, always=True)
    def default_slug(cls, v, *, values):
        return v or slugify(str(values["path"].stem))

    @pydantic.validator("slug", pre=True, always=True)
    def index_slug_is_empty(cls, v, *, values):
        if v == "index":
            return ""
        return v

    @pydantic.validator("href", pre=True, always=True)
    def default_href(cls, v, *, values):
        if v:
            return v
        return f"/{values['slug'].strip('/')}/".replace("//", "/")

    @pydantic.validator("title", pre=True, always=True)
    def title_title(cls, v, *, values):
        title = v or Path(values["path"]).stem.replace("-", " ")
        return title.title()

    @pydantic.validator("date_time", pre=True, always=True)
    def dateparser_datetime(cls, v, *, values):
        if isinstance(v, str):
            d = dateparser.parse(v)
            if d is None:
                raise ValueError(f'"{v}" is not a valid date')
        return v

    @pydantic.validator("date_time", pre=True, always=True)
    def date_is_datetime(cls, v, *, values):
        if v is None and "date" not in values:
            values["markata"].console.log(f"{values['path']} has no date")
            return datetime.datetime.now()
        if v is None and values["date"] is None:
            values["markata"].console.log(f"{values['path']} has no date")
            return datetime.datetime.now()
        if isinstance(v, datetime.datetime):
            return v
        if isinstance(values["date"], datetime.datetime):
            return values["date"]
        if isinstance(v, datetime.date):
            return datetime.datetime.combine(v, datetime.time.min)
        if isinstance(values["date"], datetime.date):
            return datetime.datetime.combine(values["date"], datetime.time.min)
        return v

    @pydantic.validator("date_time", pre=True, always=True)
    def mindate_time(cls, v, *, values):
        if v is None and "date" not in values:
            values["markata"].console.log(f"{values['path']} has no date")
            return datetime.datetime.min
        if values["date"] is None:
            values["markata"].console.log(f"{values['path']} has no date")
            return datetime.datetime.min
        if isinstance(v, datetime.datetime):
            return v
        if isinstance(values["date"], datetime.datetime):
            return values["date"]
        if isinstance(v, datetime.date):
            return datetime.datetime.combine(v, datetime.time.min)
        if isinstance(values["date"], datetime.date):
            return datetime.datetime.combine(values["date"], datetime.time.min)
        return v

    @pydantic.validator("date", pre=True, always=True)
    def dateparser_date(cls, v, *, values):
        if v is None:
            return datetime.date.min
        if isinstance(v, str):
            d = cls.markata.precache.get(v)
            if d is not None:
                return d
            d = dateparser.parse(v)
            if d is None:
                raise ValueError(f'"{v}" is not a valid date')
            d = d.date()
            with cls.markata.cache as cache:
                cache.add(v, d)
            return d
        return v

    # @pydantic.validator("date", pre=True, always=True)
    # def datetime_is_date(cls, v, *, values):
    #     if isinstance(v, datetime.date):
    #         return v
    #     if isinstance(v, datetime.datetime):
    #         return v.date()

    # @pydantic.validator("date", pre=True, always=True)
    # def mindate(cls, v, *, values):
    #     if v is None:
    #         return datetime.date.min
    #     return v


class PostModelConfig(pydantic.BaseModel):
    "Configuration for the Post model"

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

    @pydantic.validator("repr_include", pre=True, always=True)
    def repr_include_validator(cls, v, *, values):
        if v:
            return v
        return values.get("include", None)


class Config(pydantic.BaseModel):
    post_model: PostModelConfig = pydantic.Field(default_factory=PostModelConfig)


@hook_impl(trylast=True)
@register_attr("post_models")
def post_model(markata: "Markata") -> None:
    markata.post_models.append(Post)


@hook_impl(tryfirst=True)
def config_model(markata: "Markata") -> None:
    markata.config_models.append(Config)


class PostFactory(ModelFactory):
    __model__ = Post
