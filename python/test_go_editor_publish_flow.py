"""Isolated browser coverage for the Go editor's publish flow.

Set MARKATA_GO_EDITOR_BINARY to a built markata-go executable to run this
end-to-end test. It creates a temporary site repository and bare remote.
"""
import os
import re
import socket
import subprocess
import tempfile
import time
import unittest
from pathlib import Path

from playwright.sync_api import expect, sync_playwright


@unittest.skipUnless(os.environ.get("MARKATA_GO_EDITOR_BINARY"), "Set MARKATA_GO_EDITOR_BINARY to a built markata-go executable")
class GoEditorPublishFlowTest(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.temp = tempfile.TemporaryDirectory(prefix="markata-publish-browser-")
        cls.root = Path(cls.temp.name)
        cls.site = cls.root / "site"
        cls.site.mkdir()
        cls.remote = cls.root / "origin.git"
        cls.git(cls.root, "init", "--bare", str(cls.remote))
        cls.git(cls.site, "init", "-b", "main")
        cls.git(cls.site, "config", "user.name", "Markata Browser Test")
        cls.git(cls.site, "config", "user.email", "browser-test@example.invalid")
        cls.git(cls.site, "remote", "add", "origin", str(cls.remote))
        posts = cls.site / "posts"
        posts.mkdir()
        post = posts / "publish-flow.md"
        post.write_text(
            "---\ntitle: Publish flow fixture\npublished: false\ndraft: false\n"
            "tags:\n  - initial\n---\n\n# Before\n\nInitial copy.\n",
            encoding="utf-8",
        )
        long_post = posts / "long-gutter.md"
        long_post.write_text(
            "---\ntitle: Long gutter fixture\npublished: false\n---\n\n"
            + "".join(f"Line {index}: stable content\n" for index in range(600)),
            encoding="utf-8",
        )
        discard_post = posts / "discard-flow.md"
        discard_post.write_text("---\\ntitle: Discard fixture\\npublished: false\\n---\\n\\n# Original\\n", encoding="utf-8")
        external_post = posts / "external-edit.md"
        external_post.write_text("# First external version\\n", encoding="utf-8")
        clean_post = posts / "autosave-clean.md"
        clean_post.write_text("# Initial clean draft\\n", encoding="utf-8")
        cls.git(cls.site, "add", "posts/publish-flow.md", "posts/long-gutter.md", "posts/discard-flow.md", "posts/external-edit.md", "posts/autosave-clean.md")
        cls.git(cls.site, "commit", "-m", "Fixture baseline")
        cls.git(cls.site, "push", "-u", "origin", "main")

        with socket.socket() as sock:
            sock.bind(("127.0.0.1", 0))
            cls.port = sock.getsockname()[1]
        cls.url = f"http://127.0.0.1:{cls.port}/"
        cls.server = subprocess.Popen(
            [
                os.environ["MARKATA_GO_EDITOR_BINARY"], "editor",
                "--workspace", str(cls.site), "--listen", f"127.0.0.1:{cls.port}",
            ],
            stdout=subprocess.DEVNULL,
            stderr=subprocess.DEVNULL,
        )
        for _ in range(100):
            if cls.server.poll() is not None:
                raise RuntimeError("Go editor exited before the browser test started")
            try:
                with socket.create_connection(("127.0.0.1", cls.port), timeout=0.2):
                    break
            except OSError:
                time.sleep(0.1)
        else:
            cls.server.terminate()
            raise RuntimeError("Go editor did not start")

        cls.driver = sync_playwright().start()
        options = {"headless": True}
        if executable := os.environ.get("MARKATA_TEST_CHROMIUM"):
            options["executable_path"] = executable
        cls.browser = cls.driver.chromium.launch(**options)

    @classmethod
    def tearDownClass(cls):
        if hasattr(cls, "browser"):
            cls.browser.close()
            cls.driver.stop()
        if hasattr(cls, "server"):
            cls.server.terminate()
            try:
                cls.server.wait(timeout=5)
            except subprocess.TimeoutExpired:
                cls.server.kill()
        if hasattr(cls, "temp"):
            cls.temp.cleanup()

    @staticmethod
    def git(root, *args):
        return subprocess.run(
            ["git", *args], cwd=root, check=True, text=True,
            stdout=subprocess.PIPE, stderr=subprocess.PIPE,
        ).stdout.strip()

    def test_discard_restores_disk_browser_and_gutter(self):
        context = self.browser.new_context()
        self.addCleanup(context.close)
        page = context.new_page()
        errors = []
        page.on("pageerror", lambda error: errors.append(str(error)))
        page.goto(self.url)
        page.evaluate("window.MarkataSiteTree.openFile('posts/discard-flow.md')")
        expect(page.locator("#active-file")).to_contain_text("posts/discard-flow.md")
        expect(page.locator("#markata-gutter-status")).to_contain_text("No Git changes")
        original = (self.site / "posts/discard-flow.md").read_text(encoding="utf-8")
        page.evaluate("""() => {
          const editor=document.querySelector('.CodeMirror').CodeMirror;
          editor.setValue(editor.getValue().replace('# Original', '# Unpublished local edit'));
        }""")
        expect(page.locator("#markata-gutter-status")).to_contain_text("Git change region")
        page.get_by_role("button", name="Review changes").click()
        review = page.locator("#publish-review")
        expect(review).to_be_visible()
        expect(review.locator("#publish-discard")).to_be_enabled()
        review.locator("#publish-discard").click()
        expect(review.locator("#publish-discard-confirmation")).to_be_visible()
        expect(review.locator("#publish-discard-description")).to_contain_text("posts/discard-flow.md")
        review.locator("#publish-discard-accept").click()
        expect(review).not_to_be_visible(timeout=10000)
        expect(page.locator("#markata-gutter-status")).to_contain_text("No Git changes")
        self.assertEqual((self.site / "posts/discard-flow.md").read_text(encoding="utf-8"), original)
        self.assertEqual(page.evaluate("document.querySelector('.CodeMirror').CodeMirror.getValue()"), original)
        page.reload()
        expect(page.locator("#active-file")).to_contain_text("posts/discard-flow.md")
        self.assertEqual(page.evaluate("document.querySelector('.CodeMirror').CodeMirror.getValue()"), original)
        self.assertEqual(self.git(self.site, "diff", "--cached", "--name-only"), "")
        self.assertEqual(errors, [])

    def test_autosave_preserves_external_edit_and_browser_recovery(self):
        context = self.browser.new_context()
        self.addCleanup(context.close)
        page = context.new_page()
        errors = []
        page.on("pageerror", lambda error: errors.append(str(error)))
        page.goto(self.url)
        page.evaluate("window.MarkataSiteTree.openFile('posts/external-edit.md')")
        expect(page.locator("#active-file")).to_contain_text("posts/external-edit.md")
        expect(page.locator("#markata-gutter-status")).to_contain_text("No Git changes")

        # Another process changes the file after the browser has loaded it.
        post = self.site / "posts/external-edit.md"
        external = "# Changed by external editor\\n"
        post.write_text(external, encoding="utf-8")
        browser_draft = "# Uncommitted browser edit\\n"
        page.evaluate("""(text) => document.querySelector('.CodeMirror').CodeMirror.setValue(text)""", browser_draft)
        expect(page.locator("#saved")).to_contain_text("Site save failed", timeout=10000)
        expect(page.locator("#saved")).to_contain_text("browser recovery copy kept")
        self.assertEqual(post.read_text(encoding="utf-8"), external)

        # The dirty recovery copy must survive a reload without replacing disk.
        page.reload()
        expect(page.locator("#active-file")).to_contain_text("posts/external-edit.md")
        expect(page.locator("#saved")).to_contain_text("File changed on disk")
        self.assertEqual(page.evaluate("document.querySelector('.CodeMirror').CodeMirror.getValue()"), browser_draft)
        self.assertEqual(post.read_text(encoding="utf-8"), external)
        self.assertEqual(errors, [])

    def test_autosave_ack_updates_recovery_version(self):
        context = self.browser.new_context()
        self.addCleanup(context.close)
        page = context.new_page()
        errors = []
        page.on("pageerror", lambda error: errors.append(str(error)))
        page.goto(self.url)
        page.evaluate("window.MarkataSiteTree.openFile('posts/autosave-clean.md')")
        expect(page.locator("#active-file")).to_contain_text("posts/autosave-clean.md")
        post = self.site / "posts/autosave-clean.md"
        page.evaluate("""() => document.querySelector('.CodeMirror').CodeMirror.setValue('# First autosave\\n')""")
        expect(page.locator("#saved")).to_have_text("Saved to site", timeout=10000)
        self.assertEqual(post.read_text(encoding="utf-8"), "# First autosave\\n")

        page.reload()
        expect(page.locator("#active-file")).to_contain_text("posts/autosave-clean.md")
        self.assertEqual(page.evaluate("document.querySelector('.CodeMirror').CodeMirror.getValue()"), "# First autosave\\n")
        page.evaluate("""() => document.querySelector('.CodeMirror').CodeMirror.setValue('# Second autosave\\n')""")
        expect(page.locator("#saved")).to_have_text("Saved to site", timeout=10000)
        self.assertEqual(post.read_text(encoding="utf-8"), "# Second autosave\\n")
        self.assertEqual(errors, [])

    def test_frontmatter_review_commit_message_and_temporary_push(self):
        context = self.browser.new_context()
        self.addCleanup(context.close)
        page = context.new_page()
        errors = []
        page.on("pageerror", lambda error: errors.append(str(error)))
        page.goto(self.url)
        page.evaluate("window.MarkataSiteTree.openFile('posts/publish-flow.md')")
        expect(page.locator("#active-file")).to_contain_text("posts/publish-flow.md")

        page.get_by_role("button", name="Properties", exact=True).click()
        tags = page.locator('[data-field="tags"]')
        tags.get_by_role("textbox", name="New tag").fill("browser-test")
        tags.get_by_role("button", name="Add", exact=True).click()
        expect(tags.locator(".fm-tag-chip")).to_have_count(2)
        expect(tags.locator(".fm-tag-chip").nth(1)).to_contain_text("browser-test")
        published = page.locator('[data-field="published"] input[type="checkbox"]')
        published.check()
        expect(page.locator("#source")).to_have_value(re.compile(r"published: true"))
        expect(page.locator("#source")).to_have_value(re.compile(r"browser-test"))

        page.get_by_role("button", name="Back to Markdown").click()
        page.evaluate("""() => {
          const editor=document.querySelector('.CodeMirror').CodeMirror;
          editor.setValue(editor.getValue().replace('# Before\\n\\nInitial copy.', '# After\\n\\nRevised copy for browser verification.'));
        }""")
        page.get_by_role("button", name="Review changes").click()
        review = page.locator("#publish-review")
        expect(review).to_be_visible()
        expect(review.locator('[data-publish-mode="rich"]')).to_have_attribute("aria-selected", "true")
        review.get_by_role("tab", name="Source diff").click()
        expect(review.locator('[data-publish-view="source"]')).to_be_visible()
        expect(review.locator('[data-publish-view="source"]')).to_contain_text("published")
        review.get_by_role("tab", name="Commit message").click()
        message = "Publish browser fixture: revised post"
        review.locator("#publish-commit-message").fill(message)
        review.get_by_role("tab", name="Rendered diff").click()
        review.get_by_role("button", name="Commit selected page").click()
        expect(review.get_by_role("button", name="Push to origin")).to_be_visible()
        self.assertEqual(self.git(self.site, "log", "-1", "--format=%B"), message)
        committed_post = self.git(self.site, "show", "HEAD:posts/publish-flow.md")
        self.assertIn("published: true", committed_post)
        self.assertIn("browser-test", committed_post)
        self.assertIn("Revised copy for browser verification.", committed_post)

        commit = self.git(self.site, "rev-parse", "HEAD")
        review.get_by_role("button", name="Push to origin").click()
        expect(review.get_by_role("button", name="Check live status")).to_be_visible()
        self.assertEqual(self.git(self.remote, "rev-parse", "refs/heads/main"), commit)
        self.assertEqual(self.git(self.remote, "log", "-1", "--format=%B", "refs/heads/main"), message)
        self.assertEqual(errors, [])


    def test_git_gutter_large_document(self):
        context = self.browser.new_context()
        self.addCleanup(context.close)
        page = context.new_page()
        errors = []
        page.on("pageerror", lambda error: errors.append(str(error)))
        page.goto(self.url)
        page.evaluate("window.MarkataSiteTree.openFile('posts/publish-flow.md')")
        expect(page.locator("#active-file")).to_contain_text("posts/publish-flow.md")
        page.evaluate("window.MarkataSiteTree.openFile('posts/long-gutter.md')")
        expect(page.locator("#active-file")).to_contain_text("posts/long-gutter.md")
        expect(page.locator("#markata-gutter-status")).to_contain_text("No Git changes")
        page.evaluate("""() => {
          const editor=document.querySelector('.CodeMirror').CodeMirror;
          const line=300;
          editor.replaceRange(' locally edited', {line, ch:editor.getLine(line).length});
          editor.scrollIntoView({line,ch:0});
        }""")
        expect(page.locator("#markata-gutter-status")).to_contain_text("Git change region")
        marker = page.locator(".CodeMirror-gutter-wrapper .markata-change-marker")
        expect(marker.first).to_be_visible()
        self.assertGreater(marker.count(), 0)
        self.assertEqual(errors, [])


if __name__ == "__main__":
    unittest.main()
