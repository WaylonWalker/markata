"""Regression tests for the shared local/CI PR preflight."""
import contextlib
import io
import sys
import unittest
from unittest.mock import patch

import pr_check


class EmptyDiffTests(unittest.TestCase):
    def test_empty_format_check_is_a_valid_noop(self):
        with (
            patch.object(sys, "argv", ["pr_check.py", "--base", "HEAD^", "--only-format"]),
            patch.object(pr_check, "changed", return_value=[]) as changed,
            patch.object(pr_check, "format_go") as format_go,
            contextlib.redirect_stdout(io.StringIO()),
        ):
            pr_check.main()
        changed.assert_called_once_with("HEAD^")
        format_go.assert_not_called()

    def test_full_preflight_still_rejects_empty_diff(self):
        with (
            patch.object(sys, "argv", ["pr_check.py", "--base", "HEAD^"]),
            patch.object(pr_check, "changed", return_value=[]),
        ):
            with self.assertRaisesRegex(RuntimeError, "no changed paths"):
                pr_check.main()

    def test_format_check_still_formats_changed_paths(self):
        with (
            patch.object(sys, "argv", ["pr_check.py", "--base", "HEAD^", "--only-format"]),
            patch.object(pr_check, "changed", return_value=["pkg/example.go"]),
            patch.object(pr_check, "format_go") as format_go,
            contextlib.redirect_stdout(io.StringIO()),
        ):
            pr_check.main()
        format_go.assert_called_once_with(["pkg/example.go"])

    def test_format_check_does_not_swallow_invalid_base(self):
        with (
            patch.object(sys, "argv", ["pr_check.py", "--base", "invalid", "--only-format"]),
            patch.object(pr_check, "changed", side_effect=RuntimeError("bad base")),
        ):
            with self.assertRaisesRegex(RuntimeError, "bad base"):
                pr_check.main()


if __name__ == "__main__":
    unittest.main()
