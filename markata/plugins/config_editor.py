"""
Config Editor Plugin

Web-based editor for `markata.toml` that understands Pydantic config models and
preserves TOML ordering, comments, and formatting using tomlkit.

It also discovers Markata plugins (both built-in `markata.plugins.*` and local
plugins in configurable directories like `plugins/`) and surfaces their
module-level help/docstrings so users can browse and jump to each plugin's
config section.

# Installation

```bash
uv add tomlkit fastapi uvicorn
# or
pip install tomlkit fastapi uvicorn
```

Enable the plugin in your `markata.toml`:

```toml
[markata]
hooks = [
  "markata.plugins.config_model",
  "markata.plugins.base_cli",
  "markata.plugins.config_editor",  # 👈 add this
]

[markata.config_editor]
config_path = "markata.toml"  # Path to your config file
host = "127.0.0.1"            # Host for the web UI
port = 8765                   # Port for the web UI
open_browser = true           # Auto-open browser
plugin_dirs = ["plugins"]    # Local plugin directories to scan
```

# Configuration

```toml
[markata.config_editor]
config_path = "markata.toml"  # Path to the TOML file to edit
host = "127.0.0.1"            # Host bind for the web UI
port = 8765                   # Port for the web UI
open_browser = true           # Whether to open the browser automatically
plugin_dirs = ["plugins"]    # Local plugin directories to scan for plugins
```

# Usage

Launch the web UI from your project root (where `markata.toml` lives):

```bash
markata config-web
# or override host/port
markata config-web --host 0.0.0.0 --port 8001 --no-open-browser
```

Then open the printed URL in a browser.

The UI will:

- Flatten all Pydantic config models under `markata.config` into dotted paths
  like `output_dir`, `glob.glob_patterns`, `profiler.output_dir`, etc.
- Show type names and field descriptions.
- Let you edit scalar fields with text inputs.
- Let you edit list/dict fields with JSON-like text areas.
- Discover plugins and show their module-level help/docstrings.
- Let you filter config options by plugin (e.g. show only `tag_stats.*`).

Edits are applied to the in-memory `markata.config` instance **and** written
back to `markata.toml` via tomlkit so that your original key ordering and
comments are preserved as much as possible.

# Notes

- For list/dict fields, values are parsed as JSON by default. So for example:

  - List of strings: `["a", "b", "c"]`
  - Dict: `{"home": "/", "blog": "/blog/"}`

- If parsing fails, the update is rejected with an error message.
- This is intentionally conservative to avoid corrupting your config.
"""

from __future__ import annotations

from dataclasses import dataclass
from pathlib import Path
from typing import Any, Dict, Iterable, List, Optional, Sequence, Tuple, Union, get_args, get_origin, TYPE_CHECKING

import ast
import importlib
import inspect
import json
import pkgutil
import textwrap
import webbrowser

import pydantic
from pydantic import ConfigDict
from pydantic.networks import AnyUrl
from pydantic_settings import BaseSettings

from markata.hookspec import hook_impl, register_attr

from fastapi import FastAPI, Form
from fastapi.responses import HTMLResponse, RedirectResponse, JSONResponse
import tomlkit
import typer
import uvicorn

if TYPE_CHECKING:  # pragma: no cover
    from markata import Markata


MARKATA_PLUGIN_NAME = "Config Editor"
MARKATA_PLUGIN_PACKAGE_NAME = "config-editor"


class ConfigEditorSettings(pydantic.BaseModel):
    """Settings for the config editor UI."""

    config_path: Path = Path("markata.toml")
    host: str = "127.0.0.1"
    port: int = 8765
    open_browser: bool = True
    plugin_dirs: List[Path] = [Path("plugins")]

    model_config = ConfigDict(
        validate_assignment=True,
        arbitrary_types_allowed=True,
        extra="allow",
        str_strip_whitespace=True,
        validate_default=True,
        coerce_numbers_to_str=True,
        populate_by_name=True,
    )


class Config(pydantic.BaseModel):
    """Plugin config container.

    This is mounted under `markata.config.config_editor`.
    """

    config_editor: ConfigEditorSettings = ConfigEditorSettings()


@hook_impl()
@register_attr("config_models")
def config_model(markata: "Markata") -> None:
    """Register the config editor Pydantic model so it becomes part of markata.config."""

    markata.config_models.append(Config)


