#!/usr/bin/env python3
"""Summarize failure signatures without changing checks or retrying runs."""
import argparse
from pathlib import Path
import re
import subprocess
import sys

SIGNATURES = (
    ("Generated artifacts", r"is stale; run|generated.*out.of.date", "Run the generator and commit its output."),
    ("Go compiler", r"cannot use .+ as |cannot infer [A-Z]|undefined:|build failed|invalid case .+ in switch", "Build/test affected Go packages."),
    ("Lint and formatting", r"\((?:gocritic|goconst|errcheck|gofmt|goimports)\)|not properly formatted|unformatted Go", "Run just lint-new and gofmt."),
    ("Playwright timeout", r"playwright.*TimeoutError|Timeout \d+ms exceeded|waitForFunction: Timeout", "Inspect browser readiness, trace and screenshot."),
    ("Editor parity", r"AssertionError: Locator|Locator expected to be visible|element\(s\) not found|test_parity\.py", "Run focused editor browser parity."),
    ("Allocation threshold", r"allocation reduction|stream allocated .+ want below", "Measure allocation without race instrumentation."),
    ("Infrastructure", r"rate limit exceeded|HTTP 5\d\d|no space left on device|TLS handshake timeout", "Confirm infrastructure trouble before retry."),
)


def classify(log):
    out = []
    for name, pattern, hint in SIGNATURES:
        line = next((x for x in log.splitlines() if re.search(pattern, x, re.I)), None)
        if line:
            out.append((name, line.strip()[:240].replace(chr(96), "'"), hint))
    return out


def report(log, run_id="unknown"):
    output = [f"## Markata CI diagnostics: run {run_id}", "",
              "Read-only, heuristic analysis. Verify the first failed step.", ""]
    findings = classify(log)
    if not findings:
        output += ["No known signature found; inspect raw failed-job logs.", ""]
    for name, sample, hint in findings:
        output += [f"### {name}", f"- Evidence: {sample}", f"- Next: {hint}", ""]
    output.append("This report never retries jobs, overrides CI, or bypasses branch protection.")
    return "\n".join(output) + "\n"


def main():
    p = argparse.ArgumentParser(description=__doc__)
    src = p.add_mutually_exclusive_group(required=True)
    src.add_argument("--run")
    src.add_argument("--input")
    src.add_argument("--stdin", action="store_true")
    p.add_argument("--repo", default="WaylonWalker/markata")
    p.add_argument("--summary")
    args = p.parse_args()
    if args.run:
        result = subprocess.run(["gh", "run", "view", args.run, "-R", args.repo, "--log-failed"],
                                capture_output=True, text=True)
        if result.returncode and not result.stdout.strip():
            print(result.stderr, file=sys.stderr)
            return 1
        log = result.stdout
    elif args.input:
        log = Path(args.input).read_text(encoding="utf-8", errors="replace")
    else:
        log = sys.stdin.read()
    output = report(log, args.run or "saved log")
    print(output, end="")
    if args.summary:
        with Path(args.summary).open("a", encoding="utf-8") as stream:
            stream.write(output)
    return 0


if __name__ == "__main__":
    sys.exit(main())
