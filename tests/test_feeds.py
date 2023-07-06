import rich

import markata
from markata.plugins.feeds import Feed, Feeds


class DummyMarkata:
    def __init__(self):
        self.config = {"feeds": {"all-posts": {"filter": True}}}

    def map(self, *args, **kwargs):
        return []


def test_feed_map(mocker):
    mocked_markata = mocker.patch.object(markata, "Markata", DummyMarkata())
    feed = Feed(name="archive", config={}, posts=[], _m=mocked_markata)
    feed.map()


def test_feeds(mocker):
    mocked_markata = mocker.patch.object(markata, "Markata", DummyMarkata())
    Feeds(markata=mocked_markata)


def test_feeds_iter(mocker):
    mocked_markata = mocker.patch.object(markata, "Markata", DummyMarkata())
    feeds = Feeds(markata=mocked_markata)

    assert [feed for feed in feeds] == ["all-posts"]


def test_feeds_keys(mocker):
    mocked_markata = mocker.patch.object(markata, "Markata", DummyMarkata())
    feeds = Feeds(markata=mocked_markata)

    assert [feed for feed in feeds.keys()] == ["all-posts"]


def test_feeds_items(mocker):
    mocked_markata = mocker.patch.object(markata, "Markata", DummyMarkata())
    feeds = Feeds(markata=mocked_markata)
    iterate_feeds = {name: feed for name, feed in feeds.items()}

    assert list(iterate_feeds.keys()) == ["all-posts"]


def test_feeds___get_item__(mocker):
    mocked_markata = mocker.patch.object(markata, "Markata", DummyMarkata())
    feeds = Feeds(markata=mocked_markata)

    feeds["all-posts"].name == "all_posts"


def test_feeds___rich__(mocker):
    mocked_markata = mocker.patch.object(markata, "Markata", DummyMarkata())
    feeds = Feeds(markata=mocked_markata)

    assert isinstance(feeds.__rich__(), rich.table.Table)


def test_feeds_values(mocker):
    mocked_markata = mocker.patch.object(markata, "Markata", DummyMarkata())
    values = Feeds(markata=mocked_markata).values()
    assert all([isinstance(v, Feed) for v in values])
