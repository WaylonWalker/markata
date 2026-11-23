"""Offline tests for the machine-readable Build Lab evidence gate."""

from __future__ import annotations

import copy
import unittest

from buildlab_evidence import summarize


def passing_result() -> dict:
    return {
        "schema_version": 1,
        "verdict": "pass",
        "fixture": {"digest": "a" * 64},
        "candidate": {"binary_sha256": "b" * 64},
        "environment": {"goos": "linux", "goarch": "amd64"},
        "checkpoints": [
            {
                "operation": "build",
                "candidate_clean": {"successful": True, "duration_seconds": 2.0},
                "candidate_determinism": {"successful": True, "duration_seconds": 1.9},
                "candidate_incremental": {"successful": True, "duration_seconds": 0.2},
                "correctness": {
                    "deterministic_equal": True,
                    "differential_equal": True,
                    "incremental_applicable": True,
                    "incremental_equal": True,
                },
            }
        ],
    }


class BuildLabEvidenceTests(unittest.TestCase):
    def test_reports_proven_rebuild_and_timing(self):
        output = summarize(passing_result())
        self.assertIn("Verdict: **pass**", output)
        self.assertIn("2.00", output)
        self.assertIn("Identical?", output)

    def test_rejects_no_build_results(self):
        result = passing_result()
        result["checkpoints"] = []
        with self.assertRaisesRegex(ValueError, "no build checkpoints"):
            summarize(result)

    def test_rejects_false_pass_without_repeat(self):
        result = passing_result()
        del result["checkpoints"][0]["candidate_determinism"]
        with self.assertRaisesRegex(ValueError, "no repeated clean-build evidence"):
            summarize(result)

    def test_rejects_changed_bytes_even_with_pass_verdict(self):
        result = passing_result()
        result["checkpoints"][0]["correctness"]["deterministic_equal"] = False
        with self.assertRaisesRegex(ValueError, "repeated clean outputs differ"):
            summarize(result)

    def test_rejects_failed_child(self):
        result = passing_result()
        result["checkpoints"][0]["candidate_clean"]["successful"] = False
        with self.assertRaisesRegex(ValueError, "clean build did not succeed"):
            summarize(result)

    def test_rejects_changed_incremental_output(self):
        result = passing_result()
        result["checkpoints"][0]["correctness"]["incremental_equal"] = False
        with self.assertRaisesRegex(ValueError, "incremental output differs"):
            summarize(result)

    def test_rejects_failed_verdict_and_invalid_schema(self):
        result = passing_result()
        result["verdict"] = "fail"
        with self.assertRaisesRegex(ValueError, "Build Lab verdict: fail"):
            summarize(result)
        result["schema_version"] = 0
        with self.assertRaisesRegex(ValueError, "unsupported"):
            summarize(result)


if __name__ == "__main__":
    unittest.main()
