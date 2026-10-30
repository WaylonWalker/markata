"""Render Markdown with the Go renderer bundled in official Markata wheels."""
from __future__ import annotations

import base64
import re
import asyncio
import json
import hashlib
import os
import stat
import subprocess
from dataclasses import dataclass, field
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

    @property
    def feature_head_html(self) -> str:
        """Return only conditional styles required by this rendered fragment."""
        shared = {"css/variables.css", "css/main.css", "css/components.css"}
        tags = []
        for asset in self.assets:
            if asset.get("placement") != "head" or asset.get("kind") != "style" or asset.get("id") in shared:
                continue
            content = asset.get("content")
            if content:
                tags.append("<style>" + content + "</style>")
            elif asset.get("url"):
                tags.append('<link rel="stylesheet" href="' + escape(asset["url"], quote=True) + '">')
        return "\n".join(tags)

    def document(self, title: str = "", *, lang: str = "en",
                 appearance: Appearance | None = None, palette: str | None = None,
                 mode: str | None = None, aesthetic: str | None = None,
                 fontpack: str | None = None, text_size: str | None = None,
                 appearance_bridge: bool = False) -> str:
        """Create a site-compatible article document with optional appearance assets."""
        defaults = appearance.defaults if appearance else {}
        mode = mode if mode is not None else defaults.get('data-theme', 'dark')
        aesthetic = aesthetic if aesthetic is not None else defaults.get('data-aesthetic', 'minimal')
        fontpack = fontpack if fontpack is not None else defaults.get('data-fontpack', 'system')
        text_size = text_size if text_size is not None else defaults.get('data-text-size', 'large')
        palette = palette if palette is not None else defaults.get('data-palette')
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
            '<title>' + escape(title) + '</title>' + initial_style + (self.feature_head_html if appearance else self.head_html)
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
    defaults: dict[str, str] = field(default_factory=dict)

    def asset_url(self, name: str) -> str:
        url = self.asset_base + '/' + name
        if name in self.assets:
            url += '?v=' + hashlib.sha256(self.assets[name]).hexdigest()[:16]
        return url

    @property
    def head_html(self) -> str:
        return ''.join('<link rel="stylesheet" href="' + escape(self.asset_url(name), quote=True) + '">' for name in self.styles)

    @property
    def controls_head_html(self) -> str:
        """Style the appearance controls without applying the site's page layout."""
        return ''.join('<link rel="stylesheet" href="' + escape(self.asset_url(name), quote=True) + '">' for name in self.styles if name not in {'css/main.css', 'css/components.css'})

    @property
    def body_end_html(self) -> str:
        return ''.join('<script src="' + escape(self.asset_url(name), quote=True) + '"></script>' for name in self.scripts)

    def validate(self, attributes: dict[str, str]) -> None:
        for key, value in attributes.items():
            if key in self.choices and value not in self.choices[key]:
                raise ValueError(f"Unknown appearance choice for {key}: {value}")

    @property
    def bridge_html(self) -> str:
        allowed = json.dumps(self.choices).replace('<', r'\u003c')
        return """<script>(()=>{const allowed=__MARKATA_CHOICES__;

const setAppearance=attributes=>{for(const [key,values] of Object.entries(allowed)){const value=attributes?.[key];
if(value==null)document.documentElement.removeAttribute(key);
else if(values.includes(value))document.documentElement.setAttribute(key,value);
}};

let updates=Promise.resolve(), scriptSequence=0;
const originalScripts=new WeakMap();
const scriptSource=s=>originalScripts.get(s) ?? s.textContent;

const renderMarkdown=async data=>{const article=document.querySelector('.post-content');
if(!article)return;
if(window.Chart)article.querySelectorAll('canvas').forEach(canvas=>window.Chart.getChart?.(canvas)?.destroy());
article.innerHTML=data.html;
const head=document.createElement('template');
head.innerHTML=data.feature_head_html;
const styles=[...head.content.querySelectorAll('style,link[rel=stylesheet]')];
const key=n=>n.tagName==='STYLE'?n.textContent:n.href;
const desired=new Set(styles.map(key));
document.querySelectorAll('[data-markata-render-head]').forEach(n=>{if(!desired.has(key(n)))n.remove();
});
for(const source of styles){const existing=[...document.head.querySelectorAll('style,link[rel=stylesheet]')];
if(existing.some(n=>key(n)===key(source)))continue;
const style=source.cloneNode(true);
style.dataset.markataRenderHead='true';
document.head.append(style);
}
const body=document.createElement('template');
body.innerHTML=data.body_end_html;
const sources=[...body.content.querySelectorAll('script')];
const inlineScripts=new Set(sources.filter(s=>!s.src&&s.textContent.includes('rerenderMermaid')).map(s=>s.textContent));
document.querySelectorAll('script[data-markata-generated]:not([src])').forEach(s=>{if(!inlineScripts.has(scriptSource(s))&&!s.textContent.includes('rerenderMermaid'))s.remove();
});
for(const source of sources){if(source.src&&[...document.scripts].some(s=>s.src===source.src))continue;
if(!source.src&&inlineScripts.has(source.textContent)&&[...document.scripts].some(s=>!s.src&&scriptSource(s)===source.textContent))continue;
const script=document.createElement('script');
for(const attr of source.attributes)script.setAttribute(attr.name,attr.value);
script.dataset.markataGenerated='true';
script.textContent=source.textContent;
originalScripts.set(script,source.textContent);
if(!source.src&&source.type==='module'){
  await new Promise(resolve=>{
    const event='markata-module-ready-'+(++scriptSequence);
    addEventListener(event,resolve,{once:true});
    script.addEventListener('error',resolve,{once:true});
    script.textContent+='\\n;window.dispatchEvent(new Event('+JSON.stringify(event)+'));';
    document.body.append(script);
  });
}else if(source.src){await new Promise(resolve=>{script.addEventListener('load',resolve,{once:true});
script.addEventListener('error',resolve,{once:true});
document.body.append(script);
});
}else document.body.append(script);
}dispatchEvent(new Event('view-transition-complete'));
setAppearance(data.attributes);
dispatchEvent(new Event('theme-change'));
};

addEventListener('message',e=>{if(e.source!==parent)return;
if(e.data?.type==='markdown-render'){updates=updates.then(()=>renderMarkdown(e.data)).catch(console.error);
return;
}if(e.data?.type==='appearance'){updates=updates.then(()=>{setAppearance(e.data.attributes);
dispatchEvent(new Event('theme-change'));
}).catch(console.error);
}});
})();
</script>""".replace("__MARKATA_CHOICES__", allowed)


