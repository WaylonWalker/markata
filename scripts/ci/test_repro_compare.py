"""Offline tests for the deterministic output-tree comparator."""

from __future__ import annotations

import contextlib
import io
import json
import tempfile
import unittest
from pathlib import Path

from repro_compare import compare_trees, main


class ReproCompareTests(unittest.TestCase):
    def setUp(self) -> None:
        self.temp = tempfile.TemporaryDirectory()
        self.addCleanup(self.temp.cleanup)
        self.left = Path(self.temp.name) / "left"
        self.right = Path(self.temp.name) / "right"
        self.left.mkdir()
        self.right.mkdir()

    def put(self, root: Path, name: str, data: bytes) -> None:
        path = root / name
        path.parent.mkdir(parents=True, exist_ok=True)
        path.write_bytes(data)

    def test_identical_content_and_nested_paths(self) -> None:
        for root in (self.left, self.right):
            self.put(root, "index.html", b"<h1>same</h1>")
            self.put(root, "blog/post/index.html", b"post")
        result = compare_trees(self.left, self.right)
        self.assertTrue(result["identical"])
        self.assertEqual(result["left_file_count"], 2)
        self.assertEqual(result["changed"], [])

    def test_reports_exact_changed_added_and_removed_paths(self) -> None:
        self.put(self.left, "b.css", b"before")
        self.put(self.right, "b.css", b"after")
        self.put(self.left, "only-old.md", b"old")
        self.put(self.right, "only-new.md", b"new")
        result = compare_trees(self.left, self.right)
        self.assertFalse(result["identical"])
        self.assertEqual(result["only_left"], ["only-old.md"])
        self.assertEqual(result["only_right"], ["only-new.md"])
        self.assertEqual([v["path"] for v in result["changed"]], ["b.css"])
        self.assertEqual(len(result["changed"][0]["left"]["sha256"]), 64)

    def test_exclusions_are_explicit_and_reported(self) -> None:
        # Preserve a genuine nonempty build artifact while excluding the
        # explicitly variable metadata file.
        self.put(self.left, "index.html", b"same")
        self.put(self.right, "index.html", b"same")
        self.put(self.left, "generated/build.json", b"timestamp=one")
        self.put(self.right, "generated/build.json", b"timestamp=two")
        self.assertFalse(compare_trees(self.left, self.right)["identical"])
        result = compare_trees(self.left, self.right, ["generated/*"])
        self.assertTrue(result["identical"])
        self.assertEqual(result["excluded_patterns"], ["generated/*"])

    def test_cli_exit_codes_and_machine_readable_json(self) -> None:
        self.put(self.left, "index.html", b"same")
        self.put(self.right, "index.html", b"same")
        out = io.StringIO()
        with contextlib.redirect_stdout(out):
            self.assertEqual(main([str(self.left), str(self.right)]), 0)
        self.assertTrue(json.loads(out.getvalue())["identical"])
        self.put(self.right, "index.html", b"different")
        with contextlib.redirect_stdout(io.StringIO()):
            self.assertEqual(main([str(self.left), str(self.right)]), 1)

    def test_empty_outputs_are_unverified_by_default(self) -> None:
        with self.assertRaisesRegex(ValueError, "empty generated output"):
            compare_trees(self.left, self.right)
        result = compare_trees(self.left, self.right, allow_empty=True)
        self.assertTrue(result["identical"])

        self.put(self.left, "index.html", b"ok")
        with self.assertRaisesRegex(ValueError, "empty generated output"):
            compare_trees(self.left, self.right)

        with contextlib.redirect_stderr(io.StringIO()):
            self.assertEqual(main([str(self.left), str(self.right)]), 2)
        out = io.StringIO()
        with contextlib.redirect_stdout(out):
            self.assertEqual(
                main([str(self.left), str(self.right), "--allow-empty"]), 1
            )
        self.assertFalse(json.loads(out.getvalue())["identical"])

    def test_invalid_directory_is_not_claimed_equal(self) -> None:
        with contextlib.redirect_stderr(io.StringIO()):
            self.assertEqual(main([str(self.left), str(self.right / "missing")]), 2)

    def test_symlink_is_rejected_instead_of_followed(self) -> None:
        target = Path(self.temp.name) / "outside"
        target.write_bytes(b"secret")
        (self.left / "link").symlink_to(target)
        with self.assertRaisesRegex(ValueError, "symlink not allowed"):
            compare_trees(self.left, self.right)


if __name__ == "__main__":
    unittest.main()