# ----------------- INTROSPECTION & FLATTENING -----------------


@dataclass
class ConfigItem:
    path: str                  # dotted path, e.g. "glob.glob_patterns"
    toml_path: Tuple[str, ...] # TOML table path, e.g. ("markata", "glob", "glob_patterns")
    type_name: str
    field_type: Any
    description: str
    value: Any


@dataclass
class PluginInfo:
    name: str               # Human name (MARKATA_PLUGIN_NAME or derived)
    module: str             # Module path, e.g. "markata.plugins.tag_stats"
    source: str             # "core" or "local"
    help: str               # First paragraph of module docstring
    config_root: Optional[str]  # Attribute on markata.config, e.g. "tag_stats"


def _get_model_fields(model: BaseSettings | pydantic.BaseModel) -> Dict[str, Any]:
    """Pydantic v2 vs v1 compatibility for field metadata."""

    cls = model.__class__
    if hasattr(cls, "model_fields"):
        return cls.model_fields  # type: ignore[attr-defined]
    if hasattr(cls, "__fields__"):
        return cls.__fields__  # type: ignore[attr-defined]
    return {}


def _field_annotation(field: Any) -> Any:
    """Pydantic v2 vs v1 compatibility for getting the declared type."""

    if hasattr(field, "annotation"):
        return field.annotation
    if hasattr(field, "outer_type_"):
        return field.outer_type_
    return Any


def _type_name(tp: Any) -> str:
    """Human-readable type name for UI."""

    origin = get_origin(tp) or tp
    if origin is list or origin is List:
        args = get_args(tp)
        inner = _type_name(args[0]) if args else "Any"
        return f"List[{inner}]"
    if origin is dict or origin is Dict:
        args = get_args(tp)
        key = _type_name(args[0]) if args else "Any"
        val = _type_name(args[1]) if len(args) > 1 else "Any"
        return f"Dict[{key}, {val}]"
    if origin is Union:
        args = get_args(tp)
        return " | ".join(_type_name(a) for a in args)
    if isinstance(origin, type):
        return origin.__name__
    return str(tp)


def _flatten_config_model(
    obj: BaseSettings | pydantic.BaseModel,
    base_path: Tuple[str, ...],
    toml_prefix: Tuple[str, ...],
) -> Iterable[ConfigItem]:
    """Recursively flatten all Pydantic config models into ConfigItems.

    base_path: dotted path segments (excluding "markata" root)
    toml_prefix: TOML path segments (including "markata" root)
    """

    fields = _get_model_fields(obj)
    for name, field in fields.items():
        annotation = _field_annotation(field)
        description = getattr(field, "description", "") or getattr(field, "title", "") or ""

        dotted = ".".join(base_path + (name,))
        toml_path = toml_prefix + (name,)
        value = getattr(obj, name, None)

        item = ConfigItem(
            path=dotted,
            toml_path=toml_path,
            type_name=_type_name(annotation),
            field_type=annotation,
            description=description,
            value=value,
        )
        yield item

        # Recurse into nested BaseModel configs
        if isinstance(value, (BaseSettings, pydantic.BaseModel)):
            yield from _flatten_config_model(
                value,
                base_path + (name,),
                toml_prefix + (name,),
            )


def iter_config_items(markata: "Markata") -> Iterable[ConfigItem]:
    """Yield all config items from markata.config and nested plugin config models."""

    root = markata.config  # BaseSettings from config_model.py (root of [markata])
    return _flatten_config_model(root, base_path=(), toml_prefix=("markata",))


def get_item_by_path(markata: "Markata", dotted_path: str) -> Optional[ConfigItem]:
    """Find a ConfigItem by its dotted path."""

    for item in iter_config_items(markata):
        if item.path == dotted_path:
            return item
    return None


# ----------------- VALUE COERCION -----------------


def _parse_bool(value: str) -> bool:
    v = value.strip().lower()
    return v in {"1", "true", "yes", "y", "on"}


