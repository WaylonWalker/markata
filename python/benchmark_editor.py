"""Measure existing fixtures through the installed renderer and a real HTTP server."""
import hashlib
import json
import math
from datetime import datetime, timezone
from importlib.metadata import version
import platform
import statistics
import subprocess
import time
from pathlib import Path
from urllib.request import Request, urlopen

from markata import render, rendering
from editor_test_support import ROOT, load_editor, serve_editor

SAMPLES = 15
WARMUPS = 2
FIXTURES = [ROOT / 'benchmarks/site/posts/blog/2024/01/post-001.md', *sorted((ROOT / 'rendering-fixtures').glob('*.md'))]


def summary(values):
    return {'median': round(statistics.median(values), 3),
            'p95': round(sorted(values)[math.ceil(len(values) * .95) - 1], 3),
            'min': round(min(values), 3), 'max': round(max(values), 3)}


def request(url, markdown, document):
    payload = json.dumps({'markdown': markdown, 'include_document': document}).encode()
    started = time.perf_counter()
    with urlopen(Request(url + '/render', data=payload, headers={'Content-Type': 'application/json'}), timeout=30) as response:
        raw = response.read()
    elapsed = (time.perf_counter() - started) * 1000
    result = json.loads(raw)
    return elapsed, result['render_ms'], len(raw)


def main():
    report = {'revision': subprocess.check_output(['git', 'rev-parse', 'HEAD'], cwd=ROOT, text=True).strip(),
              'platform': platform.platform(), 'python': platform.python_version(),
              'recorded_at': datetime.now(timezone.utc).isoformat(),
              'binary_sha256': hashlib.sha256((Path(rendering.__file__).parent / 'bin/markata').read_bytes()).hexdigest(),
              'wrapper_sha256': hashlib.sha256(Path(rendering.__file__).read_bytes()).hexdigest(),
              'dependencies': {name: version(name) for name in ['fastapi', 'uvicorn']},
              'samples': SAMPLES, 'warmups_per_mode': WARMUPS, 'fixtures': []}
    app = load_editor().app
    with serve_editor(app) as url:
        for path in FIXTURES:
            source = path.read_text(encoding='utf-8')
            row = {'path': str(path.relative_to(ROOT)), 'sha256': hashlib.sha256(path.read_bytes()).hexdigest(),
                   'input_bytes': len(source.encode()), 'raw': {}}
            values = []
            for index in range(WARMUPS + SAMPLES):
                started = time.perf_counter()
                result = render(source, inline_assets=True, unsafe=True)
                elapsed = (time.perf_counter() - started) * 1000
                if index >= WARMUPS:
                    values.append(elapsed)
            row['direct_render_ms'] = summary(values)
            row['raw']['direct_render_ms'] = values
            row['features'] = [key for key, value in result.features.items() if value]
            for mode, document in [('initial', True), ('fragment', False)]:
                samples = [request(url, source, document) for _ in range(WARMUPS + SAMPLES)][WARMUPS:]
                row[mode] = {'http_ms': summary([v[0] for v in samples]),
                             'server_render_ms': summary([v[1] for v in samples]),
                             'response_bytes': summary([v[2] for v in samples])}
                row['raw'][mode] = samples
            report['fixtures'].append(row)
            print(path.name, row['direct_render_ms'], flush=True)
    output = ROOT / ('benchmarks/renderer-integration-' + datetime.now(timezone.utc).strftime('%Y-%m-%d') + '.json')
    output.write_text(json.dumps(report, indent=2) + '\n')
    print(output)


if __name__ == '__main__':
    main()
