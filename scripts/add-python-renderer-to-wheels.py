#!/usr/bin/env python3
"""Add the Python API to go-to-wheel output, rebuilding wheel RECORD hashes."""
import base64
import csv
import hashlib
import io
import sys
import zipfile
from pathlib import Path


def add_renderer(wheel: Path) -> None:
    source = Path(__file__).resolve().parent.parent / "python/markata/rendering.py"
    pool_source = source.with_name("pool.py")
    with zipfile.ZipFile(wheel) as archive:
        members = [(item, archive.read(item)) for item in archive.infolist()]
    license_content = (source.parents[2] / 'pkg/plugins/standalone_icons/LICENSE').read_bytes()
    records = []
    temporary = wheel.with_suffix(".tmp")
    try:
        with zipfile.ZipFile(temporary, "w", compression=zipfile.ZIP_DEFLATED) as archive:
            for item, content in members:
                if item.filename.endswith(".dist-info/RECORD"):
                    record_path = item.filename
                    continue
                if item.filename in {"markata/rendering.py", "markata/pool.py"} or item.filename.endswith(".dist-info/licenses/lucide.txt"):
                    continue
                if item.filename == "markata/__init__.py":
                    imports = b"from .rendering import Appearance, appearance, RenderError, RenderResult, render, render_async, workspace_files"
                    if imports not in content:
                        content += b"\n" + imports + b"\n"
                    if b"from .pool import RendererPool" not in content:
                        content += b"from .pool import RendererPool\n"
                archive.writestr(item, content)
                records.append(record(item.filename, content))
            content = source.read_bytes()
            archive.writestr("markata/rendering.py", content)
            records.append(record("markata/rendering.py", content))
            pool_content = pool_source.read_bytes()
            archive.writestr("markata/pool.py", pool_content)
            records.append(record("markata/pool.py", pool_content))
            license_path = record_path.rsplit("/", 1)[0] + "/licenses/lucide.txt"
            archive.writestr(license_path, license_content)
            records.append(record(license_path, license_content))
            buffer = io.StringIO(newline="")
            writer = csv.writer(buffer)
            writer.writerows(records)
            writer.writerow((record_path, "", ""))
            archive.writestr(record_path, buffer.getvalue())
        temporary.replace(wheel)
    finally:
        temporary.unlink(missing_ok=True)


def record(name: str, content: bytes) -> tuple[str, str, str]:
    digest = base64.urlsafe_b64encode(hashlib.sha256(content).digest()).rstrip(b"=").decode()
    return name, "sha256=" + digest, str(len(content))


if __name__ == "__main__":
    wheels = sorted(Path(sys.argv[1]).glob("*.whl"))
    if not wheels:
        raise SystemExit("No wheels found")
    for wheel in wheels:
        add_renderer(wheel)
