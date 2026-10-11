#!/usr/bin/env python3
"""Serve one generated Markata site for controlled local Slow 3G comparisons.

Start two instances on different ports with the same content/assets but different
Markata builds. This is a benchmark server, not a production file server.
"""

from __future__ import annotations

import argparse
import gzip
import mimetypes
import re
from functools import partial
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer
from pathlib import Path
from urllib.parse import unquote, urlsplit

for extension, mime in {
    ".css": "text/css",
    ".js": "text/javascript",
    ".mjs": "text/javascript",
    ".svg": "image/svg+xml",
    ".woff2": "font/woff2",
    ".webp": "image/webp",
    ".avif": "image/avif",
    ".wasm": "application/wasm",
    ".mp4": "video/mp4",
}.items():
    mimetypes.add_type(mime, extension)

TEXT_EXTENSIONS = {".html", ".css", ".js", ".mjs", ".svg", ".json", ".xml", ".txt", ".webmanifest"}
RANGE_PATTERN = re.compile(r"^bytes=(\d*)-(\d*)$")


class Handler(BaseHTTPRequestHandler):
    """Read generated artifacts, never fall back to index.html for missing assets."""

    protocol_version = "HTTP/1.1"

    def __init__(self, *args: object, directory: Path, **kwargs: object) -> None:
        self.directory = directory.resolve()
        super().__init__(*args, **kwargs)

    def do_GET(self) -> None:
        self._serve(head_only=False)

    def do_HEAD(self) -> None:
        self._serve(head_only=True)

    def _resolve(self) -> Path | None:
        try:
            relative_url = unquote(urlsplit(self.path).path, errors="strict")
            if "\\x00" in relative_url or "\\" in relative_url:
                return None
            candidate = (self.directory / relative_url.lstrip("/")).resolve()
            if not candidate.is_relative_to(self.directory):
                return None
            if candidate.is_dir():
                candidate = candidate / "index.html"
            elif not candidate.exists() and not candidate.suffix:
                candidate = candidate / "index.html"
            if candidate.is_file() and candidate.is_relative_to(self.directory):
                return candidate
        except (OSError, UnicodeError, ValueError):
            return None
        return None

    def _serve(self, head_only: bool) -> None:
        path = self._resolve()
        if path is None:
            self.send_error(404, "Asset not found")
            return

        mime, _ = mimetypes.guess_type(str(path))
        mime = mime or "application/octet-stream"
        if mime.startswith("text/"):
            mime += "; charset=utf-8"

        try:
            # Compress text using the same parameters on both sides of the A/B.
            # Media and Range requests always get uncompressed byte ranges.
            wants_gzip = (
                "gzip" in self.headers.get("Accept-Encoding", "").lower()
                and path.suffix.lower() in TEXT_EXTENSIONS
                and not self.headers.get("Range")
            )
            body = path.read_bytes()
            total = len(body)
            range_value = self.headers.get("Range")
            start, end = 0, total - 1
            status = 200
            if range_value:
                match = RANGE_PATTERN.fullmatch(range_value.strip())
                if match is None or not any(match.groups()):
                    self.send_error(416, "Unsupported Range")
                    return
                first, last = match.groups()
                if first:
                    start = int(first)
                    end = min(int(last), total - 1) if last else total - 1
                else:
                    count = int(last)
                    start = max(total - count, 0)
                if start >= total or start > end:
                    self.send_response(416)
                    self.send_header("Content-Range", f"bytes */{total}")
                    self.send_header("Content-Length", "0")
                    self.end_headers()
                    return
                body = body[start : end + 1]
                status = 206
            elif wants_gzip:
                body = gzip.compress(body, compresslevel=6, mtime=0)
        except OSError:
            self.send_error(404, "Asset not found")
            return

        self.send_response(status)
        self.send_header("Content-Type", mime)
        self.send_header("Content-Length", str(len(body)))
        self.send_header("Cache-Control", "no-store")
        self.send_header("Vary", "Accept-Encoding")
        self.send_header("Accept-Ranges", "bytes")
        if status == 206:
            self.send_header("Content-Range", f"bytes {start}-{end}/{total}")
        elif wants_gzip:
            self.send_header("Content-Encoding", "gzip")
        self.end_headers()
        if not head_only:
            try:
                self.wfile.write(body)
            except (BrokenPipeError, ConnectionResetError):
                # A browser aborting a prefetch is normal during test teardown.
                pass

    def log_message(self, _format: str, *_args: object) -> None:
        pass


def serve(directory: Path, host: str = "127.0.0.1", port: int = 8000) -> ThreadingHTTPServer:
    if not directory.is_dir():
        raise ValueError(f"Site output directory not found: {directory}")
    handler = partial(Handler, directory=directory)
    return ThreadingHTTPServer((host, port), handler)


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--root", type=Path, required=True, help="Generated site output directory")
    parser.add_argument("--port", type=int, default=8000)
    parser.add_argument("--host", default="127.0.0.1")
    args = parser.parse_args()
    with serve(args.root, args.host, args.port) as server:
        print(f"Serving {args.root.resolve()} on http://{args.host}:{server.server_port}", flush=True)
        server.serve_forever()
