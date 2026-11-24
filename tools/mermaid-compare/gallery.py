#!/usr/bin/env python3
"""Make a local side-by-side browser gallery from the two native SVG outputs."""
import html
import json
from pathlib import Path
import sys

root = Path(sys.argv[1])
report = json.loads((root / "comparison.json").read_text())
parts = ['<!doctype html><html lang="en"><meta charset="utf-8">',
         '<title>Markata native Mermaid comparison</title>',
         '<style>body{font:14px system-ui;max-width:1300px;margin:2rem auto;background:#fafafa;color:#172c25}',
         'header{position:sticky;top:0;background:#fafafa;padding:1rem}section{border-top:2px solid #d0d8d4;padding:1rem 0}',
         '.pair{display:grid;grid-template-columns:1fr 1fr;gap:1rem}',
         'figure{margin:0;background:white;border:1px solid #ccc;min-width:0}figcaption{padding:.6rem}',
         'iframe{width:100%;height:400px;border:0}code{white-space:pre-wrap;overflow-wrap:anywhere}',
         '</style><header><h1>Mermaid — native Go renderers</h1>',
         '<p>Independent SVG engines, NOT an exact comparison to Mermaid.js. ',
         'Geometry, labels, wrapping, palette and SVG safety need human review.</p></header>']
for i, (name, case) in enumerate(report["cases"].items()):
    parts.append(f'<section><h2>{html.escape(name)}</h2><div class="pair">')
    for backend in ("mmdg", "go-mermaid"):
        data = case[backend]
        uri = f"{i:03d}-{backend}.svg"
        if (root / uri).exists():
            img = f'<iframe sandbox src="{uri}" title="{html.escape(name)} {backend}"></iframe>'
        else:
            img = '<p>No SVG: ' + html.escape(str(data.get("error"))) + '</p>'
        ms = data.get("p50_ms")
        timing = f"{ms:.4f} ms" if ms is not None else "failed"
        parts.append(f'<figure><figcaption><b>{backend}</b> — {timing}, '+
                     f'{data.get("elements","?")} elements, {data.get("bytes","?")} bytes</figcaption>{img}</figure>')
    parts.append('</div></section>')
parts.append('</html>')
(root / "index.html").write_text("\n".join(parts), encoding="utf-8")
print("Gallery:", root / "index.html")
