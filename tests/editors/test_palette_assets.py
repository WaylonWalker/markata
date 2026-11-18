"""Packaging regression checks without a renderer or browser dependency."""
import ast
import base64
import io
import shutil
import subprocess
from pathlib import Path
import unittest
import zipfile

ROOT = Path(__file__).resolve().parents[2]


class PaletteAssetsTest(unittest.TestCase):
    def test_python_bundle_matches_canonical_palette(self):
        tree = ast.parse((ROOT / 'examples/fastapi_editor.py').read_text())
        bundle = next(node.value.value for node in tree.body if isinstance(node, ast.Assign)
                      and any(isinstance(target, ast.Name) and target.id == 'ASSET_BUNDLE' for target in node.targets))
        with zipfile.ZipFile(io.BytesIO(base64.b64decode(bundle))) as archive:
            for name in ['favicon.ico', 'command-catalog.js', 'command-palette.js', 'command-palette.css', 'command-palette-example.md', 'site-tree.js', 'search-history.js']:
                self.assertEqual(archive.read(name), (ROOT / 'pkg/workbench/web' / name).read_bytes(),
                                 f'Regenerate the Python editor asset bundle: {name}')

    def test_python_bundle_matches_canonical_excalidraw_renderer(self):
        tree = ast.parse((ROOT / 'examples/fastapi_editor.py').read_text())
        bundle = next(node.value.value for node in tree.body if isinstance(node, ast.Assign)
                      and any(isinstance(target, ast.Name) and target.id == 'ASSET_BUNDLE' for target in node.targets))
        renderer = ROOT / 'pkg/visualizations/excalidraw'
        expected = {
            path.relative_to(renderer).as_posix(): path.read_bytes()
            for directory in ('excalidraw-editor-chunks', 'excalidraw-fonts')
            for path in (renderer / directory).rglob('*') if path.is_file()
        }
        expected.update({name: (renderer / name).read_bytes()
                         for name in ('excalidraw-editor.js', 'excalidraw-editor.css')})
        with zipfile.ZipFile(io.BytesIO(base64.b64decode(bundle))) as archive:
            bundled = {name for name in archive.namelist()
                       if name.startswith(('excalidraw-editor-chunks/', 'excalidraw-fonts/'))
                       or name in {'excalidraw-editor.js', 'excalidraw-editor.css'}}
            self.assertEqual(bundled, set(expected), 'Regenerate the Python editor asset bundle: renderer asset list')
            for name, data in expected.items():
                self.assertEqual(archive.read(name), data,
                                 f'Regenerate the Python editor asset bundle: {name}')

    @unittest.skipUnless(shutil.which('node'), 'Node.js is not installed')
    def test_app_shortcut_registry(self):
        subprocess.run(['node', str(ROOT / 'tests/editors/shortcut_registry_checks.cjs')],
                       cwd=ROOT, check=True, capture_output=True, text=True)


if __name__ == '__main__':
    unittest.main()
