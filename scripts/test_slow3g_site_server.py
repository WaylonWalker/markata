"""Integration tests for the generated-site HTTP A/B benchmark server."""

from __future__ import annotations

import gzip
import http.client
import tempfile
import threading
import unittest
from pathlib import Path

from slow3g_site_server import serve


class SiteServerTests(unittest.TestCase):
    def setUp(self) -> None:
        self.tmp = tempfile.TemporaryDirectory()
        self.addCleanup(self.tmp.cleanup)
        self.root = Path(self.tmp.name) / "site"
        (self.root / "archive").mkdir(parents=True)
        (self.root / "assets").mkdir()
        (self.root / "index.html").write_text("<h1>Home</h1>", encoding="utf-8")
        (self.root / "archive" / "index.html").write_text("<h1>Archive</h1>", encoding="utf-8")
        (self.root / "assets" / "site.css").write_text("body{color:blue}", encoding="utf-8")
        (self.root / "assets" / "app.js").write_text("window.ok=true;", encoding="utf-8")
        (self.root / "assets" / "image.webp").write_bytes(b"RIFFwebp")
        (self.root / "assets" / "demo.mp4").write_bytes(b"0123456789")
        (self.root / "assets" / "font.woff2").write_bytes(b"wOF2")

        self.server = serve(self.root, port=0)
        self.thread = threading.Thread(target=self.server.serve_forever, daemon=True)
        self.thread.start()
        self.addCleanup(self.thread.join, 5)
        self.addCleanup(self.server.server_close)
        self.addCleanup(self.server.shutdown)

    def request(self, path: str, *, method: str = "GET", headers: dict[str, str] | None = None):
        conn = http.client.HTTPConnection("127.0.0.1", self.server.server_port, timeout=5)
        try:
            conn.request(method, path, headers=headers or {})
            response = conn.getresponse()
            status, hdrs, body = response.status, dict(response.getheaders()), response.read()
            return status, hdrs, body
        finally:
            conn.close()

    def test_serves_generated_routes_and_real_assets(self) -> None:
        for url, mime, fragment in (
            ("/", "text/html", b"Home"),
            ("/archive/", "text/html", b"Archive"),
            ("/assets/site.css", "text/css", b"color"),
            ("/assets/app.js", "text/javascript", b"window.ok"),
            ("/assets/image.webp", "image/webp", b"RIFF"),
            ("/assets/font.woff2", "font/woff2", b"wOF2"),
        ):
            with self.subTest(url=url):
                status, headers, body = self.request(url)
                self.assertEqual(status, 200)
                self.assertTrue(headers["Content-Type"].startswith(mime), headers)
                self.assertIn(fragment, body)
                self.assertEqual(headers["Content-Length"], str(len(body)))
        status, _, body = self.request("/assets/missing.css")
        self.assertEqual(status, 404)
        self.assertNotIn(b"Home", body)

    def test_gzip_headers_and_head_match_get(self) -> None:
        url = "/archive/"
        status, headers, body = self.request(url, headers={"Accept-Encoding": "gzip"})
        self.assertEqual(status, 200)
        self.assertEqual(headers["Content-Encoding"], "gzip")
        self.assertEqual(headers["Vary"], "Accept-Encoding")
        self.assertEqual(headers["Content-Length"], str(len(body)))
        self.assertIn(b"Archive", gzip.decompress(body))
        head_status, head_headers, head_body = self.request(
            url, method="HEAD", headers={"Accept-Encoding": "gzip"}
        )
        self.assertEqual(head_status, 200)
        self.assertEqual(head_body, b"")
        self.assertEqual(head_headers["Content-Length"], headers["Content-Length"])

    def test_streams_range_responses_for_media(self) -> None:
        status, headers, body = self.request(
            "/assets/demo.mp4", headers={"Range": "bytes=2-5", "Accept-Encoding": "gzip"}
        )
        self.assertEqual(status, 206)
        self.assertEqual(body, b"2345")
        self.assertEqual(headers["Content-Range"], "bytes 2-5/10")
        self.assertNotIn("Content-Encoding", headers)
        status, _, body = self.request("/assets/demo.mp4", headers={"Range": "bytes=-3"})
        self.assertEqual(status, 206)
        self.assertEqual(body, b"789")
        status, headers, body = self.request("/assets/demo.mp4", headers={"Range": "bytes=100-"})
        self.assertEqual(status, 416)
        self.assertEqual(headers["Content-Range"], "bytes */10")
        self.assertEqual(body, b"")

    def test_rejects_traversal_and_symlinks(self) -> None:
        outside = Path(self.tmp.name) / "secret.txt"
        outside.write_text("secret", encoding="utf-8")
        (self.root / "assets" / "secret.txt").symlink_to(outside)
        for path in ("/%2e%2e/secret.txt", "/assets/secret.txt", "/assets/%5csecret.txt"):
            with self.subTest(path=path):
                status, _, body = self.request(path)
                self.assertNotEqual(status, 200)
                self.assertNotIn(b"secret", body)


if __name__ == "__main__":
    unittest.main()
