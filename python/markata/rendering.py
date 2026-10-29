"""Render Markdown with the Go renderer bundled in official Markata wheels."""
from __future__ import annotations

import base64
import re
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

    def document(self, title: str = "", *, lang: str = "en",
                 appearance: Appearance | None = None, palette: str | None = None,
                 mode: str = "dark", aesthetic: str = "minimal",
                 fontpack: str = "system", text_size: str = "large",
                 appearance_bridge: bool = False) -> str:
        """Create a site-compatible article document with optional appearance assets."""
        if mode not in {"light", "dark"} or text_size not in {"small", "medium", "large", "x-large"}:
            raise ValueError("Invalid color mode or text size")
        values = {"data-theme": mode, "data-aesthetic": aesthetic,
                  "data-fontpack": fontpack, "data-text-size": text_size}
        if palette:
            values["data-palette"] = palette
        if appearance:
            appearance.validate(values)
        attrs = ''.join(' ' + key + '="' + escape(value, quote=True) + '"' for key, value in values.items())
        bridge = appearance.bridge_html if appearance_bridge and appearance else ''
        background = '#131721' if mode == 'dark' else '#ffffff'
        foreground = '#bfbdb6' if mode == 'dark' else '#1f2937'
        initial_style = '<style>html{color-scheme:' + mode + ';background:' + background + '}body{background:var(--color-background,' + background + ');color:var(--color-text,' + foreground + ')}</style>'
        return (
            '<!doctype html><html lang="' + escape(lang, quote=True) + '"' + attrs + '><head>'
            '<meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1">'
            '<title>' + escape(title) + '</title>' + initial_style + self.head_html
            + (appearance.head_html if appearance else '')
            + '</head><body><main class="main-content"><div class="content-wrapper"><div class="post-shell">'
            '<article class="post h-entry"><div class="post-content e-content">' + self.html
            + '</div></article></div></div></main>' + self.body_end_html + bridge + '</body></html>'
        )


@dataclass(frozen=True)
class Appearance:
    """Reusable site controls and exact-path assets for a host application."""
    picker_html: str
    assets: dict[str, bytes]
    styles: tuple[str, ...]
    scripts: tuple[str, ...]
    asset_base: str
    choices: dict[str, list[str]]

    @property
    def head_html(self) -> str:
        return ''.join('<link rel="stylesheet" href="' + escape(self.asset_base + '/' + name, quote=True) + '">' for name in self.styles)

    @property
    def controls_head_html(self) -> str:
        """Style the appearance controls without applying the site's page layout."""
        return ''.join('<link rel="stylesheet" href="' + escape(self.asset_base + '/' + name, quote=True) + '">' for name in self.styles if name not in {'css/main.css', 'css/components.css'})

    @property
    def body_end_html(self) -> str:
        return ''.join('<script src="' + escape(self.asset_base + '/' + name, quote=True) + '"></script>' for name in self.scripts)

    def validate(self, attributes: dict[str, str]) -> None:
        for key, value in attributes.items():
            if key in self.choices and value not in self.choices[key]:
                raise ValueError(f"Unknown appearance choice for {key}: {value}")

    @property
    def bridge_html(self) -> str:
        allowed = json.dumps(self.choices).replace('<', r'\u003c')
        return '<script>(()=>{const allowed=' + allowed + ";addEventListener('message',e=>{if(e.source!==parent||e.data?.type!=='appearance')return;for(const [key,values] of Object.entries(allowed)){const value=e.data.attributes?.[key];if(value==null)document.documentElement.removeAttribute(key);else if(values.includes(value))document.documentElement.setAttribute(key,value);}dispatchEvent(new Event('theme-change'));});})();</script>"


def appearance(*, asset_base: str = "/markata-assets", timeout: float = 30) -> Appearance:
    """Load the full site appearance catalog once at application startup."""
    command = [_command(False, asset_base, False)[0], 'render', '--appearance']
    try:
        process = subprocess.run(command, capture_output=True, timeout=timeout, check=False)
        if process.returncode:
            raise RenderError(process.stderr.decode(errors='replace'))
        data = json.loads(process.stdout)
        if data['schema_version'] != 1:
            raise ValueError('unsupported appearance schema')
        assets = {name: base64.b64decode(value, validate=True) for name, value in data['assets'].items()}
        base = asset_base.rstrip('/')
        font_css = assets['css/fonts.css'].decode().replace('/assets/fonts/', base + '/assets/fonts/')
        assets['css/fonts.css'] = font_css.encode()
        choices = {'data-theme':['light','dark'], 'data-text-size':['small','medium','large','x-large']}
        for key, name, prop in [('data-palette','palette','palette'), ('data-aesthetic','aesthetic','aesthetic'), ('data-fontpack','fonts','fontpack')]:
            css = assets['css/' + name + '.css'].decode()
            match = re.search(r"--" + prop + r"-manifest:\s*'((?:\\.|[^'])*)'", css)
            if not match:
                raise ValueError('missing ' + prop + ' manifest')
            entries = json.loads(match[1].replace("\\'", "'").replace('\\\\', '\\'))
            choices[key] = [entry['name'] for entry in entries]
        return Appearance(data['picker_html'], assets, tuple(data['styles']), tuple(data['scripts']), base, choices)
    except (OSError, subprocess.TimeoutExpired, ValueError, KeyError, TypeError) as error:
        raise RenderError(f"Invalid appearance response: {error}") from error


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
