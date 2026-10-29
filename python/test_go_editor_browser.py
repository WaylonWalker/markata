"""Run against a Go workbench with MARKATA_GO_WORKBENCH_URL set."""
import os
import unittest

from playwright.sync_api import expect, sync_playwright


@unittest.skipUnless(os.environ.get("MARKATA_GO_WORKBENCH_URL"), "Set MARKATA_GO_WORKBENCH_URL to a running Go workbench")
class GoEditorBrowserTest(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.driver = sync_playwright().start()
        cls.addClassCleanup(cls.driver.stop)
        options = {"headless": True}
        if executable := os.environ.get("MARKATA_TEST_CHROMIUM"):
            options["executable_path"] = executable
        cls.browser = cls.driver.chromium.launch(**options)
        cls.addClassCleanup(cls.browser.close)

    def setUp(self):
        self.context = self.browser.new_context()
        self.addCleanup(self.context.close)
        self.page = self.context.new_page()
        self.errors, self.stylesheets, self.renders = [], [], []
        self.page.on("pageerror", lambda error: self.errors.append(str(error)))
        self.page.on("request", lambda req: self.stylesheets.append(req.url)
                     if req.resource_type == "stylesheet" else None)
        self.page.on("response", lambda response: self.renders.append(response.json())
                     if response.url.endswith("/render") and response.ok else None)
        self.page.goto(os.environ["MARKATA_GO_WORKBENCH_URL"])
        self.preview = self.page.frame_locator("#preview")
        expect(self.preview.locator("h1")).to_be_visible()
        self.frame = self.page.locator("#preview").element_handle().content_frame()
        self.frame.evaluate("window.__persistentPreview = 'retained'")

    def tearDown(self):
        self.assertEqual(self.errors, [])

    def edit(self, markdown):
        self.page.evaluate("s => document.querySelector('.CodeMirror').CodeMirror.setValue(s)", markdown)

    def test_fragments_retain_frame_and_styles(self):
        self.edit("# Plain")
        expect(self.preview.locator("h1")).to_have_text("Plain")
        baseline = list(self.stylesheets)
        for i in range(3):
            self.edit(f"# Edit {i}")
            expect(self.preview.locator("h1")).to_have_text(f"Edit {i}")
            self.assertEqual(self.frame.evaluate("window.__persistentPreview"), "retained")
            self.assertIsNone(self.renders[-1]["document"])
            self.assertEqual(self.renders[-1]["document_bytes"], 0)
            self.assertNotIn("palette.css", self.renders[-1]["feature_head_html"])
        self.assertEqual(self.stylesheets, baseline)
        self.assertIsNotNone(self.renders[0]["document"])
        self.assertEqual(self.frame.evaluate("""() => {
          const links=[...document.querySelectorAll('link[rel=stylesheet]')].map(n=>n.href);
          return links.length===new Set(links).size;
        }"""), True)

    def test_conditional_styles_are_added_reused_and_removed(self):
        self.edit("# Code\n\n```go\nfmt.Println(1)\n```")
        expect(self.preview.locator("h1")).to_have_text("Code")
        expect(self.preview.locator('head [data-markata-render-head]')).to_have_count(2)
        baseline = list(self.stylesheets)
        self.frame.evaluate("window.__codeStyle = document.querySelector('link[href$=\"/css/code.css\"]')")
        self.edit("# Code again\n\n```go\nfmt.Println(2)\n```")
        expect(self.preview.locator("h1")).to_have_text("Code again")
        self.assertEqual(self.stylesheets, baseline)
        self.assertTrue(self.frame.evaluate("window.__codeStyle === document.querySelector('link[href$=\"/css/code.css\"]')"))
        self.edit("# Plain again")
        expect(self.preview.locator("h1")).to_have_text("Plain again")
        expect(self.preview.locator('head [data-markata-render-head]')).to_have_count(0)

    def test_tables_remain_sortable_without_reloading_scripts(self):
        for label in ["First", "Second", "Third"]:
            self.edit(f"# {label}\n\n| Name | Value |\n|---|---|\n| B | 2 |\n| A | 1 |")
            expect(self.preview.locator("h1")).to_have_text(label)
            expect(self.preview.locator(".table-sort-button")).to_have_count(2)
            self.preview.get_by_role("button", name="Sort by Name").click()
            expect(self.preview.locator("tbody tr").first).to_contain_text("A")
            expect(self.preview.locator('script[src$="/js/sortable-tables.js"]')).to_have_count(1)

    def test_charts_and_mermaid_update_without_reloading(self):
        for label in ["First", "Second"]:
            self.edit(f'''# {label}

```chartjs
{{"type":"bar","data":{{"labels":["{label}"],"datasets":[{{"data":[1]}}]}}}}
```
''')
            self.frame.wait_for_function("""label => {
              const canvas=document.querySelector('canvas');
              return canvas && window.Chart?.getChart(canvas)?.data.labels[0]===label;
            }""", arg=label)
            self.assertEqual(self.frame.evaluate("Object.keys(window.Chart.instances).length"), 1)
            self.assertTrue(self.frame.evaluate("""() => {
              const sources=[...document.scripts].filter(n=>n.src).map(n=>n.src);
              return sources.length===new Set(sources).size;
            }"""))
        self.edit("# Plain")
        expect(self.preview.locator("h1")).to_have_text("Plain")
        self.assertEqual(self.frame.evaluate("Object.keys(window.Chart.instances).length"), 0)
        for label in ["First", "Second"]:
            self.edit(f"# {label}\n\n```mermaid\ngraph LR\n{label}-->End\n```")
            expect(self.preview.locator("h1")).to_have_text(label)
            expect(self.preview.locator(".mermaid svg")).to_be_visible()
            expect(self.preview.locator(".mermaid svg")).to_contain_text(label)
        self.assertEqual(self.frame.evaluate("window.__persistentPreview"), "retained")


if __name__ == "__main__":
    unittest.main()
