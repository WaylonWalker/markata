---
date: 2025-12-09
description: "The plugin automatically generates documentation pages from Python source
  code docstrings. It parses Python files, extracts docstrings and code structure,
  and\u2026"
published: false
slug: markata/plugins/docs
title: docs.py


---

---

The `markata.plugins.docs` plugin automatically generates documentation pages from Python
source code docstrings. It parses Python files, extracts docstrings and code structure,
and creates markdown documentation.

## Installation

This plugin is built-in and enabled by default through the 'default' plugin.
If you want to be explicit, you can add it to your list of plugins:

```toml
hooks = [
    "markata.plugins.docs",
]
```

## Uninstallation

Since this plugin is included in the default plugin set, to disable it you must explicitly
add it to the disabled_hooks list if you are using the 'default' plugin:

```toml
disabled_hooks = [
    "markata.plugins.docs",
]
```

## Configuration

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

## Functionality

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

---

!!! function
    <h2 id="glob" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">glob <em class="small">function</em></h2>

    finds k

    ## Parameters

    `markata` the markata object

???+ source "glob <em class='small'>source</em>"
    ```python
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
    ```
!!! function
    <h2 id="load" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">load <em class="small">function</em></h2>

    similar to [glob](../glob)

???+ source "load <em class='small'>source</em>"
    ```python
    def load(markata: "MarkataDocs") -> None:
        """
        similar to [glob](../glob)
        """
        # if "articles" not in markata.__dict__:
        #     markata.articles = []
        for py_file in markata.py_files:
            with markata.cache as cache:
                markata.articles.append(make_article(markata, py_file, cache))
    ```