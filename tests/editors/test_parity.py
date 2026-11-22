"""Run with an installed Markata wheel after building Go and WASM workbenches."""
import argparse
import base64
from contextlib import ExitStack
from functools import partial
import gzip
from http.server import SimpleHTTPRequestHandler, ThreadingHTTPServer
import json
import os
import re
import time
from pathlib import Path
import shutil
import subprocess
import sys
import tempfile
import threading

from playwright.sync_api import sync_playwright, expect

ROOT = Path(__file__).resolve().parents[2]
sys.path.insert(0, str(ROOT / 'python'))
from editor_test_support import load_editor, serve_editor
from command_palette_checks import check_commands, select_keymap


def check_excalidraw_fence_controls(page):
    scene = (ROOT / 'pkg/plugins/testdata/excalidraw/renderer-corpus/text-fonts.excalidraw').read_text()
    markdown = '# Diagram\n\n```excalidraw\n' + scene.rstrip() + '\n```\n'
    page.evaluate("""markdown => {
      const editor = document.querySelector('.CodeMirror')?.CodeMirror;
      if (!editor) throw new Error('CodeMirror is unavailable');
      editor.setValue(markdown);
      editor.setCursor({line:0,ch:0});
      editor.replaceRange('X', {line:0,ch:0});
    }""", markdown)
    original_source = markdown
    edited_source = 'X' + markdown

    # Cached-only rendering intentionally rejects a scene without an artifact.
    # Wait for this preview to settle before exercising the keyboard controls.
    expect(page.get_by_role('alert').filter(has_text='no matching cached Excalidraw SVG')).to_be_visible(timeout=60000)

    show = page.locator('.excalidraw-fence-controls').get_by_role('button', name='Show Excalidraw JSON')
    expect(show).to_be_visible()
    expect(show).to_have_attribute('aria-expanded', 'false')
    show.focus()
    page.keyboard.press('Tab')
    edit = page.locator('.excalidraw-fence-controls__button--edit')
    expect(edit).to_be_focused()
    page.keyboard.press('Enter')

    dialog = page.locator('#excalidraw-dialog')
    expect(dialog).to_be_visible()
    canvas = page.locator('#excalidraw-canvas .excalidraw')
    expect(canvas).to_be_visible(timeout=30000)
    expect(canvas).to_have_class(re.compile(r'\btheme--dark\b'))
    page.locator('#excalidraw-cancel').click()
    expect(dialog).not_to_be_visible()
    # Modal Cancel restores the keyboard launcher rather than the hidden
    # CodeMirror input. Wait for that restoration before moving to Show JSON.
    expect(edit).to_be_focused()
    show = page.locator('.excalidraw-fence-controls').get_by_role('button', name='Show Excalidraw JSON')
    expect(show).to_be_visible()

    show.focus()
    expect(show).to_be_focused()
    page.keyboard.press('Space')
    hide = page.locator('.excalidraw-fence-controls').get_by_role('button', name='Hide Excalidraw JSON')
    expect(hide).to_be_visible()
    expect(hide).to_have_attribute('aria-expanded', 'true')
    hide.focus()
    page.keyboard.press('Space')
    show = page.locator('.excalidraw-fence-controls').get_by_role('button', name='Show Excalidraw JSON')
    expect(show).to_be_focused()
    assert page.evaluate("document.querySelector('.CodeMirror').CodeMirror.getValue()") == edited_source
    page.evaluate("document.querySelector('.CodeMirror').CodeMirror.undo()")
    assert page.evaluate("document.querySelector('.CodeMirror').CodeMirror.getValue()") == original_source

    page.evaluate("document.querySelector('.CodeMirror').CodeMirror.setOption('readOnly', true)")
    expect(page.locator('.excalidraw-fence-controls__button--edit')).to_be_disabled()
    page.evaluate("document.querySelector('.CodeMirror').CodeMirror.setOption('readOnly', false)")

    page.locator('.excalidraw-fence-controls__button--edit').click()
    expect(page.locator('#excalidraw-canvas .excalidraw')).to_be_visible(timeout=30000)
    page.locator('#excalidraw-apply').click()
    expect(dialog).not_to_be_visible(timeout=30000)
    applied = page.evaluate("document.querySelector('.CodeMirror').CodeMirror.getValue()")
    scene_json = re.search(r'```excalidraw\s*([\s\S]*?)\s*```', applied).group(1)
    saved_scene = json.loads(scene_json)
    artifact = saved_scene['_markata']['excalidraw_svg']
    assert artifact['format'] == 'markata-excalidraw-svg-v1'
    assert artifact['theme'] == 'dark'
    assert re.fullmatch(r'[0-9a-f]{64}', artifact['source_sha256'])
    assert re.fullmatch(r'[0-9a-f]{64}', artifact['svg_sha256'])
    svg = gzip.decompress(base64.b64decode(artifact['svg'])).decode()
    assert 'https://esm.sh/' not in svg
    assert 'data:font/woff2;base64,' in svg
    expect(page.frame_locator('#preview').locator('.markata-visualization img')).to_be_visible(timeout=60000)
    expect(page.get_by_role('alert').filter(has_text='no matching cached Excalidraw SVG')).not_to_be_visible()


