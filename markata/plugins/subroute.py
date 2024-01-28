from pathlib import Path
from typing import Dict, Any

from rich.markdown import Markdown
import pydantic
from pydantic import Field
from pydantic import BaseModel

from markata import Markata
from markata.hookspec import hook_impl, register_attr
from pydantic import ConfigDict


class Config(BaseModel):
    subroute: Path = Path("")


class SubroutePost(BaseModel):
    markata: Any = Field(None, exclude=True)
    model_config = ConfigDict(validate_assignment=True, arbitrary_types_allowed=True)

    @pydantic.validator("slug")
    @classmethod
    def relative_to_subroute(cls, v, *, values: Dict) -> Path:
        subroute = cls.markata.config.subroute
        if subroute == Path(""):
            return v

        slug = Path(v)

        if not slug.relative_to(subroute):
            return slug.parent / subroute / slug.stem
        return v


@hook_impl()
@register_attr("config_models")
def config_model(markata: "Markata") -> None:
    markata.config_models.append(Config)


@hook_impl()
@register_attr("post_models")
def post_model(markata: Markdown) -> None:
    markata.post_models.append(SubroutePost)
