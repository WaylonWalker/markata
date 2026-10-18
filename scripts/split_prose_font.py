#!/usr/bin/env -S uv run --quiet --script
# /// script
# requires-python = ">=3.12"
# dependencies = ["fonttools[woff]", "pyyaml"]
# ///
"""Regenerate disjoint prose subsets from an existing bundled core font."""
from __future__ import annotations

import argparse
import hashlib
from pathlib import Path

import yaml
from fontTools import subset
from fontTools.ttLib import TTFont


def codepoints(ranges: list[str]) -> set[int]:
    result: set[int] = set()
    for value in ranges:
        bounds = value.removeprefix("U+").split("-")
        result.update(range(int(bounds[0], 16), int(bounds[-1], 16) + 1))
    return result


def main() -> None:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("family", nargs="?", default="source-sans-3")
    parser.add_argument("--catalog-root", type=Path, default=Path("internal/fontcatalog"))
    args = parser.parse_args()
    root = args.catalog_root
    profiles = yaml.safe_load((root / "markata-fontpacks.yaml").read_text())["subset_profiles"]
    family = root / args.family
    manifest_path = family / "manifest.yaml"
    manifest = yaml.safe_load(manifest_path.read_text())
    source = family / manifest["tiers"]["prose-core"]["file"]
    original = TTFont(source)
    original_points = set(original.getBestCmap())
    ranges = {tier: profiles[tier]["unicode"] for tier in ("prose-basic", "prose-extra")}
    points = {tier: codepoints(value) for tier, value in ranges.items()}
    if points["prose-basic"] & points["prose-extra"]:
        raise ValueError("prose partitions must have disjoint Unicode ranges")
    if points["prose-basic"] | points["prose-extra"] != codepoints(profiles["prose-core"]["unicode"]):
        raise ValueError("partition ranges must preserve the prose-core repertoire")
    result_points: set[int] = set()
    original_axes = [(axis.axisTag, axis.minValue, axis.defaultValue, axis.maxValue) for axis in original["fvar"].axes] if "fvar" in original else []
    for tier in ranges:
        font = TTFont(source)
        options = subset.Options()
        options.layout_features = ["*"]
        splitter = subset.Subsetter(options=options)
        splitter.populate(unicodes=points[tier])
        splitter.subset(font)
        font.flavor = "woff2"
        font.recalcTimestamp = False
        output = family / f"{args.family}-{tier}.woff2"
        font.save(output)
        generated = TTFont(output)
        result_points.update(generated.getBestCmap())
        axes = [(axis.axisTag, axis.minValue, axis.defaultValue, axis.maxValue) for axis in generated["fvar"].axes] if "fvar" in generated else []
        if axes != original_axes:
            raise ValueError(f"{tier}: variable axes changed")
        manifest["tiers"][tier] = {
            "file": output.name, "profile": tier,
            "sha256": hashlib.sha256(output.read_bytes()).hexdigest(),
            "bytes": output.stat().st_size, "unicode_range": ranges[tier],
        }
        print(f"{tier}: {output.stat().st_size:,} bytes, {len(generated.getBestCmap())} codepoints")
    if result_points != original_points:
        raise ValueError("generated fonts did not preserve all original mapped glyphs")
    manifest["tiers"] = dict(sorted(manifest["tiers"].items()))
    manifest_path.write_text(yaml.safe_dump(manifest, sort_keys=False))


if __name__ == "__main__":
    main()
