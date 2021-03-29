import time

from rich.layout import Layout

from markata import Markata


class MarkataCli(Markata):
    count: int = 0
    _dirhash: str = ""


def make_layout() -> Layout:
    """Define the layout."""
    layout = Layout(name="root")

    layout.split(
        Layout(name="header", size=3),
        Layout(name="main"),
    )
    layout["main"].split(
        Layout(name="side", ratio=0.5, minimum_size=30),
        Layout(name="mid"),
        Layout(name="describe", ratio=0.2, minimum_size=30),
        direction="horizontal",
    )
    layout["mid"].split(
        Layout(name="server"),
        Layout(name="runner"),
    )
    layout["side"].split(
        Layout(name="plugins"),
    )
    return layout


class RichM:
    def __init__(self, markata: MarkataCli):
        self.m = markata


def run_until_keyboard_interrupt():
    try:
        while True:
            time.sleep(0.2)
    except KeyboardInterrupt:
        pass
