"""
The `markata.plugins.docs` plugin automatically generates documentation pages from Python
source code docstrings. It parses Python files, extracts docstrings and code structure,
and creates markdown documentation.

# Installation

This plugin is built-in and enabled by default through the 'default' plugin.
If you want to be explicit, you can add it to your list of plugins:

```toml
hooks = [
    "markata.plugins.docs",
]
```

# Uninstallation

Since this plugin is included in the default plugin set, to disable it you must explicitly
add it to the disabled_hooks list if you are using the 'default' plugin:

```toml
disabled_hooks = [
    "markata.plugins.docs",
]
```

# Configuration

Configure documentation generation in your `markata.toml`:

```toml
[markata.docs]
# Directories containing Python files to document
content_directories = [
    "markata",
    "tests"
]

# Template for generated documentation
template = "docs/template.md.j2"
```

## Template Variables

The following variables are available in documentation templates:
- `module`: Module name
- `docstring`: Module docstring
- `classes`: List of class definitions and their docstrings
- `functions`: List of function definitions and their docstrings
- `source`: Path to source file
- `ast`: Abstract Syntax Tree of the module

# Functionality

## Documentation Generation

The plugin:
1. Finds Python files in configured directories
2. Parses files using Python's AST
3. Extracts docstrings and code structure
4. Applies templates to generate markdown
5. Creates documentation pages in the output directory

## AST Analysis

The plugin analyzes:
- Module-level docstrings
- Class definitions and docstrings
- Function definitions and docstrings
- Function parameters and return types
- Code structure and relationships

## Registered Attributes

The plugin adds these attributes to Markata:
- `py_files`: List of Python files being documented
- `content_directories`: List of directories being processed

## Dependencies

This plugin depends on:
- ast (Python standard library) for code parsing
- jinja2 for template rendering
"""

import ast
import datetime
from functools import lru_cache
from os import path
from pathlib import Path
import textwrap
from typing import List, TYPE_CHECKING

import frontmatter
import pydantic

from markata.hookspec import hook_impl, register_attr

if TYPE_CHECKING:
    from markata import Markata

    class MarkataDocs(Markata):
        py_files: List = []
        content_directories: List = []


def add_parents(tree: ast.AST) -> None:
    for node in ast.walk(tree):
        for child in ast.iter_child_nodes(node):
            child.parent = node
            if not hasattr(child, "parents"):
                child.parents = [node]
            child.parents.append(node)
            if isinstance(node, ast.ClassDef) and isinstance(child, ast.FunctionDef):
                child.type = "method"
            elif isinstance(child, ast.FunctionDef):
                child.type = "function"
            elif isinstance(child, ast.ClassDef):
                child.type = "class"


@hook_impl
@register_attr("content_directories", "py_files")
def glob(markata: "MarkataDocs") -> None:
    """
    finds k

    ## Parameters

    `markata` the markata object

    """

    import glob

    markata.py_files = [Path(f) for f in glob.glob("**/*.py", recursive=True)]

    content_directories = list({f.parent for f in markata.py_files})
    if "content_directories" in markata.__dict__:
        markata.content_directories.extend(content_directories)
    else:
        markata.content_directories = content_directories

    try:
        ignore = True
    except KeyError:
        ignore = True

    if ignore and (Path(".gitignore").exists() or Path(".markataignore").exists()):
        import pathspec

        lines = []

        if Path(".gitignore").exists():
            lines.extend(Path(".gitignore").read_text().splitlines())

        if Path(".markataignore").exists():
            lines.extend(Path(".markataignore").read_text().splitlines())

    spec = pathspec.PathSpec.from_lines("gitwildmatch", lines)

    markata.py_files = [
        file for file in markata.py_files if not spec.match_file(str(file))
    ]


@lru_cache
def get_template(markata: "Markata"):
    template = markata.jinja_env.from_string(
        (Path(__file__).parent / "default_doc_template.md").read_text(),
    )
    return template


def make_article(markata: "Markata", file: Path, cache) -> frontmatter.Post:
    from slugify import slugify

    with open(file) as f:
        raw_source = f.read()
    key = markata.make_hash("docs", "file", raw_source)
    slug = f"{file.parent}/{slugify(file.stem)}".lstrip("/").lstrip("./")
    output_html = markata.config.output_dir / slug / "index.html"
    edit_link = (
        str(markata.config.get("repo_url", "https://github.com/"))
        + "edit/"
        + str(markata.config.get("repo_branch", "main"))
        + "/"
        + str(file)
    )
    article_from_cache = markata.precache.get(key)
    if article_from_cache is not None:
        article = article_from_cache
    else:
        tree = ast.parse(raw_source)
        add_parents(tree)
        nodes = [
            n for n in ast.walk(tree) if isinstance(n, (ast.FunctionDef, ast.ClassDef))
        ]

        article = get_template(markata).render(
            ast=ast,
            file=file,
            slug=slug,
            edit_link=edit_link,
            tree=tree,
            datetime=datetime,
            nodes=nodes,
            raw_source=raw_source,
            indent=textwrap.indent,
        )
        cache.set(
            key,
            article,
            expire=markata.config.default_cache_expire,
        )

    try:
        article = markata.Post(
            markata=markata,
            path=str(file).replace(".py", ".md"),
            output_html=output_html,
            title=file.name,
            content=article,
            file=file,
            slug=slug,
            edit_link=edit_link,
            published=False,  # Don't include docs in sitemap by default
            date=datetime.datetime.now(),
            # enable sidebar in the future
            # sidebar="plugins",
        )

    except pydantic.ValidationError as e:
        from markata.plugins.load import ValidationError, get_models

        models = get_models(markata=markata, error=e)
        models = list(models.values())
        models = "\n".join(models)
        raise ValidationError(f"{e}\n\n{models}\nfailed to load {path}") from e

    return article


@hook_impl
@register_attr("articles")
def load(markata: "MarkataDocs") -> None:
    """
    similar to [glob](../glob)
    """
    # if "articles" not in markata.__dict__:
    #     markata.articles = []
    for py_file in markata.py_files:
        with markata.cache as cache:
            markata.articles.append(make_article(markata, py_file, cache))
