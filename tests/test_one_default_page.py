import os
import textwrap
from typing import Any

import pytest

from markata import Markata


@pytest.fixture(scope="session")
def make_index(tmp_path_factory: Any) -> Any:
    fn = tmp_path_factory.mktemp("pages") / "index.md"
    fn.write_text(
        textwrap.dedent(
            """
            ---
            templateKey: blog-post
            tags: ['python',]
            title:  My Awesome Post
            date: 2022-01-21T16:40:34
            published: False

            ---

            This is my awesome post.
            """
        )
    )
    return tmp_path_factory


def test_loaded(make_index: Any) -> None:
    os.chdir(make_index.getbasetemp())
    m = Markata()
    assert len(m.articles) == 1


@pytest.fixture(scope="session")
def test_run(make_index: Any) -> Any:
    os.chdir(make_index.getbasetemp())
    m = Markata()
    m.run()
    return make_index


def test_markout_exists(test_run: Any) -> Any:
    markout = test_run.getbasetemp() / "markout"
    assert markout.exists()


def test_index_exists(test_run: Any) -> Any:
    markout = test_run.getbasetemp() / "markout"
    index = markout / "index.html"
    assert index.exists()


def test_rss_exists(test_run: Any) -> Any:
    markout = test_run.getbasetemp() / "markout"
    rss = markout / "rss.xml"
    assert rss.exists()
