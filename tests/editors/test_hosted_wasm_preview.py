"""Exercise the authenticated hosted WASM path through a local proxy."""
import argparse
import json
import re
import shutil
import subprocess
import tempfile
import threading
import urllib.error
import urllib.request
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer
from pathlib import Path

from playwright.sync_api import sync_playwright, expect

ROOT = Path(__file__).resolve().parents[2]
PUBLIC_HOST = 'editor.example.test'


def assert_unauthorized_asset(endpoint):
    request = urllib.request.Request(
        endpoint + 'editor/wasm/markata.wasm',
        headers={'Host': PUBLIC_HOST},
    )
    try:
        urllib.request.urlopen(request, timeout=10)
    except urllib.error.HTTPError as error:
        assert error.code == 401, error.code
        assert error.headers.get('Cache-Control') == 'private, no-store'
        return
    raise AssertionError('unauthenticated WASM asset request was accepted')


def authenticated_proxy(upstream):
    class ForwardAuthProxy(BaseHTTPRequestHandler):
        protocol_version = 'HTTP/1.1'

        def log_message(self, *_args):
            pass

        def forward(self):
            if self.headers.get('Upgrade', '').lower() == 'websocket':
                self.send_error(501, 'WebSocket forwarding is outside this preview smoke test')
                return
            body = self.rfile.read(int(self.headers.get('Content-Length', '0')))
            headers = {
                key: value for key, value in self.headers.items()
                if key.lower() not in {'host', 'connection', 'content-length', 'transfer-encoding', 'upgrade'}
            }
            headers.update({
                'Host': PUBLIC_HOST,
                'Origin': 'https://' + PUBLIC_HOST,
                'X-Hlab-User-Id': 'usr_test_editor',
                'X-Hlab-Scopes': 'markata-editor:admin',
            })
            request = urllib.request.Request(
                upstream.rstrip('/') + self.path,
                data=body if body else None,
                headers=headers,
                method=self.command,
            )
            try:
                response = urllib.request.urlopen(request, timeout=30)
            except urllib.error.HTTPError as error:
                response = error
            with response:
                payload = response.read()
                self.send_response(response.status)
                for key, value in response.headers.items():
                    if key.lower() not in {'connection', 'content-length', 'transfer-encoding', 'keep-alive'}:
                        self.send_header(key, value)
                content_length = response.headers.get('Content-Length') if self.command == 'HEAD' else None
                self.send_header('Content-Length', content_length or str(len(payload)))
                self.end_headers()
                if payload:
                    try:
                        self.wfile.write(payload)
                    except (BrokenPipeError, ConnectionResetError):
                        pass

        do_GET = forward
        do_HEAD = forward
        do_POST = forward
        do_PUT = forward
        do_DELETE = forward

    server = ThreadingHTTPServer(('127.0.0.1', 0), ForwardAuthProxy)
    thread = threading.Thread(target=server.serve_forever, daemon=True)
    thread.start()
    return server


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument('--go-binary', required=True)
    args = parser.parse_args()
    with tempfile.TemporaryDirectory(prefix='markata-wasm-editor-') as root:
        site = Path(root) / 'site'
        shutil.copytree(ROOT / 'examples/editor-site', site)
        # This smoke test covers the renderer-compatible no-site-config path.
        (site / 'markata-go.toml').unlink(missing_ok=True)
        process = subprocess.Popen(
            [
                args.go_binary, 'editor', '--listen', '127.0.0.1:0', '--workspace', str(site),
                '--public-origin', 'https://' + PUBLIC_HOST,
                '--trusted-proxy-cidr', '127.0.0.1/32',
                '--required-scope', 'markata-editor:admin',
            ],
            stderr=subprocess.PIPE, text=True,
        )
        watchdog = threading.Timer(90, process.terminate)
        watchdog.start()
        proxy = None
        try:
            endpoint = None
            for line in process.stderr:
                endpoint = re.search(r'http://127\.0\.0\.1:\d+/', line)
                if endpoint:
                    break
            if not endpoint:
                raise AssertionError('Could not start hosted-capable Go editor')
            upstream = endpoint.group(0)
            assert_unauthorized_asset(upstream)
            proxy = authenticated_proxy(upstream)
            browser_endpoint = f'http://127.0.0.1:{proxy.server_port}/'

            with sync_playwright() as driver:
                browser = driver.chromium.launch(headless=True)
                try:
                    context = browser.new_context()
                    # Headless Chromium omits deviceMemory and reports the
                    # container CPU limit. Simulate a capable device so this
                    # test exercises the worker path.
                    context.add_init_script("""Object.defineProperty(navigator, 'deviceMemory', {get: () => 8});
                      Object.defineProperty(navigator, 'hardwareConcurrency', {get: () => 8});""")
                    page = context.new_page()
                    rendered_requests = []
                    errors = []
                    console_messages = []
                    failed_requests = []
                    page.on('request', lambda request: rendered_requests.append(request)
                            if request.url.endswith('/render') else None)
                    page.on('pageerror', lambda error: errors.append(str(error)))
                    page.on('console', lambda message: console_messages.append(message.text))
                    page.on('requestfailed', lambda request: failed_requests.append(
                        f'{request.method} {request.url}: {request.failure}'))
                    page.goto(browser_endpoint + '?wasm-preview=1')
                    expect(page.frame_locator('#preview').locator('h1')).to_be_visible(timeout=60000)
                    page.wait_for_function("""() => {
                      const editor = document.querySelector('.CodeMirror')?.CodeMirror;
                      return !!editor && typeof ready !== 'undefined' && ready;
                    }""", timeout=60000)
                    page.evaluate("""() => document.querySelector('.CodeMirror').CodeMirror.setValue(
                      '# Native parity\\n\\nOriginal body.')""")
                    try:
                        page.wait_for_function("""() => typeof wasmReady !== 'undefined' &&
                          wasmReady && !!wasmParity""", timeout=90000)
                    except Exception:
                        state = page.evaluate("""() => ({
                          hardwareConcurrency: navigator.hardwareConcurrency,
                          deviceMemory: navigator.deviceMemory,
                          wasmReady: typeof wasmReady === 'undefined' ? null : wasmReady,
                          worker: typeof wasmWorker === 'undefined' ? null : !!wasmWorker,
                          parity: typeof wasmParity === 'undefined' ? null : wasmParity,
                          fingerprint: window.MARKATA_PREVIEW_FINGERPRINT,
                          wasmVersion: window.MARKATA_WASM_VERSION,
                          status: document.querySelector('#status')?.textContent,
                        })""")
                        raise AssertionError(json.dumps({
                            'browser': state, 'console': console_messages,
                            'page_errors': errors, 'failed_requests': failed_requests,
                        }, indent=2))
                    page.wait_for_function("""() => document.querySelector('#status').textContent
                      .includes('Preview')""", timeout=60000)
                    count = len(rendered_requests)
                    page.evaluate("""() => document.querySelector('.CodeMirror').CodeMirror.setValue(
                      '# Local WASM\\n\\nUpdated paragraph.')""")
                    expect(page.locator('#status')).to_contain_text('Quick local preview', timeout=60000)
                    preview = page.frame_locator('#preview')
                    expect(preview.locator('h1')).to_have_text('Local WASM')
                    expect(preview.locator('p')).to_contain_text('Updated paragraph.')
                    assert len(rendered_requests) == count, 'Local-only edit made a Go /render request'
                    wasm_memory_bytes = page.evaluate(
                        'window.MarkataDiagnostics.snapshot().wasm_worker_memory_bytes')
                    assert wasm_memory_bytes > 0, 'WASM worker memory was not recorded'
                    page.evaluate("""() => {
                      const editor = document.querySelector('.CodeMirror').CodeMirror;
                      editor.focus(); editor.setCursor({line: 2, ch: 1});
                    }""")
                    expect(preview.locator('p.markata-sync-position')).to_contain_text(
                        'Updated paragraph.', timeout=10000)
                    wasm_response = page.evaluate("""async () => {
                      const response = await fetch('/editor/wasm/markata.wasm', {method: 'HEAD'});
                      return {status: response.status, version: response.headers.get('X-Markata-Wasm-Version'),
                        cache: response.headers.get('Cache-Control'), cors: response.headers.get('Access-Control-Allow-Origin')};
                    }""")
                    assert wasm_response['status'] == 200, 'authorized WASM asset request failed'
                    assert wasm_response['version'], 'WASM version header is missing'
                    assert wasm_response['cache'] == 'private, no-store'
                    assert not wasm_response['cors']

                    stale_asset_pattern = re.compile(r'/editor/wasm/markata\.wasm\?v=stale-build$')
                    context.route(stale_asset_pattern, lambda route: route.fulfill(
                        status=200,
                        headers={'Content-Type': 'application/wasm', 'X-Markata-Wasm-Version': 'current-build'},
                        body='',
                    ))
                    mismatch = page.evaluate("""async () => {
                      const worker = window.MarkataWasmPreview.create('stale-build');
                      try { await worker.ready; return ''; }
                      catch (error) { return error.message; }
                      finally { worker.close(); }
                    }""")
                    assert 'version mismatch' in mismatch, mismatch

                    page.evaluate("""() => {
                      const fence = String.fromCharCode(96).repeat(3);
                      document.querySelector('.CodeMirror').CodeMirror.setValue(
                        '# Native fallback\\n\\n' + fence + 'mermaid\\ngraph LR\\nA-->B\\n' + fence);
                    }""")
                    expect(preview.locator('h1')).to_have_text('Native fallback', timeout=60000)
                    assert len(rendered_requests) > count, 'Feature Markdown did not fall back to Go'
                    assert not errors, errors
                    context.close()
                finally:
                    browser.close()
        finally:
            if proxy is not None:
                proxy.shutdown()
                proxy.server_close()
            watchdog.cancel()
            process.terminate()
            process.wait(timeout=15)


if __name__ == '__main__':
    main()
