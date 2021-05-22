from rich.live import Live

from . import (
    Describe,
    Header,
    MarkataCli,
    Plugins,
    Runner,
    Server,
    make_layout,
    run_until_keyboard_interrupt,
)


def run() -> None:
    layout = make_layout()
    layout["header"].update(Header())

    m = MarkataCli()
    server = Server(directory=m.output_dir)
    layout["server"].update(server)

    layout["plugins"].update(Plugins(m))
    m.count = 0
    m._dirhash = ""
    m.load()
    layout["runner"].update(Runner())
    layout["describe"].update(Describe(m))

    with Live(layout, refresh_per_second=1, screen=True, console=m.console):
        run_until_keyboard_interrupt()


if __name__ == "__main__":
    run()
