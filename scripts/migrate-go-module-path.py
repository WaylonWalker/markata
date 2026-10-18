#!/usr/bin/env python3
"""One-time, reversible-in-Git Go module path migration.

Targets tracked Go sources plus the module declaration and GoReleaser's ldflags.
Leaves historical release assets, binary names, TOML configuration, Kubernetes
resource names, PyPI's original Python distribution and legacy GHCR aliases alone.

Run:
  python3 scripts/migrate-go-module-path.py --apply
  python3 scripts/migrate-go-module-path.py --check
"""
from __future__ import annotations

import argparse
import subprocess
import sys
from pathlib import Path

OLD = b"github.com/WaylonWalker/markata-go"
NEW = b"github.com/WaylonWalker/markata"
ROOT = Path(__file__).resolve().parent.parent


def tracked_files() -> list[Path]:
    raw = subprocess.check_output(
        ["git", "ls-files", "-z", "--", "*.go", "go.mod", ".goreleaser.yml"],
        cwd=ROOT,
    )
    return [ROOT / part.decode("utf-8") for part in raw.split(b"\0") if part]


def run(apply: bool) -> int:
    module_path = ROOT / "go.mod"
    module = module_path.read_bytes()
    old_decl = b"module " + OLD + b"\n"
    new_decl = b"module " + NEW + b"\n"
    if old_decl not in module and new_decl not in module:
        print("ERROR: unexpected Go module declaration; refusing to rewrite", file=sys.stderr)
        return 2

    changed: list[tuple[Path, int]] = []
    old_references = 0
    for path in tracked_files():
        data = path.read_bytes()
        count = data.count(OLD)
        if count:
            old_references += count
            changed.append((path, count))
            if apply:
                path.write_bytes(data.replace(OLD, NEW))

    if changed:
        label = "MIGRATED" if apply else "STALE"
        print(f"{label}: {len(changed)} tracked files, {old_references} occurrences")
        for path, count in changed[:35]:
            print(f"  {path.relative_to(ROOT)}: {count}")
        if len(changed) > 35:
            print(f"  ... and {len(changed) - 35} more files")
    else:
        print("PASS: tracked Go sources, go.mod, and GoReleaser use canonical module path")
    return 0 if apply or not changed else 1


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    mode = parser.add_mutually_exclusive_group(required=True)
    mode.add_argument("--apply", action="store_true", help="Rewrite tracked Go module references")
    mode.add_argument("--check", action="store_true", help="Fail if a tracked Go module reference remains old")
    args = parser.parse_args()
    return run(args.apply)


if __name__ == "__main__":
    sys.exit(main())
