import subprocess
from typing import List

import pydantic
from rich.console import RenderableType

from markata import Markata
from markata.hookspec import hook_impl, register_attr


class TuiKey(pydantic.BaseModel):
    name: str
    key: str


class TuiConfig(pydantic.BaseModel):
    new_cmd: List[str] = ["markata", "new", "post"]
    keymap: List[TuiKey] = [TuiKey(name="new", key="n")]


class Config(pydantic.BaseModel):
    tui: TuiConfig = TuiConfig()


@hook_impl()
@register_attr("config_models")
def config_model(markata: "Markata") -> None:
    markata.config_models.append(Config)


try:
    from textual.app import App
    from textual.widget import Widget
    from textual.widgets import Footer

    class MarkataWidget(Widget):
        def __init__(self, markata: Markata, widget: str = "server") -> None:
            super().__init__(widget)
            self.m = markata
            self.widget = widget
            self.renderable = getattr(self.m, self.widget)

            def render(self):
                return self.renderable

            async def update(self, renderable: RenderableType) -> None:
                self.renderable = renderable
                self.refresh()

    class MarkataApp(App):
        async def on_load(self, event):
            self.m = Markata()
            self.m.console.quiet = True
            await self.bind("q", "quit", "quit")
            await self.bind("r", "refresh", "refresh")
            self.new_cmd = self.m.config.get("tui", {}).get("new_cmd", "")
            if self.new_cmd != "":
                await self.bind("n", "new", "new")

        async def on_mount(self) -> None:
            self.server = MarkataWidget(self.m, "server")
            self.runner = MarkataWidget(self.m, "runner")
            self.plugins = MarkataWidget(self.m, "plugins")
            self.summary = MarkataWidget(self.m, "summary")
            await self.view.dock(Footer(), edge="bottom")
            await self.view.dock(self.plugins, edge="left", size=30, name="plugins")
            await self.view.dock(self.summary, edge="right", size=30, name="summary")
            await self.view.dock(self.server, self.runner, edge="top")
            self.set_interval(1, self.action_refresh)

        async def action_refresh(self) -> None:
            self.refresh()
            self.runner.refresh()
            self.server.refresh()
            self.plugins.refresh()
            self.summary.refresh()

        async def action_new(self) -> None:
            subprocess.Popen(
                self.new_cmd, stdout=subprocess.PIPE, stderr=subprocess.PIPE
            )

except ImportError:
    ...


@hook_impl()
def cli(app, markata):
    @app.command()
    def tui():
        MarkataApp.run(log="textual.log")


if __name__ == "__main__":
    MarkataApp.run(log="textual.log")
