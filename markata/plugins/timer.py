"""default timer plugin"""
from typing import Generator
from markata.hookspec import hook_impl
from bs4 import BeautifulSoup
from markata import Markata, __version__


@hook_impl(hookwrapper=True)
def glob(markata: Markata) -> Generator:
    print("before call")
    breakpoint()
    outcome = yield
    breakpoint()