def coerce_value(raw: str, field_type: Any) -> Any:
    """Convert string from UI into an appropriate Python type.

    Uses the Pydantic field annotation as a hint.
    """

    raw = raw.strip()
    origin = get_origin(field_type) or field_type

    # Optional[...] unwrap
    if origin is Union:
        args = [a for a in get_args(field_type) if a is not type(None)]  # noqa: E721
        if len(args) == 1:
            return coerce_value(raw, args[0])

    if origin in (str,):
        return raw
    if origin in (int,):
        return int(raw)
    if origin in (float,):
        return float(raw)
    if origin in (bool,):
        return _parse_bool(raw)
    if origin in (Path,):
        return Path(raw)
    if origin in (AnyUrl,):
        # Let Pydantic validate later; treat as string here.
        return raw

    if origin in (list, List):
        # Expect JSON array by default
        if not raw:
            return []
        try:
            parsed = json.loads(raw)
        except json.JSONDecodeError:
            # Fallback: comma-separated
            return [s.strip() for s in raw.split(",") if s.strip()]
        if not isinstance(parsed, list):
            raise ValueError("Expected JSON array for list field")
        return parsed

    if origin in (dict, Dict):
        if not raw:
            return {}
        parsed = json.loads(raw)
        if not isinstance(parsed, dict):
            raise ValueError("Expected JSON object for dict field")
        return parsed

    # Fallback to raw string
    return raw


# ----------------- TOMLKIT UPDATER -----------------


def update_toml_value(config_path: Path, toml_path: Sequence[str], value: Any) -> None:
    """Update a given path inside markata.toml using tomlkit.

    Preserves formatting and ordering as much as possible.

    toml_path is like ("markata", "glob", "glob_patterns").
    """

    def _to_toml_value(v: Any) -> Any:
        """Convert Python values to TOML-friendly ones.

        - Path -> str
        - Recurse into lists/dicts
        - Leave other scalar types as-is (tomlkit knows ints, floats, bools,
          datetimes, etc.).
        """

        if isinstance(v, Path):
            return str(v)
        if isinstance(v, list):
            return [_to_toml_value(i) for i in v]
        if isinstance(v, dict):
            return {k: _to_toml_value(i) for k, i in v.items()}
        return v

    text = config_path.read_text()
    doc = tomlkit.parse(text)

    cur = doc
    # Ensure nested tables exist
    for key in toml_path[:-1]:
        if key not in cur:
            cur[key] = tomlkit.table()
        cur = cur[key]

    leaf = toml_path[-1]
    cur[leaf] = tomlkit.item(_to_toml_value(value))

    config_path.write_text(tomlkit.dumps(doc))


# ----------------- PLUGIN DISCOVERY -----------------


def _first_paragraph(doc: str, max_chars: int = 400) -> str:
    if not doc:
        return ""
    para = doc.strip().split("\n\n", 1)[0].strip()
    para = " ".join(para.split())
    if len(para) > max_chars:
        para = para[: max_chars - 1].rstrip() + "…"
    return para


def _plugin_info_from_module(module_name: str, source: str, markata: "Markata") -> PluginInfo:
    mod = importlib.import_module(module_name)
    doc = inspect.getdoc(mod) or ""
    help_text = _first_paragraph(doc)

    pkg_name = getattr(mod, "MARKATA_PLUGIN_PACKAGE_NAME", module_name.split(".")[-1])
    human_name = getattr(mod, "MARKATA_PLUGIN_NAME", pkg_name.replace("-", " ").title())

    # Guess config root from package name
    config_root = pkg_name.replace("-", "_")
    if not hasattr(markata.config, config_root):
        config_root = None

    return PluginInfo(
        name=human_name,
        module=module_name,
        source=source,
        help=help_text,
        config_root=config_root,
    )


def _plugin_info_from_file(path: Path, source: str, markata: "Markata") -> PluginInfo:
    text = path.read_text(encoding="utf8")
    try:
        module_ast = ast.parse(text)
        doc = ast.get_docstring(module_ast) or ""
    except SyntaxError:
        doc = ""
    help_text = _first_paragraph(doc)

    stem = path.stem
    pkg_name = stem
    human_name = stem.replace("_", " ").replace("-", " ").title()

    config_root = stem
    if not hasattr(markata.config, config_root):
        config_root = None

    return PluginInfo(
        name=human_name,
        module=str(path),
        source=source,
        help=help_text,
        config_root=config_root,
    )


