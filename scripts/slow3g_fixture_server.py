#!/usr/bin/env python3
"""Synthetic gzip HTTP fixture for testing the filmstrip harness, not Markata performance."""

from __future__ import annotations

import argparse
import gzip
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer
from urllib.parse import urlsplit

NAV = "<aside class='feed-sidebar' aria-label='Feed navigation'><nav><h2>Feeds</h2><ul>" + "".join(
    f"<li><a href='/posts/{i}/'>Archived article {i}</a></li>" for i in range(120)
) + "</ul></nav></aside>"
MAIN = (
    "<main id='main-content' tabindex='-1'><h1>Readable main content</h1>"
    "<p>Article text must visibly render while CSS and images are loaded over HTTP.</p>"
    "<img src='/assets/diagram.svg' width='120' height='60' alt='Diagram'></main>"
)
HTML_START = (
    "<!doctype html><html lang='en'><head><meta charset='utf-8'>"
    "<meta name='viewport' content='width=device-width, initial-scale=1'>"
    "<link rel='stylesheet' href='/assets/site.css'><script defer src='/assets/site.js'></script>"
    "</head><body><header>Navigation shell</header><a href='#main-content'>Skip to main</a>"
    "<div class='page-wrapper'>"
)
CSS = (
    "body{font:18px/1.5 system-ui;background:#fff;color:#222;margin:0}"
    "header{padding:16px;background:#25425e;color:#fff}"
    ".page-wrapper{display:flex;padding:12px;gap:16px}"
    ".feed-sidebar{width:200px;max-height:150px;overflow:auto;order:-1;flex-shrink:0}"
    "main{flex:1;min-width:0}h1{font-size:25px}"
    "@media(max-width:650px){.page-wrapper{display:block}.feed-sidebar{display:none}}"
)
SVG = (
    "<svg xmlns='http://www.w3.org/2000/svg' width='120' height='60'>"
    "<rect width='120' height='60' fill='#97bfb7'/></svg>"
)

class Handler(BaseHTTPRequestHandler):
    protocol_version = "HTTP/1.1"

    def do_GET(self) -> None:
        path = urlsplit(self.path).path
        if path == "/before/archive/":
            body, mime = (HTML_START + NAV + MAIN + "</div></body></html>").encode(), "text/html; charset=utf-8"
        elif path == "/after/archive/":
            body, mime = (HTML_START + MAIN + NAV + "</div></body></html>").encode(), "text/html; charset=utf-8"
        elif path == "/bad/archive/":
            body, mime = (HTML_START.replace("/assets/site.css", "/assets/broken.css") + MAIN + NAV + "</div></body></html>").encode(), "text/html; charset=utf-8"
        elif path == "/assets/site.css":
            body, mime = CSS.encode(), "text/css"
        elif path == "/assets/broken.css":
            # Deliberately broken fixture: 200 HTML masquerading as stylesheet.
            body, mime = b"<html><body>Not a stylesheet</body></html>", "text/html"
        elif path == "/assets/site.js":
            body, mime = b"document.documentElement.dataset.fixture='ok';", "text/javascript"
        elif path == "/assets/diagram.svg":
            body, mime = SVG.encode(), "image/svg+xml"
        else:
            self.send_error(404)
            return
        compressed = "gzip" in self.headers.get("Accept-Encoding", "")
        if compressed:
            body = gzip.compress(body, compresslevel=6, mtime=0)
        self.send_response(200)
        self.send_header("Content-Type", mime)
        self.send_header("Content-Length", str(len(body)))
        self.send_header("Cache-Control", "no-store")
        self.send_header("Vary", "Accept-Encoding")
        if compressed:
            self.send_header("Content-Encoding", "gzip")
        self.end_headers()
        self.wfile.write(body)

    def log_message(self, *_args: object) -> None:
        pass


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--host", default="127.0.0.1")
    parser.add_argument("--port", type=int, default=8765)
    args = parser.parse_args()
    ThreadingHTTPServer((args.host, args.port), Handler).serve_forever()
