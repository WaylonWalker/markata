"""Exercise the shared LSP adapter against both actual example servers.

Build MARKATA_LSP_BINARY and MARKATA_WORKBENCH_BINARY from this branch first.
The Python environment must contain a standalone-renderer Markata wheel.
"""
from contextlib import contextmanager
from pathlib import Path
import os
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

    def test_link_completion_and_diagnostics(self):
        for name, url in [("Python", self.python_url), ("Go", self.go_url)]:
            with self.subTest(server=name):
                context = self.browser.new_context()
                try:
                    page = context.new_page()
                    errors = []
                    page.on("pageerror", lambda error: errors.append(str(error)))
                    page.goto(url)
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
                    expect(page.locator("#language-status")).to_contain_text("unavailable")
                    page.evaluate("document.querySelector('.CodeMirror').CodeMirror.setValue('# Still editing')")
                    expect(page.frame_locator("#preview").locator("h1")).to_have_text("Still editing")
                finally:
                    context.close()


if __name__ == "__main__":
    unittest.main()
