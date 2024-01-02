from markata import Markata
from markata.plugins.config_model import ConfigFactory
from markata.plugins.post_model import PostFactory


def test_post() -> None:
    config = ConfigFactory().build(hooks=[], disabled_hooks=[])
    post = PostFactory().build(
        markata=Markata(config=config),
        path="pages/my-post.md",
        slug=None,
        title=None,
    )

    assert post.slug == "my-post"
    assert post.title == "My Post"


def test_markata_init() -> None:
    config = ConfigFactory().build(
        hooks=[
            "markata.plugins.create_models",
            "markata.plugins.post_model",
            "markata.plugins.config_model",
        ],
        disabled_hooks=[],
    )

    markata = Markata(config=config)
    markata.run()
