import os
import textwrap
from typing import Any

import pytest

from markata import Markata


@pytest.fixture(scope="session")
def make_project(tmp_path_factory: Any) -> Any:
    project = tmp_path_factory.mktemp("project")
    module = project / "my_module.py"
    module.write_text(
        textwrap.dedent(
            """
            '''
            Module level docstring
            '''

            def my_func():
                '''
                docstring for my_func
                '''
            class MyClass:
                '''
                docstring for MyClass
                '''

                def my_method(self):
                    '''
                    docstring for my_method
                    '''

            """
        )
    )
    markta_toml = project / "markata.toml"
    markta_toml.write_text(
        textwrap.dedent(
            """
            [markata]
            hooks = [
                "markata.plugins.docs",
                "default",
                ]
            """
        )
    )

    return project


def test_loaded(make_project: Any) -> None:
    os.chdir(make_project)
    m = Markata()
    assert len(m.py_files) == 1


@pytest.fixture(scope="session")
def test_run(make_project: Any) -> Any:
    os.chdir(make_project)
    m = Markata()
    m.run()
    return make_project


def test_markout_exists(test_run: Any) -> Any:
    markout = test_run / "markout"
    assert markout.exists()


def test_index_exists(test_run: Any) -> Any:
    markout = test_run / "markout"
    index = markout / "my_module" / "index.html"
    assert index.exists()


def test_rss_exists(test_run: Any) -> Any:
    markout = test_run / "markout"
    rss = markout / "rss.xml"
    assert rss.exists()
