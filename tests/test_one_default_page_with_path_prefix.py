import os
import textwrap
from typing import Any

import pytest

from markata import Markata


@pytest.fixture()
def make_index(tmp_path: Any) -> Any:
    pages = tmp_path / "pages"
    pages.mkdir()
    fn = pages / "index.md"
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
            """,
        ),
    )
    return tmp_path


def test_loaded(make_index: Any) -> None:
    os.chdir(make_index)
    m = Markata()
    assert len(m.articles) == 1


@pytest.fixture()
def test_run(make_index: Any) -> Any:
    os.chdir(make_index)
    m = Markata()
    m.config["output_dir"] = "markout/sub-route/"
    m.config["path_prefix"] = "sub-route/"
    m.run()
    return make_index


def test_markout_exists(test_run: Any) -> Any:
    markout = test_run / "markout"
    assert markout.exists()
    sub = test_run / "markout/sub-route"
    assert sub.exists()


def test_index_exists(test_run: Any) -> Any:
    markout = test_run / "markout"
    sub = test_run / "markout/sub-route"
    markout_index = markout / "index.html"
    assert not markout_index.exists()
    sub_index = sub / "index.html"
    assert sub_index.exists()


def test_rss_exists(test_run: Any) -> Any:
    markout = test_run / "markout"
    sub = test_run / "markout/sub-route"
    rss = markout / "rss.xml"
    assert not rss.exists()
    sub_rss = sub / "rss.xml"
    assert sub_rss.exists()
