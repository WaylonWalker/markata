"""Markata Interactive TUI Plugin using Textual."""

from markata.hookspec import hook_impl

import asyncio
from livereload import Server
from rich.text import Text
import threading

from textual.app import App, ComposeResult
from textual.binding import Binding
from textual.containers import ScrollableContainer, Vertical
from textual.widgets import Footer, Header, Static

from markata import Markata
from markata.cli.server import find_port


class Output(ScrollableContainer):
    """A scrollable widget to display command outputs."""

    DEFAULT_CSS = """
    Output {
        height: 1fr;
        border: solid green;
        background: $surface;
        color: $text;
        overflow-y: scroll;
        padding: 0 1;
    }
    """

    def __init__(self):
        super().__init__()
        self._text = Text()

    def write(self, text: str) -> None:
        """Write text to the output."""
        self._text.append(text)
        widget = Static(self._text)
        self.query("Static").remove()
        self.mount(widget)
        self.scroll_end(animate=False)


class MarkataStatus(Static):
    """A widget to display the current status."""

    def __init__(self):
        super().__init__("Ready")
        self.watching = False
        self.serving_url = None
        self.showing_output = False
        self.building = False
        self.update_status()

    def update_status(self):
        """Update the status text."""
        parts = []
        if self.serving_url:
            parts.append(self.serving_url)
        if self.watching:
            parts.append("Watch")
        if self.showing_output:
            parts.append("Output")
        if self.building:
            parts.append("Building...")

        status = " | ".join(parts) if parts else "Ready"
        self.update(Text(status))

    def set_serving(self, serving, url=None):
        """Set server status."""
        self.serving_url = url if serving else None
        self.update_status()

    def set_watching(self, watching):
        """Set watch status."""
        self.watching = watching
        self.update_status()

    def set_output(self, showing):
        """Set output visibility status."""
        self.showing_output = showing
        self.update_status()

    def set_building(self):
        """Set building status."""
        self.building = True
        self.update_status()

    def set_ready(self):
        """Reset to ready state."""
        self.serving_url = None
        self.watching = False
        self.showing_output = False
        self.building = False
        self.update_status()


