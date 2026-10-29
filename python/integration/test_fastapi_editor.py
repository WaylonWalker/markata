"""Integration checks against the installed Markata wheel, not protocol mocks."""
from pathlib import Path
import sys
import unittest
from unittest.mock import patch

from fastapi.testclient import TestClient

sys.path.insert(0, str(Path(__file__).resolve().parents[1]))
from editor_test_support import load_editor


class FastAPIEditorTest(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.demo = load_editor()
        cls.client = TestClient(cls.demo.app)
        cls.client.__enter__()

    @classmethod
    def tearDownClass(cls):
        cls.client.__exit__(None, None, None)

    def test_home_serves_controls_and_editor_without_site_layout_css(self):
        response = self.client.get("/")
        self.assertEqual(response.status_code, 200)
        self.assertIn("CodeMirror.fromTextArea", response.text)
        self.assertIn('id="keymap"', response.text)
        self.assertIn('data-picker-tab="font"', response.text)
        self.assertNotIn('href="/appearance/css/main.css"', response.text)

    def test_real_render_preserves_site_content_and_initial_appearance(self):
        markdown = '# :lucide-globe: Hello\n\n```python\nprint("héllo")\n```'
        response = self.client.post("/render", json={"markdown": markdown, "appearance": {
            "data-theme": "dark", "data-palette": "ayu-dark", "data-fontpack": "system",
        }})
        self.assertEqual(response.status_code, 200)
        data = response.json()
        document = data["document"]
        self.assertIn('class="post h-entry"', document)
        self.assertIn('data-icon="lucide/globe"', document)
        self.assertIn('data-palette="ayu-dark"', document)
        self.assertLess(document.index("color-scheme:dark"), document.index('<link'))
        self.assertIn('class="chroma"', document)
        self.assertIn('Syntax highlighting', data["feature_head_html"])
        self.assertEqual(data["input_bytes"], len(markdown.encode()))
        self.assertEqual(data["document_bytes"], len(document.encode()))
        self.assertGreaterEqual(data["render_ms"], 0)
        self.assertIn("markata;dur=", response.headers["server-timing"])

    def test_fragment_renders_omit_the_document_shell(self):
        response = self.client.post("/render", json={"markdown": "# Small edit", "include_document": False})
        self.assertEqual(response.status_code, 200)
        data = response.json()
        self.assertIsNone(data["document"])
        self.assertEqual(data["document_bytes"], 0)
        self.assertGreater(data["html_bytes"], 0)

    def test_asset_routes_use_exact_lookup_and_correct_media_types(self):
        for name, data in self.demo.presentation.assets.items():
            response = self.client.get("/appearance/" + name)
            self.assertEqual(response.status_code, 200, name)
            self.assertEqual(response.content, data, name)
            expected = "text/css" if name.endswith(".css") else "font/woff2" if name.endswith(".woff2") else "text/javascript"
            self.assertTrue(response.headers["content-type"].startswith(expected), name)
        self.assertEqual(self.client.get("/appearance/codemirror/lib/codemirror.js").status_code, 200)
        for name in ["missing.css", "not-a-font.woff2", "codemirror/missing.js"]:
            self.assertEqual(self.client.get("/appearance/" + name).status_code, 404)

    def test_invalid_requests_and_renderer_failures_are_reported(self):
        for payload in [{"markdown": "x" * 200_001}, {"markdown": "x", "appearance": {"data-palette": "missing"}}]:
            self.assertEqual(self.client.post("/render", json=payload).status_code, 422)
        with patch.object(self.demo, "render_async", side_effect=self.demo.RenderError("Renderer unavailable")):
            response = self.client.post("/render", json={"markdown": "# Hello"})
        self.assertEqual(response.status_code, 422)
        self.assertEqual(response.json()["error"], "Renderer unavailable")


if __name__ == "__main__":
    unittest.main()
