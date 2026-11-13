#!/usr/bin/env python3
"""Stream `go test -json`, keeping its output and recording package timings.

This profiles test execution without rerunning tests or changing their selection.
Package elapsed times can overlap; command wall time includes build/setup work.
"""
import argparse
import json
import os
from pathlib import Path
import subprocess
import sys
import time


def read_event(line, packages, tests):
    """Update timing dictionaries and return human-readable output for a line."""
    try:
        event = json.loads(line)
    except (ValueError, TypeError):
        return line
    if not isinstance(event, dict) or "Action" not in event:
        return line

    action = event["Action"]
    package = event.get("Package")
    test = event.get("Test")
    if action in ("pass", "fail", "skip") and package:
        elapsed = event.get("Elapsed")
        if isinstance(elapsed, (int, float)) and not isinstance(elapsed, bool):
            record = {"seconds": round(float(elapsed), 3), "result": action}
            if test:
                tests[(package, test)] = record
            else:
                packages[package] = record
    if action == "output":
        return event.get("Output", "")
    return ""


def build_report(label, command, code, wall, packages, tests):
    ranked_packages = sorted(
        ({"package": name, **values} for name, values in packages.items()),
        key=lambda item: (-item["seconds"], item["package"]),
    )
    ranked_tests = sorted(
        ({"package": package, "test": name, **values}
         for (package, name), values in tests.items()),
        key=lambda item: (-item["seconds"], item["package"], item["test"]),
    )
    return {
        "schema_version": 1,
        "label": label,
        "exit_code": code,
        "wall_seconds": round(wall, 3),
        "github_sha": os.environ.get("GITHUB_SHA", ""),
        "github_job": os.environ.get("GITHUB_JOB", ""),
        "runner_os": os.environ.get("RUNNER_OS", ""),
        "packages": ranked_packages,
        "tests_observed": len(tests),
        "slowest_tests": ranked_tests[:10],
    }


def markdown_report(report):
    lines = [
        "### Go test timing: " + report["label"],
        "",
        f"- Result: {'pass' if report['exit_code'] == 0 else 'fail'} (exit {report['exit_code']})",
        f"- Command wall time: {report['wall_seconds']:.2f}s (includes compile and setup)",
        f"- Reported packages: {len(report['packages'])}; test cases observed: {report['tests_observed']}",
        "- Package elapsed times may overlap; they must not be summed as wall time.",
        "",
    ]
    if report["packages"]:
        lines.extend(["| Slowest packages | Elapsed | Result |", "| --- | ---: | --- |"])
        for item in report["packages"][:10]:
            lines.append(f"| `{item['package']}` | {item['seconds']:.2f}s | {item['result']} |")
    else:
        lines.append("No package completion events captured (check for compilation/setup failure).")
    return "\n".join(lines) + "\n"


def main(argv=None):
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--label", required=True)
    parser.add_argument("--output", required=True, help="JSON report path")
    parser.add_argument("command", nargs=argparse.REMAINDER)
    args = parser.parse_args(argv)
    command = args.command[1:] if args.command[:1] == ["--"] else args.command
    if command[:2] != ["go", "test"] or "-json" not in command:
        parser.error("expected -- go test -json [other Go test options]")

    packages, tests = {}, {}
    started = time.perf_counter()
    try:
        with subprocess.Popen(
            command, stdout=subprocess.PIPE, stderr=subprocess.STDOUT,
            text=True, encoding="utf-8", errors="replace", bufsize=1,
        ) as process:
            for line in process.stdout:
                output = read_event(line, packages, tests)
                if output:
                    sys.stdout.write(output)
                    sys.stdout.flush()
            code = process.wait()
    except OSError as error:
        print(f"Could not start Go tests: {error}", file=sys.stderr)
        code = 127

    report = build_report(args.label, command, code,
                          time.perf_counter() - started, packages, tests)
    destination = Path(args.output)
    destination.parent.mkdir(parents=True, exist_ok=True)
    destination.write_text(json.dumps(report, indent=2) + "\n", encoding="utf-8")
    summary = markdown_report(report)
    print("\n" + summary, flush=True)
    if path := os.environ.get("GITHUB_STEP_SUMMARY"):
        with open(path, "a", encoding="utf-8") as file:
            file.write(summary)
    return code


if __name__ == "__main__":
    sys.exit(main())
