from rich.console import Console

from markata import LifeCycle, Markata


def test_make_instance():
    markata = Markata()
    assert isinstance(markata, Markata)


def test_make_instance_with_console():
    console = Console()
    markata = Markata(console=console)


def test_markata_lifecycles_callable():
    markata = Markata()
    for method in LifeCycle.__members__.keys():
        # simulating markata.method ex: markata.load()
        assert callable(getattr(markata, method))


def test_markata_config_accessible():
    markata = Markata()
    assert isinstance(markata.config, dict)


def test_markata_hooks_accessible():
    markata = Markata()
    assert isinstance(markata.hooks, list)


def test_markata_has_content_dir_hash():
    markata = Markata()
    assert len(markata.content_dir_hash) == 32


def test_markata_describe():
    markata = Markata()
    "version" in markata.describe()


def test_markata_to_json():
    markata = Markata()
    isinstance(markata.to_json(), str)
