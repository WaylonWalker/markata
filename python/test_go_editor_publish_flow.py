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
            "---\ntitle: Publish flow fixture\npublished: true\ndraft: false\n"
            "tags:\n  - initial\n---\n\n# Before\n\nInitial copy.\n",
            encoding="utf-8",
        )
        cls.git(cls.site, "add", "posts/publish-flow.md")
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

    def test_frontmatter_review_commit_message_and_temporary_push(self):
        context = self.browser.new_context()
        self.addCleanup(context.close)
        page = context.new_page()
        errors = []
        page.on("pageerror", lambda error: errors.append(str(error)))
        page.goto(self.url)
        expect(page.locator("#active-file")).to_contain_text("posts/publish-flow.md")

        page.get_by_role("button", name="Properties", exact=True).click()
        tags = page.locator('[data-field="tags"]')
        tags.get_by_role("textbox", name="New tag").fill("browser-test")
        tags.get_by_role("button", name="Add", exact=True).click()
        expect(tags.locator(".fm-tag-chip")).to_have_count(2)
        expect(tags.locator(".fm-tag-chip").nth(1)).to_contain_text("browser-test")
        published = page.locator('[data-field="published"] input[type="checkbox"]')
        published.uncheck()
        published.check()
        expect(page.locator("#source")).to_have_value(re.compile(r"published: true"))
        expect(page.locator("#source")).to_have_value(re.compile(r"browser-test"))

        page.get_by_role("button", name="Back to Markdown").click()
        page.evaluate("""() => {
          const editor=document.querySelector('.CodeMirror').CodeMirror;
          editor.setValue(editor.getValue().replace('# Before\\n\\nInitial copy.', '# After\\n\\nRevised copy for browser verification.'));
        }""")
        page.get_by_role("button", name="Publish updates").click()
        review = page.locator("#publish-review")
        expect(review).to_be_visible()
        expect(review.locator('[data-publish-mode="rich"]')).to_have_attribute("aria-selected", "true")
        review.get_by_role("tab", name="Source diff").click()
        expect(review.locator('[data-publish-view="source"]')).to_be_visible()
        review.get_by_role("tab", name="Commit message").click()
        message = "Publish browser fixture: revised post"
        review.locator("#publish-commit-message").fill(message)
        review.get_by_role("tab", name="Rendered diff").click()
        review.get_by_role("button", name="Commit selected page").click()
        expect(review.get_by_role("button", name="Push to origin")).to_be_visible()
        self.assertEqual(self.git(self.site, "log", "-1", "--format=%B"), message)

        commit = self.git(self.site, "rev-parse", "HEAD")
        review.get_by_role("button", name="Push to origin").click()
        expect(review.get_by_role("button", name="Check live status")).to_be_visible()
        self.assertEqual(self.git(self.remote, "rev-parse", "refs/heads/main"), commit)
        self.assertEqual(self.git(self.remote, "log", "-1", "--format=%B", "refs/heads/main"), message)
        self.assertEqual(errors, [])


if __name__ == "__main__":
    unittest.main()
