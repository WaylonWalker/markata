"""Run after building: uv run --with playwright python tests/wasm/test_workbench.py."""
import importlib.util
import os
from pathlib import Path
import statistics
import threading
import unittest
from http.server import ThreadingHTTPServer
from functools import partial

from playwright.sync_api import sync_playwright, expect

ROOT = Path(__file__).resolve().parents[2]
spec = importlib.util.spec_from_file_location('workbench', ROOT / 'scripts/wasm-workbench.py')
workbench = importlib.util.module_from_spec(spec)
spec.loader.exec_module(workbench)


class WorkbenchTest(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        output = Path(os.environ.get('MARKATA_WASM_OUTPUT', ROOT / '.wasm-workbench'))
        cls.server = ThreadingHTTPServer(('127.0.0.1', 0), partial(workbench.Handler, directory=str(output)))
        threading.Thread(target=cls.server.serve_forever, daemon=True).start()
        cls.addClassCleanup(cls.server.server_close)
        cls.addClassCleanup(cls.server.shutdown)
        cls.url = f'http://127.0.0.1:{cls.server.server_port}/'
        cls.driver = sync_playwright().start()
        cls.addClassCleanup(cls.driver.stop)
        cls.browser = cls.driver.chromium.launch(headless=True)
        cls.addClassCleanup(cls.browser.close)

    def setUp(self):
        self.context = self.browser.new_context(viewport={'width': 1500, 'height': 1000}, accept_downloads=True)
        self.addCleanup(self.context.close)
        self.page = self.context.new_page()
        self.errors = []
        self.requests = []
        self.page.on('pageerror', lambda error: self.errors.append(str(error)))
        self.page.on('request', lambda request: self.requests.append(request))
        self.page.goto(self.url)
        self.assertEqual(self.errors, [], 'JavaScript failed during application startup')
        expect(self.page.locator('#status')).to_contain_text('rendered locally', timeout=60000)
        self.frame = self.page.frame_locator('#preview')
        expect(self.frame.locator('h1')).to_have_text('A place to think')

    def tearDown(self):
        self.assertEqual(self.errors, [])
        self.assertFalse([r.url for r in self.requests if r.method != 'GET'], 'Draft data sent over HTTP')
        self.assertFalse([r.url for r in self.requests if not r.url.startswith(self.url)], 'External network dependency')

    def edit(self, source):
        self.page.evaluate("source => document.querySelector('.CodeMirror').CodeMirror.setValue(source)", source)

    def test_offline_recovery_and_export(self):
        expect(self.page.locator('#offline')).to_have_text('Ready offline', timeout=60000)
        self.edit('# Saved offline\n\n私の草稿')
        expect(self.frame.locator('h1')).to_have_text('Saved offline')
        self.context.set_offline(True)
        self.page.reload()
        expect(self.frame.locator('h1')).to_have_text('Saved offline', timeout=60000)
        self.page.evaluate("const cm = document.querySelector('.CodeMirror').CodeMirror; cm.setCursor(cm.lastLine(), cm.getLine(cm.lastLine()).length)")
        self.page.locator('#commands').click()
        self.page.get_by_role('combobox', name='Find a command').fill('note box')
        self.page.keyboard.press('Enter')
        expect(self.frame.locator('.admonition.note')).to_contain_text('Write your note here.')
        self.edit('# Still local\n\nNo server.\n\n```chartjs\n{"type":"bar","data":{"labels":["A"],"datasets":[{"data":[2]}]}}\n```')
        expect(self.frame.locator('h1')).to_have_text('Still local')
        expect(self.frame.locator('canvas')).to_be_visible()
        self.page.evaluate("document.documentElement.dataset.theme='light'")
        expect(self.frame.locator('html')).to_have_attribute('data-theme','light')
        with self.page.expect_download() as info:
            self.page.locator('#download').click()
        self.assertIn('# Still local\n\nNo server.', Path(info.value.path()).read_text())
        expect(self.page.locator('#saved')).to_have_text('Draft saved on this device')

    def test_repeated_rich_previews_and_fast_edits(self):
        chart = '\n```chartjs\n{"type":"bar","data":{"labels":["A"],"datasets":[{"label":"Value","data":[2]}]}}\n```'
        diagram = '\n```mermaid\ngraph LR\nA-->B\n```'
        for i in range(4):
            self.edit(f'# Revision {i}\n\n:lucide-rocket:\n' + chart + diagram)
            expect(self.frame.locator('h1')).to_have_text(f'Revision {i}')
            expect(self.frame.locator('canvas')).to_be_visible()
            expect(self.frame.locator('pre.mermaid svg')).to_be_visible(timeout=15000)
            preview = self.page.frames[1]
            self.assertEqual(preview.evaluate('Object.keys(Chart.instances).length'), 1)
            expect(self.frame.locator('span.icon svg')).to_have_count(1)
        self.page.evaluate("() => {const e=document.querySelector('.CodeMirror').CodeMirror;for(let i=0;i<30;i++)e.setValue('# Latest '+i);}")
        expect(self.frame.locator('h1')).to_have_text('Latest 29')
        self.assertEqual(self.page.frames[1].evaluate('Object.keys(Chart.instances).length'), 0)
        self.edit('# Safe\n\n<script>parent.hacked=true</script>')
        expect(self.frame.locator('h1')).to_have_text('Safe')
        self.assertFalse(self.page.evaluate('Boolean(window.hacked)'))
        expect(self.frame.locator('.post-content script')).to_have_count(0)

    def test_tables_reinitialize_after_replacement(self):
        for title in ['First table', 'Second table']:
            self.edit('# '+title+'\n\n| Item | Value |\n| --- | --- |\n| B | 2 |\n| A | 1 |')
            expect(self.frame.locator('h1')).to_have_text(title)
            expect(self.frame.locator('th button')).to_have_count(2)
            self.frame.locator('th button').first.click()
            expect(self.frame.locator('tbody tr').first).to_contain_text('A')
            self.frame.locator('th button').first.click()
            expect(self.frame.locator('tbody tr').first).to_contain_text('B')

    def test_import_keymap_mobile_and_size_limit(self):
        self.page.locator('#open').focus()
        with self.page.expect_file_chooser() as chooser:
            self.page.keyboard.press('Enter')
        chooser.value.set_files({'name':'notes.md','mimeType':'text/markdown','buffer':b'# Imported'})
        expect(self.frame.locator('h1')).to_have_text('Imported')
        self.page.locator('#keymap').select_option('vim')
        self.assertEqual(self.page.evaluate("document.querySelector('.CodeMirror').CodeMirror.getOption('keyMap')"), 'vim')
        self.page.set_viewport_size({'width':390,'height':844})
        self.assertLessEqual(self.page.evaluate('document.documentElement.scrollWidth'), 390)
        self.edit('界' * 70000)
        expect(self.page.locator('#error')).to_contain_text('200,000 UTF-8 bytes')
        self.edit('# Recovered')
        expect(self.frame.locator('h1')).to_have_text('Recovered')

    def test_storage_failure_does_not_block_editor(self):
        self.page.evaluate("() => {Storage.prototype.setItem=()=>{throw new Error('quota')}}")
        self.edit('# Keep writing')
        expect(self.frame.locator('h1')).to_have_text('Keep writing')
        expect(self.page.locator('#saved')).to_contain_text('Draft not saved')

    def test_worker_api_and_warm_performance(self):
        result = self.page.evaluate('''async () => {
          const worker = new Worker('worker.js');
          const next = () => new Promise(resolve => worker.onmessage=e=>resolve(e.data));
          await next();
          async function render(source,revision) {const promise=next();worker.postMessage({source,revision});return promise;}
          const invalid=await render(42,1);
          const oversized=await render('界'.repeat(70000),2);
          const samples=[];
          for(let i=0;i<8;i++) {const r=await render('# Heading\\n\\nA paragraph.\\n'.repeat(50),i+3);samples.push(r.render_ms);}
          const final=await render('# Alive',20);worker.terminate();
          return {invalid,oversized,samples,final};
        }''')
        self.assertIn('error', result['invalid']['result'])
        self.assertIn('error', result['oversized']['result'])
        self.assertIn('<h1', result['final']['result']['html'])
        self.assertLess(len(str(result['final']['result'])), 10000, 'Bundled asset bytes leaked into every response')
        # Report measurements; do not impose machine-specific timing thresholds.
        print('Warm 50-section document median ms:', round(statistics.median(result['samples']), 2))


if __name__ == '__main__':
    unittest.main(verbosity=2)
