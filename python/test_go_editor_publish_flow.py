"""Isolated browser coverage for the Go editor's publish flow.

Set MARKATA_GO_EDITOR_BINARY to a built markata-go executable to run this
end-to-end test. It creates a temporary site repository and bare remote.
"""
import base64
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
        quick_post = posts / "quick-publish.md"
        quick_post.write_text(
            "---\npublished: false\ndraft: false\n---\n\n# Quick before\n",
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
        cls.git(cls.site, "add", "markata-go.toml", "posts/publish-flow.md", "posts/long-gutter.md", "posts/discard-flow.md", "posts/external-edit.md", "posts/autosave-clean.md", "posts/quick-publish.md")
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
        source = page.evaluate("document.querySelector('.CodeMirror').CodeMirror.getValue()")
        self.edit_source(page, source.replace("# Original", "# Unpublished local edit"))
        expect(page.locator("#saved")).to_have_text("Saved to site", timeout=10000)
        expect(page.locator("#markata-gutter-status")).to_contain_text("Git change region")
        self.open_changes(page)
        panel = page.locator("#workspace-changes-panel")
        expect(panel.locator("#workspace-changes-list")).to_contain_text("posts/discard-flow.md")
        self.open_review(page)
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
        page.get_by_role("button", name="Close changes panel").click()

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
        self.open_review(page)
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
        context = self.browser.new_context(accept_downloads=True)
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
        expect(page.locator("#workspace-conflict-review")).to_be_visible()
        page.locator("#workspace-conflict-review").click()
        compare = page.locator("#workspace-conflict-dialog")
        expect(compare).to_be_visible()
        expect(compare.locator("#workspace-conflict-browser")).to_have_text(browser_draft)
        expect(compare.locator("#workspace-conflict-disk")).to_have_text(external)
        expect(compare.locator("#workspace-conflict-keep-browser")).to_be_enabled()
        expect(compare.locator("#workspace-conflict-use-disk")).to_be_enabled()
        compare.locator("#workspace-conflict-done").click()
        expect(compare).not_to_be_visible()
        # Comparing the two copies must never mutate the filesystem.
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

    @staticmethod
    def open_changes(page):
        page.locator("#workspace-tab-changes").click()

    @classmethod
    def open_review(cls, page):
        cls.open_changes(page)
        page.locator("#publish-review-open").click()

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
        self.open_review(page)
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
        self.open_review(page)
        review = page.locator("#publish-review")
        expect(review).to_be_visible()
        page.wait_for_function("typeof window.__releaseHeldReview === 'function'")
        review.locator("#publish-back").click()
        expect(review).not_to_be_visible()

        # The previous network response is deliberately still outstanding.
        self.open_review(page)
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
        source = page.evaluate("document.querySelector('.CodeMirror').CodeMirror.getValue()")
        self.edit_source(page, source.replace(
            "# Before\n\nInitial copy.", "# After\n\nRevised copy for browser verification."
        ))
        expect(page.locator("#saved")).to_have_text("Saved to site", timeout=10000)
        self.open_changes(page)
        panel = page.locator("#workspace-changes-panel")
        expect(panel.locator("#workspace-changes-list")).to_contain_text("posts/publish-flow.md")
        self.open_review(page)
        review = page.locator("#publish-review")
        expect(review).to_be_visible()
        expect(review.locator('[data-publish-mode="source"]')).to_have_attribute("aria-selected", "true")
        review.get_by_role("tab", name="Source diff").click()
        expect(review.locator('[data-publish-view="source"]')).to_be_visible()
        expect(review.locator('[data-publish-view="source"]')).to_contain_text("published")
        # Optional Pierre must actually initialize and render. A missing module,
        # failed chunk, or asynchronously rejected render is a regression.
        review.get_by_role("tab", name="Full diff").click()
        pierre = review.locator("#publish-pierre-host")
        expect(pierre).to_have_attribute("data-state", "ready", timeout=20000)
        expect(pierre.locator("pre[data-diff]").first).to_have_attribute("data-diff-type", "single")
        expect(pierre.locator("code").last).to_contain_text(
            "Revised copy for browser verification.", timeout=20000
        )
        expect(review.locator("#publish-pierre-error")).to_be_hidden()
        review.get_by_role("tab", name="Commit message").click()
        message = "Publish browser fixture: revised post"
        review.locator("#publish-commit-message").fill(message)
        review.get_by_role("tab", name="Rendered diff").click()
        review.get_by_role("button", name="Commit selected page").click()
        expect(review.get_by_role("button", name="Push to origin")).to_be_visible()
        # Immediately update the still-open Changes panel when publication
        # creates a commit; don't wait for the 15-second polling interval.
        expect(panel.locator("#workspace-changes-list")).not_to_contain_text("posts/publish-flow.md", timeout=3000)
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


    def test_editor_remains_typable_when_workspace_is_unavailable(self):
        context = self.browser.new_context()
        self.addCleanup(context.close)
        page = context.new_page()
        # The source editor must remain available if workspace discovery
        # fails; recovery stays local until the server becomes available.
        page.route(re.compile(r"/workspace(?:\\?|$)"), lambda route: route.abort())
        page.goto(self.url, wait_until="domcontentloaded")
        page.wait_for_selector(".CodeMirror")
        page.locator(".CodeMirror").click()
        page.keyboard.type("Local-first editor probe")
        page.wait_for_function("""() => document.querySelector('.CodeMirror').CodeMirror
            .getValue().includes('Local-first editor probe')""")

    def test_changes_panel_lists_all_unopened_files(self):
        # The workspace-wide inventory must include files never opened in
        # CodeMirror, staged content and changes outside the content glob.
        staged = self.site / "posts" / "changes-staged.md"
        untracked = self.site / "posts" / "changes-untracked.md"
        outside = self.site / "changes-outside.txt"
        staged.write_text("# Staged but not opened\n", encoding="utf-8")
        untracked.write_text("# New but not opened\n", encoding="utf-8")
        outside.write_text("Non-Markdown worktree change\n", encoding="utf-8")
        self.git(self.site, "add", "--", "posts/changes-staged.md")
        try:
            context = self.browser.new_context(viewport={"width": 390, "height": 844})
            self.addCleanup(context.close)
            page = context.new_page()
            page.goto(self.url)
            self.open_changes(page)
            panel = page.locator("#workspace-changes-panel")
            expect(panel).to_be_visible()
            expect(panel.locator("#workspace-changes-list")).to_contain_text("posts/changes-staged.md")
            expect(panel.locator("#workspace-changes-list")).to_contain_text("posts/changes-untracked.md")
            expect(panel.locator("#workspace-changes-list")).to_contain_text("changes-outside.txt")
            expect(panel.locator("#workspace-changes-list")).to_contain_text("staged")
            expect(panel.locator("#workspace-changes-list")).to_contain_text("New / untracked")
            expect(panel.locator("#workspace-changes-summary")).to_contain_text("files with changes")
            # A mobile bottom sheet should never exceed the device viewport.
            assert panel.bounding_box()["width"] <= 391
            panel.locator("#workspace-changes-close").click()
            expect(panel).not_to_be_visible()
        finally:
            self.git(self.site, "reset", "-q", "HEAD", "--", "posts/changes-staged.md")
            staged.unlink(missing_ok=True)
            untracked.unlink(missing_ok=True)
            outside.unlink(missing_ok=True)

    def test_changes_panel_opens_diff_without_full_post_review(self):
        post = self.site / "posts" / "external-edit.md"
        original = post.read_text(encoding="utf-8")
        post.write_text("# Changed without opening in editor\n", encoding="utf-8")
        try:
            context = self.browser.new_context()
            self.addCleanup(context.close)
            page = context.new_page()
            page.goto(self.url)
            self.open_changes(page)
            panel = page.locator("#workspace-changes-panel")
            row = panel.locator(".workspace-change").filter(has_text="posts/external-edit.md")
            expect(row).to_be_visible()
            panel.locator("#workspace-changes-review-open").click()
            review = page.locator("#workspace-changes-review")
            expect(review).to_be_visible()
            review.locator(".workspace-review-item").filter(has_text="posts/external-edit.md").click()
            expect(review.locator("#workspace-review-diff")).to_contain_text("Changed without opening")
            expect(page.locator("#publish-review")).not_to_be_visible()
        finally:
            post.write_text(original, encoding="utf-8")

    def test_quick_publish_only_active_post(self):
        context = self.browser.new_context()
        self.addCleanup(context.close)
        page = context.new_page()
        errors = []
        page.on("pageerror", lambda error: errors.append(str(error)))
        page.goto(self.url)
        self.open_site_file(page, "posts/quick-publish.md", errors)
        source = page.evaluate("document.querySelector('.CodeMirror').CodeMirror.getValue()")
        self.edit_source(page, source.replace("# Quick before", "# Quick after"))
        expect(page.locator("#saved")).to_have_text("Saved to site", timeout=10000)
        self.open_changes(page)
        panel = page.locator("#workspace-changes-panel")
        expect(panel.locator("#workspace-changes-list")).to_contain_text("posts/quick-publish.md")
        page.locator(".CodeMirror").click()
        page.keyboard.press("Control+S")
        expect(page.locator("#markata-publish-status")).to_contain_text("Published to Git", timeout=30000)
        expect(panel.locator("#workspace-changes-list")).not_to_contain_text("posts/quick-publish.md")
        self.assertIn("# Quick after", self.git(self.site, "show", "HEAD:posts/quick-publish.md"))
        self.assertIn("published: true", self.git(self.site, "show", "HEAD:posts/quick-publish.md"))
        self.assertEqual(self.git(self.site, "rev-parse", "HEAD"),
                         self.git(self.remote, "rev-parse", "refs/heads/main"))
        publish_button = page.locator("#publish-open")
        expect(publish_button).to_have_text("Up to date", timeout=15000)
        expect(publish_button).to_be_disabled()
        source = page.evaluate("document.querySelector('.CodeMirror').CodeMirror.getValue()")
        self.edit_source(page, source + "\nA later edit to publish.\n")
        expect(publish_button).to_have_text("Publish post", timeout=15000)
        expect(publish_button).to_be_enabled()
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
        self.open_review(page)
        review = page.locator("#publish-review")
        expect(review.locator("#publish-source-diff")).to_contain_text("locally edited")
        expect(review.locator("#publish-source-diff")).not_to_contain_text("Diff too large")
        assert review.locator(".publish-diff-line").count() < 40
        review.get_by_role("tab", name="Full diff").click()
        host = review.locator("#publish-pierre-host")
        expect(host).to_have_attribute("data-state", "ready", timeout=20000)
        # A one-line edit in a 600-line post must not reserve thousands of
        # pixels for unchanged lines and leave the real edit offscreen.
        size = host.evaluate("(element) => element.scrollHeight")
        self.assertLess(size, 1800, f"Pierre compact hunk layout too tall: {size}px")
        review.locator("#publish-back").click()
        self.assertEqual(errors, [])

    def test_mobile_cms_navigation_and_desktop_mode_switcher(self):
        # Desktop mode switcher and sidebar toggle
        desktop_ctx = self.browser.new_context(viewport={"width": 1440, "height": 900})
        self.addCleanup(desktop_ctx.close)
        page = desktop_ctx.new_page()
        errors = []
        page.on("pageerror", lambda error: errors.append(str(error)))
        page.goto(self.url)
        expect(page.locator(".header-modes")).to_be_visible()
        page.locator('.header-modes [data-layout-mode="editor"]').click()
        expect(page.locator("#workbench")).to_have_attribute("data-layout", "editor")
        page.locator('.header-modes [data-layout-mode="preview"]').click()
        expect(page.locator("#workbench")).to_have_attribute("data-layout", "preview")
        page.locator('.header-modes [data-layout-mode="split"]').click()
        expect(page.locator("#workbench")).to_have_attribute("data-layout", "split")
        page.locator("#sidebar-toggle").click()
        expect(page.locator("#workbench")).to_have_class(re.compile(r"\bsidebar-hidden\b"))
        page.locator("#sidebar-toggle").click()
        expect(page.locator("#workbench")).not_to_have_class(re.compile(r"\bsidebar-hidden\b"))
        self.assertEqual(errors, [])

        # Mobile full-height CMS app shell (Content / Write / Preview / New)
        mobile_ctx = self.browser.new_context(viewport={"width": 390, "height": 844})
        self.addCleanup(mobile_ctx.close)
        mpage = mobile_ctx.new_page()
        merrors = []
        mpage.on("pageerror", lambda error: merrors.append(str(error)))
        mpage.goto(self.url)
        expect(mpage.locator("#mobile-nav")).to_be_visible()
        expect(mpage.locator("#site-sidebar")).to_be_visible()
        expect(mpage.locator("#editor-pane")).to_be_hidden()
        expect(mpage.locator("#preview-pane")).to_be_hidden()
        self.assertGreater(mpage.locator("#site-sidebar").bounding_box()["height"], 500)
        self.assertLessEqual(mpage.evaluate("document.documentElement.scrollWidth"), 390)

        mpage.locator('#mobile-nav [data-mobile-tab="write"]').click()
        expect(mpage.locator("#editor-pane")).to_be_visible()
        expect(mpage.locator(".CodeMirror")).to_be_visible()
        expect(mpage.locator("#mobile-editing-toolbar")).to_be_visible()
        expect(mpage.locator("#site-sidebar")).to_be_hidden()
        self.assertGreater(mpage.locator("#editor-pane").bounding_box()["height"], 500)

        # Mobile header command palette button opens Command Palette directly from Write
        mpage.locator("#header-commands").click()
        expect(mpage.locator("dialog.markdown-commands[open]")).to_be_visible()
        mpage.keyboard.press("Escape")

        # Tapping header document chip opens Page Properties on mobile
        mpage.locator("#header-doc").click()
        expect(mpage.locator("#frontmatter-panel")).to_be_visible()
        mpage.locator("#frontmatter-toggle").click()

        mpage.locator('#mobile-nav [data-mobile-tab="preview"]').click()
        expect(mpage.locator("#preview-pane")).to_be_visible()
        expect(mpage.locator("#editor-pane")).to_be_hidden()
        self.assertGreater(mpage.locator("#preview-pane").bounding_box()["height"], 500)

        mpage.locator("#mobile-new-post").click()
        expect(mpage.locator("#template-picker")).to_be_visible()
        self.assertEqual(merrors, [])


    def test_cms_publish_menu_opens_changes_and_theme_surfaces_are_opaque(self):
        context = self.browser.new_context(viewport={"width": 1366, "height": 850})
        self.addCleanup(context.close)
        page = context.new_page()
        errors = []
        page.on("pageerror", lambda error: errors.append(str(error)))
        page.goto(self.url)

        # Recoloring must settle instead of feeding the root style observer
        # back into itself and freezing the editor.
        page.evaluate("document.documentElement.setAttribute('data-theme', 'dark')")
        expect(page.locator("#header-doc")).to_be_visible()
        self.assertIn(
            page.evaluate("document.documentElement.style.getPropertyValue('--cms-on-accent')"),
            ("#ffffff", "#11151d"),
        )

        page.locator("#sidebar-toggle").click()
        expect(page.locator("#workbench")).to_have_class(re.compile(r"\bsidebar-hidden\b"))
        page.locator("#publish-menu-toggle").click()
        menu = page.locator("#publish-menu")
        expect(menu).to_be_visible()
        self.assertNotEqual(
            menu.evaluate("(node) => getComputedStyle(node).backgroundColor"),
            "rgba(0, 0, 0, 0)",
        )
        menu.locator('[data-publish-action="changes"]').click()
        expect(page.locator("#workbench")).not_to_have_class(re.compile(r"\bsidebar-hidden\b"))
        expect(page.locator("#workspace-tab-changes")).to_have_attribute("aria-selected", "true")
        self.assertEqual(errors, [])

    def test_cms_quick_note_post_metadata_and_distinct_filenames(self):
        context = self.browser.new_context(viewport={"width": 1280, "height": 800})
        self.addCleanup(context.close)
        page = context.new_page()
        errors = []
        seen_methods = []
        page.on("pageerror", lambda error: errors.append(str(error)))

        def metadata_route(route):
            seen_methods.append(route.request.method)
            route.fulfill(
                status=200,
                content_type="application/json",
                body='{"title":"Quick Capture Fixture","description":"Found from link metadata"}',
            )

        page.route("**/workspace/link-metadata", metadata_route)
        page.goto(self.url)
        page.locator('#cms-quick-actions [data-cms-flow="quick"]').click()
        dialog = page.locator("#cms-quick-note-dialog")
        expect(dialog).to_be_visible()
        dialog.locator("#cms-qn-url").fill("https://example.org/quick-capture")
        dialog.locator("#cms-qn-body").click()
        expect(dialog.locator("#cms-qn-heading")).to_have_value("Quick Capture Fixture")
        self.assertEqual(seen_methods, ["POST"])
        dialog.locator("#cms-qn-submit").click()
        expect(dialog).not_to_be_visible(timeout=10000)
        expect(page.locator("#active-file")).to_contain_text("quick-capture-fixture-", timeout=10000)
        first_path = page.locator("#active-file").inner_text()
        self.assertIn("quick-capture-fixture-", first_path)

        # Another capture with the same title should create a different file.
        page.locator('#cms-quick-actions [data-cms-flow="quick"]').click()
        expect(dialog).to_be_visible()
        dialog.locator("#cms-qn-heading").fill("Quick Capture Fixture")
        dialog.locator("#cms-qn-body").fill("Second independent note")
        dialog.locator("#cms-qn-submit").click()
        expect(dialog).not_to_be_visible(timeout=10000)
        # The editor reloads when a newly created source enters the tree.
        # Ignore the temporary scratch selection during that reload.
        page.wait_for_function(
            """(previous) => {
                const current = document.querySelector('#active-file')?.textContent?.trim() || '';
                return current.includes('quick-capture-fixture-') && current !== previous;
            }""",
            arg=first_path,
            timeout=10000,
        )
        second_path = page.locator("#active-file").inner_text()
        self.assertNotEqual(first_path, second_path)
        self.assertTrue((self.site / first_path).exists())
        self.assertTrue((self.site / second_path).exists())
        self.assertEqual(errors, [])


    def test_pwa_android_multi_photo_share_capture_and_save(self):
        context = self.browser.new_context(viewport={"width": 390, "height": 844})
        self.addCleanup(context.close)
        page = context.new_page()
        errors = []
        page.on("pageerror", lambda error: errors.append(str(error)))
        page.goto(self.url)
        page.evaluate("""async () => {
            const registration = await navigator.serviceWorker.ready;
            if (!registration.active) throw new Error('Share service worker not active');
            if (!navigator.serviceWorker.controller) {
                await new Promise((resolve, reject) => {
                    const timeout = setTimeout(() => reject(new Error('No controller')), 10000);
                    navigator.serviceWorker.addEventListener('controllerchange', () => {
                        clearTimeout(timeout); resolve();
                    }, {once: true});
                });
            }
        }""")
        shared_url = page.evaluate("""async () => {
            // A tiny, valid PNG sent through the same Web Share Target multipart
            // POST that Android delivers to an installed Markata PWA.
            const base64 = 'iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVQIHWP4z8DwHwAFgAI/ScLhqQAAAABJRU5ErkJggg==';
            const bytes = Uint8Array.from(atob(base64), c => c.charCodeAt(0));
            const form = new FormData();
            form.set('share_title', 'Weekend photos');
            form.set('share_text', 'Two images from Android');
            form.append('images', new File([bytes], 'first.png', {type: 'image/png'}));
            form.append('images', new File([bytes], 'second.png', {type: 'image/png'}));
            const response = await fetch('/share-target', {method: 'POST', body: form});
            const capture = await response.json();
            if (!response.ok || !capture.share_id) throw new Error('Share target failed: ' + JSON.stringify(capture));
            return location.origin + '/?share_id=' + encodeURIComponent(capture.share_id);
        }""")
        self.assertIn("share_id=", shared_url)
        page.goto(shared_url)
        dialog = page.locator("#cms-quick-note-dialog")
        expect(dialog).to_be_visible(timeout=10000)
        expect(dialog.locator("#cms-qn-heading")).to_have_value("Weekend photos")
        expect(dialog.locator("#cms-qn-body")).to_have_value("Two images from Android")
        expect(dialog.locator("#cms-qn-photo-list img")).to_have_count(2)
        self.assertFalse(list((self.site / "static" / "images").glob("weekend-photos*"))
                         if (self.site / "static" / "images").exists() else False)

        # The share remains recoverable if the PWA is closed before saving.
        dialog.locator(".cms-flow-close").click()
        page.reload()
        expect(dialog).to_be_visible(timeout=10000)
        expect(dialog.locator("#cms-qn-photo-list img")).to_have_count(2)
        dialog.locator("#cms-qn-submit").click()
        expect(dialog).not_to_be_visible(timeout=15000)
        expect(page.locator("#active-file")).to_contain_text("weekend-photos-", timeout=15000)
        self.assertNotIn("share_id=", page.url)

        post_path = page.locator("#active-file").inner_text()
        source = (self.site / post_path).read_text(encoding="utf-8")
        self.assertIn("![Shared photo 1](/images/", source)
        self.assertIn("![Shared photo 2](/images/", source)
        self.assertIn('image: "/images/', source)
        self.assertEqual(errors, [])



    def test_android_share_target_handles_real_multipart_post_navigation(self):
        context = self.browser.new_context(viewport={"width": 390, "height": 844})
        self.addCleanup(context.close)
        page = context.new_page()
        page.goto(self.url)
        page.evaluate("""async () => {
            await navigator.serviceWorker.ready;
            if (!navigator.serviceWorker.controller) {
                await new Promise((resolve, reject) => {
                    const timeout = setTimeout(() => reject(Error('No controller')), 10000);
                    navigator.serviceWorker.addEventListener('controllerchange', () => {
                        clearTimeout(timeout); resolve();
                    }, {once: true});
                });
            }
            const form = document.createElement('form');
            form.action = '/share-target';
            form.method = 'POST';
            form.enctype = 'multipart/form-data';
            const title = document.createElement('input');
            title.name = 'share_title';
            title.value = 'From Android POST navigation';
            const file = document.createElement('input');
            file.type = 'file';
            file.name = 'images';
            file.id = 'pwa-share-navigation-photo';
            form.append(title, file);
            document.body.append(form);
        }""")
        tiny_png = base64.b64decode(
            "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVQIHWP4z8DwHwAFgAI/ScLhqQAAAABJRU5ErkJggg=="
        )
        page.locator("#pwa-share-navigation-photo").set_input_files(
            {"name": "nav.png", "mimeType": "image/png", "buffer": tiny_png}
        )
        with page.expect_navigation(wait_until="domcontentloaded"):
            page.evaluate("document.querySelector('form[action=\\\"/share-target\\\"]').submit()")
        self.assertIn("share_id=", page.url)
        dialog = page.locator("#cms-quick-note-dialog")
        expect(dialog).to_be_visible(timeout=10000)
        expect(dialog.locator("#cms-qn-heading")).to_have_value("From Android POST navigation")
        expect(dialog.locator("#cms-qn-photo-list img")).to_have_count(1)
        dialog.locator("#cms-qn-discard-share").click()
        expect(dialog).not_to_be_visible()
        self.assertNotIn("share_id=", page.url)

    def test_pwa_android_single_photo_is_visible_in_post_body(self):
        context = self.browser.new_context(viewport={"width": 390, "height": 844})
        self.addCleanup(context.close)
        page = context.new_page()
        errors = []
        page.on("pageerror", lambda error: errors.append(str(error)))
        page.goto(self.url)
        page.evaluate("""async () => {
            await navigator.serviceWorker.ready;
            if (!navigator.serviceWorker.controller) {
                await new Promise((resolve, reject) => {
                    const timeout = setTimeout(() => reject(new Error('No controller')), 10000);
                    navigator.serviceWorker.addEventListener('controllerchange', () => {
                        clearTimeout(timeout); resolve();
                    }, {once: true});
                });
            }
        }""")
        shared_url = page.evaluate("""async () => {
            const base64 = 'iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVQIHWP4z8DwHwAFgAI/ScLhqQAAAABJRU5ErkJggg==';
            const bytes = Uint8Array.from(atob(base64), c => c.charCodeAt(0));
            const form = new FormData();
            form.set('share_title', 'Single shared photo');
            form.append('images', new File([bytes], 'picture.png', {type: 'image/png'}));
            const response = await fetch('/share-target', {method: 'POST', body: form});
            const capture = await response.json();
            if (!response.ok || !capture.share_id) throw Error('Share target failed: ' + JSON.stringify(capture));
            return location.origin + '/?share_id=' + encodeURIComponent(capture.share_id);
        }""")
        page.goto(shared_url)
        dialog = page.locator("#cms-quick-note-dialog")
        expect(dialog).to_be_visible(timeout=10000)
        expect(dialog.locator("#cms-qn-photo-list img")).to_have_count(1)
        dialog.locator("#cms-qn-submit").click()
        expect(dialog).not_to_be_visible(timeout=15000)
        expect(page.locator("#active-file")).to_contain_text("single-shared-photo-", timeout=15000)
        source = (self.site / page.locator("#active-file").inner_text()).read_text(encoding="utf-8")
        self.assertIn('image: "/images/', source)
        self.assertIn("![Shared photo 1](/images/", source)
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
        self.open_review(page)
        review = page.locator("#publish-review")
        expect(review.locator("#publish-no-changes")).to_be_visible()
        expect(review.locator("#publish-no-changes")).to_contain_text("All caught up")
        expect(review.locator("#publish-discard")).not_to_be_visible()
        expect(review.locator("#publish-confirm")).to_be_disabled()
        expect(review.locator(".publish-mode-tabs")).not_to_be_visible()
        review.locator("#publish-back").click()


if __name__ == "__main__":
    unittest.main()
