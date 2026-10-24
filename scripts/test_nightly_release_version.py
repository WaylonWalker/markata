#!/usr/bin/env python3
"""Tests for dated 0.x nightly prerelease version calculation."""
import importlib.util
from pathlib import Path
import unittest

spec = importlib.util.spec_from_file_location(
    "nightly_release_version", Path(__file__).with_name("nightly-release-version.py")
)
module = importlib.util.module_from_spec(spec)
spec.loader.exec_module(module)


class NightlyVersionTests(unittest.TestCase):
    def test_matching_github_and_pypi_versions(self):
        self.assertEqual(
            module.calculate(["v0.13.0", "v0.14.0", "v0.14.1-dev.20261020"], "20261022"),
            ("v0.14.1-dev.20261022", "0.14.1.dev20261022", "0.14.1-dev.20261022"),
        )

    def test_new_stable_release_advances_patch(self):
        self.assertEqual(
            module.calculate(["v0.14.0", "v0.14.1"], "20261101")[1], "0.14.2.dev20261101"
        )

    def test_no_stable_release_fails_closed(self):
        with self.assertRaisesRegex(ValueError, "No stable"):
            module.calculate(["v1.0.0", "v0.1.0-rc.1"], "20261022")

    def test_malformed_date_fails(self):
        with self.assertRaises(ValueError):
            module.calculate(["v0.14.0"], "20261345")


if __name__ == "__main__":
    unittest.main()
