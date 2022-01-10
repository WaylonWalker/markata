import time
from typing import TYPE_CHECKING

import typer
from rich.layout import Layout

if TYPE_CHECKING:
    from markata import Markata


def make_layout() -> Layout:
    """Define the layout."""
    layout = Layout(name="root")

    layout.split(
        Layout(name="header", size=3),
        Layout(name="main"),
    )
    layout["main"].split_row(
        Layout(name="side", ratio=50),
        Layout(name="mid", ratio=30),
        Layout(name="describe", ratio=20),
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
    def __init__(self, markata: "Markata"):
        self.m = markata
        self.layout = make_layout()

    def __rich__(self):

        return self.layout


def run_until_keyboard_interrupt():
    try:
        while True:
            time.sleep(0.2)
    except KeyboardInterrupt:
        pass


def version_callback(value: bool):
    if value:
        from markata import __version__

        typer.echo(f"Markata CLI Version: {__version__}")
        raise typer.Exit()


def json_callback(value: bool):
    if value:
        from markata import Markata

        typer.echo(Markata().to_json())
        raise typer.Exit()


app = typer.Typer(
    name="Markata",
    help="Awesome CLI user manager.",
)


@app.callback()
def main(
    version: bool = typer.Option(
        None, "--version", callback=version_callback, is_eager=True
    ),
    to_json: bool = typer.Option(
        None, "--to-json", callback=json_callback, is_eager=True
    ),
):
    # Do other global stuff, handle other global options here
    return


def cli():
    from markata import Markata

    m = Markata()
    m._pm.hook.cli(markata=m, app=app)
    app()


if __name__ == "__main__":
    cli()