def check(browser, url, excalidraw=False):
    context = browser.new_context(accept_downloads=True)
    context.tracing.start(screenshots=True, snapshots=True, sources=True)
    try:
        page = context.new_page()
        errors = []
        page.on('pageerror', lambda error: errors.append(str(error)))
        page.goto(url)
        favicon_url = page.locator('link[rel=icon]').evaluate('element => element.href')
        assert favicon_url.endswith('/favicon.ico'), favicon_url
        favicon = page.request.get(favicon_url)
        assert favicon.status == 200, f'{favicon_url} returned {favicon.status}'
        assert favicon.body() == (ROOT / 'pkg/workbench/web/favicon.ico').read_bytes()
        preview = page.frame_locator('#preview')
        expect(preview.locator('h1')).to_be_visible(timeout=60000)
        check_commands(page, preview)
        if excalidraw:
            check_excalidraw_fence_controls(page)
        expect(page.locator('#saved')).to_contain_text(re.compile(r'Draft saved|Saved to site'))
        assert page.locator('header h1').evaluate("e => getComputedStyle(e).fontFamily") != 'Georgia, serif'
        assert not page.locator('link[href$="/css/main.css"]').count()
        select_keymap(page, 'vim')
        page.locator('#import').set_input_files({'name': 'draft.md', 'mimeType': 'text/markdown', 'buffer': '# Imported café\n\nKeep this draft.'.encode()})
        expect(preview.locator('h1')).to_have_text('Imported café')
        page.reload()
        expect(preview.locator('h1')).to_have_text('Imported café', timeout=60000)
        expect(page.locator('#keymap')).to_have_value('vim')
        select_keymap(page, 'default')
        # Bytes, rather than characters, determine the import limit.
        page.locator('#import').set_input_files({'name': 'large.md', 'mimeType': 'text/markdown', 'buffer': ('界' * 70000).encode()})
        expect(page.locator('#error')).to_be_visible()
        expect(preview.locator('h1')).to_have_text('Imported café')
        page.locator('#import').set_input_files({'name': 'valid.md', 'mimeType': 'text/markdown', 'buffer': b'# Recovered'})
        expect(preview.locator('h1')).to_have_text('Recovered')
        expect(page.locator('#error')).to_be_hidden()
        # Storage failure must keep editing and export available.
        page.evaluate("() => { Storage.prototype.setItem = function() {throw Error('storage full')} }")
        page.evaluate("""() => {
          const editor = document.querySelector('.CodeMirror').CodeMirror;
          editor.setSelection(editor.posFromIndex(0), editor.posFromIndex(editor.getValue().length));
          editor.focus();
        }""")
        page.keyboard.type('# Unsaved draft')
        expect(page.locator('#saved')).to_contain_text('Draft not saved')
        expect(preview.locator('h1')).to_have_text('Unsaved draft')
        with page.expect_download() as download:
            if page.locator('#download').is_visible():
                page.locator('#download').click()
            else:
                page.locator('#commands').click()
                query = page.get_by_role('combobox', name='Find a command or action')
                query.fill('Download current file')
                page.keyboard.press('Enter')
        assert Path(download.value.path()).read_text() == '# Unsaved draft'
        assert not errors, errors
        typography = lambda locator: locator.evaluate('e => [getComputedStyle(e).fontFamily, getComputedStyle(e).fontSize]')
        return repr([typography(page.locator('header h1')), typography(page.locator('.CodeMirror')),
                     typography(preview.locator('.post-content')), typography(preview.locator('h1'))])
    except Exception:
        # Record the failing variant's browser state before closing its context.
        directory = os.environ.get("MARKATA_BROWSER_ARTIFACTS")
        if directory:
            try:
                folder = Path(directory)
                folder.mkdir(parents=True, exist_ok=True)
                stamp = time.time_ns()
                for index, current in enumerate(context.pages):
                    current.screenshot(path=str(folder / f"parity-{stamp}-{index}.png"), full_page=True)
                context.tracing.stop(path=str(folder / f"parity-{stamp}.zip"))
            except Exception as artifact_error:
                print(f"Could not capture browser evidence: {artifact_error}", file=sys.stderr)
        raise
    finally:
        context.close()


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--go-binary', required=True)
    parser.add_argument('--wasm-output', type=Path, required=True)
    args = parser.parse_args()
    with ExitStack() as stack:
        fixture_site = Path(__file__).resolve().parents[2] / 'examples/editor-site'
        site = Path(stack.enter_context(tempfile.TemporaryDirectory(prefix='markata-editor-parity-'))) / 'site'
        shutil.copytree(fixture_site, site)
        demo = load_editor()
        demo.LSP_WORKSPACE = site
        demo.presentation = demo.appearance(asset_base="/appearance", workspace=str(site))
        python_url = stack.enter_context(serve_editor(demo.app))
        process = subprocess.Popen([args.go_binary, '--listen', '127.0.0.1:0', '--workspace', str(site)], stderr=subprocess.PIPE, text=True)
        stack.callback(process.wait, timeout=15)
        stack.callback(process.terminate)
        startup_timeout = threading.Timer(30, process.terminate)
        startup_timeout.start()
        try:
            match = None
            for line in process.stderr:
                match = re.search(r'http://127\.0\.0\.1:\d+/', line)
                if match:
                    break
            if not match:
                raise RuntimeError('Go workbench did not start within 30 seconds')
        finally:
            startup_timeout.cancel()
        server = ThreadingHTTPServer(('127.0.0.1', 0), partial(SimpleHTTPRequestHandler, directory=str(args.wasm_output)))
        stack.callback(server.server_close)
        threading.Thread(target=server.serve_forever, daemon=True).start()
        stack.callback(server.shutdown)
        driver = stack.enter_context(sync_playwright())
        browser = driver.chromium.launch()
        stack.callback(browser.close)
        fonts = []
        for name, url in [('Python', python_url), ('Go', match[0]), ('WASM', f'http://127.0.0.1:{server.server_port}/')]:
            if name == 'Go':
                from site_tree_checks import check_site_tree
                with browser.new_context() as context:
                    check_site_tree(context.new_page(), url)
                print(name + ': workspace tree passed', flush=True)
            fonts.append(check(browser, url, excalidraw=name in ('Python', 'Go')))
            print(name + ': editor parity passed', flush=True)
            if name == 'Go':
                from shortcut_ui_checks import check_app_shortcuts
                with browser.new_context() as context:
                    check_app_shortcuts(context.new_page(), url)
                print('Go: app keyboard shortcuts passed', flush=True)
            if name == 'Python':
                from site_tree_checks import check_site_tree
                with browser.new_context() as context:
                    check_site_tree(context.new_page(), url)
                print(name + ': workspace tree passed', flush=True)
        assert len(set(fonts)) == 1, fonts


if __name__ == '__main__':
    main()
