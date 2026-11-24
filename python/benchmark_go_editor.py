"""Measure the Go workbench on the Python reference's exact fixture files."""
import argparse
from contextlib import contextmanager
from datetime import datetime, timezone
import hashlib
import json
import math
import os
from pathlib import Path
import platform
import queue
import re
import statistics
import subprocess
import sys
import threading
import time
from urllib.request import Request, urlopen

ROOT = Path(__file__).resolve().parents[1]
FIXTURES = [ROOT / "benchmarks/site/posts/blog/2024/01/post-001.md",
            *sorted((ROOT / "rendering-fixtures").glob("*.md"))]
INVALIDATOR = json.dumps({"markdown": "Go benchmark cache invalidation", "include_document": False}).encode()


def summary(values):
    return {"median": round(statistics.median(values), 3),
            "p95": round(sorted(values)[math.ceil(len(values) * .95) - 1], 3),
            "min": round(min(values), 3), "max": round(max(values), 3)}


@contextmanager
def serve(binary):
    """Start a dedicated loopback server and reap it even after a failed sample."""
    process = subprocess.Popen([str(binary), "--listen", "127.0.0.1:0"],
                               stdout=subprocess.DEVNULL, stderr=subprocess.PIPE, text=True)
    addresses = queue.Queue()

    def read_log():
        for line in process.stderr:
            if match := re.search(r"Markata Go workbench: (http://127\.0\.0\.1:\d+)/", line):
                addresses.put(match[1])
        addresses.put(None)

    thread = threading.Thread(target=read_log, daemon=True)
    thread.start()
    try:
        try:
            url = addresses.get(timeout=30)
        except queue.Empty as error:
            raise RuntimeError("Go workbench did not report a listen address within 30 seconds") from error
        if url is None:
            raise RuntimeError("Go workbench exited before reporting a listen address")
        yield url
    finally:
        process.terminate()
        try:
            process.wait(timeout=10)
        except subprocess.TimeoutExpired:
            process.kill()
            process.wait()
        thread.join(timeout=5)
        process.stderr.close()


def request(url, payload):
    # Match benchmark_editor.py: JSON encode/decode is outside the HTTP timer.
    started = time.perf_counter()
    with urlopen(Request(url + "/render", data=payload, headers={"Content-Type": "application/json"}), timeout=30) as response:
        raw = response.read()
    elapsed = (time.perf_counter() - started) * 1000
    result = json.loads(raw)
    return elapsed, result, len(raw)


def measure(url, source, document, cached, warmups, samples):
    payload = json.dumps({"markdown": source, "include_document": document}).encode()
    if cached:
        request(url, payload)
    records = []
    for index in range(warmups + samples):
        if not cached:
            # Keep the measured fixture byte-for-byte identical to the reference.
            # This separate request is excluded from both timing and response size.
            request(url, INVALIDATOR)
        elapsed, result, size = request(url, payload)
        if result["cached"] != cached:
            raise RuntimeError(f"Unexpected cache state: expected {cached}, got {result['cached']}")
        if (result["document"] is not None) != document or result["input_bytes"] != len(source.encode()):
            raise RuntimeError("Unexpected document mode or input length")
        if index >= warmups:
            records.append({"http_ms": elapsed, "server_render_ms": result["render_ms"],
                            "response_bytes": size, "html_bytes": result["html_bytes"], "cached": result["cached"]})
    return {key: summary([record[key] for record in records])
            for key in ["http_ms", "server_render_ms", "response_bytes"]}, records


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--binary", type=Path, required=True, help="compiled markata-workbench binary")
    parser.add_argument("--output", type=Path, default=ROOT / f"benchmarks/go-workbench-{datetime.now(timezone.utc):%Y-%m-%d}.json")
    parser.add_argument("--samples", type=int, default=15, help="measured calls per mode (default: 15)")
    parser.add_argument("--warmups", type=int, default=2, help="excluded calls per mode (default: 2)")
    args = parser.parse_args()
    if args.samples < 1 or args.warmups < 0:
        parser.error("samples must be positive and warmups must be non-negative")
    binary = args.binary.resolve(strict=True)
    report = {"revision": subprocess.check_output(["git", "rev-parse", "HEAD"], cwd=ROOT, text=True).strip(),
              "recorded_at": datetime.now(timezone.utc).isoformat(), "platform": platform.platform(),
              "python": platform.python_version(), "logical_cpus": os.cpu_count(),
              "go": subprocess.check_output(["go", "version", str(binary)], text=True).strip(),
              "binary_sha256": hashlib.sha256(binary.read_bytes()).hexdigest(),
              "harness_sha256": hashlib.sha256(Path(__file__).read_bytes()).hexdigest(),
              "samples": args.samples, "warmups_per_mode": args.warmups,
              "asset_policy": "linked local dependencies; separately fetched assets excluded",
              "miss_policy": "separate unmeasured request invalidates the one-entry cache before each miss",
              "fixtures": []}
    with serve(binary) as url:
        for path in FIXTURES:
            data = path.read_bytes()
            source = data.decode("utf-8")
            row = {"path": str(path.relative_to(ROOT)), "sha256": hashlib.sha256(data).hexdigest(),
                   "input_bytes": len(data), "raw": {}}
            for mode, document, cached in [("initial_miss", True, False),
                                           ("fragment_miss", False, False),
                                           ("fragment_cached", False, True)]:
                row[mode], row["raw"][mode] = measure(url, source, document, cached, args.warmups, args.samples)
            report["fixtures"].append(row)
            print(f"{path.name}: edit {row['fragment_miss']['http_ms']['median']:.3f} ms; "
                  f"cached {row['fragment_cached']['http_ms']['median']:.3f} ms", file=sys.stderr, flush=True)
    args.output.parent.mkdir(parents=True, exist_ok=True)
    args.output.write_text(json.dumps(report, indent=2) + "\n")
    print(args.output)


if __name__ == "__main__":
    try:
        main()
    except (OSError, RuntimeError, ValueError, subprocess.SubprocessError) as error:
        raise SystemExit(f"error: {error}") from error
