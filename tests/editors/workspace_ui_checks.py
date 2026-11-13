"""Exercise real editor assets against fixture HTTP responses (not Go integration).

Run: python tests/editors/workspace_ui_checks.py
CHROMIUM_EXECUTABLE optionally selects a system Chromium installation.
"""
import json
import os
from pathlib import Path
from urllib.parse import urlparse
from playwright.sync_api import sync_playwright, expect

ROOT = Path(__file__).resolve().parents[2]
WEB = ROOT / 'pkg/workbench/web'
FILES = [{'path':'posts/alpha.md','title':'Alpha post'}, {'path':'posts/beta.md','title':'Beta post'}]


def run():
    cm = ['lib/codemirror','mode/xml/xml','mode/markdown/markdown','mode/gfm/gfm',
          'addon/mode/overlay','addon/edit/matchbrackets','addon/hint/show-hint',
          'addon/search/searchcursor','addon/dialog/dialog','keymap/vim','keymap/emacs','keymap/sublime']
    head = '<link rel="stylesheet" href="/editor/codemirror/lib/codemirror.css"><script src="/editor/lsp.js"></script>'
    head += ''.join(f'<script src="/editor/codemirror/{name}.js"></script>' for name in cm)
    colors = '<style>:root{--color-background:#14151b;--color-surface:#20212b;--color-text:#e6e4ed;--color-text-muted:#a09eac;--color-primary:#b792ef}</style>'
    html = (WEB/'index.html').read_text().replace('<!--EDITOR_HEAD-->',head).replace('<!--APPEARANCE_HEAD-->',colors)
    def route(request):
        path=urlparse(request.request.url).path
        if path=='/': return request.fulfill(body=html,content_type='text/html')
        if path.startswith('/editor/'):
            file=WEB/path.removeprefix('/editor/')
            if file.is_file():
                mime='text/javascript' if file.suffix=='.js' else 'text/css' if file.suffix=='.css' else 'text/plain'
                return request.fulfill(body=file.read_bytes(),content_type=mime)
        data={}
        if path=='/workspace': data={'rootUri':'file:///fixture','files':FILES,'feeds':[{'title':'Posts','paths':[x['path'] for x in FILES]}]}
        elif path=='/workspace/file': data={'source':'# Fixture post','version':'a'*64}
        elif path=='/workspace/index-status': data={'files':2,'pending':0,'phase':'ready','ready':True,'last_load_ms':2,'frontmatter_ms':1,'links_ms':1}
        elif path=='/workspace/templates': data={'templates':[]}
        elif path=='/render': data={'document':'<!doctype html><h1>Fixture preview</h1>','html':'<h1>Fixture preview</h1>','source_lines':[],'source_anchors':[],'render_ms':1,'input_bytes':10,'html_bytes':20}
        elif path=='/lsp/config': return request.fulfill(status=503,json={'error':'fixture does not run LSP'})
        return request.fulfill(json=data)
    with sync_playwright() as p:
        browser=p.chromium.launch(headless=True,executable_path=os.environ.get('CHROMIUM_EXECUTABLE'))
        context=browser.new_context(viewport={'width':1440,'height':1000})
        context.route('http://editor.test/**',route)
        page=context.new_page(); errors=[]
        page.on('pageerror',lambda error:errors.append(str(error)))
        page.goto('http://editor.test')
        page.wait_for_function('window.MarkataSettings?.open && window.MarkataSiteTree?.isReady()')
        expect(page.locator('#active-file')).to_have_text('posts/alpha.md')
        page.locator('#settings-open').click()
        expect(page.locator('#settings-dialog')).to_be_visible()
        assert page.locator('#settings-dialog').bounding_box()['width'] > 900
        assert page.locator('#settings-search').bounding_box()['width'] > 800
        page.locator('#settings-search').fill('scroll')
        expect(page.locator('#setting-scroll-sync')).to_be_visible()
        expect(page.locator('#setting-wrap')).to_be_hidden()
        page.locator('#settings-search').fill('zzzzzz')
        expect(page.get_by_text('No matching settings.',exact=True)).to_be_visible()
        page.locator('[data-category="shortcuts"]').click()
        expect(page.locator('#shortcuts-search')).to_be_visible()
        page.locator('#shortcuts-search').fill('Open Editor Settings')
        page.locator('[data-command-id="workbench.settings"]').click()
        page.keyboard.press('Control+Alt+9')
        expect(page.locator('[data-command-id="workbench.settings"]')).to_have_text('Mod+Alt+9')
        page.locator('#settings-close').click()
        page.keyboard.press('Control+Alt+9')
        expect(page.locator('#settings-dialog')).to_be_visible()
        page.locator('[data-category="writing"]').click()
        page.locator('#setting-wrap').uncheck()
        page.locator('#settings-close').click()
        page.reload();page.wait_for_function('window.MarkataSettings?.open')
        page.locator('#settings-open').click()
        expect(page.locator('#setting-wrap')).not_to_be_checked()
        page.locator('#settings-close').click()
        page.locator('#commands').click();page.locator('#commands-query').fill('synscr')
        expect(page.locator('#commands-results')).to_contain_text('Toggle Synchronized Scrolling')
        page.keyboard.press('Enter');page.locator('#commands').click()
        expect(page.locator('#commands-results [role="option"]').first).to_contain_text('Toggle Synchronized Scrolling')
        page.keyboard.press('Escape')
        page.keyboard.press('Control+p');page.locator('#files-query').fill('beta');page.keyboard.press('Enter')
        expect(page.locator('#active-file')).to_have_text('posts/beta.md')
        page.keyboard.press('Control+p')
        expect(page.locator('#files-results [role="option"]').first).to_contain_text('beta.md')
        page.keyboard.press('Escape')
        assert page.locator('#site-sidebar #startup-index').count()==0
        page.locator('#diagnostics-toggle').click();expect(page.locator('#diagnostics-panel')).to_be_visible()
        report=page.evaluate('MarkataDiagnostics.snapshot()');assert report['measurements']['search']['count']>0
        page.get_by_role('button', name='Close diagnostics').click()
        page.locator('#settings-open').click();page.locator('[data-category="advanced"]').click()
        page.locator('#settings-clear-history').click()
        assert page.evaluate("localStorage.getItem('markata-picks-v1:commands')") is None
        page.locator('#setting-pick-history').uncheck();page.locator('#settings-close').click()
        page.locator('#commands').click();page.locator('#commands-query').fill('synscr');page.keyboard.press('Enter')
        assert page.evaluate("localStorage.getItem('markata-picks-v1:commands')") is None
        page.locator('#settings-open').click()
        page.locator('[data-category="writing"]').click()
        screenshot_dir=os.environ.get('MARKATA_BROWSER_ARTIFACTS')
        if screenshot_dir:
            Path(screenshot_dir).mkdir(parents=True,exist_ok=True)
            page.screenshot(path=str(Path(screenshot_dir)/'settings-desktop.png'))
        page.set_viewport_size({'width':390,'height':844})
        bounds=page.locator('#settings-dialog').bounding_box()
        assert bounds['width']<=390 and bounds['height']<=844
        assert page.locator('#settings-dialog').evaluate('e=>e.scrollWidth<=e.clientWidth')
        if screenshot_dir: page.screenshot(path=str(Path(screenshot_dir)/'settings-mobile.png'))
        assert not errors, errors
        print(json.dumps({'result':'passed','checks':'settings search, shortcut recording, persistence, fuzzy commands, recent picks, history clear/disable, diagnostics, mobile bounds','fixture_measurements':report['measurements']},indent=2))
        browser.close()

if __name__=='__main__': run()
