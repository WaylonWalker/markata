#!/usr/bin/env python3
"""Select a complete, deterministic top-level Go site-test shard.

Site integration helpers change the process working directory with os.Chdir,
so t.Parallel() inside that package is unsafe. Separate go test processes
have independent working directories and may run concurrently on CI runners.
"""

from __future__ import annotations

import argparse
from pathlib import Path
import re
import subprocess
import sys

# go test -list . includes runnable top-level tests, examples and fuzz seeds;
# benchmarks are not part of ordinary go test execution.
RUNNABLE_NAME = re.compile(r"(?:Test|Example|Fuzz)[A-Za-z0-9_]*\Z")


def discover_names(output: str) -> list[str]:
    names = [
        line.strip()
        for line in output.splitlines()
        if RUNNABLE_NAME.fullmatch(line.strip())
    ]
    if not names:
        raise ValueError("go test -list returned no runnable site tests")
    if len(names) != len(set(names)):
        raise ValueError("go test -list returned duplicate test names")
    return sorted(names)


def select_shard(names: list[str], index: int, count: int) -> list[str]:
    if count < 1 or index < 0 or index >= count:
        raise ValueError("shard index must be within [0, count), with count >= 1")
    if len(names) < count:
        raise ValueError("fewer runnable tests than requested shards")
    selected = [name for position, name in enumerate(sorted(names)) if position % count == index]
    if not selected:
        raise ValueError(f"site test shard {index}/{count} has no tests")
    return selected


def run_regex(names: list[str]) -> str:
    if not names:
        raise ValueError("cannot make a test regex from an empty shard")
    return "^(" + "|".join(re.escape(name) for name in names) + ")$"


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--index", required=True, type=int)
    parser.add_argument("--count", required=True, type=int)
    parser.add_argument("--output", required=True, type=Path)
    args = parser.parse_args()

    try:
        listed = subprocess.run(
            ["go", "test", "-list", ".", "./tests"],
            check=True,
            text=True,
            capture_output=True,
            timeout=300,
        )
        names = discover_names(listed.stdout)
        selected = select_shard(names, args.index, args.count)
    except (OSError, ValueError, subprocess.CalledProcessError, subprocess.TimeoutExpired) as error:
        print(f"site test shard selection failed: {error}", file=sys.stderr)
        if isinstance(error, subprocess.CalledProcessError):
            print(error.stderr or error.stdout, file=sys.stderr)
        return 1

    args.output.write_text(run_regex(selected) + "\n", encoding="utf-8")
    print(f"Site shard {args.index + 1}/{args.count}: selected {len(selected)}/{len(names)} runnable tests")
    print("Tests: " + ", ".join(selected))
    return 0


if __name__ == "__main__":
    sys.exit(main())
