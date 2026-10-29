import asyncio
import importlib.util
import json
import os
import subprocess
import sys
import tempfile
import unittest
import zipfile
from pathlib import Path
from unittest.mock import patch

ROOT = Path(__file__).resolve().parents[2]
sys.path.insert(0, str(ROOT / "python/markata"))
import rendering

PAYLOAD = dict(schema_version=1, html='<p>Hello</p>', features={'code': False},
               assets=[], external_urls=[], head_html='', body_end_html='')


class RenderingTest(unittest.TestCase):
    def test_sync_protocol_and_document_escaping(self):
        with patch.object(rendering, '_command', return_value=['renderer']) as command:
            with patch.object(subprocess, 'run', return_value=subprocess.CompletedProcess([], 0, json.dumps(PAYLOAD).encode(), b'')) as run:
                result = rendering.render('hello $(echo injected)', inline_assets=True)
        self.assertEqual(run.call_args.kwargs['input'], b'hello $(echo injected)')
        self.assertNotIn('shell', run.call_args.kwargs)
        command.assert_called_once_with(True, '/markata-assets', False)
        self.assertIn('&lt;unsafe&gt;', result.document('<unsafe>'))
        self.assertIn('lang="&quot;"', result.document(lang='"'))

    def test_bad_protocol_and_process_errors(self):
        for stdout, stderr, code in [(b'bad json', b'', 0), (b'{}', b'', 0),
                                     (b'{"schema_version":2}', b'', 0), (b'', b'failed', 1)]:
            with self.assertRaises(rendering.RenderError):
                rendering._result(stdout, stderr, code)

    def test_timeout_is_render_error(self):
        with patch.object(rendering, '_command', return_value=['renderer']):
            with patch.object(subprocess, 'run', side_effect=subprocess.TimeoutExpired('renderer', 1)):
                with self.assertRaises(rendering.RenderError):
                    rendering.render('hello')

    def test_missing_binary_is_actionable(self):
        with patch.object(rendering.Path, 'is_file', return_value=False):
            with self.assertRaisesRegex(rendering.RenderError, 'official Markata binary wheel'):
                rendering.render('hello')

    @unittest.skipIf(os.name == "nt", "POSIX executable permissions")
    def test_binary_permissions_prepared_on_first_render(self):
        with tempfile.TemporaryDirectory() as folder:
            binary = Path(folder) / 'bin/markata'
            binary.parent.mkdir()
            binary.write_bytes(b'binary')
            binary.chmod(0o644)
            with patch.object(rendering, '__file__', str(Path(folder) / 'rendering.py')):
                command = rendering._command(False, '/assets', False)
            self.assertEqual(command[0], str(binary))
            self.assertTrue(os.access(binary, os.X_OK))

    def test_wheel_injection_preserves_entrypoint_and_records(self):
        spec = importlib.util.spec_from_file_location('wheel_patch', ROOT / 'scripts/add-python-renderer-to-wheels.py')
        module = importlib.util.module_from_spec(spec)
        spec.loader.exec_module(module)
        with tempfile.TemporaryDirectory() as folder:
            wheel = Path(folder) / 'markata.whl'
            with zipfile.ZipFile(wheel, 'w') as archive:
                archive.writestr('markata/__init__.py', 'def main(): pass\n')
                archive.writestr('markata-1.dist-info/RECORD', '')
                binary_info = zipfile.ZipInfo('markata/bin/markata')
                binary_info.external_attr = 0o755 << 16
                archive.writestr(binary_info, b'binary')
            module.add_renderer(wheel)
            with zipfile.ZipFile(wheel) as archive:
                self.assertIn(b'def main()', archive.read('markata/__init__.py'))
                self.assertIn(b'from .rendering import', archive.read('markata/__init__.py'))
                self.assertIn('markata/rendering.py', archive.namelist())
                self.assertEqual(archive.getinfo('markata/bin/markata').external_attr >> 16, 0o755)
                record = archive.read('markata-1.dist-info/RECORD').decode()
                self.assertIn(','.join(module.record('markata/rendering.py', archive.read('markata/rendering.py'))), record)


