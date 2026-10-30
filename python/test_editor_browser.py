"""Browser regressions for the single-file demo using an installed Markata wheel."""
import os
import unittest
from pathlib import Path
from tempfile import TemporaryDirectory

from playwright.sync_api import expect, sync_playwright

from editor_test_support import load_editor, serve_editor

RAIL_SOURCE = """# :lucide-globe: Timeline

<style>
.post-content .md-grid.rail > .md-card { overflow: visible; }
@media (min-width:1000px) { .post-content .md-grid.rail > .md-card::before { top:-6px; } }
</style>

::: md-grid rail {data-columns="3"}

:::: md-card

_2017–2020_

### [First](#first)

First entry.

::::

:::: md-card

_2021–2023_

### [Second](#second)

Second entry.

::::

:::: md-card

_2024–2026_

### [Third](#third)

Third entry.

::::

:::
"""

def chart_source(label):
    return f'''````chartjs
{{
  "type": "bar",
  "data": {{"labels": ["{label}"], "datasets": [{{"label": "Value", "data": [1]}}]}}
}}
````'''


class EditorBrowserTest(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.workspace = TemporaryDirectory()
        cls.addClassCleanup(cls.workspace.cleanup)
        cls.demo = load_editor()
        cls.demo.LSP_WORKSPACE = Path(cls.workspace.name)
        cls.server = serve_editor(cls.demo.app)
        cls.url = cls.server.__enter__()
        cls.addClassCleanup(cls.server.__exit__, None, None, None)
        cls.driver = sync_playwright().start()
        cls.addClassCleanup(cls.driver.stop)
        options = {"headless": True}
        if executable := os.environ.get("MARKATA_TEST_CHROMIUM"):
            options["executable_path"] = executable
        cls.browser = cls.driver.chromium.launch(**options)
        cls.addClassCleanup(cls.browser.close)

    def setUp(self):
        self.context = self.browser.new_context(viewport={"width": 2400, "height": 1000})
        self.context.add_init_script("""try {
          localStorage.setItem('color-mode','dark');
          localStorage.setItem('theme-palette-dark','ayu-dark');
        } catch {}""")
        self.page = self.context.new_page()
        self.errors = []
        self.stylesheet_requests = []
        self.page.on("pageerror", lambda error: self.errors.append(str(error)))
        self.page.on("request", lambda request: self.stylesheet_requests.append(request.url)
                     if request.resource_type == "stylesheet" else None)
        self.page.goto(self.url)
        expect(self.page.locator(".CodeMirror")).to_be_visible()
        self.preview = self.page.frame_locator("#preview")
        expect(self.preview.locator(".post-content h1")).to_be_visible()

    def tearDown(self):
        self.context.close()
        self.assertEqual(self.errors, [], "Browser JavaScript errors")

    def edit(self, markdown):
        self.page.evaluate("value => document.querySelector('.CodeMirror').CodeMirror.setValue(value)", markdown)

    def frame(self):
        return self.page.locator("#preview").element_handle().content_frame()

    def test_full_width_layout_lucide_and_responsive_rail(self):
        self.edit(RAIL_SOURCE)
        expect(self.preview.locator("h1")).to_contain_text("Timeline")
        self.assertEqual(self.page.locator("#workbench").bounding_box()["width"], 2400)
        expect(self.preview.locator('[data-icon="lucide/globe"] svg')).to_have_count(1)
        self.frame().wait_for_function("""() => {
          const cards = [...document.querySelectorAll('.rail > .md-card')];
          return cards.length === 3 && cards.every(e => Math.abs(e.getBoundingClientRect().top - cards[0].getBoundingClientRect().top) < 1)
            && getComputedStyle(cards[0], '::before').top === '-6px';
        }""")
        self.page.set_viewport_size({"width": 1600, "height": 1000})
        self.frame().wait_for_function("""() => {
          const cards = [...document.querySelectorAll('.rail > .md-card')];
          const marker = getComputedStyle(cards[0], '::before');
          return cards[1].getBoundingClientRect().top > cards[0].getBoundingClientRect().bottom
            && getComputedStyle(document.querySelector('.rail'), '::before').width === '2px'
            && marker.content === '\"\"' && marker.left === '-6px' && marker.top === '20px'
            && getComputedStyle(cards[0]).overflow === 'visible';
        }""")
        self.page.set_viewport_size({"width": 390, "height": 844})
        self.assertLessEqual(self.page.evaluate("document.documentElement.scrollWidth"), 390)
        sections = self.page.locator("#workbench > section")
        self.assertGreater(sections.nth(1).bounding_box()["y"], sections.first.bounding_box()["y"])

    def test_keymaps_live_edits_and_saved_draft(self):
        self.edit("# Highlighted\n\n**Bold** and `code`.")
        expect(self.preview.locator("h1")).to_have_text("Highlighted")
        expect(self.page.locator(".cm-header").first).to_be_visible()
        expect(self.page.locator(".cm-strong").first).to_be_visible()
        self.page.locator("#keymap").select_option("vim")
        self.page.keyboard.press("i")
        self.page.keyboard.type("Vim ")
        self.page.keyboard.press("Escape")
        self.assertIn("Vim ", self.page.evaluate("document.querySelector('.CodeMirror').CodeMirror.getValue()"))
        for name in ["emacs", "sublime", "default"]:
            self.page.locator("#keymap").select_option(name)
            self.assertEqual(self.page.evaluate("document.querySelector('.CodeMirror').CodeMirror.getOption('keyMap')"), name)
        self.page.locator("#keymap").select_option("emacs")
        self.edit("# Saved draft")
        expect(self.preview.locator("h1")).to_have_text("Saved draft")
        initial_stylesheets = list(self.stylesheet_requests)
        self.edit("# Second edit")
        expect(self.preview.locator("h1")).to_have_text("Second edit")
        self.assertEqual(self.stylesheet_requests, initial_stylesheets,
                         "re-rendering in the existing frame must not request stylesheets again")
        self.edit("# Saved draft")
        expect(self.preview.locator("h1")).to_have_text("Saved draft")
        self.page.reload()
        expect(self.page.locator("#keymap")).to_have_value("emacs")
        expect(self.preview.locator("h1")).to_have_text("Saved draft")

    def test_initial_dark_document_and_picker_updates(self):
        expect(self.preview.locator("html")).to_have_attribute("data-theme", "dark")
        expect(self.preview.locator("html")).to_have_attribute("data-palette", "ayu-dark")
        self.assertNotEqual(self.page.locator("#preview").evaluate("e=>getComputedStyle(e).backgroundColor"), "rgb(255, 255, 255)")
        self.page.locator(".theme-picker-toggle").click()
        expect(self.page.locator(".theme-card").first).to_be_visible()
        for tab, card in [("style", ".style-card"), ("font", ".font-card")]:
            self.page.locator(f'[data-picker-tab="{tab}"]').click()
            expect(self.page.locator(card).first).to_be_visible()
        self.page.locator('[data-picker-tab="colors"]').click()
        self.page.locator('.theme-card[data-palette="ash"]').click()
        expect(self.preview.locator("html")).to_have_attribute("data-palette", "ash")
        self.edit("# New appearance")
        expect(self.preview.locator("h1")).to_have_text("New appearance")
        expect(self.preview.locator("html")).to_have_attribute("data-palette", "ash")

    def test_feature_styles_are_added_and_removed_with_live_content(self):
        self.edit("# Plain preview")
        expect(self.preview.locator("h1")).to_have_text("Plain preview")
        self.page.reload()
        expect(self.preview.locator("h1")).to_have_text("Plain preview")
        expect(self.preview.locator('head style[data-markata-render-head]')).to_have_count(0)
        baseline_stylesheets = list(self.stylesheet_requests)
        self.edit("# Code styles\n\n```python\nprint('hello')\n```")
        expect(self.preview.locator("pre code")).to_be_visible()
        expect(self.preview.locator('head style[data-markata-render-head]')).to_have_count(2)
        self.assertEqual(self.stylesheet_requests, baseline_stylesheets,
                         "inline feature styles must not reload the document stylesheets")
        self.edit("# Plain again")
        expect(self.preview.locator("h1")).to_have_text("Plain again")
        expect(self.preview.locator('head style[data-markata-render-head]')).to_have_count(0)

    def test_chart_config_updates_without_reloading_the_preview(self):
        self.edit(chart_source("First"))
        self.frame().wait_for_function("""() => {
          const canvas = document.querySelector('canvas[id^="chartjs-"]');
          return canvas && window.Chart?.getChart(canvas)?.data.labels[0] === 'First';
        }""")
        self.edit(chart_source("Second"))
        self.frame().wait_for_function("""() => {
          const canvas = document.querySelector('canvas[id^="chartjs-"]');
          return canvas && window.Chart?.getChart(canvas)?.data.labels[0] === 'Second';
        }""")
        self.assertEqual(self.errors, [], "Chart.js rerenders must not leave stale initializers")

    def test_mermaid_initializer_survives_removing_and_restoring_a_diagram(self):
        self.frame().wait_for_function("() => typeof window.initMermaid === 'function'")
        self.frame().evaluate("""() => {
          window.__mermaidInitCalls = 0;
          const initialize = window.initMermaid;
          window.initMermaid = () => { window.__mermaidInitCalls++; return initialize(); };
        }""")
        self.edit("# Plain text")
        expect(self.preview.locator("h1")).to_have_text("Plain text")
        self.frame().wait_for_function("() => window.__mermaidInitCalls >= 2")
        self.assertEqual(self.frame().evaluate("window.__mermaidInitCalls"), 2)
        self.edit("# Diagram\n\n```mermaid\ngraph TD\n  A --> B\n```")
        expect(self.preview.locator("pre.mermaid")).to_be_visible()
        self.frame().wait_for_function("() => window.__mermaidInitCalls >= 4")
        self.assertEqual(self.frame().evaluate("window.__mermaidInitCalls"), 4)


if __name__ == "__main__":
    unittest.main()
