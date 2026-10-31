"""Benchmark Markata's per-call process vs the FastAPI-owned warm pool.

Run: PYTHONPATH=python python python/benchmark_renderer_pool.py
Requires a locally built python/markata/bin/markata binary.
"""
import asyncio
import json
import math
import os
import platform
import statistics
import time
from pathlib import Path
from markata.pool import RendererPool
from markata.rendering import render_async

FENCE = chr(96) * 3
FIXTURES = {
    "simple": "# Hello\n\nA short **paragraph** and [link](https://example.com).\n",
    "rich": "# Rich\n\n!!! note\n    Watch this.\n\n| A | B |\n|---|---|\n| 1 | 2 |\n\n" +
            FENCE + "python\nfor i in range(8):\n    print(i)\n" + FENCE + "\n",
    "mermaid": "# Mermaid\n\n" + FENCE + "mermaid\nflowchart LR\nA[Load] --> B[Render]\nB --> C[Serve]\n" + FENCE + "\n",
}

def summary(values, wall):
    sorted_values = sorted(values)
    return {
        "p50_ms": round(statistics.median(sorted_values), 2),
        "p95_ms": round(sorted_values[math.ceil(len(values) * 0.95)-1], 2),
        "requests_per_second": round(len(values) / wall, 2),
        "wall_s": round(wall, 3),
    }

async def sample(call, source, count, concurrency, warmups):
    limiter = asyncio.Semaphore(concurrency)
    async def one():
        async with limiter:
            start = time.perf_counter()
            result = await call(source, inline_assets=False)
            return ((time.perf_counter() - start) * 1000, result)
    await asyncio.gather(*(one() for _ in range(warmups)))
    start = time.perf_counter()
    data = await asyncio.gather(*(one() for _ in range(count)))
    wall = time.perf_counter() - start
    return summary([duration for duration, _ in data], wall)

async def main():
    count, warmups = 24, 3
    report = {"environment": {"platform": platform.platform(),
                              "python": platform.python_version(), "cpus": os.cpu_count()},
              "method": {"requests_per_case": count, "warmups": warmups,
                         "concurrency": [1, 4], "inline_assets": False},
              "fixtures": {}}
    for name, source in FIXTURES.items():
        direct = await render_async(source)
        async with RendererPool(workers=1) as pool:
            pooled = await pool.render(source)
        assert direct.features == pooled.features, f"{name}: features differ"
        if name != "mermaid":
            assert direct.html == pooled.html, f"{name}: HTML differs"
        case = {}
        for concurrency in (1, 4):
            baseline = await sample(render_async, source, count, concurrency, warmups)
            begin = time.perf_counter()
            async with RendererPool(workers=concurrency) as pool:
                startup_ms = round((time.perf_counter() - begin) * 1000, 2)
                warm = await sample(pool.render, source, count, concurrency, warmups)
            data = {"one_shot": baseline, "pool": warm,
                    "pool_startup_ms": startup_ms,
                    "latency_speedup": round(baseline["p50_ms"]/max(warm["p50_ms"], 0.001), 2),
                    "throughput_speedup": round(warm["requests_per_second"]/max(baseline["requests_per_second"], 0.001), 2)}
            case[str(concurrency)] = data
            print(f"BENCH {name} concurrency={concurrency} direct={baseline} warm={warm} "
                  f"latency_speedup={data['latency_speedup']}x", flush=True)
        report["fixtures"][name] = {"source_bytes":len(source.encode()),"cases":case}
    output = Path(os.environ.get("MARKATA_BENCHMARK_OUTPUT", "/tmp/markata-render-pool-bench.json"))
    output.write_text(json.dumps(report, indent=2)+"\n")
    print("BENCHMARK_REPORT " + str(output), flush=True)

if __name__ == "__main__":
    asyncio.run(main())
