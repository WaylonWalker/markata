#!/usr/bin/env python3
"""Validate and summarize Build Lab production reproducibility evidence.

This is only a result consumer. Build Lab itself runs the pinned, isolated
production builds and produces the source-of-truth versioned JSON result.
"""

from __future__ import annotations

import argparse
import json
import sys
from pathlib import Path


def summarize(result: dict) -> str:
    if result.get("schema_version") != 1:
        raise ValueError("unsupported or missing Build Lab result schema")
    checkpoints = result.get("checkpoints")
    if not isinstance(checkpoints, list) or not checkpoints:
        raise ValueError("no build checkpoints; a missing build cannot pass")

    rows = [
        "## Build Lab: isolated production-build evidence",
        "",
        f"- Verdict: **{result.get('verdict', 'missing')}**",
        f"- Fixture SHA-256: \`{result.get('fixture', {}).get('digest', 'missing')}\`",
        f"- Builder SHA-256: \`{result.get('candidate', {}).get('binary_sha256', 'missing')}\`",
        f"- Runtime: \`{result.get('environment', {}).get('goos', '?')}/{result.get('environment', {}).get('goarch', '?')}\`",
        "",
        "| Checkpoint | Clean (s) | Repeat clean (s) | Incremental (s) | Identical? |",
        "|---|---:|---:|---:|---|",
    ]
    failures = []
    for i, item in enumerate(checkpoints, start=1):
        equality = item.get("correctness", {})
        clean = item.get("candidate_clean", {})
        repeat = item.get("candidate_determinism")
        incremental = item.get("candidate_incremental", {})
        if not isinstance(repeat, dict):
            failures.append(f"checkpoint {i}: no repeated clean-build evidence")
            repeat = {}
        if not clean.get("successful") or not repeat.get("successful"):
            failures.append(f"checkpoint {i}: a clean build did not succeed")
        if not equality.get("deterministic_equal"):
            failures.append(f"checkpoint {i}: repeated clean outputs differ")
        if not equality.get("differential_equal"):
            failures.append(f"checkpoint {i}: baseline and candidate differ")
        if equality.get("incremental_applicable") and not equality.get("incremental_equal"):
            failures.append(f"checkpoint {i}: incremental output differs from clean")
        if equality.get("incremental_applicable") and not incremental.get("successful"):
            failures.append(f"checkpoint {i}: incremental build failed")
        rows.append(
            f"| {i}: {item.get('operation', '?')} | "
            f"{clean.get('duration_seconds', 0):.2f} | "
            f"{repeat.get('duration_seconds', 0):.2f} | "
            f"{incremental.get('duration_seconds', 0):.2f} | "
            f"{'yes' if equality.get('deterministic_equal') else 'NO'} |"
        )
    rows.extend(
        [
            "",
            "These timings are observations on this runner, not a cross-machine speed guarantee.",
            "This check compares published output using Build Lab's manifest contract.",
        ]
    )
    if result.get("verdict") != "pass":
        failures.append(f"Build Lab verdict: {result.get('verdict', 'missing')}")
    if result.get("failure_class"):
        failures.append(f"Failure class: {result['failure_class']}")
    if failures:
        raise ValueError("; ".join(failures))
    return "\n".join(rows) + "\n"


def main(argv: list[str] | None = None) -> int:
    parser = argparse.ArgumentParser(description="Check Build Lab JSON and create a human-readable summary")
    parser.add_argument("result", type=Path)
    parser.add_argument("--summary", type=Path, help="append verified evidence to the GitHub run summary")
    args = parser.parse_args(argv)
    try:
        result = json.loads(args.result.read_text(encoding="utf-8"))
        rendered = summarize(result)
        print(rendered)
        if args.summary:
            with args.summary.open("a", encoding="utf-8") as target:
                target.write(rendered)
    except (OSError, ValueError, TypeError, KeyError) as error:
        print(f"Build Lab evidence invalid: {error}", file=sys.stderr)
        return 1
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