def appearance(*, asset_base: str = "/markata-assets", workspace: str | None = None, timeout: float = 30) -> Appearance:
    """Load the full site appearance catalog once at application startup."""
    command = [_command(False, asset_base, False)[0], 'render', '--appearance']
    if workspace is not None:
        command += ['--workspace', str(workspace)]
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
        return Appearance(data['picker_html'], assets, tuple(data['styles']), tuple(data['scripts']), base, choices, data.get('defaults', {}))
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
           timeout: float = 30, workspace: str | None = None) -> RenderResult:
    """Render using one bundled Go process; input is passed over stdin, never a shell.

    For a standalone response, use inline_assets=True and result.document().
    Otherwise serve bundled asset contents from result.assets below asset_base.
    Rich Markdown output is not an HTML sanitizer. Sanitize untrusted input/output.
    """
    command = _command(inline_assets, asset_base, unsafe)
    if workspace is not None:
        command += ['--workspace', str(workspace)]
    try:
        process = subprocess.run(command, input=markdown.encode("utf-8"),
                                 capture_output=True, timeout=timeout, check=False)
    except (OSError, subprocess.TimeoutExpired) as error:
        raise RenderError(f"Markata render failed: {error}") from error
    return _result(process.stdout, process.stderr, process.returncode)


async def render_async(markdown: str, *, inline_assets: bool = False,
                       asset_base: str = "/markata-assets", unsafe: bool = False,
                       timeout: float = 30, workspace: str | None = None) -> RenderResult:
    """Render without blocking an ASGI event loop; reap the process on cancellation."""
    command = _command(inline_assets, asset_base, unsafe)
    if workspace is not None:
        command += ['--workspace', str(workspace)]
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


def workspace_files(root: str = ".", *, path: str | None = None, metadata: bool = True, timeout: float = 30) -> dict:
    """Discover Markata content, or read one source, through the Go loader."""
    command = [_command(False, "/markata-assets", False)[0], "workspace", "--workspace", str(root)]
    if path is not None:
        command += ["--file", path]
    elif not metadata:
        command += ["--paths"]
    try:
        process = subprocess.run(command, capture_output=True, timeout=timeout, check=False)
        if process.returncode:
            raise RenderError(process.stderr.decode(errors="replace"))
        return json.loads(process.stdout)
    except (OSError, subprocess.TimeoutExpired, ValueError) as error:
        raise RenderError(f"Cannot load workspace: {error}") from error