def discover_plugins(markata: "Markata", settings: ConfigEditorSettings) -> List[PluginInfo]:
    """Discover Markata plugins and local plugins.

    - Core: all modules under `markata.plugins.*`
    - Local: any `*.py` file inside configured `plugin_dirs`.
    """

    plugins: List[PluginInfo] = []

    # Core markata.plugins.*
    try:
        import markata.plugins as m_plugins

        prefix = m_plugins.__name__ + "."
        for info in pkgutil.iter_modules(m_plugins.__path__, prefix):  # type: ignore[attr-defined]
            module_name = info.name
            try:
                plugins.append(_plugin_info_from_module(module_name, "core", markata))
            except Exception:
                # Don't let a broken plugin kill the UI
                continue
    except Exception:
        pass

    # Local plugin files
    for plugin_dir in settings.plugin_dirs:
        if not plugin_dir.exists() or not plugin_dir.is_dir():
            continue
        for path in sorted(plugin_dir.glob("*.py")):
            try:
                plugins.append(_plugin_info_from_file(path, "local", markata))
            except Exception:
                continue

    # Deduplicate by (module, source, config_root)
    seen: set[tuple[str, str, Optional[str]]] = set()
    unique: List[PluginInfo] = []
    for p in plugins:
        key = (p.module, p.source, p.config_root)
        if key in seen:
            continue
        seen.add(key)
        unique.append(p)

    # Sort for display
    unique.sort(key=lambda p: (p.source, (p.config_root or "zzz"), p.name))
    return unique


# ----------------- FASTAPI APP -----------------


