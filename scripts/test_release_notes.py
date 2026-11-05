import importlib.util
import tempfile
import unittest
from pathlib import Path

spec = importlib.util.spec_from_file_location("release_notes", Path(__file__).with_name("release-notes.py"))
module = importlib.util.module_from_spec(spec)
spec.loader.exec_module(module)

class ReleaseNotesTests(unittest.TestCase):
    def test_deterministic_category_order_and_pr_links(self):
        with tempfile.TemporaryDirectory() as d:
            root = Path(d)
            (root / "23.fix.md").write_text("Repair archive pagination.\n")
            (root / "12.feature.md").write_text("Add release previews.\n")
            (root / "5.breaking.md").write_text("Rename CLI. Migrate by updating the command.\n")
            result = module.build(root, "v1.2.3")
            self.assertLess(result.index("Breaking changes"), result.index("Features"))
            self.assertLess(result.index("Features"), result.index("Fixes"))
            self.assertIn("Repair archive pagination. (#23)", result)
            self.assertEqual(result, module.build(root, "v1.2.3"))

    def test_invalid_filename_fails(self):
        with tempfile.TemporaryDirectory() as d:
            (Path(d) / "not-a-fragment.md").write_text("Something\n")
            with self.assertRaises(ValueError):
                module.build(Path(d), "v1.0.0")

    def test_empty_and_heading_fragments_fail(self):
        with tempfile.TemporaryDirectory() as d:
            p = Path(d) / "10.fix.md"
            for content in ("", "# heading"):
                p.write_text(content)
                with self.assertRaises(ValueError):
                    module.build(Path(d), "v1.0.0")

if __name__ == "__main__":
    unittest.main()
