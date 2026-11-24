#!/usr/bin/env python3
"""Generate a concurrency-1 latency bar chart from actual CI benchmark JSON.

Two timing boundaries are charted as separate panels and independent axes:
1. Full Markata Python request -> Go renderer response (client-mode Mermaid).
2. Native Go source -> SVG generation, performed in the persistent candidate process.

The native benchmark uses the EXACT Mermaid diagram from benchmark_renderer_pool.py,
but this is NOT a measurement of a Markata request with native Mermaid enabled.
"""
import argparse
from html import escape
import json
from pathlib import Path

def rect(x, y, w, h, cls):
    return f'<rect x="{x:.2f}" y="{y:.2f}" width="{w:.2f}" height="{h:.2f}" rx="5" class="{cls}"/>'

def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("--pool-report", type=Path, required=True)
    ap.add_argument("--native-report", type=Path, required=True)
    ap.add_argument("--output-dir", type=Path, required=True)
    args = ap.parse_args()
    args.output_dir.mkdir(parents=True, exist_ok=True)
    pool = json.loads(args.pool_report.read_text())
    native = json.loads(args.native_report.read_text())
    source = pool["fixtures"]["mermaid"]["cases"]["1"]
    matched = native["cases"].get("pool_mermaid")
    if matched is None:
        raise SystemExit("pool_mermaid missing: cannot compare an identical Mermaid example")
    native_values = {}
    for engine in ("go-mermaid", "mmdg"):
        result = matched[engine]
        if result.get("error") or result.get("p50_ms") is None:
            raise SystemExit(f"native {engine} has no successful benchmark: {result}")
        native_values[engine] = float(result["p50_ms"])
    full_values = {
        "One-shot Markata": float(source["one_shot"]["p50_ms"]),
        "Warm-pool Markata": float(source["pool"]["p50_ms"]),
    }
    if "pool_native" not in source:
        raise SystemExit("native warm-pool measurement missing; cannot claim end-to-end improvement")
    full_values["Warm-pool Markata + go-mermaid"] = float(source["pool_native"]["p50_ms"])
    # Two scales! Comparing full Markdown+asset request latency to SVG-only cost
    # on a shared numerical axis would imply an end-to-end speedup never measured.
    sections = [
        ("Full Markata request — client vs native go-mermaid", full_values,
         "Includes Python/stdio/Go Markdown rendering and asset handling",
         "full"),
        ("Native Mermaid SVG generation — standalone candidate Go library", native_values,
         "Excludes Markata Markdown, assets, JSONL, Python, browser paint",
         "native"),
    ]
    W,H = 1100,530
    svg = [
        f'<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 {W} {H}" width="{W}" height="{H}" role="img">',
        "<title>Median render latency, concurrency one: Markata current pipeline versus native Mermaid SVG-only candidates</title>",
        "<desc>Two separate panels use different axes. The native library timings do not show end-to-end Markata render speed.</desc>",
        '<style>text{font-family:system-ui,-apple-system,sans-serif;fill:#16332d;}'
        '.title{font-size:25px;font-weight:750}.head{font-size:17px;font-weight:700}'
        '.sub{font-size:12px;fill:#546660}.label{font-size:15px}'
        '.number{font-size:16px;font-weight:750}.bg{fill:#f7faf8}'
        '.full{fill:#6d7f87}.native{fill:#278b6a}'
        '.track{fill:#dce8e1}.rule{stroke:#ccd8d0;stroke-width:1}'
        '</style>',
        rect(0,0,W,H,"bg"),
        '<text x="28" y="39" class="title">Median render latency · concurrency 1</text>',
        '<text x="28" y="61" class="sub">Same Mermaid source: flowchart LR · Load → Render → Serve · measured in CI</text>',
    ]
    for section_index,(heading, values,desc,cls) in enumerate(sections):
        top = 94 + section_index*210
        svg += [
            f'<text x="28" y="{top}" class="head">{escape(heading)}</text>',
            f'<text x="28" y="{top+19}" class="sub">{escape(desc)}</text>',
        ]
        max_v = max(values.values())
        for i,(name,ms) in enumerate(values.items()):
            y = top+38+i*49
            bar_width=max(1.5,600*ms/max_v)
            svg += [
                f'<text x="28" y="{y+18}" class="label">{escape("Native warm pool" if name == "Warm-pool Markata + go-mermaid" else name)}</text>',
                rect(240,y,600,25,"track"),
                rect(240,y,bar_width,25,cls),
                f'<text x="855" y="{y+19}" class="number">{ms:.5f} ms</text>',
            ]
        if section_index == 0:
            svg.append('<line x1="28" x2="1068" y1="282" y2="282" class="rule"/>')
    svg += [
        '<text x="28" y="493" class="sub">Markata warm-pool + go-mermaid is measured end-to-end. SVG-only bars have a separate time boundary.</text>',
        '<text x="28" y="514" class="sub">Both panels use independent scales. Warm-pool statistics: 24 requests, 3 warmups; native: 3 measured, 2 warmups.</text>',
        '</svg>',
    ]
    (args.output_dir / "median-render-latency-c1.svg").write_text("\n".join(svg)+"\n", encoding="utf-8")
    rows=[
        "# Median render latency — concurrency 1", "",
        "The Mermaid source is identical in both test paths: `flowchart LR / Load -> Render -> Serve`.",
        "First three bars share a full Python warm-pool request boundary; the last two are SVG-only.", "",
        "| Path | p50 ms | Timing boundary |",
        "|---|---:|---|",
        f"| Markata one-shot | {full_values['One-shot Markata']:.3f} | Full Python request, Go renderer, existing client-mode Mermaid |",
        f"| Markata warm pool | {full_values['Warm-pool Markata']:.3f} | Full Python request via warmed Go process, existing client-mode Mermaid |",
        f"| Markata warm pool + go-mermaid | {full_values['Warm-pool Markata + go-mermaid']:.3f} | Full Python request via warmed Go process with native SVG |",
        f"| Native go-mermaid | {native_values['go-mermaid']:.5f} | Warm Go source-to-SVG only |",
        f"| Native mmdg | {native_values['mmdg']:.5f} | Warm Go source-to-SVG only |",
        "",
        "Markata warm-pool + go-mermaid is a real full-request before/after comparison. Native SVG-only bars have an independent axis.",
    ]
    (args.output_dir / "median-render-latency-c1.md").write_text("\n".join(rows)+"\n")
    print("CONCURRENCY_1_MEDIANS",json.dumps({"current":full_values,"native":native_values}))
    print("CHART",args.output_dir / "median-render-latency-c1.svg")

if __name__=="__main__":
    main()