def create_app(markata: "Markata") -> FastAPI:
    """Build a FastAPI app bound to the given Markata instance."""

    app = FastAPI(title="Markata Config Editor")

    settings: ConfigEditorSettings = getattr(markata.config, "config_editor")
    plugin_infos = discover_plugins(markata, settings)

    @app.get("/", response_class=HTMLResponse)
    def index() -> HTMLResponse:  # type: ignore[override]
        items = sorted(iter_config_items(markata), key=lambda i: i.path)

        # Build plugin filter options (only those with config_root)
        plugin_filter_options: list[str] = []
        for p in plugin_infos:
            if not p.config_root:
                continue
            plugin_filter_options.append(
                f"<option value='{p.config_root}'>{p.name} ({p.config_root})</option>"
            )
        plugin_filter_html = "\n".join(plugin_filter_options)

        # Plugin info table rows
        plugin_rows: list[str] = []
        for p in plugin_infos:
            cfg_root_display = p.config_root or "(no config root)"
            help_html = p.help or ""
            plugin_rows.append(
                "<tr class='border-b border-slate-800 hover:bg-slate-900/60'>"
                f"<td class='px-3 py-2 align-top text-xs font-mono text-slate-100'>{p.name}</td>"
                f"<td class='px-3 py-2 align-top text-xs font-mono text-slate-300'>{p.module}</td>"
                f"<td class='px-3 py-2 align-top text-xs text-slate-400'>{p.source}</td>"
                f"<td class='px-3 py-2 align-top text-xs font-mono text-slate-300'>{cfg_root_display}</td>"
                f"<td class='px-3 py-2 align-top text-xs text-slate-200'>{help_html}</td>"
                "</tr>"
            )
        plugin_table_body = "\n".join(plugin_rows)

        # Config table rows
        config_rows: list[str] = []
        for item in items:
            # Value preview
            if isinstance(item.value, (dict, list)):
                value_display = json.dumps(item.value, indent=2, default=str)
                input_html = (
                    f"<textarea name='value' rows='4' "
                    "class='w-full font-mono text-xs rounded-md border border-slate-700 "
                    "bg-slate-900 px-2 py-1 text-slate-100 focus:outline-none focus:ring-2 "
                    "focus:ring-pink-500/70'>"
                    f"{value_display}</textarea>"
                )
            else:
                value_display = "" if item.value is None else str(item.value)
                escaped = (
                    value_display
                    .replace("&", "&amp;")
                    .replace("<", "&lt;")
                    .replace(">", "&gt;")
                    .replace("\"", "&quot;")
                )
                input_html = (
                    "<input name='value' "
                    "class='w-full font-mono text-xs rounded-md border border-slate-700 "
                    "bg-slate-900 px-2 py-1 text-slate-100 focus:outline-none focus:ring-2 "
                    "focus:ring-pink-500/70' "
                    f"value='{escaped}' />"
                )

            desc_html = item.description or ""
            plugin_root = item.path.split(".", 1)[0]
            config_rows.append(
                f"<tr data-path='{item.path}' data-plugin-root='{plugin_root}' "
                "class='border-b border-slate-800 hover:bg-slate-900/60'>"
                f"<td class='px-3 py-2 align-top text-[0.70rem] font-mono text-pink-200 w-64 break-all'>{item.path}</td>"
                f"<td class='px-3 py-2 align-top text-[0.70rem] text-slate-400 w-28'>{item.type_name}</td>"
                f"<td class='hidden md:table-cell px-3 py-2 align-top text-[0.70rem] text-slate-300 w-80'>{desc_html}</td>"
                "<td class='px-3 py-2 align-top text-[0.70rem] w-full'>"
                "<form method='post' action='/update' class='flex flex-col gap-1'>"
                f"<input type='hidden' name='path' value='{item.path}' />"
                f"{input_html}"
                "<div class='flex justify-end'>"
                "<button type='submit' "
                "class='inline-flex items-center rounded-md border border-slate-600 "
                "bg-slate-800 px-2 py-0.5 text-[0.70rem] font-medium text-slate-100 "
                "hover:bg-slate-700 hover:border-pink-500/70 transition'>Save</button>"
                "</div>"
                "</form>"
                "</td>"
                "</tr>"
            )

        config_table_body = "\n".join(config_rows)

        html = f"""
<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <title>Markata Config Editor</title>
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <script src="https://cdn.tailwindcss.com"></script>
  </head>
  <body class="bg-slate-950 text-slate-100">
    <div class="min-h-screen max-w-6xl mx-auto px-4 py-6 space-y-6">
      <header class="space-y-1">
        <h1 class="text-2xl font-semibold tracking-tight">Markata Config Editor</h1>
        <p class="text-sm text-slate-400">
          Edit <code class="font-mono text-xs text-pink-300">markata.toml</code> values and explore installed plugins.
          Lists and dicts use JSON syntax.
        </p>
      </header>

      <section class="space-y-2">
        <h2 class="text-lg font-medium tracking-tight">Plugins</h2>
        <p class="text-xs text-slate-400">
          Core plugins live in <code class="font-mono text-[0.70rem] text-pink-300">markata.plugins.*</code>. Local plugins
          come from your configured <code class="font-mono text-[0.70rem] text-pink-300">plugin_dirs</code>.
        </p>

        <div class="overflow-hidden rounded-xl border border-slate-800 bg-slate-900/60">
          <div class="overflow-x-auto">
            <table class="min-w-full table-fixed text-xs">
              <thead class="bg-slate-900/80">
                <tr>
                  <th class="px-3 py-2 text-left font-medium text-slate-400 w-40">Name</th>
                  <th class="px-3 py-2 text-left font-medium text-slate-400 w-64">Module / Path</th>
                  <th class="px-3 py-2 text-left font-medium text-slate-400 w-20">Source</th>
                  <th class="px-3 py-2 text-left font-medium text-slate-400 w-32">Config Root</th>
                  <th class="px-3 py-2 text-left font-medium text-slate-400">Help</th>
                </tr>
              </thead>
              <tbody>
                {plugin_table_body}
              </tbody>
            </table>
          </div>
        </div>
      </section>

      <section class="space-y-3">
        <div class="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
          <div>
            <h2 class="text-lg font-medium tracking-tight">Configuration</h2>
            <p class="text-xs text-slate-400">
              Filter options by plugin and edit values inline. Changes are written back to
              <code class="font-mono text-[0.70rem] text-pink-300">markata.toml</code>.
            </p>
          </div>
          <div class="flex items-center gap-2 text-xs text-slate-300">
            <label class="flex items-center gap-2">
              <span class="whitespace-nowrap">Plugin filter:</span>
              <select id="pluginFilter"
                class="rounded-md border border-slate-700 bg-slate-900 px-2 py-1 text-xs text-slate-100
                       focus:outline-none focus:ring-2 focus:ring-pink-500/70">
                <option value="__all__">All config options</option>
                {plugin_filter_html}
              </select>
            </label>
          </div>
        </div>

        <div class="overflow-hidden rounded-xl border border-slate-800 bg-slate-900/60">
          <div class="overflow-x-auto">
            <table id="configTable" class="min-w-full table-fixed text-xs">
              <thead class="bg-slate-900/80">
                <tr>
                  <th class="px-3 py-2 text-left font-medium text-slate-400 w-64">Path</th>
                  <th class="px-3 py-2 text-left font-medium text-slate-400 w-28">Type</th>
                  <th class="hidden md:table-cell px-3 py-2 text-left font-medium text-slate-400 w-80">Description</th>
                  <th class="px-3 py-2 text-left font-medium text-slate-400">Value</th>
                </tr>
              </thead>
              <tbody>
                {config_table_body}
              </tbody>
            </table>
          </div>
        </div>
      </section>
    </div>

    <script>
      const filterSelect = document.getElementById('pluginFilter');
      const rows = Array.from(document.querySelectorAll('#configTable tbody tr'));

      function applyFilter() {{
        const value = filterSelect.value;
        rows.forEach((row) => {{
          const root = row.getAttribute('data-plugin-root');
          if (value === '__all__' || root === value) {{
            row.style.display = '';
          }} else {{
            row.style.display = 'none';
          }}
        }});
      }}

      filterSelect.addEventListener('change', applyFilter);
    </script>
  </body>
</html>
"""
        return HTMLResponse(content=html)

    @app.post("/update")
    def update(  # type: ignore[override]
        path: str = Form(...),
        value: str = Form(...),
    ):
        item = get_item_by_path(markata, path)
        if item is None:
            return JSONResponse(
                {"ok": False, "error": f"Unknown config path {path!r}"},
                status_code=400,
            )

        settings_local: ConfigEditorSettings = getattr(markata.config, "config_editor")
        config_path = settings_local.config_path

        try:
            coerced = coerce_value(value, item.field_type)
        except Exception as e:  # pragma: no cover - user input errors
            return JSONResponse(
                {"ok": False, "error": f"Failed to parse value: {e}"},
                status_code=400,
            )

        # Update in-memory Pydantic config
        segments = path.split(".")
        current_obj: Any = markata.config
        for seg in segments[:-1]:
            current_obj = getattr(current_obj, seg)
        setattr(current_obj, segments[-1], coerced)

        # Update TOML file
        update_toml_value(config_path, item.toml_path, coerced)

        return RedirectResponse("/", status_code=303)

    return app

    @app.post("/update")
    def update(  # type: ignore[override]
        path: str = Form(...),
        value: str = Form(...),
    ):
        item = get_item_by_path(markata, path)
        if item is None:
            return JSONResponse(
                {"ok": False, "error": f"Unknown config path {path!r}"},
                status_code=400,
            )

        settings_local: ConfigEditorSettings = getattr(markata.config, "config_editor")
        config_path = settings_local.config_path

        try:
            coerced = coerce_value(value, item.field_type)
        except Exception as e:  # pragma: no cover - user input errors
            return JSONResponse(
                {"ok": False, "error": f"Failed to parse value: {e}"},
                status_code=400,
            )

        # Update in-memory Pydantic config
        segments = path.split(".")
        current_obj: Any = markata.config
        for seg in segments[:-1]:
            current_obj = getattr(current_obj, seg)
        setattr(current_obj, segments[-1], coerced)

        # Update TOML file
        update_toml_value(config_path, item.toml_path, coerced)

        return RedirectResponse("/", status_code=303)

    return app


