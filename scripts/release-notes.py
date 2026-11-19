#!/usr/bin/env python3
"""Assemble Markdown release notes from reviewed PR fragments, without network access."""
import argparse
from pathlib import Path
import re

CATEGORIES = ("breaking", "feature", "fix", "performance", "other")
TITLES = {"breaking": "Breaking changes", "feature": "Features", "fix": "Fixes", "performance": "Performance", "other": "Other changes"}
PATTERN = re.compile(r"^([0-9]+)\.(breaking|feature|fix|performance|other)\.md$")

def build(root: Path, version: str) -> str:
    items = {name: [] for name in CATEGORIES}
    for path in sorted(root.glob("*.md")):
        match = PATTERN.fullmatch(path.name)
        if not match:
            raise ValueError(f"Invalid fragment filename: {path.name}")
        body = path.read_text(encoding="utf-8").strip()
        if not body or body.startswith("#"):
            raise ValueError(f"Expected plain-text release note: {path}")
        number, category = match.groups()
        items[category].append((int(number), body))
    lines = [f"# Markata {version}", ""]
    for category in CATEGORIES:
        if not items[category]:
            continue
        lines.extend([f"## {TITLES[category]}", ""])
        for number, body in sorted(items[category]):
            lines.append(f"- {body} (#{number})")
        lines.append("")
    if not any(items.values()):
        lines.extend(["No reviewed release-note fragments are available for this release.", ""])
    return "\n".join(lines)

if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--version", required=True)
    parser.add_argument("--root", type=Path, default=Path("changes/unreleased"))
    parser.add_argument("--output", type=Path)
    args = parser.parse_args()
    notes = build(args.root, args.version)
    if args.output:
        args.output.parent.mkdir(parents=True, exist_ok=True)
        args.output.write_text(notes, encoding="utf-8")
    else:
        print(notes, end="")
