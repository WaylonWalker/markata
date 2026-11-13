"""Contract tests for the nonblocking Go test timing reporter."""
import io
import json
import os
from pathlib import Path
import tempfile
import unittest
from unittest import mock

import go_test_timing


def event(action, package="example.org/pkg", **kwargs):
    return json.dumps({"Action": action, "Package": package, **kwargs}) + "\n"


class TimingReporterTests(unittest.TestCase):
    def test_records_package_and_nested_test_without_double_counting(self):
        packages, tests = {}, {}
        self.assertEqual(
            go_test_timing.read_event(event("output", Output="hello\n"),
                                      packages, tests), "hello\n")
        go_test_timing.read_event(
            event("pass", Test="TestOne/subtest", Elapsed=0.21), packages, tests)
        go_test_timing.read_event(event("pass", Elapsed=2.5), packages, tests)
        self.assertEqual(packages["example.org/pkg"]["seconds"], 2.5)
        self.assertEqual(tests[("example.org/pkg", "TestOne/subtest")]["seconds"], 0.21)
        self.assertEqual(len(packages), 1)
        self.assertEqual(go_test_timing.read_event("compiler error\n", packages, tests),
                         "compiler error\n")

    def test_ranking_and_wall_time_are_distinct(self):
        report = go_test_timing.build_report(
            "test", ["go", "test", "-json", "./..."], 0, 4.2,
            {"slow": {"seconds": 3.0, "result": "pass"},
             "fast": {"seconds": 1.0, "result": "pass"}},
            {("fast", "TestOne"): {"seconds": 0.1, "result": "pass"}},
        )
        self.assertEqual(report["wall_seconds"], 4.2)
        self.assertEqual(report["packages"][0]["package"], "slow")
        self.assertEqual(report["tests_observed"], 1)
        self.assertIn("may overlap", go_test_timing.markdown_report(report))

    def test_failure_keeps_exit_code_and_writes_json_and_summary(self):
        class FakeProcess:
            stdout = iter([event("output", Output="FAIL\n"),
                           event("fail", Elapsed=0.3)])
            def __enter__(self):
                return self
            def __exit__(self, *args):
                return False
            def wait(self):
                return 9

        with tempfile.TemporaryDirectory() as root:
            target = Path(root) / "metrics" / "go.json"
            summary = Path(root) / "step.md"
            with mock.patch.object(go_test_timing.subprocess, "Popen",
                                   return_value=FakeProcess()), \
                 mock.patch.dict(os.environ, {"GITHUB_STEP_SUMMARY": str(summary)}), \
                 mock.patch("sys.stdout", new_callable=io.StringIO):
                code = go_test_timing.main([
                    "--label", "sample", "--output", str(target),
                    "--", "go", "test", "-json", "./...",
                ])
            self.assertEqual(code, 9)
            saved = json.loads(target.read_text(encoding="utf-8"))
            self.assertEqual(saved["exit_code"], 9)
            self.assertEqual(saved["packages"][0]["result"], "fail")
            self.assertIn("Result: fail", summary.read_text(encoding="utf-8"))


if __name__ == "__main__":
    unittest.main()
