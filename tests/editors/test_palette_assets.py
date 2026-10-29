"""Packaging regression checks without a renderer or browser dependency."""
import ast
import base64
import io
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
            for name in ['command-palette.js', 'command-palette.css', 'site-tree.js', 'site-tree.css', 'site-tree.LICENSE']:
                self.assertEqual(archive.read(name), (ROOT / 'pkg/workbench/web' / name).read_bytes(),
                                 f'Regenerate the Python editor asset bundle: {name}')


if __name__ == '__main__':
    unittest.main()
