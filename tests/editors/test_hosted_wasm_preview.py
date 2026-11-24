"""Hosted editor smoke test for opt-in WASM. Compile the WASM bundle before Go."""
import argparse
import re
import shutil
import subprocess
import tempfile
import threading
from pathlib import Path

from playwright.sync_api import sync_playwright, expect

ROOT = Path(__file__).resolve().parents[2]


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument('--go-binary', required=True)
    args = parser.parse_args()
    with tempfile.TemporaryDirectory(prefix='markata-wasm-editor-') as root:
        site = Path(root) / 'site'
        shutil.copytree(ROOT / 'examples/editor-site', site)
        process = subprocess.Popen(
            [args.go_binary, 'editor', '--listen', '127.0.0.1:0', '--workspace', str(site)],
            stderr=subprocess.PIPE, text=True,
        )
        watchdog = threading.Timer(60, process.terminate)
        watchdog.start()
        try:
            endpoint = None
            for line in process.stderr:
                endpoint = re.search(r'http://127\.0\.0\.1:\d+/', line)
                if endpoint:
                    break
            if not endpoint:
                raise AssertionError('Could not start hosted-capable Go editor')
        finally:
            watchdog.cancel()
        try:
            with sync_playwright() as driver:
                browser = driver.chromium.launch(headless=True)
                try:
                    context = browser.new_context()
                    page = context.new_page()
                    rendered_requests = []
                    errors = []
                    page.on('request', lambda request: rendered_requests.append(request)
                            if request.url.endswith('/render') else None)
                    page.on('pageerror', lambda error: errors.append(str(error)))
                    page.goto(endpoint.group(0) + '?wasm-preview=1')
                    assert page.request.get(endpoint.group(0) + 'editor/wasm/markata.wasm').status == 200
                    expect(page.frame_locator('#preview').locator('h1')).to_be_visible(timeout=60000)
                    page.wait_for_function("""() => {
                      const editor = document.querySelector('.CodeMirror')?.CodeMirror;
                      return !!editor && typeof ready !== 'undefined' && ready;
                    }""", timeout=60000)
                    page.evaluate("""() => document.querySelector('.CodeMirror').CodeMirror.setValue(
                      '# Native parity\\n\\nOriginal body.')""")
                    page.wait_for_function("""() => typeof wasmReady !== 'undefined' &&
                      wasmReady && !!wasmParity""", timeout=60000)
                    page.wait_for_function("""() => document.querySelector('#status').textContent
                      .includes('Preview')""", timeout=60000)
                    count = len(rendered_requests)
                    page.evaluate("""() => document.querySelector('.CodeMirror').CodeMirror.setValue(
                      '# Local WASM\\n\\nUpdated paragraph.')""")
                    expect(page.locator('#status')).to_contain_text('Quick local preview', timeout=60000)
                    expect(page.frame_locator('#preview').locator('h1')).to_have_text('Local WASM')
                    expect(page.frame_locator('#preview').locator('p')).to_contain_text('Updated paragraph.')
                    assert len(rendered_requests) == count, 'Local-only edit made a Go /render request'
                    page.evaluate("""() => {
                      const fence = String.fromCharCode(96).repeat(3);
                      document.querySelector('.CodeMirror').CodeMirror.setValue(
                        '# Native fallback\\n\\n' + fence + 'mermaid\\ngraph LR\\nA-->B\\n' + fence);
                    }""")
                    expect(page.frame_locator('#preview').locator('h1')).to_have_text('Native fallback', timeout=60000)
                    assert len(rendered_requests) > count, 'Feature Markdown did not fall back to Go'
                    assert not errors, errors
                    context.close()
                finally:
                    browser.close()
        finally:
            process.terminate()
            process.wait(timeout=15)


if __name__ == '__main__':
    main()
