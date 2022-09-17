import pytest

from markata import Markata
from markata.plugins import redirects


@pytest.mark.parametrize(
    "tmp_files, old, new",
    [
        ({"_redirects": "old new"}, "old", "new"),
    ],
    indirect=["tmp_files"],
)
def test_find_text(tmp_files, old, new):
    m = Markata()
    redirects.save(m)
