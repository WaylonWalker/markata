#!/usr/bin/env python3
"""Compute matching SemVer and PEP 440 nightly versions, always within 0.x."""
from __future__ import annotations

import argparse
from datetime import datetime, timezone
from pathlib import Path
import re
import subprocess

STABLE_TAG = re.compile(r"^v0\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)$")


def calculate(tags: list[str], date: str) -> tuple[str, str, str]:
    """Ignore prerelease tags; select the greatest stable 0.x SemVer tag."""
    datetime.strptime(date, "%Y%m%d")
    stable = [tuple(map(int, match.groups()))
              for tag in tags if (match := STABLE_TAG.fullmatch(tag.strip()))]
    if not stable:
        raise ValueError("No stable v0.x.y tags; refusing to guess a release base")
    minor, patch = max(stable)
    base = f"0.{minor}.{patch + 1}"
    return (f"v{base}-dev.{date}", f"{base}.dev{date}", f"{base}-dev.{date}")


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument("--github-output", type=Path)
    parser.add_argument("--date", help="UTC date YYYYMMDD (tests only)")
    args = parser.parse_args()
    date = args.date or datetime.now(timezone.utc).strftime("%Y%m%d")
    tags = subprocess.check_output(["git", "tag", "--list", "v0.*"], text=True).splitlines()
    github_tag, pypi_version, goreleaser_version = calculate(tags, date)
    lines = (f"github_tag={github_tag}\n"
             f"pypi_version={pypi_version}\n"
             f"goreleaser_version={goreleaser_version}\n")
    if args.github_output:
        with args.github_output.open("a", encoding="utf-8") as handle:
            handle.write(lines)
    else:
        print(lines, end="")


if __name__ == "__main__":
    main()
