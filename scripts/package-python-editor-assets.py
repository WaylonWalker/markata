#!/usr/bin/env python3
"""Refresh the self-contained Python editor archive from canonical web assets."""
import ast
import base64
import io
from pathlib import Path
import re
import zipfile

root = Path(__file__).resolve().parents[1]
page = root / 'examples/fastapi_editor.py'
source = page.read_text()
assignment = next(node for node in ast.parse(source).body if isinstance(node, ast.Assign)
                  and any(isinstance(target, ast.Name) and target.id == 'ASSET_BUNDLE' for target in node.targets))
output = io.BytesIO()
with zipfile.ZipFile(io.BytesIO(base64.b64decode(assignment.value.value))) as archive:
    asset_root = root / 'pkg/workbench/web'
    shared_excalidraw = root / 'pkg/visualizations/excalidraw'
    extra = {
        'search-history.js', 'favicon.ico', 'site-tree.js', 'site-tree.css',
        'site-tree.LICENSE', 'excalidraw-workbench.js', 'excalidraw-workbench.css',
        'excalidraw-editor.js', 'excalidraw-editor.css',
    }
    fonts = shared_excalidraw / 'excalidraw-fonts'
    if fonts.is_dir():
        extra.update(path.relative_to(shared_excalidraw).as_posix() for path in fonts.rglob('*') if path.is_file())
    shared_names = {'excalidraw-editor.js', 'excalidraw-editor.css'}
    if fonts.is_dir():
        shared_names.update(path.relative_to(shared_excalidraw).as_posix() for path in fonts.rglob('*') if path.is_file())
    names = {name for name in archive.namelist() if not name.startswith('excalidraw-fonts/')} | extra
    with zipfile.ZipFile(output, 'w', zipfile.ZIP_DEFLATED) as updated:
        for name in sorted(names):
            asset = (shared_excalidraw / name) if name in shared_names else (asset_root / name)
            info = zipfile.ZipInfo(name, (2026, 10, 29, 0, 0, 0))
            info.compress_type = zipfile.ZIP_DEFLATED
            updated.writestr(info, asset.read_bytes() if asset.is_file() else archive.read(name))
source = re.sub(r"ASSET_BUNDLE = '[^']*'", "ASSET_BUNDLE = '" + base64.b64encode(output.getvalue()).decode() + "'", source)
page.write_text(source)
