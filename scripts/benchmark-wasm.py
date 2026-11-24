"""Benchmark the served WASM worker. Run with uv run --with playwright."""
import argparse
import hashlib
import json
import math
from pathlib import Path
import platform
import statistics
from datetime import datetime, timezone
from playwright.sync_api import sync_playwright

ROOT = Path(__file__).resolve().parents[1]


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--url', default='http://127.0.0.1:8766/')
    parser.add_argument('--output', type=Path, required=True)
    args = parser.parse_args()
    fixtures = [ROOT / 'benchmarks/site/posts/blog/2024/01/post-001.md', *sorted((ROOT / 'rendering-fixtures').glob('*.md'))]
    report = {'recorded_at': datetime.now(timezone.utc).isoformat(), 'platform': platform.platform(),
              'timing': 'worker render_ms; excludes startup, transfer, browser layout and paint',
              'samples': 15, 'warmups': 2, 'fixtures': []}
    with sync_playwright() as driver:
        browser = driver.chromium.launch()
        report['browser'] = browser.version
        page = browser.new_page()
        page.goto(args.url)
        for fixture in fixtures:
            source = fixture.read_text()
            values = page.evaluate("""async source => {
                const worker = new Worker('worker.js');
                const next = () => new Promise((resolve, reject) => {
                    const timeout = setTimeout(() => reject(Error('worker timeout')), 60000);
                    worker.onmessage = event => {clearTimeout(timeout); resolve(event.data)};
                    worker.onerror = event => {clearTimeout(timeout); reject(Error(event.message))};
                });
                try {
                    await next();
                    const samples = [];
                    for (let revision=0; revision<17; revision++) {
                        const pending = next(); worker.postMessage({source, revision});
                        const response = await pending;
                        if (response.result.error) throw Error(response.result.error);
                        if (!response.result.html) throw Error('empty output');
                        if (revision>=2) samples.push(response.render_ms);
                    }
                    return samples;
                } finally {worker.terminate()}
            }""", source)
            if not all(isinstance(value, (int, float)) and value >= 0 for value in values):
                raise ValueError('invalid worker timings')
            report['fixtures'].append({'path': str(fixture.relative_to(ROOT)),
                'sha256': hashlib.sha256(fixture.read_bytes()).hexdigest(), 'render_ms': values,
                'median_ms': statistics.median(values), 'p95_ms': sorted(values)[math.ceil(len(values)*.95)-1]})
        browser.close()
    args.output.write_text(json.dumps(report, indent=2) + '\n')
    print(args.output)


if __name__ == '__main__':
    main()
