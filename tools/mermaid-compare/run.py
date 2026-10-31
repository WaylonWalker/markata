#!/usr/bin/env python3
"""Native Mermaid renderer comparison. No Python dependencies beyond stdlib.

Build: (cd tools/mermaid-compare && go build -mod=mod -o /tmp/mermaid-compare .)
Run: python tools/mermaid-compare/run.py --binary /tmp/mermaid-compare --output /tmp/mermaid-report
Use --full for all Markdown fixtures; default limits real-site examples for quick PR CI.
The Go process remains warm; per-render timings exclude IPC, compilation and startup.
"""
import argparse
from collections import Counter
import html
import json
from pathlib import Path
import platform
import re
import statistics
import tomllib
import subprocess
import xml.etree.ElementTree as ET

ROOT = Path(__file__).resolve().parents[2]
CORE = {
    "flowchart": "flowchart LR\n  A[Load] --> B{Valid?}\n  B -->|yes| C[Render]\n  B -->|no| D[Error]",
    "sequence": "sequenceDiagram\n  Alice->>Bob: Hello\n  Bob-->>Alice: Hi",
    "class": "classDiagram\n  class Animal {\n    +String name\n    +eat()\n  }",
    "state": "stateDiagram-v2\n  [*] --> Active\n  Active --> [*]",
    "er": "erDiagram\n  CUSTOMER ||--o{ ORDER : places",
    "pie": 'pie title Pets\n  "Dogs" : 75\n  "Cats" : 25',
    "gantt": "gantt\n  title Delivery\n  dateFormat YYYY-MM-DD\n  section Work\n  Build :a1, 2026-01-01, 3d",
    "gitgraph": "gitGraph\n  commit\n  branch develop\n  commit\n  checkout main\n  merge develop",
    "mindmap": "mindmap\n  root((Idea))\n    Design\n    Test",
    "timeline": "timeline\n  title History\n  2020 : Started\n  2026 : Shipped",
    "journey": "journey\n  title Write\n  section Flow\n  Draft: 5: Author",
    "quadrant": "quadrantChart\n  x-axis Low --> High\n  y-axis Low --> High\n  point A: [0.3, 0.6]",
    "sankey": "sankey-beta\nSource,Process,10\nProcess,Output,10",
    "requirement": "requirementDiagram\n  requirement safety {\n    id: req1\n    text: user data remains safe\n    risk: high\n    verifymethod: test\n  }",
    "c4": "C4Context\n  Person(user, \"Reader\")\n  System(site, \"Markata\")\n  Rel(user, site, \"Visits\")",
    "xychart": "xychart-beta\n  x-axis [jan, feb, mar]\n  y-axis \"Readers\" 0 --> 12\n  bar [4, 8, 10]",
    "block": "block-beta\n  columns 2\n  A[Markdown] B[SVG]",
    "kanban": "kanban\n  column1[To Do]\n    task1[Write docs]\n  column2[Done]\n    task2[Tests]",
    "packet": "packet-beta\n  0-7: \"Header\"\n  8-15: \"Payload\"",
    "radar": "radar-beta\n  title Skill radar\n  axis speed[Speed]\n  axis quality[Quality]\n  curve output[Markata]{speed: 5, quality: 8}",
    "architecture": "architecture-beta\n  group api(cloud)[API]\n  service db(database)[DB] in api",
    "treemap": "treemap-beta\n  \"Project\"\n    \"Docs\": 3\n    \"Tests\": 4",
    "zenuml": "zenuml\n  Alice->Bob: request\n  Bob->Alice: response",
}
# Six core families must render with both engines; other families are measured, not
# silently dropped: unsupported syntax is recorded in the report.
REQUIRED = {"flowchart", "sequence", "class", "state", "er", "pie"}
# Labels visible in SVG, not just valid XML.
EXPECTED_LABELS = {
    "flowchart": [
        "Load",
        "Valid?",
        "Render"
    ],
    "sequence": [
        "Alice",
        "Bob",
        "Hello"
    ],
    "class": [
        "Animal",
        "name"
    ],
    "state": [
        "Active"
    ],
    "er": [
        "CUSTOMER",
        "ORDER",
        "places"
    ],
    "pie": [
        "Dogs",
        "Cats"
    ],
    "gantt": [
        "Build",
        "Delivery"
    ],
    "gitgraph": [
        "develop"
    ],
    "mindmap": [
        "Idea",
        "Design"
    ],
    "timeline": [
        "2020",
        "Started"
    ],
    "journey": [
        "Write",
        "Draft"
    ],
    "quadrant": [
        "Low",
        "High"
    ],
    "sankey": [
        "Source",
        "Process",
        "Output"
    ],
    "requirement": [
        "safety",
        "user data remains safe"
    ],
    "c4": [
        "Reader",
        "Markata"
    ],
    "xychart": [
        "jan",
        "feb",
        "mar"
    ],
    "block": [
        "Markdown",
        "SVG"
    ],
    "kanban": [
        "Write docs",
        "Tests"
    ],
    "packet": [
        "Header",
        "Payload"
    ],
    "radar": [
        "Skill radar",
        "Speed",
        "Quality"
    ],
    "architecture": [
        "DB",
        "API"
    ],
    "treemap": [
        "Project",
        "Docs",
        "Tests"
    ],
    "zenuml": [
        "Alice",
        "Bob"
    ]
}


