import copy
import datetime
from pathlib import Path
from typing import TYPE_CHECKING, Any, List, Optional

import pydantic
import yaml
from dataclasses_json import config
from deepmerge import always_merger
from marshmallow import fields
from slugify import slugify

from markata import Markata
from markata.hookspec import hook_impl, register_attr


class Post(pydantic.BaseModel):
    path: Path
    published: bool = False
    slug: str = None
    description: Optional[str] = None
    content: str = None
    date: datetime.date = pydantic.Field(default_factory=datetime.date.today)
    today: datetime.date = pydantic.Field(default_factory=datetime.date.today)
    now: datetime.datetime = pydantic.Field(default_factory=datetime.datetime.utcnow)
    title: Optional[str] = None
    markata: Markata = None

    class Config:
        arbitrary_types_allowed = True
        # json_encoders = {datetime.datetime: lambda x: x.isoformat()}

    def __repr_args__(self) -> "ReprArgs":
        return [
            (key, value)
            for key, value in self.__dict__.items()
            if key in self.markata.config.post_model.repr_include
        ]

    @property
    def metadata(self):
        "for backwards compatability"
        return self.__dict__

    def to_dict(self):
        "for backwards compatability"
        return self.__dict__

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

    def json(self, include=None, all=False, **kwargs):
        """
        override function to give a default include value that will include
        user configured includes.
        """
        if all:
            return pydantic.create_model("Post", **self,)(**self).json(
                **kwargs,
            )
        if include:
            return pydantic.create_model("Post", **self,)(**self).json(
                include=include,
                **kwargs,
            )
        return pydantic.create_model("Post", **self,)(**self).json(
            include={i: True for i in self.markata.config.post_model.include},
            **kwargs,
        )

    def yaml(self):
        """
        dump model to yaml
        """
        return yaml.dump(
            self.dict(include={i: True for i in self.markata.config.post_model.include})
        )

    # def __init__(self, **data) -> None:
    #     _full_config = copy.deepcopy(data.get("config", {}))
    #     data["config"] = dict(
    #         always_merger.merge(
    #             _full_config,
    #             copy.deepcopy(
    #                 data.get(
    #                     "config_overrides",
    #                     {},
    #                 ),
    #             ),
    #         ),
    #     )
    #     super().__init__(**data)

    def dumps(self):
        """
        dumps raw article back out
        """
        return f"{self.yaml()}\n\n---\n\n{self.content}"

    @pydantic.validator("title", pre=True, always=True)
    def title_title(cls, v, *, values):
        title = v or Path(values["path"]).stem
        return title.title()

    @pydantic.validator("slug", pre=True, always=True)
    def default_slug(cls, v, *, values):
        return v or slugify(str(values.get("path")))


class PostModelConfig(pydantic.BaseModel):
    "Configuration for the Post model"

    def __init__(self, **data):
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
        include = ['date', 'description', 'published', 'slug', 'title', 'content', 'html']
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


@hook_impl
@register_attr("post_models")
def post_model(markata: "Markata") -> None:
    markata.post_models.append(Post)


@hook_impl(tryfirst=True)
def config_model(markata: "Markata") -> None:
    markata.config_models.append(Config)
