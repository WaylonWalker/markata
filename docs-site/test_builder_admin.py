"""Docs builder regression smoke checks without cluster access."""
import importlib.util
import os
import tempfile
import unittest
from pathlib import Path

spec = importlib.util.spec_from_file_location("docs_builder_admin", Path(__file__).with_name("builder_admin.py"))
module = importlib.util.module_from_spec(spec)
spec.loader.exec_module(module)


class BuilderTests(unittest.TestCase):
    def test_release_validation(self):
        self.assertTrue(module.safe_release("20261018T1630-abcdef1234"))
        for bad in ("", "..", "../secrets", ".hidden", "x/y", "x?y"):
            self.assertFalse(module.safe_release(bad))

    def test_publish_rollback_retention(self):
        with tempfile.TemporaryDirectory() as tmp:
            old_site, old_versions, old_current, old_state, old_logdir = (
                module.SITE, module.VERSIONS, module.CURRENT, module.STATE_FILE, module.LOG_DIR
            )
            try:
                module.SITE = Path(tmp) / "site"
                module.VERSIONS = module.SITE / "releases"
                module.CURRENT = module.SITE / "current"
                module.STATE_FILE = module.SITE / ".docs-builder/state.json"
                module.LOG_DIR = module.SITE / ".docs-builder/logs"
                b = module.Builder()
                generated = Path(tmp) / "generated"
                for name in ("index.html", "docs/quickstart/index.html", "docs/guides/authoring/index.html", "assets/site.css"):
                    path = generated / name
                    path.parent.mkdir(parents=True, exist_ok=True)
                    path.write_text("ok")
                name = b.publish(generated, "abc123456789012")
                self.assertTrue(module.CURRENT.is_symlink())
                self.assertEqual(module.CURRENT.resolve().name, name)
                b.rollback(name)
                self.assertEqual(module.CURRENT.resolve().name, name)
                with self.assertRaises(ValueError):
                    b.rollback("../secrets")
            finally:
                module.SITE, module.VERSIONS, module.CURRENT, module.STATE_FILE, module.LOG_DIR = (
                    old_site, old_versions, old_current, old_state, old_logdir
                )


if __name__ == "__main__":
    unittest.main()