def markata_palette(name):
    """Resolve colors from the repository's actual built-in semantic palette."""
    path = ROOT / "pkg" / "palettes" / "palettes" / (name + ".toml")
    data = tomllib.loads(path.read_text(encoding="utf-8"))["palette"]
    colors = data["colors"]
    semantic = data["semantic"]
    def color(role):
        token = semantic[role]
        return colors.get(token, token)
    return {"background": color("bg-primary"), "fill": color("bg-surface"),
            "stroke": color("accent"), "text": color("text-primary"),
            "edge": color("text-primary")}

def mermaid_blocks(path):
    try:
        lines = path.read_text(encoding="utf-8").splitlines()
    except (OSError, UnicodeError):
        return
    fence = None
    buf = []
    for line in lines:
        if fence:
            if line.strip().startswith(fence) and line.strip().strip(fence[0]).strip() == "":
                if buf:
                    yield "\n".join(buf).strip()
                fence, buf = None, []
            else:
                buf.append(line)
            continue
        m = re.match(r"^\s*(`{3,}|~{3,})\s*mermaid(?:\s+.*)?$", line, re.I)
        if m:
            fence = m.group(1)
            buf = []

def examples(limit):
    paths = []
    for dirname in ("docs", "examples", "rendering-fixtures", "benchmarks", "pages", "spec", "tests", "pkg/agentskill/bundle"):
        root = ROOT / dirname
        if root.exists():
            paths += sorted(p for p in root.rglob("*") if p.suffix.lower() in (".md", ".mdx", ".markdown")
                            and not any(x in p.parts for x in ("node_modules", "vendor", ".git")))
    seen = {v: "core:" + k for k, v in CORE.items()}
    result = [(k, src) for k, src in CORE.items()]
    total = 0
    for path in paths:
        for i, source in enumerate(mermaid_blocks(path), 1):
            total += 1
            if source in seen:
                continue
            if limit is not None and len(result) - len(CORE) >= limit:
                continue
            name = f"{path.relative_to(ROOT).as_posix()}:{i}"
            seen[source] = name
            result.append((name, source))
    # The existing Go plugin tests also contain real Mermaid examples as HTML.
    import html as htmlmod
    for file in sorted((ROOT / "pkg" / "plugins").glob("*mermaid*test.go")):
        raw_html = file.read_text(encoding="utf-8")
        blocks = re.findall(r'<pre><code class="language-mermaid"[^>]*>([\s\S]*?)</code></pre>', raw_html)
        for j, raw in enumerate(blocks, 1):
            source = htmlmod.unescape(raw).strip()
            total += 1
            if not source or source in seen:
                continue
            if limit is not None and len(result) - len(CORE) >= limit:
                continue
            name = f"{file.relative_to(ROOT).as_posix()}:go-test-{j}"
            seen[source] = name
            result.append((name, source))
    return result, total

