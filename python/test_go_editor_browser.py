"""Run against a Go workbench with MARKATA_GO_WORKBENCH_URL set."""
import os
import json
import threading
import unittest
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer

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

    def test_new_edits_do_not_wait_for_initial_image_load(self):
        requested = threading.Event()
        release = threading.Event()

        class SlowImage(BaseHTTPRequestHandler):
            def do_GET(self):
                requested.set()
                release.wait(10)
                self.send_response(200)
                self.send_header("Content-Type", "image/svg+xml")
                self.end_headers()
                self.wfile.write(b'<svg xmlns="http://www.w3.org/2000/svg"/>')

            def log_message(self, *_args):
                pass

        server = ThreadingHTTPServer(("127.0.0.1", 0), SlowImage)
        worker = threading.Thread(target=server.serve_forever, daemon=True)
        worker.start()
        try:
            context = self.browser.new_context()
            self.addCleanup(context.close)
            page = context.new_page()

            def slow_initial_document(route):
                response = route.fetch()
                data = response.json()
                if data.get("document"):
                    image = f'<img src="http://127.0.0.1:{server.server_port}/slow.svg">'
                    data["document"] = data["document"].replace("</body>", image + "</body>")
                route.fulfill(response=response, body=json.dumps(data))

            page.route("**/render", slow_initial_document)
            page.goto(os.environ["MARKATA_GO_WORKBENCH_URL"])
            for _ in range(50):
                if requested.is_set():
                    break
                page.wait_for_timeout(100)
            self.assertTrue(requested.is_set(), "initial preview did not request the slow image")
            page.evaluate("document.querySelector('.CodeMirror').CodeMirror.setValue('# Fresh before image')")
            expect(page.frame_locator("#preview").locator("h1")).to_have_text(
                "Fresh before image", timeout=3000)
            self.assertFalse(release.is_set())
        finally:
            release.set()
            server.shutdown()
            worker.join(timeout=2)

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

    def test_unrelated_edit_preserves_rendered_mermaid(self):
        diagram = "```mermaid\ngraph LR\nStart-->End\n```"
        self.edit(f"# Mermaid\n\n{diagram}\n\nBefore")
        expect(self.preview.locator("pre.mermaid svg")).to_be_visible()
        self.frame.evaluate("window.__mermaidSvg = document.querySelector('pre.mermaid svg')")

        self.edit(f"# Mermaid\n\n{diagram}\n\nAfter")
        expect(self.preview.locator(".post-content p").last).to_have_text("After")
        self.assertTrue(self.frame.evaluate(
            "window.__mermaidSvg === document.querySelector('pre.mermaid svg')"))

        self.edit("# Mermaid\n\n```mermaid\ngraph LR\nChanged-->End\n```\n\nAfter")
        expect(self.preview.locator("pre.mermaid svg")).to_contain_text("Changed")
        self.assertFalse(self.frame.evaluate(
            "window.__mermaidSvg === document.querySelector('pre.mermaid svg')"))

    def test_offscreen_mermaid_renders_when_scrolled_into_view(self):
        paragraphs = "\n\n".join(f"Paragraph {number}. " + "Read this text. " * 20
                                  for number in range(30))
        self.edit("# Long page\n\n" + paragraphs +
                  "\n\n```mermaid\ngraph LR\nLazy-->Visible\n```")
        expect(self.preview.locator("pre.mermaid")).to_have_count(1)
        self.page.wait_for_timeout(500)
        expect(self.preview.locator("pre.mermaid svg")).to_have_count(0)

        self.frame.evaluate("window.scrollTo(0, document.body.scrollHeight)")
        expect(self.preview.locator("pre.mermaid svg")).to_contain_text("Lazy")


    def test_publish_rich_diff_marks_prose_and_safely_escapes_text(self):
        result = self.page.evaluate("""() => window.MarkataRichPublishDiff.richHTML(
          '<h1>A title</h1><p>A small page.</p>',
          '<h1>A title</h1><p>A more useful page.</p>'
        )""")
        self.assertIn('mk-rich-insert', result)
        self.assertIn('mk-rich-delete', result)
        self.assertIn('<h1>A title</h1>', result)
        risky = self.page.evaluate("""() => window.MarkataRichPublishDiff.richHTML(
          '<p>Keep the old text</p>',
          '<p>&lt;script&gt;not executable&lt;/script&gt;</p>'
        )""")
        self.assertNotIn('<script>', risky)
        self.assertIn('&lt;script&gt;', risky)

    def test_continuous_typing_still_updates_preview(self):
        self.page.evaluate("""() => {
          window.__typingProbe={done:false,updates:0};
          window.addEventListener('message', event => {
            if(event.data?.type==='preview-updated' && !window.__typingProbe.done)
              window.__typingProbe.updates++;
          });
          const editor=document.querySelector('.CodeMirror').CodeMirror;
          let count=0;
          const interval=setInterval(() => {
            editor.replaceRange(`\\nLive typing ${count}`, {line:editor.lastLine(),ch:editor.getLine(editor.lastLine()).length}, null, '+input');
            count++;
            if(count===40) {clearInterval(interval);window.__typingProbe.done=true;}
          },40);
        }""")
        self.page.wait_for_function('window.__typingProbe.done', timeout=10000)
        self.assertGreater(self.page.evaluate('window.__typingProbe.updates'), 0)
        expect(self.preview.locator('.post-content')).to_contain_text('Live typing 39', timeout=10000)

    def test_preview_discards_obsolete_inflight_frame_render(self):
        self.edit("# Latest edit")
        expect(self.preview.locator("h1")).to_have_text("Latest edit")
        self.page.evaluate("""() => {
          const preview=document.querySelector('#preview').contentWindow;
          // Simulate an older completed render arriving after the editor
          // announced a newer source revision.
          preview.postMessage({type:'preview-invalidate',version:900000},'*');
          preview.postMessage({
            type:'markdown-render',version:899999,
            html:'<h1>Obsolete rendering</h1>',
            source_lines:[0],source_anchors:[],source_offset:0,line_count:1,
            feature_head_html:'',body_end_html:'',attributes:{}
          },'*');
        }""")
        self.page.wait_for_timeout(150)
        expect(self.preview.locator("h1")).to_have_text("Latest edit")
        self.assertEqual(self.frame.evaluate("window.__persistentPreview"), "retained")

    def test_publish_rich_diff_one_edit_in_large_page(self):
        result = self.page.evaluate("""() => {
          const old=Array.from({length:600},(_,i)=>'<p>Paragraph '+i+'</p>');
          const next=[...old];
          next[300]='<p>Paragraph 300 with one small edit</p>';
          return window.MarkataRichPublishDiff.richHTML(old.join(''),next.join(''));
        }""")
        self.assertIn('mk-rich-insert', result)
        self.assertIn('one', result)
        self.assertIn('Paragraph 599', result)
        self.assertNotIn('too complex', result)

    def test_publish_rich_diff_keeps_complex_blocks_as_before_after(self):
        result = self.page.evaluate("""() => window.MarkataRichPublishDiff.richHTML(
          '<pre><code>old()</code></pre>',
          '<pre><code>new()</code></pre>'
        )""")
        self.assertIn('mk-rich-remove', result)
        self.assertIn('mk-rich-add', result)
        self.assertIn('old()', result)
        self.assertIn('new()', result)

    def test_publish_review_shortcut_opens_without_writing(self):
        # No Git checkout in this temporary fixture: provide a recoverable
        # state rather than fabricating a deployed version.
        self.page.keyboard.press('Control+Shift+Enter')
        expect(self.page.locator('#publish-review')).to_be_visible()
        expect(self.page.locator('#publish-confirm')).to_be_disabled()
        self.page.locator('#publish-close').click()
        expect(self.page.locator('#publish-review')).not_to_be_visible()


if __name__ == "__main__":
    unittest.main()
