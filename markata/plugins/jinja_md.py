"""
Renders your markdown as a jinja template during pre_render.

The markata instance is passed into the template, giving you access to things
such as all of your articles, config, and this post as post.

# Examples

first we can grab a few things out of the frontmatter of this post.

``` markdown
# {{ post.title }}
{{ post.description }}
```

# one-liner list of links

This one-liner will render a list of markdown links into your markdown at build
time.  It's quite handy to pop into posts.

``` markdown
{{ '\\n'.join(markata.map('f"* [{title}]({slug})"', sort='slug')) }}
```

# jinja for to markdown list of links

Sometimes quoting things like your filters are hard to do in a one line without
running out of quote variants.  Jinja for loops can make this much easier.

``` markdown
{% for post in markata.map('post', filter='"git" in tags') %}
* [{{ post.title }}]({{ post.slug }})
{% endfor %}
```

# jinja for to html list of links

Since markdown is a superset of html, you can just render out html into your
post and it is still valid.

``` markdown
<ul>
{% for post in markata.map('post', filter='"git" in tags') %}
    <li><a href="{{ post.slug }}">{{ post.title }}</a></li>
{% endfor %}
</ul>
```

# Ignoring files

It is possible to ignore files by adding an ignore to your `markata.jinja_md`
config in your `markata.toml` file.  This ignore follows the `gitwildmatch`
rules, so think of it the same as writing a gitignore.

```toml
[markata.jinja_md]
ignore=[
'jinja_md.md',
]
```

!!note
  Docs such as this jinja_md.py file will get converted to jinja_md.md during
  build time, so use `.md` extensions instead of `.py`.

# Ignoring a single file

You can also ignore a single file right from the articles frontmatter, by
adding `jinja: false`.

```markdown
---
jinja: false

---
```

# Escaping

Sometimes you want the ability to have jinja templates in a post, but also the
ability to keep a raw jinja template.  There are a couple of techniques that
are covered mroe in the jinja docs for
[escaping](https://jinja.palletsprojects.com/en/3.1.x/templates/#escaping)

```markdown
{% raw %}
{{ '\\n'.join(markata.map('f"* [{title}]({slug})"', sort='slug')) }}
{% endraw %}

{{ '{{' }} '\\n'.join(markata.map('f"* [{title}]({slug})"', sort='slug')) {{ '}}' }}
```

# Creating a jinja extension

Here is a bit of a boilerplate example of a jinja extension.

``` python
from jinja2 import nodes
from jinja2.ext import Extension

class ExampleExtension(Extension):
    tags = {"example"}

    def __init__(self, environment):
        super().__init__(environment)

    def parse(self, parser):
        line_number = next(parser.stream).lineno
        arg = [parser.parse_expression()]
        return nodes.CallBlock(self.call_method("run", arg), [], [], "").set_lineno(
            line_number
        )

    def run(self, arg, caller):
        return f'hello {arg}'
```

So that markata picks up your extension, you will need to register an
entrypoint named `markata.jinja_md`.  Once installed markata will automatically
load this extension to its list of  jinja extensions.

``` toml
[project.entry-points."markata.jinja_md"]
markta_gh = "example_extension:ExampleExtension"
```

Once you have your extension created and ready to use you can use it in your
markdown.

``` markdown
{% example world %}
```

"""
from pathlib import Path
from typing import List, TYPE_CHECKING

import jinja2
from jinja2 import TemplateSyntaxError, Undefined, UndefinedError, nodes
from jinja2.ext import Extension
import pathspec
import pkg_resources
import pydantic

from markata import __version__
from markata.hookspec import hook_impl, register_attr


def register_jinja_extensions(config: dict) -> List[Extension]:
    """
    Gets jinja extensions from entrypoints and loads them in.

    Returns: List of jinja Extensions
    """

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
    article.content = Template( template, undefined=_SilentUndefined).render()
    ```
    """

    def _fail_with_undefined_error(self, *args, **kwargs):
        return ""


class PostTemplateSyntaxError(TemplateSyntaxError):
    """
    Custom error message for post template syntax errors.
    """


class JinjaMd(pydantic.BaseModel):
    jinja: bool = True


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
    # for article in markata.iter_articles(description="jinja_md"):
    jinja_env = jinja2.Environment(
        extensions=[IncludeRawExtension, *register_jinja_extensions(config)],
    )

    for article in markata.articles:
        if article.get("jinja", True) and not ignore_spec.match_file(article["path"]):
            try:
                key = markata.make_hash(article.content)
                content_from_cache = markata.precache.get(key)
                if content_from_cache is None:
                    article.content = jinja_env.from_string(article.content).render(
                        __version__=__version__,
                        **article,
                        post=article,
                    )
                    with markata.cache:
                        markata.cache.set(key, article.content)
                else:
                    article.content = content_from_cache
                # prevent double rendering
                article.jinja = False
            except TemplateSyntaxError as e:
                errorline = article.content.split("\n")[e.lineno - 1]
                msg = f"""
                Error while processing post {article['path']}

                {errorline}
                """

                raise PostTemplateSyntaxError(msg, lineno=e.lineno)
            except UndefinedError as e:
                raise UndefinedError(f'{e} in {article["path"]}')


class JinjaMdConfig(pydantic.BaseModel):
    ignore: List[str] = []


class Config(pydantic.BaseModel):
    jinja_md: JinjaMdConfig = JinjaMdConfig()


@hook_impl(tryfirst=True)
def config_model(markata: "Markata") -> None:
    markata.config_models.append(Config)
