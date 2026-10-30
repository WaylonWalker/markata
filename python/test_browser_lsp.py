"""Exercise the shared LSP adapter against both actual example servers.

Build MARKATA_LSP_BINARY and MARKATA_WORKBENCH_BINARY from this branch first.
The Python environment must contain a standalone-renderer Markata wheel.
"""
from contextlib import contextmanager
from pathlib import Path
import os
import json
import subprocess
import tempfile
import unittest

from playwright.sync_api import expect, sync_playwright

from editor_test_support import load_editor, serve_editor


@contextmanager
def go_editor(root):
    process = subprocess.Popen([os.environ["MARKATA_WORKBENCH_BINARY"], "--listen", "127.0.0.1:0",
                                "--workspace", str(root)], stdout=subprocess.DEVNULL, stderr=subprocess.PIPE,
                               text=True)
    try:
        # The command prints its actual OS-assigned address after initialization.
        while line := process.stderr.readline():
            if line.startswith("Markata Go workbench: "):
                yield line.split(" ", 3)[3].strip()
                return
        raise RuntimeError("Go workbench exited before startup")
    finally:
        process.terminate()
        process.communicate(timeout=10)


@unittest.skipUnless(os.environ.get("MARKATA_LSP_BINARY") and os.environ.get("MARKATA_WORKBENCH_BINARY"),
                     "Set both example binary paths to run the browser integration")
