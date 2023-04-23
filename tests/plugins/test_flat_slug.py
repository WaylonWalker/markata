import pytest

from markata import Markata
from markata.plugins import flat_slug


@pytest.mark.parametrize(
    ("article", "slug"),
    [
        ({"slug": "post_one"}, "post-one"),
        ({"slug": "post-one"}, "post-one"),
        ({"slug": "post one"}, "post-one"),
        ({"slug": "POST_ONE"}, "post-one"),
        ({"slug": "posts/post_one"}, "posts/post-one"),
        ({"slug": "posts/post-one"}, "posts/post-one"),
        ({"slug": "posts/post one"}, "posts/post-one"),
        ({"slug": "posts/POST_ONE"}, "posts/post-one"),
        ({"slug": "POSTS/post_one"}, "posts/post-one"),
        ({"slug": "/posts/post_one"}, "/posts/post-one"),
    ],
)
def test_flat_slug(article, slug) -> None:
    "ensure that the explicit workflow works"
    m = Markata()
    m.articles = [article]
    flat_slug.pre_render(m)

    assert m.articles[0]["slug"] == slug


@pytest.mark.parametrize(
    ("article", "slug"),
    [
        ({"path": "post_one.md"}, "post-one"),
        ({"path": "post-one.md"}, "post-one"),
        ({"path": "post one.md"}, "post-one"),
        ({"path": "POST_ONE.md"}, "post-one"),
        ({"path": "posts/post_one.md"}, "posts/post-one"),
        ({"path": "posts/post-one.md"}, "posts/post-one"),
        ({"path": "posts/post one.md"}, "posts/post-one"),
        ({"path": "posts/POST_ONE.md"}, "posts/post-one"),
        ({"path": "POSTS/post_one.md"}, "posts/post-one"),
        ({"path": "/posts/post_one.md"}, "/posts/post-one"),
    ],
)
def test_flat_slug_path(article, slug) -> None:
    "ensure that the second most explicit workflow works"
    m = Markata()
    m.articles = [article]
    flat_slug.pre_render(m)


@pytest.mark.parametrize(
    ("article", "slug"),
    [
        ({"title": "post_one"}, "post-one"),
        ({"title": "post-one"}, "post-one"),
        ({"title": "post one"}, "post-one"),
        ({"title": "POST_ONE"}, "post-one"),
        ({"title": "/post_one"}, "post-one"),
        ({"title": "posts/post_one"}, "posts/post-one"),
        ({"title": "posts/post-one"}, "posts/post-one"),
        ({"title": "posts/post one"}, "posts/post-one"),
        ({"title": "posts/POST_ONE"}, "posts/post-one"),
        ({"title": "POSTS/post_one"}, "posts/post-one"),
        ({"title": "/posts/post_one"}, "/posts/post-one"),
    ],
)
def test_flat_slug_title(article, slug) -> None:
    "ensure that the second most explicit workflow works"
    m = Markata()
    m.articles = [article]
    flat_slug.pre_render(m)
