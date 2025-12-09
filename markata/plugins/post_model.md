---
date: 2025-12-09
description: "The plugin defines the core Post model used throughout Markata. It provides
  robust validation, serialization, and configuration options for all post\u2026"
published: false
slug: markata/plugins/post-model
title: post_model.py


---

---

The `markata.plugins.post_model` plugin defines the core Post model used throughout
Markata. It provides robust validation, serialization, and configuration options for
all post attributes.

## Installation

This plugin is built-in and enabled by default through the 'default' plugin.
If you want to be explicit, you can add it to your list of plugins:

```toml
hooks = [
    "markata.plugins.post_model",
]
```

## Uninstallation

Since this plugin is included in the default plugin set, to disable it you must explicitly
add it to the disabled_hooks list if you are using the 'default' plugin:

```toml
disabled_hooks = [
    "markata.plugins.post_model",
]
```

Note: Disabling this plugin will break most of Markata's functionality as the Post
model is fundamental to the system.

## Configuration

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

## Functionality

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

---

!!! class
    <h2 id="PostModelConfig" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">PostModelConfig <em class="small">class</em></h2>

    Configuration for the Post model

???+ source "PostModelConfig <em class='small'>source</em>"
    ```python
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
    ```
!!! method
    <h2 id="metadata" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">metadata <em class="small">method</em></h2>

    for backwards compatability

???+ source "metadata <em class='small'>source</em>"
    ```python
    def metadata(self: "Post") -> Dict:
            "for backwards compatability"
            return self.__dict__
    ```
!!! method
    <h2 id="to_dict" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">to_dict <em class="small">method</em></h2>

    for backwards compatability

???+ source "to_dict <em class='small'>source</em>"
    ```python
    def to_dict(self: "Post") -> Dict:
            "for backwards compatability"
            return self.__dict__
    ```
!!! method
    <h2 id="__getitem__" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">__getitem__ <em class="small">method</em></h2>

    for backwards compatability

???+ source "__getitem__ <em class='small'>source</em>"
    ```python
    def __getitem__(self: "Post", item: str) -> Any:
            "for backwards compatability"
            return getattr(self, item)
    ```
!!! method
    <h2 id="__setitem__" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">__setitem__ <em class="small">method</em></h2>

    for backwards compatability

???+ source "__setitem__ <em class='small'>source</em>"
    ```python
    def __setitem__(self: "Post", key: str, item: Any) -> None:
            "for backwards compatability"
            setattr(self, key, item)
    ```
!!! method
    <h2 id="get" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">get <em class="small">method</em></h2>

    for backwards compatability

???+ source "get <em class='small'>source</em>"
    ```python
    def get(self: "Post", item: str, default: Any) -> Any:
            "for backwards compatability"
            return getattr(self, item, default)
    ```
!!! method
    <h2 id="keys" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">keys <em class="small">method</em></h2>

    for backwards compatability

???+ source "keys <em class='small'>source</em>"
    ```python
    def keys(self: "Post") -> List[str]:
            "for backwards compatability"
            return self.__dict__.keys()
    ```
!!! method
    <h2 id="yaml" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">yaml <em class="small">method</em></h2>

    dump model to yaml

???+ source "yaml <em class='small'>source</em>"
    ```python
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
    ```
!!! method
    <h2 id="markdown" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">markdown <em class="small">method</em></h2>

    dump model to markdown

???+ source "markdown <em class='small'>source</em>"
    ```python
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
    ```
!!! method
    <h2 id="dumps" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">dumps <em class="small">method</em></h2>

    dumps raw article back out

???+ source "dumps <em class='small'>source</em>"
    ```python
    def dumps(self):
            """
            dumps raw article back out
            """
            return f"---\n{self.yaml()}\n\n---\n\n{self.content}"
    ```
!!! method
    <h2 id="parse_date_time" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">parse_date_time <em class="small">method</em></h2>

    Single validator to handle all date_time parsing cases

???+ source "parse_date_time <em class='small'>source</em>"
    ```python
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
    ```
!!! method
    <h2 id="__init__" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">__init__ <em class="small">method</em></h2>

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

???+ source "__init__ <em class='small'>source</em>"
    ```python
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
    ```