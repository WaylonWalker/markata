"""
leading docstring
"""
import ast
import datetime
from functools import lru_cache
from os import path
from pathlib import Path
import textwrap
from typing import List, TYPE_CHECKING

import frontmatter
import jinja2
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
def get_template():
    jinja_env = jinja2.Environment()
    template = jinja_env.from_string(
        (Path(__file__).parent / "default_doc_template.md").read_text(),
    )
    return template


def make_article(markata: "Markata", file: Path, cache) -> frontmatter.Post:
    with open(file) as f:
        raw_source = f.read()
    key = markata.make_hash("docs", "file", raw_source)
    slug = f"{file.parent}/{file.stem}".lstrip("/").lstrip("./")
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

        article = get_template().render(
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
        cache.add(
            key,
            article,
            expire=markata.config.default_cache_expire,
        )

    try:
        article = markata.Post(
            markata=markata,
            path=str(file).replace(".py", ".md"),
            title=file.name,
            content=article,
            ast=ast,
            file=file,
            slug=slug,
            edit_link=edit_link,
            datetime=datetime,
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
    if "articles" not in markata.__dict__:
        markata.articles = []
    for py_file in markata.py_files:
        with markata.cache as cache:
            markata.articles.append(make_article(markata, py_file, cache))
