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
        # Make discovery deterministic: the default glob must not decide
        # whether direct posts/*.md children match a recursive glob.
        (cls.site / "markata-go.toml").write_text(
            '[markata-go.glob]\npatterns = ["posts/*.md"]\n',
            encoding="utf-8",
        )
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
        discard_post.write_text("---\ntitle: Discard fixture\npublished: false\n---\n\n# Original\n", encoding="utf-8")
        external_post = posts / "external-edit.md"
        external_post.write_text("# First external version\n", encoding="utf-8")
        clean_post = posts / "autosave-clean.md"
        clean_post.write_text("# Initial clean draft\n", encoding="utf-8")
        cls.git(cls.site, "add", "markata-go.toml", "posts/publish-flow.md", "posts/long-gutter.md", "posts/discard-flow.md", "posts/external-edit.md", "posts/autosave-clean.md")
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
        self.open_site_file(page, "posts/discard-flow.md", errors)
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

        # Regression: a Discard followed by fresh typing must not reuse the
        # previous review iframe or the old preview cursor/scroll geometry.
        fresh = original + "\nFirst fresh paragraph.\n\nSecond fresh paragraph.\n"
        self.edit_source(page, fresh)
        expect(page.locator("#saved")).to_have_text("Saved to site", timeout=10000)
        expect(page.frame_locator("#preview").locator(".post-content")).to_contain_text("Second fresh paragraph.")
        expect(page.locator("#markata-gutter-status")).to_contain_text("Git change region")
        page.evaluate("""() => {
          const editor=document.querySelector('.CodeMirror').CodeMirror;
          const line=editor.getValue().split('\\n').findIndex(text=>text.includes('Second fresh paragraph.'));
          editor.focus(); editor.setCursor({line,ch:4});
        }""")
        expect(page.frame_locator("#preview").locator(".markata-sync-position")).to_contain_text(
            "Second fresh paragraph.", timeout=10000
        )
        # The reverse interaction must also land on the correct fresh source
        # line, rather than an anchor retained from before Discard.
        page.frame_locator("#preview").locator("p").filter(has_text="Second fresh paragraph.").click()
        page.wait_for_function("""() => {
          const editor=document.querySelector('.CodeMirror').CodeMirror;
          return editor.getLine(editor.getCursor().line).includes('Second fresh paragraph.');
        }""")
        # Additional edits can arrive while an earlier render is still in
        # flight. The preview must eventually reflect the newest revision,
        # including during rapid typing immediately after a Discard.
        for number in range(4):
            fresh += f"\nLive revision {number}.\n"
            self.edit_source(page, fresh)
            page.wait_for_timeout(140)
        expect(page.frame_locator("#preview").locator(".post-content")).to_contain_text(
            "Live revision 3.", timeout=15000
        )
        expect(page.frame_locator("#preview").locator(".post-content")).not_to_contain_text(
            "Unpublished local edit"
        )
        page.get_by_role("button", name="Review changes").click()
        expect(review).to_be_visible()
        expect(review.locator("#publish-source-diff")).to_contain_text("Live revision 3.")
        expect(review.locator("#publish-source-diff")).to_contain_text("First fresh paragraph.")
        expect(review.locator("#publish-source-diff")).to_contain_text("Second fresh paragraph.")
        expect(review.locator("#publish-source-diff")).not_to_contain_text("Unpublished local edit")
        expect(review.locator('[data-publish-mode="source"]')).to_have_attribute("aria-selected", "true")
        # Large rendered pages are deliberately loaded only when requested.
        review.locator('[data-publish-mode="rich"]').click()
        expect(page.frame_locator("#publish-rich-frame").locator(".post-content")).to_contain_text(
            "Live revision 3.", timeout=10000
        )
        expect(page.frame_locator("#publish-rich-frame").locator(".post-content")).to_contain_text(
            "Second fresh paragraph.", timeout=10000
        )
        expect(page.frame_locator("#publish-rich-frame").locator(".post-content")).not_to_contain_text(
            "Unpublished local edit"
        )
        review.locator("#publish-back").click()
        expect(review).not_to_be_visible()
        self.assertEqual((self.site / "posts/discard-flow.md").read_text(encoding="utf-8"), fresh)
        page.reload()
        expect(page.locator("#active-file")).to_contain_text("posts/discard-flow.md")
        self.assertEqual(page.evaluate("document.querySelector('.CodeMirror').CodeMirror.getValue()"), fresh)
        self.assertEqual(self.git(self.site, "diff", "--cached", "--name-only"), "")
        self.assertEqual(errors, [])

    def test_autosave_preserves_external_edit_and_browser_recovery(self):
        context = self.browser.new_context()
        self.addCleanup(context.close)
        page = context.new_page()
        errors = []
        page.on("pageerror", lambda error: errors.append(str(error)))
        page.goto(self.url)
        self.open_site_file(page, "posts/external-edit.md", errors)
        expect(page.locator("#markata-gutter-status")).to_contain_text("No Git changes")

        # Another process changes the file after the browser has loaded it.
        post = self.site / "posts/external-edit.md"
        external = "# Changed by external editor\n"
        post.write_text(external, encoding="utf-8")
        browser_draft = "# Uncommitted browser edit\n"
        self.edit_source(page, browser_draft)
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
        self.open_site_file(page, "posts/autosave-clean.md", errors)
        post = self.site / "posts/autosave-clean.md"
        self.edit_source(page, "# First autosave\n")
        expect(page.locator("#saved")).to_have_text("Saved to site", timeout=10000)
        self.assertEqual(post.read_text(encoding="utf-8"), "# First autosave\n")

        page.reload()
        expect(page.locator("#active-file")).to_contain_text("posts/autosave-clean.md")
        self.assertEqual(page.evaluate("document.querySelector('.CodeMirror').CodeMirror.getValue()"), "# First autosave\n")
        self.edit_source(page, "# Second autosave\n")
        expect(page.locator("#saved")).to_have_text("Saved to site", timeout=10000)
        self.assertEqual(post.read_text(encoding="utf-8"), "# Second autosave\n")
        self.assertEqual(errors, [])

    @staticmethod
    def edit_source(page, content):
        # CodeMirror.setValue() is a load/reset operation and deliberately
        # does not invoke browser autosave. Use a real edit origin instead.
        page.evaluate("""source => {
          const editor=document.querySelector('.CodeMirror').CodeMirror;
          editor.replaceRange(
            source, {line:0,ch:0},
            editor.posFromIndex(editor.getValue().length), '+input');
        }""", content)

    @staticmethod
    def open_site_file(page, path, errors=None):
        # The initial /workspace fetch is async. Earlier tests asked openFile
        # before the registry was populated, so the call silently did nothing.
        try:
            page.wait_for_function(
                "(target) => Boolean(window.MarkataSiteTree?.files?.().some(item => item.path === target))",
                arg=path,
                timeout=10000,
            )
        except Exception as error:
            debug = page.evaluate("""() => ({
              files:window.MarkataSiteTree?.files?.().map(file => file.path),
              visibleError:document.querySelector('#error')?.textContent,
              activeFile:window.MarkataSiteTree?.activeFile?.(),
              readyState:document.readyState,
              codeMirror:!!document.querySelector('.CodeMirror')?.CodeMirror,
              shortcutCount:window.MarkataShortcuts?.all?.().length,
              publishDialog:!!document.querySelector('#publish-review'),
            })""")
            server = page.request.get(page.url.rstrip("/") + "/workspace")
            raise AssertionError(
                f"workspace file {path!r} never became available: browser={debug}, "
                f"javascript_errors={errors}, server={server.status} {server.text()[:1500]}"
            ) from error
        page.evaluate("(target) => window.MarkataSiteTree.openFile(target)", path)
        expect(page.locator("#active-file")).to_have_text(path)

    def test_discard_new_untracked_post_removes_source_and_only_its_recovery(self):
        context = self.browser.new_context()
        self.addCleanup(context.close)
        page = context.new_page()
        errors = []
        page.on("pageerror", lambda error: errors.append(str(error)))
        page.goto(self.url)

        path = "posts/untracked-discard.md"
        original = "# New untracked draft\\n"
        created = page.request.post(
            self.url.rstrip("/") + "/workspace/file?path=" + path,
            data=original,
            headers={"Origin": self.url.rstrip("/"), "Content-Type": "text/plain"},
        )
        self.assertEqual(created.status, 201, created.text())
        page.reload()
        self.open_site_file(page, path, errors)
        self.assertEqual(page.evaluate("document.querySelector('.CodeMirror').CodeMirror.getValue()"), original)
        # Keep a separate browser recovery copy untouched.
        page.evaluate("""() => {
          localStorage.setItem('markata-unrelated-draft-test', 'keep');
          window.MarkataSiteTree.save();
        }""")
        page.get_by_role("button", name="Review changes").click()
        review = page.locator("#publish-review")
        expect(review).to_be_visible()
        expect(review.locator("#publish-before-name")).to_have_text("New page")
        expect(review.locator("#publish-discard")).to_be_enabled()
        review.locator("#publish-discard").click()
        expect(review.locator("#publish-discard-description")).to_contain_text("Permanently delete")
        expect(review.locator("#publish-discard-description")).to_contain_text(path)
        review.locator("#publish-discard-accept").click()
        page.wait_for_function("""() => window.MarkataSiteTree?.activeFile?.() !== 'posts/untracked-discard.md'""", timeout=15000)
        self.assertFalse((self.site / path).exists())
        self.assertEqual(self.git(self.site, "diff", "--cached", "--name-only"), "")
        self.assertIsNone(page.evaluate("""() => {
          const id = window.MarkataSiteTree.workspaceId();
          return localStorage.getItem('markata-workspace-v1:' + id + ':buffer:file:posts/untracked-discard.md');
        }"""))
        self.assertEqual(page.evaluate("localStorage.getItem('markata-unrelated-draft-test')"), "keep")
        self.assertEqual(errors, [])

    def test_review_close_inflight_request_does_not_block_reopening(self):
        context = self.browser.new_context()
        self.addCleanup(context.close)
        page = context.new_page()
        errors = []
        page.on("pageerror", lambda error: errors.append(str(error)))
        page.goto(self.url)
        self.open_site_file(page, "posts/autosave-clean.md", errors)

        # Hold the first review response after the user opens the dialog.
        # A closed review must not block a fresh one or overwrite its state.
        page.evaluate("""() => {
          const realFetch = window.fetch.bind(window);
          let held = false;
          window.fetch = (input, options) => {
            const url = typeof input === 'string' ? input : input.url;
            if (!held && url.startsWith('/workspace/publish/review?')) {
              held = true;
              return new Promise((resolve, reject) => {
                window.__releaseHeldReview = () =>
                  realFetch(input, options).then(resolve, reject);
              });
            }
            return realFetch(input, options);
          };
        }""")
        page.get_by_role("button", name="Review changes").click()
        review = page.locator("#publish-review")
        expect(review).to_be_visible()
        page.wait_for_function("typeof window.__releaseHeldReview === 'function'")
        review.locator("#publish-back").click()
        expect(review).not_to_be_visible()

        # The previous network response is deliberately still outstanding.
        page.get_by_role("button", name="Review changes").click()
        expect(review).to_be_visible()
        expect(review.locator("#publish-body")).to_be_visible(timeout=10000)
        expect(review.locator("#publish-path")).to_have_text("posts/autosave-clean.md")
        page.evaluate("window.__releaseHeldReview()")
        expect(review.locator("#publish-body")).to_be_visible()
        expect(review).to_have_attribute("data-unavailable", "false")
        self.assertEqual(errors, [])
        review.locator("#publish-back").click()

    def test_frontmatter_review_commit_message_and_temporary_push(self):
        context = self.browser.new_context()
        self.addCleanup(context.close)
        page = context.new_page()
        errors = []
        page.on("pageerror", lambda error: errors.append(str(error)))
        page.goto(self.url)
        self.open_site_file(page, "posts/publish-flow.md", errors)

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
        expect(review.locator('[data-publish-mode="source"]')).to_have_attribute("aria-selected", "true")
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
        self.open_site_file(page, "posts/publish-flow.md", errors)
        self.open_site_file(page, "posts/long-gutter.md", errors)
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
        # 600-line post with a single edit: do not compare every unchanged
        # line or incorrectly announce that the diff is too large.
        page.get_by_role("button", name="Review changes").click()
        review = page.locator("#publish-review")
        expect(review.locator("#publish-source-diff")).to_contain_text("locally edited")
        expect(review.locator("#publish-source-diff")).not_to_contain_text("Diff too large")
        assert review.locator(".publish-diff-line").count() < 40
        review.locator("#publish-back").click()
        self.assertEqual(errors, [])

    def test_00_review_unchanged_file_is_immediate_and_clear(self):
        context = self.browser.new_context()
        self.addCleanup(context.close)
        page = context.new_page()
        errors = []
        page.on("pageerror", lambda error: errors.append(str(error)))
        page.goto(self.url)
        self.open_site_file(page, "posts/long-gutter.md", errors)
        expect(page.locator("#markata-gutter-status")).to_contain_text("No Git changes")
        page.get_by_role("button", name="Review changes").click()
        review = page.locator("#publish-review")
        expect(review.locator("#publish-no-changes")).to_be_visible()
        expect(review.locator("#publish-no-changes")).to_contain_text("All caught up")
        expect(review.locator("#publish-discard")).not_to_be_visible()
        expect(review.locator("#publish-confirm")).to_be_disabled()
        expect(review.locator(".publish-mode-tabs")).not_to_be_visible()
        review.locator("#publish-back").click()


if __name__ == "__main__":
    unittest.main()
