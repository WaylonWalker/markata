import pydantic


class Post(pydantic.BaseModel):
    config_overrides: PostOverrides = PostOverrides()
    template: Optional[str | Dict[str, str]] = None

    @pydantic.validator("template", pre=True, always=True)
    def default_template(cls, v, *, values):
        if v is None:
            return values["markata"].config.post_template
        if isinstance(v, str):
            v = {"index": v}
        if isinstance(values["markata"].config.post_template, str):
            config_template = {
                "index": values["markata"].config.post_template,
            }
        else:
            config_template = values["markata"].config.post_template
        return {**config_template, **v}
