"""Load and serve the demo using the installed Markata package."""
from contextlib import contextmanager
import importlib.util
from pathlib import Path
import socket
import sys
import threading
import time

import uvicorn

ROOT = Path(__file__).resolve().parents[1]


def load_editor():
    spec = importlib.util.spec_from_file_location("markata_fastapi_demo", ROOT / "examples/fastapi_editor.py")
    module = importlib.util.module_from_spec(spec)
    sys.modules[spec.name] = module
    spec.loader.exec_module(module)
    return module


@contextmanager
def serve_editor(app):
    """Keep the OS-assigned port reserved until Uvicorn closes its socket."""
    listener = socket.socket()
    listener.bind(("127.0.0.1", 0))
    listener.listen()
    port = listener.getsockname()[1]
    server = uvicorn.Server(uvicorn.Config(app, log_level="warning"))
    thread = threading.Thread(target=server.run, kwargs={"sockets": [listener]}, daemon=True)
    thread.start()
    try:
        deadline = time.monotonic() + 10
        while not server.started:
            if not thread.is_alive() or time.monotonic() >= deadline:
                raise RuntimeError("Demo test server did not start")
            time.sleep(0.01)
        yield f"http://127.0.0.1:{port}"
    finally:
        server.should_exit = True
        thread.join(timeout=10)
        listener.close()
        if thread.is_alive():
            raise RuntimeError("Demo test server did not stop")