class BrowserLSPTest(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.folder = tempfile.TemporaryDirectory(prefix="markata-browser-lsp-fixture-")
        cls.addClassCleanup(cls.folder.cleanup)
        cls.root = Path(cls.folder.name)
        (cls.root / "hello world.md").write_text(
            "---\ntitle: Hello World\nslug: hello\naliases: [greeting]\n---\nHello", encoding="utf-8")
        (cls.root / "markata-go.toml").write_text("[markata-go.glob]\npatterns=['*.md']\n", encoding="utf-8")
        (cls.root / "draft.md").write_text("---\ntitle: Workspace draft\n---\n# Real workspace draft", encoding="utf-8")
        cls.demo = load_editor()
        cls.demo.LSP_WORKSPACE = cls.root
        cls.python_server = serve_editor(cls.demo.app)
        cls.python_url = cls.python_server.__enter__()
        cls.addClassCleanup(cls.python_server.__exit__, None, None, None)
        cls.go_server = go_editor(cls.root)
        cls.go_url = cls.go_server.__enter__()
        cls.addClassCleanup(cls.go_server.__exit__, None, None, None)
        cls.driver = sync_playwright().start()
        cls.addClassCleanup(cls.driver.stop)
        options = {"headless": True}
        if executable := os.environ.get("MARKATA_TEST_CHROMIUM"):
            options["executable_path"] = executable
        cls.browser = cls.driver.chromium.launch(**options)
        cls.addClassCleanup(cls.browser.close)

    def test_real_draft_file_and_scratch_are_independent(self):
        for name, url in [("Python", self.python_url), ("Go", self.go_url)]:
            with self.subTest(server=name), self.browser.new_context() as context:
                context.add_init_script("""const attach=Element.prototype.attachShadow;
                  Element.prototype.attachShadow=function(options){return attach.call(this,{...options,mode:'open'});};""")
                page=context.new_page()
                listing=page.request.get(url.rstrip('/')+'/workspace').json()
                root_key=listing.get('rootUri') or listing['root']
                context.add_init_script("localStorage.setItem(" + json.dumps('markata-workspace-v1:'+root_key+':draft.md') + ", '# Older scratch')")
                page.goto(url)
                preview=page.frame_locator('#preview')
                expect(preview.locator('h1')).to_have_text('Real workspace draft')
                page.locator('#site-files').get_by_role('treeitem',name='Unsaved draft.md',exact=True).click()
                expect(preview.locator('h1')).to_have_text('Older scratch')
                page.locator('#import').set_input_files({'name':'draft.md','mimeType':'text/markdown','buffer':b'# Imported scratch'})
                expect(preview.locator('h1')).to_have_text('Imported scratch')
                expect(page.locator('#active-file')).to_have_text('Unsaved draft.md')
                page.locator('#site-files').get_by_role('treeitem',name='draft.md',exact=True).click()
                expect(preview.locator('h1')).to_have_text('Real workspace draft')
                page.reload()
                expect(preview.locator('h1')).to_have_text('Real workspace draft')

    def test_external_edits_refresh_clean_buffers_and_preserve_drafts(self):
        source_path = self.root / "draft.md"
        original = source_path.read_text(encoding="utf-8")
        try:
            for name, url in [("Python", self.python_url), ("Go", self.go_url)]:
                with self.subTest(server=name), self.browser.new_context() as context:
                    source_path.write_text(original, encoding="utf-8")
                    page = context.new_page()
                    page.goto(url)
                    preview = page.frame_locator('#preview')
                    expect(preview.locator('h1')).to_have_text('Real workspace draft')
                    source_path.write_text('# External clean edit', encoding="utf-8")
                    page.reload()
                    expect(preview.locator('h1')).to_have_text('External clean edit')
                    page.evaluate("document.querySelector('.CodeMirror').CodeMirror.setValue('# Private browser edit')")
                    expect(preview.locator('h1')).to_have_text('Private browser edit')
                    source_path.write_text('# External conflicting edit', encoding="utf-8")
                    page.reload()
                    expect(preview.locator('h1')).to_have_text('Private browser edit')
                    expect(page.locator('#saved')).to_contain_text('File changed on disk')
                    with self.browser.new_context() as separate:
                        other = separate.new_page()
                        other.goto(url)
                        expect(other.frame_locator('#preview').locator('h1')).to_have_text('External conflicting edit')
        finally:
            source_path.write_text(original, encoding="utf-8")

    def test_link_completion_and_diagnostics(self):
        for name, url in [("Python", self.python_url), ("Go", self.go_url)]:
            with self.subTest(server=name):
                context = self.browser.new_context()
                try:
                    page = context.new_page()
                    errors = []
                    page.on("pageerror", lambda error: errors.append(str(error)))
                    page.goto(url)
                    expect(page.frame_locator("#preview").locator("h1")).to_have_text("Real workspace draft")
                    expect(page.locator("#language-status")).to_contain_text("LSP connected")
                    for source, wanted in [("😀 [[gre]]", "😀 [[greeting]]"),
                                           ("😀 [Hello](hello)", "😀 [Hello](hello%20world.md)"),
                                           ("[[", "[[hello")]:
                        page.evaluate("""source => {
                          const cm = document.querySelector('.CodeMirror').CodeMirror;
                          cm.setValue(source);
                          cm.setCursor(0, source.indexOf(']]') >= 0 ? source.indexOf(']]') :
                            source.endsWith(')') ? source.length - 1 : source.length);
                          cm.focus();
                        }""", source)
                        page.keyboard.press("Control+Space")
                        expect(page.locator(".CodeMirror-hints")).to_be_visible()
                        page.keyboard.press("Enter")
                        self.assertEqual(page.evaluate("document.querySelector('.CodeMirror').CodeMirror.getValue()"), wanted)
                    contextual = [
                        ("::::: md-i", 0, 10, "md-icon", "::::: md-icon"),
                        ("::::: md-icon {data-s", 0, 21, "data-size", "::::: md-icon {data-size"),
                        ("::::: md-icon {data-size=\"l\"", 0, 27, "lg", "::::: md-icon {data-size=\"lg\""),
                        ("::: wa-details {sum", 0, 20, "summary", "::: wa-details {summary"),
                        ("::: md-icon {data-size=\"lg\" data-s", 0, 34, "data-size", "::: md-icon {data-size=\"lg\" data-size"),
                        ("::: road e", 0, 10, "editorial", "::: road editorial"),
                        ("```mermaid\nflowchart T", 1, 11, "TD", "```mermaid\nflowchart TD"),
                        ("```mermaid\ng", 1, 1, "graph", "```mermaid\ngraph"),
                        ("```mermaid\nclassDef note c", 1, 15, "color", "```mermaid\nclassDef note color:"),
                    ]
                    for source, line, column, suggestion, expected in contextual:
                        page.evaluate("""({source,line,column}) => {
                          const cm=document.querySelector('.CodeMirror').CodeMirror;
                          cm.setValue(source); cm.setCursor(line,column); cm.focus();
                        }""", {"source":source,"line":line,"column":column})
                        page.keyboard.press("Control+Space")
                        expect(page.locator(".CodeMirror-hints")).to_be_visible()
                        candidates = page.locator(".CodeMirror-hint").filter(has_text=suggestion).all()
                        target = next((item for item in candidates if item.inner_text().split(" — ", 1)[0] == suggestion), None)
                        self.assertIsNotNone(target, f"missing exact suggestion {suggestion}")
                        target.click()
                        self.assertEqual(page.evaluate("document.querySelector('.CodeMirror').CodeMirror.getValue()"), expected)
                    # Automatic completion after actual keyboard input.
                    page.evaluate("""() => {const cm=document.querySelector('.CodeMirror').CodeMirror;
                      cm.setValue(''); cm.setCursor(0,0); cm.focus();}""")
                    page.keyboard.type("[[gre", delay=30)
                    expect(page.locator(".CodeMirror-hints")).to_be_visible()
                    page.keyboard.press("Enter")
                    self.assertEqual(page.evaluate("document.querySelector('.CodeMirror').CodeMirror.getValue()"), "[[greeting")
                    page.evaluate("""() => document.querySelector('.CodeMirror').CodeMirror.setValue('😀 [[missing]]')""")
                    expect(page.locator(".markata-lsp-diagnostic")).to_be_visible()
                    self.assertIn("missing", page.locator(".markata-lsp-diagnostic").first.get_attribute("title"))
                    expect(page.frame_locator("#preview").locator(".post-content")).to_contain_text("missing")
                    self.assertEqual(errors, [])
                    self.assertFalse((self.root / ".markata/browser-draft.md").exists())
                finally:
                    context.close()

    def test_delayed_completions_cannot_change_newer_buffer_or_cursor(self):
        for name, url in [("Python", self.python_url), ("Go", self.go_url)]:
            with self.subTest(server=name):
                context = self.browser.new_context()
                try:
                    context.add_init_script("""(() => {
                      const original = WebSocket.prototype.addEventListener;
                      WebSocket.prototype.addEventListener = function(type, listener, options) {
                        if (type !== 'message') return original.call(this,type,listener,options);
                        return original.call(this,type,event => {
                          const message=JSON.parse(event.data);
                          if (message.result?.items) setTimeout(() => {
                            listener(event); window.__completionDelivered=true;
                          },300);
                          else listener(event);
                        },options);
                      };
                    })();""")
                    page = context.new_page()
                    page.goto(url)
                    expect(page.frame_locator("#preview").locator("h1")).to_have_text("Real workspace draft")
                    expect(page.locator("#language-status")).to_contain_text("LSP connected")
                    for move_only in [False, True]:
                        page.evaluate("""() => {
                          window.__completionDelivered=false;
                          const cm=document.querySelector('.CodeMirror').CodeMirror;
                          cm.setValue('[[gre'); cm.setCursor(0,5); cm.focus();
                        }""")
                        page.keyboard.press("Control+Space")
                        page.evaluate("""moveOnly => {
                          const cm=document.querySelector('.CodeMirror').CodeMirror;
                          if (moveOnly) cm.setCursor(0,0); else cm.setValue('# Newer draft');
                        }""", move_only)
                        page.wait_for_function("window.__completionDelivered === true")
                        expect(page.locator(".CodeMirror-hints")).to_have_count(0)
                        self.assertEqual(page.evaluate("document.querySelector('.CodeMirror').CodeMirror.getValue()"),
                                         "[[gre" if move_only else "# Newer draft")
                finally:
                    context.close()

    def test_lsp_failure_preserves_editing_and_preview(self):
        for name, url in [("Python", self.python_url), ("Go", self.go_url)]:
            with self.subTest(server=name):
                context = self.browser.new_context()
                try:
                    page = context.new_page()
                    page.route("**/lsp/config", lambda route: route.abort())
                    page.goto(url)
                    expect(page.frame_locator("#preview").locator("h1")).to_have_text("Real workspace draft")
                    expect(page.locator("#language-status")).to_contain_text("unavailable")
                    page.evaluate("document.querySelector('.CodeMirror').CodeMirror.setValue('# Still editing')")
                    expect(page.frame_locator("#preview").locator("h1")).to_have_text("Still editing")
                finally:
                    context.close()


if __name__ == "__main__":
    unittest.main()