class MarkataInteractive(
    App,
):
    """A Textual app for interacting with Markata."""

    CSS = """
    Vertical {
        height: 100%;
        width: 100%;
    }

    Header {
        dock: top;
    }

    Output {
        display: none;
        height: 1fr;
        background: $surface-darken-1;
        color: $text;
        overflow-y: scroll;
        padding: 1;
    }

    Output.show-output {
        display: block;
    }

    Output > Static {
        width: 100%;
    }

    MarkataStatus {
        height: auto;
        width: 100%;
        background: $surface;
        color: $text;
        border: solid $primary;
        padding: 0 1;
    }

    Footer {
        dock: bottom;
    }
    """

    BINDINGS = [
        Binding("q", "quit", "Quit", show=True),
        Binding("w", "toggle_watch", "Toggle Watch", show=True),
        Binding("b", "build", "Build", show=True),
        Binding("s", "toggle_serve", "Toggle Serve", show=True),
        Binding("o", "toggle_output", "Toggle Output", show=True),
    ]

    def __init__(self, markata: Markata):
        super().__init__()
        self.markata = markata
        self.markata.config.profiler.should_profile = False
        self.watching = False
        self._watch_task = None
        self.serving = False
        self._server = None
        self._server_thread = None
        self.port = None
        self.content_dir_hash = None
        self.show_output = False

    def compose(self) -> ComposeResult:
        """Create child widgets for the app."""
        self.watching = False
        self._watch_task = None  # Store watch task reference
        yield Header()
        with Vertical():
            yield MarkataStatus()
            yield Output()
        yield Footer()

    async def on_mount(self) -> None:
        """Handle app mount."""
        self.title = "Markata Interactive"
        # Hide output initially
        output = self.query_one(Output)
        output.remove_class("show-output")

        # Start server and watch by default
        await self.action_toggle_serve()
        await self.action_toggle_watch()

    def start_server(self):
        """Start the development server with hot reload."""
        if not self.serving:
            output = self.query_one(Output)
            status = self.query_one(MarkataStatus)

            try:
                # Create livereload server
                self._server = Server()
                # Watch the output directory
                self._server.watch(str(self.markata.config.output_dir))
                # Serve the directory
                self._server_thread = threading.Thread(
                    target=lambda: self._server.serve(
                        root=str(self.markata.config.output_dir),
                        port=8000,
                        host="localhost",
                        restart_delay=0,  # Instant reload
                    ),
                    daemon=True,
                )
                self._server_thread.start()

                self.serving = True
                status.set_serving(True)
                output.write("Server started at http://localhost:8000\n")
            except Exception as e:
                output.write(f"Error starting server: {str(e)}\n")
                self.serving = False
                status.set_serving(False)

    def stop_server(self):
        """Stop the development server."""
        if self.serving:
            output = self.query_one(Output)
            status = self.query_one(MarkataStatus)

            # Just mark as stopped and let the daemon thread clean up
            self.serving = False
            self.port = None
            status.set_serving(False)
            output.write("Server stopped\n")
            self._server = None
            self._server_thread = None

    async def watch_files(self):
        """Watch for file changes and rebuild."""
        status = self.query_one(MarkataStatus)
        output = self.query_one(Output)

        output.write("Watch loop started\n")

        try:
            # Initial build
            status.set_building()
            output.write("\n=== Initial Build ===\n")
            await self.rebuild()
            status.set_ready()
            output.write("\nWatching for changes...\n")

            while self.watching:
                try:
                    await asyncio.sleep(0.5)

                    # Get new hash
                    new_hash = self.markata.content_dir_hash

                    # Compare hashes
                    if self.content_dir_hash != new_hash:
                        self.content_dir_hash = new_hash
                        output.write("\n=== Change Detected ===\n")
                        output.write(f"Previous hash: {self.content_dir_hash}\n")
                        output.write(f"New hash: {new_hash}\n")
                        status.set_building()
                        await self.rebuild()
                        status.set_ready()
                        output.write("\nWatching for changes...\n")

                except Exception as e:
                    output.write(f"\nError in watch loop: {str(e)}\n")
                    # If there's a traceback, show it
                    import traceback

                    output.write(traceback.format_exc())
                    await asyncio.sleep(1)
        except Exception as e:
            output.write(f"\nWatch loop failed: {str(e)}\n")
            import traceback

            output.write(traceback.format_exc())
        finally:
            output.write("\nWatch loop ended\n")
            self.watching = False
            status.set_ready()

    async def rebuild(self):
        """Rebuild the site."""
        output = self.query_one(Output)
        # Disable profiler for interactive mode
        try:
            # Capture stdout during build
            import io
            from contextlib import redirect_stdout

            # Create a string buffer to capture output
            output_buffer = io.StringIO()

            # Redirect stdout to our buffer during the build
            with redirect_stdout(output_buffer):
                await asyncio.to_thread(self.markata.run)

            # Get the captured output
            build_output = output_buffer.getvalue()

            # Write both success message and build output
            output.write("Build started...\n")
            output.write(build_output)
            output.write("Build completed successfully\n")

        except Exception as e:
            output.write(f"Build failed: {str(e)}\n")
            # If there's a traceback, show it
            import traceback

            output.write(traceback.format_exc())

    async def action_toggle_watch(self) -> None:
        """Toggle watch mode."""
        status = self.query_one(MarkataStatus)
        output = self.query_one(Output)

        if not self.watching:
            self.watching = True
            status.set_watching(True)
            output.write("Starting watch mode...\n")
            # Cancel existing task if any
            if self._watch_task and not self._watch_task.done():
                self._watch_task.cancel()
            # Create new task
            self._watch_task = asyncio.create_task(self.watch_files())
        else:
            self.watching = False
            status.set_watching(False)  # Set watch status to False
            # Cancel the watch task
            if self._watch_task and not self._watch_task.done():
                self._watch_task.cancel()
            self._watch_task = None
            output.write("Stopped watching\n")

    async def action_build(self) -> None:
        """Build the site."""
        status = self.query_one(MarkataStatus)
        status.set_building()
        await self.rebuild()
        status.set_ready()

    async def action_toggle_serve(self) -> None:
        """Toggle the development server."""
        if not self.serving:
            output = self.query_one(Output)
            status = self.query_one(MarkataStatus)

            try:
                # Find an available port
                self.port = find_port()
                url = f"http://localhost:{self.port}"

                # Create livereload server
                self._server = Server()
                self._server.watch(str(self.markata.config.output_dir))

                self._server_thread = threading.Thread(
                    target=lambda: self._server.serve(
                        root=str(self.markata.config.output_dir),
                        port=self.port,
                        host="localhost",
                        restart_delay=0,  # Instant reload
                    ),
                    daemon=True,
                )
                self._server_thread.start()

                self.serving = True
                status.set_serving(True, url)
                output.write(f"Server started at {url}\n")
            except Exception as e:
                output.write(f"Failed to start server: {str(e)}\n")
        else:
            self.stop_server()

    async def action_toggle_output(self) -> None:
        """Toggle output visibility."""
        output = self.query_one(Output)
        status = self.query_one(MarkataStatus)
        self.show_output = not self.show_output

        if self.show_output:
            output.add_class("show-output")
            status.set_output(True)
        else:
            output.remove_class("show-output")
            status.set_output(False)


@hook_impl()
def cli(app, markata):
    @app.command()
    def i():
        """Launch interactive TUI mode."""
        # from markata.plugins.interactive import InteractivePlugin

        MarkataInteractive(markata).run()
