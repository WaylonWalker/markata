---
date: 2025-12-09
description: Tests the redirects plugin
published: false
slug: tests/plugins/test-redirects
title: test_redirects.py


---

---

Tests the redirects plugin

---

!!! function
    <h2 id="set_directory" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">set_directory <em class="small">function</em></h2>

    context manager to set the directory

???+ source "set_directory <em class='small'>source</em>"
    ```python
    def set_directory(path):
        """
        context manager to set the directory
        """
        cwd = Path.cwd()
        try:
            os.chdir(path)
            yield
        finally:
            os.chdir(cwd)
    ```
!!! function
    <h2 id="test_redirect_exists" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">test_redirect_exists <em class="small">function</em></h2>

    ensure that the default workflow works

???+ source "test_redirect_exists <em class='small'>source</em>"
    ```python
    def test_redirect_exists(tmp_files: Path, old: str, new: str) -> None:
        "ensure that the default workflow works"
        with set_directory(tmp_files):
            m = Markata()
            redirects.save(m)
            redirect_file = Path("markout") / old / "index.html"
            assert redirect_file.exists()
            assert (
                f'<meta http-equiv="Refresh" content="0; url=\'{new}\'" />'
                in redirect_file.read_text()
            )
    ```
!!! function
    <h2 id="test_redirect_ignore_splat" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">test_redirect_ignore_splat <em class="small">function</em></h2>

    splats cannot be supported statically, test that they are ignored

???+ source "test_redirect_ignore_splat <em class='small'>source</em>"
    ```python
    def test_redirect_ignore_splat(tmp_files: Path, old: str) -> None:
        "splats cannot be supported statically, test that they are ignored"
        with set_directory(tmp_files):
            m = Markata()
            redirects.save(m)
            redirect_file = Path("markout") / old / "index.html"
            assert not redirect_file.exists()
    ```
!!! function
    <h2 id="test_redirect_ignore_more_params" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">test_redirect_ignore_more_params <em class="small">function</em></h2>

    status codes cannot be supported statically as they are issued by the server

???+ source "test_redirect_ignore_more_params <em class='small'>source</em>"
    ```python
    def test_redirect_ignore_more_params(tmp_files: Path, old: str) -> None:
        "status codes cannot be supported statically as they are issued by the server"
        with set_directory(tmp_files):
            m = Markata()
            redirects.save(m)
            redirect_file = Path("markout") / old / "index.html"
            assert not redirect_file.exists()
    ```
!!! function
    <h2 id="test_redirect_configure_redirect_file" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">test_redirect_configure_redirect_file <em class="small">function</em></h2>

    ensure that the redirects file can be configured

???+ source "test_redirect_configure_redirect_file <em class='small'>source</em>"
    ```python
    def test_redirect_configure_redirect_file(
        tmp_files: Path, redirect_file: str, old: str, new: str
    ) -> None:
        "ensure that the redirects file can be configured"
        with set_directory(tmp_files):
            m = Markata()
            m.config["redirects"] = redirect_file
            redirects.save(m)
            redirect_html = Path("markout") / old / "index.html"
            assert redirect_html.exists()
            assert (
                f'<meta http-equiv="Refresh" content="0; url=\'{new}\'" />'
                in redirect_html.read_text()
            )
    ```
!!! function
    <h2 id="test_redirect_custom_template" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">test_redirect_custom_template <em class="small">function</em></h2>

    ensure the template can be configured

???+ source "test_redirect_custom_template <em class='small'>source</em>"
    ```python
    def test_redirect_custom_template(
        tmp_files: Path, redirect_template: str, old: str, new: str
    ) -> None:
        "ensure the template can be configured"
        with set_directory(tmp_files):
            m = Markata()
            m.config["redirect_template"] = redirect_template
            redirects.save(m)
            redirect_file = Path("markout") / old / "index.html"
            assert redirect_file.exists()
            assert f"{old} is now {new}" in redirect_file.read_text()
    ```
!!! function
    <h2 id="test_redirect_empty" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">test_redirect_empty <em class="small">function</em></h2>

    ensure empty redirects files work

???+ source "test_redirect_empty <em class='small'>source</em>"
    ```python
    def test_redirect_empty(tmp_files: Path, old: str, new: str) -> None:
        "ensure empty redirects files work"
        with set_directory(tmp_files):
            m = Markata()
            redirects.save(m)
    ```
!!! function
    <h2 id="test_redirect_file_missing" class="admonition-title" style="margin: 0; padding: .5rem 1rem;">test_redirect_file_missing <em class="small">function</em></h2>

    ensure missing redirects file works

???+ source "test_redirect_file_missing <em class='small'>source</em>"
    ```python
    def test_redirect_file_missing(tmpdir: Path) -> None:
        "ensure missing redirects file works"
        with set_directory(tmpdir):
            m = Markata()
            redirects.save(m)
    ```