def call(proc, source, iterations=8, **theme):
    req = {"source": source, "iterations": iterations, **theme}
    proc.stdin.write(json.dumps(req) + "\n")
    proc.stdin.flush()
    line = proc.stdout.readline()
    if not line:
        raise RuntimeError("native renderer process closed unexpectedly")
    payload = json.loads(line)
    if "results" not in payload:
        raise RuntimeError(f"native renderer protocol error: {payload}")
    return {x["backend"]: x for x in payload["results"]}

def analyze(svg):
    root = ET.fromstring(svg)
    if root.tag.split("}")[-1] != "svg":
        raise ValueError("root element is not svg")
    nodes = list(root.iter())
    forbidden = {"script", "iframe", "object", "embed"}
    if any(el.tag.split("}")[-1].lower() in forbidden for el in nodes):
        raise ValueError("active SVG content")
    for el in nodes:
        for key, value in el.attrib.items():
            if key.lower().startswith("on") or (
                key.split("}")[-1].lower() in ("href", "src")
                and re.match(r"(?is)^\s*(javascript|data):", value)
            ):
                raise ValueError(f"unsafe SVG attribute {key}")
    return {
        "bytes": len(svg.encode()),
        "elements": len(nodes),
        "tags": dict(Counter(el.tag.split("}")[-1] for el in nodes)),
        "text": " ".join((el.text or "").strip() for el in nodes if el.tag.split("}")[-1].lower() not in ("style", "script") and (el.text or "").strip())[:2000],
    }

def fmt(value):
    return f"{value:.4f}"

def perf_svg(report, path):
    # Two different timing boundaries must NEVER be plotted as one speedup.
    samples = []
    for name in CORE:
        case = report["cases"][name]
        for backend in ("mmdg", "go-mermaid"):
            x = case.get(backend, {})
            if x.get("p50_ms") is not None:
                samples.append((name, backend, x["p50_ms"]))
    rows = []
    rows.append('<svg xmlns="http://www.w3.org/2000/svg" width="1050" height="%d" viewBox="0 0 1050 %d">' %
                (195 + len(samples)*23, 195 + len(samples)*23))
    rows.append('<style>text{font:13px sans-serif;fill:#172b26}.head{font:bold 17px sans-serif}.barA{fill:#168264}.barB{fill:#64748b}</style>')
    rows.append('<rect width="100%" height="100%" fill="#fff"/>')
    rows.append('<text x="20" y="28" class="head">Markata Mermaid: previous warm-pool result</text>')
    rows.append('<text x="20" y="52">One-shot 230.33 ms    |    Warm pool 60.33 ms    |    3.8× improvement (prior session)</text>')
    rows.append('<text x="20" y="76" class="head">New native Go SVG-only timings — different operation, NOT end-to-end Markata</text>')
    rows.append('<text x="20" y="98">In-process warm render, excluding Go process startup, Markdown, Python, assets and layout in a browser</text>')
    largest=max((n for _,_,n in samples),default=1.0)
    for i,(name,backend,ms) in enumerate(samples):
        y=125+i*23
        label=html.escape(name[:22]+" / "+backend)
        rows.append(f'<text x="20" y="{y+11}">{label}</text>')
        rows.append(f'<rect x="250" y="{y}" height="15" width="{max(1,600*ms/max(largest,0.001)):.2f}" class="{"barA" if backend=="mmdg" else "barB"}"/>')
        rows.append(f'<text x="875" y="{y+11}">{ms:.4f} ms</text>')
    rows.append('</svg>')
    path.write_text("\n".join(rows), encoding="utf-8")

