"""Render Markdown with the Go renderer bundled in official Markata wheels."""
from __future__ import annotations

import asyncio
import json
import os
import stat
import subprocess
from dataclasses import dataclass
from html import escape
from pathlib import Path
from typing import Any


class RenderError(RuntimeError):
    """The bundled renderer could not produce a supported result."""


@dataclass(frozen=True)
class RenderResult:
    html: str
    features: dict[str, bool]
    assets: tuple[dict[str, Any], ...]
    external_urls: tuple[str, ...]
    head_html: str
    body_end_html: str
    schema_version: int = 1

    def document(self, title: str = "", *, lang: str = "en") -> str:
        """Wrap the fragment and asset tags in a complete HTML document."""
        return (
            '<!doctype html><html lang="' + escape(lang, quote=True) + '"><head>'
            '<meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1">'
            '<title>' + escape(title) + '</title>' + self.head_html
            + '</head><body><article class="post-content">' + self.html
            + '</article>' + self.body_end_html + '</body></html>'
        )


def _command(inline_assets: bool, asset_base: str, unsafe: bool) -> list[str]:
    binary = Path(__file__).parent / "bin" / ("markata.exe" if os.name == "nt" else "markata")
    if not binary.is_file():
        raise RenderError("Bundled Markata executable missing; install an official Markata binary wheel.")
    try:
        if os.name != "nt":
            mode = binary.stat().st_mode
            if not mode & stat.S_IXUSR:
                binary.chmod(mode | stat.S_IXUSR | stat.S_IXGRP | stat.S_IXOTH)
    except OSError as error:
        raise RenderError(f"Cannot prepare bundled Markata executable: {error}") from error
    command = [str(binary), "render", "--json", "--asset-base", asset_base]
    if inline_assets:
        command.append("--inline-assets")
    if unsafe:
        command.append("--unsafe")
    return command


def _result(stdout: bytes, stderr: bytes, returncode: int) -> RenderResult:
    if returncode:
        raise RenderError(stderr.decode("utf-8", errors="replace").strip() or f"Renderer exited with code {returncode}")
    try:
        data = json.loads(stdout)
        if data["schema_version"] != 1:
            raise ValueError("unsupported schema version")
        if not all(isinstance(data[key], str) for key in ("html", "head_html", "body_end_html")):
            raise ValueError("expected HTML strings")
        if not isinstance(data["features"], dict) or not all(type(value) is bool for value in data["features"].values()):
            raise ValueError("expected boolean feature signals")
        if not isinstance(data["assets"], list) or not all(isinstance(asset, dict) for asset in data["assets"]):
            raise ValueError("expected asset list")
        if not isinstance(data["external_urls"], list) or not all(isinstance(url, str) for url in data["external_urls"]):
            raise ValueError("expected external URL list")
        return RenderResult(
            html=data["html"], features=data["features"], assets=tuple(data["assets"]),
            external_urls=tuple(data["external_urls"]), head_html=data["head_html"],
            body_end_html=data["body_end_html"],
        )
    except (ValueError, KeyError, TypeError) as error:
        raise RenderError(f"Invalid renderer response: {error}") from error


def render(markdown: str, *, inline_assets: bool = False,
           asset_base: str = "/markata-assets", unsafe: bool = False,
           timeout: float = 30) -> RenderResult:
    """Render using one bundled Go process; input is passed over stdin, never a shell.

    For a standalone response, use inline_assets=True and result.document().
    Otherwise serve bundled asset contents from result.assets below asset_base.
    Rich Markdown output is not an HTML sanitizer. Sanitize untrusted input/output.
    """
    command = _command(inline_assets, asset_base, unsafe)
    try:
        process = subprocess.run(command, input=markdown.encode("utf-8"),
                                 capture_output=True, timeout=timeout, check=False)
    except (OSError, subprocess.TimeoutExpired) as error:
        raise RenderError(f"Markata render failed: {error}") from error
    return _result(process.stdout, process.stderr, process.returncode)


async def render_async(markdown: str, *, inline_assets: bool = False,
                       asset_base: str = "/markata-assets", unsafe: bool = False,
                       timeout: float = 30) -> RenderResult:
    """Render without blocking an ASGI event loop; reap the process on cancellation."""
    command = _command(inline_assets, asset_base, unsafe)
    try:
        process = await asyncio.create_subprocess_exec(
            *command, stdin=asyncio.subprocess.PIPE, stdout=asyncio.subprocess.PIPE,
            stderr=asyncio.subprocess.PIPE,
        )
    except OSError as error:
        raise RenderError(f"Markata render failed: {error}") from error
    try:
        stdout, stderr = await asyncio.wait_for(process.communicate(markdown.encode("utf-8")), timeout)
    except BaseException as error:
        if process.returncode is None:
            try:
                process.kill()
            except ProcessLookupError:
                pass
        await process.communicate()
        if isinstance(error, asyncio.TimeoutError):
            raise RenderError(f"Markata render timed out after {timeout}s") from error
        raise
    return _result(stdout, stderr, process.returncode)
