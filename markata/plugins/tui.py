import subprocess

from rich.console import RenderableType
from textual.app import App
from textual.containers import Container
from textual.widget import Widget
from textual.widgets import Footer, Header

from markata import Markata
from markata.hookspec import hook_impl


class MarkataWidget(Widget):
    def __init__(self, markata: Markata, widget: str = "server"):
        super().__init__(name=widget, id=widget)
        self.m = markata
        self.widget = widget
        self.renderable = getattr(self.m, self.widget)

    def render(self):
        return self.renderable

    async def update(self, renderable: RenderableType) -> None:
        self.renderable = renderable
        self.refresh()


class MarkataApp(App):
    CSS = """
    #plugins, #summary {
            width: 10%;
            }
    #main {
            layout: horizontal;
        }

    #center {
            layout: vertical;
            height: 100%;
            width: 80%;
            }
    #runner, #server {
            height: 50%;
            }
    """

    async def on_load(self, event):
        self.m = Markata()
        self.m.console.quiet = True
        self.bind("q", "quit")
        self.bind("r", "refresh")
        self.new_cmd = self.m.config.get("tui", {}).get("new_cmd", "")
        if self.new_cmd != "":
            self.bind("n", "new")
        self.footer = Footer()
        self.header = Header()

        self.server = MarkataWidget(self.m, "server")
        self.runner = MarkataWidget(self.m, "runner")
        self.plugins = MarkataWidget(self.m, "plugins")
        self.summary = MarkataWidget(self.m, "summary")

    def on_mount(self) -> None:
        self.set_interval(0.2, self.action_refresh)

    def compose(self) -> None:
        yield self.header
        yield self.footer
        yield Container(
            self.plugins,
            Container(
                self.runner,
                self.server,
                id="center",
            ),
            self.summary,
            id="main",
        )

    async def action_refresh(self) -> None:
        self.refresh()
        self.runner.refresh()
        self.server.refresh()
        self.plugins.refresh()
        self.summary.refresh()

    async def action_new(self) -> None:
        subprocess.Popen(self.new_cmd, stdout=subprocess.PIPE, stderr=subprocess.PIPE)


@hook_impl()
def cli(app, markata):
    @app.command()
    def tui():
        MarkataApp().run()


if __name__ == "__main__":
    MarkataApp().run()
