"""Check network asset-response validation without starting Chromium."""

import json
import tempfile
import unittest
from argparse import Namespace
from pathlib import Path

from slow3g_compare import asset_response_problem, write_report


class AssetResponseTests(unittest.TestCase):
    def test_valid_assets(self):
        for kind, mime in (
            ("Stylesheet", "text/css"),
            ("Stylesheet", "text/css; charset=utf-8"),
            ("Script", "application/javascript"),
            ("Script", "text/javascript"),
            ("Image", "image/svg+xml"),
            ("Image", "image/jpeg"),
            ("Font", "font/woff2"),
            ("Font", "application/octet-stream"),
        ):
            with self.subTest(kind=kind, mime=mime):
                self.assertIsNone(asset_response_problem(kind, mime, 200))

    def test_html_is_not_a_valid_asset(self):
        for kind in ("Stylesheet", "Script", "Image", "Font"):
            with self.subTest(kind=kind):
                self.assertIsNotNone(asset_response_problem(kind, "text/html", 200))

    def test_missing_assets_are_rejected(self):
        self.assertEqual(asset_response_problem("Stylesheet", "text/css", 404), "HTTP 404")

    def test_html_document_is_not_a_subresource(self):
        self.assertIsNone(asset_response_problem("Document", "text/html", 200))


class SummaryTests(unittest.TestCase):
    def test_summary_includes_cls_and_resource_totals(self):
        result = {
            "variant": "before",
            "run": 1,
            "url": "http://example.test/",
            "valid_filmstrip": True,
            "metrics": {"fcp": 100.0, "lcp": 200.0, "cls": 0.012345},
            "encoded_bytes_completed": 2048,
            "requests": [{}, {}],
            "network_failures": [],
            "invalid_asset_responses": [],
            "frames": [],
            "error": None,
        }
        with tempfile.TemporaryDirectory() as temp_dir:
            output = Path(temp_dir)
            write_report([result], output, Namespace())
            summary = json.loads((output / "summary.json").read_text())

        self.assertEqual(summary[0]["median_cls"], 0.012345)
        self.assertEqual(summary[0]["p90_cls"], 0.012345)
        self.assertEqual(summary[0]["median_encoded_bytes_completed"], 2048)
        self.assertEqual(summary[0]["p90_encoded_bytes_completed"], 2048)
        self.assertEqual(summary[0]["median_request_count"], 2)


if __name__ == "__main__":
    unittest.main()
