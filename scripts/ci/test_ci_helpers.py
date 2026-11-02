"""Preflight and failure signature contract tests."""
import importlib.util
from pathlib import Path
import unittest


def load(name):
    spec = importlib.util.spec_from_file_location(name, Path(__file__).parent / (name + ".py"))
    module = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(module)
    return module


preflight = load("pr_check")
diagnose = load("diagnose")


class Helpers(unittest.TestCase):
    def test_docs_only(self):
        self.assertTrue(preflight.docs_only(["docs/index.md", "README.md"]))
        self.assertFalse(preflight.docs_only([]))
        self.assertFalse(preflight.docs_only(["docs/index.md", "pkg/x.go"]))
        self.assertFalse(preflight.docs_only([".github/workflows/ci.yml"]))

    def test_selection(self):
        self.assertTrue(preflight.affected(["pkg/workbench/app.go"], preflight.EDITORS))
        self.assertTrue(preflight.affected(["pkg/config/settings.go"], ("pkg/config/",)))
        self.assertTrue(preflight.affected(["go.mod"], preflight.RENDER, ("go.mod",)))

    def test_classifier(self):
        log = ("settings_docs_gen.go is stale; run go generate ./pkg/config\n"
               "foo.go:1: File is not properly formatted (gofmt)\n"
               "playwright._impl._errors.TimeoutError: Timeout 30000ms exceeded")
        self.assertEqual({a[0] for a in diagnose.classify(log)},
                         {"Generated artifacts", "Lint and formatting", "Playwright timeout"})
        self.assertIn("never retries", diagnose.report(log))


if __name__ == "__main__":
    unittest.main()
