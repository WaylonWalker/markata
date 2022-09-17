import os
from contextlib import contextmanager
from pathlib import Path

import pytest

from markata import Markata
from markata.plugins import redirects


@contextmanager
def set_directory(path: Path):
    """Sets the cwd within the context

    Args:
        path (Path): The path to the cwd

    Yields:
        None
    """

    origin = Path().absolute()
    try:
        os.chdir(path)
        yield
    finally:
        os.chdir(origin)


@pytest.mark.parametrize(
    "tmp_files, old, new",
    [
        ({"static/_redirects": "old new"}, "old", "new"),
        ({"static/_redirects": "post_one post-one"}, "post_one", "post-one"),
    ],
    indirect=["tmp_files"],
)
def test_redirect_exists(tmp_files, old, new):
    with set_directory(tmp_files):
        m = Markata()
        redirects.save(m)
        redirect_file = Path("markout") / old / "index.html"
        assert redirect_file.exists()
        assert (
            f'<meta http-equiv="Refresh" content="0; url=\'{new}\'" />'
            in redirect_file.read_text()
        )


@pytest.mark.parametrize(
    "tmp_files, old",
    [
        ({"static/_redirects": "posts/* pages/:splat/"}, "posts"),
        ({"static/_redirects": "posts/*/2000 pages/:splat/2000"}, "posts"),
    ],
    indirect=["tmp_files"],
)
def test_redirect_ignore_splat(tmp_files, old):
    with set_directory(tmp_files):
        m = Markata()
        redirects.save(m)
        redirect_file = Path("markout") / old / "index.html"
        assert not redirect_file.exists()


@pytest.mark.parametrize(
    "tmp_files, old",
    [
        ({"static/_redirects": "/blog/* /blog/404.html 404"}, "posts"),
        ({"static/_redirects": "/blog/* /blog/404.html 301"}, "posts"),
        ({"static/_redirects": "/blog/* /blog/404.html 200"}, "posts"),
    ],
    indirect=["tmp_files"],
)
def test_redirect_ignore_more_params(tmp_files, old):
    with set_directory(tmp_files):
        m = Markata()
        redirects.save(m)
        redirect_file = Path("markout") / old / "index.html"
        assert not redirect_file.exists()


@pytest.mark.parametrize(
    "tmp_files, redirect_file, old, new",
    [
        ({"static/_redirects": "old new"}, "static/_redirects", "old", "new"),
        ({"_redirects": "old new"}, "_redirects", "old", "new"),
        ({"assets/_redirects": "old new"}, "assets/_redirects", "old", "new"),
    ],
    indirect=["tmp_files"],
)
def test_redirect_configure_redirect_file(tmp_files, redirect_file, old, new):
    with set_directory(tmp_files):
        m = Markata()
        m.config["redirects"] = redirect_file
        redirects.save(m)
        redirect_file = Path("markout") / old / "index.html"
        assert redirect_file.exists()
        assert (
            f'<meta http-equiv="Refresh" content="0; url=\'{new}\'" />'
            in redirect_file.read_text()
        )


@pytest.mark.parametrize(
    "tmp_files, redirect_template, old, new",
    [
        (
            {
                "static/_redirects": "old new",
                "templates/redirect_template": "{{ original }} is now {{ new }}",
            },
            "templates/redirect_template",
            "old",
            "new",
        ),
        (
            {
                "static/_redirects": "old new",
                "plugins/templates/redirect_template.html": "{{ original }} is now {{ new }}",
            },
            "plugins/templates/redirect_template.html",
            "old",
            "new",
        ),
    ],
    indirect=["tmp_files"],
)
def test_redirect_custom_template(tmp_files, redirect_template, old, new):
    with set_directory(tmp_files):
        m = Markata()
        m.config["redirect_template"] = redirect_template
        redirects.save(m)
        redirect_file = Path("markout") / old / "index.html"
        assert redirect_file.exists()
        assert f"{old} is now {new}" in redirect_file.read_text()


@pytest.mark.parametrize(
    "tmp_files, old, new",
    [
        (
            {
                "static/_redirects": "",
            },
            "old",
            "new",
        ),
    ],
    indirect=["tmp_files"],
)
def test_redirect_empty(tmp_files, old, new):
    with set_directory(tmp_files):
        m = Markata()
        redirects.save(m)


def test_redirect_file_missing(tmpdir):
    with set_directory(tmpdir):
        m = Markata()
        redirects.save(m)
