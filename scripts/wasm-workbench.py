#!/usr/bin/env python3
"""Build and optionally serve the browser-local Markata workbench."""
import argparse
import base64
import functools
import gzip
import hashlib
from http.server import SimpleHTTPRequestHandler, ThreadingHTTPServer
import json
import os
import re
from pathlib import Path
import shutil
import subprocess
import sys

ROOT = Path(__file__).resolve().parents[1]


def build(output):
    output.mkdir(parents=True, exist_ok=True)
    print('Building Markata WASM and exporting local assets…', file=sys.stderr, flush=True)
    subprocess.run(['go', 'build', '-trimpath', '-ldflags=-s -w', '-o', str(output / 'markata.wasm'), './cmd/markata-wasm'], cwd=ROOT, env={**os.environ, 'GOOS':'js', 'GOARCH':'wasm', 'CGO_ENABLED':'0'}, check=True)
    subprocess.run(['go', 'run', './scripts/wasm-assets', '--output', str(output)], cwd=ROOT, check=True)
    goroot = subprocess.check_output(['go', 'env', 'GOROOT'], text=True).strip()
    shutil.copyfile(Path(goroot) / 'lib/wasm/wasm_exec.js', output / 'wasm_exec.js')
    shutil.copytree(ROOT / 'examples/wasm-workbench', output, dirs_exist_ok=True)
    appearance = json.loads((output / 'appearance.json').read_text())
    # Opaque-origin sandbox requests bypass service workers. Inline font data in
    # the preview only; the editor uses normal, separately cacheable font files.
    font_css = (output / 'appearance/css/fonts.css').read_text()
    def font_data(match):
        path = output / 'appearance' / match.group(1).replace('../', '')
        if not path.is_file():
            raise OSError(f'Missing preview font: {path}')
        return 'url("data:font/woff2;base64,' + base64.b64encode(path.read_bytes()).decode() + '")'
    font_css = re.sub(r'url\([\"\']?(\.\./assets/fonts/[^)\"\']+)[\"\']?\)', font_data, font_css)
    (output / 'appearance/css/preview-fonts.css').write_text(font_css)

    cm = ['lib/codemirror.css', 'lib/codemirror.js', 'mode/xml/xml.js', 'mode/markdown/markdown.js', 'mode/gfm/gfm.js', 'mode/javascript/javascript.js', 'mode/python/python.js', 'mode/css/css.js', 'mode/htmlmixed/htmlmixed.js', 'addon/mode/overlay.js', 'addon/edit/matchbrackets.js', 'keymap/vim.js', 'keymap/emacs.js', 'keymap/sublime.js']
    for name in cm:
        if not (output / 'codemirror' / name).is_file():
            raise OSError(f'Missing bundled editor asset: codemirror/{name}')
    cm_tags = ''.join(f'<link rel="stylesheet" href="codemirror/{name}">' if name.endswith('.css') else f'<script src="codemirror/{name}"></script>' for name in cm)
    styles = ''.join(f'<link rel="stylesheet" href="appearance/{name}">' for name in appearance['styles'])
    scripts = ''.join(f'<script src="appearance/{name}"></script>' for name in appearance['scripts'])
    page = (output / 'index.html').read_text().replace('<!-- ASSETS -->', cm_tags + styles + scripts)
    page = page.replace('<span id="offline"', appearance['picker_html'] + '<span id="offline"', 1)
    (output / 'index.html').write_text(page)
    preview = (output / 'preview.html').read_text().replace('<!-- PREVIEW_STYLES -->', styles)
    (output / 'preview.html').write_text(preview)
    # Content-derived version covers the runtime, app, editor, fonts and theme.
    files = sorted(p for p in output.rglob('*') if p.is_file() and p.name != 'sw.js' and p.suffix != '.gz')
    digest = hashlib.sha256()
    for path in files:
        digest.update(path.relative_to(output).as_posix().encode())
        digest.update(path.read_bytes())
    version = digest.hexdigest()[:16]
    urls = ['./'] + ['./' + path.relative_to(output).as_posix() for path in files]
    sw = '''const PREFIX='markata-wasm-'+self.registration.scope+'-';
const CACHE=PREFIX+VERSION;
const FILES=URLS;
self.addEventListener('install',e=>e.waitUntil((async()=>{
 const cache=await caches.open(CACHE);
 // Sequential fetches keep peak memory bounded on mobile devices.
 for(const file of FILES) {const response=await fetch(file,{cache:'reload'});if(!response.ok)throw new Error(file);await cache.put(file,response);}
 await self.skipWaiting();
})()));
self.addEventListener('activate',e=>e.waitUntil((async()=>{
 for(const key of await caches.keys())if(key.startsWith(PREFIX)&&key!==CACHE)await caches.delete(key);
 await self.clients.claim();
})()));
self.addEventListener('fetch',e=>{
 if(e.request.method!=='GET'||new URL(e.request.url).origin!==location.origin)return;
 e.respondWith(caches.open(CACHE).then(async cache=>(await cache.match(e.request))||fetch(e.request)));
});
'''.replace('VERSION', json.dumps(version)).replace('URLS', json.dumps(urls))
    (output / 'sw.js').write_text(sw)
    for path in [*files, output / 'sw.js']:
        if path.suffix in ('.wasm','.js','.css','.html','.json','.svg'):
            path.with_name(path.name + '.gz').write_bytes(gzip.compress(path.read_bytes(), compresslevel=9, mtime=0))
    print(f'Built {output} ({(output / "markata.wasm.gz").stat().st_size / 1048576:.1f} MiB compressed WASM)')


class Handler(SimpleHTTPRequestHandler):
    def do_GET(self):
        original = self.path
        path = Path(self.translate_path(self.path))
        compressed = path.with_name(path.name + '.gz')
        if path.is_file() and compressed.is_file() and 'gzip' in self.headers.get('Accept-Encoding',''):
            self.path = self.path.split('?',1)[0] + '.gz'
            self.gzip_type = self.guess_type(str(path))
        try:
            super().do_GET()
        finally:
            self.path = original
            self.gzip_type = None

    def end_headers(self):
        self.send_header('Access-Control-Allow-Origin','*')
        self.send_header('Cache-Control','no-cache')
        self.send_header('Vary','Accept-Encoding')
        if getattr(self,'gzip_type',None):
            self.send_header('Content-Encoding','gzip')
        super().end_headers()

    def guess_type(self, path):
        return getattr(self,'gzip_type',None) or super().guess_type(path)


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--output', type=Path, default=ROOT / '.wasm-workbench', help='static output directory (default: .wasm-workbench)')
    parser.add_argument('--serve', action='store_true', help='serve after building on loopback')
    parser.add_argument('--no-build', action='store_true', help='serve an existing build')
    parser.add_argument('--port', type=int, default=8766, help='local port (default: 8766)')
    args = parser.parse_args()
    output = args.output.expanduser().resolve()
    if not args.no_build:
        build(output)
    if args.serve:
        if not (output / 'sw.js').is_file():
            parser.error('No built app found. Remove --no-build to build first.')
        handler = functools.partial(Handler, directory=str(output))
        with ThreadingHTTPServer(('127.0.0.1',args.port),handler) as server:
            print(f'Open http://127.0.0.1:{args.port}/ · Ctrl+C to stop', flush=True)
            try:
                server.serve_forever()
            except KeyboardInterrupt:
                pass


if __name__ == '__main__':
    try:
        main()
    except (OSError, subprocess.CalledProcessError) as error:
        print(f'Workbench failed: {error}',file=sys.stderr)
        sys.exit(1)