def main():
    parser=argparse.ArgumentParser()
    parser.add_argument("--binary", required=True)
    parser.add_argument("--output", type=Path, default=Path("/tmp/mermaid-comparison"))
    parser.add_argument("--full", action="store_true")
    parser.add_argument("--iterations", type=int, default=8)
    a=parser.parse_args()
    a.output.mkdir(parents=True,exist_ok=True)
    fixtures,total=examples(None if a.full else 32)
    report={"environment":{"platform":platform.platform(),"python":platform.python_version()},
            "method":{"iterations":a.iterations,"warmups":2,"boundary":"warm native Go source-to-SVG only"},
            "scanned_markdown_blocks":total,"evaluated":len(fixtures),"cases":{},"style_checks":{}}
    proc=subprocess.Popen([a.binary],stdin=subprocess.PIPE,stdout=subprocess.PIPE,text=True,bufsize=1)
    try:
        for i,(name,source) in enumerate(fixtures):
            results=call(proc,source,iterations=a.iterations)
            case={}
            for engine,row in results.items():
                item={"error":row.get("error")}
                if row.get("svg"):
                    try:
                        item.update(analyze(row["svg"]))
                        if name in EXPECTED_LABELS:
                            present = item["text"].casefold()
                            item["missing_labels"] = [label for label in EXPECTED_LABELS[name]
                                                      if label.casefold() not in present]
                        if not row.get("error"):
                            item["p50_ms"]=statistics.median(row["ns"])/1e6
                        (a.output / f"{i:03d}-{engine}.svg").write_text(row["svg"],encoding="utf-8")
                    except (ET.ParseError,ValueError) as exc:
                        item["error"]=f"invalid SVG: {exc}"
                case[engine]=item
            report["cases"][name]=case
        # A real contrast test: the same graph must respond to both font and palette
        # changes. A CSS family declaration alone is NOT proof of font metrics.
        src="flowchart LR\nA[Typography and rhythm] --> B[Express an idea quickly]"
        base=call(proc,src,iterations=1)
        styled=call(proc,src,iterations=1,font="DejaVu Sans",font_size=21,
                    background="#091c18",fill="#173b32",stroke="#65e1b1",text="#f4fff9",edge="#65e1b1")
        for backend in ("mmdg","go-mermaid"):
            old,new=base[backend],styled[backend]
            newsvg=new.get("svg","")
            report["style_checks"][backend]={
                "error":new.get("error"),
                "output_changed":bool(newsvg and newsvg!=old.get("svg")),
                "font_family_in_svg":"dejavu sans" in newsvg.lower(),
                "text_color_in_svg":"#f4fff9" in newsvg.lower(),
                "edge_color_in_svg":"#65e1b1" in newsvg.lower(),
                "font_size_in_svg":bool(re.search(r"(?i)(font-size[^;>]*21|font-size[^;>]*:21)",newsvg)),
            }
            if newsvg:
                (a.output / f"theme-font-{backend}.svg").write_text(newsvg,encoding="utf-8")
        # Genuine Markata palette source files, not hand-picked example colors.
        report["palette_checks"] = {}
        palette_src = "flowchart LR\n A[Readable contrast] --> B{Themed diagram}"
        for palette_name in ("catppuccin-mocha", "catppuccin-latte", "everforest-dark", "everforest-light"):
            pal = markata_palette(palette_name)
            outputs = call(proc, palette_src, iterations=1, **pal)
            report["palette_checks"][palette_name] = {}
            for backend, row in outputs.items():
                svg = row.get("svg", "")
                report["palette_checks"][palette_name][backend] = {
                    "error": row.get("error"),
                    "background_present": pal["background"].lower() in svg.lower(),
                    "text_present": pal["text"].lower() in svg.lower(),
                    "accent_present": pal["stroke"].lower() in svg.lower(),
                    "visible_labels": all(x in svg for x in ("Readable contrast", "Themed diagram")),
                }
                if svg:
                    (a.output / f"palette-{palette_name}-{backend}.svg").write_text(svg, encoding="utf-8")
    finally:
        proc.stdin.close()
        proc.stdout.close()
        proc.wait(timeout=15)
    failures=[f"{name}:{backend} {report['cases'][name][backend].get('error')}"
              for name in REQUIRED for backend in ("mmdg","go-mermaid")
              if report["cases"][name][backend].get("error")]
    for name in REQUIRED:
        for backend in ("mmdg", "go-mermaid"):
            missing = report["cases"][name][backend].get("missing_labels", [])
            if missing:
                failures.append(f"{name}:{backend} missing visible labels {missing}")
    report["semantic_gaps"] = {
        backend: {name: item[backend]["missing_labels"]
                  for name, item in report["cases"].items() if item[backend].get("missing_labels")}
        for backend in ("mmdg", "go-mermaid")
    }
    for palette_name, engines in report["palette_checks"].items():
        for backend, check in engines.items():
            if check.get("error") or not all(check.get(x) for x in (
                "background_present", "text_present", "accent_present", "visible_labels"
            )):
                failures.append(f"{palette_name}:{backend} palette/label check failed: {check}")
    report["required_failures"]=failures
    report["coverage"]={backend:sum(not r[backend].get("error") for r in report["cases"].values())
                       for backend in ("mmdg","go-mermaid")}
    (a.output/"comparison.json").write_text(json.dumps(report,indent=2)+"\n")
    perf_svg(report,a.output/"performance.svg")
    lines=["# Native Mermaid comparison (experimental)","",
           f"Mermaid source blocks found: {total}; unique fixtures evaluated: {len(fixtures)}.",
           "Timings are **SVG-only**, not full Markata latency. Do not compare speedups directly to the pool benchmark.","",
           "| Diagram | mmdg p50 ms | go-mermaid p50 ms | Differences |","|---|---:|---:|---|"]
    for name,case in report["cases"].items():
        left,right=case["mmdg"],case["go-mermaid"]
        l=f"{left['p50_ms']:.4f}" if "p50_ms" in left else "ERROR"
        r=f"{right['p50_ms']:.4f}" if "p50_ms" in right else "ERROR"
        difference=(f"elements {left.get('elements','?')} vs {right.get('elements','?')}; "+
                    f"bytes {left.get('bytes','?')} vs {right.get('bytes','?')}")
        lines.append(f"| {name.replace('|','/')} | {l} | {r} | {difference} |")
    lines+=["","## Theme and font changes","", "Output_changed means a difference, not correct measurements or no clipping.","",
            "```json",json.dumps(report["style_checks"],indent=2),"```",
            "","## Required failures","",*(failures or ["None"])]
    lines.extend(["", "## Missing semantic content", "",
                  "Visible label gaps are failures of fidelity even when SVG is valid.",
                  "```json", json.dumps(report["semantic_gaps"], indent=2), "```"])
    lines.extend(["", "## Real Markata built-in palettes", "",
                  "Four real semantic palette TOMLs (light and dark) were rendered by both engines.",
                  "Passing CSS color presence does not validate custom WOFF2 font metrics.",
                  "```json", json.dumps(report["palette_checks"], indent=2), "```"])
    (a.output/"comparison.md").write_text("\n".join(lines)+"\n")
    print(f"REPORT {a.output}/comparison.json")
    print(f"COVERAGE {report['coverage']} / {len(fixtures)}")
    print(f"REQUIRED_FAILURES {failures}")
    print(f"SEMANTIC_GAPS {json.dumps(report['semantic_gaps'])}")
    print(f"MARKATA_PALETTES {json.dumps(report['palette_checks'])}")
    print(f"THEME_FONT {json.dumps(report['style_checks'])}")
    # Optional full-suite compatibility gaps are reported, rather than making
    # every unrelated site fixture block the experimental research PR.
    if failures:
        raise SystemExit(1)

if __name__=="__main__":
    main()
