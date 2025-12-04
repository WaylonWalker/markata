"""
The `markata.plugins.jinja_md` plugin enables Jinja2 templating within your markdown
content. This allows you to dynamically generate content using Python expressions and
access to the full Markata context.

## Installation

This plugin is built-in and enabled by default through the 'default' plugin.
If you want to be explicit, you can add it to your list of plugins:

```toml
hooks = [
    "markata.plugins.jinja_md",
]
```

## Uninstallation

Since this plugin is included in the default plugin set, to disable it you must explicitly
add it to the disabled_hooks list if you are using the 'default' plugin:

```toml
disabled_hooks = [
    "markata.plugins.jinja_md",
]
```

## Configuration

Configure Jinja markdown settings in your `markata.toml`:

```toml
[markata.jinja_md]
# List of files to ignore for Jinja processing
ignore = [
    "README.md",
    "CHANGELOG.md"
]
```

## Functionality

## Template Variables

Your markdown has access to:
- `post`: The current post being rendered
- `markata`: The Markata instance with all configuration and posts

## Example Usage

### Access Post Metadata
```markdown
# {{ post.title }}
{{ post.description }}

Published on: {{ post.date.strftime('%Y-%m-%d') }}
```

### Generate Link Lists
```markdown
{# One-liner list of all posts #}
{{ '\\n'.join(markata.map('f"* [{title}]({slug})"', sort='slug')) }}

{# For-loop with filtering #}
{% for post in markata.map('post', filter='"git" in tags') %}
* [{{ post.title }}]({{ post.slug }})
{% endfor %}
```

### Include Raw Files
```markdown
{# Include file contents without processing #}
{% include_raw 'code/example.py' %}
```

## Jinja Extensions

The plugin supports:
1. Custom extensions via entrypoints
2. Built-in extensions like include_raw
3. Silent undefined variables
4. Custom error messages

## Error Handling

The plugin provides:
- Custom error messages for template syntax errors
- Silent handling of undefined variables
- Detailed error reporting with file and line info

## Dependencies

This plugin depends on:
- Jinja2 for template processing
- The `render_markdown` plugin for final HTML rendering
"""

from pathlib import Path
from typing import TYPE_CHECKING
from typing import Any
from typing import Dict
from typing import List
from typing import Optional

import pathspec
import pydantic
from jinja2 import TemplateSyntaxError
from jinja2 import Undefined
from jinja2 import UndefinedError
from jinja2 import nodes
from jinja2.ext import Extension
from pydantic import Field

from markata import __version__
from markata.hookspec import hook_impl
from markata.hookspec import register_attr


def register_jinja_extensions(config: dict) -> List[Extension]:
    """
    Gets jinja extensions from entrypoints and loads them in.

    Returns: List of jinja Extensions
    """

    import pkg_resources

    return [
        ep.load() for ep in pkg_resources.iter_entry_points(group="markata.jinja_md")
    ]


class IncludeRawExtension(Extension):
    tags = {"include_raw"}

    def parse(self, parser):
        line_number = next(parser.stream).lineno
        file = [parser.parse_expression()]
        return nodes.CallBlock(
            self.call_method("_read_file", file),
            [],
            [],
            "",
        ).set_lineno(line_number)

    def _read_file(self, file, caller):
        return Path(file).read_text()


if TYPE_CHECKING:
    from markata import Markata


class _SilentUndefined(Undefined):
    """
    silence undefined variable errors in jinja templates.

    # Example
    ```python
    template = '{{ variable }}'
    post.content = Template( template, undefined=_SilentUndefined).render()
    ```
    """

    def _fail_with_undefined_error(self, *args, **kwargs):
        return ""


class PostTemplateSyntaxError(TemplateSyntaxError):
    """
    Custom error message for post template syntax errors.
    """


class JinjaMd(pydantic.BaseModel):
    markata: Any = Field(None, exclude=True)
    content: str = ""
    article_html: Optional[Dict[str, str] | str] = None
    jinja: bool = Field(default=False)
    model_config = pydantic.ConfigDict(
        validate_assignment=False,
        arbitrary_types_allowed=True,
        extra="allow",
        str_strip_whitespace=True,
        validate_default=True,
        coerce_numbers_to_str=True,
        populate_by_name=True,
    )


@hook_impl()
@register_attr("post_models")
def post_model(markata: "Markata") -> None:
    markata.post_models.append(JinjaMd)


@hook_impl()
@register_attr("prevnext")
def pre_render(markata: "Markata") -> None:
    """
    jinja_md hook for markata to render your markdown post as a jinja template.

    The post itself is exposed as `post`, and the markata instance is exposed
    as `markata`.
    """

    config = markata.config.jinja_md
    ignore_spec = pathspec.PathSpec.from_lines("gitwildmatch", config.ignore)
    # for post in markata.iter_articles(description="jinja_md"):

    # jinja_env = jinja2.Environment(
    #     extensions=[IncludeRawExtension, *register_jinja_extensions(config)],
    # )
    jinja_env = markata.jinja_env
    jinja_env.undefined = _SilentUndefined

    for post in markata.filter("jinja==True"):
        if post.get("jinja", True) and not ignore_spec.match_file(post["path"]):
            try:
                # Include post metadata and markata version in cache key
                # since these affect the rendered output
                key = markata.make_hash(
                    "jinja_md", 
                    "pre_render", 
                    post.content,
                    str(post.to_dict()),  # Include all post metadata
                    __version__,  # Include markata version
                )
                content_from_cache = markata.precache.get(key)
                if content_from_cache is None and post.content is not None:
                    post.content = jinja_env.from_string(post.content).render(
                        __version__=__version__,
                        markata=markata,
                        body=post.article_html,
                        config=markata.config,
                        post=post,
                    )
                    with markata.cache:
                        markata.cache.set(key, post.content)
                else:
                    post.content = content_from_cache
                # prevent double rendering
                post.jinja = False
            except TemplateSyntaxError as e:
                errorline = post.content.split("\n")[e.lineno - 1]
                msg = f"""
                Error while processing post {post["path"]}

                {errorline}
                """

                raise PostTemplateSyntaxError(msg, lineno=e.lineno)
            except UndefinedError as e:
                raise UndefinedError(f"{e} in {post['path']}")


class JinjaMdConfig(pydantic.BaseModel):
    ignore: List[str] = []


class Config(pydantic.BaseModel):
    jinja_md: JinjaMdConfig = JinjaMdConfig()


@hook_impl(tryfirst=True)
def config_model(markata: "Markata") -> None:
    markata.config_models.append(Config)
