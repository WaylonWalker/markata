#!/usr/bin/env python3
"""Contract checks for the separate docs-site homepage using native Markata tabs."""
from __future__ import annotations

import unittest
from html.parser import HTMLParser
from pathlib import Path

ROOT = Path(__file__).resolve().parent.parent
SITE = ROOT / "docs-site" / "_site"
LABELS = ("curl", "Bash", "mise", "Go", "GitHub", "PyPI")


class NativeTabsParser(HTMLParser):
    def __init__(self) -> None:
        super().__init__()
        self.roots: list[dict[str, str | None]] = []
        self.navs: list[dict[str, str | None]] = []
        self.tabs: list[dict[str, str | None]] = []
        self.panels: list[dict[str, str | None]] = []
        self.fallback_headings: list[str] = []
        self.scripts: list[str] = []
        self._heading = False

    @staticmethod
    def attrs(values: list[tuple[str, str | None]]) -> dict[str, str | None]:
        return dict(values)

    def handle_starttag(self, tag: str, values: list[tuple[str, str | None]]) -> None:
        attrs = self.attrs(values)
        if "data-md-tabs" in attrs:
            self.roots.append(attrs)
        if attrs.get("role") == "tablist":
            self.navs.append(attrs)
        if attrs.get("role") == "tab":
            self.tabs.append(attrs)
        if attrs.get("role") == "tabpanel":
            self.panels.append(attrs)
        if tag == "script" and attrs.get("src"):
            self.scripts.append(str(attrs["src"]))
        if tag == "h3" and "md-tab-fallback-label" in (attrs.get("class") or "").split():
            self._heading = True

    def handle_data(self, data: str) -> None:
        if self._heading:
            self.fallback_headings.append(data)

    def handle_endtag(self, tag: str) -> None:
        if tag == "h3":
            self._heading = False


class DocsHomeNativeTabsTest(unittest.TestCase):
    @classmethod
    def setUpClass(cls) -> None:
        html = (SITE / "index.html").read_text(encoding="utf-8")
        cls.parser = NativeTabsParser()
        cls.parser.feed(html)
        cls.source = html

    def test_all_six_tabs_use_native_accessible_contract(self) -> None:
        p = self.parser
        self.assertEqual(len(p.roots), 1)
        self.assertEqual(p.roots[0]["data-style"], "boxed")
        self.assertIn("md-tabs", (p.roots[0].get("class") or "").split())
        self.assertEqual(len(p.navs), 1)
        self.assertEqual(p.navs[0]["aria-label"], "Installation methods")
        self.assertIn("hidden", p.navs[0])
        self.assertEqual(len(p.tabs), len(LABELS))
        self.assertEqual(len(p.panels), len(LABELS))
        self.assertEqual(p.fallback_headings, list(LABELS))
        for i, (tab, panel) in enumerate(zip(p.tabs, p.panels)):
            ident = f"md-tab-1-{i}"
            self.assertEqual(tab["id"], ident + "-tab")
            self.assertEqual(tab["aria-controls"], ident + "-panel")
            self.assertEqual(panel["id"], ident + "-panel")
            self.assertEqual(panel["aria-labelledby"], ident + "-tab")
            self.assertEqual(tab["aria-selected"], "true" if i == 0 else "false")
            self.assertEqual(tab["tabindex"], "0" if i == 0 else "-1")
            self.assertNotIn("hidden", panel)

    def test_exact_shared_script_is_served(self) -> None:
        self.assertIn("/assets/md-tabs.js", self.parser.scripts)
        shared = ROOT / "pkg" / "themes" / "default" / "static" / "js" / "md-tabs.js"
        self.assertEqual((SITE / "assets" / "md-tabs.js").read_bytes(), shared.read_bytes())
        self.assertNotIn("/assets/tabs.js", self.parser.scripts)
        self.assertFalse((SITE / "assets" / "tabs.js").exists())

    def test_homepage_remains_static_and_has_no_old_tab_hooks(self) -> None:
        self.assertNotIn("data-install-tabs", self.source)
        self.assertNotIn("data-panel=", self.source)
        self.assertNotIn("data-tab=", self.source)
        self.assertIn("installer md-tabs", self.source)
        self.assertIn("md-tab-fallback-label", self.source)


if __name__ == "__main__":
    unittest.main()