class AsyncRenderingTest(unittest.IsolatedAsyncioTestCase):
    async def test_real_subprocess_protocol(self):
        program = 'import sys; sys.stdin.read(); print(' + repr(json.dumps(PAYLOAD)) + ')'
        with patch.object(rendering, '_command', return_value=[sys.executable, '-c', program]):
            result = await rendering.render_async('hello')
        self.assertEqual(result.html, PAYLOAD['html'])

    async def test_timeout_and_cancellation_reap_process(self):
        original = asyncio.create_subprocess_exec
        processes = []
        async def spawn(*args, **kwargs):
            process = await original(*args, **kwargs)
            processes.append(process)
            return process
        command = [sys.executable, '-c', 'import time; time.sleep(30)']
        with patch.object(rendering, '_command', return_value=command):
            with patch.object(asyncio, 'create_subprocess_exec', side_effect=spawn):
                with self.assertRaises(rendering.RenderError):
                    await rendering.render_async('', timeout=0.05)
                self.assertIsNotNone(processes[0].returncode)
                task = asyncio.create_task(rendering.render_async(''))
                while len(processes) < 2:
                    await asyncio.sleep(0.01)
                task.cancel()
                with self.assertRaises(asyncio.CancelledError):
                    await task
                self.assertIsNotNone(processes[1].returncode)

class AppearanceDocumentTest(unittest.TestCase):
    def test_feature_head_html_excludes_shared_styles(self):
        result = rendering.RenderResult(
            html='', features={}, assets=(
                {'id': 'css/main.css', 'kind': 'style', 'placement': 'head', 'content': 'shared'},
                {'id': 'css/code.css', 'kind': 'style', 'placement': 'head', 'content': 'feature'},
            ), external_urls=(), head_html='', body_end_html='',
        )
        self.assertEqual(result.feature_head_html, '<style>feature</style>')

    def test_document_site_layout_and_escaping(self):
        result = rendering._result(json.dumps(PAYLOAD).encode(), b'', 0)
        document = result.document('<title>', lang='"')
        self.assertIn('class="post h-entry"', document)
        self.assertIn('class="post-content e-content"', document)
        self.assertIn('&lt;title&gt;', document)
        self.assertNotIn("addEventListener('message'", document)
        self.assertIn('color-scheme:dark', document)
        self.assertIn('color-scheme:light', result.document(mode='light'))
        with self.assertRaises(ValueError):
            result.document(mode='invalid')

    def test_catalog_validation_and_bridge(self):
        presentation = rendering.Appearance('', {}, (), (), '/assets', {
            'data-palette': ['ayu-dark'], 'data-theme': ['dark'],
            'data-aesthetic': ['minimal'], 'data-fontpack': ['system'],
            'data-text-size': ['large'],
        })
        result = rendering._result(json.dumps(PAYLOAD).encode(), b'', 0)
        document = result.document(appearance=presentation, palette='ayu-dark', appearance_bridge=True)
        self.assertIn('e.source!==parent', document)
        self.assertIn('values.includes(value)', document)
        with self.assertRaises(ValueError):
            result.document(appearance=presentation, palette='unknown')


class AppearanceAssetsTest(unittest.TestCase):
    def test_rebased_fonts_and_control_styles(self):
        import base64
        files = {
            'css/palette.css': b"--palette-manifest: '[{\"name\":\"ayu-dark\"}]';",
            'css/aesthetic.css': b"--aesthetic-manifest: '[{\"name\":\"minimal\"}]';",
            'css/fonts.css': b"--fontpack-manifest: '[{\"name\":\"system\"}]'; @font-face{src:url('/assets/fonts/a.woff2')}",
        }
        payload = dict(schema_version=1, picker_html='<button>Colors</button>',
                       assets={key: base64.b64encode(value).decode() for key, value in files.items()},
                       styles=['css/main.css', 'css/components.css', 'css/palette.css'], scripts=[])
        process = subprocess.CompletedProcess([], 0, json.dumps(payload).encode(), b'')
        with patch.object(rendering, '_command', return_value=['renderer']):
            with patch.object(subprocess, 'run', return_value=process):
                presentation = rendering.appearance(asset_base='/custom')
        self.assertIn(b"url('/custom/assets/fonts/a.woff2')", presentation.assets['css/fonts.css'])
        self.assertIn('/custom/css/main.css', presentation.head_html)
        self.assertNotIn('main.css', presentation.controls_head_html)
        self.assertIn('/custom/css/palette.css', presentation.controls_head_html)
        self.assertEqual(presentation.choices['data-palette'], ['ayu-dark'])


if __name__ == '__main__':
    unittest.main()