# ----------------- CLI HOOK -----------------


@hook_impl()
def cli(app, markata: "Markata") -> None:
    """Add a `config-web` command to the markata CLI that runs the config editor UI."""

    @app.command("config-web")
    def config_web(  # type: ignore[override]
        host: Optional[str] = typer.Option(
            None, help="Host to bind (overrides [markata.config_editor].host)",
        ),
        port: Optional[int] = typer.Option(
            None, help="Port to bind (overrides [markata.config_editor].port)",
        ),
        open_browser: Optional[bool] = typer.Option(
            None,
            "--open-browser/--no-open-browser",
            help="Open browser automatically (overrides config_editor.open_browser)",
        ),
    ) -> None:
        """Launch the web-based `markata.toml` config editor and plugin browser."""

        settings: ConfigEditorSettings = getattr(markata.config, "config_editor")
        host_ = host or settings.host
        port_ = port or settings.port
        if open_browser is None:
            open_browser_ = settings.open_browser
        else:
            open_browser_ = open_browser

        url = f"http://{host_}:{port_}/"
        markata.console.print(f"[bold]Starting config editor at[/bold] {url}")

        if open_browser_:
            try:
                webbrowser.open(url)
            except Exception:  # pragma: no cover - best-effort
                markata.console.log("Failed to open browser automatically")

        api_app = create_app(markata)
        uvicorn.run(api_app, host=host_, port=port_)

