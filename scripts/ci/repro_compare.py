#!/usr/bin/env python3
"""Compare generated site trees by content SHA-256, never by mtime.

This tool compares artifacts; it does not claim the inputs or build environment
were identical. Use it after two independent production builds of pinned inputs.
"""

from __future__ import annotations

import argparse
import fnmatch
import hashlib
import json
import sys
from pathlib import Path


def hash_file(path: Path) -> str:
    digest = hashlib.sha256()
    with path.open("rb") as file:
        for chunk in iter(lambda: file.read(1024 * 1024), b""):
            digest.update(chunk)
    return digest.hexdigest()


def inventory(root: Path, excludes: list[str]) -> dict[str, dict[str, object]]:
    if not root.is_dir() or root.is_symlink():
        raise ValueError(f"not a real output directory: {root}")
    files: dict[str, dict[str, object]] = {}
    for path in root.rglob("*"):
        relative = path.relative_to(root).as_posix()
        if path.is_symlink():
            raise ValueError(f"symlink not allowed in output: {relative}")
        if path.is_dir():
            continue
        if not path.is_file():
            raise ValueError(f"unexpected output entry: {relative}")
        if any(fnmatch.fnmatchcase(relative, pattern) for pattern in excludes):
            continue
        files[relative] = {"size": path.stat().st_size, "sha256": hash_file(path)}
    return dict(sorted(files.items()))


def compare_trees(
    left: Path, right: Path, excludes: list[str] | None = None
) -> dict[str, object]:
    exclusions = excludes or []
    previous = inventory(left, exclusions)
    current = inventory(right, exclusions)
    old_paths = set(previous)
    new_paths = set(current)
    changed = [
        {"path": name, "left": previous[name], "right": current[name]}
        for name in sorted(old_paths & new_paths)
        if previous[name] != current[name]
    ]
    return {
        "schema_version": 1,
        "identical": not (changed or old_paths != new_paths),
        "excluded_patterns": exclusions,
        "left_file_count": len(previous),
        "right_file_count": len(current),
        "only_left": sorted(old_paths - new_paths),
        "only_right": sorted(new_paths - old_paths),
        "changed": changed,
    }


def main(argv: list[str] | None = None) -> int:
    parser = argparse.ArgumentParser(
        description="Compare Markata output files by SHA-256; exit 0 identical, 1 different, 2 invalid."
    )
    parser.add_argument("left", type=Path, help="first generated output directory")
    parser.add_argument("right", type=Path, help="second generated output directory")
    parser.add_argument(
        "--exclude",
        action="append",
        default=[],
        metavar="GLOB",
        help="explicit relative-path glob to ignore (repeatable; no default exclusions)",
    )
    args = parser.parse_args(argv)
    try:
        result = compare_trees(args.left, args.right, args.exclude)
    except (OSError, ValueError) as error:
        print(f"repro_compare: {error}", file=sys.stderr)
        return 2
    print(json.dumps(result, sort_keys=True, indent=2))
    return 0 if result["identical"] else 1


if __name__ == "__main__":
    raise SystemExit(main())
