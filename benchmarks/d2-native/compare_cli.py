#!/usr/bin/env python3
"""Measure the newest pinned external D2 CLI against the native benchmark's fixtures.

CLI times intentionally include process startup (the build cost for a one-shot
renderer), whereas the native Go benchmark reuses an in-process ruler.
No third-party Python packages are required.
"""
from __future__ import annotations

import gzip
import json
import os
from pathlib import Path
import statistics
import subprocess
import sys
import tempfile
import time


def graph(nodes: int) -> str:
    return "direction: right\n" + "".join(
        f"node_{index} -> node_{index + 1}\n"
        for index in range(nodes - 1)
    )


def measure(binary: Path, nodes: int, folder: Path) -> dict[str, object]:
    source_path = folder / f"graph-{nodes}.d2"
    svg_path = folder / f"graph-{nodes}.svg"
    source_path.write_text(graph(nodes), encoding="utf-8")
    durations: list[float] = []
    svg = b""
    # One launch to warm the filesystem; timed runs still spawn a new process.
    subprocess.run(
        [str(binary), str(source_path), str(svg_path)],
        check=True, capture_output=True, timeout=90,
    )
    for _ in range(5):
        start = time.perf_counter()
        subprocess.run(
            [str(binary), str(source_path), str(svg_path)],
            check=True, capture_output=True, timeout=90,
        )
        durations.append((time.perf_counter() - start) * 1000)
        svg = svg_path.read_bytes()
    return {
        "nodes": nodes,
        "samples_ms": [round(value, 2) for value in durations],
        "median_ms": round(statistics.median(durations), 2),
        "svg_bytes": len(svg),
        "gzip_bytes": len(gzip.compress(svg, compresslevel=6)),
        "embedded_font": b"data:font/" in svg,
    }


def main() -> None:
    if len(sys.argv) != 2:
        raise SystemExit("usage: compare_cli.py /absolute/path/to/d2")
    binary = Path(sys.argv[1]).resolve(strict=True)
    version = subprocess.check_output([str(binary), "--version"], text=True).strip()
    with tempfile.TemporaryDirectory(prefix="markata-d2-cli-") as tmp:
        results = [measure(binary, n, Path(tmp)) for n in (3, 18, 72)]
    print("CLI version:", version)
    print("| Nodes | Median per-process render ms | Raw SVG bytes | gzip bytes | Embedded font |")
    print("| ---: | ---: | ---: | ---: | :---: |")
    for row in results:
        print(
            f'| {row["nodes"]} | {row["median_ms"]} | '
            f'{row["svg_bytes"]} | {row["gzip_bytes"]} | {row["embedded_font"]} |'
        )
    summary = os.environ.get("GITHUB_STEP_SUMMARY")
    if summary:
        with open(summary, "a", encoding="utf-8") as file:
            file.write("\n### External D2 CLI v0.9.0 (includes one process per render)\n")
            file.write(f"\nVersion: `{version}`\n\n")
            file.write("| Nodes | Median ms | Raw SVG bytes | gzip bytes | Embedded font |\n")
            file.write("| ---: | ---: | ---: | ---: | :---: |\n")
            for row in results:
                file.write(
                    f'| {row["nodes"]} | {row["median_ms"]} | '
                    f'{row["svg_bytes"]} | {row["gzip_bytes"]} | {row["embedded_font"]} |\n'
                )
    Path("cli-results.json").write_text(
        json.dumps({"version": version, "results": results}, indent=2) + "\n",
        encoding="utf-8",
    )


if __name__ == "__main__":
    main()
