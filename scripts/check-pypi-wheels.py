#!/usr/bin/env python3
"""Validate all of Simon Willison's Go binary wheel targets before PyPI upload."""

from __future__ import annotations

import os
import sys
import zipfile
from pathlib import Path

VERSION = os.environ.get("MARKATA_WHEEL_VERSION", "0.14.0.dev0")
NAME = "markata"
PLATFORMS = {
    "manylinux_2_17_x86_64": b"\x7fELF",
    "manylinux_2_17_aarch64": b"\x7fELF",
    "musllinux_1_2_x86_64": b"\x7fELF",
    "musllinux_1_2_aarch64": b"\x7fELF",
    "macosx_10_9_x86_64": None,  # Mach-O (or universal)
    "macosx_11_0_arm64": None,
    "win_amd64": b"MZ",
    "win_arm64": b"MZ",
}


def check_wheels(folder: Path) -> None:
    files = sorted(folder.glob("*.whl"))
    expected_names = {
        f"{NAME}-{VERSION}-py3-none-{platform}.whl" for platform in PLATFORMS
    }
    actual_names = {path.name for path in files}
    if actual_names != expected_names:
        raise AssertionError(
            f"Unexpected wheel set; missing={sorted(expected_names - actual_names)}, "
            f"extra={sorted(actual_names - expected_names)}"
        )

    for wheel in files:
        platform = wheel.stem.split("-")[-1]
        expected_magic = PLATFORMS[platform]
        binary_name = "markata.exe" if platform.startswith("win_") else "markata"
        with zipfile.ZipFile(wheel) as archive:
            filenames = archive.namelist()
            metadata_path = f"{NAME}-{VERSION}.dist-info/METADATA"
            entrypoint_path = f"{NAME}-{VERSION}.dist-info/entry_points.txt"
            wheel_data_path = f"{NAME}-{VERSION}.dist-info/WHEEL"
            binary_path = f"{NAME}/bin/{binary_name}"
            for entry in (metadata_path, entrypoint_path, wheel_data_path, binary_path, "markata/rendering.py"):
                if entry not in filenames:
                    raise AssertionError(f"{wheel.name} missing {entry}")
            metadata = archive.read(metadata_path).decode()
            entrypoints = archive.read(entrypoint_path).decode()
            wheel_data = archive.read(wheel_data_path).decode()
            binary_magic = archive.read(binary_path)[:4]
            if "Name: markata\n" not in metadata or f"Version: {VERSION}\n" not in metadata:
                raise AssertionError(f"{wheel.name} has unexpected package metadata")
            if "markata = markata:" not in entrypoints:
                raise AssertionError(f"{wheel.name} lacks canonical CLI entry point")
            if f"Tag: py3-none-{platform}" not in wheel_data:
                raise AssertionError(f"{wheel.name} has mismatched platform tag")
            if expected_magic is not None and not binary_magic.startswith(expected_magic):
                raise AssertionError(f"{wheel.name} does not embed expected OS binary")
            if platform.startswith("macosx_") and binary_magic not in {
                bytes.fromhex(magic) for magic in (
                    "feedface", "feedfacf", "cefaedfe", "cffaedfe", "cafebabe", "bebafeca"
                )
            }:
                raise AssertionError(f"{wheel.name} has no Mach-O executable")
        print(f"OK: {wheel.name}")

    print("Eight valid Go executable wheels verified; no upload in this step.")


if __name__ == "__main__":
    try:
        check_wheels(Path(sys.argv[1] if len(sys.argv) > 1 else "dist/pypi-wheels"))
    except (AssertionError, OSError, zipfile.BadZipFile) as error:
        print(f"Wheel smoke failed: {error}", file=sys.stderr)
        raise SystemExit(1) from error
